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

// alleOptionen sind die Optionen mit Vergleich und --diff-all.
func alleOptionen() services.PlayOptions {
	return services.PlayOptions{CompareResponses: true, DiffAll: true}
}

// alleText ist die Meldung mit --diff-all an Session 1 und Interaktion 1: die
// Gesamtzahl gesamt und die Teile, je Position einer.
func alleText(gesamt int, teile ...string) string {
	return fmt.Sprintf("%sSession 1, Interaktion 1: Antworten weichen von der Aufzeichnung ab, Abweichungen: %d; %s", e5004, gesamt, strings.Join(teile, "; "))
}

// allePos ist der Teil einer Position der Meldung ohne Einzelheiten.
func allePos(position int, art string) string {
	return fmt.Sprintf("Position %d, Art %s", position, art)
}

// alleMit ist der Teil einer Position der Meldung mit Einzelheiten.
func alleMit(position int, art, einzelheiten string) string {
	return allePos(position, art) + " (" + einzelheiten + ")"
}

func paramBeschreibung(oids ...uint32) model.Response {
	return model.Response{Type: model.ResponseParameterDescription, ParamTypes: oids}
}

// afall ist ein Fall der Meldung mit --diff-all: die aufgezeichneten Antworten
// soll und die des Servers ist und die erwartete Meldung (leer: keine Meldung);
// verglichen mit den Formen der Interaktion, einfach und Extended mit Sync.
type afall struct {
	name        string
	soll, ist   []model.Response
	want        string
	nurExtended bool
	nurEinfach  bool
}

// spieleAlle spielt eine Interaktion in der Form (extended oder einfach) mit den
// Optionen o ein und liefert Fehler und Zähler.
func spieleAlle(t *testing.T, extended bool, soll, ist []model.Response, o services.PlayOptions) (services.Zaehler, error) {
	t.Helper()
	in := aufgezeichnet(1, "A", soll...)
	if extended {
		in = nurSync(1, soll...)
	}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(extended, ist)}}
	return spieleZaehler(context.Background(), t, sitzung(in), ziel, o)
}

func afaelle() []afall {
	rfq := bereitMit("I")
	kopfA := zeilenkopf(spalte("a", 23))
	return []afall{
		{
			name: "mehrere Arten, beide Seiten",
			soll: []model.Response{kopfA, befehl("SELECT 1"), rfq},
			ist:  []model.Response{zeilenkopf(spalte("b", 25)), befehl("UPDATE 1"), bereitMit("T")},
			want: alleText(3,
				alleMit(1, "Spaltenbeschreibung", "Spalte 1: Name aufgezeichnet „a“, Server „b“, Typ-OID aufgezeichnet 23, Server 25"),
				alleMit(2, "Befehl", "Befehl aufgezeichnet „SELECT“, Server „UPDATE“"),
				alleMit(3, "Transaktionsstatus", "aufgezeichnet „I“, Server „T“")),
		},
		{
			name: "eine Abweichung in Listenform",
			soll: []model.Response{befehl("SELECT 1"), rfq},
			ist:  []model.Response{befehl("UPDATE 1"), rfq},
			want: alleText(1, alleMit(1, "Befehl", "Befehl aufgezeichnet „SELECT“, Server „UPDATE“")),
		},
		{name: "keine Abweichung", soll: []model.Response{befehl("SELECT 1"), rfq}, ist: []model.Response{befehl("SELECT 5"), rfq}},
		{
			name: "fehlende Nachricht verschiebt den Rest",
			soll: []model.Response{kopfA, befehl("SELECT 1"), rfq},
			ist:  []model.Response{befehl("SELECT 1"), rfq},
			want: alleText(3, allePos(1, "Nachrichtenart"), allePos(2, "Nachrichtenart"), allePos(3, "Nachrichtenart")),
		},
		{
			name: "Server länger, jede Position einzeln",
			soll: []model.Response{befehl("SELECT 1"), rfq},
			ist:  []model.Response{befehl("SELECT 1"), befehl("SELECT 1"), befehl("SELECT 1"), rfq},
			want: alleText(3, allePos(2, "Nachrichtenart"), allePos(3, "Nachrichtenart"), allePos(4, "Nachrichtenart")),
		},
		{
			name: "Aufzeichnung länger, jede Position einzeln", nurEinfach: true,
			soll: []model.Response{befehl("SELECT 1"), rfq, befehl("SELECT 1"), rfq},
			ist:  []model.Response{befehl("SELECT 1"), rfq},
			want: alleText(2, allePos(3, "Nachrichtenart"), allePos(4, "Nachrichtenart")),
		},
		{
			name: "Fehler mit dem SQLSTATE beider Seiten, ohne Text",
			soll: []model.Response{fehlerMit("42P01", "aufgezeichneter Text"), rfq},
			ist:  []model.Response{fehlerMit("23505", "geheimer Text"), bereitMit("T")},
			want: alleText(2,
				alleMit(1, "Fehler", "SQLSTATE aufgezeichnet „42P01“, Server „23505“"),
				alleMit(2, "Transaktionsstatus", "aufgezeichnet „I“, Server „T“")),
		},
		{
			name: "Fehlerantwort nur beim Server, Nachrichtenart ohne Einzelheiten",
			soll: []model.Response{befehl("SELECT 1"), rfq},
			ist:  []model.Response{fehlerMit("42P01", "geheim"), rfq},
			want: alleText(1, allePos(1, "Nachrichtenart")),
		},
		{
			name: "Fehler ohne Feld C nur beim Server, SQLSTATE leer",
			soll: []model.Response{fehlerMit("42P01", "x"), rfq},
			ist:  []model.Response{{Type: model.ResponseErrorResponse, Fields: map[string]string{"M": "b"}}, rfq},
			want: alleText(1, alleMit(1, "Fehler", "SQLSTATE aufgezeichnet „42P01“, Server „“")),
		},
		{
			name: "Befehl ohne Zahlen, Zeilen nie genannt",
			soll: []model.Response{kopfA, zeile(), befehl("SELECT 1"), rfq},
			ist:  []model.Response{kopfA, zeile(), zeile(), befehl("UPDATE 5"), rfq},
			want: alleText(1, alleMit(2, "Befehl", "Befehl aufgezeichnet „SELECT“, Server „UPDATE“")),
		},
		{
			name: "Spaltenzahl verschieden, Spalte ohne Gegenstück nicht einzeln",
			soll: []model.Response{zeilenkopf(spalte("a", 23), spalte("b", 23)), rfq},
			ist:  []model.Response{zeilenkopf(spalte("a", 23)), rfq},
			want: alleText(1, alleMit(1, "Spaltenbeschreibung", "Spalten aufgezeichnet 2, Server 1")),
		},
		{
			name: "Spaltenzahl und gemeinsame Spalte, Typ-OID",
			soll: []model.Response{zeilenkopf(spalte("a", 23), spalte("b", 23), spalte("c", 23)), rfq},
			ist:  []model.Response{zeilenkopf(spalte("a", 25), spalte("b", 23)), rfq},
			want: alleText(1, alleMit(1, "Spaltenbeschreibung", "Spalten aufgezeichnet 3, Server 2; Spalte 1: Typ-OID aufgezeichnet 23, Server 25")),
		},
		{
			name: "Spaltenname bei gleicher Zahl, aufsteigende Nummern",
			soll: []model.Response{zeilenkopf(spalte("a", 23), spalte("b", 23), spalte("c", 23)), rfq},
			ist:  []model.Response{zeilenkopf(spalte("x", 23), spalte("b", 23), spalte("y", 23)), rfq},
			want: alleText(1, alleMit(1, "Spaltenbeschreibung", "Spalte 1: Name aufgezeichnet „a“, Server „x“; Spalte 3: Name aufgezeichnet „c“, Server „y“")),
		},
		{
			name: "Parameterzahl und Typ-OID", nurExtended: true,
			soll: []model.Response{{Type: model.ResponseParseComplete}, paramBeschreibung(23, 25), rfq},
			ist:  []model.Response{{Type: model.ResponseParseComplete}, paramBeschreibung(23, 26, 1), rfq},
			want: alleText(1, alleMit(2, "Spaltenbeschreibung", "Parameter aufgezeichnet 2, Server 3; Parameter 2: Typ-OID aufgezeichnet 25, Server 26")),
		},
		{
			name: "Parameter gleicher Zahl, jeder abweichende", nurExtended: true,
			soll: []model.Response{{Type: model.ResponseParseComplete}, paramBeschreibung(23, 25), rfq},
			ist:  []model.Response{{Type: model.ResponseParseComplete}, paramBeschreibung(24, 26), rfq},
			want: alleText(1, alleMit(2, "Spaltenbeschreibung", "Parameter 1: Typ-OID aufgezeichnet 23, Server 24; Parameter 2: Typ-OID aufgezeichnet 25, Server 26")),
		},
	}
}

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Boundary, LH-FA-24/Negative — mit
// --diff-all meldet play je Interaktion mit mindestens einer Abweichung genau eine
// Meldung PGR-E5004 (Exit-Code 5) mit Session, Sequenznummer und der Gesamtzahl der
// Abweichungen, dann die Abweichungen in aufsteigender Position, je Position Art
// und Einzelheiten beider Seiten, auch bei einer einzigen; ohne Abweichung keine
// Meldung. Zuordnung Position für Position ohne Abgleich (eine fehlende Nachricht
// verschiebt den Rest, die Folgepositionen werden gemeldet), verschieden lange
// Folgen jede Position einzeln; Einzelheiten je Art: Fehler SQLSTATE beider
// Seiten, Befehl ohne Zahlen, Transaktionsstatus, Spalten und Parameter (Zahl,
// je abweichende gemeinsame Nummer), Nachrichtenart keine; nie Zeilenwerte,
// Zeilenzahl oder Text einer Fehlerantwort; bei einer einfachen Anfrage und einer
// Extended-Interaktion (LH-FA-24.a *Alle Abweichungen*).
func TestVergleichAlle(t *testing.T) {
	for _, f := range afaelle() {
		for _, extended := range []bool{false, true} {
			if (f.nurExtended && !extended) || (f.nurEinfach && extended) {
				continue
			}
			t.Run(fmt.Sprintf("%s/extended=%v", f.name, extended), func(t *testing.T) {
				z, err := spieleAlle(t, extended, f.soll, f.ist, alleOptionen())
				if f.want == "" {
					if err != nil || z != (services.Zaehler{Eingespielt: 1, Verglichen: 1}) {
						t.Fatalf("erwartet ohne Meldung: %v, Zähler %+v", err, z)
					}
					return
				}
				want := f.want
				if got := texte(err); !reflect.DeepEqual(got, []string{want}) || exitVon(err) != 5 || codeVon(err) != model.CodeResponseMismatch {
					t.Fatalf("Meldungen %q, Exit-Code %d, erwartet [%q] mit 5", got, exitVon(err), want)
				}
				if z != (services.Zaehler{Eingespielt: 1, Verglichen: 1, Abweichend: 1}) {
					t.Fatalf("Zähler %+v, erwartet {1 1 1}", z)
				}
				for _, geheim := range []string{"geheim", "aufgezeichneter Text", "SELECT 1", "UPDATE 5", "„1“"} {
					if strings.Contains(texte(err)[0], geheim) {
						t.Errorf("Meldung nennt %q: %q", geheim, texte(err)[0])
					}
				}
			})
		}
	}
}

// Abdeckung: LH-FA-24/Boundary — ohne --diff-all bleibt es bei der ersten
// Abweichung, Wort für Wort die Meldung ohne Liste, auch wenn die Interaktion
// mehrere trägt; --diff-all=false ohne Wirkung ist dasselbe (LH-FA-24.a *Alle
// Abweichungen*, Einleitung).
func TestVergleichAlleStandardErste(t *testing.T) {
	soll := []model.Response{zeilenkopf(spalte("a", 23)), befehl("SELECT 1"), bereitMit("I")}
	ist := []model.Response{zeilenkopf(spalte("b", 25)), befehl("UPDATE 1"), bereitMit("T")}
	for _, extended := range []bool{false, true} {
		_, err := spieleAlle(t, extended, soll, ist, vergleichen())
		if want := []string{e5004 + abweichungText(1, "Spaltenbeschreibung")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 5 {
			t.Errorf("extended=%v: Meldungen %q, Exit-Code %d, erwartet %q mit 5", extended, texte(err), exitVon(err), want)
		}
	}
}

// grosse ist eine Interaktion mit n Abweichungen der Art Befehl, in den Folgen
// aus n Befehlen: gleich die Aufzeichnung "A", der Server "B".
func grosse(n int) (soll, ist []model.Response) {
	for range n {
		soll, ist = append(soll, befehl("A")), append(ist, befehl("B"))
	}
	return append(soll, bereitMit("I")), append(ist, bereitMit("I"))
}

// Abdeckung: LH-FA-24/Boundary — die Meldung mit --diff-all nennt höchstens 20
// Positionen je Interaktion: Genau 20 Abweichungen stehen ohne Zusatz, 21 mit
// „und 1 weitere“, bei 25 steht die Gesamtzahl 25 im Kopf und „und 5 weitere“
// hinter der zwanzigsten; null Abweichungen sind keine Meldung (LH-FA-24.a
// *Obergrenze*).
func TestVergleichAlleObergrenzePositionen(t *testing.T) {
	for _, f := range []struct {
		n       int
		zusatz  string
		erwarte int
	}{{20, "", 20}, {21, "und 1 weitere", 20}, {25, "und 5 weitere", 20}, {1, "", 1}} {
		soll, ist := grosse(f.n)
		_, err := spieleAlle(t, false, soll, ist, alleOptionen())
		got := texte(err)
		if len(got) != 1 || !strings.Contains(got[0], fmt.Sprintf("Abweichungen: %d; ", f.n)) {
			t.Fatalf("n=%d: Meldungen %q", f.n, got)
		}
		if n := strings.Count(got[0], "Position "); n != f.erwarte {
			t.Errorf("n=%d: %d Positionen in der Meldung, erwartet %d: %s", f.n, n, f.erwarte, got[0])
		}
		if f.zusatz == "" && strings.Contains(got[0], "weitere") {
			t.Errorf("n=%d: Zusatz „weitere“ in %s", f.n, got[0])
		}
		if f.zusatz != "" && !strings.HasSuffix(got[0], "; "+f.zusatz) {
			t.Errorf("n=%d: Meldung endet nicht mit %q: %s", f.n, f.zusatz, got[0])
		}
		if !strings.Contains(got[0], fmt.Sprintf("Position %d, Art Befehl", min(f.n, 20))) {
			t.Errorf("n=%d: letzte genannte Position fehlt: %s", f.n, got[0])
		}
	}
}

// Abdeckung: LH-FA-24/Boundary — je Position nennt die Meldung mit --diff-all
// höchstens 20 abweichende Spalten oder Parameter und dahinter „und N weitere“:
// 20 stehen ohne Zusatz, 21 mit „und 1 weitere“ (LH-FA-24.a *Einzelheiten je Art,
// beide Seiten*).
func TestVergleichAlleObergrenzeEinzelne(t *testing.T) {
	spalten := func(n int, oid uint32) []model.Column {
		var cs []model.Column
		for i := range n {
			cs = append(cs, spalte(fmt.Sprintf("c%d", i+1), oid))
		}
		return cs
	}
	params := func(n int, oid uint32) []uint32 {
		var os []uint32
		for range n {
			os = append(os, oid)
		}
		return os
	}
	for _, n := range []int{20, 21} {
		for _, f := range []struct {
			name, stichwort string
			extended        bool
			soll, ist       []model.Response
		}{
			{"Spalten", "Spalte ", false, []model.Response{zeilenkopf(spalten(n, 1)...), bereitMit("I")}, []model.Response{zeilenkopf(spalten(n, 2)...), bereitMit("I")}},
			{"Parameter", "Parameter ", true, []model.Response{paramBeschreibung(params(n, 1)...), bereitMit("I")}, []model.Response{paramBeschreibung(params(n, 2)...), bereitMit("I")}},
		} {
			_, err := spieleAlle(t, f.extended, f.soll, f.ist, alleOptionen())
			got := texte(err)
			if len(got) != 1 {
				t.Fatalf("%s n=%d: Meldungen %q", f.name, n, got)
			}
			if c := strings.Count(got[0], f.stichwort+fmt.Sprint(n, ":")) + strings.Count(got[0], f.stichwort+"20:"); (n == 20 && c < 1) || (n == 21 && c != 1) {
				t.Errorf("%s n=%d: Nummer %d oder 20 nicht wie erwartet in %s", f.name, n, n, got[0])
			}
			if n == 20 && (strings.Contains(got[0], "weitere") || !strings.Contains(got[0], f.stichwort+"20:")) {
				t.Errorf("%s n=20: %s", f.name, got[0])
			}
			if n == 21 && (!strings.HasSuffix(got[0], "und 1 weitere)") || strings.Contains(got[0], f.stichwort+"21:") || !strings.Contains(got[0], f.stichwort+"20:")) {
				t.Errorf("%s n=21: %s", f.name, got[0])
			}
		}
	}
	// 25 Abweichungen: „und 5 weitere“; die Gesamtzahl der Positionen bleibt 1.
	_, err := spieleAlle(t, false, []model.Response{zeilenkopf(spalten(25, 1)...), bereitMit("I")}, []model.Response{zeilenkopf(spalten(25, 2)...), bereitMit("I")}, alleOptionen())
	if got := texte(err); len(got) != 1 || !strings.HasSuffix(got[0], "und 5 weitere)") || !strings.Contains(got[0], "Abweichungen: 1; ") {
		t.Errorf("25 Spalten: Meldungen %q", got)
	}
}

// Abdeckung: LH-FA-24/Boundary — ohne Wirkung von --diff-all auf den Vergleich
// selbst und die Meldung des ersten Signals: Das erste Signal beendet die laufende
// Interaktion mit Vergleich, ihre Liste ist vollständig; das zweite unterbricht sie
// ohne Vergleich und ohne Meldung der Interaktion (LH-FA-24.a *Alle Abweichungen*,
// *Signal*).
func TestVergleichAlleSignal(t *testing.T) {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	rec := sitzung(aufgezeichnet(1, "A", ok...), aufgezeichnet(2, "B", ok...))
	ctx, cancel := context.WithCancel(context.Background())
	s := spielerFolge(false, []model.Response{zeilenkopf(spalte("a", 23)), befehl("UPDATE 1"), bereitMit("T")})
	s.beiAnfrage = func(string) { cancel() }
	ziel := &fakeZiel{sessions: []*fakeEinspiel{s}}
	z, err := spieleZaehler(ctx, t, rec, ziel, alleOptionen())
	want := alleText(3, allePos(1, "Nachrichtenart"), allePos(2, "Nachrichtenart"), allePos(3, "Nachrichtenart"))
	if got := texte(err); !reflect.DeepEqual(got, []string{want}) || exitVon(err) != 5 {
		t.Fatalf("erstes Signal: Meldungen %q, Exit-Code %d, erwartet [%q] mit 5", got, exitVon(err), want)
	}
	if z != (services.Zaehler{Eingespielt: 1, Verglichen: 1, Abweichend: 1}) {
		t.Fatalf("erstes Signal: Zähler %+v", z)
	}

	ziel = &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{{{r: zeilenkopf(spalte("a", 23))}, {blockiert: true}}}}}}
	sv, errNeu := services.NewPlayService(context.Background(), ladeRepo{rec: rec}, "r.yaml", ziel, alleOptionen())
	if errNeu != nil {
		t.Fatal(errNeu)
	}
	ablauf := make(chan struct{})
	err = warteAufPlay(context.Background(), t, sv, ablauf, func() {
		bis(t, func() bool { return len(ziel.liste()) >= 4 }, "Play wartet nicht binnen 30 s auf die zweite Antwort")
		close(ablauf)
	})
	if err != nil || sv.Zaehler() != (services.Zaehler{Eingespielt: 1}) {
		t.Fatalf("zweites Signal: %v, Zähler %+v, erwartet ohne Meldung mit {1 0 0}", err, sv.Zaehler())
	}
}

// mehrere ist eine Aufzeichnung mit drei Interaktionen in Session 1 und einer in
// Session 2, die der Server je abweichend (zwei Abweichungen: Befehl und
// Transaktionsstatus) beantwortet, außer Interaktion 2.
func mehrere() (model.Recording, *fakeZiel) {
	gut := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	schlecht := []model.Response{befehl("UPDATE 1"), bereitMit("T")}
	rec := model.NewRecording()
	rec.Sessions = []model.Session{
		{ID: 1, Interactions: []model.Interaction{aufgezeichnet(1, "A", gut...), aufgezeichnet(2, "B", gut...), aufgezeichnet(3, "C", gut...)}},
		{ID: 2, Interactions: []model.Interaction{aufgezeichnet(1, "D", gut...)}},
	}
	return rec, &fakeZiel{sessions: []*fakeEinspiel{spielerFolge(false, schlecht, gut, schlecht), spielerFolge(false, schlecht)}}
}

// zweiAbw sind die Teile der Meldung zu mehrere für eine Interaktion.
func zweiAbw() []string {
	return []string{alleMit(1, "Befehl", "Befehl aufgezeichnet „SELECT“, Server „UPDATE“"), alleMit(2, "Transaktionsstatus", "aufgezeichnet „I“, Server „T“")}
}

// Abdeckung: LH-FA-24/Boundary, LH-FA-24/Negative — mit --diff-all steht mit
// --continue-on-error in jeder Interaktion mit Abweichung die vollständige Liste
// und das Einspielen läuft weiter (Meldungen in der Reihenfolge des Auftretens,
// Exit-Code 5); ohne --continue-on-error endet play nach der Interaktion mit der
// Abweichung, deren Liste vollständig ist; eingespielt, verglichen und abweichend
// zählen Interaktionen, abweichend eine Interaktion mit mehreren Abweichungen
// einmal; --allow-recorded-errors ändert nichts (LH-FA-24.a *Zusammenspiel*,
// *Zähler*).
func TestVergleichAlleZusammenspiel(t *testing.T) {
	for _, allow := range []bool{false, true} {
		rec, ziel := mehrere()
		o := alleOptionen()
		o.ContinueOnError, o.AllowRecordedErrors = true, allow
		z, err := spieleZaehler(context.Background(), t, rec, ziel, o)
		want := []string{
			alleText(2, zweiAbw()...),
			strings.Replace(alleText(2, zweiAbw()...), "Interaktion 1", "Interaktion 3", 1),
			strings.Replace(alleText(2, zweiAbw()...), "Session 1, Interaktion 1", "Session 2, Interaktion 1", 1),
		}
		if got := texte(err); !reflect.DeepEqual(got, want) || exitVon(err) != 5 {
			t.Errorf("allow=%v mit continue: Meldungen %q, Exit-Code %d, erwartet %q mit 5", allow, got, exitVon(err), want)
		}
		if z != (services.Zaehler{Eingespielt: 4, Verglichen: 4, Abweichend: 3}) {
			t.Errorf("allow=%v mit continue: Zähler %+v, erwartet {4 4 3}", allow, z)
		}

		rec, ziel = mehrere()
		o.ContinueOnError = false
		z, err = spieleZaehler(context.Background(), t, rec, ziel, o)
		if got, want := texte(err), []string{alleText(2, zweiAbw()...)}; !reflect.DeepEqual(got, want) || exitVon(err) != 5 {
			t.Errorf("allow=%v ohne continue: Meldungen %q, Exit-Code %d, erwartet %q mit 5", allow, got, exitVon(err), want)
		}
		if z != (services.Zaehler{Eingespielt: 1, Verglichen: 1, Abweichend: 1}) {
			t.Errorf("allow=%v ohne continue: Zähler %+v, erwartet {1 1 1}", allow, z)
		}
		if got, want := ziel.liste(), []string{"verbinde", "anfrage A", "naechste", "naechste", "schliesse"}; !reflect.DeepEqual(got, want) {
			t.Errorf("allow=%v ohne continue: Ablauf %q, erwartet %q", allow, got, want)
		}
	}
}

// Abdeckung: LH-FA-24/Boundary — ohne --compare-responses hat DiffAll keine
// Wirkung auf den Service: kein Vergleich, keine Meldung, verglichen und abweichend
// sind 0 (LH-FA-24.a *Alle Abweichungen*, *Ohne --compare-responses*: der
// CLI-Adapter lehnt die Kombination ab).
func TestVergleichAlleOhneVergleich(t *testing.T) {
	rec, ziel := mehrere()
	z, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{DiffAll: true})
	if err != nil || z != (services.Zaehler{Eingespielt: 4}) {
		t.Fatalf("%v, Zähler %+v, erwartet nil mit {4 0 0}", err, z)
	}
}

// Abdeckung: LH-FA-24/Happy, LH-FA-24/Negative, LH-FA-25/Boundary — mit --diff-all
// ist eine Fehlerantwort FATAL oder PANIC des Servers die letzte Position: Die
// Positionen davor werden vollständig gemeldet, dahinter erscheint nichts, auch
// keine Abweichung der aufgezeichneten Folge; ein erwartetes Ende (gleicher SQLSTATE
// an der aufgezeichneten Stelle) ist keine Position der Liste, Abweichungen davor
// bleiben eine Meldung; sonst ist die Position des Endes eine Abweichung der Art
// Fehler mit dem SQLSTATE beider Seiten (der aufgezeichnete leer, wo dort keine
// Fehlerantwort steht), auch bei gleichem SQLSTATE einer aufgezeichneten
// Fehlerantwort ohne Ende (LH-FA-24.a *Ende des Servers*, *Ende der Aufzeichnung*).
func TestVergleichAlleEndeDesServers(t *testing.T) {
	fatal := schwer("FATAL", "FATAL", "57P01")
	kopfA, kopfB := zeilenkopf(spalte("a", 23)), zeilenkopf(spalte("b", 23))
	geschlossen := model.Response{Type: model.ResponseConnectionClosed}
	for _, f := range []struct {
		name      string
		soll, ist []model.Response
		want      string
	}{
		{
			name: "erwartetes Ende ist keine Position",
			soll: []model.Response{fatal, geschlossen}, ist: []model.Response{fatal},
			want: "",
		},
		{
			name: "Abweichung vor dem erwarteten Ende",
			soll: []model.Response{kopfA, fatal, geschlossen}, ist: []model.Response{kopfB, fatal},
			want: alleText(1, alleMit(1, "Spaltenbeschreibung", "Spalte 1: Name aufgezeichnet „a“, Server „b“")),
		},
		{
			name: "Ende ohne Gegenstück, nichts dahinter",
			soll: []model.Response{kopfA, befehl("SELECT 1"), bereitMit("I")}, ist: []model.Response{kopfB, fatal},
			want: alleText(2,
				alleMit(1, "Spaltenbeschreibung", "Spalte 1: Name aufgezeichnet „a“, Server „b“"),
				alleMit(2, "Fehler", "SQLSTATE aufgezeichnet „“, Server „57P01“")),
		},
		{
			name: "aufgezeichnete Fehlerantwort ohne Ende, gleicher SQLSTATE",
			soll: []model.Response{kopfA, fehlerMit("57P01", "x"), bereitMit("I")}, ist: []model.Response{kopfA, fatal},
			want: alleText(1, alleMit(2, "Fehler", "SQLSTATE aufgezeichnet „57P01“, Server „57P01“")),
		},
		{
			name: "anderer SQLSTATE am aufgezeichneten Ende",
			soll: []model.Response{kopfA, fatal, geschlossen}, ist: []model.Response{kopfA, schwer("FATAL", "FATAL", "40P01")},
			want: alleText(1, alleMit(2, "Fehler", "SQLSTATE aufgezeichnet „57P01“, Server „40P01“")),
		},
	} {
		t.Run(f.name, func(t *testing.T) {
			z, err := spieleAlle(t, false, f.soll, f.ist, alleOptionen())
			if f.want == "" {
				if err != nil || z != (services.Zaehler{Eingespielt: 1, Verglichen: 1}) {
					t.Fatalf("erwartetes Ende: %v, Zähler %+v, erwartet nil mit {1 1 0}", err, z)
				}
				return
			}
			if got := texte(err); !reflect.DeepEqual(got, []string{f.want}) || exitVon(err) != 5 {
				t.Fatalf("Meldungen %q, Exit-Code %d, erwartet [%q] mit 5", got, exitVon(err), f.want)
			}
			if z != (services.Zaehler{Eingespielt: 1, Verglichen: 1, Abweichend: 1}) {
				t.Fatalf("Zähler %+v", z)
			}
		})
	}
}

// Abdeckung: LH-FA-24/Boundary, LH-FA-25/Boundary — mit --diff-all und der letzten
// Gruppe mit Flush eines aufgezeichneten Endes zählen die Positionen der
// aufgezeichneten Folge, die der Server nicht beantwortet hat, bis
// connection_closed, je eine Abweichung der Art Nachrichtenart, einzeln
// (LH-FA-24.a *Zuordnung der Positionen*).
func TestVergleichAlleFlushEnde(t *testing.T) {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	f := formFlush()
	rec := model.NewRecording()
	rec.Sessions = []model.Session{{ID: 1, Interactions: []model.Interaction{aufgezeichnet(1, "A", ok...), f.interaktion(endeAufgezeichnet())}}}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{f.session([]model.Response{befehl("SELECT 1")}, false)}}
	_, err := spieleZaehler(context.Background(), t, rec, ziel, alleOptionen())
	want := "Replay [PGR-E5004]: Session 1, Interaktion 2: Antworten weichen von der Aufzeichnung ab, Abweichungen: 2; " + allePos(3, "Nachrichtenart") + "; " + allePos(4, "Nachrichtenart")
	if got := texte(err); !reflect.DeepEqual(got, []string{want}) || exitVon(err) != 5 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet [%q] mit 5", got, exitVon(err), want)
	}
}

// Abdeckung: LH-FA-24/Boundary — mit --diff-all sind verschiedene Typen an einer
// Position (ready_for_query gegen connection_closed) eine Position der Art
// Nachrichtenart; die Stelle jeder Nachrichtenart ist die der Regeln ohne die
// Option (LH-FA-24.a *Stelle in der Folge*).
func TestVergleichAlleStelle(t *testing.T) {
	fatal := schwer("FATAL", "FATAL", "57P01")
	soll := []model.Response{fatal, {Type: model.ResponseConnectionClosed}}
	for _, f := range []struct {
		name string
		ist  []model.Response
		want string
	}{
		{"ready_for_query gegen connection_closed", []model.Response{fehlerMit("57P01", "x"), bereitMit("I")}, alleText(1, allePos(2, "Nachrichtenart"))},
		{"command_complete gegen connection_closed, danach ohne Gegenstück", []model.Response{fehlerMit("57P01", "x"), befehl("SELECT 1"), bereitMit("I")}, alleText(2, allePos(2, "Nachrichtenart"), allePos(3, "Nachrichtenart"))},
	} {
		_, err := spieleAlle(t, false, soll, f.ist, alleOptionen())
		if got := texte(err); !reflect.DeepEqual(got, []string{f.want}) || exitVon(err) != 5 {
			t.Errorf("%s: Meldungen %q, Exit-Code %d, erwartet [%q] mit 5", f.name, got, exitVon(err), f.want)
		}
	}
}

// Abdeckung: LH-FA-24/Boundary — weichen Name und Typ-OID einer Spalte ab, steht
// ein Eintrag mit beiden Feldern (Name zuerst), und die Grenze von 20 zählt
// Spalten, nicht Felder: 20 solche Spalten stehen ohne Zusatz, 21 mit „und 1
// weitere“ (LH-FA-24.a *Einzelheiten je Art, beide Seiten*).
func TestVergleichAlleObergrenzeSpaltenMitBeidenFeldern(t *testing.T) {
	spalten := func(n int, praefix string, oid uint32) []model.Column {
		var cs []model.Column
		for i := range n {
			cs = append(cs, spalte(fmt.Sprintf("%s%d", praefix, i+1), oid))
		}
		return cs
	}
	for _, n := range []int{20, 21} {
		soll := []model.Response{zeilenkopf(spalten(n, "c", 1)...), bereitMit("I")}
		ist := []model.Response{zeilenkopf(spalten(n, "x", 2)...), bereitMit("I")}
		_, err := spieleAlle(t, false, soll, ist, alleOptionen())
		got := texte(err)
		if len(got) != 1 {
			t.Fatalf("n=%d: Meldungen %q", n, got)
		}
		eintrag := "Spalte 20: Name aufgezeichnet „c20“, Server „x20“, Typ-OID aufgezeichnet 1, Server 2"
		if !strings.Contains(got[0], eintrag) || strings.Count(got[0], "Typ-OID aufgezeichnet") != 20 {
			t.Errorf("n=%d: 20 Einträge mit beiden Feldern erwartet: %s", n, got[0])
		}
		if n == 20 && strings.Contains(got[0], "weitere") {
			t.Errorf("n=20: Zusatz in %s", got[0])
		}
		if n == 21 && (!strings.HasSuffix(got[0], "und 1 weitere)") || strings.Contains(got[0], "Spalte 21:")) {
			t.Errorf("n=21: %s", got[0])
		}
	}
}
