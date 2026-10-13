//go:build integration

package integration_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// tlsPaar schreibt ein selbstsigniertes Zertifikat und seinen Schlüssel in ein
// Verzeichnis des Tests und liefert beide Pfade.
func tlsPaar(t *testing.T) (cert, key string) {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	vorlage := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "recorder"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, err := x509.CreateCertificate(rand.Reader, vorlage, vorlage, &k.PublicKey, k)
	if err != nil {
		t.Fatal(err)
	}
	roh, err := x509.MarshalPKCS8PrivateKey(k)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	cert, key = filepath.Join(dir, "server.pem"), filepath.Join(dir, "server.key")
	for pfad, block := range map[string]*pem.Block{cert: {Type: "CERTIFICATE", Bytes: der}, key: {Type: "PRIVATE KEY", Bytes: roh}} {
		if err := os.WriteFile(pfad, pem.EncodeToMemory(block), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return cert, key
}

// Abdeckung: LH-FA-23/Happy — ein Client verbindet sich mit sslmode=require über
// TLS zu record, das Ergebnis ist das der unverschlüsselten Verbindung, und die
// Aufzeichnung ist byte-gleich der einer unverschlüsselten Verbindung mit denselben
// Anfragen; mit --allow-plaintext nimmt record beide Verbindungsarten an.
// Abnahmeszenario 15.
func TestE2ERecordTLS(t *testing.T) {
	cert, key := tlsPaar(t)
	queries := []string{"SELECT 1 AS eins;", "SELECT 'a' AS text, NULL AS leer;"}
	klartext := filepath.Join(t.TempDir(), "klartext.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), klartext)
	gKlartext := beobachtenMit(t, rec.listen, "disable", true, queries...)
	rec.stop(t, 0)

	verschluesselt := filepath.Join(t.TempDir(), "tls.yaml")
	rec = startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", verschluesselt, "--tls-cert", cert, "--tls-key", key, "--allow-plaintext")
	gTLS := beobachtenMit(t, rec.listen, "require", true, queries...)
	rec.stop(t, 0)
	if gTLS != gKlartext || lies(t, verschluesselt) != lies(t, klartext) {
		t.Fatalf("Ergebnis über TLS %q gegen %q; Aufzeichnungen:\n%s\n--- gegen:\n%s", gTLS, gKlartext, lies(t, verschluesselt), lies(t, klartext))
	}

	gemischt := filepath.Join(t.TempDir(), "gemischt.yaml")
	rec = startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", gemischt, "--tls-cert", cert, "--tls-key", key, "--allow-plaintext")
	beobachtenMit(t, rec.listen, "disable", true, queries...)
	beobachtenMit(t, rec.listen, "require", true, queries...)
	rec.stop(t, 0)
	if text := lies(t, gemischt); strings.Count(text, "- id:") != 2 {
		t.Fatalf("erwartet zwei Sessions:\n%s", text)
	}
}

// Abdeckung: LH-FA-23/Happy — replay mit TLS liefert die Aufzeichnung über eine
// verschlüsselte Verbindung wie über eine unverschlüsselte aus: Beide Läufe geben
// dem Client dieselbe Sicht wie beim Aufzeichnen (Abnahmeszenario 15).
func TestE2EReplayTLS(t *testing.T) {
	cert, key := tlsPaar(t)
	queries := []string{"SELECT 1 AS eins;", "SELECT 'a' AS text, NULL AS leer;"}
	input, aufgezeichnet := aufnehmen(t, queries...)
	rep := startProzess(t, "replay", "--input", input, "--tls-cert", cert, "--tls-key", key)
	if got := beobachtenMit(t, rep.listen, "require", true, queries...); got != aufgezeichnet {
		t.Fatalf("über TLS:\n%s\n--- aufgezeichnet:\n%s", got, aufgezeichnet)
	}
	rep.stop(t, 0)
	rep = startProzess(t, "replay", "--input", input, "--tls-cert", cert, "--tls-key", key, "--allow-plaintext")
	if got := beobachtenMit(t, rep.listen, "disable", true, queries...); got != aufgezeichnet {
		t.Fatalf("unverschlüsselt:\n%s\n--- aufgezeichnet:\n%s", got, aufgezeichnet)
	}
	rep.stop(t, 0)
}

// Abdeckung: LH-FA-23/Negative — ist TLS konfiguriert und --allow-plaintext nicht
// gesetzt, weist record und replay einen unverschlüsselten Client mit PGR-E6003 ab
// und endet beim Beenden mit Exit-Code 6.
func TestE2EKlartextAbgewiesen(t *testing.T) {
	cert, key := tlsPaar(t)
	input, _ := aufnehmen(t, "SELECT 1;")
	for name, rec := range map[string]*recorder{
		"record": startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", filepath.Join(t.TempDir(), "r.yaml"), "--tls-cert", cert, "--tls-key", key),
		"replay": startProzess(t, "replay", "--input", input, "--tls-cert", cert, "--tls-key", key),
	} {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		_, err := pgconn.Connect(ctx, dsn(rec.listen))
		cancel()
		if err == nil || !strings.Contains(err.Error(), "PGR-E6003") {
			t.Fatalf("%s: erwartet Fehler mit PGR-E6003, erhalten %v", name, err)
		}
		rec.stop(t, 6)
	}
}

// Abdeckung: LH-FA-23/Negative — eine ungültige Zertifikatsdatei beendet den Start
// mit PGR-E2007 und Exit-Code 2, vor der Ausgabedatei: eine vorhandene Ausgabedatei
// bleibt unberührt und ist nicht PGR-E2002; die Ausgabe nennt weder Pfad noch Inhalt.
func TestE2ETLSStartfehler(t *testing.T) {
	_, key := tlsPaar(t)
	dir := t.TempDir()
	output := filepath.Join(dir, "vorhanden.yaml")
	if err := os.WriteFile(output, []byte("vorhanden"), 0o600); err != nil {
		t.Fatal(err)
	}
	cert := filepath.Join(dir, "GEHEIMPFAD.pem")
	if err := os.WriteFile(cert, []byte("kein PEM"), 0o600); err != nil {
		t.Fatal(err)
	}
	out, err := fuehreAus(t, "record", "--listen", freieAdresse(t), "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--tls-cert", cert, "--tls-key", key)
	if exitCodeOf(err) != 2 || !strings.Contains(out, "PGR-E2007") || strings.Contains(out, "PGR-E2002") || strings.Contains(out, "GEHEIMPFAD") {
		t.Fatalf("Exit-Code %d, Ausgabe:\n%s", exitCodeOf(err), out)
	}
	if lies(t, output) != "vorhanden" {
		t.Fatal("die vorhandene Ausgabedatei wurde berührt")
	}
}

// Abdeckung: LH-FA-23/Boundary — eine offene verschlüsselte, ruhende
// Client-Verbindung hält das Beenden nach SIGTERM nicht auf; ihre abgeschlossenen
// Interaktionen stehen in der Aufzeichnung.
func TestE2ERecordTLSBeendenMitOffenerVerbindung(t *testing.T) {
	cert, key := tlsPaar(t)
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startProzess(t, "record", "--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output, "--tls-cert", cert, "--tls-key", key)
	beobachtenMit(t, rec.listen, "require", false, "SELECT 1;")
	rec.stop(t, 0)
	if text := lies(t, output); !strings.Contains(text, "sql: SELECT 1;") {
		t.Fatalf("Interaktion der offenen Verbindung fehlt:\n%s", text)
	}
}
