//go:build integration

package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgconn/ctxwatch"
)

// Abdeckung: LH-FA-13/Happy, LH-FA-05/Happy — pgx bricht den Kontext einer laufenden
// Anfrage über `record` ab und sendet dafür einen CancelRequest ohne Schlüssel (der Recorder
// gibt BackendKeyData nicht weiter, Länge 12): Der Recorder schließt diese Verbindung
// mit PGR-W3001, ohne PGR-E6001, die Anfrage auf der Hauptverbindung endet regulär (der
// Server beantwortet sie vor der Frist des Clients von 20 s), die Hauptverbindung sendet
// Terminate, und der Lauf endet nach SIGTERM mit Exit-Code 0.
func TestE2ERecordCancelOhneSchluessel(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg, err := pgx.ParseConfig(dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	cfg.BuildContextWatcherHandler = func(pc *pgconn.PgConn) ctxwatch.Handler {
		return &pgconn.CancelRequestContextWatcherHandler{Conn: pc, DeadlineDelay: 20 * time.Second}
	}
	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("Verbindung über den Recorder: %v", err)
	}

	anfrage, abbruch := context.WithTimeout(ctx, 100*time.Millisecond)
	_, err = conn.Exec(anfrage, "SELECT pg_sleep(0.5)", pgx.QueryExecModeSimpleProtocol)
	abbruch()
	t.Logf("Anfrage mit abgebrochenem Kontext: %v", err)
	if conn.IsClosed() {
		t.Fatalf("Hauptverbindung vom Client geschlossen, Anfrage: %v", err)
	}
	var n int
	if err := conn.QueryRow(ctx, "SELECT 7", pgx.QueryExecModeSimpleProtocol).Scan(&n); err != nil || n != 7 {
		t.Fatalf("Anfrage nach dem Abbruch: %d, %v", n, err)
	}
	if err := conn.Close(ctx); err != nil {
		t.Fatal(err)
	}
	rec.stop(t, 0)

	log := rec.stderr.String()
	if !strings.Contains(log, "code=PGR-W3001") || strings.Contains(log, "PGR-E6001") || strings.Contains(log, "PGR-E4003") {
		t.Fatalf("Log ohne PGR-W3001 oder mit PGR-E6001 oder PGR-E4003:\n%s", log)
	}
	if c := strings.Count(log, "code=PGR-W3001"); c != 1 {
		t.Fatalf("%d Warnungen PGR-W3001 statt einer:\n%s", c, log)
	}
	laedt(t, output)
}
