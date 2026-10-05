package services

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

type fakeUpstream struct {
	queryErr error
	mu       sync.Mutex
	// letzte ist die zuletzt geöffnete Session.
	letzte *fakeSession
}

func (f *fakeUpstream) Open(context.Context, map[string]string) (driven.UpstreamSession, []model.Response, error) {
	session := neueFakeSession(f.queryErr)
	f.mu.Lock()
	f.letzte = session
	f.mu.Unlock()
	return session, []model.Response{
		{Type: model.ResponseParameterStatus, Name: "server_version", Value: "17.0"},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}, nil
}

// fakeSession ist eine Upstream-Session ohne Server: Send sammelt die
// Gruppen, Receive liefert die Folgen aus empfang der Reihe nach und wartet,
// solange keine bereitsteht. Close beendet ein wartendes Receive mit PGR-E4003.
type fakeSession struct {
	err     error
	mu      sync.Mutex
	gesendet [][]model.ClientMessage
	empfang chan []model.Response
	sendErr error
	zuEin   sync.Once
	zu      chan struct{}
	// queryLaeuft meldet den Beginn einer Query, die danach bis queryHalt wartet;
	// beide sind nil, wenn Query nicht wartet.
	queryLaeuft chan struct{}
	queryHalt   chan struct{}
}

func neueFakeSession(err error) *fakeSession {
	return &fakeSession{err: err, empfang: make(chan []model.Response, 16), zu: make(chan struct{})}
}

func (f *fakeSession) Query(_ context.Context, sql string) ([]model.Response, error) {
	if f.queryHalt != nil {
		f.queryLaeuft <- struct{}{}
		<-f.queryHalt
	}
	if f.err != nil && sql == "FEHLER" {
		return nil, f.err
	}
	if sql == "EXTENDED-ANTWORT" {
		return []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}, nil
	}
	return []model.Response{
		{Type: model.ResponseCommandComplete, Tag: sql},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}, nil
}

func (f *fakeSession) Send(_ context.Context, msgs []model.ClientMessage) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.sendErr != nil {
		return f.sendErr
	}
	f.gesendet = append(f.gesendet, append([]model.ClientMessage(nil), msgs...))
	return nil
}

func (f *fakeSession) gruppen() [][]model.ClientMessage {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.gesendet
}

func (f *fakeSession) Receive(context.Context) ([]model.Response, error) {
	select {
	case out := <-f.empfang:
		return out, nil
	case <-f.zu:
		return nil, model.Errorf(model.CodeConnectionLost, nil, "Upstream geschlossen")
	}
}

func (f *fakeSession) Close() error {
	f.zuEin.Do(func() { close(f.zu) })
	return nil
}

type fakeRepo struct {
	mu     sync.Mutex
	writes []model.Recording
}

func (f *fakeRepo) Prepare(context.Context, string, bool) error { return nil }
func (f *fakeRepo) Write(_ context.Context, _ string, rec model.Recording) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.writes = append(f.writes, rec)
	return nil
}
func (f *fakeRepo) Load(context.Context, string) (model.Recording, error) {
	return model.Recording{}, nil
}

func (f *fakeRepo) last(t *testing.T) model.Recording {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.writes) == 0 {
		t.Fatal("keine Aufzeichnung geschrieben")
	}
	return f.writes[len(f.writes)-1]
}

func neu(t *testing.T, up driven.Upstream) (*RecordService, *fakeRepo) {
	t.Helper()
	repo := &fakeRepo{}
	s, err := NewRecordService(context.Background(), up, repo, "rec.yaml", false)
	if err != nil {
		t.Fatal(err)
	}
	return s, repo
}

func session(t *testing.T, s *RecordService, queries ...string) model.SessionID {
	t.Helper()
	ctx := context.Background()
	id, _, err := s.OpenSession(ctx, map[string]string{"user": "app"})
	if err != nil || id == 0 {
		t.Fatalf("OpenSession: id=%d err=%v", id, err)
	}
	for _, q := range queries {
		if _, err := s.Query(ctx, id, q); err != nil {
			t.Fatalf("Query %q: %v", q, err)
		}
		s.Delivered(ctx, id)
	}
	return id
}

// Die Interaktionen einer Session stehen mit
// fortlaufender Nummer, Anfrage und Antworten in der Aufzeichnung; die Session
// trägt Startup-Parameter und Serverparameter.
func TestRecordSessionMitInteraktionen(t *testing.T) {
	s, repo := neu(t, &fakeUpstream{})
	id := session(t, s, "SELECT 1", "SELECT 2")
	if err := s.CloseSession(context.Background(), id, model.EndClosed); err != nil {
		t.Fatal(err)
	}
	rec := repo.last(t)
	if len(rec.Sessions) != 1 || rec.Sessions[0].ID != 1 {
		t.Fatalf("Sessions: %#v", rec.Sessions)
	}
	got := rec.Sessions[0]
	if got.ServerParameters["server_version"] != "17.0" || got.Startup["user"] != "app" {
		t.Fatalf("Session-Daten: %#v", got)
	}
	for i, want := range []string{"SELECT 1", "SELECT 2"} {
		in := got.Interactions[i]
		if in.Sequence != i+1 || in.Request.SQL != want || in.Responses[0].Tag != want {
			t.Fatalf("Interaktion %d: %#v", i+1, in)
		}
	}
}

// Abdeckung: LH-FA-02/Boundary — eine Verbindung ohne Anfrage wird nicht
// aufgezeichnet; ein Lauf ohne Session schreibt eine Aufzeichnung ohne Sessions.
func TestRecordSessionOhneAnfrage(t *testing.T) {
	s, repo := neu(t, &fakeUpstream{})
	id := session(t, s)
	if err := s.CloseSession(context.Background(), id, model.EndClosed); err != nil {
		t.Fatal(err)
	}
	if len(repo.writes) != 0 {
		t.Fatalf("Session ohne Anfrage geschrieben: %#v", repo.writes)
	}
	if err := s.Finish(context.Background()); err != nil {
		t.Fatal(err)
	}
	if rec := repo.last(t); len(rec.Sessions) != 0 {
		t.Fatalf("Ende des Laufs: %#v", rec)
	}
}

// Abdeckung: LH-FA-05/Negative, LH-FA-06/Negative — eine Session mit einer nicht unterstützten
// Interaktion wird nicht übernommen, weder bei einer nicht unterstützten
// Client-Nachricht noch bei einer nicht unterstützten Serverantwort.
func TestRecordSessionNichtUnterstuetzt(t *testing.T) {
	t.Run("Client-Nachricht", func(t *testing.T) {
		s, repo := neu(t, &fakeUpstream{})
		id := session(t, s, "SELECT 1")
		if err := s.CloseSession(context.Background(), id, model.EndUnsupported); err != nil {
			t.Fatal(err)
		}
		if len(repo.writes) != 0 {
			t.Fatalf("verworfene Session geschrieben: %#v", repo.writes)
		}
	})
	t.Run("Serverantwort", func(t *testing.T) {
		s, repo := neu(t, &fakeUpstream{queryErr: model.Errorf(model.CodeUnsupported, nil, "COPY")})
		id := session(t, s, "SELECT 1")
		if _, err := s.Query(context.Background(), id, "FEHLER"); err == nil {
			t.Fatal("Fehler erwartet")
		}
		if err := s.CloseSession(context.Background(), id, model.EndFailed); err != nil {
			t.Fatal(err)
		}
		if len(repo.writes) != 0 {
			t.Fatalf("verworfene Session geschrieben: %#v", repo.writes)
		}
	})
}

// Erreicht die Antwort der letzten Interaktion den
// Client nicht (EndWriteFailed), entfällt diese Interaktion, die vorherigen
// bleiben, und das Ende ist PGR-E4003. Ein Abbruch des Upstreams vor
// ReadyForQuery nimmt die Interaktion nicht auf.
func TestRecordSessionVerbindungsende(t *testing.T) {
	t.Run("Antwort nicht zugestellt", func(t *testing.T) {
		s, repo := neu(t, &fakeUpstream{})
		id := session(t, s, "SELECT 1")
		if _, err := s.Query(context.Background(), id, "SELECT 2"); err != nil {
			t.Fatal(err)
		}
		if err := s.CloseSession(context.Background(), id, model.EndWriteFailed); codeOf(err) != model.CodeConnectionLost {
			t.Fatalf("Fehler: %v", err)
		}
		got := repo.last(t).Sessions[0].Interactions
		if len(got) != 1 || got[0].Request.SQL != "SELECT 1" {
			t.Fatalf("Interaktionen: %#v", got)
		}
	})
	t.Run("Upstream bricht ab", func(t *testing.T) {
		s, repo := neu(t, &fakeUpstream{queryErr: model.Errorf(model.CodeConnectionLost, nil, "weg")})
		id := session(t, s, "SELECT 1")
		if _, err := s.Query(context.Background(), id, "FEHLER"); err == nil {
			t.Fatal("Fehler erwartet")
		}
		if err := s.CloseSession(context.Background(), id, model.EndFailed); err != nil {
			t.Fatal(err)
		}
		got := repo.last(t).Sessions[0].Interactions
		if len(got) != 1 || got[0].Request.SQL != "SELECT 1" {
			t.Fatalf("Interaktionen: %#v", got)
		}
	})
}

// Gleichzeitige Sessions erhalten lückenlose
// Kennungen in der Reihenfolge ihres Endes; keine Interaktion geht verloren.
func TestRecordGleichzeitigeSessions(t *testing.T) {
	s, repo := neu(t, &fakeUpstream{})
	const n = 20
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx := context.Background()
			id, _, err := s.OpenSession(ctx, nil)
			if err != nil {
				t.Error(err)
				return
			}
			if _, err := s.Query(ctx, id, fmt.Sprintf("SELECT %d", i)); err != nil {
				t.Error(err)
			}
			s.Delivered(ctx, id)
			if err := s.CloseSession(ctx, id, model.EndClosed); err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()
	rec := repo.last(t)
	if len(rec.Sessions) != n {
		t.Fatalf("Sessions: %d", len(rec.Sessions))
	}
	seen := map[string]bool{}
	for i, sess := range rec.Sessions {
		if sess.ID != i+1 {
			t.Fatalf("Kennung %d an Stelle %d", sess.ID, i+1)
		}
		seen[sess.Interactions[0].Request.SQL] = true
	}
	if len(seen) != n {
		t.Fatalf("verschiedene Anfragen: %d", len(seen))
	}
}
