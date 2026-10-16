package bootstrap_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
	"github.com/pt9912/pgwire-recorder/internal/testuhr"
)

// zeitUpstream ist ein Upstream ohne Server: Query antwortet mit CommandComplete und
// ReadyForQuery, eine Gruppe mit ParseComplete und ReadyForQuery.
type zeitUpstream struct{}

func (zeitUpstream) Open(context.Context, map[string]string, model.Anmeldung) (driven.UpstreamSession, []model.Response, error) {
	return zeitSession{}, []model.Response{{Type: model.ResponseReadyForQuery, TxStatus: "I"}}, nil
}

type zeitSession struct{}

func (zeitSession) Query(context.Context, string) ([]model.Response, error) {
	return []model.Response{{Type: model.ResponseCommandComplete, Tag: "SELECT 1"}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}, nil
}
func (zeitSession) Send(context.Context, []model.ClientMessage) error { return nil }
func (zeitSession) Receive(context.Context) ([]model.Response, error) {
	return []model.Response{{Type: model.ResponseParseComplete}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}}, nil
}
func (zeitSession) Close() error { return nil }

// Abdeckung: LH-FA-21/Happy, LH-FA-22/Happy — Roundtrip je Format: Eine Aufzeichnung mit
// Zeitangaben aus einfacher und Extended-Interaktion, geschrieben als yaml und als sqlite,
// lädt über den Leser des Formats mit denselben Werten von offset_ms, die der Service
// nach der Uhr gemessen hat (LH-FA-21.a *Aufzeichnen*, *Prüfen*); ohne Zeitangaben lädt sie
// ohne Feld.
func TestRecordZeitRoundtrip(t *testing.T) {
	for _, format := range []string{"yaml", "sqlite"} {
		for _, mit := range []bool{true, false} {
			path := filepath.Join(t.TempDir(), "rec."+format)
			repo, err := recording.New(format)
			if err != nil {
				t.Fatal(err)
			}
			uhr := testuhr.Neu()
			ctx := context.Background()
			s, err := services.NewRecordService(ctx, zeitUpstream{}, repo, path, false, services.Zeitangaben(uhr, mit))
			if err != nil {
				t.Fatal(err)
			}
			uhr.Vorruecken(1000)
			id, _, err := s.OpenSession(ctx, map[string]string{"user": "app"}, nil)
			if err != nil {
				t.Fatal(err)
			}
			uhr.Vorruecken(12)
			if _, err := s.Query(ctx, id, "SELECT 1"); err != nil {
				t.Fatal(err)
			}
			s.Delivered(ctx, id)
			uhr.Vorruecken(30)
			for _, m := range []model.ClientMessage{{Type: model.ClientParse, Statement: "s", SQL: "SELECT 1"}, {Type: model.ClientSync}} {
				if err := s.ClientMessage(ctx, id, m); err != nil {
					t.Fatal(err)
				}
			}
			uhr.Vorruecken(8)
			if _, err := s.AwaitServer(ctx, id); err != nil {
				t.Fatal(err)
			}
			s.Delivered(ctx, id)
			if err := s.CloseSession(ctx, id, model.EndClosed); err != nil {
				t.Fatal(err)
			}
			if err := s.Finish(ctx); err != nil {
				t.Fatal(err)
			}
			rec, err := recording.Leser{}.Load(ctx, path)
			if err != nil {
				t.Fatalf("%s: laden: %v", format, err)
			}
			if len(rec.Sessions) != 1 || len(rec.Sessions[0].Interactions) != 2 {
				t.Fatalf("%s: Sessions %#v", format, rec.Sessions)
			}
			for i, want := range []int64{12, 42} {
				got := rec.Sessions[0].Interactions[i].OffsetMS
				if !mit && got != nil {
					t.Errorf("%s ohne Zeitangaben, Interaktion %d: offset_ms %d", format, i+1, *got)
				}
				if mit && got == nil {
					t.Errorf("%s, Interaktion %d: ohne offset_ms, erwartet %d", format, i+1, want)
				} else if mit && *got != want {
					t.Errorf("%s, Interaktion %d: offset_ms %d, erwartet %d", format, i+1, *got, want)
				}
			}
		}
	}
}
