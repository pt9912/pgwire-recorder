package recording

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
)

// sqliteKopf sind die ersten 16 Bytes einer SQLite-Datei: der Text SQLite format 3 und ein
// NUL-Byte (LH-FA-22.a *Wahl und Erkennung*).
const sqliteKopf = "SQLite format 3\x00"

// Leser liest Aufzeichnungen beider Formate: Load erkennt das Format der Datei an ihren
// ersten 16 Bytes und liest sie mit SQLite.Load oder YAML.Load. Prepare und Write gehören
// nicht zu ihm (PGR-E1000); geschrieben wird über New.
type Leser struct{}

var _ driven.RecordingRepository = Leser{}

// Prepare ist ein Programmierfehler des Aufrufers: PGR-E1000 ohne Pfad.
func (Leser) Prepare(context.Context, string, bool) error {
	return model.Errorf(model.CodeInternal, nil, "Leser bereitet kein Schreiben vor")
}

// Write ist ein Programmierfehler des Aufrufers: PGR-E1000 ohne Pfad.
func (Leser) Write(context.Context, string, model.Recording) error {
	return model.Errorf(model.CodeInternal, nil, "Leser schreibt nicht")
}

// Load erkennt das Format der Datei path und liest sie (LH-FA-22.a *Wahl und Erkennung*):
// Beginnt sie mit sqliteKopf, liest SQLite.Load sie, jede andere Datei, auch eine leere
// oder kürzere, liest YAML.Load. Scheitert das Öffnen oder das Lesen der ersten 16 Bytes,
// ist das PGR-E3001 mit der Ursache. Der Pfad gilt wörtlich.
func (Leser) Load(ctx context.Context, path string) (model.Recording, error) {
	f, err := os.Open(path)
	if err != nil {
		return model.Recording{}, model.Errorf(model.CodeRecordingIO, err, "%s nicht lesbar", path)
	}
	kopf := make([]byte, len(sqliteKopf))
	n, err := io.ReadFull(f, kopf)
	err = errors.Join(ignoriereKurz(err), f.Close())
	if err != nil {
		return model.Recording{}, model.Errorf(model.CodeRecordingIO, err, "%s nicht lesbar", path)
	}
	if bytes.Equal(kopf[:n], []byte(sqliteKopf)) {
		return (&SQLite{}).Load(ctx, path)
	}
	return YAML{}.Load(ctx, path)
}

// ignoriereKurz liefert nil für eine Datei mit weniger als 16 Bytes, sonst err.
func ignoriereKurz(err error) error {
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return nil
	}
	return err
}
