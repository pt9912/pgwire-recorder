package tlsproxy_test

import (
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap/tlsproxy"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// TestHaengenderProxyZweiteVerbindung prüft, dass ein hängender Proxy auch das
// erste Byte nach `S` einer zweiten Verbindung annimmt, ohne Hello noch einmal
// zu schließen: Jede Verbindung sendet SSLRequest und ein Byte und schließt
// ihre Schreibrichtung; das Ende des Servers belegt, dass sein Handler hinter
// dem Schließen von Hello angekommen ist.
func TestHaengenderProxyZweiteVerbindung(t *testing.T) {
	ca := testpki.NeueCA(t, "Proxy-CA")
	blatt := ca.Server(t, testpki.Blatt{IPs: []net.IP{net.ParseIP("127.0.0.1")}})
	proxy := tlsproxy.Start(t, blatt, "127.0.0.1:1", true)
	for i := range 2 {
		conn, err := net.DialTimeout("tcp", proxy.Addr, 10*time.Second)
		if err != nil {
			t.Fatal(err)
		}
		_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
		kopf := make([]byte, 8)
		binary.BigEndian.PutUint32(kopf, 8)
		binary.BigEndian.PutUint32(kopf[4:], 80877103)
		if _, err := conn.Write(kopf); err != nil {
			t.Fatalf("Verbindung %d: %v", i+1, err)
		}
		antwort := make([]byte, 1)
		if _, err := conn.Read(antwort); err != nil || antwort[0] != 'S' {
			t.Fatalf("Verbindung %d: Antwort %q, %v", i+1, antwort, err)
		}
		if _, err := conn.Write([]byte{0x16}); err != nil {
			t.Fatalf("Verbindung %d: %v", i+1, err)
		}
		_ = conn.(*net.TCPConn).CloseWrite()
		if n, err := conn.Read(make([]byte, 1)); n != 0 || err == nil {
			t.Fatalf("Verbindung %d: der Proxy sendet nach S noch %d Bytes", i+1, n)
		}
		conn.Close()
		select {
		case <-proxy.Hello:
		default:
			t.Fatalf("Hello nach der Verbindung %d nicht geschlossen", i+1)
		}
	}
}
