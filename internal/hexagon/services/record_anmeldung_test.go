package services_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// anmeldeUpstream ruft beim Aufbau den Rückruf mit den Anforderungen auf, wie der
// Upstream-Adapter, und hält die Antworten.
type anmeldeUpstream struct {
	fakeUpstream
	folge     []model.Anmeldeanforderung
	antworten [][]byte
}

func (u *anmeldeUpstream) Open(ctx context.Context, startup map[string]string, anmeldung model.Anmeldung) (driven.UpstreamSession, []model.Response, error) {
	for _, a := range u.folge {
		antwort, err := anmeldung(ctx, a)
		if err != nil {
			return nil, nil, err
		}
		u.antworten = append(u.antworten, antwort)
	}
	return u.fakeUpstream.Open(ctx, startup, anmeldung)
}

// Abdeckung: LH-FA-05/Happy — OpenSession reicht den Rückruf des Aufrufers an den
// Upstream: Jede Anforderung kommt unverändert beim Aufrufer an, jede Antwort
// unverändert zurück; liefert der Rückruf einen Fehler, endet OpenSession mit
// genau ihm und ohne Session, auch mit LeereSessions (LH-FA-05.b *Vermittlung*,
// *Ende durch den Client*).
func TestRecordAnmeldungRueckruf(t *testing.T) {
	up := &anmeldeUpstream{folge: []model.Anmeldeanforderung{{Code: 5, Daten: []byte{1, 2, 3, 4}}, {Code: 12}}}
	s, _ := neu(t, up)
	var gesehen []model.Anmeldeanforderung
	id, _, err := s.OpenSession(context.Background(), nil, func(_ context.Context, a model.Anmeldeanforderung) ([]byte, error) {
		gesehen = append(gesehen, a)
		return []byte("antwort\x00"), nil
	})
	if err != nil || id == 0 {
		t.Fatalf("OpenSession: id=%d err=%v", id, err)
	}
	if !reflect.DeepEqual(gesehen, up.folge) || len(up.antworten) != 2 || string(up.antworten[0]) != "antwort\x00" {
		t.Errorf("Anforderungen %v, Antworten %q", gesehen, up.antworten)
	}

	s, repo := neu(t, up)
	services.LeereSessions(true)(s)
	id, _, err = s.OpenSession(context.Background(), nil, func(context.Context, model.Anmeldeanforderung) ([]byte, error) {
		return nil, model.ErrClientEnde
	})
	if id != 0 || err != model.ErrClientEnde {
		t.Fatalf("OpenSession: id=%d err=%v", id, err)
	}
	if err := s.Finish(context.Background()); err != nil {
		t.Fatal(err)
	}
	if n := len(repo.last(t).Sessions); n != 0 {
		t.Errorf("%d Sessions nach dem Ende des Clients", n)
	}
}

// Abdeckung: LH-FA-05/Boundary, LH-RB-01/Boundary — weder Anforderungen noch
// Antworten der Anmeldung stehen in der Aufzeichnung (LH-FA-05.b
// *Aufzeichnung*).
func TestRecordAnmeldungNichtInDerAufzeichnung(t *testing.T) {
	up := &anmeldeUpstream{folge: []model.Anmeldeanforderung{{Code: 10, Daten: []byte("SALZ")}, {Code: 12, Daten: []byte("SALZ")}}}
	s, repo := neu(t, up)
	ctx := context.Background()
	id, _, err := s.OpenSession(ctx, map[string]string{"user": "app"}, func(context.Context, model.Anmeldeanforderung) ([]byte, error) {
		return []byte("GEHEIM"), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
		t.Fatal(err)
	}
	if err := s.CloseSession(ctx, id, model.EndTerminate); err != nil {
		t.Fatal(err)
	}
	text := alleTexte(reflect.ValueOf(repo.last(t)))
	if !strings.Contains(text, "SELECT 1") || strings.Contains(text, "GEHEIM") || strings.Contains(text, "SALZ") {
		t.Errorf("Aufzeichnung: %s", text)
	}
}

// alleTexte sammelt jeden Text und jede Bytefolge unterhalb von v.
func alleTexte(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String() + "\n"
	case reflect.Slice, reflect.Array:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes()) + "\n"
		}
		var out string
		for i := range v.Len() {
			out += alleTexte(v.Index(i))
		}
		return out
	case reflect.Map:
		var out string
		for _, k := range v.MapKeys() {
			out += alleTexte(k) + alleTexte(v.MapIndex(k))
		}
		return out
	case reflect.Struct:
		var out string
		for i := range v.NumField() {
			out += alleTexte(v.Field(i))
		}
		return out
	case reflect.Pointer:
		if v.IsNil() {
			return ""
		}
		return alleTexte(v.Elem())
	}
	return ""
}
