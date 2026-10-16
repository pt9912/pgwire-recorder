//go:build integration

package integration_test

import (
	_ "embed"
	"os"
	"strings"
	"testing"
	"time"
)

// Die Aufzeichnungen sind von Hand geschrieben (Version 2); jede gehört genau
// einem Test.
//
//go:embed testdata/play-ende-grenze-ohne-option.yaml
var playEndeGrenzeOhneOptionYAML []byte

//go:embed testdata/play-ende-grenze-ohne-ende.yaml
var playEndeGrenzeOhneEndeYAML []byte

// spieleEndeNichtErwartet spielt die Aufzeichnung inhalt mit den weiteren args
// gegen die leere Datenbank db ein, deren zweite Interaktion den Backend beendet
// (FATAL 57P01), und prüft: Exit-Code 4, die Meldung trägt den Code PGR-E4003 und
// die Stelle der Interaktion (Session 1, Interaktion 2).
func spieleEndeNichtErwartet(t *testing.T, db string, inhalt []byte, args ...string) {
	t.Helper()
	leereDatenbank(t, db)
	_, stderr, code := starteBisFrist(t, 30*time.Second, append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", endeDatei(t, inhalt), "--database", db}, args...)...)
	if code != 4 || !strings.Contains(stderr, "code=PGR-E4003") || !strings.Contains(stderr, "Session 1, Interaktion 2") {
		t.Fatalf("Exit-Code %d, erwartet 4 mit code=PGR-E4003 und Session 1, Interaktion 2, stderr:\n%s", code, stderr)
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-20/Negative — play spielt gegen
// PostgreSQL (Server 17) eine handgeschriebene Aufzeichnung ein, deren zweite
// Anfrage den Backend beendet (Fehlerantwort FATAL 57P01) und die das Ende an
// dieser Stelle trägt; ohne --allow-recorded-errors ist das Ende nicht
// erwartet: PGR-E4003 mit der Stelle der Interaktion, Exit-Code 4 (LH-FA-20.a
// *Interaktion*, LH-FA-25.a).
func TestE2EPlayEndeGrenzeOhneOption(t *testing.T) {
	spieleEndeNichtErwartet(t, "play_ende_grenze_ohne_option", playEndeGrenzeOhneOptionYAML)
}

// Abdeckung: LH-FA-25/Negative, LH-FA-20/Negative — wie
// TestE2EPlayEndeGrenzeOhneOption, mit --allow-recorded-errors, aber die
// Aufzeichnung trägt an der Stelle ein ERROR und kein Ende, der Server sendet
// FATAL 57P01: Das Ende ist nicht erwartet, PGR-E4003 mit der Stelle der
// Interaktion, Exit-Code 4 (LH-FA-20.a *Interaktion*, LH-FA-25.a).
func TestE2EPlayEndeGrenzeOhneEnde(t *testing.T) {
	spieleEndeNichtErwartet(t, "play_ende_grenze_ohne_ende", playEndeGrenzeOhneEndeYAML, "--allow-recorded-errors")
}
