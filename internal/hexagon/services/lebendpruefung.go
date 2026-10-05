package services

import "strings"

// leerraum sind die Zeichen, die der Scanner von PostgreSQL als Leerraum liest
// (LH-FA-09.a).
const leerraum = " \t\n\r\f\v"

// istLebendpruefung liefert true, wenn sql nur aus Leerraum, Zeilenkommentaren
// und geschlossenen, auch verschachtelten Blockkommentaren besteht; der leere
// Text zählt dazu (LH-FA-09.a, ADR-0031). Ein Zeilenkommentar reicht von `--`
// bis zum nächsten \n oder \r oder zum Textende; im Blockkommentar öffnet jedes
// `/*` eine Ebene und schließt jedes `*/` eine. Ein offener Blockkommentar und
// jedes andere Zeichen, auch `;`, machen den Text zu keiner Lebendprüfung.
func istLebendpruefung(sql string) bool {
	tiefe := 0
	for i := 0; i < len(sql); {
		rest := sql[i:]
		switch {
		case strings.HasPrefix(rest, "/*"):
			tiefe++
			i += 2
		case tiefe > 0 && strings.HasPrefix(rest, "*/"):
			tiefe--
			i += 2
		case tiefe > 0:
			i++
		case strings.HasPrefix(rest, "--"):
			ende := strings.IndexAny(rest, "\n\r")
			if ende < 0 {
				return true
			}
			i += ende + 1
		case strings.IndexByte(leerraum, rest[0]) >= 0:
			i++
		default:
			return false
		}
	}
	return tiefe == 0
}
