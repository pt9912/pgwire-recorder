package services_test

import (
	"context"
	"math"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
	"github.com/pt9912/pgwire-recorder/internal/testuhr"
	"github.com/pt9912/pgwire-recorder/internal/testzeit"
)

// zaehlUhr ist die Fake-Uhr, die die Aufrufe zählt: Lesungen von Jetzt, Aufrufe von
// WarteBis und davon die, die mit einem beendeten Kontext begannen.
type zaehlUhr struct {
	*testuhr.Uhr
	mu                          sync.Mutex
	lesungen, aufrufe, nachEnde int
}

func (u *zaehlUhr) Jetzt() int64 {
	u.mu.Lock()
	u.lesungen++
	u.mu.Unlock()
	return u.Uhr.Jetzt()
}

func (u *zaehlUhr) WarteBis(ctx context.Context, zeitpunkt int64) error {
	u.mu.Lock()
	u.aufrufe++
	if ctx.Err() != nil {
		u.nachEnde++
	}
	u.mu.Unlock()
	return u.Uhr.WarteBis(ctx, zeitpunkt)
}

func (u *zaehlUhr) zahlen() (lesungen, aufrufe, nachEnde int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.lesungen, u.aufrufe, u.nachEnde
}

// zeitSession ist eine Session mit je einer einfachen Anfrage je Text aus sqls, deren
// offset_ms offsets sind; ein negativer Wert lässt das Feld weg.
func zeitSession(id int, sqls []string, offsets []int64) model.Session {
	s := anfragen(id, nil, sqls...)
	for i := range s.Interactions {
		if offsets[i] >= 0 {
			off := offsets[i]
			s.Interactions[i].OffsetMS = &off
		}
	}
	return s
}

// zeitRec ist eine Aufzeichnung aus sessions.
func zeitRec(sessions ...model.Session) model.Recording {
	rec := model.NewRecording()
	rec.Sessions = sessions
	return rec
}

// vorrueckend ist eine Antwort ReadyForQuery, bei der die Uhr um ms vorrückt: Die
// Antwort braucht Zeit.
func vorrueckend(uhr *testuhr.Uhr, ms int64) schritt {
	s := bereitZuEnde()
	s.vorher = func() { uhr.Vorruecken(ms) }
	return s
}

// zeitLauf ist ein laufendes Einspielen mit Uhr.
type zeitLauf struct {
	t      *testing.T
	uhr    *zaehlUhr
	ziel   *fakeZiel
	fertig chan error
}

// starteZeit startet Play mit KeepTiming nach modus und bezug in einer Goroutine.
func starteZeit(ctx context.Context, t *testing.T, ablauf <-chan struct{}, rec model.Recording, ziel *fakeZiel, uhr *zaehlUhr, o services.PlayOptions) *zeitLauf {
	t.Helper()
	s, err := services.NewPlayService(context.WithoutCancel(ctx), ladeRepo{rec: rec}, "r.yaml", ziel, o, services.Wartezeiten(uhr))
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	l := &zeitLauf{t: t, uhr: uhr, ziel: ziel, fertig: make(chan error, 1)}
	go func() { l.fertig <- s.Play(ctx, ablauf) }()
	return l
}

// absolut, relativ und ersteAnfrage sind die Optionen der drei Betriebsarten.
func relativ() services.PlayOptions {
	return services.PlayOptions{KeepTiming: true, TimingMode: services.TimingRelative}
}

func absolut() services.PlayOptions {
	return services.PlayOptions{KeepTiming: true, TimingMode: services.TimingAbsolute, TimingReference: services.ReferenceConnect}
}

func ersteAnfrage() services.PlayOptions {
	return services.PlayOptions{KeepTiming: true, TimingMode: services.TimingAbsolute, TimingReference: services.ReferenceFirstRequest}
}

// wartet wartet höchstens 10 s, bis n Aufrufe von WarteBis warten; endet Play davor,
// ist der Test rot.
func (l *zeitLauf) wartet(n int) {
	l.t.Helper()
	frist := testzeit.NeueFrist(10000)
	for l.uhr.Wartende() != n {
		select {
		case err := <-l.fertig:
			l.t.Fatalf("Play endete (%v, Ablauf %q), erwartet %d Wartende", err, l.ziel.liste(), n)
		default:
		}
		if frist.Abgelaufen() {
			l.t.Fatalf("binnen 10 s warten %d statt %d (Ablauf %q)", l.uhr.Wartende(), n, l.ziel.liste())
		}
		testzeit.Schlafe(1)
	}
}

// gesendet sind die Aufrufe, die Client-Nachrichten senden, in ihrer Reihenfolge.
func (z *fakeZiel) gesendet() []string {
	var out []string
	for _, e := range z.liste() {
		if strings.HasPrefix(e, "anfrage ") || strings.HasPrefix(e, "gruppe ") {
			out = append(out, e)
		}
	}
	return out
}

// bisGesendet wartet höchstens 10 s, bis genau want gesendet ist.
func (l *zeitLauf) bisGesendet(want ...string) {
	l.t.Helper()
	frist := testzeit.NeueFrist(10000)
	for !reflect.DeepEqual(l.ziel.gesendet(), want) {
		if frist.Abgelaufen() {
			l.t.Fatalf("binnen 10 s gesendet %q, erwartet %q", l.ziel.gesendet(), want)
		}
		testzeit.Schlafe(1)
	}
}

// wartetNoch prüft, dass nach einem Vorrücken noch n warten und nur want gesendet ist.
func (l *zeitLauf) wartetNoch(n int, want ...string) {
	l.t.Helper()
	if got := l.uhr.Wartende(); got != n {
		l.t.Fatalf("%d Wartende, erwartet %d (Ablauf %q)", got, n, l.ziel.liste())
	}
	if got := l.ziel.gesendet(); !reflect.DeepEqual(got, want) {
		l.t.Fatalf("gesendet %q, erwartet %q", got, want)
	}
}

// endet wartet höchstens 10 s auf das Ende von Play und liefert dessen Ergebnis.
func (l *zeitLauf) endet() error {
	l.t.Helper()
	select {
	case err := <-l.fertig:
		return err
	case <-testzeit.Nach(10000):
		l.t.Fatalf("Play endet binnen 10 s nicht (Ablauf %q)", l.ziel.liste())
		return nil
	}
}

// Abdeckung: LH-FA-21/Happy, LH-FA-21/Boundary — relative: Vor jeder Interaktion außer der
// ersten einer Session wartet play, bis seit dem Beginn der vorigen mindestens der
// Unterschied ihrer beiden offset_ms vergangen ist (hier 1000 und 500 ms, auch vor einer
// Extended-Interaktion); der Beginn ist das Senden ihrer ersten Client-Nachricht, nicht das
// Ende ihrer Antwort (die Antworten brauchen 7 ms); die Wartezeit endet genau beim
// Zeitpunkt, 1 ms davor wartet play noch (LH-FA-21.a *Einspielen*, *Beginn einer
// Interaktion*).
func TestPlayZeitRelative(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	sitzung := &fakeEinspiel{antworten: [][]schritt{{vorrueckend(uhr.Uhr, 7)}, {vorrueckend(uhr.Uhr, 7)}}}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{sitzung}}
	s := zeitSession(1, []string{"A", "B"}, []int64{0, 1000})
	extended := syncInteraktion(3)
	off := int64(1500)
	extended.OffsetMS = &off
	s.Interactions = append(s.Interactions, extended)
	l := starteZeit(context.Background(), t, nil, zeitRec(s), ziel, uhr, relativ())

	l.wartet(1)
	l.wartetNoch(1, "anfrage A")
	uhr.Vorruecken(992) // 7 + 992 = 999
	l.wartetNoch(1, "anfrage A")
	uhr.Vorruecken(1)
	l.bisGesendet("anfrage A", "anfrage B")
	l.wartet(1) // B begann bei 1000, die Antwort liegt bei 1007; C wartet bis 1500
	uhr.Vorruecken(492)
	l.wartetNoch(1, "anfrage A", "anfrage B")
	uhr.Vorruecken(1)
	if err := l.endet(); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if got, want := ziel.gesendet(), []string{"anfrage A", "anfrage B", "gruppe sync"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("gesendet %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-21/Boundary — relative: Die Verspätung wird nicht aufgeholt; kam die
// zweite Interaktion 40 ms nach ihrem Zeitpunkt, wartet play vor der dritten die volle
// aufgezeichnete Pause ab ihrem tatsächlichen Beginn; vor der ersten Interaktion einer
// Session wartet relative nie, gleich wie groß ihr offset_ms ist, und gleiche Werte
// hintereinander warten nicht (LH-FA-21.a *Einspielen*, *Wartezeit und
// Verbindungsaufbau*).
func TestPlayZeitRelativeVerspaetung(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	ziel := &fakeZiel{}
	rec := zeitRec(zeitSession(1, []string{"A", "B", "C", "D"}, []int64{5000, 5100, 5150, 5150}))
	l := starteZeit(context.Background(), t, nil, rec, ziel, uhr, relativ())

	l.wartet(1)
	l.wartetNoch(1, "anfrage A") // A wartete nicht
	uhr.Vorruecken(140)          // B wäre bei 100 dran
	l.bisGesendet("anfrage A", "anfrage B")
	l.wartet(1) // C: 140 + 50 = 190
	uhr.Vorruecken(49)
	l.wartetNoch(1, "anfrage A", "anfrage B")
	uhr.Vorruecken(1)
	if err := l.endet(); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if got := ziel.gesendet(); len(got) != 4 {
		t.Fatalf("gesendet %q", got)
	}
	if _, aufrufe, _ := uhr.zahlen(); aufrufe != 3 {
		t.Fatalf("%d Aufrufe von WarteBis, erwartet 3 (A wartet nie)", aufrufe)
	}
}

// Abdeckung: LH-FA-21/Happy, LH-FA-21/Boundary — absolute mit connect: Jede Interaktion,
// auch die erste, wartet bis zum Bezugspunkt plus ihrem offset_ms; der Bezugspunkt ist die
// Lesung am Beginn des Aufbaus (der Aufbau dauert hier 50 ms); eine Verspätung wird
// aufgeholt: kam die zweite 40 ms zu spät, folgt die dritte 10 ms danach statt der
// aufgezeichneten 50 (LH-FA-21.a *Einspielen*, *Wartezeit und Verbindungsaufbau*).
func TestPlayZeitAbsoluteConnect(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	ziel := &fakeZiel{beimVerbinden: func() { uhr.Vorruecken(50) }}
	rec := zeitRec(zeitSession(1, []string{"A", "B", "C"}, []int64{100, 300, 350}))
	l := starteZeit(context.Background(), t, nil, rec, ziel, uhr, absolut())

	l.wartet(1)
	l.wartetNoch(1)
	uhr.Vorruecken(49) // 99
	l.wartetNoch(1)
	uhr.Vorruecken(1) // A bei 100
	l.bisGesendet("anfrage A")
	l.wartet(1)
	uhr.Vorruecken(240) // 340: B ist 40 ms zu spät
	l.bisGesendet("anfrage A", "anfrage B")
	l.wartet(1)
	uhr.Vorruecken(9)
	l.wartetNoch(1, "anfrage A", "anfrage B")
	uhr.Vorruecken(1)
	if err := l.endet(); err != nil {
		t.Fatalf("Play: %v", err)
	}
}

// Abdeckung: LH-FA-21/Happy, LH-FA-21/Boundary — absolute mit first-request: Die erste
// Interaktion der Session wartet nie und ist der Bezugspunkt, gemessen bei ihrem Beginn
// (hier 50 ms nach der Lesung am Beginn des Aufbaus); jede weitere wartet bis zu diesem
// Beginn plus dem Unterschied ihres offset_ms zu dem der ersten (LH-FA-21.a *Einspielen*,
// *Wartezeit und Verbindungsaufbau*).
func TestPlayZeitAbsoluteErsteAnfrage(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	ziel := &fakeZiel{beimVerbinden: func() { uhr.Vorruecken(50) }}
	rec := zeitRec(zeitSession(1, []string{"A", "B", "C"}, []int64{100, 300, 350}))
	l := starteZeit(context.Background(), t, nil, rec, ziel, uhr, ersteAnfrage())

	l.bisGesendet("anfrage A") // ohne Vorruecken
	l.wartet(1)
	uhr.Vorruecken(199) // 249; B bei 50 + 200
	l.wartetNoch(1, "anfrage A")
	uhr.Vorruecken(1)
	l.bisGesendet("anfrage A", "anfrage B")
	l.wartet(1) // C bei 250 + 50 = 300
	uhr.Vorruecken(49)
	l.wartetNoch(1, "anfrage A", "anfrage B")
	uhr.Vorruecken(1)
	if err := l.endet(); err != nil {
		t.Fatalf("Play: %v", err)
	}
}

// Abdeckung: LH-FA-21/Boundary — Jede Session beginnt mit eigener Uhr: Bei absolute mit
// connect wird der Bezugspunkt am Beginn jedes Aufbaus gelesen, die zweite Session wartet
// ab dem Ende der ersten (LH-FA-21.a *Einspielen*, *Wartezeit und Verbindungsaufbau*).
func TestPlayZeitSessionAbsolute(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	ziel := &fakeZiel{}
	rec := zeitRec(zeitSession(1, []string{"A", "B"}, []int64{0, 100}), zeitSession(2, []string{"C"}, []int64{100}))
	l := starteZeit(context.Background(), t, nil, rec, ziel, uhr, absolut())
	l.wartet(1)
	uhr.Vorruecken(100)
	l.bisGesendet("anfrage A", "anfrage B")
	l.wartet(1) // C bei 100 + 100
	uhr.Vorruecken(99)
	l.wartetNoch(1, "anfrage A", "anfrage B")
	uhr.Vorruecken(1)
	if err := l.endet(); err != nil {
		t.Fatalf("Play: %v", err)
	}
}

// Abdeckung: LH-FA-21/Boundary — bei relative wartet die erste Interaktion der zweiten
// Session nicht auf die letzte der ersten, gleich wie groß ihr offset_ms ist; die zweite
// wartet ab ihrer ersten (LH-FA-21.a *Wartezeit und Verbindungsaufbau*).
func TestPlayZeitSessionRelative(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	ziel := &fakeZiel{}
	rec := zeitRec(zeitSession(1, []string{"A", "B"}, []int64{0, 100}), zeitSession(2, []string{"C", "D"}, []int64{5000, 5100}))
	l := starteZeit(context.Background(), t, nil, rec, ziel, uhr, relativ())
	l.wartet(1)
	uhr.Vorruecken(100)
	l.bisGesendet("anfrage A", "anfrage B", "anfrage C")
	l.wartet(1) // D bei 100 + 100
	uhr.Vorruecken(99)
	l.wartetNoch(1, "anfrage A", "anfrage B", "anfrage C")
	uhr.Vorruecken(1)
	if err := l.endet(); err != nil {
		t.Fatalf("Play: %v", err)
	}
}

// Abdeckung: LH-FA-21/Boundary — nach einer Fehlerantwort, mit der --continue-on-error
// weiterläuft, bleibt der Bezug von relative der Beginn der vorigen aufgezeichneten
// Interaktion, nicht das Ende ihrer Antwort und nicht der Beginn der Interaktion davor;
// hier begann B 50 ms nach ihrem Zeitpunkt und antwortete 30 ms später mit dem Fehler
// (LH-FA-21.a *Beginn einer Interaktion*).
func TestPlayZeitNachFehlerantwort(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	fehlerNachZeit := schritt{r: fehlerantwort(map[string]string{"S": "ERROR", "V": "ERROR", "C": "42P01", "M": "x"}).r, vorher: func() { uhr.Vorruecken(30) }}
	sitzung := &fakeEinspiel{antworten: [][]schritt{nil, {fehlerNachZeit, bereitZuEnde()}}}
	ziel := &fakeZiel{beimVerbinden: func() { uhr.Vorruecken(10) }, sessions: []*fakeEinspiel{sitzung}}
	rec := zeitRec(zeitSession(1, []string{"A", "B", "C"}, []int64{0, 100, 150}))
	o := relativ()
	o.ContinueOnError = true
	l := starteZeit(context.Background(), t, nil, rec, ziel, uhr, o)

	l.wartet(1) // A bei 10, B bei 110
	uhr.Vorruecken(150)
	l.bisGesendet("anfrage A", "anfrage B")
	l.wartet(1) // B bei 160, Antwort bei 190, C bei 160 + 50 = 210
	uhr.Vorruecken(19)
	l.wartetNoch(1, "anfrage A", "anfrage B")
	uhr.Vorruecken(1)
	err := l.endet()
	if want := []string{e4004("Session 1, Interaktion 2", "42P01", "x")}; !reflect.DeepEqual(texte(err), want) {
		t.Fatalf("Meldungen %q, erwartet %q", texte(err), want)
	}
	if got := ziel.gesendet(); !reflect.DeepEqual(got, []string{"anfrage A", "anfrage B", "anfrage C"}) {
		t.Fatalf("gesendet %q", got)
	}
}

// Abdeckung: LH-FA-21/Negative — bricht das Einspielen mit einem Fehler ab, gibt es keine
// weitere Wartezeit: nach der Fehlerantwort ohne --continue-on-error ruft play WarteBis
// nicht mehr auf (LH-FA-21.a *Abbruch*).
func TestPlayZeitKeineWartezeitNachAbbruch(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	sitzung := &fakeEinspiel{antworten: [][]schritt{{fehlerantwort(map[string]string{"S": "ERROR", "V": "ERROR", "C": "42P01", "M": "x"})}}}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{sitzung}}
	rec := zeitRec(zeitSession(1, []string{"A", "B"}, []int64{0, 100}))
	l := starteZeit(context.Background(), t, nil, rec, ziel, uhr, relativ())
	err := l.endet()
	if codeVon(err) != model.CodeServerError {
		t.Fatalf("Fehler %v", err)
	}
	if _, aufrufe, _ := uhr.zahlen(); aufrufe != 0 || uhr.Wartende() != 0 {
		t.Fatalf("%d Aufrufe von WarteBis, %d Wartende", aufrufe, uhr.Wartende())
	}
}

// Abdeckung: LH-FA-21/Negative — scheitert der Aufbau, gibt es keine Wartezeit, und der
// Fehler gilt wie ohne Option: Die Aufbau-Fehler PGR-E4002 mit ihrer Session kommen ohne
// Aufruf von WarteBis (LH-FA-21.a *Wartezeit und Verbindungsaufbau*).
func TestPlayZeitAufbauScheitert(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	ziel := &fakeZiel{verbindeFehler: model.Errorf(model.CodeUpstream, nil, "weg")}
	rec := zeitRec(zeitSession(1, []string{"A"}, []int64{500}))
	l := starteZeit(context.Background(), t, nil, rec, ziel, uhr, absolut())
	err := l.endet()
	if want := []string{"Netzwerk [PGR-E4002]: Session 1: weg"}; !reflect.DeepEqual(texte(err), want) {
		t.Fatalf("Meldungen %q, erwartet %q", texte(err), want)
	}
	if _, aufrufe, _ := uhr.zahlen(); aufrufe != 0 {
		t.Fatalf("%d Aufrufe von WarteBis", aufrufe)
	}
}

// Abdeckung: LH-FA-21/Boundary — große Werte haben keine Obergrenze, und der Zeitpunkt
// (Bezug plus Wert) wird bei einem Überlauf von int64 auf den größten Wert begrenzt, nie
// negativ: In jeder Betriebsart wartet play bei offset_ms math.MaxInt64 bis zum
// Abbruchsignal und sendet die Interaktion nicht (LH-FA-21.a *Randformen des Einspielens*,
// *Großer Wert*).
func TestPlayZeitGrosserWert(t *testing.T) {
	for _, f := range []struct {
		name    string
		o       services.PlayOptions
		offsets []int64
		sendet  []string
	}{
		{"relative", relativ(), []int64{0, math.MaxInt64}, []string{"anfrage A"}},
		{"absolute, connect", absolut(), []int64{math.MaxInt64}, nil},
		{"absolute, first-request", ersteAnfrage(), []int64{0, math.MaxInt64}, []string{"anfrage A"}},
	} {
		t.Run(f.name, func(t *testing.T) {
			uhr := &zaehlUhr{Uhr: testuhr.Neu()}
			uhr.Vorruecken(50)
			ziel := &fakeZiel{}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			rec := zeitRec(zeitSession(1, []string{"A", "B"}[:len(f.offsets)], f.offsets))
			l := starteZeit(ctx, t, nil, rec, ziel, uhr, f.o)
			l.wartet(1)
			uhr.Vorruecken(1_000_000)
			l.wartetNoch(1, f.sendet...)
			cancel()
			if err := l.endet(); err != nil {
				t.Fatalf("Play: %v", err)
			}
			if got := ziel.gesendet(); !reflect.DeepEqual(got, f.sendet) {
				t.Fatalf("gesendet %q, erwartet %q", got, f.sendet)
			}
		})
	}
}

// Abdeckung: LH-FA-21/Boundary — das erste Abbruchsignal beendet jede Wartezeit sofort;
// ohne --finish-session-on-interrupt wird die folgende Interaktion nicht gesendet, die
// Session schließt, ohne Fehler (LH-FA-21.a *Abbruch*, LH-FA-20.a Schritt 7).
func TestPlayZeitAbbruchWartezeit(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	ziel := &fakeZiel{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rec := zeitRec(zeitSession(1, []string{"A", "B", "C"}, []int64{0, 1000, 2000}))
	l := starteZeit(ctx, t, nil, rec, ziel, uhr, relativ())
	l.wartet(1)
	cancel()
	if err := l.endet(); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if got, want := ziel.liste(), []string{"verbinde", "anfrage A", "naechste", "schliesse"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
	if _, aufrufe, nachEnde := uhr.zahlen(); aufrufe != 1 || nachEnde != 0 || uhr.Wartende() != 0 {
		t.Fatalf("%d Aufrufe von WarteBis, davon %d nach dem Signal, %d Wartende", aufrufe, nachEnde, uhr.Wartende())
	}
}

// Abdeckung: LH-FA-21/Boundary — mit --finish-session-on-interrupt beendet das erste
// Signal die Wartezeit ebenso sofort, und play spielt die übrigen Interaktionen der Session
// ohne Wartezeit zu Ende: die unterbrochene eingeschlossen, ohne weiteren Aufruf von
// WarteBis; das gilt auch, wenn das Signal im Aufbau eintrifft, bei absolute mit
// Wartezeiten, die weit in der Zukunft liegen (LH-FA-21.a *Abbruch*).
func TestPlayZeitAbbruchMitFinish(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	ziel := &fakeZiel{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	o := relativ()
	o.FinishSessionOnInterrupt = true
	rec := zeitRec(zeitSession(1, []string{"A", "B", "C"}, []int64{0, 1000, 2000}))
	l := starteZeit(ctx, t, nil, rec, ziel, uhr, o)
	l.wartet(1)
	cancel()
	if err := l.endet(); err != nil {
		t.Fatalf("Play: %v", err)
	}
	want := []string{"verbinde", "anfrage A", "naechste", "anfrage B", "naechste", "anfrage C", "naechste", "schliesse"}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
	if _, aufrufe, nachEnde := uhr.zahlen(); aufrufe != 1 || nachEnde != 0 {
		t.Fatalf("%d Aufrufe von WarteBis, davon %d nach dem Signal, erwartet 1 und 0", aufrufe, nachEnde)
	}

	uhr = &zaehlUhr{Uhr: testuhr.Neu()}
	ctx, aus := context.WithCancel(context.Background())
	defer aus()
	ziel = &fakeZiel{beimVerbinden: aus}
	o = absolut()
	o.FinishSessionOnInterrupt = true
	rec = zeitRec(zeitSession(1, []string{"A", "B"}, []int64{100, math.MaxInt64}))
	l = starteZeit(ctx, t, nil, rec, ziel, uhr, o)
	if err := l.endet(); err != nil {
		t.Fatalf("Signal im Aufbau: %v", err)
	}
	if got := ziel.gesendet(); !reflect.DeepEqual(got, []string{"anfrage A", "anfrage B"}) {
		t.Fatalf("Signal im Aufbau: gesendet %q", got)
	}
	if _, aufrufe, _ := uhr.zahlen(); aufrufe != 0 {
		t.Fatalf("Signal im Aufbau: %d Aufrufe von WarteBis", aufrufe)
	}
}

// Abdeckung: LH-FA-21/Boundary — das zweite Signal (ablauf geschlossen) beendet eine
// Wartezeit sofort; die Verbindung ist geschlossen, wenn play zurückkehrt (das Schließen
// dauert hier 100 ms), die folgende Interaktion nicht gesendet, ohne Fehler (LH-FA-21.a
// *Abbruch*, LH-FA-20.a *Abbruchsignal*).
func TestPlayZeitZweitesSignalInWartezeit(t *testing.T) {
	uhr := &zaehlUhr{Uhr: testuhr.Neu()}
	langsam := &fakeEinspiel{beiSchliesse: func() { testzeit.Schlafe(100) }}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{langsam}}
	ablauf := make(chan struct{})
	rec := zeitRec(zeitSession(1, []string{"A", "B"}, []int64{0, 1000}))
	l := starteZeit(context.Background(), t, ablauf, rec, ziel, uhr, relativ())
	l.wartet(1)
	close(ablauf)
	if err := l.endet(); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if got, want := ziel.liste(), []string{"verbinde", "anfrage A", "naechste", "schliesse"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-21/Happy, LH-FA-21/Negative — ohne --keep-timing wartet play nie, liest
// die Uhr nie und braucht keine, auch wenn die Aufzeichnung offset_ms trägt oder an
// einigen Interaktionen fehlt (LH-FA-21.a *Einspielen*, *Randformen des Einspielens*).
func TestPlayZeitOhneKeepTiming(t *testing.T) {
	rec := zeitRec(zeitSession(1, []string{"A", "B", "C"}, []int64{0, 5000, -1}))
	for _, mitUhr := range []bool{false, true} {
		uhr := &zaehlUhr{Uhr: testuhr.Neu()}
		ziel := &fakeZiel{}
		o := services.PlayOptions{TimingMode: services.TimingAbsolute, TimingReference: services.ReferenceConnect}
		var opts []services.PlayOption
		if mitUhr {
			opts = append(opts, services.Wartezeiten(uhr))
		}
		s, err := services.NewPlayService(context.Background(), ladeRepo{rec: rec}, "r.yaml", ziel, o, opts...)
		if err != nil {
			t.Fatalf("Start: %v", err)
		}
		fertig := make(chan error, 1)
		go func() { fertig <- s.Play(context.Background(), nil) }()
		select {
		case err := <-fertig:
			if err != nil {
				t.Fatalf("Play: %v", err)
			}
		case <-testzeit.Nach(10000):
			t.Fatal("Play endet binnen 10 s nicht: wartet es ohne --keep-timing?")
		}
		if got := ziel.gesendet(); !reflect.DeepEqual(got, []string{"anfrage A", "anfrage B", "anfrage C"}) {
			t.Fatalf("gesendet %q", got)
		}
		if lesungen, aufrufe, _ := uhr.zahlen(); lesungen != 0 || aufrufe != 0 {
			t.Fatalf("%d Lesungen, %d Aufrufe von WarteBis, erwartet keine", lesungen, aufrufe)
		}
	}
}

// Abdeckung: LH-FA-21/Negative — fehlt mit --keep-timing irgendeiner Interaktion
// offset_ms, ist das PGR-E2003 (Exit-Code 2) im Start: Die Meldung nennt Session und
// Interaktion der ersten ohne Wert in der Reihenfolge der Sessions und Interaktionen (auch
// einer späteren Interaktion, einer späteren Session und hinter einer Session ohne
// Interaktion), nie einen Wert; ohne --keep-timing ist dieselbe Aufzeichnung gültig
// (LH-FA-21.a *Randformen des Einspielens*, *Prüfung der Werte*).
func TestPlayZeitFehltOffset(t *testing.T) {
	for _, f := range []struct {
		name string
		rec  model.Recording
		ort  string
	}{
		{"erste Interaktion", zeitRec(zeitSession(1, []string{"A", "B"}, []int64{-1, 4242})), "Session 1, Interaktion 1"},
		{"spätere Interaktion", zeitRec(zeitSession(1, []string{"A", "B", "C"}, []int64{0, 4242, -1})), "Session 1, Interaktion 3"},
		{"spätere Session", zeitRec(zeitSession(1, []string{"A"}, []int64{4242}), model.Session{ID: 2}, zeitSession(3, []string{"B", "C"}, []int64{4243, -1})), "Session 3, Interaktion 2"},
		{"zwei fehlen", zeitRec(zeitSession(1, []string{"A", "B"}, []int64{0, -1}), zeitSession(2, []string{"C"}, []int64{-1})), "Session 1, Interaktion 2"},
	} {
		t.Run(f.name, func(t *testing.T) {
			ziel := &fakeZiel{}
			uhr := testuhr.Neu()
			_, err := services.NewPlayService(context.Background(), ladeRepo{rec: f.rec}, "r.yaml", ziel, relativ(), services.Wartezeiten(uhr))
			ms := texte(err)
			if len(ms) != 1 || !strings.HasPrefix(ms[0], "Konfiguration [PGR-E2003]: "+f.ort+" ") || exitVon(err) != 2 || strings.Contains(ms[0], "4242") || strings.Contains(ms[0], "4243") {
				t.Fatalf("Meldungen %q, Exit-Code %d, erwartet PGR-E2003 mit %q", ms, exitVon(err), f.ort)
			}
			if got := ziel.liste(); len(got) != 0 {
				t.Fatalf("Start verbindet: %q", got)
			}
			if _, err := services.NewPlayService(context.Background(), ladeRepo{rec: f.rec}, "r.yaml", ziel, services.PlayOptions{}); err != nil {
				t.Fatalf("ohne --keep-timing: %v", err)
			}
		})
	}
}

// Abdeckung: LH-FA-21/Boundary — Sessions ohne Interaktion verlangen kein offset_ms, und
// ein Fehler beim Laden (PGR-E3003 einer beschädigten Interaktion) geht PGR-E2003 vor
// (LH-FA-21.a *Randformen des Einspielens*, *Prüfung der Werte*).
func TestPlayZeitFehltOffsetRand(t *testing.T) {
	ziel := &fakeZiel{}
	rec := zeitRec(model.Session{ID: 1}, zeitSession(2, []string{"A"}, []int64{0}), model.Session{ID: 3})
	if _, err := services.NewPlayService(context.Background(), ladeRepo{rec: rec}, "r.yaml", ziel, relativ(), services.Wartezeiten(testuhr.Neu())); err != nil {
		t.Fatalf("Sessions ohne Interaktion: %v", err)
	}
	kaputt := zeitRec(zeitSession(1, []string{"A"}, []int64{-1}))
	kaputt.Sessions[0].Interactions = append(kaputt.Sessions[0].Interactions, model.Interaction{Sequence: 2, Request: model.Request{Type: model.RequestQuery, SQL: "B"}})
	if _, err := services.NewPlayService(context.Background(), ladeRepo{rec: kaputt}, "r.yaml", ziel, relativ(), services.Wartezeiten(testuhr.Neu())); codeVon(err) != model.CodeRecordingBroken {
		t.Fatalf("beschädigte Interaktion neben fehlendem offset_ms: %v", err)
	}
}

var _ driven.Uhr = (*zaehlUhr)(nil)
