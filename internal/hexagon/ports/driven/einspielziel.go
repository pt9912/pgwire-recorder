package driven

import (
	"context"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Einspielziel ist der PostgreSQL-Server beim Einspielen (ARC-004,
// LH-FA-20.a): Der Recorder ist dort Client.
type Einspielziel interface {
	// Verbinde baut die Verbindung einer Session mit den Startup-Parametern
	// startup bis zum ersten ReadyForQuery auf. Ein Fehler ist nach
	// LH-FA-20.a *Aufbau* eingestuft (PGR-E4002, PGR-E4005); die Verbindung ist
	// dann geschlossen, ohne dass danach noch etwas gesendet wurde. Endet ctx
	// während des Aufbaus, bricht er ab: Die Verbindung wird ohne Terminate
	// geschlossen, und Verbinde liefert einen Fehler.
	Verbinde(ctx context.Context, startup map[string]string) (EinspielSession, error)
}

// EinspielSession ist eine aufgebaute Verbindung zum Server beim Einspielen.
//
// Gleichzeitigkeit: Anfrage und Naechste laufen nacheinander. Gruppe beginnt
// das Senden und kehrt zurück; das Senden der Gruppen läuft in der Session
// unabhängig davon, ob Naechste liest. Schliesse darf jederzeit laufen, auch
// während Anfrage sendet oder Gruppe sendet oder Naechste wartet, und beendet
// sie mit einem Fehler.
type EinspielSession interface {
	// Anfrage sendet eine einfache Anfrage; scheitert das Senden, ist das
	// PGR-E4003.
	Anfrage(sql string) error
	// Gruppe beginnt das Senden der Client-Nachrichten einer Gruppe einer
	// Extended-Interaktion in ihrer Reihenfolge und kehrt zurück, ohne auf das
	// Ende des Sendens zu warten; mehrere Gruppen sendet die Session nacheinander
	// in der Reihenfolge der Aufrufe. Scheitert das Senden, ist das PGR-E4003 aus
	// Naechste oder der nächsten Gruppe; ein Fehler nach Schliesse bleibt ohne
	// Folge.
	Gruppe(nachrichten []model.ClientMessage) error
	// Naechste liest die nächste Server-Nachricht und liefert sie, wenn das
	// Modell sie kennt; Nachrichten ohne Abbildung im Modell, etwa
	// NotificationResponse, liest und verwirft sie. Eine CopyInResponse,
	// CopyOutResponse oder CopyBothResponse und eine Nachricht, die sich nicht
	// lesen lässt, sind PGR-E6001, ein Verbindungsende PGR-E4003. Sie liest nie
	// über die gelieferte Nachricht hinaus.
	Naechste() (model.Response, error)
	// Schliesse sendet Terminate, soweit die Verbindung es sofort annimmt, und
	// schließt sie; ein Fehler dabei bleibt ohne Folge.
	Schliesse()
}
