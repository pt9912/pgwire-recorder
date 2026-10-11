package recording_test

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-22/Boundary, LH-FA-22/Negative — Prepare des SQLite-Adapters prüft
// den Zielpfad wie der YAML-Adapter, gleich welchen Inhalts eine vorhandene Datei
// hat: ohne --force ist sie PGR-E2002, mit --force gilt sie als ersetzbar; ein
// Verzeichnis ist PGR-E3001; ein Verzeichnis, in dem sich keine Datei anlegen lässt,
// ist PGR-E3001; Prepare lässt nichts zurück.
func TestSQLitePrepareZielpfad(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	datei := filepath.Join(dir, "rec.sqlite")
	if err := os.WriteFile(datei, []byte("kein SQLite"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := (&recording.SQLite{}).Prepare(ctx, datei, false); code(err) != model.CodeOutputExists {
		t.Errorf("vorhandene Datei ohne --force: erwartet %s, erhalten %v", model.CodeOutputExists, err)
	}
	if err := (&recording.SQLite{}).Prepare(ctx, datei, true); err != nil {
		t.Errorf("vorhandene Datei mit --force: %v", err)
	}
	verzeichnis := filepath.Join(dir, "verzeichnis")
	if err := os.Mkdir(verzeichnis, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, replace := range []bool{false, true} {
		if err := (&recording.SQLite{}).Prepare(ctx, verzeichnis, replace); code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "keine reguläre Datei") {
			t.Errorf("Verzeichnis, replace=%v: erwartet %s, erhalten %v", replace, model.CodeRecordingIO, err)
		}
	}
	if err := (&recording.SQLite{}).Prepare(ctx, filepath.Join(dir, "fehlt", "rec.sqlite"), false); code(err) != model.CodeRecordingIO {
		t.Errorf("fehlendes Verzeichnis: erwartet %s, erhalten %v", model.CodeRecordingIO, err)
	}
	err := recording.PruefeSQLiteMit(filepath.Join(dir, "neu.sqlite"), false, recording.Eingriffe{Probe: func(d, _ string) (*os.File, error) {
		return nil, &fs.PathError{Op: "open", Path: d, Err: fs.ErrPermission}
	}})
	if code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "nicht beschreibbar") {
		t.Errorf("Verzeichnis nicht beschreibbar: erwartet %s, erhalten %v", model.CodeRecordingIO, err)
	}
	if es := eintraege(t, dir); len(es) != 2 {
		t.Errorf("im Verzeichnis liegt %v", es)
	}
}

// Abdeckung: LH-FA-22/Negative — ein vorhandenes Rollback-Journal <Zielpfad>-journal ist
// beim Start PGR-E3001, auch mit --force und auch ohne Zieldatei; ein anderer Name
// daneben zählt nicht.
func TestSQLitePrepareJournal(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.sqlite")
	for _, andere := range []string{"rec.sqlite-wal", "rec.sqlite-journal.alt", ".rec.sqlite-journal", "rec.sqlite.journal", "rec-journal"} {
		if err := os.WriteFile(filepath.Join(dir, andere), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, replace := range []bool{false, true} {
		if err := (&recording.SQLite{}).Prepare(ctx, path, replace); err != nil {
			t.Fatalf("ohne Journal, replace=%v: %v", replace, err)
		}
	}
	if err := os.WriteFile(path+"-journal", []byte("Rest"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, replace := range []bool{false, true} {
		err := (&recording.SQLite{}).Prepare(ctx, path, replace)
		if code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "rec.sqlite-journal vorhanden") {
			t.Errorf("Journal ohne Zieldatei, replace=%v: erwartet %s, erhalten %v", replace, model.CodeRecordingIO, err)
		}
	}
	if err := os.WriteFile(path, []byte("alt"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := (&recording.SQLite{}).Prepare(ctx, path, false); code(err) != model.CodeOutputExists {
		t.Errorf("Zieldatei und Journal ohne --force: erwartet %s, erhalten %v", model.CodeOutputExists, err)
	}
	if err := (&recording.SQLite{}).Prepare(ctx, path, true); code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "-journal vorhanden") {
		t.Errorf("Zieldatei und Journal mit --force: erwartet %s, erhalten %v", model.CodeRecordingIO, err)
	}
}

// Abdeckung: LH-FA-22/Negative — das Journal wird nach dem Zielpfad geprüft: Ist der
// Zielpfad ein Verzeichnis, nennt die Meldung das und nicht das Journal. Vorhanden
// heißt wie beim Zielpfad: Ein Journal, das ein Verzeichnis oder eine Verknüpfung auf
// eine Datei ist, ist vorhanden; eine Verknüpfung ins Leere ist es nicht. Lässt sich sein
// Zustand nicht feststellen (hier ein zu langer Name), ist das PGR-E3001.
func TestSQLitePrepareJournalReihenfolge(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	verzeichnis := filepath.Join(dir, "ziel")
	if err := os.Mkdir(verzeichnis, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(verzeichnis+"-journal", nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := (&recording.SQLite{}).Prepare(ctx, verzeichnis, true); code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "keine reguläre Datei") || strings.Contains(err.Error(), "-journal") {
		t.Errorf("Verzeichnis und Journal: %v", err)
	}
	datei := filepath.Join(dir, "datei")
	if err := os.WriteFile(datei, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	for name, c := range map[string]struct {
		anlegen   func(string) error
		vorhanden bool
	}{
		"Verzeichnis":           {func(p string) error { return os.Mkdir(p, 0o755) }, true},
		"Verknüpfung auf Datei": {func(p string) error { return os.Symlink(datei, p) }, true},
		"Verknüpfung ins Leere": {func(p string) error { return os.Symlink(filepath.Join(dir, "fehlt"), p) }, false},
	} {
		path := filepath.Join(dir, strings.ReplaceAll(name, " ", "_")+".sqlite")
		if err := c.anlegen(path + "-journal"); err != nil {
			t.Fatal(err)
		}
		err := (&recording.SQLite{}).Prepare(ctx, path, true)
		if c.vorhanden && (code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "-journal vorhanden")) {
			t.Errorf("Journal als %s: erwartet %s, erhalten %v", name, model.CodeRecordingIO, err)
		}
		if !c.vorhanden && err != nil {
			t.Errorf("Journal als %s: %v", name, err)
		}
	}
	lang := filepath.Join(dir, strings.Repeat("a", 250))
	err := recording.PruefeSQLiteMit(lang, true, recording.Eingriffe{Probe: func(string, string) (*os.File, error) {
		return os.CreateTemp(t.TempDir(), "probe")
	}})
	if code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "nicht prüfbar") || !strings.Contains(err.Error(), syscall.ENAMETOOLONG.Error()) {
		t.Errorf("Journal nicht prüfbar: erwartet %s mit „nicht prüfbar“, erhalten %v", model.CodeRecordingIO, err)
	}
}

// Abdeckung: LH-FA-22/Boundary — die temporäre Datei des ersten Schreibvorgangs liegt im
// Verzeichnis der Zieldatei und heißt .<Name der Zieldatei>.<16 Hexziffern>.tmp; bis
// zum ersten Schreibvorgang bleibt der Zielpfad unberührt, danach ist die Datei
// verschoben und keine temporäre Datei bleibt.
func TestSQLiteNameDerTempDatei(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.sqlite")
	var alt string
	s := &recording.SQLite{}
	err := recording.SQLiteSchreibeMit(context.Background(), s, path, sessionsMit(1), recording.SQLiteEingriffe{Verschieben: func(a, n string) error {
		alt = a
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("Zielpfad vor dem Verschieben: %v", err)
		}
		return os.Rename(a, n)
	}})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(alt) != dir || !regexp.MustCompile(`^\.rec\.sqlite\.[0-9a-f]{16}\.tmp$`).MatchString(filepath.Base(alt)) {
		t.Fatalf("temporäre Datei %q", alt)
	}
	if es := eintraege(t, dir); len(es) != 1 || es[0] != "rec.sqlite" {
		t.Fatalf("im Verzeichnis liegt %v", es)
	}
}

// Abdeckung: LH-FA-22/Boundary — ist der Name der temporären Datei belegt,
// überschreibt der erste Schreibvorgang ihn nicht, sondern zieht einen neuen.
func TestSQLiteBelegterName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.sqlite")
	belegt := filepath.Join(dir, ".rec.sqlite.0000000000000000.tmp")
	if err := os.WriteFile(belegt, []byte("fremd"), 0o644); err != nil {
		t.Fatal(err)
	}
	aufrufe := 0
	err := recording.SQLiteSchreibeMit(context.Background(), &recording.SQLite{}, path, sessionsMit(1),
		recording.SQLiteEingriffe{Zufall: festerZufall(&aufrufe, func(n int) byte { return byte(n - 1) })})
	if err != nil {
		t.Fatal(err)
	}
	if got := inhalt(t, belegt); got != "fremd" {
		t.Fatalf("belegter Name überschrieben: %q", got)
	}
	if aufrufe != 2 {
		t.Fatalf("Zufall %d-mal gezogen, erwartet 2", aufrufe)
	}
	if got := sessionIDs(t, path); got != "1" {
		t.Fatalf("Sessions %q", got)
	}
}

// Abdeckung: LH-FA-22/Negative — scheitert das Anlegen der temporären Datei anders als
// an einem belegten Namen (ein zu langer Name, ein fehlendes Verzeichnis), ist das sofort
// PGR-E3001 mit dem Fehler des Betriebssystems als Ursache; findet der erste
// Schreibvorgang nach zehn Versuchen keinen freien Namen, ist das PGR-E3001 nach genau
// zehn Zügen.
func TestSQLiteAnlegenScheitert(t *testing.T) {
	for _, tc := range []struct {
		name    string
		path    string
		ursache error
	}{
		{"Name zu lang", filepath.Join(t.TempDir(), strings.Repeat("a", 240)+".sqlite"), syscall.ENAMETOOLONG},
		{"Verzeichnis fehlt", filepath.Join(t.TempDir(), "fehlt", "rec.sqlite"), syscall.ENOENT},
	} {
		aufrufe := 0
		err := recording.SQLiteSchreibeMit(context.Background(), &recording.SQLite{}, tc.path, sessionsMit(1),
			recording.SQLiteEingriffe{Zufall: festerZufall(&aufrufe, func(n int) byte { return byte(n) })})
		ms := model.Meldungen(err)
		if len(ms) != 1 || ms[0].Code != model.CodeRecordingIO || !strings.Contains(ms[0].Text, tc.ursache.Error()) || strings.Contains(ms[0].Text, "kein freier Name") {
			t.Errorf("%s: Meldungen %+v", tc.name, ms)
		}
		if aufrufe != 1 {
			t.Errorf("%s: Zufall %d-mal gezogen, erwartet 1", tc.name, aufrufe)
		}
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.sqlite")
	if err := os.WriteFile(filepath.Join(dir, ".rec.sqlite.0000000000000000.tmp"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	aufrufe := 0
	zufall := festerZufall(&aufrufe, func(n int) byte {
		if n > 10 {
			return 1
		}
		return 0
	})
	err := recording.SQLiteSchreibeMit(context.Background(), &recording.SQLite{}, path, sessionsMit(1), recording.SQLiteEingriffe{Zufall: zufall})
	if code(err) != model.CodeRecordingIO || aufrufe != 10 {
		t.Fatalf("erwartet %s nach 10 Zügen, erhalten %v nach %d", model.CodeRecordingIO, err, aufrufe)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("Zieldatei angelegt: %v", err)
	}
}

// Abdeckung: LH-FA-22/Negative — schlägt beim ersten Schreibvorgang das Setzen der Rechte,
// das Schließen der temporären Datei, das Anlegen von Schema und meta oder das
// Verschieben fehl, ist das PGR-E3001; die vorhandene Zieldatei bleibt unverändert, die
// temporäre Datei ist entfernt, und der nächste Schreibvorgang legt die Datei an, als wäre
// nichts geschehen.
func TestSQLiteErsterSchreibvorgangScheitert(t *testing.T) {
	fehler := errors.New("eingespielter Fehler")
	for _, tc := range []struct {
		name      string
		eingriffe recording.SQLiteEingriffe
	}{
		// Die Zieldatei ist vorhanden; nur dann setzt der Schreibvorgang die Rechte.
		{"Rechte", recording.SQLiteEingriffe{Rechte: func(*os.File, os.FileMode) error { return fehler }}},
		{"Schliessen", recording.SQLiteEingriffe{Schliessen: func(f *os.File) error {
			_ = f.Close()
			return fehler
		}}},
		{"Anlegen", recording.SQLiteEingriffe{Anlegen: func(context.Context, string, model.Recording) error { return fehler }}},
		{"Verschieben", recording.SQLiteEingriffe{Verschieben: func(string, string) error { return fehler }}},
	} {
		dir := t.TempDir()
		path := filepath.Join(dir, "rec.sqlite")
		if err := os.WriteFile(path, []byte("alt"), 0o644); err != nil {
			t.Fatal(err)
		}
		s := &recording.SQLite{}
		err := recording.SQLiteSchreibeMit(context.Background(), s, path, sessionsMit(1), tc.eingriffe)
		if ms := model.Meldungen(err); len(ms) != 1 || ms[0].Code != model.CodeRecordingIO || !strings.Contains(ms[0].Text, fehler.Error()) {
			t.Errorf("%s: Meldungen %+v", tc.name, ms)
		}
		if got := inhalt(t, path); got != "alt" {
			t.Errorf("%s: Zieldatei verändert: %q", tc.name, got)
		}
		if es := eintraege(t, dir); len(es) != 1 {
			t.Errorf("%s: im Verzeichnis liegt %v", tc.name, es)
		}
		if err := s.Write(context.Background(), path, sessionsMit(1)); err != nil {
			t.Errorf("%s: nächster Schreibvorgang: %v", tc.name, err)
		} else if got := sessionIDs(t, path); got != "1" {
			t.Errorf("%s: Sessions %q", tc.name, got)
		}
	}
}

// Abdeckung: LH-FA-22/Negative — scheitert nach einem Fehlschlag auch das Entfernen der
// temporären Datei, bleibt sie liegen, und der Fehler des Entfernens folgt als Ursache
// derselben Meldung PGR-E3001.
func TestSQLiteEntfernenScheitert(t *testing.T) {
	entfernen := func(string) error { return errors.New("Entfernen eingespielt gescheitert") }
	for name, e := range map[string]recording.SQLiteEingriffe{
		"Anlegen": {Anlegen: func(context.Context, string, model.Recording) error {
			return errors.New("Anlegen eingespielt gescheitert")
		}, Entfernen: entfernen},
		"Verschieben": {Verschieben: func(string, string) error { return errors.New("Verschieben eingespielt gescheitert") }, Entfernen: entfernen},
	} {
		dir := t.TempDir()
		err := recording.SQLiteSchreibeMit(context.Background(), &recording.SQLite{}, filepath.Join(dir, "rec.sqlite"), sessionsMit(1), e)
		ms := model.Meldungen(err)
		if len(ms) != 1 || ms[0].Code != model.CodeRecordingIO ||
			!strings.Contains(ms[0].Text, name+" eingespielt gescheitert") || !strings.Contains(ms[0].Text, "Entfernen eingespielt gescheitert") {
			t.Errorf("%s: Meldungen %+v", name, ms)
		}
		if es := eintraege(t, dir); len(es) != 1 || !strings.HasSuffix(es[0], ".tmp") {
			t.Errorf("%s: im Verzeichnis liegt %v, erwartet die temporäre Datei", name, es)
		}
	}
}

// Abdeckung: LH-FA-22/Boundary — eine temporäre Datei oder Probedatei eines früheren
// Laufs überschreiben, entfernen und melden Prepare und Write nicht.
func TestSQLiteUebrigGebliebeneDateien(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.sqlite")
	reste := map[string]string{
		filepath.Join(dir, ".rec.sqlite.0123456789abcdef.tmp"): "rest-tmp",
		filepath.Join(dir, ".rec.sqlite.123456.probe"):         "rest-probe",
	}
	for p, s := range reste {
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	s := &recording.SQLite{}
	if err := s.Prepare(ctx, path, false); err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 2; i++ {
		if err := s.Write(ctx, path, sessionsMit(i)); err != nil {
			t.Fatal(err)
		}
	}
	for p, s := range reste {
		if got := inhalt(t, p); got != s {
			t.Errorf("%s: %q statt %q", p, got, s)
		}
	}
	if es := eintraege(t, dir); len(es) != 3 {
		t.Errorf("im Verzeichnis liegt %v", es)
	}
}

// Abdeckung: LH-FA-22/Boundary — mit --force ersetzt der erste Schreibvorgang eine
// vorhandene Datei gleich welchen Inhalts, auch eine SQLite-Datei mit anderen Sessions,
// und ergänzt sie nie; eine symbolische Verknüpfung unter dem Zielpfad ersetzt er, nicht
// ihr Ziel; ein Zielpfad, der erst nach dem Start entsteht, ersetzt er ohne Prüfung.
func TestSQLiteErsetzt(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	_, fremd := neueDatei(t, sessionsMit(3))
	for name, vorlage := range map[string][]byte{"Text": []byte("alt"), "SQLite": mustRead(t, fremd)} {
		path := filepath.Join(dir, name+".sqlite")
		if err := os.WriteFile(path, vorlage, 0o644); err != nil {
			t.Fatal(err)
		}
		s := &recording.SQLite{}
		if err := s.Prepare(ctx, path, true); err != nil {
			t.Fatal(err)
		}
		if err := s.Write(ctx, path, sessionsMit(1)); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got := sessionIDs(t, path); got != "1" {
			t.Errorf("%s: Sessions %q, erwartet nur die der neuen Aufzeichnung", name, got)
		}
	}
	ziel := filepath.Join(dir, "ziel.txt")
	if err := os.WriteFile(ziel, []byte("alt"), 0o644); err != nil {
		t.Fatal(err)
	}
	verknuepfung := filepath.Join(dir, "verknuepfung.sqlite")
	if err := os.Symlink(ziel, verknuepfung); err != nil {
		t.Fatal(err)
	}
	if err := (&recording.SQLite{}).Write(ctx, verknuepfung, sessionsMit(1)); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Lstat(verknuepfung); err != nil || !info.Mode().IsRegular() {
		t.Errorf("Verknüpfung nach Write: %v %v", info, err)
	}
	if got := inhalt(t, ziel); got != "alt" {
		t.Errorf("Ziel der Verknüpfung verändert: %q", got)
	}
	nachStart := filepath.Join(dir, "nach-start.sqlite")
	s := &recording.SQLite{}
	if err := s.Prepare(ctx, nachStart, false); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nachStart, []byte("dazwischen"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.Write(ctx, nachStart, sessionsMit(1)); err != nil {
		t.Fatal(err)
	}
	if got := sessionIDs(t, nachStart); got != "1" {
		t.Errorf("Pfad nach dem Start: Sessions %q", got)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Abdeckung: LH-FA-22/Boundary — eine neue Datei erhält die Rechte 0666 nach der umask,
// unter der umask 022 also 0644; eine ersetzte behält ihre Zugriffsrechte, auch die
// Zieldatei einer ersetzten Verknüpfung auf eine Datei die des Ziels; das Ergänzen
// ändert die Rechte nicht.
func TestSQLiteRechte(t *testing.T) {
	ctx := context.Background()
	alt := syscall.Umask(0o022)
	defer syscall.Umask(alt)
	dir := t.TempDir()
	neu := filepath.Join(dir, "neu.sqlite")
	s := &recording.SQLite{}
	if err := s.Write(ctx, neu, sessionsMit(1)); err != nil {
		t.Fatal(err)
	}
	if err := s.Write(ctx, neu, sessionsMit(2)); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(neu); info.Mode().Perm() != 0o644 {
		t.Errorf("neue Datei unter umask 022: %v", info.Mode().Perm())
	}
	ersetzt := filepath.Join(dir, "ersetzt.sqlite")
	if err := os.WriteFile(ersetzt, []byte("alt"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := (&recording.SQLite{}).Write(ctx, ersetzt, sessionsMit(1)); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(ersetzt); info.Mode().Perm() != 0o600 {
		t.Errorf("ersetzte Datei: %v", info.Mode().Perm())
	}
	ziel := filepath.Join(dir, "ziel.txt")
	if err := os.WriteFile(ziel, []byte("alt"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(ziel, 0o640); err != nil {
		t.Fatal(err)
	}
	verknuepfung := filepath.Join(dir, "verknuepfung.sqlite")
	if err := os.Symlink(ziel, verknuepfung); err != nil {
		t.Fatal(err)
	}
	if err := (&recording.SQLite{}).Write(ctx, verknuepfung, sessionsMit(1)); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Lstat(verknuepfung); err != nil || !info.Mode().IsRegular() || info.Mode().Perm() != 0o640 {
		t.Errorf("Zieldatei nach dem Ersetzen der Verknüpfung: %v %v", info, err)
	}
}

// Abdeckung: LH-FA-22/Boundary — die Datei trägt die Kopfzeile einer SQLite-Datei: Die
// ersten 16 Bytes sind SQLite format 3 und ein NUL-Byte.
func TestSQLiteKopfzeile(t *testing.T) {
	_, path := neueDatei(t, sessionsMit(0))
	if kopf := mustRead(t, path)[:16]; !bytes.Equal(kopf, []byte("SQLite format 3\x00")) {
		t.Fatalf("erste 16 Bytes: %q", kopf)
	}
	if got := sessionIDs(t, path); got != "" {
		t.Fatalf("Aufzeichnung ohne Sessions trägt %q", got)
	}
}

// Ein Aufruf mit weniger Sessions, als die Datei trägt, ist ein innerer Fehler
// (PGR-E1000) und lässt die Datei unberührt.
func TestSQLiteWenigerSessionsAlsDieDatei(t *testing.T) {
	s, path := neueDatei(t, sessionsMit(2))
	vorher := mustRead(t, path)
	err := s.Write(context.Background(), path, sessionsMit(1))
	if code(err) != model.CodeInternal || !bytes.Equal(mustRead(t, path), vorher) {
		t.Fatalf("erwartet %s und unberührte Datei, erhalten %v", model.CodeInternal, err)
	}
}

// Abdeckung: LH-FA-22/Negative — eine Session mit einer Nummer, die die Datei schon trägt,
// ist PGR-E3001 (Primärschlüssel von session), und die Datei bleibt auf dem Stand vor dem
// Aufruf.
func TestSQLiteDoppelteSessionNummer(t *testing.T) {
	s, path := neueDatei(t, sessionsMit(1))
	vorher := mustRead(t, path)
	rec := sessionsMit(2)
	rec.Sessions[1].ID = 1
	if err := s.Write(context.Background(), path, rec); code(err) != model.CodeRecordingIO {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingIO, err)
	}
	if !bytes.Equal(mustRead(t, path), vorher) {
		t.Fatal("Datei verändert")
	}
}

// Abdeckung: LH-FA-22/Negative — trägt die Datei eine Lücke in den Nummern (hier 1 und 3),
// ist der nächste Schreibvorgang PGR-E3001, und die Datei bleibt unberührt.
func TestSQLiteLueckeInDenNummern(t *testing.T) {
	ctx := context.Background()
	rec := sessionsMit(3)
	rec.Sessions = append(rec.Sessions[:1], rec.Sessions[2:]...)
	s, path := neueDatei(t, rec)
	if got := sessionIDs(t, path); got != "1,3" {
		t.Fatalf("Sessions %q", got)
	}
	vorher := mustRead(t, path)
	for _, n := range []int{2, 3} {
		if err := s.Write(ctx, path, sessionsMit(n)); code(err) != model.CodeRecordingIO {
			t.Fatalf("Write mit %d Sessions: erwartet %s, erhalten %v", n, model.CodeRecordingIO, err)
		}
	}
	if !bytes.Equal(mustRead(t, path), vorher) {
		t.Fatal("Datei verändert")
	}
}

// Abdeckung: LH-FA-22/Boundary — eine fremde temporäre Datei samt Journal
// (.<Name>.<Hex>.tmp-journal) im Zielverzeichnis bleibt nach Prepare und Write
// unverändert, ebenso die temporäre Datei selbst.
func TestSQLiteFremdesTempJournal(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.sqlite")
	reste := map[string]string{
		filepath.Join(dir, ".rec.sqlite.0123456789abcdef.tmp"):         "rest",
		filepath.Join(dir, ".rec.sqlite.0123456789abcdef.tmp-journal"): "rest-journal",
	}
	for p, s := range reste {
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	s := &recording.SQLite{}
	if err := s.Prepare(ctx, path, false); err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 2; i++ {
		if err := s.Write(ctx, path, sessionsMit(i)); err != nil {
			t.Fatal(err)
		}
	}
	for p, want := range reste {
		if got := inhalt(t, p); got != want {
			t.Errorf("%s: %q statt %q", p, got, want)
		}
	}
}
