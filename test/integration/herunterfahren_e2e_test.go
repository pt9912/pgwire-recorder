//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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
// Ergebnis sie liest. Die Interaktion läuft danach im Recorder; die Pipeline
// bleibt offen.
func laufendeInteraktion(ctx context.Context, t *testing.T, listen string) (*pgconn.PgConn, *pgconn.Pipeline) {
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
	return conn, p
}

// fehlerantwortErhalten liest, was der Client nach dem Ende des Prozesses noch
// empfängt, und prüft, dass es eine ErrorResponse (Typ 'E') mit Schweregrad
// FATAL und dem Meldungscode code ist.
func fehlerantwortErhalten(t *testing.T, conn *pgconn.PgConn, code string) {
	t.Helper()
	_ = conn.Conn().SetReadDeadline(time.Now().Add(5 * time.Second))
	rest, err := io.ReadAll(conn.Conn())
	if len(rest) == 0 || rest[0] != 'E' || !bytes.Contains(rest, []byte("SFATAL\x00")) || !bytes.Contains(rest, []byte(code)) {
		t.Fatalf("Client erhält %q (%v) statt der Fehlerantwort %s", rest, err, code)
	}
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
	conn, _ := laufendeInteraktion(ctx, t, rec.listen)
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
	fehlerantwortErhalten(t, conn, codeFrist)
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
	conn, _ := laufendeInteraktion(ctx, t, rec.listen)
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

// aufnehmenMitFlush zeichnet eine Session auf: `SELECT 1;` und eine
// Extended-Interaktion aus einer Flush-Gruppe (Parse und Describe von s1) und
// einer Sync-Gruppe; es liefert den Pfad der Aufzeichnung.
func aufnehmenMitFlush(t *testing.T) string {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, p := laufendeInteraktion(ctx, t, rec.listen)
	p.SendPipelineSync()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	if _, err := p.GetResults(); err != nil {
		t.Fatal(err)
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
	rec.stop(t, 0)
	return output
}

// Abdeckung: LH-FA-13/Boundary — pausiert ein Client im Replay nach der ersten
// Gruppe einer Extended-Interaktion, wartet replay nach SIGTERM höchstens
// --shutdown-timeout (1s); danach schließt es die Verbindung, der Client erhält
// die Fehlerantwort PGR-E4006, das Log nennt beim Beginn sessions=1, danach
// PGR-E4006 mit Session und Interaktion vor der Warnung PGR-W2001, und der Lauf
// endet mit Exit-Code 4.
func TestE2EReplayFristLaeuftAb(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	input := aufnehmenMitFlush(t)
	rep := startProzess(t, "replay", "--input", input, "--shutdown-timeout", "1s")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, _ := laufendeInteraktion(ctx, t, rep.listen)
	defer conn.Conn().Close()

	rep.signal(t, syscall.SIGTERM)
	rep.warteEnde(t, 15*time.Second, "Replay endet nicht nach Ablauf der Frist von 1s")
	rep.pruefeExit(t, 4)

	enthaeltInReihe(t, rep.stderr.String(),
		"sessions=1",
		"code="+codeFrist, "Session 1, Interaktion 2 nicht verbraucht",
		"Session 1: 1 von 2 Interaktionen nicht verbraucht, die erste mit Nummer 2", "code=PGR-W2001",
	)
	fehlerantwortErhalten(t, conn, codeFrist)
}

// Abdeckung: LH-FA-13/Negative — mit --fail-on-unconsumed merkt replay beim
// Zwangsende PGR-E4006 vor PGR-E5002 derselben Session; beide stehen in dieser
// Reihenfolge im Log, und der Lauf endet mit Exit-Code 4 (LH-FA-03.b).
func TestE2EReplayFristVorNichtVerbraucht(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	input := aufnehmenMitFlush(t)
	rep := startProzess(t, "replay", "--input", input, "--shutdown-timeout", "1s", "--fail-on-unconsumed")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, _ := laufendeInteraktion(ctx, t, rep.listen)
	defer conn.Conn().Close()

	rep.signal(t, syscall.SIGTERM)
	rep.warteEnde(t, 15*time.Second, "Replay endet nicht nach Ablauf der Frist von 1s")
	rep.pruefeExit(t, 4)
	enthaeltInReihe(t, rep.stderr.String(), "code="+codeFrist, "code=PGR-E5002")
}
