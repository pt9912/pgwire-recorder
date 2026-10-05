//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ruhe liegt über der Ruhezeit (1 s), nach der pgxpool und database/sql mit pgx
// eine Verbindung vor der nächsten Nutzung mit `-- ping` prüfen.
const ruhe = 1500 * time.Millisecond

// poolAblauf führt über listen mit pgxpool in der Default-Konfiguration, auf
// eine Verbindung begrenzt, dreimal dieselbe Anfrage mit verschiedenen Werten
// aus, zwischen zwei Anfragen pause lang ruhend, und liefert die Ergebnisse als
// Text.
func poolAblauf(t *testing.T, listen string, pause time.Duration) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(dsn(listen))
	if err != nil {
		t.Fatal(err)
	}
	cfg.MaxConns = 1
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("Pool über %s: %v", listen, err)
	}
	defer pool.Close()
	var b strings.Builder
	for i := 1; i <= 3; i++ {
		if i > 1 {
			time.Sleep(pause)
		}
		var n int
		err := pool.QueryRow(ctx, "SELECT $1::int + 1", i).Scan(&n)
		fmt.Fprintf(&b, "%d: %d %v\n", i, n, err)
	}
	return b.String()
}

// sqlAblauf ist poolAblauf mit database/sql über den Treiber pgx (stdlib) in
// der Default-Konfiguration, auf eine Verbindung begrenzt.
func sqlAblauf(t *testing.T, listen string, pause time.Duration) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, err := sql.Open("pgx", dsn(listen))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	var b strings.Builder
	for i := 1; i <= 3; i++ {
		if i > 1 {
			time.Sleep(pause)
		}
		var n int
		err := db.QueryRowContext(ctx, "SELECT $1::int + 1", i).Scan(&n)
		fmt.Fprintf(&b, "%d: %d %v\n", i, n, err)
	}
	return b.String()
}

// lebendLagen zeichnet ablauf einmal ohne und einmal mit Ruhepausen auf und
// spielt jede Aufzeichnung ohne und mit Ruhepausen wieder ab (Lagen S4 bis S7
// der Validierung des Extended-Replay): Jede Wiedergabe liefert die Sicht beim
// Aufzeichnen, endet mit Exit-Code 0 und meldet keine Warnung. Dass die
// Ruhepausen Lebendprüfungen auslösen, zeigt die Aufzeichnung mit Pausen: Sie
// enthält mehr als die ohne Pausen (database/sql prüft auch ohne Pause bei der
// ersten Wiederverwendung einer Verbindung).
func lebendLagen(t *testing.T, ablauf func(*testing.T, string, time.Duration) string) {
	t.Helper()
	dir := t.TempDir()
	ohne, mit := filepath.Join(dir, "ohne.yaml"), filepath.Join(dir, "mit.yaml")
	rec := startRecorder(t, os.Getenv("PGR_UPSTREAM"), ohne)
	sicht := ablauf(t, rec.listen, 0)
	rec.stop(t, 0)
	rec = startRecorder(t, os.Getenv("PGR_UPSTREAM"), mit)
	if got := ablauf(t, rec.listen, ruhe); got != sicht {
		t.Fatalf("Aufzeichnen mit Pausen weicht ab:\n%s\n--- ohne Pausen:\n%s", got, sicht)
	}
	rec.stop(t, 0)
	if !strings.Contains(sicht, "3: 4 <nil>") {
		t.Fatalf("Sicht beim Aufzeichnen:\n%s", sicht)
	}
	if nOhne, nMit := strings.Count(lies(t, ohne), "-- ping"), strings.Count(lies(t, mit), "-- ping"); nMit <= nOhne {
		t.Fatalf("Aufzeichnung mit Pausen enthält %d Lebendprüfungen, ohne Pausen %d; erwartet mehr", nMit, nOhne)
	}
	for _, lage := range []struct {
		name  string
		input string
		pause time.Duration
	}{
		{"ohne Pausen aufgezeichnet und wiedergegeben", ohne, 0},
		{"Pausen beim Wiedergeben", ohne, ruhe},
		{"Pausen beim Aufzeichnen", mit, 0},
		{"Pausen bei beidem", mit, ruhe},
	} {
		t.Run(lage.name, func(t *testing.T) {
			rep := startProzess(t, "replay", "--input", lage.input)
			got := ablauf(t, rep.listen, lage.pause)
			rep.stop(t, 0)
			if got != sicht {
				t.Fatalf("%s\n--- aufgezeichnet:\n%s\n--- stderr:\n%s", got, sicht, rep.stderr.String())
			}
			if strings.Contains(rep.stderr.String(), "PGR-W") {
				t.Fatalf("Warnung\n%s", rep.stderr.String())
			}
		})
	}
}

// Abdeckung: LH-FA-18/Happy, LH-FA-09/Negative, LH-QA-02/Messung — pgxpool in der
// Default-Konfiguration, am Treiber nur Host und Port umgestellt, läuft im
// Replay ohne Abweichung und ohne Warnung, gleich ob beim Aufzeichnen, beim
// Wiedergeben oder bei beidem Pausen über 1 s Lebendprüfungen auslösen.
func TestE2EReplayLebendpruefungPgxpool(t *testing.T) {
	lebendLagen(t, poolAblauf)
}

// Abdeckung: LH-FA-18/Happy, LH-FA-09/Negative, LH-QA-02/Messung — database/sql
// mit pgx in der Default-Konfiguration, am Treiber nur Host und Port
// umgestellt, läuft im Replay ohne Abweichung und ohne Warnung, gleich ob beim
// Aufzeichnen, beim Wiedergeben oder bei beidem Pausen über 1 s
// Lebendprüfungen auslösen.
func TestE2EReplayLebendpruefungDatabaseSQL(t *testing.T) {
	lebendLagen(t, sqlAblauf)
}

// lebendTexte sind Lebendprüfungen nach LH-FA-09.a gegen jede Serverversion:
// Leerraum ohne \v, Zeilenkommentare bis Zeilen- oder Textende, auch mit \v
// darin, verschachtelte Blockkommentare.
var lebendTexte = []string{
	"", " ", "\t\n\r\f", "-- ping", "--", "-- a\r-- b\n", "-- a\fb", "-- a\vb", "-- a /* b",
	"/**/", "/* a /* b */ c */", "/* -- */", "/* a */ -- b\n",
}

// vtTexte sind Lebendprüfungen nach LH-FA-09.a nur, wenn \v nach der
// Serverversion Leerraum ist.
var vtTexte = []string{"\v", "\t\n\r\f\v", "\v-- ping\v", "/* a */ -- b\n\v"}

// vtAbVersion ist die Hauptversion, ab der \v nach LH-FA-09.a Leerraum ist.
const vtAbVersion = 17

// hauptversion liest die führenden Ziffern von server_version der Verbindung.
func hauptversion(t *testing.T, conn *pgconn.PgConn) int {
	t.Helper()
	v := conn.ParameterStatus("server_version")
	n := 0
	for n < len(v) && v[n] >= '0' && v[n] <= '9' {
		n++
	}
	major, err := strconv.Atoi(v[:n])
	if err != nil {
		t.Fatalf("server_version %q: %v", v, err)
	}
	return major
}

// antwortArt beschreibt die Antwort auf die einfache Anfrage text: Fehler mit
// SQLSTATE oder je Ergebnis Befehl und Spaltenzahl, dazu den
// Transaktionsstatus danach.
func antwortArt(t *testing.T, conn *pgconn.PgConn, text string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := conn.Exec(ctx, text).ReadAll()
	var b strings.Builder
	if err != nil {
		fmt.Fprintf(&b, "fehler %s", pgFehler(err))
	}
	for _, r := range res {
		fmt.Fprintf(&b, "[%q %d %d]", r.CommandTag.String(), len(r.FieldDescriptions), len(r.Rows))
	}
	fmt.Fprintf(&b, " tx %c", conn.TxStatus())
	return b.String()
}

// lebendAblauf führt über listen eine Transaktion mit einem Fehler aus und
// sendet, wenn mitPruefungen, vor, zwischen und nach ihren Anweisungen jede
// Lebendprüfung aus lebendTexte, nach der ersten Anweisung mit vt auch die aus
// vtTexte; es liefert die Antworten als Text. Vor der ersten Anweisung ist der
// Verbindung im Replay keine Session zugeordnet, und \v ist dort kein
// Leerraum (LH-FA-09.a).
func lebendAblauf(t *testing.T, listen string, mitPruefungen, vt bool) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgconn.Connect(ctx, dsn(listen))
	if err != nil {
		t.Fatalf("Verbindung über %s: %v", listen, err)
	}
	defer conn.Close(ctx)
	var b strings.Builder
	pruefen := func(texte []string) {
		for _, text := range texte {
			if mitPruefungen {
				fmt.Fprintf(&b, "  %q: %s\n", text, antwortArt(t, conn, text))
			}
		}
	}
	nachZuordnung := lebendTexte
	if vt {
		nachZuordnung = append(append([]string{}, lebendTexte...), vtTexte...)
	}
	pruefen(lebendTexte)
	for _, anweisung := range []string{"BEGIN", "SELECT 1/0", "ROLLBACK"} {
		fmt.Fprintf(&b, "%s: %s\n", anweisung, antwortArt(t, conn, anweisung))
		pruefen(nachZuordnung)
	}
	return b.String()
}

// Abdeckung: LH-FA-09/Negative — jede Lebendprüfung nach LH-FA-09.a erhält im
// Replay dieselbe Antwort wie von PostgreSQL, auch in einer Transaktion und
// nach einem Fehler darin, ohne in der Aufzeichnung zu stehen; PostgreSQL
// beantwortet jede wie die leere Anfrage. Texte mit \v als Leerraum beantwortet
// PostgreSQL ab Hauptversion 17 wie die leere Anfrage, davor mit einem
// Syntaxfehler; das Replay folgt der Version in der Aufzeichnung: Mit ihnen
// aufgezeichnet, liefert es ab 17 die leere Antwort außerhalb der Reihe und
// davor den aufgezeichneten Fehler samt Transaktionsstatus. Ein offener Blockkommentar, Leerraum
// außerhalb von ASCII und eine Anweisung nach `\r` sind auch für PostgreSQL
// keine leere Anfrage; ein einzelnes `;` beantwortet PostgreSQL wie die leere
// Anfrage, das Replay aber als Abweichung (PGR-E5001, Exit-Code 5).
func TestE2EReplayLebendpruefungWiePostgres(t *testing.T) {
	upstream := os.Getenv("PGR_UPSTREAM")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pg, err := pgconn.Connect(ctx, dsn(upstream))
	if err != nil {
		t.Fatal(err)
	}
	leer := antwortArt(t, pg, "")
	for _, text := range lebendTexte {
		if got := antwortArt(t, pg, text); got != leer {
			t.Errorf("PostgreSQL beantwortet die Lebendprüfung %q mit %s, die leere Anfrage mit %s", text, got, leer)
		}
	}
	vt := hauptversion(t, pg) >= vtAbVersion
	for _, text := range vtTexte {
		if got := antwortArt(t, pg, text); (got == leer) != vt {
			t.Errorf("PostgreSQL %s beantwortet %q mit %s, die leere Anfrage mit %s; \\v Leerraum laut LH-FA-09.a: %v",
				pg.ParameterStatus("server_version"), text, got, leer, vt)
		}
	}
	for _, text := range []string{"/* offen", "/* a /* b */", "\u00a0", "-- a\rSELECT 1"} {
		if got := antwortArt(t, pg, text); got == leer {
			t.Errorf("PostgreSQL beantwortet %q wie die leere Anfrage: %s", text, got)
		}
	}
	if got := antwortArt(t, pg, ";"); got != leer {
		t.Errorf("PostgreSQL beantwortet `;` mit %s, die leere Anfrage mit %s", got, leer)
	}
	_ = pg.Close(ctx)
	if t.Failed() {
		t.FailNow()
	}

	sichtPG := lebendAblauf(t, upstream, true, vt)
	input := filepath.Join(t.TempDir(), "rec.yaml")
	rec := startRecorder(t, upstream, input)
	lebendAblauf(t, rec.listen, false, vt)
	rec.stop(t, 0)
	rep := startProzess(t, "replay", "--input", input)
	got := lebendAblauf(t, rep.listen, true, vt)
	rep.stop(t, 0)
	if got != sichtPG {
		t.Fatalf("Replay:\n%s\n--- PostgreSQL:\n%s\n--- stderr:\n%s", got, sichtPG, rep.stderr.String())
	}

	// Die Texte mit \v nach der Zuordnung gegen jede Version, aufgezeichnet:
	// unter 17 sind sie Interaktionen mit Syntaxfehler, ab 17 Lebendprüfungen.
	sichtVT := lebendAblauf(t, upstream, true, true)
	inputVT := filepath.Join(t.TempDir(), "rec-vt.yaml")
	rec = startRecorder(t, upstream, inputVT)
	if got := lebendAblauf(t, rec.listen, true, true); got != sichtVT {
		t.Fatalf("Aufzeichnen mit \\v:\n%s\n--- PostgreSQL:\n%s", got, sichtVT)
	}
	rec.stop(t, 0)
	if fehler := strings.Contains(sichtVT, "42601"); fehler == vt {
		t.Fatalf("Syntaxfehler für \\v bei \\v als Leerraum %v:\n%s", vt, sichtVT)
	}
	rep = startProzess(t, "replay", "--input", inputVT)
	got = lebendAblauf(t, rep.listen, true, true)
	rep.stop(t, 0)
	if got != sichtVT {
		t.Fatalf("Replay mit \\v:\n%s\n--- PostgreSQL:\n%s\n--- stderr:\n%s", got, sichtVT, rep.stderr.String())
	}

	rep = startProzess(t, "replay", "--input", input)
	conn, err := pgconn.Connect(ctx, dsn(rep.listen))
	if err != nil {
		t.Fatal(err)
	}
	if got := antwortArt(t, conn, ";"); !strings.Contains(got, "PGR-E5001") {
		t.Fatalf("`;` im Replay: %s", got)
	}
	_ = conn.Close(ctx)
	rep.stop(t, 5)
}
