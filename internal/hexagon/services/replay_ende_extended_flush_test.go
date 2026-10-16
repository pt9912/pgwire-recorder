package services_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// endeFlushNachrichten sind die Client-Nachrichten von extendedEndeFlush mit
// davor: eine Gruppe mit Flush (parse, flush) und die Gruppe mit Ende, die mit
// Flush schließt (bind, execute, flush).
func endeFlushNachrichten() []model.ClientMessage {
	n := endeNachrichten()
	n[4] = flushNachricht()
	return n
}

// extendedEndeFlush ist extendedEnde mit einer Gruppe mit Ende, die mit Flush
// statt Sync schließt; ihre Server-Nachrichten sind dieselben.
func extendedEndeFlush(seq int, sev string, davor bool) model.Interaction {
	in := extendedEnde(seq, sev, davor)
	g := &in.Groups[len(in.Groups)-1]
	g.Client = endeFlushNachrichten()[2:]
	return in
}

// sendeFlushBis sendet die ersten n der endeFlushNachrichten, meldet die Gruppe
// davor als gesendet und liefert die Antworten der letzten Nachricht; die Gruppe
// mit Ende meldet der Aufrufer.
func sendeFlushBis(t *testing.T, s *services.ReplayService, id model.SessionID, n int) (out []model.Response) {
	t.Helper()
	for i, m := range endeFlushNachrichten()[:n] {
		var err error
		if out, err = s.ClientMessage(context.Background(), id, m); err != nil {
			t.Fatalf("Nachricht %s: %v", m.Type, err)
		}
		if i == 1 {
			s.Sent(context.Background(), id)
		}
	}
	return out
}

// Abdeckung: LH-FA-25/Happy, LH-FA-18/Boundary — Replay: Die Gruppe mit Ende, die
// mit Flush schließt, liefert ihre Server-Nachrichten in der aufgezeichneten
// Reihenfolge, die letzte ist connection_closed, ohne ready_for_query, erst mit
// dem abschließenden Flush, auch als einzige Gruppe (PANIC); die Gruppe mit Flush
// davor liefert kein connection_closed (LH-FA-25.a *Stelle von connection_closed
// im Replay*).
func TestReplayEndeExtendedFlushGruppeMitEnde(t *testing.T) {
	for _, sev := range []string{"FATAL", "PANIC"} {
		s, id := endeVerbindung(t, []model.Interaction{extendedEndeFlush(1, sev, sev == "FATAL")})
		if sev == "FATAL" {
			if out := sendeFlushBis(t, s, id, 2); !reflect.DeepEqual(out, []model.Response{{Type: model.ResponseParseComplete}}) {
				t.Fatalf("Gruppe mit Flush davor: %#v", out)
			}
		}
		for _, m := range endeFlushNachrichten()[2:4] {
			if out, err := s.ClientMessage(context.Background(), id, m); err != nil || out != nil {
				t.Fatalf("%s: Antwort vor dem abschließenden Flush: %#v, %v", sev, out, err)
			}
		}
		out, err := s.ClientMessage(context.Background(), id, flushNachricht())
		if want := endeGruppeServer(sev); err != nil || !reflect.DeepEqual(out, want) {
			t.Fatalf("%s: erwartet %#v, erhalten %#v, %v", sev, want, out, err)
		}
		if out[len(out)-1].Type != model.ResponseConnectionClosed || zahl(out, model.ResponseReadyForQuery) != 0 {
			t.Fatalf("%s: connection_closed nicht letzte oder ready_for_query dabei: %#v", sev, out)
		}
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-03/Boundary — Replay: Die Interaktion mit
// Ende, die mit Flush schließt, ist mit Sent verbraucht; ohne Sent nicht:
// PGR-W2001, mit FailOnUnconsumed PGR-E5002, beim Zwangsende PGR-E4006 mit ihrer
// Nummer.
func TestReplayEndeExtendedFlushVerbrauch(t *testing.T) {
	ctx := context.Background()
	in := []model.Interaction{extendedEndeFlush(2, "FATAL", true)}
	s, id := endeVerbindung(t, in)
	sendeFlushBis(t, s, id, 5)
	s.Sent(ctx, id)
	if err := s.Forced(ctx, id); err != nil {
		t.Fatalf("Zwangsende nach Sent: %v", err)
	}
	if w := replaySchliessen(t, s, id); w != nil {
		t.Fatalf("gesendetes Ende gilt als unverbraucht: %v", w)
	}
	for _, opts := range [][]services.ReplayOption{nil, {services.FailOnUnconsumed}} {
		s, id = endeVerbindung(t, in, opts...)
		sendeFlushBis(t, s, id, 5)
		if err := s.Forced(ctx, id); code(err) != model.CodeShutdownTimeout || !strings.Contains(err.Error(), "Interaktion 2 nicht verbraucht") {
			t.Fatalf("Zwangsende vor Sent: %v", err)
		}
		w, err := s.CloseConnection(ctx, id)
		if opts != nil && code(err) != model.CodeReplayUnconsumed || opts == nil && (err != nil || w == nil || w.Code != model.CodeUnconsumed || !strings.Contains(w.Msg, "1 von 1")) {
			t.Fatalf("nicht gesendet (FailOnUnconsumed %v): %v, %v", opts != nil, w, err)
		}
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-10/Negative — Replay: Eine Nachricht, die
// der erwarteten nicht gleicht, ist PGR-E5001 und liefert kein Ende, an jeder
// Stelle (eine Zeile je Stelle, auch Sync statt Flush, einfache Anfrage und
// Lebendprüfung mitten in und zwischen den Gruppen); der Cursor bleibt stehen, die
// richtigen Nachrichten danach liefern das Ende noch (LH-FA-25.a *Reihenfolge der
// Gründe im Replay*, LH-FA-18.a).
func TestReplayEndeExtendedFlushAbweichung(t *testing.T) {
	ctx := context.Background()
	falschesBind := bind("s1", text("a"))
	falschesBind.Statement = "zz"
	nachricht := func(m model.ClientMessage) func(*services.ReplayService, model.SessionID) ([]model.Response, error) {
		return func(s *services.ReplayService, id model.SessionID) ([]model.Response, error) {
			return s.ClientMessage(ctx, id, m)
		}
	}
	anfrage := func(sql string) func(*services.ReplayService, model.SessionID) ([]model.Response, error) {
		return func(s *services.ReplayService, id model.SessionID) ([]model.Response, error) {
			return s.Query(ctx, id, sql)
		}
	}
	for _, f := range []struct {
		name   string
		vorher int
		falsch func(*services.ReplayService, model.SessionID) ([]model.Response, error)
	}{
		{"erste Nachricht der Interaktion", 0, nachricht(parse("s1", "SELECT 2"))},
		{"Sync statt Flush in der Gruppe davor", 1, nachricht(syncNachricht())},
		{"erste Nachricht der Gruppe mit Ende", 2, nachricht(falschesBind)},
		{"mittlere Nachricht der Gruppe mit Ende", 3, nachricht(model.ClientMessage{Type: model.ClientExecute, Portal: "x"})},
		{"Sync statt Flush als abschließende", 4, nachricht(syncNachricht())},
		{"einfache Anfrage vor den Gruppen", 0, anfrage("SELECT 1")},
		{"einfache Anfrage nach der Gruppe davor", 2, anfrage("SELECT 1")},
		{"einfache Anfrage mitten in der Gruppe mit Ende", 3, anfrage("SELECT 1")},
		{"Lebendprüfung mitten in der Gruppe davor", 1, anfrage("-- ping")},
		{"Lebendprüfung nach der Gruppe davor", 2, anfrage("-- ping")},
		{"Lebendprüfung mitten in der Gruppe mit Ende", 3, anfrage("-- ping")},
	} {
		s, id := endeVerbindung(t, []model.Interaction{extendedEndeFlush(1, "FATAL", true)})
		sendeFlushBis(t, s, id, f.vorher)
		if out, err := f.falsch(s, id); code(err) != model.CodeReplayMismatch || out != nil {
			t.Fatalf("%s: erwartet %s ohne Antworten, erhalten %#v, %v", f.name, model.CodeReplayMismatch, out, err)
		}
		var out []model.Response
		for _, m := range endeFlushNachrichten()[f.vorher:] {
			out, _ = s.ClientMessage(ctx, id, m)
		}
		if !reflect.DeepEqual(out, endeGruppeServer("FATAL")) {
			t.Fatalf("%s: das Ende kommt nach der Abweichung noch: %#v", f.name, out)
		}
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-03/Boundary — Replay: `Terminate` und das
// Ende der Client-Verbindung (für den Service CloseConnection) vor, in und nach
// der Gruppe davor und mitten in der Gruppe mit Ende, die mit Flush schließt,
// lassen die Interaktion unverbraucht: PGR-W2001, mit FailOnUnconsumed PGR-E5002,
// je mit der Nummer der Interaktion (LH-FA-03.b *Zeitpunkte* 1).
func TestReplayEndeExtendedFlushTerminateUnverbraucht(t *testing.T) {
	ins := []model.Interaction{ausfuehrung(1, "s1", "a", "A"), extendedEndeFlush(7, "FATAL", true)}
	for n := 0; n <= 4; n++ {
		for _, opts := range [][]services.ReplayOption{nil, {services.FailOnUnconsumed}} {
			s, id := endeVerbindung(t, ins, opts...)
			sendeAlle(t, s, id, ins[0].Groups[0].Client...)
			s.Sent(context.Background(), id)
			sendeFlushBis(t, s, id, n)
			w, err := s.CloseConnection(context.Background(), id)
			if opts != nil && (code(err) != model.CodeReplayUnconsumed || !strings.Contains(err.Error(), "Nummer 7")) ||
				opts == nil && (err != nil || w == nil || w.Code != model.CodeUnconsumed || !strings.Contains(w.Msg, "1 von 2") || !strings.Contains(w.Msg, "Nummer 7")) {
				t.Fatalf("nach %d Nachrichten (FailOnUnconsumed %v): erwartet 1 von 2, Nummer 7: %v, %v", n, opts != nil, w, err)
			}
		}
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-13/Boundary — Replay: Das Herunterfahren
// gilt wie sonst (LH-FA-13.a): Mit dem Cursor in der Interaktion mit Ende, die mit
// Flush schließt, erlaubt es Shutdown nicht, davor und nach dem Senden erlaubt es;
// das Zwangsende ist bis zum Senden der Gruppe mit Ende PGR-E4006 mit der Nummer.
func TestReplayEndeExtendedFlushHerunterfahren(t *testing.T) {
	ctx := context.Background()
	for i, w := range []struct{ shutdown, zwang bool }{{true, false}, {false, true}, {false, true}, {false, true}, {false, true}, {true, true}, {true, false}} {
		s, id := endeVerbindung(t, []model.Interaction{extendedEndeFlush(7, "FATAL", true)})
		sendeFlushBis(t, s, id, min(i, 5))
		if i == 6 {
			s.Sent(ctx, id)
		}
		err := s.Forced(ctx, id)
		if got := s.Shutdown(ctx, id); got != w.shutdown || w.zwang != (code(err) == model.CodeShutdownTimeout && strings.Contains(err.Error(), "Interaktion 7 nicht verbraucht")) {
			t.Fatalf("Stand %d: Shutdown %v (erwartet %v), Zwangsende %v (erwartet PGR-E4006 %v)", i, got, w.shutdown, err, w.zwang)
		}
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-12/Boundary — Replay mit `first-request`:
// Die Interaktion mit Ende, die mit Flush schließt, gehört der nächsten freien
// Session; die nächste Verbindung erhält die nächste Session mit ihrem Ende; ohne
// freie Session ist es PGR-E5003, nie das Ende (LH-FA-12.a, LH-FA-25.a
// *Reihenfolge der Gründe im Replay*).
func TestReplayEndeExtendedFlushFirstRequest(t *testing.T) {
	ctx := context.Background()
	s := ladeEnde(t, endeAufzeichnung([]model.Interaction{extendedEndeFlush(1, "FATAL", false)}, []model.Interaction{extendedEndeFlush(1, "PANIC", false)}))
	a, _ := s.OpenConnection(ctx)
	if out := sendeAlle(t, s, a, endeFlushNachrichten()[2:]...); !reflect.DeepEqual(out, endeGruppeServer("FATAL")) {
		t.Fatalf("erste Session: %#v", out)
	}
	if w := unassigned(t, s); w == nil || !strings.Contains(w.Msg, "1 aufgezeichnete") || !strings.Contains(w.Msg, "Kennung 2") {
		t.Fatalf("erste Session ist zugeordnet, die zweite nicht: %v", w)
	}
	if w := gesendetSchliessen(t, s, a); w != nil {
		t.Fatalf("Session 1: %v", w)
	}
	b, _ := s.OpenConnection(ctx)
	if out := sendeAlle(t, s, b, endeFlushNachrichten()[2:]...); !reflect.DeepEqual(out, endeGruppeServer("PANIC")) {
		t.Fatalf("zweite Session: %#v", out)
	}
	c, _ := s.OpenConnection(ctx)
	for _, m := range []model.ClientMessage{endeFlushNachrichten()[2], flushNachricht()} {
		if out, err := s.ClientMessage(ctx, c, m); code(err) != model.CodeReplaySession || out != nil {
			t.Fatalf("%s ohne freie Session: erwartet %s, erhalten %#v, %v", m.Type, model.CodeReplaySession, out, err)
		}
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-12/Boundary — Replay mit `connection`: Die
// n-te Verbindung erhält die Session n mit ihrem Ende, das mit Flush schließt; das
// Schließen ändert die Zuordnung der nächsten nicht; ohne Session ist es PGR-E5003
// (LH-FA-12.a).
func TestReplayEndeExtendedFlushConnection(t *testing.T) {
	ctx := context.Background()
	rec := gekennzeichnet(sessionNr(1, "17.0", extendedEndeFlush(1, "FATAL", false)), sessionNr(2, "17.0", ausfuehrung(1, "s1", "a", "A")), sessionNr(3, "17.0", extendedEndeFlush(1, "PANIC", false)))
	rec.Version = model.VersionMitEnde
	s := nachVerbindung(t, rec)
	a, _ := s.OpenConnection(ctx)
	b, _ := s.OpenConnection(ctx)
	if out := sendeAlle(t, s, a, endeFlushNachrichten()[2:]...); !reflect.DeepEqual(out, endeGruppeServer("FATAL")) {
		t.Fatalf("Session 1 an Verbindung 1: %#v", out)
	}
	if w := gesendetSchliessen(t, s, a); w != nil {
		t.Fatalf("Session 1: %v", w)
	}
	if out := sendeAlle(t, s, b, endeNachrichten()[2:]...); zahl(out, model.ResponseConnectionClosed) != 0 || len(out) != 3 {
		t.Fatalf("Session 2 an Verbindung 2: %#v", out)
	}
	c, _ := s.OpenConnection(ctx)
	if out := sendeAlle(t, s, c, endeFlushNachrichten()[2:]...); !reflect.DeepEqual(out, endeGruppeServer("PANIC")) {
		t.Fatalf("Session 3 an Verbindung 3 nach dem Schließen der ersten: %#v", out)
	}
	d, _ := s.OpenConnection(ctx)
	if out, err := s.ClientMessage(ctx, d, endeFlushNachrichten()[2]); code(err) != model.CodeReplaySession || out != nil {
		t.Fatalf("Verbindung ohne Session: erwartet %s, erhalten %#v, %v", model.CodeReplaySession, out, err)
	}
}
