package cli

import (
	"errors"
	"os"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// errUpstreamForm ist der Grund eines Werts von --upstream, der weder Name
// einer Verbindung der gewählten Datei ist noch die Form host:port hat; er
// nennt den Wert nicht.
var errUpstreamForm = errors.New("weder Name einer Verbindung der Konfigurationsdatei noch host:port")

// upstreamRecord prüft nach der Zusammenführung jeden gesetzten Wert von
// --upstream (pruefeUpstream). Nennt der zusammengeführte Wert eine
// Verbindung, ersetzt ihn deren Adresse (adresseRecord), mit den Variablen,
// wie sie jetzt gesetzt sind; host:port bleibt, wie geschrieben.
func upstreamRecord(c *Command, g gelesen, d *datei) error {
	if err := pruefeUpstream(g, d); err != nil {
		return err
	}
	v, ok := d.verbindung(c.Record.Upstream)
	if !ok {
		return nil
	}
	adresse, err := v.adresseRecord(os.Getenv)
	if err != nil {
		return err
	}
	c.Record.Upstream = adresse
	return nil
}

// upstreamPlay prüft wie upstreamRecord. Nennt der zusammengeführte Wert eine
// Verbindung, ersetzt ihn deren Adresse (zielPlay), und Benutzer und
// Datenbank der Verbindung gelten, wo --user beziehungsweise --database
// nichts setzt (LH-FA-17.a *Wirkung einer URL*).
func upstreamPlay(c *Command, g gelesen, d *datei) error {
	if err := pruefeUpstream(g, d); err != nil {
		return err
	}
	v, ok := d.verbindung(c.Play.Upstream)
	if !ok {
		return nil
	}
	z, err := v.zielPlay(os.Getenv)
	if err != nil {
		return err
	}
	c.Play.Upstream = z.adresse
	if c.Play.User == "" {
		c.Play.User = z.benutzer
	}
	if c.Play.Database == "" {
		c.Play.Database = z.datenbank
	}
	return nil
}

// pruefeUpstream prüft jeden gesetzten Wert von --upstream, den der
// Kommandozeile vor dem der Umgebungsvariable, auch wenn er nicht gilt: Er
// ist der Name einer Verbindung der gewählten Datei oder hat die Form
// host:port (hostPortForm), sonst ist er PGR-E2001 ohne den Wert
// (LH-FA-17.a). Ohne Datei gibt es keine Namen. Den Schlüssel upstream prüft
// das Laden.
func pruefeUpstream(g gelesen, d *datei) error {
	if g.cliOk && !upstreamGueltig(g.cli, d) {
		return model.Errorf(model.CodeUsage, errUpstreamForm, "--upstream")
	}
	if g.envOk && !upstreamGueltig(g.env, d) {
		return model.Errorf(model.CodeUsage, errUpstreamForm, "Umgebungsvariable %s", envName("upstream"))
	}
	return nil
}

// upstreamGueltig meldet, ob w der Name einer Verbindung von d ist oder die
// Form host:port hat.
func upstreamGueltig(w string, d *datei) bool {
	return d.hatName(w) || hostPortForm(w)
}
