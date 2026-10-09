//go:build unix

package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
)

// Abdeckung: LH-FA-17/Negative — die gewählte Datei ist, Links gefolgt, eine
// reguläre Datei: /dev/null und eine FIFO sind PGR-E2004 „nicht lesbar“ ohne
// Pfad, geprüft vor dem Öffnen, sodass der Start nicht auf die FIFO wartet;
// --config als Link auf eine reguläre Datei wird gelesen (LH-FA-17.a).
func TestDateiKeineRegulaere(t *testing.T) {
	leere(t, "replay")
	if _, err := replayMit("--config=/dev/null"); !istDatei(err) || !strings.Contains(err.Error(), "(--config) nicht lesbar") || strings.Contains(err.Error(), "/dev/null") {
		t.Errorf("/dev/null: %v", err)
	}
	fifo := filepath.Join(t.TempDir(), "fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	ergebnis := make(chan error, 1)
	go func() {
		_, err := replayMit("--config=" + fifo)
		ergebnis <- err
	}()
	select {
	case err := <-ergebnis:
		if !istDatei(err) || !strings.Contains(err.Error(), "(--config) nicht lesbar") || strings.Contains(err.Error(), fifo) {
			t.Errorf("FIFO: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Laden einer FIFO endet nicht binnen 10 s")
	}
	link := filepath.Join(t.TempDir(), "link.yaml")
	if err := os.Symlink(schreibe(t, "log_level: debug\n"), link); err != nil {
		t.Fatal(err)
	}
	if cmd, err := replayMit("--config=" + link); err != nil || cmd.Replay.LogLevel != cli.LogDebug {
		t.Errorf("Link auf eine reguläre Datei: %#v, %v", cmd.Replay, err)
	}
}
