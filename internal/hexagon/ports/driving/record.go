package driving

import (
	"context"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Recorder ist der Record-Use-Case (ARC-003): er vermittelt Client-Sessions zum
// Upstream und zeichnet ihre Interaktionen auf. Die Signaturen tragen nur
// Typen des Domain Models, damit der Record-Service den Port erfüllt, ohne ihn
// zu importieren (architecture.md §2).
type Recorder interface {
	// OpenSession baut für eine Client-Verbindung die Upstream-Session auf.
	// Die Antworten des Verbindungsaufbaus gehen unverändert an den Client;
	// endet der Aufbau mit einer Fehlerantwort, ist die Kennung 0.
	OpenSession(ctx context.Context, startup map[string]string) (model.SessionID, []model.Response, error)
	// Query leitet eine einfache Anfrage weiter und liefert die Serverantworten
	// bis einschließlich ReadyForQuery. Die Interaktion gilt als abgeschlossen;
	// erreicht ihre Antwort den Client nicht, endet die Session mit EndLost.
	Query(ctx context.Context, id model.SessionID, sql string) ([]model.Response, error)
	// CloseSession beendet die Session aus dem genannten Grund (model.SessionEnd).
	CloseSession(ctx context.Context, id model.SessionID, end model.SessionEnd) error
}
