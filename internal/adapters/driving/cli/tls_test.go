package cli_test

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// tlsKommandos sind die Kommandos, die TLS zum Client annehmen (LH-FA-23.a).
func tlsKommandos() []string { return []string{"record", "replay"} }

// zugangVon liefert die TLS-Angaben, die das Kommando von cmd trägt.
func zugangVon(cmd cli.Command) cli.TLSZugang {
	if cmd.Name == "record" {
		return cmd.Record.TLS
	}
	return cmd.Replay.TLS
}

// tlsLese leert die Umgebungsvariablen von kommando und liest es mit seinen
// Pflichtoptionen und args.
func tlsLese(t *testing.T, kommando string, args ...string) (cli.Command, error) {
	t.Helper()
	leere(t, kommando)
	return lese(append(basis(kommando, ""), args...)...)
}

// tlsPaar schreibt Zertifikat und Schlüssel (PKCS #8) eines frischen Blattes
// und liefert beide Pfade.
func tlsPaar(t *testing.T) (cert, key string, b testpki.Blattzertifikat) {
	t.Helper()
	b = testpki.NeueCA(t, "CA").Server(t, testpki.Blatt{})
	return schreibe(t, string(b.PEM())), schreibe(t, string(b.SchluesselPEM(t, "PRIVATE KEY"))), b
}

// Abdeckung: LH-FA-23/Negative, LH-FA-23/Boundary, LH-FA-17/Negative — --tls-cert
// allein, --tls-key allein und ein wirksames true von --allow-plaintext ohne
// --tls-cert sind PGR-E2001 und nennen die erste verletzte Option in der
// Reihenfolge --tls-cert, --tls-key, --allow-plaintext, nicht den Pfad;
// --allow-plaintext=false ist auch ohne --tls-cert gültig, beide Dateien allein
// oder mit --allow-plaintext sind gültig, bei record und replay (LH-FA-23.a
// *Start*).
func TestTLSKombination(t *testing.T) {
	cert, key, _ := tlsPaar(t)
	for _, f := range []struct {
		name   string
		args   []string
		fehler string
	}{
		{"nur Zertifikat", []string{"--tls-cert=" + cert}, "--tls-cert"},
		{"nur Schlüssel", []string{"--tls-key=" + key}, "--tls-key"},
		{"nur --allow-plaintext=true", []string{"--allow-plaintext=true"}, "--allow-plaintext"},
		{"nur --allow-plaintext ohne Wert", []string{"--allow-plaintext"}, "--allow-plaintext"},
		{"Schlüssel und --allow-plaintext=true", []string{"--tls-key=" + key, "--allow-plaintext=true"}, "--tls-key"},
		{"nur --allow-plaintext=false", []string{"--allow-plaintext=false"}, ""},
		{"beide Dateien", []string{"--tls-cert=" + cert, "--tls-key=" + key}, ""},
		{"beide Dateien und --allow-plaintext=true", []string{"--tls-cert=" + cert, "--tls-key=" + key, "--allow-plaintext=true"}, ""},
		{"beide Dateien und --allow-plaintext=false", []string{"--tls-cert=" + cert, "--tls-key=" + key, "--allow-plaintext=false"}, ""},
		{"letzte Angabe gilt", []string{"--allow-plaintext=true", "--allow-plaintext=false"}, ""},
	} {
		for _, kommando := range tlsKommandos() {
			t.Run(kommando+"/"+f.name, func(t *testing.T) {
				_, err := tlsLese(t, kommando, f.args...)
				switch {
				case f.fehler == "" && err != nil:
					t.Fatalf("Fehler %v, erwartet keinen", err)
				case f.fehler != "" && (!istUsage(err) || !strings.HasPrefix(err.Error(), "Konfiguration [PGR-E2001]: "+f.fehler) || strings.Contains(err.Error(), cert) || strings.Contains(err.Error(), key)):
					t.Fatalf("Fehler %v, erwartet %s mit %s", err, model.CodeUsage, f.fehler)
				}
			})
		}
	}
}

// Abdeckung: LH-FA-23/Negative, LH-FA-17/Negative — die Kombinationen kommen vor
// jedem Fehler von --upstream, das Laden der Dateien nach ihm: Neben einem
// ungültigen --upstream meldet record die Kombination, aber bei gültigen
// Kombinationen den Fehler von --upstream vor dem der Dateien; das Zertifikat
// wird vor dem Schlüssel geladen (LH-FA-17.a *Fehler*, LH-FA-23.a *Start*).
func TestTLSReihenfolge(t *testing.T) {
	fehlt := filepath.Join(t.TempDir(), "fehlt.pem")
	_, err := lese("record", "--listen=a:1", "--output=o", "--upstream=ungueltig", "--tls-key="+fehlt)
	if !istUsage(err) || !strings.Contains(err.Error(), "--tls-key verlangt") {
		t.Errorf("Kombination neben ungültigem --upstream: %v", err)
	}
	_, err = lese("record", "--listen=a:1", "--output=o", "--upstream=ungueltig", "--tls-cert="+fehlt, "--tls-key="+fehlt)
	if !istUsage(err) || !strings.Contains(err.Error(), "--upstream") {
		t.Errorf("ungültiges --upstream neben nicht lesbaren Dateien: %v", err)
	}
	for _, kommando := range tlsKommandos() {
		_, err = tlsLese(t, kommando, "--tls-cert="+fehlt, "--tls-key="+fehlt)
		if !hatCode(err, model.CodeConfigCA) || !strings.Contains(err.Error(), "--tls-cert: ") {
			t.Errorf("%s: zwei fehlende Dateien: %v, erwartet das Zertifikat zuerst", kommando, err)
		}
		cert, _, _ := tlsPaar(t)
		_, err = tlsLese(t, kommando, "--tls-cert="+cert, "--tls-key="+fehlt)
		if !hatCode(err, model.CodeConfigCA) || !strings.Contains(err.Error(), "--tls-key: ") {
			t.Errorf("%s: gültiges Zertifikat, fehlender Schlüssel: %v", kommando, err)
		}
	}
}

// Abdeckung: LH-FA-23/Happy, LH-FA-23/Boundary, LH-FA-17/Happy, LH-FA-17/Boundary —
// die drei Optionen nehmen ihren Wert aus der Kommandozeile, der Umgebungsvariable
// und dem Schlüssel tls_cert, tls_key, allow_plaintext des Abschnitts record: oder
// replay: in dieser Priorität, auch gemischt (Zertifikat aus der Umgebung,
// Schlüssel aus der Datei); die Optionen legen Zertifikat, Schlüssel und Zulassung
// nur beim eigenen Kommando ab; ein ungültiger Wert der Umgebungsvariable von
// --allow-plaintext ist PGR-E2001, auch neben der Kommandozeile, ein leerer Pfad
// auf der Kommandozeile PGR-E2001, in der Datei PGR-E2004, die leere Variable
// gilt als nicht gesetzt (LH-FA-17.a, LH-FA-23.a *Start*).
func TestTLSQuellen(t *testing.T) {
	cert1, key1, b1 := tlsPaar(t)
	cert2, key2, b2 := tlsPaar(t)
	cert3, key3, _ := tlsPaar(t)
	for _, kommando := range tlsKommandos() {
		t.Run(kommando, func(t *testing.T) {
			leere(t, kommando)
			datei := func(zeilen string) string { return "--config=" + schreibe(t, kommando+":\n"+zeilen) }
			pruefe := func(was string, args []string, wollen testpki.Blattzertifikat, klartext bool) {
				t.Helper()
				cmd, err := lese(append(basis(kommando, ""), args...)...)
				z := zugangVon(cmd)
				if err != nil || !z.Gesetzt() || z.Zertifikat != string(wollen.PEM()) || z.Schluessel != cli.Schluessel(wollen.SchluesselPEM(t, "PRIVATE KEY")) || z.KlartextErlaubt != klartext {
					t.Errorf("%s: Fehler %v, Zugang gesetzt %v, Zertifikat gleich %v, Klartext %v, erwartet Klartext %v", was, err, z.Gesetzt(), z.Zertifikat == string(wollen.PEM()), z.KlartextErlaubt, klartext)
				}
			}
			pruefe("Option", []string{"--tls-cert=" + cert1, "--tls-key=" + key1, "--allow-plaintext"}, b1, true)
			pruefe("Option, Zulassung false", []string{"--tls-cert=" + cert1, "--tls-key=" + key1, "--allow-plaintext=false"}, b1, false)
			pruefe("Option ohne Zulassung", []string{"--tls-cert=" + cert1, "--tls-key=" + key1}, b1, false)
			pruefe("Schlüssel der Datei", []string{datei("  tls_cert: " + cert2 + "\n  tls_key: " + key2 + "\n  allow_plaintext: true\n")}, b2, true)
			t.Setenv("PGWIRE_RECORDER_TLS_CERT", cert2)
			t.Setenv("PGWIRE_RECORDER_TLS_KEY", key2)
			t.Setenv("PGWIRE_RECORDER_ALLOW_PLAINTEXT", "true")
			cmd, err := lese(basis(kommando, "")...)
			if z := zugangVon(cmd); err != nil || z.Zertifikat != string(b2.PEM()) || !z.KlartextErlaubt {
				t.Errorf("Umgebungsvariablen: %v, Klartext %v", err, z.KlartextErlaubt)
			}
			cmd, err = lese(append(basis(kommando, ""), "--tls-cert="+cert1, "--tls-key="+key1, "--allow-plaintext=false")...)
			if z := zugangVon(cmd); err != nil || z.Zertifikat != string(b1.PEM()) || z.KlartextErlaubt {
				t.Errorf("Option vor Umgebungsvariable: %v, Klartext %v", err, z.KlartextErlaubt)
			}
			cmd, err = lese(append(basis(kommando, ""), datei("  tls_cert: "+cert3+"\n  tls_key: "+key3+"\n  allow_plaintext: false\n"))...)
			if z := zugangVon(cmd); err != nil || z.Zertifikat != string(b2.PEM()) || !z.KlartextErlaubt {
				t.Errorf("Umgebungsvariable vor Datei: %v, Klartext %v", err, z.KlartextErlaubt)
			}
			t.Setenv("PGWIRE_RECORDER_TLS_KEY", "")
			pruefe("Zertifikat aus der Umgebung, Schlüssel aus der Datei", []string{datei("  tls_key: " + key2 + "\n")}, b2, true)
			t.Setenv("PGWIRE_RECORDER_ALLOW_PLAINTEXT", "ja")
			if _, err := lese(append(basis(kommando, ""), "--tls-cert="+cert1, "--tls-key="+key1, "--allow-plaintext=false")...); !istUsage(err) || !strings.Contains(err.Error(), "PGWIRE_RECORDER_ALLOW_PLAINTEXT") {
				t.Errorf("ungültige Variable neben der Option: %v", err)
			}
			t.Setenv("PGWIRE_RECORDER_ALLOW_PLAINTEXT", "")
			t.Setenv("PGWIRE_RECORDER_TLS_CERT", "")
			if cmd, err := lese(basis(kommando, "")...); err != nil || zugangVon(cmd).Gesetzt() {
				t.Errorf("leere Variablen: %v, gesetzt %v", err, zugangVon(cmd).Gesetzt())
			}
			for _, o := range []string{"--tls-cert=", "--tls-key="} {
				if _, err := lese(append(basis(kommando, ""), o)...); !istUsage(err) || !strings.Contains(err.Error(), o[2:len(o)-1]) {
					t.Errorf("%s: %v", o, err)
				}
			}
			for _, k := range []string{"tls_cert", "tls_key"} {
				if _, err := lese(append(basis(kommando, ""), datei("  "+k+": \"\"\n"))...); !istDatei(err) || !strings.Contains(err.Error(), kommando+"."+k) {
					t.Errorf("Schlüssel %s leer: %v", k, err)
				}
			}
		})
	}
	t.Run("play", func(t *testing.T) {
		leere(t, "play")
		for _, k := range []string{"tls_cert", "tls_key", "allow_plaintext"} {
			if _, err := playMit("--upstream=h:5", "--config="+schreibe(t, "play:\n  "+k+": x\n")); !istDatei(err) || !strings.Contains(err.Error(), "play."+k) {
				t.Errorf("Schlüssel %s in play: %v", k, err)
			}
		}
	})
}

// zweiSchluessel liefert den Schlüssel eines zweiten Blattes der Art art als PEM
// vom Typ typ.
func zweiSchluessel(t *testing.T, art testpki.Art, typ string) []byte {
	t.Helper()
	return testpki.NeueCA(t, "CA").Server(t, testpki.Blatt{Art: art}).SchluesselPEM(t, typ)
}

// Abdeckung: LH-FA-23/Happy, LH-FA-23/Boundary — Zertifikat und Schlüssel werden
// angenommen: der Schlüssel als PRIVATE KEY (ECDSA, RSA, Ed25519), RSA PRIVATE KEY
// und EC PRIVATE KEY; die Zertifikatsdatei mit Kette, auch mit einem abgelaufenen
// zweiten Zertifikat; ein Zertifikat ohne Namen, mit anderem Namen, mit
// Schlüsselnutzung und Verwendung eines Clients und ohne die Kette seines
// Ausstellers; ein Link auf eine reguläre Datei; die Optionen tragen die Dateien
// unverändert (LH-FA-23.a *Dateien*).
func TestTLSDateienGueltig(t *testing.T) {
	ca := testpki.NeueCA(t, "CA")
	abgelaufen := ca.Server(t, testpki.Blatt{Von: time.Now().Add(-3 * time.Hour), Bis: time.Now().Add(-2 * time.Hour)})
	ecdsa, rsa, ed := ca.Server(t, testpki.Blatt{}), ca.Server(t, testpki.Blatt{Art: testpki.RSA}), ca.Server(t, testpki.Blatt{Art: testpki.Ed25519})
	for _, f := range []struct {
		name string
		b    testpki.Blattzertifikat
		typ  string
		cert string
	}{
		{"ECDSA als PRIVATE KEY", ecdsa, "PRIVATE KEY", string(ecdsa.PEM())},
		{"RSA als PRIVATE KEY", rsa, "PRIVATE KEY", string(rsa.PEM())},
		{"Ed25519 als PRIVATE KEY", ed, "PRIVATE KEY", string(ed.PEM())},
		{"RSA als RSA PRIVATE KEY", rsa, "RSA PRIVATE KEY", string(rsa.PEM())},
		{"ECDSA als EC PRIVATE KEY", ecdsa, "EC PRIVATE KEY", string(ecdsa.PEM())},
		{"Kette mit abgelaufenem zweiten Zertifikat", ecdsa, "PRIVATE KEY", string(ecdsa.PEM()) + string(abgelaufen.PEM())},
		{"ohne Namen", ca.Server(t, testpki.Blatt{}), "PRIVATE KEY", ""},
		{"anderer Name", ca.Server(t, testpki.Blatt{DNS: []string{"anderswo.example"}}), "PRIVATE KEY", ""},
		{"Verwendung eines Clients", ca.Server(t, testpki.Blatt{Verwendung: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, Schluesselnutzung: x509.KeyUsageCertSign}), "PRIVATE KEY", ""},
	} {
		for _, kommando := range tlsKommandos() {
			t.Run(kommando+"/"+f.name, func(t *testing.T) {
				cert := f.cert
				if cert == "" {
					cert = string(f.b.PEM())
				}
				key := f.b.SchluesselPEM(t, f.typ)
				cmd, err := tlsLese(t, kommando, "--tls-cert="+schreibe(t, cert), "--tls-key="+schreibe(t, string(key)))
				if z := zugangVon(cmd); err != nil || z.Zertifikat != cert || z.Schluessel != cli.Schluessel(key) {
					t.Fatalf("Fehler %v, Dateien unverändert abgelegt: %v", err, z.Zertifikat == cert && z.Schluessel == cli.Schluessel(key))
				}
			})
		}
	}
}

// pemBlock schreibt einen PEM-Block des Typs typ mit den Bytes roh und den
// Kopfzeilen kopf.
func pemBlock(typ string, kopf map[string]string, roh []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: typ, Headers: kopf, Bytes: roh}))
}

// x25519 liefert einen X25519-Schlüssel als PKCS #8, einen Schlüssel, den weder RSA
// noch ECDSA noch Ed25519 meint.
func x25519(t *testing.T) []byte {
	t.Helper()
	k, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	roh, err := x509.MarshalPKCS8PrivateKey(k)
	if err != nil {
		t.Fatal(err)
	}
	return roh
}

// Abdeckung: LH-FA-23/Negative, LH-RB-01/Messung — eine Zertifikatsdatei ohne PEM-Block,
// mit einem Schlüsselblock (auch neben dem Zertifikat, davor und danach), mit einem
// Block CERTIFICATE ohne Zertifikat (auch als zweiter) und mit Zertifikat und
// Schlüssel ist PGR-E2007 mit der Option --tls-cert; eine Schlüsseldatei ohne
// PEM-Block, mit verschlüsseltem Schlüssel (ENCRYPTED PRIVATE KEY, Proc-Type
// 4,ENCRYPTED), mit einem anderen Blocktyp (CERTIFICATE, PUBLIC KEY), mit zwei
// Blöcken (auch Schlüssel und Zertifikat, Zertifikat und Schlüssel), mit
// unlesbarem Inhalt und mit einem Schlüssel, der nicht zum Zertifikat passt
// (gleiche und andere Art), ist PGR-E2007 mit der Option; die Ursache prüft der
// Test nur bei Dateien mit einem einzigen Mangel, bei mehreren ist der Text nicht
// zugesagt; die Meldung nennt weder Pfad noch Inhalt (LH-FA-23.a *Dateien*).
func TestTLSDateienFehler(t *testing.T) {
	ca := testpki.NeueCA(t, "CA")
	b := ca.Server(t, testpki.Blatt{})
	rsa := ca.Server(t, testpki.Blatt{Art: testpki.RSA})
	cert, key := string(b.PEM()), string(b.SchluesselPEM(t, "PRIVATE KEY"))
	rumpf := []byte("GEHEIMINHALT")
	for _, f := range []struct {
		option, name, inhalt, ursache string
	}{
		{"--tls-cert", "leer", "", "kein PEM-Block"},
		{"--tls-cert", "Schlüsselblock", key, "nicht den Typ CERTIFICATE"},
		{"--tls-cert", "Zertifikat, dann Schlüssel", cert + key, ""},
		{"--tls-cert", "Schlüssel, dann Zertifikat", key + cert, ""},
		{"--tls-cert", "Block ohne Zertifikat", pemBlock("CERTIFICATE", nil, rumpf), "kein lesbares X.509-Zertifikat"},
		{"--tls-cert", "Zertifikat, dann Block ohne Zertifikat", cert + pemBlock("CERTIFICATE", nil, rumpf), ""},
		{"--tls-key", "leer", "", "kein PEM-Block"},
		{"--tls-key", "ENCRYPTED PRIVATE KEY", pemBlock("ENCRYPTED PRIVATE KEY", nil, rumpf), "verschlüsselt"},
		{"--tls-key", "Proc-Type 4,ENCRYPTED", pemBlock("RSA PRIVATE KEY", map[string]string{"Proc-Type": "4,ENCRYPTED", "DEK-Info": "AES-128-CBC,00"}, rumpf), "verschlüsselt"},
		{"--tls-key", "Zertifikat", cert, "keinen der Typen"},
		{"--tls-key", "PUBLIC KEY", pemBlock("PUBLIC KEY", nil, rumpf), "keinen der Typen"},
		{"--tls-key", "zwei Schlüssel", key + key, "mehr als ein PEM-Block"},
		{"--tls-key", "Schlüssel, dann Zertifikat", key + cert, "mehr als ein PEM-Block"},
		{"--tls-key", "Zertifikat, dann Schlüssel", cert + key, ""},
		{"--tls-key", "X25519-Schlüssel", pemBlock("PRIVATE KEY", nil, x25519(t)), "weder RSA noch ECDSA noch Ed25519"},
		{"--tls-key", "unlesbarer Schlüssel", pemBlock("PRIVATE KEY", nil, rumpf), "keinen lesbaren Schlüssel"},
		{"--tls-key", "unlesbarer RSA-Schlüssel", pemBlock("RSA PRIVATE KEY", nil, rumpf), "keinen lesbaren Schlüssel"},
		{"--tls-key", "unlesbarer EC-Schlüssel", pemBlock("EC PRIVATE KEY", nil, rumpf), "keinen lesbaren Schlüssel"},
		{"--tls-key", "anderer Schlüssel derselben Art", string(zweiSchluessel(t, testpki.ECDSA, "PRIVATE KEY")), "passt nicht zum Zertifikat"},
		{"--tls-key", "Schlüssel einer anderen Art", string(rsa.SchluesselPEM(t, "PRIVATE KEY")), "passt nicht zum Zertifikat"},
	} {
		for _, kommando := range tlsKommandos() {
			t.Run(kommando+"/"+f.option+"/"+f.name, func(t *testing.T) {
				certInhalt, keyInhalt := cert, key
				if f.option == "--tls-cert" {
					certInhalt = f.inhalt
				} else {
					keyInhalt = f.inhalt
				}
				cp, kp := schreibe(t, certInhalt), schreibe(t, keyInhalt)
				_, err := tlsLese(t, kommando, "--tls-cert="+cp, "--tls-key="+kp)
				if !hatCode(err, model.CodeConfigCA) || !strings.HasPrefix(err.Error(), "Konfiguration [PGR-E2007]: "+f.option+": ") || (f.ursache != "" && !strings.Contains(err.Error(), f.ursache)) {
					t.Fatalf("Fehler %v, erwartet %s mit %s und %q", err, model.CodeConfigCA, f.option, f.ursache)
				}
				for _, geheim := range []string{cp, kp, "GEHEIM", "BEGIN"} {
					if strings.Contains(err.Error(), geheim) {
						t.Fatalf("die Meldung nennt %q: %v", geheim, err)
					}
				}
			})
		}
	}
	for _, kommando := range tlsKommandos() {
		t.Run(kommando+"/eine Datei für beide", func(t *testing.T) {
			beide := schreibe(t, cert+key)
			_, err := tlsLese(t, kommando, "--tls-cert="+beide, "--tls-key="+beide)
			if !hatCode(err, model.CodeConfigCA) || !strings.Contains(err.Error(), "--tls-cert: ") {
				t.Fatalf("Fehler %v", err)
			}
			beide = schreibe(t, key+cert)
			_, err = tlsLese(t, kommando, "--tls-cert="+beide, "--tls-key="+beide)
			if !hatCode(err, model.CodeConfigCA) || !strings.Contains(err.Error(), "--tls-cert: ") {
				t.Fatalf("Schlüssel zuerst: %v", err)
			}
		})
	}
}

// Abdeckung: LH-FA-23/Boundary, LH-FA-23/Negative — das Zertifikat des Servers ist
// gültig, wenn die Systemuhr beim Start zwischen Beginn und Ende seiner
// Gültigkeit liegt, beide Grenzen eingeschlossen: Eine Nanosekunde vor dem Beginn
// ist es noch nicht gültig, eine nach dem Ende abgelaufen (PGR-E2007, --tls-cert);
// bei record und replay; ohne Uhr gilt die Systemuhr (LH-FA-23.a *Dateien*).
func TestTLSGueltigkeit(t *testing.T) {
	von, bis := time.Unix(1_800_000_000, 0), time.Unix(1_800_003_600, 0)
	b := testpki.NeueCA(t, "CA").Server(t, testpki.Blatt{Von: von, Bis: bis})
	cert, key := schreibe(t, string(b.PEM())), schreibe(t, string(b.SchluesselPEM(t, "PRIVATE KEY")))
	for _, f := range []struct {
		name    string
		jetzt   time.Time
		ursache string
	}{
		{"Beginn", von, ""},
		{"Ende", bis, ""},
		{"mittendrin", von.Add(time.Minute), ""},
		{"eine Nanosekunde vor dem Beginn", von.Add(-time.Nanosecond), "noch nicht gültig"},
		{"eine Nanosekunde nach dem Ende", bis.Add(time.Nanosecond), "abgelaufen"},
	} {
		for _, kommando := range tlsKommandos() {
			t.Run(kommando+"/"+f.name, func(t *testing.T) {
				leere(t, kommando)
				args := append(basis(kommando, ""), "--tls-cert="+cert, "--tls-key="+key)
				_, err := cli.ParseMitUhr(args, io.Discard, func() time.Time { return f.jetzt })
				switch {
				case f.ursache == "" && err != nil:
					t.Fatalf("Fehler %v, erwartet keinen", err)
				case f.ursache != "" && (!hatCode(err, model.CodeConfigCA) || !strings.HasPrefix(err.Error(), "Konfiguration [PGR-E2007]: --tls-cert: ") || !strings.Contains(err.Error(), f.ursache)):
					t.Fatalf("Fehler %v, erwartet %s mit %q", err, model.CodeConfigCA, f.ursache)
				}
			})
		}
	}
	t.Run("Systemuhr", func(t *testing.T) {
		ca := testpki.NeueCA(t, "CA")
		gueltig := ca.Server(t, testpki.Blatt{})
		alt := ca.Server(t, testpki.Blatt{Von: time.Now().Add(-2 * time.Hour), Bis: time.Now().Add(-time.Hour)})
		neu := ca.Server(t, testpki.Blatt{Von: time.Now().Add(time.Hour), Bis: time.Now().Add(2 * time.Hour)})
		for _, f := range []struct {
			b       testpki.Blattzertifikat
			gueltig bool
		}{{gueltig, true}, {alt, false}, {neu, false}} {
			_, err := tlsLese(t, "record", "--tls-cert="+schreibe(t, string(f.b.PEM())), "--tls-key="+schreibe(t, string(f.b.SchluesselPEM(t, "PRIVATE KEY"))))
			if (err == nil) != f.gueltig {
				t.Errorf("gültig %v, Fehler %v", f.gueltig, err)
			}
		}
	})
}

// Abdeckung: LH-FA-23/Boundary — ein relativer Pfad von --tls-cert und --tls-key
// gilt ab dem aktuellen Verzeichnis, auch wenn er in der Konfigurationsdatei
// steht und diese in einem anderen Verzeichnis liegt (LH-FA-23.a).
func TestTLSRelativ(t *testing.T) {
	dir := t.TempDir()
	b := testpki.NeueCA(t, "CA").Server(t, testpki.Blatt{})
	for name, inhalt := range map[string][]byte{"c.pem": b.PEM(), "k.pem": b.SchluesselPEM(t, "PRIVATE KEY")} {
		if err := os.WriteFile(filepath.Join(dir, name), inhalt, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(dir)
	datei := schreibe(t, "record:\n  tls_cert: c.pem\n  tls_key: k.pem\n")
	for _, args := range [][]string{{"--tls-cert=c.pem", "--tls-key=k.pem"}, {"--config=" + datei}} {
		if cmd, err := tlsLese(t, "record", args...); err != nil || !cmd.Record.TLS.Gesetzt() {
			t.Errorf("%q: %v", args, err)
		}
	}
}

// Abdeckung: LH-FA-23/Negative, LH-RB-01/Messung — eine Formatierung (%v, %+v,
// %#v, %s, %q, %x, %d, auch über einen Zeiger) der Optionen von record und replay
// gibt den Schlüssel nicht aus, weder als Text noch als Zahlen (LH-FA-23.a *Dateien*).
func TestTLSOptionenOhneSchluessel(t *testing.T) {
	cert, key, b := tlsPaar(t)
	roh := b.SchluesselPEM(t, "PRIVATE KEY")
	zahlen := strings.TrimPrefix(fmt.Sprintf("%d", roh[:12]), "[")
	for _, kommando := range tlsKommandos() {
		cmd, err := tlsLese(t, kommando, "--tls-cert="+cert, "--tls-key="+key)
		if err != nil || zugangVon(cmd).Schluessel == "" {
			t.Fatalf("%s: %v", kommando, err)
		}
		for _, verb := range []string{"%v", "%+v", "%#v", "%s", "%q", "%x", "%d"} {
			for _, arg := range []any{cmd, &cmd, zugangVon(cmd), zugangVon(cmd).Schluessel, cmd.Record, cmd.Replay} {
				out := fmt.Sprintf(verb, arg)
				if strings.Contains(out, "PRIVATE KEY") || strings.Contains(out, string(roh[40:70])) || strings.Contains(out, zahlen) || strings.Contains(out, fmt.Sprintf("%x", roh[:12])) {
					t.Errorf("%s %s gibt den Schlüssel aus: %s", kommando, verb, out)
				}
			}
		}
	}
}

// Abdeckung: LH-FA-23/Happy — config show liest die Dateien aus tls_cert und
// tls_key nicht: Mit Pfaden, die es nicht gibt, zeigt es die Konfigurationsdatei
// und endet ohne Fehler (LH-FA-23.a *Start*, LH-FA-17.a *Anzeige*).
func TestConfigShowLiestTLSNicht(t *testing.T) {
	leere(t, "record")
	pfad := schreibe(t, "record:\n  tls_cert: /GEHEIM/fehlt.pem\n  tls_key: /GEHEIM/fehlt.key\n")
	cmd, err := lese("config", "show", "--config="+pfad)
	if err != nil || !strings.Contains(cmd.Anzeige, "tls_cert: /GEHEIM/fehlt.pem") || !strings.Contains(cmd.Anzeige, "tls_key: /GEHEIM/fehlt.key") {
		t.Fatalf("config show: %q, %v", cmd.Anzeige, err)
	}
}

// Die Hilfe von record und replay nennt --tls-cert, --tls-key, --allow-plaintext
// und ihre Umgebungsvariablen (LH-FA-01.a); die Aussagen des Texts belegen
// TestTLSKombination, TestTLSDateienFehler und TestTLSGueltigkeit.
func TestTLSHilfe(t *testing.T) {
	for _, kommando := range tlsKommandos() {
		text := strings.Join(strings.Fields(hilfe(t, kommando, "--help")), " ")
		for _, teil := range []string{"--tls-cert <datei>", "PGWIRE_RECORDER_TLS_CERT", "--tls-key <datei>", "PGWIRE_RECORDER_TLS_KEY", "--allow-plaintext[=true|false]", "PGWIRE_RECORDER_ALLOW_PLAINTEXT", "PGR-E2001", "PGR-E2007", "PGR-E6003"} {
			if !strings.Contains(text, teil) {
				t.Errorf("Hilfe von %s ohne %q:\n%s", kommando, teil, text)
			}
		}
	}
}

// Abdeckung: LH-FA-23/Boundary — Text außerhalb der PEM-Blöcke bleibt in beiden
// Dateien unbeachtet: vor dem ersten, zwischen zwei Blöcken und nach dem letzten,
// bei der Schlüsseldatei auch nach dem Schlüsselblock (LH-FA-23.a *Dateien*).
func TestTLSTextAusserhalbDerBloecke(t *testing.T) {
	ca := testpki.NeueCA(t, "CA")
	b := ca.Server(t, testpki.Blatt{})
	cert, key := string(b.PEM()), string(b.SchluesselPEM(t, "PRIVATE KEY"))
	for _, f := range []struct{ name, cert, key string }{
		{"vor", "Kommentar\n" + cert, "Kommentar\n" + key},
		{"zwischen", cert + "dazwischen\n" + string(ca.PEM()), key},
		{"nach", cert + "Ende\n", key + "Ende\n"},
	} {
		for _, kommando := range tlsKommandos() {
			t.Run(kommando+"/"+f.name, func(t *testing.T) {
				cmd, err := tlsLese(t, kommando, "--tls-cert="+schreibe(t, f.cert), "--tls-key="+schreibe(t, f.key))
				if z := zugangVon(cmd); err != nil || z.Zertifikat != f.cert || z.Schluessel != cli.Schluessel(f.key) {
					t.Fatalf("Fehler %v", err)
				}
			})
		}
	}
}
