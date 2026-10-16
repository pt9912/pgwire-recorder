//go:build integration

package integration_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// beendeSQL lässt den Server die eigene Verbindung beenden: Er antwortet mit
// RowDescription, DataRow und der Fehlerantwort FATAL 57P01 und schließt dann.
const beendeSQL = "SELECT pg_terminate_backend(pg_backend_pid())"

// felderDerFehlerantwort liefert V, S, C und M aus dem Rumpf einer ErrorResponse; F, L und R
// hängen von der Serverversion ab und zählen nicht.
func felderDerFehlerantwort(rumpf []byte) map[string]string {
	f := map[string]string{}
	for _, feld := range bytes.Split(bytes.TrimRight(rumpf, "\x00"), []byte{0}) {
		if len(feld) == 0 {
			continue
		}
		if k := string(feld[:1]); k == "V" || k == "S" || k == "C" || k == "M" {
			f[k] = string(feld[1:])
		}
	}
	return f
}

// endeSicht beschreibt, was ein Client von der Verbindung sieht: Typ und Rumpf jeder
// Nachricht, bei einer Fehlerantwort nur V, S, C und M.
func endeSicht(ns []serverNachricht) string {
	var b strings.Builder
	for _, n := range ns {
		if n.typ == 'E' {
			f := felderDerFehlerantwort(n.rumpf)
			fmt.Fprintf(&b, "E V=%s S=%s C=%s M=%s\n", f["V"], f["S"], f["C"], f["M"])
			continue
		}
		fmt.Fprintf(&b, "%c %x\n", n.typ, n.rumpf)
	}
	return b.String()
}

// endeFelder liefert V, S, C und M der error_response, die in der Aufzeichnung text
// unmittelbar vor der ersten connection_closed steht, in der Fluss- wie in der
// Blockschreibweise; ohne eine solche Stelle endet der Test.
func endeFelder(t *testing.T, text string) map[string]string {
	t.Helper()
	zeilen := strings.Split(text, "\n")
	f := map[string]string{}
	for i, z := range zeilen {
		if !strings.Contains(z, "type: connection_closed") || i == 0 {
			continue
		}
		if davor := zeilen[i-1]; strings.Contains(davor, "type: error_response") && strings.Contains(davor, "notice: {") {
			inhalt := davor[strings.Index(davor, "notice: {")+len("notice: {"):]
			for _, kv := range strings.Split(strings.TrimSuffix(strings.TrimSpace(inhalt), "}}"), ", ") {
				k, v, _ := strings.Cut(kv, ": ")
				f[k] = strings.Trim(v, `"'`)
			}
			return f
		}
		for j := i - 1; j > 0; j-- {
			w := strings.TrimSpace(zeilen[j])
			if w == "notice:" && strings.Contains(zeilen[j-1], "type: error_response") {
				return f
			}
			k, v, ok := strings.Cut(w, ": ")
			if !ok {
				break
			}
			f[k] = strings.Trim(v, `"'`)
			f = nur(f)
		}
		break
	}
	t.Fatalf("keine error_response unmittelbar vor connection_closed:\n%s", text)
	return nil
}

// nur behält von f V, S, C und M.
func nur(f map[string]string) map[string]string {
	g := map[string]string{}
	for _, k := range []string{"V", "S", "C", "M"} {
		if v, ok := f[k]; ok {
			g[k] = v
		}
	}
	return g
}

// gleicheForm prüft, dass die letzten zwei Antworten der aufgezeichneten Interaktion mit
// Ende (error_response und connection_closed) in V, S, C und M die der handgeschriebenen
// Aufzeichnung fixture tragen, und dass die Fixture alle vier Felder nennt.
func gleicheForm(t *testing.T, aufgezeichnet string, fixture []byte) {
	t.Helper()
	got, want := nur(endeFelder(t, aufgezeichnet)), nur(endeFelder(t, string(fixture)))
	for _, k := range []string{"V", "S", "C", "M"} {
		if want[k] == "" || got[k] != want[k] {
			t.Fatalf("Feld %s der Aufzeichnung %q, der Fixture %q (aufgezeichnet %v, Fixture %v)", k, got[k], want[k], got, want)
		}
	}
}

// endeAntworten führt anfrage über eine Verbindung zum Recorder aus und liefert die
// Nachrichten, die der Client bis zum Ende der Verbindung erhält.
func endeAntworten(t *testing.T, listen, sql string) []serverNachricht {
	t.Helper()
	c, r := rohVerbindung(t, listen)
	anfrage(t, c, sql)
	return bisVerbindungsende(t, r)
}

// pruefeBeendeAntworten prüft die Antworten auf beendeSQL: RowDescription, DataRow und die
// Fehlerantwort FATAL 57P01, danach das Ende der Verbindung ohne ReadyForQuery.
func pruefeBeendeAntworten(t *testing.T, ns []serverNachricht) {
	t.Helper()
	if typen(ns) != "TDE" {
		t.Fatalf("erwartet RowDescription, DataRow, ErrorResponse und das Ende, erhalten %q", typen(ns))
	}
	if f := felderDerFehlerantwort(ns[2].rumpf); f["V"] != "FATAL" || f["S"] != "FATAL" || f["C"] != "57P01" {
		t.Fatalf("Fehlerantwort: %v", f)
	}
}

// endeFormate sind die Aufzeichnungen von Abnahmeszenario 18.
func endeFormate() []string { return []string{"rec.yaml", "rec.sqlite"} }

// Abdeckung: LH-FA-25/Happy, LH-FA-02/Happy, LH-FA-22/Happy — Abnahmeszenario 18 am Binary
// gegen PostgreSQL 17, in YAML und in SQLite: Beendet der Server die Verbindung in einer
// laufenden Anfrage (pg_terminate_backend, 57P01), erhält der Client RowDescription, DataRow
// und die Fehlerantwort FATAL und danach das Ende seiner Verbindung, ohne ReadyForQuery; der
// Lauf endet mit Exit-Code 0 und ohne Meldung; die Aufzeichnung trägt die Fehlerantwort und
// connection_closed als letzte zwei Antworten der Interaktion und die Version 2 (YAML:
// version: 2, SQLite: meta.version 2). Bei YAML hat die Fehlerantwort V, S, C und M der
// handgeschriebenen Aufzeichnungen ende-fatal.yaml und ende-terminate.yaml. Die Aufzeichnung
// bleibt unter PGR_DATEN für TestE2EOhnePostgresEnde liegen.
func TestE2ERecordEndeTerminate(t *testing.T) {
	for _, name := range endeFormate() {
		t.Run(name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), name)
			rec := startProzess(t, "record", append([]string{"--upstream", os.Getenv("PGR_UPSTREAM"), "--output", output}, formatArgs(name)...)...)
			ns := endeAntworten(t, rec.listen, beendeSQL)
			pruefeBeendeAntworten(t, ns)
			rec.stop(t, 0)
			keineMeldung(t, "Ende des Servers", rec.stderr.String())

			text := lies(t, output)
			if strings.HasSuffix(name, ".sqlite") {
				if !strings.HasPrefix(text, sqliteKopf) || !strings.Contains(text, "version2") || strings.Contains(text, "version1") || !strings.Contains(text, "57P01") {
					t.Fatalf("SQLite-Aufzeichnung ohne meta.version 2 oder ohne das Ende: %.60q", text)
				}
			} else {
				if !strings.Contains(text, "version: 2\n") || strings.Count(text, "type: connection_closed") != 1 {
					t.Fatalf("Aufzeichnung ohne Version 2 oder ohne genau ein Ende:\n%s", text)
				}
				gleicheForm(t, text, endeFatalYAML)
				gleicheForm(t, text, endeTerminateYAML)
			}
			if daten := os.Getenv("PGR_DATEN"); daten != "" {
				if err := os.WriteFile(filepath.Join(daten, "ende-"+name), []byte(text), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(daten, "ende-sicht-"+name+".txt"), []byte(endeSicht(ns)), 0o644); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// Abdeckung: LH-FA-25/Happy, LH-FA-03/Happy, LH-FA-22/Happy — Abnahmeszenario 18, Replay ohne
// PostgreSQL: Aus den Aufzeichnungen von TestE2ERecordEndeTerminate (YAML und SQLite) erhält
// der Client für dieselbe Anfrage dieselben Antworten und danach das Ende seiner Verbindung;
// der Lauf endet mit Exit-Code 0 und ohne Meldung. Läuft nur in der zweiten Phase des Runners
// (PGR_OHNE_POSTGRES=1).
func TestE2EOhnePostgresEnde(t *testing.T) {
	if os.Getenv("PGR_OHNE_POSTGRES") != "1" {
		t.Skip("nur in der Phase ohne PostgreSQL")
	}
	if c, err := net.DialTimeout("tcp", os.Getenv("PGR_UPSTREAM"), 2*time.Second); err == nil {
		c.Close()
		t.Fatalf("PostgreSQL unter %s ist noch erreichbar", os.Getenv("PGR_UPSTREAM"))
	}
	daten := os.Getenv("PGR_DATEN")
	for _, name := range endeFormate() {
		t.Run(name, func(t *testing.T) {
			sicht, err := os.ReadFile(filepath.Join(daten, "ende-sicht-"+name+".txt"))
			if err != nil {
				t.Fatalf("Sicht der ersten Phase: %v", err)
			}
			rep := startProzess(t, "replay", "--input", filepath.Join(daten, "ende-"+name))
			ns := endeAntworten(t, rep.listen, beendeSQL)
			pruefeBeendeAntworten(t, ns)
			if got := endeSicht(ns); got != string(sicht) {
				t.Fatalf("Replay-Sicht:\n%s\n--- aufgezeichnet:\n%s", got, sicht)
			}
			rep.stop(t, 0)
			keineMeldung(t, "Replay des Endes", rep.stderr.String())
		})
	}
}

// ruheBisZumEnde ist die Zeit, nach der der Server eine ruhende Sitzung mit idle_session_timeout
// 300 ms beendet hat.
const ruheBisZumEnde = 1500 * time.Millisecond

// Abdeckung: LH-FA-25/Boundary, LH-FA-02/Boundary — Ende im Leerlauf gegen PostgreSQL 17
// (idle_session_timeout 300 ms): Stellt der Client nach dem Ende der Sitzung noch eine Anfrage,
// gehört die Fehlerantwort 57P05 zu ihr: Der Client erhält sie und danach das Ende seiner
// Verbindung, die Aufzeichnung trägt sie mit connection_closed als letzte zwei Antworten
// dieser Interaktion (V, S, C und M wie in ende-lebend.yaml) und die Version 2; stellt er keine
// mehr, endet die Session regulär: Die Aufzeichnung trägt weder die Fehlerantwort noch ein Ende
// und die Version 1, der Lauf endet mit Exit-Code 0 und ohne Meldung.
func TestE2ERecordEndeLeerlauf(t *testing.T) {
	ruhend := func(t *testing.T, listen string) (net.Conn, *bufio.Reader) {
		t.Helper()
		c, r := rohVerbindung(t, listen)
		anfrage(t, c, "SET idle_session_timeout = '300ms';")
		bisTyp(t, r, 'Z')
		time.Sleep(ruheBisZumEnde)
		return c, r
	}
	t.Run("mit Anfrage danach", func(t *testing.T) {
		output := filepath.Join(t.TempDir(), "rec.yaml")
		rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
		c, r := ruhend(t, rec.listen)
		anfrage(t, c, "SELECT 1;")
		ns := bisVerbindungsende(t, r)
		if typen(ns) != "E" || felderDerFehlerantwort(ns[0].rumpf)["C"] != "57P05" {
			t.Fatalf("erwartet die Fehlerantwort 57P05 und das Ende, erhalten %q", typen(ns))
		}
		rec.stop(t, 0)
		keineMeldung(t, "Ende im Leerlauf", rec.stderr.String())
		text := lies(t, output)
		if !strings.Contains(text, "version: 2\n") || strings.Count(text, "type: connection_closed") != 1 || !strings.Contains(text, "sql: SELECT 1;") {
			t.Fatalf("Aufzeichnung:\n%s", text)
		}
		gleicheForm(t, text, endeLebendYAML)
	})
	t.Run("ohne Anfrage danach", func(t *testing.T) {
		output := filepath.Join(t.TempDir(), "rec.yaml")
		rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), output)
		c, _ := ruhend(t, rec.listen)
		schreibe(t, c, []byte{'X', 0, 0, 0, 4})
		_ = c.Close()
		time.Sleep(300 * time.Millisecond)
		rec.stop(t, 0)
		keineMeldung(t, "Leerlauf ohne Anfrage", rec.stderr.String())
		text := lies(t, output)
		if !strings.Contains(text, "version: 1\n") || !strings.Contains(text, "SET idle_session_timeout") || strings.Contains(text, "connection_closed") || strings.Contains(text, "57P05") {
			t.Fatalf("Aufzeichnung:\n%s", text)
		}
	})
}

// kappProxy leitet Verbindungen an ziel weiter und schließt beide Seiten, sobald die Anfrage
// durchgeleitet ist, deren Text marke enthält. Das Flag gekappt verwirft Bytes, die das Ziel
// zwischen der Weiterleitung und dem Schließen sendet; es schließt ein Zeitfenster und hat
// keinen Test, der seinen Wegfall fängt. Er liefert die Adresse, auf der er lauscht.
func kappProxy(t *testing.T, ziel, marke string) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = l.Close() })
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			go kappVerbindung(c, ziel, marke)
		}
	}()
	return l.Addr().String()
}

// kappVerbindung leitet die Nachrichten von c an ziel weiter; die erste Anfrage mit marke ist
// die letzte, danach schließt es beide Verbindungen.
func kappVerbindung(c net.Conn, ziel, marke string) {
	defer c.Close()
	u, err := net.Dial("tcp", ziel)
	if err != nil {
		return
	}
	defer u.Close()
	var gekappt atomic.Bool
	go func() {
		puffer := make([]byte, 32<<10)
		for {
			n, err := u.Read(puffer)
			if n > 0 && !gekappt.Load() {
				if _, werr := c.Write(puffer[:n]); werr != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()
	r := bufio.NewReader(c)
	kopf := make([]byte, 4)
	if _, err := io.ReadFull(r, kopf); err != nil {
		return
	}
	start := make([]byte, binary.BigEndian.Uint32(kopf)-4)
	if _, err := io.ReadFull(r, start); err != nil {
		return
	}
	if _, err := u.Write(append(kopf, start...)); err != nil {
		return
	}
	for {
		kopf := make([]byte, 5)
		if _, err := io.ReadFull(r, kopf); err != nil {
			return
		}
		rumpf := make([]byte, binary.BigEndian.Uint32(kopf[1:])-4)
		if _, err := io.ReadFull(r, rumpf); err != nil {
			return
		}
		if kopf[0] == 'Q' && bytes.Contains(rumpf, []byte(marke)) {
			gekappt.Store(true)
		}
		if _, err := u.Write(append(kopf, rumpf...)); err != nil || gekappt.Load() {
			return
		}
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-02/Negative — Abnahmeszenario 18, Netzwerkabbruch: Schließt
// ein Proxy zwischen Recorder und PostgreSQL 17 beide Verbindungen nach einer weitergeleiteten
// Anfrage, ohne dass eine Antwort zurückgeht, erhält der Client keine Datenzeile und kein
// ReadyForQuery, der Lauf meldet PGR-E4003 und endet mit Exit-Code 4, und die Aufzeichnung
// trägt die vorherige Interaktion, kein connection_closed und die Version 1.
func TestE2ERecordNetzwerkabbruch(t *testing.T) {
	output := filepath.Join(t.TempDir(), "rec.yaml")
	proxy := kappProxy(t, os.Getenv("PGR_UPSTREAM"), "kappen")
	rec := startRecorder(t, proxy, output)
	c, r := rohVerbindung(t, rec.listen)
	anfrage(t, c, "SELECT 1;")
	bisTyp(t, r, 'Z')
	anfrage(t, c, "SELECT 'kappen';")
	ns := bisVerbindungsende(t, r)
	if strings.ContainsAny(typen(ns), "DZTC") {
		t.Fatalf("Antworten des Servers trotz Abbruch: %q", typen(ns))
	}
	rec.stop(t, 4)
	if !strings.Contains(rec.stderr.String(), "PGR-E4003") {
		t.Fatalf("PGR-E4003 fehlt in der Ausgabe:\n%s", rec.stderr.String())
	}
	text := lies(t, output)
	if !strings.Contains(text, "version: 1\n") || !strings.Contains(text, "sql: SELECT 1;") || strings.Contains(text, "kappen") || strings.Contains(text, "connection_closed") {
		t.Fatalf("Aufzeichnung nach dem Abbruch:\n%s", text)
	}
	laedt(t, output)
}
