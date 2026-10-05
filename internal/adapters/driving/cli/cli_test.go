package cli

import (
	"bytes"
	"errors"
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

// Abdeckung: LH-FA-17/Happy, LH-FA-03/Happy — --fail-on-unconsumed ist ohne
// Wert true und nimmt mit = genau true oder false; ohne Option und ohne
// Umgebungsvariable ist sie false; nennt die Kommandozeile sie mehrfach, gilt
// die letzte Angabe.
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

// Abdeckung: LH-FA-17/Negative, LH-FA-03/Negative — record kennt
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
