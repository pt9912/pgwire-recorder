package bootstrap

import (
	"context"
	"math"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

var _ driven.Uhr = (*systemuhr)(nil)

// systemuhr zählt Millisekunden ab ihrer Erzeugung aus der monotonen Zeit des
// Prozesses, nie aus der Wanduhr (LH-FA-21.a *Uhr*).
type systemuhr struct {
	ursprung time.Time
}

// NeueSystemuhr liefert die Systemuhr; ihr Ursprung ist der Zeitpunkt des Aufrufs.
func NeueSystemuhr() driven.Uhr {
	return &systemuhr{ursprung: time.Now()}
}

// Jetzt liefert die seit der Erzeugung vergangenen ganzen Millisekunden.
func (u *systemuhr) Jetzt() int64 {
	return int64(time.Since(u.ursprung) / time.Millisecond)
}

// WarteBis nach der Zusage des Ports. Ein Zeitpunkt, dessen Abstand die größte
// Dauer der Zeitbibliothek übersteigt, wird in Stufen dieser Dauer abgewartet.
func (u *systemuhr) WarteBis(ctx context.Context, zeitpunkt int64) error {
	const stufeMs = int64(math.MaxInt64 / int64(time.Millisecond))
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		jetzt := u.Jetzt()
		if zeitpunkt <= jetzt {
			return nil
		}
		rest := zeitpunkt - jetzt
		if rest > stufeMs {
			rest = stufeMs
		}
		timer := time.NewTimer(time.Duration(rest) * time.Millisecond)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}
}
