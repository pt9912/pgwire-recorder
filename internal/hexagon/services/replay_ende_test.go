package services_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// endeAntworten sind die Antworten einer einfachen Anfrage mit Ende: eine
// Fehlerantwort des Schweregrads sev, danach connection_closed.
func endeAntworten(sev string) []model.Response {
	return []model.Response{
		{Type: model.ResponseErrorResponse, Fields: map[string]string{"S": sev, "V": sev, "C": "57P01", "M": "terminating connection"}},
		{Type: model.ResponseConnectionClosed},
	}
}

// mitEnde ist eine einfache Interaktion, die mit dem Ende des Servers schließt.
func mitEnde(seq int, sql, sev string) model.Interaction {
	return model.Interaction{Sequence: seq, Request: model.Request{Type: model.RequestQuery, SQL: sql}, Responses: endeAntworten(sev)}
}

// endeAufzeichnung ist eine Aufzeichnung der Version 2 mit den Sessions.
func endeAufzeichnung(sessions ...[]model.Interaction) model.Recording {
	rec := aufzeichnung(sessions...)
	rec.Version = model.VersionMitEnde
	return rec
}

func ladeEnde(t *testing.T, rec model.Recording, opts ...services.ReplayOption) *services.ReplayService {
	t.Helper()
	s, err := services.NewReplayService(context.Background(), ladeRepo{rec: rec}, "rec.yaml", opts...)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// Abdeckung: LH-FA-25/Happy, LH-FA-25/Boundary — Replay: Eine einfache Anfrage
// mit Ende erhält die aufgezeichneten Antworten unverändert, die letzte ist
// connection_closed, ohne ready_for_query (LH-FA-25.a *Meldung des Endes im
// Replay*); die Interaktionen davor laufen wie sonst. Schweregrad FATAL und
// PANIC liefern je dasselbe.
func TestReplayEndeAntwortenMitConnectionClosed(t *testing.T) {
	for _, sev := range []string{"FATAL", "PANIC"} {
		ctx := context.Background()
		s := ladeEnde(t, endeAufzeichnung([]model.Interaction{interaktion(1, "SELECT 1", "A"), mitEnde(2, "SELECT 2", sev)}))
		id, _ := s.OpenConnection(ctx)
		if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
			t.Fatal(err)
		}
		out, err := s.Query(ctx, id, "SELECT 2")
		if err != nil || !reflect.DeepEqual(out, endeAntworten(sev)) {
			t.Fatalf("%s: erwartet %#v, erhalten %#v, %v", sev, endeAntworten(sev), out, err)
		}
		for _, r := range out {
			if r.Type == model.ResponseReadyForQuery {
				t.Fatalf("%s: ready_for_query in den Antworten des Endes", sev)
			}
		}
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-03/Boundary — Replay: Die Interaktion mit
// Ende ist verbraucht, sobald ihre Antworten gesendet sind (Sent); davor, also
// wenn das Senden scheitert, ist sie es nicht (PGR-W2001, mit FailOnUnconsumed
// PGR-E5002). Nach dem Senden meldet das Schließen nichts.
func TestReplayEndeVerbrauchtMitSent(t *testing.T) {
	ctx := context.Background()
	rec := endeAufzeichnung([]model.Interaction{mitEnde(1, "SELECT 1", "FATAL")})

	s := ladeEnde(t, rec)
	id, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if w := gesendetSchliessen(t, s, id); w != nil {
		t.Fatalf("gesendetes Ende gilt als unverbraucht: %v", w)
	}

	s = ladeEnde(t, rec)
	id, _ = s.OpenConnection(ctx)
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if w := replaySchliessen(t, s, id); w == nil || w.Code != model.CodeUnconsumed || !strings.Contains(w.Msg, "1 von 1") {
		t.Fatalf("Schreibfehler: Interaktion mit Ende muss unverbraucht sein (PGR-W2001): %v", w)
	}

	s = ladeEnde(t, rec, services.FailOnUnconsumed)
	id, _ = s.OpenConnection(ctx)
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CloseConnection(ctx, id); code(err) != model.CodeReplayUnconsumed {
		t.Fatalf("Schreibfehler mit FailOnUnconsumed: erwartet %s, erhalten %v", model.CodeReplayUnconsumed, err)
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-03/Boundary — Replay: Hat der Client vor
// der Interaktion mit Ende `Terminate` gesendet (die Verbindung endet ohne
// Query), ist sie unverbraucht: PGR-W2001, mit FailOnUnconsumed PGR-E5002; die
// Zahl nennt sie mit.
func TestReplayEndeTerminateDavorLaesstUnverbraucht(t *testing.T) {
	ctx := context.Background()
	rec := endeAufzeichnung([]model.Interaction{interaktion(1, "SELECT 1", "A"), mitEnde(2, "SELECT 2", "FATAL")})
	s := ladeEnde(t, rec)
	id, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if w := gesendetSchliessen(t, s, id); w == nil || w.Code != model.CodeUnconsumed || !strings.Contains(w.Msg, "1 von 2") || !strings.Contains(w.Msg, "Nummer 2") {
		t.Fatalf("Ende nicht erreicht: erwartet PGR-W2001 mit 1 von 2, Nummer 2: %v", w)
	}

	s = ladeEnde(t, rec, services.FailOnUnconsumed)
	id, _ = s.OpenConnection(ctx)
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	s.Sent(ctx, id)
	if _, err := s.CloseConnection(ctx, id); code(err) != model.CodeReplayUnconsumed {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeReplayUnconsumed, err)
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-10/Negative — Replay: Die Gründe einer
// Ablehnung gehen dem Ende voraus. Eine Anfrage mit anderem SQL-Text ist
// PGR-E5001 und liefert kein connection_closed; die Interaktion bleibt
// unverbraucht und mit dem richtigen Text zu haben. Ohne freie Session ist es
// PGR-E5003 vor jedem Vergleich (LH-FA-25.a *Reihenfolge der Gründe im Replay*).
func TestReplayEndeAblehnungVorDemEnde(t *testing.T) {
	ctx := context.Background()
	s := ladeEnde(t, endeAufzeichnung([]model.Interaction{mitEnde(1, "SELECT 1", "FATAL")}))
	a, _ := s.OpenConnection(ctx)
	out, err := s.Query(ctx, a, "SELECT 9")
	if code(err) != model.CodeReplayMismatch || out != nil {
		t.Fatalf("anderer Text: erwartet %s ohne Antworten, erhalten %#v, %v", model.CodeReplayMismatch, out, err)
	}
	if out, err := s.Query(ctx, a, "SELECT 1"); err != nil || !reflect.DeepEqual(out, endeAntworten("FATAL")) {
		t.Fatalf("Interaktion mit Ende nach Abweichung: %#v, %v", out, err)
	}
	// Die einzige Session ist zugeordnet; die zweite Verbindung hat keine frei und
	// erhält PGR-E5003, nicht das Ende und nicht PGR-E5001.
	b, _ := s.OpenConnection(ctx)
	out, err = s.Query(ctx, b, "SELECT 9")
	if code(err) != model.CodeReplaySession || out != nil {
		t.Fatalf("keine freie Session: erwartet %s, erhalten %#v, %v", model.CodeReplaySession, out, err)
	}
	out, err = s.Query(ctx, b, "SELECT 1")
	if code(err) != model.CodeReplaySession || out != nil {
		t.Fatalf("keine freie Session mit passendem Text: erwartet %s ohne Ende, erhalten %#v, %v", model.CodeReplaySession, out, err)
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-12/Boundary — Replay: Jede Session einer
// Aufzeichnung mit mehreren Enden hat ihr eigenes; nach dem Ende der ersten
// erhält die nächste Verbindung die nächste Session (`first-request`) mit ihren
// Antworten und ihrem Ende.
func TestReplayEndeMehrereSessionsMitEnde(t *testing.T) {
	ctx := context.Background()
	s := ladeEnde(t, endeAufzeichnung(
		[]model.Interaction{interaktion(1, "SELECT 1", "A"), mitEnde(2, "SELECT 2", "FATAL")},
		[]model.Interaction{mitEnde(1, "SELECT 3", "PANIC")},
	))
	a, _ := s.OpenConnection(ctx)
	for _, sql := range []string{"SELECT 1", "SELECT 2"} {
		if _, err := s.Query(ctx, a, sql); err != nil {
			t.Fatal(err)
		}
	}
	if w := gesendetSchliessen(t, s, a); w != nil {
		t.Fatalf("Session 1: %v", w)
	}
	if w := unassigned(t, s); w == nil || !strings.Contains(w.Msg, "1 aufgezeichnete") || !strings.Contains(w.Msg, "Kennung 2") {
		t.Fatalf("die Session 2 mit Ende ist noch nicht zugeordnet: %v", w)
	}
	b, _ := s.OpenConnection(ctx)
	if out, err := s.Query(ctx, b, "SELECT 3"); err != nil || !reflect.DeepEqual(out, endeAntworten("PANIC")) {
		t.Fatalf("Session 2: %#v, %v", out, err)
	}
	if w := gesendetSchliessen(t, s, b); w != nil {
		t.Fatalf("Session 2: %v", w)
	}
	if w := unassigned(t, s); w != nil {
		t.Fatalf("beide Sessions zugeordnet: %v", w)
	}
}

// lebendMitEnde ist eine Lebendprüfung, die mit dem Ende des Servers schließt.
func lebendMitEnde(seq int, sql string) model.Interaction { return mitEnde(seq, sql, "FATAL") }

// Abdeckung: LH-FA-09/Boundary, LH-FA-25/Boundary — Replay: Eine aufgezeichnete
// Lebendprüfung mit Ende wird nicht übersprungen: Sie ist die Interaktion am
// Cursor, ihre Antworten werden verwendet (nicht EmptyQueryResponse), und sie
// ist ohne Anfrage unverbraucht. Eine ohne Ende bleibt übersprungen (LH-FA-09.a
// *Lebendprüfung mit Ende*).
func TestReplayEndeLebendpruefungMitEndeBleibtInteraktion(t *testing.T) {
	ctx := context.Background()
	rec := endeAufzeichnung([]model.Interaction{interaktion(1, "-- ping", ""), interaktion(2, "SELECT 1", "A"), interaktion(3, "SELECT 2", "B"), lebendMitEnde(4, "-- ping")})
	s := ladeEnde(t, rec)
	id, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	s.Sent(ctx, id)
	if w := replaySchliessen(t, s, id); w == nil || !strings.Contains(w.Msg, "2 von 3") || !strings.Contains(w.Msg, "Nummer 3") {
		t.Fatalf("Lebendprüfung mit Ende zählt als unverbraucht, die ohne Ende nicht: %v", w)
	}

	s = ladeEnde(t, rec)
	id, _ = s.OpenConnection(ctx)
	// Steht der Cursor noch vor der Lebendprüfung mit Ende, vor der Zuordnung wie
	// mitten in der Session, ist die eingehende mit demselben Text eine
	// Lebendprüfung außer der Reihe.
	pruefeLebend(t, s, id, "-- ping", "I")
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	s.Sent(ctx, id)
	pruefeLebend(t, s, id, "-- ping", "I")
	if _, err := s.Query(ctx, id, "SELECT 2"); err != nil {
		t.Fatal(err)
	}
	s.Sent(ctx, id)
	if out, err := s.Query(ctx, id, "-- ping"); err != nil || !reflect.DeepEqual(out, endeAntworten("FATAL")) {
		t.Fatalf("eingehende Lebendprüfung mit demselben Text ist die Anfrage: %#v, %v", out, err)
	}
	if w := gesendetSchliessen(t, s, id); w != nil {
		t.Fatalf("verbraucht: %v", w)
	}
}

// Abdeckung: LH-FA-09/Boundary, LH-FA-25/Negative — Replay: Steht der Cursor
// auf einer Lebendprüfung mit Ende, ist nur die eingehende mit gleichem
// SQL-Text diese Anfrage; eine mit anderem Text (hier ein anderer Kommentar und
// ein Leertext) wird außer der Reihe beantwortet, ohne Cursor und ohne Ende
// (LH-FA-09.a *Lebendprüfung mit Ende*).
func TestReplayEndeLebendpruefungAndererTextAusserDerReihe(t *testing.T) {
	ctx := context.Background()
	s := ladeEnde(t, endeAufzeichnung([]model.Interaction{lebendMitEnde(1, "-- ping")}))
	id, _ := s.OpenConnection(ctx)
	for _, sql := range []string{"-- anders", ""} {
		pruefeLebend(t, s, id, sql, "I")
	}
	if out, err := s.Query(ctx, id, "-- ping"); err != nil || !reflect.DeepEqual(out, endeAntworten("FATAL")) {
		t.Fatalf("Lebendprüfungen mit anderem Text haben den Cursor bewegt: %#v, %v", out, err)
	}
}

// Abdeckung: LH-FA-09/Boundary, LH-FA-12/Boundary — Replay mit `first-request`:
// Ist die erste Interaktion der nächsten freien Session eine Lebendprüfung mit
// Ende und gleichem SQL-Text, ist die eingehende Lebendprüfung vor der
// Zuordnung diese Anfrage und ordnet zu; mit anderem Text gilt Schritt 1: außer
// der Reihe, ohne Zuordnung (LH-FA-09.a *Lebendprüfung mit Ende*). Eine Session
// nur aus dieser Lebendprüfung ist eine Session mit Interaktion.
func TestReplayEndeFirstRequestLebendpruefungVorZuordnung(t *testing.T) {
	ctx := context.Background()
	rec := endeAufzeichnung([]model.Interaction{lebendMitEnde(1, "-- ping")}, []model.Interaction{interaktion(1, "SELECT 1", "A")})
	s := ladeEnde(t, rec)
	a, _ := s.OpenConnection(ctx)
	pruefeLebend(t, s, a, "-- anders", "I")
	if w := unassigned(t, s); w == nil || !strings.Contains(w.Msg, "2 aufgezeichnete") {
		t.Fatalf("Lebendprüfung mit anderem Text hat zugeordnet: %v", w)
	}
	out, err := s.Query(ctx, a, "-- ping")
	if err != nil || !reflect.DeepEqual(out, endeAntworten("FATAL")) {
		t.Fatalf("Lebendprüfung mit gleichem Text: %#v, %v", out, err)
	}
	if w := unassigned(t, s); w == nil || !strings.Contains(w.Msg, "1 aufgezeichnete") || !strings.Contains(w.Msg, "Kennung 2") {
		t.Fatalf("erste Session ist zugeordnet, die zweite nicht: %v", w)
	}
	if w := gesendetSchliessen(t, s, a); w != nil {
		t.Fatalf("verbraucht: %v", w)
	}

	// Die zweite Session beginnt nicht mit einer Anfrage mit Ende: Eine
	// Lebendprüfung vor der Zuordnung ist außer der Reihe, auch mit dem Text der
	// Lebendprüfung mit Ende einer späteren Session.
	s = ladeEnde(t, endeAufzeichnung([]model.Interaction{interaktion(1, "SELECT 1", "A")}, []model.Interaction{lebendMitEnde(1, "-- ping")}))
	b, _ := s.OpenConnection(ctx)
	pruefeLebend(t, s, b, "-- ping", "I")
	if w := unassigned(t, s); w == nil || !strings.Contains(w.Msg, "2 aufgezeichnete") {
		t.Fatalf("Lebendprüfung hat die erste Session zugeordnet: %v", w)
	}
}

// Abdeckung: LH-FA-12/Boundary, LH-FA-25/Boundary — Replay mit `connection`: Die
// Session ist mit der Annahme zugeordnet. Die Lebendprüfung mit Ende ist die
// Anfrage, wenn der Cursor auf ihr steht; das Schließen der Verbindung nach dem
// Ende ändert die Zuordnung nicht, und die Session mit Ende zählt wie jede
// andere für Unassigned (Bindung aus slice-v1-abschluss-sessions, LH-FA-12.a).
func TestReplayConnectionEndeBleibtZugeordnet(t *testing.T) {
	ctx := context.Background()
	rec := gekennzeichnet(
		sessionNr(1, "17.0", mitEnde(1, "A", "FATAL")),
		sessionNr(2, "17.0", interaktion(1, "B", "b"), lebendMitEnde(2, "-- ping")),
		sessionNr(3, "17.0", interaktion(1, "C", "c")),
		sessionNr(4, "17.0", mitEnde(1, "D", "PANIC")),
	)
	rec.Version = model.VersionMitEnde
	s := nachVerbindung(t, rec)
	a, _ := s.OpenConnection(ctx)
	b, _ := s.OpenConnection(ctx)
	if w := unassigned(t, s); w == nil || !strings.Contains(w.Msg, "2 aufgezeichnete") || !strings.Contains(w.Msg, "Kennung 3") {
		t.Fatalf("die Sessions 3 und 4, auch die mit Ende, sind nie zugeordnet: %v", w)
	}
	if _, err := s.Query(ctx, a, "A"); err != nil {
		t.Fatal(err)
	}
	if w := gesendetSchliessen(t, s, a); w != nil {
		t.Fatalf("Session 1: %v", w)
	}
	if _, err := s.Query(ctx, b, "B"); err != nil {
		t.Fatal(err)
	}
	s.Sent(ctx, b)
	if out, err := s.Query(ctx, b, "-- ping"); err != nil || !reflect.DeepEqual(out, endeAntworten("FATAL")) {
		t.Fatalf("Lebendprüfung mit Ende bei connection: %#v, %v", out, err)
	}
	if w := gesendetSchliessen(t, s, b); w != nil {
		t.Fatalf("Session 2: %v", w)
	}
	// Die Verbindung 1 hat geschlossen; die Zuordnung der dritten ändert das nicht.
	c, _ := s.OpenConnection(ctx)
	if out, err := s.Query(ctx, c, "C"); err != nil || out[0].Tag != "c" {
		t.Fatalf("Session 3 an Verbindung 3: %#v, %v", out, err)
	}
	if w := unassigned(t, s); w == nil || !strings.Contains(w.Msg, "1 aufgezeichnete") || !strings.Contains(w.Msg, "Kennung 4") {
		t.Fatalf("die Session 4 mit Ende ist nie zugeordnet: %v", w)
	}
}
