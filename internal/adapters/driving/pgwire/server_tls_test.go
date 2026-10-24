package pgwire_test

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/pgwire"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// tlsEin ist ein Server mit TLS, dessen Zertifikat der Test erzeugt hat.
func tlsEin(t *testing.T, klartext bool, b testpki.Blatt) (*pgwire.TLS, []byte) {
	t.Helper()
	blatt := testpki.NeueCA(t, "CA").Server(t, b)
	cert := blatt.PEM()
	z, err := pgwire.NeuesTLS(cert, blatt.SchluesselPEM(t, "PRIVATE KEY"), klartext)
	if err != nil {
		t.Fatal(err)
	}
	return z, cert
}

// paar liefert die zwei Enden einer Verbindung: mit pipe ein net.Pipe, dessen
// Schreiben blockiert, bis die Gegenseite liest, sonst eine Loopback-Verbindung.
func paar(t *testing.T, pipe bool) (client, server net.Conn) {
	t.Helper()
	if pipe {
		client, server = net.Pipe()
	} else {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()
		angenommen := make(chan net.Conn, 1)
		go func() {
			c, _ := l.Accept()
			angenommen <- c
		}()
		if client, err = net.Dial("tcp", l.Addr().String()); err != nil {
			t.Fatal(err)
		}
		select {
		case server = <-angenommen:
		case <-time.After(10 * time.Second):
			t.Fatal("der Listener nimmt die Verbindung nicht binnen 10 s an")
		}
	}
	t.Cleanup(func() { client.Close(); server.Close() })
	_ = client.SetDeadline(time.Now().Add(10 * time.Second))
	return client, server
}

// tlsLauf ist eine Verbindung zu einem Server mit TLS.
type tlsLauf struct {
	client net.Conn
	s      *pgwire.Server
	log    *syncBuffer
	fertig chan struct{}
}

// tlsRecord bedient eine Verbindung mit einem Record-Server mit z (nil: ohne TLS).
func tlsRecord(ctx context.Context, t *testing.T, rec *fakeRecorder, z *pgwire.TLS, pipe bool) tlsLauf {
	t.Helper()
	log := &syncBuffer{}
	return tlsStart(ctx, t, pgwire.NewRecordServer(rec, slog.New(slog.NewTextHandler(log, nil))), z, log, pipe)
}

// tlsReplay bedient eine Verbindung mit einem Replay-Server mit z.
func tlsReplay(ctx context.Context, t *testing.T, rep *fakeReplayer, z *pgwire.TLS, pipe bool) tlsLauf {
	t.Helper()
	log := &syncBuffer{}
	return tlsStart(ctx, t, pgwire.NewReplayServer(rep, slog.New(slog.NewTextHandler(log, nil))), z, log, pipe)
}

func tlsStart(ctx context.Context, t *testing.T, s *pgwire.Server, z *pgwire.TLS, log *syncBuffer, pipe bool) tlsLauf {
	t.Helper()
	s.MitTLS(z)
	client, server := paar(t, pipe)
	fertig := make(chan struct{})
	go func() {
		pgwire.Handle(ctx, s, server)
		close(fertig)
	}()
	return tlsLauf{client, s, log, fertig}
}

// sslAnfrage ist der SSLRequest als Bytes.
func sslAnfrage() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b[0:4], 8)
	binary.BigEndian.PutUint32(b[4:8], pgwire.CodeSSLRequest)
	return b
}

// aushandeln sendet den SSLRequest, erwartet „S“ und führt den Handshake mit
// cfg aus (nil: ohne Prüfung des Zertifikats).
func (l tlsLauf) aushandeln(t *testing.T, cfg *tls.Config) *tls.Conn {
	t.Helper()
	if _, err := l.client.Write(sslAnfrage()); err != nil {
		t.Fatal(err)
	}
	antwort := make([]byte, 1)
	if _, err := io.ReadFull(l.client, antwort); err != nil || antwort[0] != 'S' {
		t.Fatalf("Antwort %q, Fehler %v, erwartet S", antwort, err)
	}
	if cfg == nil {
		cfg = &tls.Config{InsecureSkipVerify: true}
	}
	tc := tls.Client(l.client, cfg)
	if err := tc.Handshake(); err != nil {
		t.Fatalf("Handshake: %v", err)
	}
	return tc
}

// endet wartet, bis die Verbindung beendet ist.
func (l tlsLauf) endet(t *testing.T) {
	t.Helper()
	zurueckBinnen(t, l.fertig, 5*time.Second, "Verbindung")
}

// warnung prüft, dass das Log genau eine Warnung der Stufe WARN mit code
// trägt, mit remote und error, ohne PEM-Text und ohne cert.
func (l tlsLauf) warnung(t *testing.T, code string, cert []byte) {
	t.Helper()
	log := l.log.String()
	if strings.Count(log, "level=WARN") != 1 || !strings.Contains(log, "code="+code) {
		t.Fatalf("Log ohne genau eine Warnung %s:\n%s", code, log)
	}
	if code == model.CodeHandshake {
		for _, attr := range []string{"remote=127.0.0.1:", "error="} {
			if !strings.Contains(log, attr) {
				t.Errorf("Warnung ohne %s:\n%s", attr, log)
			}
		}
	}
	if strings.Contains(log, "BEGIN") || (cert != nil && strings.Contains(log, strings.Split(string(cert), "\n")[1])) {
		t.Errorf("Log nennt Zertifikat oder Schlüssel:\n%s", log)
	}
}

// Abdeckung: LH-FA-23/Happy, LH-FA-23/Boundary — mit TLS antwortet der Recorder
// auf einen SSLRequest mit S und führt den Handshake aus; danach folgt der
// Verbindungsaufbau auf der verschlüsselten Verbindung, eine Anfrage geht an den
// Use Case, ihre Antworten gehen unverändert an den Client, Terminate meldet
// EndTerminate; die Version ist mindestens TLS 1.2, der Recorder fordert kein
// Zertifikat des Clients an und handelt kein ALPN aus (LH-FA-23.a *Aushandlung*).
func TestTLSAushandlung(t *testing.T) {
	for _, klartext := range []bool{false, true} {
		z, _ := tlsEin(t, klartext, testpki.Blatt{})
		rec := &fakeRecorder{}
		l := tlsRecord(context.Background(), t, rec, z, false)
		angefordert := 0
		tc := l.aushandeln(t, &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"h2", "postgresql"},
			GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
				angefordert++
				return &tls.Certificate{}, nil
			},
		})
		if st := tc.ConnectionState(); st.Version < tls.VersionTLS12 || st.NegotiatedProtocol != "" || angefordert != 0 {
			t.Fatalf("Klartext %v: Version %x, ALPN %q, Zertifikat des Clients %d-mal angefordert", klartext, st.Version, st.NegotiatedProtocol, angefordert)
		}
		fe := pgproto3.NewFrontend(tc, tc)
		startup(t, fe)
		fe.Send(&pgproto3.Query{String: "SELECT 'x', NULL"})
		if err := fe.Flush(); err != nil {
			t.Fatal(err)
		}
		if typen := empfange(t, fe, 4); strings.Join(typen, " ") != "RowDescription DataRow CommandComplete ReadyForQuery" {
			t.Fatalf("Antworten %v", typen)
		}
		fe.Send(&pgproto3.Terminate{})
		_ = fe.Flush()
		if end := rec.lastEnd(t); end != model.EndTerminate {
			t.Fatalf("Session-Ende %v", end)
		}
		l.endet(t)
		if got := rec.protokoll(); got != "q:SELECT 'x', NULL zugestellt ende:1" {
			t.Fatalf("Protokoll %q", got)
		}
	}
}

// Abdeckung: LH-FA-23/Happy — mit TLS beantwortet der Replay-Server eine Anfrage
// auf der verschlüsselten Verbindung wie auf einer unverschlüsselten (LH-FA-23.a).
func TestTLSReplay(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	rep := &fakeReplayer{}
	l := tlsReplay(context.Background(), t, rep, z, false)
	tc := l.aushandeln(t, nil)
	fe := pgproto3.NewFrontend(tc, tc)
	startup(t, fe)
	fe.Send(&pgproto3.Query{String: "SELECT 1"})
	_ = fe.Flush()
	if typen := empfange(t, fe, 2); strings.Join(typen, " ") != "CommandComplete ReadyForQuery" {
		t.Fatalf("Antworten %v", typen)
	}
	fe.Send(&pgproto3.Terminate{})
	_ = fe.Flush()
	l.endet(t)
	if got := rep.protokoll(); got != "close" || rep.sentAufrufe() != 1 {
		t.Fatalf("Protokoll %q, Sent %d", got, rep.sentAufrufe())
	}
}

// Abdeckung: LH-FA-23/Negative — ein Client, der höchstens TLS 1.1 spricht, scheitert
// im Handshake (PGR-W3002); mit TLS 1.2 gelingt er (LH-FA-23.a *Aushandlung*).
func TestTLSMindestversion(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	l := tlsRecord(context.Background(), t, &fakeRecorder{}, z, false)
	_, _ = l.client.Write(sslAnfrage())
	_, _ = io.ReadFull(l.client, make([]byte, 1))
	if err := tls.Client(l.client, &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS10, MaxVersion: tls.VersionTLS11}).Handshake(); err == nil {
		t.Fatal("Handshake mit TLS 1.1 gelungen")
	}
	l.endet(t)
	l.warnung(t, model.CodeHandshake, nil)
	l = tlsRecord(context.Background(), t, &fakeRecorder{}, z, false)
	if v := l.aushandeln(t, &tls.Config{InsecureSkipVerify: true, MaxVersion: tls.VersionTLS12}).ConnectionState().Version; v != tls.VersionTLS12 {
		t.Fatalf("Version %x", v)
	}
}

// Abdeckung: LH-FA-23/Negative — ein gescheiterter Handshake endet mit der Warnung
// PGR-W3002, die code, remote und die Ursache trägt und weder Zertifikat noch
// Schlüssel nennt; der Prozessausgang bleibt unverändert (kein gemerkter Fehler):
// der Client bricht nach dem S ab, sendet Klartext nach dem S, oder sendet Bytes
// nach dem SSLRequest, bevor das S bei ihm ankam; dann antwortet der Recorder
// nicht mit S (LH-FA-23.a *Aushandlung*).
func TestTLSHandshakeScheitert(t *testing.T) {
	z, cert := tlsEin(t, true, testpki.Blatt{})
	for _, f := range []struct {
		name string
		tue  func(t *testing.T, l tlsLauf)
	}{
		{"Abbruch nach dem S", func(t *testing.T, l tlsLauf) {
			_, _ = l.client.Write(sslAnfrage())
			_, _ = io.ReadFull(l.client, make([]byte, 1))
			l.client.Close()
		}},
		{"Klartext nach dem S", func(t *testing.T, l tlsLauf) {
			_, _ = l.client.Write(sslAnfrage())
			_, _ = io.ReadFull(l.client, make([]byte, 1))
			fe := pgproto3.NewFrontend(l.client, l.client)
			fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
			_ = fe.Flush()
			_, _ = io.Copy(io.Discard, l.client)
		}},
		{"Bytes vor dem S", func(t *testing.T, l tlsLauf) {
			_, _ = l.client.Write(append(sslAnfrage(), 0x16, 0x03, 0x01, 0x00, 0x01))
			if n, _ := io.Copy(io.Discard, l.client); n != 0 {
				t.Errorf("%d Bytes vom Recorder, erwartet keines (kein S)", n)
			}
		}},
	} {
		t.Run(f.name, func(t *testing.T) {
			rec := &fakeRecorder{}
			l := tlsRecord(context.Background(), t, rec, z, false)
			f.tue(t, l)
			l.endet(t)
			l.warnung(t, model.CodeHandshake, cert)
			if l.s.FirstErrorCode() != "" || rec.opened != 0 {
				t.Fatalf("Fehler gemerkt: %q, Sessions %d (Klartext in der Verbindung)", l.s.FirstErrorCode(), rec.opened)
			}
		})
	}
}

// Abdeckung: LH-FA-23/Negative — ist TLS konfiguriert und Klartext nicht erlaubt,
// beantwortet der Recorder eine StartupMessage der Version 3.0 ohne SSLRequest,
// auch nach einem GSSENCRequest, mit einer ErrorResponse PGR-E6003 (FATAL,
// SQLSTATE 0A000, Exit-Code 6), beendet die Verbindung und öffnet keine Session;
// mit Klartext erlaubt öffnet sie (LH-FA-23.a *Unverschlüsselte Clients*).
func TestTLSUnverschluesselt(t *testing.T) {
	gss := func() []byte {
		b := sslAnfrage()
		binary.BigEndian.PutUint32(b[4:8], pgwire.CodeGSSEncRequest)
		return b
	}
	for _, nachGSS := range []bool{false, true} {
		z, _ := tlsEin(t, false, testpki.Blatt{})
		rec := &fakeRecorder{}
		l := tlsRecord(context.Background(), t, rec, z, false)
		if nachGSS {
			_, _ = l.client.Write(gss())
			if n, err := io.ReadFull(l.client, make([]byte, 1)); n != 1 || err != nil {
				t.Fatal(err)
			}
		}
		fe := pgproto3.NewFrontend(l.client, l.client)
		fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
		_ = fe.Flush()
		e := fehlerantwort(t, fe)
		if e.Severity != "FATAL" || e.Code != "0A000" || !strings.Contains(e.Message, model.CodePlaintext) {
			t.Fatalf("nach GSS %v: ErrorResponse %+v", nachGSS, e)
		}
		l.endet(t)
		if rec.opened != 0 || l.s.FirstErrorCode() != model.CodePlaintext || (&model.Error{Code: l.s.FirstErrorCode()}).ExitCode() != 6 {
			t.Fatalf("Sessions %d, erster Fehler %q", rec.opened, l.s.FirstErrorCode())
		}
	}
	for _, nachGSS := range []bool{false, true} {
		z, _ := tlsEin(t, true, testpki.Blatt{})
		rec := &fakeRecorder{}
		l := tlsRecord(context.Background(), t, rec, z, false)
		if nachGSS {
			_, _ = l.client.Write(gss())
			_, _ = io.ReadFull(l.client, make([]byte, 1))
		}
		fe := pgproto3.NewFrontend(l.client, l.client)
		startup(t, fe)
		if rec.opened != 1 {
			t.Fatalf("Klartext erlaubt, nach GSS %v: Sessions %d", nachGSS, rec.opened)
		}
	}
}

// Abdeckung: LH-FA-23/Boundary — die Rangfolge am Verbindungsaufbau mit TLS ohne
// Klartext: eine andere Protokollversion als 3.0 ist PGR-E6002, eine unbekannte
// Sonderanfrage PGR-E6001, eine erste Nachricht ohne Startnachricht PGR-W3003 ohne
// Antwort und ohne gemerkten Fehler, ein Klartext-CancelRequest PGR-W3001; keiner
// ist PGR-E6003 (LH-FA-23.a *Unverschlüsselte Clients*).
func TestTLSRangfolge(t *testing.T) {
	nachricht := func(laenge int, code uint32) []byte {
		b := make([]byte, laenge)
		binary.BigEndian.PutUint32(b[0:4], uint32(laenge))
		binary.BigEndian.PutUint32(b[4:8], code)
		return b
	}
	for _, f := range []struct {
		name    string
		senden  []byte
		antwort string // Code der ErrorResponse, leer: keine Antwort
		gemerkt string
		warnung string
	}{
		{"andere Version", nachricht(9, 2<<16), model.CodeProtocolVersion, model.CodeProtocolVersion, ""},
		{"unbekannte Sonderanfrage", nachricht(8, pgwire.MajorSpezial<<16|9999), model.CodeUnsupported, model.CodeUnsupported, ""},
		{"keine Startnachricht", []byte("GET / HTTP/1.1\r\nHost: x\r\n\r\n"), "", "", model.CodeForeignProtocol},
		{"CancelRequest", nachricht(16, pgwire.CodeCancelRequest), "", "", model.CodeCancelRequest},
	} {
		t.Run(f.name, func(t *testing.T) {
			z, _ := tlsEin(t, false, testpki.Blatt{})
			rec := &fakeRecorder{}
			l := tlsRecord(context.Background(), t, rec, z, false)
			_, _ = l.client.Write(f.senden)
			if f.antwort != "" {
				if e := fehlerantwort(t, pgproto3.NewFrontend(l.client, l.client)); !strings.Contains(e.Message, f.antwort) {
					t.Fatalf("ErrorResponse %+v", e)
				}
			} else if n, _ := io.Copy(io.Discard, l.client); n != 0 {
				t.Fatalf("%d Bytes Antwort", n)
			}
			l.endet(t)
			if l.s.FirstErrorCode() != f.gemerkt || rec.opened != 0 {
				t.Fatalf("erster Fehler %q, Sessions %d", l.s.FirstErrorCode(), rec.opened)
			}
			if f.warnung != "" {
				l.warnung(t, f.warnung, nil)
			}
		})
	}
}

// Abdeckung: LH-FA-23/Negative — auf der verschlüsselten Verbindung sind ein
// weiterer SSLRequest und ein GSSENCRequest nicht unterstützt (PGR-E6001)
// (LH-FA-23.a *Aushandlung*).
func TestTLSZweiterAufbau(t *testing.T) {
	for _, code := range []uint32{pgwire.CodeSSLRequest, pgwire.CodeGSSEncRequest} {
		z, _ := tlsEin(t, false, testpki.Blatt{})
		l := tlsRecord(context.Background(), t, &fakeRecorder{}, z, false)
		tc := l.aushandeln(t, nil)
		anfrage := sslAnfrage()
		binary.BigEndian.PutUint32(anfrage[4:8], code)
		_, _ = tc.Write(anfrage)
		e := fehlerantwort(t, pgproto3.NewFrontend(tc, tc))
		if !strings.Contains(e.Message, model.CodeUnsupported) {
			t.Fatalf("Code %d: ErrorResponse %+v", code, e)
		}
		l.endet(t)
	}
}

// Abdeckung: LH-FA-23/Boundary — ein Client mit direktem TLS (der Handshake als
// erste Nachricht, ohne SSLRequest) ist keine PGWire-Startnachricht (PGR-W3003),
// auch mit Klartext erlaubt; der Recorder bedient ihn nicht (LH-FA-23.a *Aushandlung*).
func TestTLSDirekt(t *testing.T) {
	for _, klartext := range []bool{false, true} {
		z, _ := tlsEin(t, klartext, testpki.Blatt{})
		l := tlsRecord(context.Background(), t, &fakeRecorder{}, z, false)
		if err := tls.Client(l.client, &tls.Config{InsecureSkipVerify: true}).Handshake(); err == nil {
			t.Fatalf("Klartext %v: direktes TLS bedient", klartext)
		}
		l.endet(t)
		l.warnung(t, model.CodeForeignProtocol, nil)
		if l.s.FirstErrorCode() != "" {
			t.Fatalf("Fehler gemerkt: %q", l.s.FirstErrorCode())
		}
	}
}

// Abdeckung: LH-FA-23/Happy, LH-FA-23/Boundary — mehrere Verbindungen nacheinander
// funktionieren, verschlüsselt und (mit Klartext erlaubt) unverschlüsselt, und ein
// hängender Handshake hält eine zweite Verbindung nicht auf (LH-FA-23.a *Aushandlung*).
func TestTLSNacheinanderUndGleichzeitig(t *testing.T) {
	z, _ := tlsEin(t, true, testpki.Blatt{})
	rec := &fakeRecorder{}
	log := &syncBuffer{}
	s := pgwire.NewRecordServer(rec, slog.New(slog.NewTextHandler(log, nil)))
	s.MitTLS(z)
	ctx, cancel := context.WithCancel(context.Background())
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	fertig := make(chan struct{})
	go func() { s.Serve(ctx, l); close(fertig) }()
	defer func() {
		cancel()
		l.Close()
		select {
		case <-fertig:
		case <-time.After(10 * time.Second):
			t.Error("Serve endet nicht binnen 10 s nach dem Schließen des Listeners")
		}
	}()
	verbinde := func() net.Conn {
		c, err := net.Dial("tcp", l.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		_ = c.SetDeadline(time.Now().Add(10 * time.Second))
		t.Cleanup(func() { c.Close() })
		return c
	}
	sitzung := func(verschluesselt bool) {
		t.Helper()
		c := verbinde()
		var rw io.ReadWriter = c
		if verschluesselt {
			_, _ = c.Write(sslAnfrage())
			_, _ = io.ReadFull(c, make([]byte, 1))
			tc := tls.Client(c, &tls.Config{InsecureSkipVerify: true})
			if err := tc.Handshake(); err != nil {
				t.Fatal(err)
			}
			rw = tc
		}
		fe := pgproto3.NewFrontend(rw, rw)
		startup(t, fe)
		fe.Send(&pgproto3.Terminate{})
		_ = fe.Flush()
	}
	sitzung(true)
	sitzung(false)
	sitzung(true)
	haengend := verbinde()
	_, _ = haengend.Write(sslAnfrage())
	_, _ = io.ReadFull(haengend, make([]byte, 1))
	sitzung(true)
	if rec.opened != 4 || s.FirstErrorCode() != "" {
		t.Fatalf("Sessions %d, erster Fehler %q", rec.opened, s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-23/Boundary — der Recorder prüft das Zertifikat beim Bedienen nicht:
// ein abgelaufenes Zertifikat ohne Namen wird bedient (Ablauf nach dem Start ändert
// laufende und neue Verbindungen nicht); ungültige Bytes sind ein Fehler (LH-FA-23.a).
func TestTLSZertifikatUngeprueft(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{Von: time.Now().Add(-3 * time.Hour), Bis: time.Now().Add(-2 * time.Hour)})
	l := tlsRecord(context.Background(), t, &fakeRecorder{}, z, false)
	tc := l.aushandeln(t, nil)
	startup(t, pgproto3.NewFrontend(tc, tc))
	if z, err := pgwire.NeuesTLS(nil, nil, false); z != nil || err != nil {
		t.Fatalf("ohne Zertifikat und Schlüssel: %v, %v", z, err)
	}
	if _, err := pgwire.NeuesTLS([]byte("ungültig"), []byte("ungültig"), false); err == nil {
		t.Fatalf("Fehler %v", err)
	}
}

// Abdeckung: LH-FA-23/Boundary — Herunterfahren und Zwangsende während der
// Aushandlung beenden sie ohne Meldung, auch ohne PGR-W3002 (LH-FA-23.a *Aushandlung*).
func TestTLSAushandlungEndetBeiHerunterfahren(t *testing.T) {
	for _, zwang := range []bool{false, true} {
		z, _ := tlsEin(t, false, testpki.Blatt{})
		ctx, cancel := context.WithCancel(context.Background())
		l := tlsRecord(ctx, t, &fakeRecorder{}, z, false)
		_, _ = l.client.Write(sslAnfrage())
		_, _ = io.ReadFull(l.client, make([]byte, 1))
		time.Sleep(100 * time.Millisecond)
		if zwang {
			l.s.Zwangsende()
		} else {
			cancel()
		}
		l.endet(t)
		cancel()
		if strings.Contains(l.log.String(), model.CodeHandshake) || l.s.FirstErrorCode() != "" {
			t.Fatalf("Zwangsende %v: Log %s, erster Fehler %q", zwang, l.log.String(), l.s.FirstErrorCode())
		}
	}
}

// Abdeckung: LH-FA-23/Boundary — Herunterfahren wirkt auf eine offene verschlüsselte
// Session wie auf eine unverschlüsselte: Gibt der Use Case das Ende frei, endet die Session
// mit EndShutdown, ohne dass der Client sendet, bei record und bei replay (LH-FA-23.a *Nach der Aushandlung*).
func TestTLSHerunterfahren(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	ctx, cancel := context.WithCancel(context.Background())
	rec := &fakeRecorder{shutdownEndet: true}
	l := tlsRecord(ctx, t, rec, z, false)
	tc := l.aushandeln(t, nil)
	startup(t, pgproto3.NewFrontend(tc, tc))
	cancel()
	l.endet(t)
	if end := rec.lastEnd(t); end != model.EndShutdown {
		t.Fatalf("Session-Ende %v", end)
	}
	ctx, cancel = context.WithCancel(context.Background())
	rep := &fakeReplayer{}
	lr := tlsReplay(ctx, t, rep, z, false)
	tr := lr.aushandeln(t, nil)
	startup(t, pgproto3.NewFrontend(tr, tr))
	cancel()
	lr.endet(t)
	if got := rep.protokoll(); got != "close" {
		t.Fatalf("Protokoll %q", got)
	}
}

// Abdeckung: LH-FA-23/Boundary — das Zwangsende wirkt auf eine verschlüsselte Session wie
// auf eine unverschlüsselte: Läuft eine Interaktion, meldet die Session EndForced und stellt
// PGR-E4006 zu; liest der Client nicht, dauert es höchstens 2 s und wartet nicht auf ihn,
// bei record und bei replay (LH-FA-23.a *Nach der Aushandlung*, LH-FA-13.a).
func TestTLSZwangsende(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	for _, liest := range []bool{true, false} {
		ctx, cancel := context.WithCancel(context.Background())
		rec := &fakeRecorder{zwangErr: e4006()}
		l := tlsRecord(ctx, t, rec, z, true)
		tc := l.aushandeln(t, nil)
		fe := pgproto3.NewFrontend(tc, tc)
		startup(t, fe)
		laufendeInteraktion(t, rec, fe, cancel)
		l.s.Zwangsende()
		if liest {
			if e := fehlerantwort(t, fe); e.Code != "08006" || e.Message != e4006().Error() {
				t.Fatalf("ErrorResponse %+v", e)
			}
		}
		zurueckBinnen(t, l.fertig, 2*time.Second, "Zwangsende record")
		if got := rec.protokoll(); !strings.HasSuffix(got, "ende:6") || l.s.FirstErrorCode() != model.CodeShutdownTimeout {
			t.Fatalf("Protokoll %q, erster Fehler %q", got, l.s.FirstErrorCode())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	rep := &fakeReplayer{zwangFehler: e4006Replay()}
	l := tlsReplay(ctx, t, rep, z, true)
	tc := l.aushandeln(t, nil)
	fe := pgproto3.NewFrontend(tc, tc)
	startup(t, fe)
	replayLaufend(t, rep, fe)
	cancel()
	l.s.Zwangsende()
	zurueckBinnen(t, l.fertig, 2*time.Second, "Zwangsende replay")
	if got := rep.protokoll(); got != "forced close" || l.s.FirstErrorCode() != model.CodeShutdownTimeout {
		t.Fatalf("Protokoll %q, erster Fehler %q", got, l.s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-23/Boundary — das Ende einer verschlüsselten Verbindung mit
// close_notify, ohne close_notify (die Verbindung wird abgebrochen) und ohne Terminate
// ist das Ende einer unverschlüsselten: EndClosed, kein gemerkter Fehler (LH-FA-23.a
// *Nach der Aushandlung*).
func TestTLSEndeOhneCloseNotify(t *testing.T) {
	for name, schliesse := range map[string]func(*tls.Conn){
		"mit close_notify":  func(tc *tls.Conn) { _ = tc.Close() },
		"ohne close_notify": func(tc *tls.Conn) { _ = tc.NetConn().Close() },
	} {
		z, _ := tlsEin(t, false, testpki.Blatt{})
		rec := &fakeRecorder{}
		l := tlsRecord(context.Background(), t, rec, z, false)
		tc := l.aushandeln(t, nil)
		startup(t, pgproto3.NewFrontend(tc, tc))
		schliesse(tc)
		if end := rec.lastEnd(t); end != model.EndClosed || l.s.FirstErrorCode() != "" {
			t.Fatalf("%s: Session-Ende %v, erster Fehler %q", name, end, l.s.FirstErrorCode())
		}
	}
}

// mitschnitt hält die Bytes fest, die der Client von der Verbindung liest.
type mitschnitt struct {
	net.Conn
	mu  sync.Mutex
	roh []byte
}

func (m *mitschnitt) Read(p []byte) (int, error) {
	n, err := m.Conn.Read(p)
	m.mu.Lock()
	m.roh = append(m.roh, p[:n]...)
	m.mu.Unlock()
	return n, err
}

func (m *mitschnitt) laenge() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.roh)
}

// alerts zählt die Records des Typs Alert (21) ab dem Byte ab; sie sind in TLS 1.2
// unverschlüsselt als Typ erkennbar.
func (m *mitschnitt) alerts(ab int) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := 0
	r := m.roh[ab:]
	for len(r) >= 5 {
		if r[0] == 21 {
			n++
		}
		laenge := 5 + int(binary.BigEndian.Uint16(r[3:5]))
		if len(r) < laenge {
			break
		}
		r = r[laenge:]
	}
	return n
}

// Abdeckung: LH-FA-23/Boundary — der Recorder sendet beim Schließen einer
// verschlüsselten Verbindung kein close_notify, auch nicht beim Herunterfahren:
// Der Client liest das Ende, aber keinen Alert (LH-FA-23.a *Aushandlung*). Das
// Zwangsende mit einem Client, der nicht liest, prüft TestTLSZwangsende.
func TestTLSEndeOhneAlert(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	ctx, cancel := context.WithCancel(context.Background())
	l := tlsRecord(ctx, t, &fakeRecorder{shutdownEndet: true}, z, false)
	tap := &mitschnitt{Conn: l.client}
	l.client = tap
	tc := l.aushandeln(t, &tls.Config{InsecureSkipVerify: true, MaxVersion: tls.VersionTLS12})
	startup(t, pgproto3.NewFrontend(tc, tc))
	ab := tap.laenge()
	cancel()
	if _, err := tc.Read(make([]byte, 1)); err == nil {
		t.Fatal("Verbindung nicht beendet")
	}
	l.endet(t)
	if n := tap.alerts(ab); n != 0 {
		t.Fatalf("%d Alerts nach dem Verbindungsaufbau, erwartet keiner", n)
	}
}

// Abdeckung: LH-FA-23/Boundary, LH-FA-05/Boundary — ein SSLRequest oder GSSENCRequest,
// dessen Längenfeld nicht 8 ist, erhält keine Antwort (weder N noch S) und ist
// PGR-E6001, auch mit TLS ohne Handshake; ist die Länge unter 8 oder über dem Höchstmaß,
// gilt als erste Prüfung PGR-W3003 ohne Antwort und ohne gemerkten Fehler (LH-FA-23.a
// *Aushandlung*, LH-FA-05.e).
func TestTLSSonderanfrageMitFalscherLaenge(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	for _, mitTLS := range []bool{false, true} {
		for _, code := range []uint32{pgwire.CodeSSLRequest, pgwire.CodeGSSEncRequest} {
			for _, laenge := range []uint32{7, 9, 12, 10000, 10001} {
				ziel := (*pgwire.TLS)(nil)
				if mitTLS {
					ziel = z
				}
				l := tlsRecord(context.Background(), t, &fakeRecorder{}, ziel, false)
				msg := make([]byte, 8)
				binary.BigEndian.PutUint32(msg[0:4], laenge)
				binary.BigEndian.PutUint32(msg[4:8], code)
				if laenge >= 8 && laenge <= 10000 {
					msg = append(msg, make([]byte, laenge-8)...)
				}
				_, _ = l.client.Write(msg)
				w3003 := laenge < 8 || laenge > 10000
				if w3003 {
					if n, _ := io.Copy(io.Discard, l.client); n != 0 {
						t.Fatalf("TLS %v, Code %d, Länge %d: %d Bytes Antwort", mitTLS, code, laenge, n)
					}
					l.endet(t)
					l.warnung(t, model.CodeForeignProtocol, nil)
					if l.s.FirstErrorCode() != "" {
						t.Fatalf("Fehler gemerkt: %q", l.s.FirstErrorCode())
					}
					continue
				}
				if e := fehlerantwort(t, pgproto3.NewFrontend(l.client, l.client)); !strings.Contains(e.Message, model.CodeUnsupported) {
					t.Fatalf("TLS %v, Code %d, Länge %d: ErrorResponse %+v", mitTLS, code, laenge, e)
				}
				l.endet(t)
				if l.s.FirstErrorCode() != model.CodeUnsupported {
					t.Fatalf("erster Fehler %q", l.s.FirstErrorCode())
				}
			}
		}
	}
}

// Abdeckung: LH-FA-23/Happy — der Server liefert die Kette der Zertifikatsdatei aus:
// Der Client sieht nach dem Handshake das Zertifikat des Servers und das zweite (LH-FA-23.a).
func TestTLSKetteWirdAusgeliefert(t *testing.T) {
	ca := testpki.NeueCA(t, "CA")
	blatt := ca.Server(t, testpki.Blatt{})
	z, err := pgwire.NeuesTLS(append(blatt.PEM(), ca.PEM()...), blatt.SchluesselPEM(t, "PRIVATE KEY"), false)
	if err != nil {
		t.Fatal(err)
	}
	l := tlsRecord(context.Background(), t, &fakeRecorder{}, z, false)
	kette := l.aushandeln(t, nil).ConnectionState().PeerCertificates
	if len(kette) != 2 || !kette[1].Equal(ca.Zertifikat) {
		t.Fatalf("%d Zertifikate ausgeliefert, erwartet 2 mit der Zertifizierungsstelle an zweiter Stelle", len(kette))
	}
}
