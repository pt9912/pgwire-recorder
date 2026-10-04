package postgres

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
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

var bereit = []pgproto3.BackendMessage{
	&pgproto3.AuthenticationOk{},
	&pgproto3.ParameterStatus{Name: "server_version", Value: "17.0"},
	&pgproto3.BackendKeyData{ProcessID: 42, SecretKey: []byte{1, 2, 3, 4}},
	&pgproto3.ReadyForQuery{TxStatus: 'I'},
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
	addr := fakeServer(t, bereit, []pgproto3.BackendMessage{
		&pgproto3.RowDescription{Fields: []pgproto3.FieldDescription{{Name: []byte("a"), DataTypeOID: 23, DataTypeSize: 4, TypeModifier: -1}}},
		&pgproto3.DataRow{Values: [][]byte{[]byte("1")}},
		&pgproto3.DataRow{Values: [][]byte{nil}},
		&pgproto3.CommandComplete{CommandTag: []byte("SELECT 2")},
		&pgproto3.ReadyForQuery{TxStatus: 'I'},
	})
	up := &Upstream{Address: addr}
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
		s, _, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
		if s != nil || code(err) != model.CodeUnsupported {
			t.Fatalf("erwartet %s, erhalten %v", model.CodeUnsupported, err)
		}
	})
	t.Run("Fehlerantwort im Aufbau", func(t *testing.T) {
		addr := fakeServer(t, []pgproto3.BackendMessage{&pgproto3.AuthenticationOk{}, &pgproto3.ErrorResponse{Severity: "FATAL", Code: "3D000", Message: "database does not exist"}}, nil)
		s, out, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
		if s != nil || err != nil || len(out) != 1 || out[0].Fields["C"] != "3D000" {
			t.Fatalf("Session %v, Antworten %#v, Fehler %v", s, out, err)
		}
	})
	t.Run("COPY", func(t *testing.T) {
		addr := fakeServer(t, bereit, []pgproto3.BackendMessage{&pgproto3.CopyOutResponse{}})
		s, _, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"})
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
	if _, _, err := (&Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"}); code(err) != model.CodeUpstream {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeUpstream, err)
	}
}
