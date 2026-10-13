package pgwire

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"net"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// TLS ist die Verschlüsselung zum Client (LH-FA-23.a): die Konfiguration des
// TLS-Servers und die Zusage, ob der Server neben ihr auch eine
// unverschlüsselte Verbindung annimmt.
type TLS struct {
	cfg      *tls.Config
	klartext bool
}

// NeuesTLS baut die Verschlüsselung aus den PEM-Bytes der Zertifikatsdatei
// (Zertifikat des Servers, danach die Kette) und der Schlüsseldatei. Es prüft
// weder Gültigkeit noch Namen des Zertifikats: Das Laden und Prüfen der Dateien
// liegt beim Start (LH-FA-23.a *Dateien*). Mindestversion ist TLS 1.2; der
// Server fordert kein Zertifikat des Clients an und handelt kein ALPN aus. Ohne
// Zertifikat und Schlüssel (kein --tls-cert) liefert es nil: keine Verschlüsselung.
func NeuesTLS(zertifikat, schluessel []byte, klartextErlaubt bool) (*TLS, error) {
	if len(zertifikat) == 0 && len(schluessel) == 0 {
		return nil, nil
	}
	paar, err := tls.X509KeyPair(zertifikat, schluessel)
	if err != nil {
		return nil, model.Errorf(model.CodeConfigCA, nil, "TLS zum Client: Zertifikat und Schlüssel nicht zu laden")
	}
	return &TLS{
		cfg:      &tls.Config{Certificates: []tls.Certificate{paar}, MinVersion: tls.VersionTLS12},
		klartext: klartextErlaubt,
	}, nil
}

// MitTLS schaltet die Verschlüsselung zum Client ein; es steht vor Serve.
func (s *Server) MitTLS(t *TLS) {
	s.tls = t
}

// kanal ist die Verbindung eines Clients im Verbindungsaufbau: die Verbindung,
// ihr Leser und der Backend-Leser darüber; nach dem Handshake sind es die der
// verschlüsselten Verbindung.
type kanal struct {
	conn           net.Conn
	br             *bufio.Reader
	be             *pgproto3.Backend
	verschluesselt bool
}

func neuerKanal(conn net.Conn) *kanal {
	br := bufio.NewReader(conn)
	return &kanal{conn: conn, br: br, be: pgproto3.NewBackend(br, conn)}
}

// sonderanfrage beantwortet einen gelesenen SSLRequest (sslAnfrage) oder
// GSSENCRequest und meldet, ob der Aufbau weitergeht. Auf einer verschlüsselten
// Verbindung ist ein GSSENCRequest PGR-E6001, sonst die Antwort „N“ (LH-FA-05.c).
func (s *Server) sonderanfrage(ctx context.Context, k *kanal, code uint32) bool {
	if code == codeSSLRequest {
		return s.sslAnfrage(ctx, k)
	}
	if k.verschluesselt {
		s.fail(k.be, model.Errorf(model.CodeUnsupported, nil, "GSSENCRequest auf einer verschlüsselten Verbindung wird nicht unterstützt"))
		return false
	}
	return antworteN(k.conn)
}

// sslAnfrage beantwortet einen SSLRequest und meldet, ob der Aufbau
// weitergeht. Ohne TLS-Konfiguration ist die Antwort „N“ (LH-FA-05.c). Mit ihr
// führt es den Handshake aus, mindestens TLS 1.2; Bytes, die der Client nach
// dem SSLRequest vor der Antwort sendet, beenden die Verbindung mit PGR-W3002,
// bevor ein Handshake beginnt. Ein gescheiterter Handshake ist die Warnung
// PGR-W3002 und ändert den Prozessausgang nicht; Herunterfahren und Zwangsende
// beenden ihn ohne Meldung. Auf einer verschlüsselten Verbindung ist ein
// weiterer SSLRequest PGR-E6001 (LH-FA-23.a *Aushandlung*).
func (s *Server) sslAnfrage(ctx context.Context, k *kanal) bool {
	if s.tls == nil {
		return antworteN(k.conn)
	}
	if k.verschluesselt {
		s.fail(k.be, model.Errorf(model.CodeUnsupported, nil, "SSLRequest auf einer verschlüsselten Verbindung wird nicht unterstützt"))
		return false
	}
	if k.br.Buffered() > 0 {
		s.handshakeFehler(k.conn, errors.New("Bytes nach dem SSLRequest vor der Antwort"))
		return false
	}
	if _, err := k.conn.Write([]byte{'S'}); err != nil {
		return false
	}
	tc := tls.Server(k.conn, s.tls.cfg)
	if err := tc.HandshakeContext(ctx); err != nil {
		if ctx.Err() != nil || s.zwang.Load() {
			s.log.Debug("TLS-Aushandlung beim Herunterfahren beendet", "grund", err.Error())
			return false
		}
		s.handshakeFehler(k.conn, err)
		return false
	}
	k.conn = tlsVerbindung{tc}
	k.br = bufio.NewReader(k.conn)
	k.be = pgproto3.NewBackend(k.br, k.conn)
	k.verschluesselt = true
	return true
}

// tlsVerbindung ist eine verschlüsselte Verbindung, die beim Schließen kein
// close_notify sendet: Der Alert könnte bis zu 5 s auf einen Client warten, der
// nicht liest, und das Zwangsende wartet nicht auf den Client (LH-FA-23.a *Nach
// der Aushandlung*).
type tlsVerbindung struct{ *tls.Conn }

// Close schließt die Verbindung darunter.
func (c tlsVerbindung) Close() error { return c.NetConn().Close() }

// antworteN lehnt eine SSL- oder GSS-Anfrage mit „N“ ab und meldet, ob der
// Aufbau weitergeht.
func antworteN(conn net.Conn) bool {
	_, err := conn.Write([]byte{'N'})
	return err == nil
}

// handshakeFehler schreibt die Warnung PGR-W3002 mit der Adresse des Clients
// und der Ursache des Handshakes; Zertifikat und Schlüssel nennt sie nie.
func (s *Server) handshakeFehler(conn net.Conn, ursache error) {
	s.log.Warn("TLS-Aushandlung gescheitert", "code", model.CodeHandshake, "remote", conn.RemoteAddr().String(), "error", ursache.Error())
}
