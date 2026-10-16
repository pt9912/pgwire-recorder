package recording

import (
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// Die Werte von --format (LH-FA-22.a *Wahl und Erkennung*).
const (
	formatYAML   = "yaml"
	formatSQLite = "sqlite"
)

// New liefert den Adapter, der Aufzeichnungen im Format format schreibt
// (formatYAML oder formatSQLite). Der
// Adapter von sqlite führt den Stand seines Laufs und gehört zu genau einem `record`-Lauf.
// Jedes andere Format ist PGR-E2001.
func New(format string) (driven.RecordingRepository, error) {
	switch format {
	case formatYAML:
		return YAML{}, nil
	case formatSQLite:
		return &SQLite{}, nil
	}
	return nil, model.Errorf(model.CodeUsage, nil, "Aufzeichnungsformat unbekannt")
}
