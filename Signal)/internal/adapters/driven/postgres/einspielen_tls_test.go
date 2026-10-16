package postgres_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// sslRequest sind die acht Bytes des SSLRequest: Länge 8 und die Nummer 80877103.
func sslRequest() []byte { return []byte{0, 0, 0, 8, 4, 210, 22, 47} }

// tlsLauf ist, was die Gegenstelle von TLS erlebte: die ersten acht Bytes des
// Clients, die Startup-Parameter und alle Bytes danach über die
// verschlüsselte Verbindung, das Ende des Lesens (nil bei io.EOF) und der
// Fehler ihrer Aushandlung.
type tlsLauf struct {
	anfrage     []byte
	startup     map[string]string
	danach      []byte
	ende        error
	aushandlung error
}

// tlsStelle ist die Test-Gegenstelle: Sie liest die acht Bytes des SSLRequest
// und sendet antwort. Ohne config liest sie danach bis zum Ende der Verbindung
// (nachHello, wenn gesetzt, sendet sie nach dem ersten Lesen und schließt;
// schliesse schließt sofort nach antwort). Mit config handelt sie TLS aus, liest
// das Startup über die verschlüsselte Verbindung und führt verlauf aus; ohne
// verlauf sendet sie einen Aufbau ohne Anmeldung und liest bis zum Ende.
type tlsStelle struct {
	antwort   []byte
	config    *tls.Config
	schliesse bool
	nachHello []byte
	// hello wird geschlossen, sobald das erste Byte nach der Antwort eintrifft
	// (ohne config).
	hello   chan struct{}
	verlauf func(c net.Conn, l *tlsLauf)
}

// liesStartup liest die Startup-Nachricht von c.
func liesStartup(c net.Conn) (map[string]string, error) {
	var laenge [4]byte
	if _, err := io.ReadFull(c, laenge[:]); err != nil {
		return nil, err
	}
	rumpf := make([]byte, binary.BigEndian.Uint32(laenge[:])-4)
	if _, err := io.ReadFull(c, rumpf); err != nil {
		return nil, err
	}
	var sm pgproto3.StartupMessage
	if err := sm.Decode(rumpf); err != nil {
		return nil, err
	}
	return sm.Parameters, nil
}

// starte lauscht auf 127.0.0.1, nimmt eine Verbindung an und liefert die
// Adresse und den Lauf; die Gegenstelle gibt nach 20 s auf.
func (g tlsStelle) starte(t *testing.T, aufbau []byte) (string, <-chan tlsLauf) {
	t.Helper()
	return g.starteAuf(t, "127.0.0.1:0", aufbau)
}

// starteAuf ist starte auf address.
func (g tlsStelle) starteAuf(t *testing.T, address string, aufbau []byte) (string, <-chan tlsLauf) {
	t.Helper()
	l, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	ergebnis := make(chan tlsLauf, 1)
	go func() {
		var lauf tlsLauf
		defer func() { ergebnis <- lauf }()
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_ = conn.SetDeadline(time.Now().Add(20 * time.Second))
		lauf.anfrage = make([]byte, 8)
		if _, err := io.ReadFull(conn, lauf.anfrage); err != nil {
			lauf.ende = err
			return
		}
		if len(g.antwort) > 0 {
			_, _ = conn.Write(g.antwort)
		}
		switch {
		case g.schliesse:
			return
		case g.nachHello != nil:
			_, _ = conn.Read(make([]byte, 5))
			_, _ = conn.Write(g.nachHello)
			return
		case g.config == nil:
			if g.hello != nil {
				erstes := make([]byte, 1)
				if _, lauf.ende = io.ReadFull(conn, erstes); lauf.ende != nil {
					return
				}
				lauf.danach = erstes
				close(g.hello)
			}
			rest, ende := io.ReadAll(conn)
			lauf.danach, lauf.ende = append(lauf.danach, rest...), ende
			return
		}
		tc := tls.Server(conn, g.config)
		if err := tc.Handshake(); err != nil {
			lauf.aushandlung = err
			return
		}
		if lauf.startup, lauf.ende = liesStartup(tc); lauf.ende != nil {
			return
		}
		if g.verlauf != nil {
			g.verlauf(tc, &lauf)
			return
		}
		_, _ = tc.Write(aufbau)
		lauf.danach, lauf.ende = io.ReadAll(tc)
	}()
	return l.Addr().String(), ergebnis
}

// serverConfig ist die Konfiguration eines TLS-Servers mit dem Zertifikat blatt.
func serverConfig(blatt testpki.Blattzertifikat) *tls.Config {
	return &tls.Config{Certificates: []tls.Certificate{{Certificate: [][]byte{blatt.DER}, PrivateKey: blatt.Schluessel}}}
}

// lokal ist die Adresse 127.0.0.1 als IP-Adresse im Zertifikat.
func lokal() []net.IP { return []net.IP{net.ParseIP("127.0.0.1")} }

// tlsErgebnis wartet höchstens 30 s auf den Lauf der Gegenstelle.
func tlsErgebnis(t *testing.T, ergebnis <-chan tlsLauf) tlsLauf {
	t.Helper()
	select {
	case l := <-ergebnis:
		return l
	case <-time.After(30 * time.Second):
		t.Fatal("die Gegenstelle meldet binnen 30 s kein Ende")
		return tlsLauf{}
	}
}

// tlsVerbinde verbindet z mit startup und liefert Session und Fehler.
func tlsVerbinde(t *testing.T, z *postgres.Einspielziel) (any, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s, err := z.Verbinde(ctx, map[string]string{"user": "u"})
	if err != nil {
		return nil, err
	}
	return s, nil
}

// leerSpeicher ist ein Zertifikatsspeicher des Systems ohne Zertifikat.
func leerSpeicher() (*x509.CertPool, error) { return x509.NewCertPool(), nil }

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary — mit TLS sendet play vor dem
// Startup das SSLRequest (acht Bytes), handelt nach `S` TLS aus, prüft das
// Zertifikat gegen die Zertifikate aus CA und sendet Startup, Anmeldung und
// Terminate nur über die verschlüsselte Verbindung (LH-FA-20.a *TLS*, *Schritte*).
func TestEinspielTLS(t *testing.T) {
	ca := testpki.NeueCA(t, "Test-CA")
	g := tlsStelle{antwort: []byte("S"), config: serverConfig(ca.Server(t, testpki.Blatt{IPs: lokal()}))}
	addr, ergebnis := g.starte(t, verbunden(t))
	z := postgres.MitSystemspeicher(&postgres.Einspielziel{Address: addr, TLS: true, CA: postgres.Zertifikate{ca.Zertifikat}}, leerSpeicher)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s, err := z.Verbinde(ctx, map[string]string{"user": "u", "database": "d"})
	if err != nil {
		t.Fatal(err)
	}
	s.Schliesse()
	l := tlsErgebnis(t, ergebnis)
	if string(l.anfrage) != string(sslRequest()) {
		t.Errorf("erste Bytes %v, erwartet das SSLRequest %v", l.anfrage, sslRequest())
	}
	if l.startup["user"] != "u" || l.startup["database"] != "d" || l.ende != nil {
		t.Errorf("Startup über TLS %v, Ende %v", l.startup, l.ende)
	}
	if string(l.danach) != "X\x00\x00\x00\x04" {
		t.Errorf("nach dem Startup über TLS %q, erwartet Terminate", l.danach)
	}
}

// Abdeckung: LH-FA-20/Boundary — ohne TLS sendet play kein SSLRequest und
// keine TLS-Nachricht, auch wenn Zertifikate aus CA gesetzt sind: Die ersten
// Bytes sind das Startup (LH-FA-20.a *Schritte*, *TLS*).
func TestEinspielOhneTLS(t *testing.T) {
	ca := testpki.NeueCA(t, "Test-CA")
	addr, ergebnis := einspielServer(t, verbunden(t))
	s, err := (&postgres.Einspielziel{Address: addr, CA: postgres.Zertifikate{ca.Zertifikat}}).Verbinde(context.Background(), map[string]string{"user": "u"})
	if err != nil {
		t.Fatal(err)
	}
	s.Schliesse()
	if l := empfangen(t, ergebnis); l.startup["user"] != "u" || string(l.danach) != "X\x00\x00\x00\x04" {
		t.Errorf("Startup %v, danach %q", l.startup, l.danach)
	}
}

// Abdeckung: LH-FA-20/Negative — ein Zertifikat, das abgelaufen
// ist, auf einen anderen Namen ausgestellt ist, von einer Zertifizierungsstelle
// stammt, die nicht im Speicher und nicht in CA steht, und eines, das den Host
// nur als DNS-Namen oder als andere IP-Adresse nennt, ist PGR-E4005; der Grund
// steht in der Meldung, und es folgt weder Startup noch Terminate (LH-FA-20.a
// *TLS*, *Abbruch im Aufbau*).
func TestEinspielTLSZertifikatFehler(t *testing.T) {
	ca, andere := testpki.NeueCA(t, "Test-CA"), testpki.NeueCA(t, "Andere-CA")
	abgelaufen := testpki.Blatt{IPs: lokal(), Von: time.Now().Add(-2 * time.Hour), Bis: time.Now().Add(-time.Hour)}
	for _, f := range []struct {
		name       string
		blatt      testpki.Blatt
		aussteller *testpki.CA
		grund      string
	}{
		{"abgelaufen", abgelaufen, ca, "expired"},
		{"noch nicht gültig", testpki.Blatt{IPs: lokal(), Von: time.Now().Add(time.Hour), Bis: time.Now().Add(2 * time.Hour)}, ca, "expired or is not yet valid"},
		{"andere Zertifizierungsstelle", testpki.Blatt{IPs: lokal()}, andere, "unknown authority"},
		{"andere IP-Adresse", testpki.Blatt{IPs: []net.IP{net.ParseIP("10.1.2.3")}}, ca, "not 127.0.0.1"},
		{"nur DNS-Name localhost", testpki.Blatt{DNS: []string{"localhost"}}, ca, "any IP SANs"},
		{"IP-Adresse als DNS-Name", testpki.Blatt{DNS: []string{"127.0.0.1"}}, ca, "any IP SANs"},
	} {
		t.Run(f.name, func(t *testing.T) {
			g := tlsStelle{antwort: []byte("S"), config: serverConfig(f.aussteller.Server(t, f.blatt))}
			addr, ergebnis := g.starte(t, verbunden(t))
			z := postgres.MitSystemspeicher(&postgres.Einspielziel{Address: addr, TLS: true, CA: postgres.Zertifikate{ca.Zertifikat}}, leerSpeicher)
			s, err := tlsVerbinde(t, z)
			if s != nil || code(err) != model.CodeLogin || !strings.Contains(err.Error(), f.grund) {
				t.Fatalf("Session %v, Fehler %v, erwartet %s mit %q", s, err, model.CodeLogin, f.grund)
			}
			if l := tlsErgebnis(t, ergebnis); l.aushandlung == nil || l.startup != nil || len(l.danach) != 0 {
				t.Fatalf("die Gegenstelle sah Aushandlung %v, Startup %v, danach %q", l.aushandlung, l.startup, l.danach)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary — der Name im Zertifikat wird
// gegen den Host der Adresse geprüft, wie er eingesetzt ist: ein DNS-Name
// gegen die DNS-Namen, eine IPv4-Adresse gegen die IP-Adressen; ein
// DNS-Name im Zertifikat trifft keinen Host, der eine IP-Adresse ist, und
// eine IP-Adresse im Zertifikat keinen Host, der ein Name ist (LH-FA-20.a
// *TLS*).
func TestEinspielTLSName(t *testing.T) {
	if _, err := net.LookupHost("localhost"); err != nil {
		t.Skipf("localhost löst nicht auf: %v", err)
	}
	ca := testpki.NeueCA(t, "Test-CA")
	for _, f := range []struct {
		name  string
		blatt testpki.Blatt
		host  string
		ok    bool
	}{
		{"IPv4 gegen IP-Adresse", testpki.Blatt{IPs: lokal()}, "127.0.0.1", true},
		{"Name gegen DNS-Namen", testpki.Blatt{DNS: []string{"localhost"}}, "localhost", true},
		{"Name gegen anderen DNS-Namen", testpki.Blatt{DNS: []string{"anders.example"}}, "localhost", false},
		{"Name gegen IP-Adresse", testpki.Blatt{IPs: lokal()}, "localhost", false},
	} {
		t.Run(f.name, func(t *testing.T) {
			g := tlsStelle{antwort: []byte("S"), config: serverConfig(ca.Server(t, f.blatt))}
			addr, ergebnis := g.starte(t, verbunden(t))
			_, port, _ := net.SplitHostPort(addr)
			z := postgres.MitSystemspeicher(&postgres.Einspielziel{Address: net.JoinHostPort(f.host, port), TLS: true, CA: postgres.Zertifikate{ca.Zertifikat}}, leerSpeicher)
			s, err := tlsVerbinde(t, z)
			if f.ok {
				if err != nil {
					t.Fatalf("Fehler %v", err)
				}
				s.(interface{ Schliesse() }).Schliesse()
				tlsErgebnis(t, ergebnis)
				return
			}
			if s != nil || code(err) != model.CodeLogin {
				t.Fatalf("Session %v, Fehler %v, erwartet %s", s, err, model.CodeLogin)
			}
			tlsErgebnis(t, ergebnis)
		})
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary — der Name, gegen den das
// Zertifikat geprüft wird, ist der Host der Adresse, wie er eingesetzt ist
// (eine IPv6-Adresse ohne Klammern, auch mit Zone); die Standardbibliothek
// prüft eine IPv6-Adresse mit und ohne Zone gegen die IP-Adressen, nie gegen
// DNS-Namen (LH-FA-20.a *TLS*).
func TestEinspielTLSServerName(t *testing.T) {
	ca := testpki.NeueCA(t, "Test-CA")
	cert, err := x509.ParseCertificate(ca.Server(t, testpki.Blatt{IPs: []net.IP{net.ParseIP("::1")}, DNS: []string{"fe80::1", "db.example"}}).DER)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range []struct {
		adresse, name string
		passt         bool
	}{
		{"[::1]:5432", "::1", true},
		{"[::1%eth0]:5432", "::1%eth0", true},
		{"[fe80::1%eth0]:5432", "fe80::1%eth0", false},
		{"127.0.0.1:5432", "127.0.0.1", false},
		{"db.example:5432", "db.example", true},
		{"ohne-port", "ohne-port", false},
	} {
		got := postgres.ServerName(f.adresse)
		if got != f.name {
			t.Errorf("ServerName(%q) = %q, erwartet %q", f.adresse, got, f.name)
		}
		if err := cert.VerifyHostname(got); (err == nil) != f.passt {
			t.Errorf("Zertifikat auf ::1 und die DNS-Namen fe80::1 und db.example gegen %q: %v, erwartet passt = %v", got, err, f.passt)
		}
	}
}

// Abdeckung: LH-FA-20/Happy — eine IPv6-Adresse mit Zone in der Adresse wird
// gegen die IP-Adressen des Zertifikats geprüft; der Test braucht IPv6 auf der
// Schleife und endet sonst ohne Prüfung (LH-FA-20.a *TLS*).
func TestEinspielTLSIPv6MitZone(t *testing.T) {
	ca := testpki.NeueCA(t, "Test-CA")
	g := tlsStelle{antwort: []byte("S"), config: serverConfig(ca.Server(t, testpki.Blatt{IPs: []net.IP{net.ParseIP("::1")}}))}
	l, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Skipf("kein IPv6 auf der Schleife: %v", err)
	}
	l.Close()
	addr, ergebnis := g.starteAuf(t, "[::1]:0", verbunden(t))
	_, port, _ := net.SplitHostPort(addr)
	z := postgres.MitSystemspeicher(&postgres.Einspielziel{Address: "[::1%lo]:" + port, TLS: true, CA: postgres.Zertifikate{ca.Zertifikat}}, leerSpeicher)
	s, err := tlsVerbinde(t, z)
	if err != nil {
		t.Fatalf("Fehler %v", err)
	}
	s.(interface{ Schliesse() }).Schliesse()
	tlsErgebnis(t, ergebnis)
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary, LH-FA-20/Negative — der
// Zertifikatsspeicher des Systems und die Zertifikate aus CA gelten
// zusammen: Ein Server, dessen Zertifizierungsstelle nur im Speicher steht,
// nur in CA steht oder in beiden, wird angenommen; steht sie in keinem, ist es
// PGR-E4005, auch wenn der Speicher leer ist, nicht ladbar oder fehlt; CA
// ergänzt den Speicher und ersetzt ihn nicht (LH-FA-20.a *TLS*).
func TestEinspielTLSZertifikatsspeicher(t *testing.T) {
	ca, andere := testpki.NeueCA(t, "Test-CA"), testpki.NeueCA(t, "Andere-CA")
	speicher := func(cas ...*testpki.CA) func() (*x509.CertPool, error) {
		return func() (*x509.CertPool, error) {
			p := x509.NewCertPool()
			for _, c := range cas {
				p.AddCert(c.Zertifikat)
			}
			return p, nil
		}
	}
	for _, f := range []struct {
		name     string
		speicher func() (*x509.CertPool, error)
		eigene   postgres.Zertifikate
		ok       bool
	}{
		{"nur im Speicher", speicher(ca), nil, true},
		{"nur in CA", leerSpeicher, postgres.Zertifikate{ca.Zertifikat}, true},
		{"im Speicher, andere in CA", speicher(ca), postgres.Zertifikate{andere.Zertifikat}, true},
		{"andere im Speicher, in CA", speicher(andere), postgres.Zertifikate{ca.Zertifikat}, true},
		{"in beiden", speicher(ca), postgres.Zertifikate{ca.Zertifikat}, true},
		{"im Speicher neben anderen, eine eigene", speicher(andere, ca), postgres.Zertifikate{andere.Zertifikat, andere.Zertifikat}, true},
		{"leerer Speicher", leerSpeicher, nil, false},
		{"leerer Speicher, andere in CA", leerSpeicher, postgres.Zertifikate{andere.Zertifikat}, false},
		{"nicht ladbarer Speicher", func() (*x509.CertPool, error) { return nil, errors.New("nicht ladbar") }, nil, false},
		{"Speicher nil ohne Fehler", func() (*x509.CertPool, error) { return nil, nil }, nil, false},
		{"nicht ladbarer Speicher, in CA", func() (*x509.CertPool, error) { return nil, errors.New("nicht ladbar") }, postgres.Zertifikate{ca.Zertifikat}, true},
		{"Speicher nil ohne Fehler, in CA", func() (*x509.CertPool, error) { return nil, nil }, postgres.Zertifikate{ca.Zertifikat}, true},
	} {
		t.Run(f.name, func(t *testing.T) {
			g := tlsStelle{antwort: []byte("S"), config: serverConfig(ca.Server(t, testpki.Blatt{IPs: lokal()}))}
			addr, ergebnis := g.starte(t, verbunden(t))
			s, err := tlsVerbinde(t, postgres.MitSystemspeicher(&postgres.Einspielziel{Address: addr, TLS: true, CA: f.eigene}, f.speicher))
			if f.ok != (err == nil) || (err != nil && code(err) != model.CodeLogin) {
				t.Fatalf("Fehler %v, erwartet Erfolg %v", err, f.ok)
			}
			if err == nil {
				s.(interface{ Schliesse() }).Schliesse()
			}
			tlsErgebnis(t, ergebnis)
		})
	}
}

// Abdeckung: LH-FA-20/Negative — auf das SSLRequest ist `N`
// PGR-E4005, jedes andere Byte (auch ein Fehlerpaket, ein Byte eines
// TLS-Records und ein Byte ähnlich `S`) und ein Verbindungsende davor
// PGR-E4002, ebenso Bytes, die der Server nach `S` vor der Aushandlung sendet;
// jeder Fehler der Aushandlung, auch ein Verbindungsende darin, Daten, die
// kein TLS sind, und eine Gegenstelle, die nur TLS 1.0 spricht, sodass die
// Aushandlung scheitert, ist PGR-E4005; danach sendet play nichts mehr
// (LH-FA-20.a *TLS*, *Abbruch im Aufbau*).
func TestEinspielTLSAntwort(t *testing.T) {
	ca := testpki.NeueCA(t, "Test-CA")
	nurTLS10 := serverConfig(ca.Server(t, testpki.Blatt{IPs: lokal()}))
	nurTLS10.MaxVersion = tls.VersionTLS10
	for _, f := range []struct {
		name string
		g    tlsStelle
		code string
	}{
		{"N", tlsStelle{antwort: []byte("N")}, model.CodeLogin},
		{"N mit Bytes dahinter", tlsStelle{antwort: []byte("Nxyz")}, model.CodeLogin},
		{"E", tlsStelle{antwort: []byte("E")}, model.CodeUpstream},
		{"Fehlerpaket", tlsStelle{antwort: kodiert(t, &pgproto3.ErrorResponse{Severity: "FATAL", Code: "08P01", Message: "nein"})}, model.CodeUpstream},
		{"x", tlsStelle{antwort: []byte("x")}, model.CodeUpstream},
		{"s", tlsStelle{antwort: []byte("s")}, model.CodeUpstream},
		{"Byte eines TLS-Records", tlsStelle{antwort: []byte{0x16}}, model.CodeUpstream},
		{"Byte 0", tlsStelle{antwort: []byte{0}}, model.CodeUpstream},
		{"Ende davor", tlsStelle{schliesse: true}, model.CodeUpstream},
		{"S mit einem Byte dahinter", tlsStelle{antwort: []byte("Sx")}, model.CodeUpstream},
		{"S mit TLS-Record dahinter", tlsStelle{antwort: append([]byte("S"), 0x16, 3, 3, 0, 0)}, model.CodeUpstream},
		{"S, dann Ende", tlsStelle{antwort: []byte("S"), schliesse: true}, model.CodeLogin},
		{"S, dann kein TLS", tlsStelle{antwort: []byte("S"), nachHello: []byte("HTTP/1.1 400 Bad Request\r\n\r\n")}, model.CodeLogin},
		{"S, dann Gegenstelle nur TLS 1.0, Aushandlung scheitert", tlsStelle{antwort: []byte("S"), config: nurTLS10}, model.CodeLogin},
	} {
		t.Run(f.name, func(t *testing.T) {
			addr, ergebnis := f.g.starte(t, verbunden(t))
			s, err := tlsVerbinde(t, &postgres.Einspielziel{Address: addr, TLS: true, CA: postgres.Zertifikate{ca.Zertifikat}})
			if s != nil || code(err) != f.code {
				t.Fatalf("Session %v, Fehler %v, erwartet %s", s, err, f.code)
			}
			l := tlsErgebnis(t, ergebnis)
			if len(l.danach) != 0 || l.startup != nil {
				t.Fatalf("nach dem Fehler gesendet: Startup %v, Bytes %q", l.startup, l.danach)
			}
		})
	}
}

// stummeConn ist eine Verbindung, die aus lesen liest und deren Senden mit
// schreibFehler scheitert, wenn er gesetzt ist; sonst merkt sie das Gesendete.
type stummeConn struct {
	net.Conn
	lesen         io.Reader
	schreibFehler error
	gesendet      []byte
}

func (c *stummeConn) Read(p []byte) (int, error) { return c.lesen.Read(p) }

func (c *stummeConn) Write(p []byte) (int, error) {
	if c.schreibFehler != nil {
		return 0, c.schreibFehler
	}
	c.gesendet = append(c.gesendet, p...)
	return len(p), nil
}

// Abdeckung: LH-FA-20/Negative — scheitert das Senden des
// SSLRequest, ist das PGR-E4002 (ein gescheitertes Senden im Aufbau, nicht in
// der Aushandlung); ein Ende des Lesens mit einem anderen Fehler als dem
// Ende des Stroms ist es ebenso; die Bytes `S` mit Bytes dahinter in einem Lesen
// sind PGR-E4002, auch deterministisch ohne Netz (LH-FA-20.a *TLS*, *Aufbau*).
func TestEinspielTLSSendenUndLesen(t *testing.T) {
	z := &postgres.Einspielziel{Address: "h:1", TLS: true}
	_, err := postgres.Aushandeln(context.Background(), z, &stummeConn{schreibFehler: errors.New("kaputt")})
	if code(err) != model.CodeUpstream {
		t.Errorf("Senden scheitert: %v, erwartet %s", err, model.CodeUpstream)
	}
	for name, lesen := range map[string]io.Reader{
		"Fehler":       errReader{errors.New("zurückgesetzt")},
		"Ende":         strings.NewReader(""),
		"S und Byte":   strings.NewReader("Sx"),
		"anderes Byte": strings.NewReader("X"),
	} {
		c := &stummeConn{lesen: lesen}
		_, err := postgres.Aushandeln(context.Background(), z, c)
		if code(err) != model.CodeUpstream {
			t.Errorf("%s: %v, erwartet %s", name, err, model.CodeUpstream)
		}
		if string(c.gesendet) != string(sslRequest()) {
			t.Errorf("%s: gesendet %v, erwartet nur das SSLRequest", name, c.gesendet)
		}
	}
}

// errReader scheitert bei jedem Lesen.
type errReader struct{ err error }

func (r errReader) Read([]byte) (int, error) { return 0, r.err }

// Abdeckung: LH-FA-20/Negative — lehnt der Server eine
// unverschlüsselte Verbindung ab (Fehlerantwort der Klasse 28 auf das Startup,
// wie pg_hba.conf mit hostssl), ist das PGR-E4005 (LH-FA-20.a Tabelle
// *Fehlerregeln beim Einspielen*).
func TestEinspielKlartextAbgelehnt(t *testing.T) {
	aufbau := kodiert(t, &pgproto3.ErrorResponse{Severity: "FATAL", SeverityUnlocalized: "FATAL", Code: "28000", Message: "pg_hba.conf rejects connection for host, no encryption"})
	addr, ergebnis := einspielServer(t, aufbau)
	s, err := (&postgres.Einspielziel{Address: addr}).Verbinde(context.Background(), map[string]string{"user": "u"})
	if s != nil || code(err) != model.CodeLogin || !strings.Contains(err.Error(), "28000") {
		t.Fatalf("Session %v, Fehler %v, erwartet %s", s, err, model.CodeLogin)
	}
	if l := geschlossen(t, ergebnis); len(l.danach) != 0 {
		t.Fatalf("nach dem Startup gesendet: %q", l.danach)
	}
}

// Abdeckung: LH-FA-20/Negative — nach der Aushandlung gilt der Abbruch im
// Aufbau über die verschlüsselte Verbindung: Nach einer Fehlerantwort sendet
// play keine weitere Nachricht, auch kein Terminate; ein Alarm der TLS-Schicht
// ist keine Nachricht (LH-FA-20.a *Abbruch im Aufbau*).
func TestEinspielTLSAbbruchImAufbau(t *testing.T) {
	ca := testpki.NeueCA(t, "Test-CA")
	g := tlsStelle{antwort: []byte("S"), config: serverConfig(ca.Server(t, testpki.Blatt{IPs: lokal()}))}
	fehler := kodiert(t, &pgproto3.ErrorResponse{Severity: "FATAL", SeverityUnlocalized: "FATAL", Code: "28P01", Message: "nein"})
	addr, ergebnis := g.starte(t, fehler)
	s, err := tlsVerbinde(t, &postgres.Einspielziel{Address: addr, TLS: true, CA: postgres.Zertifikate{ca.Zertifikat}})
	if s != nil || code(err) != model.CodeLogin {
		t.Fatalf("Session %v, Fehler %v, erwartet %s", s, err, model.CodeLogin)
	}
	if l := tlsErgebnis(t, ergebnis); len(l.danach) != 0 || l.startup["user"] != "u" {
		t.Fatalf("Startup %v, nach der Fehlerantwort gesendet: %q", l.startup, l.danach)
	}
}

// Abdeckung: LH-FA-20/Happy — nach der Aushandlung laufen Aufbau und Anmeldung
// unverändert auf der verschlüsselten Verbindung: Verlangt der Server ein
// Klartext-Passwort, sendet play es über TLS (LH-FA-20.a *TLS*, *Anmeldung*).
func TestEinspielTLSKlartextPasswort(t *testing.T) {
	ca := testpki.NeueCA(t, "Test-CA")
	g := tlsStelle{antwort: []byte("S"), config: serverConfig(ca.Server(t, testpki.Blatt{IPs: lokal()})),
		verlauf: func(c net.Conn, l *tlsLauf) {
			_, _ = c.Write(kodiert(t, &pgproto3.AuthenticationCleartextPassword{}))
			var kopf [5]byte
			if _, l.ende = io.ReadFull(c, kopf[:]); l.ende != nil {
				return
			}
			rumpf := make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)
			if _, l.ende = io.ReadFull(c, rumpf); l.ende != nil {
				return
			}
			l.danach = append([]byte{kopf[0]}, rumpf...)
			_, _ = c.Write(verbunden(t))
			_, _ = io.Copy(io.Discard, c)
		}}
	addr, ergebnis := g.starte(t, nil)
	z := &postgres.Einspielziel{Address: addr, TLS: true, CA: postgres.Zertifikate{ca.Zertifikat}, Password: "GEHEIMpw"}
	s, err := tlsVerbinde(t, z)
	if err != nil {
		t.Fatal(err)
	}
	s.(interface{ Schliesse() }).Schliesse()
	if l := tlsErgebnis(t, ergebnis); string(l.danach) != "pGEHEIMpw\x00" {
		t.Fatalf("Passwortnachricht über TLS %q", l.danach)
	}
}

// Abdeckung: LH-FA-20/Boundary — das zweite Signal in der Aushandlung schließt
// die Verbindung: Verbinde endet mit einem Fehler, der den Abbruch nennt, ohne
// dass der Server die Aushandlung beendet (LH-FA-20.a *Abbruchsignal*).
func TestEinspielTLSAbbruchInDerAushandlung(t *testing.T) {
	ca := testpki.NeueCA(t, "Test-CA")
	hello := make(chan struct{})
	addr, ergebnis := tlsStelle{antwort: []byte("S"), hello: hello}.starte(t, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fertig := make(chan error, 1)
	go func() {
		_, err := (&postgres.Einspielziel{Address: addr, TLS: true, CA: postgres.Zertifikate{ca.Zertifikat}}).Verbinde(ctx, map[string]string{"user": "u"})
		fertig <- err
	}()
	// Die Gegenstelle antwortet `S` und meldet das erste Byte des ClientHello;
	// danach schweigt sie, und die Aushandlung wartet auf sie.
	select {
	case <-hello:
	case <-time.After(10 * time.Second):
		t.Fatal("das ClientHello trifft binnen 10 s nicht ein")
	}
	cancel()
	select {
	case err := <-fertig:
		if err == nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("Verbinde nach dem Abbruch: %v, erwartet den Abbruch als Ursache", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Verbinde endet binnen 5 s nach dem Abbruch in der Aushandlung nicht")
	}
	if l := tlsErgebnis(t, ergebnis); len(l.danach) == 0 || l.danach[0] != 0x16 {
		t.Fatalf("die Gegenstelle sah vor dem Abbruch kein ClientHello: %q", l.danach)
	}
}

// Abdeckung: LH-FA-20/Negative — eine Formatierung (%v, %+v, %#v, %s, %q) eines
// Einspielziels oder seiner Zertifikate gibt weder ein Zertifikat noch das
// Passwort aus (LH-FA-20.a *Passwort*, LH-RB-01).
func TestEinspielzielOhneZertifikat(t *testing.T) {
	ca := testpki.NeueCA(t, "GEHEIM-CA")
	z := &postgres.Einspielziel{Address: "h:1", TLS: true, CA: postgres.Zertifikate{ca.Zertifikat}, Password: "GEHEIMpw"}
	for _, verb := range []string{"%v", "%+v", "%#v", "%s", "%q", "%x", "%d"} {
		for _, arg := range []any{z, *z, z.CA, &z.CA, []postgres.Einspielziel{*z}} {
			got := fmt.Sprintf(verb, arg)
			if strings.Contains(got, "GEHEIM") || strings.Contains(strings.ToLower(got), "47454845494d") {
				t.Errorf("%s gibt ein Geheimnis aus: %s", verb, got)
			}
		}
	}
	if got := fmt.Sprintf("%+v", z); !strings.Contains(got, "h:1") {
		t.Errorf("%%+v des Einspielziels ohne Adresse: %s", got)
	}
}
