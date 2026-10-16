package services_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// endeGruppeServer sind die Server-Nachrichten der Gruppe mit Ende: eine
// Bestätigung, die Fehlerantwort des Schweregrads sev, connection_closed.
func endeGruppeServer(sev string) []model.Response {
	return append([]model.Response{{Type: model.ResponseBindComplete}}, endeAntworten(sev)...)
}

// extendedEnde ist eine Extended-Interaktion, deren letzte Gruppe mit Sync
// schließt und mit dem Ende des Servers endet; mit davor geht ihr eine Gruppe
// mit Flush voraus. Ihre Client-Nachrichten sind endeNachrichten.
func extendedEnde(seq int, sev string, davor bool) model.Interaction {
	in := model.Interaction{Sequence: seq, Request: model.Request{Type: model.RequestExtended}}
	if davor {
		in.Groups = append(in.Groups, model.Group{Client: endeNachrichten()[:2], Server: []model.Response{{Type: model.ResponseParseComplete}}})
	}
	return model.Interaction{Sequence: seq, Request: in.Request, Groups: append(in.Groups, model.Group{Client: endeNachrichten()[2:], Server: endeGruppeServer(sev)})}
}

// endeNachrichten sind die Client-Nachrichten von extendedEnde mit davor: eine
// Gruppe mit Flush (parse, flush) und die Gruppe mit Sync (bind, execute, sync).
func endeNachrichten() []model.ClientMessage {
	return []model.ClientMessage{parse("s1", "SELECT $1::text"), flushNachricht(), bind("s1", text("a")), execute(), syncNachricht()}
}

// abschluss ist die Art, wie die Gruppe mit Ende schließt: ihre Interaktion, die
// Client-Nachrichten der Interaktion mit davor und das Senden der ersten n davon.
type abschluss struct {
	ende        func(seq int, sev string, davor bool) model.Interaction
	nachrichten func() []model.ClientMessage
	bis         func(*testing.T, *services.ReplayService, model.SessionID, int) []model.Response
}

func mitSync() abschluss { return abschluss{extendedEnde, endeNachrichten, sendeBis} }

// endeVerbindung lädt die Interaktionen als eine Session und öffnet die erste
// Verbindung.
func endeVerbindung(t *testing.T, ins []model.Interaction, opts ...services.ReplayOption) (*services.ReplayService, model.SessionID) {
	t.Helper()
	s := ladeEnde(t, endeAufzeichnung(ins), opts...)
	id, _ := s.OpenConnection(context.Background())
	return s, id
}

// sendeBis sendet die ersten n der endeNachrichten, meldet die Gruppe mit Flush
// als gesendet und liefert die Antworten der letzten Nachricht; die Gruppe mit
// Ende meldet der Aufrufer.
func sendeBis(t *testing.T, s *services.ReplayService, id model.SessionID, n int) (out []model.Response) {
	t.Helper()
	for _, m := range endeNachrichten()[:n] {
		var err error
		if out, err = s.ClientMessage(context.Background(), id, m); err != nil {
			t.Fatalf("Nachricht %s: %v", m.Type, err)
		}
		if m.Type == model.ClientFlush {
			s.Sent(context.Background(), id)
		}
	}
	return out
}

// zahl zählt die Antworten des Typs typ.
func zahl(rs []model.Response, typ model.ResponseType) (n int) {
	for _, r := range rs {
		if r.Type == typ {
			n++
		}
	}
	return n
}

// Abdeckung: LH-FA-25/Happy, LH-FA-18/Boundary — Replay: Die Gruppe mit Ende, die
// mit Sync schließt, liefert ihre Server-Nachrichten in der aufgezeichneten
// Reihenfolge, die letzte ist connection_closed, ohne ready_for_query; geliefert
// wird erst mit dem abschließenden Sync, auch als einzige Gruppe (PANIC); die
// Gruppe mit Flush davor liefert ihre Nachrichten ohne connection_closed
// (LH-FA-25.a *Meldung des Endes im Replay*, *Stelle von connection_closed im
// Replay*, LH-FA-18.a).
func TestReplayEndeExtendedGruppeMitEnde(t *testing.T) {
	for _, sev := range []string{"FATAL", "PANIC"} {
		s, id := endeVerbindung(t, []model.Interaction{extendedEnde(1, sev, sev == "FATAL")})
		if sev == "FATAL" {
			if out := sendeBis(t, s, id, 2); !reflect.DeepEqual(out, []model.Response{{Type: model.ResponseParseComplete}}) {
				t.Fatalf("Gruppe mit Flush davor: %#v", out)
			}
		}
		for _, m := range endeNachrichten()[2:4] {
			if out, err := s.ClientMessage(context.Background(), id, m); err != nil || out != nil {
				t.Fatalf("%s: Antwort vor dem abschließenden Sync: %#v, %v", sev, out, err)
			}
		}
		out, err := s.ClientMessage(context.Background(), id, endeNachrichten()[4])
		if want := endeGruppeServer(sev); err != nil || !reflect.DeepEqual(out, want) {
			t.Fatalf("%s: erwartet %#v, erhalten %#v, %v", sev, want, out, err)
		}
		if out[len(out)-1].Type != model.ResponseConnectionClosed || zahl(out, model.ResponseReadyForQuery) != 0 {
			t.Fatalf("%s: connection_closed nicht letzte oder ready_for_query dabei: %#v", sev, out)
		}
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-18/Boundary — Replay: Eine Interaktion mit
// Sync vor der Interaktion mit Ende liefert ihre Gruppe ohne connection_closed
// und mit ready_for_query; die Gruppe mit Ende ist danach verbraucht. Eine Gruppe
// mit Sync vor der Gruppe mit Ende gibt es nur als frühere Interaktion, weil das
// Format Sync nur am Ende der letzten Gruppe zulässt (Interaction.Validate).
func TestReplayEndeExtendedSyncDavor(t *testing.T) {
	ins := []model.Interaction{ausfuehrung(1, "s1", "a", "A"), extendedEnde(2, "FATAL", false)}
	s, id := endeVerbindung(t, ins)
	out := sendeAlle(t, s, id, ins[0].Groups[0].Client...)
	if !reflect.DeepEqual(out, ins[0].Groups[0].Server) || zahl(out, model.ResponseConnectionClosed) != 0 {
		t.Fatalf("Interaktion mit Sync davor: %#v", out)
	}
	s.Sent(context.Background(), id)
	if out := sendeAlle(t, s, id, ins[1].Groups[0].Client...); !reflect.DeepEqual(out, endeGruppeServer("FATAL")) {
		t.Fatalf("Gruppe mit Ende: %#v", out)
	}
	if w := gesendetSchliessen(t, s, id); w != nil {
		t.Fatalf("verbraucht: %v", w)
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-03/Boundary — Replay: Die Interaktion mit
// Ende ist verbraucht, sobald die Gruppe gesendet ist (Sent); davor, also wenn
// das Senden scheitert, ist sie es nicht: PGR-W2001, mit FailOnUnconsumed
// PGR-E5002, beim Zwangsende PGR-E4006 mit ihrer Nummer; nach Sent weder Meldung
// noch Zwangsfehler.
func TestReplayEndeExtendedVerbrauch(t *testing.T) { endeVerbrauch(t, mitSync()) }

func endeVerbrauch(t *testing.T, a abschluss) {
	ctx := context.Background()
	in := []model.Interaction{a.ende(2, "FATAL", true)}
	s, id := endeVerbindung(t, in)
	a.bis(t, s, id, 5)
	s.Sent(ctx, id)
	if err := s.Forced(ctx, id); err != nil {
		t.Fatalf("Zwangsende nach Sent: %v", err)
	}
	if w := replaySchliessen(t, s, id); w != nil {
		t.Fatalf("gesendetes Ende gilt als unverbraucht: %v", w)
	}
	for _, opts := range [][]services.ReplayOption{nil, {services.FailOnUnconsumed}} {
		s, id = endeVerbindung(t, in, opts...)
		a.bis(t, s, id, 5)
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
// Stelle: die erste Nachricht der Interaktion, die erste und die mittlere der
// Gruppe mit Ende, Flush statt Sync als abschließende, Sync statt Flush in der
// Gruppe davor, eine einfache Anfrage vor und mitten in den Gruppen und eine
// Lebendprüfung mitten in einer Gruppe (auch nach der Gruppe davor). Der Cursor bleibt stehen: die richtigen Nachrichten danach
// liefern das Ende noch (LH-FA-25.a *Reihenfolge der Gründe im Replay*,
// LH-FA-18.a).
func TestReplayEndeExtendedAbweichung(t *testing.T) {
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
		{"Flush statt Sync als abschließende", 4, nachricht(flushNachricht())},
		{"einfache Anfrage vor den Gruppen", 0, anfrage("SELECT 1")},
		{"einfache Anfrage mitten in der Gruppe", 3, anfrage("SELECT 1")},
		{"Lebendprüfung mitten in der Gruppe davor", 1, anfrage("-- ping")},
		{"Lebendprüfung nach der Gruppe davor", 2, anfrage("-- ping")},
		{"Lebendprüfung mitten in der Gruppe mit Ende", 3, anfrage("-- ping")},
	} {
		s, id := endeVerbindung(t, []model.Interaction{extendedEnde(1, "FATAL", true)})
		sendeBis(t, s, id, f.vorher)
		if out, err := f.falsch(s, id); code(err) != model.CodeReplayMismatch || out != nil {
			t.Fatalf("%s: erwartet %s ohne Antworten, erhalten %#v, %v", f.name, model.CodeReplayMismatch, out, err)
		}
		var out []model.Response
		for _, m := range endeNachrichten()[f.vorher:] {
			out, _ = s.ClientMessage(ctx, id, m)
		}
		if !reflect.DeepEqual(out, endeGruppeServer("FATAL")) {
			t.Fatalf("%s: das Ende kommt nach der Abweichung noch: %#v", f.name, out)
		}
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-09/Boundary — Replay: Steht der Cursor vor
// der Interaktion mit Ende (Extended), ist eine Lebendprüfung, auch die leere,
// eine Lebendprüfung außer der Reihe: keine Abweichung, der Cursor bleibt; die
// Interaktion mit Ende ist keine einfache Anfrage mit Ende (LH-FA-09.a).
func TestReplayEndeExtendedLebendpruefungVorDerGruppe(t *testing.T) {
	s, id := endeVerbindung(t, []model.Interaction{extendedEnde(1, "FATAL", true)})
	pruefeLebend(t, s, id, "", "I")
	pruefeLebend(t, s, id, "-- ping", "I")
	if out := sendeBis(t, s, id, 5); !reflect.DeepEqual(out, endeGruppeServer("FATAL")) {
		t.Fatalf("Ende nach Lebendprüfungen: %#v", out)
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-03/Boundary — Replay: `Terminate` und das
// Ende der Client-Verbindung (für den Service dasselbe: CloseConnection) vor der
// Gruppe mit Ende, mitten in der Gruppe davor, nach ihr und mitten in der Gruppe
// mit Ende lassen die Interaktion unverbraucht: PGR-W2001 mit der Nummer der
// Interaktion, mit FailOnUnconsumed PGR-E5002 (LH-FA-03.b *Zeitpunkte* 1).
func TestReplayEndeExtendedTerminateUnverbraucht(t *testing.T) {
	endeTerminate(t, mitSync())
}

func endeTerminate(t *testing.T, a abschluss) {
	ins := []model.Interaction{ausfuehrung(1, "s1", "a", "A"), a.ende(7, "FATAL", true)}
	for n := 0; n <= 4; n++ {
		for _, opts := range [][]services.ReplayOption{nil, {services.FailOnUnconsumed}} {
			s, id := endeVerbindung(t, ins, opts...)
			sendeAlle(t, s, id, ins[0].Groups[0].Client...)
			s.Sent(context.Background(), id)
			a.bis(t, s, id, n)
			w, err := s.CloseConnection(context.Background(), id)
			if opts != nil && (code(err) != model.CodeReplayUnconsumed || !strings.Contains(err.Error(), "Nummer 7")) ||
				opts == nil && (err != nil || w == nil || w.Code != model.CodeUnconsumed || !strings.Contains(w.Msg, "1 von 2") || !strings.Contains(w.Msg, "Nummer 7")) {
				t.Fatalf("nach %d Nachrichten (FailOnUnconsumed %v): erwartet 1 von 2, Nummer 7: %v, %v", n, opts != nil, w, err)
			}
		}
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-13/Boundary — Replay: Das Herunterfahren
// gilt wie sonst (LH-FA-13.a): Mit dem Cursor in der Interaktion mit Ende, in der
// Gruppe davor, nach ihr und in der Gruppe mit Ende, erlaubt es Shutdown nicht,
// davor und nach dem Senden erlaubt es; das Zwangsende ist von der ersten
// Nachricht bis zum Senden der Gruppe mit Ende PGR-E4006 mit der Nummer der
// Interaktion.
func TestReplayEndeExtendedHerunterfahren(t *testing.T) { endeHerunterfahren(t, mitSync()) }

func endeHerunterfahren(t *testing.T, a abschluss) {
	ctx := context.Background()
	for i, w := range []struct{ shutdown, zwang bool }{{true, false}, {false, true}, {false, true}, {false, true}, {false, true}, {true, true}, {true, false}} {
		s, id := endeVerbindung(t, []model.Interaction{a.ende(7, "FATAL", true)})
		a.bis(t, s, id, min(i, 5))
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
// Die Interaktion mit Ende ist die erste Nachricht der nächsten freien Session;
// nach dem Ende der ersten Verbindung erhält die nächste die nächste Session mit
// ihrem Ende; eine Verbindung ohne freie Session erhält PGR-E5003 für die
// richtige wie die abweichende Nachricht, nie das Ende (LH-FA-12.a, LH-FA-25.a
// *Reihenfolge der Gründe im Replay*).
func TestReplayEndeExtendedFirstRequest(t *testing.T) { endeFirstRequest(t, mitSync()) }

func endeFirstRequest(t *testing.T, ab abschluss) {
	ctx := context.Background()
	s := ladeEnde(t, endeAufzeichnung([]model.Interaction{ab.ende(1, "FATAL", false)}, []model.Interaction{ab.ende(1, "PANIC", false)}))
	a, _ := s.OpenConnection(ctx)
	if out := sendeAlle(t, s, a, ab.nachrichten()[2:]...); !reflect.DeepEqual(out, endeGruppeServer("FATAL")) {
		t.Fatalf("erste Session: %#v", out)
	}
	if w := unassigned(t, s); w == nil || !strings.Contains(w.Msg, "1 aufgezeichnete") || !strings.Contains(w.Msg, "Kennung 2") {
		t.Fatalf("erste Session ist zugeordnet, die zweite nicht: %v", w)
	}
	if w := gesendetSchliessen(t, s, a); w != nil {
		t.Fatalf("Session 1: %v", w)
	}
	b, _ := s.OpenConnection(ctx)
	if out := sendeAlle(t, s, b, ab.nachrichten()[2:]...); !reflect.DeepEqual(out, endeGruppeServer("PANIC")) {
		t.Fatalf("zweite Session: %#v", out)
	}
	c, _ := s.OpenConnection(ctx)
	for _, m := range []model.ClientMessage{ab.nachrichten()[2], ab.nachrichten()[4]} {
		if out, err := s.ClientMessage(ctx, c, m); code(err) != model.CodeReplaySession || out != nil {
			t.Fatalf("%s ohne freie Session: erwartet %s, erhalten %#v, %v", m.Type, model.CodeReplaySession, out, err)
		}
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-12/Boundary — Replay mit `connection`: Die
// n-te Verbindung erhält die Session n mit ihrem Ende; das Schließen nach dem
// Ende ändert die Zuordnung der nächsten nicht; eine Verbindung ohne Session
// erhält PGR-E5003, auch für die Nachricht, die der Interaktion mit Ende
// entspräche (LH-FA-12.a).
func TestReplayEndeExtendedConnection(t *testing.T) { endeConnection(t, mitSync()) }

func endeConnection(t *testing.T, ab abschluss) {
	ctx := context.Background()
	rec := gekennzeichnet(sessionNr(1, "17.0", ab.ende(1, "FATAL", false)), sessionNr(2, "17.0", ausfuehrung(1, "s1", "a", "A")), sessionNr(3, "17.0", ab.ende(1, "PANIC", false)))
	rec.Version = model.VersionMitEnde
	s := nachVerbindung(t, rec)
	a, _ := s.OpenConnection(ctx)
	b, _ := s.OpenConnection(ctx)
	if out := sendeAlle(t, s, a, ab.nachrichten()[2:]...); !reflect.DeepEqual(out, endeGruppeServer("FATAL")) {
		t.Fatalf("Session 1 an Verbindung 1: %#v", out)
	}
	if w := gesendetSchliessen(t, s, a); w != nil {
		t.Fatalf("Session 1: %v", w)
	}
	if out := sendeAlle(t, s, b, endeNachrichten()[2:]...); zahl(out, model.ResponseConnectionClosed) != 0 || len(out) != 3 {
		t.Fatalf("Session 2 an Verbindung 2: %#v", out)
	}
	c, _ := s.OpenConnection(ctx)
	if out := sendeAlle(t, s, c, ab.nachrichten()[2:]...); !reflect.DeepEqual(out, endeGruppeServer("PANIC")) {
		t.Fatalf("Session 3 an Verbindung 3 nach dem Schließen der ersten: %#v", out)
	}
	d, _ := s.OpenConnection(ctx)
	if out, err := s.ClientMessage(ctx, d, ab.nachrichten()[2]); code(err) != model.CodeReplaySession || out != nil {
		t.Fatalf("Verbindung ohne Session: erwartet %s, erhalten %#v, %v", model.CodeReplaySession, out, err)
	}
}
