//go:build integration

// Package integration_test prüft das gebaute Binary Ende-zu-Ende gegen eine reale
// PostgreSQL-Instanz. Den Lauf startet `make test-integration`
// (tools/test/run-integration-tests.sh); er setzt PGR_BINARY und PGR_UPSTREAM.
//
// Jeder Test TestE2E* trägt direkt darüber eine Abdeckungs-Deklaration
// `// Abdeckung: <Anforderung>/<Pfad>, … — <Kurzbeschreibung>`; `make abdeckung`
// bildet daraus docs/user/abdeckung-e2e.md.
package integration_test

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
	if os.Getenv("PGR_OHNE_POSTGRES") == "1" {
		// Zweite Phase des Runners: PostgreSQL ist gestoppt.
		os.Exit(m.Run())
	}
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

// Abdeckung: LH-FA-02/Happy — ein Client führt
// `SELECT 1;` über `record` gegen eine reale PostgreSQL-Instanz aus und erhält
// deren Ergebnis; nach dem Beenden des Laufs steht die Interaktion geordnet in
// einer Aufzeichnung mit Formatkennung und Version, ohne Adressen und Pfade des
// Laufs.
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
	for _, fremd := range []string{rec.listen, os.Getenv("PGR_UPSTREAM"), output, filepath.Dir(output)} {
		if strings.Contains(text, fremd) {
			t.Errorf("Aufzeichnung enthält die rechnerspezifische Angabe %q:\n%s", fremd, text)
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

// Abdeckung: LH-FA-02/Negative, LH-FA-13/Boundary, LH-FA-13/Negative — ist der Upstream nicht
// erreichbar, erhält der Client eine Fehlerantwort mit PGR-E4002; der Lauf geht
// weiter und endet beim Beenden mit dem Exit-Code der Klasse dieses Fehlers (4)
// und einer gültigen Aufzeichnung ohne Session.
func TestE2ERecordUpstreamNichtErreichbar(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, "127.0.0.1:1", output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err == nil || !strings.Contains(err.Error(), codeUpstream) {
		t.Fatalf("erwartet Fehler mit %s, erhalten %v", codeUpstream, err)
	}

	rec.stop(t, 4)

	text := lies(t, output)
	if !strings.Contains(text, "format: pgwire-recorder") || !strings.Contains(text, "sessions: []") {
		t.Fatalf("erwartet eine gültige Aufzeichnung ohne Session:\n%s", text)
	}
}

// Abdeckung: LH-FA-13/Happy — mehrere Interaktionen mit
// mehreren Spalten und NULL stehen in Reihenfolge und mit ihren Werten in der
// Aufzeichnung; ein fehlerfreier Lauf endet nach SIGTERM mit Exit-Code 0.
func TestE2ERecordMehrereInteraktionen(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatalf("Verbindung über den Recorder: %v", err)
	}
	for _, q := range []string{"SELECT 'a' AS spalte_a, NULL AS spalte_b;", "SELECT 2 AS z;"} {
		if _, err := conn.Exec(ctx, q).ReadAll(); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	_ = conn.Close(ctx)
	rec.stop(t, 0)

	text := lies(t, output)
	reihe := []string{
		"- sequence: 1", "sql: SELECT 'a' AS spalte_a, NULL AS spalte_b;", "name: spalte_a", "name: spalte_b", "- text: a", `- "null": true`,
		"- sequence: 2", "sql: SELECT 2 AS z;", "name: z", `- text: "2"`,
	}
	pos := 0
	for _, z := range reihe {
		i := strings.Index(text[pos:], z)
		if i < 0 {
			t.Fatalf("%q fehlt oder steht nicht in der Reihenfolge:\n%s", z, text)
		}
		pos += i + len(z)
	}
}

// Abdeckung: LH-FA-05/Boundary — fordert der Client TLS an (sslmode=prefer),
// lehnt der Recorder ab, und der Client verbindet sich unverschlüsselt.
func TestE2ERecordSSLAbgelehnt(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, strings.Replace(dsn(rec.listen), "sslmode=disable", "sslmode=prefer", 1))
	if err != nil {
		t.Fatalf("Verbindung mit sslmode=prefer: %v", err)
	}
	if _, err := conn.Exec(ctx, "SELECT 1;").ReadAll(); err != nil {
		t.Fatal(err)
	}
	_ = conn.Close(ctx)
	rec.stop(t, 0)
}

// Abdeckung: LH-FA-05/Negative, LH-FA-06/Negative, LH-FA-13/Negative — eine nicht unterstützte
// Interaktion (`COPY … TO STDOUT`) beendet die Verbindung mit PGR-E6001; die
// Session wird nicht aufgezeichnet, auch nicht ihre vorherige Interaktion, und
// der Lauf endet mit Exit-Code 6.
func TestE2ERecordNichtUnterstuetzt(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(ctx, "SELECT 1;").ReadAll(); err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec(ctx, "COPY (SELECT 1) TO STDOUT;").ReadAll()
	if err == nil || !strings.Contains(err.Error(), "PGR-E6001") {
		t.Fatalf("erwartet PGR-E6001, erhalten %v", err)
	}
	_ = conn.Close(ctx)
	rec.stop(t, 6)

	if text := lies(t, output); !strings.Contains(text, "sessions: []") {
		t.Fatalf("verworfene Session aufgezeichnet:\n%s", text)
	}
}

// Abdeckung: LH-FA-13/Boundary — eine offene, ruhende Client-Verbindung hält das
// Beenden nach SIGTERM nicht auf; ihre abgeschlossenen Interaktionen stehen in
// der Aufzeichnung.
func TestE2ERecordBeendenMitOffenerVerbindung(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			t.Error(err)
		}
	}()
	if _, err := conn.Exec(ctx, "SELECT 1;").ReadAll(); err != nil {
		t.Fatal(err)
	}
	rec.stop(t, 0)

	if text := lies(t, output); !strings.Contains(text, "sql: SELECT 1;") {
		t.Fatalf("Interaktion der offenen Verbindung fehlt:\n%s", text)
	}
}

// Abdeckung: LH-FA-02/Boundary — ein Client, der sich anmeldet und ohne Anfrage
// trennt, wird nicht aufgezeichnet; der Recorder bleibt nicht hängen, und die
// Aufzeichnung ist gültig.
func TestE2ERecordVerbindungOhneAnfrage(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	_ = conn.Close(ctx)
	rec.stop(t, 0)

	if text := lies(t, output); !strings.Contains(text, "sessions: []") {
		t.Fatalf("Verbindung ohne Anfrage aufgezeichnet:\n%s", text)
	}
}

// Abdeckung: LH-FA-13/Happy — trennt ein Client nach seiner Anfrage ohne
// Terminate, ist das ein reguläres Ende: die Session wird aufgezeichnet, und der
// Lauf endet mit Exit-Code 0 (LH-FA-02.b).
func TestE2ERecordEndeOhneTerminate(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(ctx, "SELECT 1;").ReadAll(); err != nil {
		t.Fatal(err)
	}
	// Socket schließen, ohne Terminate zu senden.
	_ = conn.Conn().Close()
	time.Sleep(200 * time.Millisecond)
	rec.stop(t, 0)

	if text := lies(t, output); !strings.Contains(text, "sql: SELECT 1;") {
		t.Fatalf("Session ohne Terminate fehlt:\n%s", text)
	}
}

// Abdeckung: LH-FA-13/Negative — scheitert das Schreiben der Aufzeichnung beim
// Ende des Laufs, endet er mit Exit-Code 3 (Klasse Recording), auch ohne
// vorherigen Verbindungsfehler.
func TestE2ERecordSchreibfehlerAmEnde(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ziel")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), filepath.Join(dir, "rec.yaml"))
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	rec.stop(t, 3)
	if !strings.Contains(rec.stderr.String(), "PGR-E3001") {
		t.Fatalf("erwartet PGR-E3001 in der Ausgabe:\n%s", rec.stderr.String())
	}
}

// Abdeckung: LH-FA-05/Boundary, LH-FA-13/Happy — eine HTTP-Anfrage an den Port
// des Recorders (etwa ein Gesundheitscheck) wird ohne Antwort geschlossen; sie
// zählt nicht als Fehler, und der Lauf endet mit Exit-Code 0.
func TestE2ERecordFremdeAnfrage(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	c, err := net.DialTimeout("tcp", rec.listen, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_ = c.SetDeadline(time.Now().Add(5 * time.Second))
	if _, err := c.Write([]byte("GET /health HTTP/1.1\r\nHost: x\r\n\r\n")); err != nil {
		t.Fatal(err)
	}
	if n, err := c.Read(make([]byte, 1)); err == nil {
		t.Fatalf("Antwort erhalten (%d Bytes) statt geschlossener Verbindung", n)
	}
	c.Close()
	rec.stop(t, 0)
	if !strings.Contains(rec.stderr.String(), "PGR-W3003") {
		t.Fatalf("Warnung PGR-W3003 fehlt:\n%s", rec.stderr.String())
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

// recorder ist ein gestarteter Prozess des Binaries. Wait ruft nur die
// Goroutine aus startProzess: Sie legt den Fehler von Wait in waitErr ab und
// schließt danach beendet. warteEnde, signal, nachKill und der Cleanup aus
// startProzess warten auf beendet; pruefeExit liest cmd.ProcessState und
// waitErr nach warteEnde.
type recorder struct {
	cmd     *exec.Cmd
	listen  string
	stderr  *strings.Builder
	beendet chan struct{}
	waitErr error
}

func startRecorder(t *testing.T, upstream, output string) *recorder {
	t.Helper()
	return startProzess(t, "record", "--upstream", upstream, "--output", output)
}

// startProzess startet das Binary mit einem Kommando, einer freien
// --listen-Adresse und den übrigen Argumenten und wartet höchstens 10 s, bis es
// lauscht; sonst endet es über nachKill mit t.Fatalf. Ist der Prozess beim
// Cleanup des Tests nicht beendet, ruft der Cleanup Kill und wartet höchstens
// 5 s auf das Ende; danach t.Errorf.
func startProzess(t *testing.T, kommando string, args ...string) *recorder {
	t.Helper()
	return startProzessIn(t, "", kommando, args...)
}

// startProzessIn ist startProzess im Verzeichnis dir; "" ist das Verzeichnis
// des Tests.
func startProzessIn(t *testing.T, dir, kommando string, args ...string) *recorder {
	t.Helper()
	listen := freieAdresse(t)
	var stderr strings.Builder
	cmd := exec.Command(os.Getenv("PGR_BINARY"), append([]string{kommando, "--listen", listen}, args...)...)
	cmd.Dir = dir
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("Recorder starten: %v", err)
	}
	r := &recorder{cmd: cmd, listen: listen, stderr: &stderr, beendet: make(chan struct{})}
	go func() {
		r.waitErr = cmd.Wait()
		close(r.beendet)
	}()
	t.Cleanup(func() {
		select {
		case <-r.beendet:
			return
		default:
		}
		_ = cmd.Process.Kill()
		select {
		case <-r.beendet:
		case <-time.After(5 * time.Second):
			t.Errorf("Prozess %s endet im Cleanup auch nach Kill nicht binnen 5 s", kommando)
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
			t.Fatal(r.nachKill(fmt.Sprintf("Recorder lauscht nicht auf %s: %v", listen, err)))
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (r *recorder) stop(t *testing.T, wantExit int) {
	t.Helper()
	r.signal(t, syscall.SIGTERM)
	r.warteEnde(t, 15*time.Second, "Recorder endet nicht nach SIGTERM")
	r.pruefeExit(t, wantExit)
}

// signal sendet sig an den Prozess. Liefert Signal einen Fehler, wartet es
// höchstens 5 s auf das Ende: Endet der Prozess, meldet t.Fatalf „Prozess
// endete vor dem Signal“ mit ProcessState.String(), dem Fehler von Wait und
// stderr, sonst den Fehler von Signal und dass der Prozess binnen 5 s nicht
// endete.
func (r *recorder) signal(t *testing.T, sig os.Signal) {
	t.Helper()
	err := r.cmd.Process.Signal(sig)
	if err == nil {
		return
	}
	select {
	case <-r.beendet:
		t.Fatalf("Prozess endete vor dem Signal %v (%s, Wait: %v)\n%s", sig, r.cmd.ProcessState.String(), r.waitErr, r.stderr.String())
	case <-time.After(5 * time.Second):
		t.Fatalf("Signal %v: %v; Prozess endete auch binnen 5 s nicht", sig, err)
	}
}

// warteEnde wartet höchstens frist auf das Ende des Prozesses. Läuft die Frist
// ab, endet es über nachKill mit t.Fatalf.
func (r *recorder) warteEnde(t *testing.T, frist time.Duration, meldung string) {
	t.Helper()
	select {
	case <-r.beendet:
		return
	case <-time.After(frist):
	}
	t.Fatal(r.nachKill(meldung))
}

// nachKill ruft Kill und wartet höchstens 5 s auf das Ende. Es liefert meldung
// und stderr, oder meldung und den Hinweis, dass der Prozess auch nach Kill
// nicht endet.
func (r *recorder) nachKill(meldung string) string {
	_ = r.cmd.Process.Kill()
	select {
	case <-r.beendet:
		return meldung + "\n" + r.stderr.String()
	case <-time.After(5 * time.Second):
		return meldung + "; endet auch nach Kill nicht binnen 5 s"
	}
}

// pruefeExit vergleicht nach warteEnde den Exit-Code mit wantExit und nennt bei
// Abweichung ProcessState.String() und den Fehler von Wait.
func (r *recorder) pruefeExit(t *testing.T, wantExit int) {
	t.Helper()
	if code := r.cmd.ProcessState.ExitCode(); code != wantExit {
		t.Fatalf("Exit-Code %d, erwartet %d (%s, Wait: %v)\n%s", code, wantExit, r.cmd.ProcessState.String(), r.waitErr, r.stderr.String())
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
