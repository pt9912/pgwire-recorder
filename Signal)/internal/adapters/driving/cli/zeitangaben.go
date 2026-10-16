package cli

import (
	"errors"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Die Werte von --timing-mode und --timing-reference (LH-FA-21.a *Einspielen*).
const (
	zeitRelative   = "relative"
	zeitAbsolute   = "absolute"
	bezugConnect   = "connect"
	bezugErsteAnfr = "first-request"
)

// zeitmodus ist der Wert von --timing-mode: genau relative oder absolute in
// Kleinbuchstaben; jeder andere Wert, auch der leere, ist ein Fehler.
type zeitmodus struct{ wert *string }

func (z zeitmodus) Set(v string) error {
	if v != zeitRelative && v != zeitAbsolute {
		return errors.New("erlaubt sind relative und absolute")
	}
	*z.wert = v
	return nil
}

// zeitbezug ist der Wert von --timing-reference: genau connect oder
// first-request in Kleinbuchstaben; jeder andere Wert, auch der leere, ist ein
// Fehler.
type zeitbezug struct{ wert *string }

func (z zeitbezug) Set(v string) error {
	if v != bezugConnect && v != bezugErsteAnfr {
		return errors.New("erlaubt sind connect und first-request")
	}
	*z.wert = v
	return nil
}

// zeitPlay prüft nach der Zusammenführung die Kombinationen der drei Optionen
// des zeitgetreuen Einspielens und setzt die Werte im Command (LH-FA-21.a
// *Randformen des Einspielens*, *Gesetzt*). Gesetzt ist eine Option, wenn
// irgendeine Quelle sie nennt, auch mit dem Wert des Standards. --timing-mode
// und --timing-reference ohne wirksames --keep-timing sind PGR-E2001, ebenso
// --timing-reference mit dem wirksamen Modus relative; die Reihenfolge der
// Optionen ändert nichts. Die Meldung nennt Optionen, keinen Wert.
func zeitPlay(c *Command, _ gelesen, _ *datei, q quellen) error {
	v := &c.Play.Einspielen
	if !v.KeepTiming {
		for _, name := range []string{"timing-mode", "timing-reference"} {
			if q.gesetzt(name) {
				return model.Errorf(model.CodeUsage, nil, "--%s verlangt --keep-timing", name)
			}
		}
		return nil
	}
	v.TimingMode = q.wert["timing-mode"]
	if v.TimingMode == zeitAbsolute {
		v.TimingReference = q.wert["timing-reference"]
	} else if q.gesetzt("timing-reference") {
		return model.Errorf(model.CodeUsage, nil, "--timing-reference verlangt --timing-mode absolute")
	}
	return nil
}
