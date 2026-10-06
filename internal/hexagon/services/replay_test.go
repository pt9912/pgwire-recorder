package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

type ladeRepo struct {
	rec model.Recording
	err error
}

func (l ladeRepo) Prepare(context.Context, string, bool) error           { return nil }
func (l ladeRepo) Write(context.Context, string, model.Recording) error  { return nil }
func (l ladeRepo) Load(context.Context, string) (model.Recording, error) { return l.rec, l.err }

func interaktion(seq int, sql, tag string) model.Interaction {
	return model.Interaction{Sequence: seq, Request: model.Request{Type: model.RequestQuery, SQL: sql}, Responses: []model.Response{
		{Type: model.ResponseCommandComplete, Tag: tag},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}}
}

func aufzeichnung(sessions ...[]model.Interaction) model.Recording {
	rec := model.NewRecording()
	for i, in := range sessions {
		rec.Sessions = append(rec.Sessions, model.Session{ID: i + 1, ServerParameters: map[string]string{"server_version": "17.0", "client_encoding": "UTF8"}, Interactions: in})
	}
	return rec
}

// replaySchliessen beendet die Verbindung ohne Fehler und liefert die Warnung.
func replaySchliessen(t *testing.T, s *ReplayService, id model.SessionID) *model.Warning {
	t.Helper()
	w, err := s.CloseConnection(context.Background(), id)
	if err != nil {
		t.Fatalf("Fehler statt Warnung: %v", err)
	}
	return w
}

// gesendetSchliessen meldet die zuletzt gelieferten Antworten als gesendet, wie
// es der Adapter nach jedem Senden tut, und beendet dann die Verbindung.
func gesendetSchliessen(t *testing.T, s *ReplayService, id model.SessionID) *model.Warning {
	t.Helper()
	s.Sent(context.Background(), id)
	return replaySchliessen(t, s, id)
}

// unassigned liefert die Warnung über nie zugeordnete Sessions; ein Fehler
// bricht ab.
func unassigned(t *testing.T, s *ReplayService) *model.Warning {
	t.Helper()
	w, err := s.Unassigned()
	if err != nil {
		t.Fatalf("Fehler statt Warnung: %v", err)
	}
	return w
}

func code(err error) string {
	var me *model.Error
	if errors.As(err, &me) {
		return me.Code
	}
	return ""
}

// Abdeckung: LH-FA-09/Boundary — eine Anfrage, die der aufgezeichneten am Cursor
// gleicht, erhält deren Antworten; der Cursor rückt vor, und auch identische
// Anfragen werden der Reihe nach je einmal verbraucht.
func TestReplayStrictSequential(t *testing.T) {
	ctx := context.Background()
	s, err := NewReplayService(ctx, ladeRepo{rec: aufzeichnung([]model.Interaction{
		interaktion(1, "SELECT 1", "A"), interaktion(2, "SELECT 1", "B"),
	})}, "rec.yaml")
	if err != nil {
		t.Fatal(err)
	}
	id, handshake := s.OpenConnection(ctx)
	if len(handshake) != 3 || handshake[0].Name != "client_encoding" || handshake[2].Type != model.ResponseReadyForQuery {
		t.Fatalf("Handshake: %#v", handshake)
	}
	for _, want := range []string{"A", "B"} {
		out, err := s.Query(ctx, id, "SELECT 1")
		if err != nil || out[0].Tag != want {
			t.Fatalf("erwartet %s, erhalten %#v, %v", want, out, err)
		}
	}
	if w := gesendetSchliessen(t, s, id); w != nil {
		t.Fatalf("Warnung trotz verbrauchter Session: %v", w)
	}
}

// Abdeckung: LH-FA-09/Negative, LH-FA-10/Boundary, LH-FA-10/Negative — eine abweichende Anfrage,
// auch nur im Whitespace, ist PGR-E5001, erhält keine Antworten, auch nicht die der
// aufgezeichneten Anfrage, und verbraucht die Interaktion nicht; eine Anfrage über
// die aufgezeichneten hinaus ebenso.
func TestReplayMismatch(t *testing.T) {
	ctx := context.Background()
	s, _ := NewReplayService(ctx, ladeRepo{rec: aufzeichnung([]model.Interaction{interaktion(1, "SELECT 1", "A")})}, "rec.yaml")
	id, _ := s.OpenConnection(ctx)
	out, err := s.Query(ctx, id, "SELECT  1")
	if code(err) != model.CodeReplayMismatch || out != nil {
		t.Fatalf("erwartet %s ohne Antworten, erhalten %#v, %v", model.CodeReplayMismatch, out, err)
	}
	// Diagnose nach LH-FA-10.a: Session, erwartete Nummer, erwartete und
	// empfangene Anfrage.
	for _, teil := range []string{"Session 1", "Interaktion 1", `"SELECT 1"`, `"SELECT  1"`} {
		if !strings.Contains(err.Error(), teil) {
			t.Fatalf("Diagnose ohne %s: %v", teil, err)
		}
	}
	if out, err := s.Query(ctx, id, "SELECT 1"); err != nil || out[0].Tag != "A" {
		t.Fatalf("Cursor nach Mismatch verschoben: %#v, %v", out, err)
	}
	if out, err := s.Query(ctx, id, "SELECT 1"); code(err) != model.CodeReplayMismatch || out != nil {
		t.Fatalf("Anfrage über die Aufzeichnung hinaus: %#v, %v", out, err)
	}
}

// Sessions werden Verbindungen in der Reihenfolge
// ihrer ersten Anfrage zugeordnet; eine Verbindung ohne Anfrage zählt nicht, und
// eine Verbindung ohne freie Session erhält PGR-E5003; nicht verbrauchte
// Interaktionen und nie zugeordnete Sessions sind die Warnung PGR-W2001.
func TestReplaySessionZuordnung(t *testing.T) {
	ctx := context.Background()
	s, _ := NewReplayService(ctx, ladeRepo{rec: aufzeichnung(
		[]model.Interaction{interaktion(1, "S1", "eins"), interaktion(2, "S1b", "eins-b")},
		[]model.Interaction{interaktion(1, "S2", "zwei")},
		[]model.Interaction{interaktion(1, "S3", "drei")},
	)}, "rec.yaml")
	leer, _ := s.OpenConnection(ctx)
	a, _ := s.OpenConnection(ctx)
	b, _ := s.OpenConnection(ctx)
	if out, err := s.Query(ctx, b, "S1"); err != nil || out[0].Tag != "eins" {
		t.Fatalf("erste Anfrage erhält Session 1: %#v %v", out, err)
	}
	if out, err := s.Query(ctx, a, "S2"); err != nil || out[0].Tag != "zwei" {
		t.Fatalf("zweite Anfrage erhält Session 2: %#v %v", out, err)
	}
	if w := gesendetSchliessen(t, s, b); w == nil || w.Code != model.CodeUnconsumed {
		t.Fatalf("unverbrauchte Interaktion ohne Warnung: %v", w)
	}
	if w := replaySchliessen(t, s, leer); w != nil {
		t.Fatalf("Verbindung ohne Anfrage gewarnt: %v", w)
	}
	if w := unassigned(t, s); w == nil || w.Code != model.CodeUnconsumed {
		t.Fatalf("nie zugeordnete Session ohne Warnung: %v", w)
	}
	c, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, c, "S3"); err != nil {
		t.Fatal(err)
	}
	d, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, d, "S4"); code(err) != model.CodeReplaySession {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeReplaySession, err)
	}
}

// Abdeckung: LH-FA-03/Negative — eine Aufzeichnung ohne Session mit Interaktion
// ist PGR-E3004; ein Ladefehler geht unverändert zurück.
func TestReplayStartfehler(t *testing.T) {
	ctx := context.Background()
	if _, err := NewReplayService(ctx, ladeRepo{rec: model.NewRecording()}, "rec.yaml"); code(err) != model.CodeRecordingNoSession {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingNoSession, err)
	}
	ladefehler := model.Errorf(model.CodeRecordingBroken, nil, "kaputt")
	if _, err := NewReplayService(ctx, ladeRepo{err: ladefehler}, "rec.yaml"); code(err) != model.CodeRecordingBroken {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingBroken, err)
	}
}

// Der Handshake trägt die Serverparameter der nächsten noch nicht zugeordneten
// Session, gibt es keine mehr, der letzten; die Parameter stehen sortiert nach
// Namen (LH-FA-12.a, SPEC-004).
func TestReplayHandshake(t *testing.T) {
	ctx := context.Background()
	rec := model.NewRecording()
	for i, v := range []string{"16", "17"} {
		rec.Sessions = append(rec.Sessions, model.Session{ID: i + 1,
			ServerParameters: map[string]string{"server_version": v, "a": "1", "z": "2", "m": "3", "c": "4"},
			Interactions:     []model.Interaction{interaktion(1, "S", "x")}})
	}
	s, _ := NewReplayService(ctx, ladeRepo{rec: rec}, "rec.yaml")
	version := func(out []model.Response) string {
		for _, r := range out {
			if r.Name == "server_version" {
				return r.Value
			}
		}
		return ""
	}
	a, hs := s.OpenConnection(ctx)
	var namen []string
	for _, r := range hs[:len(hs)-1] {
		namen = append(namen, r.Name)
	}
	if strings.Join(namen, ",") != "a,c,m,server_version,z" {
		t.Fatalf("Reihenfolge: %v", namen)
	}
	if version(hs) != "16" {
		t.Fatalf("erste Verbindung: %v", hs)
	}
	if _, err := s.Query(ctx, a, "S"); err != nil {
		t.Fatal(err)
	}
	b, hs := s.OpenConnection(ctx)
	if version(hs) != "17" {
		t.Fatalf("zweite Verbindung: %v", hs)
	}
	if _, err := s.Query(ctx, b, "S"); err != nil {
		t.Fatal(err)
	}
	if _, hs := s.OpenConnection(ctx); version(hs) != "17" {
		t.Fatalf("ohne freie Session: %v", hs)
	}
}

// Steht am Cursor eine Extended-Interaktion, ist jede einfache Anfrage außer
// einer Lebendprüfung eine Abweichung (PGR-E5001), auch `;`; die leere Anfrage,
// deren SQL-Text dem leeren der Extended-Interaktion gleicht, ist eine
// Lebendprüfung und erhält EmptyQueryResponse (LH-FA-09.a). Der Cursor bleibt
// in beiden Fällen stehen.
func TestReplayExtendedAmCursor(t *testing.T) {
	ctx := context.Background()
	extended := model.Interaction{Sequence: 1, Request: model.Request{Type: model.RequestExtended}, Groups: []model.Group{{
		Client: []model.ClientMessage{{Type: model.ClientSync}},
		Server: []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}},
	}}}
	s, err := NewReplayService(ctx, ladeRepo{rec: aufzeichnung([]model.Interaction{extended})}, "rec.yaml")
	if err != nil {
		t.Fatal(err)
	}
	id, _ := s.OpenConnection(ctx)
	pruefeLebend(t, s, id, "", "I")
	for _, sql := range []string{";", "SELECT 1"} {
		if out, err := s.Query(ctx, id, sql); code(err) != model.CodeReplayMismatch {
			t.Fatalf("Anfrage %q: erwartet %s, erhalten %#v, %v", sql, model.CodeReplayMismatch, out, err)
		}
	}
	if w := gesendetSchliessen(t, s, id); w == nil {
		t.Fatal("Cursor ist vorgerückt: keine Warnung über die nicht verbrauchte Interaktion")
	}
}
