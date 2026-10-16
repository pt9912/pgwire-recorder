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

// unterschied ist die erste Abweichung einer Interaktion: die Position in der
// normalisierten Folge (ab 1), die Art und, bei der Art Fehler, der SQLSTATE
// beider Seiten.
type unterschied struct {
	position   int
	art        string
	sqlstateIs string
	sqlstateSo string
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
	u := v.gewoehnlich()
	if u == nil || ende == 0 || u.position < ende {
		return u
	}
	f := &unterschied{position: ende, art: artFehler, sqlstateIs: v.ist[ende-1].Fields["C"]}
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

// gewoehnlich liefert die erste verschiedene Position der Folge des Servers und
// der aufgezeichneten, sonst die Position hinter der kürzeren, sonst nil.
func (v *vergleich) gewoehnlich() *unterschied {
	grenze := min(len(v.ist), len(v.soll))
	for i := range grenze {
		if u := nachrichtenUnterschied(v.ist[i], v.soll[i]); u != nil {
			u.position = i + 1
			return u
		}
	}
	if len(v.ist) != len(v.soll) {
		return &unterschied{position: grenze + 1, art: artNachricht}
	}
	return nil
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
