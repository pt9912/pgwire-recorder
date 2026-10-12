package services

import (
	"context"
	"math"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// pruefeZeitangaben verlangt offset_ms an jeder Interaktion jeder Session mit
// Interaktionen: PGR-E2003 nennt Session und Interaktion der ersten ohne, in der
// Reihenfolge der Aufzeichnung, keinen Wert (LH-FA-21.a *Randformen des
// Einspielens*).
func (s *PlayService) pruefeZeitangaben() error {
	for _, sess := range s.sessions {
		for _, in := range sess.Interactions {
			if in.OffsetMS == nil {
				return model.Errorf(model.CodeTimingMissing, nil, "Session %d, Interaktion %d trägt kein offset_ms, das --keep-timing an jeder Interaktion verlangt", sess.ID, in.Sequence)
			}
		}
	}
	return nil
}

// takt hält die Zeitpunkte einer Session, nach denen die Wartezeiten bemessen
// sind (LH-FA-21.a *Einspielen*); ohne KeepTiming ist uhr nil, und takt wartet
// nie und liest die Uhr nie.
type takt struct {
	uhr          driven.Uhr
	absolut      bool
	ersteAnfrage bool
	// verbindung ist die Lesung am Beginn des Aufbaus.
	verbindung int64
	// begonnen sagt, dass eine Interaktion der Session begonnen hat; die beiden
	// folgenden Paare tragen Beginn und offset_ms der ersten und der
	// vorigen.
	begonnen                bool
	ersterBeginn, ersterOff int64
	vorigerBeginn, vorigOff int64
}

// neuerTakt liest, mit KeepTiming, die Uhr als Bezugspunkt connect; der Aufruf
// steht vor dem Aufbau der Verbindung.
func (s *PlayService) neuerTakt() *takt {
	if !s.optionen.KeepTiming {
		return &takt{}
	}
	return &takt{
		uhr:          s.uhr,
		absolut:      s.optionen.TimingMode == TimingAbsolute,
		ersteAnfrage: s.optionen.TimingReference == ReferenceFirstRequest,
		verbindung:   s.uhr.Jetzt(),
	}
}

// ziel ist der Zeitpunkt, vor dem die Interaktion in nicht beginnt; false heißt,
// sie wartet nicht: ohne KeepTiming, relative und absolute mit first-request vor
// der ersten Interaktion der Session.
func (t *takt) ziel(in model.Interaction) (int64, bool) {
	if t.uhr == nil {
		return 0, false
	}
	off := *in.OffsetMS
	switch {
	case !t.absolut && t.begonnen:
		return plus(t.vorigerBeginn, off-t.vorigOff), true
	case t.absolut && !t.ersteAnfrage:
		return plus(t.verbindung, off), true
	case t.absolut && t.begonnen:
		return plus(t.ersterBeginn, off-t.ersterOff), true
	}
	return 0, false
}

// beginnt hält den Beginn der Interaktion fest, die jetzt ihre erste
// Client-Nachricht sendet: die Lesung der Uhr (LH-FA-21.a *Randformen des
// Einspielens*, *Beginn einer Interaktion*).
func (t *takt) beginnt(in model.Interaction) {
	if t.uhr == nil {
		return
	}
	jetzt := t.uhr.Jetzt()
	if !t.begonnen {
		t.begonnen, t.ersterBeginn, t.ersterOff = true, jetzt, *in.OffsetMS
	}
	t.vorigerBeginn, t.vorigOff = jetzt, *in.OffsetMS
}

// plus addiert d zum Zeitpunkt a und begrenzt das Ergebnis bei einem Überlauf
// von int64 auf den größten Wert, statt negativ zu werden.
func plus(a, d int64) int64 {
	if d > 0 && a > math.MaxInt64-d {
		return math.MaxInt64
	}
	return a + d
}

// warte wartet vor der Interaktion in, bis die Uhr den Zeitpunkt des Takts
// erreicht hat, und meldet, ob die Interaktion gesendet wird. Das erste
// Abbruchsignal (ctx) beendet die Wartezeit sofort; ohne FinishSessionOnInterrupt
// wird die Interaktion dann nicht gesendet, mit ihr schon, und nach dem Signal
// ruft warte die Uhr nicht mehr. Das zweite Signal (abbruch) beendet die
// Wartezeit ebenso, und die Interaktion wird nicht gesendet (LH-FA-21.a
// *Randformen des Einspielens*, *Abbruch*).
func (s *PlayService) warte(ctx, abbruch context.Context, t *takt, in model.Interaction) bool {
	zeitpunkt, wartet := t.ziel(in)
	if !wartet || ctx.Err() != nil {
		return true
	}
	wartezeit, ende := context.WithCancel(ctx)
	defer ende()
	stop := context.AfterFunc(abbruch, ende)
	defer stop()
	if err := t.uhr.WarteBis(wartezeit, zeitpunkt); err == nil {
		return true
	}
	return s.optionen.FinishSessionOnInterrupt && ctx.Err() != nil && abbruch.Err() == nil
}
