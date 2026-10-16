package pgwire_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// endeReplayer liefert auf jede Anfrage antworten und zählt die Aufrufe von
// Query und ClientMessage; Sent und CloseConnection protokolliert der
// eingebettete fakeReplayer.
type endeReplayer struct {
	*fakeReplayer
	antworten []model.Response
	anfragen  int
}

func (e *endeReplayer) Query(context.Context, model.SessionID, string) ([]model.Response, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.anfragen++
	e.ereignisse = append(e.ereignisse, "query")
	return e.antworten, nil
}

func (e *endeReplayer) ClientMessage(_ context.Context, id model.SessionID, m model.ClientMessage) ([]model.Response, error) {
	e.mu.Lock()
	e.anfragen++
	e.mu.Unlock()
	if m.Type == model.ClientSync {
		return e.antworten, nil
	}
	return nil, nil
}

func (e *endeReplayer) Sent(ctx context.Context, id model.SessionID) {
	e.fakeReplayer.Sent(ctx, id)
	e.mu.Lock()
	defer e.mu.Unlock()
	e.ereignisse = append(e.ereignisse, "sent")
}

func (e *endeReplayer) zahl() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.anfragen
}

// endeFolge sind die Antworten eines Endes: die Fehlerantwort des Servers, danach
// connection_closed.
func endeFolge() []model.Response {
	return []model.Response{
		{Type: model.ResponseCommandComplete, Tag: "SELECT 1"},
		{Type: model.ResponseErrorResponse, Fields: map[string]string{"S": "FATAL", "V": "FATAL", "C": "57P01", "M": "terminating connection"}},
		{Type: model.ResponseConnectionClosed},
	}
}

// bisEndeGelesen liest vom Client, bis die Verbindung endet, und liefert die
// Nachrichten davor als Typnamen; Fehler außer dem Ende der Verbindung machen
// den Test rot. Die Zahl der Nachrichten bis zum Ende zählt der Aufrufer.
func bisEndeGelesen(t *testing.T, fe *pgproto3.Frontend) []string {
	t.Helper()
	var typen []string
	for {
		msg, err := fe.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.ErrClosedPipe) {
				return typen
			}
			t.Fatalf("Lesen bis zum Ende: %v", err)
		}
		typen = append(typen, typName(msg))
	}
}

// typName nennt die Art einer Server-Nachricht für die Prüfung der Folge.
func typName(m pgproto3.BackendMessage) string {
	switch m.(type) {
	case *pgproto3.CommandComplete:
		return "CommandComplete"
	case *pgproto3.ErrorResponse:
		return "ErrorResponse"
	case *pgproto3.ReadyForQuery:
		return "ReadyForQuery"
	}
	return fmt.Sprintf("%T", m)
}

// Abdeckung: LH-FA-25/Happy, LH-FA-25/Negative — Replay im Adapter: Endet die
// Folge der Antworten mit connection_closed, sendet der Adapter die
// Nachrichten davor, danach weder ReadyForQuery noch eine Meldung noch eine
// weitere Fehlerantwort, und schließt die Verbindung (der Client liest bis zum
// Ende); ein PGWire-Byte für connection_closed gibt es nicht (LH-FA-25.a
// *Meldung des Endes im Replay*). Gilt für Query wie für den abschließenden Sync.
func TestReplayEndeSendetDavorUndSchliesst(t *testing.T) {
	for _, weg := range []string{"query", "sync"} {
		rep := &endeReplayer{fakeReplayer: &fakeReplayer{}, antworten: endeFolge()}
		s, _, fe, log, fertig := replayVerbindung(t, rep)
		if weg == "query" {
			fe.Send(&pgproto3.Query{String: "SELECT 1"})
		} else {
			fe.Send(&pgproto3.Parse{Name: "s", Query: "SELECT 1"})
			fe.Send(&pgproto3.Sync{})
		}
		_ = fe.Flush()
		typen := bisEndeGelesen(t, fe)
		warte(t, fertig)
		if got := strings.Join(typen, ","); got != "CommandComplete,ErrorResponse" {
			t.Fatalf("%s: Nachrichten bis zum Ende %q, erwartet CommandComplete,ErrorResponse", weg, got)
		}
		if text := log.String(); strings.Contains(text, "level=ERROR") || strings.Contains(text, "level=WARN") || s.FirstErrorCode() != "" {
			t.Fatalf("%s: Meldung beim Ende (erster Fehler %q): %s", weg, s.FirstErrorCode(), text)
		}
	}
}

// Abdeckung: LH-FA-25/Negative — Replay im Adapter: Nach dem Ende liest der
// Adapter nichts mehr. Eine Nachricht, die der Client schon vor dem Ende der
// Verbindung gesendet hat, erreicht den Use Case nicht (kein zweiter Aufruf von
// Query oder ClientMessage); ein Terminate danach ändert nichts.
func TestReplayEndeLiestNichtsMehr(t *testing.T) {
	rep := &endeReplayer{fakeReplayer: &fakeReplayer{}, antworten: endeFolge()}
	_, _, fe, _, fertig := replayVerbindung(t, rep)
	fe.Send(&pgproto3.Query{String: "SELECT 1"})
	fe.Send(&pgproto3.Query{String: "SELECT 2"})
	fe.Send(&pgproto3.Terminate{})
	_ = fe.Flush()
	_ = bisEndeGelesen(t, fe)
	warte(t, fertig)
	if n := rep.zahl(); n != 1 {
		t.Fatalf("Use Case %d-mal gefragt, erwartet 1", n)
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-03/Boundary — Replay im Adapter: Der
// Adapter meldet die Antworten vor dem Ende als gesendet (Sent) und schließt
// danach (CloseConnection): Reihenfolge query, sent, close, je einmal.
func TestReplayEndeSentVorSchliessen(t *testing.T) {
	rep := &endeReplayer{fakeReplayer: &fakeReplayer{}, antworten: endeFolge()}
	_, _, fe, _, fertig := replayVerbindung(t, rep)
	fe.Send(&pgproto3.Query{String: "SELECT 1"})
	_ = fe.Flush()
	_ = bisEndeGelesen(t, fe)
	warte(t, fertig)
	if got := rep.protokoll(); got != "query sent close" {
		t.Fatalf("Ablauf %q, erwartet \"query sent close\"", got)
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-03/Boundary — Replay im Adapter: Scheitert
// das Senden der Antworten vor dem Ende (Schreibfehler, der Client hat die
// Verbindung geschlossen), meldet der Adapter sie nicht als gesendet (die
// Interaktion bleibt unverbraucht) und merkt PGR-E4003 wie bei jeder Interaktion.
func TestReplayEndeSchreibfehlerBleibtUnverbraucht(t *testing.T) {
	rep := &endeReplayer{fakeReplayer: &fakeReplayer{}, antworten: endeFolge()}
	s, client, fe, _, fertig := replayVerbindung(t, rep)
	fe.Send(&pgproto3.Query{String: "SELECT 1"})
	_ = fe.Flush()
	client.Close()
	warte(t, fertig)
	if n := rep.sentAufrufe(); n != 0 {
		t.Fatalf("Sent nach gescheitertem Senden gemeldet (%d)", n)
	}
	if s.FirstErrorCode() != model.CodeConnectionLost {
		t.Fatalf("erster Fehler %q, erwartet %s", s.FirstErrorCode(), model.CodeConnectionLost)
	}
}

// Abdeckung: LH-FA-25/Negative — Replay im Adapter: connection_closed ist nur
// als letzte Antwort eine Auslassung; an anderer Stelle ist sie wie jeder Typ
// ohne PGWire-Nachricht ein interner Fehler (kein Senden des Rests, kein Sent).
func TestReplayEndeNurAlsLetzteAntwort(t *testing.T) {
	mitte := []model.Response{{Type: model.ResponseConnectionClosed}, {Type: model.ResponseCommandComplete, Tag: "SELECT 1"}}
	rep := &endeReplayer{fakeReplayer: &fakeReplayer{}, antworten: mitte}
	s, _, fe, _, fertig := replayVerbindung(t, rep)
	fe.Send(&pgproto3.Query{String: "SELECT 1"})
	_ = fe.Flush()
	_ = bisEndeGelesen(t, fe)
	warte(t, fertig)
	if n := rep.sentAufrufe(); n != 0 {
		t.Fatalf("Sent bei nicht abbildbarer Antwort gemeldet (%d)", n)
	}
	if s.FirstErrorCode() != model.CodeInternal {
		t.Fatalf("erster Fehler %q, erwartet %s", s.FirstErrorCode(), model.CodeInternal)
	}
}
