//go:build integration

package integration_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// versionsZeile liefert die Zeile, die das Kommando version des gebauten Binaries ausgibt, und
// bricht ab, wenn sie nicht genau eine Zeile "pgwire-recorder <Version>" auf stdout ist.
func versionsZeile(t *testing.T) string {
	t.Helper()
	stdout, stderr, code := starteBis(t, t.TempDir(), "version")
	if code != 0 || stderr != "" || !regexp.MustCompile(`^pgwire-recorder \S+\n$`).MatchString(stdout) {
		t.Fatalf("version: Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	return stdout
}

// Abdeckung: LH-FA-01/Happy, LH-FA-14/Boundary — version, --version und -v als einziges Argument
// schreiben auf stdout dieselbe Zeile, nichts auf stderr, und enden mit Exit-Code 0; sie lesen
// weder Umgebungsvariablen noch Konfigurationsdatei: Jede Umgebungsvariable PGWIRE_RECORDER_*, die
// die allgemeine Hilfe nennt, ist ungültig gesetzt, PGWIRE_RECORDER_CONFIG nennt eine fehlende
// Datei, und die Standarddatei im Arbeitsverzeichnis ist ungültig (LH-FA-01.a *Version*,
// LH-FA-14.a).
func TestE2EVersionsAngabe(t *testing.T) {
	zeile := versionsZeile(t)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("log_level: INFO\nunbekannt: [\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	hilfe, _, _ := starteBis(t, dir, "--help")
	namen := regexp.MustCompile(`PGWIRE_RECORDER_[A-Z_]+`).FindAllString(hilfe, -1)
	if len(namen) < 20 {
		t.Fatalf("die allgemeine Hilfe nennt nur %d Umgebungsvariablen: %q", len(namen), namen)
	}
	for _, name := range namen {
		t.Setenv(name, "%ungültig%")
	}
	t.Setenv("PGWIRE_RECORDER_PASSWORD", "%ungültig%")
	t.Setenv("PGWIRE_RECORDER_GIBT_ES_NICHT", "%ungültig%")
	t.Setenv("PGWIRE_RECORDER_CONFIG", filepath.Join(dir, "fehlt.yaml"))
	for _, args := range [][]string{{"version"}, {"--version"}, {"-v"}} {
		stdout, stderr, code := starteBis(t, dir, args...)
		if code != 0 || stdout != zeile || stderr != "" {
			t.Errorf("%q: Exit-Code %d, stdout %q (erwartet %q), stderr %q", args, code, stdout, zeile, stderr)
		}
	}
}

// Abdeckung: LH-FA-01/Negative — andere Schreibweisen (-version, --v, -V, --Version), ein Argument
// mit = (--version=x, -v=x) und --version oder -v mit einem weiteren Argument, davor oder danach
// (--version record, -v record, record --version, record -v, version --version, -- --version),
// sind PGR-E2001 mit Exit-Code 2, nichts auf stdout; als einziges Argument sind die ersten sechs
// ein unbekanntes Kommando (LH-FA-01.a *Version*).
func TestE2EKeineVersionsAngabe(t *testing.T) {
	dir := t.TempDir()
	for _, args := range [][]string{
		{"-version"}, {"--v"}, {"-V"}, {"--Version"}, {"--version=x"}, {"-v=x"},
		{"--version", "record"}, {"-v", "record"}, {"record", "--version"}, {"record", "-v"},
		{"version", "--version"}, {"--", "--version"},
	} {
		stdout, stderr, code := starteBis(t, dir, args...)
		einzeln := len(args) == 1
		if code != 2 || stdout != "" || !strings.Contains(stderr, "[PGR-E2001]") || (einzeln && !strings.Contains(stderr, "unbekanntes Kommando")) {
			t.Errorf("%q: Exit-Code %d, stdout %q, stderr %q", args, code, stdout, stderr)
		}
	}
}

// Abdeckung: LH-FA-01/Negative — bei record, replay und play mit allen Pflichtoptionen sind
// --version und -v eine unbekannte Option (PGR-E2001, Exit-Code 2, nichts auf stdout)
// (LH-FA-01.a *Version*).
func TestE2EVersionNachKommando(t *testing.T) {
	dir := t.TempDir()
	datei := filepath.Join(dir, "r.yaml")
	for _, basis := range [][]string{
		{"record", "--listen", "127.0.0.1:0", "--upstream", "127.0.0.1:1", "--output", datei},
		{"replay", "--listen", "127.0.0.1:0", "--input", datei},
		{"play", "--upstream", "127.0.0.1:1", "--input", datei},
	} {
		for _, option := range []string{"--version", "-v"} {
			args := append(append([]string{}, basis...), option)
			stdout, stderr, code := starteBis(t, dir, args...)
			if code != 2 || stdout != "" || !strings.Contains(stderr, "[PGR-E2001]") || !strings.Contains(stderr, "not defined: -"+strings.TrimLeft(option, "-")) {
				t.Errorf("%s %s: Exit-Code %d, stdout %q, stderr %q", basis[0], option, code, stdout, stderr)
			}
		}
	}
}

// Abdeckung: LH-FA-01/Boundary — neben --version oder -v gilt Hilfe vor Prüfung: --version --help
// und -v -h geben die allgemeine Hilfe aus, die --version und -v nennt, mit Exit-Code 0 und
// nichts auf stderr; die Hilfe von record, replay, play, config und version nennt beide nicht
// (LH-FA-01.a *Version*).
func TestE2EVersionHilfe(t *testing.T) {
	dir := t.TempDir()
	for _, args := range [][]string{{"--version", "--help"}, {"-v", "-h"}, {"--help", "-v"}} {
		stdout, stderr, code := starteBis(t, dir, args...)
		if code != 0 || stderr != "" || !strings.Contains(stdout, "Kommandos:") || !strings.Contains(stdout, "\n  --version, -v\n") {
			t.Errorf("%q: Exit-Code %d, stdout %q, stderr %q", args, code, stdout, stderr)
		}
	}
	nennt := regexp.MustCompile(`--version|(^|[\s,])-v([\s,]|$)`)
	for _, kommando := range []string{"record", "replay", "play", "config", "version"} {
		stdout, stderr, code := starteBis(t, dir, kommando, "--help")
		if code != 0 || stderr != "" || stdout == "" || nennt.MatchString(stdout) {
			t.Errorf("%s --help: Exit-Code %d, stdout %q, stderr %q", kommando, code, stdout, stderr)
		}
	}
}
