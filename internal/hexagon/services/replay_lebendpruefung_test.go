package services

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// mitStatus ist eine einfache Interaktion, deren ReadyForQuery den
// Transaktionsstatus tx trägt.
func mitStatus(seq int, sql, tx string) model.Interaction {
	in := interaktion(seq, sql, sql)
	in.Responses[len(in.Responses)-1].TxStatus = tx
	return in
}

// lebendAntwort ist die Antwort des Replay auf eine Lebendprüfung im Status tx.
func lebendAntwort(tx string) []model.Response {
	return []model.Response{{Type: model.ResponseEmptyQueryResponse}, {Type: model.ResponseReadyForQuery, TxStatus: tx}}
}

// pruefeLebend sendet sql als Lebendprüfung und erwartet die Antwort im Status
// tx.
func pruefeLebend(t *testing.T, s *ReplayService, id model.SessionID, sql, tx string) {
	t.Helper()
	out, err := s.Query(context.Background(), id, sql)
	if err != nil || !reflect.DeepEqual(out, lebendAntwort(tx)) {
		t.Fatalf("Lebendprüfung %q: erwartet %#v, erhalten %#v, %v", sql, lebendAntwort(tx), out, err)
	}
}

// Abdeckung: LH-FA-09/Negative — eine Lebendprüfung zwischen zwei
// Interaktionen, auch vor der ersten und nach der letzten, beantwortet das
// Replay mit EmptyQueryResponse und ReadyForQuery im Transaktionsstatus des
// letzten ReadyForQuery der Verbindung (I nach dem Handshake, T in der
// Transaktion, E nach einem Fehler darin); der Cursor bleibt stehen, und sie
// ordnet keine Session zu, auch nicht, wenn keine mehr frei ist.
func TestReplayLebendpruefungAusserDerReihe(t *testing.T) {
	ctx := context.Background()
	s, err := NewReplayService(ctx, ladeRepo{rec: aufzeichnung([]model.Interaction{
		mitStatus(1, "BEGIN", "T"), mitStatus(2, "SELECT x", "E"), mitStatus(3, "ROLLBACK", "I"), mitStatus(4, "BEGIN", "T"),
	})}, "rec.yaml")
	if err != nil {
		t.Fatal(err)
	}
	a, _ := s.OpenConnection(ctx)
	pruefeLebend(t, s, a, "-- ping", "I")
	b, _ := s.OpenConnection(ctx)
	pruefeLebend(t, s, b, "", "I")
	for _, sch := range []struct{ sql, tx string }{{"BEGIN", "T"}, {"SELECT x", "E"}, {"ROLLBACK", "I"}, {"BEGIN", "T"}} {
		if _, err := s.Query(ctx, b, sch.sql); err != nil {
			t.Fatalf("%s nach Lebendprüfung: %v", sch.sql, err)
		}
		pruefeLebend(t, s, b, "-- ping", sch.tx)
		pruefeLebend(t, s, b, "/* a */ -- b", sch.tx)
	}
	if w := s.CloseConnection(ctx, b); w != nil {
		t.Fatalf("Lebendprüfungen haben den Cursor bewegt: %v", w)
	}
	if _, err := s.Query(ctx, a, "BEGIN"); code(err) != model.CodeReplaySession {
		t.Fatalf("die erste Verbindung hätte keine Session mehr frei haben dürfen: %v", err)
	}
	c, _ := s.OpenConnection(ctx)
	pruefeLebend(t, s, c, "-- ping", "I")
}

// Abdeckung: LH-FA-09/Negative — das Replay liest eine
// Aufzeichnung, als stünden die Lebendprüfungen nicht darin: Der Cursor
// überspringt sie, sie sind nicht unverbraucht, eine Session nur aus
// Lebendprüfungen wird bei der Zuordnung übersprungen und ist nie unzugeordnet,
// eine Aufzeichnung nur aus Lebendprüfungen ist PGR-E3004; Diagnosen nennen die
// aufgezeichneten Nummern. Andere Anfragen, auch `;` und ein offener
// Blockkommentar, bleiben in der Aufzeichnung.
func TestReplayLebendpruefungAufgezeichnet(t *testing.T) {
	ctx := context.Background()
	rec := aufzeichnung(
		[]model.Interaction{interaktion(1, "-- ping", ""), interaktion(2, "SELECT 1", "A"), interaktion(3, "/* x */", ""), interaktion(4, "SELECT 2", "B"), interaktion(5, "-- ping", "")},
		[]model.Interaction{interaktion(1, "-- ping", ""), interaktion(2, "\n", "")},
		[]model.Interaction{interaktion(1, ";", "C"), interaktion(2, "/* offen", "D")},
	)
	s, err := NewReplayService(ctx, ladeRepo{rec: rec}, "rec.yaml")
	if err != nil {
		t.Fatal(err)
	}
	a, _ := s.OpenConnection(ctx)
	_, err = s.Query(ctx, a, "SELECT 2")
	if code(err) != model.CodeReplayMismatch || !strings.Contains(err.Error(), "Interaktion 2") {
		t.Fatalf("Diagnose mit aufgezeichneter Nummer: %v", err)
	}
	for _, want := range []string{"A", "B"} {
		sql := map[string]string{"A": "SELECT 1", "B": "SELECT 2"}[want]
		if out, err := s.Query(ctx, a, sql); err != nil || out[0].Tag != want {
			t.Fatalf("%s: %#v, %v", sql, out, err)
		}
	}
	_, err = s.Query(ctx, a, "SELECT 3")
	if code(err) != model.CodeReplayMismatch || !strings.Contains(err.Error(), "nach 4") {
		t.Fatalf("nach der letzten Interaktion: %v", err)
	}
	if w := s.CloseConnection(ctx, a); w != nil {
		t.Fatalf("aufgezeichnete Lebendprüfung gilt als unverbraucht: %v", w)
	}
	b, _ := s.OpenConnection(ctx)
	if out, err := s.Query(ctx, b, ";"); err != nil || out[0].Tag != "C" {
		t.Fatalf("zweite Verbindung erhält Session 3 mit `;`: %#v, %v", out, err)
	}
	if out, err := s.Query(ctx, b, "/* offen"); err != nil || out[0].Tag != "D" {
		t.Fatalf("offener Blockkommentar bleibt Interaktion: %#v, %v", out, err)
	}
	if w := s.Unassigned(); w != nil {
		t.Fatalf("Session nur aus Lebendprüfungen gilt als unzugeordnet: %v", w)
	}

	s, _ = NewReplayService(ctx, ladeRepo{rec: aufzeichnung(
		[]model.Interaction{interaktion(1, "-- ping", ""), interaktion(2, "SELECT 1", "A"), interaktion(3, "-- ping", ""), interaktion(4, "SELECT 2", "B")},
	)}, "rec.yaml")
	c, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, c, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if w := s.CloseConnection(ctx, c); w == nil || !strings.Contains(w.Msg, "1 von 2") {
		t.Fatalf("Warnung zählt die Lebendprüfungen mit: %v", w)
	}

	if _, err := NewReplayService(ctx, ladeRepo{rec: aufzeichnung(
		[]model.Interaction{interaktion(1, "-- ping", "")}, []model.Interaction{interaktion(1, "", "")},
	)}, "rec.yaml"); code(err) != model.CodeRecordingNoSession {
		t.Fatalf("Aufzeichnung nur aus Lebendprüfungen: erwartet %s, erhalten %v", model.CodeRecordingNoSession, err)
	}
}

// Abdeckung: LH-FA-18/Negative, LH-FA-09/Negative — Extended: eine
// Lebendprüfung vor einer Extended-Interaktion und nach ihrem Sync beantwortet
// das Replay mit dem Status des ReadyForQuery ihrer Sync-Gruppe; mitten in der
// Interaktion, nach der ersten Nachricht wie nach einer Flush-Gruppe, ist sie
// eine Abweichung (PGR-E5001), ebenso jede andere einfache Anfrage am Cursor.
func TestReplayLebendpruefungExtended(t *testing.T) {
	ctx := context.Background()
	in := vorbereitung(1)
	in.Groups[1].Server[len(in.Groups[1].Server)-1] = model.Response{Type: model.ResponseReadyForQuery, TxStatus: "T"}
	s, id := replayMit(t, []model.Interaction{in})
	pruefeLebend(t, s, id, "-- ping", "I")
	if _, err := s.Query(ctx, id, "SELECT 1"); code(err) != model.CodeReplayMismatch {
		t.Fatalf("Anweisung am Cursor einer Extended-Interaktion: %v", err)
	}
	sendeAlle(t, s, id, in.Groups[0].Client[0])
	for _, sql := range []string{"-- ping", ""} {
		_, err := s.Query(ctx, id, sql)
		if code(err) != model.CodeReplayMismatch || !strings.Contains(err.Error(), "Gruppe 1, Nachricht 2") {
			t.Fatalf("Lebendprüfung %q nach der ersten Nachricht: %v", sql, err)
		}
	}
	sendeAlle(t, s, id, in.Groups[0].Client[1:]...)
	if _, err := s.Query(ctx, id, "-- ping"); code(err) != model.CodeReplayMismatch {
		t.Fatalf("Lebendprüfung nach der Flush-Gruppe: %v", err)
	}
	sendeAlle(t, s, id, in.Groups[1].Client...)
	pruefeLebend(t, s, id, "-- ping", "T")
	if w := s.CloseConnection(ctx, id); w != nil {
		t.Fatalf("Warnung trotz verbrauchter Session: %v", w)
	}
}
