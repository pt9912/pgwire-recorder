package postgres

import (
	"io"
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
	_, err := aufbau(conn, startup, zugang{})
	return err
}

// AufbauMit reicht an aufbau mit passwort weiter; zufall liefert die Bytes der
// Nonce, ohne zufall gilt crypto/rand; conn stellt der Test.
func AufbauMit(conn net.Conn, startup map[string]string, passwort string, zufall io.Reader) error {
	z := zugang{passwort: Passwort(passwort)}
	if zufall != nil {
		z.zufall = func(b []byte) { _, _ = io.ReadFull(zufall, b) }
	}
	_, err := aufbau(conn, startup, z)
	return err
}

// ScramBeweis reicht an scramBeweis weiter.
func ScramBeweis(passwort string, salz []byte, iterationen int, authMessage string) (beweis, serverSig []byte, err error) {
	return scramBeweis(passwort, salz, iterationen, authMessage)
}

// LeseServerErste reicht an leseServerErste weiter.
func LeseServerErste(nonce, text string) (servernonce string, salz []byte, iterationen int, grund string) {
	return leseServerErste(nonce, text)
}

// Md5Antwort reicht an md5Antwort weiter.
func Md5Antwort(passwort, benutzer string, salz []byte) string {
	return md5Antwort(passwort, benutzer, salz)
}
