package pgwire

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// fakeRecorder bildet den Record-Use-Case nach. ereignisse protokolliert die
// Aufrufe in ihrer Reihenfolge: „q:<sql>“, „c:<typ>“, „zugestellt“,
// „shutdown“ und „ende:<grund>“. AwaitServer liefert die Folgen aus server und
// wartet sonst bis CloseSession (dann ErrSessionEnded); ein geschlossener Kanal
// server ist PGR-E4003. halt[typ] lässt ClientMessage für diesen Typ warten, bis
// der Kanal geschlossen ist oder CloseSession läuft.
type fakeRecorder struct {
	mu           sync.Mutex
	opened       int
	ends         []model.SessionEnd
	err          error
	aufbauFehler bool

	ereignisse      []string
	client          []model.ClientMessage
	server          chan []model.Response
	halt            map[model.ClientMessageType]chan struct{}
	queryHalt       chan struct{}
	clientErr       error
	closeErr        error
	shutdownEndet   bool
	zugestelltEndet bool
	// shutdownAntworten beantwortet die ersten Shutdown-Aufrufe der Reihe nach,
	// danach gilt shutdownEndet; shutdownHalt[n] lässt den n-ten Aufruf (ab 1)
	// warten, bis der Kanal geschlossen ist.
	shutdownAntworten []bool
	shutdownHalt      map[int]chan struct{}
	shutdownAufrufe   int
	// neueNachShutdown: nach dem ersten Shutdown liefert ClientMessage für
	// parse ErrShutdown, wie der Use Case für eine neue Interaktion.
	neueNachShutdown bool

	zuEin sync.Once
	zuCh  chan struct{}
}

func (f *fakeRecorder) zu() chan struct{} {
	f.zuEin.Do(func() { f.zuCh = make(chan struct{}) })
	return f.zuCh
}

func (f *fakeRecorder) protokolliere(e string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ereignisse = append(f.ereignisse, e)
}

func (f *fakeRecorder) protokoll() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return strings.Join(f.ereignisse, " ")
}

func (f *fakeRecorder) OpenSession(context.Context, map[string]string) (model.SessionID, []model.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return 0, nil, f.err
	}
	if f.aufbauFehler {
		return 0, []model.Response{{Type: model.ResponseErrorResponse, Fields: map[string]string{"S": "FATAL", "C": "3D000", "M": "database does not exist"}}}, nil
	}
	f.opened++
	return 1, []model.Response{
		{Type: model.ResponseParameterStatus, Name: "server_version", Value: "17.0"},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}, nil
}

func (f *fakeRecorder) Query(_ context.Context, _ model.SessionID, sql string) ([]model.Response, error) {
	f.protokolliere("q:" + sql)
	if f.queryHalt != nil {
		<-f.queryHalt
	}
	return []model.Response{
		{Type: model.ResponseRowDescription, Columns: []model.Column{{Name: "a", TypeOID: 25, TypeSize: -1, TypeModifier: -1}, {Name: "b", TypeOID: 25, TypeSize: -1, TypeModifier: -1}}},
		{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte("x")}, {Null: true}}},
		{Type: model.ResponseCommandComplete, Tag: "SELECT 1"},
		{Type: model.ResponseReadyForQuery, TxStatus: "T"},
	}, nil
}

func (f *fakeRecorder) ClientMessage(_ context.Context, _ model.SessionID, m model.ClientMessage) error {
	f.mu.Lock()
	f.ereignisse = append(f.ereignisse, "c:"+string(m.Type))
	f.client = append(f.client, m)
	halt := f.halt[m.Type]
	err := f.clientErr
	if f.neueNachShutdown && f.shutdownAufrufe > 0 && m.Type == model.ClientParse {
		err = model.ErrShutdown
	}
	f.mu.Unlock()
	if halt != nil {
		select {
		case <-halt:
		case <-f.zu():
			return model.ErrSessionEnded
		}
	}
	return err
}

func (f *fakeRecorder) AwaitServer(context.Context, model.SessionID) ([]model.Response, error) {
	select {
	case out, ok := <-f.server:
		if !ok {
			return nil, model.Errorf(model.CodeConnectionLost, nil, "Upstream weg")
		}
		return out, nil
	case <-f.zu():
		return nil, model.ErrSessionEnded
	}
}

func (f *fakeRecorder) Delivered(context.Context, model.SessionID) bool {
	f.protokolliere("zugestellt")
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.zugestelltEndet
}

func (f *fakeRecorder) Shutdown(context.Context, model.SessionID) bool {
	f.protokolliere("shutdown")
	f.mu.Lock()
	f.shutdownAufrufe++
	n := f.shutdownAufrufe
	halt := f.shutdownHalt[n]
	f.mu.Unlock()
	if halt != nil {
		<-halt
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if n <= len(f.shutdownAntworten) {
		return f.shutdownAntworten[n-1]
	}
	return f.shutdownEndet
}

func (f *fakeRecorder) CloseSession(_ context.Context, _ model.SessionID, end model.SessionEnd) error {
	f.mu.Lock()
	f.ends = append(f.ends, end)
	f.ereignisse = append(f.ereignisse, fmt.Sprintf("ende:%d", end))
	err := f.closeErr
	f.mu.Unlock()
	f.zuEin.Do(func() { f.zuCh = make(chan struct{}) })
	select {
	case <-f.zuCh:
	default:
		close(f.zuCh)
	}
	return err
}

func (f *fakeRecorder) lastEnd(t *testing.T) model.SessionEnd {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		f.mu.Lock()
		if len(f.ends) > 0 {
			e := f.ends[len(f.ends)-1]
			f.mu.Unlock()
			return e
		}
		f.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("CloseSession nicht aufgerufen")
	return 0
}

// verbinde startet den Server-Handler auf einer Seite eines net.Pipe.
func verbinde(t *testing.T, rec *fakeRecorder) (net.Conn, *Server) {
	t.Helper()
	return verbindeMitLog(t, rec, io.Discard)
}

func verbindeMitLog(t *testing.T, rec *fakeRecorder, log io.Writer) (net.Conn, *Server) {
	t.Helper()
	client, serverSeite := net.Pipe()
	s := NewRecordServer(rec, slog.New(slog.NewTextHandler(log, nil)))
	go s.handle(context.Background(), serverSeite)
	t.Cleanup(func() { client.Close() })
	_ = client.SetDeadline(time.Now().Add(5 * time.Second))
	return client, s
}

func startup(t *testing.T, fe *pgproto3.Frontend) {
	t.Helper()
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
	if err := fe.Flush(); err != nil {
		t.Fatal(err)
	}
	for {
		msg, err := fe.Receive()
		if err != nil {
			t.Fatalf("Startup: %v", err)
		}
		if _, ok := msg.(*pgproto3.ReadyForQuery); ok {
			return
		}
	}
}

func fehlerantwort(t *testing.T, fe *pgproto3.Frontend) *pgproto3.ErrorResponse {
	t.Helper()
	msg, err := fe.Receive()
	if err != nil {
		t.Fatalf("Antwort lesen: %v", err)
	}
	e, ok := msg.(*pgproto3.ErrorResponse)
	if !ok {
		t.Fatalf("erwartet ErrorResponse, erhalten %T", msg)
	}
	return e
}

// Abdeckung: LH-FA-05/Boundary — auf SSLRequest und GSSENCRequest antwortet der
// Recorder mit „N“ und erwartet danach einen unverschlüsselten Aufbau
// (LH-FA-05.c).
func TestSSLUndGSSMitN(t *testing.T) {
	for _, code := range []uint32{codeSSLRequest, codeGSSEncRequest} {
		rec := &fakeRecorder{}
		client, _ := verbinde(t, rec)
		anfrage := make([]byte, 8)
		binary.BigEndian.PutUint32(anfrage[0:4], 8)
		binary.BigEndian.PutUint32(anfrage[4:8], code)
		if _, err := client.Write(anfrage); err != nil {
			t.Fatal(err)
		}
		antwort := make([]byte, 1)
		if _, err := io.ReadFull(client, antwort); err != nil || antwort[0] != 'N' {
			t.Fatalf("Code %d: Antwort %q, Fehler %v", code, antwort, err)
		}
		startup(t, pgproto3.NewFrontend(client, client))
	}
}

// Abdeckung: LH-FA-02/Happy — eine Anfrage geht an den Use Case, die Antworten
// gehen in Reihenfolge und Inhalt unverändert an den Client und werden danach
// als zugestellt gemeldet; Terminate meldet das Ende EndTerminate.
func TestQueryUndTerminate(t *testing.T) {
	rec := &fakeRecorder{}
	client, _ := verbinde(t, rec)
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	fe.Send(&pgproto3.Query{String: "SELECT 'x', NULL"})
	if err := fe.Flush(); err != nil {
		t.Fatal(err)
	}
	var typen []string
	for {
		msg, err := fe.Receive()
		if err != nil {
			t.Fatal(err)
		}
		switch m := msg.(type) {
		case *pgproto3.RowDescription:
			typen = append(typen, "row_description:"+string(m.Fields[0].Name)+","+string(m.Fields[1].Name))
		case *pgproto3.DataRow:
			if string(m.Values[0]) != "x" || m.Values[1] != nil {
				t.Fatalf("Werte: %q", m.Values)
			}
			typen = append(typen, "data_row")
		case *pgproto3.CommandComplete:
			typen = append(typen, "command_complete:"+string(m.CommandTag))
		case *pgproto3.ReadyForQuery:
			typen = append(typen, "ready_for_query:"+string(m.TxStatus))
		}
		if _, ok := msg.(*pgproto3.ReadyForQuery); ok {
			break
		}
	}
	want := "row_description:a,b data_row command_complete:SELECT 1 ready_for_query:T"
	if strings.Join(typen, " ") != want {
		t.Fatalf("Antworten: %v", typen)
	}
	fe.Send(&pgproto3.Terminate{})
	_ = fe.Flush()
	if end := rec.lastEnd(t); end != model.EndTerminate {
		t.Fatalf("Session-Ende: %v", end)
	}
	if got := rec.protokoll(); got != "q:SELECT 'x', NULL zugestellt ende:1" {
		t.Fatalf("Protokoll %q", got)
	}
}

// Abdeckung: LH-FA-05/Negative, LH-FA-06/Negative — eine nicht unterstützte Client-Nachricht
// beendet die Verbindung mit einer ErrorResponse (PGR-E6001, SQLSTATE 0A000),
// und die Session wird verworfen.
func TestNichtUnterstuetzteNachricht(t *testing.T) {
	rec := &fakeRecorder{}
	client, s := verbinde(t, rec)
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	fe.Send(&pgproto3.FunctionCall{Function: 1})
	_ = fe.Flush()
	e := fehlerantwort(t, fe)
	if e.Code != "0A000" || !strings.Contains(e.Message, model.CodeUnsupported) {
		t.Fatalf("ErrorResponse: %+v", e)
	}
	if end := rec.lastEnd(t); end != model.EndUnsupported {
		t.Fatalf("Session-Ende: %v", end)
	}
	if s.FirstErrorCode() != model.CodeUnsupported {
		t.Fatalf("erster Fehler: %q", s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-05/Negative — eine andere Protokollversion als 3.0 erhält eine
// ErrorResponse mit PGR-E6002 (LH-FA-05.e).
func TestAndereProtokollversion(t *testing.T) {
	for _, version := range []uint32{2 << 16, 3<<16 | 1, 3<<16 | 2, 4 << 16, 0} {
		rec := &fakeRecorder{}
		client, _ := verbinde(t, rec)
		msg := make([]byte, 9)
		binary.BigEndian.PutUint32(msg[0:4], 9)
		binary.BigEndian.PutUint32(msg[4:8], version)
		if _, err := client.Write(msg); err != nil {
			t.Fatal(err)
		}
		e := fehlerantwort(t, pgproto3.NewFrontend(client, client))
		if !strings.Contains(e.Message, model.CodeProtocolVersion) || e.Code != "0A000" {
			t.Fatalf("Version %x: %+v", version, e)
		}
		if rec.opened != 0 {
			t.Fatalf("Version %x: Session geöffnet", version)
		}
	}
}

// Abdeckung: LH-FA-05/Boundary — ein CancelRequest wird nicht weitergeleitet; die
// Verbindung endet ohne Session, und die Warnung trägt PGR-W3001.
func TestCancelRequest(t *testing.T) {
	rec := &fakeRecorder{}
	var log syncBuffer
	client, _ := verbindeMitLog(t, rec, &log)
	msg := make([]byte, 16)
	binary.BigEndian.PutUint32(msg[0:4], 16)
	binary.BigEndian.PutUint32(msg[4:8], codeCancelRequest)
	if _, err := client.Write(msg); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Read(make([]byte, 1)); err == nil {
		t.Fatal("Verbindung nicht beendet")
	}
	if rec.opened != 0 {
		t.Fatal("Session geöffnet")
	}
	time.Sleep(50 * time.Millisecond)
	if !strings.Contains(log.String(), model.CodeCancelRequest) {
		t.Fatalf("Warnung ohne %s: %s", model.CodeCancelRequest, log.String())
	}
}

type syncBuffer struct {
	mu sync.Mutex
	b  strings.Builder
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

// Eine unbekannte Sonderanfrage ist PGR-E6001.
func TestUnbekannteSonderanfrage(t *testing.T) {
	sonder := make([]byte, 8)
	binary.BigEndian.PutUint32(sonder[0:4], 8)
	binary.BigEndian.PutUint32(sonder[4:8], majorSpezial<<16|9999)
	rec := &fakeRecorder{}
	client, _ := verbinde(t, rec)
	if _, err := client.Write(sonder); err != nil {
		t.Fatal(err)
	}
	e := fehlerantwort(t, pgproto3.NewFrontend(client, client))
	if !strings.Contains(e.Message, model.CodeUnsupported) {
		t.Fatalf("ErrorResponse: %+v", e)
	}
}

// Abdeckung: LH-FA-05/Boundary — eine erste Nachricht, die keine
// PGWire-Startnachricht ist (HTTP), schließt die Verbindung ohne Antwort; der
// Lauf zählt keinen Fehler, die Warnung trägt PGR-W3003.
func TestFremdeErsteNachricht(t *testing.T) {
	rec := &fakeRecorder{}
	var log syncBuffer
	client, s := verbindeMitLog(t, rec, &log)
	if _, err := client.Write([]byte("GET / HTTP/1.1\r\nHost: x\r\n\r\n")); err != nil {
		t.Fatal(err)
	}
	if n, err := client.Read(make([]byte, 1)); err == nil {
		t.Fatalf("Antwort erhalten (%d Bytes) statt geschlossener Verbindung", n)
	}
	time.Sleep(50 * time.Millisecond)
	if s.FirstErrorCode() != "" || rec.opened != 0 {
		t.Fatalf("Fehler gemerkt %q, Sessions %d", s.FirstErrorCode(), rec.opened)
	}
	if !strings.Contains(log.String(), model.CodeForeignProtocol) {
		t.Fatalf("Warnung ohne %s: %s", model.CodeForeignProtocol, log.String())
	}
}

// Eine Verbindung ohne erste Nachricht endet still: kein Fehler, keine Warnung.
func TestVerbindungOhneNachricht(t *testing.T) {
	rec := &fakeRecorder{}
	var log syncBuffer
	client, s := verbindeMitLog(t, rec, &log)
	client.Close()
	time.Sleep(50 * time.Millisecond)
	if s.FirstErrorCode() != "" || strings.Contains(log.String(), "level=WARN") {
		t.Fatalf("Fehler %q, Log %s", s.FirstErrorCode(), log.String())
	}
}

// Endet der Aufbau beim Server mit einer Fehlerantwort, erhält der Client genau
// diese, ohne vorheriges AuthenticationOk.
func TestFehlerantwortImAufbau(t *testing.T) {
	rec := &fakeRecorder{aufbauFehler: true}
	client, _ := verbinde(t, rec)
	fe := pgproto3.NewFrontend(client, client)
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
	_ = fe.Flush()
	e := fehlerantwort(t, fe)
	if e.Code != "3D000" {
		t.Fatalf("ErrorResponse: %+v", e)
	}
}

// Trennt der Client ohne Terminate, meldet der Adapter EndClosed; einen
// Fehler merkt er nur, wenn der Use Case einen liefert (LH-FA-02.b).
func TestEndeOhneTerminate(t *testing.T) {
	rec := &fakeRecorder{}
	client, s := verbinde(t, rec)
	startup(t, pgproto3.NewFrontend(client, client))
	client.Close()
	if end := rec.lastEnd(t); end != model.EndClosed {
		t.Fatalf("Session-Ende: %v", end)
	}
	if s.FirstErrorCode() != "" {
		t.Fatalf("Fehler gemerkt: %q", s.FirstErrorCode())
	}
}

// Erreicht die Antwort einer Anfrage den Client nicht, meldet der Adapter
// EndWriteFailed und merkt den Fehler, den der Use Case daraus ableitet
// (LH-FA-02.b).
func TestAntwortNichtZugestellt(t *testing.T) {
	rec := &fakeRecorder{closeErr: model.Errorf(model.CodeConnectionLost, nil, "Antwort nicht zugestellt")}
	client, s := verbinde(t, rec)
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	fe.Send(&pgproto3.Query{String: "SELECT 1"})
	_ = fe.Flush()
	client.Close()
	if end := rec.lastEnd(t); end != model.EndWriteFailed {
		t.Fatalf("Session-Ende: %v", end)
	}
	if s.FirstErrorCode() != model.CodeConnectionLost {
		t.Fatalf("erster Fehler: %q", s.FirstErrorCode())
	}
}

// Eine nicht lesbare Client-Nachricht (unbekannter Typ) erhält eine
// ErrorResponse mit PGR-E6001, und die Session wird verworfen.
func TestUnlesbareNachricht(t *testing.T) {
	rec := &fakeRecorder{}
	client, _ := verbinde(t, rec)
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	if _, err := client.Write([]byte{'z', 0, 0, 0, 4}); err != nil {
		t.Fatal(err)
	}
	e := fehlerantwort(t, fe)
	if !strings.Contains(e.Message, model.CodeUnsupported) {
		t.Fatalf("ErrorResponse: %+v", e)
	}
	if end := rec.lastEnd(t); end != model.EndUnsupported {
		t.Fatalf("Session-Ende: %v", end)
	}
}

// Der Lauf merkt sich den ersten Verbindungsfehler, nicht den letzten.
func TestErsterFehlerZaehlt(t *testing.T) {
	s := &Server{log: slog.New(slog.NewTextHandler(io.Discard, nil))}
	s.note(model.Errorf(model.CodeUpstream, nil, "erster"))
	s.note(model.Errorf(model.CodeUnsupported, nil, "zweiter"))
	if s.FirstErrorCode() != model.CodeUpstream {
		t.Fatalf("erster Fehler: %q", s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-02/Negative — ist der Upstream nicht erreichbar, erhält der
// Client eine ErrorResponse mit PGR-E4002 und SQLSTATE 08006; der Lauf merkt sich
// den Fehler für den Exit-Code (LH-FA-13.b).
func TestUpstreamNichtErreichbar(t *testing.T) {
	rec := &fakeRecorder{err: model.Errorf(model.CodeUpstream, nil, "weg")}
	client, s := verbinde(t, rec)
	fe := pgproto3.NewFrontend(client, client)
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
	_ = fe.Flush()
	e := fehlerantwort(t, fe)
	if e.Code != "08006" || !strings.Contains(e.Message, model.CodeUpstream) {
		t.Fatalf("ErrorResponse: %+v", e)
	}
	if s.FirstErrorCode() != model.CodeUpstream {
		t.Fatalf("erster Fehler: %q", s.FirstErrorCode())
	}
}

// Eine Antwort ohne Transaktionsstatus oder mit unbekanntem Typ wird nicht
// erfunden, sondern ist ein Fehler.
func TestToMessageErfindetNichts(t *testing.T) {
	if _, err := toMessage(model.Response{Type: model.ResponseReadyForQuery}); err == nil {
		t.Fatal("ReadyForQuery ohne Status angenommen")
	}
	if _, err := toMessage(model.Response{Type: "bogus"}); err == nil {
		t.Fatal("unbekannter Typ angenommen")
	}
}

type fakeReplayer struct {
	mu       sync.Mutex
	closed   int
	warnung  *model.Warning
	extended []model.ClientMessage
}

func (f *fakeReplayer) OpenConnection(context.Context) (model.SessionID, []model.Response) {
	return 7, []model.Response{
		{Type: model.ResponseParameterStatus, Name: "server_version", Value: "17.0"},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}
}

func (f *fakeReplayer) Query(_ context.Context, _ model.SessionID, sql string) ([]model.Response, error) {
	if sql != "SELECT 1" {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "erwartet %q, empfangen %q", "SELECT 1", sql)
	}
	return []model.Response{{Type: model.ResponseCommandComplete, Tag: "SELECT 1"}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}, nil
}

// ClientMessage merkt sich jede Extended-Nachricht; ein Sync erhält eine
// vollständige Antwort, ein Execute auf das Portal "abweichend" PGR-E5001, jede
// andere Nachricht keine Antwort.
func (f *fakeReplayer) ClientMessage(_ context.Context, _ model.SessionID, m model.ClientMessage) ([]model.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.extended = append(f.extended, m)
	switch {
	case m.Type == model.ClientExecute && m.Portal == "abweichend":
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "erwartet execute, empfangen execute, abweichend in portal")
	case m.Type == model.ClientSync:
		return []model.Response{
			{Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete},
			{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte("a")}}},
			{Type: model.ResponseCommandComplete, Tag: "SELECT 1"}, {Type: model.ResponseReadyForQuery, TxStatus: "I"},
		}, nil
	}
	return nil, nil
}

func (f *fakeReplayer) nachrichten() []model.ClientMessage {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]model.ClientMessage(nil), f.extended...)
}

func (f *fakeReplayer) CloseConnection(context.Context, model.SessionID) *model.Warning {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closed++
	return f.warnung
}

// Im Replay-Modus beantwortet der Adapter Startup und Anfragen aus dem
// Replay-Use-Case; eine Abweichung erhält eine ErrorResponse mit PGR-E5001 und
// zählt als Fehler der Klasse 5; eine unverbrauchte Session ist nur eine Warnung.
func TestReplayModus(t *testing.T) {
	rep := &fakeReplayer{warnung: model.Warnf(model.CodeUnconsumed, "nicht verbraucht")}
	client, serverSeite := net.Pipe()
	var log syncBuffer
	s := NewReplayServer(rep, slog.New(slog.NewTextHandler(&log, nil)))
	go s.handle(context.Background(), serverSeite)
	defer client.Close()
	_ = client.SetDeadline(time.Now().Add(5 * time.Second))

	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	fe.Send(&pgproto3.Query{String: "SELECT 2"})
	_ = fe.Flush()
	e := fehlerantwort(t, fe)
	if !strings.Contains(e.Message, model.CodeReplayMismatch) {
		t.Fatalf("ErrorResponse: %+v", e)
	}
	time.Sleep(50 * time.Millisecond)
	if s.FirstErrorCode() != model.CodeReplayMismatch {
		t.Fatalf("erster Fehler: %q", s.FirstErrorCode())
	}
	if !strings.Contains(log.String(), model.CodeUnconsumed) {
		t.Fatalf("Warnung %s fehlt: %s", model.CodeUnconsumed, log.String())
	}
}
