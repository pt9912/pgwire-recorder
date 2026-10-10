package postgres

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Passwort ist das Passwort der Anmeldung am Server (LH-FA-20.a *Passwort*).
// Jede Formatierung gibt festen Text aus, nie den Wert: Ein Aufruf mit %v, %+v,
// %#v, %s, %q oder %x auf einem Passwort, auf einem Einspielziel, einem Zugang,
// einer Anmeldung oder einem SCRAM-Austausch, auch über einen Zeiger, verrät es
// nicht. Die Strukturen, die das Passwort in einem Feld halten, tragen ein
// eigenes Format, weil fmt die Methode eines unexportierten Felds nicht ruft.
type Passwort string

// Format gibt für jedes Verb den festen Text aus.
func (Passwort) Format(f fmt.State, _ rune) {
	_, _ = io.WriteString(f, "***")
}

// Format gibt für jedes Verb den festen Text aus, nie das Passwort. anmeldung
// bettet zugang ein und erbt diese Methode.
func (zugang) Format(f fmt.State, _ rune) {
	_, _ = io.WriteString(f, "zugang{***}")
}

// Format gibt für jedes Verb den festen Text aus, nie das Passwort, die Nonce
// oder die erwartete Signatur.
func (scramAustausch) Format(f fmt.State, _ rune) {
	_, _ = io.WriteString(f, "scramAustausch{***}")
}

// Codes der Anmelde-Nachricht R (LH-FA-20.a *Anmelde-Nachrichten*):
// AuthenticationOk, die Fortsetzungen GSSContinue, SASLContinue und SASLFinal
// und die Anforderungen der drei Verfahren, die play unterstützt; jeder andere
// Code fordert ein Verfahren an, das play nicht unterstützt.
const (
	anmeldungOk            = 0
	anmeldungKlartext      = 3
	anmeldungMD5           = 5
	anmeldungGSSWeiter     = 8
	anmeldungSASL          = 10
	anmeldungSASLWeiter    = 11
	anmeldungSASLAbschluss = 12
)

// zugang ist, was die Anmeldung von außen bekommt: das Passwort und die Quelle
// des Zufalls für die Nonce; ohne zufall gilt crypto/rand.
type zugang struct {
	passwort Passwort
	zufall   func([]byte)
}

// anmeldung ist der Stand der Anmeldung eines Aufbaus (LH-FA-20.a
// *Anmeldung*, *Verfahren*, *Anmelde-Nachrichten*): ob eine Anforderung
// beantwortet ist, ob AuthenticationOk gelesen ist und der laufende
// SCRAM-Austausch.
type anmeldung struct {
	zugang
	benutzer    string
	beantwortet bool
	ok          bool
	scram       *scramAustausch
}

func neueAnmeldung(startup map[string]string, z zugang) *anmeldung {
	return &anmeldung{zugang: z, benutzer: startup["user"]}
}

// nachricht stuft eine Nachricht R nach ihrem Code ein (LH-FA-20.a
// *Anmelde-Nachrichten*) und liefert die Antwort an den Server, nil, wenn
// keine zu senden ist: ohne vollständigen Code PGR-E4002; nach
// AuthenticationOk jede PGR-E4002; im laufenden SCRAM-Austausch gilt dessen
// Folge (scramAustausch.schritt); AuthenticationOk ist ohne Fehler, wenn es
// sich lesen lässt; eine Fortsetzung ohne laufenden Austausch PGR-E4002; jeder
// andere Code ist die Anforderung eines Verfahrens (anforderung).
func (a *anmeldung) nachricht(rumpf []byte) (pgproto3.FrontendMessage, error) {
	if len(rumpf) < 4 {
		return nil, model.Errorf(model.CodeUpstream, nil, "Anmelde-Nachricht ohne vollständigen Code")
	}
	code := binary.BigEndian.Uint32(rumpf)
	switch {
	case a.ok:
		return nil, model.Errorf(model.CodeUpstream, nil, "Anmelde-Nachricht (Code %d) nach AuthenticationOk", code)
	case a.scram != nil && a.scram.laeuft():
		return a.scram.schritt(code, rumpf[4:])
	case code == anmeldungOk:
		if err := (&pgproto3.AuthenticationOk{}).Decode(rumpf); err != nil {
			return nil, model.Errorf(model.CodeUpstream, err, "AuthenticationOk nicht lesbar")
		}
		a.ok = true
		return nil, nil
	case code == anmeldungGSSWeiter || code == anmeldungSASLWeiter || code == anmeldungSASLAbschluss:
		return nil, model.Errorf(model.CodeUpstream, nil, "Fortsetzung einer Anmeldung (Code %d) ohne laufenden Austausch", code)
	}
	return a.anforderung(code, rumpf[4:])
}

// bearbeite stuft eine Nachricht R ein (nachricht) und sendet die Antwort an
// den Server; ein gescheitertes Senden im Aufbau ist PGR-E4002.
func (a *anmeldung) bearbeite(conn net.Conn, rumpf []byte) error {
	antwort, err := a.nachricht(rumpf)
	if err != nil || antwort == nil {
		return err
	}
	kodiert, err := antwort.Encode(nil)
	if err != nil {
		return model.Errorf(model.CodeUpstream, nil, "Antwort der Anmeldung nicht zu kodieren")
	}
	if _, err := conn.Write(kodiert); err != nil {
		return model.Errorf(model.CodeUpstream, err, "Antwort der Anmeldung nicht zu senden")
	}
	return nil
}

// anforderung beantwortet die Anforderung eines Verfahrens mit genau diesem
// Verfahren (LH-FA-20.a *Verfahren*): Klartext, MD5 oder SASL. Eine weitere
// Anforderung dieser drei Codes, nachdem play eine beantwortet hat, ist
// PGR-E4002, auch wenn sich ihr Rest nicht lesen ließe; ein Verfahren, das play
// nicht unterstützt, PGR-E4005, auch nach einer Antwort.
func (a *anmeldung) anforderung(code uint32, rest []byte) (pgproto3.FrontendMessage, error) {
	switch code {
	case anmeldungKlartext, anmeldungMD5, anmeldungSASL:
	default:
		return nil, model.Errorf(model.CodeLogin, nil, "der Server verlangt ein Anmeldeverfahren (Code %d), das play nicht unterstützt", code)
	}
	if a.beantwortet {
		return nil, model.Errorf(model.CodeUpstream, nil, "weitere Anforderung eines Anmeldeverfahrens (Code %d) nach der Antwort", code)
	}
	switch code {
	case anmeldungKlartext:
		return a.klartext(rest)
	case anmeldungMD5:
		return a.md5(rest)
	}
	return a.sasl(rest)
}

// fehlendesPasswort ist PGR-E4005 für eine Anforderung, die das Passwort
// braucht, wenn es keines gibt; play sendet dann nichts.
func (a *anmeldung) fehlendesPasswort() error {
	if a.passwort == "" {
		return model.Errorf(model.CodeLogin, nil, "der Server verlangt ein Passwort, play hat keines")
	}
	return nil
}

// klartext antwortet auf Code 3, dessen Rumpf nur der Code ist, mit dem
// Passwort.
func (a *anmeldung) klartext(rest []byte) (pgproto3.FrontendMessage, error) {
	if len(rest) != 0 {
		return nil, model.Errorf(model.CodeLogin, nil, "Anforderung des Klartext-Passworts nicht lesbar")
	}
	if err := a.fehlendesPasswort(); err != nil {
		return nil, err
	}
	a.beantwortet = true
	return &pgproto3.PasswordMessage{Password: string(a.passwort)}, nil
}

// md5 antwortet auf Code 5, dessen Rumpf der Code und genau vier Byte Salz
// sind, mit md5 und dem hexadezimalen MD5 aus dem hexadezimalen MD5 von
// Passwort und Benutzer, gefolgt vom Salz.
func (a *anmeldung) md5(rest []byte) (pgproto3.FrontendMessage, error) {
	if len(rest) != 4 {
		return nil, model.Errorf(model.CodeLogin, nil, "Anforderung von MD5 nicht lesbar")
	}
	if err := a.fehlendesPasswort(); err != nil {
		return nil, err
	}
	a.beantwortet = true
	return &pgproto3.PasswordMessage{Password: md5Antwort(string(a.passwort), a.benutzer, rest)}, nil
}

// md5Antwort ist "md5" und das hexadezimale MD5 aus dem hexadezimalen MD5 von
// passwort und benutzer, gefolgt vom salz.
func md5Antwort(passwort, benutzer string, salz []byte) string {
	innen := md5.Sum([]byte(passwort + benutzer))
	aussen := md5.Sum(append([]byte(hex.EncodeToString(innen[:])), salz...))
	return "md5" + hex.EncodeToString(aussen[:])
}

// scramName ist das Verfahren, das play wählt, in genau dieser Schreibweise.
const scramName = "SCRAM-SHA-256"

// sasl antwortet auf Code 10 mit dem Beginn des SCRAM-Austauschs, wenn die
// Liste der Verfahren SCRAM-SHA-256 enthält.
func (a *anmeldung) sasl(rest []byte) (pgproto3.FrontendMessage, error) {
	verfahren, lesbar := leseVerfahren(rest)
	switch {
	case !lesbar:
		return nil, model.Errorf(model.CodeLogin, nil, "Liste der SASL-Verfahren nicht lesbar")
	case !slices.Contains(verfahren, scramName):
		return nil, model.Errorf(model.CodeLogin, nil, "der Server bietet kein %s an", scramName)
	}
	if err := a.fehlendesPasswort(); err != nil {
		return nil, err
	}
	a.beantwortet = true
	var erste []byte
	a.scram, erste = neuerScramAustausch(string(a.passwort), a.zufall)
	return &pgproto3.SASLInitialResponse{AuthMechanism: scramName, Data: erste}, nil
}

// leseVerfahren liest die Namen der Verfahren einer Anforderung von SASL: je
// Name mit NUL abgeschlossen, danach ein leerer Name als Ende und kein Byte
// weiter. Sonst ist lesbar falsch.
func leseVerfahren(rest []byte) (namen []string, lesbar bool) {
	for {
		i := bytes.IndexByte(rest, 0)
		if i < 0 {
			return nil, false
		}
		if i == 0 {
			return namen, len(rest) == 1
		}
		namen = append(namen, string(rest[:i]))
		rest = rest[i+1:]
	}
}

// Stufen des SCRAM-Austauschs.
const (
	scramErsteErwartet     = iota + 1 // erste Nachricht gesendet, Code 11 erwartet
	scramAbschlussErwartet            // Antwort gesendet, Code 12 erwartet
	scramBeendet                      // Abschluss geprüft
)

// iterationenMax ist die Obergrenze der Iterationen, die der Server nennen
// darf (LH-FA-20.a *SCRAM-Austausch*).
const iterationenMax = 10_000_000

// scramAustausch ist der Stand eines SCRAM-Austauschs nach RFC 5802 und
// RFC 7677 ohne Channel Binding (LH-FA-20.a *SCRAM-Austausch*).
type scramAustausch struct {
	passwort   string
	nonce      string
	stufe      int
	serverSig  []byte
	authBeginn string
}

// neuerScramAustausch beginnt einen Austausch mit einer Nonce aus 18
// Zufallsbytes, Base64 kodiert, und liefert die erste Nachricht des Clients:
// den Kopf n,, den leeren Benutzer und die Nonce.
func neuerScramAustausch(passwort string, zufall func([]byte)) (*scramAustausch, []byte) {
	roh := make([]byte, 18)
	if zufall == nil {
		_, _ = rand.Read(roh)
	} else {
		zufall(roh)
	}
	s := &scramAustausch{passwort: passwort, nonce: base64.StdEncoding.EncodeToString(roh), stufe: scramErsteErwartet}
	s.authBeginn = "n=,r=" + s.nonce
	return s, []byte("n,," + s.authBeginn)
}

// laeuft meldet, ob der Austausch begonnen und noch nicht beendet ist.
func (s *scramAustausch) laeuft() bool { return s.stufe != scramBeendet }

// schritt wertet eine Nachricht R im laufenden Austausch: Code 11 an der
// ersten Stelle, Code 12 an der zweiten; jede andere ist ein Fehler im
// Austausch (PGR-E4005), auch AuthenticationOk und eine Anforderung.
func (s *scramAustausch) schritt(code uint32, daten []byte) (pgproto3.FrontendMessage, error) {
	switch {
	case s.stufe == scramErsteErwartet && code == anmeldungSASLWeiter:
		return s.antworte(string(daten))
	case s.stufe == scramAbschlussErwartet && code == anmeldungSASLAbschluss:
		return nil, s.pruefeAbschluss(string(daten))
	}
	return nil, scramFehler(fmt.Sprintf("Nachricht (Code %d) an dieser Stelle nicht vorgesehen", code))
}

// scramFehler ist PGR-E4005 mit dem Grund in eigenen Worten.
func scramFehler(grund string) error {
	return model.Errorf(model.CodeLogin, nil, "Fehler im SCRAM-Austausch: %s", grund)
}

// antworte liest die erste Nachricht des Servers und liefert die Antwort des
// Clients mit dem Beweis.
func (s *scramAustausch) antworte(vomServer string) (pgproto3.FrontendMessage, error) {
	servernonce, salz, iterationen, grund := leseServerErste(s.nonce, vomServer)
	if grund != "" {
		return nil, scramFehler(grund)
	}
	ohneBeweis := "c=biws,r=" + servernonce
	beweis, serverSig, err := scramBeweis(s.passwort, salz, iterationen, s.authBeginn+","+vomServer+","+ohneBeweis)
	if err != nil {
		return nil, scramFehler("Schlüsselableitung nicht möglich")
	}
	s.serverSig = serverSig
	s.stufe = scramAbschlussErwartet
	return &pgproto3.SASLResponse{Data: []byte(ohneBeweis + ",p=" + base64.StdEncoding.EncodeToString(beweis))}, nil
}

// pruefeAbschluss prüft den Abschluss des Servers: genau v= und die
// Base64-kodierte Signatur des Servers.
func (s *scramAustausch) pruefeAbschluss(vomServer string) error {
	kodiert, ok := strings.CutPrefix(vomServer, "v=")
	if !ok {
		return scramFehler("Abschluss des Servers ohne Signatur")
	}
	sig, gueltig := dekodiereBase64(kodiert)
	if !gueltig || !hmac.Equal(sig, s.serverSig) {
		return scramFehler("Signatur des Servers stimmt nicht")
	}
	s.stufe = scramBeendet
	return nil
}

// leseServerErste liest die erste Nachricht des Servers: genau r=, s= und i= in
// dieser Reihenfolge, durch Komma getrennt. r= beginnt mit nonce, ist länger
// und trägt nur Zeichen von 0x21 bis 0x7E (printable in RFC 5802, das Komma
// trennt schon), s= ist gültiges Base64 (auch leer), i= eine Dezimalzahl aus
// Ziffern ohne führende Null von 1 bis iterationenMax. Bei einem Fehler nennt
// grund ihn, sonst ist er "".
func leseServerErste(nonce, text string) (servernonce string, salz []byte, iterationen int, grund string) {
	teile := strings.Split(text, ",")
	if len(teile) != 3 {
		return "", nil, 0, "erste Nachricht des Servers hat nicht genau drei Attribute"
	}
	r, rOK := strings.CutPrefix(teile[0], "r=")
	s, sOK := strings.CutPrefix(teile[1], "s=")
	i, iOK := strings.CutPrefix(teile[2], "i=")
	if !rOK || !sOK || !iOK {
		return "", nil, 0, "erste Nachricht des Servers ohne r=, s= und i= in dieser Reihenfolge"
	}
	if !strings.HasPrefix(r, nonce) || len(r) <= len(nonce) {
		return "", nil, 0, "Nonce des Servers beginnt nicht mit der eigenen oder ist nicht länger"
	}
	for i := 0; i < len(r); i++ {
		if r[i] < 0x21 || r[i] > 0x7E {
			return "", nil, 0, "Nonce des Servers trägt ein Zeichen außerhalb der druckbaren ASCII-Zeichen"
		}
	}
	salz, gueltig := dekodiereBase64(s)
	if !gueltig {
		return "", nil, 0, "Salz des Servers ist kein Base64"
	}
	n, gueltig := leseIterationen(i)
	if !gueltig {
		return "", nil, 0, "Iterationszahl des Servers ist keine Zahl von 1 bis 10000000"
	}
	return r, salz, n, ""
}

// dekodiereBase64 dekodiert gültiges Base64 mit Auffüllung, auch mit gesetzten
// Restbits (AB==); Zeilenumbrüche, die die Bibliothek überliest, gehören nicht
// zum Alphabet.
func dekodiereBase64(text string) ([]byte, bool) {
	if strings.ContainsAny(text, "\r\n") {
		return nil, false
	}
	b, err := base64.StdEncoding.DecodeString(text)
	return b, err == nil
}

// leseIterationen liest eine Dezimalzahl aus Ziffern ohne Vorzeichen und ohne
// führende Null von 1 bis iterationenMax; der leere Text und ein Wert darüber
// sind ungültig, ohne zu überlaufen.
func leseIterationen(text string) (int, bool) {
	if strings.HasPrefix(text, "0") {
		return 0, false
	}
	n := 0
	for _, c := range text {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
		if n > iterationenMax {
			return 0, false
		}
	}
	return n, n >= 1
}

// scramBeweis rechnet nach RFC 5802 den Beweis des Clients und die erwartete
// Signatur des Servers aus Passwort, Salz, Iterationen und der Auth-Message.
func scramBeweis(passwort string, salz []byte, iterationen int, authMessage string) (beweis, serverSig []byte, err error) {
	gesalzen, err := pbkdf2.Key(sha256.New, passwort, salz, iterationen, sha256.Size)
	if err != nil {
		return nil, nil, err
	}
	clientKey := hmacSHA256(gesalzen, "Client Key")
	gespeichert := sha256.Sum256(clientKey)
	clientSig := hmacSHA256(gespeichert[:], authMessage)
	beweis = make([]byte, len(clientKey))
	for i := range clientKey {
		beweis[i] = clientKey[i] ^ clientSig[i]
	}
	return beweis, hmacSHA256(hmacSHA256(gesalzen, "Server Key"), authMessage), nil
}

// hmacSHA256 ist HMAC-SHA-256 von text unter schluessel.
func hmacSHA256(schluessel []byte, text string) []byte {
	h := hmac.New(sha256.New, schluessel)
	_, _ = io.WriteString(h, text)
	return h.Sum(nil)
}
