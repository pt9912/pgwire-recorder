package services

import (
	"context"
	"errors"
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

// laufend ist eine offene Session. Ihre Felder berührt nur die Goroutine der
// zugehörigen Client-Verbindung; die Map schützt mu. AwaitServer liest nur
// upstream.
type laufend struct {
	upstream    driven.UpstreamSession
	session     model.Session
	unsupported bool
	// offen ist die laufende Extended-Interaktion von ihrer ersten
	// Client-Nachricht bis zu ihrem ReadyForQuery; sonst nil.
	offen *offeneInteraktion
}

// offeneInteraktion ist eine Extended-Interaktion im Aufbau (LH-FA-18.a).
type offeneInteraktion struct {
	groups []model.Group
	// aufbau: die letzte Gruppe nimmt noch Client-Nachrichten auf, ihr Flush
	// oder Sync steht aus.
	aufbau bool
	// ziel ist der Index der Gruppe, der eintreffende Server-Nachrichten
	// zugeordnet werden: die zuletzt begonnene. Eine Flush-Gruppe nimmt damit
	// Server-Nachrichten auf, bis die nächste Client-Nachricht eintrifft.
	ziel int
	// sync: das Sync der Interaktion ist an den Upstream gegangen.
	sync bool
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
// des Servers Teil der Session (LH-FA-02.b). Eine nicht unterstützte
// Serverantwort, eine Anfrage während einer laufenden Extended-Interaktion und
// eine Interaktion, die Validate ablehnt, sind PGR-E6001 und markieren die
// Session als nicht übernehmbar (LH-FA-18.a); die Antworten gehen dann nicht
// zurück.
func (s *RecordService) Query(ctx context.Context, id model.SessionID, sql string) ([]model.Response, error) {
	l, err := s.laufende(id)
	if err != nil {
		return nil, err
	}
	if l.offen != nil {
		l.unsupported = true
		return nil, model.Errorf(model.CodeUnsupported, nil, "einfache Anfrage während einer laufenden Extended-Interaktion")
	}
	responses, err := l.upstream.Query(ctx, sql)
	if err != nil {
		var me *model.Error
		if errors.As(err, &me) && me.Code == model.CodeUnsupported {
			l.unsupported = true
		}
		return responses, err
	}
	in := model.Interaction{
		Sequence:  len(l.session.Interactions) + 1,
		Request:   model.Request{Type: model.RequestQuery, SQL: sql},
		Responses: responses,
	}
	if err := l.uebernehmen(in); err != nil {
		return nil, err
	}
	return responses, nil
}

// ClientMessage nimmt die Nachricht in die laufende Extended-Interaktion auf;
// die erste Nachricht nach einem Flush beginnt eine neue Gruppe, und ab ihr
// gehen eintreffende Server-Nachrichten an diese Gruppe (LH-FA-18.a). Mit Flush
// oder Sync sendet es die Client-Nachrichten der Gruppe an den Upstream. Eine
// Nachricht nach dem Sync der Interaktion ist ein Fehler des Aufrufers
// (PGR-E1000).
func (s *RecordService) ClientMessage(ctx context.Context, id model.SessionID, m model.ClientMessage) error {
	l, err := s.laufende(id)
	if err != nil {
		return err
	}
	if l.offen == nil {
		l.offen = &offeneInteraktion{}
	}
	o := l.offen
	if o.sync {
		return model.Errorf(model.CodeInternal, nil, "Client-Nachricht %q nach dem Sync vor dessen ReadyForQuery", m.Type)
	}
	if !o.aufbau {
		o.groups = append(o.groups, model.Group{})
		o.ziel = len(o.groups) - 1
		o.aufbau = true
	}
	g := &o.groups[len(o.groups)-1]
	g.Client = append(g.Client, m)
	if m.Type != model.ClientFlush && m.Type != model.ClientSync {
		return nil
	}
	o.aufbau = false
	o.sync = m.Type == model.ClientSync
	return l.upstream.Send(ctx, g.Client)
}

// AwaitServer liest die nächsten Server-Nachrichten vom Upstream, ohne den
// Zustand der Session zu berühren.
func (s *RecordService) AwaitServer(ctx context.Context, id model.SessionID) ([]model.Response, error) {
	l, err := s.laufende(id)
	if err != nil {
		return nil, err
	}
	return l.upstream.Receive(ctx)
}

// ServerMessage hängt die Nachricht an die Gruppe, die Antworten aufnimmt. Mit
// ReadyForQuery ist die Interaktion abgeschlossen und wird nach Validate Teil
// der Session; lehnt Validate sie ab, ist das PGR-E6001, und die Session wird
// nicht übernommen. Eine Nachricht ohne laufende Interaktion mit gesendeter
// Gruppe ist ein Fehler des Aufrufers (PGR-E1000).
func (s *RecordService) ServerMessage(_ context.Context, id model.SessionID, r model.Response) (bool, error) {
	l, err := s.laufende(id)
	if err != nil {
		return false, err
	}
	o := l.offen
	if o == nil || (len(o.groups) == 1 && o.aufbau) {
		return false, model.Errorf(model.CodeInternal, nil, "Server-Nachricht %q ohne gesendete Gruppe", r.Type)
	}
	o.groups[o.ziel].Server = append(o.groups[o.ziel].Server, r)
	if r.Type != model.ResponseReadyForQuery {
		return false, nil
	}
	l.offen = nil
	in := model.Interaction{
		Sequence: len(l.session.Interactions) + 1,
		Request:  model.Request{Type: model.RequestExtended},
		Groups:   o.groups,
	}
	if err := l.uebernehmen(in); err != nil {
		return false, err
	}
	return true, nil
}

// uebernehmen hängt eine abgeschlossene Interaktion an die Session, wenn
// Validate sie annimmt; sonst ist die Session nicht übernehmbar (PGR-E6001).
func (l *laufend) uebernehmen(in model.Interaction) error {
	if err := in.Validate(); err != nil {
		l.unsupported = true
		return model.Errorf(model.CodeUnsupported, err, "Interaktion %d nicht in der Form der Aufzeichnung", in.Sequence)
	}
	l.session.Interactions = append(l.session.Interactions, in)
	return nil
}

// CloseSession beendet die Session. Übernommen wird sie, wenn sie nach dem
// Grund des Endes mindestens eine Interaktion trägt und keine nicht
// unterstützte Interaktion enthielt; die Kennung in der Aufzeichnung zählt
// lückenlos ab 1 (LH-FA-12.a). Eine laufende Extended-Interaktion ist nicht
// abgeschlossen und steht nie in der Session (LH-FA-18.a).
func (s *RecordService) CloseSession(ctx context.Context, id model.SessionID, end model.SessionEnd) error {
	l, err := s.laufende(id)
	if err != nil {
		return err
	}
	upErr := l.upstream.Close()

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)

	interactions := l.session.Interactions
	if end == model.EndLost && len(interactions) > 0 {
		interactions = interactions[:len(interactions)-1]
	}
	if end == model.EndUnsupported || l.unsupported || len(interactions) == 0 {
		return upErr
	}
	l.session.Interactions = interactions
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
