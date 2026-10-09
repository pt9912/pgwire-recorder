package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// werteJeArt liefert zwei verschiedene gültige Werte je Wertemenge des
// allgemeinen Lesers und einen ungültigen nicht leeren ("" für eine
// Wertemenge, die jeden nicht leeren Wert annimmt). Der leere Wert ist in
// jeder Wertemenge ungültig (LH-FA-17.a).
func werteJeArt() map[string][3]string {
	return map[string][3]string{
		"text":          {"wert-a", "wert-b", ""},
		"wahrheitswert": {"true", "false", "1"},
		"dauer":         {"1s", "2m", "5"},
		"stufe":         {"warn", "debug", "INFO"},
	}
}

// tabelle ist je Option von record und replay der Default aus der
// Optionstabelle in LH-FA-17.a; "" heißt Pflicht.
func tabelle() map[string]string {
	return map[string]string{
		"listen":             "",
		"upstream":           "",
		"output":             "",
		"force":              "false",
		"input":              "",
		"fail-on-unconsumed": "false",
		"shutdown-timeout":   "5s",
		"log-level":          "info",
	}
}

// alleOptionen sind die Namen aller Optionen der Optionstabelle in
// LH-FA-17.a, gleich welchen Kommandos.
func alleOptionen() []string {
	return []string{
		"listen", "upstream", "output", "force", "format", "record-timing", "record-empty-sessions",
		"input", "fail-on-unconsumed", "shutdown-timeout", "session-assignment", "user", "database",
		"continue-on-error", "allow-recorded-errors", "upstream-tls", "finish-session-on-interrupt",
		"keep-timing", "timing-mode", "timing-reference", "tls-cert", "tls-key", "allow-plaintext",
		"upstream-ca", "compare-responses", "config", "log-level",
	}
}

// leserKommandos sind die Kommandos am allgemeinen Leser.
func leserKommandos() []string { return []string{"record", "replay"} }

// leere setzt die Umgebungsvariablen aller Optionen von kommando und
// PGWIRE_RECORDER_CONFIG leer, also nicht gesetzt (LH-FA-17.a).
func leere(t *testing.T, kommando string) {
	t.Helper()
	for _, o := range cli.Optionen(kommando) {
		t.Setenv(o.Env, "")
	}
	t.Setenv(cli.EnvConfig, "")
}

// schreibe legt eine Datei mit inhalt in einem eigenen Verzeichnis des Tests
// an und liefert ihren Pfad.
func schreibe(t *testing.T, inhalt string) string {
	t.Helper()
	pfad := filepath.Join(t.TempDir(), "konfiguration.yaml")
	if err := os.WriteFile(pfad, []byte(inhalt), 0o600); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// schluesselIn ist der Text einer Datei, die name auf wert setzt: auf der
// obersten Ebene, wenn abschnitt leer ist, sonst im Abschnitt.
func schluesselIn(abschnitt, name, wert string) string {
	if abschnitt == "" {
		return name + ": " + wert + "\n"
	}
	return abschnitt + ":\n  " + name + ": " + wert + "\n"
}

// ortDerDatei ist nach LH-FA-17.a der Abschnitt und der Schlüssel einer Option
// in der Datei: log_level auf der obersten Ebene, jede andere im Abschnitt des
// Kommandos; der Schlüssel ist der Name mit _ statt -.
func ortDerDatei(kommando, name string) (string, string) {
	schluessel := strings.ReplaceAll(name, "-", "_")
	if name == "log-level" {
		return "", schluessel
	}
	return kommando, schluessel
}

// istDatei meldet, ob err PGR-E2004 mit Exit-Code 2 ist.
func istDatei(err error) bool {
	return hatCode(err, model.CodeConfigFile)
}

// basis sind kommando und seine Pflichtoptionen außer ohne.
func basis(kommando, ohne string) []string {
	args := []string{kommando}
	for _, o := range cli.Optionen(kommando) {
		if o.Pflicht && o.Name != ohne {
			args = append(args, "--"+o.Name+"=pflicht")
		}
	}
	return args
}

func lese(args ...string) (cli.Command, error) {
	return cli.Parse(args, &bytes.Buffer{})
}

// Abdeckung: LH-FA-17/Boundary, LH-FA-17/Negative — für jede Option am
// allgemeinen Leser von record und replay gilt Kommandozeile vor
// Umgebungsvariable vor Konfigurationsdatei vor Standardwert: Umgebungsvariable
// und Datei setzen den Wert wie die Option, die Option geht beiden vor, die
// Umgebungsvariable der Datei, ohne alle drei gilt der Default der
// Optionstabelle, und eine Pflichtoption ohne alle drei ist PGR-E2001; eine
// gesetzte Umgebungsvariable mit ungültigem Wert ist PGR-E2001, ein ungültiger
// oder leerer Wert in der Datei PGR-E2004, der den Schlüssel nennt und nicht den
// Wert, beide auch, wenn die Kommandozeile die Option setzt; der Schlüssel an
// der anderen Ebene ist PGR-E2004; ein leerer Wert auf der Kommandozeile ist
// PGR-E2001 und nennt die Option, auch neben gesetzter Umgebungsvariable
// (LH-FA-17.a, SPEC-007).
func TestLeserAlleOptionen(t *testing.T) {
	for _, kommando := range leserKommandos() {
		for _, o := range cli.Optionen(kommando) {
			werte, ok := werteJeArt()[o.Art]
			if !ok {
				t.Fatalf("%s --%s: Wertemenge %q ohne Testwerte", kommando, o.Name, o.Art)
			}
			standard, ok := tabelle()[o.Name]
			if !ok {
				t.Fatalf("%s --%s: kein Default aus der Optionstabelle im Test", kommando, o.Name)
			}
			if o.Pflicht != (standard == "") {
				t.Errorf("%s --%s: Pflicht %v gegen die Optionstabelle", kommando, o.Name, o.Pflicht)
			}
			a, b, ungueltig := werte[0], werte[1], werte[2]
			args := basis(kommando, o.Name)
			leere(t, kommando)

			mitA, err := lese(append(args, "--"+o.Name+"="+a)...)
			if err != nil {
				t.Fatalf("%s --%s=%s: %v", kommando, o.Name, a, err)
			}
			mitB, err := lese(append(args, "--"+o.Name+"="+b)...)
			if err != nil || mitA == mitB {
				t.Fatalf("%s --%s: %s und %s ergeben dasselbe oder einen Fehler: %#v, %v", kommando, o.Name, a, b, mitB, err)
			}

			ohne, err := lese(args...)
			if o.Pflicht {
				if !istUsage(err) || !strings.Contains(err.Error(), "--"+o.Name) {
					t.Errorf("%s ohne Pflichtoption --%s: %#v, %v", kommando, o.Name, ohne, err)
				}
			} else if mitStandard, err2 := lese(append(args, "--"+o.Name+"="+standard)...); err != nil || err2 != nil || ohne != mitStandard {
				t.Errorf("%s ohne --%s: %#v, %v, erwartet Default %q der Tabelle: %#v, %v", kommando, o.Name, ohne, err, standard, mitStandard, err2)
			}

			t.Setenv(o.Env, a)
			if got, err := lese(args...); err != nil || got != mitA {
				t.Errorf("%s, %s=%s: %#v, %v, erwartet %#v", kommando, o.Env, a, got, err, mitA)
			}
			if got, err := lese(append(args, "--"+o.Name+"="+b)...); err != nil || got != mitB {
				t.Errorf("%s, %s=%s, --%s=%s: %#v, %v, erwartet die Kommandozeile %#v", kommando, o.Env, a, o.Name, b, got, err, mitB)
			}
			_, err = lese(append(args, "--"+o.Name+"=")...)
			if !istUsage(err) || !strings.Contains(err.Error(), o.Name) || strings.Contains(err.Error(), o.Env) || strings.Contains(err.Error(), "Pflichtoption") {
				t.Errorf("%s, %s=%s, --%s=: erwartet %s mit der Option, erhalten %v", kommando, o.Env, a, o.Name, model.CodeUsage, err)
			}

			if ungueltig != "" {
				t.Setenv(o.Env, ungueltig)
				if _, err := lese(append(args, "--"+o.Name+"="+b)...); !istUsage(err) || !strings.Contains(err.Error(), o.Env) {
					t.Errorf("%s, %s=%s neben --%s=%s: erwartet %s mit %s, erhalten %v", kommando, o.Env, ungueltig, o.Name, b, model.CodeUsage, o.Env, err)
				}
			}
			t.Setenv(o.Env, "")
			dreiQuellen(t, kommando, o, args, mitA, mitB, b, ungueltig)
		}
	}
}

// dreiQuellen prüft die Datei als dritte Quelle für Option o: Die Datei mit
// dem ersten gültigen Wert ergibt mitA, mit Umgebungsvariable oder
// Kommandozeile auf dem zweiten mitB; ein ungültiger oder leerer Wert der
// Datei ist PGR-E2004 mit dem Schlüssel und ohne den Wert, auch neben der
// Kommandozeile; derselbe Schlüssel an der anderen Ebene ist PGR-E2004.
func dreiQuellen(t *testing.T, kommando string, o cli.Option, args []string, mitA, mitB cli.Command, b, ungueltig string) {
	t.Helper()
	abschnitt, schluessel := ortDerDatei(kommando, o.Name)
	stelle := schluessel
	if abschnitt != "" {
		stelle = abschnitt + "." + schluessel
	}
	a := werteJeArt()[o.Art][0]
	mitDatei := func(inhalt string, extra ...string) (cli.Command, error) {
		return lese(append(append(append([]string{}, args...), "--config="+schreibe(t, inhalt)), extra...)...)
	}
	if got, err := mitDatei(schluesselIn(abschnitt, schluessel, a)); err != nil || got != mitA {
		t.Errorf("%s, Datei %s: %s: %#v, %v, erwartet %#v", kommando, stelle, a, got, err, mitA)
	}
	t.Setenv(o.Env, b)
	if got, err := mitDatei(schluesselIn(abschnitt, schluessel, a)); err != nil || got != mitB {
		t.Errorf("%s, %s=%s, Datei %s: %s: %#v, %v, erwartet die Umgebungsvariable %#v", kommando, o.Env, b, stelle, a, got, err, mitB)
	}
	t.Setenv(o.Env, "")
	if got, err := mitDatei(schluesselIn(abschnitt, schluessel, a), "--"+o.Name+"="+b); err != nil || got != mitB {
		t.Errorf("%s, --%s=%s, Datei %s: %s: %#v, %v, erwartet die Kommandozeile %#v", kommando, o.Name, b, stelle, a, got, err, mitB)
	}
	schlecht := []string{`""`}
	if ungueltig != "" {
		schlecht = append(schlecht, ungueltig+"GEHEIM")
	}
	for _, w := range schlecht {
		_, err := mitDatei(schluesselIn(abschnitt, schluessel, w), "--"+o.Name+"="+b)
		if !istDatei(err) || !strings.Contains(err.Error(), stelle) || strings.Contains(err.Error(), "GEHEIM") {
			t.Errorf("%s, Datei %s: %s neben --%s=%s: erwartet %s mit %s ohne den Wert, erhalten %v", kommando, stelle, w, o.Name, b, model.CodeConfigFile, stelle, err)
		}
	}
	anderer := kommando
	if abschnitt != "" {
		anderer = ""
	}
	if _, err := mitDatei(schluesselIn(anderer, schluessel, a)); !istDatei(err) || !strings.Contains(err.Error(), schluessel) {
		t.Errorf("%s, Schlüssel %s auf der falschen Ebene: erwartet %s, erhalten %v", kommando, schluessel, model.CodeConfigFile, err)
	}
}

// Abdeckung: LH-FA-17/Negative — ein leerer Wert auf der Kommandozeile ist
// gesetzt und ungültig, geprüft mit der Kommandozeile vor den
// Umgebungsvariablen: Die Meldung nennt die Option, weder eine fehlende
// Pflichtoption noch eine Umgebungsvariable (LH-FA-17.a).
func TestLeserLeererWert(t *testing.T) {
	leere(t, "record")
	t.Setenv("PGWIRE_RECORDER_LISTEN", "127.0.0.1:1")
	t.Setenv(cli.EnvShutdownTimeout, "x")
	_, err := lese("record", "--listen=", "--upstream", "h:1", "--output", "r.yaml")
	if !istUsage(err) || !strings.Contains(err.Error(), "listen") || strings.Contains(err.Error(), "Pflichtoption") || strings.Contains(err.Error(), "PGWIRE_RECORDER_") {
		t.Fatalf("--listen= neben PGWIRE_RECORDER_LISTEN: erwartet %s mit listen, erhalten %v", model.CodeUsage, err)
	}
}

// Abdeckung: LH-FA-17/Boundary — der Name der Umgebungsvariable einer Option ist
// PGWIRE_RECORDER_ und der Optionsname in Großbuchstaben mit _ statt -; die
// Optionen stehen in der Reihenfolge der Tabelle in LH-FA-17.a.
func TestLeserOptionen(t *testing.T) {
	for kommando, want := range map[string][]string{
		"record": {"listen", "upstream", "output", "force", "shutdown-timeout", "log-level"},
		"replay": {"listen", "input", "fail-on-unconsumed", "shutdown-timeout", "log-level"},
	} {
		var got []string
		for _, o := range cli.Optionen(kommando) {
			got = append(got, o.Name)
			if env := "PGWIRE_RECORDER_" + strings.ToUpper(strings.ReplaceAll(o.Name, "-", "_")); o.Env != env {
				t.Errorf("%s --%s: Umgebungsvariable %s, erwartet %s", kommando, o.Name, o.Env, env)
			}
		}
		if strings.Join(got, " ") != strings.Join(want, " ") {
			t.Errorf("%s: Optionen %v, erwartet %v", kommando, got, want)
		}
	}
	for _, env := range []string{cli.EnvFailOnUnconsumed, cli.EnvLogLevel, cli.EnvShutdownTimeout} {
		if !strings.HasPrefix(env, "PGWIRE_RECORDER_") {
			t.Errorf("Umgebungsvariable %s ohne Präfix", env)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — die Kommandozeile kennt genau die Optionen am
// allgemeinen Leser und --config: Jede andere Option der Optionstabelle ist bei
// record und replay eine unbekannte Option (PGR-E2001), mit und ohne Wert,
// sodass keine Option am Leser vorbei angemeldet ist.
func TestLeserNurAngemeldete(t *testing.T) {
	for _, kommando := range leserKommandos() {
		leere(t, kommando)
		angemeldet := map[string]bool{"config": true}
		for _, o := range cli.Optionen(kommando) {
			angemeldet[o.Name] = true
		}
		for _, name := range alleOptionen() {
			if angemeldet[name] {
				continue
			}
			for _, arg := range []string{"--" + name, "--" + name + "=wert"} {
				_, err := lese(append(basis(kommando, ""), arg)...)
				if !istUsage(err) || !strings.Contains(err.Error(), "flag provided but not defined") {
					t.Errorf("%s %s: erwartet unbekannte Option, erhalten %v", kommando, arg, err)
				}
			}
		}
	}
}

// Abdeckung: LH-FA-17/Negative — der Leser bricht beim ersten Fehler ab und
// prüft zuerst die Kommandozeile, auch ein unerwartetes Argument, dann die
// Umgebungsvariablen in der Reihenfolge der Tabelle, zuletzt die
// Pflichtoptionen (LH-FA-17.a *Fehler*).
func TestLeserReihenfolge(t *testing.T) {
	leere(t, "replay")
	t.Setenv(cli.EnvShutdownTimeout, "5")
	t.Setenv(cli.EnvLogLevel, "INFO")
	if _, err := lese("replay", "--listen=x", "--input=r.yaml", "--log-level=trace"); !istUsage(err) || strings.Contains(err.Error(), "Umgebungsvariable") {
		t.Errorf("Kommandozeile vor Umgebung: %v", err)
	}
	for _, args := range [][]string{
		{"replay", "--listen=x", "--input=r.yaml", "zusatz"},
		{"replay", "--listen=x", "--input=r.yaml", "--", "zusatz"},
	} {
		if _, err := lese(args...); !istUsage(err) || !strings.Contains(err.Error(), `unerwartetes Argument "zusatz"`) {
			t.Errorf("%q: unerwartetes Argument vor Umgebung: %v", args, err)
		}
	}
	_, err := lese("replay", "--listen=x", "--input=r.yaml")
	if !istUsage(err) || !strings.Contains(err.Error(), cli.EnvShutdownTimeout) || strings.Contains(err.Error(), cli.EnvLogLevel) {
		t.Errorf("Umgebung in der Reihenfolge der Tabelle: %v", err)
	}
	t.Setenv(cli.EnvShutdownTimeout, "")
	t.Setenv(cli.EnvFailOnUnconsumed, "1")
	if _, err := lese("replay"); !istUsage(err) || !strings.Contains(err.Error(), cli.EnvFailOnUnconsumed) {
		t.Errorf("Umgebung vor Pflichtoption: %v", err)
	}
	t.Setenv(cli.EnvFailOnUnconsumed, "")
	t.Setenv(cli.EnvLogLevel, "")
	if _, err := lese("replay"); !istUsage(err) || !strings.Contains(err.Error(), "--listen") {
		t.Errorf("Pflichtoptionen in der Reihenfolge der Tabelle: %v", err)
	}
}

// Abdeckung: LH-FA-17/Boundary — eine Umgebungsvariable mit dem Präfix, deren
// Name keiner Option des Kommandos entspricht, bleibt unbeachtet, auch mit
// einem Wert, der für eine Option ungültig wäre (LH-FA-17.a).
func TestLeserFremdeUmgebung(t *testing.T) {
	for _, kommando := range leserKommandos() {
		leere(t, kommando)
		for _, env := range []string{"PGWIRE_RECORDER_BOGUS", "PGWIRE_RECORDER_FORMAT", "PGWIRE_RECORDER_SESSION_ASSIGNMENT", "PGWIRE_RECORDER_KEEP_TIMING", "PGWIRE_RECORDER_"} {
			t.Setenv(env, "ungültig")
		}
		if _, err := lese(basis(kommando, "")...); err != nil {
			t.Errorf("%s: fremde Umgebungsvariable ausgewertet: %v", kommando, err)
		}
	}
}

// Die Hilfe von record und replay nennt die Umgebungsvariablen der Optionen, ihre
// Priorität und dass ein leerer Wert auf der Kommandozeile ungültig ist, ohne
// Ausnahme (LH-FA-01.a). Bei record gilt das auch für --output und --force:
// PGWIRE_RECORDER_OUTPUT allein setzt --output, PGWIRE_RECORDER_FORCE=ja und
// --force=1 sind PGR-E2001 (LH-FA-17.a).
func TestLeserHilfe(t *testing.T) {
	for _, kommando := range leserKommandos() {
		text := strings.Join(strings.Fields(hilfe(t, kommando, "--help")), " ")
		for _, satz := range []string{
			"Jede Option ist auch über ihre Umgebungsvariable setzbar: PGWIRE_RECORDER_ und der Name in Großbuchstaben mit _ statt -",
			"die Kommandozeile geht ihr vor",
			"Ein leerer Wert auf der Kommandozeile ist ungültig",
		} {
			if !strings.Contains(text, satz) {
				t.Errorf("%s: Hilfe ohne %q:\n%s", kommando, satz, text)
			}
		}
		if strings.Contains(text, "außer") {
			t.Errorf("%s: Hilfe nennt eine Ausnahme:\n%s", kommando, text)
		}
	}
	leere(t, "record")
	t.Setenv("PGWIRE_RECORDER_OUTPUT", "aus-der-umgebung.yaml")
	cmd, err := lese("record", "--listen=x", "--upstream=pg:5432")
	if err != nil || cmd.Record.Output != "aus-der-umgebung.yaml" {
		t.Errorf("PGWIRE_RECORDER_OUTPUT allein: %#v, %v", cmd, err)
	}
	t.Setenv("PGWIRE_RECORDER_FORCE", "ja")
	if _, err := lese("record", "--listen=x", "--upstream=pg:5432", "--output=r.yaml", "--force"); !istUsage(err) || !strings.Contains(err.Error(), "PGWIRE_RECORDER_FORCE") {
		t.Errorf("PGWIRE_RECORDER_FORCE=ja: %v", err)
	}
	t.Setenv("PGWIRE_RECORDER_FORCE", "")
	if _, err := lese("record", "--listen=x", "--upstream=pg:5432", "--output=r.yaml", "--force=1"); !istUsage(err) {
		t.Errorf("--force=1: %v", err)
	}
}
