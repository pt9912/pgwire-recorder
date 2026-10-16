package offsetlese

import (
	"context"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
)

// Fehlt steht für eine Interaktion ohne offset_ms.
const Fehlt int64 = -1

// Laden lädt die Aufzeichnung path mit dem Leser, der das Format erkennt, und
// liefert je Session offset_ms jeder Interaktion in der Reihenfolge der Datei, Fehlt
// für eine Interaktion ohne das Feld.
func Laden(path string) ([][]int64, error) {
	rec, err := recording.Leser{}.Load(context.Background(), path)
	if err != nil {
		return nil, err
	}
	var out [][]int64
	for _, s := range rec.Sessions {
		var werte []int64
		for _, in := range s.Interactions {
			if in.OffsetMS == nil {
				werte = append(werte, Fehlt)
				continue
			}
			werte = append(werte, *in.OffsetMS)
		}
		out = append(out, werte)
	}
	return out, nil
}
