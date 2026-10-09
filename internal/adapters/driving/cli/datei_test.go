package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// replayMit liest replay mit --listen, --input und den übrigen Argumenten.
func replayMit(args ...string) (cli.Command, error) {
	return lese(append([]string{"replay", "--listen=x", "--input=r.yaml"}, args...)...)
}

// Abdeckung: LH-FA-17/Happy, LH-FA-17/Boundary — es gilt genau eine Datei:
// die aus --config vor der aus PGWIRE_RECORDER_CONFIG vor
// .pgwire-recorder.yaml im aktuellen Verzeichnis; nennt --config eine Datei,
// wird die aus PGWIRE_RECORDER_CONFIG nicht gelesen, auch wenn sie fehlt; eine
// leere PGWIRE_RECORDER_CONFIG gilt als nicht gesetzt; ohne Datei wird keine
// gelesen (LH-FA-17.a).
func TestDateiWahl(t *testing.T) {
	leere(t, "replay")
	t.Chdir(t.TempDir())
	ausConfig := schreibe(t, "log_level: error\n")
	ausUmgebung := schreibe(t, "log_level: warn\n")
	stufe := func(args ...string) string {
		t.Helper()
		cmd, err := replayMit(args...)
		if err != nil {
			t.Fatalf("%q: %v", args, err)
		}
		return cmd.Replay.LogLevel
	}
	if got := stufe(); got != cli.LogInfo {
		t.Errorf("ohne Datei: %s", got)
	}
	if err := os.WriteFile(cli.StandardDatei, []byte("log_level: debug\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := stufe(); got != cli.LogDebug {
		t.Errorf("Standarddatei: %s", got)
	}
	t.Setenv(cli.EnvConfig, ausUmgebung)
	if got := stufe(); got != cli.LogWarn {
		t.Errorf("PGWIRE_RECORDER_CONFIG vor der Standarddatei: %s", got)
	}
	if got := stufe("--config", ausConfig); got != cli.LogError {
		t.Errorf("--config vor PGWIRE_RECORDER_CONFIG: %s", got)
	}
	t.Setenv(cli.EnvConfig, filepath.Join(t.TempDir(), "fehlt.yaml"))
	if got := stufe("--config", ausConfig); got != cli.LogError {
		t.Errorf("--config neben fehlender Datei aus PGWIRE_RECORDER_CONFIG: %s", got)
	}
	t.Setenv(cli.EnvConfig, "")
	if got := stufe(); got != cli.LogDebug {
		t.Errorf("leere PGWIRE_RECORDER_CONFIG: %s", got)
	}
}

// Abdeckung: LH-FA-17/Negative — eine mit --config oder PGWIRE_RECORDER_CONFIG
// genannte Datei, die fehlt oder nicht lesbar ist, und eine vorhandene, aber
// nicht lesbare Standarddatei (Verzeichnis, Schleife eines symbolischen Links)
// sind PGR-E2004; die Meldung nennt die Quelle,
// nicht den Pfad; ein leerer Wert von --config ist PGR-E2001 (LH-FA-17.a).
func TestDateiNichtLesbar(t *testing.T) {
	leere(t, "replay")
	t.Chdir(t.TempDir())
	fehlt := filepath.Join(t.TempDir(), "GEHEIM.yaml")
	verzeichnis := t.TempDir()
	for _, f := range []struct {
		env    string
		args   []string
		quelle string
	}{
		{"", []string{"--config", fehlt}, "--config"},
		{"", []string{"--config", verzeichnis}, "--config"},
		{fehlt, nil, cli.EnvConfig},
		{verzeichnis, nil, cli.EnvConfig},
	} {
		t.Setenv(cli.EnvConfig, f.env)
		_, err := replayMit(f.args...)
		if !istDatei(err) || !strings.Contains(err.Error(), f.quelle) || strings.Contains(err.Error(), "GEHEIM") || strings.Contains(err.Error(), verzeichnis) {
			t.Errorf("Umgebung %q, %q: erwartet %s mit %s ohne Pfad, erhalten %v", f.env, f.args, model.CodeConfigFile, f.quelle, err)
		}
	}
	t.Setenv(cli.EnvConfig, "")
	if err := os.Symlink(cli.StandardDatei, cli.StandardDatei); err != nil {
		t.Fatal(err)
	}
	if _, err := replayMit(); !istDatei(err) || !strings.Contains(err.Error(), cli.StandardDatei) {
		t.Errorf("Standarddatei als Schleife eines symbolischen Links: %v", err)
	}
	if err := os.Remove(cli.StandardDatei); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(cli.StandardDatei, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := replayMit(); !istDatei(err) || !strings.Contains(err.Error(), cli.StandardDatei) {
		t.Errorf("Standarddatei nicht lesbar: %v", err)
	}
	if _, err := replayMit("--config="); !istUsage(err) {
		t.Errorf("--config=: %v", err)
	}
}

// Abdeckung: LH-FA-17/Boundary — gültig sind eine leere Datei, eine nur mit
// Kommentaren, ein einzelnes --- vor dem Dokument, leere Abbildungen {} auf
// jeder Ebene und Werte, deren Text in der Wertemenge liegt, gleich ob mit
// oder ohne Anführungszeichen geschrieben: true und "true", 0 und "5s"; ein
// Schlüssel in Anführungszeichen zählt mit seinem Text; "null" und '~' in
// Anführungszeichen sind Text, bei output also ein Pfad (LH-FA-17.a).
func TestDateiGueltig(t *testing.T) {
	leere(t, "replay")
	standard := cli.ReplayOptions{Listen: "x", Input: "r.yaml", LogLevel: cli.LogInfo, ShutdownTimeout: cli.StandardFrist}
	mit := func(aendern func(*cli.ReplayOptions)) cli.ReplayOptions {
		o := standard
		aendern(&o)
		return o
	}
	for _, f := range []struct {
		inhalt string
		want   cli.ReplayOptions
	}{
		{"", standard},
		{"# nur ein Kommentar\n\n# noch einer\n", standard},
		{"{}\n", standard},
		{"replay: {}\nrecord: {}\nconnections: {}\n", standard},
		{"---\nreplay:\n  fail_on_unconsumed: true\n", mit(func(o *cli.ReplayOptions) { o.FailOnUnconsumed = true })},
		{"replay:\n  fail_on_unconsumed: \"true\"\n", mit(func(o *cli.ReplayOptions) { o.FailOnUnconsumed = true })},
		{"replay:\n  fail_on_unconsumed: 'true'\n", mit(func(o *cli.ReplayOptions) { o.FailOnUnconsumed = true })},
		{"replay:\n  shutdown_timeout: 0\n", mit(func(o *cli.ReplayOptions) { o.ShutdownTimeout = 0 })},
		{"replay:\n  shutdown_timeout: \"7s\"\n", mit(func(o *cli.ReplayOptions) { o.ShutdownTimeout = 7 * time.Second })},
		{"\"log_level\": debug\n", mit(func(o *cli.ReplayOptions) { o.LogLevel = cli.LogDebug })},
		{"replay:\n  input: ./relativ.yaml # Kommentar\n", mit(func(o *cli.ReplayOptions) { o.Input = "./relativ.yaml" })},
	} {
		args := []string{"replay", "--config=" + schreibe(t, f.inhalt), "--listen=x"}
		if f.want.Input == standard.Input {
			args = append(args, "--input=r.yaml")
		}
		cmd, err := lese(args...)
		if err != nil || cmd.Replay != f.want {
			t.Errorf("%q: %#v, %v, erwartet %#v", f.inhalt, cmd.Replay, err, f.want)
		}
	}
	leere(t, "record")
	for inhalt, pfad := range map[string]string{
		"record:\n  output: \"null\"\n": "null",
		"record:\n  output: '~'\n":      "~",
	} {
		cmd, err := lese("record", "--listen=x", "--upstream=h:1", "--config="+schreibe(t, inhalt))
		if err != nil || cmd.Record.Output != pfad {
			t.Errorf("%q: %#v, %v, erwartet den Pfad %q", inhalt, cmd.Record, err, pfad)
		}
	}
}

// ungueltigeDateien sind Dateien, die PGR-E2004 sind, mit der Stelle, die die
// Meldung nennt (LH-FA-17.a). GEHEIM steht für einen Wert, den keine Meldung
// nennt.
func ungueltigeDateien() []struct{ name, inhalt, stelle string } {
	return []struct{ name, inhalt, stelle string }{
		{"zwei Dokumente", "log_level: info\n---\nlog_level: warn\n", "mehr als ein Dokument"},
		{"zweites leeres Dokument", "log_level: info\n---\n", "mehr als ein Dokument"},
		{"Syntax", "replay:\n  listen: \"GEHEIM\n", "Konfigurationsdatei: ungültiges YAML in Zeile 2"},
		{"Parser, Folge", "log_level: info\nrecord:\n  force: true\nd: [x\n", "Konfigurationsdatei: ungültiges YAML in Zeile 4"},
		{"Parser, Abbildung unter record", "log_level: info\nrecord:\n  force: true\n  output: x\n  - y\n", "Konfigurationsdatei: ungültiges YAML in Zeile 3"},
		{"Parser, Grenze", "a: 1\nb: 2\n]\n", "Konfigurationsdatei: ungültiges YAML in Zeile 3"},
		{"Scanner", "a: 1\nb: c: d\n", "Konfigurationsdatei: ungültiges YAML in Zeile 2"},
		{"ohne Zahl", "]\n", "Konfigurationsdatei: ungültiges YAML in Zeile 1"},
		{"doppelt oben", "log_level: info\nlog_level: warn\n", "ungültiges YAML in Zeile 2, Schlüssel doppelt"},
		{"doppelt im Abschnitt", "replay:\n  listen: GEHEIM\n  listen: b\n", "ungültiges YAML in Zeile 3, Schlüssel doppelt"},
		{"doppelt in Anführungszeichen", "replay:\n  listen: a\n  \"listen\": b\n", "ungültiges YAML in Zeile 3, Schlüssel doppelt"},
		{"doppelte Verbindung", "connections:\n  a: x\n  a: y\n", "ungültiges YAML in Zeile 3, Schlüssel doppelt"},
		{"doppelt vor unbekanntem Schlüssel", "bogus: 1\nreplay:\n  listen: a\n  listen: b\n", "ungültiges YAML in Zeile 4, Schlüssel doppelt"},
		{"doppelt an Wertstelle mit Wert", "connections:\n  a: {postgresql://u:GEHEIM@h/db: 1, postgresql://u:GEHEIM@h/db: 2}\n", "ungültiges YAML in Zeile 2, Schlüssel doppelt"},
		{"doppelt in dritter Ebene", "replay:\n  listen:\n    a: 1\n    a: 2\n", "ungültiges YAML in Zeile 4, Schlüssel doppelt"},
		{"doppelt in vierter Ebene", "replay:\n  listen:\n    a:\n      b: 1\n      b: 2\n", "ungültiges YAML in Zeile 5, Schlüssel doppelt"},
		{"doppelt in einer Liste", "replay:\n  listen:\n    - a: 1\n      a: 2\n", "ungültiges YAML in Zeile 4, Schlüssel doppelt"},
		{"leerer Schlüssel oben", "\"\": 1\n", "Konfigurationsdatei: oberste Ebene: unbekannter Schlüssel"},
		{"Schlüssel kein Skalar oben", "? [a]\n: 1\n", "Konfigurationsdatei: oberste Ebene: unbekannter Schlüssel"},
		{"leerer Schlüssel im Abschnitt", "record:\n  \"\": 1\n", "Konfigurationsdatei: record: unbekannter Schlüssel"},
		{"Schlüssel kein Skalar im Abschnitt", "record:\n  ? [a]\n  : 1\n", "Konfigurationsdatei: record: unbekannter Schlüssel"},
		{"oberste Ebene Liste", "- log_level\n", "oberste Ebene"},
		{"oberste Ebene Skalar", "GEHEIM\n", "oberste Ebene"},
		{"nur ---", "---\n", "oberste Ebene"},
		{"Abschnitt ohne Inhalt", "replay:\n", "replay"},
		{"Abschnitt nur Kommentar", "replay:\n  # listen: x\n", "replay"},
		{"Abschnitt null", "replay: null\n", "replay"},
		{"Abschnitt Liste", "replay:\n  - listen\n", "replay"},
		{"connections ohne Inhalt", "connections:\n", "connections"},
		{"connections null", "connections: ~\n", "connections"},
		{"unbekannter Schlüssel", "bogus: 1\n", "bogus"},
		{"Groß- und Kleinschreibung", "Log_Level: info\n", "Log_Level"},
		{"Schlüssel config oben", "config: x.yaml\n", "config"},
		{"Schlüssel config im Abschnitt", "replay:\n  config: x.yaml\n", "replay.config"},
		{"Abschnitt play", "play:\n  input: x\n", "play"},
		{"leerer Abschnitt play", "play: {}\n", "play"},
		{"Alias ohne Anker", "replay:\n  listen: *GEHEIM\n", "Konfigurationsdatei: ungültiges YAML, Alias ohne Anker"},
		{"Option des anderen Kommandos", "record:\n  input: GEHEIM\n", "record.input"},
		{"Option des anderen Kommandos replay", "replay:\n  output: GEHEIM\n", "replay.output"},
		{"Option, die der Stand nicht kennt", "record:\n  format: yaml\n", "record.format"},
		{"Merge-Schlüssel", "replay:\n  <<: {listen: GEHEIM}\n", "replay.<<"},
		{"anderer Abschnitt mitgeprüft", "record:\n  force: GEHEIM\n", "record.force"},
		{"leerer Wert", "replay:\n  listen:\n", "replay.listen"},
		{"null", "replay:\n  listen: null\n", "replay.listen"},
		{"Tilde", "replay:\n  listen: ~\n", "replay.listen"},
		{"Null", "record:\n  output: Null\n", "record.output"},
		{"NULL", "record:\n  output: NULL\n", "record.output"},
		{"leer in Anführungszeichen", "replay:\n  listen: \"\"\n", "replay.listen"},
		{"Liste als Wert", "replay:\n  listen: [GEHEIM]\n", "replay.listen"},
		{"Abbildung als Wert", "replay:\n  listen: {a: GEHEIM}\n", "replay.listen"},
		{"Anker", "replay:\n  listen: &a GEHEIM\n", "replay.listen"},
		{"Alias", "replay:\n  listen: &a GEHEIM\n  input: *a\n", "replay.listen"},
		{"Tag !!str", "replay:\n  listen: !!str GEHEIM\n", "replay.listen"},
		{"Tag !x", "replay:\n  listen: !x GEHEIM\n", "replay.listen"},
		{"Tag !", "replay:\n  listen: ! GEHEIM\n", "replay.listen"},
		{"Tag am Abschnitt", "replay: !!map\n  listen: GEHEIM\n", "replay"},
		{"Tag am Schlüssel", "replay:\n  !!str listen: GEHEIM\n", "replay"},
		{"Anker am Abschnitt", "replay: &r\n  listen: GEHEIM\n", "replay"},
		{"Tag oben", "!!map\nlog_level: info\n", "oberste Ebene"},
		{"Wahrheitswert True", "replay:\n  fail_on_unconsumed: True\n", "replay.fail_on_unconsumed"},
		{"Wahrheitswert yes", "replay:\n  fail_on_unconsumed: yes\n", "replay.fail_on_unconsumed"},
		{"Wahrheitswert 1", "replay:\n  fail_on_unconsumed: 1\n", "replay.fail_on_unconsumed"},
		{"Dauer ohne Einheit", "replay:\n  shutdown_timeout: 5\n", "replay.shutdown_timeout"},
		{"Dauer 00", "replay:\n  shutdown_timeout: 00\n", "replay.shutdown_timeout"},
		{"log_level groß", "log_level: INFO\n", "log_level"},
		{"log_level im Abschnitt", "replay:\n  log_level: info\n", "replay.log_level"},
		{"Option oben", "listen: GEHEIM\n", "listen"},
		{"Name einer Verbindung leer", "connections:\n  \"\": x\n", "connections"},
		{"Name einer Verbindung null", "connections:\n  ~: x\n", "connections"},
		{"Name einer Verbindung mit Tag", "connections:\n  !!str a: x\n", "connections"},
		{"Name einer Verbindung mit Steuerzeichen", "connections:\n  \"a\\tb\": x\n", "connections"},
		{"Name einer Verbindung mit DEL", "connections:\n  \"a\\x7fb\": x\n", "Name einer Verbindung mit Steuerzeichen"},
		{"Name einer Verbindung mit C1", "connections:\n  \"a\\x9fb\": x\n", "Name einer Verbindung mit Steuerzeichen"},
		{"Wert einer Verbindung leer", "connections:\n  a:\n", "connections.a"},
		{"Wert einer Verbindung Liste", "connections:\n  a: [GEHEIM]\n", "connections.a"},
	}
}

// Abdeckung: LH-FA-17/Negative — eine ungültige Datei ist PGR-E2004 (zu
// ungültigem YAML mit genauer Meldung und Zeile ab 1, bei einem Fehler des
// Parsers dessen Zahl plus 1, ohne Zahl Zeile 1): mehr als
// ein Dokument, ungültiges YAML, ein Schlüssel zweimal in derselben Abbildung,
// oberste Ebene, Abschnitt oder connections: keine Abbildung oder ohne Inhalt,
// ein unbekannter Schlüssel (auch config, ein anders geschriebener, einer im
// falschen Abschnitt, einer des Abschnitts play: und einer, den der Stand noch
// nicht kennt), ein Merge-Schlüssel, ein leerer Wert, null, Liste oder
// Abbildung als Wert, Anker, Alias und jeder ausdrücklich geschriebene Tag, ein
// Wert außerhalb der Wertemenge, ein ungültiger Name einer Verbindung; geprüft
// wird auch der Abschnitt des anderen Kommandos; die Meldung nennt die Stelle
// und nie den Wert (LH-FA-17.a).
func TestDateiUngueltig(t *testing.T) {
	leere(t, "replay")
	for _, f := range ungueltigeDateien() {
		_, err := replayMit("--config=" + schreibe(t, f.inhalt))
		if !istDatei(err) || !strings.Contains(err.Error(), f.stelle) || strings.Contains(err.Error(), "GEHEIM") {
			t.Errorf("%s (%q): erwartet %s mit %q ohne Wert, erhalten %v", f.name, f.inhalt, model.CodeConfigFile, f.stelle, err)
		}
		if genau := "Konfiguration [PGR-E2004]: " + f.stelle; strings.HasPrefix(f.stelle, "Konfigurationsdatei: ungültiges YAML") && err != nil && err.Error() != genau {
			t.Errorf("%s: Meldung %q, erwartet genau %q", f.name, err.Error(), genau)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — der Start endet beim ersten Fehler: zuerst
// die Kommandozeile, dann die Umgebungsvariablen, dann die Datei in ihrer
// Reihenfolge, zuletzt die Pflichtoptionen (LH-FA-17.a *Fehler*).
func TestDateiReihenfolge(t *testing.T) {
	leere(t, "replay")
	ungueltig := schreibe(t, "replay:\n  shutdown_timeout: 5\n")
	if _, err := lese("replay", "--config="+ungueltig, "--log-level=trace"); !istUsage(err) {
		t.Errorf("Kommandozeile vor Datei: %v", err)
	}
	t.Setenv(cli.EnvLogLevel, "trace")
	if _, err := lese("replay", "--config="+ungueltig); !istUsage(err) || !strings.Contains(err.Error(), cli.EnvLogLevel) {
		t.Errorf("Umgebung vor Datei: %v", err)
	}
	t.Setenv(cli.EnvLogLevel, "")
	if _, err := lese("replay", "--config="+ungueltig); !istDatei(err) {
		t.Errorf("Datei vor Pflichtoptionen: %v", err)
	}
	zwei := schreibe(t, "replay:\n  fail_on_unconsumed: 1\n  shutdown_timeout: 5\n")
	if _, err := lese("replay", "--config="+zwei); !istDatei(err) || !strings.Contains(err.Error(), "fail_on_unconsumed") || strings.Contains(err.Error(), "shutdown_timeout") {
		t.Errorf("Datei in ihrer Reihenfolge: %v", err)
	}
	vorn := schreibe(t, "record:\n  force: ja\nlog_level: TRACE\n")
	if _, err := lese("replay", "--config="+vorn); !istDatei(err) || !strings.Contains(err.Error(), "record.force") {
		t.Errorf("Abschnitt des anderen Kommandos in der Reihenfolge der Datei: %v", err)
	}
	if _, err := lese("replay", "--config="+schreibe(t, "replay:\n  input: r.yaml\n")); !istUsage(err) || !strings.Contains(err.Error(), "--listen") {
		t.Errorf("Pflichtoption nach der Datei: %v", err)
	}
}

// konfigurationZeigen ruft config show und liefert die Anzeige.
func konfigurationZeigen(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	cmd, err := cli.Parse(append([]string{"config", "show"}, args...), &out)
	if out.Len() > 0 {
		t.Fatalf("config show schreibt selbst: %q", out.String())
	}
	if err == nil && cmd.Name != "config show" {
		t.Fatalf("Kommando %q", cmd.Name)
	}
	return cmd.Anzeige, err
}

// Abdeckung: LH-FA-17/Happy — config show nennt in der ersten Zeile den Pfad
// der gewählten Datei, danach ihren Inhalt als YAML mit zwei Leerzeichen Einzug
// in der Reihenfolge der Datei, ohne Kommentare, Platzhalter unaufgelöst,
// danach die Namen der gesetzten, nicht leeren Umgebungsvariablen mit dem
// Präfix, auch ohne passende Option, nach Namen sortiert und ohne Werte
// (LH-FA-17.a *Anzeige*).
func TestConfigShow(t *testing.T) {
	leere(t, "replay")
	leere(t, "record")
	t.Setenv("PGWIRE_RECORDER_ZZZ", "wert-z")
	t.Setenv("PGWIRE_RECORDER_BOGUS", "wert-bogus")
	t.Setenv("PGWIRE_RECORDER_LEER", "")
	t.Setenv("X_PGWIRE_RECORDER_MITTE", "wert-mitte")
	pfad := schreibe(t, "# Kopf\nreplay:    # Abschnitt\n    listen: \":5432\"\n    input:   ./r.yaml\nlog_level: warn\nconnections:\n    lokal: \"postgresql://${BENUTZER}@h/db\"\n# Ende\n")
	got, err := konfigurationZeigen(t, "--config", pfad)
	want := pfad + "\nreplay:\n  listen: \":5432\"\n  input: ./r.yaml\nlog_level: warn\nconnections:\n  lokal: \"postgresql://${BENUTZER}@h/db\"\nPGWIRE_RECORDER_BOGUS\nPGWIRE_RECORDER_ZZZ\n"
	if err != nil || got != want {
		t.Errorf("Anzeige:\n%s\n%v\nerwartet:\n%s", got, err, want)
	}
	if strings.Contains(got, "wert-") {
		t.Errorf("Anzeige nennt einen Wert einer Umgebungsvariable:\n%s", got)
	}
}

// Abdeckung: LH-FA-17/Boundary — ohne Datei sagt config show das in der ersten
// Zeile und listet die aktiven Umgebungsvariablen; eine leere Datei zeigt nur
// den Pfad; PGWIRE_RECORDER_CONFIG wählt die Datei, und von den
// Umgebungsvariablen prüft config show nur sie: eine ungültige
// PGWIRE_RECORDER_LOG_LEVEL bleibt ungeprüft; die erste Zeile ist der Pfad, wie
// --config ihn nennt, auch relativ; eine Standarddatei, die ein Link auf ein
// fehlendes Ziel ist, wählt keine Datei (LH-FA-17.a).
func TestConfigShowOhneDatei(t *testing.T) {
	leere(t, "replay")
	t.Chdir(t.TempDir())
	t.Setenv(cli.EnvLogLevel, "TRACE")
	if got, err := konfigurationZeigen(t); err != nil || got != "keine Konfigurationsdatei gefunden\n"+cli.EnvLogLevel+"\n" {
		t.Errorf("ohne Datei: %q, %v", got, err)
	}
	t.Setenv(cli.EnvLogLevel, "")
	leer := schreibe(t, "# nur ein Kommentar\n")
	t.Setenv(cli.EnvConfig, leer)
	if got, err := konfigurationZeigen(t); err != nil || got != leer+"\n"+cli.EnvConfig+"\n" {
		t.Errorf("leere Datei aus PGWIRE_RECORDER_CONFIG: %q, %v", got, err)
	}
	t.Setenv(cli.EnvConfig, "")
	if err := os.WriteFile(cli.StandardDatei, []byte("log_level: info\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := konfigurationZeigen(t); err != nil || got != cli.StandardDatei+"\nlog_level: info\n" {
		t.Errorf("Standarddatei: %q, %v", got, err)
	}
	if err := os.WriteFile("relativ.yaml", []byte("log_level: warn\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := konfigurationZeigen(t, "--config", "relativ.yaml"); err != nil || got != "relativ.yaml\nlog_level: warn\n" {
		t.Errorf("--config mit relativem Pfad: %q, %v", got, err)
	}
	if err := os.Remove(cli.StandardDatei); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("fehlt.yaml", cli.StandardDatei); err != nil {
		t.Fatal(err)
	}
	if got, err := konfigurationZeigen(t); err != nil || got != "keine Konfigurationsdatei gefunden\n" {
		t.Errorf("Standarddatei als Link auf ein fehlendes Ziel: %q, %v", got, err)
	}
}

// Abdeckung: LH-FA-17/Negative — ist die Datei ungültig, endet config show mit
// PGR-E2004 und zeigt nichts, auch bei einer Datei, die nur im Abschnitt eines
// Kommandos ungültig ist; eine fehlende Datei aus PGWIRE_RECORDER_CONFIG ist
// PGR-E2004; config ohne show, ein anderes Unterkommando, ein Argument, --
// mit Argument und eine Option außer --config sind PGR-E2001 (LH-FA-17.a,
// LH-FA-01.a).
func TestConfigShowFehler(t *testing.T) {
	leere(t, "replay")
	for _, f := range ungueltigeDateien() {
		got, err := konfigurationZeigen(t, "--config="+schreibe(t, f.inhalt))
		if !istDatei(err) || got != "" {
			t.Errorf("%s: %q, %v", f.name, got, err)
		}
	}
	t.Setenv(cli.EnvConfig, filepath.Join(t.TempDir(), "fehlt.yaml"))
	if _, err := konfigurationZeigen(t); !istDatei(err) {
		t.Errorf("fehlende Datei aus PGWIRE_RECORDER_CONFIG: %v", err)
	}
	t.Setenv(cli.EnvConfig, "")
	for _, args := range [][]string{
		{"config"}, {"config", "zeige"}, {"config", "show", "x"},
		{"config", "show", "--", "x"}, {"config", "show", "--log-level=info"}, {"config", "show", "--config="},
	} {
		if _, err := keineHilfe(t, args...); !istUsage(err) {
			t.Errorf("%q: erwartet %s, erhalten %v", args, model.CodeUsage, err)
		}
	}
}

// Abdeckung: LH-FA-01/Negative, LH-FA-17/Negative — die Hilfe geht jeder
// Prüfung vor, auch der Konfigurationsdatei: Bei record, replay und config show
// geben -h und --help die Hilfe des Kommandos aus, auch mit fehlender Datei aus
// --config oder PGWIRE_RECORDER_CONFIG, ungültiger Standarddatei und ungültigem
// Wert von --config; config --help gibt die Hilfe von config show aus; nach
// -- ist eine Hilfe-Angabe ein gewöhnliches Argument (LH-FA-01.a). Die Hilfe
// von record und replay nennt --config, seine Umgebungsvariable, die
// Standarddatei und die Datei als Quelle der übrigen Optionen mit ihrer
// Priorität.
func TestDateiHilfeVorPruefung(t *testing.T) {
	leere(t, "replay")
	t.Chdir(t.TempDir())
	if err := os.WriteFile(cli.StandardDatei, []byte("bogus: 1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(cli.EnvConfig, filepath.Join(t.TempDir(), "fehlt.yaml"))
	for _, args := range [][]string{
		{"replay", "--help"},
		{"record", "--config", "fehlt.yaml", "-h"},
		{"replay", "--config=", "--help"},
		{"config", "show", "--help"},
		{"config", "show", "--config", "-h"},
		{"config", "show", "--bogus", "--help"},
		{"config", "--help"},
		{"config", "bogus", "-h"},
	} {
		text := hilfe(t, args...)
		if strings.Contains(text, "Kommandos:") {
			t.Errorf("%q: globale statt Hilfe des Kommandos:\n%s", args, text)
		}
		if args[0] == "config" && !strings.Contains(text, "Optionen von config show") {
			t.Errorf("%q: nicht die Hilfe von config show:\n%s", args, text)
		}
	}
	if _, err := keineHilfe(t, "config", "show", "--", "--help"); !istUsage(err) || !strings.Contains(err.Error(), `unerwartetes Argument "--help"`) {
		t.Errorf("config show -- --help: %v", err)
	}
	if text := hilfe(t, "--help"); !strings.Contains(text, "config show") || !strings.Contains(text, "--config") {
		t.Errorf("globale Hilfe ohne config show:\n%s", text)
	}
	for _, kommando := range []string{"record", "replay"} {
		text := strings.Join(strings.Fields(hilfe(t, kommando, "-h")), " ")
		for _, satz := range []string{"--config <datei>", "Standard .pgwire-recorder.yaml im aktuellen Verzeichnis", cli.EnvConfig, "In der Konfigurationsdatei sind alle Optionen ohne --config setzbar", "die Umgebungsvariable geht der Datei vor, die Datei dem Standardwert"} {
			if !strings.Contains(text, satz) {
				t.Errorf("%s: Hilfe ohne %q:\n%s", kommando, satz, text)
			}
		}
	}
}

// Abdeckung: LH-FA-17/Negative — jeder ausdrücklich geschriebene Tag ist
// PGR-E2004, auch der nicht spezifische Tag ! allein, mit führendem BOM, mit
// dem Zeilenende \r allein und mit \r\n; einen Tag, den die Bibliothek als
// TaggedStyle markiert, lehnt form auch ohne Text der Datei ab; ein ! in einem
// Kommentar macht eine gültige Datei mit \r nicht ungültig (LH-FA-17.a).
func TestDateiTag(t *testing.T) {
	leere(t, "replay")
	if err := cli.FormMitTaggedStyle(); err == nil || !strings.Contains(err.Error(), "Tag ist ungültig") {
		t.Errorf("TaggedStyle ohne Text: %v", err)
	}
	for _, tag := range []string{"!", "!!str"} {
		for name, inhalt := range map[string]string{
			"BOM":    "\uFEFFlog_level: " + tag + " info\n",
			"\\r":    "record: {}\rlog_level: " + tag + " info\r",
			"\\r\\n": "record: {}\r\nlog_level: " + tag + " info\r\n",
		} {
			if _, err := replayMit("--config=" + schreibe(t, inhalt)); !istDatei(err) || !strings.Contains(err.Error(), "log_level: Tag ist ungültig") {
				t.Errorf("Tag %s mit %s: %v", tag, name, err)
			}
		}
	}
	cmd, err := replayMit("--config=" + schreibe(t, "record:\r  force: true\nlog_level: debug\n#          !\n"))
	if err != nil || cmd.Replay.LogLevel != cli.LogDebug {
		t.Errorf("! im Kommentar nach \\r: %#v, %v", cmd.Replay, err)
	}
}

// utf16 kodiert text als UTF-16, little oder big endian, mit oder ohne BOM.
func utf16(text string, little, mitBOM bool) []byte {
	var out []byte
	if mitBOM {
		text = "\uFEFF" + text
	}
	for _, r := range text {
		if little {
			out = append(out, byte(r), byte(r>>8))
		} else {
			out = append(out, byte(r>>8), byte(r))
		}
	}
	return out
}

// Abdeckung: LH-FA-17/Boundary, LH-FA-17/Negative — die Datei ist UTF-8: ein
// BOM als erstes Zeichen wird übergangen, an anderer Stelle (Anfang von Zeile 2,
// Wert in Anführungszeichen) ist es PGR-E2004 mit der Zeile; UTF-16 mit und ohne
// BOM und eine ungültige UTF-8-Folge sind PGR-E2004, die Meldung nennt kein Byte
// der Datei. Zeilenenden sind \n, \r\n und \r; U+0085, U+2028 und U+2029 sind
// PGR-E2004 mit der Zeile, auch in einem Kommentar; ebenso jedes Zeichen, das
// YAML 1.2 nicht als druckbar zulässt (NUL, U+0001, U+007F, U+0080, U+FFFE,
// U+FFFF), auch UTF-16LE ohne BOM in Zeile 1; ein Tabulator bleibt gültig; ein
// Alias ohne Anker nennt die Ursache ohne Zeile und Namen (LH-FA-17.a).
func TestDateiKodierung(t *testing.T) {
	leere(t, "replay")
	for _, inhalt := range []string{
		"\uFEFFreplay:\n  fail_on_unconsumed: true\n",
		"replay:\r  fail_on_unconsumed: true\r",
		"replay:\r\n  fail_on_unconsumed: true\r\n",
	} {
		if cmd, err := replayMit("--config=" + schreibe(t, inhalt)); err != nil || !cmd.Replay.FailOnUnconsumed {
			t.Errorf("%q: %#v, %v", inhalt, cmd.Replay, err)
		}
	}
	for _, f := range []struct{ inhalt, zeile string }{
		{"log_level: info\n\uFEFFreplay: {}\n", "Zeile 2"},
		{"log_level: info\nrecord: {output: \"a\uFEFFb\"}\n", "Zeile 2"},
		{"log_level: info\nrecord: {output: \"a\u0085b\"}\n", "Zeile 2"},
		{"log_level: info\r\nrecord: {output: \"a\u2028b\"}\n", "Zeile 2"},
		{"log_level: info\rrecord: {output: \"a\u2029b\"}\n", "Zeile 2"},
		{"log_level: info\n# a\u2028b\n", "Zeile 2"},
		{"log_level: info\n# \xff\n", "Zeile 2"},
		{"log_level: info\nrecord: {output: a\x00b}\n", "Zeile 2"},
		{"log_level: info\nrecord: {}\nreplay: {listen: \"a\x01b\"}\n", "Zeile 3"},
		{"log_level: info\nrecord: {output: \"a\x7fb\"}\n", "Zeile 2"},
		{"log_level: info\nrecord: {output: \"a\u0080b\"}\n", "Zeile 2"},
		{"log_level: info\nrecord: {output: \"a\uFFFEb\"}\n", "Zeile 2"},
		{"log_level: info\nrecord: {output: \"a\uFFFFb\"}\n", "Zeile 2"},
	} {
		_, err := replayMit("--config=" + schreibe(t, f.inhalt))
		if !istDatei(err) || !strings.Contains(err.Error(), "ungültiges YAML in "+f.zeile) || strings.ContainsAny(err.Error(), "\uFEFF\u0085\u2028\u2029\uFFFDÿ") || strings.Contains(err.Error(), string([]byte{0xff})) {
			t.Errorf("%q: erwartet %s mit %s ohne Byte der Datei, erhalten %v", f.inhalt, model.CodeConfigFile, f.zeile, err)
		}
	}
	text := "log_level: info\n"
	for name, roh := range map[string][]byte{
		"UTF-16LE mit BOM":  utf16(text, true, true),
		"UTF-16BE mit BOM":  utf16(text, false, true),
		"UTF-16LE ohne BOM": utf16(text, true, false),
	} {
		pfad := filepath.Join(t.TempDir(), "k.yaml")
		if err := os.WriteFile(pfad, roh, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := replayMit("--config=" + pfad); !istDatei(err) || !strings.Contains(err.Error(), "ungültiges YAML in Zeile 1") {
			t.Errorf("%s: %v", name, err)
		}
	}
	if cmd, err := replayMit("--config=" + schreibe(t, "record:\n  output: \"a\tb\"\n")); err != nil || cmd.Name != "replay" {
		t.Errorf("Tabulator in Anführungszeichen: %v", err)
	}
	_, err := replayMit("--config=" + schreibe(t, "a: *x\n"))
	if !istDatei(err) || err.Error() != "Konfiguration [PGR-E2004]: Konfigurationsdatei: ungültiges YAML, Alias ohne Anker" {
		t.Errorf("Alias ohne Anker: %v", err)
	}
}

// Abdeckung: LH-FA-17/Boundary — config show gibt die Datei ohne BOM und mit
// \n als Zeilenende aus, auch wenn sie mit BOM und \r geschrieben ist
// (LH-FA-17.a *Anzeige*).
func TestConfigShowOhneBOM(t *testing.T) {
	leere(t, "replay")
	leere(t, "record")
	pfad := schreibe(t, "\uFEFFlog_level: warn\r# Kommentar\rreplay:\r  listen: x\r")
	got, err := konfigurationZeigen(t, "--config", pfad)
	if err != nil || got != pfad+"\nlog_level: warn\nreplay:\n  listen: x\n" {
		t.Errorf("Anzeige %q, %v", got, err)
	}
}
