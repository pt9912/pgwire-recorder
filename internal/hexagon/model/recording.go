package model

// Formatkennung und Formatversion einer Aufzeichnung (SPEC-010).
const (
	Format  = "pgwire-recorder"
	Version = 1
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

// Interaction ist eine abgeschlossene Anfrage mit allen Serverantworten bis
// einschließlich ReadyForQuery (LH-FA-02.b, LH-FA-18.a). Request.Type nennt die
// Art: Eine einfache Anfrage trägt SQL und Responses, eine Extended-Interaktion
// nur Groups (SPEC-041). Validate prüft diese Form.
type Interaction struct {
	Sequence int
	// OffsetMS ist der Abstand zum Beginn der Session in Millisekunden, nur mit
	// Zeitaufzeichnung (LH-FA-21.a); sonst nil.
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

// SessionEnd ist der Grund, aus dem eine Session endet.
type SessionEnd int

const (
	// EndNormal: Die Verbindung endete nach einem ReadyForQuery; die Session
	// wird mit allen abgeschlossenen Interaktionen übernommen.
	EndNormal SessionEnd = iota
	// EndLost: Die Antwort der zuletzt gelieferten Interaktion erreichte den
	// Client nicht (PGR-E4003); diese Interaktion entfällt, die vorherigen
	// bleiben (LH-FA-02.b).
	EndLost
	// EndUnsupported: In der Session trat eine nicht unterstützte Interaktion
	// auf; die Session wird nicht übernommen (LH-FA-05.a).
	EndUnsupported
)
