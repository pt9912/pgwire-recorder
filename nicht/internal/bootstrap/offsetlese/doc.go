// Package offsetlese liefert den Integrationstests die Werte von offset_ms einer
// Aufzeichnung beider Formate. Es liegt im Wurzelverzeichnis der Verdrahtung, weil
// dort der Recording-Adapter importiert werden darf, die Integrationstests unter
// `test/` aber keine Adapter (`.a-check.yml`). Kein Paket unter `./cmd/...`
// importiert es, auch nicht mittelbar (`TestBinaryOhneTesthilfen`).
package offsetlese
