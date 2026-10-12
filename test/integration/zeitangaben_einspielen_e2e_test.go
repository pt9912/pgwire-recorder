//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// markenSQL schreibt die Marke n mit der Zeit der Ausführung in die Tabelle marken.
func markenSQL(n int, zeit string) string {
	return fmt.Sprintf("INSERT INTO marken (n, t) VALUES (%d, %s)", n, zeit)
}

// ankerSQL schreibt die Marke 0 mit dem Beginn des Backends der Verbindung, auf der sie läuft:
// Beim Einspielen liegt er nach dem Beginn des Verbindungsaufbaus und vor der ersten
// Anfrage.
const ankerSQL = "INSERT INTO marken (n, t) VALUES (0, (SELECT backend_start FROM pg_stat_activity WHERE pid = pg_backend_pid()))"

// markenDatenbank legt die Datenbank name mit der Tabelle marken (n, t) an und liefert
// ihre Verbindung.
func markenDatenbank(t *testing.T, name string) *pgconn.PgConn {
	t.Helper()
	conn := leereDatenbank(t, name)
	ausfuehren(t, conn, "CREATE TABLE marken (n int, t timestamptz)")
	return conn
}

// markenAblauf führt über den Recorder unter listen in der Datenbank db drei Anfragen aus: den
// Anker (Marke 0), nach pause1 eine einfache Anfrage (Marke 1) und nach pause2 eine
// Extended-Interaktion (Marke 2, Parse, Bind, Describe, Execute, Sync in einem Aufruf).
func markenAblauf(t *testing.T, listen, db string, pause1, pause2 time.Duration) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsnDatenbank(listen, db))
	if err != nil {
		t.Fatalf("Verbindung über den Recorder: %v", err)
	}
	if _, err := conn.Exec(ctx, ankerSQL).ReadAll(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(pause1)
	if _, err := conn.Exec(ctx, markenSQL(1, "clock_timestamp()")).ReadAll(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(pause2)
	if r := conn.ExecParams(ctx, "INSERT INTO marken (n, t) VALUES ($1::int, clock_timestamp())", [][]byte{[]byte("2")}, nil, nil, nil).Read(); r.Err != nil {
		t.Fatal(r.Err)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
}

// marken liefert die Zeiten der Marken 0 bis 2 in Millisekunden (seit 1970).
func marken(t *testing.T, conn *pgconn.PgConn) [3]float64 {
	t.Helper()
	text := wert(t, conn, "SELECT string_agg(n || '=' || (extract(epoch FROM t) * 1000)::text, ' ' ORDER BY n) FROM marken")
	var out [3]float64
	felder := strings.Fields(text)
	if len(felder) != 3 {
		t.Fatalf("Marken %q, erwartet drei", text)
	}
	for i, f := range felder {
		_, v, _ := strings.Cut(f, "=")
		x, err := strconv.ParseFloat(v, 64)
		if err != nil {
			t.Fatalf("Marke %q: %v", f, err)
		}
		out[i] = x
	}
	return out
}

// toleranz ist der Abzug von den Untergrenzen der Pausen in Millisekunden: die Streuung der
// Zeit zwischen dem Senden durch play und der Ausführung auf dem Server und, bei absolute mit
// connect, der Abstand zwischen dem Beginn des Verbindungsaufbaus und dem Beginn des Backends.
const toleranz = 50.0

// Abdeckung: LH-FA-21/Happy, LH-FA-22/Happy — Abnahmeszenario 13: aufgezeichnet mit record
// --record-timing (einfache Anfragen und eine Extended-Interaktion, Pausen des Clients von
// 300 und 200 ms), in yaml und in sqlite, spielt play --keep-timing die Aufzeichnung gegen eine
// leere Datenbank ein, in jeder Betriebsart (relative, absolute mit connect und mit
// first-request): Die Zeiten der Ausführung auf dem Server halten die Untergrenze der
// aufgezeichneten Abstände ein (offset_ms aus der Datei, abzüglich der Toleranz), Exit-Code 0.
func TestE2EPlayKeepTiming(t *testing.T) {
	for _, name := range []string{"rec.yaml", "rec.sqlite"} {
		quelle := fmt.Sprintf("play_zeit_quelle_%s", filepath.Ext(name)[1:])
		markenDatenbank(t, quelle)
		output := filepath.Join(t.TempDir(), name)
		rec := startProzess(t, "record", append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-timing"}, formatArgs(name)...)...)
		markenAblauf(t, rec.listen, quelle, 300*time.Millisecond, 200*time.Millisecond)
		rec.stop(t, 0)
		o := offsetsAus(t, output)
		if len(o) != 3 || o[1]-o[0] < 300 || o[2]-o[1] < 200 {
			t.Fatalf("%s: offset_ms %v", name, o)
		}
		for i, f := range []struct {
			modus string
			args  []string
		}{
			{"relative", []string{"--timing-mode", "relative"}},
			{"absolute, connect", []string{"--timing-mode", "absolute"}},
			{"absolute, first-request", []string{"--timing-mode", "absolute", "--timing-reference", "first-request"}},
		} {
			db := fmt.Sprintf("play_zeit_%s_%d", filepath.Ext(name)[1:], i)
			conn := markenDatenbank(t, db)
			args := append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", output, "--database", db, "--keep-timing"}, f.args...)
			if stdout, stderr, code := starteBis(t, "", args...); code != 0 || stdout != "" || strings.Contains(stderr, "level=ERROR") {
				t.Fatalf("%s, %s: Exit-Code %d, stdout %q, stderr:\n%s", name, f.modus, code, stdout, stderr)
			}
			m := marken(t, conn)
			d01, d02, d12 := m[1]-m[0], m[2]-m[0], m[2]-m[1]
			switch f.modus {
			case "relative":
				if d01 < float64(o[1]-o[0])-toleranz || d12 < float64(o[2]-o[1])-toleranz {
					t.Errorf("%s, relative: Abstände %.0f und %.0f ms, erwartet mindestens %d und %d abzüglich %.0f", name, d01, d12, o[1]-o[0], o[2]-o[1], toleranz)
				}
			case "absolute, connect":
				if d01 < float64(o[1])-toleranz || d02 < float64(o[2])-toleranz {
					t.Errorf("%s, absolute, connect: Abstände %.0f und %.0f ms zum Beginn des Backends, erwartet mindestens %d und %d abzüglich %.0f", name, d01, d02, o[1], o[2], toleranz)
				}
			default:
				if d01 < float64(o[1]-o[0])-toleranz || d02 < float64(o[2]-o[0])-toleranz {
					t.Errorf("%s, absolute, first-request: Abstände %.0f und %.0f ms, erwartet mindestens %d und %d abzüglich %.0f", name, d01, d02, o[1]-o[0], o[2]-o[0], toleranz)
				}
			}
		}
	}
}

// zeitAufzeichnung schreibt eine Aufzeichnung mit einer Session aus je einer einfachen
// Anfrage je Eintrag von offsets (Marke i + 1 mit clock_timestamp()), deren offset_ms
// offsets sind, und liefert ihren Pfad.
func zeitAufzeichnung(t *testing.T, offsets ...int64) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    startup:\n      user: postgres\n      database: postgres\n    interactions:\n")
	for i, off := range offsets {
		fmt.Fprintf(&b, "      - sequence: %d\n        offset_ms: %d\n        request:\n          type: query\n          sql: %s\n", i+1, off, strconv.Quote(markenSQL(i, "clock_timestamp()")))
		b.WriteString("        responses:\n          - type: command_complete\n            tag: \"OK\"\n          - type: ready_for_query\n            tx_status: \"I\"\n")
	}
	pfad := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(pfad, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// luecke ist der Abstand der Marken 0 und 1 einer Datenbank in Millisekunden.
func luecke(t *testing.T, conn *pgconn.PgConn) float64 {
	t.Helper()
	text := wert(t, conn, "SELECT (extract(epoch FROM (SELECT t FROM marken WHERE n = 1) - (SELECT t FROM marken WHERE n = 0)) * 1000)::text")
	v, err := strconv.ParseFloat(text, 64)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// Abdeckung: LH-FA-21/Happy, LH-FA-17/Boundary — --keep-timing, --timing-mode und
// --timing-reference gelten aus der Option, aus PGWIRE_RECORDER_KEEP_TIMING,
// PGWIRE_RECORDER_TIMING_MODE und PGWIRE_RECORDER_TIMING_REFERENCE und aus den Schlüsseln
// keep_timing, timing_mode und timing_reference im Abschnitt play: (je eine Pause von 400 ms,
// Untergrenze 400 abzüglich der Toleranz); --keep-timing=false auf der Kommandozeile geht der
// Umgebungsvariable und der Datei vor, dann wartet play nicht (Pause von 3000 ms nicht
// eingehalten).
func TestE2EPlayKeepTimingQuellen(t *testing.T) {
	verzeichnis := func(inhalt string) string {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte(inhalt), 0o600); err != nil {
			t.Fatal(err)
		}
		return dir
	}
	alle := verzeichnis("play:\n  keep_timing: true\n  timing_mode: absolute\n  timing_reference: first-request\n")
	nurKeep := verzeichnis("play:\n  keep_timing: true\n")
	ohneDatei := t.TempDir()
	for i, f := range []struct {
		quelle string
		dir    string
		args   []string
		env    map[string]string
		pause  int64
		wartet bool
	}{
		{"Option", ohneDatei, []string{"--keep-timing"}, nil, 400, true},
		{"Umgebung", ohneDatei, nil, map[string]string{"PGWIRE_RECORDER_KEEP_TIMING": "true", "PGWIRE_RECORDER_TIMING_MODE": "absolute", "PGWIRE_RECORDER_TIMING_REFERENCE": "first-request"}, 400, true},
		{"Datei", alle, nil, nil, 400, true},
		{"Option false vor Datei", nurKeep, []string{"--keep-timing=false"}, nil, 3000, false},
		{"Option false vor Umgebung", ohneDatei, []string{"--keep-timing=false"}, map[string]string{"PGWIRE_RECORDER_KEEP_TIMING": "true"}, 3000, false},
	} {
		t.Run(f.quelle, func(t *testing.T) {
			db := fmt.Sprintf("play_zeit_quellen_%d", i)
			conn := markenDatenbank(t, db)
			for k, v := range f.env {
				t.Setenv(k, v)
			}
			input := zeitAufzeichnung(t, 0, f.pause)
			args := append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", db}, f.args...)
			if stdout, stderr, code := starteBis(t, f.dir, args...); code != 0 || stdout != "" || strings.Contains(stderr, "level=ERROR") {
				t.Fatalf("Exit-Code %d, stdout %q, stderr:\n%s", code, stdout, stderr)
			}
			got := luecke(t, conn)
			if f.wartet && got < float64(f.pause)-toleranz {
				t.Errorf("Abstand %.0f ms, erwartet mindestens %d abzüglich %.0f", got, f.pause, toleranz)
			}
			if !f.wartet && got >= float64(f.pause)-toleranz {
				t.Errorf("Abstand %.0f ms trotz --keep-timing=false, die Pause der Aufzeichnung ist %d", got, f.pause)
			}
		})
	}
}

// Abdeckung: LH-FA-21/Negative, LH-FA-17/Negative — ungültige Verwendung am Binary ist
// PGR-E2001 mit Exit-Code 2 und ohne Verbindung: --timing-mode oder --timing-reference ohne
// --keep-timing, --timing-reference mit relative (auch mit dem Wert des Standards), ein Wert
// außerhalb der Aufzählung (Großschreibung, Unterstrich, leer) aus Option und Umgebungsvariable,
// im Schlüssel der Datei PGR-E2004; die Meldung der Umgebungsvariable und der Datei nennt
// den Wert nicht; record und replay kennen die Optionen nicht, play nicht --record-timing.
func TestE2EPlayKeepTimingUngueltig(t *testing.T) {
	input := zeitAufzeichnung(t, 0, 10)
	play := []string{"play", "--upstream", freieAdresse(t), "--input", input}
	dir := t.TempDir()
	for _, f := range []struct {
		name string
		args []string
		env  map[string]string
		code string
		text string
	}{
		{"Modus ohne --keep-timing", []string{"--timing-mode", "absolute"}, nil, "PGR-E2001", "--timing-mode verlangt --keep-timing"},
		{"Modus relative ohne --keep-timing", []string{"--timing-mode=relative"}, nil, "PGR-E2001", "--timing-mode verlangt --keep-timing"},
		{"Bezug ohne --keep-timing", []string{"--timing-reference=first-request"}, nil, "PGR-E2001", "--timing-reference verlangt --keep-timing"},
		{"Bezug mit relative", []string{"--keep-timing", "--timing-reference=connect"}, nil, "PGR-E2001", "--timing-reference verlangt --timing-mode absolute"},
		{"Modus aus der Umgebung ohne --keep-timing", nil, map[string]string{"PGWIRE_RECORDER_TIMING_MODE": "absolute"}, "PGR-E2001", "--timing-mode verlangt --keep-timing"},
		{"Wert Absolute", []string{"--keep-timing", "--timing-mode=Absolute"}, nil, "PGR-E2001", "timing-mode"},
		{"Wert first_request", []string{"--keep-timing", "--timing-mode=absolute", "--timing-reference=first_request"}, nil, "PGR-E2001", "timing-reference"},
		{"leerer Wert", []string{"--keep-timing", "--timing-mode="}, nil, "PGR-E2001", "timing-mode"},
		{"Wert aus der Umgebung", []string{"--keep-timing"}, map[string]string{"PGWIRE_RECORDER_TIMING_MODE": "GEHEIM"}, "PGR-E2001", "PGWIRE_RECORDER_TIMING_MODE"},
		{"Keep-Timing aus der Umgebung", nil, map[string]string{"PGWIRE_RECORDER_KEEP_TIMING": "GEHEIM"}, "PGR-E2001", "PGWIRE_RECORDER_KEEP_TIMING"},
	} {
		for k, v := range f.env {
			t.Setenv(k, v)
		}
		stdout, stderr, code := starteBis(t, dir, append(append([]string{}, play...), f.args...)...)
		if code != 2 || stdout != "" || !strings.Contains(stderr, "["+f.code+"]") || !strings.Contains(stderr, f.text) || strings.Contains(stderr, "level=") || strings.Contains(stderr, "GEHEIM") {
			t.Errorf("%s: Exit-Code %d, stdout %q, stderr %q", f.name, code, stdout, stderr)
		}
		for k := range f.env {
			t.Setenv(k, "")
		}
	}
	for schluessel, wert := range map[string]string{"keep_timing": "GEHEIM", "timing_mode": "Absolute", "timing_reference": "first_request"} {
		if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("play:\n  "+schluessel+": "+wert+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		stdout, stderr, code := starteBis(t, dir, play...)
		if code != 2 || stdout != "" || !strings.Contains(stderr, "[PGR-E2004]") || !strings.Contains(stderr, "play."+schluessel) || strings.Contains(stderr, wert) {
			t.Errorf("Schlüssel %s: Exit-Code %d, stdout %q, stderr %q", schluessel, code, stdout, stderr)
		}
	}
	if err := os.Remove(filepath.Join(dir, ".pgwire-recorder.yaml")); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"record", "--listen", "127.0.0.1:0", "--upstream", "pg:5432", "--output", filepath.Join(dir, "r.yaml"), "--keep-timing"},
		{"replay", "--listen", "127.0.0.1:0", "--input", input, "--timing-mode=absolute"},
		{"play", "--upstream", "pg:5432", "--input", input, "--record-timing"},
	} {
		stdout, stderr, code := starteBis(t, dir, args...)
		if code != 2 || stdout != "" || !strings.Contains(stderr, "[PGR-E2001]") {
			t.Errorf("%s: Exit-Code %d, stdout %q, stderr %q", args[0], code, stdout, stderr)
		}
	}
}

// Abdeckung: LH-FA-21/Negative — fehlt mit --keep-timing irgendeiner Interaktion einer Session
// offset_ms, ist das PGR-E2003 mit Exit-Code 2, die Zeile beim Prozessende ohne Log-Zeile,
// vor der ersten Verbindung (die Adresse des Servers ist nicht erreichbar, sonst wäre es
// PGR-E4002) und nennt Session und Interaktion der ersten ohne Wert, keinen Wert; ohne
// --keep-timing spielt dieselbe Aufzeichnung ein.
func TestE2EPlayKeepTimingOhneOffset(t *testing.T) {
	conn := markenDatenbank(t, "play_zeit_ohne_offset")
	input := zeitAufzeichnung(t, 7001, 7002)
	text, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	ohne := regexp.MustCompile(`(?m)^        offset_ms: 7002\n`).ReplaceAllString(string(text), "")
	if ohne == string(text) {
		t.Fatal("offset_ms der zweiten Interaktion nicht entfernt")
	}
	if err := os.WriteFile(input, []byte(ohne), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := starteBis(t, "", "play", "--upstream", freieAdresse(t), "--input", input, "--keep-timing")
	if code != 2 || stdout != "" || !strings.HasPrefix(stderr, "Konfiguration [PGR-E2003]: Session 1, Interaktion 2 ") || strings.Count(stderr, "\n") != 1 || strings.Contains(stderr, "7001") || strings.Contains(stderr, "level=") {
		t.Errorf("Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	if stdout, stderr, code := starteBis(t, "", "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "play_zeit_ohne_offset"); code != 0 || stdout != "" || strings.Contains(stderr, "level=ERROR") {
		t.Errorf("ohne --keep-timing: Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	if got := wert(t, conn, "SELECT count(*) FROM marken"); got != "2" {
		t.Errorf("ohne --keep-timing eingespielt: %s Zeilen, erwartet 2", got)
	}
}

// Abdeckung: LH-FA-01/Happy, LH-FA-17/Happy — die Hilfe von play nennt --keep-timing,
// --timing-mode, --timing-reference und ihre Umgebungsvariablen; config show nennt die
// Schlüssel keep_timing, timing_mode und timing_reference der Datei und, bei gesetzten
// Umgebungsvariablen, deren Namen ohne Wert.
func TestE2EPlayKeepTimingHilfeUndConfigShow(t *testing.T) {
	dir := t.TempDir()
	stdout, stderr, code := starteBis(t, dir, "play", "--help")
	if code != 0 || stderr != "" {
		t.Fatalf("play --help: Exit-Code %d, stderr %q", code, stderr)
	}
	for _, teil := range []string{"--keep-timing", "--timing-mode", "--timing-reference", "PGWIRE_RECORDER_KEEP_TIMING", "PGWIRE_RECORDER_TIMING_MODE", "PGWIRE_RECORDER_TIMING_REFERENCE"} {
		if !strings.Contains(stdout, teil) {
			t.Errorf("play --help ohne %q", teil)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("play:\n  keep_timing: true\n  timing_mode: absolute\n  timing_reference: first-request\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PGWIRE_RECORDER_TIMING_MODE", "relative")
	stdout, stderr, code = starteBis(t, dir, "config", "show")
	if code != 0 || stderr != "" || strings.Contains(stdout, "PGWIRE_RECORDER_TIMING_MODE=") || !strings.Contains(stdout, "PGWIRE_RECORDER_TIMING_MODE\n") {
		t.Fatalf("config show: Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	for _, zeile := range []string{"  keep_timing: true\n", "  timing_mode: absolute\n", "  timing_reference: first-request\n"} {
		if !strings.Contains(stdout, zeile) {
			t.Errorf("config show ohne %q:\n%s", zeile, stdout)
		}
	}
}

// Abdeckung: LH-FA-21/Boundary, LH-FA-07/Boundary — replay ignoriert Zeitangaben immer: Eine
// mit --record-timing aufgezeichnete Aufzeichnung, deren offset_ms alle auf eine Stunde
// stehen, beantwortet replay mit derselben Sicht des Clients wie beim Aufzeichnen, ohne eine
// Pause einzuhalten (die Anfragen haben 30 s).
func TestE2EReplayIgnoriertZeitangaben(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-timing")
	queries := []string{"SELECT 1 AS eins;", "SELECT 'a' AS text;"}
	aufgezeichnet := beobachten(t, rec.listen, queries...)
	rec.stop(t, 0)
	text := lies(t, output)
	stunde := regexp.MustCompile(`(?m)^(\s+offset_ms: )\d+$`).ReplaceAllString(text, "${1}3600000")
	if stunde == text || strings.Count(stunde, "offset_ms: 3600000") != 2 {
		t.Fatalf("offset_ms nicht auf eine Stunde gesetzt:\n%s", stunde)
	}
	if err := os.WriteFile(output, []byte(stunde), 0o644); err != nil {
		t.Fatal(err)
	}
	rep := startProzess(t, "replay", "--input", output)
	got := beobachten(t, rep.listen, queries...)
	rep.stop(t, 0)
	if got != aufgezeichnet {
		t.Fatalf("Sicht beim Replay:\n%s\n--- aufgezeichnet:\n%s", got, aufgezeichnet)
	}
}
