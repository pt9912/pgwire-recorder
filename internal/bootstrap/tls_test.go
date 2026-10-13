package bootstrap_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// tlsStartnachricht ist eine StartupMessage der Version 3.0 mit dem Benutzer app.
func tlsStartnachricht() []byte {
	rumpf := []byte("user\x00app\x00\x00")
	b := make([]byte, 8, 8+len(rumpf))
	binary.BigEndian.PutUint32(b[0:4], uint32(8+len(rumpf)))
	binary.BigEndian.PutUint32(b[4:8], 196608)
	return append(b, rumpf...)
}

// Abdeckung: LH-FA-23/Happy, LH-FA-23/Negative — replay reicht Zertifikat, Schlüssel
// und --allow-plaintext an den Server durch: Ein Client mit SSLRequest erhält S und
// einen Handshake, einer ohne SSLRequest den Aufbau (mit --allow-plaintext) oder
// PGR-E6003 mit Exit-Code 6 (ohne); ohne --tls-cert antwortet der Server mit N
// (LH-FA-23.a).
func TestRunReplayTLS(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	blatt := testpki.NeueCA(t, "CA").Server(t, testpki.Blatt{})
	dir := t.TempDir()
	cert, key := filepath.Join(dir, "c.pem"), filepath.Join(dir, "k.pem")
	for pfad, inhalt := range map[string][]byte{cert: blatt.PEM(), key: blatt.SchluesselPEM(t, "PRIVATE KEY")} {
		if err := os.WriteFile(pfad, inhalt, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []struct {
		name      string
		args      []string
		antwort   byte // Antwort auf den SSLRequest
		klartext  byte // erste Antwort auf eine StartupMessage ohne SSLRequest
		exitCode  int
		verbindet bool
	}{
		{"TLS", []string{"--tls-cert=" + cert, "--tls-key=" + key}, 'S', 'E', 6, true},
		{"TLS mit Klartext", []string{"--tls-cert=" + cert, "--tls-key=" + key, "--allow-plaintext"}, 'S', 'R', 0, true},
		{"ohne TLS", nil, 'N', 'R', 0, false},
	} {
		t.Run(f.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			listen := freieAdresse(t)
			var stderr bytes.Buffer
			ende := make(chan int, 1)
			args := append([]string{"replay", "--listen", listen, "--input", aufzeichnung(t)}, f.args...)
			go func() { ende <- bootstrap.Run(ctx, nil, args, "dev", io.Discard, &stderr) }()
			dial := func() net.Conn {
				var c net.Conn
				var err error
				for i := 0; i < 100; i++ {
					if c, err = net.Dial("tcp", listen); err == nil {
						_ = c.SetDeadline(time.Now().Add(10 * time.Second))
						return c
					}
					time.Sleep(50 * time.Millisecond)
				}
				t.Fatal(err)
				return nil
			}
			c := dial()
			ssl := make([]byte, 8)
			binary.BigEndian.PutUint32(ssl[0:4], 8)
			binary.BigEndian.PutUint32(ssl[4:8], 80877103)
			_, _ = c.Write(ssl)
			antwort := make([]byte, 1)
			if _, err := io.ReadFull(c, antwort); err != nil || antwort[0] != f.antwort {
				t.Fatalf("Antwort auf SSLRequest %q, %v, erwartet %q", antwort, err, f.antwort)
			}
			if f.verbindet {
				if err := tls.Client(c, &tls.Config{InsecureSkipVerify: true}).Handshake(); err != nil {
					t.Fatalf("Handshake: %v", err)
				}
			}
			c.Close()
			c = dial()
			defer c.Close()
			_, _ = c.Write(tlsStartnachricht())
			if _, err := io.ReadFull(c, antwort); err != nil || antwort[0] != f.klartext {
				t.Fatalf("erste Antwort ohne SSLRequest %q, %v, erwartet %q", antwort, err, f.klartext)
			}
			c.Close()
			cancel()
			if code := <-ende; code != f.exitCode {
				t.Fatalf("Exit-Code %d, erwartet %d, stderr:\n%s", code, f.exitCode, stderr.String())
			}
			if f.exitCode == 6 && !strings.Contains(stderr.String(), "PGR-E6003") {
				t.Fatalf("stderr ohne PGR-E6003:\n%s", stderr.String())
			}
		})
	}
}

// Abdeckung: LH-FA-23/Negative, LH-RB-01/Messung — mit --log-level=debug erscheint der
// Schlüssel in keiner Zeile: nicht beim Start, nicht bei einem gelungenen und einem
// gescheiterten Handshake, nicht bei einem abgewiesenen Klartext-Client, und ein
// Start mit nicht passendem Schlüssel nennt ihn nicht (LH-FA-23.a *Dateien*).
func TestRunTLSLogOhneSchluessel(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	ca := testpki.NeueCA(t, "CA")
	blatt, anderes := ca.Server(t, testpki.Blatt{}), ca.Server(t, testpki.Blatt{})
	dir := t.TempDir()
	dateien := map[string][]byte{"c.pem": blatt.PEM(), "k.pem": blatt.SchluesselPEM(t, "PRIVATE KEY"), "x.pem": anderes.SchluesselPEM(t, "PRIVATE KEY")}
	for name, inhalt := range dateien {
		if err := os.WriteFile(filepath.Join(dir, name), inhalt, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	cert, key, fremd := filepath.Join(dir, "c.pem"), filepath.Join(dir, "k.pem"), filepath.Join(dir, "x.pem")
	pruefe := func(was, log string) {
		t.Helper()
		zeilen := strings.Split(string(dateien["k.pem"]), "\n")
		for _, geheim := range []string{"PRIVATE KEY", zeilen[1], string(dateien["x.pem"])[40:70]} {
			if strings.Contains(log, geheim) {
				t.Errorf("%s: das Log nennt Schlüsselinhalt %q:\n%s", was, geheim, log)
			}
		}
	}
	var stderr bytes.Buffer
	if code := bootstrap.Run(context.Background(), nil, []string{"replay", "--listen", "127.0.0.1:0", "--input", aufzeichnung(t), "--tls-cert", cert, "--tls-key", fremd, "--log-level=debug"}, "dev", io.Discard, &stderr); code != 2 || !strings.Contains(stderr.String(), "PGR-E2007") {
		t.Fatalf("nicht passender Schlüssel: Exit-Code %d, stderr:\n%s", code, stderr.String())
	}
	pruefe("Start mit nicht passendem Schlüssel", stderr.String())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listen := freieAdresse(t)
	stderr.Reset()
	ende := make(chan int, 1)
	args := []string{"replay", "--listen", listen, "--input", aufzeichnung(t), "--tls-cert", cert, "--tls-key", key, "--log-level=debug"}
	go func() { ende <- bootstrap.Run(ctx, nil, args, "dev", io.Discard, &stderr) }()
	dial := func() net.Conn {
		for i := 0; i < 100; i++ {
			if c, err := net.Dial("tcp", listen); err == nil {
				_ = c.SetDeadline(time.Now().Add(10 * time.Second))
				return c
			}
			time.Sleep(50 * time.Millisecond)
		}
		t.Fatal("Server lauscht nicht")
		return nil
	}
	ssl := []byte{0, 0, 0, 8, 0x04, 0xd2, 0x16, 0x2f}
	c := dial()
	_, _ = c.Write(ssl)
	_, _ = io.ReadFull(c, make([]byte, 1))
	if err := tls.Client(c, &tls.Config{InsecureSkipVerify: true}).Handshake(); err != nil {
		t.Fatal(err)
	}
	c.Close()
	c = dial()
	_, _ = c.Write(ssl)
	_, _ = io.ReadFull(c, make([]byte, 1))
	c.Close()
	c = dial()
	_, _ = c.Write(tlsStartnachricht())
	_, _ = io.Copy(io.Discard, c)
	c.Close()
	time.Sleep(200 * time.Millisecond)
	cancel()
	<-ende
	log := stderr.String()
	if !strings.Contains(log, "PGR-W3002") || !strings.Contains(log, "PGR-E6003") {
		t.Fatalf("Log ohne die erwarteten Zeilen:\n%s", log)
	}
	pruefe("Lauf mit Handshakes", log)
}
