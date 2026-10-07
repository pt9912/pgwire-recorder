package postgres

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// Upstream verbindet sich je Session mit dem realen PostgreSQL-Server.
// Vermittelt wird ein Verbindungsaufbau ohne Passwort-Anmeldung; verlangt der
// Server ein Anmeldeverfahren, endet der Aufbau mit PGR-E6001.
type Upstream struct {
	Address string
	Dialer  net.Dialer
}

var _ driven.Upstream = (*Upstream)(nil)

// Open baut die Verbindung auf und liefert die Antworten bis zum ersten
// ReadyForQuery oder bis zur Fehlerantwort des Servers.
func (u *Upstream) Open(ctx context.Context, startup map[string]string) (driven.UpstreamSession, []model.Response, error) {
	conn, err := u.Dialer.DialContext(ctx, "tcp", u.Address)
	if err != nil {
		return nil, nil, model.Errorf(model.CodeUpstream, err, "Upstream %s nicht erreichbar", u.Address)
	}
	s := newSession(conn)
	fe := s.fe
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: startup})
	if err := fe.Flush(); err != nil {
		conn.Close()
		return nil, nil, model.Errorf(model.CodeUpstream, err, "Startup an %s nicht zu senden", u.Address)
	}

	var responses []model.Response
	for {
		msg, err := fe.Receive()
		if err != nil {
			conn.Close()
			return nil, responses, model.Errorf(model.CodeConnectionLost, err, "Verbindung zum Upstream im Verbindungsaufbau beendet")
		}
		switch m := msg.(type) {
		case *pgproto3.AuthenticationOk:
			// Ohne Passwort-Anmeldung; der Client erhält AuthenticationOk vom Adapter.
		case *pgproto3.AuthenticationCleartextPassword, *pgproto3.AuthenticationMD5Password, *pgproto3.AuthenticationSASL:
			conn.Close()
			return nil, responses, model.Errorf(model.CodeUnsupported, nil, "der Upstream verlangt ein Anmeldeverfahren (%T), das dieser Stand nicht vermittelt", m)
		case *pgproto3.ParameterStatus:
			responses = append(responses, model.Response{Type: model.ResponseParameterStatus, Name: m.Name, Value: m.Value})
		case *pgproto3.BackendKeyData:
			// Die Abbruchkennung des Servers geht nicht an den Client und nicht in
			// die Aufzeichnung (SPEC-004).
		case *pgproto3.NoticeResponse:
			responses = append(responses, model.Response{Type: model.ResponseNoticeResponse, Fields: noticeFields((*pgproto3.ErrorResponse)(m))})
		case *pgproto3.ErrorResponse:
			conn.Close()
			return nil, append(responses, model.Response{Type: model.ResponseErrorResponse, Fields: noticeFields(m)}), nil
		case *pgproto3.ReadyForQuery:
			responses = append(responses, model.Response{Type: model.ResponseReadyForQuery, TxStatus: string(m.TxStatus)})
			return s, responses, nil
		default:
			conn.Close()
			return nil, responses, model.Errorf(model.CodeUnsupported, nil, "Nachricht %T im Verbindungsaufbau wird nicht unterstützt", m)
		}
	}
}

// session ist eine Verbindung zum Server. schreiben hält den Schreibpuffer des
// Frontends für Query, Send und das Terminate von Close; der Lesepuffer gehört
// dem einen laufenden Query oder Receive. Das Frontend hält beide Puffer
// getrennt, darum laufen Send und Receive gleichzeitig.
type session struct {
	conn      net.Conn
	fe        *pgproto3.Frontend
	schreiben sync.Mutex
	// merker schützt fehler: Die Leserichtung setzt es, Lese- und
	// Senderichtung stufen danach ein.
	merker sync.Mutex
	// fehler sind die Felder der letzten ErrorResponse seit dem letzten
	// ReadyForQuery, sonst nil.
	fehler map[string]string
}

// newSession legt die Session zu conn mit ihrem Frontend an.
func newSession(conn net.Conn) *session {
	return &session{conn: conn, fe: pgproto3.NewFrontend(conn, conn)}
}

// terminateFrist begrenzt das Senden von Terminate in Close.
const terminateFrist = 100 * time.Millisecond

func (s *session) Query(_ context.Context, sql string) ([]model.Response, error) {
	s.schreiben.Lock()
	s.fe.Send(&pgproto3.Query{String: sql})
	err := s.fe.Flush()
	s.schreiben.Unlock()
	if err != nil {
		return nil, model.Errorf(model.CodeConnectionLost, err, "Anfrage an den Upstream nicht zu senden")
	}
	var responses []model.Response
	for {
		r, err := s.lies()
		if err != nil {
			return responses, err
		}
		responses = append(responses, r)
		if r.Type == model.ResponseReadyForQuery {
			return responses, nil
		}
	}
}

// Send bildet die Client-Nachrichten auf PGWire ab und sendet sie in einem
// Schreibvorgang; es blockiert, solange der Server nicht liest. Scheitert das
// Senden, stuft es wie lies nach einer gelesenen ErrorResponse ein
// (LH-FA-02.b).
func (s *session) Send(_ context.Context, msgs []model.ClientMessage) error {
	s.schreiben.Lock()
	defer s.schreiben.Unlock()
	for _, m := range msgs {
		msg, err := toFrontendMessage(m)
		if err != nil {
			return err
		}
		s.fe.Send(msg)
	}
	if err := s.fe.Flush(); err != nil {
		return s.verbindungsende(err, "Nachrichten an den Upstream nicht zu senden")
	}
	return nil
}

// Receive wartet auf eine Server-Nachricht und liest weiter, solange der
// Lesepuffer schon Bytes der nächsten enthält; ein ReadyForQuery beendet die
// Folge.
func (s *session) Receive(context.Context) ([]model.Response, error) {
	var out []model.Response
	for {
		r, err := s.lies()
		if err != nil {
			return nil, err
		}
		out = append(out, r)
		if r.Type == model.ResponseReadyForQuery || s.fe.ReadBufferLen() == 0 {
			return out, nil
		}
	}
}

// lies liest die nächste Server-Nachricht und bildet sie ab. Eine
// ErrorResponse merkt es sich bis zum nächsten ReadyForQuery. Endet die
// Verbindung (istVerbindungsende), stuft verbindungsende ein; jeder andere
// Lesefehler ist eine nicht lesbare Serverantwort, PGR-E6001 ohne Blick auf
// den Merker (LH-FA-02.b, LH-FA-05.a).
func (s *session) lies() (model.Response, error) {
	msg, err := s.fe.Receive()
	if err != nil {
		if istVerbindungsende(err) {
			return model.Response{}, s.verbindungsende(err, "Verbindung zum Upstream vor ReadyForQuery beendet")
		}
		return model.Response{}, model.Errorf(model.CodeUnsupported, err, "Serverantwort des Upstreams nicht lesbar")
	}
	r, err := toResponse(msg)
	if err != nil {
		return model.Response{}, err
	}
	switch r.Type {
	case model.ResponseErrorResponse:
		s.merke(r.Fields)
	case model.ResponseReadyForQuery:
		s.merke(nil)
	}
	return r, nil
}

func (s *session) merke(fehler map[string]string) {
	s.merker.Lock()
	defer s.merker.Unlock()
	s.fehler = fehler
}

// verbindungsende stuft das Ende der Verbindung ein: mit gemerkter
// ErrorResponse PGR-E6001 mit deren SQLSTATE und Meldung, sonst PGR-E4003
// mit text (LH-FA-02.b).
func (s *session) verbindungsende(err error, text string) error {
	s.merker.Lock()
	fehler := s.fehler
	s.merker.Unlock()
	if fehler != nil {
		return model.Errorf(model.CodeUnsupported, err,
			"Verbindung zum Upstream nach der Fehlerantwort %s „%s“ vor ReadyForQuery beendet", fehler["C"], fehler["M"])
	}
	return model.Errorf(model.CodeConnectionLost, err, "%s", text)
}

// istVerbindungsende ist wahr für das Ende des Datenstroms, auch mitten in
// einer Nachricht, und für einen Fehler der Verbindung selbst.
func istVerbindungsende(err error) bool {
	var ne net.Error
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, net.ErrClosed) || errors.As(err, &ne)
}

// Close sendet Terminate nur, wenn gerade kein Send oder Query schreibt, und
// höchstens terminateFrist lang; ein Fehler dabei ist ohne Folge, denn die
// Verbindung endet ohnehin. Das Schließen der Verbindung beendet ein wartendes
// Send, Receive oder Query mit einem Fehler.
func (s *session) Close() error {
	if s.schreiben.TryLock() {
		_ = s.conn.SetWriteDeadline(time.Now().Add(terminateFrist))
		s.fe.Send(&pgproto3.Terminate{})
		_ = s.fe.Flush()
		s.schreiben.Unlock()
	}
	return s.conn.Close()
}

// toResponse bildet eine Serverantwort auf das Modell ab; jede andere als die
// hier genannten ist ein Fehler mit CodeUnsupported.
func toResponse(msg pgproto3.BackendMessage) (model.Response, error) {
	switch m := msg.(type) {
	case *pgproto3.RowDescription:
		return model.Response{Type: model.ResponseRowDescription, Columns: spalten(m.Fields)}, nil
	case *pgproto3.DataRow:
		return model.Response{Type: model.ResponseDataRow, Values: werte(m.Values)}, nil
	case *pgproto3.CommandComplete:
		return model.Response{Type: model.ResponseCommandComplete, Tag: string(m.CommandTag)}, nil
	case *pgproto3.ErrorResponse:
		return model.Response{Type: model.ResponseErrorResponse, Fields: noticeFields(m)}, nil
	case *pgproto3.NoticeResponse:
		return model.Response{Type: model.ResponseNoticeResponse, Fields: noticeFields((*pgproto3.ErrorResponse)(m))}, nil
	case *pgproto3.ParameterStatus:
		return model.Response{Type: model.ResponseParameterStatus, Name: m.Name, Value: m.Value}, nil
	case *pgproto3.ReadyForQuery:
		return model.Response{Type: model.ResponseReadyForQuery, TxStatus: string(m.TxStatus)}, nil
	case *pgproto3.ParameterDescription:
		return model.Response{Type: model.ResponseParameterDescription, ParamTypes: kopieOderNil(m.ParameterOIDs)}, nil
	}
	if typ, ok := ohneFelder(msg); ok {
		return model.Response{Type: typ}, nil
	}
	return model.Response{}, model.Errorf(model.CodeUnsupported, nil, "Serverantwort %T wird nicht unterstützt", msg)
}

// ohneFelder liefert den Typ einer Serverantwort, die außer ihrem Typ nichts
// trägt; ok ist falsch für jede andere.
func ohneFelder(msg pgproto3.BackendMessage) (typ model.ResponseType, ok bool) {
	switch msg.(type) {
	case *pgproto3.EmptyQueryResponse:
		return model.ResponseEmptyQueryResponse, true
	case *pgproto3.ParseComplete:
		return model.ResponseParseComplete, true
	case *pgproto3.BindComplete:
		return model.ResponseBindComplete, true
	case *pgproto3.CloseComplete:
		return model.ResponseCloseComplete, true
	case *pgproto3.NoData:
		return model.ResponseNoData, true
	case *pgproto3.PortalSuspended:
		return model.ResponsePortalSuspended, true
	}
	return "", false
}

// spalten bildet die Felder einer RowDescription auf Spalten ab.
func spalten(fields []pgproto3.FieldDescription) []model.Column {
	cols := make([]model.Column, len(fields))
	for i, f := range fields {
		cols[i] = model.Column{
			Name:         string(f.Name),
			TableOID:     f.TableOID,
			ColumnNumber: f.TableAttributeNumber,
			TypeOID:      f.DataTypeOID,
			TypeSize:     f.DataTypeSize,
			TypeModifier: f.TypeModifier,
			Format:       f.Format,
		}
	}
	return cols
}

// werte bildet die Werte einer DataRow ab: nil wird zu NULL, jeder andere Wert
// zu einer Kopie seiner Bytes.
func werte(values [][]byte) []model.Value {
	vals := make([]model.Value, len(values))
	for i, v := range values {
		if v == nil {
			vals[i] = model.Value{Null: true}
		} else {
			vals[i] = model.Value{Bytes: append([]byte(nil), v...)}
		}
	}
	return vals
}

// toFrontendMessage bildet eine Client-Nachricht einer Extended-Interaktion
// auf PGWire ab; ein Parameter mit Value.Null wird zu NULL, jeder andere zu
// seinen Bytes, auch leer. Eine Zielart außer statement und portal und ein
// unbekannter Typ sind PGR-E1000.
func toFrontendMessage(m model.ClientMessage) (pgproto3.FrontendMessage, error) {
	switch m.Type {
	case model.ClientParse:
		return &pgproto3.Parse{Name: m.Statement, Query: m.SQL, ParameterOIDs: m.ParamTypes}, nil
	case model.ClientBind:
		params := make([][]byte, len(m.Params))
		for i, v := range m.Params {
			if !v.Null {
				params[i] = append([]byte{}, v.Bytes...)
			}
		}
		return &pgproto3.Bind{
			DestinationPortal: m.Portal, PreparedStatement: m.Statement,
			ParameterFormatCodes: m.ParamFormats, Parameters: params, ResultFormatCodes: m.ResultFormats,
		}, nil
	case model.ClientDescribe, model.ClientClose:
		ziel, ok := zielart[m.Target]
		if !ok {
			return nil, model.Errorf(model.CodeInternal, nil, "%s mit Zielart %q", m.Type, m.Target)
		}
		if m.Type == model.ClientDescribe {
			return &pgproto3.Describe{ObjectType: ziel, Name: m.Name}, nil
		}
		return &pgproto3.Close{ObjectType: ziel, Name: m.Name}, nil
	case model.ClientExecute:
		return &pgproto3.Execute{Portal: m.Portal, MaxRows: m.MaxRows}, nil
	case model.ClientFlush:
		return &pgproto3.Flush{}, nil
	case model.ClientSync:
		return &pgproto3.Sync{}, nil
	default:
		return nil, model.Errorf(model.CodeInternal, nil, "Client-Nachricht %q ohne Abbildung", m.Type)
	}
}

// zielart ist das PGWire-Byte der Zielart von Describe und Close.
var zielart = map[model.Target]byte{model.TargetStatement: 'S', model.TargetPortal: 'P'}

// kopieOderNil kopiert eine Liste; eine leere wird nil, wie sie der Leser der
// Aufzeichnung liefert.
func kopieOderNil[T any](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	return append([]T(nil), s...)
}

// noticeFields bildet die Felder einer Fehler- oder Hinweisantwort auf ihre
// PGWire-Feldcodes ab; leere Felder entfallen.
func noticeFields(e *pgproto3.ErrorResponse) map[string]string {
	f := map[string]string{}
	set := func(code, v string) {
		if v != "" {
			f[code] = v
		}
	}
	setInt := func(code string, v int32) {
		if v != 0 {
			f[code] = fmt.Sprint(v)
		}
	}
	set("S", e.Severity)
	set("V", e.SeverityUnlocalized)
	set("C", e.Code)
	set("M", e.Message)
	set("D", e.Detail)
	set("H", e.Hint)
	setInt("P", e.Position)
	setInt("p", e.InternalPosition)
	set("q", e.InternalQuery)
	set("W", e.Where)
	set("s", e.SchemaName)
	set("t", e.TableName)
	set("c", e.ColumnName)
	set("d", e.DataTypeName)
	set("n", e.ConstraintName)
	set("F", e.File)
	setInt("L", e.Line)
	set("R", e.Routine)
	for k, v := range e.UnknownFields {
		f[string(k)] = v
	}
	return f
}
