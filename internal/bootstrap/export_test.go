package bootstrap

import (
	"io"
	"log/slog"
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
