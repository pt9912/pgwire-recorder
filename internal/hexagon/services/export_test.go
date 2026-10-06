package services

import "github.com/pt9912/pgwire-recorder/internal/hexagon/model"

// Export-Test-Brücke (SPEC-049 Punkt 7): reicht an Unexportiertes weiter, ohne
// eigenen Zustand auf Paketebene.

// IstLebendpruefung reicht an istLebendpruefung weiter.
func IstLebendpruefung(sql string, vt bool) bool { return istLebendpruefung(sql, vt) }

// VTLeerraum reicht an vtLeerraum weiter.
func VTLeerraum(params map[string]string) bool { return vtLeerraum(params) }

// Abweichung reicht an abweichung weiter.
func Abweichung(empfangen, erwartet model.ClientMessage) string {
	return abweichung(empfangen, erwartet)
}

// ParameterStelle reicht an parameterStelle weiter.
func ParameterStelle(empfangen, erwartet []model.Value) string {
	return parameterStelle(empfangen, erwartet)
}

// LetzteNummer ist letzteNummer eines Cursors auf der übergebenen Session.
func LetzteNummer(session *model.Session) int {
	return (&cursor{session: session}).letzteNummer()
}
