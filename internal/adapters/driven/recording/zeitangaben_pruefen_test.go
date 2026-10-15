package recording_test

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// zeitAufzeichnung baut eine aufnahme mit einer Session je Eintrag von werte und je Wert einer
// einfachen Interaktion; nil ist eine Interaktion ohne offset_ms.
func zeitAufzeichnung(werte ...[]*int64) model.Recording {
	rec := model.NewRecording()
	for si, sw := range werte {
		s := model.Session{ID: si + 1}
		for ii, w := range sw {
			s.Interactions = append(s.Interactions, model.Interaction{
				Sequence: ii + 1,
				OffsetMS: w,
				Request:  model.Request{Type: model.RequestQuery, SQL: "SELECT 1"},
				Responses: []model.Response{
					{Type: model.ResponseCommandComplete, Tag: "SELECT 1"},
					{Type: model.ResponseReadyForQuery, TxStatus: "I"},
				},
			})
		}
		rec.Sessions = append(rec.Sessions, s)
	}
	rec.EmptySessions = true
	return rec
}

func ms(v ...int64) []*int64 {
	out := make([]*int64, len(v))
	for i := range v {
		out[i] = &v[i]
	}
	return out
}

// ladeBeide schreibt rec in beiden Formaten und liefert je Format das Ergebnis des Ladens.
func ladeBeide(t *testing.T, rec model.Recording) map[string]struct {
	rec model.Recording
	err error
} {
	t.Helper()
	ctx := context.Background()
	ergebnis := map[string]struct {
		rec model.Recording
		err error
	}{}
	yamlPfad := filepath.Join(t.TempDir(), "rec.yaml")
	if err := (recording.YAML{}).Write(ctx, yamlPfad, rec); err != nil {
		t.Fatal(err)
	}
	r, err := recording.Leser{}.Load(ctx, yamlPfad)
	ergebnis["yaml"] = struct {
		rec model.Recording
		err error
	}{r, err}
	r, err = laden(sqliteDatei(t, rec))
	ergebnis["sqlite"] = struct {
		rec model.Recording
		err error
	}{r, err}
	return ergebnis
}

// Abdeckung: LH-FA-21/Negative — Werte von offset_ms, die innerhalb einer Session absteigen, machen
// die Aufzeichnung in yaml und in sqlite beschädigt (PGR-E3003); der Vergleich gilt je Session
// gegen den vorigen vorhandenen Wert, auch über eine Interaktion ohne Wert hinweg und auch bei
// einem Wert, der nur um 1 kleiner ist. Die Meldung nennt Session und Interaktion des kleineren
// Werts und keinen der Werte.
func TestZeitangabenAbsteigendBeschaedigt(t *testing.T) {
	for _, c := range []struct {
		name  string
		rec   model.Recording
		ort   string
		werte []string
	}{
		{"zwei Werte", zeitAufzeichnung(ms(4711, 4709)), "Session 1, Interaktion 2", []string{"4711", "4709"}},
		{"nur um 1 kleiner", zeitAufzeichnung(ms(4711, 4710)), "Session 1, Interaktion 2", []string{"4711", "4710"}},
		{"nicht benachbart", zeitAufzeichnung(ms(4100, 4900, 4200)), "Session 1, Interaktion 3", []string{"4100", "4900", "4200"}},
		{"über eine Interaktion ohne Wert", zeitAufzeichnung([]*int64{ms(4900)[0], nil, ms(4300)[0]}), "Session 1, Interaktion 3", []string{"4900", "4300"}},
		{"in der zweiten Session", zeitAufzeichnung(ms(4100, 4200), ms(4800, 4900, 4300)), "Session 2, Interaktion 3", []string{"4900", "4300"}},
		{"nach gleichen Werten", zeitAufzeichnung(ms(4800, 4800, 4799)), "Session 1, Interaktion 3", []string{"4800", "4799"}},
	} {
		for format, e := range ladeBeide(t, c.rec) {
			t.Run(c.name+", "+format, func(t *testing.T) {
				if codeVon(e.err) != model.CodeRecordingBroken {
					t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingBroken, e.err)
				}
				text := e.err.Error()
				if !strings.Contains(text, c.ort) {
					t.Errorf("Meldung ohne %q: %s", c.ort, text)
				}
				for _, w := range c.werte {
					if strings.Contains(text, w) {
						t.Errorf("Meldung nennt den Wert %s: %s", w, text)
					}
				}
			})
		}
	}
}

// Abdeckung: LH-FA-21/Boundary — gleiche Werte hintereinander, ein kleinerer Wert am Anfang der
// nächsten Session, eine Interaktion ohne Wert zwischen zwei Werten, eine Session mit einer oder
// ohne Interaktion, 0 und der größte Wert von int64 laden in yaml und in sqlite ohne Fehler und
// kommen unverändert zurück.
func TestZeitangabenGueltig(t *testing.T) {
	for _, c := range []struct {
		name string
		rec  model.Recording
	}{
		{"gleiche Werte", zeitAufzeichnung(ms(5, 5, 5))},
		{"nur 0", zeitAufzeichnung(ms(0, 0))},
		{"aufsteigend", zeitAufzeichnung(ms(0, 1, 40000))},
		{"kleinerer Wert in der nächsten Session", zeitAufzeichnung(ms(900, 900), ms(3, 3))},
		{"ohne Wert dazwischen, gleich", zeitAufzeichnung([]*int64{ms(7)[0], nil, ms(7)[0]})},
		{"ohne Wert dazwischen, größer", zeitAufzeichnung([]*int64{ms(7)[0], nil, ms(8)[0]})},
		{"ohne Wert am Anfang und am Ende", zeitAufzeichnung([]*int64{nil, ms(7)[0], nil})},
		{"keine Werte", zeitAufzeichnung([]*int64{nil, nil})},
		{"eine Interaktion", zeitAufzeichnung(ms(9))},
		{"Session ohne Interaktion", zeitAufzeichnung(ms(9), nil, ms(2))},
		{"größter Wert", zeitAufzeichnung(ms(math.MaxInt64-1, math.MaxInt64))},
	} {
		for format, e := range ladeBeide(t, c.rec) {
			t.Run(c.name+", "+format, func(t *testing.T) {
				if e.err != nil {
					t.Fatalf("erwartet ohne Fehler: %v", e.err)
				}
				if !reflect.DeepEqual(e.rec, c.rec) {
					t.Errorf("gelesen\n%#v\nerwartet\n%#v", e.rec, c.rec)
				}
			})
		}
	}
}

// Abdeckung: LH-FA-21/Negative — in yaml ist jede Schreibweise von offset_ms, die keine kanonische
// Dezimalzahl ist, beschädigt (PGR-E3003): Dezimalpunkt, Exponent, Hex, Oktal, führende Null,
// Vorzeichen, Unterstriche, Text in Anführungszeichen, ausdrücklicher Typ, Wahrheitswert, null,
// leerer Wert, Liste, Abbildung und ein Wert über int64; die Meldung nennt Session und
// Interaktion, keinen Wert und keine Zeile. Kanonische Zahlen laden.
func TestZeitangabenFormYAML(t *testing.T) {
	kopf := "format: pgwire-recorder\nversion: 1\nsessions:\n"
	eins := func(nr int, offset string) string {
		z := fmt.Sprintf("      - sequence: %d\n", nr)
		if offset != "" {
			z += "        offset_ms: " + offset + "\n"
		}
		return z + "        request: {type: query, sql: x}\n        responses:\n          - {type: command_complete, tag: x}\n          - {type: ready_for_query, tx_status: I}\n"
	}
	datei := func(offset string) string {
		return kopf + "  - id: 1\n    interactions:\n" + eins(1, "") + "  - id: 2\n    interactions:\n" + eins(1, "0") + eins(2, "0") + eins(3, offset)
	}
	for _, c := range []struct{ name, offset, wert string }{
		{"Dezimalpunkt 1.0", "4711.0", "4711.0"},
		{"Dezimalpunkt 1.5", "4711.5", "4711.5"},
		{"Exponent", "4711e3", "4711e3"},
		{"Hex", "0x4711", "0x4711"},
		{"Oktal", "0o4711", "0o4711"},
		{"führende Null", "04711", "04711"},
		{"Plus", "+4711", "+4711"},
		{"Minus", "-4711", "4711"},
		{"minus 0", "-0", "-0"},
		{"Unterstrich", "4_711", "4_711"},
		{"doppelte Anführungszeichen", `"4711"`, "4711"},
		{"einfache Anführungszeichen", `'4711'`, "4711"},
		{"ausdrücklicher Typ", "!!int 4711", "4711"},
		{"ausdrücklicher Typ Text", "!!str 4711", "4711"},
		{"true", "true", "true"},
		{"null", "null", "null"},
		{"Tilde", "~", "~"},
		{"leerer Wert", "", ""},
		{"Liste", "[4711]", "4711"},
		{"Abbildung", "{a: 4711}", "4711"},
		{"über int64", "9223372036854775808", "9223372036854775808"},
		{"weit über int64", "9223372036854775808000", "9223372036854775808000"},
	} {
		t.Run(c.name, func(t *testing.T) {
			offset := c.offset
			text := datei("X")
			if offset == "" {
				text = strings.Replace(text, "offset_ms: X", "offset_ms:", 1)
			} else {
				text = strings.Replace(text, "offset_ms: X", "offset_ms: "+offset, 1)
			}
			_, err := recording.Unmarshal([]byte(text))
			if code(err) != model.CodeRecordingBroken {
				t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingBroken, err)
			}
			msg := err.Error()
			if !strings.Contains(msg, "Session 2, Interaktion 3") {
				t.Errorf("Meldung ohne Ort: %s", msg)
			}
			if c.wert != "" && strings.Contains(msg, c.wert) {
				t.Errorf("Meldung nennt den Wert %q: %s", c.wert, msg)
			}
			if strings.Contains(strings.ToLower(msg), "line") || strings.Contains(msg, "Zeile") {
				t.Errorf("Meldung nennt eine Zeile: %s", msg)
			}
		})
	}
	for _, offset := range []string{"0", "4711", "9223372036854775807"} {
		if _, err := recording.Unmarshal([]byte(datei(offset))); err != nil {
			t.Errorf("offset_ms %s: %v", offset, err)
		}
	}
	// Ein Wert in der ersten Interaktion derselben Session und ein Wert in einer späteren
	// Session sind keine Ausnahme der Form.
	for _, text := range []string{
		kopf + "  - id: 1\n    interactions:\n" + eins(1, "1.5"),
		kopf + "  - id: 1\n    interactions:\n" + eins(1, "7") + "  - id: 2\n    interactions:\n" + eins(1, "7") + eins(2, "7") + eins(3, "true"),
	} {
		if _, err := recording.Unmarshal([]byte(text)); code(err) != model.CodeRecordingBroken {
			t.Errorf("erwartet %s, erhalten %v", model.CodeRecordingBroken, err)
		}
	}
}

// Abdeckung: LH-FA-21/Negative — in sqlite ist ein offset_ms, der negativ ist oder keine ganze
// Zahl, beschädigt (PGR-E3003); die Meldung nennt Session und Interaktion und keinen Wert.
func TestZeitangabenFormSQLite(t *testing.T) {
	basis := sqliteDatei(t, zeitAufzeichnung(ms(10, 20), ms(30, 40, 50)))
	lockere(t, basis)
	for _, c := range []struct{ name, sql, wert string }{
		{"negativ", `UPDATE interaction SET offset_ms = -4711 WHERE session_id = 2 AND sequence = 3`, "4711"},
		{"Text", `UPDATE interaction SET offset_ms = 'x4711' WHERE session_id = 2 AND sequence = 3`, "x4711"},
		{"REAL", `UPDATE interaction SET offset_ms = 4711.5 WHERE session_id = 2 AND sequence = 3`, "4711.5"},
		{"BLOB", `UPDATE interaction SET offset_ms = x'34373131' WHERE session_id = 2 AND sequence = 3`, "4711"},
	} {
		t.Run(c.name, func(t *testing.T) {
			path := kopiere(t, basis)
			rohSchreibe(t, path, c.sql)
			_, err := laden(path)
			if codeVon(err) != model.CodeRecordingBroken {
				t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingBroken, err)
			}
			if msg := err.Error(); !strings.Contains(msg, "Session 2, Interaktion 3") || strings.Contains(msg, c.wert) {
				t.Errorf("Meldung ohne Ort oder mit Wert: %s", msg)
			}
		})
	}
}

// Abdeckung: LH-FA-21/Negative — Lesbarkeit, Formatkennung und Version gehen dem Inhalt voran: eine
// Aufzeichnung mit kanonischer Version außer 1 und 2 (auch 0, negativ, über int64) ist PGR-E3002, auch
// wenn der Inhalt Fehler trägt; eine fehlende Version oder ein Wert, der keine kanonische
// Dezimalzahl ist, ist PGR-E3003 ohne Wert in der Meldung; eine falsche Formatkennung ist
// PGR-E3003, auch bei unbekannter Version.
func TestZeitangabenVersionVorForm(t *testing.T) {
	rumpf := "sessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        offset_ms: \"4711\"\n        request: {type: query, sql: x}\n        responses:\n          - {type: ready_for_query, tx_status: I}\n"
	const fremd = "format: pgwire-recorder\n"
	for _, c := range []struct{ name, text, code, enthaelt, nicht, wert string }{
		{"Version 3", fremd + "version: 3\n" + rumpf, model.CodeRecordingVersion, "Version 3", "", ""},
		{"Version 3, unbekannter Schlüssel", fremd + "version: 3\nunbekannt: 4711\n" + rumpf, model.CodeRecordingVersion, "Version 3", "", ""},
		{"Version 0", fremd + "version: 0\n" + rumpf, model.CodeRecordingVersion, "Version 0", "", ""},
		{"Version -1", fremd + "version: -1\n" + rumpf, model.CodeRecordingVersion, "Version -1", "", ""},
		{"Version über int64", fremd + "version: 9223372036854775808\n" + rumpf, model.CodeRecordingVersion, "Version 9223372036854775808", "", ""},
		{"andere Formatkennung, Version 2", "format: pgwire-recorden\nversion: 2\n" + rumpf, model.CodeRecordingBroken, "Formatkennung", "Version", ""},
		{"fehlende Formatkennung, Version 2", "version: 2\n" + rumpf, model.CodeRecordingBroken, "Formatkennung", "Version", ""},
		{"Version 2", fremd + "version: 2\n" + rumpf, model.CodeRecordingBroken, "Session 1, Interaktion 1: offset_ms", "", ""},
		{"Version 1", fremd + "version: 1\n" + rumpf, model.CodeRecordingBroken, "Session 1, Interaktion 1: offset_ms", "", ""},
		{"Version abc", fremd + "version: abc\n" + rumpf, model.CodeRecordingBroken, "version", "Version", "abc"},
		{"Version fehlt", fremd + rumpf, model.CodeRecordingBroken, "version", "Version", ""},
		{"Version 01", fremd + "version: 01\n" + rumpf, model.CodeRecordingBroken, "version", "Version", "01"},
		{"Version -0", fremd + "version: -0\n" + rumpf, model.CodeRecordingBroken, "version", "Version", "-0"},
		{"Version 1.0", fremd + "version: 1.0\n" + rumpf, model.CodeRecordingBroken, "version", "Version", "1.0"},
		{"Version in Anführungszeichen", fremd + "version: \"1\"\n" + rumpf, model.CodeRecordingBroken, "version", "Version", "1"},
		{"Version null", fremd + "version: null\n" + rumpf, model.CodeRecordingBroken, "version", "Version", "null"},
		{"Version 0x2", fremd + "version: 0x2\n" + rumpf, model.CodeRecordingBroken, "version", "Version", "0x2"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := recording.Unmarshal([]byte(c.text))
			if code(err) != c.code || !strings.Contains(err.Error(), c.enthaelt) || c.nicht != "" && strings.Contains(err.Error(), c.nicht) || c.wert != "" && strings.Contains(err.Error(), c.wert) {
				t.Errorf("erwartet %s mit %q ohne %q und ohne den Wert %q, erhalten %v", c.code, c.enthaelt, c.nicht, c.wert, err)
			}
		})
	}
}
