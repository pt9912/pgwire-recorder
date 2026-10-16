package recording_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

const (
	endeNachricht   = "{type: connection_closed}"
	fatalNachricht  = "{type: error_response, notice: {V: FATAL, S: FATAL, C: 57P01, M: x}}"
	fehlerNachricht = "{type: error_response, notice: {V: ERROR, S: ERROR, C: 42P01, M: x}}"
)

// endeKopf ist der Kopf einer Aufzeichnung der Version version; zusatz steht danach
// auf oberster Ebene (zum Beispiel empty_sessions), sessions folgt.
func endeKopf(version, zusatz string) string {
	return "format: pgwire-recorder\nversion: " + version + "\n" + zusatz + "sessions:\n"
}

// endeSession ist eine Session mit den Interaktionen in der Einrückung von einfach und
// erweitert.
func endeSession(id int, interaktionen ...string) string {
	return "  - id: " + strconv.Itoa(id) + "\n    interactions:\n" + strings.Join(interaktionen, "")
}

// einfach ist eine einfache Anfrage mit den Antworten; zusatz steht neben sequence.
func einfach(nummer int, zusatz string, antworten ...string) string {
	t := "      - sequence: " + strconv.Itoa(nummer) + "\n" + zusatz + "        request: {type: query, sql: x}\n        responses:\n"
	for _, a := range antworten {
		t += "          - " + a + "\n"
	}
	return t
}

// erweitert ist eine Extended-Interaktion mit Gruppen in der Einrückung von gruppe.
func erweitert(nummer int, gruppen ...string) string {
	return "      - sequence: " + strconv.Itoa(nummer) + "\n        type: extended\n        groups:\n" + strings.Join(gruppen, "")
}

// mitEnde ist eine Aufzeichnung der Version version mit einer einfachen Anfrage und einem Ende.
func mitEnde(version string) string {
	return endeKopf(version, "") + endeSession(1, einfach(1, "", fatalNachricht, endeNachricht))
}

// gruppeEnde ist eine Gruppe mit parse und abschluss als Client-Nachrichten.
func gruppeEnde(abschluss string, server ...string) string {
	return gruppe([]string{parseNachricht, abschluss}, server...)
}

// Abdeckung: LH-FA-25/Happy — der Leser lädt connection_closed ohne Felder als letzte Antwort
// der letzten Interaktion einer Session unmittelbar nach einer error_response FATAL oder PANIC
// (Schweregrad aus V, ohne V aus S), in einer einfachen Interaktion und in einer Extended-Interaktion,
// deren letzte Gruppe mit sync oder flush endet; auch mit offset_ms, empty_sessions und
// Version 2 ohne Ende. Das Feld version des Modells trägt die Version der Datei.
func TestUnmarshalEnde(t *testing.T) {
	cases := []struct {
		name, text string
		version    int
		ende       bool
	}{
		{"einfach, FATAL", mitEnde("2"), 2, true},
		{"einfach, PANIC", endeKopf("2", "") + endeSession(1, einfach(1, "", "{type: error_response, notice: {V: PANIC}}", endeNachricht)), 2, true},
		{"einfach, FATAL nur in S", endeKopf("2", "") + endeSession(1, einfach(1, "", "{type: error_response, notice: {S: FATAL}}", endeNachricht)), 2, true},
		{"einfach, PANIC nur in S", endeKopf("2", "") + endeSession(1, einfach(1, "", "{type: error_response, notice: {S: PANIC}}", endeNachricht)), 2, true},
		{"einfach, V FATAL gegen S ERROR", endeKopf("2", "") + endeSession(1, einfach(1, "", "{type: error_response, notice: {V: FATAL, S: ERROR}}", endeNachricht)), 2, true},
		{"einfach, Ende in Blockform", endeKopf("2", "") + endeSession(1, einfach(1, "", fatalNachricht, "type: connection_closed")), 2, true},
		{"extended, sync", endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde(syncNachricht, fatalNachricht, endeNachricht))), 2, true},
		{"extended, flush", endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde("{type: flush}", fatalNachricht, endeNachricht))), 2, true},
		{"extended, vordere Flush-Gruppe", endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde("{type: flush}", "{type: parse_complete}"), gruppeEnde(syncNachricht, fatalNachricht, endeNachricht))), 2, true},
		{"Ende in der letzten von zwei Interaktionen", endeKopf("2", "") + endeSession(1, einfach(1, "", "{type: command_complete, tag: x}", rfqNachricht), einfach(2, "", fatalNachricht, endeNachricht)), 2, true},
		{"mit offset_ms", endeKopf("2", "") + endeSession(1, einfach(1, "        offset_ms: 5\n", "{type: command_complete, tag: x}", rfqNachricht), einfach(2, "        offset_ms: 9\n", fatalNachricht, endeNachricht)), 2, true},
		{"mit empty_sessions", endeKopf("2", "empty_sessions: true\n") + endeSession(1, einfach(1, "", fatalNachricht, endeNachricht)) + "  - id: 2\n    interactions: []\n", 2, true},
		{"Version 2 ohne Ende, ready_for_query mitten in einer einfachen Anfrage", endeKopf("2", "") + endeSession(1, einfach(1, "", "{type: command_complete, tag: x}", rfqNachricht, "{type: command_complete, tag: x}", rfqNachricht)), 2, false},
		{"Version 2 ohne Ende", endeKopf("2", "") + endeSession(1, einfach(1, "", "{type: command_complete, tag: x}", rfqNachricht)), 2, false},
		{"Version 1 ohne Ende", endeKopf("1", "") + endeSession(1, einfach(1, "", "{type: command_complete, tag: x}", rfqNachricht)), 1, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec, err := recording.Unmarshal([]byte(c.text))
			if err != nil {
				t.Fatal(err)
			}
			if rec.Version != c.version {
				t.Errorf("Version %d statt %d", rec.Version, c.version)
			}
			hat := false
			for _, s := range rec.Sessions {
				for _, i := range s.Interactions {
					hat = hat || endeIn(i)
				}
			}
			if hat != c.ende {
				t.Errorf("Ende geladen: %v, erwartet %v", hat, c.ende)
			}
		})
	}
}

// endeIn sagt, ob die letzte Server-Nachricht der Interaktion connection_closed ist, ohne Felder.
func endeIn(i model.Interaction) bool {
	rs := i.Responses
	if len(i.Groups) > 0 {
		rs = i.Groups[len(i.Groups)-1].Server
	}
	return len(rs) > 0 && rs[len(rs)-1].Type == model.ResponseConnectionClosed && rs[len(rs)-1].TxStatus == "" && rs[len(rs)-1].Fields == nil
}

// Abdeckung: LH-FA-25/Negative — der Leser lehnt mit PGR-E3003 ab: connection_closed in Version 1, nicht
// als letzte Nachricht, nicht in der letzten Interaktion einer Session, zweimal in einer Session, mit
// einem Schlüssel neben type (mit Wert, leer, null), nicht nach FATAL oder PANIC (nach ERROR, ohne
// Fehlerantwort, V vor S, leeres V, fatal in Kleinschreibung), mit einer Nachricht zwischen
// Fehlerantwort und Ende, in einer Extended-Gruppe mit ready_for_query, in einer vorderen Gruppe und
// bei einer letzten Gruppe ohne Ende mit flush, mit ready_for_query an jeder Stelle einer einfachen Interaktion
// mit Ende (am Anfang, in der Mitte, direkt vor und nach der Fehlerantwort) und mit connection_closed an
// einer falschen Stelle (mitten in der letzten Interaktion, einfach und Extended, in einer Interaktion vor
// der letzten); die Meldung nennt die verletzte Regel, wo die Zeile sie nennt.
func TestUnmarshalEndeAbgelehnt(t *testing.T) {
	einfachEnde := func(antworten ...string) string {
		return endeKopf("2", "") + endeSession(1, einfach(1, "", antworten...))
	}
	const ccN = "{type: command_complete, tag: x}"
	extendedEndeText := func(server ...string) string {
		return endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde(syncNachricht, server...)))
	}
	cases := []struct{ name, text, want string }{
		{"ready_for_query am Anfang, Ende statt ready_for_query", einfachEnde(rfqNachricht, ccN, fatalNachricht, endeNachricht), "ready_for_query in einer Interaktion mit connection_closed"},
		{"ready_for_query in der Mitte, Ende statt ready_for_query", einfachEnde(ccN, rfqNachricht, ccN, fatalNachricht, endeNachricht), "ready_for_query in einer Interaktion mit connection_closed"},
		{"ready_for_query direkt vor der Fehlerantwort, Ende statt ready_for_query", einfachEnde(ccN, rfqNachricht, fatalNachricht, endeNachricht), "ready_for_query in einer Interaktion mit connection_closed"},
		{"ready_for_query zwischen Fehlerantwort und Ende", einfachEnde(fatalNachricht, rfqNachricht, endeNachricht), "nach einer error_response FATAL oder PANIC"},
		{"Ende mitten in der letzten Interaktion, einfach", endeKopf("2", "") + endeSession(1, einfach(1, "", ccN, rfqNachricht), einfach(2, "", fatalNachricht, endeNachricht, ccN)), "nur als letzte Nachricht"},
		{"Ende mitten in der letzten Gruppe, extended", endeKopf("2", "") + endeSession(1, einfach(1, "", ccN, rfqNachricht), erweitert(2, gruppeEnde(syncNachricht, fatalNachricht, endeNachricht, "{type: bind_complete}"))), "nur als letzte Nachricht"},
		{"Ende in einer vorderen Gruppe, extended, mit Nachricht danach", endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde("{type: flush}", fatalNachricht, endeNachricht, "{type: bind_complete}"), gruppeEnde(syncNachricht, ccN, rfqNachricht))), ""},
		{"falsch platziertes Ende in einer Interaktion vor der letzten, einfach", endeKopf("2", "") + endeSession(1, einfach(1, "", fatalNachricht, endeNachricht, ccN), einfach(2, "", ccN, rfqNachricht)), ""},
		{"falsch platziertes Ende in einer Interaktion vor der letzten, extended", endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde(syncNachricht, fatalNachricht, endeNachricht, "{type: bind_complete}")), einfach(2, "", ccN, rfqNachricht)), ""},
		{"Version 1, einfach", mitEnde("1"), "connection_closed in Version 1"},
		{"Version 1, extended", endeKopf("1", "") + endeSession(1, erweitert(1, gruppeEnde(syncNachricht, fatalNachricht, endeNachricht))), "connection_closed in Version 1"},
		{"nicht letzte Nachricht, einfach", einfachEnde(fatalNachricht, endeNachricht, rfqNachricht), "nur als letzte Nachricht"},
		{"nicht letzte Nachricht, extended", extendedEndeText(fatalNachricht, endeNachricht, "{type: notice_response}"), "nur als letzte Nachricht"},
		{"nicht letzte Interaktion, einfach", endeKopf("2", "") + endeSession(1, einfach(1, "", fatalNachricht, endeNachricht), einfach(2, "", "{type: command_complete, tag: x}", rfqNachricht)), "Session 1: Interaktion 1: connection_closed nur in der letzten Interaktion"},
		{"nicht letzte Interaktion, extended", endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde(syncNachricht, fatalNachricht, endeNachricht)), einfach(2, "", "{type: command_complete, tag: x}", rfqNachricht)), "Session 1: Interaktion 1: connection_closed nur in der letzten Interaktion"},
		{"zwei Enden in einer Session", endeKopf("2", "") + endeSession(1, einfach(1, "", fatalNachricht, endeNachricht), einfach(2, "", fatalNachricht, endeNachricht)), "Session 1: Interaktion 1: connection_closed nur in der letzten Interaktion"},
		{"Schlüssel mit Wert", einfachEnde(fatalNachricht, "{type: connection_closed, tag: x}"), "Schlüssel neben type"},
		{"Schlüssel leer", einfachEnde(fatalNachricht, "{type: connection_closed, tag: \"\"}"), "Schlüssel neben type"},
		{"Schlüssel null", einfachEnde(fatalNachricht, "{type: connection_closed, notice: null}"), "Schlüssel neben type"},
		{"Schlüssel unbekannt", einfachEnde(fatalNachricht, "{type: connection_closed, extra: 1}"), "Schlüssel neben type"},
		{"Schlüssel leer, extended", extendedEndeText(fatalNachricht, "{type: connection_closed, tx_status: \"\"}"), "Schlüssel neben type"},
		{"nach ERROR", einfachEnde(fehlerNachricht, endeNachricht), "nach einer error_response FATAL oder PANIC"},
		{"nach ERROR, extended", extendedEndeText(fehlerNachricht, endeNachricht), "nach einer error_response FATAL oder PANIC"},
		{"ohne Fehlerantwort", einfachEnde("{type: command_complete, tag: x}", endeNachricht), "nach einer error_response FATAL oder PANIC"},
		{"nur Ende", einfachEnde(endeNachricht), "nach einer error_response FATAL oder PANIC"},
		{"V ERROR gegen S FATAL", einfachEnde("{type: error_response, notice: {V: ERROR, S: FATAL}}", endeNachricht), "nach einer error_response FATAL oder PANIC"},
		{"leeres V gegen S FATAL", einfachEnde("{type: error_response, notice: {V: '', S: FATAL}}", endeNachricht), "nach einer error_response FATAL oder PANIC"},
		{"fatal klein", einfachEnde("{type: error_response, notice: {V: fatal}}", endeNachricht), "nach einer error_response FATAL oder PANIC"},
		{"Hinweis zwischen Fehler und Ende", einfachEnde(fatalNachricht, "{type: notice_response}", endeNachricht), "nach einer error_response FATAL oder PANIC"},
		{"ready_for_query in der Gruppe mit Ende", extendedEndeText(fatalNachricht, rfqNachricht, endeNachricht), "ready_for_query nur als letzte"},
		{"Ende in einer vorderen Gruppe", endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde("{type: flush}", fatalNachricht, endeNachricht), gruppeEnde(syncNachricht, "{type: command_complete, tag: x}", rfqNachricht))), "connection_closed nur in der letzten Gruppe"},
		{"vordere Gruppe mit sync vor Ende", endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde(syncNachricht, "{type: parse_complete}"), gruppeEnde(syncNachricht, fatalNachricht, endeNachricht))), "Gruppe 1: sync steht"},
		{"letzte Gruppe mit flush ohne Ende", endeKopf("2", "") + endeSession(1, erweitert(1, gruppeEnde("{type: flush}", "{type: parse_complete}"))), "Gruppe 1: sync steht"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := recording.Unmarshal([]byte(c.text))
			if code(err) != model.CodeRecordingBroken || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("erwartet %s mit %q, erhalten %v", model.CodeRecordingBroken, c.want, err)
			}
		})
	}
}

// Abdeckung: LH-FA-25/Boundary — ein Ende je Session ist erlaubt: zwei Sessions mit je einem Ende in
// ihrer letzten Interaktion laden.
func TestUnmarshalEndeJeSession(t *testing.T) {
	text := endeKopf("2", "") + endeSession(1, einfach(1, "", fatalNachricht, endeNachricht)) + endeSession(2, einfach(1, "", "{type: command_complete, tag: x}", rfqNachricht), einfach(2, "", fatalNachricht, endeNachricht))
	rec, err := recording.Unmarshal([]byte(text))
	if err != nil || len(rec.Sessions) != 2 {
		t.Fatalf("erwartet zwei Sessions, erhalten %d, %v", len(rec.Sessions), err)
	}
}

// Abdeckung: LH-FA-25/Negative — der Leser nimmt Version 1 und 2 an; jede andere kanonische Zahl ist
// PGR-E3002 (auch 3, 0, eine negative, eine über 64 Bit), vor dem Inhalt, auch bei connection_closed
// und bei Fehlern im Inhalt; Fehler im Inhalt einer Version 2 sind PGR-E3003.
func TestUnmarshalVersionen(t *testing.T) {
	rumpf := func(version, session string) string { return endeKopf(version, "") + session }
	gut := endeSession(1, einfach(1, "", "{type: command_complete, tag: x}", rfqNachricht))
	kaputt := endeSession(1, einfach(1, "", "{type: command_complete, tag: x}"))
	ende := endeSession(1, einfach(1, "", fatalNachricht, endeNachricht))
	cases := []struct{ name, text, code string }{
		{"Version 1", rumpf("1", gut), ""},
		{"Version 2", rumpf("2", gut), ""},
		{"Version 3", rumpf("3", gut), model.CodeRecordingVersion},
		{"Version 0", rumpf("0", gut), model.CodeRecordingVersion},
		{"Version -1", rumpf("-1", gut), model.CodeRecordingVersion},
		{"Version 99999999999999999999", rumpf("99999999999999999999", gut), model.CodeRecordingVersion},
		{"Version 3 mit Fehler im Inhalt", rumpf("3", kaputt), model.CodeRecordingVersion},
		{"Version 3 mit Ende", rumpf("3", ende), model.CodeRecordingVersion},
		{"Version 3 mit Schlüssel neben type", rumpf("3", endeSession(1, einfach(1, "", fatalNachricht, "{type: connection_closed, tag: x}"))), model.CodeRecordingVersion},
		{"Version 2 mit Fehler im Inhalt", rumpf("2", kaputt), model.CodeRecordingBroken},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := recording.Unmarshal([]byte(c.text))
			if code(err) != c.code {
				t.Fatalf("erwartet %q, erhalten %v", c.code, err)
			}
		})
	}
}
