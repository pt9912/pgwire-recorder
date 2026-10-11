// Package testzeit gibt den Tests des Kerns Warten, Frist und Zeitgrenzen in ganzen
// Millisekunden. Es ist die einzige Stelle dieser Tests, die die Zeitbibliothek
// einbindet; Domain Model, Ports und Services binden sie nicht ein (LH-FA-21.a *Uhr*).
package testzeit
