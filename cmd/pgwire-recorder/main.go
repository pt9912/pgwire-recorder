// Command pgwire-recorder zeichnet PostgreSQL-Kommunikation auf, gibt sie ohne
// Datenbank wieder und spielt sie in eine Datenbank ein (ARC-009).
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
)

// version setzt der Build per -ldflags (Dockerfile, Stufe build); das Kommando
// `version` gibt sie aus.
var version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// Nach dem ersten Signal gilt wieder das Standardverhalten: ein zweites
	// beendet den Prozess sofort.
	go func() {
		<-ctx.Done()
		stop()
	}()
	code := bootstrap.Run(ctx, os.Args[1:], version, os.Stdout, os.Stderr)
	stop()
	os.Exit(code)
}
