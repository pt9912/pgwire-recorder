package bootstrap_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Verhalten des Fake-Servers nach der Antwort mit dem Passwort.
const (
	anmAntworten = "antworten" // Passwort prüfen, bei Erfolg Aufbau abschließen und Anfragen beantworten
	anmWarten    = "warten"    // wie antworten, aber erst nach frei
	anmSchweigen = "schweigen" // nichts mehr senden, bis die Verbindung endet
)

// anmeldeServer fordert von jedem Client das Klartext-Passwort an und meldet
// das gelesene Passwort auf passwoerter. Ein anderes als erwartet beantwortet
// er mit einer Fehlerantwort der Klasse 28; sonst verhält er sich nach
// verhalten. Was der Client nach dem Passwort sendet, bis die Verbindung
// endet, steht als Typen der Nachrichten auf danach; sendet der Client kein
// Passwort, steht dort "".
func anmeldeServer(t *testing.T, erwartet, verhalten string) (addr string, passwoerter <-chan string, frei chan struct{}, danach <-chan string) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	pw := make(chan string, 8)
	frei = make(chan struct{})
	typen := make(chan string, 8)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
				var laenge [4]byte
				if _, err := io.ReadFull(conn, laenge[:]); err != nil {
					return
				}
				if _, err := io.ReadFull(conn, make([]byte, binary.BigEndian.Uint32(laenge[:])-4)); err != nil {
					return
				}
				_, _ = conn.Write(nachricht('R', 0, 0, 0, 3))
				var kopf [5]byte
				if _, err := io.ReadFull(conn, kopf[:]); err != nil {
					typen <- ""
					return
				}
				rumpf := make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)
				if _, err := io.ReadFull(conn, rumpf); err != nil {
					typen <- ""
					return
				}
				gesendet := strings.TrimSuffix(string(rumpf), "\x00")
				pw <- gesendet
				switch {
				case gesendet != erwartet:
					_, _ = conn.Write(fehlerantwortRoh("FATAL", "28P01", "Passwort falsch"))
					return
				case verhalten == anmSchweigen:
				default:
					if verhalten == anmWarten {
						select {
						case <-frei:
						case <-time.After(30 * time.Second):
							return
						}
					}
					_, _ = conn.Write(bereitOhneAnmeldung())
				}
				var gelesen []byte
				for {
					if _, err := io.ReadFull(conn, kopf[:]); err != nil {
						typen <- string(gelesen)
						return
					}
					if _, err := io.ReadFull(conn, make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)); err != nil {
						typen <- string(gelesen)
						return
					}
					gelesen = append(gelesen, kopf[0])
					if kopf[0] == 'Q' {
						_, _ = conn.Write(append(nachricht('C', []byte("SELECT 0\x00")...), nachricht('Z', 'I')...))
					}
				}
			}()
		}
	}()
	return l.Addr().String(), pw, frei, typen
}

// anmeldePasswort wartet höchstens 30 s auf das Passwort, das der Fake-Server
// gelesen hat.
func anmeldePasswort(t *testing.T, pw <-chan string) string {
	t.Helper()
	select {
	case s := <-pw:
		return s
	case <-time.After(30 * time.Second):
		t.Fatal("der Fake-Server liest binnen 30 s kein Passwort")
		return ""
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Negative — play
// meldet sich mit dem Passwort aus PGWIRE_RECORDER_PASSWORD an, bei host:port,
// und mit dem Passwort aus dem Platzhalter der benutzten Verbindung, das der
// Variable vorgeht; ein falsches Passwort ist PGR-E4005 mit Exit-Code 4, der
// SQLSTATE steht in der Zeile, und weder das richtige noch das falsche
// Passwort steht in stdout oder stderr (LH-FA-20.a *Passwort*, *Anmeldung*).
func TestRunPlayPasswort(t *testing.T) {
	leerePlay(t)
	addr, passwoerter, _, danach := anmeldeServer(t, "GEHEIMPW", anmAntworten)
	input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A")})
	host, port, _ := net.SplitHostPort(addr)
	datei := filepath.Join(t.TempDir(), "k.yaml")
	if err := os.WriteFile(datei, []byte(fmt.Sprintf("connections:\n  v: \"postgresql://u:${PGR_BT_PW}@%s:%s/db\"\n", host, port)), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, f := range []struct {
		name string
		env  string
		pw   string
		args []string
		code int
		want string
	}{
		{"Variable bei host:port", "GEHEIMPW", "", []string{"--upstream", addr}, 0, "GEHEIMPW"},
		{"Platzhalter vor der Variable", "GEHEIMFALSCH", "GEHEIMPW", []string{"--upstream", "v", "--config", datei}, 0, "GEHEIMPW"},
		{"falsche Variable", "GEHEIMFALSCH", "", []string{"--upstream", addr}, 4, "GEHEIMFALSCH"},
		{"falscher Platzhalter", "GEHEIMPW", "GEHEIMFALSCH", []string{"--upstream", "v", "--config", datei}, 4, "GEHEIMFALSCH"},
	} {
		t.Run(f.name, func(t *testing.T) {
			t.Setenv("PGWIRE_RECORDER_PASSWORD", f.env)
			t.Setenv("PGR_BT_PW", f.pw)
			var stdout, stderr bytes.Buffer
			code := bootstrap.Run(context.Background(), nil, append([]string{"play", "--input", input}, f.args...), "dev", &stdout, &stderr)
			if got := anmeldePasswort(t, passwoerter); got != f.want {
				t.Fatalf("der Server liest das Passwort %q, erwartet %q", got, f.want)
			}
			if code != f.code || strings.Contains(stdout.String()+stderr.String(), "GEHEIM") {
				t.Fatalf("Exit-Code %d, erwartet %d; Ausgabe ohne GEHEIM erwartet, stdout %q, stderr %q", code, f.code, stdout.String(), stderr.String())
			}
			if f.code == 4 && (!strings.Contains(stderr.String(), "PGR-E4005") || !strings.Contains(stderr.String(), "28P01")) {
				t.Fatalf("stderr ohne PGR-E4005 und SQLSTATE: %q", stderr.String())
			}
			if f.code == 0 {
				if got := <-danach; got != "QX" {
					t.Fatalf("Client-Nachrichten %q, erwartet QX", got)
				}
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — ein Passwort, das der Server nicht verlangt,
// ist kein Fehler, und die Variable bleibt ohne Wirkung, wenn der Server kein
// Passwort fordert; fordert er eines und es gibt keines, ist das PGR-E4005 mit
// Exit-Code 4, und play sendet nichts (LH-FA-20.a *Passwort*, *Anmeldung*).
func TestRunPlayOhnePasswort(t *testing.T) {
	leerePlay(t)
	addr, _, _, danach := anmeldeServer(t, "x", anmAntworten)
	input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A")})
	var stdout, stderr bytes.Buffer
	if code := bootstrap.Run(context.Background(), nil, []string{"play", "--upstream", addr, "--input", input}, "dev", &stdout, &stderr); code != 4 || !strings.Contains(stderr.String(), "PGR-E4005") {
		t.Fatalf("Exit-Code %d, stderr %q", code, stderr.String())
	}
	select {
	case got := <-danach:
		if got != "" {
			t.Fatalf("play sendet %q", got)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("der Fake-Server sieht binnen 30 s kein Ende der Verbindung")
	}
	// Der erste Server verlangt ein Passwort; ein Server, der keines verlangt,
	// ist playServer, der AuthenticationOk sofort sendet.
	ohne, typen := playServer(t, append(nachricht('C', []byte("SELECT 0\x00")...), nachricht('Z', 'I')...))
	t.Setenv("PGWIRE_RECORDER_PASSWORD", "GEHEIMPW")
	if code := bootstrap.Run(context.Background(), nil, []string{"play", "--upstream", ohne, "--input", input}, "dev", &stdout, &stderr); code != 0 {
		t.Fatalf("Server ohne Passwort: Exit-Code %d, stderr %q", code, stderr.String())
	}
	if got := typenBis(t, typen, 'X'); got != "QX" {
		t.Fatalf("Client-Nachrichten %q, erwartet QX", got)
	}
}

// Abdeckung: LH-FA-20/Boundary — das erste Signal während der Anmeldung lässt
// den Aufbau zu Ende laufen: Wartet play auf die Antwort nach dem Passwort und
// kommt das Signal, steht die Verbindung danach, es beginnt keine Interaktion,
// Terminate geht an den Server, Exit-Code 0 (LH-FA-20.a *Abbruchsignal*).
func TestRunPlayErstesSignalInAnmeldung(t *testing.T) {
	leerePlay(t)
	t.Setenv("PGWIRE_RECORDER_PASSWORD", "GEHEIMPW")
	addr, passwoerter, frei, danach := anmeldeServer(t, "GEHEIMPW", anmWarten)
	ctx, erstes := context.WithCancel(context.Background())
	input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A")})
	fertig := make(chan int, 1)
	var stdout, stderr bytes.Buffer
	go func() {
		fertig <- bootstrap.Run(ctx, nil, []string{"play", "--upstream", addr, "--input", input}, "dev", &stdout, &stderr)
	}()
	anmeldePasswort(t, passwoerter)
	erstes()
	select {
	case code := <-fertig:
		t.Fatalf("play endet nach dem ersten Signal vor der Antwort des Servers: Exit-Code %d, stderr %q", code, stderr.String())
	case <-time.After(500 * time.Millisecond):
	}
	close(frei)
	select {
	case code := <-fertig:
		if code != 0 || strings.Contains(stderr.String(), "level=ERROR") {
			t.Fatalf("Exit-Code %d, stderr %q", code, stderr.String())
		}
	case <-time.After(30 * time.Second):
		t.Fatal("play endet binnen 30 s nach der Antwort des Servers nicht")
	}
	if got := <-danach; got != "X" {
		t.Fatalf("Client-Nachrichten %q, erwartet nur Terminate", got)
	}
}

// Abdeckung: LH-FA-20/Boundary — das zweite Signal beendet play während der
// Anmeldung sofort, wenn der Server nach dem Passwort schweigt: Die Verbindung
// wird ohne Terminate geschlossen, der abgebrochene Aufbau ist kein Fehler,
// Exit-Code 0 (LH-FA-20.a *Abbruchsignal*, *Abbruch im Aufbau*).
func TestRunPlayZweitesSignalInAnmeldung(t *testing.T) {
	leerePlay(t)
	t.Setenv("PGWIRE_RECORDER_PASSWORD", "GEHEIMPW")
	addr, passwoerter, _, danach := anmeldeServer(t, "GEHEIMPW", anmSchweigen)
	ctx, erstes := context.WithCancel(context.Background())
	ablauf := make(chan struct{})
	input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A")})
	fertig := make(chan int, 1)
	var stdout, stderr bytes.Buffer
	go func() {
		fertig <- bootstrap.Run(ctx, ablauf, []string{"play", "--upstream", addr, "--input", input}, "dev", &stdout, &stderr)
	}()
	anmeldePasswort(t, passwoerter)
	erstes()
	close(ablauf)
	select {
	case code := <-fertig:
		if code != 0 || strings.Contains(stderr.String(), "level=ERROR") {
			t.Fatalf("Exit-Code %d, stderr %q", code, stderr.String())
		}
	case <-time.After(30 * time.Second):
		t.Fatal("play endet binnen 30 s nach dem zweiten Signal nicht")
	}
	select {
	case got := <-danach:
		if got != "" {
			t.Fatalf("Client-Nachrichten %q nach dem Passwort, erwartet keine", got)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("der Fake-Server sieht binnen 30 s kein Ende der Verbindung")
	}
}
