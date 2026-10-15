package postgres_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-05/Negative — jede der Serverantworten CopyInResponse,
// CopyOutResponse, CopyBothResponse und NotificationResponse ist PGR-E6001, gleich ob sie
// vor dem ReadyForQuery einer Anfrage eintrifft oder im Verbindungsaufbau; die
// Antwort der Anfrage kommt nicht zurück (LH-FA-05.e).
func TestServerNachrichtenNichtUnterstuetzt(t *testing.T) {
	for _, n := range []struct {
		name string
		msg  pgproto3.BackendMessage
	}{
		{"CopyInResponse", &pgproto3.CopyInResponse{}},
		{"CopyOutResponse", &pgproto3.CopyOutResponse{}},
		{"CopyBothResponse", &pgproto3.CopyBothResponse{}},
		{"NotificationResponse", &pgproto3.NotificationResponse{PID: 1, Channel: "k", Payload: "p"}},
	} {
		t.Run(n.name+"/Anfrage", func(t *testing.T) {
			addr := fakeServer(t, bereit(), []pgproto3.BackendMessage{n.msg, &pgproto3.CommandComplete{CommandTag: []byte("SELECT 1")}, &pgproto3.ReadyForQuery{TxStatus: 'I'}})
			s, _, err := (&postgres.Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"}, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer s.Close()
			if out, err := s.Query(ctx(t), "SELECT 1"); code(err) != model.CodeUnsupported || len(out) != 0 {
				t.Fatalf("erwartet %s ohne Antworten, erhalten %v, %#v", model.CodeUnsupported, err, out)
			}
		})
		t.Run(n.name+"/Aufbau", func(t *testing.T) {
			addr := fakeServer(t, append([]pgproto3.BackendMessage{&pgproto3.AuthenticationOk{}, n.msg}, bereit()[1:]...), nil)
			if s, _, err := (&postgres.Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"}, nil); s != nil || code(err) != model.CodeUnsupported {
				t.Fatalf("erwartet %s ohne Session, erhalten %v, %v", model.CodeUnsupported, s, err)
			}
		})
	}
}
