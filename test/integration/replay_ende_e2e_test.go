//go:build integration

package integration_test

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Die Aufzeichnungen sind von Hand geschrieben (Version 2, ohne PostgreSQL);
// jede gehört einem Fall.
//
//go:embed testdata/ende-fatal.yaml
var endeFatalYAML []byte

//go:embed testdata/ende-panic.yaml
var endePanicYAML []byte

//go:embed testdata/ende-lebend.yaml
var endeLebendYAML []byte

// endeDatei schreibt inhalt in eine Datei im Verzeichnis des Tests und liefert
// ihren Pfad.
func endeDatei(t *testing.T, inhalt []byte) string {
	t.Helper()
	pfad := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(pfad, inhalt, 0o644); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// serverNachricht ist eine Nachricht des Servers: Typ und Rumpf.
type serverNachricht struct {
	typ   byte
	rumpf []byte
}

// rohVerbindung verbindet sich ohne Treiber, baut die Verbindung auf und liest
// bis zum ersten ReadyForQuery.
func rohVerbindung(t *testing.T, listen string) (net.Conn, *bufio.Reader) {
	t.Helper()
	c, err := net.DialTimeout("tcp", listen, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	_ = c.SetDeadline(time.Now().Add(10 * time.Second))
	r := bufio.NewReader(c)
	start := []byte{0, 0, 0, 0, 0, 3, 0, 0}
	start = append(start, "user\x00postgres\x00database\x00postgres\x00\x00"...)
	binary.BigEndian.PutUint32(start, uint32(len(start)))
	schreibe(t, c, start)
	bisTyp(t, r, 'Z')
	return c, r
}

// anfrage sendet eine einfache Anfrage.
func anfrage(t *testing.T, c net.Conn, sql string) {
	t.Helper()
	q := append([]byte{'Q', 0, 0, 0, 0}, sql+"\x00"...)
	binary.BigEndian.PutUint32(q[1:], uint32(len(q)-1))
	schreibe(t, c, q)
}

// bisVerbindungsende liest die Nachrichten des Servers, bis er die Verbindung
// schließt; jedes andere Ende des Lesens ist ein Fehler des Tests.
func bisVerbindungsende(t *testing.T, r *bufio.Reader) []serverNachricht {
	t.Helper()
	var out []serverNachricht
	for {
		kopf := make([]byte, 5)
		if _, err := io.ReadFull(r, kopf); err != nil {
			if errors.Is(err, io.EOF) {
				return out
			}
			t.Fatalf("Verbindung endet nicht sauber: %v", err)
		}
		rumpf := make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)
		if _, err := io.ReadFull(r, rumpf); err != nil {
			t.Fatalf("Nachricht lesen: %v", err)
		}
		out = append(out, serverNachricht{kopf[0], rumpf})
	}
}

// typen nennt die Typen der Nachrichten als Text.
func typen(ns []serverNachricht) string {
	var b strings.Builder
	for _, n := range ns {
		b.WriteByte(n.typ)
	}
	return b.String()
}

// keineMeldung prüft, dass stderr weder eine Warnung noch einen Fehler trägt.
func keineMeldung(t *testing.T, was, stderr string) {
	t.Helper()
	if strings.Contains(stderr, "PGR-") {
		t.Fatalf("%s: Meldung im Log:\n%s", was, stderr)
	}
}

// Abdeckung: LH-FA-25/Happy, LH-FA-03/Boundary, LH-FA-12/Boundary — replay
// mit einer handgeschriebenen Aufzeichnung: Der Client erhält für die Anfrage
// mit Ende die aufgezeichnete Fehlerantwort (FATAL) und danach das Ende der
// Verbindung, ohne ReadyForQuery, und der Recorder meldet nichts; die nächste
// Verbindung erhält die nächste Session; mit --fail-on-unconsumed ist der
// Exit-Code 0, weil die Interaktion mit Ende verbraucht ist.
func TestE2EReplayEndeFatal(t *testing.T) {
	rep := startProzess(t, "replay", "--input", endeDatei(t, endeFatalYAML), "--fail-on-unconsumed")

	c, r := rohVerbindung(t, rep.listen)
	anfrage(t, c, "SELECT 1;")
	bisTyp(t, r, 'Z')
	anfrage(t, c, "SELECT 2;")
	ns := bisVerbindungsende(t, r)
	if typen(ns) != "E" || !bytes.Contains(ns[0].rumpf, []byte("SFATAL\x00")) || !bytes.Contains(ns[0].rumpf, []byte("C57P01\x00")) {
		t.Fatalf("erwartet genau die aufgezeichnete Fehlerantwort FATAL 57P01 und das Ende, erhalten %q %q", typen(ns), ns)
	}

	c2, r2 := rohVerbindung(t, rep.listen)
	anfrage(t, c2, "SELECT 3;")
	bisTyp(t, r2, 'Z')
	schreibe(t, c2, []byte{'X', 0, 0, 0, 4})
	rep.stop(t, 0)
	keineMeldung(t, "Ende, nächste Session", rep.stderr.String())
}

// Abdeckung: LH-FA-25/Boundary — replay mit einer Aufzeichnung, deren Ende
// `PANIC` ist: Der Client erhält die aufgezeichnete Fehlerantwort mit dem
// Schweregrad PANIC und das Ende der Verbindung, der Recorder meldet nichts
// (der Schweregrad ändert die Wiedergabe nicht).
func TestE2EReplayEndePanic(t *testing.T) {
	rep := startProzess(t, "replay", "--input", endeDatei(t, endePanicYAML), "--fail-on-unconsumed")
	c, r := rohVerbindung(t, rep.listen)
	anfrage(t, c, "SELECT 1;")
	ns := bisVerbindungsende(t, r)
	if typen(ns) != "E" || !bytes.Contains(ns[0].rumpf, []byte("SPANIC\x00")) {
		t.Fatalf("erwartet genau die aufgezeichnete Fehlerantwort PANIC und das Ende, erhalten %q %q", typen(ns), ns)
	}
	rep.stop(t, 0)
	keineMeldung(t, "PANIC", rep.stderr.String())
}

// Abdeckung: LH-FA-25/Negative, LH-FA-03/Boundary — replay: Sendet der Client
// vor der Interaktion mit Ende `Terminate`, ist sie nicht verbraucht: PGR-W2001
// nennt sie (1 von 2, Nummer 2), mit --fail-on-unconsumed PGR-E5002 und
// Exit-Code 5.
func TestE2EReplayEndeTerminateDavor(t *testing.T) {
	for _, streng := range []bool{false, true} {
		args, exit, code := []string{"--input", endeDatei(t, endeFatalYAML)}, 0, "PGR-W2001"
		if streng {
			args, exit, code = append(args, "--fail-on-unconsumed"), 5, "PGR-E5002"
		}
		rep := startProzess(t, "replay", args...)
		c, r := rohVerbindung(t, rep.listen)
		anfrage(t, c, "SELECT 1;")
		bisTyp(t, r, 'Z')
		schreibe(t, c, []byte{'X', 0, 0, 0, 4})
		bisVerbindungsende(t, r)
		rep.stop(t, exit)
		out := rep.stderr.String()
		if !strings.Contains(out, code) || !strings.Contains(out, "Session 1: 1 von 2 Interaktionen nicht verbraucht, die erste mit Nummer 2") {
			t.Fatalf("streng %v: erwartet %s mit Nummer 2:\n%s", streng, code, out)
		}
	}
}

// Abdeckung: LH-FA-25/Happy, LH-FA-09/Boundary, LH-FA-12/Boundary — replay: Eine
// aufgezeichnete Lebendprüfung `-- ping` mit Ende wird nicht übersprungen und
// nicht außer der Reihe beantwortet: Die erste Verbindung erhält mit der
// Lebendprüfung die Fehlerantwort und das Ende, die nächste die zweite Session;
// mit --fail-on-unconsumed ist der Exit-Code 0.
func TestE2EReplayEndeLebendpruefung(t *testing.T) {
	rep := startProzess(t, "replay", "--input", endeDatei(t, endeLebendYAML), "--fail-on-unconsumed")
	c, r := rohVerbindung(t, rep.listen)
	anfrage(t, c, "-- ping")
	ns := bisVerbindungsende(t, r)
	if typen(ns) != "E" || !bytes.Contains(ns[0].rumpf, []byte("C57P05\x00")) {
		t.Fatalf("erwartet die aufgezeichnete Fehlerantwort 57P05 und das Ende, erhalten %q %q", typen(ns), ns)
	}
	c2, r2 := rohVerbindung(t, rep.listen)
	anfrage(t, c2, "SELECT 3;")
	bisTyp(t, r2, 'Z')
	schreibe(t, c2, []byte{'X', 0, 0, 0, 4})
	rep.stop(t, 0)
	keineMeldung(t, "Lebendprüfung mit Ende", rep.stderr.String())
}
