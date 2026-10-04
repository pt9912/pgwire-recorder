package cli

import (
	"bytes"
	"errors"
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

func TestParseRecord(t *testing.T) {
	cmd, err := Parse([]string{"record", "--listen", ":15432", "--upstream", "pg:5432", "--output", "r.yaml", "--force"}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	want := RecordOptions{Listen: ":15432", Upstream: "pg:5432", Output: "r.yaml", Force: true}
	if cmd.Name != "record" || cmd.Record != want {
		t.Fatalf("erhalten %#v", cmd)
	}
}

func TestParseFehler(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"unbekannt"},
		{"record", "--listen", ":1", "--upstream", "pg:5432"},
		{"record", "--listen", ":1", "--upstream", "pg:5432", "--output", "r.yaml", "zusatz"},
		{"record", "--gibtsnicht"},
	} {
		_, err := Parse(args, &bytes.Buffer{})
		var me *model.Error
		if !errors.As(err, &me) || me.Code != model.CodeUsage || me.ExitCode() != 2 {
			t.Errorf("%v: erwartet %s mit Exit-Code 2, erhalten %v", args, model.CodeUsage, err)
		}
	}
}

func TestParseHilfe(t *testing.T) {
	var out bytes.Buffer
	_, err := Parse([]string{"--help"}, &out)
	if !errors.Is(err, ErrHelp) || out.Len() == 0 {
		t.Fatalf("Hilfe: err=%v ausgabe=%q", err, out.String())
	}
}
