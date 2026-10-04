package model

// ResponseType ist die Art einer Serverantwort; die Namen folgen den
// PGWire-Nachrichten in Kleinbuchstaben mit Unterstrich (SPEC-041).
type ResponseType string

// Die Serverantworten einer einfachen Anfrage; eine Extended-Interaktion kennt sie
// ebenfalls.
const (
	ResponseRowDescription     ResponseType = "row_description"
	ResponseDataRow            ResponseType = "data_row"
	ResponseCommandComplete    ResponseType = "command_complete"
	ResponseEmptyQueryResponse ResponseType = "empty_query_response"
	ResponseErrorResponse      ResponseType = "error_response"
	ResponseNoticeResponse     ResponseType = "notice_response"
	ResponseParameterStatus    ResponseType = "parameter_status"
	ResponseReadyForQuery      ResponseType = "ready_for_query"
)

// Die Serverantworten, die nur eine Extended-Interaktion kennt (LH-FA-18.a).
const (
	ResponseParseComplete        ResponseType = "parse_complete"
	ResponseBindComplete         ResponseType = "bind_complete"
	ResponseCloseComplete        ResponseType = "close_complete"
	ResponseParameterDescription ResponseType = "parameter_description"
	ResponseNoData               ResponseType = "no_data"
	ResponsePortalSuspended      ResponseType = "portal_suspended"
)

// simpleResponses sind die Serverantworten einer einfachen Anfrage,
// extendedResponses die einer Gruppe einer Extended-Interaktion; Validate lehnt
// jeden anderen Typ ab (SPEC-001).
var (
	simpleResponses = map[ResponseType]bool{
		ResponseRowDescription:     true,
		ResponseDataRow:            true,
		ResponseCommandComplete:    true,
		ResponseEmptyQueryResponse: true,
		ResponseErrorResponse:      true,
		ResponseNoticeResponse:     true,
		ResponseParameterStatus:    true,
		ResponseReadyForQuery:      true,
	}
	extendedResponses = map[ResponseType]bool{
		ResponseParseComplete:        true,
		ResponseBindComplete:         true,
		ResponseCloseComplete:        true,
		ResponseParameterDescription: true,
		ResponseRowDescription:       true,
		ResponseNoData:               true,
		ResponseDataRow:              true,
		ResponseCommandComplete:      true,
		ResponseEmptyQueryResponse:   true,
		ResponsePortalSuspended:      true,
		ResponseErrorResponse:        true,
		ResponseNoticeResponse:       true,
		ResponseParameterStatus:      true,
		ResponseReadyForQuery:        true,
	}
)

// Response ist eine Serverantwort. Welche Felder belegt sind, hängt vom Typ ab:
//
//	row_description              Columns
//	data_row                     Values
//	command_complete             Tag
//	error_response, notice_...   Fields (Feldcode → Wert, zum Beispiel "C" → SQLSTATE)
//	parameter_status             Name, Value
//	parameter_description        ParamTypes (Typ-OIDs der Parameter)
//	ready_for_query              TxStatus
//
// Die übrigen Typen tragen keine Felder.
type Response struct {
	Type       ResponseType
	Columns    []Column
	Values     []Value
	Tag        string
	Fields     map[string]string
	Name       string
	Value      string
	TxStatus   string
	ParamTypes []uint32
}

// Column ist eine Spaltenbeschreibung einer RowDescription.
type Column struct {
	Name         string
	TableOID     uint32
	ColumnNumber uint16
	TypeOID      uint32
	TypeSize     int16
	TypeModifier int32
	Format       int16
}

// Value ist ein Feldwert. Null unterscheidet SQL-NULL von leeren Bytes.
type Value struct {
	Null  bool
	Bytes []byte
}
