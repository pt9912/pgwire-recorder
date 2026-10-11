package recording_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go.yaml.in/yaml/v3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// neueDatei schreibt rec über einen neuen SQLite-Adapter nach rec.sqlite in einem
// neuen Verzeichnis und liefert den Adapter und den Pfad.
func neueDatei(t *testing.T, rec model.Recording) (*recording.SQLite, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rec.sqlite")
	s := &recording.SQLite{}
	if err := s.Write(context.Background(), path, rec); err != nil {
		t.Fatal(err)
	}
	return s, path
}

// Abdeckung: LH-FA-22/Happy — die Datei von sqlite trägt die Tabellen und Zeilen
// nach SPEC-043, gelesen mit rohem SQL und ohne den Leser des Adapters: meta mit
// Formatkennung und Version, je Session Startup-Parameter und Server-Parameter als
// JSON-Objekt, je Interaktion Art, SQL (nur bei query) und offset_ms (NULL ohne
// Zeitangaben, 0 ist 0), je Nachricht position über beide Richtungen, direction,
// group_no (nur bei extended, ab 1) und kind, in jeder Gruppe die Client-Nachrichten
// vor den Server-Nachrichten, und fields als JSON-Objekt ohne den Schlüssel type;
// PRAGMA user_version bleibt 0.
func TestSQLiteInhalt(t *testing.T) {
	rec := extendedBeispiel()
	null := int64(0)
	rec.Sessions[0].Interactions[0].OffsetMS = &null
	_, path := neueDatei(t, rec)
	db := roh(t, path)

	check := func(name, anfrage string, want ...string) {
		t.Helper()
		if got := zeilen(t, db, anfrage); !reflect.DeepEqual(got, want) {
			t.Errorf("%s:\n got %q\nwant %q", name, got, want)
		}
	}
	check("meta", `SELECT key, value FROM meta ORDER BY key`, "format|pgwire-recorder", "version|1")
	check("user_version", `PRAGMA user_version`, "0")
	check("session", `SELECT id, startup, server_parameters FROM session`,
		`1|{"database":"app","user":"app"}|{"server_version":"17.0"}`)
	check("interaction", `SELECT session_id, sequence, type, offset_ms, sql FROM interaction ORDER BY sequence`,
		`1|1|query|0|SELECT 1, NULL, '\x00ff'::bytea`,
		`1|2|extended|NULL|NULL`)
	check("message", `SELECT session_id, sequence, position, direction, group_no, kind, fields FROM message ORDER BY sequence, position`,
		`1|1|1|server|NULL|row_description|{"fields":[{"name":"?column?","table_oid":0,"column_number":0,"type_oid":23,"type_size":4,"type_modifier":-1,"format":0}]}`,
		`1|1|2|server|NULL|data_row|{"values":[{"text":"1"},{"null":true},{"base64":"AP8="}]}`,
		`1|1|3|server|NULL|command_complete|{"tag":"SELECT 1"}`,
		`1|1|4|server|NULL|notice_response|{"notice":{"C":"00000","M":"Hinweis","S":"NOTICE"}}`,
		`1|1|5|server|NULL|ready_for_query|{"tx_status":"I"}`,
		`1|2|1|client|1|parse|{"statement":"s1","sql":"SELECT name FROM users WHERE id = $1","param_types":[23]}`,
		`1|2|2|client|1|describe|{"target":"statement","name":"s1"}`,
		`1|2|3|client|1|flush|{}`,
		`1|2|4|server|1|parse_complete|{}`,
		`1|2|5|server|1|parameter_description|{"param_types":[23]}`,
		`1|2|6|server|1|no_data|{}`,
		`1|2|7|client|2|close|{"target":"portal","name":"p0"}`,
		`1|2|8|client|2|flush|{}`,
		`1|2|9|client|3|bind|{"portal":"p1","statement":"s1","param_formats":[0,1],"params":[{"null":true},{"base64":"AP8="}],"result_formats":[1]}`,
		`1|2|10|client|3|execute|{"portal":"p1","max_rows":5}`,
		`1|2|11|client|3|sync|{}`,
		`1|2|12|server|3|close_complete|{}`,
		`1|2|13|server|3|bind_complete|{}`,
		`1|2|14|server|3|data_row|{"values":[{"text":"alice"}]}`,
		`1|2|15|server|3|portal_suspended|{}`,
		`1|2|16|server|3|parameter_status|{"name":"application_name","value":"x"}`,
		`1|2|17|server|3|ready_for_query|{"tx_status":"T"}`)
}

// Abdeckung: LH-FA-22/Boundary — eine Nachricht ohne Felder steht mit {}, nie mit
// NULL; eine einfache Anfrage ohne Text (leere Anfrage) trägt in interaction.sql den
// leeren Text und nicht NULL; offset_ms ohne Zeitangaben ist NULL und mit Wert 7 die
// Zahl 7; Startup und Server-Parameter ohne Einträge, auch als nil, stehen als {}.
func TestSQLiteLeereWerte(t *testing.T) {
	rec := model.NewRecording()
	sieben := int64(7)
	rec.Sessions = []model.Session{{
		ID: 1,
		Interactions: []model.Interaction{
			{Sequence: 1, Request: model.Request{Type: model.RequestQuery}, Responses: []model.Response{
				{Type: model.ResponseEmptyQueryResponse}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}},
			{Sequence: 2, OffsetMS: &sieben, Request: model.Request{Type: model.RequestQuery, SQL: "x"}, Responses: []model.Response{
				{Type: model.ResponseReadyForQuery, TxStatus: "I"}}},
		},
	}}
	_, path := neueDatei(t, rec)
	db := roh(t, path)
	for name, c := range map[string]struct {
		anfrage string
		want    []string
	}{
		"session":   {`SELECT startup, server_parameters FROM session`, []string{`{}|{}`}},
		"sql":       {`SELECT sequence, typeof(sql), sql, typeof(offset_ms), offset_ms FROM interaction ORDER BY sequence`, []string{`1|text||null|NULL`, `2|text|x|integer|7`}},
		"fields":    {`SELECT kind, typeof(fields), fields FROM message WHERE sequence = 1 AND position = 1`, []string{`empty_query_response|text|{}`}},
		"ohne NULL": {`SELECT count(*) FROM message WHERE fields IS NULL`, []string{`0`}},
		"group_no":  {`SELECT count(*) FROM message WHERE group_no IS NOT NULL`, []string{`0`}},
		"direction": {`SELECT DISTINCT direction FROM message`, []string{`server`}},
	} {
		if got := zeilen(t, db, c.anfrage); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s:\n got %q\nwant %q", name, got, c.want)
		}
	}
}

// Abdeckung: LH-FA-22/Boundary — meta trägt die Zeile empty_sessions mit dem Text
// true nur bei einer Aufzeichnung mit leeren Sessions; die leere Session steht
// ohne Interaktionen in session.
func TestSQLiteEmptySessions(t *testing.T) {
	for _, leer := range []bool{false, true} {
		rec := sessionsMit(1)
		rec.EmptySessions = leer
		if leer {
			rec.Sessions = append(rec.Sessions, model.Session{ID: 2})
		}
		_, path := neueDatei(t, rec)
		db := roh(t, path)
		want := []string{"format|pgwire-recorder", "version|1"}
		if leer {
			want = []string{"empty_sessions|true", "format|pgwire-recorder", "version|1"}
		}
		if got := zeilen(t, db, `SELECT key, value FROM meta ORDER BY key`); !reflect.DeepEqual(got, want) {
			t.Errorf("EmptySessions=%v: meta %q, erwartet %q", leer, got, want)
		}
		if got := zeilen(t, db, `SELECT count(*) FROM interaction WHERE session_id = 2`); leer && got[0] != "0" {
			t.Errorf("leere Session mit Interaktionen: %q", got)
		}
	}
}

// Abdeckung: LH-FA-22/Boundary — fields trägt dieselben Schlüssel und Werte wie die
// YAML-Darstellung derselben Nachricht (ohne type), für jede Art einer Server- und
// Client-Nachricht der Beispiele.
func TestSQLiteFieldsGleichYAML(t *testing.T) {
	rec := extendedBeispiel()
	_, path := neueDatei(t, rec)

	yamlText, err := recording.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Sessions []struct {
			Interactions []struct {
				Responses []map[string]any `yaml:"responses"`
				Groups    []struct {
					Client []map[string]any `yaml:"client"`
					Server []map[string]any `yaml:"server"`
				} `yaml:"groups"`
			} `yaml:"interactions"`
		} `yaml:"sessions"`
	}
	if err := yaml.Unmarshal(yamlText, &doc); err != nil {
		t.Fatal(err)
	}
	var erwartet []string
	for _, i := range doc.Sessions[0].Interactions {
		nachrichten := i.Responses
		for _, g := range i.Groups {
			nachrichten = append(append(nachrichten, g.Client...), g.Server...)
		}
		for _, n := range nachrichten {
			delete(n, "type")
			b, err := json.Marshal(n)
			if err != nil {
				t.Fatal(err)
			}
			erwartet = append(erwartet, string(b))
		}
	}
	var got []string
	for _, f := range zeilen(t, roh(t, path), `SELECT fields FROM message ORDER BY sequence, position`) {
		var v any
		if err := json.Unmarshal([]byte(f), &v); err != nil {
			t.Fatalf("%q: %v", f, err)
		}
		b, _ := json.Marshal(v)
		got = append(got, string(b))
	}
	if !reflect.DeepEqual(got, erwartet) {
		t.Fatalf("fields gegen die YAML-Darstellung:\n got %q\nwant %q", got, erwartet)
	}
	if len(got) != 22 {
		t.Fatalf("%d Nachrichten, erwartet 22", len(got))
	}
}

// Abdeckung: LH-FA-22/Happy — der erste Schreibvorgang legt eine gültige Aufzeichnung
// ohne Sessions an (Schema, meta); jeder weitere ergänzt genau die Sessions, die die
// Datei noch nicht trägt, in der Reihenfolge ihrer Nummern, ohne die Datei zu
// ersetzen (dieselbe Datei) und ohne eine schon geschriebene Session zu verändern;
// ein Aufruf ohne neue Session ändert nichts.
func TestSQLiteErgaenzt(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "rec.sqlite")
	s := &recording.SQLite{}
	if err := s.Prepare(ctx, path, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("Zielpfad nach Prepare: %v", err)
	}
	var vorher os.FileInfo
	for n, want := range []string{"", "1", "1,2", "1,2", "1,2,3"} {
		rec := sessionsMit([]int{0, 1, 2, 2, 3}[n])
		if err := s.Write(ctx, path, rec); err != nil {
			t.Fatalf("Write %d: %v", n, err)
		}
		if got := sessionIDs(t, path); got != want {
			t.Fatalf("Write %d: Sessions %q, erwartet %q", n, got, want)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if vorher != nil && !os.SameFile(vorher, info) {
			t.Fatalf("Write %d ersetzt die Datei", n)
		}
		vorher = info
	}
	db := roh(t, path)
	if got := zeilen(t, db, `SELECT session_id, sql FROM interaction ORDER BY session_id`); !reflect.DeepEqual(got, []string{"1|SELECT 1", "2|SELECT 2", "3|SELECT 3"}) {
		t.Fatalf("Interaktionen: %q", got)
	}
	if got := zeilen(t, db, `SELECT count(*) FROM meta`); got[0] != "2" {
		t.Fatalf("meta: %q", got)
	}
}

// Abdeckung: LH-FA-22/Negative — die Datei enthält zu jedem Zeitpunkt nur
// vollständige Sessions: scheitert das Einfügen mitten in einer Session (zwei
// Interaktionen mit derselben Nummer), ist das PGR-E3001, die Transaktion ist
// zurückgerollt, die Datei trägt nur die zuvor ergänzten Sessions mit allen
// Zeilen, die Sessions danach sind nicht ergänzt, und der nächste Schreibvorgang
// ergänzt, was fehlt.
func TestSQLiteHalbeSessionNie(t *testing.T) {
	ctx := context.Background()
	s, path := neueDatei(t, sessionsMit(1))
	rec := sessionsMit(3)
	rec.Sessions[1].Interactions = append(rec.Sessions[1].Interactions, rec.Sessions[1].Interactions[0])
	err := s.Write(ctx, path, rec)
	if code(err) != model.CodeRecordingIO {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingIO, err)
	}
	db := roh(t, path)
	for tabelle, want := range map[string]string{"session": "1", "interaction": "1", "message": "2"} {
		if got := zeilen(t, db, `SELECT count(*) FROM `+tabelle); got[0] != want {
			t.Errorf("nach dem Fehlschlag trägt %s %s Zeilen, erwartet %s", tabelle, got[0], want)
		}
	}
	if err := s.Write(ctx, path, sessionsMit(3)); err != nil {
		t.Fatalf("Nachholen: %v", err)
	}
	if got := sessionIDs(t, path); got != "1,2,3" {
		t.Fatalf("nach dem Nachholen Sessions %q", got)
	}
	for tabelle, want := range map[string]string{"session": "3", "interaction": "3", "message": "6"} {
		if got := zeilen(t, db, `SELECT count(*) FROM `+tabelle); got[0] != want {
			t.Errorf("nach dem Nachholen trägt %s %s Zeilen, erwartet %s", tabelle, got[0], want)
		}
	}
}

// Abdeckung: LH-FA-22/Negative — jede Session hat ihre eigene Transaktion: Scheitert die
// dritte von drei fehlenden Sessions, bleibt die zweite, die vor ihr abgeschlossen war,
// in der Datei.
func TestSQLiteTransaktionJeSession(t *testing.T) {
	ctx := context.Background()
	s, path := neueDatei(t, sessionsMit(1))
	rec := sessionsMit(4)
	rec.Sessions[2].Interactions = append(rec.Sessions[2].Interactions, rec.Sessions[2].Interactions[0])
	if err := s.Write(ctx, path, rec); code(err) != model.CodeRecordingIO {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingIO, err)
	}
	if got := sessionIDs(t, path); got != "1,2" {
		t.Fatalf("Sessions %q, erwartet 1,2", got)
	}
}

// Abdeckung: LH-FA-22/Negative — die Datei wird vor jedem Ergänzen geprüft, auch
// wenn nichts zu ergänzen ist: Fehlt sie, ist das PGR-E3001, und Write legt sie
// nicht an; trägt session nicht genau die Nummern 1 bis zur zuletzt ergänzten Session
// (eine andere Datei, eine Session zu viel oder zu wenig, eine Lücke), ist das
// PGR-E3001, und die Datei bleibt unberührt.
func TestSQLiteDateiGeprueft(t *testing.T) {
	ctx := context.Background()
	for name, c := range map[string]struct {
		eingriff func(t *testing.T, path string)
		rec      model.Recording
	}{
		"fehlt": {func(t *testing.T, path string) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
		}, sessionsMit(3)},
		"fehlt ohne Ergänzung": {func(t *testing.T, path string) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
		}, sessionsMit(2)},
		"andere Datei": {func(t *testing.T, path string) {
			ersetze(t, path, sessionsMit(1))
		}, sessionsMit(3)},
		"andere Datei ohne Ergänzung": {func(t *testing.T, path string) {
			ersetze(t, path, sessionsMit(3))
		}, sessionsMit(2)},
		"Session zu viel": {func(t *testing.T, path string) {
			ersetze(t, path, sessionsMit(3))
		}, sessionsMit(2)},
		"Session zu wenig": {func(t *testing.T, path string) {
			ersetze(t, path, sessionsMit(1))
		}, sessionsMit(2)},
		"Lücke": {func(t *testing.T, path string) {
			rohSchreibe(t, path, `DELETE FROM session WHERE id = 1`)
		}, sessionsMit(2)},
		"Nummer verschoben": {func(t *testing.T, path string) {
			rohSchreibe(t, path, `UPDATE session SET id = 3 WHERE id = 2`)
		}, sessionsMit(2)},
	} {
		s, path := neueDatei(t, sessionsMit(2))
		c.eingriff(t, path)
		vorher, _ := os.ReadFile(path)
		err := s.Write(ctx, path, c.rec)
		if code(err) != model.CodeRecordingIO {
			t.Errorf("%s: erwartet %s, erhalten %v", name, model.CodeRecordingIO, err)
		}
		nachher, _ := os.ReadFile(path)
		if string(vorher) != string(nachher) {
			t.Errorf("%s: Datei verändert", name)
		}
		if name == "fehlt" || name == "fehlt ohne Ergänzung" {
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Errorf("%s: Write legt die Datei an: %v", name, err)
			}
		}
	}
}
