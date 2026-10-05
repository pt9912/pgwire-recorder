//go:build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// extendedAblauf führt über listen einen festen Ablauf mit pgx in seinem
// Standardmodus (Prepared Statements über das Extended Query Protocol) aus und
// liefert das beobachtbare Verhalten als Text: dasselbe Statement mit
// verschiedenen und wiederholten Werten, einen NULL-Parameter, mehrere Zeilen,
// einen Fehler, einen Batch, in dem nach einem Fehler die folgende Anweisung bis
// zum Sync verworfen wird, und eine einfache Anfrage.
func extendedAblauf(t *testing.T, listen string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn(listen))
	if err != nil {
		t.Fatalf("Verbindung über %s: %v", listen, err)
	}
	defer conn.Close(ctx)
	var b strings.Builder
	for _, wert := range []string{"eins", "zwei", "eins"} {
		var got string
		err := conn.QueryRow(ctx, "SELECT $1::text || '!' AS t", wert).Scan(&got)
		fmt.Fprintf(&b, "wert %s: %q %v\n", wert, got, err)
	}
	var istNull bool
	err = conn.QueryRow(ctx, "SELECT $1::text IS NULL AS n", nil).Scan(&istNull)
	fmt.Fprintf(&b, "null: %v %v\n", istNull, err)
	rows, err := conn.Query(ctx, "SELECT g, g * 2 AS doppelt FROM generate_series(1, $1::int) g", 3)
	if err != nil {
		t.Fatalf("Zeilen: %v", err)
	}
	for _, f := range rows.FieldDescriptions() {
		fmt.Fprintf(&b, "spalte %s %d\n", f.Name, f.DataTypeOID)
	}
	for rows.Next() {
		var g, d int
		if err := rows.Scan(&g, &d); err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&b, "zeile %d %d\n", g, d)
	}
	fmt.Fprintf(&b, "befehl %s %v\n", rows.CommandTag(), rows.Err())
	_, err = conn.Exec(ctx, "SELECT 1/$1::int", 0)
	fmt.Fprintf(&b, "fehler: %s\n", pgFehler(err))

	batch := &pgx.Batch{}
	batch.Queue("SELECT 10/$1::int", 0)
	batch.Queue("SELECT $1::int + 1", 41)
	br := conn.SendBatch(ctx, batch)
	_, err1 := br.Exec()
	_, err2 := br.Exec()
	err3 := br.Close()
	fmt.Fprintf(&b, "batch: %s | %s | %s\n", pgFehler(err1), pgFehler(err2), pgFehler(err3))

	var n int
	err = conn.QueryRow(ctx, "SELECT 7", pgx.QueryExecModeSimpleProtocol).Scan(&n)
	fmt.Fprintf(&b, "einfach: %d %v\n", n, err)
	return b.String()
}

// pgFehler beschreibt einen Fehler mit SQLSTATE und Meldung des Servers, sonst
// mit seinem Text.
func pgFehler(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code + " " + pgErr.Message
	}
	return fmt.Sprint(err)
}

// aufnehmenExtended zeichnet extendedAblauf über record gegen die reale Instanz
// in output auf und liefert die Sicht des Clients beim Aufzeichnen.
func aufnehmenExtended(t *testing.T, output string) string {
	t.Helper()
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	sicht := extendedAblauf(t, rec.listen)
	rec.stop(t, 0)
	for _, z := range []string{`wert eins: "eins!" <nil>`, `wert zwei: "zwei!" <nil>`, "null: true <nil>", "zeile 3 6", "fehler: 22012 division by zero", "batch: 22012 division by zero", "einfach: 7 <nil>"} {
		if !strings.Contains(sicht, z) {
			t.Fatalf("Sicht beim Aufzeichnen ohne %q:\n%s", z, sicht)
		}
	}
	return sicht
}

// Abdeckung: LH-FA-18/Happy, LH-FA-18/Boundary, LH-FA-09/Happy, LH-FA-11/Happy —
// pgx im Standardmodus erhält über replay aus einer von record geschriebenen
// Aufzeichnung dieselbe Sicht wie beim Aufzeichnen: dasselbe Statement mit
// verschiedenen und wiederholten Werten je in seiner Position, NULL-Parameter,
// Spalten und Zeilen, eine Fehlerantwort und im Batch das Verwerfen bis zum
// Sync; drei aufeinanderfolgende Läufe zeigen dieselbe Sicht, und jeder endet
// mit Exit-Code 0.
func TestE2EReplayExtendedPgx(t *testing.T) {
	input := filepath.Join(t.TempDir(), "rec.yaml")
	aufgezeichnet := aufnehmenExtended(t, input)
	for lauf := 1; lauf <= 3; lauf++ {
		rep := startProzess(t, "replay", "--input", input)
		got := extendedAblauf(t, rep.listen)
		rep.stop(t, 0)
		if got != aufgezeichnet {
			t.Fatalf("Lauf %d weicht von der Sicht beim Aufzeichnen ab:\n%s\n--- aufgezeichnet:\n%s\n--- stderr:\n%s", lauf, got, aufgezeichnet, rep.stderr.String())
		}
	}
}

// pipelineAblauf führt über listen eine Pipeline aus: eine Flush-Gruppe
// (Parse, Describe), auf deren Antwort der Client wartet, danach eine
// Sync-Gruppe mit zwei Ausführungen. Es liefert das beobachtbare Verhalten als
// Text.
func pipelineAblauf(t *testing.T, listen string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(listen))
	if err != nil {
		t.Fatalf("Verbindung über %s: %v", listen, err)
	}
	defer conn.Close(ctx)
	var b strings.Builder
	p := conn.StartPipeline(ctx)
	p.SendPrepare("s1", "SELECT $1::text AS t", nil)
	p.SendFlushRequest()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	r, err := p.GetResults()
	if sd, ok := r.(*pgconn.StatementDescription); ok {
		fmt.Fprintf(&b, "prepare: %v %d\n", sd.ParamOIDs, len(sd.Fields))
	} else {
		t.Fatalf("Prepare: %#v, %v", r, err)
	}
	p.SendQueryPrepared("s1", [][]byte{[]byte("a")}, nil, nil)
	p.SendQueryPrepared("s1", [][]byte{nil}, nil, nil)
	p.SendPipelineSync()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		r, err := p.GetResults()
		if err != nil {
			t.Fatalf("Ausführung %d: %v", i+1, err)
		}
		res := r.(*pgconn.ResultReader).Read()
		fmt.Fprintf(&b, "ausführung %d: %q %s %v\n", i+1, res.Rows, res.CommandTag, res.Err)
	}
	if _, err := p.GetResults(); err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// Abdeckung: LH-FA-18/Happy — eine Pipeline aus einer Flush-Gruppe, auf deren
// Antwort der Client wartet, und einer Sync-Gruppe erhält über replay dieselbe
// Sicht wie beim Aufzeichnen; die Antwort der Flush-Gruppe geht nach deren
// Flush an den Client.
func TestE2EReplayExtendedPipeline(t *testing.T) {
	input := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), input)
	aufgezeichnet := pipelineAblauf(t, rec.listen)
	rec.stop(t, 0)
	if !strings.Contains(aufgezeichnet, "prepare: [25] 1") || !strings.Contains(aufgezeichnet, `ausführung 1: [["a"]] SELECT 1 <nil>`) {
		t.Fatalf("Sicht beim Aufzeichnen:\n%s", aufgezeichnet)
	}
	rep := startProzess(t, "replay", "--input", input)
	got := pipelineAblauf(t, rep.listen)
	rep.stop(t, 0)
	if got != aufgezeichnet {
		t.Fatalf("Replay-Sicht:\n%s\n--- aufgezeichnet:\n%s\n--- stderr:\n%s", got, aufgezeichnet, rep.stderr.String())
	}
}

// Abdeckung: LH-FA-18/Happy — Replay-Hälfte: Gegendruck in beiden Richtungen
// (gegendruckAblauf), über `record` aufgezeichnet, läuft über `replay` binnen
// des Zeitlimits durch, auch die verzögerte große Ausgabe der Flush-Gruppe, die
// der Client nicht abwartet; der Lauf endet mit Exit-Code 0.
func TestE2EReplayExtendedGegendruck(t *testing.T) {
	input := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), input)
	gegendruckAblauf(t, rec.listen)
	rec.stop(t, 0)
	rep := startProzess(t, "replay", "--input", input)
	gegendruckAblauf(t, rep.listen)
	rep.stop(t, 0)
}

// Abdeckung: LH-FA-13/Boundary — Replay: kommt SIGTERM, nachdem die
// Flush-Gruppe einer Pipeline beantwortet ist und bevor ihre Sync-Gruppe
// eintrifft, beantwortet replay die Sync-Gruppe noch wie aufgezeichnet (die
// Sicht des Clients gleicht der beim Aufzeichnen), schließt danach die
// Verbindung und endet von selbst mit Exit-Code 0.
func TestE2EReplayExtendedSigtermMittenInFolge(t *testing.T) {
	input := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), input)
	aufgezeichnet := pipelineAblauf(t, rec.listen)
	rec.stop(t, 0)
	rep := startProzess(t, "replay", "--input", input)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rep.listen))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(ctx)
	p := conn.StartPipeline(ctx)
	p.SendPrepare("s1", "SELECT $1::text AS t", nil)
	p.SendFlushRequest()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	r, err := p.GetResults()
	if sd, ok := r.(*pgconn.StatementDescription); ok {
		fmt.Fprintf(&b, "prepare: %v %d\n", sd.ParamOIDs, len(sd.Fields))
	} else {
		t.Fatalf("Prepare: %#v, %v", r, err)
	}
	if err := rep.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)
	p.SendQueryPrepared("s1", [][]byte{[]byte("a")}, nil, nil)
	p.SendQueryPrepared("s1", [][]byte{nil}, nil, nil)
	p.SendPipelineSync()
	if err := p.Flush(); err != nil {
		t.Fatalf("Sync-Gruppe nach SIGTERM: %v", err)
	}
	for i := 0; i < 2; i++ {
		r, err := p.GetResults()
		if err != nil {
			t.Fatalf("Ausführung %d nach SIGTERM: %v", i+1, err)
		}
		res := r.(*pgconn.ResultReader).Read()
		fmt.Fprintf(&b, "ausführung %d: %q %s %v\n", i+1, res.Rows, res.CommandTag, res.Err)
	}
	if _, err := p.GetResults(); err != nil {
		t.Fatalf("Sync nach SIGTERM: %v", err)
	}
	if b.String() != aufgezeichnet {
		t.Fatalf("Sicht nach SIGTERM:\n%s\n--- aufgezeichnet:\n%s", b.String(), aufgezeichnet)
	}

	done := make(chan error, 1)
	go func() { done <- rep.cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		_ = rep.cmd.Process.Kill()
		t.Fatalf("replay endet nach dem Sync nicht von selbst\n%s", rep.stderr.String())
	}
	if code := rep.cmd.ProcessState.ExitCode(); code != 0 {
		t.Fatalf("Exit-Code %d\n%s", code, rep.stderr.String())
	}
	if strings.Contains(rep.stderr.String(), "PGR-W2001") {
		t.Fatalf("Interaktion nicht verbraucht:\n%s", rep.stderr.String())
	}
}

// Abdeckung: LH-FA-18/Negative, LH-FA-10/Happy — führt pgx im Replay dasselbe
// Prepared Statement mit einem anderen Parameterwert aus als aufgezeichnet,
// erhält es einen eindeutigen Fehler mit PGR-E5001 und keine Zeile; die
// Diagnose nennt den Wert nicht, und der Lauf endet mit Exit-Code 5.
func TestE2EReplayExtendedAbweichung(t *testing.T) {
	input := filepath.Join(t.TempDir(), "rec.yaml")
	aufnehmenExtended(t, input)
	rep := startProzess(t, "replay", "--input", input)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn(rep.listen))
	if err != nil {
		t.Fatal(err)
	}
	var got string
	err = conn.QueryRow(ctx, "SELECT $1::text || '!' AS t", "anderer-wert").Scan(&got)
	if err == nil || !strings.Contains(err.Error(), "PGR-E5001") || got != "" {
		t.Fatalf("erwartet PGR-E5001 ohne Zeile, erhalten %q, %v", got, err)
	}
	_ = conn.Close(ctx)
	rep.stop(t, 5)
	for _, quelle := range []string{err.Error(), rep.stderr.String()} {
		if strings.Contains(quelle, "anderer-wert") {
			t.Fatalf("Diagnose nennt den Parameterwert:\n%s", quelle)
		}
	}
}

// Abdeckung: LH-FA-18/Happy — Vorbereitung von Abnahmeszenario 7: pgx im
// Standardmodus wird über record gegen die reale Instanz aufgezeichnet; die
// Aufzeichnung und die Sicht des Clients liegen danach unter PGR_DATEN.
func TestE2EVorbereitungExtendedOhnePostgres(t *testing.T) {
	daten := os.Getenv("PGR_DATEN")
	if daten == "" {
		t.Skip("ohne PGR_DATEN keine zweite Phase")
	}
	sicht := aufnehmenExtended(t, filepath.Join(daten, "rec-extended.yaml"))
	if err := os.WriteFile(filepath.Join(daten, "sicht-extended.txt"), []byte(sicht), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-18/Happy, LH-FA-04/Happy — Abnahmeszenario 7: bei gestoppter
// PostgreSQL-Instanz erhält pgx im Standardmodus, nur mit Host und Port des
// Recorders, aus der in der ersten Phase geschriebenen Aufzeichnung dieselbe
// Sicht wie beim Aufzeichnen, und der Lauf endet mit Exit-Code 0. Läuft nur in
// der zweiten Phase des Runners (PGR_OHNE_POSTGRES=1).
func TestE2EOhnePostgresExtendedReplay(t *testing.T) {
	if os.Getenv("PGR_OHNE_POSTGRES") != "1" {
		t.Skip("nur in der Phase ohne PostgreSQL")
	}
	if c, err := net.DialTimeout("tcp", os.Getenv("PGR_UPSTREAM"), 2*time.Second); err == nil {
		c.Close()
		t.Fatalf("PostgreSQL unter %s ist noch erreichbar", os.Getenv("PGR_UPSTREAM"))
	}
	daten := os.Getenv("PGR_DATEN")
	aufgezeichnet, err := os.ReadFile(filepath.Join(daten, "sicht-extended.txt"))
	if err != nil {
		t.Fatalf("Sicht der ersten Phase: %v", err)
	}
	rep := startProzess(t, "replay", "--input", filepath.Join(daten, "rec-extended.yaml"))
	got := extendedAblauf(t, rep.listen)
	rep.stop(t, 0)
	if got != string(aufgezeichnet) {
		t.Fatalf("Replay-Sicht:\n%s\n--- aufgezeichnet:\n%s\n--- stderr:\n%s", got, aufgezeichnet, rep.stderr.String())
	}
}
