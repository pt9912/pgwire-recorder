package bootstrap_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// mitOffset ist die Interaktion in mit offset_ms ms.
func mitOffset(in model.Interaction, ms int64) model.Interaction {
	in.OffsetMS = &ms
	return in
}

// Abdeckung: LH-FA-21/Happy — play verdrahtet die Systemuhr: Mit --keep-timing liegen
// zwei Interaktionen mit offset_ms 0 und 400 mindestens 350 ms Echtzeit auseinander (die
// Untergrenze der Pause), Exit-Code 0 (LH-FA-21.a *Einspielen*, *Uhr*).
func TestRunPlayKeepTiming(t *testing.T) {
	leerePlay(t)
	addr, typen := playServer(t, append(nachricht('C', []byte("SELECT 0\x00")...), nachricht('Z', 'I')...))
	input := einfacheAufzeichnung(t, []model.Interaction{mitOffset(anfrage(1, "A"), 0), mitOffset(anfrage(2, "B"), 400)})
	var stdout, stderr bytes.Buffer
	beginn := time.Now()
	code := bootstrap.Run(context.Background(), nil, []string{"play", "--upstream", addr, "--input", input, "--keep-timing"}, "dev", &stdout, &stderr)
	dauer := time.Since(beginn)
	if code != 0 || stdout.Len() > 0 {
		t.Fatalf("Exit-Code %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
	if got := typenBis(t, typen, 'X'); got != "QQX" {
		t.Fatalf("Client-Nachrichten %q", got)
	}
	if dauer < 350*time.Millisecond {
		t.Fatalf("play dauerte %v, erwartet mindestens 350 ms Pause", dauer)
	}
}

// Abdeckung: LH-FA-21/Negative — fehlt mit --keep-timing einer Interaktion offset_ms, ist
// das ein Startfehler PGR-E2003 mit Exit-Code 2, die Zeile beim Prozessende ohne Log-Zeile
// davor, und play verbindet sich nicht (LH-FA-21.a *Randformen des Einspielens*).
func TestRunPlayKeepTimingOhneOffset(t *testing.T) {
	leerePlay(t)
	addr, typen := playServer(t, nil)
	input := einfacheAufzeichnung(t, []model.Interaction{mitOffset(anfrage(1, "A"), 0), anfrage(2, "B")})
	var stdout, stderr bytes.Buffer
	code := bootstrap.Run(context.Background(), nil, []string{"play", "--upstream", addr, "--input", input, "--keep-timing"}, "dev", &stdout, &stderr)
	if code != 2 || stdout.Len() > 0 || stderr.String() != "Konfiguration [PGR-E2003]: Session 1, Interaktion 2 trägt kein offset_ms, das --keep-timing an jeder Interaktion verlangt\n" {
		t.Fatalf("Exit-Code %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
	select {
	case b := <-typen:
		t.Fatalf("Startfehler verbindet sich: %q", b)
	default:
	}
}
