package pgwire_test

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/pgwire"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// anmeldeVerbindung startet den Handler mit ctx auf einer Seite einer
// Verbindung (net.Pipe, mit tcp eine TCP-Verbindung auf localhost) und sendet
// die Startnachricht. ende schließt, wenn der Handler zurückkehrt; log nimmt die
// Log-Zeilen aller Stufen auf.
func anmeldeVerbindung(ctx context.Context, t *testing.T, rec *fakeRecorder, log io.Writer, tcp bool) (net.Conn, *pgproto3.Frontend, *pgwire.Server, chan struct{}) {
	t.Helper()
	client, serverSeite := net.Pipe()
	if tcp {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { l.Close() })
		if client, err = net.Dial("tcp", l.Addr().String()); err != nil {
			t.Fatal(err)
		}
		if serverSeite, err = l.Accept(); err != nil {
			t.Fatal(err)
		}
	}
	s := pgwire.NewRecordServer(rec, slog.New(slog.NewTextHandler(log, &slog.HandlerOptions{Level: slog.LevelDebug})))
	ende := make(chan struct{})
	go func() {
		defer close(ende)
		pgwire.Handle(ctx, s, serverSeite)
	}()
	t.Cleanup(func() { client.Close() })
	_ = client.SetDeadline(time.Now().Add(5 * time.Second))
	fe := pgproto3.NewFrontend(client, client)
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
	if err := fe.Flush(); err != nil {
		t.Fatal(err)
	}
	return client, fe, s, ende
}

// rohP ist eine Nachricht p mit dem Rumpf rumpf.
func rohP(rumpf []byte) []byte {
	b := binary.BigEndian.AppendUint32([]byte{'p'}, uint32(len(rumpf)+4))
	return append(b, rumpf...)
}

// anforderungVon liest die nächste Nachricht des Servers als Anforderung der
// Anmeldung.
func anforderungVon(t *testing.T, fe *pgproto3.Frontend) pgproto3.BackendMessage {
	t.Helper()
	msg, err := fe.Receive()
	if err != nil {
		t.Fatalf("Anforderung lesen: %v", err)
	}
	if _, ok := msg.(pgproto3.AuthenticationResponseMessage); !ok {
		t.Fatalf("erwartet eine Nachricht R, erhalten %T", msg)
	}
	return msg
}

// bisAufbauEnde liest bis zum ReadyForQuery und liefert die Typen der Nachrichten.
func bisAufbauEnde(t *testing.T, fe *pgproto3.Frontend) []string {
	t.Helper()
	var typen []string
	for {
		msg, err := fe.Receive()
		if err != nil {
			t.Fatalf("Aufbau lesen: %v", err)
		}
		typen = append(typen, reflect.TypeOf(msg).Elem().Name())
		if _, ok := msg.(*pgproto3.ReadyForQuery); ok {
			return typen
		}
	}
}

// Abdeckung: LH-FA-05/Happy — der Recorder gibt jede Anforderung des Servers
// (Code 3, 5, 10, 11, 12) unverändert an den Client und dessen Antwort, den
// Rumpf der Nachricht p, unverändert an den Use Case: auch eine leere, eine mit
// beliebigen Bytes und eine von 70000 Bytes; auf Code 12 erwartet er keine; die
// Anmeldung endet mit AuthenticationOk vor den übrigen Antworten des Aufbaus, und
// weder Anforderung noch Antwort steht im Log einer Stufe (LH-FA-05.b
// *Vermittlung*, *Antwort des Clients*, *Log und Meldungen*).
func TestAnmeldungVermittelt(t *testing.T) {
	gross := make([]byte, 70000)
	gross[0], gross[69999] = 'G', 'E'
	for _, c := range []struct {
		name     string
		anfrage  model.Anmeldeanforderung
		erwartet pgproto3.BackendMessage
		antwort  []byte
	}{
		{"Klartext", model.Anmeldeanforderung{Code: 3}, &pgproto3.AuthenticationCleartextPassword{}, []byte("GEHEIM\x00")},
		{"MD5", model.Anmeldeanforderung{Code: 5, Daten: []byte{1, 2, 3, 4}}, &pgproto3.AuthenticationMD5Password{Salt: [4]byte{1, 2, 3, 4}}, []byte("md5GEHEIM\x00")},
		{"SASL", model.Anmeldeanforderung{Code: 10, Daten: []byte("SCRAM-SHA-256\x00\x00")}, &pgproto3.AuthenticationSASL{AuthMechanisms: []string{"SCRAM-SHA-256"}}, []byte("SCRAM-SHA-256\x00\x00\x00\x00\x02GEHEIM")},
		{"SASL-Fortsetzung", model.Anmeldeanforderung{Code: 11, Daten: []byte("r=SALZ")}, &pgproto3.AuthenticationSASLContinue{Data: []byte("r=SALZ")}, []byte("c=GEHEIM,\xff\x00")},
		{"leere Antwort", model.Anmeldeanforderung{Code: 3}, &pgproto3.AuthenticationCleartextPassword{}, []byte{}},
		{"Antwort von 70000 Bytes", model.Anmeldeanforderung{Code: 11, Daten: []byte("r=SALZ")}, &pgproto3.AuthenticationSASLContinue{Data: []byte("r=SALZ")}, gross},
		{"SASL-Abschluss ohne Antwort", model.Anmeldeanforderung{Code: 12, Daten: []byte("v=SALZ")}, &pgproto3.AuthenticationSASLFinal{Data: []byte("v=SALZ")}, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			rec := &fakeRecorder{anmeldeFolge: []model.Anmeldeanforderung{c.anfrage}}
			var log syncBuffer
			client, fe, _, _ := anmeldeVerbindung(context.Background(), t, rec, &log, false)
			if got := anforderungVon(t, fe); !reflect.DeepEqual(got, c.erwartet) {
				t.Fatalf("Anforderung %#v, erwartet %#v", got, c.erwartet)
			}
			if c.antwort != nil {
				if _, err := client.Write(rohP(c.antwort)); err != nil {
					t.Fatal(err)
				}
			}
			if got := bisAufbauEnde(t, fe); !reflect.DeepEqual(got, []string{"AuthenticationOk", "ParameterStatus", "ReadyForQuery"}) {
				t.Fatalf("Aufbau nach der Anmeldung: %v", got)
			}
			rec.mu.Lock()
			defer rec.mu.Unlock()
			if len(rec.anmeldeAntworten) != 1 || (c.antwort == nil) != (rec.anmeldeAntworten[0] == nil) || string(rec.anmeldeAntworten[0]) != string(c.antwort) {
				t.Errorf("Antwort an den Use Case: %.20q, erwartet %.20q", rec.anmeldeAntworten, c.antwort)
			}
			if l := log.String(); strings.Contains(l, "GEHEIM") || strings.Contains(l, "SALZ") || strings.Contains(l, "md5") {
				t.Errorf("das Log nennt die Anmeldung:\n%s", l)
			}
		})
	}
}

// Abdeckung: LH-FA-05/Negative — jede andere Nachricht des Clients als p und
// Terminate in der Anmeldung (eine Anfrage, eine Extended-Nachricht, eine
// Funktion, eine nicht lesbare) ist PGR-E6001 mit Fehlerantwort an den Client;
// sie geht nicht an den Use Case, und weder Meldung noch Log nennen ihren Inhalt;
// eine Nachricht p außerhalb der Anmeldung ist wie zuvor nicht unterstützt
// (LH-FA-05.b *Antwort des Clients*, *Log und Meldungen*).
func TestAnmeldungAndereNachricht(t *testing.T) {
	for _, c := range []struct {
		name string
		roh  []byte
	}{
		{"Anfrage", mustEncode(&pgproto3.Query{String: "SELECT 'GEHEIM'"})},
		{"Extended-Nachricht", mustEncode(&pgproto3.Parse{Query: "SELECT 'GEHEIM'"})},
		{"Sync", mustEncode(&pgproto3.Sync{})},
		{"Funktion", mustEncode(&pgproto3.FunctionCall{Function: 1})},
		{"nicht lesbar", append([]byte{'z', 0, 0, 0, 10}, "GEHEIM"...)},
	} {
		t.Run(c.name, func(t *testing.T) {
			rec := &fakeRecorder{anmeldeFolge: []model.Anmeldeanforderung{{Code: 3}}}
			var log syncBuffer
			client, fe, s, _ := anmeldeVerbindung(context.Background(), t, rec, &log, false)
			anforderungVon(t, fe)
			if _, err := client.Write(c.roh); err != nil {
				t.Fatal(err)
			}
			e := fehlerantwort(t, fe)
			if e.Code != "0A000" || !strings.Contains(e.Message, model.CodeUnsupported) || strings.Contains(e.Message, "GEHEIM") {
				t.Errorf("ErrorResponse: %+v", e)
			}
			rec.mu.Lock()
			defer rec.mu.Unlock()
			if rec.anmeldeAntworten[0] != nil || s.FirstErrorCode() != model.CodeUnsupported || strings.Contains(log.String(), "GEHEIM") {
				t.Errorf("Antwort %q, erster Fehler %q, Log:\n%s", rec.anmeldeAntworten[0], s.FirstErrorCode(), log.String())
			}
		})
	}
	for _, mitAnmeldung := range []bool{false, true} {
		rec := &fakeRecorder{}
		if mitAnmeldung {
			rec.anmeldeFolge = []model.Anmeldeanforderung{{Code: 3}}
		}
		client, fe, _, _ := anmeldeVerbindung(context.Background(), t, rec, io.Discard, false)
		if mitAnmeldung {
			anforderungVon(t, fe)
			_, _ = client.Write(rohP([]byte("x\x00")))
		}
		bisAufbauEnde(t, fe)
		_, _ = client.Write(rohP([]byte("x\x00")))
		if e := fehlerantwort(t, fe); !strings.Contains(e.Message, model.CodeUnsupported) || !strings.Contains(e.Message, "PasswordMessage") {
			t.Errorf("Nachricht p nach dem Aufbau (mit Anmeldung %v): %+v", mitAnmeldung, e)
		}
	}
}

func mustEncode(m pgproto3.FrontendMessage) []byte {
	b, err := m.Encode(nil)
	if err != nil {
		panic(err)
	}
	return b
}

// Abdeckung: LH-FA-05/Boundary — beendet der Client die Verbindung in der
// Anmeldung (Terminate, Ende des Datenstroms vor und nach der Anforderung,
// Reset), ist das kein Verbindungsfehler: keine Fehlerantwort, keine Meldung, kein
// Log, keine Session und keine Interaktion (LH-FA-05.b *Ende durch den Client*).
func TestAnmeldungEndeDesClients(t *testing.T) {
	for _, c := range []struct {
		name   string
		tcp    bool
		aktion func(net.Conn)
	}{
		{"Terminate", false, func(client net.Conn) { _, _ = client.Write(mustEncode(&pgproto3.Terminate{})) }},
		{"Ende des Datenstroms nach der Anforderung", false, func(client net.Conn) { client.Close() }},
		{"Ende des Datenstroms vor dem Lesen der Anforderung", false, nil},
		{"Reset", true, func(client net.Conn) {
			_ = client.(*net.TCPConn).SetLinger(0)
			client.Close()
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			rec := &fakeRecorder{anmeldeFolge: []model.Anmeldeanforderung{{Code: 3}}, anmeldeEnde: make(chan struct{})}
			var log syncBuffer
			client, fe, s, ende := anmeldeVerbindung(context.Background(), t, rec, &log, c.tcp)
			if c.aktion == nil {
				client.Close()
			} else {
				anforderungVon(t, fe)
				c.aktion(client)
			}
			warteBis(t, ende, "Handle kehrt nicht zurück")
			rec.mu.Lock()
			defer rec.mu.Unlock()
			if !errors.Is(rec.anmeldeFehler, model.ErrClientEnde) || rec.opened != 0 || len(rec.ends) != 0 || s.FirstErrorCode() != "" || log.String() != "" {
				t.Errorf("Fehler %v, Sessions %d, Enden %v, erster Fehler %q, Log:\n%s", rec.anmeldeFehler, rec.opened, rec.ends, s.FirstErrorCode(), log.String())
			}
			if _, err := fe.Receive(); err == nil {
				t.Error("der Client erhält noch eine Nachricht")
			}
		})
	}
}

func warteBis(t *testing.T, kanal chan struct{}, meldung string) {
	t.Helper()
	select {
	case <-kanal:
	case <-time.After(2 * time.Second):
		t.Fatal(meldung)
	}
}

// Abdeckung: LH-FA-05/Boundary, LH-FA-13/Boundary — das erste Signal unterbricht
// die Anmeldung nicht, weder beim Schreiben der Anforderung an einen Client, der
// noch nicht liest, noch beim Warten auf seine Antwort: Sie läuft zu Ende, und
// der Aufbau endet mit AuthenticationOk (LH-FA-05.b *Abbruch und
// Herunterfahren*).
func TestAnmeldungErstesSignal(t *testing.T) {
	for _, schreibend := range []bool{true, false} {
		name := "beim Warten auf die Antwort"
		if schreibend {
			name = "beim Schreiben der Anforderung"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			rec := &fakeRecorder{anmeldeFolge: []model.Anmeldeanforderung{{Code: 3}}}
			client, fe, _, _ := anmeldeVerbindung(ctx, t, rec, io.Discard, false)
			if !schreibend {
				anforderungVon(t, fe)
			}
			cancel()
			time.Sleep(100 * time.Millisecond)
			if schreibend {
				anforderungVon(t, fe)
			}
			if _, err := client.Write(rohP([]byte("pw\x00"))); err != nil {
				t.Fatal(err)
			}
			if got := bisAufbauEnde(t, fe); got[0] != "AuthenticationOk" {
				t.Fatalf("Aufbau nach der Anmeldung: %v", got)
			}
			rec.mu.Lock()
			defer rec.mu.Unlock()
			if len(rec.anmeldeAntworten) != 1 || string(rec.anmeldeAntworten[0]) != "pw\x00" || rec.anmeldeFehler != nil {
				t.Errorf("Antworten %q, Fehler %v", rec.anmeldeAntworten, rec.anmeldeFehler)
			}
		})
	}
}

// Abdeckung: LH-FA-05/Boundary, LH-FA-13/Boundary — das Zwangsende beendet die
// Anmeldung sofort, auch beim Schreiben der Anforderung an einen Client, der
// nicht liest, und beim Warten auf seine Antwort: Der Aufbau endet mit dem Ende
// des Clients, ohne Fehlerantwort, Meldung und Log, und die Verbindung ist zu
// (LH-FA-05.b *Abbruch und Herunterfahren*).
func TestAnmeldungZwangsende(t *testing.T) {
	for _, wartetAufAntwort := range []bool{false, true} {
		name := "beim Schreiben der Anforderung"
		if wartetAufAntwort {
			name = "beim Warten auf die Antwort"
		}
		t.Run(name, func(t *testing.T) {
			rec := &fakeRecorder{anmeldeFolge: []model.Anmeldeanforderung{{Code: 3}}, anmeldeEnde: make(chan struct{})}
			var log syncBuffer
			_, fe, s, ende := anmeldeVerbindung(context.Background(), t, rec, &log, false)
			if wartetAufAntwort {
				anforderungVon(t, fe)
			}
			time.Sleep(50 * time.Millisecond)
			s.Zwangsende()
			warteBis(t, rec.anmeldeEnde, "die Anmeldung endet nicht")
			warteBis(t, ende, "Handle kehrt nicht zurück")
			rec.mu.Lock()
			defer rec.mu.Unlock()
			if !errors.Is(rec.anmeldeFehler, model.ErrClientEnde) || s.FirstErrorCode() != "" || log.String() != "" {
				t.Errorf("Fehler %v, erster Fehler %q, Log:\n%s", rec.anmeldeFehler, s.FirstErrorCode(), log.String())
			}
			if _, err := fe.Receive(); err == nil {
				t.Error("der Client erhält noch eine Nachricht")
			}
		})
	}
}

// Abdeckung: LH-FA-05/Negative — eine Fehlerantwort des Servers in der Anmeldung
// (nach der Antwort des Clients und ohne Anforderung) geht unverändert und ohne
// AuthenticationOk an den Client; es gibt weder Log-Zeile noch Meldung noch
// Session (LH-FA-05.b *Fehlerantwort des Servers*).
func TestAnmeldungFehlerantwortDesServers(t *testing.T) {
	for _, mitAnforderung := range []bool{false, true} {
		rec := &fakeRecorder{aufbauFehler: true}
		if mitAnforderung {
			rec.anmeldeFolge = []model.Anmeldeanforderung{{Code: 3}}
		}
		var log syncBuffer
		client, fe, s, ende := anmeldeVerbindung(context.Background(), t, rec, &log, false)
		if mitAnforderung {
			anforderungVon(t, fe)
			_, _ = client.Write(rohP([]byte("falsch\x00")))
		}
		if e := fehlerantwort(t, fe); e.Code != "3D000" || e.Severity != "FATAL" || e.Message != "database does not exist" {
			t.Errorf("ErrorResponse (mit Anforderung %v): %+v", mitAnforderung, e)
		}
		warteBis(t, ende, "Handle kehrt nicht zurück")
		if s.FirstErrorCode() != "" || log.String() != "" || rec.opened != 0 {
			t.Errorf("erster Fehler %q, Sessions %d, Log:\n%s", s.FirstErrorCode(), rec.opened, log.String())
		}
	}
}
