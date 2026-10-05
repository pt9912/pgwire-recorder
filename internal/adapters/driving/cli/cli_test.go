package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

func TestParseRecord(t *testing.T) {
	cmd, err := Parse([]string{"record", "--listen", ":15432", "--upstream", "pg:5432", "--output", "r.yaml", "--force"}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	want := RecordOptions{Listen: ":15432", Upstream: "pg:5432", Output: "r.yaml", Force: true}
	if cmd.Name != "record" || cmd.Record != want {
		t.Fatalf("erhalten %#v", cmd)
	}
}

func TestParseReplay(t *testing.T) {
	t.Setenv(envFailOnUnconsumed, "")
	cmd, err := Parse([]string{"replay", "--listen", ":15432", "--input", "r.yaml"}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Name != "replay" || cmd.Replay != (ReplayOptions{Listen: ":15432", Input: "r.yaml"}) {
		t.Fatalf("erhalten %#v", cmd)
	}
	if _, err := Parse([]string{"replay", "--listen", ":1"}, &bytes.Buffer{}); err == nil {
		t.Fatal("fehlendes --input angenommen")
	}
}

func TestParseFehler(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"unbekannt"},
		{"record", "--listen", ":1", "--upstream", "pg:5432"},
		{"record", "--listen", ":1", "--upstream", "pg:5432", "--output", "r.yaml", "zusatz"},
		{"record", "--gibtsnicht"},
	} {
		_, err := Parse(args, &bytes.Buffer{})
		var me *model.Error
		if !errors.As(err, &me) || me.Code != model.CodeUsage || me.ExitCode() != 2 {
			t.Errorf("%v: erwartet %s mit Exit-Code 2, erhalten %v", args, model.CodeUsage, err)
		}
	}
}

func TestParseHilfe(t *testing.T) {
	var out bytes.Buffer
	_, err := Parse([]string{"--help"}, &out)
	if !errors.Is(err, ErrHelp) || out.Len() == 0 {
		t.Fatalf("Hilfe: err=%v ausgabe=%q", err, out.String())
	}
}

// replayFail liest replay mit den Zusatzargumenten und liefert
// FailOnUnconsumed oder den Fehler.
func replayFail(args ...string) (bool, error) {
	cmd, err := Parse(append([]string{"replay", "--listen", ":1", "--input", "r.yaml"}, args...), &bytes.Buffer{})
	return cmd.Replay.FailOnUnconsumed, err
}

func istUsage(err error) bool {
	var me *model.Error
	return errors.As(err, &me) && me.Code == model.CodeUsage && me.ExitCode() == 2
}

// --fail-on-unconsumed ist ohne Wert true und nimmt mit = genau true oder
// false; ohne Option und ohne Umgebungsvariable ist sie false; nennt die
// Kommandozeile sie mehrfach, gilt die letzte Angabe.
func TestParseFailOnUnconsumed(t *testing.T) {
	t.Setenv(envFailOnUnconsumed, "")
	for _, fall := range []struct {
		args []string
		want bool
	}{
		{nil, false},
		{[]string{"--fail-on-unconsumed"}, true},
		{[]string{"--fail-on-unconsumed=true"}, true},
		{[]string{"--fail-on-unconsumed=false"}, false},
		{[]string{"--fail-on-unconsumed", "--fail-on-unconsumed=false"}, false},
		{[]string{"--fail-on-unconsumed=false", "--fail-on-unconsumed"}, true},
	} {
		got, err := replayFail(fall.args...)
		if err != nil || got != fall.want {
			t.Errorf("%v: erhalten %v, %v, erwartet %v", fall.args, got, err, fall.want)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — jeder andere Wert von --fail-on-unconsumed,
// auch der leere und 1, ist PGR-E2001.
func TestParseFailOnUnconsumedWerte(t *testing.T) {
	t.Setenv(envFailOnUnconsumed, "")
	for _, wert := range []string{"", "1", "0", "TRUE", "True", "t", "yes"} {
		if _, err := replayFail("--fail-on-unconsumed=" + wert); !istUsage(err) {
			t.Errorf("Wert %q: erwartet %s, erhalten %v", wert, model.CodeUsage, err)
		}
	}
}

// Abdeckung: LH-FA-17/Boundary — PGWIRE_RECORDER_FAIL_ON_UNCONSUMED nimmt true
// oder false, leer gilt sie als nicht gesetzt, jeder andere Wert ist
// PGR-E2001; die Option der Kommandozeile geht ihr vor.
func TestParseFailOnUnconsumedUmgebung(t *testing.T) {
	for _, fall := range []struct {
		env  string
		args []string
		want bool
	}{
		{"true", nil, true},
		{"false", nil, false},
		{"", nil, false},
		{"true", []string{"--fail-on-unconsumed=false"}, false},
		{"false", []string{"--fail-on-unconsumed"}, true},
	} {
		t.Setenv(envFailOnUnconsumed, fall.env)
		got, err := replayFail(fall.args...)
		if err != nil || got != fall.want {
			t.Errorf("Umgebung %q, %v: erhalten %v, %v, erwartet %v", fall.env, fall.args, got, err, fall.want)
		}
	}
	for _, wert := range []string{"1", "TRUE", " true", "ja"} {
		t.Setenv(envFailOnUnconsumed, wert)
		if _, err := replayFail(); !istUsage(err) {
			t.Errorf("Umgebung %q: erwartet %s, erhalten %v", wert, model.CodeUsage, err)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — eine gesetzte
// PGWIRE_RECORDER_FAIL_ON_UNCONSUMED mit ungültigem Wert ist PGR-E2001, auch
// wenn die Kommandozeile die Option setzt; eine leere gilt daneben als nicht
// gesetzt.
func TestParseFailOnUnconsumedUmgebungNebenOption(t *testing.T) {
	for _, args := range [][]string{{"--fail-on-unconsumed"}, {"--fail-on-unconsumed=true"}, {"--fail-on-unconsumed=false"}} {
		t.Setenv(envFailOnUnconsumed, "1")
		if _, err := replayFail(args...); !istUsage(err) {
			t.Errorf("Umgebung \"1\", %v: erwartet %s, erhalten %v", args, model.CodeUsage, err)
		}
		t.Setenv(envFailOnUnconsumed, "")
		if _, err := replayFail(args...); err != nil {
			t.Errorf("leere Umgebung, %v: %v", args, err)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — record kennt
// --fail-on-unconsumed nicht (PGR-E2001) und lässt ihre Umgebungsvariable
// unbeachtet, auch mit ungültigem Wert.
func TestParseFailOnUnconsumedRecord(t *testing.T) {
	t.Setenv(envFailOnUnconsumed, "1")
	record := []string{"record", "--listen", ":1", "--upstream", "pg:5432", "--output", "r.yaml"}
	if _, err := Parse(append(record, "--fail-on-unconsumed"), &bytes.Buffer{}); !istUsage(err) {
		t.Fatalf("record mit --fail-on-unconsumed: %v", err)
	}
	if _, err := Parse(record, &bytes.Buffer{}); err != nil {
		t.Fatalf("record wertet die Umgebungsvariable aus: %v", err)
	}
}

// hilfe ruft Parse und liefert die Ausgabe, wenn Parse die Hilfe meldet; sonst
// bricht es ab.
func hilfe(t *testing.T, args ...string) string {
	t.Helper()
	var out bytes.Buffer
	cmd, err := Parse(args, &out)
	if !errors.Is(err, ErrHelp) || out.Len() == 0 {
		t.Fatalf("%q: erwartet Hilfe, erhalten %#v, %v, Ausgabe %q", args, cmd, err, out.String())
	}
	return out.String()
}

// keineHilfe ruft Parse und bricht ab, wenn Parse die Hilfe meldet oder etwas
// ausgibt; es liefert Kommando und Fehler.
func keineHilfe(t *testing.T, args ...string) (Command, error) {
	t.Helper()
	var out bytes.Buffer
	cmd, err := Parse(args, &out)
	if errors.Is(err, ErrHelp) || out.Len() > 0 {
		t.Fatalf("%q: Hilfe statt Prüfung: %v, Ausgabe %q", args, err, out.String())
	}
	return cmd, err
}

// Eine Hilfe-Angabe geht jeder Prüfung vor (LH-FA-01.a): Bei
// record und replay geben -h, --h, -help und --help, auch mit beliebigem
// =-Wert, an jeder Stelle vor dem ersten --, auch an der Stelle eines
// Optionswerts, die Hilfe des Kommandos aus, ohne die ungültige
// Umgebungsvariable PGWIRE_RECORDER_FAIL_ON_UNCONSUMED, eine unbekannte Option
// davor, einen ungültigen Optionswert oder fehlende Pflichtoptionen zu prüfen.
func TestParseHilfeVorPruefung(t *testing.T) {
	t.Setenv(envFailOnUnconsumed, "1")
	for kommando, args := range map[string][][]string{
		"replay": {
			{"--help"}, {"-h"}, {"--h"}, {"-help"},
			{"--bogus", "--help"},
			{"--fail-on-unconsumed=1", "-h"},
			{"--help=false"}, {"-h=x"}, {"--help="}, {"-help=false"},
			{"--listen", "x", "--input", "-help"},
		},
		"record": {
			{"--help"}, {"--h"},
			{"--bogus", "--help"},
			{"--force=yes", "-h"},
			{"--listen", "x", "--output", "--h=1"},
		},
	} {
		for _, a := range args {
			text := hilfe(t, append([]string{kommando}, a...)...)
			if !strings.Contains(text, "Optionen von "+kommando) || strings.Contains(text, "Kommandos:") {
				t.Errorf("%s %q: nicht die Hilfe von %s:\n%s", kommando, a, kommando, text)
			}
		}
	}
	if text := hilfe(t, "version", "--help"); !strings.Contains(text, "pgwire-recorder version") || strings.Contains(text, "Kommandos:") {
		t.Errorf("version --help: nicht die Hilfe von version:\n%s", text)
	}
}

// Vor dem Kommando, ohne Kommando und nach
// einem unbekannten Kommando gibt eine Hilfe-Angabe die globale Hilfe aus.
func TestParseHilfeGlobal(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"-h"}, {"bogus", "--help"}, {"--listen", "x", "--help"}, {"--h", "replay"}} {
		if text := hilfe(t, args...); !strings.Contains(text, "Kommandos:") {
			t.Errorf("%q: nicht die globale Hilfe:\n%s", args, text)
		}
	}
}

// Abdeckung: LH-FA-01/Negative — keine Hilfe-Angabe sind andere Schreibweisen
// (--Help, ---help, -H), eine Hilfe-Angabe als Teil eines Arguments
// (--input=-help, --input=--help) und eine Hilfe-Angabe nach --; sie werden
// geprüft wie jedes andere Argument, und ein ungültiger Aufruf ist PGR-E2001.
func TestParseKeineHilfe(t *testing.T) {
	t.Setenv(envFailOnUnconsumed, "")
	for _, args := range [][]string{
		{"replay", "--Help"}, {"replay", "---help"}, {"replay", "-H"},
		{"replay", "--listen", "x", "--input", "r.yaml", "--", "--help"},
		{"--", "--help"},
	} {
		if _, err := keineHilfe(t, args...); !istUsage(err) {
			t.Errorf("%q: erwartet %s, erhalten %v", args, model.CodeUsage, err)
		}
	}
	for _, wert := range []string{"-help", "--help"} {
		cmd, err := keineHilfe(t, "replay", "--input="+wert, "--listen", "x")
		if err != nil || cmd.Replay.Input != wert {
			t.Errorf("--input=%s: erhalten %#v, %v", wert, cmd, err)
		}
	}
}

// Abdeckung: LH-FA-01/Negative — das erste -- beendet die Optionen an jeder
// Stelle, auch an der eines Optionswerts: Fehlt einer Option davor dadurch
// ihr Wert, ist das PGR-E2001, ohne die Hilfe-Anforderung des Parsers zu
// nennen; ein Argument danach ist PGR-E2001 (unerwartetes Argument); der Wert
// -- geht nur mit =.
func TestParseEndeDerOptionen(t *testing.T) {
	t.Setenv(envFailOnUnconsumed, "")
	for _, args := range [][]string{
		{"replay", "--input", "--", "--help"},
		{"replay", "--listen", "x", "--input", "--"},
		{"record", "--listen", "x", "--upstream", "pg:5432", "--output", "--"},
	} {
		_, err := keineHilfe(t, args...)
		if !istUsage(err) || strings.Contains(err.Error(), "help requested") || !strings.Contains(err.Error(), "needs an argument") {
			t.Errorf("%q: erwartet %s ohne Wert der Option, erhalten %v", args, model.CodeUsage, err)
		}
	}
	for _, args := range [][]string{
		{"replay", "--listen", "x", "--input", "y", "--", "z"},
		{"record", "--listen", "x", "--upstream", "pg:5432", "--output", "y", "--", "z"},
	} {
		_, err := keineHilfe(t, args...)
		if !istUsage(err) || !strings.Contains(err.Error(), `unerwartetes Argument "z"`) {
			t.Errorf("%q: erwartet %s mit unerwartetem Argument, erhalten %v", args, model.CodeUsage, err)
		}
	}
	cmd, err := keineHilfe(t, "replay", "--listen", "x", "--input=--")
	if err != nil || cmd.Replay.Input != "--" {
		t.Errorf("--input=--: erhalten %#v, %v", cmd, err)
	}
}

// Abdeckung: LH-FA-01/Negative — eine Option --version gibt es nicht, weder
// als Kommando noch nach version (PGR-E2001).
func TestParseVersion(t *testing.T) {
	if cmd, err := keineHilfe(t, "version"); err != nil || cmd.Name != "version" {
		t.Fatalf("version: %#v, %v", cmd, err)
	}
	for _, args := range [][]string{{"--version"}, {"version", "--version"}} {
		if _, err := keineHilfe(t, args...); !istUsage(err) {
			t.Errorf("%q: erwartet %s, erhalten %v", args, model.CodeUsage, err)
		}
	}
}
