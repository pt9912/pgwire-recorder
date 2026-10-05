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
type UpstreamSession interface {
	// Query sendet eine einfache Anfrage und liefert die Antworten bis
	// einschließlich ReadyForQuery.
	Query(ctx context.Context, sql string) ([]model.Response, error)
	// Send sendet die Client-Nachrichten einer Gruppe einer Extended-Interaktion
	// unverändert und in ihrer Reihenfolge an den Server.
	Send(ctx context.Context, msgs []model.ClientMessage) error
	// Receive wartet auf die nächste Server-Nachricht einer Extended-Interaktion
	// und liefert sie zusammen mit den Nachrichten, die schon dahinter empfangen
	// sind, höchstens bis einschließlich ReadyForQuery. Receive darf gleichzeitig
	// mit Send und Close laufen, nicht mit Query oder einem zweiten Receive;
	// Close beendet es mit einem Fehler.
	Receive(ctx context.Context) ([]model.Response, error)
	Close() error
}
