package pgwire_test

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/pgwire"
)

// fehlerLeser liefert bei jedem Lesen denselben Fehler.
type fehlerLeser struct{ err error }

func (f fehlerLeser) Read([]byte) (int, error) { return 0, f.err }

// TestLeserEndeDesStroms: Jeder der fünf Fehler des Datenstroms (auch eingepackt) ist
// das Ende der Client-Verbindung, jeder andere Fehler nicht.
func TestLeserEndeDesStroms(t *testing.T) {
	ende := map[string]error{
		"io.EOF":                 io.EOF,
		"io.ErrUnexpectedEOF":    io.ErrUnexpectedEOF,
		"net.ErrClosed":          net.ErrClosed,
		"syscall.ECONNRESET":     syscall.ECONNRESET,
		"os.ErrDeadlineExceeded": os.ErrDeadlineExceeded,
	}
	for name, err := range ende {
		t.Run(name, func(t *testing.T) {
			if got, _ := pgwire.LeseFehler(fehlerLeser{fmt.Errorf("lesen: %w", err)}); !got {
				t.Errorf("%s ist nicht als Ende der Verbindung gekennzeichnet", name)
			}
		})
	}
	t.Run("anderer Fehler", func(t *testing.T) {
		andere := errors.New("anderer Fehler")
		if got, err := pgwire.LeseFehler(fehlerLeser{andere}); got || !errors.Is(err, andere) {
			t.Errorf("ende = %v, err = %v; ein anderer Fehler ist kein Ende und bleibt erhalten", got, err)
		}
	})
}
