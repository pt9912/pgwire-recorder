//go:build integration

package integration_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// inReihe prüft, dass die Zeilen in dieser Reihenfolge im Text stehen.
func inReihe(t *testing.T, text string, zeilen ...string) {
	t.Helper()
	pos := 0
	for _, z := range zeilen {
		i := strings.Index(text[pos:], z)
		if i < 0 {
			t.Fatalf("%q fehlt oder steht nicht in der Reihenfolge %q:\n%s", z, zeilen, text)
		}
		pos += i + len(z)
	}
}

// laedt startet replay mit der Aufzeichnung: es lauscht nur, wenn der Leser sie
// geladen und jede Interaktion Validate bestanden hat (sonst Exit-Code 3).
func laedt(t *testing.T, input string) {
	t.Helper()
	rep := startProzess(t, "replay", "--input", input)
	rep.stop(t, 0)
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: pgx in seinem Standardmodus
// (Prepared Statements über das Extended Query Protocol) führt über `record`
// gegen eine reale PostgreSQL-Instanz dasselbe Statement mit verschiedenen
// Werten, einen NULL-Parameter, einen Fehler und danach eine einfache Anfrage
// aus und erhält die Ergebnisse des Servers; die Aufzeichnung trägt die
// Extended-Interaktionen mit Parse, Bind, Execute und Sync in Reihenfolge,
// lädt mit dem Leser, und der Lauf endet mit Exit-Code 0.
func TestE2ERecordExtendedPgx(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatalf("Verbindung über den Recorder: %v", err)
	}
	for _, wert := range []string{"eins", "zwei", "drei"} {
		var got string
		if err := conn.QueryRow(ctx, "SELECT $1::text || '!' AS t", wert).Scan(&got); err != nil || got != wert+"!" {
			t.Fatalf("%s: %q, %v", wert, got, err)
		}
	}
	var istNull bool
	if err := conn.QueryRow(ctx, "SELECT $1::text IS NULL AS n", nil).Scan(&istNull); err != nil || !istNull {
		t.Fatalf("NULL-Parameter: %v, %v", istNull, err)
	}
	var pgErr *pgconn.PgError
	if _, err := conn.Exec(ctx, "SELECT 1/$1::int", 0); !errors.As(err, &pgErr) || pgErr.Code != "22012" {
		t.Fatalf("erwartet 22012, erhalten %v", err)
	}
	var n int
	if err := conn.QueryRow(ctx, "SELECT 7", pgx.QueryExecModeSimpleProtocol).Scan(&n); err != nil || n != 7 {
		t.Fatalf("einfache Anfrage: %d, %v", n, err)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
	rec.stop(t, 0)

	text := lies(t, output)
	if c := strings.Count(text, "sql: SELECT $1::text || '!' AS t"); c != 1 {
		t.Errorf("Parse des wiederholten Statements %d-mal statt einmal:\n%s", c, text)
	}
	inReihe(t, text,
		"type: extended", "- type: parse", "sql: SELECT $1::text || '!' AS t", "- type: sync", "- type: parse_complete", "- type: parameter_description", "- type: ready_for_query",
		"- type: bind", "- text: eins", "- type: execute", "- type: sync", "- text: eins!", "- type: ready_for_query",
		"- text: zwei", "- text: zwei!",
		"- text: drei", "- text: drei!",
		`"null": true`, "- type: ready_for_query",
		"- type: error_response", "22012", "- type: ready_for_query",
		"type: query", "sql: SELECT 7",
	)
	laedt(t, output)
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: eine Pipeline aus einer
// Flush-Gruppe, auf deren Antwort der Client wartet, und danach zwei ohne
// Warten gesendeten Sync-Gruppen (zwei Ausführungen und Close, dann ein
// Fehler) liefert dem Client die Ergebnisse des Servers; die Aufzeichnung
// trägt die Gruppen mit ihren Server-Nachrichten in Reihenfolge, das Ende einer
// Interaktion je Sync, und lädt mit dem Leser.
func TestE2ERecordExtendedPipeline(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	p := conn.StartPipeline(ctx)
	p.SendPrepare("s1", "SELECT $1::text AS t", nil)
	p.SendFlushRequest()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	if r, err := p.GetResults(); err != nil {
		t.Fatalf("Prepare: %v", err)
	} else if sd, ok := r.(*pgconn.StatementDescription); !ok || len(sd.ParamOIDs) != 1 {
		t.Fatalf("Prepare: %#v", r)
	}
	p.SendQueryPrepared("s1", [][]byte{[]byte("a")}, nil, nil)
	p.SendQueryPrepared("s1", [][]byte{nil}, nil, nil)
	p.SendDeallocate("s1")
	p.SendPipelineSync()
	p.SendQueryParams("SELECT 1/0", nil, nil, nil, nil)
	p.SendPipelineSync()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	var zeilen []string
	for i := 0; i < 2; i++ {
		r, err := p.GetResults()
		if err != nil {
			t.Fatalf("Ausführung %d: %v", i+1, err)
		}
		res := r.(*pgconn.ResultReader).Read()
		if res.Err != nil || len(res.Rows) != 1 {
			t.Fatalf("Ausführung %d: %#v", i+1, res)
		}
		zeilen = append(zeilen, string(res.Rows[0][0]))
	}
	if zeilen[0] != "a" || zeilen[1] != "" {
		t.Fatalf("Zeilen: %q", zeilen)
	}
	for _, schritt := range []string{"Close", "Sync"} {
		if _, err := p.GetResults(); err != nil {
			t.Fatalf("%s: %v", schritt, err)
		}
	}
	var pgErr *pgconn.PgError
	if r, err := p.GetResults(); err == nil {
		if res := r.(*pgconn.ResultReader).Read(); !errors.As(res.Err, &pgErr) {
			t.Fatalf("erwartet 22012, erhalten %#v", res)
		}
	} else if !errors.As(err, &pgErr) {
		t.Fatalf("erwartet 22012, erhalten %v", err)
	}
	if pgErr.Code != "22012" {
		t.Fatalf("erwartet 22012, erhalten %s", pgErr.Code)
	}
	if _, err := p.GetResults(); err != nil {
		t.Fatalf("Sync nach dem Fehler: %v", err)
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	_ = conn.Close(ctx)
	rec.stop(t, 0)

	text := lies(t, output)
	inReihe(t, text,
		"- sequence: 1", "type: extended", "groups:",
		"- client:", "- type: parse", "statement: s1", "- type: describe", "target: statement", "- type: flush",
		"server:", "- type: parse_complete", "- type: parameter_description", "- type: row_description",
		"- client:", "- type: bind", "- text: a", "- type: execute", "- type: bind", `"null": true`, "- type: execute", "- type: close", "- type: sync",
		"server:", "- type: bind_complete", "- text: a", "- type: command_complete", "- type: bind_complete", "- type: command_complete", "- type: close_complete", "- type: ready_for_query",
		"- sequence: 2", "type: extended", "- type: parse", "sql: SELECT 1/0", "- type: sync",
		"server:", "- type: error_response", "22012", "- type: ready_for_query",
	)
	if c := strings.Count(text, "- client:"); c != 3 {
		t.Errorf("%d Gruppen statt drei:\n%s", c, text)
	}
	laedt(t, output)
}

// Abdeckung: LH-FA-13/Negative — endet die Client-Verbindung
// nach einem Flush ohne Sync, ist das ein unerwartetes Verbindungsende
// (PGR-E4003): die unvollständige Extended-Interaktion steht nicht in der
// Aufzeichnung, die vorherige Interaktion der Session bleibt, und der Lauf endet
// mit dem Exit-Code der Klasse Netzwerk (4).
func TestE2ERecordExtendedAbbruch(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
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
	_ = conn.Conn().Close()
	time.Sleep(200 * time.Millisecond)
	rec.stop(t, 4)

	text := lies(t, output)
	if !strings.Contains(text, "sql: SELECT 1;") || strings.Contains(text, "type: extended") || strings.Contains(text, "SELECT 2") {
		t.Fatalf("Aufzeichnung nach dem Abbruch:\n%s", text)
	}
	if !strings.Contains(rec.stderr.String(), "PGR-E4003") {
		t.Fatalf("PGR-E4003 fehlt in der Ausgabe:\n%s", rec.stderr.String())
	}
	laedt(t, output)
}

// Größen der Gegendruck-Tests: Ausgabe und Parameter übersteigen die Puffer
// der Verbindungen weit.
const (
	zeilenMB    = 32
	parameterMB = 16
)

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: Gegendruck in beiden Richtungen
// (gegendruckAblauf) läuft über `record` binnen des Zeitlimits durch, und der
// Lauf endet mit Exit-Code 0.
func TestE2ERecordExtendedGegendruck(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	gegendruckAblauf(t, rec.listen)
	rec.stop(t, 0)
}

// gegendruckAblauf: pgx im Standardmodus sendet einen Batch (eine Sync-Gruppe)
// aus einer Anfrage mit 32 MB Ausgabe und einer mit 16 MB Parameter, und eine
// Pipeline sendet nach einer Flush-Gruppe mit verzögerter großer Ausgabe ohne
// Warten eine Sync-Gruppe mit großem Parameter; der Client erhält jedes
// Ergebnis binnen des Zeitlimits.
func gegendruckAblauf(t *testing.T, listen string) {
	t.Helper()
	gross := strings.Repeat("y", parameterMB<<20)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn(listen))
	if err != nil {
		t.Fatal(err)
	}
	batch := &pgx.Batch{}
	batch.Queue("SELECT repeat('x', 1000000) FROM generate_series(1, $1::int)", zeilenMB)
	batch.Queue("SELECT length($1::text)", gross)
	br := conn.SendBatch(ctx, batch)
	rows, err := br.Query()
	if err != nil {
		t.Fatalf("Batch, Ausgabe: %v", err)
	}
	n := 0
	for rows.Next() {
		n++
	}
	if rows.Err() != nil || n != zeilenMB {
		t.Fatalf("Batch, Ausgabe: %d Zeilen, %v", n, rows.Err())
	}
	var laenge int
	if err := br.QueryRow().Scan(&laenge); err != nil || laenge != parameterMB<<20 {
		t.Fatalf("Batch, Parameter: %d, %v", laenge, err)
	}
	if err := br.Close(); err != nil {
		t.Fatal(err)
	}
	_ = conn.Close(ctx)

	pc, err := pgconn.Connect(ctx, dsn(listen))
	if err != nil {
		t.Fatal(err)
	}
	p := pc.StartPipeline(ctx)
	p.SendQueryParams("WITH s AS MATERIALIZED (SELECT pg_sleep(1)) SELECT repeat('x', 1000000) FROM s, generate_series(1, 32)", nil, nil, nil, nil)
	p.SendFlushRequest()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	p.SendQueryParams("SELECT length($1::text)", [][]byte{[]byte(gross)}, nil, nil, nil)
	p.SendPipelineSync()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		r, err := p.GetResults()
		if err != nil {
			t.Fatalf("Pipeline, Ergebnis %d: %v", i+1, err)
		}
		if res := r.(*pgconn.ResultReader).Read(); res.Err != nil {
			t.Fatalf("Pipeline, Ergebnis %d: %v", i+1, res.Err)
		}
	}
	if _, err := p.GetResults(); err != nil {
		t.Fatalf("Pipeline, Sync: %v", err)
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	_ = pc.Close(ctx)
}

// Abdeckung: LH-FA-13/Boundary — liest ein Client die Antworten auf eine große
// Gruppe nicht und schließt dann die Verbindung, endet die Session, obwohl
// der Recorder beim Senden an den Server und beim Schreiben an den Client
// stand; der Lauf endet danach auf SIGTERM mit dem Exit-Code der Klasse
// Netzwerk (4), weil die Interaktion unvollständig abbrach.
func TestE2ERecordExtendedSigtermNachBlockade(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	pc, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	p := pc.StartPipeline(ctx)
	p.SendQueryParams("SELECT repeat('x', 1000000) FROM generate_series(1, 64)", nil, nil, nil, nil)
	p.SendQueryParams("SELECT length($1::text)", [][]byte{[]byte(strings.Repeat("y", parameterMB<<20))}, nil, nil, nil)
	p.SendPipelineSync()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	// Die Antworten bleiben ungelesen; nach einer Weile staut sich alles.
	time.Sleep(2 * time.Second)
	_ = pc.Conn().Close()
	time.Sleep(500 * time.Millisecond)
	rec.stop(t, 4)
	if !strings.Contains(rec.stderr.String(), "PGR-E4003") {
		t.Fatalf("PGR-E4003 fehlt:\n%s", rec.stderr.String())
	}
}

// Abdeckung: LH-FA-13/Boundary — ein Client, der ohne Pause pipelinet (zwei
// Interaktionen unterwegs, nach jedem ReadyForQuery die nächste), hält das
// Herunterfahren nicht auf: Nach SIGTERM beginnt keine neue Interaktion, der
// Prozess endet binnen 5 s mit Exit-Code 0, und die Aufzeichnung trägt höchstens
// die zwei beim Signal unterwegs gewesenen Interaktionen mehr, als der Client bis
// dahin abgeschlossen hatte.
func TestE2ERecordExtendedSigtermBeimPipelining(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	pc, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := pc.Close(context.Background()); err != nil {
			t.Error(err)
		}
	}()
	p := pc.StartPipeline(ctx)
	senden := func() error {
		p.SendQueryParams("SELECT pg_sleep(0.05)", nil, nil, nil, nil)
		p.SendPipelineSync()
		return p.Flush()
	}
	var fertig atomic.Int64
	go func() {
		if senden() != nil {
			return
		}
		if senden() != nil {
			return
		}
		for {
			for i := 0; i < 2; i++ {
				r, err := p.GetResults()
				if err != nil {
					return
				}
				if rr, ok := r.(*pgconn.ResultReader); ok {
					if rr.Read().Err != nil {
						return
					}
				}
			}
			fertig.Add(1)
			if senden() != nil {
				return
			}
		}
	}()
	time.Sleep(time.Second)
	if err := rec.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	vorSignal := fertig.Load()
	rec.warteEnde(t, 5*time.Second, "Recorder endet nicht binnen 5 s nach SIGTERM, obwohl der Client weiter pipelinet")
	rec.pruefeExit(t, 0)
	if vorSignal < 3 {
		t.Fatalf("vor dem Signal nur %d Interaktionen abgeschlossen", vorSignal)
	}
	if n := int64(strings.Count(lies(t, output), "type: extended")); n > vorSignal+2 {
		t.Fatalf("%d Interaktionen aufgezeichnet, vor dem Signal %d abgeschlossen", n, vorSignal)
	}
}
