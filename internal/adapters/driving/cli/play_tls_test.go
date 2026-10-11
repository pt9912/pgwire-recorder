package cli_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// tlsDatei ist der Text einer Konfigurationsdatei mit der Verbindung v mit
// dem sslmode (leer: ohne Parameter) und dem Abschnitt play: plan.
func tlsDatei(sslmode, play string) string {
	url := "postgresql://u@h/db"
	if sslmode != "" {
		url += "?sslmode=" + sslmode
	}
	return mitVerbindung(url) + play
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary, LH-FA-17/Happy — play nutzt
// TLS zum Server, wenn --upstream-tls gesetzt ist (Option, Umgebungsvariable
// oder Schlüssel), sonst wenn die benutzte Verbindung sslmode=require
// verlangt; ein gesetztes --upstream-tls geht sslmode vor, auch mit false und
// aus jeder Quelle; ohne beides, mit sslmode=disable und bei host:port gilt
// kein TLS; die Kommandozeile geht der Umgebungsvariable vor und diese dem
// Schlüssel (LH-FA-17.a *Wirkung einer URL*, LH-FA-20.a).
func TestPlayUpstreamTLS(t *testing.T) {
	for _, f := range []struct {
		name     string
		upstream string
		sslmode  string
		cli      []string
		env      string
		datei    string
		want     bool
	}{
		{"host:port ohne Angabe", "h:5", "", nil, "", "", false},
		{"Option ohne Wert", "h:5", "", []string{"--upstream-tls"}, "", "", true},
		{"Option true", "h:5", "", []string{"--upstream-tls=true"}, "", "", true},
		{"Option false", "h:5", "", []string{"--upstream-tls=false"}, "", "", false},
		{"Umgebungsvariable true", "h:5", "", nil, "true", "", true},
		{"Umgebungsvariable false", "h:5", "", nil, "false", "", false},
		{"Schlüssel true", "h:5", "", nil, "", "play:\n  upstream_tls: true\n", true},
		{"Schlüssel false", "h:5", "", nil, "", "play:\n  upstream_tls: false\n", false},
		{"Option false vor Umgebungsvariable true", "h:5", "", []string{"--upstream-tls=false"}, "true", "", false},
		{"Umgebungsvariable false vor Schlüssel true", "h:5", "", nil, "false", "play:\n  upstream_tls: true\n", false},
		{"Option true vor Umgebungsvariable false", "h:5", "", []string{"--upstream-tls=true"}, "false", "", true},
		{"Verbindung ohne sslmode", "v", "", nil, "", "", false},
		{"Verbindung disable", "v", "disable", nil, "", "", false},
		{"Verbindung require", "v", "require", nil, "", "", true},
		{"require, Option false", "v", "require", []string{"--upstream-tls=false"}, "", "", false},
		{"require, Umgebungsvariable false", "v", "require", nil, "false", "", false},
		{"require, Schlüssel false", "v", "require", nil, "", "play:\n  upstream_tls: false\n", false},
		{"require, Option true", "v", "require", []string{"--upstream-tls"}, "", "", true},
		{"require, Option false vor Umgebungsvariable true", "v", "require", []string{"--upstream-tls=false"}, "true", "", false},
		{"disable, Option true", "v", "disable", []string{"--upstream-tls"}, "", "", true},
		{"disable, Umgebungsvariable true", "v", "disable", nil, "true", "", true},
		{"disable, Schlüssel true", "v", "disable", nil, "", "play:\n  upstream_tls: true\n", true},
		{"ohne sslmode, Option true", "v", "", []string{"--upstream-tls=true"}, "", "", true},
	} {
		t.Run(f.name, func(t *testing.T) {
			leere(t, "play")
			t.Setenv("PGWIRE_RECORDER_UPSTREAM_TLS", f.env)
			args := append([]string{"--upstream=" + f.upstream, "--config=" + schreibe(t, tlsDatei(f.sslmode, f.datei))}, f.cli...)
			got, err := playMit(args...)
			if err != nil || got.UpstreamTLS != f.want {
				t.Fatalf("UpstreamTLS %v, Fehler %v, erwartet %v", got.UpstreamTLS, err, f.want)
			}
		})
	}
}

// ohneTLSFehler ist die Meldung zu --upstream-ca ohne TLS.
const ohneTLSFehler = "Konfiguration [PGR-E2001]: --upstream-ca verlangt TLS zum Server, mit --upstream-tls oder sslmode=require der Verbindung"

// Abdeckung: LH-FA-20/Negative, LH-FA-17/Negative —
// --upstream-ca ohne TLS ist PGR-E2001 und nennt die Option, nicht den Pfad,
// aus jeder Quelle (Option, Umgebungsvariable, Schlüssel) und für jede Weise,
// ohne TLS zu sein: ohne Angabe bei host:port, bei einer Verbindung ohne und
// mit sslmode=disable, und bei sslmode=require mit --upstream-tls=false aus
// Option, Umgebungsvariable und Schlüssel; mit TLS ist es gültig. Die Prüfung
// folgt --upstream und steht vor den Variablen der Platzhalter und vor dem Lesen
// der Datei (LH-FA-17.a *Fehler*).
func TestPlayUpstreamCAOhneTLS(t *testing.T) {
	ca := schreibe(t, string(testpki.NeueCA(t, "CA").PEM()))
	envCA, envTLS := "PGWIRE_RECORDER_UPSTREAM_CA", "PGWIRE_RECORDER_UPSTREAM_TLS"
	for _, f := range []struct {
		name     string
		upstream string
		sslmode  string
		cli      []string
		env      map[string]string
		datei    string
		fehler   bool
	}{
		{"host:port, Option", "h:5", "", []string{"--upstream-ca=" + ca}, nil, "", true},
		{"host:port, Umgebungsvariable", "h:5", "", nil, map[string]string{envCA: ca}, "", true},
		{"host:port, Schlüssel", "h:5", "", nil, nil, "play:\n  upstream_ca: " + ca + "\n", true},
		{"Verbindung ohne sslmode", "v", "", []string{"--upstream-ca=" + ca}, nil, "", true},
		{"Verbindung disable", "v", "disable", []string{"--upstream-ca=" + ca}, nil, "", true},
		{"require, Option false", "v", "require", []string{"--upstream-ca=" + ca, "--upstream-tls=false"}, nil, "", true},
		{"require, Umgebungsvariable false", "v", "require", []string{"--upstream-ca=" + ca}, map[string]string{envTLS: "false"}, "", true},
		{"require, Schlüssel false", "v", "require", []string{"--upstream-ca=" + ca}, nil, "play:\n  upstream_tls: false\n", true},
		{"host:port, Option TLS", "h:5", "", []string{"--upstream-ca=" + ca, "--upstream-tls"}, nil, "", false},
		{"host:port, Umgebungsvariable TLS", "h:5", "", []string{"--upstream-ca=" + ca}, map[string]string{envTLS: "true"}, "", false},
		{"host:port, Schlüssel TLS", "h:5", "", []string{"--upstream-ca=" + ca}, nil, "play:\n  upstream_tls: true\n", false},
		{"require", "v", "require", []string{"--upstream-ca=" + ca}, nil, "", false},
	} {
		t.Run(f.name, func(t *testing.T) {
			leere(t, "play")
			for k, v := range f.env {
				t.Setenv(k, v)
			}
			args := append([]string{"--upstream=" + f.upstream, "--config=" + schreibe(t, tlsDatei(f.sslmode, f.datei))}, f.cli...)
			_, err := playMit(args...)
			switch {
			case f.fehler && (err == nil || err.Error() != ohneTLSFehler || !istUsage(err)):
				t.Fatalf("Fehler %v, erwartet %q", err, ohneTLSFehler)
			case !f.fehler && err != nil:
				t.Fatalf("Fehler %v, erwartet keinen", err)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-17/Negative — die Prüfung von
// --upstream-ca ohne TLS steht nach den Pflichtoptionen und nach --upstream,
// vor den Variablen der Platzhalter und vor dem Lesen der Datei: Eine
// fehlende Pflichtoption und ein ungültiges --upstream gehen ihr vor, sie geht
// einer nicht gesetzten Variablen und einer ungültigen Datei vor (LH-FA-17.a
// *Fehler*).
func TestPlayUpstreamCAOhneTLSReihenfolge(t *testing.T) {
	leere(t, "play")
	variablen(t, map[string]*string{"PGR_T_H": nil})
	ungueltig := schreibe(t, "kein PEM")
	url := mitVerbindung("postgresql://u@${PGR_T_H}/db")
	if _, err := lese("play", "--upstream=h:5", "--upstream-ca="+ungueltig); !istUsage(err) || !strings.Contains(err.Error(), "Pflichtoption --input") {
		t.Errorf("fehlende Pflichtoption --input: %v", err)
	}
	if _, err := playMit("--upstream=GEHEIM", "--upstream-ca="+ungueltig); !istUsage(err) || !strings.Contains(err.Error(), "--upstream:") {
		t.Errorf("ungültiges --upstream: %v", err)
	}
	if _, err := playMit("--upstream=v", "--config="+schreibe(t, url), "--upstream-ca="+ungueltig); err == nil || err.Error() != ohneTLSFehler {
		t.Errorf("--upstream-ca ohne TLS neben nicht gesetzter Variable und ungültiger Datei: %v", err)
	}
	_, err := playMit("--upstream=v", "--config="+schreibe(t, url), "--upstream-ca="+ungueltig, "--upstream-tls")
	if !hatCode(err, model.CodeConfigVariable) {
		t.Errorf("mit TLS gehen die Variablen der Datei vor: %v", err)
	}
	_, err = playMit("--upstream=v", "--config="+schreibe(t, mitVerbindung("postgresql://u@h:${PGR_T_P}/db")), "--upstream-ca="+ungueltig, "--upstream-tls")
	if !hatCode(err, model.CodeConfigVariable) {
		t.Errorf("mit TLS und fehlendem Port: %v", err)
	}
	variablen(t, map[string]*string{"PGR_T_P": w("kein-port")})
	_, err = playMit("--upstream=v", "--config="+schreibe(t, mitVerbindung("postgresql://u@h:${PGR_T_P}/db")), "--upstream-ca="+ungueltig, "--upstream-tls")
	if !istDatei(err) || !strings.Contains(err.Error(), "Port nach dem Einsetzen") {
		t.Errorf("mit TLS und ungültigem Port geht der Port der Datei vor: %v", err)
	}
	if _, err := playMit("--upstream=h:5", "--upstream-ca="+ungueltig, "--upstream-tls"); !hatCode(err, model.CodeConfigCA) {
		t.Errorf("mit TLS und ungültiger Datei: %v, erwartet %s", err, model.CodeConfigCA)
	}
}

// pemDatei schreibt eine Datei mit inhalt und liefert ihren Pfad.
func pemDatei(t *testing.T, inhalt string) string {
	t.Helper()
	return schreibe(t, inhalt)
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary — die Datei aus --upstream-ca
// enthält einen oder mehrere PEM-Blöcke vom Typ CERTIFICATE mit lesbarem
// X.509-Zertifikat; Text vor, zwischen und nach den Blöcken bleibt unbeachtet;
// play legt die Zertifikate in dieser Reihenfolge in den Optionen ab (LH-FA-20.a
// *Start*).
func TestPlayUpstreamCADatei(t *testing.T) {
	a, b := testpki.NeueCA(t, "CA-A"), testpki.NeueCA(t, "CA-B")
	for _, f := range []struct {
		name   string
		inhalt string
		want   []*testpki.CA
	}{
		{"ein Zertifikat", string(a.PEM()), []*testpki.CA{a}},
		{"zwei Zertifikate", string(a.PEM()) + string(b.PEM()), []*testpki.CA{a, b}},
		{"zwei Zertifikate andersherum", string(b.PEM()) + string(a.PEM()), []*testpki.CA{b, a}},
		{"Text davor", "Kommentar\nzweite Zeile\n" + string(a.PEM()), []*testpki.CA{a}},
		{"Text zwischen den Blöcken", string(a.PEM()) + "dazwischen\n" + string(b.PEM()), []*testpki.CA{a, b}},
		{"Text danach", string(a.PEM()) + "Ende der Datei", []*testpki.CA{a}},
		{"Windows-Zeilenenden", strings.ReplaceAll(string(a.PEM()), "\n", "\r\n"), []*testpki.CA{a}},
	} {
		t.Run(f.name, func(t *testing.T) {
			leere(t, "play")
			got, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="+pemDatei(t, f.inhalt))
			if err != nil || len(got.UpstreamCA) != len(f.want) {
				t.Fatalf("Zertifikate %d, Fehler %v, erwartet %d", len(got.UpstreamCA), err, len(f.want))
			}
			for i, ca := range f.want {
				if !got.UpstreamCA[i].Equal(ca.Zertifikat) {
					t.Errorf("Zertifikat %d ist nicht das erwartete", i)
				}
			}
		})
	}
}

// Abdeckung: LH-FA-20/Negative, LH-RB-01/Messung — eine Datei
// aus --upstream-ca ohne PEM-Block, mit einem Block eines anderen Typs (auch
// mit dem Inhalt eines gültigen Zertifikats, auch neben einem Zertifikat, davor
// und danach), mit einem Block CERTIFICATE ohne
// lesbares X.509-Zertifikat (auch als zweiter), leer, nur Text, nur mit
// unvollständigem Block, ein Verzeichnis und eine fehlende Datei sind
// PGR-E2007 mit Exit-Code 2; die Meldung nennt die Option, weder Pfad noch
// Inhalt (LH-FA-20.a *Start*).
func TestPlayUpstreamCADateiFehler(t *testing.T) {
	ca := testpki.NeueCA(t, "CA")
	zert := string(ca.PEM())
	block := func(typ, inhalt string) string {
		return "-----BEGIN " + typ + "-----\n" + inhalt + "\n-----END " + typ + "-----\n"
	}
	// base64 von "GEHEIMINHALT" als Rumpf eines Blocks
	rumpf := "R0VIRUlNSU5IQUxU"
	for _, f := range []struct {
		name, inhalt string
	}{
		{"leer", ""},
		{"nur Text", "GEHEIMINHALT\n"},
		{"nur Zeilenumbrüche", "\n\n\n"},
		{"Block ohne Ende", "-----BEGIN CERTIFICATE-----\n" + rumpf + "\n"},
		{"Block mit ungültigem base64", block("CERTIFICATE", "!!!GEHEIM!!!")},
		{"privater Schlüssel", block("PRIVATE KEY", rumpf)},
		{"Zertifikat und privater Schlüssel danach", zert + block("PRIVATE KEY", rumpf)},
		{"privater Schlüssel und Zertifikat danach", block("PRIVATE KEY", rumpf) + zert},
		{"Zertifikatsanforderung", block("CERTIFICATE REQUEST", rumpf)},
		{"Zertifikat im Block TRUSTED CERTIFICATE", strings.ReplaceAll(zert, "CERTIFICATE", "TRUSTED CERTIFICATE")},
		{"Zertifikat, dann Zertifikat im Block X509 CERTIFICATE", zert + strings.ReplaceAll(zert, "CERTIFICATE", "X509 CERTIFICATE")},
		{"Zertifikat im Block X509 CERTIFICATE, dann Zertifikat", strings.ReplaceAll(zert, "CERTIFICATE", "X509 CERTIFICATE") + zert},
		{"X509 CERTIFICATE", block("X509 CERTIFICATE", rumpf)},
		{"Block CERTIFICATE ohne Zertifikat", block("CERTIFICATE", rumpf)},
		{"Zertifikat, dann Block CERTIFICATE ohne Zertifikat", zert + block("CERTIFICATE", rumpf)},
		{"Block CERTIFICATE ohne Zertifikat, dann Zertifikat", block("CERTIFICATE", rumpf) + zert},
		{"abgeschnittenes Zertifikat", zert[:len(zert)/2]},
	} {
		t.Run(f.name, func(t *testing.T) {
			leere(t, "play")
			pfad := pemDatei(t, f.inhalt)
			_, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="+pfad)
			if !hatCode(err, model.CodeConfigCA) || !strings.HasPrefix(err.Error(), "Konfiguration [PGR-E2007]: --upstream-ca: ") {
				t.Fatalf("Fehler %v, erwartet %s mit der Option", err, model.CodeConfigCA)
			}
			if strings.Contains(err.Error(), "GEHEIM") || strings.Contains(err.Error(), pfad) || strings.Contains(err.Error(), rumpf) {
				t.Fatalf("die Meldung nennt Pfad oder Inhalt: %v", err)
			}
		})
	}
	t.Run("Verzeichnis", func(t *testing.T) {
		leere(t, "play")
		dir := filepath.Join(t.TempDir(), "GEHEIMVERZ")
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		_, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="+dir)
		if !hatCode(err, model.CodeConfigCA) || !strings.Contains(err.Error(), "keine reguläre Datei") || strings.Contains(err.Error(), "GEHEIM") {
			t.Fatalf("Verzeichnis: %v", err)
		}
	})
	t.Run("Pfad unterhalb einer Datei", func(t *testing.T) {
		leere(t, "play")
		pfad := filepath.Join(pemDatei(t, zert), "GEHEIMUNTER.pem")
		_, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="+pfad)
		if !hatCode(err, model.CodeConfigCA) || !strings.Contains(err.Error(), "nicht lesbar") || strings.Contains(err.Error(), "GEHEIM") {
			t.Fatalf("Pfad unterhalb einer Datei: %v", err)
		}
	})
	t.Run("fehlende Datei", func(t *testing.T) {
		leere(t, "play")
		pfad := filepath.Join(t.TempDir(), "GEHEIMFEHLT.pem")
		_, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="+pfad)
		if !hatCode(err, model.CodeConfigCA) || !strings.Contains(err.Error(), "nicht vorhanden") || strings.Contains(err.Error(), "GEHEIM") {
			t.Fatalf("fehlende Datei: %v", err)
		}
	})
}

// Abdeckung: LH-FA-20/Happy, LH-FA-17/Happy, LH-FA-17/Boundary — --upstream-ca
// nimmt den Pfad aus der Option, der Umgebungsvariable PGWIRE_RECORDER_UPSTREAM_CA
// und dem Schlüssel upstream_ca des Abschnitts play:, in dieser Priorität; ein
// relativer Pfad gilt ab dem aktuellen Verzeichnis, auch in der Datei; die
// Datei wird nur gelesen, wenn die Option gesetzt ist (LH-FA-17.a, LH-FA-20.a
// *Start*).
func TestPlayUpstreamCAQuellen(t *testing.T) {
	eins, zwei, drei := testpki.NeueCA(t, "CA-1"), testpki.NeueCA(t, "CA-2"), testpki.NeueCA(t, "CA-3")
	dir := t.TempDir()
	schreibeCA := func(name string, cas ...*testpki.CA) string {
		var b []byte
		for _, c := range cas {
			b = append(b, c.PEM()...)
		}
		pfad := filepath.Join(dir, name)
		if err := os.WriteFile(pfad, b, 0o600); err != nil {
			t.Fatal(err)
		}
		return pfad
	}
	deins, dzwei, ddrei := schreibeCA("eins.pem", eins), schreibeCA("zwei.pem", zwei), schreibeCA("drei.pem", drei, drei)
	leere(t, "play")
	nur := func(args ...string) int {
		t.Helper()
		got, err := playMit(append([]string{"--upstream=h:5", "--upstream-tls"}, args...)...)
		if err != nil {
			t.Fatalf("%q: %v", args, err)
		}
		return len(got.UpstreamCA)
	}
	if n := nur(); n != 0 {
		t.Errorf("ohne Option %d Zertifikate", n)
	}
	datei := func(pfad string) string { return "--config=" + schreibe(t, "play:\n  upstream_ca: "+pfad+"\n") }
	if n := nur("--upstream-ca=" + deins); n != 1 {
		t.Errorf("Option: %d Zertifikate", n)
	}
	if n := nur(datei(ddrei)); n != 2 {
		t.Errorf("Schlüssel: %d Zertifikate", n)
	}
	t.Setenv("PGWIRE_RECORDER_UPSTREAM_CA", dzwei)
	got, err := playMit("--upstream=h:5", "--upstream-tls")
	if err != nil || len(got.UpstreamCA) != 1 || !got.UpstreamCA[0].Equal(zwei.Zertifikat) {
		t.Errorf("Umgebungsvariable: %d Zertifikate, %v", len(got.UpstreamCA), err)
	}
	if got, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="+deins); err != nil || len(got.UpstreamCA) != 1 || !got.UpstreamCA[0].Equal(eins.Zertifikat) {
		t.Errorf("Option vor Umgebungsvariable: %d Zertifikate, %v", len(got.UpstreamCA), err)
	}
	if got, err := playMit("--upstream=h:5", "--upstream-tls", datei(ddrei)); err != nil || len(got.UpstreamCA) != 1 || !got.UpstreamCA[0].Equal(zwei.Zertifikat) {
		t.Errorf("Umgebungsvariable vor Schlüssel: %d Zertifikate, %v", len(got.UpstreamCA), err)
	}
	t.Setenv("PGWIRE_RECORDER_UPSTREAM_CA", "")
	t.Chdir(dir)
	if n := nur("--upstream-ca=eins.pem"); n != 1 {
		t.Errorf("relativer Pfad der Option: %d Zertifikate", n)
	}
	if n := nur("--config=" + schreibe(t, "play:\n  upstream_ca: drei.pem\n")); n != 2 {
		t.Errorf("relativer Pfad des Schlüssels: %d Zertifikate", n)
	}
	if _, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca=fehlt.pem"); !hatCode(err, model.CodeConfigCA) {
		t.Errorf("relativer Pfad einer fehlenden Datei: %v", err)
	}
}

// Abdeckung: LH-FA-17/Negative, LH-FA-20/Negative — ein leerer Wert von
// --upstream-ca ist PGR-E2001 und nennt die Option, in der Datei PGR-E2004; eine
// leere Umgebungsvariable gilt als nicht gesetzt (LH-FA-17.a).
func TestPlayUpstreamCALeer(t *testing.T) {
	leere(t, "play")
	if _, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="); !istUsage(err) || !strings.Contains(err.Error(), "upstream-ca") {
		t.Errorf("--upstream-ca=: %v", err)
	}
	if got, err := playMit("--upstream=h:5", "--upstream-tls"); err != nil || len(got.UpstreamCA) != 0 {
		t.Errorf("leere Umgebungsvariable: %v, %d Zertifikate", err, len(got.UpstreamCA))
	}
	if _, err := playMit("--upstream=h:5", "--upstream-tls", "--config="+schreibe(t, "play:\n  upstream_ca: \"\"\n")); !istDatei(err) || !strings.Contains(err.Error(), "play.upstream_ca") {
		t.Errorf("Schlüssel leer: %v", err)
	}
}

// Abdeckung: LH-FA-17/Happy, LH-FA-20/Boundary — config show liest die Datei aus
// upstream_ca nicht: Mit einem Pfad, den es nicht gibt, zeigt es die
// Konfigurationsdatei und endet ohne Fehler (LH-FA-20.a *TLS*, LH-FA-17.a
// *Anzeige*).
func TestConfigShowLiestCANicht(t *testing.T) {
	leere(t, "play")
	pfad := schreibe(t, "play:\n  upstream_tls: true\n  upstream_ca: /GEHEIM/fehlt.pem\n")
	cmd, err := lese("config", "show", "--config="+pfad)
	if err != nil || !strings.Contains(cmd.Anzeige, "upstream_ca: /GEHEIM/fehlt.pem") {
		t.Fatalf("config show: %q, %v", cmd.Anzeige, err)
	}
}

// Abdeckung: LH-FA-20/Negative, LH-RB-01/Messung — eine Formatierung (%v, %+v,
// %#v, %s, %q, auch über einen Zeiger) der Optionen von play gibt kein Zertifikat
// aus (LH-FA-20.a *TLS*).
func TestPlayOptionenOhneZertifikat(t *testing.T) {
	leere(t, "play")
	ca := testpki.NeueCA(t, "GEHEIM-CA")
	got, err := playMit("--upstream=h:5", "--upstream-tls", "--upstream-ca="+pemDatei(t, string(ca.PEM())))
	if err != nil || len(got.UpstreamCA) != 1 {
		t.Fatalf("play: %d Zertifikate, %v", len(got.UpstreamCA), err)
	}
	for _, verb := range []string{"%v", "%+v", "%#v", "%s", "%q", "%x", "%d"} {
		for _, arg := range []any{got, &got, got.UpstreamCA, []cli.PlayOptions{got}} {
			if out := fmt.Sprintf(verb, arg); strings.Contains(out, "GEHEIM") || strings.Contains(strings.ToLower(out), "47454845494d") || strings.Contains(out, "CERTIFICATE") {
				t.Errorf("%s gibt ein Zertifikat aus: %s", verb, out)
			}
		}
	}
}

// Die Hilfe von play nennt --upstream-tls, --upstream-ca und ihre
// Umgebungsvariablen (LH-FA-01.a).
func TestPlayHilfeTLS(t *testing.T) {
	text := strings.Join(strings.Fields(hilfe(t, "play", "--help")), " ")
	for _, teil := range []string{"--upstream-tls[=true|false]", "PGWIRE_RECORDER_UPSTREAM_TLS", "--upstream-ca <datei>", "PGWIRE_RECORDER_UPSTREAM_CA", "PGR-E2001", "PGR-E2007"} {
		if !strings.Contains(text, teil) {
			t.Errorf("Hilfe von play ohne %q:\n%s", teil, text)
		}
	}
}
