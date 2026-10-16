//go:build integration

package integration_test

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"net"
	"testing"
)

// Die Aufzeichnung ist von Hand geschrieben (Version 2, ohne PostgreSQL) und
// gehört genau einem Test.
//
//go:embed testdata/ende-extended-sync.yaml
var endeExtendedSyncYAML []byte

// rohNachricht ist eine Client-Nachricht: Typ, Länge und Rumpf.
func rohNachricht(typ byte, rumpf ...byte) []byte {
	m := append([]byte{typ, 0, 0, 0, 0}, rumpf...)
	binary.BigEndian.PutUint32(m[1:], uint32(len(m)-1))
	return m
}

// parseSync sendet Parse (unbenanntes Statement, sql, keine Parametertypen) und
// danach abschluss, Sync (S) oder Flush (H).
func parseSync(t *testing.T, c net.Conn, sql string, abschluss byte) {
	t.Helper()
	schreibe(t, c, rohNachricht('P', append([]byte("\x00"+sql+"\x00"), 0, 0)...))
	schreibe(t, c, rohNachricht(abschluss))
}

// Abdeckung: LH-FA-25/Happy, LH-FA-18/Boundary, LH-FA-03/Boundary — replay mit
// einer handgeschriebenen Extended-Aufzeichnung: Der Client erhält für die
// Gruppe mit Ende, die mit Sync schließt, die Nachricht vor dem Ende und die
// aufgezeichnete Fehlerantwort (FATAL), danach das Ende der Verbindung, ohne
// ReadyForQuery, und der Recorder meldet nichts; die Interaktion mit Sync und die
// Gruppe mit Flush davor werden bis ReadyForQuery bzw. ParseComplete gelesen
// (ihren Inhalt prüfen die Tests des Service); die nächste Verbindung erhält die
// nächste Session; mit --fail-on-unconsumed ist der Exit-Code 0.
func TestE2EReplayEndeExtendedSync(t *testing.T) {
	rep := startProzess(t, "replay", "--input", endeDatei(t, endeExtendedSyncYAML), "--fail-on-unconsumed")
	c, r := rohVerbindung(t, rep.listen)
	parseSync(t, c, "SELECT 1", 'S')
	bisTyp(t, r, 'Z')
	parseSync(t, c, "SELECT 2", 'H')
	bisTyp(t, r, '1')
	parseSync(t, c, "SELECT 3", 'S')
	ns := bisVerbindungsende(t, r)
	if typen(ns) != "1E" || !bytes.Contains(ns[1].rumpf, []byte("SFATAL\x00")) || !bytes.Contains(ns[1].rumpf, []byte("C57P01\x00")) {
		t.Fatalf("erwartet die Bestätigung, die aufgezeichnete Fehlerantwort FATAL 57P01 und das Ende, erhalten %q %q", typen(ns), ns)
	}

	c2, r2 := rohVerbindung(t, rep.listen)
	anfrage(t, c2, "SELECT 4;")
	bisTyp(t, r2, 'Z')
	schreibe(t, c2, []byte{'X', 0, 0, 0, 4})
	rep.stop(t, 0)
	keineMeldung(t, "Ende der Extended-Interaktion, nächste Session", rep.stderr.String())
}
