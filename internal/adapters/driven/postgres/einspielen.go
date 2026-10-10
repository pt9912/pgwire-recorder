package postgres

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// Einspielziel verbindet sich beim Einspielen je Session als Client mit dem
// Server unter Address, ohne TLS, und meldet sich mit Password an, wenn der
// Server ein Passwort verlangt; ohne Password (leer) gibt es keines
// (LH-FA-20.a *Anmeldung*, *Passwort*).
type Einspielziel struct {
	Address  string
	Password Passwort
	Dialer   net.Dialer
}

var _ driven.Einspielziel = (*Einspielziel)(nil)

// Verbinde baut die Verbindung auf und liest den Aufbau nach LH-FA-20.a
// *Aufbau* und *Anmelde-Nachrichten* (aufbau). Bei jedem Fehler schließt es die
// Verbindung ohne Terminate (*Abbruch im Aufbau*). Endet ctx, schließt es die
// Verbindung oder bricht den Verbindungsversuch ab und liefert einen Fehler.
func (z *Einspielziel) Verbinde(ctx context.Context, startup map[string]string) (driven.EinspielSession, error) {
	conn, err := z.Dialer.DialContext(ctx, "tcp", z.Address)
	if err != nil {
		return nil, model.Errorf(model.CodeUpstream, err, "Server %s nicht erreichbar", z.Address)
	}
	stop := context.AfterFunc(ctx, func() { _ = conn.Close() })
	s, err := aufbau(conn, startup, zugang{passwort: z.Password})
	if !stop() {
		_ = conn.Close()
		return nil, model.Errorf(model.CodeUpstream, ctx.Err(), "Aufbau zu %s abgebrochen", z.Address)
	}
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return s, nil
}

// aufbau sendet das Startup, meldet sich an und liest bis zum ersten
// ReadyForQuery. Es liest jede Nachricht selbst, damit der Code einer
// Anmelde-Nachricht auch dort gilt, wo pgproto3 ihn nicht dekodiert. Der Puffer
// des Lesens geht mit der Session weiter, sodass nichts verloren geht, was der
// Server hinter dem ReadyForQuery sendet.
func aufbau(conn net.Conn, startup map[string]string, z zugang) (*einspielSession, error) {
	start, err := (&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: startup}).Encode(nil)
	if err != nil {
		return nil, model.Errorf(model.CodeUpstream, err, "Startup nicht zu kodieren")
	}
	if _, err := conn.Write(start); err != nil {
		return nil, model.Errorf(model.CodeUpstream, err, "Startup nicht zu senden")
	}
	r := bufio.NewReader(conn)
	login := neueAnmeldung(startup, z)
	for {
		typ, rumpf, err := liesRoh(r)
		if err != nil {
			return nil, model.Errorf(model.CodeUpstream, err, "Verbindung im Aufbau beendet oder Nachricht nicht lesbar")
		}
		switch typ {
		case 'R':
			if err := login.bearbeite(conn, rumpf); err != nil {
				return nil, err
			}
		case 'E':
			return nil, fehlerImAufbau(rumpf)
		case 'Z':
			if !login.ok {
				return nil, model.Errorf(model.CodeUpstream, nil, "ReadyForQuery vor AuthenticationOk")
			}
			if err := (&pgproto3.ReadyForQuery{}).Decode(rumpf); err != nil {
				return nil, model.Errorf(model.CodeUpstream, err, "ReadyForQuery im Aufbau nicht lesbar")
			}
			return neueEinspielSession(conn, r), nil
		case 'K', 'S', 'N', 'A':
			if err := verworfen(typ).Decode(rumpf); err != nil {
				return nil, model.Errorf(model.CodeUpstream, err, "Nachricht %q im Aufbau nicht lesbar", typ)
			}
		default:
			return nil, model.Errorf(model.CodeUpstream, nil, "Nachricht %q ist im Aufbau nicht vorgesehen", typ)
		}
	}
}

// liesRoh liest eine Server-Nachricht: Typ, Länge und Rumpf. Eine Länge unter 4
// ist nicht lesbar.
func liesRoh(r io.Reader) (byte, []byte, error) {
	var kopf [5]byte
	if _, err := io.ReadFull(r, kopf[:]); err != nil {
		return 0, nil, err
	}
	laenge := int(int32(binary.BigEndian.Uint32(kopf[1:])))
	if laenge < 4 {
		return 0, nil, errors.New("Länge einer Nachricht unter 4")
	}
	rumpf := make([]byte, laenge-4)
	if _, err := io.ReadFull(r, rumpf); err != nil {
		return 0, nil, err
	}
	return kopf[0], rumpf, nil
}

// fehlerImAufbau stuft eine Fehlerantwort im Aufbau ein: SQLSTATE-Klasse 28
// ist PGR-E4005, jede andere PGR-E4002, gleich welcher Schweregrad; eine
// Fehlerantwort, die sich nicht lesen lässt, ist PGR-E4002. Die Meldung nennt
// SQLSTATE und Meldung der Fehlerantwort (LH-FA-20.a *Aufbau*, *Meldungen*).
func fehlerImAufbau(rumpf []byte) error {
	var e pgproto3.ErrorResponse
	if err := e.Decode(rumpf); err != nil {
		return model.Errorf(model.CodeUpstream, err, "Fehlerantwort im Aufbau nicht lesbar")
	}
	code := model.CodeUpstream
	if strings.HasPrefix(e.Code, "28") {
		code = model.CodeLogin
	}
	return model.Errorf(code, nil, "Fehlerantwort im Aufbau %s „%s“", e.Code, e.Message)
}

// verworfen ist die Nachricht, in die eine der im Aufbau verworfenen
// Nachrichten dekodiert wird: BackendKeyData, ParameterStatus,
// NoticeResponse oder NotificationResponse.
func verworfen(typ byte) pgproto3.BackendMessage {
	switch typ {
	case 'K':
		return &pgproto3.BackendKeyData{}
	case 'S':
		return &pgproto3.ParameterStatus{}
	case 'N':
		return &pgproto3.NoticeResponse{}
	}
	return &pgproto3.NotificationResponse{}
}

// einspielSession ist eine aufgebaute Verbindung beim Einspielen. schreiben
// hält das Senden von Anfrage und einer Gruppe und das Terminate von Schliesse
// auseinander. Die Gruppen sendet ein eigener Sender nacheinander aus
// warteschlange, unabhängig vom Lesen in Naechste; zustand schützt
// warteschlange, senderLaeuft, geschlossen und sendefehler. weck weckt den
// wartenden Sender nach einer eingereihten Gruppe, ende beendet ihn mit
// Schliesse.
type einspielSession struct {
	conn      net.Conn
	fe        *pgproto3.Frontend
	schreiben sync.Mutex

	zustand       sync.Mutex
	warteschlange [][]pgproto3.FrontendMessage
	senderLaeuft  bool
	weck          chan struct{}
	ende          chan struct{}
	geschlossen   bool
	sendefehler   error
}

// neueEinspielSession ist die Session auf conn; sie liest aus r, das die
// Bytes hinter dem Aufbau noch trägt, und schreibt auf conn.
func neueEinspielSession(conn net.Conn, r io.Reader) *einspielSession {
	return &einspielSession{
		conn: conn, fe: pgproto3.NewFrontend(r, conn),
		weck: make(chan struct{}, 1), ende: make(chan struct{}),
	}
}

func (s *einspielSession) Anfrage(sql string) error {
	s.schreiben.Lock()
	defer s.schreiben.Unlock()
	s.fe.Send(&pgproto3.Query{String: sql})
	if err := s.fe.Flush(); err != nil {
		return model.Errorf(model.CodeConnectionLost, err, "Anfrage an den Server nicht zu senden")
	}
	return nil
}

// Gruppe bildet die Nachrichten auf PGWire ab und reiht sie für den Sender
// ein; es wartet nicht auf das Senden. Ein früherer Fehler des Sendens ist
// PGR-E4003.
func (s *einspielSession) Gruppe(nachrichten []model.ClientMessage) error {
	if err := s.fehlerDesSendens(); err != nil {
		return err
	}
	msgs := make([]pgproto3.FrontendMessage, 0, len(nachrichten))
	for _, m := range nachrichten {
		msg, err := toFrontendMessage(m)
		if err != nil {
			return err
		}
		msgs = append(msgs, msg)
	}
	s.zustand.Lock()
	if s.geschlossen {
		s.zustand.Unlock()
		return model.Errorf(model.CodeConnectionLost, nil, "Gruppe nach dem Schließen der Verbindung")
	}
	s.warteschlange = append(s.warteschlange, msgs)
	if !s.senderLaeuft {
		s.senderLaeuft = true
		go s.sender()
	}
	s.zustand.Unlock()
	select {
	case s.weck <- struct{}{}:
	default:
	}
	return nil
}

// sender sendet die eingereihten Gruppen nacheinander, bis Schliesse ihn
// beendet oder ein Senden scheitert.
func (s *einspielSession) sender() {
	for {
		s.zustand.Lock()
		var gruppe []pgproto3.FrontendMessage
		if len(s.warteschlange) > 0 {
			gruppe, s.warteschlange = s.warteschlange[0], s.warteschlange[1:]
		}
		s.zustand.Unlock()
		if gruppe == nil {
			select {
			case <-s.weck:
				continue
			case <-s.ende:
				return
			}
		}
		if !s.sendeGruppe(gruppe) {
			return
		}
	}
}

// sendeGruppe schreibt die Nachrichten einer Gruppe. Es liefert falsch, wenn
// der Sender enden soll: nach Schliesse, oder weil das Senden scheiterte. Dann
// ist der Fehler PGR-E4003 für Naechste und die nächste Gruppe gemerkt und die
// Verbindung geschlossen, damit ein wartendes Naechste endet; nach Schliesse
// bleibt ein Fehler ohne Folge.
func (s *einspielSession) sendeGruppe(gruppe []pgproto3.FrontendMessage) bool {
	s.schreiben.Lock()
	defer s.schreiben.Unlock()
	s.zustand.Lock()
	zu := s.geschlossen
	s.zustand.Unlock()
	if zu {
		return false
	}
	for _, m := range gruppe {
		s.fe.Send(m)
	}
	err := s.fe.Flush()
	if err == nil {
		return true
	}
	s.zustand.Lock()
	if !s.geschlossen {
		s.sendefehler = model.Errorf(model.CodeConnectionLost, err, "Gruppe an den Server nicht zu senden")
	}
	s.zustand.Unlock()
	_ = s.conn.Close()
	return false
}

// fehlerDesSendens ist der gemerkte Fehler eines gescheiterten Sendens, sonst nil.
func (s *einspielSession) fehlerDesSendens() error {
	s.zustand.Lock()
	defer s.zustand.Unlock()
	return s.sendefehler
}

func (s *einspielSession) Naechste() (model.Response, error) {
	for {
		if err := s.fehlerDesSendens(); err != nil {
			return model.Response{}, err
		}
		msg, err := s.fe.Receive()
		if err != nil {
			if err := s.fehlerDesSendens(); err != nil {
				return model.Response{}, err
			}
			if istVerbindungsende(err) {
				return model.Response{}, model.Errorf(model.CodeConnectionLost, err, "Verbindung zum Server beendet")
			}
			return model.Response{}, model.Errorf(model.CodeUnsupported, err, "Serverantwort nicht lesbar")
		}
		switch msg.(type) {
		case *pgproto3.CopyInResponse, *pgproto3.CopyOutResponse, *pgproto3.CopyBothResponse:
			return model.Response{}, model.Errorf(model.CodeUnsupported, nil, "Serverantwort %T kann play nicht bedienen", msg)
		}
		if r, err := toResponse(msg); err == nil {
			return r, nil
		}
	}
}

// Schliesse beendet den Sender, sendet Terminate nur, wenn gerade nichts
// sendet, und höchstens terminateFrist lang; danach schließt es die
// Verbindung.
func (s *einspielSession) Schliesse() {
	s.zustand.Lock()
	if !s.geschlossen {
		s.geschlossen = true
		close(s.ende)
	}
	s.zustand.Unlock()
	if s.schreiben.TryLock() {
		_ = s.conn.SetWriteDeadline(time.Now().Add(terminateFrist))
		s.fe.Send(&pgproto3.Terminate{})
		_ = s.fe.Flush()
		s.schreiben.Unlock()
	}
	_ = s.conn.Close()
}
