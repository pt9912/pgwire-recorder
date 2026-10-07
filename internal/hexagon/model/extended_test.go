package model_test

import (
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// extendedBeispiel ist die Interaktion aus SPEC-041 mit einer vorangestellten
// Flush-Gruppe: parse/describe/flush, danach bind/execute/sync.
func extendedBeispiel() model.Interaction {
	return model.Interaction{
		Sequence: 1,
		Request:  model.Request{Type: model.RequestExtended},
		Groups: []model.Group{
			{
				Client: []model.ClientMessage{
					{Type: model.ClientParse, Statement: "s1", SQL: "SELECT name FROM users WHERE id = $1", ParamTypes: []uint32{23}},
					{Type: model.ClientDescribe, Target: model.TargetStatement, Name: "s1"},
					{Type: model.ClientFlush},
				},
				Server: []model.Response{
					{Type: model.ResponseParseComplete},
					{Type: model.ResponseParameterDescription, ParamTypes: []uint32{23}},
					{Type: model.ResponseRowDescription, Columns: []model.Column{{Name: "name", TypeOID: 25, TypeSize: -1, TypeModifier: -1}}},
				},
			},
			{
				Client: []model.ClientMessage{
					{Type: model.ClientBind, Statement: "s1", ParamFormats: []int16{0}, Params: []model.Value{{Bytes: []byte("1")}}, ResultFormats: []int16{0}},
					{Type: model.ClientExecute, MaxRows: 0},
					{Type: model.ClientSync},
				},
				Server: []model.Response{
					{Type: model.ResponseBindComplete},
					{Type: model.ResponseDataRow, Values: []model.Value{{Bytes: []byte("alice")}}},
					{Type: model.ResponseCommandComplete, Tag: "SELECT 1"},
					{Type: model.ResponseReadyForQuery, TxStatus: "I"},
				},
			},
		},
	}
}

func queryBeispiel() model.Interaction {
	return model.Interaction{
		Sequence:  1,
		Request:   model.Request{Type: model.RequestQuery, SQL: "SELECT 1"},
		Responses: []model.Response{{Type: model.ResponseCommandComplete, Tag: "SELECT 1"}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}},
	}
}

// Gültige Formen: einfache Anfrage, auch mit jedem Antworttyp des Simple Query
// Protocol (LH-FA-05.a), Extended mit Flush- und Sync-Gruppe, Extended
// mit nur einer Sync-Gruppe, eine Flush-Gruppe ohne Server-Nachricht, jede
// Client- und Server-Nachricht aus LH-FA-18.a einschließlich parameter_status.
func TestValidateGueltig(t *testing.T) {
	nurSync := model.Interaction{Request: model.Request{Type: model.RequestExtended}, Groups: []model.Group{{
		Client: []model.ClientMessage{{Type: model.ClientSync}},
		Server: []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}},
	}}}
	alle := model.Interaction{Request: model.Request{Type: model.RequestExtended}, Groups: []model.Group{
		{
			Client: []model.ClientMessage{{Type: model.ClientClose, Target: model.TargetPortal}, {Type: model.ClientFlush}},
		},
		{
			Client: []model.ClientMessage{
				{Type: model.ClientParse}, {Type: model.ClientBind}, {Type: model.ClientDescribe, Target: model.TargetPortal},
				{Type: model.ClientExecute, MaxRows: 1}, {Type: model.ClientClose, Target: model.TargetStatement}, {Type: model.ClientSync},
			},
			Server: []model.Response{
				{Type: model.ResponseCloseComplete}, {Type: model.ResponseParseComplete}, {Type: model.ResponseBindComplete},
				{Type: model.ResponseNoData}, {Type: model.ResponseParameterDescription}, {Type: model.ResponseRowDescription},
				{Type: model.ResponseDataRow}, {Type: model.ResponsePortalSuspended}, {Type: model.ResponseCommandComplete},
				{Type: model.ResponseEmptyQueryResponse}, {Type: model.ResponseNoticeResponse}, {Type: model.ResponseParameterStatus},
				{Type: model.ResponseErrorResponse}, {Type: model.ResponseCloseComplete}, {Type: model.ResponseReadyForQuery},
			},
		},
	}}
	queryAlle := queryBeispiel()
	queryAlle.Responses = []model.Response{
		{Type: model.ResponseRowDescription}, {Type: model.ResponseDataRow}, {Type: model.ResponseCommandComplete},
		{Type: model.ResponseEmptyQueryResponse}, {Type: model.ResponseErrorResponse}, {Type: model.ResponseNoticeResponse},
		{Type: model.ResponseParameterStatus}, {Type: model.ResponseReadyForQuery},
	}
	for name, i := range map[string]model.Interaction{
		"query":                       queryBeispiel(),
		"query mit jedem Antworttyp":  queryAlle,
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
		mutate func(*model.Interaction)
		want   string
	}{
		{"unbekannter Anfrage-Typ", func(i *model.Interaction) { i.Request.Type = "bogus" }, "Anfrage-Typ"},
		{"leerer Anfrage-Typ", func(i *model.Interaction) { i.Request.Type = "" }, "Anfrage-Typ"},
		{"extended mit SQL", func(i *model.Interaction) { i.Request.SQL = "SELECT 1" }, "außerhalb der Gruppen"},
		{"extended mit Antworten", func(i *model.Interaction) {
			i.Responses = []model.Response{{Type: model.ResponseReadyForQuery}}
		}, "außerhalb der Gruppen"},
		{"extended ohne Gruppe", func(i *model.Interaction) { i.Groups = nil }, "ohne Gruppe"},
		{"Gruppe ohne Client-Nachricht", func(i *model.Interaction) { i.Groups[0].Client = nil }, "ohne Client-Nachricht"},
		{"unbekannte Client-Nachricht", func(i *model.Interaction) { i.Groups[0].Client[0].Type = "copy_data" }, "Client-Nachricht \"copy_data\" unbekannt"},
		{"flush mitten in der Gruppe", func(i *model.Interaction) { i.Groups[0].Client[1] = model.ClientMessage{Type: model.ClientFlush} }, "nur als letzte Nachricht"},
		{"sync mitten in der Gruppe", func(i *model.Interaction) { i.Groups[1].Client[0] = model.ClientMessage{Type: model.ClientSync} }, "nur als letzte Nachricht"},
		{"Gruppe ohne flush oder sync am Ende", func(i *model.Interaction) {
			i.Groups[0].Client = i.Groups[0].Client[:2]
		}, "nur als letzte Nachricht"},
		{"vordere Gruppe endet mit sync", func(i *model.Interaction) { i.Groups[0].Client[2].Type = model.ClientSync }, "Gruppe 1: sync steht"},
		{"letzte Gruppe endet mit flush (abgeschnitten)", func(i *model.Interaction) { i.Groups[1].Client[2].Type = model.ClientFlush }, "Gruppe 2: sync steht"},
		{"describe ohne Zielart", func(i *model.Interaction) { i.Groups[0].Client[1].Target = "" }, "Zielart"},
		{"close mit falscher Zielart", func(i *model.Interaction) {
			i.Groups[0].Client[1] = model.ClientMessage{Type: model.ClientClose, Target: "tabelle"}
		}, "Zielart"},
		{"unbekannte Server-Nachricht", func(i *model.Interaction) { i.Groups[1].Server[1].Type = "copy_out_response" }, "Server-Nachricht \"copy_out_response\" unbekannt"},
		{"ready_for_query in der Flush-Gruppe", func(i *model.Interaction) {
			i.Groups[0].Server = append(i.Groups[0].Server, model.Response{Type: model.ResponseReadyForQuery})
		}, "Gruppe 1: Server-Nachricht 4: ready_for_query"},
		{"ready_for_query vor dem Ende der letzten Gruppe", func(i *model.Interaction) {
			i.Groups[1].Server[2].Type = model.ResponseReadyForQuery
		}, "Gruppe 2: Server-Nachricht 3: ready_for_query"},
		{"letzte Gruppe ohne ready_for_query (abgeschnitten)", func(i *model.Interaction) {
			i.Groups[1].Server = i.Groups[1].Server[:3]
		}, "endet nicht mit ready_for_query"},
		{"letzte Gruppe ohne Server-Nachricht", func(i *model.Interaction) { i.Groups[1].Server = nil }, "endet nicht mit ready_for_query"},
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
		mutate func(*model.Interaction)
		want   string
	}{
		{"ohne Antwort", func(i *model.Interaction) { i.Responses = nil }, "endet nicht mit ready_for_query"},
		{"ohne ready_for_query am Ende", func(i *model.Interaction) { i.Responses = i.Responses[:1] }, "endet nicht mit ready_for_query"},
		{"unbekannter Antwort-Typ", func(i *model.Interaction) { i.Responses[0].Type = "bogus" }, "Antwort-Typ \"bogus\""},
		{"Extended-Antwort in einfacher Anfrage", func(i *model.Interaction) { i.Responses[0].Type = model.ResponseParseComplete }, "Antwort-Typ \"parse_complete\""},
		{"mit Gruppen", func(i *model.Interaction) { i.Groups = extendedBeispiel().Groups }, "mit Gruppen"},
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

// Verletzt eine Gruppe zwei Formregeln zugleich, meldet Validate genau eine. Der
// Test hält fest, welche der Bestand meldet; eine Zusage ist das nicht, und wer
// die Reihenfolge bewusst ändert, passt ihn an. Im Bestand gehen die
// Client-Nachrichten der Reihe nach vor, je Nachricht erst der Typ, dann die
// Stellung von flush und sync, dann die Zielart; danach sync oder flush am Ende
// der Gruppe; danach die Server-Nachrichten der Reihe nach, je Nachricht erst der
// Typ, dann die Stellung von ready_for_query; zuletzt ready_for_query am Ende der
// letzten Gruppe. Erwartet ist der ganze Fehlertext.
func TestValidateFehlerReihenfolge(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*model.Interaction)
		want   string
	}{
		{"ohne Client-Nachricht vor einer unbekannten Server-Nachricht", func(i *model.Interaction) {
			i.Groups[0].Client = nil
			i.Groups[0].Server[0].Type = "copy_out_response"
		}, "Gruppe 1: ohne Client-Nachricht"},
		{"Typ vor Stellung in derselben Nachricht", func(i *model.Interaction) {
			i.Groups[0].Client[2].Type = "copy_data"
		}, "Gruppe 1: Client-Nachricht \"copy_data\" unbekannt"},
		{"Stellung vor Zielart in derselben Nachricht", func(i *model.Interaction) {
			i.Groups[0].Client = []model.ClientMessage{{Type: model.ClientParse}, {Type: model.ClientDescribe}}
		}, "Gruppe 1: Client-Nachricht 2: flush oder sync nur als letzte Nachricht der Gruppe"},
		{"Zielart einer früheren vor Typ einer späteren Nachricht", func(i *model.Interaction) {
			i.Groups[0].Client[1].Target = ""
			i.Groups[0].Client[2].Type = "copy_data"
		}, "Gruppe 1: Client-Nachricht 2: Zielart \"\" statt statement oder portal"},
		{"Stellung einer früheren vor Typ einer späteren Nachricht", func(i *model.Interaction) {
			i.Groups[0].Client[0] = model.ClientMessage{Type: model.ClientFlush}
			i.Groups[0].Client[1].Type = "copy_data"
		}, "Gruppe 1: Client-Nachricht 1: flush oder sync nur als letzte Nachricht der Gruppe"},
		{"Zielart vor sync am Ende einer vorderen Gruppe", func(i *model.Interaction) {
			i.Groups[0].Client[1].Target = ""
			i.Groups[0].Client[2].Type = model.ClientSync
		}, "Gruppe 1: Client-Nachricht 2: Zielart \"\" statt statement oder portal"},
		{"sync am Ende einer vorderen Gruppe vor einer unbekannten Server-Nachricht", func(i *model.Interaction) {
			i.Groups[0].Client[2].Type = model.ClientSync
			i.Groups[0].Server[0].Type = "copy_out_response"
		}, "Gruppe 1: sync steht genau am Ende der letzten Gruppe, flush am Ende jeder anderen"},
		{"flush am Ende der letzten Gruppe vor fehlendem ready_for_query", func(i *model.Interaction) {
			i.Groups[1].Client[2].Type = model.ClientFlush
			i.Groups[1].Server = i.Groups[1].Server[:3]
		}, "Gruppe 2: sync steht genau am Ende der letzten Gruppe, flush am Ende jeder anderen"},
		{"Client-Nachricht vor Server-Nachricht", func(i *model.Interaction) {
			i.Groups[1].Client[0].Type = "copy_data"
			i.Groups[1].Server[1].Type = "copy_out_response"
		}, "Gruppe 2: Client-Nachricht \"copy_data\" unbekannt"},
		{"Stellung einer früheren vor Typ einer späteren Server-Nachricht", func(i *model.Interaction) {
			i.Groups[1].Server[0].Type = model.ResponseReadyForQuery
			i.Groups[1].Server[1].Type = "copy_out_response"
		}, "Gruppe 2: Server-Nachricht 1: ready_for_query nur als letzte Nachricht der letzten Gruppe"},
		{"Typ einer früheren vor Stellung einer späteren Server-Nachricht", func(i *model.Interaction) {
			i.Groups[0].Server[0].Type = "copy_out_response"
			i.Groups[0].Server[1].Type = model.ResponseReadyForQuery
		}, "Gruppe 1: Server-Nachricht \"copy_out_response\" unbekannt in einer Extended-Interaktion"},
		{"unbekannte Server-Nachricht vor fehlendem ready_for_query", func(i *model.Interaction) {
			i.Groups[1].Server = i.Groups[1].Server[:3]
			i.Groups[1].Server[1].Type = "copy_out_response"
		}, "Gruppe 2: Server-Nachricht \"copy_out_response\" unbekannt in einer Extended-Interaktion"},
		{"Stellung von ready_for_query vor fehlendem ready_for_query am Ende", func(i *model.Interaction) {
			i.Groups[1].Server = []model.Response{{Type: model.ResponseReadyForQuery}, {Type: model.ResponseCommandComplete}}
		}, "Gruppe 2: Server-Nachricht 1: ready_for_query nur als letzte Nachricht der letzten Gruppe"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := extendedBeispiel()
			c.mutate(&i)
			err := i.Validate()
			if err == nil || err.Error() != c.want {
				t.Fatalf("erwartet Fehler %q, erhalten %v", c.want, err)
			}
		})
	}
}
