package bootstrap_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-13/Boundary, LH-FA-13/Negative — nach einem kontrollierten
// Herunterfahren ist der Exit-Code 0 ohne gemerkten Fehler (SPEC-013), sonst
// der der Klasse des gemerkten Codes: 1 (SPEC-015), 2 (SPEC-014), 3
// (SPEC-016), 4 (SPEC-017), 5 (SPEC-018), 6 (SPEC-019) (LH-FA-13.b, SPEC-038
// Exit-Code-Zuordnung).
func TestExitCodeJeKlasse(t *testing.T) {
	for code, want := range map[string]int{
		"":                        0,
		model.CodeInternal:        1,
		model.CodeUsage:           2,
		model.CodeRecordingIO:     3,
		model.CodeShutdownTimeout: 4,
		model.CodeConnectionLost:  4,
		model.CodeReplayMismatch:  5,
		model.CodeUnsupported:     6,
	} {
		if got := bootstrap.ExitCode(code); got != want {
			t.Errorf("Code %q: Exit-Code %d, erwartet %d", code, got, want)
		}
	}
}

// startnachricht ist die Startnachricht der Protokollversion 3.0 mit
// user=app, von Hand kodiert: Länge, Version, Parameter, abschließendes
// Nullbyte.
func startnachricht() []byte {
	rumpf := append(binary.BigEndian.AppendUint32(nil, 196608), "user\x00app\x00\x00"...)
	return append(binary.BigEndian.AppendUint32(nil, uint32(4+len(rumpf))), rumpf...)
}

// Abdeckung: LH-FA-13/Negative — merkt record einen Verbindungsfehler (der
// Upstream ist nicht erreichbar, PGR-E4002), endet der Lauf nach dem Signal mit
// dessen Exit-Code 4; scheitert dann das Schreiben der Aufzeichnung, ist der
// Exit-Code 3 und hat Vorrang vor der gemerkten Klasse (LH-FA-13.b).
func TestRunRecordSchreibfehlerVorGemerkterKlasse(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_LOG_LEVEL", "")
	t.Setenv("PGWIRE_RECORDER_SHUTDOWN_TIMEOUT", "")
	for _, f := range []struct {
		name        string
		schreiben   bool
		want        int
		letzteZeile string
	}{
		{"Aufzeichnung geschrieben", true, 4, "msg=\"record beendet\""},
		{"Schreibfehler am Ende", false, 3, "Recording [PGR-E3001]: "},
	} {
		t.Run(f.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "ziel")
			if err := os.Mkdir(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			listen := freieAdresse(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var stderr syncPuffer
			ende := make(chan int, 1)
			go func() {
				var stdout bytes.Buffer
				ende <- bootstrap.Run(ctx, nil, []string{"record", "--listen", listen, "--upstream", "127.0.0.1:1", "--output", filepath.Join(dir, "rec.yaml")}, "dev", &stdout, &stderr)
			}()
			var client net.Conn
			for deadline := time.Now().Add(5 * time.Second); ; {
				c, err := net.Dial("tcp", listen)
				if err == nil {
					client = c
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("record lauscht nicht binnen 5 s")
				}
				time.Sleep(20 * time.Millisecond)
			}
			defer client.Close()
			if _, err := client.Write(startnachricht()); err != nil {
				t.Fatal(err)
			}
			// Die Fehlerantwort mit PGR-E4002 kommt, danach schließt record die
			// Verbindung.
			_ = client.SetReadDeadline(time.Now().Add(5 * time.Second))
			antwort, err := io.ReadAll(client)
			if err != nil || !bytes.Contains(antwort, []byte("PGR-E4002")) {
				t.Fatalf("Antwort %q, %v", antwort, err)
			}
			if !f.schreiben {
				if err := os.RemoveAll(dir); err != nil {
					t.Fatal(err)
				}
			}
			cancel()
			var code int
			select {
			case code = <-ende:
			case <-time.After(5 * time.Second):
				t.Fatal("Run endet nicht binnen 5 s")
			}
			zeilen := strings.Split(strings.TrimSuffix(stderr.String(), "\n"), "\n")
			if code != f.want || !strings.Contains(zeilen[len(zeilen)-1], f.letzteZeile) {
				t.Fatalf("Exit-Code %d, erwartet %d; stderr\n%s", code, f.want, stderr.String())
			}
		})
	}
}
