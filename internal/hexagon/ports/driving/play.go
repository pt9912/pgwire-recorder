package driving

import "context"

// Player ist der Play-Use-Case (ARC-003): Er spielt die Anfragen einer
// Aufzeichnung gegen einen Server ein; die CLI startet ihn, keine
// Client-Verbindung (LH-FA-20.a).
type Player interface {
	// Play spielt ein, bis die Aufzeichnung zu Ende ist, ein Fehler abbricht
	// oder ein Abbruchsignal es beendet: Endet ctx (erstes Signal), endet es
	// nach der laufenden Interaktion; wird ablauf geschlossen (zweites Signal),
	// sofort. Es liefert den Fehler, nach dem es abgebrochen hat, sonst nil.
	Play(ctx context.Context, ablauf <-chan struct{}) error
}
