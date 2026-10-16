package recording_test

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// roh öffnet die Datei path nur lesend mit der Bibliothek selbst, ohne den
// Adapter und ohne sqlite3: Die Tests lesen mit rohem SQL gegen die Liste in
// SPEC-043. Der Pfad enthält keines der Zeichen ?, # und %.
func roh(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// zeilen führt anfrage aus und liefert je Zeile ihre Spalten, mit | verbunden:
// NULL als NULL, Bytes als x'..', alles andere im Format %v.
func zeilen(t *testing.T, db *sql.DB, anfrage string, args ...any) []string {
	t.Helper()
	rows, err := db.Query(anfrage, args...)
	if err != nil {
		t.Fatalf("%s: %v", anfrage, err)
	}
	defer func() { _ = rows.Close() }()
	cols, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for rows.Next() {
		werte := make([]any, len(cols))
		ziele := make([]any, len(cols))
		for i := range werte {
			ziele[i] = &werte[i]
		}
		if err := rows.Scan(ziele...); err != nil {
			t.Fatal(err)
		}
		teile := make([]string, len(cols))
		for i, w := range werte {
			switch v := w.(type) {
			case nil:
				teile[i] = "NULL"
			case []byte:
				teile[i] = fmt.Sprintf("x'%x'", v)
			default:
				teile[i] = fmt.Sprintf("%v", v)
			}
		}
		out = append(out, strings.Join(teile, "|"))
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

// sessionIDs liefert die Nummern in der Tabelle session, aufsteigend.
func sessionIDs(t *testing.T, path string) string {
	t.Helper()
	return strings.Join(zeilen(t, roh(t, path), `SELECT id FROM session ORDER BY id`), ",")
}

// sessionsMit liefert eine Aufzeichnung mit n Sessions, jede mit einer
// einfachen Anfrage SELECT <Nummer>.
func sessionsMit(n int) model.Recording {
	rec := model.NewRecording()
	for i := 1; i <= n; i++ {
		rec.Sessions = append(rec.Sessions, model.Session{
			ID:      i,
			Startup: map[string]string{"user": "app"},
			Interactions: []model.Interaction{{
				Sequence: 1,
				Request:  model.Request{Type: model.RequestQuery, SQL: fmt.Sprintf("SELECT %d", i)},
				Responses: []model.Response{
					{Type: model.ResponseCommandComplete, Tag: "SELECT 1"},
					{Type: model.ResponseReadyForQuery, TxStatus: "I"},
				},
			}},
		})
	}
	return rec
}

// rohSchreibe führt anweisung mit der Bibliothek selbst in der Datei path aus.
func rohSchreibe(t *testing.T, path, anweisung string) {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+path+"?mode=rw")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	if _, err := db.Exec(anweisung); err != nil {
		t.Fatalf("%s: %v", anweisung, err)
	}
}

// ersetze ersetzt die Datei path durch eine neue SQLite-Datei mit rec.
func ersetze(t *testing.T, path string, rec model.Recording) {
	t.Helper()
	_, neu := neueDatei(t, rec)
	if err := os.Rename(neu, path); err != nil {
		t.Fatal(err)
	}
}
