package recording

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// YAML speichert Aufzeichnungen im Format yaml (SPEC-001, SPEC-002).
type YAML struct{}

var _ driven.RecordingRepository = YAML{}

// Prepare prüft den Zielpfad (LH-FA-07.a).
func (YAML) Prepare(_ context.Context, path string, replace bool) error {
	_, err := os.Stat(path)
	switch {
	case err == nil && !replace:
		return model.Errorf(model.CodeOutputExists, nil, "%s existiert bereits; --force ersetzt die Datei", path)
	case err == nil, errors.Is(err, fs.ErrNotExist):
	default:
		return model.Errorf(model.CodeRecordingIO, err, "%s nicht prüfbar", path)
	}
	// Das Zielverzeichnis muss beim Start beschreibbar sein, nicht erst nach der
	// ersten Session.
	probe, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*.probe")
	if err != nil {
		return model.Errorf(model.CodeRecordingIO, err, "Verzeichnis von %s nicht beschreibbar", path)
	}
	name := probe.Name()
	return errors.Join(closeErr(probe.Close()), removeErr(os.Remove(name)))
}

func closeErr(err error) error {
	if err != nil {
		return model.Errorf(model.CodeRecordingIO, err, "Probedatei nicht zu schließen")
	}
	return nil
}

func removeErr(err error) error {
	if err != nil {
		return model.Errorf(model.CodeRecordingIO, err, "Probedatei nicht zu entfernen")
	}
	return nil
}

// Write schreibt in eine temporäre Datei im Zielverzeichnis und benennt sie
// danach um; die Zieldatei ist damit vollständig oder unverändert. Eine neue
// Datei erhält die Rechte nach der umask des Prozesses (0666 vor der umask),
// eine ersetzte behält ihre bisherigen (SPEC-033).
func (YAML) Write(_ context.Context, path string, rec model.Recording) error {
	data, err := Marshal(rec)
	if err != nil {
		return err
	}
	tmp, err := neueTempDatei(path)
	if err != nil {
		return model.Errorf(model.CodeRecordingIO, err, "temporäre Datei für %s nicht anzulegen", path)
	}
	var merr error
	if info, err := os.Stat(path); err == nil {
		merr = tmp.Chmod(info.Mode().Perm())
	}
	_, werr := tmp.Write(data)
	serr := tmp.Sync()
	cerr := tmp.Close()
	if err := errors.Join(werr, merr, serr, cerr); err != nil {
		return model.Errorf(model.CodeRecordingIO, err, "temporäre Datei für %s nicht zu schreiben", path)
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return model.Errorf(model.CodeRecordingIO, err, "%s nicht zu schreiben", path)
	}
	return nil
}

// neueTempDatei legt im Verzeichnis von path eine neue Datei an; die Rechte
// 0666 schränkt die umask ein (anders als os.CreateTemp mit fest 0600).
func neueTempDatei(path string) (*os.File, error) {
	var zufall [8]byte
	for i := 0; i < 10; i++ {
		if _, err := rand.Read(zufall[:]); err != nil {
			return nil, err
		}
		name := filepath.Join(filepath.Dir(path), "."+filepath.Base(path)+"."+hex.EncodeToString(zufall[:])+".tmp")
		f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
		if errors.Is(err, fs.ErrExist) {
			continue
		}
		return f, err
	}
	return nil, errors.New("kein freier Name für die temporäre Datei")
}

// Load liest eine Aufzeichnung.
func (YAML) Load(_ context.Context, path string) (model.Recording, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Recording{}, model.Errorf(model.CodeRecordingIO, err, "%s nicht lesbar", path)
	}
	return Unmarshal(data)
}

// Marshal liefert die YAML-Darstellung; gleiche Aufzeichnungen ergeben gleiche
// Bytes (SPEC-004).
func Marshal(rec model.Recording) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(toDTO(rec)); err != nil {
		return nil, model.Errorf(model.CodeInternal, err, "Aufzeichnung nicht zu serialisieren")
	}
	if err := enc.Close(); err != nil {
		return nil, model.Errorf(model.CodeInternal, err, "Aufzeichnung nicht zu serialisieren")
	}
	return buf.Bytes(), nil
}

// Unmarshal liest die YAML-Darstellung. Eine fremde Formatkennung oder ein
// unlesbarer Inhalt ist PGR-E3003, eine unbekannte Version PGR-E3002.
func Unmarshal(data []byte) (model.Recording, error) {
	var d recordingDTO
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&d); err != nil {
		return model.Recording{}, model.Errorf(model.CodeRecordingBroken, err, "Aufzeichnung beschädigt")
	}
	if d.Format != model.Format {
		return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Formatkennung %q statt %q", d.Format, model.Format)
	}
	if d.Version != model.Version {
		return model.Recording{}, model.Errorf(model.CodeRecordingVersion, nil, "Version %d wird nicht unterstützt", d.Version)
	}
	return fromDTO(d)
}

type recordingDTO struct {
	Format        string        `yaml:"format"`
	Version       int           `yaml:"version"`
	EmptySessions bool          `yaml:"empty_sessions,omitempty"`
	Sessions      *[]sessionDTO `yaml:"sessions"`
}

type sessionDTO struct {
	ID               int               `yaml:"id"`
	Startup          map[string]string `yaml:"startup,omitempty"`
	ServerParameters map[string]string `yaml:"server_parameters,omitempty"`
	Interactions     []interactionDTO  `yaml:"interactions"`
}

// interactionDTO trägt beide Arten einer Interaktion: eine einfache Anfrage mit
// request und responses (SPEC-002), eine Extended-Interaktion mit
// `type: extended` und groups (SPEC-041). Ein anderer Wert von type ist
// beschädigt (SPEC-001).
type interactionDTO struct {
	Sequence  int           `yaml:"sequence"`
	OffsetMS  *int64        `yaml:"offset_ms,omitempty"`
	Type      string        `yaml:"type,omitempty"`
	Request   *requestDTO   `yaml:"request,omitempty"`
	Responses []responseDTO `yaml:"responses,omitempty"`
	Groups    []groupDTO    `yaml:"groups,omitempty"`
}

// interactionExtended ist der Wert von type einer Extended-Interaktion.
const interactionExtended = "extended"

type groupDTO struct {
	Client []clientMessageDTO `yaml:"client"`
	Server []responseDTO      `yaml:"server"`
}

// clientFelder nennt je Client-Nachricht die Schlüssel neben type in der
// Reihenfolge der Ausgabe (SPEC-041). Der Schreiber schreibt genau diese, auch
// leer; der Leser lehnt jeden anderen Schlüssel und jeden anderen Typ ab.
var clientFelder = map[string][]string{
	"parse":    {"statement", "sql", "param_types"},
	"bind":     {"portal", "statement", "param_formats", "params", "result_formats"},
	"describe": {"target", "name"},
	"execute":  {"portal", "max_rows"},
	"close":    {"target", "name"},
	"flush":    nil,
	"sync":     nil,
}

// clientMessageDTO ist eine Client-Nachricht; welche Felder zu einem Typ
// gehören, legt clientFelder fest.
type clientMessageDTO struct {
	Type          string     `yaml:"type"`
	Statement     string     `yaml:"statement"`
	Portal        string     `yaml:"portal"`
	SQL           string     `yaml:"sql"`
	ParamTypes    []uint32   `yaml:"param_types,flow"`
	ParamFormats  []int16    `yaml:"param_formats,flow"`
	Params        []valueDTO `yaml:"params"`
	ResultFormats []int16    `yaml:"result_formats,flow"`
	Target        string     `yaml:"target"`
	Name          string     `yaml:"name"`
	MaxRows       uint32     `yaml:"max_rows"`
}

// clientMessageRoh ist clientMessageDTO ohne die eigenen YAML-Methoden.
type clientMessageRoh clientMessageDTO

// MarshalYAML schreibt type und danach die Schlüssel aus clientFelder.
func (m clientMessageDTO) MarshalYAML() (any, error) {
	felder, ok := clientFelder[m.Type]
	if !ok {
		return nil, fmt.Errorf("Client-Nachricht %q unbekannt", m.Type)
	}
	var alle yaml.Node
	if err := alle.Encode(clientMessageRoh(m)); err != nil {
		return nil, err
	}
	werte := map[string]*yaml.Node{}
	for i := 0; i+1 < len(alle.Content); i += 2 {
		werte[alle.Content[i].Value] = alle.Content[i+1]
	}
	out := &yaml.Node{Kind: yaml.MappingNode}
	for _, f := range append([]string{"type"}, felder...) {
		out.Content = append(out.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: f}, werte[f])
	}
	return out, nil
}

// UnmarshalYAML liest eine Client-Nachricht; ein unbekannter Typ oder ein
// Schlüssel, der nicht zu ihm gehört, ist ein Fehler.
func (m *clientMessageDTO) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("Client-Nachricht ist keine Abbildung")
	}
	typ := ""
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == "type" {
			typ = node.Content[i+1].Value
		}
	}
	felder, ok := clientFelder[typ]
	if !ok {
		return fmt.Errorf("Client-Nachricht %q unbekannt", typ)
	}
	erlaubt := map[string]bool{"type": true}
	for _, f := range felder {
		erlaubt[f] = true
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if k := node.Content[i].Value; !erlaubt[k] {
			return fmt.Errorf("Schlüssel %q gehört nicht zu %s", k, typ)
		}
	}
	return node.Decode((*clientMessageRoh)(m))
}

type requestDTO struct {
	Type string `yaml:"type"`
	SQL  string `yaml:"sql"`
}

type responseDTO struct {
	Type     string            `yaml:"type"`
	Fields   []columnDTO       `yaml:"fields,omitempty"`
	Values   []valueDTO        `yaml:"values,omitempty"`
	Tag      string            `yaml:"tag,omitempty"`
	Notice   map[string]string `yaml:"notice,omitempty"`
	Name     string            `yaml:"name,omitempty"`
	Value    string            `yaml:"value,omitempty"`
	TxStatus string            `yaml:"tx_status,omitempty"`
	// ParamTypes trägt die Typ-OIDs einer parameter_description.
	ParamTypes []uint32 `yaml:"param_types,omitempty,flow"`
}

type columnDTO struct {
	Name         string `yaml:"name"`
	TableOID     uint32 `yaml:"table_oid"`
	ColumnNumber uint16 `yaml:"column_number"`
	TypeOID      uint32 `yaml:"type_oid"`
	TypeSize     int16  `yaml:"type_size"`
	TypeModifier int32  `yaml:"type_modifier"`
	Format       int16  `yaml:"format"`
}

// valueDTO trägt einen Wert als Text, als NULL oder als Base64, wenn die Bytes
// kein gültiges UTF-8 sind (SPEC-003).
type valueDTO struct {
	Text   *string `yaml:"text,omitempty"`
	Null   bool    `yaml:"null,omitempty"`
	Base64 *string `yaml:"base64,omitempty"`
}

// UnmarshalYAML liest einen Wert. Der Schlüssel `null` darf ungequotet stehen
// (`- null: true`, SPEC-041); YAML liest ihn dann als Null-Schlüssel.
func (v *valueDTO) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("Wert ist keine Abbildung")
	}
	gesehen := map[string]bool{}
	for i := 0; i+1 < len(node.Content); i += 2 {
		key, val := node.Content[i], node.Content[i+1]
		name := key.Value
		if key.Tag == "!!null" {
			name = "null"
		}
		if gesehen[name] {
			return fmt.Errorf("Schlüssel %q doppelt in einem Wert", name)
		}
		gesehen[name] = true
		switch name {
		case "text":
			var s string
			if err := val.Decode(&s); err != nil {
				return err
			}
			v.Text = &s
		case "base64":
			var s string
			if err := val.Decode(&s); err != nil {
				return err
			}
			v.Base64 = &s
		case "null":
			if err := val.Decode(&v.Null); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unbekannter Schlüssel %q in einem Wert", name)
		}
	}
	return nil
}

func toDTO(rec model.Recording) recordingDTO {
	sessions := []sessionDTO{}
	d := recordingDTO{Format: rec.Format, Version: rec.Version, EmptySessions: rec.EmptySessions, Sessions: &sessions}
	for _, s := range rec.Sessions {
		sd := sessionDTO{ID: s.ID, Startup: s.Startup, ServerParameters: s.ServerParameters, Interactions: []interactionDTO{}}
		for _, i := range s.Interactions {
			sd.Interactions = append(sd.Interactions, interactionToDTO(i))
		}
		sessions = append(sessions, sd)
	}
	return d
}

// interactionToDTO schreibt eine Extended-Interaktion mit type und groups, jede
// andere mit request und responses.
func interactionToDTO(i model.Interaction) interactionDTO {
	d := interactionDTO{Sequence: i.Sequence, OffsetMS: i.OffsetMS}
	if i.Request.Type == model.RequestExtended {
		d.Type = interactionExtended
		for _, g := range i.Groups {
			var gd groupDTO
			for _, m := range g.Client {
				gd.Client = append(gd.Client, clientMessageToDTO(m))
			}
			for _, r := range g.Server {
				gd.Server = append(gd.Server, responseToDTO(r))
			}
			d.Groups = append(d.Groups, gd)
		}
		return d
	}
	d.Request = &requestDTO{Type: string(i.Request.Type), SQL: i.Request.SQL}
	for _, r := range i.Responses {
		d.Responses = append(d.Responses, responseToDTO(r))
	}
	return d
}

func clientMessageToDTO(m model.ClientMessage) clientMessageDTO {
	d := clientMessageDTO{
		Type: string(m.Type), Statement: m.Statement, Portal: m.Portal, SQL: m.SQL,
		ParamTypes: m.ParamTypes, ParamFormats: m.ParamFormats, ResultFormats: m.ResultFormats,
		Target: string(m.Target), Name: m.Name, MaxRows: m.MaxRows,
	}
	for _, v := range m.Params {
		d.Params = append(d.Params, valueToDTO(v))
	}
	return d
}

func responseToDTO(r model.Response) responseDTO {
	d := responseDTO{Type: string(r.Type), Tag: r.Tag, Notice: r.Fields, Name: r.Name, Value: r.Value, TxStatus: r.TxStatus, ParamTypes: r.ParamTypes}
	for _, c := range r.Columns {
		d.Fields = append(d.Fields, columnDTO(c))
	}
	for _, v := range r.Values {
		d.Values = append(d.Values, valueToDTO(v))
	}
	return d
}

func valueToDTO(v model.Value) valueDTO {
	if v.Null {
		return valueDTO{Null: true}
	}
	if utf8.Valid(v.Bytes) {
		s := string(v.Bytes)
		return valueDTO{Text: &s}
	}
	b := base64.StdEncoding.EncodeToString(v.Bytes)
	return valueDTO{Base64: &b}
}

// fromDTO prüft die Struktur, an der eine abgeschnittene Datei auffällt: Die
// Liste der Sessions fehlt nicht, jede Session trägt Interaktionen mit
// fortlaufender Nummer, und jede Interaktion hat die Form, die
// model.Interaction.Validate verlangt — eine einfache endet mit ready_for_query,
// eine Extended-Interaktion mit der Sync-Gruppe und deren ready_for_query
// (LH-FA-07, LH-FA-18.a). Eine Datei, die genau an einer Session-Grenze endet,
// bleibt eine gültige kürzere Aufzeichnung.
func fromDTO(d recordingDTO) (model.Recording, error) {
	rec := model.Recording{Format: d.Format, Version: d.Version, EmptySessions: d.EmptySessions}
	if d.Sessions == nil {
		return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Liste der Sessions fehlt")
	}
	for si, sd := range *d.Sessions {
		if sd.ID != si+1 {
			return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Session %d trägt die Kennung %d", si+1, sd.ID)
		}
		if len(sd.Interactions) == 0 && !d.EmptySessions {
			return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Session %d ohne Interaktion", sd.ID)
		}
		s := model.Session{ID: sd.ID, Startup: sd.Startup, ServerParameters: sd.ServerParameters}
		for ii, id := range sd.Interactions {
			if id.Sequence != ii+1 {
				return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Session %d: Interaktion %d trägt die Nummer %d", sd.ID, ii+1, id.Sequence)
			}
			if id.OffsetMS != nil && *id.OffsetMS < 0 {
				return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Session %d, Interaktion %d: offset_ms negativ", sd.ID, id.Sequence)
			}
			i, err := interactionFromDTO(id)
			if err == nil {
				err = i.Validate()
			}
			if err != nil {
				return model.Recording{}, model.Errorf(model.CodeRecordingBroken, err, "Session %d, Interaktion %d", sd.ID, id.Sequence)
			}
			s.Interactions = append(s.Interactions, i)
		}
		rec.Sessions = append(rec.Sessions, s)
	}
	return rec, nil
}

// interactionFromDTO liest die Art aus type: fehlt es, ist es eine einfache
// Anfrage mit request und responses, bei `extended` eine Extended-Interaktion
// mit groups. Ein anderer Wert und jede Mischung der beiden Formen sind ein
// Fehler (SPEC-001, SPEC-041).
func interactionFromDTO(d interactionDTO) (model.Interaction, error) {
	i := model.Interaction{Sequence: d.Sequence, OffsetMS: d.OffsetMS}
	switch d.Type {
	case "":
		if d.Request == nil || d.Groups != nil {
			return model.Interaction{}, fmt.Errorf("einfache Anfrage braucht request und trägt keine groups")
		}
		if model.RequestType(d.Request.Type) != model.RequestQuery {
			return model.Interaction{}, fmt.Errorf("Anfrage-Typ %q unbekannt", d.Request.Type)
		}
		i.Request = model.Request{Type: model.RequestQuery, SQL: d.Request.SQL}
		for _, rd := range d.Responses {
			r, err := responseFromDTO(rd)
			if err != nil {
				return model.Interaction{}, err
			}
			i.Responses = append(i.Responses, r)
		}
	case interactionExtended:
		if d.Request != nil || d.Responses != nil {
			return model.Interaction{}, fmt.Errorf("Extended-Interaktion trägt weder request noch responses")
		}
		i.Request = model.Request{Type: model.RequestExtended}
		for _, gd := range d.Groups {
			g, err := groupFromDTO(gd)
			if err != nil {
				return model.Interaction{}, err
			}
			i.Groups = append(i.Groups, g)
		}
	default:
		return model.Interaction{}, fmt.Errorf("Interaktions-Typ %q unbekannt", d.Type)
	}
	return i, nil
}

func groupFromDTO(d groupDTO) (model.Group, error) {
	var g model.Group
	for _, md := range d.Client {
		m := model.ClientMessage{
			Type: model.ClientMessageType(md.Type), Statement: md.Statement, Portal: md.Portal, SQL: md.SQL,
			ParamTypes: leerAlsNil(md.ParamTypes), ParamFormats: leerAlsNil(md.ParamFormats),
			ResultFormats: leerAlsNil(md.ResultFormats),
			Target:        model.Target(md.Target), Name: md.Name, MaxRows: md.MaxRows,
		}
		for _, vd := range md.Params {
			v, err := valueFromDTO(vd)
			if err != nil {
				return model.Group{}, err
			}
			m.Params = append(m.Params, v)
		}
		g.Client = append(g.Client, m)
	}
	for _, rd := range d.Server {
		r, err := responseFromDTO(rd)
		if err != nil {
			return model.Group{}, err
		}
		g.Server = append(g.Server, r)
	}
	return g, nil
}

// leerAlsNil liest eine leere Liste als nil; der Schreiber gibt nil als []
// aus.
func leerAlsNil[T any](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	return s
}

// responseFromDTO übernimmt die Felder; welche Typen in welcher Art zulässig
// sind, prüft model.Interaction.Validate.
func responseFromDTO(d responseDTO) (model.Response, error) {
	r := model.Response{Type: model.ResponseType(d.Type), Tag: d.Tag, Fields: d.Notice, Name: d.Name, Value: d.Value, TxStatus: d.TxStatus, ParamTypes: leerAlsNil(d.ParamTypes)}
	for _, c := range d.Fields {
		r.Columns = append(r.Columns, model.Column(c))
	}
	for _, vd := range d.Values {
		v, err := valueFromDTO(vd)
		if err != nil {
			return model.Response{}, err
		}
		r.Values = append(r.Values, v)
	}
	return r, nil
}

func valueFromDTO(d valueDTO) (model.Value, error) {
	gesetzt := 0
	for _, b := range []bool{d.Null, d.Text != nil, d.Base64 != nil} {
		if b {
			gesetzt++
		}
	}
	if gesetzt > 1 {
		return model.Value{}, model.Errorf(model.CodeRecordingBroken, nil, "Wert mit mehr als einem von text, null und base64")
	}
	switch {
	case d.Null:
		return model.Value{Null: true}, nil
	case d.Text != nil:
		return model.Value{Bytes: []byte(*d.Text)}, nil
	case d.Base64 != nil:
		b, err := base64.StdEncoding.DecodeString(*d.Base64)
		if err != nil {
			return model.Value{}, model.Errorf(model.CodeRecordingBroken, err, "Base64-Wert nicht lesbar")
		}
		return model.Value{Bytes: b}, nil
	default:
		return model.Value{}, model.Errorf(model.CodeRecordingBroken, nil, "Wert ohne text, null oder base64")
	}
}
