package pgwire

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driving"
)

// Startcodes der ersten Client-Nachricht (PGWire).
const (
	codeProtocol30    = 196608   // 3.0
	codeCancelRequest = 80877102 // 1234.5678
	codeSSLRequest    = 80877103 // 1234.5679
	codeGSSEncRequest = 80877104 // 1234.5680
	majorSpezial      = 1234     // Hauptnummer der Sonderanfragen

	// maxStartLaenge begrenzt die Länge einer Startnachricht (SPEC-045); eine
	// längere erste Nachricht ist keine PGWire-Startnachricht.
	maxStartLaenge = 10000
)

// Server nimmt PostgreSQL-Clients an und übersetzt ihre Nachrichten in Aufrufe
// des Record- oder des Replay-Use-Cases (ARC-006). NewRecordServer und
// NewReplayServer setzen genau einen der beiden.
type Server struct {
	recorder driving.Recorder
	replayer driving.Replayer
	log      *slog.Logger

	wg sync.WaitGroup

	mu        sync.Mutex
	firstCode string
}

// NewRecordServer liefert einen Server für den Record-Modus.
func NewRecordServer(r driving.Recorder, log *slog.Logger) *Server {
	return &Server{recorder: r, log: log}
}

// NewReplayServer liefert einen Server für den Replay-Modus.
func NewReplayServer(r driving.Replayer, log *slog.Logger) *Server {
	return &Server{replayer: r, log: log}
}

// Listen öffnet den TCP-Endpunkt; ein nicht zu öffnender Port ist PGR-E4001.
func Listen(address string) (net.Listener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, model.Errorf(model.CodeListen, err, "Adresse %s nicht nutzbar", address)
	}
	return l, nil
}

// FirstErrorCode liefert den Meldungscode des ersten Verbindungsfehlers des
// Laufs oder "" (LH-FA-13.b).
func (s *Server) FirstErrorCode() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.firstCode
}

// Serve nimmt Verbindungen an, bis der Listener geschlossen wird. Endet ctx,
// endet jede Verbindung nach ihrer laufenden Interaktion; Serve kehrt zurück,
// wenn alle Verbindungen beendet sind.
func (s *Server) Serve(ctx context.Context, l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.wg.Wait()
				return
			}
			s.log.Error("Verbindung nicht anzunehmen", "code", model.CodeNetwork, "error", err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
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

	// Endet ctx im Verbindungsaufbau, bricht das Lesen der Startnachricht ab.
	// Danach wacht sitzung selbst über ctx.
	fertig := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.SetReadDeadline(time.Now())
		case <-fertig:
		}
	}()

	br := bufio.NewReader(conn)
	be := pgproto3.NewBackend(br, conn)

	startup, ok := s.startup(conn, br, be)
	if !ok {
		close(fertig)
		return
	}

	// Der Verbindungsaufbau zum Upstream läuft auch bei endendem ctx zu Ende;
	// er ist Teil der laufenden Interaktion.
	id, responses, err := s.open(context.WithoutCancel(ctx), startup.Parameters)
	close(fertig)
	if err != nil {
		s.fail(be, err)
		return
	}
	if id == 0 {
		// Der Server hat den Aufbau mit einer Fehlerantwort beendet; sie geht
		// ohne AuthenticationOk an den Client.
		if err := s.send(be, responses); err != nil {
			s.sendFailed(err)
		}
		return
	}
	be.Send(&pgproto3.AuthenticationOk{})
	if err := s.send(be, responses); err != nil {
		s.sendFailed(err)
		s.close(ctx, id, model.EndNormal)
		return
	}
	s.sitzung(ctx, be, id)
}

// eingang ist eine gelesene Client-Nachricht, abgebildet, bevor der Leser die
// nächste liest (pgproto3 überschreibt die vorige): genau eines von query,
// terminate, extended oder fremd ist belegt, oder err.
type eingang struct {
	query     *string
	terminate bool
	extended  *model.ClientMessage
	fremd     error
	err       error
}

// ausgang sind Server-Nachrichten einer Extended-Interaktion vom Upstream.
type ausgang struct {
	responses []model.Response
	err       error
}

// sitzung verarbeitet die Client-Nachrichten einer Session bis zu ihrem Ende.
//
// Ein Leser liest je Anforderung genau eine Client-Nachricht; nach einem Sync
// fordert sitzung die nächste erst an, wenn das ReadyForQuery der Interaktion
// an den Client gegangen ist (LH-FA-18.a). Hat eine Extended-Interaktion eine
// Gruppe an den Upstream gesendet, wartet AwaitServer nebenher auf
// Server-Nachrichten, bis ihr ReadyForQuery eintrifft; Client- und
// Server-Nachrichten gehen in der Reihenfolge, in der sitzung sie entgegennimmt,
// an den Record-Use-Case.
//
// Endet ctx, endet die Session, sobald keine Extended-Interaktion läuft
// (LH-FA-13.a). Endet die Client-Verbindung oder kommt Terminate, während eine
// läuft, ist das PGR-E4003 (LH-FA-18.a).
func (s *Server) sitzung(ctx context.Context, be *pgproto3.Backend, id model.SessionID) {
	anfrage := make(chan struct{})
	clientCh := make(chan eingang, 1)
	defer close(anfrage)
	go func() {
		for range anfrage {
			msg, err := be.Receive()
			clientCh <- lese(msg, err)
		}
	}()
	serverCh := make(chan ausgang, 1)

	var (
		clientAngefordert, serverAngefordert bool
		// offen: eine Extended-Interaktion läuft; gesendet: eine ihrer Gruppen
		// ist an den Upstream gegangen; sync: ihr Sync ist gesendet.
		offen, gesendet, sync bool
	)
	for {
		if !offen && ctx.Err() != nil {
			s.close(ctx, id, model.EndNormal)
			return
		}
		if !clientAngefordert && !sync {
			anfrage <- struct{}{}
			clientAngefordert = true
		}
		if gesendet && !serverAngefordert {
			serverAngefordert = true
			go func() {
				out, err := s.recorder.AwaitServer(context.WithoutCancel(ctx), id)
				serverCh <- ausgang{out, err}
			}()
		}
		var stopp <-chan struct{}
		if !offen {
			stopp = ctx.Done()
		}

		select {
		case <-stopp:
			s.close(ctx, id, model.EndNormal)
			return

		case a := <-serverCh:
			serverAngefordert = false
			if a.err != nil {
				s.fail(be, a.err)
				s.close(ctx, id, endFor(a.err))
				return
			}
			abgeschlossen := false
			for _, r := range a.responses {
				ende, err := s.recorder.ServerMessage(ctx, id, r)
				if err != nil {
					s.fail(be, err)
					s.close(ctx, id, endFor(err))
					return
				}
				abgeschlossen = abgeschlossen || ende
			}
			if err := s.send(be, a.responses); err != nil {
				s.sendFailed(err)
				end := endForSend(err)
				if end == model.EndLost && !abgeschlossen {
					// Die laufende Interaktion verwirft CloseSession ohnehin;
					// die vorige hat ihr ReadyForQuery erreicht.
					end = model.EndNormal
				}
				s.close(ctx, id, end)
				return
			}
			if abgeschlossen {
				offen, gesendet, sync = false, false, false
			}

		case e := <-clientCh:
			clientAngefordert = false
			switch {
			case e.err != nil:
				if !verbindungsende(e.err) {
					s.fail(be, model.Errorf(model.CodeUnsupported, e.err, "Client-Nachricht nicht lesbar"))
					s.close(ctx, id, model.EndUnsupported)
					return
				}
				// Ende nach einem ReadyForQuery, mit oder ohne Terminate, ist
				// regulär (LH-FA-02.b).
				if offen {
					s.note(model.Errorf(model.CodeConnectionLost, e.err, "Client-Verbindung vor dem ReadyForQuery der Extended-Interaktion beendet"))
				}
				s.close(ctx, id, model.EndNormal)
				return
			case e.terminate:
				if offen {
					s.note(model.Errorf(model.CodeConnectionLost, nil, "Terminate vor dem ReadyForQuery der Extended-Interaktion"))
				}
				s.close(ctx, id, model.EndNormal)
				return
			case e.query != nil:
				out, err := s.query(ctx, id, *e.query)
				if err != nil {
					s.fail(be, err)
					s.close(ctx, id, endFor(err))
					return
				}
				if err := s.send(be, out); err != nil {
					s.sendFailed(err)
					s.close(ctx, id, endForSend(err))
					return
				}
			case e.extended != nil && s.recorder != nil:
				if err := s.recorder.ClientMessage(ctx, id, *e.extended); err != nil {
					s.fail(be, err)
					s.close(ctx, id, endFor(err))
					return
				}
				offen = true
				switch e.extended.Type {
				case model.ClientSync:
					gesendet, sync = true, true
				case model.ClientFlush:
					gesendet = true
				}
			default:
				err := e.fremd
				if err == nil {
					err = model.Errorf(model.CodeUnsupported, nil, "Client-Nachricht %s wird im Replay nicht unterstützt", e.extended.Type)
				}
				s.fail(be, err)
				s.close(ctx, id, model.EndUnsupported)
				return
			}
		}
	}
}

// lese bildet eine gelesene Client-Nachricht ab. Extended-Nachrichten werden zu
// Client-Nachrichten des Domain Models; eine Zielart von Describe oder Close
// außer 'S' und 'P' und jede andere Nachricht sind PGR-E6001 (LH-FA-05.e,
// LH-FA-18.a).
func lese(msg pgproto3.FrontendMessage, err error) eingang {
	if err != nil {
		return eingang{err: err}
	}
	switch m := msg.(type) {
	case *pgproto3.Query:
		sql := m.String
		return eingang{query: &sql}
	case *pgproto3.Terminate:
		return eingang{terminate: true}
	}
	cm, err := toClientMessage(msg)
	if err != nil {
		return eingang{fremd: err}
	}
	return eingang{extended: &cm}
}

// toClientMessage belegt je Typ genau dessen Felder (SPEC-041) und kopiert jede
// Liste und jeden Wert; eine leere Liste wird nil, wie sie der Leser der
// Aufzeichnung liefert. Ein Parameter NULL wird Value.Null, jeder andere seine
// Bytes, auch leer.
func toClientMessage(msg pgproto3.FrontendMessage) (model.ClientMessage, error) {
	switch m := msg.(type) {
	case *pgproto3.Parse:
		return model.ClientMessage{Type: model.ClientParse, Statement: m.Name, SQL: m.Query, ParamTypes: kopieOderNil(m.ParameterOIDs)}, nil
	case *pgproto3.Bind:
		var params []model.Value
		for _, p := range m.Parameters {
			if p == nil {
				params = append(params, model.Value{Null: true})
			} else {
				params = append(params, model.Value{Bytes: append([]byte{}, p...)})
			}
		}
		return model.ClientMessage{
			Type: model.ClientBind, Portal: m.DestinationPortal, Statement: m.PreparedStatement,
			ParamFormats: kopieOderNil(m.ParameterFormatCodes), Params: params, ResultFormats: kopieOderNil(m.ResultFormatCodes),
		}, nil
	case *pgproto3.Describe:
		ziel, err := zielart(m.ObjectType, "Describe")
		return model.ClientMessage{Type: model.ClientDescribe, Target: ziel, Name: m.Name}, err
	case *pgproto3.Close:
		ziel, err := zielart(m.ObjectType, "Close")
		return model.ClientMessage{Type: model.ClientClose, Target: ziel, Name: m.Name}, err
	case *pgproto3.Execute:
		return model.ClientMessage{Type: model.ClientExecute, Portal: m.Portal, MaxRows: m.MaxRows}, nil
	case *pgproto3.Flush:
		return model.ClientMessage{Type: model.ClientFlush}, nil
	case *pgproto3.Sync:
		return model.ClientMessage{Type: model.ClientSync}, nil
	default:
		return model.ClientMessage{}, model.Errorf(model.CodeUnsupported, nil, "Client-Nachricht %T wird nicht unterstützt", m)
	}
}

func zielart(b byte, typ string) (model.Target, error) {
	switch b {
	case 'S':
		return model.TargetStatement, nil
	case 'P':
		return model.TargetPortal, nil
	default:
		return "", model.Errorf(model.CodeUnsupported, nil, "%s mit Zielart %q wird nicht unterstützt", typ, b)
	}
}

// kopieOderNil kopiert eine Liste; eine leere wird nil.
func kopieOderNil[T any](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	return append([]T(nil), s...)
}

// startup liest die erste Client-Nachricht. SSL- und GSS-Anfragen beantwortet
// er mit „N“ (LH-FA-05.c), ein CancelRequest schließt die Verbindung
// (PGR-W3001). Jede andere Protokollversion als 3.0 ist PGR-E6002, eine
// unbekannte Sonderanfrage (Hauptnummer 1234) PGR-E6001. Eine erste Nachricht,
// die keine PGWire-Startnachricht sein kann, schließt die Verbindung ohne Antwort
// und nur mit der Warnung PGR-W3003; eine Verbindung ohne erste Nachricht endet
// still (LH-FA-05.e). Der Startcode wird vor pgproto3 gelesen, weil die
// Bibliothek unbekannte Codes ohne Antwort ablehnt.
func (s *Server) startup(conn net.Conn, br *bufio.Reader, be *pgproto3.Backend) (*pgproto3.StartupMessage, bool) {
	for {
		kopf, err := br.Peek(8)
		if err != nil {
			// Eine Verbindung ohne Startnachricht, zum Beispiel eine TCP-Probe.
			s.log.Debug("Verbindung ohne Startnachricht beendet", "error", err.Error())
			return nil, false
		}
		laenge := binary.BigEndian.Uint32(kopf[0:4])
		code := binary.BigEndian.Uint32(kopf[4:8])
		switch {
		case laenge < 8 || laenge > maxStartLaenge:
			s.log.Warn("erste Nachricht ist keine PGWire-Startnachricht", "code", model.CodeForeignProtocol, "remote", conn.RemoteAddr().String())
			return nil, false
		case code == codeSSLRequest || code == codeGSSEncRequest || code == codeCancelRequest || code == codeProtocol30:
		case code>>16 == majorSpezial:
			s.fail(be, model.Errorf(model.CodeUnsupported, nil, "unbekannte Sonderanfrage (Code %d)", code))
			return nil, false
		default:
			s.fail(be, model.Errorf(model.CodeProtocolVersion, nil, "Protokollversion %d.%d wird nicht unterstützt", code>>16, code&0xffff))
			return nil, false
		}

		msg, err := be.ReceiveStartupMessage()
		if err != nil {
			s.fail(be, model.Errorf(model.CodeUnsupported, err, "Startnachricht nicht lesbar"))
			return nil, false
		}
		switch m := msg.(type) {
		case *pgproto3.SSLRequest, *pgproto3.GSSEncRequest:
			if _, err := conn.Write([]byte{'N'}); err != nil {
				return nil, false
			}
		case *pgproto3.CancelRequest:
			s.log.Warn("CancelRequest empfangen und nicht weitergeleitet", "code", model.CodeCancelRequest)
			return nil, false
		case *pgproto3.StartupMessage:
			return m, true
		default:
			s.fail(be, model.Errorf(model.CodeUnsupported, nil, "Startnachricht %T wird nicht unterstützt", m))
			return nil, false
		}
	}
}

// verbindungsende meldet, ob ein Lesefehler das Ende der Client-Verbindung ist
// und nicht eine unlesbare Nachricht.
func verbindungsende(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, net.ErrClosed) || errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, os.ErrDeadlineExceeded)
}

// abbildungsfehler kennzeichnet eine Antwort, die nicht auf eine PGWire-Nachricht
// abbildbar ist; send liefert ihn statt eines Schreibfehlers.
type abbildungsfehler struct{ err error }

func (a abbildungsfehler) Error() string { return a.err.Error() }
func (a abbildungsfehler) Unwrap() error { return a.err }

func endForSend(err error) model.SessionEnd {
	var a abbildungsfehler
	if errors.As(err, &a) {
		return model.EndUnsupported
	}
	return model.EndLost
}

// sendFailed merkt sich einen gescheiterten Versand: Eine nicht abbildbare
// Antwort ist ein interner Fehler, ein Schreibfehler ein unerwartetes Ende der
// Client-Verbindung (PGR-E4003).
func (s *Server) sendFailed(err error) {
	var a abbildungsfehler
	if errors.As(err, &a) {
		s.note(a.err)
		return
	}
	s.note(model.Errorf(model.CodeConnectionLost, err, "Antwort nicht an den Client zu senden"))
}

func endFor(err error) model.SessionEnd {
	var me *model.Error
	if errors.As(err, &me) && me.Code == model.CodeUnsupported {
		return model.EndUnsupported
	}
	return model.EndNormal
}

func (s *Server) open(ctx context.Context, startup map[string]string) (model.SessionID, []model.Response, error) {
	if s.replayer != nil {
		id, out := s.replayer.OpenConnection(ctx)
		return id, out, nil
	}
	return s.recorder.OpenSession(ctx, startup)
}

func (s *Server) query(ctx context.Context, id model.SessionID, sql string) ([]model.Response, error) {
	if s.replayer != nil {
		return s.replayer.Query(ctx, id, sql)
	}
	return s.recorder.Query(ctx, id, sql)
}

// close beendet die Session. Im Replay ist eine unverbrauchte Session eine
// Warnung (PGR-W2001), kein Verbindungsfehler.
func (s *Server) close(ctx context.Context, id model.SessionID, end model.SessionEnd) {
	if s.replayer != nil {
		if w := s.replayer.CloseConnection(context.WithoutCancel(ctx), id); w != nil {
			s.log.Warn(w.Msg, "code", w.Code)
		}
		return
	}
	if err := s.recorder.CloseSession(context.WithoutCancel(ctx), id, end); err != nil {
		s.note(err)
	}
}

// fail merkt sich den Fehler und stellt ihn dem Client als FATAL-ErrorResponse
// mit dem Meldungscode im Text zu (LH-FA-13.b).
func (s *Server) fail(be *pgproto3.Backend, err error) {
	code := s.note(err)
	be.Send(&pgproto3.ErrorResponse{Severity: "FATAL", SeverityUnlocalized: "FATAL", Code: sqlstate(code), Message: err.Error()})
	_ = be.Flush()
}

// note protokolliert einen Verbindungsfehler, merkt sich den ersten und
// liefert seinen Meldungscode.
func (s *Server) note(err error) string {
	code := model.CodeInternal
	var me *model.Error
	if errors.As(err, &me) {
		code = me.Code
	}
	s.mu.Lock()
	if s.firstCode == "" {
		s.firstCode = code
	}
	s.mu.Unlock()
	s.log.Error("Fehler", "code", code, "error", err.Error())
	return code
}

// sqlstate wählt den SQLSTATE der Fehlerantwort nach der Klasse des
// Meldungscodes: nicht unterstützt 0A000, Netzwerk 08006, sonst XX000.
func sqlstate(code string) string {
	switch {
	case len(code) >= 6 && code[5] == '6':
		return "0A000"
	case len(code) >= 6 && code[5] == '4':
		return "08006"
	default:
		return "XX000"
	}
}

func (s *Server) send(be *pgproto3.Backend, responses []model.Response) error {
	for _, r := range responses {
		msg, err := toMessage(r)
		if err != nil {
			return abbildungsfehler{err}
		}
		be.Send(msg)
	}
	return be.Flush()
}

func toMessage(r model.Response) (pgproto3.BackendMessage, error) {
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
		return &pgproto3.RowDescription{Fields: fields}, nil
	case model.ResponseDataRow:
		vals := make([][]byte, len(r.Values))
		for i, v := range r.Values {
			if !v.Null {
				vals[i] = append([]byte{}, v.Bytes...)
			}
		}
		return &pgproto3.DataRow{Values: vals}, nil
	case model.ResponseCommandComplete:
		return &pgproto3.CommandComplete{CommandTag: []byte(r.Tag)}, nil
	case model.ResponseEmptyQueryResponse:
		return &pgproto3.EmptyQueryResponse{}, nil
	case model.ResponseErrorResponse:
		e := errorResponse(r.Fields)
		return &e, nil
	case model.ResponseNoticeResponse:
		e := pgproto3.NoticeResponse(errorResponse(r.Fields))
		return &e, nil
	case model.ResponseParameterStatus:
		return &pgproto3.ParameterStatus{Name: r.Name, Value: r.Value}, nil
	case model.ResponseParseComplete:
		return &pgproto3.ParseComplete{}, nil
	case model.ResponseBindComplete:
		return &pgproto3.BindComplete{}, nil
	case model.ResponseCloseComplete:
		return &pgproto3.CloseComplete{}, nil
	case model.ResponseParameterDescription:
		return &pgproto3.ParameterDescription{ParameterOIDs: append([]uint32{}, r.ParamTypes...)}, nil
	case model.ResponseNoData:
		return &pgproto3.NoData{}, nil
	case model.ResponsePortalSuspended:
		return &pgproto3.PortalSuspended{}, nil
	case model.ResponseReadyForQuery:
		if len(r.TxStatus) != 1 {
			return nil, model.Errorf(model.CodeInternal, nil, "ReadyForQuery ohne Transaktionsstatus")
		}
		return &pgproto3.ReadyForQuery{TxStatus: r.TxStatus[0]}, nil
	default:
		return nil, model.Errorf(model.CodeInternal, nil, "Antworttyp %q ohne Abbildung", r.Type)
	}
}
