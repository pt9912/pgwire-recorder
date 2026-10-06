package services_test

import (
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// Abdeckung: LH-FA-09/Negative — eine Anfrage ist eine Lebendprüfung, wenn ihr
// Text nur aus Leerraum, Zeilenkommentaren bis Zeilen- oder Textende und
// geschlossenen, auch verschachtelten Blockkommentaren besteht, der leere Text
// eingeschlossen; ein offener Blockkommentar, `;`, Leerraum außerhalb von ASCII
// und jedes andere Zeichen vor, zwischen oder nach Kommentaren machen sie zu
// keiner. Die Fälle ohne \v gelten gleich, ob \v Leerraum ist oder nicht.
func TestLebendpruefungErkennung(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want bool
	}{
		// Leerraum
		{"", true},
		{" ", true},
		{"\t", true},
		{"\n", true},
		{"\r", true},
		{"\f", true},
		{" \t\n\r\f ", true},
		{" ", false},
		{"\x00", false},
		{"\b", false},
		// Zeilenkommentar
		{"-- ping", true},
		{"--", true},
		{"--ping", true},
		{"-- ping\n", true},
		{"-- a\n-- b", true},
		{"-- ping\nSELECT 1", false},
		{"-- ping\rSELECT 1", false},
		{"-- ping\fSELECT 1", true},
		{"-- SELECT 1", true},
		{"-- a /* b", true},
		{"-", false},
		{"- - ping", false},
		{" -- ping", true},
		// Blockkommentar
		{"/**/", true},
		{"/* ping */", true},
		{"/* a */ /* b */", true},
		{"/* a /* b */ c */", true},
		{"/* /* /* */ */ */", true},
		{"/* a -- b */", true},
		{"/* -- */ SELECT", false},
		{"/* ping", false},
		{"/*", false},
		{"/*/", false},
		{"/* a /* b */", false},
		{"/* a */ */", false},
		{"*/", false},
		{"/* a **/", true},
		{"/* a */*", false},
		{"/", false},
		{"/ * a */", false},
		// gemischt
		{"/* a */ -- b\n\t/* c */", true},
		{"-- a\r/* b */", true},
		{"-- a\n/* b", false},
		// Anweisung oder Trennzeichen
		{";", false},
		{" ; ", false},
		{"-- ping\n;", false},
		{"/* */;", false},
		{"SELECT 1", false},
		{"x-- ping", false},
		{"-- ping\nx", false},
	} {
		for _, vt := range []bool{false, true} {
			if got := services.IstLebendpruefung(tc.sql, vt); got != tc.want {
				t.Errorf("istLebendpruefung(%q, %v) = %v, erwartet %v", tc.sql, vt, got, tc.want)
			}
		}
	}
}

// Abdeckung: LH-FA-09/Negative — der vertikale Tabulator (\v) ist Leerraum
// einer Lebendprüfung nur, wenn er es nach der Serverversion ist (vt); in einem
// Kommentar zählt er in beiden Varianten zum Kommentar.
func TestLebendpruefungVertikalerTabulator(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		ohneVT bool
		mitVT  bool
	}{
		{"\v", false, true},
		{" \t\n\r\f\v ", false, true},
		{"\v-- ping\v", false, true},
		{"/* a */\v", false, true},
		{"-- ping\vSELECT 1", true, true},
		{"/* \v */", true, true},
		{"\v;", false, false},
		{"\v/* offen", false, false},
	} {
		if got := services.IstLebendpruefung(tc.sql, false); got != tc.ohneVT {
			t.Errorf("istLebendpruefung(%q, false) = %v, erwartet %v", tc.sql, got, tc.ohneVT)
		}
		if got := services.IstLebendpruefung(tc.sql, true); got != tc.mitVT {
			t.Errorf("istLebendpruefung(%q, true) = %v, erwartet %v", tc.sql, got, tc.mitVT)
		}
	}
}

// Abdeckung: LH-FA-09/Negative — \v ist Leerraum ab der Hauptversion 17 in
// server_version, gelesen aus den führenden Ziffern; fehlt der Parameter, ist
// er leer, beginnt er nicht mit einer Ziffer oder ist die Zahl nicht lesbar,
// ist \v kein Leerraum.
func TestLebendpruefungServerversion(t *testing.T) {
	for _, tc := range []struct {
		version string
		fehlt   bool
		want    bool
	}{
		{"17.0", false, true},
		{"17", false, true},
		{"18.6", false, true},
		{"18beta1", false, true},
		{"17.2 (Debian 17.2-1.pgdg120+1)", false, true},
		{"100.1", false, true},
		{"16.15", false, false},
		{"16", false, false},
		{"14.24", false, false},
		{"9.6.24", false, false},
		{"1", false, false},
		{"", false, false},
		{"", true, false},
		{"v17", false, false},
		{" 17", false, false},
		{"99999999999999999999.0", false, false},
	} {
		params := map[string]string{"client_encoding": "UTF8"}
		if !tc.fehlt {
			params["server_version"] = tc.version
		}
		if got := services.VTLeerraum(params); got != tc.want {
			t.Errorf("vtLeerraum(server_version %q, fehlt %v) = %v, erwartet %v", tc.version, tc.fehlt, got, tc.want)
		}
	}
}
