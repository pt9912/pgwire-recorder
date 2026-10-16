package services_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// parseDann ist eine Extended-Interaktion aus zwei Gruppen: Parse mit Flush
// (aufgezeichnet ParseComplete) und Bind, Execute mit schluss (Sync oder
// Flush), deren aufgezeichnete Server-Nachrichten letzte sind.
func parseDann(seq int, schluss model.ClientMessageType, letzte ...model.Response) model.Interaction {
	return extendedInteraktion(seq,
		gruppe([]model.Response{{Type: model.ResponseParseComplete}}, cm(model.ClientParse), cm(model.ClientFlush)),
		gruppe(letzte, cm(model.ClientBind), cm(model.ClientExecute), cm(schluss)),
	)
}

// endeExtended ist parseDann mit dem aufgezeichneten Ende in der letzten Gruppe:
// BindComplete, FATAL 57P01, connection_closed.
func endeExtended(seq int, schluss model.ClientMessageType) model.Interaction {
	return parseDann(seq, schluss, append([]model.Response{{Type: model.ResponseBindComplete}}, endeAufgezeichnet()...)...)
}

// endeExtRec ist eine Aufzeichnung aus zwei Sessions: Session 1 mit einer
// Extended-Interaktion nur mit Sync und danach in, Session 2 mit einer einfachen
// Anfrage C. Jeder Aufruf baut seine Aufzeichnung selbst.
func endeExtRec(in model.Interaction) model.Recording {
	rec := model.NewRecording()
	rec.Sessions = []model.Session{
		{ID: 1, Interactions: []model.Interaction{syncInteraktion(1), in}},
		{ID: 2, Interactions: []model.Interaction{aufgezeichnet(1, "C", befehl("SELECT 1"), bereitMit("I"))}},
	}
	return rec
}

// endeFall ist ein Lauf gegen den Fake: Aufzeichnung, Optionen, die Folgen der
// ersten Session je Gruppe (server), ob eine zweite Session mit der einfachen
// Anfrage C läuft (folge), die erwarteten Meldungen mit Exit-Code, der Ablauf
// der ersten Session, ob sie ohne Terminate endet (ohne) und die Zahl der
// eingespielten Interaktionen; verglichen und abweichend sind 0.
type endeFall struct {
	name        string
	rec         model.Recording
	o           services.PlayOptions
	server      [][]model.Response
	folge       bool
	want        []string
	exit        int
	ablauf      []string
	ohne        bool
	eingespielt int
}

// pruefeEndeExtended spielt f ein und prüft Meldungen, Ablauf, Schließen und Zähler.
func pruefeEndeExtended(t *testing.T, f endeFall) {
	t.Helper()
	sitz := spielerFolge(true, f.server...)
	sessions := []*fakeEinspiel{sitz}
	var zweite *fakeEinspiel
	ablauf := f.ablauf
	if f.folge {
		zweite = spielerFolge(false, []model.Response{befehl("SELECT 1"), bereitMit("I")})
		sessions = append(sessions, zweite)
		ablauf = append(append([]string(nil), ablauf...), "verbinde", "anfrage C", "naechste", "naechste", "schliesse")
	}
	ziel := &fakeZiel{sessions: sessions}
	z, err := spieleZaehler(context.Background(), t, f.rec, ziel, f.o)
	if got := texte(err); !reflect.DeepEqual(got, f.want) || exitVon(err) != f.exit {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit %d", got, exitVon(err), f.want, f.exit)
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, ablauf) {
		t.Fatalf("Ablauf %q, erwartet %q", got, ablauf)
	}
	if zu, ohne := sitz.war(); !zu || ohne != f.ohne {
		t.Fatalf("Session 1: geschlossen %v, ohne Terminate %v, erwartet geschlossen mit ohne Terminate %v", zu, ohne, f.ohne)
	}
	if zweite != nil {
		if zu, ohne := zweite.war(); !zu || ohne {
			t.Fatalf("Session 2: geschlossen %v, ohne Terminate %v, erwartet mit Terminate", zu, ohne)
		}
	}
	if want := (services.Zaehler{Eingespielt: f.eingespielt}); z != want {
		t.Fatalf("Zähler %+v, erwartet %+v (verglichen 0, abweichend 0)", z, want)
	}
}

// folgenEnde sind die Folgen der ersten Session für endeExtRec: die Antwort auf
// die Interaktion nur mit Sync, ParseComplete und die Antworten der letzten
// Gruppe, deren zweite die Fehlerantwort r ist (BindComplete davor).
func folgenEnde(r model.Response) [][]model.Response {
	return [][]model.Response{bereit(), {{Type: model.ResponseParseComplete}}, {{Type: model.ResponseBindComplete}, r}}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-25/Happy — ohne Vergleich ist eine
// Fehlerantwort FATAL oder PANIC (Feld V, ohne es S) in einer Extended-Interaktion
// mit --allow-recorded-errors und einem aufgezeichneten Ende in der letzten Gruppe
// das erwartete Ende, in einer Gruppe mit Sync wie in einer mit Flush und in einer
// früheren Gruppe mit Flush: keine Meldung, Exit-Code 0, auch ohne und mit
// --continue-on-error; play sendet danach weder weitere Gruppen noch Terminate und
// liest nichts mehr, die Session endet, die nächste läuft; eingespielt zählt die
// Interaktion, verglichen und abweichend bleiben 0 (LH-FA-20.a *Interaktion*,
// *Gruppen*, *Ende einer Session*, LH-FA-25.a).
func TestPlayEndeExtendedErwartet(t *testing.T) {
	kopf := []string{"verbinde", "gruppe sync", "naechste", "gruppe parse,flush", "naechste"}
	for _, f := range fatalFaelle() {
		for _, fortsetzen := range []bool{false, true} {
			o := services.PlayOptions{AllowRecordedErrors: true, ContinueOnError: fortsetzen}
			for _, v := range []endeFall{
				{
					name: "Sync", rec: endeExtRec(endeExtended(2, model.ClientSync)), server: folgenEnde(f.r),
					ablauf: append(append([]string(nil), kopf...), "gruppe bind,execute,sync", "naechste", "naechste", "schliesse"),
				},
				{
					name: "Flush", rec: endeExtRec(endeExtended(2, model.ClientFlush)), server: folgenEnde(f.r),
					ablauf: append(append([]string(nil), kopf...), "gruppe bind,execute,flush", "naechste", "naechste", "schliesse"),
				},
				{
					name: "frühere Gruppe", rec: endeExtRec(parseDann(2, model.ClientSync, endeAufgezeichnet()...)), server: [][]model.Response{bereit(), {f.r}},
					ablauf: []string{"verbinde", "gruppe sync", "naechste", "gruppe parse,flush", "naechste", "schliesse"},
				},
			} {
				v.o, v.folge, v.ohne, v.eingespielt = o, true, true, 3
				name := f.name + "/" + v.name
				if fortsetzen {
					name += " mit --continue-on-error"
				}
				t.Run(name, func(t *testing.T) { pruefeEndeExtended(t, v) })
			}
		}
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-25/Negative — ohne Vergleich ist eine
// Fehlerantwort FATAL oder PANIC in einer Extended-Interaktion PGR-E4003, sofort
// und auch mit --continue-on-error, die Verbindung ohne Terminate geschlossen,
// wenn eine Bedingung des erwarteten Endes fehlt: ohne --allow-recorded-errors
// (Gruppe mit Sync wie mit Flush), ohne aufgezeichnetes Ende in der letzten Gruppe
// (auch bei aufgezeichneter Fehlerantwort ERROR, auch bei Fehlerantwort des Servers
// in einer früheren Gruppe) oder bei Ende des Servers in einer früheren
// Interaktion als der aufgezeichneten; eine nächste Session läuft nicht
// (LH-FA-20.a *Interaktion*, Tabelle der Fehlerregeln).
func TestPlayEndeExtendedNichtErwartet(t *testing.T) {
	fatal := schwer("FATAL", "FATAL", "57P01")
	mitOption := services.PlayOptions{AllowRecordedErrors: true}
	fortsetzen := services.PlayOptions{AllowRecordedErrors: true, ContinueOnError: true}
	e4003A := "Netzwerk [PGR-E4003]: Session 1, Interaktion 1: Fehlerantwort des Servers 57P01 „geheim“, die Verbindung endet"
	letzteGruppe := []string{"verbinde", "gruppe sync", "naechste", "gruppe parse,flush", "naechste", "gruppe bind,execute,sync", "naechste", "naechste", "schliesse"}
	letzteFlush := []string{"verbinde", "gruppe sync", "naechste", "gruppe parse,flush", "naechste", "gruppe bind,execute,flush", "naechste", "naechste", "schliesse"}
	fruehere := []string{"verbinde", "gruppe sync", "naechste", "gruppe parse,flush", "naechste", "schliesse"}
	ohneEnde := func(letzte ...model.Response) model.Recording {
		return endeExtRec(parseDann(2, model.ClientSync, letzte...))
	}
	for _, f := range []endeFall{
		{name: "ohne Option, Sync", rec: endeExtRec(endeExtended(2, model.ClientSync)), server: folgenEnde(fatal), want: []string{e4003B}, ablauf: letzteGruppe},
		{name: "ohne Option mit --continue-on-error, Sync", rec: endeExtRec(endeExtended(2, model.ClientSync)), o: services.PlayOptions{ContinueOnError: true}, server: folgenEnde(fatal), want: []string{e4003B}, ablauf: letzteGruppe},
		{name: "ohne Option, Flush", rec: endeExtRec(endeExtended(2, model.ClientFlush)), server: folgenEnde(fatal), want: []string{e4003B}, ablauf: letzteFlush},
		{name: "ohne Ende an der Stelle", rec: ohneEnde(bereit()...), o: mitOption, server: folgenEnde(fatal), want: []string{e4003B}, ablauf: letzteGruppe},
		{
			name: "ohne Ende an der Stelle, aufgezeichnet ERROR", o: fortsetzen, server: folgenEnde(fatal), want: []string{e4003B}, ablauf: letzteGruppe,
			rec: ohneEnde(append([]model.Response{fehlerMit("42P01", "aufgezeichnet")}, bereit()...)...),
		},
		{
			name: "ohne Ende an der Stelle, Fehlerantwort in früherer Gruppe", o: mitOption, server: [][]model.Response{bereit(), {fatal}},
			rec: ohneEnde(bereit()...), want: []string{e4003B}, ablauf: fruehere,
		},
		{
			name: "Ende in früherer Interaktion", o: mitOption, server: [][]model.Response{{fatal}},
			rec: endeExtRec(endeExtended(2, model.ClientSync)), want: []string{e4003A}, ablauf: []string{"verbinde", "gruppe sync", "naechste", "schliesse"}, eingespielt: 1,
		},
	} {
		f.ohne, f.exit = true, 4
		if f.eingespielt == 0 {
			f.eingespielt = 2
		}
		t.Run(f.name, func(t *testing.T) { pruefeEndeExtended(t, f) })
	}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-25/Boundary — nach einer Fehlerantwort
// ERROR, die nicht abbricht, wird eine Fehlerantwort FATAL beim Lesen bis zum
// ReadyForQuery wie jede andere eingestuft: mit --allow-recorded-errors und einem
// aufgezeichneten Ende ist sie das erwartete Ende, play sendet die Gruppen nach
// dem ERROR ohne Warten dazwischen, nichts nach dem Ende (LH-FA-20.a *Interaktion*,
// *Gruppen*).
func TestPlayEndeExtendedNachError(t *testing.T) {
	rec := sitzung(extendedInteraktion(1,
		gruppe(nil, cm(model.ClientParse), cm(model.ClientFlush)),
		gruppe(nil, cm(model.ClientBind), cm(model.ClientFlush)),
		gruppe(endeAufgezeichnet(), cm(model.ClientExecute), cm(model.ClientSync)),
	))
	pruefeEndeExtended(t, endeFall{
		rec: rec, o: services.PlayOptions{AllowRecordedErrors: true},
		server: [][]model.Response{{fehlerMit("42P01", "x")}, nil, {schwer("FATAL", "FATAL", "57P01")}},
		ablauf: []string{"verbinde", "gruppe parse,flush", "naechste", "gruppe bind,flush", "gruppe execute,sync", "naechste", "schliesse"},
		ohne:   true, eingespielt: 1,
	})
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-25/Boundary — endet die letzte Gruppe einer
// Interaktion mit Ende mit Flush, kommt kein ReadyForQuery: play liest die
// Antworten dieser Gruppe nur bis zur Zahl, die sie erwartet, oder bis zur
// Fehlerantwort, ohne ein weiteres Lesen, und die Session endet mit Terminate;
// das gilt für eine normale Antwort, ein ERROR und ein ERROR in einer früheren
// Gruppe (dann liest play nichts mehr), mit und ohne --allow-recorded-errors
// (LH-FA-20.a *Gruppen*).
func TestPlayEndeExtendedLetzteFlush(t *testing.T) {
	in := func() model.Recording { return sitzung(endeExtended(1, model.ClientFlush)) }
	mit := services.PlayOptions{AllowRecordedErrors: true}
	ohneOption := services.PlayOptions{}
	normal := [][]model.Response{{{Type: model.ResponseParseComplete}}, {{Type: model.ResponseBindComplete}, befehl("SELECT 1")}}
	imFehler := [][]model.Response{{{Type: model.ResponseParseComplete}}, {fehlerMit("42P01", "x"), befehl("SELECT 1")}}
	nachFehler := [][]model.Response{{fehlerMit("42P01", "x")}, nil}
	e4004Fehler := []string{e4004("Session 1, Interaktion 1", "42P01", "x")}
	letzte := []string{"verbinde", "gruppe parse,flush", "naechste", "gruppe bind,execute,flush"}
	for _, f := range []endeFall{
		{name: "normale Antwort mit Option", o: mit, server: normal, ablauf: append(append([]string(nil), letzte...), "naechste", "naechste", "schliesse")},
		{name: "normale Antwort ohne Option", o: ohneOption, server: normal, ablauf: append(append([]string(nil), letzte...), "naechste", "naechste", "schliesse")},
		{name: "ERROR mit Option", o: mit, server: imFehler, ablauf: append(append([]string(nil), letzte...), "naechste", "schliesse")},
		{name: "ERROR ohne Option", o: ohneOption, server: imFehler, want: e4004Fehler, exit: 4, ablauf: append(append([]string(nil), letzte...), "naechste", "schliesse")},
		{name: "früherer ERROR mit Option", o: mit, server: nachFehler},
		{
			name: "früherer ERROR ohne Option mit --continue-on-error", o: services.PlayOptions{ContinueOnError: true},
			server: nachFehler, want: e4004Fehler, exit: 4,
		},
	} {
		f.rec, f.eingespielt = in(), 1
		if f.ablauf == nil {
			f.ablauf = append(append([]string(nil), letzte...), "schliesse")
		}
		t.Run(f.name, func(t *testing.T) { pruefeEndeExtended(t, f) })
	}
}
