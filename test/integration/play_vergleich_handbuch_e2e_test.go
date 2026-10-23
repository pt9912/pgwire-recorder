//go:build integration

package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// handbuchMeldung ist der Anfang der Meldung mit --diff-all in einer Interaktion:
// Session 1, die Nummer der Interaktion und die Gesamtzahl der Abweichungen.
func handbuchMeldung(interaktion, anzahl int) string {
	return fmt.Sprintf("Session 1, Interaktion %d: Antworten weichen von der Aufzeichnung ab, Abweichungen: %d; ", interaktion, anzahl)
}

// handbuchLauf spielt input mit --compare-responses --diff-all und den weiteren
// args ein und verlangt Exit-Code 5 mit genau einer Zeile PGR-E5004; es liefert
// stderr.
func handbuchLauf(t *testing.T, db, input string, args ...string) string {
	t.Helper()
	stderr, code := vergleichLauf(t, db, input, append([]string{"--diff-all"}, args...)...)
	if code != 5 || strings.Count(stderr, "code=PGR-E5004") != 1 {
		t.Fatalf("Exit-Code %d, erwartet 5 mit einer Zeile PGR-E5004:\n%s", code, stderr)
	}
	return stderr
}

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Negative — gegen die reale Instanz nennt
// die Meldung mit --diff-all je Art die Einzelheiten beider Seiten, wie das
// Benutzerhandbuch sie beschreibt: bei Spaltenbeschreibung die Typ-OID je Spalte
// und die Zahl der Spalten, bei Befehl beide Befehle, bei Fehler beide SQLSTATE, bei
// Transaktionsstatus beide Zustände, bei Nachrichtenart nichts; auch bei einer
// einzigen Abweichung hat die Meldung die Form mit der Gesamtzahl;
// nie ein Zeilenwert (Abnahmeszenario 16, LH-FA-24.a *Alle Abweichungen*).
func TestE2EPlayHandbuchDiffAllEinzelheiten(t *testing.T) {
	leereDatenbank(t, "hb_quelle")
	input := nimmAufVergleich(t, "hb_quelle")
	typen := geaendert(t, geaendert(t, input, "id int", "id bigint"), "name text", "name varchar")
	leereDatenbank(t, "hb_quelle_t")
	transaktion := nimmAufEinfach(t, "hb_quelle_t", "COMMIT")
	for _, f := range []struct {
		name, datei, erwartet, zaehler string
	}{
		{"Typ-OID", typen, handbuchMeldung(3, 1) + "Position 1, Art Spaltenbeschreibung (Spalte 1: Typ-OID aufgezeichnet 23, Server 20; Spalte 2: Typ-OID aufgezeichnet 25, Server 1043)", "3/3/1"},
		{"Zahl der Spalten", geaendert(t, input, "SELECT * FROM vgl", "SELECT id FROM vgl"), handbuchMeldung(3, 1) + "Position 1, Art Spaltenbeschreibung (Spalten aufgezeichnet 2, Server 1)", "3/3/1"},
		{"Befehl", geaendert(t, input, "INSERT INTO vgl VALUES (1, 'a'), (2, 'b')", "UPDATE vgl SET id = id"), handbuchMeldung(2, 1) + "Position 1, Art Befehl (Befehl aufgezeichnet „INSERT“, Server „UPDATE“)", "2/2/1"},
		{"Fehler", geaendert(t, input, "SELECT * FROM fehlt", "SELECT 1/0"), handbuchMeldung(4, 1) + "Position 1, Art Fehler (SQLSTATE aufgezeichnet „42P01“, Server „22012“)", "4/4/1"},
		{"Transaktionsstatus", geaendert(t, transaktion, "COMMIT", "BEGIN"), handbuchMeldung(1, 2) + "Position 1, Art Befehl (Befehl aufgezeichnet „COMMIT“, Server „BEGIN“); Position 2, Art Transaktionsstatus (aufgezeichnet „I“, Server „T“)", "1/1/1"},
		{"Nachrichtenart", geaendert(t, input, "INSERT INTO vgl", "INSERT INTO fehlt"), handbuchMeldung(2, 1) + "Position 1, Art Nachrichtenart", "2/2/1"},
	} {
		t.Run(f.name, func(t *testing.T) {
			stderr := handbuchLauf(t, "hb_einzelheiten", f.datei)
			if !strings.Contains(stderr, "Replay [PGR-E5004]: "+f.erwartet+"\"") || strings.Contains(stderr, "(1, 'a')") || strings.Contains(stderr, "division by zero") || strings.Contains(stderr, "does not exist") || zaehler(t, stderr) != f.zaehler {
				t.Fatalf("erwartet %q in:\n%s", f.erwartet, stderr)
			}
		})
	}
}

// spaltenAbfrage ist eine Abfrage mit n Spalten, je Spalte der Name praefix und
// ihre Nummer; mit typ bigint ist jede Spalte ein bigint, sonst ein int.
func spaltenAbfrage(n int, praefix, typ string) string {
	var s []string
	for i := 1; i <= n; i++ {
		s = append(s, fmt.Sprintf("1%s AS %s%d", typ, praefix, i))
	}
	return "SELECT " + strings.Join(s, ", ")
}

// Abdeckung: LH-FA-24/Boundary — gegen die reale Instanz nennt die Meldung mit
// --diff-all je Position höchstens 20 abweichende Spalten, dahinter „und N weitere“,
// und die Grenze zählt Spalten, nicht Felder: 20 Spalten mit abweichendem Namen und
// Typ-OID (40 Felder) stehen alle in der Meldung, die 21. ist „und 1 weitere“; ebenso
// höchstens 20 Positionen je Interaktion, 20 Abweichungen stehen alle in der Meldung,
// bei 21 bleibt eine als „und 1 weitere“ (Abnahmeszenario 16, LH-FA-24.a *Alle
// Abweichungen*, *Obergrenze*).
func TestE2EPlayHandbuchDiffAllObergrenzen(t *testing.T) {
	for _, n := range []int{20, 21, 23} {
		t.Run(fmt.Sprintf("%d Spalten", n), func(t *testing.T) {
			quelle := spaltenAbfrage(n, "c", "")
			leereDatenbank(t, "hb_spalten_q")
			input := nimmAufEinfach(t, "hb_spalten_q", quelle)
			stderr := handbuchLauf(t, "hb_spalten", geaendert(t, input, quelle, spaltenAbfrage(n, "d", "::bigint")))
			var einzelne []string
			for i := 1; i <= 20; i++ {
				einzelne = append(einzelne, fmt.Sprintf("Spalte %d: Name aufgezeichnet „c%d“, Server „d%d“, Typ-OID aufgezeichnet 23, Server 20", i, i, i))
			}
			if n > 20 {
				einzelne = append(einzelne, fmt.Sprintf("und %d weitere", n-20))
			}
			erwartet := handbuchMeldung(1, 1) + "Position 1, Art Spaltenbeschreibung (" + strings.Join(einzelne, "; ") + ")"
			if !strings.Contains(stderr, erwartet) || strings.Contains(stderr, "Spalte 21:") {
				t.Fatalf("erwartet %q in:\n%s", erwartet, stderr)
			}
		})
	}
	for _, n := range []int{19, 20, 22} {
		t.Run(fmt.Sprintf("%d Abweichungen", n+1), func(t *testing.T) {
			leereDatenbank(t, "hb_pos_q")
			input := nimmAufEinfach(t, "hb_pos_q", "SELECT 1")
			stderr := handbuchLauf(t, "hb_pos", geaendert(t, input, "SELECT 1", "SELECT 1"+strings.Repeat("; DO $$ BEGIN NULL; END $$", n)))
			zeile := stderr[strings.Index(stderr, "code=PGR-E5004"):]
			zeile = zeile[:strings.Index(zeile, "\n")]
			angezeigt := strings.Count(zeile, ", Art Nachrichtenart")
			weitere := ""
			if n+1 > 20 {
				weitere = fmt.Sprintf("; und %d weitere", n+1-20)
			}
			if !strings.Contains(zeile, handbuchMeldung(1, n+1)) || angezeigt != 20 || strings.Contains(zeile, "weitere") != (weitere != "") || !strings.HasSuffix(zeile, weitere+"\"") {
				t.Fatalf("Positionen %d, erwartet 20 und %q am Ende:\n%s", angezeigt, weitere, zeile)
			}
		})
	}
}

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Negative, LH-FA-17/Boundary — gegen die reale
// Instanz gilt --diff-all aus der Option, der Umgebungsvariable und dem Schlüssel
// diff_all unter play:, wie das Benutzerhandbuch sie nennt; ohne --compare-responses
// ist es PGR-E2001 mit Exit-Code 2, aus jeder der drei Quellen, auch bei
// --compare-responses=false und auch, wenn die Datei den Schlüssel trägt (nicht
// PGR-E2004); ein ungültiger Wert des Schlüssels unter play: ist PGR-E2004; bei record und replay ist die Option und die Umgebungsvariable
// PGR-E2001 und der Schlüssel PGR-E2004 (Abnahmeszenario 16, LH-FA-24.a *Alle
// Abweichungen*, LH-FA-17.a).
func TestE2EPlayHandbuchDiffAllQuellen(t *testing.T) {
	leereDatenbank(t, "hb_quelle_q")
	input := geaendert(t, geaendert(t, nimmAufVergleich(t, "hb_quelle_q"), "id int", "id bigint"), "name text", "name varchar")
	erwartet := handbuchMeldung(3, 1) + "Position 1, Art Spaltenbeschreibung (Spalte 1: Typ-OID aufgezeichnet 23, Server 20; Spalte 2: Typ-OID aufgezeichnet 25, Server 1043)"
	dir := t.TempDir()
	schreibe := func(name, text string) string {
		pfad := filepath.Join(dir, name)
		if err := os.WriteFile(pfad, []byte(text), 0o600); err != nil {
			t.Fatal(err)
		}
		return pfad
	}
	beide := schreibe("beide.yaml", "play:\n  compare_responses: true\n  diff_all: true\n")
	nurAlle := schreibe("nur.yaml", "play:\n  diff_all: true\n")
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "")
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "")
	play := func(db string, extra ...string) (string, int) {
		leereDatenbank(t, db)
		_, stderr, code := starteBisFrist(t, 30*time.Second, append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", db}, extra...)...)
		return stderr, code
	}
	stderr, code := play("hb_q_option", "--compare-responses", "--diff-all")
	if code != 5 || !strings.Contains(stderr, erwartet) {
		t.Fatalf("Option: Exit-Code %d:\n%s", code, stderr)
	}
	stderr, code = play("hb_q_datei", "--config", beide)
	if code != 5 || !strings.Contains(stderr, erwartet) {
		t.Fatalf("Schlüssel: Exit-Code %d:\n%s", code, stderr)
	}
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "true")
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "true")
	stderr, code = play("hb_q_env")
	if code != 5 || !strings.Contains(stderr, erwartet) {
		t.Fatalf("Umgebungsvariable: Exit-Code %d:\n%s", code, stderr)
	}
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "")
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "")
	for _, f := range []struct {
		name, env string
		args      []string
	}{
		{"Option", "", []string{"--diff-all"}},
		{"Umgebungsvariable", "true", nil},
		{"Datei", "", []string{"--config", nurAlle}},
		{"compare-responses=false", "", []string{"--diff-all", "--compare-responses=false"}},
	} {
		t.Run(f.name, func(t *testing.T) {
			t.Setenv("PGWIRE_RECORDER_DIFF_ALL", f.env)
			stderr, code := play("hb_q_ohne_"+strconv.Itoa(len(f.name)), f.args...)
			if code != 2 || strings.Count(stderr, "PGR-E2001") != 1 || strings.Contains(stderr, "PGR-E2004") || !strings.Contains(stderr, "--diff-all") || !strings.Contains(stderr, "--compare-responses") {
				t.Fatalf("Exit-Code %d:\n%s", code, stderr)
			}
		})
	}
	for _, f := range []struct {
		name, env, meldung string
		args               []string
	}{
		{"record Option", "", "PGR-E2001", []string{"record", "--upstream", "x:5432", "--output", filepath.Join(dir, "o.yaml"), "--diff-all"}},
		{"replay Option", "", "PGR-E2001", []string{"replay", "--input", input, "--diff-all"}},
		{"replay Umgebungsvariable", "true", "PGR-E2001", []string{"replay", "--input", input}},
		{"record Umgebungsvariable", "true", "PGR-E2001", []string{"record", "--upstream", "x:5432", "--output", filepath.Join(dir, "o.yaml")}},
		{"play Schlüssel ungültig", "", "PGR-E2004", []string{"play", "--upstream", "x:5432", "--input", input, "--config", schreibe("p.yaml", "play:\n  diff_all: 1\n")}},
		{"replay Schlüssel", "", "PGR-E2004", []string{"replay", "--input", input, "--config", schreibe("r.yaml", "replay:\n  diff_all: true\n")}},
		{"record Schlüssel", "", "PGR-E2004", []string{"record", "--upstream", "x:5432", "--output", filepath.Join(dir, "o.yaml"), "--config", schreibe("rec.yaml", "record:\n  diff_all: true\n")}},
	} {
		t.Run(f.name, func(t *testing.T) {
			t.Setenv("PGWIRE_RECORDER_DIFF_ALL", f.env)
			_, stderr, code := starteBisFrist(t, 30*time.Second, f.args...)
			if code != 2 || !strings.Contains(stderr, f.meldung) {
				t.Fatalf("Exit-Code %d, erwartet 2 mit %s:\n%s", code, f.meldung, stderr)
			}
		})
	}
}

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Boundary — gegen die reale Instanz steht mit
// --diff-all und --continue-on-error in jeder Interaktion mit Abweichung die
// vollständige Liste (eine Zeile PGR-E5004 je Interaktion), und die Zeile zum Ende
// zählt jede Interaktion mit Abweichung einmal; ohne --continue-on-error endet play
// nach der ersten Interaktion mit Abweichung mit der vollständigen Liste dieser
// Interaktion; der Exit-Code ist 5 (Abnahmeszenario 16, LH-FA-24.a *Alle
// Abweichungen*, *Zusammenspiel*).
func TestE2EPlayHandbuchDiffAllFortsetzung(t *testing.T) {
	leereDatenbank(t, "hb_quelle_f")
	input := geaendert(t, geaendert(t, nimmAufVergleich(t, "hb_quelle_f"), "id int", "id bigint"), "INSERT INTO vgl VALUES (1, 'a'), (2, 'b')", "UPDATE vgl SET id = id")
	stderr, code := vergleichLauf(t, "hb_fort", input, "--diff-all", "--continue-on-error")
	if code != 5 || strings.Count(stderr, "code=PGR-E5004") != 2 || !strings.Contains(stderr, handbuchMeldung(2, 1)+"Position 1, Art Befehl") || !strings.Contains(stderr, handbuchMeldung(3, 1)+"Position 1, Art Spaltenbeschreibung") || zaehler(t, stderr) != "5/5/2" {
		t.Fatalf("mit --continue-on-error: Exit-Code %d, Zähler %s:\n%s", code, zaehler(t, stderr), stderr)
	}
	stderr, code = vergleichLauf(t, "hb_ohne_fort", input, "--diff-all")
	if code != 5 || strings.Count(stderr, "code=PGR-E5004") != 1 || !strings.Contains(stderr, handbuchMeldung(2, 1)+"Position 1, Art Befehl") || zaehler(t, stderr) != "2/2/1" {
		t.Fatalf("ohne --continue-on-error: Exit-Code %d, Zähler %s:\n%s", code, zaehler(t, stderr), stderr)
	}
}
