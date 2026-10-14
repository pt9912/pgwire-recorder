//go:build integration

package integration_test

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// dsnMit ist die Verbindung zu address als Benutzer name mit dem Passwort pw.
func dsnMit(address, name, pw, optionen string) string {
	u := url.URL{Scheme: "postgres", User: url.UserPassword(name, pw), Host: address, Path: "/postgres", RawQuery: "sslmode=disable&connect_timeout=10" + optionen}
	return u.String()
}

// Abdeckung: LH-FA-05/Happy, LH-FA-02/Happy, LH-RB-01/Boundary — ein Client
// meldet sich über record an einem Server mit scram-sha-256, md5 und password in
// pg_hba.conf an und führt eine Anfrage aus; die Aufzeichnung trägt die Anfrage
// und den Benutzer der Startnachricht, aber weder Passwort noch Anmeldenachrichten,
// und auch das Log der Stufe debug nennt kein Passwort (LH-FA-05.b *Vermittlung*,
// *Aufzeichnung*, *Log und Meldungen*).
func TestE2ERecordAnmeldung(t *testing.T) {
	for _, r := range anmeldeRollen() {
		t.Run(r.verfahren, func(t *testing.T) {
			pw := r.passwort(t)
			output := filepath.Join(t.TempDir(), "rec.yaml")
			rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--log-level", "debug")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			conn, err := pgconn.Connect(ctx, dsnMit(rec.listen, r.name, pw, ""))
			if err != nil {
				t.Fatalf("Anmeldung über den Recorder: %v", err)
			}
			results, err := conn.Exec(ctx, "SELECT current_user;").ReadAll()
			if err != nil || len(results) != 1 || string(results[0].Rows[0][0]) != r.name {
				t.Fatalf("Anfrage: %v, %#v", err, results)
			}
			_ = conn.Close(ctx)
			rec.stop(t, 0)
			text := lies(t, output)
			for _, zeile := range []string{"sql: SELECT current_user;", "user: " + r.name} {
				if !strings.Contains(text, zeile) {
					t.Errorf("Aufzeichnung enthält %q nicht:\n%s", zeile, text)
				}
			}
			ohneGeheimnis(t, text, rec.stderr.String())
		})
	}
}

// Abdeckung: LH-FA-05/Negative — eine Fehlerantwort des Servers in der Anmeldung
// (falsches Passwort nach scram-sha-256, md5 und password, und eine Zugriffsregel
// reject ohne Anforderung) geht mit allen Feldern unverändert an den Client, wie
// bei einer Verbindung unmittelbar zum Server; der Lauf meldet nichts, endet mit
// Exit-Code 0 und schreibt keine Session (LH-FA-05.b *Fehlerantwort des Servers*).
func TestE2ERecordAnmeldungFehlerantwort(t *testing.T) {
	rollen := append(anmeldeRollen(), anmeldeRolle{name: "play_reject", verfahren: "reject"})
	for _, r := range rollen {
		t.Run(r.verfahren, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "rec.yaml")
			rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-empty-sessions")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			var ueber, direkt *pgconn.PgError
			_, errUeber := pgconn.Connect(ctx, dsnMit(rec.listen, r.name, "GEHEIMfalsch", ""))
			_, errDirekt := pgconn.Connect(ctx, dsnMit(os.Getenv("PGR_UPSTREAM"), r.name, "GEHEIMfalsch", ""))
			if !errors.As(errUeber, &ueber) || !errors.As(errDirekt, &direkt) || !reflect.DeepEqual(ueber, direkt) {
				t.Fatalf("über den Recorder %v, direkt %v", errUeber, errDirekt)
			}
			rec.stop(t, 0)
			text := lies(t, output)
			if !strings.Contains(text, "sessions: []") {
				t.Errorf("erwartet eine Aufzeichnung ohne Session:\n%s", text)
			}
			ohneGeheimnis(t, text, rec.stderr.String())
			if strings.Contains(rec.stderr.String(), "level=ERROR") || strings.Contains(rec.stderr.String(), "PGR-") {
				t.Errorf("der Lauf meldet etwas:\n%s", rec.stderr.String())
			}
		})
	}
}

// Abdeckung: LH-FA-05/Boundary — beendet der Client die Anmeldung selbst (ein
// Client, der Channel Binding verlangt, bekommt es nicht angeboten), ist das kein
// Verbindungsfehler: keine Meldung, Exit-Code 0, auch mit --record-empty-sessions
// keine Session (LH-FA-05.b *Verfahren*, *Ende durch den Client*).
func TestE2ERecordAnmeldungClientBeendet(t *testing.T) {
	r := anmeldeRollen()[0]
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--record-empty-sessions")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := pgconn.Connect(ctx, dsnMit(rec.listen, r.name, r.passwort(t), "&channel_binding=require"))
	if err == nil || !strings.Contains(err.Error(), "channel binding required") {
		t.Fatalf("erwartet den Abbruch des Clients, erhalten %v", err)
	}
	rec.stop(t, 0)
	text := lies(t, output)
	if !strings.Contains(text, "sessions: []") {
		t.Errorf("erwartet eine Aufzeichnung ohne Session:\n%s", text)
	}
	ohneGeheimnis(t, text, rec.stderr.String())
	if strings.Contains(rec.stderr.String(), "level=ERROR") || strings.Contains(rec.stderr.String(), "PGR-") {
		t.Errorf("der Lauf meldet etwas:\n%s", rec.stderr.String())
	}
}
