package recording

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"go.yaml.in/yaml/v3"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Load liest die Aufzeichnung aus der Datei path (LH-FA-22.a *Lesen von `sqlite`*): Die
// Datei wird nur lesend geöffnet und in einer einzigen Lesetransaktion gelesen, nach der
// Formatkennung, der Version und dem Inhalt in dieser Reihenfolge geprüft. Der Inhalt
// wird zu dem Baum, den der Leser von yaml liest, und geht denselben Weg wie dort
// (Unmarshal): Schlüssel, Formen und Regeln der Aufzeichnung prüft nur dieser eine Weg.
// Die Fehlerklasse eines Fehlers der Bibliothek folgt aus seinem primären Ergebniscode
// (stufeCode), jeder Fehler des Inhalts ist PGR-E3003, eine unbekannte Version PGR-E3002.
func (*SQLite) Load(ctx context.Context, path string) (model.Recording, error) {
	text, err := leseAlsYAML(ctx, path, nil)
	if err != nil {
		return model.Recording{}, err
	}
	return Unmarshal(text)
}

// stufeCode liefert den Meldungscode für den Ergebniscode code der Bibliothek: ein
// abgeschnittenes Lesen (SQLITE_IOERR_SHORT_READ) ist PGR-E3003, sonst entscheidet der
// primäre Code (die unteren acht Bit): PERM, BUSY, LOCKED, NOMEM, READONLY, INTERRUPT,
// IOERR, FULL, CANTOPEN und PROTOCOL sind PGR-E3001, jeder andere PGR-E3003.
func stufeCode(code int) string {
	if code == sqlite3.SQLITE_IOERR_SHORT_READ {
		return model.CodeRecordingBroken
	}
	switch code & 0xff {
	case sqlite3.SQLITE_PERM, sqlite3.SQLITE_BUSY, sqlite3.SQLITE_LOCKED, sqlite3.SQLITE_NOMEM,
		sqlite3.SQLITE_READONLY, sqlite3.SQLITE_INTERRUPT, sqlite3.SQLITE_IOERR, sqlite3.SQLITE_FULL,
		sqlite3.SQLITE_CANTOPEN, sqlite3.SQLITE_PROTOCOL:
		return model.CodeRecordingIO
	}
	return model.CodeRecordingBroken
}

// bibliotheksfehler stuft einen Fehler err der Bibliothek beim Lesen von path ein: ein
// Fehler der Bibliothek nach stufeCode, jeder andere (etwa der Abbruch von ctx) ist
// PGR-E3001.
func bibliotheksfehler(path string, err error) error {
	var se *sqlite.Error
	if errors.As(err, &se) {
		return model.Errorf(stufeCode(se.Code()), err, "%s nicht lesbar", path)
	}
	return model.Errorf(model.CodeRecordingIO, err, "%s nicht lesbar", path)
}

// beschaedigt ist der Fehler PGR-E3003 für eine Datei, die als Aufzeichnung ungültig ist.
func beschaedigt(format string, args ...any) error {
	return model.Errorf(model.CodeRecordingBroken, nil, format, args...)
}

// lesung ist ein Lesen der Datei path in der Lesetransaktion tx. nach ist nil oder wird
// nach dem Lesen jeder Tabelle mit ihrem Namen gerufen, noch in der Transaktion.
type lesung struct {
	tx   *sql.Tx
	path string
	nach func(tx *sql.Tx, tabelle string)
}

// leseAlsYAML liest die Datei path in einer Lesetransaktion und liefert die Aufzeichnung
// als Text, den Unmarshal liest. nach ist der Eingriff von lesung.
func leseAlsYAML(ctx context.Context, path string, nach func(*sql.Tx, string)) ([]byte, error) {
	uri, err := uriVon(path)
	if err != nil {
		return nil, model.Errorf(model.CodeRecordingIO, err, "%s nicht lesbar", path)
	}
	db, err := sql.Open("sqlite", fmt.Sprintf("%s?mode=ro&_pragma=busy_timeout(%d)", uri, sperrfristMS))
	if err != nil {
		return nil, bibliotheksfehler(path, err)
	}
	defer func() { _ = db.Close() }()
	db.SetMaxOpenConns(1)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, bibliotheksfehler(path, err)
	}
	defer func() { _ = tx.Rollback() }()
	wurzel, err := (lesung{tx, path, nach}).liesAufzeichnung(ctx)
	if err != nil {
		return nil, err
	}
	text, err := yaml.Marshal(&yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{wurzel}})
	if err != nil {
		return nil, model.Errorf(model.CodeInternal, err, "Aufzeichnung aus %s nicht zu übertragen", path)
	}
	return text, nil
}

// Die Schlüssel von meta (SPEC-043).
const (
	metaFormat        = "format"
	metaVersion       = "version"
	metaEmptySessions = "empty_sessions"
)

// liesAufzeichnung liest meta und, wenn Formatkennung und Version stimmen, die Tabellen
// session, interaction und message; es liefert den Baum der Aufzeichnung.
func (l lesung) liesAufzeichnung(ctx context.Context) (*yaml.Node, error) {
	leere, err := l.liesMeta(ctx)
	if err != nil {
		return nil, err
	}
	sessions, err := l.liesSessions(ctx)
	if err != nil {
		return nil, err
	}
	wurzel := abbildung()
	wurzel.Content = append(wurzel.Content,
		schluessel("format"), textKnoten(model.Format),
		schluessel("version"), &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"})
	if leere {
		wurzel.Content = append(wurzel.Content, schluessel("empty_sessions"), &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "true"})
	}
	liste := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, s := range sessions {
		n, err := s.knoten()
		if err != nil {
			return nil, err
		}
		liste.Content = append(liste.Content, n)
	}
	wurzel.Content = append(wurzel.Content, schluessel("sessions"), liste)
	return wurzel, nil
}

// zeilen liefert die Zeilen der Anfrage an die Tabelle tabelle mit spalten Spalten, jede
// Spalte als Wert der Bibliothek (nil, int64, float64, string oder []byte).
func (l lesung) zeilen(ctx context.Context, tabelle, anfrage string, spalten int) ([][]any, error) {
	rows, err := l.tx.QueryContext(ctx, anfrage)
	if err != nil {
		return nil, bibliotheksfehler(l.path, err)
	}
	defer func() { _ = rows.Close() }()
	var out [][]any
	for rows.Next() {
		werte := make([]any, spalten)
		ziele := make([]any, spalten)
		for i := range werte {
			ziele[i] = &werte[i]
		}
		if err := rows.Scan(ziele...); err != nil {
			return nil, bibliotheksfehler(l.path, err)
		}
		out = append(out, werte)
	}
	if err := errors.Join(rows.Err(), rows.Close()); err != nil {
		return nil, bibliotheksfehler(l.path, err)
	}
	if l.nach != nil {
		l.nach(l.tx, tabelle)
	}
	return out, nil
}

// ganzzahl und zeichen liefern den Wert einer Spalte, wenn seine gespeicherte Art die
// des Spaltentyps ist (INTEGER oder TEXT).
func ganzzahl(w any) (int64, bool) { v, ok := w.(int64); return v, ok }
func zeichen(w any) (string, bool) { v, ok := w.(string); return v, ok }

// liesMeta liest meta (SPEC-043); es liefert, ob die Aufzeichnung leere Sessions trägt.
func (l lesung) liesMeta(ctx context.Context) (bool, error) {
	zeilen, err := l.zeilen(ctx, "meta", `SELECT "key", "value" FROM "meta"`, 2)
	if err != nil {
		return false, err
	}
	meta := map[string]string{}
	for _, z := range zeilen {
		k, kok := zeichen(z[0])
		v, vok := zeichen(z[1])
		if !kok || !vok {
			return false, beschaedigt("meta: Schlüssel und Wert sind kein TEXT")
		}
		if _, doppelt := meta[k]; doppelt {
			return false, beschaedigt("meta: Schlüssel %q doppelt", k)
		}
		meta[k] = v
	}
	return pruefeMeta(meta)
}

// pruefeMeta prüft meta in der Reihenfolge Formatkennung, Version, übrige Schlüssel und
// liefert, ob die Aufzeichnung leere Sessions trägt.
func pruefeMeta(meta map[string]string) (bool, error) {
	format, ok := meta[metaFormat]
	if !ok {
		return false, beschaedigt("meta ohne %s", metaFormat)
	}
	if format != model.Format {
		return false, beschaedigt("Formatkennung %q statt %q", format, model.Format)
	}
	version, ok := meta[metaVersion]
	if !ok {
		return false, beschaedigt("meta ohne %s", metaVersion)
	}
	if !kanonischeZahl.MatchString(version) {
		return false, beschaedigt("meta: version %q ist keine kanonische ganze Zahl", version)
	}
	if version != "1" {
		return false, model.Errorf(model.CodeRecordingVersion, nil, "Version %s wird nicht unterstützt", version)
	}
	for k, v := range meta {
		switch k {
		case metaFormat, metaVersion:
		case metaEmptySessions:
			if v != "true" {
				return false, beschaedigt("meta: %s trägt %q statt true", k, v)
			}
		default:
			return false, beschaedigt("meta: unbekannter Schlüssel %q", k)
		}
	}
	_, leere := meta[metaEmptySessions]
	return leere, nil
}

// Die Knoten, aus denen der Baum der Aufzeichnung entsteht.
func abbildung() *yaml.Node { return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"} }

func schluessel(name string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: name}
}

func zahl(n int64) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprint(n)}
}

// sessionBau ist eine Session der Datei mit ihren Interaktionen.
type sessionBau struct {
	id               int64
	startup          *yaml.Node
	serverParameters *yaml.Node
	interaktionen    []*interaktionBau
}

// interaktionBau ist eine Interaktion der Datei mit ihren Nachrichten.
type interaktionBau struct {
	sequence    int64
	art         string
	offsetMS    *int64
	sql         string
	nachrichten []nachrichtBau
}

type nachrichtBau struct {
	session, sequence, position int64
	richtung                    string
	gruppe                      *int64
	art                         string
	felder                      *yaml.Node
}

// liesSessions liest die Tabellen session, interaction und message und liefert die
// Sessions nach ihren Nummern, deren Interaktionen nach ihren Nummern und deren
// Nachrichten nach ihren Positionen: Die Reihenfolge der Zeilen in der Datei hat keine
// Bedeutung (LH-FA-22.a *Lesen von `sqlite`*).
func (l lesung) liesSessions(ctx context.Context) ([]*sessionBau, error) {
	sessions, err := l.liesSessionTabelle(ctx)
	if err != nil {
		return nil, err
	}
	nachNummer := map[int64]*sessionBau{}
	for _, s := range sessions {
		nachNummer[s.id] = s
	}
	interaktionen, err := l.liesInteraktionen(ctx, nachNummer)
	if err != nil {
		return nil, err
	}
	if err := l.liesNachrichten(ctx, interaktionen); err != nil {
		return nil, err
	}
	return sessions, nil
}

// liesStartup liest die Spalte text als JSON-Objekt Name → String, ohne Einträge als
// fehlend (nil): Das Modell von yaml trägt dann keine Einträge.
func liesStartup(spalte string, text any) (*yaml.Node, error) {
	s, ok := zeichen(text)
	if !ok {
		return nil, beschaedigt("session.%s ist kein TEXT", spalte)
	}
	n, err := jsonKnoten(s)
	if err != nil {
		return nil, beschaedigt("session.%s: %v", spalte, err)
	}
	if n.Kind != yaml.MappingNode {
		return nil, beschaedigt("session.%s ist kein JSON-Objekt", spalte)
	}
	for i := 1; i < len(n.Content); i += 2 {
		if t := n.Content[i].Tag; n.Content[i].Kind != yaml.ScalarNode || (t != "!!str" && t != "!!binary") {
			return nil, beschaedigt("session.%s: ein Wert ist kein String", spalte)
		}
	}
	if len(n.Content) == 0 {
		return nil, nil
	}
	return n, nil
}

func (l lesung) liesSessionTabelle(ctx context.Context) ([]*sessionBau, error) {
	zeilen, err := l.zeilen(ctx, "session", `SELECT "id", "startup", "server_parameters" FROM "session"`, 3)
	if err != nil {
		return nil, err
	}
	var sessions []*sessionBau
	for _, z := range zeilen {
		id, ok := ganzzahl(z[0])
		if !ok {
			return nil, beschaedigt("session.id ist keine ganze Zahl")
		}
		startup, err := liesStartup("startup", z[1])
		if err != nil {
			return nil, err
		}
		server, err := liesStartup("server_parameters", z[2])
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &sessionBau{id: id, startup: startup, serverParameters: server})
	}
	sort.SliceStable(sessions, func(i, j int) bool { return sessions[i].id < sessions[j].id })
	return sessions, nil
}

func (l lesung) liesInteraktionen(ctx context.Context, sessions map[int64]*sessionBau) (map[[2]int64]*interaktionBau, error) {
	zeilen, err := l.zeilen(ctx, "interaction", `SELECT "session_id", "sequence", "type", "offset_ms", "sql" FROM "interaction"`, 5)
	if err != nil {
		return nil, err
	}
	type zeile struct {
		session int64
		i       *interaktionBau
	}
	var alle []zeile
	for _, z := range zeilen {
		sid, i, err := interaktionVon(z)
		if err != nil {
			return nil, err
		}
		alle = append(alle, zeile{sid, i})
	}
	sort.SliceStable(alle, func(a, b int) bool {
		if alle[a].session != alle[b].session {
			return alle[a].session < alle[b].session
		}
		return alle[a].i.sequence < alle[b].i.sequence
	})
	nachSchluessel := map[[2]int64]*interaktionBau{}
	for _, z := range alle {
		s, ok := sessions[z.session]
		if !ok {
			return nil, beschaedigt("Interaktion %d ohne Session %d", z.i.sequence, z.session)
		}
		s.interaktionen = append(s.interaktionen, z.i)
		nachSchluessel[[2]int64{z.session, z.i.sequence}] = z.i
	}
	return nachSchluessel, nil
}

// interaktionVon prüft die Zeile z der Tabelle interaction gegen die Typen der Spalten und
// die Regeln von type und sql (SPEC-043); es liefert die Nummer der Session und die
// Interaktion.
func interaktionVon(z []any) (int64, *interaktionBau, error) {
	sid, sok := ganzzahl(z[0])
	seq, qok := ganzzahl(z[1])
	art, aok := zeichen(z[2])
	if !sok || !qok || !aok {
		return 0, nil, beschaedigt("interaction: session_id und sequence sind keine ganzen Zahlen oder type ist kein TEXT")
	}
	i := &interaktionBau{sequence: seq, art: art}
	if z[3] != nil {
		off, ok := ganzzahl(z[3])
		if !ok {
			return 0, nil, beschaedigt("Session %d, Interaktion %d: offset_ms weder NULL noch eine ganze Zahl", sid, seq)
		}
		i.offsetMS = &off
	}
	switch art {
	case "query":
		s, ok := zeichen(z[4])
		if !ok {
			return 0, nil, beschaedigt("Session %d, Interaktion %d: sql einer einfachen Anfrage ist kein TEXT", sid, seq)
		}
		i.sql = s
	case "extended":
		if z[4] != nil {
			return 0, nil, beschaedigt("Session %d, Interaktion %d: sql einer Extended-Interaktion ist nicht NULL", sid, seq)
		}
	default:
		return 0, nil, beschaedigt("Session %d, Interaktion %d: type %q unbekannt", sid, seq, art)
	}
	return sid, i, nil
}

func (l lesung) liesNachrichten(ctx context.Context, interaktionen map[[2]int64]*interaktionBau) error {
	zeilen, err := l.zeilen(ctx, "message", `SELECT "session_id", "sequence", "position", "direction", "group_no", "kind", "fields" FROM "message"`, 7)
	if err != nil {
		return err
	}
	var alle []nachrichtBau
	for _, z := range zeilen {
		n, err := nachrichtVon(z)
		if err != nil {
			return err
		}
		alle = append(alle, n)
	}
	sort.SliceStable(alle, func(a, b int) bool {
		x, y := alle[a], alle[b]
		if x.session != y.session {
			return x.session < y.session
		}
		if x.sequence != y.sequence {
			return x.sequence < y.sequence
		}
		return x.position < y.position
	})
	for _, n := range alle {
		i, ok := interaktionen[[2]int64{n.session, n.sequence}]
		if !ok {
			return beschaedigt("Nachricht %d ohne Interaktion %d in Session %d", n.position, n.sequence, n.session)
		}
		if n.position != int64(len(i.nachrichten))+1 {
			return beschaedigt("Session %d, Interaktion %d: Nachricht %d trägt die Position %d", n.session, n.sequence, len(i.nachrichten)+1, n.position)
		}
		i.nachrichten = append(i.nachrichten, n)
	}
	return nil
}

// nachrichtVon prüft die Zeile z der Tabelle message gegen die Typen der Spalten und liest
// fields als Objekt (SPEC-043).
func nachrichtVon(z []any) (nachrichtBau, error) {
	sid, sok := ganzzahl(z[0])
	seq, qok := ganzzahl(z[1])
	pos, pok := ganzzahl(z[2])
	richtung, rok := zeichen(z[3])
	art, aok := zeichen(z[5])
	text, fok := zeichen(z[6])
	if !sok || !qok || !pok || !rok || !aok || !fok {
		return nachrichtBau{}, beschaedigt("message: eine Spalte hat nicht die Art ihres Typs")
	}
	n := nachrichtBau{session: sid, sequence: seq, position: pos, richtung: richtung, art: art}
	if z[4] != nil {
		g, ok := ganzzahl(z[4])
		if !ok {
			return nachrichtBau{}, beschaedigt("Session %d, Interaktion %d, Nachricht %d: group_no weder NULL noch eine ganze Zahl", sid, seq, pos)
		}
		n.gruppe = &g
	}
	if richtung != richtungClient && richtung != richtungServer {
		return nachrichtBau{}, beschaedigt("Session %d, Interaktion %d, Nachricht %d: direction %q unbekannt", sid, seq, pos, richtung)
	}
	felder, err := jsonKnoten(text)
	if err != nil {
		return nachrichtBau{}, beschaedigt("Session %d, Interaktion %d, Nachricht %d: fields: %v", sid, seq, pos, err)
	}
	if felder.Kind != yaml.MappingNode {
		return nachrichtBau{}, beschaedigt("Session %d, Interaktion %d, Nachricht %d: fields ist kein JSON-Objekt", sid, seq, pos)
	}
	n.felder = felder
	return n, nil
}
