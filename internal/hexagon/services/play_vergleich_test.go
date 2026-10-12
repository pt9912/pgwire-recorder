package services_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// vergleichen sind die Optionen mit Vergleich.
func vergleichen() services.PlayOptions { return services.PlayOptions{CompareResponses: true} }

func spalte(name string, oid uint32) model.Column {
	return model.Column{Name: name, TypeOID: oid}
}

func zeilenkopf(spalten ...model.Column) model.Response {
	return model.Response{Type: model.ResponseRowDescription, Columns: spalten}
}

func zeile() model.Response {
	return model.Response{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte("1")}}}
}

func befehl(tag string) model.Response {
	return model.Response{Type: model.ResponseCommandComplete, Tag: tag}
}

func leer() model.Response { return model.Response{Type: model.ResponseEmptyQueryResponse} }

func hinweis() model.Response { return model.Response{Type: model.ResponseNoticeResponse} }

func parameter() model.Response { return model.Response{Type: model.ResponseParameterStatus} }

func bereitMit(tx string) model.Response {
	return model.Response{Type: model.ResponseReadyForQuery, TxStatus: tx}
}

// fehlerMit ist eine Fehlerantwort mit Schweregrad ERROR, SQLSTATE code und
// Meldung text.
func fehlerMit(code, text string) model.Response {
	return model.Response{Type: model.ResponseErrorResponse, Fields: map[string]string{"S": "ERROR", "V": "ERROR", "C": code, "M": text}}
}

// schwer ist eine Fehlerantwort mit den Feldern S und V und dem SQLSTATE code.
func schwer(s, v, code string) model.Response {
	felder := map[string]string{"C": code, "M": "geheim"}
	if s != "" {
		felder["S"] = s
	}
	if v != "" {
		felder["V"] = v
	}
	return model.Response{Type: model.ResponseErrorResponse, Fields: felder}
}

func schritte(rs ...model.Response) []schritt {
	var out []schritt
	for _, r := range rs {
		out = append(out, schritt{r: r})
	}
	return out
}

// aufgezeichnet ist eine einfache Anfrage sql mit den Antworten soll.
func aufgezeichnet(seq int, sql string, soll ...model.Response) model.Interaction {
	return model.Interaction{Sequence: seq, Request: model.Request{Type: model.RequestQuery, SQL: sql}, Responses: soll}
}

// nurSync ist eine Extended-Interaktion aus einer Gruppe nur mit Sync, deren
// aufgezeichnete Antworten soll sind.
func nurSync(seq int, soll ...model.Response) model.Interaction {
	return extendedInteraktion(seq, gruppe(soll, cm(model.ClientSync)))
}

// spielerFolge ist eine Session des Servers, die auf die erste Anfrage oder
// Gruppe ist liefert, auf die weiteren je die nächste Folge aus weitere.
func spielerFolge(extended bool, folgen ...[]model.Response) *fakeEinspiel {
	s := &fakeEinspiel{}
	for _, f := range folgen {
		if extended {
			s.gruppen = append(s.gruppen, schritte(f...))
		} else {
			s.antworten = append(s.antworten, schritte(f...))
		}
	}
	return s
}

// spieleZaehler spielt rec gegen ziel ein und liefert den Fehler und die Zähler.
func spieleZaehler(ctx context.Context, t *testing.T, rec model.Recording, ziel *fakeZiel, o services.PlayOptions) (services.Zaehler, error) {
	t.Helper()
	s, err := services.NewPlayService(context.WithoutCancel(ctx), ladeRepo{rec: rec}, "r.yaml", ziel, o)
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	err = s.Play(ctx, nil)
	return s.Zaehler(), err
}

// abweichungText ist die Meldung einer Abweichung an Session 1 und Interaktion 1
// ohne den Kopf.
func abweichungText(position int, art string) string {
	return fmt.Sprintf("Session 1, Interaktion 1: Antwort weicht von der Aufzeichnung ab: Position %d, Art %s", position, art)
}

const e5004 = "Replay [PGR-E5004]: "

// vfall ist ein Fall des Vergleichs einer Interaktion: die aufgezeichneten
// Antworten soll, die des Servers ist und die erste Abweichung (Position 0:
// keine) mit ihrer Art und dem Zusatz der Meldung.
type vfall struct {
	name        string
	soll, ist   []model.Response
	position    int
	art, zusatz string
	nurExtended bool
	nurEinfach  bool
}

func vfaelle() []vfall {
	rfq := bereitMit("I")
	basis := []model.Response{zeilenkopf(spalte("a", 23)), zeile(), befehl("SELECT 1"), rfq}
	return []vfall{
		{name: "gleiche Struktur", soll: basis, ist: basis},
		{name: "Zeilenwerte und Zeilenzahl", soll: basis, ist: []model.Response{zeilenkopf(spalte("a", 23)), zeile(), zeile(), zeile(), befehl("SELECT 3"), rfq}},
		{name: "keine Zeile", soll: basis, ist: []model.Response{zeilenkopf(spalte("a", 23)), befehl("SELECT 0"), rfq}},
		{name: "Hinweise und Parameterstatus nur beim Server", soll: basis, ist: []model.Response{hinweis(), zeilenkopf(spalte("a", 23)), parameter(), zeile(), hinweis(), befehl("SELECT 1"), rfq}},
		{name: "Hinweise und Parameterstatus nur in der Aufzeichnung", soll: []model.Response{hinweis(), parameter(), zeilenkopf(spalte("a", 23)), zeile(), befehl("SELECT 1"), hinweis(), rfq}, ist: basis},
		{name: "Nachrichtenart", soll: []model.Response{befehl("SELECT 1"), rfq}, ist: []model.Response{leer(), rfq}, position: 1, art: "Nachrichtenart"},
		{name: "fehlende Nachricht", soll: basis, ist: []model.Response{zeile(), befehl("SELECT 1"), rfq}, position: 1, art: "Nachrichtenart"},
		{name: "zusätzliche Nachricht", soll: []model.Response{befehl("SELECT 1"), rfq}, ist: []model.Response{befehl("SELECT 1"), befehl("SELECT 1"), rfq}, position: 2, art: "Nachrichtenart"},
		{name: "Position zählt in der normalisierten Folge", soll: []model.Response{hinweis(), zeilenkopf(spalte("a", 23)), zeile(), parameter(), befehl("SELECT 1"), rfq}, ist: []model.Response{zeilenkopf(spalte("a", 23)), zeile(), hinweis(), zeile(), befehl("UPDATE 1"), rfq}, position: 2, art: "Befehl"},
		{name: "Spaltenzahl", soll: basis, ist: []model.Response{zeilenkopf(spalte("a", 23), spalte("b", 23)), zeile(), befehl("SELECT 1"), rfq}, position: 1, art: "Spaltenbeschreibung"},
		{name: "Spaltenname", soll: basis, ist: []model.Response{zeilenkopf(spalte("b", 23)), zeile(), befehl("SELECT 1"), rfq}, position: 1, art: "Spaltenbeschreibung"},
		{name: "Typ-OID der Spalte", soll: basis, ist: []model.Response{zeilenkopf(spalte("a", 25)), zeile(), befehl("SELECT 1"), rfq}, position: 1, art: "Spaltenbeschreibung"},
		{name: "zweite Spalte", soll: []model.Response{zeilenkopf(spalte("a", 23), spalte("b", 23)), befehl("SELECT 0"), rfq}, ist: []model.Response{zeilenkopf(spalte("a", 23), spalte("b", 25)), befehl("SELECT 0"), rfq}, position: 1, art: "Spaltenbeschreibung"},
		{name: "übrige Felder der Spalte", soll: basis, ist: []model.Response{zeilenkopf(model.Column{Name: "a", TypeOID: 23, TableOID: 7, ColumnNumber: 3, TypeSize: 4, TypeModifier: 9, Format: 1}), zeile(), befehl("SELECT 1"), rfq}},
		{name: "Zahlen am Ende des Befehls", soll: []model.Response{befehl("INSERT 0 1"), rfq}, ist: []model.Response{befehl("INSERT 0 7"), rfq}},
		{name: "Befehl mit und ohne Zahl", soll: []model.Response{befehl("SELECT 5"), rfq}, ist: []model.Response{befehl("SELECT"), rfq}},
		{name: "Tag aus Ziffern allein", soll: []model.Response{befehl("5"), rfq}, ist: []model.Response{befehl(""), rfq}},
		{name: "anderer Befehl", soll: []model.Response{befehl("INSERT 0 1"), rfq}, ist: []model.Response{befehl("UPDATE 1"), rfq}, position: 1, art: "Befehl"},
		{name: "Wort des Befehls", soll: []model.Response{befehl("CREATE TABLE"), rfq}, ist: []model.Response{befehl("CREATE INDEX"), rfq}, position: 1, art: "Befehl"},
		{name: "Befehl kürzer um ein Wort", soll: []model.Response{befehl("CREATE TABLE"), rfq}, ist: []model.Response{befehl("CREATE"), rfq}, position: 1, art: "Befehl"},
		{name: "Wort mit Ziffer ist keine Zahl", soll: []model.Response{befehl("T1"), rfq}, ist: []model.Response{befehl("T"), rfq}, position: 1, art: "Befehl"},
		{name: "Zahl vor einem Wort bleibt", soll: []model.Response{befehl("A 1 B"), rfq}, ist: []model.Response{befehl("A 2 B"), rfq}, position: 1, art: "Befehl"},
		{name: "gleicher SQLSTATE, anderer Text und Felder", soll: []model.Response{{Type: model.ResponseErrorResponse, Fields: map[string]string{"S": "ERROR", "C": "42P01", "M": "a", "D": "x"}}, rfq}, ist: []model.Response{fehlerMit("42P01", "b"), rfq}},
		{name: "anderer Schweregrad, gleicher SQLSTATE", soll: []model.Response{schwer("WARNING", "WARNING", "42P01"), rfq}, ist: []model.Response{schwer("ERROR", "ERROR", "42P01"), rfq}},
		{name: "anderer SQLSTATE", soll: []model.Response{fehlerMit("42P01", "x"), rfq}, ist: []model.Response{fehlerMit("23505", "y"), rfq}, position: 1, art: "Fehler", zusatz: " (SQLSTATE aufgezeichnet „42P01“, Server „23505“)"},
		{name: "Fehler ohne Feld C auf beiden Seiten", soll: []model.Response{{Type: model.ResponseErrorResponse, Fields: map[string]string{"M": "a"}}, rfq}, ist: []model.Response{{Type: model.ResponseErrorResponse, Fields: map[string]string{"M": "b"}}, rfq}},
		{name: "Fehler ohne Feld C nur beim Server", soll: []model.Response{fehlerMit("42P01", "x"), rfq}, ist: []model.Response{{Type: model.ResponseErrorResponse, Fields: map[string]string{"M": "b"}}, rfq}, position: 1, art: "Fehler", zusatz: " (SQLSTATE aufgezeichnet „42P01“, Server „“)"},
		{name: "Fehler nur beim Server", soll: []model.Response{befehl("SELECT 1"), rfq}, ist: []model.Response{fehlerMit("42P01", "geheim"), rfq}, position: 1, art: "Nachrichtenart"},
		{name: "Aufzeichnung länger, Anfang gleich", nurEinfach: true, soll: []model.Response{befehl("SELECT 1"), rfq, befehl("SELECT 1"), rfq}, ist: []model.Response{befehl("SELECT 1"), rfq}, position: 3, art: "Nachrichtenart"},
		{name: "Fehler nur in der Aufzeichnung", soll: []model.Response{fehlerMit("42P01", "geheim"), rfq}, ist: []model.Response{befehl("SELECT 1"), rfq}, position: 1, art: "Nachrichtenart"},
		{name: "Transaktionsstatus", soll: []model.Response{befehl("BEGIN"), bereitMit("I")}, ist: []model.Response{befehl("BEGIN"), bereitMit("T")}, position: 2, art: "Transaktionsstatus"},
		{name: "erste von mehreren Abweichungen", soll: []model.Response{zeilenkopf(spalte("a", 23)), befehl("SELECT 1"), bereitMit("I")}, ist: []model.Response{zeilenkopf(spalte("b", 23)), befehl("UPDATE 1"), bereitMit("T")}, position: 1, art: "Spaltenbeschreibung"},
		{name: "Parameterbeschreibung gleich", nurExtended: true, soll: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseParameterDescription, ParamTypes: []uint32{23, 25}}, rfq}, ist: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseParameterDescription, ParamTypes: []uint32{23, 25}}, rfq}},
		{name: "Parameterbeschreibung, OID", nurExtended: true, soll: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseParameterDescription, ParamTypes: []uint32{23, 25}}, rfq}, ist: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseParameterDescription, ParamTypes: []uint32{23, 26}}, rfq}, position: 2, art: "Spaltenbeschreibung"},
		{name: "Parameterbeschreibung, Anzahl", nurExtended: true, soll: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseParameterDescription, ParamTypes: []uint32{23, 25}}, rfq}, ist: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseParameterDescription, ParamTypes: []uint32{23}}, rfq}, position: 2, art: "Spaltenbeschreibung"},
		{name: "Nachrichten ohne Felder", nurExtended: true, soll: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete}, rfq}, ist: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete}, rfq}},
		{name: "Nachrichtenart in Extended", nurExtended: true, soll: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete}, rfq}, ist: []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseCloseComplete}, rfq}, position: 2, art: "Nachrichtenart"},
	}
}

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Boundary, LH-FA-24/Negative — mit
// Vergleich endet das Einspielen einer einfachen Anfrage und einer
// Extended-Interaktion bei gleicher Antwortstruktur ohne Fehler; Unterschiede
// nur in Zeilenwerten, Zeilenzahl, Hinweisen, Parameterstatus, den
// abschließenden Zahlen des Befehls, dem Text und den weiteren Feldern einer
// Fehlerantwort, ihrem Schweregrad und den übrigen Feldern einer Spalte sind
// keine Abweichung; jede andere Art ist PGR-E5004 mit Session, Sequenznummer,
// Position ab 1 in der normalisierten Folge und Art (Nachrichtenart,
// Spaltenbeschreibung, Befehl, Fehler, Transaktionsstatus), bei einer Fehlerantwort mit dem
// SQLSTATE beider Seiten und sonst ohne Wert; Exit-Code 5; je Interaktion die erste
// Abweichung (LH-FA-24.a).
func TestVergleichStruktur(t *testing.T) {
	for _, f := range vfaelle() {
		for _, extended := range []bool{false, true} {
			if (f.nurExtended && !extended) || (f.nurEinfach && extended) {
				continue
			}
			t.Run(fmt.Sprintf("%s/extended=%v", f.name, extended), func(t *testing.T) {
				in := aufgezeichnet(1, "A", f.soll...)
				if extended {
					in = nurSync(1, f.soll...)
				}
				rec := sitzung(in)
				ziel := &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(extended, f.ist)}}
				z, err := spieleZaehler(context.Background(), t, rec, ziel, vergleichen())
				abw := 0
				if f.position > 0 {
					abw = 1
				}
				if want := (services.Zaehler{Eingespielt: 1, Verglichen: 1, Abweichend: abw}); z != want {
					t.Errorf("Zähler %+v, erwartet %+v", z, want)
				}
				if f.position == 0 {
					if err != nil {
						t.Fatalf("erwartet ohne Abweichung, erhalten %v", err)
					}
					return
				}
				want := e5004 + abweichungText(f.position, f.art) + f.zusatz
				if got := texte(err); !reflect.DeepEqual(got, []string{want}) || exitVon(err) != 5 || codeVon(err) != model.CodeResponseMismatch {
					t.Fatalf("Meldungen %q, Exit-Code %d, erwartet [%q] mit 5", got, exitVon(err), want)
				}
			})
		}
	}
}

// Abdeckung: LH-FA-24/Negative — bei einer Abweichung liest play die
// Interaktion bis zu ihrem ReadyForQuery und bricht dann vor der nächsten
// Interaktion ab, auch ohne --continue-on-error; die Abweichung ist ein Fehler
// PGR-E5004 mit Exit-Code 5 (LH-FA-20.a Schritt 6, Tabelle; LH-FA-24.a).
func TestVergleichAbbruchNachInteraktion(t *testing.T) {
	rec := sitzung(aufgezeichnet(1, "A", befehl("SELECT 1"), bereitMit("I")), aufgezeichnet(2, "B", befehl("SELECT 1"), bereitMit("I")))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(false, []model.Response{zeilenkopf(spalte("a", 23)), zeile(), befehl("SELECT 1"), bereitMit("I")})}}
	z, err := spieleZaehler(context.Background(), t, rec, ziel, vergleichen())
	if got, want := texte(err), []string{e5004 + abweichungText(1, "Nachrichtenart")}; !reflect.DeepEqual(got, want) || exitVon(err) != 5 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 5", got, exitVon(err), want)
	}
	if got, want := ziel.liste(), []string{"verbinde", "anfrage A", "naechste", "naechste", "naechste", "naechste", "schliesse"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
	if want := (services.Zaehler{Eingespielt: 1, Verglichen: 1, Abweichend: 1}); z != want {
		t.Fatalf("Zähler %+v, erwartet %+v", z, want)
	}
}

// Abdeckung: LH-FA-24/Boundary — mit --continue-on-error läuft das Einspielen
// nach einer Abweichung mit der nächsten Interaktion und der nächsten Session
// weiter, je Interaktion mit Abweichung eine Meldung PGR-E5004 in der
// Reihenfolge des Auftretens, Exit-Code 5; eingespielt, verglichen und
// abweichend zählen nach den Regeln von LH-FA-24.a (LH-FA-20.a Schritt 6,
// Tabelle).
func TestVergleichFortsetzung(t *testing.T) {
	gut := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	rec := model.NewRecording()
	rec.Sessions = []model.Session{
		{ID: 1, Interactions: []model.Interaction{aufgezeichnet(1, "A", gut...), aufgezeichnet(2, "B", gut...), aufgezeichnet(3, "C", gut...)}},
		{ID: 2, Interactions: []model.Interaction{aufgezeichnet(1, "D", gut...)}},
	}
	schlecht := []model.Response{befehl("UPDATE 1"), bereitMit("I")}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{
		spielerFolge(false, schlecht, gut, []model.Response{befehl("SELECT 1"), bereitMit("T")}),
		spielerFolge(false, schlecht),
	}}
	z, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{CompareResponses: true, ContinueOnError: true})
	want := []string{
		e5004 + abweichungText(1, "Befehl"),
		e5004 + "Session 1, Interaktion 3: Antwort weicht von der Aufzeichnung ab: Position 2, Art Transaktionsstatus",
		e5004 + "Session 2, Interaktion 1: Antwort weicht von der Aufzeichnung ab: Position 1, Art Befehl",
	}
	if got := texte(err); !reflect.DeepEqual(got, want) || exitVon(err) != 5 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 5", got, exitVon(err), want)
	}
	if z != (services.Zaehler{Eingespielt: 4, Verglichen: 4, Abweichend: 3}) {
		t.Fatalf("Zähler %+v", z)
	}
	ablauf := []string{
		"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "naechste", "anfrage C", "naechste", "naechste", "schliesse",
		"verbinde", "anfrage D", "naechste", "naechste", "schliesse",
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, ablauf) {
		t.Fatalf("Ablauf %q, erwartet %q", got, ablauf)
	}
}

// Abdeckung: LH-FA-24/Boundary — ohne --compare-responses findet kein Vergleich
// statt: Eine abweichende Struktur ist kein Fehler, eingespielt zählt, verglichen
// und abweichend bleiben 0; ein Lauf ohne Interaktion zählt 0 (LH-FA-24.a
// *Meldung*).
func TestVergleichOhneOption(t *testing.T) {
	rec := sitzung(aufgezeichnet(1, "A", befehl("SELECT 1"), bereitMit("I")), aufgezeichnet(2, "B", befehl("SELECT 1"), bereitMit("I")))
	abweichend := []model.Response{zeilenkopf(spalte("a", 23)), befehl("UPDATE 1"), bereitMit("T")}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(false, abweichend, abweichend)}}
	z, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{})
	if err != nil || z != (services.Zaehler{Eingespielt: 2}) {
		t.Fatalf("%v, Zähler %+v, erwartet nil mit {2 0 0}", err, z)
	}
	z, err = spieleZaehler(context.Background(), t, model.NewRecording(), &fakeZiel{}, vergleichen())
	if err != nil || z != (services.Zaehler{}) {
		t.Fatalf("Lauf ohne Interaktion: %v, Zähler %+v", err, z)
	}
}

// Abdeckung: LH-FA-24/Negative — mit Vergleich ist ein Serverfehler ohne Vorbild
// in der Aufzeichnung PGR-E5004 und nicht PGR-E4004, mit und ohne
// --continue-on-error; kein PGR-E4004 entsteht (LH-FA-24.a, LH-FA-20.a Tabelle).
func TestVergleichFehlerOhneVorbild(t *testing.T) {
	for _, weiter := range []bool{false, true} {
		rec := sitzung(aufgezeichnet(1, "A", befehl("SELECT 1"), bereitMit("I")), aufgezeichnet(2, "B", befehl("SELECT 1"), bereitMit("I")))
		ziel := &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(false, []model.Response{fehlerMit("42P01", "geheim"), bereitMit("I")}, []model.Response{befehl("SELECT 1"), bereitMit("I")})}}
		_, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{CompareResponses: true, ContinueOnError: weiter})
		if want := []string{e5004 + abweichungText(1, "Nachrichtenart")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 5 {
			t.Fatalf("mit continue=%v: Meldungen %q, Exit-Code %d, erwartet %q mit 5", weiter, texte(err), exitVon(err), want)
		}
	}
}

// Abdeckung: LH-FA-24/Boundary — mit Vergleich gilt eine Fehlerantwort mit
// gleichem SQLSTATE an der aufgezeichneten Stelle als erwartet: keine Meldung,
// kein Abbruch, die nächste Interaktion läuft, auch ohne --allow-recorded-errors;
// ein anderer SQLSTATE ist PGR-E5004, auch mit --allow-recorded-errors, das
// mit Vergleich ohne Wirkung ist (LH-FA-24.a, LH-FA-20.a Optionstabelle).
func TestVergleichErwarteterFehler(t *testing.T) {
	mitFehler := []model.Response{fehlerMit("42P01", "aufgezeichnet"), bereitMit("I")}
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	for _, allow := range []bool{false, true} {
		o := services.PlayOptions{CompareResponses: true, AllowRecordedErrors: allow}
		rec := sitzung(aufgezeichnet(1, "A", mitFehler...), aufgezeichnet(2, "B", ok...))
		ziel := &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(false, []model.Response{fehlerMit("42P01", "anderer Text"), bereitMit("I")}, ok)}}
		z, err := spieleZaehler(context.Background(), t, rec, ziel, o)
		if err != nil || z != (services.Zaehler{Eingespielt: 2, Verglichen: 2}) {
			t.Errorf("allow=%v, gleicher SQLSTATE: %v, Zähler %+v", allow, err, z)
		}
		if got := ziel.liste(); !strings.Contains(strings.Join(got, ","), "anfrage B") {
			t.Errorf("allow=%v: nächste Interaktion nicht gesendet: %q", allow, got)
		}

		ziel = &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(false, []model.Response{fehlerMit("23505", "x"), bereitMit("I")}, ok)}}
		_, err = spieleZaehler(context.Background(), t, rec, ziel, o)
		want := []string{e5004 + abweichungText(1, "Fehler") + " (SQLSTATE aufgezeichnet „42P01“, Server „23505“)"}
		if !reflect.DeepEqual(texte(err), want) || exitVon(err) != 5 {
			t.Errorf("allow=%v, anderer SQLSTATE: Meldungen %q, Exit-Code %d, erwartet %q mit 5", allow, texte(err), exitVon(err), want)
		}
	}
	rec := sitzung(aufgezeichnet(1, "A", ok...))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(false, []model.Response{fehlerMit("42P01", "x"), bereitMit("I")})}}
	_, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{CompareResponses: true, AllowRecordedErrors: true})
	if codeVon(err) != model.CodeResponseMismatch {
		t.Errorf("Fehler ohne aufgezeichnete Fehlerantwort mit --allow-recorded-errors: %v", err)
	}
}

// fatalFall ist eine Fehlerantwort FATAL oder PANIC und die Felder, die den
// Schweregrad nennen.
type fatalFall struct {
	name string
	r    model.Response
}

func fatalFaelle() []fatalFall {
	return []fatalFall{
		{"V FATAL", schwer("SCHWER", "FATAL", "57P01")},
		{"S FATAL ohne V", schwer("FATAL", "", "57P01")},
		{"V PANIC", schwer("x", "PANIC", "57P01")},
		{"S PANIC ohne V", schwer("PANIC", "", "57P01")},
	}
}

// Abdeckung: LH-FA-24/Negative, LH-FA-24/Boundary — mit Vergleich ist eine
// Fehlerantwort FATAL oder PANIC (Feld V, ohne es S) sofort PGR-E5004 und nicht
// PGR-E4003, auch bei gleichem SQLSTATE wie die aufgezeichnete: Art Fehler mit
// dem SQLSTATE beider Seiten, an ihrer Position in der normalisierten Folge;
// play liest nichts mehr, sendet nichts mehr und schließt die Verbindung, die
// Interaktion zählt als verglichen und abweichend; Exit-Code 5 (LH-FA-24.a
// *Ende der Aufzeichnung*, LH-FA-20.a Tabelle).
func TestVergleichFatal(t *testing.T) {
	for _, f := range fatalFaelle() {
		for _, extended := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/extended=%v", f.name, extended), func(t *testing.T) {
				soll := []model.Response{hinweis(), zeilenkopf(spalte("a", 23)), fehlerMit("57P01", "aufgezeichnet"), bereitMit("I")}
				in := aufgezeichnet(1, "A", soll...)
				naechsteB := aufgezeichnet(2, "B", befehl("SELECT 1"), bereitMit("I"))
				if extended {
					in, naechsteB = nurSync(1, soll...), nurSync(2, befehl("SELECT 1"), bereitMit("I"))
				}
				rec := sitzung(in, naechsteB)
				ziel := &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(extended, []model.Response{parameter(), zeilenkopf(spalte("a", 23)), zeile(), f.r, befehl("SELECT 1"), bereitMit("I")})}}
				z, err := spieleZaehler(context.Background(), t, rec, ziel, vergleichen())
				want := e5004 + abweichungText(2, "Fehler") + " (SQLSTATE aufgezeichnet „57P01“, Server „57P01“)"
				if got := texte(err); !reflect.DeepEqual(got, []string{want}) || exitVon(err) != 5 {
					t.Fatalf("Meldungen %q, Exit-Code %d, erwartet [%q] mit 5", got, exitVon(err), want)
				}
				senden := "anfrage A"
				if extended {
					senden = "gruppe sync"
				}
				if got, wantAblauf := ziel.liste(), []string{"verbinde", senden, "naechste", "naechste", "naechste", "naechste", "schliesse"}; !reflect.DeepEqual(got, wantAblauf) {
					t.Fatalf("Ablauf %q, erwartet %q", got, wantAblauf)
				}
				if z != (services.Zaehler{Eingespielt: 1, Verglichen: 1, Abweichend: 1}) {
					t.Fatalf("Zähler %+v", z)
				}
			})
		}
	}
}

// Abdeckung: LH-FA-24/Boundary — eine Fehlerantwort mit Schweregrad ERROR in V
// bei FATAL in S ist keine FATAL-Fehlerantwort, und eine FATAL-Fehlerantwort an
// einer Position, an der die Aufzeichnung keine Fehlerantwort hat, ist stets
// eine Abweichung der Art Fehler, ohne SQLSTATE der Aufzeichnung; eine frühere
// Abweichung der Interaktion geht ihr vor (LH-FA-24.a *Ende der Aufzeichnung*,
// *Position*, *Art*).
func TestVergleichFatalRand(t *testing.T) {
	soll := aufgezeichnet(1, "A", zeilenkopf(spalte("a", 23)), befehl("SELECT 1"), bereitMit("I"))
	sitz := func(ist ...model.Response) (services.Zaehler, error) {
		ziel := &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(false, ist)}}
		return spieleZaehler(context.Background(), t, sitzung(soll), ziel, vergleichen())
	}
	z, err := sitz(zeilenkopf(spalte("a", 23)), schwer("FATAL", "ERROR", "42P01"), befehl("SELECT 1"), bereitMit("I"))
	if want := []string{e5004 + abweichungText(2, "Nachrichtenart")}; !reflect.DeepEqual(texte(err), want) || z.Abweichend != 1 {
		t.Errorf("V ERROR vor S FATAL: Meldungen %q, Zähler %+v, erwartet %q", texte(err), z, want)
	}
	_, err = sitz(zeilenkopf(spalte("a", 23)), schwer("FATAL", "FATAL", "57P01"))
	if want := []string{e5004 + abweichungText(2, "Fehler") + " (SQLSTATE aufgezeichnet „“, Server „57P01“)"}; !reflect.DeepEqual(texte(err), want) {
		t.Errorf("FATAL an der Stelle von command_complete: Meldungen %q, erwartet %q", texte(err), want)
	}
	_, err = sitz(zeilenkopf(spalte("b", 23)), schwer("FATAL", "FATAL", "57P01"))
	if want := []string{e5004 + abweichungText(1, "Spaltenbeschreibung")}; !reflect.DeepEqual(texte(err), want) {
		t.Errorf("frühere Abweichung vor FATAL: Meldungen %q, erwartet %q", texte(err), want)
	}
}

// Abdeckung: LH-FA-24/Boundary — mit Vergleich und --continue-on-error endet
// die Session bei einer Fehlerantwort FATAL: play sendet ihr nichts mehr, die
// nächste Session läuft; die Meldungen stehen in der Reihenfolge des
// Auftretens, Exit-Code 5 (LH-FA-20.a Tabelle).
func TestVergleichFatalFortsetzung(t *testing.T) {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	rec := model.NewRecording()
	rec.Sessions = []model.Session{
		{ID: 1, Interactions: []model.Interaction{aufgezeichnet(1, "A", ok...), aufgezeichnet(2, "B", ok...)}},
		{ID: 2, Interactions: []model.Interaction{aufgezeichnet(1, "C", ok...)}},
	}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{
		spielerFolge(false, []model.Response{schwer("FATAL", "FATAL", "57P01"), bereitMit("I")}),
		spielerFolge(false, []model.Response{befehl("UPDATE 1"), bereitMit("I")}),
	}}
	z, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{CompareResponses: true, ContinueOnError: true})
	want := []string{
		e5004 + abweichungText(1, "Fehler") + " (SQLSTATE aufgezeichnet „“, Server „57P01“)",
		e5004 + "Session 2, Interaktion 1: Antwort weicht von der Aufzeichnung ab: Position 1, Art Befehl",
	}
	if got := texte(err); !reflect.DeepEqual(got, want) || exitVon(err) != 5 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 5", got, exitVon(err), want)
	}
	ablauf := []string{"verbinde", "anfrage A", "naechste", "schliesse", "verbinde", "anfrage C", "naechste", "naechste", "schliesse"}
	if got := ziel.liste(); !reflect.DeepEqual(got, ablauf) {
		t.Fatalf("Ablauf %q, erwartet %q", got, ablauf)
	}
	if z != (services.Zaehler{Eingespielt: 2, Verglichen: 2, Abweichend: 2}) {
		t.Fatalf("Zähler %+v", z)
	}
}

// Abdeckung: LH-FA-24/Negative — ein Verbindungsverlust ohne Fehlerantwort ist
// PGR-E4003 und keine Abweichung, auch mit --continue-on-error sofort; die
// unterbrochene Interaktion zählt nicht als verglichen; folgt er auf eine
// frühere Abweichung, steht seine Meldung zuerst und der Exit-Code ist 4
// (LH-FA-20.a Tabelle, Exit-Code).
func TestVergleichVerbindungsverlust(t *testing.T) {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	verlust := schritt{err: model.Errorf(model.CodeConnectionLost, nil, "weg")}
	rec := sitzung(aufgezeichnet(1, "A", ok...), aufgezeichnet(2, "B", ok...), aufgezeichnet(3, "C", ok...))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{schritte(befehl("UPDATE 1"), bereitMit("I")), append(schritte(zeilenkopf(spalte("a", 23))), verlust)}}}}
	z, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{CompareResponses: true, ContinueOnError: true})
	want := []string{"Netzwerk [PGR-E4003]: Session 1, Interaktion 2: weg", e5004 + abweichungText(1, "Befehl")}
	if got := texte(err); !reflect.DeepEqual(got, want) || exitVon(err) != 4 {
		t.Errorf("Meldungen %q, Exit-Code %d, erwartet %q mit 4", got, exitVon(err), want)
	}
	if got, wantAblauf := ziel.liste(), []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "naechste", "schliesse"}; !reflect.DeepEqual(got, wantAblauf) {
		t.Errorf("Ablauf %q, erwartet %q", got, wantAblauf)
	}
	if z != (services.Zaehler{Eingespielt: 2, Verglichen: 1, Abweichend: 1}) {
		t.Errorf("Zähler %+v, erwartet {2 1 1}", z)
	}
	for _, weiter := range []bool{false, true} {
		ziel := &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{append(schritte(befehl("SELECT 1")), verlust)}}}}
		z, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{CompareResponses: true, ContinueOnError: weiter})
		if want := []string{"Netzwerk [PGR-E4003]: Session 1, Interaktion 1: weg"}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 4 {
			t.Errorf("continue=%v: Meldungen %q, Exit-Code %d, erwartet %q mit 4", weiter, texte(err), exitVon(err), want)
		}
		if z != (services.Zaehler{Eingespielt: 1}) {
			t.Errorf("continue=%v: Zähler %+v, erwartet {1 0 0}", weiter, z)
		}
	}
}

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Negative — bei einer Extended-Interaktion
// ist die Folge aller Server-Nachrichten ihrer Gruppen ohne die Gruppengrenzen
// zu vergleichen: Liegen die Grenzen der Aufzeichnung anders als die der
// Ankunft, ist das keine Abweichung; eine Abweichung in einer späteren Gruppe
// wird mit ihrer Position in der ganzen Folge gemeldet, auch bei Gruppen mit
// Flush (LH-FA-24.a, LH-FA-20.a *Gruppen*).
func TestVergleichExtendedGruppen(t *testing.T) {
	g1 := gruppe([]model.Response{{Type: model.ResponseParseComplete}}, cm(model.ClientParse), beschreibe(model.TargetStatement), cm(model.ClientFlush))
	g2 := gruppe([]model.Response{
		{Type: model.ResponseParameterDescription, ParamTypes: []uint32{23}}, {Type: model.ResponseNoData}, {Type: model.ResponseBindComplete}, befehl("INSERT 0 1"), bereitMit("I"),
	}, cm(model.ClientBind), cm(model.ClientExecute), cm(model.ClientSync))
	rec := sitzung(extendedInteraktion(1, g1, g2))
	gruppen := func(oid uint32, tag string) *fakeEinspiel {
		return &fakeEinspiel{gruppen: [][]schritt{
			schritte(model.Response{Type: model.ResponseParseComplete}, model.Response{Type: model.ResponseParameterDescription, ParamTypes: []uint32{oid}}, model.Response{Type: model.ResponseNoData}),
			schritte(model.Response{Type: model.ResponseBindComplete}, befehl(tag), bereitMit("I")),
		}}
	}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{gruppen(23, "INSERT 0 5")}}
	z, err := spieleZaehler(context.Background(), t, rec, ziel, vergleichen())
	if err != nil || z != (services.Zaehler{Eingespielt: 1, Verglichen: 1}) {
		t.Fatalf("andere Gruppengrenzen: %v, Zähler %+v", err, z)
	}
	want := []string{"verbinde", "gruppe parse,describe,flush", "naechste", "naechste", "naechste", "gruppe bind,execute,sync", "naechste", "naechste", "naechste", "schliesse"}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
	ziel = &fakeZiel{sessions: []*fakeEinspiel{gruppen(25, "INSERT 0 5")}}
	_, err = spieleZaehler(context.Background(), t, rec, ziel, vergleichen())
	if wantText := []string{e5004 + abweichungText(2, "Spaltenbeschreibung")}; !reflect.DeepEqual(texte(err), wantText) {
		t.Fatalf("Abweichung in der ersten Gruppe: Meldungen %q, erwartet %q", texte(err), wantText)
	}
	ziel = &fakeZiel{sessions: []*fakeEinspiel{gruppen(23, "UPDATE 1")}}
	_, err = spieleZaehler(context.Background(), t, rec, ziel, vergleichen())
	if wantText := []string{e5004 + abweichungText(5, "Befehl")}; !reflect.DeepEqual(texte(err), wantText) {
		t.Fatalf("Abweichung in der zweiten Gruppe: Meldungen %q, erwartet %q", texte(err), wantText)
	}
}

// Abdeckung: LH-FA-24/Boundary — ein Fehler des Servers mitten in einer
// Extended-Interaktion ändert den Vergleich nicht: play sendet auch mit
// Vergleich nach der Fehlerantwort die übrigen Gruppen bis zum Sync ohne Warten,
// liest bis zum ReadyForQuery und vergleicht die Folge ab dem Fehler wie
// aufgezeichnet; bei gleichem SQLSTATE ist es keine Abweichung, bei einem
// anderen PGR-E5004 am Ende der Interaktion (LH-FA-24.a, LH-FA-20.a Schritt 6).
func TestVergleichExtendedFehlerMitten(t *testing.T) {
	g1 := gruppe([]model.Response{{Type: model.ResponseParseComplete}, fehlerMit("42P01", "a")}, cm(model.ClientParse), cm(model.ClientExecute), cm(model.ClientFlush))
	g2 := gruppe(bereit(), cm(model.ClientSync))
	rec := sitzung(extendedInteraktion(1, g1, g2), nurSync(2, bereitMit("I")))
	server := func(code string) *fakeEinspiel {
		return &fakeEinspiel{gruppen: [][]schritt{
			schritte(model.Response{Type: model.ResponseParseComplete}, fehlerMit(code, "b"), model.Response{Type: model.ResponseParseComplete}),
			schritte(bereitMit("I")),
		}}
	}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{server("42P01")}}
	_, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{CompareResponses: true})
	if wantText := []string{e5004 + abweichungText(3, "Nachrichtenart")}; !reflect.DeepEqual(texte(err), wantText) {
		t.Fatalf("zusätzliche Nachricht nach dem Fehler: Meldungen %q, erwartet %q", texte(err), wantText)
	}
	rec = sitzung(extendedInteraktion(1, g1, g2))
	ziel = &fakeZiel{sessions: []*fakeEinspiel{{gruppen: [][]schritt{schritte(model.Response{Type: model.ResponseParseComplete}, fehlerMit("42P01", "b")), schritte(bereitMit("I"))}}}}
	if z, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{CompareResponses: true}); err != nil || z.Abweichend != 0 {
		t.Fatalf("gleicher SQLSTATE: %v, Zähler %+v", err, z)
	}
	want := []string{"verbinde", "gruppe parse,execute,flush", "naechste", "naechste", "gruppe sync", "naechste", "schliesse"}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
	ziel = &fakeZiel{sessions: []*fakeEinspiel{{gruppen: [][]schritt{schritte(model.Response{Type: model.ResponseParseComplete}, fehlerMit("23505", "b")), schritte(bereitMit("I"))}}}}
	_, err = spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{CompareResponses: true})
	if wantText := []string{e5004 + abweichungText(2, "Fehler") + " (SQLSTATE aufgezeichnet „42P01“, Server „23505“)"}; !reflect.DeepEqual(texte(err), wantText) {
		t.Fatalf("anderer SQLSTATE: Meldungen %q, erwartet %q", texte(err), wantText)
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf bei Abweichung %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-24/Boundary — nach einer Fehlerantwort mitten in einer
// Extended-Interaktion sendet play auch mit Vergleich die übrigen Gruppen ohne
// Warten dazwischen, auch eine mit Flush, und liest dann bis zum ReadyForQuery
// (LH-FA-20.a *Gruppen*, Schritt 6; LH-FA-24.a).
func TestVergleichExtendedFehlerOhneWarten(t *testing.T) {
	g1 := gruppe([]model.Response{{Type: model.ResponseParseComplete}, fehlerMit("42P01", "a")}, cm(model.ClientParse), cm(model.ClientExecute), cm(model.ClientFlush))
	g2 := gruppe(nil, cm(model.ClientBind), cm(model.ClientExecute), cm(model.ClientFlush))
	g3 := gruppe(bereit(), cm(model.ClientSync))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{gruppen: [][]schritt{schritte(model.Response{Type: model.ResponseParseComplete}, fehlerMit("42P01", "b")), nil, schritte(bereitMit("I"))}}}}
	z, err := spieleZaehler(context.Background(), t, sitzung(extendedInteraktion(1, g1, g2, g3)), ziel, vergleichen())
	if err != nil || z != (services.Zaehler{Eingespielt: 1, Verglichen: 1}) {
		t.Fatalf("%v, Zähler %+v", err, z)
	}
	want := []string{"verbinde", "gruppe parse,execute,flush", "naechste", "naechste", "gruppe bind,execute,flush", "gruppe sync", "naechste", "schliesse"}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-24/Negative — eine Fehlerantwort FATAL in der Antwort auf
// eine Gruppe mit Flush ist mit Vergleich sofort PGR-E5004: play sendet keine
// weitere Gruppe und liest nichts mehr (LH-FA-20.a Schritt 6; LH-FA-24.a *Ende
// der Aufzeichnung*).
func TestVergleichFatalExtendedGruppen(t *testing.T) {
	g1 := gruppe([]model.Response{fehlerMit("57P01", "a")}, cm(model.ClientParse), cm(model.ClientFlush))
	g2 := gruppe(bereit(), cm(model.ClientSync))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{gruppen: [][]schritt{schritte(schwer("FATAL", "FATAL", "57P01"), bereitMit("I")), schritte(bereitMit("I"))}}}}
	z, err := spieleZaehler(context.Background(), t, sitzung(extendedInteraktion(1, g1, g2)), ziel, vergleichen())
	want := []string{e5004 + abweichungText(1, "Fehler") + " (SQLSTATE aufgezeichnet „57P01“, Server „57P01“)"}
	if got := texte(err); !reflect.DeepEqual(got, want) || exitVon(err) != 5 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 5", got, exitVon(err), want)
	}
	if got, wantAblauf := ziel.liste(), []string{"verbinde", "gruppe parse,flush", "naechste", "schliesse"}; !reflect.DeepEqual(got, wantAblauf) {
		t.Fatalf("Ablauf %q, erwartet %q", got, wantAblauf)
	}
	if z != (services.Zaehler{Eingespielt: 1, Verglichen: 1, Abweichend: 1}) {
		t.Fatalf("Zähler %+v", z)
	}
}

// Abdeckung: LH-FA-24/Boundary — eingespielt zählt jede Interaktion, deren erste
// Client-Nachricht play zu senden begann, auch wenn das Senden scheitert, ohne
// und mit Vergleich, bei einer einfachen Anfrage und bei einer
// Extended-Interaktion; eine Interaktion, die nicht mehr beginnt, zählt nicht
// (LH-FA-24.a *Meldung*, LH-FA-20.a Schritt 4).
func TestVergleichZaehlerEingespielt(t *testing.T) {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	for _, o := range []services.PlayOptions{{}, vergleichen()} {
		rec := sitzung(aufgezeichnet(1, "A", ok...), aufgezeichnet(2, "B", ok...))
		ziel := &fakeZiel{sessions: []*fakeEinspiel{{anfrageFehler: model.Errorf(model.CodeConnectionLost, nil, "senden")}}}
		z, err := spieleZaehler(context.Background(), t, rec, ziel, o)
		if codeVon(err) != model.CodeConnectionLost || z != (services.Zaehler{Eingespielt: 1}) {
			t.Errorf("Senden scheitert, %+v: %v, Zähler %+v, erwartet {1 0 0}", o, err, z)
		}
		rec = sitzung(nurSync(1, ok...), nurSync(2, ok...))
		ziel = &fakeZiel{sessions: []*fakeEinspiel{{gruppeFehler: model.Errorf(model.CodeConnectionLost, nil, "senden")}}}
		z, err = spieleZaehler(context.Background(), t, rec, ziel, o)
		if codeVon(err) != model.CodeConnectionLost || z != (services.Zaehler{Eingespielt: 1}) {
			t.Errorf("Senden der Gruppe scheitert, %+v: %v, Zähler %+v, erwartet {1 0 0}", o, err, z)
		}
	}
}

// Abdeckung: LH-FA-24/Boundary — das erste Abbruchsignal beendet die laufende
// Interaktion mit Vergleich wie ohne Signal: Eine Abweichung darin wird
// festgestellt und gezählt, danach beginnt keine Interaktion; mit
// --continue-on-error bleibt der Exit-Code nach einer früheren Abweichung 5
// (LH-FA-24.a *Signal*, LH-FA-20.a Schritt 7, Tabelle).
func TestVergleichErstesSignal(t *testing.T) {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	rec := sitzung(aufgezeichnet(1, "A", ok...), aufgezeichnet(2, "B", ok...), aufgezeichnet(3, "C", ok...))

	ctx, cancel := context.WithCancel(context.Background())
	s := spielerFolge(false, []model.Response{zeilenkopf(spalte("a", 23)), befehl("SELECT 1"), bereitMit("I")})
	s.beiAnfrage = func(string) { cancel() }
	ziel := &fakeZiel{sessions: []*fakeEinspiel{s}}
	z, err := spieleZaehler(ctx, t, rec, ziel, vergleichen())
	if want := []string{e5004 + abweichungText(1, "Nachrichtenart")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 5 {
		t.Fatalf("Abweichung in der Interaktion des Signals: Meldungen %q, Exit-Code %d, erwartet %q mit 5", texte(err), exitVon(err), want)
	}
	if z != (services.Zaehler{Eingespielt: 1, Verglichen: 1, Abweichend: 1}) {
		t.Fatalf("Zähler %+v", z)
	}

	ctx, cancel = context.WithCancel(context.Background())
	s = spielerFolge(false, []model.Response{befehl("UPDATE 1"), bereitMit("I")}, ok)
	s.beiAnfrage = func(sql string) {
		if sql == "B" {
			cancel()
		}
	}
	ziel = &fakeZiel{sessions: []*fakeEinspiel{s}}
	z, err = spieleZaehler(ctx, t, rec, ziel, services.PlayOptions{CompareResponses: true, ContinueOnError: true})
	if want := []string{e5004 + abweichungText(1, "Befehl")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 5 {
		t.Fatalf("Signal nach früherer Abweichung: Meldungen %q, Exit-Code %d, erwartet %q mit 5", texte(err), exitVon(err), want)
	}
	if z != (services.Zaehler{Eingespielt: 2, Verglichen: 2, Abweichend: 1}) {
		t.Fatalf("Zähler %+v", z)
	}
	if got, want := ziel.liste(), []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "naechste", "schliesse"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-24/Boundary — das zweite Signal unterbricht die laufende
// Interaktion ohne Vergleich: Sie zählt nicht als verglichen, auch wenn die bis
// dahin gelesenen Antworten abwichen, und die Unterbrechung ist kein Fehler;
// eine frühere Abweichung bleibt, der Exit-Code ist 5 (LH-FA-24.a *Signal*,
// *Meldung*, LH-FA-20.a Tabelle *Abbruchsignal*).
func TestVergleichZweitesSignal(t *testing.T) {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	rec := sitzung(aufgezeichnet(1, "A", ok...), aufgezeichnet(2, "B", ok...))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{
		schritte(befehl("UPDATE 1"), bereitMit("I")),
		{{r: zeilenkopf(spalte("a", 23))}, {blockiert: true}},
	}}}}
	s, err := services.NewPlayService(context.Background(), ladeRepo{rec: rec}, "r.yaml", ziel, services.PlayOptions{CompareResponses: true, ContinueOnError: true})
	if err != nil {
		t.Fatal(err)
	}
	ablauf := make(chan struct{})
	err = warteAufPlay(context.Background(), t, s, ablauf, func() {
		bis(t, func() bool { return len(ziel.liste()) >= 7 }, "Play wartet nicht binnen 30 s auf die zweite Antwort der zweiten Interaktion")
		close(ablauf)
	})
	if want := []string{e5004 + abweichungText(1, "Befehl")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 5 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 5", texte(err), exitVon(err), want)
	}
	if z := s.Zaehler(); z != (services.Zaehler{Eingespielt: 2, Verglichen: 1, Abweichend: 1}) {
		t.Fatalf("Zähler %+v, erwartet {2 1 1}", z)
	}
}
