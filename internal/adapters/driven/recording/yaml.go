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

type interactionDTO struct {
	Sequence  int           `yaml:"sequence"`
	OffsetMS  *int64        `yaml:"offset_ms,omitempty"`
	Request   requestDTO    `yaml:"request"`
	Responses []responseDTO `yaml:"responses"`
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
			id := interactionDTO{Sequence: i.Sequence, OffsetMS: i.OffsetMS, Request: requestDTO{Type: string(i.Request.Type), SQL: i.Request.SQL}}
			for _, r := range i.Responses {
				id.Responses = append(id.Responses, responseToDTO(r))
			}
			sd.Interactions = append(sd.Interactions, id)
		}
		sessions = append(sessions, sd)
	}
	return d
}

func responseToDTO(r model.Response) responseDTO {
	d := responseDTO{Type: string(r.Type), Tag: r.Tag, Notice: r.Fields, Name: r.Name, Value: r.Value, TxStatus: r.TxStatus}
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

// bekannteAntworten sind die Antworttypen, die dieser Leser kennt; ein anderer
// Typ macht die Aufzeichnung zu einer beschädigten (SPEC-001).
var bekannteAntworten = map[model.ResponseType]bool{
	model.ResponseRowDescription:     true,
	model.ResponseDataRow:            true,
	model.ResponseCommandComplete:    true,
	model.ResponseEmptyQueryResponse: true,
	model.ResponseErrorResponse:      true,
	model.ResponseNoticeResponse:     true,
	model.ResponseParameterStatus:    true,
	model.ResponseReadyForQuery:      true,
}

// fromDTO prüft die Struktur, an der eine abgeschnittene Datei auffällt: Die
// Liste der Sessions fehlt nicht, jede Session trägt Interaktionen mit
// fortlaufender Nummer, und jede Interaktion endet mit ready_for_query
// (LH-FA-07). Eine Datei, die genau an einer Session-Grenze endet, bleibt eine
// gültige kürzere Aufzeichnung.
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
			if model.RequestType(id.Request.Type) != model.RequestQuery {
				return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Anfrage-Typ %q unbekannt", id.Request.Type)
			}
			if id.Sequence != ii+1 {
				return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Session %d: Interaktion %d trägt die Nummer %d", sd.ID, ii+1, id.Sequence)
			}
			if n := len(id.Responses); n == 0 || id.Responses[n-1].Type != string(model.ResponseReadyForQuery) {
				return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Session %d, Interaktion %d endet nicht mit ready_for_query", sd.ID, id.Sequence)
			}
			if id.OffsetMS != nil && *id.OffsetMS < 0 {
				return model.Recording{}, model.Errorf(model.CodeRecordingBroken, nil, "Session %d, Interaktion %d: offset_ms negativ", sd.ID, id.Sequence)
			}
			i := model.Interaction{Sequence: id.Sequence, OffsetMS: id.OffsetMS, Request: model.Request{Type: model.RequestType(id.Request.Type), SQL: id.Request.SQL}}
			for _, rd := range id.Responses {
				r, err := responseFromDTO(rd)
				if err != nil {
					return model.Recording{}, err
				}
				i.Responses = append(i.Responses, r)
			}
			s.Interactions = append(s.Interactions, i)
		}
		rec.Sessions = append(rec.Sessions, s)
	}
	return rec, nil
}

func responseFromDTO(d responseDTO) (model.Response, error) {
	if !bekannteAntworten[model.ResponseType(d.Type)] {
		return model.Response{}, model.Errorf(model.CodeRecordingBroken, nil, "Antwort-Typ %q unbekannt", d.Type)
	}
	r := model.Response{Type: model.ResponseType(d.Type), Tag: d.Tag, Fields: d.Notice, Name: d.Name, Value: d.Value, TxStatus: d.TxStatus}
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
