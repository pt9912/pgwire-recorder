package bootstrap

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// Fordert der Aufruf die Hilfe an, endet Run mit Exit-Code 0 und der Hilfe auf
// stdout, auch mit ungültiger PGWIRE_RECORDER_FAIL_ON_UNCONSUMED; stderr
// bleibt leer (LH-FA-01.a).
func TestRunHilfe(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "1")
	var stdout, stderr bytes.Buffer
	if code := Run(context.Background(), []string{"replay", "--help"}, "dev", &stdout, &stderr); code != 0 {
		t.Fatalf("Exit-Code %d, stderr %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Optionen von replay") || stderr.Len() > 0 {
		t.Fatalf("stdout %q, stderr %q", stdout.String(), stderr.String())
	}
}

// version gibt die Version aus und endet mit Exit-Code 0, auch mit gesetzter
// ungültiger PGWIRE_RECORDER_FAIL_ON_UNCONSUMED: Es liest keine
// Umgebungsvariable (LH-FA-01.a).
func TestRunVersion(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "1")
	var stdout, stderr bytes.Buffer
	if code := Run(context.Background(), []string{"version"}, "1.2.3", &stdout, &stderr); code != 0 || stdout.String() != "pgwire-recorder 1.2.3\n" || stderr.Len() > 0 {
		t.Fatalf("Exit-Code %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
}

// Das erste -- beendet die Optionen auch an der Stelle eines Werts: Mit
// --input -- fehlt der Wert (PGR-E2001, Exit-Code 2), und keine Datei wird
// gelesen; mit --input=-- ist -- der Wert, und erst das Laden scheitert
// (PGR-E3001, Exit-Code 3) (LH-FA-01.a).
func TestRunEndeDerOptionen(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	t.Chdir(t.TempDir())
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"replay", "--listen", "127.0.0.1:0", "--input", "--"}, "dev", &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "PGR-E2001") || strings.Contains(stderr.String(), "PGR-E3") {
		t.Fatalf("--input --: Exit-Code %d, stderr %q", code, stderr.String())
	}
	stderr.Reset()
	code = Run(context.Background(), []string{"replay", "--listen", "127.0.0.1:0", "--input=--"}, "dev", &stdout, &stderr)
	if code != 3 || !strings.Contains(stderr.String(), "PGR-E3001") {
		t.Fatalf("--input=--: Exit-Code %d, stderr %q", code, stderr.String())
	}
}
