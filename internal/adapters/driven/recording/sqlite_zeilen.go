package recording

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Die Richtungen einer Nachricht in der Tabelle message (SPEC-043).
const (
	richtungClient = "client"
	richtungServer = "server"
)

// sessionZeilen sind die Zeilen einer Session in den Tabellen session, interaction
// und message (SPEC-043), fertig für das Einfügen.
type sessionZeilen struct {
	id               int
	startup          string
	serverParameters string
	interaktionen    []interaktionZeile
	nachrichten      []nachrichtZeile
}

type interaktionZeile struct {
	sequence int
	art      model.RequestType
	offsetMS sql.NullInt64
	sql      sql.NullString
}

type nachrichtZeile struct {
	sequence int
	position int
	richtung string
	gruppe   sql.NullInt64
	art      string
	felder   string
}

// zeilenDerSession bildet die Zeilen der Session s. Startup und Server-Parameter
// stehen als JSON-Objekt, ohne Einträge als {}; die Nachrichten einer Interaktion
// zählen ab 1 in der Reihenfolge der Interaktion (SPEC-043).
func zeilenDerSession(s model.Session) (sessionZeilen, error) {
	z := sessionZeilen{id: s.ID}
	var err error
	if z.startup, err = jsonText(s.Startup, ""); err != nil {
		return z, err
	}
	if z.serverParameters, err = jsonText(s.ServerParameters, ""); err != nil {
		return z, err
	}
	for _, i := range s.Interactions {
		iz := interaktionZeile{sequence: i.Sequence, art: i.Request.Type}
		if i.OffsetMS != nil {
			iz.offsetMS = sql.NullInt64{Int64: *i.OffsetMS, Valid: true}
		}
		var ns []nachrichtZeile
		if i.Request.Type == model.RequestExtended {
			ns, err = nachrichtenExtended(i)
		} else {
			iz.sql = sql.NullString{String: i.Request.SQL, Valid: true}
			ns, err = nachrichtenEinfach(i)
		}
		if err != nil {
			return z, err
		}
		z.interaktionen = append(z.interaktionen, iz)
		z.nachrichten = append(z.nachrichten, ns...)
	}
	return z, nil
}

// nachrichtenEinfach bildet die Server-Nachrichten einer einfachen Anfrage; die
// Anfrage selbst steht in interaction.sql.
func nachrichtenEinfach(i model.Interaction) ([]nachrichtZeile, error) {
	var ns []nachrichtZeile
	for _, r := range i.Responses {
		z, err := serverZeile(i.Sequence, len(ns)+1, sql.NullInt64{}, r)
		if err != nil {
			return nil, err
		}
		ns = append(ns, z)
	}
	return ns, nil
}

// nachrichtenExtended bildet die Nachrichten einer Extended-Interaktion: je Gruppe
// (group_no ab 1) zuerst die Client-Nachrichten, dann die Server-Nachrichten;
// position zählt über alle Gruppen und beide Richtungen.
func nachrichtenExtended(i model.Interaction) ([]nachrichtZeile, error) {
	var ns []nachrichtZeile
	for gi, g := range i.Groups {
		gruppe := sql.NullInt64{Int64: int64(gi + 1), Valid: true}
		for _, m := range g.Client {
			felder, err := jsonText(clientMessageToDTO(m), "type")
			if err != nil {
				return nil, err
			}
			ns = append(ns, nachrichtZeile{sequence: i.Sequence, position: len(ns) + 1, richtung: richtungClient, gruppe: gruppe, art: string(m.Type), felder: felder})
		}
		for _, r := range g.Server {
			z, err := serverZeile(i.Sequence, len(ns)+1, gruppe, r)
			if err != nil {
				return nil, err
			}
			ns = append(ns, z)
		}
	}
	return ns, nil
}

// serverZeile bildet die Zeile der Server-Nachricht r.
func serverZeile(sequence, position int, gruppe sql.NullInt64, r model.Response) (nachrichtZeile, error) {
	felder, err := jsonText(responseToDTO(r), "type")
	if err != nil {
		return nachrichtZeile{}, err
	}
	return nachrichtZeile{sequence: sequence, position: position, richtung: richtungServer, gruppe: gruppe, art: string(r.Type), felder: felder}, nil
}

// ergaenzeSession fügt die Zeilen z in einer Transaktion ein (BEGIN IMMEDIATE); bei
// einem Fehler ist die Transaktion zurückgerollt und die Datei unverändert. Alle
// Werte gelangen als Parameter vorbereiteter Anweisungen in die Datenbank.
func ergaenzeSession(ctx context.Context, db *sql.DB, z sessionZeilen) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fuegeEin(ctx, tx, z); err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

func fuegeEin(ctx context.Context, tx *sql.Tx, z sessionZeilen) error {
	if _, err := tx.ExecContext(ctx, `INSERT INTO "session" ("id", "startup", "server_parameters") VALUES (?, ?, ?)`,
		z.id, z.startup, z.serverParameters); err != nil {
		return err
	}
	for _, i := range z.interaktionen {
		if _, err := tx.ExecContext(ctx, `INSERT INTO "interaction" ("session_id", "sequence", "type", "offset_ms", "sql") VALUES (?, ?, ?, ?, ?)`,
			z.id, i.sequence, string(i.art), i.offsetMS, i.sql); err != nil {
			return err
		}
	}
	for _, n := range z.nachrichten {
		if _, err := tx.ExecContext(ctx, `INSERT INTO "message" ("session_id", "sequence", "position", "direction", "group_no", "kind", "fields") VALUES (?, ?, ?, ?, ?, ?, ?)`,
			z.id, n.sequence, n.position, n.richtung, n.gruppe, n.art, n.felder); err != nil {
			return err
		}
	}
	return nil
}
