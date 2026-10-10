package bootstrap_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/bootstrap"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// nachricht ist eine PGWire-Nachricht aus Typ und Rumpf.
func nachricht(typ byte, rumpf ...byte) []byte {
	out := []byte{typ, 0, 0, 0, 0}
	binary.BigEndian.PutUint32(out[1:], uint32(len(rumpf)+4))
	return append(out, rumpf...)
}

// bereitOhneAnmeldung ist AuthenticationOk und ReadyForQuery.
func bereitOhneAnmeldung() []byte {
	return append(nachricht('R', 0, 0, 0, 0), nachricht('Z', 'I')...)
}

// fehlerantwortRoh ist eine ErrorResponse mit Schweregrad, SQLSTATE und Meldung.
func fehlerantwortRoh(schwere, code, text string) []byte {
	var r []byte
	for _, f := range [][2]string{{"S", schwere}, {"V", schwere}, {"C", code}, {"M", text}} {
		r = append(append(append(r, f[0][0]), f[1]...), 0)
	}
	return nachricht('E', append(r, 0)...)
}

// playServer nimmt Verbindungen an; je Verbindung liest er das Startup, sendet
// bereitOhneAnmeldung und beantwortet jede Query mit antwort, ohne antwort gar
// nicht. Jede Client-Nachricht außer Query und alles nach einer
// unbeantworteten Query meldet er mit ihrem Typ auf typen.
func playServer(t *testing.T, antwort []byte) (string, <-chan byte) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	typen := make(chan byte, 64)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
				var laenge [4]byte
				if _, err := io.ReadFull(conn, laenge[:]); err != nil {
					return
				}
				if _, err := io.ReadFull(conn, make([]byte, binary.BigEndian.Uint32(laenge[:])-4)); err != nil {
					return
				}
				_, _ = conn.Write(bereitOhneAnmeldung())
				for {
					var kopf [5]byte
					if _, err := io.ReadFull(conn, kopf[:]); err != nil {
						return
					}
					if _, err := io.ReadFull(conn, make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)); err != nil {
						return
					}
					typen <- kopf[0]
					if kopf[0] == 'Q' && antwort != nil {
						_, _ = conn.Write(antwort)
					}
				}
			}()
		}
	}()
	return l.Addr().String(), typen
}

// einfacheAufzeichnung schreibt eine Aufzeichnung mit einer Session aus den
// Anfragen und liefert ihren Pfad.
func einfacheAufzeichnung(t *testing.T, sessions ...[]model.Interaction) string {
	t.Helper()
	rec := model.NewRecording()
	for i, in := range sessions {
		rec.Sessions = append(rec.Sessions, model.Session{ID: i + 1, Startup: map[string]string{"user": "u"}, Interactions: in})
	}
	path := filepath.Join(t.TempDir(), "rec.yaml")
	if err := (recording.YAML{}).Write(context.Background(), path, rec); err != nil {
		t.Fatal(err)
	}
	return path
}

func anfrage(seq int, sql string) model.Interaction {
	return model.Interaction{Sequence: seq, Request: model.Request{Type: model.RequestQuery, SQL: sql}, Responses: []model.Response{
		{Type: model.ResponseCommandComplete, Tag: "SELECT 0"},
		{Type: model.ResponseReadyForQuery, TxStatus: "I"},
	}}
}

// leerePlay setzt die Umgebungsvariablen von play leer.
func leerePlay(t *testing.T) {
	for _, n := range []string{"UPSTREAM", "INPUT", "USER", "DATABASE", "LOG_LEVEL", "CONFIG", "PASSWORD"} {
		t.Setenv("PGWIRE_RECORDER_"+n, "")
	}
	t.Chdir(t.TempDir())
}

// typenBis liest die Typen der Client-Nachrichten, bis ende kommt, höchstens
// 30 s lang.
func typenBis(t *testing.T, typen <-chan byte, ende byte) string {
	t.Helper()
	var out []byte
	for {
		select {
		case b := <-typen:
			out = append(out, b)
			if b == ende {
				return string(out)
			}
		case <-time.After(30 * time.Second):
			t.Fatalf("binnen 30 s keine Client-Nachricht %q, bisher %q", ende, out)
		}
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-14/Boundary — play spielt die Anfragen ein
// und endet mit Exit-Code 0; die Zeilen der Stufe info nennen Start mit
// Adresse und Aufzeichnung, ohne Benutzer und Datenbank, und Ende; jede
// Session endet mit Terminate (LH-FA-20.a *Meldungen*, *Ende einer Session*).
func TestRunPlay(t *testing.T) {
	leerePlay(t)
	addr, typen := playServer(t, append(nachricht('C', []byte("SELECT 0\x00")...), nachricht('Z', 'I')...))
	input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A"), anfrage(2, "B")}, []model.Interaction{anfrage(1, "C")})
	var stdout, stderr bytes.Buffer
	code := bootstrap.Run(context.Background(), nil, []string{"play", "--upstream", addr, "--input", input, "--user", "GEHEIMU", "--database", "GEHEIMD"}, "dev", &stdout, &stderr)
	if code != 0 || stdout.Len() > 0 {
		t.Fatalf("Exit-Code %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
	if got := typenBis(t, typen, 'X'); got != "QQX" {
		t.Fatalf("erste Session: %q", got)
	}
	if got := typenBis(t, typen, 'X'); got != "QX" {
		t.Fatalf("zweite Session: %q", got)
	}
	if got := stufenImLog(t, stderr.String()); strings.Join(got, ",") != "INFO,INFO" {
		t.Fatalf("Zeilen %v\n%s", got, stderr.String())
	}
	if !strings.Contains(stderr.String(), "upstream="+addr) || !strings.Contains(stderr.String(), "input="+input) || strings.Contains(stderr.String(), "GEHEIM") {
		t.Fatalf("Zeile beim Start: %q", stderr.String())
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-14/Boundary — eine Fehlerantwort des
// Servers bricht ab: Exit-Code 4, eine Zeile der Stufe error mit code
// PGR-E4004 und dem Fehlertext mit Session, Interaktion, SQLSTATE und
// Meldung; keine weitere Anfrage, Terminate (LH-FA-20.a).
func TestRunPlayFehlerantwort(t *testing.T) {
	leerePlay(t)
	addr, typen := playServer(t, append(fehlerantwortRoh("ERROR", "42P01", "fehlt"), nachricht('Z', 'I')...))
	input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A"), anfrage(2, "B")})
	var stdout, stderr bytes.Buffer
	code := bootstrap.Run(context.Background(), nil, []string{"play", "--upstream", addr, "--input", input}, "dev", &stdout, &stderr)
	if code != 4 || !strings.Contains(stderr.String(), `error="Netzwerk [PGR-E4004]: Session 1, Interaktion 1: Fehlerantwort des Servers 42P01 „fehlt“"`) {
		t.Fatalf("Exit-Code %d, stderr %q", code, stderr.String())
	}
	if got := stufenImLog(t, stderr.String()); strings.Join(got, ",") != "INFO,ERROR PGR-E4004,INFO" {
		t.Fatalf("Zeilen %v", got)
	}
	if got := typenBis(t, typen, 'X'); got != "QX" {
		t.Fatalf("Client-Nachrichten %q", got)
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-14/Negative — ein Startfehler ist die
// Zeile beim Prozessende ohne Log-Zeile davor, auch nach einem Abbruchsignal
// während des Starts: eine nicht ladbare Aufzeichnung Exit-Code 3; play
// verbindet sich dabei nicht (LH-FA-20.a *Start*).
func TestRunPlayStartfehler(t *testing.T) {
	leerePlay(t)
	addr, typen := playServer(t, nil)
	beendet, aus := context.WithCancel(context.Background())
	aus()
	for _, ctx := range []context.Context{context.Background(), beendet} {
		var stdout, stderr bytes.Buffer
		if code := bootstrap.Run(ctx, nil, []string{"play", "--upstream", addr, "--input", filepath.Join(t.TempDir(), "fehlt.yaml")}, "dev", &stdout, &stderr); code != 3 || strings.Contains(stderr.String(), "level=") {
			t.Fatalf("nicht ladbar: Exit-Code %d, stderr %q", code, stderr.String())
		}
	}
	select {
	case b := <-typen:
		t.Fatalf("Startfehler verbindet sich: %q", b)
	default:
	}
}

// Abdeckung: LH-FA-20/Boundary — eine Aufzeichnung ohne Session mit
// Interaktion endet ohne Verbindung mit Exit-Code 0 (LH-FA-20.a *Start*).
func TestRunPlayOhneInteraktion(t *testing.T) {
	leerePlay(t)
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	var stdout, stderr bytes.Buffer
	if code := bootstrap.Run(context.Background(), nil, []string{"play", "--upstream", addr, "--input", einfacheAufzeichnung(t)}, "dev", &stdout, &stderr); code != 0 {
		t.Fatalf("Exit-Code %d, stderr %q", code, stderr.String())
	}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-14/Boundary — ein Abbruchsignal während
// des Starts bricht ihn nicht ab; danach beginnt keine Session, die Zeile der
// Stufe info zum ersten Signal steht im Log, Exit-Code 0 (LH-FA-20.a
// *Start*, LH-FA-14.a).
func TestRunPlaySignalImStart(t *testing.T) {
	leerePlay(t)
	addr, typen := playServer(t, nil)
	ctx, aus := context.WithCancel(context.Background())
	aus()
	var stdout, stderr bytes.Buffer
	code := bootstrap.Run(ctx, nil, []string{"play", "--upstream", addr, "--input", einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A")})}, "dev", &stdout, &stderr)
	if code != 0 || !strings.Contains(stderr.String(), "Abbruchsignal") {
		t.Fatalf("Exit-Code %d, stderr %q", code, stderr.String())
	}
	if got := stufenImLog(t, stderr.String()); strings.Join(got, ",") != "INFO,INFO,INFO" {
		t.Fatalf("Zeilen %v", got)
	}
	select {
	case b := <-typen:
		t.Fatalf("Session nach dem Signal im Start: %q", b)
	default:
	}
}

// Abdeckung: LH-FA-20/Boundary — das zweite Signal beendet play sofort, auch
// während es auf eine Antwort wartet: Terminate geht an den Server, die
// unterbrochene Interaktion ist kein Fehler, Exit-Code 0 (LH-FA-20.a
// *Abbruchsignal*).
func TestRunPlayZweitesSignal(t *testing.T) {
	leerePlay(t)
	addr, typen := playServer(t, nil)
	ctx, erstes := context.WithCancel(context.Background())
	ablauf := make(chan struct{})
	fertig := make(chan int, 1)
	var stdout, stderr bytes.Buffer
	input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A"), anfrage(2, "B")})
	go func() {
		fertig <- bootstrap.Run(ctx, ablauf, []string{"play", "--upstream", addr, "--input", input}, "dev", &stdout, &stderr)
	}()
	if got := typenBis(t, typen, 'Q'); got != "Q" {
		t.Fatalf("Client-Nachrichten %q", got)
	}
	erstes()
	close(ablauf)
	select {
	case code := <-fertig:
		if code != 0 {
			t.Fatalf("Exit-Code %d, stderr %q", code, stderr.String())
		}
	case <-time.After(30 * time.Second):
		t.Fatal("play endet binnen 30 s nach dem zweiten Signal nicht")
	}
	if got := typenBis(t, typen, 'X'); got != "X" {
		t.Fatalf("nach dem zweiten Signal %q, erwartet Terminate", got)
	}
	if strings.Contains(stderr.String(), "level=ERROR") {
		t.Fatalf("Fehler nach dem zweiten Signal: %q", stderr.String())
	}
}

// startupServer nimmt eine Verbindung an, liefert ihre Startup-Parameter auf
// dem Kanal und lehnt den Aufbau mit einer Fehlerantwort ab.
func startupServer(t *testing.T) (string, <-chan map[string]string) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	parameter := make(chan map[string]string, 1)
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
		var laenge [4]byte
		if _, err := io.ReadFull(conn, laenge[:]); err != nil {
			return
		}
		rumpf := make([]byte, binary.BigEndian.Uint32(laenge[:])-4)
		if _, err := io.ReadFull(conn, rumpf); err != nil {
			return
		}
		// Nach der Protokollversion folgen Name und Wert, je mit NUL beendet,
		// und ein abschließendes NUL.
		teile := strings.Split(string(rumpf[4:]), "\x00")
		p := map[string]string{}
		for i := 0; i+1 < len(teile) && teile[i] != ""; i += 2 {
			p[teile[i]] = teile[i+1]
		}
		parameter <- p
		_, _ = conn.Write(fehlerantwortRoh("FATAL", "3D000", "nein"))
	}()
	return l.Addr().String(), parameter
}

// Abdeckung: LH-FA-20/Boundary — --user und --database erreichen den
// Play-Service und gehen als user und database im Startup an den Server,
// --user statt des Benutzers der Aufzeichnung (LH-FA-20.a *Startup-Daten*).
func TestRunPlayOptionen(t *testing.T) {
	leerePlay(t)
	addr, parameter := startupServer(t)
	input := einfacheAufzeichnung(t, []model.Interaction{anfrage(1, "A")})
	var stdout, stderr bytes.Buffer
	code := bootstrap.Run(context.Background(), nil, []string{"play", "--upstream", addr, "--input", input, "--user", "optu", "--database", "optd"}, "dev", &stdout, &stderr)
	if code != 4 {
		t.Fatalf("Exit-Code %d, stderr %q", code, stderr.String())
	}
	select {
	case p := <-parameter:
		if p["user"] != "optu" || p["database"] != "optd" {
			t.Fatalf("Startup %v, erwartet user optu und database optd", p)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("der Fake-Server empfängt binnen 30 s kein Startup")
	}
}
