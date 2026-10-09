package postgres

import (
	"net"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// Export-Test-Brücke nach SPEC-049 Punkt 7: Sie reicht nur an Unexportiertes
// des Pakets weiter.

// ToFrontendMessage reicht an toFrontendMessage weiter.
func ToFrontendMessage(m model.ClientMessage) (pgproto3.FrontendMessage, error) {
	return toFrontendMessage(m)
}

// NeueSession reicht an newSession weiter; conn stellt der Test.
func NeueSession(conn net.Conn) driven.UpstreamSession {
	return newSession(conn)
}

// Verbindung liefert die Verbindung der übergebenen Session, die Open oder
// NeueSession erzeugt hat.
func Verbindung(s driven.UpstreamSession) net.Conn {
	return s.(*session).conn
}

// NeueEinspielSession reicht an neueEinspielSession weiter, lesend und
// schreibend auf conn; conn stellt der Test.
func NeueEinspielSession(conn net.Conn) driven.EinspielSession {
	return neueEinspielSession(conn, conn)
}

// Aufbau reicht an aufbau weiter; conn stellt der Test.
func Aufbau(conn net.Conn, startup map[string]string) error {
	_, err := aufbau(conn, startup)
	return err
}
