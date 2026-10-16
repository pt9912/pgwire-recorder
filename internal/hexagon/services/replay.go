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
// den Verbindungen zugeordnet: mit `first-request` (Standard) in der Reihenfolge
// der ersten Anfragen, mit `connection` (ZuordnungNachVerbindung) mit der
// Annahme der n-ten Verbindung die Session mit der Kennung n (LH-FA-12.a).
// Lebendprüfungen beantwortet er außerhalb der Reihe, aufgezeichnete ohne Ende
// überspringt er, eine mit Ende ist eine Interaktion wie jede andere (ADR-0031,
// LH-FA-09.a).
type ReplayService struct {
	// streng meldet nicht verbrauchte Interaktionen und nie zugeordnete
	// Sessions als Fehler PGR-E5002 statt als Warnung PGR-W2001 (LH-FA-03.b).
	streng bool
	// nachVerbindung ordnet mit der Annahme der n-ten Verbindung die Session mit
	// der Kennung n zu (`connection`), sonst bei der ersten Anfrage die nächste
	// freie (`first-request`).
	nachVerbindung bool

	mu sync.Mutex
	// frei sind die Sessions mit Interaktionen in der Reihenfolge ihrer Kennung.
	// Mit `first-request` bleiben die, die noch keiner Verbindung zugeordnet
	// sind; mit `connection` liest `OpenConnection` die Session aus `alle`.
	frei []model.Session
	// alle sind die Sessions der Aufzeichnung nach ihrer Kennung, auch die ohne
	// Interaktion; nur `connection` liest sie.
	alle       map[int]model.Session
	letzte     model.Session
	next       model.SessionID
	verbindung map[model.SessionID]*cursor
}

// cursor ist der Replay-Cursor einer Verbindung; session ist nil, solange die
// Verbindung keiner Session zugeordnet ist: mit `first-request` bis zur ersten
// Anfrage, mit `connection` dauerhaft, wenn es keine Session mit der Kennung der
// Verbindung gibt. pos ist die nächste erwartete
// Interaktion, gruppe und nachricht innerhalb einer Extended-Interaktion die
// nächste erwartete Gruppe und Client-Nachricht (ARC-002, LH-FA-18.a). status
// ist der Transaktionsstatus des letzten ReadyForQuery, das die Verbindung
// erhalten hat, nach dem Handshake "I"; vt sagt, ob \v für Lebendprüfungen als
// Leerraum zählt: nach der Zuordnung nach der Version der Session, vorher nicht
// (LH-FA-09.a). ungesendet ist wahr, solange die Antworten der Interaktion vor
// pos geliefert, aber noch nicht als gesendet gemeldet sind (Sent); diese
// Interaktion ist dann nicht verbraucht (LH-FA-03.b).
type cursor struct {
	session    *model.Session
	pos        int
	gruppe     int
	nachricht  int
	status     string
	vt         bool
	ungesendet bool
}

// ReplayOption stellt einen ReplayService ein.
type ReplayOption func(*ReplayService)

// FailOnUnconsumed meldet nicht verbrauchte Interaktionen und nie zugeordnete
// Sessions als Fehler PGR-E5002 statt als Warnung PGR-W2001 (LH-FA-03.b).
func FailOnUnconsumed(s *ReplayService) { s.streng = true }

// ZuordnungNachVerbindung ordnet mit der Annahme der n-ten Verbindung die
// Session mit der Kennung n zu, auch eine ohne Interaktion (`--session-assignment
// connection`, LH-FA-12.a). Die Aufzeichnung muss `empty_sessions` tragen,
// sonst ist NewReplayService PGR-E2001.
func ZuordnungNachVerbindung(s *ReplayService) { s.nachVerbindung = true }

// NewReplayService lädt die Aufzeichnung. Jede Interaktion muss Validate
// bestehen, sonst ist die Aufzeichnung beschädigt (PGR-E3003); auf diese Form
// verlassen sich Query und ClientMessage. Interaktionen, deren Anfrage nach
// der Version ihrer Session eine Lebendprüfung ohne Ende ist, nimmt es nicht in
// die Sessions auf; eine Session ohne
// andere Interaktion ist damit eine Session ohne Interaktion (LH-FA-09.a). Eine
// Aufzeichnung ohne Session mit Interaktion ist PGR-E3004 (LH-FA-03.a); mit
// ZuordnungNachVerbindung und ohne `empty_sessions` in der Aufzeichnung ist es
// danach PGR-E2001 (LH-FA-12.a).
func NewReplayService(ctx context.Context, repo driven.RecordingRepository, path string, opts ...ReplayOption) (*ReplayService, error) {
	rec, err := repo.Load(ctx, path)
	if err != nil {
		return nil, err
	}
	s := &ReplayService{verbindung: map[model.SessionID]*cursor{}, alle: map[int]model.Session{}}
	for _, o := range opts {
		o(s)
	}
	var letzteAlle model.Session
	for _, sess := range rec.Sessions {
		for _, in := range sess.Interactions {
			if err := in.Validate(); err != nil {
				return nil, model.Errorf(model.CodeRecordingBroken, err, "%s: Session %d, Interaktion %d", path, sess.ID, in.Sequence)
			}
		}
		sess.Interactions = ohneLebendpruefungen(sess.Interactions, vtLeerraum(sess.ServerParameters))
		s.alle[sess.ID] = sess
		letzteAlle = sess
		if len(sess.Interactions) > 0 {
			s.frei = append(s.frei, sess)
		}
	}
	if len(s.frei) == 0 {
		return nil, model.Errorf(model.CodeRecordingNoSession, nil, "%s enthält keine Session mit Interaktion", path)
	}
	s.letzte = s.frei[len(s.frei)-1]
	if s.nachVerbindung {
		if !rec.EmptySessions {
			return nil, model.Errorf(model.CodeUsage, nil, "%s trägt keine Kennzeichnung empty_sessions, die --session-assignment connection verlangt", path)
		}
		s.letzte = letzteAlle
	}
	return s, nil
}

// ohneLebendpruefungen liefert die Interaktionen ohne die einfachen Anfragen,
// die Lebendprüfungen ohne Ende sind, mit vt auch solche mit \v, in einer neuen
// Liste; Sequence bleibt die aufgezeichnete Nummer. Eine Lebendprüfung mit Ende
// bleibt (LH-FA-09.a): Hier entscheidet das Replay, was eine Interaktion ist.
func ohneLebendpruefungen(ins []model.Interaction, vt bool) []model.Interaction {
	out := make([]model.Interaction, 0, len(ins))
	for _, in := range ins {
		if in.Request.Type == model.RequestQuery && istLebendpruefung(in.Request.SQL, vt) && !endetEinfach(in) {
			continue
		}
		out = append(out, in)
	}
	return out
}

// endetEinfach liefert true, wenn die Antworten einer einfachen Anfrage mit
// connection_closed enden (LH-FA-25.a).
func endetEinfach(in model.Interaction) bool {
	n := len(in.Responses)
	return n > 0 && in.Responses[n-1].Type == model.ResponseConnectionClosed
}

// OpenConnection liefert den Handshake mit den Serverparametern der nächsten
// noch nicht zugeordneten Session, gibt es keine mehr, der letzten. Mit
// ZuordnungNachVerbindung ordnet es der n-ten Verbindung die Session mit der
// Kennung n zu und liefert deren Serverparameter; gibt es keine, bleibt die
// Verbindung ohne Session und erhält die der letzten (LH-FA-12.a).
func (s *ReplayService) OpenConnection(context.Context) (model.SessionID, []model.Response) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	c := &cursor{status: "I"}
	s.verbindung[s.next] = c
	vorlage := s.letzte
	if s.nachVerbindung {
		if sess, ok := s.alle[int(s.next)]; ok {
			c.session, c.vt = &sess, vtLeerraum(sess.ServerParameters)
			vorlage = sess
		}
	} else if len(s.frei) > 0 {
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
// (Leerraum nach cursor.vt) zwischen zwei Interaktionen beantwortet es mit EmptyQueryResponse und
// ReadyForQuery im Status des letzten ReadyForQuery, ohne Cursor und Zuordnung zu
// berühren (LH-FA-09.a); steht am Cursor, bei `first-request` für eine noch nicht
// zugeordnete Verbindung am Anfang der nächsten freien Session, eine Anfrage mit
// Ende und demselben SQL-Text, ist sie diese Anfrage. Trägt die gelieferte
// Interaktion ein Ende, ist connection_closed die letzte Antwort (LH-FA-25.a). Erwartet der Cursor eine Extended-Nachricht, ist jede
// andere Anfrage eine Abweichung, mitten in einer Extended-Interaktion auch
// eine Lebendprüfung (LH-FA-18.a).
func (s *ReplayService) Query(_ context.Context, id model.SessionID, sql string) ([]model.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.verbindung[id]; ok && !c.mitten() && istLebendpruefung(sql, c.vt) && !s.anfrageMitEnde(c, sql) {
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
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d: %s; empfangen %q",
			c.session.ID, c.keineWeitere(), sql)
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
	c.ungesendet = true
	c.merke(erwartet.Responses)
	return erwartet.Responses, nil
}

// anfrageMitEnde liefert true, wenn die als Nächste erwartete Interaktion eine
// einfache Anfrage mit Ende und dem SQL-Text sql ist: die am Cursor, bei einer
// noch nicht zugeordneten Verbindung mit `first-request` die erste der nächsten
// freien Session (LH-FA-09.a *Lebendprüfung mit Ende*).
func (s *ReplayService) anfrageMitEnde(c *cursor, sql string) bool {
	var erwartet []model.Interaction
	switch {
	case c.session != nil:
		erwartet = c.session.Interactions[min(c.pos, len(c.session.Interactions)):]
	case !s.nachVerbindung && len(s.frei) > 0:
		erwartet = s.frei[0].Interactions
	}
	return len(erwartet) > 0 && erwartet[0].Request.Type == model.RequestQuery &&
		erwartet[0].Request.SQL == sql && endetEinfach(erwartet[0])
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

// keineWeitere nennt in der Diagnose, was die Aufzeichnung am Ende der Session
// erwartet: nach der letzten erwarteten Interaktion, aufgezeichnete
// Lebendprüfungen nicht gezählt, keine weitere; eine Session ohne Interaktion
// (mit `connection` zugeordnet) erwartet keine (LH-FA-12.a).
func (c *cursor) keineWeitere() string {
	if n := len(c.session.Interactions); n > 0 {
		return fmt.Sprintf("nach Interaktion %d erwartet die Aufzeichnung keine weitere", c.session.Interactions[n-1].Sequence)
	}
	return "die Aufzeichnung erwartet keine Interaktion"
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
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "Session %d: %s; empfangen %s",
			c.session.ID, c.keineWeitere(), c.nachrichtText(m))
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
	case "params":
		return nil, model.Errorf(model.CodeReplayMismatch, nil, "%s: erwartet %s, empfangen %s, abweichend in params%s%s",
			c.stelle(), e.Type, m.Type, parameterStelle(m.Params, e.Params), c.anweisungen(e, m))
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
		c.ungesendet = true
	}
	c.merke(g.Server)
	return g.Server, nil
}

// zuordnen liefert den Cursor der Verbindung und ordnet ihr mit `first-request`
// bei der ersten Anfrage die nächste freie Session zu; ist keine mehr frei, ist
// das PGR-E5003. Mit `connection` ist die Zuordnung mit der Annahme geschehen
// (OpenConnection); eine Verbindung ohne Session ist PGR-E5003.
func (s *ReplayService) zuordnen(id model.SessionID, empfangen string) (*cursor, error) {
	c, ok := s.verbindung[id]
	if !ok {
		return nil, model.Errorf(model.CodeInternal, nil, "unbekannte Verbindung %d", id)
	}
	if c.session == nil {
		if s.nachVerbindung {
			return nil, model.Errorf(model.CodeReplaySession, nil, "Verbindung %d: keine aufgezeichnete Session mit der Kennung %d für %s", id, id, empfangen)
		}
		if len(s.frei) == 0 {
			return nil, model.Errorf(model.CodeReplaySession, nil, "Verbindung %d: keine aufgezeichnete Session mehr frei für %s", id, empfangen)
		}
		sess := s.frei[0]
		s.frei = s.frei[1:]
		c.session = &sess
		c.vt = vtLeerraum(sess.ServerParameters)
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
// Lebendprüfungen ohne Ende stehen nicht in der Session (NewReplayService); dass
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
		halt := func(gi, ni int) bool {
			return pi == c.pos && (gi > c.gruppe || gi == c.gruppe && ni >= c.nachricht)
		}
		if !o.extendedNachspielen(in, halt) {
			return o
		}
		o.transaktionsende(in.Groups[len(in.Groups)-1].Server)
	}
	return o
}

// extendedNachspielen spielt die Client-Nachrichten einer Extended-Interaktion,
// deren Art in der Interaktion bestätigt ist, in Sendereihenfolge nach, bis halt
// für die Nachricht ni der Gruppe gi wahr ist; diese und alle späteren bleiben
// aus. Das Ergebnis ist wahr, wenn halt für keine Nachricht wahr war.
//
// Die Bestätigungen einer Art zählt es über die ganze Interaktion, weil eine
// späte Bestätigung in der folgenden Gruppe steht (LH-FA-18.a); eine Art ohne
// Bestätigung in der Interaktion wird nicht nachgespielt.
func (o objekte) extendedNachspielen(in model.Interaction, halt func(gi, ni int) bool) bool {
	bestaetigt := map[model.ClientMessageType]int{}
	for _, g := range in.Groups {
		for _, r := range g.Server {
			bestaetigt[bestaetigung[r.Type]]++
		}
	}
	for gi, g := range in.Groups {
		for ni, m := range g.Client {
			if halt(gi, ni) {
				return false
			}
			if bestaetigt[m.Type] > 0 {
				bestaetigt[m.Type]--
				o.nachspielen(m)
			}
		}
	}
	return true
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

// Sent meldet, dass die zuletzt gelieferten Antworten gesendet sind; erst
// damit ist eine Interaktion, deren letzte Antworten sie waren, verbraucht
// (LH-FA-03.b).
func (s *ReplayService) Sent(_ context.Context, id model.SessionID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c := s.verbindung[id]; c != nil {
		c.ungesendet = false
	}
}

// Forced stuft das Zwangsende einer Verbindung ein (LH-FA-13.a *Zwangsende*):
// Ist eine Interaktion der zugeordneten Session begonnen und nicht verbraucht,
// ist das PGR-E4006 mit der Kennung der Session und der aufgezeichneten Nummer
// der Interaktion. Begonnen und nicht verbraucht ist die Interaktion vor dem
// Cursor, solange ihre Antworten nicht als gesendet gemeldet sind, sonst die am
// Cursor, wenn der Cursor innerhalb einer Extended-Interaktion steht
// (LH-FA-03.b *Verbraucht*). Eine Verbindung ohne Session liefert nil.
func (s *ReplayService) Forced(_ context.Context, id model.SessionID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.verbindung[id]
	if c == nil || c.session == nil {
		return nil
	}
	var in model.Interaction
	switch {
	case c.ungesendet:
		in = c.session.Interactions[c.pos-1]
	case c.mitten():
		in = c.session.Interactions[c.pos]
	default:
		return nil
	}
	return model.Errorf(model.CodeShutdownTimeout, nil, "Frist beim Herunterfahren abgelaufen: Session %d, Interaktion %d nicht verbraucht",
		c.session.ID, in.Sequence)
}

// CloseConnection beendet die Verbindung. Hat sie eine Session mit nicht
// verbrauchten Interaktionen, liefert sie eine Meldung, die die Session, die
// Zahl der nicht verbrauchten und aller Interaktionen und die aufgezeichnete
// Nummer der ersten nicht verbrauchten nennt: mit FailOnUnconsumed als Fehler
// PGR-E5002, sonst als Warnung PGR-W2001 (LH-FA-03.b). Eine Verbindung ohne
// Session meldet nichts.
func (s *ReplayService) CloseConnection(_ context.Context, id model.SessionID) (*model.Warning, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.verbindung[id]
	delete(s.verbindung, id)
	if c == nil || c.session == nil {
		return nil, nil
	}
	verbraucht := c.pos
	if c.ungesendet {
		verbraucht--
	}
	alle := len(c.session.Interactions)
	if verbraucht >= alle {
		return nil, nil
	}
	return s.meldung("Session %d: %d von %d Interaktionen nicht verbraucht, die erste mit Nummer %d",
		c.session.ID, alle-verbraucht, alle, c.session.Interactions[verbraucht].Sequence)
}

// Unassigned meldet die Sessions mit Interaktionen, die bis zum Ende des Laufs
// keiner Verbindung zugeordnet wurden, mit ihrer Zahl und der Kennung der
// ersten: mit FailOnUnconsumed als Fehler PGR-E5002, sonst als Warnung
// PGR-W2001 (LH-FA-03.b). Mit ZuordnungNachVerbindung sind das die Sessions mit
// Interaktionen, deren Kennung die Zahl der angenommenen Verbindungen übersteigt.
// Ohne solche Session liefert es nil, nil.
func (s *ReplayService) Unassigned() (*model.Warning, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nie := s.frei
	if s.nachVerbindung {
		nie = nil
		for _, sess := range s.frei {
			if sess.ID > int(s.next) {
				nie = append(nie, sess)
			}
		}
	}
	if len(nie) == 0 {
		return nil, nil
	}
	return s.meldung("%d aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung %d",
		len(nie), nie[0].ID)
}

// meldung liefert denselben Text als Fehler PGR-E5002, wenn der Service
// streng ist, sonst als Warnung PGR-W2001.
func (s *ReplayService) meldung(format string, args ...any) (*model.Warning, error) {
	if s.streng {
		return nil, model.Errorf(model.CodeReplayUnconsumed, nil, format, args...)
	}
	return model.Warnf(model.CodeUnconsumed, format, args...), nil
}
