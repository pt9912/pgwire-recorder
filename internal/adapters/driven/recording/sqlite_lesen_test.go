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
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// laden liest path mit dem Leser, der das Format selbst erkennt.
func laden(path string) (model.Recording, error) {
	return recording.Leser{}.Load(context.Background(), path)
}

// ladeOK liest path und beendet den Test bei einem Fehler.
func ladeOK(t *testing.T, path string) model.Recording {
	t.Helper()
	rec, err := laden(path)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return rec
}

// lesefixture ist die Aufzeichnung, an der die Tests des Lesens ändern: Session 1 mit einer
// einfachen und einer Extended-Interaktion (drei Gruppen), Session 2 mit einer einfachen.
func lesefixture() model.Recording {
	rec := extendedBeispiel()
	sieben := int64(7)
	rec.Sessions[0].Interactions[0].OffsetMS = &sieben
	rec.Sessions = append(rec.Sessions, sessionsMit(2).Sessions[1])
	return rec
}

// sqliteDatei schreibt rec nach rec.sqlite in einem neuen Verzeichnis.
func sqliteDatei(t *testing.T, rec model.Recording) string {
	t.Helper()
	_, path := neueDatei(t, rec)
	return path
}

// kopiere kopiert die Datei path nach rec.sqlite in einem neuen Verzeichnis.
func kopiere(t *testing.T, path string) string {
	t.Helper()
	neu := filepath.Join(t.TempDir(), "rec.sqlite")
	if err := os.WriteFile(neu, mustRead(t, path), 0o644); err != nil {
		t.Fatal(err)
	}
	return neu
}

// lockere ersetzt die Tabellen der Datei path durch Kopien ohne Schlüssel und ohne
// Bedingungen, damit ein Test eine Datei schreiben kann, die das Schema der Aufzeichnung
// verletzt (eine fremde Datei mit denselben Tabellen).
func lockere(t *testing.T, path string) {
	t.Helper()
	var b strings.Builder
	for _, tabelle := range []string{"message", "interaction", "session", "meta"} {
		b.WriteString(`CREATE TABLE "` + tabelle + `_neu" AS SELECT * FROM "` + tabelle + `"; DROP TABLE "` + tabelle + `"; ALTER TABLE "` + tabelle + `_neu" RENAME TO "` + tabelle + `";`)
	}
	rohSchreibe(t, path, b.String())
}

// codeVon liefert den Meldungscode von err oder "" ohne Fehler.
func codeVon(err error) string {
	if err == nil {
		return ""
	}
	return code(err)
}

// Abdeckung: LH-FA-22/Happy — dieselbe Aufzeichnung aus yaml und aus sqlite liest der Leser als
// dieselbe Aufzeichnung: einfache und Extended-Interaktionen (die Nachrichtenarten von beispiel
// und extendedBeispiel), mit und ohne Zeitangaben, mit leeren Sessions, mit und ohne Startup-Parameter und ohne
// Sessions; die Aufzeichnung aus sqlite ist auch gleich der geschriebenen.
func TestLesenRoundtrip(t *testing.T) {
	leer := model.NewRecording()
	leereSessions := sessionsMit(1)
	leereSessions.EmptySessions = true
	leereSessions.Sessions = append(leereSessions.Sessions, model.Session{ID: 2})
	ohneStartup := sessionsMit(2)
	ohneStartup.Sessions[0].Startup = nil
	ohneStartup.Sessions[1].ServerParameters = map[string]string{"a": "b"}
	for name, c := range map[string]struct {
		rec      model.Recording
		sessions int
	}{
		"einfach":          {beispiel(), 1},
		"extended":         {extendedBeispiel(), 1},
		"gemischt":         {lesefixture(), 2},
		"ohne Sessions":    {leer, 0},
		"leere Session":    {leereSessions, 2},
		"ohne Startup":     {ohneStartup, 2},
		"mehr als 9 Zeile": {sessionsMit(12), 12},
	} {
		t.Run(name, func(t *testing.T) {
			ys := filepath.Join(t.TempDir(), "rec.yaml")
			if err := (recording.YAML{}).Write(context.Background(), ys, c.rec); err != nil {
				t.Fatal(err)
			}
			ausYAML := ladeOK(t, ys)
			ausSQLite := ladeOK(t, sqliteDatei(t, c.rec))
			if len(ausSQLite.Sessions) != c.sessions {
				t.Fatalf("%d Sessions, erwartet %d", len(ausSQLite.Sessions), c.sessions)
			}
			if !reflect.DeepEqual(ausSQLite, ausYAML) {
				t.Errorf("aus sqlite:\n%#v\naus yaml:\n%#v", ausSQLite, ausYAML)
			}
			if c.sessions > 0 && !reflect.DeepEqual(ausSQLite, c.rec) {
				t.Errorf("aus sqlite:\n%#v\ngeschrieben:\n%#v", ausSQLite, c.rec)
			}
		})
	}
}

// Abdeckung: LH-FA-22/Happy — Load des SQLite-Adapters liest die Datei ohne Erkennung, und der
// Leser liefert mit beiden Formaten dieselbe Aufzeichnung (Load von YAML, SQLite und Leser).
func TestLesenAdapterUndLeser(t *testing.T) {
	rec := lesefixture()
	path := sqliteDatei(t, rec)
	got, err := (&recording.SQLite{}).Load(context.Background(), path)
	if err != nil || !reflect.DeepEqual(got, rec) {
		t.Fatalf("SQLite.Load: %v\n%#v", err, got)
	}
	ys := filepath.Join(t.TempDir(), "rec.yaml")
	if err := (recording.YAML{}).Write(context.Background(), ys, rec); err != nil {
		t.Fatal(err)
	}
	gotY, err := (recording.YAML{}).Load(context.Background(), ys)
	if err != nil || !reflect.DeepEqual(gotY, got) {
		t.Fatalf("YAML.Load: %v", err)
	}
	if _, err := (recording.YAML{}).Load(context.Background(), path); code(err) != model.CodeRecordingBroken {
		t.Errorf("YAML.Load einer SQLite-Datei: %v", err)
	}
}

// Abdeckung: LH-FA-22/Happy — die Reihenfolge der Zeilen in der Datei hat keine Bedeutung: Sind die
// Zeilen von session, interaction oder message in umgekehrter Reihenfolge eingefügt (je Tabelle
// allein), liest der Leser dieselbe Aufzeichnung.
func TestLesenReihenfolgeDerZeilen(t *testing.T) {
	rec := lesefixture()
	for _, c := range []struct{ tabelle, ordnung string }{
		{"session", "id"},
		{"interaction", "session_id, sequence"},
		{"message", "session_id, sequence, position"},
	} {
		t.Run(c.tabelle, func(t *testing.T) {
			path := sqliteDatei(t, rec)
			if n := zeilenAn(t, path, `SELECT count(*) FROM `+c.tabelle); n == "0" || n == "1" {
				t.Fatalf("zu wenige Zeilen in %s: %s", c.tabelle, n)
			}
			rohSchreibe(t, path, `CREATE TABLE neu AS SELECT * FROM `+c.tabelle+` ORDER BY `+strings.ReplaceAll(c.ordnung, ", ", " DESC, ")+` DESC; DROP TABLE `+c.tabelle+`; ALTER TABLE neu RENAME TO `+c.tabelle)
			if got := zeilenAn(t, path, `SELECT `+strings.Split(c.ordnung, ",")[0]+` FROM `+c.tabelle+` LIMIT 1`); got == "1" {
				t.Fatalf("die Reihenfolge der Datei ist nicht umgekehrt: %s", got)
			}
			if got := ladeOK(t, path); !reflect.DeepEqual(got, rec) {
				t.Errorf("aus der umgestellten Datei:\n%#v", got)
			}
		})
	}
}

// zeilenAn liefert die erste Zeile der Anfrage an die Datei path.
func zeilenAn(t *testing.T, path, anfrage string) string {
	t.Helper()
	got := zeilen(t, roh(t, path), anfrage)
	if len(got) == 0 {
		t.Fatalf("%s: keine Zeile", anfrage)
	}
	return got[0]
}

// Abdeckung: LH-FA-22/Boundary — Zusätze der Datei ignoriert der Leser: eine zusätzliche Tabelle, Spalte,
// ein Index, ein Trigger, eine View und PRAGMA user_version ändern nichts.
func TestLesenIgnoriertZusaetze(t *testing.T) {
	rec := lesefixture()
	path := sqliteDatei(t, rec)
	rohSchreibe(t, path, `CREATE TABLE fremd (a TEXT, b INTEGER);
INSERT INTO fremd VALUES ('x', 1);
ALTER TABLE message ADD COLUMN extra TEXT;
ALTER TABLE session ADD COLUMN zusatz INTEGER DEFAULT 5;
UPDATE message SET extra = 'x';
CREATE INDEX fremd_idx ON message (kind);
CREATE VIEW sicht AS SELECT id FROM session;
CREATE TRIGGER fremd_trigger AFTER INSERT ON fremd BEGIN SELECT 1; END;
PRAGMA user_version = 42`)
	if got := ladeOK(t, path); !reflect.DeepEqual(got, rec) {
		t.Errorf("mit Zusätzen:\n%#v", got)
	}
	if got := zeilenAn(t, path, `PRAGMA user_version`); got != "42" {
		t.Errorf("user_version %s", got)
	}
}

// Abdeckung: LH-FA-22/Negative — fehlt eine Tabelle oder eine Spalte der Tabellen nach SPEC-043, ist
// die Datei beschädigt (PGR-E3003), ebenso eine fremde SQLite-Datei ohne diese Tabellen.
func TestLesenFehlendeTabelleOderSpalte(t *testing.T) {
	lose := sqliteDatei(t, lesefixture())
	lockere(t, lose)
	for name, ddl := range map[string]string{
		"meta":             `DROP TABLE meta`,
		"session":          `DROP TABLE session`,
		"interaction":      `DROP TABLE interaction`,
		"message":          `DROP TABLE message`,
		"meta.value":       `ALTER TABLE meta DROP COLUMN value`,
		"session.startup":  `ALTER TABLE session DROP COLUMN startup`,
		"session.server":   `ALTER TABLE session RENAME COLUMN server_parameters TO sp`,
		"interaction.type": `ALTER TABLE interaction RENAME COLUMN type TO art`,
		"interaction.sql":  `ALTER TABLE interaction DROP COLUMN sql`,
		"interaction.off":  `ALTER TABLE interaction DROP COLUMN offset_ms`,
		"message.fields":   `ALTER TABLE message RENAME COLUMN fields TO felder`,
		"message.group_no": `ALTER TABLE message DROP COLUMN group_no`,
		"message.kind":     `ALTER TABLE message DROP COLUMN kind`,
	} {
		path := kopiere(t, lose)
		rohSchreibe(t, path, ddl)
		if _, err := laden(path); code(err) != model.CodeRecordingBroken {
			t.Errorf("%s: %v", name, err)
		}
	}
	fremd := filepath.Join(t.TempDir(), "fremd.sqlite")
	if err := os.WriteFile(fremd, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	rohSchreibe(t, fremd, `CREATE TABLE a (b TEXT)`)
	if _, err := laden(fremd); code(err) != model.CodeRecordingBroken {
		t.Errorf("fremde Datei: %v", err)
	}
}

// Abdeckung: LH-FA-22/Boundary — die Datei wird nur lesend geöffnet und nie verändert: Nach dem Lesen
// sind ihre Bytes gleich, im Verzeichnis liegt nur sie, und ein Schreibversuch in der
// Lesetransaktion scheitert (attempt to write a readonly database).
func TestLesenNurLesend(t *testing.T) {
	path := sqliteDatei(t, lesefixture())
	vorher := mustRead(t, path)
	infoVorher, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	var schreibfehler error
	got, err := recording.LeseMit(context.Background(), path, func(tx *sql.Tx, tabelle string) {
		if tabelle != "meta" {
			return
		}
		_, schreibfehler = tx.Exec(`CREATE TABLE neu (a TEXT)`)
	})
	if err != nil || !reflect.DeepEqual(got, lesefixture()) {
		t.Fatalf("Lesen: %v", err)
	}
	if schreibfehler == nil || !strings.Contains(schreibfehler.Error(), "readonly") {
		t.Errorf("Schreibversuch in der Lesetransaktion: %v", schreibfehler)
	}
	infoNachher, _ := os.Stat(path)
	if nachher := mustRead(t, path); string(nachher) != string(vorher) || !infoNachher.ModTime().Equal(infoVorher.ModTime()) {
		t.Error("Datei nach dem Lesen verändert")
	}
	if es := eintraege(t, filepath.Dir(path)); !reflect.DeepEqual(es, []string{"rec.sqlite"}) {
		t.Errorf("im Verzeichnis liegt %v", es)
	}
}

// Abdeckung: LH-FA-22/Boundary — die Datei wird in einer einzigen Lesetransaktion gelesen: Zwischen
// dem Lesen zweier Tabellen kann kein anderer Prozess einen Schreibvorgang abschließen (die
// Sperre der Lesetransaktion hält bis zum Ende); eine Session, die nach dem Lesen ergänzt
// wird, trägt der nächste Lauf.
func TestLesenEineLesetransaktion(t *testing.T) {
	ctx := context.Background()
	s, path := neueDatei(t, sessionsMit(1))
	var fehler error
	got, err := recording.LeseMit(ctx, path, func(_ *sql.Tx, tabelle string) {
		if tabelle != "session" {
			return
		}
		schreiber, err := sql.Open("sqlite", "file:"+path+"?mode=rw&_txlock=immediate&_pragma=busy_timeout(0)")
		if err != nil {
			fehler = err
			return
		}
		defer func() { _ = schreiber.Close() }()
		_, fehler = schreiber.Exec(`UPDATE meta SET value = value`)
		if fehler == nil {
			fehler = errors.New("Schreiber kam während des Lesens durch")
		}
	})
	if err != nil || len(got.Sessions) != 1 {
		t.Fatalf("Lesen: %v, %d Sessions", err, len(got.Sessions))
	}
	if fehler == nil || !strings.Contains(fehler.Error(), "locked") {
		t.Errorf("Schreiben während des Lesens: %v", fehler)
	}
	if err := s.Write(ctx, path, sessionsMit(2)); err != nil {
		t.Fatalf("Nachholen: %v", err)
	}
	if got := ladeOK(t, path); len(got.Sessions) != 2 {
		t.Errorf("nach dem Nachholen %d Sessions", len(got.Sessions))
	}
}

// Abdeckung: LH-FA-22/Boundary — ein Schreiber mit offener, nicht abgeschlossener Transaktion hält das
// Lesen nicht auf, und seine Session fehlt im Gelesenen.
func TestLesenOffeneSchreibtransaktion(t *testing.T) {
	ctx := context.Background()
	path := sqliteDatei(t, sessionsMit(1))
	schreiber, err := sql.Open("sqlite", "file:"+path+"?mode=rw&_pragma=busy_timeout(0)")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = schreiber.Close() }()
	tx, err := schreiber.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(`INSERT INTO session (id, startup, server_parameters) VALUES (2, '{}', '{}')`); err != nil {
		t.Fatal(err)
	}
	got := ladeOK(t, path)
	if len(got.Sessions) != 1 {
		t.Errorf("%d Sessions neben einer offenen Schreibtransaktion", len(got.Sessions))
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-22/Negative — ist die Datei länger als fünf Sekunden gesperrt, scheitert das Lesen
// nach etwa fünf Sekunden mit PGR-E3001 und der Ursache der Bibliothek; ohne Sperre liest es.
func TestLesenSperrfrist(t *testing.T) {
	ctx := context.Background()
	path := sqliteDatei(t, sessionsMit(1))
	ex, err := sql.Open("sqlite", "file:"+path+"?mode=rw&_txlock=exclusive&_pragma=busy_timeout(0)")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ex.Close() }()
	etx, err := ex.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	_, err = laden(path)
	dauer := time.Since(start)
	if code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "locked") {
		t.Fatalf("erwartet %s mit der Ursache der Bibliothek, erhalten %v", model.CodeRecordingIO, err)
	}
	if dauer < 4500*time.Millisecond || dauer > 8*time.Second {
		t.Errorf("gescheitert nach %v, erwartet etwa 5 s", dauer)
	}
	if err := etx.Rollback(); err != nil {
		t.Fatal(err)
	}
	ladeOK(t, path)
}

// Abdeckung: LH-FA-22/Negative — der Meldungscode eines Fehlers der Bibliothek folgt aus seinem primären
// Ergebniscode: PERM, BUSY, LOCKED, NOMEM, READONLY, INTERRUPT, IOERR, FULL, CANTOPEN und PROTOCOL
// sind PGR-E3001, auch mit erweitertem Code, jeder andere PGR-E3003; ein abgeschnittenes
// Lesen (IOERR_SHORT_READ) ist PGR-E3003.
func TestLesenStufeNachErgebniscode(t *testing.T) {
	io, inhalt := model.CodeRecordingIO, model.CodeRecordingBroken
	for _, c := range []struct {
		name string
		code int
		want string
	}{
		{"PERM", 3, io}, {"BUSY", 5, io}, {"LOCKED", 6, io}, {"NOMEM", 7, io}, {"READONLY", 8, io},
		{"INTERRUPT", 9, io}, {"IOERR", 10, io}, {"FULL", 13, io}, {"CANTOPEN", 14, io}, {"PROTOCOL", 15, io},
		{"BUSY_SNAPSHOT", 517, io}, {"LOCKED_SHAREDCACHE", 262, io}, {"IOERR_READ", 266, io},
		{"CANTOPEN_NOTEMPDIR", 270, io}, {"READONLY_ROLLBACK", 776, io},
		{"IOERR_SHORT_READ", 522, inhalt},
		{"ERROR", 1, inhalt}, {"INTERNAL", 2, inhalt}, {"ABORT", 4, inhalt}, {"NOTFOUND", 12, inhalt},
		{"CORRUPT", 11, inhalt}, {"SCHEMA", 17, inhalt}, {"TOOBIG", 18, inhalt}, {"CONSTRAINT", 19, inhalt},
		{"MISMATCH", 20, inhalt}, {"MISUSE", 21, inhalt}, {"NOLFS", 22, inhalt}, {"AUTH", 23, inhalt},
		{"RANGE", 25, inhalt}, {"NOTADB", 26, inhalt}, {"CORRUPT_VTAB", 267, inhalt},
	} {
		if got := recording.StufeCode(c.code); got != c.want {
			t.Errorf("%s (%d): %s, erwartet %s", c.name, c.code, got, c.want)
		}
	}
}

// Abdeckung: LH-FA-22/Negative — eine Datei, die nicht lesbar ist, ist PGR-E3001: eine fehlende Datei
// beim Lesen mit SQLite.Load, der Abbruch von ctx vor dem Lesen; eine Datei mit dem Kopf
// einer SQLite-Datei und anderem Inhalt ist PGR-E3003.
func TestLesenNichtLesbarUndKopfOhneInhalt(t *testing.T) {
	ctx := context.Background()
	fehlt := filepath.Join(t.TempDir(), "fehlt.sqlite")
	if _, err := (&recording.SQLite{}).Load(ctx, fehlt); code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), fehlt) {
		t.Errorf("fehlende Datei: %v", err)
	}
	if _, err := os.Stat(fehlt); err == nil {
		t.Error("Lesen legt die fehlende Datei an")
	}
	abgebrochen, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := (&recording.SQLite{}).Load(abgebrochen, sqliteDatei(t, sessionsMit(1))); code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("abgebrochen: %v", err)
	}
	for name, inhalt := range map[string]string{
		"nur der Kopf":    "SQLite format 3\x00",
		"Kopf mit Müll":   "SQLite format 3\x00" + strings.Repeat("müll", 200),
		"Kopf, 4096 Null": "SQLite format 3\x00" + strings.Repeat("\x00", 4080),
	} {
		path := filepath.Join(t.TempDir(), "x.sqlite")
		if err := os.WriteFile(path, []byte(inhalt), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := laden(path); code(err) != model.CodeRecordingBroken {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// Abdeckung: LH-FA-22/Negative — eine abgeschnittene Datei ist PGR-E3003, an der Seitengrenze wie mitten
// in einer Seite, in der Kopfzeile wie mitten im Inhalt, und sie wird nie als Aufzeichnung
// gelesen; auch die Datei, die um genau eine Seite oder ein Byte zu kurz ist.
func TestLesenAbgeschnitten(t *testing.T) {
	rec := sessionsMit(40)
	voll := mustRead(t, sqliteDatei(t, rec))
	if len(voll) < 5*4096 {
		t.Fatalf("Datei nur %d Bytes lang", len(voll))
	}
	laengen := []int{16, 17, 50, 99, 100, 101, 511, 512, 513, 4095, 4096, 4097, 4096 + 100, 2 * 4096, 2*4096 + 1, 3*4096 - 1, 3 * 4096, 3*4096 + 2000, len(voll) - 4096, len(voll) - 1}
	for _, n := range laengen {
		path := filepath.Join(t.TempDir(), "kurz.sqlite")
		if err := os.WriteFile(path, voll[:n], 0o644); err != nil {
			t.Fatal(err)
		}
		if got, err := laden(path); code(err) != model.CodeRecordingBroken {
			t.Errorf("auf %d Bytes abgeschnitten: %v (%d Sessions)", n, err, len(got.Sessions))
		}
	}
	for _, n := range []int{0, 1, 15} {
		path := filepath.Join(t.TempDir(), "kurz.sqlite")
		if err := os.WriteFile(path, voll[:n], 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := laden(path); code(err) != model.CodeRecordingBroken {
			t.Errorf("auf %d Bytes abgeschnitten: %v", n, err)
		}
	}
	if got := ladeOK(t, sqliteDatei(t, rec)); len(got.Sessions) != 40 {
		t.Errorf("unversehrt: %d Sessions", len(got.Sessions))
	}
}

// Abdeckung: LH-FA-22/Boundary — der Eingabepfad gilt wörtlich auch in der Form mit führendem //, die unter Linux
// derselbe Pfad ist wie mit einem /: Schreiben, Ergänzen und Lesen finden die Datei unter diesem Pfad
// (ohne den absoluten Pfad würde die URI die Authority tmp lesen und abgelehnt).
func TestLesenPfadMitFuehrendemDoppelstrich(t *testing.T) {
	ctx := context.Background()
	dir, err := filepath.Abs(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := "/" + filepath.Join(dir, "rec.sqlite")
	s := &recording.SQLite{}
	if err := s.Write(ctx, path, sessionsMit(1)); err != nil {
		t.Fatalf("Schreiben: %v", err)
	}
	if err := s.Write(ctx, path, sessionsMit(2)); err != nil {
		t.Fatalf("Ergänzen: %v", err)
	}
	for name, lesen := range map[string]func() (model.Recording, error){
		"Leser":  func() (model.Recording, error) { return laden(path) },
		"SQLite": func() (model.Recording, error) { return (&recording.SQLite{}).Load(ctx, path) },
	} {
		if got, err := lesen(); err != nil || !reflect.DeepEqual(got, sessionsMit(2)) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// Abdeckung: LH-FA-22/Negative — bricht der Lauf ab, während die Datei gelesen wird (nach dem Lesen einer
// Tabelle, vor der nächsten), ist das PGR-E3001 mit der Ursache und dem Pfad.
func TestLesenAbbruchWaehrendDesLesens(t *testing.T) {
	path := sqliteDatei(t, lesefixture())
	for _, tabelle := range []string{"meta", "session", "interaction"} {
		ctx, cancel := context.WithCancel(context.Background())
		_, err := recording.LeseMit(ctx, path, func(_ *sql.Tx, t string) {
			if t == tabelle {
				cancel()
			}
		})
		cancel()
		if code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "context canceled") || !strings.Contains(err.Error(), path) {
			t.Errorf("Abbruch nach %s: %v", tabelle, err)
		}
	}
}
