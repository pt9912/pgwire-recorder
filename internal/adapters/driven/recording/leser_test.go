package recording_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// yamlDatei schreibt rec als YAML nach path.
func yamlDatei(t *testing.T, path string, rec model.Recording) {
	t.Helper()
	if err := (recording.YAML{}).Write(context.Background(), path, rec); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-22/Happy, LH-FA-22/Boundary — der Leser erkennt das Format an den ersten 16 Bytes, nicht
// an der Endung: eine SQLite-Datei mit der Endung .yaml und eine YAML-Datei mit der Endung
// .sqlite lesen als das, was sie sind.
func TestLeserErkennungOhneEndung(t *testing.T) {
	rec := lesefixture()
	dir := t.TempDir()
	sqliteAlsYAML := filepath.Join(dir, "a.yaml")
	if err := os.WriteFile(sqliteAlsYAML, mustRead(t, sqliteDatei(t, rec)), 0o644); err != nil {
		t.Fatal(err)
	}
	yamlAlsSQLite := filepath.Join(dir, "b.sqlite")
	yamlDatei(t, yamlAlsSQLite, rec)
	for _, path := range []string{sqliteAlsYAML, yamlAlsSQLite} {
		if got := ladeOK(t, path); !reflect.DeepEqual(got, rec) {
			t.Errorf("%s: %#v", path, got)
		}
	}
}

// Abdeckung: LH-FA-22/Negative, LH-FA-22/Boundary — sqlite ist nur eine Datei, die mit SQLite format 3 und einem
// NUL-Byte beginnt; jede andere Datei liest YAML, auch eine leere, eine mit weniger als 16 Bytes
// und eine, die den Kopf bis auf ein Byte trägt (das Lesen als YAML meldet eine beschädigte
// Aufzeichnung, das Lesen als sqlite eine Datei, die sich nicht lesen lässt).
func TestLeserErkennungAnDenErsten16Bytes(t *testing.T) {
	for name, c := range map[string]struct {
		inhalt string
		yaml   bool
	}{
		"leer":                  {"", true},
		"ein Byte":              {"S", true},
		"15 Bytes":              {"SQLite format 3", true},
		"16 Bytes, letztes X":   {"SQLite format 3X", true},
		"16 Bytes, erstes s":    {"sQLite format 3\x00", true},
		"16 Bytes, andere Zahl": {"SQLite format 2\x00", true},
		"16 Bytes, Leerzeichen": {"SQLite format 3 ", true},
		"16 Bytes mit NUL":      {"SQLite format 3\x00", false},
		"mehr als 16 Bytes":     {"SQLite format 3\x00 und mehr", false},
		"YAML":                  {"format: other\n", true},
		"YAML mit Kopf später":  {"# SQLite format 3\x00", true},
	} {
		path := filepath.Join(t.TempDir(), "x")
		if err := os.WriteFile(path, []byte(c.inhalt), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := laden(path)
		if codeVon(err) != beschaedigt {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if gelesenAlsYAML := strings.Contains(err.Error(), "Aufzeichnung beschädigt") || strings.Contains(err.Error(), "Formatkennung"); gelesenAlsYAML != c.yaml {
			t.Errorf("%s: als YAML gelesen %v, erwartet %v (%v)", name, gelesenAlsYAML, c.yaml, err)
		}
	}
}

// Abdeckung: LH-FA-22/Negative — scheitert das Öffnen der Datei, ist das PGR-E3001 mit dem Pfad und
// der Ursache: eine fehlende Datei, ein Verzeichnis; der Leser legt nichts an.
func TestLeserNichtLesbar(t *testing.T) {
	dir := t.TempDir()
	fehlt := filepath.Join(dir, "fehlt")
	for name, path := range map[string]string{"fehlt": fehlt, "Verzeichnis": dir, "Verzeichnis mit Schrägstrich": dir + "/"} {
		_, err := laden(path)
		if codeVon(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), path) {
			t.Errorf("%s: %v", name, err)
		}
	}
	if es := eintraege(t, dir); len(es) != 0 {
		t.Errorf("im Verzeichnis liegt %v", es)
	}
}

// Abdeckung: LH-FA-22/Boundary — der Pfad der Eingabedatei gilt wörtlich, für die Erkennung wie für
// beide Leser: Namen mit ?, #, %, %41, Leerzeichen und ?mode=rw lesen als Datei dieses
// Namens, absolut wie relativ, in einem Verzeichnis mit solchen Zeichen im Namen; der Pfad
// einer vorhandenen Datei mit angehängtem ?mode=rw nennt eine Datei, die es nicht gibt
// (PGR-E3001).
func TestLesenSonderzeichenImPfad(t *testing.T) {
	rec := lesefixture()
	for _, name := range []string{"a?b#c%d e.sqlite", "%41 %3F.sqlite", "x?mode=rw.sqlite", "#.sqlite", "+&=.sqlite", "a.yaml?x", "100%.yaml"} {
		dir := filepath.Join(t.TempDir(), "ver?zeich#nis %25")
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, relativ := range []bool{false, true} {
			path := filepath.Join(dir, name)
			if relativ {
				t.Chdir(dir)
				path = name
			}
			if strings.Contains(name, "sqlite") {
				if err := (&recording.SQLite{}).Write(context.Background(), path, rec); err != nil {
					t.Fatalf("%q: %v", name, err)
				}
			} else {
				yamlDatei(t, path, rec)
			}
			if got, err := laden(path); err != nil || !reflect.DeepEqual(got, rec) {
				t.Errorf("%q relativ=%v: %v", name, relativ, err)
			}
			if strings.Contains(name, "sqlite") {
				if got, err := (&recording.SQLite{}).Load(context.Background(), path); err != nil || !reflect.DeepEqual(got, rec) {
					t.Errorf("%q relativ=%v, SQLite.Load: %v", name, relativ, err)
				}
			}
			if err := os.Remove(filepath.Join(dir, name)); err != nil {
				t.Fatal(err)
			}
		}
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "x.sqlite")
	rohDatei := sqliteDatei(t, rec)
	if err := os.WriteFile(path, mustRead(t, rohDatei), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := laden(path + "?mode=rw"); codeVon(err) != model.CodeRecordingIO {
		t.Errorf("Pfad mit Parametern: %v", err)
	}
}

// Abdeckung: LH-FA-22/Negative — Prepare und Write des Lesers sind ein Programmierfehler des Aufrufers:
// PGR-E1000 ohne den Pfad, und es entsteht keine Datei.
func TestLeserPrepareUndWrite(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "geheim-pfad.yaml")
	for name, err := range map[string]error{
		"Prepare": recording.Leser{}.Prepare(ctx, path, true),
		"Write":   recording.Leser{}.Write(ctx, path, lesefixture()),
	} {
		if codeVon(err) != model.CodeInternal || strings.Contains(err.Error(), "geheim-pfad") {
			t.Errorf("%s: %v", name, err)
		}
	}
	if es := eintraege(t, filepath.Dir(path)); len(es) != 0 {
		t.Errorf("im Verzeichnis liegt %v", es)
	}
}

// Abdeckung: LH-FA-22/Negative, LH-QA-06/Messung — dieselbe Aufzeichnung aus yaml und aus sqlite bekommt dieselbe
// Ablehnung: ein unbekannter Schlüssel, null an einem Schlüssel, der einen Wert verlangt, ein
// doppelter Schlüssel und eine unbekannte Version sind in beiden PGR-E3003 beziehungsweise
// PGR-E3002; die Aufzeichnung ohne den Fehler liest in beiden gleich.
func TestLesenGleicheAblehnungBeiderFormate(t *testing.T) {
	yamlText := func(antwort, parse string) string {
		return fmt.Sprintf(`format: pgwire-recorder
version: 1
sessions:
  - id: 1
    interactions:
      - sequence: 1
        request:
          type: query
          sql: SELECT 1
        responses:
          - type: command_complete
            tag: SELECT 1%s
          - type: ready_for_query
            tx_status: I
      - sequence: 2
        type: extended
        groups:
          - client:
              - type: parse
                statement: %s
                sql: SELECT 1
                param_types: []
              - type: sync
            server:
              - type: parse_complete
              - type: ready_for_query
                tx_status: I
`, antwort, parse)
	}
	rec := model.NewRecording()
	rec.Sessions = []model.Session{{ID: 1, Interactions: []model.Interaction{
		{Sequence: 1, Request: model.Request{Type: model.RequestQuery, SQL: "SELECT 1"}, Responses: []model.Response{{Type: model.ResponseCommandComplete, Tag: "SELECT 1"}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}},
		{Sequence: 2, Request: model.Request{Type: model.RequestExtended}, Groups: []model.Group{{
			Client: []model.ClientMessage{{Type: model.ClientParse, Statement: "s", SQL: "SELECT 1"}, {Type: model.ClientSync}},
			Server: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}},
		}}},
	}}}
	gueltig := sqliteDatei(t, rec)
	lockere(t, gueltig)
	setze := func(position, felder string) string {
		return `UPDATE message SET fields = '` + felder + `' WHERE session_id = 1 AND sequence = ` + position
	}
	for _, c := range []struct {
		name   string
		yaml   string
		sqlite string
		want   string
	}{
		{"gültig", yamlText("", "s"), setze("1 AND position = 1", `{"tag":"SELECT 1"}`), ""},
		{"unbekannter Schlüssel", yamlText("\n            bogus: 1", "s"), setze("1 AND position = 1", `{"tag":"SELECT 1","bogus":1}`), beschaedigt},
		{"unbekannter Schlüssel an Client", yamlText("", "s\n                bogus: 1"), setze("2 AND position = 1", `{"statement":"s","sql":"SELECT 1","param_types":[],"bogus":1}`), beschaedigt},
		{"null an Client", yamlText("", "null"), setze("2 AND position = 1", `{"statement":null,"sql":"SELECT 1","param_types":[]}`), beschaedigt},
		{"doppelter Schlüssel", yamlText("\n            tag: andere", "s"), setze("1 AND position = 1", `{"tag":"SELECT 1","tag":"andere"}`), beschaedigt},
		{"doppelter Schlüssel an Client", yamlText("", "s\n                statement: t"), setze("2 AND position = 1", `{"statement":"s","statement":"t","sql":"SELECT 1","param_types":[]}`), beschaedigt},
	} {
		ys := filepath.Join(t.TempDir(), "rec.yaml")
		if err := os.WriteFile(ys, []byte(c.yaml), 0o644); err != nil {
			t.Fatal(err)
		}
		sq := kopiere(t, gueltig)
		rohSchreibe(t, sq, c.sqlite)
		gotY, errY := laden(ys)
		gotS, errS := laden(sq)
		if codeVon(errY) != c.want || codeVon(errS) != c.want {
			t.Errorf("%s: yaml %v, sqlite %v, erwartet %q", c.name, errY, errS, c.want)
		}
		if c.want == "" && !reflect.DeepEqual(gotY, gotS) {
			t.Errorf("%s: aus yaml\n%#v\naus sqlite\n%#v", c.name, gotY, gotS)
		}
	}
	// Eine unbekannte Version: in beiden Formaten PGR-E3002.
	ys := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(ys, []byte(strings.Replace(yamlText("", "s"), "version: 1", "version: 2", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	sq := kopiere(t, gueltig)
	rohSchreibe(t, sq, `UPDATE meta SET value = '2' WHERE key = 'version'`)
	if _, err := laden(ys); codeVon(err) != unbekannteV {
		t.Errorf("yaml: %v", err)
	}
	if _, err := laden(sq); codeVon(err) != unbekannteV {
		t.Errorf("sqlite: %v", err)
	}
}
