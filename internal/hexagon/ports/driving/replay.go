package driving

import (
	"context"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Replayer ist der Replay-Use-Case (ARC-003): er beantwortet Anfragen aus einer
// Aufzeichnung, ohne Upstream. Die Signaturen tragen nur Typen des Domain
// Models (architecture.md §2).
type Replayer interface {
	// OpenConnection nimmt eine Client-Verbindung an und liefert die Antworten
	// des Verbindungsaufbaus (Serverparameter und ReadyForQuery) nach
	// LH-FA-12.a.
	OpenConnection(ctx context.Context) (model.SessionID, []model.Response)
	// Query ordnet der Verbindung bei ihrer ersten Anfrage eine Session zu und
	// liefert die aufgezeichneten Antworten der Interaktion am Cursor; eine
	// abweichende Anfrage ist PGR-E5001, eine Anfrage ohne freie Session
	// PGR-E5003 (LH-FA-09.a, LH-FA-10.a).
	Query(ctx context.Context, id model.SessionID, sql string) ([]model.Response, error)
	// ClientMessage ordnet wie Query zu und vergleicht eine Client-Nachricht
	// des Extended Query Protocol mit der erwarteten am Cursor. Schließt sie
	// ihre Gruppe ab (Flush oder Sync), liefert es die aufgezeichneten
	// Server-Nachrichten der Gruppe, vorher keine; eine abweichende Nachricht
	// ist PGR-E5001 (LH-FA-18.a).
	ClientMessage(ctx context.Context, id model.SessionID, m model.ClientMessage) ([]model.Response, error)
	// Shutdown fragt, ob die Verbindung beim Herunterfahren enden darf: true,
	// wenn keine Extended-Interaktion läuft, deren Sync noch aussteht
	// (LH-FA-13.a). Es ändert keinen Zustand und darf mehrfach gerufen werden;
	// dass danach keine neue Interaktion beginnt, hält der Aufrufer ein.
	Shutdown(ctx context.Context, id model.SessionID) bool
	// CloseConnection beendet die Verbindung. Bleiben Interaktionen der
	// zugeordneten Session unverbraucht, liefert sie die Warnung PGR-W2001.
	CloseConnection(ctx context.Context, id model.SessionID) *model.Warning
}
