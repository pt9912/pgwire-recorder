package postgres_test

import (
	"io"
	"net"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// fehlerMit ist eine ErrorResponse mit dem Schweregrad v im Feld V und s im Feld S.
func fehlerMit(v, s string) *pgproto3.ErrorResponse {
	return &pgproto3.ErrorResponse{Severity: s, SeverityUnlocalized: v, Code: "57P01", Message: "terminating connection due to administrator command"}
}

// schweregradFall ist eine Fehlerantwort mit Schweregrad v im Feld V und s im Feld S und die
// Antwort, ob sie ein Ende des Servers ist.
type schweregradFall struct {
	name string
	v, s string
	ende bool
}

// schweregrade sind die Fehlerantworten, an denen der Adapter ein Ende des Servers
// erkennt (ende) oder nicht: genau FATAL und PANIC, gelesen aus V, ohne V (auch
// leer) aus S.
func schweregrade() []schweregradFall {
	return []schweregradFall{
		{"FATAL", "FATAL", "FATAL", true},
		{"PANIC", "PANIC", "PANIC", true},
		{"V leer, S FATAL", "", "FATAL", true},
		{"V leer, S PANIC", "", "PANIC", true},
		{"V FATAL, S lokalisiert", "FATAL", "SCHWERER FEHLER", true},
		{"V ERROR, S FATAL", "ERROR", "FATAL", false},
		{"V klein geschrieben", "fatal", "FATAL", false},
		{"V leer, S klein geschrieben", "", "panic", false},
		{"V leer, S leer", "", "", false},
		{"ERROR", "ERROR", "ERROR", false},
	}
}

// kodiere ist die Leitungsform von msgs, gefolgt von roh.
func kodiere(t *testing.T, msgs []pgproto3.BackendMessage, roh ...[]byte) []byte {
	t.Helper()
	var b []byte
	for _, m := range msgs {
		var err error
		if b, err = m.Encode(b); err != nil {
			t.Fatal(err)
		}
	}
	for _, r := range roh {
		b = append(b, r...)
	}
	return b
}

// endeAufbau beschreibt den Server von endeServer.
type endeAufbau struct {
	vorab     []byte // sofort nach dem Aufbau, ohne Anfrage
	antwort   []byte // in einem Schreibvorgang auf die erste Anfrage (Query, Flush, Sync)
	schliesst bool   // schließt nach der Antwort; sonst bleibt die Verbindung offen
}

// endeServer sendet nach dem Aufbau vorab, beantwortet die erste Anfrage mit
// antwort und schließt danach die Verbindung nur, wenn schliesst gilt.
func endeServer(t *testing.T, a endeAufbau) string {
	t.Helper()
	return rohServer(t, func(be *pgproto3.Backend, conn net.Conn) {
		if len(a.vorab) > 0 {
			if _, err := conn.Write(a.vorab); err != nil {
				return
			}
		}
		for {
			msg, err := be.Receive()
			if err != nil {
				return
			}
			switch msg.(type) {
			case *pgproto3.Query, *pgproto3.Flush, *pgproto3.Sync:
				if a.antwort == nil {
					continue
				}
				_, _ = conn.Write(a.antwort)
				if a.schliesst {
					return
				}
				a.antwort = nil
			}
		}
	})
}

// weg ist ein Weg, auf dem der Adapter die Antworten einer Interaktion liefert:
// Query, oder Send einer Gruppe und Receive, bis ReadyForQuery, connection_closed
// oder ein Fehler kommt.
type weg struct {
	name string
	los  func(t *testing.T, s driven.UpstreamSession) ([]model.Response, error)
}

// ueber sendet die Gruppe Parse und abschluss und empfängt.
func ueber(abschluss model.ClientMessage) func(*testing.T, driven.UpstreamSession) ([]model.Response, error) {
	return func(t *testing.T, s driven.UpstreamSession) ([]model.Response, error) {
		if err := s.Send(ctx(t), []model.ClientMessage{{Type: model.ClientParse, SQL: "SELECT 1"}, abschluss}); err != nil {
			return nil, err
		}
		var out []model.Response
		for {
			rs, err := s.Receive(ctx(t))
			if err != nil {
				return out, err
			}
			out = append(out, rs...)
			if letzte := out[len(out)-1].Type; letzte == model.ResponseReadyForQuery || letzte == model.ResponseConnectionClosed {
				return out, nil
			}
		}
	}
}

func wege() []weg {
	return []weg{
		{"Query", func(t *testing.T, s driven.UpstreamSession) ([]model.Response, error) {
			return s.Query(ctx(t), "SELECT 1")
		}},
		{"Receive nach Flush", ueber(model.ClientMessage{Type: model.ClientFlush})},
		{"Receive nach Sync", ueber(model.ClientMessage{Type: model.ClientSync})},
	}
}

// endeLauf öffnet die Session zu einem Server nach a und geht den Weg w; er
// meldet einen Fehler, wenn der Weg nicht binnen zwei Sekunden endet, etwa weil
// er auf das Verbindungsende des offenen Servers wartet.
func endeLauf(t *testing.T, a endeAufbau, w weg) ([]model.Response, error) {
	t.Helper()
	s := oeffne(t, endeServer(t, a))
	var out []model.Response
	var err error
	fertigBinnen(t, 2*time.Second, w.name+" wartet auf das Verbindungsende", func() { out, err = w.los(t, s) })
	return out, err
}

func typen(rs []model.Response) []model.ResponseType {
	out := make([]model.ResponseType, len(rs))
	for i, r := range rs {
		out[i] = r.Type
	}
	return out
}

// bisZumEnde sind die Antwortarten von ergebnis() und der Fehlerantwort, gefolgt
// von connection_closed, mit n Fehlerantworten.
func bisZumEnde(n int) []model.ResponseType {
	want := []model.ResponseType{model.ResponseRowDescription, model.ResponseDataRow, model.ResponseCommandComplete}
	for range n {
		want = append(want, model.ResponseErrorResponse)
	}
	return append(want, model.ResponseConnectionClosed)
}

// pruefeEnde prüft, dass rs genau die Antwortarten want trägt und dass die letzte
// Fehlerantwort davor SQLSTATE und Meldung des Servers trägt.
func pruefeEnde(t *testing.T, rs []model.Response, err error, want []model.ResponseType) {
	t.Helper()
	if err != nil {
		t.Fatalf("Fehler statt Ende: %v", err)
	}
	if !reflect.DeepEqual(typen(rs), want) {
		t.Fatalf("Antworten %v, erwartet %v", typen(rs), want)
	}
	if f := rs[len(rs)-2].Fields; f["C"] != "57P01" || f["M"] != "terminating connection due to administrator command" {
		t.Fatalf("Fehlerantwort vor dem Ende: %v", f)
	}
}

// Abdeckung: LH-FA-25/Happy, LH-FA-25/Boundary — Upstream-Adapter, Ende des Servers: Liest
// Query, Receive nach Flush oder Receive nach Sync eine ErrorResponse FATAL oder
// PANIC, liefert er die Antworten bis einschließlich dieser Fehlerantwort und danach
// connection_closed, auch wenn der Server die Verbindung offen lässt; der Schweregrad
// ist V, ein leeres V gilt als fehlend, dann S, und genau FATAL und PANIC sind ein
// Ende. Jeder andere Wert, auch fatal und panic klein geschrieben und V ERROR mit S
// FATAL, ist kein Ende: Endet die Verbindung danach, ist es PGR-E6001 mit SQLSTATE
// und Meldung der Fehlerantwort.
func TestEndeSchweregrad(t *testing.T) {
	for _, c := range schweregrade() {
		for _, w := range wege() {
			t.Run(c.name+"/"+w.name, func(t *testing.T) {
				f := fehlerMit(c.v, c.s)
				rs, err := endeLauf(t, endeAufbau{antwort: kodiere(t, mit(ergebnis(), f)), schliesst: !c.ende}, w)
				if c.ende {
					pruefeEnde(t, rs, err, bisZumEnde(1))
					return
				}
				pruefeAbbruch(t, err, model.CodeUnsupported, f)
			})
		}
	}
}

// Abdeckung: LH-FA-25/Boundary — Upstream-Adapter: Folgen auf eine ErrorResponse ERROR
// eine ErrorResponse FATAL, ist die zweite das Ende und die erste Teil der Antworten;
// die letzte gelesene Fehlerantwort entscheidet.
func TestEndeNachFehlerAntwort(t *testing.T) {
	for _, w := range wege() {
		t.Run(w.name, func(t *testing.T) {
			antwort := kodiere(t, mit(ergebnis(), fehlerAntwort("XX001", "erste"), beendet()))
			rs, err := endeLauf(t, endeAufbau{antwort: antwort}, w)
			pruefeEnde(t, rs, err, bisZumEnde(2))
			if rs[3].Fields["C"] != "XX001" {
				t.Fatalf("erste Fehlerantwort fehlt: %v", rs[3].Fields)
			}
		})
	}
}

// Abdeckung: LH-FA-25/Boundary — Upstream-Adapter: Was der Server hinter der
// Fehlerantwort des Endes sendet, liest er nicht, auch nicht, wenn es mit der
// Fehlerantwort in einem Lesevorgang ankommt (ein ReadyForQuery und eine
// Nachricht, die er nicht lesen könnte); die Antworten enden mit connection_closed.
func TestEndeLiestNichtWeiter(t *testing.T) {
	antwort := kodiere(t, mit(ergebnis(), beendet()), kodiere(t, []pgproto3.BackendMessage{&pgproto3.ReadyForQuery{TxStatus: 'I'}}), rohNachricht('Y', ""))
	for _, w := range wege() {
		t.Run(w.name, func(t *testing.T) {
			rs, err := endeLauf(t, endeAufbau{antwort: antwort}, w)
			pruefeEnde(t, rs, err, bisZumEnde(1))
		})
	}
}

// Abdeckung: LH-FA-25/Boundary — Upstream-Adapter, Ende im Leerlauf: Eine
// ErrorResponse FATAL, die der Server zwischen zwei Interaktionen gesendet hat,
// gehört zur nächsten Anfrage, nach den Nachrichten zwischen den Interaktionen
// (NoticeResponse, ParameterStatus), mit connection_closed danach.
func TestEndeImLeerlauf(t *testing.T) {
	vorab := kodiere(t, []pgproto3.BackendMessage{
		&pgproto3.NoticeResponse{Severity: "NOTICE", Code: "00000", Message: "zwischen"},
		&pgproto3.ParameterStatus{Name: "application_name", Value: "zwischen"},
		beendet(),
	})
	want := []model.ResponseType{model.ResponseNoticeResponse, model.ResponseParameterStatus, model.ResponseErrorResponse, model.ResponseConnectionClosed}
	for _, w := range wege() {
		t.Run(w.name, func(t *testing.T) {
			rs, err := endeLauf(t, endeAufbau{vorab: vorab}, w)
			if err != nil || !reflect.DeepEqual(typen(rs), want) || rs[2].Fields["C"] != "57P01" {
				t.Fatalf("Antworten %v, Fehler %v", typen(rs), err)
			}
		})
	}
}

// Abdeckung: LH-FA-25/Negative — Upstream-Adapter, Ende im Leerlauf: Scheitert das
// Senden der Anfrage, bevor der Adapter die Fehlerantwort FATAL gelesen hat, die der
// Server gesendet hat, ist es PGR-E4003 ohne Antworten; die Fehlerantwort im Puffer
// zählt nicht, bei Query wie bei Send.
func TestEndeSendefehlerImLeerlauf(t *testing.T) {
	for _, w := range wege() {
		t.Run(w.name, func(t *testing.T) {
			s := oeffne(t, endeServer(t, endeAufbau{vorab: kodiere(t, []pgproto3.BackendMessage{beendet()})}))
			schreibhaelfteZu(t, s)
			rs, err := w.los(t, s)
			pruefeAbbruch(t, err, model.CodeConnectionLost, nil)
			if len(rs) != 0 {
				t.Fatalf("Antworten trotz Sendefehler: %v", typen(rs))
			}
		})
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-25/Negative — Upstream-Adapter, Senden nach dem
// Ende: Hat die Leserichtung eine ErrorResponse FATAL oder PANIC gelesen (auch nach
// einer ErrorResponse ERROR), liefert Send keinen Fehler, auch wenn das Schreiben
// scheitert und auch wenn es gelingt; nach einer anderen gelesenen ErrorResponse ist
// ein Sendefehler PGR-E6001 mit deren SQLSTATE.
func TestSendNachDemEnde(t *testing.T) {
	parse := []model.ClientMessage{{Type: model.ClientParse, SQL: "SELECT 1"}, {Type: model.ClientFlush}}
	sync := []model.ClientMessage{{Type: model.ClientSync}}
	type fall struct {
		name    string
		fehler  []pgproto3.BackendMessage
		ende    bool
		letzter *pgproto3.ErrorResponse // die Fehlerantwort, die der Sendefehler nennt
	}
	faelle := []fall{{"ERROR dann FATAL", []pgproto3.BackendMessage{fehlerAntwort("XX001", "erste"), beendet()}, true, nil}}
	for _, c := range schweregrade() {
		f := fehlerMit(c.v, c.s)
		faelle = append(faelle, fall{c.name, []pgproto3.BackendMessage{f}, c.ende, f})
	}
	for _, c := range faelle {
		for _, scheitert := range []bool{false, true} {
			name := c.name + "/Schreiben gelingt"
			if scheitert {
				name = c.name + "/Schreiben scheitert"
			}
			t.Run(name, func(t *testing.T) {
				s := oeffne(t, endeServer(t, endeAufbau{antwort: kodiere(t, c.fehler)}))
				if err := s.Send(ctx(t), parse); err != nil {
					t.Fatal(err)
				}
				if _, err := s.Receive(ctx(t)); err != nil {
					t.Fatal(err)
				}
				if scheitert {
					schreibhaelfteZu(t, s)
				}
				err := s.Send(ctx(t), sync)
				switch {
				case c.ende:
					if err != nil {
						t.Fatalf("Send nach dem Ende: %v", err)
					}
				case scheitert:
					pruefeAbbruch(t, err, model.CodeUnsupported, c.letzter)
				case err != nil:
					t.Fatalf("Send ohne Sendefehler: %v", err)
				}
			})
		}
	}
}

// zaehlendeVerbindung zählt die Aufrufe von Read.
type zaehlendeVerbindung struct {
	net.Conn
	gelesen atomic.Int64
}

func (z *zaehlendeVerbindung) Read(b []byte) (int, error) {
	z.gelesen.Add(1)
	return z.Conn.Read(b)
}

// Abdeckung: LH-FA-25/Boundary — Upstream-Adapter, Senden nach dem Ende: Hat die Leserichtung
// eine ErrorResponse FATAL gelesen, ruft Send Read der Verbindung nicht auf, auch wenn eine
// weitere Nachricht des Servers wartet, und kehrt binnen zwei Sekunden zurück; die wartende
// Nachricht bleibt unverbraucht.
func TestSendLiestNichtNachDemEnde(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = l.Close() })
	weiter := make(chan struct{})
	go func() {
		c, err := l.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		if _, err := c.Write(kodiere(t, []pgproto3.BackendMessage{beendet()})); err != nil {
			return
		}
		<-weiter
		_, _ = c.Write(kodiere(t, []pgproto3.BackendMessage{&pgproto3.ReadyForQuery{TxStatus: 'I'}}))
		_, _ = io.Copy(io.Discard, c)
	}()
	roh, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	z := &zaehlendeVerbindung{Conn: roh}
	s := postgres.NeueSession(z)
	t.Cleanup(func() { _ = s.Close() })
	rs, err := s.Receive(ctx(t))
	if err != nil || len(rs) != 2 || rs[1].Type != model.ResponseConnectionClosed {
		t.Fatalf("Receive: %v, %v", typen(rs), err)
	}
	vorher := z.gelesen.Load()
	close(weiter)
	var sendErr error
	fertigBinnen(t, 2*time.Second, "Send nach dem Ende", func() {
		sendErr = s.Send(ctx(t), []model.ClientMessage{{Type: model.ClientSync}})
	})
	if sendErr != nil {
		t.Fatalf("Send nach dem Ende: %v", sendErr)
	}
	if nachher := z.gelesen.Load(); nachher != vorher {
		t.Fatalf("Send ruft Read auf: %d Aufrufe vor, %d nach dem Senden", vorher, nachher)
	}
}
