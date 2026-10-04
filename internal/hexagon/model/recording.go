package model

// Formatkennung und Formatversion einer Aufzeichnung (SPEC-010).
const (
	Format  = "pgwire-recorder"
	Version = 1
)

// Recording ist eine Aufzeichnung: geordnete Sessions mit ihren Interaktionen
// (SPEC-002).
type Recording struct {
	Format   string
	Version  int
	Sessions []Session
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

// Interaction ist eine abgeschlossene Anfrage mit allen Serverantworten bis
// einschließlich ReadyForQuery (LH-FA-02.b).
type Interaction struct {
	Sequence  int
	Request   Request
	Responses []Response
}

// RequestType ist die Art einer Anfrage.
type RequestType string

// RequestQuery ist eine einfache Anfrage (Simple Query Protocol).
const RequestQuery RequestType = "query"

// Request ist die Anfrage einer Interaktion.
type Request struct {
	Type RequestType
	SQL  string
}

// SessionID kennzeichnet eine laufende Session innerhalb eines Laufs; 0 heißt
// „keine Session“.
type SessionID int64
