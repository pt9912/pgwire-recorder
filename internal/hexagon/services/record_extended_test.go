package services_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// Testdaten: c… sind Client-, s… Server-Nachrichten.
func cParse() model.ClientMessage {
	return model.ClientMessage{Type: model.ClientParse, Statement: "s1", SQL: "SELECT $1"}
}
func cBind() model.ClientMessage {
	return model.ClientMessage{Type: model.ClientBind, Statement: "s1", Params: []model.Value{{Bytes: []byte("1")}}}
}
func cDescribe() model.ClientMessage {
	return model.ClientMessage{Type: model.ClientDescribe, Target: model.TargetPortal}
}
func cExecute() model.ClientMessage { return model.ClientMessage{Type: model.ClientExecute} }
func cFlush() model.ClientMessage   { return model.ClientMessage{Type: model.ClientFlush} }
func cSync() model.ClientMessage    { return model.ClientMessage{Type: model.ClientSync} }
func sParse() model.Response        { return model.Response{Type: model.ResponseParseComplete} }
func sBind() model.Response         { return model.Response{Type: model.ResponseBindComplete} }
func sRow() model.Response {
	return model.Response{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte("1")}}}
}
func sCmd() model.Response {
	return model.Response{Type: model.ResponseCommandComplete, Tag: "SELECT 1"}
}
func sRFQ() model.Response { return model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"} }
func sErr() model.Response {
	return model.Response{Type: model.ResponseErrorResponse, Fields: map[string]string{"C": "22012"}}
}

func codeOf(err error) string {
	var me *model.Error
	if errors.As(err, &me) {
		return me.Code
	}
	return ""
}

// warte läuft f nebenher und meldet einen Fehler, wenn es nicht binnen zwei
// Sekunden zurückkehrt.
func warte(t *testing.T, was string, f func()) {
	t.Helper()
	fertig := make(chan struct{})
	go func() {
		defer close(fertig)
		f()
	}()
	select {
	case <-fertig:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s kehrt nicht zurück", was)
	}
}

// client übergibt die Nachrichten der Reihe nach.
func client(t *testing.T, s *services.RecordService, id model.SessionID, msgs ...model.ClientMessage) {
	t.Helper()
	for _, m := range msgs {
		if err := s.ClientMessage(context.Background(), id, m); err != nil {
			t.Fatalf("ClientMessage %s: %v", m.Type, err)
		}
	}
}

// server legt die Nachrichten als eine Folge in den Upstream, holt sie mit
// AwaitServer und meldet sie als zugestellt.
func server(t *testing.T, s *services.RecordService, up *fakeUpstream, id model.SessionID, rs ...model.Response) {
	t.Helper()
	up.letzte.empfang <- rs
	var got []model.Response
	var err error
	warte(t, "AwaitServer", func() { got, err = s.AwaitServer(context.Background(), id) })
	if err != nil || !reflect.DeepEqual(got, rs) {
		t.Fatalf("AwaitServer: %#v, %v", got, err)
	}
	s.Delivered(context.Background(), id)
}

func schliessen(t *testing.T, s *services.RecordService, id model.SessionID, end model.SessionEnd) {
	t.Helper()
	if err := s.CloseSession(context.Background(), id, end); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: die Client-Nachrichten einer
// Sync-Gruppe gehen erst mit dem Sync, dann in einem Send an den Upstream; mit
// dem ReadyForQuery steht die Interaktion als type extended mit Client- und
// Server-Nachrichten in Reihenfolge in der Aufzeichnung, zwischen einfachen
// Anfragen in der Reihenfolge der Session.
func TestRecordExtendedSyncGruppe(t *testing.T) {
	up := &fakeUpstream{}
	s, repo := neu(t, up)
	id := session(t, s, "SELECT 0")
	client(t, s, id, cParse(), cBind(), cDescribe(), cExecute())
	if g := up.letzte.gruppen(); len(g) != 0 {
		t.Fatalf("vor dem Sync gesendet: %v", g)
	}
	client(t, s, id, cSync())
	if g := up.letzte.gruppen(); !reflect.DeepEqual(g, [][]model.ClientMessage{{cParse(), cBind(), cDescribe(), cExecute(), cSync()}}) {
		t.Fatalf("an den Upstream: %#v", g)
	}
	server(t, s, up, id, sParse(), sBind(), sRow(), sCmd(), sRFQ())
	if _, err := s.Query(context.Background(), id, "SELECT 9"); err != nil {
		t.Fatal(err)
	}
	s.Delivered(context.Background(), id)
	schliessen(t, s, id, model.EndClosed)

	got := repo.last(t).Sessions[0].Interactions
	want := model.Interaction{
		Sequence: 2,
		Request:  model.Request{Type: model.RequestExtended},
		Groups: []model.Group{{
			Client: []model.ClientMessage{cParse(), cBind(), cDescribe(), cExecute(), cSync()},
			Server: []model.Response{sParse(), sBind(), sRow(), sCmd(), sRFQ()},
		}},
	}
	if len(got) != 3 || !reflect.DeepEqual(got[1], want) || got[0].Request.SQL != "SELECT 0" || got[2].Request.SQL != "SELECT 9" || got[2].Sequence != 3 {
		t.Fatalf("Interaktionen: %#v", got)
	}
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: eine Flush-Gruppe nimmt die
// Server-Nachrichten auf, die vor der nächsten Client-Nachricht eintreffen; eine
// danach eintreffende gehört zur folgenden Gruppe. Jede Gruppe geht mit ihrem
// Flush beziehungsweise Sync einzeln an den Upstream, und die Interaktion
// besteht Validate.
func TestRecordExtendedFlushGruppen(t *testing.T) {
	up := &fakeUpstream{}
	s, repo := neu(t, up)
	id := session(t, s)
	client(t, s, id, cParse(), cFlush())
	server(t, s, up, id, sParse())
	client(t, s, id, cBind())
	server(t, s, up, id, sErr()) // trifft nach der nächsten Client-Nachricht ein
	client(t, s, id, cExecute(), cFlush())
	client(t, s, id, cSync())
	server(t, s, up, id, sRFQ())
	schliessen(t, s, id, model.EndClosed)

	want := []model.Group{
		{Client: []model.ClientMessage{cParse(), cFlush()}, Server: []model.Response{sParse()}},
		{Client: []model.ClientMessage{cBind(), cExecute(), cFlush()}, Server: []model.Response{sErr()}},
		{Client: []model.ClientMessage{cSync()}, Server: []model.Response{sRFQ()}},
	}
	got := repo.last(t).Sessions[0].Interactions
	if len(got) != 1 || !reflect.DeepEqual(got[0].Groups, want) {
		t.Fatalf("Gruppen: %#v", got)
	}
	wantGesendet := [][]model.ClientMessage{{cParse(), cFlush()}, {cBind(), cExecute(), cFlush()}, {cSync()}}
	if g := up.letzte.gruppen(); !reflect.DeepEqual(g, wantGesendet) {
		t.Fatalf("an den Upstream: %#v", g)
	}
	if err := got[0].Validate(); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: Client-Nachrichten nach einem Sync
// beginnen die nächste Interaktion, auch bevor dessen ReadyForQuery eintrifft
// (Pipelining); die Server-Nachrichten gehören bis zum ReadyForQuery der
// ersten. Auch eine Interaktion aus nur einem Sync zählt; nach einem Fehler
// stehen die empfangenen Client-Nachrichten und die Server-Nachrichten der
// Gruppe in der Aufzeichnung.
func TestRecordExtendedPipelining(t *testing.T) {
	up := &fakeUpstream{}
	s, repo := neu(t, up)
	id := session(t, s)
	client(t, s, id, cParse(), cBind(), cExecute(), cSync(), cSync())
	server(t, s, up, id, sParse(), sErr(), sRFQ())
	server(t, s, up, id, sRFQ())
	schliessen(t, s, id, model.EndClosed)

	got := repo.last(t).Sessions[0].Interactions
	if len(got) != 2 || got[0].Sequence != 1 || got[1].Sequence != 2 {
		t.Fatalf("Interaktionen: %#v", got)
	}
	if !reflect.DeepEqual(got[0].Groups, []model.Group{{Client: []model.ClientMessage{cParse(), cBind(), cExecute(), cSync()}, Server: []model.Response{sParse(), sErr(), sRFQ()}}}) ||
		!reflect.DeepEqual(got[1].Groups, []model.Group{{Client: []model.ClientMessage{cSync()}, Server: []model.Response{sRFQ()}}}) {
		t.Fatalf("Gruppen: %#v", got)
	}
}

// Abdeckung: LH-FA-13/Negative — das gemeldete Ende entscheidet der Service:
// Verbindungsende oder Terminate während einer laufenden Interaktion ist
// PGR-E4003, danach kein Fehler; ein Schreibfehler zum Client ist PGR-E4003, und
// die nicht zugestellte Interaktion entfällt; Herunterfahren beendet ohne
// Fehler. Eine laufende Interaktion wird nie übernommen, die vorherigen
// bleiben.
func TestRecordExtendedEnde(t *testing.T) {
	cases := []struct {
		name    string
		ablauf  func(s *services.RecordService, up *fakeUpstream, id model.SessionID)
		end     model.SessionEnd
		wantErr string
		wantSQL []string
	}{
		{"Verbindungsende in laufender Interaktion", func(s *services.RecordService, up *fakeUpstream, id model.SessionID) {
			client(t, s, id, cParse(), cFlush())
			server(t, s, up, id, sParse())
		}, model.EndClosed, model.CodeConnectionLost, []string{"SELECT 1"}},
		{"Terminate in laufender Interaktion", func(s *services.RecordService, up *fakeUpstream, id model.SessionID) {
			client(t, s, id, cParse())
		}, model.EndTerminate, model.CodeConnectionLost, []string{"SELECT 1"}},
		{"Verbindungsende danach", func(s *services.RecordService, up *fakeUpstream, id model.SessionID) {
			client(t, s, id, cSync())
			server(t, s, up, id, sRFQ())
		}, model.EndClosed, "", []string{"SELECT 1", ""}},
		{"ReadyForQuery nicht zugestellt", func(s *services.RecordService, up *fakeUpstream, id model.SessionID) {
			client(t, s, id, cSync())
			up.letzte.empfang <- []model.Response{sRFQ()}
			if _, err := s.AwaitServer(context.Background(), id); err != nil {
				t.Fatal(err)
			}
		}, model.EndWriteFailed, model.CodeConnectionLost, []string{"SELECT 1"}},
		{"Herunterfahren mit offener Interaktion", func(s *services.RecordService, up *fakeUpstream, id model.SessionID) {
			client(t, s, id, cParse())
		}, model.EndShutdown, "", []string{"SELECT 1"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			up := &fakeUpstream{}
			s, repo := neu(t, up)
			id := session(t, s, "SELECT 1")
			c.ablauf(s, up, id)
			if err := s.CloseSession(context.Background(), id, c.end); codeOf(err) != c.wantErr {
				t.Fatalf("Fehler %v, erwartet %q", err, c.wantErr)
			}
			var sqls []string
			for _, in := range repo.last(t).Sessions[0].Interactions {
				sqls = append(sqls, in.Request.SQL)
			}
			if !reflect.DeepEqual(sqls, c.wantSQL) {
				t.Fatalf("Interaktionen %q, erwartet %q", sqls, c.wantSQL)
			}
		})
	}
}

// Abdeckung: LH-FA-06/Negative — eine einfache Anfrage während einer laufenden
// Extended-Interaktion vor deren Sync und eine Interaktion, die Validate
// ablehnt (ein ReadyForQuery ohne Sync, eine Antwort des Extended Query
// Protocol in einer einfachen Anfrage), sind PGR-E6001; die Session wird nicht
// übernommen, und die Antworten der einfachen Anfrage gehen nicht zurück.
func TestRecordExtendedNichtDarstellbar(t *testing.T) {
	t.Run("Query in laufender Interaktion", func(t *testing.T) {
		s, repo := neu(t, &fakeUpstream{})
		id := session(t, s, "SELECT 1")
		client(t, s, id, cParse(), cFlush())
		out, err := s.Query(context.Background(), id, "SELECT 2")
		if codeOf(err) != model.CodeUnsupported || out != nil {
			t.Fatalf("Antworten %v, Fehler %v", out, err)
		}
		schliessen(t, s, id, model.EndFailed)
		if len(repo.writes) != 0 {
			t.Fatalf("Session übernommen: %#v", repo.writes)
		}
	})
	t.Run("ReadyForQuery ohne Sync", func(t *testing.T) {
		up := &fakeUpstream{}
		s, repo := neu(t, up)
		id := session(t, s, "SELECT 1")
		client(t, s, id, cParse(), cFlush())
		up.letzte.empfang <- []model.Response{sRFQ()}
		_, err := s.AwaitServer(context.Background(), id)
		if codeOf(err) != model.CodeUnsupported || !strings.Contains(err.Error(), "sync") {
			t.Fatalf("Fehler: %v", err)
		}
		schliessen(t, s, id, model.EndFailed)
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
		schliessen(t, s, id, model.EndFailed)
		if len(repo.writes) != 0 {
			t.Fatalf("Session übernommen: %#v", repo.writes)
		}
	})
}

// Abdeckung: LH-FA-11/Negative — liefert der Empfang vom Upstream in einer
// Extended-Interaktion PGR-E6001 (Fehlerantwort vor dem Abbruch), gibt
// AwaitServer den Fehler weiter, und die Session wird nicht übernommen, auch
// nicht ihre vorherige Interaktion.
func TestRecordExtendedFehlerantwortVorDemAbbruch(t *testing.T) {
	up := &fakeUpstream{}
	s, repo := neu(t, up)
	id := session(t, s, "SELECT 1")
	up.letzte.receiveErr = model.Errorf(model.CodeUnsupported, nil, "Fehlerantwort 57P01 vor dem Abbruch")
	client(t, s, id, cParse(), cFlush())
	if _, err := s.AwaitServer(context.Background(), id); codeOf(err) != model.CodeUnsupported {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeUnsupported, err)
	}
	schliessen(t, s, id, model.EndFailed)
	if len(repo.writes) != 0 {
		t.Fatalf("Session übernommen: %#v", repo.writes)
	}
}

// Abdeckung: LH-FA-11/Negative — hat die Server-Richtung eine ErrorResponse
// gelesen und scheitern danach Senden und Empfangen mit PGR-E6001, wird die
// Session nicht übernommen, gleich welche Richtung das Ende zuerst meldet:
// meldet es die Client-Richtung, ist die Session nicht übernehmbar, bevor
// ClientMessage den Fehler liefert und CloseSession läuft; meldet es die
// Server-Richtung, liefert ein späteres ClientMessage ErrSessionEnded.
func TestRecordExtendedSendefehlerNachFehlerantwort(t *testing.T) {
	e6001 := model.Errorf(model.CodeUnsupported, nil, "Fehlerantwort 57P01 vor dem Abbruch")
	vorbereiten := func(t *testing.T) (*services.RecordService, *fakeRepo, *fakeUpstream, model.SessionID) {
		t.Helper()
		up := &fakeUpstream{}
		s, repo := neu(t, up)
		id := session(t, s, "SELECT 1")
		client(t, s, id, cParse(), cFlush())
		server(t, s, up, id, sErr())
		up.letzte.mu.Lock()
		up.letzte.sendErr = e6001
		up.letzte.mu.Unlock()
		return s, repo, up, id
	}
	t.Run("Client-Richtung zuerst", func(t *testing.T) {
		s, repo, up, id := vorbereiten(t)
		empfangen := make(chan error, 1)
		go func() {
			_, err := s.AwaitServer(context.Background(), id)
			empfangen <- err
		}()
		if err := s.ClientMessage(context.Background(), id, cSync()); codeOf(err) != model.CodeUnsupported {
			t.Fatalf("ClientMessage: erwartet %s, erhalten %v", model.CodeUnsupported, err)
		}
		if err := s.CloseSession(context.Background(), id, model.EndFailed); err != nil {
			t.Fatal(err)
		}
		up.letzte.empfangsFehler <- e6001
		if err := <-empfangen; err != model.ErrSessionEnded {
			t.Fatalf("AwaitServer nach CloseSession: %v", err)
		}
		if len(repo.writes) != 0 {
			t.Fatalf("Session übernommen: %#v", repo.writes)
		}
	})
	t.Run("Server-Richtung zuerst", func(t *testing.T) {
		s, repo, up, id := vorbereiten(t)
		up.letzte.empfangsFehler <- e6001
		if _, err := s.AwaitServer(context.Background(), id); codeOf(err) != model.CodeUnsupported {
			t.Fatalf("AwaitServer: erwartet %s, erhalten %v", model.CodeUnsupported, err)
		}
		if err := s.CloseSession(context.Background(), id, model.EndFailed); err != nil {
			t.Fatal(err)
		}
		if err := s.ClientMessage(context.Background(), id, cSync()); err != model.ErrSessionEnded {
			t.Fatalf("ClientMessage nach CloseSession: %v", err)
		}
		if len(repo.writes) != 0 {
			t.Fatalf("Session übernommen: %#v", repo.writes)
		}
	})
}

// Eine Server-Nachricht nach dem ReadyForQuery der letzten Interaktion liest
// AwaitServer nicht: Es wartet auf eine gesendete Gruppe, bis CloseSession es
// mit ErrSessionEnded beendet, und die Session trägt nur ihre Interaktionen
// (LH-FA-05.a §Nachrichten zwischen Interaktionen).
func TestRecordNachDerLetztenInteraktion(t *testing.T) {
	up := &fakeUpstream{}
	s, repo := neu(t, up)
	id := session(t, s, "SELECT 1")
	client(t, s, id, cParse(), cSync())
	server(t, s, up, id, sParse(), sRFQ())
	hinweis := model.Response{Type: model.ResponseNoticeResponse, Fields: map[string]string{"M": "danach"}}
	up.letzte.empfang <- []model.Response{hinweis}
	empfangen := make(chan error, 1)
	go func() {
		_, err := s.AwaitServer(context.Background(), id)
		empfangen <- err
	}()
	select {
	case err := <-empfangen:
		t.Fatalf("AwaitServer ohne laufende Interaktion zurückgekehrt: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	schliessen(t, s, id, model.EndClosed)
	if err := <-empfangen; err != model.ErrSessionEnded {
		t.Fatalf("AwaitServer nach CloseSession: %v", err)
	}
	got := repo.last(t).Sessions[0].Interactions
	if len(got) != 2 || !reflect.DeepEqual(got[1].Groups[0].Server, []model.Response{sParse(), sRFQ()}) {
		t.Fatalf("Interaktionen: %#v", got)
	}
}

// Eine einfache Anfrage nach dem Sync einer laufenden Extended-Interaktion
// wartet, bis deren Antworten zugestellt sind, und läuft nie gleichzeitig mit
// dem Empfang der Gegenrichtung; sie wird die nächste Interaktion.
func TestRecordQueryWartetAufExtended(t *testing.T) {
	up := &fakeUpstream{}
	s, repo := neu(t, up)
	id := session(t, s)
	client(t, s, id, cSync())
	fertig := make(chan error, 1)
	go func() {
		_, err := s.Query(context.Background(), id, "SELECT 2")
		fertig <- err
	}()
	time.Sleep(50 * time.Millisecond)
	select {
	case err := <-fertig:
		t.Fatalf("Query lief vor dem ReadyForQuery: %v", err)
	default:
	}
	up.letzte.empfang <- []model.Response{sRFQ()}
	if _, err := s.AwaitServer(context.Background(), id); err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
	select {
	case err := <-fertig:
		t.Fatalf("Query lief vor der Zustellung: %v", err)
	default:
	}
	s.Delivered(context.Background(), id)
	select {
	case err := <-fertig:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Query läuft nach der Zustellung nicht")
	}
	s.Delivered(context.Background(), id)
	// AwaitServer wartet nun, statt neben einer einfachen Anfrage zu empfangen.
	gewartet := make(chan error, 1)
	go func() {
		_, err := s.AwaitServer(context.Background(), id)
		gewartet <- err
	}()
	time.Sleep(50 * time.Millisecond)
	select {
	case err := <-gewartet:
		t.Fatalf("AwaitServer ohne gesendete Gruppe: %v", err)
	default:
	}
	schliessen(t, s, id, model.EndClosed)
	if err := <-gewartet; !errors.Is(err, model.ErrSessionEnded) {
		t.Fatalf("AwaitServer nach CloseSession: %v", err)
	}
	got := repo.last(t).Sessions[0].Interactions
	if len(got) != 2 || got[1].Request.SQL != "SELECT 2" {
		t.Fatalf("Interaktionen: %#v", got)
	}
}

// Abdeckung: LH-FA-13/Boundary — nach Shutdown nimmt die laufende
// Extended-Interaktion noch ihre Nachrichten bis zum Sync an, eine neue
// Interaktion (Extended oder einfach) beginnt nicht (ErrShutdown, nichts geht an
// den Upstream). Shutdown gibt das Ende frei, wenn keine Interaktion läuft und
// keine Antwort auf Zustellung wartet; sonst gibt Delivered es nach dem
// ReadyForQuery frei, auch bei einer laufenden einfachen Anfrage.
func TestRecordHerunterfahren(t *testing.T) {
	ctx := context.Background()
	up := &fakeUpstream{}
	s, _ := neu(t, up)
	id := session(t, s)
	client(t, s, id, cParse())
	if s.Shutdown(ctx, id) {
		t.Fatal("Ende freigegeben, während eine Interaktion läuft")
	}
	client(t, s, id, cSync())
	// Nach dem Beginn beginnt keine neue Interaktion: weder Extended noch
	// einfach, und nichts geht an den Upstream.
	if err := s.ClientMessage(ctx, id, cParse()); !errors.Is(err, model.ErrShutdown) {
		t.Fatalf("neue Extended-Interaktion: %v", err)
	}
	if _, err := s.Query(ctx, id, "SELECT 2"); !errors.Is(err, model.ErrShutdown) {
		t.Fatalf("neue einfache Anfrage: %v", err)
	}
	if g := up.letzte.gruppen(); len(g) != 1 {
		t.Fatalf("an den Upstream: %#v", g)
	}
	up.letzte.empfang <- []model.Response{sParse(), sRFQ()}
	if _, err := s.AwaitServer(ctx, id); err != nil {
		t.Fatal(err)
	}
	if s.Shutdown(ctx, id) {
		t.Fatal("Ende freigegeben, bevor das ReadyForQuery zugestellt ist")
	}
	if !s.Delivered(ctx, id) {
		t.Fatal("Ende nach dem ReadyForQuery nicht freigegeben")
	}
	if !s.Shutdown(ctx, id) {
		t.Fatal("Ende ohne laufende Interaktion nicht freigegeben")
	}
	s2, _ := neu(t, &fakeUpstream{})
	id2 := session(t, s2)
	if s2.Delivered(ctx, id2) {
		t.Fatal("Ende ohne Herunterfahren freigegeben")
	}

	// Eine laufende einfache Anfrage hält das Ende ebenso auf.
	up3 := &fakeUpstream{}
	s3, _ := neu(t, up3)
	id3 := session(t, s3)
	up3.letzte.queryLaeuft = make(chan struct{})
	up3.letzte.queryHalt = make(chan struct{})
	fertig := make(chan error, 1)
	go func() {
		_, err := s3.Query(ctx, id3, "SELECT 1")
		fertig <- err
	}()
	<-up3.letzte.queryLaeuft
	if s3.Shutdown(ctx, id3) {
		t.Fatal("Ende freigegeben, während eine einfache Anfrage läuft")
	}
	close(up3.letzte.queryHalt)
	if err := <-fertig; err != nil {
		t.Fatal(err)
	}
	if !s3.Delivered(ctx, id3) {
		t.Fatal("Ende nach der Zustellung der einfachen Anfrage nicht freigegeben")
	}
}

// Fehler des Aufrufers: AwaitServer, bevor die vorigen Nachrichten zugestellt
// sind, ist PGR-E1000; ein Sendefehler des Upstreams geht unverändert zurück;
// nach CloseSession liefert jeder Aufruf ErrSessionEnded.
func TestRecordExtendedAufruferfehler(t *testing.T) {
	ctx := context.Background()
	up := &fakeUpstream{}
	s, _ := neu(t, up)
	id := session(t, s)
	client(t, s, id, cFlush())
	up.letzte.empfang <- []model.Response{sParse()}
	if _, err := s.AwaitServer(ctx, id); err != nil {
		t.Fatal(err)
	}
	if _, err := s.AwaitServer(ctx, id); codeOf(err) != model.CodeInternal {
		t.Fatalf("vor der Zustellung: %v", err)
	}
	up.letzte.mu.Lock()
	up.letzte.sendErr = model.Errorf(model.CodeConnectionLost, nil, "weg")
	up.letzte.mu.Unlock()
	if err := s.ClientMessage(ctx, id, cSync()); codeOf(err) != model.CodeConnectionLost {
		t.Fatalf("Sendefehler: %v", err)
	}
	schliessen(t, s, id, model.EndFailed)
	if err := s.ClientMessage(ctx, id, cSync()); !errors.Is(err, model.ErrSessionEnded) {
		t.Fatalf("ClientMessage nach dem Ende: %v", err)
	}
	if _, err := s.Query(ctx, id, "SELECT 1"); !errors.Is(err, model.ErrSessionEnded) {
		t.Fatalf("Query nach dem Ende: %v", err)
	}
	if _, err := s.AwaitServer(ctx, id); !errors.Is(err, model.ErrSessionEnded) {
		t.Fatalf("AwaitServer nach dem Ende: %v", err)
	}
	if err := s.CloseSession(ctx, id, model.EndClosed); err != nil {
		t.Fatalf("zweites CloseSession: %v", err)
	}
}

// engpassUpstream bildet einen Server mit begrenzten Puffern nach: er liest
// eine Client-Nachricht erst, wenn er die Antworten auf die vorige in seinen
// Ausgabepuffer geschrieben hat, und dieser fasst nur zwei Nachrichten. Liest
// niemand die Ausgabe, liest er nicht mehr, und Send blockiert.
type engpassUpstream struct {
	letzte *engpassSession
}

func (e *engpassUpstream) Open(context.Context, map[string]string) (driven.UpstreamSession, []model.Response, error) {
	e.letzte = &engpassSession{ein: make(chan model.ClientMessage, 2), aus: make(chan model.Response, 2), zu: make(chan struct{})}
	go e.letzte.server()
	return e.letzte, []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}, nil
}

type engpassSession struct {
	ein   chan model.ClientMessage
	aus   chan model.Response
	zu    chan struct{}
	zuEin sync.Once
}

// zeilenJeExecute ist die Zahl der Zeilen, die der Server auf ein Execute
// schreibt.
const zeilenJeExecute = 50

func (e *engpassSession) server() {
	for {
		var m model.ClientMessage
		select {
		case m = <-e.ein:
		case <-e.zu:
			return
		}
		var out []model.Response
		switch m.Type {
		case model.ClientParse:
			out = []model.Response{sParse()}
		case model.ClientBind:
			out = []model.Response{sBind()}
		case model.ClientExecute:
			for i := 0; i < zeilenJeExecute; i++ {
				out = append(out, sRow())
			}
			out = append(out, sCmd())
		case model.ClientSync:
			out = []model.Response{sRFQ()}
		}
		for _, r := range out {
			select {
			case e.aus <- r:
			case <-e.zu:
				return
			}
		}
	}
}

func (e *engpassSession) Query(context.Context, string) ([]model.Response, error) {
	return nil, errors.New("nicht vorgesehen")
}

func (e *engpassSession) Send(_ context.Context, msgs []model.ClientMessage) error {
	for _, m := range msgs {
		select {
		case e.ein <- m:
		case <-e.zu:
			return model.Errorf(model.CodeConnectionLost, nil, "geschlossen")
		}
	}
	return nil
}

func (e *engpassSession) Receive(context.Context) ([]model.Response, error) {
	select {
	case r := <-e.aus:
		return []model.Response{r}, nil
	case <-e.zu:
		return nil, model.Errorf(model.CodeConnectionLost, nil, "geschlossen")
	}
}

func (e *engpassSession) Close() error {
	e.zuEin.Do(func() { close(e.zu) })
	return nil
}

// grosseGruppe ist eine Sync-Gruppe mit vielen Ausführungen, deren Ausgabe
// die Puffer des engpassUpstream weit übersteigt.
func grosseGruppe() []model.ClientMessage {
	g := []model.ClientMessage{cParse()}
	for i := 0; i < 20; i++ {
		g = append(g, cBind(), cExecute())
	}
	return append(g, cSync())
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: eine Gruppe, deren Ausgabe die
// Puffer des Servers übersteigt, läuft durch, weil AwaitServer empfängt,
// während ClientMessage noch an einen Server sendet, der erst nach dem Lesen
// seiner Ausgabe weiterliest; ClientMessage hält dabei keine Sperre, die
// AwaitServer aufhält.
func TestRecordGegendruck(t *testing.T) {
	ctx := context.Background()
	up := &engpassUpstream{}
	s, repo := neu(t, up)
	id, _, err := s.OpenSession(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	gruppe := grosseGruppe()
	clientFertig := make(chan error, 1)
	go func() {
		for _, m := range gruppe {
			if err := s.ClientMessage(ctx, id, m); err != nil {
				clientFertig <- err
				return
			}
		}
		clientFertig <- nil
	}()
	var empfangen int
	warte(t, "Empfang der großen Gruppe", func() {
		for {
			rs, err := s.AwaitServer(ctx, id)
			if err != nil {
				t.Error(err)
				return
			}
			empfangen += len(rs)
			s.Delivered(ctx, id)
			if rs[len(rs)-1].Type == model.ResponseReadyForQuery {
				return
			}
		}
	})
	if err := <-clientFertig; err != nil {
		t.Fatal(err)
	}
	if want := 1 + 20*(1+zeilenJeExecute+1) + 1; empfangen != want {
		t.Fatalf("%d Server-Nachrichten, erwartet %d", empfangen, want)
	}
	schliessen(t, s, id, model.EndClosed)
	if got := repo.last(t).Sessions[0].Interactions; len(got) != 1 || len(got[0].Groups[0].Client) != len(gruppe) {
		t.Fatalf("Interaktionen: %#v", got)
	}
}

// CloseSession beendet jeden wartenden Aufruf der Session: ein ClientMessage,
// das an einem nicht lesenden Server sendet, ein AwaitServer, das auf den Server
// wartet, und eine Query, die auf eine laufende Extended-Interaktion wartet;
// jeder liefert ErrSessionEnded.
func TestRecordCloseBeendetWartende(t *testing.T) {
	ctx := context.Background()
	up := &engpassUpstream{}
	s, _ := neu(t, up)
	id, _, err := s.OpenSession(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	gesendet := make(chan error, 1)
	go func() {
		var err error
		for _, m := range grosseGruppe() {
			if err = s.ClientMessage(ctx, id, m); err != nil {
				break
			}
		}
		gesendet <- err
	}()
	// Niemand liest die Ausgabe: der Server blockiert beim Schreiben und liest
	// nicht weiter, Send blockiert.
	time.Sleep(100 * time.Millisecond)
	select {
	case err := <-gesendet:
		t.Fatalf("Send kehrte zurück, obwohl der Server nicht liest: %v", err)
	default:
	}
	warte(t, "CloseSession", func() { _ = s.CloseSession(ctx, id, model.EndClosed) })
	select {
	case err := <-gesendet:
		if !errors.Is(err, model.ErrSessionEnded) {
			t.Fatalf("ClientMessage nach CloseSession: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ClientMessage endet nicht nach CloseSession")
	}

	// AwaitServer, das am Upstream wartet, und Query, die auf eine laufende
	// Interaktion wartet.
	up2 := &fakeUpstream{}
	s2, _ := neu(t, up2)
	id2 := session(t, s2)
	client(t, s2, id2, cSync())
	await := make(chan error, 1)
	go func() {
		_, err := s2.AwaitServer(ctx, id2)
		await <- err
	}()
	query := make(chan error, 1)
	go func() {
		_, err := s2.Query(ctx, id2, "SELECT 1")
		query <- err
	}()
	time.Sleep(50 * time.Millisecond)
	warte(t, "CloseSession", func() { _ = s2.CloseSession(ctx, id2, model.EndClosed) })
	for name, ch := range map[string]chan error{"AwaitServer": await, "Query": query} {
		select {
		case err := <-ch:
			if !errors.Is(err, model.ErrSessionEnded) {
				t.Fatalf("%s nach CloseSession: %v", name, err)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s endet nicht nach CloseSession", name)
		}
	}
}

// CloseSession beendet auch einen Aufruf, dessen Upstream-Aufruf nach dem Ende
// erfolgreich zurückkommt: Query und AwaitServer liefern dann ErrSessionEnded
// statt Antworten, und nichts wird übernommen.
func TestRecordEndeNachErfolgreichemUpstream(t *testing.T) {
	t.Run("Query", func(t *testing.T) {
		ctx := t.Context()
		up := &fakeUpstream{}
		s, repo := neu(t, up)
		id := session(t, s, "SELECT 1")
		up.letzte.queryLaeuft = make(chan struct{})
		up.letzte.queryHalt = make(chan struct{})
		fertig := make(chan error, 1)
		go func() {
			_, err := s.Query(ctx, id, "SELECT 2")
			fertig <- err
		}()
		<-up.letzte.queryLaeuft
		if err := s.CloseSession(ctx, id, model.EndClosed); codeOf(err) != model.CodeConnectionLost {
			t.Fatalf("CloseSession: %v", err)
		}
		close(up.letzte.queryHalt)
		if err := <-fertig; !errors.Is(err, model.ErrSessionEnded) {
			t.Fatalf("Query: %v", err)
		}
		if got := repo.last(t).Sessions[0].Interactions; len(got) != 1 {
			t.Fatalf("Interaktionen: %#v", got)
		}
	})
	t.Run("AwaitServer", func(t *testing.T) {
		ctx := t.Context()
		up := &fakeUpstream{}
		s, _ := neu(t, up)
		id := session(t, s)
		client(t, s, id, cSync())
		up.letzte.receiveHalt = make(chan struct{})
		up.letzte.empfang <- []model.Response{sRFQ()}
		fertig := make(chan error, 1)
		go func() {
			_, err := s.AwaitServer(ctx, id)
			fertig <- err
		}()
		time.Sleep(50 * time.Millisecond)
		if err := s.CloseSession(ctx, id, model.EndClosed); codeOf(err) != model.CodeConnectionLost {
			t.Fatalf("CloseSession: %v", err)
		}
		close(up.letzte.receiveHalt)
		if err := <-fertig; !errors.Is(err, model.ErrSessionEnded) {
			t.Fatalf("AwaitServer: %v", err)
		}
	})
}
