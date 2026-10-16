package services_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
	"github.com/pt9912/pgwire-recorder/internal/testuhr"
	"github.com/pt9912/pgwire-recorder/internal/testzeit"
)

// neuMitUhr liefert den Service mit Zeitaufzeichnung nach der Uhr.
func neuMitUhr(t *testing.T, up driven.Upstream, uhr driven.Uhr) (*services.RecordService, *fakeRepo) {
	t.Helper()
	repo := &fakeRepo{}
	s, err := services.NewRecordService(context.Background(), up, repo, "rec.yaml", false, services.Zeitangaben(uhr, true))
	if err != nil {
		t.Fatal(err)
	}
	return s, repo
}

// offsets liefert offset_ms jeder Interaktion der Session n der letzten Aufzeichnung;
// -1 steht für ein fehlendes Feld.
func offsets(t *testing.T, repo *fakeRepo, n int) []int64 {
	t.Helper()
	var out []int64
	for _, in := range repo.last(t).Sessions[n].Interactions {
		if in.OffsetMS == nil {
			out = append(out, -1)
			continue
		}
		out = append(out, *in.OffsetMS)
	}
	return out
}

func gleichZahlen(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// frage führt eine einfache Anfrage aus und meldet sie als zugestellt.
func frage(t *testing.T, s *services.RecordService, id model.SessionID, sql string) {
	t.Helper()
	if _, err := s.Query(context.Background(), id, sql); err != nil {
		t.Fatalf("Query %q: %v", sql, err)
	}
	s.Delivered(context.Background(), id)
}

// uhrUpstream rückt die Uhr während Open vor: Der Aufbau der Upstream-Session
// kostet Zeit.
type uhrUpstream struct {
	driven.Upstream
	uhr *testuhr.Uhr
	ms  int64
}

func (u uhrUpstream) Open(ctx context.Context, startup map[string]string, anmeldung model.Anmeldung) (driven.UpstreamSession, []model.Response, error) {
	u.uhr.Vorruecken(u.ms)
	return u.Upstream.Open(ctx, startup, anmeldung)
}

// Abdeckung: LH-FA-21/Happy — mit Zeitaufzeichnung trägt jede einfache Anfrage
// offset_ms: der Unterschied der Lesung beim Aufruf von Query und der Lesung beim
// Aufruf von OpenSession, in ganzen Millisekunden, ohne Rundung (LH-FA-21.a
// *Messpunkt*); die Uhr stand beim Beginn nicht auf 0.
func TestRecordZeitEinfach(t *testing.T) {
	uhr := testuhr.Neu()
	s, repo := neuMitUhr(t, &fakeUpstream{}, uhr)
	uhr.Vorruecken(100)
	id := session(t, s)
	uhr.Vorruecken(40)
	frage(t, s, id, "SELECT 1")
	uhr.Vorruecken(5)
	frage(t, s, id, "SELECT 2")
	schliessen(t, s, id, model.EndClosed)
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{40, 45}) {
		t.Fatalf("offset_ms %v, erwartet [40 45]", got)
	}
}

// Abdeckung: LH-FA-21/Boundary — die erste Interaktion kann offset_ms 0 tragen: Das
// Feld steht dann mit dem Wert 0 da und fehlt nicht (LH-FA-21.a *Aufzeichnen*).
func TestRecordZeitNullSteht(t *testing.T) {
	uhr := testuhr.Neu()
	s, repo := neuMitUhr(t, &fakeUpstream{}, uhr)
	uhr.Vorruecken(7)
	id := session(t, s)
	frage(t, s, id, "SELECT 1")
	schliessen(t, s, id, model.EndClosed)
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{0}) {
		t.Fatalf("offset_ms %v, erwartet [0]", got)
	}
}

// Abdeckung: LH-FA-21/Boundary — Beginn der Session ist der Aufruf von OpenSession vor dem
// Aufbau der Upstream-Session: Kostet der Aufbau 30 ms, trägt die erste Anfrage
// danach 30 (LH-FA-21.a *Messpunkt*).
func TestRecordZeitBeginnVorAufbau(t *testing.T) {
	uhr := testuhr.Neu()
	s, repo := neuMitUhr(t, uhrUpstream{Upstream: &fakeUpstream{}, uhr: uhr, ms: 30}, uhr)
	id := session(t, s)
	frage(t, s, id, "SELECT 1")
	schliessen(t, s, id, model.EndClosed)
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{30}) {
		t.Fatalf("offset_ms %v, erwartet [30]", got)
	}
}

// Abdeckung: LH-FA-21/Boundary — Ankunft einer einfachen Anfrage ist der Aufruf von Query,
// gelesen vor jedem Warten: Wartet der Upstream 500 ms auf die Antwort, ändert das den
// Wert nicht (LH-FA-21.a *Messpunkt*).
func TestRecordZeitAnkunftVorUpstream(t *testing.T) {
	uhr := testuhr.Neu()
	up := &fakeUpstream{}
	s, repo := neuMitUhr(t, up, uhr)
	id := session(t, s)
	uhr.Vorruecken(10)
	up.letzte.queryLaeuft = make(chan struct{})
	up.letzte.queryHalt = make(chan struct{})
	fertig := make(chan error, 1)
	go func() {
		_, err := s.Query(context.Background(), id, "SELECT 1")
		fertig <- err
	}()
	<-up.letzte.queryLaeuft
	uhr.Vorruecken(500)
	close(up.letzte.queryHalt)
	if err := <-fertig; err != nil {
		t.Fatal(err)
	}
	s.Delivered(context.Background(), id)
	schliessen(t, s, id, model.EndClosed)
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{10}) {
		t.Fatalf("offset_ms %v, erwartet [10]", got)
	}
}

// meldendeUhr meldet jede Lesung auf gelesen.
type meldendeUhr struct {
	*testuhr.Uhr
	gelesen chan struct{}
}

func (u meldendeUhr) Jetzt() int64 {
	v := u.Uhr.Jetzt()
	select {
	case u.gelesen <- struct{}{}:
	default:
	}
	return v
}

// Abdeckung: LH-FA-21/Boundary — Ankunft einer einfachen Anfrage ist gelesen, bevor Query auf
// etwas wartet: Wartet Query auf die Antworten einer Extended-Interaktion mit gesendetem
// Sync, ändert die Wartezeit den Wert nicht (LH-FA-21.a *Messpunkt*).
func TestRecordZeitAnkunftVorWarten(t *testing.T) {
	basis := testuhr.Neu()
	uhr := meldendeUhr{Uhr: basis, gelesen: make(chan struct{}, 64)}
	up := &fakeUpstream{}
	s, repo := neuMitUhr(t, up, uhr)
	id := session(t, s)
	client(t, s, id, cParse(), cSync())
	basis.Vorruecken(10)
	for len(uhr.gelesen) > 0 {
		<-uhr.gelesen
	}
	fertig := make(chan error, 1)
	go func() {
		_, err := s.Query(context.Background(), id, "SELECT 1")
		fertig <- err
	}()
	select {
	case <-uhr.gelesen:
	case <-testzeit.Nach(2000):
		t.Fatal("Query liest die Uhr nicht, bevor sie auf die Extended-Interaktion wartet")
	}
	basis.Vorruecken(300)
	server(t, s, up, id, sParse(), sRFQ())
	if err := <-fertig; err != nil {
		t.Fatal(err)
	}
	s.Delivered(context.Background(), id)
	schliessen(t, s, id, model.EndClosed)
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{0, 10}) {
		t.Fatalf("offset_ms %v, erwartet [0 10]", got)
	}
}

// ersteArten ist die erste Client-Nachricht einer Extended-Interaktion je Art und die
// Nachrichten, die die Interaktion danach bis zum Sync vervollständigen.
func ersteArten() map[string][]model.ClientMessage {
	return map[string][]model.ClientMessage{
		"Parse":    {cParse(), cSync()},
		"Bind":     {cBind(), cSync()},
		"Describe": {cDescribe(), cSync()},
		"Execute":  {cExecute(), cSync()},
		"Close":    {{Type: model.ClientClose, Target: model.TargetStatement}, cSync()},
		"Flush":    {cFlush(), cSync()},
		"Sync":     {cSync()},
	}
}

// Abdeckung: LH-FA-21/Happy, LH-FA-21/Boundary — Ankunft einer Extended-Interaktion ist
// die erste Client-Nachricht der ersten Gruppe, gleich welche Art sie hat (Parse, Bind,
// Describe, Execute, Close, Flush, Sync): offset_ms ist der Unterschied zu dieser Lesung,
// nicht zur Lesung bei der zweiten Nachricht, beim Sync oder beim ReadyForQuery (LH-FA-21.a
// *Messpunkt*).
func TestRecordZeitExtendedErsteNachricht(t *testing.T) {
	for art, msgs := range ersteArten() {
		uhr := testuhr.Neu()
		up := &fakeUpstream{}
		s, repo := neuMitUhr(t, up, uhr)
		id := session(t, s)
		uhr.Vorruecken(20)
		for _, m := range msgs {
			client(t, s, id, m)
			uhr.Vorruecken(5)
		}
		uhr.Vorruecken(100)
		server(t, s, up, id, sRFQ())
		schliessen(t, s, id, model.EndClosed)
		if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{20}) {
			t.Errorf("erste Nachricht %s: offset_ms %v, erwartet [20]", art, got)
		}
	}
}

// Abdeckung: LH-FA-21/Boundary — das Feld einer Extended-Interaktion mit mehreren Gruppen
// steht auf der Ebene der Interaktion und trägt die Lesung der ersten Nachricht der ersten
// Gruppe, nicht die der ersten Nachricht einer folgenden Gruppe (LH-FA-21.a *Messpunkt*).
func TestRecordZeitExtendedMehrereGruppen(t *testing.T) {
	uhr := testuhr.Neu()
	up := &fakeUpstream{}
	s, repo := neuMitUhr(t, up, uhr)
	id := session(t, s)
	uhr.Vorruecken(10)
	client(t, s, id, cParse(), cFlush())
	server(t, s, up, id, sParse())
	uhr.Vorruecken(10)
	client(t, s, id, cBind(), cSync())
	server(t, s, up, id, sBind(), sRFQ())
	schliessen(t, s, id, model.EndClosed)
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{10}) {
		t.Fatalf("offset_ms %v, erwartet [10]", got)
	}
}

// Abdeckung: LH-FA-21/Happy, LH-FA-21/Boundary — die Werte einer Session steigen nicht ab:
// Eine einfache Anfrage, zwei Extended-Interaktionen, deren zweite beginnt, bevor die
// Antworten der ersten da sind (Pipelining), und eine einfache Anfrage danach tragen die
// Lesungen ihrer Ankunft in der Reihenfolge ihrer Nummern (LH-FA-21.a *Messpunkt*).
func TestRecordZeitNichtAbsteigend(t *testing.T) {
	uhr := testuhr.Neu()
	up := &fakeUpstream{}
	s, repo := neuMitUhr(t, up, uhr)
	id := session(t, s)
	uhr.Vorruecken(10)
	frage(t, s, id, "SELECT 1")
	uhr.Vorruecken(10)
	client(t, s, id, cParse(), cSync())
	uhr.Vorruecken(10)
	client(t, s, id, cBind(), cSync())
	uhr.Vorruecken(100)
	server(t, s, up, id, sParse(), sRFQ())
	uhr.Vorruecken(100)
	server(t, s, up, id, sBind(), sRFQ())
	uhr.Vorruecken(10)
	frage(t, s, id, "SELECT 2")
	schliessen(t, s, id, model.EndClosed)
	got := offsets(t, repo, 0)
	if !gleichZahlen(got, []int64{10, 20, 30, 240}) {
		t.Fatalf("offset_ms %v, erwartet [10 20 30 240]", got)
	}
}

// Abdeckung: LH-FA-21/Happy — die Lebendprüfung trägt das Feld wie jede Anfrage
// (LH-FA-21.a *Messpunkt*, LH-FA-09.a).
func TestRecordZeitLebendpruefung(t *testing.T) {
	uhr := testuhr.Neu()
	s, repo := neuMitUhr(t, &fakeUpstream{}, uhr)
	id := session(t, s)
	uhr.Vorruecken(3)
	frage(t, s, id, "-- ping")
	uhr.Vorruecken(4)
	frage(t, s, id, "")
	schliessen(t, s, id, model.EndClosed)
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{3, 7}) {
		t.Fatalf("offset_ms %v, erwartet [3 7]", got)
	}
}

// Abdeckung: LH-FA-21/Boundary — jede Session misst von ihrem eigenen Beginn, auch wenn
// zwei Sessions gleichzeitig offen sind (LH-FA-21.a *Messpunkt*).
func TestRecordZeitJedeSessionEigenerBeginn(t *testing.T) {
	uhr := testuhr.Neu()
	s, repo := neuMitUhr(t, &fakeUpstream{}, uhr)
	a := session(t, s)
	uhr.Vorruecken(100)
	b := session(t, s)
	uhr.Vorruecken(5)
	frage(t, s, b, "SELECT 1")
	uhr.Vorruecken(5)
	frage(t, s, a, "SELECT 2")
	schliessen(t, s, a, model.EndClosed)
	schliessen(t, s, b, model.EndClosed)
	rec := repo.last(t)
	if len(rec.Sessions) != 2 {
		t.Fatalf("Sessions: %d", len(rec.Sessions))
	}
	// Die Session a wird zuerst geschlossen und steht an erster Stelle.
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{110}) {
		t.Errorf("Session a: offset_ms %v, erwartet [110]", got)
	}
	if got := offsets(t, repo, 1); !gleichZahlen(got, []int64{5}) {
		t.Errorf("Session b: offset_ms %v, erwartet [5]", got)
	}
}

// Abdeckung: LH-FA-21/Negative — eine Interaktion, die unvollständig endet, wird nicht
// übernommen und hat keinen Wert; die vollständigen davor tragen ihren (LH-FA-21.a
// *Messpunkt*, LH-FA-13.a).
func TestRecordZeitUnvollstaendigOhneWert(t *testing.T) {
	uhr := testuhr.Neu()
	s, repo := neuMitUhr(t, &fakeUpstream{}, uhr)
	id := session(t, s)
	uhr.Vorruecken(10)
	frage(t, s, id, "SELECT 1")
	uhr.Vorruecken(10)
	client(t, s, id, cParse(), cFlush())
	if err := s.CloseSession(context.Background(), id, model.EndClosed); codeOf(err) != model.CodeConnectionLost {
		t.Fatalf("CloseSession: %v, erwartet %s", err, model.CodeConnectionLost)
	}
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{10}) {
		t.Fatalf("offset_ms %v, erwartet [10]", got)
	}
}

// zaehlendeUhr zählt die Lesungen und lässt den Test scheitern, wenn jemand wartet.
type zaehlendeUhr struct {
	*testuhr.Uhr
	t        *testing.T
	lesungen atomic.Int64
}

func (u *zaehlendeUhr) Jetzt() int64 {
	u.lesungen.Add(1)
	return u.Uhr.Jetzt()
}

func (u *zaehlendeUhr) WarteBis(context.Context, int64) error {
	u.t.Error("WarteBis aufgerufen")
	return nil
}

// Abdeckung: LH-FA-21/Boundary — mit Zeitaufzeichnung liest der Service die Uhr und wartet
// nie auf sie (LH-FA-21.a *Aufzeichnen*).
func TestRecordZeitWartetNie(t *testing.T) {
	uhr := &zaehlendeUhr{Uhr: testuhr.Neu(), t: t}
	up := &fakeUpstream{}
	s, _ := neuMitUhr(t, up, uhr)
	id := session(t, s, "SELECT 1")
	client(t, s, id, cParse(), cSync())
	server(t, s, up, id, sParse(), sRFQ())
	schliessen(t, s, id, model.EndClosed)
	if uhr.lesungen.Load() == 0 {
		t.Fatal("die Uhr wurde nie gelesen")
	}
}

// Abdeckung: LH-FA-21/Boundary, LH-FA-21/Negative — ohne Zeitaufzeichnung fehlt offset_ms an jeder
// Interaktion, einfach und Extended, und der Service liest die Uhr nie, auch nicht, wenn
// er eine bekommt (LH-FA-21.a *Aufzeichnen*).
func TestRecordOhneZeitangaben(t *testing.T) {
	uhr := &zaehlendeUhr{Uhr: testuhr.Neu(), t: t}
	repo := &fakeRepo{}
	up := &fakeUpstream{}
	s, err := services.NewRecordService(context.Background(), up, repo, "rec.yaml", false, services.Zeitangaben(uhr, false))
	if err != nil {
		t.Fatal(err)
	}
	id := session(t, s)
	uhr.Vorruecken(10)
	frage(t, s, id, "SELECT 1")
	client(t, s, id, cParse(), cSync())
	server(t, s, up, id, sParse(), sRFQ())
	schliessen(t, s, id, model.EndClosed)
	if got := offsets(t, repo, 0); !gleichZahlen(got, []int64{-1, -1}) {
		t.Fatalf("offset_ms %v, erwartet ohne Feld", got)
	}
	if n := uhr.lesungen.Load(); n != 0 {
		t.Fatalf("%d Lesungen ohne Zeitaufzeichnung", n)
	}
}
