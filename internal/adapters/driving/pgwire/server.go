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
	"sync/atomic"
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
	// Danach wachen recordSitzung und replaySitzung selbst über ctx.
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
	sendErr := s.send(be, responses)
	if s.replayer != nil {
		if sendErr != nil {
			s.sendFailed(sendErr)
			s.closeReplay(ctx, id)
			return
		}
		s.replaySitzung(ctx, conn, be, id)
		return
	}
	if sendErr != nil {
		s.closeRecord(ctx, id, endBeiSchreibfehler(s, sendErr))
		return
	}
	s.recordSitzung(ctx, conn, be, id)
}

// eingang ist eine gelesene Client-Nachricht, abgebildet, bevor der Leser die
// nächste liest (pgproto3 überschreibt die vorige): genau eines von query,
// terminate, extended oder fremd ist belegt.
type eingang struct {
	query     *string
	terminate bool
	extended  *model.ClientMessage
	fremd     error
}

// replaySitzung beantwortet die Client-Nachrichten einer Replay-Session
// nacheinander: Anfragen und Extended-Nachrichten gehen an den Replay-Use-Case,
// und was er liefert, geht an den Client.
//
// Endet ctx, weckt ein Wächter das Lesen. Vor jedem weiteren Lesen fragt die
// Sitzung dann den Use Case (Shutdown): Läuft keine Interaktion, endet die
// Session, ohne die nächste Nachricht zu lesen; läuft eine Extended-Interaktion,
// liest die Sitzung weiter und beantwortet ihre Nachrichten bis zu ihrem Sync
// (LH-FA-13.a). Eine schon gelesene Nachricht wird vorher noch beantwortet.
// Das Warten ist nicht begrenzt; die Frist --shutdown-timeout (LH-FA-13.a) ist
// hier nicht verdrahtet. Ein Schließen von conn beendet das Lesen wie jedes
// Verbindungsende, und die Session endet regulär: Wer die Frist so durchsetzt,
// merkt PGR-E4006 für eine unvollständige Interaktion selbst.
func (s *Server) replaySitzung(ctx context.Context, conn net.Conn, be *pgproto3.Backend, id model.SessionID) {
	fertig := make(chan struct{})
	defer close(fertig)
	// geweckt ist geschlossen, sobald der Wächter die Lesefrist gesetzt hat;
	// erst danach setzt die Sitzung sie zurück.
	geweckt := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.SetReadDeadline(time.Now())
			close(geweckt)
		case <-fertig:
		}
	}()
	defer s.closeReplay(ctx, id)

	herunterfahren := false
	for {
		if ctx.Err() != nil {
			if s.replayer.Shutdown(ctx, id) {
				return
			}
			if !herunterfahren {
				herunterfahren = true
				<-geweckt
				_ = conn.SetReadDeadline(time.Time{})
			}
		}
		msg, err := be.Receive()
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) && ctx.Err() != nil && !herunterfahren {
				// Geweckt durch den Wächter: oben entscheidet der Use Case.
				continue
			}
			if verbindungsende(err) {
				// Ein Ende der Client-Verbindung ist im Replay regulär, auch
				// mitten in einer Extended-Interaktion; was unverbraucht
				// bleibt, meldet closeReplay (LH-FA-03.b).
				return
			}
			s.fail(be, model.Errorf(model.CodeUnsupported, err, "Client-Nachricht nicht lesbar"))
			return
		}
		e := lese(msg)
		switch {
		case e.query != nil:
			out, err := s.replayer.Query(ctx, id, *e.query)
			if err != nil {
				s.fail(be, err)
				return
			}
			if err := s.send(be, out); err != nil {
				s.sendFailed(err)
				return
			}
			s.replayer.Sent(ctx, id)
		case e.terminate:
			return
		case e.fremd != nil:
			s.fail(be, e.fremd)
			return
		default:
			out, err := s.replayer.ClientMessage(ctx, id, *e.extended)
			if err != nil {
				s.fail(be, err)
				return
			}
			if len(out) == 0 {
				continue
			}
			if err := s.send(be, out); err != nil {
				s.sendFailed(err)
				return
			}
			s.replayer.Sent(ctx, id)
		}
	}
}

// recordSitzung vermittelt eine Record-Session in zwei Richtungen, die
// unabhängig voneinander blockieren (LH-FA-18.a): clientRichtung liest
// Client-Nachrichten und übergibt sie dem Record-Use-Case, serverRichtung holt
// Server-Nachrichten über AwaitServer und schreibt sie an den Client. Beide
// schreiben an den Client unter schreiben. Den Interaktionszustand und den
// Grund des Session-Endes führt der Use Case; die Sitzung meldet Ereignisse.
// Wer zuerst ein Ende meldet, beendet die Session und schließt die
// Client-Verbindung; das beendet blockiertes Lesen und Schreiben der anderen
// Richtung, CloseSession ihre wartenden Aufrufe des Use Case.
//
// Endet ctx, weckt ein Wächter die Client-Richtung; sie meldet das
// Herunterfahren selbst, bevor sie die nächste Nachricht liest. Eine schon
// gelesene Nachricht verarbeitet sie also vorher (LH-FA-13.a).
func (s *Server) recordSitzung(ctx context.Context, conn net.Conn, be *pgproto3.Backend, id model.SessionID) {
	r := &richtungen{s: s, conn: conn, be: be, id: id, ctx: context.WithoutCancel(ctx),
		weck: make(chan struct{}, 1), ende: make(chan struct{})}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			r.wecke()
		case <-r.ende:
		}
	}()
	go func() {
		defer wg.Done()
		r.serverRichtung()
	}()
	r.clientRichtung(ctx)
	wg.Wait()
}

// richtungen hält, was beide Richtungen einer Record-Session teilen.
type richtungen struct {
	s    *Server
	conn net.Conn
	be   *pgproto3.Backend
	id   model.SessionID
	ctx  context.Context

	schreiben sync.Mutex
	beendet   atomic.Bool
	einmal    sync.Once
	// ende ist geschlossen, sobald die Session beendet ist.
	ende chan struct{}

	// frist ordnet das Wecken (Signal in weck, Lesefrist sofort) gegen das
	// Zurücksetzen der Lesefrist: wer zurücksetzt, sieht unter derselben Sperre,
	// ob seither geweckt wurde, und kein Wecken geht verloren.
	frist sync.Mutex
	weck  chan struct{}
}

// meldeFrist begrenzt das Schreiben der Fehlerantwort beim Ende einer Session,
// auch wenn der Client nicht liest.
const meldeFrist = time.Second

// wecke bricht das laufende oder nächste Lesen vom Client ab und hinterlegt
// ein Signal in weck.
func (r *richtungen) wecke() {
	r.frist.Lock()
	defer r.frist.Unlock()
	select {
	case r.weck <- struct{}{}:
	default:
	}
	_ = r.conn.SetReadDeadline(time.Now())
}

// weiterlesen setzt die Lesefrist zurück, wenn seit dem letzten Wecken keines
// hinzukam; sonst lässt es sie stehen und meldet das Wecken.
func (r *richtungen) weiterlesen() (geweckt bool) {
	r.frist.Lock()
	defer r.frist.Unlock()
	select {
	case <-r.weck:
		return true
	default:
	}
	_ = r.conn.SetReadDeadline(time.Time{})
	return false
}

// beende meldet das Ende der Session genau einmal. Mit meldung stellt es dem
// Client vorher eine Fehlerantwort zu und merkt sie (LH-FA-13.b); das Schreiben
// dauert höchstens meldeFrist, auch wenn der Client nicht liest. Danach schließt
// es die Client-Verbindung.
func (r *richtungen) beende(end model.SessionEnd, meldung error) {
	r.einmal.Do(func() {
		r.beendet.Store(true)
		if meldung != nil {
			_ = r.conn.SetWriteDeadline(time.Now().Add(meldeFrist))
			r.schreiben.Lock()
			r.s.fail(r.be, meldung)
			r.schreiben.Unlock()
		}
		r.s.closeRecord(r.ctx, r.id, end)
		_ = r.conn.Close()
		close(r.ende)
	})
}

// schreibe schreibt Antworten an den Client und meldet sie danach als
// zugestellt; gibt der Use Case damit das Ende frei, weckt es die
// Client-Richtung. Scheitert das Schreiben, beendet es die Session.
func (r *richtungen) schreibe(rs []model.Response) bool {
	r.schreiben.Lock()
	err := r.s.send(r.be, rs)
	r.schreiben.Unlock()
	if err != nil {
		r.beende(endBeiSchreibfehler(r.s, err), nil)
		return false
	}
	if r.s.recorder.Delivered(r.ctx, r.id) {
		r.wecke()
	}
	return true
}

// fehler beendet die Session nach einem Fehler des Use Case. Nach dem Ende der
// Session ist nichts mehr zu tun; beim Herunterfahren wartet die
// Client-Richtung auf das Ende (wartenAufEnde).
func (r *richtungen) fehler(err error) {
	if errors.Is(err, model.ErrSessionEnded) {
		return
	}
	if errors.Is(err, model.ErrShutdown) {
		r.wartenAufEnde()
		return
	}
	r.beende(model.EndFailed, err)
}

// wartenAufEnde liest nicht weiter und beendet die Session mit EndShutdown,
// sobald der Use Case das Ende freigibt; es fragt nach jedem Wecken erneut.
func (r *richtungen) wartenAufEnde() {
	for {
		if r.beendet.Load() {
			return
		}
		if r.s.recorder.Shutdown(r.ctx, r.id) {
			r.beende(model.EndShutdown, nil)
			return
		}
		select {
		case <-r.weck:
		case <-r.ende:
			return
		}
	}
}

func (r *richtungen) serverRichtung() {
	for {
		rs, err := r.s.recorder.AwaitServer(r.ctx, r.id)
		if err != nil {
			if !errors.Is(err, model.ErrSessionEnded) {
				r.beende(model.EndFailed, err)
			}
			return
		}
		if !r.schreibe(rs) {
			return
		}
	}
}

// clientRichtung liest Client-Nachrichten, bis die Session endet. Ist ctx
// beendet, meldet sie vor jedem Lesen das Herunterfahren; gibt der Use Case
// das Ende frei, endet die Session mit EndShutdown. Eine Lesefrist aus dem
// Wecken setzt sie über weiterlesen zurück und fragt erneut.
func (r *richtungen) clientRichtung(ctx context.Context) {
	for {
		if r.beendet.Load() {
			return
		}
		if ctx.Err() != nil && r.s.recorder.Shutdown(r.ctx, r.id) {
			r.beende(model.EndShutdown, nil)
			return
		}
		msg, err := r.be.Receive()
		if err != nil {
			switch {
			case r.beendet.Load():
				return
			case errors.Is(err, os.ErrDeadlineExceeded):
				r.weiterlesen()
				continue
			case verbindungsende(err):
				r.beende(model.EndClosed, nil)
			default:
				r.beende(model.EndUnsupported, model.Errorf(model.CodeUnsupported, err, "Client-Nachricht nicht lesbar"))
			}
			return
		}
		e := lese(msg)
		switch {
		case e.query != nil:
			out, err := r.s.recorder.Query(r.ctx, r.id, *e.query)
			if err != nil {
				r.fehler(err)
				return
			}
			if !r.schreibe(out) {
				return
			}
		case e.terminate:
			r.beende(model.EndTerminate, nil)
			return
		case e.extended != nil:
			if err := r.s.recorder.ClientMessage(r.ctx, r.id, *e.extended); err != nil {
				r.fehler(err)
				return
			}
		default:
			r.beende(model.EndUnsupported, e.fremd)
			return
		}
	}
}

// lese bildet eine gelesene Client-Nachricht ab. Extended-Nachrichten werden zu
// Client-Nachrichten des Domain Models; eine Zielart von Describe oder Close
// außer 'S' und 'P' und jede andere Nachricht sind PGR-E6001 (LH-FA-05.e,
// LH-FA-18.a).
func lese(msg pgproto3.FrontendMessage) eingang {
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

// endBeiSchreibfehler meldet das Ereignis eines gescheiterten Versands an den
// Client: Eine nicht abbildbare Antwort ist nicht unterstützt und wird als
// interner Fehler gemerkt; ein Schreibfehler ist EndWriteFailed, und seine
// Einstufung trifft der Use Case.
func endBeiSchreibfehler(s *Server, err error) model.SessionEnd {
	var a abbildungsfehler
	if errors.As(err, &a) {
		s.note(a.err)
		return model.EndUnsupported
	}
	return model.EndWriteFailed
}

// sendFailed merkt sich einen gescheiterten Versand im Replay und im
// Verbindungsaufbau ohne Session: Eine nicht abbildbare Antwort ist ein
// interner Fehler, ein Schreibfehler ein unerwartetes Ende der
// Client-Verbindung (PGR-E4003). Im Record leitet der Use Case die Einstufung
// aus EndWriteFailed ab (endBeiSchreibfehler).
func (s *Server) sendFailed(err error) {
	var a abbildungsfehler
	if errors.As(err, &a) {
		s.note(a.err)
		return
	}
	s.note(model.Errorf(model.CodeConnectionLost, err, "Antwort nicht an den Client zu senden"))
}

func (s *Server) open(ctx context.Context, startup map[string]string) (model.SessionID, []model.Response, error) {
	if s.replayer != nil {
		id, out := s.replayer.OpenConnection(ctx)
		return id, out, nil
	}
	return s.recorder.OpenSession(ctx, startup)
}

// closeReplay beendet eine Replay-Verbindung und meldet, was der Use Case
// über nicht verbrauchte Interaktionen liefert: eine Warnung (PGR-W2001) im
// Log, einen Fehler (PGR-E5002) als Verbindungsfehler, der nach einem
// Fehler, der die Verbindung beendet hat, gemerkt und dem Client nicht
// zugestellt wird (LH-FA-03.b, LH-FA-13.b).
func (s *Server) closeReplay(ctx context.Context, id model.SessionID) {
	w, err := s.replayer.CloseConnection(context.WithoutCancel(ctx), id)
	if err != nil {
		s.note(err)
		return
	}
	if w != nil {
		s.log.Warn(w.Msg, "code", w.Code)
	}
}

// closeRecord meldet dem Record-Use-Case das Ende der Session und merkt den
// Fehler, den er daraus ableitet.
func (s *Server) closeRecord(ctx context.Context, id model.SessionID, end model.SessionEnd) {
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
