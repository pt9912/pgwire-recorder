package services_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// cm ist eine Client-Nachricht des Typs typ ohne Felder.
func cm(typ model.ClientMessageType) model.ClientMessage {
	return model.ClientMessage{Type: typ}
}

// beschreibe ist ein Describe der Zielart ziel.
func beschreibe(ziel model.Target) model.ClientMessage {
	return model.ClientMessage{Type: model.ClientDescribe, Target: ziel}
}

// gruppe ist eine Gruppe aus den Client-Nachrichten client und den
// aufgezeichneten Server-Nachrichten server.
func gruppe(server []model.Response, client ...model.ClientMessage) model.Group {
	return model.Group{Client: client, Server: server}
}

// bereit ist ein ReadyForQuery als aufgezeichnete Server-Nachricht.
func bereit() []model.Response {
	return []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
}

// extendedInteraktion ist eine Extended-Interaktion aus den Gruppen.
func extendedInteraktion(seq int, gruppen ...model.Group) model.Interaction {
	return model.Interaction{Sequence: seq, Request: model.Request{Type: model.RequestExtended}, Groups: gruppen}
}

// syncInteraktion ist eine Extended-Interaktion aus einer Gruppe nur mit Sync.
func syncInteraktion(seq int) model.Interaction {
	return extendedInteraktion(seq, gruppe(bereit(), cm(model.ClientSync)))
}

// antwort ist eine Server-Antwort des Typs typ als Schritt.
func antwort(typ model.ResponseType) schritt {
	return schritt{r: model.Response{Type: typ, TxStatus: "I"}}
}

// sitzung ist eine Aufzeichnung mit einer Session aus den Interaktionen.
func sitzung(interaktionen ...model.Interaction) model.Recording {
	rec := model.NewRecording()
	rec.Sessions = []model.Session{{ID: 1, Interactions: interaktionen}}
	return rec
}

// Abdeckung: LH-FA-20/Happy, LH-FA-18/Happy — play sendet die Gruppen einer
// Extended-Interaktion in der aufgezeichneten Reihenfolge; nach einer Gruppe mit
// Flush liest es, bevor es die nächste sendet, die Antworten auf jede
// Client-Nachricht der Gruppe: eine auf Parse, Bind, Close, Execute und auf
// Describe eines Portals, zwei auf Describe einer Anweisung, keine auf Flush;
// DataRow, NoticeResponse und ParameterStatus zählen nicht, und die
// aufgezeichneten Server-Nachrichten bestimmen das Warten nicht; nach der Gruppe
// mit Sync liest es bis zum ReadyForQuery, bevor die nächste Interaktion
// beginnt (LH-FA-20.a *Gruppen*).
func TestPlayExtendedWarten(t *testing.T) {
	g1 := gruppe([]model.Response{{Type: model.ResponseParseComplete}}, cm(model.ClientParse), beschreibe(model.TargetStatement), cm(model.ClientFlush))
	g2 := gruppe(nil, cm(model.ClientBind), cm(model.ClientExecute), beschreibe(model.TargetPortal), model.ClientMessage{Type: model.ClientClose, Target: model.TargetPortal}, cm(model.ClientFlush))
	g3 := gruppe(bereit(), cm(model.ClientSync))
	rec := sitzung(extendedInteraktion(1, g1, g2, g3), interaktion(2, "B", "SELECT 1"))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{gruppen: [][]schritt{
		{antwort(model.ResponseParseComplete), antwort(model.ResponseNoticeResponse), antwort(model.ResponseParameterDescription), antwort(model.ResponseParameterStatus), antwort(model.ResponseNoData)},
		{antwort(model.ResponseBindComplete), antwort(model.ResponseDataRow), antwort(model.ResponseRowDescription), antwort(model.ResponseDataRow), antwort(model.ResponseCommandComplete), antwort(model.ResponseCloseComplete)},
		{antwort(model.ResponseReadyForQuery)},
	}}}}
	if err := spiele(context.Background(), t, nil, rec, ziel, "", ""); err != nil {
		t.Fatalf("Play: %v", err)
	}
	want := []string{
		"verbinde",
		"gruppe parse,describe,flush", "naechste", "naechste", "naechste", "naechste", "naechste",
		"gruppe bind,execute,describe,close,flush", "naechste", "naechste", "naechste", "naechste", "naechste", "naechste",
		"gruppe sync", "naechste",
		"anfrage B", "naechste", "schliesse",
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Boundary — eine Antwort anderer Art als erwartet ändert
// das Warten nicht, gezählt wird je Antwort (LH-FA-20.a *Gruppen*).
func TestPlayExtendedAndereArt(t *testing.T) {
	g1 := gruppe(nil, cm(model.ClientExecute), cm(model.ClientFlush))
	g2 := gruppe(bereit(), cm(model.ClientSync))
	rec := sitzung(extendedInteraktion(1, g1, g2))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{gruppen: [][]schritt{
		{antwort(model.ResponsePortalSuspended)},
		{antwort(model.ResponseReadyForQuery)},
	}}}}
	if err := spiele(context.Background(), t, nil, rec, ziel, "", ""); err != nil {
		t.Fatalf("Play: %v", err)
	}
	want := []string{"verbinde", "gruppe execute,flush", "naechste", "gruppe sync", "naechste", "schliesse"}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Boundary — eine Gruppe nur mit Flush wartet auf keine
// Antwort (LH-FA-20.a *Gruppen*: auf Flush nichts).
func TestPlayExtendedNurFlush(t *testing.T) {
	g1 := gruppe(nil, cm(model.ClientFlush))
	g2 := gruppe(bereit(), cm(model.ClientSync))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{gruppen: [][]schritt{nil, {antwort(model.ResponseReadyForQuery)}}}}}
	if err := spiele(context.Background(), t, nil, sitzung(extendedInteraktion(1, g1, g2)), ziel, "", ""); err != nil {
		t.Fatalf("Play: %v", err)
	}
	want := []string{"verbinde", "gruppe flush", "gruppe sync", "naechste", "schliesse"}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// zweiGruppen ist eine Extended-Interaktion aus einer Flush-Gruppe mit Parse
// und Bind (zwei Antworten), einer mit Execute (eine Antwort) und der Gruppe
// mit Sync, in der aufgezeichneten Antworten stehen.
func zweiGruppen(aufgezeichnet ...model.Response) model.Interaction {
	return extendedInteraktion(1,
		gruppe(aufgezeichnet, cm(model.ClientParse), cm(model.ClientBind), cm(model.ClientFlush)),
		gruppe(nil, cm(model.ClientExecute), cm(model.ClientFlush)),
		gruppe(bereit(), cm(model.ClientSync)),
	)
}

// Abdeckung: LH-FA-20/Negative — eine Fehlerantwort in einer
// Extended-Interaktion bricht ab, wenn nichts anderes gilt: PGR-E4004 mit
// Session, Interaktion, SQLSTATE und Meldung, in der ersten, einer späteren
// Gruppe mit Flush und der Gruppe mit Sync; danach keine weitere Antwort
// gelesen, keine weitere Gruppe und keine weitere Interaktion, die Verbindung mit
// Schliesse beendet (LH-FA-20.a Schritt 6, *Interaktion*).
func TestPlayExtendedFehlerBrichtAb(t *testing.T) {
	for _, f := range []struct {
		name   string
		gruppe [][]schritt
		ablauf []string
	}{
		{
			"erste Gruppe",
			[][]schritt{{fehler("42P01", "a")}},
			[]string{"verbinde", "gruppe parse,bind,flush", "naechste", "schliesse"},
		},
		{
			"spätere Gruppe",
			[][]schritt{{antwort(model.ResponseParseComplete), antwort(model.ResponseBindComplete)}, {fehler("42P01", "a")}},
			[]string{"verbinde", "gruppe parse,bind,flush", "naechste", "naechste", "gruppe execute,flush", "naechste", "schliesse"},
		},
		{
			"Gruppe mit Sync",
			[][]schritt{{antwort(model.ResponseParseComplete), antwort(model.ResponseBindComplete)}, {antwort(model.ResponseCommandComplete)}, {fehler("42P01", "a")}},
			[]string{"verbinde", "gruppe parse,bind,flush", "naechste", "naechste", "gruppe execute,flush", "naechste", "gruppe sync", "naechste", "schliesse"},
		},
	} {
		t.Run(f.name, func(t *testing.T) {
			rec := sitzung(zweiGruppen(), interaktion(2, "B", "SELECT 1"))
			ziel := &fakeZiel{sessions: []*fakeEinspiel{{gruppen: f.gruppe}}}
			err := spiele(context.Background(), t, nil, rec, ziel, "", "")
			if want := []string{e4004("Session 1, Interaktion 1", "42P01", "a")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 4 {
				t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 4", texte(err), exitVon(err), want)
			}
			if got := ziel.liste(); !reflect.DeepEqual(got, f.ablauf) {
				t.Fatalf("Ablauf %q, erwartet %q", got, f.ablauf)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — mit --continue-on-error
// zählt play nach einer Fehlerantwort in einer Extended-Interaktion nicht mehr:
// Es sendet die übrigen Gruppen ohne Warten dazwischen und liest dann bis zum
// ReadyForQuery, auch nach einer Fehlerantwort in der Gruppe mit Sync; die
// nächste Interaktion läuft, je Fehlerantwort ein PGR-E4004, Exit-Code 4
// (LH-FA-20.a Schritt 6, *Gruppen*).
func TestPlayExtendedFortsetzung(t *testing.T) {
	for _, f := range []struct {
		name   string
		gruppe [][]schritt
		ablauf []string
	}{
		{
			"Fehler in der ersten Gruppe",
			[][]schritt{{fehler("42P01", "a")}, nil, {antwort(model.ResponseNoticeResponse), antwort(model.ResponseReadyForQuery)}},
			[]string{"verbinde", "gruppe parse,bind,flush", "naechste", "gruppe execute,flush", "gruppe sync", "naechste", "naechste", "anfrage B", "naechste", "schliesse"},
		},
		{
			"Fehler in einer späteren Gruppe",
			[][]schritt{{antwort(model.ResponseParseComplete), antwort(model.ResponseBindComplete)}, {fehler("42P01", "a")}, {antwort(model.ResponseReadyForQuery)}},
			[]string{"verbinde", "gruppe parse,bind,flush", "naechste", "naechste", "gruppe execute,flush", "naechste", "gruppe sync", "naechste", "anfrage B", "naechste", "schliesse"},
		},
		{
			"Fehler in der Gruppe mit Sync",
			[][]schritt{{antwort(model.ResponseParseComplete), antwort(model.ResponseBindComplete)}, {antwort(model.ResponseCommandComplete)}, {fehler("42P01", "a"), antwort(model.ResponseReadyForQuery)}},
			[]string{"verbinde", "gruppe parse,bind,flush", "naechste", "naechste", "gruppe execute,flush", "naechste", "gruppe sync", "naechste", "naechste", "anfrage B", "naechste", "schliesse"},
		},
	} {
		t.Run(f.name, func(t *testing.T) {
			rec := sitzung(zweiGruppen(), interaktion(2, "B", "SELECT 1"))
			ziel := &fakeZiel{sessions: []*fakeEinspiel{{gruppen: f.gruppe}}}
			err := spieleMit(context.Background(), t, nil, rec, ziel, services.PlayOptions{ContinueOnError: true})
			if want := []string{e4004("Session 1, Interaktion 1", "42P01", "a")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 4 {
				t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 4", texte(err), exitVon(err), want)
			}
			if got := ziel.liste(); !reflect.DeepEqual(got, f.ablauf) {
				t.Fatalf("Ablauf %q, erwartet %q", got, f.ablauf)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — ein Fehler, der abbricht, bricht in einer
// Extended-Interaktion auch mit --continue-on-error ab: eine Fehlerantwort mit
// FATAL ist PGR-E4003, auch wenn die Aufzeichnung eine error_response trägt, und
// es wird keine weitere Gruppe gesendet und keine weitere Antwort gelesen
// (LH-FA-20.a *Interaktion*).
func TestPlayExtendedFatal(t *testing.T) {
	fatal := fehlerantwort(map[string]string{"S": "FATAL", "V": "FATAL", "C": "57P01", "M": "beendet"})
	rec := sitzung(zweiGruppen(model.Response{Type: model.ResponseErrorResponse, Fields: map[string]string{"C": "XX000", "M": "aufgezeichnet"}}))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{gruppen: [][]schritt{{antwort(model.ResponseParseComplete), fatal}}}}}
	err := spieleMit(context.Background(), t, nil, rec, ziel, services.PlayOptions{ContinueOnError: true, AllowRecordedErrors: true})
	if want := []string{"Netzwerk [PGR-E4003]: Session 1, Interaktion 1: Fehlerantwort des Servers 57P01 „beendet“, die Verbindung endet"}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 4 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 4", texte(err), exitVon(err), want)
	}
	if got, want := ziel.liste(), []string{"verbinde", "gruppe parse,bind,flush", "naechste", "naechste", "schliesse"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Boundary — mit --allow-recorded-errors
// ist eine Fehlerantwort in einer Extended-Interaktion erwartet, wenn die
// aufgezeichnete Interaktion in irgendeiner Gruppe eine error_response trägt: in
// der ersten, einer mittleren oder der letzten, gleich in welcher Gruppe der
// Server den Fehler sendet und mit welchem SQLSTATE; keine Meldung, Exit-Code 0,
// und die nächste Interaktion läuft ohne --continue-on-error. Trägt keine Gruppe
// eine, ist sie PGR-E4004 (LH-FA-20.a *Interaktion*).
func TestPlayExtendedErwarteterFehler(t *testing.T) {
	fehlerAufgezeichnet := []model.Response{{Type: model.ResponseErrorResponse, Fields: map[string]string{"C": "XX000", "M": "anders"}}}
	drei := func(aufgezeichnet1, aufgezeichnet2, aufgezeichnet3 []model.Response) model.Interaction {
		return extendedInteraktion(1,
			gruppe(aufgezeichnet1, cm(model.ClientParse), cm(model.ClientFlush)),
			gruppe(aufgezeichnet2, cm(model.ClientBind), cm(model.ClientFlush)),
			gruppe(append(aufgezeichnet3, bereit()...), cm(model.ClientSync)),
		)
	}
	for _, f := range []struct {
		name string
		in   model.Interaction
		want []string
	}{
		{"erste Gruppe", drei(fehlerAufgezeichnet, nil, nil), nil},
		{"mittlere Gruppe", drei(nil, fehlerAufgezeichnet, nil), nil},
		{"letzte Gruppe", drei(nil, nil, fehlerAufgezeichnet), nil},
		{"keine Gruppe", drei(nil, nil, nil), []string{e4004("Session 1, Interaktion 1", "42P01", "a")}},
	} {
		t.Run(f.name, func(t *testing.T) {
			rec := sitzung(f.in, interaktion(2, "B", "SELECT 1"))
			ziel := &fakeZiel{sessions: []*fakeEinspiel{{gruppen: [][]schritt{{fehler("42P01", "a")}, nil, {antwort(model.ResponseReadyForQuery)}}}}}
			err := spieleMit(context.Background(), t, nil, rec, ziel, services.PlayOptions{AllowRecordedErrors: true})
			if !reflect.DeepEqual(texte(err), f.want) {
				t.Fatalf("Meldungen %q, erwartet %q", texte(err), f.want)
			}
			ablauf := ziel.liste()
			if f.want == nil {
				if want := []string{"verbinde", "gruppe parse,flush", "naechste", "gruppe bind,flush", "gruppe sync", "naechste", "anfrage B", "naechste", "schliesse"}; !reflect.DeepEqual(ablauf, want) {
					t.Fatalf("Ablauf %q, erwartet %q", ablauf, want)
				}
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — das erste Abbruchsignal (ctx endet) in einer
// Extended-Interaktion lässt sie bis zu ihrem ReadyForQuery laufen: alle
// übrigen Gruppen werden gesendet und gelesen; danach beginnt keine weitere
// Interaktion, die Verbindung wird geschlossen, ohne Fehler (LH-FA-20.a
// *Abbruchsignal*).
func TestPlayExtendedSignal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rec := sitzung(zweiGruppen(), interaktion(2, "B", "SELECT 1"))
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{
		beiGruppe: func(n int) {
			if n == 1 {
				cancel()
			}
		},
		gruppen: [][]schritt{{antwort(model.ResponseParseComplete), antwort(model.ResponseBindComplete)}, {antwort(model.ResponseCommandComplete)}, {antwort(model.ResponseReadyForQuery)}},
	}}}
	if err := spiele(ctx, t, nil, rec, ziel, "", ""); err != nil {
		t.Fatalf("Play: %v", err)
	}
	want := []string{"verbinde", "gruppe parse,bind,flush", "naechste", "naechste", "gruppe execute,flush", "naechste", "gruppe sync", "naechste", "schliesse"}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Negative — scheitert das Senden einer Gruppe oder das
// Lesen in einer Extended-Interaktion, bricht play mit dem Code des Fehlers ab
// (PGR-E4003 beim Senden und beim Verbindungsende, PGR-E6001 bei einer nicht
// lesbaren Antwort) und mit dem Ort; es sendet keine weitere Gruppe und liest
// nicht weiter (LH-FA-20.a *Interaktion*).
func TestPlayExtendedSendenLesenScheitert(t *testing.T) {
	zweite := &fakeEinspiel{gruppen: [][]schritt{{antwort(model.ResponseParseComplete), antwort(model.ResponseBindComplete)}}}
	zweite.beiGruppe = func(n int) {
		if n == 2 {
			zweite.gruppeFehler = model.Errorf(model.CodeConnectionLost, nil, "senden")
		}
	}
	for _, f := range []struct {
		name   string
		sess   *fakeEinspiel
		code   string
		ablauf []string
	}{
		{
			"Senden der ersten Gruppe",
			&fakeEinspiel{gruppeFehler: model.Errorf(model.CodeConnectionLost, nil, "senden")},
			model.CodeConnectionLost,
			[]string{"verbinde", "gruppe parse,bind,flush", "schliesse"},
		},
		{
			"Senden der zweiten Gruppe",
			zweite,
			model.CodeConnectionLost,
			[]string{"verbinde", "gruppe parse,bind,flush", "naechste", "naechste", "gruppe execute,flush", "schliesse"},
		},
		{
			"Verbindungsende beim Lesen",
			&fakeEinspiel{gruppen: [][]schritt{{{err: model.Errorf(model.CodeConnectionLost, nil, "senden")}}}},
			model.CodeConnectionLost,
			[]string{"verbinde", "gruppe parse,bind,flush", "naechste", "schliesse"},
		},
		{
			"Antwort nicht lesbar",
			&fakeEinspiel{gruppen: [][]schritt{{{err: model.Errorf(model.CodeUnsupported, nil, "senden")}}}},
			model.CodeUnsupported,
			[]string{"verbinde", "gruppe parse,bind,flush", "naechste", "schliesse"},
		},
	} {
		t.Run(f.name, func(t *testing.T) {
			ziel := &fakeZiel{sessions: []*fakeEinspiel{f.sess}}
			err := spiele(context.Background(), t, nil, sitzung(zweiGruppen(), interaktion(2, "B", "SELECT 1")), ziel, "", "")
			if codeVon(err) != f.code || !strings.Contains(err.Error(), "Session 1, Interaktion 1: senden") {
				t.Fatalf("%v, erwartet %s mit Ort", err, f.code)
			}
			if got := ziel.liste(); !reflect.DeepEqual(got, f.ablauf) {
				t.Fatalf("Ablauf %q, erwartet %q", got, f.ablauf)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Happy — einfache Anfragen und Extended-Interaktionen in
// einer Session laufen in der aufgezeichneten Reihenfolge nacheinander, über
// mehrere Sessions je eine eigene Verbindung (LH-FA-20.a Schritt 3).
func TestPlayExtendedGemischt(t *testing.T) {
	rec := model.NewRecording()
	rec.Sessions = []model.Session{
		{ID: 1, Interactions: []model.Interaction{interaktion(1, "A", "SELECT 1"), syncInteraktion(2), interaktion(3, "C", "SELECT 1")}},
		{ID: 2, Interactions: []model.Interaction{syncInteraktion(1)}},
	}
	ziel := &fakeZiel{}
	if err := spiele(context.Background(), t, nil, rec, ziel, "", ""); err != nil {
		t.Fatalf("Play: %v", err)
	}
	want := []string{
		"verbinde", "anfrage A", "naechste", "gruppe sync", "naechste", "anfrage C", "naechste", "schliesse",
		"verbinde", "gruppe sync", "naechste", "schliesse",
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}
