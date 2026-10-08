package services_test

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// Abdeckung: LH-FA-13/Boundary — Record im Kern: meldet der Adapter das
// Zwangsende (EndForced), ist eine laufende Interaktion, einfach oder Extended,
// PGR-E4006 der Klasse 4; die Meldung nennt die Kennung, unter der die Session
// geschrieben wird, oder dass sie ohne abgeschlossene Interaktion nicht
// geschrieben wird, und die Nummer, die die verworfene Interaktion getragen
// hätte. Die abgeschlossenen Interaktionen bleiben, die laufende nicht. Ohne
// laufende Interaktion endet die Session ohne Fehler (LH-FA-13.a *Zwangsende*).
func TestRecordZwangsende(t *testing.T) {
	cases := []struct {
		name     string
		vorher   []string
		ablauf   func(t *testing.T, s *services.RecordService, up *fakeUpstream, id model.SessionID)
		wantErr  string
		wantText string
		wantSQL  []string
	}{
		{"Extended-Interaktion vor dem Sync", []string{"SELECT 1"}, func(t *testing.T, s *services.RecordService, _ *fakeUpstream, id model.SessionID) {
			client(t, s, id, cParse(), cBind())
		}, model.CodeShutdownTimeout, "Interaktion 2 nicht abgeschlossen und verworfen; die Session wird als Session 1 geschrieben", []string{"SELECT 1"}},
		{"Extended-Interaktion nach dem Sync", []string{"SELECT 1", "SELECT 2"}, func(t *testing.T, s *services.RecordService, _ *fakeUpstream, id model.SessionID) {
			client(t, s, id, cParse(), cSync())
		}, model.CodeShutdownTimeout, "Interaktion 3 nicht abgeschlossen und verworfen; die Session wird als Session 1 geschrieben", []string{"SELECT 1", "SELECT 2"}},
		{"ohne abgeschlossene Interaktion", nil, func(t *testing.T, s *services.RecordService, _ *fakeUpstream, id model.SessionID) {
			client(t, s, id, cParse())
		}, model.CodeShutdownTimeout, "Interaktion 1 nicht abgeschlossen und verworfen; die Session wird ohne abgeschlossene Interaktion nicht geschrieben", nil},
		{"ohne laufende Interaktion", []string{"SELECT 1"}, func(*testing.T, *services.RecordService, *fakeUpstream, model.SessionID) {}, "", "", []string{"SELECT 1"}},
		{"abgeschlossen, nicht zugestellt", []string{"SELECT 1"}, func(t *testing.T, s *services.RecordService, up *fakeUpstream, id model.SessionID) {
			client(t, s, id, cSync())
			up.letzte.empfang <- []model.Response{sRFQ()}
			if _, err := s.AwaitServer(context.Background(), id); err != nil {
				t.Fatal(err)
			}
		}, "", "", []string{"SELECT 1", ""}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			up := &fakeUpstream{}
			s, repo := neu(t, up)
			id := session(t, s, c.vorher...)
			c.ablauf(t, s, up, id)
			err := s.CloseSession(context.Background(), id, model.EndForced)
			if codeOf(err) != c.wantErr || (c.wantText != "" && !strings.Contains(err.Error(), c.wantText)) {
				t.Fatalf("Fehler %v, erwartet %q mit %q", err, c.wantErr, c.wantText)
			}
			if c.wantErr != "" {
				if ms := model.Meldungen(err); ms[0].ExitCode() != 4 || !strings.HasPrefix(ms[0].Text, "Netzwerk [PGR-E4006]: ") {
					t.Fatalf("Meldung %#v", ms[0])
				}
			}
			select {
			case <-up.letzte.zu:
			default:
				t.Fatal("Upstream nicht geschlossen")
			}
			if c.wantSQL == nil {
				repo.mu.Lock()
				n := len(repo.writes)
				repo.mu.Unlock()
				if n != 0 {
					t.Fatalf("Session ohne abgeschlossene Interaktion geschrieben: %#v", repo.last(t))
				}
				return
			}
			var sqls []string
			for _, in := range repo.last(t).Sessions[0].Interactions {
				sqls = append(sqls, in.Request.SQL)
			}
			if !reflect.DeepEqual(sqls, c.wantSQL) {
				t.Fatalf("Interaktionen %q, erwartet %q", sqls, c.wantSQL)
			}
		})
	}
}

// Abdeckung: LH-FA-13/Boundary — Record im Kern: wartet eine einfache Anfrage
// beim Zwangsende auf den Upstream, ist sie die laufende Interaktion
// (PGR-E4006); die Anfrage endet mit ErrSessionEnded, und die Session wird mit
// der nächsten Kennung geschrieben, ohne die Anfrage.
func TestRecordZwangsendeEinfacheAnfrage(t *testing.T) {
	up := &fakeUpstream{}
	s, repo := neu(t, up)
	erste := session(t, s, "SELECT 0")
	schliessen(t, s, erste, model.EndClosed)
	id := session(t, s, "SELECT 1")
	up.letzte.queryLaeuft = make(chan struct{})
	up.letzte.queryHalt = make(chan struct{})
	ergebnis := make(chan error, 1)
	go func() {
		_, err := s.Query(context.Background(), id, "SELECT pg_sleep(60)")
		ergebnis <- err
	}()
	<-up.letzte.queryLaeuft
	var err error
	warte(t, "CloseSession", func() { err = s.CloseSession(context.Background(), id, model.EndForced) })
	if codeOf(err) != model.CodeShutdownTimeout || !strings.Contains(err.Error(), "Interaktion 2 nicht abgeschlossen und verworfen; die Session wird als Session 2 geschrieben") {
		t.Fatalf("Fehler %v", err)
	}
	close(up.letzte.queryHalt)
	select {
	case qerr := <-ergebnis:
		if qerr != model.ErrSessionEnded {
			t.Fatalf("Query nach dem Zwangsende: %v", qerr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Query kehrt nach dem Zwangsende nicht zurück")
	}
	rec := repo.last(t)
	if len(rec.Sessions) != 2 || len(rec.Sessions[1].Interactions) != 1 || rec.Sessions[1].Interactions[0].Request.SQL != "SELECT 1" {
		t.Fatalf("Aufzeichnung %#v", rec.Sessions)
	}
}
