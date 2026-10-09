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

// recordMit liest record mit --listen und --output aus der Datei inhalt, ohne
// --upstream auf der Kommandozeile.
func recordMit(t *testing.T, inhalt string, args ...string) (cli.Command, error) {
	t.Helper()
	return lese(append([]string{"record", "--listen=x", "--output=r.yaml", "--config=" + schreibe(t, inhalt)}, args...)...)
}

// Abdeckung: LH-FA-17/Happy, LH-FA-17/Boundary — in jedem Wert der Datei
// außerhalb einer URL ist $$ ein $, von links gelesen, ein $ vor einem anderen
// Zeichen bleibt stehen, und die Wertemenge prüft den Text nach $$; config show
// zeigt $$ wie geschrieben; Kommandozeile und Umgebung kennen weder $$ noch
// Platzhalter (LH-FA-17.a *Geheimnisse*).
func TestDateiDollar(t *testing.T) {
	leere(t, "record")
	leere(t, "replay")
	for geschrieben, pfad := range map[string]string{
		`"a$$b"`:   "a$b",
		`"$${X}"`:  "${X}",
		`"a$b"`:    "a$b",
		`"$$$"`:    "$$",
		`"$$$$"`:   "$$",
		`'x$$$$y'`: "x$$y",
	} {
		cmd, err := lese("record", "--listen=x", "--upstream=h:1", "--config="+schreibe(t, "record:\n  output: "+geschrieben+"\n"))
		if err != nil || cmd.Record.Output != pfad {
			t.Errorf("output: %s: %q, %v, erwartet %q", geschrieben, cmd.Record.Output, err, pfad)
		}
	}
	zeigen := schreibe(t, "record:\n  output: \"a$$b\"\n  upstream: \"h$$:1\"\n")
	if got, err := konfigurationZeigen(t, "--config", zeigen); err != nil || got != zeigen+"\nrecord:\n  output: \"a$$b\"\n  upstream: \"h$$:1\"\n" {
		t.Errorf("config show zeigt $$ wie geschrieben: %q, %v", got, err)
	}
	cmd, err := lese("record", "--listen=x", "--upstream=h:1", "--output=a$${X}")
	if err != nil || cmd.Record.Output != "a$${X}" {
		t.Errorf("--output=a$${X}: %q, %v", cmd.Record.Output, err)
	}
	t.Setenv("PGWIRE_RECORDER_OUTPUT", "${X}$$")
	cmd, err = lese("record", "--listen=x", "--upstream=h:1")
	if err != nil || cmd.Record.Output != "${X}$$" {
		t.Errorf("PGWIRE_RECORDER_OUTPUT=${X}$$: %q, %v", cmd.Record.Output, err)
	}
}

// Abdeckung: LH-FA-17/Negative — außerhalb einer URL ist jedes ${, das nicht
// aus $$ hervorgeht, PGR-E2004 an seinem Schlüssel, gleich ob seine Form gültig
// ist, auch in einem Abschnitt eines anderen Kommandos und bei config show; die
// Meldung nennt den Wert nicht (LH-FA-17.a *Geheimnisse*).
func TestDateiPlatzhalterAusserhalb(t *testing.T) {
	leere(t, "replay")
	for _, inhalt := range []string{
		"record:\n  output: \"${GEHEIM}\"\n",
		"record:\n  output: \"a$$${GEHEIM}\"\n",
		"record:\n  output: \"x${1GEHEIM\"\n",
		"record:\n  upstream: \"${GEHEIM}\"\n",
		"replay:\n  listen: \"${GEHEIM}:1\"\n",
		"log_level: \"${GEHEIM}\"\n",
	} {
		pfad := schreibe(t, inhalt)
		_, err := replayMit("--config=" + pfad)
		if !istDatei(err) || !strings.Contains(err.Error(), ": Platzhalter außerhalb einer URL") || strings.Contains(err.Error(), "GEHEIM") {
			t.Errorf("%q: %v", inhalt, err)
		}
		if got, err := konfigurationZeigen(t, "--config", pfad); !istDatei(err) || got != "" {
			t.Errorf("config show %q: %q, %v", inhalt, got, err)
		}
	}
	if _, err := replayMit("--config=" + schreibe(t, "log_level: \"in$$fo\"\n")); !istDatei(err) || !strings.Contains(err.Error(), "log_level") {
		t.Errorf("Wertemenge nach $$: %v", err)
	}
}

// Abdeckung: LH-FA-17/Happy, LH-FA-17/Boundary — der Schlüssel upstream ist
// der Name einer gültigen Verbindung der Datei, auch einer, die nach ihm steht,
// oder hat die Form host:port: Host nicht leer, IPv6 in eckigen Klammern, Port
// wie geschrieben mit führenden Nullen (LH-FA-17.a).
func TestDateiUpstream(t *testing.T) {
	leere(t, "record")
	for inhalt, want := range map[string]string{
		"record:\n  upstream: staging\nconnections:\n  staging: postgresql://h/db\n": "staging",
		"connections:\n  staging: postgresql://h/db\nrecord:\n  upstream: staging\n": "staging",
		"record:\n  upstream: h:5432\n":                                              "h:5432",
		"record:\n  upstream: \"[::1]:5432\"\n":                                      "[::1]:5432",
		"record:\n  upstream: h:05432\n":                                             "h:05432",
		"record:\n  upstream: \"h%41:1\"\n":                                          "h%41:1",
		"record:\n  upstream: \"h$$:1\"\n":                                           "h$:1",
	} {
		cmd, err := recordMit(t, inhalt)
		if err != nil || cmd.Record.Upstream != want {
			t.Errorf("%q: %q, %v, erwartet %q", inhalt, cmd.Record.Upstream, err, want)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — ein Schlüssel upstream, der weder eine gültige
// Verbindung der Datei nennt noch die Form host:port hat, ist PGR-E2004 an
// record.upstream ohne den Wert, auch bei replay und config show: ein
// unbekannter Name, einer in anderer Schreibweise, einer mit $$, der nach $$
// keinen Namen trifft, leerer Host, kein Port, ein Port außerhalb 1 bis 65535
// oder dekodiert geschrieben, IPv6 ohne oder mit offener Klammer; ein
// ungültiger Name zählt nicht als Name, auch einer mit $, den ein Wert nach $$
// wörtlich trifft; der erste Fehler ist der, der in der Datei zuerst steht
// (LH-FA-17.a, L1, L4).
func TestDateiUpstreamUngueltig(t *testing.T) {
	leere(t, "record")
	leere(t, "replay")
	verbindungen := "connections:\n  staging: postgresql://h/db\n  ab: postgresql://h/db\n"
	for _, wert := range []string{
		"GEHEIM", "Staging", "\"a$$b\"", "\":5432\"", "\"GEHEIM:\"", "GEHEIM.example", "GEHEIM:0", "GEHEIM:65536",
		"GEHEIM:%35", "\"::1:5432\"", "\"[::1]\"", "\"[::1:5432\"", "\"[::1]x:1\"", "\"[]:1\"", "\"a]b:1\"", "GEHEIM:1:2",
	} {
		inhalt := verbindungen + "record:\n  upstream: " + wert + "\n"
		genau := "Konfiguration [PGR-E2004]: Konfigurationsdatei: record.upstream: weder Name einer Verbindung der Datei noch host:port"
		for name, lauf := range map[string]func() error{
			"record": func() error { _, err := recordMit(t, inhalt); return err },
			"replay": func() error { _, err := replayMit("--config=" + schreibe(t, inhalt)); return err },
			"config show": func() error {
				got, err := konfigurationZeigen(t, "--config", schreibe(t, inhalt))
				if got != "" {
					t.Errorf("config show zeigt %q", got)
				}
				return err
			},
		} {
			if err := lauf(); err == nil || err.Error() != genau || !istDatei(err) {
				t.Errorf("%s, upstream: %s: %v, erwartet %q", name, wert, err, genau)
			}
		}
	}
	for inhalt, stelle := range map[string]string{
		"record:\n  upstream: a@b\nconnections:\n  a@b: postgresql://h/db\n":             "Konfigurationsdatei: record.upstream: weder Name",
		"connections:\n  a@b: postgresql://h/db\nrecord:\n  upstream: a@b\n":             "Konfigurationsdatei: connections: Name einer Verbindung",
		"record:\n  upstream: \"a$$$$b\"\nconnections:\n  \"a$$b\": postgresql://h/db\n": "Konfigurationsdatei: record.upstream: weder Name",
		"record:\n  upstream: v\nconnections:\n  v: postgres://h/db\n":                   "Konfigurationsdatei: connections.v: Schema",
	} {
		if _, err := recordMit(t, inhalt); !istDatei(err) || !strings.Contains(err.Error(), stelle) {
			t.Errorf("%q: %v, erwartet %q", inhalt, err, stelle)
		}
	}
}

// Abdeckung: LH-FA-17/Negative — ein Klartext-Passwort ist PGR-E2006 an
// connections.<Name> ohne den Wert, in jeder Verbindung der Datei, auch bei
// config show, das dann nichts zeigt: ein Passwortteil, dessen geschriebener
// Text nicht genau ein ${VAR} ist (leer, wörtlich, $${VAR}, ein fehlerhafter
// Platzhalter, ein Platzhalter mit Text, zwei Platzhalter, ein ungültiges
// Escape), getrennt am ersten :; ein Parameter password, auch dekodiert
// geschrieben, mit Platzhalter, leerem oder ungültig kodiertem Wert. Geprüft
// in der Reihenfolge der URL: nach Benutzer, vor Host, Port, Datenbank, den
// späteren Parametern und dem Fragment; ein Parameter davor geht vor
// (LH-FA-17.a *Geheimnisse*).
func TestVerbindungKlartext(t *testing.T) {
	leere(t, "replay")
	passwort := "Konfiguration [PGR-E2006]: Konfigurationsdatei: connections.v: Klartext-Passwort, erlaubt ist im Passwort nur genau ein Platzhalter ${VAR}"
	parameter := "Konfiguration [PGR-E2006]: Konfigurationsdatei: connections.v: Parameter password ist ein Klartext-Passwort"
	for url, genau := range map[string]string{
		"postgresql://u:GEHEIM@h/db":                              passwort,
		"postgresql://u:@h/db":                                    passwort,
		"postgresql://u:$${GEHEIM}@h/db":                          passwort,
		"postgresql://u:${1GEHEIM}@h/db":                          passwort,
		"postgresql://u:${GEHEIM@h/db":                            passwort,
		"postgresql://u:x${GEHEIM}@h/db":                          passwort,
		"postgresql://u:${P}${GEHEIM}@h/db":                       passwort,
		"postgresql://u:%ZZGEHEIM@h/db":                           passwort,
		"postgresql://u:%47EHEIM@h/db":                            passwort,
		"postgresql://${U}:${P}:${GEHEIM}@h/db":                   passwort,
		"postgresql://u:GEHEIM@[x/db?a=b":                         passwort,
		"postgresql://u:GEHEIM@h:0/db":                            passwort,
		"postgresql://u:GEHEIM@h/db?foo=1#x":                      passwort,
		"postgresql://h/db?password=GEHEIM":                       parameter,
		"postgresql://h/db?pass%77ord=GEHEIM":                     parameter,
		"postgresql://h/db?password=${GEHEIM}":                    parameter,
		"postgresql://h/db?password=":                             parameter,
		"postgresql://h/db?password=%ZZGEHEIM":                    parameter,
		"postgresql://h/db?sslmode=disable&password=GEHEIM&foo=1": parameter,
		"postgresql://h/db?password=GEHEIM#x":                     parameter,
	} {
		pfad := schreibe(t, mitVerbindung(url))
		_, err := replayMit("--config=" + pfad)
		if err == nil || err.Error() != genau || !hatCode(err, model.CodeConfigPassword) {
			t.Errorf("%s: %v, erwartet %q", url, err, genau)
		}
		if got, err := konfigurationZeigen(t, "--config", pfad); !hatCode(err, model.CodeConfigPassword) || got != "" {
			t.Errorf("config show %s: %q, %v", url, got, err)
		}
	}
	for url, grund := range map[string]string{
		"postgresql://:GEHEIM@h/db":               "connections.v: Benutzer ist leer",
		"postgresql://u:GEHEIM\\x01@h/db":         "connections.v: Steuerzeichen in der URL",
		"postgres://u:GEHEIM@h/db":                "connections.v: Schema",
		"postgresql://h/db?foo=1&password=GEHEIM": "connections.v: unbekannter Parameter",
		"postgresql://h/db?PASSWORD=GEHEIM":       "connections.v: unbekannter Parameter",
		"postgresql://h/db?pass%ZZword=GEHEIM":    "connections.v: Name eines Parameters: ungültiges Escape",
		"postgresql://u:GEHEIM/db@h/db":           "connections.v: Port ist keine Zahl",
	} {
		_, err := replayMit("--config=" + schreibe(t, mitVerbindung(url)))
		if !istDatei(err) || !strings.Contains(err.Error(), grund) || strings.Contains(err.Error(), "GEHEIM") {
			t.Errorf("%s: %v, erwartet %s mit %q", url, err, model.CodeConfigFile, grund)
		}
	}
	for inhalt, code := range map[string]string{
		"connections:\n  a: postgresql://h/db\n  v: postgresql://u:GEHEIM@h/db\nbogus: 1\n": model.CodeConfigPassword,
		"bogus: 1\nconnections:\n  v: postgresql://u:GEHEIM@h/db\n":                         model.CodeConfigFile,
		"connections:\n  v: postgresql://u:GEHEIM@h/db\n  w: postgres://h/db\n":             model.CodeConfigPassword,
	} {
		if _, err := replayMit("--config=" + schreibe(t, inhalt)); !hatCode(err, code) || strings.Contains(err.Error(), "GEHEIM") {
			t.Errorf("%q: %v, erwartet %s", inhalt, err, code)
		}
	}
	for _, url := range []string{"postgresql://u:${PGWIRE_RECORDER_PASSWORD}@h/db", "postgresql://u@h/db", "postgresql://h/db?sslmode=require"} {
		if _, err := replayMit("--config=" + schreibe(t, mitVerbindung(url))); err != nil {
			t.Errorf("%s: %v", url, err)
		}
	}
}
