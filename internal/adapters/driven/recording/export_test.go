package recording

import (
	"os"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Eingriffe ersetzt für SchreibeMit einzelne Operationen von Write; ein Feld mit
// nil ist die Operation des Betriebssystems.
type Eingriffe struct {
	Zufall          func([]byte) (int, error)
	Schreiben       func(*os.File, []byte) (int, error)
	Synchronisieren func(*os.File) error
	Verschieben     func(alt, neu string) error
	Entfernen       func(name string) error
}

// SchreibeMit schreibt rec wie Write, mit den Operationen aus e.
func SchreibeMit(path string, rec model.Recording, e Eingriffe) error {
	ops := betriebssystem()
	if e.Zufall != nil {
		ops.zufall = e.Zufall
	}
	if e.Schreiben != nil {
		ops.schreiben = e.Schreiben
	}
	if e.Synchronisieren != nil {
		ops.synchronisieren = e.Synchronisieren
	}
	if e.Verschieben != nil {
		ops.verschieben = e.Verschieben
	}
	if e.Entfernen != nil {
		ops.entfernen = e.Entfernen
	}
	return schreibe(path, rec, ops)
}
