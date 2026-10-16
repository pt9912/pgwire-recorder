//go:build integration

package integration_test

import (
	_ "embed"
	"os"
	"strings"
	"testing"
	"time"
)

// Die Aufzeichnung ist von Hand geschrieben (Version 2); sie gehört genau
// diesem Test.
//
//go:embed testdata/play-ende-erwartet.yaml
var playEndeErwartetYAML []byte

// Abdeckung: LH-FA-25/Happy, LH-FA-20/Boundary — play spielt gegen PostgreSQL
// eine handgeschriebene Aufzeichnung mit Ende ein, deren Anfrage
// pg_terminate_backend(pg_backend_pid()) den Backend beendet (Fehlerantwort
// FATAL 57P01, Server 17): Mit --allow-recorded-errors endet play mit Exit-Code
// 0 ohne Meldung, die nächste Session läuft (ihre Wirkung steht in der
// Datenbank), und die Zeile zum Ende zählt 3 eingespielt, 0 verglichen, 0
// abweichend (LH-FA-20.a *Interaktion*, LH-FA-25.a, *Ende einer
// Session*).
func TestE2EPlayEndeErwartet(t *testing.T) {
	conn := leereDatenbank(t, "play_ende_erwartet")
	input := endeDatei(t, playEndeErwartetYAML)
	_, stderr, code := starteBisFrist(t, 30*time.Second, "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", "play_ende_erwartet", "--allow-recorded-errors")
	if code != 0 || strings.Contains(stderr, "level=ERROR") || strings.Contains(stderr, "level=WARN") || zaehler(t, stderr) != "3/0/0" {
		t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
	}
	if got := wert(t, conn, "SELECT to_regclass('ende_folge') IS NOT NULL"); got != "t" {
		t.Fatalf("die Session nach dem Ende lief nicht, to_regclass nicht gesetzt: %s", got)
	}
}
