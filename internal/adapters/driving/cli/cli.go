package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// RecordOptions sind die Optionen von `record` (LH-FA-02.a), soweit dieser Stand
// sie kennt.
type RecordOptions struct {
	Listen   string
	Upstream string
	Output   string
	Force    bool
	LogLevel string
}

// ReplayOptions sind die Optionen von `replay` (LH-FA-03.a), soweit dieser Stand
// sie kennt.
type ReplayOptions struct {
	Listen           string
	Input            string
	FailOnUnconsumed bool
	LogLevel         string
}

// envFailOnUnconsumed ist die Umgebungsvariable von --fail-on-unconsumed
// (LH-FA-17.a).
const envFailOnUnconsumed = "PGWIRE_RECORDER_FAIL_ON_UNCONSUMED"

// envLogLevel ist die Umgebungsvariable von --log-level (LH-FA-14.a).
const envLogLevel = "PGWIRE_RECORDER_LOG_LEVEL"

// Stufen von --log-level (LH-FA-14.a); Standard ist LogInfo (SPEC-005).
const (
	LogError = "error"
	LogWarn  = "warn"
	LogInfo  = "info"
	LogDebug = "debug"
)

// Command ist das gewählte Kommando mit seinen Optionen.
type Command struct {
	Name   string
	Record RecordOptions
	Replay ReplayOptions
}

const optionenRecord = `Optionen von record:
  --listen    Adresse, auf der Clients angenommen werden (Pflicht)
  --upstream  Adresse des PostgreSQL-Servers, host:port (Pflicht)
  --output    Zieldatei der Aufzeichnung (Pflicht)
  --force     vorhandene Zieldatei ersetzen
` + optionLogLevel

const optionenReplay = `Optionen von replay:
  --listen    Adresse, auf der Clients angenommen werden (Pflicht)
  --input     Aufzeichnung (Pflicht)
  --fail-on-unconsumed[=true|false]
              nicht verbrauchte Interaktionen und nie zugeordnete Sessions
              sind ein Fehler (PGR-E5002, Exit-Code 5) statt einer Warnung;
              Umgebungsvariable PGWIRE_RECORDER_FAIL_ON_UNCONSUMED
` + optionLogLevel

const optionLogLevel = `  --log-level error|warn|info|debug
              Stufe der Log-Zeilen auf stderr, Standard info; eine Stufe zeigt
              auch die strengeren; Umgebungsvariable PGWIRE_RECORDER_LOG_LEVEL
`

// usage ist die globale Hilfe.
const usage = `Aufruf: pgwire-recorder <kommando> [optionen]

Kommandos:
  record   vermittelt Clients zu PostgreSQL und zeichnet die Kommunikation auf
  replay   beantwortet Anfragen aus einer Aufzeichnung, ohne PostgreSQL
  version  gibt die Programmversion aus

` + optionenRecord + `
` + optionenReplay

// hilfen ist die Hilfe je bekanntem Kommando (LH-FA-01.a).
var hilfen = map[string]string{
	"record":  "Aufruf: pgwire-recorder record [optionen]\n\n" + optionenRecord,
	"replay":  "Aufruf: pgwire-recorder replay [optionen]\n\n" + optionenReplay,
	"version": "Aufruf: pgwire-recorder version\n\nGibt die Programmversion aus.\n",
}

// ErrHelp meldet, dass die Hilfe angefordert und ausgegeben wurde.
var ErrHelp = errors.New("Hilfe ausgegeben")

// Parse liest Kommando und Optionen. Eine ungültige Verwendung ist PGR-E2001.
// Steht vor dem ersten "--" eine Hilfe-Angabe (hilfeAngabe), an welcher Stelle
// auch immer, gibt Parse vor jeder anderen Prüfung die Hilfe auf out aus und
// liefert ErrHelp: die des Kommandos, wenn das erste Argument ein bekanntes
// Kommando ist, sonst die globale (LH-FA-01.a). Das Kommando help gibt die
// globale Hilfe aus.
func Parse(args []string, out io.Writer) (Command, error) {
	if hilfeVerlangt(args) {
		text := usage
		if len(args) > 0 {
			if h, ok := hilfen[args[0]]; ok {
				text = h
			}
		}
		fmt.Fprint(out, text)
		return Command{}, ErrHelp
	}
	if len(args) == 0 {
		return Command{}, model.Errorf(model.CodeUsage, nil, "kein Kommando angegeben; --help zeigt die Kommandos")
	}
	switch args[0] {
	case "help":
		fmt.Fprint(out, usage)
		return Command{}, ErrHelp
	case "record":
		return parseRecord(args[1:])
	case "replay":
		return parseReplay(args[1:])
	case "version":
		if len(args) > 1 {
			return Command{}, model.Errorf(model.CodeUsage, nil, "version nimmt keine Argumente")
		}
		return Command{Name: "version"}, nil
	default:
		return Command{}, model.Errorf(model.CodeUsage, nil, "unbekanntes Kommando %q; --help zeigt die Kommandos", args[0])
	}
}

func parseRecord(args []string) (Command, error) {
	fs := flag.NewFlagSet("record", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var o RecordOptions
	fs.StringVar(&o.Listen, "listen", "", "")
	fs.StringVar(&o.Upstream, "upstream", "", "")
	fs.StringVar(&o.Output, "output", "", "")
	fs.BoolVar(&o.Force, "force", false, "")
	level, err := logLevelOption(fs)
	if err != nil {
		return Command{}, err
	}
	optionen, rest := endeDerOptionen(args)
	if err := fs.Parse(optionen); err != nil {
		return Command{}, model.Errorf(model.CodeUsage, err, "ungültige Verwendung von record")
	}
	if rest = append(fs.Args(), rest...); len(rest) > 0 {
		return Command{}, model.Errorf(model.CodeUsage, nil, "unerwartetes Argument %q", rest[0])
	}
	for _, p := range []struct{ name, value string }{{"--listen", o.Listen}, {"--upstream", o.Upstream}, {"--output", o.Output}} {
		if p.value == "" {
			return Command{}, model.Errorf(model.CodeUsage, nil, "Pflichtoption %s fehlt", p.name)
		}
	}
	o.LogLevel = *level.wert
	return Command{Name: "record", Record: o}, nil
}

func parseReplay(args []string) (Command, error) {
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
	level, err := logLevelOption(fs)
	if err != nil {
		return Command{}, err
	}
	optionen, rest := endeDerOptionen(args)
	if err := fs.Parse(optionen); err != nil {
		return Command{}, model.Errorf(model.CodeUsage, err, "ungültige Verwendung von replay")
	}
	if rest = append(fs.Args(), rest...); len(rest) > 0 {
		return Command{}, model.Errorf(model.CodeUsage, nil, "unerwartetes Argument %q", rest[0])
	}
	for _, p := range []struct{ name, value string }{{"--listen", o.Listen}, {"--input", o.Input}} {
		if p.value == "" {
			return Command{}, model.Errorf(model.CodeUsage, nil, "Pflichtoption %s fehlt", p.name)
		}
	}
	o.LogLevel = *level.wert
	return Command{Name: "replay", Replay: o}, nil
}

// logLevelOption meldet --log-level an fs an: Standard info, davor die
// Umgebungsvariable, wenn sie nicht leer ist, danach die Kommandozeile, deren
// letzte Angabe gilt. Eine Umgebungsvariable mit ungültigem Wert ist
// PGR-E2001, auch wenn die Kommandozeile die Option setzt (LH-FA-17.a).
func logLevelOption(fs *flag.FlagSet) (stufe, error) {
	l := stufe{new(string)}
	*l.wert = LogInfo
	if v := os.Getenv(envLogLevel); v != "" {
		if err := l.Set(v); err != nil {
			return stufe{}, model.Errorf(model.CodeUsage, err, "Umgebungsvariable %s", envLogLevel)
		}
	}
	fs.Var(l, "log-level", "")
	return l, nil
}

// stufe ist der Wert von --log-level: genau error, warn, info oder debug in
// Kleinbuchstaben; jeder andere Wert, auch der leere, ist ein Fehler
// (LH-FA-14.a).
type stufe struct{ wert *string }

func (l stufe) String() string {
	if l.wert == nil {
		return ""
	}
	return *l.wert
}

func (l stufe) Set(v string) error {
	switch v {
	case LogError, LogWarn, LogInfo, LogDebug:
		*l.wert = v
		return nil
	}
	return errors.New("erlaubt sind error, warn, info und debug")
}

// wahrheitswert ist der Wert einer booleschen Option nach LH-FA-17.a: ohne
// Wert true, mit Wert genau "true" oder "false"; jeder andere Wert, auch der
// leere und "1", ist ein Fehler. Nennt die Kommandozeile die Option mehrfach,
// gilt die letzte Angabe.
type wahrheitswert struct{ wert *bool }

// IsBoolFlag lässt die Option ohne Wert zu; flag setzt dann "true".
func (wahrheitswert) IsBoolFlag() bool { return true }

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

// hilfeAngabe meldet, ob ein ganzes Argument eine Hilfe-Angabe ist: genau -h,
// --h, -help oder --help, ohne oder mit "=" und beliebigem Wert (LH-FA-01.a).
func hilfeAngabe(arg string) bool {
	name, _, _ := strings.Cut(arg, "=")
	switch name {
	case "-h", "--h", "-help", "--help":
		return true
	}
	return false
}

// hilfeVerlangt meldet, ob unter den Argumenten vor dem ersten "--" eine
// Hilfe-Angabe steht; danach ist sie ein gewöhnliches Argument (LH-FA-01.a).
func hilfeVerlangt(args []string) bool {
	for _, a := range args {
		if a == "--" {
			return false
		}
		if hilfeAngabe(a) {
			return true
		}
	}
	return false
}

// endeDerOptionen teilt die Argumente am ersten "--", an jeder Stelle, auch an
// der eines Optionswerts: optionen sind die Argumente davor, rest die danach,
// gewöhnliche Argumente; "--" selbst gehört zu keinem (LH-FA-01.a).
func endeDerOptionen(args []string) (optionen, rest []string) {
	for i, a := range args {
		if a == "--" {
			return args[:i], args[i+1:]
		}
	}
	return args, nil
}
