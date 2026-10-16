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

// fuerBeideArten führt f einmal mit Schliesse und einmal mit
// SchliesseOhneTerminate als Teilfall aus.
func fuerBeideArten(t *testing.T, f func(t *testing.T, schliesse func(driven.EinspielSession))) {
	t.Helper()
	for _, a := range []struct {
		name      string
		schliesse func(driven.EinspielSession)
	}{
		{"Schliesse", func(s driven.EinspielSession) { s.Schliesse() }},
		{"SchliesseOhneTerminate", func(s driven.EinspielSession) { s.SchliesseOhneTerminate() }},
	} {
		t.Run(a.name, func(t *testing.T) { f(t, a.schliesse) })
	}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-25/Happy — SchliesseOhneTerminate
// schließt die Verbindung, ohne dass der Server nach der Anfrage ein Byte liest:
// kein Terminate auf der Leitung, der Server sieht das Verbindungsende
// (LH-FA-20.a *Ende einer Session*).
func TestEinspielSchliesseOhneTerminate(t *testing.T) {
	antwort := kodiert(t, &pgproto3.ErrorResponse{Severity: "FATAL", Code: "57P01", Message: "beendet"})
	addr, ergebnis := einspielServer(t, verbunden(t), antwort)
	s := verbinde(t, addr)
	if err := s.Anfrage("SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if r, err := s.Naechste(); err != nil || r.Type != model.ResponseErrorResponse {
		t.Fatalf("Naechste: %v, %v, erwartet die Fehlerantwort", r.Type, err)
	}
	s.SchliesseOhneTerminate()
	if l := empfangen(t, ergebnis); len(l.danach) != 0 {
		t.Fatalf("nach der Anfrage beim Schließen gesendet %q, erwartet nichts", l.danach)
	}
}

// zaehlVerbindung hält fest, was auf die Verbindung geschrieben wird; sie
// nimmt jedes Write an. Nur Write, Close und SetWriteDeadline sind belegt.
type zaehlVerbindung struct {
	net.Conn
	mu          sync.Mutex
	geschrieben []byte
}

func (c *zaehlVerbindung) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.geschrieben = append(c.geschrieben, b...)
	return len(b), nil
}

func (c *zaehlVerbindung) Close() error { return nil }

func (c *zaehlVerbindung) SetWriteDeadline(time.Time) error { return nil }

// Abdeckung: LH-FA-20/Boundary, LH-FA-25/Boundary — nur das erste Schließen
// wirkt und nur dort sendet Schliesse ein Terminate: Nach Schliesse oder
// SchliesseOhneTerminate sendet ein weiteres Schließen, gleich welcher Art,
// nichts; SchliesseOhneTerminate sendet nie eines (LH-FA-20.a *Ende einer
// Session*).
func TestEinspielSchliesseNurEinmal(t *testing.T) {
	terminate := kodiertFrontend(t, &pgproto3.Terminate{})
	for _, f := range []struct {
		name   string
		ablauf []string
		want   []byte
	}{
		{"Schliesse, Schliesse", []string{"mit", "mit"}, terminate},
		{"Schliesse, SchliesseOhneTerminate", []string{"mit", "ohne"}, terminate},
		{"SchliesseOhneTerminate, Schliesse", []string{"ohne", "mit"}, nil},
		{"SchliesseOhneTerminate, SchliesseOhneTerminate", []string{"ohne", "ohne"}, nil},
	} {
		t.Run(f.name, func(t *testing.T) {
			conn := &zaehlVerbindung{}
			s := postgres.NeueEinspielSession(conn)
			for _, art := range f.ablauf {
				if art == "mit" {
					s.Schliesse()
				} else {
					s.SchliesseOhneTerminate()
				}
			}
			conn.mu.Lock()
			defer conn.mu.Unlock()
			if !bytes.Equal(conn.geschrieben, f.want) {
				t.Fatalf("geschrieben %q, erwartet %q", conn.geschrieben, f.want)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — Schliesse und SchliesseOhneTerminate beenden
// ein wartendes Naechste: Es kehrt binnen 5 s mit dem Verbindungsende
// (PGR-E4003) zurück (LH-FA-20.a *Ende einer Session*).
func TestEinspielSchliesseBeendetNaechste(t *testing.T) {
	fuerBeideArten(t, func(t *testing.T, schliesse func(driven.EinspielSession)) {
		conn := neueSteuerConn(func(b []byte) (int, error) { return len(b), nil })
		s := postgres.NeueEinspielSession(conn)
		fertig := make(chan error, 1)
		go func() {
			_, err := s.Naechste()
			fertig <- err
		}()
		time.Sleep(50 * time.Millisecond)
		schliesse(s)
		select {
		case err := <-fertig:
			if code(err) != model.CodeConnectionLost {
				t.Fatalf("Naechste nach dem Schließen: %v, erwartet %s", err, model.CodeConnectionLost)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Naechste endet binnen 5 s nach dem Schließen nicht")
		}
	})
}
