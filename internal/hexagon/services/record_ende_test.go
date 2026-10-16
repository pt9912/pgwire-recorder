package services_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
	"github.com/pt9912/pgwire-recorder/internal/testuhr"
	"github.com/pt9912/pgwire-recorder/internal/testzeit"
)

// endeFehler ist eine ErrorResponse mit den Schweregrad-Feldern v und s; ein leerer
// Wert lässt das Feld weg.
func endeFehler(v, s string) model.Response {
	f := map[string]string{"C": "57P01", "M": "terminating connection"}
	if v != "" {
		f["V"] = v
	}
	if s != "" {
		f["S"] = s
	}
	return model.Response{Type: model.ResponseErrorResponse, Fields: f}
}

func sEnde() model.Response { return model.Response{Type: model.ResponseConnectionClosed} }

func gSync() []model.ClientMessage {
	return []model.ClientMessage{cParse(), cBind(), cExecute(), cSync()}
}
func gFlush() []model.ClientMessage { return []model.ClientMessage{cParse(), cFlush()} }

// endeLauf ist eine Session mit der Interaktion „SELECT 1“, auf die eine mit Ende folgt.
type endeLauf struct {
	t    *testing.T
	s    *services.RecordService
	repo *fakeRepo
	up   *fakeUpstream
	id   model.SessionID
}

// beginne öffnet die Session; die Anfrage „ENDE“ beantwortet der Upstream mit rs.
func beginne(t *testing.T, rs []model.Response, opts ...services.RecordOption) *endeLauf {
	t.Helper()
	return oeffne(t, rs, []string{"SELECT 1"}, opts...)
}

// oeffne öffnet die Session und führt vorher die Anfragen aus.
func oeffne(t *testing.T, rs []model.Response, vorher []string, opts ...services.RecordOption) *endeLauf {
	t.Helper()
	up := &fakeUpstream{antworten: map[string][]model.Response{"ENDE": rs}}
	repo := &fakeRepo{}
	s, err := services.NewRecordService(context.Background(), up, repo, "rec.yaml", false, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return &endeLauf{t: t, s: s, repo: repo, up: up, id: session(t, s, vorher...)}
}

func (l *endeLauf) einfach() ([]model.Response, error) {
	return l.s.Query(context.Background(), l.id, "ENDE")
}

// erweitert übergibt die Gruppen des Clients und lässt den Upstream rs liefern.
func (l *endeLauf) erweitert(rs []model.Response, gruppen ...[]model.ClientMessage) ([]model.Response, error) {
	l.t.Helper()
	for _, g := range gruppen {
		client(l.t, l.s, l.id, g...)
	}
	l.up.letzte.empfang <- rs
	var out []model.Response
	var err error
	warte(l.t, "AwaitServer", func() { out, err = l.s.AwaitServer(context.Background(), l.id) })
	return out, err
}

// lies ruft AwaitServer und liefert den Fehler; der Aufruf kehrt binnen zwei Sekunden
// zurück.
func (l *endeLauf) lies() error {
	var err error
	warte(l.t, "AwaitServer", func() { _, err = l.s.AwaitServer(context.Background(), l.id) })
	return err
}

func (l *endeLauf) schliessen(end model.SessionEnd) error {
	return l.s.CloseSession(context.Background(), l.id, end)
}

// jeWeg fährt f einmal mit einer einfachen Anfrage und einmal mit einer Sync-Gruppe;
// der Upstream liefert in beiden rs.
func jeWeg(t *testing.T, rs []model.Response, f func(t *testing.T, l *endeLauf, out []model.Response, err error)) {
	t.Helper()
	t.Run("einfach", func(t *testing.T) {
		l := beginne(t, rs)
		out, err := l.einfach()
		f(t, l, out, err)
	})
	t.Run("erweitert", func(t *testing.T) {
		l := beginne(t, nil)
		out, err := l.erweitert(rs, gSync())
		f(t, l, out, err)
	})
}

// antwortenDer sind die Antworten der Interaktion: bei einer Extended-Interaktion
// die Server-Nachrichten der letzten Gruppe.
func antwortenDer(in model.Interaction) []model.Response {
	if in.Request.Type == model.RequestExtended {
		return in.Groups[len(in.Groups)-1].Server
	}
	return in.Responses
}

// endeAngenommen prüft, dass der Use Case want zurückgab und die Session mit der
// Interaktion mit Ende als zweiter und letzter geschrieben wird.
func endeAngenommen(want []model.Response) func(*testing.T, *endeLauf, []model.Response, error) {
	return func(t *testing.T, l *endeLauf, out []model.Response, err error) {
		t.Helper()
		if err != nil || !reflect.DeepEqual(out, want) {
			t.Fatalf("Antworten %#v, Fehler %v, erwartet %#v", out, err, want)
		}
		l.s.Delivered(context.Background(), l.id)
		if err := l.schliessen(model.EndServer); err != nil {
			t.Fatalf("Ende des Servers meldet %v", err)
		}
		sess := l.repo.last(t).Sessions[0]
		if len(sess.Interactions) != 2 || sess.Interactions[1].Sequence != 2 || !reflect.DeepEqual(antwortenDer(sess.Interactions[1]), want) {
			t.Fatalf("Interaktionen %#v", sess.Interactions)
		}
		if err := sess.ValidateEnd(); err != nil {
			t.Fatalf("Form der Session: %v", err)
		}
	}
}

// endeAbgelehnt prüft PGR-E6001 ohne Antworten und dass keine Session geschrieben wird,
// auch nicht die vorherige Interaktion.
func endeAbgelehnt(t *testing.T, l *endeLauf, out []model.Response, err error) {
	t.Helper()
	if codeOf(err) != model.CodeUnsupported || out != nil {
		t.Fatalf("Antworten %#v, Fehler %v, erwartet %s", out, err, model.CodeUnsupported)
	}
	if err := l.schliessen(model.EndFailed); err != nil {
		t.Fatal(err)
	}
	if len(l.repo.writes) != 0 {
		t.Fatalf("Session geschrieben: %#v", l.repo.writes)
	}
}

// Abdeckung: LH-FA-25/Happy, LH-FA-25/Boundary — Record im Kern, einfache Anfrage und
// Extended-Interaktion: Eine Fehlerantwort FATAL oder PANIC (Feld V, ohne V das
// Feld S, ein leeres V wie ein fehlendes) mit connection_closed als letzter
// Antwort wird als letzte Interaktion der Session übernommen, und der Use Case
// liefert dieselben Antworten zurück (LH-FA-25.a *Aufzeichnen*, *Meldung des Endes
// beim Aufzeichnen*); das Ende meldet keinen Fehler (EndServer). Eine Fehlerantwort,
// deren V keine ist, obwohl S FATAL trägt, ist keines.
func TestRecordEndeSchweregrad(t *testing.T) {
	cases := []struct {
		name string
		fehl model.Response
		ok   bool
	}{
		{"V FATAL, S FATAL", endeFehler("FATAL", "FATAL"), true},
		{"V PANIC, S PANIC", endeFehler("PANIC", "PANIC"), true},
		{"V fehlt, S FATAL", endeFehler("", "FATAL"), true},
		{"V leer, S PANIC", endeFehler("", "PANIC"), true},
		{"V FATAL, S ERROR", endeFehler("FATAL", "ERROR"), true},
		{"V ERROR, S FATAL", endeFehler("ERROR", "FATAL"), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rs := []model.Response{c.fehl, sEnde()}
			if c.name == "V leer, S PANIC" {
				rs[0].Fields["V"] = ""
			}
			if c.ok {
				jeWeg(t, rs, endeAngenommen(rs))
				return
			}
			jeWeg(t, rs, endeAbgelehnt)
		})
	}
}

// Abdeckung: LH-FA-25/Negative — Record im Kern: Ein Ende ohne FATAL- oder
// PANIC-Fehlerantwort unmittelbar davor ist nicht aufzeichenbar, je Stelle von
// connection_closed (Anfang der Antworten, hinter einer anderen Antwort, nach einer
// ErrorResponse ERROR, nach einer Fehlerantwort mit einer Nachricht dazwischen, mit
// ReadyForQuery davor oder danach): PGR-E6001, die Session wird nicht übernommen
// (LH-FA-25.a *Stelle von connection_closed beim Aufzeichnen*).
func TestRecordEndeStelle(t *testing.T) {
	fatal := endeFehler("FATAL", "FATAL")
	cases := []struct {
		name string
		rs   []model.Response
		// nurEinfach: In einer Extended-Interaktion schließt ein ReadyForQuery davor die
		// Interaktion, und was danach folgt, gehört zu keiner (PGR-E1000, unverändert).
		nurEinfach bool
	}{
		{"Anfang", []model.Response{sEnde()}, false},
		{"hinter einer anderen Antwort", []model.Response{sCmd(), sEnde()}, false},
		{"nach ERROR", []model.Response{endeFehler("ERROR", "ERROR"), sEnde()}, false},
		{"Nachricht dazwischen", []model.Response{fatal, {Type: model.ResponseNoticeResponse}, sEnde()}, false},
		{"ReadyForQuery davor", []model.Response{sRFQ(), fatal, sEnde()}, true},
		{"ReadyForQuery dazwischen", []model.Response{fatal, sRFQ(), sEnde()}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.nurEinfach {
				l := beginne(t, c.rs)
				out, err := l.einfach()
				endeAbgelehnt(t, l, out, err)
				return
			}
			jeWeg(t, c.rs, endeAbgelehnt)
		})
	}
}

// Abdeckung: LH-FA-25/Boundary — Record im Kern: Was der Upstream hinter
// connection_closed liefert, übernimmt der Use Case nicht: Die Antworten enden mit
// connection_closed, die Aufzeichnung trägt nichts dahinter (LH-FA-25.a *Meldung des
// Endes beim Aufzeichnen*); steht connection_closed hinter einer anderen Antwort als
// einer Fehlerantwort FATAL, bleibt es nicht aufzeichenbar, gleich was dahinter folgt.
func TestRecordEndeNichtsDahinter(t *testing.T) {
	want := []model.Response{endeFehler("FATAL", "FATAL"), sEnde()}
	rs := append(append([]model.Response{}, want...), sRow(), sRFQ())
	jeWeg(t, rs, endeAngenommen(want))
	jeWeg(t, []model.Response{sCmd(), sEnde(), sRFQ()}, endeAbgelehnt)
}

// Abdeckung: LH-FA-25/Happy, LH-FA-25/Negative, LH-FA-18/Boundary — Record im Kern,
// Extended-Interaktion: Das Ende liegt in der letzten begonnenen Gruppe, die mit Sync
// oder mit Flush abgeschlossen ist, auch mit einer Gruppe davor (Flush) und mit
// ERROR und FATAL in derselben Interaktion, in der Reihenfolge des Upstreams; ist die
// Gruppe nicht abgeschlossen, ist das PGR-E6001 und die Session nicht übernommen
// (LH-FA-25.a *Stelle von connection_closed beim Aufzeichnen*).
func TestRecordEndeGruppe(t *testing.T) {
	ende := []model.Response{endeFehler("FATAL", "FATAL"), sEnde()}
	t.Run("Sync-Gruppe", func(t *testing.T) {
		l := beginne(t, nil)
		out, err := l.erweitert(ende, gSync())
		endeAngenommen(ende)(t, l, out, err)
	})
	t.Run("Flush-Gruppe", func(t *testing.T) {
		l := beginne(t, nil)
		out, err := l.erweitert(ende, gFlush())
		endeAngenommen(ende)(t, l, out, err)
	})
	t.Run("Gruppe davor", func(t *testing.T) {
		l := beginne(t, nil)
		client(t, l.s, l.id, gFlush()...)
		server(t, l.s, l.up, l.id, sParse())
		out, err := l.erweitert(ende, []model.ClientMessage{cBind(), cExecute(), cFlush()})
		endeAngenommen(ende)(t, l, out, err)
		in := l.repo.last(t).Sessions[0].Interactions[1]
		if len(in.Groups) != 2 || !reflect.DeepEqual(in.Groups[0].Server, []model.Response{sParse()}) {
			t.Fatalf("Gruppen %#v", in.Groups)
		}
	})
	t.Run("ERROR vor FATAL", func(t *testing.T) {
		l := beginne(t, nil)
		rs := []model.Response{endeFehler("ERROR", "ERROR"), endeFehler("FATAL", "FATAL"), sEnde()}
		out, err := l.erweitert(rs, gFlush())
		endeAngenommen(rs)(t, l, out, err)
	})
	t.Run("nicht abgeschlossene Gruppe", func(t *testing.T) {
		l := beginne(t, nil)
		out, err := l.erweitert(ende, gFlush(), []model.ClientMessage{cBind()})
		endeAbgelehnt(t, l, out, err)
	})
}

// Abdeckung: LH-FA-25/Negative — Record im Kern, Reihenfolge der Gründe: Ein Fehler des
// Upstream-Ports geht der Prüfung der Form vor; er bleibt PGR-E4003, die Session ist
// übernehmbar (LH-FA-25.a *Reihenfolge der Gründe beim Aufzeichnen*).
func TestRecordEndePortfehlerVorForm(t *testing.T) {
	up := &fakeUpstream{queryErr: model.Errorf(model.CodeConnectionLost, nil, "weg")}
	s, repo := neu(t, up)
	id := session(t, s, "SELECT 1")
	out, err := s.Query(context.Background(), id, "MIT-FEHLER")
	if codeOf(err) != model.CodeConnectionLost || len(out) != 1 {
		t.Fatalf("Antworten %#v, Fehler %v, erwartet %s", out, err, model.CodeConnectionLost)
	}
	schliessen(t, s, id, model.EndFailed)
	if got := repo.last(t).Sessions[0].Interactions; len(got) != 1 {
		t.Fatalf("Interaktionen %#v", got)
	}
}

// Abdeckung: LH-FA-25/Negative — Record im Kern: Nach der Interaktion mit Ende nimmt
// der Use Case nichts mehr an (ErrSessionEnded), je Aufruf: eine einfache Anfrage
// (nicht an den Upstream weitergegeben), eine Client-Nachricht, das Lesen des Upstreams
// (nicht gelesen, auch nicht, was dort wartet), eine Anfrage, die bei dem Ende schon
// wartete, und eine Interaktion, die der Client schon vollständig bis zum Sync in die
// Pipeline gestellt hat; die Session trägt das Ende einmal, in der letzten Interaktion
// (LH-FA-25.a *Stelle von connection_closed beim Aufzeichnen*).
func TestRecordEndeDanachNichts(t *testing.T) {
	ende := []model.Response{endeFehler("FATAL", "FATAL"), sEnde()}
	// zustellen liest die Antworten des Endes und meldet sie als zugestellt.
	zustellen := func(t *testing.T, l *endeLauf) {
		t.Helper()
		if err := l.lies(); err != nil {
			t.Fatal(err)
		}
		l.s.Delivered(context.Background(), l.id)
	}
	// nachEnde lässt den Upstream das Ende liefern und lässt die Session offen; mit
	// herunterfahren hat der Adapter danach Shutdown gemeldet.
	nachEnde := func(t *testing.T, herunterfahren bool) *endeLauf {
		t.Helper()
		l := beginne(t, nil)
		client(t, l.s, l.id, gSync()...)
		l.up.letzte.empfang <- ende
		zustellen(t, l)
		if herunterfahren {
			l.s.Shutdown(context.Background(), l.id)
		}
		return l
	}
	for _, herunterfahren := range []bool{false, true} {
		t.Run(fmt.Sprintf("Anfrage, herunterfahren %v", herunterfahren), func(t *testing.T) {
			l := nachEnde(t, herunterfahren)
			if _, err := l.s.Query(context.Background(), l.id, "SELECT 2"); err != model.ErrSessionEnded {
				t.Fatalf("Query: %v", err)
			}
			if got := l.up.letzte.anfragen; !reflect.DeepEqual(got, []string{"SELECT 1"}) {
				t.Fatalf("Anfragen am Upstream %q", got)
			}
		})
		t.Run(fmt.Sprintf("Client-Nachricht, herunterfahren %v", herunterfahren), func(t *testing.T) {
			l := nachEnde(t, herunterfahren)
			if err := l.s.ClientMessage(context.Background(), l.id, cParse()); err != model.ErrSessionEnded {
				t.Fatalf("ClientMessage: %v", err)
			}
		})
	}
	t.Run("Lesen des Upstreams", func(t *testing.T) {
		l := nachEnde(t, false)
		l.up.letzte.empfang <- []model.Response{sRFQ()}
		if err := l.lies(); err != model.ErrSessionEnded || len(l.up.letzte.empfang) != 1 {
			t.Fatalf("AwaitServer: %v, wartende Folgen %d", err, len(l.up.letzte.empfang))
		}
	})
	t.Run("wartende Anfrage", func(t *testing.T) {
		l := beginne(t, nil)
		client(t, l.s, l.id, gSync()...)
		fertig := make(chan error, 1)
		go func() {
			_, err := l.s.Query(context.Background(), l.id, "SELECT 2")
			fertig <- err
		}()
		select {
		case err := <-fertig:
			t.Fatalf("Query kehrt vor dem Ende zurück: %v", err)
		case <-testzeit.Nach(100):
		}
		l.up.letzte.empfang <- ende
		zustellen(t, l)
		if err := <-fertig; err != model.ErrSessionEnded {
			t.Fatalf("wartende Query: %v", err)
		}
		if got := l.up.letzte.anfragen; !reflect.DeepEqual(got, []string{"SELECT 1"}) {
			t.Fatalf("Anfragen am Upstream %q", got)
		}
	})
	t.Run("Pipeline", func(t *testing.T) {
		l := beginne(t, nil)
		client(t, l.s, l.id, gSync()...)
		client(t, l.s, l.id, gSync()...)
		l.up.letzte.empfang <- ende
		zustellen(t, l)
		l.up.letzte.empfang <- []model.Response{sRFQ()}
		if err := l.lies(); err != model.ErrSessionEnded || len(l.up.letzte.empfang) != 1 {
			t.Fatalf("AwaitServer: %v, wartende Folgen %d", err, len(l.up.letzte.empfang))
		}
		if err := l.schliessen(model.EndServer); err != nil {
			t.Fatal(err)
		}
		if got := l.repo.last(t).Sessions[0].Interactions; len(got) != 2 {
			t.Fatalf("Interaktionen %#v", got)
		}
	})
}

// Abdeckung: LH-FA-25/Happy, LH-FA-13/Boundary — Record im Kern: Das Ende des Servers
// und das Zwangsende: Nach der Interaktion mit Ende läuft keine mehr, auch keine
// in die Pipeline gestellte; EndServer, EndClosed und EndForced melden keinen Fehler
// (kein PGR-E4003, kein PGR-E4006), die Session wird mit der Interaktion geschrieben
// und der Upstream geschlossen; Shutdown meldet mitten in der Interaktion nicht das Ende, Delivered nach dem Ende
// gibt es frei (LH-FA-25.a *Herunterfahren beim Aufzeichnen*).
func TestRecordEndeEnden(t *testing.T) {
	ende := []model.Response{endeFehler("FATAL", "FATAL"), sEnde()}
	for _, mitPipeline := range []bool{false, true} {
		for _, end := range []model.SessionEnd{model.EndServer, model.EndClosed, model.EndForced} {
			t.Run(fmt.Sprintf("Pipeline %v, Ende %d", mitPipeline, end), func(t *testing.T) {
				l := beginne(t, nil)
				client(t, l.s, l.id, gSync()...)
				if mitPipeline {
					client(t, l.s, l.id, gSync()...)
				}
				if l.s.Shutdown(context.Background(), l.id) {
					t.Fatal("Shutdown gibt das Ende mitten in der Interaktion frei")
				}
				l.up.letzte.empfang <- ende
				if err := l.lies(); err != nil {
					t.Fatal(err)
				}
				if !l.s.Delivered(context.Background(), l.id) {
					t.Error("Delivered gibt das Ende nach dem Ende des Servers nicht frei")
				}
				if err := l.schliessen(end); err != nil {
					t.Fatalf("Ende %d meldet %v", end, err)
				}
				if got := l.repo.last(t).Sessions[0].Interactions; len(got) != 2 {
					t.Fatalf("Interaktionen %#v", got)
				}
				select {
				case <-l.up.letzte.zu:
				default:
					t.Error("Upstream nach dem Ende nicht geschlossen")
				}
			})
		}
	}
	t.Run("Zwangsende vor dem Lesen", func(t *testing.T) {
		l := beginne(t, nil)
		client(t, l.s, l.id, gFlush()...)
		l.up.letzte.empfang <- ende
		if err := l.schliessen(model.EndForced); codeOf(err) != model.CodeShutdownTimeout {
			t.Fatalf("Zwangsende vor dem Lesen des Endes: %v", err)
		}
		if got := l.repo.last(t).Sessions[0].Interactions; len(got) != 1 {
			t.Fatalf("Interaktionen %#v", got)
		}
	})
	t.Run("einfache Anfrage", func(t *testing.T) {
		l := beginne(t, ende)
		l.up.letzte.queryLaeuft = make(chan struct{})
		l.up.letzte.queryHalt = make(chan struct{})
		fertig := make(chan error, 1)
		go func() {
			_, err := l.einfach()
			fertig <- err
		}()
		<-l.up.letzte.queryLaeuft
		if l.s.Shutdown(context.Background(), l.id) {
			t.Fatal("Shutdown gibt das Ende mitten in der Anfrage frei")
		}
		close(l.up.letzte.queryHalt)
		if err := <-fertig; err != nil || !l.s.Delivered(context.Background(), l.id) {
			t.Fatalf("Query: %v", err)
		}
	})
}

// Abdeckung: LH-FA-25/Negative, LH-FA-02/Negative — Record im Kern: Scheitert die
// Zustellung der Antworten mit Ende (EndWriteFailed), entfällt die Interaktion wie
// jede nicht zugestellte: PGR-E4003, die Session wird mit den vorherigen
// Interaktionen übernommen, die Aufzeichnung trägt kein Ende (LH-FA-25.a *Meldung des
// Endes beim Aufzeichnen*).
func TestRecordEndeSchreibfehler(t *testing.T) {
	rs := []model.Response{endeFehler("FATAL", "FATAL"), sEnde()}
	jeWeg(t, rs, func(t *testing.T, l *endeLauf, out []model.Response, err error) {
		t.Helper()
		if err != nil || len(out) != 2 {
			t.Fatalf("Antworten %#v, Fehler %v", out, err)
		}
		if err := l.schliessen(model.EndWriteFailed); codeOf(err) != model.CodeConnectionLost {
			t.Fatalf("Fehler %v, erwartet %s", err, model.CodeConnectionLost)
		}
		sess := l.repo.last(t).Sessions[0]
		if len(sess.Interactions) != 1 || sess.Interactions[0].Request.SQL != "SELECT 1" {
			t.Fatalf("Interaktionen %#v", sess.Interactions)
		}
	})
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-21/Boundary — Record im Kern: Mit Zeitaufzeichnung
// trägt die Interaktion mit Ende offset_ms (Unterschied der Lesungen der Uhr beim
// Eintreffen der Anfrage und beim Beginn der Session), und mit
// LeereSessions ändert das Ende nichts: Die Session wird einmal geschrieben, mit der
// Kennung 1 und der Kennzeichnung empty_sessions (LH-FA-21.a *Aufzeichnen*, LH-FA-12.a).
func TestRecordEndeZeitUndLeereSessions(t *testing.T) {
	rs := []model.Response{endeFehler("FATAL", "FATAL"), sEnde()}
	for _, weg := range []string{"einfach", "erweitert"} {
		t.Run(weg, func(t *testing.T) {
			uhr := testuhr.Neu()
			var l *endeLauf
			if weg == "einfach" {
				l = beginne(t, rs, services.Zeitangaben(uhr, true), services.LeereSessions(true))
				uhr.Vorruecken(25)
				if _, err := l.einfach(); err != nil {
					t.Fatal(err)
				}
			} else {
				l = beginne(t, nil, services.Zeitangaben(uhr, true), services.LeereSessions(true))
				uhr.Vorruecken(25)
				if _, err := l.erweitert(rs, gSync()); err != nil {
					t.Fatal(err)
				}
			}
			l.s.Delivered(context.Background(), l.id)
			if err := l.schliessen(model.EndServer); err != nil {
				t.Fatal(err)
			}
			rec := l.repo.last(t)
			if got := offsets(t, l.repo, 0); !gleichZahlen(got, []int64{0, 25}) || len(rec.Sessions) != 1 || rec.Sessions[0].ID != 1 || !rec.EmptySessions {
				t.Fatalf("offset_ms %v, Sessions %#v, empty_sessions %v", got, rec.Sessions, rec.EmptySessions)
			}
		})
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-12/Boundary — Record im Kern: Ist die Interaktion mit
// Ende die einzige der Session, wird die Session geschrieben, auch beim Zwangsende (kein
// Zwangsende im Aufbau, LH-FA-12.a *Angenommene Verbindung*), mit und ohne
// LeereSessions (eine Session mit Ende ist nie leer).
func TestRecordEndeEinzigeInteraktion(t *testing.T) {
	ende := []model.Response{endeFehler("FATAL", "FATAL"), sEnde()}
	for _, leere := range []bool{false, true} {
		for _, end := range []model.SessionEnd{model.EndServer, model.EndForced} {
			t.Run(fmt.Sprintf("leere %v, Ende %d", leere, end), func(t *testing.T) {
				l := oeffne(t, nil, nil, services.LeereSessions(leere))
				if _, err := l.erweitert(ende, gSync()); err != nil {
					t.Fatal(err)
				}
				l.s.Delivered(context.Background(), l.id)
				if err := l.schliessen(end); err != nil {
					t.Fatalf("Ende %d meldet %v", end, err)
				}
				if got := l.repo.last(t).Sessions; len(got) != 1 || got[0].ID != 1 || len(got[0].Interactions) != 1 {
					t.Fatalf("Sessions %#v", got)
				}
			})
		}
	}
}
