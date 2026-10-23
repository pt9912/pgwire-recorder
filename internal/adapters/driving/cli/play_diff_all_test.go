package cli_test

import (
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// diffAll liest play mit --upstream, --input und args und liefert die
// Einspielvorgaben.
func diffAll(args ...string) (cli.Einspielvorgaben, error) {
	got, err := playMit(append([]string{"--upstream=pg:1"}, args...)...)
	return got.Einspielen, err
}

// Abdeckung: LH-FA-24/Happy, LH-FA-17/Boundary, LH-FA-17/Negative — --diff-all setzt
// mit --compare-responses genau DiffAll: ohne Wert, mit =true, aus der
// Umgebungsvariable PGWIRE_RECORDER_DIFF_ALL und aus dem Schlüssel diff_all in
// play:; mit =false und ohne Quelle ist es false; die Kommandozeile geht der
// Umgebung vor, die Umgebung der Datei; ein anderer Wert (1, True, yes, leer) ist
// als Option oder Umgebungsvariable PGR-E2001, als Schlüssel PGR-E2004, auch wenn
// eine Quelle davor die Option setzt (LH-FA-24.a *Alle Abweichungen*, *Quellen und
// Werte*).
func TestPlayDiffAllQuellen(t *testing.T) {
	want := cli.Einspielvorgaben{CompareResponses: true, DiffAll: true}
	datei := func(s string) string { return "--config=" + schreibe(t, "play:\n"+s) }
	leere(t, "play")
	for _, args := range [][]string{{"--compare-responses", "--diff-all"}, {"--compare-responses", "--diff-all=true"}, {datei("  compare_responses: true\n  diff_all: true\n")}} {
		if got, err := diffAll(args...); err != nil || got != want {
			t.Errorf("%q: %#v, %v, erwartet %#v", args, got, err, want)
		}
	}
	for _, args := range [][]string{{"--compare-responses"}, {"--compare-responses", "--diff-all=false"}, {"--diff-all=false"}, nil} {
		if got, err := diffAll(args...); err != nil || got.DiffAll || got.CompareResponses != (len(args) > 0 && args[0] == "--compare-responses") {
			t.Errorf("%q: %#v, %v, erwartet DiffAll false", args, got, err)
		}
	}
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "true")
	if got, err := diffAll("--compare-responses"); err != nil || got != want {
		t.Errorf("PGWIRE_RECORDER_DIFF_ALL=true: %#v, %v", got, err)
	}
	if got, err := diffAll("--compare-responses", "--diff-all=false"); err != nil || got.DiffAll {
		t.Errorf("Option false vor Umgebung true: %#v, %v", got, err)
	}
	if got, err := diffAll(datei("  compare_responses: true\n  diff_all: false\n")); err != nil || got != want {
		t.Errorf("Umgebung true vor Datei false: %#v, %v", got, err)
	}
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "false")
	if got, err := diffAll("--compare-responses", datei("  diff_all: true\n")); err != nil || got.DiffAll {
		t.Errorf("Umgebung false vor Datei true: %#v, %v", got, err)
	}
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "")
	for _, w := range []string{"1", "True", "yes", ""} {
		if _, err := diffAll("--compare-responses", "--diff-all="+w); !istUsage(err) || !strings.Contains(err.Error(), "diff-all") {
			t.Errorf("--diff-all=%q: erwartet %s mit dem Namen, erhalten %v", w, model.CodeUsage, err)
		}
	}
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "GEHEIM")
	if _, err := diffAll("--compare-responses", "--diff-all=false"); !istUsage(err) || !strings.Contains(err.Error(), "PGWIRE_RECORDER_DIFF_ALL") || strings.Contains(err.Error(), "GEHEIM") {
		t.Errorf("ungültige Umgebungsvariable neben der Option: %v", err)
	}
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "")
	for _, w := range []string{"yes", "1", `""`} {
		if _, err := diffAll("--compare-responses", "--diff-all=false", datei("  diff_all: "+w+"\n")); !istDatei(err) || !strings.Contains(err.Error(), "play.diff_all") {
			t.Errorf("diff_all: %s in der Datei: erwartet %s mit dem Schlüssel, erhalten %v", w, model.CodeConfigFile, err)
		}
	}
}

// Abdeckung: LH-FA-24/Negative, LH-FA-17/Boundary — ist nach der Zusammenführung
// der Quellen --diff-all wirksam true und --compare-responses wirksam false,
// gleich ob ausdrücklich oder aus welcher Quelle, ist das PGR-E2001 (kein
// PGR-E2004, auch wenn die Datei --diff-all setzt), und die Meldung nennt beide
// Optionen, nie einen Wert; --diff-all=false ist auch ohne --compare-responses
// gültig; die Kombination steht vor jeder anderen Kombination von play
// (LH-FA-24.a *Ohne --compare-responses*).
func TestPlayDiffAllOhneVergleich(t *testing.T) {
	datei := func(s string) string { return "--config=" + schreibe(t, "play:\n"+s) }
	leere(t, "play")
	pruefe := func(name string, args ...string) {
		t.Helper()
		_, err := diffAll(args...)
		if !istUsage(err) || istDatei(err) || !strings.Contains(err.Error(), "--diff-all") || !strings.Contains(err.Error(), "--compare-responses") || strings.Contains(err.Error(), "true") || strings.Contains(err.Error(), "false") {
			t.Errorf("%s: erwartet %s mit beiden Optionen und ohne Wert, erhalten %v", name, model.CodeUsage, err)
		}
	}
	pruefe("Option allein", "--diff-all")
	pruefe("compare-responses=false", "--compare-responses=false", "--diff-all")
	pruefe("diff-all aus der Datei", datei("  diff_all: true\n"))
	pruefe("compare-responses=false aus der Datei", datei("  compare_responses: false\n"), "--diff-all")
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "true")
	pruefe("diff-all aus der Umgebung", "--compare-responses=false")
	pruefe("diff-all aus der Umgebung, compare-responses aus der Datei false", datei("  compare_responses: false\n"))
	t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "")
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "false")
	pruefe("compare-responses=false aus der Umgebung", "--diff-all")
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "")
	if got, err := diffAll("--diff-all=false"); err != nil || got.DiffAll {
		t.Errorf("--diff-all=false ohne --compare-responses: %#v, %v", got, err)
	}
	if got, err := diffAll(datei("  diff_all: false\n")); err != nil || got.DiffAll {
		t.Errorf("diff_all: false ohne compare_responses: %#v, %v", got, err)
	}
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "true")
	if got, err := diffAll("--diff-all"); err != nil || !got.DiffAll || !got.CompareResponses {
		t.Errorf("compare-responses aus der Umgebung: %#v, %v", got, err)
	}
	t.Setenv("PGWIRE_RECORDER_COMPARE_RESPONSES", "")

	// Reihenfolge: vor --upstream-ca ohne TLS, vor einer Zeitkombination und vor
	// einem Fehler von --upstream.
	for name, andere := range map[string][]string{
		"--upstream-ca ohne TLS":  {"--upstream-ca=ca.pem"},
		"--timing-mode ohne Zeit": {"--timing-mode=absolute"},
		"--upstream ohne Eintrag": {"--upstream=name", datei("  continue_on_error: false\n")},
	} {
		_, err := playMit(append([]string{"--upstream=pg:1", "--diff-all"}, andere...)...)
		if name == "--upstream ohne Eintrag" {
			_, err = playMit(append([]string{"--diff-all"}, andere...)...)
		}
		if !istUsage(err) || !strings.Contains(err.Error(), "--diff-all") {
			t.Errorf("%s: erwartet die Meldung von --diff-all zuerst, erhalten %v", name, err)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — record und replay kennen --diff-all nicht: auf
// der Kommandozeile ist es eine unbekannte Option (PGR-E2001), die
// Umgebungsvariable PGWIRE_RECORDER_DIFF_ALL ist PGR-E2001 mit dem Namen und ohne
// den Wert, der Schlüssel diff_all in ihrem Abschnitt PGR-E2004 (LH-FA-24.a *Andere
// Kommandos*).
func TestRecordReplayDiffAll(t *testing.T) {
	for _, k := range []string{"record", "replay"} {
		leere(t, k)
		args := append(basis(k, ""), "--diff-all")
		if _, err := lese(args...); !istUsage(err) || !strings.Contains(err.Error(), "flag provided but not defined") {
			t.Errorf("%s --diff-all: %v", k, err)
		}
		t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "GEHEIM")
		if _, err := lese(basis(k, "")...); !istUsage(err) || !strings.Contains(err.Error(), "PGWIRE_RECORDER_DIFF_ALL") || strings.Contains(err.Error(), "GEHEIM") {
			t.Errorf("%s mit PGWIRE_RECORDER_DIFF_ALL: %v", k, err)
		}
		t.Setenv("PGWIRE_RECORDER_DIFF_ALL", "")
		if _, err := lese(append(basis(k, ""), "--config="+schreibe(t, k+":\n  diff_all: true\n"))...); !istDatei(err) || !strings.Contains(err.Error(), k+".diff_all") {
			t.Errorf("Schlüssel diff_all in %s: %v", k, err)
		}
	}
}

// Abdeckung: LH-FA-01/Boundary, LH-FA-20/Boundary — die Hilfe von play trägt zu
// --diff-all drei Aussagen und die Umgebungsvariable: alle Abweichungen statt der
// ersten, höchstens 20 Positionen, nur mit --compare-responses; die globale Hilfe
// und die Hilfe von record und replay nennen die Option nicht (LH-FA-20.a, Hilfe
// von play).
func TestParsePlayDiffAllHilfe(t *testing.T) {
	text := strings.Join(strings.Fields(hilfe(t, "play", "--help")), " ")
	i := strings.Index(text, "--diff-all")
	if i < 0 {
		t.Fatalf("Hilfe von play ohne --diff-all:\n%s", text)
	}
	abschnitt := text[i:]
	if j := strings.Index(abschnitt, " --log-level"); j > 0 {
		abschnitt = abschnitt[:j]
	}
	for _, teil := range []string{"meldet play je Interaktion alle Abweichungen statt der ersten", "höchstens 20 Positionen", "nur mit --compare-responses", "PGWIRE_RECORDER_DIFF_ALL"} {
		if !strings.Contains(abschnitt, teil) {
			t.Errorf("Hilfe von play zu --diff-all ohne %q: %s", teil, abschnitt)
		}
	}
	for name, args := range map[string][]string{"global": {"--help"}, "record": {"record", "--help"}, "replay": {"replay", "--help"}} {
		if t2 := hilfe(t, args...); strings.Contains(t2, "diff-all") || strings.Contains(t2, "DIFF_ALL") {
			t.Errorf("Hilfe %s nennt --diff-all:\n%s", name, t2)
		}
	}
}
