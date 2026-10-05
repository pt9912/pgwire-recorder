package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

var (
	cParse    = model.ClientMessage{Type: model.ClientParse, Statement: "s1", SQL: "SELECT $1"}
	cBind     = model.ClientMessage{Type: model.ClientBind, Statement: "s1", Params: []model.Value{{Bytes: []byte("1")}}}
	cDescribe = model.ClientMessage{Type: model.ClientDescribe, Target: model.TargetPortal}
	cExecute  = model.ClientMessage{Type: model.ClientExecute}
	cFlush    = model.ClientMessage{Type: model.ClientFlush}
	cSync     = model.ClientMessage{Type: model.ClientSync}

	sParse = model.Response{Type: model.ResponseParseComplete}
	sBind  = model.Response{Type: model.ResponseBindComplete}
	sRow   = model.Response{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte("1")}}}
	sCmd   = model.Response{Type: model.ResponseCommandComplete, Tag: "SELECT 1"}
	sRFQ   = model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"}
	sErr   = model.Response{Type: model.ResponseErrorResponse, Fields: map[string]string{"C": "22012"}}
)

func codeOf(err error) string {
	var me *model.Error
	if errors.As(err, &me) {
		return me.Code
	}
	return ""
}

// client übergibt die Nachrichten der Reihe nach.
func client(t *testing.T, s *RecordService, id model.SessionID, msgs ...model.ClientMessage) {
	t.Helper()
	for _, m := range msgs {
		if err := s.ClientMessage(context.Background(), id, m); err != nil {
			t.Fatalf("ClientMessage %s: %v", m.Type, err)
		}
	}
}

// server übergibt die Nachrichten der Reihe nach und liefert, ob die letzte die
// Interaktion abschloss; keine davor darf es.
func server(t *testing.T, s *RecordService, id model.SessionID, rs ...model.Response) bool {
	t.Helper()
	ende := false
	for i, r := range rs {
		var err error
		ende, err = s.ServerMessage(context.Background(), id, r)
		if err != nil {
			t.Fatalf("ServerMessage %s: %v", r.Type, err)
		}
		if ende && i != len(rs)-1 {
			t.Fatalf("Interaktion vor der letzten Nachricht abgeschlossen (bei %d)", i)
		}
	}
	return ende
}

func schliessen(t *testing.T, s *RecordService, id model.SessionID, end model.SessionEnd) {
	t.Helper()
	if err := s.CloseSession(context.Background(), id, end); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: die Client-Nachrichten einer
// Sync-Gruppe gehen erst mit dem Sync, dann in einem Send an den Upstream; mit
// dem ReadyForQuery steht die Interaktion steht sie als type extended mit
// Client- und Server-Nachrichten in Reihenfolge in der Aufzeichnung, zwischen
// einfachen Anfragen in der Reihenfolge der Session.
func TestRecordExtendedSyncGruppe(t *testing.T) {
	up := &fakeUpstream{}
	s, repo := neu(t, up)
	id := session(t, s, "SELECT 0")
	client(t, s, id, cParse, cBind, cDescribe, cExecute)
	if len(up.letzte.gesendet) != 0 {
		t.Fatalf("vor dem Sync gesendet: %v", up.letzte.gesendet)
	}
	client(t, s, id, cSync)
	if !reflect.DeepEqual(up.letzte.gesendet, [][]model.ClientMessage{{cParse, cBind, cDescribe, cExecute, cSync}}) {
		t.Fatalf("an den Upstream: %#v", up.letzte.gesendet)
	}
	if !server(t, s, id, sParse, sBind, sRow, sCmd, sRFQ) {
		t.Fatal("ReadyForQuery schließt die Interaktion nicht ab")
	}
	if _, err := s.Query(context.Background(), id, "SELECT 9"); err != nil {
		t.Fatal(err)
	}
	schliessen(t, s, id, model.EndNormal)

	got := repo.last(t).Sessions[0].Interactions
	want := model.Interaction{
		Sequence: 2,
		Request:  model.Request{Type: model.RequestExtended},
		Groups: []model.Group{{
			Client: []model.ClientMessage{cParse, cBind, cDescribe, cExecute, cSync},
			Server: []model.Response{sParse, sBind, sRow, sCmd, sRFQ},
		}},
	}
	if len(got) != 3 || !reflect.DeepEqual(got[1], want) || got[0].Request.SQL != "SELECT 0" || got[2].Request.SQL != "SELECT 9" || got[2].Sequence != 3 {
		t.Fatalf("Interaktionen: %#v", got)
	}
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: eine Flush-Gruppe nimmt die Server-Nachrichten
// auf, die vor der nächsten Client-Nachricht eintreffen; eine danach
// eintreffende gehört zur folgenden Gruppe. Jede Gruppe geht mit ihrem Flush
// beziehungsweise Sync einzeln an den Upstream, und die Interaktion besteht
// Validate.
func TestRecordExtendedFlushGruppen(t *testing.T) {
	up := &fakeUpstream{}
	s, repo := neu(t, up)
	id := session(t, s)
	client(t, s, id, cParse, cFlush)
	server(t, s, id, sParse)
	client(t, s, id, cBind)
	server(t, s, id, sErr) // trifft nach der nächsten Client-Nachricht ein
	client(t, s, id, cExecute, cFlush)
	client(t, s, id, cSync)
	if !server(t, s, id, sRFQ) {
		t.Fatal("nicht abgeschlossen")
	}
	schliessen(t, s, id, model.EndNormal)

	want := []model.Group{
		{Client: []model.ClientMessage{cParse, cFlush}, Server: []model.Response{sParse}},
		{Client: []model.ClientMessage{cBind, cExecute, cFlush}, Server: []model.Response{sErr}},
		{Client: []model.ClientMessage{cSync}, Server: []model.Response{sRFQ}},
	}
	got := repo.last(t).Sessions[0].Interactions
	if len(got) != 1 || !reflect.DeepEqual(got[0].Groups, want) {
		t.Fatalf("Gruppen: %#v", got)
	}
	wantGesendet := [][]model.ClientMessage{{cParse, cFlush}, {cBind, cExecute, cFlush}, {cSync}}
	if !reflect.DeepEqual(up.letzte.gesendet, wantGesendet) {
		t.Fatalf("an den Upstream: %#v", up.letzte.gesendet)
	}
	if err := got[0].Validate(); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: mehrere Interaktionen nacheinander, auch eine
// aus nur einem Sync, zählen fortlaufend; nach einem Fehler stehen die
// empfangenen Client-Nachrichten und die Server-Nachrichten der Gruppe in der
// Aufzeichnung.
func TestRecordExtendedMehrereInteraktionen(t *testing.T) {
	s, repo := neu(t, &fakeUpstream{})
	id := session(t, s)
	client(t, s, id, cParse, cBind, cExecute, cSync)
	server(t, s, id, sParse, sErr, sRFQ)
	client(t, s, id, cSync)
	server(t, s, id, sRFQ)
	schliessen(t, s, id, model.EndNormal)

	got := repo.last(t).Sessions[0].Interactions
	if len(got) != 2 || got[0].Sequence != 1 || got[1].Sequence != 2 {
		t.Fatalf("Interaktionen: %#v", got)
	}
	if !reflect.DeepEqual(got[0].Groups[0].Server, []model.Response{sParse, sErr, sRFQ}) ||
		!reflect.DeepEqual(got[1].Groups, []model.Group{{Client: []model.ClientMessage{cSync}, Server: []model.Response{sRFQ}}}) {
		t.Fatalf("Gruppen: %#v", got)
	}
}

// Eine Extended-Interaktion ohne
// ReadyForQuery wird beim Ende der Session nicht übernommen, die vorherige
// bleibt; erreicht das ReadyForQuery einer abgeschlossenen den Client nicht
// (EndLost), entfällt sie.
func TestRecordExtendedAbbruch(t *testing.T) {
	t.Run("offen", func(t *testing.T) {
		s, repo := neu(t, &fakeUpstream{})
		id := session(t, s, "SELECT 1")
		client(t, s, id, cParse, cFlush)
		server(t, s, id, sParse)
		schliessen(t, s, id, model.EndNormal)
		got := repo.last(t).Sessions[0].Interactions
		if len(got) != 1 || got[0].Request.SQL != "SELECT 1" {
			t.Fatalf("Interaktionen: %#v", got)
		}
	})
	t.Run("ReadyForQuery nicht zugestellt", func(t *testing.T) {
		s, repo := neu(t, &fakeUpstream{})
		id := session(t, s, "SELECT 1")
		client(t, s, id, cSync)
		server(t, s, id, sRFQ)
		schliessen(t, s, id, model.EndLost)
		got := repo.last(t).Sessions[0].Interactions
		if len(got) != 1 || got[0].Request.SQL != "SELECT 1" {
			t.Fatalf("Interaktionen: %#v", got)
		}
	})
}

// Abdeckung: LH-FA-06/Negative — eine einfache Anfrage
// während einer laufenden Extended-Interaktion und eine Interaktion, die
// Validate ablehnt (ein ReadyForQuery ohne Sync, eine Antwort des Extended
// Query Protocol in einer einfachen Anfrage), sind PGR-E6001; die Session wird
// nicht übernommen, und die Antworten der einfachen Anfrage gehen nicht zurück.
func TestRecordExtendedNichtDarstellbar(t *testing.T) {
	t.Run("Query in laufender Interaktion", func(t *testing.T) {
		up := &fakeUpstream{}
		s, repo := neu(t, up)
		id := session(t, s, "SELECT 1")
		client(t, s, id, cParse, cFlush)
		out, err := s.Query(context.Background(), id, "SELECT 2")
		if codeOf(err) != model.CodeUnsupported || out != nil {
			t.Fatalf("Antworten %v, Fehler %v", out, err)
		}
		schliessen(t, s, id, model.EndNormal)
		if len(repo.writes) != 0 {
			t.Fatalf("Session übernommen: %#v", repo.writes)
		}
	})
	t.Run("ReadyForQuery ohne Sync", func(t *testing.T) {
		s, repo := neu(t, &fakeUpstream{})
		id := session(t, s, "SELECT 1")
		client(t, s, id, cParse, cFlush)
		_, err := s.ServerMessage(context.Background(), id, sRFQ)
		if codeOf(err) != model.CodeUnsupported || !strings.Contains(err.Error(), "sync") {
			t.Fatalf("Fehler: %v", err)
		}
		schliessen(t, s, id, model.EndNormal)
		if len(repo.writes) != 0 {
			t.Fatalf("Session übernommen: %#v", repo.writes)
		}
	})
	t.Run("Extended-Antwort in einfacher Anfrage", func(t *testing.T) {
		s, repo := neu(t, &fakeUpstream{})
		id := session(t, s, "SELECT 1")
		out, err := s.Query(context.Background(), id, "EXTENDED-ANTWORT")
		if codeOf(err) != model.CodeUnsupported || out != nil || !strings.Contains(err.Error(), "parse_complete") {
			t.Fatalf("Antworten %v, Fehler %v", out, err)
		}
		schliessen(t, s, id, model.EndNormal)
		if len(repo.writes) != 0 {
			t.Fatalf("Session übernommen: %#v", repo.writes)
		}
	})
}

// Eine Client-Nachricht nach dem Sync vor dessen ReadyForQuery und eine
// Server-Nachricht ohne gesendete Gruppe sind Fehler des Aufrufers (PGR-E1000);
// ein Sendefehler des Upstreams geht unverändert zurück.
func TestRecordExtendedAufruferfehler(t *testing.T) {
	ctx := context.Background()
	s, _ := neu(t, &fakeUpstream{})
	id := session(t, s)
	if _, err := s.ServerMessage(ctx, id, sParse); codeOf(err) != model.CodeInternal {
		t.Fatalf("ohne Interaktion: %v", err)
	}
	client(t, s, id, cParse)
	if _, err := s.ServerMessage(ctx, id, sParse); codeOf(err) != model.CodeInternal {
		t.Fatalf("ohne gesendete Gruppe: %v", err)
	}
	client(t, s, id, cSync)
	if err := s.ClientMessage(ctx, id, cParse); codeOf(err) != model.CodeInternal {
		t.Fatalf("nach dem Sync: %v", err)
	}

	up := &fakeUpstream{}
	s2, _ := neu(t, up)
	id2 := session(t, s2)
	up.letzte.sendErr = model.Errorf(model.CodeConnectionLost, nil, "weg")
	if err := s2.ClientMessage(ctx, id2, cSync); codeOf(err) != model.CodeConnectionLost {
		t.Fatalf("Sendefehler: %v", err)
	}
}

// AwaitServer liefert, was der Upstream empfängt, und ändert die Aufzeichnung
// nicht.
func TestRecordAwaitServer(t *testing.T) {
	up := &fakeUpstream{}
	s, _ := neu(t, up)
	id := session(t, s)
	up.letzte.empfang = [][]model.Response{{sParse, sRFQ}}
	out, err := s.AwaitServer(context.Background(), id)
	if err != nil || !reflect.DeepEqual(out, []model.Response{sParse, sRFQ}) {
		t.Fatalf("Antworten %v, Fehler %v", out, err)
	}
	if _, err := s.AwaitServer(context.Background(), id); codeOf(err) != model.CodeConnectionLost {
		t.Fatalf("Fehler: %v", err)
	}
	if _, err := s.AwaitServer(context.Background(), 99); codeOf(err) != model.CodeInternal {
		t.Fatalf("unbekannte Session: %v", err)
	}
}
