//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// einspielAufzeichnung schreibt eine Aufzeichnung mit je einer Session je
// Eintrag von sessions, deren Anfragen einfache Anfragen sind, mit dem
// Startup user postgres und database postgres, und liefert ihren Pfad.
func einspielAufzeichnung(t *testing.T, sessions ...[]string) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("format: pgwire-recorder\nversion: 1\nsessions:\n")
	for i, anfragen := range sessions {
		fmt.Fprintf(&b, "  - id: %d\n    startup:\n      user: postgres\n      database: postgres\n    interactions:\n", i+1)
		for j, sql := range anfragen {
			fmt.Fprintf(&b, "      - sequence: %d\n        request:\n          type: query\n          sql: %s\n", j+1, strconv.Quote(sql))
			b.WriteString("        responses:\n          - type: command_complete\n            tag: \"OK\"\n          - type: ready_for_query\n            tx_status: \"I\"\n")
		}
	}
	pfad := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(pfad, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// leereDatenbank legt eine neue Datenbank name auf der Instanz an und
// liefert eine Verbindung zu ihr, die der Cleanup schließt.
func leereDatenbank(t *testing.T, name string) *pgconn.PgConn {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	admin, err := pgconn.Connect(ctx, dsn(os.Getenv("PGR_UPSTREAM")))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = admin.Close(ctx) }()
	for _, sql := range []string{"DROP DATABASE IF EXISTS " + name, "CREATE DATABASE " + name} {
		if _, err := admin.Exec(ctx, sql).ReadAll(); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
	}
	conn, err := pgconn.Connect(ctx, fmt.Sprintf("postgres://postgres@%s/%s?sslmode=disable&connect_timeout=10", os.Getenv("PGR_UPSTREAM"), name))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close(context.Background()) })
	return conn
}

// wert liefert den ersten Wert der ersten Zeile der Anfrage.
func wert(t *testing.T, conn *pgconn.PgConn, sql string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := conn.Exec(ctx, sql).ReadAll()
	if err != nil || len(res) == 0 || len(res[0].Rows) == 0 {
		t.Fatalf("%s: %v", sql, err)
	}
	return string(res[0].Rows[0][0])
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary — play spielt eine Aufzeichnung
// mit DDL und DML als einfache Anfragen gegen eine leere Datenbank ohne
// Passwort und TLS ein (--database vor der Aufzeichnung): Die Datenbank
// enthält danach deren Wirkung, die Sessions liefen nacheinander in der
// Reihenfolge der Aufzeichnung, jede über eine eigene Verbindung; Exit-Code 0
// (Abnahmeszenario 12).
func TestE2EPlayDDLDML(t *testing.T) {
	conn := leereDatenbank(t, "play_ddl_dml")
	input := einspielAufzeichnung(t,
		[]string{"CREATE TABLE log (id serial PRIMARY KEY, n int, pid int)", "INSERT INTO log (n, pid) VALUES (1, pg_backend_pid())", "INSERT INTO log (n, pid) VALUES (2, pg_backend_pid())"},
		[]string{"INSERT INTO log (n, pid) VALUES (3, pg_backend_pid())", "UPDATE log SET n = n * 10 WHERE n = 2"},
	)
	stdout, stderr, code := starteBis(t, "", "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "play_ddl_dml")
	if code != 0 || stdout != "" || strings.Contains(stderr, "level=ERROR") {
		t.Fatalf("Exit-Code %d, stdout %q, stderr:\n%s", code, stdout, stderr)
	}
	if got := wert(t, conn, "SELECT string_agg(n::text, ',' ORDER BY id) FROM log"); got != "1,20,3" {
		t.Fatalf("Wirkung %q, erwartet 1,20,3", got)
	}
	if got := wert(t, conn, "SELECT count(DISTINCT pid) FROM log"); got != "2" {
		t.Fatalf("Verbindungen %s, erwartet 2", got)
	}
}

// Abdeckung: LH-FA-20/Negative — beantwortet der Server eine Anfrage mit einem
// Fehler, bricht play ab: PGR-E4004 mit SQLSTATE im Log, Exit-Code 4, keine
// weitere Anfrage (LH-FA-20.a).
func TestE2EPlayFehlerantwort(t *testing.T) {
	conn := leereDatenbank(t, "play_fehler")
	input := einspielAufzeichnung(t,
		[]string{"CREATE TABLE t (n int)", "SELECT * FROM fehlt", "INSERT INTO t VALUES (1)"},
		[]string{"INSERT INTO t VALUES (2)"},
	)
	_, stderr, code := starteBis(t, "", "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "play_fehler")
	if code != 4 || !strings.Contains(stderr, "code=PGR-E4004") || !strings.Contains(stderr, "42P01") || !strings.Contains(stderr, "Session 1, Interaktion 2") {
		t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
	}
	if got := wert(t, conn, "SELECT count(*) FROM t"); got != "0" {
		t.Fatalf("nach dem Abbruch eingespielt: %s Zeilen", got)
	}
}

// Abdeckung: LH-FA-20/Negative — im Aufbau sind ein nicht erreichbarer Server
// und eine fehlende Datenbank PGR-E4002, ein unbekannter Benutzer PGR-E4005,
// jeweils Exit-Code 4; bei einer Fehlerantwort nennt der Fehlertext deren
// SQLSTATE und Meldung (LH-FA-20.a *Aufbau*, *Meldungen*).
func TestE2EPlayAufbau(t *testing.T) {
	input := einspielAufzeichnung(t, []string{"SELECT 1"})
	for _, f := range []struct {
		name string
		args []string
		code string
		// meldung steht im Fehlertext: SQLSTATE und Meldung der Fehlerantwort.
		meldung string
	}{
		{"nicht erreichbar", []string{"--upstream", freieAdresse(t)}, "PGR-E4002", "nicht erreichbar"},
		{"fehlende Datenbank", []string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--database", "gibt_es_nicht"}, "PGR-E4002", `3D000 „database \"gibt_es_nicht\" does not exist“`},
		{"unbekannter Benutzer", []string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--user", "gibt_es_nicht"}, "PGR-E4005", `28000 „role \"gibt_es_nicht\" does not exist“`},
	} {
		_, stderr, code := starteBis(t, "", append([]string{"play", "--input", input}, f.args...)...)
		if code != 4 || !strings.Contains(stderr, "code="+f.code) || !strings.Contains(stderr, f.meldung) {
			t.Errorf("%s: Exit-Code %d, erwartet %s mit %s, stderr:\n%s", f.name, code, f.code, f.meldung, stderr)
		}
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-17/Happy — die Optionen von play wirken aus
// dem Abschnitt play: und einer benannten Verbindung, deren Platzhalter aus
// der Umgebung eingesetzt werden; die Datenbank der Verbindung geht der
// Aufzeichnung vor; --fail-on-unconsumed ist bei play unbekannt (PGR-E2001)
// (LH-FA-17.a *Wirkung einer URL*, LH-FA-03.b).
func TestE2EPlayKonfiguration(t *testing.T) {
	conn := leereDatenbank(t, "play_konfiguration")
	dir := t.TempDir()
	input := einspielAufzeichnung(t, []string{"CREATE TABLE aus_datei (n int)"})
	host, port, _ := strings.Cut(os.Getenv("PGR_UPSTREAM"), ":")
	inhalt := fmt.Sprintf("log_level: warn\nconnections:\n  ziel: \"postgresql://${PGR_E2E_BENUTZER}@%s:%s/${PGR_E2E_DB}\"\nplay:\n  upstream: ziel\n  input: %s\n", host, port, input)
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte(inhalt), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PGR_E2E_BENUTZER", "postgres")
	t.Setenv("PGR_E2E_DB", "play_konfiguration")
	stdout, stderr, code := starteBis(t, dir, "play")
	if code != 0 || stdout != "" || stderr != "" {
		t.Fatalf("Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	if got := wert(t, conn, "SELECT count(*) FROM pg_tables WHERE tablename = 'aus_datei'"); got != "1" {
		t.Fatalf("Tabelle aus der Aufzeichnung fehlt: %s", got)
	}
	if _, stderr, code := starteBis(t, dir, "play", "--fail-on-unconsumed"); code != 2 || !strings.Contains(stderr, "PGR-E2001") {
		t.Fatalf("--fail-on-unconsumed: Exit-Code %d, stderr %q", code, stderr)
	}
}

// aufzeichnungMit schreibt eine Aufzeichnung mit einer Session aus einer
// einfachen Anfrage sql und dem Startup user und database und liefert ihren
// Pfad.
func aufzeichnungMit(t *testing.T, user, database, sql string) string {
	t.Helper()
	inhalt := fmt.Sprintf("format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    startup:\n      user: %s\n      database: %s\n    interactions:\n      - sequence: 1\n        request:\n          type: query\n          sql: %s\n        responses:\n          - type: command_complete\n            tag: \"OK\"\n          - type: ready_for_query\n            tx_status: \"I\"\n", user, database, strconv.Quote(sql))
	pfad := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(pfad, []byte(inhalt), 0o644); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// ausfuehren führt sql auf conn aus.
func ausfuehren(t *testing.T, conn *pgconn.PgConn, sql string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := conn.Exec(ctx, sql).ReadAll(); err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
}

// spurVon liefert je Zeile des Falls fall in den Tabellen spur der
// Datenbanken von spur Benutzer und Datenbank als "wer/datenbank".
func spurVon(t *testing.T, spur map[string]*pgconn.PgConn, fall string) []string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var gefunden []string
	for db, conn := range spur {
		res, err := conn.Exec(ctx, "SELECT wer FROM spur WHERE fall = '"+fall+"'").ReadAll()
		if err != nil {
			t.Fatal(err)
		}
		for _, z := range res[0].Rows {
			gefunden = append(gefunden, string(z[0])+"/"+db)
		}
	}
	return gefunden
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-17/Boundary — gegen die reale Instanz
// gehen --user und --database, aus Kommandozeile oder Umgebung, Benutzer und
// Datenbank der benutzten Verbindung vor, und diese der Aufzeichnung; einen
// Benutzer, den die URL nicht schreibt, nimmt play aus der Aufzeichnung
// (LH-FA-17.a *Wirkung einer URL*, LH-FA-20.a *Startup-Daten*).
func TestE2EPlayVorrang(t *testing.T) {
	const (
		aufgez, url, opt       = "play_vr_aufgez", "play_vr_url", "play_vr_opt"
		aufgezdb, urldb, optdb = "play_vr_aufgezdb", "play_vr_urldb", "play_vr_optdb"
	)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	admin, err := pgconn.Connect(ctx, dsn(os.Getenv("PGR_UPSTREAM")))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = admin.Close(ctx) }()
	for _, rolle := range []string{aufgez, url, opt} {
		if _, err := admin.Exec(ctx, "DROP ROLE IF EXISTS "+rolle+"; CREATE ROLE "+rolle+" LOGIN").ReadAll(); err != nil {
			t.Fatalf("Rolle %s: %v", rolle, err)
		}
	}
	// spur hält je Datenbank, unter welchem Benutzer welcher Fall einspielte.
	spur := map[string]*pgconn.PgConn{}
	for _, db := range []string{aufgezdb, urldb, optdb} {
		conn := leereDatenbank(t, db)
		if _, err := conn.Exec(ctx, "CREATE TABLE spur (wer text, fall text); GRANT INSERT ON spur TO PUBLIC").ReadAll(); err != nil {
			t.Fatalf("spur in %s: %v", db, err)
		}
		spur[db] = conn
	}
	dir := t.TempDir()
	host, port, _ := strings.Cut(os.Getenv("PGR_UPSTREAM"), ":")
	inhalt := fmt.Sprintf("connections:\n  url: \"postgresql://%[3]s@%[1]s:%[2]s/%[4]s\"\n  ohneuser: \"postgresql://%[1]s:%[2]s/%[4]s\"\n", host, port, url, urldb)
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte(inhalt), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, f := range []struct {
		name    string
		args    []string
		env     map[string]string
		wer, wo string
	}{
		{"URL", []string{"--upstream", "url"}, nil, url, urldb},
		{"--user vor URL", []string{"--upstream", "url", "--user", opt}, nil, opt, urldb},
		{"--database vor URL", []string{"--upstream", "url", "--database", optdb}, nil, url, optdb},
		{"Umgebung vor URL", []string{"--upstream", "url"}, map[string]string{"PGWIRE_RECORDER_USER": opt, "PGWIRE_RECORDER_DATABASE": optdb}, opt, optdb},
		{"Aufzeichnung", []string{"--upstream", os.Getenv("PGR_UPSTREAM")}, nil, aufgez, aufgezdb},
		{"Optionen vor Aufzeichnung", []string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--user", opt, "--database", optdb}, nil, opt, optdb},
		{"URL ohne Benutzer", []string{"--upstream", "ohneuser"}, nil, aufgez, urldb},
	} {
		t.Run(f.name, func(t *testing.T) {
			for k, v := range f.env {
				t.Setenv(k, v)
			}
			input := aufzeichnungMit(t, aufgez, aufgezdb, "INSERT INTO spur VALUES (current_user, '"+f.name+"')")
			_, stderr, code := starteBis(t, dir, append([]string{"play", "--input", input}, f.args...)...)
			if code != 0 {
				t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
			}
			gefunden := spurVon(t, spur, f.name)
			if want := f.wer + "/" + f.wo; len(gefunden) != 1 || gefunden[0] != want {
				t.Fatalf("eingespielt als %v, erwartet %s", gefunden, want)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Negative — bis zu den Folge-Slices spielt play gegen die
// reale Instanz nichts ein, wenn die benutzte Verbindung sslmode=require
// verlangt (PGR-E2004, Exit-Code 2) und mit --upstream-tls (PGR-E2001,
// Exit-Code 2); keiner der Fälle schreibt eine Log-Zeile (LH-FA-17.a *Wirkung
// einer URL*).
func TestE2EPlayZwischenstand(t *testing.T) {
	conn := leereDatenbank(t, "play_zwischenstand")
	einfach := einspielAufzeichnung(t, []string{"CREATE TABLE vorher (n int)"})
	dir := t.TempDir()
	host, port, _ := strings.Cut(os.Getenv("PGR_UPSTREAM"), ":")
	datei := fmt.Sprintf("connections:\n  tls: \"postgresql://postgres@%s:%s/play_zwischenstand?sslmode=require\"\n", host, port)
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte(datei), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, f := range []struct {
		name    string
		args    []string
		exit    int
		meldung []string
	}{
		{"sslmode=require", []string{"--upstream", "tls", "--input", einfach}, 2, []string{"PGR-E2004", "connections.tls", "sslmode=require"}},
		{"--upstream-tls", []string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--input", einfach, "--database", "play_zwischenstand", "--upstream-tls=require"}, 2, []string{"PGR-E2001", "upstream-tls"}},
	} {
		t.Run(f.name, func(t *testing.T) {
			// Ein Fall beginnt ohne die Tabelle, die ein Fall davor eingespielt
			// haben kann.
			ausfuehren(t, conn, "DROP TABLE IF EXISTS vorher")
			_, stderr, code := starteBis(t, dir, append([]string{"play"}, f.args...)...)
			if code != f.exit || strings.Contains(stderr, "level=") {
				t.Fatalf("Exit-Code %d, erwartet %d ohne Log-Zeile, stderr:\n%s", code, f.exit, stderr)
			}
			for _, m := range f.meldung {
				if !strings.Contains(stderr, m) {
					t.Fatalf("stderr ohne %q:\n%s", m, stderr)
				}
			}
			if got := wert(t, conn, "SELECT count(*) FROM pg_tables WHERE tablename = 'vorher'"); got != "0" {
				t.Fatalf("eingespielt: Tabelle vorher %s-mal", got)
			}
		})
	}
}
