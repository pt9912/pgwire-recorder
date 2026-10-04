//go:build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// aufnehmen zeichnet die Anfragen über record gegen die reale Instanz auf und
// liefert den Pfad der Aufzeichnung.
func aufnehmen(t *testing.T, queries ...string) string {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range queries {
		if _, err := conn.Exec(ctx, q).ReadAll(); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	_ = conn.Close(ctx)
	rec.stop(t, 0)
	return output
}

// beobachten führt die Anfragen über replay aus und liefert das beobachtbare
// Verhalten als Text: Spaltennamen, Zeilen und Befehlsabschluss je Ergebnis.
func beobachten(t *testing.T, listen string, queries ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(listen))
	if err != nil {
		t.Fatalf("Verbindung zum Replay: %v", err)
	}
	defer conn.Close(ctx)
	var b strings.Builder
	for _, q := range queries {
		results, err := conn.Exec(ctx, q).ReadAll()
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		for _, r := range results {
			for _, f := range r.FieldDescriptions {
				fmt.Fprintf(&b, "spalte %s %d\n", f.Name, f.DataTypeOID)
			}
			for _, row := range r.Rows {
				fmt.Fprintf(&b, "zeile %q\n", row)
			}
			fmt.Fprintf(&b, "befehl %s\n", r.CommandTag.String())
		}
	}
	return b.String()
}

// Abdeckung: LH-FA-03/Happy, LH-FA-07/Happy, LH-FA-09/Happy, LH-QA-01/Messung —
// eine von record geschriebene Aufzeichnung beantwortet replay ohne Upstream
// mit dem Verhalten der realen Instanz; zehn aufeinanderfolgende Läufe zeigen
// dasselbe beobachtbare Verhalten.
func TestE2EReplaySelect1(t *testing.T) {
	queries := []string{"SELECT 1 AS eins;", "SELECT 'a' AS text, NULL AS leer;"}
	input := aufnehmen(t, queries...)

	var erstes string
	for lauf := 1; lauf <= 10; lauf++ {
		rep := startProzess(t, "replay", "--input", input)
		got := beobachten(t, rep.listen, queries...)
		rep.stop(t, 0)
		if lauf == 1 {
			erstes = got
			if !strings.Contains(got, `zeile ["1"]`) || !strings.Contains(got, "befehl SELECT 1") {
				t.Fatalf("Replay-Ergebnis:\n%s", got)
			}
			continue
		}
		if got != erstes {
			t.Fatalf("Lauf %d weicht ab:\n%s\n--- erster Lauf:\n%s", lauf, got, erstes)
		}
	}
}

// Abdeckung: LH-FA-10/Happy, LH-FA-13/Negative — eine Anfrage, die nicht zur
// Aufzeichnung passt, erhält einen eindeutigen Fehler mit PGR-E5001 und keine
// Antwort; der Lauf endet mit Exit-Code 5.
func TestE2EReplayAbweichung(t *testing.T) {
	input := aufnehmen(t, "SELECT 1;")
	rep := startProzess(t, "replay", "--input", input)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rep.listen))
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec(ctx, "SELECT 2;").ReadAll()
	if err == nil || !strings.Contains(err.Error(), "PGR-E5001") {
		t.Fatalf("erwartet PGR-E5001, erhalten %v", err)
	}
	_ = conn.Close(ctx)
	rep.stop(t, 5)
}

// Abdeckung: LH-FA-03/Negative — eine beschädigte Aufzeichnung beendet den Start
// mit Exit-Code 3 und PGR-E3003.
func TestE2EReplayBeschaedigt(t *testing.T) {
	input := filepath.Join(t.TempDir(), "rec.yaml")
	if err := os.WriteFile(input, []byte("format: pgwire-recorder\nversion: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := exec_(t, "replay", "--listen", freieAdresse(t), "--input", input)
	if code := exitCodeOf(err); code != 3 || !strings.Contains(out, "PGR-E3003") {
		t.Fatalf("Exit-Code %d, Ausgabe:\n%s", code, out)
	}
}

func exec_(t *testing.T, args ...string) (string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, os.Getenv("PGR_BINARY"), args...).CombinedOutput()
	return string(out), err
}

func exitCodeOf(err error) int {
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	if err != nil {
		return -1
	}
	return 0
}
