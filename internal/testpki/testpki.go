package testpki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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

// Blatt beschreibt ein Serverzertifikat: die DNS-Namen und IP-Adressen, auf
// die es ausgestellt ist, und seine Gültigkeit; ohne Gültigkeit gilt es von
// einer Stunde vor bis eine Stunde nach jetzt.
type Blatt struct {
	DNS      []string
	IPs      []net.IP
	Von, Bis time.Time
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
type Blattzertifikat struct {
	DER        []byte
	Schluessel *ecdsa.PrivateKey
}

// Server stellt ein Serverzertifikat nach b aus.
func (ca *CA) Server(t testing.TB, b Blatt) Blattzertifikat {
	t.Helper()
	schluessel, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	von, bis := b.Von, b.Bis
	if von.IsZero() {
		von = time.Now().Add(-time.Hour)
	}
	if bis.IsZero() {
		bis = time.Now().Add(time.Hour)
	}
	vorlage := &x509.Certificate{
		SerialNumber: seriennummer(t),
		Subject:      pkix.Name{CommonName: "server"},
		NotBefore:    von,
		NotAfter:     bis,
		DNSNames:     b.DNS,
		IPAddresses:  b.IPs,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, vorlage, ca.Zertifikat, &schluessel.PublicKey, ca.schluessel)
	if err != nil {
		t.Fatal(err)
	}
	return Blattzertifikat{DER: der, Schluessel: schluessel}
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
