package services_test

import (
	"context"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// describeS1 ist das Describe des Statements s1 aus vorbereitung.
func describeS1() model.ClientMessage {
	return model.ClientMessage{Type: model.ClientDescribe, Target: model.TargetStatement, Name: "s1"}
}

// Abdeckung: LH-FA-13/Boundary — Replay im Kern: meldet der Adapter das
// Zwangsende (Forced), ist eine begonnene, nicht verbrauchte Interaktion
// PGR-E4006 der Klasse 4 mit der Kennung der zugeordneten Session und der
// aufgezeichneten Nummer der Interaktion: eine Extended-Interaktion vor dem
// Sync ihrer letzten Gruppe, auch nach beantworteten Gruppen, und eine
// Interaktion, deren Antworten nicht als gesendet gemeldet sind, einfach oder
// Extended. Ohne solche Interaktion und ohne zugeordnete Session liefert Forced
// nil; es ändert keinen Zustand (LH-FA-13.a *Zwangsende*, LH-FA-03.b
// *Verbraucht*).
func TestReplayZwangsende(t *testing.T) {
	sitzung := []model.Interaction{interaktion(1, "SELECT 1", "SELECT 1"), vorbereitung(3)}
	cases := []struct {
		name   string
		ablauf func(t *testing.T, s *services.ReplayService, id model.SessionID)
		want   string
	}{
		{"ohne Session", func(*testing.T, *services.ReplayService, model.SessionID) {}, ""},
		{"verbraucht", func(t *testing.T, s *services.ReplayService, id model.SessionID) {
			if _, err := s.Query(context.Background(), id, "SELECT 1"); err != nil {
				t.Fatal(err)
			}
			s.Sent(context.Background(), id)
		}, ""},
		{"einfache Interaktion nicht gesendet", func(t *testing.T, s *services.ReplayService, id model.SessionID) {
			if _, err := s.Query(context.Background(), id, "SELECT 1"); err != nil {
				t.Fatal(err)
			}
		}, "Session 1, Interaktion 1 nicht verbraucht"},
		{"Extended vor der ersten Gruppe beantwortet", func(t *testing.T, s *services.ReplayService, id model.SessionID) {
			if _, err := s.Query(context.Background(), id, "SELECT 1"); err != nil {
				t.Fatal(err)
			}
			s.Sent(context.Background(), id)
			sendeAlle(t, s, id, parse("s1", "SELECT $1::text"))
		}, "Session 1, Interaktion 3 nicht verbraucht"},
		{"Extended nach beantworteter Gruppe", func(t *testing.T, s *services.ReplayService, id model.SessionID) {
			if _, err := s.Query(context.Background(), id, "SELECT 1"); err != nil {
				t.Fatal(err)
			}
			s.Sent(context.Background(), id)
			sendeAlle(t, s, id, parse("s1", "SELECT $1::text"), describeS1(), flushNachricht())
			s.Sent(context.Background(), id)
		}, "Session 1, Interaktion 3 nicht verbraucht"},
		{"Extended beantwortet, nicht gesendet", func(t *testing.T, s *services.ReplayService, id model.SessionID) {
			if _, err := s.Query(context.Background(), id, "SELECT 1"); err != nil {
				t.Fatal(err)
			}
			s.Sent(context.Background(), id)
			sendeAlle(t, s, id, parse("s1", "SELECT $1::text"), describeS1(), flushNachricht())
			s.Sent(context.Background(), id)
			sendeAlle(t, s, id, bind("s1", text("a")), execute(), syncNachricht())
		}, "Session 1, Interaktion 3 nicht verbraucht"},
		{"Extended verbraucht", func(t *testing.T, s *services.ReplayService, id model.SessionID) {
			if _, err := s.Query(context.Background(), id, "SELECT 1"); err != nil {
				t.Fatal(err)
			}
			s.Sent(context.Background(), id)
			sendeAlle(t, s, id, parse("s1", "SELECT $1::text"), describeS1(), flushNachricht())
			s.Sent(context.Background(), id)
			sendeAlle(t, s, id, bind("s1", text("a")), execute(), syncNachricht())
			s.Sent(context.Background(), id)
		}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, id := replayMit(t, sitzung)
			c.ablauf(t, s, id)
			for i := 0; i < 2; i++ {
				err := s.Forced(context.Background(), id)
				if c.want == "" {
					if err != nil {
						t.Fatalf("Aufruf %d: Fehler %v ohne begonnene Interaktion", i+1, err)
					}
					continue
				}
				ms := model.Meldungen(err)
				if len(ms) != 1 || ms[0].Code != model.CodeShutdownTimeout || ms[0].ExitCode() != 4 ||
					!strings.HasPrefix(ms[0].Text, "Netzwerk [PGR-E4006]: ") || !strings.Contains(ms[0].Text, c.want) {
					t.Fatalf("Aufruf %d: Meldungen %#v, erwartet PGR-E4006 mit %q", i+1, ms, c.want)
				}
			}
		})
	}
}

// Abdeckung: LH-FA-13/Negative — Replay im Kern: mit --fail-on-unconsumed
// liefert das Zwangsende PGR-E4006 und das Ende derselben Verbindung danach
// PGR-E5002 für die nicht verbrauchten Interaktionen, mit der unvollständigen
// als erster; der Adapter merkt PGR-E4006 zuerst (LH-FA-03.b *Fehlerebene*).
func TestReplayZwangsendeNichtVerbraucht(t *testing.T) {
	s, err := services.NewReplayService(context.Background(), ladeRepo{rec: aufzeichnung([]model.Interaction{interaktion(1, "SELECT 1", "SELECT 1"), vorbereitung(2)})}, "rec.yaml", services.FailOnUnconsumed)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := s.OpenConnection(context.Background())
	if _, err := s.Query(context.Background(), id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	s.Sent(context.Background(), id)
	sendeAlle(t, s, id, parse("s1", "SELECT $1::text"))
	if err := s.Forced(context.Background(), id); codeOf(err) != model.CodeShutdownTimeout {
		t.Fatalf("Forced: %v", err)
	}
	w, err := s.CloseConnection(context.Background(), id)
	if w != nil || codeOf(err) != model.CodeReplayUnconsumed || !strings.Contains(err.Error(), "Session 1: 1 von 2 Interaktionen nicht verbraucht, die erste mit Nummer 2") {
		t.Fatalf("CloseConnection: %v, %v", w, err)
	}
}
