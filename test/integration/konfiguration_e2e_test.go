//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"errors"
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Getenv("PGR_BINARY"), args...)
	cmd.Dir = dir
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
