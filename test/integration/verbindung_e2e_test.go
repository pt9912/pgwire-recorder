//go:build integration

package integration_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// verbindungsDatei legt eine Konfigurationsdatei mit inhalt im Verzeichnis
// des Tests an und liefert ihren Pfad.
func verbindungsDatei(t *testing.T, inhalt string) string {
	t.Helper()
	pfad := filepath.Join(t.TempDir(), "konfiguration.yaml")
	if err := os.WriteFile(pfad, []byte(inhalt), 0o600); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// Abdeckung: LH-FA-17/Happy, LH-FA-17/Boundary — record mit --upstream und dem
// Namen einer Verbindung, deren Host und Port aus Platzhaltern kommen, verbindet
// zu Host und Port der Verbindung: Ein Client erhält das Ergebnis von
// `SELECT 1;`, und die Zeile beim Start nennt upstream= mit der eingesetzten
// Adresse, nicht den Namen und weder Benutzer, Passwort noch Datenbank der URL;
// deren Variablen bleiben bei record unbeachtet.
func TestE2ERecordVerbindungAusDatei(t *testing.T) {
	host, port, err := net.SplitHostPort(os.Getenv("PGR_UPSTREAM"))
	if err != nil {
		t.Fatalf("PGR_UPSTREAM: %v", err)
	}
	t.Setenv("PGR_E2E_HOST", host)
	t.Setenv("PGR_E2E_PORT", port)
	t.Setenv("PGR_E2E_PW", "GEHEIMPW")
	t.Setenv("PGR_E2E_DB", "")
	pfad := verbindungsDatei(t, "connections:\n  staging: \"postgresql://NUTZERX:${PGR_E2E_PW}@${PGR_E2E_HOST}:${PGR_E2E_PORT}/DBNAMEX${PGR_E2E_DB}\"\n")
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", "staging", "--output", output, "--config", pfad)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatalf("Verbindung über den Recorder: %v", err)
	}
	results, err := conn.Exec(ctx, "SELECT 1;").ReadAll()
	if err != nil || len(results) != 1 || len(results[0].Rows) != 1 || string(results[0].Rows[0][0]) != "1" {
		t.Fatalf("SELECT 1 über den Recorder: %#v, %v", results, err)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatalf("Verbindung schließen: %v", err)
	}
	rec.stop(t, 0)

	log := rec.stderr.String()
	if !strings.Contains(log, "upstream="+host+":"+port+"\n") {
		t.Errorf("Zeile beim Start ohne upstream=%s:%s:\n%s", host, port, log)
	}
	for _, fremd := range []string{"staging", "NUTZERX", "GEHEIMPW", "DBNAMEX"} {
		if strings.Contains(log, fremd) {
			t.Errorf("stderr nennt %q:\n%s", fremd, log)
		}
	}
	if !strings.Contains(lies(t, output), "sql: SELECT 1;") {
		t.Errorf("Aufzeichnung ohne SELECT 1:\n%s", lies(t, output))
	}
}

// Abdeckung: LH-FA-17/Boundary — record mit einer Verbindung, deren Host eine
// IPv6-Adresse in eckigen Klammern ist, verbindet zu ihr wieder in eckigen
// Klammern: Die Zeile beim Start nennt upstream=[::1]:1, der Client erhält
// PGR-E4002 mit dieser Adresse, und der Lauf endet beim Beenden mit Exit-Code 4.
func TestE2ERecordVerbindungIPv6(t *testing.T) {
	pfad := verbindungsDatei(t, "connections:\n  sechs: \"postgresql://[::1]:1/db\"\n")
	rec := startProzess(t, "record", "--upstream", "sechs", "--output", filepath.Join(t.TempDir(), "rec.yaml"), "--config", pfad)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err == nil || !strings.Contains(err.Error(), codeUpstream) || !strings.Contains(err.Error(), "[::1]:1") {
		t.Fatalf("erwartet %s mit [::1]:1, erhalten %v", codeUpstream, err)
	}
	rec.stop(t, 4)
	if log := rec.stderr.String(); !strings.Contains(log, "upstream=[::1]:1\n") {
		t.Errorf("Zeile beim Start ohne upstream=[::1]:1:\n%s", log)
	}
}
