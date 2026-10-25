package postgres_test

import (
	"bytes"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// fristVerbindung schreibt jedes Write auf Wunsch von hook und hält die
// Schreibfristen fest, die SetWriteDeadline bekommt, mit dem Ende jedes Write.
type fristVerbindung struct {
	net.Conn
	mu          sync.Mutex
	geschrieben []byte
	fristen     []time.Time
	hook        func(n int)
	writes      int
}

func (c *fristVerbindung) Write(b []byte) (int, error) {
	c.mu.Lock()
	c.geschrieben = append(c.geschrieben, b...)
	c.writes++
	n := c.writes
	c.mu.Unlock()
	if c.hook != nil {
		c.hook(n)
	}
	return len(b), nil
}

func (c *fristVerbindung) SetWriteDeadline(t time.Time) error {
	c.mu.Lock()
	c.fristen = append(c.fristen, t)
	c.mu.Unlock()
	return nil
}

func (c *fristVerbindung) Close() error { return nil }

func (c *fristVerbindung) stand() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]byte(nil), c.geschrieben...)
}

// geschriebenSignal ist ein hook, der beim ersten Write geschrieben schließt.
func geschriebenSignal() (func(int), <-chan struct{}) {
	ch := make(chan struct{})
	return func(n int) {
		if n == 1 {
			close(ch)
		}
	}, ch
}

func bisSignal(t *testing.T, ch <-chan struct{}, was string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatal(was)
	}
}

// langeFrist ist die Gesamtfrist der Tests, die nur die Reihenfolge prüfen: Sie
// liegt um Größenordnungen über der Zeit zwischen Write und Ende des
// Schreibvorgangs, damit Last sie nicht reißt.
const langeFrist = time.Second

// Abdeckung: LH-FA-20/Boundary — ist die Anfrage vollständig geschrieben, sendet
// das zweite Signal Terminate hinter die Anfrage, sobald der Schreibvorgang
// endet; Schliesse wartet darauf höchstens die Gesamtfrist (hier eine Sekunde,
// LH-FA-20.a *Abbruchsignal*: auch bei nicht abgeschlossener Bearbeitung der
// Anfrage).
func TestEinspielSchliesseSendetNachVollstaendigerAnfrage(t *testing.T) {
	hook, geschrieben := geschriebenSignal()
	conn := &fristVerbindung{hook: hook}
	s := postgres.NeueEinspielSessionMitFrist(conn, langeFrist)
	go func() { _ = s.Anfrage("SELECT 1") }()
	bisSignal(t, geschrieben, "Anfrage schreibt binnen 5 s nicht")
	s.Schliesse()
	want := append(kodiertFrontend(t, &pgproto3.Query{String: "SELECT 1"}), kodiertFrontend(t, &pgproto3.Terminate{})...)
	if got := conn.stand(); !bytes.Equal(got, want) {
		t.Fatalf("geschrieben %q, erwartet die Anfrage und Terminate %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Boundary — wie die Anfrage: ist eine eingereihte Gruppe
// vollständig geschrieben, sendet das zweite Signal Terminate hinter sie,
// sobald der Schreibvorgang des Senders endet; Schliesse wartet darauf
// höchstens die Gesamtfrist (hier eine Sekunde, LH-FA-20.a *Abbruchsignal*).
func TestEinspielSchliesseSendetNachVollstaendigerGruppe(t *testing.T) {
	hook, geschrieben := geschriebenSignal()
	conn := &fristVerbindung{hook: hook}
	s := postgres.NeueEinspielSessionMitFrist(conn, langeFrist)
	if err := gruppeBinnen(t, s, "Gruppe", model.ClientMessage{Type: model.ClientSync}); err != nil {
		t.Fatal(err)
	}
	bisSignal(t, geschrieben, "der Sender schreibt die Gruppe binnen 5 s nicht")
	s.Schliesse()
	want := append(kodiertFrontend(t, &pgproto3.Sync{}), kodiertFrontend(t, &pgproto3.Terminate{})...)
	if got := conn.stand(); !bytes.Equal(got, want) {
		t.Fatalf("geschrieben %q, erwartet die Gruppe und Terminate %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Boundary — Warten auf einen laufenden Schreibvorgang und
// Schreiben von Terminate teilen sich eine Gesamtfrist ab dem zweiten Signal:
// Endet der Vorgang 200 ms, nachdem das Schließen begonnen hat, liegt die Frist
// des Terminate nicht nach dem Ende des Vorgangs plus der Gesamtfrist (hier
// eine Sekunde), und Terminate wird gesendet (LH-FA-20.a *Abbruchsignal*:
// Gesamtfrist für Warten und Schreiben zusammen).
func TestEinspielSchliesseGesamtfrist(t *testing.T) {
	frei := make(chan struct{})
	beginnt := make(chan struct{})
	conn := &fristVerbindung{}
	conn.hook = func(n int) {
		if n == 1 {
			close(beginnt)
			<-frei
		}
	}
	s := postgres.NeueEinspielSessionMitFrist(conn, langeFrist)
	go func() { _ = s.Anfrage("SELECT 1") }()
	bisSignal(t, beginnt, "Anfrage beginnt binnen 5 s nicht zu schreiben")
	fertig := make(chan struct{})
	go func() {
		s.Schliesse()
		close(fertig)
	}()
	// Eine Gruppe nach dem Schließen wird abgelehnt: Das Schließen hat dann
	// begonnen und seine Frist gesetzt.
	bis := time.Now().Add(5 * time.Second)
	for s.Gruppe([]model.ClientMessage{{Type: model.ClientSync}}) == nil {
		if time.Now().After(bis) {
			t.Fatal("das Schließen beginnt binnen 5 s nicht")
		}
		time.Sleep(time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond)
	ende := time.Now()
	close(frei)
	bisSignal(t, fertig, "Schliesse endet binnen 5 s nicht")
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if len(conn.fristen) != 1 {
		t.Fatalf("Terminate mit %d Schreibfristen, erwartet 1 (Schreiben nach dem Warten)", len(conn.fristen))
	}
	if grenze := ende.Add(langeFrist); conn.fristen[0].After(grenze) {
		t.Fatalf("Frist des Terminate %v, erwartet nicht nach %v: das Warten verbraucht die Gesamtfrist nicht", conn.fristen[0], grenze)
	}
	if want := kodiertFrontend(t, &pgproto3.Terminate{}); !bytes.HasSuffix(conn.geschrieben, want) {
		t.Fatalf("geschrieben %q, erwartet Terminate am Ende", conn.geschrieben)
	}
}

// Abdeckung: LH-FA-20/Boundary — ist der Schreibvorgang beendet, bevor das
// Warten beginnt, entscheidet nicht die Frist: Mit abgelaufener Frist liefert
// das Warten wahr, wenn der Vorgang beendet ist, sonst falsch; ohne Vorgang
// wahr (LH-FA-20.a *Abbruchsignal*: Terminate, sobald die Nachricht vollständig
// geschrieben ist).
func TestEinspielWartePunktlich(t *testing.T) {
	vorbei := time.Now().Add(-time.Second)
	spaeter := time.Now().Add(time.Hour)
	beendet := make(chan struct{})
	close(beendet)
	for i := 0; i < 200; i++ {
		if !postgres.WartePunktlich(beendet, vorbei) {
			t.Fatal("beendeter Vorgang bei abgelaufener Frist: falsch, erwartet wahr")
		}
	}
	if !postgres.WartePunktlich(beendet, spaeter) {
		t.Fatal("beendeter Vorgang vor der Frist: falsch, erwartet wahr")
	}
	if !postgres.WartePunktlich(nil, vorbei) {
		t.Fatal("kein Vorgang: falsch, erwartet wahr")
	}
	laufend := make(chan bool, 1)
	go func() { laufend <- postgres.WartePunktlich(make(chan struct{}), vorbei) }()
	select {
	case got := <-laufend:
		if got {
			t.Fatal("laufender Vorgang bei abgelaufener Frist: wahr, erwartet falsch")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("das Warten endet binnen 5 s nicht, obwohl die Frist abgelaufen ist")
	}
}

// Abdeckung: LH-FA-20/Boundary — nach dem Schließen schreibt Anfrage nichts
// mehr: kein Byte der Anfrage hinter das Terminate, der Fehler ist PGR-E4003
// (LH-FA-20.a *Abbruchsignal*).
func TestEinspielAnfrageNachSchliessen(t *testing.T) {
	fuerBeideArten(t, func(t *testing.T, schliesse func(driven.EinspielSession)) {
		conn := &fristVerbindung{}
		s := postgres.NeueEinspielSession(conn)
		schliesse(s)
		vorher := conn.stand()
		if err := s.Anfrage("SELECT 1"); code(err) != model.CodeConnectionLost {
			t.Fatalf("Anfrage nach dem Schließen: %v, erwartet %s", err, model.CodeConnectionLost)
		}
		if got := conn.stand(); !bytes.Equal(got, vorher) {
			t.Fatalf("nach dem Schließen geschrieben %q, erwartet nichts", got[len(vorher):])
		}
	})
}
