//go:build unix

package bootstrap_test

import (
	"context"
	"errors"
	"math"
	"syscall"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
)

// prozessCPU liefert die bisher verbrauchte CPU-Zeit des Prozesses.
func prozessCPU(t *testing.T) time.Duration {
	t.Helper()
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		t.Fatal(err)
	}
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

// Das Warten auf einen fernen Zeitpunkt, dessen Abstand die größte Dauer der
// Zeitbibliothek übersteigt, verbraucht keine CPU-Zeit: es schläft, statt zu
// drehen.
func TestSystemuhrFernerZeitpunktDrehtNicht(t *testing.T) {
	u := bootstrap.NeueSystemuhr()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	vorher := prozessCPU(t)
	err := wartenBinnen(ctx, t, u, math.MaxInt64, 5*time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WarteBis: %v", err)
	}
	if cpu := prozessCPU(t) - vorher; cpu > 150*time.Millisecond {
		t.Fatalf("das Warten von 300 ms verbrauchte %v CPU-Zeit", cpu)
	}
}
