package driven

import (
	"context"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// RecordingRepository ist der Port der Aufzeichnung (ARC-004).
type RecordingRepository interface {
	// Prepare prüft den Zielpfad vor dem ersten Schreiben: ein vorhandener Pfad
	// ohne replace ist PGR-E2002.
	Prepare(ctx context.Context, path string, replace bool) error
	// Write schreibt die Aufzeichnung als Ganzes; die Zieldatei ist danach
	// vollständig oder unverändert.
	Write(ctx context.Context, path string, rec model.Recording) error
	// Load liest eine Aufzeichnung.
	Load(ctx context.Context, path string) (model.Recording, error)
}
