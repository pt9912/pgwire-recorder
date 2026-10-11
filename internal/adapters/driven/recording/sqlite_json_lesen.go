package recording

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"
)

// kanonischeZahl ist der Text einer ganzen Zahl in kanonischer Dezimalschreibweise
// (SPEC-043): 0 oder Ziffern ohne führende Null, bei einer negativen Zahl mit - davor,
// ohne + und ohne Leerraum; -0 ist es nicht.
var kanonischeZahl = regexp.MustCompile(`^(0|-?[1-9][0-9]*)$`)

// jsonKnoten liest den JSON-Text text einer JSON-Spalte als yaml.Node-Baum, den der
// Leser von yaml ebenso liest (SPEC-043): ein Objekt wird zur Abbildung, eine Liste zur
// Liste, ein String zu !!str in Anführungszeichen (ohne gültiges UTF-8 nach der Umkehr
// der Umschreibung zu !!binary), eine Zahl zu !!int, ein Wahrheitswert zu !!bool und null
// zu !!null. Ungültiger JSON-Text, Text ohne gültiges UTF-8, eine Zahl, die nicht kanonisch
// ganz im Bereich der 64-Bit-Zahlen mit Vorzeichen liegt, und ungültiges Base64 sind ein
// Fehler. Die Umwandlung kennt nur JSON-Typen, keine Nachrichtenart; einen doppelten
// Schlüssel, den Schlüssel type in fields und jede Form lehnt der Weg über Unmarshal ab.
func jsonKnoten(text string) (*yaml.Node, error) {
	if !utf8.ValidString(text) || !json.Valid([]byte(text)) {
		return nil, fmt.Errorf("kein gültiger JSON-Text")
	}
	dec := json.NewDecoder(strings.NewReader(text))
	dec.UseNumber()
	return leseJSON(dec)
}

// leseJSON liest den nächsten Wert aus dec; der Text ist als gültig geprüft.
func leseJSON(dec *json.Decoder) (*yaml.Node, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		if t == '{' {
			return leseObjekt(dec)
		}
		return leseListe(dec)
	case string:
		return stringKnoten(t)
	case json.Number:
		return zahlKnoten(t.String())
	case bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(t)}, nil
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, nil
}

func leseObjekt(dec *json.Decoder) (*yaml.Node, error) {
	n := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		schluessel, err := stringKnoten(tok.(string))
		if err != nil {
			return nil, fmt.Errorf("Schlüssel: %w", err)
		}
		wert, err := leseJSON(dec)
		if err != nil {
			return nil, err
		}
		n.Content = append(n.Content, schluessel, wert)
	}
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	return n, nil
}

func leseListe(dec *json.Decoder) (*yaml.Node, error) {
	n := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for dec.More() {
		wert, err := leseJSON(dec)
		if err != nil {
			return nil, err
		}
		n.Content = append(n.Content, wert)
	}
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	return n, nil
}

// zahlKnoten liefert den !!int-Knoten der JSON-Zahl text; ein Bruch, ein Exponent, -0
// und jede Zahl außerhalb der 64-Bit-Zahlen mit Vorzeichen sind ein Fehler.
func zahlKnoten(text string) (*yaml.Node, error) {
	if !kanonischeZahl.MatchString(text) {
		return nil, fmt.Errorf("Zahl %q ist nicht kanonisch ganz", text)
	}
	if _, err := strconv.ParseInt(text, 10, 64); err != nil {
		return nil, fmt.Errorf("Zahl %q außerhalb von 64 Bit", text)
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: text}, nil
}

// stringKnoten kehrt die Umschreibung eines Strings um (SPEC-043): Beginnt s mit dem
// Fluchtzeichen, entfällt dieses eine Zeichen und der Rest gilt wörtlich; beginnt s sonst
// mit praefixBase64, sind die Bytes des Rests dekodiert (kein gültiges Base64 ist ein
// Fehler); jeder andere String gilt wörtlich. Sind die Bytes gültiges UTF-8, ist das
// Ergebnis ein !!str in Anführungszeichen, sonst ein !!binary mit den Bytes in Base64.
func stringKnoten(s string) (*yaml.Node, error) {
	var inhalt string
	switch {
	case strings.HasPrefix(s, fluchtzeichen):
		inhalt = s[len(fluchtzeichen):]
	case strings.HasPrefix(s, praefixBase64):
		roh, err := base64.StdEncoding.DecodeString(s[len(praefixBase64):])
		if err != nil {
			return nil, fmt.Errorf("Base64: %w", err)
		}
		inhalt = string(roh)
	default:
		inhalt = s
	}
	return textKnoten(inhalt), nil
}

// textKnoten liefert den Knoten für den String s: !!str in Anführungszeichen,
// ohne gültiges UTF-8 !!binary mit den Bytes in Base64.
func textKnoten(s string) *yaml.Node {
	if !utf8.ValidString(s) {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!binary", Value: base64.StdEncoding.EncodeToString([]byte(s))}
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Style: yaml.DoubleQuotedStyle, Value: s}
}
