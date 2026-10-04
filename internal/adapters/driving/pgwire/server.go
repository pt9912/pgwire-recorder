package pgwire

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driving"
)

// Server nimmt PostgreSQL-Clients an und übersetzt ihre Nachrichten in Aufrufe
// des Record-Use-Cases (ARC-006).
type Server struct {
	Recorder driving.Recorder
	Log      *slog.Logger

	wg sync.WaitGroup
}

// Listen öffnet den TCP-Endpunkt; ein nicht zu öffnender Port ist PGR-E4001.
func Listen(address string) (net.Listener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, model.Errorf(model.CodeListen, err, "Adresse %s nicht nutzbar", address)
	}
	return l, nil
}

// Serve nimmt Verbindungen an, bis der Listener geschlossen wird, und wartet
// danach, bis alle Verbindungen beendet sind.
func (s *Server) Serve(ctx context.Context, l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			s.wg.Wait()
			return
		}
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handle(ctx, conn)
		}()
	}
}

func (s *Server) handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	be := pgproto3.NewBackend(conn, conn)

	startup, ok := s.startup(conn, be)
	if !ok {
		return
	}

	id, responses, err := s.Recorder.OpenSession(ctx, startup.Parameters)
	if err != nil {
		s.fail(be, err)
		return
	}
	if id == 0 {
		// Der Server hat den Aufbau mit einer Fehlerantwort beendet.
		s.send(be, responses)
		return
	}
	be.Send(&pgproto3.AuthenticationOk{})
	s.send(be, responses)

	for {
		msg, err := be.Receive()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				s.Log.Warn("Client-Verbindung beendet", "error", err.Error())
			}
			s.close(ctx, id)
			return
		}
		switch m := msg.(type) {
		case *pgproto3.Query:
			out, err := s.Recorder.Query(ctx, id, m.String)
			if err != nil {
				s.fail(be, err)
				s.close(ctx, id)
				return
			}
			s.send(be, out)
		case *pgproto3.Terminate:
			s.close(ctx, id)
			return
		default:
			s.fail(be, model.Errorf(model.CodeUnsupported, nil, "Client-Nachricht %T wird nicht unterstützt", m))
			s.close(ctx, id)
			return
		}
	}
}

// startup liest die StartupMessage. SSL- und GSS-Anfragen beantwortet er mit
// „N“ (LH-FA-05.c); ein CancelRequest schließt die Verbindung (PGR-W3001).
func (s *Server) startup(conn net.Conn, be *pgproto3.Backend) (*pgproto3.StartupMessage, bool) {
	for {
		msg, err := be.ReceiveStartupMessage()
		if err != nil {
			s.Log.Warn("Verbindungsaufbau des Clients gescheitert", "error", err.Error())
			return nil, false
		}
		switch m := msg.(type) {
		case *pgproto3.SSLRequest, *pgproto3.GSSEncRequest:
			if _, err := conn.Write([]byte{'N'}); err != nil {
				return nil, false
			}
		case *pgproto3.CancelRequest:
			s.Log.Warn("CancelRequest empfangen und nicht weitergeleitet", "code", model.CodeCancelRequest)
			return nil, false
		case *pgproto3.StartupMessage:
			if m.ProtocolVersion != pgproto3.ProtocolVersionNumber {
				s.fail(be, model.Errorf(model.CodeProtocolVersion, nil, "Protokollversion %d wird nicht unterstützt", m.ProtocolVersion))
				return nil, false
			}
			return m, true
		default:
			s.fail(be, model.Errorf(model.CodeUnsupported, nil, "Startup-Nachricht %T wird nicht unterstützt", m))
			return nil, false
		}
	}
}

func (s *Server) close(ctx context.Context, id model.SessionID) {
	if err := s.Recorder.CloseSession(ctx, id); err != nil {
		s.logError(err)
	}
}

// fail protokolliert den Fehler und stellt ihn dem Client als FATAL-ErrorResponse
// mit dem Meldungscode im Text zu.
func (s *Server) fail(be *pgproto3.Backend, err error) {
	s.logError(err)
	be.Send(&pgproto3.ErrorResponse{Severity: "FATAL", SeverityUnlocalized: "FATAL", Code: "08006", Message: err.Error()})
	_ = be.Flush()
}

func (s *Server) logError(err error) {
	code := model.CodeInternal
	var me *model.Error
	if errors.As(err, &me) {
		code = me.Code
	}
	s.Log.Error("Verbindungsfehler", "code", code, "error", err.Error())
}

func (s *Server) send(be *pgproto3.Backend, responses []model.Response) {
	for _, r := range responses {
		if msg := toMessage(r); msg != nil {
			be.Send(msg)
		}
	}
	if err := be.Flush(); err != nil {
		s.Log.Warn("Antwort an den Client nicht zu senden", "error", err.Error())
	}
}

func toMessage(r model.Response) pgproto3.BackendMessage {
	switch r.Type {
	case model.ResponseRowDescription:
		fields := make([]pgproto3.FieldDescription, len(r.Columns))
		for i, c := range r.Columns {
			fields[i] = pgproto3.FieldDescription{
				Name:                 []byte(c.Name),
				TableOID:             c.TableOID,
				TableAttributeNumber: c.ColumnNumber,
				DataTypeOID:          c.TypeOID,
				DataTypeSize:         c.TypeSize,
				TypeModifier:         c.TypeModifier,
				Format:               c.Format,
			}
		}
		return &pgproto3.RowDescription{Fields: fields}
	case model.ResponseDataRow:
		vals := make([][]byte, len(r.Values))
		for i, v := range r.Values {
			if !v.Null {
				vals[i] = append([]byte{}, v.Bytes...)
			}
		}
		return &pgproto3.DataRow{Values: vals}
	case model.ResponseCommandComplete:
		return &pgproto3.CommandComplete{CommandTag: []byte(r.Tag)}
	case model.ResponseEmptyQueryResponse:
		return &pgproto3.EmptyQueryResponse{}
	case model.ResponseErrorResponse:
		e := errorResponse(r.Fields)
		return &e
	case model.ResponseNoticeResponse:
		e := pgproto3.NoticeResponse(errorResponse(r.Fields))
		return &e
	case model.ResponseParameterStatus:
		return &pgproto3.ParameterStatus{Name: r.Name, Value: r.Value}
	case model.ResponseReadyForQuery:
		tx := byte('I')
		if r.TxStatus != "" {
			tx = r.TxStatus[0]
		}
		return &pgproto3.ReadyForQuery{TxStatus: tx}
	default:
		return nil
	}
}

func errorResponse(f map[string]string) pgproto3.ErrorResponse {
	num := func(code string) int32 {
		n, _ := strconv.ParseInt(f[code], 10, 32)
		return int32(n)
	}
	e := pgproto3.ErrorResponse{
		Severity:            f["S"],
		SeverityUnlocalized: f["V"],
		Code:                f["C"],
		Message:             f["M"],
		Detail:              f["D"],
		Hint:                f["H"],
		Position:            num("P"),
		InternalPosition:    num("p"),
		InternalQuery:       f["q"],
		Where:               f["W"],
		SchemaName:          f["s"],
		TableName:           f["t"],
		ColumnName:          f["c"],
		DataTypeName:        f["d"],
		ConstraintName:      f["n"],
		File:                f["F"],
		Line:                num("L"),
		Routine:             f["R"],
	}
	known := "SVCMDHPpqWstcdnFLR"
	for k, v := range f {
		if len(k) == 1 && !containsByte(known, k[0]) {
			if e.UnknownFields == nil {
				e.UnknownFields = map[byte]string{}
			}
			e.UnknownFields[k[0]] = v
		}
	}
	return e
}

func containsByte(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}
