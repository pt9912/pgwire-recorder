package pgwire

import (
	"context"
	"encoding/binary"

	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// codeSASLAbschluss ist der Code der Nachricht R mit SASLFinal: Auf sie erwartet
// der Client keine Antwort (LH-FA-05.b *Vermittlung*).
const codeSASLAbschluss = pgproto3.AuthTypeSASLFinal

// rohAnforderung ist eine Nachricht R des Servers, die der Adapter weder liest
// noch verändert.
type rohAnforderung model.Anmeldeanforderung

// Backend kennzeichnet die Nachricht als Backend-Nachricht.
func (rohAnforderung) Backend() {}

// Decode liest nichts: Die Nachricht wird nur gesendet.
func (rohAnforderung) Decode([]byte) error { return nil }

// Encode schreibt 'R', die Länge, den Code und den Rest der Nachricht.
func (r rohAnforderung) Encode(dst []byte) ([]byte, error) {
	dst = append(dst, 'R')
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(r.Daten)+8))
	dst = binary.BigEndian.AppendUint32(dst, r.Code)
	return append(dst, r.Daten...), nil
}

// anmeldung liefert den Rückruf, über den der Record-Use-Case die Anmeldung mit
// dem Client führt (LH-FA-05.b *Vermittlung*): Die Anforderung des Servers geht
// unverändert an den Client, sein Rumpf `p` unverändert zurück. Eine leere oder
// inhaltlich falsche Antwort gilt wie jede andere; ob sie stimmt, entscheidet
// der Server. Endet der Client die Verbindung, liefert er model.ErrClientEnde.
func (s *Server) anmeldung(be *pgproto3.Backend) model.Anmeldung {
	return func(_ context.Context, a model.Anmeldeanforderung) ([]byte, error) {
		be.Send(rohAnforderung(a))
		if err := be.Flush(); err != nil {
			return nil, model.ErrClientEnde
		}
		if a.Code == codeSASLAbschluss {
			return nil, nil
		}
		rumpf, err := antwort(be)
		if err != nil && s.zwang.Load() {
			// Das Zwangsende hat die Verbindung geschlossen, gleich wie das Lesen scheitert.
			return nil, model.ErrClientEnde
		}
		return rumpf, err
	}
}

// antwort liest die Antwort des Clients auf eine Anforderung. Die Nachricht p
// liest der Adapter als SASLResponse: Sie nimmt jeden Rumpf unverändert, auch
// einen leeren. Terminate und das Ende der Verbindung beenden die Anmeldung
// (model.ErrClientEnde); jede andere Nachricht ist PGR-E6001 und geht nicht an
// den Server.
func antwort(be *pgproto3.Backend) ([]byte, error) {
	_ = be.SetAuthType(pgproto3.AuthTypeSASLContinue)
	msg, err := be.Receive()
	if err != nil {
		if verbindungsende(err) {
			return nil, model.ErrClientEnde
		}
		return nil, model.Errorf(model.CodeUnsupported, err, "Client-Nachricht in der Anmeldung nicht lesbar")
	}
	switch m := msg.(type) {
	case *pgproto3.SASLResponse:
		return append([]byte{}, m.Data...), nil
	case *pgproto3.Terminate:
		return nil, model.ErrClientEnde
	}
	return nil, model.Errorf(model.CodeUnsupported, nil, "Client-Nachricht %T in der Anmeldung wird nicht unterstützt", msg)
}
