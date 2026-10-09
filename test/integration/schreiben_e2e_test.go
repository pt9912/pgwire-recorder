//go:build integration

package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// zielMitAlt legt in einem neuen Verzeichnis die Zieldatei rec.yaml mit dem
// Inhalt "alt" an und liefert Verzeichnis und Pfad.
func zielMitAlt(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	output := filepath.Join(dir, "rec.yaml")
	if err := os.WriteFile(output, []byte("alt"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir, output
}

// forceQuelle setzt den Wert v von --force aus der Quelle quelle für den
// Aufruf in dir und liefert die Argumente, die dafür auf die Kommandozeile
// gehören.
func forceQuelle(t *testing.T, dir, quelle, v string) []string {
	t.Helper()
	switch quelle {
	case "Kommandozeile":
		return []string{"--force=" + v}
	case "Umgebung":
		t.Setenv("PGWIRE_RECORDER_FORCE", v)
	case "Datei":
		if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("record:\n  force: "+v+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return nil
}

// verzeichnisEnthaelt prüft, dass in dir genau die Einträge namen liegen.
func verzeichnisEnthaelt(t *testing.T, dir string, namen ...string) {
	t.Helper()
	es, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range es {
		got = append(got, e.Name())
	}
	if strings.Join(got, ",") != strings.Join(namen, ",") {
		t.Fatalf("im Verzeichnis liegt %v, erwartet %v", got, namen)
	}
}

// Abdeckung: LH-FA-08/Boundary, LH-FA-17/Happy — eine vorhandene --output-Datei
// lehnt record ohne --force mit Exit-Code 2 (PGR-E2002) ab und lässt sie
// unverändert, ob --force fehlt oder false aus der Kommandozeile, der Umgebung
// oder der Konfigurationsdatei kommt; mit true aus jeder dieser Quellen ersetzt
// record sie, und der Lauf endet mit Exit-Code 0.
func TestE2ERecordVorhandeneZieldatei(t *testing.T) {
	for _, quelle := range []string{"keine", "Kommandozeile", "Umgebung", "Datei"} {
		t.Run(quelle+"/false", func(t *testing.T) {
			dir, output := zielMitAlt(t)
			args := append([]string{"record", "--listen", freieAdresse(t), "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output}, forceQuelle(t, dir, quelle, "false")...)
			stdout, stderr, code := starteBis(t, dir, args...)
			if code != 2 || !strings.HasPrefix(stderr, "Konfiguration [PGR-E2002]: ") || stdout != "" {
				t.Fatalf("Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
			}
			if text := lies(t, output); text != "alt" {
				t.Fatalf("Zieldatei verändert: %q", text)
			}
		})
	}
	for _, quelle := range []string{"Kommandozeile", "Umgebung", "Datei"} {
		t.Run(quelle+"/true", func(t *testing.T) {
			dir, output := zielMitAlt(t)
			args := append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output}, forceQuelle(t, dir, quelle, "true")...)
			rec := startProzessIn(t, dir, "record", args...)
			rec.stop(t, 0)
			if text := lies(t, output); !strings.Contains(text, "format: pgwire-recorder") || !strings.Contains(text, "sessions: []") {
				t.Fatalf("Zieldatei nicht ersetzt:\n%s", text)
			}
		})
	}
}

// Abdeckung: LH-FA-07/Negative — ist --output ein Verzeichnis, endet record beim
// Start mit Exit-Code 3 (PGR-E3001), ohne und mit --force; das Verzeichnis
// bleibt leer.
func TestE2ERecordZielIstVerzeichnis(t *testing.T) {
	for _, extra := range [][]string{nil, {"--force"}} {
		ziel := filepath.Join(t.TempDir(), "rec.yaml")
		if err := os.Mkdir(ziel, 0o755); err != nil {
			t.Fatal(err)
		}
		args := append([]string{"record", "--listen", freieAdresse(t), "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", ziel}, extra...)
		stdout, stderr, code := starteBis(t, t.TempDir(), args...)
		if code != 3 || !strings.HasPrefix(stderr, "Recording [PGR-E3001]: ") || stdout != "" {
			t.Fatalf("%v: Exit-Code %d, stdout %q, stderr %q", extra, code, stdout, stderr)
		}
		verzeichnisEnthaelt(t, ziel)
	}
}

// Abdeckung: LH-FA-07/Boundary — mit --force über einer vorhandenen Datei steht
// nach dem Zwangsende des Herunterfahrens unter --output eine vollständige,
// ladbare Aufzeichnung, und im Verzeichnis liegt keine temporäre Datei.
func TestE2ERecordZwangsendeSchreibtVollstaendig(t *testing.T) {
	dir, output := zielMitAlt(t)
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--force", "--shutdown-timeout", "1s")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, _ := laufendeInteraktion(ctx, t, rec.listen)
	defer conn.Conn().Close()

	rec.signal(t, syscall.SIGTERM)
	rec.warteEnde(t, 15*time.Second, "Recorder endet nicht nach Ablauf der Frist von 1s")
	rec.pruefeExit(t, 4)

	verzeichnisEnthaelt(t, dir, "rec.yaml")
	nurErsteInteraktion(t, output)
}

// Abdeckung: LH-FA-07/Negative, LH-FA-13/Negative — steht nach dem Start an
// --output ein Verzeichnis, scheitert das Schreiben nach dem Zwangsende des
// Herunterfahrens: der Lauf endet mit Exit-Code 3 (PGR-E3001), das Verzeichnis
// bleibt, und neben ihm liegt keine temporäre Datei.
func TestE2ERecordSchreibfehlerNachZwangsende(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--shutdown-timeout", "1s")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, _ := laufendeInteraktion(ctx, t, rec.listen)
	defer conn.Conn().Close()
	if err := os.Mkdir(output, 0o755); err != nil {
		t.Fatal(err)
	}

	rec.signal(t, syscall.SIGTERM)
	rec.warteEnde(t, 15*time.Second, "Recorder endet nicht nach Ablauf der Frist von 1s")
	rec.pruefeExit(t, 3)

	if !strings.Contains(rec.stderr.String(), "Recording [PGR-E3001]: ") {
		t.Fatalf("PGR-E3001 fehlt in der Ausgabe:\n%s", rec.stderr.String())
	}
	verzeichnisEnthaelt(t, dir, "rec.yaml")
	if info, err := os.Stat(output); err != nil || !info.IsDir() {
		t.Fatalf("Verzeichnis an --output: %v %v", info, err)
	}
}
