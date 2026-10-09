package cli_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// variablen setzt je Name den Wert für die Dauer des Tests; nil heißt nicht
// gesetzt, auch nicht leer.
func variablen(t *testing.T, werte map[string]*string) {
	t.Helper()
	for name, wert := range werte {
		t.Setenv(name, "")
		if wert == nil {
			if err := os.Unsetenv(name); err != nil {
				t.Fatal(err)
			}
			continue
		}
		t.Setenv(name, *wert)
	}
}

// w ist ein gesetzter Wert für variablen.
func w(s string) *string { return &s }

// mitUpstream liest record mit --upstream=v und der Datei mit der einen
// Verbindung v (mitVerbindung) und liefert die Adresse.
func mitUpstream(t *testing.T, url string) (string, error) {
	t.Helper()
	return recordUpstream("--upstream=v", "--config="+schreibe(t, mitVerbindung(url)))
}

// Abdeckung: LH-FA-17/Happy — nennt der zusammengeführte Wert von --upstream,
// PGWIRE_RECORDER_UPSTREAM oder des Schlüssels upstream eine Verbindung, ist
// die Adresse von record Host und Port dieser Verbindung: der Host dekodiert,
// IPv6 wieder in eckigen Klammern, auch mit Zone, ohne Port 5432, der Port wie
// geschrieben; Benutzer, Passwort und Datenbank und ihre Variablen bleiben
// unbeachtet (LH-FA-17.a *Benannte Verbindungen*).
func TestUpstreamVerbindung(t *testing.T) {
	leere(t, "record")
	variablen(t, map[string]*string{"PGR_T_U": nil, "PGR_T_PW": nil, "PGR_T_DB": nil})
	for url, want := range map[string]string{
		"postgresql://u@h.example:6543/db":                           "h.example:6543",
		"postgresql://h/db":                                          "h:5432",
		"postgresql://h:05432/db":                                    "h:05432",
		"postgresql://h%41/db":                                       "hA:5432",
		"postgresql://[::1]/db":                                      "[::1]:5432",
		"postgresql://[::1]:6/db":                                    "[::1]:6",
		"postgresql://[fe80::1%25eth0]:5/db":                         "[fe80::1%eth0]:5",
		"postgresql://${PGR_T_U}:${PGR_T_PW}@h:7/${PGR_T_DB}":        "h:7",
		"postgresql://u${PGR_T_U}@h/a${PGR_T_DB}?sslmode=disable":    "h:5432",
		"postgresql://h$$x/db":                                       "h$x:5432",
		"postgresql://127.0.0.1:5/db":                                "127.0.0.1:5",
		"postgresql://${PGR_T_U}@h/db%2F${PGR_T_DB}?sslmode=disable": "h:5432",
	} {
		if got, err := mitUpstream(t, url); err != nil || got != want {
			t.Errorf("%s: %q, %v, erwartet %q", url, got, err, want)
		}
	}
	datei := "connections:\n  a: postgresql://a.example/db\n  b: postgresql://b.example:2/db\n"
	pfad := schreibe(t, datei)
	for args, want := range map[string]string{
		"--upstream=b":   "b.example:2",
		"--upstream=a":   "a.example:5432",
		"--upstream=k:1": "k:1",
	} {
		if got, err := recordUpstream(args, "--config="+pfad); err != nil || got != want {
			t.Errorf("%s: %q, %v, erwartet %q", args, got, err, want)
		}
	}
	t.Setenv(envUpstream, "a")
	if got, err := recordUpstream("--upstream=b", "--config="+pfad); err != nil || got != "b.example:2" {
		t.Errorf("--upstream=b neben %s=a: %q, %v", envUpstream, got, err)
	}
	if got, err := recordUpstream("--config=" + pfad); err != nil || got != "a.example:5432" {
		t.Errorf("%s=a: %q, %v", envUpstream, got, err)
	}
	t.Setenv(envUpstream, "")
	if got, err := recordUpstream("--config=" + schreibe(t, datei+"record:\n  upstream: b\n")); err != nil || got != "b.example:2" {
		t.Errorf("Schlüssel upstream: b: %q, %v", got, err)
	}
}

// Abdeckung: LH-FA-17/Boundary — die Variablen in Host und Port der benutzten
// Verbindung werden einmal eingesetzt: Der Wert steht unverändert in seinem
// Teil, weder dekodiert noch erneut ausgewertet, auch mit @, /, % oder ${;
// ein Host mit : steht danach in eckigen Klammern, auch einer, der keine
// IPv6-Adresse ist; der Port gilt, wie er eingesetzt ist, auch mit führenden
// Nullen und neben wörtlichen Ziffern (LH-FA-17.a *Geheimnisse*, U4, U5).
func TestUpstreamEinsetzen(t *testing.T) {
	leere(t, "record")
	for _, fall := range []struct {
		url, h, p, y, want string
	}{
		{"postgresql://${PGR_T_H}/db", "::1", "", "", "[::1]:5432"},
		{"postgresql://${PGR_T_H}/db", "a:b", "", "", "[a:b]:5432"},
		{"postgresql://${PGR_T_H}/db", "h", "", "", "h:5432"},
		{"postgresql://${PGR_T_H}/db", "fe80::1%eth0", "", "", "[fe80::1%eth0]:5432"},
		{"postgresql://${PGR_T_H}/db", "x${PGR_T_Y}", "", "gesetzt", "x${PGR_T_Y}:5432"},
		{"postgresql://${PGR_T_H}/db", "a@b/%41", "", "", "a@b/%41:5432"},
		{"postgresql://vor${PGR_T_H}nach/db", "-", "", "", "vor-nach:5432"},
		{"postgresql://h:${PGR_T_P}/db", "", "05432", "", "h:05432"},
		{"postgresql://h:5${PGR_T_P}/db", "", "5432", "", "h:55432"},
		{"postgresql://h:${PGR_T_P}/db", "", "65535", "", "h:65535"},
		{"postgresql://${PGR_T_H}:${PGR_T_P}/db", "db.example", "1", "", "db.example:1"},
	} {
		variablen(t, map[string]*string{"PGR_T_H": w(fall.h), "PGR_T_P": w(fall.p), "PGR_T_Y": w(fall.y)})
		if got, err := mitUpstream(t, fall.url); err != nil || got != fall.want {
			t.Errorf("%s mit H=%q P=%q: %q, %v, erwartet %q", fall.url, fall.h, fall.p, got, err, fall.want)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — eine nicht gesetzte oder leere Variable in
// Host oder Port der benutzten Verbindung ist PGR-E2005 an connections.<Name>
// mit dem Namen der ersten in der Reihenfolge der URL und ohne einen Wert, im
// Port vor dessen Form; die Variablen einer nicht benutzten Verbindung bleiben
// unbeachtet, und config show meldet kein PGR-E2005 (LH-FA-17.a, A12, U3).
func TestUpstreamVariableFehlt(t *testing.T) {
	leere(t, "record")
	fehlt := func(name string) string {
		return "Konfiguration [PGR-E2005]: Konfigurationsdatei: connections.v: Umgebungsvariable " + name + " eines Platzhalters nicht gesetzt"
	}
	for _, fall := range []struct {
		url   string
		h, p  *string
		erste string
	}{
		{"postgresql://${PGR_T_H}/db", nil, nil, "PGR_T_H"},
		{"postgresql://${PGR_T_H}/db", w(""), nil, "PGR_T_H"},
		{"postgresql://h:${PGR_T_P}/db", nil, nil, "PGR_T_P"},
		{"postgresql://h:${PGR_T_P}/db", nil, w(""), "PGR_T_P"},
		{"postgresql://${PGR_T_H}:${PGR_T_P}/db", nil, nil, "PGR_T_H"},
		{"postgresql://${PGR_T_H}:${PGR_T_P}/db", w("GEHEIM"), nil, "PGR_T_P"},
		{"postgresql://h${PGR_T_P}${PGR_T_H}/db", w("GEHEIM"), nil, "PGR_T_P"},
	} {
		variablen(t, map[string]*string{"PGR_T_H": fall.h, "PGR_T_P": fall.p})
		_, err := mitUpstream(t, fall.url)
		if err == nil || err.Error() != fehlt(fall.erste) || !hatCode(err, model.CodeConfigVariable) {
			t.Errorf("%s: %v, erwartet %q", fall.url, err, fehlt(fall.erste))
		}
	}
	variablen(t, map[string]*string{"PGR_T_H": nil})
	andere := "connections:\n  w: postgresql://${PGR_T_H}:${PGR_T_H}/db\n  v: postgresql://h/db\n"
	if got, err := recordUpstream("--upstream=v", "--config="+schreibe(t, andere)); err != nil || got != "h:5432" {
		t.Errorf("Variable einer nicht benutzten Verbindung: %q, %v", got, err)
	}
	zeigen := schreibe(t, "connections:\n  v: postgresql://${PGR_T_H}/db\nrecord:\n  upstream: v\n")
	if got, err := konfigurationZeigen(t, "--config", zeigen); err != nil || !strings.Contains(got, "${PGR_T_H}") {
		t.Errorf("config show mit nicht gesetzter Variable der benutzten Verbindung: %q, %v", got, err)
	}
}

// Abdeckung: LH-FA-17/Negative — hat der Port der benutzten Verbindung nach
// dem Einsetzen nicht die Form eines Ports, ist das PGR-E2004 an
// connections.<Name> ohne den Wert (LH-FA-17.a *Geheimnisse*, Rückgabe 8, U3).
func TestUpstreamPortNachEinsetzen(t *testing.T) {
	leere(t, "record")
	genau := "Konfiguration [PGR-E2004]: Konfigurationsdatei: connections.v: Port nach dem Einsetzen ist keine Zahl von 1 bis 65535"
	for _, fall := range []struct{ url, p string }{
		{"postgresql://h:${PGR_T_P}/db", "GEHEIM"},
		{"postgresql://h:${PGR_T_P}/db", "0"},
		{"postgresql://h:${PGR_T_P}/db", "65536"},
		{"postgresql://h:${PGR_T_P}/db", "5 GEHEIM"},
		{"postgresql://h:${PGR_T_P}/db", "-5"},
		{"postgresql://h:${PGR_T_P}/db", "99999999999999999999"},
		{"postgresql://h:9${PGR_T_P}/db", "0000"},
		{"postgresql://h:${PGR_T_P}/db", "5\x01"},
	} {
		variablen(t, map[string]*string{"PGR_T_P": w(fall.p)})
		_, err := mitUpstream(t, fall.url)
		if err == nil || err.Error() != genau || !istDatei(err) {
			t.Errorf("%s mit P=%q: %v, erwartet %q", fall.url, fall.p, err, genau)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — sslmode=require der benutzten Verbindung ist
// bei record PGR-E2004 an connections.<Name>; das einer nicht benutzten nicht.
// Am Ende gilt die Reihenfolge --upstream, sslmode=require, Variablen, Port
// (LH-FA-17.a *Fehler*, A9).
func TestUpstreamReihenfolgeAmEnde(t *testing.T) {
	leere(t, "record")
	ssl := "Konfiguration [PGR-E2004]: Konfigurationsdatei: connections.v: sslmode=require ist bei record ungültig, record verbindet ohne TLS zum Upstream"
	variablen(t, map[string]*string{"PGR_T_H": nil, "PGR_T_P": w("x")})
	for _, url := range []string{
		"postgresql://h/db?sslmode=require",
		"postgresql://${PGR_T_H}:${PGR_T_P}/db?sslmode=require",
	} {
		if _, err := mitUpstream(t, url); err == nil || err.Error() != ssl {
			t.Errorf("%s: %v, erwartet %q", url, err, ssl)
		}
	}
	zwei := schreibe(t, "connections:\n  v: postgresql://h/db?sslmode=require\n  w: postgresql://k/db\n")
	if got, err := recordUpstream("--upstream=w", "--config="+zwei); err != nil || got != "k:5432" {
		t.Errorf("sslmode=require einer nicht benutzten Verbindung: %q, %v", got, err)
	}
	t.Setenv(envUpstream, "GEHEIM")
	if _, err := recordUpstream("--upstream=v", "--config="+zwei); !istUsage(err) || !strings.Contains(err.Error(), envUpstream) {
		t.Errorf("--upstream vor sslmode=require: %v", err)
	}
	t.Setenv(envUpstream, "")
	if _, err := mitUpstream(t, "postgresql://${PGR_T_H}:${PGR_T_P}/db"); !hatCode(err, model.CodeConfigVariable) {
		t.Errorf("Variablen vor dem Port: %v", err)
	}
}

// Abdeckung: LH-FA-17/Boundary — die Variablen der benutzten Verbindung liest
// der Start einmal, nach der Zusammenführung: Eine spätere Änderung ändert die
// gelesene Adresse nicht (LH-FA-17.a, U2).
func TestUpstreamEinmalGelesen(t *testing.T) {
	leere(t, "record")
	variablen(t, map[string]*string{"PGR_T_H": w("erst")})
	cmd, err := lese("record", "--listen=x", "--output=r.yaml", "--upstream=v", "--config="+schreibe(t, mitVerbindung("postgresql://${PGR_T_H}/db")))
	if err != nil || cmd.Record.Upstream != "erst:5432" {
		t.Fatalf("%q, %v", cmd.Record.Upstream, err)
	}
	t.Setenv("PGR_T_H", "danach")
	if cmd.Record.Upstream != "erst:5432" {
		t.Errorf("Adresse nach der Änderung: %q", cmd.Record.Upstream)
	}
}

// Abdeckung: LH-FA-17/Boundary, LH-FA-17/Negative — das Einsetzen gilt für
// jeden Teil der URL: in der Reihenfolge der URL, Benutzer, Passwort, Host,
// Port, Datenbank, je Teil mit seinem wörtlichen Text; die erste nicht
// gesetzte Variable in dieser Reihenfolge ist PGR-E2005 (LH-FA-17.a, U8).
func TestEinsetzenAlleTeile(t *testing.T) {
	pfad := schreibe(t, mitVerbindung("postgresql://a${PGR_T_U}b:${PGR_T_PW}@${PGR_T_H}:${PGR_T_P}/d${PGR_T_DB}"))
	alle := map[string]string{"PGR_T_U": "u", "PGR_T_PW": "p@w", "PGR_T_H": "h", "PGR_T_P": "5", "PGR_T_DB": "db/%41"}
	got, err := cli.SetzeEin(pfad, func(n string) string { return alle[n] })
	if want := []string{"aub", "p@w", "h", "5", "ddb/%41"}; err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("alle gesetzt: %q, %v, erwartet %q", got, err, want)
	}
	for _, name := range []string{"PGR_T_U", "PGR_T_PW", "PGR_T_H", "PGR_T_P", "PGR_T_DB"} {
		ohne := func(n string) string {
			if n == name {
				return ""
			}
			return alle[n]
		}
		_, err := cli.SetzeEin(pfad, ohne)
		if !hatCode(err, model.CodeConfigVariable) || !strings.HasSuffix(err.Error(), "Umgebungsvariable "+name+" eines Platzhalters nicht gesetzt") {
			t.Errorf("ohne %s: %v", name, err)
		}
	}
	nichts := func(string) string { return "" }
	if _, err := cli.SetzeEin(pfad, nichts); err == nil || !strings.Contains(err.Error(), "PGR_T_U eines") {
		t.Errorf("keine gesetzt: %v, erwartet die erste, PGR_T_U", err)
	}
	ohneBenutzer := schreibe(t, mitVerbindung("postgresql://${PGR_T_H}/${PGR_T_DB}"))
	if _, err := cli.SetzeEin(ohneBenutzer, nichts); err == nil || !strings.Contains(err.Error(), "PGR_T_H eines") {
		t.Errorf("ohne Benutzer, keine gesetzt: %v, erwartet PGR_T_H", err)
	}
	if got, err := cli.SetzeEin(ohneBenutzer, func(n string) string { return alle[n] }); err != nil || !reflect.DeepEqual(got, []string{"", "", "h", "5432", "db/%41"}) {
		t.Errorf("ohne Benutzer: %q, %v", got, err)
	}
}
