package bootstrap

import (
	"io"
	"log/slog"
	"strings"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// Logger reicht an logger weiter.
func Logger(stderr io.Writer, stufe string) *slog.Logger {
	return logger(stderr, stufe)
}

// Fail reicht an fail weiter.
func Fail(stderr io.Writer, err error) int {
	return fail(stderr, err)
}

// ExitCode reicht an exitCode weiter.
func ExitCode(code string) int {
	return exitCode(code)
}

// UrsprungMonoton meldet, ob der Ursprung der Systemuhr einen Messwert der
// monotonen Zeit trägt.
func UrsprungMonoton(u driven.Uhr) bool {
	return strings.Contains(u.(*systemuhr).ursprung.String(), " m=")
}
