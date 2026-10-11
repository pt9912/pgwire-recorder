package testuhr_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/testuhr"
	"github.com/pt9912/pgwire-recorder/internal/testzeit"
)

// wartet startet WarteBis in einer Goroutine und liefert den Kanal mit ihrem
// Ergebnis, sobald der Aufruf wartet.
func wartet(ctx context.Context, t *testing.T, u *testuhr.Uhr, zeitpunkt int64) <-chan error {
	t.Helper()
	vorher := u.Wartende()
	ergebnis := make(chan error, 1)
	go func() { ergebnis <- u.WarteBis(ctx, zeitpunkt) }()
	frist := testzeit.NeueFrist(2000)
	for u.Wartende() == vorher {
		if frist.Abgelaufen() {
			t.Fatal("WarteBis wartet nicht binnen 2 s")
		}
		testzeit.Schlafe(1)
	}
	return ergebnis
}

// ergebnisBinnen liest das Ergebnis eines Aufrufs, der zurückkehren muss.
func ergebnisBinnen(t *testing.T, ergebnis <-chan error, was string) error {
	t.Helper()
	select {
	case err := <-ergebnis:
		return err
	case <-testzeit.Nach(2000):
		t.Fatalf("%s kehrt nicht binnen 2 s zurück", was)
		return nil
	}
}

// bleibt prüft, dass ein Aufruf nach kurzer Zeit noch wartet.
func bleibt(t *testing.T, ergebnis <-chan error, was string) {
	t.Helper()
	select {
	case err := <-ergebnis:
		t.Fatalf("%s kehrte zurück: %v", was, err)
	case <-testzeit.Nach(50):
	}
}

// aufruf führt WarteBis aus und meldet, wenn es nicht zurückkehrt.
func aufruf(ctx context.Context, t *testing.T, u *testuhr.Uhr, zeitpunkt int64) error {
	t.Helper()
	ergebnis := make(chan error, 1)
	go func() { ergebnis <- u.WarteBis(ctx, zeitpunkt) }()
	return ergebnisBinnen(t, ergebnis, "WarteBis")
}

// Die Uhr steht auf 0 und rückt nur durch Vorruecken vor, um genau diese Zahl.
func TestUhrRueckenNurAufAnweisung(t *testing.T) {
	u := testuhr.Neu()
	if u.Jetzt() != 0 {
		t.Fatalf("Jetzt einer neuen Uhr: %d", u.Jetzt())
	}
	testzeit.Schlafe(20)
	if u.Jetzt() != 0 {
		t.Fatalf("Jetzt nach Wartezeit ohne Vorruecken: %d", u.Jetzt())
	}
	u.Vorruecken(7)
	u.Vorruecken(0)
	u.Vorruecken(3)
	if u.Jetzt() != 10 {
		t.Fatalf("Jetzt nach 7, 0 und 3: %d", u.Jetzt())
	}
}

// Ein erreichter oder vergangener Zeitpunkt kehrt sofort mit nil zurück.
func TestUhrErreichterZeitpunktSofort(t *testing.T) {
	u := testuhr.Neu()
	u.Vorruecken(10)
	for _, z := range []int64{10, 9, 0, -5} {
		if err := aufruf(context.Background(), t, u, z); err != nil {
			t.Fatalf("WarteBis(%d) bei Jetzt 10: %v", z, err)
		}
	}
}

// Ein Zeitpunkt in der Zukunft wartet, bis Vorruecken ihn erreicht, nicht früher.
func TestUhrWartetBisZumZeitpunkt(t *testing.T) {
	u := testuhr.Neu()
	ergebnis := wartet(context.Background(), t, u, 5)
	u.Vorruecken(4)
	bleibt(t, ergebnis, "WarteBis(5) bei Jetzt 4")
	u.Vorruecken(1)
	if err := ergebnisBinnen(t, ergebnis, "WarteBis(5) bei Jetzt 5"); err != nil {
		t.Fatal(err)
	}
	if u.Wartende() != 0 {
		t.Fatalf("Wartende nach dem Wecken: %d", u.Wartende())
	}
}

// Ein Vorruecken über den Zeitpunkt hinaus weckt ebenfalls.
func TestUhrVorrueckenUeberZeitpunkt(t *testing.T) {
	u := testuhr.Neu()
	ergebnis := wartet(context.Background(), t, u, 5)
	u.Vorruecken(100)
	if err := ergebnisBinnen(t, ergebnis, "WarteBis(5) bei Jetzt 100"); err != nil {
		t.Fatal(err)
	}
}

// Vorruecken weckt jeden Wartenden, dessen Zeitpunkt erreicht ist, und keinen
// anderen.
func TestUhrWecktAlleErreichten(t *testing.T) {
	u := testuhr.Neu()
	a := wartet(context.Background(), t, u, 5)
	b := wartet(context.Background(), t, u, 5)
	c := wartet(context.Background(), t, u, 7)
	d := wartet(context.Background(), t, u, 10)
	u.Vorruecken(7)
	for name, ch := range map[string]<-chan error{"5 (a)": a, "5 (b)": b, "7": c} {
		if err := ergebnisBinnen(t, ch, "WarteBis("+name+") bei Jetzt 7"); err != nil {
			t.Fatalf("WarteBis(%s): %v", name, err)
		}
	}
	bleibt(t, d, "WarteBis(10) bei Jetzt 7")
	if u.Wartende() != 1 {
		t.Fatalf("Wartende nach Vorruecken auf 7: %d", u.Wartende())
	}
	u.Vorruecken(3)
	if err := ergebnisBinnen(t, d, "WarteBis(10) bei Jetzt 10"); err != nil {
		t.Fatal(err)
	}
}

// Ein beendeter Kontext geht vor, auch bei erreichtem oder vergangenem Zeitpunkt;
// der Aufruf liefert den Fehler des Kontexts.
func TestUhrBeendeterKontextVorDemWarten(t *testing.T) {
	abgebrochen, cancel := context.WithCancel(context.Background())
	cancel()
	frist, cancelFrist := context.WithTimeout(context.Background(), 0)
	defer cancelFrist()
	for _, tc := range []struct {
		name string
		ctx  func() context.Context
		want error
	}{
		{"abgebrochen", func() context.Context { return abgebrochen }, context.Canceled},
		{"Frist abgelaufen", func() context.Context { return frist }, context.DeadlineExceeded},
	} {
		for _, z := range []int64{-1, 0, 5, 100} {
			u := testuhr.Neu()
			u.Vorruecken(5)
			if err := aufruf(tc.ctx(), t, u, z); !errors.Is(err, tc.want) {
				t.Fatalf("%s, WarteBis(%d) bei Jetzt 5: %v, erwartet %v", tc.name, z, err, tc.want)
			}
			if u.Wartende() != 0 {
				t.Fatalf("%s, WarteBis(%d): %d Wartende", tc.name, z, u.Wartende())
			}
		}
	}
}

// Endet der Kontext während des Wartens, kehrt der Aufruf mit dessen Fehler
// zurück und wartet nicht mehr.
func TestUhrAbbruchWaehrendDesWartens(t *testing.T) {
	u := testuhr.Neu()
	ctx, cancel := context.WithCancel(context.Background())
	ergebnis := wartet(ctx, t, u, 5)
	cancel()
	if err := ergebnisBinnen(t, ergebnis, "WarteBis nach Abbruch"); !errors.Is(err, context.Canceled) {
		t.Fatalf("WarteBis nach Abbruch: %v", err)
	}
	if u.Wartende() != 0 {
		t.Fatalf("Wartende nach Abbruch: %d", u.Wartende())
	}

	frist, cancelFrist := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancelFrist()
	if err := ergebnisBinnen(t, wartet(frist, t, u, 5), "WarteBis nach Fristablauf"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WarteBis nach Fristablauf: %v", err)
	}
}

// Vorruecken um eine negative Zahl ist ein Programmierfehler der Tests: Panik,
// und die Uhr bleibt stehen.
func TestUhrVorrueckenNegativPanik(t *testing.T) {
	u := testuhr.Neu()
	u.Vorruecken(5)
	func() {
		defer func() {
			if recover() == nil {
				t.Error("Vorruecken(-1) löst keine Panik aus")
			}
		}()
		u.Vorruecken(-1)
	}()
	if u.Jetzt() != 5 {
		t.Fatalf("Jetzt nach Vorruecken(-1): %d", u.Jetzt())
	}
}
