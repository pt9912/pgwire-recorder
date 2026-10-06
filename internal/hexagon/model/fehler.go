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
// Fehler seine Meldung und die Ursache des inneren, bei einer reinen
// Zusammenfassung (zusammenfassung) die Ursachen ihrer Teile durch "; "
// getrennt, sonst der ganze Text des Fehlers, in dem jeder innere Fehler mit
// seiner Ursache statt seines Textes steht, also ohne inneren Kopf.
func ursache(err error) string {
	if x, ok := err.(*Error); ok {
		if x.Err == nil {
			return x.Msg
		}
		return x.Msg + ": " + ursache(x.Err)
	}
	if teile, ok := zusammenfassung(err); ok {
		texte := make([]string, len(teile))
		for i, e := range teile {
			texte[i] = ursache(e)
		}
		return strings.Join(texte, "; ")
	}
	text := err.Error()
	for _, inner := range innere(err) {
		text = strings.Replace(text, inner.Error(), ursache(inner), 1)
	}
	return text
}

// innere liefert die nicht leeren inneren Fehler, ob über Unwrap() error oder
// Unwrap() []error, in ihrer Reihenfolge.
func innere(err error) []error {
	switch x := err.(type) {
	case interface{ Unwrap() error }:
		if e := x.Unwrap(); e != nil {
			return []error{e}
		}
	case interface{ Unwrap() []error }:
		var out []error
		for _, e := range x.Unwrap() {
			if e != nil {
				out = append(out, e)
			}
		}
		return out
	}
	return nil
}

// zusammenfassung meldet, ob ein Fehler mehrere andere ohne eigenen Text
// zusammenfasst, wie errors.Join: Er hat Unwrap() []error, und sein Text ist
// genau der seiner nicht leeren Teile, durch Zeilenumbrüche getrennt. Eine Hülle
// mit eigenem Text um mehrere Ursachen, etwa fmt.Errorf mit mehreren %w, ist
// keine; sie ist eine Kette (SPEC-034 §Ausgabe *Fehlerkette*).
func zusammenfassung(err error) ([]error, bool) {
	if _, ok := err.(interface{ Unwrap() []error }); !ok {
		return nil, false
	}
	teile := innere(err)
	texte := make([]string, len(teile))
	for i, e := range teile {
		texte[i] = e.Error()
	}
	if err.Error() != strings.Join(texte, "\n") {
		return nil, false
	}
	return teile, true
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
// mindestens eine Meldung. Jeder klassifizierte Teil einer reinen
// Zusammenfassung (errors.Join) ist eine eigene Meldung mit eigenem Kopf, in
// ihrer Reihenfolge; die erste ist die gemerkte. Ein nicht klassifizierter
// daneben folgt als Ursache der ersten klassifizierten; sind alle nicht
// klassifiziert, ist es eine Meldung PGR-E1000, die Texte durch "; " getrennt.
// Jeder andere Fehler ist eine Kette: eine Meldung mit seinem ganzen Text und
// dem Code des ersten klassifizierten Fehlers unter ihm, außen nach innen und
// in der Reihenfolge seiner Ursachen, sonst PGR-E1000. Für nil liefert sie nil.
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

// nebeneinander zerlegt verschachtelte reine Zusammenfassungen
// (zusammenfassung) in ihre nicht leeren Bestandteile, in deren Reihenfolge;
// jeder andere Fehler, auch ein klassifizierter mit innerem Fehler und eine
// Hülle mit eigenem Text um mehrere Ursachen, ist ein Bestandteil.
func nebeneinander(err error) []error {
	teile, ok := zusammenfassung(err)
	if !ok {
		return []error{err}
	}
	var out []error
	for _, e := range teile {
		out = append(out, nebeneinander(e)...)
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
