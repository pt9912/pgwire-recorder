package cli_test

import (
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
)

// schalterQuellen prüft die Quellen einer Option, die auf der Kommandozeile ohne
// Wert true ist: ohne Quelle false; die Option ohne Wert, mit true, die
// Umgebungsvariable und der Schlüssel der Datei im Abschnitt setzen sie; die
// Option mit false auf der Kommandozeile geht der Umgebungsvariable vor.
func schalterQuellen(t *testing.T, kommando, option, env, schluessel string, wert func(cli.Command) bool) {
	t.Helper()
	leere(t, kommando)
	args := basis(kommando, "")
	mit := func(extra ...string) bool {
		t.Helper()
		cmd, err := lese(append(append([]string{}, args...), extra...)...)
		if err != nil {
			t.Fatalf("%s %v: %v", kommando, extra, err)
		}
		return wert(cmd)
	}
	if mit() {
		t.Errorf("%s ohne Option: gesetzt", kommando)
	}
	datei := "--config=" + schreibe(t, kommando+":\n  "+schluessel+": true\n")
	for _, extra := range []string{"--" + option, "--" + option + "=true", datei} {
		if !mit(extra) {
			t.Errorf("%s %s: nicht gesetzt", kommando, extra)
		}
	}
	t.Setenv(env, "true")
	if !mit() {
		t.Errorf("%s mit %s=true: nicht gesetzt", kommando, env)
	}
	if mit("--" + option + "=false") {
		t.Errorf("%s mit %s=true und --%s=false: gesetzt", kommando, env, option)
	}
}

// Abdeckung: LH-FA-12/Happy, LH-FA-17/Boundary — --record-empty-sessions, seine
// Umgebungsvariable und der Schlüssel record_empty_sessions im Abschnitt record:
// setzen RecordEmptySessions; ohne alle drei ist es false, und die Option auf der
// Kommandozeile geht der Umgebungsvariable vor (LH-FA-12.a, LH-FA-17.a).
func TestRecordEmptySessionsQuellen(t *testing.T) {
	schalterQuellen(t, "record", "record-empty-sessions", "PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS", "record_empty_sessions",
		func(c cli.Command) bool { return c.Record.RecordEmptySessions })
}

// Abdeckung: LH-FA-12/Happy, LH-FA-17/Boundary — --session-assignment, seine
// Umgebungsvariable und der Schlüssel session_assignment im Abschnitt replay:
// setzen SessionAssignment auf connection; ohne alle drei gilt first-request, und
// die Option auf der Kommandozeile geht der Umgebungsvariable vor (LH-FA-12.a,
// LH-FA-17.a).
func TestSessionAssignmentQuellen(t *testing.T) {
	leere(t, "replay")
	args := basis("replay", "")
	cmd, err := lese(args...)
	if err != nil || cmd.Replay.SessionAssignment != cli.ZuordnungErsteAnfrage {
		t.Fatalf("ohne Option: %q, %v", cmd.Replay.SessionAssignment, err)
	}
	datei := schreibe(t, "replay:\n  session_assignment: connection\n")
	for _, extra := range []string{"--session-assignment=connection", "--config=" + datei} {
		cmd, err := lese(append(append([]string{}, args...), extra)...)
		if err != nil || cmd.Replay.SessionAssignment != cli.ZuordnungVerbindung {
			t.Errorf("%s: %q, %v", extra, cmd.Replay.SessionAssignment, err)
		}
	}
	// Die Option ist kein Schalter: ohne Wert ist sie ungültig.
	if _, err := lese(append(append([]string{}, args...), "--session-assignment")...); !istUsage(err) {
		t.Errorf("--session-assignment ohne Wert: %v", err)
	}
	t.Setenv("PGWIRE_RECORDER_SESSION_ASSIGNMENT", "connection")
	if cmd, err := lese(args...); err != nil || cmd.Replay.SessionAssignment != cli.ZuordnungVerbindung {
		t.Errorf("Umgebungsvariable: %q, %v", cmd.Replay.SessionAssignment, err)
	}
	if cmd, err := lese(append(append([]string{}, args...), "--session-assignment=first-request")...); err != nil || cmd.Replay.SessionAssignment != cli.ZuordnungErsteAnfrage {
		t.Errorf("Option first-request vor Umgebungsvariable connection: %q, %v", cmd.Replay.SessionAssignment, err)
	}
}

// Abdeckung: LH-FA-12/Negative, LH-FA-17/Negative — --session-assignment nimmt
// genau first-request und connection, klein geschrieben; jeder andere Wert, auch
// Connection, CONNECTION, first_request, ein Alias, einer mit Leerraum und der
// leere, ist auf der Kommandozeile und in der Umgebungsvariable PGR-E2001 und in
// der Datei PGR-E2004; die Meldung zur Umgebungsvariable und zur Datei nennt die
// Stelle, nicht den Wert (LH-FA-12.a *Wertemenge*, LH-FA-17.a).
func TestSessionAssignmentWertemenge(t *testing.T) {
	leere(t, "replay")
	args := basis("replay", "")
	for _, wert := range []string{"Connection", "CONNECTION", "First-Request", "first_request", "firstrequest", "first", "conn", "per-connection", "connection ", " connection", "1", "true"} {
		geheim := wert + "GEHEIM"
		if _, err := lese(append(append([]string{}, args...), "--session-assignment="+wert)...); !istUsage(err) || !strings.Contains(err.Error(), "session-assignment") {
			t.Errorf("Kommandozeile %q: %v", wert, err)
		}
		t.Setenv("PGWIRE_RECORDER_SESSION_ASSIGNMENT", wert)
		if _, err := lese(args...); !istUsage(err) || !strings.Contains(err.Error(), "PGWIRE_RECORDER_SESSION_ASSIGNMENT") {
			t.Errorf("Umgebungsvariable %q: %v", wert, err)
		}
		t.Setenv("PGWIRE_RECORDER_SESSION_ASSIGNMENT", geheim)
		if _, err := lese(args...); err == nil || strings.Contains(err.Error(), "GEHEIM") {
			t.Errorf("Umgebungsvariable %q: %v", geheim, err)
		}
		t.Setenv("PGWIRE_RECORDER_SESSION_ASSIGNMENT", "")
		datei := schreibe(t, "replay:\n  session_assignment: \""+geheim+"\"\n")
		if _, err := lese(append(append([]string{}, args...), "--config="+datei)...); !istDatei(err) || !strings.Contains(err.Error(), "replay.session_assignment") || strings.Contains(err.Error(), "GEHEIM") {
			t.Errorf("Datei %q: %v", geheim, err)
		}
		datei = schreibe(t, "replay:\n  session_assignment: \""+wert+"\"\n")
		if _, err := lese(append(append([]string{}, args...), "--config="+datei)...); !istDatei(err) {
			t.Errorf("Datei %q: %v", wert, err)
		}
	}
	if _, err := lese(append(append([]string{}, args...), "--session-assignment=")...); !istUsage(err) {
		t.Errorf("leerer Wert auf der Kommandozeile: %v", err)
	}
	datei := schreibe(t, "replay:\n  session_assignment: \"\"\n")
	if _, err := lese(append(append([]string{}, args...), "--config="+datei)...); !istDatei(err) {
		t.Errorf("leerer Wert in der Datei: %v", err)
	}
}

// Abdeckung: LH-FA-12/Negative, LH-FA-17/Negative — --session-assignment gibt es
// nur bei replay und --record-empty-sessions nur bei record: Bei jedem anderen
// Kommando ist die Option auf der Kommandozeile PGR-E2001 (LH-FA-12.a
// *Wertemenge*, LH-FA-17.a).
func TestSessionsOptionenNurBeiEigenemKommando(t *testing.T) {
	for _, f := range []struct{ kommando, option string }{
		{"record", "--session-assignment=connection"},
		{"play", "--session-assignment=connection"},
		{"replay", "--record-empty-sessions"},
		{"play", "--record-empty-sessions"},
	} {
		leere(t, f.kommando)
		if _, err := lese(append(basis(f.kommando, ""), f.option)...); !istUsage(err) {
			t.Errorf("%s %s: %v", f.kommando, f.option, err)
		}
	}
}

// Abdeckung: LH-FA-01/Happy, LH-FA-12/Happy — die Hilfe von record nennt
// --record-empty-sessions und PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS, die von
// replay und play nicht (die Hilfe von replay nennt die Option im Text zu
// --session-assignment, aber keine eigene Zeile und keine Variable); die Hilfe
// von replay nennt --session-assignment mit beiden Werten und
// PGWIRE_RECORDER_SESSION_ASSIGNMENT, die von record und play nicht (LH-FA-12.a,
// LH-FA-01.a).
func TestSessionsHilfe(t *testing.T) {
	for kommando, erwartet := range map[string]bool{"record": true, "replay": false, "play": false} {
		text := hilfe(t, kommando, "--help")
		for _, teil := range []string{"\n  --record-empty-sessions", "PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS"} {
			if strings.Contains(text, teil) != erwartet {
				t.Errorf("%s: Hilfe nennt %q: %v, erwartet %v", kommando, teil, !erwartet, erwartet)
			}
		}
	}
	for kommando, erwartet := range map[string]bool{"record": false, "replay": true, "play": false} {
		text := hilfe(t, kommando, "--help")
		teile := []string{"--session-assignment", "PGWIRE_RECORDER_SESSION_ASSIGNMENT", "first-request", "connection"}
		for _, teil := range teile {
			if kommando != "replay" && (teil == "first-request" || teil == "connection") {
				// first-request nennt auch --timing-reference bei play.
				continue
			}
			if strings.Contains(text, teil) != erwartet {
				t.Errorf("%s: Hilfe nennt %q: %v, erwartet %v", kommando, teil, !erwartet, erwartet)
			}
		}
	}
}
