package services

import (
	"context"
	"fmt"
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
// Verbindung keine Anfrage gestellt hat. pos ist die nächste erwartete
// Interaktion, gruppe und nachricht innerhalb einer Extended-Interaktion die
// nächste erwartete Gruppe und Client-Nachricht (ARC-002, LH-FA-18.a).
type cursor struct {
	session   *model.Session
	pos       int
	gruppe    int
	nachricht int
}

// NewReplayService lädt die Aufzeichnung. Jede Interaktion muss Validate
// bestehen, sonst ist die Aufzeichnung beschädigt (PGR-E3003); auf diese Form
// verlassen sich Query und ClientMessage. Eine Aufzeichnung ohne Session mit
// Interaktion ist PGR-E3004 (LH-FA-03.a).
func NewReplayService(ctx context.Context, repo driven.RecordingRepository, path string) (*ReplayService, error) {
	rec, err := repo.Load(ctx, path)
	if err != nil {
		return nil, err
	}
	s := &ReplayService{verbindung: map[model.SessionID]*cursor{}}
	for _, sess := range rec.Sessions {
		for _, in := range sess.Interactions {
			if err := in.Validate(); err != nil {
				return nil, model.Errorf(model.CodeRecordingBroken, err, "%s: Session %d, Interaktion %d", path, sess.ID, in.Sequence)
			}
		}
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
// Antworten; der Cursor rückt nur bei Gleichheit vor. Erwartet der Cursor eine
// Extended-Nachricht, ist jede Anfrage eine Abweichung (LH-FA-18.a).
func (s *ReplayService) Query(_ context.Context, id model.SessionID, sql string) ([]model.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, err := s.zuordnen(id, fmt.Sprintf("Anfrage %q", sql))
	if err != nil {
		return nil, err
	}
	if c.pos >= len(c.session.Interactions) {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d: keine aufgezeichnete Interaktion mehr nach %d; empfangen %q",
			c.session.ID, c.pos, sql)
	}
	erwartet := c.session.Interactions[c.pos]
	if erwartet.Request.Type != model.RequestQuery {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "%s: erwartet Client-Nachricht %s, empfangen Anfrage %q",
			c.stelle(), c.erwarteteNachricht().Type, sql)
	}
	if erwartet.Request.SQL != sql {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d, Interaktion %d: erwartet %q, empfangen %q",
			c.session.ID, erwartet.Sequence, erwartet.Request.SQL, sql)
	}
	c.pos++
	return erwartet.Responses, nil
}

// ClientMessage vergleicht eine Client-Nachricht des Extended Query Protocol
// mit der erwarteten am Cursor (LH-FA-18.a). Bei Gleichheit rückt der Cursor
// um eine Nachricht vor; schließt sie ihre Gruppe ab, liefert ClientMessage die
// aufgezeichneten Server-Nachrichten der Gruppe, sonst keine. Die Abweichung
// nennt das abweichende Feld und das SQL der betroffenen Anweisungen
// (anweisungen), Parameterwerte nennt sie nicht (SPEC-033); der Cursor bleibt
// dann stehen.
func (s *ReplayService) ClientMessage(_ context.Context, id model.SessionID, m model.ClientMessage) ([]model.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, err := s.zuordnen(id, fmt.Sprintf("Client-Nachricht %s", m.Type))
	if err != nil {
		return nil, err
	}
	if c.pos >= len(c.session.Interactions) {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d: keine aufgezeichnete Interaktion mehr nach %d; empfangen Client-Nachricht %s",
			c.session.ID, c.pos, m.Type)
	}
	erwartet := c.session.Interactions[c.pos]
	if erwartet.Request.Type != model.RequestExtended {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d, Interaktion %d: erwartet Anfrage %q, empfangen Client-Nachricht %s",
			c.session.ID, erwartet.Sequence, erwartet.Request.SQL, m.Type)
	}
	e := c.erwarteteNachricht()
	switch feld := abweichung(m, e); feld {
	case "":
	case "sql":
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "%s: %s weicht in sql ab, erwartet %q, empfangen %q",
			c.stelle(), e.Type, e.SQL, m.SQL)
	default:
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "%s: erwartet %s, empfangen %s, abweichend in %s%s",
			c.stelle(), e.Type, m.Type, feld, c.anweisungen(e, m))
	}
	g := erwartet.Groups[c.gruppe]
	c.nachricht++
	if c.nachricht < len(g.Client) {
		return nil, nil
	}
	c.nachricht = 0
	c.gruppe++
	if c.gruppe == len(erwartet.Groups) {
		c.gruppe = 0
		c.pos++
	}
	return g.Server, nil
}

// zuordnen liefert den Cursor der Verbindung und ordnet ihr bei der ersten
// Anfrage die nächste freie Session zu; ist keine mehr frei, ist das PGR-E5003.
func (s *ReplayService) zuordnen(id model.SessionID, empfangen string) (*cursor, error) {
	c, ok := s.verbindung[id]
	if !ok {
		return nil, model.Errorf(model.CodeInternal, nil, "unbekannte Verbindung %d", id)
	}
	if c.session == nil {
		if len(s.frei) == 0 {
			return nil, model.Errorf(model.CodeReplaySession, nil, "Verbindung %d: keine aufgezeichnete Session mehr frei für %s", id, empfangen)
		}
		sess := s.frei[0]
		s.frei = s.frei[1:]
		c.session = &sess
	}
	return c, nil
}

// anweisungen nennt für die Diagnose das SQL der Anweisungen, auf die sich die
// erwartete und die empfangene Nachricht beziehen (LH-FA-10.a): bei parse ihr
// SQL, bei bind das des Statements, bei execute das des Statements hinter dem
// Portal, bei describe und close das ihres Ziels. Gesucht wird in den
// aufgezeichneten Client-Nachrichten der Session vor dem Cursor; ist keine
// Anweisung bekannt, steht dort „unbekannt“, bezieht sich eine Nachricht auf
// keine, „keine“. Bezieht sich keine der beiden auf eine Anweisung, liefert es
// "".
func (c *cursor) anweisungen(e, m model.ClientMessage) string {
	vorher := c.vorher()
	se, okE := anweisung(vorher, e)
	sm, okM := anweisung(vorher, m)
	if !okE && !okM {
		return ""
	}
	if !okE {
		se = "keine"
	}
	if !okM {
		sm = "keine"
	}
	return fmt.Sprintf(", Anweisung erwartet %s, empfangen %s", se, sm)
}

// vorher sind die aufgezeichneten Client-Nachrichten der Session vor dem
// Cursor in Sendereihenfolge.
func (c *cursor) vorher() []model.ClientMessage {
	var out []model.ClientMessage
	for pi, in := range c.session.Interactions[:c.pos+1] {
		for gi, g := range in.Groups {
			for ni, m := range g.Client {
				if pi == c.pos && (gi > c.gruppe || gi == c.gruppe && ni >= c.nachricht) {
					return out
				}
				out = append(out, m)
			}
		}
	}
	return out
}

// anweisung liefert das SQL der Anweisung, auf die sich m bezieht, gequotet
// oder „unbekannt“; ok ist falsch, wenn sich m auf keine Anweisung bezieht.
func anweisung(vorher []model.ClientMessage, m model.ClientMessage) (string, bool) {
	switch {
	case m.Type == model.ClientParse:
		return fmt.Sprintf("%q", m.SQL), true
	case m.Type == model.ClientBind:
		return statementSQL(vorher, len(vorher), m.Statement), true
	case m.Type == model.ClientExecute:
		return portalSQL(vorher, m.Portal), true
	case (m.Type == model.ClientDescribe || m.Type == model.ClientClose) && m.Target == model.TargetStatement:
		return statementSQL(vorher, len(vorher), m.Name), true
	case m.Type == model.ClientDescribe || m.Type == model.ClientClose:
		return portalSQL(vorher, m.Name), true
	}
	return "", false
}

// statementSQL sucht das letzte parse des Statements name vor bis.
func statementSQL(vorher []model.ClientMessage, bis int, name string) string {
	for i := bis - 1; i >= 0; i-- {
		if vorher[i].Type == model.ClientParse && vorher[i].Statement == name {
			return fmt.Sprintf("%q", vorher[i].SQL)
		}
	}
	return "unbekannt"
}

// portalSQL sucht das letzte bind des Portals und das Statement dahinter.
func portalSQL(vorher []model.ClientMessage, portal string) string {
	for i := len(vorher) - 1; i >= 0; i-- {
		if vorher[i].Type == model.ClientBind && vorher[i].Portal == portal {
			return statementSQL(vorher, i, vorher[i].Statement)
		}
	}
	return "unbekannt"
}

// erwarteteNachricht ist die Client-Nachricht am Cursor; die Interaktion am
// Cursor ist eine Extended-Interaktion.
func (c *cursor) erwarteteNachricht() model.ClientMessage {
	return c.session.Interactions[c.pos].Groups[c.gruppe].Client[c.nachricht]
}

// stelle nennt Session, Interaktion, Gruppe und Nachricht am Cursor, Gruppe und
// Nachricht ab 1 gezählt.
func (c *cursor) stelle() string {
	return fmt.Sprintf("Session %d, Interaktion %d, Gruppe %d, Nachricht %d",
		c.session.ID, c.session.Interactions[c.pos].Sequence, c.gruppe+1, c.nachricht+1)
}

// Shutdown liefert true, wenn die Verbindung beim Herunterfahren enden darf:
// Der Cursor steht nicht innerhalb einer Extended-Interaktion (LH-FA-13.a).
func (s *ReplayService) Shutdown(_ context.Context, id model.SessionID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.verbindung[id]
	return c == nil || c.gruppe == 0 && c.nachricht == 0
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
