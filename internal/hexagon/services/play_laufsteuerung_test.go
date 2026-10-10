package services_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// spieleMit startet den PlayService mit rec, ziel und den Optionen o und
// spielt ein.
func spieleMit(ctx context.Context, t *testing.T, ablauf <-chan struct{}, rec model.Recording, ziel *fakeZiel, o services.PlayOptions) error {
	t.Helper()
	s, err := services.NewPlayService(context.WithoutCancel(ctx), ladeRepo{rec: rec}, "r.yaml", ziel, o)
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	return s.Play(ctx, ablauf)
}

// texte sind die Fehlertexte der Meldungen von err in ihrer Reihenfolge, wie
// der Bootstrap sie als Zeilen error schreibt (model.Meldungen).
func texte(err error) []string {
	var out []string
	for _, m := range model.Meldungen(err) {
		out = append(out, m.Text)
	}
	return out
}

// exitVon ist der Exit-Code, den der Bootstrap aus err bildet: der der ersten
// Meldung, ohne Meldung 0.
func exitVon(err error) int {
	ms := model.Meldungen(err)
	if len(ms) == 0 {
		return 0
	}
	return ms[0].ExitCode()
}

// fehler ist eine Fehlerantwort mit dem Schweregrad ERROR.
func fehler(code, text string) schritt {
	return fehlerantwort(map[string]string{"S": "ERROR", "V": "ERROR", "C": code, "M": text})
}

// e4004 ist der Fehlertext eines PGR-E4004 an Session und Interaktion.
func e4004(ort, code, text string) string {
	return "Netzwerk [PGR-E4004]: " + ort + ": Fehlerantwort des Servers " + code + " „" + text + "“"
}

// drei ist eine Session mit drei Anfragen und eine zweite mit einer.
func drei() model.Recording {
	rec := model.NewRecording()
	rec.Sessions = []model.Session{anfragen(1, nil, "A", "B", "C"), anfragen(2, nil, "D")}
	return rec
}

// Abdeckung: LH-FA-20/Boundary — mit --continue-on-error läuft das Einspielen
// nach einer Fehlerantwort weiter: Es liest die Interaktion bis zum
// ReadyForQuery, sendet die nächste Interaktion und die nächste Session; je
// Fehlerantwort ein PGR-E4004 mit Session, Interaktion, SQLSTATE und Meldung,
// auch zwei in einer Interaktion, in der Reihenfolge ihres Auftretens;
// Exit-Code 4 (LH-FA-20.a Schritt 6, *Meldungen*).
func TestPlayFortsetzung(t *testing.T) {
	ziel := &fakeZiel{sessions: []*fakeEinspiel{
		{antworten: [][]schritt{{fehler("42P01", "a1"), {r: model.Response{Type: model.ResponseNoticeResponse}}, fehler("42P02", "a2"), bereitZuEnde()}, nil, {fehler("22012", "c")}}},
		{antworten: [][]schritt{{fehler("23505", "d")}}},
	}}
	err := spieleMit(context.Background(), t, nil, drei(), ziel, services.PlayOptions{ContinueOnError: true})
	want := []string{
		e4004("Session 1, Interaktion 1", "42P01", "a1"),
		e4004("Session 1, Interaktion 1", "42P02", "a2"),
		e4004("Session 1, Interaktion 3", "22012", "c"),
		e4004("Session 2, Interaktion 1", "23505", "d"),
	}
	if got := texte(err); !reflect.DeepEqual(got, want) || exitVon(err) != 4 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 4", got, exitVon(err), want)
	}
	ablauf := []string{
		"verbinde", "anfrage A", "naechste", "naechste", "naechste", "naechste", "anfrage B", "naechste", "anfrage C", "naechste", "naechste", "schliesse",
		"verbinde", "anfrage D", "naechste", "naechste", "schliesse",
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, ablauf) {
		t.Fatalf("Ablauf %q, erwartet %q", got, ablauf)
	}
}

// Abdeckung: LH-FA-20/Negative — ein Fehler, der abbricht, bricht auch mit
// --continue-on-error ab: beim Lesen (PGR-E6001, Exit-Code 6), eine
// Fehlerantwort mit FATAL (PGR-E4003), das Verbindungsende beim Weiterlesen
// nach einer Fehlerantwort (PGR-E4003) und der Aufbau einer späteren Session
// (PGR-E4002); seine Meldung steht zuerst, danach die der früheren PGR-E4004
// in der Reihenfolge ihres Auftretens, und der Exit-Code ist der seiner
// Klasse; danach keine weitere Anfrage und keine weitere Session (LH-FA-20.a
// *Meldungen*, *Exit-Code*).
func TestPlayFortsetzungAbbruch(t *testing.T) {
	fatal := fehlerantwort(map[string]string{"S": "FATAL", "V": "FATAL", "C": "57P01", "M": "beendet"})
	for _, f := range []struct {
		name      string
		sessionen []*fakeEinspiel
		aufbau2   bool
		want      []string
		exit      int
		ablauf    []string
	}{
		{
			name:      "PGR-E6001 in Session 2",
			sessionen: []*fakeEinspiel{{antworten: [][]schritt{{fehler("42P01", "a")}, {fehler("42P02", "b")}}}, {antworten: [][]schritt{{{err: model.Errorf(model.CodeUnsupported, nil, "lesen")}}}}},
			want:      []string{"nicht unterstützt [PGR-E6001]: Session 2, Interaktion 1: lesen", e4004("Session 1, Interaktion 1", "42P01", "a"), e4004("Session 1, Interaktion 2", "42P02", "b")},
			exit:      6,
			ablauf:    []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "naechste", "anfrage C", "naechste", "schliesse", "verbinde", "anfrage D", "naechste", "schliesse"},
		},
		{
			name:      "FATAL",
			sessionen: []*fakeEinspiel{{antworten: [][]schritt{{fehler("42P01", "a")}, {fatal}}}},
			want:      []string{"Netzwerk [PGR-E4003]: Session 1, Interaktion 2: Fehlerantwort des Servers 57P01 „beendet“, die Verbindung endet", e4004("Session 1, Interaktion 1", "42P01", "a")},
			exit:      4,
			ablauf:    []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "schliesse"},
		},
		{
			name:      "Verbindungsende nach der Fehlerantwort",
			sessionen: []*fakeEinspiel{{antworten: [][]schritt{{fehler("42P01", "a"), {err: model.Errorf(model.CodeConnectionLost, nil, "weg")}}}}},
			want:      []string{"Netzwerk [PGR-E4003]: Session 1, Interaktion 1: weg", e4004("Session 1, Interaktion 1", "42P01", "a")},
			exit:      4,
			ablauf:    []string{"verbinde", "anfrage A", "naechste", "naechste", "schliesse"},
		},
		{
			name:      "Aufbau von Session 2",
			sessionen: []*fakeEinspiel{{antworten: [][]schritt{nil, {fehler("42P01", "b")}}}},
			aufbau2:   true,
			want:      []string{"Netzwerk [PGR-E4002]: Session 2: weg", e4004("Session 1, Interaktion 2", "42P01", "b")},
			exit:      4,
			ablauf:    []string{"verbinde", "anfrage A", "naechste", "anfrage B", "naechste", "naechste", "anfrage C", "naechste", "schliesse", "verbinde"},
		},
	} {
		t.Run(f.name, func(t *testing.T) {
			ziel := &fakeZiel{sessions: f.sessionen}
			if f.aufbau2 {
				f.sessionen[0].beiAnfrage = func(string) { ziel.verbindeFehler = model.Errorf(model.CodeUpstream, nil, "weg") }
			}
			err := spieleMit(context.Background(), t, nil, drei(), ziel, services.PlayOptions{ContinueOnError: true})
			if got := texte(err); !reflect.DeepEqual(got, f.want) || exitVon(err) != f.exit {
				t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit %d", got, exitVon(err), f.want, f.exit)
			}
			if got := ziel.liste(); !reflect.DeepEqual(got, f.ablauf) {
				t.Fatalf("Ablauf %q, erwartet %q", got, f.ablauf)
			}
		})
	}
}

// mitAufgezeichnetemFehler ist eine Session aus A, deren Aufzeichnung an
// zweiter Stelle eine error_response mit einem anderen SQLSTATE trägt, als der
// Server sendet, und B und C ohne error_response.
func mitAufgezeichnetemFehler() model.Recording {
	a := model.Interaction{Sequence: 1, Request: model.Request{Type: model.RequestQuery, SQL: "A"}, Responses: []model.Response{
		{Type: model.ResponseRowDescription},
		{Type: model.ResponseErrorResponse, Fields: map[string]string{"S": "ERROR", "C": "XX000", "M": "anders"}},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}}
	rec := model.NewRecording()
	rec.Sessions = []model.Session{{ID: 1, Interactions: []model.Interaction{a, interaktion(2, "B", "SELECT 1"), interaktion(3, "C", "SELECT 1")}}}
	return rec
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-20/Negative — mit
// --allow-recorded-errors ist eine Fehlerantwort in einer Interaktion, deren
// Aufzeichnung an irgendeiner Stelle eine error_response trägt, erwartet,
// gleich welcher SQLSTATE: keine Meldung, kein Abbruch, das Lesen bis zum
// ReadyForQuery und die nächste Interaktion laufen, auch ohne
// --continue-on-error. Trägt die Aufzeichnung keine, ist sie PGR-E4004 wie
// ohne Option; ohne Option ist auch die erste PGR-E4004; FATAL ist PGR-E4003,
// auch wenn sie erwartet wäre; mit --continue-on-error zählt nur die
// unerwartete (LH-FA-20.a *Interaktion*, Schritt 6).
func TestPlayErwarteterFehler(t *testing.T) {
	fatal := fehlerantwort(map[string]string{"S": "FATAL", "C": "57P01", "M": "beendet"})
	for _, f := range []struct {
		name      string
		o         services.PlayOptions
		antworten [][]schritt
		want      []string
		ablauf    []string
	}{
		{
			name:      "erwartet",
			o:         services.PlayOptions{AllowRecordedErrors: true},
			antworten: [][]schritt{{fehler("42P01", "a"), bereitZuEnde()}},
			ablauf:    []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "anfrage C", "naechste", "schliesse"},
		},
		{
			name:      "ohne aufgezeichnete error_response",
			o:         services.PlayOptions{AllowRecordedErrors: true},
			antworten: [][]schritt{nil, {fehler("42P02", "b"), bereitZuEnde()}},
			want:      []string{e4004("Session 1, Interaktion 2", "42P02", "b")},
			ablauf:    []string{"verbinde", "anfrage A", "naechste", "anfrage B", "naechste", "schliesse"},
		},
		{
			name:      "ohne Option",
			antworten: [][]schritt{{fehler("42P01", "a"), bereitZuEnde()}},
			want:      []string{e4004("Session 1, Interaktion 1", "42P01", "a")},
			ablauf:    []string{"verbinde", "anfrage A", "naechste", "schliesse"},
		},
		{
			name:      "FATAL",
			o:         services.PlayOptions{AllowRecordedErrors: true, ContinueOnError: true},
			antworten: [][]schritt{{fatal}},
			want:      []string{"Netzwerk [PGR-E4003]: Session 1, Interaktion 1: Fehlerantwort des Servers 57P01 „beendet“, die Verbindung endet"},
			ablauf:    []string{"verbinde", "anfrage A", "naechste", "schliesse"},
		},
		{
			name:      "mit --continue-on-error",
			o:         services.PlayOptions{AllowRecordedErrors: true, ContinueOnError: true},
			antworten: [][]schritt{{fehler("42P01", "a")}, {fehler("42P02", "b")}},
			want:      []string{e4004("Session 1, Interaktion 2", "42P02", "b")},
			ablauf:    []string{"verbinde", "anfrage A", "naechste", "naechste", "anfrage B", "naechste", "naechste", "anfrage C", "naechste", "schliesse"},
		},
	} {
		t.Run(f.name, func(t *testing.T) {
			ziel := &fakeZiel{sessions: []*fakeEinspiel{{antworten: f.antworten}}}
			err := spieleMit(context.Background(), t, nil, mitAufgezeichnetemFehler(), ziel, f.o)
			if got := texte(err); !reflect.DeepEqual(got, f.want) {
				t.Fatalf("Meldungen %q, erwartet %q", got, f.want)
			}
			if got := ziel.liste(); !reflect.DeepEqual(got, f.ablauf) {
				t.Fatalf("Ablauf %q, erwartet %q", got, f.ablauf)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — mit --finish-session-on-interrupt laufen nach
// dem ersten Signal alle Interaktionen der laufenden Session, auch nach einem
// Signal im Aufbau, dessen ctx nicht endet; danach beginnt keine weitere
// Session, auch wenn das Signal in ihrer letzten Interaktion eintraf; ohne
// Fehler ist der Exit-Code 0 (LH-FA-20.a Schritt 7, *Abbruchsignal*).
func TestPlayFinishSession(t *testing.T) {
	o := services.PlayOptions{FinishSessionOnInterrupt: true}
	ganz := []string{"verbinde", "anfrage A", "naechste", "anfrage B", "naechste", "anfrage C", "naechste", "schliesse"}

	ctx, cancel := context.WithCancel(context.Background())
	sitzung := &fakeEinspiel{beiAnfrage: func(sql string) {
		if sql == "A" {
			cancel()
		}
	}}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{sitzung}}
	if err := spieleMit(ctx, t, nil, drei(), ziel, o); err != nil || exitVon(err) != 0 {
		t.Fatalf("Signal in der Interaktion: %v", err)
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, ganz) {
		t.Fatalf("Signal in der Interaktion: Ablauf %q, erwartet %q", got, ganz)
	}

	ctx, cancel = context.WithCancel(context.Background())
	ziel = &fakeZiel{beimVerbinden: cancel}
	if err := spieleMit(ctx, t, nil, drei(), ziel, o); err != nil {
		t.Fatalf("Signal im Aufbau: %v", err)
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, ganz) || ziel.verbindeCtxEnde {
		t.Fatalf("Signal im Aufbau: Ablauf %q, erwartet %q; ctx des Aufbaus beendet: %v", got, ganz, ziel.verbindeCtxEnde)
	}

	ctx, cancel = context.WithCancel(context.Background())
	sitzung = &fakeEinspiel{antworten: [][]schritt{nil, nil, {{r: model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"}, vorher: cancel}}}}
	ziel = &fakeZiel{sessions: []*fakeEinspiel{sitzung}}
	if err := spieleMit(ctx, t, nil, drei(), ziel, o); err != nil {
		t.Fatalf("Signal in der letzten Interaktion: %v", err)
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, ganz) {
		t.Fatalf("Signal in der letzten Interaktion: Ablauf %q, erwartet %q", got, ganz)
	}
}

// Abdeckung: LH-FA-20/Negative — in der Session, die mit
// --finish-session-on-interrupt nach dem ersten Signal zu Ende läuft, gelten
// die Fehlerregeln wie ohne Signal: Ohne --continue-on-error bricht eine
// Fehlerantwort ab, mit ihm läuft die Session weiter und zählt; der Exit-Code
// ist 4 (LH-FA-20.a *Abbruchsignal*, *Exit-Code*).
func TestPlayFinishSessionFehler(t *testing.T) {
	for _, f := range []struct {
		name   string
		o      services.PlayOptions
		ablauf []string
	}{
		{"ohne --continue-on-error", services.PlayOptions{FinishSessionOnInterrupt: true}, []string{"verbinde", "anfrage A", "naechste", "anfrage B", "naechste", "schliesse"}},
		{"mit --continue-on-error", services.PlayOptions{FinishSessionOnInterrupt: true, ContinueOnError: true}, []string{"verbinde", "anfrage A", "naechste", "anfrage B", "naechste", "naechste", "anfrage C", "naechste", "schliesse"}},
	} {
		t.Run(f.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			sitzung := &fakeEinspiel{antworten: [][]schritt{nil, {fehler("42P01", "b")}}, beiAnfrage: func(sql string) {
				if sql == "A" {
					cancel()
				}
			}}
			ziel := &fakeZiel{sessions: []*fakeEinspiel{sitzung}}
			err := spieleMit(ctx, t, nil, drei(), ziel, f.o)
			if want := []string{e4004("Session 1, Interaktion 2", "42P01", "b")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 4 {
				t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 4", texte(err), exitVon(err), want)
			}
			if got := ziel.liste(); !reflect.DeepEqual(got, f.ablauf) {
				t.Fatalf("Ablauf %q, erwartet %q", got, f.ablauf)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — trifft das erste Signal mit
// --continue-on-error zwischen einer Fehlerantwort und der Fortsetzung ein,
// liest play die Interaktion bis zum ReadyForQuery, beginnt keine weitere,
// und der Fehler zählt: Exit-Code 4 (LH-FA-20.a Tabellenzeile
// *Abbruchsignal*).
func TestPlaySignalNachFehlerantwort(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	f := fehler("42P01", "a")
	f.vorher = cancel
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{{f, {r: model.Response{Type: model.ResponseNoticeResponse}}, bereitZuEnde()}}}}}
	err := spieleMit(ctx, t, nil, drei(), ziel, services.PlayOptions{ContinueOnError: true})
	if want := []string{e4004("Session 1, Interaktion 1", "42P01", "a")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 4 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 4", texte(err), exitVon(err), want)
	}
	if got, want := ziel.liste(), []string{"verbinde", "anfrage A", "naechste", "naechste", "naechste", "schliesse"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Boundary — eine Fehlerantwort, die mit
// --continue-on-error vor dem zweiten Signal eintraf, zählt, auch wenn das
// Signal ihre Interaktion danach unterbricht; die Unterbrechung selbst ist
// kein Fehler (LH-FA-20.a *Meldungen*, *Abbruchsignal*).
func TestPlayZweitesSignalNachFehlerantwort(t *testing.T) {
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{{fehler("42P01", "a"), {blockiert: true}}}}}}
	s, err := services.NewPlayService(context.Background(), ladeRepo{rec: drei()}, "r.yaml", ziel, services.PlayOptions{ContinueOnError: true})
	if err != nil {
		t.Fatal(err)
	}
	ablauf := make(chan struct{})
	err = warteAufPlay(context.Background(), t, s, ablauf, func() {
		bis(t, func() bool { return len(ziel.liste()) >= 4 }, "Play wartet nicht binnen 30 s nach der Fehlerantwort auf die nächste Antwort")
		close(ablauf)
	})
	if want := []string{e4004("Session 1, Interaktion 1", "42P01", "a")}; !reflect.DeepEqual(texte(err), want) || exitVon(err) != 4 {
		t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit 4", texte(err), exitVon(err), want)
	}
	if got, want := ziel.liste(), []string{"verbinde", "anfrage A", "naechste", "naechste", "schliesse"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Boundary — auch mit --finish-session-on-interrupt
// schließt das zweite Signal die laufende Verbindung sofort, in einer
// Interaktion, die auf eine Antwort wartet, und bricht einen laufenden Aufbau
// ab; die unterbrochene Interaktion und der unterbrochene Aufbau sind kein
// Fehler, und keine weitere Interaktion oder Session beginnt. Ohne früheren
// Fehler ist der Exit-Code 0, nach einem früheren PGR-E4004 mit
// --continue-on-error 4 mit dessen Meldung als einziger (LH-FA-20.a
// *Abbruchsignal*, *Meldungen*).
func TestPlayZweitesSignalFinishSession(t *testing.T) {
	for _, f := range []struct {
		name   string
		o      services.PlayOptions
		ziel   func(erstesSignal func()) *fakeZiel
		warte  int
		want   []string
		exit   int
		ablauf []string
	}{
		{
			name: "in der Interaktion",
			o:    services.PlayOptions{FinishSessionOnInterrupt: true},
			ziel: func(erstesSignal func()) *fakeZiel {
				return &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{{{blockiert: true}}}, beiAnfrage: func(string) { erstesSignal() }}}}
			},
			warte:  3,
			ablauf: []string{"verbinde", "anfrage A", "naechste", "schliesse"},
		},
		{
			name: "nach früherem Fehler",
			o:    services.PlayOptions{FinishSessionOnInterrupt: true, ContinueOnError: true},
			ziel: func(erstesSignal func()) *fakeZiel {
				return &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{{fehler("42P01", "a"), {blockiert: true}}}, beiAnfrage: func(string) { erstesSignal() }}}}
			},
			warte:  4,
			want:   []string{e4004("Session 1, Interaktion 1", "42P01", "a")},
			exit:   4,
			ablauf: []string{"verbinde", "anfrage A", "naechste", "naechste", "schliesse"},
		},
		{
			name: "im Aufbau",
			o:    services.PlayOptions{FinishSessionOnInterrupt: true},
			ziel: func(erstesSignal func()) *fakeZiel {
				return &fakeZiel{verbindeBlockiert: true, beimVerbinden: erstesSignal}
			},
			warte:  1,
			ablauf: []string{"verbinde"},
		},
	} {
		t.Run(f.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ziel := f.ziel(cancel)
			s, err := services.NewPlayService(context.Background(), ladeRepo{rec: drei()}, "r.yaml", ziel, f.o)
			if err != nil {
				t.Fatal(err)
			}
			ablauf := make(chan struct{})
			err = warteAufPlay(ctx, t, s, ablauf, func() {
				bis(t, func() bool { return len(ziel.liste()) >= f.warte }, "Play wartet nicht binnen 30 s auf den Server")
				close(ablauf)
			})
			if !reflect.DeepEqual(texte(err), f.want) || exitVon(err) != f.exit {
				t.Fatalf("Meldungen %q, Exit-Code %d, erwartet %q mit %d", texte(err), exitVon(err), f.want, f.exit)
			}
			if got := ziel.liste(); !reflect.DeepEqual(got, f.ablauf) {
				t.Fatalf("Ablauf %q, erwartet %q", got, f.ablauf)
			}
		})
	}
}
