//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

// rohAblauf meldet sich über listen ohne Treiber an und sendet die Anfragen
// einzeln als einfache Anfragen. Es liefert jede Server-Nachricht nach dem
// Verbindungsaufbau als Zeile mit Typ und Inhalt, so wie der Client sie
// empfängt, in Reihenfolge.
func rohAblauf(t *testing.T, listen string, anfragen ...string) string {
	t.Helper()
	conn, err := net.DialTimeout("tcp", listen, 10*time.Second)
	if err != nil {
		t.Fatalf("Verbindung über %s: %v", listen, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
	fe := pgproto3.NewFrontend(conn, conn)
	fe.Send(&pgproto3.StartupMessage{ProtocolVersion: pgproto3.ProtocolVersionNumber,
		Parameters: map[string]string{"user": "postgres", "database": "postgres"}})
	if err := fe.Flush(); err != nil {
		t.Fatal(err)
	}
	for {
		msg, err := fe.Receive()
		if err != nil {
			t.Fatalf("Verbindungsaufbau über %s: %v", listen, err)
		}
		if e, ok := msg.(*pgproto3.ErrorResponse); ok {
			t.Fatalf("Verbindungsaufbau über %s: %s %s", listen, e.Code, e.Message)
		}
		if _, ok := msg.(*pgproto3.ReadyForQuery); ok {
			break
		}
	}
	var b strings.Builder
	for _, q := range anfragen {
		fmt.Fprintf(&b, "> %q\n", q)
		fe.Send(&pgproto3.Query{String: q})
		if err := fe.Flush(); err != nil {
			t.Fatal(err)
		}
		for {
			msg, err := fe.Receive()
			if err != nil {
				fmt.Fprintf(&b, "Ende: %v\n", err)
				return b.String()
			}
			j, err := json.Marshal(msg)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Fprintf(&b, "%s\n", j)
			if _, ok := msg.(*pgproto3.ReadyForQuery); ok {
				break
			}
		}
	}
	fe.Send(&pgproto3.Terminate{})
	_ = fe.Flush()
	return b.String()
}

// dreiSichten führt die Anfragen direkt gegen PostgreSQL, über record und über
// replay aus der eben geschriebenen Aufzeichnung aus und verlangt dieselbe
// Sicht des Clients (rohAblauf); record und replay enden mit Exit-Code 0 und
// ohne Warnung. Es liefert die Sicht.
func dreiSichten(t *testing.T, anfragen ...string) string {
	t.Helper()
	upstream := os.Getenv("PGR_UPSTREAM")
	direkt := rohAblauf(t, upstream, anfragen...)
	input := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, upstream, input)
	aufgezeichnet := rohAblauf(t, rec.listen, anfragen...)
	rec.stop(t, 0)
	if aufgezeichnet != direkt {
		t.Fatalf("Sicht beim Aufzeichnen:\n%s\n--- direkt gegen PostgreSQL:\n%s\n--- stderr:\n%s", aufgezeichnet, direkt, rec.stderr.String())
	}
	rep := startProzess(t, "replay", "--input", input)
	wiedergegeben := rohAblauf(t, rep.listen, anfragen...)
	rep.stop(t, 0)
	if wiedergegeben != direkt {
		t.Fatalf("Sicht im Replay:\n%s\n--- direkt gegen PostgreSQL:\n%s\n--- stderr:\n%s", wiedergegeben, direkt, rep.stderr.String())
	}
	for _, p := range []*recorder{rec, rep} {
		if strings.Contains(p.stderr.String(), "PGR-") {
			t.Fatalf("Meldung im Lauf:\n%s", p.stderr.String())
		}
	}
	return direkt
}

// enthaeltInFolge verlangt, dass die Teiltexte in dieser Reihenfolge in sicht
// stehen.
func enthaeltInFolge(t *testing.T, sicht string, teile ...string) {
	t.Helper()
	rest := sicht
	for _, teil := range teile {
		i := strings.Index(rest, teil)
		if i < 0 {
			t.Fatalf("Sicht ohne %q nach dem vorigen Teil:\n%s", teil, sicht)
		}
		rest = rest[i+len(teil):]
	}
}

// Abdeckung: LH-FA-11/Happy, LH-FA-11/Boundary — Abnahmeszenario 6 für einfache
// Anfragen: eine einzelne Fehlerantwort, ein Fehler in der zweiten Anweisung
// einer Anfrage nach dem Ergebnis der ersten und ein Fehler in einer
// Transaktion mit ReadyForQuery im Status E und anschließendem ROLLBACK, jeweils
// zwischen fehlerfreien Interaktionen, zeigen im Replay dieselben
// Server-Nachrichten in derselben Reihenfolge (Ergebnisse, SQLSTATE, Meldung,
// Transaktionsstatus) wie beim Aufzeichnen und wie direkt gegen PostgreSQL.
func TestE2EFehlerreplayEinfach(t *testing.T) {
	sicht := dreiSichten(t,
		"SELECT 1 AS a",
		"SELECT 1/0",
		"SELECT 2 AS b; SELECT 1/0",
		"BEGIN",
		"SELECT 1/0",
		"SELECT 3",
		"ROLLBACK",
		"SELECT 4 AS c",
	)
	enthaeltInFolge(t, sicht,
		`> "SELECT 1/0"`, `"Code":"22012"`, `"Message":"division by zero"`, `"TxStatus":"I"`,
		`> "SELECT 2 AS b; SELECT 1/0"`, `"CommandTag":"SELECT 1"`, `"Code":"22012"`, `"TxStatus":"I"`,
		`> "BEGIN"`, `"TxStatus":"T"`,
		`> "SELECT 1/0"`, `"Code":"22012"`, `"TxStatus":"E"`,
		`> "SELECT 3"`, `"Code":"25P02"`, `"TxStatus":"E"`,
		`> "ROLLBACK"`, `"CommandTag":"ROLLBACK"`, `"TxStatus":"I"`,
		`> "SELECT 4 AS c"`, `"CommandTag":"SELECT 1"`, `"TxStatus":"I"`,
	)
}

// Abdeckung: LH-FA-06/Happy, LH-FA-06/Boundary, LH-FA-12/Happy — einfache
// Anfragen mit mehrzeiliger und leerer Ergebnismenge, Befehlen ohne Zeilen
// (CREATE TABLE, INSERT, UPDATE), mehreren Ergebnismengen einer Anfrage, zwei
// NoticeResponse, ParameterStatus nach SET und Binärwerten über einen
// Binär-Cursor zeigen beim Aufzeichnen dieselben Server-Nachrichten in
// derselben Reihenfolge wie direkt gegen PostgreSQL und im Replay dieselben wie
// beim Aufzeichnen.
func TestE2EErgebnisartenEinfach(t *testing.T) {
	sicht := dreiSichten(t,
		"SELECT g, 'z' || g AS t FROM generate_series(1, 3) g",
		"SELECT 1 AS a WHERE false",
		"CREATE TEMP TABLE t_ergebnis (a int, b text)",
		"INSERT INTO t_ergebnis VALUES (1, 'x'), (2, NULL)",
		"UPDATE t_ergebnis SET b = 'y' WHERE a = 2",
		"SELECT 1 AS a, NULL::text AS n; SELECT 'x' AS b, 2.5::numeric AS c",
		"DO $$BEGIN RAISE NOTICE 'eins'; RAISE NOTICE 'zwei'; END$$",
		"SET application_name = 'ergebnisarten'",
		"BEGIN",
		`DECLARE c BINARY CURSOR FOR SELECT 258::int4 AS i, '\x00ff'::bytea AS b`,
		"FETCH ALL FROM c",
		"COMMIT",
	)
	enthaeltInFolge(t, sicht,
		`"CommandTag":"SELECT 3"`,
		`> "SELECT 1 AS a WHERE false"`, `"RowDescription"`, `"CommandTag":"SELECT 0"`,
		`"CommandTag":"CREATE TABLE"`, `"CommandTag":"INSERT 0 2"`, `"CommandTag":"UPDATE 1"`,
		`"Name":"a"`, `"CommandTag":"SELECT 1"`, `"Name":"b"`, `"CommandTag":"SELECT 1"`,
		`"Message":"eins"`, `"Message":"zwei"`, `"CommandTag":"DO"`,
		`"CommandTag":"SET"`, `"Name":"application_name","Value":"ergebnisarten"`,
		`> "FETCH ALL FROM c"`, `"Format":1`, `"binary":"00000102"`, `"binary":"00ff"`, `"CommandTag":"FETCH 1"`,
	)
}

// Abdeckung: LH-FA-05/Negative, LH-FA-06/Negative — folgt in einer einfachen
// Anfrage auf das Ergebnis der ersten Anweisung eine nicht unterstützte
// Antwort (`COPY … TO STDOUT`), erhält der Client keine Antwort der Anfrage,
// nur die Fehlerantwort mit PGR-E6001; die Session wird nicht aufgezeichnet,
// und der Lauf endet mit Exit-Code 6.
func TestE2ERecordWeitergabeNichtUnterstuetzt(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(rec.listen))
	if err != nil {
		t.Fatal(err)
	}
	results, err := conn.Exec(ctx, "SELECT 1; COPY (SELECT 1) TO STDOUT").ReadAll()
	if err == nil || !strings.Contains(err.Error(), "PGR-E6001") || len(results) != 0 {
		t.Fatalf("erwartet nur PGR-E6001, erhalten %d Ergebnisse, %v", len(results), err)
	}
	_ = conn.Close(ctx)
	rec.stop(t, 6)

	if text := lies(t, output); !strings.Contains(text, "sessions: []") {
		t.Fatalf("verworfene Session aufgezeichnet:\n%s", text)
	}
}
