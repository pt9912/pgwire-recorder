package services_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// endeAufgezeichnet sind die aufgezeichneten Antworten einer einfachen Anfrage,
// die mit dem Ende des Servers endet: FATAL 57P01, dann connection_closed.
func endeAufgezeichnet() []model.Response {
	return []model.Response{schwer("FATAL", "FATAL", "57P01"), {Type: model.ResponseConnectionClosed}}
}

// endeRec ist eine Aufzeichnung aus zwei Sessions: Session 1 mit den
// Interaktionen A (ohne Ende) und B, Session 2 mit C. B trägt das Ende, wenn
// mitEnde gesetzt ist, sonst eine normale Antwort. Jeder Aufruf baut seine
// Aufzeichnung selbst.
func endeRec(mitEnde bool) model.Recording {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	b := aufgezeichnet(2, "B", ok...)
	if mitEnde {
		b = aufgezeichnet(2, "B", endeAufgezeichnet()...)
	}
	rec := model.NewRecording()
	rec.Sessions = []model.Session{
		{ID: 1, Interactions: []model.Interaction{aufgezeichnet(1, "A", ok...), b}},
		{ID: 2, Interactions: []model.Interaction{aufgezeichnet(1, "C", ok...)}},
	}
	return rec
}

const e4003B = "Netzwerk [PGR-E4003]: Session 1, Interaktion 2: Fehlerantwort des Servers 57P01 „geheim“, die Verbindung endet"

// Abdeckung: LH-FA-20/Boundary, LH-FA-25/Happy — ohne Vergleich ist eine
// Fehlerantwort FATAL oder PANIC (Feld V, ohne es S) auf eine einfache Anfrage
// mit --allow-recorded-errors und einem aufgezeichneten Ende in der Interaktion
// das erwartete Ende: keine Meldung, Exit-Code 0, play liest und sendet nichts
// mehr, schließt die Verbindung ohne Terminate, die Session endet und die nächste
// läuft, auch ohne und mit --continue-on-error; eingespielt zählt die
// Interaktion, verglichen und abweichend bleiben 0 (LH-FA-20.a *Interaktion*,
// *Ende einer Session*, LH-FA-25.a).
func TestPlayEndeErwartet(t *testing.T) {
	for _, f := range fatalFaelle() {
		for _, fortsetzen := range []bool{false, true} {
			name := f.name
			if fortsetzen {
				name += " mit --continue-on-error"
			}
			t.Run(name, func(t *testing.T) {
				sitz := spielerFolge(false, []model.Response{befehl("SELECT 1"), bereitMit("I")}, []model.Response{f.r, bereitMit("I")})
				folge := spielerFolge(false, []model.Response{befehl("SELECT 1"), bereitMit("I")})
				ziel := &fakeZiel{sessions: []*fakeEinspiel{sitz, folge}}
				o := services.PlayOptions{AllowRecordedErrors: true, ContinueOnError: fortsetzen}
				z, err := spieleZaehler(context.Background(), t, endeRec(true), ziel, o)
				if err != nil {
					t.Fatalf("Meldungen %q, erwartet keine", texte(err))
				}
				want := []string{
					"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "schliesse",
					"verbinde", "anfrage C", "naechste", "naechste", "schliesse",
				}
				if got := ziel.liste(); !reflect.DeepEqual(got, want) {
					t.Fatalf("Ablauf %q, erwartet %q", got, want)
				}
				if zu, ohne := sitz.war(); !zu || !ohne {
					t.Fatalf("Session 1: geschlossen %v, ohne Terminate %v, erwartet beides", zu, ohne)
				}
				if zu, ohne := folge.war(); !zu || ohne {
					t.Fatalf("Session 2: geschlossen %v, ohne Terminate %v, erwartet mit Terminate", zu, ohne)
				}
				if z != (services.Zaehler{Eingespielt: 3}) {
					t.Fatalf("Zähler %+v, erwartet eingespielt 3, verglichen 0, abweichend 0", z)
				}
			})
		}
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-25/Negative — ohne Vergleich ist eine
// Fehlerantwort FATAL oder PANIC auf eine einfache Anfrage PGR-E4003, sofort und
// auch mit --continue-on-error, die Verbindung ohne Terminate geschlossen, wenn
// eine Bedingung des erwarteten Endes fehlt: ohne --allow-recorded-errors, ohne
// aufgezeichnetes Ende an der Stelle (auch bei aufgezeichneter Fehlerantwort
// ERROR) oder bei Ende des Servers in einer früheren Interaktion als der
// aufgezeichneten; eine nächste Session läuft nicht (LH-FA-20.a *Interaktion*,
// Tabelle der Fehlerregeln).
func TestPlayEndeNichtErwartet(t *testing.T) {
	for _, f := range []struct {
		name    string
		o       services.PlayOptions
		mitEnde bool
		aufErr  bool
		folge   [][]model.Response
		want    string
		ablauf  []string
	}{
		{
			name: "ohne Option", o: services.PlayOptions{}, mitEnde: true,
			folge:  [][]model.Response{{befehl("SELECT 1"), bereitMit("I")}, {schwer("FATAL", "FATAL", "57P01")}},
			want:   e4003B,
			ablauf: []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "schliesse"},
		},
		{
			name: "ohne Option mit --continue-on-error", o: services.PlayOptions{ContinueOnError: true}, mitEnde: true,
			folge:  [][]model.Response{{befehl("SELECT 1"), bereitMit("I")}, {schwer("FATAL", "FATAL", "57P01")}},
			want:   e4003B,
			ablauf: []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "schliesse"},
		},
		{
			name: "ohne Ende an der Stelle", o: services.PlayOptions{AllowRecordedErrors: true},
			folge:  [][]model.Response{{befehl("SELECT 1"), bereitMit("I")}, {schwer("FATAL", "FATAL", "57P01")}},
			want:   e4003B,
			ablauf: []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "schliesse"},
		},
		{
			name: "ohne Ende an der Stelle, aufgezeichnet ERROR", o: services.PlayOptions{AllowRecordedErrors: true, ContinueOnError: true}, aufErr: true,
			folge:  [][]model.Response{{befehl("SELECT 1"), bereitMit("I")}, {schwer("FATAL", "FATAL", "57P01")}},
			want:   e4003B,
			ablauf: []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "schliesse"},
		},
		{
			name: "FATAL in früherer Interaktion", o: services.PlayOptions{AllowRecordedErrors: true}, mitEnde: true,
			folge:  [][]model.Response{{schwer("FATAL", "FATAL", "57P01")}},
			want:   "Netzwerk [PGR-E4003]: Session 1, Interaktion 1: Fehlerantwort des Servers 57P01 „geheim“, die Verbindung endet",
			ablauf: []string{"verbinde", "anfrage A", "naechste", "schliesse"},
		},
	} {
		t.Run(f.name, func(t *testing.T) {
			rec := endeRec(f.mitEnde)
			if f.aufErr {
				rec.Sessions[0].Interactions[1] = aufgezeichnet(2, "B", fehlerMit("42P01", "aufgezeichnet"), bereitMit("I"))
			}
			sitz := spielerFolge(false, f.folge...)
			ziel := &fakeZiel{sessions: []*fakeEinspiel{sitz}}
			z, err := spieleZaehler(context.Background(), t, rec, ziel, f.o)
			if got := texte(err); !reflect.DeepEqual(got, []string{f.want}) || exitVon(err) != 4 {
				t.Fatalf("Meldungen %q, Exit-Code %d, erwartet [%q] mit 4", got, exitVon(err), f.want)
			}
			if got := ziel.liste(); !reflect.DeepEqual(got, f.ablauf) {
				t.Fatalf("Ablauf %q, erwartet %q", got, f.ablauf)
			}
			if zu, ohne := sitz.war(); !zu || !ohne {
				t.Fatalf("geschlossen %v, ohne Terminate %v, erwartet beides", zu, ohne)
			}
			if z.Verglichen != 0 || z.Abweichend != 0 {
				t.Fatalf("Zähler %+v, erwartet verglichen 0, abweichend 0", z)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-25/Boundary — ohne Vergleich bleibt ein
// aufgezeichnetes Ende ohne Fehlerantwort des Servers erwartet-frei: Antwortet
// der Server normal, ist das kein Fehler, die Session endet mit Terminate; ein
// Verbindungsverlust trotz aufgezeichnetem Ende ist PGR-E4003 und die
// Verbindung wird mit Terminate geschlossen; eine Fehlerantwort ERROR in einer
// Interaktion mit Ende ist mit --allow-recorded-errors erwartet, die Interaktion
// läuft bis ReadyForQuery und die Session endet mit Terminate (LH-FA-20.a
// *Interaktion*, Tabelle der Fehlerregeln, LH-FA-25.a).
func TestPlayEndeOhneFehlerantwort(t *testing.T) {
	ok := []model.Response{befehl("SELECT 1"), bereitMit("I")}
	for _, f := range []struct {
		name   string
		folge  []schritt
		want   []string
		exit   int
		ablauf []string
	}{
		{
			name:   "normale Antwort",
			folge:  schritte(ok...),
			ablauf: []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "naechste", "schliesse"},
		},
		{
			name:   "Verbindungsverlust",
			folge:  []schritt{{err: model.Errorf(model.CodeConnectionLost, nil, "weg")}},
			want:   []string{"Netzwerk [PGR-E4003]: Session 1, Interaktion 2: weg"},
			exit:   4,
			ablauf: []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "schliesse"},
		},
		{
			name:   "ERROR",
			folge:  schritte(fehlerMit("42P01", "x"), bereitMit("I")),
			ablauf: []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "naechste", "schliesse"},
		},
	} {
		t.Run(f.name, func(t *testing.T) {
			sitz := &fakeEinspiel{antworten: [][]schritt{schritte(ok...), f.folge}}
			ziel := &fakeZiel{sessions: []*fakeEinspiel{sitz}}
			rec := endeRec(true)
			rec.Sessions = rec.Sessions[:1]
			z, err := spieleZaehler(context.Background(), t, rec, ziel, services.PlayOptions{AllowRecordedErrors: true})
			if got := texte(err); !reflect.DeepEqual(got, f.want) || exitVon(err) != f.exit {
				t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit %d", got, exitVon(err), f.want, f.exit)
			}
			if got := ziel.liste(); !reflect.DeepEqual(got, f.ablauf) {
				t.Fatalf("Ablauf %q, erwartet %q", got, f.ablauf)
			}
			if zu, ohne := sitz.war(); !zu || ohne {
				t.Fatalf("geschlossen %v, ohne Terminate %v, erwartet mit Terminate", zu, ohne)
			}
			if z.Verglichen != 0 || z.Abweichend != 0 {
				t.Fatalf("Zähler %+v, erwartet verglichen 0, abweichend 0", z)
			}
		})
	}
}
