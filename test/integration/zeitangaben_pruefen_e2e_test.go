//go:build integration

package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-21/Negative — replay und play beenden den Start mit Exit-Code 3 und PGR-E3003,
// wenn offset_ms innerhalb einer Session absteigt, in einer Aufzeichnung aus yaml wie aus sqlite,
// und wenn offset_ms in yaml keine kanonische Dezimalzahl ist; die Meldung nennt Session und
// Interaktion und keinen Wert, und es läuft nichts an (play verbindet sich nicht).
func TestE2ENichtAbsteigendeZeitangaben(t *testing.T) {
	absteigend := func(off ...int64) model.Recording {
		rec := model.NewRecording()
		s := model.Session{ID: 1}
		for i := range off {
			s.Interactions = append(s.Interactions, model.Interaction{
				Sequence:  i + 1,
				OffsetMS:  &off[i],
				Request:   model.Request{Type: model.RequestQuery, SQL: "SELECT 1"},
				Responses: []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}},
			})
		}
		rec.Sessions = append(rec.Sessions, s)
		return rec
	}
	dir := t.TempDir()
	sqlite := filepath.Join(dir, "rec.sqlite")
	if err := (&recording.SQLite{}).Write(context.Background(), sqlite, absteigend(4711, 4709)); err != nil {
		t.Fatal(err)
	}
	yamlText := func(offset string) string {
		return "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n" +
			"      - sequence: 1\n        offset_ms: 4711\n        request: {type: query, sql: x}\n        responses:\n          - {type: ready_for_query, tx_status: I}\n" +
			"      - sequence: 2\n        offset_ms: " + offset + "\n        request: {type: query, sql: x}\n        responses:\n          - {type: ready_for_query, tx_status: I}\n"
	}
	for name, c := range map[string]struct {
		datei  string
		inhalt []byte
		wert   string
	}{
		"sqlite absteigend": {"rec.sqlite", nil, "4709"},
		"yaml absteigend":   {"abst.yaml", []byte(yamlText("4709")), "4709"},
		"yaml 1.5":          {"punkt.yaml", []byte(yamlText("4711.5")), "4711.5"},
		"yaml Text":         {"text.yaml", []byte(yamlText(`"4712"`)), "4712"},
	} {
		input := filepath.Join(dir, c.datei)
		if c.inhalt != nil {
			if err := os.WriteFile(input, c.inhalt, 0o644); err != nil {
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
				t.Errorf("%s, %s: Exit-Code %d, stdout %q, stderr %q", name, args[0], code, stdout, stderr)
			}
		}
	}
}
