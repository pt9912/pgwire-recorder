package recording_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

func sEnde() model.Response { return model.Response{Type: model.ResponseConnectionClosed} }

func sFatal() model.Response {
	return model.Response{Type: model.ResponseErrorResponse, Fields: map[string]string{"V": "FATAL", "S": "FATAL", "C": "57P01", "M": "x"}}
}

func sCC() model.Response { return model.Response{Type: model.ResponseCommandComplete, Tag: "x"} }

func sRFQ() model.Response { return model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"} }

// sEinfach ist eine einfache Interaktion mit den Antworten.
func sEinfach(nummer int, rs ...model.Response) model.Interaction {
	return model.Interaction{Sequence: nummer, Request: model.Request{Type: model.RequestQuery, SQL: "x"}, Responses: rs}
}

// sGruppe ist eine Gruppe, deren einzige Client-Nachricht abschluss (sync oder flush) ist.
func sGruppe(abschluss model.ClientMessageType, server ...model.Response) model.Group {
	return model.Group{Client: []model.ClientMessage{{Type: abschluss}}, Server: server}
}

func sExtended(nummer int, gs ...model.Group) model.Interaction {
	return model.Interaction{Sequence: nummer, Request: model.Request{Type: model.RequestExtended}, Groups: gs}
}

// sSession ist eine Session; user kennzeichnet die Aufzeichnung für den Fall.
func sSession(id int, user string, is ...model.Interaction) model.Session {
	return model.Session{ID: id, Startup: map[string]string{"user": user}, Interactions: is}
}

// sRec ist eine Aufzeichnung mit dem Feld version des Modells.
func sRec(feld int, ss ...model.Session) model.Recording {
	return model.Recording{Format: model.Format, Version: feld, Sessions: ss}
}

// schreibeUndLade schreibt rec mit dem Schreiber in eine Datei und liest Text und Laden.
func schreibeUndLade(t *testing.T, rec model.Recording) (string, model.Recording, error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rec.yaml")
	if err := (recording.YAML{}).Write(context.Background(), path, rec); err != nil {
		t.Fatalf("Write: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := recording.YAML{}.Load(context.Background(), path)
	return string(data), got, err
}

// Abdeckung: LH-FA-25/Happy — der Schreiber wählt die Version nach dem Inhalt: 2 mit einem Ende
// (Responses oder Gruppe einer Extended-Interaktion, in einer späteren oder einer nicht letzten
// Session), sonst 1 (auch ohne Sessions und nur mit empty_sessions); das Feld version des Modells
// bestimmt sie nicht (1 mit Ende schreibt 2, 2 ohne Ende schreibt 1, 0, 3 und ein negativer Wert ändern nichts).
func TestSchreibenVersionNachInhalt(t *testing.T) {
	mitEnde := sEinfach(1, sFatal(), sEnde())
	ohne := sEinfach(1, sCC(), sRFQ())
	cases := []struct {
		name string
		rec  model.Recording
		want string
	}{
		{"Feld 1 mit Ende", sRec(1, sSession(1, "u", mitEnde)), "2"},
		{"Feld 1 ohne Ende", sRec(1, sSession(1, "u", ohne)), "1"},
		{"Feld 2 ohne Ende", sRec(2, sSession(1, "u", ohne)), "1"},
		{"Feld 2 mit Ende", sRec(2, sSession(1, "u", mitEnde)), "2"},
		{"Feld 0 mit Ende", sRec(0, sSession(1, "u", mitEnde)), "2"},
		{"Feld 0 ohne Ende", sRec(0, sSession(1, "u", ohne)), "1"},
		{"Feld 3 mit Ende", sRec(3, sSession(1, "u", mitEnde)), "2"},
		{"Feld 3 ohne Ende", sRec(3, sSession(1, "u", ohne)), "1"},
		{"negatives Feld mit Ende", sRec(-1, sSession(1, "u", mitEnde)), "2"},
		{"negatives Feld ohne Ende", sRec(-1, sSession(1, "u", ohne)), "1"},
		{"Ende nur in der letzten Session", sRec(2, sSession(1, "u", ohne), sSession(2, "u", mitEnde)), "2"},
		{"Ende nur in der ersten von zwei Sessions", sRec(2, sSession(1, "u", mitEnde), sSession(2, "u", ohne)), "2"},
		{"ohne Sessions", sRec(1), "1"},
		{"nur empty_sessions", func() model.Recording {
			r := sRec(1, sSession(1, "u"))
			r.EmptySessions = true
			return r
		}(), "1"},
		{"Ende in einer Extended-Gruppe", sRec(2, sSession(1, "u", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde())))), "2"},
		{"Ende in einer einfachen Interaktion", sRec(2, sSession(1, "u", ohne, mitEnde)), "2"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := recording.Marshal(c.rec)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(string(b), "format: pgwire-recorder\nversion: "+c.want+"\n") {
				t.Fatalf("erwartet Version %s:\n%s", c.want, b)
			}
		})
	}
}

// Abdeckung: LH-FA-25/Negative — der Schreiber prüft die Form des Endes nicht: Eine Aufzeichnung mit
// dem Ende an einer falschen Stelle (Anfang, Mitte, nicht letzte Interaktion, nicht letzte Gruppe,
// mitten in der letzten Gruppe), zweimal in einer Session, mit ready_for_query daneben oder mit einem
// Feld neben type schreibt Version 2 ohne Fehler, und das Laden lehnt sie mit PGR-E3003 ab; das Ende
// an der richtigen Stelle, einfach, in der letzten Gruppe und in einer nicht letzten Session, lädt.
func TestSchreibenNimmtFormverletzungAn(t *testing.T) {
	cases := []struct {
		name string
		rec  model.Recording
		lade bool
	}{
		{"richtige Stelle, einfach", sRec(2, sSession(1, "u", sEinfach(1, sFatal(), sEnde()))), true},
		{"richtige Stelle, letzte Gruppe", sRec(2, sSession(1, "u", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde())))), true},
		{"richtige Stelle, nicht letzte Session", sRec(2, sSession(1, "u", sEinfach(1, sFatal(), sEnde())), sSession(2, "u", sEinfach(1, sCC(), sRFQ()))), true},
		{"Anfang der Folge", sRec(2, sSession(1, "u", sEinfach(1, sEnde(), sFatal(), sCC()))), false},
		{"Mitte der Folge", sRec(2, sSession(1, "u", sEinfach(1, sFatal(), sEnde(), sCC()))), false},
		{"nicht letzte Interaktion", sRec(2, sSession(1, "u", sEinfach(1, sFatal(), sEnde()), sEinfach(2, sCC(), sRFQ()))), false},
		{"nicht letzte Gruppe", sRec(2, sSession(1, "u", sExtended(1, sGruppe(model.ClientFlush, sFatal(), sEnde()), sGruppe(model.ClientSync, sCC(), sRFQ())))), false},
		{"mitten in der letzten Gruppe", sRec(2, sSession(1, "u", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde(), sCC())))), false},
		{"zwei Enden in einer Session", sRec(2, sSession(1, "u", sEinfach(1, sFatal(), sEnde()), sEinfach(2, sFatal(), sEnde()))), false},
		{"ready_for_query neben dem Ende", sRec(2, sSession(1, "u", sEinfach(1, sCC(), sRFQ(), sFatal(), sEnde()))), false},
		{"Feld neben type", sRec(2, sSession(1, "u", sEinfach(1, sFatal(), model.Response{Type: model.ResponseConnectionClosed, Tag: "x"}))), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			text, _, err := schreibeUndLade(t, c.rec)
			if !strings.HasPrefix(text, "format: pgwire-recorder\nversion: 2\n") {
				t.Fatalf("erwartet Version 2:\n%s", text)
			}
			if c.lade && err != nil {
				t.Fatalf("Laden: %v", err)
			}
			if !c.lade && code(err) != model.CodeRecordingBroken {
				t.Fatalf("erwartet %s beim Laden, erhalten %v", model.CodeRecordingBroken, err)
			}
		})
	}
}

// Abdeckung: LH-FA-25/Happy — das Ende steht in Blockform (ein Listenelement mit dem einzigen
// Schlüssel type, eingerückt wie die übrigen Nachrichten) in responses und in server einer Gruppe;
// ein Feld neben type im Modell wird geschrieben, nicht verworfen; eine Aufzeichnung ohne Ende
// ergibt dieselben Bytes wie ihr Text beim Laden und Schreiben.
func TestSchreibenEndeDarstellung(t *testing.T) {
	const kopf = "format: pgwire-recorder\nversion: 2\nsessions:\n  - id: 1\n    startup:\n      user: u\n    interactions:\n"
	const fehler = "                  C: 57P01\n                  M: x\n                  S: FATAL\n                  V: FATAL\n"
	const einfach = kopf + "      - sequence: 1\n        request:\n          type: query\n          sql: x\n        responses:\n" +
		"          - type: error_response\n            notice:\n" + "              C: 57P01\n              M: x\n              S: FATAL\n              V: FATAL\n" +
		"          - type: connection_closed\n"
	const erweitert = kopf + "      - sequence: 1\n        type: extended\n        groups:\n          - client:\n              - type: sync\n            server:\n" +
		"              - type: error_response\n                notice:\n" + fehler + "              - type: connection_closed\n"
	cases := []struct {
		name string
		rec  model.Recording
		want string
	}{
		{"in responses", sRec(2, sSession(1, "u", sEinfach(1, sFatal(), sEnde()))), einfach},
		{"in server einer Gruppe", sRec(2, sSession(1, "u", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde())))), erweitert},
		{"mit Feld neben type", sRec(2, sSession(1, "u", sEinfach(1, model.Response{Type: model.ResponseConnectionClosed, Tag: "x"}))),
			kopf + "      - sequence: 1\n        request:\n          type: query\n          sql: x\n        responses:\n          - type: connection_closed\n            tag: x\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := recording.Marshal(c.rec)
			if err != nil || string(b) != c.want {
				t.Fatalf("%v, erhalten:\n%s\nerwartet:\n%s", err, b, c.want)
			}
		})
	}
	t.Run("ohne Ende, Bytes bleiben", func(t *testing.T) {
		const text = "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    startup:\n      user: u\n    interactions:\n" +
			"      - sequence: 1\n        request:\n          type: query\n          sql: x\n        responses:\n" +
			"          - type: command_complete\n            tag: x\n          - type: ready_for_query\n            tx_status: I\n"
		rec, err := recording.Unmarshal([]byte(text))
		if err != nil {
			t.Fatal(err)
		}
		b, err := recording.Marshal(rec)
		if err != nil || string(b) != text {
			t.Fatalf("%v, erhalten:\n%s", err, b)
		}
	})
}

// Abdeckung: LH-FA-25/Happy, LH-FA-22/Happy, LH-QA-06/Messung — Schreiben und Laden einer Aufzeichnung
// und, von der geladenen Aufzeichnung aus, Schreiben und erneutes Laden liefern dieselbe Aufzeichnung mit Ende: einfach, Extended mit Sync-Gruppe, mit Flush-Gruppe (nach
// einer Gruppe mit flush), mit offset_ms an der Interaktion mit Ende, mit empty_sessions; eine Datei
// der Version 2 ohne Ende ist gleich bis auf version, die Version 1 schreibt und liest.
func TestRoundtripMitEnde(t *testing.T) {
	off := int64(7)
	mitOffset := sEinfach(1, sFatal(), sEnde())
	mitOffset.OffsetMS = &off
	leer := sRec(2, sSession(1, "rt-leer", sEinfach(1, sFatal(), sEnde())), sSession(2, "rt-leer"))
	leer.EmptySessions = true
	cases := []struct {
		name    string
		rec     model.Recording
		version int
	}{
		{"einfach", sRec(2, sSession(1, "rt-einfach", sEinfach(1, sCC(), sRFQ()), sEinfach(2, sFatal(), sEnde()))), 2},
		{"Sync-Gruppe", sRec(2, sSession(1, "rt-sync", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde())))), 2},
		{"Flush-Gruppe", sRec(2, sSession(1, "rt-flush", sExtended(1, sGruppe(model.ClientFlush, sCC()), sGruppe(model.ClientFlush, sFatal(), sEnde())))), 2},
		{"offset_ms", sRec(2, sSession(1, "rt-offset", mitOffset)), 2},
		{"empty_sessions", leer, 2},
		{"Version 2 ohne Ende", sRec(2, sSession(1, "rt-v2", sEinfach(1, sCC(), sRFQ()))), 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, got, err := schreibeUndLade(t, c.rec)
			if err != nil {
				t.Fatal(err)
			}
			want := c.rec
			want.Version = c.version
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Roundtrip weicht ab:\n got %#v\nwant %#v", got, want)
			}
			_, again, err := schreibeUndLade(t, got)
			if err != nil {
				t.Fatalf("erneutes Laden: %v", err)
			}
			if !reflect.DeepEqual(again, got) {
				t.Fatalf("Laden, Schreiben, Laden weicht ab:\n got %#v\nwant %#v", again, got)
			}
		})
	}
}
