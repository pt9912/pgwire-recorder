package testzeit

import "time"

// Nach liefert einen Kanal, der nach ms Millisekunden ein Zeichen erhält.
func Nach(ms int64) <-chan time.Time {
	return time.After(dauer(ms))
}

// Schlafe hält den Aufrufer ms Millisekunden an.
func Schlafe(ms int64) {
	time.Sleep(dauer(ms))
}

// Frist ist ein Zeitpunkt, der ms Millisekunden nach ihrer Erzeugung liegt.
type Frist struct{ ende time.Time }

// NeueFrist liefert eine Frist, die ms Millisekunden nach dem Aufruf abläuft.
func NeueFrist(ms int64) Frist {
	return Frist{ende: time.Now().Add(dauer(ms))}
}

// Abgelaufen meldet, ob die Frist verstrichen ist.
func (f Frist) Abgelaufen() bool {
	return time.Now().After(f.ende)
}

func dauer(ms int64) time.Duration {
	return time.Duration(ms) * time.Millisecond
}
