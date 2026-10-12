//go:build integration

package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// nimmAufVergleich zeichnet über record gegen die Datenbank quelle einfache
// Anfragen in einer Session auf und liefert den Pfad der Aufzeichnung: DDL, DML
// mit zwei Zeilen, eine Abfrage mit Spalten, eine Anfrage, die mit
// SQLSTATE 42P01 endet, und eine Abfrage ohne Spalten der Tabelle.
func nimmAufVergleich(t *testing.T, quelle string) string {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsnDatenbank(rec.listen, quelle))
	if err != nil {
		t.Fatal(err)
	}
	for i, sql := range []string{
		"CREATE TABLE vgl (id int, name text)",
		"INSERT INTO vgl VALUES (1, 'a'), (2, 'b')",
		"SELECT * FROM vgl",
		"SELECT * FROM fehlt",
		"SELECT 1",
	} {
		_, err := conn.Exec(ctx, sql).ReadAll()
		if (i == 3) != (err != nil) {
			t.Fatalf("%s: %v", sql, err)
		}
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
	rec.stop(t, 0)
	return output
}

// geaendert schreibt die Aufzeichnung input mit ersetze statt vorher in eine
// neue Datei und liefert deren Pfad; kommt vorher nicht vor, wird der Test rot.
func geaendert(t *testing.T, input, vorher, nachher string) string {
	t.Helper()
	text := lies(t, input)
	if !strings.Contains(text, vorher) {
		t.Fatalf("die Aufzeichnung enthält %q nicht:\n%s", vorher, text)
	}
	pfad := filepath.Join(t.TempDir(), "geaendert.yaml")
	if err := os.WriteFile(pfad, []byte(strings.Replace(text, vorher, nachher, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// zeileEnde ist die Zeile zum Ende des Einspielens in stderr.
var zeileEnde = regexp.MustCompile(`msg="play beendet" eingespielt=(\d+) verglichen=(\d+) abweichend=(\d+)`)

// zaehler liefert eingespielt, verglichen und abweichend der Zeile zum Ende,
// als "e/v/a"; ohne genau eine solche Zeile wird der Test rot.
func zaehler(t *testing.T, stderr string) string {
	t.Helper()
	treffer := zeileEnde.FindAllStringSubmatch(stderr, -1)
	if len(treffer) != 1 || strings.Count(stderr, "play beendet") != 1 {
		t.Fatalf("keine einzelne Zeile zum Ende mit Zählern:\n%s", stderr)
	}
	return treffer[0][1] + "/" + treffer[0][2] + "/" + treffer[0][3]
}

// vergleichLauf spielt input gegen eine leere Datenbank db ein, mit
// --compare-responses und den weiteren args, und liefert stderr und den
// Exit-Code.
func vergleichLauf(t *testing.T, db, input string, args ...string) (string, int) {
	t.Helper()
	leereDatenbank(t, db)
	_, stderr, code := starteBisFrist(t, 30*time.Second, append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", db, "--compare-responses"}, args...)...)
	return stderr, code
}

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Boundary — gegen die
// reale Instanz mit gleicher Antwortstruktur endet play mit --compare-responses
// mit Exit-Code 0, auch mit --allow-recorded-errors; der Fehler, den die
// Aufzeichnung genauso enthält, ist erwartet; unterscheiden sich nur Zeilenzahl
// und Zeilenwerte, ist das keine Abweichung; die Zeile zum Ende trägt
// eingespielt, verglichen und abweichend (Abnahmeszenario 16, LH-FA-24.a).
func TestE2EPlayVergleichGleich(t *testing.T) {
	leereDatenbank(t, "vgl_quelle")
	input := nimmAufVergleich(t, "vgl_quelle")
	if !strings.Contains(lies(t, input), "42P01") {
		t.Fatalf("die Aufzeichnung trägt den Fehler nicht:\n%s", lies(t, input))
	}
	stderr, code := vergleichLauf(t, "vgl_gleich", input)
	if code != 0 || strings.Contains(stderr, "level=ERROR") || zaehler(t, stderr) != "5/5/0" {
		t.Fatalf("Exit-Code %d, Zähler %s, stderr:\n%s", code, zaehler(t, stderr), stderr)
	}
	stderr, code = vergleichLauf(t, "vgl_gleich_allow", input, "--allow-recorded-errors")
	if code != 0 || zaehler(t, stderr) != "5/5/0" {
		t.Fatalf("mit --allow-recorded-errors: Exit-Code %d, stderr:\n%s", code, stderr)
	}
	weniger := geaendert(t, input, "VALUES (1, 'a'), (2, 'b')", "VALUES (1, 'a')")
	stderr, code = vergleichLauf(t, "vgl_weniger", weniger)
	if code != 0 || strings.Contains(stderr, "level=ERROR") || zaehler(t, stderr) != "5/5/0" {
		t.Fatalf("eine Zeile weniger: Exit-Code %d, stderr:\n%s", code, stderr)
	}
}

// Abdeckung: LH-FA-24/Negative — gegen eine Instanz, deren Antwort in der Struktur
// abweicht, meldet play PGR-E5004 mit Session, Sequenznummer, Position und Art
// und endet mit Exit-Code 5; weicht der SQLSTATE eines Fehlers ab, nennt die
// Meldung beide, auch mit --allow-recorded-errors; ein Serverfehler ohne
// Vorbild in der Aufzeichnung ist PGR-E5004 und nicht PGR-E4004; play bricht vor
// der nächsten Interaktion ab (Abnahmeszenario 16, LH-FA-24.a, LH-FA-20.a
// Tabelle).
func TestE2EPlayVergleichAbweichung(t *testing.T) {
	leereDatenbank(t, "vgl_quelle_a")
	input := nimmAufVergleich(t, "vgl_quelle_a")
	for _, f := range []struct {
		name, vorher, nachher string
		args                  []string
		meldung, zaehler      string
	}{
		{"Spaltenbeschreibung", "id int", "id bigint", nil, "Session 1, Interaktion 3: Antwort weicht von der Aufzeichnung ab: Position 1, Art Spaltenbeschreibung", "3/3/1"},
		{"Serverfehler ohne Vorbild", "INSERT INTO vgl", "INSERT INTO fehlt", nil, "Session 1, Interaktion 2: Antwort weicht von der Aufzeichnung ab: Position 1, Art Nachrichtenart", "2/2/1"},
		{"anderer Befehl", "INSERT INTO vgl VALUES (1, 'a'), (2, 'b')", "UPDATE vgl SET id = id", nil, "Session 1, Interaktion 2: Antwort weicht von der Aufzeichnung ab: Position 1, Art Befehl", "2/2/1"},
		{"anderer SQLSTATE", "SELECT * FROM fehlt", "SELECT 1/0", []string{"--allow-recorded-errors"}, "Session 1, Interaktion 4: Antwort weicht von der Aufzeichnung ab: Position 1, Art Fehler (SQLSTATE aufgezeichnet „42P01“, Server „22012“)", "4/4/1"},
	} {
		t.Run(f.name, func(t *testing.T) {
			datei := geaendert(t, input, f.vorher, f.nachher)
			stderr, code := vergleichLauf(t, "vgl_abw", datei, f.args...)
			if code != 5 || strings.Count(stderr, "code=PGR-E5004") != 1 || strings.Contains(stderr, "PGR-E4004") || !strings.Contains(stderr, f.meldung) || zaehler(t, stderr) != f.zaehler {
				t.Fatalf("Exit-Code %d, Zähler %s, erwartet 5 mit einer Meldung %q und Zähler %s:\n%s", code, zaehler(t, stderr), f.meldung, f.zaehler, stderr)
			}
		})
	}
}

// Abdeckung: LH-FA-24/Boundary — mit --continue-on-error läuft play nach einer
// Abweichung weiter und endet mit Exit-Code 5: Eine Abweichung ist eine
// Meldung, die weiteren Interaktionen (auch der erwartete Fehler) laufen und
// werden verglichen; ein Fehler ohne Vorbild ist weiter PGR-E5004 (Abnahmeszenario
// 16, LH-FA-20.a Schritt 6).
func TestE2EPlayVergleichFortsetzung(t *testing.T) {
	leereDatenbank(t, "vgl_quelle_f")
	input := nimmAufVergleich(t, "vgl_quelle_f")
	datei := geaendert(t, input, "INSERT INTO vgl", "INSERT INTO fehlt")
	stderr, code := vergleichLauf(t, "vgl_fortsetzung", datei, "--continue-on-error")
	if code != 5 || strings.Count(stderr, "code=PGR-E5004") != 1 || strings.Contains(stderr, "PGR-E4004") || zaehler(t, stderr) != "5/5/1" {
		t.Fatalf("Exit-Code %d, Zähler %s, stderr:\n%s", code, zaehler(t, stderr), stderr)
	}
}

// Abdeckung: LH-FA-24/Negative — eine Fehlerantwort FATAL ist mit Vergleich
// PGR-E5004 und nicht PGR-E4003, auch mit --continue-on-error ohne weitere
// Interaktion der Session; die Meldung nennt den SQLSTATE beider Seiten;
// Exit-Code 5 (Abnahmeszenario 16, LH-FA-24.a *Ende der Aufzeichnung*).
func TestE2EPlayVergleichFatal(t *testing.T) {
	leereDatenbank(t, "vgl_quelle_x")
	input := nimmAufVergleich(t, "vgl_quelle_x")
	datei := geaendert(t, input, "SELECT * FROM fehlt", "DO $$ BEGIN PERFORM pg_terminate_backend(pg_backend_pid()); PERFORM pg_sleep(5); END $$")
	for _, f := range []struct {
		db   string
		args []string
	}{{"vgl_fatal_a", nil}, {"vgl_fatal_b", []string{"--continue-on-error"}}} {
		stderr, code := vergleichLauf(t, f.db, datei, f.args...)
		meldung := "Session 1, Interaktion 4: Antwort weicht von der Aufzeichnung ab: Position 1, Art Fehler (SQLSTATE aufgezeichnet „42P01“, Server „57P01“)"
		if code != 5 || strings.Count(stderr, "code=PGR-E5004") != 1 || strings.Contains(stderr, "PGR-E4003") || !strings.Contains(stderr, meldung) || zaehler(t, stderr) != "4/4/1" {
			t.Fatalf("%q: Exit-Code %d, Zähler %s, stderr:\n%s", f.args, code, zaehler(t, stderr), stderr)
		}
	}
}

// Abdeckung: LH-FA-24/Boundary, LH-FA-17/Boundary — ohne --compare-responses
// findet kein Vergleich statt: Ein Serverfehler ohne Vorbild ist PGR-E4004 mit
// Exit-Code 4, und die Zeile zum Ende zählt verglichen und abweichend 0; die
// Option ist über die Umgebungsvariable und den Schlüssel compare_responses
// unter play: gesetzt, die Kommandozeile mit =false geht beiden vor; ein
// ungültiger Wert ist PGR-E2001 beziehungsweise PGR-E2004 (LH-FA-24.a,
// LH-FA-17.a).
func TestE2EPlayVergleichQuellen(t *testing.T) {
	leereDatenbank(t, "vgl_quelle_q")
	input := nimmAufVergleich(t, "vgl_quelle_q")
	datei := geaendert(t, input, "INSERT INTO vgl", "INSERT INTO fehlt")
	dir := t.TempDir()
	konfiguration := filepath.Join(dir, "k.yaml")
	if err := os.WriteFile(konfiguration, []byte("play:\n  compare_responses: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	basis := func(db string, extra ...string) []string {
		leereDatenbank(t, db)
		return append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", datei, "--database", db}, extra...)
	}
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "")
	_, stderr, code := starteBisFrist(t, 30*time.Second, basis("vgl_q_ohne")...)
	if code != 4 || !strings.Contains(stderr, "PGR-E4004") || strings.Contains(stderr, "PGR-E5004") || zaehler(t, stderr) != "2/0/0" {
		t.Fatalf("ohne Option: Exit-Code %d, Zähler %s, stderr:\n%s", code, zaehler(t, stderr), stderr)
	}
	_, stderr, code = starteBisFrist(t, 30*time.Second, basis("vgl_q_datei", "--config", konfiguration)...)
	if code != 5 || !strings.Contains(stderr, "PGR-E5004") || zaehler(t, stderr) != "2/2/1" {
		t.Fatalf("Schlüssel in der Datei: Exit-Code %d, Zähler %s, stderr:\n%s", code, zaehler(t, stderr), stderr)
	}
	_, stderr, code = starteBisFrist(t, 30*time.Second, basis("vgl_q_aus", "--config", konfiguration, "--compare-responses=false")...)
	if code != 4 || zaehler(t, stderr) != "2/0/0" {
		t.Fatalf("Kommandozeile false vor der Datei: Exit-Code %d, stderr:\n%s", code, stderr)
	}
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "true")
	_, stderr, code = starteBisFrist(t, 30*time.Second, basis("vgl_q_env")...)
	if code != 5 || !strings.Contains(stderr, "PGR-E5004") || zaehler(t, stderr) != "2/2/1" {
		t.Fatalf("Umgebungsvariable: Exit-Code %d, Zähler %s, stderr:\n%s", code, zaehler(t, stderr), stderr)
	}
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "1")
	if _, stderr, code := starteBisFrist(t, 30*time.Second, basis("vgl_q_ungueltig")...); code != 2 || !strings.Contains(stderr, "PGR-E2001") || !strings.Contains(stderr, "PGWIRE_RECORDER_COMPARE_RESPONSES") {
		t.Fatalf("ungültige Umgebungsvariable: Exit-Code %d, stderr:\n%s", code, stderr)
	}
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "")
	ungueltig := filepath.Join(dir, "u.yaml")
	if err := os.WriteFile(ungueltig, []byte("play:\n  compare_responses: 1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, stderr, code := starteBisFrist(t, 30*time.Second, basis("vgl_q_datei_ungueltig", "--config", ungueltig)...); code != 2 || !strings.Contains(stderr, "PGR-E2004") || !strings.Contains(stderr, "play.compare_responses") {
		t.Fatalf("ungültiger Schlüssel in der Datei: Exit-Code %d, stderr:\n%s", code, stderr)
	}
}

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Negative, LH-FA-18/Happy — mit Vergleich
// spielt play eine Aufzeichnung mit Extended-Interaktionen (Gruppen mit Sync
// und mit Flush, gemischt mit einer einfachen Anfrage) gegen eine leere
// Datenbank mit Exit-Code 0 und ohne Abweichung ein; weicht ein Parametertyp ab
// (die Spalte hat einen anderen Typ), ist es PGR-E5004 der Art Spaltenbeschreibung
// mit Exit-Code 5 (Abnahmeszenario 16, LH-FA-24.a).
func TestE2EPlayVergleichExtended(t *testing.T) {
	leereDatenbank(t, "vgl_ext_quelle")
	input := nimmAufEinspielAblauf(t, "vgl_ext_quelle")
	stderr, code := vergleichLauf(t, "vgl_ext_gleich", input)
	z := strings.Split(zaehler(t, stderr), "/")
	if code != 0 || strings.Contains(stderr, "level=ERROR") || z[0] != z[1] || z[2] != "0" || z[0] == "0" {
		t.Fatalf("Exit-Code %d, Zähler %v, stderr:\n%s", code, z, stderr)
	}
	datei := geaendert(t, input, "zahl int)", "zahl bigint)")
	stderr, code = vergleichLauf(t, "vgl_ext_abw", datei)
	if code != 5 || strings.Count(stderr, "code=PGR-E5004") != 1 || !strings.Contains(stderr, "Art Spaltenbeschreibung") || strings.Contains(stderr, "PGR-E4004") {
		t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
	}
}

// Abdeckung: LH-FA-24/Boundary, LH-FA-24/Negative — endet eine Extended-Interaktion
// gegen die reale Instanz mit einer Fehlerantwort, die die Aufzeichnung genauso
// trägt (ein Batch, dessen zweite Anweisung durch null teilt), sendet play mit
// Vergleich die Gruppe bis zum Sync und endet mit Exit-Code 0 ohne Meldung; die
// Anweisung nach dem Fehler läuft nicht, die danach folgende Anfrage schon. Trägt
// die Aufzeichnung einen anderen SQLSTATE, ist es PGR-E5004 der Art Fehler mit
// beiden SQLSTATE und Exit-Code 5 (Abnahmeszenario 16, LH-FA-24.a, LH-FA-20.a
// Schritt 6).
func TestE2EPlayVergleichExtendedFehler(t *testing.T) {
	leereDatenbank(t, "vgl_extf_quelle")
	input := nimmAufEinspielFehler(t, "vgl_extf_quelle")
	conn := leereDatenbank(t, "vgl_extf_gleich")
	_, stderr, code := starteBisFrist(t, 30*time.Second, "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "vgl_extf_gleich", "--compare-responses")
	z := strings.Split(zaehler(t, stderr), "/")
	if code != 0 || strings.Contains(stderr, "level=ERROR") || z[0] != z[1] || z[2] != "0" || z[0] == "0" {
		t.Fatalf("Exit-Code %d, Zähler %v, stderr:\n%s", code, z, stderr)
	}
	if got := wert(t, conn, "SELECT coalesce(string_agg(n::text, ',' ORDER BY n), '') FROM ergebnis"); got != "4" {
		t.Fatalf("Zeilen %q, erwartet 4", got)
	}
	datei := geaendert(t, input, "22012", "22003")
	stderr, code = vergleichLauf(t, "vgl_extf_abw", datei)
	if code != 5 || strings.Count(stderr, "code=PGR-E5004") != 1 || !strings.Contains(stderr, "Art Fehler (SQLSTATE aufgezeichnet „22003“, Server „22012“)") {
		t.Fatalf("anderer SQLSTATE: Exit-Code %d, stderr:\n%s", code, stderr)
	}
}

// nimmAufEinfach zeichnet über record gegen die Datenbank quelle die
// einfachen Anfragen sqls in einer Session auf und liefert den Pfad der
// Aufzeichnung.
func nimmAufEinfach(t *testing.T, quelle string, sqls ...string) string {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsnDatenbank(rec.listen, quelle))
	if err != nil {
		t.Fatal(err)
	}
	for _, sql := range sqls {
		if _, err := conn.Exec(ctx, sql).ReadAll(); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
	rec.stop(t, 0)
	return output
}

// Abdeckung: LH-FA-24/Boundary — gegen die reale Instanz sind Hinweise
// (NoticeResponse) und Parameterstatus (ParameterStatus) kein Teil des
// Vergleichs, in beide Richtungen: Trägt nur die Aufzeichnung einen Hinweis und
// einen Parameterstatus, trägt sie nur der Server, oder beide, endet play mit
// Exit-Code 0 und ohne Abweichung (LH-FA-24.a, Normalisierung).
func TestE2EPlayVergleichHinweise(t *testing.T) {
	leereDatenbank(t, "vgl_hinw_quelle")
	mit := nimmAufEinfach(t, "vgl_hinw_quelle", "SET application_name = 'a'", "DO $$ BEGIN RAISE NOTICE 'x'; END $$")
	text := lies(t, mit)
	if !strings.Contains(text, "type: parameter_status") || !strings.Contains(text, "type: notice_response") {
		t.Fatalf("die Aufzeichnung trägt Parameterstatus und Hinweis nicht:\n%s", text)
	}
	ohne := nimmAufEinfach(t, "vgl_hinw_quelle", "SET work_mem = '5MB'", "DO $$ BEGIN NULL; END $$")
	if text := lies(t, ohne); strings.Contains(text, "type: parameter_status") || strings.Contains(text, "type: notice_response") {
		t.Fatalf("die Aufzeichnung trägt Parameterstatus oder Hinweis:\n%s", text)
	}
	nurAufgezeichnet := geaendert(t, geaendert(t, mit, "SET application_name = 'a'", "SET work_mem = '5MB'"), "RAISE NOTICE 'x'", "NULL")
	nurServer := geaendert(t, geaendert(t, ohne, "SET work_mem = '5MB'", "SET application_name = 'b'"), "BEGIN NULL;", "BEGIN RAISE NOTICE 'y';")
	for i, f := range []struct{ name, input string }{{"nur aufgezeichnet", nurAufgezeichnet}, {"nur Server", nurServer}, {"beide", mit}} {
		stderr, code := vergleichLauf(t, "vgl_hinw_"+string(rune('a'+i)), f.input)
		if code != 0 || strings.Contains(stderr, "level=ERROR") || zaehler(t, stderr) != "2/2/0" {
			t.Fatalf("%s: Exit-Code %d, stderr:\n%s", f.name, code, stderr)
		}
	}
}
