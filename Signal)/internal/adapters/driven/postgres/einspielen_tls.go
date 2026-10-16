package postgres

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Zertifikate sind die Zertifikate einer zusätzlichen Zertifizierungsstelle
// (LH-FA-20.a *Start*). Jede Formatierung gibt festen Text aus, nie ein
// Zertifikat: Ein Aufruf mit %v, %+v, %#v, %s oder %q auf einem Zertifikat
// oder einem Einspielziel verrät keines.
type Zertifikate []*x509.Certificate

// Format gibt für jedes Verb den festen Text aus.
func (z Zertifikate) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprintf(f, "Zertifikate{%d}", len(z))
}

// aushandeln sendet das SSLRequest und handelt TLS aus (LH-FA-20.a *TLS*):
// `S` auf das SSLRequest ist der Beginn der Aushandlung, `N` ist PGR-E4005;
// jedes andere Byte, ein Verbindungsende davor und Bytes, die der Server nach
// `S` vor der Aushandlung sendet, sind PGR-E4002. Jeder Fehler der Aushandlung
// selbst, auch ein Verbindungsende darin, ist PGR-E4005. Die Antwort liest es
// mit einem Lesen; Bytes, die darin auf das erste folgen, sind die Bytes nach
// `S` (Grenze: Bytes, die erst nach diesem Lesen eintreffen, erkennt es nicht).
// Endet ctx in der Aushandlung, bricht sie ab und schließt conn. Es liefert die
// verschlüsselte Verbindung über conn; bei einem Fehler schließt der Aufrufer
// conn, ohne dass danach etwas gesendet wird.
func (z *Einspielziel) aushandeln(ctx context.Context, conn net.Conn) (net.Conn, error) {
	anfrage, err := (&pgproto3.SSLRequest{}).Encode(nil)
	if err != nil {
		return nil, model.Errorf(model.CodeUpstream, err, "SSLRequest nicht zu kodieren")
	}
	if _, err := conn.Write(anfrage); err != nil {
		return nil, model.Errorf(model.CodeUpstream, err, "SSLRequest nicht zu senden")
	}
	var antwort [512]byte
	n, err := io.ReadAtLeast(conn, antwort[:], 1)
	if err != nil {
		return nil, model.Errorf(model.CodeUpstream, err, "Verbindung vor der Antwort auf das SSLRequest beendet")
	}
	switch antwort[0] {
	case 'S':
	case 'N':
		return nil, model.Errorf(model.CodeLogin, nil, "Server %s lehnt TLS ab", z.Address)
	default:
		return nil, model.Errorf(model.CodeUpstream, nil, "Antwort %q auf das SSLRequest ist weder S noch N", antwort[0])
	}
	if n > 1 {
		return nil, model.Errorf(model.CodeUpstream, nil, "Bytes des Servers nach S vor der Aushandlung")
	}
	verschluesselt := tls.Client(conn, z.tlsKonfiguration())
	if err := verschluesselt.HandshakeContext(ctx); err != nil {
		return nil, model.Errorf(model.CodeLogin, err, "TLS mit %s nicht auszuhandeln", z.Address)
	}
	return verschluesselt, nil
}

// tlsKonfiguration ist die Konfiguration des Clients: Das Serverzertifikat
// wird gegen den Zertifikatsspeicher des Systems und die Zertifikate aus CA
// geprüft, der Name gegen den Host der Adresse (serverName); ein Überspringen
// der Prüfung gibt es nicht (LH-FA-20.a *Schritte* 2, *TLS*).
func (z *Einspielziel) tlsKonfiguration() *tls.Config {
	wurzeln := z.systemspeicher
	if wurzeln == nil {
		wurzeln = x509.SystemCertPool
	}
	pool, err := wurzeln()
	if err != nil || pool == nil {
		pool = x509.NewCertPool()
	}
	for _, c := range z.CA {
		pool.AddCert(c)
	}
	return &tls.Config{RootCAs: pool, ServerName: serverName(z.Address)}
}

// serverName ist der Name, gegen den das Zertifikat geprüft wird: der Host von
// address, wie er eingesetzt ist. Die Standardbibliothek prüft eine
// IP-Adresse (IPv4 oder IPv6, eine IPv6-Adresse ohne ihre Zone) nur gegen die
// IP-Adressen des Zertifikats, nie gegen die DNS-Namen.
func serverName(address string) string {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return address
	}
	return host
}
