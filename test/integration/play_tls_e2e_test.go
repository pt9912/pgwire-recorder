//go:build integration

package integration_test

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pgwire-recorder/internal/bootstrap/tlsproxy"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// zertifikatFuerLokal ist ein Serverzertifikat der Zertifizierungsstelle ca auf
// die Adresse 127.0.0.1.
func zertifikatFuerLokal(t *testing.T, ca *testpki.CA) testpki.Blattzertifikat {
	t.Helper()
	return ca.Server(t, testpki.Blatt{IPs: []net.IP{net.ParseIP("127.0.0.1")}})
}

// caDatei schreibt das Zertifikat von ca als ca.pem in dir und liefert den
// Pfad.
func caDatei(t *testing.T, dir string, ca *testpki.CA) string {
	t.Helper()
	pfad := filepath.Join(dir, "ca.pem")
	if err := os.WriteFile(pfad, ca.PEM(), 0o600); err != nil {
		t.Fatal(err)
	}
	return pfad
}

// verzeichnisMit legt ein Verzeichnis mit der Konfigurationsdatei datei (leer:
// keine) an und liefert es.
func verzeichnisMit(t *testing.T, datei string) string {
	t.Helper()
	dir := t.TempDir()
	if datei != "" {
		if err := os.WriteFile(filepath.Join(dir, ".pgwire-recorder.yaml"), []byte(datei), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// Abdeckung: LH-FA-20/Happy, LH-FA-17/Happy, LH-FA-17/Boundary — play baut die
// Verbindung zum Server mit TLS auf und prüft das Zertifikat des Servers gegen
// die Zertifizierungsstelle aus --upstream-ca (Pfad absolut oder relativ zum
// aktuellen Verzeichnis), wenn --upstream-tls gesetzt ist (Option,
// Umgebungsvariable oder Schlüssel) und wenn die benutzte Verbindung
// sslmode=require verlangt: Die Eingespielten Anfragen wirken in der Datenbank,
// der Proxy sah genau eine TLS-Verbindung je Session, keine ohne TLS, und der
// Name 127.0.0.1 passt zur IP-Adresse des Zertifikats. Das Zertifikat der
// Zertifizierungsstelle liegt nicht im Speicher des Systems; gegen die Instanz
// ohne TLS ist die Verbindung nicht möglich (LH-FA-20.a *TLS*, LH-FA-17.a
// *Wirkung einer URL*).
func TestE2EPlayTLS(t *testing.T) {
	conn := anmeldeDatenbank(t, "play_tls")
	ca := testpki.NeueCA(t, "E2E-CA")
	proxy := tlsproxy.Start(t, zertifikatFuerLokal(t, ca), os.Getenv("PGR_UPSTREAM"), false)
	host, port, _ := strings.Cut(proxy.Addr, ":")
	verbindung := fmt.Sprintf("connections:\n  tls: \"postgresql://postgres@%s:%s/play_tls?sslmode=require\"\n", host, port)
	for _, f := range []struct {
		name  string
		datei string
		env   map[string]string
		args  []string
		sess  int
	}{
		{"Optionen", "", nil, []string{"--upstream", proxy.Addr, "--database", "play_tls", "--upstream-tls", "--upstream-ca", "CA"}, 1},
		{"sslmode=require", verbindung, nil, []string{"--upstream", "tls", "--upstream-ca", "CA"}, 1},
		{"Umgebungsvariablen", "", map[string]string{"PGWIRE_RECORDER_UPSTREAM_TLS": "true", "PGWIRE_RECORDER_UPSTREAM_CA": "CA"}, []string{"--upstream", proxy.Addr, "--database", "play_tls"}, 1},
		{"Schlüssel mit relativem Pfad", "play:\n  upstream: " + proxy.Addr + "\n  database: play_tls\n  upstream_tls: true\n  upstream_ca: ca.pem\n", nil, nil, 1},
		{"sslmode=require, Option TLS", verbindung, nil, []string{"--upstream", "tls", "--upstream-tls=true", "--upstream-ca", "ca.pem"}, 1},
	} {
		t.Run(f.name, func(t *testing.T) {
			dir := verzeichnisMit(t, f.datei)
			pfad := caDatei(t, dir, ca)
			for k, v := range f.env {
				t.Setenv(k, strings.ReplaceAll(v, "CA", pfad))
			}
			args := make([]string, len(f.args))
			for i, a := range f.args {
				args[i] = strings.ReplaceAll(a, "CA", pfad)
			}
			vorher := proxy.TLSOk()
			input := aufzeichnungMit(t, "postgres", "play_tls", "INSERT INTO spur VALUES (current_user, '"+f.name+"')")
			_, stderr, code := starteBis(t, dir, append([]string{"play", "--input", input}, args...)...)
			if code != 0 || strings.Contains(stderr, "level=ERROR") {
				t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
			}
			if got := proxy.TLSOk() - vorher; got != f.sess || proxy.Klartext() != 0 {
				t.Fatalf("TLS-Verbindungen %d, erwartet %d; ohne TLS %d", got, f.sess, proxy.Klartext())
			}
			if gefunden := spurVon(t, map[string]*pgconn.PgConn{"play_tls": conn}, f.name); len(gefunden) != 1 {
				t.Fatalf("eingespielt: %v", gefunden)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Happy, LH-RB-01/Messung — nach der Aushandlung laufen
// Aufbau und Anmeldung unverändert auf der verschlüsselten Verbindung: Über den
// TLS-Proxy meldet sich play an der Instanz mit scram-sha-256, md5 und password
// an (das Klartext-Passwort geht über TLS) und spielt als der Benutzer der Rolle
// ein; in keiner Ausgabe steht ein Passwort, und der Proxy sah je Lauf genau
// eine TLS-Verbindung (LH-FA-20.a *TLS*, *Anmeldung*).
func TestE2EPlayTLSAnmeldung(t *testing.T) {
	conn := anmeldeDatenbank(t, "play_tls_anm")
	ca := testpki.NeueCA(t, "E2E-CA")
	proxy := tlsproxy.Start(t, zertifikatFuerLokal(t, ca), os.Getenv("PGR_UPSTREAM"), false)
	dir := t.TempDir()
	pfad := caDatei(t, dir, ca)
	for _, r := range anmeldeRollen() {
		t.Run(r.verfahren, func(t *testing.T) {
			t.Setenv("PGWIRE_RECORDER_PASSWORD", r.passwort(t))
			vorher := proxy.TLSOk()
			fall := r.name + " tls"
			input := aufzeichnungMit(t, "postgres", "play_tls_anm", "INSERT INTO spur VALUES (current_user, '"+fall+"')")
			stdout, stderr, code := starteBis(t, dir, "play", "--input", input, "--upstream", proxy.Addr, "--user", r.name, "--database", "play_tls_anm", "--upstream-tls", "--upstream-ca", pfad)
			if code != 0 || strings.Contains(stderr, "level=ERROR") {
				t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
			}
			ohneGeheimnis(t, stdout, stderr)
			if got := proxy.TLSOk() - vorher; got != 1 || proxy.Klartext() != 0 {
				t.Fatalf("TLS-Verbindungen %d, erwartet 1; ohne TLS %d", got, proxy.Klartext())
			}
			if gefunden := spurVon(t, map[string]*pgconn.PgConn{"play_tls_anm": conn}, fall); len(gefunden) != 1 || gefunden[0] != r.name+"/play_tls_anm" {
				t.Fatalf("eingespielt als %v, erwartet %s", gefunden, r.name)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-17/Boundary — antwortet
// der Server auf das SSLRequest mit `N` (die reale Instanz ohne ssl), ist das
// PGR-E4005 mit Exit-Code 4 und spielt nichts ein, gleich ob TLS aus
// --upstream-tls (Option, Umgebungsvariable, Schlüssel) oder aus
// sslmode=require folgt; ein ausdrücklich gesetztes --upstream-tls=false (aus
// Option, Umgebungsvariable und Schlüssel) geht sslmode=require vor, und play
// verbindet ohne TLS und spielt ein (LH-FA-20.a *TLS*, LH-FA-17.a *Wirkung
// einer URL*).
func TestE2EPlayTLSAbgelehnt(t *testing.T) {
	conn := anmeldeDatenbank(t, "play_tls_n")
	host, port, _ := strings.Cut(os.Getenv("PGR_UPSTREAM"), ":")
	verbindung := fmt.Sprintf("connections:\n  tls: \"postgresql://postgres@%s:%s/play_tls_n?sslmode=require\"\n", host, port)
	for _, f := range []struct {
		name   string
		datei  string
		env    map[string]string
		args   []string
		spielt bool
	}{
		{"Option", "", nil, []string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--database", "play_tls_n", "--upstream-tls"}, false},
		{"Umgebungsvariable", "", map[string]string{"PGWIRE_RECORDER_UPSTREAM_TLS": "true"}, []string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--database", "play_tls_n"}, false},
		{"Schlüssel", "play:\n  upstream_tls: true\n", nil, []string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--database", "play_tls_n"}, false},
		{"sslmode=require", verbindung, nil, []string{"--upstream", "tls"}, false},
		{"require, Option false", verbindung, nil, []string{"--upstream", "tls", "--upstream-tls=false"}, true},
		{"require, Umgebungsvariable false", verbindung, map[string]string{"PGWIRE_RECORDER_UPSTREAM_TLS": "false"}, []string{"--upstream", "tls"}, true},
		{"require, Schlüssel false", verbindung + "play:\n  upstream_tls: false\n", nil, []string{"--upstream", "tls"}, true},
	} {
		t.Run(f.name, func(t *testing.T) {
			dir := verzeichnisMit(t, f.datei)
			for k, v := range f.env {
				t.Setenv(k, v)
			}
			input := aufzeichnungMit(t, "postgres", "play_tls_n", "INSERT INTO spur VALUES (current_user, '"+f.name+"')")
			_, stderr, code := starteBis(t, dir, append([]string{"play", "--input", input}, f.args...)...)
			gefunden := spurVon(t, map[string]*pgconn.PgConn{"play_tls_n": conn}, f.name)
			if f.spielt {
				if code != 0 || strings.Contains(stderr, "level=ERROR") || len(gefunden) != 1 {
					t.Fatalf("Exit-Code %d, eingespielt %v, stderr:\n%s", code, gefunden, stderr)
				}
				return
			}
			if code != 4 || !strings.Contains(stderr, "code=PGR-E4005") || !strings.Contains(stderr, "lehnt TLS ab") || len(gefunden) != 0 {
				t.Fatalf("Exit-Code %d, eingespielt %v, stderr:\n%s", code, gefunden, stderr)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Negative — ein abgelaufenes Zertifikat,
// eines auf eine andere IP-Adresse, eines auf den Namen localhost statt der
// IP-Adresse, eines einer Zertifizierungsstelle, die --upstream-ca nicht nennt,
// und das richtige Zertifikat ohne --upstream-ca (es liegt nicht im Speicher des
// Systems) sind PGR-E4005 mit Exit-Code 4: Der Proxy reicht nichts an die Instanz
// weiter, nichts wird eingespielt, und die Meldung nennt den Grund (LH-FA-20.a
// *TLS*).
func TestE2EPlayTLSZertifikatFehler(t *testing.T) {
	conn := anmeldeDatenbank(t, "play_tls_z")
	ca, andere := testpki.NeueCA(t, "E2E-CA"), testpki.NeueCA(t, "E2E-Andere")
	for _, f := range []struct {
		name       string
		aussteller *testpki.CA
		blatt      testpki.Blatt
		mitCA      bool
		grund      string
	}{
		{"abgelaufen", ca, testpki.Blatt{IPs: []net.IP{net.ParseIP("127.0.0.1")}, Von: time.Now().Add(-2 * time.Hour), Bis: time.Now().Add(-time.Hour)}, true, "expired"},
		{"andere IP-Adresse", ca, testpki.Blatt{IPs: []net.IP{net.ParseIP("10.1.2.3")}}, true, "not 127.0.0.1"},
		{"Name statt IP-Adresse", ca, testpki.Blatt{DNS: []string{"localhost"}}, true, "any IP SANs"},
		{"andere Zertifizierungsstelle", andere, testpki.Blatt{IPs: []net.IP{net.ParseIP("127.0.0.1")}}, true, "unknown authority"},
		{"ohne --upstream-ca", ca, testpki.Blatt{IPs: []net.IP{net.ParseIP("127.0.0.1")}}, false, "unknown authority"},
	} {
		t.Run(f.name, func(t *testing.T) {
			proxy := tlsproxy.Start(t, f.aussteller.Server(t, f.blatt), os.Getenv("PGR_UPSTREAM"), false)
			dir := t.TempDir()
			args := []string{"play", "--upstream", proxy.Addr, "--database", "play_tls_z", "--upstream-tls"}
			if f.mitCA {
				args = append(args, "--upstream-ca", caDatei(t, dir, ca))
			}
			input := aufzeichnungMit(t, "postgres", "play_tls_z", "INSERT INTO spur VALUES (current_user, '"+f.name+"')")
			_, stderr, code := starteBis(t, dir, append(args, "--input", input)...)
			if code != 4 || !strings.Contains(stderr, "code=PGR-E4005") || !strings.Contains(stderr, f.grund) {
				t.Fatalf("Exit-Code %d, erwartet 4 mit PGR-E4005 und %q, stderr:\n%s", code, f.grund, stderr)
			}
			if proxy.Weitergabe() != 0 || len(spurVon(t, map[string]*pgconn.PgConn{"play_tls_z": conn}, f.name)) != 0 {
				t.Fatalf("trotz Fehler weitergereicht oder eingespielt")
			}
		})
	}
}

// Abdeckung: LH-FA-20/Negative — lehnt der Server eine
// unverschlüsselte Verbindung ab (der Proxy antwortet auf das Startup wie ein
// Server mit hostssl mit der Fehlerantwort 28000), ist das PGR-E4005 mit
// Exit-Code 4, und es entsteht keine TLS-Verbindung (LH-FA-20.a Tabelle
// *Fehlerregeln beim Einspielen*).
func TestE2EPlayKlartextAbgelehnt(t *testing.T) {
	ca := testpki.NeueCA(t, "E2E-CA")
	proxy := tlsproxy.Start(t, zertifikatFuerLokal(t, ca), os.Getenv("PGR_UPSTREAM"), false)
	input := aufzeichnungMit(t, "postgres", "postgres", "SELECT 1")
	_, stderr, code := starteBis(t, "", "play", "--upstream", proxy.Addr, "--input", input)
	if code != 4 || !strings.Contains(stderr, "code=PGR-E4005") || !strings.Contains(stderr, "28000") {
		t.Fatalf("Exit-Code %d, stderr:\n%s", code, stderr)
	}
	if proxy.Klartext() != 1 || proxy.TLSOk() != 0 {
		t.Fatalf("Verbindungen ohne TLS %d, mit TLS %d", proxy.Klartext(), proxy.TLSOk())
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-17/Negative, LH-RB-01/Messung —
// --upstream-ca ohne TLS ist PGR-E2001; eine Datei aus --upstream-ca, die
// keine reguläre Datei ist, keinen PEM-Block, einen Block eines anderen Typs
// oder kein lesbares X.509-Zertifikat enthält, ist PGR-E2007: beides Startfehler
// mit Exit-Code 2 als Zeile beim Prozessende, vor dem Laden der Aufzeichnung (eine
// fehlende Aufzeichnung ist nicht PGR-E3001) und vor jeder Verbindung; die
// Meldung nennt weder Pfad noch Inhalt (LH-FA-20.a *Start*, LH-FA-17.a *Fehler*).
func TestE2EPlayUpstreamCAFehler(t *testing.T) {
	ca := testpki.NeueCA(t, "E2E-CA")
	proxy := tlsproxy.Start(t, zertifikatFuerLokal(t, ca), os.Getenv("PGR_UPSTREAM"), false)
	dir := t.TempDir()
	gueltig := caDatei(t, dir, ca)
	schreibe := func(name, inhalt string) string {
		pfad := filepath.Join(dir, name)
		if err := os.WriteFile(pfad, []byte(inhalt), 0o600); err != nil {
			t.Fatal(err)
		}
		return pfad
	}
	verz := filepath.Join(dir, "GEHEIMVERZ")
	if err := os.Mkdir(verz, 0o700); err != nil {
		t.Fatal(err)
	}
	fehlt := filepath.Join(dir, "GEHEIMFEHLT.pem")
	fehlt2001 := []string{"PGR-E2001", "--upstream-ca"}
	for _, f := range []struct {
		name    string
		args    []string
		code    int
		meldung []string
	}{
		{"ohne TLS", []string{"--upstream-ca", gueltig}, 2, fehlt2001},
		{"ohne TLS, Datei ungültig", []string{"--upstream-ca", schreibe("GEHEIMkein.pem", "GEHEIMINHALT")}, 2, fehlt2001},
		{"Datei ohne PEM-Block", []string{"--upstream-tls", "--upstream-ca", schreibe("GEHEIMkein.pem", "GEHEIMINHALT")}, 2, []string{"PGR-E2007", "--upstream-ca"}},
		{"Block eines anderen Typs", []string{"--upstream-tls", "--upstream-ca", schreibe("GEHEIMschluessel.pem", "-----BEGIN PRIVATE KEY-----\nR0VIRUlNSU5IQUxU\n-----END PRIVATE KEY-----\n")}, 2, []string{"PGR-E2007", "Typ CERTIFICATE"}},
		{"Block ohne Zertifikat", []string{"--upstream-tls", "--upstream-ca", schreibe("GEHEIMleer.pem", "-----BEGIN CERTIFICATE-----\nR0VIRUlNSU5IQUxU\n-----END CERTIFICATE-----\n")}, 2, []string{"PGR-E2007", "X.509"}},
		{"Verzeichnis", []string{"--upstream-tls", "--upstream-ca", verz}, 2, []string{"PGR-E2007", "keine reguläre Datei"}},
		{"fehlende Datei", []string{"--upstream-tls", "--upstream-ca", fehlt}, 2, []string{"PGR-E2007", "nicht vorhanden"}},
	} {
		t.Run(f.name, func(t *testing.T) {
			args := append([]string{"play", "--upstream", proxy.Addr, "--input", filepath.Join(dir, "fehlt.yaml")}, f.args...)
			stdout, stderr, code := starteBis(t, dir, args...)
			if code != f.code || stdout != "" || strings.Contains(stderr, "level=") {
				t.Fatalf("Exit-Code %d, stdout %q, stderr:\n%s", code, stdout, stderr)
			}
			for _, m := range f.meldung {
				if !strings.Contains(stderr, m) {
					t.Fatalf("stderr ohne %q:\n%s", m, stderr)
				}
			}
			if strings.Contains(stderr, "GEHEIM") || strings.Contains(stderr, dir) || strings.Contains(stderr, "PGR-E3001") {
				t.Fatalf("stderr nennt Pfad oder Inhalt oder lädt die Aufzeichnung:\n%s", stderr)
			}
			if proxy.TLSOk() != 0 || proxy.Klartext() != 0 {
				t.Fatalf("Startfehler verbindet sich")
			}
		})
	}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-14/Boundary — das zweite SIGINT oder
// SIGTERM, während die Aushandlung von TLS wartet (der Server antwortet `S` und
// schweigt), schließt die Verbindung: Das erste Signal lässt den Aufbau weiterlaufen,
// play endet binnen 10 s nach dem zweiten mit Exit-Code 0, der unterbrochene
// Aufbau ist kein Fehler, und es steht keine Zeile error im Log (LH-FA-20.a
// *Abbruchsignal*, *TLS*).
func TestE2EPlayTLSAbbruch(t *testing.T) {
	ca := testpki.NeueCA(t, "E2E-CA")
	proxy := tlsproxy.Start(t, zertifikatFuerLokal(t, ca), os.Getenv("PGR_UPSTREAM"), true)
	dir := t.TempDir()
	input := aufzeichnungMit(t, "postgres", "postgres", "SELECT 1")
	stderr := &stderrMitSignal{abbruch: make(chan struct{})}
	cmd := exec.Command(os.Getenv("PGR_BINARY"), "play", "--upstream", proxy.Addr, "--input", input, "--upstream-tls", "--upstream-ca", caDatei(t, dir, ca))
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	beendet := make(chan error, 1)
	go func() { beendet <- cmd.Wait() }()
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	select {
	case <-proxy.Hello:
	case <-time.After(30 * time.Second):
		t.Fatalf("das ClientHello trifft binnen 30 s nicht ein:\n%s", stderr.String())
	}
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("erstes Signal: %v", err)
	}
	select {
	case <-stderr.abbruch:
	case <-time.After(30 * time.Second):
		t.Fatalf("play schreibt binnen 30 s nach dem ersten Signal keine Zeile zum Abbruchsignal:\n%s", stderr.String())
	}
	select {
	case err := <-beendet:
		t.Fatalf("play endet nach dem ersten Signal mitten in der Aushandlung: %v\n%s", err, stderr.String())
	case <-time.After(1 * time.Second):
	}
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("zweites Signal: %v", err)
	}
	select {
	case err := <-beendet:
		zeilen := logZeilen(t, stderr.String())
		if code := exitCodeOf(err); code != 0 || form(zeilen) != "INFO play gestartet | INFO Abbruchsignal, play endet vorzeitig | INFO play beendet" {
			t.Fatalf("Exit-Code %d, Zeilen %s:\n%s", code, form(zeilen), stderr.String())
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("play endet binnen 10 s nach dem zweiten Signal nicht:\n%s", stderr.String())
	}
}
