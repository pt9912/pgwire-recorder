package recording_test

import (
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

// eintraege liefert die Namen im Verzeichnis dir.
func eintraege(t *testing.T, dir string) []string {
	t.Helper()
	es, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range es {
		out = append(out, e.Name())
	}
	return out
}

// inhalt liefert den Inhalt der Datei path.
func inhalt(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// Abdeckung: LH-FA-07/Negative — ist der vorhandene --output-Pfad keine
// reguläre Datei (ein Verzeichnis, eine benannte Pipe, eine Verknüpfung auf ein
// Verzeichnis), ist das beim Start PGR-E3001, ohne und mit --force; der Pfad
// bleibt, was er war.
func TestPrepareKeineRegulaereDatei(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	verzeichnis := filepath.Join(dir, "verzeichnis")
	if err := os.Mkdir(verzeichnis, 0o755); err != nil {
		t.Fatal(err)
	}
	pipe := filepath.Join(dir, "pipe")
	if err := syscall.Mkfifo(pipe, 0o644); err != nil {
		t.Fatal(err)
	}
	verknuepfung := filepath.Join(dir, "verknuepfung")
	if err := os.Symlink(verzeichnis, verknuepfung); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{verzeichnis, pipe, verknuepfung} {
		for _, replace := range []bool{false, true} {
			err := recording.YAML{}.Prepare(ctx, path, replace)
			if code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "keine reguläre Datei") {
				t.Errorf("%s, replace=%v: erwartet %s mit „keine reguläre Datei“, erhalten %v", filepath.Base(path), replace, model.CodeRecordingIO, err)
			}
		}
	}
	if info, err := os.Stat(verzeichnis); err != nil || !info.IsDir() {
		t.Fatalf("Verzeichnis nach Prepare: %v %v", info, err)
	}
}

// Abdeckung: LH-FA-08/Boundary — eine symbolische Verknüpfung als --output zählt
// nach ihrem Ziel: auf eine vorhandene Datei ist sie vorhanden (PGR-E2002 ohne
// --force), ins Leere ist sie nicht vorhanden (kein Fehler ohne --force).
func TestPrepareVerknuepfung(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	ziel := filepath.Join(dir, "ziel.yaml")
	if err := os.WriteFile(ziel, []byte("alt"), 0o644); err != nil {
		t.Fatal(err)
	}
	aufDatei := filepath.Join(dir, "auf-datei.yaml")
	if err := os.Symlink(ziel, aufDatei); err != nil {
		t.Fatal(err)
	}
	if err := (recording.YAML{}).Prepare(ctx, aufDatei, false); code(err) != model.CodeOutputExists {
		t.Errorf("Verknüpfung auf eine Datei: erwartet %s, erhalten %v", model.CodeOutputExists, err)
	}
	insLeere := filepath.Join(dir, "ins-leere.yaml")
	if err := os.Symlink(filepath.Join(dir, "fehlt.yaml"), insLeere); err != nil {
		t.Fatal(err)
	}
	if err := (recording.YAML{}).Prepare(ctx, insLeere, false); err != nil {
		t.Errorf("Verknüpfung ins Leere: %v", err)
	}
}

// Abdeckung: LH-FA-07/Negative — lässt sich nicht feststellen, ob der
// --output-Pfad besteht (eine Verknüpfung auf sich selbst, ein Pfad unter einer
// Datei), ist das beim Start PGR-E3001.
func TestPrepareNichtPruefbar(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	schleife := filepath.Join(dir, "schleife.yaml")
	if err := os.Symlink(schleife, schleife); err != nil {
		t.Fatal(err)
	}
	datei := filepath.Join(dir, "datei")
	if err := os.WriteFile(datei, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{schleife, filepath.Join(datei, "rec.yaml")} {
		err := recording.YAML{}.Prepare(ctx, path, true)
		if code(err) != model.CodeRecordingIO || !strings.Contains(err.Error(), "nicht prüfbar") {
			t.Errorf("%s: erwartet %s mit „nicht prüfbar“, erhalten %v", path, model.CodeRecordingIO, err)
		}
	}
}

// Abdeckung: LH-FA-07/Negative — fehlt das Verzeichnis von --output, ist das
// beim Start PGR-E3001, und Prepare legt es nicht an; in einem vorhandenen
// Verzeichnis bleibt nach Prepare keine Probedatei zurück.
func TestPrepareVerzeichnis(t *testing.T) {
	ctx := context.Background()
	fehlt := filepath.Join(t.TempDir(), "fehlt")
	if err := (recording.YAML{}).Prepare(ctx, filepath.Join(fehlt, "rec.yaml"), false); code(err) != model.CodeRecordingIO {
		t.Fatalf("fehlendes Verzeichnis: erwartet %s, erhalten %v", model.CodeRecordingIO, err)
	}
	if _, err := os.Stat(fehlt); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("fehlendes Verzeichnis angelegt: %v", err)
	}
	dir := t.TempDir()
	if err := (recording.YAML{}).Prepare(ctx, filepath.Join(dir, "rec.yaml"), false); err != nil {
		t.Fatal(err)
	}
	if es := eintraege(t, dir); len(es) != 0 {
		t.Fatalf("nach Prepare liegt %v im Verzeichnis", es)
	}
}

// Abdeckung: LH-FA-07/Negative — lässt sich im Verzeichnis von --output keine
// Datei anlegen, ist das beim Start PGR-E3001.
func TestPrepareVerzeichnisNichtBeschreibbar(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rec.yaml")
	err := recording.PruefeMit(path, false, recording.Eingriffe{Probe: func(dir, _ string) (*os.File, error) {
		return nil, &fs.PathError{Op: "open", Path: dir, Err: fs.ErrPermission}
	}})
	if code(err) != model.CodeRecordingIO {
		t.Fatalf("erwartet %s, erhalten %v", model.CodeRecordingIO, err)
	}
}

// Abdeckung: LH-FA-07/Boundary — die Probedatei liegt im Verzeichnis der
// Zieldatei und heißt .<Name der Zieldatei>.<Zufallsteil>.probe; Prepare
// entfernt sie nach dem Anlegen.
func TestPrepareProbedatei(t *testing.T) {
	dir := t.TempDir()
	var entfernt []string
	err := recording.PruefeMit(filepath.Join(dir, "rec.yaml"), false, recording.Eingriffe{Entfernen: func(name string) error {
		entfernt = append(entfernt, name)
		return nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	es := eintraege(t, dir)
	if len(entfernt) != 1 || len(es) != 1 || entfernt[0] != filepath.Join(dir, es[0]) ||
		!regexp.MustCompile(`^\.rec\.yaml\.[^.]+\.probe$`).MatchString(es[0]) {
		t.Fatalf("entfernt %v, im Verzeichnis %v", entfernt, es)
	}
}

// Abdeckung: LH-FA-08/Boundary — mit --force ersetzt das Schreiben eine
// symbolische Verknüpfung unter --output durch die Aufzeichnung, nicht ihr Ziel;
// ins Leere bleibt das Ziel fehlend.
func TestWriteVerknuepfungErsetzt(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	ziel := filepath.Join(dir, "ziel.yaml")
	if err := os.WriteFile(ziel, []byte("alt"), 0o644); err != nil {
		t.Fatal(err)
	}
	fehlt := filepath.Join(dir, "fehlt.yaml")
	for name, auf := range map[string]string{"auf-datei.yaml": ziel, "ins-leere.yaml": fehlt} {
		path := filepath.Join(dir, name)
		if err := os.Symlink(auf, path); err != nil {
			t.Fatal(err)
		}
		if err := (recording.YAML{}).Write(ctx, path, beispiel()); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if info, err := os.Lstat(path); err != nil || !info.Mode().IsRegular() {
			t.Errorf("%s nach Write: %v %v", name, info, err)
		}
	}
	if got := inhalt(t, ziel); got != "alt" {
		t.Errorf("Ziel der Verknüpfung verändert: %q", got)
	}
	if _, err := os.Lstat(fehlt); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Ziel der Verknüpfung ins Leere angelegt: %v", err)
	}
}

// Abdeckung: LH-FA-07/Boundary — die temporäre Datei liegt im Verzeichnis der
// Zieldatei und heißt .<Name der Zieldatei>.<16 Hexziffern>.tmp.
func TestWriteNameDerTempDatei(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.yaml")
	var alt string
	err := recording.SchreibeMit(path, beispiel(), recording.Eingriffe{Verschieben: func(a, n string) error {
		alt = a
		return os.Rename(a, n)
	}})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(alt) != dir || !regexp.MustCompile(`^\.rec\.yaml\.[0-9a-f]{16}\.tmp$`).MatchString(filepath.Base(alt)) {
		t.Fatalf("temporäre Datei %q", alt)
	}
}

// festerZufall liefert Zufall, der beim n-ten Aufruf (ab 1) jedes Byte auf
// wert(n) setzt, und zählt die Aufrufe in aufrufe.
func festerZufall(aufrufe *int, wert func(int) byte) func([]byte) (int, error) {
	return func(b []byte) (int, error) {
		*aufrufe++
		for i := range b {
			b[i] = wert(*aufrufe)
		}
		return len(b), nil
	}
}

// Abdeckung: LH-FA-07/Boundary — ist der Name der temporären Datei belegt,
// überschreibt Write ihn nicht, sondern zieht einen neuen.
func TestWriteBelegterName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.yaml")
	belegt := filepath.Join(dir, ".rec.yaml.0000000000000000.tmp")
	if err := os.WriteFile(belegt, []byte("fremd"), 0o644); err != nil {
		t.Fatal(err)
	}
	aufrufe := 0
	err := recording.SchreibeMit(path, beispiel(), recording.Eingriffe{Zufall: festerZufall(&aufrufe, func(n int) byte { return byte(n - 1) })})
	if err != nil {
		t.Fatal(err)
	}
	if got := inhalt(t, belegt); got != "fremd" {
		t.Fatalf("belegter Name überschrieben: %q", got)
	}
	if aufrufe != 2 {
		t.Fatalf("Zufall %d-mal gezogen, erwartet 2", aufrufe)
	}
	if _, err := (recording.YAML{}).Load(context.Background(), path); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-07/Negative — findet Write nach zehn Versuchen keinen freien
// Namen für die temporäre Datei, ist das PGR-E3001; es zieht nicht ein elftes
// Mal.
func TestWriteZehnBelegteNamen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.yaml")
	if err := os.WriteFile(filepath.Join(dir, ".rec.yaml.0000000000000000.tmp"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	aufrufe := 0
	// Ab dem elften Zug wäre der Name frei.
	zufall := festerZufall(&aufrufe, func(n int) byte {
		if n > 10 {
			return 1
		}
		return 0
	})
	err := recording.SchreibeMit(path, beispiel(), recording.Eingriffe{Zufall: zufall})
	if code(err) != model.CodeRecordingIO || aufrufe != 10 {
		t.Fatalf("erwartet %s nach 10 Zügen, erhalten %v nach %d", model.CodeRecordingIO, err, aufrufe)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Zieldatei angelegt: %v", err)
	}
}

// Abdeckung: LH-FA-07/Negative — schlägt das Anlegen, Schreiben (auch nach einem
// Teil der Daten), Synchronisieren oder Verschieben fehl, ist das PGR-E3001; die
// Zieldatei bleibt unverändert, und die temporäre Datei ist entfernt.
func TestWriteFehlschlag(t *testing.T) {
	fehler := errors.New("eingespielter Fehler")
	lang := strings.Repeat("a", 240) + ".yaml"
	for _, tc := range []struct {
		name      string
		datei     string
		eingriffe recording.Eingriffe
	}{
		// Der Name der temporären Datei ist länger als 255 Bytes.
		{"Anlegen", lang, recording.Eingriffe{}},
		{"Schreiben", "rec.yaml", recording.Eingriffe{Schreiben: func(f *os.File, b []byte) (int, error) {
			n, _ := f.Write(b[:len(b)/2])
			return n, fehler
		}}},
		{"Synchronisieren", "rec.yaml", recording.Eingriffe{Synchronisieren: func(*os.File) error { return fehler }}},
		{"Verschieben", "rec.yaml", recording.Eingriffe{Verschieben: func(string, string) error { return fehler }}},
	} {
		dir := t.TempDir()
		path := filepath.Join(dir, tc.datei)
		if err := os.WriteFile(path, []byte("alt"), 0o644); err != nil {
			t.Fatal(err)
		}
		err := recording.SchreibeMit(path, beispiel(), tc.eingriffe)
		if code(err) != model.CodeRecordingIO {
			t.Errorf("%s: erwartet %s, erhalten %v", tc.name, model.CodeRecordingIO, err)
		}
		if got := inhalt(t, path); got != "alt" {
			t.Errorf("%s: Zieldatei verändert: %q", tc.name, got)
		}
		if es := eintraege(t, dir); len(es) != 1 {
			t.Errorf("%s: im Verzeichnis liegt %v", tc.name, es)
		}
	}
}

// Abdeckung: LH-FA-07/Negative — scheitert nach einem Fehlschlag des Schreibens
// oder Verschiebens auch das Entfernen der temporären Datei, bleibt sie liegen,
// und der Fehler des Entfernens folgt als Ursache derselben Meldung PGR-E3001.
func TestWriteEntfernenScheitert(t *testing.T) {
	entfernen := func(string) error { return errors.New("Entfernen eingespielt gescheitert") }
	for name, e := range map[string]recording.Eingriffe{
		"Schreiben":   {Schreiben: func(*os.File, []byte) (int, error) { return 0, errors.New("Schreiben eingespielt gescheitert") }, Entfernen: entfernen},
		"Verschieben": {Verschieben: func(string, string) error { return errors.New("Verschieben eingespielt gescheitert") }, Entfernen: entfernen},
	} {
		dir := t.TempDir()
		path := filepath.Join(dir, "rec.yaml")
		err := recording.SchreibeMit(path, beispiel(), e)
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

// Abdeckung: LH-FA-07/Boundary — eine temporäre Datei oder Probedatei eines
// früheren Laufs überschreiben, entfernen und melden Prepare und Write nicht.
func TestUebrigGebliebeneDateien(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "rec.yaml")
	reste := map[string]string{
		filepath.Join(dir, ".rec.yaml.0123456789abcdef.tmp"): "rest-tmp",
		filepath.Join(dir, ".rec.yaml.123456.probe"):         "rest-probe",
	}
	for p, s := range reste {
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := (recording.YAML{}).Prepare(ctx, path, false); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := (recording.YAML{}).Write(ctx, path, beispiel()); err != nil {
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

// Abdeckung: LH-FA-08/Boundary — ein --output-Pfad, der erst nach dem Start
// entsteht, ersetzt der nächste Schreibvorgang ohne Prüfung.
func TestWritePfadNachDemStart(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "rec.yaml")
	if err := (recording.YAML{}).Prepare(ctx, path, false); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("dazwischen"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := (recording.YAML{}).Write(ctx, path, beispiel()); err != nil {
		t.Fatal(err)
	}
	if _, err := (recording.YAML{}).Load(ctx, path); err != nil {
		t.Fatal(err)
	}
}

// Abdeckung: LH-FA-07/Boundary — eine neue Aufzeichnung erhält die Rechte 0666
// nach der umask, unter der umask 022 also 0644.
func TestWriteRechteUnterUmask022(t *testing.T) {
	alt := syscall.Umask(0o022)
	defer syscall.Umask(alt)
	path := filepath.Join(t.TempDir(), "rec.yaml")
	if err := (recording.YAML{}).Write(context.Background(), path, beispiel()); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o644 {
		t.Fatalf("neue Datei unter umask 022: %v", info.Mode().Perm())
	}
}
