package postgres_test

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// fakeServer nimmt eine Verbindung an, liest die StartupMessage und schickt die
// gegebenen Nachrichten; danach beantwortet er jede Query mit antwort.
func fakeServer(t *testing.T, aufbau []pgproto3.BackendMessage, antwort []pgproto3.BackendMessage) string {
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
		for _, m := range aufbau {
			be.Send(m)
		}
		_ = be.Flush()
		for {
			msg, err := be.Receive()
			if err != nil {
				return
			}
			if _, ok := msg.(*pgproto3.Query); ok {
				for _, m := range antwort {
					be.Send(m)
				}
				_ = be.Flush()
			}
		}
	}()
	return l.Addr().String()
}

func code(err error) string {
	var me *model.Error
	if errors.As(err, &me) {
		return me.Code
	}
	return ""
}

// bereit ist ein Verbindungsaufbau ohne Anmeldeverfahren bis zum ersten
// ReadyForQuery.
func bereit() []pgproto3.BackendMessage {
	return []pgproto3.BackendMessage{
		&pgproto3.AuthenticationOk{},
		&pgproto3.ParameterStatus{Name: "server_version", Value: "17.0"},
		&pgproto3.BackendKeyData{ProcessID: 42, SecretKey: []byte{1, 2, 3, 4}},
		&pgproto3.ReadyForQuery{TxStatus: 'I'},
	}
}

func ctx(t *testing.T) context.Context {
	t.Helper()
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return c
}

// Abdeckung: LH-FA-02/Happy — Verbindungsaufbau und Anfrage
// liefern die Serverantworten in Reihenfolge; die Abbruchkennung des Servers
// (BackendKeyData) ist nicht darunter (SPEC-004).
func TestOpenUndQuery(t *testing.T) {
	addr := fakeServer(t, bereit(), []pgproto3.BackendMessage{
		&pgproto3.RowDescription{Fields: []pgproto3.FieldDescription{{Name: []byte("a"), DataTypeOID: 23, DataTypeSize: 4, TypeModifier: -1}}},
		&pgproto3.DataRow{Values: [][]byte{[]byte("1")}},
		&pgproto3.DataRow{Values: [][]byte{nil}},
		&pgproto3.CommandComplete{CommandTag: []byte("SELECT 2")},
		&pgproto3.ReadyForQuery{TxStatus: 'I'},
	})
	up := &postgres.Upstream{Address: addr}
	s, aufbau, err := up.Open(ctx(t), map[string]string{"user": "app"})
	if err != nil || s == nil {
		t.Fatalf("Open: %v", err)
	}
	for _, r := range aufbau {
		if r.Type != model.ResponseParameterStatus && r.Type != model.ResponseReadyForQuery {
			t.Fatalf("unerwartete Aufbau-Antwort: %#v", r)
		}
	}
	out, err := s.Query(ctx(t), "SELECT a")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 5 || string(out[1].Values[0].Bytes) != "1" || !out[2].Values[0].Null || out[3].Tag != "SELECT 2" || out[4].TxStatus != "I" {
		t.Fatalf("Antworten: %#v", out)
	}
	_ = s.Close()
}

// Abdeckung: LH-FA-05/Negative — verlangt der Server ein Anmeldeverfahren, endet
// der Aufbau mit PGR-E6001; eine Fehlerantwort im Aufbau geht ohne Session an
// den Client; eine nicht unterstützte Serverantwort ist PGR-E6001.
func TestNichtVermittelbar(t *testing.T) {
	t.Run("Anmeldeverfahren", func(t *testing.T) {
		addr := fakeServer(t, []pgproto3.BackendMessage{&pgproto3.AuthenticationCleartextPassword{}}, nil)
		s, _, err := (&postgres.Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
		if s != nil || code(err) != model.CodeUnsupported {
			t.Fatalf("erwartet %s, erhalten %v", model.CodeUnsupported, err)
		}
	})
	t.Run("Fehlerantwort im Aufbau", func(t *testing.T) {
		addr := fakeServer(t, []pgproto3.BackendMessage{&pgproto3.AuthenticationOk{}, &pgproto3.ErrorResponse{Severity: "FATAL", Code: "3D000", Message: "database does not exist"}}, nil)
		s, out, err := (&postgres.Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
		if s != nil || err != nil || len(out) != 1 || out[0].Fields["C"] != "3D000" {
			t.Fatalf("Session %v, Antworten %#v, Fehler %v", s, out, err)
		}
	})
	t.Run("COPY", func(t *testing.T) {
		addr := fakeServer(t, bereit(), []pgproto3.BackendMessage{&pgproto3.CopyOutResponse{}})
		s, _, err := (&postgres.Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := s.Query(ctx(t), "COPY t TO STDOUT"); code(err) != model.CodeUnsupported {
			t.Fatalf("erwartet %s, erhalten %v", model.CodeUnsupported, err)
		}
		_ = s.Close()
	})
}

// Abdeckung: LH-FA-02/Negative — ein nicht erreichbarer Server ist PGR-E4002.
func TestUpstreamNichtErreichbar(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	if _, _, err := (&postgres.Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"}); code(err) != model.CodeUpstream {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeUpstream, err)
	}
}

// abbruchServer schickt nach dem Aufbau die Nachrichten aus zwischen und
// beantwortet danach jede Query mit der nächsten Folge aus antworten; nach der
// letzten schließt er die Verbindung.
func abbruchServer(t *testing.T, zwischen []pgproto3.BackendMessage, antworten ...[]pgproto3.BackendMessage) string {
	t.Helper()
	return rohServer(t, func(be *pgproto3.Backend, _ net.Conn) {
		for _, m := range zwischen {
			be.Send(m)
		}
		_ = be.Flush()
		for _, antwort := range antworten {
			for {
				msg, err := be.Receive()
				if err != nil {
					return
				}
				if _, ok := msg.(*pgproto3.Query); ok {
					break
				}
			}
			for _, m := range antwort {
				be.Send(m)
			}
			_ = be.Flush()
		}
	})
}

// beendet ist die Fehlerantwort, mit der PostgreSQL eine Verbindung auf
// Anweisung des Administrators beendet.
func beendet() *pgproto3.ErrorResponse {
	return &pgproto3.ErrorResponse{Severity: "FATAL", SeverityUnlocalized: "FATAL", Code: "57P01", Message: "terminating connection due to administrator command"}
}

// ergebnis sind die Antworten auf SELECT 1 ohne ReadyForQuery.
func ergebnis() []pgproto3.BackendMessage {
	return []pgproto3.BackendMessage{
		&pgproto3.RowDescription{Fields: []pgproto3.FieldDescription{{Name: []byte("a"), DataTypeOID: 23, DataTypeSize: 4, TypeModifier: -1}}},
		&pgproto3.DataRow{Values: [][]byte{[]byte("1")}},
		&pgproto3.CommandComplete{CommandTag: []byte("SELECT 1")},
	}
}

func mit(folge []pgproto3.BackendMessage, weitere ...pgproto3.BackendMessage) []pgproto3.BackendMessage {
	return append(append([]pgproto3.BackendMessage{}, folge...), weitere...)
}

// pruefeAbbruch prüft den Meldungscode von err und bei PGR-E6001, dass der
// Text SQLSTATE und Meldung der Fehlerantwort nennt.
func pruefeAbbruch(t *testing.T, err error, want string, fehler *pgproto3.ErrorResponse) {
	t.Helper()
	if code(err) != want {
		t.Fatalf("erwartet %s, erhalten %v", want, err)
	}
	if want == model.CodeUnsupported && (!strings.Contains(err.Error(), fehler.Code) || !strings.Contains(err.Error(), fehler.Message)) {
		t.Fatalf("Text ohne SQLSTATE %s und Meldung %q: %v", fehler.Code, fehler.Message, err)
	}
}

// oeffne baut die Session zu addr auf und begrenzt ihr Lesen auf fünf
// Sekunden: Ein Lesen über das Ende der Antworten hinaus endet so mit einem
// Fehler, statt auf den Server zu warten.
func oeffne(t *testing.T, addr string) driven.UpstreamSession {
	t.Helper()
	s, _, err := (&postgres.Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	if err := postgres.Verbindung(s).SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	return s
}

// Abdeckung: LH-FA-11/Negative — endet die Verbindung zum Upstream nach einer
// ErrorResponse vor dem ReadyForQuery, liefert der Upstream-Adapter PGR-E6001
// mit SQLSTATE und Meldung der letzten ErrorResponse im Text: bei FATAL, ERROR und
// PANIC, nach Ergebnissen einer einfachen Anfrage, als Nachricht zwischen zwei
// Interaktionen und im Extended Query Protocol nach einer schon gelieferten
// ErrorResponse. Ohne ErrorResponse und nach einer ErrorResponse, auf die ein
// ReadyForQuery folgte, ist es PGR-E4003.
func TestFehlerantwortVorDemAbbruch(t *testing.T) {
	for _, c := range []struct {
		name   string
		fehler *pgproto3.ErrorResponse
	}{
		{"nach Ergebnissen, FATAL", beendet()},
		{"Schweregrad ERROR", &pgproto3.ErrorResponse{Severity: "ERROR", SeverityUnlocalized: "ERROR", Code: "XX000", Message: "interner Fehler"}},
		{"Schweregrad PANIC", &pgproto3.ErrorResponse{Severity: "PANIC", SeverityUnlocalized: "PANIC", Code: "XX000", Message: "Panik"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			s := oeffne(t, abbruchServer(t, nil, mit(ergebnis(), c.fehler)))
			out, err := s.Query(ctx(t), "SELECT 1; SELECT pg_terminate_backend(pg_backend_pid())")
			pruefeAbbruch(t, err, model.CodeUnsupported, c.fehler)
			if len(out) != 4 || out[3].Type != model.ResponseErrorResponse {
				t.Fatalf("Antworten vor dem Abbruch: %#v", out)
			}
		})
	}
	t.Run("zwei Fehlerantworten, die letzte zählt", func(t *testing.T) {
		erste := &pgproto3.ErrorResponse{Severity: "ERROR", Code: "XX001", Message: "erste"}
		s := oeffne(t, abbruchServer(t, nil, mit(ergebnis(), erste, beendet())))
		_, err := s.Query(ctx(t), "SELECT 1")
		pruefeAbbruch(t, err, model.CodeUnsupported, beendet())
		if strings.Contains(err.Error(), "XX001") {
			t.Fatalf("Text nennt die erste Fehlerantwort: %v", err)
		}
	})
	t.Run("ohne Fehlerantwort", func(t *testing.T) {
		s := oeffne(t, abbruchServer(t, nil, ergebnis()))
		_, err := s.Query(ctx(t), "SELECT 1")
		pruefeAbbruch(t, err, model.CodeConnectionLost, nil)
	})
	t.Run("ReadyForQuery nach der Fehlerantwort", func(t *testing.T) {
		s := oeffne(t, abbruchServer(t, nil, []pgproto3.BackendMessage{beendet(), &pgproto3.ReadyForQuery{TxStatus: 'I'}}, nil))
		if _, err := s.Query(ctx(t), "SELECT 1/0"); err != nil {
			t.Fatal(err)
		}
		_, err := s.Query(ctx(t), "SELECT 2")
		pruefeAbbruch(t, err, model.CodeConnectionLost, nil)
	})
	t.Run("zwischen zwei Interaktionen", func(t *testing.T) {
		s := oeffne(t, abbruchServer(t, []pgproto3.BackendMessage{beendet()}, nil))
		out, err := s.Query(ctx(t), "SELECT 1")
		pruefeAbbruch(t, err, model.CodeUnsupported, beendet())
		if len(out) != 1 || out[0].Fields["C"] != "57P01" {
			t.Fatalf("Antworten vor dem Abbruch: %#v", out)
		}
	})
	t.Run("Extended", func(t *testing.T) {
		addr, _ := extendedServer(t, [][]pgproto3.BackendMessage{{beendet()}, nil})
		s := oeffne(t, addr)
		if err := s.Send(ctx(t), []model.ClientMessage{{Type: model.ClientParse, SQL: "SELECT 1"}, {Type: model.ClientFlush}}); err != nil {
			t.Fatal(err)
		}
		rs, err := s.Receive(ctx(t))
		if err != nil || len(rs) != 1 || rs[0].Type != model.ResponseErrorResponse || rs[0].Fields["C"] != "57P01" {
			t.Fatalf("Receive: %#v, %v", rs, err)
		}
		if err := s.Send(ctx(t), []model.ClientMessage{{Type: model.ClientSync}}); err != nil {
			t.Fatal(err)
		}
		_, err = s.Receive(ctx(t))
		pruefeAbbruch(t, err, model.CodeUnsupported, beendet())
	})
}

// Eine NoticeResponse und ein ParameterStatus, die der Server nach dem
// ReadyForQuery einer Interaktion sendet, liefert der Upstream-Adapter mit der
// nächsten Interaktion vor deren übrigen Antworten, bei einer einfachen Anfrage
// wie im Extended Query Protocol (LH-FA-05.a).
func TestNachrichtenZwischenInteraktionen(t *testing.T) {
	hinweis := &pgproto3.NoticeResponse{Severity: "NOTICE", Code: "00000", Message: "zwischen"}
	status := &pgproto3.ParameterStatus{Name: "application_name", Value: "zwischen"}
	wantZwischen := []model.Response{
		{Type: model.ResponseNoticeResponse, Fields: map[string]string{"S": "NOTICE", "C": "00000", "M": "zwischen"}},
		{Type: model.ResponseParameterStatus, Name: "application_name", Value: "zwischen"},
	}
	rfq := model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"}
	t.Run("einfache Anfrage", func(t *testing.T) {
		addr := abbruchServer(t, nil,
			[]pgproto3.BackendMessage{&pgproto3.CommandComplete{CommandTag: []byte("SET")}, &pgproto3.ReadyForQuery{TxStatus: 'I'}, hinweis, status},
			[]pgproto3.BackendMessage{&pgproto3.CommandComplete{CommandTag: []byte("SET")}, &pgproto3.ReadyForQuery{TxStatus: 'I'}})
		s := oeffne(t, addr)
		erste, err := s.Query(ctx(t), "SET a = 1")
		if err != nil {
			t.Fatal(err)
		}
		zweite, err := s.Query(ctx(t), "SET b = 2")
		if err != nil {
			t.Fatal(err)
		}
		set := model.Response{Type: model.ResponseCommandComplete, Tag: "SET"}
		if want := []model.Response{set, rfq}; !reflect.DeepEqual(erste, want) {
			t.Fatalf("erste Interaktion: %#v", erste)
		}
		if want := append(append([]model.Response{}, wantZwischen...), set, rfq); !reflect.DeepEqual(zweite, want) {
			t.Fatalf("zweite Interaktion: %#v", zweite)
		}
	})
	t.Run("Extended", func(t *testing.T) {
		addr, _ := extendedServer(t, [][]pgproto3.BackendMessage{
			{&pgproto3.ParseComplete{}, &pgproto3.ReadyForQuery{TxStatus: 'I'}, hinweis, status},
			{&pgproto3.ParseComplete{}, &pgproto3.ReadyForQuery{TxStatus: 'I'}},
		})
		s := oeffne(t, addr)
		gruppe := []model.ClientMessage{{Type: model.ClientParse, SQL: "SELECT 1"}, {Type: model.ClientSync}}
		lies := func() []model.Response {
			var out []model.Response
			for len(out) == 0 || out[len(out)-1].Type != model.ResponseReadyForQuery {
				rs, err := s.Receive(ctx(t))
				if err != nil {
					t.Fatal(err)
				}
				out = append(out, rs...)
			}
			return out
		}
		if err := s.Send(ctx(t), gruppe); err != nil {
			t.Fatal(err)
		}
		parse := model.Response{Type: model.ResponseParseComplete}
		if erste := lies(); !reflect.DeepEqual(erste, []model.Response{parse, rfq}) {
			t.Fatalf("erste Interaktion: %#v", erste)
		}
		if err := s.Send(ctx(t), gruppe); err != nil {
			t.Fatal(err)
		}
		if zweite, want := lies(), append(append([]model.Response{}, wantZwischen...), parse, rfq); !reflect.DeepEqual(zweite, want) {
			t.Fatalf("zweite Interaktion: %#v", zweite)
		}
	})
}

// rohFelder ist der Rumpf einer Fehler- oder Hinweisantwort mit leerer Meldung,
// leerem Detail, einer Position 0, einer internen Position ohne Zahl, einer
// leeren Zeilennummer und dem unbekannten Feldcode Y mit leerem Wert.
const rohFelder = "SERROR\x00VERROR\x00CXX000\x00M\x00D\x00P0\x00pabc\x00L\x00Rfn\x00Y\x00\x00"

func rohNachricht(typ byte, rumpf string) []byte {
	b := []byte{typ, 0, 0, 0, 0}
	binary.BigEndian.PutUint32(b[1:], uint32(4+len(rumpf)))
	return append(b, rumpf...)
}

// Ein benanntes Diagnosefeld mit leerem Wert und ein Zahlenfeld mit 0 oder
// ohne Zahl fehlen in den Feldern, die der Upstream-Adapter für ErrorResponse
// und NoticeResponse liefert; ein unbekannter Feldcode steht auch mit leerem
// Wert darin (LH-FA-11.a).
func TestDiagnosefelder(t *testing.T) {
	addr := rohServer(t, func(be *pgproto3.Backend, conn net.Conn) {
		if _, err := be.Receive(); err != nil {
			return
		}
		var b []byte
		b = append(b, rohNachricht('E', rohFelder)...)
		b = append(b, rohNachricht('N', rohFelder)...)
		b = append(b, rohNachricht('Z', "I")...)
		_, _ = conn.Write(b)
		_, _ = be.Receive()
	})
	s := oeffne(t, addr)
	out, err := s.Query(ctx(t), "SELECT 1")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"S": "ERROR", "V": "ERROR", "C": "XX000", "R": "fn", "Y": ""}
	if len(out) != 3 || out[0].Type != model.ResponseErrorResponse || out[1].Type != model.ResponseNoticeResponse ||
		!reflect.DeepEqual(out[0].Fields, want) || !reflect.DeepEqual(out[1].Fields, want) {
		t.Fatalf("Antworten: %#v", out)
	}
}

// schreibhaelfteZu schließt die Schreibhälfte der Upstream-Verbindung, sodass
// das nächste Senden sicher scheitert.
func schreibhaelfteZu(t *testing.T, s driven.UpstreamSession) {
	t.Helper()
	if err := postgres.Verbindung(s).(*net.TCPConn).CloseWrite(); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-11/Negative — scheitert im Extended Query Protocol das
// Senden, nachdem die Leserichtung eine ErrorResponse der laufenden
// Interaktion gelesen hat, liefert Send PGR-E6001 mit deren SQLSTATE; ohne
// gelesene ErrorResponse und nach einem ReadyForQuery ist es PGR-E4003.
func TestSendNachFehlerantwort(t *testing.T) {
	parse := []model.ClientMessage{{Type: model.ClientParse, SQL: "SELECT 1"}, {Type: model.ClientFlush}}
	sync := []model.ClientMessage{{Type: model.ClientSync}}
	lies := func(t *testing.T, s driven.UpstreamSession, bis model.ResponseType) {
		t.Helper()
		for {
			rs, err := s.Receive(ctx(t))
			if err != nil {
				t.Fatal(err)
			}
			if rs[len(rs)-1].Type == bis {
				return
			}
		}
	}
	t.Run("ErrorResponse gelesen", func(t *testing.T) {
		addr, _ := extendedServer(t, [][]pgproto3.BackendMessage{{beendet()}})
		s := oeffne(t, addr)
		if err := s.Send(ctx(t), parse); err != nil {
			t.Fatal(err)
		}
		lies(t, s, model.ResponseErrorResponse)
		schreibhaelfteZu(t, s)
		pruefeAbbruch(t, s.Send(ctx(t), sync), model.CodeUnsupported, beendet())
	})
	t.Run("ohne ErrorResponse", func(t *testing.T) {
		addr, _ := extendedServer(t, [][]pgproto3.BackendMessage{{&pgproto3.ParseComplete{}}})
		s := oeffne(t, addr)
		if err := s.Send(ctx(t), parse); err != nil {
			t.Fatal(err)
		}
		lies(t, s, model.ResponseParseComplete)
		schreibhaelfteZu(t, s)
		pruefeAbbruch(t, s.Send(ctx(t), sync), model.CodeConnectionLost, nil)
	})
	t.Run("ErrorResponse und ReadyForQuery gelesen", func(t *testing.T) {
		addr, _ := extendedServer(t, [][]pgproto3.BackendMessage{{beendet(), &pgproto3.ReadyForQuery{TxStatus: 'I'}}})
		s := oeffne(t, addr)
		if err := s.Send(ctx(t), append(parse[:1:1], sync...)); err != nil {
			t.Fatal(err)
		}
		lies(t, s, model.ResponseReadyForQuery)
		schreibhaelfteZu(t, s)
		pruefeAbbruch(t, s.Send(ctx(t), parse), model.CodeConnectionLost, nil)
	})
}

// Abdeckung: LH-FA-05/Negative — eine Serverantwort, die die Bibliothek nicht
// lesen kann, während die Verbindung besteht (unbekannter Nachrichtentyp,
// ReadyForQuery ohne Status), ist PGR-E6001 mit „nicht lesbar“ im Text, auch
// nach einer ErrorResponse derselben Interaktion, deren SQLSTATE der Text dann
// nicht nennt.
func TestNichtLesbareServerantwort(t *testing.T) {
	for _, c := range []struct {
		name   string
		vorher []pgproto3.BackendMessage
		roh    []byte
	}{
		{"unbekannter Typ", nil, rohNachricht('Y', "")},
		{"unbekannter Typ nach ErrorResponse", []pgproto3.BackendMessage{beendet()}, rohNachricht('Y', "")},
		{"ReadyForQuery ohne Status", nil, rohNachricht('Z', "")},
	} {
		t.Run(c.name, func(t *testing.T) {
			addr := rohServer(t, func(be *pgproto3.Backend, conn net.Conn) {
				if _, err := be.Receive(); err != nil {
					return
				}
				for _, m := range c.vorher {
					be.Send(m)
				}
				_ = be.Flush()
				_, _ = conn.Write(c.roh)
				_, _ = be.Receive()
			})
			s := oeffne(t, addr)
			_, err := s.Query(ctx(t), "SELECT 1")
			if code(err) != model.CodeUnsupported || !strings.Contains(err.Error(), "nicht lesbar") || strings.Contains(err.Error(), "57P01") {
				t.Fatalf("erwartet %s „nicht lesbar“ ohne 57P01, erhalten %v", model.CodeUnsupported, err)
			}
		})
	}
}
