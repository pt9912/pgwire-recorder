package model

import "fmt"

// Formatkennung und Formatversion einer Aufzeichnung (SPEC-010).
const (
	Format  = "pgwire-recorder"
	Version = 1
	// VersionMitEnde ist die Formatversion, die das Ende des Servers
	// (ResponseConnectionClosed) kennt (SPEC-001); Version bleibt die der übrigen
	// Formate.
	VersionMitEnde = 2
)

// Recording ist eine Aufzeichnung: geordnete Sessions mit ihren Interaktionen
// (SPEC-002).
type Recording struct {
	Format  string
	Version int
	// EmptySessions kennzeichnet eine Aufzeichnung, die auch Verbindungen ohne
	// Anfrage als Session trägt (LH-FA-12.a).
	EmptySessions bool
	Sessions      []Session
}

// NewRecording liefert eine leere Aufzeichnung mit Formatkennung und Version.
func NewRecording() Recording {
	return Recording{Format: Format, Version: Version}
}

// Session ist eine aufgezeichnete Client-Verbindung. ID zählt ab 1 lückenlos in
// der Reihenfolge, in der die Sessions in die Aufzeichnung übernommen werden.
type Session struct {
	ID int
	// Startup trägt die Parameter der StartupMessage des Clients unverändert.
	Startup map[string]string
	// ServerParameters trägt die ParameterStatus-Werte, die der Server beim
	// Verbindungsaufbau gesendet hat.
	ServerParameters map[string]string
	Interactions     []Interaction
}

// ValidateEnd prüft die Form der Session für das Ende des Servers: connection_closed
// steht höchstens einmal, in der letzten Interaktion (SPEC-002). Die Stellen innerhalb
// einer Interaktion prüft Interaction.Validate; Aufrufe von ValidateEnd und Validate
// sind Zustandsprüfungen ohne Ordnung über mehrere Aufrufe.
func (s Session) ValidateEnd() error {
	letzte := len(s.Interactions) - 1
	for ii, i := range s.Interactions {
		if ii != letzte && i.endet() {
			return fmt.Errorf("Interaktion %d: connection_closed nur in der letzten Interaktion der Session", ii+1)
		}
	}
	return nil
}

// Interaction ist eine abgeschlossene Anfrage mit allen Serverantworten bis
// einschließlich ReadyForQuery (LH-FA-02.b, LH-FA-18.a). Request.Type nennt die
// Art: Eine einfache Anfrage trägt SQL und Responses, eine Extended-Interaktion
// nur Groups (SPEC-041). Validate prüft diese Form.
type Interaction struct {
	Sequence int
	// OffsetMS ist der Unterschied zweier Lesungen der Uhr in Millisekunden: der
	// Lesung bei der Ankunft der ersten Client-Nachricht der Interaktion und der
	// Lesung beim Beginn der Session (LH-FA-21.a *Messpunkt*); nur mit
	// Zeitaufzeichnung, sonst nil.
	OffsetMS  *int64
	Request   Request
	Responses []Response
	Groups    []Group
}

// RequestType ist die Art einer Anfrage.
type RequestType string

const (
	// RequestQuery ist eine einfache Anfrage (Simple Query Protocol).
	RequestQuery RequestType = "query"
	// RequestExtended ist eine Extended-Interaktion (Extended Query Protocol,
	// LH-FA-18.a); ihre Nachrichten stehen in Interaction.Groups.
	RequestExtended RequestType = "extended"
)

// Request ist die Anfrage einer Interaktion; SQL ist nur bei RequestQuery
// belegt.
type Request struct {
	Type RequestType
	SQL  string
}

// SessionID kennzeichnet eine offene Client-Verbindung innerhalb eines Laufs;
// im Record ist sie zugleich die laufende Session, im Replay die Verbindung, der
// der Service eine aufgezeichnete Session zuordnet. 0 heißt „keine“.
type SessionID int64

// SessionEnd ist das Ereignis der Verbindung, mit dem eine Session endet. Der
// Driving Adapter meldet, was geschah; welche Interaktionen übernommen werden
// und ob das Ende ein Verbindungsfehler ist, entscheidet der Use Case.
type SessionEnd int

const (
	// EndClosed heißt: Die Client-Verbindung endete ohne Terminate.
	EndClosed SessionEnd = iota
	// EndTerminate heißt: Der Client sandte Terminate.
	EndTerminate
	// EndWriteFailed heißt: Eine Antwort ließ sich nicht an den Client schreiben.
	EndWriteFailed
	// EndShutdown heißt: Der Lauf endet, und der Use Case hat das Ende der Session
	// freigegeben.
	EndShutdown
	// EndUnsupported heißt: Der Client sandte eine nicht unterstützte oder nicht
	// lesbare Nachricht, oder eine Antwort ließ sich nicht auf eine
	// PGWire-Nachricht abbilden.
	EndUnsupported
	// EndFailed heißt: Ein Aufruf des Use Case endete mit einem Fehler.
	EndFailed
	// EndForced heißt: Die Frist beim Herunterfahren ist abgelaufen, und der
	// Adapter hat die Verbindung zwangsweise beendet (LH-FA-13.a). Ohne laufende
	// und ohne abgeschlossene Interaktion schreibt der Record die Session nicht,
	// auch mit LeereSessions (LH-FA-12.a).
	EndForced
)
