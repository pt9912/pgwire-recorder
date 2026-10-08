package pgwire_test

import (
	"context"
	"io"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/pgwire"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// e4006 liefert den Fehler, den der Use Case beim Zwangsende einer laufenden
// Interaktion liefert.
func e4006() *model.Error {
	return model.Errorf(model.CodeShutdownTimeout, nil, "Frist beim Herunterfahren abgelaufen: Interaktion 2 nicht abgeschlossen und verworfen; die Session wird als Session 1 geschrieben")
}

// recordMitZwang startet eine Record-Verbindung über net.Pipe mit ctx und liest
// den Verbindungsaufbau, wenn aufbau gesetzt ist; fertig ist geschlossen, wenn
// Handle zurückkehrt.
func recordMitZwang(ctx context.Context, t *testing.T, rec *fakeRecorder, log io.Writer, aufbau bool) (net.Conn, *pgproto3.Frontend, *pgwire.Server, chan struct{}) {
	t.Helper()
	client, serverSeite := net.Pipe()
	s := pgwire.NewRecordServer(rec, slog.New(slog.NewTextHandler(log, nil)))
	fertig := make(chan struct{})
	go func() {
		pgwire.Handle(ctx, s, serverSeite)
		close(fertig)
	}()
	t.Cleanup(func() { client.Close() })
	_ = client.SetDeadline(time.Now().Add(10 * time.Second))
	fe := pgproto3.NewFrontend(client, client)
	if aufbau {
		startup(t, fe)
	}
	return client, fe, s, fertig
}

// ereignisBinnen wartet höchstens frist, bis der Prüfling ch schließt; sonst
// wird der Test rot und nennt das ausgebliebene Ereignis (SPEC-038 *Warten in
// Tests*).
func ereignisBinnen(t *testing.T, ch chan struct{}, frist time.Duration, ereignis string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(frist):
		t.Fatalf("%s bleibt binnen %v aus", ereignis, frist)
	}
}

// zurueckBinnen wartet höchstens frist, bis fertig geschlossen ist.
func zurueckBinnen(t *testing.T, fertig chan struct{}, frist time.Duration, was string) {
	t.Helper()
	select {
	case <-fertig:
	case <-time.After(frist):
		t.Fatalf("%s: Verbindung endet nicht binnen %v", was, frist)
	}
}

// laufendeInteraktion beginnt eine Extended-Interaktion ohne Sync und fährt
// herunter: Der Use Case gibt das Ende nicht frei, die Sitzung wartet.
func laufendeInteraktion(t *testing.T, rec *fakeRecorder, fe *pgproto3.Frontend, cancel context.CancelFunc) {
	t.Helper()
	sende := nebenher(fe)
	sende(&pgproto3.Parse{Query: "SELECT 2"})
	warteAuf(t, rec, "c:parse")
	cancel()
	warteAuf(t, rec, "shutdown")
}

// Abdeckung: LH-FA-13/Boundary — Record im Adapter: läuft beim Zwangsende eine
// Interaktion, meldet die Sitzung dem Use Case EndForced, stellt dem Client
// den Fehler PGR-E4006 zu (FATAL, SQLSTATE 08006), merkt ihn als ersten Fehler
// und schreibt ihn als Log-Zeile der Stufe error (LH-FA-13.a *Meldung*).
func TestRecordZwangsendeZustellung(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rec := &fakeRecorder{zwangErr: e4006()}
	log := &syncBuffer{}
	_, fe, s, fertig := recordMitZwang(ctx, t, rec, log, true)
	laufendeInteraktion(t, rec, fe, cancel)

	s.Zwangsende()
	e := fehlerantwort(t, fe)
	if e.Severity != "FATAL" || e.Code != "08006" || e.Message != e4006().Error() {
		t.Fatalf("ErrorResponse %+v", e)
	}
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende")
	if got := rec.protokoll(); !strings.HasSuffix(got, "ende:6") || strings.Count(got, "ende:") != 1 {
		t.Fatalf("Protokoll %q, erwartet genau ein Ende EndForced (6)", got)
	}
	if s.FirstErrorCode() != model.CodeShutdownTimeout {
		t.Fatalf("erster Fehler %q", s.FirstErrorCode())
	}
	if !strings.Contains(log.String(), "level=ERROR") || !strings.Contains(log.String(), "code=PGR-E4006") {
		t.Fatalf("Log %s", log.String())
	}
}

// Abdeckung: LH-FA-13/Negative — Record im Adapter: liest der Client beim
// Zwangsende nicht, dauert das Schreiben der Fehlerantwort höchstens 1 s
// (SPEC-051); die Sitzung kehrt binnen 2 s zurück und merkt PGR-E4006.
func TestRecordZwangsendeClientLiestNicht(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rec := &fakeRecorder{zwangErr: e4006()}
	_, fe, s, fertig := recordMitZwang(ctx, t, rec, io.Discard, true)
	laufendeInteraktion(t, rec, fe, cancel)

	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende ohne lesenden Client")
	if s.FirstErrorCode() != model.CodeShutdownTimeout {
		t.Fatalf("erster Fehler %q", s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-13/Boundary — Record im Adapter: ohne laufende Interaktion
// meldet die Sitzung beim Zwangsende nur EndForced; das Schließen der
// Verbindung durch das Zwangsende ist kein weiteres Verbindungsende, der Client
// erhält keine Fehlerantwort, und kein Fehler wird gemerkt (LH-FA-13.a
// *Zwangsende*).
func TestRecordZwangsendeOhneMeldung(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rec := &fakeRecorder{}
	_, fe, s, fertig := recordMitZwang(ctx, t, rec, io.Discard, true)
	laufendeInteraktion(t, rec, fe, cancel)

	s.Zwangsende()
	if msg, err := fe.Receive(); err == nil {
		t.Fatalf("Nachricht beim Zwangsende ohne Meldung: %#v", msg)
	}
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende")
	if got := rec.protokoll(); !strings.HasSuffix(got, "ende:6") || strings.Count(got, "ende:") != 1 {
		t.Fatalf("Protokoll %q, erwartet genau ein Ende EndForced (6)", got)
	}
	if s.FirstErrorCode() != "" {
		t.Fatalf("Fehler gemerkt: %q", s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-13/Boundary — Record im Adapter: hat die Sitzung das Ende
// der Verbindung vor dem Ablauf bemerkt (hier Terminate), endet sie nicht
// zwangsweise, auch wenn ihr Ende beim Zwangsende noch läuft: Der Use Case
// erhält genau ein Ende, EndTerminate (LH-FA-13.a *Was die Frist begrenzt*).
func TestRecordZwangsendeNachBemerktemEnde(t *testing.T) {
	rec := &fakeRecorder{closeHalt: make(chan struct{}), closeLaeuft: make(chan struct{}), zwangErr: e4006()}
	_, fe, s, fertig := recordMitZwang(context.Background(), t, rec, io.Discard, true)
	sende := nebenher(fe)
	sende(&pgproto3.Terminate{})
	ereignisBinnen(t, rec.closeLaeuft, 2*time.Second, "CloseSession nach Terminate")

	s.Zwangsende()
	time.Sleep(100 * time.Millisecond)
	close(rec.closeHalt)
	zurueckBinnen(t, fertig, 2*time.Second, "Ende nach Terminate")
	time.Sleep(100 * time.Millisecond)
	if got := rec.protokoll(); got != "ende:1" {
		t.Fatalf("Protokoll %q, erwartet genau ein Ende EndTerminate (1)", got)
	}
	if s.FirstErrorCode() != "" {
		t.Fatalf("Fehler gemerkt: %q", s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-13/Boundary — Record im Adapter: ein zweiter Aufruf von
// Zwangsende ändert nichts. Ohne laufende Interaktion meldet die Sitzung
// beim Zwangsende nur EndForced; das Schließen der Verbindung durch das
// Zwangsende ist kein weiteres Verbindungsende.
func TestZwangsendeZweimal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rec := &fakeRecorder{}
	_, fe, s, fertig := recordMitZwang(ctx, t, rec, io.Discard, true)
	laufendeInteraktion(t, rec, fe, cancel)
	s.Zwangsende()
	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "zweimal Zwangsende")
	if got := rec.protokoll(); strings.Count(got, "ende:") != 1 {
		t.Fatalf("Protokoll %q", got)
	}
}

// Abdeckung: LH-FA-13/Boundary — Record im Adapter: die Verbindungen enden
// beim Zwangsende nebeneinander. Drei Sessions, deren Clients nicht lesen,
// brauchen je bis zu 1 s für die Fehlerantwort (SPEC-051); zusammen enden sie
// binnen 2 s (SPEC-051, Begründung: die Frist addiert sich nicht).
func TestZwangsendeNebeneinander(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rec := &fakeRecorder{zwangErr: e4006()}
	s := pgwire.NewRecordServer(rec, slog.New(slog.NewTextHandler(io.Discard, nil)))
	var alle []chan struct{}
	for i := 0; i < 3; i++ {
		client, serverSeite := net.Pipe()
		t.Cleanup(func() { client.Close() })
		fertig := make(chan struct{})
		go func() {
			pgwire.Handle(ctx, s, serverSeite)
			close(fertig)
		}()
		_ = client.SetDeadline(time.Now().Add(10 * time.Second))
		fe := pgproto3.NewFrontend(client, client)
		startup(t, fe)
		sende := nebenher(fe)
		sende(&pgproto3.Parse{Query: "SELECT 2"})
		alle = append(alle, fertig)
	}
	deadline := time.Now().Add(2 * time.Second)
	for strings.Count(rec.protokoll(), "c:parse") < 3 {
		if time.Now().After(deadline) {
			t.Fatalf("Protokoll %q", rec.protokoll())
		}
		time.Sleep(5 * time.Millisecond)
	}
	cancel()
	start := time.Now()
	s.Zwangsende()
	for _, fertig := range alle {
		zurueckBinnen(t, fertig, 2*time.Second-time.Since(start), "drei Sessions ohne lesenden Client")
	}
}

// Abdeckung: LH-FA-13/Boundary — Record im Adapter: im Verbindungsaufbau zum
// Upstream wartet das Zwangsende nicht auf den Aufbau; die Verbindung endet
// ohne Fehlerantwort und ohne gemerkten Fehler. Eine Session, die der Aufbau
// danach noch liefert, beendet der Adapter mit EndForced (LH-FA-13.a
// *Zwangsende*).
func TestRecordZwangsendeImAufbau(t *testing.T) {
	rec := &fakeRecorder{openHalt: make(chan struct{}), openLaeuft: make(chan struct{})}
	client, fe, s, fertig := recordMitZwang(context.Background(), t, rec, io.Discard, false)
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
	go func() { _ = fe.Flush() }()
	ereignisBinnen(t, rec.openLaeuft, 2*time.Second, "OpenSession nach der Startnachricht")

	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende im Aufbau")
	if n, err := client.Read(make([]byte, 1)); err == nil {
		t.Fatalf("Antwort erhalten (%d Bytes) statt geschlossener Verbindung", n)
	}
	if s.FirstErrorCode() != "" {
		t.Fatalf("Fehler gemerkt: %q", s.FirstErrorCode())
	}
	close(rec.openHalt)
	if end := rec.lastEnd(t); end != model.EndForced {
		t.Fatalf("späte Session endet mit %v, erwartet EndForced", end)
	}
}

// Abdeckung: LH-FA-13/Boundary — Offen zählt die angenommenen Verbindungen,
// deren Behandlung nicht beendet ist (Attribut sessions, LH-FA-13.a): zwei
// offene Verbindungen sind 2, nach dem Ende der einen 1, nach dem der anderen 0.
func TestOffeneVerbindungen(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := pgwire.NewRecordServer(&fakeRecorder{}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	served := make(chan struct{})
	go func() {
		s.Serve(context.Background(), l)
		close(served)
	}()
	defer func() {
		l.Close()
		zurueckBinnen(t, served, 5*time.Second, "Serve")
	}()
	warteOffen := func(want int) {
		t.Helper()
		deadline := time.Now().Add(2 * time.Second)
		for s.Offen() != want {
			if time.Now().After(deadline) {
				t.Fatalf("Offen %d, erwartet %d", s.Offen(), want)
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
	warteOffen(0)
	a, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	b, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	warteOffen(2)
	a.Close()
	warteOffen(1)
	b.Close()
	warteOffen(0)
}

// Abdeckung: LH-FA-13/Boundary — Record im Adapter: wartet die Client-Richtung
// beim Zwangsende in einer Anfrage auf den Use Case, endet sie mit dessen Ende;
// die Sitzung kehrt erst zurück, wenn der Fehler PGR-E4006 gemerkt und
// zugestellt ist, auch wenn der Use Case ihn erst danach liefert.
func TestRecordZwangsendeWaehrendAnfrage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rec := &fakeRecorder{queryHalt: make(chan struct{}), zwangErr: e4006(), zwangVerzoegerung: 200 * time.Millisecond}
	_, fe, s, fertig := recordMitZwang(ctx, t, rec, io.Discard, true)
	sende := nebenher(fe)
	sende(&pgproto3.Query{String: "SELECT pg_sleep(60)"})
	warteAuf(t, rec, "q:SELECT pg_sleep(60)")
	cancel()
	antwort := make(chan pgproto3.BackendMessage, 1)
	go func() {
		msg, _ := fe.Receive()
		antwort <- msg
	}()

	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende während einer Anfrage")
	if s.FirstErrorCode() != model.CodeShutdownTimeout {
		t.Fatalf("Sitzung kehrt zurück, bevor PGR-E4006 gemerkt ist: %q", s.FirstErrorCode())
	}
	select {
	case msg := <-antwort:
		if e, ok := msg.(*pgproto3.ErrorResponse); !ok || !strings.Contains(e.Message, model.CodeShutdownTimeout) {
			t.Fatalf("Client erhält %#v statt PGR-E4006", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Client erhält binnen 2 s keine Nachricht")
	}
}

// Abdeckung: LH-FA-13/Boundary — Record im Adapter: scheitert das Senden des
// Verbindungsaufbaus an den Client, weil das Zwangsende die Verbindung
// geschlossen hat, ist das kein weiteres Verbindungsende: Mit Session meldet die
// Sitzung EndForced statt EndWriteFailed, ohne Session (Fehlerantwort des
// Servers im Aufbau) merkt sie kein PGR-E4003 (LH-FA-13.a *Zwangsende*).
func TestRecordZwangsendeBeimSendenDesAufbaus(t *testing.T) {
	for _, f := range []struct {
		name         string
		aufbauFehler bool
		wantEnde     string
	}{
		{"mit Session", false, "ende:6"},
		{"Fehlerantwort des Servers", true, ""},
	} {
		t.Run(f.name, func(t *testing.T) {
			rec := &fakeRecorder{aufbauFehler: f.aufbauFehler}
			_, fe, s, fertig := recordMitZwang(context.Background(), t, rec, io.Discard, false)
			fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
			go func() { _ = fe.Flush() }()
			// Der Client liest den Aufbau nicht: das Senden blockiert.
			time.Sleep(100 * time.Millisecond)

			s.Zwangsende()
			zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende beim Senden des Aufbaus")
			if got := rec.protokoll(); got != f.wantEnde {
				t.Fatalf("Protokoll %q, erwartet %q", got, f.wantEnde)
			}
			if s.FirstErrorCode() != "" {
				t.Fatalf("Fehler gemerkt: %q", s.FirstErrorCode())
			}
		})
	}
}

// Abdeckung: LH-FA-13/Boundary — Record im Adapter: das Zwangsende bricht den
// Kontext des Verbindungsaufbaus beim Use Case ab, und dessen Fehler wird nicht
// gemerkt (LH-FA-13.a *Zwangsende*).
func TestRecordZwangsendeBrichtAufbauAb(t *testing.T) {
	rec := &fakeRecorder{openMitKontext: true, openLaeuft: make(chan struct{}), openAbgebrochen: make(chan struct{})}
	_, fe, s, fertig := recordMitZwang(context.Background(), t, rec, io.Discard, false)
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
	go func() { _ = fe.Flush() }()
	ereignisBinnen(t, rec.openLaeuft, 2*time.Second, "OpenSession nach der Startnachricht")

	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende im Aufbau")
	select {
	case <-rec.openAbgebrochen:
	case <-time.After(2 * time.Second):
		t.Fatal("Kontext des Aufbaus nicht abgebrochen")
	}
	time.Sleep(50 * time.Millisecond)
	if s.FirstErrorCode() != "" {
		t.Fatalf("Fehler gemerkt: %q", s.FirstErrorCode())
	}
}

// e4006Replay liefert den Fehler, den der Replay-Use-Case beim Zwangsende einer
// begonnenen, nicht verbrauchten Interaktion liefert.
func e4006Replay() *model.Error {
	return model.Errorf(model.CodeShutdownTimeout, nil, "Frist beim Herunterfahren abgelaufen: Session 1, Interaktion 2 nicht verbraucht")
}

// replayMitZwang startet eine Replay-Verbindung über net.Pipe mit ctx und liest
// den Verbindungsaufbau, wenn aufbau gesetzt ist; fertig ist geschlossen, wenn
// Handle zurückkehrt.
func replayMitZwang(ctx context.Context, t *testing.T, rep *fakeReplayer, aufbau bool) (*pgproto3.Frontend, *pgwire.Server, *syncBuffer, chan struct{}) {
	t.Helper()
	client, serverSeite := net.Pipe()
	log := &syncBuffer{}
	s := pgwire.NewReplayServer(rep, slog.New(slog.NewTextHandler(log, nil)))
	fertig := make(chan struct{})
	go func() {
		pgwire.Handle(ctx, s, serverSeite)
		close(fertig)
	}()
	t.Cleanup(func() { client.Close() })
	_ = client.SetDeadline(time.Now().Add(10 * time.Second))
	fe := pgproto3.NewFrontend(client, client)
	if aufbau {
		startup(t, fe)
	}
	return fe, s, log, fertig
}

// replayLaufend sendet ein Parse ohne Sync und wartet, bis der Use Case es hat:
// Danach läuft eine Extended-Interaktion.
func replayLaufend(t *testing.T, rep *fakeReplayer, fe *pgproto3.Frontend) {
	t.Helper()
	sende := nebenher(fe)
	sende(&pgproto3.Parse{Query: "SELECT 2"})
	deadline := time.Now().Add(2 * time.Second)
	for len(rep.nachrichten()) == 0 {
		if time.Now().After(deadline) {
			t.Fatal("Parse erreicht den Use Case nicht")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// liestWieder wartet, bis die Sitzung beim Herunterfahren den Use Case gefragt
// hat, und gibt ihr Zeit, wieder im Lesen zu blockieren: Danach beendet nur
// noch die Lesefrist des Zwangsendes das Lesen.
func liestWieder(t *testing.T, rep *fakeReplayer) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		rep.mu.Lock()
		n := rep.shutdowns
		rep.mu.Unlock()
		if n > 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("Sitzung fragt den Use Case beim Herunterfahren nicht")
		}
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
}

// lieseNebenher liest die nächste Nachricht des Clients nebenher.
func lieseNebenher(fe *pgproto3.Frontend) chan pgproto3.BackendMessage {
	antwort := make(chan pgproto3.BackendMessage, 1)
	go func() {
		msg, _ := fe.Receive()
		antwort <- msg
	}()
	return antwort
}

// Abdeckung: LH-FA-13/Boundary — Replay im Adapter (V-95):
// läuft beim Zwangsende eine Extended-Interaktion, beendet die Frist die
// Replay-Sitzung, und die Sitzung meldet dem Use Case das Zwangsende (Forced);
// den Fehler PGR-E4006, den er liefert, stellt sie dem Client zu (FATAL,
// SQLSTATE 08006), merkt ihn als ersten Fehler und schreibt ihn als Log-Zeile
// der Stufe error, danach beendet sie die Verbindung (CloseConnection).
func TestReplayZwangsende(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rep := &fakeReplayer{zwangFehler: e4006Replay()}
	fe, s, log, fertig := replayMitZwang(ctx, t, rep, true)
	replayLaufend(t, rep, fe)
	cancel()
	liestWieder(t, rep)
	antwort := lieseNebenher(fe)

	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende im Replay")
	select {
	case msg := <-antwort:
		if e, ok := msg.(*pgproto3.ErrorResponse); !ok || e.Severity != "FATAL" || e.Code != "08006" || e.Message != e4006Replay().Error() {
			t.Fatalf("Client erhält %#v statt PGR-E4006", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Client erhält binnen 2 s keine Nachricht")
	}
	if got := rep.protokoll(); got != "forced close" {
		t.Fatalf("Protokoll %q, erwartet „forced close“", got)
	}
	if s.FirstErrorCode() != model.CodeShutdownTimeout || !strings.Contains(log.String(), "level=ERROR") || !strings.Contains(log.String(), "code=PGR-E4006") {
		t.Fatalf("erster Fehler %q, Log %s", s.FirstErrorCode(), log.String())
	}
}

// Abdeckung: LH-FA-13/Negative — Replay im Adapter: mit --fail-on-unconsumed
// merkt die Sitzung beim Zwangsende PGR-E4006 vor PGR-E5002; beide stehen in
// dieser Reihenfolge im Log, und der Client erhält nur PGR-E4006 (LH-FA-03.b
// *Fehlerebene*).
func TestReplayZwangsendeVorNichtVerbraucht(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	unverbraucht := model.Errorf(model.CodeReplayUnconsumed, nil, "Session 1: 1 von 2 Interaktionen nicht verbraucht, die erste mit Nummer 2")
	rep := &fakeReplayer{zwangFehler: e4006Replay(), fehler: unverbraucht}
	fe, s, log, fertig := replayMitZwang(ctx, t, rep, true)
	replayLaufend(t, rep, fe)
	cancel()
	liestWieder(t, rep)
	antwort := lieseNebenher(fe)

	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende im Replay")
	select {
	case msg := <-antwort:
		if e, ok := msg.(*pgproto3.ErrorResponse); !ok || !strings.Contains(e.Message, model.CodeShutdownTimeout) {
			t.Fatalf("Client erhält %#v", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Client erhält binnen 2 s keine Nachricht")
	}
	if msg, err := fe.Receive(); err == nil {
		t.Fatalf("zweite Nachricht an den Client: %#v", msg)
	}
	text := log.String()
	if s.FirstErrorCode() != model.CodeShutdownTimeout || strings.Index(text, "PGR-E4006") > strings.Index(text, "PGR-E5002") || !strings.Contains(text, "PGR-E5002") {
		t.Fatalf("erster Fehler %q, Log %s", s.FirstErrorCode(), text)
	}
}

// Abdeckung: LH-FA-13/Boundary — Replay im Adapter: ohne begonnene Interaktion
// liefert der Use Case beim Zwangsende nichts; die Sitzung endet ohne
// Fehlerantwort, und das Schließen durch das Zwangsende ist kein weiteres
// Verbindungsende (kein PGR-E4003).
func TestReplayZwangsendeOhneMeldung(t *testing.T) {
	rep := &fakeReplayer{}
	fe, s, _, fertig := replayMitZwang(context.Background(), t, rep, true)

	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende im Replay")
	if msg, err := fe.Receive(); err == nil {
		t.Fatalf("Nachricht beim Zwangsende ohne Meldung: %#v", msg)
	}
	if got := rep.protokoll(); got != "forced close" {
		t.Fatalf("Protokoll %q, erwartet „forced close“", got)
	}
	if s.FirstErrorCode() != "" {
		t.Fatalf("Fehler gemerkt: %q", s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-13/Negative — Replay im Adapter: liest der Client beim
// Zwangsende die Antworten nicht, bricht die Schreibfrist von 1 s (SPEC-051)
// das Senden ab; die Interaktion ist dann nicht verbraucht, die Sitzung meldet
// das Zwangsende, merkt PGR-E4006 statt PGR-E4003 und kehrt binnen 2 s zurück.
func TestReplayZwangsendeClientLiestNicht(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rep := &fakeReplayer{zwangFehler: e4006Replay()}
	fe, s, log, fertig := replayMitZwang(ctx, t, rep, true)
	sende := nebenher(fe)
	sende(&pgproto3.Sync{})
	deadline := time.Now().Add(2 * time.Second)
	for len(rep.nachrichten()) == 0 {
		if time.Now().After(deadline) {
			t.Fatal("Sync erreicht den Use Case nicht")
		}
		time.Sleep(5 * time.Millisecond)
	}
	cancel()

	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende ohne lesenden Client")
	if got := rep.protokoll(); got != "forced close" || rep.sentAufrufe() != 0 {
		t.Fatalf("Protokoll %q, Sent %d", got, rep.sentAufrufe())
	}
	if s.FirstErrorCode() != model.CodeShutdownTimeout || strings.Contains(log.String(), "PGR-E4003") {
		t.Fatalf("erster Fehler %q, Log %s", s.FirstErrorCode(), log.String())
	}
}

// Abdeckung: LH-FA-13/Boundary — Replay im Adapter: trifft das Zwangsende die
// Sitzung, während sie beim Herunterfahren den Use Case fragt, hebt das
// Zurücksetzen der Lesefrist danach das Zwangsende nicht auf; die Sitzung
// endet binnen 2 s mit Forced.
func TestReplayZwangsendeWaehrendShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rep := &fakeReplayer{shutdownHalt: make(chan struct{}), shutdownLaeuft: make(chan struct{})}
	fe, s, _, fertig := replayMitZwang(ctx, t, rep, true)
	replayLaufend(t, rep, fe)
	cancel()
	ereignisBinnen(t, rep.shutdownLaeuft, 2*time.Second, "Shutdown nach dem Ende von ctx")

	s.Zwangsende()
	time.Sleep(100 * time.Millisecond)
	close(rep.shutdownHalt)
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende während Shutdown")
	if got := rep.protokoll(); got != "forced close" {
		t.Fatalf("Protokoll %q, erwartet „forced close“", got)
	}
}

// Abdeckung: LH-FA-13/Boundary — Replay im Adapter: scheitert das Senden des
// Verbindungsaufbaus, weil das Zwangsende die Verbindung geschlossen hat, merkt
// die Sitzung kein PGR-E4003 und beendet die Verbindung (LH-FA-13.a
// *Zwangsende*).
func TestReplayZwangsendeBeimSendenDesAufbaus(t *testing.T) {
	rep := &fakeReplayer{}
	fe, s, _, fertig := replayMitZwang(context.Background(), t, rep, false)
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber, Parameters: map[string]string{"user": "app"}})
	go func() { _ = fe.Flush() }()
	// Der Client liest den Aufbau nicht: das Senden blockiert.
	time.Sleep(100 * time.Millisecond)

	s.Zwangsende()
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende beim Senden des Aufbaus")
	if got := rep.protokoll(); got != "close" {
		t.Fatalf("Protokoll %q, erwartet „close“", got)
	}
	if s.FirstErrorCode() != "" {
		t.Fatalf("Fehler gemerkt: %q", s.FirstErrorCode())
	}
}

// Abdeckung: LH-FA-13/Negative — Record im Adapter: liefert der Use Case beim
// Zwangsende einen anderen Fehler als PGR-E4006 (hier das Schreiben der
// Aufzeichnung, PGR-E3001), merkt die Sitzung ihn und schreibt ihn als
// Log-Zeile der Stufe error, stellt ihn dem Client aber nicht zu.
func TestRecordZwangsendeAndererFehler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	schreibfehler := model.Errorf(model.CodeRecordingIO, nil, "rec.yaml nicht schreibbar")
	rec := &fakeRecorder{zwangErr: schreibfehler}
	log := &syncBuffer{}
	_, fe, s, fertig := recordMitZwang(ctx, t, rec, log, true)
	laufendeInteraktion(t, rec, fe, cancel)

	s.Zwangsende()
	if msg, err := fe.Receive(); err == nil {
		t.Fatalf("Fehler dem Client zugestellt: %#v", msg)
	}
	zurueckBinnen(t, fertig, 2*time.Second, "Zwangsende")
	if s.FirstErrorCode() != model.CodeRecordingIO || !strings.Contains(log.String(), "level=ERROR") || !strings.Contains(log.String(), "code=PGR-E3001") {
		t.Fatalf("erster Fehler %q, Log %s", s.FirstErrorCode(), log.String())
	}
}
