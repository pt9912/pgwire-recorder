package model

// ResponseType ist die Art einer Serverantwort; die Namen folgen den
// PGWire-Nachrichten in Kleinbuchstaben mit Unterstrich (SPEC-041).
type ResponseType string

// Die Serverantworten einer einfachen Anfrage.
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

// Response ist eine Serverantwort. Welche Felder belegt sind, hängt vom Typ ab:
//
//	row_description              Columns
//	data_row                     Values
//	command_complete             Tag
//	error_response, notice_...   Fields (Feldcode → Wert, zum Beispiel "C" → SQLSTATE)
//	parameter_status             Name, Value
//	ready_for_query              TxStatus
type Response struct {
	Type     ResponseType
	Columns  []Column
	Values   []Value
	Tag      string
	Fields   map[string]string
	Name     string
	Value    string
	TxStatus string
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
