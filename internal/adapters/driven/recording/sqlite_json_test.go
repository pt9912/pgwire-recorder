package recording_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// umgekehrt wendet die Regel des Lesers aus SPEC-043 auf s an: Beginnt s mit \,
// entfällt dieses eine Zeichen und der Rest gilt wörtlich; beginnt s sonst mit
// base64:, sind die Bytes des Rests dekodiert; jeder andere String gilt wörtlich.
func umgekehrt(t *testing.T, s string) string {
	t.Helper()
	switch {
	case strings.HasPrefix(s, `\`):
		return s[1:]
	case strings.HasPrefix(s, "base64:"):
		b, err := base64.StdEncoding.DecodeString(s[len("base64:"):])
		if err != nil {
			t.Fatalf("%q: %v", s, err)
		}
		return string(b)
	}
	return s
}

// startupText schreibt eine Aufzeichnung mit einer Session, deren Startup-Parameter
// startup und deren Server-Parameter server sind, und liefert die beiden Spalten.
func startupText(t *testing.T, startup, server map[string]string) (string, string) {
	t.Helper()
	rec := sessionsMit(1)
	rec.Sessions[0].Startup, rec.Sessions[0].ServerParameters = startup, server
	_, path := neueDatei(t, rec)
	got := zeilen(t, roh(t, path), `SELECT startup, server_parameters FROM session`)
	teile := strings.SplitN(got[0], "|", 2)
	return teile[0], teile[1]
}

// Abdeckung: LH-FA-22/Happy — ein String ohne gültiges UTF-8 steht in den JSON-Spalten
// als base64: und seine Bytes in Base64 (Standard-Alphabet mit Auffüllung): als Wert und
// als Schlüssel der Startup-Parameter, der Server-Parameter und der Felder einer
// Nachricht (ein Fehlerfeld, ein Schlüssel eines Hinweises); ein gültiger String bleibt
// unverändert.
func TestSQLiteStringOhneUTF8(t *testing.T) {
	lang := strings.Repeat("\xe9", 200)
	startup, server := startupText(t,
		map[string]string{"application_name": "caf\xe9", "client_encoding": "LATIN1", "na\xe9me": "v", "lang": lang},
		map[string]string{"server\xe9": "17.0", "wert": "\xff\xfe"})
	if want := `{"application_name":"base64:Y2Fm6Q==","client_encoding":"LATIN1","lang":"base64:` + base64.StdEncoding.EncodeToString([]byte(lang)) + `","base64:bmHpbWU=":"v"}`; startup != want {
		t.Errorf("Startup:\n got %s\nwant %s", startup, want)
	}
	if want := `{"base64:c2VydmVy6Q==":"17.0","wert":"base64://4="}`; server != want {
		t.Errorf("Server-Parameter:\n got %s\nwant %s", server, want)
	}

	rec := model.NewRecording()
	rec.Sessions = []model.Session{{ID: 1, Interactions: []model.Interaction{{
		Sequence: 1, Request: model.Request{Type: model.RequestQuery, SQL: "x"},
		Responses: []model.Response{
			{Type: model.ResponseErrorResponse, Fields: map[string]string{"M": "caf\xe9", "S": "ERROR"}},
			{Type: model.ResponseNoticeResponse, Fields: map[string]string{"na\xe9me": "x"}},
			{Type: model.ResponseReadyForQuery, TxStatus: "I"},
		},
	}}}}
	_, path := neueDatei(t, rec)
	want := []string{
		`{"notice":{"M":"base64:Y2Fm6Q==","S":"ERROR"}}`,
		`{"notice":{"base64:bmHpbWU=":"x"}}`,
		`{"tx_status":"I"}`,
	}
	if got := zeilen(t, roh(t, path), `SELECT fields FROM message ORDER BY position`); !reflect.DeepEqual(got, want) {
		t.Errorf("Felder:\n got %q\nwant %q", got, want)
	}
}

// Abdeckung: LH-FA-22/Happy — ein gültiger String, der mit \ oder mit base64: beginnt,
// steht mit einem weiteren \ davor, als Wert wie als Schlüssel (Startup, Server-Parameter
// und Felder); ein String, der \ oder base64: nur im Inneren trägt, bleibt
// unverändert, ebenso der Wert eines Schlüssels base64 einer Nachricht.
func TestSQLiteFluchtzeichen(t *testing.T) {
	startup, server := startupText(t,
		map[string]string{`\a`: `\\b`, "base64:k": "base64:AAAA", "x": `a\b`, "y": "ab base64:", "z": "base64"},
		map[string]string{`\`: "base64:"})
	if want := `{"\\\\a":"\\\\\\b","\\base64:k":"\\base64:AAAA","x":"a\\b","y":"ab base64:","z":"base64"}`; startup != want {
		t.Errorf("Startup:\n got %s\nwant %s", startup, want)
	}
	if want := `{"\\\\":"\\base64:"}`; server != want {
		t.Errorf("Server-Parameter:\n got %s\nwant %s", server, want)
	}

	rec := model.NewRecording()
	rec.Sessions = []model.Session{{ID: 1, Interactions: []model.Interaction{{
		Sequence: 1, Request: model.Request{Type: model.RequestQuery, SQL: `\SELECT`},
		Responses: []model.Response{
			{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte(`\x`)}, {Bytes: []byte("base64:")}, {Bytes: []byte{0x00, 0xff}}, {Bytes: []byte("ok")}}},
			{Type: model.ResponseCommandComplete, Tag: `\tag`},
			{Type: model.ResponseNoticeResponse, Fields: map[string]string{`\k`: "base64:v"}},
			{Type: model.ResponseReadyForQuery, TxStatus: "I"},
		},
	}}}}
	_, path := neueDatei(t, rec)
	db := roh(t, path)
	want := []string{
		`{"values":[{"text":"\\\\x"},{"text":"\\base64:"},{"base64":"AP8="},{"text":"ok"}]}`,
		`{"tag":"\\\\tag"}`,
		`{"notice":{"\\\\k":"\\base64:v"}}`,
		`{"tx_status":"I"}`,
	}
	if got := zeilen(t, db, `SELECT fields FROM message ORDER BY position`); !reflect.DeepEqual(got, want) {
		t.Errorf("Felder:\n got %q\nwant %q", got, want)
	}
	if got := zeilen(t, db, `SELECT sql FROM interaction`); !reflect.DeepEqual(got, []string{`\SELECT`}) {
		t.Errorf("SQL-Text: %q", got)
	}
}

// Abdeckung: LH-FA-22/Boundary — jeder String kommt durch die Umschreibung unverändert
// zurück: Die Regel des Lesers (\ entfernen, sonst base64: dekodieren, sonst wörtlich),
// auf Schlüssel und Wert angewandt, ergibt für Strings mit ungültigem UTF-8, mit
// Anführungszeichen, Steuerzeichen, Zeilenenden, <, >, & und U+2028 dieselben Bytes; die
// Spalte ist gültiges JSON.
func TestSQLiteUmschreibungKehrtZurueck(t *testing.T) {
	korpus := []string{
		"", "a", `\`, `\\`, `\base64:`, "base64:", "base64:AAAA", "base64:!!!", "base64", "Base64:x",
		"\xe9", "caf\xe9", "\xff\xfe\x00", "\xc3", "\xed\xa0\x80", "\xf4\x90\x80\x80",
		`"`, `a"b`, "a\nb", "a\r\nb", "a\tb", "\x00", "\x1f", "\x7f", "<a>&", "  ", "\u0085", "ä€𝄞",
		" führend", "nachfolgend ", "null", "true", "123", "~", "- x", "? x", "a: b", "# c", "{}", "[]", `'`, "%", "!!binary x",
	}
	werte := map[string]string{}
	for i, k := range korpus {
		werte[k+"/"+string(rune('a'+i%26))] = korpus[(i+7)%len(korpus)]
	}
	for _, k := range korpus {
		werte[k] = k
	}
	startup, _ := startupText(t, werte, nil)
	var roh map[string]string
	if err := json.Unmarshal([]byte(startup), &roh); err != nil {
		t.Fatalf("Startup ist kein JSON-Objekt aus Strings: %v\n%s", err, startup)
	}
	got := map[string]string{}
	for k, v := range roh {
		got[umgekehrt(t, k)] = umgekehrt(t, v)
	}
	if !reflect.DeepEqual(got, werte) {
		for k, v := range werte {
			if got[k] != v {
				t.Errorf("%q: %q statt %q", k, got[k], v)
			}
		}
		t.Errorf("%d statt %d Einträge", len(got), len(werte))
	}
}

// Abdeckung: LH-FA-22/Happy — interaction.sql trägt die Bytes des SQL-Textes
// unverändert, auch ohne gültiges UTF-8 (das Byte 0xE9) und mit führendem \ oder
// base64: (keine Umschreibung); die Spalte hat den Typ text.
func TestSQLiteSQLUnveraendert(t *testing.T) {
	rec := sessionsMit(1)
	texte := []string{"SELECT 'caf\xe9'", `\SELECT`, "base64:SELECT", "SELECT '\xff'"}
	rec.Sessions[0].Interactions = nil
	for i, sqlText := range texte {
		rec.Sessions[0].Interactions = append(rec.Sessions[0].Interactions, model.Interaction{
			Sequence: i + 1, Request: model.Request{Type: model.RequestQuery, SQL: sqlText},
			Responses: []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}},
		})
	}
	_, path := neueDatei(t, rec)
	db := roh(t, path)
	var want []string
	for i, sqlText := range texte {
		want = append(want, strings.Join([]string{string(rune('1' + i)), "text", strings.ToLower(hexVon(sqlText))}, "|"))
	}
	got := zeilen(t, db, `SELECT sequence, typeof(sql), hex(CAST(sql AS BLOB)) FROM interaction ORDER BY sequence`)
	for i := range got {
		got[i] = strings.ToLower(got[i])
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Bytes des SQL-Textes:\n got %q\nwant %q", got, want)
	}
}

func hexVon(s string) string {
	const ziffern = "0123456789ABCDEF"
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		b.WriteByte(ziffern[s[i]>>4])
		b.WriteByte(ziffern[s[i]&15])
	}
	return b.String()
}

// Abdeckung: LH-FA-22/Boundary — den Stand eines Laufs (welche Sessions die Datei trägt) führt
// der Adapter, nicht das Paket: Zwei Adapter ergänzen zwei Dateien unabhängig voneinander,
// auch im selben Verzeichnis.
func TestSQLiteStandJeAdapter(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	a, b := filepath.Join(dir, "a.sqlite"), filepath.Join(dir, "b.sqlite")
	sa, sb := &recording.SQLite{}, &recording.SQLite{}
	for i := 1; i <= 2; i++ {
		if err := sa.Write(ctx, a, sessionsMit(i)); err != nil {
			t.Fatal(err)
		}
	}
	if err := sb.Write(ctx, b, sessionsMit(1)); err != nil {
		t.Fatal(err)
	}
	if sessionIDs(t, a) != "1,2" || sessionIDs(t, b) != "1" {
		t.Fatalf("a: %s, b: %s", sessionIDs(t, a), sessionIDs(t, b))
	}
}

// Abdeckung: LH-FA-22/Boundary — die Umwandlung in JSON kennt Abbildung, Liste, String,
// Zahl (auch negativ), Wahrheitswert und null, schreibt die Schlüssel in der Reihenfolge
// der YAML-Darstellung, lässt den genannten Schlüssel der obersten Abbildung aus und
// maskiert <, > und & nicht; jeder andere Wert ist ein Fehler.
func TestSQLiteJSONTypen(t *testing.T) {
	got, err := recording.JSONText(map[string]any{
		"s": "<a&b>", "i": -5, "g": 0, "t": true, "f": false, "n": nil, "l": []any{1, "x", nil}, "m": map[string]any{}, "typ": "weg",
	}, "typ")
	want := `{"f":false,"g":0,"i":-5,"l":[1,"x",null],"m":{},"n":null,"s":"<a&b>","t":true}`
	if err != nil || got != want {
		t.Errorf("JSON:\n got %s (%v)\nwant %s", got, err, want)
	}
	if got, err := recording.JSONText(map[string]any{"typ": map[string]any{"typ": 1}, "a": map[string]any{"typ": 2}}, "typ"); err != nil || got != `{"a":{"typ":2}}` {
		t.Errorf("Schlüssel ohne nur oben: %s (%v)", got, err)
	}
	if got, err := recording.JSONText(map[any]any{nil: 1}, ""); err != nil || got != `{"null":1}` {
		t.Errorf("Schlüssel null: %s (%v)", got, err)
	}
	for name, wert := range map[string]any{"Fließkommazahl": 1.5, "Zeit": time.Unix(0, 0).UTC()} {
		if got, err := recording.JSONText(map[string]any{"x": wert}, ""); err == nil {
			t.Errorf("%s: erwartet einen Fehler, erhalten %s", name, got)
		}
	}
}
