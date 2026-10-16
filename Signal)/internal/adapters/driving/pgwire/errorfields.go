package pgwire

import (
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgproto3"
)

// bekannteFelder sind die Feldcodes, die pgproto3.ErrorResponse benennt; alle
// anderen gehen in UnknownFields.
const bekannteFelder = "SVCMDHPpqWstcdnFLR"

// errorResponse bildet die Felder einer Fehler- oder Hinweisantwort
// (Feldcode → Wert) auf die PGWire-Nachricht ab.
func errorResponse(f map[string]string) pgproto3.ErrorResponse {
	num := func(code string) int32 {
		n, _ := strconv.ParseInt(f[code], 10, 32)
		return int32(n)
	}
	e := pgproto3.ErrorResponse{
		Severity:            f["S"],
		SeverityUnlocalized: f["V"],
		Code:                f["C"],
		Message:             f["M"],
		Detail:              f["D"],
		Hint:                f["H"],
		Position:            num("P"),
		InternalPosition:    num("p"),
		InternalQuery:       f["q"],
		Where:               f["W"],
		SchemaName:          f["s"],
		TableName:           f["t"],
		ColumnName:          f["c"],
		DataTypeName:        f["d"],
		ConstraintName:      f["n"],
		File:                f["F"],
		Line:                num("L"),
		Routine:             f["R"],
	}
	for k, v := range f {
		if len(k) == 1 && !strings.Contains(bekannteFelder, k) {
			if e.UnknownFields == nil {
				e.UnknownFields = map[byte]string{}
			}
			e.UnknownFields[k[0]] = v
		}
	}
	return e
}
