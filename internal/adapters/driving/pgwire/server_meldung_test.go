package pgwire

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// logZeilen liest Log-Zeilen im Format logfmt: je Zeile Schlüssel und Wert, ein
// Wert in Anführungszeichen mit Escapes.
func logZeilen(t *testing.T, text string) []map[string]string {
	t.Helper()
	var out []map[string]string
	for _, zeile := range strings.Split(strings.TrimSpace(text), "\n") {
		if zeile == "" {
			continue
		}
		felder := map[string]string{}
		for zeile != "" {
			k, rest, ok := strings.Cut(zeile, "=")
			if !ok {
				t.Fatalf("keine logfmt-Zeile: %q", zeile)
			}
			var v string
			if strings.HasPrefix(rest, `"`) {
				ende := 1
				for ende < len(rest) && rest[ende] != '"' {
					if rest[ende] == '\\' {
						ende++
					}
					ende++
				}
				u, err := strconv.Unquote(rest[:ende+1])
				if err != nil {
					t.Fatalf("Wert %q: %v", rest[:ende+1], err)
				}
				v, rest = u, rest[ende+1:]
			} else {
				v, rest, _ = strings.Cut(rest, " ")
				rest = " " + rest
			}
			felder[k] = v
			zeile = strings.TrimPrefix(rest, " ")
		}
		out = append(out, felder)
	}
	return out
}

// Abdeckung: LH-FA-14/Negative, LH-FA-13/Negative — je Fehlerklasse trägt die
// ErrorResponse eines Verbindungsfehlers FATAL, den SQLSTATE der Klasse
// (0A000, 08006, sonst XX000), als Meldungstext genau den Fehlertext des
// Attributs error der Log-Zeile, die auf der Stufe error steht, und keine
// weiteren Felder; der Lauf merkt den Code (SPEC-034 §Ausgabe *Zustellung an
// den Client*).
func TestFehlerantwortJeKlasse(t *testing.T) {
	for _, f := range []struct {
		err      error
		code     string
		sqlstate string
	}{
		{errors.New("nicht\neingeordnet"), model.CodeInternal, "XX000"},
		{model.Errorf(model.CodeUsage, nil, "x"), model.CodeUsage, "XX000"},
		{model.Errorf(model.CodeRecordingIO, errors.New("voll"), "x"), model.CodeRecordingIO, "XX000"},
		{model.Errorf(model.CodeUpstream, nil, "x"), model.CodeUpstream, "08006"},
		{model.Errorf(model.CodeReplayMismatch, nil, "x"), model.CodeReplayMismatch, "XX000"},
		{model.Errorf(model.CodeUnsupported, nil, "x"), model.CodeUnsupported, "0A000"},
	} {
		var log, netz bytes.Buffer
		s := &Server{log: slog.New(slog.NewTextHandler(&log, nil))}
		s.fail(pgproto3.NewBackend(&bytes.Buffer{}, &netz), f.err)
		msg, err := pgproto3.NewFrontend(&netz, io.Discard).Receive()
		if err != nil {
			t.Fatal(err)
		}
		e, ok := msg.(*pgproto3.ErrorResponse)
		if !ok {
			t.Fatalf("%s: %T", f.code, msg)
		}
		zeilen := logZeilen(t, log.String())
		if len(zeilen) != 1 || zeilen[0]["level"] != "ERROR" || zeilen[0]["code"] != f.code {
			t.Fatalf("%s: Log %q", f.code, log.String())
		}
		text := zeilen[0]["error"]
		if !strings.Contains(text, " ["+f.code+"]: ") || strings.Count(text, "[PGR-") != 1 || strings.Contains(text, "\n") {
			t.Fatalf("%s: Fehlertext %q", f.code, text)
		}
		want := pgproto3.ErrorResponse{Severity: "FATAL", SeverityUnlocalized: "FATAL", Code: f.sqlstate, Message: text}
		if !reflect.DeepEqual(*e, want) {
			t.Fatalf("%s: ErrorResponse %+v, erwartet %+v", f.code, *e, want)
		}
		if s.FirstErrorCode() != f.code {
			t.Fatalf("%s: gemerkt %q", f.code, s.FirstErrorCode())
		}
	}
}

// Abdeckung: LH-FA-13/Negative, LH-FA-14/Negative — entstehen beim Ende einer
// Verbindung mehrere Fehler nebeneinander, ist jeder klassifizierte eine eigene
// Log-Zeile mit eigenem Kopf, in der Reihenfolge des Entstehens, und der erste
// wird gemerkt; ein nicht klassifizierter folgt der ersten als Ursache
// (SPEC-034 §Ausgabe *Gleichrangige Fehler*, LH-FA-13.b).
func TestNoteGleichrangig(t *testing.T) {
	var log bytes.Buffer
	s := &Server{log: slog.New(slog.NewTextHandler(&log, nil))}
	s.note(errors.Join(
		model.Errorf(model.CodeConnectionLost, nil, "Client weg"),
		model.Errorf(model.CodeRecordingIO, nil, "r.yaml nicht zu schreiben"),
		errors.New("upstream close"),
	))
	zeilen := logZeilen(t, log.String())
	if len(zeilen) != 2 ||
		zeilen[0]["code"] != model.CodeConnectionLost || zeilen[0]["error"] != "Netzwerk [PGR-E4003]: Client weg: upstream close" ||
		zeilen[1]["code"] != model.CodeRecordingIO || zeilen[1]["error"] != "Recording [PGR-E3001]: r.yaml nicht zu schreiben" {
		t.Fatalf("Log: %s", log.String())
	}
	if s.FirstErrorCode() != model.CodeConnectionLost {
		t.Fatalf("gemerkt %q", s.FirstErrorCode())
	}
}

// annahmeFehler liefert beim ersten Accept einen Fehler, danach net.ErrClosed.
type annahmeFehler struct {
	net.Listener
	mu    sync.Mutex
	runde int
}

func (a *annahmeFehler) Accept() (net.Conn, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.runde++
	if a.runde == 1 {
		return nil, errors.New("accept: too many open files")
	}
	return nil, net.ErrClosed
}

// Abdeckung: LH-FA-13/Negative — eine Verbindung, die der Listener nicht
// annehmen kann, ist der Verbindungsfehler PGR-E4000 mit Kopf, er wird gemerkt,
// und Serve nimmt danach weiter an (LH-FA-13.b).
func TestAnnahmefehler(t *testing.T) {
	var log syncBuffer
	s := NewReplayServer(&fakeReplayer{}, slog.New(slog.NewTextHandler(&log, nil)))
	l := &annahmeFehler{}
	s.Serve(context.Background(), l)
	if l.runde != 2 {
		t.Fatalf("Serve nach dem Fehler nicht weiter: %d Runden", l.runde)
	}
	zeilen := logZeilen(t, log.String())
	if len(zeilen) != 1 || zeilen[0]["level"] != "ERROR" || zeilen[0]["code"] != model.CodeNetwork ||
		zeilen[0]["error"] != "Netzwerk [PGR-E4000]: Verbindung nicht anzunehmen: accept: too many open files" {
		t.Fatalf("Log: %s", log.String())
	}
	if s.FirstErrorCode() != model.CodeNetwork {
		t.Fatalf("gemerkt %q", s.FirstErrorCode())
	}
}

// Eine Verbindung ohne Startnachricht endet still; die Zeile der Stufe debug
// nennt den Text der Bibliothek unter grund, nicht unter error (LH-FA-14.a
// §Zeilenform).
func TestOhneStartnachrichtGrund(t *testing.T) {
	var log syncBuffer
	s := NewReplayServer(&fakeReplayer{}, slog.New(slog.NewTextHandler(&log, &slog.HandlerOptions{Level: slog.LevelDebug})))
	client, serverSeite := net.Pipe()
	_ = client.Close()
	fertig := make(chan struct{})
	go func() {
		s.handle(context.Background(), serverSeite)
		close(fertig)
	}()
	select {
	case <-fertig:
	case <-time.After(5 * time.Second):
		t.Fatal("handle endet nicht")
	}
	zeilen := logZeilen(t, log.String())
	if len(zeilen) != 1 || zeilen[0]["level"] != "DEBUG" || zeilen[0]["grund"] == "" {
		t.Fatalf("Log: %s", log.String())
	}
	if _, ok := zeilen[0]["error"]; ok {
		t.Fatalf("Attribut error an einer Zeile der Stufe debug: %s", log.String())
	}
	if s.FirstErrorCode() != "" {
		t.Fatalf("gemerkt %q", s.FirstErrorCode())
	}
}

// bisRFQ liest Server-Nachrichten bis einschließlich ReadyForQuery.
func bisRFQ(t *testing.T, fe *pgproto3.Frontend) {
	t.Helper()
	for {
		msg, err := fe.Receive()
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := msg.(*pgproto3.ReadyForQuery); ok {
			return
		}
	}
}

// ohneZeileUnterWarn bricht ab, wenn das Log eine Zeile der Stufe info oder
// debug trägt.
func ohneZeileUnterWarn(t *testing.T, modus, log string) {
	t.Helper()
	for _, z := range logZeilen(t, log) {
		if z["level"] == "INFO" || z["level"] == "DEBUG" {
			t.Fatalf("%s: Zeile der Stufe %s je Verbindung oder Anfrage:\n%s", modus, z["level"], log)
		}
	}
}

// Abdeckung: LH-FA-14/Boundary — auf der Stufe info schreibt der PGWire-Adapter
// keine Zeile je Verbindung oder Interaktion, weder in record noch in replay,
// für einfache Anfragen und Extended-Interaktionen (LH-FA-14.a §Log-Level).
func TestInfoOhneZeileJeVerbindung(t *testing.T) {
	var recLog syncBuffer
	rec := &fakeRecorder{}
	client, _ := verbindeMitLog(t, rec, &recLog)
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	for _, sql := range []string{"SELECT 1", "SELECT 2"} {
		fe.Send(&pgproto3.Query{String: sql})
		_ = fe.Flush()
		bisRFQ(t, fe)
	}
	fe.Send(&pgproto3.Terminate{})
	_ = fe.Flush()
	if end := rec.lastEnd(t); end != model.EndTerminate {
		t.Fatalf("Session-Ende: %v", end)
	}
	ohneZeileUnterWarn(t, "record", recLog.String())

	_, _, fe, repLog, fertig := replayVerbindung(t, &fakeReplayer{})
	fe.Send(&pgproto3.Query{String: "SELECT 1"})
	_ = fe.Flush()
	bisRFQ(t, fe)
	fe.Send(&pgproto3.Parse{Name: "s", Query: "SELECT 1"})
	fe.Send(&pgproto3.Sync{})
	_ = fe.Flush()
	bisRFQ(t, fe)
	fe.Send(&pgproto3.Terminate{})
	_ = fe.Flush()
	warte(t, fertig)
	ohneZeileUnterWarn(t, "replay", repLog.String())
}

// Abdeckung: LH-FA-13/Negative, LH-FA-14/Negative — ergibt ein Fehler mehrere
// Meldungen, stellt fail genau eine ErrorResponse zu, mit Text und SQLSTATE der
// ersten, gemerkten Meldung; jede Meldung steht als eigene Zeile der Stufe
// error im Log (SPEC-034 §Ausgabe *Zustellung an den Client*).
func TestFailGleichrangig(t *testing.T) {
	var log, netz bytes.Buffer
	s := &Server{log: slog.New(slog.NewTextHandler(&log, nil))}
	s.fail(pgproto3.NewBackend(&bytes.Buffer{}, &netz), errors.Join(
		model.Errorf(model.CodeUnsupported, nil, "nicht unterstützt"),
		model.Errorf(model.CodeConnectionLost, nil, "weg"),
	))
	fe := pgproto3.NewFrontend(&netz, io.Discard)
	msg, err := fe.Receive()
	if err != nil {
		t.Fatal(err)
	}
	want := pgproto3.ErrorResponse{Severity: "FATAL", SeverityUnlocalized: "FATAL", Code: "0A000", Message: "nicht unterstützt [PGR-E6001]: nicht unterstützt"}
	if e, ok := msg.(*pgproto3.ErrorResponse); !ok || !reflect.DeepEqual(*e, want) {
		t.Fatalf("ErrorResponse %#v, erwartet %+v", msg, want)
	}
	if netz.Len() > 0 {
		t.Fatalf("nach der ersten ErrorResponse noch %d Bytes", netz.Len())
	}
	zeilen := logZeilen(t, log.String())
	if len(zeilen) != 2 || zeilen[0]["level"] != "ERROR" || zeilen[0]["code"] != model.CodeUnsupported ||
		zeilen[1]["level"] != "ERROR" || zeilen[1]["code"] != model.CodeConnectionLost {
		t.Fatalf("Log: %s", log.String())
	}
	if s.FirstErrorCode() != model.CodeUnsupported {
		t.Fatalf("gemerkt %q", s.FirstErrorCode())
	}
}
