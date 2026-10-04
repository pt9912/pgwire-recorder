package recording

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

func beispiel() model.Recording {
	rec := model.NewRecording()
	rec.Sessions = []model.Session{{
		ID:               1,
		Startup:          map[string]string{"user": "app", "database": "app"},
		ServerParameters: map[string]string{"server_version": "17.0"},
		Interactions: []model.Interaction{{
			Sequence: 1,
			Request:  model.Request{Type: model.RequestQuery, SQL: "SELECT 1, NULL, '\\x00ff'::bytea"},
			Responses: []model.Response{
				{Type: model.ResponseRowDescription, Columns: []model.Column{{Name: "?column?", TypeOID: 23, TypeSize: 4, TypeModifier: -1}}},
				{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte("1")}, {Null: true}, {Bytes: []byte{0x00, 0xff}}}},
				{Type: model.ResponseCommandComplete, Tag: "SELECT 1"},
				{Type: model.ResponseNoticeResponse, Fields: map[string]string{"S": "NOTICE", "C": "00000", "M": "Hinweis"}},
				{Type: model.ResponseReadyForQuery, TxStatus: "I"},
			},
		}},
	}}
	return rec
}

// Abdeckung: LH-FA-07/Happy, LH-QA-06/Happy — die Aufzeichnung trägt
// Formatkennung und Version und lässt sich per Roundtrip laden; Binärwerte und
// NULL bleiben erhalten, die Datei ist danach lesbar für andere (0644) und kein
// temporärer Rest bleibt liegen.
func TestRoundtrip(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "rec.yaml")
	repo := YAML{}
	want := beispiel()

	if err := repo.Prepare(ctx, path, false); err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if err := repo.Write(ctx, path, want); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := repo.Load(ctx, path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Roundtrip weicht ab:\n got %#v\nwant %#v", got, want)
	}

	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o644 {
		t.Fatalf("Dateirechte: %v %v", info.Mode(), err)
	}
	reste, _ := filepath.Glob(filepath.Join(filepath.Dir(path), ".*"))
	if len(reste) != 0 {
		t.Fatalf("temporäre Reste: %v", reste)
	}

	data, _ := os.ReadFile(path)
	text := string(data)
	for _, s := range []string{"format: pgwire-recorder", "version: 1", "base64: AP8=", "\"null\": true"} {
		if !strings.Contains(text, s) {
			t.Errorf("Aufzeichnung enthält %q nicht:\n%s", s, text)
		}
	}
}

// Abdeckung: LH-QA-01/Boundary — gleiche Aufzeichnungen ergeben gleiche Bytes
// (SPEC-004).
func TestMarshalDeterministisch(t *testing.T) {
	a, err := Marshal(beispiel())
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		b, err := Marshal(beispiel())
		if err != nil {
			t.Fatal(err)
		}
		if string(a) != string(b) {
			t.Fatalf("Ausgabe nicht deterministisch:\n%s\n---\n%s", a, b)
		}
	}
}

// Abdeckung: LH-FA-07/Negative — eine vorhandene Zieldatei ohne --force und ein
// nicht beschreibbares Zielverzeichnis sind Startfehler.
func TestPrepareVorhandeneDatei(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := YAML{}.Prepare(ctx, path, false)
	if code(err) != model.CodeOutputExists {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeOutputExists, err)
	}
	if err := (YAML{}).Prepare(ctx, path, true); err != nil {
		t.Fatalf("mit --force: %v", err)
	}
	fehlt := filepath.Join(t.TempDir(), "gibt-es-nicht", "rec.yaml")
	if err := (YAML{}).Prepare(ctx, fehlt, false); code(err) != model.CodeRecordingIO {
		t.Fatalf("fehlendes Verzeichnis: erwartet %s, erhalten %v", model.CodeRecordingIO, err)
	}
}

// Die Form aus SPEC-041 mit ungequotetem `null:` wird gelesen.
func TestUnmarshalNullUngequotet(t *testing.T) {
	data := `format: pgwire-recorder
version: 1
sessions:
  - id: 1
    interactions:
      - sequence: 1
        request: {type: query, sql: SELECT NULL}
        responses:
          - type: data_row
            values:
              - null: true
              - text: a
          - type: ready_for_query
            tx_status: I
`
	rec, err := Unmarshal([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	v := rec.Sessions[0].Interactions[0].Responses[0].Values
	if !v[0].Null || string(v[1].Bytes) != "a" {
		t.Fatalf("Werte: %#v", v)
	}
}

// Abdeckung: LH-FA-07/Negative, LH-QA-06/Negative — eine fremde Formatkennung,
// eine unbekannte Version, ein unbekannter Typ oder ein widersprüchlicher Wert
// machen die Aufzeichnung zu einer beschädigten.
func TestUnmarshalFehler(t *testing.T) {
	cases := []struct {
		name, data, want string
	}{
		{"fremdes Format", "format: anderes\nversion: 1\nsessions: []\n", model.CodeRecordingBroken},
		{"unbekannte Version", "format: pgwire-recorder\nversion: 2\nsessions: []\n", model.CodeRecordingVersion},
		{"kein YAML", "format: [\n", model.CodeRecordingBroken},
		{"unbekannter Schlüssel", "format: pgwire-recorder\nversion: 1\nsessions: []\nextra: 1\n", model.CodeRecordingBroken},
		{"unbekannter Anfrage-Typ", "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        request: {type: bogus, sql: x}\n        responses: []\n", model.CodeRecordingBroken},
		{"unbekannter Antwort-Typ", "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        request: {type: query, sql: x}\n        responses:\n          - type: bogus\n", model.CodeRecordingBroken},
		{"text und base64", "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        request: {type: query, sql: x}\n        responses:\n          - type: data_row\n            values:\n              - {text: a, base64: YQ==}\n", model.CodeRecordingBroken},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Unmarshal([]byte(c.data))
			if code(err) != c.want {
				t.Fatalf("erwartet %s, erhalten %v", c.want, err)
			}
		})
	}
}

func code(err error) string {
	var me *model.Error
	if errors.As(err, &me) {
		return me.Code
	}
	return ""
}
