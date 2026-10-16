package pgwire_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/pgwire"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// recordEndeFolge sind die Antworten einer Interaktion mit Ende: eine Antwort davor, die
// Fehlerantwort des Servers, connection_closed.
func recordEndeFolge(erste model.Response) []model.Response {
	return []model.Response{
		erste,
		{Type: model.ResponseErrorResponse, Fields: map[string]string{"S": "FATAL", "V": "FATAL", "C": "57P01", "M": "terminating connection"}},
		{Type: model.ResponseConnectionClosed},
	}
}

// recordWeg lässt die Antworten rs auf dem Weg einfach (Query) oder erweitert
// (Parse, Sync, Antworten über AwaitServer) beim Client ankommen und liefert den
// Recorder, den Server, die Client-Seite und das Log.
func recordWeg(t *testing.T, erweitert bool, rs []model.Response, zugestelltEndet bool) (*fakeRecorder, *pgwire.Server, *pgproto3.Frontend, *syncBuffer) {
	t.Helper()
	fake := &fakeRecorder{server: make(chan []model.Response, 1), zugestelltEndet: zugestelltEndet, antworten: rs}
	log := &syncBuffer{}
	client, s := verbindeMitLog(t, fake, log)
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	if erweitert {
		fe.Send(&pgproto3.Parse{Name: "s", Query: "SELECT 1"})
		fe.Send(&pgproto3.Sync{})
		_ = fe.Flush()
		warteAuf(t, fake, "c:sync")
		fake.server <- rs
		return fake, s, fe, log
	}
	fe.Send(&pgproto3.Query{String: "SELECT 1"})
	fe.Send(&pgproto3.Query{String: "SELECT 2"})
	fe.Send(&pgproto3.Terminate{})
	_ = fe.Flush()
	return fake, s, fe, log
}

// Abdeckung: LH-FA-25/Happy, LH-FA-25/Negative — Record im Adapter, einfache Anfrage und
// Extended-Interaktion, auch beim Herunterfahren: Liefert der Use Case Antworten mit
// connection_closed als letzter, sendet der Adapter die Nachrichten davor, danach
// weder ReadyForQuery noch eine Meldung noch eine weitere Fehlerantwort, meldet sie
// als zugestellt, beendet die Session mit EndServer (ohne Fehler und ohne Eintrag im
// Log) und schließt die Verbindung; eine schon gesendete weitere Anfrage erreicht den
// Use Case nicht (LH-FA-25.a *Meldung des Endes beim Aufzeichnen*).
func TestRecordEndeSendetDavorUndSchliesst(t *testing.T) {
	for _, erweitert := range []bool{false, true} {
		for _, herunterfahren := range []bool{false, true} {
			t.Run(fmt.Sprintf("erweitert %v, herunterfahren %v", erweitert, herunterfahren), func(t *testing.T) {
				erste, want := model.Response{Type: model.ResponseCommandComplete, Tag: "SELECT 1"}, "CommandComplete,ErrorResponse"
				if erweitert {
					erste, want = model.Response{Type: model.ResponseParseComplete}, "*pgproto3.ParseComplete,ErrorResponse"
				}
				rec, s, fe, log := recordWeg(t, erweitert, recordEndeFolge(erste), herunterfahren)
				typen := bisEndeGelesen(t, fe)
				if got := strings.Join(typen, ","); got != want {
					t.Fatalf("Nachrichten bis zum Ende %q, erwartet %q", got, want)
				}
				if end := rec.lastEnd(t); end != model.EndServer {
					t.Fatalf("Session-Ende %v, erwartet EndServer", end)
				}
				ablauf := rec.protokoll()
				if !strings.HasSuffix(ablauf, fmt.Sprintf("zugestellt ende:%d", model.EndServer)) || strings.Count(ablauf, "q:") > 1 {
					t.Fatalf("Ablauf %q, erwartet Zustellung vor dem Ende und höchstens eine Anfrage", ablauf)
				}
				if text := log.String(); strings.Contains(text, "level=ERROR") || strings.Contains(text, "level=WARN") || s.FirstErrorCode() != "" {
					t.Fatalf("Meldung beim Ende (erster Fehler %q): %s", s.FirstErrorCode(), text)
				}
			})
		}
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-02/Negative — Record im Adapter, einfache Anfrage:
// Scheitert das Senden der Antworten mit Ende (der Client hat die Verbindung
// geschlossen), meldet der Adapter sie nicht als zugestellt und beendet die Session mit
// EndWriteFailed, nicht mit EndServer; den Fehler, den der Use Case daraus ableitet,
// merkt er (LH-FA-25.a *Meldung des Endes beim Aufzeichnen*).
func TestRecordEndeSchreibfehler(t *testing.T) {
	fake := &fakeRecorder{closeErr: model.Errorf(model.CodeConnectionLost, nil, "Antwort nicht zugestellt"),
		antworten: recordEndeFolge(model.Response{Type: model.ResponseCommandComplete, Tag: "SELECT 1"})}
	client, s := verbinde(t, fake)
	fe := pgproto3.NewFrontend(client, client)
	startup(t, fe)
	fe.Send(&pgproto3.Query{String: "SELECT 1"})
	_ = fe.Flush()
	client.Close()
	if end := fake.lastEnd(t); end != model.EndWriteFailed || strings.Contains(fake.protokoll(), "zugestellt") {
		t.Fatalf("Session-Ende %v, Ablauf %q, erwartet EndWriteFailed ohne Zustellung", end, fake.protokoll())
	}
	for frist := time.Now().Add(2 * time.Second); s.FirstErrorCode() == "" && time.Now().Before(frist); {
		time.Sleep(5 * time.Millisecond)
	}
	if s.FirstErrorCode() != model.CodeConnectionLost {
		t.Fatalf("erster Fehler %q, erwartet %s", s.FirstErrorCode(), model.CodeConnectionLost)
	}
}
