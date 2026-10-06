package model

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Meldungscodes (SPEC-034), soweit dieser Stand sie erzeugt. Die erste Ziffer
// ist der Exit-Code der Fehlerklasse.
const (
	CodeInternal           = "PGR-E1000"
	CodeUsage              = "PGR-E2001"
	CodeOutputExists       = "PGR-E2002"
	CodeRecordingIO        = "PGR-E3001"
	CodeRecordingVersion   = "PGR-E3002"
	CodeRecordingBroken    = "PGR-E3003"
	CodeRecordingNoSession = "PGR-E3004"
	CodeNetwork            = "PGR-E4000"
	CodeListen             = "PGR-E4001"
	CodeUpstream           = "PGR-E4002"
	CodeConnectionLost     = "PGR-E4003"
	CodeReplayMismatch     = "PGR-E5001"
	CodeReplayUnconsumed   = "PGR-E5002"
	CodeReplaySession      = "PGR-E5003"
	CodeUnsupported        = "PGR-E6001"
	CodeProtocolVersion    = "PGR-E6002"
	CodeUnconsumed         = "PGR-W2001"
	CodeCancelRequest      = "PGR-W3001"
	CodeForeignProtocol    = "PGR-W3003"
)

// Klassen je Exit-Code (SPEC-034, Kopf „<klasse> [<code>]: <Ursache>“).
var klassen = map[int]string{
	1: "sonstiger Fehler",
	2: "Konfiguration",
	3: "Recording",
	4: "Netzwerk",
	5: "Replay",
	6: "nicht unterstützt",
}

// Error ist ein klassifizierter Fehler mit Meldungscode.
type Error struct {
	Code string
	Msg  string
	Err  error
}

// Errorf liefert einen Fehler mit Meldungscode.
func Errorf(code string, err error, format string, args ...any) *Error {
	return &Error{Code: code, Msg: fmt.Sprintf(format, args...), Err: err}
}

// Error ist der Fehlertext nach SPEC-034 §Ausgabe: eine Zeile, die mit dem Kopf
// dieses Fehlers beginnt; ein innerer klassifizierter Fehler der Kette trägt nur
// seine Ursache bei, ohne eigenen Kopf.
func (e *Error) Error() string {
	return einzeilig(kopf(e.Code) + ursache(e))
}

func (e *Error) Unwrap() error { return e.Err }

// ExitCode ist der Exit-Code der Fehlerklasse: die erste Ziffer des Codes.
func (e *Error) ExitCode() int { return exitCode(e.Code) }

func exitCode(code string) int {
	if len(code) >= 6 && code[5] >= '1' && code[5] <= '6' {
		return int(code[5] - '0')
	}
	return 1
}

// kopf ist der Kopf „<klasse> [<code>]: “ (SPEC-034 §Ausgabe).
func kopf(code string) string {
	return fmt.Sprintf("%s [%s]: ", klassen[exitCode(code)], code)
}

// ursache ist der Text eines Fehlers ohne Kopf: bei einem klassifizierten
// Fehler seine Meldung und die Ursache des inneren, bei nebeneinander
// entstandenen Fehlern (errors.Join) deren Ursachen durch "; " getrennt, sonst
// der Text des Fehlers, in dem ein darin eingebetteter klassifizierter Fehler
// ohne Kopf steht.
func ursache(err error) string {
	switch x := err.(type) {
	case *Error:
		if x.Err == nil {
			return x.Msg
		}
		return x.Msg + ": " + ursache(x.Err)
	case interface{ Unwrap() []error }:
		var teile []string
		for _, e := range x.Unwrap() {
			if e != nil {
				teile = append(teile, ursache(e))
			}
		}
		return strings.Join(teile, "; ")
	}
	text := err.Error()
	var me *Error
	if inner := errors.Unwrap(err); inner != nil && errors.As(inner, &me) {
		text = strings.Replace(text, inner.Error(), ursache(inner), 1)
	}
	return text
}

// umbruch ist ein Zeilenumbruch (LF, CR LF, einzelnes CR) mit dem Leerraum
// danach: Leerzeichen, Tabulatoren und weitere Zeilenumbrüche (SPEC-034
// §Ausgabe *Eine Zeile*). VT, FF, U+0085, U+2028 und U+2029 gehören nicht dazu.
var umbruch = regexp.MustCompile(`(\r\n|\r|\n)[ \t\r\n]*`)

// einzeilig ersetzt jeden Zeilenumbruch samt Leerraum danach durch ein
// Leerzeichen.
func einzeilig(s string) string { return umbruch.ReplaceAllString(s, " ") }

// Meldung ist eine Meldung des Recorders zu einem Fehler: Code und Fehlertext
// (SPEC-034 §Ausgabe).
type Meldung struct {
	Code string
	Text string
}

// ExitCode ist der Exit-Code der Klasse des Codes.
func (m Meldung) ExitCode() int { return exitCode(m.Code) }

// Meldungen ordnet einen Fehler nach SPEC-034 §Ausgabe ein und liefert
// mindestens eine Meldung. Jeder klassifizierte unter nebeneinander
// entstandenen Fehlern (errors.Join) ist eine eigene Meldung mit eigenem Kopf,
// in der Reihenfolge des Join; die erste ist die gemerkte. Ein nicht
// klassifizierter daneben folgt als Ursache der ersten klassifizierten; sind
// alle nicht klassifiziert, ist es eine Meldung PGR-E1000, die Texte durch "; "
// getrennt. In einem Fehler, der einen klassifizierten umhüllt, gilt dessen
// Code. Für nil liefert sie nil.
func Meldungen(err error) []Meldung {
	if err == nil {
		return nil
	}
	var out []Meldung
	var fremd []string
	for _, b := range nebeneinander(err) {
		var me *Error
		if errors.As(b, &me) {
			out = append(out, Meldung{Code: me.Code, Text: kopf(me.Code) + ursache(b)})
			continue
		}
		fremd = append(fremd, ursache(b))
	}
	switch {
	case len(out) == 0:
		out = []Meldung{{Code: CodeInternal, Text: kopf(CodeInternal) + strings.Join(fremd, "; ")}}
	case len(fremd) > 0:
		out[0].Text += ": " + strings.Join(fremd, "; ")
	}
	for i := range out {
		out[i].Text = einzeilig(out[i].Text)
	}
	return out
}

// nebeneinander zerlegt verschachtelte errors.Join in ihre nicht leeren
// Bestandteile, in deren Reihenfolge; jeder andere Fehler, auch ein
// klassifizierter mit innerem Fehler, ist ein Bestandteil.
func nebeneinander(err error) []error {
	j, ok := err.(interface{ Unwrap() []error })
	if !ok {
		return []error{err}
	}
	var out []error
	for _, e := range j.Unwrap() {
		if e != nil {
			out = append(out, nebeneinander(e)...)
		}
	}
	return out
}

// Warning ist eine Warnung mit Meldungscode (`PGR-W…`). Sie trägt keine
// Fehlerklasse und ändert den Exit-Code nicht (SPEC-034).
type Warning struct {
	Code string
	Msg  string
}

// Warnf liefert eine Warnung mit Meldungscode.
func Warnf(code string, format string, args ...any) *Warning {
	return &Warning{Code: code, Msg: fmt.Sprintf(format, args...)}
}

// ErrSessionEnded liefert ein Aufruf des Record-Use-Case, der auf eine schon
// beendete Session trifft oder durch ihr Beenden abgebrochen wird.
var ErrSessionEnded = errors.New("Session beendet")

// ErrShutdown liefert ein Aufruf des Record-Use-Case, wenn die Client-Nachricht
// nach dem Beginn des Herunterfahrens eine neue Interaktion begänne; sie wird
// nicht weitergeleitet, und die Session endet nach der laufenden Interaktion
// (LH-FA-13.a).
var ErrShutdown = errors.New("Herunterfahren: keine neue Interaktion")
