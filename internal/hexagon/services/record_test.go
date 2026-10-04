package services

import (
	"context"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

type fakeUpstream struct{ sessions int }

func (f *fakeUpstream) Open(context.Context, map[string]string) (driven.UpstreamSession, []model.Response, error) {
	f.sessions++
	return fakeSession{}, []model.Response{
		{Type: model.ResponseParameterStatus, Name: "server_version", Value: "17.0"},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}, nil
}

type fakeSession struct{}

func (fakeSession) Query(_ context.Context, sql string) ([]model.Response, error) {
	return []model.Response{
		{Type: model.ResponseCommandComplete, Tag: "SELECT 1"},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}, nil
}
func (fakeSession) Close() error { return nil }

type fakeRepo struct{ writes []model.Recording }

func (f *fakeRepo) Prepare(context.Context, string, bool) error { return nil }
func (f *fakeRepo) Write(_ context.Context, _ string, rec model.Recording) error {
	f.writes = append(f.writes, rec)
	return nil
}
func (f *fakeRepo) Load(context.Context, string) (model.Recording, error) {
	return model.Recording{}, nil
}

// LH-FA-02, LH-FA-06: Die Interaktion steht geordnet in der Session; die Session
// wird beim Ende übernommen und die Aufzeichnung geschrieben.
func TestRecordSessionMitInteraktion(t *testing.T) {
	ctx := context.Background()
	repo := &fakeRepo{}
	s, err := NewRecordService(ctx, &fakeUpstream{}, repo, "rec.yaml", false)
	if err != nil {
		t.Fatal(err)
	}
	id, startup, err := s.OpenSession(ctx, map[string]string{"user": "app"})
	if err != nil || id == 0 {
		t.Fatalf("OpenSession: id=%d err=%v", id, err)
	}
	if len(startup) != 2 {
		t.Fatalf("Antworten des Verbindungsaufbaus: %d", len(startup))
	}
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if err := s.CloseSession(ctx, id); err != nil {
		t.Fatal(err)
	}
	if len(repo.writes) != 1 {
		t.Fatalf("Schreibvorgänge: %d", len(repo.writes))
	}
	rec := repo.writes[0]
	if len(rec.Sessions) != 1 || rec.Sessions[0].ID != 1 {
		t.Fatalf("Sessions: %#v", rec.Sessions)
	}
	got := rec.Sessions[0]
	if got.ServerParameters["server_version"] != "17.0" || got.Startup["user"] != "app" {
		t.Fatalf("Session-Daten: %#v", got)
	}
	if len(got.Interactions) != 1 || got.Interactions[0].Sequence != 1 || got.Interactions[0].Request.SQL != "SELECT 1" {
		t.Fatalf("Interaktionen: %#v", got.Interactions)
	}
}

// LH-FA-07.a: Eine Session ohne Anfrage wird nicht aufgezeichnet.
func TestRecordSessionOhneAnfrage(t *testing.T) {
	ctx := context.Background()
	repo := &fakeRepo{}
	s, _ := NewRecordService(ctx, &fakeUpstream{}, repo, "rec.yaml", false)
	id, _, _ := s.OpenSession(ctx, nil)
	if err := s.CloseSession(ctx, id); err != nil {
		t.Fatal(err)
	}
	if len(repo.writes) != 0 {
		t.Fatalf("Session ohne Anfrage geschrieben: %#v", repo.writes)
	}
	if err := s.Finish(ctx); err != nil {
		t.Fatal(err)
	}
	if len(repo.writes) != 1 || len(repo.writes[0].Sessions) != 0 {
		t.Fatalf("Ende des Laufs: %#v", repo.writes)
	}
}
