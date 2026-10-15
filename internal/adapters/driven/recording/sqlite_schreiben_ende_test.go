package recording_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// feldNeutral ist ein Wert des Feldes version im Modell, den kein Fall der Wahl der Version nennt:
// Nur die Fälle „Feld N“ legen die Version über das Feld fest.
const feldNeutral = 5

// versionVon liefert den Wert von meta.version der Datei path.
func versionVon(t *testing.T, path string) string {
	t.Helper()
	z := zeilen(t, roh(t, path), `SELECT value FROM meta WHERE key = 'version'`)
	if len(z) != 1 {
		t.Fatalf("meta.version: %q", z)
	}
	return z[0]
}

// zaehle liefert die Zahl der Zeilen der Tabelle tabelle in der Datei path.
func zaehle(t *testing.T, path, tabelle string) string {
	t.Helper()
	return zeilen(t, roh(t, path), `SELECT count(*) FROM `+tabelle)[0]
}

// ixOhne ist eine Interaktion ohne Ende mit der Nummer nummer.
func ixOhne(nummer int) model.Interaction { return sEinfach(nummer, sCC(), sRFQ()) }

// ixEnde ist eine Interaktion mit Ende mit der Nummer nummer.
func ixEnde(nummer int) model.Interaction { return sEinfach(nummer, sFatal(), sEnde()) }

// ausloeserAufMeta legt in der Datei path einen Auslöser an, der jede Änderung einer
// Zeile von meta mit einem Fehler abbricht.
func ausloeserAufMeta(t *testing.T, path string) {
	t.Helper()
	rohSchreibe(t, path, `CREATE TRIGGER meta_gesperrt BEFORE UPDATE ON "meta" BEGIN SELECT RAISE(ABORT, 'meta gesperrt'); END;`)
}

// Abdeckung: LH-FA-25/Happy, LH-FA-22/Happy — der Schreiber von sqlite wählt die Version nach dem Inhalt: 2 mit einem
// Ende (in responses einer einfachen Interaktion, in server einer Gruppe, in der ersten, in einer späteren
// und in einer nicht letzten Session), sonst 1 (auch ohne Sessions und nur mit empty_sessions); das Feld
// version des Modells bestimmt sie nicht (1 mit Ende schreibt 2, 2 ohne Ende schreibt 1, ebenso 0, 3 und ein
// negativer Wert); PRAGMA user_version bleibt 0.
func TestSQLiteVersionNachInhalt(t *testing.T) {
	leer := sRec(feldNeutral, sSession(1, "u"))
	leer.EmptySessions = true
	cases := []struct {
		name string
		rec  model.Recording
		want string
	}{
		{"Feld 1 mit Ende", sRec(1, sSession(1, "feld1", ixEnde(1))), "2"},
		{"Feld 1 ohne Ende", sRec(1, sSession(1, "feld1o", ixOhne(1))), "1"},
		{"Feld 2 mit Ende", sRec(2, sSession(1, "feld2", ixEnde(1))), "2"},
		{"Feld 2 ohne Ende", sRec(2, sSession(1, "feld2o", ixOhne(1))), "1"},
		{"Feld 0 mit Ende", sRec(0, sSession(1, "feld0", ixEnde(1))), "2"},
		{"Feld 0 ohne Ende", sRec(0, sSession(1, "feld0o", ixOhne(1))), "1"},
		{"Feld 3 mit Ende", sRec(3, sSession(1, "feld3", ixEnde(1))), "2"},
		{"Feld 3 ohne Ende", sRec(3, sSession(1, "feld3o", ixOhne(1))), "1"},
		{"Feld -1 mit Ende", sRec(-1, sSession(1, "feldm1", ixEnde(1))), "2"},
		{"Feld -1 ohne Ende", sRec(-1, sSession(1, "feldm1o", ixOhne(1))), "1"},
		{"ohne Sessions", sRec(feldNeutral), "1"},
		{"nur empty_sessions", leer, "1"},
		{"Ende nur in einer späteren Session", sRec(feldNeutral, sSession(1, "u", ixOhne(1)), sSession(2, "u", ixEnde(1))), "2"},
		{"Ende nur in einer nicht letzten Session", sRec(feldNeutral, sSession(1, "u", ixEnde(1)), sSession(2, "u", ixOhne(1))), "2"},
		{"Ende in einer Extended-Gruppe, in keiner einfachen Interaktion", sRec(feldNeutral, sSession(1, "u", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde())))), "2"},
		{"Ende in einer einfachen Interaktion, in keiner Gruppe", sRec(feldNeutral, sSession(1, "u", ixOhne(1), ixEnde(2))), "2"},
		{"Flush-Gruppe mit Ende, zweite Session ohne", sRec(feldNeutral, sSession(1, "u", sExtended(1, sGruppe(model.ClientFlush, sFatal(), sEnde()))), sSession(2, "u", ixOhne(1))), "2"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, path := neueDatei(t, c.rec)
			if got := versionVon(t, path); got != c.want {
				t.Fatalf("meta.version %q, erwartet %q", got, c.want)
			}
			if got := zeilen(t, roh(t, path), `PRAGMA user_version`); !reflect.DeepEqual(got, []string{"0"}) {
				t.Fatalf("user_version %q", got)
			}
		})
	}
}

// Abdeckung: LH-FA-25/Happy, LH-FA-22/Happy — hat der erste Schreibvorgang kein Ende, bleibt die Version 1, und der
// Schreibvorgang mit der ersten Session mit Ende setzt sie auf 2; die Transaktion einer zweiten Session mit
// Ende und die einer Session ohne Ende nach der Version 2 ändern meta nicht (ein Auslöser auf meta bricht jede
// Änderung ab, der Schreibvorgang gelingt); mehrere Sessions mit Ende in einem Schreibvorgang setzen die
// Version einmal; die Datei lädt danach als dieselbe Aufzeichnung mit Version 2.
func TestSQLiteVersionWerSetzt(t *testing.T) {
	t.Run("spätere Session setzt, frühere nicht", func(t *testing.T) {
		s, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixOhne(1))))
		if got := versionVon(t, path); got != "1" {
			t.Fatalf("nach dem ersten Schreibvorgang %q", got)
		}
		rec := sRec(feldNeutral, sSession(1, "u", ixOhne(1)), sSession(2, "u", ixEnde(1)))
		if err := s.Write(context.Background(), path, rec); err != nil {
			t.Fatal(err)
		}
		if got := versionVon(t, path); got != "2" {
			t.Fatalf("nach der Session mit Ende %q", got)
		}
	})
	t.Run("zweite Session mit Ende im selben Schreibvorgang", func(t *testing.T) {
		_, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixEnde(1)), sSession(2, "u", ixEnde(1)), sSession(3, "u", ixOhne(1))))
		if got := versionVon(t, path); got != "2" {
			t.Fatalf("meta.version %q", got)
		}
		if got := sessionIDs(t, path); got != "1,2,3" {
			t.Fatalf("Sessions %q", got)
		}
	})
	t.Run("zweite Session mit Ende ändert meta nicht", func(t *testing.T) {
		s, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixEnde(1))))
		ausloeserAufMeta(t, path)
		if err := s.Write(context.Background(), path, sRec(feldNeutral, sSession(1, "u", ixEnde(1)), sSession(2, "u", ixEnde(1)))); err != nil {
			t.Fatalf("Schreibvorgang: %v", err)
		}
		if got, want := sessionIDs(t, path), "1,2"; got != want {
			t.Fatalf("Sessions %q, erwartet %q", got, want)
		}
		if got := versionVon(t, path); got != "2" {
			t.Fatalf("meta.version %q", got)
		}
	})
	t.Run("Session ohne Ende nach der Version 2 ändert meta nicht", func(t *testing.T) {
		s, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixEnde(1))))
		ausloeserAufMeta(t, path)
		if err := s.Write(context.Background(), path, sRec(feldNeutral, sSession(1, "u", ixEnde(1)), sSession(2, "u", ixOhne(1)))); err != nil {
			t.Fatalf("Schreibvorgang: %v", err)
		}
		if got, want := sessionIDs(t, path), "1,2"; got != want {
			t.Fatalf("Sessions %q, erwartet %q", got, want)
		}
		if got := versionVon(t, path); got != "2" {
			t.Fatalf("meta.version %q", got)
		}
	})
	t.Run("Session ohne Ende ändert meta nicht", func(t *testing.T) {
		s, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixOhne(1))))
		ausloeserAufMeta(t, path)
		if err := s.Write(context.Background(), path, sRec(feldNeutral, sSession(1, "u", ixOhne(1)), sSession(2, "u", ixOhne(1)))); err != nil {
			t.Fatalf("Schreibvorgang: %v", err)
		}
		if got := versionVon(t, path); got != "1" {
			t.Fatalf("meta.version %q", got)
		}
	})
	t.Run("Laden nach späteren Schreibvorgängen", func(t *testing.T) {
		s, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixEnde(1))))
		rec := sRec(feldNeutral, sSession(1, "u", ixEnde(1)), sSession(2, "u", ixOhne(1)), sSession(3, "u", ixEnde(1)))
		if err := s.Write(context.Background(), path, rec); err != nil {
			t.Fatal(err)
		}
		want := rec
		want.Version = model.VersionMitEnde
		got, err := laden(path)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("%v\n got %#v\nwant %#v", err, got, want)
		}
	})
}

// Abdeckung: LH-FA-25/Negative, LH-FA-22/Negative — steht in der Datei die Zeile version von meta nicht
// oder nicht mit dem Wert 1, ist das Setzen der Version in der Transaktion der ersten Session mit Ende
// PGR-E3001 und die Session nicht ergänzt (keine Zeile der Session, meta unverändert); trägt meta
// keine Zeile version, aber eine andere mit dem Wert 1, ändert die Transaktion diese nicht.
func TestSQLiteVersionZeileUnpassend(t *testing.T) {
	cases := []struct {
		name     string
		eingriff string
		meta     string
	}{
		{"Zeile fehlt", `DELETE FROM meta WHERE key = 'version'`, "format|pgwire-recorder"},
		{"Wert 2", `UPDATE meta SET value = '2' WHERE key = 'version'`, "format|pgwire-recorder version|2"},
		{"Zeile fehlt, andere Zeile mit Wert 1", `UPDATE meta SET key = 'anders' WHERE key = 'version'`, "anders|1 format|pgwire-recorder"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixOhne(1))))
			rohSchreibe(t, path, c.eingriff)
			err := s.Write(context.Background(), path, sRec(feldNeutral, sSession(1, "u", ixOhne(1)), sSession(2, "u", ixEnde(1))))
			if code(err) != model.CodeRecordingIO {
				t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingIO, err)
			}
			if got := sessionIDs(t, path); got != "1" {
				t.Fatalf("Sessions %q", got)
			}
			for tabelle, want := range map[string]string{"interaction": "1", "message": "2"} {
				if got := zaehle(t, path, tabelle); got != want {
					t.Errorf("%s trägt %s Zeilen, erwartet %s", tabelle, got, want)
				}
			}
			if got := strings.Join(zeilen(t, roh(t, path), `SELECT key || '|' || value FROM meta ORDER BY key`), " "); got != c.meta {
				t.Fatalf("meta %q, erwartet %q", got, c.meta)
			}
		})
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-22/Negative — scheitert die Transaktion der ersten Session mit Ende beim Setzen der
// Version (ein Auslöser auf meta bricht ab, nachdem die Zeilen der Session eingefügt sind), ist das PGR-E3001 mit
// der Ursache im Text, die Datei trägt nur die zuvor ergänzten Sessions mit allen ihren Zeilen, die Version
// bleibt 1, und der nächste Schreibvorgang (nach dem Entfernen des Auslösers) ergänzt die Session und setzt 2;
// scheitert im ersten Schreibvorgang das Einfügen der Session (doppelte sequence), steht die Datei mit der
// Version 1 und ohne Session da, und der nächste Schreibvorgang ergänzt und setzt 2.
func TestSQLiteVersionFehlerInDerTransaktion(t *testing.T) {
	t.Run("Auslöser beim Setzen", func(t *testing.T) {
		s, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixOhne(1))))
		ausloeserAufMeta(t, path)
		rec := sRec(feldNeutral, sSession(1, "u", ixOhne(1)), sSession(2, "u", ixEnde(1)))
		err := s.Write(context.Background(), path, rec)
		ms := model.Meldungen(err)
		if len(ms) != 1 || ms[0].Code != model.CodeRecordingIO || !strings.Contains(ms[0].Text, "meta gesperrt") {
			t.Fatalf("Meldungen %+v", ms)
		}
		if got := sessionIDs(t, path); got != "1" {
			t.Fatalf("Sessions %q", got)
		}
		for tabelle, want := range map[string]string{"session": "1", "interaction": "1", "message": "2"} {
			if got := zaehle(t, path, tabelle); got != want {
				t.Errorf("%s trägt %s Zeilen, erwartet %s", tabelle, got, want)
			}
		}
		if got := versionVon(t, path); got != "1" {
			t.Fatalf("meta.version %q", got)
		}
		rohSchreibe(t, path, `DROP TRIGGER meta_gesperrt`)
		if err := s.Write(context.Background(), path, rec); err != nil {
			t.Fatalf("nächster Schreibvorgang: %v", err)
		}
		if got := sessionIDs(t, path); got != "1,2" {
			t.Fatalf("Sessions %q", got)
		}
		if got := versionVon(t, path); got != "2" {
			t.Fatalf("meta.version %q", got)
		}
	})
	t.Run("doppelte sequence im ersten Schreibvorgang", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "rec.sqlite")
		s := &recording.SQLite{}
		rec := sRec(feldNeutral, sSession(1, "u", ixEnde(1), ixOhne(1)))
		if code(s.Write(context.Background(), path, rec)) != model.CodeRecordingIO {
			t.Fatal("erwartet PGR-E3001")
		}
		if got := versionVon(t, path); got != "1" {
			t.Fatalf("meta.version %q", got)
		}
		for _, tabelle := range []string{"session", "interaction", "message"} {
			if got := zaehle(t, path, tabelle); got != "0" {
				t.Errorf("%s trägt %s Zeilen", tabelle, got)
			}
		}
		rec.Sessions[0].Interactions[1].Sequence = 2
		if err := s.Write(context.Background(), path, rec); err != nil {
			t.Fatalf("nächster Schreibvorgang: %v", err)
		}
		if got := versionVon(t, path); got != "2" {
			t.Fatalf("meta.version %q", got)
		}
	})
}

// Abdeckung: LH-FA-25/Negative, LH-FA-22/Negative — scheitert das Schließen der Datenbank nach der Transaktion der
// ersten Session mit Ende, ist das PGR-E3001, die Version 2 steht, und der nächste Schreibvorgang ändert sie
// nicht (eine weitere Session mit Ende gelingt trotz eines Auslösers auf meta) und ergänzt die Session nicht
// doppelt.
func TestSQLiteVersionSchliessenScheitert(t *testing.T) {
	ctx := context.Background()
	s, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixOhne(1))))
	fehler := errors.New("Schließen eingespielt gescheitert")
	rec := sRec(feldNeutral, sSession(1, "u", ixOhne(1)), sSession(2, "u", ixEnde(1)))
	err := recording.SQLiteSchreibeMit(ctx, s, path, rec, recording.SQLiteEingriffe{SchliesseDB: func(db *sql.DB) error {
		return errors.Join(db.Close(), fehler)
	}})
	ms := model.Meldungen(err)
	if len(ms) != 1 || ms[0].Code != model.CodeRecordingIO || !strings.Contains(ms[0].Text, fehler.Error()) {
		t.Fatalf("Meldungen %+v", ms)
	}
	if got := versionVon(t, path); got != "2" {
		t.Fatalf("meta.version %q", got)
	}
	ausloeserAufMeta(t, path)
	rec.Sessions = append(rec.Sessions, sSession(3, "u", ixEnde(1)))
	if err := s.Write(ctx, path, rec); err != nil {
		t.Fatalf("nächster Schreibvorgang: %v", err)
	}
	if got := sessionIDs(t, path); got != "1,2,3" {
		t.Fatalf("Sessions %q", got)
	}
	if got := versionVon(t, path); got != "2" {
		t.Fatalf("meta.version %q", got)
	}
}

// Abdeckung: LH-FA-25/Boundary, LH-FA-22/Boundary — ein abgebrochener Kontext unterbricht die Transaktion mit dem
// Ende nicht: Der Schreibvorgang legt die Datei an, ergänzt die Session mit Ende vollständig und setzt die Version 2.
func TestSQLiteVersionAbbruchUnterbrichtNicht(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	path := filepath.Join(t.TempDir(), "rec.sqlite")
	if err := (&recording.SQLite{}).Write(ctx, path, sRec(feldNeutral, sSession(1, "u", ixOhne(1)), sSession(2, "u", ixEnde(1)))); err != nil {
		t.Fatal(err)
	}
	if got := sessionIDs(t, path); got != "1,2" {
		t.Fatalf("Sessions %q", got)
	}
	if got := versionVon(t, path); got != "2" {
		t.Fatalf("meta.version %q", got)
	}
	if got := zaehle(t, path, "message"); got != "4" {
		t.Fatalf("message trägt %s Zeilen", got)
	}
}

// Abdeckung: LH-FA-25/Negative, LH-FA-22/Negative — eine Session mit Ende und falscher Nummer im ersten Schreibvorgang ist
// PGR-E1000, und die Datei entsteht nicht.
func TestSQLiteEndeMitFalscherNummer(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rec.sqlite")
	err := (&recording.SQLite{}).Write(context.Background(), path, sRec(feldNeutral, sSession(2, "u", ixEnde(1))))
	if code(err) != model.CodeInternal {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeInternal, err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("Zielpfad: %v", err)
	}
}

// Abdeckung: LH-FA-25/Happy, LH-FA-22/Happy — das Ende steht als Zeile von message mit kind connection_closed, direction
// server, der sequence der Interaktion und fields genau {} (Text mit zwei Zeichen): in einer einfachen
// Interaktion mit group_no NULL und position nach der Fehlerantwort, in einer Sync-Gruppe und in einer
// Flush-Gruppe mit der Nummer der Gruppe und position über beide Richtungen, unter den Server-Nachrichten der
// Gruppe; offset_ms steht an der Interaktion mit Ende; eine Aufzeichnung ohne Ende schreibt dieselben
// Zeilen wie bisher mit den meta-Zeilen format und version 1.
func TestSQLiteEndeDarstellung(t *testing.T) {
	off := int64(7)
	mitOffset := ixEnde(1)
	mitOffset.OffsetMS = &off
	cases := []struct {
		name string
		rec  model.Recording
		want []string
	}{
		{"einfach", sRec(feldNeutral, sSession(1, "u", ixOhne(1), ixEnde(2))),
			[]string{"2|2|server|NULL|connection_closed|text|2|{}"}},
		{"Sync-Gruppe", sRec(feldNeutral, sSession(1, "u", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde())))),
			[]string{"1|3|server|1|connection_closed|text|2|{}"}},
		{"Flush-Gruppe", sRec(feldNeutral, sSession(1, "u", sExtended(1, sGruppe(model.ClientFlush, sCC()), sGruppe(model.ClientFlush, sFatal(), sEnde())))),
			[]string{"1|5|server|2|connection_closed|text|2|{}"}},
		{"mit offset_ms", sRec(feldNeutral, sSession(1, "u", mitOffset)),
			[]string{"1|2|server|NULL|connection_closed|text|2|{}"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, path := neueDatei(t, c.rec)
			got := zeilen(t, roh(t, path), `SELECT sequence, position, direction, group_no, kind, typeof(fields), length(fields), fields FROM message WHERE kind = 'connection_closed'`)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Zeilen des Endes %q, erwartet %q", got, c.want)
			}
		})
	}
	t.Run("offset_ms an der Interaktion mit Ende", func(t *testing.T) {
		_, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", mitOffset)))
		if got := zeilen(t, roh(t, path), `SELECT offset_ms FROM interaction`); !reflect.DeepEqual(got, []string{"7"}) {
			t.Fatalf("offset_ms %q", got)
		}
	})
	t.Run("ohne Ende", func(t *testing.T) {
		_, path := neueDatei(t, sRec(feldNeutral, sSession(1, "u", ixOhne(1), sExtended(2, sGruppe(model.ClientSync, sCC(), sRFQ())))))
		db := roh(t, path)
		want := map[string][]string{
			`SELECT key || '|' || value FROM meta ORDER BY key`:                 {"format|pgwire-recorder", "version|1"},
			`SELECT sequence || '|' || type FROM interaction ORDER BY sequence`: {"1|query", "2|extended"},
			`SELECT sequence || '|' || position || '|' || direction || '|' || ifnull(group_no, 'NULL') || '|' || kind FROM message ORDER BY sequence, position`: {
				"1|1|server|NULL|command_complete", "1|2|server|NULL|ready_for_query",
				"2|1|client|1|sync", "2|2|server|1|command_complete", "2|3|server|1|ready_for_query"},
		}
		for anfrage, w := range want {
			if got := zeilen(t, db, anfrage); !reflect.DeepEqual(got, w) {
				t.Errorf("%s:\n got %q\nwant %q", anfrage, got, w)
			}
		}
	})
	t.Run("Feld neben type steht in fields", func(t *testing.T) {
		rec := sRec(feldNeutral, sSession(1, "u", sEinfach(1, sFatal(), model.Response{Type: model.ResponseConnectionClosed, Tag: "x"})))
		_, path := neueDatei(t, rec)
		if got := zeilen(t, roh(t, path), `SELECT fields FROM message WHERE kind = 'connection_closed'`); !reflect.DeepEqual(got, []string{`{"tag":"x"}`}) {
			t.Fatalf("fields %q", got)
		}
		if _, err := laden(path); code(err) != model.CodeRecordingBroken {
			t.Fatalf("erwartet %s beim Laden, erhalten %v", model.CodeRecordingBroken, err)
		}
	})
}

// Abdeckung: LH-FA-25/Negative, LH-FA-22/Negative — der Schreiber von sqlite prüft die Form des Endes nicht: Eine Aufzeichnung mit
// dem Ende an einer falschen Stelle (Anfang, Mitte, nicht letzte Interaktion, nicht letzte Gruppe, mitten in der
// letzten Gruppe), zweimal in einer Interaktion, zweimal in zwei Interaktionen einer Session, mit
// ready_for_query vor oder nach dem Ende oder mit einem Feld neben type schreibt ohne Fehler (Version 2, ein
// Aufruf, err gleich nil), und das Laden lehnt sie mit PGR-E3003 ab; das Ende an der richtigen Stelle
// (direkt nach der Fehlerantwort als letzte Nachricht, in der letzten Gruppe, in der letzten Interaktion
// und in einer nicht letzten Session) lädt.
func TestSQLiteSchreibenNimmtFormverletzungAn(t *testing.T) {
	cases := []struct {
		name string
		rec  model.Recording
		lade bool
	}{
		{"direkt nach der Fehlerantwort, letzte Nachricht", sRec(feldNeutral, sSession(1, "stelle-richtig", sEinfach(1, sFatal(), sEnde()))), true},
		{"letzte Gruppe", sRec(feldNeutral, sSession(1, "stelle-gruppe", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde())))), true},
		{"letzte Interaktion einer Session", sRec(feldNeutral, sSession(1, "stelle-letzte-interaktion", ixOhne(1), ixEnde(2))), true},
		{"nicht letzte Session", sRec(feldNeutral, sSession(1, "stelle-nicht-letzte-session", ixEnde(1)), sSession(2, "u", ixOhne(1))), true},
		{"Anfang der Folge", sRec(feldNeutral, sSession(1, "stelle-anfang", sEinfach(1, sEnde(), sFatal(), sCC()))), false},
		{"Mitte der Folge", sRec(feldNeutral, sSession(1, "stelle-mitte", sEinfach(1, sFatal(), sEnde(), sCC()))), false},
		{"nicht letzte Interaktion", sRec(feldNeutral, sSession(1, "stelle-nicht-letzte-interaktion", ixEnde(1), ixOhne(2))), false},
		{"nicht letzte Gruppe", sRec(feldNeutral, sSession(1, "stelle-nicht-letzte-gruppe", sExtended(1, sGruppe(model.ClientFlush, sFatal(), sEnde()), sGruppe(model.ClientSync, sCC(), sRFQ())))), false},
		{"mitten in der letzten Gruppe", sRec(feldNeutral, sSession(1, "stelle-mitte-gruppe", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde(), sCC())))), false},
		{"zwei Enden in einer Interaktion", sRec(feldNeutral, sSession(1, "stelle-zwei-in-einer", sEinfach(1, sFatal(), sEnde(), sFatal(), sEnde()))), false},
		{"zwei Enden in zwei Interaktionen einer Session", sRec(feldNeutral, sSession(1, "stelle-zwei-in-zwei", ixEnde(1), ixEnde(2))), false},
		{"ready_for_query vor dem Ende", sRec(feldNeutral, sSession(1, "stelle-rfq-vor", sEinfach(1, sCC(), sRFQ(), sFatal(), sEnde()))), false},
		{"ready_for_query nach dem Ende", sRec(feldNeutral, sSession(1, "stelle-rfq-nach", sEinfach(1, sFatal(), sEnde(), sRFQ()))), false},
		{"Feld neben type", sRec(feldNeutral, sSession(1, "stelle-feld", sEinfach(1, sFatal(), model.Response{Type: model.ResponseConnectionClosed, Tag: "x"}))), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "rec.sqlite")
			err := (&recording.SQLite{}).Write(context.Background(), path, c.rec)
			if err != nil {
				t.Fatalf("Schreiben: %v", err)
			}
			if got := versionVon(t, path); got != "2" {
				t.Fatalf("meta.version %q, erwartet 2", got)
			}
			_, err = laden(path)
			if c.lade && err != nil {
				t.Fatalf("Laden: %v", err)
			}
			if !c.lade && code(err) != model.CodeRecordingBroken {
				t.Fatalf("erwartet %s beim Laden, erhalten %v", model.CodeRecordingBroken, err)
			}
		})
	}
}

// Abdeckung: LH-FA-25/Happy, LH-FA-22/Happy, LH-QA-06/Messung — dieselbe Aufzeichnung mit Ende, als yaml und als sqlite geschrieben, liest gleich
// (auch in version) und gleich der geschriebenen, und von der geladenen Aufzeichnung aus ergibt Schreiben und
// erneutes Laden dieselbe: einfach, Extended mit Sync-Gruppe, mit Flush-Gruppe, mit offset_ms an der Interaktion mit
// Ende, mit empty_sessions; eine Datei der Version 2 ohne Ende liest mit Version 2 und ergibt beim Schreiben die
// Version 1, sonst gleich; ein späterer Schreibvorgang an eine Datei der Version 2 ergänzt Sessions ohne Änderung der
// Version.
func TestSQLiteRoundtripMitEnde(t *testing.T) {
	off := int64(7)
	mitOffset := ixEnde(1)
	mitOffset.OffsetMS = &off
	leer := sRec(feldNeutral, sSession(1, "rt-leer", ixEnde(1)), sSession(2, "rt-leer"))
	leer.EmptySessions = true
	cases := []struct {
		name string
		rec  model.Recording
	}{
		{"einfach", sRec(feldNeutral, sSession(1, "rt-einfach", ixOhne(1), ixEnde(2)))},
		{"Sync-Gruppe", sRec(feldNeutral, sSession(1, "rt-sync", sExtended(1, sGruppe(model.ClientSync, sFatal(), sEnde()))))},
		{"Flush-Gruppe", sRec(feldNeutral, sSession(1, "rt-flush", sExtended(1, sGruppe(model.ClientFlush, sCC()), sGruppe(model.ClientFlush, sFatal(), sEnde()))))},
		{"offset_ms", sRec(feldNeutral, sSession(1, "rt-offset", mitOffset))},
		{"empty_sessions", leer},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			yamlPfad := filepath.Join(t.TempDir(), "rec.yaml")
			if err := (recording.YAML{}).Write(context.Background(), yamlPfad, c.rec); err != nil {
				t.Fatal(err)
			}
			ausYAML, err := laden(yamlPfad)
			if err != nil {
				t.Fatal(err)
			}
			_, pfad := neueDatei(t, c.rec)
			ausSQLite, err := laden(pfad)
			if err != nil {
				t.Fatal(err)
			}
			want := c.rec
			want.Version = model.VersionMitEnde
			if !reflect.DeepEqual(ausSQLite, ausYAML) || !reflect.DeepEqual(ausSQLite, want) {
				t.Fatalf("Formate weichen ab:\nsqlite %#v\n  yaml %#v\n  want %#v", ausSQLite, ausYAML, want)
			}
			_, nochmal := neueDatei(t, ausSQLite)
			again, err := laden(nochmal)
			if err != nil || !reflect.DeepEqual(again, ausSQLite) {
				t.Fatalf("Laden, Schreiben, Laden: %v\n got %#v\nwant %#v", err, again, ausSQLite)
			}
		})
	}
	t.Run("Version 2 ohne Ende", func(t *testing.T) {
		_, pfad := neueDatei(t, sRec(feldNeutral, sSession(1, "rt-v2", ixOhne(1))))
		rohSchreibe(t, pfad, `UPDATE meta SET value = '2' WHERE key = 'version'`)
		geladen, err := laden(pfad)
		if err != nil || geladen.Version != 2 {
			t.Fatalf("%v, Version %d", err, geladen.Version)
		}
		_, neu := neueDatei(t, geladen)
		if got := versionVon(t, neu); got != "1" {
			t.Fatalf("meta.version %q", got)
		}
		again, err := laden(neu)
		geladen.Version = 1
		if err != nil || !reflect.DeepEqual(again, geladen) {
			t.Fatalf("%v\n got %#v\nwant %#v", err, again, geladen)
		}
	})
}
