package cli

import (
	"net/netip"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// schema ist der Anfang jeder URL einer Verbindung (LH-FA-17.a *Benannte
// Verbindungen*).
const schema = "postgresql://"

// standardPort gilt, wenn die URL keinen Port nennt.
const standardPort = "5432"

// stueck ist ein Stück eines Teils der URL: wörtlicher Text nach $$ und, wo der
// Teil dekodiert wird, nach der Prozent-Dekodierung, oder mit variable ungleich
// "" ein Platzhalter ${variable}.
type stueck struct {
	text     string
	variable string
}

// teil ist ein Teil der URL als Folge seiner Stücke in der Reihenfolge der URL.
type teil []stueck

// verbindung ist eine geladene Verbindung, zerlegt vor dem Einsetzen
// (LH-FA-17.a *Geheimnisse*): je Teil die Stücke, Platzhalter in der
// Reihenfolge der URL. mitBenutzer und mitPasswort sagen, ob die URL den Teil
// schreibt; ein Passwort ist genau ein Platzhalter. port ist ohne Angabe 5432,
// sslmode ohne Parameter disable.
type verbindung struct {
	name        string
	mitBenutzer bool
	benutzer    teil
	mitPasswort bool
	passwort    teil
	host        teil
	port        teil
	datenbank   teil
	sslmode     string
}

// zerlegung ist der Stand des Zerlegers einer URL: die Stelle für die Meldung
// und die Verbindung, soweit sie gelesen ist.
type zerlegung struct {
	stelle string
	v      verbindung
}

// fehler ist PGR-E2004 an der Verbindung der Zerlegung.
func (z *zerlegung) fehler(grund string) error {
	return fehlerDatei(z.stelle, grund)
}

// zerlegeURL zerlegt und prüft den geschriebenen Text der URL einer
// Verbindung nach LH-FA-17.a in dieser Reihenfolge: Steuerzeichen, Schema,
// Benutzer, Passwort, Host, Port, Datenbank, Parameter, Fragment; je Teil die
// Form seiner Platzhalter vor dem Rest. Ein Fehler ist PGR-E2004 an
// connections.<name>, ein Klartext-Passwort PGR-E2006; keiner nennt einen Wert.
func zerlegeURL(name, text string) (verbindung, error) {
	z := &zerlegung{stelle: unter("connections", name), v: verbindung{name: name, sslmode: "disable"}}
	if strings.IndexFunc(text, unicode.IsControl) >= 0 {
		return verbindung{}, z.fehler("Steuerzeichen in der URL")
	}
	rest, ok := strings.CutPrefix(text, schema)
	if !ok {
		return verbindung{}, z.fehler("Schema ist nicht " + schema)
	}
	ende := strings.IndexAny(rest, "/?#")
	if ende < 0 {
		ende = len(rest)
	}
	if err := z.anmeldung(rest[:ende]); err != nil {
		return verbindung{}, err
	}
	if err := z.pfad(rest[ende:]); err != nil {
		return verbindung{}, err
	}
	return z.v, nil
}

// anmeldung liest den Teil mit Benutzer, Host und Port: Das letzte @ trennt den
// Benutzerteil ab, darin das erste : Benutzer und Passwort.
func (z *zerlegung) anmeldung(text string) error {
	hostPort := text
	if at := strings.LastIndex(text, "@"); at >= 0 {
		if err := z.benutzerteil(text[:at]); err != nil {
			return err
		}
		hostPort = text[at+1:]
	}
	host, hinten, klammern, zu := trenneHost(hostPort)
	h, err := liesTeil(host, true)
	if err != nil {
		return z.fehler("Host: " + err.Error())
	}
	if klammern && !zu {
		return z.fehler("Host mit [ ohne ]")
	}
	if host == "" {
		return z.fehler("Host ist leer")
	}
	if err := hostInhalt(h, klammern); err != nil {
		return z.fehler(err.Error())
	}
	z.v.host = h
	port, mitPort := strings.CutPrefix(hinten, ":")
	if klammern && hinten != "" && !mitPort {
		return z.fehler("hinter ] steht nicht : mit Port")
	}
	return z.port(port, mitPort)
}

// hostInhalt prüft den gelesenen Host (LH-FA-17.a *Benannte Verbindungen*): in
// eckigen Klammern wörtlich eine IPv6-Adresse, auch mit Zone; ohne Klammern in
// den wörtlichen Stücken, nach der Dekodierung, kein Leerraum und keines von
// @ : / ? # [ ] %. Den eingesetzten Wert eines Platzhalters prüft es nicht.
func hostInhalt(h teil, klammern bool) error {
	switch {
	case klammern && h.platzhalter():
		return fehlerText("Host in Klammern mit Platzhalter")
	case klammern && !ipv6(h.woertlich()):
		return fehlerText("Host in Klammern ist keine IPv6-Adresse")
	case !klammern && unzulaessigImHost(h.woertlich()):
		return fehlerText("Host mit Leerraum oder einem der Zeichen @ : / ? # [ ] %")
	}
	return nil
}

// ipv6 meldet, ob text eine IPv6-Adresse in Textform ist, auch als IPv4-Adresse
// in IPv6-Form und mit Zone hinter dem ersten %; eine IPv4-Adresse ist keine.
// Eine Zone ist nicht leer und besteht unzulaessigImHost, ein weiteres % ist
// damit ungültig (LH-FA-17.a, V-122).
func ipv6(text string) bool {
	adresse, zone, mitZone := strings.Cut(text, "%")
	if mitZone && (zone == "" || unzulaessigImHost(zone)) {
		return false
	}
	a, err := netip.ParseAddr(adresse)
	return err == nil && a.Is6()
}

// unzulaessigImHost meldet, ob text ein Zeichen enthält, das ein Host ohne
// Klammern nicht enthält: Leerraum (Unicode White_Space), ein Steuerzeichen
// oder eines von @ : / ? # [ ] %.
func unzulaessigImHost(text string) bool {
	return strings.IndexFunc(text, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) >= 0 ||
		strings.ContainsAny(text, "@:/?#[]%")
}

// benutzerteil liest Benutzer und Passwort; ein leerer Benutzer ist ungültig.
// Ein Passwort, dessen geschriebener Text nicht genau ein ${VAR} ist, ist ein
// Klartext-Passwort (PGR-E2006), auch leer, mit fehlerhaftem Platzhalter oder
// ungültigem Escape; es wird nicht dekodiert.
func (z *zerlegung) benutzerteil(text string) error {
	benutzer, passwort, mitPasswort := strings.Cut(text, ":")
	b, err := liesTeil(benutzer, true)
	if err != nil {
		return z.fehler("Benutzer: " + err.Error())
	}
	if benutzer == "" {
		return z.fehler("Benutzer ist leer")
	}
	z.v.mitBenutzer, z.v.benutzer = true, b
	if mitPasswort {
		m := genauEinPlatzhalter.FindStringSubmatch(passwort)
		if m == nil {
			return klartext(z.stelle, "Klartext-Passwort, erlaubt ist im Passwort nur genau ein Platzhalter ${VAR}")
		}
		z.v.mitPasswort, z.v.passwort = true, teil{{variable: m[1]}}
	}
	return nil
}

// genauEinPlatzhalter ist ein geschriebener Text, der genau ein ${VAR} ist.
var genauEinPlatzhalter = regexp.MustCompile(`^\$\{([A-Za-z_][A-Za-z0-9_]*)\}$`)

// fehlerText ist ein Grund ohne Wert, den ein Hilfsschritt des Zerlegers
// meldet; die Stelle setzt der Aufrufer davor.
type fehlerText string

func (f fehlerText) Error() string { return string(f) }

// trenneHost trennt den Host vom Rest: Ein Host in eckigen Klammern (klammern)
// endet an ], hinten ist der Text danach; fehlt ], ist zu false und host der
// ganze Text hinter [. Sonst endet der Host am ersten :, hinten beginnt mit
// ihm.
func trenneHost(text string) (host, hinten string, klammern, zu bool) {
	if innen, ok := strings.CutPrefix(text, "["); ok {
		host, hinten, zu = strings.Cut(innen, "]")
		return host, hinten, true, zu
	}
	if i := strings.IndexByte(text, ':'); i >= 0 {
		return text[:i], text[i:], false, false
	}
	return text, "", false, false
}

// port prüft den Port, wie geschrieben und nicht dekodiert: ohne Angabe 5432;
// mit Platzhalter sind die wörtlichen Zeichen Ziffern, Form und Wert prüft der
// Start nach dem Einsetzen; sonst Ziffern mit Wert 1 bis 65535.
func (z *zerlegung) port(text string, mitPort bool) error {
	if !mitPort {
		z.v.port = teil{{text: standardPort}}
		return nil
	}
	if text == "" {
		return z.fehler(": ohne Port")
	}
	p, err := liesTeil(text, false)
	if err != nil {
		return z.fehler("Port: " + err.Error())
	}
	if !p.platzhalter() {
		if !portForm(p.woertlich()) {
			return z.fehler("Port ist keine Zahl von 1 bis 65535")
		}
	} else if !ziffern(p.woertlich()) {
		return z.fehler("Port mit Zeichen außer Ziffern")
	}
	z.v.port = p
	return nil
}

// pfad liest, was auf den Teil mit Benutzer, Host und Port folgt: / mit der
// Datenbank bis ? oder #, die Parameter bis #, dann ein Fragment, das ungültig
// ist.
func (z *zerlegung) pfad(text string) error {
	db, ok := strings.CutPrefix(text, "/")
	if !ok {
		return z.fehler("Datenbank fehlt")
	}
	ende := strings.IndexAny(db, "?#")
	if ende < 0 {
		ende = len(db)
	}
	d, err := liesTeil(db[:ende], true)
	if err != nil {
		return z.fehler("Datenbank: " + err.Error())
	}
	if ende == 0 {
		return z.fehler("Datenbank ist leer")
	}
	z.v.datenbank = d
	rest := db[ende:]
	if param, ok := strings.CutPrefix(rest, "?"); ok {
		hinten := strings.IndexByte(param, '#')
		if hinten < 0 {
			hinten = len(param)
		}
		if err := z.parameter(param[:hinten]); err != nil {
			return err
		}
		rest = param[hinten:]
	}
	if rest != "" {
		return z.fehler("Fragment ist ungültig")
	}
	return nil
}

// parameter liest die Parameter, getrennt durch &, in ihrer Reihenfolge; je
// Parameter zuerst der Name, dekodiert und genau in der Schreibweise
// verglichen, ohne Platzhalter; ein Parameter password ist ein
// Klartext-Passwort (PGR-E2006), auch ohne = und mit leerem Wert, sein Wert
// wird weder dekodiert noch geprüft; danach das =, ob der Name bekannt ist (nur
// sslmode), ob er zweimal steht, zuletzt der Wert.
func (z *zerlegung) parameter(text string) error {
	gesehen := false
	for _, p := range strings.Split(text, "&") {
		name, wert, mitWert := strings.Cut(p, "=")
		n, err := liesTeil(name, true)
		if err != nil {
			return z.fehler("Name eines Parameters: " + err.Error())
		}
		if n.platzhalter() {
			return z.fehler("Platzhalter im Namen eines Parameters")
		}
		if n.woertlich() == "password" {
			return klartext(z.stelle, "Parameter password ist ein Klartext-Passwort")
		}
		if !mitWert {
			return z.fehler("Parameter ohne =")
		}
		if n.woertlich() != "sslmode" {
			return z.fehler("unbekannter Parameter")
		}
		if gesehen {
			return z.fehler("Parameter sslmode doppelt")
		}
		gesehen = true
		if err := z.sslmode(wert); err != nil {
			return err
		}
	}
	return nil
}

// sslmode prüft den Wert von sslmode vor dem Einsetzen: dekodiert disable oder
// require; ein Platzhalter an irgendeiner Stelle ist ein ungültiger sslmode.
func (z *zerlegung) sslmode(text string) error {
	w, err := liesTeil(text, true)
	if err != nil {
		return z.fehler("sslmode: " + err.Error())
	}
	if m := w.woertlich(); !w.platzhalter() && (m == "disable" || m == "require") {
		z.v.sslmode = m
		return nil
	}
	return z.fehler("ungültiger sslmode, erlaubt sind disable und require")
}

// platzhalterName ist die Form von VAR in ${VAR} (LH-FA-17.a *Geheimnisse*).
var platzhalterName = regexp.MustCompile(`^\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// liesTeil liest den geschriebenen Text eines Teils von links: $$ ist ein $,
// ${VAR} ein Platzhalter, jedes andere ${ ein ungültiger Platzhalter; erst
// wenn alle Platzhalter die Form haben, dekodiert es mit dekodieren die
// wörtlichen Stücke.
func liesTeil(text string, dekodieren bool) (teil, error) {
	var t teil
	var w strings.Builder
	for i := 0; i < len(text); {
		switch {
		case strings.HasPrefix(text[i:], "$$"):
			w.WriteByte('$')
			i += 2
		case strings.HasPrefix(text[i:], "${"):
			m := platzhalterName.FindStringSubmatch(text[i:])
			if m == nil {
				return nil, fehlerText("ungültiger Platzhalter")
			}
			t = t.mitText(w.String())
			w.Reset()
			t = append(t, stueck{variable: m[1]})
			i += len(m[0])
		default:
			w.WriteByte(text[i])
			i++
		}
	}
	t = t.mitText(w.String())
	if !dekodieren {
		return t, nil
	}
	for i, s := range t {
		if s.variable != "" {
			continue
		}
		d, err := dekodiere(s.text)
		if err != nil {
			return nil, err
		}
		t[i].text = d
	}
	return t, nil
}

// mitText hängt ein wörtliches Stück an, wenn text nicht leer ist.
func (t teil) mitText(text string) teil {
	if text == "" {
		return t
	}
	return append(t, stueck{text: text})
}

// platzhalter meldet, ob der Teil einen Platzhalter enthält.
func (t teil) platzhalter() bool {
	for _, s := range t {
		if s.variable != "" {
			return true
		}
	}
	return false
}

// woertlich ist der wörtliche Text des Teils ohne die Platzhalter.
func (t teil) woertlich() string {
	var b strings.Builder
	for _, s := range t {
		b.WriteString(s.text)
	}
	return b.String()
}

// dekodiere dekodiert die Prozent-Escapes eines wörtlichen Stücks; ein
// ungültiges Escape ist ein Fehler, ebenso ein Ergebnis, das kein gültiges
// UTF-8 ist oder ein Steuerzeichen (C0, U+007F, C1) enthält.
func dekodiere(text string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(text); i++ {
		if text[i] != '%' {
			b.WriteByte(text[i])
			continue
		}
		if i+2 >= len(text) || !hex(text[i+1]) || !hex(text[i+2]) {
			return "", fehlerText("ungültiges Escape")
		}
		n, _ := strconv.ParseUint(text[i+1:i+3], 16, 8)
		b.WriteByte(byte(n))
		i += 2
	}
	out := b.String()
	if !utf8.ValidString(out) {
		return "", fehlerText("Escape ergibt kein gültiges UTF-8")
	}
	if strings.IndexFunc(out, unicode.IsControl) >= 0 {
		return "", fehlerText("Escape ergibt ein Steuerzeichen")
	}
	return out, nil
}

// hex meldet, ob c eine Hexadezimalziffer ist.
func hex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// ziffern meldet, ob text nur aus den Ziffern 0 bis 9 besteht.
func ziffern(text string) bool {
	for _, r := range text {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// portForm meldet, ob text die Form eines Ports hat: Ziffern mit einem Wert
// von 1 bis 65535, führende Nullen erlaubt; eine Zahl über dem Wertebereich
// von int liefert Atoi als dessen Grenze, den leeren Text als 0, beide also
// außerhalb.
func portForm(text string) bool {
	if !ziffern(text) {
		return false
	}
	n, _ := strconv.Atoi(text)
	return n >= 1 && n <= 65535
}

// hostPortForm meldet, ob v die Form host:port hat (LH-FA-17.a), wie
// geschrieben, ohne Dekodierung: ein Host in eckigen Klammern ist eine
// IPv6-Adresse, auch mit Zone hinter %; einer ohne Klammern ist nicht leer und
// besteht unzulaessigImHost; dahinter : und ein Port der Form portForm.
func hostPortForm(v string) bool {
	if innen, ok := strings.CutPrefix(v, "["); ok {
		host, hinten, _ := strings.Cut(innen, "]")
		port, ok := strings.CutPrefix(hinten, ":")
		return ok && ipv6(host) && portForm(port)
	}
	host, port, _ := strings.Cut(v, ":")
	return host != "" && !unzulaessigImHost(host) && portForm(port)
}

// nameFehler ist der Grund, aus dem name kein Name einer Verbindung ist, oder
// "": nicht leer, ohne Steuerzeichen und weder :, @ noch $ (LH-FA-17.a). Der
// Grund nennt den Namen nicht.
func nameFehler(name string) string {
	switch {
	case name == "":
		return "Name einer Verbindung ist leer"
	case strings.IndexFunc(name, unicode.IsControl) >= 0:
		return "Name einer Verbindung mit Steuerzeichen"
	case strings.ContainsAny(name, ":@$"):
		return "Name einer Verbindung mit :, @ oder $"
	}
	return ""
}

// klartext ist PGR-E2006 an einer Verbindung; der Grund nennt keinen Wert.
func klartext(stelle, grund string) error {
	return model.Errorf(model.CodeConfigPassword, nil, "Konfigurationsdatei: %s: %s", stelle, grund)
}

// einsetzen setzt die Variablen der Platzhalter in teile ein, die Teile in der
// gegebenen Reihenfolge, je Teil seine Stücke in der Reihenfolge der URL, und
// liefert je Teil seinen Text (LH-FA-17.a *Geheimnisse*): Ein Platzhalter wird
// einmal durch den Wert aus wert ersetzt, der unverändert, weder dekodiert noch
// erneut ausgewertet, in seinem Teil steht. Ein leerer Wert gilt als nicht
// gesetzt; die erste nicht gesetzte Variable ist PGR-E2005 an
// connections.<Name> mit ihrem Namen, ohne einen Wert.
func (v verbindung) einsetzen(wert func(string) string, teile ...teil) ([]string, error) {
	out := make([]string, len(teile))
	for i, t := range teile {
		var b strings.Builder
		for _, s := range t {
			if s.variable == "" {
				b.WriteString(s.text)
				continue
			}
			w := wert(s.variable)
			if w == "" {
				return nil, model.Errorf(model.CodeConfigVariable, nil, "Konfigurationsdatei: %s: Umgebungsvariable %s eines Platzhalters nicht gesetzt", unter("connections", v.name), s.variable)
			}
			b.WriteString(w)
		}
		out[i] = b.String()
	}
	return out, nil
}

// adresseRecord ist die Adresse host:port, zu der record verbindet
// (LH-FA-17.a), in dieser Reihenfolge: sslmode=require ist PGR-E2004; dann
// werden Host und Port eingesetzt, Benutzer, Passwort und Datenbank nicht; hat
// der Port danach nicht die Form eines Ports, ist das PGR-E2004. Jede Meldung
// nennt connections.<Name> und keinen Wert. wert liefert die Variablen.
func (v verbindung) adresseRecord(wert func(string) string) (string, error) {
	if err := v.ohneTLS(); err != nil {
		return "", err
	}
	hp, err := v.einsetzen(wert, v.host, v.port)
	if err != nil {
		return "", err
	}
	return v.adresse(hp[0], hp[1])
}

// ziel ist, was play aus einer Verbindung nimmt: die Adresse host:port und
// Benutzer, Passwort und Datenbank nach dem Einsetzen; ein Benutzer, den die URL
// nicht schreibt, ist "". mitPasswort sagt, ob die URL ein Passwort schreibt;
// dann ist passwort dessen eingesetzter Wert, sonst "".
type ziel struct {
	adresse, benutzer, datenbank string
	mitPasswort                  bool
	passwort                     string
}

// zielPlay ist das Ziel von play (LH-FA-17.a *Wirkung einer URL*): alle Teile
// werden eingesetzt, in der Reihenfolge der URL Benutzer, Passwort, Host, Port,
// Datenbank, die erste nicht gesetzte Variable ist PGR-E2005; zuletzt der Port
// wie bei adresseRecord. sslmode ändert das Ziel nicht; ob TLS gilt, entscheidet
// upstreamPlay.
func (v verbindung) zielPlay(wert func(string) string) (ziel, error) {
	t, err := v.einsetzen(wert, v.benutzer, v.passwort, v.host, v.port, v.datenbank)
	if err != nil {
		return ziel{}, err
	}
	adresse, err := v.adresse(t[2], t[3])
	if err != nil {
		return ziel{}, err
	}
	z := ziel{adresse: adresse, benutzer: t[0], datenbank: t[4]}
	if v.mitPasswort {
		z.mitPasswort, z.passwort = true, t[1]
	}
	return z, nil
}

// ohneTLS ist PGR-E2004 an connections.<Name>, wenn die Verbindung
// sslmode=require verlangt; record verbindet ohne TLS zum Upstream.
func (v verbindung) ohneTLS() error {
	if v.sslmode == "require" {
		return fehlerDatei(unter("connections", v.name), "sslmode=require ist bei record ungültig, record verbindet ohne TLS zum Upstream")
	}
	return nil
}

// adresse setzt eingesetzten Host und Port zusammen; hat der Port nicht die
// Form eines Ports, ist das PGR-E2004 an connections.<Name>.
func (v verbindung) adresse(host, port string) (string, error) {
	if !portForm(port) {
		return "", fehlerDatei(unter("connections", v.name), "Port nach dem Einsetzen ist keine Zahl von 1 bis 65535")
	}
	return zusammensetzen(host, port), nil
}

// zusammensetzen ist host:port mit dem Port wie gegeben; ein Host mit : steht
// in eckigen Klammern, auch einer, der keine IPv6-Adresse ist (LH-FA-17.a).
func zusammensetzen(host, port string) string {
	if strings.Contains(host, ":") {
		return "[" + host + "]:" + port
	}
	return host + ":" + port
}
