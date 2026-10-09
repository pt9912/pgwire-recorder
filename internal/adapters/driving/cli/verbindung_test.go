package cli_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// mitVerbindung ist der Text einer Datei mit der einen Verbindung v, deren URL
// url in doppelten Anführungszeichen steht (YAML-Escapes wie \x01 gelten).
func mitVerbindung(url string) string {
	return "connections:\n  v: \"" + url + "\"\n"
}

// Abdeckung: LH-FA-17/Happy, LH-FA-17/Boundary — eine URL wird von links
// zerlegt: das letzte @ trennt den Benutzerteil ab, darin das erste : das
// Passwort; ein Host in eckigen Klammern ist IPv6 ohne die Klammern; ohne Port
// gilt 5432, ein Port gilt wie geschrieben mit führenden Nullen, 1 und 65535
// sind gültig; die Datenbank reicht bis zum ?, auch mit /; wörtliche Teile und
// Name und Wert eines Parameters werden prozent-dekodiert; $$ ist ein $;
// Platzhalter stehen je Teil in der Reihenfolge der URL, im Port mit
// wörtlichen Ziffern; sslmode ist disable ohne Parameter, sonst disable oder
// require (LH-FA-17.a *Benannte Verbindungen*).
func TestVerbindungZerlegung(t *testing.T) {
	leere(t, "replay")
	for _, f := range []struct {
		url  string
		want cli.Zerlegt
	}{
		{"postgresql://dev@localhost:5432/myapp", cli.Zerlegt{Benutzer: []string{"dev"}, Host: []string{"localhost"}, Port: []string{"5432"}, Datenbank: []string{"myapp"}, SSLMode: "disable"}},
		{"postgresql://app:${STAGING_PASSWORD}@staging.example.com:5432/myapp?sslmode=require", cli.Zerlegt{Benutzer: []string{"app"}, Passwort: []string{"<STAGING_PASSWORD>"}, Host: []string{"staging.example.com"}, Port: []string{"5432"}, Datenbank: []string{"myapp"}, SSLMode: "require"}},
		{"postgresql://h/db", cli.Zerlegt{Host: []string{"h"}, Port: []string{"5432"}, Datenbank: []string{"db"}, SSLMode: "disable"}},
		{"postgresql://[::1]:6432/db?sslmode=disable", cli.Zerlegt{Host: []string{"::1"}, Port: []string{"6432"}, Datenbank: []string{"db"}, SSLMode: "disable"}},
		{"postgresql://[fe80::1%25eth0]/db", cli.Zerlegt{Host: []string{"fe80::1%eth0"}, Port: []string{"5432"}, Datenbank: []string{"db"}, SSLMode: "disable"}},
		{"postgresql://h:05432/db", cli.Zerlegt{Host: []string{"h"}, Port: []string{"05432"}, Datenbank: []string{"db"}, SSLMode: "disable"}},
		{"postgresql://h:1/db", cli.Zerlegt{Host: []string{"h"}, Port: []string{"1"}, Datenbank: []string{"db"}, SSLMode: "disable"}},
		{"postgresql://h:65535/db", cli.Zerlegt{Host: []string{"h"}, Port: []string{"65535"}, Datenbank: []string{"db"}, SSLMode: "disable"}},
		{"postgresql://a@b@h/a/b", cli.Zerlegt{Benutzer: []string{"a@b"}, Host: []string{"h"}, Port: []string{"5432"}, Datenbank: []string{"a/b"}, SSLMode: "disable"}},
		{"postgresql://%75ser@h%2Dx/d%C3%A4?ssl%6Dode=%72equire", cli.Zerlegt{Benutzer: []string{"user"}, Host: []string{"h-x"}, Port: []string{"5432"}, Datenbank: []string{"dä"}, SSLMode: "require"}},
		{"postgresql://${U}:${P}@${H}:${PORT}/${DB}", cli.Zerlegt{Benutzer: []string{"<U>"}, Passwort: []string{"<P>"}, Host: []string{"<H>"}, Port: []string{"<PORT>"}, Datenbank: []string{"<DB>"}, SSLMode: "disable"}},
		{"postgresql://x${U}y@db-${N}.example:5${P}/a${D}", cli.Zerlegt{Benutzer: []string{"x", "<U>", "y"}, Host: []string{"db-", "<N>", ".example"}, Port: []string{"5", "<P>"}, Datenbank: []string{"a", "<D>"}, SSLMode: "disable"}},
		{"postgresql://a$$b@h/d$${X}", cli.Zerlegt{Benutzer: []string{"a$b"}, Host: []string{"h"}, Port: []string{"5432"}, Datenbank: []string{"d${X}"}, SSLMode: "disable"}},
		{"postgresql://a$$$${X}@h/$x", cli.Zerlegt{Benutzer: []string{"a$${X}"}, Host: []string{"h"}, Port: []string{"5432"}, Datenbank: []string{"$x"}, SSLMode: "disable"}},
		{"postgresql://$$${U}@h/db", cli.Zerlegt{Benutzer: []string{"$", "<U>"}, Host: []string{"h"}, Port: []string{"5432"}, Datenbank: []string{"db"}, SSLMode: "disable"}},
	} {
		got, err := cli.Zerlege(schreibe(t, mitVerbindung(f.url)))
		if err != nil || !reflect.DeepEqual(got, f.want) {
			t.Errorf("%s: %#v, %v, erwartet %#v", f.url, got, err, f.want)
		}
	}
	for _, name := range []string{"staging", "a.b-c", "mit Leerraum"} {
		if _, err := replayMit("--config=" + schreibe(t, "connections:\n  "+name+": postgresql://h/db\n")); err != nil {
			t.Errorf("Name %q: %v", name, err)
		}
	}
}

// ungueltigeURLs sind URLs, die PGR-E2004 an connections.v sind, mit dem Grund,
// den die Meldung nennt; GEHEIM steht für einen Wert, den keine Meldung nennt.
func ungueltigeURLs() []struct{ url, grund string } {
	return []struct{ url, grund string }{
		{`postgresql://h\x01GEHEIM/db`, "Steuerzeichen in der URL"},
		{`postgresql://h/d\tGEHEIM`, "Steuerzeichen in der URL"},
		{`postgresql://u:p\x9fGEHEIM@h/db`, "Steuerzeichen in der URL"},
		{`postgresql://h/d\x7fGEHEIM`, "Steuerzeichen in der URL"},
		{"postgres://u@GEHEIM/db", "Schema ist nicht postgresql://"},
		{"POSTGRESQL://u@GEHEIM/db", "Schema ist nicht postgresql://"},
		{"postgres://:p@/", "Schema ist nicht postgresql://"},
		{"postgresql://@GEHEIM/db", "Benutzer ist leer"},
		{"postgresql://:${PW}@GEHEIM/db", "Benutzer ist leer"},
		{"postgresql://@/db", "Benutzer ist leer"},
		{"postgresql://${1}@GEHEIM/db", "Benutzer: ungültiger Platzhalter"},
		{"postgresql://%ZZ@GEHEIM/db", "Benutzer: ungültiges Escape"},
		{"postgresql://%C2%80@GEHEIM/db", "Benutzer: Escape ergibt ein Steuerzeichen"},
		{"postgresql:///db", "Host ist leer"},
		{"postgresql://u@/db", "Host ist leer"},
		{"postgresql://:5432/db", "Host ist leer"},
		{"postgresql://::1/db", "Host ist leer"},
		{"postgresql://[]:5432/db", "Host ist leer"},
		{"postgresql://[abc/db", "Host mit [ ohne ]"},
		{"postgresql://[::1/db", "Host mit [ ohne ]"},
		{"postgresql://a]b/db", "Host mit [ oder ] an falscher Stelle"},
		{"postgresql://a[b/db", "Host mit [ oder ] an falscher Stelle"},
		{"postgresql://[a[b]/db", "Host mit [ oder ] an falscher Stelle"},
		{"postgresql://[::1]x/db", "hinter ] steht nicht : mit Port"},
		{"postgresql://%09/db", "Host: Escape ergibt ein Steuerzeichen"},
		{"postgresql://h${/db", "Host: ungültiger Platzhalter"},
		{"postgresql://h:/db", ": ohne Port"},
		{"postgresql://h:0/db", "Port ist keine Zahl von 1 bis 65535"},
		{"postgresql://h:00000/db", "Port ist keine Zahl von 1 bis 65535"},
		{"postgresql://h:65536/db", "Port ist keine Zahl von 1 bis 65535"},
		{"postgresql://h:99999999999999999999/db", "Port ist keine Zahl von 1 bis 65535"},
		{"postgresql://h:%35/db", "Port ist keine Zahl von 1 bis 65535"},
		{"postgresql://h:+5/db", "Port ist keine Zahl von 1 bis 65535"},
		{"postgresql://h:1:2/db", "Port ist keine Zahl von 1 bis 65535"},
		{"postgresql://h1:5432,h2:5432/db", "Port ist keine Zahl von 1 bis 65535"},
		{"postgresql://h:x${P}/db", "Port mit Zeichen außer Ziffern"},
		{"postgresql://h:$$${P}/db", "Port mit Zeichen außer Ziffern"},
		{"postgresql://h:${P/db", "Port: ungültiger Platzhalter"},
		{"postgresql://GEHEIM", "Datenbank fehlt"},
		{"postgresql://GEHEIM?sslmode=disable", "Datenbank fehlt"},
		{"postgresql://GEHEIM#x", "Datenbank fehlt"},
		{"postgresql://h/", "Datenbank ist leer"},
		{"postgresql://h/?sslmode=disable", "Datenbank ist leer"},
		{"postgresql://h/%ZZ", "Datenbank: ungültiges Escape"},
		{"postgresql://h/%4", "Datenbank: ungültiges Escape"},
		{"postgresql://h/%00", "Datenbank: Escape ergibt ein Steuerzeichen"},
		{"postgresql://h/%FF", "Datenbank: Escape ergibt kein gültiges UTF-8"},
		{"postgresql://h/${A-B}", "Datenbank: ungültiger Platzhalter"},
		{"postgresql://h/db?", "Parameter ohne ="},
		{"postgresql://h/db?sslmode", "Parameter ohne ="},
		{"postgresql://h/db?sslmode=disable&", "Parameter ohne ="},
		{"postgresql://h/db?foo=GEHEIM", "unbekannter Parameter"},
		{"postgresql://h/db?SSLMODE=disable", "unbekannter Parameter"},
		{"postgresql://h/db?PASSWORD=GEHEIM", "unbekannter Parameter"},
		{"postgresql://h/db?${N}=disable", "unbekannter Parameter"},
		{"postgresql://h/db?pass%ZZ=GEHEIM", "Name eines Parameters: ungültiges Escape"},
		{"postgresql://h/db?sslmode=prefer", "ungültiger sslmode, erlaubt sind disable und require"},
		{"postgresql://h/db?sslmode=", "ungültiger sslmode, erlaubt sind disable und require"},
		{"postgresql://h/db?sslmode=Require", "ungültiger sslmode, erlaubt sind disable und require"},
		{"postgresql://h/db?sslmode=${M}", "ungültiger sslmode, erlaubt sind disable und require"},
		{"postgresql://h/db?sslmode=%ZZ", "sslmode: ungültiges Escape"},
		{"postgresql://h/db?sslmode=disable&sslmode=disable", "Parameter sslmode doppelt"},
		{"postgresql://h/db#GEHEIM", "Fragment ist ungültig"},
		{"postgresql://h/db?sslmode=disable#GEHEIM", "Fragment ist ungültig"},
	}
}

// Abdeckung: LH-FA-17/Negative — was die Grammatik der URL nicht zulässt, ist
// PGR-E2004 an connections.<Name> mit genauem Grund und ohne den Wert: ein
// Steuerzeichen im geschriebenen Text (auch Tabulator, DEL und C1), ein
// anderes Schema, ein leerer Benutzer, ein leerer Host, [ ohne ], [ oder ] an
// anderer Stelle des Hosts, Text hinter ], ein : ohne Port, ein Port außerhalb
// 1 bis 65535 oder nicht aus Ziffern (auch dekodiert geschrieben), ein Port mit
// Platzhalter und wörtlichem Zeichen außer Ziffern, eine fehlende oder leere
// Datenbank, ein Parameter ohne =, ein unbekannter (auch in anderer
// Schreibweise), ein ungültiger sslmode (auch als Platzhalter), sslmode
// zweimal, ein Fragment, ein ungültiges Escape, eines, das ein Steuerzeichen
// oder kein gültiges UTF-8 ergibt, und ein ungültiger Platzhalter in jedem Teil;
// geprüft in der Reihenfolge Steuerzeichen, Schema, Benutzer, Host, Port,
// Datenbank, Parameter, Fragment; auch bei config show, das dann nichts zeigt
// (LH-FA-17.a *Benannte Verbindungen*).
func TestVerbindungUngueltig(t *testing.T) {
	leere(t, "replay")
	for _, f := range ungueltigeURLs() {
		pfad := schreibe(t, mitVerbindung(f.url))
		genau := "Konfiguration [PGR-E2004]: Konfigurationsdatei: connections.v: " + f.grund
		if _, err := replayMit("--config=" + pfad); err == nil || err.Error() != genau || !istDatei(err) {
			t.Errorf("%s: %v, erwartet %q", f.url, err, genau)
		}
		if got, err := konfigurationZeigen(t, "--config", pfad); !istDatei(err) || got != "" {
			t.Errorf("config show %s: %q, %v", f.url, got, err)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — ein Name einer Verbindung mit :, @ oder $ ist
// PGR-E2004, die Meldung nennt als Stelle nur connections, nicht den Namen und
// keinen Teil davon; der Name gilt wörtlich, ohne $$; die Verbindungen werden
// in der Reihenfolge der Datei geprüft, auch nach einer gültigen und vor dem
// Abschnitt eines Kommandos (LH-FA-17.a).
func TestVerbindungName(t *testing.T) {
	leere(t, "replay")
	for _, name := range []string{"postgresql://u:GEHEIM@h/db", "a:b", "a@b", "a$b", "a$$b", "${X}"} {
		pfad := schreibe(t, "connections:\n  \""+name+"\": postgresql://h/db\n")
		_, err := replayMit("--config=" + pfad)
		genau := "Konfiguration [PGR-E2004]: Konfigurationsdatei: connections: Name einer Verbindung mit :, @ oder $"
		if err == nil || err.Error() != genau || strings.Contains(err.Error(), "GEHEIM") || strings.Contains(err.Error(), name) {
			t.Errorf("Name %q: %v, erwartet %q", name, err, genau)
		}
		if got, err := konfigurationZeigen(t, "--config", pfad); !istDatei(err) || got != "" {
			t.Errorf("config show mit Name %q: %q, %v", name, got, err)
		}
	}
	for inhalt, stelle := range map[string]string{
		"connections:\n  a: postgresql://h/db\n  b: postgres://h/db\n  c: postgresql://h\n": "connections.b: Schema",
		"connections:\n  a: postgresql://h\nreplay:\n  listen: \"\"\n":                      "connections.a: Datenbank fehlt",
		"replay:\n  listen: \"\"\nconnections:\n  a: postgresql://h\n":                      "replay.listen",
	} {
		if _, err := replayMit("--config=" + schreibe(t, inhalt)); !istDatei(err) || !strings.Contains(err.Error(), stelle) {
			t.Errorf("%q: erwartet %s an %s, erhalten %v", inhalt, model.CodeConfigFile, stelle, err)
		}
	}
}
