package postgres_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// roh ist eine Server-Nachricht aus Typ und Rumpf.
func roh(typ byte, rumpf ...byte) []byte {
	out := []byte{typ, 0, 0, 0, 0}
	binary.BigEndian.PutUint32(out[1:], uint32(len(rumpf)+4))
	return append(out, rumpf...)
}

// anmeldung ist eine Nachricht R mit code und den Bytes danach.
func anmeldung(code uint32, rest ...byte) []byte {
	r := binary.BigEndian.AppendUint32(nil, code)
	return roh('R', append(r, rest...)...)
}

// kodiert ist die Folge der kodierten Nachrichten.
func kodiert(t *testing.T, msgs ...pgproto3.BackendMessage) []byte {
	t.Helper()
	var out []byte
	for _, m := range msgs {
		var err error
		if out, err = m.Encode(out); err != nil {
			t.Fatal(err)
		}
	}
	return out
}

// verbunden ist ein Aufbau ohne Anmeldung bis zum ersten ReadyForQuery.
func verbunden(t *testing.T) []byte {
	return kodiert(t, &pgproto3.AuthenticationOk{}, &pgproto3.ParameterStatus{Name: "server_version", Value: "17.0"},
		&pgproto3.BackendKeyData{ProcessID: 1, SecretKey: []byte{1, 2, 3, 4}}, &pgproto3.ReadyForQuery{TxStatus: 'I'})
}

// lauf ist, was der Fake-Server vom Client empfing: die Startup-Parameter und
// alle Bytes danach bis zum Verbindungsende. ende ist nil, wenn der Client die
// Verbindung geschlossen hat, sonst der Grund, aus dem das Lesen endete, etwa
// die Frist des Fake-Servers.
type lauf struct {
	startup map[string]string
	danach  []byte
	ende    error
}

// einspielServer nimmt eine Verbindung an, liest das Startup, sendet aufbau,
// liest dann je Eintrag von jeAnfrage eine Client-Nachricht und sendet den
// Eintrag; danach liest er bis zum Verbindungsende oder seiner Frist von 20 s
// und liefert den lauf. Die Bytes der gelesenen Client-Nachrichten stehen nicht
// in lauf.danach. Ein Eintrag nil schließt die Verbindung, ohne zu lesen.
func einspielServer(t *testing.T, aufbau []byte, jeAnfrage ...[]byte) (string, <-chan lauf) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	ergebnis := make(chan lauf, 1)
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_ = conn.SetDeadline(time.Now().Add(20 * time.Second))
		var laenge [4]byte
		if _, err := io.ReadFull(conn, laenge[:]); err != nil {
			return
		}
		rumpf := make([]byte, binary.BigEndian.Uint32(laenge[:])-4)
		if _, err := io.ReadFull(conn, rumpf); err != nil {
			return
		}
		var sm pgproto3.StartupMessage
		if err := sm.Decode(rumpf); err != nil {
			return
		}
		_, _ = conn.Write(aufbau)
		for _, antwort := range jeAnfrage {
			if antwort == nil {
				ergebnis <- lauf{startup: sm.Parameters}
				return
			}
			var kopf [5]byte
			if _, err := io.ReadFull(conn, kopf[:]); err != nil {
				return
			}
			if _, err := io.ReadFull(conn, make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)); err != nil {
				return
			}
			_, _ = conn.Write(antwort)
		}
		rest, err := io.ReadAll(conn)
		ergebnis <- lauf{startup: sm.Parameters, danach: rest, ende: err}
	}()
	return l.Addr().String(), ergebnis
}

// empfangen wartet höchstens 30 s auf den lauf des Fake-Servers und verlangt,
// dass der Client die Verbindung geschlossen hat, nicht die Frist des
// Fake-Servers das Lesen beendete.
func empfangen(t *testing.T, ergebnis <-chan lauf) lauf {
	t.Helper()
	select {
	case l := <-ergebnis:
		if l.ende != nil {
			t.Fatalf("der Client schließt die Verbindung nicht, das Lesen des Fake-Servers endet mit %v", l.ende)
		}
		return l
	case <-time.After(30 * time.Second):
		t.Fatal("der Fake-Server meldet binnen 30 s kein Verbindungsende")
		return lauf{}
	}
}

// geschlossen wartet nach einem Fehler von Verbinde höchstens 5 s darauf, dass
// der Fake-Server das Schließen durch den Client sieht: Verbinde schließt die
// Verbindung, bevor es zurückkehrt (Abbruch im Aufbau), und die Frist des
// Fake-Servers von 20 s liegt dahinter.
func geschlossen(t *testing.T, ergebnis <-chan lauf) lauf {
	t.Helper()
	select {
	case l := <-ergebnis:
		if l.ende != nil {
			t.Fatalf("der Client schließt die Verbindung nach dem Fehler im Aufbau nicht, das Lesen des Fake-Servers endet mit %v", l.ende)
		}
		return l
	case <-time.After(5 * time.Second):
		t.Fatal("der Fake-Server sieht binnen 5 s nach dem Fehler im Aufbau kein Schließen der Verbindung")
		return lauf{}
	}
}

// Abdeckung: LH-FA-20/Negative — jeder Fehler im Aufbau ist eingestuft und
// schließt die Verbindung, ohne dass play nach dem Startup etwas sendet, auch
// kein Terminate: eine Fehlerantwort der SQLSTATE-Klasse 28 ist PGR-E4005,
// jede andere PGR-E4002, gleich welcher Schweregrad; die Anforderung eines
// Verfahrens (Code 2, 3, 5, 6, 7, 9, 10 und unbekannte) PGR-E4005; eine
// Fortsetzung (8, 11, 12), ein R ohne vollständigen Code, ein R nach
// AuthenticationOk, ReadyForQuery vor AuthenticationOk, jede nicht
// vorgesehene Nachricht, eine nicht lesbare und ein Verbindungsende vor dem
// ersten ReadyForQuery PGR-E4002 (LH-FA-20.a *Aufbau*, *Anmelde-Nachrichten*,
// *Abbruch im Aufbau*).
func TestEinspielAufbauFehler(t *testing.T) {
	ok := kodiert(t, &pgproto3.AuthenticationOk{})
	// mitEnde hängt ein gültiges ReadyForQuery an: Ein Mutant, der die
	// Nachricht davor annähme, endete ohne Fehler statt mit PGR-E4002.
	mitEnde := func(teile ...[]byte) []byte {
		var out []byte
		for _, x := range teile {
			out = append(out, x...)
		}
		return append(out, kodiert(t, &pgproto3.ReadyForQuery{TxStatus: 'I'})...)
	}
	fehler := func(schwere, code string) []byte {
		return kodiert(t, &pgproto3.ErrorResponse{Severity: schwere, SeverityUnlocalized: schwere, Code: code, Message: "nein"})
	}
	for _, f := range []struct {
		name   string
		aufbau []byte
		code   string
		// ende schließt die Verbindung nach aufbau.
		ende bool
	}{
		{"Klasse 28 FATAL", fehler("FATAL", "28P01"), model.CodeLogin, false},
		{"Klasse 28 ERROR", fehler("ERROR", "28000"), model.CodeLogin, false},
		{"Klasse 28 nach AuthenticationOk", append(append([]byte{}, ok...), fehler("FATAL", "28000")...), model.CodeLogin, false},
		{"andere Klasse FATAL", fehler("FATAL", "3D000"), model.CodeUpstream, false},
		{"andere Klasse ERROR", fehler("ERROR", "53300"), model.CodeUpstream, false},
		{"Klasse 2 ähnlich", fehler("FATAL", "2BP01"), model.CodeUpstream, false},
		{"Fehlerantwort nicht lesbar", roh('E', 'S', 'F'), model.CodeUpstream, false},
		{"Kerberos", anmeldung(2), model.CodeLogin, false},
		{"Klartext", kodiert(t, &pgproto3.AuthenticationCleartextPassword{}), model.CodeLogin, false},
		{"MD5", kodiert(t, &pgproto3.AuthenticationMD5Password{Salt: [4]byte{1, 2, 3, 4}}), model.CodeLogin, false},
		{"SCM", anmeldung(6), model.CodeLogin, false},
		{"GSS", anmeldung(7), model.CodeLogin, false},
		{"SSPI", anmeldung(9), model.CodeLogin, false},
		{"SASL", kodiert(t, &pgproto3.AuthenticationSASL{AuthMechanisms: []string{"SCRAM-SHA-256"}}), model.CodeLogin, false},
		{"unbekannt 4", anmeldung(4), model.CodeLogin, false},
		{"unbekannt 99", anmeldung(99), model.CodeLogin, false},
		{"Anforderung mit unlesbarem Rest", anmeldung(5, 1), model.CodeLogin, false},
		{"GSSContinue", mitEnde(anmeldung(8)), model.CodeUpstream, false},
		{"SASLContinue", mitEnde(anmeldung(11)), model.CodeUpstream, false},
		{"SASLFinal", mitEnde(anmeldung(12)), model.CodeUpstream, false},
		{"R ohne vollständigen Code", mitEnde(roh('R', 0, 0, 0)), model.CodeUpstream, false},
		{"AuthenticationOk nicht lesbar", mitEnde(anmeldung(0, 1)), model.CodeUpstream, false},
		{"zweites AuthenticationOk", mitEnde(ok, ok), model.CodeUpstream, false},
		{"R nach ParameterStatus", append(kodiert(t, &pgproto3.AuthenticationOk{}, &pgproto3.ParameterStatus{Name: "a", Value: "b"}), anmeldung(3)...), model.CodeUpstream, false},
		{"ReadyForQuery vor AuthenticationOk", mitEnde(), model.CodeUpstream, false},
		{"NegotiateProtocolVersion", mitEnde(ok, kodiert(t, &pgproto3.NegotiateProtocolVersion{NewestMinorProtocol: 0})), model.CodeUpstream, false},
		{"DataRow", mitEnde(ok, kodiert(t, &pgproto3.DataRow{})), model.CodeUpstream, false},
		{"ParameterStatus nicht lesbar", mitEnde(ok, roh('S', 'a')), model.CodeUpstream, false},
		{"NoticeResponse nicht lesbar", mitEnde(ok, roh('N', 'S')), model.CodeUpstream, false},
		{"BackendKeyData nicht lesbar", mitEnde(ok, roh('K', 1)), model.CodeUpstream, false},
		{"NotificationResponse nicht lesbar", mitEnde(ok, roh('A', 1)), model.CodeUpstream, false},
		{"ReadyForQuery nicht lesbar", append(append([]byte{}, ok...), roh('Z')...), model.CodeUpstream, false},
		{"Länge unter 4", append(append([]byte{}, ok...), 'S', 0, 0, 0, 3), model.CodeUpstream, false},
		{"Ende vor ReadyForQuery", ok, model.CodeUpstream, true},
		{"Ende ohne Nachricht", nil, model.CodeUpstream, true},
	} {
		t.Run(f.name, func(t *testing.T) {
			var jeAnfrage [][]byte
			if f.ende {
				jeAnfrage = [][]byte{nil}
			}
			addr, ergebnis := einspielServer(t, f.aufbau, jeAnfrage...)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			s, err := (&postgres.Einspielziel{Address: addr}).Verbinde(ctx, map[string]string{"user": "u"})
			if s != nil || code(err) != f.code {
				t.Fatalf("Session %v, Fehler %v, erwartet %s", s, err, f.code)
			}
			if l := geschlossen(t, ergebnis); len(l.danach) != 0 {
				t.Fatalf("nach dem Startup gesendet: %q", l.danach)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — im Aufbau liest play BackendKeyData,
// ParameterStatus, NoticeResponse und NotificationResponse an jeder Stelle und
// verwirft sie; die Startup-Parameter gehen unverändert an den Server; was der
// Server hinter dem ersten ReadyForQuery sendet, geht nicht verloren
// (LH-FA-20.a *Aufbau*, *Startup-Daten*).
func TestEinspielAufbauVerworfen(t *testing.T) {
	asynchron := []pgproto3.BackendMessage{
		&pgproto3.NoticeResponse{Severity: "NOTICE", Code: "00000", Message: "hallo"},
		&pgproto3.ParameterStatus{Name: "a", Value: "b"},
		&pgproto3.NotificationResponse{PID: 1, Channel: "k", Payload: "p"},
		&pgproto3.BackendKeyData{ProcessID: 1, SecretKey: []byte{1, 2, 3, 4}},
	}
	aufbau := kodiert(t, asynchron...)
	aufbau = append(aufbau, kodiert(t, &pgproto3.AuthenticationOk{})...)
	aufbau = append(aufbau, kodiert(t, asynchron...)...)
	aufbau = append(aufbau, kodiert(t, &pgproto3.ReadyForQuery{TxStatus: 'I'}, &pgproto3.CommandComplete{CommandTag: []byte("DANACH")})...)
	addr, ergebnis := einspielServer(t, aufbau)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	startup := map[string]string{"user": "u", "database": "d", "application_name": "x", "replication": "false"}
	s, err := (&postgres.Einspielziel{Address: addr}).Verbinde(ctx, startup)
	if err != nil {
		t.Fatalf("Verbinde: %v", err)
	}
	if r, err := s.Naechste(); err != nil || r.Type != model.ResponseCommandComplete || r.Tag != "DANACH" {
		t.Fatalf("Nachricht hinter dem ReadyForQuery: %#v, %v", r, err)
	}
	s.Schliesse()
	l := empfangen(t, ergebnis)
	if !reflect.DeepEqual(l.startup, startup) {
		t.Fatalf("Startup %v, erwartet %v", l.startup, startup)
	}
	if !bytes.Equal(l.danach, kodiertFrontend(t, &pgproto3.Terminate{})) {
		t.Fatalf("nach dem Aufbau gesendet %q, erwartet nur Terminate", l.danach)
	}
}

// kodiertFrontend ist die kodierte Client-Nachricht.
func kodiertFrontend(t *testing.T, m pgproto3.FrontendMessage) []byte {
	t.Helper()
	out, err := m.Encode(nil)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Abdeckung: LH-FA-20/Negative — scheitert das Senden des Startup, ist das
// PGR-E4002 (LH-FA-20.a *Aufbau*).
func TestEinspielStartupNichtGesendet(t *testing.T) {
	client, server := net.Pipe()
	server.Close()
	if err := postgres.Aufbau(client, map[string]string{"user": "u"}); code(err) != model.CodeUpstream {
		t.Fatalf("Startup auf geschlossener Verbindung: %v, erwartet %s", err, model.CodeUpstream)
	}
}

// Abdeckung: LH-FA-20/Negative — ein Server, der nicht erreichbar ist, ist
// PGR-E4002 (LH-FA-20.a).
func TestEinspielNichtErreichbar(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := (&postgres.Einspielziel{Address: addr}).Verbinde(ctx, map[string]string{"user": "u"}); code(err) != model.CodeUpstream {
		t.Fatalf("nicht erreichbar: %v", err)
	}
}

// Abdeckung: LH-FA-20/Boundary — endet ctx während des Aufbaus (zweites
// Signal), bricht Verbinde ab und schließt die Verbindung, ohne Terminate
// (LH-FA-20.a *Abbruchsignal*, *Abbruch im Aufbau*).
func TestEinspielAufbauAbgebrochen(t *testing.T) {
	addr, ergebnis := einspielServer(t, kodiert(t, &pgproto3.AuthenticationOk{}))
	ctx, cancel := context.WithCancel(context.Background())
	fertig := make(chan error, 1)
	go func() {
		_, err := (&postgres.Einspielziel{Address: addr}).Verbinde(ctx, map[string]string{"user": "u"})
		fertig <- err
	}()
	time.Sleep(100 * time.Millisecond)
	cancel()
	select {
	case err := <-fertig:
		if err == nil {
			t.Fatal("Verbinde nach dem Abbruch ohne Fehler")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Verbinde endet binnen 5 s nach dem Abbruch nicht")
	}
	if l := geschlossen(t, ergebnis); len(l.danach) != 0 {
		t.Fatalf("nach dem Startup gesendet: %q", l.danach)
	}
}

// verbinde baut gegen den Fake-Server eine Session auf.
func verbinde(t *testing.T, addr string) interface {
	Anfrage(sql string) error
	Naechste() (model.Response, error)
	Schliesse()
} {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s, err := (&postgres.Einspielziel{Address: addr}).Verbinde(ctx, map[string]string{"user": "u"})
	if err != nil {
		t.Fatalf("Verbinde: %v", err)
	}
	return s
}

// Abdeckung: LH-FA-20/Boundary — nach dem Aufbau liefert Naechste je Aufruf
// die nächste Nachricht, die das Modell kennt, und verwirft
// NotificationResponse und andere ohne Abbildung; nach einer Fehlerantwort
// liest sie nicht weiter, auch wenn dahinter Bytes stehen, die sich nicht
// lesen lassen; Schliesse sendet nach der Anfrage Terminate (LH-FA-20.a
// *Interaktion*, *Ende einer Session*).
func TestEinspielNaechste(t *testing.T) {
	antwort := kodiert(t,
		&pgproto3.NotificationResponse{PID: 1, Channel: "k", Payload: "p"},
		&pgproto3.NoticeResponse{Severity: "NOTICE", Code: "00000", Message: "n"},
		&pgproto3.FunctionCallResponse{Result: []byte("x")},
		&pgproto3.DataRow{Values: [][]byte{[]byte("1")}},
		&pgproto3.ErrorResponse{Severity: "ERROR", Code: "42P01", Message: "fehlt"},
	)
	antwort = append(antwort, 0xff, 0, 0, 0, 1)
	addr, ergebnis := einspielServer(t, verbunden(t), antwort)
	s := verbinde(t, addr)
	if err := s.Anfrage("SELECT 1"); err != nil {
		t.Fatal(err)
	}
	var typen []model.ResponseType
	for {
		r, err := s.Naechste()
		if err != nil {
			t.Fatalf("Naechste nach %v: %v", typen, err)
		}
		typen = append(typen, r.Type)
		if r.Type == model.ResponseErrorResponse {
			break
		}
	}
	want := []model.ResponseType{model.ResponseNoticeResponse, model.ResponseDataRow, model.ResponseErrorResponse}
	if !reflect.DeepEqual(typen, want) {
		t.Fatalf("Antworten %v, erwartet %v", typen, want)
	}
	s.Schliesse()
	if l := empfangen(t, ergebnis); !bytes.Equal(l.danach, kodiertFrontend(t, &pgproto3.Terminate{})) {
		t.Fatalf("nach der Anfrage beim Schließen gesendet %q, erwartet nur Terminate", l.danach)
	}
}

// Abdeckung: LH-FA-20/Negative — nach dem Aufbau sind CopyInResponse,
// CopyOutResponse und CopyBothResponse und eine Nachricht, die sich nicht
// lesen lässt, PGR-E6001, ein Verbindungsende PGR-E4003; ein gescheitertes
// Senden ist PGR-E4003 (LH-FA-20.a *Interaktion*).
func TestEinspielNaechsteFehler(t *testing.T) {
	for _, f := range []struct {
		name    string
		antwort []byte
		code    string
	}{
		{"CopyInResponse", kodiert(t, &pgproto3.CopyInResponse{}), model.CodeUnsupported},
		{"CopyOutResponse", kodiert(t, &pgproto3.CopyOutResponse{}), model.CodeUnsupported},
		{"CopyBothResponse", kodiert(t, &pgproto3.CopyBothResponse{}), model.CodeUnsupported},
		{"unbekannter Typ", roh('!', 1), model.CodeUnsupported},
		{"nicht lesbar", roh('C'), model.CodeUnsupported},
	} {
		t.Run(f.name, func(t *testing.T) {
			addr, ergebnis := einspielServer(t, verbunden(t), f.antwort)
			s := verbinde(t, addr)
			if err := s.Anfrage("SELECT 1"); err != nil {
				t.Fatal(err)
			}
			if _, err := s.Naechste(); code(err) != f.code {
				t.Fatalf("%v, erwartet %s", err, f.code)
			}
			s.Schliesse()
			empfangen(t, ergebnis)
		})
	}
	t.Run("Verbindungsende", func(t *testing.T) {
		addr, ergebnis := einspielServer(t, verbunden(t), nil)
		s := verbinde(t, addr)
		if err := s.Anfrage("SELECT 1"); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Naechste(); code(err) != model.CodeConnectionLost {
			t.Fatalf("Verbindungsende: %v, erwartet %s", err, model.CodeConnectionLost)
		}
		s.Schliesse()
		empfangen(t, ergebnis)
	})
	t.Run("Senden scheitert", func(t *testing.T) {
		addr, ergebnis := einspielServer(t, verbunden(t))
		s := verbinde(t, addr)
		s.Schliesse()
		gesendet := make(chan error, 1)
		go func() { gesendet <- s.Anfrage("SELECT 1") }()
		select {
		case err := <-gesendet:
			if code(err) != model.CodeConnectionLost {
				t.Fatalf("Senden auf geschlossener Verbindung: %v, erwartet %s", err, model.CodeConnectionLost)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Anfrage nach Schliesse endet binnen 5 s nicht")
		}
		empfangen(t, ergebnis)
	})
}

// Abdeckung: LH-FA-20/Negative — die Meldung einer Fehlerantwort im Aufbau
// nennt deren SQLSTATE und Meldung (M), keine weiteren Felder, gleich welche
// Klasse und welcher Schweregrad (LH-FA-20.a *Aufbau*, *Meldungen*).
func TestEinspielAufbauFehlerMeldung(t *testing.T) {
	for _, f := range []struct {
		schwere, code, meldung string
	}{
		{"FATAL", "28P01", "Passwort falsch"},
		{"ERROR", "28000", "Rolle fehlt"},
		{"FATAL", "3D000", "Datenbank fehlt"},
	} {
		t.Run(f.code, func(t *testing.T) {
			aufbau := kodiert(t, &pgproto3.ErrorResponse{Severity: f.schwere, SeverityUnlocalized: f.schwere, Code: f.code, Message: f.meldung, Detail: "DETAIL", Hint: "HINT"})
			addr, ergebnis := einspielServer(t, aufbau)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			_, err := (&postgres.Einspielziel{Address: addr}).Verbinde(ctx, map[string]string{"user": "u"})
			want := ": Fehlerantwort im Aufbau " + f.code + " „" + f.meldung + "“"
			if err == nil || !strings.HasSuffix(err.Error(), want) {
				t.Fatalf("Meldung %v, erwartet mit Ende %q", err, want)
			}
			geschlossen(t, ergebnis)
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — Schliesse (zweites Signal) endet, während
// Anfrage an einen Server sendet, der nicht liest, ohne auf das Senden zu
// warten, und das Senden scheitert danach (LH-FA-20.a *Abbruchsignal*).
func TestEinspielSchliesseBeimSenden(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	angefangen, stop := make(chan struct{}), make(chan struct{})
	t.Cleanup(func() { close(stop) })
	aufbau := verbunden(t)
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		var laenge [4]byte
		if _, err := io.ReadFull(conn, laenge[:]); err != nil {
			return
		}
		if _, err := io.ReadFull(conn, make([]byte, binary.BigEndian.Uint32(laenge[:])-4)); err != nil {
			return
		}
		_, _ = conn.Write(aufbau)
		// Ein Byte der Anfrage zeigt, dass Anfrage sendet; danach liest der
		// Server nichts mehr, bis der Test endet.
		if _, err := io.ReadFull(conn, make([]byte, 1)); err != nil {
			return
		}
		close(angefangen)
		<-stop
	}()
	s := verbinde(t, l.Addr().String())
	gesendet := make(chan error, 1)
	go func() { gesendet <- s.Anfrage(strings.Repeat("x", 64<<20)) }()
	select {
	case <-angefangen:
	case <-time.After(5 * time.Second):
		t.Fatal("der Fake-Server empfängt binnen 5 s kein Byte der Anfrage")
	}
	geschlossen := make(chan struct{})
	go func() {
		s.Schliesse()
		close(geschlossen)
	}()
	select {
	case <-geschlossen:
	case <-time.After(5 * time.Second):
		t.Fatal("Schliesse endet binnen 5 s nicht, während Anfrage sendet")
	}
	select {
	case err := <-gesendet:
		if code(err) != model.CodeConnectionLost {
			t.Fatalf("Anfrage nach Schliesse: %v, erwartet %s", err, model.CodeConnectionLost)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Anfrage endet binnen 5 s nach Schliesse nicht")
	}
}

// Abdeckung: LH-FA-20/Boundary — nimmt die Verbindung kein Terminate an,
// endet Schliesse dennoch und schließt sie (LH-FA-20.a *Abbruchsignal*: soweit
// die Verbindung es sofort annimmt).
func TestEinspielSchliesseOhneAnnahme(t *testing.T) {
	client, server := net.Pipe()
	t.Cleanup(func() { server.Close() })
	s := postgres.NeueEinspielSession(client)
	geschlossen := make(chan struct{})
	go func() {
		s.Schliesse()
		close(geschlossen)
	}()
	select {
	case <-geschlossen:
	case <-time.After(5 * time.Second):
		t.Fatal("Schliesse endet binnen 5 s nicht, wenn die Verbindung nichts annimmt")
	}
	_ = server.SetReadDeadline(time.Now().Add(5 * time.Second))
	if n, err := server.Read(make([]byte, 1)); n != 0 || err != io.EOF {
		t.Fatalf("nach Schliesse %d Byte, %v, erwartet das Verbindungsende", n, err)
	}
}

// Abdeckung: LH-FA-20/Boundary — endet ctx während des Verbindungsversuchs
// (zweites Signal), bricht Verbinde ihn ab und liefert einen Fehler
// (LH-FA-20.a *Abbruchsignal*: auch während des Verbindungsversuchs).
func TestEinspielVersuchAbgebrochen(t *testing.T) {
	imVersuch := make(chan struct{})
	ziel := &postgres.Einspielziel{Address: "127.0.0.1:9", Dialer: net.Dialer{
		ControlContext: func(ctx context.Context, _, _ string, _ syscall.RawConn) error {
			close(imVersuch)
			<-ctx.Done()
			return ctx.Err()
		},
	}}
	ctx, cancel := context.WithCancel(context.Background())
	fertig := make(chan error, 1)
	go func() {
		_, err := ziel.Verbinde(ctx, map[string]string{"user": "u"})
		fertig <- err
	}()
	select {
	case <-imVersuch:
	case <-time.After(5 * time.Second):
		t.Fatal("Verbinde beginnt binnen 5 s keinen Verbindungsversuch")
	}
	cancel()
	select {
	case err := <-fertig:
		if err == nil {
			t.Fatal("Verbinde nach dem Abbruch des Versuchs ohne Fehler")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Verbinde endet binnen 5 s nach dem Abbruch des Verbindungsversuchs nicht")
	}
}

// sendeSperre ist eine Verbindung, deren erstes Write blockiert, bis Close
// sie schließt; ein Write, das beginnt, während ein anderes läuft, zählt sie in
// gleichzeitig und endet sofort mit einem Fehler.
type sendeSperre struct {
	net.Conn
	mu           sync.Mutex
	laufend      int
	gleichzeitig int
	schreibt     chan struct{}
	zu           chan struct{}
	zuEinmal     sync.Once
}

func (c *sendeSperre) Write(b []byte) (int, error) {
	c.mu.Lock()
	c.laufend++
	erstes := c.laufend == 1
	if !erstes {
		c.gleichzeitig++
	}
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		c.laufend--
		c.mu.Unlock()
	}()
	if !erstes {
		return 0, errors.New("gleichzeitiges Write")
	}
	close(c.schreibt)
	<-c.zu
	return 0, net.ErrClosed
}

func (c *sendeSperre) Close() error {
	c.zuEinmal.Do(func() { close(c.zu) })
	return nil
}

func (c *sendeSperre) SetWriteDeadline(time.Time) error { return nil }

// Abdeckung: LH-FA-20/Boundary — Schliesse (zweites Signal) schreibt nichts,
// während Anfrage sendet: kein Terminate zwischen die Bytes der Anfrage; es
// schließt die Verbindung, und das Senden scheitert (LH-FA-20.a
// *Abbruchsignal*).
func TestEinspielSchliesseSchreibtNichtBeimSenden(t *testing.T) {
	conn := &sendeSperre{schreibt: make(chan struct{}), zu: make(chan struct{})}
	s := postgres.NeueEinspielSession(conn)
	gesendet := make(chan error, 1)
	go func() { gesendet <- s.Anfrage("SELECT 1") }()
	select {
	case <-conn.schreibt:
	case <-time.After(5 * time.Second):
		t.Fatal("Anfrage beginnt binnen 5 s nicht zu senden")
	}
	geschlossen := make(chan struct{})
	go func() {
		s.Schliesse()
		close(geschlossen)
	}()
	select {
	case <-geschlossen:
	case <-time.After(5 * time.Second):
		t.Fatal("Schliesse endet binnen 5 s nicht, während Anfrage sendet")
	}
	select {
	case err := <-gesendet:
		if code(err) != model.CodeConnectionLost {
			t.Fatalf("Anfrage nach Schliesse: %v, erwartet %s", err, model.CodeConnectionLost)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Anfrage endet binnen 5 s nach Schliesse nicht")
	}
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.gleichzeitig != 0 {
		t.Fatalf("Schliesse schreibt %d-mal, während Anfrage sendet", conn.gleichzeitig)
	}
}
