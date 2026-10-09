package recording

import (
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
