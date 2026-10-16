package recording

import (
	"go.yaml.in/yaml/v3"
)

// knoten liefert den Baum der Session, wie ihn die Datei von yaml trägt (SPEC-002):
// Startup und Server-Parameter ohne Einträge fehlen.
func (s *sessionBau) knoten() (*yaml.Node, error) {
	n := abbildung()
	n.Content = append(n.Content, schluessel("id"), zahl(s.id))
	if s.startup != nil {
		n.Content = append(n.Content, schluessel("startup"), s.startup)
	}
	if s.serverParameters != nil {
		n.Content = append(n.Content, schluessel("server_parameters"), s.serverParameters)
	}
	liste := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, i := range s.interaktionen {
		in, err := i.knoten(s.id)
		if err != nil {
			return nil, err
		}
		liste.Content = append(liste.Content, in)
	}
	n.Content = append(n.Content, schluessel("interactions"), liste)
	return n, nil
}

// knoten liefert den Baum der Interaktion der Session session. Eine einfache Anfrage
// trägt request und die Server-Nachrichten als responses; eine Extended-Interaktion
// trägt type und ihre Nachrichten als groups (SPEC-041, SPEC-043).
func (i *interaktionBau) knoten(session int64) (*yaml.Node, error) {
	n := abbildung()
	n.Content = append(n.Content, schluessel("sequence"), zahl(i.sequence))
	if i.offsetMS != nil {
		n.Content = append(n.Content, schluessel("offset_ms"), zahl(*i.offsetMS))
	}
	if i.art == "extended" {
		gruppen, err := i.gruppen(session)
		if err != nil {
			return nil, err
		}
		n.Content = append(n.Content, schluessel("type"), schluessel("extended"))
		if len(gruppen.Content) > 0 {
			n.Content = append(n.Content, schluessel("groups"), gruppen)
		}
		return n, nil
	}
	anfrage := abbildung()
	anfrage.Content = append(anfrage.Content, schluessel("type"), schluessel("query"), schluessel("sql"), textKnoten(i.sql))
	n.Content = append(n.Content, schluessel("request"), anfrage)
	antworten := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, m := range i.nachrichten {
		if m.richtung != richtungServer || m.gruppe != nil {
			return nil, beschaedigt("Session %d, Interaktion %d, Nachricht %d: eine einfache Anfrage trägt nur Server-Nachrichten ohne group_no", session, i.sequence, m.position)
		}
		antworten.Content = append(antworten.Content, m.knoten())
	}
	if len(antworten.Content) > 0 {
		n.Content = append(n.Content, schluessel("responses"), antworten)
	}
	return n, nil
}

// gruppen liefert die Gruppen einer Extended-Interaktion: group_no zählt ab 1 lückenlos
// in der Reihenfolge der Nachrichten, in jeder Gruppe stehen die Client-Nachrichten vor
// den Server-Nachrichten (SPEC-043).
func (i *interaktionBau) gruppen(session int64) (*yaml.Node, error) {
	liste := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	var client, server *yaml.Node
	var aktuell int64
	for _, m := range i.nachrichten {
		if m.gruppe == nil {
			return nil, beschaedigt("Session %d, Interaktion %d, Nachricht %d: group_no einer Extended-Interaktion ist NULL", session, i.sequence, m.position)
		}
		switch {
		case *m.gruppe == aktuell && aktuell != 0:
		case *m.gruppe == aktuell+1:
			aktuell++
			client = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			server = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			g := abbildung()
			g.Content = append(g.Content, schluessel("client"), client, schluessel("server"), server)
			liste.Content = append(liste.Content, g)
		default:
			return nil, beschaedigt("Session %d, Interaktion %d, Nachricht %d: group_no %d folgt auf %d", session, i.sequence, m.position, *m.gruppe, aktuell)
		}
		switch {
		case m.richtung == richtungClient && len(server.Content) > 0:
			return nil, beschaedigt("Session %d, Interaktion %d, Nachricht %d: Client-Nachricht nach einer Server-Nachricht der Gruppe %d", session, i.sequence, m.position, aktuell)
		case m.richtung == richtungClient:
			client.Content = append(client.Content, m.knoten())
		default:
			server.Content = append(server.Content, m.knoten())
		}
	}
	return liste, nil
}

// knoten liefert den Baum der Nachricht: type mit der Art, danach die Felder.
func (m nachrichtBau) knoten() *yaml.Node {
	n := abbildung()
	n.Content = append(n.Content, schluessel("type"), textKnoten(m.art))
	n.Content = append(n.Content, m.felder.Content...)
	return n
}
