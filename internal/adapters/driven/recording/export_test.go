package recording

import (
	"context"
	"database/sql"
	"io/fs"
	"os"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Eingriffe ersetzt für PruefeMit und SchreibeMit einzelne Operationen von
// Prepare und Write; ein Feld mit nil ist die Operation des Betriebssystems.
type Eingriffe struct {
	Probe           func(dir, muster string) (*os.File, error)
	Zufall          func([]byte) (int, error)
	Rechte          func(*os.File, os.FileMode) error
	Schreiben       func(*os.File, []byte) (int, error)
	Synchronisieren func(*os.File) error
	Schliessen      func(*os.File) error
	Verschieben     func(alt, neu string) error
	Entfernen       func(name string) error
}

// PruefeMit prüft path wie Prepare, mit den Operationen aus e.
func PruefeMit(path string, replace bool, e Eingriffe) error {
	return pruefe(path, replace, e.ops())
}

// SchreibeMit schreibt rec wie Write, mit den Operationen aus e.
func SchreibeMit(path string, rec model.Recording, e Eingriffe) error {
	return schreibe(path, rec, e.ops())
}

// ops liefert die Operationen des Betriebssystems, ersetzt durch die Felder von
// e, die nicht nil sind.
func (e Eingriffe) ops() dateiOps {
	ops := betriebssystem()
	if e.Probe != nil {
		ops.probe = e.Probe
	}
	if e.Zufall != nil {
		ops.zufall = e.Zufall
	}
	if e.Rechte != nil {
		ops.rechte = e.Rechte
	}
	if e.Schreiben != nil {
		ops.schreiben = e.Schreiben
	}
	if e.Synchronisieren != nil {
		ops.synchronisieren = e.Synchronisieren
	}
	if e.Schliessen != nil {
		ops.schliessen = e.Schliessen
	}
	if e.Verschieben != nil {
		ops.verschieben = e.Verschieben
	}
	if e.Entfernen != nil {
		ops.entfernen = e.Entfernen
	}
	return ops
}

// SQLiteEingriffe ersetzt für SQLiteSchreibeMit einzelne Operationen von
// SQLite.Write; ein Feld mit nil ist die Operation des Betriebssystems und der
// Bibliothek.
type SQLiteEingriffe struct {
	Zufall      func([]byte) (int, error)
	Rechte      func(*os.File, fs.FileMode) error
	Schliessen  func(*os.File) error
	Anlegen     func(ctx context.Context, tmp string, rec model.Recording) error
	Verschieben func(alt, neu string) error
	Entfernen   func(name string) error
	SchliesseDB func(*sql.DB) error
}

// SQLiteSchreibeMit schreibt rec wie s.Write, mit den Operationen aus e.
func SQLiteSchreibeMit(ctx context.Context, s *SQLite, path string, rec model.Recording, e SQLiteEingriffe) error {
	ops := betriebssystemSQLite()
	if e.Zufall != nil {
		ops.zufall = e.Zufall
	}
	if e.Rechte != nil {
		ops.rechte = e.Rechte
	}
	if e.Schliessen != nil {
		ops.schliessen = e.Schliessen
	}
	if e.Anlegen != nil {
		ops.anlegen = e.Anlegen
	}
	if e.Verschieben != nil {
		ops.verschieben = e.Verschieben
	}
	if e.Entfernen != nil {
		ops.entfernen = e.Entfernen
	}
	if e.SchliesseDB != nil {
		ops.schliesseDB = e.SchliesseDB
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.schreibe(ctx, path, rec, ops)
}

// PruefeSQLiteMit prüft path wie SQLite.Prepare, mit den Operationen aus e.
func PruefeSQLiteMit(path string, replace bool, e Eingriffe) error {
	return pruefeSQLite(path, replace, e.ops())
}

// SchemaSQL ist das eingebettete SQL des Schemas.
func SchemaSQL() string { return schemaSQL }

// DSN liefert die file:-URI, mit der SQLite.Write path öffnet.
func DSN(path string) (string, error) { return dsnVon(path) }

// Oeffne öffnet path wie SQLite.Write.
func Oeffne(ctx context.Context, path string) (*sql.DB, error) { return oeffne(ctx, path) }

// JSONText wandelt wert wie SQLite.Write für eine JSON-Spalte in JSON-Text, ohne den
// Schlüssel ohne der obersten Abbildung.
func JSONText(wert any, ohne string) (string, error) { return jsonText(wert, ohne) }

// LeseMit liest path wie SQLite.Load; nach wird nach dem Lesen jeder Tabelle mit der
// Lesetransaktion und dem Namen der Tabelle gerufen.
func LeseMit(ctx context.Context, path string, nach func(tx *sql.Tx, tabelle string)) (model.Recording, error) {
	text, err := leseAlsYAML(ctx, path, nach)
	if err != nil {
		return model.Recording{}, err
	}
	return ladeInhalt(path, text)
}

// StufeCode liefert den Meldungscode, den SQLite.Load für den Ergebniscode code der
// Bibliothek meldet.
func StufeCode(code int) string { return stufeCode(code) }

// JSONKnoten liest den JSON-Text text wie SQLite.Load und liefert Tag und Wert des
// Skalars, oder "" und "" für einen Knoten, der kein Skalar ist.
func JSONKnoten(text string) (tag, wert string, err error) {
	n, err := jsonKnoten(text)
	if err != nil {
		return "", "", err
	}
	return n.Tag, n.Value, nil
}
