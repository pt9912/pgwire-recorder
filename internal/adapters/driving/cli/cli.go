package cli

import (
	"crypto/x509"
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
// sie kennt. RecordEmptySessions ist --record-empty-sessions (LH-FA-12.a). TLS
// sind Zertifikat und Schlüssel aus --tls-cert und --tls-key, gelesen beim
// Start, und --allow-plaintext (LH-FA-23.a).
type RecordOptions struct {
	Listen              string
	Upstream            string
	Output              string
	Force               bool
	Format              string
	RecordTiming        bool
	RecordEmptySessions bool
	LogLevel            string
	ShutdownTimeout     time.Duration
	TLS                 TLSZugang
}

// ReplayOptions sind die Optionen von `replay` (LH-FA-03.a), soweit dieser Stand
// sie kennt. SessionAssignment ist der Wert von --session-assignment,
// ZuordnungErsteAnfrage oder ZuordnungVerbindung (LH-FA-12.a). TLS wie bei
// RecordOptions.
type ReplayOptions struct {
	Listen            string
	Input             string
	FailOnUnconsumed  bool
	SessionAssignment string
	LogLevel          string
	ShutdownTimeout   time.Duration
	TLS               TLSZugang
}

// Die Werte von --session-assignment (LH-FA-12.a); Standard ist
// ZuordnungErsteAnfrage.
const (
	ZuordnungErsteAnfrage = "first-request"
	ZuordnungVerbindung   = "connection"
)

// PlayOptions sind die Optionen von `play` (LH-FA-20.a), soweit dieser Stand
// sie kennt. Nennt --upstream eine benannte Verbindung, ist Upstream deren
// Adresse host:port. Einspielen trägt die Optionen, die der Play-Service
// auswertet. Passwort ist das Passwort der Anmeldung am Server (LH-FA-20.a
// *Passwort*); "" heißt, es gibt keines. UpstreamTLS sagt, ob die Verbindung
// zum Server TLS nutzt: das gesetzte --upstream-tls, sonst sslmode=require der
// benutzten Verbindung (LH-FA-17.a *Wirkung einer URL*). UpstreamCA sind die
// Zertifikate der Datei aus --upstream-ca, gelesen beim Start; ohne die Option
// sind es keine.
type PlayOptions struct {
	Upstream    string
	Input       string
	Einspielen  Einspielvorgaben
	Passwort    Passwort
	UpstreamTLS bool
	UpstreamCA  Zertifikate
	LogLevel    string
}

// Zertifikate sind die Zertifikate einer Datei aus --upstream-ca. Jede
// Formatierung gibt festen Text aus, nie ein Zertifikat: Ein Aufruf mit %v,
// %+v, %#v, %s oder %q auf den Optionen von play verrät keines.
type Zertifikate []*x509.Certificate

// Format gibt für jedes Verb den festen Text aus.
func (z Zertifikate) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprintf(f, "Zertifikate{%d}", len(z))
}

// Passwort ist ein Passwort der Anmeldung. Jede Formatierung gibt festen Text
// aus, nie den Wert: Ein Aufruf mit %v, %+v, %#v, %s oder %q auf den Optionen
// von play verrät es nicht (LH-FA-20.a *Passwort*).
type Passwort string

// Format gibt für jedes Verb den festen Text aus.
func (Passwort) Format(f fmt.State, _ rune) {
	_, _ = io.WriteString(f, "***")
}

// Einspielvorgaben sind die Optionen von `play`, die der Play-Service
// auswertet. User und Database tragen ohne --user beziehungsweise --database
// Benutzer und Datenbank der benutzten Verbindung (LH-FA-17.a *Wirkung einer
// URL*); "" heißt, die Startup-Daten der Session gelten. ContinueOnError,
// AllowRecordedErrors und FinishSessionOnInterrupt sind --continue-on-error,
// --allow-recorded-errors und --finish-session-on-interrupt (LH-FA-20.a).
// CompareResponses ist --compare-responses (LH-FA-24.a).
// KeepTiming ist --keep-timing; mit ihm sind TimingMode (relative oder absolute)
// und, bei absolute, TimingReference (connect oder first-request) der Wert der
// zusammengeführten Optionen (LH-FA-21.a *Einspielen*), sonst sind beide leer,
// ebenso TimingReference bei relative.
//
// Kopplung: Der Play-Service führt einen Typ mit denselben Feldern in
// derselben Reihenfolge, und der Bootstrap konvertiert diesen in jenen; eine
// Option kommt in beiden Typen hinzu, sonst baut der Bootstrap nicht.
type Einspielvorgaben struct {
	User                     string
	Database                 string
	ContinueOnError          bool
	AllowRecordedErrors      bool
	FinishSessionOnInterrupt bool
	KeepTiming               bool
	TimingMode               string
	TimingReference          string
	CompareResponses         bool
}

// envFailOnUnconsumed ist die Umgebungsvariable von --fail-on-unconsumed
// (LH-FA-17.a).
const envFailOnUnconsumed = "PGWIRE_RECORDER_FAIL_ON_UNCONSUMED"

// envPassword ist die Umgebungsvariable des Passworts bei play (LH-FA-17.a,
// LH-FA-20.a *Passwort*); eine Option dafür gibt es nicht.
const envPassword = "PGWIRE_RECORDER_PASSWORD"

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

// Command ist das gewählte Kommando mit seinen Optionen; bei config show
// trägt Anzeige die Ausgabe für stdout (LH-FA-17.a *Anzeige*).
type Command struct {
	Name    string
	Record  RecordOptions
	Replay  ReplayOptions
	Play    PlayOptions
	Anzeige string
}

const optionenRecord = `Optionen von record:
  --listen    Adresse, auf der Clients angenommen werden (Pflicht)
  --upstream  PostgreSQL-Server als host:port oder Name einer Verbindung der
              Konfigurationsdatei (Pflicht)
  --output    Zieldatei der Aufzeichnung (Pflicht)
  --force     vorhandene Zieldatei ersetzen
  --format yaml|sqlite
              Format der Aufzeichnung, Standard yaml; Umgebungsvariable
              PGWIRE_RECORDER_FORMAT
  --record-timing[=true|false]
              trägt an jeder Interaktion offset_ms ein, die Millisekunden
              vom Beginn der Session bis zum Eintreffen der ersten
              Client-Nachricht der Interaktion; Standard false;
              Umgebungsvariable PGWIRE_RECORDER_RECORD_TIMING
  --record-empty-sessions[=true|false]
              zeichnet auch eine Verbindung ohne Anfrage als Session ohne
              Interaktionen auf; Standard false; Umgebungsvariable
              PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS
` + optionShutdownTimeout + optionTLS + optionLogLevel + optionConfig + optionUmgebung

const optionenReplay = `Optionen von replay:
  --listen    Adresse, auf der Clients angenommen werden (Pflicht)
  --input     Aufzeichnung (Pflicht)
  --fail-on-unconsumed[=true|false]
              nicht verbrauchte Interaktionen und nie zugeordnete Sessions
              sind ein Fehler (PGR-E5002, Exit-Code 5) statt einer Warnung;
              Umgebungsvariable PGWIRE_RECORDER_FAIL_ON_UNCONSUMED
  --session-assignment first-request|connection
              first-request ordnet die Sessions den Verbindungen in der
              Reihenfolge der ersten Anfragen zu, connection der n-ten
              Verbindung die Session mit der Kennung n (verlangt die
              Kennzeichnung empty_sessions, die record
              --record-empty-sessions schreibt, sonst PGR-E2001);
              Standard first-request; Umgebungsvariable
              PGWIRE_RECORDER_SESSION_ASSIGNMENT
` + optionShutdownTimeout + optionTLS + optionLogLevel + optionConfig + optionUmgebung

const optionenPlay = `Optionen von play:
  --upstream  PostgreSQL-Server als host:port oder Name einer Verbindung der
              Konfigurationsdatei (Pflicht); verlangt der Server ein
              Passwort, nimmt play es aus dem Platzhalter der Verbindung,
              sonst aus der Umgebungsvariable PGWIRE_RECORDER_PASSWORD, und
              meldet sich mit Klartext-Passwort, md5 oder scram-sha-256 an;
              ein falsches oder fehlendes Passwort endet play mit
              PGR-E4005
  --input     Aufzeichnung (Pflicht)
  --user      Benutzer für jede Session, statt dem der Verbindung und der
              Aufzeichnung
  --database  Datenbank für jede Session, statt der der Verbindung und der
              Aufzeichnung
  --continue-on-error[=true|false]
              nach einer Fehlerantwort des Servers (PGR-E4004) macht play mit
              der nächsten Anfrage weiter; Umgebungsvariable
              PGWIRE_RECORDER_CONTINUE_ON_ERROR
  --allow-recorded-errors[=true|false]
              eine Fehlerantwort mit dem Schweregrad ERROR ist kein Fehler,
              wenn die aufgezeichnete Anfrage ebenfalls eine trägt; ebenso ein
              Ende des Servers (FATAL oder PANIC), wenn die aufgezeichnete
              Anfrage es ebenfalls trägt; dann läuft die nächste Session; mit
              --compare-responses ohne Wirkung; Umgebungsvariable
              PGWIRE_RECORDER_ALLOW_RECORDED_ERRORS
  --upstream-tls[=true|false]
              die Verbindung zum Server nutzt TLS und prüft sein Zertifikat
              und seinen Namen; ohne die Option entscheidet sslmode der
              Verbindung, und die Option geht vor, auch mit false;
              Umgebungsvariable PGWIRE_RECORDER_UPSTREAM_TLS
  --finish-session-on-interrupt[=true|false]
              nach dem ersten SIGINT oder SIGTERM endet play erst nach der
              laufenden Session; Umgebungsvariable
              PGWIRE_RECORDER_FINISH_SESSION_ON_INTERRUPT
  --keep-timing[=true|false]
              wartet vor jeder Interaktion auf ihre aufgezeichnete Zeit
              (offset_ms); fehlt sie einer Interaktion, endet play mit
              PGR-E2003; Standard false; Umgebungsvariable
              PGWIRE_RECORDER_KEEP_TIMING
  --timing-mode relative|absolute
              relative wartet seit dem Beginn der vorigen Interaktion,
              absolute seit dem Bezugspunkt; Standard relative; nur mit
              --keep-timing (sonst PGR-E2001); Umgebungsvariable
              PGWIRE_RECORDER_TIMING_MODE
  --timing-reference connect|first-request
              Bezugspunkt von absolute, der Beginn des Verbindungsaufbaus
              oder der Beginn der ersten Interaktion der Session; Standard
              connect; nur mit --keep-timing und absolute (sonst PGR-E2001);
              Umgebungsvariable PGWIRE_RECORDER_TIMING_REFERENCE
  --upstream-ca <datei>
              PEM-Datei mit Zertifikaten einer Zertifizierungsstelle, die
              zusätzlich zum Zertifikatsspeicher des Systems gilt; nur mit
              TLS (sonst PGR-E2001); eine Datei, die keine Zertifikate
              enthält, endet play mit PGR-E2007; Umgebungsvariable
              PGWIRE_RECORDER_UPSTREAM_CA
  --compare-responses[=true|false]
              vergleicht die Struktur der Serverantworten jeder Interaktion
              mit der Aufzeichnung; eine Abweichung ist PGR-E5004
              (Exit-Code 5); ein Ende des Servers (FATAL oder PANIC) ist
              erwartet, wenn die Aufzeichnung es an derselben Stelle mit
              demselben SQLSTATE zeigt, sonst eine Abweichung; Standard false;
              Umgebungsvariable PGWIRE_RECORDER_COMPARE_RESPONSES
` + optionLogLevel + optionConfig + optionUmgebung

const optionUmgebung = `
Jede Option ist auch über ihre Umgebungsvariable setzbar: PGWIRE_RECORDER_ und
der Name in Großbuchstaben mit _ statt -, etwa PGWIRE_RECORDER_LISTEN; die
Kommandozeile geht ihr vor. Ein leerer Wert auf der Kommandozeile ist ungültig.
In der Konfigurationsdatei sind alle Optionen ohne --config setzbar: Schlüssel
wie die Option mit _ statt -, log_level auf der obersten Ebene, jeder andere im
Abschnitt des Kommandos (record:, replay:, play:); die Umgebungsvariable geht der
Datei vor, die Datei dem Standardwert.
`

const optionTLS = `  --tls-cert <datei>
              PEM-Datei mit dem Zertifikat des Servers und seiner Kette; mit
              --tls-cert und --tls-key nimmt der Recorder TLS von Clients an;
              jede Option allein ist PGR-E2001; eine nicht lesbare, ungültige
              oder abgelaufene Datei endet den Start mit PGR-E2007;
              Umgebungsvariable PGWIRE_RECORDER_TLS_CERT
  --tls-key <datei>
              PEM-Datei mit dem privaten Schlüssel zum Zertifikat, ohne
              Passwort; eine verschlüsselte oder nicht zum Zertifikat
              passende Datei endet den Start mit PGR-E2007; Umgebungsvariable
              PGWIRE_RECORDER_TLS_KEY
  --allow-plaintext[=true|false]
              mit TLS auch unverschlüsselte Clients zulassen, sonst weist der
              Recorder sie mit PGR-E6003 ab; true nur mit --tls-cert (sonst
              PGR-E2001); Standard false; Umgebungsvariable
              PGWIRE_RECORDER_ALLOW_PLAINTEXT
`

const optionConfig = `  --config <datei>
              Konfigurationsdatei, Standard .pgwire-recorder.yaml im aktuellen
              Verzeichnis, sofern vorhanden; Umgebungsvariable
              PGWIRE_RECORDER_CONFIG
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
  play     spielt die Anfragen einer Aufzeichnung in PostgreSQL ein
  config show
           zeigt die gewählte Konfigurationsdatei
  version  gibt die Programmversion aus

Programmversion, nur als einziges Argument:
  --version, -v
           gibt die Programmversion aus, wie version

` + optionenRecord + `
` + optionenReplay + `
` + optionenPlay + `
` + optionenConfigShow

const optionenConfigShow = `Optionen von config show:
` + optionConfig + `
Gibt auf stdout den Pfad der gewählten Konfigurationsdatei aus, danach ihren
Inhalt ohne Kommentare, Platzhalter unaufgelöst, danach die Namen der gesetzten
Umgebungsvariablen PGWIRE_RECORDER_*, ohne Werte. Ist die Datei ungültig, zeigt
der Befehl nichts.
`

// hilfen ist die Hilfe je bekanntem Kommando (LH-FA-01.a).
var hilfen = map[string]string{
	"record":  "Aufruf: pgwire-recorder record [optionen]\n\n" + optionenRecord,
	"replay":  "Aufruf: pgwire-recorder replay [optionen]\n\n" + optionenReplay,
	"play":    "Aufruf: pgwire-recorder play [optionen]\n\n" + optionenPlay,
	"config":  "Aufruf: pgwire-recorder config show [optionen]\n\n" + optionenConfigShow,
	"version": "Aufruf: pgwire-recorder version\n\nGibt die Programmversion aus.\n",
}

// ErrHelp meldet, dass die Hilfe angefordert und ausgegeben wurde.
var ErrHelp = errors.New("Hilfe ausgegeben")

// Parse liest Kommando und Optionen. Eine ungültige Verwendung ist PGR-E2001.
// Steht vor dem ersten "--" eine Hilfe-Angabe (hilfeAngabe), an welcher Stelle
// auch immer, gibt Parse vor jeder anderen Prüfung die Hilfe auf out aus und
// liefert ErrHelp: die des Kommandos, wenn das erste Argument ein bekanntes
// Kommando ist, sonst die globale (LH-FA-01.a). Das Kommando help gibt die
// globale Hilfe aus. Eine Versions-Angabe (versionsAngabe) liefert dasselbe
// Kommando wie version.
func Parse(args []string, out io.Writer) (Command, error) {
	return parse(args, out, time.Now)
}

// parse ist Parse mit der Systemuhr jetzt.
func parse(args []string, out io.Writer, jetzt func() time.Time) (Command, error) {
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
	if versionsAngabe(args) {
		return Command{Name: "version"}, nil
	}
	if len(args) == 0 {
		return Command{}, model.Errorf(model.CodeUsage, nil, "kein Kommando angegeben; --help zeigt die Kommandos")
	}
	switch args[0] {
	case "help":
		fmt.Fprint(out, usage)
		return Command{}, ErrHelp
	case "record":
		return lies("record", args[1:], jetzt)
	case "replay":
		return lies("replay", args[1:], jetzt)
	case "play":
		return lies("play", args[1:], jetzt)
	case "config":
		return parseConfig(args[1:])
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
// Der Schlüssel der Option in der Konfigurationsdatei steht im Abschnitt des
// Kommandos, mit oben auf der obersten Ebene (datei.wert). zuletzt prüft die
// Option nach der Zusammenführung aller Optionen, mit dem Stand ihrer Quellen
// und der gewählten Datei (LH-FA-17.a *Fehler*), und darf ihren Wert im
// Command ersetzen: upstreamRecord setzt c.Record.Upstream auf die Adresse der
// benutzten Verbindung. Ein späterer zuletzt-Aufruf sieht den ersetzten Wert;
// nil prüft nichts. Eine Option ohne setze trägt ihren Wert nur in quellen
// (upstream-ca, deren Datei zuletzt gelesen wird).
type option struct {
	name     string
	art      art
	pflicht  bool
	standard string
	oben     bool
	setze    func(*Command, string)
	zuletzt  func(*Command, gelesen, *datei, quellen) error
}

// art ist die Wertemenge einer Option: pruefe lehnt jeden Wert außerhalb ab,
// in jeder Wertemenge auch den leeren (LH-FA-17.a);
// schalter lässt die Option auf der Kommandozeile ohne Wert zu, dann gilt
// "true".
type art struct {
	name     string
	pruefe   func(string) error
	schalter bool
}

func artText() art {
	return art{name: "text", pruefe: func(v string) error {
		if v == "" {
			return errors.New("ein leerer Wert ist ungültig")
		}
		return nil
	}}
}

func artWahrheitswert() art {
	return art{name: "wahrheitswert", pruefe: func(v string) error { return wahrheitswert{new(bool)}.Set(v) }, schalter: true}
}

func artDauer() art {
	return art{name: "dauer", pruefe: func(v string) error { return dauer{new(time.Duration)}.Set(v) }}
}

func artFormat() art {
	return art{name: "format", pruefe: func(v string) error { return format{new(string)}.Set(v) }}
}

func artSitzungszuordnung() art {
	return art{name: "sitzungszuordnung", pruefe: func(v string) error { return sitzungszuordnung{new(string)}.Set(v) }}
}

func artZeitmodus() art {
	return art{name: "zeitmodus", pruefe: func(v string) error { return zeitmodus{new(string)}.Set(v) }}
}

func artZeitbezug() art {
	return art{name: "zeitbezug", pruefe: func(v string) error { return zeitbezug{new(string)}.Set(v) }}
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
// Leser die Umgebungsvariablen.
func optionen(kommando string) []option {
	switch kommando {
	case "record":
		return []option{
			{name: "listen", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Record.Listen = v }},
			{name: "upstream", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Record.Upstream = v }, zuletzt: upstreamRecord},
			{name: "output", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Record.Output = v }},
			{name: "force", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Record.Force = v == "true" }},
			{name: "format", art: artFormat(), standard: formatYAML, setze: func(c *Command, v string) { c.Record.Format = v }},
			{name: "record-timing", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Record.RecordTiming = v == "true" }},
			{name: "record-empty-sessions", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Record.RecordEmptySessions = v == "true" }},
			{name: "shutdown-timeout", art: artDauer(), standard: StandardFrist.String(), setze: func(c *Command, v string) { setzeDauer(&c.Record.ShutdownTimeout, v) }},

			{name: "tls-cert", art: artText(), zuletzt: tlsZertifikat},
			{name: "tls-key", art: artText(), zuletzt: tlsSchluessel},
			{name: "allow-plaintext", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Record.TLS.KlartextErlaubt = v == "true" }},
			{name: "log-level", art: artStufe(), standard: LogInfo, oben: true, setze: func(c *Command, v string) { c.Record.LogLevel = v }},
		}
	case "replay":
		return []option{
			{name: "listen", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Replay.Listen = v }},
			{name: "input", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Replay.Input = v }},
			{name: "fail-on-unconsumed", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Replay.FailOnUnconsumed = v == "true" }},
			{name: "shutdown-timeout", art: artDauer(), standard: StandardFrist.String(), setze: func(c *Command, v string) { setzeDauer(&c.Replay.ShutdownTimeout, v) }},
			{name: "session-assignment", art: artSitzungszuordnung(), standard: ZuordnungErsteAnfrage, setze: func(c *Command, v string) { c.Replay.SessionAssignment = v }},
			{name: "tls-cert", art: artText(), zuletzt: tlsZertifikat},
			{name: "tls-key", art: artText(), zuletzt: tlsSchluessel},
			{name: "allow-plaintext", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Replay.TLS.KlartextErlaubt = v == "true" }},
			{name: "log-level", art: artStufe(), standard: LogInfo, oben: true, setze: func(c *Command, v string) { c.Replay.LogLevel = v }},
		}
	case "play":
		return []option{
			{name: "upstream", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Play.Upstream = v }, zuletzt: upstreamPlay},
			{name: "input", art: artText(), pflicht: true, setze: func(c *Command, v string) { c.Play.Input = v }},
			{name: "user", art: artText(), setze: func(c *Command, v string) { c.Play.Einspielen.User = v }},
			{name: "database", art: artText(), setze: func(c *Command, v string) { c.Play.Einspielen.Database = v }},
			{name: "continue-on-error", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Play.Einspielen.ContinueOnError = v == "true" }},
			{name: "allow-recorded-errors", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Play.Einspielen.AllowRecordedErrors = v == "true" }},
			{name: "upstream-tls", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Play.UpstreamTLS = v == "true" }},
			{name: "finish-session-on-interrupt", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Play.Einspielen.FinishSessionOnInterrupt = v == "true" }},
			{name: "keep-timing", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Play.Einspielen.KeepTiming = v == "true" }},
			{name: "timing-mode", art: artZeitmodus(), standard: zeitRelative},
			{name: "timing-reference", art: artZeitbezug(), standard: bezugConnect, zuletzt: zeitPlay},
			{name: "upstream-ca", art: artText(), zuletzt: upstreamCA},
			{name: "compare-responses", art: artWahrheitswert(), standard: "false", setze: func(c *Command, v string) { c.Play.Einspielen.CompareResponses = v == "true" }},
			{name: "log-level", art: artStufe(), standard: LogInfo, oben: true, setze: func(c *Command, v string) { c.Play.LogLevel = v }},
		}
	}
	return nil
}

// leserKommandos sind die Kommandos am allgemeinen Leser; je eines hat einen
// Abschnitt in der Konfigurationsdatei.
func leserKommandos() []string { return []string{"record", "replay", "play"} }

func istLeserKommando(name string) bool {
	for _, k := range leserKommandos() {
		if k == name {
			return true
		}
	}
	return false
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

// gelesen ist der Stand einer Option nach dem Lesen der Quellen: ob die
// Kommandozeile, die Umgebungsvariable und die Datei sie setzen, und die
// Werte von Kommandozeile und Umgebungsvariable.
type gelesen struct {
	cli, env              string
	cliOk, envOk, dateiOk bool
}

// gesetzt meldet, ob irgendeine Quelle die Option setzt, auch mit dem Wert des
// Standardwerts.
func (g gelesen) gesetzt() bool { return g.cliOk || g.envOk || g.dateiOk }

// quellen ist, was ein zuletzt-Aufruf von der Zusammenführung sieht: je Option,
// nach ihrem Namen, den Stand ihrer Quellen und den zusammengeführten Wert.
type quellen struct {
	stand map[string]gelesen
	wert  map[string]string
	// jetzt liest die Systemuhr, gegen die der Start die Gültigkeit des
	// Serverzertifikats prüft (LH-FA-23.a *Dateien*).
	jetzt func() time.Time
}

// gesetzt meldet, ob irgendeine Quelle die Option name setzt.
func (q quellen) gesetzt(name string) bool { return q.stand[name].gesetzt() }

// lies liest kommando nach LH-FA-17.a: zuerst die Kommandozeile, an deren
// FlagSet genau die Optionen aus optionen und --config angemeldet sind, dann
// die Umgebungsvariablen in einem Durchlauf über die Zeilen von zuordnung (die
// Variable einer Option, die das Kommando nicht kennt, ist PGR-E2001), dann die
// gewählte Konfigurationsdatei; jeder gesetzte Wert wird geprüft, ein ungültiger
// der Kommandozeile oder einer Umgebungsvariable ist PGR-E2001, einer der Datei
// PGR-E2004, auch wenn eine Quelle davor dieselbe Option setzt. Danach
// übernimmt es je Option den Wert nach der Priorität (SPEC-007); eine
// Pflichtoption, die keine Quelle setzt, ist PGR-E2001. Zuletzt ruft es je
// Option in ihrer Reihenfolge option.zuletzt, das prüft und den Wert ersetzen
// kann (upstreamRecord setzt c.Record.Upstream auf die Adresse); davor prüft es
// die Kombinationen von --tls-cert, --tls-key und --allow-plaintext
// (tlsKombination), die vor jedem Fehler von --upstream stehen.
func lies(kommando string, args []string, jetzt func() time.Time) (Command, error) {
	opts := optionen(kommando)
	stand := make([]gelesen, len(opts))
	var config gelesen
	if err := liesKommandozeile(kommando, opts, stand, &config, args); err != nil {
		return Command{}, err
	}
	if err := liesUmgebung(kommando, opts, stand); err != nil {
		return Command{}, err
	}
	_, d, err := ladeGewaehlte(config.cli)
	if err != nil {
		return Command{}, err
	}
	cmd, q, err := zusammenfuehren(kommando, opts, stand, d)
	if err != nil {
		return Command{}, err
	}
	q.jetzt = jetzt
	if err := tlsKombination(kommando, q); err != nil {
		return Command{}, err
	}
	for i, o := range opts {
		if o.zuletzt == nil {
			continue
		}
		if err := o.zuletzt(&cmd, stand[i], d, q); err != nil {
			return Command{}, err
		}
	}
	return cmd, nil
}

// zeile ist eine Zeile der Optionstabelle in LH-FA-17.a: eine Option und die
// Kommandos der Spalte *Kommando*.
type zeile struct {
	name      string
	kommandos []string
}

// zuordnung ist die Spalte *Kommando* der Optionstabelle in LH-FA-17.a für jede
// Option der Tabelle, auch die, die der Stand noch nicht liest, in der
// Reihenfolge der Tabelle; --config und PGWIRE_RECORDER_PASSWORD sind keine
// Zeilen.
func zuordnung() []zeile {
	rr, r, rp, p := []string{"record", "replay"}, []string{"record"}, []string{"replay"}, []string{"play"}
	return []zeile{
		{"listen", rr}, {"upstream", []string{"record", "play"}}, {"output", r}, {"force", r}, {"format", r},
		{"record-timing", r}, {"record-empty-sessions", r}, {"input", []string{"replay", "play"}},
		{"fail-on-unconsumed", rp}, {"shutdown-timeout", rr}, {"session-assignment", rp}, {"user", p},
		{"database", p}, {"continue-on-error", p}, {"allow-recorded-errors", p}, {"upstream-tls", p},
		{"finish-session-on-interrupt", p}, {"keep-timing", p}, {"timing-mode", p}, {"timing-reference", p},
		{"tls-cert", rr}, {"tls-key", rr}, {"allow-plaintext", rr}, {"upstream-ca", p}, {"compare-responses", p},
		{"log-level", leserKommandos()},
	}
}

// kennt meldet, ob kommando in der Spalte *Kommando* der Zeile steht.
func (z zeile) kennt(kommando string) bool {
	for _, k := range z.kommandos {
		if k == kommando {
			return true
		}
	}
	return false
}

// liesUmgebung liest die Umgebungsvariablen in einem Durchlauf über die Zeilen
// der Zuordnung in der Reihenfolge der Tabelle: bei einer Option, die das
// Kommando liest, prüft es den gesetzten Wert; bei einer Option, die das
// Kommando kennt und der Stand noch nicht liest, bleibt die Variable
// unbeachtet; bei einer Option, die das Kommando nicht kennt, ist eine gesetzte,
// nicht leere Variable PGR-E2001 und nennt ihren Namen, nicht den Wert. Eine
// leere Variable gilt als nicht gesetzt (LH-FA-17.a).
func liesUmgebung(kommando string, opts []option, stand []gelesen) error {
	index := map[string]int{}
	for i, o := range opts {
		index[o.name] = i
	}
	for _, z := range zuordnung() {
		v := os.Getenv(envName(z.name))
		if v == "" {
			continue
		}
		if !z.kennt(kommando) {
			return model.Errorf(model.CodeUsage, nil, "Umgebungsvariable %s gehört zu einer Option, die %s nicht kennt", envName(z.name), kommando)
		}
		i, ok := index[z.name]
		if !ok {
			continue
		}
		if err := opts[i].art.pruefe(v); err != nil {
			return model.Errorf(model.CodeUsage, err, "Umgebungsvariable %s", envName(z.name))
		}
		stand[i].env, stand[i].envOk = v, true
	}
	return nil
}

// zusammenfuehren übernimmt je Option den Wert nach der Priorität (SPEC-007) in
// das Command, hält in stand fest, ob die Datei die Option setzt, und liefert
// den Stand der Quellen mit den zusammengeführten Werten; eine Pflichtoption,
// die keine Quelle setzt, ist PGR-E2001.
func zusammenfuehren(kommando string, opts []option, stand []gelesen, d *datei) (Command, quellen, error) {
	cmd := Command{Name: kommando}
	q := quellen{stand: map[string]gelesen{}, wert: map[string]string{}}
	for i, o := range opts {
		v := o.standard
		if w, ok := d.wert(kommando, o); ok {
			v, stand[i].dateiOk = w, true
		}
		switch {
		case stand[i].cliOk:
			v = stand[i].cli
		case stand[i].envOk:
			v = stand[i].env
		}
		if o.pflicht && v == "" {
			return Command{}, quellen{}, model.Errorf(model.CodeUsage, nil, "Pflichtoption --%s fehlt", o.name)
		}
		q.stand[o.name], q.wert[o.name] = stand[i], v
		if o.setze != nil {
			o.setze(&cmd, v)
		}
	}
	return cmd, q, nil
}

// liesKommandozeile liest die Argumente bis zum ersten "--" mit einem FlagSet,
// an dem die Optionen aus opts und --config angemeldet sind; jedes Argument
// danach oder ohne Option ist unerwartet (LH-FA-01.a).
func liesKommandozeile(kommando string, opts []option, stand []gelesen, config *gelesen, args []string) error {
	fs := flag.NewFlagSet(kommando, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	for i, o := range opts {
		fs.Var(kommandozeile{o.art, &stand[i].cli, &stand[i].cliOk}, o.name, "")
	}
	fs.Var(kommandozeile{artText(), &config.cli, &config.cliOk}, "config", "")
	vorne, rest := endeDerOptionen(args)
	if err := fs.Parse(vorne); err != nil {
		return model.Errorf(model.CodeUsage, err, "ungültige Verwendung von %s", kommando)
	}
	if rest = append(fs.Args(), rest...); len(rest) > 0 {
		return model.Errorf(model.CodeUsage, nil, "unerwartetes Argument %q", rest[0])
	}
	return nil
}

// parseConfig liest config show: auf der Kommandozeile nur --config, von den
// Umgebungsvariablen nur PGWIRE_RECORDER_CONFIG; es lädt und prüft die
// gewählte Datei wie jedes Kommando und liefert die Anzeige (LH-FA-17.a).
func parseConfig(args []string) (Command, error) {
	if len(args) == 0 || args[0] != "show" {
		return Command{}, model.Errorf(model.CodeUsage, nil, "config kennt nur das Kommando config show; --help zeigt die Kommandos")
	}
	var config gelesen
	if err := liesKommandozeile("config show", nil, nil, &config, args[1:]); err != nil {
		return Command{}, err
	}
	pfad, d, err := ladeGewaehlte(config.cli)
	if err != nil {
		return Command{}, err
	}
	text, err := anzeige(pfad, d)
	if err != nil {
		return Command{}, err
	}
	return Command{Name: "config show", Anzeige: text}, nil
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

func (l stufe) Set(v string) error {
	switch v {
	case LogError, LogWarn, LogInfo, LogDebug:
		*l.wert = v
		return nil
	}
	return errors.New("erlaubt sind error, warn, info und debug")
}

// Die Werte von --format (LH-FA-22.a *Wahl und Erkennung*).
const (
	formatYAML   = "yaml"
	formatSQLite = "sqlite"
)

// Formate liefert die Werte von --format in der Reihenfolge der Hilfe (LH-FA-22.a
// *Wahl und Erkennung*).
func Formate() []string { return []string{formatYAML, formatSQLite} }

// format ist der Wert von --format: genau yaml oder sqlite in Kleinbuchstaben;
// jeder andere Wert, auch der leere, ist ein Fehler (LH-FA-22.a).
type format struct{ wert *string }

func (f format) Set(v string) error {
	for _, erlaubt := range Formate() {
		if v == erlaubt {
			*f.wert = v
			return nil
		}
	}
	return errors.New("erlaubt sind yaml und sqlite")
}

// sitzungszuordnung ist der Wert von --session-assignment: genau first-request
// oder connection in Kleinbuchstaben, kein Alias; jeder andere Wert, auch der
// leere, ist ein Fehler (LH-FA-12.a *Wertemenge*).
type sitzungszuordnung struct{ wert *string }

func (z sitzungszuordnung) Set(v string) error {
	if v != ZuordnungErsteAnfrage && v != ZuordnungVerbindung {
		return errors.New("erlaubt sind first-request und connection")
	}
	*z.wert = v
	return nil
}

// wahrheitswert ist der Wert einer booleschen Option nach LH-FA-17.a: genau
// "true" oder "false"; jeder andere Wert, auch der leere und "1", ist ein
// Fehler. Dass die Option auf der Kommandozeile ohne Wert "true" ist, trägt
// art.schalter.
type wahrheitswert struct{ wert *bool }

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

// versionsAngabe meldet, ob args genau ein Argument sind und es genau --version
// oder -v lautet (LH-FA-01.a, Absatz Version). Jede andere Schreibweise, ein
// Argument mit "=" und ein weiteres Argument vor oder nach ihm ergeben keine
// Versions-Angabe.
func versionsAngabe(args []string) bool {
	return len(args) == 1 && (args[0] == "--version" || args[0] == "-v")
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
