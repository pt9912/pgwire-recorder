package bootstrap_test

import (
	"bytes"
	"context"
	"net"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// beendeterLauf führt kommando mit einem schon beendeten ctx: Der Lauf fährt
// nach dem Start sofort herunter.
func beendeterLauf(t *testing.T, args ...string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var stdout, stderr bytes.Buffer
	code := bootstrap.Run(ctx, nil, args, "dev", &stdout, &stderr)
	return code, stderr.String()
}

// Abdeckung: LH-FA-12/Happy, LH-FA-22/Happy — record reicht --record-empty-sessions
// an den Service: Ein Lauf ohne Verbindung schreibt eine Aufzeichnung ohne
// Sessions, die mit der Option empty_sessions trägt und ohne nicht, in beiden
// Formaten (LH-FA-12.a *Kennzeichnung*, LH-FA-22.a).
func TestRunRecordEmptySessionsKennzeichnung(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS", "")
	for _, format := range []string{"yaml", "sqlite"} {
		for _, mit := range []bool{true, false} {
			out := filepath.Join(t.TempDir(), "rec."+format)
			args := []string{"record", "--listen", "127.0.0.1:0", "--upstream", "127.0.0.1:1", "--output", out, "--format", format}
			if mit {
				args = append(args, "--record-empty-sessions")
			}
			if code, stderr := beendeterLauf(t, args...); code != 0 {
				t.Fatalf("%s, Option %v: Exit-Code %d, stderr %q", format, mit, code, stderr)
			}
			rec, err := recording.Leser{}.Load(context.Background(), out)
			if err != nil || rec.EmptySessions != mit || len(rec.Sessions) != 0 {
				t.Errorf("%s, Option %v: %#v, %v", format, mit, rec, err)
			}
		}
	}
}

// Abdeckung: LH-FA-12/Negative, LH-FA-13/Negative — replay mit
// --session-assignment connection und einer Aufzeichnung ohne Kennzeichnung
// empty_sessions endet als Startfehler PGR-E2001 mit Exit-Code 2, bevor der Port
// geöffnet wird: Der Port ist belegt, und die Meldung ist trotzdem nicht der Fehler
// des Öffnens; mit Kennzeichnung läuft der Start bis zum Öffnen des Ports (hier
// PGR-E4001, weil er belegt ist). Mit der Zuordnung first-request ist die
// Aufzeichnung ohne Kennzeichnung gültig (LH-FA-12.a *Kennzeichnung bei
// `connection`*).
func TestRunReplayConnectionKennzeichnung(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_SESSION_ASSIGNMENT", "")
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	belegt, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer belegt.Close()
	ohne := aufzeichnung(t)
	rec := model.NewRecording()
	rec.EmptySessions = true
	rec.Sessions = []model.Session{{ID: 1, Interactions: []model.Interaction{{
		Sequence:  1,
		Request:   model.Request{Type: model.RequestQuery, SQL: "SELECT 1"},
		Responses: []model.Response{{Type: model.ResponseCommandComplete, Tag: "SELECT 0"}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}},
	}}}}
	mit := filepath.Join(t.TempDir(), "mit.yaml")
	if err := (recording.YAML{}).Write(context.Background(), mit, rec); err != nil {
		t.Fatal(err)
	}
	for _, f := range []struct {
		name string
		args []string
		code int
		want string
	}{
		{"connection ohne Kennzeichnung", []string{"--input", ohne, "--session-assignment", "connection"}, 2, "[PGR-E2001]"},
		{"connection mit Kennzeichnung", []string{"--input", mit, "--session-assignment", "connection"}, 4, "[PGR-E4001]"},
		{"first-request ohne Kennzeichnung", []string{"--input", ohne, "--session-assignment", "first-request"}, 4, "[PGR-E4001]"},
	} {
		code, stderr := beendeterLauf(t, append([]string{"replay", "--listen", belegt.Addr().String()}, f.args...)...)
		if code != f.code || !strings.Contains(stderr, f.want) || strings.Contains(stderr, "replay gestartet") {
			t.Errorf("%s: Exit-Code %d, stderr %q, erwartet %d mit %s", f.name, code, stderr, f.code, f.want)
		}
	}
}

// Abdeckung: LH-FA-12/Boundary, LH-FA-03/Boundary — replay reicht
// --session-assignment connection an den Service: Mit der Kennzeichnung läuft der
// Start durch und das Prozessende meldet die Session, die keine Verbindung
// erhielt (PGR-W2001, mit --fail-on-unconsumed PGR-E5002 und Exit-Code 5)
// (LH-FA-12.a, LH-FA-03.b).
func TestRunReplayConnectionLaeuft(t *testing.T) {
	t.Setenv("PGWIRE_RECORDER_FAIL_ON_UNCONSUMED", "")
	rec := model.NewRecording()
	rec.EmptySessions = true
	rec.Sessions = []model.Session{{ID: 1, Interactions: []model.Interaction{{
		Sequence:  1,
		Request:   model.Request{Type: model.RequestQuery, SQL: "SELECT 1"},
		Responses: []model.Response{{Type: model.ResponseCommandComplete, Tag: "SELECT 0"}, {Type: model.ResponseReadyForQuery, TxStatus: "I"}},
	}}}}
	mit := filepath.Join(t.TempDir(), "mit.yaml")
	if err := (recording.YAML{}).Write(context.Background(), mit, rec); err != nil {
		t.Fatal(err)
	}
	code, stderr := beendeterLauf(t, "replay", "--listen", "127.0.0.1:0", "--input", mit, "--session-assignment", "connection")
	if code != 0 || !strings.Contains(stderr, "replay gestartet") || !strings.Contains(stderr, "code=PGR-W2001") {
		t.Errorf("Exit-Code %d, stderr %q", code, stderr)
	}
	code, stderr = beendeterLauf(t, "replay", "--listen", "127.0.0.1:0", "--input", mit, "--session-assignment", "connection", "--fail-on-unconsumed")
	if code != 5 || !strings.Contains(stderr, "PGR-E5002") {
		t.Errorf("mit --fail-on-unconsumed: Exit-Code %d, stderr %q", code, stderr)
	}
}
