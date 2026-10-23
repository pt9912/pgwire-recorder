package services

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Zaehler sind die Zähler der Zeile zum Ende von play (LH-FA-24.a *Meldung*,
// LH-FA-20.a Schritt 4): Eingespielt zählt die Interaktionen, deren erste
// Client-Nachricht der Service zu senden begann; Verglichen die Interaktionen,
// deren Antworten mit Vergleich geprüft wurden, bis zum ReadyForQuery, bis zur
// Fehlerantwort FATAL oder PANIC oder, bei der letzten Gruppe mit Flush, bis zur
// Antwort auf jede ihrer Client-Nachrichten; Abweichend die Interaktionen mit einer
// Abweichung, höchstens eine je Interaktion. Ohne CompareResponses sind
// Verglichen und Abweichend 0.
type Zaehler struct {
	Eingespielt int
	Verglichen  int
	Abweichend  int
}

// Die Arten einer Abweichung (LH-FA-24.a *Meldung*, *Art*).
const (
	artNachricht = "Nachrichtenart"
	artSpalten   = "Spaltenbeschreibung"
	artBefehl    = "Befehl"
	artFehler    = "Fehler"
	artStatus    = "Transaktionsstatus"
)

// vergleich hält die normalisierte Folge der Serverantworten einer Interaktion
// (ist) und die der Aufzeichnung (soll), bis die Interaktion zu Ende ist
// (LH-FA-20.a Schritt 4).
type vergleich struct {
	ist, soll []model.Response
}

// neuerVergleich bereitet den Vergleich einer Interaktion vor: soll ist die
// normalisierte Folge der aufgezeichneten Antworten, bei einer
// Extended-Interaktion die aller Gruppen ohne die Gruppengrenzen.
func neuerVergleich(in model.Interaction) *vergleich {
	var alle []model.Response
	if in.Request.Type == model.RequestExtended {
		for _, g := range in.Groups {
			alle = append(alle, g.Server...)
		}
	} else {
		alle = in.Responses
	}
	v := &vergleich{}
	for _, r := range alle {
		v.soll = nimmNormalisiert(v.soll, r)
	}
	return v
}

// nimm hängt eine gelesene Serverantwort an die Folge des Servers an, sofern die
// Normalisierung sie nicht entfernt.
func (v *vergleich) nimm(r model.Response) {
	v.ist = nimmNormalisiert(v.ist, r)
}

// nimmNormalisiert hängt r an folge an, außer r ist ein data_row, ein
// notice_response oder ein parameter_status (LH-FA-24.a, Einleitung).
func nimmNormalisiert(folge []model.Response, r model.Response) []model.Response {
	switch r.Type {
	case model.ResponseDataRow, model.ResponseNoticeResponse, model.ResponseParameterStatus:
		return folge
	}
	return append(folge, r)
}

// unterschied ist eine Abweichung einer Interaktion: die Position in der
// normalisierten Folge (ab 1), die Art, bei der Art Fehler der SQLSTATE beider
// Seiten und, wo beide Seiten an der Position eine Nachricht desselben Typs
// tragen, diese beiden Nachrichten (ist: der Server, soll: die Aufzeichnung) für
// die Einzelheiten der Meldung mit --diff-all.
type unterschied struct {
	position   int
	art        string
	sqlstateIs string
	sqlstateSo string
	ist, soll  model.Response
}

// text ist die Meldung der Abweichung: Position und Art, bei der Art Fehler der
// SQLSTATE beider Seiten; nie ein Wert, nie der Text einer Fehlerantwort
// (LH-FA-24.a *Meldung*). Session und Interaktion stellt der Aufrufer voran.
func (u unterschied) text() string {
	s := fmt.Sprintf("Antwort weicht von der Aufzeichnung ab: Position %d, Art %s", u.position, u.art)
	if u.art == artFehler {
		s += fmt.Sprintf(" (SQLSTATE aufgezeichnet „%s“, Server „%s“)", u.sqlstateSo, u.sqlstateIs)
	}
	return s
}

// maxPositionen und maxEinzelne sind die Obergrenzen der Meldung mit --diff-all:
// Positionen je Interaktion und abweichende Spalten oder Parameter je Position
// (LH-FA-24.a *Alle Abweichungen*, *Obergrenze*).
const (
	maxPositionen = 20
	maxEinzelne   = 20
)

// listeText ist die Meldung mit --diff-all für us, die Abweichungen einer
// Interaktion in aufsteigender Position und nicht leer: die Gesamtzahl, dann
// höchstens maxPositionen Positionen mit Art und Einzelheiten beider Seiten, dahinter
// „und N weitere“ (LH-FA-24.a *Form der Meldung*, *Obergrenze*). Session und
// Interaktion stellt der Aufrufer voran.
func listeText(us []unterschied) string {
	teile := make([]string, 0, maxPositionen+1)
	for i, u := range us {
		if i == maxPositionen {
			teile = append(teile, fmt.Sprintf("und %d weitere", len(us)-maxPositionen))
			break
		}
		t := fmt.Sprintf("Position %d, Art %s", u.position, u.art)
		if d := u.einzelheiten(); d != "" {
			t += " (" + d + ")"
		}
		teile = append(teile, t)
	}
	return fmt.Sprintf("Antworten weichen von der Aufzeichnung ab, Abweichungen: %d; %s", len(us), strings.Join(teile, "; "))
}

// einzelheiten nennt je Art beide Seiten (LH-FA-24.a *Einzelheiten je Art, beide
// Seiten*); die Art Nachrichtenart hat keine. Nie ein Zeilenwert und nie der Text
// einer Fehlerantwort.
func (u unterschied) einzelheiten() string {
	switch u.art {
	case artFehler:
		return fmt.Sprintf("SQLSTATE aufgezeichnet „%s“, Server „%s“", u.sqlstateSo, u.sqlstateIs)
	case artBefehl:
		return fmt.Sprintf("Befehl aufgezeichnet „%s“, Server „%s“", befehl(u.soll.Tag), befehl(u.ist.Tag))
	case artStatus:
		return fmt.Sprintf("aufgezeichnet „%s“, Server „%s“", u.soll.TxStatus, u.ist.TxStatus)
	case artSpalten:
		if u.ist.Type == model.ResponseParameterDescription {
			return parameterEinzelheiten(u.ist.ParamTypes, u.soll.ParamTypes)
		}
		return spaltenEinzelheiten(u.ist.Columns, u.soll.Columns)
	}
	return ""
}

// spaltenEinzelheiten nennt die Zahl der Spalten beider Seiten, wenn sie
// verschieden ist, und je abweichende Spalte der gemeinsamen Nummern Nummer, Name
// und Typ-OID beider Seiten, soweit sie abweichen, der Name zuerst; höchstens
// maxEinzelne Spalten, dahinter „und N weitere“.
func spaltenEinzelheiten(ist, soll []model.Column) string {
	var teile []string
	if len(ist) != len(soll) {
		teile = append(teile, fmt.Sprintf("Spalten aufgezeichnet %d, Server %d", len(soll), len(ist)))
	}
	var spalten []string
	for i := range min(len(ist), len(soll)) {
		var felder []string
		if ist[i].Name != soll[i].Name {
			felder = append(felder, fmt.Sprintf("Name aufgezeichnet „%s“, Server „%s“", soll[i].Name, ist[i].Name))
		}
		if ist[i].TypeOID != soll[i].TypeOID {
			felder = append(felder, fmt.Sprintf("Typ-OID aufgezeichnet %d, Server %d", soll[i].TypeOID, ist[i].TypeOID))
		}
		if len(felder) > 0 {
			spalten = append(spalten, fmt.Sprintf("Spalte %d: %s", i+1, strings.Join(felder, ", ")))
		}
	}
	return strings.Join(append(teile, begrenze(spalten)...), "; ")
}

// parameterEinzelheiten nennt die Zahl der Parameter beider Seiten, wenn sie
// verschieden ist, und je abweichenden Parameter der gemeinsamen Nummern die
// Typ-OID beider Seiten; höchstens maxEinzelne Parameter, dahinter „und N weitere“.
func parameterEinzelheiten(ist, soll []uint32) string {
	var teile []string
	if len(ist) != len(soll) {
		teile = append(teile, fmt.Sprintf("Parameter aufgezeichnet %d, Server %d", len(soll), len(ist)))
	}
	var parameter []string
	for i := range min(len(ist), len(soll)) {
		if ist[i] != soll[i] {
			parameter = append(parameter, fmt.Sprintf("Parameter %d: Typ-OID aufgezeichnet %d, Server %d", i+1, soll[i], ist[i]))
		}
	}
	return strings.Join(append(teile, begrenze(parameter)...), "; ")
}

// begrenze kürzt einzelne auf maxEinzelne Einträge und hängt „und N weitere“ an,
// wenn es mehr gibt.
func begrenze(einzelne []string) []string {
	if len(einzelne) <= maxEinzelne {
		return einzelne
	}
	return append(einzelne[:maxEinzelne:maxEinzelne], fmt.Sprintf("und %d weitere", len(einzelne)-maxEinzelne))
}

// erste liefert die erste Abweichung der Folge des Servers von der
// aufgezeichneten oder nil (LH-FA-24.a *Randformen des Vergleichs*): die erste
// Position, an der die Nachrichten abweichen; stimmen alle Positionen der
// kürzeren Folge überein und die Folgen sind verschieden lang, ist es die
// Position hinter der kürzeren, die Art Nachrichtenart. ende ist die Position der
// Fehlerantwort FATAL oder PANIC des Servers, hinter der connection_closed
// angefügt ist, sonst 0. Weicht die Folge an ihr oder dahinter ab, ist die
// Abweichung die Art Fehler an dieser Position, mit dem SQLSTATE beider Seiten
// (der aufgezeichnete leer, wenn dort keine Fehlerantwort steht); weicht eine
// frühere Position ab, gilt diese. Ohne Abweichung ist das Ende erwartet: gleicher
// SQLSTATE und aufgezeichnetes Ende an dieser Position (LH-FA-24.a *Ende der
// Aufzeichnung*).
func (v *vergleich) erste(ende int) *unterschied {
	us := v.alle(ende)
	if len(us) == 0 {
		return nil
	}
	return &us[0]
}

// alle liefert alle Abweichungen der Folge des Servers von der aufgezeichneten in
// aufsteigender Position, je Position höchstens eine (LH-FA-24.a *Alle
// Abweichungen*, *Zuordnung der Positionen*): jede Position, an der die Nachrichten
// abweichen, und bei verschieden langen Folgen jede Position hinter der kürzeren
// bis zum Ende der längeren, die Art Nachrichtenart. Mit ende (siehe erste) ist die
// Position des Endes die letzte: Abweichungen davor bleiben, weicht die Folge an
// ihr oder dahinter ab, ist die Position des Endes eine Abweichung der Art Fehler,
// und dahinter erscheint nichts.
func (v *vergleich) alle(ende int) []unterschied {
	var us []unterschied
	for _, u := range v.gewoehnlich() {
		if ende != 0 && u.position >= ende {
			return append(us, v.endeFehler(ende))
		}
		us = append(us, u)
	}
	return us
}

// endeFehler ist die Abweichung der Art Fehler an der Position des Endes des
// Servers.
func (v *vergleich) endeFehler(ende int) unterschied {
	f := unterschied{position: ende, art: artFehler, sqlstateIs: v.ist[ende-1].Fields["C"]}
	// Eine Aufzeichnung, die Validate besteht, endet mit ready_for_query oder
	// connection_closed, die Folge des Servers trägt beides nicht vor der
	// Fehlerantwort: soll reicht bis zu ihrer Position, der Schutz der Länge greift
	// nur bei einer Folge, die Validate nicht besteht. Den SQLSTATE liefert nur eine
	// error_response (LH-FA-24.a *Fehler*).
	if ende <= len(v.soll) && v.soll[ende-1].Type == model.ResponseErrorResponse {
		f.sqlstateSo = v.soll[ende-1].Fields["C"]
	}
	return f
}

// gewoehnlich liefert alle Abweichungen der Folge des Servers von der
// aufgezeichneten nach den gewöhnlichen Regeln, ohne Ende: je Position beider
// Folgen, an der die Nachrichten verschieden sind, und je Position hinter der
// kürzeren Folge bis zum Ende der längeren.
func (v *vergleich) gewoehnlich() []unterschied {
	var us []unterschied
	grenze := min(len(v.ist), len(v.soll))
	for i := range grenze {
		if u := nachrichtenUnterschied(v.ist[i], v.soll[i]); u != nil {
			u.position = i + 1
			u.ist, u.soll = v.ist[i], v.soll[i]
			us = append(us, *u)
		}
	}
	for p := grenze + 1; p <= max(len(v.ist), len(v.soll)); p++ {
		us = append(us, unterschied{position: p, art: artNachricht})
	}
	return us
}

// nachrichtenUnterschied vergleicht zwei Nachrichten an derselben Position nach
// der Tabelle in LH-FA-24.a; die Position setzt der Aufrufer.
func nachrichtenUnterschied(ist, soll model.Response) *unterschied {
	if ist.Type != soll.Type {
		return &unterschied{art: artNachricht}
	}
	switch ist.Type {
	case model.ResponseRowDescription:
		if !gleicheSpalten(ist.Columns, soll.Columns) {
			return &unterschied{art: artSpalten}
		}
	case model.ResponseParameterDescription:
		if !slices.Equal(ist.ParamTypes, soll.ParamTypes) {
			return &unterschied{art: artSpalten}
		}
	case model.ResponseCommandComplete:
		if befehl(ist.Tag) != befehl(soll.Tag) {
			return &unterschied{art: artBefehl}
		}
	case model.ResponseErrorResponse:
		if ist.Fields["C"] != soll.Fields["C"] {
			return &unterschied{art: artFehler, sqlstateIs: ist.Fields["C"], sqlstateSo: soll.Fields["C"]}
		}
	case model.ResponseReadyForQuery:
		if ist.TxStatus != soll.TxStatus {
			return &unterschied{art: artStatus}
		}
	}
	return nil
}

// gleicheSpalten meldet, ob beide Beschreibungen gleich viele Spalten haben und
// je Spalte Name und Typ-OID gleich sind.
func gleicheSpalten(a, b []model.Column) bool {
	return slices.EqualFunc(a, b, func(x, y model.Column) bool {
		return x.Name == y.Name && x.TypeOID == y.TypeOID
	})
}

// befehl ist der Tag ohne die abschließenden Zahlen: die Wörter am Ende, durch
// Leerzeichen getrennt, die nur aus Ziffern bestehen (LH-FA-24.a *Befehl*).
func befehl(tag string) string {
	woerter := strings.Split(tag, " ")
	for len(woerter) > 0 && nurZiffern(woerter[len(woerter)-1]) {
		woerter = woerter[:len(woerter)-1]
	}
	return strings.Join(woerter, " ")
}

// nurZiffern meldet, ob w nicht leer ist und nur aus den Ziffern 0 bis 9
// besteht.
func nurZiffern(w string) bool {
	return w != "" && strings.Trim(w, "0123456789") == ""
}
