package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// RecordOptions sind die Optionen von `record` (LH-FA-02.a), soweit dieser Stand
// sie kennt.
type RecordOptions struct {
	Listen          string
	Upstream        string
	Output          string
	Force           bool
	LogLevel        string
	ShutdownTimeout time.Duration
}

// ReplayOptions sind die Optionen von `replay` (LH-FA-03.a), soweit dieser Stand
// sie kennt.
type ReplayOptions struct {
	Listen           string
	Input            string
	FailOnUnconsumed bool
	LogLevel         string
	ShutdownTimeout  time.Duration
}

// envFailOnUnconsumed ist die Umgebungsvariable von --fail-on-unconsumed
// (LH-FA-17.a).
const envFailOnUnconsumed = "PGWIRE_RECORDER_FAIL_ON_UNCONSUMED"

// envLogLevel ist die Umgebungsvariable von --log-level (LH-FA-14.a).
const envLogLevel = "PGWIRE_RECORDER_LOG_LEVEL"

// envShutdownTimeout ist die Umgebungsvariable von --shutdown-timeout
// (LH-FA-17.a).
const envShutdownTimeout = "PGWIRE_RECORDER_SHUTDOWN_TIMEOUT"

// StandardFrist ist der Standardwert von --shutdown-timeout (SPEC-046).
const StandardFrist = 5 * time.Second

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
` + optionShutdownTimeout + optionLogLevel + `
Jede Option außer --output und --force ist auch über ihre Umgebungsvariable
setzbar: PGWIRE_RECORDER_ und der Name in Großbuchstaben mit _ statt -, etwa
PGWIRE_RECORDER_LISTEN; die Kommandozeile geht ihr vor.
`

const optionenReplay = `Optionen von replay:
  --listen    Adresse, auf der Clients angenommen werden (Pflicht)
  --input     Aufzeichnung (Pflicht)
  --fail-on-unconsumed[=true|false]
              nicht verbrauchte Interaktionen und nie zugeordnete Sessions
              sind ein Fehler (PGR-E5002, Exit-Code 5) statt einer Warnung;
              Umgebungsvariable PGWIRE_RECORDER_FAIL_ON_UNCONSUMED
` + optionShutdownTimeout + optionLogLevel + `
Jede Option ist auch über ihre Umgebungsvariable setzbar: PGWIRE_RECORDER_ und
der Name in Großbuchstaben mit _ statt -, etwa PGWIRE_RECORDER_LISTEN; die
Kommandozeile geht ihr vor.
`

const optionShutdownTimeout = `  --shutdown-timeout 0|<zahl>ms|<zahl>s|<zahl>m
              Frist ab dem ersten SIGINT oder SIGTERM, Standard 5s; danach
              endet jede noch laufende Verbindung zwangsweise, und eine
              unvollständige Interaktion ist PGR-E4006 (Exit-Code 4); 0 ohne
              Frist; ein zweites Signal lässt die Frist sofort ablaufen;
              Umgebungsvariable PGWIRE_RECORDER_SHUTDOWN_TIMEOUT
`

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

// option ist eine Option von record oder replay am allgemeinen Leser: Er liest
// sie von der Kommandozeile und aus ihrer Umgebungsvariable (envName), prüft
// jeden gesetzten Wert mit art und übernimmt den Wert nach der Priorität
// Kommandozeile vor Umgebungsvariable vor Standardwert (LH-FA-17.a, SPEC-007).
type option struct {
	name     string
	art      art
	pflicht  bool
	standard string
	setze    func(*Command, string)
}

// art ist die Wertemenge einer Option: pruefe lehnt jeden Wert außerhalb ab;
// schalter lässt die Option auf der Kommandozeile ohne Wert zu, dann gilt
// "true".
type art struct {
	name     string
	pruefe   func(string) error
	schalter bool
}

func artText() art { return art{name: "text", pruefe: func(string) error { return nil }} }

func artWahrheitswert() art {
	return art{name: "wahrheitswert", pruefe: func(v string) error { return wahrheitswert{new(bool)}.Set(v) }, schalter: true}
}

func artDauer() art {
	return art{name: "dauer", pruefe: func(v string) error { return dauer{new(time.Duration)}.Set(v) }}
}

func artStufe() art {
	return art{name: "stufe", pruefe: func(v string) error { return stufe{new(string)}.Set(v) }}
}

// envName ist der Name der Umgebungsvariable einer Option: das Präfix
// PGWIRE_RECORDER_ (SPEC-008), dann der Name in Großbuchstaben mit _ statt -
// (LH-FA-17.a).
func envName(name string) string {
	return "PGWIRE_RECORDER_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
}

// optionen liefert die Optionen eines Kommandos am allgemeinen Leser in der
// Reihenfolge der Tabelle in LH-FA-17.a; in dieser Reihenfolge prüft der
// Leser die Umgebungsvariablen. --output und --force von record stehen nicht
// darunter, sie liest nur die Kommandozeile (parseRecord).
func optionen(kommando string) []option {
	switch kommando {
	case "record":
		return []option{
			{name: "listen", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Record.Listen = v }},
			{name: "upstream", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Record.Upstream = v }},
			{name: "shutdown-timeout", art: artDauer(), standard: StandardFrist.String(), setze: func(c *Command, v string) { setzeDauer(&c.Record.ShutdownTimeout, v) }},
			{name: "log-level", art: artStufe(), standard: LogInfo, setze: func(c *Command, v string) { c.Record.LogLevel = v }},
		}
	case "replay":
		return []option{
			{name: "listen", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Replay.Listen = v }},
			{name: "input", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Replay.Input = v }},
			{name: "fail-on-unconsumed", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Replay.FailOnUnconsumed = v == "true" }},
			{name: "shutdown-timeout", art: artDauer(), standard: StandardFrist.String(), setze: func(c *Command, v string) { setzeDauer(&c.Replay.ShutdownTimeout, v) }},
			{name: "log-level", art: artStufe(), standard: LogInfo, setze: func(c *Command, v string) { c.Replay.LogLevel = v }},
		}
	}
	return nil
}

// setzeDauer setzt einen Wert von --shutdown-timeout, den artDauer geprüft hat;
// Set scheitert an ihm nicht.
func setzeDauer(ziel *time.Duration, v string) {
	_ = dauer{ziel}.Set(v)
}

// kommandozeile ist der Wert einer Option auf der Kommandozeile: Set prüft ihn
// mit der Wertemenge und merkt die letzte Angabe.
type kommandozeile struct {
	art     art
	wert    *string
	gesetzt *bool
}

func (k kommandozeile) String() string {
	if k.wert == nil {
		return ""
	}
	return *k.wert
}

func (k kommandozeile) Set(v string) error {
	if err := k.art.pruefe(v); err != nil {
		return err
	}
	*k.wert, *k.gesetzt = v, true
	return nil
}

// IsBoolFlag lässt eine Option der Art schalter ohne Wert zu; flag setzt dann
// "true".
func (k kommandozeile) IsBoolFlag() bool { return k.art.schalter }

// gelesen ist der Stand einer Option nach dem Lesen der Quellen.
type gelesen struct {
	cli, env     string
	cliOk, envOk bool
}

// lies liest kommando nach LH-FA-17.a: zuerst die Kommandozeile mit fs, an dem
// der Aufrufer weitere Optionen angemeldet haben kann, dann die
// Umgebungsvariablen der Optionen in ihrer Reihenfolge; jeder gesetzte Wert
// wird geprüft, ein ungültiger ist PGR-E2001, auch wenn die Kommandozeile
// dieselbe Option setzt. Danach übernimmt es je Option den Wert nach der
// Priorität; eine Pflichtoption ohne Wert oder mit leerem Wert ist PGR-E2001.
func lies(kommando string, fs *flag.FlagSet, args []string) (Command, error) {
	opts := optionen(kommando)
	stand := make([]gelesen, len(opts))
	for i, o := range opts {
		fs.Var(kommandozeile{o.art, &stand[i].cli, &stand[i].cliOk}, o.name, "")
	}
	vorne, rest := endeDerOptionen(args)
	if err := fs.Parse(vorne); err != nil {
		return Command{}, model.Errorf(model.CodeUsage, err, "ungültige Verwendung von %s", kommando)
	}
	if rest = append(fs.Args(), rest...); len(rest) > 0 {
		return Command{}, model.Errorf(model.CodeUsage, nil, "unerwartetes Argument %q", rest[0])
	}
	for i, o := range opts {
		v := os.Getenv(envName(o.name))
		if v == "" {
			continue
		}
		if err := o.art.pruefe(v); err != nil {
			return Command{}, model.Errorf(model.CodeUsage, err, "Umgebungsvariable %s", envName(o.name))
		}
		stand[i].env, stand[i].envOk = v, true
	}
	cmd := Command{Name: kommando}
	for i, o := range opts {
		v := o.standard
		switch {
		case stand[i].cliOk:
			v = stand[i].cli
		case stand[i].envOk:
			v = stand[i].env
		}
		if o.pflicht && v == "" {
			return Command{}, model.Errorf(model.CodeUsage, nil, "Pflichtoption --%s fehlt", o.name)
		}
		o.setze(&cmd, v)
	}
	return cmd, nil
}

func parseRecord(args []string) (Command, error) {
	fs := flag.NewFlagSet("record", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var output string
	var force bool
	fs.StringVar(&output, "output", "", "")
	fs.BoolVar(&force, "force", false, "")
	cmd, err := lies("record", fs, args)
	if err != nil {
		return Command{}, err
	}
	if output == "" {
		return Command{}, model.Errorf(model.CodeUsage, nil, "Pflichtoption --output fehlt")
	}
	cmd.Record.Output, cmd.Record.Force = output, force
	return cmd, nil
}

func parseReplay(args []string) (Command, error) {
	fs := flag.NewFlagSet("replay", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return lies("replay", fs, args)
}

// dauerForm ist eine ganze Zahl ohne Vorzeichen mit genau einer Einheit.
var dauerForm = regexp.MustCompile(`^([0-9]+)(ms|s|m)$`)

// einheit ist die Dauer einer Einheit ms, s oder m (LH-FA-17.a *Dauer*);
// dauerForm lässt keine andere zu.
func einheit(name string) time.Duration {
	switch name {
	case "ms":
		return time.Millisecond
	case "s":
		return time.Second
	default:
		return time.Minute
	}
}

// dauer ist der Wert von --shutdown-timeout nach LH-FA-17.a *Dauer*: 0 oder
// eine ganze Zahl ohne Vorzeichen mit genau einer Einheit ms, s oder m in
// Kleinbuchstaben, auch mit führenden Nullen; 0 mit Einheit ist 0. Jeder andere Wert ist ein Fehler, auch
// der leere und einer, der länger ist als die längste Dauer, die time.Duration
// darstellt.
type dauer struct{ wert *time.Duration }

func (d dauer) String() string {
	if d.wert == nil {
		return ""
	}
	return d.wert.String()
}

func (d dauer) Set(v string) error {
	if v == "0" {
		*d.wert = 0
		return nil
	}
	m := dauerForm.FindStringSubmatch(v)
	if m == nil {
		return errors.New("erlaubt sind 0 und eine ganze Zahl mit genau einer Einheit ms, s oder m, etwa 5s")
	}
	n, err := strconv.ParseInt(m[1], 10, 64)
	e := einheit(m[2])
	if err != nil || n > math.MaxInt64/int64(e) {
		return errors.New("länger als die längste darstellbare Dauer")
	}
	*d.wert = time.Duration(n) * e
	return nil
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
