package recording

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
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

// Abdeckung: LH-QA-06/Messung — die Aufzeichnung trägt
// Formatkennung und Version und lässt sich per Roundtrip laden; Binärwerte und
// NULL bleiben erhalten, und kein temporärer Rest bleibt liegen.
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

// Gleiche Aufzeichnungen ergeben gleiche Bytes
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

// Eine vorhandene Zieldatei ohne --force und ein
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

// Eine neue Aufzeichnung erhält die Rechte nach der umask, eine ersetzte behält
// ihre bisherigen (SPEC-033); ersetzt wird per Umbenennen, also auch eine
// schreibgeschützte Datei in einem beschreibbaren Verzeichnis (LH-FA-07.a).
func TestWriteRechteUndAtomar(t *testing.T) {
	ctx := context.Background()
	alt := syscall.Umask(0o077)
	defer syscall.Umask(alt)

	neuPfad := filepath.Join(t.TempDir(), "neu.yaml")
	if err := (YAML{}).Write(ctx, neuPfad, beispiel()); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(neuPfad); info.Mode().Perm() != 0o600 {
		t.Fatalf("neue Datei unter umask 077: %v", info.Mode().Perm())
	}

	vorhanden := filepath.Join(t.TempDir(), "vorhanden.yaml")
	if err := os.WriteFile(vorhanden, []byte("alt"), 0o400); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(vorhanden, 0o440); err != nil {
		t.Fatal(err)
	}
	if err := (YAML{}).Write(ctx, vorhanden, beispiel()); err != nil {
		t.Fatalf("schreibgeschützte Datei nicht ersetzt: %v", err)
	}
	if info, _ := os.Stat(vorhanden); info.Mode().Perm() != 0o440 {
		t.Fatalf("Rechte der ersetzten Datei: %v", info.Mode().Perm())
	}
	if got, err := (YAML{}).Load(ctx, vorhanden); err != nil || len(got.Sessions) != 1 {
		t.Fatalf("ersetzte Datei: %v %v", got, err)
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

// Abdeckung: LH-FA-07/Negative, LH-QA-06/Messung — eine fremde Formatkennung,
// eine unbekannte Version, ein unbekannter Typ, ein widersprüchlicher Wert oder
// eine abgeschnittene Datei machen die Aufzeichnung zu einer beschädigten.
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
		{"Sessions fehlen (abgeschnitten nach version)", "format: pgwire-recorder\nversion: 1\n", model.CodeRecordingBroken},
		{"abgeschnitten nach data_row", "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        request: {type: query, sql: x}\n        responses:\n          - type: data_row\n            values:\n              - text: a\n", model.CodeRecordingBroken},
		{"abgeschnitten nach request", "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        request: {type: query, sql: x}\n", model.CodeRecordingBroken},
		{"Session ohne Interaktion", "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions: []\n", model.CodeRecordingBroken},
		{"Lücke in den Kennungen", "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 2\n    interactions:\n      - sequence: 1\n        request: {type: query, sql: x}\n        responses:\n          - type: ready_for_query\n            tx_status: I\n", model.CodeRecordingBroken},
		{"doppelter Schlüssel im Wert", "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        request: {type: query, sql: x}\n        responses:\n          - type: data_row\n            values:\n              - {text: a, text: b}\n", model.CodeRecordingBroken},
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

// Die Felder offset_ms und empty_sessions der Version 1 überstehen den
// Roundtrip; ein negativer offset_ms ist beschädigt; eine Session ohne
// Interaktion ist nur mit empty_sessions gültig (LH-FA-12.a, LH-FA-21.a).
func TestFelderDerVersion1(t *testing.T) {
	ctx := context.Background()
	rec := beispiel()
	off := int64(42)
	rec.EmptySessions = true
	rec.Sessions[0].Interactions[0].OffsetMS = &off
	rec.Sessions = append(rec.Sessions, model.Session{ID: 2})
	path := filepath.Join(t.TempDir(), "rec.yaml")
	if err := (YAML{}).Write(ctx, path, rec); err != nil {
		t.Fatal(err)
	}
	got, err := YAML{}.Load(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.EmptySessions || got.Sessions[0].Interactions[0].OffsetMS == nil || *got.Sessions[0].Interactions[0].OffsetMS != 42 || len(got.Sessions) != 2 {
		t.Fatalf("Felder: %#v", got)
	}
	data, _ := os.ReadFile(path)
	for _, z := range []string{"empty_sessions: true", "offset_ms: 42"} {
		if !strings.Contains(string(data), z) {
			t.Fatalf("%q fehlt:\n%s", z, data)
		}
	}
	for name, text := range map[string]string{
		"negativer offset_ms":              "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        offset_ms: -1\n        request: {type: query, sql: x}\n        responses:\n          - type: ready_for_query\n            tx_status: I\n",
		"leere Session ohne Kennzeichnung": "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions: []\n",
	} {
		if _, err := Unmarshal([]byte(text)); code(err) != model.CodeRecordingBroken {
			t.Fatalf("%s: erwartet %s, erhalten %v", name, model.CodeRecordingBroken, err)
		}
	}
	if _, err := Unmarshal([]byte("format: pgwire-recorder\nversion: 1\nempty_sessions: true\nsessions:\n  - id: 1\n    interactions: []\n")); err != nil {
		t.Fatalf("leere Session mit Kennzeichnung: %v", err)
	}
}
