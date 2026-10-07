//go:build integration

package integration_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// rohAblauf meldet sich über listen ohne Treiber an und sendet die Anfragen
// einzeln als einfache Anfragen. Es liefert jede Server-Nachricht nach dem
// Verbindungsaufbau als Zeile (nachrichtText), so wie der Client sie empfängt,
// in Reihenfolge. Der Test liest PGWire selbst; die PGWire-Bibliothek gehört
// den Adaptern (.a-check.yml).
func rohAblauf(t *testing.T, listen string, anfragen ...string) string {
	t.Helper()
	conn, err := net.DialTimeout("tcp", listen, 10*time.Second)
	if err != nil {
		t.Fatalf("Verbindung über %s: %v", listen, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
	r := bufio.NewReader(conn)
	start := binary.BigEndian.AppendUint32(nil, 196608) // Protokoll 3.0
	start = append(start, "user\x00postgres\x00database\x00postgres\x00\x00"...)
	if _, err := conn.Write(append(binary.BigEndian.AppendUint32(nil, uint32(4+len(start))), start...)); err != nil {
		t.Fatal(err)
	}
	for {
		typ, rumpf, err := leseNachricht(r)
		if err != nil {
			t.Fatalf("Verbindungsaufbau über %s: %v", listen, err)
		}
		if typ == 'E' {
			t.Fatalf("Verbindungsaufbau über %s: %s", listen, lesbar(t, typ, rumpf))
		}
		if typ == 'Z' {
			break
		}
	}
	var b strings.Builder
	for _, q := range anfragen {
		fmt.Fprintf(&b, "> %q\n", q)
		if _, err := conn.Write(clientNachricht('Q', q+"\x00")); err != nil {
			t.Fatal(err)
		}
		for {
			typ, rumpf, err := leseNachricht(r)
			if err != nil {
				fmt.Fprintf(&b, "Ende: %v\n", err)
				return b.String()
			}
			fmt.Fprintf(&b, "%s\n", lesbar(t, typ, rumpf))
			if typ == 'Z' {
				break
			}
		}
	}
	_, _ = conn.Write(clientNachricht('X', ""))
	return b.String()
}

func clientNachricht(typ byte, rumpf string) []byte {
	b := binary.BigEndian.AppendUint32([]byte{typ}, uint32(4+len(rumpf)))
	return append(b, rumpf...)
}

// leseNachricht liest eine Server-Nachricht: Typ, Länge, Rumpf.
func leseNachricht(r *bufio.Reader) (byte, []byte, error) {
	kopf := make([]byte, 5)
	if _, err := io.ReadFull(r, kopf); err != nil {
		return 0, nil, err
	}
	n := binary.BigEndian.Uint32(kopf[1:])
	if n < 4 {
		return 0, nil, fmt.Errorf("Länge %d", n)
	}
	rumpf := make([]byte, n-4)
	_, err := io.ReadFull(r, rumpf)
	return kopf[0], rumpf, err
}

// lesbar ist nachrichtText; eine zu kurze Nachricht beendet den Test mit
// Diagnose.
func lesbar(t *testing.T, typ byte, rumpf []byte) (text string) {
	t.Helper()
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("Nachricht %c %q nicht lesbar: %v", typ, rumpf, p)
		}
	}()
	return nachrichtText(typ, rumpf)
}

// nachrichtText beschreibt eine Server-Nachricht mit Name und Inhalt.
// ErrorResponse und NoticeResponse nennen ihre Felder nach Feldcode sortiert,
// weil die Reihenfolge der Felder nicht Teil der Aufzeichnung ist
// (LH-FA-11.a); jede andere Nachricht nennt ihren Rumpf vollständig.
func nachrichtText(typ byte, rumpf []byte) string {
	cstr := func(b []byte) (string, []byte) {
		i := bytes.IndexByte(b, 0)
		if i < 0 {
			return string(b), nil
		}
		return string(b[:i]), b[i+1:]
	}
	switch typ {
	case 'E', 'N':
		var felder []string
		for len(rumpf) > 0 && rumpf[0] != 0 {
			code := rumpf[0]
			var v string
			v, rumpf = cstr(rumpf[1:])
			felder = append(felder, fmt.Sprintf("%c=%s", code, v))
		}
		sort.Strings(felder)
		name := map[byte]string{'E': "ErrorResponse", 'N': "NoticeResponse"}[typ]
		return name + " " + strings.Join(felder, " | ")
	case 'C':
		tag, _ := cstr(rumpf)
		return "CommandComplete " + tag
	case 'Z':
		return "ReadyForQuery " + string(rumpf)
	case 'S':
		name, rest := cstr(rumpf)
		wert, _ := cstr(rest)
		return "ParameterStatus " + name + "=" + wert
	case 'T':
		var spalten []string
		for n, rest := int(binary.BigEndian.Uint16(rumpf)), rumpf[2:]; n > 0; n-- {
			var name string
			name, rest = cstr(rest)
			spalten = append(spalten, fmt.Sprintf("%s tabelle=%d spalte=%d typ=%d größe=%d mod=%d format=%d", name,
				binary.BigEndian.Uint32(rest), binary.BigEndian.Uint16(rest[4:]), binary.BigEndian.Uint32(rest[6:]),
				int16(binary.BigEndian.Uint16(rest[10:])), int32(binary.BigEndian.Uint32(rest[12:])), binary.BigEndian.Uint16(rest[16:])))
			rest = rest[18:]
		}
		return "RowDescription " + strings.Join(spalten, " | ")
	case 'D':
		var werte []string
		for n, rest := int(binary.BigEndian.Uint16(rumpf)), rumpf[2:]; n > 0; n-- {
			l := int32(binary.BigEndian.Uint32(rest))
			rest = rest[4:]
			if l < 0 {
				werte = append(werte, "NULL")
				continue
			}
			werte = append(werte, fmt.Sprintf("%q", rest[:l]))
			rest = rest[l:]
		}
		return "DataRow " + strings.Join(werte, " ")
	case 'I':
		return "EmptyQueryResponse"
	default:
		return fmt.Sprintf("Nachricht %c %q", typ, rumpf)
	}
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
// Transaktion mit ReadyForQuery im Status E und anschließendem ROLLBACK, in
// einer Folge von Interaktionen, die mit fehlerfreien beginnt und endet, zeigen im Replay dieselben
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
		`> "SELECT 1/0"`, "ErrorResponse C=22012 | ", " | M=division by zero | ", "ReadyForQuery I",
		`> "SELECT 2 AS b; SELECT 1/0"`, "RowDescription b ", `DataRow "2"`, "CommandComplete SELECT 1", "ErrorResponse C=22012 | ", "ReadyForQuery I",
		`> "BEGIN"`, "ReadyForQuery T",
		`> "SELECT 1/0"`, "ErrorResponse C=22012 | ", "ReadyForQuery E",
		`> "SELECT 3"`, "ErrorResponse C=25P02 | ", "ReadyForQuery E",
		`> "ROLLBACK"`, "CommandComplete ROLLBACK", "ReadyForQuery I",
		`> "SELECT 4 AS c"`, "CommandComplete SELECT 1", "ReadyForQuery I",
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
		`DataRow "1" "z1"`, `DataRow "2" "z2"`, `DataRow "3" "z3"`, "CommandComplete SELECT 3",
		`> "SELECT 1 AS a WHERE false"`, "RowDescription a ", "CommandComplete SELECT 0",
		"CommandComplete CREATE TABLE", "CommandComplete INSERT 0 2", "CommandComplete UPDATE 1",
		"RowDescription a ", `DataRow "1" NULL`, "CommandComplete SELECT 1", "RowDescription b ", `DataRow "x" "2.5"`, "CommandComplete SELECT 1",
		"NoticeResponse ", " | M=eins | ", "NoticeResponse ", " | M=zwei | ", "CommandComplete DO",
		"CommandComplete SET", "ParameterStatus application_name=ergebnisarten",
		`> "FETCH ALL FROM c"`, "format=1 | b ", "format=1", `DataRow "\x00\x00\x01\x02" "\x00\xff"`, "CommandComplete FETCH 1",
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
