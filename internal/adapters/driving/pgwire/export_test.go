package pgwire

import (
	"context"
	"net"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Codes der Sonderanfragen beim Verbindungsaufbau, die Hauptnummer, aus der sie
// gebildet sind, und die Frist für die Fehlerantwort beim Ende einer Session.
const (
	CodeCancelRequest = codeCancelRequest
	CodeSSLRequest    = codeSSLRequest
	CodeGSSEncRequest = codeGSSEncRequest
	MajorSpezial      = majorSpezial
	MeldeFrist        = meldeFrist
)

// Handle bedient conn mit s, wie Serve eine angenommene Verbindung bedient.
func Handle(ctx context.Context, s *Server, conn net.Conn) {
	s.handle(ctx, conn)
}

// Fail reicht an die Methode fail von s weiter.
func Fail(s *Server, be *pgproto3.Backend, err error) {
	s.fail(be, err)
}

// Note reicht an die Methode note von s weiter.
func Note(s *Server, err error) model.Meldung {
	return s.note(err)
}

// ToClientMessage reicht an toClientMessage weiter.
func ToClientMessage(msg pgproto3.FrontendMessage) (model.ClientMessage, error) {
	return toClientMessage(msg)
}

// ToMessage reicht an toMessage weiter.
func ToMessage(r model.Response) (pgproto3.BackendMessage, error) {
	return toMessage(r)
}
