package services_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// vform ist die Form der Interaktion B, an deren Stelle das Ende des Servers
// liegt: eine einfache Anfrage oder eine Extended-Interaktion aus der Gruppe
// Parse, Flush und einer letzten Gruppe Bind, Execute mit schluss. Die
// aufgezeichneten und die gelesenen Nachrichten der Fälle sind die der letzten
// Gruppe hinter BindComplete; vorspann zählt die Nachrichten davor in der
// normalisierten Folge (ParseComplete, BindComplete).
type vform struct {
	name     string
	vorspann int
	schluss  model.ClientMessageType
}

func formEinfach() vform { return vform{name: "einfach"} }

func formSync() vform { return vform{name: "Extended Sync", vorspann: 2, schluss: model.ClientSync} }

func formFlush() vform { return vform{name: "Extended Flush", vorspann: 2, schluss: model.ClientFlush} }

func alleFormen() []vform { return []vform{formEinfach(), formSync(), formFlush()} }

// ohneFlush sind die Formen, deren letzte Gruppe bis zum ReadyForQuery liest.
func ohneFlush() []vform { return []vform{formEinfach(), formSync()} }

// interaktion ist B mit den aufgezeichneten Antworten soll.
func (f vform) interaktion(soll []model.Response) model.Interaction {
	if f.schluss == "" {
		return aufgezeichnet(2, "B", soll...)
	}
	return parseDann(2, f.schluss, append([]model.Response{{Type: model.ResponseBindComplete}}, soll...)...)
}

// session ist die Session des Servers: A normal, dann B mit den Antworten ist
// und, mit verlust, einem Verbindungsverlust dahinter.
func (f vform) session(ist []model.Response, verlust bool) *fakeEinspiel {
	folge := schritte(ist...)
	if verlust {
		folge = append(folge, schritt{err: model.Errorf(model.CodeConnectionLost, nil, "weg")})
	}
	a := schritte(befehl("SELECT 1"), bereitMit("I"))
	if f.schluss == "" {
		return &fakeEinspiel{antworten: [][]schritt{a, folge}}
	}
	return &fakeEinspiel{antworten: [][]schritt{a}, gruppen: [][]schritt{
		schritte(model.Response{Type: model.ResponseParseComplete}),
		append(schritte(model.Response{Type: model.ResponseBindComplete}), folge...),
	}}
}

// ablauf ist der Ablauf der Session mit B: gelesen ist die Zahl der Nachrichten
// hinter BindComplete, die play liest.
func (f vform) ablauf(gelesen int) []string {
	a := []string{"verbinde", "anfrage A", "naechste", "naechste"}
	if f.schluss == "" {
		a = append(a, "anfrage B")
	} else {
		a = append(a, "gruppe parse,flush", "naechste", "gruppe bind,execute,"+string(f.schluss))
		gelesen++
	}
	for range gelesen {
		a = append(a, "naechste")
	}
	return append(a, "schliesse")
}

// vendeFall ist ein Fall des Endes im Vergleich: die aufgezeichneten Antworten
// soll und die des Servers ist von B, die Zahl der gelesenen Nachrichten von ist
// (gelesenFlush statt gelesen bei Extended Flush, wenn nicht 0), die erste
// Abweichung (position 0: keine; ab 1 in der einfachen Form, in Extended um
// vorspann verschoben) mit Art und Zusatz der Meldung, ob play die Verbindung
// ohne Terminate schließt (ohne) und ob der Server die Verbindung nach ist
// verliert (verlust).
type vendeFall struct {
	name         string
	formen       []vform
	soll, ist    []model.Response
	gelesen      int
	gelesenFlush int
	position     int
	art, zusatz  string
	ohne         bool
	verlust      bool
}

// dann ist r, gefolgt von den Antworten, die der Server nach der Fehlerantwort
// noch senden würde: play liest sie nicht.
func dann(r model.Response) []model.Response {
	return []model.Response{r, befehl("SELECT 1"), bereitMit("I")}
}

func fehlerOhneC(schwere string) model.Response {
	return model.Response{Type: model.ResponseErrorResponse, Fields: map[string]string{"S": schwere, "V": schwere, "M": "geheim"}}
}

func sqlstates(aufgezeichnet, server string) string {
	return fmt.Sprintf(" (SQLSTATE aufgezeichnet „%s“, Server „%s“)", aufgezeichnet, server)
}

// laufVergleichEnde spielt jeden Fall in jeder seiner Formen ein, ohne und mit
// --continue-on-error, ohne und mit --allow-recorded-errors.
func laufVergleichEnde(t *testing.T, faelle []vendeFall) {
	for _, c := range faelle {
		for _, f := range c.formen {
			for _, fortsetzen := range []bool{false, true} {
				for _, allow := range []bool{false, true} {
					o := services.PlayOptions{CompareResponses: true, ContinueOnError: fortsetzen, AllowRecordedErrors: allow}
					name := fmt.Sprintf("%s/%s/continue=%v/allow=%v", c.name, f.name, fortsetzen, allow)
					t.Run(name, func(t *testing.T) { pruefeVergleichEnde(t, c, f, o) })
				}
			}
		}
	}
}

func pruefeVergleichEnde(t *testing.T, c vendeFall, f vform, o services.PlayOptions) {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	abweichung := c.position > 0
	abbruch := c.verlust || abweichung && !o.ContinueOnError
	rec := model.NewRecording()
	rec.Sessions = []model.Session{
		{ID: 1, Interactions: []model.Interaction{aufgezeichnet(1, "A", ok...), f.interaktion(c.soll)}},
		{ID: 2, Interactions: []model.Interaction{aufgezeichnet(1, "C", ok...)}},
	}
	gelesen := c.gelesen
	if f.schluss == model.ClientFlush && c.gelesenFlush != 0 {
		gelesen = c.gelesenFlush
	}
	sitz := f.session(c.ist, c.verlust)
	sessions, ablauf := []*fakeEinspiel{sitz}, f.ablauf(gelesen)
	want := services.Zaehler{Eingespielt: 2, Verglichen: 2}
	var zweite *fakeEinspiel
	if !abbruch {
		zweite = spielerFolge(false, ok)
		sessions = append(sessions, zweite)
		ablauf = append(ablauf, "verbinde", "anfrage C", "naechste", "naechste", "schliesse")
		want = services.Zaehler{Eingespielt: 3, Verglichen: 3}
	}
	if c.verlust {
		want.Verglichen = 1
	}
	if abweichung {
		want.Abweichend = 1
	}
	ziel := &fakeZiel{sessions: sessions}
	z, err := spieleZaehler(t.Context(), t, rec, ziel, o)
	var meldungen []string
	exit := 0
	switch {
	case c.verlust:
		meldungen, exit = []string{"Netzwerk [PGR-E4003]: Session 1, Interaktion 2: weg"}, 4
	case abweichung:
		text := fmt.Sprintf("Session 1, Interaktion 2: Antwort weicht von der Aufzeichnung ab: Position %d, Art %s", c.position+f.vorspann, c.art)
		meldungen, exit = []string{e5004 + text + c.zusatz}, 5
	}
	if got := texte(err); !reflect.DeepEqual(got, meldungen) || exitVon(err) != exit {
		t.Errorf("Meldungen %q, Exit-Code %d, erwartet %q mit %d", got, exitVon(err), meldungen, exit)
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, ablauf) {
		t.Errorf("Ablauf %q, erwartet %q", got, ablauf)
	}
	if zu, ohne := sitz.war(); !zu || ohne != c.ohne {
		t.Errorf("Session 1: geschlossen %v, ohne Terminate %v, erwartet geschlossen mit ohne Terminate %v", zu, ohne, c.ohne)
	}
	if zweite != nil {
		if zu, ohne := zweite.war(); !zu || ohne {
			t.Errorf("Session 2: geschlossen %v, ohne Terminate %v, erwartet mit Terminate", zu, ohne)
		}
	}
	if z != want {
		t.Errorf("Zähler %+v, erwartet %+v", z, want)
	}
}

// Abdeckung: LH-FA-24/Happy, LH-FA-25/Happy — mit Vergleich ist eine
// Fehlerantwort FATAL oder PANIC (Feld V, ohne es S) des Servers, wenn die
// Aufzeichnung an der Stelle eine error_response FATAL oder PANIC mit
// connection_closed dahinter und denselben SQLSTATE (Feld C) zeigt, das erwartete
// Ende: keine Meldung, Exit-Code 0, keine Abweichung; play liest nichts mehr,
// schließt die Verbindung ohne Terminate, die Session endet und die nächste
// läuft, ohne und mit --continue-on-error und --allow-recorded-errors, bei
// einfacher Anfrage, Extended mit Sync und Extended mit Flush. Der Schweregrad
// beider Seiten ist gleich welcher (FATAL oder PANIC auf jeder Seite), zwei
// Fehlerantworten ohne Feld C sind gleich; verglichen zählt die Interaktion,
// abweichend nicht (LH-FA-24.a *Ende der Aufzeichnung*, *Randformen des
// Vergleichs*, LH-FA-20.a Tabelle).
func TestVergleichEndeErwartet(t *testing.T) {
	var faelle []vendeFall
	for _, f := range fatalFaelle() {
		faelle = append(faelle, vendeFall{name: "Server " + f.name, formen: alleFormen(), soll: endeAufgezeichnet(), ist: dann(f.r), gelesen: 1, ohne: true})
	}
	panikAufgezeichnet := []model.Response{schwer("PANIC", "PANIC", "57P01"), {Type: model.ResponseConnectionClosed}}
	ohneC := []model.Response{fehlerOhneC("FATAL"), {Type: model.ResponseConnectionClosed}}
	faelle = append(faelle,
		vendeFall{name: "aufgezeichnet PANIC", formen: alleFormen(), soll: panikAufgezeichnet, ist: dann(schwer("FATAL", "FATAL", "57P01")), gelesen: 1, ohne: true},
		vendeFall{name: "ohne Feld C", formen: alleFormen(), soll: ohneC, ist: dann(fehlerOhneC("PANIC")), gelesen: 1, ohne: true},
	)
	laufVergleichEnde(t, faelle)
}

// Abdeckung: LH-FA-24/Negative, LH-FA-24/Boundary — mit Vergleich ist eine
// Fehlerantwort FATAL oder PANIC des Servers die Abweichung PGR-E5004 der Art
// Fehler an ihrer Position, wenn eine Bedingung des erwarteten Endes fehlt:
// anderer SQLSTATE als das aufgezeichnete Ende; gleicher SQLSTATE wie eine
// aufgezeichnete Fehlerantwort ERROR an der Stelle, die kein Ende ist
// (Normalisierung vorher, mit Position 2); keine Fehlerantwort der Aufzeichnung an
// der Stelle (aufgezeichneter SQLSTATE leer), auch wenn die Aufzeichnung später
// ein Ende zeigt (Server früher). Eine frühere Abweichung geht vor, ebenso ein
// aufgezeichnetes Ende an früherer Position als das des Servers (Server später).
// Ohne --continue-on-error bricht der Lauf ab (Exit-Code 5), mit ihm endet die
// Session und die nächste läuft; --allow-recorded-errors ändert nichts;
// verglichen und abweichend zählen die Interaktion einmal (LH-FA-24.a *Ende der
// Aufzeichnung*, *Position*, *Art*, LH-FA-20.a Tabelle).
func TestVergleichEndeAbweichung(t *testing.T) {
	var faelle []vendeFall
	for _, f := range fatalFaelle() {
		faelle = append(faelle, vendeFall{
			name: "aufgezeichnet ERROR ohne Ende, Server " + f.name, formen: ohneFlush(), gelesen: 4,
			soll:     []model.Response{hinweis(), zeilenkopf(spalte("a", 23)), fehlerMit("57P01", "aufgezeichnet"), bereitMit("I")},
			ist:      append([]model.Response{parameter(), zeilenkopf(spalte("a", 23)), zeile()}, dann(f.r)...),
			position: 2, art: "Fehler", zusatz: sqlstates("57P01", "57P01"), ohne: true,
		})
	}
	nichtErwartet := func(name string, formen []vform, soll, ist []model.Response, gelesen, position int, art, zusatz string) vendeFall {
		return vendeFall{name: name, formen: formen, soll: soll, ist: ist, gelesen: gelesen, position: position, art: art, zusatz: zusatz, ohne: true}
	}
	ohneFehler := []model.Response{zeilenkopf(spalte("a", 23)), befehl("SELECT 1"), bereitMit("I")}
	fatal := schwer("FATAL", "FATAL", "57P01")
	faelle = append(faelle,
		nichtErwartet("anderer SQLSTATE", alleFormen(), endeAufgezeichnet(), dann(schwer("FATAL", "FATAL", "40P01")), 1, 1, "Fehler", sqlstates("57P01", "40P01")),
		nichtErwartet("ohne Fehlerantwort der Aufzeichnung", ohneFlush(), ohneFehler, append([]model.Response{zeilenkopf(spalte("a", 23))}, dann(fatal)...), 2, 2, "Fehler", sqlstates("", "57P01")),
		nichtErwartet("Server früher", alleFormen(), []model.Response{befehl("SELECT 1"), fatal, {Type: model.ResponseConnectionClosed}}, dann(fatal), 1, 1, "Fehler", sqlstates("", "57P01")),
		nichtErwartet("Server später", ohneFlush(), endeAufgezeichnet(), append([]model.Response{befehl("SELECT 1")}, dann(fatal)...), 2, 1, "Nachrichtenart", ""),
		nichtErwartet("frühere Abweichung", ohneFlush(), ohneFehler, append([]model.Response{zeilenkopf(spalte("b", 23))}, dann(fatal)...), 2, 1, "Spaltenbeschreibung", ""),
	)
	laufVergleichEnde(t, faelle)
}

// Abdeckung: LH-FA-24/Boundary — mit Vergleich vergleicht sich die Folge nach den
// gewöhnlichen Regeln, wenn der Server keine Fehlerantwort FATAL oder PANIC sendet:
// zeigt die Aufzeichnung ein Ende und antwortet der Server normal oder mit ERROR
// (gleicher oder anderer SQLSTATE), ist die Abweichung die erste verschiedene
// Position, bei ERROR mit gleichem SQLSTATE die des connection_closed, mit der Art
// Nachrichtenart; eine Fehlerantwort ERROR mit FATAL im Feld S und ERROR im Feld V
// ist keine FATAL-Fehlerantwort. Bei Extended Flush liest play nur bis zur Antwort
// auf jede Client-Nachricht der letzten Gruppe oder bis zur Fehlerantwort, ohne
// ReadyForQuery, und die Interaktion zählt als verglichen und abweichend
// (LH-FA-24.a *Ende der Aufzeichnung*, *Zähler*, LH-FA-20.a *Gruppen*).
func TestVergleichEndeGewoehnlich(t *testing.T) {
	ohneFehler := []model.Response{zeilenkopf(spalte("a", 23)), befehl("SELECT 1"), bereitMit("I")}
	laufVergleichEnde(t, []vendeFall{
		{
			name: "aufgezeichnetes Ende, Server ERROR gleicher SQLSTATE", formen: alleFormen(), soll: endeAufgezeichnet(),
			ist: []model.Response{fehlerMit("57P01", "x"), befehl("SELECT 1"), bereitMit("I")}, gelesen: 3, gelesenFlush: 1, position: 2, art: "Nachrichtenart",
		},
		{
			name: "aufgezeichnetes Ende, Server ERROR anderer SQLSTATE", formen: alleFormen(), soll: endeAufgezeichnet(),
			ist: []model.Response{fehlerMit("23505", "x"), bereitMit("I")}, gelesen: 2, gelesenFlush: 1, position: 1, art: "Fehler", zusatz: sqlstates("57P01", "23505"),
		},
		{
			name: "aufgezeichnetes Ende, normale Antwort", formen: alleFormen(), soll: endeAufgezeichnet(),
			ist: []model.Response{befehl("SELECT 1"), bereitMit("I")}, gelesen: 2, gelesenFlush: 1, position: 1, art: "Nachrichtenart",
		},
		{
			name: "V ERROR bei S FATAL", formen: ohneFlush(), soll: ohneFehler, gelesen: 4, position: 2, art: "Nachrichtenart",
			ist: []model.Response{zeilenkopf(spalte("a", 23)), schwer("FATAL", "ERROR", "42P01"), befehl("SELECT 1"), bereitMit("I")},
		},
	})
}

// Abdeckung: LH-FA-24/Negative — ein Verbindungsverlust ohne Fehlerantwort ist
// auch bei aufgezeichnetem Ende PGR-E4003 und keine Abweichung, ohne und mit
// --continue-on-error sofort, die Verbindung mit Terminate geschlossen; die
// Interaktion wird nicht verglichen und zählt nicht als verglichen; eine nächste
// Session läuft nicht (LH-FA-24.a *Ende der Aufzeichnung*, *Zähler*, LH-FA-20.a
// Tabelle).
func TestVergleichEndeVerbindungsverlust(t *testing.T) {
	laufVergleichEnde(t, []vendeFall{{name: "Verbindungsverlust", formen: alleFormen(), soll: endeAufgezeichnet(), gelesen: 1, verlust: true}})
}
