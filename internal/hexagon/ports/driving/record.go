package driving

import (
	"context"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Recorder ist der Record-Use-Case (ARC-003): er vermittelt Client-Sessions zum
// Upstream und zeichnet ihre Interaktionen auf. Die Signaturen tragen nur
// Typen des Domain Models, damit der Record-Service den Port erfüllt, ohne ihn
// zu importieren (architecture.md §2).
//
// Der Driving Adapter betreibt je Session zwei Richtungen: Client → Server ruft
// Query, ClientMessage und Delivered; Server → Client ruft AwaitServer und
// Delivered. Aufrufe der beiden Richtungen dürfen für dieselbe Session
// gleichzeitig laufen. Den Zustand der Interaktionen führt allein der Use Case;
// der Adapter meldet Ereignisse und führt aus, was zurückkommt. CloseSession
// beendet jeden wartenden Aufruf der Session mit model.ErrSessionEnded, und
// jeder spätere Aufruf liefert ihn ebenso, auch wenn sein Upstream-Aufruf nach
// dem Ende erfolgreich zurückkam.
//
// Nach Shutdown beginnt keine neue Interaktion: Query und ClientMessage liefern
// für eine Nachricht, die eine neue begänne, model.ErrShutdown und leiten sie
// nicht weiter. Der Adapter liest dann nicht weiter und beendet die Session mit
// EndShutdown, sobald Shutdown oder Delivered das Ende freigibt.
type Recorder interface {
	// OpenSession baut für eine Client-Verbindung die Upstream-Session auf.
	// Die Antworten des Verbindungsaufbaus gehen unverändert an den Client;
	// endet der Aufbau mit einer Fehlerantwort, ist die Kennung 0.
	OpenSession(ctx context.Context, startup map[string]string) (model.SessionID, []model.Response, error)
	// Query leitet eine einfache Anfrage weiter und liefert die Serverantworten
	// bis einschließlich ReadyForQuery. Laufen Extended-Interaktionen, deren
	// Sync schon gesendet ist, wartet Query, bis ihre Antworten zugestellt sind.
	// Danach gehen die Antworten an den Client, und der Adapter meldet Delivered.
	Query(ctx context.Context, id model.SessionID, sql string) ([]model.Response, error)
	// ClientMessage übernimmt eine Client-Nachricht einer Extended-Interaktion
	// in Ankunftsreihenfolge (LH-FA-18.a); mit Flush oder Sync sendet es ihre
	// Gruppe an den Upstream und blockiert dabei nur, solange der Upstream nicht
	// liest.
	ClientMessage(ctx context.Context, id model.SessionID, m model.ClientMessage) error
	// AwaitServer wartet, bis eine Gruppe gesendet ist, liest die nächsten
	// Server-Nachrichten vom Upstream und ordnet sie ihren Gruppen zu. Der
	// Adapter schreibt sie an den Client und meldet danach Delivered, bevor er
	// AwaitServer erneut ruft.
	AwaitServer(ctx context.Context, id model.SessionID) ([]model.Response, error)
	// Delivered meldet, dass die zuletzt gelieferten Antworten (aus Query oder
	// AwaitServer) beim Client sind. endet sagt, dass die Session herunterfährt
	// und keine Interaktion mehr läuft: der Adapter beendet sie dann mit
	// EndShutdown.
	Delivered(ctx context.Context, id model.SessionID) (endet bool)
	// Shutdown meldet das Herunterfahren; die Client-Richtung ruft es vor dem
	// Lesen der nächsten Nachricht, damit eine schon gelesene noch verarbeitet
	// wird (LH-FA-13.a). endet sagt, dass keine Interaktion läuft und keine
	// Antwort auf Zustellung wartet; der Adapter beendet die Session dann mit
	// EndShutdown. Sonst gibt Delivered das Ende später frei.
	Shutdown(ctx context.Context, id model.SessionID) (endet bool)
	// CloseSession beendet die Session mit dem gemeldeten Ereignis
	// (model.SessionEnd) und schließt den Upstream. Der Fehler nennt einen
	// Verbindungsfehler, den der Use Case aus dem Ereignis ableitet (etwa
	// PGR-E4003, beim Zwangsende model.EndForced PGR-E4006), und einen Fehler
	// beim Schreiben der Aufzeichnung.
	CloseSession(ctx context.Context, id model.SessionID, end model.SessionEnd) error
}
