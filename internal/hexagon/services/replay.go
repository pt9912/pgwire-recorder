package services

import (
	"context"
	"sort"
	"sync"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// ReplayService erfüllt den Replay-Use-Case: Er beantwortet Anfragen aus einer
// Aufzeichnung mit strict sequential matching (LH-FA-09.a). Die Sessions werden
// den Verbindungen in der Reihenfolge ihrer ersten Anfrage zugeordnet
// (`first-request`, LH-FA-12.a).
type ReplayService struct {
	mu sync.Mutex
	// frei sind die Sessions mit Interaktionen, die noch keiner Verbindung
	// zugeordnet sind, in der Reihenfolge ihrer Kennung.
	frei       []model.Session
	letzte     model.Session
	next       model.SessionID
	verbindung map[model.SessionID]*cursor
}

// cursor ist der Replay-Cursor einer Verbindung; session ist nil, solange die
// Verbindung keine Anfrage gestellt hat.
type cursor struct {
	session *model.Session
	pos     int
}

// NewReplayService lädt die Aufzeichnung. Eine Aufzeichnung ohne Session mit
// Interaktion ist PGR-E3004 (LH-FA-03.a).
func NewReplayService(ctx context.Context, repo driven.RecordingRepository, path string) (*ReplayService, error) {
	rec, err := repo.Load(ctx, path)
	if err != nil {
		return nil, err
	}
	s := &ReplayService{verbindung: map[model.SessionID]*cursor{}}
	for _, sess := range rec.Sessions {
		if len(sess.Interactions) > 0 {
			s.frei = append(s.frei, sess)
		}
	}
	if len(s.frei) == 0 {
		return nil, model.Errorf(model.CodeRecordingNoSession, nil, "%s enthält keine Session mit Interaktion", path)
	}
	s.letzte = s.frei[len(s.frei)-1]
	return s, nil
}

// OpenConnection liefert den Handshake mit den Serverparametern der nächsten
// noch nicht zugeordneten Session, gibt es keine mehr, der letzten.
func (s *ReplayService) OpenConnection(context.Context) (model.SessionID, []model.Response) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	s.verbindung[s.next] = &cursor{}
	vorlage := s.letzte
	if len(s.frei) > 0 {
		vorlage = s.frei[0]
	}
	namen := make([]string, 0, len(vorlage.ServerParameters))
	for n := range vorlage.ServerParameters {
		namen = append(namen, n)
	}
	sort.Strings(namen)
	out := make([]model.Response, 0, len(namen)+1)
	for _, n := range namen {
		out = append(out, model.Response{Type: model.ResponseParameterStatus, Name: n, Value: vorlage.ServerParameters[n]})
	}
	out = append(out, model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"})
	return s.next, out
}

// Query vergleicht die Anfrage mit der Interaktion am Cursor und liefert deren
// Antworten; der Cursor rückt nur bei Gleichheit vor.
func (s *ReplayService) Query(_ context.Context, id model.SessionID, sql string) ([]model.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.verbindung[id]
	if !ok {
		return nil, model.Errorf(model.CodeInternal, nil, "unbekannte Verbindung %d", id)
	}
	if c.session == nil {
		if len(s.frei) == 0 {
			return nil, model.Errorf(model.CodeReplaySession, nil, "Verbindung %d: keine aufgezeichnete Session mehr frei für %q", id, sql)
		}
		sess := s.frei[0]
		s.frei = s.frei[1:]
		c.session = &sess
	}
	if c.pos >= len(c.session.Interactions) {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d: keine aufgezeichnete Interaktion mehr nach %d; empfangen %q",
			c.session.ID, c.pos, sql)
	}
	erwartet := c.session.Interactions[c.pos]
	if erwartet.Request.Type != model.RequestQuery || erwartet.Request.SQL != sql {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d, Interaktion %d: erwartet %q, empfangen %q",
			c.session.ID, erwartet.Sequence, erwartet.Request.SQL, sql)
	}
	c.pos++
	return erwartet.Responses, nil
}

// CloseConnection beendet die Verbindung; unverbrauchte Interaktionen der
// zugeordneten Session sind die Warnung PGR-W2001 (LH-FA-03.b).
func (s *ReplayService) CloseConnection(_ context.Context, id model.SessionID) *model.Warning {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.verbindung[id]
	delete(s.verbindung, id)
	if c == nil || c.session == nil || c.pos >= len(c.session.Interactions) {
		return nil
	}
	return model.Warnf(model.CodeUnconsumed, "Session %d: %d von %d Interaktionen nicht verbraucht",
		c.session.ID, len(c.session.Interactions)-c.pos, len(c.session.Interactions))
}

// Unassigned liefert die Warnung PGR-W2001 für Sessions mit Interaktionen, die
// bis zum Ende des Laufs keiner Verbindung zugeordnet wurden, oder nil.
func (s *ReplayService) Unassigned() *model.Warning {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.frei) == 0 {
		return nil
	}
	return model.Warnf(model.CodeUnconsumed, "%d aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung %d",
		len(s.frei), s.frei[0].ID)
}
