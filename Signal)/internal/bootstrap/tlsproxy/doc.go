// Package tlsproxy ist ein TLS-Proxy vor einer PostgreSQL-Instanz für die
// Integrationstests von TLS beim Einspielen. Er liegt im Wurzelverzeichnis der
// Verdrahtung, weil allein dort und in den PGWire-Adaptern die Technik des TLS
// (`crypto/tls`) importiert werden darf (`.a-check.yml`). Kein Paket unter
// `./cmd/...` importiert es, auch nicht mittelbar (`TestBinaryOhneTesthilfen`);
// für Testdateien außerhalb von `./cmd/...` prüft das nichts.
package tlsproxy
