package services

import (
	"bytes"

	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
)

// abweichung ist der Strict Matcher für Extended-Nachrichten (LH-FA-18.a): Er
// vergleicht die empfangene Client-Nachricht in allen Feldern mit der
// erwarteten und liefert den Namen des ersten abweichenden Felds oder "" bei
// Gleichheit. Namen werden nicht normalisiert. Eine leere Liste und nil sind
// gleich; Parameterwerte vergleicht er bytegenau, und NULL ist von einem
// leeren Wert verschieden.
func abweichung(empfangen, erwartet model.ClientMessage) string {
	switch {
	case empfangen.Type != erwartet.Type:
		return "type"
	case empfangen.Statement != erwartet.Statement:
		return "statement"
	case empfangen.Portal != erwartet.Portal:
		return "portal"
	case empfangen.SQL != erwartet.SQL:
		return "sql"
	case !gleicheListe(empfangen.ParamTypes, erwartet.ParamTypes):
		return "param_types"
	case !gleicheListe(empfangen.ParamFormats, erwartet.ParamFormats):
		return "param_formats"
	case !gleicheWerte(empfangen.Params, erwartet.Params):
		return "params"
	case !gleicheListe(empfangen.ResultFormats, erwartet.ResultFormats):
		return "result_formats"
	case empfangen.Target != erwartet.Target:
		return "target"
	case empfangen.Name != erwartet.Name:
		return "name"
	case empfangen.MaxRows != erwartet.MaxRows:
		return "max_rows"
	}
	return ""
}

func gleicheListe[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func gleicheWerte(a, b []model.Value) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Null != b[i].Null || !bytes.Equal(a[i].Bytes, b[i].Bytes) {
			return false
		}
	}
	return true
}
