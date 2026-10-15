package pgwire

import (
	"errors"
	"io"
	"net"
	"os"
	"syscall"
)

// streamFehler kennzeichnet einen Fehler, den der Datenstrom der Client-Verbindung
// selbst liefert und der sein Ende ist (Schließen, Abbruch, Frist). Ein Fehler, den
// die PGWire-Bibliothek beim Lesen eines vollständig gelesenen Rumpfs erzeugt, trägt
// ihn nicht, auch wenn er io.EOF ist (LH-FA-05.e, *Ungültiger Rumpf und abgebrochene
// Nachricht*).
type streamFehler struct{ err error }

// Error liefert den Text des Fehlers des Datenstroms.
func (e *streamFehler) Error() string { return e.err.Error() }

// Unwrap liefert den Fehler des Datenstroms.
func (e *streamFehler) Unwrap() error { return e.err }

// leser liest von der Client-Verbindung und kennzeichnet deren Ende als
// streamFehler.
type leser struct{ r io.Reader }

// Read liest von r; ein Fehler, der das Ende der Verbindung ist, kommt als
// streamFehler zurück.
func (l leser) Read(p []byte) (int, error) {
	n, err := l.r.Read(p)
	if err != nil && endeDesStroms(err) {
		err = &streamFehler{err}
	}
	return n, err
}

// endeDesStroms meldet, ob ein Lesefehler der Verbindung ihr Ende ist.
func endeDesStroms(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, net.ErrClosed) || errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, os.ErrDeadlineExceeded)
}

// verbindungsende meldet, ob ein Lesefehler von Receive das Ende der Client-Verbindung
// ist und nicht eine unlesbare Nachricht.
func verbindungsende(err error) bool {
	var s *streamFehler
	return errors.As(err, &s)
}
