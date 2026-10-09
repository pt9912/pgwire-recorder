//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// starteBis führt das Binary mit args in dir aus und wartet höchstens 15 s auf
// sein Ende; läuft die Frist ab, wird der Test rot und nennt das ausgebliebene
// Ende. Es liefert stdout, stderr und den Exit-Code.
func starteBis(t *testing.T, dir string, args ...string) (string, string, int) {
	t.Helper()
	return starteMitEingabe(t, dir, nil, args...)
}

// starteMitEingabe ist starteBis mit stdin aus eingabe; nil ist kein stdin.
func starteMitEingabe(t *testing.T, dir string, eingabe io.Reader, args ...string) (string, string, int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Getenv("PGR_BINARY"), args...)
	cmd.Dir = dir
	cmd.Stdin = eingabe
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("%q endet nicht binnen 15 s", args)
	}
	return stdout.String(), stderr.String(), exitCodeOf(err)
}

// Abdeckung: LH-FA-17/Happy — config show nennt im Verzeichnis mit
// .pgwire-recorder.yaml die Datei in der ersten Zeile und ihren Inhalt ohne
// Kommentare auf stdout, mit Platzhaltern unaufgelöst, nichts auf stderr, und
// endet mit Exit-Code 0; --config wählt eine andere Datei.
func TestE2EConfigShow(t *testing.T) {
	dir := t.TempDir()
	inhalt := "# Kommentar\nconnections:\n  lokal: \"postgresql://${PGR_E2E_NUTZER}@localhost/db\"\nreplay:\n  listen: 127.0.0.1:0\n"
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte(inhalt), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := starteBis(t, dir, "config", "show")
	want := ".pgwire-recorder.yaml\nconnections:\n  lokal: \"postgresql://${PGR_E2E_NUTZER}@localhost/db\"\nreplay:\n  listen: 127.0.0.1:0\n"
	if code != 0 || stdout != want || stderr != "" {
		t.Fatalf("Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	andere := filepath.Join(t.TempDir(), "andere.yaml")
	if err := os.WriteFile(andere, []byte("log_level: warn\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout, _, code = starteBis(t, dir, "config", "show", "--config", andere)
	if code != 0 || stdout != andere+"\nlog_level: warn\n" {
		t.Fatalf("--config: Exit-Code %d, stdout %q", code, stdout)
	}
}

// Abdeckung: LH-FA-17/Negative — eine ungültige .pgwire-recorder.yaml beendet
// replay mit PGR-E2004 und Exit-Code 2, ohne die Aufzeichnung zu lesen; die
// Meldung nennt den Schlüssel, nicht den Wert.
func TestE2EReplayKonfigurationsdateiUngueltig(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("replay:\n  shutdown_timeout: GEHEIM\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := starteBis(t, dir, "replay", "--listen", "127.0.0.1:0", "--input", "fehlt.yaml")
	if code != 2 || stdout != "" || !strings.Contains(stderr, "PGR-E2004") || !strings.Contains(stderr, "replay.shutdown_timeout") || strings.Contains(stderr, "GEHEIM") || strings.Contains(stderr, "PGR-E3") {
		t.Fatalf("Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
}

// Abdeckung: LH-FA-17/Negative, LH-FA-17/Boundary — config show --config
// /dev/stdin mit stdin aus einer offenen Pipe ohne Eingabe endet binnen der
// Frist mit PGR-E2004 und Exit-Code 2, statt auf Eingabe zu warten; mit stdin
// aus einer regulären Datei
// zeigt es deren Inhalt.
func TestE2EConfigShowStdin(t *testing.T) {
	dir := t.TempDir()
	lesen, schreiben, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = lesen.Close() }()
	defer func() { _ = schreiben.Close() }()
	stdout, stderr, code := starteMitEingabe(t, dir, lesen, "config", "show", "--config", "/dev/stdin")
	if code != 2 || stdout != "" || !strings.Contains(stderr, "PGR-E2004") || !strings.Contains(stderr, "nicht lesbar") {
		t.Fatalf("Pipe: Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	pfad := filepath.Join(dir, "eingabe.yaml")
	if err := os.WriteFile(pfad, []byte("log_level: warn\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	datei, err := os.Open(pfad)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = datei.Close() }()
	stdout, stderr, code = starteMitEingabe(t, dir, datei, "config", "show", "--config", "/dev/stdin")
	if code != 0 || stdout != "/dev/stdin\nlog_level: warn\n" || stderr != "" {
		t.Fatalf("reguläre Datei: Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
}
