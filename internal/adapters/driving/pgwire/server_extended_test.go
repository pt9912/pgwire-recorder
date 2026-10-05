package pgwire

import (
	"context"
	"io"
	"log/slog"
	"net"
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
func verbindeExtended(t *testing.T, ctx context.Context) (net.Conn, *pgproto3.Frontend, *fakeRecorder, *Server) {
	t.Helper()
	rec := &fakeRecorder{server: make(chan []model.Response, 4)}
	client, serverSeite := net.Pipe()
	s := NewRecordServer(rec, slog.New(slog.NewTextHandler(io.Discard, nil)))
	go s.handle(ctx, serverSeite)
	t.Cleanup(func() { client.Close() })
	_ = client.SetDeadline(time.Now().Add(5 * time.Second))
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	return client, fe, rec, s
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

// warteAuf wartet, bis das Protokoll des Recorders mit want endet.
func warteAuf(t *testing.T, rec *fakeRecorder, want string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !strings.HasSuffix(rec.protokoll(), want) {
		if time.Now().After(deadline) {
			t.Fatalf("Protokoll %q endet nicht mit %q", rec.protokoll(), want)
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

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: die Nachrichten einer Sync-Gruppe gehen in
// Ankunftsreihenfolge an den Use Case, jede mit genau den Feldern ihres Typs;
// die Server-Nachrichten gehen in Reihenfolge an den Client. Eine ohne Warten
// nachgeschickte Gruppe erreicht den Use Case erst nach dem ReadyForQuery der
// ersten (Pipelining).
func TestExtendedSyncGruppe(t *testing.T) {
	_, fe, rec, s := verbindeExtended(t, context.Background())
	sende(fe,
		&pgproto3.Parse{Name: "s1", Query: "SELECT $1, $2, $3", ParameterOIDs: []uint32{23, 25, 25}},
		&pgproto3.Bind{DestinationPortal: "p", PreparedStatement: "s1", ParameterFormatCodes: []int16{0}, Parameters: [][]byte{[]byte("1"), nil, {}}, ResultFormatCodes: []int16{1}},
		&pgproto3.Describe{ObjectType: 'P', Name: "p"},
		&pgproto3.Execute{Portal: "p", MaxRows: 5},
		&pgproto3.Sync{},
		&pgproto3.Close{ObjectType: 'S', Name: "s1"},
		&pgproto3.Sync{},
	)
	warteAuf(t, rec, "c:sync")
	time.Sleep(50 * time.Millisecond)
	if got := rec.protokoll(); got != "c:parse c:bind c:describe c:execute c:sync" {
		t.Fatalf("vor dem ReadyForQuery übergeben: %q", got)
	}
	want := []model.ClientMessage{
		{Type: model.ClientParse, Statement: "s1", SQL: "SELECT $1, $2, $3", ParamTypes: []uint32{23, 25, 25}},
		{Type: model.ClientBind, Portal: "p", Statement: "s1", ParamFormats: []int16{0}, Params: []model.Value{{Bytes: []byte("1")}, {Null: true}, {Bytes: []byte{}}}, ResultFormats: []int16{1}},
		{Type: model.ClientDescribe, Target: model.TargetPortal, Name: "p"},
		{Type: model.ClientExecute, Portal: "p", MaxRows: 5},
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
	warteAuf(t, rec, "s:ready_for_query c:close c:sync")
	rec.server <- []model.Response{{Type: model.ResponseCloseComplete}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}
	if got := strings.Join(empfange(t, fe, 2), " "); got != "CloseComplete ReadyForQuery" {
		t.Fatalf("an den Client: %s", got)
	}
	rec.mu.Lock()
	closeMsg := rec.client[5]
	rec.mu.Unlock()
	if !reflect.DeepEqual(closeMsg, model.ClientMessage{Type: model.ClientClose, Target: model.TargetStatement, Name: "s1"}) {
		t.Fatalf("Close: %#v", closeMsg)
	}
	sende(fe, &pgproto3.Terminate{})
	if end := rec.lastEnd(t); end != model.EndNormal || s.FirstErrorCode() != "" {
		t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: nach einem Flush gehen die Server-Nachrichten
// an den Client, ohne dass ein Sync folgt; die nächste Client-Nachricht
// erreicht den Use Case nach ihnen.
func TestExtendedFlush(t *testing.T) {
	_, fe, rec, _ := verbindeExtended(t, context.Background())
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
	sende(fe, &pgproto3.Sync{})
	warteAuf(t, rec, "c:sync")
	rec.server <- []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
	empfange(t, fe, 1)
	want := "c:parse c:describe c:flush s:parse_complete s:parameter_description s:no_data c:sync s:ready_for_query"
	if got := rec.protokoll(); got != want {
		t.Fatalf("Protokoll %q", got)
	}
}

// Abdeckung: LH-FA-13/Negative — endet die Client-Verbindung
// oder kommt Terminate, bevor das ReadyForQuery der laufenden
// Extended-Interaktion verarbeitet ist, ist das PGR-E4003; nach dem
// ReadyForQuery ist das Ende regulär. Erreicht das ReadyForQuery den Client
// nicht, endet die Session mit EndLost, eine frühere Server-Nachricht mit
// EndNormal; ein Abbruch des Upstreams ist PGR-E4003.
func TestExtendedAbbruch(t *testing.T) {
	cases := []struct {
		name     string
		ablauf   func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder)
		wantEnd  model.SessionEnd
		wantCode string
	}{
		{"Verbindungsende nach Parse", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			sende(fe, &pgproto3.Parse{Query: "SELECT 1"})
			warteAuf(t, rec, "c:parse")
			client.Close()
		}, model.EndNormal, model.CodeConnectionLost},
		{"Terminate nach Flush", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			sende(fe, &pgproto3.Parse{Query: "SELECT 1"}, &pgproto3.Flush{})
			warteAuf(t, rec, "c:flush")
			sende(fe, &pgproto3.Terminate{})
		}, model.EndNormal, model.CodeConnectionLost},
		{"Verbindungsende nach ReadyForQuery", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			sende(fe, &pgproto3.Sync{})
			warteAuf(t, rec, "c:sync")
			rec.server <- []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
			empfange(t, fe, 1)
			client.Close()
		}, model.EndNormal, ""},
		{"ReadyForQuery nicht zugestellt", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			sende(fe, &pgproto3.Sync{})
			warteAuf(t, rec, "c:sync")
			client.Close()
			rec.server <- []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
		}, model.EndLost, model.CodeConnectionLost},
		{"frühere Server-Nachricht nicht zugestellt", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			sende(fe, &pgproto3.Parse{Query: "SELECT 1"}, &pgproto3.Sync{})
			warteAuf(t, rec, "c:sync")
			client.Close()
			rec.server <- []model.Response{{Type: model.ResponseParseComplete}}
		}, model.EndNormal, model.CodeConnectionLost},
		{"Upstream bricht ab", func(t *testing.T, client net.Conn, fe *pgproto3.Frontend, rec *fakeRecorder) {
			sende(fe, &pgproto3.Sync{})
			warteAuf(t, rec, "c:sync")
			close(rec.server)
			e := fehlerantwort(t, fe)
			if !strings.Contains(e.Message, model.CodeConnectionLost) {
				t.Fatalf("ErrorResponse: %+v", e)
			}
		}, model.EndNormal, model.CodeConnectionLost},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, fe, rec, s := verbindeExtended(t, context.Background())
			c.ablauf(t, client, fe, rec)
			if end := rec.lastEnd(t); end != c.wantEnd {
				t.Fatalf("Session-Ende %v, erwartet %v", end, c.wantEnd)
			}
			if got := s.FirstErrorCode(); got != c.wantCode {
				t.Fatalf("erster Fehler %q, erwartet %q", got, c.wantCode)
			}
		})
	}
}

// Abdeckung: LH-FA-06/Negative — ein Describe mit einer
// anderen Zielart als 'S' oder 'P', eine Server-Nachricht, die der Use Case
// als nicht unterstützt ablehnt, und eine Extended-Nachricht im Replay-Modus
// beenden die Verbindung mit PGR-E6001 (SQLSTATE 0A000); im Record-Modus endet
// die Session mit EndUnsupported.
func TestExtendedNichtUnterstuetzt(t *testing.T) {
	t.Run("Zielart", func(t *testing.T) {
		client, fe, rec, s := verbindeExtended(t, context.Background())
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
	t.Run("Server-Nachricht abgelehnt", func(t *testing.T) {
		_, fe, rec, s := verbindeExtended(t, context.Background())
		rec.serverErr = map[model.ResponseType]error{model.ResponseReadyForQuery: model.Errorf(model.CodeUnsupported, nil, "Form")}
		sende(fe, &pgproto3.Sync{})
		warteAuf(t, rec, "c:sync")
		rec.server <- []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
		e := fehlerantwort(t, fe)
		if e.Code != "0A000" || !strings.Contains(e.Message, model.CodeUnsupported) {
			t.Fatalf("ErrorResponse: %+v", e)
		}
		if end := rec.lastEnd(t); end != model.EndUnsupported || s.FirstErrorCode() != model.CodeUnsupported {
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

// Abdeckung: LH-FA-13/Boundary — endet der Lauf, während eine
// Extended-Interaktion läuft, endet die Session erst nach deren ReadyForQuery,
// ohne Verbindungsfehler; eine ruhende Session endet sofort (LH-FA-13.a).
func TestExtendedHerunterfahren(t *testing.T) {
	t.Run("laufende Interaktion", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		_, fe, rec, s := verbindeExtended(t, ctx)
		sende(fe, &pgproto3.Parse{Query: "SELECT 1"})
		warteAuf(t, rec, "c:parse")
		cancel()
		time.Sleep(50 * time.Millisecond)
		sende(fe, &pgproto3.Bind{}, &pgproto3.Execute{}, &pgproto3.Sync{})
		warteAuf(t, rec, "c:sync")
		rec.server <- []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}
		empfange(t, fe, 3)
		if end := rec.lastEnd(t); end != model.EndNormal || s.FirstErrorCode() != "" {
			t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
		}
		if got := rec.protokoll(); !strings.HasSuffix(got, "s:ready_for_query ende") {
			t.Fatalf("Protokoll %q", got)
		}
	})
	t.Run("ruhend", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		_, _, rec, s := verbindeExtended(t, ctx)
		cancel()
		if end := rec.lastEnd(t); end != model.EndNormal || s.FirstErrorCode() != "" {
			t.Fatalf("Ende %v, Fehler %q", end, s.FirstErrorCode())
		}
	})
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
