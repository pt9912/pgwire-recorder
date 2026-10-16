package recording

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Präfixe der Umschreibung eines Strings in den JSON-Spalten (SPEC-043): Ein String
// ohne gültiges UTF-8 steht als praefixBase64 und seine Bytes in Base64; ein gültiger
// String, der mit praefixBase64 oder mit dem Fluchtzeichen beginnt, steht mit einem
// weiteren Fluchtzeichen davor.
const (
	praefixBase64 = "base64:"
	fluchtzeichen = `\`
)

// jsonText liefert den JSON-Text von wert für eine JSON-Spalte (SPEC-043): wert wird
// wie für die YAML-Datei zu einem yaml.Node, und dessen Baum wird zu JSON. Die
// Schlüssel einer Abbildung stehen in der Reihenfolge des Baums. ohne nennt einen
// Schlüssel der obersten Abbildung, der entfällt; "" lässt keinen entfallen. Die
// Umwandlung kennt nur Abbildung, Liste, String, Zahl, Wahrheitswert und null.
func jsonText(wert any, ohne string) (string, error) {
	var n yaml.Node
	if err := n.Encode(wert); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := schreibeJSON(&buf, &n, ohne); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// schreibeJSON schreibt den Baum n als JSON nach buf; ohne gilt nur für n selbst.
func schreibeJSON(buf *bytes.Buffer, n *yaml.Node, ohne string) error {
	switch n.Kind {
	case yaml.MappingNode:
		return schreibeObjekt(buf, n, ohne)
	case yaml.SequenceNode:
		buf.WriteByte('[')
		for i, c := range n.Content {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := schreibeJSON(buf, c, ""); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
		return nil
	case yaml.ScalarNode:
		return schreibeSkalar(buf, n)
	}
	return fmt.Errorf("Knoten der Art %d lässt sich nicht als JSON schreiben", n.Kind)
}

func schreibeObjekt(buf *bytes.Buffer, n *yaml.Node, ohne string) error {
	buf.WriteByte('{')
	erster := true
	for i := 0; i+1 < len(n.Content); i += 2 {
		schluessel, wert := n.Content[i], n.Content[i+1]
		if ohne != "" && schluessel.Value == ohne {
			continue
		}
		if !erster {
			buf.WriteByte(',')
		}
		erster = false
		text, err := zeichenkette(schluessel)
		if err != nil {
			return fmt.Errorf("Schlüssel: %w", err)
		}
		schreibeString(buf, text)
		buf.WriteByte(':')
		if err := schreibeJSON(buf, wert, ""); err != nil {
			return err
		}
	}
	buf.WriteByte('}')
	return nil
}

func schreibeSkalar(buf *bytes.Buffer, n *yaml.Node) error {
	switch n.Tag {
	case "!!str", "!!binary":
		text, err := zeichenkette(n)
		if err != nil {
			return err
		}
		schreibeString(buf, text)
	case "!!int":
		if !json.Valid([]byte(n.Value)) {
			return fmt.Errorf("Zahl %q lässt sich nicht als JSON schreiben", n.Value)
		}
		buf.WriteString(n.Value)
	case "!!bool":
		buf.WriteString(n.Value)
	case "!!null":
		buf.WriteString("null")
	default:
		return fmt.Errorf("Wert vom Tag %s lässt sich nicht als JSON schreiben", n.Tag)
	}
	return nil
}

// zeichenkette liefert den String des Skalars n in der Umschreibung der JSON-Spalten
// (SPEC-043): ein !!binary-Skalar steht als praefixBase64 und Base64, ein !!str-Skalar,
// der mit praefixBase64 oder dem Fluchtzeichen beginnt, mit einem weiteren
// Fluchtzeichen davor. Ein Skalar ohne Tag !!str, !!binary oder !!null ist ein
// Fehler; null steht als sein Text.
func zeichenkette(n *yaml.Node) (string, error) {
	if n.Kind != yaml.ScalarNode {
		return "", fmt.Errorf("Skalar erwartet")
	}
	switch n.Tag {
	case "!!binary":
		// Die Bibliothek bricht lange Base64-Texte um; DecodeString überspringt die
		// Zeilenumbrüche.
		roh, err := base64.StdEncoding.DecodeString(n.Value)
		if err != nil {
			return "", err
		}
		return praefixBase64 + base64.StdEncoding.EncodeToString(roh), nil
	case "!!str", "!!null":
		return maskiere(n.Value), nil
	}
	return "", fmt.Errorf("Tag %s ist kein String", n.Tag)
}

// maskiere liefert den gültigen String s als Text einer JSON-Spalte: beginnt er mit
// praefixBase64 oder dem Fluchtzeichen, mit einem weiteren Fluchtzeichen davor, sonst
// unverändert. Die Bibliothek schreibt einen String ohne gültiges UTF-8 als
// !!binary-Skalar, den zeichenkette umschreibt; ein !!str-Skalar ist gültiges UTF-8.
func maskiere(s string) string {
	if strings.HasPrefix(s, praefixBase64) || strings.HasPrefix(s, fluchtzeichen) {
		return fluchtzeichen + s
	}
	return s
}

// schreibeString schreibt s als JSON-String nach buf, ohne <, > und & zu maskieren.
func schreibeString(buf *bytes.Buffer, s string) {
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	// Encode hängt ein Zeilenende an; es gehört nicht zum JSON-Text. s ist gültiges
	// UTF-8 und ein String, den Encode nie ablehnt.
	_ = enc.Encode(s)
	buf.Truncate(buf.Len() - 1)
}
