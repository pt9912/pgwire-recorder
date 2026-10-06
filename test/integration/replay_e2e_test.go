//go:build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// aufnehmen zeichnet die Anfragen über record gegen die reale Instanz auf und
// liefert den Pfad der Aufzeichnung und die Sicht des Clients beim Aufzeichnen.
func aufnehmen(t *testing.T, queries ...string) (string, string) {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	sicht := beobachten(t, rec.listen, queries...)
	rec.stop(t, 0)
	return output, sicht
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

// Abdeckung: LH-FA-07/Happy, LH-FA-09/Happy — eine von record geschriebene
// Aufzeichnung beantwortet replay mit derselben Sicht des Clients wie beim
// Aufzeichnen (Spaltennamen, Typ-OIDs, Zeilen, Befehlsabschluss); zehn
// aufeinanderfolgende Läufe zeigen dieselbe Sicht.
func TestE2EReplaySelect1(t *testing.T) {
	queries := []string{"SELECT 1 AS eins;", "SELECT 'a' AS text, NULL AS leer;"}
	input, aufgezeichnet := aufnehmen(t, queries...)
	if !strings.Contains(aufgezeichnet, `zeile ["1"]`) || !strings.Contains(aufgezeichnet, "spalte eins 23") {
		t.Fatalf("Sicht beim Aufzeichnen:\n%s", aufgezeichnet)
	}
	for lauf := 1; lauf <= 10; lauf++ {
		rep := startProzess(t, "replay", "--input", input)
		got := beobachten(t, rep.listen, queries...)
		rep.stop(t, 0)
		if got != aufgezeichnet {
			t.Fatalf("Lauf %d weicht von der Sicht beim Aufzeichnen ab:\n%s\n--- aufgezeichnet:\n%s", lauf, got, aufgezeichnet)
		}
	}
}

// Abdeckung: LH-FA-02/Happy — Vorbereitung der Phase ohne PostgreSQL: zwei
// Verbindungen werden über record aufgezeichnet; die Aufzeichnung und die Sicht
// der ersten Verbindung liegen danach unter PGR_DATEN.
func TestE2EVorbereitungOhnePostgres(t *testing.T) {
	daten := os.Getenv("PGR_DATEN")
	if daten == "" {
		t.Skip("ohne PGR_DATEN keine zweite Phase")
	}
	output := filepath.Join(daten, "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	sicht := beobachten(t, rec.listen, "SELECT 1;")
	_ = beobachten(t, rec.listen, "SELECT 2;")
	rec.stop(t, 0)
	if err := os.WriteFile(filepath.Join(daten, "sicht.txt"), []byte(sicht), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-10/Happy, LH-FA-10/Negative, LH-FA-13/Negative — eine
// Anfrage, die nicht zur Aufzeichnung passt, erhält einen eindeutigen Fehler mit
// PGR-E5001 und kein Ergebnis, auch nicht das der aufgezeichneten Anfrage; der
// Lauf endet mit Exit-Code 5.
func TestE2EReplayAbweichung(t *testing.T) {
	input, _ := aufnehmen(t, "SELECT 1;")
	rep := startProzess(t, "replay", "--input", input)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rep.listen))
	if err != nil {
		t.Fatal(err)
	}
	ergebnisse, err := conn.Exec(ctx, "SELECT 2;").ReadAll()
	if err == nil || !strings.Contains(err.Error(), "PGR-E5001") || len(ergebnisse) != 0 {
		t.Fatalf("erwartet PGR-E5001 ohne Ergebnis, erhalten %d Ergebnisse, %v", len(ergebnisse), err)
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

// Abdeckung: LH-FA-03/Happy — bei gestoppter PostgreSQL-Instanz beantwortet
// replay `SELECT 1;` aus der in der ersten Phase von record geschriebenen
// Aufzeichnung mit derselben Sicht wie beim Aufzeichnen; die nie zugeordnete
// zweite Session meldet der Lauf am Ende als Warnung PGR-W2001. Läuft nur in
// der zweiten Phase des Runners (PGR_OHNE_POSTGRES=1).
func TestE2EOhnePostgresReplay(t *testing.T) {
	if os.Getenv("PGR_OHNE_POSTGRES") != "1" {
		t.Skip("nur in der Phase ohne PostgreSQL")
	}
	if c, err := net.DialTimeout("tcp", os.Getenv("PGR_UPSTREAM"), 2*time.Second); err == nil {
		c.Close()
		t.Fatalf("PostgreSQL unter %s ist noch erreichbar", os.Getenv("PGR_UPSTREAM"))
	}
	daten := os.Getenv("PGR_DATEN")
	aufgezeichnet, err := os.ReadFile(filepath.Join(daten, "sicht.txt"))
	if err != nil {
		t.Fatalf("Sicht der ersten Phase: %v", err)
	}
	rep := startProzess(t, "replay", "--input", filepath.Join(daten, "rec.yaml"))
	got := beobachten(t, rep.listen, "SELECT 1;")
	rep.stop(t, 0)
	if got != string(aufgezeichnet) {
		t.Fatalf("Replay-Sicht:\n%s\n--- aufgezeichnet:\n%s", got, aufgezeichnet)
	}
	if !strings.Contains(rep.stderr.String(), "PGR-W2001") {
		t.Fatalf("Warnung PGR-W2001 fehlt:\n%s", rep.stderr.String())
	}
}

// Abdeckung: LH-FA-10/Negative — eine Anfrage nach dem Ende der Aufzeichnung
// ist eine Abweichung: Der Client erhält PGR-E5001 und kein Ergebnis, auch
// nicht das der letzten aufgezeichneten Anfrage; der Lauf endet mit Exit-Code 5
// (LH-FA-10.a).
func TestE2EReplayNachDemEnde(t *testing.T) {
	input, _ := aufnehmen(t, "SELECT 1;")
	rep := startProzess(t, "replay", "--input", input)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rep.listen))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(ctx, "SELECT 1;").ReadAll(); err != nil {
		t.Fatalf("aufgezeichnete Anfrage: %v", err)
	}
	ergebnisse, err := conn.Exec(ctx, "SELECT 1;").ReadAll()
	if err == nil || !strings.Contains(err.Error(), "PGR-E5001") || len(ergebnisse) != 0 {
		t.Fatalf("erwartet PGR-E5001 ohne Ergebnis, erhalten %d Ergebnisse, %v", len(ergebnisse), err)
	}
	_ = conn.Close(ctx)
	rep.stop(t, 5)
}
