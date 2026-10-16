package cli_test

import (
	"strings"
	"testing"
)

// optionsblock liefert den Block der Option name aus der Hilfe text, von ihrer
// Zeile bis zur nächsten Option, den Leerraum zu einem Leerzeichen zusammengezogen;
// fehlt die Option, wird der Test rot.
func optionsblock(t *testing.T, text, name string) string {
	t.Helper()
	var block []string
	drin := false
	for _, zeile := range strings.Split(text, "\n") {
		if strings.HasPrefix(zeile, "  --") {
			if drin {
				break
			}
			drin = strings.HasPrefix(zeile, "  "+name)
		}
		if drin {
			block = append(block, zeile)
		}
	}
	if len(block) == 0 {
		t.Fatalf("die Hilfe nennt %s nicht:\n%s", name, text)
	}
	return strings.Join(strings.Fields(strings.Join(block, " ")), " ")
}

// Abdeckung: LH-FA-24/Happy, LH-FA-20/Happy — die Hilfe von play sagt zu
// --allow-recorded-errors, dass eine Fehlerantwort ERROR erlaubt ist, wenn die
// Aufzeichnung sie trägt, ebenso ein Ende des Servers (FATAL oder PANIC), dann
// die nächste Session läuft und die Option mit --compare-responses ohne Wirkung
// ist; zu --compare-responses, dass ein Ende des Servers (FATAL oder PANIC)
// erwartet ist, wenn die Aufzeichnung es an derselben Stelle mit demselben
// SQLSTATE zeigt, sonst eine Abweichung; die Aussagen stehen im Block ihrer
// Option (LH-FA-20.a, LH-FA-24.a).
func TestPlayEndeHilfe(t *testing.T) {
	text := hilfe(t, "play", "--help")
	for _, f := range []struct {
		option, aussage string
	}{
		{"--allow-recorded-errors", "eine Fehlerantwort mit dem Schweregrad ERROR ist kein Fehler, wenn die aufgezeichnete Anfrage ebenfalls eine trägt"},
		{"--allow-recorded-errors", "ebenso ein Ende des Servers (FATAL oder PANIC), wenn die aufgezeichnete Anfrage es ebenfalls trägt"},
		{"--allow-recorded-errors", "dann läuft die nächste Session"},
		{"--allow-recorded-errors", "mit --compare-responses ohne Wirkung"},
		{"--compare-responses", "ein Ende des Servers (FATAL oder PANIC) ist erwartet, wenn die Aufzeichnung es"},
		{"--compare-responses", "an derselben Stelle"},
		{"--compare-responses", "mit demselben SQLSTATE zeigt"},
		{"--compare-responses", "sonst eine Abweichung"},
	} {
		t.Run(f.option+": "+f.aussage, func(t *testing.T) {
			if block := optionsblock(t, text, f.option); !strings.Contains(block, f.aussage) {
				t.Errorf("der Block von %s nennt %q nicht:\n%s", f.option, f.aussage, block)
			}
		})
	}
}
