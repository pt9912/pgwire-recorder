package tlsproxy

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// sslRequestNummer ist die Nummer des SSLRequest im Startpaket.
const sslRequestNummer = 80877103

// Proxy ist ein TLS-Proxy vor einer PostgreSQL-Instanz (der gepinnte Server
// enthält kein openssl und läuft ohne ssl): Auf ein SSLRequest antwortet er
// `S`, handelt TLS mit seinem Zertifikat aus und reicht die entschlüsselte
// Verbindung an die Instanz weiter. Eine Verbindung ohne SSLRequest lehnt er
// wie ein Server mit hostssl ab (Fehlerantwort 28000). Hängt er, antwortet er
// `S` und schweigt danach, sodass die Aushandlung des Clients wartet. Die
// Zähler nennen die Verbindungen, die TLS aushandelten (TLSOk), die ohne
// SSLRequest kamen (Klartext) und die an die Instanz weitergereicht wurden
// (Weitergabe).
type Proxy struct {
	// Addr ist die Adresse, auf der der Proxy lauscht.
	Addr string
	// Hello wird geschlossen, sobald ein hängender Proxy das erste Byte nach
	// seiner Antwort `S` empfängt.
	Hello chan struct{}

	ziel                        string
	haengt                      bool
	blatt                       tls.Certificate
	tlsOk, klartext, weitergabe atomic.Int32
}

// TLSOk ist die Zahl der Verbindungen, die TLS aushandelten.
func (p *Proxy) TLSOk() int { return int(p.tlsOk.Load()) }

// Klartext ist die Zahl der Verbindungen ohne SSLRequest.
func (p *Proxy) Klartext() int { return int(p.klartext.Load()) }

// Weitergabe ist die Zahl der Verbindungen, die an die Instanz gingen.
func (p *Proxy) Weitergabe() int { return int(p.weitergabe.Load()) }

// Start lauscht auf 127.0.0.1 und reicht entschlüsselte Verbindungen an ziel
// (host:port der Instanz) weiter; blatt ist sein Zertifikat, haengt lässt ihn
// nach `S` schweigen. Der Cleanup des Tests schließt den Proxy.
func Start(t testing.TB, blatt testpki.Blattzertifikat, ziel string, haengt bool) *Proxy {
	t.Helper()
	l, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	p := &Proxy{
		Addr: l.Addr().String(), Hello: make(chan struct{}), ziel: ziel, haengt: haengt,
		blatt: tls.Certificate{Certificate: [][]byte{blatt.DER}, PrivateKey: blatt.Schluessel},
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go p.bediene(conn)
		}
	}()
	return p
}

// fehlerantwort28000 ist die Fehlerantwort FATAL 28000 eines Servers, der eine
// Verbindung ohne Verschlüsselung ablehnt.
func fehlerantwort28000() []byte {
	var rumpf []byte
	for _, f := range [][2]string{{"S", "FATAL"}, {"V", "FATAL"}, {"C", "28000"}, {"M", "pg_hba.conf rejects connection, no encryption"}} {
		rumpf = append(append(append(rumpf, f[0][0]), f[1]...), 0)
	}
	rumpf = append(rumpf, 0)
	kopf := []byte{'E', 0, 0, 0, 0}
	binary.BigEndian.PutUint32(kopf[1:], uint32(len(rumpf)+4))
	return append(kopf, rumpf...)
}

// bediene führt eine Verbindung des Clients.
func (p *Proxy) bediene(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(60 * time.Second))
	kopf := make([]byte, 8)
	if _, err := io.ReadFull(conn, kopf); err != nil {
		return
	}
	if binary.BigEndian.Uint32(kopf[4:]) != sslRequestNummer {
		p.klartext.Add(1)
		_, _ = conn.Write(fehlerantwort28000())
		return
	}
	_, _ = conn.Write([]byte("S"))
	if p.haengt {
		if _, err := io.ReadFull(conn, make([]byte, 1)); err == nil {
			close(p.Hello)
		}
		_, _ = io.Copy(io.Discard, conn)
		return
	}
	tc := tls.Server(conn, &tls.Config{Certificates: []tls.Certificate{p.blatt}})
	if tc.HandshakeContext(context.Background()) != nil {
		return
	}
	p.tlsOk.Add(1)
	backend, err := (&net.Dialer{}).DialContext(context.Background(), "tcp", p.ziel)
	if err != nil {
		return
	}
	defer backend.Close()
	p.weitergabe.Add(1)
	go func() {
		_, _ = io.Copy(backend, tc)
		_ = backend.Close()
	}()
	_, _ = io.Copy(tc, backend)
}
