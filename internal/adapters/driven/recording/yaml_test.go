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

// extendedBeispiel trägt eine Extended-Interaktion mit Flush- und Sync-Gruppe
// nach einer einfachen Anfrage; sie nutzt jedes Feld einer Client-Nachricht,
// einen NULL- und einen Binärparameter und eine parameter_description.
func extendedBeispiel() model.Recording {
	rec := beispiel()
	rec.Sessions[0].Interactions = append(rec.Sessions[0].Interactions, model.Interaction{
		Sequence: 2,
		Request:  model.Request{Type: model.RequestExtended},
		Groups: []model.Group{
			{
				Client: []model.ClientMessage{
					{Type: model.ClientParse, Statement: "s1", SQL: "SELECT name FROM users WHERE id = $1", ParamTypes: []uint32{23}},
					{Type: model.ClientDescribe, Target: model.TargetStatement, Name: "s1"},
					{Type: model.ClientFlush},
				},
				Server: []model.Response{
					{Type: model.ResponseParseComplete},
					{Type: model.ResponseParameterDescription, ParamTypes: []uint32{23}},
					{Type: model.ResponseNoData},
				},
			},
			{
				Client: []model.ClientMessage{{Type: model.ClientClose, Target: model.TargetPortal, Name: "p0"}, {Type: model.ClientFlush}},
			},
			{
				Client: []model.ClientMessage{
					{Type: model.ClientBind, Portal: "p1", Statement: "s1", ParamFormats: []int16{0, 1}, Params: []model.Value{{Null: true}, {Bytes: []byte{0x00, 0xff}}}, ResultFormats: []int16{1}},
					{Type: model.ClientExecute, Portal: "p1", MaxRows: 5},
					{Type: model.ClientSync},
				},
				Server: []model.Response{
					{Type: model.ResponseCloseComplete},
					{Type: model.ResponseBindComplete},
					{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte("alice")}}},
					{Type: model.ResponsePortalSuspended},
					{Type: model.ResponseParameterStatus, Name: "application_name", Value: "x"},
					{Type: model.ResponseReadyForQuery, TxStatus: "T"},
				},
			},
		},
	})
	return rec
}

// Eine Extended-Interaktion übersteht den Roundtrip in allen Feldern; die Datei
// trägt die Form aus SPEC-041 (type: extended, groups mit client und server, die
// Felder jeder Client-Nachricht auch leer), bleibt deterministisch, und die
// einfache Anfrage davor behält ihre Form aus SPEC-002.
func TestExtendedRoundtrip(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "rec.yaml")
	want := extendedBeispiel()
	if err := (YAML{}).Write(ctx, path, want); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := YAML{}.Load(ctx, path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Roundtrip weicht ab:\n got %#v\nwant %#v", got, want)
	}
	data, _ := os.ReadFile(path)
	text := string(data)
	for _, s := range []string{
		"      - sequence: 1\n        request:\n",
		"      - sequence: 2\n        type: extended\n        groups:\n          - client:\n              - type: parse\n                statement: s1\n                sql: SELECT name FROM users WHERE id = $1\n                param_types: [23]\n",
		"              - type: describe\n                target: statement\n                name: s1\n              - type: flush\n            server:\n              - type: parse_complete\n              - type: parameter_description\n                param_types: [23]\n",
		"              - type: flush\n            server: []\n",
		"              - type: bind\n                portal: p1\n                statement: s1\n                param_formats: [0, 1]\n                params:\n                  - \"null\": true\n                  - base64: AP8=\n                result_formats: [1]\n",
		"              - type: execute\n                portal: p1\n                max_rows: 5\n              - type: sync\n",
	} {
		if !strings.Contains(text, s) {
			t.Errorf("Aufzeichnung enthält nicht:\n%s\n---\n%s", s, text)
		}
	}
	nochmal, err := Marshal(got)
	if err != nil || string(nochmal) != text {
		t.Fatalf("zweites Schreiben weicht ab (%v):\n%s\n---\n%s", err, nochmal, text)
	}
}

// Leere Felder einer Client-Nachricht stehen in der Datei (`portal: ""`,
// `param_types: []`, `max_rows: 0` wie in SPEC-041), ebenso die leeren
// `param_types` einer parameter_description, und eine leere Liste kommt als nil
// zurück; der Roundtrip bleibt damit gleich.
func TestExtendedLeereFelder(t *testing.T) {
	want := model.NewRecording()
	want.Sessions = []model.Session{{ID: 1, Interactions: []model.Interaction{{Sequence: 1, Request: model.Request{Type: model.RequestExtended}, Groups: []model.Group{{
		Client: []model.ClientMessage{{Type: model.ClientParse}, {Type: model.ClientBind}, {Type: model.ClientExecute}, {Type: model.ClientSync}},
		Server: []model.Response{{Type: model.ResponseParameterDescription}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}},
	}}}}}}
	leer, err := Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"statement: \"\"", "sql: \"\"", "param_types: []", "portal: \"\"", "param_formats: []", "params: []", "result_formats: []", "max_rows: 0", "- type: parameter_description\n                param_types: []\n"} {
		if !strings.Contains(string(leer), s) {
			t.Errorf("leere Felder: %q fehlt:\n%s", s, leer)
		}
	}
	got, err := Unmarshal(leer)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Roundtrip weicht ab:\n got %#v\nwant %#v", got, want)
	}
}

// Das Beispiel aus SPEC-041 wird gelesen, wie es dort steht.
func TestUnmarshalSpec041(t *testing.T) {
	data := `format: pgwire-recorder
version: 1
sessions:
  - id: 1
    interactions:
      - sequence: 1
        request: {type: query, sql: BEGIN}
        responses:
          - type: ready_for_query
            tx_status: T
      - sequence: 2
        type: extended
        groups:
          - client:
              - type: parse
                statement: "s1"
                sql: "SELECT name FROM users WHERE id = $1"
                param_types: [23]
              - type: bind
                portal: ""
                statement: "s1"
                param_formats: [0]
                params:
                  - text: "1"
                result_formats: [0]
              - type: describe
                target: portal
                name: ""
              - type: execute
                portal: ""
                max_rows: 0
              - type: sync
            server:
              - type: parse_complete
              - type: bind_complete
              - type: row_description
                # protokollrelevante Felder
              - type: data_row
                values:
                  - text: "alice"
              - type: command_complete
                tag: "SELECT 1"
              - type: ready_for_query
                tx_status: "I"
`
	rec, err := Unmarshal([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	i := rec.Sessions[0].Interactions[1]
	if i.Request.Type != model.RequestExtended || len(i.Groups) != 1 {
		t.Fatalf("Interaktion: %#v", i)
	}
	g := i.Groups[0]
	bind := g.Client[1]
	if len(g.Client) != 5 || len(g.Server) != 6 || g.Client[0].SQL != "SELECT name FROM users WHERE id = $1" ||
		!reflect.DeepEqual(g.Client[0].ParamTypes, []uint32{23}) || string(bind.Params[0].Bytes) != "1" ||
		!reflect.DeepEqual(bind.ResultFormats, []int16{0}) || g.Client[2].Target != model.TargetPortal ||
		g.Server[5].TxStatus != "I" {
		t.Fatalf("Gruppe: %#v", g)
	}
}

// extendedText ist eine Aufzeichnung mit einer Extended-Interaktion; gruppen
// steht eingerückt unter `groups:`.
func extendedText(gruppen string) string {
	return "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        type: extended\n        groups:\n" + gruppen
}

const syncGruppe = "          - client:\n              - type: sync\n            server:\n              - type: ready_for_query\n                tx_status: I\n"

// gruppe ist eine Gruppe in der Einrückung von extendedText; client und server
// sind die Einträge der beiden Listen, ohne server fehlt der Schlüssel.
func gruppe(client []string, server ...string) string {
	g := "          - client:\n"
	for _, c := range client {
		g += "              - " + c + "\n"
	}
	if server != nil {
		g += "            server:\n"
		for _, s := range server {
			g += "              - " + s + "\n"
		}
	}
	return g
}

// einfacherText ist eine Aufzeichnung mit einer einfachen Anfrage; zusatz steht
// in der Einrückung ihrer Schlüssel am Ende.
func einfacherText(antwort, zusatz string) string {
	return "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        request: {type: query, sql: x}\n        responses:\n          - " + antwort + "\n          - {type: ready_for_query, tx_status: I}\n" + zusatz
}

const (
	rfqNachricht   = "{type: ready_for_query, tx_status: I}"
	syncNachricht  = "{type: sync}"
	parseNachricht = "{type: parse, statement: s, sql: x, param_types: []}"
)

// Abdeckung: LH-FA-07/Negative — eine Extended-Interaktion, die abgeschnitten
// ist, eine unbekannte Nachricht, einen fremden, fehlenden oder mit null
// belegten Schlüssel oder eine Mischung mit der einfachen Form trägt, macht die
// Aufzeichnung zu einer beschädigten (PGR-E3003), ebenso param_types an einer
// anderen Antwort als parameter_description; die Meldung nennt die verletzte
// Regel.
func TestUnmarshalExtendedFehler(t *testing.T) {
	syncGruppe := gruppe([]string{syncNachricht}, rfqNachricht)
	cases := []struct{ name, data, want string }{
		{"unbekannter Interaktions-Typ", strings.Replace(extendedText(syncGruppe), "type: extended", "type: bogus", 1), "Interaktions-Typ \"bogus\""},
		{"type: query auf Interaktionsebene", strings.Replace(extendedText(syncGruppe), "type: extended", "type: query", 1), "Interaktions-Typ \"query\""},
		{"Anfrage-Typ extended in request", strings.Replace(einfacherText("{type: command_complete, tag: x}", ""), "type: query", "type: extended", 1), "Anfrage-Typ \"extended\""},
		{"extended mit request", extendedText(syncGruppe) + "        request: {type: query, sql: x}\n", "weder request noch responses"},
		{"extended mit responses", extendedText(syncGruppe) + "        responses: []\n", "weder request noch responses"},
		{"extended mit request: null", extendedText(syncGruppe) + "        request: null\n", "request: null"},
		{"extended mit responses: null", extendedText(syncGruppe) + "        responses:\n", "responses: null"},
		{"einfache Anfrage mit groups", einfacherText("{type: command_complete, tag: x}", "        groups: []\n"), "keine groups"},
		{"einfache Anfrage mit groups: null", einfacherText("{type: command_complete, tag: x}", "        groups: ~\n"), "groups: null"},
		{"einfache Anfrage ohne request", "format: pgwire-recorder\nversion: 1\nsessions:\n  - id: 1\n    interactions:\n      - sequence: 1\n        responses:\n          - " + rfqNachricht + "\n", "braucht request"},
		{"ohne groups", strings.TrimSuffix(extendedText(""), "        groups:\n") + "\n", "ohne Gruppe"},
		{"unbekannte Client-Nachricht", extendedText(gruppe([]string{"{type: copy_data}", syncNachricht}, rfqNachricht)), "Client-Nachricht \"copy_data\""},
		{"fremder Schlüssel an sync", extendedText(gruppe([]string{"{type: sync, sql: x}"}, rfqNachricht)), "\"sql\" gehört nicht zu sync"},
		{"Schlüssel von bind an parse", extendedText(gruppe([]string{"{type: parse, statement: s, sql: x, param_types: [], portal: p}", syncNachricht}, rfqNachricht)), "\"portal\" gehört nicht zu parse"},
		{"parse ohne sql", extendedText(gruppe([]string{"{type: parse, statement: s, param_types: []}", syncNachricht}, rfqNachricht)), "\"sql\" fehlt an parse"},
		{"parse mit sql: null", extendedText(gruppe([]string{"{type: parse, statement: s, sql: null, param_types: []}", syncNachricht}, rfqNachricht)), "\"sql\" fehlt an parse"},
		{"bind ohne params", extendedText(gruppe([]string{"{type: bind, portal: '', statement: s, param_formats: [], result_formats: []}", syncNachricht}, rfqNachricht)), "\"params\" fehlt an bind"},
		{"execute ohne max_rows", extendedText(gruppe([]string{"{type: execute, portal: ''}", syncNachricht}, rfqNachricht)), "\"max_rows\" fehlt an execute"},
		{"describe ohne name", extendedText(gruppe([]string{"{type: describe, target: portal}", syncNachricht}, rfqNachricht)), "\"name\" fehlt an describe"},
		{"close ohne target", extendedText(gruppe([]string{"{type: close, name: s}", syncNachricht}, rfqNachricht)), "\"target\" fehlt an close"},
		{"doppelter Schlüssel", extendedText(gruppe([]string{"{type: parse, statement: s, sql: a, sql: b, param_types: []}", syncNachricht}, rfqNachricht)), "already defined"},
		{"Zielart leer", extendedText(gruppe([]string{"{type: describe, target: '', name: s}", syncNachricht}, rfqNachricht)), "Zielart"},
		{"unbekannte Server-Nachricht", extendedText(gruppe([]string{syncNachricht}, "{type: copy_out_response}", rfqNachricht)), "Server-Nachricht \"copy_out_response\""},
		{"widersprüchlicher Parameter", extendedText(gruppe([]string{"{type: bind, portal: '', statement: s, param_formats: [], params: [{text: a, null: true}], result_formats: []}", syncNachricht}, rfqNachricht)), "mehr als einem"},
		{"parameter_description ohne param_types", extendedText(gruppe([]string{parseNachricht, syncNachricht}, "{type: parameter_description}", rfqNachricht)), "\"param_types\" fehlt an parameter_description"},
		{"param_types an parse_complete", extendedText(gruppe([]string{parseNachricht, syncNachricht}, "{type: parse_complete, param_types: [23]}", rfqNachricht)), "\"param_types\" gehört nicht zu parse_complete"},
		{"param_types in einfacher Anfrage", einfacherText("{type: command_complete, tag: x, param_types: [23]}", ""), "\"param_types\" gehört nicht zu command_complete"},
		{"abgeschnitten nach der Flush-Gruppe", extendedText(gruppe([]string{"{type: flush}"}, []string{}...)), "sync steht"},
		{"abgeschnitten nach den Client-Nachrichten", extendedText(gruppe([]string{syncNachricht})), "endet nicht mit ready_for_query"},
		{"abgeschnitten in den Server-Nachrichten", extendedText(gruppe([]string{syncNachricht}, "{type: command_complete, tag: x}")), "endet nicht mit ready_for_query"},
		{"abgeschnitten in den Client-Nachrichten", extendedText(gruppe([]string{parseNachricht})), "nur als letzte Nachricht"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Unmarshal([]byte(c.data))
			if code(err) != model.CodeRecordingBroken || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("erwartet %s mit %q, erhalten %v", model.CodeRecordingBroken, c.want, err)
			}
		})
	}
	for name, text := range map[string]string{
		"Sync-Gruppe":                          extendedText(syncGruppe),
		"parameter_description mit Typen":      extendedText(gruppe([]string{parseNachricht, syncNachricht}, "{type: parameter_description, param_types: [23]}", rfqNachricht)),
		"parameter_description ohne Parameter": extendedText(gruppe([]string{parseNachricht, syncNachricht}, "{type: parameter_description, param_types: []}", rfqNachricht)),
	} {
		if _, err := Unmarshal([]byte(text)); err != nil {
			t.Fatalf("Gegenstück %s ohne Fehler erwartet: %v", name, err)
		}
	}
}
