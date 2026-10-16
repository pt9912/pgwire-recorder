package services_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

// sessionNr ist eine Session der Aufzeichnung mit der Kennung id, der Version
// version und application_name "s<id>", an dem der Handshake sie erkennen lässt.
func sessionNr(id int, version string, ins ...model.Interaction) model.Session {
	return model.Session{ID: id, ServerParameters: map[string]string{"server_version": version, "application_name": fmt.Sprintf("s%d", id)}, Interactions: ins}
}

// gekennzeichnet ist eine Aufzeichnung mit empty_sessions und den Sessions.
func gekennzeichnet(sessions ...model.Session) model.Recording {
	rec := model.NewRecording()
	rec.EmptySessions = true
	rec.Sessions = sessions
	return rec
}

// nachVerbindung lädt rec mit der Zuordnung connection.
func nachVerbindung(t *testing.T, rec model.Recording, opts ...services.ReplayOption) *services.ReplayService {
	t.Helper()
	s, err := services.NewReplayService(context.Background(), ladeRepo{rec: rec}, "rec.yaml", append([]services.ReplayOption{services.ZuordnungNachVerbindung}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// anwendung liest application_name aus dem Handshake.
func anwendung(handshake []model.Response) string {
	for _, r := range handshake {
		if r.Name == "application_name" {
			return r.Value
		}
	}
	return ""
}

// gemischt ist eine Aufzeichnung mit fünf Sessions: 2 und 5 bestehen nur aus
// Lebendprüfungen (5 aus \v bei Version 17), 3 hat keine Interaktion, 4 trägt
// Version 16 und zwei Interaktionen; die letzte Session, 5, hat für das Replay
// keine Interaktion.
func gemischt() model.Recording {
	return gekennzeichnet(
		sessionNr(1, "17.0", interaktion(1, "A", "a")),
		sessionNr(2, "17.0", interaktion(1, "-- ping", "x")),
		sessionNr(3, "17.0"),
		sessionNr(4, "16.4", interaktion(1, "B", "b"), interaktion(2, "C", "c")),
		sessionNr(5, "17.0", interaktion(1, "\v", "x")),
	)
}

// Abdeckung: LH-FA-12/Boundary — mit ZuordnungNachVerbindung erhält die n-te
// angenommene Verbindung die Session mit der Kennung n, auch eine ohne
// Interaktion und eine nur aus Lebendprüfungen: Der Handshake der Verbindung
// trägt die Serverparameter ihrer Session, und die n-te Verbindung fragt die
// Interaktionen der Session n, nicht der n-ten mit Interaktion (LH-FA-12.a
// *Replay*).
func TestReplayConnectionNteVerbindungNteSession(t *testing.T) {
	ctx := context.Background()
	s := nachVerbindung(t, gemischt())
	var ids []model.SessionID
	for n := 1; n <= 5; n++ {
		id, handshake := s.OpenConnection(ctx)
		if want := fmt.Sprintf("s%d", n); anwendung(handshake) != want {
			t.Fatalf("Verbindung %d: Handshake von %q, erwartet %s", n, anwendung(handshake), want)
		}
		ids = append(ids, id)
	}
	if out, err := s.Query(ctx, ids[0], "A"); err != nil || out[0].Tag != "a" {
		t.Fatalf("Verbindung 1: %#v, %v", out, err)
	}
	for _, f := range []struct {
		i   int
		sql string
		tag string
	}{{3, "B", "b"}, {3, "C", "c"}} {
		if out, err := s.Query(ctx, ids[f.i], f.sql); err != nil || out[0].Tag != f.tag {
			t.Fatalf("Verbindung %d, %s: %#v, %v", f.i+1, f.sql, out, err)
		}
	}
	for _, i := range []int{0, 3} {
		if w := gesendetSchliessen(t, s, ids[i]); w != nil {
			t.Fatalf("Verbindung %d verbraucht ihre Session, Warnung: %v", i+1, w)
		}
	}
}

// Abdeckung: LH-FA-12/Boundary, LH-FA-03/Boundary — mit ZuordnungNachVerbindung
// ist eine Session ohne Interaktion, auch eine nur aus Lebendprüfungen, weder
// nicht verbraucht noch nie zugeordnet, mit und ohne FailOnUnconsumed: Ihre
// Verbindung meldet am Ende nichts, und Unassigned nennt sie nicht, auch nicht
// ohne Verbindung (LH-FA-03.b *Verbraucht*, LH-FA-12.a).
func TestReplayConnectionSessionOhneInteraktionNichtZuMelden(t *testing.T) {
	for _, streng := range []bool{false, true} {
		t.Run(fmt.Sprintf("streng %v", streng), func(t *testing.T) {
			ctx := context.Background()
			var opts []services.ReplayOption
			if streng {
				opts = append(opts, services.FailOnUnconsumed)
			}
			rec := gekennzeichnet(
				sessionNr(1, "17.0", interaktion(1, "A", "a")),
				sessionNr(2, "17.0", interaktion(1, "-- ping", "x")),
				sessionNr(3, "17.0"),
			)
			s := nachVerbindung(t, rec, opts...)
			a, _ := s.OpenConnection(ctx)
			leerLebend, _ := s.OpenConnection(ctx)
			leer, _ := s.OpenConnection(ctx)
			if _, err := s.Query(ctx, a, "A"); err != nil {
				t.Fatal(err)
			}
			s.Sent(ctx, a)
			for _, id := range []model.SessionID{a, leerLebend, leer} {
				if w, err := s.CloseConnection(ctx, id); w != nil || err != nil {
					t.Fatalf("Verbindung %d: %v, %v", id, w, err)
				}
			}
			if w, err := s.Unassigned(); w != nil || err != nil {
				t.Fatalf("Unassigned: %v, %v", w, err)
			}
			// Ohne Verbindung nennt Unassigned nur die Session mit Interaktion.
			w, err := nachVerbindung(t, rec, opts...).Unassigned()
			if text := meldungstext(streng, w, err); !strings.Contains(text, "1 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 1") {
				t.Fatalf("Unassigned ohne Verbindung: %v, %v", w, err)
			}
		})
	}
}

// Abdeckung: LH-FA-12/Boundary, LH-FA-03/Boundary — mit ZuordnungNachVerbindung
// ist eine Session mit Interaktionen, deren Verbindung nichts anfragt, mit der
// Annahme zugeordnet und am Ende der Verbindung nicht verbraucht (PGR-W2001,
// mit FailOnUnconsumed PGR-E5002); nie zugeordnet ist nur eine Session mit
// Interaktionen, deren Kennung die Zahl der angenommenen Verbindungen übersteigt
// (LH-FA-12.a *Zuordnung bei `connection` mit der Annahme*).
func TestReplayConnectionZuordnungMitDerAnnahme(t *testing.T) {
	for _, streng := range []bool{false, true} {
		t.Run(fmt.Sprintf("streng %v", streng), func(t *testing.T) {
			ctx := context.Background()
			var opts []services.ReplayOption
			if streng {
				opts = append(opts, services.FailOnUnconsumed)
			}
			rec := gekennzeichnet(
				sessionNr(1, "17.0", interaktion(1, "A", "a")),
				sessionNr(2, "17.0"),
				sessionNr(3, "17.0", interaktion(1, "B", "b"), interaktion(2, "C", "c")),
				sessionNr(4, "17.0", interaktion(1, "D", "d")),
				sessionNr(5, "17.0", interaktion(1, "E", "e")),
			)
			s := nachVerbindung(t, rec, opts...)
			for i := 0; i < 3; i++ {
				s.OpenConnection(ctx)
			}
			// Die dritte Verbindung hat nichts angefragt: ihre Session ist nicht verbraucht.
			w, err := s.CloseConnection(ctx, 3)
			text := meldungstext(streng, w, err)
			if text == "" || !strings.Contains(text, "Session 3: 2 von 2 Interaktionen nicht verbraucht, die erste mit Nummer 1") {
				t.Fatalf("Ende der dritten Verbindung: %v, %v", w, err)
			}
			if streng != (err != nil) {
				t.Fatalf("Fehler %v bei streng %v", err, streng)
			}
			// Sessions 4 und 5 haben mehr Verbindungen als angenommen.
			w, err = s.Unassigned()
			if text := meldungstext(streng, w, err); !strings.Contains(text, "2 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 4") {
				t.Fatalf("Unassigned: %v, %v", w, err)
			}
			s.OpenConnection(ctx)
			w, err = s.Unassigned()
			if text := meldungstext(streng, w, err); !strings.Contains(text, "1 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 5") {
				t.Fatalf("Unassigned nach der vierten Verbindung: %v, %v", w, err)
			}
			s.OpenConnection(ctx)
			if w, err := s.Unassigned(); w != nil || err != nil {
				t.Fatalf("Unassigned nach der fünften Verbindung: %v, %v", w, err)
			}
		})
	}
}

// meldungstext liefert den Text der Meldung, als Fehler bei streng, sonst als
// Warnung; "" ohne Meldung.
func meldungstext(streng bool, w *model.Warning, err error) string {
	switch {
	case streng && err != nil:
		return err.Error()
	case !streng && w != nil:
		return w.Msg
	}
	return ""
}

// Abdeckung: LH-FA-09/Boundary, LH-FA-12/Boundary — mit ZuordnungNachVerbindung
// richtet sich eine Lebendprüfung vom Verbindungsaufbau an nach der Serverversion
// der zugeordneten Session: \v ist ab Version 17 Leerraum. Die erste Anfrage der
// Verbindung zu einer Session der Version 17 an einer mittleren Kennung wird
// beantwortet; bei Version 16 davor und danach ist \v eine abweichende Anfrage
// (LH-FA-12.a *Verhalten je Serverversion*, LH-FA-09.a).
func TestReplayConnectionServerversionJeSession(t *testing.T) {
	ctx := context.Background()
	s := nachVerbindung(t, gekennzeichnet(
		sessionNr(1, "16.4", interaktion(1, "A", "a")),
		sessionNr(2, "17.0", interaktion(1, "B", "b")),
		sessionNr(3, "16.4", interaktion(1, "C", "c")),
	))
	eins, _ := s.OpenConnection(ctx)
	zwei, _ := s.OpenConnection(ctx)
	drei, _ := s.OpenConnection(ctx)
	pruefeLebend(t, s, zwei, "\v", "I")
	for _, id := range []model.SessionID{eins, drei} {
		if out, err := s.Query(ctx, id, "\v"); code(err) != model.CodeReplayMismatch || out != nil {
			t.Fatalf("\\v bei Version 16, Verbindung %d: %#v, %v", id, out, err)
		}
	}
	if out, err := s.Query(ctx, zwei, "B"); err != nil || out[0].Tag != "b" {
		t.Fatalf("Verbindung 2: %#v, %v", out, err)
	}
}

// Abdeckung: LH-FA-12/Boundary — mit ZuordnungNachVerbindung erhält die n-te
// Verbindung mit n über der Zahl der Sessions den Handshake der letzten Session
// der Aufzeichnung, auch wenn diese keine Interaktion hat; eine Anfrage, die keine
// Lebendprüfung ist, ist PGR-E5003, eine Lebendprüfung wird beantwortet, \v nicht,
// weil die Version nicht gesetzt ist (LH-FA-12.a *Mehr Verbindungen als
// Sessions*).
func TestReplayConnectionMehrVerbindungenAlsSessions(t *testing.T) {
	ctx := context.Background()
	s := nachVerbindung(t, gemischt())
	for i := 0; i < 5; i++ {
		s.OpenConnection(ctx)
	}
	id, handshake := s.OpenConnection(ctx)
	if anwendung(handshake) != "s5" {
		t.Fatalf("Handshake von %q, erwartet s5, die letzte Session", anwendung(handshake))
	}
	pruefeLebend(t, s, id, "-- ping", "I")
	for _, sql := range []string{"A", "\v"} {
		out, err := s.Query(ctx, id, sql)
		if code(err) != model.CodeReplaySession || out != nil {
			t.Fatalf("Anfrage %q: %#v, %v, erwartet %s", sql, out, err, model.CodeReplaySession)
		}
	}
	if _, err := s.ClientMessage(ctx, id, model.ClientMessage{Type: model.ClientSync}); code(err) != model.CodeReplaySession {
		t.Fatalf("Extended-Nachricht: %v", err)
	}
	if w, err := s.CloseConnection(ctx, id); w != nil || err != nil {
		t.Fatalf("Ende der Verbindung ohne Session: %v, %v", w, err)
	}
}

// Abdeckung: LH-FA-12/Negative, LH-FA-10/Negative — eine Anfrage, die keine
// Lebendprüfung ist, an eine Session ohne Interaktion ist PGR-E5001, einfach wie
// Extended, auch wenn die Session nur Lebendprüfungen trug; die Diagnose nennt
// die Session und dass die Aufzeichnung keine Interaktion erwartet, und es gibt
// keinen Zugriff außerhalb der Liste. Eine Lebendprüfung wird beantwortet. An eine
// Session, deren Interaktionen alle verbraucht sind, nennt die Diagnose die letzte
// Nummer (LH-FA-12.a *Anfrage an eine Session ohne Interaktion*).
func TestReplayConnectionAnfrageAnSessionOhneInteraktion(t *testing.T) {
	ctx := context.Background()
	s := nachVerbindung(t, gemischt())
	eins, _ := s.OpenConnection(ctx)
	lebend, _ := s.OpenConnection(ctx)
	leer, _ := s.OpenConnection(ctx)
	for _, f := range []struct {
		name string
		id   model.SessionID
		sess int
	}{{"nur Lebendprüfungen", lebend, 2}, {"leer", leer, 3}} {
		t.Run(f.name, func(t *testing.T) {
			if out, err := s.Query(ctx, f.id, "-- ping"); err != nil || !reflect.DeepEqual(out, lebendAntwort("I")) {
				t.Fatalf("Lebendprüfung: %#v, %v", out, err)
			}
			_, err := s.Query(ctx, f.id, "A")
			want := fmt.Sprintf("Session %d: die Aufzeichnung erwartet keine Interaktion; empfangen \"A\"", f.sess)
			if code(err) != model.CodeReplayMismatch || !strings.Contains(err.Error(), want) {
				t.Fatalf("Query: %v, erwartet %s mit %q", err, model.CodeReplayMismatch, want)
			}
			_, err = s.ClientMessage(ctx, f.id, model.ClientMessage{Type: model.ClientParse, Statement: "s1", SQL: "SELECT 1"})
			want = fmt.Sprintf("Session %d: die Aufzeichnung erwartet keine Interaktion; empfangen Client-Nachricht parse (Anweisung \"SELECT 1\")", f.sess)
			if code(err) != model.CodeReplayMismatch || !strings.Contains(err.Error(), want) {
				t.Fatalf("ClientMessage: %v, erwartet %s mit %q", err, model.CodeReplayMismatch, want)
			}
		})
	}
	if _, err := s.Query(ctx, eins, "A"); err != nil {
		t.Fatal(err)
	}
	_, err := s.Query(ctx, eins, "A")
	if want := "Session 1: nach Interaktion 1 erwartet die Aufzeichnung keine weitere; empfangen \"A\""; code(err) != model.CodeReplayMismatch || !strings.Contains(err.Error(), want) {
		t.Fatalf("nach der letzten Interaktion: %v", err)
	}
}

// Abdeckung: LH-FA-12/Negative, LH-FA-03/Negative — mit ZuordnungNachVerbindung
// verlangt die Aufzeichnung die Kennzeichnung empty_sessions, sonst ist
// NewReplayService PGR-E2001 mit Exit-Code 2; die
// Kennzeichnung genügt, eine Session ohne Interaktion verlangt es nicht. Die
// Fehler des Ladens gehen vor: ein Ladefehler bleibt er selbst, und eine
// Aufzeichnung ohne Session mit Interaktion ist PGR-E3004, auch ohne
// Kennzeichnung (LH-FA-12.a *Kennzeichnung bei `connection`*).
func TestReplayConnectionKennzeichnung(t *testing.T) {
	ctx := context.Background()
	mit := func(rec model.Recording, opts ...services.ReplayOption) error {
		_, err := services.NewReplayService(ctx, ladeRepo{rec: rec}, "rec.yaml", opts...)
		return err
	}
	ohne := aufzeichnung([]model.Interaction{interaktion(1, "A", "a")})
	err := mit(ohne, services.ZuordnungNachVerbindung)
	if code(err) != model.CodeUsage || model.Meldungen(err)[0].ExitCode() != 2 {
		t.Fatalf("ohne Kennzeichnung: %v", err)
	}
	if err := mit(ohne); err != nil {
		t.Fatalf("first-request ohne Kennzeichnung: %v", err)
	}
	if err := mit(gekennzeichnet(sessionNr(1, "17.0", interaktion(1, "A", "a"))), services.ZuordnungNachVerbindung); err != nil {
		t.Fatalf("Kennzeichnung ohne leere Session: %v", err)
	}
	// Eine Aufzeichnung ohne Session mit Interaktion: PGR-E3004 vor PGR-E2001.
	for name, rec := range map[string]model.Recording{
		"ohne Sessions":       aufzeichnung(),
		"nur Lebendprüfungen": aufzeichnung([]model.Interaction{interaktion(1, "-- ping", "x")}),
	} {
		if err := mit(rec, services.ZuordnungNachVerbindung); code(err) != model.CodeRecordingNoSession {
			t.Errorf("%s ohne Kennzeichnung: %v, erwartet %s", name, err, model.CodeRecordingNoSession)
		}
	}
	if err := mit(gekennzeichnet(sessionNr(1, "17.0")), services.ZuordnungNachVerbindung); code(err) != model.CodeRecordingNoSession {
		t.Errorf("gekennzeichnet ohne Session mit Interaktion: %v", err)
	}
	kaputt := ladeRepo{err: model.Errorf(model.CodeRecordingBroken, nil, "kaputt")}
	if _, err := services.NewReplayService(ctx, kaputt, "rec.yaml", services.ZuordnungNachVerbindung); code(err) != model.CodeRecordingBroken {
		t.Errorf("Ladefehler: %v", err)
	}
}

// Abdeckung: LH-FA-12/Boundary — mit der Zuordnung first-request ist eine
// gekennzeichnete Aufzeichnung eine wie jede andere: Sessions ohne Interaktion,
// auch die nur aus Lebendprüfungen, werden übersprungen; die erste Verbindung
// erhält die Session 1, die zweite die Session 4, und nie zugeordnet ist keine
// (LH-FA-12.a *Kennzeichnung bei `connection`*).
func TestReplayFirstRequestGekennzeichneteAufzeichnung(t *testing.T) {
	ctx := context.Background()
	s, err := services.NewReplayService(ctx, ladeRepo{rec: gemischt()}, "rec.yaml")
	if err != nil {
		t.Fatal(err)
	}
	eins, handshake := s.OpenConnection(ctx)
	if anwendung(handshake) != "s1" {
		t.Fatalf("Handshake von %q, erwartet s1", anwendung(handshake))
	}
	zwei, _ := s.OpenConnection(ctx)
	if out, err := s.Query(ctx, eins, "A"); err != nil || out[0].Tag != "a" {
		t.Fatalf("erste Verbindung: %#v, %v", out, err)
	}
	for _, sql := range []string{"B", "C"} {
		if out, err := s.Query(ctx, zwei, sql); err != nil || out[0].Tag != strings.ToLower(sql) {
			t.Fatalf("zweite Verbindung, %s: %#v, %v", sql, out, err)
		}
	}
	if w, err := s.Unassigned(); w != nil || err != nil {
		t.Fatalf("Unassigned: %v, %v", w, err)
	}
}

// Abdeckung: LH-FA-13/Boundary, LH-FA-12/Boundary — beim Herunterfahren und
// Prozessende verhält sich connection wie first-request: Shutdown gibt eine
// Verbindung ohne laufende Extended-Interaktion frei; Forced meldet PGR-E4006
// mit der Kennung der Session für eine begonnene, nicht verbrauchte Interaktion
// und nichts für eine zugeordnete Session ohne Interaktion und für eine
// Verbindung ohne Anfrage zu einer Session mit Interaktionen (LH-FA-13.a
// *Zwangsende*, LH-FA-12.a).
func TestReplayConnectionHerunterfahren(t *testing.T) {
	ctx := context.Background()
	s := nachVerbindung(t, gemischt())
	var ids []model.SessionID
	for i := 0; i < 4; i++ {
		id, _ := s.OpenConnection(ctx)
		ids = append(ids, id)
	}
	for i, id := range ids {
		if !s.Shutdown(ctx, id) {
			t.Errorf("Verbindung %d: Shutdown verweigert", i+1)
		}
	}
	if _, err := s.Query(ctx, ids[3], "B"); err != nil {
		t.Fatal(err)
	}
	err := s.Forced(ctx, ids[3])
	if code(err) != model.CodeShutdownTimeout || !strings.Contains(err.Error(), "Session 4, Interaktion 1 nicht verbraucht") {
		t.Fatalf("Forced mit unverbrauchter Interaktion: %v", err)
	}
	s.Sent(ctx, ids[3])
	if err := s.Forced(ctx, ids[3]); err != nil {
		t.Fatalf("Forced nach dem Senden: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := s.Forced(ctx, ids[i]); err != nil {
			t.Errorf("Forced der Verbindung %d ohne Anfrage: %v", i+1, err)
		}
	}
	// Ohne Anfrage zugeordnet, aber nicht verbraucht: die Meldung des Endes der Verbindung.
	if w, err := s.CloseConnection(ctx, ids[0]); err != nil || w == nil || !strings.Contains(w.Msg, "Session 1: 1 von 1 Interaktionen nicht verbraucht") {
		t.Fatalf("Ende der Verbindung 1: %v, %v", w, err)
	}
}

// Abdeckung: LH-FA-12/Boundary — mit ZuordnungNachVerbindung besteht der
// Handshake der n-ten Verbindung aus den Serverparametern der Session n und dem
// ReadyForQuery, auch nachdem eine frühere Verbindung angefragt hat (LH-FA-12.a
// *Replay*).
func TestReplayConnectionHandshakeUnabhaengigVonAnfragen(t *testing.T) {
	ctx := context.Background()
	s := nachVerbindung(t, gemischt())
	eins, _ := s.OpenConnection(ctx)
	if _, err := s.Query(ctx, eins, "A"); err != nil {
		t.Fatal(err)
	}
	_, h := s.OpenConnection(ctx)
	if want := []model.Response{
		{Type: model.ResponseParameterStatus, Name: "application_name", Value: "s2"},
		{Type: model.ResponseParameterStatus, Name: "server_version", Value: "17.0"},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}; !reflect.DeepEqual(h, want) {
		t.Fatalf("Handshake %#v, erwartet %#v", h, want)
	}
}
