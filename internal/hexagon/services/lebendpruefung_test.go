package services

import "testing"

// Abdeckung: LH-FA-09/Negative — eine Anfrage ist eine Lebendprüfung, wenn ihr
// Text nur aus Leerraum, Zeilenkommentaren bis Zeilen- oder Textende und
// geschlossenen, auch verschachtelten Blockkommentaren besteht, der leere Text
// eingeschlossen; ein offener Blockkommentar, `;`, Leerraum außerhalb von ASCII
// und jedes andere Zeichen vor, zwischen oder nach Kommentaren machen sie zu
// keiner.
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
		{"\v", true},
		{" \t\n\r\f\v ", true},
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
		{"-- ping\vSELECT 1", true},
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
		if got := istLebendpruefung(tc.sql); got != tc.want {
			t.Errorf("istLebendpruefung(%q) = %v, erwartet %v", tc.sql, got, tc.want)
		}
	}
}
