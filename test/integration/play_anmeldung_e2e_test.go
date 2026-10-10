//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// anmeldeRolle ist ein Benutzer der Instanz, dessen Anmeldung pg_hba.conf auf
// ein Verfahren festlegt (tools/test/run-integration-tests.sh): Name,
// Verfahren und die Umgebungsvariable, die das Passwort an den Testlauf gibt.
type anmeldeRolle struct {
	name, verfahren, envPasswort string
}

// anmeldeRollen sind die drei Rollen, die der Runner anlegt.
func anmeldeRollen() []anmeldeRolle {
	return []anmeldeRolle{
		{"play_scram", "scram-sha-256", "PGR_PW_SCRAM"},
		{"play_md5", "md5", "PGR_PW_MD5"},
		{"play_pw", "password", "PGR_PW_KLARTEXT"},
	}
}

// passwort liefert das Passwort der Rolle aus dem Testlauf; es enthält GEHEIM.
func (r anmeldeRolle) passwort(t *testing.T) string {
	t.Helper()
	pw := os.Getenv(r.envPasswort)
	if !strings.Contains(pw, "GEHEIM") {
		t.Fatalf("%s fehlt oder enthält kein GEHEIM: das Skript run-integration-tests.sh legt es an", r.envPasswort)
	}
	return pw
}

// anmeldeDatenbank legt die Datenbank name mit der Tabelle spur an, in die jede
// Rolle schreiben darf, und liefert die Verbindung des Administrators zu ihr.
func anmeldeDatenbank(t *testing.T, name string) *pgconn.PgConn {
	t.Helper()
	conn := leereDatenbank(t, name)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := conn.Exec(ctx, "CREATE TABLE spur (wer text, fall text); GRANT INSERT ON spur TO PUBLIC").ReadAll(); err != nil {
		t.Fatalf("spur in %s: %v", name, err)
	}
	return conn
}

// ohneGeheimnis verlangt, dass keine Ausgabe GEHEIM nennt: kein Passwort einer
// Rolle und kein falsches Passwort.
func ohneGeheimnis(t *testing.T, stdout, stderr string) {
	t.Helper()
	if strings.Contains(stdout+stderr, "GEHEIM") {
		t.Errorf("die Ausgabe nennt ein Passwort:\nstdout %q\nstderr:\n%s", stdout, stderr)
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-17/Happy — play meldet sich an einem Server
// mit scram-sha-256, md5 und password in pg_hba.conf an und spielt eine
// Aufzeichnung ein, als der Benutzer der Rolle: Das Passwort kommt aus
// PGWIRE_RECORDER_PASSWORD bei host:port und aus dem Platzhalter der benutzten
// Verbindung, der der Variable vorgeht; in keiner Ausgabe steht ein Passwort
// (LH-FA-20.a *Anmeldung*, *Passwort*, LH-RB-01).
func TestE2EPlayAnmeldung(t *testing.T) {
	conn := anmeldeDatenbank(t, "play_anmeldung")
	host, port, _ := strings.Cut(os.Getenv("PGR_UPSTREAM"), ":")
	dir := t.TempDir()
	datei := fmt.Sprintf("connections:\n  mitpw: \"postgresql://${PGR_E2E_BENUTZER}:${PGR_E2E_PW}@%s:%s/play_anmeldung\"\n", host, port)
	if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte(datei), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, r := range anmeldeRollen() {
		pw := r.passwort(t)
		for _, f := range []struct {
			name string
			args []string
			env  map[string]string
		}{
			{"Variable", []string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--user", r.name, "--database", "play_anmeldung"}, map[string]string{"PGWIRE_RECORDER_PASSWORD": pw}},
			{"Platzhalter vor der Variable", []string{"--upstream", "mitpw", "--log-level", "debug"}, map[string]string{"PGWIRE_RECORDER_PASSWORD": "GEHEIMfalsch", "PGR_E2E_BENUTZER": r.name, "PGR_E2E_PW": pw}},
		} {
			t.Run(r.verfahren+" "+f.name, func(t *testing.T) {
				for k, v := range f.env {
					t.Setenv(k, v)
				}
				fall := r.name + " " + f.name
				input := aufzeichnungMit(t, "postgres", "postgres", "INSERT INTO spur VALUES (current_user, '"+fall+"')")
				stdout, stderr, code := starteBis(t, dir, append([]string{"play", "--input", input}, f.args...)...)
				if code != 0 || strings.Contains(stderr, "level=ERROR") {
					t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
				}
				ohneGeheimnis(t, stdout, stderr)
				if gefunden := spurVon(t, map[string]*pgconn.PgConn{"play_anmeldung": conn}, fall); len(gefunden) != 1 || gefunden[0] != r.name+"/play_anmeldung" {
					t.Fatalf("eingespielt als %v, erwartet %s", gefunden, r.name)
				}
			})
		}
	}
}

// Abdeckung: LH-FA-20/Boundary — verlangt der Server kein Passwort (trust),
// sendet play keines, auch wenn PGWIRE_RECORDER_PASSWORD gesetzt ist: Ein
// unverlangtes Passwort wäre für den Server eine ungültige Nachricht, und das
// Einspielen gelänge nicht (LH-FA-20.a *Anmeldung*, *Passwort*).
func TestE2EPlayAnmeldungOhneVerlangen(t *testing.T) {
	conn := anmeldeDatenbank(t, "play_anmeldung_trust")
	t.Setenv("PGWIRE_RECORDER_PASSWORD", "GEHEIMunverlangt")
	input := aufzeichnungMit(t, "postgres", "play_anmeldung_trust", "INSERT INTO spur VALUES (current_user, 'trust')")
	stdout, stderr, code := starteBis(t, "", "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--input", input)
	if code != 0 || strings.Contains(stderr, "level=ERROR") {
		t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
	}
	ohneGeheimnis(t, stdout, stderr)
	if gefunden := spurVon(t, map[string]*pgconn.PgConn{"play_anmeldung_trust": conn}, "trust"); len(gefunden) != 1 {
		t.Fatalf("eingespielt: %v", gefunden)
	}
}

// Abdeckung: LH-FA-20/Negative — ein falsches Passwort (SQLSTATE-Klasse 28) und
// ein fehlendes Passwort, wenn der Server eines verlangt, sind PGR-E4005 mit
// Exit-Code 4 und spielen nichts ein, nach scram-sha-256, md5 und password
// gleich; eine leere PGWIRE_RECORDER_PASSWORD gilt als nicht gesetzt; weder
// das richtige noch das falsche Passwort steht in einer Ausgabe (LH-FA-20.a
// *Anmeldung*, *Passwort*, LH-QA-05, LH-RB-01).
func TestE2EPlayAnmeldungFehler(t *testing.T) {
	conn := anmeldeDatenbank(t, "play_anmeldung_fehler")
	for _, r := range anmeldeRollen() {
		for _, f := range []struct {
			name, passwort, meldung string
		}{
			{"falsches Passwort", "GEHEIMfalsch", "28P01"},
			{"leeres Passwort", "", "kein"},
			{"Passwort mit dem Anfang des richtigen", r.passwort(t)[:len(r.passwort(t))-1], "28P01"},
		} {
			t.Run(r.verfahren+" "+f.name, func(t *testing.T) {
				t.Setenv("PGWIRE_RECORDER_PASSWORD", f.passwort)
				fall := r.name + " " + f.name
				input := aufzeichnungMit(t, "postgres", "postgres", "INSERT INTO spur VALUES (current_user, '"+fall+"')")
				stdout, stderr, code := starteBis(t, "", "play", "--upstream", os.Getenv("PGR_UPSTREAM"), "--user", r.name, "--database", "play_anmeldung_fehler", "--input", input)
				if code != 4 || !strings.Contains(stderr, "code=PGR-E4005") || !strings.Contains(stderr, f.meldung) {
					t.Fatalf("Exit-Code %d, erwartet 4 mit PGR-E4005 und %q, stderr:\n%s", code, f.meldung, stderr)
				}
				ohneGeheimnis(t, stdout, stderr)
				if gefunden := spurVon(t, map[string]*pgconn.PgConn{"play_anmeldung_fehler": conn}, fall); len(gefunden) != 0 {
					t.Fatalf("trotz Fehler der Anmeldung eingespielt: %v", gefunden)
				}
			})
		}
	}
}
