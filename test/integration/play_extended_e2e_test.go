//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// dsnDatenbank ist die Verbindungs-URL zu address und der Datenbank db.
func dsnDatenbank(address, db string) string {
	return fmt.Sprintf("postgres://postgres@%s/%s?sslmode=disable&connect_timeout=10", address, db)
}

// starteBisFrist führt das Binary mit args aus und wartet höchstens frist auf
// sein Ende; läuft die Frist ab, wird der Test rot und nennt das ausgebliebene
// Ende. Es liefert stdout, stderr und den Exit-Code.
func starteBisFrist(t *testing.T, frist time.Duration, args ...string) (string, string, int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), frist)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Getenv("PGR_BINARY"), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if ctx.Err() != nil {
		t.Fatalf("%q endet nicht binnen %s\n%s", args, frist, stderr.String())
	}
	return stdout.String(), stderr.String(), exitCodeOf(err)
}

// personen fasst die Zeilen der Tabelle personen einer Datenbank in einen Text:
// je Zeile id, name, notiz (NULL als <NULL>) und zahl, nach id geordnet.
func personen(t *testing.T, conn *pgconn.PgConn) string {
	t.Helper()
	return wert(t, conn, "SELECT coalesce(string_agg(id || ':' || name || ':' || coalesce(notiz, '<NULL>') || ':' || zahl, '|' ORDER BY id), '') FROM personen")
}

// nimmAufEinspielAblauf zeichnet über record gegen die Datenbank quelle einen
// Ablauf in drei Sessions auf und liefert den Pfad der Aufzeichnung: pgx im
// Standardmodus (Extended Query) mit DDL und DML mit Parametern (NULL, leer,
// Unicode, Anführungszeichen), ein Batch (eine Sync-Gruppe), eine einfache
// Anfrage; eine Pipeline mit einer Flush-Gruppe vor zwei ohne Warten gesendeten
// Ausführungen und einem Sync; eine Transaktion.
func nimmAufEinspielAblauf(t *testing.T, quelle string) string {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsnDatenbank(rec.listen, quelle))
	if err != nil {
		t.Fatal(err)
	}
	einfuegen := "INSERT INTO personen (id, name, notiz, zahl) VALUES ($1, $2, $3, $4)"
	for _, sql := range []struct {
		sql  string
		args []any
	}{
		{"CREATE TABLE personen (id int PRIMARY KEY, name text, notiz text, zahl int)", nil},
		{einfuegen, []any{1, "Ada", nil, 10}},
		{einfuegen, []any{2, "Bob", "", 20}},
		{einfuegen, []any{3, "Cäsar 日本", "mit 'Anführungszeichen'", 30}},
	} {
		if _, err := conn.Exec(ctx, sql.sql, sql.args...); err != nil {
			t.Fatalf("%s: %v", sql.sql, err)
		}
	}
	batch := &pgx.Batch{}
	batch.Queue(einfuegen, 4, "Dora", "batch", 40)
	batch.Queue("UPDATE personen SET zahl = zahl + $1 WHERE id = $2", 5, 2)
	if err := conn.SendBatch(ctx, batch).Close(); err != nil {
		t.Fatalf("Batch: %v", err)
	}
	if _, err := conn.Exec(ctx, "UPDATE personen SET notiz = 'einfach' WHERE id = 1", pgx.QueryExecModeSimpleProtocol); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}

	pc, err := pgconn.Connect(ctx, dsnDatenbank(rec.listen, quelle))
	if err != nil {
		t.Fatal(err)
	}
	p := pc.StartPipeline(ctx)
	p.SendPrepare("ins", einfuegen, nil)
	p.SendFlushRequest()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	if _, err := p.GetResults(); err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	p.SendQueryPrepared("ins", [][]byte{[]byte("5"), []byte("Eva"), nil, []byte("50")}, nil, nil)
	p.SendQueryPrepared("ins", [][]byte{[]byte("6"), []byte(""), []byte("x"), []byte("60")}, nil, nil)
	p.SendDeallocate("ins")
	p.SendPipelineSync()
	if err := p.Flush(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		r, err := p.GetResults()
		if err != nil {
			t.Fatalf("Pipeline, Ergebnis %d: %v", i+1, err)
		}
		if rr, ok := r.(*pgconn.ResultReader); ok {
			if res := rr.Read(); res.Err != nil {
				t.Fatalf("Pipeline, Ergebnis %d: %v", i+1, res.Err)
			}
		}
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	_ = pc.Close(ctx)

	conn, err = pgx.Connect(ctx, dsnDatenbank(rec.listen, quelle))
	if err != nil {
		t.Fatal(err)
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, einfuegen, 7, "Fred", nil, 70); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
	rec.stop(t, 0)
	return output
}

// Abdeckung: LH-FA-20/Happy, LH-FA-18/Happy — play spielt eine Aufzeichnung mit
// Extended-Interaktionen (DDL und DML mit Parametern einschließlich NULL und
// leer, Gruppen mit Sync und mit Flush, gemischt mit einer einfachen Anfrage,
// über drei Sessions) gegen eine leere Datenbank ohne Passwort und TLS ein
// (--database vor der Aufzeichnung): Die Datenbank enthält danach dieselben
// Zeilen wie die, gegen die aufgezeichnet wurde, NULL und leerer Wert
// unterschieden; Exit-Code 0, keine Log-Zeile der Stufe error (Abnahmeszenario
// 12, LH-FA-20.a Schritt 3).
func TestE2EPlayExtended(t *testing.T) {
	quelle := leereDatenbank(t, "play_ext_quelle")
	input := nimmAufEinspielAblauf(t, "play_ext_quelle")
	text := lies(t, input)
	inReihe(t, text, "type: extended", "- type: flush", "- type: sync", "type: query")
	if c := strings.Count(text, "- id: "); c != 3 {
		t.Fatalf("%d Sessions aufgezeichnet statt drei:\n%s", c, text)
	}
	erwartet := personen(t, quelle)
	if !strings.Contains(erwartet, "1:Ada:einfach:10") || !strings.Contains(erwartet, "2:Bob::25") || !strings.Contains(erwartet, "5:Eva:<NULL>:50") || !strings.Contains(erwartet, "6::x:60") || !strings.Contains(erwartet, "3:Cäsar 日本:mit 'Anführungszeichen':30") || !strings.Contains(erwartet, "7:Fred:<NULL>:70") {
		t.Fatalf("Aufzeichnung gegen die Quelle: %q", erwartet)
	}

	ziel := leereDatenbank(t, "play_ext_ziel")
	stdout, stderr, code := starteBisFrist(t, 30*time.Second, "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "play_ext_ziel")
	if code != 0 || stdout != "" || strings.Contains(stderr, "level=ERROR") {
		t.Fatalf("Exit-Code %d, stdout %q, stderr:\n%s", code, stdout, stderr)
	}
	if got := personen(t, ziel); got != erwartet {
		t.Fatalf("Wirkung\n%s\nerwartet\n%s", got, erwartet)
	}
}

// nimmAufEinspielFehler zeichnet über record gegen die Datenbank quelle einen
// Ablauf auf, dessen Batch (eine Sync-Gruppe) in der zweiten Anweisung mit
// einem Fehler endet (Division durch null), gefolgt von einer weiteren
// Anfrage, und liefert den Pfad der Aufzeichnung.
func nimmAufEinspielFehler(t *testing.T, quelle string) string {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsnDatenbank(rec.listen, quelle))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(ctx, "CREATE TABLE ergebnis (n int)"); err != nil {
		t.Fatal(err)
	}
	batch := &pgx.Batch{}
	batch.Queue("INSERT INTO ergebnis VALUES ($1)", 1)
	batch.Queue("SELECT 1/$1::int", 0)
	batch.Queue("INSERT INTO ergebnis VALUES ($1)", 3)
	if err := conn.SendBatch(ctx, batch).Close(); err == nil {
		t.Fatal("der Batch endet ohne Fehler")
	}
	if _, err := conn.Exec(ctx, "INSERT INTO ergebnis VALUES ($1)", 4); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
	rec.stop(t, 0)
	return output
}

// Abdeckung: LH-FA-20/Negative — endet eine Extended-Interaktion
// gegen die reale Instanz mit einer Fehlerantwort, bricht play ohne Optionen ab
// (PGR-E4004 mit SQLSTATE, Exit-Code 4, die weitere Anfrage nicht eingespielt);
// mit --continue-on-error liest es nur noch bis zum ReadyForQuery, spielt die
// weitere Anfrage ein und endet mit Exit-Code 4; mit --allow-recorded-errors
// gilt der Fehler als erwartet, weil die Aufzeichnung ihn trägt (Exit-Code 0,
// keine Meldung) (LH-FA-20.a Schritt 6, *Interaktion*).
func TestE2EPlayExtendedFehler(t *testing.T) {
	leereDatenbank(t, "play_extf_quelle")
	input := nimmAufEinspielFehler(t, "play_extf_quelle")
	if !strings.Contains(lies(t, input), "- type: error_response") {
		t.Fatalf("die Aufzeichnung trägt keinen Fehler:\n%s", lies(t, input))
	}
	for _, f := range []struct {
		name  string
		db    string
		opt   []string
		exit  int
		meldg int
		zeile string
	}{
		{"ohne Optionen", "play_extf_a", nil, 4, 1, ""},
		{"--continue-on-error", "play_extf_b", []string{"--continue-on-error"}, 4, 1, "4"},
		{"--allow-recorded-errors", "play_extf_c", []string{"--allow-recorded-errors"}, 0, 0, "4"},
	} {
		t.Run(f.name, func(t *testing.T) {
			conn := leereDatenbank(t, f.db)
			args := append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", f.db}, f.opt...)
			_, stderr, code := starteBisFrist(t, 30*time.Second, args...)
			if code != f.exit || strings.Count(stderr, "code=PGR-E4004") != f.meldg || (f.meldg > 0 && !strings.Contains(stderr, "22012")) {
				t.Fatalf("Exit-Code %d, %d-mal PGR-E4004, erwartet %d mit %d, stderr:\n%s", code, strings.Count(stderr, "code=PGR-E4004"), f.exit, f.meldg, stderr)
			}
			got := wert(t, conn, "SELECT coalesce(string_agg(n::text, ',' ORDER BY n), '') FROM ergebnis")
			if got != f.zeile {
				t.Fatalf("Zeilen %q, erwartet %q", got, f.zeile)
			}
		})
	}
}

// gegendruckAufzeichnung schreibt eine Aufzeichnung mit einer Session aus zwei
// Extended-Interaktionen und liefert ihren Pfad. Die erste ist eine Gruppe mit
// Sync: eine Anfrage mit 32 MB Ausgabe, danach eine Anfrage mit 16 MB
// Parameter. Die zweite hat eine Gruppe mit Flush, deren Ausgabe (32 MB)
// verzögert eintrifft, und eine Gruppe mit Sync.
func gegendruckAufzeichnung(t *testing.T) string {
	t.Helper()
	gross := strings.Repeat("y", parameterMB<<20)
	var b strings.Builder
	b.WriteString("format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    startup:\n      user: postgres\n      database: postgres\n    interactions:\n")
	parse := func(sql, typen string) string {
		return fmt.Sprintf("              - type: parse\n                statement: \"\"\n                sql: %q\n                param_types: %s\n", sql, typen)
	}
	bind := func(werte ...string) string {
		s := "              - type: bind\n                portal: \"\"\n                statement: \"\"\n                param_formats: []\n"
		if len(werte) == 0 {
			s += "                params: []\n"
		} else {
			s += "                params:\n"
			for _, w := range werte {
				s += fmt.Sprintf("                  - text: %q\n", w)
			}
		}
		return s + "                result_formats: []\n"
	}
	const execute = "              - type: execute\n                portal: \"\"\n                max_rows: 0\n"
	const flush = "              - type: flush\n"
	const sync = "              - type: sync\n"
	const bereit = "            server:\n              - type: ready_for_query\n                tx_status: \"I\"\n"
	b.WriteString("      - sequence: 1\n        type: extended\n        groups:\n          - client:\n")
	b.WriteString(parse("SELECT repeat('x', 1000000) FROM generate_series(1, $1::int)", "[23]") + bind(fmt.Sprint(zeilenMB)) + execute)
	b.WriteString(parse("INSERT INTO gross (n) VALUES (length($1::text))", "[25]") + bind(gross) + execute + sync + bereit)
	b.WriteString("      - sequence: 2\n        type: extended\n        groups:\n          - client:\n")
	b.WriteString(parse("WITH s AS MATERIALIZED (SELECT pg_sleep(1)) SELECT repeat('x', 1000000) FROM s, generate_series(1, 32)", "[]") + bind() + execute + flush + "            server: []\n")
	b.WriteString("          - client:\n")
	b.WriteString(parse("INSERT INTO gross (n) VALUES ($1::int)", "[23]") + bind("7") + execute + sync + bereit)
	pfad := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(pfad, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-18/Boundary — eine Gruppe, deren
// Nachrichten (16 MB Parameter) und Antworten (32 MB Ausgabe) größer sind als
// die Puffer der Verbindung, verklemmt gegen die reale Instanz nicht: play sendet
// und liest unabhängig voneinander und endet binnen 60 s mit Exit-Code 0; eine
// Gruppe mit Flush, deren große Ausgabe verzögert eintrifft, wird abgewartet, bevor
// die nächste Gruppe geht (LH-FA-20.a *Gruppen*).
func TestE2EPlayExtendedGegendruck(t *testing.T) {
	conn := leereDatenbank(t, "play_ext_gross")
	ausfuehren(t, conn, "CREATE TABLE gross (n bigint)")
	input := gegendruckAufzeichnung(t)
	stdout, stderr, code := starteBisFrist(t, 60*time.Second, "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "play_ext_gross")
	if code != 0 || stdout != "" || strings.Contains(stderr, "level=ERROR") {
		t.Fatalf("Exit-Code %d, stdout %q, stderr:\n%s", code, stdout, stderr)
	}
	if got := wert(t, conn, "SELECT string_agg(n::text, ',' ORDER BY n DESC) FROM gross"); got != fmt.Sprintf("%d,7", parameterMB<<20) {
		t.Fatalf("Zeilen %q, erwartet die Länge des Parameters und 7", got)
	}
}
