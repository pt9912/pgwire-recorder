package postgres

import (
	"context"
	"errors"
	"fmt"
	"net"

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
	fe := pgproto3.NewFrontend(conn, conn)
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
			return &session{conn: conn, fe: fe}, responses, nil
		default:
			conn.Close()
			return nil, responses, model.Errorf(model.CodeUnsupported, nil, "Nachricht %T im Verbindungsaufbau wird nicht unterstützt", m)
		}
	}
}

type session struct {
	conn net.Conn
	fe   *pgproto3.Frontend
}

func (s *session) Query(ctx context.Context, sql string) ([]model.Response, error) {
	s.fe.Send(&pgproto3.Query{String: sql})
	if err := s.fe.Flush(); err != nil {
		return nil, model.Errorf(model.CodeConnectionLost, err, "Anfrage an den Upstream nicht zu senden")
	}
	var responses []model.Response
	for {
		msg, err := s.fe.Receive()
		if err != nil {
			return responses, model.Errorf(model.CodeConnectionLost, err, "Verbindung zum Upstream vor ReadyForQuery beendet")
		}
		r, err := toResponse(msg)
		if err != nil {
			return responses, err
		}
		responses = append(responses, r)
		if r.Type == model.ResponseReadyForQuery {
			return responses, nil
		}
	}
}

func (s *session) Close() error {
	s.fe.Send(&pgproto3.Terminate{})
	flushErr := s.fe.Flush()
	closeErr := s.conn.Close()
	return errors.Join(flushErr, closeErr)
}

func toResponse(msg pgproto3.BackendMessage) (model.Response, error) {
	switch m := msg.(type) {
	case *pgproto3.RowDescription:
		cols := make([]model.Column, len(m.Fields))
		for i, f := range m.Fields {
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
		return model.Response{Type: model.ResponseRowDescription, Columns: cols}, nil
	case *pgproto3.DataRow:
		vals := make([]model.Value, len(m.Values))
		for i, v := range m.Values {
			if v == nil {
				vals[i] = model.Value{Null: true}
			} else {
				vals[i] = model.Value{Bytes: append([]byte(nil), v...)}
			}
		}
		return model.Response{Type: model.ResponseDataRow, Values: vals}, nil
	case *pgproto3.CommandComplete:
		return model.Response{Type: model.ResponseCommandComplete, Tag: string(m.CommandTag)}, nil
	case *pgproto3.EmptyQueryResponse:
		return model.Response{Type: model.ResponseEmptyQueryResponse}, nil
	case *pgproto3.ErrorResponse:
		return model.Response{Type: model.ResponseErrorResponse, Fields: noticeFields(m)}, nil
	case *pgproto3.NoticeResponse:
		return model.Response{Type: model.ResponseNoticeResponse, Fields: noticeFields((*pgproto3.ErrorResponse)(m))}, nil
	case *pgproto3.ParameterStatus:
		return model.Response{Type: model.ResponseParameterStatus, Name: m.Name, Value: m.Value}, nil
	case *pgproto3.ReadyForQuery:
		return model.Response{Type: model.ResponseReadyForQuery, TxStatus: string(m.TxStatus)}, nil
	default:
		return model.Response{}, model.Errorf(model.CodeUnsupported, nil, "Serverantwort %T wird nicht unterstützt", m)
	}
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
