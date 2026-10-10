package cli

import "go.yaml.in/yaml/v3"

// Die Umgebungsvariablen von --fail-on-unconsumed, --log-level und
// --shutdown-timeout.
const (
	EnvFailOnUnconsumed = envFailOnUnconsumed
	EnvLogLevel         = envLogLevel
	EnvShutdownTimeout  = envShutdownTimeout
	EnvPassword         = envPassword
)

// Option beschreibt eine Option am allgemeinen Leser für Tests: Name,
// Umgebungsvariable, Name der Wertemenge, Pflicht und Standardwert.
type Option struct {
	Name, Env, Art string
	Pflicht        bool
	Standard       string
}

// Optionen liefert die Optionen eines Kommandos am allgemeinen Leser in ihrer
// Reihenfolge.
func Optionen(kommando string) []Option {
	var out []Option
	for _, o := range optionen(kommando) {
		out = append(out, Option{Name: o.name, Env: envName(o.name), Art: o.art.name, Pflicht: o.pflicht, Standard: o.standard})
	}
	return out
}

// EnvConfig ist die Umgebungsvariable von --config, StandardDatei die Datei im
// aktuellen Verzeichnis, die ohne --config und PGWIRE_RECORDER_CONFIG gilt.
const (
	EnvConfig     = envConfig
	StandardDatei = standardDatei
)

// FormMitTaggedStyle ruft form mit einem Knoten, den die Bibliothek als
// TaggedStyle markiert, an einer Zeile außerhalb des Texts der Datei.
func FormMitTaggedStyle() error {
	return form(&yaml.Node{Kind: yaml.ScalarNode, Style: yaml.TaggedStyle, Line: 5, Column: 1}, "x", nil)
}

// Zerlegt ist eine zerlegte URL für Tests: je Teil die Stücke in der
// Reihenfolge der URL, ein Platzhalter als <VAR>; ein Teil, den die URL nicht
// schreibt, ist nil.
type Zerlegt struct {
	Benutzer, Passwort, Host, Port, Datenbank []string
	SSLMode                                   string
}

// stuecke ist ein Teil als Stücke für Zerlegt.
func stuecke(t teil) []string {
	out := []string{}
	for _, s := range t {
		if s.variable != "" {
			out = append(out, "<"+s.variable+">")
		} else {
			out = append(out, s.text)
		}
	}
	return out
}

// Zerlege lädt die Datei pfad und liefert ihre erste Verbindung zerlegt.
func Zerlege(pfad string) (Zerlegt, error) {
	d, err := ladeDatei(pfad, "--config")
	if err != nil {
		return Zerlegt{}, err
	}
	v := d.verbindungen[0]
	z := Zerlegt{Host: stuecke(v.host), Port: stuecke(v.port), Datenbank: stuecke(v.datenbank), SSLMode: v.sslmode}
	if v.mitBenutzer {
		z.Benutzer = stuecke(v.benutzer)
	}
	if v.mitPasswort {
		z.Passwort = stuecke(v.passwort)
	}
	return z, nil
}

// SetzeEin lädt die Datei pfad und setzt in alle Teile ihrer ersten Verbindung
// in der Reihenfolge der URL (Benutzer, Passwort, Host, Port, Datenbank) die
// Variablen aus wert ein; Benutzer, Passwort und Datenbank, die die URL nicht
// schreibt, sind "", ein nicht geschriebener Port ist 5432 (das Laden legt ihn
// so ab).
func SetzeEin(pfad string, wert func(string) string) ([]string, error) {
	d, err := ladeDatei(pfad, "--config")
	if err != nil {
		return nil, err
	}
	v := d.verbindungen[0]
	return v.einsetzen(wert, v.benutzer, v.passwort, v.host, v.port, v.datenbank)
}
