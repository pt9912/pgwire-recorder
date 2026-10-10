//go:build integration

package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// anfrageE2E ist eine einfache Anfrage einer Aufzeichnung; fehler zeichnet
// an ihr nach command_complete eine error_response mit SQLSTATE XX000 auf.
type anfrageE2E struct {
	sql    string
	fehler bool
}

// aufzeichnungE2E schreibt eine Aufzeichnung mit je einer Session je Eintrag
// von sessions, Startup user postgres und database postgres, und liefert
// ihren Pfad.
func aufzeichnungE2E(t *testing.T, sessions ...[]anfrageE2E) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("format: pgwire-recorder\nversion: 1\nsessions:\n")
	for i, anfragen := range sessions {
		fmt.Fprintf(&b, "  - id: %d\n    startup:\n      user: postgres\n      database: postgres\n    interactions:\n", i+1)
		for j, a := range anfragen {
			fmt.Fprintf(&b, "      - sequence: %d\n        request:\n          type: query\n          sql: %s\n", j+1, strconv.Quote(a.sql))
			b.WriteString("        responses:\n          - type: command_complete\n            tag: \"OK\"\n")
			if a.fehler {
				b.WriteString("          - type: error_response\n            notice:\n              S: ERROR\n              C: XX000\n              M: aufgezeichnet\n")
			}
			b.WriteString("          - type: ready_for_query\n            tx_status: \"I\"\n")
		}
	}
	pfad := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(pfad, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// q ist eine Anfrage ohne aufgezeichnete error_response.
func q(sql string) anfrageE2E { return anfrageE2E{sql: sql} }

// logZeile ist eine Log-Zeile von play: Stufe, Nachricht, Code und die ganze
// Zeile.
type logZeile struct {
	stufe, msg, code, text string
}

var (
	stufeMuster = regexp.MustCompile(`level=(\S+)`)
	msgMuster   = regexp.MustCompile(`msg=("(?:[^"\\]|\\.)*"|\S+)`)
	codeMuster  = regexp.MustCompile(`code=(\S+)`)
)

// logZeilen zerlegt stderr in Log-Zeilen; jede Zeile ohne level= ist ein
// Fehler des Tests.
func logZeilen(t *testing.T, stderr string) []logZeile {
	t.Helper()
	var out []logZeile
	for _, z := range strings.Split(strings.TrimSuffix(stderr, "\n"), "\n") {
		s := stufeMuster.FindStringSubmatch(z)
		m := msgMuster.FindStringSubmatch(z)
		if s == nil || m == nil {
			t.Fatalf("keine Log-Zeile: %q\n%s", z, stderr)
		}
		l := logZeile{stufe: s[1], msg: strings.Trim(m[1], `"`), text: z}
		if c := codeMuster.FindStringSubmatch(z); c != nil {
			l.code = c[1]
		}
		out = append(out, l)
	}
	return out
}

// form ist die Folge der Zeilen als "STUFE msg" oder, mit Code, "STUFE code".
func form(zeilen []logZeile) string {
	var teile []string
	for _, z := range zeilen {
		if z.code != "" {
			teile = append(teile, z.stufe+" "+z.code)
		} else {
			teile = append(teile, z.stufe+" "+z.msg)
		}
	}
	return strings.Join(teile, " | ")
}

// enthaelt prüft, dass die Zeile alle teile enthält.
func enthaelt(t *testing.T, z logZeile, teile ...string) {
	t.Helper()
	for _, teil := range teile {
		if !strings.Contains(z.text, teil) {
			t.Fatalf("Zeile ohne %q: %s", teil, z.text)
		}
	}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-17/Boundary, LH-FA-14/Boundary — gegen
// die reale Instanz läuft play mit --continue-on-error, gesetzt über die
// Kommandozeile, die Umgebung oder den Abschnitt play:, nach jeder
// Fehlerantwort mit der nächsten Interaktion und der nächsten Session weiter
// und endet mit Exit-Code 4; je PGR-E4004 eine Zeile error mit Session,
// Interaktion und SQLSTATE in der Reihenfolge ihres Auftretens, nach der
// Zeile zum Start und vor der zum Ende. Ein ungültiger Wert in der Umgebung
// ist PGR-E2001, in der Datei PGR-E2004 (LH-FA-20.a Schritt 6, *Meldungen*,
// LH-FA-17.a).
func TestE2EPlayFortsetzung(t *testing.T) {
	conn := leereDatenbank(t, "play_fortsetzung")
	input := aufzeichnungE2E(t,
		[]anfrageE2E{q("CREATE TABLE t (n int)"), q("SELECT * FROM fehlt"), q("INSERT INTO t VALUES (1)"), q("SELECT 1/0")},
		[]anfrageE2E{q("INSERT INTO t VALUES (2)")},
	)
	dir := t.TempDir()
	basis := []string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "play_fortsetzung"}
	datei := func(wert string) string {
		pfad := filepath.Join(dir, "k-"+wert+".yaml")
		if err := os.WriteFile(pfad, []byte("play:\n  continue_on_error: "+wert+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		return "--config=" + pfad
	}
	for _, f := range []struct {
		name string
		args []string
		env  string
	}{
		{"Kommandozeile", []string{"--continue-on-error"}, ""},
		{"Umgebung", nil, "true"},
		{"Datei", []string{datei("true")}, ""},
	} {
		t.Run(f.name, func(t *testing.T) {
			t.Setenv("PGWIRE_RECORDER_CONTINUE_ON_ERROR", f.env)
			ausfuehren(t, conn, "DROP TABLE IF EXISTS t")
			_, stderr, code := starteBis(t, "", append(append([]string{}, basis...), f.args...)...)
			zeilen := logZeilen(t, stderr)
			if code != 4 || form(zeilen) != "INFO play gestartet | ERROR PGR-E4004 | ERROR PGR-E4004 | INFO play beendet" {
				t.Fatalf("Exit-Code %d, Zeilen %s, erwartet 4 mit zwei PGR-E4004:\n%s", code, form(zeilen), stderr)
			}
			enthaelt(t, zeilen[1], "Session 1, Interaktion 2", "42P01")
			enthaelt(t, zeilen[2], "Session 1, Interaktion 4", "22012")
			if got := wert(t, conn, "SELECT string_agg(n::text, ',' ORDER BY n) FROM t"); got != "1,2" {
				t.Fatalf("Wirkung %q, erwartet 1,2", got)
			}
		})
	}
	t.Setenv("PGWIRE_RECORDER_CONTINUE_ON_ERROR", "1")
	if _, stderr, code := starteBis(t, "", basis...); code != 2 || !strings.Contains(stderr, "PGR-E2001") || !strings.Contains(stderr, "PGWIRE_RECORDER_CONTINUE_ON_ERROR") {
		t.Fatalf("ungültige Umgebungsvariable: Exit-Code %d, stderr %q", code, stderr)
	}
	t.Setenv("PGWIRE_RECORDER_CONTINUE_ON_ERROR", "")
	if _, stderr, code := starteBis(t, "", append(append([]string{}, basis...), datei("1"))...); code != 2 || !strings.Contains(stderr, "PGR-E2004") || !strings.Contains(stderr, "play.continue_on_error") {
		t.Fatalf("ungültiger Wert in der Datei: Exit-Code %d, stderr %q", code, stderr)
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-14/Boundary — gegen die reale Instanz
// bricht ein Fehler, der abbricht, auch mit --continue-on-error ab: eine
// COPY-Antwort (PGR-E6001, Exit-Code 6) und eine Fehlerantwort mit FATAL
// (PGR-E4003, Exit-Code 4); seine Zeile error steht zuerst, danach die des
// früheren PGR-E4004, und keine weitere Anfrage läuft (LH-FA-20.a
// *Meldungen*, *Exit-Code*).
func TestE2EPlayFortsetzungAbbruch(t *testing.T) {
	conn := leereDatenbank(t, "play_fortsetzung_abbruch")
	ausfuehren(t, conn, "CREATE TABLE t (n int)")
	for _, f := range []struct {
		name, sql string
		exit      int
		code, ort string
	}{
		{"PGR-E6001", "COPY t FROM STDIN", 6, "PGR-E6001", "Session 1, Interaktion 2"},
		{"PGR-E4003", "SELECT pg_terminate_backend(pg_backend_pid())", 4, "PGR-E4003", "Session 1, Interaktion 2"},
	} {
		t.Run(f.name, func(t *testing.T) {
			ausfuehren(t, conn, "TRUNCATE t")
			input := aufzeichnungE2E(t,
				[]anfrageE2E{q("SELECT * FROM fehlt"), q(f.sql), q("INSERT INTO t VALUES (1)")},
				[]anfrageE2E{q("INSERT INTO t VALUES (2)")},
			)
			_, stderr, code := starteBis(t, "", "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "play_fortsetzung_abbruch", "--continue-on-error")
			zeilen := logZeilen(t, stderr)
			if want := "INFO play gestartet | ERROR " + f.code + " | ERROR PGR-E4004 | INFO play beendet"; code != f.exit || form(zeilen) != want {
				t.Fatalf("Exit-Code %d, Zeilen %s, erwartet %d mit %s:\n%s", code, form(zeilen), f.exit, want, stderr)
			}
			enthaelt(t, zeilen[1], f.ort)
			enthaelt(t, zeilen[2], "Session 1, Interaktion 1", "42P01")
			if got := wert(t, conn, "SELECT count(*) FROM t"); got != "0" {
				t.Fatalf("nach dem Abbruch eingespielt: %s Zeilen", got)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-20/Negative — gegen
// die reale Instanz ist mit --allow-recorded-errors, gesetzt über die
// Kommandozeile, die Umgebung oder den Abschnitt play:, eine Fehlerantwort
// erwartet, wenn die aufgezeichnete Interaktion eine error_response trägt,
// auch mit anderem SQLSTATE: keine Zeile error, das Einspielen läuft weiter,
// Exit-Code 0. Trägt sie keine, ist es PGR-E4004 mit Abbruch und Exit-Code 4
// (LH-FA-20.a *Interaktion*, Schritt 6).
func TestE2EPlayErwarteterFehler(t *testing.T) {
	conn := leereDatenbank(t, "play_erwartet")
	ausfuehren(t, conn, "CREATE TABLE t (n int)")
	erwartet := aufzeichnungE2E(t,
		[]anfrageE2E{{sql: "SELECT * FROM fehlt", fehler: true}, q("INSERT INTO t VALUES (1)")},
		[]anfrageE2E{q("INSERT INTO t VALUES (2)")},
	)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("play:\n  allow_recorded_errors: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, f := range []struct {
		name, dir, env string
		args           []string
	}{
		{"Kommandozeile", "", "", []string{"--allow-recorded-errors"}},
		{"Umgebung", "", "true", nil},
		{"Datei", dir, "", nil},
	} {
		t.Run(f.name, func(t *testing.T) {
			t.Setenv("PGWIRE_RECORDER_ALLOW_RECORDED_ERRORS", f.env)
			ausfuehren(t, conn, "TRUNCATE t")
			_, stderr, code := starteBis(t, f.dir, append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", erwartet, "--database", "play_erwartet"}, f.args...)...)
			if code != 0 || form(logZeilen(t, stderr)) != "INFO play gestartet | INFO play beendet" {
				t.Fatalf("Exit-Code %d, erwartet 0 ohne Zeile error:\n%s", code, stderr)
			}
			if got := wert(t, conn, "SELECT string_agg(n::text, ',' ORDER BY n) FROM t"); got != "1,2" {
				t.Fatalf("Wirkung %q, erwartet 1,2", got)
			}
		})
	}
	ausfuehren(t, conn, "TRUNCATE t")
	unerwartet := aufzeichnungE2E(t, []anfrageE2E{{sql: "SELECT 1", fehler: true}, q("SELECT * FROM fehlt"), q("INSERT INTO t VALUES (1)")})
	_, stderr, code := starteBis(t, "", "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", unerwartet, "--database", "play_erwartet", "--allow-recorded-errors")
	zeilen := logZeilen(t, stderr)
	if code != 4 || form(zeilen) != "INFO play gestartet | ERROR PGR-E4004 | INFO play beendet" {
		t.Fatalf("ohne aufgezeichnete error_response: Exit-Code %d, Zeilen %s:\n%s", code, form(zeilen), stderr)
	}
	enthaelt(t, zeilen[1], "Session 1, Interaktion 2", "42P01")
	if got := wert(t, conn, "SELECT count(*) FROM t"); got != "0" {
		t.Fatalf("nach dem Abbruch eingespielt: %s Zeilen", got)
	}
}

// warteAufSchlaf fragt über conn höchstens 30 s lang, bis in der Datenbank db
// eine Anfrage SELECT pg_sleep läuft.
func warteAufSchlaf(t *testing.T, conn *pgconn.PgConn, db string) {
	t.Helper()
	frist := time.Now().Add(30 * time.Second)
	for {
		if wert(t, conn, "SELECT count(*) FROM pg_stat_activity WHERE datname = '"+db+"' AND state = 'active' AND query LIKE 'SELECT pg_sleep(%'") == "1" {
			return
		}
		if time.Now().After(frist) {
			t.Fatalf("binnen 30 s läuft in %s kein pg_sleep", db)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-14/Boundary — gegen die reale Instanz
// beendet SIGINT oder SIGTERM play ohne --finish-session-on-interrupt nach der
// laufenden Interaktion, mit der Option aus Kommandozeile oder Umgebung nach
// der laufenden Session; eine weitere Session beginnt nicht. Der Exit-Code ist
// 0 ohne vorherigen Fehler und 4 nach einem früheren PGR-E4004 mit
// --continue-on-error, dessen Zeile error nach der Zeile zum Abbruchsignal und
// vor der zum Ende steht (LH-FA-20.a Schritt 7, *Abbruchsignal*, *Meldungen*).
func TestE2EPlayFinishSession(t *testing.T) {
	conn := leereDatenbank(t, "play_finish")
	ausfuehren(t, conn, "CREATE TABLE t (n int)")
	ohneFehler := aufzeichnungE2E(t,
		[]anfrageE2E{q("INSERT INTO t VALUES (1)"), q("SELECT pg_sleep(2)"), q("INSERT INTO t VALUES (2)")},
		[]anfrageE2E{q("INSERT INTO t VALUES (3)")},
	)
	mitFehler := aufzeichnungE2E(t,
		[]anfrageE2E{q("INSERT INTO t VALUES (1)"), q("SELECT * FROM fehlt"), q("SELECT pg_sleep(2)"), q("INSERT INTO t VALUES (2)")},
		[]anfrageE2E{q("INSERT INTO t VALUES (3)")},
	)
	for _, f := range []struct {
		name   string
		input  string
		sig    syscall.Signal
		args   []string
		env    string
		exit   int
		zeilen string
		wirkt  string
	}{
		{"ohne Option", ohneFehler, syscall.SIGINT, nil, "", 0, "INFO play gestartet | INFO Abbruchsignal, play endet vorzeitig | INFO play beendet", "1"},
		{"Kommandozeile", ohneFehler, syscall.SIGTERM, []string{"--finish-session-on-interrupt"}, "", 0, "INFO play gestartet | INFO Abbruchsignal, play endet vorzeitig | INFO play beendet", "1,2"},
		{"Umgebung nach früherem Fehler", mitFehler, syscall.SIGINT, []string{"--continue-on-error"}, "true", 4, "INFO play gestartet | INFO Abbruchsignal, play endet vorzeitig | ERROR PGR-E4004 | INFO play beendet", "1,2"},
	} {
		t.Run(f.name, func(t *testing.T) {
			t.Setenv("PGWIRE_RECORDER_FINISH_SESSION_ON_INTERRUPT", f.env)
			ausfuehren(t, conn, "TRUNCATE t")
			var stderr bytes.Buffer
			cmd := exec.Command(os.Getenv("PGR_BINARY"), append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", f.input, "--database", "play_finish"}, f.args...)...)
			cmd.Stderr = &stderr
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			beendet := make(chan error, 1)
			go func() { beendet <- cmd.Wait() }()
			t.Cleanup(func() { _ = cmd.Process.Kill() })
			warteAufSchlaf(t, conn, "play_finish")
			if err := cmd.Process.Signal(f.sig); err != nil {
				t.Fatalf("Signal: %v", err)
			}
			var code int
			select {
			case err := <-beendet:
				code = exitCodeOf(err)
			case <-time.After(30 * time.Second):
				t.Fatal("play endet binnen 30 s nach dem Signal nicht")
			}
			zeilen := logZeilen(t, stderr.String())
			if code != f.exit || form(zeilen) != f.zeilen {
				t.Fatalf("Exit-Code %d, Zeilen %s, erwartet %d mit %s:\n%s", code, form(zeilen), f.exit, f.zeilen, stderr.String())
			}
			if got := wert(t, conn, "SELECT string_agg(n::text, ',' ORDER BY n) FROM t"); got != f.wirkt {
				t.Fatalf("Wirkung %q, erwartet %s", got, f.wirkt)
			}
		})
	}
}
