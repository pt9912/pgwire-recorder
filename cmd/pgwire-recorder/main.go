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

// version setzt der Build per -ldflags (Dockerfile, Stufe build).
var version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := bootstrap.Run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	stop()
	os.Exit(code)
}
