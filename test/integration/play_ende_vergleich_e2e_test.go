//go:build integration

package integration_test

import (
	_ "embed"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// Die Aufzeichnungen sind von Hand geschrieben (Version 2), je eine pro Fall; sie
// unterscheiden sich im SQLSTATE des aufgezeichneten Endes.
//
//go:embed testdata/play-ende-vergleich-gleich.yaml
var playEndeVergleichGleichYAML []byte

//go:embed testdata/play-ende-vergleich-anders.yaml
var playEndeVergleichAndersYAML []byte

// endeVergleichLauf spielt inhalt gegen eine leere Datenbank db mit
// --compare-responses und den weiteren args ein und liefert die Verbindung zur
// Datenbank, stderr und den Exit-Code.
func endeVergleichLauf(t *testing.T, db string, inhalt []byte, args ...string) (*pgconn.PgConn, string, int) {
	t.Helper()
	conn := leereDatenbank(t, db)
	input := endeDatei(t, inhalt)
	_, stderr, code := starteBisFrist(t, 30*time.Second, append([]string{"play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input, "--database", db, "--compare-responses"}, args...)...)
	return conn, stderr, code
}

// Abdeckung: LH-FA-24/Boundary, LH-FA-25/Happy — play spielt gegen PostgreSQL
// eine handgeschriebene Aufzeichnung mit Ende (Fehlerantwort FATAL 57P01 auf
// einen DO-Block mit pg_terminate_backend(pg_backend_pid()) und pg_sleep, Server
// 17) mit --compare-responses ohne --allow-recorded-errors ein: Der Server führt
// das Ende an derselben Stelle mit demselben SQLSTATE herbei, das Ende ist
// erwartet, play endet mit Exit-Code 0 ohne Meldung, die nächste Session läuft
// (ihre Wirkung steht in der Datenbank), und die Zeile zum Ende zählt 3
// eingespielt, 3 verglichen, 0 abweichend (Abnahmeszenario 16, dritter Satz,
// LH-FA-24.a *Ende der Aufzeichnung*); dass der Recorder ein solches Ende genauso
// aufzeichnet, belegt dieser Test nicht.
func TestE2EPlayEndeVergleichGleich(t *testing.T) {
	conn, stderr, code := endeVergleichLauf(t, "play_ende_vergleich_gleich", playEndeVergleichGleichYAML)
	if code != 0 || strings.Contains(stderr, "level=ERROR") || strings.Contains(stderr, "level=WARN") || strings.Contains(stderr, "PGR-E4003") || zaehler(t, stderr) != "3/3/0" {
		t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
	}
	if got := wert(t, conn, "SELECT to_regclass('ende_folge') IS NOT NULL"); got != "t" {
		t.Fatalf("die Session nach dem Ende lief nicht, to_regclass nicht gesetzt: %s", got)
	}
}

// Abdeckung: LH-FA-24/Negative — dieselbe handgeschriebene Aufzeichnung, doch
// das aufgezeichnete Ende trägt 57P02 statt 57P01: Mit --compare-responses, ohne
// und mit --allow-recorded-errors und beide Male ohne --continue-on-error, endet
// play mit PGR-E5004 der Art Fehler, die beide SQLSTATE nennt, genau einmal, ohne
// PGR-E4003 und mit Exit-Code 5; die Zeile zum Ende zählt 2 eingespielt, 2
// verglichen, 1 abweichend, und die nächste Session läuft nicht (Abnahmeszenario
// 16, dritter Satz, LH-FA-24.a *Ende der Aufzeichnung*, LH-FA-20.a Tabelle).
func TestE2EPlayEndeVergleichAnders(t *testing.T) {
	meldung := "Session 1, Interaktion 2: Antwort weicht von der Aufzeichnung ab: Position 1, Art Fehler (SQLSTATE aufgezeichnet „57P02“, Server „57P01“)"
	for _, f := range []struct {
		db   string
		args []string
	}{{"play_ende_vergleich_anders", nil}, {"play_ende_vergleich_anders_allow", []string{"--allow-recorded-errors"}}} {
		t.Run(f.db, func(t *testing.T) {
			conn, stderr, code := endeVergleichLauf(t, f.db, playEndeVergleichAndersYAML, f.args...)
			if code != 5 || strings.Count(stderr, "code=PGR-E5004") != 1 || strings.Contains(stderr, "PGR-E4003") || !strings.Contains(stderr, meldung) || zaehler(t, stderr) != "2/2/1" {
				t.Fatalf("%q: Exit-Code %d, Zähler %s, stderr:\n%s", f.args, code, zaehler(t, stderr), stderr)
			}
			if got := wert(t, conn, "SELECT to_regclass('ende_folge') IS NULL"); got != "t" {
				t.Fatalf("die Session nach der Abweichung lief, to_regclass gesetzt: %s", got)
			}
		})
	}
}
