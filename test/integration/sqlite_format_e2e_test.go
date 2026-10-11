//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// sqliteKopf sind die ersten 16 Bytes einer SQLite-Datei: der Text SQLite format 3
// und ein NUL-Byte. Der Test liest die Datei nur als Bytes und nimmt keinen
// SQLite-Treiber und kein sqlite3 zu Hilfe.
const sqliteKopf = "SQLite format 3\x00"

// wartetAufInhalt wartet höchstens 15 s, bis die Datei path den Text teil enthält, und
// liefert ihren Inhalt; sonst endet der Test.
func wartetAufInhalt(t *testing.T, path, teil string) []byte {
	t.Helper()
	for deadline := time.Now().Add(15 * time.Second); ; time.Sleep(50 * time.Millisecond) {
		data, err := os.ReadFile(path)
		if err == nil && bytes.Contains(data, []byte(teil)) {
			return data
		}
		if time.Now().After(deadline) {
			t.Fatalf("%q steht nicht binnen 15 s in %s (%v)", teil, path, err)
		}
	}
}

// ausfuehrenJeVerbindung führt je Anfrage eine eigene Verbindung über den Recorder aus.
func ausfuehrenJeVerbindung(t *testing.T, listen string, anfragen ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for _, q := range anfragen {
		conn, err := pgconn.Connect(ctx, dsn(listen))
		if err != nil {
			t.Fatalf("Verbindung über den Recorder: %v", err)
		}
		if _, err := conn.Exec(ctx, q).ReadAll(); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		if err := conn.Close(ctx); err != nil {
			t.Fatal(err)
		}
	}
}

// Abdeckung: LH-FA-22/Happy, LH-FA-22/Boundary — record --format sqlite schreibt eine
// SQLite-Datei (die ersten 16 Bytes sind SQLite format 3 und ein NUL-Byte) und ergänzt
// sie nach jeder beendeten Session, nicht erst beim Beenden: Die Anfrage der ersten
// Session steht schon vor dem Beenden in der Datei, die der zweiten danach; im Verzeichnis
// liegt neben der Datei weder ein Journal noch eine temporäre Datei, und der Lauf endet mit
// Exit-Code 0.
func TestE2ERecordFormatSQLite(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "rec.sqlite")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--format", "sqlite")
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("Zieldatei vor der ersten Session: %v", err)
	}

	ausfuehrenJeVerbindung(t, rec.listen, "SELECT 'erste_session';")
	data := wartetAufInhalt(t, output, "erste_session")
	if !bytes.HasPrefix(data, []byte(sqliteKopf)) {
		t.Fatalf("Datei beginnt mit %q", data[:16])
	}
	if bytes.Contains(data, []byte("zweite_session")) {
		t.Fatal("zweite Session vor ihrer Ausführung in der Datei")
	}
	ausfuehrenJeVerbindung(t, rec.listen, "SELECT 'zweite_session';")
	wartetAufInhalt(t, output, "zweite_session")

	rec.stop(t, 0)
	text := lies(t, output)
	if !strings.HasPrefix(text, sqliteKopf) || !strings.Contains(text, "erste_session") || !strings.Contains(text, "pgwire-recorder") {
		t.Fatalf("Datei nach dem Beenden: %.40q", text)
	}
	verzeichnisEnthaelt(t, dir, "rec.sqlite")
}

// Abdeckung: LH-FA-22/Boundary — das Format kommt gleichwertig aus der Umgebung
// (PGWIRE_RECORDER_FORMAT) und aus dem Abschnitt record: der Konfigurationsdatei
// (format), die Kommandozeile geht beiden vor: yaml auf der Kommandozeile schreibt YAML trotz
// sqlite in Umgebung und Datei; ohne Session beim Beenden steht eine gültige Aufzeichnung
// ohne Sessions in der Datei.
func TestE2ERecordFormatQuellen(t *testing.T) {
	for _, tc := range []struct {
		quelle string
		args   []string
		sqlite bool
	}{
		{"Umgebung", nil, true},
		{"Datei", nil, true},
		{"Umgebung mit Kommandozeile", []string{"--format", "yaml"}, false},
		{"Datei mit Kommandozeile", []string{"--format=yaml"}, false},
	} {
		dir := t.TempDir()
		output := filepath.Join(dir, "rec.out")
		if strings.HasPrefix(tc.quelle, "Umgebung") {
			t.Setenv("PGWIRE_RECORDER_FORMAT", "sqlite")
		} else if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("record:\n  format: sqlite\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		rec := startProzessIn(t, dir, "record", append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output}, tc.args...)...)
		rec.stop(t, 0)
		data := lies(t, output)
		if got := strings.HasPrefix(data, sqliteKopf); got != tc.sqlite {
			t.Errorf("%s: SQLite-Datei %v, erwartet %v: %.40q", tc.quelle, got, tc.sqlite, data)
		}
		if !tc.sqlite && !strings.Contains(data, "format: pgwire-recorder") {
			t.Errorf("%s: keine YAML-Aufzeichnung: %.40q", tc.quelle, data)
		}
		t.Setenv("PGWIRE_RECORDER_FORMAT", "")
	}
}

// Abdeckung: LH-FA-22/Negative — ein Rollback-Journal <Zielpfad>-journal eines früheren
// Laufs beendet record beim Start mit Exit-Code 3 (PGR-E3001), auch mit --force und auch
// ohne Zieldatei; Zieldatei und Journal bleiben unverändert; mit vorhandener Zieldatei
// ohne --force ist es Exit-Code 2 (PGR-E2002); YAML ist vom Journal unberührt.
func TestE2ERecordFormatSQLiteJournal(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "rec.sqlite")
	if err := os.WriteFile(output+"-journal", []byte("Rest"), 0o644); err != nil {
		t.Fatal(err)
	}
	args := func(extra ...string) []string {
		return append([]string{"record", "--listen", freieAdresse(t), "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--format", "sqlite"}, extra...)
	}
	for _, extra := range [][]string{nil, {"--force"}} {
		stdout, stderr, code := starteBis(t, dir, args(extra...)...)
		if code != 3 || !strings.HasPrefix(stderr, "Recording [PGR-E3001]: ") || !strings.Contains(stderr, "rec.sqlite-journal") || stdout != "" {
			t.Fatalf("%v: Exit-Code %d, stdout %q, stderr %q", extra, code, stdout, stderr)
		}
	}
	verzeichnisEnthaelt(t, dir, "rec.sqlite-journal")

	if err := os.WriteFile(output, []byte("alt"), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := starteBis(t, dir, args()...)
	if code != 2 || !strings.HasPrefix(stderr, "Konfiguration [PGR-E2002]: ") || stdout != "" {
		t.Fatalf("Zieldatei und Journal ohne --force: Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	_, stderr, code = starteBis(t, dir, args("--force")...)
	if code != 3 || !strings.HasPrefix(stderr, "Recording [PGR-E3001]: ") {
		t.Fatalf("Zieldatei und Journal mit --force: Exit-Code %d, stderr %q", code, stderr)
	}
	if got := lies(t, output); got != "alt" || lies(t, output+"-journal") != "Rest" {
		t.Fatalf("Dateien verändert: %q, %q", got, lies(t, output+"-journal"))
	}

	yamlZiel := filepath.Join(dir, "rec.yaml")
	if err := os.WriteFile(yamlZiel+"-journal", []byte("Rest"), 0o644); err != nil {
		t.Fatal(err)
	}
	rec := startProzessIn(t, dir, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", yamlZiel)
	rec.stop(t, 0)
}

// Abdeckung: LH-FA-22/Negative — jeder Wert von --format außerhalb von yaml und sqlite
// (der leere, SQLite, sqlite3) endet record mit Exit-Code 2 (PGR-E2001) und legt die
// Zieldatei nicht an; ungültig aus der Umgebung ist es PGR-E2001 mit dem Namen der Variable,
// aus der Datei PGR-E2004 mit dem Schlüssel und ohne den Wert; replay und play kennen die
// Option nicht (PGR-E2001, Exit-Code 2).
func TestE2EFormatUngueltig(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "rec.out")
	record := []string{"record", "--listen", freieAdresse(t), "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output}
	for _, wert := range []string{"", "SQLite", "sqlite3"} {
		stdout, stderr, code := starteBis(t, dir, append(append([]string{}, record...), "--format="+wert)...)
		if code != 2 || !strings.HasPrefix(stderr, "Konfiguration [PGR-E2001]: ") || !strings.Contains(stderr, "erlaubt sind yaml und sqlite") || stdout != "" {
			t.Errorf("--format=%q: Exit-Code %d, stdout %q, stderr %q", wert, code, stdout, stderr)
		}
	}
	t.Setenv("PGWIRE_RECORDER_FORMAT", "GEHEIM")
	_, stderr, code := starteBis(t, dir, record...)
	if code != 2 || !strings.HasPrefix(stderr, "Konfiguration [PGR-E2001]: ") || !strings.Contains(stderr, "PGWIRE_RECORDER_FORMAT") || !strings.Contains(stderr, "erlaubt sind yaml und sqlite") {
		t.Errorf("Umgebung: Exit-Code %d, stderr %q", code, stderr)
	}
	t.Setenv("PGWIRE_RECORDER_FORMAT", "")
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("record:\n  format: GEHEIM\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, stderr, code = starteBis(t, dir, record...)
	if code != 2 || !strings.HasPrefix(stderr, "Konfiguration [PGR-E2004]: ") || !strings.Contains(stderr, "record.format") || strings.Contains(stderr, "GEHEIM") {
		t.Errorf("Datei: Exit-Code %d, stderr %q", code, stderr)
	}
	if err := os.Remove(filepath.Join(dir, ".pgwire-recorder.yaml")); err != nil {
		t.Fatal(err)
	}
	verzeichnisEnthaelt(t, dir)

	for _, args := range [][]string{
		{"replay", "--listen", "127.0.0.1:0", "--input", output, "--format", "sqlite"},
		{"replay", "--listen", "127.0.0.1:0", "--input", output, "--format=yaml"},
		{"play", "--upstream", "127.0.0.1:1", "--input", output, "--format", "sqlite"},
	} {
		stdout, stderr, code := starteBis(t, dir, args...)
		if code != 2 || !strings.HasPrefix(stderr, "Konfiguration [PGR-E2001]: ") || !strings.Contains(stderr, "format") || stdout != "" {
			t.Errorf("%v: Exit-Code %d, stdout %q, stderr %q", args, code, stdout, stderr)
		}
	}
}
