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
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
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

// Abdeckung: LH-FA-20/Negative, LH-RB-01/Messung — die Datei aus --upstream-ca ist,
// Links gefolgt, eine reguläre Datei: Eine FIFO ist PGR-E2007 ohne Pfad, geprüft
// vor dem Öffnen, sodass der Start nicht auf sie wartet; ein Link auf eine
// reguläre Datei mit einem Zertifikat wird gelesen, ein Link auf eine FIFO ist
// PGR-E2007 (LH-FA-20.a *Start*).
func TestPlayUpstreamCAKeineRegulaere(t *testing.T) {
	leere(t, "play")
	dir := t.TempDir()
	fifo := filepath.Join(dir, "GEHEIMFIFO")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	linkAufFifo := filepath.Join(dir, "GEHEIMLINK")
	if err := os.Symlink(fifo, linkAufFifo); err != nil {
		t.Fatal(err)
	}
	for name, pfad := range map[string]string{"FIFO": fifo, "Link auf eine FIFO": linkAufFifo, "Gerät": "/dev/null"} {
		ergebnis := make(chan error, 1)
		go func() {
			_, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="+pfad)
			ergebnis <- err
		}()
		select {
		case err := <-ergebnis:
			if !hatCode(err, model.CodeConfigCA) || !strings.Contains(err.Error(), "keine reguläre Datei") || strings.Contains(err.Error(), "GEHEIM") || strings.Contains(err.Error(), "/dev/null") {
				t.Errorf("%s: %v", name, err)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("%s: Lesen von --upstream-ca endet nicht binnen 10 s", name)
		}
	}
	link := filepath.Join(dir, "link.pem")
	if err := os.Symlink(schreibe(t, string(testpki.NeueCA(t, "CA").PEM())), link); err != nil {
		t.Fatal(err)
	}
	if got, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="+link); err != nil || len(got.UpstreamCA) != 1 {
		t.Errorf("Link auf eine reguläre Datei: %d Zertifikate, %v", len(got.UpstreamCA), err)
	}
}
