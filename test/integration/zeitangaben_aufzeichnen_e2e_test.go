//go:build integration

package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
)

// zeitAblauf führt über den Recorder unter listen eine einfache Anfrage, nach pause1 eine
// Extended-Interaktion (Parse, Bind, Describe, Execute, Sync in einem Aufruf) und nach
// pause2 eine einfache Anfrage aus; die Pausen liegen beim Client.
func zeitAblauf(t *testing.T, listen string, pause1, pause2 time.Duration) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(listen))
	if err != nil {
		t.Fatalf("Verbindung über den Recorder: %v", err)
	}
	if _, err := conn.Exec(ctx, "SELECT 1").ReadAll(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(pause1)
	if r := conn.ExecParams(ctx, "SELECT $1::int", [][]byte{[]byte("2")}, nil, nil, nil).Read(); r.Err != nil {
		t.Fatal(r.Err)
	}
	time.Sleep(pause2)
	if _, err := conn.Exec(ctx, "SELECT 3").ReadAll(); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
}

// offsetsAus lädt die Aufzeichnung path mit dem Leser des Formats und liefert offset_ms der
// Interaktionen der ersten Session; -1 steht für ein fehlendes Feld.
func offsetsAus(t *testing.T, path string) []int64 {
	t.Helper()
	rec, err := recording.Leser{}.Load(context.Background(), path)
	if err != nil {
		t.Fatalf("laden: %v", err)
	}
	if len(rec.Sessions) != 1 {
		t.Fatalf("%d Sessions, erwartet 1", len(rec.Sessions))
	}
	var out []int64
	for _, in := range rec.Sessions[0].Interactions {
		if in.OffsetMS == nil {
			out = append(out, -1)
			continue
		}
		out = append(out, *in.OffsetMS)
	}
	return out
}

// Abdeckung: LH-FA-21/Happy, LH-FA-22/Happy — Abnahmeszenario 13, Aufzeichnen-Hälfte: record
// --record-timing trägt an jeder Interaktion offset_ms ein, einfach wie Extended, in yaml und in
// sqlite; die Aufzeichnung lädt mit dem Leser des Formats, die Werte steigen nicht ab, und der
// Unterschied der Werte zweier Interaktionen ist mindestens die Pause des Clients zwischen ihnen
// (300 ms und 200 ms; ohne die verdrahtete Systemuhr wäre er 0). In yaml stehen dieselben
// Werte als Text in der Datei.
func TestE2ERecordTiming(t *testing.T) {
	for _, name := range []string{"rec.yaml", "rec.sqlite"} {
		output := filepath.Join(t.TempDir(), name)
		rec := startProzess(t, "record", append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-timing"}, formatArgs(name)...)...)
		zeitAblauf(t, rec.listen, 300*time.Millisecond, 200*time.Millisecond)
		rec.stop(t, 0)

		got := offsetsAus(t, output)
		if len(got) != 3 || got[0] < 0 || got[1]-got[0] < 300 || got[2]-got[1] < 200 {
			t.Errorf("%s: offset_ms %v, erwartet drei Werte, nicht negativ, mit Abständen von mindestens 300 und 200", name, got)
		}
		if strings.HasSuffix(name, ".yaml") {
			var text []int64
			for _, m := range regexp.MustCompile(`(?m)^\s+offset_ms: (\d+)$`).FindAllStringSubmatch(lies(t, output), -1) {
				v, err := strconv.ParseInt(m[1], 10, 64)
				if err != nil {
					t.Fatal(err)
				}
				text = append(text, v)
			}
			if len(text) != len(got) || text[0] != got[0] || text[1] != got[1] || text[2] != got[2] {
				t.Errorf("%s: offset_ms im Text %v, geladen %v", name, text, got)
			}
		}
	}
}

// Abdeckung: LH-FA-21/Boundary — ohne --record-timing fehlt offset_ms an jeder Interaktion in
// yaml wie in sqlite, und die Datei in yaml ist byte-gleich, gleich wie lang der Client
// zwischen den Anfragen wartet.
func TestE2ERecordOhneTiming(t *testing.T) {
	var yaml []string
	for _, pause := range []time.Duration{0, 300 * time.Millisecond} {
		output := filepath.Join(t.TempDir(), "rec.yaml")
		rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output)
		zeitAblauf(t, rec.listen, pause, pause)
		rec.stop(t, 0)
		yaml = append(yaml, lies(t, output))
		if got := offsetsAus(t, output); len(got) != 3 || got[0] != -1 || got[1] != -1 || got[2] != -1 {
			t.Errorf("yaml: offset_ms %v, erwartet ohne Feld", got)
		}
		if strings.Contains(lies(t, output), "offset_ms") {
			t.Errorf("yaml nennt offset_ms ohne die Option")
		}
	}
	if yaml[0] != yaml[1] {
		t.Errorf("yaml unterscheidet sich zwischen den Läufen:\n%s\n---\n%s", yaml[0], yaml[1])
	}
	output := filepath.Join(t.TempDir(), "rec.sqlite")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--format", "sqlite")
	zeitAblauf(t, rec.listen, 0, 0)
	rec.stop(t, 0)
	if got := offsetsAus(t, output); len(got) != 3 || got[0] != -1 || got[1] != -1 || got[2] != -1 {
		t.Errorf("sqlite: offset_ms %v, erwartet ohne Feld", got)
	}
}

// Abdeckung: LH-FA-17/Boundary, LH-FA-21/Happy — PGWIRE_RECORDER_RECORD_TIMING und der Schlüssel
// record_timing im Abschnitt record: der Konfigurationsdatei schalten die Zeitangaben ein wie
// die Option; --record-timing=false auf der Kommandozeile geht beiden vor.
func TestE2ERecordTimingQuellen(t *testing.T) {
	for _, tc := range []struct {
		quelle string
		args   []string
		mit    bool
	}{
		{"Umgebung", nil, true},
		{"Datei", nil, true},
		{"Umgebung mit false auf der Kommandozeile", []string{"--record-timing=false"}, false},
		{"Datei mit false auf der Kommandozeile", []string{"--record-timing=false"}, false},
	} {
		dir := t.TempDir()
		output := filepath.Join(dir, "rec.yaml")
		if strings.HasPrefix(tc.quelle, "Umgebung") {
			t.Setenv("PGWIRE_RECORDER_RECORD_TIMING", "true")
		} else if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("record:\n  record_timing: true\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		rec := startProzessIn(t, dir, "record", append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output}, tc.args...)...)
		zeitAblauf(t, rec.listen, 0, 0)
		rec.stop(t, 0)
		if got := strings.Contains(lies(t, output), "offset_ms"); got != tc.mit {
			t.Errorf("%s: offset_ms in der Aufzeichnung %v, erwartet %v", tc.quelle, got, tc.mit)
		}
		t.Setenv("PGWIRE_RECORDER_RECORD_TIMING", "")
	}
}

// Abdeckung: LH-FA-17/Negative, LH-FA-21/Negative — ein ungültiger Wert von --record-timing
// ist PGR-E2001 mit dem Namen der Option, der der Umgebungsvariable PGR-E2001 mit ihrem Namen
// und der der Datei PGR-E2004 mit dem Schlüssel, je Exit-Code 2, die beiden letzten ohne den
// Wert; bei replay und play ist die
// Option, ihre Umgebungsvariable nicht und ihr Schlüssel in der Datei unbekannt.
func TestE2ERecordTimingUngueltig(t *testing.T) {
	dir := t.TempDir()
	record := []string{"record", "--listen", "127.0.0.1:0", "--upstream", "pg:5432", "--output", filepath.Join(dir, "r.yaml")}
	_, stderr, code := starteBis(t, dir, append(append([]string{}, record...), "--record-timing=GEHEIM")...)
	if code != 2 || !strings.Contains(stderr, "[PGR-E2001]") || !strings.Contains(stderr, "record-timing") {
		t.Errorf("Kommandozeile: Exit-Code %d, stderr %q", code, stderr)
	}
	t.Setenv("PGWIRE_RECORDER_RECORD_TIMING", "GEHEIM")
	_, stderr, code = starteBis(t, dir, record...)
	if code != 2 || !strings.Contains(stderr, "[PGR-E2001]") || !strings.Contains(stderr, "PGWIRE_RECORDER_RECORD_TIMING") || strings.Contains(stderr, "GEHEIM") {
		t.Errorf("Umgebung: Exit-Code %d, stderr %q", code, stderr)
	}
	t.Setenv("PGWIRE_RECORDER_RECORD_TIMING", "")
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("record:\n  record_timing: GEHEIM\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, stderr, code = starteBis(t, dir, record...)
	if code != 2 || !strings.Contains(stderr, "[PGR-E2004]") || !strings.Contains(stderr, "record.record_timing") || strings.Contains(stderr, "GEHEIM") {
		t.Errorf("Datei: Exit-Code %d, stderr %q", code, stderr)
	}
	for _, args := range [][]string{
		{"replay", "--listen", "127.0.0.1:0", "--input", "fehlt.yaml", "--record-timing"},
		{"play", "--upstream", "pg:5432", "--input", "fehlt.yaml", "--record-timing"},
	} {
		d := t.TempDir()
		stdout, stderr, code := starteBis(t, d, args...)
		if code != 2 || stdout != "" || !strings.Contains(stderr, "[PGR-E2001]") || !strings.Contains(stderr, "record-timing") {
			t.Errorf("%s: Exit-Code %d, stdout %q, stderr %q", args[0], code, stdout, stderr)
		}
	}
}

// Abdeckung: LH-FA-01/Happy, LH-FA-17/Happy — die Hilfe von record nennt --record-timing und
// PGWIRE_RECORDER_RECORD_TIMING; config show nennt den Schlüssel record_timing der Datei und,
// bei gesetzter Umgebungsvariable, deren Namen ohne Wert.
func TestE2ERecordTimingHilfeUndConfigShow(t *testing.T) {
	dir := t.TempDir()
	stdout, stderr, code := starteBis(t, dir, "record", "--help")
	if code != 0 || stderr != "" || !strings.Contains(stdout, "--record-timing") || !strings.Contains(stdout, "PGWIRE_RECORDER_RECORD_TIMING") {
		t.Errorf("record --help: Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte("record:\n  record_timing: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PGWIRE_RECORDER_RECORD_TIMING", "true")
	stdout, stderr, code = starteBis(t, dir, "config", "show")
	if code != 0 || stderr != "" || !strings.Contains(stdout, "  record_timing: true\n") || !strings.Contains(stdout, "PGWIRE_RECORDER_RECORD_TIMING\n") || strings.Contains(stdout, "PGWIRE_RECORDER_RECORD_TIMING=") {
		t.Errorf("config show: Exit-Code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
}
