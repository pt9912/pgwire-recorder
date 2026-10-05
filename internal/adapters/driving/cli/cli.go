package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// RecordOptions sind die Optionen von `record` (LH-FA-02.a), soweit dieser Stand
// sie kennt.
type RecordOptions struct {
	Listen   string
	Upstream string
	Output   string
	Force    bool
}

// ReplayOptions sind die Optionen von `replay` (LH-FA-03.a), soweit dieser Stand
// sie kennt.
type ReplayOptions struct {
	Listen           string
	Input            string
	FailOnUnconsumed bool
}

// envFailOnUnconsumed ist die Umgebungsvariable von --fail-on-unconsumed
// (LH-FA-17.a).
const envFailOnUnconsumed = "PGWIRE_RECORDER_FAIL_ON_UNCONSUMED"

// Command ist das gewählte Kommando mit seinen Optionen.
type Command struct {
	Name   string
	Record RecordOptions
	Replay ReplayOptions
}

const usage = `Aufruf: pgwire-recorder <kommando> [optionen]

Kommandos:
  record   vermittelt Clients zu PostgreSQL und zeichnet die Kommunikation auf
  replay   beantwortet Anfragen aus einer Aufzeichnung, ohne PostgreSQL
  version  gibt die Programmversion aus

Optionen von record:
  --listen    Adresse, auf der Clients angenommen werden (Pflicht)
  --upstream  Adresse des PostgreSQL-Servers, host:port (Pflicht)
  --output    Zieldatei der Aufzeichnung (Pflicht)
  --force     vorhandene Zieldatei ersetzen

Optionen von replay:
  --listen    Adresse, auf der Clients angenommen werden (Pflicht)
  --input     Aufzeichnung (Pflicht)
  --fail-on-unconsumed[=true|false]
              nicht verbrauchte Interaktionen und nie zugeordnete Sessions
              sind ein Fehler (PGR-E5002, Exit-Code 5) statt einer Warnung;
              Umgebungsvariable PGWIRE_RECORDER_FAIL_ON_UNCONSUMED
`

// ErrHelp meldet, dass die Hilfe angefordert und ausgegeben wurde.
var ErrHelp = errors.New("Hilfe ausgegeben")

// Parse liest Kommando und Optionen. Eine ungültige Verwendung ist PGR-E2001;
// --help und -h geben die Hilfe auf out aus und liefern ErrHelp.
func Parse(args []string, out io.Writer) (Command, error) {
	if len(args) == 0 {
		return Command{}, model.Errorf(model.CodeUsage, nil, "kein Kommando angegeben; --help zeigt die Kommandos")
	}
	switch args[0] {
	case "-h", "--help", "help":
		fmt.Fprint(out, usage)
		return Command{}, ErrHelp
	case "record":
		return parseRecord(args[1:], out)
	case "replay":
		return parseReplay(args[1:], out)
	case "version":
		if len(args) > 1 {
			return Command{}, model.Errorf(model.CodeUsage, nil, "version nimmt keine Argumente")
		}
		return Command{Name: "version"}, nil
	default:
		return Command{}, model.Errorf(model.CodeUsage, nil, "unbekanntes Kommando %q; --help zeigt die Kommandos", args[0])
	}
}

func parseRecord(args []string, out io.Writer) (Command, error) {
	fs := flag.NewFlagSet("record", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var o RecordOptions
	fs.StringVar(&o.Listen, "listen", "", "")
	fs.StringVar(&o.Upstream, "upstream", "", "")
	fs.StringVar(&o.Output, "output", "", "")
	fs.BoolVar(&o.Force, "force", false, "")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(out, usage)
			return Command{}, ErrHelp
		}
		return Command{}, model.Errorf(model.CodeUsage, err, "ungültige Verwendung von record")
	}
	if fs.NArg() > 0 {
		return Command{}, model.Errorf(model.CodeUsage, nil, "unerwartetes Argument %q", fs.Arg(0))
	}
	for _, p := range []struct{ name, value string }{{"--listen", o.Listen}, {"--upstream", o.Upstream}, {"--output", o.Output}} {
		if p.value == "" {
			return Command{}, model.Errorf(model.CodeUsage, nil, "Pflichtoption %s fehlt", p.name)
		}
	}
	return Command{Name: "record", Record: o}, nil
}

func parseReplay(args []string, out io.Writer) (Command, error) {
	fs := flag.NewFlagSet("replay", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var o ReplayOptions
	fs.StringVar(&o.Listen, "listen", "", "")
	fs.StringVar(&o.Input, "input", "", "")
	fail := wahrheitswert{&o.FailOnUnconsumed}
	// Die Umgebungsvariable setzt den Wert vor der Kommandozeile, die ihn
	// danach überschreibt (CLI vor Umgebungsvariable); leer gilt sie als nicht
	// gesetzt (LH-FA-17.a).
	if v := os.Getenv(envFailOnUnconsumed); v != "" {
		if err := fail.Set(v); err != nil {
			return Command{}, model.Errorf(model.CodeUsage, err, "Umgebungsvariable %s", envFailOnUnconsumed)
		}
	}
	fs.Var(fail, "fail-on-unconsumed", "")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(out, usage)
			return Command{}, ErrHelp
		}
		return Command{}, model.Errorf(model.CodeUsage, err, "ungültige Verwendung von replay")
	}
	if fs.NArg() > 0 {
		return Command{}, model.Errorf(model.CodeUsage, nil, "unerwartetes Argument %q", fs.Arg(0))
	}
	for _, p := range []struct{ name, value string }{{"--listen", o.Listen}, {"--input", o.Input}} {
		if p.value == "" {
			return Command{}, model.Errorf(model.CodeUsage, nil, "Pflichtoption %s fehlt", p.name)
		}
	}
	return Command{Name: "replay", Replay: o}, nil
}

// wahrheitswert ist der Wert einer booleschen Option nach LH-FA-17.a: ohne
// Wert true, mit Wert genau "true" oder "false"; jeder andere Wert, auch der
// leere und "1", ist ein Fehler. Nennt die Kommandozeile die Option mehrfach,
// gilt die letzte Angabe.
type wahrheitswert struct{ wert *bool }

// IsBoolFlag lässt die Option ohne Wert zu; flag setzt dann "true".
func (w wahrheitswert) IsBoolFlag() bool { return true }

func (w wahrheitswert) String() string {
	if w.wert == nil {
		return "false"
	}
	return fmt.Sprint(*w.wert)
}

func (w wahrheitswert) Set(v string) error {
	switch v {
	case "true":
		*w.wert = true
	case "false":
		*w.wert = false
	default:
		return errors.New("erlaubt sind true und false")
	}
	return nil
}
