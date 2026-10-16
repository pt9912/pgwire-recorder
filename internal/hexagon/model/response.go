package model

import "errors"

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

// ResponseConnectionClosed ist das Ende der Verbindung durch den Server (LH-FA-25.a,
// Version 2 des Formats): ohne Felder, als letzte Nachricht der letzten Interaktion
// einer Session, unmittelbar nach einer error_response FATAL oder PANIC (SPEC-002).
const ResponseConnectionClosed ResponseType = "connection_closed"

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
		ResponseConnectionClosed:   true,
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
		ResponseConnectionClosed:     true,
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
// Die übrigen Typen tragen keine Felder; connection_closed trägt keines, und Validate
// lehnt es mit einem Feld ab.
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

// ohneFelder ist wahr, wenn keines der Felder von Response belegt ist.
func (r Response) ohneFelder() bool {
	return len(r.Columns) == 0 && len(r.Values) == 0 && r.Tag == "" && len(r.Fields) == 0 &&
		r.Name == "" && r.Value == "" && r.TxStatus == "" && len(r.ParamTypes) == 0
}

// beendetVerbindung ist wahr für eine error_response, deren Schweregrad FATAL oder
// PANIC ist: der Wert des Feldes "V"; fehlt "V" oder ist es leer, zählt "S". Der
// Vergleich ist genau (SPEC-002).
func (r Response) beendetVerbindung() bool {
	if r.Type != ResponseErrorResponse {
		return false
	}
	schweregrad := r.Fields["V"]
	if schweregrad == "" {
		schweregrad = r.Fields["S"]
	}
	return schweregrad == "FATAL" || schweregrad == "PANIC"
}

// pruefeEnde prüft die Stellen von connection_closed in einer Nachrichtenfolge (die
// Antworten einer einfachen Anfrage oder die Server-Nachrichten einer Gruppe):
// connection_closed steht nur als letzte Nachricht, trägt keine Felder und folgt
// unmittelbar auf eine error_response FATAL oder PANIC.
func pruefeEnde(antworten []Response) error {
	for idx, r := range antworten {
		if r.Type != ResponseConnectionClosed {
			continue
		}
		if idx != len(antworten)-1 {
			return errors.New("connection_closed nur als letzte Nachricht")
		}
		if !r.ohneFelder() {
			return errors.New("connection_closed trägt Felder")
		}
		if idx == 0 || !antworten[idx-1].beendetVerbindung() {
			return errors.New("connection_closed nur unmittelbar nach einer error_response FATAL oder PANIC")
		}
	}
	return nil
}

// ohneReadyForQuery prüft die Antworten einer einfachen Anfrage mit Ende: Sie tragen an
// keiner Stelle ein ready_for_query, das Ende steht statt seiner (SPEC-002). Ohne
// connection_closed in den Antworten prüft sie nichts.
func ohneReadyForQuery(antworten []Response) error {
	ende := false
	for _, r := range antworten {
		ende = ende || r.Type == ResponseConnectionClosed
	}
	if !ende {
		return nil
	}
	for _, r := range antworten {
		if r.Type == ResponseReadyForQuery {
			return errors.New("ready_for_query in einer Interaktion mit connection_closed")
		}
	}
	return nil
}

// endet ist wahr, wenn die Interaktion eine Antwort connection_closed trägt, an
// welcher Stelle auch immer.
func (i Interaction) endet() bool {
	for _, r := range i.Responses {
		if r.Type == ResponseConnectionClosed {
			return true
		}
	}
	for _, g := range i.Groups {
		if g.endet() {
			return true
		}
	}
	return false
}

// endet ist wahr, wenn die Server-Nachrichten der Gruppe connection_closed tragen.
func (g Group) endet() bool {
	for _, r := range g.Server {
		if r.Type == ResponseConnectionClosed {
			return true
		}
	}
	return false
}
