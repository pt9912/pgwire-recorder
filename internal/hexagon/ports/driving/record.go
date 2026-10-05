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

	// ClientMessage übernimmt eine Client-Nachricht einer Extended-Interaktion
	// in Ankunftsreihenfolge (LH-FA-18.a); mit Flush oder Sync geht ihre Gruppe
	// an den Upstream. Nach einem Sync nimmt die Session erst wieder eine
	// Client-Nachricht an, wenn ServerMessage die Interaktion abgeschlossen hat.
	ClientMessage(ctx context.Context, id model.SessionID, m model.ClientMessage) error
	// AwaitServer wartet auf die nächsten Server-Nachrichten einer
	// Extended-Interaktion; es ändert die Aufzeichnung nicht und darf
	// gleichzeitig mit ClientMessage laufen. Jede gelieferte Nachricht geht
	// danach einzeln an ServerMessage, bevor sie an den Client geht.
	AwaitServer(ctx context.Context, id model.SessionID) ([]model.Response, error)
	// ServerMessage ordnet eine Server-Nachricht der Gruppe zu, die gerade
	// Antworten aufnimmt. abgeschlossen meldet das ReadyForQuery, mit dem die
	// Interaktion Teil der Session wurde; erreicht es den Client nicht, endet die
	// Session mit EndLost.
	ServerMessage(ctx context.Context, id model.SessionID, r model.Response) (abgeschlossen bool, err error)

	// CloseSession beendet die Session aus dem genannten Grund (model.SessionEnd).
	// Eine nicht abgeschlossene Extended-Interaktion wird nie übernommen.
	CloseSession(ctx context.Context, id model.SessionID, end model.SessionEnd) error
}
