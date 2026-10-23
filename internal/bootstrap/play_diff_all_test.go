package bootstrap_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Negative, LH-FA-20/Boundary — --diff-all
// (Kommandozeile oder Umgebungsvariable) erreicht den Play-Service: Bei einer
// Interaktion mit zwei Abweichungen meldet play eine Zeile PGR-E5004 mit beiden und
// Exit-Code 5, ohne die Option nur die erste; die Zeile zum Ende zählt die
// Interaktion einmal; ohne --compare-responses ist die Option PGR-E2001 mit
// Exit-Code 2 und ohne Verbindung (LH-FA-24.a *Alle Abweichungen*).
func TestRunPlayDiffAll(t *testing.T) {
	leerePlay(t)
	anders := append(nachricht('C', []byte("UPDATE 1\x00")...), nachricht('Z', 'T')...)
	alle := "Replay [PGR-E5004]: Session 1, Interaktion 1: Antworten weichen von der Aufzeichnung ab, Abweichungen: 2; " +
		"Position 1, Art Befehl (Befehl aufgezeichnet „SELECT“, Server „UPDATE“); Position 2, Art Transaktionsstatus (aufgezeichnet „I“, Server „T“)"
	erste := "Replay [PGR-E5004]: Session 1, Interaktion 1: Antwort weicht von der Aufzeichnung ab: Position 1, Art Befehl"
	for _, f := range []struct {
		name string
		args []string
		env  string
		want string
	}{
		{"Option", []string{"--compare-responses", "--diff-all"}, "", alle},
		{"Umgebungsvariable", []string{"--compare-responses"}, "true", alle},
		{"ohne Option", []string{"--compare-responses"}, "", erste},
		{"false", []string{"--compare-responses", "--diff-all=false"}, "", erste},
	} {
		t.Run(f.name, func(t *testing.T) {
			t.Setenv("PGWIRE_RECORDER_DIFF_ALL", f.env)
			addr, _ := playServer(t, anders)
			input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A")})
			var stdout, stderr bytes.Buffer
			code := bootstrap.Run(context.Background(), nil, append([]string{"play", "--upstream", addr, "--input", input}, f.args...), "dev", &stdout, &stderr)
			zeilen := beendetZeilen(stderr.String())
			if code != 5 || !strings.Contains(stderr.String(), `code=PGR-E5004 error="`+f.want+`"`) || len(zeilen) != 1 || !strings.HasSuffix(zeilen[0], "eingespielt=1 verglichen=1 abweichend=1") {
				t.Fatalf("Exit-Code %d, Zeilen zum Ende %q, erwartet 5 mit %q:\n%s", code, zeilen, f.want, stderr.String())
			}
		})
	}
	addr, typen := playServer(t, anders)
	var stdout, stderr bytes.Buffer
	code := bootstrap.Run(context.Background(), nil, []string{"play", "--upstream", addr, "--input", einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A")}), "--diff-all"}, "dev", &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "PGR-E2001") || !strings.Contains(stderr.String(), "--diff-all") {
		t.Fatalf("--diff-all ohne --compare-responses: Exit-Code %d, stderr %q", code, stderr.String())
	}
	select {
	case b := <-typen:
		t.Fatalf("play verbindet sich trotz ungültiger Kombination: %q", b)
	default:
	}
}
