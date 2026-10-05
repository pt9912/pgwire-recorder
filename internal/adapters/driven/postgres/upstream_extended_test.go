package postgres

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// extendedServer nimmt eine Verbindung an, schickt den Aufbau und liest
// danach Client-Nachrichten. Nach jedem Flush oder Sync übergibt er die seit
// dem vorigen gelesenen Nachrichten an empfangen und sendet die nächste Folge
// aus antworten in einem Schreibvorgang; ist sie nil, schließt er die
// Verbindung.
func extendedServer(t *testing.T, antworten [][]pgproto3.BackendMessage) (string, <-chan []pgproto3.FrontendMessage) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	empfangen := make(chan []pgproto3.FrontendMessage, len(antworten))
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		be := pgproto3.NewBackend(conn, conn)
		if _, err := be.ReceiveStartupMessage(); err != nil {
			return
		}
		for _, m := range bereit {
			be.Send(m)
		}
		_ = be.Flush()
		var gruppe []pgproto3.FrontendMessage
		for _, antwort := range antworten {
			for {
				msg, err := be.Receive()
				if err != nil {
					return
				}
				// Die Nachricht gilt nur bis zum nächsten Receive.
				gruppe = append(gruppe, kopieFrontend(t, msg))
				_, flush := msg.(*pgproto3.Flush)
				_, sync := msg.(*pgproto3.Sync)
				if flush || sync {
					break
				}
			}
			empfangen <- gruppe
			gruppe = nil
			if antwort == nil {
				return
			}
			for _, m := range antwort {
				be.Send(m)
			}
			_ = be.Flush()
		}
	}()
	return l.Addr().String(), empfangen
}

func kopieFrontend(t *testing.T, msg pgproto3.FrontendMessage) pgproto3.FrontendMessage {
	t.Helper()
	b, err := msg.Encode(nil)
	if err != nil {
		t.Error(err)
		return nil
	}
	kopie := reflect.New(reflect.TypeOf(msg).Elem()).Interface().(pgproto3.FrontendMessage)
	if err := kopie.Decode(b[5:]); err != nil {
		t.Error(err)
	}
	return kopie
}

// Abdeckung: LH-FA-18/Happy — Record-Hälfte: Send bringt die Client-Nachrichten einer Gruppe
// unverändert und in Reihenfolge auf die Leitung, NULL und leerer Wert eines
// Parameters bleiben unterschieden; Receive bildet die Server-Nachrichten des
// Extended Query Protocol ab, param_types nur an parameter_description.
func TestSendUndReceive(t *testing.T) {
	addr, empfangen := extendedServer(t, [][]pgproto3.BackendMessage{
		{&pgproto3.ParseComplete{}, &pgproto3.ParameterDescription{ParameterOIDs: []uint32{23, 25}}, &pgproto3.NoData{}},
		{&pgproto3.BindComplete{}, &pgproto3.DataRow{Values: [][]byte{[]byte("1")}}, &pgproto3.PortalSuspended{}, &pgproto3.CloseComplete{},
			&pgproto3.ParameterDescription{}, &pgproto3.ReadyForQuery{TxStatus: 'T'}},
	})
	s, _, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	gruppe1 := []model.ClientMessage{
		{Type: model.ClientParse, Statement: "s1", SQL: "SELECT $1, $2", ParamTypes: []uint32{23, 25}},
		{Type: model.ClientDescribe, Target: model.TargetStatement, Name: "s1"},
		{Type: model.ClientFlush},
	}
	if err := s.Send(ctx(t), gruppe1); err != nil {
		t.Fatal(err)
	}
	want1 := []pgproto3.FrontendMessage{
		&pgproto3.Parse{Name: "s1", Query: "SELECT $1, $2", ParameterOIDs: []uint32{23, 25}},
		&pgproto3.Describe{ObjectType: 'S', Name: "s1"},
		&pgproto3.Flush{},
	}
	if got := <-empfangen; !reflect.DeepEqual(got, want1) {
		t.Fatalf("Gruppe 1 auf der Leitung: %#v", got)
	}
	var out []model.Response
	for len(out) < 3 {
		rs, err := s.Receive(ctx(t))
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, rs...)
	}
	wantOut1 := []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseParameterDescription, ParamTypes: []uint32{23, 25}}, {Type: model.ResponseNoData}}
	if !reflect.DeepEqual(out, wantOut1) {
		t.Fatalf("Antworten 1: %#v", out)
	}

	gruppe2 := []model.ClientMessage{
		{Type: model.ClientBind, Portal: "p", Statement: "s1", ParamFormats: []int16{0, 1}, Params: []model.Value{{Null: true}, {}}, ResultFormats: []int16{1}},
		{Type: model.ClientExecute, Portal: "p", MaxRows: 1},
		{Type: model.ClientClose, Target: model.TargetPortal, Name: "p"},
		{Type: model.ClientSync},
	}
	if err := s.Send(ctx(t), gruppe2); err != nil {
		t.Fatal(err)
	}
	want2 := []pgproto3.FrontendMessage{
		&pgproto3.Bind{DestinationPortal: "p", PreparedStatement: "s1", ParameterFormatCodes: []int16{0, 1}, Parameters: [][]byte{nil, {}}, ResultFormatCodes: []int16{1}},
		&pgproto3.Execute{Portal: "p", MaxRows: 1},
		&pgproto3.Close{ObjectType: 'P', Name: "p"},
		&pgproto3.Sync{},
	}
	if got := <-empfangen; !reflect.DeepEqual(got, want2) {
		t.Fatalf("Gruppe 2 auf der Leitung: %#v", got)
	}
	out = nil
	for len(out) == 0 || out[len(out)-1].Type != model.ResponseReadyForQuery {
		rs, err := s.Receive(ctx(t))
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, rs...)
	}
	wantOut2 := []model.Response{
		{Type: model.ResponseBindComplete},
		{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte("1")}}},
		{Type: model.ResponsePortalSuspended},
		{Type: model.ResponseCloseComplete},
		{Type: model.ResponseParameterDescription},
		{Type: model.ResponseReadyForQuery, TxStatus: "T"},
	}
	if !reflect.DeepEqual(out, wantOut2) {
		t.Fatalf("Antworten 2: %#v", out)
	}
}

// Receive liest, was schon im Puffer liegt, in einem Aufruf, aber nicht über
// ReadyForQuery hinaus.
func TestReceiveEndetMitReadyForQuery(t *testing.T) {
	addr, _ := extendedServer(t, [][]pgproto3.BackendMessage{
		{&pgproto3.ParseComplete{}, &pgproto3.ReadyForQuery{TxStatus: 'I'}, &pgproto3.ParseComplete{}},
	})
	s, _, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Send(ctx(t), []model.ClientMessage{{Type: model.ClientSync}}); err != nil {
		t.Fatal(err)
	}
	var out []model.Response
	for len(out) == 0 || out[len(out)-1].Type != model.ResponseReadyForQuery {
		rs, err := s.Receive(ctx(t))
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, rs...)
	}
	if len(out) != 2 {
		t.Fatalf("über ReadyForQuery hinaus gelesen: %#v", out)
	}
}

// Eine Server-Nachricht, die
// eine Extended-Interaktion nicht trägt (COPY), ist PGR-E6001; endet die
// Verbindung vor ReadyForQuery, ist das PGR-E4003; eine Client-Nachricht ohne
// gültige Zielart wird nicht gesendet (PGR-E1000).
func TestReceiveFehler(t *testing.T) {
	t.Run("COPY", func(t *testing.T) {
		addr, _ := extendedServer(t, [][]pgproto3.BackendMessage{{&pgproto3.CopyInResponse{}}})
		s, _, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
		if err != nil {
			t.Fatal(err)
		}
		defer s.Close()
		if err := s.Send(ctx(t), []model.ClientMessage{{Type: model.ClientExecute}, {Type: model.ClientSync}}); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Receive(ctx(t)); code(err) != model.CodeUnsupported {
			t.Fatalf("erwartet %s, erhalten %v", model.CodeUnsupported, err)
		}
	})
	t.Run("Verbindungsende", func(t *testing.T) {
		addr, _ := extendedServer(t, [][]pgproto3.BackendMessage{nil})
		s, _, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
		if err != nil {
			t.Fatal(err)
		}
		defer s.Close()
		if err := s.Send(ctx(t), []model.ClientMessage{{Type: model.ClientSync}}); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Receive(ctx(t)); code(err) != model.CodeConnectionLost {
			t.Fatalf("erwartet %s, erhalten %v", model.CodeConnectionLost, err)
		}
	})
	t.Run("Zielart", func(t *testing.T) {
		if _, err := toFrontendMessage(model.ClientMessage{Type: model.ClientDescribe}); code(err) != model.CodeInternal {
			t.Fatalf("erwartet %s, erhalten %v", model.CodeInternal, err)
		}
		if _, err := toFrontendMessage(model.ClientMessage{Type: "bogus"}); code(err) != model.CodeInternal {
			t.Fatalf("erwartet %s, erhalten %v", model.CodeInternal, err)
		}
	})
}

// rohServer nimmt eine Verbindung an, schickt den Aufbau und übergibt die
// Verbindung danach an weiter. Er liefert die Adresse.
func rohServer(t *testing.T, weiter func(be *pgproto3.Backend, conn net.Conn)) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		be := pgproto3.NewBackend(conn, conn)
		if _, err := be.ReceiveStartupMessage(); err != nil {
			return
		}
		for _, m := range bereit {
			be.Send(m)
		}
		_ = be.Flush()
		weiter(be, conn)
	}()
	return l.Addr().String()
}

// grosserBind ist eine Gruppe mit einem Parameter, der die Puffer der
// Verbindung weit übersteigt.
func grosserBind() []model.ClientMessage {
	return []model.ClientMessage{
		{Type: model.ClientBind, Params: []model.Value{{Bytes: make([]byte, 32<<20)}}},
		{Type: model.ClientSync},
	}
}

// fertigBinnen meldet einen Fehler, wenn f nicht binnen d zurückkehrt.
func fertigBinnen(t *testing.T, d time.Duration, was string, f func()) {
	t.Helper()
	fertig := make(chan struct{})
	go func() {
		defer close(fertig)
		f()
	}()
	select {
	case <-fertig:
	case <-time.After(d):
		t.Fatalf("%s kehrt nicht binnen %v zurück", was, d)
	}
}

// Send und Receive laufen gleichzeitig: Der Server schreibt zuerst eine
// Ausgabe, die die Puffer der Verbindung übersteigt, und liest erst danach;
// ein Send, das bis zu ihrem Lesen wartete, käme nie zurück.
func TestSendUndReceiveGleichzeitig(t *testing.T) {
	addr := rohServer(t, func(be *pgproto3.Backend, conn net.Conn) {
		zeile := &pgproto3.DataRow{Values: [][]byte{make([]byte, 1<<20)}}
		for i := 0; i < 32; i++ {
			be.Send(zeile)
			if err := be.Flush(); err != nil {
				return
			}
		}
		be.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		_ = be.Flush()
		for {
			if _, err := be.Receive(); err != nil {
				return
			}
		}
	})
	s, _, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	sendeFehler := make(chan error, 1)
	go func() { sendeFehler <- s.Send(ctx(t), grosserBind()) }()
	fertigBinnen(t, 10*time.Second, "Receive bis ReadyForQuery", func() {
		for {
			rs, err := s.Receive(ctx(t))
			if err != nil {
				t.Error(err)
				return
			}
			if rs[len(rs)-1].Type == model.ResponseReadyForQuery {
				return
			}
		}
	})
	fertigBinnen(t, 10*time.Second, "Send", func() {
		if err := <-sendeFehler; err != nil {
			t.Error(err)
		}
	})
}

// Close kehrt zurück, auch wenn der Server nicht liest, und beendet ein Send,
// das an diesem Server wartet, und ein Receive, das auf ihn wartet, mit einem
// Fehler.
func TestCloseBeendetWartende(t *testing.T) {
	stumm := make(chan struct{})
	t.Cleanup(func() { close(stumm) })
	addr := rohServer(t, func(*pgproto3.Backend, net.Conn) { <-stumm })
	s, _, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
	if err != nil {
		t.Fatal(err)
	}
	sendeFehler := make(chan error, 1)
	go func() { sendeFehler <- s.Send(ctx(t), grosserBind()) }()
	empfangsFehler := make(chan error, 1)
	go func() {
		_, err := s.Receive(ctx(t))
		empfangsFehler <- err
	}()
	time.Sleep(100 * time.Millisecond)
	select {
	case err := <-sendeFehler:
		t.Fatalf("Send kehrte zurück, obwohl der Server nicht liest: %v", err)
	default:
	}
	fertigBinnen(t, 2*time.Second, "Close", func() { _ = s.Close() })
	for name, ch := range map[string]chan error{"Send": sendeFehler, "Receive": empfangsFehler} {
		select {
		case err := <-ch:
			if err == nil {
				t.Fatalf("%s nach Close ohne Fehler", name)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s endet nicht nach Close", name)
		}
	}
}

// Close kehrt auch zurück, wenn es die Schreibsperre bekommt, der Server aber
// nicht liest: das Terminate hat eine Schreibfrist.
func TestCloseMitSchreibfrist(t *testing.T) {
	a, b := net.Pipe()
	defer b.Close()
	s := &session{conn: a, fe: pgproto3.NewFrontend(a, a)}
	fertigBinnen(t, 2*time.Second, "Close", func() { _ = s.Close() })
}
