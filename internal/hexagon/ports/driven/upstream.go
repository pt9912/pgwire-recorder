package driven

import (
	"context"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Upstream ist der PostgreSQL-Upstream-Port (ARC-004).
type Upstream interface {
	// Open baut eine Verbindung zum Server auf und meldet sich mit den
	// Startup-Parametern an. Die Antworten des Aufbaus gehen zurück; endet der
	// Aufbau mit einer Fehlerantwort, ist die Session nil.
	Open(ctx context.Context, startup map[string]string) (UpstreamSession, []model.Response, error)
}

// UpstreamSession ist eine offene Verbindung zum Server.
//
// Gleichzeitigkeit: Send und Receive dürfen gleichzeitig laufen; jedes
// blockiert nur an seiner Richtung (Send, solange der Server nicht liest;
// Receive, solange er nichts sendet). Query läuft nie gleichzeitig mit Send oder
// Receive, und von jeder Operation läuft höchstens ein Aufruf. Close darf
// jederzeit laufen, blockiert nicht an einem Server, der nicht liest, und
// beendet ein wartendes Send, Receive oder Query mit einem Fehler.
type UpstreamSession interface {
	// Query sendet eine einfache Anfrage und liefert die Antworten bis
	// einschließlich ReadyForQuery.
	Query(ctx context.Context, sql string) ([]model.Response, error)
	// Send sendet die Client-Nachrichten einer Gruppe einer Extended-Interaktion
	// unverändert und in ihrer Reihenfolge an den Server.
	Send(ctx context.Context, msgs []model.ClientMessage) error
	// Receive wartet auf die nächste Server-Nachricht einer Extended-Interaktion
	// und liefert sie zusammen mit den Nachrichten, die schon dahinter empfangen
	// sind, höchstens bis einschließlich ReadyForQuery.
	Receive(ctx context.Context) ([]model.Response, error)
	Close() error
}
