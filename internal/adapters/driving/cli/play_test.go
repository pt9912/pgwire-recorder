package cli_test

import (
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// playMit liest play mit --input und den übrigen Argumenten.
func playMit(args ...string) (cli.PlayOptions, error) {
	cmd, err := lese(append([]string{"play", "--input=r.yaml"}, args...)...)
	return cmd.Play, err
}

// Abdeckung: LH-FA-20/Happy, LH-FA-17/Happy — play liest --upstream, --input,
// --user, --database und --log-level; ohne --user und --database bleiben beide
// leer, und die Startup-Daten der Session gelten (LH-FA-17.a, LH-FA-20.a).
func TestParsePlay(t *testing.T) {
	leere(t, "play")
	cmd, err := lese("play", "--upstream", "pg:5432", "--input", "r.yaml")
	if err != nil || cmd.Name != "play" || cmd.Play != (cli.PlayOptions{Upstream: "pg:5432", Input: "r.yaml", LogLevel: cli.LogInfo}) {
		t.Fatalf("play ohne --user und --database: %#v, %v", cmd, err)
	}
	cmd, err = lese("play", "--upstream=pg:5432", "--input=r.yaml", "--user=u", "--database=d", "--log-level=debug")
	if err != nil || cmd.Play != (cli.PlayOptions{Upstream: "pg:5432", Input: "r.yaml", User: "u", Database: "d", LogLevel: cli.LogDebug}) {
		t.Fatalf("play mit allen Optionen: %#v, %v", cmd, err)
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-01/Negative — play nimmt kein
// gewöhnliches Argument an, auch nicht nach -- (PGR-E2001, LH-FA-01.a).
func TestParsePlayArgument(t *testing.T) {
	leere(t, "play")
	for _, args := range [][]string{{"zusatz"}, {"--", "zusatz"}} {
		if _, err := playMit(append([]string{"--upstream=pg:1"}, args...)...); !istUsage(err) || !strings.Contains(err.Error(), `unerwartetes Argument "zusatz"`) {
			t.Errorf("%q: erwartet unerwartetes Argument, erhalten %v", args, err)
		}
	}
}

// Abdeckung: LH-FA-03/Negative, LH-FA-17/Negative — play kennt
// --fail-on-unconsumed nicht (PGR-E2001) und lässt
// PGWIRE_RECORDER_FAIL_ON_UNCONSUMED unbeachtet, auch mit ungültigem Wert
// (LH-FA-03.b §Andere Kommandos).
func TestParsePlayFailOnUnconsumed(t *testing.T) {
	leere(t, "play")
	t.Setenv(cli.EnvFailOnUnconsumed, "1")
	if _, err := playMit("--upstream=pg:1", "--fail-on-unconsumed"); !istUsage(err) || !strings.Contains(err.Error(), "flag provided but not defined") {
		t.Fatalf("play mit --fail-on-unconsumed: %v", err)
	}
	if _, err := playMit("--upstream=pg:1"); err != nil {
		t.Fatalf("play wertet PGWIRE_RECORDER_FAIL_ON_UNCONSUMED aus: %v", err)
	}
	if _, err := lese("play", "--upstream=pg:1", "--input=r.yaml", "--config="+schreibe(t, "play:\n  fail_on_unconsumed: true\n")); !istDatei(err) || !strings.Contains(err.Error(), "play.fail_on_unconsumed") {
		t.Fatalf("Schlüssel fail_on_unconsumed in play: %v", err)
	}
}

// Abdeckung: LH-FA-17/Boundary — bei play bleiben die Umgebungsvariablen der
// Optionen unbeachtet, die dieser Stand von play nicht kennt (Laufsteuerung,
// TLS, Zeitangaben, Vergleich), auch mit ungültigem Wert, ebenso
// PGWIRE_RECORDER_PASSWORD (LH-FA-17.a).
func TestParsePlayFremdeUmgebung(t *testing.T) {
	leere(t, "play")
	for _, name := range []string{"CONTINUE_ON_ERROR", "ALLOW_RECORDED_ERRORS", "FINISH_SESSION_ON_INTERRUPT", "UPSTREAM_TLS", "UPSTREAM_CA", "KEEP_TIMING", "TIMING_MODE", "TIMING_REFERENCE", "COMPARE_RESPONSES", "PASSWORD", "LISTEN", "SHUTDOWN_TIMEOUT"} {
		t.Setenv("PGWIRE_RECORDER_"+name, "ungültig")
	}
	if got, err := playMit("--upstream=pg:1"); err != nil || got != (cli.PlayOptions{Upstream: "pg:1", Input: "r.yaml", LogLevel: cli.LogInfo}) {
		t.Fatalf("fremde Umgebungsvariable ausgewertet: %#v, %v", got, err)
	}
}

// Abdeckung: LH-FA-17/Happy, LH-FA-20/Happy — nennt --upstream bei play eine
// Verbindung, ist Upstream ihre Adresse host:port, und Benutzer und Datenbank
// der URL gelten, eingesetzt; --user und --database gehen ihnen vor, auch aus
// der Umgebung; ohne Benutzer in der URL bleibt er leer (LH-FA-17.a *Wirkung
// einer URL*).
func TestPlayVerbindung(t *testing.T) {
	leere(t, "play")
	variablen(t, map[string]*string{"PGR_T_U": w("app"), "PGR_T_PW": w("GEHEIM"), "PGR_T_H": w("db.example"), "PGR_T_P": w("6543"), "PGR_T_DB": w("shop")})
	url := "postgresql://${PGR_T_U}:${PGR_T_PW}@${PGR_T_H}:${PGR_T_P}/${PGR_T_DB}"
	config := "--config=" + schreibe(t, mitVerbindung(url))
	if got, err := playMit("--upstream=v", config); err != nil || got != (cli.PlayOptions{Upstream: "db.example:6543", Input: "r.yaml", User: "app", Database: "shop", LogLevel: cli.LogInfo}) {
		t.Errorf("Verbindung: %#v, %v", got, err)
	}
	if got, err := playMit("--upstream=v", config, "--user=ich", "--database=meine"); err != nil || got.User != "ich" || got.Database != "meine" || got.Upstream != "db.example:6543" {
		t.Errorf("--user und --database vor der URL: %#v, %v", got, err)
	}
	t.Setenv("PGWIRE_RECORDER_USER", "env-u")
	t.Setenv("PGWIRE_RECORDER_DATABASE", "env-d")
	if got, err := playMit("--upstream=v", config); err != nil || got.User != "env-u" || got.Database != "env-d" {
		t.Errorf("Umgebung vor der URL: %#v, %v", got, err)
	}
	t.Setenv("PGWIRE_RECORDER_USER", "")
	t.Setenv("PGWIRE_RECORDER_DATABASE", "")
	if got, err := playMit("--upstream=v", "--config="+schreibe(t, mitVerbindung("postgresql://[::1]/db"))); err != nil || got.User != "" || got.Database != "db" || got.Upstream != "[::1]:5432" {
		t.Errorf("URL ohne Benutzer: %#v, %v", got, err)
	}
	if got, err := playMit("--upstream=h:7", "--user=u"); err != nil || got.Upstream != "h:7" || got.User != "u" || got.Database != "" {
		t.Errorf("host:port: %#v, %v", got, err)
	}
}

// Abdeckung: LH-FA-17/Negative, LH-FA-20/Negative — bei play ist eine nicht
// gesetzte oder leere Variable in jedem Teil der benutzten Verbindung
// PGR-E2005 mit der ersten in der Reihenfolge der URL (Benutzer, Passwort,
// Host, Port, Datenbank), auch in einem Teil, den --user oder --database
// überschreibt, und im Passwort (U8); sslmode=require ist PGR-E2004 vor den
// Variablen, ein Port, der nach dem Einsetzen keiner ist, PGR-E2004 nach ihnen
// (LH-FA-17.a).
func TestPlayVariablen(t *testing.T) {
	leere(t, "play")
	fehlt := func(name string) string {
		return "Konfiguration [PGR-E2005]: Konfigurationsdatei: connections.v: Umgebungsvariable " + name + " eines Platzhalters nicht gesetzt"
	}
	url := "postgresql://${PGR_T_U}:${PGR_T_PW}@${PGR_T_H}:${PGR_T_P}/${PGR_T_DB}"
	alle := []string{"PGR_T_U", "PGR_T_PW", "PGR_T_H", "PGR_T_P", "PGR_T_DB"}
	for i, erste := range alle {
		werte := map[string]*string{}
		for j, name := range alle {
			switch {
			case j < i:
				werte[name] = w("1")
			case j == i && i%2 == 0:
				werte[name] = w("")
			default:
				werte[name] = nil
			}
		}
		variablen(t, werte)
		_, err := playMit("--upstream=v", "--user=u", "--database=d", "--config="+schreibe(t, mitVerbindung(url)))
		if err == nil || err.Error() != fehlt(erste) || !hatCode(err, model.CodeConfigVariable) {
			t.Errorf("erste fehlende %s: %v, erwartet %q", erste, err, fehlt(erste))
		}
	}
	variablen(t, map[string]*string{"PGR_T_U": nil, "PGR_T_PW": nil, "PGR_T_H": nil, "PGR_T_P": w("GEHEIM"), "PGR_T_DB": nil})
	_, err := playMit("--upstream=v", "--config="+schreibe(t, mitVerbindung(url+"?sslmode=require")))
	if !istDatei(err) || err.Error() != "Konfiguration [PGR-E2004]: Konfigurationsdatei: connections.v: sslmode=require ist bei play ungültig, play verbindet ohne TLS zum Upstream" {
		t.Errorf("sslmode=require vor den Variablen: %v", err)
	}
	variablen(t, map[string]*string{"PGR_T_U": w("1"), "PGR_T_PW": w("1"), "PGR_T_H": w("h"), "PGR_T_P": w("GEHEIM"), "PGR_T_DB": nil})
	_, err = playMit("--upstream=v", "--config="+schreibe(t, mitVerbindung(url)))
	if err == nil || err.Error() != fehlt("PGR_T_DB") {
		t.Errorf("Variablen vor dem Port: %v", err)
	}
	variablen(t, map[string]*string{"PGR_T_DB": w("d")})
	_, err = playMit("--upstream=v", "--config="+schreibe(t, mitVerbindung(url)))
	if !istDatei(err) || strings.Contains(err.Error(), "GEHEIM") || !strings.Contains(err.Error(), "Port nach dem Einsetzen") {
		t.Errorf("Port nach dem Einsetzen: %v", err)
	}
}

// Die Hilfe von play nennt seine Optionen und den Abschnitt play: der
// Konfigurationsdatei (LH-FA-01.a); die globale Hilfe nennt play.
func TestParsePlayHilfe(t *testing.T) {
	text := hilfe(t, "play", "--help")
	for _, teil := range []string{"Optionen von play:", "--upstream", "--input", "--user", "--database", "--log-level", "--config", "play:"} {
		if !strings.Contains(text, teil) {
			t.Errorf("Hilfe von play ohne %q:\n%s", teil, text)
		}
	}
	if global := hilfe(t, "--help"); !strings.Contains(global, "  play     ") || !strings.Contains(global, "Optionen von play:") {
		t.Errorf("globale Hilfe ohne play:\n%s", global)
	}
}
