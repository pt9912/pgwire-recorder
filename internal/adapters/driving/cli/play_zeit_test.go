package cli_test

import (
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// zeit sind die Einspielvorgaben mit KeepTiming und den Werten modus und bezug.
func zeit(modus, bezug string) cli.Einspielvorgaben {
	return cli.Einspielvorgaben{KeepTiming: true, TimingMode: modus, TimingReference: bezug}
}

// Abdeckung: LH-FA-21/Happy, LH-FA-17/Boundary — --keep-timing, --timing-mode und
// --timing-reference setzen je ihr Feld der Einspielvorgaben, jede aus Option,
// Umgebungsvariable und Schlüssel im Abschnitt play:; mit --keep-timing allein gilt
// relative ohne Bezugspunkt, mit absolute gilt connect, wo --timing-reference nichts
// setzt; ohne --keep-timing sind alle drei Felder leer (LH-FA-21.a *Einspielen*,
// LH-FA-17.a).
func TestPlayZeitWerte(t *testing.T) {
	for _, f := range []struct {
		name string
		args []string
		want cli.Einspielvorgaben
	}{
		{"ohne Optionen", nil, cli.Einspielvorgaben{}},
		{"--keep-timing=false", []string{"--keep-timing=false"}, cli.Einspielvorgaben{}},
		{"--keep-timing", []string{"--keep-timing"}, zeit("relative", "")},
		{"--keep-timing=true", []string{"--keep-timing=true"}, zeit("relative", "")},
		{"relative", []string{"--keep-timing", "--timing-mode=relative"}, zeit("relative", "")},
		{"absolute", []string{"--keep-timing", "--timing-mode=absolute"}, zeit("absolute", "connect")},
		{"absolute, connect", []string{"--keep-timing", "--timing-mode=absolute", "--timing-reference=connect"}, zeit("absolute", "connect")},
		{"absolute, first-request", []string{"--keep-timing", "--timing-mode=absolute", "--timing-reference=first-request"}, zeit("absolute", "first-request")},
		{"Reihenfolge der Optionen", []string{"--timing-reference=first-request", "--timing-mode", "absolute", "--keep-timing"}, zeit("absolute", "first-request")},
	} {
		leere(t, "play")
		got, err := playMit(append([]string{"--upstream=pg:1"}, f.args...)...)
		if err != nil || got.Einspielen != f.want {
			t.Errorf("%s: %#v, %v, erwartet %#v", f.name, got.Einspielen, err, f.want)
		}
	}
}

// Abdeckung: LH-FA-17/Boundary, LH-FA-21/Happy — jede der drei Optionen gilt aus
// Umgebungsvariable und Schlüssel wie aus der Option; die Option geht der Umgebungsvariable
// vor, diese dem Schlüssel; eine leere Umgebungsvariable gilt als nicht gesetzt
// (LH-FA-17.a, SPEC-007).
func TestPlayZeitQuellen(t *testing.T) {
	for _, f := range []struct {
		option, env, schluessel string
		basis                   []string
		a, b                    string
		wantA, wantB            cli.Einspielvorgaben
	}{
		{"keep-timing", "PGWIRE_RECORDER_KEEP_TIMING", "keep_timing", nil, "true", "false", zeit("relative", ""), cli.Einspielvorgaben{}},
		{"timing-mode", "PGWIRE_RECORDER_TIMING_MODE", "timing_mode", []string{"--keep-timing"}, "absolute", "relative", zeit("absolute", "connect"), zeit("relative", "")},
		{"timing-reference", "PGWIRE_RECORDER_TIMING_REFERENCE", "timing_reference", []string{"--keep-timing", "--timing-mode=absolute"}, "first-request", "connect", zeit("absolute", "first-request"), zeit("absolute", "connect")},
	} {
		t.Run(f.option, func(t *testing.T) {
			leere(t, "play")
			lies := func(extra ...string) cli.Einspielvorgaben {
				t.Helper()
				got, err := playMit(append(append([]string{"--upstream=pg:1"}, f.basis...), extra...)...)
				if err != nil {
					t.Fatalf("%q: %v", extra, err)
				}
				return got.Einspielen
			}
			datei := func(wert string) string {
				return "--config=" + schreibe(t, "play:\n  "+f.schluessel+": "+wert+"\n")
			}
			if got := lies("--" + f.option + "=" + f.a); got != f.wantA {
				t.Errorf("Option %s: %#v, erwartet %#v", f.a, got, f.wantA)
			}
			if got := lies(datei(f.a)); got != f.wantA {
				t.Errorf("Schlüssel %s: %#v, erwartet %#v", f.a, got, f.wantA)
			}
			t.Setenv(f.env, f.a)
			if got := lies(); got != f.wantA {
				t.Errorf("Umgebungsvariable %s: %#v, erwartet %#v", f.a, got, f.wantA)
			}
			if got := lies(datei(f.b)); got != f.wantA {
				t.Errorf("Umgebungsvariable %s vor Schlüssel %s: %#v, erwartet %#v", f.a, f.b, got, f.wantA)
			}
			if got := lies("--" + f.option + "=" + f.b); got != f.wantB {
				t.Errorf("Option %s vor Umgebungsvariable %s: %#v, erwartet %#v", f.b, f.a, got, f.wantB)
			}
			t.Setenv(f.env, "")
			if got := lies(datei(f.b)); got != f.wantB {
				t.Errorf("leere Umgebungsvariable neben Schlüssel %s: %#v, erwartet %#v", f.b, got, f.wantB)
			}
		})
	}
}

// Abdeckung: LH-FA-21/Negative, LH-FA-17/Negative — ungültige Verwendung ist PGR-E2001:
// --timing-mode oder --timing-reference ohne --keep-timing, --timing-reference zusammen mit
// relative; „gesetzt“ ist eine Option, wenn irgendeine Quelle (Option, Umgebungsvariable,
// Schlüssel) sie nennt, auch mit dem Wert des Standards (--timing-reference connect mit
// relative, --timing-mode relative ohne --keep-timing, --keep-timing=false mit
// --timing-mode absolute); die Reihenfolge der Optionen ändert nichts; die Meldung nennt
// genau die Optionen, auch bei einem Wert aus Umgebungsvariable oder Datei keinen Wert
// (LH-FA-21.a *Randformen des Einspielens*, *Gesetzt*).
func TestPlayZeitVerwendung(t *testing.T) {
	const ohneKeep = "verlangt --keep-timing"
	const ohneAbsolute = "--timing-reference verlangt --timing-mode absolute"
	for _, f := range []struct {
		name    string
		args    []string
		env     map[string]string
		datei   string
		meldung string
	}{
		{"Modus absolute ohne --keep-timing", []string{"--timing-mode=absolute"}, nil, "", "--timing-mode " + ohneKeep},
		{"Modus relative ohne --keep-timing", []string{"--timing-mode=relative"}, nil, "", "--timing-mode " + ohneKeep},
		{"Bezug connect ohne --keep-timing", []string{"--timing-reference=connect"}, nil, "", "--timing-reference " + ohneKeep},
		{"Bezug first-request ohne --keep-timing", []string{"--timing-reference=first-request"}, nil, "", "--timing-reference " + ohneKeep},
		{"Modus absolute mit --keep-timing=false", []string{"--keep-timing=false", "--timing-mode=absolute"}, nil, "", "--timing-mode " + ohneKeep},
		{"Reihenfolge umgekehrt", []string{"--timing-mode=absolute", "--keep-timing=false"}, nil, "", "--timing-mode " + ohneKeep},
		{"Bezug mit Modus relative", []string{"--keep-timing", "--timing-mode=relative", "--timing-reference=first-request"}, nil, "", ohneAbsolute},
		{"Bezug connect mit Modus relative", []string{"--keep-timing", "--timing-mode=relative", "--timing-reference=connect"}, nil, "", ohneAbsolute},
		{"Bezug connect mit Standardmodus", []string{"--keep-timing", "--timing-reference=connect"}, nil, "", ohneAbsolute},
		{"Modus aus der Umgebung ohne --keep-timing", nil, map[string]string{"PGWIRE_RECORDER_TIMING_MODE": "relative"}, "", "--timing-mode " + ohneKeep},
		{"Bezug aus der Umgebung ohne --keep-timing", nil, map[string]string{"PGWIRE_RECORDER_TIMING_REFERENCE": "connect"}, "", "--timing-reference " + ohneKeep},
		{"Modus aus der Datei ohne --keep-timing", nil, nil, "play:\n  timing_mode: relative\n", "--timing-mode " + ohneKeep},
		{"Bezug aus der Datei ohne --keep-timing", nil, nil, "play:\n  timing_reference: connect\n", "--timing-reference " + ohneKeep},
		{"Bezug aus der Datei mit Modus relative aus der Umgebung", []string{"--keep-timing"}, map[string]string{"PGWIRE_RECORDER_TIMING_MODE": "relative"}, "play:\n  timing_reference: connect\n", ohneAbsolute},
		{"--keep-timing=false vor der Umgebung", []string{"--keep-timing=false"}, map[string]string{"PGWIRE_RECORDER_KEEP_TIMING": "true", "PGWIRE_RECORDER_TIMING_MODE": "absolute"}, "", "--timing-mode " + ohneKeep},
	} {
		t.Run(f.name, func(t *testing.T) {
			leere(t, "play")
			for k, v := range f.env {
				t.Setenv(k, v)
			}
			args := append([]string{"--upstream=pg:1"}, f.args...)
			if f.datei != "" {
				args = append(args, "--config="+schreibe(t, f.datei))
			}
			_, err := playMit(args...)
			if want := "Konfiguration [PGR-E2001]: " + f.meldung; !istUsage(err) || err.Error() != want {
				t.Fatalf("erwartet %s mit %q, erhalten %v", model.CodeUsage, want, err)
			}
		})
	}
}

// Abdeckung: LH-FA-21/Negative, LH-FA-17/Negative — ein Wert außerhalb der Aufzählungen,
// auch der leere, Großschreibung und Unterstrich (Absolute, Connect, first_request), ist
// auf der Kommandozeile und in der Umgebungsvariable PGR-E2001, im Schlüssel PGR-E2004;
// die Meldung der Umgebungsvariable und des Schlüssels nennt ihren Namen und nie den
// Wert (LH-FA-21.a *Gesetzt*, LH-FA-17.a).
func TestPlayZeitWerteUngueltig(t *testing.T) {
	for _, f := range []struct {
		option, env, schluessel string
		basis                   []string
		schlecht                []string
	}{
		{"timing-mode", "PGWIRE_RECORDER_TIMING_MODE", "timing_mode", []string{"--keep-timing"}, []string{"", "Absolute", "RELATIVE", "GEHEIM"}},
		{"timing-reference", "PGWIRE_RECORDER_TIMING_REFERENCE", "timing_reference", []string{"--keep-timing", "--timing-mode=absolute"}, []string{"", "Connect", "first_request", "First-Request", "GEHEIM"}},
	} {
		for _, w := range f.schlecht {
			leere(t, "play")
			args := append([]string{"--upstream=pg:1"}, f.basis...)
			if _, err := playMit(append(args, "--"+f.option+"="+w)...); !istUsage(err) || !strings.Contains(err.Error(), f.option) {
				t.Errorf("--%s=%q: erwartet %s mit der Option, erhalten %v", f.option, w, model.CodeUsage, err)
			}
			if w == "" {
				// Eine leere Umgebungsvariable gilt als nicht gesetzt; ein leerer
				// Schlüssel ist ungültig.
				if _, err := playMit(append(args, "--config="+schreibe(t, "play:\n  "+f.schluessel+": \"\"\n"))...); !istDatei(err) || !strings.Contains(err.Error(), "play."+f.schluessel) {
					t.Errorf("leerer Schlüssel %s: %v", f.schluessel, err)
				}
				continue
			}
			t.Setenv(f.env, w)
			_, err := playMit(args...)
			if !istUsage(err) || !strings.Contains(err.Error(), f.env) || strings.Contains(err.Error(), w) {
				t.Errorf("%s=%q: erwartet %s mit dem Namen der Variable ohne den Wert, erhalten %v", f.env, w, model.CodeUsage, err)
			}
			t.Setenv(f.env, "")
			_, err = playMit(append(args, "--config="+schreibe(t, "play:\n  "+f.schluessel+": "+w+"\n"))...)
			if !istDatei(err) || !strings.Contains(err.Error(), "play."+f.schluessel) || strings.Contains(err.Error(), w) {
				t.Errorf("Schlüssel %s: %q: erwartet %s mit dem Schlüssel ohne den Wert, erhalten %v", f.schluessel, w, model.CodeConfigFile, err)
			}
		}
	}
}

// Abdeckung: LH-FA-01/Happy, LH-FA-21/Happy — die Hilfe von play nennt --keep-timing,
// --timing-mode und --timing-reference mit ihren Werten, Standards und Umgebungsvariablen
// und dass die beiden letzten --keep-timing verlangen; die Hilfe von record und replay
// nennt sie nicht (LH-FA-21.a *Einspielen*, LH-FA-01.a).
func TestPlayZeitHilfe(t *testing.T) {
	text := strings.Join(strings.Fields(hilfe(t, "play", "--help")), " ")
	for _, teil := range []string{
		"--keep-timing[=true|false]", "PGWIRE_RECORDER_KEEP_TIMING", "fehlt sie einer Interaktion, endet play mit PGR-E2003",
		"--timing-mode relative|absolute", "PGWIRE_RECORDER_TIMING_MODE", "Standard relative",
		"--timing-reference connect|first-request", "PGWIRE_RECORDER_TIMING_REFERENCE", "Standard connect",
		"nur mit --keep-timing (sonst PGR-E2001)", "nur mit --keep-timing und absolute (sonst PGR-E2001)",
	} {
		if !strings.Contains(text, teil) {
			t.Errorf("Hilfe von play ohne %q:\n%s", teil, text)
		}
	}
	for _, kommando := range []string{"record", "replay"} {
		text := hilfe(t, kommando, "--help")
		for _, teil := range []string{"--keep-timing", "--timing-mode", "--timing-reference", "PGWIRE_RECORDER_KEEP_TIMING", "PGWIRE_RECORDER_TIMING_MODE", "PGWIRE_RECORDER_TIMING_REFERENCE"} {
			if strings.Contains(text, teil) {
				t.Errorf("Hilfe von %s nennt %q", kommando, teil)
			}
		}
	}
}
