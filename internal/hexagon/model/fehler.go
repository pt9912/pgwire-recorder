package model

import "fmt"

// Meldungscodes (SPEC-034), soweit dieser Stand sie erzeugt. Die erste Ziffer
// ist der Exit-Code der Fehlerklasse.
const (
	CodeInternal         = "PGR-E1000"
	CodeUsage            = "PGR-E2001"
	CodeOutputExists     = "PGR-E2002"
	CodeRecordingIO      = "PGR-E3001"
	CodeRecordingVersion = "PGR-E3002"
	CodeRecordingBroken  = "PGR-E3003"
	CodeNetwork          = "PGR-E4000"
	CodeListen           = "PGR-E4001"
	CodeUpstream         = "PGR-E4002"
	CodeConnectionLost   = "PGR-E4003"
	CodeUnsupported      = "PGR-E6001"
	CodeProtocolVersion  = "PGR-E6002"
	CodeCancelRequest    = "PGR-W3001"
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

func (e *Error) Error() string {
	kopf := fmt.Sprintf("%s [%s]: %s", klassen[e.ExitCode()], e.Code, e.Msg)
	if e.Err != nil {
		return kopf + ": " + e.Err.Error()
	}
	return kopf
}

func (e *Error) Unwrap() error { return e.Err }

// ExitCode ist der Exit-Code der Fehlerklasse: die erste Ziffer des Codes.
func (e *Error) ExitCode() int {
	if len(e.Code) >= 6 && e.Code[5] >= '1' && e.Code[5] <= '6' {
		return int(e.Code[5] - '0')
	}
	return 1
}
