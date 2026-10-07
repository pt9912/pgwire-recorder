#!/usr/bin/env bash
# lint-gegenprobe — prüft das Lint-Gate (Stufe `lint` des Dockerfile mit
# tools/harness/lint.sh und .golangci.yml) gegen seinen Vertrag
# (spec/spezifikation.md SPEC-049) an Mutanten in Kopien des Arbeitsbaums.
# Aufruf: über `make lint-gegenprobe` (harness/mk/lint.mk), das DOCKER_BUILD
# übergibt; auf dem Host mit bash, make, docker, grep, sed, awk und coreutils.
#
# Je Kopie ein Verzeichnis unter `mktemp -d`, kopiert mit `cp -R` ohne `-p` aus dem
# Arbeitsbaum ohne .git; jede mutierte Datei bekommt danach ein `touch`, damit der
# Build-Kontext sie neu überträgt (Grenze von SPEC-049). Je Kopie ein Lauf von
# `docker build --target lint`.
#
# LINT_GEGENPROBE_PARALLEL — Höchstzahl gleichzeitiger Läufe der Stufe. Nicht
# gesetzt: 6. Gesetzt: eine positive ganze Zahl ohne Vorzeichen, Leerraum und
# führende Null (`^[1-9][0-9]*$`). Jeder andere Wert, auch der leere, bricht vor der
# ersten Kopie ab: Ausgang 2 und auf stderr
# `lint-gegenprobe: LINT_GEGENPROBE_PARALLEL ist keine positive ganze Zahl: '<wert>'`.
# GEPRÜFT DURCH einstellung-<wert> (0, leer, -1, abc, ` 3`); den Default fährt jeder
# Lauf von `make gates`.
#
# Ein roter Fall prüft neben dem Ausgang seine eigene Zeile: einen Befund von
# golangci-lint an Pfad und Zeile mit dem Namen des Linters oder eine `lint:`-Zeile
# im Wortlaut. Die Zeile eines Befunds findet die Gegenprobe über einen Anker im
# Quelltext des Falls. Die Fälle der Kopie `sammel` teilen sich einen Lauf, jeder an
# seiner eigenen Zeile erkannt; allein läuft ein Fall, dessen Zusage der Ausgang
# ist oder der andere ausschließt.
#
# GEPRÜFT WIRD — Fälle je Punkt von SPEC-049:
#   Grundlauf — p0-grundlauf: die unveränderte Kopie endet mit Ausgang 0 ohne
#       `lint:`-Zeile, also auch ohne ungenutzte Regel; das ist der grüne Fall je
#       dauerhafter Ausnahme (Punkt 8). Punkt 2 hat keinen eigenen Fall (OFFEN).
#   (1) Gegenstand — p1-integration-build-tag, p1-cmd, p1-test,
#       p1-pfad-anker-global, p1-pfad-anker-forbidigo
#   (3) Linter — p3-genau-diese und p3-default-none (die Liste im Profil); je
#       aktivem Linter ein roter Fall: p3-<linter> für die Linter ohne Schwelle (17),
#       p4-<linter>-rot für die acht mit Schwelle, p5-revive-* für revive,
#       p5-contextcheck-test für contextcheck, p7-whitebox-* für testpackage
#   (4) Schwellen — je Schwelle ein grüner und ein roter Fall an der Grenze:
#       p4-cyclop-*, p4-gocyclo-*, p4-gocognit-*, p4-nestif-*,
#       p4-funlen-anweisungen-*, p4-funlen-zeilen-*, p4-maintidx-*, p4-dupl-*,
#       p4-interfacebloat-*
#   (5) Einstellungen — p5-errcheck-frei-* und p5-errcheck-rot,
#       p5-ireturn-frei-*, p5-forbidigo-* und p1-cmd, p5-gomodguard-frei-*,
#       p5-revive-genau-diese, p5-revive-<regel> je Regel, p5-revive-exported-form
#       und p9-lint-zeile-allein (ein deutscher Doc-Kommentar der richtigen Form ohne
#       Befund), für testpackage p7-internal-test, p7-endet-auf-export-test und
#       p7-bruecke-kein-testpackage,
#       p5-contextcheck-test, p5-gochecknoglobals-frei-*
#   (6) Kein `//nolint` — p6-doppelstrich, p6-leerzeichen-zusatz, p6-block,
#       p6-schreibweise, p6-drei-striche, p6-leerzeichen-dann-strich, p6-tab,
#       p6-zweites-zeichen, p6-string, p6-cmd, p6-test, p6-testdatei; kein Befund:
#       p6-siehe-nolint, p6-wortgrenze
#   (7) Export-Test-Brücke — erste Bedingung von Messmethode 3 je Paketgruppe:
#       p7-whitebox-kern, p7-whitebox-driven, p7-whitebox-pgwire,
#       p7-whitebox-einstieg; zweite Bedingung: p7-internal-test,
#       p7-endet-auf-export-test,
#       p7-test-in-bruecke-<art> (Test, Benchmark, Example, Fuzz, tab),
#       p7-variable-in-bruecke; kein Befund: p7-bruecke-kein-testpackage,
#       p7-bruecke-weiterreichen
#   (8) Ausnahmen — Regeln: p8-ausnahme-name-<name> je benannter Globale (gilt
#       nur für den Namen), p8-ausnahme-test-<linter> (gilt für jede Testdatei, auch
#       unter test/), p8-ausnahme-st1005; Why: p8-ohne-why, p8-why-nicht-erste,
#       p8-why-leerzeile, p8-why-ohne-doppelpunkt, p8-why-vorhanden; ungenutzt: p8-ungenutzt-feldfolge (mit
#       `\` und `"` im Wert), p8-ungenutzt-pfad-ausser; feste Form, rot:
#       p8-form-einzug, p8-form-andere-stelle, p8-form-fluss-exclusions,
#       p8-form-fluss-rules, p8-form-fluss-anderswo, p8-form-listenpunkt,
#       p8-form-quote-exclusions, p8-form-quote-rules, p8-form-ohne-strich,
#       p8-form-fortsetzung-7, p8-form-eintrag-8 (die vollständige Liste der Zeilen),
#       p8-form-eintrag-2 (nur die erste Eintragszeile und Ausgang 1),
#       p8-form-strich-allein-8 (`-` allein mit Einzug 8); feste Form, kein Befund:
#       p8-kommentar-nach-exclusions, p8-kommentar-nach-rules,
#       p8-schluessel-nach-rules, p8-reihenfolge (`exclusions` vor `settings`),
#       p8-leerraum (Leerzeichen und Tab nach `exclusions:` und `rules:`),
#       p8-strich-allein (`-` allein mit `# Why:` darüber); `-` allein ohne
#       Kommentarblock: p8-strich-allein-ohne-why (Why-Zeile, kein `Form nicht erkannt`)
#   (9) Ausgabe und Ausgang — p9-profil-fehlt-*, p9-schema, p9-schema-meldung,
#       p9-max-same-issues-*, p9-max-issues-per-linter-*, p9-uniq-by-line-*,
#       p9-laden (rot ohne `lint:`-Zeile); Ausgang 1 je Art einer Zeile allein:
#       p9-lint-zeile-allein, p9-bruecke-allein, p9-form-allein, p9-fehlt-allein,
#       p9-schema-allein, p9-ungenutzt-allein, p9-golangci-allein
#   (10) Gate — p10-gate-checks (`lint` und `lint-gegenprobe` an GATE_CHECKS),
#       p10-gates-rot (`make gates` endet bei rotem `make lint` mit Fehlerstatus),
#       p10-gegenprobe-scharf (`make lint-gegenprobe` endet bei rotem Skript mit
#       Fehlerstatus)
#   Grenze — kein Befund: neg-main-test (Testdatei im Paket main), neg-generiert
#       (Markierung für generierten Code)
#
# OFFEN — Zusagen ohne Mutation, die ein Lauf hier unterscheidet:
#   - `relative-path-mode: cfg` (Punkt 1): Profil, Modulwurzel und
#     Arbeitsverzeichnis sind in der Stufe dasselbe /src. Den Anker `^` der Pfade
#     unterscheiden p1-pfad-anker-*.
#   - `--network=none` und die Module aus deps (Punkt 2): keine Prüfung braucht
#     Netz; ob die Module aus deps kommen, sähe nur ein Lauf mit Netz und ohne sie.
#   - Pin und Plattform des Images (Punkt 2): Ein anderes Image ließe sich nur mit
#     Netz ziehen. Eine Ausgabe der Version in der Stufe diente nur der Gegenprobe,
#     ein Textvergleich mit der `FROM`-Zeile wäre eine zweite Quelle für den Pin, den
#     Punkt 2 nur im Dockerfile führt. Eine Anhebung ist ein Commit an dieser Zeile;
#     sie fängt der Re-Evaluierungs-Trigger von ADR-0034 und das Review.
#   - `GOFLAGS=-mod=readonly` (Punkt 2): ohne vendor/ wählt go denselben Modus.
#   - `GOTOOLCHAIN=local` und die Go-Version gegen die `go`-Zeile (Punkt 2): eine
#     höhere `go`-Zeile lässt schon die Stufe deps scheitern, bevor `lint` läuft.
#   - `-c` statt Default-Suche (Punkt 9): der Build-Kontext führt nur .golangci.yml,
#     die Default-Suche fände dieselbe Datei.
#   - die Zeile `Regel ohne Befund: Felder nicht erkannt` (Punkt 9): Die gepinnte
#     Version von golangci-lint nennt in jeder Warnung mindestens zwei der Felder;
#     ein anderer Logtext käme nur mit einer anderen Version, also mit einem Commit am
#     Pin (wie oben).
#   - nichts in den Arbeitsbaum (Kopf von SPEC-049), weder durch `make lint` noch
#     durch diese Gegenprobe: kein Fall vergleicht den Arbeitsbaum vor und nach
#     einem Lauf. Ebenso, dass die Gegenprobe ihre Kopien löscht und bei einem
#     Abbruch über `set -e` laufende Builds beendet und eine Zeile `ROT` schreibt.
#   - Punkt 8, zulässige Regeln: ob eine Regel zu den zulässigen gehört und mehr als
#     Bestand ausblendet, prüft das Werkzeug nicht (Grenze von SPEC-049); kein Fall.
#   Ob ein `Why:` zutrifft, ob die Brücke nur weiterreicht, ob ein Wert eines
#   unexportierten Typs nur aus dem Produkt-Code stammt, ob eine Einstellung ihren
#   Grund trägt und welche Schlüssel unter `exclusions` stehen, prüft das Werkzeug
#   nicht (Grenze von SPEC-049); dafür gibt es keinen Fall.
#
# Ausgang: 0, wenn jeder Fall wie erwartet endet, sonst 1.
set -euo pipefail
export LC_ALL=C

wurzel="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
docker_build="${DOCKER_BUILD:-docker build}"
if [ -z "${LINT_GEGENPROBE_PARALLEL+gesetzt}" ]; then
  parallel=6
elif [[ "$LINT_GEGENPROBE_PARALLEL" =~ ^[1-9][0-9]*$ ]]; then
  parallel="$LINT_GEGENPROBE_PARALLEL"
else
  echo "lint-gegenprobe: LINT_GEGENPROBE_PARALLEL ist keine positive ganze Zahl: '$LINT_GEGENPROBE_PARALLEL'" >&2
  exit 2
fi
arbeit="$(mktemp -d)"
fehler=0
aufraeumen() {
  local status=$?
  local j
  j="$(jobs -p)"
  if [ -n "$j" ]; then kill $j 2>/dev/null || true; wait 2>/dev/null || true; fi
  rm -rf "$arbeit"
  if [ "$status" -ne 0 ] && [ "$fehler" -eq 0 ]; then
    echo "lint-gegenprobe: ROT — Abbruch mit Ausgang $status" >&2
  fi
}
trap aufraeumen EXIT
G=internal/gegenprobe
P=.golangci.yml

melde() { echo "lint-gegenprobe: ROT — $*" >&2; fehler=1; }

# kopie <name> — legt eine Kopie des Arbeitsbaums ohne .git an; k ist ihre Wurzel.
kopie() {
  k="$arbeit/$1"
  mkdir -p "$k"
  local e
  for e in "$wurzel"/* "$wurzel"/.[!.]*; do
    [ -e "$e" ] || continue
    [ "${e##*/}" = .git ] && continue
    cp -R "$e" "$k/"
  done
}

# datei <pfad> — schreibt stdin nach <pfad> in der Kopie k.
datei() {
  mkdir -p "$(dirname "$k/$1")"
  cat > "$k/$1"
  touch "$k/$1"
}

# aendere <pfad> <sed-ausdruck> — ändert eine Datei der Kopie k mit sed.
aendere() {
  sed -i "$2" "$k/$1"
  touch "$k/$1"
}

# starte <name> [<befehl>…] — fährt im Hintergrund in der Kopie <name> den Build der
# Stufe lint oder den genannten Befehl; Ausgabe nach <name>.log, Ausgang nach
# <name>.exit.
starte() {
  local name="$1"
  shift
  while [ "$(jobs -rp | wc -l)" -ge "$parallel" ]; do wait -n || true; done
  (
    cd "$arbeit/$name"
    set +e
    if [ "$#" -eq 0 ]; then
      $docker_build --progress=plain --target lint . >"$arbeit/$name.log" 2>&1 &
    else
      "$@" >"$arbeit/$name.log" 2>&1 &
    fi
    kind=$!
    trap 'kill "$kind" 2>/dev/null' TERM
    wait "$kind"
    { echo $? >"$arbeit/$name.exit"; } 2>/dev/null
  ) &
}

# n <kopie> — die Ausgabe ohne das Präfix `#<schritt> <sekunden> ` von BuildKit.
n() { echo "$arbeit/$1.n"; }

# zeile_von <kopie> <pfad> <anker> — Nummer der ersten Zeile mit <anker>.
zeile_von() { grep -n -m1 -F -- "$3" "$arbeit/$1/$2" | cut -d: -f1 || true; }

# befund_muster <kopie> <pfad> <anker> <linter> — ERE des Befunds an der Zeile.
befund_muster() {
  local z
  z="$(zeile_von "$1" "$2" "$3")"
  [ -n "$z" ] || return 1
  printf '^%s:%s(:[0-9]+)?: .*\\(%s\\)$' "${2//./\\.}" "$z" "$4"
}

# befund <fall> <kopie> <pfad> <anker> <linter> [<text-ere>] — erwartet einen Befund
# des Linters an der Zeile mit <anker>, mit <text-ere> im Text, falls genannt.
befund() {
  local m
  m="$(befund_muster "$2" "$3" "$4" "$5")" || { melde "Fall '$1': Anker '$4' fehlt in $3"; return; }
  if ! grep -E -- "$m" "$(n "$2")" | grep -E -- "${6:-.}" >/dev/null; then
    melde "Fall '$1': erwartet einen Befund ($5) an $3, Zeile mit '$4'"
  fi
}

# kein_befund <fall> <kopie> <pfad> <anker> <linter> — erwartet keinen Befund des
# Linters an der Zeile mit <anker>.
kein_befund() {
  local m
  m="$(befund_muster "$2" "$3" "$4" "$5")" || { melde "Fall '$1': Anker '$4' fehlt in $3"; return; }
  if grep -qE -- "$m" "$(n "$2")"; then
    melde "Fall '$1': unerwarteter Befund ($5) an $3, Zeile mit '$4'"
  fi
}

# lint_zeile <fall> <kopie> <zeile> — erwartet die `lint:`-Zeile im Wortlaut.
lint_zeile() {
  grep -qxF -- "$3" "$(n "$2")" || melde "Fall '$1': erwartet die Zeile '$3'"
}

# keine_lint_zeile <fall> <kopie> [<zeile>] — erwartet die genannte `lint:`-Zeile
# nicht, ohne <zeile> gar keine.
keine_lint_zeile() {
  if [ "$#" -ge 3 ]; then
    if grep -qxF -- "$3" "$(n "$2")"; then melde "Fall '$1': unerwartete Zeile '$3'"; fi
  elif grep -q '^lint: ' "$(n "$2")"; then
    melde "Fall '$1': unerwartete Zeile '$(grep -m1 '^lint: ' "$(n "$2")")'"
  fi
}

# rot <fall> <kopie> — erwartet, dass die Stufe mit Ausgang 1 endet.
rot() {
  if [ "$(cat "$arbeit/$2.exit")" -eq 0 ] \
    || ! grep -qF 'did not complete successfully: exit code: 1' "$(n "$2")"; then
    melde "Fall '$1': erwartet die Stufe mit Ausgang 1, Lauf endete mit $(cat "$arbeit/$2.exit")"
  fi
}

# Generatoren für die Schwellen, je <name> <n>.
# wenn_kette — n Abfragen nacheinander: zyklomatisch n+1, kognitiv n.
wenn_kette() {
  printf 'func %s(a int) int {\n' "$1"
  local i
  for i in $(seq 1 "$2"); do printf '\tif a == %d {\n\t\treturn %d\n\t}\n' "$i" "$i"; done
  printf '\treturn 0\n}\n'
}
# wenn_geschwister — eine Abfrage mit n Abfragen darin: Komplexität für nestif n.
wenn_geschwister() {
  printf 'func %s(a int) int {\n\tif a > 0 { // %s\n' "$1" "$1"
  local i
  for i in $(seq 1 "$2"); do printf '\t\tif a > %d {\n\t\t\ta++\n\t\t}\n' "$i"; done
  printf '\t}\n\treturn a\n}\n'
}
# anweisungen — n Anweisungen und ein return: n+1 Anweisungen.
anweisungen() {
  printf 'func %s(a int) int {\n' "$1"
  local i
  for i in $(seq 1 "$2"); do printf '\ta += %d\n' "$i"; done
  printf '\treturn a\n}\n'
}
# zeilen — eine Anweisung, der Rumpf hat n+2 Zeilen.
zeilen() {
  printf 'func %s() []int {\n\treturn []int{\n' "$1"
  local i
  for i in $(seq 1 "$2"); do printf '\t\t%d,\n' "$i"; done
  printf '\t}\n}\n'
}
# doppelt <name> <schluss> — 45 Anweisungen einer Form, die sonst kein Fall nutzt,
# dann <schluss> (für dupl). Gemessen mit golangci-lint v2.14.0: mit `{}` ist ein
# Paar 149 groß, mit `a++` 150.
doppelt() {
  printf 'func %s(a int) int {\n' "$1"
  local i
  for i in $(seq 1 45); do printf '\ta -= %d\n' "$i"; done
  printf '\t%s\n\treturn a\n}\n' "$2"
}
# schnittstelle — eine Schnittstelle mit n Methoden.
schnittstelle() {
  printf 'type %s interface {\n' "$1"
  local i
  for i in $(seq 1 "$2"); do printf '\tM%d()\n' "$i"; done
  printf '}\n'
}

# Zeilen des Ausnahme-Abschnitts im Profil des Arbeitsbaums: `exclusions:`,
# `rules:`, der erste Eintrag, dessen erste Fortsetzung ohne `- ` und die letzte
# Zeile des Abschnitts.
exkl="$(grep -n -m1 '^  exclusions:$' "$wurzel/$P" | cut -d: -f1 || true)"
regeln="$(grep -n -m1 '^    rules:$' "$wurzel/$P" | cut -d: -f1 || true)"
erste="$(awk -v r="${regeln:-0}" 'r && NR > r && /^      - / { print NR; exit }' "$wurzel/$P")"
fortsetzung="$(awk -v r="${erste:-0}" 'r && NR > r && /^        [^ -]/ { print NR; exit }' "$wurzel/$P")"
ende="$(awk -v e="${exkl:-0}" 'NR > e && /^[^ #]/ { print NR - 1; f = 1; exit } END { if (!f) print NR }' "$wurzel/$P")"
if [ -z "$exkl" ] || [ -z "$regeln" ] || [ -z "$erste" ] || [ -z "$fortsetzung" ]; then
  melde "das Profil des Arbeitsbaums hat nicht die feste Form (exclusions '$exkl', rules '$regeln', erster Eintrag '$erste')"
  exit 1
fi

# --- Profil: aktive Linter und Regeln von revive (Punkt 3 und 5) --------------
# Verglichen mit dem Text des Profils; dass jeder Linter und jede Regel wirkt,
# zeigen die roten Fälle unten.
linter_soll="containedctx contextcheck cyclop dupl errcheck fatcontext forbidigo funlen
gochecknoglobals gochecknoinits gocognit gocyclo gomodguard_v2 govet iface inamedparam
ineffassign interfacebloat ireturn maintidx nestif noctx reassign revive staticcheck
testpackage unparam unused"
linter_ist="$(awk '/^  enable:$/ { f = 1; next } f && /^    - / { print $2; next } f { exit }' "$wurzel/$P" | sort | tr '\n' ' ')"
[ "$linter_ist" = "$(printf '%s\n' $linter_soll | sort | tr '\n' ' ')" ] \
  || melde "Fall 'p3-genau-diese': aktiv sind $linter_ist"
grep -qx '  default: none' "$wurzel/$P" || melde "Fall 'p3-default-none': linters.default ist nicht none"
revive_soll="blank-imports context-as-argument context-keys-type dot-imports empty-block
error-naming error-return error-strings errorf exported if-return increment-decrement
indent-error-flow package-comments range receiver-naming redefines-builtin-id
superfluous-else time-naming unexported-return unreachable-code unused-parameter
var-declaration var-naming unused-receiver"
revive_ist="$(awk '/^    revive:$/ { r = 1; next } r && /^      rules:$/ { f = 1; next } f && /^        / { print; next } f { exit }' "$wurzel/$P" | sort | tr '\n' '|')"
[ "$revive_ist" = "$(printf '        - name: %s\n' $revive_soll | sort | tr '\n' '|')" ] \
  || melde "Fall 'p5-revive-genau-diese': die Regeln von revive sind $revive_ist"

# --- LINT_GEGENPROBE_PARALLEL: ungültige Werte -------------------------------
# Das Skript läuft als Kopie in einem Baum ohne Profil, damit es nach der Prüfung
# der Einstellung nichts baut; ohne die Prüfung endet es dort mit Ausgang 1.
mkdir -p "$arbeit/einstellung/tools/harness"
cp "${BASH_SOURCE[0]}" "$arbeit/einstellung/tools/harness/lint-gegenprobe.sh"
for wert in 0 '' -1 abc ' 3'; do
  set +e
  aus="$(LINT_GEGENPROBE_PARALLEL="$wert" DOCKER_BUILD=false timeout 20 \
    bash "$arbeit/einstellung/tools/harness/lint-gegenprobe.sh" 2>&1)"
  code=$?
  set -e
  soll="lint-gegenprobe: LINT_GEGENPROBE_PARALLEL ist keine positive ganze Zahl: '$wert'"
  if [ "$code" -ne 2 ] || [ "$aus" != "$soll" ]; then
    melde "Fall 'einstellung-$wert': erwartet Ausgang 2 und nur die Zeile '$soll', bekam $code: $aus"
  fi
done

# --- Grundlauf ---------------------------------------------------------------
kopie grundlauf
starte grundlauf

# --- sammel: die Fälle mit eigener Zeile -------------------------------------
kopie sammel
L="$G/linter/fall.go"
datei "$L" <<'EOF'
package linter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driven"
	_ "go.yaml.in/yaml/v3"
	_ "golang.org/x/sync/errgroup"
)

func errcheckFall() { os.Remove("errcheck-fall") }
func govetFall() string { return fmt.Sprintf("%d", "govet-fall") }
func ineffassignFall() int { x := 1; x = 2; return x }
func staticcheckFall(s string) string { return strings.Replace(s, "a", "b", 0) }
func unusedFall() {}
type containedctxFall struct{ ctx context.Context }
type fatcontextSchluessel struct{}
func fatcontextFall(ctx context.Context) context.Context {
	for i := 0; i < 3; i++ { ctx = context.WithValue(ctx, fatcontextSchluessel{}, i) }
	return ctx
}
var gochecknoglobalsFall = 1
var version = "frei"
var _ = 2
var ErrFrei = errors.New("frei")
func init() {}
type ifaceErste interface{ Methode() }
type ifaceZweite interface{ Methode() }
type inamedparamFall interface{ Methode(int) }
type ireturnEigene interface{ Methode() }
func ireturnFall() ireturnEigene { return nil }
func ireturnError() error { return nil }
func ireturnLeer() interface{} { return nil }
func ireturnAnonym() interface{ Methode() } { return nil }
func ireturnStdlib() io.Reader { return nil }
func ireturnGenerisch[T ireturnEigene](t T) T { return t }
func ireturnPorts() driven.Upstream { return nil }
func ireturnPgproto3() pgproto3.FrontendMessage { return nil }
func noctxFall() { r, err := http.Get("http://noctx.invalid"); if err == nil { _ = r.Body.Close() } }
func reassignFall() { io.EOF = nil }
func unparamLeer(a int) int { if true { return 2 }; return 3 }
func UnparamFall() int { return unparamLeer(4) + unparamLeer(5) }
func errcheckFrei(w io.Writer, c net.Conn, l net.Listener) {
	fmt.Fprint(w, "fprint")
	fmt.Fprintf(w, "fprintf")
	fmt.Fprintln(w, "fprintln")
	c.Close()
	l.Close()
	w.Write(nil)
}
func forbidigoFall() {
	fmt.Print("print")
	fmt.Printf("printf")
	fmt.Println("println")
	print("builtin-print")
	println("builtin-println")
	_ = os.Stdout
	_ = os.Stderr
}
func gleich() {
	os.Chdir("gleich-1")
	os.Chdir("gleich-2")
	os.Chdir("gleich-3")
	os.Chdir("gleich-4")
}
var zeileA, zeileB = 1, 2
func st1005Frei() error { return errors.New("Großbuchstabe am Anfang") }
EOF
{
  echo "package linter"
  wenn_kette zyklomatisch15 14
  wenn_kette zyklomatisch16 15
  wenn_kette kognitiv20 20
  wenn_kette kognitiv21 21
  wenn_kette wartbar20 55
  wenn_kette wartbar19 56
  wenn_geschwister verschachtelt4 4
  wenn_geschwister verschachtelt5 5
  anweisungen anweisungen60 59
  anweisungen anweisungen61 60
  zeilen zeilen100 98
  zeilen zeilen101 99
  doppelt doppeltKurzA '{}'
  doppelt doppeltKurzB '{}'
  schnittstelle breit10 10
  schnittstelle breit11 11
  for i in $(seq 1 51); do echo "var viele$i = $i"; done
} | datei "$G/linter/schwellen.go"
{ echo "package doppelt"; doppelt langA 'a++'; doppelt langB 'a++'; } | datei "$G/doppelt/fall.go"
# Die Ausnahmen für Testdateien gelten für jede Testdatei, auch unter test/.
{
  printf 'package gegenprobe_test\n\nimport "net/http"\n\n'
  wenn_kette testZyklomatisch16 15
  wenn_kette testKognitiv21 21
  wenn_geschwister testVerschachtelt5 5
  anweisungen testAnweisungen61 60
  zeilen testZeilen101 99
  printf 'func testNoctx() { r, err := http.Get("http://noctx.invalid"); if err == nil { _ = r.Body.Close() } }\n'
  printf 'func testUnparam(a int) int { if true { return 2 }; return 3 }\n'
  printf 'func testParameter(a int) {}\n'
  printf 'type testTyp struct{}\n'
  printf 'func (t testTyp) ohneEmpfaenger() int { return 1 }\n'
} | datei test/gegenprobe/ausnahme_test.go
datei "$G/kontext/fall_test.go" <<'EOF'
package kontext_test

import (
	"context"
	"testing"
)

func benutzt(ctx context.Context) { _ = ctx }
func ohneKontext()                 { benutzt(context.Background()) }
func mitKontext(ctx context.Context) {
	benutzt(ctx)
	ohneKontext() // kontext-fall
}
func TestKontext(t *testing.T) { mitKontext(t.Context()) }
EOF
R="$G/revive/fall.go"
datei "$R" <<'EOF'
// Package revive trägt je Regel von revive einen Verstoß.
package revive

import (
	"context"
	"errors"
	"fmt"
	_ "strings"
	"time"
	. "unicode"
)

func punktImport() bool { return IsUpper('A') }
func kontextHinten(a int, ctx context.Context) { _, _ = a, ctx }
func schluessel(ctx context.Context) context.Context { return context.WithValue(ctx, "k", 1) }
func leer(a bool) { if a { } }
var fehlerWert = errors.New("fehler-wert")
func fehlerVorn() (error, int) { return fehlerWert, 0 }
func punktAmEnde() error { return errors.New("endet mit punkt.") }
func formatiert(a int) error { return errors.New(fmt.Sprintf("formatiert %d", a)) }
func OhneKommentar() {}
// liefert nichts.
func FalscheForm() {}
func weiter(f func() error) error { if err := f(); err != nil { return err }; return nil }
func erhoehe(a int) int { a += 1; return a }
func sonst(a bool) int { if a { return 1 } else { return 2 } }
func bereich(s []int) (n int) { for i, _ := range s { n += i }; return n }
type typ struct{ wert int }
func (this typ) methode() int { return this.wert }
func (t typ) ohneEmpfaenger() int { return 1 }
func eingebaut() int { len := 1; return len }
func schleife(s []int) (n int) { for _, v := range s { if v > 0 { continue } else { n++ } }; return n }
func zeit() time.Duration { var dauerSecs time.Duration; return dauerSecs }
type intern struct{}
// Liefert liefert einen unexportierten Typ.
func Liefert() intern { return intern{} }
func toter() {
	return // toter
	_ = 2
}
func unbenutzt(a int) {}
func deklaration() int { var x int = 0; return x }
func benennung() int { meinId := 1; return meinId }
EOF
datei "$G/paketkommentar/fall.go" <<<'package paketkommentar'
datei cmd/gegenprobe/main.go <<'EOF'
package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = os.Stdout, os.Stderr
	fmt.Println("cmd-println")
	_ = 1 // nolint-zeile: //nolint
}
EOF
datei cmd/gegenprobe/main_test.go <<'EOF'
package main

import "testing"

func TestMain(t *testing.T) { _ = t }
EOF
datei test/gegenprobe/fall.go <<'EOF'
// Package gegenprobe liegt unter test/.
package gegenprobe

import "os"

var testGlobal = os.Stdout

func schreibe() { _ = os.Stderr }

var _ = 1 // nolint-zeile: //nolint
EOF
datei "$G/test/fall.go" <<'EOF'
// Package test liegt nicht unter test/.
package test

import "os"

var _ = os.Stdout
EOF
datei test/anker/internal/hexagon/model/fehler.go <<'EOF'
// Package model liegt nicht unter internal/hexagon/model.
package model

var klassen = map[int]string{}
EOF
# Eine Regel für eine benannte Globale gilt nur für ihren Namen: je Regel eine
# Globale in ihrer Datei, deren Name auf den ausgenommenen endet.
namen="$(awk '
  /^      - / { if (pfad != "" && name != "") print pfad, name; pfad = ""; name = "" }
  /^        path: / { pfad = $2 }
  /^        text: \^?[A-Za-z]+ is a global variable\$?$/ { name = $2; sub(/^\^/, "", name) }
  END { if (pfad != "" && name != "") print pfad, name }
' "$wurzel/$P")"
[ -n "$namen" ] || melde "Fall 'p8-ausnahme-name': das Profil nennt keine Globale"
while read -r pfad name; do
  d="${pfad#^}"; d="${d%\$}"; d="${d//\\./.}"
  printf '\nvar gegenprobe%s = 1\n' "$name" >> "$k/$d"
  touch "$k/$d"
done <<<"$namen"
datei test/integration/gegenprobe_e2e_test.go <<'EOF'
//go:build integration

package integration_test

var gegenprobeIntegration = 1
EOF
datei "$G/generiert/fall.go" <<'EOF'
// Code generated by gegenprobe. DO NOT EDIT.

package generiert

var generiertGlobal = 1
EOF
NF="$G/nolint/fall.go"
datei "$NF" <<'EOF'
package nolint

func wert() int {
	a := 1 //nolint
	a++    // nolint:errcheck // Begruendung
	a++    /* nolint */
	a++    //NoLint
	a++    ///nolint
	a++    //	nolint
	a++    // x //nolint
	a++    // siehe nolint
	a++    //nolintx
	a++    // /nolint
	_ = "//nolint"
	return a
}
EOF
datei "$G/nolint/fall_test.go" <<'EOF'
package nolint_test

import "testing"

func TestDirektive(t *testing.T) { _ = t } //nolint
EOF
for p in internal/hexagon/model internal/adapters/driven/recording \
  internal/adapters/driving/pgwire internal/adapters/driving/cli; do
  printf 'package %s\n\nimport "testing"\n\nfunc TestWeiss(t *testing.T) { _ = t }\n' "${p##*/}" \
    | datei "$p/weiss_test.go"
done
B="$G/bruecke"
datei "$B/code.go" <<'EOF'
// Package bruecke hat eine Brücke.
package bruecke

func intern() int { return 1 }
EOF
datei "$B/export_test.go" <<'EOF'
package bruecke

import "testing"

var Zustand = 1

func TestBruecke(t *testing.T)      { _ = t }
func BenchmarkBruecke(b *testing.B) { _ = b }
func ExampleIntern()                {}
func FuzzBruecke(f *testing.F)      { _ = f }
func	TestTab(t *testing.T)         { _ = t }
func Intern() int                   { return intern() + Zustand }
EOF
datei "$B/fooexport_test.go" <<'EOF'
package bruecke

import "testing"

func TestFooExport(t *testing.T) { _ = intern() }
EOF
datei "$B/internal_test.go" <<'EOF'
package bruecke

import "testing"

func TestIntern(t *testing.T) { _ = intern() }
EOF
datei "$B/bruecke_test.go" <<'EOF'
package bruecke_test

import (
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/gegenprobe/bruecke"
)

func TestBlackBox(t *testing.T) { _ = bruecke.Intern() }
EOF
# Regeln im Profil: ohne Why, Why nicht zuerst, Why durch Leerzeile getrennt,
# ungenutzt mit allen Feldern außer Pfad außer, ungenutzt mit Pfad außer; dazu ein
# Schlüssel, den das Schema ablehnt.
datei "$G/ohnewhy/fall.go" <<'EOF'
// Package ohnewhy hat ein init, das die Regel ohne Why ausblendet.
package ohnewhy

func init() {}
EOF
aendere "$P" "${exkl}s/\$/ # Kommentar nach dem Schlüssel/;${regeln}s/\$/ # Kommentar/"
n0="$(wc -l < "$k/$P")"
cat >> "$k/$P" <<'EOF'

      - linters:
          - gochecknoinits
        path: ^internal/gegenprobe/ohnewhy/

      # Grund: steht vor dem Why.
      # Why: zweite Zeile des Blocks.
      - linters:
          - gochecknoinits
        path: ^internal/gegenprobe/nirgends-b/

      # Why: durch eine Leerzeile getrennt.

      - linters:
          - gochecknoinits
        path: ^internal/gegenprobe/nirgends-c/

      # Why: Feldfolge der ungenutzten Regel.
      - source: ^gegenprobe-quelle$
        text: ^gegenprobe\.text"$
        path: ^internal/gegenprobe/nirgends/
        linters:
          - revive

      # Why: Pfad außer in der ungenutzten Regel.
      - text: ^gegenprobe-ausser$
        path-except: ^internal/gegenprobe/
        linters:
          - revive

      # Why kein Doppelpunkt.
      - linters:
          - gochecknoinits
        path: ^internal/gegenprobe/nirgends-d/

    paths: []
gegenprobe-unbekannt: true
EOF
touch "$k/$P"
starte sammel

# --- allein: Profil fehlt ----------------------------------------------------
kopie profil-fehlt
rm "$k/$P"
datei "$G/fehlt/fall.go" <<'EOF'
package fehlt

import "os"

func fehlt() { os.Remove("fehlt") }

var _ = 1 //nolint
EOF
datei "$G/fehlt/export_test.go" <<'EOF'
package fehlt

func TestFehlt() {}
EOF
starte profil-fehlt

# --- allein: Ablehnung erst beim Laden ---------------------------------------
kopie laden
cat >> "$k/$P" <<'EOF'

      # Why: path und path-except zusammen nimmt config verify an.
      - path: ^internal/gegenprobe/
        path-except: ^internal/gegenprobe/laden/
EOF
touch "$k/$P"
starte laden

# --- allein: ein Befund von golangci-lint ohne `lint:`-Zeile -----------------
kopie golangci-allein
datei "$G/allein/fall.go" <<'EOF'
// Package allein hat ein init.
package allein

func init() {}
EOF
starte golangci-allein

# --- allein: Form des Profils ------------------------------------------------
kopie form-einzug
aendere "$P" "${exkl},${ende}s/^\\(.\\)/  \\1/"
starte form-einzug

kopie form-andere-stelle
printf 'formatters:\n  exclusions:\n    paths: []\n' >> "$k/$P"
touch "$k/$P"
andere="$(grep -n '^  exclusions:$' "$k/$P" | tail -1 | cut -d: -f1)"
starte form-andere-stelle

kopie form-fluss-exclusions
aendere "$P" "${exkl},${ende}c\\  exclusions: {warn-unused: true, rules: []}"
starte form-fluss-exclusions

kopie form-fluss-rules
aendere "$P" "${regeln},${ende}c\\    rules: []"
starte form-fluss-rules

kopie form-eintrag-8
aendere "$P" "$((regeln + 1)),${ende}s/^\\(.\\)/  \\1/"
starte form-eintrag-8

kopie form-eintrag-2
aendere "$P" "${erste}s/^      - /  - /"
starte form-eintrag-2

kopie form-fortsetzung-7
aendere "$P" "${fortsetzung}s/^        /       /"
starte form-fortsetzung-7

kopie form-ohne-strich
aendere "$P" "${regeln}a\\      gegenprobe: 1"
starte form-ohne-strich

kopie form-quote-exclusions
aendere "$P" "${exkl}s/exclusions/\"exclusions\"/"
starte form-quote-exclusions

kopie form-quote-rules
aendere "$P" "${regeln}s/rules/'rules'/"
starte form-quote-rules

kopie form-fluss-anderswo
printf 'formatters: {exclusions: {paths: []}}\n' >> "$k/$P"
touch "$k/$P"
starte form-fluss-anderswo

kopie form-listenpunkt
printf 'gegenprobe-liste:\n  - exclusions: 1\n' >> "$k/$P"
touch "$k/$P"
starte form-listenpunkt

# Grün: `exclusions` vor `settings` ist dieselbe feste Form.
kopie form-reihenfolge
einst="$(grep -n -m1 '^  settings:$' "$wurzel/$P" | cut -d: -f1 || true)"
if [ -n "$einst" ] && [ "$einst" -lt "$exkl" ]; then
  {
    sed -n "1,$((einst - 1))p" "$wurzel/$P"
    sed -n "${exkl},${ende}p" "$wurzel/$P"
    sed -n "${einst},$((exkl - 1))p" "$wurzel/$P"
    sed -n "$((ende + 1)),\$p" "$wurzel/$P"
  } | datei "$P"
else
  melde "Fall 'p8-reihenfolge': settings steht nicht vor exclusions"
fi
starte form-reihenfolge

# Grün: Leerzeichen und Tab nach `exclusions:` und `rules:` sind Blockform.
kopie form-leerraum
aendere "$P" "${exkl}s/\$/ \t/;${regeln}s/\$/\t /"
starte form-leerraum

# `-` allein auf der Zeile ist ein Eintrag, sein Inhalt steht auf Fortsetzungszeilen:
# grün mit dem `# Why:`-Block darüber, rot ohne ihn.
kopie form-strich-allein
aendere "$P" "${erste}s/^      - /      -\n        /"
starte form-strich-allein

kopie form-strich-ohne-why
why_anfang="$(awk -v z="$erste" 'NR < z && /^[ \t]*#/ { if (!a) a = NR; next } NR < z { a = 0 } END { print a }' "$wurzel/$P")"
aendere "$P" "${erste}s/^      - /      -\n        /;${why_anfang},$((erste - 1))d"
strich_zeile="$why_anfang"
starte form-strich-ohne-why

# Rot: `-` allein mit Einzug 8 bleibt `Form nicht erkannt`.
kopie form-strich-allein-8
aendere "$P" "${regeln}a\\        -\n          linters:\n            - revive"
starte form-strich-allein-8

# --- allein: je Art einer `lint:`-Zeile ohne andere Zeile (Punkt 9) ----------
kopie fehlt-allein
rm "$k/$P"
starte fehlt-allein

kopie schema-allein
printf 'gegenprobe-unbekannt: true\n' >> "$k/$P"
touch "$k/$P"
starte schema-allein

kopie ungenutzt-allein
cat >> "$k/$P" <<'EOF'

      # Why: ungenutzt.
      - path: ^internal/gegenprobe/nirgends/
        linters:
          - revive
EOF
touch "$k/$P"
starte ungenutzt-allein

kopie bruecke-allein
datei "$G/bruecke/code.go" <<'EOF'
// Package bruecke hat eine Brücke.
package bruecke

func intern() int { return 1 }
EOF
datei "$G/bruecke/export_test.go" <<'EOF'
package bruecke

import "testing"

func TestBruecke(t *testing.T) { _ = intern() }
EOF
starte bruecke-allein

# --- Gate: eine Kopie nur mit dem Fragment lint.mk ---------------------------
# Die Kopie trägt eine `lint:`-Zeile allein; ihre Gegenprobe ist ein Ersatzskript.
kopie gate
find "$k/harness/mk" -name '*.mk' ! -name lint.mk -delete
printf 'exit 0\n' | datei tools/harness/lint-gegenprobe.sh
datei "$G/allein/fall.go" <<'EOF'
// Package allein trägt eine Direktive.
package allein

// Wert liefert eins.
func Wert() int {
	return 1 //nolint
}
EOF
starte gate make -s DOCKER_BUILD="$docker_build --progress=plain" gates

wait
for f in "$arbeit"/*.log; do sed -E 's/^#[0-9]+ [0-9]+\.[0-9]+ //' "$f" > "${f%.log}.n"; done

# --- Auswertung --------------------------------------------------------------
[ "$(cat "$arbeit/grundlauf.exit")" -eq 0 ] \
  || melde "Fall 'p0-grundlauf': erwartet Ausgang 0, bekam $(cat "$arbeit/grundlauf.exit"): $(grep -m5 -E '^(lint: |ERROR)' "$(n grundlauf)")"
keine_lint_zeile p0-grundlauf grundlauf

S=sammel
SW="$G/linter/schwellen.go"
RF="$R"
rot p0-sammel "$S"
# (1) Gegenstand.
befund p1-integration-build-tag "$S" test/integration/gegenprobe_e2e_test.go 'var gegenprobeIntegration' gochecknoglobals
befund p1-cmd "$S" cmd/gegenprobe/main.go 'fmt.Println(' forbidigo
befund p1-test "$S" test/gegenprobe/fall.go 'var testGlobal' gochecknoglobals
befund p1-pfad-anker-global "$S" test/anker/internal/hexagon/model/fehler.go 'var klassen' gochecknoglobals 'klassen is a global variable'
befund p1-pfad-anker-forbidigo "$S" "$G/test/fall.go" 'os.Stdout' forbidigo
# (3) Linter.
befund p3-errcheck "$S" "$L" 'func errcheckFall' errcheck
befund p3-govet "$S" "$L" 'func govetFall' govet
befund p3-ineffassign "$S" "$L" 'func ineffassignFall' ineffassign
befund p3-staticcheck "$S" "$L" 'func staticcheckFall' staticcheck 'SA1018'
befund p3-unused "$S" "$L" 'func unusedFall' unused
befund p3-containedctx "$S" "$L" 'type containedctxFall' containedctx
befund p3-fatcontext "$S" "$L" 'ctx = context.WithValue' fatcontext
befund p3-forbidigo "$S" "$L" 'fmt.Print("print")' forbidigo
befund p3-gochecknoglobals "$S" "$L" 'var gochecknoglobalsFall' gochecknoglobals
befund p3-gochecknoinits "$S" "$L" 'func init()' gochecknoinits
befund p3-gomodguard "$S" "$L" 'golang.org/x/sync/errgroup' gomodguard_v2
befund p3-iface "$S" "$L" 'type ifaceErste' iface
befund p3-inamedparam "$S" "$L" 'type inamedparamFall' inamedparam
befund p3-ireturn "$S" "$L" 'func ireturnFall' ireturn
befund p3-noctx "$S" "$L" 'func noctxFall' noctx
befund p3-reassign "$S" "$L" 'func reassignFall' reassign
befund p3-unparam "$S" "$L" 'func unparamLeer' unparam
# (4) Schwellen: grün an der Schwelle, rot einen Schritt darüber.
kein_befund p4-cyclop-gruen "$S" "$SW" 'func zyklomatisch15(' cyclop
befund p4-cyclop-rot "$S" "$SW" 'func zyklomatisch16(' cyclop 'is 16, max is 15'
kein_befund p4-gocyclo-gruen "$S" "$SW" 'func zyklomatisch15(' gocyclo
befund p4-gocyclo-rot "$S" "$SW" 'func zyklomatisch16(' gocyclo 'complexity 16 '
kein_befund p4-gocognit-gruen "$S" "$SW" 'func kognitiv20(' gocognit
befund p4-gocognit-rot "$S" "$SW" 'func kognitiv21(' gocognit 'complexity 21 '
kein_befund p4-maintidx-gruen "$S" "$SW" 'func wartbar20(' maintidx
befund p4-maintidx-rot "$S" "$SW" 'func wartbar19(' maintidx 'Maintainability Index: 19'
kein_befund p4-nestif-gruen "$S" "$SW" '// verschachtelt4' nestif
befund p4-nestif-rot "$S" "$SW" '// verschachtelt5' nestif 'complexity: 5'
kein_befund p4-funlen-anweisungen-gruen "$S" "$SW" 'func anweisungen60(' funlen
befund p4-funlen-anweisungen-rot "$S" "$SW" 'func anweisungen61(' funlen '61 > 60'
kein_befund p4-funlen-zeilen-gruen "$S" "$SW" 'func zeilen100(' funlen
befund p4-funlen-zeilen-rot "$S" "$SW" 'func zeilen101(' funlen '101 > 100'
kein_befund p4-dupl-gruen "$S" "$SW" 'func doppeltKurzA(' dupl
befund p4-dupl-rot "$S" "$G/doppelt/fall.go" 'func langA(' dupl
kein_befund p4-interfacebloat-gruen "$S" "$SW" 'type breit10 interface' interfacebloat
befund p4-interfacebloat-rot "$S" "$SW" 'type breit11 interface' interfacebloat
# (5) Einstellungen.
for a in Fprint Fprintf Fprintln; do
  kein_befund "p5-errcheck-frei-$a" "$S" "$L" "fmt.$a(w" errcheck
done
kein_befund p5-errcheck-frei-conn "$S" "$L" 'c.Close()' errcheck
kein_befund p5-errcheck-frei-listener "$S" "$L" 'l.Close()' errcheck
befund p5-errcheck-rot "$S" "$L" 'w.Write(nil)' errcheck
for a in Error Leer Anonym Stdlib Generisch Ports Pgproto3; do
  kein_befund "p5-ireturn-frei-$a" "$S" "$L" "func ireturn$a" ireturn
done
for a in 'fmt.Printf(' 'fmt.Println(' 'print("builtin' 'println("builtin' '_ = os.Stdout' '_ = os.Stderr'; do
  befund "p5-forbidigo-$a" "$S" "$L" "$a" forbidigo
done
kein_befund p5-forbidigo-cmd "$S" cmd/gegenprobe/main.go 'os.Stdout, os.Stderr' forbidigo
kein_befund p5-forbidigo-test-stdout "$S" test/gegenprobe/fall.go 'var testGlobal' forbidigo
kein_befund p5-forbidigo-test-stderr "$S" test/gegenprobe/fall.go 'os.Stderr' forbidigo
A=test/gegenprobe/ausnahme_test.go
kein_befund p8-ausnahme-test-cyclop "$S" "$A" 'func testZyklomatisch16(' cyclop
kein_befund p8-ausnahme-test-gocyclo "$S" "$A" 'func testZyklomatisch16(' gocyclo
kein_befund p8-ausnahme-test-gocognit "$S" "$A" 'func testKognitiv21(' gocognit
kein_befund p8-ausnahme-test-nestif "$S" "$A" '// testVerschachtelt5' nestif
kein_befund p8-ausnahme-test-funlen-anweisungen "$S" "$A" 'func testAnweisungen61(' funlen
kein_befund p8-ausnahme-test-funlen-zeilen "$S" "$A" 'func testZeilen101(' funlen
kein_befund p8-ausnahme-test-noctx "$S" "$A" 'func testNoctx' noctx
kein_befund p8-ausnahme-test-unparam "$S" "$A" 'func testUnparam' unparam
kein_befund p8-ausnahme-test-unused-parameter "$S" "$A" 'func testParameter' revive
kein_befund p8-ausnahme-test-unused-receiver "$S" "$A" 'func (t testTyp)' revive
kein_befund p8-ausnahme-st1005 "$S" "$L" 'func st1005Frei' staticcheck
kein_befund p5-gomodguard-frei-pgx "$S" "$L" 'pgx/v5/pgproto3' gomodguard_v2
kein_befund p5-gomodguard-frei-yaml "$S" "$L" 'go.yaml.in/yaml/v3' gomodguard_v2
befund p5-revive-blank-imports "$S" "$RF" '_ "strings"' revive 'blank-imports'
befund p5-revive-context-as-argument "$S" "$RF" 'func kontextHinten' revive 'context-as-argument'
befund p5-revive-context-keys-type "$S" "$RF" 'func schluessel' revive 'context-keys-type'
befund p5-revive-dot-imports "$S" "$RF" '. "unicode"' revive 'dot-imports'
befund p5-revive-empty-block "$S" "$RF" 'func leer' revive 'empty-block'
befund p5-revive-error-naming "$S" "$RF" 'var fehlerWert' revive 'error-naming'
befund p5-revive-error-return "$S" "$RF" 'func fehlerVorn' revive 'error-return'
befund p5-revive-error-strings "$S" "$RF" 'func punktAmEnde' revive 'error-strings'
befund p5-revive-errorf "$S" "$RF" 'func formatiert' revive 'errorf'
befund p5-revive-exported "$S" "$RF" 'func OhneKommentar' revive 'exported'
befund p5-revive-exported-form "$S" "$RF" '// liefert nichts.' revive 'exported'
befund p5-revive-if-return "$S" "$RF" 'func weiter' revive 'if-return'
befund p5-revive-increment-decrement "$S" "$RF" 'func erhoehe' revive 'increment-decrement'
befund p5-revive-indent-error-flow "$S" "$RF" 'func sonst' revive 'indent-error-flow'
befund p5-revive-package-comments "$S" "$G/paketkommentar/fall.go" 'package paketkommentar' revive 'package-comments'
befund p5-revive-range "$S" "$RF" 'func bereich' revive 'range'
befund p5-revive-receiver-naming "$S" "$RF" 'func (this typ)' revive 'receiver-naming'
befund p5-revive-redefines-builtin-id "$S" "$RF" 'func eingebaut' revive 'redefines-builtin-id'
befund p5-revive-superfluous-else "$S" "$RF" 'func schleife' revive 'superfluous-else'
befund p5-revive-time-naming "$S" "$RF" 'func zeit' revive 'time-naming'
befund p5-revive-unexported-return "$S" "$RF" 'func Liefert' revive 'unexported-return'
befund p5-revive-unreachable-code "$S" "$RF" '// toter' revive 'unreachable-code'
befund p5-revive-unused-parameter "$S" "$RF" 'func unbenutzt' revive 'unused-parameter'
befund p5-revive-unused-receiver "$S" "$RF" 'func (t typ)' revive 'unused-receiver'
befund p5-revive-var-declaration "$S" "$RF" 'func deklaration' revive 'var-declaration'
befund p5-revive-var-naming "$S" "$RF" 'func benennung' revive 'var-naming'
befund p5-contextcheck-test "$S" "$G/kontext/fall_test.go" '// kontext-fall' contextcheck
for a in 'var version' 'var _ = 2' 'var ErrFrei'; do
  kein_befund "p5-gochecknoglobals-frei-$a" "$S" "$L" "$a" gochecknoglobals
done
while read -r pfad name; do
  d="${pfad#^}"; d="${d%\$}"; d="${d//\\./.}"
  befund "p8-ausnahme-name-$name" "$S" "$d" "var gegenprobe$name =" gochecknoglobals "gegenprobe$name is a global"
done <<<"$namen"
# (6) Kein `//nolint`: je Schreibweise eine Zeile.
nolint_rot() { lint_zeile "$1" "$S" "lint: $2:$(zeile_von "$S" "$2" "$3"): Direktive nolint"; }
nolint_rot p6-doppelstrich "$NF" 'a := 1 //nolint'
nolint_rot p6-leerzeichen-zusatz "$NF" '// nolint:errcheck'
nolint_rot p6-block "$NF" '/* nolint */'
nolint_rot p6-schreibweise "$NF" '//NoLint'
nolint_rot p6-drei-striche "$NF" '///nolint'
nolint_rot p6-leerzeichen-dann-strich "$NF" '// /nolint'
nolint_rot p6-tab "$NF" $'//\tnolint'
nolint_rot p6-zweites-zeichen "$NF" '// x //nolint'
nolint_rot p6-string "$NF" '_ = "//nolint"'
nolint_rot p6-cmd cmd/gegenprobe/main.go 'nolint-zeile'
nolint_rot p6-test test/gegenprobe/fall.go 'nolint-zeile'
nolint_rot p6-testdatei "$G/nolint/fall_test.go" 'func TestDirektive'
for a in 'siehe-nolint:// siehe nolint' 'wortgrenze://nolintx'; do
  keine_lint_zeile "p6-${a%%:*}" "$S" "lint: $NF:$(zeile_von "$S" "$NF" "${a#*:}"): Direktive nolint"
done
# (7) Export-Test-Brücke.
befund p7-whitebox-kern "$S" internal/hexagon/model/weiss_test.go 'package model' testpackage
befund p7-whitebox-driven "$S" internal/adapters/driven/recording/weiss_test.go 'package recording' testpackage
befund p7-whitebox-pgwire "$S" internal/adapters/driving/pgwire/weiss_test.go 'package pgwire' testpackage
befund p7-whitebox-einstieg "$S" internal/adapters/driving/cli/weiss_test.go 'package cli' testpackage
befund p7-internal-test "$S" "$B/internal_test.go" 'package bruecke' testpackage
befund p7-endet-auf-export-test "$S" "$B/fooexport_test.go" 'package bruecke' testpackage
kein_befund p7-bruecke-kein-testpackage "$S" "$B/export_test.go" 'package bruecke' testpackage
for a in TestBruecke BenchmarkBruecke ExampleIntern FuzzBruecke; do
  lint_zeile "p7-test-in-bruecke-$a" "$S" "lint: $B/export_test.go:$(zeile_von "$S" "$B/export_test.go" "func $a"): Testfunktion in der Brücke export_test.go"
done
lint_zeile p7-test-in-bruecke-tab "$S" "lint: $B/export_test.go:$(zeile_von "$S" "$B/export_test.go" $'func\tTestTab'): Testfunktion in der Brücke export_test.go"
keine_lint_zeile p7-bruecke-weiterreichen "$S" "lint: $B/export_test.go:$(zeile_von "$S" "$B/export_test.go" 'func Intern'): Testfunktion in der Brücke export_test.go"
befund p7-variable-in-bruecke "$S" "$B/export_test.go" 'var Zustand' gochecknoglobals
# (8) Ausnahmen.
warum="Regel ohne Kommentarblock \"# Why:\" unmittelbar darüber"
lint_zeile p8-ohne-why "$S" "lint: $P:$((n0 + 2)): $warum"
lint_zeile p8-why-nicht-erste "$S" "lint: $P:$((n0 + 8)): $warum"
lint_zeile p8-why-leerzeile "$S" "lint: $P:$((n0 + 14)): $warum"
lint_zeile p8-why-ohne-doppelpunkt "$S" "lint: $P:$(($(grep -n '^      # Why kein Doppelpunkt\.$' "$arbeit/$S/$P" | cut -d: -f1) + 1)): $warum"
keine_lint_zeile p8-why-vorhanden "$S" "lint: $P:$((n0 + 19)): $warum"
lint_zeile p8-ungenutzt-feldfolge "$S" 'lint: .golangci.yml: Regel ohne Befund: Linter: revive, Pfad: ^internal/gegenprobe/nirgends/, Text: ^gegenprobe\.text"$, Quelle: ^gegenprobe-quelle$'
lint_zeile p8-ungenutzt-pfad-ausser "$S" 'lint: .golangci.yml: Regel ohne Befund: Linter: revive, Pfad außer: ^internal/gegenprobe/, Text: ^gegenprobe-ausser$'
form() { lint_zeile "$1" "$2" "lint: $P:$3: Form nicht erkannt"; }
keine_form() { keine_lint_zeile "$1" "$2" "lint: $P:$3: Form nicht erkannt"; }
keine_form p8-kommentar-nach-exclusions "$S" "$exkl"
keine_form p8-kommentar-nach-rules "$S" "$regeln"
keine_form p8-schluessel-nach-rules "$S" "$(grep -n '^    paths: \[\]$' "$arbeit/$S/$P" | cut -d: -f1)"
rot p8-form-einzug form-einzug
form p8-form-einzug form-einzug "$exkl"
rot p8-form-andere-stelle form-andere-stelle
form p8-form-andere-stelle form-andere-stelle "$andere"
rot p8-form-fluss-exclusions form-fluss-exclusions
form p8-form-fluss-exclusions form-fluss-exclusions "$exkl"
rot p8-form-fluss-rules form-fluss-rules
form p8-form-fluss-rules form-fluss-rules "$regeln"
rot p8-form-eintrag-8 form-eintrag-8
soll="$(awk -v r="$regeln" -v e="$ende" -v p="$P" 'NR > r && NR <= e && !/^[ \t]*(#|$)/ { print "lint: " p ":" NR ": Form nicht erkannt" }' "$arbeit/form-eintrag-8/$P")"
ist="$(grep '^lint: ' "$(n form-eintrag-8)" || true)"
[ "$ist" = "$soll" ] || melde "Fall 'p8-form-eintrag-8' erwartet:"$'\n'"$soll"$'\n'"bekam:"$'\n'"$ist"
rot p8-form-eintrag-2 form-eintrag-2
form p8-form-eintrag-2 form-eintrag-2 "$erste"
rot p8-form-fortsetzung-7 form-fortsetzung-7
form p8-form-fortsetzung-7 form-fortsetzung-7 "$fortsetzung"
rot p8-form-ohne-strich form-ohne-strich
form p8-form-ohne-strich form-ohne-strich "$((regeln + 1))"
rot p8-form-quote-exclusions form-quote-exclusions
form p8-form-quote-exclusions form-quote-exclusions "$exkl"
rot p8-form-quote-rules form-quote-rules
form p8-form-quote-rules form-quote-rules "$regeln"
rot p8-form-fluss-anderswo form-fluss-anderswo
form p8-form-fluss-anderswo form-fluss-anderswo "$(wc -l < "$arbeit/form-fluss-anderswo/$P")"
rot p8-form-listenpunkt form-listenpunkt
form p8-form-listenpunkt form-listenpunkt "$(wc -l < "$arbeit/form-listenpunkt/$P")"
[ "$(cat "$arbeit/form-reihenfolge.exit")" -eq 0 ] || melde "Fall 'p8-reihenfolge': erwartet Ausgang 0"
keine_lint_zeile p8-reihenfolge form-reihenfolge
for c in leerraum strich-allein; do
  [ "$(cat "$arbeit/form-$c.exit")" -eq 0 ] || melde "Fall 'p8-$c': erwartet Ausgang 0"
  keine_lint_zeile "p8-$c" "form-$c"
done
ist="$(grep '^lint: ' "$(n form-strich-ohne-why)" || true)"
[ "$ist" = "lint: $P:$strich_zeile: $warum" ] \
  || melde "Fall 'p8-strich-allein-ohne-why': erwartet genau die Why-Zeile an $strich_zeile, bekam: $ist"
rot p8-strich-allein-ohne-why form-strich-ohne-why
rot p8-form-strich-allein-8 form-strich-allein-8
form p8-form-strich-allein-8 form-strich-allein-8 "$((regeln + 1))"
# (9) Ausgabe und Ausgang.
# allein <fall> <kopie> <zeile> — genau diese `lint:`-Zeile, golangci-lint ohne
# Befund, die Stufe mit Ausgang 1.
allein() {
  local ist
  ist="$(grep '^lint: ' "$(n "$2")" || true)"
  [ "$ist" = "$3" ] || melde "Fall '$1': erwartet genau die Zeile '$3', bekam: $ist"
  rot "$1" "$2"
}
allein p9-fehlt-allein fehlt-allein 'lint: .golangci.yml: fehlt'
allein p9-schema-allein schema-allein 'lint: .golangci.yml: von golangci-lint abgelehnt'
grep -qx '0 issues\.' "$(n schema-allein)" || melde "Fall 'p9-schema-allein': golangci-lint meldet etwas"
allein p9-ungenutzt-allein ungenutzt-allein 'lint: .golangci.yml: Regel ohne Befund: Linter: revive, Pfad: ^internal/gegenprobe/nirgends/'
grep -qx '0 issues\.' "$(n ungenutzt-allein)" || melde "Fall 'p9-ungenutzt-allein': golangci-lint meldet etwas"
allein p9-bruecke-allein bruecke-allein "lint: $G/bruecke/export_test.go:5: Testfunktion in der Brücke export_test.go"
grep -qx '0 issues\.' "$(n bruecke-allein)" || melde "Fall 'p9-bruecke-allein': golangci-lint meldet etwas"
allein p9-form-allein form-andere-stelle "lint: $P:$andere: Form nicht erkannt"
grep -qx '0 issues\.' "$(n form-andere-stelle)" || melde "Fall 'p9-form-allein': golangci-lint meldet etwas"
rot p9-profil-fehlt profil-fehlt
lint_zeile p9-profil-fehlt profil-fehlt 'lint: .golangci.yml: fehlt'
lint_zeile p9-profil-fehlt-nolint profil-fehlt "lint: $G/fehlt/fall.go:7: Direktive nolint"
lint_zeile p9-profil-fehlt-bruecke profil-fehlt "lint: $G/fehlt/export_test.go:3: Testfunktion in der Brücke export_test.go"
keine_lint_zeile p9-profil-fehlt-ohne-schema profil-fehlt 'lint: .golangci.yml: von golangci-lint abgelehnt'
if grep -q "^$G/fehlt/fall\\.go:5:" "$(n profil-fehlt)"; then
  melde "Fall 'p9-profil-fehlt-ohne-golangci': golangci-lint lief ohne Profil"
fi
lint_zeile p9-schema "$S" 'lint: .golangci.yml: von golangci-lint abgelehnt'
grep -v '^lint: ' "$(n "$S")" | grep -F 'gegenprobe-unbekannt' >/dev/null \
  || melde "Fall 'p9-schema-meldung': die Meldung von config verify fehlt"
for i in 1 2 3 4; do
  befund "p9-max-same-issues-$i" "$S" "$L" "os.Chdir(\"gleich-$i\")" errcheck
done
for i in $(seq 1 51); do
  befund "p9-max-issues-per-linter-$i" "$S" "$SW" "var viele$i =" gochecknoglobals
done
befund p9-uniq-by-line-a "$S" "$L" 'var zeileA' gochecknoglobals 'zeileA is a global'
befund p9-uniq-by-line-b "$S" "$L" 'var zeileA' gochecknoglobals 'zeileB is a global'
[ "$(cat "$arbeit/laden.exit")" -ne 0 ] || melde "Fall 'p9-laden': erwartet Ausgang ungleich 0"
keine_lint_zeile p9-laden laden
grep -qF "can't load config" "$(n laden)" || melde "Fall 'p9-laden': die Meldung von golangci-lint fehlt"
rot p9-golangci-allein golangci-allein
keine_lint_zeile p9-golangci-allein golangci-allein
befund p9-golangci-allein "golangci-allein" "$G/allein/fall.go" 'func init()' gochecknoinits
ist="$(grep '^lint: ' "$(n gate)" || true)"
[ "$ist" = "lint: $G/allein/fall.go:6: Direktive nolint" ] \
  || melde "Fall 'p9-lint-zeile-allein': erwartet genau die Zeile der Direktive, bekam: $ist"
grep -qx '0 issues\.' "$(n gate)" || melde "Fall 'p9-lint-zeile-allein': golangci-lint meldet etwas"
rot p9-lint-zeile-allein gate
# (10) Gate.
[ "$(cat "$arbeit/gate.exit")" -ne 0 ] || melde "Fall 'p10-gates-rot': make gates endet bei rotem make lint mit Exit 0"
datenbank="$(make -C "$arbeit/gate" -pq help 2>/dev/null || true)"
for ziel in lint lint-gegenprobe; do
  if ! printf '%s\n' "$datenbank" | awk -v z="$ziel" '/^GATE_CHECKS :?= / { for (i = 3; i <= NF; i++) if ($i == z) t = 1 } END { exit !t }'; then
    melde "Fall 'p10-gate-checks': $ziel hängt nicht an GATE_CHECKS"
  fi
done
printf 'exit 1\n' > "$arbeit/gate/tools/harness/lint-gegenprobe.sh"
if make -s -C "$arbeit/gate" lint-gegenprobe >/dev/null 2>&1; then
  melde "Fall 'p10-gegenprobe-scharf': make lint-gegenprobe endet bei rotem Skript mit Exit 0"
fi
# Grenze.
kein_befund neg-main-test "$S" cmd/gegenprobe/main_test.go 'package main' testpackage
kein_befund neg-generiert "$S" "$G/generiert/fall.go" 'var generiertGlobal' gochecknoglobals

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "lint-gegenprobe: gruen — Grundlauf ohne Befund, jeder Fall je Punkt von SPEC-049 wie erwartet"
