package cli_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// fristKommandos liefert die Kommandos mit --shutdown-timeout.
func fristKommandos() []string { return []string{"record", "replay"} }

// fristVon liest kommando mit seinen Pflichtoptionen und den Zusatzargumenten
// und liefert den Wert von --shutdown-timeout oder den Fehler.
func fristVon(kommando string, args ...string) (time.Duration, error) {
	basis := map[string][]string{
		"record": {"record", "--listen", ":1", "--upstream", "pg:5432", "--output", "r.yaml"},
		"replay": {"replay", "--listen", ":1", "--input", "r.yaml"},
	}[kommando]
	cmd, err := cli.Parse(append(append([]string{}, basis...), args...), &bytes.Buffer{})
	if kommando == "record" {
		return cmd.Record.ShutdownTimeout, err
	}
	return cmd.Replay.ShutdownTimeout, err
}

// Abdeckung: LH-FA-17/Boundary, LH-FA-13/Boundary — --shutdown-timeout nimmt 0
// oder eine ganze Zahl mit genau einer Einheit ms, s oder m, 0 mit Einheit ist
// 0; ohne Option gilt 5s (SPEC-046), mehrfach genannt die letzte Angabe, und
// die längste darstellbare Dauer ist zulässig (LH-FA-17.a *Dauer*).
func TestParseShutdownTimeout(t *testing.T) {
	t.Setenv(cli.EnvShutdownTimeout, "")
	t.Setenv(cli.EnvLogLevel, "")
	for _, kommando := range fristKommandos() {
		for _, f := range []struct {
			args []string
			want time.Duration
		}{
			{nil, 5 * time.Second},
			{[]string{"--shutdown-timeout=0"}, 0},
			{[]string{"--shutdown-timeout", "0ms"}, 0},
			{[]string{"--shutdown-timeout=0s"}, 0},
			{[]string{"--shutdown-timeout=0m"}, 0},
			{[]string{"--shutdown-timeout=250ms"}, 250 * time.Millisecond},
			{[]string{"--shutdown-timeout=7s"}, 7 * time.Second},
			{[]string{"--shutdown-timeout=2m"}, 2 * time.Minute},
			{[]string{"--shutdown-timeout=1s", "--shutdown-timeout=3s"}, 3 * time.Second},
			{[]string{"--shutdown-timeout=9223372036854ms"}, 9223372036854 * time.Millisecond},
		} {
			got, err := fristVon(kommando, f.args...)
			if err != nil || got != f.want {
				t.Errorf("%s %v: %v, %v, erwartet %v", kommando, f.args, got, err, f.want)
			}
		}
	}
}

// Abdeckung: LH-FA-17/Negative — jeder andere Wert von --shutdown-timeout ist
// PGR-E2001: leer, negativ, ohne Einheit außer genau 0 (auch `00`), mit
// Nachkommastellen, mit Leerraum, großgeschriebener, unbekannter oder
// zusammengesetzter Einheit und länger als die längste darstellbare Dauer;
// führende Nullen vor einer Einheit sind erlaubt (`05s` ist 5 s, `00s` ist 0)
// (LH-FA-17.a *Dauer*).
func TestParseShutdownTimeoutWerte(t *testing.T) {
	t.Setenv(cli.EnvShutdownTimeout, "")
	t.Setenv(cli.EnvLogLevel, "")
	for _, kommando := range fristKommandos() {
		for _, wert := range []string{
			"", "-1s", "+1s", "5", "00", "1.5s", "1,5s", " 5s", "5s ", "5 s", "5S", "5MS", "5h", "5us",
			"1m30s", "s", "ms", "9223372036855ms", "9223372036854775808s", "99999999999999999999ms",
		} {
			if _, err := fristVon(kommando, "--shutdown-timeout="+wert); !istUsage(err) {
				t.Errorf("%s --shutdown-timeout=%q: erwartet %s, erhalten %v", kommando, wert, model.CodeUsage, err)
			}
		}
		// Führende Nullen vor einer Einheit sind erlaubt, ohne Einheit gilt nur
		// genau 0 (`00` steht oben unter den ungültigen).
		for wert, want := range map[string]time.Duration{"05s": 5 * time.Second, "00s": 0} {
			if got, err := fristVon(kommando, "--shutdown-timeout="+wert); err != nil || got != want {
				t.Errorf("%s --shutdown-timeout=%q: %v, %v, erwartet %v", kommando, wert, got, err, want)
			}
		}
	}
}

// Abdeckung: LH-FA-17/Boundary — PGWIRE_RECORDER_SHUTDOWN_TIMEOUT setzt die
// Frist, leer gilt sie als nicht gesetzt, die Option der Kommandozeile geht ihr
// vor (LH-FA-17.a).
func TestParseShutdownTimeoutUmgebung(t *testing.T) {
	t.Setenv(cli.EnvLogLevel, "")
	for _, kommando := range fristKommandos() {
		for _, f := range []struct {
			env  string
			args []string
			want time.Duration
		}{
			{"2s", nil, 2 * time.Second},
			{"0", nil, 0},
			{"", nil, 5 * time.Second},
			{"", []string{"--shutdown-timeout=1m"}, time.Minute},
			{"2s", []string{"--shutdown-timeout=30ms"}, 30 * time.Millisecond},
		} {
			t.Setenv(cli.EnvShutdownTimeout, f.env)
			got, err := fristVon(kommando, f.args...)
			if err != nil || got != f.want {
				t.Errorf("%s, Umgebung %q, %v: %v, %v, erwartet %v", kommando, f.env, f.args, got, err, f.want)
			}
		}
	}
}

// Abdeckung: LH-FA-17/Negative — eine gesetzte
// PGWIRE_RECORDER_SHUTDOWN_TIMEOUT mit ungültigem Wert ist PGR-E2001, auch wenn
// die Kommandozeile eine gültige Frist setzt (LH-FA-17.a).
func TestParseShutdownTimeoutUmgebungUngueltig(t *testing.T) {
	t.Setenv(cli.EnvLogLevel, "")
	for _, kommando := range fristKommandos() {
		for _, wert := range []string{"5", "-1s", "1m30s", "5S"} {
			t.Setenv(cli.EnvShutdownTimeout, wert)
			for _, args := range [][]string{nil, {"--shutdown-timeout=1s"}} {
				if _, err := fristVon(kommando, args...); !istUsage(err) {
					t.Errorf("%s, Umgebung %q, %v: erwartet %s, erhalten %v", kommando, wert, args, model.CodeUsage, err)
				}
			}
		}
	}
}

// Die Hilfe von record und replay nennt --shutdown-timeout und seine Umgebungsvariable,
// auch wenn diese ungültig ist (LH-FA-01.a).
func TestParseShutdownTimeoutHilfe(t *testing.T) {
	t.Setenv(cli.EnvShutdownTimeout, "ungültig")
	for _, kommando := range fristKommandos() {
		if text := hilfe(t, kommando, "--help"); !strings.Contains(text, "--shutdown-timeout") || !strings.Contains(text, cli.EnvShutdownTimeout) {
			t.Errorf("%s: Hilfe ohne --shutdown-timeout: %q", kommando, text)
		}
	}
}
