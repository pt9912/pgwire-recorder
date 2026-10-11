// Package testpki erzeugt zur Testzeit Zertifizierungsstellen und
// Serverzertifikate für die Tests von TLS; nichts davon wird eingecheckt.
// Kein Paket unter `./cmd/...` importiert es, auch nicht mittelbar
// (`TestBinaryOhneTesthilfen` im Paket `internal/bootstrap`); für Testdateien
// außerhalb von `./cmd/...` prüft das nichts.
package testpki
