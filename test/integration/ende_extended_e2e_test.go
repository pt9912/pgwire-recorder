//go:build integration

package integration_test

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// endeFlushSQL ist der Text der Parse-Nachricht der Gruppe mit Ende, die mit Flush schließt;
// die zweite Verbindung findet den Backend in pg_stat_activity an diesem Text.
const endeFlushSQL = "SELECT 3 /* ende-extended-flush */"

// endeAblauf sendet den Ablauf einer Extended-Aufzeichnung über c und liefert die Nachrichten
// der Gruppe mit Ende, die noch nicht gelesen sind, bis zum Ende der Verbindung.
type endeAblauf func(t *testing.T, c net.Conn, r *bufio.Reader) []serverNachricht

// endeSyncAblauf sendet den Ablauf der Session 1 von ende-extended-sync.yaml: eine Interaktion
// mit Sync, die Gruppe mit Flush und die Gruppe mit Ende, die mit Sync schließt (Parse, Bind und
// Execute von beendeSQL, Sync).
func endeSyncAblauf(t *testing.T, c net.Conn, r *bufio.Reader) []serverNachricht {
	t.Helper()
	parseSync(t, c, "SELECT 1", 'S')
	bisTyp(t, r, 'Z')
	parseSync(t, c, "SELECT 2", 'H')
	bisTyp(t, r, '1')
	parseBindExecuteSync(t, c, beendeSQL)
	return bisVerbindungsende(t, r)
}

// endeFlushAblauf sendet den Ablauf von ende-extended-flush.yaml mit endeFlushSQL als Text der
// Gruppe mit Ende: zwei Gruppen, die mit Flush schließen. Hat der Client die Bestätigung der
// zweiten gelesen, ruft es beende auf; erst danach liest es das Ende.
func endeFlushAblauf(beende func(t *testing.T)) endeAblauf {
	return func(t *testing.T, c net.Conn, r *bufio.Reader) []serverNachricht {
		t.Helper()
		parseSync(t, c, "SELECT 2", 'H')
		bisTyp(t, r, '1')
		parseSync(t, c, endeFlushSQL, 'H')
		bisTyp(t, r, '1')
		beende(t)
		return bisVerbindungsende(t, r)
	}
}

// beendeFlushBackend beendet über eine zweite Verbindung zum Upstream den Backend, dessen
// letzte Anfrage endeFlushSQL ist; gibt es nicht genau einen solchen Backend, endet der Test.
func beendeFlushBackend(t *testing.T) {
	t.Helper()
	conn := verbinde(t, os.Getenv("PGR_UPSTREAM"))
	defer schliesse(t, conn)
	sql := "SELECT count(pg_terminate_backend(pid)) FROM pg_stat_activity WHERE query = '" + endeFlushSQL + "' AND pid <> pg_backend_pid()"
	if got := wert(t, conn, sql); got != "1" {
		t.Fatalf("beendete Backends mit der Anfrage %q: %s, erwartet 1", endeFlushSQL, got)
	}
}

// typZeile findet den Wert eines type in Block- wie in Flussschreibweise.
var typZeile = regexp.MustCompile(`\btype: ([a-z_]+)`)

// typenfolge liefert die Werte aller type der Session 1 in der Reihenfolge der Datei, durch
// Leerzeichen getrennt; ab "- id: 2" liest es nicht.
func typenfolge(text string) string {
	if i := strings.Index(text, "- id: 2"); i >= 0 {
		text = text[:i]
	}
	var folge []string
	for _, m := range typZeile.FindAllStringSubmatch(text, -1) {
		folge = append(folge, m[1])
	}
	return strings.Join(folge, " ")
}

// gleicheTypenfolge prüft, dass die Typenfolge der Session 1 der Aufzeichnung die der
// handgeschriebenen Aufzeichnung fixture ist; die Folge steht nicht im Test, die Abweichung
// ist der Befund.
func gleicheTypenfolge(t *testing.T, aufgezeichnet string, fixture []byte) {
	t.Helper()
	got, want := typenfolge(aufgezeichnet), typenfolge(string(fixture))
	if want == "" || got != want {
		t.Fatalf("Typenfolge der Session 1:\naufgezeichnet %s\nFixture       %s", got, want)
	}
}

// aufnehmenEndeExtended zeichnet den Ablauf am Binary gegen PostgreSQL auf, prüft den Client und
// die Aufzeichnung gegen fixture und legt sie mit der Sicht des Clients unter PGR_DATEN ab
// (ende-extended-<art>.yaml und ende-extended-<art>-sicht.txt).
func aufnehmenEndeExtended(t *testing.T, art string, fixture []byte, ablauf endeAblauf) {
	t.Helper()
	output := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
	c, r := rohVerbindung(t, rec.listen)
	ns := ablauf(t, c, r)
	if len(ns) == 0 || strings.Contains(typen(ns), "Z") {
		t.Fatalf("erwartet die Fehlerantwort als letzte Nachricht und kein ReadyForQuery, erhalten %q", typen(ns))
	}
	if letzte := ns[len(ns)-1]; letzte.typ != 'E' || felderDerFehlerantwort(letzte.rumpf)["V"] != "FATAL" || felderDerFehlerantwort(letzte.rumpf)["C"] != "57P01" {
		t.Fatalf("letzte Nachricht %c, erwartet die Fehlerantwort FATAL 57P01: %v", letzte.typ, felderDerFehlerantwort(letzte.rumpf))
	}
	rec.stop(t, 0)
	keineMeldung(t, "Ende der Extended-Interaktion", rec.stderr.String())

	text := lies(t, output)
	if !strings.Contains(text, "version: 2\n") || strings.Count(text, "type: connection_closed") != 1 || strings.Count(text, "- id:") != 1 {
		t.Fatalf("Aufzeichnung ohne Version 2, ohne genau ein Ende oder mit mehr als einer Session:\n%s", text)
	}
	gleicheTypenfolge(t, text, fixture)
	gleicheForm(t, text, fixture)
	if daten := os.Getenv("PGR_DATEN"); daten != "" {
		if err := os.WriteFile(filepath.Join(daten, "ende-extended-"+art+".yaml"), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(daten, "ende-extended-"+art+"-sicht.txt"), []byte(endeSicht(ns)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// nurOhnePostgres lässt den Test in der Phase mit PostgreSQL aus und endet, wenn in der Phase
// ohne PostgreSQL noch eine Instanz unter PGR_UPSTREAM antwortet.
func nurOhnePostgres(t *testing.T) {
	t.Helper()
	if os.Getenv("PGR_OHNE_POSTGRES") != "1" {
		t.Skip("nur in der Phase ohne PostgreSQL")
	}
	if c, err := net.DialTimeout("tcp", os.Getenv("PGR_UPSTREAM"), 2*time.Second); err == nil {
		c.Close()
		t.Fatalf("PostgreSQL unter %s ist noch erreichbar", os.Getenv("PGR_UPSTREAM"))
	}
}

// wiedergabeEndeExtended führt replay --fail-on-unconsumed mit der Aufzeichnung aus, die
// aufnehmenEndeExtended unter PGR_DATEN abgelegt hat, sendet denselben Ablauf und prüft, dass
// der Client dieselbe Sicht erhält wie beim Aufzeichnen und der Lauf mit Exit-Code 0 und ohne
// Meldung endet.
func wiedergabeEndeExtended(t *testing.T, art string, ablauf endeAblauf) {
	t.Helper()
	nurOhnePostgres(t)
	daten := os.Getenv("PGR_DATEN")
	sicht, err := os.ReadFile(filepath.Join(daten, "ende-extended-"+art+"-sicht.txt"))
	if err != nil {
		t.Fatalf("Sicht der ersten Phase: %v", err)
	}
	rep := startProzess(t, "replay", "--input", filepath.Join(daten, "ende-extended-"+art+".yaml"), "--fail-on-unconsumed")
	c, r := rohVerbindung(t, rep.listen)
	ns := ablauf(t, c, r)
	if got := endeSicht(ns); got != string(sicht) {
		t.Fatalf("Replay-Sicht:\n%s\n--- aufgezeichnet:\n%s", got, sicht)
	}
	rep.stop(t, 0)
	keineMeldung(t, "Replay des Endes einer Extended-Interaktion", rep.stderr.String())
}

// Abdeckung: LH-FA-25/Happy, LH-FA-18/Happy, LH-FA-02/Happy — Aufzeichnung am Binary gegen
// PostgreSQL 17 (Gruppe mit Ende, die mit Sync schließt): Beendet der Server die Verbindung in
// Execute von pg_terminate_backend (Parse, Bind, Execute, Sync, 57P01), erhält der Client als
// letzte Nachricht die Fehlerantwort FATAL 57P01 und danach das Ende seiner Verbindung, ohne
// ReadyForQuery; der Lauf endet mit Exit-Code 0 und ohne Meldung; die Aufzeichnung trägt die
// Version 2, genau ein connection_closed und eine Session, die Typenfolge der Session 1 (darin
// vor der Fehlerantwort die Bestätigungen von Parse und Bind und die Datenzeile) und V, S, C und
// M der Fehlerantwort von ende-extended-sync.yaml. Die Aufzeichnung bleibt unter PGR_DATEN für
// TestE2EOhnePostgresEndeExtendedSync liegen.
func TestE2ERecordEndeExtendedSync(t *testing.T) {
	aufnehmenEndeExtended(t, "sync", endeExtendedSyncYAML, endeSyncAblauf)
}

// Abdeckung: LH-FA-25/Happy, LH-FA-18/Happy, LH-FA-02/Happy — Aufzeichnung am Binary gegen
// PostgreSQL 17 (Gruppe mit Ende, die mit Flush schließt, die Gruppe davor ebenso): Beendet eine
// zweite Verbindung den Backend, nachdem der Client die Bestätigung von Parse der Gruppe mit
// Ende gelesen hat (genau ein Backend mit dieser Anfrage in pg_stat_activity), erhält der Client
// als letzte Nachricht die Fehlerantwort FATAL 57P01 und danach das Ende seiner Verbindung,
// ohne ReadyForQuery; der Lauf endet mit Exit-Code 0 und ohne Meldung; die Aufzeichnung trägt
// die Version 2, genau ein connection_closed und eine Session, die Typenfolge der Session 1 und
// V, S, C und M der Fehlerantwort von ende-extended-flush.yaml. Die Aufzeichnung bleibt unter
// PGR_DATEN für TestE2EOhnePostgresEndeExtendedFlush liegen.
func TestE2ERecordEndeExtendedFlush(t *testing.T) {
	aufnehmenEndeExtended(t, "flush", endeExtendedFlushYAML, endeFlushAblauf(beendeFlushBackend))
}

// Abdeckung: LH-FA-25/Happy, LH-FA-18/Happy, LH-FA-03/Happy — Replay ohne PostgreSQL: Aus der
// Aufzeichnung von TestE2ERecordEndeExtendedSync erhält der Client für denselben Ablauf
// dieselben Antworten der Gruppe mit Ende (Typ und Rumpf, bei der Fehlerantwort V, S, C und M)
// und danach das Ende seiner Verbindung; der Lauf endet mit --fail-on-unconsumed mit Exit-Code 0
// und ohne Meldung. Läuft nur in der zweiten Phase des Runners (PGR_OHNE_POSTGRES=1).
func TestE2EOhnePostgresEndeExtendedSync(t *testing.T) {
	wiedergabeEndeExtended(t, "sync", endeSyncAblauf)
}

// Abdeckung: LH-FA-25/Happy, LH-FA-18/Happy, LH-FA-03/Happy — Replay ohne PostgreSQL: Aus der
// Aufzeichnung von TestE2ERecordEndeExtendedFlush erhält der Client für denselben Ablauf
// dieselbe Fehlerantwort (V, S, C und M) und danach das Ende seiner Verbindung; der Lauf endet
// mit --fail-on-unconsumed mit Exit-Code 0 und ohne Meldung. Läuft nur in der zweiten Phase des
// Runners (PGR_OHNE_POSTGRES=1).
func TestE2EOhnePostgresEndeExtendedFlush(t *testing.T) {
	wiedergabeEndeExtended(t, "flush", endeFlushAblauf(func(*testing.T) {}))
}
