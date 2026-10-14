package postgres

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// stand ist der Stand der Anmeldung eines Aufbaus im Record-Modus
// (LH-FA-05.b *Reihenfolge des Servers*).
type stand int

const (
	// standOffen: keine Anforderung gelesen.
	standOffen stand = iota
	// standBeantwortet: Klartext-Passwort oder MD5 angefordert und beantwortet.
	standBeantwortet
	// standSASL: Code 10 beantwortet, Code 11 steht aus.
	standSASL
	// standSASLLaeuft: mindestens ein Code 11 beantwortet.
	standSASLLaeuft
	// standSASLAbschluss: Code 12 gelesen, AuthenticationOk steht aus.
	standSASLAbschluss
	// standOk: AuthenticationOk gelesen.
	standOk
)

// folge liefert den Stand nach der Nachricht R mit code; ok ist falsch, wenn sie
// in diesem Stand außerhalb der Reihenfolge liegt.
func (s stand) folge(code uint32) (stand, bool) {
	switch {
	case code == anmeldungOk && (s == standOffen || s == standBeantwortet || s == standSASLAbschluss):
		return standOk, true
	case code == anmeldungKlartext && s == standOffen, code == anmeldungMD5 && s == standOffen:
		return standBeantwortet, true
	case code == anmeldungSASL && s == standOffen:
		return standSASL, true
	case code == anmeldungSASLWeiter && (s == standSASL || s == standSASLLaeuft):
		return standSASLLaeuft, true
	case code == anmeldungSASLAbschluss && s == standSASLLaeuft:
		return standSASLAbschluss, true
	}
	return s, false
}

// inReihenfolge sagt, ob code eine Nachricht R ist, deren Platz in der
// Reihenfolge der Anmeldung feststeht; jeder andere Code ist kein vermitteltes
// Verfahren.
func inReihenfolge(code uint32) bool {
	switch code {
	case anmeldungOk, anmeldungKlartext, anmeldungMD5, anmeldungGSSWeiter, anmeldungSASL, anmeldungSASLWeiter, anmeldungSASLAbschluss:
		return true
	}
	return false
}

// vermittlung führt die Anmeldung eines Aufbaus: Sie reicht jede Anforderung des
// Servers an den Client weiter (rueckruf) und dessen Antwort an den Server,
// ohne etwas zu rechnen oder zu lesen (LH-FA-05.b *Vermittlung*).
type vermittlung struct {
	stand    stand
	rueckruf model.Anmeldung
}

// anmeldecode liefert den Code einer Nachricht R, die pgproto3 liest.
func anmeldecode(msg pgproto3.BackendMessage) (uint32, bool) {
	switch msg.(type) {
	case *pgproto3.AuthenticationOk:
		return anmeldungOk, true
	case *pgproto3.AuthenticationCleartextPassword:
		return anmeldungKlartext, true
	case *pgproto3.AuthenticationMD5Password:
		return anmeldungMD5, true
	case *pgproto3.AuthenticationGSS:
		return pgproto3.AuthTypeGSS, true
	case *pgproto3.AuthenticationGSSContinue:
		return anmeldungGSSWeiter, true
	case *pgproto3.AuthenticationSASL:
		return anmeldungSASL, true
	case *pgproto3.AuthenticationSASLContinue:
		return anmeldungSASLWeiter, true
	case *pgproto3.AuthenticationSASLFinal:
		return anmeldungSASLAbschluss, true
	}
	return 0, false
}

// behandle prüft die Nachricht R gegen die Reihenfolge und gibt eine
// Anforderung an den Client weiter; die Antwort geht unverändert als Nachricht
// p an den Server. Auf Code 12 erwartet der Client keine Antwort.
func (v *vermittlung) behandle(ctx context.Context, fe *pgproto3.Frontend, code uint32, msg pgproto3.BackendMessage) error {
	neu, ok := v.stand.folge(code)
	if !ok {
		return v.abgelehnt(code)
	}
	v.stand = neu
	if code == anmeldungOk {
		return nil
	}
	roh, err := msg.Encode(nil)
	if err != nil {
		return model.Errorf(model.CodeInternal, err, "Anforderung mit Code %d nicht abzubilden", code)
	}
	// Hinter 'R', der Länge und dem Code beginnt der Rest der Nachricht.
	antwort, err := v.rueckruf(ctx, model.Anmeldeanforderung{Code: code, Daten: roh[9:]})
	if err != nil || code == anmeldungSASLAbschluss {
		return err
	}
	fe.Send(rohAntwort(antwort))
	if err := fe.Flush(); err != nil {
		return model.Errorf(model.CodeConnectionLost, err, "Antwort der Anmeldung nicht an den Upstream zu senden")
	}
	return nil
}

// abgelehnt liefert den Fehler einer Nachricht R, die der Recorder nicht
// weitergibt: Nach AuthenticationOk ist jede Nachricht R und sonst jeder Code in
// falscher Reihenfolge PGR-E4002, jeder andere Code ein nicht vermitteltes
// Verfahren (PGR-E6001). Die Meldung nennt den Code, nie den Inhalt.
func (v *vermittlung) abgelehnt(code uint32) error {
	if v.stand != standOk && !inReihenfolge(code) {
		return model.Errorf(model.CodeUnsupported, nil, "der Upstream verlangt ein Anmeldeverfahren (Code %d), das record nicht vermittelt", code)
	}
	return model.Errorf(model.CodeUpstream, nil, "Nachricht der Anmeldung mit Code %d außerhalb der Reihenfolge", code)
}

// lesefehler stuft einen Fehler beim Lesen einer Nachricht des Servers im Aufbau
// ein: Das Ende der Verbindung ist PGR-E4003, eine Nachricht R, die pgproto3
// nicht liest, folgt abgelehnt, jede andere nicht lesbare Nachricht ist
// PGR-E6001.
func (v *vermittlung) lesefehler(fe *pgproto3.Frontend, err error) error {
	if istVerbindungsende(err) {
		return model.Errorf(model.CodeConnectionLost, err, "Verbindung zum Upstream im Verbindungsaufbau beendet")
	}
	if code, ok := unlesbarerCode(fe, err); ok {
		return v.abgelehnt(code)
	}
	return model.Errorf(model.CodeUnsupported, nil, "Nachricht des Upstreams im Verbindungsaufbau nicht lesbar")
}

// unlesbarerCode liefert den Code einer Nachricht R, die pgproto3 nicht liest:
// aus dem Fehler für einen unbekannten Code, sonst aus dem Frontend für Code 6
// und 9, die es als nicht umgesetzt meldet.
func unlesbarerCode(fe *pgproto3.Frontend, err error) (uint32, bool) {
	var code uint32
	if _, e := fmt.Sscanf(err.Error(), "unknown authentication type: %d", &code); e == nil {
		return code, true
	}
	if strings.HasSuffix(err.Error(), "is unimplemented") {
		return fe.GetAuthType(), true
	}
	return 0, false
}

// rohAntwort ist eine Nachricht p, deren Rumpf der Recorder nicht liest und
// nicht verändert.
type rohAntwort []byte

// Frontend kennzeichnet die Nachricht als Frontend-Nachricht.
func (rohAntwort) Frontend() {}

// Decode liest nichts: Die Nachricht wird nur gesendet.
func (rohAntwort) Decode([]byte) error { return nil }

// Encode schreibt 'p', die Länge und den Rumpf.
func (r rohAntwort) Encode(dst []byte) ([]byte, error) {
	dst = append(dst, 'p')
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(r)+4))
	return append(dst, r...), nil
}
