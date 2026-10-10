package cli_test

import (
	"fmt"
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
	if err != nil || cmd.Play != (cli.PlayOptions{Upstream: "pg:5432", Input: "r.yaml", Einspielen: cli.Einspielvorgaben{User: "u", Database: "d"}, LogLevel: cli.LogDebug}) {
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
// Optionen unbeachtet, die dieser Stand von play nicht kennt (TLS,
// Zeitangaben, Vergleich), auch mit ungültigem Wert (LH-FA-17.a).
func TestParsePlayFremdeUmgebung(t *testing.T) {
	leere(t, "play")
	t.Setenv(cli.EnvPassword, "")
	for _, name := range []string{"UPSTREAM_TLS", "UPSTREAM_CA", "KEEP_TIMING", "TIMING_MODE", "TIMING_REFERENCE", "COMPARE_RESPONSES", "LISTEN", "SHUTDOWN_TIMEOUT"} {
		t.Setenv("PGWIRE_RECORDER_"+name, "ungültig")
	}
	if got, err := playMit("--upstream=pg:1"); err != nil || got != (cli.PlayOptions{Upstream: "pg:1", Input: "r.yaml", LogLevel: cli.LogInfo}) {
		t.Fatalf("fremde Umgebungsvariable ausgewertet: %#v, %v", got, err)
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-17/Boundary — --continue-on-error,
// --allow-recorded-errors und --finish-session-on-interrupt setzen je genau ihr
// Feld der Einspielvorgaben: ohne Wert, mit =true, aus der Umgebungsvariable
// und aus dem Abschnitt play: true, mit =false und ohne Quelle false; ohne
// Wert geht die Option einer Umgebungsvariable mit false vor; 1, True und yes
// sind auf der Kommandozeile PGR-E2001 (LH-FA-17.a, LH-FA-20.a).
func TestParsePlayLaufsteuerung(t *testing.T) {
	for _, f := range []struct {
		name, env, schluessel string
		want                  cli.Einspielvorgaben
	}{
		{"continue-on-error", "PGWIRE_RECORDER_CONTINUE_ON_ERROR", "continue_on_error", cli.Einspielvorgaben{ContinueOnError: true}},
		{"allow-recorded-errors", "PGWIRE_RECORDER_ALLOW_RECORDED_ERRORS", "allow_recorded_errors", cli.Einspielvorgaben{AllowRecordedErrors: true}},
		{"finish-session-on-interrupt", "PGWIRE_RECORDER_FINISH_SESSION_ON_INTERRUPT", "finish_session_on_interrupt", cli.Einspielvorgaben{FinishSessionOnInterrupt: true}},
	} {
		t.Run(f.name, func(t *testing.T) {
			leere(t, "play")
			for _, args := range [][]string{{"--" + f.name}, {"--" + f.name + "=true"}, {"--config=" + schreibe(t, "play:\n  "+f.schluessel+": true\n")}} {
				if got, err := playMit(append([]string{"--upstream=pg:1"}, args...)...); err != nil || got.Einspielen != f.want {
					t.Errorf("%q: %#v, %v, erwartet %#v", args, got.Einspielen, err, f.want)
				}
			}
			for _, args := range [][]string{nil, {"--" + f.name + "=false"}} {
				if got, err := playMit(append([]string{"--upstream=pg:1"}, args...)...); err != nil || got.Einspielen != (cli.Einspielvorgaben{}) {
					t.Errorf("%q: %#v, %v, erwartet ohne Wirkung", args, got.Einspielen, err)
				}
			}
			for _, w := range []string{"1", "True", "yes"} {
				if _, err := playMit("--upstream=pg:1", "--"+f.name+"="+w); !istUsage(err) {
					t.Errorf("--%s=%s: erwartet %s, erhalten %v", f.name, w, model.CodeUsage, err)
				}
			}
			t.Setenv(f.env, "true")
			if got, err := playMit("--upstream=pg:1"); err != nil || got.Einspielen != f.want {
				t.Errorf("%s=true: %#v, %v", f.env, got.Einspielen, err)
			}
			t.Setenv(f.env, "false")
			if got, err := playMit("--upstream=pg:1", "--"+f.name); err != nil || got.Einspielen != f.want {
				t.Errorf("%s=false neben --%s: %#v, %v", f.env, f.name, got.Einspielen, err)
			}
		})
	}
}

// Abdeckung: LH-FA-17/Happy, LH-FA-20/Happy — nennt --upstream bei play eine
// Verbindung, ist Upstream ihre Adresse host:port, und Benutzer und Datenbank
// der URL gelten, eingesetzt; --user und --database gehen ihnen vor, auch aus
// der Umgebung; ohne Benutzer in der URL bleibt er leer (LH-FA-17.a *Wirkung
// einer URL*).
func TestPlayVerbindung(t *testing.T) {
	leere(t, "play")
	t.Setenv(cli.EnvPassword, "")
	variablen(t, map[string]*string{"PGR_T_U": w("app"), "PGR_T_PW": w("GEHEIM"), "PGR_T_H": w("db.example"), "PGR_T_P": w("6543"), "PGR_T_DB": w("shop")})
	url := "postgresql://${PGR_T_U}:${PGR_T_PW}@${PGR_T_H}:${PGR_T_P}/${PGR_T_DB}"
	config := "--config=" + schreibe(t, mitVerbindung(url))
	if got, err := playMit("--upstream=v", config); err != nil || got != (cli.PlayOptions{Upstream: "db.example:6543", Input: "r.yaml", Einspielen: cli.Einspielvorgaben{User: "app", Database: "shop"}, Passwort: "GEHEIM", LogLevel: cli.LogInfo}) {
		t.Errorf("Verbindung: %#v, %v", got, err)
	}
	if got, err := playMit("--upstream=v", config, "--user=ich", "--database=meine"); err != nil || got.Einspielen.User != "ich" || got.Einspielen.Database != "meine" || got.Upstream != "db.example:6543" {
		t.Errorf("--user und --database vor der URL: %#v, %v", got, err)
	}
	t.Setenv("PGWIRE_RECORDER_USER", "env-u")
	t.Setenv("PGWIRE_RECORDER_DATABASE", "env-d")
	if got, err := playMit("--upstream=v", config); err != nil || got.Einspielen.User != "env-u" || got.Einspielen.Database != "env-d" {
		t.Errorf("Umgebung vor der URL: %#v, %v", got, err)
	}
	t.Setenv("PGWIRE_RECORDER_USER", "")
	t.Setenv("PGWIRE_RECORDER_DATABASE", "")
	if got, err := playMit("--upstream=v", "--config="+schreibe(t, mitVerbindung("postgresql://[::1]/db"))); err != nil || got.Einspielen.User != "" || got.Einspielen.Database != "db" || got.Upstream != "[::1]:5432" {
		t.Errorf("URL ohne Benutzer: %#v, %v", got, err)
	}
	if got, err := playMit("--upstream=h:7", "--user=u"); err != nil || got.Upstream != "h:7" || got.Einspielen.User != "u" || got.Einspielen.Database != "" {
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
	for _, teil := range []string{"Optionen von play:", "--upstream", "--input", "--user", "--database", "--continue-on-error", "PGWIRE_RECORDER_CONTINUE_ON_ERROR", "--allow-recorded-errors", "PGWIRE_RECORDER_ALLOW_RECORDED_ERRORS", "--finish-session-on-interrupt", "PGWIRE_RECORDER_FINISH_SESSION_ON_INTERRUPT", "--log-level", "--config", "play:", "PGWIRE_RECORDER_PASSWORD", "scram-sha-256"} {
		if !strings.Contains(text, teil) {
			t.Errorf("Hilfe von play ohne %q:\n%s", teil, text)
		}
	}
	if global := hilfe(t, "--help"); !strings.Contains(global, "  play     ") || !strings.Contains(global, "Optionen von play:") {
		t.Errorf("globale Hilfe ohne play:\n%s", global)
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary, LH-FA-17/Boundary — das
// Passwort von play ist der eingesetzte Passwortteil der benutzten
// Verbindung; schreibt sie keinen, und bei host:port, ist es der Wert von
// PGWIRE_RECORDER_PASSWORD; eine leere Variable gilt als nicht gesetzt. Schreibt
// die Verbindung einen Passwortteil, bleibt die Variable unbeachtet, auch wenn
// sie gesetzt ist, und fehlt die Variable des Platzhalters, ist das PGR-E2005,
// ohne dass die Variable einspringt. Das Passwort gilt unverändert, mit
// Leerraum, ohne Prozent-Dekodierung und mit $ (LH-FA-20.a *Passwort*).
func TestPlayPasswort(t *testing.T) {
	leere(t, "play")
	roh := "  p%41$ss wort \t"
	for _, f := range []struct {
		name string
		args []string
		// url ist die benutzte Verbindung, "" heißt host:port.
		url string
		// platzhalter ist der Wert von PGR_T_PW, nil nicht gesetzt.
		platzhalter *string
		// env ist der Wert von PGWIRE_RECORDER_PASSWORD.
		env  string
		want cli.Passwort
	}{
		{"host:port mit Variable", nil, "", nil, roh, cli.Passwort(roh)},
		{"host:port mit leerer Variable", nil, "", nil, "", ""},
		{"Verbindung ohne Passwortteil mit Variable", nil, "postgresql://app@h/db", nil, "AUSENV", "AUSENV"},
		{"Verbindung ohne Passwortteil und ohne Variable", nil, "postgresql://app@h/db", nil, "", ""},
		{"Verbindung ohne Benutzer mit Variable", nil, "postgresql://h/db", nil, "AUSENV", "AUSENV"},
		{"Platzhalter und Variable gesetzt", nil, "postgresql://app:${PGR_T_PW}@h/db", w(roh), "AUSENV", cli.Passwort(roh)},
		{"Platzhalter, Variable leer", nil, "postgresql://app:${PGR_T_PW}@h/db", w(roh), "", cli.Passwort(roh)},
		{"Platzhalter und --user", []string{"--user=ich"}, "postgresql://app:${PGR_T_PW}@h/db", w("AUSPLATZ"), "AUSENV", "AUSPLATZ"},
	} {
		t.Run(f.name, func(t *testing.T) {
			variablen(t, map[string]*string{"PGR_T_PW": f.platzhalter})
			t.Setenv(cli.EnvPassword, f.env)
			args := append([]string{"--upstream=h:7"}, f.args...)
			if f.url != "" {
				args = append([]string{"--upstream=v", "--config=" + schreibe(t, mitVerbindung(f.url))}, f.args...)
			}
			got, err := playMit(args...)
			if err != nil || got.Passwort != f.want {
				t.Fatalf("Passwort %q, Fehler %v, erwartet %q", string(got.Passwort), err, string(f.want))
			}
		})
	}
	variablen(t, map[string]*string{"PGR_T_PW": nil})
	t.Setenv(cli.EnvPassword, "AUSENV")
	if _, err := playMit("--upstream=v", "--config="+schreibe(t, mitVerbindung("postgresql://app:${PGR_T_PW}@h/db"))); err == nil || !hatCode(err, model.CodeConfigVariable) || strings.Contains(err.Error(), "AUSENV") {
		t.Errorf("Platzhalter ohne Variable neben gesetzter PGWIRE_RECORDER_PASSWORD: %v, erwartet %s", err, model.CodeConfigVariable)
	}
}

// Abdeckung: LH-FA-20/Negative — eine Formatierung der
// Optionen von play (%v, %+v, %#v, %s, %q, auch über einen Zeiger) gibt das
// Passwort nicht aus (LH-FA-20.a *Passwort*).
func TestPlayOptionenOhnePasswort(t *testing.T) {
	o := cli.PlayOptions{Upstream: "h:1", Input: "r.yaml", Passwort: "GEHEIMPW"}
	for _, verb := range []string{"%v", "%+v", "%#v", "%s", "%q", "%x", "%d"} {
		for _, arg := range []any{o, &o, []cli.PlayOptions{o}, map[string]cli.PlayOptions{"k": o}} {
			if got := fmt.Sprintf(verb, arg); strings.Contains(got, "GEHEIMPW") || strings.Contains(strings.ToLower(got), "47454845494d5057") {
				t.Errorf("%s gibt das Passwort aus: %s", verb, got)
			}
		}
	}
	if got := fmt.Sprintf("%+v", o); !strings.Contains(got, "h:1") {
		t.Errorf("%%+v der Optionen ohne Upstream: %s", got)
	}
}
