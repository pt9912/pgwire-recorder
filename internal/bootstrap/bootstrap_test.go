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
