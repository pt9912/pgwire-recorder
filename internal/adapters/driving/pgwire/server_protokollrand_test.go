package pgwire_test

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/pgwire"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// randServer ist ein Server im Record- oder Replay-Modus auf einer Loopback-Adresse,
// optional mit TLS; senden schickt rohe Bytes und liest bis zum Ende der Verbindung.
type randServer struct {
	t      *testing.T
	replay bool
	rec    *fakeRecorder
	rep    *fakeReplayer
	s      *pgwire.Server
	l      net.Listener
	tls    *pgwire.TLS
	log    *syncBuffer
	fertig chan struct{}
}

func neuerRandServer(t *testing.T, replay bool, z *pgwire.TLS) *randServer {
	t.Helper()
	r := &randServer{t: t, replay: replay, rec: &fakeRecorder{}, rep: &fakeReplayer{}, tls: z, log: &syncBuffer{}, fertig: make(chan struct{})}
	log := slog.New(slog.NewTextHandler(r.log, nil))
	if replay {
		r.s = pgwire.NewReplayServer(r.rep, log)
	} else {
		r.s = pgwire.NewRecordServer(r.rec, log)
	}
	r.s.MitTLS(z)
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	r.l = l
	go func() {
		r.s.Serve(context.Background(), l)
		close(r.fertig)
	}()
	t.Cleanup(func() { _ = l.Close() })
	return r
}

// verbinde öffnet eine Verbindung; mit TLS und ohne klartext ist sie ausgehandelt.
func (r *randServer) verbinde(klartext bool) (*net.TCPConn, net.Conn) {
	r.t.Helper()
	c, err := net.DialTimeout("tcp", r.l.Addr().String(), 5*time.Second)
	if err != nil {
		r.t.Fatal(err)
	}
	r.t.Cleanup(func() { c.Close() })
	_ = c.SetDeadline(time.Now().Add(10 * time.Second))
	if r.tls == nil || klartext {
		return c.(*net.TCPConn), c
	}
	_, _ = c.Write(sslAnfrage())
	if b := make([]byte, 1); func() bool { _, err := io.ReadFull(c, b); return err != nil || b[0] != 'S' }() {
		r.t.Fatalf("keine Annahme des SSLRequest")
	}
	tc := tls.Client(c, &tls.Config{InsecureSkipVerify: true})
	if err := tc.Handshake(); err != nil {
		r.t.Fatal(err)
	}
	return c.(*net.TCPConn), tc
}

// ende wartet, bis der Server alle Verbindungen beendet hat.
func (r *randServer) ende() {
	r.t.Helper()
	_ = r.l.Close()
	select {
	case <-r.fertig:
	case <-time.After(10 * time.Second):
		r.t.Fatal("der Server beendet seine Verbindungen nicht")
	}
}

// senden schickt daten auf einer Verbindung, schließt sie (mit reset als Reset, sonst
// ihre Schreibseite), liest die Antwort bis zum Ende und liefert den Befund.
func (r *randServer) senden(daten []byte, reset, klartext bool) randBefund {
	r.t.Helper()
	tcp, c := r.verbinde(klartext)
	_, _ = c.Write(daten)
	var roh []byte
	if reset {
		_ = tcp.SetLinger(0)
	} else {
		if tc, ok := c.(*tls.Conn); ok {
			_ = tc.CloseWrite()
		} else {
			_ = tcp.CloseWrite()
		}
		roh, _ = io.ReadAll(c)
	}
	_ = tcp.Close()
	return r.befund(roh)
}

// randBefund ist, was der Client vom Server erhielt und was der Server danach meldet:
// n führende „N“, typen die Nachrichten (R für AuthenticationOk, r für jede andere
// Anmeldung), fehler der Text einer ErrorResponse, sessions die geöffneten Sessions
// und ende (ein Plus eins) das Ende der letzten Record-Session, 0 ohne.
type randBefund struct {
	n        int
	typen    string
	fehler   string
	sessions int
	erster   string
	ende     int
}

func (r *randServer) befund(roh []byte) randBefund {
	r.t.Helper()
	r.ende()
	b := randBefund{erster: r.s.FirstErrorCode()}
	for b.n < len(roh) && roh[b.n] == 'N' {
		b.n++
	}
	for roh = roh[b.n:]; len(roh) >= 5; {
		l := int(binary.BigEndian.Uint32(roh[1:5]))
		if l < 4 || 1+l > len(roh) {
			b.typen += "?"
			break
		}
		rumpf := roh[5 : 1+l]
		switch roh[0] {
		case 'R':
			b.typen += map[bool]string{true: "R", false: "r"}[binary.BigEndian.Uint32(rumpf) == 0]
		case 'E':
			b.typen += "E"
			for f := rumpf; len(f) > 1; {
				i := strings.IndexByte(string(f[1:]), 0)
				if f[0] == 'M' {
					b.fehler = string(f[1 : 1+i])
				}
				f = f[i+2:]
			}
		default:
			b.typen += string(roh[0])
		}
		roh = roh[1+l:]
	}
	if r.replay {
		r.rep.mu.Lock()
		defer r.rep.mu.Unlock()
		b.sessions = r.rep.geoeffnet
		return b
	}
	r.rec.mu.Lock()
	defer r.rec.mu.Unlock()
	b.sessions = r.rec.opened
	if len(r.rec.ends) > 0 {
		b.ende = 1 + int(r.rec.ends[len(r.rec.ends)-1])
	}
	return b
}

func rohKopf(laenge int, code uint32, rest ...byte) []byte {
	return append(binary.BigEndian.AppendUint32(binary.BigEndian.AppendUint32(nil, uint32(laenge)), code), rest...)
}

// rohStart ist eine Startnachricht 3.0 mit den Paaren aus Name und Wert.
func rohStart(paare ...string) []byte {
	var rumpf []byte
	for _, p := range paare {
		rumpf = append(append(rumpf, p...), 0)
	}
	rumpf = append(rumpf, 0)
	return rohKopf(8+len(rumpf), pgproto3.ProtocolVersionNumber, rumpf...)
}

func rohNachricht(typ byte, rumpf string) []byte {
	return append(binary.BigEndian.AppendUint32([]byte{typ}, uint32(len(rumpf)+4)), rumpf...)
}

func sonder(nebennummer uint32) uint32 { return pgwire.MajorSpezial<<16 | nebennummer }

// randFall ist eine Zeile der Tabelle in LH-FA-05.e: die Bytes des Clients und was der
// Server darauf tut. fehler ist der Meldungscode in der ErrorResponse, leer ohne; ende
// ist das erwartete Ende der Record-Session (ein Plus eins), 0 ohne Erwartung.
type randFall struct {
	name     string
	daten    []byte
	reset    bool
	n        int
	typen    string
	fehler   string
	erster   string
	sessions int
	ende     int
}

const (
	aufbauOk     = "RSZ"
	endeClosed   = 1 + int(model.EndClosed)
	endeUnsupp   = 1 + int(model.EndUnsupported)
	endeTerminat = 1 + int(model.EndTerminate)
)

// randLauf führt die Zeilen in beiden Modi aus, mit z als TLS des Servers (nil ohne).
func randLauf(t *testing.T, faelle []randFall, z *pgwire.TLS) {
	t.Helper()
	for _, replay := range []bool{false, true} {
		for _, f := range faelle {
			if z != nil && f.reset {
				continue
			}
			t.Run(fmt.Sprintf("replay=%v/tls=%v/%s", replay, z != nil, f.name), func(t *testing.T) {
				b := neuerRandServer(t, replay, z).senden(f.daten, f.reset, false)
				fehlerOk := (f.fehler == "") == (b.fehler == "") && strings.Contains(b.fehler, f.fehler)
				if b.n != f.n || b.typen != f.typen || !fehlerOk || b.sessions != f.sessions || b.erster != f.erster || (!replay && f.ende != 0 && b.ende != f.ende) {
					t.Fatalf("Befund %+v, erwartet n=%d typen=%q fehler=%q erster=%q sessions=%d ende=%d", b, f.n, f.typen, f.fehler, f.erster, f.sessions, f.ende)
				}
			})
		}
	}
}

func fehlerFall(name string, daten []byte, code string) randFall {
	return randFall{name: name, daten: daten, typen: "E", fehler: code, erster: code}
}

func stillFall(name string, daten []byte, reset bool) randFall {
	return randFall{name: name, daten: daten, reset: reset}
}

func startFaelle() []randFall {
	fs := []randFall{
		{name: "Benutzer, Datenbank und Parameter ungeprüft", daten: rohStart("user", "gibtesnicht", "database", "nix", "options", "-c x=1"), typen: aufbauOk, sessions: 1, ende: endeClosed},
		{name: "_pq_ ohne Punkt ist keine Erweiterung", daten: rohStart("user", "app", "_pq_x", "1"), typen: aufbauOk, sessions: 1},
		{name: "pq_. ohne Unterstrich ist keine Erweiterung", daten: rohStart("user", "app", "pq_.x", "1"), typen: aufbauOk, sessions: 1},
		{name: "_pq_. nicht am Anfang ist keine Erweiterung", daten: rohStart("user", "app", "x_pq_.y", "1"), typen: aufbauOk, sessions: 1},
		fehlerFall("_pq_.", rohStart("user", "app", "_pq_.", "1"), model.CodeUnsupported),
		fehlerFall("_pq_.foo", rohStart("user", "app", "_pq_.foo", "1"), model.CodeUnsupported),
		fehlerFall("Rumpf ohne Paare, Länge 8", rohKopf(8, pgproto3.ProtocolVersionNumber), model.CodeUnsupported),
		fehlerFall("Rumpf mit einem Nullbyte, Länge 9", rohKopf(9, pgproto3.ProtocolVersionNumber, 0), model.CodeUnsupported),
		fehlerFall("Rumpf ohne abschließendes Nullbyte", rohKopf(13, pgproto3.ProtocolVersionNumber, []byte("us\x00x\x00")...), model.CodeUnsupported),
		fehlerFall("Kopf, Rest fehlt", rohKopf(40, pgproto3.ProtocolVersionNumber), model.CodeUnsupported),
		fehlerFall("Sonderanfrage 5677", rohKopf(8, sonder(5677)), model.CodeUnsupported),
		fehlerFall("Sonderanfrage 5681", rohKopf(8, sonder(5681)), model.CodeUnsupported),
		fehlerFall("Sonderanfrage 1", rohKopf(8, sonder(1)), model.CodeUnsupported),
		fehlerFall("CancelRequest Länge 8", rohKopf(8, pgwire.CodeCancelRequest), model.CodeUnsupported),
		stillFall("Verbindung ohne Byte", nil, false),
		stillFall("Verbindung ohne Byte, Reset", nil, true),
	}
	for _, v := range []struct {
		name string
		code uint32
	}{{"0.0", 0}, {"2.0", 2 << 16}, {"3.1", 3<<16 | 1}, {"3.2", 3<<16 | 2}, {"4.0", 4 << 16}, {"65535.1", 65535<<16 | 1}} {
		fs = append(fs, fehlerFall("Version "+v.name, rohKopf(8, v.code), model.CodeProtocolVersion))
	}
	for k := 1; k <= 7; k++ {
		fs = append(fs, stillFall(fmt.Sprintf("%d Bytes", k), rohStart()[:k], false), stillFall(fmt.Sprintf("%d Bytes, Reset", k), rohStart()[:k], true))
	}
	return fs
}

// Abdeckung: LH-FA-05/Negative, LH-FA-05/Boundary — die Zeilen der Tabelle in LH-FA-05.e
// zur Startnachricht, in beiden Modi und mit und ohne TLS (Aushandlung vor den Bytes):
// jede andere Version als 3.0 (0.0, 2.0, 3.1, 3.2, 4.0, 65535.1) ist PGR-E6002 schon
// bei Länge 8; ein Rumpf, den die Bibliothek nicht liest (Länge 8, Länge 9, ohne
// abschließendes Nullbyte) und ein Parameter mit dem Präfix _pq_. sind PGR-E6001, ohne
// dass eine Session entsteht, ein Name ohne dieses Präfix nicht; eine Sonderanfrage mit
// anderer Nebennummer (5677, 5681, 1) und ein CancelRequest der Länge 8 sind PGR-E6001;
// Benutzer, Datenbank und options werden bedient; eine Verbindung mit 0 bis 7 Bytes,
// geschlossen oder zurückgesetzt (der Reset nur ohne TLS), endet still, und eine mit
// Kopf ohne Rest ist PGR-E6001.
func TestProtokollrandStartnachricht(t *testing.T) {
	randLauf(t, startFaelle(), nil)
	z, _ := tlsEin(t, false, testpki.Blatt{})
	randLauf(t, startFaelle(), z)
}

// startcodes liefert die Codes, für die die Prüfreihenfolge der ersten Nachricht gilt: die
// Version 3.0, zwei andere Versionen und jede Sonderanfrage.
func startcodes() []struct {
	name string
	code uint32
} {
	return []struct {
		name string
		code uint32
	}{
		{"3.0", pgproto3.ProtocolVersionNumber}, {"4.0", 4 << 16}, {"65535.1", 65535<<16 | 1},
		{"SSLRequest", pgwire.CodeSSLRequest}, {"GSSENCRequest", pgwire.CodeGSSEncRequest},
		{"CancelRequest", pgwire.CodeCancelRequest}, {"Sonderanfrage 1", sonder(1)},
	}
}

// laengeFaelle sind Köpfe mit einer Länge außerhalb einer Startnachricht (unter 8,
// über 10000), je Code, ohne Rest; der Server endet still, ohne ein Byte Antwort.
func laengeFaelle() []randFall {
	var fs []randFall
	for _, c := range startcodes() {
		for _, l := range []int{4, 7, 10001} {
			fs = append(fs, stillFall(fmt.Sprintf("%s, Länge %d", c.name, l), rohKopf(l, c.code), false))
		}
	}
	return fs
}

// restFaelle sind Köpfe mit gültiger Länge (40), deren Rest fehlt: Ein abgelehnter Code
// ist schon am Kopf abgelehnt, der Rest wird nicht gelesen (sonst wäre es PGR-E6001
// „Startnachricht nicht lesbar“).
func restFaelle() []randFall {
	return []randFall{
		fehlerFall("Version 4.0, Länge 40 ohne Rest", rohKopf(40, 4<<16), model.CodeProtocolVersion),
		fehlerFall("Version 65535.1, Länge 40 ohne Rest", rohKopf(40, 65535<<16|1), model.CodeProtocolVersion),
		{name: "Sonderanfrage 1, Länge 40 ohne Rest", daten: rohKopf(40, sonder(1)), typen: "E", fehler: "unbekannte Sonderanfrage", erster: model.CodeUnsupported},
		{name: "Sonderanfrage 5677, Länge 40 ohne Rest", daten: rohKopf(40, sonder(5677)), typen: "E", fehler: "unbekannte Sonderanfrage", erster: model.CodeUnsupported},
	}
}

// Abdeckung: LH-FA-05/Boundary, LH-FA-05/Negative — Prüfreihenfolge der ersten Nachricht
// (LH-FA-05.e): erst die Länge, dann der Code. Eine Länge unter 8 oder über 10000
// beendet die Verbindung ohne Antwort und ohne Meldung, gleich welcher Code folgt
// (3.0, 4.0, 65535.1, SSLRequest, GSSENCRequest, CancelRequest, Sonderanfrage); bei
// gültiger Länge ist ein abgelehnter Code (Version 4.0, 65535.1, Sonderanfrage 1 und
// 5677) schon am Kopf abgelehnt, ohne dass der fehlende Rest gelesen wird; in beiden
// Modi und mit und ohne TLS.
func TestProtokollrandPruefreihenfolge(t *testing.T) {
	fs := append(laengeFaelle(), restFaelle()...)
	randLauf(t, fs, nil)
	z, _ := tlsEin(t, false, testpki.Blatt{})
	randLauf(t, fs, z)
}

// zaehler liefert die Zahl der Sessions, die der Server in seinem Modus geöffnet hat.
func (r *randServer) zaehler() int {
	if r.replay {
		r.rep.mu.Lock()
		defer r.rep.mu.Unlock()
		return r.rep.geoeffnet
	}
	r.rec.mu.Lock()
	defer r.rec.mu.Unlock()
	return r.rec.opened
}

// Abdeckung: LH-FA-12/Boundary, LH-FA-05/Negative — eine Verbindung, deren Startnachricht
// der Recorder ablehnt (Version außer 3.0, nicht lesbarer Rumpf, Parameter _pq_.,
// unverschlüsselt trotz TLS-Pflicht), zählt nicht: Sie öffnet keine Session, und die
// nächste reguläre Verbindung desselben Servers öffnet die erste, in beiden Modi
// (LH-FA-05.e *Zählung der Verbindungen*).
func TestProtokollrandAbgelehnteVerbindungZaehltNicht(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	for _, f := range []struct {
		name    string
		daten   []byte
		tls     *pgwire.TLS
		klartxt bool
		code    string
	}{
		{"Version 4.0", rohKopf(8, 4<<16), nil, false, model.CodeProtocolVersion},
		{"Rumpf nicht lesbar", rohKopf(9, pgproto3.ProtocolVersionNumber, 0), nil, false, model.CodeUnsupported},
		{"_pq_.", rohStart("user", "app", "_pq_.x", "1"), nil, false, model.CodeUnsupported},
		{"Klartext trotz TLS-Pflicht", rohStart("user", "app"), z, true, model.CodePlaintext},
	} {
		for _, replay := range []bool{false, true} {
			t.Run(fmt.Sprintf("replay=%v/%s", replay, f.name), func(t *testing.T) {
				r := neuerRandServer(t, replay, f.tls)
				_, c := r.verbinde(f.klartxt)
				_, _ = c.Write(f.daten)
				if roh, _ := io.ReadAll(c); !strings.Contains(string(roh), f.code) {
					t.Fatalf("Antwort %q ohne %s", roh, f.code)
				}
				if n := r.zaehler(); n != 0 {
					t.Fatalf("abgelehnte Verbindung öffnete %d Sessions", n)
				}
				_, reg := r.verbinde(false)
				_, _ = reg.Write(append(append(rohStart("user", "app"), rohNachricht('Q', "SELECT 1\x00")...), rohNachricht('X', "")...))
				roh, _ := io.ReadAll(reg)
				if !strings.Contains(string(roh), "SELECT 1") || strings.Contains(string(roh), "PGR-") {
					t.Fatalf("Antwort der regulären Verbindung %q", roh)
				}
				if n := r.zaehler(); n != 1 {
					t.Fatalf("nach der regulären Verbindung %d Sessions, erwartet 1", n)
				}
			})
		}
	}
}

func clientFaelle() []randFall {
	start := rohStart("user", "app")
	nach := func(name string, m []byte, erster string) randFall {
		f := randFall{name: name, daten: append(append([]byte{}, start...), m...), typen: aufbauOk + "E", fehler: erster, erster: erster, sessions: 1, ende: endeUnsupp}
		return f
	}
	parse := rohNachricht('P', "\x00SELECT 1\x00\x00\x00")
	return []randFall{
		nach("CopyData", rohNachricht('d', "x"), model.CodeUnsupported),
		nach("CopyDone", rohNachricht('c', ""), model.CodeUnsupported),
		nach("CopyFail", rohNachricht('f', "x\x00"), model.CodeUnsupported),
		nach("FunctionCall nach dem Aufbau", rohNachricht('F', "\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00"), model.CodeUnsupported),
		nach("FunctionCall in einer Extended-Sequenz", append(parse, rohNachricht('F', "\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00")...), model.CodeUnsupported),
		nach("Nachricht p", rohNachricht('p', "pw\x00"), model.CodeUnsupported),
		nach("Nachricht p ohne Nullbyte", rohNachricht('p', "abc"), model.CodeUnsupported),
		nach("Parse ohne Nullbyte", rohNachricht('P', "abc"), model.CodeUnsupported),
		{name: "Ende nach dem Aufbau", daten: start, typen: aufbauOk, sessions: 1, ende: endeClosed},
		{name: "Terminate", daten: append(append([]byte{}, start...), rohNachricht('X', "")...), typen: aufbauOk, sessions: 1, ende: endeTerminat},
		{name: "Abbruch mitten in einer Query", daten: append(append([]byte{}, start...), append([]byte{'Q', 0, 0, 0, 100}, "SELEC"...)...), typen: aufbauOk, sessions: 1, ende: endeClosed},
		{name: "Abbruch im Kopf einer Query", daten: append(append([]byte{}, start...), 'Q', 0, 0), typen: aufbauOk, sessions: 1, ende: endeClosed},
	}
}

// Abdeckung: LH-FA-05/Negative, LH-FA-05/Boundary — die Zeilen der Tabelle in LH-FA-05.e
// zu Client-Nachrichten nach dem Aufbau, in beiden Modi und mit und ohne TLS: CopyData,
// CopyDone, CopyFail, FunctionCall (nach dem Aufbau und in einer Extended-Sequenz) und
// eine Nachricht p sind PGR-E6001 mit Fehlerantwort, im Record endet die Session als
// nicht unterstützt; ein Rumpf, den die Bibliothek nicht liest (p und Parse ohne
// Nullbyte), ist PGR-E6001 und nicht das Ende des Clients; ein Abbruch mitten in einer
// Query oder in ihrem Kopf ist das Ende des Clients ohne Meldung.
func TestProtokollrandClientNachrichten(t *testing.T) {
	randLauf(t, clientFaelle(), nil)
	z, _ := tlsEin(t, false, testpki.Blatt{})
	randLauf(t, clientFaelle(), z)
}

// Abdeckung: LH-FA-05/Happy, LH-FA-05/Boundary — ohne TLS-Konfiguration beantwortet der
// Recorder in beiden Modi jede Folge aus zwei Anfragen (GSSENCRequest und SSLRequest in
// jeder Reihenfolge) mit je einem N, auch wenn die Startnachricht im selben
// Schreibvorgang folgt, und bedient danach die Startnachricht (LH-FA-05.e).
func TestProtokollrandAnfragenFolgen(t *testing.T) {
	gss, ssl := rohKopf(8, pgwire.CodeGSSEncRequest), rohKopf(8, pgwire.CodeSSLRequest)
	var fs []randFall
	for _, f := range []struct {
		name string
		a, b []byte
	}{{"GSS GSS", gss, gss}, {"GSS SSL", gss, ssl}, {"SSL GSS", ssl, gss}, {"SSL SSL", ssl, ssl}} {
		fs = append(fs, randFall{name: f.name, daten: append(append(append([]byte{}, f.a...), f.b...), rohStart("user", "app")...), n: 2, typen: aufbauOk, sessions: 1})
	}
	randLauf(t, fs, nil)
}

// Abdeckung: LH-FA-23/Happy, LH-FA-05/Boundary — mit TLS-Konfiguration ist die Antwort
// auf einen GSSENCRequest N und ein darauf folgender SSLRequest wird angenommen, in
// beiden Modi (LH-FA-05.e).
func TestProtokollrandGSSdannSSLMitTLS(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	for _, replay := range []bool{false, true} {
		var l tlsLauf
		if replay {
			l = tlsReplay(context.Background(), t, &fakeReplayer{}, z, false)
		} else {
			l = tlsRecord(context.Background(), t, &fakeRecorder{}, z, false)
		}
		_, _ = l.client.Write(rohKopf(8, pgwire.CodeGSSEncRequest))
		if b := make([]byte, 1); func() bool { _, err := io.ReadFull(l.client, b); return err != nil || b[0] != 'N' }() {
			t.Fatalf("replay=%v: keine Antwort N auf den GSSENCRequest", replay)
		}
		tc := l.aushandeln(t, nil)
		startup(t, pgproto3.NewFrontend(tc, tc))
	}
}

// Abdeckung: LH-FA-05/Negative, LH-FA-23/Boundary — ein Parameter _pq_. in einer
// unverschlüsselten Startnachricht an einen Server mit TLS ist PGR-E6001 und nicht
// PGR-E6003 (Rangfolge, LH-FA-23.a *Unverschlüsselte Clients*), in beiden Modi.
func TestProtokollrandErweiterungVorKlartext(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	for _, replay := range []bool{false, true} {
		b := neuerRandServer(t, replay, z).senden(rohStart("user", "app", "_pq_.x", "1"), false, true)
		if b.typen != "E" || b.erster != model.CodeUnsupported || b.sessions != 0 {
			t.Fatalf("replay=%v: Befund %+v", replay, b)
		}
	}
}

// Abdeckung: LH-FA-12/Boundary, LH-FA-05/Boundary — jede Verbindung wird für sich
// behandelt: Eine Verbindung am Protokollrand (hier eine Sonderanfrage, die erst nach
// einer regulären Verbindung vollständig wird) beendet nur sich, die reguläre läuft
// vollständig durch, und eine dritte wird danach bedient, in beiden Modi
// (LH-FA-05.e *Modus*).
func TestProtokollrandGleichzeitig(t *testing.T) {
	for _, replay := range []bool{false, true} {
		r := neuerRandServer(t, replay, nil)
		randTCP, randConn := r.verbinde(false)
		_, _ = randConn.Write(rohKopf(8, sonder(1))[:4])
		for i := 0; i < 2; i++ {
			regulaer, _ := r.verbinde(false)
			_, _ = regulaer.Write(append(append(rohStart("user", "app"), rohNachricht('Q', "SELECT 1\x00")...), rohNachricht('X', "")...))
			roh, _ := io.ReadAll(regulaer)
			if !strings.Contains(string(roh), "SELECT 1") || strings.Contains(string(roh), "PGR-") {
				t.Fatalf("replay=%v, Verbindung %d: Antwort %q", replay, i, roh)
			}
			if i == 0 {
				_, _ = randConn.Write(rohKopf(8, sonder(1))[4:])
				_ = randTCP.CloseWrite()
				if b, _ := io.ReadAll(randConn); len(b) == 0 || b[0] != 'E' {
					t.Fatalf("replay=%v: Antwort der Randverbindung %q", replay, b)
				}
			}
		}
		if b := r.befund(nil); b.erster != model.CodeUnsupported || b.sessions != 2 {
			t.Fatalf("replay=%v: Befund %+v", replay, b)
		}
	}
}
