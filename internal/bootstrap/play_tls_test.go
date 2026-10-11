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
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// tlsPlayServer nimmt Verbindungen an: Auf ein SSLRequest antwortet er `S`,
// handelt TLS mit blatt aus, liest das Startup, sendet bereitOhneAnmeldung und
// meldet jede Client-Nachricht mit ihrem Typ auf typen; jede Verbindung ohne
// SSLRequest meldet er mit dem Typ 'P' und schließt sie.
func tlsPlayServer(t *testing.T, blatt testpki.Blattzertifikat) (string, <-chan byte) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	typen := make(chan byte, 64)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
				kopf := make([]byte, 8)
				if _, err := io.ReadFull(conn, kopf); err != nil {
					return
				}
				if binary.BigEndian.Uint32(kopf[4:]) != 80877103 {
					typen <- 'P'
					return
				}
				_, _ = conn.Write([]byte("S"))
				tc := tls.Server(conn, &tls.Config{Certificates: []tls.Certificate{{Certificate: [][]byte{blatt.DER}, PrivateKey: blatt.Schluessel}}})
				if tc.Handshake() != nil {
					typen <- 'H'
					return
				}
				var laenge [4]byte
				if _, err := io.ReadFull(tc, laenge[:]); err != nil {
					return
				}
				if _, err := io.ReadFull(tc, make([]byte, binary.BigEndian.Uint32(laenge[:])-4)); err != nil {
					return
				}
				typen <- 'S'
				_, _ = tc.Write(bereitOhneAnmeldung())
				for {
					var k [5]byte
					if _, err := io.ReadFull(tc, k[:]); err != nil {
						return
					}
					if _, err := io.ReadFull(tc, make([]byte, binary.BigEndian.Uint32(k[1:])-4)); err != nil {
						return
					}
					typen <- k[0]
					if k[0] == 'Q' {
						_, _ = tc.Write(append(nachricht('C', []byte("SELECT 0\x00")...), nachricht('Z', 'I')...))
					}
				}
			}()
		}
	}()
	return l.Addr().String(), typen
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Negative — der Bootstrap
// reicht UpstreamTLS und die Zertifikate aus --upstream-ca an den Server
// weiter: Mit --upstream-tls und --upstream-ca spielt play über TLS ein und
// endet mit Exit-Code 0; mit --upstream-tls ohne --upstream-ca ist die
// Zertifizierungsstelle unbekannt (PGR-E4005, Exit-Code 4); ohne --upstream-tls
// verbindet play ohne SSLRequest (PGR-E4002, Exit-Code 4); die Meldungen nennen
// weder Pfad noch Zertifikat (LH-FA-20.a *TLS*, *Fehlerregeln beim Einspielen*).
func TestRunPlayTLS(t *testing.T) {
	ca := testpki.NeueCA(t, "Test-CA")
	caDatei := filepath.Join(t.TempDir(), "GEHEIMCA.pem")
	if err := os.WriteFile(caDatei, ca.PEM(), 0o600); err != nil {
		t.Fatal(err)
	}
	input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A")})
	for _, f := range []struct {
		name    string
		args    []string
		code    int
		meldung string
		typen   string
	}{
		{"TLS mit Zertifizierungsstelle", []string{"--upstream-tls", "--upstream-ca", caDatei}, 0, "", "SQX"},
		{"TLS ohne Zertifizierungsstelle", []string{"--upstream-tls"}, 4, "PGR-E4005", "H"},
		{"ohne TLS", nil, 4, "PGR-E4002", "P"},
	} {
		t.Run(f.name, func(t *testing.T) {
			leerePlay(t)
			addr, typen := tlsPlayServer(t, ca.Server(t, testpki.Blatt{IPs: []net.IP{net.ParseIP("127.0.0.1")}}))
			var stdout, stderr bytes.Buffer
			args := append([]string{"play", "--upstream", addr, "--input", input}, f.args...)
			code := bootstrap.Run(context.Background(), nil, args, "dev", &stdout, &stderr)
			if code != f.code || !strings.Contains(stderr.String(), f.meldung) || strings.Contains(stderr.String(), "GEHEIM") {
				t.Fatalf("Exit-Code %d, erwartet %d mit %q, stderr %q", code, f.code, f.meldung, stderr.String())
			}
			var got []byte
			for len(got) < len(f.typen) {
				select {
				case b := <-typen:
					got = append(got, b)
				case <-time.After(30 * time.Second):
					t.Fatalf("Client-Nachrichten %q, erwartet %q", got, f.typen)
				}
			}
			if string(got) != f.typen {
				t.Fatalf("Client-Nachrichten %q, erwartet %q", got, f.typen)
			}
		})
	}
}
