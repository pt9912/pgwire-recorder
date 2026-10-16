package testzeit_test

import (
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/testzeit"
)

// Nach und Schlafe zählen in ganzen Millisekunden und kehren nicht früher zurück.
func TestWartenInMillisekunden(t *testing.T) {
	start := time.Now()
	testzeit.Schlafe(30)
	if d := time.Since(start); d < 30*time.Millisecond || d > 5*time.Second {
		t.Fatalf("Schlafe(30) dauerte %v", d)
	}
	start = time.Now()
	<-testzeit.Nach(30)
	if d := time.Since(start); d < 30*time.Millisecond || d > 5*time.Second {
		t.Fatalf("Nach(30) lieferte nach %v", d)
	}
}

// Eine Frist ist erst nach ihren Millisekunden abgelaufen.
func TestFristLaeuftNachMillisekundenAb(t *testing.T) {
	f := testzeit.NeueFrist(60)
	if f.Abgelaufen() {
		t.Fatal("die Frist ist sofort abgelaufen")
	}
	testzeit.Schlafe(80)
	if !f.Abgelaufen() {
		t.Fatal("die Frist ist nach 80 ms nicht abgelaufen")
	}
}
