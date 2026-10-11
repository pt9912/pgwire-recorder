// Package tlsproxy ist ein TLS-Proxy vor einer PostgreSQL-Instanz für die
// Integrationstests von TLS beim Einspielen. Er liegt im Wurzelverzeichnis der
// Verdrahtung, weil allein dort und in den PGWire-Adaptern die Technik des TLS
// (`crypto/tls`) importiert werden darf (`.a-check.yml`); nur Testcode importiert
// das Paket.
package tlsproxy
