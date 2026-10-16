package driven

import (
	"context"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Upstream ist der PostgreSQL-Upstream-Port (ARC-004).
type Upstream interface {
	// Open baut eine Verbindung zum Server auf und meldet sich mit den
	// Startup-Parametern an. Verlangt der Server eine Anmeldung, gibt Open jede
	// Anforderung über anmeldung an den Client und dessen Antwort an den Server
	// weiter; anmeldung ist nicht nil. Die Aufrufe von anmeldung kommen
	// nacheinander, nie gleichzeitig (LH-FA-05.b *Aufruf-Voraussetzung*); liefert
	// sie einen Fehler, endet Open mit ihm und schließt die Verbindung zum Server
	// ohne Terminate. Die Antworten des Aufbaus gehen zurück; endet der Aufbau
	// mit einer Fehlerantwort, ist die Session nil.
	Open(ctx context.Context, startup map[string]string, anmeldung model.Anmeldung) (UpstreamSession, []model.Response, error)
}

// UpstreamSession ist eine offene Verbindung zum Server.
//
// Gleichzeitigkeit: Send und Receive dürfen gleichzeitig laufen; jedes
// blockiert nur an seiner Richtung (Send, solange der Server nicht liest;
// Receive, solange er nichts sendet). Query läuft nie gleichzeitig mit Send oder
// Receive, und von jeder Operation läuft höchstens ein Aufruf. Close darf
// jederzeit laufen, blockiert nicht an einem Server, der nicht liest, und
// beendet ein wartendes Send, Receive oder Query mit einem Fehler.
//
// Ende des Servers (LH-FA-25.a): Liest die Session in der laufenden Interaktion
// eine Fehlerantwort FATAL oder PANIC, endet die Antwortfolge von Query und
// Receive mit dieser Fehlerantwort und danach connection_closed statt mit
// ReadyForQuery; die Session wartet nicht auf das Verbindungsende und liest
// nichts dahinter. Send liefert danach keinen Fehler, auch wenn das Schreiben
// scheitert, und liest nichts; ohne gelesene Fehlerantwort ist ein Sendefehler
// PGR-E4003, nach einer anderen gelesenen Fehlerantwort PGR-E6001 (LH-FA-02.b).
// Was Query und Receive nach connection_closed liefern, ist nicht zugesagt.
type UpstreamSession interface {
	// Query sendet eine einfache Anfrage und liefert die Antworten bis
	// einschließlich ReadyForQuery oder, bei einem Ende des Servers, bis
	// einschließlich connection_closed.
	Query(ctx context.Context, sql string) ([]model.Response, error)
	// Send sendet die Client-Nachrichten einer Gruppe einer Extended-Interaktion
	// unverändert und in ihrer Reihenfolge an den Server.
	Send(ctx context.Context, msgs []model.ClientMessage) error
	// Receive wartet auf die nächste Server-Nachricht einer Extended-Interaktion
	// und liefert sie zusammen mit den Nachrichten, die schon dahinter empfangen
	// sind, höchstens bis einschließlich ReadyForQuery oder connection_closed.
	Receive(ctx context.Context) ([]model.Response, error)
	Close() error
}
