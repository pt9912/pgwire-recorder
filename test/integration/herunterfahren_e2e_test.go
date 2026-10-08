//go:build integration

package integration_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

// codeFrist ist der Meldungscode einer Interaktion, die die Frist beim
// Herunterfahren unvollständig beendet.
const codeFrist = "PGR-E4006"

// warteAbgelehnt wartet höchstens frist, bis listen keine Verbindung mehr
// annimmt; dann hat der Prozess das erste Signal behandelt (LH-FA-13.a
// *Zweites Signal*). Läuft die Frist ab, endet der Test über nachKill.
func (r *recorder) warteAbgelehnt(t *testing.T, frist time.Duration) {
	t.Helper()
	deadline := time.Now().Add(frist)
	for {
		c, err := net.DialTimeout("tcp", r.listen, 200*time.Millisecond)
		if err != nil {
			return
		}
		c.Close()
		if time.Now().After(deadline) {
			t.Fatal(r.nachKill("Port nimmt nach dem ersten Signal weiter Verbindungen an"))
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// laufendeInteraktion verbindet sich über den Recorder, schließt `SELECT 1;`
// ab und beginnt eine Extended-Interaktion ohne Sync: Prepare mit Flush, deren
// ParseComplete sie liest. Die Interaktion läuft danach im Recorder.
func laufendeInteraktion(ctx context.Context, t *testing.T, listen string) *pgconn.PgConn {
	t.Helper()
	conn, err := pgconn.Connect(ctx, dsn(listen))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(ctx, "SELECT 1;").ReadAll(); err != nil {
		t.Fatal(err)
	}
	p := conn.StartPipeline(ctx)
	p.SendPrepare("s1", "SELECT 2", nil)
	p.SendFlushRequest()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	if _, err := p.GetResults(); err != nil {
		t.Fatal(err)
	}
	return conn
}

// nurErsteInteraktion prüft, dass die Aufzeichnung `SELECT 1;` trägt und keine
// Extended-Interaktion.
func nurErsteInteraktion(t *testing.T, output string) {
	t.Helper()
	if text := lies(t, output); !strings.Contains(text, "sql: SELECT 1;") || strings.Contains(text, "type: extended") || strings.Contains(text, "SELECT 2") {
		t.Fatalf("Aufzeichnung nach dem Zwangsende:\n%s", text)
	}
	laedt(t, output)
}

// Abdeckung: LH-FA-13/Boundary — läuft beim SIGTERM eine Interaktion, wartet
// record höchstens --shutdown-timeout (1s); danach endet die Session
// zwangsweise: die laufende Interaktion fehlt in der Aufzeichnung, die
// abgeschlossene bleibt, der Client erhält die Fehlerantwort PGR-E4006, das Log
// nennt Session und verworfene Interaktion und beim Beginn sessions=1, und der
// Lauf endet mit Exit-Code 4.
func TestE2ERecordFristLaeuftAb(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--shutdown-timeout", "1s")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn := laufendeInteraktion(ctx, t, rec.listen)
	defer conn.Conn().Close()

	rec.signal(t, syscall.SIGTERM)
	rec.warteEnde(t, 15*time.Second, "Recorder endet nicht nach Ablauf der Frist von 1s")
	rec.pruefeExit(t, 4)

	nurErsteInteraktion(t, output)
	text := rec.stderr.String()
	for _, want := range []string{"sessions=1", "level=ERROR", "code=" + codeFrist, "Interaktion 2 nicht abgeschlossen und verworfen", "als Session 1 geschrieben"} {
		if !strings.Contains(text, want) {
			t.Fatalf("%q fehlt im Log:\n%s", want, text)
		}
	}
	_ = conn.Conn().SetReadDeadline(time.Now().Add(5 * time.Second))
	msg, err := pgproto3.NewFrontend(conn.Conn(), conn.Conn()).Receive()
	if e, ok := msg.(*pgproto3.ErrorResponse); err != nil || !ok || e.Severity != "FATAL" || !strings.Contains(e.Message, codeFrist) {
		t.Fatalf("Client erhält %#v, %v statt der Fehlerantwort %s", msg, err, codeFrist)
	}
}

// Abdeckung: LH-FA-13/Boundary — mit --shutdown-timeout 0 wartet record nach
// dem ersten SIGTERM ohne Frist auf die laufende Interaktion; das zweite SIGTERM
// lässt die Frist sofort ablaufen, die Session endet zwangsweise, und der Lauf
// endet mit Exit-Code 4 und der Aufzeichnung der abgeschlossenen Interaktion.
func TestE2ERecordZweitesSignal(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--shutdown-timeout", "0")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn := laufendeInteraktion(ctx, t, rec.listen)
	defer conn.Conn().Close()

	rec.signal(t, syscall.SIGTERM)
	rec.warteAbgelehnt(t, 10*time.Second)
	// Ohne Frist endet der Prozess nicht von selbst; endete er doch, meldet
	// signal das Ende vor dem zweiten Signal.
	time.Sleep(time.Second)
	rec.signal(t, syscall.SIGTERM)
	rec.warteEnde(t, 15*time.Second, "Recorder endet nicht nach dem zweiten SIGTERM")
	rec.pruefeExit(t, 4)

	nurErsteInteraktion(t, output)
	if !strings.Contains(rec.stderr.String(), codeFrist) {
		t.Fatalf("%s fehlt im Log:\n%s", codeFrist, rec.stderr.String())
	}
}

// Abdeckung: LH-FA-13/Boundary — ein drittes SIGTERM bleibt ohne Wirkung:
// Während das Zwangsende nach dem zweiten Signal an einem Client schreibt, der
// eine große Ausgabe nicht liest (bis zu 1 s), trifft das dritte ein; der
// Prozess führt das Zwangsende zu Ende, schreibt die Aufzeichnung mit der
// abgeschlossenen Interaktion und endet mit Exit-Code 4 (LH-FA-13.a *Weitere
// Signale*).
func TestE2ERecordDrittesSignal(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--shutdown-timeout", "60s")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Conn().Close()
	if _, err := conn.Exec(ctx, "SELECT 1;").ReadAll(); err != nil {
		t.Fatal(err)
	}
	p := conn.StartPipeline(ctx)
	p.SendQueryParams("SELECT repeat('x', 1000000) FROM generate_series(1, 64)", nil, nil, nil, nil)
	p.SendFlushRequest()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	// Die Ausgabe bleibt ungelesen; nach einer Weile staut sie sich, und das
	// Schreiben an den Client blockiert.
	time.Sleep(2 * time.Second)

	rec.signal(t, syscall.SIGTERM)
	rec.warteAbgelehnt(t, 10*time.Second)
	rec.signal(t, syscall.SIGTERM)
	time.Sleep(300 * time.Millisecond)
	rec.signal(t, syscall.SIGTERM)
	rec.warteEnde(t, 15*time.Second, "Recorder endet nicht nach dem dritten SIGTERM")
	rec.pruefeExit(t, 4)

	nurErsteInteraktion(t, output)
	if !strings.Contains(rec.stderr.String(), codeFrist) {
		t.Fatalf("%s fehlt im Log:\n%s", codeFrist, rec.stderr.String())
	}
}
