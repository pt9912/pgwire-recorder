//go:build unix

package cli_test

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-23/Negative, LH-FA-23/Boundary — jede der beiden Dateien ist,
// Links gefolgt, eine reguläre Datei: Ein Verzeichnis, eine FIFO, ein Gerät (/dev/null), ein Socket und ein Link auf
// ein Verzeichnis sind PGR-E2007 „keine reguläre Datei“ ohne Pfad, geprüft vor
// dem Öffnen, sodass der Start nicht auf die FIFO wartet; ein Link auf eine
// reguläre Datei wird gelesen (LH-FA-23.a *Dateien*).
func TestTLSDateiKeineRegulaere(t *testing.T) {
	cert, key, _ := tlsPaar(t)
	dir := t.TempDir()
	fifo, verz, linkDatei, linkVerz := filepath.Join(dir, "GEHEIMfifo"), filepath.Join(dir, "GEHEIMverz"), filepath.Join(dir, "GEHEIMlink"), filepath.Join(dir, "GEHEIMlinkverz")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	sock := filepath.Join(dir, "GEHEIMsock")
	l, err := net.Listen("unix", sock)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	for _, e := range []struct{ was, ziel, pfad string }{{"mkdir", "", verz}, {"symlink", key, linkDatei}, {"symlink", verz, linkVerz}} {
		var err error
		if e.was == "mkdir" {
			err = os.Mkdir(e.pfad, 0o700)
		} else {
			err = os.Symlink(e.ziel, e.pfad)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, kommando := range tlsKommandos() {
		for _, pfad := range []string{fifo, verz, linkVerz, "/dev/null", sock} {
			for _, option := range []string{"--tls-cert", "--tls-key"} {
				args := []string{"--tls-cert=" + cert, "--tls-key=" + key}
				if option == "--tls-cert" {
					args[0] = "--tls-cert=" + pfad
				} else {
					args[1] = "--tls-key=" + pfad
				}
				_, err := tlsLese(t, kommando, args...)
				if !hatCode(err, model.CodeConfigCA) || !strings.HasPrefix(err.Error(), "Konfiguration [PGR-E2007]: "+option+": ") || !strings.Contains(err.Error(), "keine reguläre Datei") || strings.Contains(err.Error(), "GEHEIM") {
					t.Errorf("%s %s=%s: %v", kommando, option, filepath.Base(pfad), err)
				}
			}
		}
		if _, err := tlsLese(t, kommando, "--tls-cert="+cert, "--tls-key="+linkDatei); err != nil {
			t.Errorf("%s: Link auf eine reguläre Datei: %v", kommando, err)
		}
	}
}
