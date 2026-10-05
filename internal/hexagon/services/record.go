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

// laufend ist eine offene Session. Ihre Felder schützt mu; die beiden
// Richtungen des Adapters rufen gleichzeitig. Kein Aufruf hält mu über einem
// blockierenden Aufruf des Upstreams; cond weckt Wartende bei jeder
// Zustandsänderung.
type laufend struct {
	upstream driven.UpstreamSession

	mu          sync.Mutex
	cond        *sync.Cond
	session     model.Session
	unsupported bool
	// offen sind die laufenden Extended-Interaktionen in Ankunftsreihenfolge:
	// die erste nimmt Server-Nachrichten auf, die letzte Client-Nachrichten,
	// solange ihr Sync aussteht (LH-FA-18.a).
	offen []*offeneInteraktion
	// einfach: eine einfache Anfrage läuft am Upstream.
	einfach bool
	// empfaengt: AwaitServer liest vom Upstream, oder seine Nachrichten sind
	// noch nicht als zugestellt gemeldet.
	empfaengt bool
	// unzugestellt zählt die abgeschlossenen Interaktionen am Ende der Session,
	// deren Antworten noch nicht als zugestellt gemeldet sind.
	unzugestellt   int
	herunterfahren bool
	beendet        bool
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
	// gesendet: mindestens eine Gruppe ist an den Upstream gegangen.
	gesendet bool
	// sync: das Sync der Interaktion ist an den Upstream gegangen.
	sync bool
}

// laeuft sagt, ob eine Interaktion läuft: eine einfache Anfrage am Upstream
// oder eine Extended-Interaktion vor ihrem ReadyForQuery. Aufrufer hält mu.
func (l *laufend) laeuft() bool {
	return l.einfach || len(l.offen) > 0
}

// beendetErr liefert ErrSessionEnded, wenn die Session beendet ist, sonst
// err. Aufrufer hält mu.
func (l *laufend) beendetErr(err error) error {
	if l.beendet {
		return model.ErrSessionEnded
	}
	return err
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
	l := &laufend{upstream: up, session: session}
	l.cond = sync.NewCond(&l.mu)
	s.sessions[s.next] = l
	return s.next, responses, nil
}

// Query leitet die Anfrage weiter; eine Interaktion wird mit dem ReadyForQuery
// des Servers Teil der Session (LH-FA-02.b). Laufen Extended-Interaktionen mit
// gesendetem Sync, wartet Query, bis ihre Antworten zugestellt sind. Die
// Gegenrichtung empfängt darum nie gleichzeitig mit einer einfachen Anfrage:
// AwaitServer braucht eine laufende Extended-Interaktion, und eine neue beginnt
// nur über ClientMessage derselben Richtung, die gerade in Query steht. Eine
// nicht unterstützte Serverantwort, eine Anfrage während einer
// Extended-Interaktion vor deren Sync und eine Interaktion, die Validate
// ablehnt, sind PGR-E6001 und markieren die Session als nicht übernehmbar
// (LH-FA-18.a); die Antworten gehen dann nicht zurück.
func (s *RecordService) Query(ctx context.Context, id model.SessionID, sql string) ([]model.Response, error) {
	l, err := s.laufende(id)
	if err != nil {
		return nil, err
	}
	l.mu.Lock()
	if n := len(l.offen); n > 0 && !l.offen[n-1].sync {
		l.unsupported = true
		l.mu.Unlock()
		return nil, model.Errorf(model.CodeUnsupported, nil, "einfache Anfrage während einer laufenden Extended-Interaktion")
	}
	for !l.beendet && (len(l.offen) > 0 || l.empfaengt) {
		l.cond.Wait()
	}
	if l.beendet {
		l.mu.Unlock()
		return nil, model.ErrSessionEnded
	}
	l.einfach = true
	l.mu.Unlock()

	responses, err := l.upstream.Query(ctx, sql)

	l.mu.Lock()
	defer l.mu.Unlock()
	defer l.cond.Broadcast()
	l.einfach = false
	if err != nil {
		var me *model.Error
		if errors.As(err, &me) && me.Code == model.CodeUnsupported {
			l.unsupported = true
		}
		return responses, l.beendetErr(err)
	}
	if l.beendet {
		return nil, model.ErrSessionEnded
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

// ClientMessage nimmt die Nachricht in die Extended-Interaktion auf, die
// Client-Nachrichten aufnimmt; nach einem Sync beginnt die nächste Nachricht
// eine neue Interaktion. Die erste Nachricht nach einem Flush beginnt eine neue
// Gruppe, und ab ihr gehen eintreffende Server-Nachrichten an diese Gruppe
// (LH-FA-18.a). Mit Flush oder Sync sendet es die Client-Nachrichten der Gruppe
// an den Upstream, ohne mu zu halten.
func (s *RecordService) ClientMessage(ctx context.Context, id model.SessionID, m model.ClientMessage) error {
	l, err := s.laufende(id)
	if err != nil {
		return err
	}
	l.mu.Lock()
	if l.beendet {
		l.mu.Unlock()
		return model.ErrSessionEnded
	}
	n := len(l.offen)
	if n == 0 || l.offen[n-1].sync {
		l.offen = append(l.offen, &offeneInteraktion{})
		n++
	}
	o := l.offen[n-1]
	if !o.aufbau {
		o.groups = append(o.groups, model.Group{})
		o.ziel = len(o.groups) - 1
		o.aufbau = true
	}
	g := &o.groups[len(o.groups)-1]
	g.Client = append(g.Client, m)
	if m.Type != model.ClientFlush && m.Type != model.ClientSync {
		l.mu.Unlock()
		return nil
	}
	o.aufbau = false
	o.gesendet = true
	o.sync = m.Type == model.ClientSync
	gruppe := append([]model.ClientMessage(nil), g.Client...)
	l.cond.Broadcast()
	l.mu.Unlock()

	if err := l.upstream.Send(ctx, gruppe); err != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		return l.beendetErr(err)
	}
	return nil
}

// AwaitServer wartet, bis die älteste laufende Extended-Interaktion eine
// gesendete Gruppe hat, liest dann ohne mu vom
// Upstream und ordnet jede Nachricht der Gruppe zu, die Antworten aufnimmt. Mit
// ReadyForQuery ist die Interaktion abgeschlossen und wird nach Validate Teil
// der Session; lehnt Validate sie ab, ist das PGR-E6001, und die Session wird
// nicht übernommen. Ein Aufruf, bevor die vorigen Nachrichten als zugestellt
// gemeldet sind, ist ein Fehler des Aufrufers (PGR-E1000).
func (s *RecordService) AwaitServer(ctx context.Context, id model.SessionID) ([]model.Response, error) {
	l, err := s.laufende(id)
	if err != nil {
		return nil, err
	}
	l.mu.Lock()
	if l.empfaengt {
		l.mu.Unlock()
		return nil, model.Errorf(model.CodeInternal, nil, "AwaitServer vor der Zustellung der vorigen Server-Nachrichten")
	}
	for !l.beendet && (len(l.offen) == 0 || !l.offen[0].gesendet) {
		l.cond.Wait()
	}
	if l.beendet {
		l.mu.Unlock()
		return nil, model.ErrSessionEnded
	}
	l.empfaengt = true
	l.mu.Unlock()

	rs, err := l.upstream.Receive(ctx)

	l.mu.Lock()
	defer l.mu.Unlock()
	defer l.cond.Broadcast()
	if err != nil || l.beendet {
		return nil, l.beendetErr(err)
	}
	for _, r := range rs {
		if len(l.offen) == 0 {
			return nil, model.Errorf(model.CodeInternal, nil, "Server-Nachricht %q ohne laufende Interaktion", r.Type)
		}
		o := l.offen[0]
		o.groups[o.ziel].Server = append(o.groups[o.ziel].Server, r)
		if r.Type != model.ResponseReadyForQuery {
			continue
		}
		l.offen = l.offen[1:]
		in := model.Interaction{
			Sequence: len(l.session.Interactions) + 1,
			Request:  model.Request{Type: model.RequestExtended},
			Groups:   o.groups,
		}
		if err := l.uebernehmen(in); err != nil {
			return nil, err
		}
	}
	return rs, nil
}

// Delivered gibt die zuletzt gelieferten Antworten als zugestellt frei. endet
// ist wahr, wenn die Session herunterfährt und keine Interaktion mehr läuft.
func (s *RecordService) Delivered(_ context.Context, id model.SessionID) bool {
	l, err := s.laufende(id)
	if err != nil {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.empfaengt = false
	l.unzugestellt = 0
	l.cond.Broadcast()
	return l.herunterfahren && !l.laeuft()
}

// Shutdown merkt das Herunterfahren vor. Die Session darf sofort enden, wenn
// keine Interaktion läuft und keine Antwort auf ihre Zustellung wartet; sonst
// endet sie nach dem ReadyForQuery der laufenden Interaktion, auch wenn deren
// Sync noch aussteht (LH-FA-13.a, LH-FA-18.a).
func (s *RecordService) Shutdown(_ context.Context, id model.SessionID) bool {
	l, err := s.laufende(id)
	if err != nil {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.herunterfahren = true
	return !l.laeuft() && l.unzugestellt == 0
}

// uebernehmen hängt eine abgeschlossene Interaktion an die Session, wenn
// Validate sie annimmt; sonst ist die Session nicht übernehmbar (PGR-E6001).
// Aufrufer hält mu.
func (l *laufend) uebernehmen(in model.Interaction) error {
	if err := in.Validate(); err != nil {
		l.unsupported = true
		return model.Errorf(model.CodeUnsupported, err, "Interaktion %d nicht in der Form der Aufzeichnung", in.Sequence)
	}
	l.session.Interactions = append(l.session.Interactions, in)
	l.unzugestellt++
	return nil
}

// CloseSession beendet die Session: jeder wartende und jeder spätere Aufruf
// liefert ErrSessionEnded, der Upstream wird geschlossen. Aus dem Ereignis
// folgt (LH-FA-02.b, LH-FA-18.a):
//
//   - EndClosed, EndTerminate: läuft eine Interaktion, ist das PGR-E4003.
//   - EndWriteFailed: PGR-E4003; die Interaktionen, deren Antworten nicht als
//     zugestellt gemeldet sind, entfallen.
//   - EndUnsupported: die Session wird nicht übernommen.
//   - EndShutdown, EndFailed: kein weiterer Fehler.
//
// Eine laufende Interaktion steht nie in der Session. Übernommen wird die
// Session, wenn sie danach mindestens eine Interaktion trägt und keine nicht
// unterstützte enthielt; die Kennung zählt lückenlos ab 1 (LH-FA-12.a).
func (s *RecordService) CloseSession(ctx context.Context, id model.SessionID, end model.SessionEnd) error {
	l, err := s.laufende(id)
	if err != nil {
		return nil
	}
	l.mu.Lock()
	if l.beendet {
		l.mu.Unlock()
		return nil
	}
	l.beendet = true
	l.cond.Broadcast()
	var verbindung error
	interactions := l.session.Interactions
	switch end {
	case model.EndClosed, model.EndTerminate:
		if l.laeuft() {
			verbindung = model.Errorf(model.CodeConnectionLost, nil, "Client-Verbindung vor dem ReadyForQuery der laufenden Interaktion beendet")
		}
	case model.EndWriteFailed:
		interactions = interactions[:len(interactions)-l.unzugestellt]
		verbindung = model.Errorf(model.CodeConnectionLost, nil, "Antwort nicht an den Client zu senden")
	}
	verwerfen := end == model.EndUnsupported || l.unsupported || len(interactions) == 0
	session := l.session
	session.Interactions = interactions
	l.mu.Unlock()

	upErr := l.upstream.Close()

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	if verwerfen {
		return errors.Join(verbindung, upErr)
	}
	session.ID = len(s.rec.Sessions) + 1
	s.rec.Sessions = append(s.rec.Sessions, session)
	return errors.Join(verbindung, s.repo.Write(ctx, s.path, s.rec), upErr)
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
		// Kennungen vergibt nur OpenSession; eine unbekannte ist beendet.
		return nil, model.ErrSessionEnded
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
