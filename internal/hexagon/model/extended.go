package model

import (
	"errors"
	"fmt"
)

// ClientMessageType ist die Art einer Client-Nachricht des Extended Query
// Protocol; die Namen folgen den PGWire-Nachrichten in Kleinbuchstaben (SPEC-041).
type ClientMessageType string

// Die Client-Nachrichten einer Extended-Interaktion (LH-FA-18.a).
const (
	ClientParse    ClientMessageType = "parse"
	ClientBind     ClientMessageType = "bind"
	ClientDescribe ClientMessageType = "describe"
	ClientExecute  ClientMessageType = "execute"
	ClientClose    ClientMessageType = "close"
	ClientFlush    ClientMessageType = "flush"
	ClientSync     ClientMessageType = "sync"
)

var clientMessages = map[ClientMessageType]bool{
	ClientParse: true, ClientBind: true, ClientDescribe: true, ClientExecute: true,
	ClientClose: true, ClientFlush: true, ClientSync: true,
}

// Target ist die Zielart von describe und close.
type Target string

// Die Zielarten von describe und close.
const (
	TargetStatement Target = "statement"
	TargetPortal    Target = "portal"
)

// ClientMessage ist eine Client-Nachricht einer Extended-Interaktion. Welche
// Felder belegt sind, hängt vom Typ ab:
//
//	parse              Statement, SQL, ParamTypes (Typ-OIDs)
//	bind               Portal, Statement, ParamFormats, Params, ResultFormats
//	describe, close    Target, Name
//	execute            Portal, MaxRows
//	flush, sync        —
//
// Ein Parameter mit Value.Null ist SQL-NULL.
type ClientMessage struct {
	Type          ClientMessageType
	Statement     string
	Portal        string
	SQL           string
	ParamTypes    []uint32
	ParamFormats  []int16
	Params        []Value
	ResultFormats []int16
	Target        Target
	Name          string
	MaxRows       uint32
}

// Group ist eine Gruppe einer Extended-Interaktion: die Client-Nachrichten bis
// einschließlich Flush oder Sync in Sendereihenfolge, danach die
// Server-Nachrichten, die auf sie antworten, in Empfangsreihenfolge (LH-FA-18.a).
type Group struct {
	Client []ClientMessage
	Server []Response
}

// Validate prüft die Form einer Interaktion nach ihrer Art.
//
// Eine einfache Anfrage trägt keine Gruppen, nur Antworttypen des Simple Query
// Protocol, und ihre letzte Antwort ist ready_for_query (LH-FA-02.b).
//
// Eine Extended-Interaktion trägt weder SQL noch Antworten außerhalb ihrer
// Gruppen und mindestens eine Gruppe (LH-FA-18.a). Jede Gruppe hat mindestens
// eine Client-Nachricht; die letzte ist flush oder sync, und keine davor ist es.
// Nur die letzte Gruppe endet mit sync. Ihre letzte Server-Nachricht ist
// ready_for_query, und keine andere Server-Nachricht der Interaktion ist es.
// describe und close tragen die Zielart statement oder portal.
//
// Ein Nachrichten- oder Anfragetyp, den die Art nicht kennt, ist ein Fehler
// (SPEC-001).
func (i Interaction) Validate() error {
	switch i.Request.Type {
	case RequestQuery:
		return i.validateQuery()
	case RequestExtended:
		return i.validateExtended()
	default:
		return fmt.Errorf("Anfrage-Typ %q unbekannt", i.Request.Type)
	}
}

func (i Interaction) validateQuery() error {
	if len(i.Groups) > 0 {
		return errors.New("einfache Anfrage mit Gruppen")
	}
	for _, r := range i.Responses {
		if !simpleResponses[r.Type] {
			return fmt.Errorf("Antwort-Typ %q unbekannt in einer einfachen Anfrage", r.Type)
		}
	}
	if n := len(i.Responses); n == 0 || i.Responses[n-1].Type != ResponseReadyForQuery {
		return errors.New("endet nicht mit ready_for_query")
	}
	return nil
}

func (i Interaction) validateExtended() error {
	if i.Request.SQL != "" || len(i.Responses) > 0 {
		return errors.New("Extended-Interaktion mit SQL oder Antworten außerhalb der Gruppen")
	}
	if len(i.Groups) == 0 {
		return errors.New("Extended-Interaktion ohne Gruppe")
	}
	letzte := len(i.Groups) - 1
	for gi, g := range i.Groups {
		if err := g.validate(gi == letzte); err != nil {
			return fmt.Errorf("Gruppe %d: %w", gi+1, err)
		}
	}
	return nil
}

// validate prüft eine Gruppe; letzte sagt, ob sie die Interaktion abschließt.
// Die Prüfungen laufen in dieser Reihenfolge, und die erste verletzte meldet
// ihren Fehler: die Client-Nachrichten, sync oder flush am Ende der Gruppe, die
// Server-Nachrichten, ready_for_query am Ende der letzten Gruppe.
func (g Group) validate(letzte bool) error {
	if err := g.validateClient(); err != nil {
		return err
	}
	endetMitSync := g.Client[len(g.Client)-1].Type == ClientSync
	if endetMitSync != letzte {
		return errors.New("sync steht genau am Ende der letzten Gruppe, flush am Ende jeder anderen")
	}
	if err := g.validateServer(letzte); err != nil {
		return err
	}
	if letzte && (len(g.Server) == 0 || g.Server[len(g.Server)-1].Type != ResponseReadyForQuery) {
		return errors.New("letzte Gruppe endet nicht mit ready_for_query")
	}
	return nil
}

// validateClient prüft die Client-Nachrichten der Reihe nach, je Nachricht den
// Typ, die Stellung von flush und sync und die Zielart von describe und close.
func (g Group) validateClient() error {
	if len(g.Client) == 0 {
		return errors.New("ohne Client-Nachricht")
	}
	for ci, m := range g.Client {
		if !clientMessages[m.Type] {
			return fmt.Errorf("Client-Nachricht %q unbekannt", m.Type)
		}
		ende := m.Type == ClientFlush || m.Type == ClientSync
		if ende != (ci == len(g.Client)-1) {
			return fmt.Errorf("Client-Nachricht %d: flush oder sync nur als letzte Nachricht der Gruppe", ci+1)
		}
		if (m.Type == ClientDescribe || m.Type == ClientClose) && m.Target != TargetStatement && m.Target != TargetPortal {
			return fmt.Errorf("Client-Nachricht %d: Zielart %q statt statement oder portal", ci+1, m.Target)
		}
	}
	return nil
}

// validateServer prüft die Server-Nachrichten der Reihe nach, je Nachricht den
// Typ und die Stellung von ready_for_query; letzte sagt, ob die Gruppe die
// Interaktion abschließt.
func (g Group) validateServer(letzte bool) error {
	for si, r := range g.Server {
		if !extendedResponses[r.Type] {
			return fmt.Errorf("Server-Nachricht %q unbekannt in einer Extended-Interaktion", r.Type)
		}
		if r.Type == ResponseReadyForQuery && (!letzte || si != len(g.Server)-1) {
			return fmt.Errorf("Server-Nachricht %d: ready_for_query nur als letzte Nachricht der letzten Gruppe", si+1)
		}
	}
	return nil
}
