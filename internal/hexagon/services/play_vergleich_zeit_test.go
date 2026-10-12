package services_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
	"github.com/pt9912/pgwire-recorder/internal/testuhr"
)

// Abdeckung: LH-FA-21/Boundary — Mit --keep-timing und --compare-responses zugleich
// ändert der Vergleich die Wartezeit nicht: Die zweite Interaktion wartet bis zum
// aufgezeichneten Abstand (999 ms davor noch nicht, 1000 ms gesendet), und die Antwort
// wird nach der Interaktion verglichen (die Abweichung der zweiten ist PGR-E5004, beide
// Interaktionen zählen als verglichen) (LH-FA-21.a *Mit `--compare-responses`*).
func TestPlayZeitMitVergleich(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	gut := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	sitzungA := spielerFolge(false, gut, []model.Response{befehl("UPDATE 1"), bereitMit("I")})
	ziel := &fakeZiel{sessions: []*fakeEinspiel{sitzungA}}
	s := model.Session{ID: 1, Interactions: []model.Interaction{aufgezeichnet(1, "A", gut...), aufgezeichnet(2, "B", gut...)}}
	for i, off := range []int64{0, 1000} {
		s.Interactions[i].OffsetMS = &off
	}
	o := relativ()
	o.CompareResponses = true
	rec := model.NewRecording()
	rec.Sessions = []model.Session{s}
	svc, err := services.NewPlayService(context.Background(), ladeRepo{rec: rec}, "r.yaml", ziel, o, services.Wartezeiten(uhr))
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	l := &zeitLauf{t: t, uhr: uhr, ziel: ziel, fertig: make(chan error, 1)}
	go func() { l.fertig <- svc.Play(context.Background(), nil) }()

	l.wartet(1)
	l.wartetNoch(1, "anfrage A")
	uhr.Vorruecken(999)
	l.wartetNoch(1, "anfrage A")
	uhr.Vorruecken(1)
	err = l.endet()
	if got, want := ziel.gesendet(), []string{"anfrage A", "anfrage B"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("gesendet %q, erwartet %q", got, want)
	}
	if got, want := texte(err), []string{e5004 + "Session 1, Interaktion 2: Antwort weicht von der Aufzeichnung ab: Position 1, Art Befehl"}; !reflect.DeepEqual(got, want) || exitVon(err) != 5 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 5", got, exitVon(err), want)
	}
	if want := (services.Zaehler{Eingespielt: 2, Verglichen: 2, Abweichend: 1}); svc.Zaehler() != want {
		t.Fatalf("Zähler %+v, erwartet %+v", svc.Zaehler(), want)
	}
}
