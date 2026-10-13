package cli

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// TLSZugang sind die Dateien, mit denen record und replay TLS von Clients
// annehmen (LH-FA-23.a): die Zertifikatsdatei als PEM-Text, die Schlüsseldatei
// und die Zulassung unverschlüsselter Clients. Ohne --tls-cert sind Zertifikat
// und Schlüssel leer.
type TLSZugang struct {
	Zertifikat      string
	Schluessel      Schluessel
	KlartextErlaubt bool
}

// Gesetzt sagt, ob TLS konfiguriert ist.
func (z TLSZugang) Gesetzt() bool { return len(z.Zertifikat) > 0 }

// Schluessel ist die PEM-Datei eines privaten Schlüssels. Jede Formatierung
// gibt festen Text aus, nie den Inhalt: Ein Aufruf mit %v, %+v, %#v, %s oder %q
// auf den Optionen von record oder replay verrät ihn nicht (LH-FA-23.a
// *Dateien*).
type Schluessel string

// Format gibt für jedes Verb den festen Text aus.
func (Schluessel) Format(f fmt.State, _ rune) {
	_, _ = io.WriteString(f, "***")
}

// tlsKombination prüft nach der Zusammenführung die Kombinationen von
// --tls-cert, --tls-key und --allow-plaintext bei record und replay: jede der
// beiden Dateien allein und ein wirksames true von --allow-plaintext ohne
// --tls-cert sind PGR-E2001; bei mehreren Verstößen nennt die Meldung den ersten
// in der Reihenfolge --tls-cert, --tls-key, --allow-plaintext (LH-FA-23.a *Start*).
func tlsKombination(kommando string, q quellen) error {
	if kommando != "record" && kommando != "replay" {
		return nil
	}
	cert, key := q.wert["tls-cert"] != "", q.wert["tls-key"] != ""
	switch {
	case cert && !key:
		return model.Errorf(model.CodeUsage, nil, "--tls-cert verlangt --tls-key")
	case key && !cert:
		return model.Errorf(model.CodeUsage, nil, "--tls-key verlangt --tls-cert")
	case q.wert["allow-plaintext"] == "true" && !cert:
		return model.Errorf(model.CodeUsage, nil, "--allow-plaintext=true verlangt --tls-cert und --tls-key")
	}
	return nil
}

// zugang liefert die TLS-Angaben des Kommandos von c.
func zugang(c *Command) *TLSZugang {
	if c.Name == "record" {
		return &c.Record.TLS
	}
	return &c.Replay.TLS
}

// tlsZertifikat liest die Zertifikatsdatei aus --tls-cert und legt ihre Bytes in
// den Optionen ab. Sie ist die erste der beiden letzten Prüfungen des Starts
// (LH-FA-23.a *Start*): jeder Block hat den Typ CERTIFICATE und ein lesbares
// X.509-Zertifikat, und das erste, das des Servers, gilt, wenn die Systemuhr
// zwischen Beginn und Ende seiner Gültigkeit liegt, beide Grenzen eingeschlossen.
// Namen, Verwendung und Aussteller prüft es nicht. Ohne die Option liest es nichts.
func tlsZertifikat(c *Command, _ gelesen, _ *datei, q quellen) error {
	pfad := q.wert["tls-cert"]
	if pfad == "" {
		return nil
	}
	roh, err := liesDatei("--tls-cert", pfad)
	if err != nil {
		return err
	}
	kette, err := zertifikate("--tls-cert", roh)
	if err != nil {
		return err
	}
	jetzt := q.jetzt()
	switch {
	case jetzt.Before(kette[0].NotBefore):
		return dateiUngueltig("--tls-cert", "das Zertifikat ist noch nicht gültig")
	case jetzt.After(kette[0].NotAfter):
		return dateiUngueltig("--tls-cert", "das Zertifikat ist abgelaufen")
	}
	zugang(c).Zertifikat = string(roh)
	return nil
}

// tlsSchluessel liest die Schlüsseldatei aus --tls-key und legt ihre Bytes in
// den Optionen ab; es läuft nach tlsZertifikat. Die Datei enthält genau einen
// PEM-Block vom Typ PRIVATE KEY, RSA PRIVATE KEY oder EC PRIVATE KEY (RSA, ECDSA
// oder Ed25519), unverschlüsselt, und sein öffentlicher Teil ist der des
// Zertifikats (LH-FA-23.a *Dateien*).
func tlsSchluessel(c *Command, _ gelesen, _ *datei, q quellen) error {
	pfad := q.wert["tls-key"]
	if pfad == "" {
		return nil
	}
	roh, err := liesDatei("--tls-key", pfad)
	if err != nil {
		return err
	}
	schluessel, err := privaterSchluessel(roh)
	if err != nil {
		return err
	}
	z := zugang(c)
	if !passt(z.Zertifikat, schluessel) {
		return dateiUngueltig("--tls-key", "der Schlüssel passt nicht zum Zertifikat")
	}
	z.Schluessel = Schluessel(roh)
	return nil
}

// passt sagt, ob der öffentliche Schlüssel des ersten Zertifikats in der
// Zertifikatsdatei der öffentliche Schlüssel ist, den der private liefert.
func passt(zertifikatsdatei string, schluessel crypto.PublicKey) bool {
	block, _ := pem.Decode([]byte(zertifikatsdatei))
	if block == nil {
		return false
	}
	server, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false
	}
	gleich, ok := server.PublicKey.(interface{ Equal(x crypto.PublicKey) bool })
	return ok && gleich.Equal(schluessel)
}

// privaterSchluessel liefert den öffentlichen Teil des privaten Schlüssels in roh.
func privaterSchluessel(roh []byte) (crypto.PublicKey, error) {
	block, rest := pem.Decode(roh)
	if block == nil {
		return nil, dateiUngueltig("--tls-key", "kein PEM-Block")
	}
	if block.Type == "ENCRYPTED PRIVATE KEY" || block.Headers["Proc-Type"] == "4,ENCRYPTED" {
		return nil, dateiUngueltig("--tls-key", "der Schlüssel ist verschlüsselt")
	}
	var parse func([]byte) (any, error)
	switch block.Type {
	case "PRIVATE KEY":
		parse = x509.ParsePKCS8PrivateKey
	case "RSA PRIVATE KEY":
		parse = func(b []byte) (any, error) { return x509.ParsePKCS1PrivateKey(b) }
	case "EC PRIVATE KEY":
		parse = func(b []byte) (any, error) { return x509.ParseECPrivateKey(b) }
	default:
		return nil, dateiUngueltig("--tls-key", "der PEM-Block hat keinen der Typen PRIVATE KEY, RSA PRIVATE KEY und EC PRIVATE KEY")
	}
	if weiterer, _ := pem.Decode(rest); weiterer != nil {
		return nil, dateiUngueltig("--tls-key", "mehr als ein PEM-Block")
	}
	schluessel, err := parse(block.Bytes)
	if err != nil {
		return nil, dateiUngueltig("--tls-key", "der PEM-Block enthält keinen lesbaren Schlüssel")
	}
	switch k := schluessel.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey, nil
	case *ecdsa.PrivateKey:
		return &k.PublicKey, nil
	case ed25519.PrivateKey:
		return k.Public(), nil
	}
	return nil, dateiUngueltig("--tls-key", "der Schlüssel ist weder RSA noch ECDSA noch Ed25519")
}
