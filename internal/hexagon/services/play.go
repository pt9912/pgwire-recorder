package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// PlayService erfüllt den Play-Use-Case (ARC-002, LH-FA-20.a): Er spielt die
// einfachen Anfragen einer Aufzeichnung Session für Session über eigene
// Verbindungen gegen einen Server ein und bricht beim ersten Fehler ab.
type PlayService struct {
	ziel     driven.Einspielziel
	optionen PlayOptions
	// sessions sind die Sessions mit Interaktionen in der Reihenfolge der
	// Aufzeichnung.
	sessions []model.Session
}

// PlayOptions sind die Optionen von play, die der Play-Service auswertet
// (LH-FA-20.a): User und Database ersetzen, wenn nicht leer, die
// Startup-Parameter user und database jeder Session (*Startup-Daten*).
//
// Kopplung: Der CLI-Adapter liest die Optionen in einen Typ mit denselben
// Feldern in derselben Reihenfolge, und der Bootstrap konvertiert ihn in
// diesen; eine Option kommt in beiden Typen hinzu, sonst baut der Bootstrap
// nicht.
type PlayOptions struct {
	User     string
	Database string
}

// NewPlayService ist der Start des Einspielens (LH-FA-20.a *Start*): Es lädt
// die Aufzeichnung; ein Fehler dabei und eine Interaktion, die Validate nicht
// besteht (PGR-E3003), gehen vor. Danach ist die erste Interaktion, die keine
// einfache Anfrage ist, in der Reihenfolge der Sessions und ihrer
// Interaktionen, PGR-E6001 mit Session, Nummer und Art (*Art der
// Interaktion*). o sind die Optionen des Einspielens.
func NewPlayService(ctx context.Context, repo driven.RecordingRepository, path string, ziel driven.Einspielziel, o PlayOptions) (*PlayService, error) {
	rec, err := repo.Load(ctx, path)
	if err != nil {
		return nil, err
	}
	for _, sess := range rec.Sessions {
		for _, in := range sess.Interactions {
			if err := in.Validate(); err != nil {
				return nil, model.Errorf(model.CodeRecordingBroken, err, "%s: Session %d, Interaktion %d", path, sess.ID, in.Sequence)
			}
		}
	}
	s := &PlayService{ziel: ziel, optionen: o}
	for _, sess := range rec.Sessions {
		for _, in := range sess.Interactions {
			if in.Request.Type != model.RequestQuery {
				return nil, model.Errorf(model.CodeUnsupported, nil, "Session %d, Interaktion %d: eine Interaktion der Art %s spielt play nicht ein", sess.ID, in.Sequence, in.Request.Type)
			}
		}
		if len(sess.Interactions) > 0 {
			s.sessions = append(s.sessions, sess)
		}
	}
	return s, nil
}

// Play spielt die Sessions ein (LH-FA-20.a). Endet ctx (erstes
// Abbruchsignal), beginnt keine weitere Interaktion und keine weitere
// Session; ein laufender Aufbau läuft zu Ende, eine laufende Interaktion bis
// zu ihrem ReadyForQuery. Wird ablauf geschlossen (zweites Signal), schließt
// es die laufende Verbindung sofort; der unterbrochene Aufbau oder die
// unterbrochene Interaktion ist kein Fehler. Play liefert den ersten Fehler,
// nach dem es abgebrochen hat, sonst nil; ein nil-Kanal ablauf wird nie
// geschlossen.
func (s *PlayService) Play(ctx context.Context, ablauf <-chan struct{}) error {
	abbruch, abbrechen := context.WithCancel(context.WithoutCancel(ctx))
	defer abbrechen()
	go func() {
		select {
		case <-ablauf:
			abbrechen()
		case <-abbruch.Done():
		}
	}()
	for _, sess := range s.sessions {
		if ctx.Err() != nil {
			return nil
		}
		err := s.session(ctx, abbruch, sess)
		if abbruch.Err() != nil {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// session spielt eine Session über eine eigene Verbindung ein. Am Ende,
// regulär oder nach einem Fehler, sendet sie Terminate und schließt die
// Verbindung (Schliesse); endet abbruch, schließt sie sofort.
func (s *PlayService) session(ctx, abbruch context.Context, sess model.Session) error {
	us, err := s.ziel.Verbinde(abbruch, s.startup(sess))
	if err != nil {
		return eingeordnet(err, "Session %d", sess.ID)
	}
	stop := context.AfterFunc(abbruch, us.Schliesse)
	defer func() {
		if stop() {
			us.Schliesse()
		}
	}()
	for _, in := range sess.Interactions {
		if ctx.Err() != nil {
			return nil
		}
		if err := interaktion(us, in); err != nil {
			return eingeordnet(err, "Session %d, Interaktion %d", sess.ID, in.Sequence)
		}
	}
	return nil
}

// startup sind die aufgezeichneten Startup-Parameter der Session, user und
// database ersetzt durch User und Database der Optionen, wenn diese nicht
// leer sind (LH-FA-20.a *Startup-Daten*).
func (s *PlayService) startup(sess model.Session) map[string]string {
	out := make(map[string]string, len(sess.Startup)+2)
	for k, v := range sess.Startup {
		out[k] = v
	}
	if s.optionen.User != "" {
		out["user"] = s.optionen.User
	}
	if s.optionen.Database != "" {
		out["database"] = s.optionen.Database
	}
	return out
}

// interaktion sendet die Anfrage und liest bis zum ReadyForQuery. Eine
// Fehlerantwort beendet das Lesen: mit dem Schweregrad FATAL oder PANIC
// (Feld V, ohne es S) PGR-E4003, sonst PGR-E4004; die Meldung nennt von der
// Fehlerantwort SQLSTATE und Meldung, keine weiteren Felder (LH-FA-20.a
// *Interaktion*, *Meldungen*). Jede andere Antwort verwirft sie.
func interaktion(us driven.EinspielSession, in model.Interaction) error {
	if err := us.Anfrage(in.Request.SQL); err != nil {
		return err
	}
	for {
		r, err := us.Naechste()
		if err != nil {
			return err
		}
		switch r.Type {
		case model.ResponseReadyForQuery:
			return nil
		case model.ResponseErrorResponse:
			schwere := r.Fields["V"]
			if schwere == "" {
				schwere = r.Fields["S"]
			}
			if schwere == "FATAL" || schwere == "PANIC" {
				return model.Errorf(model.CodeConnectionLost, nil, "Fehlerantwort des Servers %s „%s“, die Verbindung endet", r.Fields["C"], r.Fields["M"])
			}
			return model.Errorf(model.CodeServerError, nil, "Fehlerantwort des Servers %s „%s“", r.Fields["C"], r.Fields["M"])
		}
	}
}

// eingeordnet stellt den Ort vor einen Fehler und behält dessen Code; ein
// Fehler ohne Code bleibt ohne.
func eingeordnet(err error, format string, args ...any) error {
	var me *model.Error
	if errors.As(err, &me) {
		return model.Errorf(me.Code, err, format, args...)
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}
