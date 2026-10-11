package bootstrap_test

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
	"github.com/pt9912/pgwire-recorder/internal/testuhr"
)

// Die Fake-Uhr erfüllt den Port; der Compiler prüft es.
var _ driven.Uhr = (*testuhr.Uhr)(nil)

// wartenBinnen führt WarteBis aus und meldet, wenn der Aufruf nicht binnen
// grenze zurückkehrt.
func wartenBinnen(ctx context.Context, t *testing.T, u driven.Uhr, zeitpunkt int64, grenze time.Duration) error {
	t.Helper()
	ergebnis := make(chan error, 1)
	go func() { ergebnis <- u.WarteBis(ctx, zeitpunkt) }()
	select {
	case err := <-ergebnis:
		return err
	case <-time.After(grenze):
		t.Fatalf("WarteBis(%d) kehrt nicht binnen %v zurück", zeitpunkt, grenze)
		return nil
	}
}

// Der Ursprung der Systemuhr liegt bei ihrer Erzeugung: die erste Lesung ist
// klein, und sie wächst mit der Zeit.
func TestSystemuhrUrsprungBeiErzeugung(t *testing.T) {
	u := bootstrap.NeueSystemuhr()
	if j := u.Jetzt(); j < 0 || j > 5000 {
		t.Fatalf("erste Lesung %d ms", j)
	}
	time.Sleep(30 * time.Millisecond)
	if j := u.Jetzt(); j < 30 {
		t.Fatalf("Lesung 30 ms nach der Erzeugung: %d", j)
	}
}

// Die Lesung ist die abgeschnittene Zahl ganzer Millisekunden seit der
// Erzeugung, gerahmt von zwei Messungen um sie herum.
func TestSystemuhrLesungAbgeschnitten(t *testing.T) {
	vorher := time.Now()
	u := bootstrap.NeueSystemuhr()
	nachErzeugung := time.Now()
	for i := range 300 {
		time.Sleep(time.Duration(300+i*7%900) * time.Microsecond)
		davor := time.Now()
		j := u.Jetzt()
		danach := time.Now()
		if obere := danach.Sub(vorher).Milliseconds(); j > obere {
			t.Fatalf("Lesung %d ms über der Obergrenze %d ms (Lesung %d)", j, obere, i)
		}
		if untere := davor.Sub(nachErzeugung).Milliseconds(); j < untere {
			t.Fatalf("Lesung %d ms unter der Untergrenze %d ms (Lesung %d)", j, untere, i)
		}
	}
}

// Die Lesung wird nie kleiner.
func TestSystemuhrMonoton(t *testing.T) {
	u := bootstrap.NeueSystemuhr()
	vorige := u.Jetzt()
	for i := range 2000 {
		if i%100 == 0 {
			time.Sleep(time.Millisecond)
		}
		j := u.Jetzt()
		if j < vorige {
			t.Fatalf("Lesung %d nach %d", j, vorige)
		}
		vorige = j
	}
}

// Die Systemuhr zählt aus der monotonen Zeit des Prozesses: ihr Ursprung trägt
// deren Messwert. Dass sie einer Änderung der Wanduhr nicht folgt, prüft kein
// Test.
func TestSystemuhrUrsprungTraegtMonotonenMesswert(t *testing.T) {
	if !bootstrap.UrsprungMonoton(bootstrap.NeueSystemuhr()) {
		t.Fatal("der Ursprung der Systemuhr trägt keine monotone Zeit")
	}
}

// Warten kehrt erst zurück, wenn die Lesung den Zeitpunkt erreicht hat.
func TestSystemuhrWartenUntergrenze(t *testing.T) {
	u := bootstrap.NeueSystemuhr()
	for _, ms := range []int64{1, 50} {
		zeitpunkt := u.Jetzt() + ms
		if err := wartenBinnen(context.Background(), t, u, zeitpunkt, 5*time.Second); err != nil {
			t.Fatalf("WarteBis: %v", err)
		}
		if j := u.Jetzt(); j < zeitpunkt {
			t.Fatalf("WarteBis(%d) kehrte bei Lesung %d zurück", zeitpunkt, j)
		}
	}
}

// Ein erreichter oder vergangener Zeitpunkt kehrt sofort mit nil zurück, auch
// der kleinste und ein Zeitpunkt, dessen Abstand die Zahl überläuft.
func TestSystemuhrVergangenerZeitpunktSofort(t *testing.T) {
	u := bootstrap.NeueSystemuhr()
	time.Sleep(5 * time.Millisecond)
	for _, z := range []int64{u.Jetzt() - 1, 0, -1, math.MinInt64} {
		start := time.Now()
		if err := wartenBinnen(context.Background(), t, u, z, 5*time.Second); err != nil {
			t.Fatalf("WarteBis(%d): %v", z, err)
		}
		if d := time.Since(start); d > 500*time.Millisecond {
			t.Fatalf("WarteBis(%d) brauchte %v", z, d)
		}
	}
}

// Ein beendeter Kontext geht vor, auch bei vergangenem Zeitpunkt, und der Aufruf
// liefert dessen Fehler.
func TestSystemuhrBeendeterKontextVorDemWarten(t *testing.T) {
	abgebrochen, cancel := context.WithCancel(context.Background())
	cancel()
	frist, cancelFrist := context.WithTimeout(context.Background(), 0)
	defer cancelFrist()
	u := bootstrap.NeueSystemuhr()
	for _, tc := range []struct {
		name string
		ctx  func() context.Context
		want error
	}{
		{"abgebrochen", func() context.Context { return abgebrochen }, context.Canceled},
		{"Frist abgelaufen", func() context.Context { return frist }, context.DeadlineExceeded},
	} {
		for _, z := range []int64{-1, 0, u.Jetzt() + 60000, math.MaxInt64} {
			if err := wartenBinnen(tc.ctx(), t, u, z, 5*time.Second); !errors.Is(err, tc.want) {
				t.Fatalf("%s, WarteBis(%d): %v, erwartet %v", tc.name, z, err, tc.want)
			}
		}
	}
}

// Endet der Kontext während des Wartens, kehrt der Aufruf mit dessen Fehler
// zurück, lange vor dem Zeitpunkt; das gilt auch für einen Zeitpunkt, dessen
// Abstand über der größten Dauer der Zeitbibliothek liegt.
func TestSystemuhrAbbruchWaehrendDesWartens(t *testing.T) {
	u := bootstrap.NeueSystemuhr()
	for _, z := range []int64{u.Jetzt() + 600000, math.MaxInt64} {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(40 * time.Millisecond)
			cancel()
		}()
		start := time.Now()
		if err := wartenBinnen(ctx, t, u, z, 5*time.Second); !errors.Is(err, context.Canceled) {
			t.Fatalf("WarteBis(%d) nach Abbruch: %v", z, err)
		}
		if d := time.Since(start); d < 30*time.Millisecond {
			t.Fatalf("WarteBis(%d) kehrte nach %v zurück, vor dem Abbruch", z, d)
		}
		fristCtx, cancelFrist := context.WithTimeout(context.Background(), 40*time.Millisecond)
		err := wartenBinnen(fristCtx, t, u, z, 5*time.Second)
		cancelFrist()
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("WarteBis(%d) nach Fristablauf: %v", z, err)
		}
	}
}
