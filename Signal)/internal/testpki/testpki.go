package testpki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"
)

// CA ist eine Zertifizierungsstelle, die ein Test erzeugt hat.
type CA struct {
	Zertifikat *x509.Certificate
	schluessel *ecdsa.PrivateKey
}

// Art ist die Schlüsselart eines Serverzertifikats; die Voreinstellung ist ECDSA.
type Art int

// Die Schlüsselarten.
const (
	ECDSA Art = iota
	RSA
	Ed25519
)

// Blatt beschreibt ein Serverzertifikat: die DNS-Namen und IP-Adressen, auf
// die es ausgestellt ist, seine Gültigkeit und die Schlüsselart; ohne
// Gültigkeit gilt es von einer Stunde vor bis eine Stunde nach jetzt. Ohne
// Verwendung ist es für Server ausgestellt, ohne Schlüsselnutzung für
// Signaturen.
type Blatt struct {
	DNS               []string
	IPs               []net.IP
	Von, Bis          time.Time
	Art               Art
	Verwendung        []x509.ExtKeyUsage
	Schluesselnutzung x509.KeyUsage
}

// NeueCA erzeugt eine Zertifizierungsstelle mit dem Namen name.
func NeueCA(t testing.TB, name string) *CA {
	t.Helper()
	schluessel, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	vorlage := &x509.Certificate{
		SerialNumber:          seriennummer(t),
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, vorlage, vorlage, &schluessel.PublicKey, schluessel)
	if err != nil {
		t.Fatal(err)
	}
	zertifikat, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return &CA{Zertifikat: zertifikat, schluessel: schluessel}
}

// PEM ist das Zertifikat der Zertifizierungsstelle als PEM-Block.
func (ca *CA) PEM() []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.Zertifikat.Raw})
}

// Blattzertifikat ist ein ausgestelltes Serverzertifikat: das Zertifikat in
// DER-Form und sein Schlüssel. Der Test baut daraus die Zertifikatsform seines
// TLS-Servers, damit dieses Paket keine Technik des TLS braucht.
//
// Schluessel ist der ECDSA-Schlüssel und nil bei einer anderen Art; Privat ist
// der Schlüssel jeder Art.
type Blattzertifikat struct {
	DER        []byte
	Schluessel *ecdsa.PrivateKey
	Privat     crypto.Signer
}

// Server stellt ein Serverzertifikat nach b aus.
func (ca *CA) Server(t testing.TB, b Blatt) Blattzertifikat {
	t.Helper()
	privat := neuerSchluessel(t, b.Art)
	von, bis := b.Von, b.Bis
	if von.IsZero() {
		von = time.Now().Add(-time.Hour)
	}
	if bis.IsZero() {
		bis = time.Now().Add(time.Hour)
	}
	verwendung, nutzung := b.Verwendung, b.Schluesselnutzung
	if verwendung == nil {
		verwendung = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}
	if nutzung == 0 {
		nutzung = x509.KeyUsageDigitalSignature
	}
	vorlage := &x509.Certificate{
		SerialNumber: seriennummer(t),
		Subject:      pkix.Name{CommonName: "server"},
		NotBefore:    von,
		NotAfter:     bis,
		DNSNames:     b.DNS,
		IPAddresses:  b.IPs,
		KeyUsage:     nutzung,
		ExtKeyUsage:  verwendung,
	}
	der, err := x509.CreateCertificate(rand.Reader, vorlage, ca.Zertifikat, privat.Public(), ca.schluessel)
	if err != nil {
		t.Fatal(err)
	}
	ec, _ := privat.(*ecdsa.PrivateKey)
	return Blattzertifikat{DER: der, Schluessel: ec, Privat: privat}
}

// neuerSchluessel erzeugt einen Schlüssel der Art art.
func neuerSchluessel(t testing.TB, art Art) crypto.Signer {
	t.Helper()
	var (
		k   crypto.Signer
		err error
	)
	switch art {
	case RSA:
		k, err = rsa.GenerateKey(rand.Reader, 2048)
	case Ed25519:
		_, k, err = ed25519.GenerateKey(rand.Reader)
	default:
		k, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// PEM ist das Zertifikat als PEM-Block.
func (b Blattzertifikat) PEM() []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b.DER})
}

// SchluesselPEM ist der Schlüssel als PEM-Block vom Typ typ: PRIVATE KEY
// (PKCS #8, jede Art), RSA PRIVATE KEY (PKCS #1, nur RSA) oder EC PRIVATE KEY
// (SEC 1, nur ECDSA).
func (b Blattzertifikat) SchluesselPEM(t testing.TB, typ string) []byte {
	t.Helper()
	var (
		roh []byte
		err error
	)
	switch k := b.Privat.(type) {
	case *rsa.PrivateKey:
		if typ == "RSA PRIVATE KEY" {
			roh = x509.MarshalPKCS1PrivateKey(k)
		}
	case *ecdsa.PrivateKey:
		if typ == "EC PRIVATE KEY" {
			roh, err = x509.MarshalECPrivateKey(k)
		}
	}
	if typ == "PRIVATE KEY" {
		roh, err = x509.MarshalPKCS8PrivateKey(b.Privat)
	}
	if err != nil || roh == nil {
		t.Fatalf("Schlüssel nicht als %s zu schreiben: %v", typ, err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: roh})
}

// seriennummer ist eine zufällige Seriennummer.
func seriennummer(t testing.TB) *big.Int {
	t.Helper()
	n, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 62))
	if err != nil {
		t.Fatal(err)
	}
	return n
}
