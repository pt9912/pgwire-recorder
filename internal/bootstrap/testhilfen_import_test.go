package bootstrap_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// modulWurzel liefert das Verzeichnis mit der go.mod über dem Verzeichnis des
// Tests.
func modulWurzel(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		eltern := filepath.Dir(dir)
		if eltern == dir {
			t.Fatal("keine go.mod über dem Verzeichnis des Tests")
		}
		dir = eltern
	}
}

// TestBinaryOhneTesthilfen prüft, dass keines der Pakete unter ./cmd/... die
// Testhilfen tlsproxy und testpki oder das Paket testing, auch nicht mittelbar,
// importiert: `go list -deps ./cmd/...` läuft ohne Netz (GOPROXY=off) und die
// Menge der Pakete wird gegen die Verbotsliste gehalten.
func TestBinaryOhneTesthilfen(t *testing.T) {
	wurzel := modulWurzel(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "list", "-deps", "./cmd/...")
	cmd.Dir = wurzel
	cmd.Env = append(os.Environ(), "GOPROXY=off", "GOFLAGS=-mod=readonly")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("go list -deps ./cmd/...: %v\n%s", err, stderr.String())
	}
	const modul = "github.com/pt9912/pgwire-recorder"
	verboten := map[string]bool{
		"testing":                              true,
		modul + "/internal/testpki":            true,
		modul + "/internal/bootstrap/tlsproxy": true,
	}
	pakete := strings.Fields(stdout.String())
	gesehen := false
	for _, p := range pakete {
		if verboten[p] {
			t.Errorf("das Binary linkt %s", p)
		}
		if p == modul+"/internal/bootstrap" {
			gesehen = true
		}
	}
	if !gesehen {
		t.Fatalf("die Liste nennt %s/internal/bootstrap nicht (%d Pakete): sie sagt über das Binary nichts", modul, len(pakete))
	}
}
