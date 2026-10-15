package recording_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// nm ist eine Zeile von message ohne ihre Nummern (session_id, sequence, position stammen aus
// der Stellung im Aufruf); gruppe 0 ist NULL.
type nm struct {
	richtung, art, felder string
	gruppe                int
}

// ia ist eine Interaktion der Datei: erw für type extended, sonst query mit sql 'x'; offset ist
// der Wert von offset_ms ("" für NULL).
type ia struct {
	erw    bool
	offset string
	msgs   []nm
}

const (
	fatalFelder = `{"notice":{"V":"FATAL","S":"FATAL","C":"57P01","M":"x"}}`
	fehlerFeld  = `{"notice":{"V":"ERROR","S":"ERROR","C":"42P01","M":"x"}}`
)

func srv(art, felder string) nm         { return nm{"server", art, felder, 0} }
func srvG(g int, art, felder string) nm { return nm{"server", art, felder, g} }
func cli(g int, art, felder string) nm  { return nm{"client", art, felder, g} }
func einf(msgs ...nm) ia                { return ia{msgs: msgs} }
func erw(msgs ...nm) ia                 { return ia{erw: true, msgs: msgs} }
func fatal() nm                         { return srv("error_response", fatalFelder) }
func ende() nm                          { return srv("connection_closed", `{}`) }
func fertig() []nm {
	return []nm{srv("command_complete", `{"tag":"x"}`), srv("ready_for_query", `{"tx_status":"I"}`)}
}
func parse(g int) nm         { return cli(g, "parse", `{"statement":"s","sql":"x","param_types":[]}`) }
func sync(g int) nm          { return cli(g, "sync", `{}`) }
func flush(g int) nm         { return cli(g, "flush", `{}`) }
func fertigeInteraktion() ia { return einf(fertig()...) }
func fertigeGruppe() []nm {
	return []nm{parse(1), sync(1), srvG(1, "ready_for_query", `{"tx_status":"I"}`)}
}
func endeInGruppe(abschluss func(int) nm) ia {
	return erw(parse(1), abschluss(1), srvG(1, "error_response", fatalFelder), srvG(1, "connection_closed", `{}`))
}

// endeSQL liefert die Anweisungen, die den Inhalt der Datei durch die Sessions ersetzen (je
// Session mit der Nummer ihrer Stellung) und meta.version auf version setzen. rueck schreibt
// die Zeilen von message in umgekehrter Reihenfolge.
func endeSQL(version string, rueck bool, sessions ...[]ia) string {
	var b strings.Builder
	b.WriteString(`DELETE FROM message; DELETE FROM interaction; DELETE FROM session; UPDATE meta SET value = '` + version + `' WHERE key = 'version';`)
	for si, s := range sessions {
		fmt.Fprintf(&b, `INSERT INTO session VALUES (%d, '{}', '{}');`, si+1)
		for ii, i := range s {
			art, sqlText, offset := "query", "'x'", "NULL"
			if i.erw {
				art, sqlText = "extended", "NULL"
			}
			if i.offset != "" {
				offset = i.offset
			}
			fmt.Fprintf(&b, `INSERT INTO interaction VALUES (%d, %d, '%s', %s, %s);`, si+1, ii+1, art, offset, sqlText)
			var zeilen []string
			for mi, m := range i.msgs {
				gruppe := "NULL"
				if m.gruppe > 0 {
					gruppe = fmt.Sprint(m.gruppe)
				}
				zeilen = append(zeilen, fmt.Sprintf(`INSERT INTO message VALUES (%d, %d, %d, '%s', %s, '%s', '%s');`, si+1, ii+1, mi+1, m.richtung, gruppe, m.art, m.felder))
			}
			if rueck {
				for k := len(zeilen) - 1; k >= 0; k-- {
					b.WriteString(zeilen[k])
				}
			} else {
				b.WriteString(strings.Join(zeilen, ""))
			}
		}
	}
	return b.String()
}

// endeFall ist eine Datei, die der Test auf einer gültigen Datei mit SQL aufbaut, und was der
// Leser aus ihr macht: version ist der Text von meta.version, nachher folgt als SQL, enden die
// Zahl der Interaktionen mit Ende in der geladenen Aufzeichnung und want der Meldungscode.
type endeFall struct {
	name, version, nachher string
	rueck                  bool
	sessions               [][]ia
	enden                  int
	want                   string
}

// einSession ist die Liste der Sessions eines Falls mit einer Session.
func einSession(is ...ia) [][]ia { return [][]ia{is} }

func pruefeEndeFaelle(t *testing.T, faelle []endeFall) {
	t.Helper()
	basis := sqliteDatei(t, sessionsMit(1))
	for _, c := range faelle {
		path := kopiere(t, basis)
		rohSchreibe(t, path, endeSQL(c.version, c.rueck, c.sessions...)+c.nachher)
		rec, err := laden(path)
		if got := codeVon(err); got != c.want {
			t.Errorf("%s: Meldungscode %q, erwartet %q (%v)", c.name, got, c.want, err)
			continue
		}
		if c.want == unbekannteV && strings.Contains(err.Error(), "Version "+c.version) {
			t.Errorf("%s: die Meldung nennt die Version: %v", c.name, err)
		}
		if c.want != "" {
			continue
		}
		hat := 0
		for _, s := range rec.Sessions {
			for _, i := range s.Interactions {
				if endeIn(i) {
					hat++
				}
			}
		}
		if hat != c.enden || fmt.Sprint(rec.Version) != c.version {
			t.Errorf("%s: %d Enden, Version %d; erwartet %d Enden, Version %s", c.name, hat, rec.Version, c.enden, c.version)
		}
	}
}

// Abdeckung: LH-FA-25/Happy — der Leser von sqlite lädt das Ende des Servers aus einer Zeile von message
// (kind connection_closed, direction server, fields {}) als letzte Antwort der letzten Interaktion
// einer Session: in einer einfachen Interaktion nach FATAL, PANIC und nach S FATAL ohne V, in einer
// Extended-Interaktion mit Sync- und mit Flush-Gruppe, mit offset_ms, in einer Session, die nicht die letzte
// ist, und auch dann, wenn die Zeile des Endes in der Datei vor der der Fehlerantwort steht (die
// Stelle folgt den Nummern); die Aufzeichnung trägt die Version der Datei (1 oder 2, auch 2 ohne Ende),
// PRAGMA user_version ändert sie nicht.
func TestLesenEndeAngenommen(t *testing.T) {
	panikFeld := `{"notice":{"V":"PANIC","S":"PANIC","C":"XX000","M":"x"}}`
	nurS := `{"notice":{"S":"FATAL","C":"57P01","M":"x"}}`
	pruefeEndeFaelle(t, []endeFall{
		{"einfach", "2", "", false, einSession(einf(fatal(), ende())), 1, ""},
		{"einfach, PANIC", "2", "", false, einSession(einf(srv("error_response", panikFeld), ende())), 1, ""},
		{"einfach, V fehlt, S FATAL", "2", "", false, einSession(einf(srv("error_response", nurS), ende())), 1, ""},
		{"extended, Sync-Gruppe", "2", "", false, einSession(endeInGruppe(sync)), 1, ""},
		{"extended, Flush-Gruppe", "2", "", false, einSession(endeInGruppe(flush)), 1, ""},
		{"extended, vordere Flush-Gruppe", "2", "", false, einSession(erw(parse(1), flush(1), srvG(1, "parse_complete", `{}`), parse(2), sync(2), srvG(2, "error_response", fatalFelder), srvG(2, "connection_closed", `{}`))), 1, ""},
		{"mit offset_ms", "2", "", false, einSession(ia{offset: "5", msgs: fertig()}, ia{offset: "9", msgs: []nm{fatal(), ende()}}), 1, ""},
		{"Ende in der letzten von zwei Interaktionen", "2", "", false, einSession(fertigeInteraktion(), einf(fatal(), ende())), 1, ""},
		{"Ende in einer Session, die nicht die letzte ist", "2", "", false, [][]ia{{einf(fatal(), ende())}, {fertigeInteraktion()}}, 1, ""},
		{"Zeile des Endes vor der Fehlerantwort", "2", "", true, einSession(einf(fatal(), ende())), 1, ""},
		{"Zeile des Endes vor der Fehlerantwort, extended", "2", "", true, einSession(endeInGruppe(sync)), 1, ""},
		{"Version 2 ohne Ende", "2", "", false, einSession(fertigeInteraktion()), 0, ""},
		{"Version 1 ohne Ende", "1", "", false, einSession(fertigeInteraktion()), 0, ""},
		{"Version 2 ohne Ende, extended", "2", "", false, einSession(erw(fertigeGruppe()...)), 0, ""},
		{"user_version 1, meta.version 2, Ende", "2", `PRAGMA user_version = 1;`, false, einSession(einf(fatal(), ende())), 1, ""},
		{"user_version 2, meta.version 1 ohne Ende", "1", `PRAGMA user_version = 2;`, false, einSession(fertigeInteraktion()), 0, ""},
	})
}

// Abdeckung: LH-FA-25/Negative — der Leser von sqlite lehnt mit PGR-E3003 ab: fields des Endes mit
// irgendeinem Schlüssel (auch mit Wert null, mit leerem Namen und mit dem Schlüssel type), direction
// client bei kind connection_closed (in einer einfachen und in einer Extended-Interaktion), kind in
// anderer Schreibweise (Großbuchstaben, nachgestelltes Leerzeichen), connection_closed in einer Datei
// der Version 1 (einfach und extended) und bei PRAGMA user_version 2 mit meta.version 1.
func TestLesenEndeAbgelehntDarstellung(t *testing.T) {
	mitFelder := func(felder string) [][]ia {
		return einSession(einf(fatal(), srv("connection_closed", felder)))
	}
	pruefeEndeFaelle(t, []endeFall{
		{"fields mit tag", "2", "", false, mitFelder(`{"tag":"x"}`), 0, beschaedigt},
		{"fields mit Wert null", "2", "", false, mitFelder(`{"x":null}`), 0, beschaedigt},
		{"fields mit leerem Namen", "2", "", false, mitFelder(`{"":1}`), 0, beschaedigt},
		{"fields mit type", "2", "", false, mitFelder(`{"type":"x"}`), 0, beschaedigt},
		{"direction client, einfach", "2", "", false, einSession(einf(fatal(), nm{"client", "connection_closed", `{}`, 0})), 0, beschaedigt},
		{"direction client, extended", "2", "", false, einSession(erw(parse(1), cli(1, "connection_closed", `{}`), sync(1), srvG(1, "ready_for_query", `{"tx_status":"I"}`))), 0, beschaedigt},
		{"kind in Großbuchstaben", "2", "", false, einSession(einf(fatal(), srv("Connection_Closed", `{}`))), 0, beschaedigt},
		{"kind mit Leerzeichen", "2", "", false, einSession(einf(fatal(), srv("connection_closed ", `{}`))), 0, beschaedigt},
		{"Version 1, einfach", "1", "", false, einSession(einf(fatal(), ende())), 0, beschaedigt},
		{"Version 1, extended", "1", "", false, einSession(endeInGruppe(sync)), 0, beschaedigt},
		{"user_version 2, meta.version 1, Ende", "1", `PRAGMA user_version = 2;`, false, einSession(einf(fatal(), ende())), 0, beschaedigt},
	})
}

// Abdeckung: LH-FA-25/Negative — der Leser von sqlite lehnt mit PGR-E3003 ab: connection_closed an
// jeder falschen Stelle, je Stelle eine Datei von Hand gebaut (in der ersten Nachricht, mitten in der
// Interaktion, vor der Fehlerantwort, nach einer Meldung zwischen Fehlerantwort und Ende, nach ERROR,
// nach V leer mit S FATAL, nach Fatal in Kleinschreibung, mit ready_for_query vor und nach dem Ende,
// zweimal in derselben Interaktion, in der Interaktion 1 von 2, in zwei Interaktionen einer Session,
// bei Extended nicht in der letzten Gruppe, nicht als letzte Server-Nachricht der Gruppe, mit
// ready_for_query in der Gruppe, hinter einer Sync-Gruppe und in einer letzten Gruppe ohne Sync und
// Flush).
func TestLesenEndeAbgelehntStelle(t *testing.T) {
	meldung := srv("notice_response", `{}`)
	rfq := srv("ready_for_query", `{"tx_status":"I"}`)
	f := func(name string, s ...ia) endeFall {
		return endeFall{name, "2", "", false, einSession(s...), 0, beschaedigt}
	}
	pruefeEndeFaelle(t, []endeFall{
		f("Anfang", einf(ende(), fatal())),
		f("Mitte", einf(fatal(), ende(), meldung)),
		f("vor der Fehlerantwort", einf(srv("command_complete", `{"tag":"x"}`), ende(), fatal())),
		f("Meldung zwischen Fehlerantwort und Ende", einf(fatal(), meldung, ende())),
		f("nach ERROR", einf(srv("error_response", fehlerFeld), ende())),
		f("V leer, S FATAL", einf(srv("error_response", `{"notice":{"V":"","S":"FATAL"}}`), ende())),
		f("Fatal in Kleinschreibung", einf(srv("error_response", `{"notice":{"V":"Fatal"}}`), ende())),
		f("ready_for_query vor dem Ende", einf(fatal(), rfq, ende())),
		f("ready_for_query nach dem Ende", einf(fatal(), ende(), rfq)),
		f("zweites Ende in derselben Interaktion", einf(fatal(), ende(), ende())),
		f("Ende in der Interaktion 1 von 2", einf(fatal(), ende()), fertigeInteraktion()),
		f("Ende in zwei Interaktionen", einf(fatal(), ende()), einf(fatal(), ende())),
		f("extended, nicht in der letzten Gruppe", erw(parse(1), flush(1), srvG(1, "error_response", fatalFelder), srvG(1, "connection_closed", `{}`), parse(2), sync(2), srvG(2, "ready_for_query", `{"tx_status":"I"}`))),
		f("extended, nicht die letzte Server-Nachricht der Gruppe", erw(parse(1), sync(1), srvG(1, "error_response", fatalFelder), srvG(1, "connection_closed", `{}`), srvG(1, "notice_response", `{}`))),
		f("extended, ready_for_query in der Gruppe", erw(parse(1), sync(1), srvG(1, "error_response", fatalFelder), srvG(1, "ready_for_query", `{"tx_status":"I"}`), srvG(1, "connection_closed", `{}`))),
		f("extended, vorletzte Gruppe mit Sync", erw(parse(1), sync(1), srvG(1, "parse_complete", `{}`), parse(2), sync(2), srvG(2, "error_response", fatalFelder), srvG(2, "connection_closed", `{}`))),
		f("extended, letzte Gruppe ohne Sync und Flush", erw(parse(1), srvG(1, "error_response", fatalFelder), srvG(1, "connection_closed", `{}`))),
	})
}

// Abdeckung: LH-FA-25/Negative — der Leser von sqlite prüft die Version vor dem Inhalt und vergleicht sie als Zahl:
// jede kanonische Zahl außer 1 und 2 (0, 3, -1, 10, 20, eine über 64 Bit) ist PGR-E3002, auch bei einem Ende an
// falscher Stelle, und ihre Meldung enthält nicht "Version" mit der Zahl; andere Schreibweisen von 2 (02, +2, mit
// Leerzeichen davor oder danach) sind PGR-E3003.
func TestLesenEndeVersionen(t *testing.T) {
	gut := einSession(fertigeInteraktion())
	vorFatal := einSession(einf(ende(), fatal()))
	var faelle []endeFall
	for _, v := range []string{"0", "3", "-1", "10", "20", "99999999999999999999"} {
		faelle = append(faelle, endeFall{"Version " + v, v, "", false, gut, 0, unbekannteV})
	}
	for _, v := range []string{"02", "+2", " 2", "2 "} {
		faelle = append(faelle, endeFall{fmt.Sprintf("Version %q", v), v, "", false, gut, 0, beschaedigt})
	}
	faelle = append(faelle,
		endeFall{"Version 3, Ende vor der Fehlerantwort", "3", "", false, vorFatal, 0, unbekannteV},
		endeFall{"Version 3, Ende vor der Fehlerantwort, Zeilen umgekehrt", "3", "", true, vorFatal, 0, unbekannteV},
		endeFall{"Version 2, Ende vor der Fehlerantwort", "2", "", false, vorFatal, 0, beschaedigt},
	)
	pruefeEndeFaelle(t, faelle)
}

// Abdeckung: LH-FA-25/Negative — der Leser von sqlite lädt die Aufzeichnung mit Ende über Unmarshal, das
// Session.ValidateEnd je Session ruft: ein Ende in der Interaktion 1 von 2 und je ein Ende in zwei Interaktionen einer Session
// sind PGR-E3003, ein Ende in jeder von zwei Sessions wird geladen.
func TestLesenEndeJeSession(t *testing.T) {
	basis := sqliteDatei(t, sessionsMit(1))
	for _, c := range []struct {
		name     string
		sessions [][]ia
		want     string
	}{
		{"Interaktion 1 von 2", einSession(einf(fatal(), ende()), fertigeInteraktion()), model.CodeRecordingBroken},
		{"zwei Interaktionen", einSession(einf(fatal(), ende()), einf(fatal(), ende())), model.CodeRecordingBroken},
		{"je Session eines", [][]ia{{einf(fatal(), ende())}, {fertigeInteraktion(), einf(fatal(), ende())}}, ""},
	} {
		path := kopiere(t, basis)
		rohSchreibe(t, path, endeSQL("2", false, c.sessions...))
		_, err := laden(path)
		if codeVon(err) != c.want {
			t.Errorf("%s: Meldungscode %q, erwartet %q (%v)", c.name, codeVon(err), c.want, err)
		}
	}
}
