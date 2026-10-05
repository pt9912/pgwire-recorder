package pgwire

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// verbindeExtended startet den Handler im Record-Modus mit einem Recorder, der
// Server-Nachrichten aus seinem Kanal liefert.
func verbindeExtended(t *testing.T, ctx context.Context, rec *fakeRecorder) (net.Conn, *pgproto3.Frontend, *Server) {
	t.Helper()
	if rec.server == nil {
		rec.server = make(chan []model.Response, 4)
	}
	client, serverSeite := net.Pipe()
	s := NewRecordServer(rec, slog.New(slog.NewTextHandler(io.Discard, nil)))
	go s.handle(ctx, serverSeite)
	t.Cleanup(func() { client.Close() })
	_ = client.SetDeadline(time.Now().Add(5 * time.Second))
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	return client, fe, s
}

// sendeMu hält den Schreibpuffer des Frontends von Send bis zum Ende des
// nebenher laufenden Flush.
var sendeMu sync.Mutex

// sende schreibt die Nachrichten nebenher; net.Pipe blockiert, bis der Server
// liest. Ein weiterer Aufruf wartet, bis der Server die vorigen gelesen hat.
func sende(fe *pgproto3.Frontend, msgs ...pgproto3.FrontendMessage) {
	sendeMu.Lock()
	for _, m := range msgs {
		fe.Send(m)
	}
	go func() {
		defer sendeMu.Unlock()
		_ = fe.Flush()
	}()
}

// warteAuf wartet, bis das Protokoll des Recorders want enthält.
func warteAuf(t *testing.T, rec *fakeRecorder, want string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !strings.Contains(rec.protokoll(), want) {
		if time.Now().After(deadline) {
			t.Fatalf("Protokoll %q enthält %q nicht", rec.protokoll(), want)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// empfange liest n Nachrichten und liefert ihre Go-Typen.
func empfange(t *testing.T, fe *pgproto3.Frontend, n int) []string {
	t.Helper()
	var typen []string
	for i := 0; i < n; i++ {
		msg, err := fe.Receive()
		if err != nil {
			t.Fatalf("Nachricht %d: %v", i+1, err)
		}
		typen = append(typen, strings.TrimPrefix(reflect.TypeOf(msg).String(), "*pgproto3."))
	}
	return typen
}

func endeName(e model.SessionEnd) string { return fmt.Sprintf("ende:%d", e) }

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: die Client-Nachrichten gehen in
// Ankunftsreihenfolge an den Use Case, jede mit genau den Feldern ihres Typs,
// auch die ohne Warten nach dem Sync gesendeten; die Server-Nachrichten gehen
// in Reihenfolge an den Client und werden danach als zugestellt gemeldet.
func TestExtendedSyncGruppe(t *testing.T) {
	rec := &fakeRecorder{}
	_, fe, s := verbindeExtended(t, context.Background(), rec)
	sende(fe,
		&pgproto3.Parse{Name: "s1", Query: "SELECT $1, $2, $3", ParameterOIDs: []uint32{23, 25, 25}},
		&pgproto3.Bind{DestinationPortal: "p", PreparedStatement: "s1", ParameterFormatCodes: []int16{0}, Parameters: [][]byte{[]byte("1"), nil, {}}, ResultFormatCodes: []int16{1}},
		&pgproto3.Describe{ObjectType: 'P', Name: "p"},
		&pgproto3.Execute{Portal: "p", MaxRows: 5},
		&pgproto3.Sync{},
		&pgproto3.Close{ObjectType: 'S', Name: "s1"},
		&pgproto3.Sync{},
	)
	warteAuf(t, rec, "c:parse c:bind c:describe c:execute c:sync c:close c:sync")
	want := []model.ClientMessage{
		{Type: model.ClientParse, Statement: "s1", SQL: "SELECT $1, $2, $3", ParamTypes: []uint32{23, 25, 25}},
		{Type: model.ClientBind, Portal: "p", Statement: "s1", ParamFormats: []int16{0}, Params: []model.Value{{Bytes: []byte("1")}, {Null: true}, {Bytes: []byte{}}}, ResultFormats: []int16{1}},
		{Type: model.ClientDescribe, Target: model.TargetPortal, Name: "p"},
		{Type: model.ClientExecute, Portal: "p", MaxRows: 5},
		{Type: model.ClientSync},
		{Type: model.ClientClose, Target: model.TargetStatement, Name: "s1"},
		{Type: model.ClientSync},
	}
	rec.mu.Lock()
	gotClient := rec.client
	rec.mu.Unlock()
	if !reflect.DeepEqual(gotClient, want) {
		t.Fatalf("Client-Nachrichten:\n%#v\nerwartet\n%#v", gotClient, want)
	}

	rec.server <- []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete}}
	rec.server <- []model.Response{
		{Type: model.ResponseRowDescription, Columns: []model.Column{{Name: "a", TypeOID: 23, TypeSize: 4, TypeModifier: -1, Format: 1}}},
		{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte{0, 0, 0, 1}}}},
		{Type: model.ResponsePortalSuspended},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}
	if got := strings.Join(empfange(t, fe, 6), " "); got != "ParseComplete BindComplete RowDescription DataRow PortalSuspended ReadyForQuery" {
		t.Fatalf("an den Client: %s", got)
	}
	rec.server <- []model.Response{{Type: model.ResponseCloseComplete}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}
	if got := strings.Join(empfange(t, fe, 2), " "); got != "CloseComplete ReadyForQuery" {
		t.Fatalf("an den Client: %s", got)
	}
	warteAuf(t, rec, "zugestellt zugestellt zugestellt")
	sende(fe, &pgproto3.Terminate{})
	if end := rec.lastEnd(t); end != model.EndTerminate || s.FirstErrorCode() != "" {
		t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: die beiden Richtungen blockieren
// unabhängig: solange der Use Case beim Senden einer Gruppe wartet (der Server
// liest nicht), gehen Server-Nachrichten an den Client, und weitere
// Client-Nachrichten warten nicht auf das Lesen der Antworten.
func TestExtendedZweiRichtungen(t *testing.T) {
	halt := make(chan struct{})
	rec := &fakeRecorder{halt: map[model.ClientMessageType]chan struct{}{model.ClientSync: halt}}
	_, fe, _ := verbindeExtended(t, context.Background(), rec)
	sende(fe, &pgproto3.Parse{Query: "SELECT 1"}, &pgproto3.Bind{}, &pgproto3.Execute{}, &pgproto3.Sync{})
	warteAuf(t, rec, "c:sync")
	rec.server <- []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete}}
	if got := strings.Join(empfange(t, fe, 2), " "); got != "ParseComplete BindComplete" {
		t.Fatalf("an den Client: %s", got)
	}
	warteAuf(t, rec, "zugestellt")
	close(halt)
	rec.server <- []model.Response{{Type: model.ResponseCommandComplete, Tag: "SELECT 1"}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}
	empfange(t, fe, 2)
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: nach einem Flush gehen die
// Server-Nachrichten an den Client, ohne dass ein Sync folgt.
func TestExtendedFlush(t *testing.T) {
	rec := &fakeRecorder{}
	_, fe, _ := verbindeExtended(t, context.Background(), rec)
	sende(fe, &pgproto3.Parse{Name: "s1", Query: "SELECT $1"}, &pgproto3.Describe{ObjectType: 'S', Name: "s1"}, &pgproto3.Flush{})
	warteAuf(t, rec, "c:flush")
	rec.server <- []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseParameterDescription, ParamTypes: []uint32{23}}, {Type: model.ResponseNoData}}
	var oids []uint32
	for i := 0; i < 3; i++ {
		msg, err := fe.Receive()
		if err != nil {
			t.Fatal(err)
		}
		if pd, ok := msg.(*pgproto3.ParameterDescription); ok {
			oids = append(oids, pd.ParameterOIDs...)
		}
		if i == 2 {
			if _, ok := msg.(*pgproto3.NoData); !ok {
				t.Fatalf("dritte Nachricht %T", msg)
			}
		}
	}
	if !reflect.DeepEqual(oids, []uint32{23}) {
		t.Fatalf("ParameterDescription: %v", oids)
	}
}

// Abdeckung: LH-FA-13/Negative — der Adapter meldet, was an der Verbindung
// geschah, und merkt den Fehler, den der Use Case daraus ableitet:
// Verbindungsende EndClosed, Terminate EndTerminate, ein Schreibfehler zum
// Client EndWriteFailed, ein Fehler aus AwaitServer geht als FATAL an den Client
// und endet mit EndFailed; danach schließt der Adapter die Verbindung, auch wenn
// die Client-Richtung gerade liest. Jede Session endet genau einmal.
func TestExtendedEreignisse(t *testing.T) {
	verloren := model.Errorf(model.CodeConnectionLost, nil, "vom Use Case")
	cases := []struct {
		name     string
		ablauf   func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder)
		wantEnd  model.SessionEnd
		wantCode string
	}{
		{"Verbindungsende", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			sende(fe, &pgproto3.Parse{Query: "SELECT 1"})
			warteAuf(t, rec, "c:parse")
			client.Close()
		}, model.EndClosed, model.CodeConnectionLost},
		{"Terminate", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			sende(fe, &pgproto3.Parse{Query: "SELECT 1"}, &pgproto3.Flush{})
			warteAuf(t, rec, "c:flush")
			sende(fe, &pgproto3.Terminate{})
		}, model.EndTerminate, model.CodeConnectionLost},
		{"Schreibfehler, während ClientMessage wartet", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			rec.mu.Lock()
			rec.halt = map[model.ClientMessageType]chan struct{}{model.ClientSync: make(chan struct{})}
			rec.mu.Unlock()
			sende(fe, &pgproto3.Sync{})
			warteAuf(t, rec, "c:sync")
			client.Close()
			rec.server <- []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
		}, model.EndWriteFailed, model.CodeConnectionLost},
		{"Fehler aus AwaitServer", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			sende(fe, &pgproto3.Sync{})
			warteAuf(t, rec, "c:sync")
			close(rec.server)
			e := fehlerantwort(t, fe)
			if !strings.Contains(e.Message, model.CodeConnectionLost) {
				t.Fatalf("ErrorResponse: %+v", e)
			}
			// Die Client-Richtung endet mit: der Adapter schließt die Verbindung.
			if _, err := fe.Receive(); !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatalf("Verbindung nach dem Ende nicht geschlossen: %v", err)
			}
		}, model.EndFailed, model.CodeConnectionLost},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := &fakeRecorder{}
			if c.wantEnd != model.EndFailed {
				rec.closeErr = verloren
			}
			client, fe, s := verbindeExtended(t, context.Background(), rec)
			c.ablauf(t, client, fe, rec)
			if end := rec.lastEnd(t); end != c.wantEnd {
				t.Fatalf("Session-Ende %v, erwartet %v", end, c.wantEnd)
			}
			time.Sleep(20 * time.Millisecond)
			if got := s.FirstErrorCode(); got != c.wantCode {
				t.Fatalf("erster Fehler %q, erwartet %q", got, c.wantCode)
			}
			rec.mu.Lock()
			n := len(rec.ends)
			rec.mu.Unlock()
			if n != 1 {
				t.Fatalf("%d Session-Enden gemeldet", n)
			}
		})
	}
}

// Abdeckung: LH-FA-06/Negative — ein Describe mit einer anderen Zielart als 'S'
// oder 'P', ein Fehler des Use Case zu einer Client-Nachricht und eine
// Extended-Nachricht im Replay-Modus beenden die Verbindung mit PGR-E6001
// (SQLSTATE 0A000); im Record-Modus meldet der Adapter EndUnsupported
// beziehungsweise EndFailed.
func TestExtendedNichtUnterstuetzt(t *testing.T) {
	t.Run("Zielart", func(t *testing.T) {
		rec := &fakeRecorder{}
		client, fe, s := verbindeExtended(t, context.Background(), rec)
		go func() { _, _ = client.Write([]byte{'D', 0, 0, 0, 6, 'X', 0}) }()
		e := fehlerantwort(t, fe)
		if e.Code != "0A000" || !strings.Contains(e.Message, model.CodeUnsupported) {
			t.Fatalf("ErrorResponse: %+v", e)
		}
		if end := rec.lastEnd(t); end != model.EndUnsupported || s.FirstErrorCode() != model.CodeUnsupported {
			t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
		}
		if strings.Contains(rec.protokoll(), "c:") {
			t.Fatalf("an den Use Case übergeben: %s", rec.protokoll())
		}
	})
	t.Run("Fehler des Use Case", func(t *testing.T) {
		rec := &fakeRecorder{clientErr: model.Errorf(model.CodeUnsupported, nil, "Form")}
		_, fe, s := verbindeExtended(t, context.Background(), rec)
		sende(fe, &pgproto3.Sync{})
		e := fehlerantwort(t, fe)
		if e.Code != "0A000" || !strings.Contains(e.Message, model.CodeUnsupported) {
			t.Fatalf("ErrorResponse: %+v", e)
		}
		if end := rec.lastEnd(t); end != model.EndFailed || s.FirstErrorCode() != model.CodeUnsupported {
			t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
		}
	})
	t.Run("Replay", func(t *testing.T) {
		rep := &fakeReplayer{}
		client, serverSeite := net.Pipe()
		s := NewReplayServer(rep, slog.New(slog.NewTextHandler(io.Discard, nil)))
		go s.handle(context.Background(), serverSeite)
		defer client.Close()
		_ = client.SetDeadline(time.Now().Add(5 * time.Second))
		fe := pgproto3.NewFrontend(client, client)
		startup(t, fe)
		sende(fe, &pgproto3.Parse{Query: "SELECT 1"})
		e := fehlerantwort(t, fe)
		if e.Code != "0A000" || !strings.Contains(e.Message, model.CodeUnsupported) || !strings.Contains(e.Message, "Replay") {
			t.Fatalf("ErrorResponse: %+v", e)
		}
	})
}

// Abdeckung: LH-FA-13/Boundary — beim Herunterfahren fragt der Adapter den Use
// Case: gibt er das Ende frei, endet die Session mit EndShutdown; sonst liest
// der Adapter weiter, bis Delivered das Ende freigibt. Eine Anfrage, die beim
// Beginn schon gelesen ist, wird noch beantwortet (LH-FA-13.a).
func TestExtendedHerunterfahren(t *testing.T) {
	t.Run("laufende Interaktion", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		rec := &fakeRecorder{}
		_, fe, s := verbindeExtended(t, ctx, rec)
		sende(fe, &pgproto3.Parse{Query: "SELECT 1"})
		warteAuf(t, rec, "c:parse")
		cancel()
		warteAuf(t, rec, "shutdown")
		time.Sleep(50 * time.Millisecond)
		sende(fe, &pgproto3.Bind{}, &pgproto3.Execute{}, &pgproto3.Sync{})
		warteAuf(t, rec, "c:sync")
		if strings.Contains(rec.protokoll(), "ende") {
			t.Fatalf("Session vor dem ReadyForQuery beendet: %s", rec.protokoll())
		}
		rec.mu.Lock()
		rec.zugestelltEndet = true
		rec.shutdownEndet = true
		rec.mu.Unlock()
		rec.server <- []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}
		empfange(t, fe, 3)
		if end := rec.lastEnd(t); end != model.EndShutdown || s.FirstErrorCode() != "" {
			t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
		}
		if got := rec.protokoll(); !strings.HasSuffix(got, "zugestellt shutdown "+endeName(model.EndShutdown)) {
			t.Fatalf("Protokoll %q", got)
		}
	})
	t.Run("ruhend", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		rec := &fakeRecorder{shutdownEndet: true}
		_, _, s := verbindeExtended(t, ctx, rec)
		cancel()
		if end := rec.lastEnd(t); end != model.EndShutdown || s.FirstErrorCode() != "" {
			t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
		}
	})
	t.Run("gelesene Anfrage", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		halt := make(chan struct{})
		rec := &fakeRecorder{queryHalt: halt, shutdownEndet: true, zugestelltEndet: true}
		_, fe, _ := verbindeExtended(t, ctx, rec)
		sende(fe, &pgproto3.Query{String: "SELECT 1"})
		warteAuf(t, rec, "q:SELECT 1")
		cancel()
		time.Sleep(50 * time.Millisecond)
		close(halt)
		if got := strings.Join(empfange(t, fe, 4), " "); got != "RowDescription DataRow CommandComplete ReadyForQuery" {
			t.Fatalf("an den Client: %s", got)
		}
		if end := rec.lastEnd(t); end != model.EndShutdown {
			t.Fatalf("Ende %v", end)
		}
		if got := rec.protokoll(); got != "q:SELECT 1 zugestellt shutdown "+endeName(model.EndShutdown) {
			t.Fatalf("Protokoll %q: Herunterfahren nicht erst nach der gelesenen Anfrage gemeldet", got)
		}
	})
	t.Run("neue Interaktion nach dem Beginn", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		rec := &fakeRecorder{neueNachShutdown: true}
		_, fe, s := verbindeExtended(t, ctx, rec)
		sende(fe, &pgproto3.Sync{})
		warteAuf(t, rec, "c:sync")
		cancel()
		warteAuf(t, rec, "shutdown")
		sende(fe, &pgproto3.Parse{Query: "SELECT 2"}, &pgproto3.Bind{}, &pgproto3.Sync{})
		warteAuf(t, rec, "c:parse")
		rec.mu.Lock()
		rec.shutdownEndet = true
		rec.zugestelltEndet = true
		rec.mu.Unlock()
		rec.server <- []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
		empfange(t, fe, 1)
		if end := rec.lastEnd(t); end != model.EndShutdown || s.FirstErrorCode() != "" {
			t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
		}
		if got := rec.protokoll(); strings.Contains(got, "c:bind") {
			t.Fatalf("nach der abgelehnten Nachricht weitergelesen: %q", got)
		}
	})
}

// Ein Wecken während des Fragens geht nicht verloren: Gibt Shutdown das Ende
// nicht frei, und gibt Delivered es frei, während die Client-Richtung noch
// fragt, endet die Session mit EndShutdown, ohne dass der Client weiter sendet.
func TestHerunterfahrenWeckenNichtVerloren(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	halt := make(chan struct{})
	rec := &fakeRecorder{shutdownAntworten: []bool{false, false}, shutdownHalt: map[int]chan struct{}{2: halt}}
	_, fe, s := verbindeExtended(t, ctx, rec)
	sende(fe, &pgproto3.Sync{})
	warteAuf(t, rec, "c:sync")
	cancel()
	warteAuf(t, rec, "shutdown shutdown")
	rec.mu.Lock()
	rec.zugestelltEndet = true
	rec.shutdownEndet = true
	rec.mu.Unlock()
	rec.server <- []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
	empfange(t, fe, 1)
	warteAuf(t, rec, "zugestellt")
	close(halt)
	if end := rec.lastEnd(t); end != model.EndShutdown || s.FirstErrorCode() != "" {
		t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
	}
}

// weiterlesen setzt die Lesefrist nicht zurück, wenn vorher geweckt wurde;
// ohne Wecken setzt es sie zurück.
func TestWeiterlesenNachWecken(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	r := &richtungen{conn: a, weck: make(chan struct{}, 1)}
	r.wecke()
	if !r.weiterlesen() {
		t.Fatal("Wecken nicht gemeldet")
	}
	fehler := make(chan error, 1)
	go func() {
		_, err := a.Read(make([]byte, 1))
		fehler <- err
	}()
	select {
	case err := <-fehler:
		if !errors.Is(err, os.ErrDeadlineExceeded) {
			t.Fatalf("Lesen: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Lesefrist nach dem Wecken zurückgesetzt")
	}
	if r.weiterlesen() {
		t.Fatal("Wecken ohne Signal gemeldet")
	}
	go func() {
		_, err := a.Read(make([]byte, 1))
		fehler <- err
	}()
	select {
	case err := <-fehler:
		t.Fatalf("Lesefrist nicht zurückgesetzt: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
}

// Abdeckung: LH-FA-13/Boundary — endet eine Session, schließt der Adapter die
// Client-Verbindung; ein Schreiben an einen Client, der nicht liest, endet
// damit, und die Sitzung kehrt zurück.
func TestEndeSchliesstVerbindung(t *testing.T) {
	rec := &fakeRecorder{server: make(chan []model.Response, 1)}
	client, serverSeite := net.Pipe()
	defer client.Close()
	s := NewRecordServer(rec, slog.New(slog.NewTextHandler(io.Discard, nil)))
	zurueck := make(chan struct{})
	go func() {
		s.handle(context.Background(), serverSeite)
		close(zurueck)
	}()
	_ = client.SetDeadline(time.Now().Add(5 * time.Second))
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	sende(fe, &pgproto3.Sync{})
	warteAuf(t, rec, "c:sync")
	// Der Client liest nicht: das Schreiben der Server-Richtung blockiert.
	rec.server <- []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
	time.Sleep(50 * time.Millisecond)
	sende(fe, &pgproto3.Terminate{})
	select {
	case <-zurueck:
	case <-time.After(2 * time.Second):
		t.Fatal("Sitzung kehrt nach Terminate nicht zurück, solange der Client nicht liest")
	}
	if end := rec.lastEnd(t); end != model.EndTerminate {
		t.Fatalf("Ende %v", end)
	}
}

// Die Abbildung belegt je Client-Nachricht genau die Felder ihres Typs
// (SPEC-041), eine leere Liste als nil, und kopiert Werte aus dem Puffer der
// Bibliothek.
func TestToClientMessage(t *testing.T) {
	param := []byte("x")
	cases := []struct {
		in   pgproto3.FrontendMessage
		want model.ClientMessage
	}{
		{&pgproto3.Parse{Name: "s", Query: "q", ParameterOIDs: []uint32{}}, model.ClientMessage{Type: model.ClientParse, Statement: "s", SQL: "q"}},
		{&pgproto3.Bind{DestinationPortal: "p", PreparedStatement: "s", ParameterFormatCodes: []int16{}, Parameters: [][]byte{param}, ResultFormatCodes: []int16{}},
			model.ClientMessage{Type: model.ClientBind, Portal: "p", Statement: "s", Params: []model.Value{{Bytes: []byte("x")}}}},
		{&pgproto3.Describe{ObjectType: 'S', Name: "n"}, model.ClientMessage{Type: model.ClientDescribe, Target: model.TargetStatement, Name: "n"}},
		{&pgproto3.Close{ObjectType: 'P', Name: "n"}, model.ClientMessage{Type: model.ClientClose, Target: model.TargetPortal, Name: "n"}},
		{&pgproto3.Execute{Portal: "p", MaxRows: 3}, model.ClientMessage{Type: model.ClientExecute, Portal: "p", MaxRows: 3}},
		{&pgproto3.Flush{}, model.ClientMessage{Type: model.ClientFlush}},
		{&pgproto3.Sync{}, model.ClientMessage{Type: model.ClientSync}},
	}
	for _, c := range cases {
		got, err := toClientMessage(c.in)
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Fatalf("%T: %#v, %v; erwartet %#v", c.in, got, err, c.want)
		}
	}
	got, _ := toClientMessage(&pgproto3.Bind{Parameters: [][]byte{param}})
	param[0] = 'y'
	if string(got.Params[0].Bytes) != "x" {
		t.Fatal("Parameterwert nicht kopiert")
	}
	if _, err := toClientMessage(&pgproto3.Close{ObjectType: 'X'}); err == nil {
		t.Fatal("Close mit Zielart X angenommen")
	}
}

// Die Server-Nachrichten des Extended Query Protocol werden auf ihre
// PGWire-Nachricht abgebildet.
func TestToMessageExtended(t *testing.T) {
	cases := map[model.ResponseType]pgproto3.BackendMessage{
		model.ResponseParseComplete:        &pgproto3.ParseComplete{},
		model.ResponseBindComplete:         &pgproto3.BindComplete{},
		model.ResponseCloseComplete:        &pgproto3.CloseComplete{},
		model.ResponseParameterDescription: &pgproto3.ParameterDescription{ParameterOIDs: []uint32{23, 25}},
		model.ResponseNoData:               &pgproto3.NoData{},
		model.ResponsePortalSuspended:      &pgproto3.PortalSuspended{},
	}
	for typ, want := range cases {
		got, err := toMessage(model.Response{Type: typ, ParamTypes: []uint32{23, 25}})
		if typ != model.ResponseParameterDescription {
			got, err = toMessage(model.Response{Type: typ})
		}
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: %#v, %v", typ, got, err)
		}
	}
}
