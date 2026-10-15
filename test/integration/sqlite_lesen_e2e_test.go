//go:build integration

package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// formatArgs sind die Argumente von record für die Aufzeichnung name: bei der Endung
// .sqlite --format sqlite, sonst keine (das Format yaml ist der Default).
func formatArgs(name string) []string {
	if strings.HasSuffix(name, ".sqlite") {
		return []string{"--format", "sqlite"}
	}
	return nil
}

// aufnehmenAls zeichnet die Anfragen über record in einer Datei name auf, bei der Endung
// .sqlite im Format sqlite, und liefert den Pfad der Aufzeichnung und die Sicht des Clients
// beim Aufzeichnen.
func aufnehmenAls(t *testing.T, name string, queries ...string) (string, string) {
	t.Helper()
	output := filepath.Join(t.TempDir(), name)
	rec := startProzess(t, "record", append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output}, formatArgs(name)...)...)
	sicht := beobachten(t, rec.listen, queries...)
	rec.stop(t, 0)
	return output, sicht
}

// unveraendert hält den Stand einer Datei fest; pruefe prüft, dass sie nach einem Lauf
// dieselben Bytes und dieselbe Änderungszeit hat wie bei merke und im Verzeichnis neben ihr
// nichts liegt.
type unveraendert struct {
	path   string
	bytes  []byte
	mtimeN int64
}

func merke(t *testing.T, path string) unveraendert {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return unveraendert{path, []byte(lies(t, path)), info.ModTime().UnixNano()}
}

func (u unveraendert) pruefe(t *testing.T) {
	t.Helper()
	info, err := os.Stat(u.path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(lies(t, u.path)), u.bytes) || info.ModTime().UnixNano() != u.mtimeN {
		t.Errorf("%s nach dem Lauf verändert", u.path)
	}
	verzeichnisEnthaelt(t, filepath.Dir(u.path), filepath.Base(u.path))
}

// Abdeckung: LH-FA-22/Happy, LH-FA-07/Happy — Abnahmeszenario 14 für einfache Anfragen: dieselben Anfragen,
// einmal mit record --format sqlite und einmal als yaml aufgezeichnet (Zeichen wie base64: und \
// in den Werten), liefern über replay, das das Format selbst erkennt, beide dieselbe Sicht wie
// beim Aufzeichnen, mehrere Läufe hintereinander; die SQLite-Datei bleibt dabei unverändert,
// und neben ihr liegt nichts.
func TestE2EReplaySQLite(t *testing.T) {
	queries := []string{"SELECT 1 AS eins;", "SELECT 'a' AS text, NULL AS leer;", `SELECT 'base64:abc' AS b, E'\\x41' AS fluchtzeichen, '' AS leerer_text;`, "SELECT 'café 日本' AS unicode;"}
	sqlite, aufgezeichnet := aufnehmenAls(t, "rec.sqlite", queries...)
	yaml, ausYAML := aufnehmenAls(t, "rec.yaml", queries...)
	if !strings.HasPrefix(lies(t, sqlite), sqliteKopf) || strings.HasPrefix(lies(t, yaml), sqliteKopf) {
		t.Fatal("Formate der Aufzeichnungen vertauscht")
	}
	if !strings.Contains(aufgezeichnet, `zeile ["1"]`) || aufgezeichnet != ausYAML {
		t.Fatalf("Sicht beim Aufzeichnen:\n%s\nals yaml:\n%s", aufgezeichnet, ausYAML)
	}
	vorher := merke(t, sqlite)
	for lauf := 1; lauf <= 3; lauf++ {
		for _, input := range []string{sqlite, yaml} {
			rep := startProzess(t, "replay", "--input", input)
			got := beobachten(t, rep.listen, queries...)
			rep.stop(t, 0)
			if got != aufgezeichnet {
				t.Fatalf("Lauf %d, %s weicht von der Sicht beim Aufzeichnen ab:\n%s\n--- aufgezeichnet:\n%s", lauf, input, got, aufgezeichnet)
			}
		}
	}
	vorher.pruefe(t)
}

// Abdeckung: LH-FA-22/Happy, LH-FA-18/Happy — Abnahmeszenario 14 für Extended Query: pgx im Standardmodus
// erhält über replay aus einer mit --format sqlite geschriebenen Aufzeichnung dieselbe Sicht wie
// beim Aufzeichnen, ebenso aus derselben Aufzeichnung als yaml; die SQLite-Datei bleibt
// unverändert.
func TestE2EReplayExtendedSQLite(t *testing.T) {
	sqlite, yaml := filepath.Join(t.TempDir(), "rec.sqlite"), filepath.Join(t.TempDir(), "rec.yaml")
	aufgezeichnet := aufnehmenExtended(t, sqlite)
	if ausYAML := aufnehmenExtended(t, yaml); ausYAML != aufgezeichnet {
		t.Fatalf("Sicht als yaml weicht ab:\n%s\n--- sqlite:\n%s", ausYAML, aufgezeichnet)
	}
	vorher := merke(t, sqlite)
	for _, input := range []string{sqlite, yaml, sqlite} {
		rep := startProzess(t, "replay", "--input", input)
		got := extendedAblauf(t, rep.listen)
		rep.stop(t, 0)
		if got != aufgezeichnet {
			t.Fatalf("%s weicht von der Sicht beim Aufzeichnen ab:\n%s\n--- aufgezeichnet:\n%s\n--- stderr:\n%s", input, got, aufgezeichnet, rep.stderr.String())
		}
	}
	vorher.pruefe(t)
}

// Abdeckung: LH-FA-22/Happy, LH-FA-20/Happy — Abnahmeszenario 14 für das Einspielen: play spielt eine mit
// --format sqlite geschriebene Aufzeichnung (Extended-Interaktionen mit Parametern einschließlich
// NULL und leer, Batch, Pipeline, einfache Anfrage) ein: Die Datenbank
// enthält danach dieselben Zeilen wie die, gegen die aufgezeichnet wurde; die SQLite-Datei bleibt
// dabei unverändert.
func TestE2EPlaySQLite(t *testing.T) {
	quelle := leereDatenbank(t, "play_sql_quelle")
	input := nimmAufEinspielAblaufAls(t, "play_sql_quelle", "rec.sqlite")
	if !strings.HasPrefix(lies(t, input), sqliteKopf) {
		t.Fatalf("keine SQLite-Datei: %.20q", lies(t, input))
	}
	erwartet := personen(t, quelle)
	if !strings.Contains(erwartet, "2:Bob::25") || !strings.Contains(erwartet, "5:Eva:<NULL>:50") || !strings.Contains(erwartet, "3:Cäsar 日本:mit 'Anführungszeichen':30") {
		t.Fatalf("Aufzeichnung gegen die Quelle: %q", erwartet)
	}
	vorher := merke(t, input)
	ziel := leereDatenbank(t, "play_sql_ziel")
	stdout, stderr, code := starteBis(t, "", "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "play_sql_ziel")
	if code != 0 || stdout != "" || strings.Contains(stderr, "level=ERROR") {
		t.Fatalf("Exit-Code %d, stdout %q, stderr:\n%s", code, stdout, stderr)
	}
	if got := personen(t, ziel); got != erwartet {
		t.Fatalf("Wirkung\n%s\nerwartet\n%s", got, erwartet)
	}
	vorher.pruefe(t)
}

// Abdeckung: LH-FA-22/Negative, LH-QA-06/Messung — replay und play lehnen eine Datei, die in ihrem Format
// keine gültige Aufzeichnung ist, in beiden Formaten gleich ab: eine unbekannte Version mit
// PGR-E3002, eine andere Formatkennung, eine abgeschnittene Datei und eine Datei mit dem Kopf
// einer SQLite-Datei und anderem Inhalt mit PGR-E3003, eine fehlende Datei mit PGR-E3001; der
// Exit-Code ist 3, und es läuft nichts an (play verbindet sich nicht).
func TestE2ELesenAbgelehnt(t *testing.T) {
	sqlite, _ := aufnehmenAls(t, "rec.sqlite", "SELECT 1;")
	yaml, _ := aufnehmenAls(t, "rec.yaml", "SELECT 1;")
	roh := []byte(lies(t, sqlite))
	text := lies(t, yaml)
	ersetzt := func(b []byte, alt, neu string) []byte {
		if !bytes.Contains(b, []byte(alt)) {
			t.Fatalf("%q nicht in der Datei", alt)
		}
		return bytes.Replace(b, []byte(alt), []byte(neu), 1)
	}
	faelle := []struct {
		name   string
		inhalt []byte
		code   string
	}{
		{"sqlite, Version 2", ersetzt(roh, "version1", "version2"), "PGR-E3002"},
		{"yaml, Version 3", ersetzt([]byte(text), "version: 1", "version: 3"), "PGR-E3002"},
		{"sqlite, andere Formatkennung", ersetzt(roh, "pgwire-recorder", "pgwire-recorden"), "PGR-E3003"},
		{"yaml, andere Formatkennung", ersetzt([]byte(text), "pgwire-recorder", "pgwire-recorden"), "PGR-E3003"},
		{"sqlite, abgeschnitten an der Seitengrenze", roh[:len(roh)/4096/2*4096], "PGR-E3003"},
		{"sqlite, abgeschnitten in einer Seite", roh[:len(roh)/4096/2*4096+100], "PGR-E3003"},
		{"yaml, abgeschnitten", []byte(text[:len(text)/2]), "PGR-E3003"},
		{"sqlite, nur der Kopf", []byte(sqliteKopf), "PGR-E3003"},
		{"sqlite, Kopf und Müll", []byte(sqliteKopf + strings.Repeat("müll", 100)), "PGR-E3003"},
	}
	if len(roh) < 3*4096 {
		t.Fatalf("SQLite-Datei nur %d Bytes", len(roh))
	}
	for _, f := range faelle {
		input := filepath.Join(t.TempDir(), "rec.in")
		if err := os.WriteFile(input, f.inhalt, 0o644); err != nil {
			t.Fatal(err)
		}
		for _, args := range [][]string{
			{"replay", "--listen", freieAdresse(t), "--input", input},
			{"play", "--upstream", freieAdresse(t), "--input", input},
		} {
			stdout, stderr, code := starteBis(t, "", args...)
			if code != 3 || !strings.Contains(stderr, "["+f.code+"]") || stdout != "" {
				t.Errorf("%s, %s: Exit-Code %d, stdout %q, stderr %q", f.name, args[0], code, stdout, stderr)
			}
		}
	}
	for _, args := range [][]string{
		{"replay", "--listen", freieAdresse(t), "--input", filepath.Join(t.TempDir(), "fehlt")},
		{"play", "--upstream", freieAdresse(t), "--input", filepath.Join(t.TempDir(), "fehlt")},
	} {
		if _, stderr, code := starteBis(t, "", args...); code != 3 || !strings.Contains(stderr, "[PGR-E3001]") {
			t.Errorf("fehlende Datei, %s: Exit-Code %d, stderr %q", args[0], code, stderr)
		}
	}
}
