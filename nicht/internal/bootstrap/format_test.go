package bootstrap_test

import (
	"bytes"
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-22/Boundary — jeder Wert, den die Kommandozeile für --format annimmt,
// nimmt auch die Fabrik der Aufzeichnungsadapter an.
func TestFormateDerKommandozeileKenntDieFabrik(t *testing.T) {
	werte := cli.Formate()
	if len(werte) < 2 {
		t.Fatalf("Formate: %v", werte)
	}
	for _, f := range werte {
		if _, err := recording.New(f); err != nil {
			t.Errorf("--format=%s: die Fabrik lehnt ab: %v", f, err)
		}
	}
}

// Abdeckung: LH-FA-22/Happy, LH-FA-22/Boundary — record schreibt das Format aus --format
// in die Zieldatei: ohne die Option und mit yaml die Textaufzeichnung, mit sqlite eine
// Datei mit der Kopfzeile einer SQLite-Datei (SQLite format 3 und ein NUL-Byte), auch
// ohne Session beim Beenden.
func TestRunRecordFormat(t *testing.T) {
	for name, tc := range map[string]struct {
		args  []string
		kopf  string
		nicht string
	}{
		"Standard": {nil, "format: pgwire-recorder", "SQLite format 3"},
		"yaml":     {[]string{"--format", "yaml"}, "format: pgwire-recorder", "SQLite format 3"},
		"sqlite":   {[]string{"--format", "sqlite"}, "SQLite format 3\x00", "format: pgwire-recorder"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("PGWIRE_RECORDER_LOG_LEVEL", "")
			t.Setenv("PGWIRE_RECORDER_SHUTDOWN_TIMEOUT", "")
			t.Setenv("PGWIRE_RECORDER_FORMAT", "")
			upstream, _ := stummerUpstream(t)
			listen := freieAdresse(t)
			output := filepath.Join(t.TempDir(), "rec.out")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ende := make(chan int, 1)
			stderr := &syncPuffer{}
			go func() {
				var stdout bytes.Buffer
				ende <- bootstrap.Run(ctx, nil, append([]string{"record", "--listen", listen, "--upstream", upstream, "--output", output}, tc.args...), "dev", &stdout, stderr)
			}()
			for deadline := time.Now().Add(5 * time.Second); ; {
				c, err := net.Dial("tcp", listen)
				if err == nil {
					c.Close()
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("record lauscht nicht binnen 5 s")
				}
				time.Sleep(20 * time.Millisecond)
			}
			cancel()
			select {
			case code := <-ende:
				if code != 0 {
					t.Fatalf("Exit-Code %d\n%s", code, stderr.String())
				}
			case <-time.After(5 * time.Second):
				t.Fatalf("Run endet nicht binnen 5 s\n%s", stderr.String())
			}
			data, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(string(data), tc.kopf) || strings.HasPrefix(string(data), tc.nicht) {
				t.Fatalf("Datei beginnt mit %q", data[:min(len(data), 32)])
			}
		})
	}
}

// Abdeckung: LH-FA-22/Negative — ein Wert außerhalb von yaml und sqlite in --format endet
// record beim Start mit Exit-Code 2 (PGR-E2001), ohne die Zieldatei anzulegen.
func TestRunRecordFormatUngueltig(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FORMAT", "")
	output := filepath.Join(t.TempDir(), "rec.out")
	var stdout bytes.Buffer
	stderr := &syncPuffer{}
	code := bootstrap.Run(context.Background(), nil, []string{"record", "--listen", "127.0.0.1:1", "--upstream", "h:1", "--output", output, "--format", "SQLite"}, "dev", &stdout, stderr)
	if code != 2 || !strings.Contains(stderr.String(), model.CodeUsage) {
		t.Fatalf("Exit-Code %d\n%s", code, stderr.String())
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("Zieldatei: %v", err)
	}
}
