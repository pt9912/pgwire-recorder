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
// Server unter Address, ohne Passwort und ohne TLS (LH-FA-20.a).
type Einspielziel struct {
	Address string
	Dialer  net.Dialer
}

var _ driven.Einspielziel = (*Einspielziel)(nil)

// Codes der Anmelde-Nachricht R (LH-FA-20.a *Anmelde-Nachrichten*):
// AuthenticationOk und die Fortsetzungen GSSContinue, SASLContinue und
// SASLFinal; jeder andere Code ist die Anforderung eines Verfahrens.
const (
	anmeldungOk            = 0
	anmeldungGSSWeiter     = 8
	anmeldungSASLWeiter    = 11
	anmeldungSASLAbschluss = 12
)

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
	s, err := aufbau(conn, startup)
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

// aufbau sendet das Startup und liest bis zum ersten ReadyForQuery. Es liest
// jede Nachricht selbst, damit der Code einer Anmelde-Nachricht auch dort
// gilt, wo pgproto3 ihn nicht dekodiert. Der Puffer des Lesens geht mit der
// Session weiter, sodass nichts verloren geht, was der Server hinter dem
// ReadyForQuery sendet.
func aufbau(conn net.Conn, startup map[string]string) (*einspielSession, error) {
	start, err := (&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: startup}).Encode(nil)
	if err != nil {
		return nil, model.Errorf(model.CodeUpstream, err, "Startup nicht zu kodieren")
	}
	if _, err := conn.Write(start); err != nil {
		return nil, model.Errorf(model.CodeUpstream, err, "Startup nicht zu senden")
	}
	r := bufio.NewReader(conn)
	angemeldet := false
	for {
		typ, rumpf, err := liesRoh(r)
		if err != nil {
			return nil, model.Errorf(model.CodeUpstream, err, "Verbindung im Aufbau beendet oder Nachricht nicht lesbar")
		}
		switch typ {
		case 'R':
			if err := anmeldeNachricht(rumpf, angemeldet); err != nil {
				return nil, err
			}
			angemeldet = true
		case 'E':
			return nil, fehlerImAufbau(rumpf)
		case 'Z':
			if !angemeldet {
				return nil, model.Errorf(model.CodeUpstream, nil, "ReadyForQuery vor AuthenticationOk")
			}
			if err := (&pgproto3.ReadyForQuery{}).Decode(rumpf); err != nil {
				return nil, model.Errorf(model.CodeUpstream, err, "ReadyForQuery im Aufbau nicht lesbar")
			}
			return &einspielSession{conn: conn, fe: pgproto3.NewFrontend(r, conn)}, nil
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

// anmeldeNachricht stuft eine Nachricht R nach ihrem Code ein (LH-FA-20.a
// *Anmelde-Nachrichten*): ohne vollständigen Code PGR-E4002; nach
// AuthenticationOk jede PGR-E4002; AuthenticationOk ist ohne Fehler, wenn es
// sich lesen lässt; eine Fortsetzung ohne laufenden Austausch PGR-E4002; jeder
// andere Code fordert ein Verfahren an, das dieser Stand nicht unterstützt,
// und ist PGR-E4005.
func anmeldeNachricht(rumpf []byte, angemeldet bool) error {
	if len(rumpf) < 4 {
		return model.Errorf(model.CodeUpstream, nil, "Anmelde-Nachricht ohne vollständigen Code")
	}
	code := binary.BigEndian.Uint32(rumpf)
	switch {
	case angemeldet:
		return model.Errorf(model.CodeUpstream, nil, "Anmelde-Nachricht (Code %d) nach AuthenticationOk", code)
	case code == anmeldungOk:
		if err := (&pgproto3.AuthenticationOk{}).Decode(rumpf); err != nil {
			return model.Errorf(model.CodeUpstream, err, "AuthenticationOk nicht lesbar")
		}
		return nil
	case code == anmeldungGSSWeiter || code == anmeldungSASLWeiter || code == anmeldungSASLAbschluss:
		return model.Errorf(model.CodeUpstream, nil, "Fortsetzung einer Anmeldung (Code %d) ohne laufenden Austausch", code)
	}
	return model.Errorf(model.CodeLogin, nil, "der Server verlangt ein Anmeldeverfahren (Code %d), das play nicht unterstützt", code)
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
// hält das Senden von Anfrage und das Terminate von Schliesse auseinander.
type einspielSession struct {
	conn      net.Conn
	fe        *pgproto3.Frontend
	schreiben sync.Mutex
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

func (s *einspielSession) Naechste() (model.Response, error) {
	for {
		msg, err := s.fe.Receive()
		if err != nil {
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

// Schliesse sendet Terminate nur, wenn gerade keine Anfrage sendet, und
// höchstens terminateFrist lang; danach schließt es die Verbindung.
func (s *einspielSession) Schliesse() {
	if s.schreiben.TryLock() {
		_ = s.conn.SetWriteDeadline(time.Now().Add(terminateFrist))
		s.fe.Send(&pgproto3.Terminate{})
		_ = s.fe.Flush()
		s.schreiben.Unlock()
	}
	_ = s.conn.Close()
}
