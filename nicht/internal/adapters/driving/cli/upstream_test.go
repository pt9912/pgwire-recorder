package cli_test

import (
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
)

// envUpstream ist die Umgebungsvariable von --upstream.
const envUpstream = "PGWIRE_RECORDER_UPSTREAM"

// upstreamUngueltig sind Werte, die weder Name einer Verbindung noch host:port
// sind (LH-FA-17.a, A7, F-521): ohne :, leerer Host, kein Port, Port außerhalb
// 1 bis 65535, Host mit @, Leerraum, /, %, [ oder Steuerzeichen, in Klammern
// keine IPv6-Adresse, IPv6 ohne Klammern, ein zweites :.
func upstreamUngueltig() []string {
	return []string{
		"GEHEIM", ":5432", "GEHEIM:", "GEHEIM:0", "GEHEIM:65536", "GEHEIM@h:5", "GEHEIM h:5", "a/GEHEIM:5",
		"h%41GEHEIM:1", "a[GEHEIM:1", "h\x01GEHEIM:5", "[GEHEIM]:5", "[1.2.3.4]:5", "::1:5432", "[::1]", "GEHEIM:1:2",
	}
}

// recordUpstream liest record mit --listen und --output und den übrigen
// Argumenten.
func recordUpstream(args ...string) (string, error) {
	cmd, err := lese(append([]string{"record", "--listen=x", "--output=r.yaml"}, args...)...)
	return cmd.Record.Upstream, err
}

// Abdeckung: LH-FA-17/Negative — ein Wert von --upstream, der weder Name einer
// Verbindung der gewählten Datei ist noch die Form host:port hat, ist
// PGR-E2001; die Meldung zur Umgebungsvariable nennt den Wert nicht. Geprüft
// wird jeder gesetzte Wert nach der Zusammenführung, die Kommandozeile vor der
// Umgebungsvariable, auch der Wert, der nicht gilt (LH-FA-17.a, A7, A8, F-521).
func TestUpstreamUngueltig(t *testing.T) {
	leere(t, "record")
	ausCli := "Konfiguration [PGR-E2001]: --upstream: weder Name einer Verbindung der Konfigurationsdatei noch host:port"
	ausEnv := "Konfiguration [PGR-E2001]: Umgebungsvariable " + envUpstream + ": weder Name einer Verbindung der Konfigurationsdatei noch host:port"
	for _, w := range upstreamUngueltig() {
		t.Setenv(envUpstream, "")
		if _, err := recordUpstream("--upstream=" + w); err == nil || err.Error() != ausCli || !istUsage(err) {
			t.Errorf("--upstream=%q: %v, erwartet %q", w, err, ausCli)
		}
		t.Setenv(envUpstream, w)
		if _, err := recordUpstream(); err == nil || err.Error() != ausEnv || !istUsage(err) {
			t.Errorf("%s=%q: %v, erwartet %q", envUpstream, w, err, ausEnv)
		}
		if _, err := recordUpstream("--upstream=h:1"); err == nil || err.Error() != ausEnv {
			t.Errorf("%s=%q neben gültigem --upstream: %v, erwartet %q", envUpstream, w, err, ausEnv)
		}
		if _, err := recordUpstream("--upstream=ohne-port"); err == nil || err.Error() != ausCli {
			t.Errorf("--upstream=ohne-port neben %s=%q: %v, erwartet zuerst die Kommandozeile", envUpstream, w, err)
		}
	}
}

// Abdeckung: LH-FA-17/Happy, LH-FA-17/Boundary — --upstream und
// PGWIRE_RECORDER_UPSTREAM nehmen host:port wie geschrieben: IPv6 in eckigen
// Klammern, auch mit Zone hinter %, Port mit führenden Nullen; die
// Kommandozeile geht der Umgebungsvariable vor (LH-FA-17.a, A7).
func TestUpstreamHostPort(t *testing.T) {
	leere(t, "record")
	for _, w := range []string{"h:5432", "h:05432", "h:1", "h:65535", "[::1]:5432", "[fe80::1%eth0]:5", "127.0.0.1:5432", "h.example:5"} {
		t.Setenv(envUpstream, "")
		if got, err := recordUpstream("--upstream=" + w); err != nil || got != w {
			t.Errorf("--upstream=%s: %q, %v", w, got, err)
		}
		t.Setenv(envUpstream, w)
		if got, err := recordUpstream(); err != nil || got != w {
			t.Errorf("%s=%s: %q, %v", envUpstream, w, got, err)
		}
		if got, err := recordUpstream("--upstream=k:2"); err != nil || got != "k:2" {
			t.Errorf("%s=%s neben --upstream=k:2: %q, %v", envUpstream, w, got, err)
		}
	}
}

// Abdeckung: LH-FA-17/Boundary, LH-FA-17/Negative — Namen für --upstream und
// PGWIRE_RECORDER_UPSTREAM kommen nur aus der gewählten Datei: Ohne Datei ist
// ein Name PGR-E2001, mit einer Datei, die ihn führt, gültig, nur bei genauer
// Übereinstimmung (LH-FA-17.a, U1).
func TestUpstreamNameNurAusDatei(t *testing.T) {
	leere(t, "record")
	if _, err := recordUpstream("--upstream=staging"); !istUsage(err) || !strings.Contains(err.Error(), "--upstream") {
		t.Errorf("ohne Datei --upstream=staging: %v", err)
	}
	pfad := schreibe(t, "connections:\n  staging: postgresql://h/db\n")
	if _, err := recordUpstream("--upstream=staging", "--config="+pfad); err != nil {
		t.Errorf("mit Datei --upstream=staging: %v", err)
	}
	t.Setenv(envUpstream, "staging")
	if _, err := recordUpstream("--config=" + pfad); err != nil {
		t.Errorf("mit Datei %s=staging: %v", envUpstream, err)
	}
	if _, err := recordUpstream(); !istUsage(err) || !strings.Contains(err.Error(), envUpstream) {
		t.Errorf("ohne Datei %s=staging: %v", envUpstream, err)
	}
	t.Setenv(envUpstream, "")
	if _, err := recordUpstream("--upstream=Staging", "--config="+pfad); !istUsage(err) || !strings.Contains(err.Error(), "--upstream") {
		t.Errorf("--upstream=Staging: %v", err)
	}
}

// Abdeckung: LH-FA-17/Negative — die Prüfung von --upstream steht nach der
// Zusammenführung: nach einer ungültigen Umgebungsvariable einer anderen
// Option, nach der Datei und nach den Pflichtoptionen (LH-FA-17.a *Fehler*).
func TestUpstreamReihenfolge(t *testing.T) {
	leere(t, "record")
	t.Setenv(cli.EnvShutdownTimeout, "x")
	if _, err := recordUpstream("--upstream=GEHEIM"); !istUsage(err) || !strings.Contains(err.Error(), cli.EnvShutdownTimeout) {
		t.Errorf("Umgebung vor --upstream: %v", err)
	}
	t.Setenv(cli.EnvShutdownTimeout, "")
	if _, err := recordUpstream("--upstream=GEHEIM", "--config="+schreibe(t, "bogus: 1\n")); !istDatei(err) {
		t.Errorf("Datei vor --upstream: %v", err)
	}
	if _, err := lese("record", "--listen=x", "--upstream=GEHEIM"); !istUsage(err) || !strings.Contains(err.Error(), "Pflichtoption --output") {
		t.Errorf("Pflichtoption vor --upstream: %v", err)
	}
}

// Die Hilfe von record nennt bei --upstream host:port und den Namen einer
// Verbindung, wie die Optionstabelle in LH-FA-17.a (U7).
func TestUpstreamHilfe(t *testing.T) {
	text := strings.Join(strings.Fields(hilfe(t, "record", "--help")), " ")
	if !strings.Contains(text, "--upstream PostgreSQL-Server als host:port oder Name einer Verbindung der Konfigurationsdatei (Pflicht)") {
		t.Errorf("Hilfe von --upstream:\n%s", text)
	}
	if !strings.Contains(strings.Join(strings.Fields(hilfe(t, "--help")), " "), "host:port oder Name einer Verbindung") {
		t.Errorf("globale Hilfe ohne den Namen einer Verbindung bei --upstream")
	}
}
