package postgres_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/pbkdf2"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// anmNonce ist die Nonce von play in den Tests mit festem Zufall: das Base64
// der Bytes 1 bis 18.
const anmNonce = "AQIDBAUGBwgJCgsMDQ4PEBES"

// anmZufall liefert die Bytes 1 bis 18.
func anmZufall() io.Reader {
	b := make([]byte, 18)
	for i := range b {
		b[i] = byte(i + 1)
	}
	return bytes.NewReader(b)
}

// anmeldenMit führt den Aufbau gegen den Fake-Server unter addr mit passwort
// und der festen Nonce aus und schließt die Verbindung danach wie Verbinde,
// ohne Terminate.
func anmeldenMit(t *testing.T, addr string, startup map[string]string, passwort string) error {
	t.Helper()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(20 * time.Second))
	return postgres.AufbauMit(conn, startup, passwort, anmZufall())
}

// anmLauf führt einen Aufbau mit passwort gegen einen Fake-Server aus, der
// aufbau sendet und je Eintrag von jeAnfrage eine Client-Nachricht liest und
// den Eintrag sendet, und liefert den lauf des Servers und den Fehler des
// Aufbaus.
func anmLauf(t *testing.T, passwort string, aufbau []byte, jeAnfrage ...[]byte) (lauf, error) {
	t.Helper()
	addr, ergebnis := einspielServer(t, aufbau, jeAnfrage...)
	err := anmeldenMit(t, addr, map[string]string{"user": "u"}, passwort)
	return geschlossen(t, ergebnis), err
}

// anmKette hängt Nachrichten des Servers aneinander, ohne eine zu verändern.
func anmKette(teile ...[]byte) []byte {
	var out []byte
	for _, t := range teile {
		out = append(out, t...)
	}
	return out
}

// anmSasl ist die Anforderung von SASL (Code 10) mit der Liste der Verfahren.
func anmSasl(namen ...string) []byte {
	var rest []byte
	for _, n := range namen {
		rest = append(append(rest, n...), 0)
	}
	return anmeldung(10, append(rest, 0)...)
}

// anmMD5 ist die Anforderung von MD5 (Code 5) mit dem Salz.
func anmMD5(salz ...byte) []byte { return anmeldung(5, salz...) }

// anmKlartext ist die Anforderung des Klartext-Passworts (Code 3).
func anmKlartext() []byte { return anmeldung(3) }

// anmBereit ist AuthenticationOk und ReadyForQuery.
func anmBereit(t *testing.T) []byte {
	return kodiert(t, &pgproto3.AuthenticationOk{}, &pgproto3.ReadyForQuery{TxStatus: 'I'})
}

// anmFehler ist eine ErrorResponse mit SQLSTATE und Meldung.
func anmFehler(t *testing.T, sqlstate, text string) []byte {
	return kodiert(t, &pgproto3.ErrorResponse{Severity: "FATAL", SeverityUnlocalized: "FATAL", Code: sqlstate, Message: text})
}

// anmPasswort ist die Client-Nachricht mit einem Passwort.
func anmPasswort(passwort string) clientNachricht {
	return clientNachricht{'p', []byte(passwort + "\x00")}
}

// anmErste ist die SASLInitialResponse von play mit der festen Nonce.
func anmErste() clientNachricht {
	data := "n,,n=,r=" + anmNonce
	rumpf := append([]byte("SCRAM-SHA-256\x00"), binary.BigEndian.AppendUint32(nil, uint32(len(data)))...)
	return clientNachricht{'p', append(rumpf, data...)}
}

// anmNichts prüft, dass der Client nach dem Startup nichts gesendet hat: weder
// eine Antwort noch Terminate.
func anmNichts(t *testing.T, l lauf) {
	t.Helper()
	if len(l.gelesen) != 0 || len(l.danach) != 0 {
		t.Fatalf("der Client sendet nach dem Startup %v und %q", l.gelesen, l.danach)
	}
}

// anmGesendet prüft die Nachrichten, die der Client nach dem Startup gesendet
// hat, und dass danach nichts mehr kam, auch kein Terminate.
func anmGesendet(t *testing.T, l lauf, want ...clientNachricht) {
	t.Helper()
	if fmt.Sprint(l.gelesen) != fmt.Sprint(want) || len(l.danach) != 0 {
		t.Fatalf("der Client sendet %q und danach %q, erwartet %q", l.gelesen, l.danach, want)
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary — verlangt der Server das
// Klartext-Passwort (Code 3), antwortet play mit genau dem Passwort,
// unverändert (Leerraum, Prozentzeichen und Nicht-ASCII-Zeichen bleiben), und
// sendet sonst nichts; hinter AuthenticationOk und ReadyForQuery steht die
// Session (LH-FA-20.a *Anmeldung*, *Verfahren*, *Passwort*).
func TestAnmeldungKlartext(t *testing.T) {
	for _, pw := range []string{"GEHEIMpw", "  ä%41$ ss "} {
		addr, ergebnis := einspielServer(t, anmKlartext(), anmBereit(t))
		if err := anmeldenMit(t, addr, map[string]string{"user": "u"}, pw); err != nil {
			t.Fatalf("Klartext mit Passwort %q: %v", pw, err)
		}
		anmGesendet(t, geschlossen(t, ergebnis), anmPasswort(pw))
	}
}

// Abdeckung: LH-FA-20/Happy, LH-FA-20/Boundary — verlangt der Server MD5 (Code
// 5 mit vier Byte Salz), antwortet play mit md5 und dem hexadezimalen MD5 aus
// dem hexadezimalen MD5 von Passwort und Benutzer der Startup-Daten, gefolgt
// vom Salz; ohne user in den Startup-Daten ist der Benutzer leer, und ein
// anderer Benutzer, ein anderes Salz und ein anderes Passwort geben eine
// andere Antwort (LH-FA-20.a *Verfahren*).
func TestAnmeldungMD5(t *testing.T) {
	salz := []byte{1, 2, 3, 4}
	for _, f := range []struct {
		name, user, passwort string
		salz                 []byte
		want                 string
	}{
		{"Benutzer u", "u", "GEHEIMpw", salz, "md51a3979a40cfaef91fe04aec853d6cbe9"},
		{"anderer Benutzer", "anna", "GEHEIMpw", salz, "md5e27c805331ad780908591df5c2feadf4"},
		{"ohne Benutzer", "", "GEHEIMpw", salz, "md58eb3e4f9190c1c892463b5151bb18d93"},
		{"anderes Salz", "u", "GEHEIMpw", []byte{5, 6, 7, 8}, "md565a7b07442432cf6d96b0d1dbed379db"},
		{"Passwort unverändert", "u", "  ä%41$ ", salz, "md51d548beef516b0ca8d6baec919633988"},
	} {
		t.Run(f.name, func(t *testing.T) {
			addr, ergebnis := einspielServer(t, anmMD5(f.salz...), anmBereit(t))
			startup := map[string]string{"user": f.user}
			if f.user == "" {
				startup = map[string]string{"database": "d"}
			}
			if err := anmeldenMit(t, addr, startup, f.passwort); err != nil {
				t.Fatalf("MD5: %v", err)
			}
			anmGesendet(t, geschlossen(t, ergebnis), anmPasswort(f.want))
			if got := postgres.Md5Antwort(f.passwort, f.user, f.salz); got != f.want {
				t.Errorf("md5Antwort %s, erwartet %s", got, f.want)
			}
		})
	}
}

// Abdeckung: LH-FA-20/Happy — Beweis und Serversignatur von SCRAM-SHA-256
// stimmen mit dem Beispiel aus RFC 7677 überein (Benutzer user, Passwort
// pencil, 4096 Iterationen) (LH-FA-20.a *SCRAM-Austausch*).
func TestScramBeweisRFC7677(t *testing.T) {
	salz, _ := base64.StdEncoding.DecodeString("W22ZaJ0SNY7soEsUEjb6gQ==")
	auth := "n=user,r=rOprNGfwEbeRWgbNEkqO,r=rOprNGfwEbeRWgbNEkqO%hvYDpWUa2RaTCAfuxFIlj)hNlF$k0,s=W22ZaJ0SNY7soEsUEjb6gQ==,i=4096,c=biws,r=rOprNGfwEbeRWgbNEkqO%hvYDpWUa2RaTCAfuxFIlj)hNlF$k0"
	beweis, sig, err := postgres.ScramBeweis("pencil", salz, 4096, auth)
	if err != nil {
		t.Fatal(err)
	}
	if got := base64.StdEncoding.EncodeToString(beweis); got != "dHzbZapWIk4jUhN+Ute9ytag9zjfMHgsqmmiz7AndVQ=" {
		t.Errorf("Beweis %s", got)
	}
	if got := base64.StdEncoding.EncodeToString(sig); got != "6rriTRBi23WpRR/wtup+mMhUZUn/dB5nLTJRsjl95G4=" {
		t.Errorf("Serversignatur %s", got)
	}
}

// Abdeckung: LH-FA-20/Boundary, LH-FA-20/Negative — die erste Nachricht des
// Servers im SCRAM-Austausch besteht genau aus r=, s= und i= in dieser
// Reihenfolge; r= beginnt mit der Nonce von play und ist länger, s= ist
// gültiges Base64 (ohne Zeilenumbruch), i= eine Dezimalzahl aus Ziffern ohne
// führende Null von 1 bis 10 000 000 (LH-FA-20.a *SCRAM-Austausch*).
func TestScramServerErste(t *testing.T) {
	const r = "r=NONCEserver"
	const s = "s=c2FsegAB"
	gut := func(i string) string { return r + "," + s + "," + i }
	for _, f := range []struct {
		name, text string
		gueltig    bool
	}{
		{"Iterationen 1", gut("i=1"), true},
		{"Iterationen 4096", gut("i=4096"), true},
		{"Iterationen genau an der Obergrenze", gut("i=10000000"), true},
		{"Iterationen 0", gut("i=0"), false},
		{"Iterationen über der Obergrenze", gut("i=10000001"), false},
		{"Iterationen weit über der Grenze", gut("i=99999999999999999999999"), false},
		{"Iterationen 2 hoch 64 plus 5", gut("i=18446744073709551621"), false},
		{"Iterationen negativ", gut("i=-1"), false},
		{"Iterationen mit Vorzeichen", gut("i=+5"), false},
		{"Iterationen leer", gut("i="), false},
		{"Iterationen mit führenden Nullen", gut("i=0004096"), false},
		{"Iterationen mit einer führenden Null", gut("i=01"), false},
		{"Iterationen 10 mit führender Null", gut("i=010"), false},
		{"Iterationen keine Zahl", gut("i=abc"), false},
		{"Iterationen mit Leerraum dahinter", gut("i=4096 "), false},
		{"Iterationen mit Exponent", gut("i=1e3"), false},
		{"zwei Attribute", r + "," + s, false},
		{"vier Attribute", gut("i=1") + ",x=1", false},
		{"Komma am Ende", gut("i=1") + ",", false},
		{"andere Reihenfolge", s + "," + r + ",i=1", false},
		{"r ohne eigene Nonce", "r=ANDERE" + "server" + "," + s + ",i=1", false},
		{"r gleich der eigenen Nonce", "r=NONCE," + s + ",i=1", false},
		{"r kürzer als die eigene Nonce", "r=NONC," + s + ",i=1", false},
		{"s kein Base64", r + ",s=!!!,i=1", false},
		{"s ohne Auffüllung", r + ",s=c2FsegA,i=1", false},
		{"s mit Zeilenumbruch", r + ",s=c2Fs\negAB,i=1", false},
		{"r= fehlt", "x=NONCEserver," + s + ",i=1", false},
		{"s= fehlt", r + ",x=c2FsegAB,i=1", false},
		{"i= fehlt", r + "," + s + ",x=1", false},
	} {
		t.Run(f.name, func(t *testing.T) {
			servernonce, salz, n, grund := postgres.LeseServerErste("NONCE", f.text)
			if f.gueltig {
				if grund != "" || servernonce != "NONCEserver" || !bytes.Equal(salz, []byte("salz\x00\x01")) || n < 1 {
					t.Fatalf("gültige Nachricht abgelehnt oder falsch gelesen: %q %q %d %q", servernonce, salz, n, grund)
				}
				return
			}
			if grund == "" {
				t.Fatalf("Nachricht angenommen: %q %q %d", servernonce, salz, n)
			}
		})
	}
	if _, _, n, _ := postgres.LeseServerErste("NONCE", gut("i=10000000")); n != 10_000_000 {
		t.Errorf("Iterationen %d, erwartet 10000000", n)
	}
}

// scramRunde sind die Nachrichten eines SCRAM-Austauschs aus Sicht des Servers,
// unabhängig vom Code von play berechnet (RFC 5802): die erste Nachricht des
// Servers, die erwartete Antwort von play und der Abschluss.
type scramRunde struct {
	erste, antwort, abschluss string
}

func b64(b []byte) string { return base64.StdEncoding.EncodeToString(b) }

// anmScramRunde rechnet die Runde für passwort, salz, iterationen und die
// Nonce des Servers (Anhang an die Nonce von play).
func anmScramRunde(t *testing.T, passwort string, salz []byte, iterationen int, serverTeil string) scramRunde {
	t.Helper()
	return anmScramRundeFuer(t, anmNonce, passwort, salz, iterationen, serverTeil)
}

// anmScramRundeFuer rechnet die Runde für die Nonce nonce von play.
func anmScramRundeFuer(t *testing.T, nonce, passwort string, salz []byte, iterationen int, serverTeil string) scramRunde {
	t.Helper()
	servernonce := nonce + serverTeil
	erste := fmt.Sprintf("r=%s,s=%s,i=%d", servernonce, b64(salz), iterationen)
	ohneBeweis := "c=biws,r=" + servernonce
	auth := "n=,r=" + nonce + "," + erste + "," + ohneBeweis
	gesalzen, err := pbkdf2.Key(sha256.New, passwort, salz, iterationen, 32)
	if err != nil {
		t.Fatal(err)
	}
	mac := func(schluessel []byte, text string) []byte {
		h := hmac.New(sha256.New, schluessel)
		h.Write([]byte(text))
		return h.Sum(nil)
	}
	clientKey := mac(gesalzen, "Client Key")
	gespeichert := sha256.Sum256(clientKey)
	signatur := mac(gespeichert[:], auth)
	beweis := make([]byte, len(clientKey))
	for i := range beweis {
		beweis[i] = clientKey[i] ^ signatur[i]
	}
	return scramRunde{
		erste:     erste,
		antwort:   ohneBeweis + ",p=" + b64(beweis),
		abschluss: "v=" + b64(mac(mac(gesalzen, "Server Key"), auth)),
	}
}

// anmCode ist eine Nachricht R mit code und Text.
func anmCode(code uint32, text string) []byte { return anmeldung(code, []byte(text)...) }

// anmAntwort ist die SASLResponse von play mit text.
func anmAntwort(text string) clientNachricht { return clientNachricht{'p', []byte(text)} }

// Abdeckung: LH-FA-20/Happy — play meldet sich mit SCRAM-SHA-256 an: Es wählt
// das Verfahren in genau dieser Schreibweise aus der Liste, auch wenn es nicht
// das erste ist, sendet als erste Nachricht den Kopf n,, den leeren Benutzer
// und die Nonce aus 18 Bytes (Base64), antwortet auf die erste Nachricht des
// Servers mit c=biws, der Nonce des Servers und dem Beweis, prüft die
// Signatur im Abschluss und sendet danach nichts (LH-FA-20.a *Verfahren*,
// *SCRAM-Austausch*); mit 1 Iteration und mit einem Salz, das kein Text ist,
// ebenso.
func TestAnmeldungScram(t *testing.T) {
	for _, f := range []struct {
		name         string
		liste        []string
		salz         []byte
		iterationen  int
		serverTeil   string
		passwort     string
		wirdGesendet bool
	}{
		{"SCRAM allein", []string{"SCRAM-SHA-256"}, []byte("salzsalzsalz"), 4096, "SERVER", "GEHEIMpw", true},
		{"SCRAM nach PLUS", []string{"SCRAM-SHA-256-PLUS", "SCRAM-SHA-256"}, []byte("salzsalzsalz"), 4096, "SERVER", "GEHEIMpw", true},
		{"eine Iteration", []string{"SCRAM-SHA-256"}, []byte{0, 1, 2, 255}, 1, "x", "  ä%41$ ", true},
	} {
		t.Run(f.name, func(t *testing.T) {
			r := anmScramRunde(t, f.passwort, f.salz, f.iterationen, f.serverTeil)
			l, err := anmLauf(t, f.passwort, anmSasl(f.liste...),
				anmCode(11, r.erste), anmKette(anmCode(12, r.abschluss), anmBereit(t)))
			if err != nil {
				t.Fatalf("SCRAM: %v", err)
			}
			anmGesendet(t, l, anmErste(), anmAntwort(r.antwort))
		})
	}
}

// Abdeckung: LH-FA-20/Boundary — ein leeres Salz (s=) in der ersten Nachricht
// des Servers ist gültig: play liest es als leeres Salz und rechnet den Beweis
// darüber, gegen einen Vektor, den Python (hashlib, hmac) für das Passwort
// GEHEIMpw, die feste Nonce, 4096 Iterationen und die Servernonce SERVER
// gerechnet hat (LH-FA-20.a *SCRAM-Austausch*).
func TestAnmeldungScramLeeresSalz(t *testing.T) {
	const erste = "r=" + anmNonce + "SERVER,s=,i=4096"
	servernonce, salz, n, grund := postgres.LeseServerErste(anmNonce, erste)
	if grund != "" || servernonce != anmNonce+"SERVER" || len(salz) != 0 || n != 4096 {
		t.Fatalf("leeres Salz nicht gelesen: %q %q %d %q", servernonce, salz, n, grund)
	}
	const antwort = "c=biws,r=" + anmNonce + "SERVER,p=y253gaxpbQlQX+qMApvbQJrWp72N1nTYBsi6Lyeq+r8="
	const abschluss = "v=gN80RGhfm/TVcpyOZOO/BehrjyCE0J9An8DL0PIVzlI="
	l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"),
		anmCode(11, erste), anmKette(anmCode(12, abschluss), anmBereit(t)))
	if err != nil {
		t.Fatalf("SCRAM mit leerem Salz: %v", err)
	}
	anmGesendet(t, l, anmErste(), anmAntwort(antwort))
}

// Abdeckung: LH-FA-20/Negative — jeder Fehler im SCRAM-Austausch ist
// PGR-E4005 und bricht sofort ab, ohne dass play danach etwas sendet oder
// Terminate: eine erste Nachricht des Servers, die nicht passt (zu wenige
// oder zu viele Attribute, andere Reihenfolge, fremde oder gleich lange Nonce,
// ungültiges Salz, Iterationen außerhalb von 1 bis 10 000 000), ein Abschluss
// mit falscher Signatur (auch eine kürzere und eine längere), mit e=, ohne v=,
// mit weiterem Attribut oder ohne gültiges Base64 (LH-FA-20.a
// *SCRAM-Austausch*).
func TestAnmeldungScramFehler(t *testing.T) {
	r := anmScramRunde(t, "GEHEIMpw", []byte("salzsalzsalz"), 4096, "SERVER")
	salz := b64([]byte("salzsalzsalz"))
	servernonce := anmNonce + "SERVER"
	sig, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(r.abschluss, "v="))
	for _, f := range []struct{ name, erste string }{
		{"zwei Attribute", "r=" + servernonce + ",s=" + salz},
		{"vier Attribute", r.erste + ",x=1"},
		{"Komma am Ende", r.erste + ","},
		{"andere Reihenfolge", "s=" + salz + ",r=" + servernonce + ",i=4096"},
		{"fremde Nonce", "r=FREMDE" + servernonce + ",s=" + salz + ",i=4096"},
		{"Nonce nicht länger", "r=" + anmNonce + ",s=" + salz + ",i=4096"},
		{"Salz kein Base64", "r=" + servernonce + ",s=!!!,i=4096"},
		{"Iterationen 0", "r=" + servernonce + ",s=" + salz + ",i=0"},
		{"Iterationen über der Obergrenze", "r=" + servernonce + ",s=" + salz + ",i=10000001"},
		{"Iterationen mit Vorzeichen", "r=" + servernonce + ",s=" + salz + ",i=+4096"},
	} {
		t.Run("erste Nachricht: "+f.name, func(t *testing.T) {
			l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"), anmCode(11, f.erste))
			if code(err) != model.CodeLogin {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
			}
			anmGesendet(t, l, anmErste())
		})
	}
	for _, f := range []struct{ name, abschluss string }{
		{"falsche Signatur", "v=" + b64(bytes.Repeat([]byte{7}, 32))},
		{"Signatur ein Byte kürzer", "v=" + b64(sig[:31])},
		{"Signatur ein Byte länger", "v=" + b64(append(append([]byte{}, sig...), 0))},
		{"e= statt v=", "e=invalid-proof"},
		{"ohne v=", strings.TrimPrefix(r.abschluss, "v=")},
		{"v= mit weiterem Attribut", r.abschluss + ",x=1"},
		{"v= kein Base64", "v=!!!"},
		{"v= leer", "v="},
		{"leer", ""},
	} {
		t.Run("Abschluss: "+f.name, func(t *testing.T) {
			l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"), anmCode(11, r.erste), anmKette(anmCode(12, f.abschluss), anmBereit(t)))
			if code(err) != model.CodeLogin {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
			}
			anmGesendet(t, l, anmErste(), anmAntwort(r.antwort))
		})
	}
}

// Abdeckung: LH-FA-20/Negative — im SCRAM-Austausch ist jede Nachricht R, die
// an dieser Stelle nicht vorgesehen ist, PGR-E4005: vor der ersten Nachricht
// des Servers ein Abschluss (Code 12), AuthenticationOk, die Anforderung eines
// Verfahrens (3, 5, 10, 7, ein unbekannter Code) und GSSContinue (8); nach der
// Antwort ebenso und eine zweite erste Nachricht (Code 11). Play sendet
// danach nichts (LH-FA-20.a *SCRAM-Austausch*, *Anmelde-Nachrichten*).
func TestAnmeldungScramUnvorgesehen(t *testing.T) {
	r := anmScramRunde(t, "GEHEIMpw", []byte("salzsalzsalz"), 4096, "SERVER")
	stoerer := []struct {
		name string
		msg  []byte
	}{
		{"AuthenticationOk", anmeldung(0)},
		{"Klartext", anmKlartext()},
		{"MD5", anmMD5(1, 2, 3, 4)},
		{"SASL", anmSasl("SCRAM-SHA-256")},
		{"GSSContinue", anmCode(8, "x")},
		{"SSPI", anmeldung(9)},
		{"unbekannt", anmeldung(99)},
	}
	for _, s := range stoerer {
		t.Run("vor der ersten Nachricht: "+s.name, func(t *testing.T) {
			l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"), anmKette(s.msg, anmBereit(t)))
			if code(err) != model.CodeLogin {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
			}
			anmGesendet(t, l, anmErste())
		})
		t.Run("nach der Antwort: "+s.name, func(t *testing.T) {
			l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"), anmCode(11, r.erste), anmKette(s.msg, anmBereit(t)))
			if code(err) != model.CodeLogin {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
			}
			anmGesendet(t, l, anmErste(), anmAntwort(r.antwort))
		})
	}
	t.Run("Abschluss vor der ersten Nachricht", func(t *testing.T) {
		l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"), anmKette(anmCode(12, r.abschluss), anmBereit(t)))
		if code(err) != model.CodeLogin {
			t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
		}
		anmGesendet(t, l, anmErste())
	})
	// Code und Inhalt sind getrennt: Ein Abschluss mit dem Inhalt einer ersten
	// Nachricht (Code 12) und eine erste Nachricht (Code 11) mit dem Inhalt eines
	// Abschlusses sind an der Stelle des anderen ein Fehler, obwohl der Inhalt
	// dort gültig wäre.
	t.Run("Code 12 mit dem Inhalt der ersten Nachricht", func(t *testing.T) {
		l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"), anmKette(anmCode(12, r.erste), anmBereit(t)))
		if code(err) != model.CodeLogin {
			t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
		}
		anmGesendet(t, l, anmErste())
	})
	t.Run("Code 11 mit dem Inhalt des Abschlusses", func(t *testing.T) {
		l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"), anmCode(11, r.erste), anmKette(anmCode(11, r.abschluss), anmBereit(t)))
		if code(err) != model.CodeLogin {
			t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
		}
		anmGesendet(t, l, anmErste(), anmAntwort(r.antwort))
	})
	t.Run("zweite erste Nachricht nach der Antwort", func(t *testing.T) {
		l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"), anmCode(11, r.erste), anmKette(anmCode(11, r.erste), anmBereit(t)))
		if code(err) != model.CodeLogin {
			t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
		}
		anmGesendet(t, l, anmErste(), anmAntwort(r.antwort))
	})
}

// Abdeckung: LH-FA-20/Negative — nach dem Abschluss des SCRAM-Austauschs ist
// eine weitere Fortsetzung (Code 11, 12, 8) und eine weitere Anforderung (3, 5,
// 10) PGR-E4002, nicht PGR-E4005; play sendet danach nichts (LH-FA-20.a
// *Anmelde-Nachrichten*, *Verfahren*).
func TestAnmeldungNachScramAbschluss(t *testing.T) {
	r := anmScramRunde(t, "GEHEIMpw", []byte("salzsalzsalz"), 4096, "SERVER")
	for name, msg := range map[string][]byte{
		"Code 11": anmCode(11, r.erste), "Code 12": anmCode(12, r.abschluss), "Code 8": anmCode(8, "x"),
		"Klartext": anmKlartext(), "MD5": anmMD5(1, 2, 3, 4), "SASL": anmSasl("SCRAM-SHA-256"),
	} {
		t.Run(name, func(t *testing.T) {
			l, err := anmLauf(t, "GEHEIMpw", anmSasl("SCRAM-SHA-256"), anmCode(11, r.erste), anmKette(anmCode(12, r.abschluss), msg, anmBereit(t)))
			if code(err) != model.CodeUpstream {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeUpstream)
			}
			anmGesendet(t, l, anmErste(), anmAntwort(r.antwort))
		})
	}
}

// Abdeckung: LH-FA-20/Negative — verlangt der Server ein Passwort und play hat
// keines, ist das PGR-E4005, ohne dass play etwas sendet: für Klartext, MD5 und
// SCRAM-SHA-256 (LH-FA-20.a *Anmeldung*, *Passwort*).
func TestAnmeldungOhnePasswort(t *testing.T) {
	for name, anforderung := range map[string][]byte{
		"Klartext": anmKlartext(), "MD5": anmMD5(1, 2, 3, 4), "SCRAM": anmSasl("SCRAM-SHA-256"),
	} {
		t.Run(name, func(t *testing.T) {
			l, err := anmLauf(t, "", anmKette(anforderung, anmBereit(t)))
			if code(err) != model.CodeLogin {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
			}
			anmNichts(t, l)
		})
	}
}

// Abdeckung: LH-FA-20/Negative — ein Verfahren, das play nicht unterstützt,
// ist PGR-E4005, ohne dass play etwas sendet, auch mit einem Passwort: Kerberos
// (2), SCM (6), GSS (7), SSPI (9), ein unbekannter Code; SASL ohne
// SCRAM-SHA-256 (nur SCRAM-SHA-256-PLUS, leere Liste, andere Schreibweise,
// anderes Verfahren) (LH-FA-20.a *Anmeldung*, *Verfahren*).
func TestAnmeldungNichtUnterstuetzt(t *testing.T) {
	for name, anforderung := range map[string][]byte{
		"Kerberos": anmeldung(2), "SCM": anmeldung(6), "GSS": anmeldung(7), "SSPI": anmeldung(9), "unbekannt": anmeldung(99),
		"nur PLUS":              anmSasl("SCRAM-SHA-256-PLUS"),
		"leere Liste":           anmSasl(),
		"Kleinbuchstaben":       anmSasl("scram-sha-256"),
		"anderes Verfahren":     anmSasl("OAUTHBEARER", "SCRAM-SHA-1"),
		"mit Anhängsel":         anmSasl("SCRAM-SHA-256x"),
		"Präfix des Verfahrens": anmSasl("SCRAM-SHA-25"),
	} {
		t.Run(name, func(t *testing.T) {
			l, err := anmLauf(t, "GEHEIMpw", anmKette(anforderung, anmBereit(t)))
			if code(err) != model.CodeLogin {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
			}
			anmNichts(t, l)
		})
	}
}

// Abdeckung: LH-FA-20/Negative, LH-FA-20/Boundary — lässt sich der Rest der
// Anforderung eines unterstützten Verfahrens nicht lesen, ist das PGR-E4005,
// ohne dass play etwas sendet: Klartext mit einem Byte hinter dem Code, MD5 mit
// weniger oder mehr als vier Byte Salz, SASL ohne NUL am Namen, ohne leeren
// Namen als Ende, mit Bytes dahinter und ohne Liste (LH-FA-20.a
// *Verfahren*, *Anmelde-Nachrichten*).
func TestAnmeldungAnforderungNichtLesbar(t *testing.T) {
	for name, anforderung := range map[string][]byte{
		"Klartext mit Rest": anmeldung(3, 1), "MD5 ohne Salz": anmMD5(), "MD5 mit drei Byte": anmMD5(1, 2, 3),
		"MD5 mit fünf Byte":     anmMD5(1, 2, 3, 4, 5),
		"SASL ohne Liste":       anmeldung(10),
		"SASL ohne NUL":         anmeldung(10, []byte("SCRAM-SHA-256")...),
		"SASL ohne Ende":        anmeldung(10, []byte("SCRAM-SHA-256\x00")...),
		"SASL mit Bytes danach": anmeldung(10, []byte("SCRAM-SHA-256\x00\x00x")...),
		"SASL mit zwei Enden":   anmeldung(10, []byte("SCRAM-SHA-256\x00\x00\x00")...),
	} {
		t.Run(name, func(t *testing.T) {
			l, err := anmLauf(t, "GEHEIMpw", anmKette(anforderung, anmBereit(t)))
			if code(err) != model.CodeLogin {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeLogin)
			}
			anmNichts(t, l)
		})
	}
}

// Abdeckung: LH-FA-20/Negative — play wechselt das Verfahren nicht: Eine
// weitere Anforderung von Klartext, MD5 oder SASL, nachdem es eine Anforderung
// von Klartext oder MD5 beantwortet hat, ist PGR-E4002; es hat dann genau die
// erste beantwortet (LH-FA-20.a *Verfahren*).
func TestAnmeldungWeitereAnforderung(t *testing.T) {
	ersten := map[string]struct {
		anforderung []byte
		antwort     clientNachricht
	}{
		"Klartext": {anmKlartext(), anmPasswort("GEHEIMpw")},
		"MD5":      {anmMD5(1, 2, 3, 4), anmPasswort(postgres.Md5Antwort("GEHEIMpw", "u", []byte{1, 2, 3, 4}))},
	}
	zweiten := map[string][]byte{"Klartext": anmKlartext(), "MD5": anmMD5(1, 2, 3, 4), "SASL": anmSasl("SCRAM-SHA-256")}
	for en, e := range ersten {
		for zn, z := range zweiten {
			t.Run(en+" dann "+zn, func(t *testing.T) {
				l, err := anmLauf(t, "GEHEIMpw", e.anforderung, anmKette(z, anmBereit(t)))
				if code(err) != model.CodeUpstream {
					t.Fatalf("Fehler %v, erwartet %s", err, model.CodeUpstream)
				}
				anmGesendet(t, l, e.antwort)
			})
		}
	}
}

// Abdeckung: LH-FA-20/Negative — nach einer Antwort auf Klartext oder MD5
// bestimmt der Code die Art: Eine weitere Anforderung der Codes 3, 5 oder 10 ist
// PGR-E4002, auch wenn sich ihr Rest nicht lesen ließe; eine Anforderung eines
// Verfahrens, das play nicht unterstützt (Code 7, 9, unbekannt), bleibt
// PGR-E4005 (LH-FA-20.a *Verfahren*).
func TestAnmeldungWeitereAnforderungArt(t *testing.T) {
	ersten := map[string][]byte{"Klartext": anmKlartext(), "MD5": anmMD5(1, 2, 3, 4)}
	zweiten := []struct {
		name string
		z    []byte
		want string
	}{
		{"MD5 mit drei Byte Salz", anmMD5(1, 2, 3), model.CodeUpstream},
		{"MD5 ohne Salz", anmMD5(), model.CodeUpstream},
		{"Klartext mit Rest", anmeldung(3, 1), model.CodeUpstream},
		{"SASL ohne Ende", anmeldung(10, 'S'), model.CodeUpstream},
		{"Code 7", anmeldung(7), model.CodeLogin},
		{"Code 9", anmeldung(9), model.CodeLogin},
		{"unbekannter Code", anmeldung(99), model.CodeLogin},
	}
	for en, e := range ersten {
		for _, z := range zweiten {
			t.Run(en+" dann "+z.name, func(t *testing.T) {
				antwort := anmPasswort("GEHEIMpw")
				if en == "MD5" {
					antwort = anmPasswort(postgres.Md5Antwort("GEHEIMpw", "u", []byte{1, 2, 3, 4}))
				}
				l, err := anmLauf(t, "GEHEIMpw", e, anmKette(z.z, anmBereit(t)))
				if code(err) != z.want {
					t.Fatalf("Fehler %v, erwartet %s", err, z.want)
				}
				anmGesendet(t, l, antwort)
			})
		}
	}
}

// Abdeckung: LH-FA-20/Negative — nach einer Antwort auf Klartext oder MD5 ist
// eine Fortsetzung (Code 8, 11, 12) ohne laufenden Austausch PGR-E4002
// (LH-FA-20.a *Anmelde-Nachrichten*).
func TestAnmeldungFortsetzungOhneAustausch(t *testing.T) {
	for _, code8 := range []uint32{8, 11, 12} {
		t.Run(fmt.Sprint("Code ", code8), func(t *testing.T) {
			l, err := anmLauf(t, "GEHEIMpw", anmKlartext(), anmKette(anmCode(code8, "x"), anmBereit(t)))
			if code(err) != model.CodeUpstream {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeUpstream)
			}
			anmGesendet(t, l, anmPasswort("GEHEIMpw"))
		})
	}
}

// Abdeckung: LH-FA-20/Negative — ein falsches Passwort ist PGR-E4005 (die
// Fehlerantwort des Servers hat die SQLSTATE-Klasse 28), nach Klartext, MD5 und
// SCRAM-SHA-256 gleich; eine Fehlerantwort anderer Klasse, ein Verbindungsende
// und ein ReadyForQuery vor AuthenticationOk nach dem Senden der Antwort sind
// PGR-E4002; die Meldung nennt SQLSTATE und Meldung des Servers und kein
// Passwort (LH-FA-20.a *Aufbau*, *Anmeldung*).
func TestAnmeldungAbgelehnt(t *testing.T) {
	r := anmScramRunde(t, "GEHEIMpw", []byte("salzsalzsalz"), 4096, "SERVER")
	flows := []struct {
		name   string
		aufbau []byte
		// vorher sind die Antworten des Servers, bevor er abweist.
		vorher []clientNachricht
		jeAnf  func(abweisung []byte) [][]byte
	}{
		{"Klartext", anmKlartext(), []clientNachricht{anmPasswort("GEHEIMpw")}, func(a []byte) [][]byte { return [][]byte{a} }},
		{"MD5", anmMD5(1, 2, 3, 4), []clientNachricht{anmPasswort(postgres.Md5Antwort("GEHEIMpw", "u", []byte{1, 2, 3, 4}))}, func(a []byte) [][]byte { return [][]byte{a} }},
		{"SCRAM nach der ersten Nachricht", anmSasl("SCRAM-SHA-256"), []clientNachricht{anmErste()}, func(a []byte) [][]byte { return [][]byte{a} }},
		{"SCRAM nach der Antwort", anmSasl("SCRAM-SHA-256"), []clientNachricht{anmErste(), anmAntwort(r.antwort)}, func(a []byte) [][]byte { return [][]byte{anmCode(11, r.erste), a} }},
	}
	for _, f := range flows {
		for _, a := range []struct {
			name string
			msg  []byte
			code string
			text string
		}{
			{"Klasse 28", anmFehler(t, "28P01", "Passwort falsch"), model.CodeLogin, "28P01 „Passwort falsch“"},
			{"andere Klasse", anmFehler(t, "53300", "zu viele"), model.CodeUpstream, "53300 „zu viele“"},
			{"ReadyForQuery vor AuthenticationOk", kodiert(t, &pgproto3.ReadyForQuery{TxStatus: 'I'}), model.CodeUpstream, ""},
		} {
			t.Run(f.name+": "+a.name, func(t *testing.T) {
				l, err := anmLauf(t, "GEHEIMpw", f.aufbau, f.jeAnf(a.msg)...)
				if code(err) != a.code || !strings.Contains(err.Error(), a.text) || strings.Contains(err.Error(), "GEHEIM") {
					t.Fatalf("Fehler %v, erwartet %s mit %q", err, a.code, a.text)
				}
				anmGesendet(t, l, f.vorher...)
			})
		}
		t.Run(f.name+": Verbindungsende", func(t *testing.T) {
			jeAnf := f.jeAnf(nil)
			jeAnf[len(jeAnf)-1] = []byte{}
			addr, ergebnis := einspielServer(t, f.aufbau, append(jeAnf, nil)...)
			err := anmeldenMit(t, addr, map[string]string{"user": "u"}, "GEHEIMpw")
			if code(err) != model.CodeUpstream {
				t.Fatalf("Fehler %v, erwartet %s", err, model.CodeUpstream)
			}
			anmGesendet(t, geschlossen(t, ergebnis), f.vorher...)
		})
	}
}

// Abdeckung: LH-FA-20/Negative — scheitert das Senden der Antwort auf eine
// Anforderung, ist das PGR-E4002 (LH-FA-20.a *Aufbau*).
func TestAnmeldungSendenScheitert(t *testing.T) {
	client, server := net.Pipe()
	go func() {
		defer server.Close()
		var laenge [4]byte
		if _, err := io.ReadFull(server, laenge[:]); err != nil {
			return
		}
		_, _ = io.ReadFull(server, make([]byte, binary.BigEndian.Uint32(laenge[:])-4))
		_, _ = server.Write(anmKlartext())
	}()
	defer client.Close()
	_ = client.SetDeadline(time.Now().Add(20 * time.Second))
	// Der Grund ist das Senden, nicht das Lesen danach: beide enden in
	// PGR-E4002, die Meldung unterscheidet sie.
	if err := postgres.AufbauMit(client, map[string]string{"user": "u"}, "GEHEIMpw", nil); code(err) != model.CodeUpstream || !strings.Contains(err.Error(), "nicht zu senden") {
		t.Fatalf("Fehler %v, erwartet %s mit „nicht zu senden“", err, model.CodeUpstream)
	}
}

// Abdeckung: LH-FA-20/Boundary — verlangt der Server kein Passwort
// (AuthenticationOk ohne Anforderung), sendet play keines, auch wenn es eines
// hat (LH-FA-20.a *Anmeldung*, *Passwort*).
func TestAnmeldungKeinPasswortVerlangt(t *testing.T) {
	l, err := anmLauf(t, "GEHEIMpw", anmBereit(t))
	if err != nil {
		t.Fatalf("Aufbau: %v", err)
	}
	anmNichts(t, l)
}

// Abdeckung: LH-FA-20/Boundary — Einspielziel.Verbinde reicht Password an die
// Anmeldung: Mit Klartext, MD5 und SCRAM-SHA-256 steht die Session, sie sendet
// nach der Antwort nur Terminate beim Schließen, und ohne Password endet
// Verbinde mit PGR-E4005 (LH-FA-20.a *Anmeldung*).
func TestAnmeldungVerbinde(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	addr, ergebnis := einspielServer(t, anmKlartext(), anmBereit(t))
	s, err := (&postgres.Einspielziel{Address: addr, Password: "GEHEIMpw"}).Verbinde(ctx, map[string]string{"user": "u"})
	if err != nil {
		t.Fatalf("Verbinde: %v", err)
	}
	s.Schliesse()
	l := empfangen(t, ergebnis)
	if fmt.Sprint(l.gelesen) != fmt.Sprint([]clientNachricht{anmPasswort("GEHEIMpw")}) || !bytes.Equal(l.danach, kodiertFrontend(t, &pgproto3.Terminate{})) {
		t.Fatalf("der Client sendet %q und danach %q", l.gelesen, l.danach)
	}
	addr, ergebnis = einspielServer(t, anmKlartext())
	if _, err := (&postgres.Einspielziel{Address: addr}).Verbinde(ctx, map[string]string{"user": "u"}); code(err) != model.CodeLogin {
		t.Fatalf("ohne Password: %v, erwartet %s", err, model.CodeLogin)
	}
	anmNichts(t, geschlossen(t, ergebnis))
}

// Abdeckung: LH-FA-20/Boundary — die Nonce von play besteht aus 18
// Zufallsbytes, Base64 kodiert, und ist je Verbindung eine andere (LH-FA-20.a
// *SCRAM-Austausch*); mit festem Zufall ist sie das Base64 genau dieser Bytes
// (TestAnmeldungScram).
func TestAnmeldungNonce(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var nonces []string
	for range 2 {
		addr, ergebnis := einspielServer(t, anmSasl("SCRAM-SHA-256"), anmFehler(t, "28P01", "nein"))
		if _, err := (&postgres.Einspielziel{Address: addr, Password: "pw"}).Verbinde(ctx, map[string]string{"user": "u"}); code(err) != model.CodeLogin {
			t.Fatalf("Verbinde: %v, erwartet %s", err, model.CodeLogin)
		}
		l := geschlossen(t, ergebnis)
		if len(l.gelesen) != 1 {
			t.Fatalf("der Client sendet %q, erwartet die erste Nachricht", l.gelesen)
		}
		_, daten, _ := bytes.Cut(l.gelesen[0].rumpf, []byte("SCRAM-SHA-256\x00"))
		nonce, ok := strings.CutPrefix(string(daten[4:]), "n,,n=,r=")
		roh, err := base64.StdEncoding.DecodeString(nonce)
		if !ok || err != nil || len(roh) != 18 {
			t.Fatalf("erste Nachricht %q: Nonce %q ist kein Base64 aus 18 Bytes", daten[4:], nonce)
		}
		nonces = append(nonces, nonce)
	}
	if nonces[0] == nonces[1] {
		t.Fatalf("zwei Verbindungen mit derselben Nonce %s", nonces[0])
	}
}

// Abdeckung: LH-FA-20/Boundary — endet ctx (zweites Signal), während play auf
// die Antwort des Servers nach dem Passwort wartet, bricht Verbinde ab: Es
// meldet einen Fehler, schließt die Verbindung und sendet nichts mehr, auch
// kein Terminate; ebenso im SCRAM-Austausch nach der Antwort (LH-FA-20.a
// *Abbruchsignal*, *Abbruch im Aufbau*).
func TestAnmeldungSignal(t *testing.T) {
	for _, f := range []struct {
		name string
		// start liefert Adresse und lauf des Fake-Servers, der nach der Antwort
		// schweigt.
		start func(t *testing.T) (string, <-chan lauf)
		// want ist, was der Client bis dahin gesendet hat; bei SCRAM kennt der
		// Test die Nonce nicht und prüft nur die Zahl der Nachrichten.
		want int
	}{
		{"nach dem Passwort", func(t *testing.T) (string, <-chan lauf) {
			return einspielServer(t, anmKlartext(), []byte{})
		}, 1},
		{"im SCRAM-Austausch", anmScramStumm, 2},
	} {
		t.Run(f.name, func(t *testing.T) {
			addr, ergebnis := f.start(t)
			ctx, cancel := context.WithCancel(context.Background())
			fertig := make(chan error, 1)
			go func() {
				_, err := (&postgres.Einspielziel{Address: addr, Password: "GEHEIMpw"}).Verbinde(ctx, map[string]string{"user": "u"})
				fertig <- err
			}()
			time.Sleep(500 * time.Millisecond)
			select {
			case err := <-fertig:
				t.Fatalf("Verbinde endet vor dem Signal: %v", err)
			default:
			}
			cancel()
			select {
			case err := <-fertig:
				if code(err) != model.CodeUpstream {
					t.Fatalf("Fehler %v, erwartet %s", err, model.CodeUpstream)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("Verbinde endet binnen 5 s nach dem Signal nicht")
			}
			l := geschlossen(t, ergebnis)
			if len(l.gelesen) != f.want || len(l.danach) != 0 {
				t.Fatalf("der Client sendet %q und danach %q, erwartet %d Nachrichten und nichts danach", l.gelesen, l.danach, f.want)
			}
		})
	}
}

// anmScramStumm ist ein Fake-Server, der SCRAM-SHA-256 anbietet, auf die erste
// Nachricht mit einer zu ihrer Nonce passenden Antwort folgt und nach der
// Antwort des Clients schweigt.
func anmScramStumm(t *testing.T) (string, <-chan lauf) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	ergebnis := make(chan lauf, 1)
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_ = conn.SetDeadline(time.Now().Add(20 * time.Second))
		lies := func() (clientNachricht, bool) {
			var kopf [5]byte
			if _, err := io.ReadFull(conn, kopf[:]); err != nil {
				return clientNachricht{}, false
			}
			rumpf := make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)
			_, err := io.ReadFull(conn, rumpf)
			return clientNachricht{kopf[0], rumpf}, err == nil
		}
		var laenge [4]byte
		if _, err := io.ReadFull(conn, laenge[:]); err != nil {
			return
		}
		_, _ = io.ReadFull(conn, make([]byte, binary.BigEndian.Uint32(laenge[:])-4))
		_, _ = conn.Write(anmSasl("SCRAM-SHA-256"))
		erste, ok := lies()
		if !ok {
			return
		}
		_, daten, _ := bytes.Cut(erste.rumpf, []byte("SCRAM-SHA-256\x00"))
		nonce := strings.TrimPrefix(string(daten[4:]), "n,,n=,r=")
		_, _ = conn.Write(anmCode(11, anmScramRundeFuer(t, nonce, "GEHEIMpw", []byte("salzsalzsalz"), 4096, "SERVER").erste))
		zweite, ok := lies()
		if !ok {
			return
		}
		rest, err := io.ReadAll(conn)
		ergebnis <- lauf{gelesen: []clientNachricht{erste, zweite}, danach: rest, ende: err}
	}()
	return l.Addr().String(), ergebnis
}

// Abdeckung: LH-FA-20/Negative — kein Fehler der Anmeldung nennt das Passwort,
// einen daraus abgeleiteten Wert (Beweis, Signatur, Salz, Nonce des Servers,
// Antwort bei MD5) oder den Inhalt einer Nachricht des Servers im Austausch,
// weder im Text noch in der Ursachenkette noch in einer Formatierung (%v, %+v,
// %#v), und das Einspielziel gibt sein Passwort bei einer Formatierung nicht
// aus (LH-FA-20.a *Passwort*).
func TestAnmeldungOhneGeheimnis(t *testing.T) {
	r := anmScramRunde(t, "GEHEIMpw", []byte("salzsalzsalz"), 4096, "SERVER")
	salz := b64([]byte("salzsalzsalz"))
	md5 := postgres.Md5Antwort("GEHEIMpw", "u", []byte{1, 2, 3, 4})
	verboten := []string{"GEHEIM", "MARKERINHALT", salz, "SERVER", anmNonce, strings.TrimPrefix(r.abschluss, "v="), r.antwort[strings.Index(r.antwort, "p=")+2:], md5, "salzsalzsalz"}
	for _, f := range []struct {
		name   string
		aufbau []byte
		jeAnf  [][]byte
	}{
		{"Klartext abgelehnt", anmKlartext(), [][]byte{anmFehler(t, "28P01", "nein")}},
		{"MD5 abgelehnt", anmMD5(1, 2, 3, 4), [][]byte{anmFehler(t, "28P01", "nein")}},
		{"SCRAM erste Nachricht mit Inhalt", anmSasl("SCRAM-SHA-256"), [][]byte{anmCode(11, r.erste+",x=MARKERINHALT")}},
		{"SCRAM falsche Signatur", anmSasl("SCRAM-SHA-256"), [][]byte{anmCode(11, r.erste), anmCode(12, "v="+b64([]byte("MARKERINHALT")))}},
		{"SCRAM Abschluss ohne v=", anmSasl("SCRAM-SHA-256"), [][]byte{anmCode(11, r.erste), anmCode(12, "e=MARKERINHALT")}},
		{"SCRAM unvorgesehene Nachricht", anmSasl("SCRAM-SHA-256"), [][]byte{anmCode(11, r.erste), anmCode(8, "MARKERINHALT")}},
		{"ohne Verfahren", anmSasl("MARKERINHALT"), nil},
	} {
		t.Run(f.name, func(t *testing.T) {
			_, err := anmLauf(t, "GEHEIMpw", f.aufbau, f.jeAnf...)
			if err == nil {
				t.Fatal("Anmeldung ohne Fehler")
			}
			text := fmt.Sprintf("%v|%+v|%#v|%s", err, err, err, err.Error())
			for _, m := range model.Meldungen(err) {
				text += "|" + m.Text
			}
			for _, v := range verboten {
				if strings.Contains(text, v) {
					t.Errorf("der Fehler nennt %q: %s", v, text)
				}
			}
		})
	}
	ziel := &postgres.Einspielziel{Address: "h:1", Password: "GEHEIMpw"}
	for _, verb := range []string{"%v", "%+v", "%#v", "%s", "%q", "%x"} {
		if got := fmt.Sprintf(verb, ziel) + fmt.Sprintf(verb, *ziel); strings.Contains(got, "GEHEIM") || strings.Contains(strings.ToLower(got), "47454845494d7077") {
			t.Errorf("%s des Einspielziels nennt das Passwort: %s", verb, got)
		}
	}
}
