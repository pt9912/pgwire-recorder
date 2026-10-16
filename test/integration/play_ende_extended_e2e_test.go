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
//go:embed testdata/play-ende-extended-sync.yaml
var playEndeExtendedSyncYAML []byte

//go:embed testdata/play-ende-extended-flush.yaml
var playEndeExtendedFlushYAML []byte

// spieleEndeExtended spielt die Aufzeichnung inhalt mit --allow-recorded-errors
// gegen die leere Datenbank db ein und prüft: Exit-Code 0 ohne Meldung, die Zeile
// zum Ende zählt 3 eingespielt, 0 verglichen, 0 abweichend, und die Session nach
// dem Ende lief (ihre Wirkung steht in der Datenbank).
func spieleEndeExtended(t *testing.T, db string, inhalt []byte) {
	t.Helper()
	conn := leereDatenbank(t, db)
	_, stderr, code := starteBisFrist(t, 30*time.Second, "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", endeDatei(t, inhalt), "--database", db, "--allow-recorded-errors")
	if code != 0 || strings.Contains(stderr, "level=ERROR") || strings.Contains(stderr, "level=WARN") || zaehler(t, stderr) != "3/0/0" {
		t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
	}
	if got := wert(t, conn, "SELECT to_regclass('ende_folge') IS NOT NULL"); got != "t" {
		t.Fatalf("die Session nach dem Ende lief nicht, to_regclass nicht gesetzt: %s", got)
	}
}

// Abdeckung: LH-FA-25/Happy, LH-FA-20/Boundary, LH-FA-18/Boundary — play spielt
// gegen PostgreSQL (Server 17) eine handgeschriebene Extended-Aufzeichnung mit
// Ende ein, deren Ausführung pg_terminate_backend(pg_backend_pid()) in der Gruppe
// mit Sync den Backend beendet (Fehlerantwort FATAL 57P01): Mit
// --allow-recorded-errors endet play mit Exit-Code 0 ohne Meldung, die nächste
// Session läuft, und die Zeile zum Ende zählt 3 eingespielt, 0 verglichen, 0
// abweichend (LH-FA-20.a *Interaktion*, *Gruppen*, LH-FA-25.a).
func TestE2EPlayEndeExtendedSync(t *testing.T) {
	spieleEndeExtended(t, "play_ende_extended_sync", playEndeExtendedSyncYAML)
}

// Abdeckung: LH-FA-25/Happy, LH-FA-20/Boundary, LH-FA-18/Boundary — wie
// TestE2EPlayEndeExtendedSync, die Gruppe mit dem Ende schließt mit Flush: Die
// Fehlerantwort FATAL 57P01 ist die Antwort auf Execute, play endet mit Exit-Code
// 0 ohne Meldung, die nächste Session läuft (LH-FA-20.a *Gruppen*, LH-FA-25.a).
func TestE2EPlayEndeExtendedFlush(t *testing.T) {
	spieleEndeExtended(t, "play_ende_extended_flush", playEndeExtendedFlushYAML)
}
