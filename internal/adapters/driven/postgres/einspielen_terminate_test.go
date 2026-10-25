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

// Abdeckung: LH-FA-20/Boundary — ist die Anfrage vollständig geschrieben, sendet
// das zweite Signal Terminate, auch wenn Anfrage noch nicht zurückgekehrt ist:
// das Schließen sofort nach dem Write sendet Terminate hinter die Anfrage, ohne
// dass der Schreibvorgang als laufend gilt (LH-FA-20.a *Abbruchsignal*: auch bei
// nicht abgeschlossener Bearbeitung der Anfrage).
func TestEinspielSchliesseSendetNachVollstaendigerAnfrage(t *testing.T) {
	hook, geschrieben := geschriebenSignal()
	conn := &fristVerbindung{hook: hook}
	s := postgres.NeueEinspielSession(conn)
	go func() { _ = s.Anfrage("SELECT 1") }()
	bisSignal(t, geschrieben, "Anfrage schreibt binnen 5 s nicht")
	s.Schliesse()
	want := append(kodiertFrontend(t, &pgproto3.Query{String: "SELECT 1"}), kodiertFrontend(t, &pgproto3.Terminate{})...)
	if got := conn.stand(); !bytes.Equal(got, want) {
		t.Fatalf("geschrieben %q, erwartet die Anfrage und Terminate %q", got, want)
	}
}

// Abdeckung: LH-FA-20/Boundary — wie die Anfrage: ist eine eingereihte Gruppe
// vollständig geschrieben, sendet das zweite Signal Terminate hinter sie, auch
// wenn der Sender seinen Vorgang noch nicht beendet hat (LH-FA-20.a
// *Abbruchsignal*).
func TestEinspielSchliesseSendetNachVollstaendigerGruppe(t *testing.T) {
	hook, geschrieben := geschriebenSignal()
	conn := &fristVerbindung{hook: hook}
	s := postgres.NeueEinspielSession(conn)
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
// Schreiben von Terminate teilen sich eine Gesamtfrist von 100 ms ab dem zweiten
// Signal: Endet der Vorgang 80 ms nach dem Aufruf, liegt die Frist des
// Terminate vor dem Ende des Vorgangs plus 100 ms; Terminate wird gesendet
// (LH-FA-20.a *Abbruchsignal*: Gesamtfrist für Warten und Schreiben zusammen).
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
	s := postgres.NeueEinspielSession(conn)
	go func() { _ = s.Anfrage("SELECT 1") }()
	bisSignal(t, beginnt, "Anfrage beginnt binnen 5 s nicht zu schreiben")
	fertig := make(chan struct{})
	go func() {
		s.Schliesse()
		close(fertig)
	}()
	time.Sleep(80 * time.Millisecond)
	ende := time.Now()
	close(frei)
	bisSignal(t, fertig, "Schliesse endet binnen 5 s nicht")
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if len(conn.fristen) != 1 {
		t.Fatalf("Terminate mit %d Schreibfristen, erwartet 1 (Schreiben nach dem Warten)", len(conn.fristen))
	}
	if !conn.fristen[0].Before(ende.Add(100 * time.Millisecond)) {
		t.Fatalf("Frist des Terminate %v, erwartet vor %v: das Warten verbraucht die Gesamtfrist", conn.fristen[0], ende.Add(100*time.Millisecond))
	}
	if want := kodiertFrontend(t, &pgproto3.Terminate{}); !bytes.HasSuffix(conn.geschrieben, want) {
		t.Fatalf("geschrieben %q, erwartet Terminate am Ende", conn.geschrieben)
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
