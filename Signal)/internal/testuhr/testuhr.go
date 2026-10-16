package testuhr

import (
	"context"
	"sync"
)

// Uhr zählt Millisekunden ab 0 und rückt nur durch Vorruecken vor.
type Uhr struct {
	mu       sync.Mutex
	jetzt    int64
	wartende map[*wartender]struct{}
}

type wartender struct {
	zeitpunkt int64
	geweckt   chan struct{}
}

// Neu liefert eine Uhr, die auf 0 steht.
func Neu() *Uhr {
	return &Uhr{wartende: map[*wartender]struct{}{}}
}

// Jetzt liefert den Stand der Uhr.
func (u *Uhr) Jetzt() int64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.jetzt
}

// WarteBis nach der Zusage des Ports: ein beendeter ctx geht vor, auch wenn
// Vorruecken den Zeitpunkt zugleich erreicht; ein erreichter Zeitpunkt kehrt
// sofort mit nil zurück, sonst wartet der Aufruf auf Vorruecken oder das Ende von
// ctx.
func (u *Uhr) WarteBis(ctx context.Context, zeitpunkt int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	u.mu.Lock()
	if zeitpunkt <= u.jetzt {
		u.mu.Unlock()
		return nil
	}
	w := &wartender{zeitpunkt: zeitpunkt, geweckt: make(chan struct{})}
	u.wartende[w] = struct{}{}
	u.mu.Unlock()
	select {
	case <-w.geweckt:
		return ctx.Err()
	case <-ctx.Done():
		u.mu.Lock()
		delete(u.wartende, w)
		u.mu.Unlock()
		return ctx.Err()
	}
}

// Vorruecken rückt die Uhr um ms Millisekunden vor und weckt jeden Wartenden,
// dessen Zeitpunkt erreicht ist. Eine negative Zahl ist ein Programmierfehler
// des Tests und löst eine Panik aus.
func (u *Uhr) Vorruecken(ms int64) {
	if ms < 0 {
		panic("testuhr: Vorruecken um eine negative Zahl")
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	u.jetzt += ms
	for w := range u.wartende {
		if w.zeitpunkt <= u.jetzt {
			delete(u.wartende, w)
			close(w.geweckt)
		}
	}
}

// Wartende liefert die Zahl der Aufrufe von WarteBis, die gerade warten.
func (u *Uhr) Wartende() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return len(u.wartende)
}
