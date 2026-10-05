//go:build integration

package integration

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// aufnehmenSessions zeichnet je Liste von Anfragen eine Verbindung auf und
// liefert den Pfad der Aufzeichnung.
func aufnehmenSessions(t *testing.T, sessions ...[]string) string {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	for _, queries := range sessions {
		_ = beobachten(t, rec.listen, queries...)
	}
	rec.stop(t, 0)
	return output
}

// enthaeltInReihe prüft, dass die Teile in dieser Reihenfolge im Text stehen.
func enthaeltInReihe(t *testing.T, text string, teile ...string) {
	t.Helper()
	rest := text
	for _, teil := range teile {
		i := strings.Index(rest, teil)
		if i < 0 {
			t.Fatalf("%q fehlt oder steht nicht in Reihe:\n%s", teil, text)
		}
		rest = rest[i+len(teil):]
	}
}

// Abdeckung: LH-FA-13/Happy, LH-FA-13/Negative, LH-FA-17/Boundary — endet eine
// Verbindung vor dem Verbrauch aller Interaktionen, warnt replay ohne
// --fail-on-unconsumed mit PGR-W2001 und endet mit Exit-Code 0; mit der Option
// oder PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=true ist es PGR-E5002 mit demselben
// Text und Exit-Code 5; die Option der Kommandozeile geht der
// Umgebungsvariable vor; eine verbrauchte Session meldet auch mit der Option
// nichts.
func TestE2EReplayNichtVerbraucht(t *testing.T) {
	input := aufnehmenSessions(t, []string{"SELECT 1;", "SELECT 2;"})
	const text = "Session 1: 1 von 2 Interaktionen nicht verbraucht, die erste mit Nummer 2"

	for _, fall := range []struct {
		env  string
		args []string
		exit int
		code string
	}{
		{"", nil, 0, "PGR-W2001"},
		{"", []string{"--fail-on-unconsumed"}, 5, "PGR-E5002"},
		{"true", nil, 5, "PGR-E5002"},
		{"true", []string{"--fail-on-unconsumed=false"}, 0, "PGR-W2001"},
	} {
		t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", fall.env)
		rep := startProzess(t, "replay", append([]string{"--input", input}, fall.args...)...)
		_ = beobachten(t, rep.listen, "SELECT 1;")
		rep.stop(t, fall.exit)
		out := rep.stderr.String()
		if !strings.Contains(out, "code="+fall.code) || !strings.Contains(out, text) {
			t.Fatalf("Umgebung %q, %v: erwartet %s mit %q:\n%s", fall.env, fall.args, fall.code, text, out)
		}
		anderer := map[string]string{"PGR-W2001": "PGR-E5002", "PGR-E5002": "PGR-W2001"}[fall.code]
		if strings.Contains(out, anderer) {
			t.Fatalf("Umgebung %q, %v: %s neben %s:\n%s", fall.env, fall.args, anderer, fall.code, out)
		}
	}

	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	rep := startProzess(t, "replay", "--input", input, "--fail-on-unconsumed")
	_ = beobachten(t, rep.listen, "SELECT 1;", "SELECT 2;")
	rep.stop(t, 0)
	if out := rep.stderr.String(); strings.Contains(out, "PGR-E5002") || strings.Contains(out, "PGR-W2001") {
		t.Fatalf("verbrauchte Session gemeldet:\n%s", out)
	}
}

// Abdeckung: LH-FA-13/Negative — mit --fail-on-unconsumed
// bestimmen nie zugeordnete Sessions den Exit-Code 5, wenn kein
// Verbindungsfehler auftrat; endet eine Verbindung durch einen anderen
// Verbindungsfehler (PGR-E6001), steht er vor PGR-E5002 ihrer Session und
// vor PGR-E5002 der nie zugeordneten Sessions im Log, und der Exit-Code ist
// seiner (6).
func TestE2EReplayNichtVerbrauchtRangfolge(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	input := aufnehmenSessions(t, []string{"SELECT 1;", "SELECT 2;"}, []string{"SELECT 3;"})
	const nie = "1 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 2"

	rep := startProzess(t, "replay", "--input", input, "--fail-on-unconsumed")
	_ = beobachten(t, rep.listen, "SELECT 1;", "SELECT 2;")
	rep.stop(t, 5)
	if out := rep.stderr.String(); !strings.Contains(out, "PGR-E5002") || !strings.Contains(out, nie) {
		t.Fatalf("nie zugeordnete Session ohne PGR-E5002:\n%s", out)
	}

	rep = startProzess(t, "replay", "--input", input, "--fail-on-unconsumed")
	funktionsaufrufNachAnfrage(t, rep.listen, "SELECT 1;")
	rep.stop(t, 6)
	enthaeltInReihe(t, rep.stderr.String(),
		"PGR-E6001",
		"PGR-E5002", "Session 1: 1 von 2 Interaktionen nicht verbraucht, die erste mit Nummer 2",
		"PGR-E5002", nie,
	)
}

// funktionsaufrufNachAnfrage verbindet sich ohne Treiber, stellt die Anfrage,
// liest bis ReadyForQuery und sendet dann einen FunctionCall, den das Werkzeug
// nicht unterstützt; es liest, bis der Server die Verbindung schließt.
func funktionsaufrufNachAnfrage(t *testing.T, listen, sql string) {
	t.Helper()
	c, err := net.DialTimeout("tcp", listen, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	_ = c.SetDeadline(time.Now().Add(10 * time.Second))
	r := bufio.NewReader(c)

	start := []byte{0, 0, 0, 0, 0, 3, 0, 0}
	start = append(start, "user\x00postgres\x00database\x00postgres\x00\x00"...)
	binary.BigEndian.PutUint32(start, uint32(len(start)))
	schreibe(t, c, start)
	bisTyp(t, r, 'Z')

	query := append([]byte{'Q', 0, 0, 0, 0}, sql+"\x00"...)
	binary.BigEndian.PutUint32(query[1:], uint32(len(query)-1))
	schreibe(t, c, query)
	bisTyp(t, r, 'Z')

	// FunctionCall: OID 0, keine Argumentformate, keine Argumente, Ergebnis
	// im Textformat.
	schreibe(t, c, []byte{'F', 0, 0, 0, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	if typ := bisTyp(t, r, 'E'); typ != 'E' {
		t.Fatalf("erwartet ErrorResponse, erhalten %q", typ)
	}
	if _, err := io.Copy(io.Discard, r); err != nil {
		t.Fatalf("Verbindung endet nicht: %v", err)
	}
}

func schreibe(t *testing.T, c net.Conn, b []byte) {
	t.Helper()
	if _, err := c.Write(b); err != nil {
		t.Fatal(err)
	}
}

// bisTyp liest Nachrichten bis zu einer des Typs typ oder einer ErrorResponse
// und liefert deren Typ; eine ErrorResponse vor typ bricht ab.
func bisTyp(t *testing.T, r *bufio.Reader, typ byte) byte {
	t.Helper()
	for {
		kopf := make([]byte, 5)
		if _, err := io.ReadFull(r, kopf); err != nil {
			t.Fatalf("Nachricht lesen: %v", err)
		}
		rumpf := make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)
		if _, err := io.ReadFull(r, rumpf); err != nil {
			t.Fatalf("Nachricht lesen: %v", err)
		}
		switch kopf[0] {
		case typ:
			return typ
		case 'E':
			t.Fatalf("ErrorResponse statt %q: %q", typ, rumpf)
		}
	}
}

// Abdeckung: LH-FA-13/Negative — nach einem Startfehler
// (Listen-Port belegt, PGR-E4001) prüft replay auch mit --fail-on-unconsumed
// keine nie zugeordneten Sessions: Exit-Code 4, weder PGR-E5002 noch
// PGR-W2001.
func TestE2EReplayNichtVerbrauchtStartfehler(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	input := aufnehmenSessions(t, []string{"SELECT 1;"})
	belegt, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer belegt.Close()
	out, err := exec_(t, "replay", "--listen", belegt.Addr().String(), "--input", input, "--fail-on-unconsumed")
	if code := exitCodeOf(err); code != 4 || !strings.Contains(out, "PGR-E4001") {
		t.Fatalf("Exit-Code %d, Ausgabe:\n%s", code, out)
	}
	if strings.Contains(out, "PGR-E5002") || strings.Contains(out, "PGR-W2001") {
		t.Fatalf("nach dem Startfehler geprüft:\n%s", out)
	}
}

// Abdeckung: LH-FA-13/Boundary — beendet das Herunterfahren
// eine offene, ruhende Verbindung, deren Session nicht verbraucht ist, meldet
// replay mit --fail-on-unconsumed PGR-E5002 für diese Session und endet mit
// Exit-Code 5.
func TestE2EReplayNichtVerbrauchtHerunterfahren(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	input := aufnehmenSessions(t, []string{"SELECT 1;", "SELECT 2;"})
	rep := startProzess(t, "replay", "--input", input, "--fail-on-unconsumed")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rep.listen))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(ctx)
	if _, err := conn.Exec(ctx, "SELECT 1;").ReadAll(); err != nil {
		t.Fatal(err)
	}
	rep.stop(t, 5)
	enthaeltInReihe(t, rep.stderr.String(), "PGR-E5002", "Session 1: 1 von 2 Interaktionen nicht verbraucht, die erste mit Nummer 2")
}
