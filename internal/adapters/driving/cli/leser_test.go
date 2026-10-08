package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// werteJeArt liefert zwei verschiedene gültige Werte je Wertemenge des
// allgemeinen Lesers und einen ungültigen ("" für eine Wertemenge ohne
// ungültigen Wert).
func werteJeArt() map[string][3]string {
	return map[string][3]string{
		"text":          {"wert-a", "wert-b", ""},
		"wahrheitswert": {"true", "false", "1"},
		"dauer":         {"1s", "2m", "5"},
		"stufe":         {"warn", "debug", "INFO"},
	}
}

// leserKommandos sind die Kommandos am allgemeinen Leser.
func leserKommandos() []string { return []string{"record", "replay"} }

// leere setzt die Umgebungsvariablen aller Optionen von kommando leer, also
// nicht gesetzt (LH-FA-17.a).
func leere(t *testing.T, kommando string) {
	t.Helper()
	for _, o := range cli.Optionen(kommando) {
		t.Setenv(o.Env, "")
	}
}

// basis sind kommando und seine Pflichtoptionen außer ohne.
func basis(kommando, ohne string) []string {
	args := []string{kommando}
	for _, o := range cli.Optionen(kommando) {
		if o.Pflicht && o.Name != ohne {
			args = append(args, "--"+o.Name+"=pflicht")
		}
	}
	if kommando == "record" {
		args = append(args, "--output=r.yaml")
	}
	return args
}

func lese(args ...string) (cli.Command, error) {
	return cli.Parse(args, &bytes.Buffer{})
}

// Abdeckung: LH-FA-17/Boundary — für jede Option am allgemeinen Leser von
// record und replay gilt Kommandozeile vor Umgebungsvariable vor Standardwert:
// Die Umgebungsvariable setzt den Wert wie die Option, die Option geht ihr
// vor, ohne beide gilt der Standardwert, und eine Pflichtoption ohne beide ist
// PGR-E2001; eine gesetzte Umgebungsvariable mit ungültigem Wert ist PGR-E2001,
// auch wenn die Kommandozeile die Option setzt (LH-FA-17.a, SPEC-007).
func TestLeserAlleOptionen(t *testing.T) {
	for _, kommando := range leserKommandos() {
		for _, o := range cli.Optionen(kommando) {
			werte, ok := werteJeArt()[o.Art]
			if !ok {
				t.Fatalf("%s --%s: Wertemenge %q ohne Testwerte", kommando, o.Name, o.Art)
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
			} else if standard, err2 := lese(append(args, "--"+o.Name+"="+o.Standard)...); err != nil || err2 != nil || ohne != standard {
				t.Errorf("%s ohne --%s: %#v, %v, erwartet Standard %q: %#v, %v", kommando, o.Name, ohne, err, o.Standard, standard, err2)
			}

			t.Setenv(o.Env, a)
			if got, err := lese(args...); err != nil || got != mitA {
				t.Errorf("%s, %s=%s: %#v, %v, erwartet %#v", kommando, o.Env, a, got, err, mitA)
			}
			if got, err := lese(append(args, "--"+o.Name+"="+b)...); err != nil || got != mitB {
				t.Errorf("%s, %s=%s, --%s=%s: %#v, %v, erwartet die Kommandozeile %#v", kommando, o.Env, a, o.Name, b, got, err, mitB)
			}

			if ungueltig != "" {
				t.Setenv(o.Env, ungueltig)
				if _, err := lese(append(args, "--"+o.Name+"="+b)...); !istUsage(err) || !strings.Contains(err.Error(), o.Env) {
					t.Errorf("%s, %s=%s neben --%s=%s: erwartet %s mit %s, erhalten %v", kommando, o.Env, ungueltig, o.Name, b, model.CodeUsage, o.Env, err)
				}
			}
			t.Setenv(o.Env, "")
		}
	}
}

// Abdeckung: LH-FA-17/Boundary — der Name der Umgebungsvariable einer Option ist
// PGWIRE_RECORDER_ und der Optionsname in Großbuchstaben mit _ statt -; die
// Optionen stehen in der Reihenfolge der Tabelle in LH-FA-17.a.
func TestLeserOptionen(t *testing.T) {
	for kommando, want := range map[string][]string{
		"record": {"listen", "upstream", "shutdown-timeout", "log-level"},
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

// Abdeckung: LH-FA-17/Negative — der Leser bricht beim ersten Fehler ab und
// prüft zuerst die Kommandozeile, dann die Umgebungsvariablen in der
// Reihenfolge der Tabelle, zuletzt die Pflichtoptionen (LH-FA-17.a *Fehler*).
func TestLeserReihenfolge(t *testing.T) {
	leere(t, "replay")
	t.Setenv(cli.EnvShutdownTimeout, "5")
	t.Setenv(cli.EnvLogLevel, "INFO")
	if _, err := lese("replay", "--listen=x", "--input=r.yaml", "--log-level=trace"); !istUsage(err) || strings.Contains(err.Error(), "Umgebungsvariable") {
		t.Errorf("Kommandozeile vor Umgebung: %v", err)
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

// Die Hilfe von record und replay nennt die Umgebungsvariablen der Optionen und
// ihre Priorität; bei record nimmt sie --output und --force aus, die nur die
// Kommandozeile liest (LH-FA-01.a).
func TestLeserHilfe(t *testing.T) {
	for kommando, ausnahme := range map[string]bool{"record": true, "replay": false} {
		text := strings.Join(strings.Fields(hilfe(t, kommando, "--help")), " ")
		if !strings.Contains(text, "über ihre Umgebungsvariable setzbar: PGWIRE_RECORDER_ und der Name in Großbuchstaben mit _ statt -") || !strings.Contains(text, "die Kommandozeile geht ihr vor") {
			t.Errorf("%s: Hilfe ohne Umgebungsvariablen:\n%s", kommando, text)
		}
		if strings.Contains(text, "außer --output und --force") != ausnahme {
			t.Errorf("%s: Ausnahme --output und --force: %v erwartet:\n%s", kommando, ausnahme, text)
		}
	}
	leere(t, "record")
	t.Setenv("PGWIRE_RECORDER_OUTPUT", "r.yaml")
	t.Setenv("PGWIRE_RECORDER_FORCE", "ungültig")
	if _, err := lese("record", "--listen=x", "--upstream=pg:5432"); !istUsage(err) || !strings.Contains(err.Error(), "--output") {
		t.Errorf("record liest --output aus der Umgebung: %v", err)
	}
	if _, err := lese("record", "--listen=x", "--upstream=pg:5432", "--output=r.yaml"); err != nil {
		t.Errorf("record liest --force aus der Umgebung: %v", err)
	}
}
