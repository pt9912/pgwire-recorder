package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// PlayService erfüllt den Play-Use-Case (ARC-002, LH-FA-20.a): Er spielt die
// einfachen Anfragen und die Extended-Interaktionen einer Aufzeichnung Session
// für Session über eigene Verbindungen gegen einen Server ein und bricht beim
// ersten Fehler ab, der nach PlayOptions nicht weiterläuft.
type PlayService struct {
	ziel     driven.Einspielziel
	optionen PlayOptions
	// sessions sind die Sessions mit Interaktionen in der Reihenfolge der
	// Aufzeichnung.
	sessions []model.Session
	// uhr ist die Uhr der Wartezeiten (LH-FA-21.a *Einspielen*); nil bleibt
	// ohne Wartezeiten und ohne KeepTiming unbenutzt.
	uhr driven.Uhr
	// zaehler sind die Zähler der Zeile zum Ende.
	zaehler Zaehler
}

// PlayOption stellt einen PlayService ein.
type PlayOption func(*PlayService)

// Wartezeiten stellt die Uhr ein, nach der KeepTiming die Pausen der Aufzeichnung
// herstellt (LH-FA-21.a *Einspielen*). Ohne KeepTiming liest der Service die Uhr
// nie; mit KeepTiming ist die Option Pflicht des Aufrufers, ohne sie endet der
// Lauf in einem Nil-Zugriff (nicht zugesagt).
func Wartezeiten(uhr driven.Uhr) PlayOption {
	return func(s *PlayService) { s.uhr = uhr }
}

// PlayOptions sind die Optionen von play, die der Play-Service auswertet
// (LH-FA-20.a): User und Database ersetzen, wenn nicht leer, die
// Startup-Parameter user und database jeder Session (*Startup-Daten*).
// ContinueOnError ist --continue-on-error: Nach einem PGR-E4004 läuft das
// Einspielen weiter. AllowRecordedErrors ist --allow-recorded-errors: Eine
// Fehlerantwort in einer Interaktion, deren Aufzeichnung eine error_response
// trägt, ist erwartet. FinishSessionOnInterrupt ist
// --finish-session-on-interrupt: Nach dem ersten Abbruchsignal laufen die
// übrigen Interaktionen der laufenden Session.
//
// KeepTiming ist --keep-timing: Vor jeder Interaktion wartet der Service nach
// TimingMode (TimingRelative oder TimingAbsolute) und, bei TimingAbsolute,
// TimingReference (ReferenceConnect oder ReferenceFirstRequest) auf die
// aufgezeichnete Zeit (LH-FA-21.a *Einspielen*). Der Service erwartet genau
// diese Werte; ein anderer ist beim Aufrufer ungültig (der CLI-Adapter lehnt ihn
// mit PGR-E2001 ab), und was der Service dann tut, ist nicht zugesagt.
//
// CompareResponses ist --compare-responses: Der Service vergleicht die
// Antworten jeder Interaktion mit der Aufzeichnung (LH-FA-24.a); mit ihm hat
// AllowRecordedErrors keine Wirkung.
//
// Kopplung: Der CLI-Adapter liest die Optionen in einen Typ mit denselben
// Feldern in derselben Reihenfolge, und der Bootstrap konvertiert ihn in
// diesen; eine Option kommt in beiden Typen hinzu, sonst baut der Bootstrap
// nicht.
type PlayOptions struct {
	User                     string
	Database                 string
	ContinueOnError          bool
	AllowRecordedErrors      bool
	FinishSessionOnInterrupt bool
	KeepTiming               bool
	TimingMode               string
	TimingReference          string
	CompareResponses         bool
}

// Werte von PlayOptions.TimingMode und PlayOptions.TimingReference.
const (
	TimingRelative        = "relative"
	TimingAbsolute        = "absolute"
	ReferenceConnect      = "connect"
	ReferenceFirstRequest = "first-request"
)

// Zaehler liefert die Zähler des Laufs für die Zeile zum Ende; der Aufruf folgt
// auf Play (LH-FA-24.a *Meldung*).
func (s *PlayService) Zaehler() Zaehler { return s.zaehler }

// NewPlayService ist der Start des Einspielens (LH-FA-20.a *Start*): Es lädt
// die Aufzeichnung; ein Fehler dabei und eine Interaktion, die Validate nicht
// besteht (PGR-E3003), gehen vor. Jede Interaktion, die Validate besteht, ist
// eine einfache Anfrage oder eine Extended-Interaktion und wird eingespielt.
// Mit KeepTiming prüft er danach, dass jede Interaktion jeder Session offset_ms
// trägt (PGR-E2003, LH-FA-21.a *Randformen des Einspielens*). o sind die
// Optionen des Einspielens.
func NewPlayService(ctx context.Context, repo driven.RecordingRepository, path string, ziel driven.Einspielziel, o PlayOptions, opts ...PlayOption) (*PlayService, error) {
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
	for _, opt := range opts {
		opt(s)
	}
	for _, sess := range rec.Sessions {
		if len(sess.Interactions) > 0 {
			s.sessions = append(s.sessions, sess)
		}
	}
	if o.KeepTiming {
		if err := s.pruefeZeitangaben(); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// Play spielt die Sessions ein (LH-FA-20.a). Endet ctx (erstes
// Abbruchsignal), beginnt keine weitere Session und, ohne
// FinishSessionOnInterrupt, keine weitere Interaktion; ein laufender Aufbau
// läuft zu Ende, eine laufende Interaktion bis zu ihrem ReadyForQuery, eine
// Wartezeit (KeepTiming) endet sofort, und danach wartet Play nicht mehr. Wird
// ablauf geschlossen (zweites Signal), schließt es die laufende Verbindung
// sofort; der unterbrochene Aufbau oder die unterbrochene Interaktion ist kein
// Fehler. Ein nil-Kanal ablauf wird nie geschlossen.
//
// Play liefert die Fehler des Laufs als reine Zusammenfassung (errors.Join):
// den Fehler, nach dem es abgebrochen hat, zuerst, danach je PGR-E4004 und je
// PGR-E5004, nach dem es weiterlief, einen in der Reihenfolge ihres Auftretens,
// auch wenn das zweite Signal die Interaktion danach unterbricht; ohne Fehler
// nil (LH-FA-20.a *Meldungen*, *Exit-Code*).
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
	var frueher []error
	for _, sess := range s.sessions {
		if ctx.Err() != nil {
			break
		}
		err := s.session(ctx, abbruch, sess, &frueher)
		if abbruch.Err() != nil {
			break
		}
		if err != nil {
			return errors.Join(append([]error{err}, frueher...)...)
		}
	}
	return errors.Join(frueher...)
}

// session spielt eine Session über eine eigene Verbindung ein. Am Ende,
// regulär oder nach einem Fehler, sendet sie Terminate und schließt die
// Verbindung (Schliesse); endet abbruch, schließt sie sofort und kehrt erst nach
// dem Schließen zurück. Mit KeepTiming wartet sie vor jeder Interaktion (warte).
// Ein PGR-E4004 oder PGR-E5004, nach dem das Einspielen weiterläuft, hängt sie
// mit seinem Ort an frueher an; die Session endet ohne Fehler, wenn eine
// Fehlerantwort FATAL oder PANIC mit Vergleich weiter gereicht ist
// (errSessionEnde).
func (s *PlayService) session(ctx, abbruch context.Context, sess model.Session, frueher *[]error) error {
	zeit := s.neuerTakt()
	us, err := s.ziel.Verbinde(abbruch, s.startup(sess))
	if err != nil {
		return eingeordnet(err, "Session %d", sess.ID)
	}
	geschlossen := make(chan struct{})
	stop := context.AfterFunc(abbruch, func() {
		us.Schliesse()
		close(geschlossen)
	})
	defer func() {
		if stop() {
			us.Schliesse()
			return
		}
		<-geschlossen
	}()
	for _, in := range sess.Interactions {
		if ctx.Err() != nil && !s.optionen.FinishSessionOnInterrupt {
			return nil
		}
		if !s.warte(ctx, abbruch, zeit, in) {
			return nil
		}
		zeit.beginnt(in)
		ort := func(err error) error {
			return eingeordnet(err, "Session %d, Interaktion %d", sess.ID, in.Sequence)
		}
		weiter := func(err error) { *frueher = append(*frueher, ort(err)) }
		if err := s.interaktion(us, in, weiter); err != nil {
			if errors.Is(err, errSessionEnde) {
				return nil
			}
			return ort(err)
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

// interaktion spielt eine Interaktion ein: eine einfache Anfrage mit anfrage,
// eine Extended-Interaktion mit extended. Ohne CompareResponses verwerfen beide
// jede Antwort außer einer Fehlerantwort. Eine Fehlerantwort mit dem Schweregrad
// FATAL oder PANIC (Feld V, ohne es S) ist PGR-E4003 und bricht ab, auch wenn
// sie erwartet wäre. Jede andere ist mit AllowRecordedErrors erwartet, wenn die
// Aufzeichnung der Interaktion eine error_response trägt, und ohne Meldung;
// sonst ist sie PGR-E4004, mit ContinueOnError an weiter gereicht, ohne ihn der
// Fehler, nach dem die Interaktion ohne weiteres Lesen abbricht. Nach einer
// erwarteten oder weitergereichten Fehlerantwort liest sie weiter bis zum
// ReadyForQuery. Die Meldung nennt von der Fehlerantwort SQLSTATE und Meldung,
// keine weiteren Felder (LH-FA-20.a *Interaktion*, *Meldungen*, Schritt 6).
//
// Mit CompareResponses hält sie die Antworten bis zum ReadyForQuery und
// vergleicht sie dann mit der Aufzeichnung (LH-FA-24.a): Eine Abweichung ist
// PGR-E5004, mit ContinueOnError an weiter gereicht, ohne ihn der Fehler. Eine
// Fehlerantwort FATAL oder PANIC ist sofort die Abweichung PGR-E5004; mit
// ContinueOnError endet dann die Session (errSessionEnde), die nächste läuft.
// Ein PGR-E4004 und ein erwarteter Fehler gibt es dann nicht.
func (s *PlayService) interaktion(us driven.EinspielSession, in model.Interaction, weiter func(error)) error {
	s.zaehler.Eingespielt++
	l := &lauf{us: us, weiter: weiter}
	if s.optionen.CompareResponses {
		l.v = neuerVergleich(in)
	} else {
		l.erwartet = s.optionen.AllowRecordedErrors && mitFehlerantwort(in)
	}
	var err error
	if in.Request.Type == model.RequestExtended {
		err = s.extended(l, in.Groups)
	} else {
		err = s.anfrage(l, in.Request.SQL)
	}
	if err != nil || l.v == nil {
		return err
	}
	return s.abschluss(l)
}

// lauf ist der Stand einer Interaktion beim Lesen der Antworten: die Session,
// ob eine Fehlerantwort erwartet ist (nur ohne Vergleich), die Funktion für
// Fehler, nach denen das Einspielen weiterläuft, und mit Vergleich dessen Stand
// (sonst nil).
type lauf struct {
	us       driven.EinspielSession
	erwartet bool
	weiter   func(error)
	v        *vergleich
}

// naechste liest die nächste Serverantwort; mit Vergleich merkt sie sie.
func (l *lauf) naechste() (model.Response, error) {
	r, err := l.us.Naechste()
	if err == nil && l.v != nil {
		l.v.nimm(r)
	}
	return r, err
}

// anfrage sendet eine einfache Anfrage und liest bis zum ReadyForQuery.
func (s *PlayService) anfrage(l *lauf, sql string) error {
	if err := l.us.Anfrage(sql); err != nil {
		return err
	}
	return s.bisBereit(l)
}

// extended sendet die Gruppen einer Extended-Interaktion in ihrer Reihenfolge
// (LH-FA-20.a *Gruppen*). Nach einer Gruppe mit Flush liest sie, bevor sie die
// nächste sendet, die Antworten auf jede Client-Nachricht der Gruppe
// (erwarteteAntworten); nach der letzten Gruppe, die mit Sync endet, liest sie
// bis zum ReadyForQuery. Nach einer Fehlerantwort, die nicht abbricht, zählt
// sie nicht mehr: Sie sendet die übrigen Gruppen ohne Warten dazwischen und
// liest dann bis zum ReadyForQuery.
func (s *PlayService) extended(l *lauf, gruppen []model.Group) error {
	fehlerGesehen := false
	for _, g := range gruppen {
		if err := l.us.Gruppe(g.Client); err != nil {
			return err
		}
		if fehlerGesehen || g.Client[len(g.Client)-1].Type != model.ClientFlush {
			continue
		}
		for n := erwarteteAntworten(g.Client); n > 0; {
			r, err := l.naechste()
			if err != nil {
				return err
			}
			switch r.Type {
			case model.ResponseDataRow, model.ResponseNoticeResponse, model.ResponseParameterStatus:
				continue
			case model.ResponseErrorResponse:
				if err := s.fehlerantwort(l, r); err != nil {
					return err
				}
				fehlerGesehen = true
				n = 0
			default:
				n--
			}
		}
	}
	return s.bisBereit(l)
}

// erwarteteAntworten ist die Zahl der Antworten, auf die der Server nach einer
// Gruppe mit Flush antwortet: eine je Parse, Bind, Close und Execute und je
// Describe eines Portals, zwei je Describe einer Anweisung
// (ParameterDescription und RowDescription oder NoData), keine auf Flush.
func erwarteteAntworten(nachrichten []model.ClientMessage) int {
	n := 0
	for _, m := range nachrichten {
		switch {
		case m.Type == model.ClientFlush || m.Type == model.ClientSync:
		case m.Type == model.ClientDescribe && m.Target == model.TargetStatement:
			n += 2
		default:
			n++
		}
	}
	return n
}

// bisBereit liest bis zum ReadyForQuery und verwirft jede Antwort außer einer
// Fehlerantwort.
func (s *PlayService) bisBereit(l *lauf) error {
	for {
		r, err := l.naechste()
		if err != nil {
			return err
		}
		switch r.Type {
		case model.ResponseReadyForQuery:
			return nil
		case model.ResponseErrorResponse:
			if err := s.fehlerantwort(l, r); err != nil {
				return err
			}
		}
	}
}

// fehlerantwort stuft eine Fehlerantwort ein. Ohne Vergleich: FATAL und PANIC
// sind PGR-E4003; eine erwartete ist keine Meldung; jede andere ist PGR-E4004,
// mit ContinueOnError an weiter gereicht, sonst der Fehler, den sie liefert. Mit
// Vergleich: FATAL und PANIC sind sofort die Abweichung (fatal), jede andere
// geht in den Vergleich am Ende der Interaktion ein.
func (s *PlayService) fehlerantwort(l *lauf, r model.Response) error {
	schwere := r.Fields["V"]
	if schwere == "" {
		schwere = r.Fields["S"]
	}
	fatal := schwere == "FATAL" || schwere == "PANIC"
	if l.v != nil {
		if fatal {
			return s.fatal(l)
		}
		return nil
	}
	if fatal {
		return model.Errorf(model.CodeConnectionLost, nil, "Fehlerantwort des Servers %s „%s“, die Verbindung endet", r.Fields["C"], r.Fields["M"])
	}
	if l.erwartet {
		return nil
	}
	fehler := model.Errorf(model.CodeServerError, nil, "Fehlerantwort des Servers %s „%s“", r.Fields["C"], r.Fields["M"])
	if !s.optionen.ContinueOnError {
		return fehler
	}
	l.weiter(fehler)
	return nil
}

// errSessionEnde meldet der Session, dass ihre Verbindung durch eine
// Fehlerantwort FATAL oder PANIC endete, deren Abweichung weiter gereicht ist:
// Die Session endet ohne Fehler, die nächste läuft.
var errSessionEnde = errors.New("Session endet")

// fatal stellt die Abweichung einer Fehlerantwort FATAL oder PANIC fest: Die
// Interaktion zählt als verglichen und abweichend. Ohne ContinueOnError ist sie
// der Fehler, der abbricht, mit ihm wird sie weiter gereicht und die Session
// endet (LH-FA-20.a Tabelle, *Interaktion*).
func (s *PlayService) fatal(l *lauf) error {
	s.zaehler.Verglichen++
	s.zaehler.Abweichend++
	fehler := model.Errorf(model.CodeResponseMismatch, nil, "%s", l.v.erste(true).text())
	if !s.optionen.ContinueOnError {
		return fehler
	}
	l.weiter(fehler)
	return errSessionEnde
}

// abschluss vergleicht die Antworten einer Interaktion, die bis zu ihrem
// ReadyForQuery gelesen sind, mit der Aufzeichnung (LH-FA-24.a): Die Interaktion
// zählt als verglichen; weicht sie ab, zählt sie als abweichend, und die erste
// Abweichung ist PGR-E5004, mit ContinueOnError an weiter gereicht, sonst der
// Fehler, der abbricht.
func (s *PlayService) abschluss(l *lauf) error {
	s.zaehler.Verglichen++
	u := l.v.erste(false)
	if u == nil {
		return nil
	}
	s.zaehler.Abweichend++
	fehler := model.Errorf(model.CodeResponseMismatch, nil, "%s", u.text())
	if !s.optionen.ContinueOnError {
		return fehler
	}
	l.weiter(fehler)
	return nil
}

// mitFehlerantwort meldet, ob die Aufzeichnung der Interaktion an irgendeiner
// Stelle eine error_response trägt, in den Antworten einer einfachen Anfrage
// oder in den Server-Nachrichten jeder Gruppe (LH-FA-20.a *Interaktion*).
func mitFehlerantwort(in model.Interaction) bool {
	if trifftFehlerantwort(in.Responses) {
		return true
	}
	for _, g := range in.Groups {
		if trifftFehlerantwort(g.Server) {
			return true
		}
	}
	return false
}

// trifftFehlerantwort meldet, ob eine der Antworten eine error_response ist.
func trifftFehlerantwort(antworten []model.Response) bool {
	for _, r := range antworten {
		if r.Type == model.ResponseErrorResponse {
			return true
		}
	}
	return false
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
