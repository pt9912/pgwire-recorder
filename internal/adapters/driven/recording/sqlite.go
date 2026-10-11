package recording

import (
	"context"
	"database/sql"
	_ "embed" // schemaSQL
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	// Der Treiber meldet sich unter dem Namen "sqlite" an (SPEC-044).
	_ "modernc.org/sqlite"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// schemaSQL ist die Ausgabe von `make schema-generate` für tools/schema/schema.yaml;
// wer das Schema ändert, erzeugt die Datei schema.sql neu. TestSQLiteSchema prüft
// dieses SQL gegen die Liste in SPEC-043, nicht gegen schema.yaml.
//
//go:embed schema.sql
var schemaSQL string

// sperrfristMS ist die Frist in Millisekunden, die ein Zugriff auf die Sperre einer
// anderen Verbindung wartet, beim Schreiben wie beim Lesen (LH-FA-22.a *Schreiben von
// `sqlite`*, *Lesen von `sqlite`*).
const sperrfristMS = 5000

// SQLite speichert Aufzeichnungen im Format sqlite (SPEC-043, LH-FA-22.a). Der Wert
// führt den Stand eines Laufs und ist nur als Zeiger zu benutzen.
type SQLite struct {
	mu sync.Mutex
	// angelegt: der erste Schreibvorgang hat die Datei auf den Zielpfad verschoben.
	angelegt bool
	// ergaenzt zählt die Sessions, die die Datei trägt: die Nummern 1 bis ergaenzt.
	ergaenzt int
}

var _ driven.RecordingRepository = (*SQLite)(nil)

// Prepare prüft den Zielpfad wie YAML.Prepare (LH-FA-07.a *Zielpfad beim Start*),
// gleich welchen Inhalts eine vorhandene Datei hat. Danach, nie davor, ist ein
// vorhandener Pfad <path>-journal PGR-E3001, auch mit replace und auch ohne
// Zieldatei (LH-FA-22.a *Schreiben von `sqlite`*). Vorhanden heißt wie beim Zielpfad:
// Der Pfad besteht, eine symbolische Verknüpfung zählt nach ihrem Ziel, eine
// Verknüpfung ins Leere ist nicht vorhanden; lässt sich der Zustand nicht
// feststellen, ist das ebenfalls PGR-E3001.
func (*SQLite) Prepare(_ context.Context, path string, replace bool) error {
	return pruefeSQLite(path, replace, betriebssystem())
}

// pruefeSQLite ist Prepare mit den Operationen ops.
func pruefeSQLite(path string, replace bool, ops dateiOps) error {
	if err := pruefe(path, replace, ops); err != nil {
		return err
	}
	journal := journalName(path)
	_, err := os.Stat(journal)
	switch {
	case err == nil:
		return model.Errorf(model.CodeRecordingIO, nil, "%s vorhanden: Rollback-Journal eines früheren Laufs", journal)
	case errors.Is(err, fs.ErrNotExist):
		return nil
	}
	return model.Errorf(model.CodeRecordingIO, err, "%s nicht prüfbar", journal)
}

// journalName ist der Name des Rollback-Journals der Datei path.
func journalName(path string) string { return path + "-journal" }

// Write ergänzt die Datei um die Sessions von rec, die sie noch nicht trägt, jede in
// einer eigenen Transaktion (LH-FA-22.a *Schreiben von `sqlite`*). Der erste Aufruf
// legt die Datei über eine temporäre Datei im Zielverzeichnis an und verschiebt sie
// auf path; Name, Rechte und Fehlschlag der temporären Datei gelten wie bei
// YAML.Write. Ein Abbruch von ctx unterbricht keine Transaktion. Jeder Fehler der
// Bibliothek ist PGR-E3001.
func (s *SQLite) Write(ctx context.Context, path string, rec model.Recording) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.schreibe(context.WithoutCancel(ctx), path, rec, betriebssystemSQLite())
}

// sqliteOps sind die Operationen, über die schreibe die temporäre Datei benennt, mit
// Rechten versieht, schließt, mit Schema und meta füllt, verschiebt und entfernt und
// die Datenbank schließt.
type sqliteOps struct {
	zufall      func([]byte) (int, error)
	rechte      func(*os.File, fs.FileMode) error
	schliessen  func(*os.File) error
	anlegen     func(ctx context.Context, tmp string, rec model.Recording) error
	verschieben func(alt, neu string) error
	entfernen   func(name string) error
	schliesseDB func(*sql.DB) error
}

// betriebssystemSQLite liefert die Operationen des Betriebssystems und der Bibliothek.
func betriebssystemSQLite() sqliteOps {
	return sqliteOps{
		zufall:      betriebssystem().zufall,
		rechte:      (*os.File).Chmod,
		schliessen:  (*os.File).Close,
		anlegen:     legeSchemaAn,
		verschieben: os.Rename,
		entfernen:   os.Remove,
		schliesseDB: (*sql.DB).Close,
	}
}

// schreibe ist Write mit den Operationen ops. Zuerst prüft es, dass die Sessions der
// Aufzeichnung ab 1 lückenlos zählen (sonst PGR-E1000, ein Fehler des Aufrufers); die
// Zeilen der Sessions, die noch fehlen, entstehen vor dem ersten Zugriff auf die Datei.
// Ein Fehler darin lässt sie unberührt, eine noch nicht angelegte Datei entsteht nicht.
func (s *SQLite) schreibe(ctx context.Context, path string, rec model.Recording, ops sqliteOps) error {
	for i, session := range rec.Sessions {
		if session.ID != i+1 {
			return model.Errorf(model.CodeInternal, nil, "Session %d trägt die Nummer %d", i+1, session.ID)
		}
	}
	if s.ergaenzt > len(rec.Sessions) {
		return model.Errorf(model.CodeInternal, nil, "Aufzeichnung trägt weniger Sessions als die Datei")
	}
	offen := rec.Sessions[s.ergaenzt:]
	zeilen := make([]sessionZeilen, len(offen))
	for i, session := range offen {
		z, err := zeilenDerSession(session)
		if err != nil {
			return model.Errorf(model.CodeInternal, err, "Session %d nicht zu serialisieren", session.ID)
		}
		zeilen[i] = z
	}
	if !s.angelegt {
		if err := legeAn(ctx, path, rec, ops); err != nil {
			return err
		}
		s.angelegt = true
	}
	db, err := oeffne(ctx, path)
	if err != nil {
		return model.Errorf(model.CodeRecordingIO, err, "%s nicht zu öffnen", path)
	}
	err = errors.Join(s.ergaenze(ctx, db, zeilen), ops.schliesseDB(db))
	if err != nil {
		return model.Errorf(model.CodeRecordingIO, err, "%s nicht zu schreiben", path)
	}
	return nil
}

// ergaenze prüft, dass die Datei genau die Sessions 1 bis s.ergaenzt trägt, und
// ergänzt danach die Sessions aus zeilen, jede in einer Transaktion; nach jeder
// abgeschlossenen zählt s.ergaenzt, bei einem Fehler endet es ohne die übrigen.
func (s *SQLite) ergaenze(ctx context.Context, db *sql.DB, zeilen []sessionZeilen) error {
	if err := pruefeSessions(ctx, db, s.ergaenzt); err != nil {
		return err
	}
	for _, z := range zeilen {
		if err := ergaenzeSession(ctx, db, z); err != nil {
			return err
		}
		s.ergaenzt++
	}
	return nil
}

// legeAn legt die Datei auf path an: In einer temporären Datei im Verzeichnis von
// path entstehen Schema und meta einer Aufzeichnung ohne Sessions, danach wird sie
// auf path verschoben. Scheitert ein Schritt, entfernt legeAn die temporäre Datei;
// scheitert auch das Entfernen, bleibt sie liegen, und der Fehler des Entfernens
// folgt als Ursache derselben Meldung PGR-E3001 (LH-FA-07.a *Fehlschlag*).
func legeAn(ctx context.Context, path string, rec model.Recording, ops sqliteOps) error {
	tmp, err := neueTempDatei(path, ops.zufall)
	if err != nil {
		return model.Errorf(model.CodeRecordingIO, err, "temporäre Datei für %s nicht anzulegen", path)
	}
	var rerr error
	if info, err := os.Stat(path); err == nil {
		rerr = ops.rechte(tmp, info.Mode().Perm())
	}
	name := tmp.Name()
	err = errors.Join(rerr, ops.schliessen(tmp))
	if err == nil {
		err = ops.anlegen(ctx, name, rec)
	}
	if err != nil {
		return model.Errorf(model.CodeRecordingIO, errors.Join(err, ops.entfernen(name)), "temporäre Datei für %s nicht zu schreiben", path)
	}
	if err := ops.verschieben(name, path); err != nil {
		return model.Errorf(model.CodeRecordingIO, errors.Join(err, ops.entfernen(name)), "%s nicht zu schreiben", path)
	}
	return nil
}

// legeSchemaAn legt in der Datei tmp in einer Transaktion Schema und meta an: die
// Formatkennung, die Version und bei rec.EmptySessions die Zeile empty_sessions.
func legeSchemaAn(ctx context.Context, tmp string, rec model.Recording) (err error) {
	db, err := oeffne(ctx, tmp)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, db.Close()) }()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	meta := [][2]string{{"format", rec.Format}, {"version", strconv.Itoa(rec.Version)}}
	if rec.EmptySessions {
		meta = append(meta, [2]string{"empty_sessions", "true"})
	}
	if _, err := tx.ExecContext(ctx, schemaSQL); err != nil {
		return errors.Join(err, tx.Rollback())
	}
	for _, m := range meta {
		if _, err := tx.ExecContext(ctx, `INSERT INTO "meta" ("key", "value") VALUES (?, ?)`, m[0], m[1]); err != nil {
			return errors.Join(err, tx.Rollback())
		}
	}
	return tx.Commit()
}

// oeffne öffnet die vorhandene Datei path mit den Einstellungen von LH-FA-22.a
// *Schreiben von `sqlite`*: mode=rw (ohne legt SQLite die Datei still an), das
// Rollback-Journal, foreign_keys=ON, synchronous=FULL, eine Sperrfrist von
// sperrfristMS und BEGIN IMMEDIATE für jede Transaktion. Es prüft, dass sich die
// Datenbank öffnen lässt.
func oeffne(ctx context.Context, path string) (*sql.DB, error) {
	dsn, err := dsnVon(path)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Join(err, db.Close())
	}
	return db, nil
}

// dsnVon liefert die file:-URI von path mit den Einstellungen von oeffne.
func dsnVon(path string) (string, error) {
	uri, err := uriVon(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?mode=rw&_txlock=immediate&_pragma=busy_timeout(%d)&_pragma=foreign_keys(ON)&_pragma=journal_mode(DELETE)&_pragma=synchronous(FULL)", uri, sperrfristMS), nil
}

// uriVon liefert die file:-URI von path ohne Parameter, für Schreiben und Lesen. Der
// Pfad gilt wörtlich: %, ? und # stehen als %XX, damit sie nichts anderes bedeuten
// (TestSQLiteSonderzeichenImPfad, TestLesenSonderzeichenImPfad).
func uriVon(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	abs = filepath.ToSlash(abs)
	if runtime.GOOS == "windows" {
		abs = "/" + abs
	}
	var b strings.Builder
	b.WriteString("file:")
	for i := 0; i < len(abs); i++ {
		switch c := abs[i]; {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', strings.IndexByte("-._~/", c) >= 0:
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String(), nil
}

// pruefeSessions prüft, dass die Tabelle session genau die Nummern 1 bis erwartet
// trägt; jede andere Datei ist ein Fehler (eine fehlende oder ersetzte Datei).
func pruefeSessions(ctx context.Context, db *sql.DB, erwartet int) error {
	rows, err := db.QueryContext(ctx, `SELECT "id" FROM "session" ORDER BY "id"`)
	if err != nil {
		return err
	}
	n := 0
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return errors.Join(err, rows.Close())
		}
		n++
		if id != int64(n) {
			return errors.Join(fmt.Errorf("Datei trägt nicht die Sessions 1 bis %d", erwartet), rows.Close())
		}
	}
	if err := errors.Join(rows.Err(), rows.Close()); err != nil {
		return err
	}
	if n != erwartet {
		return fmt.Errorf("Datei trägt nicht die Sessions 1 bis %d", erwartet)
	}
	return nil
}
