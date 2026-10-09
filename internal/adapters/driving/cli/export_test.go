package cli

// Die Umgebungsvariablen von --fail-on-unconsumed, --log-level und
// --shutdown-timeout.
const (
	EnvFailOnUnconsumed = envFailOnUnconsumed
	EnvLogLevel         = envLogLevel
	EnvShutdownTimeout  = envShutdownTimeout
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
