//go:build integration

package integration_test

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// absteigendSQLite ist eine Aufzeichnung im Format sqlite mit einer Session und zwei
// Interaktionen, deren offset_ms 4711 und 4709 sind. Herkunft: mit `record --format
// sqlite` zwei Anfragen in einer Verbindung aufgenommen, dann offset_ms der beiden
// Interaktionen auf 4711 und 4709 gesetzt. Eine Änderung des Schemas, die die Datei
// ungültig macht, lässt den Test rot werden, weil "Session 1, Interaktion 2" dann
// in der Meldung fehlt.
//
//go:embed testdata/zeitangaben-absteigend.sqlite
var absteigendSQLite []byte

// Abdeckung: LH-FA-21/Negative — replay und play beenden den Start mit Exit-Code 3 und PGR-E3003,
// wenn offset_ms innerhalb einer Session absteigt, in einer Aufzeichnung aus yaml wie aus sqlite,
// und wenn offset_ms in yaml keine kanonische Dezimalzahl ist; die Meldung nennt Session und
// Interaktion und keinen Wert, und es läuft nichts an (play verbindet sich nicht).
func TestE2ENichtAbsteigendeZeitangaben(t *testing.T) {
	dir := t.TempDir()
	sqlite := filepath.Join(dir, "rec.sqlite")
	if err := os.WriteFile(sqlite, absteigendSQLite, 0o644); err != nil {
		t.Fatal(err)
	}
	yamlText := func(offset string) string {
		return "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n" +
			"      - sequence: 1\n        offset_ms: 4711\n        request: {type: query, sql: x}\n        responses:\n          - {type: ready_for_query, tx_status: I}\n" +
			"      - sequence: 2\n        offset_ms: " + offset + "\n        request: {type: query, sql: x}\n        responses:\n          - {type: ready_for_query, tx_status: I}\n"
	}
	faelle := []struct{ name, datei, inhalt, wert string }{
		{"sqlite absteigend", sqlite, "", "4709"},
		{"yaml absteigend", "abst.yaml", yamlText("4709"), "4709"},
		{"yaml 4711.5", "punkt.yaml", yamlText("4711.5"), "4711.5"},
		{"yaml Text", "text.yaml", yamlText(`"4712"`), "4712"},
	}
	for _, c := range faelle {
		input := c.datei
		if c.inhalt != "" {
			input = filepath.Join(dir, c.datei)
			if err := os.WriteFile(input, []byte(c.inhalt), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		for _, args := range [][]string{
			{"replay", "--listen", freieAdresse(t), "--input", input},
			{"play", "--upstream", freieAdresse(t), "--input", input},
		} {
			stdout, stderr, code := starteBis(t, "", args...)
			if code != 3 || !strings.Contains(stderr, "[PGR-E3003]") || stdout != "" ||
				!strings.Contains(stderr, "Session 1, Interaktion 2") || strings.Contains(stderr, c.wert) {
				t.Errorf("%s, %s: Exit-Code %d, stdout %q, stderr %q", c.name, args[0], code, stdout, stderr)
			}
		}
	}
}
