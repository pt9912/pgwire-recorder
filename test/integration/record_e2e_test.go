//go:build integration

// Package integration prüft das gebaute Binary Ende-zu-Ende gegen eine reale
// PostgreSQL-Instanz. Den Lauf startet `make test-integration`
// (tools/test/run-integration-tests.sh); er setzt PGR_BINARY und PGR_UPSTREAM.
//
// Jeder Test TestE2E* trägt direkt darüber eine Zeile
// `// Abdeckung: <Kennungen> — <Kurzbeschreibung>`; der Runner bildet daraus
// docs/user/e2e-abdeckung.md.
package integration

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// Der Test prüft das Binary und seine Aufzeichnung von außen; er importiert
// nichts aus internal/.
const codeUpstream = "PGR-E4002"

func TestMain(m *testing.M) {
	if os.Getenv("PGR_BINARY") == "" || os.Getenv("PGR_UPSTREAM") == "" {
		fmt.Fprintln(os.Stderr, "PGR_BINARY und PGR_UPSTREAM fehlen; der Lauf gehört zu make test-integration")
		os.Exit(2)
	}
	if err := warteAufUpstream(os.Getenv("PGR_UPSTREAM"), 90*time.Second); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// Abdeckung: LH-FA-02, LH-FA-06, LH-FA-07 — ein Client führt `SELECT 1;` über `record` gegen eine reale PostgreSQL-Instanz aus und erhält deren Ergebnis; nach dem Beenden des Laufs steht die Interaktion geordnet in einer Aufzeichnung mit Formatkennung und Version.
func TestE2ERecordSelect1(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatalf("Verbindung über den Recorder: %v", err)
	}
	results, err := conn.Exec(ctx, "SELECT 1;").ReadAll()
	if err != nil {
		t.Fatalf("SELECT 1: %v", err)
	}
	if len(results) != 1 || len(results[0].Rows) != 1 || string(results[0].Rows[0][0]) != "1" {
		t.Fatalf("Ergebnis über den Recorder: %#v", results)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatalf("Verbindung schließen: %v", err)
	}

	rec.stop(t, 0)

	text := lies(t, output)
	for _, zeile := range []string{
		"format: pgwire-recorder",
		"version: 1",
		"- id: 1",
		"- sequence: 1",
		"type: query",
		"sql: SELECT 1;",
		"- type: row_description",
		"- type: data_row",
		"- text: \"1\"",
		"- type: command_complete",
		"tag: SELECT 1",
		"- type: ready_for_query",
	} {
		if !strings.Contains(text, zeile) {
			t.Errorf("Aufzeichnung enthält %q nicht:\n%s", zeile, text)
		}
	}
	if strings.Count(text, "- sequence:") != 1 || strings.Count(text, "- id:") != 1 {
		t.Errorf("erwartet genau eine Session mit einer Interaktion:\n%s", text)
	}
	reihenfolge := []string{"row_description", "data_row", "command_complete", "ready_for_query"}
	pos := 0
	for _, typ := range reihenfolge {
		i := strings.Index(text[pos:], "type: "+typ)
		if i < 0 {
			t.Fatalf("Antwort %s fehlt oder steht nicht in der Reihenfolge %v:\n%s", typ, reihenfolge, text)
		}
		pos += i
	}
}

// Abdeckung: LH-FA-02, LH-FA-13 — ist der Upstream nicht erreichbar, erhält der Client eine Fehlerantwort mit PGR-E4002; der Lauf geht weiter und endet beim Beenden mit Exit-Code 0 und einer gültigen Aufzeichnung ohne Session.
func TestE2ERecordUpstreamNichtErreichbar(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, "127.0.0.1:1", output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err == nil || !strings.Contains(err.Error(), codeUpstream) {
		t.Fatalf("erwartet Fehler mit %s, erhalten %v", codeUpstream, err)
	}

	rec.stop(t, 0)

	text := lies(t, output)
	if !strings.Contains(text, "format: pgwire-recorder") || !strings.Contains(text, "sessions: []") {
		t.Fatalf("erwartet eine gültige Aufzeichnung ohne Session:\n%s", text)
	}
}

func lies(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Aufzeichnung lesen: %v", err)
	}
	return string(data)
}

type recorder struct {
	cmd    *exec.Cmd
	listen string
	stderr *strings.Builder
}

func startRecorder(t *testing.T, upstream, output string) *recorder {
	t.Helper()
	listen := freieAdresse(t)
	var stderr strings.Builder
	cmd := exec.Command(os.Getenv("PGR_BINARY"), "record", "--listen", listen, "--upstream", upstream, "--output", output)
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("Recorder starten: %v", err)
	}
	r := &recorder{cmd: cmd, listen: listen, stderr: &stderr}
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	deadline := time.Now().Add(10 * time.Second)
	for {
		c, err := net.DialTimeout("tcp", listen, 200*time.Millisecond)
		if err == nil {
			c.Close()
			return r
		}
		if time.Now().After(deadline) {
			t.Fatalf("Recorder lauscht nicht auf %s: %v\n%s", listen, err, stderr.String())
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (r *recorder) stop(t *testing.T, wantExit int) {
	t.Helper()
	if err := r.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("SIGTERM: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- r.cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		_ = r.cmd.Process.Kill()
		t.Fatalf("Recorder endet nicht nach SIGTERM\n%s", r.stderr.String())
	}
	if code := r.cmd.ProcessState.ExitCode(); code != wantExit {
		t.Fatalf("Exit-Code %d, erwartet %d\n%s", code, wantExit, r.stderr.String())
	}
}

func freieAdresse(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func dsn(address string) string {
	return fmt.Sprintf("postgres://postgres@%s/postgres?sslmode=disable&connect_timeout=10", address)
}

func warteAufUpstream(address string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := pgconn.Connect(ctx, dsn(address))
		cancel()
		if err == nil {
			_ = conn.Close(context.Background())
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("PostgreSQL unter %s nicht bereit: %w", address, err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
