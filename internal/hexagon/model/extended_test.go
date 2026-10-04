package model

import (
	"strings"
	"testing"
)

// extendedBeispiel ist die Interaktion aus SPEC-041 mit einer vorangestellten
// Flush-Gruppe: parse/describe/flush, danach bind/execute/sync.
func extendedBeispiel() Interaction {
	return Interaction{
		Sequence: 1,
		Request:  Request{Type: RequestExtended},
		Groups: []Group{
			{
				Client: []ClientMessage{
					{Type: ClientParse, Statement: "s1", SQL: "SELECT name FROM users WHERE id = $1", ParamTypes: []uint32{23}},
					{Type: ClientDescribe, Target: TargetStatement, Name: "s1"},
					{Type: ClientFlush},
				},
				Server: []Response{
					{Type: ResponseParseComplete},
					{Type: ResponseParameterDescription, ParamTypes: []uint32{23}},
					{Type: ResponseRowDescription, Columns: []Column{{Name: "name", TypeOID: 25, TypeSize: -1, TypeModifier: -1}}},
				},
			},
			{
				Client: []ClientMessage{
					{Type: ClientBind, Statement: "s1", ParamFormats: []int16{0}, Params: []Value{{Bytes: []byte("1")}}, ResultFormats: []int16{0}},
					{Type: ClientExecute, MaxRows: 0},
					{Type: ClientSync},
				},
				Server: []Response{
					{Type: ResponseBindComplete},
					{Type: ResponseDataRow, Values: []Value{{Bytes: []byte("alice")}}},
					{Type: ResponseCommandComplete, Tag: "SELECT 1"},
					{Type: ResponseReadyForQuery, TxStatus: "I"},
				},
			},
		},
	}
}

func queryBeispiel() Interaction {
	return Interaction{
		Sequence:  1,
		Request:   Request{Type: RequestQuery, SQL: "SELECT 1"},
		Responses: []Response{{Type: ResponseCommandComplete, Tag: "SELECT 1"}, {Type: ResponseReadyForQuery, TxStatus: "I"}},
	}
}

// Gültige Formen: einfache Anfrage, Extended mit Flush- und Sync-Gruppe, Extended
// mit nur einer Sync-Gruppe, eine Flush-Gruppe ohne Server-Nachricht, jede
// Client- und Server-Nachricht aus LH-FA-18.a einschließlich parameter_status.
func TestValidateGueltig(t *testing.T) {
	nurSync := Interaction{Request: Request{Type: RequestExtended}, Groups: []Group{{
		Client: []ClientMessage{{Type: ClientSync}},
		Server: []Response{{Type: ResponseReadyForQuery, TxStatus: "I"}},
	}}}
	alle := Interaction{Request: Request{Type: RequestExtended}, Groups: []Group{
		{
			Client: []ClientMessage{{Type: ClientClose, Target: TargetPortal}, {Type: ClientFlush}},
		},
		{
			Client: []ClientMessage{
				{Type: ClientParse}, {Type: ClientBind}, {Type: ClientDescribe, Target: TargetPortal},
				{Type: ClientExecute, MaxRows: 1}, {Type: ClientClose, Target: TargetStatement}, {Type: ClientSync},
			},
			Server: []Response{
				{Type: ResponseCloseComplete}, {Type: ResponseParseComplete}, {Type: ResponseBindComplete},
				{Type: ResponseNoData}, {Type: ResponseParameterDescription}, {Type: ResponseRowDescription},
				{Type: ResponseDataRow}, {Type: ResponsePortalSuspended}, {Type: ResponseCommandComplete},
				{Type: ResponseEmptyQueryResponse}, {Type: ResponseNoticeResponse}, {Type: ResponseParameterStatus},
				{Type: ResponseErrorResponse}, {Type: ResponseCloseComplete}, {Type: ResponseReadyForQuery},
			},
		},
	}}
	for name, i := range map[string]Interaction{
		"query":                       queryBeispiel(),
		"extended mit Flush und Sync": extendedBeispiel(),
		"extended nur Sync":           nurSync,
		"alle Nachrichtentypen":       alle,
	} {
		if err := i.Validate(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// Jede Formregel einer Extended-Interaktion hat einen Fall, der sie verletzt und
// abgelehnt wird; der erwartete Text belegt, dass die gemeinte Regel greift.
func TestValidateFehler(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Interaction)
		want   string
	}{
		{"unbekannter Anfrage-Typ", func(i *Interaction) { i.Request.Type = "bogus" }, "Anfrage-Typ"},
		{"leerer Anfrage-Typ", func(i *Interaction) { i.Request.Type = "" }, "Anfrage-Typ"},
		{"extended mit SQL", func(i *Interaction) { i.Request.SQL = "SELECT 1" }, "außerhalb der Gruppen"},
		{"extended mit Antworten", func(i *Interaction) {
			i.Responses = []Response{{Type: ResponseReadyForQuery}}
		}, "außerhalb der Gruppen"},
		{"extended ohne Gruppe", func(i *Interaction) { i.Groups = nil }, "ohne Gruppe"},
		{"Gruppe ohne Client-Nachricht", func(i *Interaction) { i.Groups[0].Client = nil }, "ohne Client-Nachricht"},
		{"unbekannte Client-Nachricht", func(i *Interaction) { i.Groups[0].Client[0].Type = "copy_data" }, "Client-Nachricht \"copy_data\" unbekannt"},
		{"flush mitten in der Gruppe", func(i *Interaction) { i.Groups[0].Client[1] = ClientMessage{Type: ClientFlush} }, "nur als letzte Nachricht"},
		{"sync mitten in der Gruppe", func(i *Interaction) { i.Groups[1].Client[0] = ClientMessage{Type: ClientSync} }, "nur als letzte Nachricht"},
		{"Gruppe ohne flush oder sync am Ende", func(i *Interaction) {
			i.Groups[0].Client = i.Groups[0].Client[:2]
		}, "nur als letzte Nachricht"},
		{"vordere Gruppe endet mit sync", func(i *Interaction) { i.Groups[0].Client[2].Type = ClientSync }, "Gruppe 1: sync steht"},
		{"letzte Gruppe endet mit flush (abgeschnitten)", func(i *Interaction) { i.Groups[1].Client[2].Type = ClientFlush }, "Gruppe 2: sync steht"},
		{"describe ohne Zielart", func(i *Interaction) { i.Groups[0].Client[1].Target = "" }, "Zielart"},
		{"close mit falscher Zielart", func(i *Interaction) {
			i.Groups[0].Client[1] = ClientMessage{Type: ClientClose, Target: "tabelle"}
		}, "Zielart"},
		{"unbekannte Server-Nachricht", func(i *Interaction) { i.Groups[1].Server[1].Type = "copy_out_response" }, "Server-Nachricht \"copy_out_response\" unbekannt"},
		{"ready_for_query in der Flush-Gruppe", func(i *Interaction) {
			i.Groups[0].Server = append(i.Groups[0].Server, Response{Type: ResponseReadyForQuery})
		}, "Gruppe 1: Server-Nachricht 4: ready_for_query"},
		{"ready_for_query vor dem Ende der letzten Gruppe", func(i *Interaction) {
			i.Groups[1].Server[2].Type = ResponseReadyForQuery
		}, "Gruppe 2: Server-Nachricht 3: ready_for_query"},
		{"letzte Gruppe ohne ready_for_query (abgeschnitten)", func(i *Interaction) {
			i.Groups[1].Server = i.Groups[1].Server[:3]
		}, "endet nicht mit ready_for_query"},
		{"letzte Gruppe ohne Server-Nachricht", func(i *Interaction) { i.Groups[1].Server = nil }, "endet nicht mit ready_for_query"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := extendedBeispiel()
			c.mutate(&i)
			err := i.Validate()
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("erwartet Fehler mit %q, erhalten %v", c.want, err)
			}
		})
	}
}

// Die Simple-Regeln gelten weiter, und eine Nachricht, die nur das Extended
// Query Protocol kennt, ist in einer einfachen Anfrage unbekannt.
func TestValidateQueryFehler(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Interaction)
		want   string
	}{
		{"ohne Antwort", func(i *Interaction) { i.Responses = nil }, "endet nicht mit ready_for_query"},
		{"ohne ready_for_query am Ende", func(i *Interaction) { i.Responses = i.Responses[:1] }, "endet nicht mit ready_for_query"},
		{"unbekannter Antwort-Typ", func(i *Interaction) { i.Responses[0].Type = "bogus" }, "Antwort-Typ \"bogus\""},
		{"Extended-Antwort in einfacher Anfrage", func(i *Interaction) { i.Responses[0].Type = ResponseParseComplete }, "Antwort-Typ \"parse_complete\""},
		{"mit Gruppen", func(i *Interaction) { i.Groups = extendedBeispiel().Groups }, "mit Gruppen"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := queryBeispiel()
			c.mutate(&i)
			err := i.Validate()
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("erwartet Fehler mit %q, erhalten %v", c.want, err)
			}
		})
	}
}
