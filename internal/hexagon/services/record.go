package services

import (
	"context"
	"sync"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// RecordService erfüllt den Record-Use-Case: Es vermittelt Sessions zum
// Upstream und hält die Aufzeichnung. Nach jeder übernommenen Session schreibt
// es die Aufzeichnung als Ganzes (LH-FA-07.a).
type RecordService struct {
	upstream driven.Upstream
	repo     driven.RecordingRepository
	path     string

	mu       sync.Mutex
	rec      model.Recording
	next     model.SessionID
	sessions map[model.SessionID]*laufend
}

type laufend struct {
	upstream driven.UpstreamSession
	session  model.Session
}

// NewRecordService prüft den Zielpfad und liefert den Service.
func NewRecordService(ctx context.Context, upstream driven.Upstream, repo driven.RecordingRepository, path string, replace bool) (*RecordService, error) {
	if err := repo.Prepare(ctx, path, replace); err != nil {
		return nil, err
	}
	return &RecordService{
		upstream: upstream,
		repo:     repo,
		path:     path,
		rec:      model.NewRecording(),
		sessions: map[model.SessionID]*laufend{},
	}, nil
}

// OpenSession baut die Upstream-Session auf.
func (s *RecordService) OpenSession(ctx context.Context, startup map[string]string) (model.SessionID, []model.Response, error) {
	up, responses, err := s.upstream.Open(ctx, startup)
	if err != nil || up == nil {
		return 0, responses, err
	}
	session := model.Session{
		Startup:          kopie(startup),
		ServerParameters: map[string]string{},
	}
	for _, r := range responses {
		if r.Type == model.ResponseParameterStatus {
			session.ServerParameters[r.Name] = r.Value
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	s.sessions[s.next] = &laufend{upstream: up, session: session}
	return s.next, responses, nil
}

// Query leitet die Anfrage weiter; eine Interaktion wird mit dem ReadyForQuery
// des Servers Teil der Session (LH-FA-02.b).
func (s *RecordService) Query(ctx context.Context, id model.SessionID, sql string) ([]model.Response, error) {
	l, err := s.laufende(id)
	if err != nil {
		return nil, err
	}
	responses, err := l.upstream.Query(ctx, sql)
	if err != nil {
		return responses, err
	}
	l.session.Interactions = append(l.session.Interactions, model.Interaction{
		Sequence:  len(l.session.Interactions) + 1,
		Request:   model.Request{Type: model.RequestQuery, SQL: sql},
		Responses: responses,
	})
	return responses, nil
}

// CloseSession beendet die Session und übernimmt sie, wenn sie mindestens eine
// Interaktion trägt; die Kennung in der Aufzeichnung zählt lückenlos ab 1.
func (s *RecordService) CloseSession(ctx context.Context, id model.SessionID) error {
	l, err := s.laufende(id)
	if err != nil {
		return err
	}
	upErr := l.upstream.Close()

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	if len(l.session.Interactions) == 0 {
		return upErr
	}
	l.session.ID = len(s.rec.Sessions) + 1
	s.rec.Sessions = append(s.rec.Sessions, l.session)
	if err := s.repo.Write(ctx, s.path, s.rec); err != nil {
		return err
	}
	return upErr
}

// Finish schreibt die Aufzeichnung beim Ende des Laufs; ein Lauf ohne Session
// schreibt eine gültige Aufzeichnung ohne Sessions.
func (s *RecordService) Finish(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.repo.Write(ctx, s.path, s.rec)
}

func (s *RecordService) laufende(id model.SessionID) (*laufend, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.sessions[id]
	if !ok {
		return nil, model.Errorf(model.CodeInternal, nil, "unbekannte Session %d", id)
	}
	return l, nil
}

func kopie(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
