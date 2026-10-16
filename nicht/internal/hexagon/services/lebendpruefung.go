package services

import (
	"strconv"
	"strings"
)

const (
	// leerraum ist der Leerraum einer Lebendprüfung ohne den vertikalen
	// Tabulator (LH-FA-09.a).
	leerraum = " \t\n\r\f"
	// vtAbVersion ist die erste Hauptversion von PostgreSQL, ab der auch der
	// vertikale Tabulator (\v) zum Leerraum einer Lebendprüfung zählt
	// (LH-FA-09.a).
	vtAbVersion = 17
)

// istLebendpruefung liefert true, wenn sql nur aus Leerraum, Zeilenkommentaren
// und geschlossenen, auch verschachtelten Blockkommentaren besteht; der leere
// Text zählt dazu (LH-FA-09.a, ADR-0031). Leerraum ist leerraum, mit vt auch
// \v. Ein Zeilenkommentar reicht von `--` bis zum nächsten \n oder \r oder zum
// Textende; im Blockkommentar öffnet jedes `/*` eine Ebene und schließt jedes
// `*/` eine. Ein offener Blockkommentar und jedes andere Zeichen, auch `;`,
// machen den Text zu keiner Lebendprüfung.
func istLebendpruefung(sql string, vt bool) bool {
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
		case strings.IndexByte(leerraum, rest[0]) >= 0, vt && rest[0] == '\v':
			i++
		default:
			return false
		}
	}
	return tiefe == 0
}

// vtLeerraum liefert true, wenn die Hauptversion im Serverparameter
// server_version, die führenden Ziffern, mindestens vtAbVersion ist. Fehlt der
// Parameter oder beginnt er nicht mit einer Ziffer, liefert es false
// (LH-FA-09.a).
func vtLeerraum(params map[string]string) bool {
	v := params["server_version"]
	n := 0
	for n < len(v) && v[n] >= '0' && v[n] <= '9' {
		n++
	}
	major, err := strconv.Atoi(v[:n])
	return err == nil && major >= vtAbVersion
}
