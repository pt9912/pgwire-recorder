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

// main behandelt SIGINT und SIGTERM, unter Windows auch den Konsolenabbruch
// (os.Interrupt), für die ganze Laufzeit (LH-FA-13.a): Das erste Signal beendet
// ctx und beginnt das Herunterfahren, das zweite schließt ablauf und lässt die
// Frist sofort ablaufen, jedes weitere bleibt ohne Wirkung. Das erste Signal ist
// aus dem Kanal gelesen, bevor ctx endet; ein zweites, das danach eintrifft,
// wartet darum im Kanal.
func main() {
	signale := make(chan os.Signal, 1)
	signal.Notify(signale, os.Interrupt, syscall.SIGTERM)
	ctx, herunterfahren := context.WithCancel(context.Background())
	ablauf := make(chan struct{})
	go func() {
		<-signale
		herunterfahren()
		<-signale
		close(ablauf)
	}()
	code := bootstrap.Run(ctx, ablauf, os.Args[1:], version, os.Stdout, os.Stderr)
	herunterfahren()
	os.Exit(code)
}
