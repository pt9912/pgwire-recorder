package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// zeileSpec ist eine Zeile der Optionstabelle in LH-FA-17.a: der Name der Option
// und die Kommandos der Spalte *Kommando* (record, replay, play). Die Tabelle
// steht hart kodiert im Test: Der Build-Kontext der Stufe test schließt spec/
// aus, der Test kann sie nicht lesen.
type zeileSpec struct {
	name      string
	kommandos string
}

// tabelleSpec ist die Optionstabelle ohne --config in der Reihenfolge der
// Spezifikation.
func tabelleSpec() []zeileSpec {
	return []zeileSpec{
		{"listen", "record replay"}, {"upstream", "record play"}, {"output", "record"}, {"force", "record"},
		{"format", "record"}, {"record-timing", "record"}, {"record-empty-sessions", "record"},
		{"input", "replay play"}, {"fail-on-unconsumed", "replay"}, {"shutdown-timeout", "record replay"},
		{"session-assignment", "replay"}, {"user", "play"}, {"database", "play"}, {"continue-on-error", "play"},
		{"allow-recorded-errors", "play"}, {"upstream-tls", "play"}, {"finish-session-on-interrupt", "play"},
		{"keep-timing", "play"}, {"timing-mode", "play"}, {"timing-reference", "play"},
		{"tls-cert", "record replay"}, {"tls-key", "record replay"}, {"allow-plaintext", "record replay"},
		{"upstream-ca", "play"}, {"compare-responses", "play"}, {"log-level", "record replay play"},
	}
}

// nichtGelesen sind die Optionen der Tabelle, die der Stand noch nicht liest:
// bei den eigenen Kommandos bleibt ihre Variable unbeachtet. Der Stand liest
// jede Option der Tabelle.
func nichtGelesen() map[string]bool {
	return map[string]bool{}
}

func (z zeileSpec) kennt(kommando string) bool {
	for _, k := range strings.Fields(z.kommandos) {
		if k == kommando {
			return true
		}
	}
	return false
}

func envVon(name string) string {
	return "PGWIRE_RECORDER_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
}

// leereAlle setzt die Variablen aller Optionen der Tabelle, PASSWORD und CONFIG
// leer.
func leereAlle(t *testing.T) {
	t.Helper()
	for _, z := range tabelleSpec() {
		t.Setenv(envVon(z.name), "")
	}
	t.Setenv(cli.EnvPassword, "")
	t.Setenv(cli.EnvConfig, "")
}

// meldetFremd meldet, ob err PGR-E2001 ist, den Namen env nennt und den Wert
// wert nicht.
func meldetFremd(err error, env, wert string) bool {
	return istUsage(err) && strings.Contains(err.Error(), env) && (wert == "" || !strings.Contains(err.Error(), wert))
}

// Abdeckung: LH-FA-17/Negative — bei record, replay und play ist eine gesetzte,
// nicht leere Umgebungsvariable jeder Option, die das Kommando laut Spalte
// Kommando der Tabelle nicht kennt, PGR-E2001, gleich welchen Wert sie hat
// (gültig wie true und 5s oder ungültig), und die Meldung nennt den Namen der
// Variable, nicht den Wert; eine leere Variable gilt als nicht gesetzt (LH-FA-17.a).
func TestLeserUmgebungAndererKommandos(t *testing.T) {
	for _, kommando := range leserKommandos() {
		for _, z := range tabelleSpec() {
			if z.kennt(kommando) {
				continue
			}
			env := envVon(z.name)
			for _, wert := range []string{"true", "5s", "pfad.yaml", "GEHEIM"} {
				leereAlle(t)
				t.Setenv(env, wert)
				if _, err := lese(basis(kommando, "")...); !meldetFremd(err, env, wert) {
					t.Errorf("%s mit %s=%s: erwartet %s mit dem Namen und ohne den Wert, erhalten %v", kommando, env, wert, model.CodeUsage, err)
				}
			}
			leereAlle(t)
			if _, err := lese(basis(kommando, "")...); err != nil {
				t.Errorf("%s mit leerer %s: %v", kommando, env, err)
			}
		}
	}
}

// Abdeckung: LH-FA-17/Boundary — die Umgebungsvariable einer Option, die das
// Kommando kennt und der Stand nicht liest (nichtGelesen), bleibt beim eigenen
// Kommando unbeachtet, auch mit ungültigem Wert; beim anderen Kommando ist sie
// PGR-E2001 (LH-FA-17.a).
func TestLeserUmgebungNichtGelesenerOptionen(t *testing.T) {
	for _, kommando := range leserKommandos() {
		for _, z := range tabelleSpec() {
			if !nichtGelesen()[z.name] {
				continue
			}
			env := envVon(z.name)
			leereAlle(t)
			t.Setenv(env, "ungültig")
			_, err := lese(basis(kommando, "")...)
			if z.kennt(kommando) && err != nil {
				t.Errorf("%s mit %s: erwartet unbeachtet, erhalten %v", kommando, env, err)
			}
			if !z.kennt(kommando) && !meldetFremd(err, env, "ungültig") {
				t.Errorf("%s mit %s: erwartet %s, erhalten %v", kommando, env, model.CodeUsage, err)
			}
		}
	}
}

// Abdeckung: LH-FA-17/Negative — die Zuordnung der Optionen zu den Kommandos
// stimmt mit der Spalte Kommando überein: Jede Option, die ein Kommando am
// allgemeinen Leser liest, steht in der Tabelle mit diesem Kommando, jede
// Option der Tabelle mit diesem Kommando ist gelesen oder eine der drei noch
// nicht gelesenen, und --config ist keine Zeile (LH-FA-17.a).
func TestLeserZuordnungGleichTabelle(t *testing.T) {
	for _, kommando := range leserKommandos() {
		gelesen := map[string]bool{}
		for _, o := range cli.Optionen(kommando) {
			gelesen[o.Name] = true
			ok := false
			for _, z := range tabelleSpec() {
				if z.name == o.Name && z.kennt(kommando) {
					ok = true
				}
			}
			if !ok {
				t.Errorf("%s liest --%s, die Tabelle ordnet sie dem Kommando nicht zu", kommando, o.Name)
			}
		}
		for _, z := range tabelleSpec() {
			if z.kennt(kommando) && !gelesen[z.name] && !nichtGelesen()[z.name] {
				t.Errorf("%s: --%s steht in der Tabelle, wird nicht gelesen und ist nicht als noch nicht gelesen benannt", kommando, z.name)
			}
		}
	}
}

// Abdeckung: LH-FA-17/Negative — das Kommando prüft die Variablen in einem
// Durchlauf über die Optionen in der Reihenfolge der Tabelle und meldet nur den
// ersten Fehler: zwei fremde Variablen melden die in der Tabelle erste; eine
// fremde Variable vor und nach einer ungültigen eigenen melden die in der
// Tabelle erste, gleich ob sie fremd oder eigen ist; die Kommandozeile geht vor, die fremde Variable
// (PGR-E2001) geht der Datei vor (LH-FA-17.a *Fehler*).
func TestLeserUmgebungReihenfolge(t *testing.T) {
	melde := func(kommando string, env map[string]string, args ...string) error {
		leereAlle(t)
		for k, v := range env {
			t.Setenv(k, v)
		}
		_, err := lese(append(basis(kommando, ""), args...)...)
		return err
	}
	// zwei fremde Variablen: bei play steht LISTEN vor SHUTDOWN_TIMEOUT.
	err := melde("play", map[string]string{"PGWIRE_RECORDER_SHUTDOWN_TIMEOUT": "1", "PGWIRE_RECORDER_LISTEN": "1"})
	if !meldetFremd(err, "PGWIRE_RECORDER_LISTEN", "") || strings.Contains(err.Error(), "SHUTDOWN_TIMEOUT") {
		t.Errorf("zwei fremde Variablen: erwartet nur LISTEN, erhalten %v", err)
	}
	// ungültige eigene Variable vor einer fremden: bei play steht CONTINUE_ON_ERROR vor TLS_CERT.
	err = melde("play", map[string]string{"PGWIRE_RECORDER_CONTINUE_ON_ERROR": "GEHEIM", "PGWIRE_RECORDER_TLS_CERT": "x"})
	if !meldetFremd(err, "PGWIRE_RECORDER_CONTINUE_ON_ERROR", "GEHEIM") || strings.Contains(err.Error(), "TLS_CERT") {
		t.Errorf("ungültige eigene vor fremder: erwartet nur CONTINUE_ON_ERROR, erhalten %v", err)
	}
	// fremde Variable vor einer ungültigen eigenen: bei play steht SHUTDOWN_TIMEOUT vor CONTINUE_ON_ERROR.
	err = melde("play", map[string]string{"PGWIRE_RECORDER_CONTINUE_ON_ERROR": "GEHEIM", "PGWIRE_RECORDER_SHUTDOWN_TIMEOUT": "5s"})
	if !meldetFremd(err, "PGWIRE_RECORDER_SHUTDOWN_TIMEOUT", "5s") || strings.Contains(err.Error(), "CONTINUE_ON_ERROR") {
		t.Errorf("fremde vor ungültiger eigener: erwartet nur SHUTDOWN_TIMEOUT, erhalten %v", err)
	}
	// replay: FORCE (fremd) steht in der Tabelle vor FAIL_ON_UNCONSUMED (eigen, ungültig).
	err = melde("replay", map[string]string{"PGWIRE_RECORDER_FAIL_ON_UNCONSUMED": "GEHEIM", "PGWIRE_RECORDER_FORCE": "true"})
	if !meldetFremd(err, "PGWIRE_RECORDER_FORCE", "true") || strings.Contains(err.Error(), "FAIL_ON_UNCONSUMED") {
		t.Errorf("replay: erwartet nur FORCE, erhalten %v", err)
	}
	// record: FORMAT (eigen, ungültig) steht vor INPUT (fremd).
	err = melde("record", map[string]string{"PGWIRE_RECORDER_FORMAT": "GEHEIM", "PGWIRE_RECORDER_INPUT": "r.yaml"})
	if !meldetFremd(err, "PGWIRE_RECORDER_FORMAT", "GEHEIM") || strings.Contains(err.Error(), "INPUT") {
		t.Errorf("record: erwartet nur FORMAT, erhalten %v", err)
	}
	// Kommandozeile vor Umgebung: eine ungültige Option geht der fremden Variable vor.
	err = melde("play", map[string]string{"PGWIRE_RECORDER_LISTEN": "x"}, "--user=")
	if !istUsage(err) || strings.Contains(err.Error(), "PGWIRE_RECORDER_LISTEN") {
		t.Errorf("Kommandozeile vor fremder Variable: %v", err)
	}
	// fremde Variable vor der Datei: die ungültige Datei wird nicht gelesen.
	datei := schreibe(t, "kein: [gültiges yaml\n")
	err = melde("play", map[string]string{"PGWIRE_RECORDER_LISTEN": "x"}, "--config="+datei)
	if !meldetFremd(err, "PGWIRE_RECORDER_LISTEN", "") || istDatei(err) {
		t.Errorf("fremde Variable vor ungültiger Datei: erwartet %s, erhalten %v", model.CodeUsage, err)
	}
}

// Abdeckung: LH-FA-17/Boundary — die Variable wird unabhängig von der Priorität
// geprüft: Setzt die Kommandozeile dieselbe eigene Option, ist eine ungültige
// Variable dennoch PGR-E2001, und die Option eines anderen Kommandos auf der
// Kommandozeile ist eine unbekannte Option vor der Umgebung (LH-FA-17.a).
func TestLeserUmgebungNebenKommandozeile(t *testing.T) {
	leereAlle(t)
	t.Setenv("PGWIRE_RECORDER_KEEP_TIMING", "GEHEIM")
	if _, err := lese(append(basis("play", ""), "--keep-timing=false")...); !meldetFremd(err, "PGWIRE_RECORDER_KEEP_TIMING", "GEHEIM") {
		t.Errorf("ungültige eigene Variable neben der Option: %v", err)
	}
	leereAlle(t)
	t.Setenv("PGWIRE_RECORDER_LISTEN", "x")
	_, err := lese(append(basis("play", ""), "--listen=x")...)
	if !istUsage(err) || !strings.Contains(err.Error(), "flag provided but not defined") || strings.Contains(err.Error(), "PGWIRE_RECORDER_LISTEN") {
		t.Errorf("Option eines anderen Kommandos auf der Kommandozeile: %v", err)
	}
}

// Abdeckung: LH-FA-17/Boundary, LH-FA-01/Boundary — config show prüft von den
// Variablen nur PGWIRE_RECORDER_CONFIG: Variablen jeder Option, gültig oder
// nicht, ändern es nicht; die Hilfe von record, replay, play und config show
// liest und prüft nichts (LH-FA-17.a, LH-FA-01.a).
func TestLeserUmgebungConfigShowUndHilfe(t *testing.T) {
	leereAlle(t)
	for _, z := range tabelleSpec() {
		t.Setenv(envVon(z.name), "ungültig")
	}
	var out bytes.Buffer
	if cmd, err := cli.Parse([]string{"config", "show"}, &out); err != nil || cmd.Name != "config show" {
		t.Errorf("config show mit Variablen jeder Option: %#v, %v", cmd, err)
	}
	t.Setenv(cli.EnvConfig, "/gibt/es/nicht.yaml")
	if _, err := cli.Parse([]string{"config", "show"}, &out); !istDatei(err) {
		t.Errorf("config show mit PGWIRE_RECORDER_CONFIG auf eine fehlende Datei: erwartet %s, erhalten %v", model.CodeConfigFile, err)
	}
	t.Setenv(cli.EnvConfig, "")
	for _, args := range [][]string{{"record", "--help"}, {"replay", "-h"}, {"play", "--help"}, {"config", "show", "--help"}} {
		hilfe(t, args...)
	}
}
