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
// (`first-request`, LH-FA-12.a). Lebendprüfungen beantwortet er außerhalb der
// Reihe, aufgezeichnete überspringt er (ADR-0031).
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
// nächste erwartete Gruppe und Client-Nachricht (ARC-002, LH-FA-18.a). status
// ist der Transaktionsstatus des letzten ReadyForQuery, das die Verbindung
// erhalten hat, nach dem Handshake "I" (LH-FA-09.a).
type cursor struct {
	session   *model.Session
	pos       int
	gruppe    int
	nachricht int
	status    string
}

// NewReplayService lädt die Aufzeichnung. Jede Interaktion muss Validate
// bestehen, sonst ist die Aufzeichnung beschädigt (PGR-E3003); auf diese Form
// verlassen sich Query und ClientMessage. Interaktionen, deren Anfrage eine
// Lebendprüfung ist, nimmt es nicht in die Sessions auf; eine Session ohne
// andere Interaktion ist damit eine Session ohne Interaktion (LH-FA-09.a). Eine
// Aufzeichnung ohne Session mit Interaktion ist PGR-E3004 (LH-FA-03.a).
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
		sess.Interactions = ohneLebendpruefungen(sess.Interactions)
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

// ohneLebendpruefungen liefert die Interaktionen ohne die einfachen Anfragen,
// die Lebendprüfungen sind, in einer neuen Liste; Sequence bleibt die
// aufgezeichnete Nummer.
func ohneLebendpruefungen(ins []model.Interaction) []model.Interaction {
	out := make([]model.Interaction, 0, len(ins))
	for _, in := range ins {
		if in.Request.Type == model.RequestQuery && istLebendpruefung(in.Request.SQL) {
			continue
		}
		out = append(out, in)
	}
	return out
}

// OpenConnection liefert den Handshake mit den Serverparametern der nächsten
// noch nicht zugeordneten Session, gibt es keine mehr, der letzten.
func (s *ReplayService) OpenConnection(context.Context) (model.SessionID, []model.Response) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	s.verbindung[s.next] = &cursor{status: "I"}
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
// Antworten; der Cursor rückt nur bei Gleichheit vor. Eine Lebendprüfung
// zwischen zwei Interaktionen beantwortet es mit EmptyQueryResponse und
// ReadyForQuery im Status des letzten ReadyForQuery, ohne Cursor und Zuordnung zu
// berühren (LH-FA-09.a). Erwartet der Cursor eine Extended-Nachricht, ist jede
// andere Anfrage eine Abweichung, mitten in einer Extended-Interaktion auch
// eine Lebendprüfung (LH-FA-18.a).
func (s *ReplayService) Query(_ context.Context, id model.SessionID, sql string) ([]model.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.verbindung[id]; ok && !c.mitten() && istLebendpruefung(sql) {
		return []model.Response{
			{Type: model.ResponseEmptyQueryResponse},
			{Type: model.ResponseReadyForQuery, TxStatus: c.status},
		}, nil
	}
	c, err := s.zuordnen(id, fmt.Sprintf("Anfrage %q", sql))
	if err != nil {
		return nil, err
	}
	if c.pos >= len(c.session.Interactions) {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d: keine aufgezeichnete Interaktion mehr nach %d; empfangen %q",
			c.session.ID, c.letzteNummer(), sql)
	}
	erwartet := c.session.Interactions[c.pos]
	if erwartet.Request.Type != model.RequestQuery {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "%s: erwartet %s, empfangen Anfrage %q",
			c.stelle(), c.nachrichtText(c.erwarteteNachricht()), sql)
	}
	if erwartet.Request.SQL != sql {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d, Interaktion %d: erwartet %q, empfangen %q",
			c.session.ID, erwartet.Sequence, erwartet.Request.SQL, sql)
	}
	c.pos++
	c.merke(erwartet.Responses)
	return erwartet.Responses, nil
}

// mitten liefert true, wenn der Cursor innerhalb einer Extended-Interaktion
// steht: nach ihrer ersten Nachricht und vor ihrem Sync (LH-FA-18.a).
func (c *cursor) mitten() bool {
	return c.gruppe > 0 || c.nachricht > 0
}

// merke hält den Transaktionsstatus des letzten ReadyForQuery in rs fest.
func (c *cursor) merke(rs []model.Response) {
	for _, r := range rs {
		if r.Type == model.ResponseReadyForQuery {
			c.status = r.TxStatus
		}
	}
}

// letzteNummer ist die aufgezeichnete Nummer der letzten Interaktion der
// Session; die Session hat mindestens eine.
func (c *cursor) letzteNummer() int {
	return c.session.Interactions[len(c.session.Interactions)-1].Sequence
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
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d: keine aufgezeichnete Interaktion mehr nach %d; empfangen %s",
			c.session.ID, c.letzteNummer(), c.nachrichtText(m))
	}
	erwartet := c.session.Interactions[c.pos]
	if erwartet.Request.Type != model.RequestExtended {
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d, Interaktion %d: erwartet Anfrage %q, empfangen %s",
			c.session.ID, erwartet.Sequence, erwartet.Request.SQL, c.nachrichtText(m))
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
	c.merke(g.Server)
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

// nachrichtText nennt eine Client-Nachricht für die Diagnose einer
// Abweichung der Art oder nach dem Ende der Aufzeichnung: ihren Typ und, wenn
// sie sich auf eine Anweisung bezieht, deren SQL wie in anweisungen
// (LH-FA-10.a), etwa `Client-Nachricht parse (Anweisung "SELECT 1")`.
func (c *cursor) nachrichtText(m model.ClientMessage) string {
	if sql, ok := c.objekte().anweisung(m); ok {
		return fmt.Sprintf("Client-Nachricht %s (Anweisung %s)", m.Type, sql)
	}
	return fmt.Sprintf("Client-Nachricht %s", m.Type)
}

// anweisungen nennt für die Diagnose das SQL der Anweisungen, auf die sich die
// erwartete und die empfangene Nachricht beziehen (LH-FA-10.a): bei parse ihr
// SQL, bei bind das des Statements, bei execute das des Portals, bei describe
// und close das ihres Ziels. Was ein Statement oder Portal vor dem Cursor
// bedeutet, ergibt der Verlauf der Aufzeichnung (objekte); ein dort nicht
// bestehendes Objekt ist „unbekannt“. Bezieht sich eine Nachricht auf keine
// Anweisung, steht „keine“; gilt das für beide, liefert es "".
func (c *cursor) anweisungen(e, m model.ClientMessage) string {
	o := c.objekte()
	se, okE := o.anweisung(e)
	sm, okM := o.anweisung(m)
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

// objekte sind die Statements und Portale einer Session mit dem SQL ihrer
// Anweisung; ein Portal trägt das SQL seines Statements zum Zeitpunkt des
// bind.
type objekte struct {
	statements map[string]string
	portale    map[string]string
}

// objekte spielt die aufgezeichneten Protokollnachrichten der Session vor dem
// Cursor nach, soweit der Server sie bestätigt hat: parse legt ein Statement an
// oder überschreibt das unbenannte, bind ein Portal, close entfernt sein Ziel;
// ein abgelehntes oder bis zum Sync verworfenes parse oder bind legt nichts an.
// Eine einfache Anfrage entfernt das unbenannte Statement und das unbenannte
// Portal. Endet eine Interaktion mit ready_for_query im Status "I" (keine
// Transaktion), enden alle Portale. SQL-Befehle, die Objekte beenden
// (DEALLOCATE, DISCARD ALL, CLOSE, ROLLBACK TO SAVEPOINT), bildet es nicht
// nach; danach kann die Diagnose ein beendetes Objekt als bestehend nennen.
// Lebendprüfungen stehen nicht in der Session (NewReplayService); dass
// PostgreSQL bei ihnen das unbenannte Statement entfernt, bildet es ebenfalls
// nicht nach.
func (c *cursor) objekte() objekte {
	o := objekte{statements: map[string]string{}, portale: map[string]string{}}
	for pi, in := range c.session.Interactions {
		if pi > c.pos || pi == c.pos && in.Request.Type == model.RequestQuery {
			return o
		}
		if in.Request.Type == model.RequestQuery {
			delete(o.statements, "")
			delete(o.portale, "")
			o.transaktionsende(in.Responses)
			continue
		}
		// Der Server bestätigt jedes angenommene parse, bind und close; nach
		// einer Ablehnung verwirft er bis zum Sync. Die ersten n Nachrichten
		// einer Art mit n Bestätigungen in der Interaktion gelten als
		// angenommen; gezählt wird je Interaktion, weil eine späte Bestätigung
		// in der folgenden Gruppe steht (LH-FA-18.a).
		bestaetigt := map[model.ClientMessageType]int{}
		for _, g := range in.Groups {
			for _, r := range g.Server {
				bestaetigt[bestaetigung[r.Type]]++
			}
		}
		for gi, g := range in.Groups {
			for ni, m := range g.Client {
				if pi == c.pos && (gi > c.gruppe || gi == c.gruppe && ni >= c.nachricht) {
					return o
				}
				if bestaetigt[m.Type] > 0 {
					bestaetigt[m.Type]--
					o.nachspielen(m)
				}
			}
		}
		o.transaktionsende(in.Groups[len(in.Groups)-1].Server)
	}
	return o
}

// bestaetigung ordnet jeder Bestätigung des Servers die Client-Nachricht zu,
// die sie bestätigt.
var bestaetigung = map[model.ResponseType]model.ClientMessageType{
	model.ResponseParseComplete: model.ClientParse,
	model.ResponseBindComplete:  model.ClientBind,
	model.ResponseCloseComplete: model.ClientClose,
}

// nachspielen wendet eine vom Server angenommene Nachricht an.
func (o objekte) nachspielen(m model.ClientMessage) {
	switch m.Type {
	case model.ClientParse:
		o.statements[m.Statement] = m.SQL
	case model.ClientBind:
		if sql, ok := o.statements[m.Statement]; ok {
			o.portale[m.Portal] = sql
		} else {
			delete(o.portale, m.Portal)
		}
	case model.ClientClose:
		if m.Target == model.TargetStatement {
			delete(o.statements, m.Name)
		} else {
			delete(o.portale, m.Name)
		}
	}
}

// transaktionsende entfernt alle Portale, wenn die Antworten mit
// ready_for_query im Status "I" enden.
func (o objekte) transaktionsende(rs []model.Response) {
	if n := len(rs); n > 0 && rs[n-1].Type == model.ResponseReadyForQuery && rs[n-1].TxStatus == "I" {
		clear(o.portale)
	}
}

// anweisung liefert das SQL der Anweisung, auf die sich m bezieht, gequotet
// oder „unbekannt“; ok ist falsch, wenn sich m auf keine Anweisung bezieht.
func (o objekte) anweisung(m model.ClientMessage) (string, bool) {
	nachschlagen := func(tabelle map[string]string, name string) (string, bool) {
		if sql, ok := tabelle[name]; ok {
			return fmt.Sprintf("%q", sql), true
		}
		return "unbekannt", true
	}
	switch {
	case m.Type == model.ClientParse:
		return fmt.Sprintf("%q", m.SQL), true
	case m.Type == model.ClientBind:
		return nachschlagen(o.statements, m.Statement)
	case m.Type == model.ClientExecute:
		return nachschlagen(o.portale, m.Portal)
	case (m.Type == model.ClientDescribe || m.Type == model.ClientClose) && m.Target == model.TargetStatement:
		return nachschlagen(o.statements, m.Name)
	case m.Type == model.ClientDescribe || m.Type == model.ClientClose:
		return nachschlagen(o.portale, m.Name)
	}
	return "", false
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
