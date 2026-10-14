package model

import (
	"context"
	"errors"
)

// Anmeldeanforderung ist eine Anforderung des Servers im Anmeldeaustausch
// (LH-FA-05.b *Vermittlung*): der Code der Anforderung und der Rest der
// Nachricht dahinter. Der Kern liest beides nicht.
type Anmeldeanforderung struct {
	Code  uint32
	Daten []byte
}

// Anmeldung gibt eine Anforderung des Servers unverändert an den Client weiter
// und liefert den Rumpf seiner Antwort unverändert zurück; auf eine Anforderung,
// die keine Antwort erwartet, wartet sie nicht und liefert nil. Endet der Client
// die Verbindung in der Anmeldung, liefert sie ErrClientEnde.
//
// Aufrufe einer Verbindung kommen nacheinander: Ein Aufruf beginnt erst, wenn
// der vorige zurück ist (LH-FA-05.b *Aufruf-Voraussetzung*).
type Anmeldung func(ctx context.Context, a Anmeldeanforderung) ([]byte, error)

// ErrClientEnde liefert Anmeldung, wenn der Client die Verbindung in der
// Anmeldung beendet (Ende des Datenstroms, Reset, Terminate, Zwangsende). Es ist
// kein Verbindungsfehler und erzeugt weder Meldung noch Session (LH-FA-05.b
// *Ende durch den Client*).
var ErrClientEnde = errors.New("Client beendet die Verbindung in der Anmeldung")
