package bootstrap

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
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

// aufzeichnung schreibt eine gültige Aufzeichnung mit einer Session und
// liefert ihren Pfad.
func aufzeichnung(t *testing.T) string {
	t.Helper()
	rec := model.NewRecording()
	rec.Sessions = []model.Session{{ID: 1, Interactions: []model.Interaction{{
		Sequence: 1,
		Request:  model.Request{Type: model.RequestQuery, SQL: "SELECT 1"},
		Responses: []model.Response{
			{Type: model.ResponseCommandComplete, Tag: "SELECT 0"},
			{Type: model.ResponseReadyForQuery, TxStatus: "I"},
		},
	}}}}
	path := filepath.Join(t.TempDir(), "rec.yaml")
	if err := (recording.YAML{}).Write(context.Background(), path, rec); err != nil {
		t.Fatal(err)
	}
	return path
}

// replayBeendet führt replay mit ctx, der schon beendet ist, und den
// Zusatzargumenten aus; die Session der Aufzeichnung bleibt nie zugeordnet
// (PGR-W2001, mit --fail-on-unconsumed PGR-E5002).
func replayBeendet(t *testing.T, args ...string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var stdout, stderr bytes.Buffer
	code := Run(ctx, append([]string{"replay", "--listen", "127.0.0.1:0", "--input", aufzeichnung(t)}, args...), "dev", &stdout, &stderr)
	if stdout.Len() > 0 {
		t.Fatalf("stdout %q", stdout.String())
	}
	return code, stderr.String()
}

// stufenImLog liefert die Werte von level und code je Log-Zeile.
func stufenImLog(t *testing.T, text string) []string {
	t.Helper()
	var out []string
	for _, zeile := range strings.Split(strings.TrimSpace(text), "\n") {
		if zeile == "" {
			continue
		}
		m := logZeile.FindStringSubmatch(zeile)
		if m == nil {
			t.Fatalf("Zeile nicht in der Zeilenform: %q", zeile)
		}
		eintrag := m[1]
		if c := codeAttribut.FindStringSubmatch(zeile); c != nil {
			eintrag += " " + c[1]
		}
		out = append(out, eintrag)
	}
	return out
}

var (
	// logZeile ist die Zeilenform aus LH-FA-14.a: time (RFC 3339 mit
	// Millisekunden und Zonenversatz), level, msg, dann Attribute.
	logZeile     = regexp.MustCompile(`^time=\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d\.\d{3}(?:Z|[+-]\d\d:\d\d) level=(DEBUG|INFO|WARN|ERROR) msg=`)
	codeAttribut = regexp.MustCompile(` code=(PGR-[EW]\d{4})`)
)

// Abdeckung: LH-FA-14/Boundary — --log-level und PGWIRE_RECORDER_LOG_LEVEL
// wirken bei replay: info zeigt Start, Warnung und Ende, warn nur die Warnung
// mit Code, error nur Fehler; jede Zeile hat die Zeilenform logfmt mit time,
// level und msg (LH-FA-14.a).
func TestRunLogLevel(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	t.Setenv("PGWIRE_RECORDER_LOG_LEVEL", "")
	for _, f := range []struct {
		env  string
		args []string
		want []string
	}{
		{"", nil, []string{"INFO", "WARN PGR-W2001", "INFO"}},
		{"", []string{"--log-level=info"}, []string{"INFO", "WARN PGR-W2001", "INFO"}},
		{"", []string{"--log-level=debug"}, []string{"INFO", "WARN PGR-W2001", "INFO"}},
		{"", []string{"--log-level=warn"}, []string{"WARN PGR-W2001"}},
		{"warn", nil, []string{"WARN PGR-W2001"}},
		{"debug", []string{"--log-level=warn"}, []string{"WARN PGR-W2001"}},
		{"", []string{"--log-level=error"}, nil},
		{"", []string{"--log-level=error", "--fail-on-unconsumed"}, []string{"ERROR PGR-E5002"}},
		{"", []string{"--log-level=warn", "--fail-on-unconsumed"}, []string{"ERROR PGR-E5002"}},
	} {
		t.Setenv("PGWIRE_RECORDER_LOG_LEVEL", f.env)
		code, stderr := replayBeendet(t, f.args...)
		want := 0
		if len(f.want) > 0 && strings.HasPrefix(f.want[0], "ERROR") {
			want = 5
		}
		if got := stufenImLog(t, stderr); code != want || strings.Join(got, ",") != strings.Join(f.want, ",") {
			t.Errorf("Umgebung %q, %v: Exit-Code %d, Zeilen %v, erwartet %v\n%s", f.env, f.args, code, got, f.want, stderr)
		}
	}
	_, stderr := replayBeendet(t, "--log-level=error", "--fail-on-unconsumed")
	if !strings.Contains(stderr, ` error="Replay [PGR-E5002]: `) {
		t.Fatalf("Fehlerzeile ohne Attribut error mit Kopf: %q", stderr)
	}
}

// Abdeckung: LH-FA-14/Boundary — eine Stufe zeigt ihre Zeilen und die aller
// strengeren, nicht die der milderen (LH-FA-14.a).
func TestLoggerSchwelle(t *testing.T) {
	for stufe, want := range map[string]string{
		cli.LogError: "ERROR",
		cli.LogWarn:  "WARN,ERROR",
		cli.LogInfo:  "INFO,WARN,ERROR",
		cli.LogDebug: "DEBUG,INFO,WARN,ERROR",
	} {
		var b bytes.Buffer
		log := logger(&b, stufe)
		log.Debug("d")
		log.Info("i")
		log.Warn("w", "code", model.CodeUnconsumed)
		log.Error("e", "code", model.CodeInternal, "error", "x")
		var got []string
		for _, s := range stufenImLog(t, b.String()) {
			got = append(got, strings.Fields(s)[0])
		}
		if strings.Join(got, ",") != want {
			t.Errorf("%s: %v, erwartet %s", stufe, got, want)
		}
	}
}

// Abdeckung: LH-FA-14/Negative, LH-FA-13/Negative — ein Startfehler erscheint
// auf jeder Stufe als Zeile beim Prozessende, genau der Fehlertext, auch wenn
// er den Wert von --log-level oder seiner Umgebungsvariable betrifft; vorher
// schreibt der Prozess nichts anderes, und der Exit-Code ist der der Klasse
// (LH-FA-14.a).
func TestRunStartfehlerJeStufe(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	fehlt := filepath.Join(t.TempDir(), "fehlt.yaml")
	for _, f := range []struct {
		env  string
		args []string
		code string
		exit int
	}{
		{"", []string{"--log-level=error", "--input", fehlt}, model.CodeRecordingIO, 3},
		{"", []string{"--log-level=warn", "--input", fehlt}, model.CodeRecordingIO, 3},
		{"", []string{"--log-level=debug", "--input", fehlt}, model.CodeRecordingIO, 3},
		{"error", []string{"--input", fehlt}, model.CodeRecordingIO, 3},
		{"", []string{"--log-level=INFO", "--input", fehlt}, model.CodeUsage, 2},
		{"", []string{"--log-level=", "--input", fehlt}, model.CodeUsage, 2},
		{"off", []string{"--log-level=error", "--input", fehlt}, model.CodeUsage, 2},
	} {
		t.Setenv("PGWIRE_RECORDER_LOG_LEVEL", f.env)
		var stdout, stderr bytes.Buffer
		exit := Run(context.Background(), append([]string{"replay", "--listen", "127.0.0.1:0"}, f.args...), "dev", &stdout, &stderr)
		zeilen := strings.Split(strings.TrimSuffix(stderr.String(), "\n"), "\n")
		if exit != f.exit || len(zeilen) != 1 || !strings.Contains(zeilen[0], "["+f.code+"]: ") || strings.HasPrefix(zeilen[0], "time=") || stdout.Len() > 0 {
			t.Errorf("Umgebung %q, %v: Exit-Code %d, stderr %q", f.env, f.args, exit, stderr.String())
		}
	}
}

// Abdeckung: LH-FA-13/Negative, LH-FA-14/Negative — die Zeile beim
// Prozessende ist je Fehlerklasse genau der Fehlertext mit Kopf und der
// Exit-Code der der Klasse; ein nicht eingeordneter Fehler ist PGR-E1000 mit
// Exit-Code 1; nebeneinander entstandene Fehler sind je eine Zeile, der
// Exit-Code ist der des ersten (SPEC-034 §Ausgabe).
func TestFailJeKlasse(t *testing.T) {
	for _, f := range []struct {
		err  error
		want string
		exit int
	}{
		{errors.New("intern\nkaputt"), "sonstiger Fehler [PGR-E1000]: intern kaputt\n", 1},
		{model.Errorf(model.CodeUsage, nil, "u"), "Konfiguration [PGR-E2001]: u\n", 2},
		{model.Errorf(model.CodeRecordingBroken, model.Errorf(model.CodeRecordingBroken, nil, "innen"), "aussen"), "Recording [PGR-E3003]: aussen: innen\n", 3},
		{model.Errorf(model.CodeListen, errors.New("in use"), "l"), "Netzwerk [PGR-E4001]: l: in use\n", 4},
		{model.Errorf(model.CodeReplayUnconsumed, nil, "r"), "Replay [PGR-E5002]: r\n", 5},
		{model.Errorf(model.CodeProtocolVersion, nil, "p"), "nicht unterstützt [PGR-E6002]: p\n", 6},
		{errors.Join(model.Errorf(model.CodeRecordingIO, nil, "a"), model.Errorf(model.CodeNetwork, nil, "b")), "Recording [PGR-E3001]: a\nNetzwerk [PGR-E4000]: b\n", 3},
	} {
		var stderr bytes.Buffer
		if exit := fail(&stderr, f.err); exit != f.exit || stderr.String() != f.want {
			t.Errorf("Exit-Code %d, stderr %q, erwartet %d, %q", exit, stderr.String(), f.exit, f.want)
		}
	}
}

// Abdeckung: LH-FA-14/Boundary — --log-level wirkt bei record: info zeigt Start
// und Ende, warn und error zeigen davon nichts (LH-FA-14.a).
func TestRunRecordLogLevel(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_LOG_LEVEL", "")
	for args, want := range map[string]string{
		"":                  "INFO,INFO",
		"--log-level=info":  "INFO,INFO",
		"--log-level=warn":  "",
		"--log-level=error": "",
	} {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		aufruf := []string{"record", "--listen", "127.0.0.1:0", "--upstream", "127.0.0.1:1", "--output", filepath.Join(t.TempDir(), "rec.yaml")}
		if args != "" {
			aufruf = append(aufruf, args)
		}
		var stdout, stderr bytes.Buffer
		code := Run(ctx, aufruf, "dev", &stdout, &stderr)
		if got := strings.Join(stufenImLog(t, stderr.String()), ","); code != 0 || got != want || stdout.Len() > 0 {
			t.Errorf("%q: Exit-Code %d, Zeilen %q, erwartet %q\n%s", args, code, got, want, stderr.String())
		}
	}
}

// Abdeckung: LH-FA-14/Boundary — time einer Log-Zeile ist Ortszeit mit
// Millisekunden und dem Zonenversatz der Ortszeit (LH-FA-14.a §Zeilenform).
func TestLoggerOrtszeit(t *testing.T) {
	alt := time.Local
	time.Local = time.FixedZone("Test", 2*60*60)
	defer func() { time.Local = alt }()
	var b bytes.Buffer
	logger(&b, cli.LogInfo).Info("x")
	if !regexp.MustCompile(`^time=\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d\.\d{3}\+02:00 level=INFO msg=x\n$`).MatchString(b.String()) {
		t.Fatalf("Zeile %q", b.String())
	}
}
