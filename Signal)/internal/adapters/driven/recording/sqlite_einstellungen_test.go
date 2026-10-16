package recording_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-22/Boundary — die Datenbank wird mit dem Rollback-Journal
// (journal_mode=DELETE; die Kopfzeile der Datei trägt Schreib- und Leseversion 1,
// nicht die 2 des WAL), mit eingeschalteten Fremdschlüsseln, synchronous=FULL und
// einer Sperrfrist von 5000 ms geöffnet; nach dem Schreiben bleibt neben der Datei
// kein Journal und keine WAL-Datei.
func TestSQLiteEinstellungen(t *testing.T) {
	ctx := context.Background()
	s, path := neueDatei(t, sessionsMit(1))
	if err := s.Write(ctx, path, sessionsMit(2)); err != nil {
		t.Fatal(err)
	}
	db, err := recording.Oeffne(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	for pragma, want := range map[string]string{
		"journal_mode": "delete", "foreign_keys": "1", "synchronous": "2", "busy_timeout": "5000",
	} {
		var got string
		if err := db.QueryRowContext(ctx, "PRAGMA "+pragma).Scan(&got); err != nil || got != want {
			t.Errorf("PRAGMA %s: %q, %v, erwartet %q", pragma, got, err, want)
		}
	}
	kopf, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if kopf[18] != 1 || kopf[19] != 1 {
		t.Errorf("Schreib- und Leseversion der Datei: %d, %d, erwartet 1, 1", kopf[18], kopf[19])
	}
	if es := eintraege(t, filepath.Dir(path)); !reflect.DeepEqual(es, []string{"rec.sqlite"}) {
		t.Errorf("im Verzeichnis liegt %v", es)
	}
}

// Abdeckung: LH-FA-22/Boundary — die Fremdschlüssel sind eingeschaltet: eine
// Interaktion ohne Session lehnt die geöffnete Datenbank ab.
func TestSQLiteFremdschluessel(t *testing.T) {
	ctx := context.Background()
	_, path := neueDatei(t, sessionsMit(1))
	db, err := recording.Oeffne(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	_, err = db.ExecContext(ctx, `INSERT INTO interaction (session_id, sequence, type) VALUES (99, 1, 'query')`)
	if err == nil || !strings.Contains(err.Error(), "FOREIGN KEY") {
		t.Fatalf("Interaktion ohne Session: %v", err)
	}
}

// Abdeckung: LH-FA-22/Boundary — jede Transaktion beginnt mit BEGIN IMMEDIATE: Solange
// eine geöffnete Transaktion läuft, scheitert eine zweite Verbindung schon beim
// Beginn, noch bevor sie schreibt.
func TestSQLiteTransaktionImmediate(t *testing.T) {
	ctx := context.Background()
	_, path := neueDatei(t, sessionsMit(1))
	db, err := recording.Oeffne(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	andere, err := sql.Open("sqlite", "file:"+path+"?mode=rw&_pragma=busy_timeout(0)")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = andere.Close() }()
	if _, err := andere.ExecContext(ctx, `BEGIN IMMEDIATE`); err == nil || !strings.Contains(err.Error(), "locked") {
		t.Fatalf("BEGIN IMMEDIATE neben einer laufenden Transaktion: %v", err)
	}
}

// Abdeckung: LH-FA-22/Negative — ist die Datei länger als fünf Sekunden gesperrt,
// scheitert das Ergänzen nach etwa fünf Sekunden mit PGR-E3001, die Transaktion ist
// zurückgerollt, und der nächste Schreibvorgang ergänzt, was fehlt.
func TestSQLiteSperrfrist(t *testing.T) {
	ctx := context.Background()
	s, path := neueDatei(t, sessionsMit(1))
	sperre, err := sql.Open("sqlite", "file:"+path+"?mode=rw&_txlock=immediate&_pragma=busy_timeout(0)")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sperre.Close() }()
	tx, err := sperre.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	err = s.Write(ctx, path, sessionsMit(2))
	dauer := time.Since(start)
	if code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "locked") {
		t.Fatalf("erwartet %s mit der Ursache der Bibliothek, erhalten %v", model.CodeRecordingIO, err)
	}
	if dauer < 4500*time.Millisecond || dauer > 8*time.Second {
		t.Errorf("gescheitert nach %v, erwartet etwa 5 s", dauer)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	if got := sessionIDs(t, path); got != "1" {
		t.Fatalf("nach dem Fehlschlag Sessions %q", got)
	}
	if err := s.Write(ctx, path, sessionsMit(2)); err != nil {
		t.Fatalf("Nachholen: %v", err)
	}
	if got := sessionIDs(t, path); got != "1,2" {
		t.Fatalf("nach dem Nachholen Sessions %q", got)
	}
}

// Abdeckung: LH-FA-22/Boundary — der Pfad der Zieldatei gilt wörtlich, auch mit ?, #,
// % (auch als %41) und Leerzeichen, absolut wie relativ: Die Datei liegt genau unter
// diesem Pfad, und im Verzeichnis liegt sonst nichts.
func TestSQLiteSonderzeichenImPfad(t *testing.T) {
	ctx := context.Background()
	for _, name := range []string{"a?b#c%d e.sqlite", "%41 %3F.sqlite", "x?mode=ro.sqlite", "#.sqlite", "+&=.sqlite"} {
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
			s := &recording.SQLite{}
			if err := s.Write(ctx, path, sessionsMit(1)); err != nil {
				t.Fatalf("%q relativ=%v: %v", name, relativ, err)
			}
			if es := eintraege(t, dir); !reflect.DeepEqual(es, []string{name}) {
				t.Fatalf("%q relativ=%v: im Verzeichnis liegt %q", name, relativ, es)
			}
			if err := s.Write(ctx, path, sessionsMit(2)); err != nil {
				t.Fatalf("%q relativ=%v, Ergänzen: %v", name, relativ, err)
			}
			db, err := sql.Open("sqlite", "file:"+strings.NewReplacer("%", "%25", "?", "%3F", "#", "%23").Replace(filepath.Join(dir, name))+"?mode=ro")
			if err != nil {
				t.Fatal(err)
			}
			if got := zeilen(t, db, `SELECT id FROM session ORDER BY id`); !reflect.DeepEqual(got, []string{"1", "2"}) {
				t.Fatalf("%q relativ=%v: Sessions %q", name, relativ, got)
			}
			_ = db.Close()
			if err := os.Remove(filepath.Join(dir, name)); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// Abdeckung: LH-FA-22/Negative — jeder Wert gelangt als Parameter in die Datenbank,
// nie als Text im SQL: SQL-Text, Startup-Parameter und Nachrichtenfelder mit
// Anführungszeichen, Semikolon und SQL-Befehlen stehen wörtlich in der Datei, und
// keine Tabelle fehlt.
func TestSQLiteWerteAlsParameter(t *testing.T) {
	boese := `'); DROP TABLE session; --"` + "\x01"
	rec := model.NewRecording()
	rec.Sessions = []model.Session{{
		ID: 1, Startup: map[string]string{"user": boese},
		Interactions: []model.Interaction{{
			Sequence: 1, Request: model.Request{Type: model.RequestQuery, SQL: boese},
			Responses: []model.Response{{Type: model.ResponseCommandComplete, Tag: boese}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}},
		}},
	}}
	_, path := neueDatei(t, rec)
	db := roh(t, path)
	if got := zeilen(t, db, `SELECT sql FROM interaction`); !reflect.DeepEqual(got, []string{boese}) {
		t.Errorf("SQL-Text: %q", got)
	}
	if got := zeilen(t, db, `SELECT json_extract(startup, '$.user'), (SELECT json_extract(fields, '$.tag') FROM message WHERE position = 1) FROM session`); !reflect.DeepEqual(got, []string{boese + "|" + boese}) {
		t.Errorf("Startup und Feld: %q", got)
	}
	if got := zeilen(t, db, `SELECT count(*) FROM sqlite_master WHERE name IN ('meta', 'session', 'interaction', 'message')`); got[0] != "4" {
		t.Errorf("Tabellen: %q", got)
	}
}

// Abdeckung: LH-FA-22/Boundary — ein Abbruchsignal (abgebrochener Kontext) unterbricht
// keine Transaktion: Der Aufruf mit abgebrochenem Kontext legt die Datei an und
// ergänzt die Session vollständig.
func TestSQLiteAbbruchUnterbrichtNicht(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	path := filepath.Join(t.TempDir(), "rec.sqlite")
	s := &recording.SQLite{}
	if err := s.Write(ctx, path, sessionsMit(2)); err != nil {
		t.Fatal(err)
	}
	if got := sessionIDs(t, path); got != "1,2" {
		t.Fatalf("Sessions %q", got)
	}
	if err := s.Write(ctx, path, sessionsMit(3)); err != nil {
		t.Fatal(err)
	}
	if got := sessionIDs(t, path); got != "1,2,3" {
		t.Fatalf("Sessions %q", got)
	}
}

// Abdeckung: LH-FA-22/Negative — scheitert das Schließen der Datenbank, nachdem die
// Transaktion abgeschlossen ist, ist das PGR-E3001 mit dem Fehler als Ursache; die
// Session ist ergänzt, und der nächste Schreibvorgang trägt sie nicht noch einmal ein.
func TestSQLiteSchliessenScheitert(t *testing.T) {
	ctx := context.Background()
	s, path := neueDatei(t, sessionsMit(1))
	fehler := errors.New("Schließen eingespielt gescheitert")
	err := recording.SQLiteSchreibeMit(ctx, s, path, sessionsMit(2), recording.SQLiteEingriffe{SchliesseDB: func(db *sql.DB) error {
		return errors.Join(db.Close(), fehler)
	}})
	ms := model.Meldungen(err)
	if len(ms) != 1 || ms[0].Code != model.CodeRecordingIO || !strings.Contains(ms[0].Text, fehler.Error()) {
		t.Fatalf("Meldungen %+v", ms)
	}
	if got := sessionIDs(t, path); got != "1,2" {
		t.Fatalf("Sessions %q", got)
	}
	if err := s.Write(ctx, path, sessionsMit(2)); err != nil {
		t.Fatalf("nächster Schreibvorgang: %v", err)
	}
	if err := s.Write(ctx, path, sessionsMit(3)); err != nil {
		t.Fatalf("Ergänzen danach: %v", err)
	}
	if got := sessionIDs(t, path); got != "1,2,3" {
		t.Fatalf("Sessions %q", got)
	}
}

// Abdeckung: LH-FA-22/Boundary — das eingebettete SQL legt in einer In-Memory-Datenbank
// genau die Tabellen, Spalten (Name, Typ, NOT NULL) und Primärschlüssel der Liste in
// SPEC-043 an, und die Fremdschlüssel von interaction.session_id und
// message.session_id auf session.id mit ON DELETE CASCADE; message.bytes gibt es nicht.
func TestSQLiteSchema(t *testing.T) {
	db, err := sql.Open("sqlite", "file::memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(recording.SchemaSQL()); err != nil {
		t.Fatal(err)
	}
	// je Tabelle: Spalten als name:typ:notnull:primärschlüsselposition in Tabellenreihenfolge
	for tabelle, want := range map[string]string{
		"meta":        "key:TEXT:1:1 value:TEXT:1:0",
		"session":     "id:INTEGER:1:1 startup:TEXT:1:0 server_parameters:TEXT:1:0",
		"interaction": "session_id:INTEGER:1:1 sequence:INTEGER:1:2 type:TEXT:1:0 offset_ms:INTEGER:0:0 sql:TEXT:0:0",
		"message":     "session_id:INTEGER:1:1 sequence:INTEGER:1:2 position:INTEGER:1:3 direction:TEXT:1:0 group_no:INTEGER:0:0 kind:TEXT:1:0 fields:TEXT:1:0",
	} {
		got := strings.Join(zeilen(t, db, `SELECT name || ':' || type || ':' || "notnull" || ':' || pk FROM pragma_table_info('`+tabelle+`')`), " ")
		if got != want {
			t.Errorf("Tabelle %s:\n got %s\nwant %s", tabelle, got, want)
		}
	}
	if got := zeilen(t, db, `SELECT name FROM sqlite_master WHERE type = 'table' ORDER BY name`); !reflect.DeepEqual(got, []string{"interaction", "message", "meta", "session"}) {
		t.Errorf("Tabellen: %q", got)
	}
	for _, tabelle := range []string{"interaction", "message"} {
		got := zeilen(t, db, `SELECT "table" || '.' || "from" || '->' || "to" || ' ' || on_delete FROM pragma_foreign_key_list('`+tabelle+`')`)
		if want := []string{"session.session_id->id CASCADE"}; !reflect.DeepEqual(got, want) {
			t.Errorf("Fremdschlüssel von %s: %q, erwartet %q", tabelle, got, want)
		}
	}
	if got := zeilen(t, db, `SELECT count(*) FROM pragma_foreign_key_list('session') UNION ALL SELECT count(*) FROM pragma_foreign_key_list('meta')`); !reflect.DeepEqual(got, []string{"0", "0"}) {
		t.Errorf("Fremdschlüssel von session und meta: %q", got)
	}
}

// Abdeckung: LH-FA-22/Boundary, LH-FA-22/Negative — die Fabrik New liefert für yaml den
// YAML-Adapter, für sqlite den SQLite-Adapter und für jeden anderen Wert, auch den
// leeren, SQLite und sqlite3, PGR-E2001.
func TestNewFormat(t *testing.T) {
	if r, err := recording.New("yaml"); err != nil || !reflect.DeepEqual(r, recording.YAML{}) {
		t.Errorf("yaml: %#v, %v", r, err)
	}
	if r, err := recording.New("sqlite"); err != nil || reflect.TypeOf(r) != reflect.TypeOf(&recording.SQLite{}) {
		t.Errorf("sqlite: %T, %v", r, err)
	}
	for _, f := range []string{"", "SQLite", "sqlite3", "YAML", " yaml"} {
		if r, err := recording.New(f); code(err) != model.CodeUsage || r != nil || strings.Contains(err.Error(), "sqlite3") {
			t.Errorf("%q: %v, %v", f, r, err)
		}
	}
}
