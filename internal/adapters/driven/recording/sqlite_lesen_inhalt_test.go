package recording_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// fall ist eine Änderung an der gelockerten Datei der lesefixture und die Klasse des
// Fehlers, den der Leser meldet ("" liest ohne Fehler).
type fall struct {
	name string
	sql  string
	want string
}

// pruefeFaelle wendet jeden Fall auf eine Kopie der gelockerten Datei an und prüft die
// Fehlerklasse des Lesens; die Datei ohne Änderung liest ohne Fehler.
func pruefeFaelle(t *testing.T, faelle []fall) {
	t.Helper()
	basis := sqliteDatei(t, lesefixture())
	lockere(t, basis)
	if got, err := laden(basis); err != nil || !reflect.DeepEqual(got, lesefixture()) {
		t.Fatalf("gelockerte Datei ohne Änderung: %v", err)
	}
	for _, f := range faelle {
		path := kopiere(t, basis)
		rohSchreibe(t, path, f.sql)
		_, err := laden(path)
		if got := codeVon(err); got != f.want {
			t.Errorf("%s: Meldungscode %q, erwartet %q (%v)", f.name, got, f.want, err)
		}
	}
}

const (
	beschaedigt  = model.CodeRecordingBroken
	unbekannteV  = model.CodeRecordingVersion
	ausNachricht = `session_id = 1 AND sequence = 2`
)

// Abdeckung: LH-FA-22/Negative, LH-QA-06/Messung — jede Spalte mit einer anderen gespeicherten Art als der ihres
// Typs, jede Nummer, die nicht ab 1 lückenlos zählt, jede verwaiste und jede doppelte Zeile und
// jede Zeile gegen die Regeln von SPEC-043 macht die Aufzeichnung beschädigt (PGR-E3003);
// die zugehörigen Fälle, die gültig sind, lesen ohne Fehler.
func TestLesenInhaltVerletzt(t *testing.T) {
	pruefeFaelle(t, []fall{
		// session
		{"session.id TEXT", `UPDATE session SET id = 'x' WHERE id = 2`, beschaedigt},
		{"session.id REAL", `UPDATE session SET id = 2.5 WHERE id = 2`, beschaedigt},
		{"session.id NULL", `UPDATE session SET id = NULL WHERE id = 2`, beschaedigt},
		{"Session-Nummer Lücke", `UPDATE session SET id = 3 WHERE id = 2; UPDATE interaction SET session_id = 3 WHERE session_id = 2; UPDATE message SET session_id = 3 WHERE session_id = 2`, beschaedigt},
		{"Session-Nummer beginnt bei 2", `UPDATE session SET id = id + 1; UPDATE interaction SET session_id = session_id + 1; UPDATE message SET session_id = session_id + 1`, beschaedigt},
		{"Session doppelt", `INSERT INTO session SELECT * FROM session WHERE id = 1`, beschaedigt},
		{"Session ohne Interaktion", `DELETE FROM message WHERE session_id = 2; DELETE FROM interaction WHERE session_id = 2`, beschaedigt},
		{"Session ohne Interaktion, empty_sessions", `DELETE FROM message WHERE session_id = 2; DELETE FROM interaction WHERE session_id = 2; INSERT INTO meta VALUES ('empty_sessions', 'true')`, ""},
		// interaction
		{"interaction.sequence TEXT", `UPDATE interaction SET sequence = 'x' WHERE session_id = 2`, beschaedigt},
		{"interaction.sequence REAL", `UPDATE interaction SET sequence = 1.5 WHERE session_id = 2`, beschaedigt},
		{"interaction.sequence NULL", `UPDATE interaction SET sequence = NULL WHERE session_id = 2`, beschaedigt},
		{"interaction.session_id TEXT", `UPDATE interaction SET session_id = 'x' WHERE session_id = 2`, beschaedigt},
		{"Interaktion ohne Session", `INSERT INTO interaction VALUES (9, 1, 'query', NULL, 'SELECT 1')`, beschaedigt},
		{"Interaktion und Nachrichten ohne Session", `UPDATE interaction SET session_id = 9 WHERE session_id = 2; UPDATE message SET session_id = 9 WHERE session_id = 2`, beschaedigt},
		{"Interaktions-Nummer Lücke", `UPDATE interaction SET sequence = 3 WHERE ` + ausNachricht + `; UPDATE message SET sequence = 3 WHERE ` + ausNachricht, beschaedigt},
		{"Interaktion doppelt", `INSERT INTO interaction SELECT * FROM interaction WHERE session_id = 2`, beschaedigt},
		{"interaction.type NULL", `UPDATE interaction SET type = NULL WHERE session_id = 2`, beschaedigt},
		{"interaction.type unbekannt", `UPDATE interaction SET type = 'foo' WHERE session_id = 2`, beschaedigt},
		{"interaction.type Großschreibung", `UPDATE interaction SET type = 'QUERY' WHERE session_id = 2`, beschaedigt},
		{"offset_ms TEXT", `UPDATE interaction SET offset_ms = 'x' WHERE session_id = 2`, beschaedigt},
		{"offset_ms negativ", `UPDATE interaction SET offset_ms = -1 WHERE session_id = 2`, beschaedigt},
		{"offset_ms REAL", `UPDATE interaction SET offset_ms = 1.5 WHERE session_id = 2`, beschaedigt},
		{"offset_ms BLOB", `UPDATE interaction SET offset_ms = x'31' WHERE session_id = 2`, beschaedigt},
		{"offset_ms 0", `UPDATE interaction SET offset_ms = 0 WHERE session_id = 2`, ""},
		{"sql NULL bei query", `UPDATE interaction SET sql = NULL WHERE session_id = 2`, beschaedigt},
		{"sql BLOB bei query", `UPDATE interaction SET sql = x'4142' WHERE session_id = 2`, beschaedigt},
		{"sql leer bei query", `UPDATE interaction SET sql = '' WHERE session_id = 2`, ""},
		{"sql bei extended", `UPDATE interaction SET sql = 'x' WHERE type = 'extended'`, beschaedigt},
		{"sql leer bei extended", `UPDATE interaction SET sql = '' WHERE type = 'extended'`, beschaedigt},
		// message
		{"message.position TEXT", `UPDATE message SET position = 'x' WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"message.position REAL", `UPDATE message SET position = 1.5 WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"message.position NULL", `UPDATE message SET position = NULL WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"Position Lücke", `UPDATE message SET position = 3 WHERE session_id = 2 AND position = 2`, beschaedigt},
		{"Position beginnt bei 2", `UPDATE message SET position = position + 1 WHERE session_id = 2`, beschaedigt},
		{"Position doppelt", `INSERT INTO message SELECT * FROM message WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"message.session_id TEXT", `UPDATE message SET session_id = 'x' WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"message.sequence TEXT", `UPDATE message SET sequence = 'x' WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"Nachricht ohne Interaktion", `UPDATE message SET sequence = 9 WHERE session_id = 2`, beschaedigt},
		{"zusätzliche Nachricht ohne Interaktion", `INSERT INTO message SELECT session_id, 9, position, direction, group_no, kind, fields FROM message WHERE session_id = 2`, beschaedigt},
		{"zusätzliche Nachricht ohne Session", `INSERT INTO message SELECT 9, sequence, position, direction, group_no, kind, fields FROM message WHERE session_id = 2`, beschaedigt},
		{"Nachricht ohne Session", `UPDATE message SET session_id = 9 WHERE session_id = 2`, beschaedigt},
		{"direction NULL", `UPDATE message SET direction = NULL WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"direction unbekannt", `UPDATE message SET direction = 'x' WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"direction Großschreibung", `UPDATE message SET direction = 'Server' WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"direction BLOB", `UPDATE message SET direction = x'736572766572' WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"direction unbekannt bei extended", `UPDATE message SET direction = 'x' WHERE ` + ausNachricht + ` AND position = 4`, beschaedigt},
		{"Client-Nachricht bei query", `UPDATE message SET direction = 'client' WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"group_no bei query", `UPDATE message SET group_no = 1 WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"group_no NULL bei extended", `UPDATE message SET group_no = NULL WHERE ` + ausNachricht + ` AND position = 1`, beschaedigt},
		{"group_no NULL mitten in der Gruppe", `UPDATE message SET group_no = NULL WHERE ` + ausNachricht + ` AND position = 2`, beschaedigt},
		{"group_no TEXT", `UPDATE message SET group_no = 'x' WHERE ` + ausNachricht + ` AND position = 1`, beschaedigt},
		{"group_no REAL", `UPDATE message SET group_no = 1.5 WHERE ` + ausNachricht + ` AND position = 1`, beschaedigt},
		{"group_no 0", `UPDATE message SET group_no = 0 WHERE ` + ausNachricht + ` AND group_no = 1`, beschaedigt},
		{"group_no beginnt bei 2", `UPDATE message SET group_no = group_no + 1 WHERE ` + ausNachricht, beschaedigt},
		{"group_no Lücke", `UPDATE message SET group_no = 4 WHERE ` + ausNachricht + ` AND group_no = 3`, beschaedigt},
		{"group_no zurück", `UPDATE message SET group_no = 1 WHERE ` + ausNachricht + ` AND position = 9`, beschaedigt},
		{"Client nach Server in der Gruppe", `UPDATE message SET position = 4 WHERE ` + ausNachricht + ` AND position = 3; UPDATE message SET position = 3 WHERE ` + ausNachricht + ` AND position = 4 AND kind = 'parse_complete'`, beschaedigt},
		{"message.kind NULL", `UPDATE message SET kind = NULL WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"message.kind BLOB", `UPDATE message SET kind = x'41' WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"Server-Art unbekannt", `UPDATE message SET kind = 'bogus' WHERE session_id = 2 AND position = 1`, beschaedigt},
		{"Client-Art unbekannt", `UPDATE message SET kind = 'bogus' WHERE ` + ausNachricht + ` AND position = 1`, beschaedigt},
		{"Client-Art als Server", `UPDATE message SET kind = 'parse' WHERE session_id = 2 AND position = 1`, beschaedigt},
	})
}

// Abdeckung: LH-FA-22/Negative, LH-QA-06/Messung — fields und die JSON-Spalten von session: kein TEXT, kein gültiger
// JSON-Text, kein gültiges UTF-8, kein Objekt, ein Wert in startup, der kein String ist, ein
// doppelter Schlüssel (auch erst nach der Umkehr des Fluchtzeichens), der Schlüssel type in
// fields und ein Schlüssel, den die Nachricht nicht kennt, sind PGR-E3003; {} in startup und
// server_parameters liest der Leser wie das fehlende Feld.
func TestLesenJSONVerletzt(t *testing.T) {
	setzeFelder := func(wert, wo string) string { return `UPDATE message SET fields = ` + wert + ` WHERE ` + wo }
	tag := `session_id = 2 AND position = 1`
	pruefeFaelle(t, []fall{
		{"fields NULL", setzeFelder(`NULL`, tag), beschaedigt},
		{"fields BLOB", setzeFelder(`x'7b7d'`, tag), beschaedigt},
		{"fields kein JSON", setzeFelder(`'x'`, tag), beschaedigt},
		{"fields leer", setzeFelder(`''`, tag), beschaedigt},
		{"fields offen", setzeFelder(`'{'`, tag), beschaedigt},
		{"fields danach Müll", setzeFelder(`'{"tag":"x"} x'`, tag), beschaedigt},
		{"fields Liste", setzeFelder(`'[]'`, tag), beschaedigt},
		{"fields String", setzeFelder(`'"s"'`, tag), beschaedigt},
		{"fields null", setzeFelder(`'null'`, tag), beschaedigt},
		{"fields Zahl", setzeFelder(`'1'`, tag), beschaedigt},
		{"fields Wahrheitswert", setzeFelder(`'true'`, tag), beschaedigt},
		{"fields ohne UTF-8", setzeFelder(`CAST(x'7b22746167223a22e9227d' AS TEXT)`, tag), beschaedigt},
		{"fields doppelter Schlüssel", setzeFelder(`'{"tag":"a","tag":"b"}'`, tag), beschaedigt},
		{"fields doppelt nach Umkehr", setzeFelder(`'{"tag":"a","\\tag":"b"}'`, tag), beschaedigt},
		{"fields doppelt nach base64", setzeFelder(`'{"tag":"a","base64:dGFn":"b"}'`, tag), beschaedigt},
		{"fields mit type", setzeFelder(`'{"type":"command_complete","tag":"a"}'`, tag), beschaedigt},
		{"fields unbekannter Schlüssel", setzeFelder(`'{"tag":"a","bogus":1}'`, tag), beschaedigt},
		{"fields Bruch", setzeFelder(`'{"tag":1.5}'`, tag), beschaedigt},
		{"fields gültig", setzeFelder(`'{"tag":"a"}'`, tag), ""},
		{"startup NULL", `UPDATE session SET startup = NULL WHERE id = 2`, beschaedigt},
		{"startup BLOB", `UPDATE session SET startup = x'7b7d' WHERE id = 2`, beschaedigt},
		{"startup kein JSON", `UPDATE session SET startup = '{' WHERE id = 2`, beschaedigt},
		{"startup null", `UPDATE session SET startup = 'null' WHERE id = 2`, beschaedigt},
		{"startup Liste", `UPDATE session SET startup = '[]' WHERE id = 2`, beschaedigt},
		{"startup String", `UPDATE session SET startup = '"x"' WHERE id = 2`, beschaedigt},
		{"startup Zahl als Wert", `UPDATE session SET startup = '{"a":1}' WHERE id = 2`, beschaedigt},
		{"startup null als Wert", `UPDATE session SET startup = '{"a":null}' WHERE id = 2`, beschaedigt},
		{"startup Wahrheitswert als Wert", `UPDATE session SET startup = '{"a":true}' WHERE id = 2`, beschaedigt},
		{"startup Objekt als Wert", `UPDATE session SET startup = '{"a":{}}' WHERE id = 2`, beschaedigt},
		{"startup Liste als Wert", `UPDATE session SET startup = '{"a":[]}' WHERE id = 2`, beschaedigt},
		{"startup doppelter Schlüssel", `UPDATE session SET startup = '{"a":"1","a":"2"}' WHERE id = 2`, beschaedigt},
		{"startup doppelt nach Umkehr", `UPDATE session SET startup = '{"a":"1","\\a":"2"}' WHERE id = 2`, beschaedigt},
		{"startup {}", `UPDATE session SET startup = '{}' WHERE id = 2`, ""},
		{"server_parameters NULL", `UPDATE session SET server_parameters = NULL WHERE id = 2`, beschaedigt},
		{"server_parameters Liste", `UPDATE session SET server_parameters = '[]' WHERE id = 2`, beschaedigt},
		{"server_parameters Zahl als Wert", `UPDATE session SET server_parameters = '{"a":1}' WHERE id = 2`, beschaedigt},
		{"server_parameters {}", `UPDATE session SET server_parameters = '{}' WHERE id = 2`, ""},
	})
}

// Abdeckung: LH-FA-22/Boundary — {} in startup und server_parameters liest der Leser wie das fehlende Feld
// in yaml: Das Modell trägt keine Einträge (nil, nicht eine leere Abbildung).
func TestLesenLeereObjekteSindKeineEintraege(t *testing.T) {
	path := sqliteDatei(t, lesefixture())
	rohSchreibe(t, path, `UPDATE session SET startup = '{}', server_parameters = '{}' WHERE id = 1`)
	got := ladeOK(t, path)
	if got.Sessions[0].Startup != nil || got.Sessions[0].ServerParameters != nil {
		t.Errorf("Startup %#v, Server-Parameter %#v", got.Sessions[0].Startup, got.Sessions[0].ServerParameters)
	}
}

// Abdeckung: LH-FA-22/Negative, LH-QA-06/Messung — meta: ohne format oder version, mit anderer Formatkennung, mit
// einem unbekannten Schlüssel, mit einem Wert von empty_sessions außer true oder mit einer
// Version, die keine kanonische Zahl ist, ist die Datei beschädigt (PGR-E3003); jede kanonische
// Zahl außer 1, auch 0, eine negative und eine zu große für 64 Bit, ist PGR-E3002; geprüft
// wird in der Reihenfolge Formatkennung, Version, Inhalt.
func TestLesenMeta(t *testing.T) {
	version := func(v string) string { return `UPDATE meta SET value = '` + v + `' WHERE key = 'version'` }
	kaputt := `UPDATE message SET fields = 'x'; DROP TABLE interaction`
	faelle := []fall{
		{"ohne format", `DELETE FROM meta WHERE key = 'format'`, beschaedigt},
		{"format NULL", `UPDATE meta SET value = NULL WHERE key = 'format'`, beschaedigt},
		{"format BLOB", `UPDATE meta SET value = x'70' WHERE key = 'format'`, beschaedigt},
		{"format anders", `UPDATE meta SET value = 'x' WHERE key = 'format'`, beschaedigt},
		{"format Großschreibung", `UPDATE meta SET value = 'PGWIRE-RECORDER' WHERE key = 'format'`, beschaedigt},
		{"ohne version", `DELETE FROM meta WHERE key = 'version'`, beschaedigt},
		{"version BLOB", `UPDATE meta SET value = x'31' WHERE key = 'version'`, beschaedigt},
		{"unbekannter Schlüssel", `INSERT INTO meta VALUES ('foo', 'bar')`, beschaedigt},
		{"Schlüssel doppelt", `INSERT INTO meta VALUES ('version', '1')`, beschaedigt},
		{"empty_sessions false", `INSERT INTO meta VALUES ('empty_sessions', 'false')`, beschaedigt},
		{"empty_sessions 1", `INSERT INTO meta VALUES ('empty_sessions', '1')`, beschaedigt},
		{"empty_sessions leer", `INSERT INTO meta VALUES ('empty_sessions', '')`, beschaedigt},
		{"empty_sessions TRUE", `INSERT INTO meta VALUES ('empty_sessions', 'TRUE')`, beschaedigt},
		{"empty_sessions true", `INSERT INTO meta VALUES ('empty_sessions', 'true')`, ""},
		{"Formatkennung vor Version", `UPDATE meta SET value = 'x' WHERE key = 'format'; ` + version("2"), beschaedigt},
		{"Version vor Inhalt", version("2") + `; ` + kaputt, unbekannteV},
		{"Version vor fehlender Tabelle", version("2") + `; DROP TABLE session`, unbekannteV},
		{"Version vor unbekanntem Schlüssel", version("2") + `; INSERT INTO meta VALUES ('foo', 'bar')`, unbekannteV},
		{"unlesbare Version vor Inhalt", version("abc") + `; ` + kaputt, beschaedigt},
	}
	for _, v := range []string{"abc", "01", "+1", " 1", "1 ", "-0", "", "1.0", "0x1", "1e0", "00", "-01", "٣"} {
		faelle = append(faelle, fall{fmt.Sprintf("version %q", v), version(v), beschaedigt})
	}
	for _, v := range []string{"0", "2", "-1", "10", "99999999999999999999", "-99999999999999999999", "9223372036854775808"} {
		faelle = append(faelle, fall{fmt.Sprintf("version %q", v), version(v), unbekannteV})
	}
	faelle = append(faelle, fall{"version 1", version("1"), ""})
	pruefeFaelle(t, faelle)
}

// Abdeckung: LH-FA-22/Negative — die Version ist Sache von meta: PRAGMA user_version wird weder gelesen noch
// erzwungen, eine Datei mit user_version 1, 2 oder -1 liest wie eine ohne.
func TestLesenUserVersionIrrelevant(t *testing.T) {
	for _, v := range []string{"1", "2", "-1", "4294967295"} {
		path := sqliteDatei(t, lesefixture())
		rohSchreibe(t, path, `PRAGMA user_version = `+v)
		if got := ladeOK(t, path); !reflect.DeepEqual(got, lesefixture()) {
			t.Errorf("user_version %s: %#v", v, got)
		}
	}
}

// Abdeckung: LH-FA-22/Happy — ein String ohne gültiges UTF-8 steht in der Datei als base64: und kommt als
// dieselben Bytes zurück, als Wert wie als Schlüssel in Startup-Parametern, Server-Parametern
// und Nachrichtenfeldern (ein Fehlerfeld, ein Hinweis, ein Parameterwert); ein gültiger String,
// der mit \ oder base64: beginnt, kommt unverändert zurück; interaction.sql kommt unverändert
// zurück, auch ohne gültiges UTF-8; dieselbe Aufzeichnung aus yaml und sqlite ist gleich.
func TestLesenStringOhneUTF8(t *testing.T) {
	lang := strings.Repeat("\xe9", 200)
	rec := lesefixture()
	rec.Sessions[0].Startup = map[string]string{
		"application_name": "caf\xe9", "na\xe9me": "v", "lang": lang, "leer": "",
		"fluchtzeichen": `\x`, "nur_fluchtzeichen": `\`, "doppelt": `\\x`, "praefix": "base64:Zm9v", "nur_praefix": "base64:", "praefix_kaputt": "base64:!!",
		`\schluessel`: "1", "base64:schluessel": "2", "nul": "a\x00b", "unicode": "日本語 \U0001F600 \u2028 \ufeff", "steuer": "\x07\x1b\n\t",
	}
	rec.Sessions[0].ServerParameters = map[string]string{"server\xe9": "17.0", "wert": "\xff\xfe"}
	rec.Sessions[0].Interactions[0].Request.SQL = "SELECT '\xe9\xff', 1"
	rec.Sessions[0].Interactions[0].Responses[3].Fields = map[string]string{"M": "caf\xe9", "S\xe9": "NOTICE", "C": "base64:x"}
	rec.Sessions[0].Interactions[0].Responses[2].Tag = "\xff"
	rec.Sessions[0].Interactions[1].Groups[2].Client[0].Params[1] = model.Value{Bytes: []byte("\xe9")}
	rec.Sessions[0].Interactions[1].Groups[0].Client[0].Statement = "\xe9"
	path := sqliteDatei(t, rec)
	if got := ladeOK(t, path); !reflect.DeepEqual(got, rec) {
		t.Errorf("aus sqlite:\n%#v\ngeschrieben:\n%#v", got, rec)
	}
}

// Abdeckung: LH-FA-22/Negative — die Umkehr der Umschreibung folgt SPEC-043: Beginnt ein String mit \,
// entfällt dieses eine Zeichen; beginnt er sonst mit base64:, ist der Rest Base64 im
// Standard-Alphabet mit Auffüllung (sonst PGR-E3003); jeder andere String gilt wörtlich.
func TestLesenUmkehrDerUmschreibung(t *testing.T) {
	path := sqliteDatei(t, sessionsMit(1))
	lockere(t, path)
	for _, c := range []struct {
		json string
		want map[string]string
		err  bool
	}{
		{`{"a":"base64:Y2Fm6Q=="}`, map[string]string{"a": "caf\xe9"}, false},
		{`{"a":"base64:Zm9v"}`, map[string]string{"a": "foo"}, false},
		{`{"a":"base64:"}`, map[string]string{"a": ""}, false},
		{`{"base64:Y2Fm6Q==":"v"}`, map[string]string{"caf\xe9": "v"}, false},
		{`{"a":"\\base64:x"}`, map[string]string{"a": "base64:x"}, false},
		{`{"a":"\\\\x"}`, map[string]string{"a": `\x`}, false},
		{`{"a":"\\x"}`, map[string]string{"a": "x"}, false},
		{`{"a":"\\"}`, map[string]string{"a": ""}, false},
		{`{"a":"x\\y"}`, map[string]string{"a": `x\y`}, false},
		{`{"a":"xbase64:"}`, map[string]string{"a": "xbase64:"}, false},
		{`{"a":"BASE64:Zm9v"}`, map[string]string{"a": "BASE64:Zm9v"}, false},
		{`{"a":"base64:!!"}`, nil, true},
		{`{"a":"base64:Y2Fm6Q"}`, nil, true},
		{`{"a":"base64: Y2Fm6Q=="}`, nil, true},
		{`{"a":"base64:Y2Fm6Q==="}`, nil, true},
		{`{"a":"base64:Y2Fm-Q=="}`, nil, true},
		{`{"base64:!":"v"}`, nil, true},
	} {
		p := kopiere(t, path)
		rohSchreibe(t, p, `UPDATE session SET startup = '`+strings.ReplaceAll(c.json, `'`, `''`)+`'`)
		got, err := laden(p)
		switch {
		case c.err && codeVon(err) != beschaedigt:
			t.Errorf("%s: %v", c.json, err)
		case !c.err && (err != nil || !reflect.DeepEqual(got.Sessions[0].Startup, c.want)):
			t.Errorf("%s: %v, %q", c.json, err, got.Sessions[0].Startup)
		}
	}
}

// Abdeckung: LH-FA-22/Negative, LH-QA-06/Messung — eine Zahl in fields ist eine ganze Zahl in kanonischer
// Dezimalschreibweise im Bereich der 64-Bit-Zahlen mit Vorzeichen, und die Bibliothek führt sie
// nicht über float64: 9007199254740993 bleibt 9007199254740993; 1.0, 1e2, -0, 01, +1, ein Bruch,
// ein Exponent und 9223372036854775808 sind abgelehnt.
func TestLesenJSONZahlen(t *testing.T) {
	for text, want := range map[string]string{
		"0": "0", "1": "1", "-1": "-1", "42": "42",
		"9007199254740993": "9007199254740993", "9007199254740992": "9007199254740992",
		"9223372036854775807": "9223372036854775807", "-9223372036854775808": "-9223372036854775808",
		"[1]": "", "{}": "",
	} {
		tag, wert, err := recording.JSONKnoten(text)
		if err != nil || (want != "" && (tag != "!!int" || wert != want)) {
			t.Errorf("%s: %s %q, %v", text, tag, wert, err)
		}
	}
	for _, text := range []string{"1.0", "1e2", "1E2", "-0", "01", "+1", "1.5", "-1.5", "1e0", "0.0", "-0.0", "1e-2", "9223372036854775808", "-9223372036854775809", "99999999999999999999", "00", "-", "--1", ".5", "1.", "0x10", "NaN", "Infinity"} {
		if _, _, err := recording.JSONKnoten(text); err == nil {
			t.Errorf("%s: abgelehnt erwartet", text)
		}
	}
	// In einer Nachricht: ein Wert im Bereich von uint32 liest sich genau.
	path := sqliteDatei(t, lesefixture())
	zeile := func(oid string) string {
		return `UPDATE message SET fields = '{"fields":[{"name":"x","table_oid":` + oid + `,"column_number":0,"type_oid":23,"type_size":4,"type_modifier":-1,"format":0}]}' WHERE session_id = 1 AND sequence = 1 AND position = 1`
	}
	for oid, want := range map[string]uint32{"0": 0, "4294967295": 4294967295} {
		p := kopiere(t, path)
		rohSchreibe(t, p, zeile(oid))
		if got := ladeOK(t, p); got.Sessions[0].Interactions[0].Responses[0].Columns[0].TableOID != want {
			t.Errorf("table_oid %s: %d", oid, got.Sessions[0].Interactions[0].Responses[0].Columns[0].TableOID)
		}
	}
	for _, oid := range []string{"1.0", "1e2", "-0", "01", "+1", "9007199254740993", "4294967296", "-1", "1.5", "9223372036854775808", `"7"`, "true"} {
		p := kopiere(t, path)
		rohSchreibe(t, p, zeile(oid))
		if _, err := laden(p); codeVon(err) != beschaedigt {
			t.Errorf("table_oid %s: %v", oid, err)
		}
	}
}

// Abdeckung: LH-FA-22/Negative — ein JSON-Text mit sehr tiefer Verschachtelung ist abgelehnt
// (PGR-E3003), ohne den Prozess zu gefährden.
func TestLesenJSONTief(t *testing.T) {
	tief := strings.Repeat("[", 100000) + strings.Repeat("]", 100000)
	if _, _, err := recording.JSONKnoten(tief); err == nil {
		t.Error("tief verschachtelter Text gelesen")
	}
	path := sqliteDatei(t, lesefixture())
	rohSchreibe(t, path, `UPDATE message SET fields = '{"tag":`+tief+`}' WHERE session_id = 2 AND position = 1`)
	if _, err := laden(path); codeVon(err) != beschaedigt {
		t.Errorf("tief verschachtelt: %v", err)
	}
}
