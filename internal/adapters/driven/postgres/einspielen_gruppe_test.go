package postgres_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"runtime/pprof"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// verbindeGruppen baut gegen den Fake-Server eine Session auf, die Gruppe
// kennt.
func verbindeGruppen(t *testing.T, addr string) driven.EinspielSession {
	t.Helper()
	return verbinde(t, addr).(driven.EinspielSession)
}

// gruppenServer nimmt eine Verbindung an, liest das Startup, sendet aufbau und
// sofort antwort; Client-Nachrichten liest er erst, wenn lesen geschlossen ist
// (nil: sofort), und zwar erwartet Bytes, bevor er ende sendet; danach liest er
// bis zum Verbindungsende oder seiner Frist von 30 s. lauf.danach sind alle
// Bytes nach dem Startup.
func gruppenServer(t *testing.T, antwort []byte, lesen <-chan struct{}, erwartet int, ende []byte) (string, <-chan lauf) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	ergebnis := make(chan lauf, 1)
	aufbau := verbunden(t)
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
		var laenge [4]byte
		if _, err := io.ReadFull(conn, laenge[:]); err != nil {
			return
		}
		if _, err := io.ReadFull(conn, make([]byte, binary.BigEndian.Uint32(laenge[:])-4)); err != nil {
			return
		}
		_, _ = conn.Write(append(append([]byte{}, aufbau...), antwort...))
		if lesen != nil {
			<-lesen
		}
		erstes := make([]byte, erwartet)
		n, err := io.ReadFull(conn, erstes)
		if err != nil {
			ergebnis <- lauf{danach: erstes[:n], ende: err}
			return
		}
		_, _ = conn.Write(ende)
		rest, err := io.ReadAll(conn)
		ergebnis <- lauf{danach: append(erstes, rest...), ende: err}
	}()
	return l.Addr().String(), ergebnis
}

// frontend ist die Folge der kodierten Client-Nachrichten.
func frontend(t *testing.T, msgs ...pgproto3.FrontendMessage) []byte {
	t.Helper()
	var out []byte
	for _, m := range msgs {
		out = append(out, kodiertFrontend(t, m)...)
	}
	return out
}

// naechsteBinnen ruft Naechste und wartet höchstens 5 s auf die Rückkehr.
func naechsteBinnen(t *testing.T, s driven.EinspielSession, was string) (model.Response, error) {
	t.Helper()
	type ergebnis struct {
		r   model.Response
		err error
	}
	fertig := make(chan ergebnis, 1)
	go func() {
		r, err := s.Naechste()
		fertig <- ergebnis{r, err}
	}()
	select {
	case e := <-fertig:
		return e.r, e.err
	case <-time.After(5 * time.Second):
		t.Fatalf("Naechste kehrt binnen 5 s nicht zurück: %s", was)
		return model.Response{}, nil
	}
}

// gruppeBinnen ruft Gruppe und wartet höchstens 5 s auf die Rückkehr.
func gruppeBinnen(t *testing.T, s driven.EinspielSession, was string, msgs ...model.ClientMessage) error {
	t.Helper()
	fertig := make(chan error, 1)
	go func() { fertig <- s.Gruppe(msgs) }()
	select {
	case err := <-fertig:
		return err
	case <-time.After(5 * time.Second):
		t.Fatalf("Gruppe kehrt binnen 5 s nicht zurück: %s", was)
		return nil
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-18/Boundary — Gruppe sendet die
// Client-Nachrichten wie aufgezeichnet auf die Verbindung: Parse, Bind,
// Describe, Execute, Close, Flush und Sync mit ihren Feldern, die Gruppen in der
// Reihenfolge der Aufrufe; ein Parameter mit Null ist SQL-NULL (Länge −1), ein
// leerer Wert ist ein leerer Wert, gleich ob die geladenen Bytes nil oder leer
// sind; param_formats und result_formats unverändert, auch leer (LH-FA-20.a
// *Gruppen*, SPEC-041).
func TestEinspielGruppeNachrichten(t *testing.T) {
	g1 := []model.ClientMessage{
		{Type: model.ClientParse, Statement: "s1", SQL: "SELECT $1, $2, $3, $4", ParamTypes: []uint32{25, 25, 25, 25}},
		{Type: model.ClientBind, Portal: "p1", Statement: "s1", ParamFormats: []int16{0, 1, 0, 0}, Params: []model.Value{{Null: true}, {}, {Bytes: []byte{}}, {Bytes: []byte("x")}}, ResultFormats: []int16{1}},
		{Type: model.ClientBind, Portal: "", Statement: "", Params: []model.Value{}},
		{Type: model.ClientDescribe, Target: model.TargetStatement, Name: "s1"},
		{Type: model.ClientFlush},
	}
	g2 := []model.ClientMessage{
		{Type: model.ClientDescribe, Target: model.TargetPortal, Name: "p1"},
		{Type: model.ClientExecute, Portal: "p1", MaxRows: 5},
		{Type: model.ClientClose, Target: model.TargetPortal, Name: "p1"},
		{Type: model.ClientClose, Target: model.TargetStatement, Name: "s1"},
		{Type: model.ClientSync},
	}
	erwartet := frontend(t,
		&pgproto3.Parse{Name: "s1", Query: "SELECT $1, $2, $3, $4", ParameterOIDs: []uint32{25, 25, 25, 25}},
		&pgproto3.Bind{DestinationPortal: "p1", PreparedStatement: "s1", ParameterFormatCodes: []int16{0, 1, 0, 0}, Parameters: [][]byte{nil, {}, {}, []byte("x")}, ResultFormatCodes: []int16{1}},
		&pgproto3.Bind{Parameters: [][]byte{}},
		&pgproto3.Describe{ObjectType: 'S', Name: "s1"},
		&pgproto3.Flush{},
		&pgproto3.Describe{ObjectType: 'P', Name: "p1"},
		&pgproto3.Execute{Portal: "p1", MaxRows: 5},
		&pgproto3.Close{ObjectType: 'P', Name: "p1"},
		&pgproto3.Close{ObjectType: 'S', Name: "s1"},
		&pgproto3.Sync{},
	)
	addr, ergebnis := gruppenServer(t, nil, nil, len(erwartet), kodiert(t, &pgproto3.ReadyForQuery{TxStatus: 'I'}))
	s := verbindeGruppen(t, addr)
	if err := gruppeBinnen(t, s, "erste Gruppe", g1...); err != nil {
		t.Fatal(err)
	}
	if err := gruppeBinnen(t, s, "zweite Gruppe", g2...); err != nil {
		t.Fatal(err)
	}
	// Das ReadyForQuery sendet der Server erst, wenn alle Bytes der Gruppen bei
	// ihm sind; danach schließt Schliesse mit Terminate.
	gelesen := make(chan model.Response, 1)
	go func() {
		r, _ := s.Naechste()
		gelesen <- r
	}()
	select {
	case r := <-gelesen:
		if r.Type != model.ResponseReadyForQuery {
			t.Fatalf("Naechste: %v", r)
		}
	case <-time.After(5 * time.Second):
		s.Schliesse()
		l := <-ergebnis
		t.Fatalf("der Server empfängt binnen 5 s nicht alle %d Byte, sondern\n%q", len(erwartet), l.danach)
	}
	s.Schliesse()
	if l := empfangen(t, ergebnis); !bytes.Equal(l.danach, append(erwartet, kodiertFrontend(t, &pgproto3.Terminate{})...)) {
		t.Fatalf("gesendet\n%q\nerwartet\n%q", l.danach, append(erwartet, kodiertFrontend(t, &pgproto3.Terminate{})...))
	}
}

// Abdeckung: LH-FA-20/Negative — eine Client-Nachricht, die sich nicht auf
// PGWire abbilden lässt (Zielart außer statement und portal), liefert Gruppe
// als Fehler PGR-E1000, und von der Gruppe geht nichts auf die Verbindung
// (LH-FA-20.a *Interaktion*).
func TestEinspielGruppeNichtAbbildbar(t *testing.T) {
	addr, ergebnis := einspielServer(t, verbunden(t))
	s := verbindeGruppen(t, addr)
	err := gruppeBinnen(t, s, "nicht abbildbare Gruppe",
		model.ClientMessage{Type: model.ClientParse, SQL: "SELECT 1"},
		model.ClientMessage{Type: model.ClientDescribe, Target: "", Name: "x"})
	if code(err) != model.CodeInternal {
		t.Fatalf("Gruppe: %v, erwartet %s", err, model.CodeInternal)
	}
	s.Schliesse()
	if l := empfangen(t, ergebnis); !bytes.Equal(l.danach, kodiertFrontend(t, &pgproto3.Terminate{})) {
		t.Fatalf("gesendet %q, erwartet nur Terminate", l.danach)
	}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-18/Boundary — Gruppe kehrt zurück, ohne
// auf das Senden zu warten, und das Lesen ist unabhängig vom Senden: Solange
// das Senden einer Gruppe steht, weil der Server nichts liest, liefert Naechste
// seine Antworten, und weitere Gruppen werden eingereiht und nach ihr gesendet,
// in der Reihenfolge der Aufrufe (LH-FA-20.a *Gruppen*).
func TestEinspielGruppeUnabhaengig(t *testing.T) {
	gross := model.ClientMessage{Type: model.ClientParse, SQL: strings.Repeat("x", 64<<20)}
	zweite := model.ClientMessage{Type: model.ClientDescribe, Target: model.TargetPortal, Name: "zweite"}
	dritte := model.ClientMessage{Type: model.ClientExecute, Portal: "dritte"}
	erwartet := frontend(t,
		&pgproto3.Parse{Query: strings.Repeat("x", 64<<20)},
		&pgproto3.Describe{ObjectType: 'P', Name: "zweite"},
		&pgproto3.Execute{Portal: "dritte"},
	)
	lesen := make(chan struct{})
	antwort := kodiert(t, &pgproto3.ParseComplete{}, &pgproto3.ReadyForQuery{TxStatus: 'I'})
	addr, ergebnis := gruppenServer(t, antwort, lesen, len(erwartet), kodiert(t, &pgproto3.CloseComplete{}))
	s := verbindeGruppen(t, addr)
	if err := gruppeBinnen(t, s, "große Gruppe, der Server liest nicht", gross); err != nil {
		t.Fatal(err)
	}
	if r, err := naechsteBinnen(t, s, "während das Senden steht"); err != nil || r.Type != model.ResponseParseComplete {
		t.Fatalf("Naechste während des Sendens: %v, %v", r, err)
	}
	if err := gruppeBinnen(t, s, "zweite Gruppe hinter der stehenden", zweite); err != nil {
		t.Fatal(err)
	}
	if err := gruppeBinnen(t, s, "dritte Gruppe hinter der stehenden", dritte); err != nil {
		t.Fatal(err)
	}
	close(lesen)
	if r, err := naechsteBinnen(t, s, "ReadyForQuery"); err != nil || r.Type != model.ResponseReadyForQuery {
		t.Fatalf("Naechste: %v, %v", r, err)
	}
	if r, err := naechsteBinnen(t, s, "Antwort nach allen Gruppen"); err != nil || r.Type != model.ResponseCloseComplete {
		t.Fatalf("Naechste nach den Gruppen: %v, %v", r, err)
	}
	s.Schliesse()
	l := empfangen(t, ergebnis)
	if got := l.danach; !bytes.HasPrefix(got, erwartet) || len(got) != len(erwartet)+len(kodiertFrontend(t, &pgproto3.Terminate{})) {
		t.Fatalf("gesendet %d Byte (Anfang gleich: %v), erwartet %d Byte der drei Gruppen in Reihenfolge und Terminate", len(got), bytes.HasPrefix(got, erwartet), len(erwartet))
	}
}

// steuerConn ist eine Verbindung, deren Read erst daten liefert, dann bis Close
// blockiert und dann net.ErrClosed liefert; Write ruft write. schreibt und
// schliesst zählen die Aufrufe von Write und Close.
type steuerConn struct {
	net.Conn
	daten     []byte
	zu        chan struct{}
	zuEinmal  sync.Once
	write     func(b []byte) (int, error)
	mu        sync.Mutex
	schreibt  int
	schliesst int
}

func neueSteuerConn(write func(b []byte) (int, error)) *steuerConn {
	return &steuerConn{zu: make(chan struct{}), write: write}
}

func (c *steuerConn) Read(b []byte) (int, error) {
	c.mu.Lock()
	if len(c.daten) > 0 {
		n := copy(b, c.daten)
		c.daten = c.daten[n:]
		c.mu.Unlock()
		return n, nil
	}
	c.mu.Unlock()
	<-c.zu
	return 0, net.ErrClosed
}

func (c *steuerConn) Write(b []byte) (int, error) {
	c.mu.Lock()
	c.schreibt++
	c.mu.Unlock()
	return c.write(b)
}

func (c *steuerConn) Close() error {
	c.mu.Lock()
	c.schliesst++
	c.mu.Unlock()
	c.zuEinmal.Do(func() { close(c.zu) })
	return nil
}

func (c *steuerConn) SetWriteDeadline(time.Time) error { return nil }

func (c *steuerConn) zaehler() (schreibt, schliesst int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.schreibt, c.schliesst
}

// Abdeckung: LH-FA-20/Negative — scheitert das Senden einer Gruppe, kehrt Gruppe
// zuvor zurück, und der Fehler PGR-E4003 kommt aus Naechste, auch wenn diese
// gerade wartet, und aus der nächsten Operation; die Verbindung ist danach
// geschlossen (LH-FA-20.a *Interaktion*, *Gruppen*).
func TestEinspielGruppeSendefehler(t *testing.T) {
	conn := neueSteuerConn(func([]byte) (int, error) { return 0, errors.New("kaputt") })
	s := postgres.NeueEinspielSession(conn)
	wartet := make(chan error, 1)
	go func() {
		_, err := s.Naechste()
		wartet <- err
	}()
	if err := gruppeBinnen(t, s, "Senden scheitert", model.ClientMessage{Type: model.ClientSync}); err != nil {
		t.Fatalf("Gruppe: %v, erwartet die Rückkehr ohne den Fehler des Sendens", err)
	}
	pruefe := func(was string, err error) {
		t.Helper()
		if code(err) != model.CodeConnectionLost || err == nil || !strings.Contains(err.Error(), "Gruppe an den Server nicht zu senden") {
			t.Fatalf("%s: %v, erwartet %s mit dem Fehler des Sendens", was, err, model.CodeConnectionLost)
		}
	}
	select {
	case err := <-wartet:
		pruefe("wartendes Naechste", err)
	case <-time.After(5 * time.Second):
		t.Fatal("das wartende Naechste endet binnen 5 s nach dem gescheiterten Senden nicht")
	}
	_, err := naechsteBinnen(t, s, "nach dem gescheiterten Senden")
	pruefe("nächstes Naechste", err)
	pruefe("nächste Gruppe", gruppeBinnen(t, s, "nach dem gescheiterten Senden", model.ClientMessage{Type: model.ClientSync}))
	s.Schliesse()

	// Liegt hinter dem gescheiterten Senden noch eine Antwort bereit, liest
	// Naechste sie nicht mehr.
	conn = neueSteuerConn(func([]byte) (int, error) { return 0, errors.New("kaputt") })
	conn.daten = kodiert(t, &pgproto3.ParseComplete{})
	s = postgres.NeueEinspielSession(conn)
	if err := gruppeBinnen(t, s, "Senden scheitert", model.ClientMessage{Type: model.ClientSync}); err != nil {
		t.Fatal(err)
	}
	bis := time.Now().Add(5 * time.Second)
	for {
		if _, schliesst := conn.zaehler(); schliesst > 0 {
			break
		}
		if time.Now().After(bis) {
			t.Fatal("das gescheiterte Senden schließt die Verbindung binnen 5 s nicht")
		}
		time.Sleep(time.Millisecond)
	}
	r, err := naechsteBinnen(t, s, "Antwort hinter dem gescheiterten Senden")
	if err == nil {
		t.Fatalf("Naechste liefert %v hinter dem gescheiterten Senden", r)
	}
	pruefe("Naechste mit bereitliegender Antwort", err)
	s.Schliesse()
}

// Abdeckung: LH-FA-20/Boundary — Schliesse beendet den Sender: Es endet, während
// eine Gruppe sendet, ohne auf das Senden zu warten und ohne Terminate zwischen
// die Bytes zu schreiben; das gescheiterte Senden danach bleibt ohne Folge: Es
// ist kein Fehler des Sendens für Naechste, und weitere Gruppen nimmt die
// Session nicht mehr an (LH-FA-20.a *Gruppen*, *Abbruchsignal*).
func TestEinspielGruppeSchliesseBeimSenden(t *testing.T) {
	angefangen := make(chan struct{})
	var conn *steuerConn
	conn = neueSteuerConn(func([]byte) (int, error) {
		close(angefangen)
		<-conn.zu
		return 0, net.ErrClosed
	})
	s := postgres.NeueEinspielSession(conn)
	if err := gruppeBinnen(t, s, "Senden steht", model.ClientMessage{Type: model.ClientSync}); err != nil {
		t.Fatal(err)
	}
	select {
	case <-angefangen:
	case <-time.After(5 * time.Second):
		t.Fatal("die Gruppe beginnt binnen 5 s nicht zu senden")
	}
	geschlossen := make(chan struct{})
	go func() {
		s.Schliesse()
		close(geschlossen)
	}()
	select {
	case <-geschlossen:
	case <-time.After(5 * time.Second):
		t.Fatal("Schliesse endet binnen 5 s nicht, während eine Gruppe sendet")
	}
	// Close ruft Schliesse und, nach dem gescheiterten Senden, der Sender; erst
	// danach ist entschieden, ob das Senden als Fehler gemerkt wurde.
	bis := time.Now().Add(5 * time.Second)
	for {
		if _, schliesst := conn.zaehler(); schliesst >= 2 {
			break
		}
		if time.Now().After(bis) {
			t.Fatal("der Sender endet binnen 5 s nach Schliesse nicht")
		}
		time.Sleep(time.Millisecond)
	}
	if schreibt, _ := conn.zaehler(); schreibt != 1 {
		t.Fatalf("%d Schreibvorgänge, erwartet 1: Schliesse schreibt nichts, während die Gruppe sendet", schreibt)
	}
	_, err := naechsteBinnen(t, s, "nach Schliesse")
	if code(err) != model.CodeConnectionLost || err == nil || strings.Contains(err.Error(), "nicht zu senden") {
		t.Fatalf("Naechste nach Schliesse: %v, erwartet das Verbindungsende ohne Fehler des Sendens", err)
	}
	if err := gruppeBinnen(t, s, "nach Schliesse", model.ClientMessage{Type: model.ClientSync}); code(err) != model.CodeConnectionLost {
		t.Fatalf("Gruppe nach Schliesse: %v, erwartet %s", err, model.CodeConnectionLost)
	}
}

// senderZahl ist die Zahl der Goroutinen im Sender einer Session.
func senderZahl() int {
	var b strings.Builder
	_ = pprof.Lookup("goroutine").WriteTo(&b, 2)
	return strings.Count(b.String(), "(*einspielSession).sender(")
}

// Abdeckung: LH-FA-20/Boundary — Schliesse beendet den Sender auch, wenn er
// nichts sendet: Nach der letzten Gruppe läuft ein Sender, nach Schliesse keiner
// mehr (LH-FA-20.a *Gruppen*).
func TestEinspielGruppeSchliesseBeendetSender(t *testing.T) {
	// Sender früherer Tests enden nach ihrem Schliesse, aber nicht synchron.
	bis := time.Now().Add(5 * time.Second)
	for senderZahl() != 0 {
		if time.Now().After(bis) {
			t.Fatalf("%d Sender früherer Tests laufen binnen 5 s weiter", senderZahl())
		}
		time.Sleep(time.Millisecond)
	}
	const vorher = 0
	addr, ergebnis := einspielServer(t, verbunden(t))
	s := verbindeGruppen(t, addr)
	if err := gruppeBinnen(t, s, "kleine Gruppe", model.ClientMessage{Type: model.ClientSync}); err != nil {
		t.Fatal(err)
	}
	bis = time.Now().Add(5 * time.Second)
	for senderZahl() != vorher+1 {
		if time.Now().After(bis) {
			t.Fatalf("%d Sender, erwartet %d: Gruppe startet keinen", senderZahl(), vorher+1)
		}
		time.Sleep(time.Millisecond)
	}
	s.Schliesse()
	bis = time.Now().Add(5 * time.Second)
	for senderZahl() != vorher {
		if time.Now().After(bis) {
			t.Fatalf("%d Sender nach Schliesse, erwartet %d", senderZahl(), vorher)
		}
		time.Sleep(time.Millisecond)
	}
	empfangen(t, ergebnis)
}

// warteAufKeinenSender wartet höchstens 5 s darauf, dass kein Sender mehr
// läuft; sonst endet der Test mit meldung.
func warteAufKeinenSender(t *testing.T, meldung string) {
	t.Helper()
	bis := time.Now().Add(5 * time.Second)
	for senderZahl() != 0 {
		if time.Now().After(bis) {
			t.Fatalf("%s: %d Sender laufen binnen 5 s weiter", meldung, senderZahl())
		}
		time.Sleep(time.Millisecond)
	}
}

// Abdeckung: LH-FA-20/Boundary — Schliesse verwirft die eingereihten Gruppen:
// Nach Schliesse sendet der Sender keine weitere Gruppe, auch wenn das Senden der
// laufenden erst danach gelingt (LH-FA-20.a *Gruppen*).
func TestEinspielGruppeSchliesseVerwirftWarteschlange(t *testing.T) {
	warteAufKeinenSender(t, "Sender früherer Tests")
	angefangen, frei := make(chan struct{}), make(chan struct{})
	var einmal sync.Once
	conn := neueSteuerConn(func(b []byte) (int, error) {
		einmal.Do(func() { close(angefangen) })
		<-frei
		return len(b), nil
	})
	s := postgres.NeueEinspielSession(conn)
	if err := gruppeBinnen(t, s, "erste Gruppe", model.ClientMessage{Type: model.ClientSync}); err != nil {
		t.Fatal(err)
	}
	select {
	case <-angefangen:
	case <-time.After(5 * time.Second):
		t.Fatal("die erste Gruppe beginnt binnen 5 s nicht zu senden")
	}
	if err := gruppeBinnen(t, s, "zweite Gruppe hinter der sendenden", model.ClientMessage{Type: model.ClientFlush}); err != nil {
		t.Fatal(err)
	}
	s.Schliesse()
	close(frei)
	warteAufKeinenSender(t, "nach Schliesse")
	if schreibt, _ := conn.zaehler(); schreibt != 1 {
		t.Fatalf("%d Schreibvorgänge, erwartet 1: Schliesse verwirft die zweite Gruppe", schreibt)
	}
}
