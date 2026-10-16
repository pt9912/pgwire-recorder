package model_test

import (
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

func fehler(felder map[string]string) model.Response {
	return model.Response{Type: model.ResponseErrorResponse, Fields: felder}
}

func ende() model.Response { return model.Response{Type: model.ResponseConnectionClosed} }

// queryEnde ist eine einfache Anfrage, die mit den Antworten rs endet.
func queryEnde(rs ...model.Response) model.Interaction {
	return model.Interaction{
		Sequence:  1,
		Request:   model.Request{Type: model.RequestQuery, SQL: "SELECT 1"},
		Responses: rs,
	}
}

// extendedEnde ist eine Extended-Interaktion aus einer Flush-Gruppe und einer
// letzten Gruppe, die mit dem Client-Nachrichtentyp letzter endet und die
// Server-Nachrichten rs trägt.
func extendedEnde(letzter model.ClientMessageType, rs ...model.Response) model.Interaction {
	i := extendedBeispiel()
	i.Groups[1].Client[2].Type = letzter
	i.Groups[1].Server = rs
	return i
}

func fatalFelder() map[string]string {
	return map[string]string{"V": "FATAL", "S": "FATAL", "C": "57P01"}
}

// Gültige Formen des Endes (LH-FA-25.a): FATAL und PANIC, der Schweregrad aus V oder, ohne V
// oder mit leerem V, aus S, einfache Anfrage und Extended mit sync oder flush am Ende der letzten Gruppe.
func TestValidateEndeGueltig(t *testing.T) {
	for name, i := range map[string]model.Interaction{
		"query FATAL":                             queryEnde(fehler(fatalFelder()), ende()),
		"query PANIC":                             queryEnde(fehler(map[string]string{"V": "PANIC"}), ende()),
		"query FATAL nur in S":                    queryEnde(fehler(map[string]string{"S": "FATAL"}), ende()),
		"query PANIC nur in S":                    queryEnde(fehler(map[string]string{"S": "PANIC"}), ende()),
		"query V FATAL gegen S ERROR":             queryEnde(fehler(map[string]string{"V": "FATAL", "S": "ERROR"}), ende()),
		"query leeres V gegen S FATAL":            queryEnde(fehler(map[string]string{"V": "", "S": "FATAL"}), ende()),
		"query nach Hinweis und Daten":            queryEnde(model.Response{Type: model.ResponseNoticeResponse}, model.Response{Type: model.ResponseDataRow}, fehler(fatalFelder()), ende()),
		"extended mit sync":                       extendedEnde(model.ClientSync, fehler(fatalFelder()), ende()),
		"extended mit flush":                      extendedEnde(model.ClientFlush, fehler(fatalFelder()), ende()),
		"query ohne Ende, ready_for_query mitten": queryEnde(model.Response{Type: model.ResponseCommandComplete}, model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"}, model.Response{Type: model.ResponseCommandComplete}, model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"}),
		"extended PANIC":                          extendedEnde(model.ClientSync, model.Response{Type: model.ResponseBindComplete}, fehler(map[string]string{"V": "PANIC"}), ende()),
	} {
		if err := i.Validate(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// Jede Regel des Endes hat einen Fall, der nur sie verletzt; der erwartete Text belegt, dass
// die gemeinte Regel greift.
func TestValidateEndeFehler(t *testing.T) {
	const (
		letzte   = "connection_closed nur als letzte Nachricht"
		felder   = "connection_closed trägt Felder"
		fatalNur = "connection_closed nur unmittelbar nach einer error_response FATAL oder PANIC"
		rfqEnde  = "ready_for_query in einer Interaktion mit connection_closed"
	)
	cc := model.Response{Type: model.ResponseCommandComplete}
	mitFeld := func(r model.Response) model.Response {
		r.Type = model.ResponseConnectionClosed
		return r
	}
	rfq := model.Response{Type: model.ResponseReadyForQuery, TxStatus: "I"}
	cases := []struct {
		name string
		i    model.Interaction
		want string
	}{
		{"query nicht letzte Nachricht", queryEnde(fehler(fatalFelder()), ende(), rfq), letzte},
		{"query Ende in der Mitte, Nachricht ohne ready_for_query danach", queryEnde(fehler(fatalFelder()), ende(), cc), letzte},
		{"query ready_for_query am Anfang, Ende statt ready_for_query", queryEnde(rfq, cc, fehler(fatalFelder()), ende()), rfqEnde},
		{"query ready_for_query in der Mitte, Ende statt ready_for_query", queryEnde(cc, rfq, cc, fehler(fatalFelder()), ende()), rfqEnde},
		{"query ready_for_query direkt vor der Fehlerantwort, Ende statt ready_for_query", queryEnde(cc, rfq, fehler(fatalFelder()), ende()), rfqEnde},
		{"query ready_for_query zwischen Fehlerantwort und Ende", queryEnde(fehler(fatalFelder()), rfq, ende()), fatalNur},
		{"query nicht letzte Nachricht, davor kein Fehler", queryEnde(ende(), fehler(fatalFelder()), ende()), letzte},
		{"query mit Spalten", queryEnde(fehler(fatalFelder()), mitFeld(model.Response{Columns: []model.Column{{Name: "a"}}})), felder},
		{"query mit Werten", queryEnde(fehler(fatalFelder()), mitFeld(model.Response{Values: []model.Value{{Null: true}}})), felder},
		{"query mit Tag", queryEnde(fehler(fatalFelder()), mitFeld(model.Response{Tag: "x"})), felder},
		{"query mit Feldern", queryEnde(fehler(fatalFelder()), mitFeld(model.Response{Fields: map[string]string{"V": "FATAL"}})), felder},
		{"query mit Name", queryEnde(fehler(fatalFelder()), mitFeld(model.Response{Name: "x"})), felder},
		{"query mit Value", queryEnde(fehler(fatalFelder()), mitFeld(model.Response{Value: "x"})), felder},
		{"query mit TxStatus", queryEnde(fehler(fatalFelder()), mitFeld(model.Response{TxStatus: "I"})), felder},
		{"query mit Parametertypen", queryEnde(fehler(fatalFelder()), mitFeld(model.Response{ParamTypes: []uint32{23}})), felder},
		{"query nur connection_closed", queryEnde(ende()), fatalNur},
		{"query ohne Fehlerantwort davor", queryEnde(model.Response{Type: model.ResponseCommandComplete}, ende()), fatalNur},
		{"query Hinweis mit FATAL statt Fehlerantwort", queryEnde(model.Response{Type: model.ResponseNoticeResponse, Fields: fatalFelder()}, ende()), fatalNur},
		{"query nach ERROR", queryEnde(fehler(map[string]string{"V": "ERROR", "S": "ERROR"}), ende()), fatalNur},
		{"query ohne Schweregrad", queryEnde(fehler(map[string]string{"C": "57P01"}), ende()), fatalNur},
		{"query ohne Felder", queryEnde(fehler(nil), ende()), fatalNur},
		{"query V ERROR gegen S FATAL", queryEnde(fehler(map[string]string{"V": "ERROR", "S": "FATAL"}), ende()), fatalNur},
		{"query fatal klein in V", queryEnde(fehler(map[string]string{"V": "fatal"}), ende()), fatalNur},
		{"query fatal klein in S", queryEnde(fehler(map[string]string{"S": "fatal"}), ende()), fatalNur},
		{"query panic klein in V", queryEnde(fehler(map[string]string{"V": "panic"}), ende()), fatalNur},
		{"query Hinweis zwischen Fehler und Ende", queryEnde(fehler(fatalFelder()), model.Response{Type: model.ResponseNoticeResponse}, ende()), fatalNur},
		{"query FATAL mit Leerzeichen", queryEnde(fehler(map[string]string{"V": "FATAL "}), ende()), fatalNur},
		{"extended nicht letzte Nachricht", extendedEnde(model.ClientSync, fehler(fatalFelder()), ende(), model.Response{Type: model.ResponseNoticeResponse}), letzte},
		{"extended Ende mitten in der letzten Gruppe", extendedEnde(model.ClientSync, fehler(fatalFelder()), ende(), model.Response{Type: model.ResponseBindComplete}), letzte},
		{"extended mit Feld", extendedEnde(model.ClientSync, fehler(fatalFelder()), mitFeld(model.Response{Tag: "x"})), felder},
		{"extended nach ERROR", extendedEnde(model.ClientSync, fehler(map[string]string{"V": "ERROR"}), ende()), fatalNur},
		{"extended ohne Fehlerantwort davor", extendedEnde(model.ClientSync, model.Response{Type: model.ResponseBindComplete}, ende()), fatalNur},
		{"extended Hinweis zwischen Fehler und Ende", extendedEnde(model.ClientSync, fehler(fatalFelder()), model.Response{Type: model.ResponseNoticeResponse}, ende()), fatalNur},
		{"extended mit ready_for_query in der Gruppe mit Ende", extendedEnde(model.ClientSync, fehler(fatalFelder()), rfq, ende()), "ready_for_query nur als letzte Nachricht der letzten Gruppe"},
		{"extended mit ready_for_query nach dem Ende", extendedEnde(model.ClientSync, fehler(fatalFelder()), ende(), rfq), letzte},
		{"extended Ende in einer vorderen Gruppe", func() model.Interaction {
			i := extendedEnde(model.ClientSync, fehler(fatalFelder()), ende())
			i.Groups[0].Server = []model.Response{fehler(fatalFelder()), ende()}
			return i
		}(), "Gruppe 1: connection_closed nur in der letzten Gruppe"},
		{"extended Ende in der vorderen Gruppe, letzte Gruppe regulär", func() model.Interaction {
			i := extendedBeispiel()
			i.Groups[0].Server = []model.Response{fehler(fatalFelder()), ende()}
			return i
		}(), "Gruppe 1: connection_closed nur in der letzten Gruppe"},
		{"extended vordere Gruppe mit sync, letzte mit Ende", func() model.Interaction {
			i := extendedEnde(model.ClientSync, fehler(fatalFelder()), ende())
			i.Groups[0].Client[2].Type = model.ClientSync
			return i
		}(), "Gruppe 1: sync steht am Ende"},
		{"extended ohne Ende, letzte Gruppe mit flush", extendedEnde(model.ClientFlush, model.Response{Type: model.ResponseBindComplete}), "Gruppe 2: sync steht am Ende"},
		{"extended ohne Ende und ohne ready_for_query", extendedEnde(model.ClientSync, model.Response{Type: model.ResponseBindComplete}), "Gruppe 2: letzte Gruppe endet nicht mit ready_for_query"},
		{"query ohne Ende und ohne ready_for_query", queryEnde(fehler(fatalFelder())), "endet nicht mit ready_for_query"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.i.Validate()
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("erwartet Fehler mit %q, erhalten %v", c.want, err)
			}
		})
	}
}

// Die Formregel einer Session (SPEC-002): connection_closed steht höchstens einmal, in der
// letzten Interaktion; sie gilt für einfache und Extended-Interaktionen.
func TestSessionValidateEnd(t *testing.T) {
	einfach := queryEnde(fehler(fatalFelder()), ende())
	normal := queryBeispiel()
	gruppen := extendedEnde(model.ClientSync, fehler(fatalFelder()), ende())
	erweitert := extendedBeispiel()
	cases := []struct {
		name string
		ints []model.Interaction
		want string
	}{
		{"ohne Interaktion", nil, ""},
		{"ohne Ende", []model.Interaction{normal, erweitert}, ""},
		{"Ende in der einzigen Interaktion", []model.Interaction{einfach}, ""},
		{"Ende in der letzten von zwei", []model.Interaction{normal, einfach}, ""},
		{"Extended-Ende in der letzten von drei", []model.Interaction{normal, erweitert, gruppen}, ""},
		{"einfaches Ende in der ersten von zwei", []model.Interaction{einfach, normal}, "Interaktion 1:"},
		{"Extended-Ende in der ersten von zwei", []model.Interaction{gruppen, normal}, "Interaktion 1:"},
		{"einfaches Ende in der mittleren von dreien", []model.Interaction{normal, einfach, normal}, "Interaktion 2:"},
		{"Extended-Ende in der mittleren von dreien", []model.Interaction{normal, gruppen, normal}, "Interaktion 2:"},
		{"zwei Enden in aufeinanderfolgenden Interaktionen", []model.Interaction{einfach, einfach}, "Interaktion 1:"},
		{"zwei Enden, vorletztes Extended", []model.Interaction{normal, gruppen, einfach}, "Interaktion 2:"},
		{"falsch platziertes Ende in der ersten von zwei, einfach", []model.Interaction{queryEnde(fehler(fatalFelder()), ende(), model.Response{Type: model.ResponseCommandComplete}), normal}, "Interaktion 1:"},
		{"falsch platziertes Ende am Anfang der ersten von zwei, einfach", []model.Interaction{queryEnde(ende(), fehler(fatalFelder()), model.Response{Type: model.ResponseCommandComplete}), normal}, "Interaktion 1:"},
		{"falsch platziertes Ende in der letzten Gruppe der ersten von zwei, Extended", []model.Interaction{extendedEnde(model.ClientSync, fehler(fatalFelder()), ende(), model.Response{Type: model.ResponseBindComplete}), normal}, "Interaktion 1:"},
		{"falsch platziertes Ende in einer vorderen Gruppe der ersten von zwei, Extended", []model.Interaction{func() model.Interaction {
			i := extendedBeispiel()
			i.Groups[0].Server = []model.Response{fehler(fatalFelder()), ende(), model.Response{Type: model.ResponseBindComplete}}
			return i
		}(), normal}, "Interaktion 1:"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := model.Session{ID: 1, Interactions: c.ints}.ValidateEnd()
			switch {
			case c.want == "" && err != nil:
				t.Fatalf("erwartet kein Fehler, erhalten %v", err)
			case c.want != "" && (err == nil || !strings.Contains(err.Error(), c.want)):
				t.Fatalf("erwartet Fehler mit %q, erhalten %v", c.want, err)
			}
		})
	}
}

// Version 2 ist eine eigene Konstante; model.Version bleibt die der übrigen Formate.
func TestVersionMitEnde(t *testing.T) {
	if model.Version != 1 || model.VersionMitEnde != 2 || model.NewRecording().Version != 1 {
		t.Fatalf("Version %d, VersionMitEnde %d, NewRecording %d", model.Version, model.VersionMitEnde, model.NewRecording().Version)
	}
}
