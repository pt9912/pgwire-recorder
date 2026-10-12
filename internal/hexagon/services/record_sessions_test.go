package services_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// neuLeer liefert einen RecordService mit --record-empty-sessions (an) oder
// ohne.
func neuLeer(t *testing.T, up *fakeUpstream, an bool) (*services.RecordService, *fakeRepo) {
	t.Helper()
	repo := &fakeRepo{}
	s, err := services.NewRecordService(context.Background(), up, repo, "rec.yaml", false, services.LeereSessions(an))
	if err != nil {
		t.Fatal(err)
	}
	return s, repo
}

// schreibungen zählt die Schreibungen der Aufzeichnung.
func schreibungen(r *fakeRepo) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.writes)
}

// Abdeckung: LH-FA-12/Boundary — mit LeereSessions wird jede eröffnete Session
// geschrieben, auch eine ohne abgeschlossene Interaktion, nach jedem dieser
// Enden: Ende durch den Client mit und ohne Terminate, Schreibfehler zum Client,
// Fehler des Use Case, Herunterfahren (das Zwangsende mit laufender Interaktion
// prüft TestRecordZwangsendeLeereSession); ohne die Option wird sie nicht geschrieben
// (LH-FA-12.a *Recording*, *Mit der Option*).
func TestRecordLeereSessionJedesEnde(t *testing.T) {
	enden := []struct {
		name string
		end  model.SessionEnd
	}{
		{"Client ohne Terminate", model.EndClosed},
		{"Client mit Terminate", model.EndTerminate},
		{"Schreibfehler zum Client", model.EndWriteFailed},
		{"Fehler des Use Case", model.EndFailed},
		{"Herunterfahren", model.EndShutdown},
	}
	for _, e := range enden {
		t.Run(e.name+" mit Option", func(t *testing.T) {
			s, repo := neuLeer(t, &fakeUpstream{}, true)
			id := session(t, s)
			err := s.CloseSession(context.Background(), id, e.end)
			if e.end == model.EndWriteFailed {
				if codeOf(err) != model.CodeConnectionLost {
					t.Fatalf("Fehler: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			rec := repo.last(t)
			if len(rec.Sessions) != 1 || rec.Sessions[0].ID != 1 || len(rec.Sessions[0].Interactions) != 0 {
				t.Fatalf("Sessions: %#v", rec.Sessions)
			}
			if !rec.EmptySessions {
				t.Fatal("Kennzeichnung empty_sessions fehlt")
			}
		})
		t.Run(e.name+" ohne Option", func(t *testing.T) {
			s, repo := neuLeer(t, &fakeUpstream{}, false)
			id := session(t, s)
			_ = s.CloseSession(context.Background(), id, e.end)
			if n := schreibungen(repo); n != 0 {
				t.Fatalf("Session ohne Interaktion geschrieben: %#v", repo.last(t))
			}
		})
	}
}

// Abdeckung: LH-FA-12/Boundary — eine Session, deren einzige Antwort den Client
// nicht erreicht, verliert diese Interaktion und wird mit der Option als Session
// ohne Interaktionen geschrieben (LH-FA-12.a *Mit der Option*).
func TestRecordLeereSessionNachSchreibfehler(t *testing.T) {
	s, repo := neuLeer(t, &fakeUpstream{}, true)
	id := session(t, s)
	if _, err := s.Query(context.Background(), id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if err := s.CloseSession(context.Background(), id, model.EndWriteFailed); codeOf(err) != model.CodeConnectionLost {
		t.Fatalf("Fehler: %v", err)
	}
	if rec := repo.last(t); len(rec.Sessions) != 1 || len(rec.Sessions[0].Interactions) != 0 {
		t.Fatalf("Sessions: %#v", rec.Sessions)
	}
}

// Abdeckung: LH-FA-12/Negative — eine Session mit einer nicht unterstützten
// Interaktion wird verworfen, mit und ohne LeereSessions, ob der Adapter
// EndUnsupported meldet oder der Use Case sie selbst als nicht unterstützt
// markiert hat, auch wenn sie keine abgeschlossene Interaktion trägt; sie
// erhält keine Kennung (LH-FA-12.a *Mit der Option*).
func TestRecordLeereSessionNichtUnterstuetztVerworfen(t *testing.T) {
	for _, an := range []bool{true, false} {
		t.Run(fmt.Sprintf("EndUnsupported, Option %v", an), func(t *testing.T) {
			s, repo := neuLeer(t, &fakeUpstream{}, an)
			if err := s.CloseSession(context.Background(), session(t, s), model.EndUnsupported); err != nil {
				t.Fatal(err)
			}
			if n := schreibungen(repo); n != 0 {
				t.Fatalf("verworfene Session geschrieben: %#v", repo.last(t))
			}
		})
		t.Run(fmt.Sprintf("Serverantwort nicht unterstützt, Option %v", an), func(t *testing.T) {
			s, repo := neuLeer(t, &fakeUpstream{queryErr: model.Errorf(model.CodeUnsupported, nil, "COPY")}, an)
			id := session(t, s)
			if _, err := s.Query(context.Background(), id, "FEHLER"); err == nil {
				t.Fatal("Fehler erwartet")
			}
			if err := s.CloseSession(context.Background(), id, model.EndFailed); err != nil {
				t.Fatal(err)
			}
			if n := schreibungen(repo); n != 0 {
				t.Fatalf("verworfene Session geschrieben: %#v", repo.last(t))
			}
		})
	}
}

// Abdeckung: LH-FA-12/Boundary — die Kennzeichnung empty_sessions steht in der
// Aufzeichnung, sobald der Lauf LeereSessions hat, auch ohne Session; ohne die
// Option steht sie nicht, auch mit Session (LH-FA-12.a *Kennzeichnung*).
func TestRecordKennzeichnung(t *testing.T) {
	s, repo := neuLeer(t, &fakeUpstream{}, true)
	if err := s.Finish(context.Background()); err != nil {
		t.Fatal(err)
	}
	if rec := repo.last(t); !rec.EmptySessions || len(rec.Sessions) != 0 {
		t.Fatalf("Lauf mit Option ohne Session: %#v", rec)
	}
	s, repo = neuLeer(t, &fakeUpstream{}, false)
	if err := s.CloseSession(context.Background(), session(t, s, "SELECT 1"), model.EndClosed); err != nil {
		t.Fatal(err)
	}
	if err := s.Finish(context.Background()); err != nil {
		t.Fatal(err)
	}
	if rec := repo.last(t); rec.EmptySessions || len(rec.Sessions) != 1 {
		t.Fatalf("Lauf ohne Option mit Session: %#v", rec)
	}
}

// Abdeckung: LH-FA-12/Boundary — die Kennung ist eine fortlaufende Zahl ab 1 ohne
// Lücken, vergeben beim Schreiben der Session in der Reihenfolge der Enden; eine
// verworfene Session erhält keine. Die Verbindungen enden hier nacheinander: Die
// zuletzt geöffnete endet zuerst (LH-FA-12.a *Recording*).
func TestRecordKennungNachReihenfolgeDerEnden(t *testing.T) {
	for _, an := range []bool{true, false} {
		t.Run(fmt.Sprintf("Option %v", an), func(t *testing.T) {
			s, repo := neuLeer(t, &fakeUpstream{}, an)
			ctx := context.Background()
			a := session(t, s, "A")
			b := session(t, s, "B")
			verworfen := session(t, s, "V")
			c := session(t, s, "C")
			ende := func(id model.SessionID, e model.SessionEnd) {
				t.Helper()
				if err := s.CloseSession(ctx, id, e); err != nil {
					t.Fatal(err)
				}
			}
			ende(c, model.EndClosed)
			ende(verworfen, model.EndUnsupported)
			ende(a, model.EndClosed)
			ende(b, model.EndClosed)
			var got []string
			for _, sess := range repo.last(t).Sessions {
				got = append(got, fmt.Sprintf("%d:%s", sess.ID, sess.Interactions[0].Request.SQL))
			}
			if want := "1:C 2:A 3:B"; strings.Join(got, " ") != want {
				t.Fatalf("Sessions %v, erwartet %s", got, want)
			}
		})
	}
}

// Abdeckung: LH-FA-12/Boundary — enden Sessions gleichzeitig, vergibt jedes Ende
// genau eine Kennung: Die Kennungen sind eindeutig und ohne Lücken, mit und ohne
// Interaktion, welche Session welche erhält, prüft der Test nicht
// (LH-FA-12.a *Gleichzeitiges Ende*).
func TestRecordGleichzeitigesEnde(t *testing.T) {
	s, repo := neuLeer(t, &fakeUpstream{}, true)
	ctx := context.Background()
	const n = 24
	ids := make([]model.SessionID, n)
	for i := range ids {
		if i%2 == 0 {
			ids[i] = session(t, s, fmt.Sprintf("SELECT %d", i))
		} else {
			ids[i] = session(t, s)
		}
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if err := s.CloseSession(ctx, id, model.EndClosed); err != nil {
				t.Error(err)
			}
		}()
	}
	close(start)
	wg.Wait()
	rec := repo.last(t)
	if len(rec.Sessions) != n {
		t.Fatalf("Sessions: %d, erwartet %d", len(rec.Sessions), n)
	}
	gesehen := map[int]bool{}
	for i, sess := range rec.Sessions {
		if sess.ID != i+1 || gesehen[sess.ID] {
			t.Fatalf("Kennung %d an Stelle %d", sess.ID, i+1)
		}
		gesehen[sess.ID] = true
	}
}

// Abdeckung: LH-FA-12/Happy — zwei gleichzeitig offene Verbindungen sind zwei
// Sessions mit je der Reihenfolge ihrer Interaktionen, auch wenn die Anfragen
// abwechselnd kommen; eine dritte Verbindung ohne Anfrage ist keine Session
// (LH-FA-12.a *Recording*).
func TestRecordParalleleSessions(t *testing.T) {
	s, repo := neu(t, &fakeUpstream{})
	ctx := context.Background()
	a, _, err := s.OpenSession(ctx, map[string]string{"user": "a"})
	if err != nil {
		t.Fatal(err)
	}
	b, _, err := s.OpenSession(ctx, map[string]string{"user": "b"})
	if err != nil {
		t.Fatal(err)
	}
	ohne, _, err := s.OpenSession(ctx, map[string]string{"user": "ohne"})
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range []struct {
		id  model.SessionID
		sql string
	}{{a, "A1"}, {b, "B1"}, {b, "B2"}, {a, "A2"}, {a, "A3"}, {b, "B3"}} {
		if _, err := s.Query(ctx, q.id, q.sql); err != nil {
			t.Fatal(err)
		}
		s.Delivered(ctx, q.id)
	}
	for _, id := range []model.SessionID{ohne, b, a} {
		if err := s.CloseSession(ctx, id, model.EndTerminate); err != nil {
			t.Fatal(err)
		}
	}
	rec := repo.last(t)
	if len(rec.Sessions) != 2 {
		t.Fatalf("Sessions: %#v", rec.Sessions)
	}
	for i, want := range []struct {
		user string
		sqls string
	}{{"b", "B1 B2 B3"}, {"a", "A1 A2 A3"}} {
		got := rec.Sessions[i]
		var sqls []string
		for j, in := range got.Interactions {
			if in.Sequence != j+1 {
				t.Errorf("Session %d: Interaktion %d trägt die Nummer %d", got.ID, j+1, in.Sequence)
			}
			sqls = append(sqls, in.Request.SQL)
		}
		if got.Startup["user"] != want.user || strings.Join(sqls, " ") != want.sqls {
			t.Errorf("Session %d: Benutzer %q, Anfragen %v; erwartet %q, %s", got.ID, got.Startup["user"], sqls, want.user, want.sqls)
		}
	}
}

// Abdeckung: LH-FA-12/Boundary, LH-FA-13/Boundary — beim Zwangsende mit
// LeereSessions nennt PGR-E4006 für eine Session ohne abgeschlossene
// Interaktion die Kennung, unter der sie geschrieben wird (hier die zweite, weil
// schon eine Session geschrieben ist); die Session steht in der Aufzeichnung.
// Ist die laufende Interaktion nicht unterstützt, nennt die Meldung diesen Grund
// und keine Kennung, und die Session wird nicht geschrieben (LH-FA-12.a *Mit der
// Option*, LH-FA-13.a *Meldung*). Ohne laufende und ohne abgeschlossene Interaktion
// ist das Zwangsende das im Aufbau: Die Session wird auch mit der Option nicht
// geschrieben, und es gibt keine Meldung (LH-FA-12.a *Angenommene Verbindung*);
// mit einer abgeschlossenen Interaktion ohne laufende wird sie geschrieben.
func TestRecordZwangsendeLeereSession(t *testing.T) {
	t.Run("Kennung genannt", func(t *testing.T) {
		up := &fakeUpstream{}
		s, repo := neuLeer(t, up, true)
		if err := s.CloseSession(context.Background(), session(t, s, "SELECT 1"), model.EndClosed); err != nil {
			t.Fatal(err)
		}
		id := session(t, s)
		client(t, s, id, cParse())
		err := s.CloseSession(context.Background(), id, model.EndForced)
		if codeOf(err) != model.CodeShutdownTimeout || !strings.Contains(err.Error(), "Interaktion 1 nicht abgeschlossen und verworfen; die Session wird als Session 2 geschrieben") {
			t.Fatalf("Fehler: %v", err)
		}
		if rec := repo.last(t); len(rec.Sessions) != 2 || len(rec.Sessions[1].Interactions) != 0 {
			t.Fatalf("Sessions: %#v", rec.Sessions)
		}
	})
	t.Run("im Aufbau", func(t *testing.T) {
		s, repo := neuLeer(t, &fakeUpstream{}, true)
		id := session(t, s)
		if err := s.CloseSession(context.Background(), id, model.EndForced); err != nil {
			t.Fatalf("Meldung beim Zwangsende im Aufbau: %v", err)
		}
		if n := schreibungen(repo); n != 0 {
			t.Fatalf("Session im Aufbau geschrieben: %#v", repo.last(t))
		}
		if err := s.Finish(context.Background()); err != nil {
			t.Fatal(err)
		}
		if rec := repo.last(t); len(rec.Sessions) != 0 || !rec.EmptySessions {
			t.Fatalf("Aufzeichnung: %#v", rec)
		}
	})
	t.Run("nicht unterstützt", func(t *testing.T) {
		up := &fakeUpstream{}
		s, repo := neuLeer(t, up, true)
		id := session(t, s)
		up.letzte.sendErr = model.Errorf(model.CodeUnsupported, nil, "COPY")
		if err := s.ClientMessage(context.Background(), id, cSync()); codeOf(err) != model.CodeUnsupported {
			t.Fatalf("ClientMessage: %v", err)
		}
		err := s.CloseSession(context.Background(), id, model.EndForced)
		if codeOf(err) != model.CodeShutdownTimeout || !strings.Contains(err.Error(), "die Session wird wegen einer nicht unterstützten Interaktion (PGR-E6001) nicht geschrieben") {
			t.Fatalf("Fehler: %v", err)
		}
		if n := schreibungen(repo); n != 0 {
			t.Fatalf("verworfene Session geschrieben: %#v", repo.last(t))
		}
	})
}

// aufbauFehler ist ein Upstream, dessen Aufbau scheitert (fehler) oder mit einer
// Fehlerantwort statt einer Session endet (fehler nil).
type aufbauFehler struct{ fehler error }

func (a aufbauFehler) Open(context.Context, map[string]string) (driven.UpstreamSession, []model.Response, error) {
	if a.fehler != nil {
		return nil, nil, a.fehler
	}
	return nil, []model.Response{{Type: model.ResponseErrorResponse, Fields: map[string]string{"C": "28000"}}}, nil
}

// Abdeckung: LH-FA-12/Negative — scheitert der Aufbau zum Upstream, mit Fehler
// oder mit einer Fehlerantwort des Servers, ist die Kennung 0 und keine Verbindung
// eine Session, auch mit LeereSessions: Ein CloseSession darauf ändert nichts, und
// die Aufzeichnung hat keine Session, aber die Kennzeichnung (LH-FA-12.a
// *Angenommene Verbindung*).
func TestRecordAufbauFehlerKeineSession(t *testing.T) {
	for name, up := range map[string]aufbauFehler{
		"Fehler":        {fehler: model.Errorf(model.CodeUpstream, nil, "nicht erreichbar")},
		"Fehlerantwort": {},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			repo := &fakeRepo{}
			s, err := services.NewRecordService(ctx, up, repo, "rec.yaml", false, services.LeereSessions(true))
			if err != nil {
				t.Fatal(err)
			}
			id, _, err := s.OpenSession(ctx, map[string]string{"user": "app"})
			if id != 0 || (up.fehler != nil) != (err != nil) {
				t.Fatalf("OpenSession: Kennung %d, Fehler %v", id, err)
			}
			if err := s.CloseSession(ctx, id, model.EndClosed); err != nil {
				t.Fatal(err)
			}
			if n := schreibungen(repo); n != 0 {
				t.Fatalf("Session ohne Aufbau geschrieben: %#v", repo.last(t))
			}
			if err := s.Finish(ctx); err != nil {
				t.Fatal(err)
			}
			if rec := repo.last(t); len(rec.Sessions) != 0 || !rec.EmptySessions {
				t.Fatalf("Aufzeichnung: %#v", rec)
			}
		})
	}
}
