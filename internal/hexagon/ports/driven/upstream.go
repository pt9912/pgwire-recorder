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
	Close() error
}
