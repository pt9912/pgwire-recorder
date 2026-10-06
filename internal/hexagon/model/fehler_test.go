package model_test

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// Abdeckung: LH-FA-14/Negative, LH-FA-13/Negative — je Fehlerklasse beginnt der
// Fehlertext mit dem Kopf „<klasse> [<code>]: “, Klassenname nach SPEC-034
// §Ausgabe, und der Exit-Code ist der der Klasse.
func TestFehlerKopfJeKlasse(t *testing.T) {
	for _, f := range []struct {
		code, klasse string
		exit         int
	}{
		{model.CodeInternal, "sonstiger Fehler", 1},
		{model.CodeUsage, "Konfiguration", 2},
		{model.CodeRecordingBroken, "Recording", 3},
		{model.CodeNetwork, "Netzwerk", 4},
		{model.CodeReplayMismatch, "Replay", 5},
		{model.CodeUnsupported, "nicht unterstützt", 6},
	} {
		err := model.Errorf(f.code, nil, "Ursache")
		if got, want := err.Error(), f.klasse+" ["+f.code+"]: Ursache"; got != want {
			t.Errorf("%s: Fehlertext %q, erwartet %q", f.code, got, want)
		}
		if err.ExitCode() != f.exit {
			t.Errorf("%s: Exit-Code %d, erwartet %d", f.code, err.ExitCode(), f.exit)
		}
		ms := model.Meldungen(err)
		if len(ms) != 1 || ms[0].Code != f.code || ms[0].Text != err.Error() || ms[0].ExitCode() != f.exit {
			t.Errorf("%s: Meldungen %#v", f.code, ms)
		}
	}
}

// Abdeckung: LH-FA-14/Negative — eine Fehlerkette trägt einen Kopf, mit Code
// und Klasse des äußersten klassifizierten Fehlers; der innere trägt nur seine
// Ursache bei, auch durch einen nicht klassifizierten Fehler hindurch
// (SPEC-034 §Ausgabe *Fehlerkette*).
func TestFehlerKette(t *testing.T) {
	innen := model.Errorf(model.CodeRecordingBroken, errors.New("Zeile 3"), "Base64-Wert nicht lesbar")
	aussen := model.Errorf(model.CodeRecordingIO, innen, "Session 1, Interaktion 2")
	if got, want := aussen.Error(), "Recording [PGR-E3001]: Session 1, Interaktion 2: Base64-Wert nicht lesbar: Zeile 3"; got != want {
		t.Fatalf("Kette: %q, erwartet %q", got, want)
	}
	gehuellt := fmt.Errorf("Kontext: %w", model.Errorf(model.CodeUpstream, nil, "weg"))
	ms := model.Meldungen(model.Errorf(model.CodeConnectionLost, gehuellt, "aussen"))
	if len(ms) != 1 || ms[0].Text != "Netzwerk [PGR-E4003]: aussen: Kontext: weg" {
		t.Fatalf("Kette durch einen nicht klassifizierten Fehler: %#v", ms)
	}
	ms = model.Meldungen(gehuellt)
	if len(ms) != 1 || ms[0].Code != model.CodeUpstream || ms[0].Text != "Netzwerk [PGR-E4002]: Kontext: weg" {
		t.Fatalf("nicht klassifizierte Hülle um einen klassifizierten: %#v", ms)
	}
	if n := strings.Count(aussen.Error(), "[PGR-"); n != 1 {
		t.Fatalf("%d Köpfe: %q", n, aussen.Error())
	}

	// Eine Hülle mit eigenem Text um mehrere Ursachen ist eine Kette: eine
	// Meldung mit ihrem ganzen Text, Code des ersten klassifizierten Fehlers
	// unter ihr, kein innerer Kopf (SPEC-034 §Ausgabe *Fehlerkette*).
	a := model.Errorf(model.CodeUnsupported, nil, "a")
	b := model.Errorf(model.CodeRecordingIO, nil, "b")
	x := errors.New("x")
	for _, f := range []struct {
		name string
		err  error
		want model.Meldung
	}{
		{"S1", fmt.Errorf("Kontext %w und %w", a, b), model.Meldung{Code: model.CodeUnsupported, Text: "nicht unterstützt [PGR-E6001]: Kontext a und b"}},
		{"S1 fremd zuerst", fmt.Errorf("Kontext %w und %w", x, b), model.Meldung{Code: model.CodeRecordingIO, Text: "Recording [PGR-E3001]: Kontext x und b"}},
		{"S1 ohne Klasse", fmt.Errorf("K %w / %w", x, errors.New("y")), model.Meldung{Code: model.CodeInternal, Text: "sonstiger Fehler [PGR-E1000]: K x / y"}},
		{"S2", model.Errorf(model.CodeUpstream, fmt.Errorf("K %w / %w", a, x), "aussen"), model.Meldung{Code: model.CodeUpstream, Text: "Netzwerk [PGR-E4002]: aussen: K a / x"}},
		{"S3", fmt.Errorf("Kontext: %w", errors.Join(b, a)), model.Meldung{Code: model.CodeRecordingIO, Text: "Recording [PGR-E3001]: Kontext: b; a"}},
		// V-57: Tiefensuche wie errors.As, ein tiefer klassifizierter Fehler in
		// der ersten Ursache geht einem flachen in der zweiten vor.
		{"Tiefensuche", fmt.Errorf("K %w / %w", fmt.Errorf("x: %w", b), model.Errorf(model.CodeUpstream, nil, "u")), model.Meldung{Code: model.CodeRecordingIO, Text: "Recording [PGR-E3001]: K x: b / u"}},
		{"Hülle innen", fmt.Errorf("K %w / %w", fmt.Errorf("i %w", x), model.Errorf(model.CodeUsage, a, "u")), model.Meldung{Code: model.CodeUsage, Text: "Konfiguration [PGR-E2001]: K i x / u: a"}},
	} {
		if got := model.Meldungen(f.err); !reflect.DeepEqual(got, []model.Meldung{f.want}) {
			t.Errorf("%s: %#v, erwartet %#v", f.name, got, f.want)
		}
	}
	// Ohne eigenen Text ist auch eine Hülle mit mehreren %w eine Zusammenfassung.
	if got := model.Meldungen(fmt.Errorf("%w\n%w", a, b)); len(got) != 2 {
		t.Errorf("Hülle ohne eigenen Text: %#v", got)
	}
}

// Abdeckung: LH-FA-14/Negative — ein nicht eingeordneter Fehler ist PGR-E1000
// mit Kopf und Exit-Code 1 (SPEC-034 §Ausgabe *Fremde Fehler*).
func TestFehlerFremd(t *testing.T) {
	ms := model.Meldungen(errors.New("connection reset by peer"))
	if len(ms) != 1 || ms[0].Code != model.CodeInternal || ms[0].Text != "sonstiger Fehler [PGR-E1000]: connection reset by peer" || ms[0].ExitCode() != 1 {
		t.Fatalf("fremder Fehler: %#v", ms)
	}
	if model.Meldungen(nil) != nil {
		t.Fatal("Meldungen(nil) ist nicht nil")
	}
}

// Abdeckung: LH-FA-14/Negative — der Fehlertext ist eine Zeile: LF, CR LF und
// ein einzelnes CR samt folgendem Leerraum (Leerzeichen, Tabulatoren, weitere
// Umbrüche) werden ein Leerzeichen, auch im Text einer Bibliothek; VT, FF,
// U+0085, U+2028 und U+2029 bleiben stehen (SPEC-034 §Ausgabe *Eine Zeile*).
func TestFehlerEineZeile(t *testing.T) {
	for _, f := range []struct{ ursache, want string }{
		{"a\nb", "a b"},
		{"a\r\nb", "a b"},
		{"a\rb", "a b"},
		{"a\n \t\n\r\n  b", "a b"},
		{"a \nb", "a  b"},
		{"a\vb\fc\u0085d e f", "a\vb\fc\u0085d e f"},
	} {
		err := model.Errorf(model.CodeInternal, errors.New(f.ursache), "x")
		if got, want := err.Error(), "sonstiger Fehler [PGR-E1000]: x: "+f.want; got != want {
			t.Errorf("%q: %q, erwartet %q", f.ursache, got, want)
		}
		ms := model.Meldungen(errors.New(f.ursache))
		if got, want := ms[0].Text, "sonstiger Fehler [PGR-E1000]: "+f.want; got != want {
			t.Errorf("fremd %q: %q, erwartet %q", f.ursache, got, want)
		}
	}
	if got := model.Errorf(model.CodeUsage, nil, "Zeile 1\nZeile 2").Error(); got != "Konfiguration [PGR-E2001]: Zeile 1 Zeile 2" {
		t.Fatalf("Umbruch in der eigenen Meldung: %q", got)
	}
}

// Abdeckung: LH-FA-13/Negative, LH-FA-14/Negative — nebeneinander entstandene
// Fehler (errors.Join): jeder klassifizierte ist eine eigene Meldung mit
// eigenem Kopf in der Reihenfolge des Entstehens, die erste ist die gemerkte;
// ein nicht klassifizierter folgt als Ursache der ersten klassifizierten; sind
// alle nicht klassifiziert, ist es eine PGR-E1000 mit den Texten durch "; "
// (SPEC-034 §Ausgabe *Gleichrangige Fehler*).
func TestFehlerGleichrangig(t *testing.T) {
	verbindung := model.Errorf(model.CodeConnectionLost, nil, "Client weg")
	schreiben := model.Errorf(model.CodeRecordingIO, errors.New("disk\nfull"), "r.yaml nicht zu schreiben")
	upstream := errors.New("close tcp: use of closed network connection")

	got := model.Meldungen(errors.Join(verbindung, schreiben, upstream))
	want := []model.Meldung{
		{Code: model.CodeConnectionLost, Text: "Netzwerk [PGR-E4003]: Client weg: close tcp: use of closed network connection"},
		{Code: model.CodeRecordingIO, Text: "Recording [PGR-E3001]: r.yaml nicht zu schreiben: disk full"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Join: %#v", got)
	}
	got = model.Meldungen(errors.Join(nil, upstream, errors.Join(schreiben, verbindung)))
	if len(got) != 2 || got[0].Code != model.CodeRecordingIO || got[1].Code != model.CodeConnectionLost ||
		!strings.HasSuffix(got[0].Text, ": "+upstream.Error()) {
		t.Fatalf("verschachtelter Join, nicht klassifizierter zuerst: %#v", got)
	}
	got = model.Meldungen(errors.Join(errors.New("eins"), errors.New("zwei")))
	if len(got) != 1 || got[0].Code != model.CodeInternal || got[0].Text != "sonstiger Fehler [PGR-E1000]: eins; zwei" {
		t.Fatalf("nur nicht klassifizierte: %#v", got)
	}
	// Ein Join in einer Kette ist eine Ursache, keine eigenen Meldungen.
	got = model.Meldungen(model.Errorf(model.CodeRecordingIO, errors.Join(errors.New("a"), model.Errorf(model.CodeInternal, nil, "b")), "Datei"))
	if len(got) != 1 || got[0].Text != "Recording [PGR-E3001]: Datei: a; b" {
		t.Fatalf("Join in der Kette: %#v", got)
	}
}

// Abdeckung: LH-FA-14/Negative — jeder Code der Tabelle im Quelltext hat die
// Form PGR-[EWI][0-9]{4}; ein Fehlercode ergibt mit der ersten Ziffer seine
// Klasse (1 bis 6) mit Namen, eine Warnung hat einen Bereich (1, 2, 3, 5)
// (SPEC-034).
func TestCodeTabelle(t *testing.T) {
	codes := codesImQuelltext(t)
	if len(codes) < 10 {
		t.Fatalf("nur %d Codes gefunden: %v", len(codes), codes)
	}
	form := regexp.MustCompile(`^PGR-[EWI][0-9]{4}$`)
	gesehen := map[string]string{}
	for name, code := range codes {
		if !form.MatchString(code) {
			t.Errorf("%s = %q hat nicht die Form PGR-[EWI][0-9]{4}", name, code)
			continue
		}
		if anderer, ok := gesehen[code]; ok {
			t.Errorf("%s und %s belegen %s doppelt", name, anderer, code)
		}
		gesehen[code] = name
		ziffer := int(code[5] - '0')
		switch code[4] {
		case 'E':
			if ziffer < 1 || ziffer > 6 || model.Klasse(ziffer) == "" || (&model.Error{Code: code}).ExitCode() != ziffer {
				t.Errorf("%s = %s ergibt keine Klasse", name, code)
			}
		case 'W':
			if !strings.ContainsRune("1235", rune(code[5])) {
				t.Errorf("%s = %s hat keinen Bereich", name, code)
			}
		default:
			t.Errorf("%s = %s: Schwere I ist reserviert", name, code)
		}
	}
}

// codesImQuelltext liest alle Konstanten Code… aus fehler.go.
func codesImQuelltext(t *testing.T) map[string]string {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), "fehler.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	codes := map[string]string{}
	for _, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok || g.Tok != token.CONST {
			continue
		}
		for _, spec := range g.Specs {
			v := spec.(*ast.ValueSpec)
			for i, n := range v.Names {
				if !strings.HasPrefix(n.Name, "Code") {
					continue
				}
				lit, ok := v.Values[i].(*ast.BasicLit)
				if !ok {
					t.Fatalf("%s ist kein Literal", n.Name)
				}
				s, _ := strconv.Unquote(lit.Value)
				codes[n.Name] = s
			}
		}
	}
	return codes
}
