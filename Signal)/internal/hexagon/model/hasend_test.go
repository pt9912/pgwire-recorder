package model_test

import (
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-25/Boundary — HasEnd ist wahr, wenn irgendeine Antwort irgendeiner
// Interaktion irgendeiner Session connection_closed ist, in Responses wie in Server einer
// Gruppe, an jeder Stelle (auch einer falschen) und in jeder Session; ohne Sessions, ohne
// Interaktion und ohne diese Antwort ist es falsch.
func TestHasEnd(t *testing.T) {
	ende := model.Response{Type: model.ResponseConnectionClosed}
	cc := model.Response{Type: model.ResponseCommandComplete, Tag: "x"}
	einfach := func(rs ...model.Response) model.Interaction { return model.Interaction{Responses: rs} }
	gruppen := func(server ...[]model.Response) model.Interaction {
		var i model.Interaction
		for _, s := range server {
			i.Groups = append(i.Groups, model.Group{Server: s})
		}
		return i
	}
	sess := func(is ...model.Interaction) model.Session { return model.Session{Interactions: is} }
	cases := []struct {
		name string
		rec  model.Recording
		want bool
	}{
		{"ohne Sessions", model.Recording{}, false},
		{"Session ohne Interaktion", model.Recording{Sessions: []model.Session{{}}}, false},
		{"ohne Ende", model.Recording{Sessions: []model.Session{sess(einfach(cc), gruppen([]model.Response{cc}))}}, false},
		{"Ende in Responses", model.Recording{Sessions: []model.Session{sess(einfach(cc, ende))}}, true},
		{"Ende in Responses, nicht letzte Nachricht", model.Recording{Sessions: []model.Session{sess(einfach(ende, cc))}}, true},
		{"Ende in Server einer Gruppe", model.Recording{Sessions: []model.Session{sess(gruppen([]model.Response{cc, ende}))}}, true},
		{"Ende in Server einer Gruppe, nicht letzte Nachricht", model.Recording{Sessions: []model.Session{sess(gruppen([]model.Response{ende, cc}))}}, true},
		{"Ende in Server einer nicht letzten Gruppe", model.Recording{Sessions: []model.Session{sess(gruppen([]model.Response{ende}, []model.Response{cc}))}}, true},
		{"Ende in einer nicht letzten Interaktion", model.Recording{Sessions: []model.Session{sess(einfach(ende), einfach(cc))}}, true},
		{"Ende in der letzten Interaktion", model.Recording{Sessions: []model.Session{sess(einfach(cc), einfach(ende))}}, true},
		{"Ende nur in der letzten Session", model.Recording{Sessions: []model.Session{sess(einfach(cc)), sess(einfach(ende))}}, true},
		{"Ende nur in der ersten Session", model.Recording{Sessions: []model.Session{sess(einfach(ende)), sess(einfach(cc))}}, true},
		{"Ende in der mittleren von drei Sessions", model.Recording{Sessions: []model.Session{sess(einfach(cc)), sess(einfach(ende)), sess(einfach(cc))}}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.rec.HasEnd(); got != c.want {
				t.Fatalf("HasEnd %v, erwartet %v", got, c.want)
			}
		})
	}
}
