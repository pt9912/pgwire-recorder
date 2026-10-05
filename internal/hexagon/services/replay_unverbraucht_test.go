package services

import (
	"context"
	"errors"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// ladenMit lädt die Sessions mit den Optionen.
func ladenMit(t *testing.T, opts []ReplayOption, sessions ...[]model.Interaction) *ReplayService {
	t.Helper()
	s, err := NewReplayService(context.Background(), ladeRepo{rec: aufzeichnung(sessions...)}, "rec.yaml", opts...)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

var streng = []ReplayOption{FailOnUnconsumed}

// meldungBeimSchliessen beendet die Verbindung und liefert Code, Text und
// Exit-Code ihrer Meldung (zerlege).
func meldungBeimSchliessen(t *testing.T, s *ReplayService, id model.SessionID) (string, string, int) {
	t.Helper()
	w, err := s.CloseConnection(context.Background(), id)
	return zerlege(t, w, err)
}

// zerlege liefert Code, Text und Exit-Code einer Meldung, ohne Meldung "";
// Warnung und Fehler zugleich, ein Fehler ohne Meldungscode oder ein
// Fehlertext ohne den Kopf „Replay [<code>]: “ vor dem Text brechen ab.
func zerlege(t *testing.T, w *model.Warning, err error) (string, string, int) {
	t.Helper()
	switch {
	case w != nil && err != nil:
		t.Fatalf("Warnung und Fehler zugleich: %v, %v", w, err)
	case w != nil:
		return w.Code, w.Msg, 0
	case err != nil:
		var me *model.Error
		if !errors.As(err, &me) {
			t.Fatalf("Fehler ohne Meldungscode: %v", err)
		}
		if want := "Replay [" + me.Code + "]: " + me.Msg; err.Error() != want {
			t.Fatalf("Fehlertext %q, erwartet Kopf und Text %q", err.Error(), want)
		}
		return me.Code, me.Msg, me.ExitCode()
	}
	return "", "", 0
}

// Abdeckung: LH-FA-13/Negative — mit FailOnUnconsumed ist eine Session, deren
// Verbindung vor dem Verbrauch aller Interaktionen endet, der Fehler PGR-E5002
// der Klasse 5. Dazu: Endet eine Verbindung vor dem Verbrauch aller
// Interaktionen ihrer Session, nennt die Meldung die Session, die Zahl der
// nicht verbrauchten und aller Interaktionen ohne aufgezeichnete
// Lebendprüfungen und die aufgezeichnete Nummer der ersten nicht verbrauchten;
// ohne FailOnUnconsumed ist sie die Warnung PGR-W2001, mit ihr der Fehler
// PGR-E5002 der Klasse 5 mit demselben Text.
func TestReplayNichtVerbrauchtMeldung(t *testing.T) {
	ctx := context.Background()
	session := []model.Interaction{interaktion(1, "-- ping", ""), interaktion(2, "SELECT 1", "A"), interaktion(3, "SELECT 2", "B"), interaktion(4, "SELECT 3", "C")}
	const text = "Session 1: 2 von 3 Interaktionen nicht verbraucht, die erste mit Nummer 3"
	for _, fall := range []struct {
		opts       []ReplayOption
		code       string
		exit       int
		beschreibt string
	}{
		{nil, model.CodeUnconsumed, 0, "ohne Option"},
		{streng, model.CodeReplayUnconsumed, 5, "mit Option"},
	} {
		s := ladenMit(t, fall.opts, session)
		id, _ := s.OpenConnection(ctx)
		if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
			t.Fatal(err)
		}
		s.Sent(ctx, id)
		c, msg, exit := meldungBeimSchliessen(t, s, id)
		if c != fall.code || msg != text || exit != fall.exit {
			t.Fatalf("%s: %s %q Exit %d, erwartet %s %q Exit %d", fall.beschreibt, c, msg, exit, fall.code, text, fall.exit)
		}
	}
}

// Verbraucht ist eine einfache Interaktion erst,
// wenn ihre Antworten als gesendet gemeldet sind, eine Extended-Interaktion
// erst mit den gesendeten Antworten ihrer letzten Gruppe; eine vor dem Sync
// ihrer letzten Gruppe beendete Extended-Interaktion ist nicht verbraucht,
// auch wenn eine Gruppe schon beantwortet ist.
func TestReplayVerbraucht(t *testing.T) {
	ctx := context.Background()

	s := ladenMit(t, streng, []model.Interaction{interaktion(1, "SELECT 1", "A")})
	id, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if c, msg, _ := meldungBeimSchliessen(t, s, id); c != model.CodeReplayUnconsumed || msg != "Session 1: 1 von 1 Interaktionen nicht verbraucht, die erste mit Nummer 1" {
		t.Fatalf("Antworten nicht gesendet: %s %q", c, msg)
	}

	s = ladenMit(t, streng, []model.Interaction{interaktion(1, "SELECT 1", "A"), vorbereitung(2)})
	id, _ = s.OpenConnection(ctx)
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	s.Sent(ctx, id)
	in := vorbereitung(2)
	sendeAlle(t, s, id, in.Groups[0].Client...)
	s.Sent(ctx, id)
	if c, msg, _ := meldungBeimSchliessen(t, s, id); c != model.CodeReplayUnconsumed || msg != "Session 1: 1 von 2 Interaktionen nicht verbraucht, die erste mit Nummer 2" {
		t.Fatalf("vor dem Sync der letzten Gruppe beendet: %s %q", c, msg)
	}

	s = ladenMit(t, streng, []model.Interaction{vorbereitung(1)})
	id, _ = s.OpenConnection(ctx)
	sendeAlle(t, s, id, in.Groups[0].Client...)
	s.Sent(ctx, id)
	sendeAlle(t, s, id, in.Groups[1].Client...)
	if c, msg, _ := meldungBeimSchliessen(t, s, id); c != model.CodeReplayUnconsumed || msg != "Session 1: 1 von 1 Interaktionen nicht verbraucht, die erste mit Nummer 1" {
		t.Fatalf("Antworten der letzten Gruppe nicht gesendet: %s %q", c, msg)
	}

	s = ladenMit(t, streng, []model.Interaction{vorbereitung(1)})
	id, _ = s.OpenConnection(ctx)
	sendeAlle(t, s, id, in.Groups[0].Client...)
	s.Sent(ctx, id)
	sendeAlle(t, s, id, in.Groups[1].Client...)
	s.Sent(ctx, id)
	if c, msg, _ := meldungBeimSchliessen(t, s, id); c != "" {
		t.Fatalf("verbrauchte Session gemeldet: %s %q", c, msg)
	}
}

// Abdeckung: LH-FA-13/Negative — nie zugeordnete Sessions mit Interaktionen
// ergeben eine Meldung mit ihrer Zahl und der Kennung der ersten: ohne
// FailOnUnconsumed die Warnung PGR-W2001, mit ihr der Fehler PGR-E5002 mit
// demselben Text.
func TestReplayNieZugeordnetMeldung(t *testing.T) {
	ctx := context.Background()
	const text = "2 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 2"
	for _, fall := range []struct {
		opts []ReplayOption
		code string
		exit int
	}{
		{nil, model.CodeUnconsumed, 0},
		{streng, model.CodeReplayUnconsumed, 5},
	} {
		s := ladenMit(t, fall.opts,
			[]model.Interaction{interaktion(1, "S1", "A")},
			[]model.Interaction{interaktion(1, "S2", "B")},
			[]model.Interaction{interaktion(1, "S3", "C")},
		)
		id, _ := s.OpenConnection(ctx)
		if _, err := s.Query(ctx, id, "S1"); err != nil {
			t.Fatal(err)
		}
		s.Sent(ctx, id)
		if c, _, _ := meldungBeimSchliessen(t, s, id); c != "" {
			t.Fatalf("verbrauchte Session gemeldet: %s", c)
		}
		w, err := s.Unassigned()
		if c, msg, exit := zerlege(t, w, err); c != fall.code || msg != text || exit != fall.exit {
			t.Fatalf("%s %q Exit %d, erwartet %s %q Exit %d", c, msg, exit, fall.code, text, fall.exit)
		}
	}
}

// Auch mit FailOnUnconsumed meldet eine
// Verbindung ohne zugeordnete Session nichts, auch nach Lebendprüfungen, und
// eine Session ohne Interaktion oder nur aus aufgezeichneten Lebendprüfungen
// ist weder nicht verbraucht noch nie zugeordnet.
func TestReplayNichtsZuMelden(t *testing.T) {
	ctx := context.Background()
	s := ladenMit(t, streng,
		[]model.Interaction{interaktion(1, "SELECT 1", "A")},
		nil,
		[]model.Interaction{interaktion(1, "-- ping", ""), interaktion(2, "", "")},
	)
	leer, _ := s.OpenConnection(ctx)
	lebend, _ := s.OpenConnection(ctx)
	pruefeLebend(t, s, lebend, "-- ping", "I")
	s.Sent(ctx, lebend)
	for _, id := range []model.SessionID{leer, lebend} {
		if c, msg, _ := meldungBeimSchliessen(t, s, id); c != "" {
			t.Fatalf("Verbindung %d ohne Session meldet %s %q", id, c, msg)
		}
	}
	id, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	s.Sent(ctx, id)
	if c, msg, _ := meldungBeimSchliessen(t, s, id); c != "" {
		t.Fatalf("verbrauchte Session meldet %s %q", c, msg)
	}
	w, err := s.Unassigned()
	if c, msg, _ := zerlege(t, w, err); c != "" {
		t.Fatalf("Sessions ohne Interaktion gelten als nie zugeordnet: %s %q", c, msg)
	}
}

// Sent meldet nur die Antworten der eigenen Verbindung als gesendet: Meldet
// eine zweite Verbindung Sent, bleibt die zuletzt gelieferte Interaktion der
// ersten nicht verbraucht.
func TestReplaySentJeVerbindung(t *testing.T) {
	ctx := context.Background()
	s := ladenMit(t, streng,
		[]model.Interaction{interaktion(1, "S1", "A")},
		[]model.Interaction{interaktion(1, "S2", "B")},
	)
	a, _ := s.OpenConnection(ctx)
	b, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, a, "S1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Query(ctx, b, "S2"); err != nil {
		t.Fatal(err)
	}
	s.Sent(ctx, b)
	if c, msg, _ := meldungBeimSchliessen(t, s, a); c != model.CodeReplayUnconsumed || msg != "Session 1: 1 von 1 Interaktionen nicht verbraucht, die erste mit Nummer 1" {
		t.Fatalf("Sent der zweiten Verbindung hat die erste verbraucht: %s %q", c, msg)
	}
	if c, msg, _ := meldungBeimSchliessen(t, s, b); c != "" {
		t.Fatalf("gesendete Session gemeldet: %s %q", c, msg)
	}
}
