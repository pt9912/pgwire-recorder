package bootstrap_test

import (
	"bytes"
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
)

// stummerUpstream nimmt TCP-Verbindungen an und antwortet nie; angenommen
// erhält je angenommener Verbindung einen Wert. Der Cleanup schließt Listener
// und Verbindungen.
func stummerUpstream(t *testing.T) (string, chan struct{}) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	angenommen := make(chan struct{}, 16)
	var mu sync.Mutex
	var conns []net.Conn
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			conns = append(conns, c)
			mu.Unlock()
			angenommen <- struct{}{}
		}
	}()
	t.Cleanup(func() {
		l.Close()
		mu.Lock()
		defer mu.Unlock()
		for _, c := range conns {
			c.Close()
		}
	})
	return l.Addr().String(), angenommen
}

// recordLauf ist ein laufender Aufruf von record über Run.
type recordLauf struct {
	cancel context.CancelFunc
	ablauf chan struct{}
	ende   chan int
	stderr *syncPuffer
}

// starteRecord ruft Run mit record und den Zusatzargumenten nebenher auf.
// Mit verbinden baut ein Client eine Verbindung auf und sendet die
// Startnachricht, und starteRecord wartet höchstens 5 s, bis der Recorder den
// Upstream erreicht hat, der nie antwortet: Die Verbindung steht dann im
// Verbindungsaufbau zum Upstream.
func starteRecord(t *testing.T, verbinden bool, args ...string) *recordLauf {
	t.Helper()
	t.Setenv("PGWIRE_RECORDER_LOG_LEVEL", "")
	t.Setenv("PGWIRE_RECORDER_SHUTDOWN_TIMEOUT", "")
	upstream, angenommen := stummerUpstream(t)
	listen := freieAdresse(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	r := &recordLauf{cancel: cancel, ablauf: make(chan struct{}), ende: make(chan int, 1), stderr: &syncPuffer{}}
	aufruf := append([]string{"record", "--listen", listen, "--upstream", upstream, "--output", filepath.Join(t.TempDir(), "rec.yaml")}, args...)
	go func() {
		var stdout bytes.Buffer
		r.ende <- bootstrap.Run(ctx, r.ablauf, aufruf, "dev", &stdout, r.stderr)
	}()
	var client net.Conn
	for deadline := time.Now().Add(5 * time.Second); ; {
		c, err := net.Dial("tcp", listen)
		if err == nil {
			client = c
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("record lauscht nicht binnen 5 s")
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Cleanup(func() { client.Close() })
	if !verbinden {
		client.Close()
		return r
	}
	if _, err := client.Write(startnachricht()); err != nil {
		t.Fatal(err)
	}
	select {
	case <-angenommen:
	case <-time.After(5 * time.Second):
		t.Fatal("record erreicht den Upstream nicht binnen 5 s")
	}
	return r
}

// endetBinnen wartet höchstens frist auf das Ende von Run und liefert den
// Exit-Code.
func (r *recordLauf) endetBinnen(t *testing.T, frist time.Duration, was string) int {
	t.Helper()
	select {
	case code := <-r.ende:
		return code
	case <-time.After(frist):
		t.Fatalf("%s: Run endet nicht binnen %v\n%s", was, frist, r.stderr.String())
		return 0
	}
}

// laeuftNoch prüft, dass Run binnen dauer nicht endet.
func (r *recordLauf) laeuftNoch(t *testing.T, dauer time.Duration, was string) {
	t.Helper()
	select {
	case code := <-r.ende:
		t.Fatalf("%s: Run endete mit Exit-Code %d\n%s", was, code, r.stderr.String())
	case <-time.After(dauer):
	}
}

// Abdeckung: LH-FA-13/Boundary — nach dem ersten Signal wartet record
// höchstens --shutdown-timeout auf eine Verbindung, die im Aufbau zum Upstream
// steht; danach endet sie zwangsweise ohne Meldung, und der Lauf endet mit
// Exit-Code 0. Die Zeile der Stufe info beim Beginn des Herunterfahrens nennt
// sessions=1 (LH-FA-13.a).
func TestRunRecordFristLaeuftAb(t *testing.T) {
	r := starteRecord(t, true, "--shutdown-timeout=200ms")
	r.cancel()
	code := r.endetBinnen(t, 3*time.Second, "Frist 200ms")
	text := r.stderr.String()
	if code != 0 || !strings.Contains(text, "level=INFO msg=\"Herunterfahren begonnen\" sessions=1") || strings.Contains(text, "level=ERROR") {
		t.Fatalf("Exit-Code %d, stderr\n%s", code, text)
	}
}

// Abdeckung: LH-FA-13/Boundary — mit --shutdown-timeout=0 wartet record nach
// dem ersten Signal ohne Frist; das zweite Signal (ablauf) lässt die Frist
// sofort ablaufen, auch beim Wert 0, und der Lauf endet (LH-FA-13.a *Zweites
// Signal*).
func TestRunRecordZweitesSignalOhneFrist(t *testing.T) {
	r := starteRecord(t, true, "--shutdown-timeout=0")
	r.cancel()
	r.laeuftNoch(t, time.Second, "Frist 0 nach dem ersten Signal")
	close(r.ablauf)
	if code := r.endetBinnen(t, 3*time.Second, "zweites Signal"); code != 0 {
		t.Fatalf("Exit-Code %d\n%s", code, r.stderr.String())
	}
}

// Abdeckung: LH-FA-13/Boundary — das zweite Signal lässt auch eine lange Frist
// sofort ablaufen (LH-FA-13.a *Zweites Signal*).
func TestRunRecordZweitesSignalLangeFrist(t *testing.T) {
	r := starteRecord(t, true, "--shutdown-timeout=60s")
	r.cancel()
	close(r.ablauf)
	if code := r.endetBinnen(t, 3*time.Second, "zweites Signal bei Frist 60s"); code != 0 {
		t.Fatalf("Exit-Code %d\n%s", code, r.stderr.String())
	}
}

// Abdeckung: LH-FA-13/Boundary — ist beim Signal keine Verbindung offen, endet
// record ohne zu warten, auch bei langer Frist, schreibt die Aufzeichnung und
// nennt sessions=0 (LH-FA-13.a *Ohne offene Verbindung*).
func TestRunRecordOhneOffeneVerbindung(t *testing.T) {
	r := starteRecord(t, false, "--shutdown-timeout=60s")
	time.Sleep(100 * time.Millisecond)
	r.cancel()
	code := r.endetBinnen(t, 2*time.Second, "ohne offene Verbindung")
	if text := r.stderr.String(); code != 0 || !strings.Contains(text, "sessions=0") || !strings.Contains(text, "record beendet") {
		t.Fatalf("Exit-Code %d, stderr\n%s", code, text)
	}
}

// Abdeckung: LH-FA-13/Negative — ein Signal in der Startphase bricht die
// Startprüfungen nicht ab: Ein Startfehler (vorhandenes --output ohne --force,
// PGR-E2002) beendet record mit Exit-Code 2; ohne Startfehler endet record wie
// ohne offene Verbindung mit Exit-Code 0 und sessions=0 (LH-FA-13.a
// *Startphase*).
func TestRunRecordSignalInDerStartphase(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_LOG_LEVEL", "")
	t.Setenv("PGWIRE_RECORDER_SHUTDOWN_TIMEOUT", "")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(output, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := bootstrap.Run(ctx, nil, []string{"record", "--listen", "127.0.0.1:0", "--upstream", "127.0.0.1:1", "--output", output}, "dev", &stdout, &stderr); code != 2 || !strings.Contains(stderr.String(), "PGR-E2002") {
		t.Fatalf("Startfehler: Exit-Code %d, stderr %q", code, stderr.String())
	}
	stderr.Reset()
	neu := filepath.Join(t.TempDir(), "neu.yaml")
	if code := bootstrap.Run(ctx, nil, []string{"record", "--listen", "127.0.0.1:0", "--upstream", "127.0.0.1:1", "--output", neu, "--shutdown-timeout=60s"}, "dev", &stdout, &stderr); code != 0 || !strings.Contains(stderr.String(), "sessions=0") {
		t.Fatalf("ohne Startfehler: Exit-Code %d, stderr %q", code, stderr.String())
	}
	if _, err := os.Stat(neu); err != nil {
		t.Fatalf("Aufzeichnung nicht geschrieben: %v", err)
	}
}
