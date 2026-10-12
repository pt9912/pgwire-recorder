//go:build integration

package integration_test

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

// verbinde verbindet sich über listen mit dem Werkzeug und liefert die Verbindung.
func verbinde(t *testing.T, listen string) *pgconn.PgConn {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(listen))
	if err != nil {
		t.Fatalf("Verbindung über %s: %v", listen, err)
	}
	return conn
}

// frage führt sql auf conn aus und liefert den Fehler.
func frage(conn *pgconn.PgConn, sql string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := conn.Exec(ctx, sql).ReadAll()
	return err
}

// schliesse schließt conn mit Terminate.
func schliesse(t *testing.T, conn *pgconn.PgConn) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
}

// wartetAufSessions wartet höchstens 15 s, bis die YAML-Aufzeichnung path n Sessions
// trägt, und liefert ihren Text.
func wartetAufSessions(t *testing.T, path string, n int) string {
	t.Helper()
	for deadline := time.Now().Add(15 * time.Second); ; time.Sleep(50 * time.Millisecond) {
		data, err := os.ReadFile(path)
		if err == nil && strings.Count(string(data), "- id:") == n {
			return string(data)
		}
		if time.Now().After(deadline) {
			t.Fatalf("%d Sessions stehen nicht binnen 15 s in %s (%v)", n, path, err)
		}
	}
}

// Die Anfangsbytes der Verbindungen ohne Startnachricht: SSLRequest und
// GSSENCRequest (Länge 8) und CancelRequest (Länge 16, Prozess und Schlüssel).
func sslRequest() []byte { return []byte{0, 0, 0, 8, 4, 210, 22, 47} }
func gssRequest() []byte { return []byte{0, 0, 0, 8, 4, 210, 22, 48} }
func cancelAnfr() []byte { return []byte{0, 0, 0, 16, 4, 210, 22, 46, 0, 0, 0, 1, 0, 0, 0, 2} }

// ohneStartnachricht öffnet eine Verbindung zu listen, sendet die Bytes (keine, wenn
// leer), liest, soweit der Server antwortet (bei SSLRequest und GSSENCRequest ein
// Byte N), und schließt sie; dass der Server keine weitere Nachricht sendet, ist
// bei den Anfragen mit Antwort das Ende der Verbindung durch ihn.
func ohneStartnachricht(t *testing.T, listen string, daten []byte, antwort int) {
	t.Helper()
	c, err := net.DialTimeout("tcp", listen, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	_ = c.SetDeadline(time.Now().Add(10 * time.Second))
	if len(daten) > 0 {
		schreibe(t, c, daten)
	}
	if antwort > 0 {
		got := make([]byte, antwort)
		if _, err := io.ReadFull(c, got); err != nil || string(got) != "N" {
			t.Fatalf("Antwort auf %v: %q, %v, erwartet N", daten, got, err)
		}
	}
}

// Abdeckung: LH-FA-12/Happy, LH-FA-12/Boundary, LH-FA-22/Happy — record
// --record-empty-sessions zeichnet eine Verbindung ohne Anfrage als Session ohne
// Interaktionen auf, in yaml und in sqlite: Die Aufzeichnung trägt die
// Kennzeichnung empty_sessions, die Sessions stehen in der Reihenfolge ihres
// Endes (hier 1 und 2 mit Anfrage, 3 ohne), und replay mit
// --session-assignment connection liest beide Dateien gleich: Die dritte
// Verbindung trifft eine Session 3 ohne Interaktion (PGR-E5001) und nicht die
// Fehlermeldung ohne Session (PGR-E5003). Ohne die Option steht weder die
// Kennzeichnung noch die Session in der Aufzeichnung.
func TestE2ERecordEmptySessions(t *testing.T) {
	for _, name := range []string{"rec.yaml", "rec.sqlite"} {
		output := filepath.Join(t.TempDir(), name)
		rec := startProzess(t, "record", append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-empty-sessions"}, formatArgs(name)...)...)
		erste := verbinde(t, rec.listen)
		if err := frage(erste, "SELECT 'eins' AS a;"); err != nil {
			t.Fatal(err)
		}
		schliesse(t, erste)
		wartetAufInhalt(t, output, "eins")
		zweite := verbinde(t, rec.listen)
		if err := frage(zweite, "SELECT 'zwei' AS b;"); err != nil {
			t.Fatal(err)
		}
		schliesse(t, zweite)
		wartetAufInhalt(t, output, "zwei")
		schliesse(t, verbinde(t, rec.listen))
		rec.stop(t, 0)

		data := lies(t, output)
		if !strings.Contains(data, "empty_sessions") {
			t.Errorf("%s: Kennzeichnung empty_sessions fehlt", name)
		}
		if name == "rec.yaml" {
			if !strings.Contains(data, "empty_sessions: true") || strings.Count(data, "- id:") != 3 || !strings.Contains(data, "- id: 3\n") || !strings.Contains(data, "interactions: []") {
				t.Errorf("%s:\n%s", name, data)
			}
		}

		rep := startProzess(t, "replay", "--input", output, "--session-assignment", "connection", "--fail-on-unconsumed")
		for i, sql := range []string{"SELECT 'eins' AS a;", "SELECT 'zwei' AS b;"} {
			conn := verbinde(t, rep.listen)
			if err := frage(conn, sql); err != nil {
				t.Fatalf("%s: Verbindung %d, %s: %v", name, i+1, sql, err)
			}
			schliesse(t, conn)
		}
		dritte := verbinde(t, rep.listen)
		err := frage(dritte, "SELECT 'drei';")
		if err == nil || !strings.Contains(err.Error(), "PGR-E5001") || !strings.Contains(err.Error(), "Session 3") {
			t.Errorf("%s: dritte Verbindung: %v, erwartet PGR-E5001 für Session 3", name, err)
		}
		_ = dritte.Close(context.Background())
		rep.stop(t, 5)
		if strings.Contains(rep.stderr.String(), "PGR-E5003") || strings.Contains(rep.stderr.String(), "PGR-E5002") {
			t.Errorf("%s: Ausgabe des Replay:\n%s", name, rep.stderr.String())
		}

		ohne := filepath.Join(t.TempDir(), name)
		rec = startProzess(t, "record", append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", ohne}, formatArgs(name)...)...)
		conn := verbinde(t, rec.listen)
		if err := frage(conn, "SELECT 'eins' AS a;"); err != nil {
			t.Fatal(err)
		}
		schliesse(t, conn)
		wartetAufInhalt(t, ohne, "eins")
		schliesse(t, verbinde(t, rec.listen))
		rec.stop(t, 0)
		if data := lies(t, ohne); strings.Contains(data, "empty_sessions") || (name == "rec.yaml" && strings.Count(data, "- id:") != 1) {
			t.Errorf("%s ohne Option:\n%s", name, data)
		}
	}
}

// Abdeckung: LH-FA-12/Happy, LH-FA-12/Boundary — zwei gleichzeitig offene
// Verbindungen werden zu zwei Sessions mit je der Reihenfolge ihrer Interaktionen,
// auch wenn die Anfragen abwechselnd kommen; die Kennungen folgen der Reihenfolge der
// Enden (die zweite Verbindung endet, nachdem das Ende der ersten in der Aufzeichnung
// steht), und eine dritte Verbindung ohne Anfrage ist keine Session, mit
// --record-empty-sessions die dritte. Replay ordnet die Sessions der Reihe nach
// den Verbindungen mit einer Anfrage zu.
func TestE2ERecordParalleleVerbindungen(t *testing.T) {
	for _, leer := range []bool{false, true} {
		output := filepath.Join(t.TempDir(), "rec.yaml")
		args := []string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output}
		if leer {
			args = append(args, "--record-empty-sessions")
		}
		rec := startProzess(t, "record", args...)
		a, b, ohne := verbinde(t, rec.listen), verbinde(t, rec.listen), verbinde(t, rec.listen)
		for _, f := range []struct {
			conn *pgconn.PgConn
			sql  string
		}{{a, "SELECT 'a1';"}, {b, "SELECT 'b1';"}, {a, "SELECT 'a2';"}, {b, "SELECT 'b2';"}, {b, "SELECT 'b3';"}} {
			if err := frage(f.conn, f.sql); err != nil {
				t.Fatalf("%s: %v", f.sql, err)
			}
		}
		schliesse(t, b)
		wartetAufSessions(t, output, 1)
		schliesse(t, a)
		wartetAufSessions(t, output, 2)
		schliesse(t, ohne)
		rec.stop(t, 0)

		text := lies(t, output)
		wantSessions := 2
		if leer {
			wantSessions = 3
		}
		if strings.Count(text, "- id:") != wantSessions {
			t.Errorf("leer %v: erwartet %d Sessions:\n%s", leer, wantSessions, text)
		}
		enthaeltInReihe(t, text, "- id: 1", "sql: SELECT 'b1';", "sql: SELECT 'b2';", "sql: SELECT 'b3';", "- id: 2", "sql: SELECT 'a1';", "sql: SELECT 'a2';")
		if strings.Contains(text[strings.Index(text, "- id: 2"):], "'b") {
			t.Errorf("leer %v: Anfrage der ersten Session in der zweiten:\n%s", leer, text)
		}

		rep := startProzess(t, "replay", "--input", output)
		for i, qs := range [][]string{{"SELECT 'b1';", "SELECT 'b2';", "SELECT 'b3';"}, {"SELECT 'a1';", "SELECT 'a2';"}} {
			conn := verbinde(t, rep.listen)
			for _, q := range qs {
				if err := frage(conn, q); err != nil {
					t.Errorf("leer %v: Verbindung %d, %s: %v", leer, i+1, q, err)
				}
			}
			schliesse(t, conn)
		}
		rep.stop(t, 0)
	}
}

// aufzeichnungMitLeerer schreibt mit record --record-empty-sessions eine
// Aufzeichnung mit drei Sessions: 1 mit "SELECT 1;", 2 ohne Anfrage und 3 mit
// "SELECT 3;", und liefert den Pfad.
func aufzeichnungMitLeerer(t *testing.T) string {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-empty-sessions")
	eins := verbinde(t, rec.listen)
	if err := frage(eins, "SELECT 1;"); err != nil {
		t.Fatal(err)
	}
	schliesse(t, eins)
	wartetAufSessions(t, output, 1)
	schliesse(t, verbinde(t, rec.listen))
	wartetAufSessions(t, output, 2)
	drei := verbinde(t, rec.listen)
	if err := frage(drei, "SELECT 3;"); err != nil {
		t.Fatal(err)
	}
	schliesse(t, drei)
	rec.stop(t, 0)
	return output
}

// Abdeckung: LH-FA-12/Happy, LH-FA-12/Boundary, LH-FA-03/Boundary — replay
// --session-assignment connection ordnet der n-ten Verbindung die Session n zu, auch
// einer ohne Interaktion (die zweite Verbindung stellt keine Anfrage, die dritte erhält
// die dritte Session); mit --fail-on-unconsumed ist eine Session ohne Interaktion weder
// nicht verbraucht noch nie zugeordnet (Exit-Code 0). Mit der Zuordnung first-request
// erhält die zweite Verbindung mit einer Anfrage die dritte Session. Bleibt eine Session
// mit Interaktionen ohne Verbindung, ist sie nie zugeordnet (PGR-W2001, mit
// --fail-on-unconsumed PGR-E5002 und Exit-Code 5), und eine Verbindung ohne Anfrage
// zu einer Session mit Interaktionen lässt diese nicht verbraucht.
func TestE2EReplayConnection(t *testing.T) {
	input := aufzeichnungMitLeerer(t)

	rep := startProzess(t, "replay", "--input", input, "--session-assignment", "connection", "--fail-on-unconsumed")
	for _, sql := range []string{"SELECT 1;", "", "SELECT 3;"} {
		conn := verbinde(t, rep.listen)
		if sql != "" {
			if err := frage(conn, sql); err != nil {
				t.Fatalf("%s: %v", sql, err)
			}
		}
		schliesse(t, conn)
	}
	rep.stop(t, 0)
	if rep.stderr.String() != "" && (strings.Contains(rep.stderr.String(), "PGR-W") || strings.Contains(rep.stderr.String(), "PGR-E")) {
		t.Errorf("Meldung trotz verbrauchter Sessions:\n%s", rep.stderr.String())
	}

	rep = startProzess(t, "replay", "--input", input)
	for _, sql := range []string{"SELECT 1;", "SELECT 3;"} {
		conn := verbinde(t, rep.listen)
		if err := frage(conn, sql); err != nil {
			t.Fatalf("first-request, %s: %v", sql, err)
		}
		schliesse(t, conn)
	}
	rep.stop(t, 0)

	for _, streng := range []bool{false, true} {
		args := []string{"--input", input, "--session-assignment", "connection"}
		exit, code := 0, "PGR-W2001"
		if streng {
			args, exit, code = append(args, "--fail-on-unconsumed"), 5, "PGR-E5002"
		}
		rep = startProzess(t, "replay", args...)
		conn := verbinde(t, rep.listen)
		if err := frage(conn, "SELECT 1;"); err != nil {
			t.Fatal(err)
		}
		schliesse(t, conn)
		schliesse(t, verbinde(t, rep.listen))
		rep.stop(t, exit)
		if out := rep.stderr.String(); !strings.Contains(out, code) || !strings.Contains(out, "1 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 3") {
			t.Errorf("streng %v: nie zugeordnete Session 3 ohne %s:\n%s", streng, code, out)
		}
	}

	rep = startProzess(t, "replay", "--input", input, "--session-assignment", "connection")
	schliesse(t, verbinde(t, rep.listen))
	rep.stop(t, 0)
	enthaeltInReihe(t, rep.stderr.String(), "Session 1: 1 von 1 Interaktionen nicht verbraucht", "PGR-W2001", "1 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 3", "PGR-W2001")
}

// Abdeckung: LH-FA-12/Negative, LH-FA-13/Negative — replay mit
// --session-assignment connection und einer Aufzeichnung ohne die Kennzeichnung
// empty_sessions endet als Startfehler PGR-E2001 mit Exit-Code 2, bevor der Port
// geöffnet wird (der Port ist belegt, und die Meldung ist nicht die des Öffnens);
// die Fehler des Ladens gehen vor (PGR-E3003 für eine beschädigte Aufzeichnung mit
// Exit-Code 3); die Kennzeichnung genügt (sie verlangt keine leere Session).
func TestE2EReplayConnectionOhneKennzeichnung(t *testing.T) {
	ohne, _ := aufnehmen(t, "SELECT 1;")
	belegt, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer belegt.Close()
	dir := t.TempDir()
	_, stderr, code := starteBis(t, dir, "replay", "--listen", belegt.Addr().String(), "--input", ohne, "--session-assignment", "connection")
	if code != 2 || !strings.Contains(stderr, "[PGR-E2001]") {
		t.Errorf("ohne Kennzeichnung: Exit-Code %d, stderr %q", code, stderr)
	}
	kaputt := filepath.Join(dir, "kaputt.yaml")
	if err := os.WriteFile(kaputt, []byte("format: pgwire-recorder\nversion: 1\nsessions: [x]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, stderr, code = starteBis(t, dir, "replay", "--listen", belegt.Addr().String(), "--input", kaputt, "--session-assignment", "connection")
	if code != 3 {
		t.Errorf("beschädigte Aufzeichnung: Exit-Code %d, stderr %q", code, stderr)
	}
	mit := aufzeichnungMitLeerer(t)
	_, stderr, code = starteBis(t, dir, "replay", "--listen", belegt.Addr().String(), "--input", mit, "--session-assignment", "connection")
	if code != 4 || !strings.Contains(stderr, "[PGR-E4001]") {
		t.Errorf("mit Kennzeichnung: Exit-Code %d, stderr %q, erwartet das Öffnen des belegten Ports", code, stderr)
	}
	rep := startProzess(t, "replay", "--input", mit, "--session-assignment", "connection")
	rep.stop(t, 0)
}

// Abdeckung: LH-FA-12/Boundary — eine Verbindung zählt erst mit der gelesenen
// Startnachricht: SSLRequest, GSSENCRequest, CancelRequest und eine Verbindung, die ohne
// Bytes endet, zählen im Replay mit
// --session-assignment connection nicht; die erste Verbindung danach erhält die
// Session 1.
func TestE2EReplayZaehlungDerVerbindungen(t *testing.T) {
	input := aufzeichnungMitLeerer(t)
	rep := startProzess(t, "replay", "--input", input, "--session-assignment", "connection")
	ohneStartnachricht(t, rep.listen, sslRequest(), 1)
	ohneStartnachricht(t, rep.listen, gssRequest(), 1)
	ohneStartnachricht(t, rep.listen, cancelAnfr(), 0)
	ohneStartnachricht(t, rep.listen, nil, 0)
	for _, sql := range []string{"SELECT 1;", "", "SELECT 3;"} {
		conn := verbinde(t, rep.listen)
		if sql != "" {
			if err := frage(conn, sql); err != nil {
				t.Fatalf("%s: %v", sql, err)
			}
		}
		schliesse(t, conn)
	}
	rep.stop(t, 0)
}

// Abdeckung: LH-FA-12/Boundary — im Record ist eine Verbindung keine Session, solange
// ihre Startnachricht nicht gelesen und die Session zum Upstream nicht eröffnet ist,
// auch mit --record-empty-sessions: SSLRequest, GSSENCRequest, CancelRequest, eine
// Verbindung, die ohne Bytes endet, und eine, deren
// Aufbau zum Upstream scheitert (PGR-E4002, Exit-Code 4); nach ihnen erhalten die
// nächsten Verbindungen die Kennungen 1 und 2, und die Aufzeichnung trägt die
// Kennzeichnung auch ohne Session.
func TestE2ERecordKeineVerbindung(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-empty-sessions")
	ohneStartnachricht(t, rec.listen, sslRequest(), 1)
	ohneStartnachricht(t, rec.listen, gssRequest(), 1)
	ohneStartnachricht(t, rec.listen, cancelAnfr(), 0)
	ohneStartnachricht(t, rec.listen, nil, 0)
	rec.stop(t, 0)
	if text := lies(t, output); !strings.Contains(text, "empty_sessions: true") || !strings.Contains(text, "sessions: []") {
		t.Errorf("Aufzeichnung ohne Session:\n%s", text)
	}

	output = filepath.Join(t.TempDir(), "rec.yaml")
	rec = startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-empty-sessions")
	ohneStartnachricht(t, rec.listen, sslRequest(), 1)
	ohneStartnachricht(t, rec.listen, cancelAnfr(), 0)
	ohneStartnachricht(t, rec.listen, nil, 0)
	conn := verbinde(t, rec.listen)
	if err := frage(conn, "SELECT 1;"); err != nil {
		t.Fatal(err)
	}
	schliesse(t, conn)
	wartetAufSessions(t, output, 1)
	ohneStartnachricht(t, rec.listen, gssRequest(), 1)
	schliesse(t, verbinde(t, rec.listen))
	rec.stop(t, 0)
	text := lies(t, output)
	if strings.Count(text, "- id:") != 2 || !strings.Contains(text, "- id: 2\n") {
		t.Errorf("erwartet die Kennungen 1 und 2:\n%s", text)
	}

	output = filepath.Join(t.TempDir(), "rec.yaml")
	rec = startProzess(t, "record", "--upstream", "127.0.0.1:1", "--output", output, "--record-empty-sessions")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := pgconn.Connect(ctx, dsn(rec.listen)); err == nil || !strings.Contains(err.Error(), codeUpstream) {
		t.Fatalf("erwartet Fehler mit %s, erhalten %v", codeUpstream, err)
	}
	rec.stop(t, 4)
	if text := lies(t, output); !strings.Contains(text, "empty_sessions: true") || !strings.Contains(text, "sessions: []") {
		t.Errorf("Aufzeichnung nach gescheitertem Aufbau:\n%s", text)
	}
}

// Abdeckung: LH-FA-12/Boundary, LH-FA-13/Boundary — trifft das Zwangsende eine Verbindung
// im Aufbau (der Upstream antwortet nicht), ist sie keine Session, auch mit
// --record-empty-sessions: Der Lauf endet nach der Frist ohne Meldung mit Exit-Code 0
// und einer Aufzeichnung ohne Session mit der Kennzeichnung.
func TestE2ERecordZwangsendeImAufbau(t *testing.T) {
	stumm, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer stumm.Close()
	angenommen := make(chan net.Conn, 4)
	go func() {
		for {
			c, err := stumm.Accept()
			if err != nil {
				return
			}
			angenommen <- c
		}
	}()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", stumm.Addr().String(), "--output", output, "--record-empty-sessions", "--shutdown-timeout", "300ms")
	c, err := net.DialTimeout("tcp", rec.listen, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	start := binary.BigEndian.AppendUint32(nil, 196608)
	start = append(start, "user\x00postgres\x00database\x00postgres\x00\x00"...)
	schreibe(t, c, append(binary.BigEndian.AppendUint32(nil, uint32(4+len(start))), start...))
	select {
	case up := <-angenommen:
		defer up.Close()
	case <-time.After(10 * time.Second):
		t.Fatal("der Aufbau zum Upstream beginnt nicht")
	}
	rec.stop(t, 0)
	if text := lies(t, output); !strings.Contains(text, "empty_sessions: true") || !strings.Contains(text, "sessions: []") {
		t.Errorf("Aufzeichnung nach dem Zwangsende im Aufbau:\n%s", text)
	}
	if strings.Contains(rec.stderr.String(), "PGR-E4006") {
		t.Errorf("Meldung für eine Verbindung im Aufbau:\n%s", rec.stderr.String())
	}
}

// Abdeckung: LH-FA-13/Negative, LH-FA-12/Boundary — ein Verbindungsfehler beendet nur
// seine Verbindung, und der Prozess merkt sich die Klasse: Im Replay endet die erste
// Verbindung mit einer abweichenden Anfrage (PGR-E5001), die zweite, schon offene,
// läuft vollständig durch, und der Lauf endet mit Exit-Code 5. Im Record endet die
// erste Verbindung mit einer nicht unterstützten Interaktion (PGR-E6001), die zweite
// wird aufgezeichnet; ein später folgender Verbindungsfehler anderer Klasse (das
// Zwangsende einer laufenden Interaktion, PGR-E4006) ändert den Exit-Code nicht: Er
// bleibt 6, der des ersten.
func TestE2EVerbindungsfehlerParallel(t *testing.T) {
	input := aufnehmenSessions(t, []string{"SELECT 1;"}, []string{"SELECT 2;"})
	rep := startProzess(t, "replay", "--input", input)
	a, b := verbinde(t, rep.listen), verbinde(t, rep.listen)
	if err := frage(a, "SELECT 99;"); err == nil || !strings.Contains(err.Error(), "PGR-E5001") {
		t.Fatalf("erste Verbindung: %v, erwartet PGR-E5001", err)
	}
	if err := frage(b, "SELECT 2;"); err != nil {
		t.Fatalf("zweite Verbindung nach dem Fehler der ersten: %v", err)
	}
	schliesse(t, b)
	_ = a.Close(context.Background())
	rep.stop(t, 5)

	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--shutdown-timeout", "500ms")
	a, b = verbinde(t, rec.listen), verbinde(t, rec.listen)
	if err := frage(a, "COPY (SELECT 1) TO STDOUT;"); err == nil || !strings.Contains(err.Error(), "PGR-E6001") {
		t.Fatalf("erste Verbindung: %v, erwartet PGR-E6001", err)
	}
	if err := frage(b, "SELECT 'zweite';"); err != nil {
		t.Fatalf("zweite Verbindung nach dem Fehler der ersten: %v", err)
	}
	schliesse(t, b)
	wartetAufSessions(t, output, 1)
	_ = a.Close(context.Background())

	lang, err := net.DialTimeout("tcp", rec.listen, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer lang.Close()
	_ = lang.SetDeadline(time.Now().Add(10 * time.Second))
	start := binary.BigEndian.AppendUint32(nil, 196608)
	start = append(start, "user\x00postgres\x00database\x00postgres\x00\x00"...)
	schreibe(t, lang, append(binary.BigEndian.AppendUint32(nil, uint32(4+len(start))), start...))
	bisTyp(t, bufio.NewReader(lang), 'Z')
	schreibe(t, lang, clientNachricht('Q', "SELECT pg_sleep(30);\x00"))
	time.Sleep(300 * time.Millisecond)
	rec.stop(t, 6)
	enthaeltInReihe(t, rec.stderr.String(), "PGR-E6001", "PGR-E4006")
	if text := lies(t, output); strings.Count(text, "- id:") != 1 || !strings.Contains(text, "sql: SELECT 'zweite';") {
		t.Errorf("erwartet nur die Session der zweiten Verbindung:\n%s", text)
	}
}

// Abdeckung: LH-FA-17/Boundary, LH-FA-12/Happy — PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS
// und der Schlüssel record_empty_sessions im Abschnitt record: der Konfigurationsdatei
// schalten --record-empty-sessions ein; --record-empty-sessions=false auf der
// Kommandozeile geht beiden vor. PGWIRE_RECORDER_SESSION_ASSIGNMENT und der Schlüssel
// session_assignment im Abschnitt replay: wählen connection wie die Option;
// --session-assignment=first-request geht beiden vor.
func TestE2ESessionsQuellen(t *testing.T) {
	for _, tc := range []struct {
		quelle string
		args   []string
		mit    bool
	}{
		{"Umgebung", nil, true},
		{"Datei", nil, true},
		{"Umgebung mit false auf der Kommandozeile", []string{"--record-empty-sessions=false"}, false},
		{"Datei mit false auf der Kommandozeile", []string{"--record-empty-sessions=false"}, false},
	} {
		dir := t.TempDir()
		output := filepath.Join(dir, "rec.yaml")
		if strings.HasPrefix(tc.quelle, "Umgebung") {
			t.Setenv("PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS", "true")
		} else if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("record:\n  record_empty_sessions: true\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		rec := startProzessIn(t, dir, "record", append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output}, tc.args...)...)
		schliesse(t, verbinde(t, rec.listen))
		rec.stop(t, 0)
		if got := strings.Contains(lies(t, output), "empty_sessions: true"); got != tc.mit {
			t.Errorf("%s: Kennzeichnung %v, erwartet %v", tc.quelle, got, tc.mit)
		}
		t.Setenv("PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS", "")
	}

	input := aufzeichnungMitLeerer(t)
	for _, tc := range []struct {
		quelle string
		args   []string
		dritte string
	}{
		{"Umgebung", nil, "SELECT 3;"},
		{"Datei", nil, "SELECT 3;"},
		{"Umgebung mit first-request auf der Kommandozeile", []string{"--session-assignment=first-request"}, "SELECT 3;"},
	} {
		dir := t.TempDir()
		if strings.HasPrefix(tc.quelle, "Umgebung") {
			t.Setenv("PGWIRE_RECORDER_SESSION_ASSIGNMENT", "connection")
		} else if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("replay:\n  session_assignment: connection\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		rep := startProzessIn(t, dir, "replay", append([]string{"--input", input}, tc.args...)...)
		// Mit connection erhält die zweite Verbindung die Session ohne Interaktion, mit
		// first-request die Session 3.
		eins, zwei := verbinde(t, rep.listen), verbinde(t, rep.listen)
		if err := frage(eins, "SELECT 1;"); err != nil {
			t.Fatalf("%s: %v", tc.quelle, err)
		}
		err := frage(zwei, tc.dritte)
		if connection := len(tc.args) == 0; connection != (err != nil) {
			t.Errorf("%s: zweite Verbindung: %v, erwartet Fehler %v", tc.quelle, err, connection)
		}
		schliesse(t, eins)
		_ = zwei.Close(context.Background())
		if len(tc.args) == 0 {
			rep.stop(t, 5)
		} else {
			rep.stop(t, 0)
		}
		t.Setenv("PGWIRE_RECORDER_SESSION_ASSIGNMENT", "")
	}
}

// Abdeckung: LH-FA-17/Negative, LH-FA-12/Negative — ein ungültiger Wert von
// --session-assignment (auch Connection, first_request und der leere) ist PGR-E2001
// mit dem Namen der Option, der der Umgebungsvariable PGR-E2001 mit ihrem Namen und der
// der Datei PGR-E2004 mit dem Schlüssel, je Exit-Code 2, die beiden letzten ohne den
// Wert; --session-assignment bei record und play und --record-empty-sessions bei replay
// und play sind PGR-E2001, ebenso ihre Umgebungsvariablen beim anderen Kommando; ihre
// Schlüssel in einem anderen Abschnitt sind PGR-E2004.
func TestE2ESessionsUngueltig(t *testing.T) {
	dir := t.TempDir()
	replay := []string{"replay", "--listen", "127.0.0.1:0", "--input", "r.yaml"}
	for _, wert := range []string{"GEHEIM", "Connection", "first_request", ""} {
		_, stderr, code := starteBis(t, dir, append(append([]string{}, replay...), "--session-assignment="+wert)...)
		if code != 2 || !strings.Contains(stderr, "[PGR-E2001]") || !strings.Contains(stderr, "session-assignment") {
			t.Errorf("Kommandozeile %q: Exit-Code %d, stderr %q", wert, code, stderr)
		}
	}
	t.Setenv("PGWIRE_RECORDER_SESSION_ASSIGNMENT", "GEHEIM")
	_, stderr, code := starteBis(t, dir, replay...)
	if code != 2 || !strings.Contains(stderr, "[PGR-E2001]") || !strings.Contains(stderr, "PGWIRE_RECORDER_SESSION_ASSIGNMENT") || strings.Contains(stderr, "GEHEIM") {
		t.Errorf("Umgebung: Exit-Code %d, stderr %q", code, stderr)
	}
	for _, kommando := range [][]string{
		{"record", "--listen", "127.0.0.1:0", "--upstream", "pg:5432", "--output", filepath.Join(dir, "r.yaml")},
		{"play", "--upstream", "pg:5432", "--input", "r.yaml"},
	} {
		_, stderr, code = starteBis(t, dir, kommando...)
		if code != 2 || !strings.Contains(stderr, "[PGR-E2001]") || !strings.Contains(stderr, "PGWIRE_RECORDER_SESSION_ASSIGNMENT") || strings.Contains(stderr, "GEHEIM") {
			t.Errorf("%s mit der Umgebungsvariable von replay: Exit-Code %d, stderr %q", kommando[0], code, stderr)
		}
		_, stderr, code = starteBis(t, dir, append(append([]string{}, kommando...), "--session-assignment=connection")...)
		if code != 2 || !strings.Contains(stderr, "[PGR-E2001]") {
			t.Errorf("%s mit --session-assignment: Exit-Code %d, stderr %q", kommando[0], code, stderr)
		}
	}
	t.Setenv("PGWIRE_RECORDER_SESSION_ASSIGNMENT", "")
	t.Setenv("PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS", "true")
	for _, kommando := range [][]string{replay, {"play", "--upstream", "pg:5432", "--input", "r.yaml"}} {
		_, stderr, code = starteBis(t, dir, kommando...)
		if code != 2 || !strings.Contains(stderr, "[PGR-E2001]") || !strings.Contains(stderr, "PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS") {
			t.Errorf("%s mit der Umgebungsvariable von record: Exit-Code %d, stderr %q", kommando[0], code, stderr)
		}
		_, stderr, code = starteBis(t, dir, append(append([]string{}, kommando...), "--record-empty-sessions")...)
		if code != 2 || !strings.Contains(stderr, "[PGR-E2001]") {
			t.Errorf("%s mit --record-empty-sessions: Exit-Code %d, stderr %q", kommando[0], code, stderr)
		}
	}
	t.Setenv("PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS", "")
	for _, f := range []struct{ datei, stelle string }{
		{"replay:\n  session_assignment: GEHEIM\n", "replay.session_assignment"},
		{"replay:\n  session_assignment: Connection\n", "replay.session_assignment"},
		{"record:\n  session_assignment: connection\n", "record.session_assignment"},
		{"replay:\n  record_empty_sessions: true\n", "replay.record_empty_sessions"},
		{"record:\n  record_empty_sessions: GEHEIM\n", "record.record_empty_sessions"},
	} {
		if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte(f.datei), 0o600); err != nil {
			t.Fatal(err)
		}
		_, stderr, code = starteBis(t, dir, replay...)
		if code != 2 || !strings.Contains(stderr, "[PGR-E2004]") || !strings.Contains(stderr, f.stelle) || strings.Contains(stderr, "GEHEIM") {
			t.Errorf("Datei %q: Exit-Code %d, stderr %q", f.datei, code, stderr)
		}
	}
}
