package cli

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// standardDatei ist die Konfigurationsdatei im aktuellen Verzeichnis; sie gilt,
// wenn weder --config noch PGWIRE_RECORDER_CONFIG eine Datei nennt und sie
// existiert (LH-FA-17.a).
const standardDatei = ".pgwire-recorder.yaml"

// envConfig ist die Umgebungsvariable von --config (LH-FA-17.a).
const envConfig = "PGWIRE_RECORDER_CONFIG"

// praefix ist das Präfix der Umgebungsvariablen (SPEC-008).
const praefix = "PGWIRE_RECORDER_"

// datei ist eine geladene und geprüfte Konfigurationsdatei: werte trägt den
// Text jedes Werts je Abschnitt ("" ist die oberste Ebene) und Schlüssel,
// inhalt das Dokument für config show (nil bei einer leeren Datei).
type datei struct {
	werte  map[string]map[string]string
	inhalt *yaml.Node
}

// wert liefert den Text des Schlüssels einer Option von kommando, wenn die
// Datei ihn setzt: log_level auf der obersten Ebene, jeder andere im Abschnitt
// des Kommandos. Eine nil-Datei setzt nichts.
func (d *datei) wert(kommando string, o option) (string, bool) {
	if d == nil {
		return "", false
	}
	abschnitt := kommando
	if o.oben {
		abschnitt = ""
	}
	v, ok := d.werte[abschnitt][schluesselName(o.name)]
	return v, ok
}

// schluesselName ist der Schlüssel einer Option in der Datei: ihr Name mit _
// statt - (LH-FA-17.a).
func schluesselName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

// waehleDatei liefert die eine Konfigurationsdatei nach LH-FA-17.a: die aus
// --config (cli), sonst die aus PGWIRE_RECORDER_CONFIG, sonst standardDatei,
// sofern sie existiert; pfad "" heißt, keine Datei wird gelesen. quelle nennt,
// woher der Pfad kommt, für die Meldung, die den Pfad selbst nicht nennt.
func waehleDatei(cli string) (pfad, quelle string, err error) {
	if cli != "" {
		return cli, "--config", nil
	}
	if v := os.Getenv(envConfig); v != "" {
		return v, envConfig, nil
	}
	_, err = os.Stat(standardDatei)
	switch {
	case err == nil:
		return standardDatei, standardDatei, nil
	case errors.Is(err, fs.ErrNotExist):
		return "", "", nil
	}
	return "", "", nichtLesbar(standardDatei, err)
}

// ladeGewaehlte wählt die Datei (waehleDatei) und lädt sie; ohne gewählte
// Datei liefert es pfad "" und eine nil-Datei, die nichts setzt.
func ladeGewaehlte(cli string) (string, *datei, error) {
	pfad, quelle, err := waehleDatei(cli)
	if err != nil || pfad == "" {
		return "", nil, err
	}
	d, err := ladeDatei(pfad, quelle)
	if err != nil {
		return "", nil, err
	}
	return pfad, d, nil
}

// ladeDatei liest und prüft die Datei pfad ganz, unabhängig vom Kommando
// (LH-FA-17.a): zuerst, ob sie, Links gefolgt, eine reguläre Datei ist, vor
// dem Öffnen; dann Lesen, Kodierung und YAML, dann Schlüssel und Werte in der
// Reihenfolge der Datei. Jeder Fehler ist PGR-E2004 und nennt die Stelle, nie
// einen Wert und nicht den Pfad; eine Meldung zu ungültigem YAML nennt die
// Zeile, außer zu einem Alias ohne Anker (Grenze: die Bibliothek nennt dafür
// keine Stelle).
func ladeDatei(pfad, quelle string) (*datei, error) {
	info, err := os.Stat(pfad)
	if err != nil {
		return nil, nichtLesbar(quelle, err)
	}
	if !info.Mode().IsRegular() {
		return nil, model.Errorf(model.CodeConfigFile, nil, "Konfigurationsdatei (%s) nicht lesbar: keine reguläre Datei", quelle)
	}
	roh, err := os.ReadFile(pfad)
	if err != nil {
		return nil, nichtLesbar(quelle, err)
	}
	data, err := kodierung(roh)
	if err != nil {
		return nil, err
	}
	doc, err := leseYAML(data)
	if err != nil {
		return nil, err
	}
	d := &datei{werte: map[string]map[string]string{"": {}}, inhalt: doc}
	if doc == nil {
		return d, nil
	}
	z := zeilen(data)
	if err := d.pruefe(doc.Content[0], z); err != nil {
		return nil, err
	}
	return d, nil
}

// nichtLesbar ist der Fehler einer Datei, die fehlt oder nicht gelesen werden
// kann; er trägt den Grund des Betriebssystems ohne den Pfad.
func nichtLesbar(quelle string, err error) error {
	if errors.Is(err, fs.ErrNotExist) {
		return model.Errorf(model.CodeConfigFile, nil, "Konfigurationsdatei (%s) nicht vorhanden", quelle)
	}
	var pe *fs.PathError
	if errors.As(err, &pe) {
		err = pe.Err
	}
	return model.Errorf(model.CodeConfigFile, err, "Konfigurationsdatei (%s) nicht lesbar", quelle)
}

// fehlerDatei ist PGR-E2004 an einer Stelle der Datei.
func fehlerDatei(stelle, grund string) error {
	return model.Errorf(model.CodeConfigFile, nil, "Konfigurationsdatei: %s: %s", stelle, grund)
}

// fehlerZeile ist PGR-E2004 für ungültiges YAML in einer Zeile der Datei.
func fehlerZeile(zeile int) error {
	return model.Errorf(model.CodeConfigFile, nil, "Konfigurationsdatei: ungültiges YAML in Zeile %d", zeile)
}

// kodierung prüft den Text der Datei vor dem YAML (LH-FA-17.a): UTF-8, ein BOM
// nur als erstes Zeichen, kein U+0085, U+2028 oder U+2029 und nur Zeichen, die
// YAML 1.2 als druckbar zulässt (nichtDruckbar). Es liefert den Text
// ohne das führende BOM; ein Fehler nennt die Zeile, gezählt nach den
// Zeilenenden \n, \r\n und \r, und kein Zeichen der Datei.
func kodierung(roh []byte) ([]byte, error) {
	data := bytes.TrimPrefix(roh, []byte("\uFEFF"))
	zeile := 1
	for i := 0; i < len(data); {
		r, n := utf8.DecodeRune(data[i:])
		switch {
		case r == utf8.RuneError && n == 1, r == '\uFEFF', r == '\u0085', r == '\u2028', r == '\u2029', nichtDruckbar(r):
			return nil, fehlerZeile(zeile)
		case r == '\r' && i+1 < len(data) && data[i+1] == '\n':
			n = 2
			zeile++
		case r == '\n', r == '\r':
			zeile++
		}
		i += n
	}
	return data, nil
}

// nichtDruckbar meldet die Zeichen, die YAML 1.2 nicht als druckbar zulässt:
// C0 außer Tabulator, \n und \r (auch NUL), U+007F, C1, U+FFFE und U+FFFF.
func nichtDruckbar(r rune) bool {
	switch {
	case r == '\t', r == '\n', r == '\r':
		return false
	case r < 0x20, r >= 0x7F && r <= 0x9F, r == 0xFFFE, r == 0xFFFF:
		return true
	}
	return false
}

// yamlFehler zerlegt den Text eines Fehlers der YAML-Bibliothek in die
// genannte Zeile (leer, wenn er keine nennt) und die Ursache.
var yamlFehler = regexp.MustCompile(`(?s)^yaml: (?:line ([0-9]+): )?(.*)$`)

// aliasOhneAnker ist der Fehler der Bibliothek zu einem Alias ohne Anker, ein
// Fehler aus dem Bezug zwischen Knoten ohne Stelle in der Datei.
var aliasOhneAnker = regexp.MustCompile(`(?s)^unknown anchor '.*' referenced$`)

// parserFehler meldet, ob ursache ein Fehler des Parsers der Bibliothek ist:
// einer der elf Texte aus parserc.go, im genauen Vergleich, weil Texte des
// Scanners ebenso beginnen. Kopplung an go.yaml.in/yaml/v3 v3.0.5: Der Parser
// nennt die Zeile ab 0, der Scanner ab 1 (TestDateiUngueltig hält beides).
func parserFehler(ursache string) bool {
	switch ursache {
	case "did not find expected <stream-start>", "did not find expected <document start>",
		"did not find expected node content", "did not find expected '-' indicator",
		"did not find expected key", "did not find expected ',' or ']'",
		"did not find expected ',' or '}'", "found undefined tag handle",
		"found duplicate %YAML directive", "found incompatible YAML document",
		"found duplicate %TAG directive":
		return true
	}
	return false
}

// ungueltigesYAML ist PGR-E2004 für ungültiges YAML mit der Zeile ab 1
// (LH-FA-17.a *Fehler*): die Zahl im Text der Bibliothek, bei einem Fehler des
// Parsers plus 1; ohne Zahl Zeile 1. Ein Alias ohne Anker nennt die Ursache
// ohne Zeile und ohne den Namen. Vom Text der Bibliothek übernimmt es sonst
// nichts, er kann Inhalt der Datei tragen.
func ungueltigesYAML(err error) error {
	m := yamlFehler.FindStringSubmatch(err.Error())
	if m == nil {
		return fehlerZeile(1)
	}
	ursache := m[2]
	if m[1] == "" {
		if aliasOhneAnker.MatchString(ursache) {
			return model.Errorf(model.CodeConfigFile, nil, "Konfigurationsdatei: ungültiges YAML, Alias ohne Anker")
		}
		return fehlerZeile(1)
	}
	zeile, _ := strconv.Atoi(m[1])
	if parserFehler(ursache) {
		zeile++
	}
	return fehlerZeile(zeile)
}

// leseYAML liest höchstens ein YAML-Dokument; ein zweites ist ungültiges YAML,
// ebenso ein Schlüssel, der in derselben Abbildung zweimal steht. Eine leere
// Datei oder eine nur mit Kommentaren liefert nil (LH-FA-17.a).
func leseYAML(data []byte) (*yaml.Node, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var doc yaml.Node
	err := dec.Decode(&doc)
	if errors.Is(err, io.EOF) {
		return nil, nil
	}
	if err != nil {
		return nil, ungueltigesYAML(err)
	}
	var weiteres yaml.Node
	if err := dec.Decode(&weiteres); !errors.Is(err, io.EOF) {
		if err != nil {
			return nil, ungueltigesYAML(err)
		}
		return nil, model.Errorf(model.CodeConfigFile, nil, "Konfigurationsdatei: ungültiges YAML, mehr als ein Dokument")
	}
	if len(doc.Content) == 0 {
		return nil, nil
	}
	if err := doppelte(doc.Content[0]); err != nil {
		return nil, err
	}
	return &doc, nil
}

// doppelte meldet mit seiner Zeile den ersten Schlüssel, der in einer
// Abbildung zweimal steht, in Abbildungen und Listen jeder Tiefe; verglichen
// wird der Text des Schlüssels. Die Meldung nennt den Schlüssel nicht, er kann
// an der Stelle eines Werts stehen (LH-FA-17.a *Fehler*).
func doppelte(n *yaml.Node) error {
	gesehen := map[string]bool{}
	for i, c := range n.Content {
		istSchluessel := n.Kind == yaml.MappingNode && i%2 == 0
		if istSchluessel && c.Kind == yaml.ScalarNode {
			if gesehen[c.Value] {
				return model.Errorf(model.CodeConfigFile, nil, "Konfigurationsdatei: ungültiges YAML in Zeile %d, Schlüssel doppelt", c.Line)
			}
			gesehen[c.Value] = true
		}
		if err := doppelte(c); err != nil {
			return err
		}
	}
	return nil
}

// stelleVon ist die Stelle eines unbekannten Schlüssels name in der
// Abbildung abschnitt: der Schlüssel unter abschnitt, bei leerem Text (ein
// leerer Schlüssel oder einer, der kein Skalar ist) die Abbildung selbst,
// abbildung (LH-FA-17.a *Fehler*).
func stelleVon(abbildung, abschnitt, name string) string {
	if name == "" {
		return abbildung
	}
	return unter(abschnitt, name)
}

// unter ist die Stelle eines Schlüssels unter stelle, durch einen Punkt
// getrennt.
func unter(stelle, schluessel string) string {
	if stelle == "" {
		return schluessel
	}
	return stelle + "." + schluessel
}

// zeilen ist der Text der Datei ohne führendes BOM (kodierung) je Zeile in
// Zeichen, getrennt nach \n, \r\n und \r.
func zeilen(data []byte) [][]rune {
	text := strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(string(data))
	var out [][]rune
	for _, z := range strings.Split(text, "\n") {
		out = append(out, []rune(z))
	}
	return out
}

// form prüft, was an keinem Knoten der Datei stehen darf: Anker und ein
// ausdrücklich geschriebener Tag (LH-FA-17.a). Einen Tag mit Namen markiert die
// Bibliothek als TaggedStyle; den nicht spezifischen Tag "!" allein markiert
// sie nicht, ihn erkennt form am ersten Zeichen des Knotens in zeilen, an
// Zeile und Spalte der Bibliothek. Ein Alias verweist auf einen Anker davor in
// der Datei, und form lehnt diesen zuerst ab.
func form(n *yaml.Node, stelle string, z [][]rune) error {
	switch {
	case n.Anchor != "":
		return fehlerDatei(stelle, "Anker ist ungültig")
	case n.Style&yaml.TaggedStyle != 0 || beginntMit(n, z, '!'):
		return fehlerDatei(stelle, "Tag ist ungültig")
	}
	return nil
}

// beginntMit meldet, ob der Text der Datei an der Stelle des Knotens mit c
// beginnt.
func beginntMit(n *yaml.Node, z [][]rune, c rune) bool {
	if n.Line < 1 || n.Line > len(z) || n.Column < 1 || n.Column > len(z[n.Line-1]) {
		return false
	}
	return z[n.Line-1][n.Column-1] == c
}

// skalar liefert den Text eines Werts (LH-FA-17.a): ein Skalar, gleich ob mit
// oder ohne Anführungszeichen geschrieben; ungültig sind Liste, Abbildung,
// null und ein leerer Wert ohne Anführungszeichen sowie alles, was form
// ablehnt.
func skalar(n *yaml.Node, stelle string, z [][]rune) (string, error) {
	if err := form(n, stelle, z); err != nil {
		return "", err
	}
	if n.Kind != yaml.ScalarNode {
		return "", fehlerDatei(stelle, "erwartet ein Skalar, keine Liste oder Abbildung")
	}
	if n.Style&(yaml.DoubleQuotedStyle|yaml.SingleQuotedStyle|yaml.LiteralStyle|yaml.FoldedStyle) == 0 && n.Tag == "!!null" {
		return "", fehlerDatei(stelle, "ein leerer Wert und null sind ungültig")
	}
	return n.Value, nil
}

// abbildung prüft, dass ein Knoten eine Abbildung ist: die oberste Ebene, ein
// Abschnitt oder connections:; ohne Inhalt oder null ist ungültig, eine leere
// Abbildung {} setzt nichts (LH-FA-17.a).
func abbildung(n *yaml.Node, stelle string, z [][]rune) error {
	if err := form(n, stelle, z); err != nil {
		return err
	}
	if n.Kind != yaml.MappingNode {
		return fehlerDatei(stelle, "erwartet eine Abbildung")
	}
	return nil
}

// schluessel liefert den Text eines Schlüssels, den form zulässt; ein
// Schlüssel, der kein Skalar ist, hat den leeren Text, und keine Option heißt
// so.
func schluessel(k *yaml.Node, stelle string, z [][]rune) (string, error) {
	if err := form(k, stelle, z); err != nil {
		return "", err
	}
	return k.Value, nil
}

// pruefe prüft die oberste Ebene in der Reihenfolge der Datei: log_level,
// connections und je Kommando am allgemeinen Leser ein Abschnitt; jeder andere
// Schlüssel ist unbekannt (LH-FA-17.a).
func (d *datei) pruefe(top *yaml.Node, z [][]rune) error {
	if err := abbildung(top, "oberste Ebene", z); err != nil {
		return err
	}
	oben := obenOptionen()
	for i := 0; i+1 < len(top.Content); i += 2 {
		name, err := schluessel(top.Content[i], "oberste Ebene", z)
		if err != nil {
			return err
		}
		v := top.Content[i+1]
		switch o, ist := oben[name]; {
		case ist:
			err = d.setze("", name, o, v, z)
		case name == "connections":
			err = verbindungen(v, z)
		case istLeserKommando(name):
			err = d.abschnitt(name, v, z)
		default:
			err = fehlerDatei(stelleVon("oberste Ebene", "", name), "unbekannter Schlüssel")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// abschnitt prüft den Abschnitt eines Kommandos: jeder Schlüssel ist der einer
// Option des Kommandos, die nicht auf der obersten Ebene steht.
func (d *datei) abschnitt(kommando string, n *yaml.Node, z [][]rune) error {
	if err := abbildung(n, kommando, z); err != nil {
		return err
	}
	d.werte[kommando] = map[string]string{}
	opts := map[string]option{}
	for _, o := range optionen(kommando) {
		if !o.oben {
			opts[schluesselName(o.name)] = o
		}
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		name, err := schluessel(n.Content[i], kommando, z)
		if err != nil {
			return err
		}
		o, ok := opts[name]
		if !ok {
			return fehlerDatei(stelleVon(kommando, kommando, name), "unbekannter Schlüssel")
		}
		if err := d.setze(kommando, name, o, n.Content[i+1], z); err != nil {
			return err
		}
	}
	return nil
}

// setze prüft den Wert eines Schlüssels mit der Wertemenge seiner Option und
// merkt ihn.
func (d *datei) setze(abschnitt, name string, o option, n *yaml.Node, z [][]rune) error {
	stelle := unter(abschnitt, name)
	text, err := skalar(n, stelle, z)
	if err != nil {
		return err
	}
	if err := o.art.pruefe(text); err != nil {
		return model.Errorf(model.CodeConfigFile, err, "Konfigurationsdatei: %s", stelle)
	}
	d.werte[abschnitt][name] = text
	return nil
}

// verbindungen prüft connections: eine Abbildung von Namen auf Werte; ein
// Name hat die Form eines Werts und enthält kein Steuerzeichen (LH-FA-17.a).
func verbindungen(n *yaml.Node, z [][]rune) error {
	if err := abbildung(n, "connections", z); err != nil {
		return err
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		name, err := skalar(n.Content[i], "connections", z)
		if err != nil {
			return err
		}
		if name == "" {
			return fehlerDatei("connections", "Name einer Verbindung ist leer")
		}
		if strings.IndexFunc(name, unicode.IsControl) >= 0 {
			return fehlerDatei("connections", "Name einer Verbindung mit Steuerzeichen")
		}
		if _, err := skalar(n.Content[i+1], unter("connections", name), z); err != nil {
			return err
		}
	}
	return nil
}

// obenOptionen sind die Optionen am allgemeinen Leser, deren Schlüssel auf der
// obersten Ebene der Datei steht, nach Schlüssel.
func obenOptionen() map[string]option {
	out := map[string]option{}
	for _, k := range leserKommandos() {
		for _, o := range optionen(k) {
			if o.oben {
				out[schluesselName(o.name)] = o
			}
		}
	}
	return out
}

// anzeige ist die Ausgabe von config show (LH-FA-17.a *Anzeige*): in der
// ersten Zeile der Pfad der gewählten Datei oder dass keine gefunden wurde,
// danach ihr Inhalt als YAML mit zwei Leerzeichen Einzug in der Reihenfolge
// der Datei, ohne Kommentare, danach die Namen der aktiven
// Umgebungsvariablen mit dem Präfix, nach Namen sortiert.
func anzeige(pfad string, d *datei) (string, error) {
	var b strings.Builder
	if pfad == "" {
		b.WriteString("keine Konfigurationsdatei gefunden\n")
	} else {
		b.WriteString(pfad + "\n")
	}
	if d != nil && d.inhalt != nil {
		ohneKommentare(d.inhalt)
		enc := yaml.NewEncoder(&b)
		enc.SetIndent(2)
		if err := enc.Encode(d.inhalt); err != nil {
			return "", model.Errorf(model.CodeInternal, err, "Konfigurationsdatei anzeigen")
		}
		if err := enc.Close(); err != nil {
			return "", model.Errorf(model.CodeInternal, err, "Konfigurationsdatei anzeigen")
		}
	}
	for _, name := range aktiveUmgebung() {
		b.WriteString(name + "\n")
	}
	return b.String(), nil
}

// ohneKommentare entfernt die Kommentare aus einem Knoten und allen darunter.
func ohneKommentare(n *yaml.Node) {
	n.HeadComment, n.LineComment, n.FootComment = "", "", ""
	for _, c := range n.Content {
		ohneKommentare(c)
	}
}

// aktiveUmgebung sind die Namen der gesetzten und nicht leeren
// Umgebungsvariablen mit dem Präfix PGWIRE_RECORDER_, auch ohne passende
// Option, nach Namen sortiert.
func aktiveUmgebung() []string {
	var out []string
	for _, e := range os.Environ() {
		name, wert, _ := strings.Cut(e, "=")
		if strings.HasPrefix(name, praefix) && wert != "" {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}
