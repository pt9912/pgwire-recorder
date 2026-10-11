package cli

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/fs"
	"os"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// upstreamCA liest die Datei aus --upstream-ca und legt ihre Zertifikate in
// den Optionen von play ab. Es läuft als letzte Prüfung des Starts, nach allen
// Prüfungen von upstreamPlay und vor dem Laden der Aufzeichnung (LH-FA-17.a
// *Fehler*, LH-FA-20.a *Start*); ohne die Option liest es nichts.
func upstreamCA(c *Command, _ gelesen, _ *datei, q quellen) error {
	pfad := q.wert["upstream-ca"]
	if pfad == "" {
		return nil
	}
	z, err := liesZertifizierungsstelle(pfad)
	if err != nil {
		return err
	}
	c.Play.UpstreamCA = z
	return nil
}

// liesZertifizierungsstelle liest die PEM-Datei pfad (LH-FA-20.a *Start*): Sie
// ist, Links gefolgt, eine reguläre Datei, geprüft vor dem Öffnen, damit der
// Start nie auf eine FIFO wartet. Sie enthält mindestens einen PEM-Block, jeder
// Block hat den Typ CERTIFICATE und ein lesbares X.509-Zertifikat; Text
// außerhalb der Blöcke bleibt unbeachtet. Jeder Fehler ist PGR-E2007 und nennt
// die Option und die Ursache, weder den Pfad noch den Inhalt (LH-RB-01).
func liesZertifizierungsstelle(pfad string) (Zertifikate, error) {
	info, err := os.Stat(pfad)
	if err != nil {
		return nil, caNichtLesbar(err)
	}
	if !info.Mode().IsRegular() {
		return nil, model.Errorf(model.CodeConfigCA, nil, "--upstream-ca: Datei nicht lesbar: keine reguläre Datei")
	}
	roh, err := os.ReadFile(pfad)
	if err != nil {
		return nil, caNichtLesbar(err)
	}
	return zertifikate(roh)
}

// zertifikate liest die PEM-Blöcke von roh als X.509-Zertifikate.
func zertifikate(roh []byte) (Zertifikate, error) {
	var out Zertifikate
	for {
		var block *pem.Block
		block, roh = pem.Decode(roh)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			return nil, caUngueltig("ein PEM-Block hat nicht den Typ CERTIFICATE")
		}
		c, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, caUngueltig("ein PEM-Block enthält kein lesbares X.509-Zertifikat")
		}
		out = append(out, c)
	}
	if len(out) == 0 {
		return nil, caUngueltig("kein PEM-Block")
	}
	return out, nil
}

// caNichtLesbar ist der Fehler einer Datei, die fehlt oder nicht gelesen
// werden kann; er trägt den Grund des Betriebssystems ohne den Pfad.
func caNichtLesbar(err error) error {
	if errors.Is(err, fs.ErrNotExist) {
		return model.Errorf(model.CodeConfigCA, nil, "--upstream-ca: Datei nicht vorhanden")
	}
	var pe *fs.PathError
	if errors.As(err, &pe) {
		err = pe.Err
	}
	return model.Errorf(model.CodeConfigCA, err, "--upstream-ca: Datei nicht lesbar")
}

// caUngueltig ist PGR-E2007 für eine lesbare Datei, die kein gültiges PEM mit
// Zertifikaten ist; der Grund nennt keinen Inhalt.
func caUngueltig(grund string) error {
	return model.Errorf(model.CodeConfigCA, nil, "--upstream-ca: %s", grund)
}
