package postgres_test

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// Codes der Nachricht R, wie PGWire sie festlegt.
const (
	rOk, rKlartext, rMD5, rGSSWeiter, rSASL, rSASLWeiter, rSASLAbschluss = 0, 3, 5, 8, 10, 11, 12
)

// rohR ist eine Nachricht R mit code und dem Rest daten.
func rohR(code uint32, daten string) []byte {
	return rohNachricht('R', string(binary.BigEndian.AppendUint32(nil, code))+daten)
}

// datenZu ist der Rest, den der Server bei code sendet.
func datenZu(code uint32) string {
	switch code {
	case rMD5:
		return "\x01\xff\x00\x7f"
	case rSASL:
		return "SCRAM-SHA-256\x00\x00"
	case rSASLWeiter:
		return "r=nonce,s=c2FseiE=,i=4096"
	case rSASLAbschluss:
		return "v=signatur"
	case rGSSWeiter:
		return "gss"
	}
	return ""
}

// erwartetAntwort sagt, ob der Server auf code eine Nachricht p erwartet.
func erwartetAntwort(code uint32) bool {
	return code == rKlartext || code == rMD5 || code == rSASL || code == rSASLWeiter
}

// skriptServer ist die Seite des Servers in einer Verbindung, die ein Skript
// treibt. gelesen sind die Rümpfe der Nachrichten p, die das Skript las, und
// fehler, was es dabei Unerwartetes sah; beides gilt nach fertig.
type skriptServer struct {
	conn    net.Conn
	br      *bufio.Reader
	fertig  chan struct{}
	gelesen [][]byte
	fehler  string
}

// lies liest eine Nachricht des Clients und liefert Typ und Rumpf.
func (s *skriptServer) lies() (byte, []byte, error) {
	kopf := make([]byte, 5)
	if _, err := io.ReadFull(s.br, kopf); err != nil {
		return 0, nil, err
	}
	rumpf := make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)
	_, err := io.ReadFull(s.br, rumpf)
	return kopf[0], rumpf, err
}

// rest liest, was der Recorder nach dem Skript noch sendet, bis er die Verbindung
// schließt, und liefert die Typen der Nachrichten; schließt er sie nicht, ist
// der Test rot.
func (s *skriptServer) rest(t *testing.T) string {
	t.Helper()
	<-s.fertig
	if s.fehler != "" {
		t.Error(s.fehler)
	}
	_ = s.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	var typen string
	for {
		typ, _, err := s.lies()
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return typen
		}
		if err != nil {
			t.Fatalf("der Recorder schließt die Verbindung zum Server nicht: %v", err)
		}
		typen += string(typ)
	}
}

// anmeldeServer nimmt eine Verbindung an, liest die Startnachricht und führt
// skript aus.
func anmeldeServer(t *testing.T, skript func(s *skriptServer)) (string, *skriptServer) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	srv := &skriptServer{fertig: make(chan struct{})}
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		t.Cleanup(func() { conn.Close() })
		srv.conn, srv.br = conn, bufio.NewReader(conn)
		if _, err := pgproto3.NewBackend(srv.br, conn).ReceiveStartupMessage(); err != nil {
			return
		}
		defer close(srv.fertig)
		skript(srv)
		// Der Server sendet nichts mehr: Ein Recorder, der weiterliest, sieht das Ende.
		_ = conn.(*net.TCPConn).CloseWrite()
	}()
	return l.Addr().String(), srv
}

// folgeSkript sendet die Codes der Folge der Reihe nach; auf einen, der eine
// Antwort erwartet, liest es die Nachricht p, außer nach dem letzten. Danach
// sendet es danach.
func folgeSkript(folge []uint32, danach ...[]byte) func(*skriptServer) {
	return func(s *skriptServer) {
		for i, code := range folge {
			_, _ = s.conn.Write(rohR(code, datenZu(code)))
			if !erwartetAntwort(code) || i == len(folge)-1 && len(danach) == 0 {
				continue
			}
			typ, rumpf, err := s.lies()
			if err != nil || typ != 'p' {
				s.fehler = fmt.Sprintf("auf Code %d keine Nachricht p: %c, %v", code, typ, err)
				return
			}
			s.gelesen = append(s.gelesen, rumpf)
		}
		for _, d := range danach {
			_, _ = s.conn.Write(d)
		}
	}
}

// aufbauOk ist der Rest des Aufbaus nach der Anmeldung.
func aufbauOk() [][]byte {
	var out [][]byte
	for _, m := range bereit() {
		b, _ := m.Encode(nil)
		out = append(out, b)
	}
	return out
}

// spion ist der Rückruf der Anmeldung: Er hält die Anforderungen und antwortet
// der Reihe nach mit antworten, danach mit „antwort“; auf Code 12 antwortet er
// nicht. Beim Aufruf i liefert er fehler[i], wenn gesetzt.
type spion struct {
	mu        sync.Mutex
	aufrufe   []model.Anmeldeanforderung
	antworten [][]byte
	fehler    map[int]error
}

func (r *spion) rufe(_ context.Context, a model.Anmeldeanforderung) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := len(r.aufrufe)
	r.aufrufe = append(r.aufrufe, a)
	if err := r.fehler[i]; err != nil {
		return nil, err
	}
	if a.Code == rSASLAbschluss {
		return nil, nil
	}
	if i < len(r.antworten) {
		return r.antworten[i], nil
	}
	return []byte("antwort"), nil
}

// codes liefert die Codes der Aufrufe.
func (r *spion) codes() []uint32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []uint32
	for _, a := range r.aufrufe {
		out = append(out, a.Code)
	}
	return out
}

// anmelden öffnet die Session zu addr mit dem Rückruf r.
func anmelden(t *testing.T, addr string, r *spion) (driven.UpstreamSession, []model.Response, error) {
	t.Helper()
	return (&postgres.Upstream{Address: addr}).Open(ctx(t), map[string]string{"user": "app"}, r.rufe)
}

// Abdeckung: LH-FA-05/Happy — der Recorder vermittelt Klartext-Passwort, MD5 und
// SCRAM-SHA-256: Er gibt jede Anforderung des Servers mit Code und Rest unverändert
// an den Rückruf und dessen Antwort unverändert als Nachricht p an den Server,
// auch eine leere und eine mit beliebigen Bytes; auf Code 12 erwartet er keine
// Antwort und sendet keine; AuthenticationOk geht nicht an den Rückruf und
// steht nicht in den Antworten des Aufbaus (LH-FA-05.b *Vermittlung*, *Verfahren*,
// *Antwort des Clients*).
func TestVermittlungVerfahren(t *testing.T) {
	for _, c := range []struct {
		name      string
		folge     []uint32
		antworten [][]byte
	}{
		{"Klartext-Passwort", []uint32{rKlartext}, [][]byte{[]byte("geheim\x00")}},
		{"MD5", []uint32{rMD5}, [][]byte{[]byte("md5\xff\x00\x01\x00")}},
		{"SCRAM-SHA-256", []uint32{rSASL, rSASLWeiter, rSASLAbschluss}, [][]byte{[]byte("SCRAM-SHA-256\x00\x00\x00\x00\x03n,,"), []byte("c=biws,r=\xfe\xff")}},
		{"leere Antwort", []uint32{rKlartext}, [][]byte{{}}},
		{"Antwort mit NUL-Bytes", []uint32{rMD5}, [][]byte{{0, 0, 0}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			addr, srv := anmeldeServer(t, folgeSkript(c.folge, append([][]byte{rohR(rOk, "")}, aufbauOk()[1:]...)...))
			r := &spion{antworten: c.antworten}
			s, out, err := anmelden(t, addr, r)
			if err != nil || s == nil {
				t.Fatalf("Open: %v", err)
			}
			_ = s.Close()
			rest := srv.rest(t)
			if got, want := fmt.Sprint(r.codes()), fmt.Sprint(c.folge); got != want {
				t.Fatalf("Aufrufe %s, erwartet %s", got, want)
			}
			for i, code := range c.folge {
				a := r.aufrufe[i]
				if a.Code != code || string(a.Daten) != datenZu(code) {
					t.Errorf("Anforderung %d: %+v, erwartet Code %d und %q", i, a, code, datenZu(code))
				}
			}
			if got, want := fmt.Sprintf("%q", srv.gelesen), fmt.Sprintf("%q", c.antworten); got != want {
				t.Errorf("der Server liest %s, erwartet %s", got, want)
			}
			if len(out) != 2 || out[0].Type != model.ResponseParameterStatus || out[1].Type != model.ResponseReadyForQuery {
				t.Errorf("Antworten des Aufbaus: %#v", out)
			}
			if rest != "X" {
				t.Errorf("nach dem Aufbau sendet der Recorder %q, erwartet nur Terminate", rest)
			}
		})
	}
}

// Abdeckung: LH-FA-05/Negative — jeder andere Code einer Anforderung (2, 6, 7, 9
// und ein unbekannter) und eine nicht lesbare Nachricht des Servers im Aufbau
// sind PGR-E6001; der Rückruf läuft nicht, der Recorder sendet dem Server
// weder eine weitere Nachricht noch Terminate, und die Meldung nennt den Code,
// nicht den Typ der Bibliothek (LH-FA-05.b *Nicht vermittelt*).
func TestVermittlungNichtVermittelt(t *testing.T) {
	for _, c := range []struct {
		name, nachricht, code string
	}{
		{"Code 2", string(rohR(2, "GEHEIM")), "2"},
		{"Code 6", string(rohR(6, "GEHEIM")), "6"},
		{"Code 7", string(rohR(7, "GEHEIM")), "7"},
		{"Code 9", string(rohR(9, "GEHEIM")), "9"},
		{"unbekannter Code", string(rohR(4242, "GEHEIM")), "4242"},
		{"nicht lesbare Nachricht", string(rohNachricht('z', "x")), ""},
		{"Nachricht R ohne Code", string(rohNachricht('R', "")), ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			addr, srv := anmeldeServer(t, func(s *skriptServer) { _, _ = s.conn.Write([]byte(c.nachricht)) })
			r := &spion{}
			s, _, err := anmelden(t, addr, r)
			if s != nil || code(err) != model.CodeUnsupported {
				t.Fatalf("erwartet %s, erhalten %v", model.CodeUnsupported, err)
			}
			if c.code != "" && !strings.Contains(err.Error(), "Code "+c.code+")") {
				t.Errorf("Meldung ohne den Code %q: %v", c.code, err)
			}
			for _, bibliothek := range []string{"pgproto3", "unknown authentication", "unimplemented", "AuthType", "dieser Stand", "GEHEIM"} {
				if strings.Contains(err.Error(), bibliothek) {
					t.Errorf("Meldung nennt die Bibliothek (%q): %v", bibliothek, err)
				}
			}
			if len(r.aufrufe) != 0 {
				t.Errorf("Rückruf gerufen: %v", r.aufrufe)
			}
			if rest := srv.rest(t); rest != "" {
				t.Errorf("der Recorder sendet nach dem Abbruch %q", rest)
			}
		})
	}
}

// Abdeckung: LH-FA-05/Negative — der Recorder lehnt jede Nachricht R des Servers
// außerhalb der Reihenfolge mit PGR-E4002 ab, die den Code nennt: eine
// Fortsetzung (8, 11, 12) ohne laufenden Austausch ihres Verfahrens, Code 12
// vor Code 11, AuthenticationOk im laufenden SASL-Austausch, eine weitere
// Anforderung (3, 5, 10) nach einer beantworteten und im SASL-Austausch, jede
// Nachricht R nach AuthenticationOk und nach SASLFinal; der Server bekommt
// danach nichts mehr (LH-FA-05.b *Reihenfolge des Servers*).
func TestVermittlungReihenfolge(t *testing.T) {
	for _, folge := range [][]uint32{
		{rGSSWeiter}, {rSASLWeiter}, {rSASLAbschluss},
		{rKlartext, rSASLWeiter}, {rKlartext, rSASLAbschluss}, {rKlartext, rGSSWeiter}, {rMD5, rSASLWeiter},
		{rSASL, rSASLAbschluss}, {rSASL, rOk}, {rSASL, rSASLWeiter, rOk},
		{rKlartext, rKlartext}, {rKlartext, rMD5}, {rKlartext, rSASL}, {rMD5, rKlartext},
		{rSASL, rKlartext}, {rSASL, rMD5}, {rSASL, rSASL}, {rSASL, rSASLWeiter, rKlartext},
		{rSASL, rSASLWeiter, rSASLAbschluss, rSASLAbschluss}, {rSASL, rSASLWeiter, rSASLAbschluss, rSASLWeiter}, {rSASL, rSASLWeiter, rSASLAbschluss, rKlartext},
		{rOk, rKlartext}, {rOk, rOk}, {rOk, rSASLWeiter},
	} {
		t.Run(fmt.Sprint(folge), func(t *testing.T) {
			addr, srv := anmeldeServer(t, folgeSkript(folge))
			r := &spion{}
			s, _, err := anmelden(t, addr, r)
			falsch := folge[len(folge)-1]
			if s != nil || code(err) != model.CodeUpstream || !strings.Contains(err.Error(), fmt.Sprintf("Code %d ", falsch)) {
				t.Fatalf("erwartet %s mit Code %d, erhalten %v", model.CodeUpstream, falsch, err)
			}
			if d := datenZu(falsch); d != "" && (strings.Contains(err.Error(), d) || strings.Contains(err.Error(), hex.EncodeToString([]byte(d)))) {
				t.Errorf("die Meldung nennt den Inhalt der Nachricht: %v", err)
			}
			if rest := srv.rest(t); rest != "" {
				t.Errorf("der Recorder sendet nach dem Abbruch %q", rest)
			}
		})
	}
}

// Abdeckung: LH-FA-05/Boundary — liefert der Rückruf einen Fehler (Ende des
// Clients, PGR-E6001 aus dem Client), endet Open mit genau diesem Fehler, ohne
// Session, bei jeder Anforderung, auch bei Code 12; der Recorder schließt die
// Verbindung zum Server und sendet weder eine weitere Nachricht noch
// Terminate (LH-FA-05.b *Ende durch den Client*, *Antwort des Clients*).
func TestVermittlungFehlerDesRueckrufs(t *testing.T) {
	unterstuetzt := model.Errorf(model.CodeUnsupported, nil, "Client-Nachricht in der Anmeldung")
	for _, c := range []struct {
		name  string
		folge []uint32
		i     int
		err   error
	}{
		{"Ende bei Klartext", []uint32{rKlartext}, 0, model.ErrClientEnde},
		{"Ende bei SASL", []uint32{rSASL}, 0, model.ErrClientEnde},
		{"Ende mitten im SCRAM", []uint32{rSASL, rSASLWeiter}, 1, model.ErrClientEnde},
		{"Ende bei Code 12", []uint32{rSASL, rSASLWeiter, rSASLAbschluss}, 2, model.ErrClientEnde},
		{"PGR-E6001 des Clients", []uint32{rMD5}, 0, unterstuetzt},
	} {
		t.Run(c.name, func(t *testing.T) {
			addr, srv := anmeldeServer(t, folgeSkript(c.folge))
			r := &spion{fehler: map[int]error{c.i: c.err}}
			s, out, err := anmelden(t, addr, r)
			if s != nil || out != nil || !errors.Is(err, c.err) {
				t.Fatalf("Session %v, Antworten %v, Fehler %v; erwartet %v", s, out, err, c.err)
			}
			if rest := srv.rest(t); rest != "" {
				t.Errorf("der Recorder sendet nach dem Abbruch %q", rest)
			}
		})
	}
}

// Abdeckung: LH-FA-05/Negative — eine Fehlerantwort des Servers, am Anfang
// (reject) und nach jeder Antwort (Klartext, SASL, mitten im SCRAM), geht mit
// allen Feldern unverändert zurück, ohne Session und ohne Fehler; der Recorder
// schließt die Verbindung zum Server ohne weitere Nachricht
// (LH-FA-05.b *Fehlerantwort des Servers*).
func TestVermittlungFehlerantwortDesServers(t *testing.T) {
	fehler, _ := (&pgproto3.ErrorResponse{
		Severity: "FATAL", SeverityUnlocalized: "FATAL", Code: "28P01", Message: "m", Detail: "d", Hint: "h", Position: 1, InternalPosition: 2,
		InternalQuery: "q", Where: "w", SchemaName: "s", TableName: "t", ColumnName: "c", DataTypeName: "dt", ConstraintName: "n",
		File: "f", Line: 3, Routine: "r", UnknownFields: map[byte]string{'Z': "z"},
	}).Encode(nil)
	want := map[string]string{
		"S": "FATAL", "V": "FATAL", "C": "28P01", "M": "m", "D": "d", "H": "h", "P": "1", "p": "2", "q": "q", "W": "w",
		"s": "s", "t": "t", "c": "c", "d": "dt", "n": "n", "F": "f", "L": "3", "R": "r", "Z": "z",
	}
	for _, c := range []struct {
		name  string
		folge []uint32
	}{
		{"am Anfang", nil},
		{"nach Klartext", []uint32{rKlartext}},
		{"nach Code 10", []uint32{rSASL}},
		{"mitten im SCRAM", []uint32{rSASL, rSASLWeiter}},
	} {
		t.Run(c.name, func(t *testing.T) {
			addr, srv := anmeldeServer(t, folgeSkript(c.folge, fehler))
			r := &spion{}
			s, out, err := anmelden(t, addr, r)
			if s != nil || err != nil || len(out) != 1 || out[0].Type != model.ResponseErrorResponse {
				t.Fatalf("Session %v, Antworten %#v, Fehler %v", s, out, err)
			}
			if fmt.Sprint(out[0].Fields) != fmt.Sprint(want) {
				t.Errorf("Felder %v, erwartet %v", out[0].Fields, want)
			}
			if rest := srv.rest(t); rest != "" {
				t.Errorf("der Recorder sendet nach dem Abbruch %q", rest)
			}
		})
	}
}

// Abdeckung: LH-FA-02/Boundary — die Session aus Open trägt den Lesepuffer des
// Aufbaus: Schickt der Server hinter dem ReadyForQuery in demselben Schreibvorgang
// eine NoticeResponse, liefert der erste Receive nach Open sie (auf TCP-Loopback
// kommen beide Nachrichten in einem Lesevorgang an; zugesagt ist das nicht).
func TestOpenSessionTraegtLesepufferDesAufbaus(t *testing.T) {
	addr, srv := anmeldeServer(t, func(s *skriptServer) {
		var alles []byte
		for _, d := range aufbauOk() {
			alles = append(alles, d...)
		}
		notice, _ := (&pgproto3.NoticeResponse{Severity: "NOTICE", Code: "00000", Message: "hinweis"}).Encode(nil)
		_, _ = s.conn.Write(append(alles, notice...))
	})
	s, _, err := anmelden(t, addr, &spion{})
	if err != nil {
		t.Fatal(err)
	}
	if err := postgres.Verbindung(s).SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	out, err := s.Receive(ctx(t))
	if err != nil || len(out) != 1 || out[0].Type != model.ResponseNoticeResponse || out[0].Fields["M"] != "hinweis" {
		t.Errorf("Receive nach Open: %#v, %v", out, err)
	}
	_ = s.Close()
	srv.rest(t)
}

// Abdeckung: LH-FA-05/Negative — eine Fehlerantwort des Servers in der Anmeldung
// geht gleich welcher SQLSTATE zurück, ohne Session und ohne Fehler: Anmeldung
// (28P01, 28000), Datenbank (3D000), Überlast (53300, 57P03) und ein Fehler
// außerhalb dieser Klassen (XX000) (LH-FA-05.b *Fehlerantwort des Servers*).
func TestVermittlungFehlerantwortJederSQLSTATE(t *testing.T) {
	for _, sqlstate := range []string{"28P01", "28000", "3D000", "53300", "57P03", "XX000"} {
		t.Run(sqlstate, func(t *testing.T) {
			fehler, _ := (&pgproto3.ErrorResponse{Severity: "FATAL", Code: sqlstate, Message: "m"}).Encode(nil)
			addr, srv := anmeldeServer(t, folgeSkript([]uint32{rSASL, rSASLWeiter}, fehler))
			s, out, err := anmelden(t, addr, &spion{})
			if s != nil || err != nil || len(out) != 1 || out[0].Fields["C"] != sqlstate {
				t.Fatalf("Session %v, Antworten %#v, Fehler %v", s, out, err)
			}
			if rest := srv.rest(t); rest != "" {
				t.Errorf("der Recorder sendet nach dem Abbruch %q", rest)
			}
		})
	}
}
