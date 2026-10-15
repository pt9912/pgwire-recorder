//go:build integration

package integration_test

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

// rohSitzung verbindet sich mit listen als Client auf der Ebene der Nachrichten und
// liefert die Verbindung und ihren Leser; mit startup führt es den Aufbau bis zum
// ReadyForQuery aus.
func rohSitzung(t *testing.T, listen string, startup bool) (net.Conn, *pgproto3.Frontend) {
	t.Helper()
	c, err := net.DialTimeout("tcp", listen, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	_ = c.SetDeadline(time.Now().Add(15 * time.Second))
	fe := pgproto3.NewFrontend(c, c)
	if !startup {
		return c, fe
	}
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "postgres", "database": "postgres"}})
	if err := fe.Flush(); err != nil {
		t.Fatal(err)
	}
	for {
		msg, err := fe.Receive()
		if err != nil {
			t.Fatalf("Aufbau: %v", err)
		}
		if _, ok := msg.(*pgproto3.ReadyForQuery); ok {
			return c, fe
		}
	}
}

// fehlerAuf liest die nächste Nachricht und liefert den Text der ErrorResponse; jede
// andere Nachricht beendet den Test.
func fehlerAuf(t *testing.T, fe *pgproto3.Frontend) string {
	t.Helper()
	msg, err := fe.Receive()
	if err != nil {
		t.Fatalf("Fehlerantwort lesen: %v", err)
	}
	e, ok := msg.(*pgproto3.ErrorResponse)
	if !ok {
		t.Fatalf("erwartet eine ErrorResponse, erhalten %T", msg)
	}
	return e.Message
}

// rohStartnachricht ist die Startnachricht 3.0 mit dem Benutzer postgres und den Paaren.
func rohStartnachricht(paare ...string) []byte {
	rumpf := []byte("user\x00postgres\x00")
	for _, p := range paare {
		rumpf = append(append(rumpf, p...), 0)
	}
	rumpf = append(rumpf, 0)
	b := binary.BigEndian.AppendUint32(binary.BigEndian.AppendUint32(nil, uint32(8+len(rumpf))), 3<<16)
	return append(b, rumpf...)
}

// aufnehmenMitLeeren zeichnet je Anfrage eine Session auf (record --record-empty-sessions,
// wie --session-assignment connection es von der Aufzeichnung verlangt) nach output.
func aufnehmenMitLeeren(t *testing.T, output string, anfragen ...string) {
	t.Helper()
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-empty-sessions")
	for _, q := range anfragen {
		beobachten(t, rec.listen, q)
	}
	rec.stop(t, 0)
}

// Abdeckung: LH-FA-05/Negative, LH-FA-12/Boundary — am gebauten Binary in beiden Modi:
// eine Startnachricht der Version 3.2 ist PGR-E6002, eine Sonderanfrage 1234.5677, ein
// Parameter _pq_. in der Startnachricht und eine Nachricht p ohne Nullbyte, CopyData
// und FunctionCall nach dem Aufbau sind PGR-E6001 mit Fehlerantwort; die Randverbindung
// beendet nur sich: Eine reguläre Verbindung danach läuft durch, im Replay mit der
// ersten Session, wenn die Randverbindung vor dem Aufbau endete, sonst mit der zweiten;
// der Lauf endet mit Exit-Code 6 (LH-FA-05.e).
func TestE2EProtokollrandBeideModi(t *testing.T) {
	aufnahme := filepath.Join(t.TempDir(), "zwei.yaml")
	aufnehmenMitLeeren(t, aufnahme, "SELECT 1;", "SELECT 2;")
	kopf := func(laenge int, code uint32) []byte {
		return binary.BigEndian.AppendUint32(binary.BigEndian.AppendUint32(nil, uint32(laenge)), code)
	}
	nachricht := func(typ byte, rumpf string) []byte {
		return append(binary.BigEndian.AppendUint32([]byte{typ}, uint32(len(rumpf)+4)), rumpf...)
	}
	for _, f := range []struct {
		name       string
		startup    bool
		daten      []byte
		code       string
		verbraucht bool
	}{
		{"Version 3.2", false, kopf(8, 3<<16|2), "PGR-E6002", false},
		{"Sonderanfrage 1234.5677", false, kopf(8, 1234<<16|5677), "PGR-E6001", false},
		{"Parameter _pq_.", false, rohStartnachricht("_pq_.x", "1"), "PGR-E6001", false},
		{"Nachricht p ohne Nullbyte", true, nachricht('p', "abc"), "PGR-E6001", true},
		{"CopyData", true, nachricht('d', "x"), "PGR-E6001", true},
		{"FunctionCall", true, nachricht('F', "\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00"), "PGR-E6001", true},
	} {
		for _, modus := range []string{"record", "replay"} {
			t.Run(modus+"/"+f.name, func(t *testing.T) {
				var p *recorder
				if modus == "record" {
					p = startRecorder(t, os.Getenv("PGR_UPSTREAM"), filepath.Join(t.TempDir(), "r.yaml"))
				} else {
					p = startProzess(t, "replay", "--input", aufnahme, "--session-assignment", "connection")
				}
				c, fe := rohSitzung(t, p.listen, f.startup)
				if _, err := c.Write(f.daten); err != nil {
					t.Fatal(err)
				}
				if m := fehlerAuf(t, fe); !strings.Contains(m, f.code) {
					t.Fatalf("Fehlerantwort %q ohne %s", m, f.code)
				}
				frage := "SELECT 1;"
				if f.verbraucht && modus == "replay" {
					frage = "SELECT 2;"
				}
				if sicht := beobachten(t, p.listen, frage); !strings.Contains(sicht, "befehl SELECT 1") {
					t.Fatalf("reguläre Verbindung nach der Randverbindung: %s", sicht)
				}
				p.stop(t, 6)
			})
		}
	}
}

// Abdeckung: LH-FA-05/Negative — ein Parameter _pq_. in der Startnachricht ist
// PGR-E6001 mit Exit-Code 6, auch wenn der Upstream nicht erreichbar ist: Der Recorder
// baut keine Verbindung zum Upstream auf (sonst PGR-E4002 und Exit-Code 4)
// (LH-FA-05.e *Startnachricht 3.0*).
func TestE2ERecordErweiterungOhneUpstream(t *testing.T) {
	rec := startRecorder(t, "127.0.0.1:1", filepath.Join(t.TempDir(), "rec.yaml"))
	c, fe := rohSitzung(t, rec.listen, false)
	if _, err := c.Write(rohStartnachricht("_pq_.x", "1")); err != nil {
		t.Fatal(err)
	}
	if m := fehlerAuf(t, fe); !strings.Contains(m, "PGR-E6001") {
		t.Fatalf("Fehlerantwort %q", m)
	}
	rec.stop(t, 6)
	if strings.Contains(rec.stderr.String(), "PGR-E4002") {
		t.Fatalf("Verbindung zum Upstream versucht:\n%s", rec.stderr.String())
	}
}

// Abdeckung: LH-FA-05/Boundary, LH-FA-12/Boundary — eine Verbindung, die ohne ein Byte,
// nach 1 oder nach 7 Bytes endet, durch Schließen oder durch Reset, ist in beiden Modi
// kein Fehler (Exit-Code 0) und keine Session: im Record auch mit
// --record-empty-sessions nur die eine reguläre Verbindung, im Replay erhält sie danach
// die erste Session (LH-FA-05.e *Verbindung ohne erste Nachricht*).
func TestE2EProtokollrandTCPProbe(t *testing.T) {
	for _, modus := range []string{"record", "replay"} {
		output := filepath.Join(t.TempDir(), "rec.yaml")
		var p *recorder
		if modus == "record" {
			p = startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-empty-sessions")
		} else {
			aufnahme := filepath.Join(t.TempDir(), "eine.yaml")
			aufnehmenMitLeeren(t, aufnahme, "SELECT 1;")
			p = startProzess(t, "replay", "--input", aufnahme, "--session-assignment", "connection")
		}
		for _, n := range []int{0, 1, 7} {
			for _, reset := range []bool{false, true} {
				c, err := net.DialTimeout("tcp", p.listen, 5*time.Second)
				if err != nil {
					t.Fatal(err)
				}
				if n > 0 {
					schreibe(t, c, []byte{0, 0, 0, 40, 0, 3, 0}[:n])
				}
				if reset {
					_ = c.(*net.TCPConn).SetLinger(0)
				}
				_ = c.Close()
			}
		}
		if sicht := beobachten(t, p.listen, "SELECT 1;"); !strings.Contains(sicht, "befehl SELECT 1") {
			t.Fatalf("%s: reguläre Verbindung nach den Proben: %s", modus, sicht)
		}
		p.stop(t, 0)
		if modus == "record" {
			if text := lies(t, output); strings.Count(text, "- id:") != 1 {
				t.Errorf("record: erwartet eine Session:\n%s", text)
			}
		}
	}
}

// Abdeckung: LH-FA-05/Happy, LH-FA-05/Negative, LH-FA-06/Negative — gegen PostgreSQL 17:
// LISTEN allein und NOTIFY allein werden aufgezeichnet; LISTEN und NOTIFY auf denselben
// Kanal in einer Anfrage, in zwei Anfragen derselben Session und, bei zwei Sessions,
// bei der nächsten Anfrage der hörenden Session sind PGR-E6001 (Exit-Code 6): Die
// betroffene Session wird nicht aufgezeichnet, die sendende zweier Sessions schon; trifft
// die Benachrichtigung nach der letzten Interaktion der hörenden Session ein, ist das kein
// Fehler (Exit-Code 0), und beide Sessions stehen in der Aufzeichnung (LH-FA-05.e
// *LISTEN und NOTIFY*, *NotificationResponse*).
func TestE2ERecordListenNotify(t *testing.T) {
	for _, f := range []struct {
		name    string
		lauf    func(t *testing.T, listen, output string) error
		sieht   []string
		nicht   []string
		exit    int
		meldung bool
	}{
		{"LISTEN allein", func(t *testing.T, l, _ string) error {
			c := verbinde(t, l)
			defer schliesse(t, c)
			return frage(c, "LISTEN kanal_a")
		}, []string{"LISTEN kanal_a"}, nil, 0, false},
		{"NOTIFY allein", func(t *testing.T, l, _ string) error {
			c := verbinde(t, l)
			defer schliesse(t, c)
			return frage(c, "NOTIFY kanal_b")
		}, []string{"NOTIFY kanal_b"}, nil, 0, false},
		{"beides in einer Anfrage", func(t *testing.T, l, _ string) error {
			return frage(verbinde(t, l), "LISTEN kanal_c; NOTIFY kanal_c")
		}, nil, []string{"kanal_c"}, 6, true},
		{"zwei Anfragen derselben Session", func(t *testing.T, l, _ string) error {
			c := verbinde(t, l)
			if err := frage(c, "LISTEN kanal_d"); err != nil {
				t.Fatal(err)
			}
			return frage(c, "NOTIFY kanal_d")
		}, nil, []string{"kanal_d"}, 6, true},
		{"Benachrichtigung nach der letzten Interaktion", func(t *testing.T, l, output string) error {
			hoerend := verbinde(t, l)
			if err := frage(hoerend, "LISTEN kanal_f"); err != nil {
				t.Fatal(err)
			}
			sendend := verbinde(t, l)
			if err := frage(sendend, "NOTIFY kanal_f"); err != nil {
				t.Fatal(err)
			}
			schliesse(t, sendend)
			wartetAufInhalt(t, output, "NOTIFY kanal_f")
			schliesse(t, hoerend)
			return nil
		}, []string{"NOTIFY kanal_f", "LISTEN kanal_f"}, nil, 0, false},
		{"zwei Sessions", func(t *testing.T, l, output string) error {
			hoerend := verbinde(t, l)
			if err := frage(hoerend, "LISTEN kanal_e"); err != nil {
				t.Fatal(err)
			}
			sendend := verbinde(t, l)
			if err := frage(sendend, "NOTIFY kanal_e"); err != nil {
				t.Fatal(err)
			}
			schliesse(t, sendend)
			wartetAufInhalt(t, output, "NOTIFY kanal_e")
			return frage(hoerend, "SELECT 1")
		}, []string{"NOTIFY kanal_e"}, []string{"LISTEN kanal_e"}, 6, true},
	} {
		t.Run(f.name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "rec.yaml")
			rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
			err := f.lauf(t, rec.listen, output)
			if f.meldung != (err != nil && strings.Contains(err.Error(), "PGR-E6001")) {
				t.Fatalf("Fehler der Anfrage %v, erwartet PGR-E6001: %v", err, f.meldung)
			}
			if !f.meldung && err != nil {
				t.Fatal(err)
			}
			rec.stop(t, f.exit)
			text := lies(t, output)
			for _, s := range f.sieht {
				if !strings.Contains(text, s) {
					t.Errorf("%q fehlt in der Aufzeichnung:\n%s", s, text)
				}
			}
			for _, s := range f.nicht {
				if strings.Contains(text, s) {
					t.Errorf("%q steht in der Aufzeichnung:\n%s", s, text)
				}
			}
		})
	}
}

// Abdeckung: LH-FA-05/Happy, LH-FA-05/Boundary — replay sendet AuthenticationOk als erste
// Nachricht ohne Anforderung und prüft nichts: Ein Benutzer, eine Datenbank und eine
// Option, die die Aufzeichnung nicht kennt, werden bedient; ein Client mit
// require_auth=scram-sha-256 (pgconn) bricht nach dem AuthenticationOk selbst ab, und
// der Lauf meldet das nicht (Exit-Code 0, keine Meldung PGR-) (LH-FA-05.b *Auswahl im
// Replay*, ohne TLS).
func TestE2EReplayAuswahl(t *testing.T) {
	aufnahme, _ := aufnehmen(t, "SELECT 1;")
	rep := startProzess(t, "replay", "--input", aufnahme)
	_, fe := rohSitzung(t, rep.listen, false)
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "gibtesnicht", "database": "nix", "options": "-c x=1"}})
	if err := fe.Flush(); err != nil {
		t.Fatal(err)
	}
	msg, err := fe.Receive()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := msg.(*pgproto3.AuthenticationOk); !ok {
		t.Fatalf("erste Nachricht %T statt AuthenticationOk ohne Anforderung", msg)
	}
	for {
		if msg, err = fe.Receive(); err != nil {
			t.Fatal(err)
		}
		if _, ok := msg.(*pgproto3.ReadyForQuery); ok {
			break
		}
	}
	fe.Send(&pgproto3.Query{String: "SELECT 1;"})
	fe.Send(&pgproto3.Terminate{})
	if err := fe.Flush(); err != nil {
		t.Fatal(err)
	}
	if typen := empfangeBis(t, fe); !strings.Contains(typen, "CommandComplete") {
		t.Fatalf("Antworten %s", typen)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, strings.Replace(dsn(rep.listen), "sslmode=disable", "sslmode=disable&require_auth=scram-sha-256", 1))
	if err == nil {
		_ = conn.Close(ctx)
		t.Fatal("require_auth=scram-sha-256 hat nicht abgebrochen")
	}
	rep.stop(t, 0)
	if strings.Contains(rep.stderr.String(), "PGR-") {
		t.Fatalf("Meldung nach dem Abbruch des Clients:\n%s", rep.stderr.String())
	}
}

// empfangeBis liest bis zum Ende der Verbindung und liefert die Typen der Nachrichten.
func empfangeBis(t *testing.T, fe *pgproto3.Frontend) string {
	t.Helper()
	var typen []string
	for {
		msg, err := fe.Receive()
		if err != nil {
			return strings.Join(typen, " ")
		}
		typen = append(typen, strings.TrimPrefix(fmt.Sprintf("%T", msg), "*pgproto3."))
	}
}
