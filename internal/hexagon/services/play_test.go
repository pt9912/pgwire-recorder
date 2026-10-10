package services_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// schritt ist eine Antwort von Naechste: eine Nachricht oder ein Fehler.
// blockiert wartet vorher, bis Schliesse läuft, und liefert dann einen Fehler
// PGR-E4003. vorher läuft, wenn nicht nil, bevor Naechste den Schritt liefert.
type schritt struct {
	r         model.Response
	err       error
	blockiert bool
	vorher    func()
}

// fakeZiel ist das Einspielziel der Tests. Je Verbinde nimmt es die nächste
// Session aus sessions (fehlt sie, eine ohne Antworten); verbindeFehler liefert
// je Aufruf einen Fehler statt einer Session. beimVerbinden läuft während
// Verbinde; blockiert Verbinde wartet bis ctx endet. ereignisse hält die
// Aufrufe in ihrer Reihenfolge.
type fakeZiel struct {
	mu                sync.Mutex
	ereignisse        []string
	startups          []map[string]string
	sessions          []*fakeEinspiel
	verbindeFehler    error
	beimVerbinden     func()
	verbindeBlockiert bool
	verbindeCtxEnde   bool
}

func (z *fakeZiel) notiere(e string) {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.ereignisse = append(z.ereignisse, e)
}

func (z *fakeZiel) liste() []string {
	z.mu.Lock()
	defer z.mu.Unlock()
	return append([]string(nil), z.ereignisse...)
}

func (z *fakeZiel) Verbinde(ctx context.Context, startup map[string]string) (driven.EinspielSession, error) {
	z.notiere("verbinde")
	z.mu.Lock()
	z.startups = append(z.startups, startup)
	z.mu.Unlock()
	if z.beimVerbinden != nil {
		z.beimVerbinden()
	}
	if z.verbindeBlockiert {
		<-ctx.Done()
		return nil, model.Errorf(model.CodeUpstream, ctx.Err(), "abgebrochen")
	}
	z.mu.Lock()
	z.verbindeCtxEnde = ctx.Err() != nil
	z.mu.Unlock()
	if z.verbindeFehler != nil {
		return nil, z.verbindeFehler
	}
	z.mu.Lock()
	defer z.mu.Unlock()
	s := &fakeEinspiel{ziel: z, zu: make(chan struct{})}
	if len(z.sessions) > 0 {
		s = z.sessions[0]
		s.ziel, s.zu = z, make(chan struct{})
		z.sessions = z.sessions[1:]
	}
	return s, nil
}

// fakeEinspiel ist eine Session des fakeZiel: je Anfrage die nächste Folge
// aus antworten; beiAnfrage läuft mit dem Text der Anfrage, anfrageFehler
// liefert Anfrage als Fehler. Je Gruppe hängt Gruppe die nächste Folge aus
// gruppen an die noch ungelesenen Antworten an; beiGruppe läuft mit der Nummer
// der Gruppe seit dem Aufbau (ab 1), gruppeFehler liefert Gruppe als Fehler.
type fakeEinspiel struct {
	ziel          *fakeZiel
	antworten     [][]schritt
	gruppen       [][]schritt
	folge         []schritt
	beiAnfrage    func(string)
	anfrageFehler error
	beiGruppe     func(int)
	gruppeFehler  error
	nGruppen      int
	zu            chan struct{}
	einmal        sync.Once
}

func (s *fakeEinspiel) Gruppe(nachrichten []model.ClientMessage) error {
	var typen []string
	for _, m := range nachrichten {
		typen = append(typen, string(m.Type))
	}
	s.ziel.notiere("gruppe " + strings.Join(typen, ","))
	s.nGruppen++
	if s.beiGruppe != nil {
		s.beiGruppe(s.nGruppen)
	}
	if s.gruppeFehler != nil {
		return s.gruppeFehler
	}
	if len(s.gruppen) > 0 {
		s.folge = append(s.folge, s.gruppen[0]...)
		s.gruppen = s.gruppen[1:]
	}
	return nil
}

func (s *fakeEinspiel) Anfrage(sql string) error {
	s.ziel.notiere("anfrage " + sql)
	if s.beiAnfrage != nil {
		s.beiAnfrage(sql)
	}
	if s.anfrageFehler != nil {
		return s.anfrageFehler
	}
	if len(s.antworten) > 0 {
		s.folge, s.antworten = s.antworten[0], s.antworten[1:]
	}
	return nil
}

func (s *fakeEinspiel) Naechste() (model.Response, error) {
	s.ziel.notiere("naechste")
	if len(s.folge) == 0 {
		return model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"}, nil
	}
	x := s.folge[0]
	s.folge = s.folge[1:]
	if x.vorher != nil {
		x.vorher()
	}
	if x.blockiert {
		<-s.zu
		return model.Response{}, model.Errorf(model.CodeConnectionLost, nil, "geschlossen")
	}
	return x.r, x.err
}

func (s *fakeEinspiel) Schliesse() {
	s.ziel.notiere("schliesse")
	s.einmal.Do(func() { close(s.zu) })
}

func bereitZuEnde() schritt {
	return schritt{r: model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"}}
}

func fehlerantwort(felder map[string]string) schritt {
	return schritt{r: model.Response{Type: model.ResponseErrorResponse, Fields: felder}}
}

// anfragen ist eine Session mit je einer einfachen Anfrage je Text.
func anfragen(id int, startup map[string]string, sqls ...string) model.Session {
	s := model.Session{ID: id, Startup: startup}
	for i, sql := range sqls {
		s.Interactions = append(s.Interactions, interaktion(i+1, sql, "SELECT 1"))
	}
	return s
}

// spiele startet den PlayService mit rec und ziel und spielt ein.
func spiele(ctx context.Context, t *testing.T, ablauf <-chan struct{}, rec model.Recording, ziel *fakeZiel, benutzer, datenbank string) error {
	t.Helper()
	s, err := services.NewPlayService(context.WithoutCancel(ctx), ladeRepo{rec: rec}, "r.yaml", ziel, services.PlayOptions{User: benutzer, Database: datenbank})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	return s.Play(ctx, ablauf)
}

func codeVon(err error) string {
	var me *model.Error
	if errors.As(err, &me) {
		return me.Code
	}
	return ""
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary — play spielt jede Session mit
// Interaktionen über eine eigene Verbindung ein, die Sessions und ihre
// Anfragen in der Reihenfolge der Aufzeichnung, liest je Anfrage bis zum
// ReadyForQuery, verwirft die übrigen Antworten und schließt jede Verbindung
// nach ihrer letzten Interaktion; eine Session ohne Interaktion bekommt keine
// Verbindung; eine Session nur aus Lebendprüfungen wird eingespielt
// (LH-FA-20.a, LH-FA-12.a).
func TestPlayReihenfolge(t *testing.T) {
	rec := model.NewRecording()
	rec.Sessions = []model.Session{anfragen(1, nil, "A", "B"), {ID: 2}, anfragen(3, nil, "-- lebt", "")}
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{{
		{r: model.Response{Type: model.ResponseNoticeResponse}},
		{r: model.Response{Type: model.ResponseRowDescription}},
		{r: model.Response{Type: model.ResponseDataRow}},
		{r: model.Response{Type: model.ResponseParameterStatus}},
		bereitZuEnde(),
	}}}}}
	if err := spiele(context.Background(), t, nil, rec, ziel, "", ""); err != nil {
		t.Fatalf("Play: %v", err)
	}
	want := []string{
		"verbinde", "anfrage A", "naechste", "naechste", "naechste", "naechste", "naechste", "anfrage B", "naechste", "schliesse",
		"verbinde", "anfrage -- lebt", "naechste", "anfrage ", "naechste", "schliesse",
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Boundary — play sendet die aufgezeichneten
// Startup-Parameter unverändert; nicht leere benutzer und datenbank ersetzen
// user und database, leere lassen sie, wie aufgezeichnet, ein fehlender
// Parameter bleibt fehlend (LH-FA-20.a *Startup-Daten*).
func TestPlayStartup(t *testing.T) {
	aufgezeichnet := map[string]string{"user": "alt", "database": "altdb", "application_name": "a", "replication": "false"}
	for _, f := range []struct {
		benutzer, datenbank string
		start, want         map[string]string
	}{
		{"", "", aufgezeichnet, aufgezeichnet},
		{"neu", "", aufgezeichnet, map[string]string{"user": "neu", "database": "altdb", "application_name": "a", "replication": "false"}},
		{"", "neudb", aufgezeichnet, map[string]string{"user": "alt", "database": "neudb", "application_name": "a", "replication": "false"}},
		{"neu", "neudb", map[string]string{"x": "y"}, map[string]string{"user": "neu", "database": "neudb", "x": "y"}},
		{"", "", map[string]string{"x": "y"}, map[string]string{"x": "y"}},
		{"", "", nil, map[string]string{}},
	} {
		rec := model.NewRecording()
		rec.Sessions = []model.Session{anfragen(1, f.start, "A")}
		ziel := &fakeZiel{}
		if err := spiele(context.Background(), t, nil, rec, ziel, f.benutzer, f.datenbank); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ziel.startups[0], f.want) {
			t.Errorf("%q/%q mit %v: %v, erwartet %v", f.benutzer, f.datenbank, f.start, ziel.startups[0], f.want)
		}
	}
	if aufgezeichnet["user"] != "alt" || len(aufgezeichnet) != 4 {
		t.Fatalf("Aufzeichnung verändert: %v", aufgezeichnet)
	}
}

// Abdeckung: LH-FA-20/Negative — eine Fehlerantwort auf eine Anfrage bricht
// ab: PGR-E4004 mit Session, Interaktion, SQLSTATE und Meldung, keine weitere
// Antwort gelesen, keine weitere Anfrage und keine weitere Session, die
// Verbindung mit Schliesse beendet; mit dem Schweregrad FATAL oder PANIC aus
// V, ohne V aus S, ist sie PGR-E4003; die Meldung nennt von der Fehlerantwort
// SQLSTATE und Meldung, keine weiteren Felder, auch nicht den Schweregrad
// (LH-FA-20.a *Interaktion*, *Meldungen*).
func TestPlayFehlerantwort(t *testing.T) {
	for _, f := range []struct {
		name   string
		felder map[string]string
		code   string
	}{
		{"ERROR", map[string]string{"S": "ERROR", "V": "ERROR", "C": "42P01", "M": "fehlt"}, model.CodeServerError},
		{"ohne Schweregrad", map[string]string{"C": "42P01", "M": "fehlt"}, model.CodeServerError},
		{"V ERROR vor S FATAL", map[string]string{"S": "FATAL", "V": "ERROR", "C": "42P01", "M": "fehlt"}, model.CodeServerError},
		{"V FATAL", map[string]string{"S": "SCHWER", "V": "FATAL", "C": "57P01", "M": "fehlt"}, model.CodeConnectionLost},
		{"S FATAL ohne V", map[string]string{"S": "FATAL", "C": "57P01", "M": "fehlt"}, model.CodeConnectionLost},
		{"V PANIC", map[string]string{"S": "x", "V": "PANIC", "C": "XX000", "M": "fehlt"}, model.CodeConnectionLost},
		{"S PANIC ohne V", map[string]string{"S": "PANIC", "C": "XX000", "M": "fehlt"}, model.CodeConnectionLost},
	} {
		t.Run(f.name, func(t *testing.T) {
			rec := model.NewRecording()
			rec.Sessions = []model.Session{anfragen(1, nil, "A", "B"), anfragen(2, nil, "C")}
			ziel := &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{{{r: model.Response{Type: model.ResponseDataRow}}, fehlerantwort(f.felder), bereitZuEnde()}}}}}
			err := spiele(context.Background(), t, nil, rec, ziel, "", "")
			text := "[" + f.code + "]: Session 1, Interaktion 1: Fehlerantwort des Servers " + f.felder["C"] + " „fehlt“"
			if f.code == model.CodeConnectionLost {
				text += ", die Verbindung endet"
			}
			if codeVon(err) != f.code || !strings.HasSuffix(err.Error(), text) {
				t.Fatalf("%v, erwartet %s mit Ort, SQLSTATE und Meldung: %q", err, f.code, text)
			}
			want := []string{"verbinde", "anfrage A", "naechste", "naechste", "schliesse"}
			if got := ziel.liste(); !reflect.DeepEqual(got, want) {
				t.Fatalf("Ablauf %q, erwartet %q", got, want)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Negative — ein Fehler im Aufbau, beim Senden oder beim
// Lesen bricht ab und behält seinen Code; die Meldung nennt die Session, bei
// einem Fehler in einer Interaktion auch ihre Nummer; nach einem Fehler beim
// Senden liest play nicht weiter (LH-FA-20.a *Interaktion*, *Meldungen*).
func TestPlayFehler(t *testing.T) {
	rec := model.NewRecording()
	rec.Sessions = []model.Session{anfragen(4, nil, "A", "B"), anfragen(5, nil, "C")}
	ziel := &fakeZiel{verbindeFehler: model.Errorf(model.CodeLogin, nil, "Klasse 28")}
	err := spiele(context.Background(), t, nil, rec, ziel, "", "")
	if codeVon(err) != model.CodeLogin || !strings.HasPrefix(err.Error(), "Netzwerk [PGR-E4005]: Session 4: Klasse 28") {
		t.Errorf("Aufbau: %v", err)
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, []string{"verbinde"}) {
		t.Errorf("nach dem Fehler im Aufbau: %q", got)
	}

	ziel = &fakeZiel{sessions: []*fakeEinspiel{{anfrageFehler: model.Errorf(model.CodeConnectionLost, nil, "senden")}}}
	err = spiele(context.Background(), t, nil, rec, ziel, "", "")
	if codeVon(err) != model.CodeConnectionLost || !strings.Contains(err.Error(), "Session 4, Interaktion 1: senden") {
		t.Errorf("Senden: %v", err)
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, []string{"verbinde", "anfrage A", "schliesse"}) {
		t.Errorf("nach dem Fehler beim Senden: %q", got)
	}

	for _, code := range []string{model.CodeUnsupported, model.CodeConnectionLost} {
		ziel = &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{nil, {{err: model.Errorf(code, nil, "lesen")}}}}}}
		err = spiele(context.Background(), t, nil, rec, ziel, "", "")
		if codeVon(err) != code || !strings.Contains(err.Error(), "Session 4, Interaktion 2: lesen") {
			t.Errorf("Lesen %s: %v", code, err)
		}
		if got := ziel.liste(); !reflect.DeepEqual(got, []string{"verbinde", "anfrage A", "naechste", "anfrage B", "naechste", "schliesse"}) {
			t.Errorf("nach dem Fehler beim Lesen %s: %q", code, got)
		}
	}
}

// Abdeckung: LH-FA-20/Negative — der Start lädt die Aufzeichnung; ein
// Ladefehler und eine beschädigte Interaktion gehen vor; eine Extended-Interaktion,
// auch hinter einfachen Anfragen und in einer späteren Session, ist kein
// Startfehler; der Start verbindet sich dabei nie (LH-FA-20.a *Start*).
func TestPlayStart(t *testing.T) {
	rec := model.NewRecording()
	s2 := anfragen(2, nil, "A", "B")
	s2.Interactions = append(s2.Interactions, syncInteraktion(3))
	s3 := model.Session{ID: 3, Interactions: []model.Interaction{syncInteraktion(1)}}
	rec.Sessions = []model.Session{anfragen(1, nil, "X"), s2, s3}
	ziel := &fakeZiel{}
	if _, err := services.NewPlayService(context.Background(), ladeRepo{rec: rec}, "r.yaml", ziel, services.PlayOptions{}); err != nil {
		t.Errorf("Extended: %v", err)
	}
	kaputt := rec
	kaputt.Sessions = append([]model.Session{}, rec.Sessions...)
	kaputt.Sessions[2] = model.Session{ID: 3, Interactions: []model.Interaction{{Sequence: 1, Request: model.Request{Type: model.RequestQuery, SQL: "A"}}}}
	if _, err := services.NewPlayService(context.Background(), ladeRepo{rec: kaputt}, "r.yaml", ziel, services.PlayOptions{}); codeVon(err) != model.CodeRecordingBroken {
		t.Errorf("beschädigt vor Extended: %v", err)
	}
	if _, err := services.NewPlayService(context.Background(), ladeRepo{err: model.Errorf(model.CodeRecordingIO, nil, "weg")}, "r.yaml", ziel, services.PlayOptions{}); codeVon(err) != model.CodeRecordingIO {
		t.Errorf("Ladefehler: %v", err)
	}
	if got := ziel.liste(); len(got) != 0 {
		t.Errorf("Start verbindet: %q", got)
	}
}

// Abdeckung: LH-FA-20/Boundary — eine Aufzeichnung ohne Session mit
// Interaktion spielt play ohne Verbindung und ohne Fehler ein (LH-FA-20.a
// *Start*).
func TestPlayOhneInteraktion(t *testing.T) {
	for _, sessions := range [][]model.Session{nil, {{ID: 1}, {ID: 2}}} {
		rec := model.NewRecording()
		rec.Sessions = sessions
		ziel := &fakeZiel{}
		if err := spiele(context.Background(), t, nil, rec, ziel, "", ""); err != nil || len(ziel.liste()) != 0 {
			t.Errorf("%v: %v, %q", sessions, err, ziel.liste())
		}
	}
}

// zwei sind zwei Sessions mit je zwei Anfragen.
func zwei() model.Recording {
	rec := model.NewRecording()
	rec.Sessions = []model.Session{anfragen(1, nil, "A", "B"), anfragen(2, nil, "C", "D")}
	return rec
}

// Abdeckung: LH-FA-20/Boundary — das erste Abbruchsignal (ctx endet) beendet
// das Einspielen nach der laufenden Interaktion: sie liest bis zum
// ReadyForQuery, danach keine weitere Interaktion und keine weitere Session,
// die Verbindung wird geschlossen, ohne Fehler; endet ctx vor dem Einspielen,
// beginnt keine Session; ein Fehler der laufenden Interaktion zählt
// (LH-FA-20.a Schritt 7, *Abbruchsignal*).
func TestPlayErstesSignal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sitzung := &fakeEinspiel{antworten: [][]schritt{{{r: model.Response{Type: model.ResponseDataRow}}, bereitZuEnde()}}}
	sitzung.beiAnfrage = func(string) { cancel() }
	ziel := &fakeZiel{sessions: []*fakeEinspiel{sitzung}}
	if err := spiele(ctx, t, nil, zwei(), ziel, "", ""); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if got, want := ziel.liste(), []string{"verbinde", "anfrage A", "naechste", "naechste", "schliesse"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ablauf %q, erwartet %q", got, want)
	}

	vorher, aus := context.WithCancel(context.Background())
	aus()
	ziel = &fakeZiel{}
	if err := spiele(vorher, t, nil, zwei(), ziel, "", ""); err != nil || len(ziel.liste()) != 0 {
		t.Fatalf("Signal vor dem Einspielen: %v, %q", err, ziel.liste())
	}

	ctx, cancel = context.WithCancel(context.Background())
	sitzung = &fakeEinspiel{antworten: [][]schritt{{fehlerantwort(map[string]string{"C": "42P01", "M": "x"})}}}
	sitzung.beiAnfrage = func(string) { cancel() }
	ziel = &fakeZiel{sessions: []*fakeEinspiel{sitzung}}
	if err := spiele(ctx, t, nil, zwei(), ziel, "", ""); codeVon(err) != model.CodeServerError {
		t.Fatalf("Fehler nach dem ersten Signal: %v", err)
	}
}

// Abdeckung: LH-FA-20/Boundary — trifft das erste Signal im Aufbau ein, läuft
// er zu Ende, ohne dass sein ctx endet; danach keine Interaktion, die
// Verbindung wird geschlossen; ein Fehler im Aufbau zählt (LH-FA-20.a
// *Abbruchsignal*).
func TestPlayErstesSignalImAufbau(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ziel := &fakeZiel{beimVerbinden: cancel}
	if err := spiele(ctx, t, nil, zwei(), ziel, "", ""); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if got, want := ziel.liste(), []string{"verbinde", "schliesse"}; !reflect.DeepEqual(got, want) || ziel.verbindeCtxEnde {
		t.Fatalf("Ablauf %q, erwartet %q; ctx des Aufbaus beendet: %v", got, want, ziel.verbindeCtxEnde)
	}
	ctx, cancel = context.WithCancel(context.Background())
	ziel = &fakeZiel{beimVerbinden: cancel, verbindeFehler: model.Errorf(model.CodeUpstream, nil, "weg")}
	if err := spiele(ctx, t, nil, zwei(), ziel, "", ""); codeVon(err) != model.CodeUpstream {
		t.Fatalf("Fehler im Aufbau nach dem ersten Signal: %v", err)
	}
}

// bis prüft bedingung, bis sie gilt, höchstens 30 s lang; sonst endet der Test
// mit meldung.
func bis(t *testing.T, bedingung func() bool, meldung string) {
	t.Helper()
	frist := time.Now().Add(30 * time.Second)
	for !bedingung() {
		if time.Now().After(frist) {
			t.Fatal(meldung)
		}
		time.Sleep(time.Millisecond)
	}
}

// warteAufPlay startet Play in einer Goroutine, ruft danach dann und wartet
// höchstens 30 s auf das Ende.
func warteAufPlay(ctx context.Context, t *testing.T, s *services.PlayService, ablauf <-chan struct{}, dann func()) error {
	t.Helper()
	fertig := make(chan error, 1)
	go func() { fertig <- s.Play(ctx, ablauf) }()
	dann()
	select {
	case err := <-fertig:
		return err
	case <-time.After(30 * time.Second):
		t.Fatal("Play endet binnen 30 s nach dem zweiten Signal nicht")
		return nil
	}
}

// Abdeckung: LH-FA-20/Boundary — das zweite Signal (ablauf geschlossen)
// schließt die laufende Verbindung sofort, auch in einer Interaktion, die auf
// eine Antwort wartet, und bricht einen laufenden Aufbau ab; die
// unterbrochene Interaktion und der unterbrochene Aufbau sind kein Fehler, und
// keine weitere Session beginnt (LH-FA-20.a *Abbruchsignal*).
func TestPlayZweitesSignal(t *testing.T) {
	ziel := &fakeZiel{sessions: []*fakeEinspiel{{antworten: [][]schritt{{{blockiert: true}}}}}}
	s, err := services.NewPlayService(context.Background(), ladeRepo{rec: zwei()}, "r.yaml", ziel, services.PlayOptions{})
	if err != nil {
		t.Fatal(err)
	}
	ablauf := make(chan struct{})
	err = warteAufPlay(context.Background(), t, s, ablauf, func() {
		bis(t, func() bool { return len(ziel.liste()) >= 3 }, "Play wartet nicht binnen 30 s auf eine Antwort")
		close(ablauf)
	})
	if err != nil {
		t.Fatalf("zweites Signal in der Interaktion: %v", err)
	}
	if got := ziel.liste(); !reflect.DeepEqual(got, []string{"verbinde", "anfrage A", "naechste", "schliesse"}) {
		t.Fatalf("Ablauf %q", got)
	}

	ziel = &fakeZiel{verbindeBlockiert: true}
	s, err = services.NewPlayService(context.Background(), ladeRepo{rec: zwei()}, "r.yaml", ziel, services.PlayOptions{})
	if err != nil {
		t.Fatal(err)
	}
	ablauf = make(chan struct{})
	err = warteAufPlay(context.Background(), t, s, ablauf, func() {
		bis(t, func() bool { return len(ziel.liste()) >= 1 }, "Play beginnt binnen 30 s keinen Aufbau")
		close(ablauf)
	})
	if err != nil || !reflect.DeepEqual(ziel.liste(), []string{"verbinde"}) {
		t.Fatalf("zweites Signal im Aufbau: %v, %q", err, ziel.liste())
	}
}
