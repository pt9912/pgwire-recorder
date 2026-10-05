package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

func text(s string) model.Value { return model.Value{Bytes: []byte(s)} }

func parse(name, sql string) model.ClientMessage {
	return model.ClientMessage{Type: model.ClientParse, Statement: name, SQL: sql, ParamTypes: []uint32{25}}
}

func bind(name string, werte ...model.Value) model.ClientMessage {
	return model.ClientMessage{Type: model.ClientBind, Statement: name, ParamFormats: []int16{0}, Params: werte, ResultFormats: []int16{0}}
}

var (
	execute        = model.ClientMessage{Type: model.ClientExecute}
	flushNachricht = model.ClientMessage{Type: model.ClientFlush}
	syncNachricht  = model.ClientMessage{Type: model.ClientSync}
	rfq            = model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"}
)

func tag(t string) model.Response { return model.Response{Type: model.ResponseCommandComplete, Tag: t} }

// ausfuehrung ist eine Extended-Interaktion aus einer Sync-Gruppe, die das
// Statement name mit dem Wert ausführt und mit dem Befehlsabschluss t antwortet.
func ausfuehrung(seq int, name, wert, t string) model.Interaction {
	return model.Interaction{Sequence: seq, Request: model.Request{Type: model.RequestExtended}, Groups: []model.Group{{
		Client: []model.ClientMessage{bind(name, text(wert)), execute, syncNachricht},
		Server: []model.Response{{Type: model.ResponseBindComplete}, tag(t), rfq},
	}}}
}

// vorbereitung ist eine Extended-Interaktion aus einer Flush-Gruppe (Parse,
// Describe) und einer Sync-Gruppe (Bind, Execute).
func vorbereitung(seq int) model.Interaction {
	return model.Interaction{Sequence: seq, Request: model.Request{Type: model.RequestExtended}, Groups: []model.Group{
		{
			Client: []model.ClientMessage{parse("s1", "SELECT $1::text"), {Type: model.ClientDescribe, Target: model.TargetStatement, Name: "s1"}, flushNachricht},
			Server: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseParameterDescription, ParamTypes: []uint32{25}}, {Type: model.ResponseRowDescription}},
		},
		{
			Client: []model.ClientMessage{bind("s1", text("a")), execute, syncNachricht},
			Server: []model.Response{{Type: model.ResponseBindComplete}, {Type: model.ResponseDataRow, Values: []model.Value{text("a")}}, tag("SELECT 1"), rfq},
		},
	}}
}

func replayMit(t *testing.T, sessions ...[]model.Interaction) (*ReplayService, model.SessionID) {
	t.Helper()
	ctx := context.Background()
	s, err := NewReplayService(ctx, ladeRepo{rec: aufzeichnung(sessions...)}, "rec.yaml")
	if err != nil {
		t.Fatal(err)
	}
	id, _ := s.OpenConnection(ctx)
	return s, id
}

// sendeAlle übergibt die Nachrichten der Reihe nach und liefert die Antworten
// der letzten; jede frühere muss ohne Antwort bleiben.
func sendeAlle(t *testing.T, s *ReplayService, id model.SessionID, msgs ...model.ClientMessage) []model.Response {
	t.Helper()
	var out []model.Response
	for i, m := range msgs {
		var err error
		out, err = s.ClientMessage(context.Background(), id, m)
		if err != nil {
			t.Fatalf("Nachricht %d (%s): %v", i+1, m.Type, err)
		}
		if i < len(msgs)-1 && out != nil {
			t.Fatalf("Nachricht %d (%s) vor dem Ende der Gruppe beantwortet: %#v", i+1, m.Type, out)
		}
	}
	return out
}

// Abdeckung: LH-FA-18/Happy — Replay-Hälfte: jede Client-Nachricht einer Gruppe
// wird mit der erwarteten verglichen; die Server-Nachrichten einer Gruppe gehen
// erst nach ihrer letzten Client-Nachricht (Flush oder Sync) zurück, vorher
// keine; nach dem Sync der letzten Gruppe steht der Cursor auf der nächsten
// Interaktion, auch einer einfachen, und die Session ist verbraucht.
func TestReplayExtendedGruppen(t *testing.T) {
	s, id := replayMit(t, []model.Interaction{vorbereitung(1), interaktion(2, "SELECT 7", "SELECT 1")})
	in := vorbereitung(1)
	for gi, g := range in.Groups {
		out := sendeAlle(t, s, id, g.Client...)
		if !reflect.DeepEqual(out, g.Server) {
			t.Fatalf("Gruppe %d: %#v statt %#v", gi+1, out, g.Server)
		}
	}
	if out, err := s.Query(context.Background(), id, "SELECT 7"); err != nil || out[0].Tag != "SELECT 1" {
		t.Fatalf("einfache Anfrage nach der Extended-Interaktion: %#v, %v", out, err)
	}
	if w := s.CloseConnection(context.Background(), id); w != nil {
		t.Fatalf("Warnung trotz verbrauchter Session: %v", w)
	}
}

// Abdeckung: LH-FA-18/Boundary, LH-FA-09/Boundary — dasselbe Prepared Statement,
// mehrfach mit verschiedenen und mit gleichen Werten ausgeführt, erhält je
// Ausführung die Antwort seiner Position; eine Ausführung in anderer
// Reihenfolge ist eine Abweichung.
func TestReplayExtendedWiederholung(t *testing.T) {
	sitzung := []model.Interaction{ausfuehrung(1, "s1", "a", "A"), ausfuehrung(2, "s1", "b", "B"), ausfuehrung(3, "s1", "a", "C")}
	s, id := replayMit(t, sitzung)
	for i, w := range []struct{ wert, tag string }{{"a", "A"}, {"b", "B"}, {"a", "C"}} {
		out := sendeAlle(t, s, id, bind("s1", text(w.wert)), execute, syncNachricht)
		if len(out) != 3 || out[1].Tag != w.tag {
			t.Fatalf("Ausführung %d: %#v", i+1, out)
		}
	}

	s, id = replayMit(t, sitzung)
	if _, err := s.ClientMessage(context.Background(), id, bind("s1", text("b"))); code(err) != model.CodeReplayMismatch {
		t.Fatalf("Ausführung außer der Reihe: %v", err)
	}
}

// Abdeckung: LH-FA-10/Boundary — eine Client-Nachricht, die der erwarteten nur
// in einem Feld nicht gleicht, ist eine Abweichung, für jedes Feld nach
// LH-FA-18.a; Namen werden nicht normalisiert, ein NULL-Parameter gleicht
// keinem leeren Wert. Eine leere Liste gleicht nil.
func TestReplayExtendedFelder(t *testing.T) {
	basis := model.ClientMessage{
		Type: model.ClientBind, Statement: "s1", Portal: "p1", SQL: "SELECT $1",
		ParamTypes: []uint32{25}, ParamFormats: []int16{0}, Params: []model.Value{text("x"), {Null: true}},
		ResultFormats: []int16{1}, Target: model.TargetPortal, Name: "n", MaxRows: 5,
	}
	if f := abweichung(basis, basis); f != "" {
		t.Fatalf("gleiche Nachricht weicht in %s ab", f)
	}
	faelle := []struct {
		feld  string
		aendr func(*model.ClientMessage)
	}{
		{"type", func(m *model.ClientMessage) { m.Type = model.ClientExecute }},
		{"statement", func(m *model.ClientMessage) { m.Statement = "S1" }},
		{"portal", func(m *model.ClientMessage) { m.Portal = "p1 " }},
		{"sql", func(m *model.ClientMessage) { m.SQL = "SELECT  $1" }},
		{"param_types", func(m *model.ClientMessage) { m.ParamTypes = []uint32{23} }},
		{"param_types", func(m *model.ClientMessage) { m.ParamTypes = nil }},
		{"param_formats", func(m *model.ClientMessage) { m.ParamFormats = []int16{1} }},
		{"params", func(m *model.ClientMessage) { m.Params = []model.Value{text("y"), {Null: true}} }},
		{"params", func(m *model.ClientMessage) { m.Params = []model.Value{text("x"), {}} }},
		{"params", func(m *model.ClientMessage) { m.Params = []model.Value{text("x")} }},
		{"result_formats", func(m *model.ClientMessage) { m.ResultFormats = []int16{0} }},
		{"target", func(m *model.ClientMessage) { m.Target = model.TargetStatement }},
		{"name", func(m *model.ClientMessage) { m.Name = "" }},
		{"max_rows", func(m *model.ClientMessage) { m.MaxRows = 0 }},
	}
	for _, f := range faelle {
		m := basis
		f.aendr(&m)
		if got := abweichung(m, basis); got != f.feld {
			t.Errorf("Änderung an %s: abweichend in %q", f.feld, got)
		}
	}

	leer := model.ClientMessage{Type: model.ClientBind, ParamTypes: []uint32{}, ParamFormats: []int16{}, Params: []model.Value{{Bytes: []byte{}}}, ResultFormats: []int16{}}
	null := model.ClientMessage{Type: model.ClientBind, Params: []model.Value{{}}}
	if f := abweichung(leer, null); f != "" {
		t.Fatalf("leere Listen und leerer Wert weichen von nil ab: %s", f)
	}
}

// Abdeckung: LH-FA-10/Negative, LH-FA-18/Negative — eine abweichende
// Extended-Nachricht ist PGR-E5001 und erhält keine Antwort, auch nicht die
// Server-Nachrichten ihrer Gruppe; der Cursor bleibt stehen. Die Diagnose nennt
// Session, Interaktion, Gruppe, Nachricht und beide Nachrichtentypen, das
// abweichende Feld, bei SQL beide Texte, aber keine Parameterwerte.
func TestReplayExtendedAbweichung(t *testing.T) {
	ctx := context.Background()
	s, id := replayMit(t, []model.Interaction{vorbereitung(1)})
	in := vorbereitung(1)
	sendeAlle(t, s, id, in.Groups[0].Client...)
	if _, err := s.ClientMessage(ctx, id, bind("s1", text("a"))); err != nil {
		t.Fatal(err)
	}
	out, err := s.ClientMessage(ctx, id, syncNachricht)
	if code(err) != model.CodeReplayMismatch || out != nil {
		t.Fatalf("Sync statt Execute: %#v, %v", out, err)
	}
	for _, teil := range []string{"Session 1", "Interaktion 1", "Gruppe 2", "Nachricht 2", "erwartet execute", "empfangen sync", "abweichend in type"} {
		if !strings.Contains(err.Error(), teil) {
			t.Fatalf("Diagnose ohne %q: %v", teil, err)
		}
	}
	if out := sendeAlle(t, s, id, execute, syncNachricht); !reflect.DeepEqual(out, in.Groups[1].Server) {
		t.Fatalf("Cursor nach der Abweichung verschoben: %#v", out)
	}

	s, id = replayMit(t, []model.Interaction{vorbereitung(1)})
	_, err = s.ClientMessage(ctx, id, parse("s1", "SELECT $1::int"))
	if code(err) != model.CodeReplayMismatch || !strings.Contains(err.Error(), `erwartet "SELECT $1::text", empfangen "SELECT $1::int"`) {
		t.Fatalf("abweichender SQL-Text: %v", err)
	}
	sendeAlle(t, s, id, in.Groups[0].Client...)
	_, err = s.ClientMessage(ctx, id, bind("s1", text("geheim-wert")))
	var me *model.Error
	if !errors.As(err, &me) || me.Code != model.CodeReplayMismatch ||
		me.Msg != "Session 1, Interaktion 1, Gruppe 2, Nachricht 1: erwartet bind, empfangen bind, abweichend in params" || me.Err != nil {
		t.Fatalf("abweichender Parameterwert: die Diagnose nennt mehr als Stelle, Typen und Feld: %v", err)
	}
}

// Abdeckung: LH-FA-10/Negative — steht am Cursor eine Extended-Interaktion,
// mitten in ihr oder an ihrem Anfang, ist eine einfache Anfrage eine
// Abweichung, die die erwartete Nachricht nennt; steht dort eine einfache
// Anfrage oder keine Interaktion mehr, ist es eine Extended-Nachricht.
// Die erste Extended-Nachricht einer Verbindung ordnet ihr eine Session zu;
// ist keine mehr frei, ist das PGR-E5003.
func TestReplayExtendedFalscheArt(t *testing.T) {
	ctx := context.Background()
	s, id := replayMit(t, []model.Interaction{vorbereitung(1), interaktion(2, "SELECT 7", "SELECT 1")})
	if _, err := s.ClientMessage(ctx, id, parse("s1", "SELECT $1::text")); err != nil {
		t.Fatal(err)
	}
	_, err := s.Query(ctx, id, "SELECT 7")
	if code(err) != model.CodeReplayMismatch || !strings.Contains(err.Error(), "Gruppe 1, Nachricht 2: erwartet Client-Nachricht describe") {
		t.Fatalf("Anfrage mitten in der Extended-Interaktion: %v", err)
	}
	in := vorbereitung(1)
	sendeAlle(t, s, id, in.Groups[0].Client[1:]...)
	sendeAlle(t, s, id, in.Groups[1].Client...)
	_, err = s.ClientMessage(ctx, id, syncNachricht)
	if code(err) != model.CodeReplayMismatch || !strings.Contains(err.Error(), `erwartet Anfrage "SELECT 7"`) {
		t.Fatalf("Extended-Nachricht statt einfacher Anfrage: %v", err)
	}
	if _, err := s.Query(ctx, id, "SELECT 7"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ClientMessage(ctx, id, syncNachricht); code(err) != model.CodeReplayMismatch {
		t.Fatalf("Extended-Nachricht nach dem Ende der Aufzeichnung: %v", err)
	}

	s, a := replayMit(t, []model.Interaction{ausfuehrung(1, "s1", "a", "A")})
	b, _ := s.OpenConnection(ctx)
	if _, err := s.ClientMessage(ctx, b, bind("s1", text("a"))); err != nil {
		t.Fatalf("erste Extended-Nachricht der zweiten Verbindung erhält Session 1: %v", err)
	}
	if _, err := s.ClientMessage(ctx, a, bind("s1", text("a"))); code(err) != model.CodeReplaySession {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeReplaySession, err)
	}
	if w := s.CloseConnection(ctx, b); w == nil || w.Code != model.CodeUnconsumed {
		t.Fatalf("mitten in der Interaktion beendet ohne Warnung: %v", w)
	}
}

// Abdeckung: LH-FA-11/Happy — Extended: eine aufgezeichnete ErrorResponse
// mitten in einer Gruppe wird samt den danach verworfenen Client-Nachrichten
// bis zum Sync reproduziert; die verworfenen Nachrichten werden verglichen wie
// jede andere, und die Gruppe antwortet nach dem Sync mit den aufgezeichneten
// Server-Nachrichten, ohne eigene Fehlerlogik.
func TestReplayExtendedFehlerantwort(t *testing.T) {
	fehler := model.Response{Type: model.ResponseErrorResponse, Fields: map[string]string{"S": "ERROR", "C": "22012", "M": "division by zero"}}
	in := model.Interaction{Sequence: 1, Request: model.Request{Type: model.RequestExtended}, Groups: []model.Group{{
		Client: []model.ClientMessage{parse("", "SELECT 1/$1::int"), bind("", text("0")), execute, parse("", "SELECT 2"), bind(""), execute, syncNachricht},
		Server: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete}, fehler, rfq},
	}}}
	s, id := replayMit(t, []model.Interaction{in})
	if out := sendeAlle(t, s, id, in.Groups[0].Client...); !reflect.DeepEqual(out, in.Groups[0].Server) {
		t.Fatalf("Gruppe mit Fehler: %#v", out)
	}

	s, id = replayMit(t, []model.Interaction{in})
	sendeAlle(t, s, id, in.Groups[0].Client[:3]...)
	if _, err := s.ClientMessage(context.Background(), id, parse("", "SELECT 3")); code(err) != model.CodeReplayMismatch {
		t.Fatalf("abweichende verworfene Nachricht: %v", err)
	}
}

// Abdeckung: LH-FA-13/Boundary — Replay: beim Herunterfahren darf eine
// Verbindung enden, solange ihr Cursor nicht innerhalb einer
// Extended-Interaktion steht: vor der ersten Nachricht und nach dem Sync ja,
// nach der ersten Nachricht und zwischen zwei Gruppen nein.
func TestReplayExtendedHerunterfahren(t *testing.T) {
	ctx := context.Background()
	s, id := replayMit(t, []model.Interaction{vorbereitung(1)})
	if !s.Shutdown(ctx, id) || !s.Shutdown(ctx, 99) {
		t.Fatal("ohne Anfrage kein Ende freigegeben")
	}
	in := vorbereitung(1)
	schritte := []struct {
		msg  model.ClientMessage
		ende bool
	}{
		{in.Groups[0].Client[0], false}, {in.Groups[0].Client[1], false}, {in.Groups[0].Client[2], false},
		{in.Groups[1].Client[0], false}, {in.Groups[1].Client[1], false}, {in.Groups[1].Client[2], true},
	}
	for i, sch := range schritte {
		if _, err := s.ClientMessage(ctx, id, sch.msg); err != nil {
			t.Fatal(err)
		}
		if got := s.Shutdown(ctx, id); got != sch.ende {
			t.Fatalf("nach Nachricht %d (%s): Ende freigegeben %v", i+1, sch.msg.Type, got)
		}
	}
}

// Abdeckung: LH-FA-03/Negative — eine Aufzeichnung, deren Interaktion die Form
// nach SPEC-041 verfehlt (Extended ohne Gruppe, Gruppe ohne Client-Nachricht),
// ist beim Start PGR-E3003, statt im Replay an der fehlenden Gruppe zu enden.
func TestReplayStartformExtended(t *testing.T) {
	ctx := context.Background()
	for name, in := range map[string]model.Interaction{
		"ohne Gruppe":           {Sequence: 1, Request: model.Request{Type: model.RequestExtended}},
		"ohne Client-Nachricht": {Sequence: 1, Request: model.Request{Type: model.RequestExtended}, Groups: []model.Group{{Server: []model.Response{rfq}}}},
	} {
		if _, err := NewReplayService(ctx, ladeRepo{rec: aufzeichnung([]model.Interaction{in})}, "rec.yaml"); code(err) != model.CodeRecordingBroken {
			t.Errorf("%s: erwartet %s, erhalten %v", name, model.CodeRecordingBroken, err)
		}
	}
}
