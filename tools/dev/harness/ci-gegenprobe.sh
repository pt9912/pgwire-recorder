#!/usr/bin/env bash
# ci-gegenprobe — prüft die Zusagen von SPEC-064 Teil 1 und 2 (harness/vertraege/SPEC-064-ci.md, Punkt 14): den Text von .github/workflows/ci.yml
# (Fälle p1-*, geprüft von der Funktion pruefe_ci, die den Text hält) und tools/harness/ci-range.sh (Fälle p6-* bis p9-*, kleine Git-Repositories in
# Temp-Verzeichnissen mit einem Remote-Zweig refs/remotes/origin/main und einem flachen Klon). Ein Fall besteht, wenn Exit-Code, stdout und stderr
# des Skripts vollständig dem erwarteten Text entsprechen bzw. pruefe_ci das erwartete Urteil fällt; jeder Fall hat seinen Temp-Baum. Je Mutation
# (Textersatz in pruefe_ci bzw. in einer Kopie von ci-range.sh) müssen genau die benannten Fälle rot werden; die unveränderte Kopie besteht jeden Fall.
# Der Name jedes Falls trägt die Nummer des Punktes; ein Name in der Fallliste ohne Funktion f_<name> zählt als roter Fall.
# Abbruch (Punkt 18): Bei Ende, SIGINT und SIGTERM beendet das Skript erst seine Nachkommen und wartet auf sie, entfernt dann sein Temp-Verzeichnis und
# endet mit 130 bzw. 143; die Fälle p18-sigint und p18-sigterm starten dafür eine Kopie des Skripts mit eigenem TMPDIR (die Kopie startet ohne GP_MUT,
# GP_ENTDECKE und GP_FALL und mit GP_PARALLEL=2, und GP_INNEN=1 schaltet in jedem ihrer Läufe die p18-Fälle und die Mutationen m-p18-* ab; ein Lauf mit
# GP_INNEN=1 und GP_MUT=m-p18 fährt nichts und endet mit Exit 1, Lauf `rekursion-schutz`) und senden das Signal, sobald dort ein Fall läuft; Mutationen
# der Kopie (Textersatz im Skript): die Falle fehlt, die Nachkommen werden nicht beendet, das Verzeichnis wird vor dem Ende der Nachkommen entfernt.
# SIGKILL und SIGHUP sind nicht zugesagt.
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Läufe, deren Name passt (die unveränderte Kopie nur, wenn das Muster `ok` trifft); GP_ENTDECKE=1 nennt je
# Mutation die roten Fälle; GP_PARALLEL=<n> (Vorgabe 4) setzt die Zahl der gleichzeitigen Läufe; GP_FALL=<name> gibt für die unveränderte Kopie das
# Protokoll dieses Falls aus.
# Braucht bash (mit `wait -n`), git, pgrep, setsid, timeout und die Werkzeuge von coreutils (mktemp, sort, tr, wc, cat, cp, chmod, rm, env, find, ls, seq,
# sleep, mkdir, dirname, basename), GNU sed und grep; ohne Docker und ohne make. Eine Version prüft das Skript nicht, zugesagt ist die Plattform von SPEC-064.
# Kein Fall ändert Dateirechte. Aufrufarten von `make` (-n, -q, -t, -i, MAKEFLAGS) haben keinen Fall: kein Hook und kein Stand liest das
# Ergebnis, und ci-range.sh läuft nie unter `make` (SPEC-064 Punkt 14).
set -uo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="$(cd "$hier/../../.." && pwd)"
arbeit="$(mktemp -d)"
selbst="$hier/$(basename "${BASH_SOURCE[0]}")"
raeume_auf() { chmod -R u+rwX "${arbeit:?}" 2> /dev/null; rm -rf "${arbeit:?}"; }
# nachkommen <pid> — alle Nachkommen, Eltern vor Kindern.
nachkommen() { local k; for k in $(pgrep -P "$1" 2> /dev/null); do echo "$k"; nachkommen "$k"; done; }
# beende_kinder — hält die Nachkommen an (bis die Menge feststeht), beendet sie und wartet, bis keiner mehr läuft.
beende_kinder() {
  local i p="" n k
  for i in 1 2 3 4 5; do
    n="$(nachkommen $$ | sort -u | tr '\n' ' ')"
    [ "$n" != "$p" ] || break
    p="$n"
    # shellcheck disable=SC2086
    kill -s STOP $p 2> /dev/null
  done
  # shellcheck disable=SC2086
  [ -z "$p" ] || { kill -s TERM $p 2> /dev/null; kill -s CONT $p 2> /dev/null; }
  for k in $p; do for i in $(seq 1 100); do kill -0 "$k" 2> /dev/null || break; sleep 0.1; done; done
  wait 2> /dev/null
}
abbruch_int() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 130; }
abbruch_term() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 143; }
trap raeume_auf EXIT
trap abbruch_int INT
trap abbruch_term TERM
: > "$arbeit/fehler"
: > "$arbeit/gefahren"
melde() { echo "ci-gegenprobe: ROT — $*" >> "$arbeit/fehler"; }

for t in git sed grep mktemp sort tr wc cat cp chmod rm env find ls seq sleep mkdir dirname basename pgrep setsid timeout; do
  command -v "$t" > /dev/null 2>&1 || { echo "ci-gegenprobe: $t fehlt" >&2; exit 2; }
done
for f in .github/workflows/ci.yml tools/harness/ci-range.sh; do
  [ -f "$repo/$f" ] || { echo "ci-gegenprobe: $f fehlt im Arbeitsbaum" >&2; exit 2; }
done
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1
export GIT_AUTHOR_NAME=gp GIT_AUTHOR_EMAIL=gp@example.invalid GIT_COMMITTER_NAME=gp GIT_COMMITTER_EMAIL=gp@example.invalid

# --- Punkt 1: der Text -------------------------------------------------------------------------------------------------------------------------------
# SOLL ist der Text von Punkt 1, Zeile für Zeile; die Zeile mit <SHA> und <Tag> prüft PIN_RE. Die Gegenprobe trägt den Text wörtlich als eigene Kopie.
mapfile -t SOLL << 'EOT'
name: ci
on:
  pull_request:
  push:
    branches:
      - '**'
    tags-ignore:
      - '**'
permissions: {}
jobs:
  gates:
    name: ci-gates
    runs-on: ubuntu-24.04
    timeout-minutes: 30
    permissions:
      contents: read
    env:
      CI_EVENT: ${{ github.event_name }}
      CI_HEAD: ${{ github.sha }}
      CI_PR_BASE: ${{ github.event.pull_request.base.sha }}
      CI_PUSH_BEFORE: ${{ github.event.before }}
      CI_DEFAULT_BRANCH: ${{ github.event.repository.default_branch }}
    steps:
      - name: Checkout
        uses: actions/checkout@<SHA> # <Tag>
        with:
          fetch-depth: 0
          persist-credentials: false
      - name: Commit-Messages der Range
        run: |
          range="$(bash tools/harness/ci-range.sh)"
          make doc-commits RANGE="$range"
      - name: ADR-Immutabilität der Range
        if: ${{ !cancelled() }}
        run: |
          range="$(bash tools/harness/ci-range.sh)"
          make doc-immutable RANGE="$range"
      - name: Gates
        run: make gates
EOT
N=${#SOLL[@]}
PIN=-1
for ((i = 0; i < N; i++)); do if [[ ${SOLL[i]} == *"actions/checkout@<SHA>"* ]]; then PIN=$i; fi; done
PIN_RE="^${SOLL[PIN]%%<SHA>*}[0-9a-f]{40} # [^ ].*\$"
SHA_PROBE=0123456789abcdef0123456789abcdef01234567
# Anfang der drei Schritte nach dem Checkout (Zeilen in SOLL)
for ((i = 0; i < N; i++)); do
  case "${SOLL[i]}" in
    *"name: Commit-Messages der Range") S_C=$i ;;
    *"name: ADR-Immutabilität der Range") S_A=$i ;;
    *"name: Gates") S_G=$i ;;
  esac
done

# pruefe_ci <datei> — Exit 0, wenn die Datei außer Leerzeilen und Kommentarzeilen genau aus den Zeilen von SOLL besteht (Pin der Zeile PIN nach PIN_RE).
pruefe_ci() {
  local datei="$1" z i j=0
  local -a ist=()
  while IFS= read -r z || [ -n "$z" ]; do
    [[ $z =~ ^\ *(#|$) ]] && continue
    ist+=("$z")
  done < "$datei"
  for ((i = 0; i < ${#SOLL[@]}; i++)); do
    if [ "$i" -eq "$PIN" ]; then
      [[ ${ist[j]-} =~ $PIN_RE ]] || return 1
    else
      [ "${ist[j]-}" = "${SOLL[i]}" ] || return 1
    fi
    j=$((j + 1))
  done
  [ "${#ist[@]}" -eq "$j" ]
}

# --- Bausteine der Fälle ---------------------------------------------------------------------------------------------------------------------------
# text — der Text von Punkt 1 mit einem gültigen Pin, eine Zeile je Zeile.
text() {
  local i
  for i in "${!SOLL[@]}"; do
    if [ "$i" -eq "$PIN" ]; then echo "${SOLL[i]/<SHA> # <Tag>/$SHA_PROBE # v7.0.1}"; else echo "${SOLL[i]}"; fi
  done
}
# urteil <ok|rot> — schreibt stdin nach $w/ci.yml und vergleicht das Urteil von pruefe_ci.
urteil() {
  local soll="$1" r=0
  mkdir -p "$w"
  cat > "$w/ci.yml"
  pruefe_ci "$w/ci.yml" || r=$?
  if { [ "$soll" = ok ] && [ "$r" -ne 0 ]; } || { [ "$soll" = rot ] && [ "$r" -eq 0 ]; }; then
    bad=1
    printf 'erwartet: %s\ngelesen: Exit %s\n' "$soll" "$r" >> "$run/log.$fallname"
  fi
}
neu() { w="$run/t.$fallname"; mapfile -t t < <(text); }
zeilen() { printf '%s\n' "$@"; }
f_p1_gleich() { neu; mkdir -p "$w"; pruefe_ci "$repo/.github/workflows/ci.yml" || { bad=1; echo "ci.yml des Arbeitsbaums weicht von Punkt 1 ab" >> "$run/log.$fallname"; }; }
f_p1_text_gleich() { neu; urteil ok < <(zeilen "${t[@]}"); }
p1_fehlt() { neu; unset "t[$(($1 - 1))]"; urteil rot < <(zeilen "${t[@]}"); }
p1_geaendert() {
  local l lead
  neu
  l="${t[$(($1 - 1))]}"; lead="${l%%[! ]*}"
  t[$(($1 - 1))]="${lead}Z${l:${#lead}+1}"
  urteil rot < <(zeilen "${t[@]}")
}
for ((n = 1; n <= N; n++)); do eval "f_p1_fehlt_$n() { p1_fehlt $n; }; f_p1_geaendert_$n() { p1_geaendert $n; }"; done
# Zusätze am Ende der Datei: <Wort, das die Mutation durchlässt>|<Zeile>
ZUSATZ=(
  "concurrency:|concurrency: ci"
  "workflow_dispatch:|workflow_dispatch:"
  "pull_request_target:|pull_request_target:"
  "workflow_run:|workflow_run:"
  "strategy:|    strategy:"
  "continue-on-error: true|        continue-on-error: true"
  'secrets.|      X: ${{ secrets.X }}'
  "actions/cache@|      - uses: actions/cache@$SHA_PROBE # v4"
  "actions/upload-artifact@|      - uses: actions/upload-artifact@$SHA_PROBE # v4"
)
ZUSATZ_NAMEN=(concurrency workflow-dispatch pull-request-target workflow-run strategy continue-on-error secrets schritt-cache schritt-artefakt)
zusatz() { neu; urteil rot < <(zeilen "${t[@]}" "$1"); }
for k in "${!ZUSATZ[@]}"; do eval "f_p1_zusatz_${ZUSATZ_NAMEN[k]//-/_}() { zusatz \"\${ZUSATZ[$k]#*|}\"; }"; done
pin() { neu; t[PIN]="        $1"; urteil rot < <(zeilen "${t[@]}"); }
f_p1_pin_kurz() { pin "uses: actions/checkout@${SHA_PROBE:1} # v7.0.1"; }
f_p1_pin_gross() { pin "uses: actions/checkout@${SHA_PROBE^^} # v7.0.1"; }
f_p1_pin_tag_leer() { pin "uses: actions/checkout@$SHA_PROBE # "; }
f_p1_pin_ohne_tag() { pin "uses: actions/checkout@$SHA_PROBE"; }
f_p1_pin_andere_action() { pin "uses: actions/cache@$SHA_PROBE # v4"; }
# umgestellt <C|A|G>… — der Text mit den Schritten in der genannten Reihenfolge (C: Commit-Messages, A: ADR-Immutabilität, G: Gates; der Checkout bleibt vorn).
umgestellt() {
  local b
  zeilen "${t[@]:0:S_C}"
  for b in "$@"; do
    case "$b" in
      C) zeilen "${t[@]:S_C:S_A-S_C}" ;;
      A) zeilen "${t[@]:S_A:S_G-S_A}" ;;
      G) zeilen "${t[@]:S_G}" ;;
    esac
  done
}
f_p1_reihenfolge_schritte() { neu; urteil rot < <(umgestellt A C G); }
f_p1_reihenfolge_gates() { neu; urteil rot < <(umgestellt G C A); }
f_p1_leerzeile() { neu; urteil ok < <(zeilen "${t[@]:0:3}" "" "${t[@]:3}"); }
f_p1_kommentarzeile() { neu; urteil ok < <(zeilen "${t[@]:0:3}" "# Kommentar" "${t[@]:3}"); }
f_p1_kommentarzeile_eingerueckt() { neu; urteil ok < <(zeilen "${t[@]:0:3}" "    # Kommentar" "${t[@]:3}"); }

# --- Punkte 6 bis 9: ci-range.sh ---------------------------------------------------------------------------------------------------------------------
# welt <dir> — legt ein Repository quelle/ (main: c1 c2 c3, Zweig side mit s1 aus c1, Merge-Commit m, Tag-Objekt tg auf c3; refs/remotes/origin/main auf s1: der Standard-Zweig lief nach dem Abzweig bei c1 weiter) und den flachen Klon flach/ an.
welt() {
  local q="$1/quelle"
  mkdir -p "$q"
  git -C "$q" init -q -b main
  git -C "$q" commit -q --allow-empty -m c1
  c1="$(git -C "$q" rev-parse HEAD)"
  git -C "$q" checkout -q -b side
  git -C "$q" commit -q --allow-empty -m s1
  s1="$(git -C "$q" rev-parse HEAD)"
  git -C "$q" checkout -q main
  git -C "$q" commit -q --allow-empty -m c2
  c2="$(git -C "$q" rev-parse HEAD)"
  git -C "$q" commit -q --allow-empty -m c3
  c3="$(git -C "$q" rev-parse HEAD)"
  git -C "$q" merge -q --no-ff --no-edit side
  m="$(git -C "$q" rev-parse HEAD)"
  git -C "$q" tag -a -m tag t3 "$c3"
  tg="$(git -C "$q" rev-parse t3)"
  git -C "$q" update-ref refs/remotes/origin/main "$s1"
  git clone -q --depth 1 "file://$q" "$1/flach"
  fh="$(git -C "$1/flach" rev-parse HEAD)"
  EINS=1111111111111111111111111111111111111111
  ZWEI=2222222222222222222222222222222222222222
  NULL=0000000000000000000000000000000000000000
}
# kopie <quelle|flach> — der Baum des Falls (eine Kopie der Welt).
kopie() { neu; rm -rf "$w"; mkdir -p "$w"; cp -a "$run/welt/$1/." "$w"; }
# umg <Ereignis> [NAME=Wert | -NAME]… — setzt E: alle fünf Variablen gültig (die unbenutzten als Köder), dann die Abweichungen.
umg() {
  local -A v=([CI_EVENT]="$1" [CI_HEAD]="$c3" [CI_PR_BASE]="$c1" [CI_PUSH_BEFORE]="$c2" [CI_DEFAULT_BRANCH]=main)
  local a
  shift
  for a in "$@"; do if [[ $a == -* ]]; then unset "v[${a#-}]"; else v[${a%%=*}]="${a#*=}"; fi; done
  E=()
  for a in "${!v[@]}"; do E+=("$a=${v[$a]}"); done
}
# lauf [Verzeichnis] — führt ci-range.sh mit E in w aus und setzt code, out, err (stdout und stderr unverändert, mit Zeilenende).
lauf() {
  ( cd "${1:-$w}" && env -i PATH="$PATH" GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 "${E[@]}" bash "$rangesh" ) > "$w.out" 2> "$w.err"
  code=$?
  out="$(cat "$w.out"; printf x)"; out="${out%x}"
  err="$(cat "$w.err"; printf x)"; err="${err%x}"
}
# erwarte <Exit> <stdout-Zeile oder ""> <stderr-Zeile oder "">… — stdout und stderr vollständig (jede Zeile mit Zeilenende).
erwarte() {
  local eo="" ee=""
  [ -z "$2" ] || eo="$2"$'\n'
  [ -z "$3" ] || ee="$3"$'\n'
  [ -z "${4-}" ] || ee+="$4"$'\n'
  if [ "$code" -ne "$1" ] || [ "$out" != "$eo" ] || [ "$err" != "$ee" ]; then
    bad=1
    printf 'erwartet: Exit %s\nstdout:\n%s\nstderr:\n%s\ngelesen: Exit %s\nstdout:\n%s\nstderr:\n%s\n' "$1" "$eo" "$ee" "$code" "$out" "$err" >> "$run/log.$fallname"
  fi
}
# erwarte_grund <Exit> <Teiltext> — Exit-Code, leere stdout und eine stderr, deren Zeilen alle mit `ci-range:` beginnen und die den Teiltext nennen (der Wortlaut der
# Meldungen ist nicht zugesagt, SPEC-064 Punkt 8; der Teiltext ist der Ort: Variable oder fetch-depth).
erwarte_grund() {
  local z ok=1
  [ "$code" -eq "$1" ] && [ -z "$out" ] && [ -n "$err" ] || ok=0
  [[ $err == *"$2"* ]] || ok=0
  while IFS= read -r z; do [[ $z == "ci-range: "* ]] || ok=0; done <<< "${err%$'\n'}"
  if [ "$ok" -eq 0 ]; then
    bad=1
    printf 'erwartet: Exit %s, stdout leer, stderr mit Präfix und %s\ngelesen: Exit %s\nstdout:\n%s\nstderr:\n%s\n' "$1" "$2" "$code" "$out" "$err" >> "$run/log.$fallname"
  fi
}
UNGUELTIG="fehlt oder ist ungültig"
# Punkt 6
f_p6_event_fehlt() { kopie quelle; umg pull_request -CI_EVENT; lauf; erwarte_grund 2 CI_EVENT; }
f_p6_event_leer() { kopie quelle; umg "" ; lauf; erwarte_grund 2 CI_EVENT; }
f_p6_event_unbekannt() { kopie quelle; umg workflow_dispatch; lauf; erwarte_grund 2 CI_EVENT; }
f_p6_event_schedule() { kopie quelle; umg schedule; lauf; erwarte_grund 2 CI_EVENT; }
FORMEN=(fehlt ungueltig kurz lang gross)
wert() { case "$1" in ungueltig) printf 'g%.0s' $(seq 1 40) ;; kurz) printf 'a%.0s' $(seq 1 39) ;; lang) printf 'a%.0s' $(seq 1 41) ;; gross) printf 'A%.0s' $(seq 1 40) ;; esac; }
ungueltige_variable() { # <Ereignis> <Variable> <Form>
  if [ "$3" = fehlt ]; then umg "$1" "-$2"; else umg "$1" "$2=$(wert "$3")"; fi
  kopie quelle; lauf; erwarte_grund 2 "$2"
}
for f in "${FORMEN[@]}"; do
  eval "f_p6_head_$f() { ungueltige_variable pull_request CI_HEAD $f; }
        f_p6_prbase_$f() { ungueltige_variable pull_request CI_PR_BASE $f; }
        f_p6_before_$f() { ungueltige_variable push CI_PUSH_BEFORE $f; }"
done
f_p6_default_fehlt() { umg push "CI_PUSH_BEFORE=$NULL" -CI_DEFAULT_BRANCH; kopie quelle; lauf; erwarte_grund 2 CI_DEFAULT_BRANCH; }
f_p6_default_ungueltig() { umg push "CI_PUSH_BEFORE=$NULL" CI_DEFAULT_BRANCH=a..b; kopie quelle; lauf; erwarte_grund 2 CI_DEFAULT_BRANCH; }
f_p6_default_vorgaenger_fehlt_fehlt() { umg push "CI_PUSH_BEFORE=$EINS" -CI_DEFAULT_BRANCH; kopie quelle; lauf; erwarte_grund 2 CI_DEFAULT_BRANCH; }
f_p6_default_vorgaenger_fehlt_ungueltig() { umg push "CI_PUSH_BEFORE=$EINS" CI_DEFAULT_BRANCH=a..b; kopie quelle; lauf; erwarte_grund 2 CI_DEFAULT_BRANCH; }
f_p6_ungenutzt() { umg pull_request CI_PUSH_BEFORE=x; kopie quelle; lauf; erwarte 0 "$c1..$c3" ""; }
f_p6_ungenutzt_prbase() { umg push CI_PR_BASE=x; kopie quelle; lauf; erwarte 0 "$c2..$c3" ""; }
f_p6_ungenutzt_default() { umg push CI_DEFAULT_BRANCH=a..b; kopie quelle; lauf; erwarte 0 "$c2..$c3" ""; }
# Punkt 7
f_p7_pr() { umg pull_request; kopie quelle; lauf; erwarte 0 "$c1..$c3" ""; }
f_p7_pr_merge_kopf() { umg pull_request "CI_HEAD=$m"; kopie quelle; lauf; erwarte 0 "$c1..$m" ""; }
f_p7_pr_basis_fehlt_im_klon() { umg pull_request "CI_PR_BASE=$EINS"; kopie quelle; lauf; erwarte_grund 2 ""; }
f_p7_pr_kopf_fehlt_im_klon() { umg pull_request "CI_HEAD=$ZWEI"; kopie quelle; lauf; erwarte_grund 2 ""; }
f_p7_push() { umg push; kopie quelle; lauf; erwarte 0 "$c2..$c3" ""; }
f_p7_push_neu() { umg push "CI_PUSH_BEFORE=$NULL"; kopie quelle; lauf; erwarte 0 "$c1..$c3" ""; }
f_p7_pr_tag_objekt() { umg pull_request "CI_HEAD=$tg"; kopie quelle; lauf; erwarte 0 "$c1..$c3" ""; }
# Zweig ohne gemeinsamen Vorfahr mit dem Standard-Zweig (entkoppelte Historie): Exit 2.
f_p7_push_neu_ohne_vorfahr() {
  local o
  kopie quelle
  git -C "$w" checkout -q --orphan fremd
  git -C "$w" commit -q --allow-empty -m fremd
  o="$(git -C "$w" rev-parse HEAD)"
  umg push "CI_PUSH_BEFORE=$NULL" "CI_HEAD=$o"; lauf; erwarte_grund 2 ""
}
f_p7_push_neu_ohne_standard() {
  umg push "CI_PUSH_BEFORE=$NULL"; kopie quelle; git -C "$w" update-ref -d refs/remotes/origin/main; lauf
  erwarte_grund 2 ""
}
f_p7_push_vorgaenger_fehlt() {
  umg push "CI_PUSH_BEFORE=$EINS"; kopie quelle; lauf
  erwarte 0 "$c1..$c3" "ci-range: Vorgänger $EINS nicht im Klon, Basis ist der Standard-Zweig"
}
f_p7_leer() { umg push "CI_PUSH_BEFORE=$NULL" "CI_HEAD=$c1"; kopie quelle; lauf; erwarte 1 "" "ci-range: Range $c1..$c1 ist leer"; }
f_p7_flach() { umg push "CI_PUSH_BEFORE=$fh" "CI_HEAD=$fh"; kopie flach; lauf; erwarte_grund 2 "fetch-depth: 0"; }
# Punkt 8
f_p8_stdout_eine_zeile() {
  umg push "CI_PUSH_BEFORE=$EINS"; kopie quelle; lauf
  [ "$(printf '%s' "$out" | wc -l)" -eq 1 ] && [ "$out" = "$c1..$c3"$'\n' ] && [ "$err" = "ci-range: Vorgänger $EINS nicht im Klon, Basis ist der Standard-Zweig"$'\n' ] || {
    bad=1; printf 'stdout:\n%s\nstderr:\n%s\n' "$out" "$err" >> "$run/log.$fallname"; }
}
# Jede Form mit Exit ungleich 0: stdout leer, jede Zeile von stderr mit dem Präfix.
f_p8_fehler_stdout_leer() {
  local form r
  for form in event kopf basis kopf-fehlt standard flach leer git; do
    case "$form" in
      event) umg pull_request -CI_EVENT; kopie quelle ;;
      kopf) umg pull_request -CI_HEAD; kopie quelle ;;
      basis) umg pull_request "CI_PR_BASE=$EINS"; kopie quelle ;;
      kopf-fehlt) umg pull_request "CI_HEAD=$ZWEI"; kopie quelle ;;
      standard) umg push "CI_PUSH_BEFORE=$NULL"; kopie quelle; git -C "$w" update-ref -d refs/remotes/origin/main ;;
      flach) umg push "CI_PUSH_BEFORE=$fh" "CI_HEAD=$fh"; kopie flach ;;
      leer) umg push "CI_PUSH_BEFORE=$NULL" "CI_HEAD=$c1"; kopie quelle ;;
      git) umg pull_request; kopie quelle; rm -rf "$w/.git" ;;
    esac
    lauf
    r=0
    [ "$code" -ne 0 ] && [ -z "$out" ] && [ -n "$err" ] || r=1
    [ -z "$(printf '%s' "$err" | grep -v '^ci-range: ')" ] || r=1
    if [ "$r" -ne 0 ]; then bad=1; printf 'Form %s: Exit %s\nstdout:\n%s\nstderr:\n%s\n' "$form" "$code" "$out" "$err" >> "$run/log.$fallname"; fi
  done
}
# Basis und Kopf sind aufzulösen, aber ein Commit dazwischen fehlt (Objekt gelöscht): git rev-list --count schlägt fehl.
f_p8_rev_list_fehler() {
  local o
  umg pull_request; kopie quelle
  o="$w/.git/objects/${c2:0:2}/${c2:2}"
  [ -f "$o" ] || { bad=1; echo "Objekt von c2 nicht lose: $o" >> "$run/log.$fallname"; return; }
  rm -f "$o"; lauf; erwarte_grund 2 ""
}
f_p8_git_fehler() { umg pull_request; kopie quelle; rm -rf "$w/.git"; lauf; erwarte_grund 2 ""; }
# Punkt 9
f_p9_ereignis_vor_head() { umg pull_request -CI_EVENT -CI_HEAD; kopie quelle; lauf; erwarte_grund 2 CI_EVENT; }
f_p9_head_vor_prbase() { umg pull_request -CI_HEAD -CI_PR_BASE; kopie quelle; lauf; erwarte_grund 2 CI_HEAD; }
f_p9_head_vor_before() { umg push -CI_HEAD -CI_PUSH_BEFORE; kopie quelle; lauf; erwarte_grund 2 CI_HEAD; }
f_p9_eingabe_vor_flach() { umg pull_request -CI_HEAD; kopie flach; lauf; erwarte_grund 2 CI_HEAD; }
f_p9_flach_vor_aufloesung() { umg pull_request "CI_PR_BASE=$EINS" "CI_HEAD=$fh"; kopie flach; lauf; erwarte_grund 2 "fetch-depth: 0"; }
f_p9_aufloesung_vor_leer() { umg pull_request "CI_PR_BASE=$EINS" "CI_HEAD=$EINS"; kopie quelle; lauf; erwarte_grund 2 ""; }

# --- Punkt 18: Abbruch ------------------------------------------------------------------------------------------------------------------------------
warte_auf() { local i; for i in $(seq 1 300); do if "$@"; then return 0; fi; sleep 0.1; done; return 1; }
fall_laeuft() { [ "$(find "$t" -mindepth 3 -maxdepth 3 -name 't.*' 2> /dev/null | wc -l)" -ge 1 ]; }
reste_da() { [ -n "$(ls -A "$t" 2> /dev/null)" ] || pgrep -f -- "$t" > /dev/null 2>&1; }
# warte_auf_reste — bis zu 8 s auf ein Verzeichnis oder einen Prozess, den ein nicht beendeter Nachkomme hinterlässt.
warte_auf_reste() { local i; for i in $(seq 1 80); do if reste_da; then return 0; fi; sleep 0.1; done; return 0; }
# abbruchfall <Signal> <Exit> — fährt $skript18 mit eigenem TMPDIR in eigener Sitzung, sendet das Signal, sobald dort ein Fall läuft, und
# prüft danach Exit, leeren TMPDIR und Prozesse mit dem Temp-Pfad in der Befehlszeile.
abbruchfall() {
  local sig="$1" soll="$2" pid=0 rc=0 bg ok=1 reste
  mkdir -p "$run"
  t="$(mktemp -d "$run/tmp.XXXXXX")" || exit 1
  ( trap - INT TERM; exec env -u GP_MUT -u GP_ENTDECKE -u GP_FALL TMPDIR="$t" GP_PARALLEL=2 GP_INNEN=1 setsid bash -c 'echo $$ > "$1"; exec bash "$2"' _ "$t.pid" "$skript18" > "$t.out" 2> "$t.err" ) &
  bg=$!
  warte_auf fall_laeuft || ok=0
  pid="$(cat "$t.pid" 2> /dev/null || echo 0)"
  if [ "$pid" -gt 1 ]; then kill -s "$sig" "$pid" 2> /dev/null; else ok=0; fi
  wait "$bg" || rc=$?
  warte_auf_reste
  [ "$rc" -eq "$soll" ] || ok=0
  reste="$(ls -A "$t" 2> /dev/null)"
  ! reste_da || ok=0
  if [ "$ok" -eq 0 ]; then
    printf 'erwartet: Exit %s, TMPDIR leer, kein Prozess mit dem Pfad\ngelesen: Exit %s, TMPDIR: %s\nstderr:\n%s\n' "$soll" "$rc" "$reste" "$(cat "$t.err" 2> /dev/null)" >> "$run/log.$fallname"
    bad=1
  fi
  if [ "$pid" -gt 1 ]; then kill -s KILL -- "-$pid" 2> /dev/null || true; fi
}
# Normales Ende: eine Kopie mit einem einzigen schnellen Lauf endet mit Exit 0 und lässt nichts in ihrem TMPDIR.
f_p18_ende() {
  local rc=0 reste
  mkdir -p "$run"
  t="$(mktemp -d "$run/tmp.XXXXXX")" || exit 1
  env -u GP_FALL -u GP_ENTDECKE TMPDIR="$t" GP_INNEN=1 GP_MUT='^fall-ohne-funktion$' bash "$skript18" > "$t.out" 2> "$t.err" || rc=$?
  reste="$(ls -A "$t")"
  if [ "$rc" -ne 0 ] || [ -n "$reste" ]; then
    bad=1; printf 'erwartet: Exit 0, TMPDIR leer\ngelesen: Exit %s, TMPDIR: %s\n' "$rc" "$reste" >> "$run/log.$fallname"
  fi
}
f_p18_sigint() { abbruchfall INT 130; }
f_p18_sigterm() { abbruchfall TERM 143; }

# --- Läufe ------------------------------------------------------------------------------------------------------------------------------------------
FAELLE_P1=(p1-gleich p1-text-gleich p1-reihenfolge-schritte p1-reihenfolge-gates p1-leerzeile p1-kommentarzeile p1-kommentarzeile-eingerueckt
  p1-pin-kurz p1-pin-gross p1-pin-ohne-tag p1-pin-tag-leer p1-pin-andere-action)
for ((n = 1; n <= N; n++)); do FAELLE_P1+=("p1-fehlt-$n" "p1-geaendert-$n"); done
for k in "${ZUSATZ_NAMEN[@]}"; do FAELLE_P1+=("p1-zusatz-$k"); done
FAELLE_P6=(p6-default-vorgaenger-fehlt-fehlt p6-default-vorgaenger-fehlt-ungueltig p6-event-fehlt p6-event-leer p6-event-unbekannt p6-event-schedule p6-default-fehlt p6-default-ungueltig p6-ungenutzt p6-ungenutzt-prbase p6-ungenutzt-default)
for f in "${FORMEN[@]}"; do FAELLE_P6+=("p6-head-$f" "p6-prbase-$f" "p6-before-$f"); done
FAELLE_RANGE=("${FAELLE_P6[@]}" p7-pr p7-pr-merge-kopf p7-pr-basis-fehlt-im-klon p7-pr-kopf-fehlt-im-klon p7-pr-tag-objekt p7-push p7-push-neu p7-push-neu-ohne-vorfahr p7-push-neu-ohne-standard
  p7-push-vorgaenger-fehlt p7-leer p7-flach p8-stdout-eine-zeile p8-fehler-stdout-leer p8-git-fehler p8-rev-list-fehler p9-ereignis-vor-head p9-head-vor-prbase
  p9-head-vor-before p9-eingabe-vor-flach p9-flach-vor-aufloesung p9-aufloesung-vor-leer)
FAELLE_P18=(p18-sigint p18-sigterm p18-ende)

# faelle <gruppe> — schreibt die Namen der roten Fälle der Gruppe (p1, range, p18, ok = p1 und range und p18), einen je Zeile. FALL_ZUSATZ nennt einen
# weiteren Namen (Probe: ein Name ohne Funktion ist ein roter Fall).
faelle() {
  local n
  case "$1" in
    p1) set -- "${FAELLE_P1[@]}" ;;
    range) set -- "${FAELLE_RANGE[@]}" ;;
    p18) if [ -n "${GP_INNEN-}" ]; then set --; else set -- "${FAELLE_P18[@]}"; fi ;;
    ok) set -- "${FAELLE_P1[@]}" "${FAELLE_RANGE[@]}"; [ -n "${GP_INNEN-}" ] || set -- "$@" "${FAELLE_P18[@]}" ;;
  esac
  for n in "$@" ${FALL_ZUSATZ-}; do
    fallname="$n"
    mkdir -p "$run"
    declare -F "f_${n//-/_}" > /dev/null || { echo "$n"; continue; }
    ( bad=0; "f_${n//-/_}"; exit "$bad" ) > /dev/null 2>> "$run/stderr" || echo "$n"
  done
}
protokoll_zeigen() { [ -z "${GP_FALL-}" ] || [ ! -f "$run/log.$GP_FALL" ] || cat "$run/log.$GP_FALL" >&2; }

maxjobs="${GP_PARALLEL:-4}"
wartejob() { while [ "$(jobs -rp | wc -l)" -ge "$maxjobs" ]; do wait -n || true; done; }
# lauf_welt — legt die Welt des Laufs an (nur für die Gruppen range und ok).
lauf_welt() { mkdir -p "$run/welt"; welt "$run/welt"; }
# bewerte <name> <rote Fälle> <Gruppe> — fährt die Gruppe und vergleicht die roten Fälle mit den erwarteten.
bewerte() {
  local ist erwartet
  ist="$(faelle "$3" | LC_ALL=C sort -u | tr '\n' ' ')"
  erwartet="$(printf '%s\n' $2 | LC_ALL=C sort -u | tr '\n' ' ')"
  if [ -n "${GP_ENTDECKE-}" ]; then echo "$1 => ${ist:-nichts}" >&2; fi
  echo "$1" >> "$arbeit/gefahren"
  [ "$ist" = "$erwartet" ] || melde "Mutation '$1': erwartet rot: ${2:-nichts}; rot sind: ${ist:-nichts}"
}
# fn_ersetze <alt> <neu> — Textersatz in pruefe_ci (Quelltext aus declare -f), wirkt im aktuellen Lauf.
fn_ersetze() {
  local src alt="$1" rest
  src="$(declare -f pruefe_ci)"
  if [[ $src != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
  rest="${src#*"$alt"}"
  if [[ $rest == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
  eval "${src/"$alt"/"$2"}"
}
# mut1 <name> <rote Fälle> <Code> — Mutation von pruefe_ci: der Code läuft im Lauf vor den Fällen (fn_ersetze, Zuweisung an PIN_RE, SOLL).
mut1() {
  name="$1"
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  ( run="$arbeit/$name"; mkdir -p "$run"; eval "$3"; bewerte "$name" "$2" p1 ) &
}
# mutr <name> <rote Fälle> <alt> <neu> [<alt> <neu>…] — Mutation einer Kopie von ci-range.sh; rot sind genau die benannten Fälle der Gruppe range.
mutr() {
  local name="$1" rote="$2" inhalt neu alt rest
  shift 2
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  (
    run="$arbeit/$name"; mkdir -p "$run"; rangesh="$run/ci-range.sh"
    inhalt="$(cat "$repo/tools/harness/ci-range.sh"; printf x)"; inhalt="${inhalt%x}"; neu="$inhalt"
    while [ "$#" -ge 2 ]; do
      alt="$1"
      if [[ $neu != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
      rest="${neu#*"$alt"}"
      if [[ $rest == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
      neu="${neu/"$alt"/"$2"}"
      shift 2
    done
    if [[ $neu == "$inhalt" ]]; then melde "Mutation '$name': Ersatz ändert nichts"; exit 0; fi
    printf '%s' "$neu" > "$rangesh"
    lauf_welt
    bewerte "$name" "$rote" range
  ) &
}

if [[ ok =~ ${GP_MUT:-.} ]]; then
  (
    run="$arbeit/ok"; mkdir -p "$run"; rangesh="$repo/tools/harness/ci-range.sh"; skript18="$selbst"
    lauf_welt
    ist="$(faelle ok | LC_ALL=C sort -u | tr '\n' ' ')"
    echo ok >> "$arbeit/gefahren"
    [ -z "$ist" ] || melde "unveränderte Kopie ist nicht grün: $ist"
    protokoll_zeigen
    if [ -n "${GP_FALL-}" ]; then sed 's/^/stderr: /' "$run/stderr" >&2; fi
  ) &
fi

# mut18 <name> <rote Fälle> <alt> <neu> — Textersatz im Skript (Kopie in einem Baum mit den Dateien, die es braucht); rot sind genau die benannten p18-Fälle.
mut18() {
  local name="$1" rote="$2" alt="$3" inhalt neu m f
  [ -z "${GP_INNEN-}" ] || return 0
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  (
    run="$arbeit/$name"; m="$run/baum"; mkdir -p "$m/tools/harness" "$m/tools/dev/harness" "$m/.github/workflows"
    for f in .github/workflows/ci.yml tools/harness/ci-range.sh; do cp "$repo/$f" "$m/$f"; done
    inhalt="$(cat "$selbst"; printf x)"; inhalt="${inhalt%x}"
    if [[ $inhalt != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
    if [[ ${inhalt#*"$alt"} == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
    neu="${inhalt/"$alt"/"$4"}"
    printf '%s' "$neu" > "$m/tools/dev/harness/ci-gegenprobe.sh"
    skript18="$m/tools/dev/harness/ci-gegenprobe.sh"
    bewerte "$name" "$rote" p18
  ) &
}

# Schutz gegen Rekursion: Eine Kopie (GP_INNEN=1), deren GP_MUT auf m-p18 passt, fährt nichts und endet mit Exit 1.
if [[ rekursion-schutz =~ ${GP_MUT:-.} ]] && [ -z "${GP_INNEN-}" ]; then
  (
    run="$arbeit/rekursion-schutz"; mkdir -p "$run"
    rc=0; ausgabe="$(env -u GP_FALL GP_INNEN=1 GP_MUT=m-p18 TMPDIR="$run" timeout 60 bash "$selbst" 2>&1)" || rc=$?
    echo rekursion-schutz >> "$arbeit/gefahren"
    [[ $ausgabe == *"keine Mutation und keine Kopie gefahren"* ]] && [ "$rc" -eq 1 ] || melde "Rekursionsschutz: erwartet Exit 1 ohne Lauf; Exit $rc: $ausgabe"
  ) &
fi

# Probe der Fallliste: ein Name ohne Funktion ist rot.
if [[ fall-ohne-funktion =~ ${GP_MUT:-.} ]]; then
  (
    run="$arbeit/fall-ohne-funktion"; mkdir -p "$run"
    ist="$(FALL_ZUSATZ=p1-tippfehler faelle p1 | LC_ALL=C sort -u | tr '\n' ' ')"
    echo fall-ohne-funktion >> "$arbeit/gefahren"
    [ "$ist" = "p1-tippfehler " ] || melde "Fall ohne Funktion: erwartet rot: p1-tippfehler; rot sind: ${ist:-nichts}"
  ) &
fi

# --- Mutationen von pruefe_ci (Punkt 1) -------------------------------------------------------------------------------------------------------------
# Je Zeile n zwei Mutationen: die Zeile ist optional (fehlt sie, rückt die Prüfung nicht weiter) und ihr Inhalt wird nicht verglichen.
for ((n = 1; n <= N; n++)); do
  i=$((n - 1))
  if [ "$i" -eq "$PIN" ]; then igel='[[ ${ist[j]-} =~ $PIN_RE ]]'; else igel='[ "${ist[j]-}" = "${SOLL[i]}" ]'; fi
  mut1 "m-p1-zeile-$n-optional" "p1-fehlt-$n" "fn_ersetze 'if [ \"\$i\" -eq \"\$PIN\" ]; then' 'if [ \"\$i\" -eq $i ] && ! $igel; then continue; fi; if [ \"\$i\" -eq \"\$PIN\" ]; then'"
  rote="p1-geaendert-$n"
  if [ "$i" -eq "$PIN" ]; then rote+=" p1-pin-kurz p1-pin-gross p1-pin-ohne-tag p1-pin-tag-leer p1-pin-andere-action"; fi
  mut1 "m-p1-zeile-$n-nicht-verglichen" "$rote" "fn_ersetze 'if [ \"\$i\" -eq \"\$PIN\" ]; then' 'if [ \"\$i\" -eq $i ]; then j=\$((j + 1)); continue; fi; if [ \"\$i\" -eq \"\$PIN\" ]; then'"
done
for k in "${!ZUSATZ[@]}"; do
  wort="${ZUSATZ[k]%%|*}"
  mut1 "m-p1-zusatz-${ZUSATZ_NAMEN[k]}-erlaubt" "p1-zusatz-${ZUSATZ_NAMEN[k]}" "fn_ersetze 'ist+=(\"\$z\")' '[[ \$z == *\"$wort\"* ]] || ist+=(\"\$z\")'"
done
mut1 m-p1-zusatzzeilen-erlaubt "$(printf 'p1-zusatz-%s ' "${ZUSATZ_NAMEN[@]}")" "fn_ersetze '[ \"\${#ist[@]}\" -eq \"\$j\" ]' 'true'"
mut1 m-p1-pin-grossbuchstaben "p1-pin-gross" 'PIN_RE="${PIN_RE/0-9a-f/0-9a-fA-F}"'
mut1 m-p1-pin-39-ziffern "p1-pin-kurz" 'PIN_RE="${PIN_RE/"{40}"/"{39,40}"}"'
mut1 m-p1-pin-ohne-kommentar "p1-pin-ohne-tag" 'PIN_RE="${PIN_RE/ # \[^ \].*\$/( # [^ ].*)?\$}"'
mut1 m-p1-pin-kommentar-ohne-text "p1-pin-tag-leer" 'PIN_RE="${PIN_RE/"[^ ]"/}"'
mut1 m-p1-pin-andere-action "p1-pin-andere-action" 'PIN_RE="${PIN_RE/actions?checkout@/actions/[a-z-]+@}"'
mut1 m-p1-kommentarzeilen-nicht-ignoriert "p1-kommentarzeile p1-kommentarzeile-eingerueckt" "fn_ersetze '^\\ *(#|\$)' '^\\ *\$'"
mut1 m-p1-leerzeilen-nicht-ignoriert "p1-leerzeile" "fn_ersetze '^\\ *(#|\$)' '^\\ *#'"
mut1 m-p1-kommentar-irgendwo-in-der-zeile "p1-gleich p1-text-gleich p1-leerzeile p1-kommentarzeile p1-kommentarzeile-eingerueckt" "fn_ersetze '^\\ *(#|\$)' '^.*#|^\\ *\$'"
mut1 m-p1-kommentar-nur-am-zeilenanfang "p1-kommentarzeile-eingerueckt" "fn_ersetze '^\\ *(#|\$)' '^(#|\$)'"
# Reihenfolge: die Prüfung nimmt zusätzlich die Folge mit vertauschten Schritten an.
PERM_AC=("${SOLL[@]:0:S_C}" "${SOLL[@]:S_A:S_G-S_A}" "${SOLL[@]:S_C:S_A-S_C}" "${SOLL[@]:S_G}")
PERM_G=("${SOLL[@]:0:S_C}" "${SOLL[@]:S_G}" "${SOLL[@]:S_C:S_G-S_C}")
tauschbar='src="$(declare -f pruefe_ci)"; eval "${src/pruefe_ci/pruefe_ci_orig}"; pruefe_ci() { pruefe_ci_orig "$1" && return 0; local -a s=("${SOLL[@]}"); local r; SOLL=("${%s[@]}"); pruefe_ci_orig "$1"; r=$?; SOLL=("${s[@]}"); return $r; }'
# shellcheck disable=SC2059
mut1 m-p1-reihenfolge-commit-adr-tauschbar "p1-reihenfolge-schritte" "$(printf "$tauschbar" PERM_AC)"
# shellcheck disable=SC2059
mut1 m-p1-reihenfolge-gates-zuerst "p1-reihenfolge-gates" "$(printf "$tauschbar" PERM_G)"

# --- Mutationen von ci-range.sh (Punkte 6 bis 9) ----------------------------------------------------------------------------------------------------
EV='p6-event-fehlt p6-event-leer p6-event-unbekannt p6-event-schedule'
mutr m-p6-event-ungeprueft "$EV p8-fehler-stdout-leer p9-ereignis-vor-head" '  *) ungueltig CI_EVENT ;;' '  *) ;;'
mutr m-p6-event-leer-erlaubt "p6-event-fehlt p6-event-leer p8-fehler-stdout-leer p9-ereignis-vor-head" '  pull_request | push) ;;' "  pull_request | push | '') ;;"
mutr m-p6-event-dispatch-erlaubt "p6-event-unbekannt" '  pull_request | push) ;;' '  pull_request | push | workflow_dispatch) ;;'
mutr m-p6-event-schedule-erlaubt "p6-event-schedule" '  pull_request | push) ;;' '  pull_request | push | schedule) ;;'
mutr m-p6-head-ungeprueft "$(printf 'p6-head-%s ' "${FORMEN[@]}") p9-head-vor-prbase p9-head-vor-before p9-eingabe-vor-flach" '[[ $kopf_roh =~ $HEX ]] || ungueltig CI_HEAD' ':'
mutr m-p6-prbase-ungeprueft "$(printf 'p6-prbase-%s ' "${FORMEN[@]}")" '  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PR_BASE' '  :'
mutr m-p6-before-ungeprueft "$(printf 'p6-before-%s ' "${FORMEN[@]}")" '  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PUSH_BEFORE' '  :'
mutr m-p6-default-ungeprueft "p6-default-fehlt p6-default-ungueltig p6-default-vorgaenger-fehlt-fehlt p6-default-vorgaenger-fehlt-ungueltig" '  [ -n "$standard" ] && [ "$(git check-ref-format --branch "$standard" 2> /dev/null)" = "$standard" ] || ungueltig CI_DEFAULT_BRANCH' '  :'
mutr m-p6-default-leer-erlaubt "p6-default-fehlt p6-default-vorgaenger-fehlt-fehlt" '[ -n "$standard" ] && [ "$(git' '[ "$(git'
mutr m-p6-default-format-ungeprueft "p6-default-ungueltig p6-default-vorgaenger-fehlt-ungueltig" '[ "$(git check-ref-format --branch "$standard" 2> /dev/null)" = "$standard" ]' 'true'
mutr m-p6-pr-liest-before "p6-ungenutzt" '  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PR_BASE' $'  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PR_BASE\n  [[ ${CI_PUSH_BEFORE-} =~ $HEX ]] || ungueltig CI_PUSH_BEFORE'
mutr m-p6-push-liest-prbase "p6-ungenutzt-prbase" '  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PUSH_BEFORE' $'  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PUSH_BEFORE\n  [[ ${CI_PR_BASE-} =~ $HEX ]] || ungueltig CI_PR_BASE'
mutr m-p6-push-liest-standard "p6-ungenutzt-default" '  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PUSH_BEFORE' $'  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PUSH_BEFORE\n  pruefe_standard'
mutr m-p6-form-39-ziffern "p6-head-kurz p6-prbase-kurz p6-before-kurz" "HEX='^[0-9a-f]{40}\$'" "HEX='^[0-9a-f]{39,40}\$'"
mutr m-p6-form-41-ziffern "p6-head-lang p6-prbase-lang p6-before-lang" "HEX='^[0-9a-f]{40}\$'" "HEX='^[0-9a-f]{40,41}\$'"
mutr m-p6-form-grossbuchstaben "p6-head-gross p6-prbase-gross p6-before-gross" "HEX='^[0-9a-f]{40}\$'" "HEX='^[0-9a-fA-F]{40}\$'"
mutr m-p6-form-nicht-hex "p6-head-ungueltig p6-prbase-ungueltig p6-before-ungueltig" "HEX='^[0-9a-f]{40}\$'" "HEX='^[0-9a-z]{40}\$'"
mutr m-p7-pr-nimmt-vorgaenger "$(printf 'p6-prbase-%s ' "${FORMEN[@]}") p6-ungenutzt p7-pr p7-pr-merge-kopf p7-pr-basis-fehlt-im-klon p7-pr-tag-objekt p8-fehler-stdout-leer" '  basis_roh="${CI_PR_BASE-}"' '  basis_roh="${CI_PUSH_BEFORE-}"'
mutr m-p7-push-nimmt-pr-basis "$(printf 'p6-before-%s ' "${FORMEN[@]}") p6-default-fehlt p6-default-ungueltig p6-ungenutzt-prbase p6-ungenutzt-default p7-push p7-push-neu-ohne-standard p7-push-neu-ohne-vorfahr p7-push-vorgaenger-fehlt p8-fehler-stdout-leer p8-stdout-eine-zeile p6-default-vorgaenger-fehlt-fehlt p6-default-vorgaenger-fehlt-ungueltig" '  basis_roh="${CI_PUSH_BEFORE-}"' '  basis_roh="${CI_PR_BASE-}"'
mutr m-p7-nullen-kein-neuer-zweig "p7-push-neu p7-leer" 'if [[ $basis_roh =~ $NULLEN ]]; then' 'if false; then'
mutr m-p7-merge-base-fehler-ohne-exit "p7-push-neu-ohne-standard p7-push-neu-ohne-vorfahr" '|| fehler "keine Merge-Base von $kopf_roh und refs/remotes/origin/$standard (Zweig fehlt im Klon oder kein gemeinsamer Vorfahr)"' '|| true'
mutr m-p7-spitze-statt-merge-base "p7-push-neu p7-push-vorgaenger-fehlt p7-push-neu-ohne-vorfahr p7-leer p8-stdout-eine-zeile" 'git merge-base "$kopf_roh" "refs/remotes/origin/$standard"' 'git rev-parse --verify --quiet "refs/remotes/origin/$standard^{commit}"'
mutr m-p7-vorgaenger-fehlt-exit-2 "p7-push-vorgaenger-fehlt p8-stdout-eine-zeile p6-default-vorgaenger-fehlt-fehlt p6-default-vorgaenger-fehlt-ungueltig" '  [ "$event" = push ] || fehler "Basis $basis_roh nicht im Klon"' '  fehler "Basis $basis_roh nicht im Klon"'
mutr m-p7-basis-fehlt-ohne-exit "p7-pr-basis-fehlt-im-klon p8-fehler-stdout-leer" '  [ "$event" = push ] || fehler "Basis $basis_roh nicht im Klon"' '  [ "$event" = push ] || true'
mutr m-p7-kopf-fehlt-ohne-exit "p7-pr-kopf-fehlt-im-klon p8-fehler-stdout-leer" '|| fehler "Kopf $kopf_roh nicht im Klon"' '|| true'
mutr m-p7-leer-exit-0 "p7-leer p8-fehler-stdout-leer" $'ist leer" >&2\n  exit 1' $'ist leer" >&2\n  exit 0'
mutr m-p7-flach-nicht-erkannt "p7-flach p9-flach-vor-aufloesung" '  true) fehler "flacher Klon, der Checkout braucht fetch-depth: 0" ;;' '  true) ;;'
mutr m-p7-merge-kopf-abgelehnt "p7-pr-merge-kopf p8-rev-list-fehler" $'kopf="$(aufloesen "$kopf_roh")" || fehler "Kopf $kopf_roh nicht im Klon"' $'kopf="$(aufloesen "$kopf_roh")" || fehler "Kopf $kopf_roh nicht im Klon"\n[ "$(git rev-list --parents -n 1 "$kopf" | wc -w)" -le 2 ] || fehler "Kopf ist ein Merge-Commit"'
mutr m-p7-namen-statt-shas "p7-pr-tag-objekt" 'git rev-parse --verify --quiet "$1^{commit}" 2> /dev/null' 'git rev-parse --verify --quiet "$1^{commit}" > /dev/null 2>&1 && echo "$1"'
mutr m-p8-hinweis-auf-stdout "p7-push-vorgaenger-fehlt p8-stdout-eine-zeile p6-default-vorgaenger-fehlt-fehlt p6-default-vorgaenger-fehlt-ungueltig" 'Basis ist der Standard-Zweig" >&2' 'Basis ist der Standard-Zweig"'
mutr m-p8-rev-list-fehler-ignoriert "p8-rev-list-fehler" '|| fehler "git rev-list --count $basis..$kopf schlug fehl"' '|| true'
mutr m-p6-default-vorgaenger-ungeprueft "p6-default-vorgaenger-fehlt-fehlt p6-default-vorgaenger-fehlt-ungueltig p7-push-vorgaenger-fehlt p8-stdout-eine-zeile" $'Standard-Zweig" >&2\n  pruefe_standard\n' $'Standard-Zweig" >&2\n'
mutr m-p8-teil-range-auf-stdout "p7-pr-kopf-fehlt-im-klon p8-fehler-stdout-leer" '|| fehler "Kopf $kopf_roh nicht im Klon"' '|| { echo "$basis.."; fehler "Kopf $kopf_roh nicht im Klon"; }'
mutr m-p9-ereignis-nach-head "p9-ereignis-vor-head" $'# (1) Ereignis\n' $'# (1) Ereignis\n[[ ${CI_HEAD-} =~ $HEX ]] || ungueltig CI_HEAD\n'
mutr m-p9-head-nach-prbase "p9-head-vor-prbase" $'kopf_roh="${CI_HEAD-}"\n' $'[ "$event" != pull_request ] || [[ ${CI_PR_BASE-} =~ $HEX ]] || ungueltig CI_PR_BASE\nkopf_roh="${CI_HEAD-}"\n'
mutr m-p9-head-nach-before "p9-head-vor-before" $'kopf_roh="${CI_HEAD-}"\n' $'[ "$event" != push ] || [[ ${CI_PUSH_BEFORE-} =~ $HEX ]] || ungueltig CI_PUSH_BEFORE\nkopf_roh="${CI_HEAD-}"\n'
mutr m-p9-flach-vor-eingabe "p9-eingabe-vor-flach p8-fehler-stdout-leer p8-git-fehler" $'# (2) Variablen, die das Ereignis liest\n' $'[ "$(git rev-parse --is-shallow-repository)" != true ] || fehler "flacher Klon, der Checkout braucht fetch-depth: 0"\n'
mutr m-p9-aufloesung-vor-flach "p9-flach-vor-aufloesung" $'# (3) flacher Klon\n' $'[ "$event" != pull_request ] || aufloesen "$basis_roh" > /dev/null || fehler "Basis $basis_roh nicht im Klon"\n'
mutr m-p9-leer-vor-aufloesung "p9-aufloesung-vor-leer" $'# (4) Auflösung von Basis und Kopf\n' $'[ "$basis_roh" != "$kopf_roh" ] || { echo "ci-range: Range ist leer" >&2; exit 1; }\n'
# Punkt 18
mut18 m-p18-ende-falle-fehlt "p18-ende" $'\ntrap raeume_auf EXIT\n' $'\n'
mut18 m-p18-int-falle-fehlt "p18-sigint" $'\ntrap abbruch_int INT\n' $'\n'
mut18 m-p18-term-falle-fehlt "p18-sigterm" $'\ntrap abbruch_term TERM\n' $'\n'
mut18 m-p18-int-kinder-bleiben "p18-sigint" $'\nabbruch_int() { trap - EXIT INT TERM; beende_kinder; ' $'\nabbruch_int() { trap - EXIT INT TERM; '
mut18 m-p18-term-kinder-bleiben "p18-sigterm" $'\nabbruch_term() { trap - EXIT INT TERM; beende_kinder; ' $'\nabbruch_term() { trap - EXIT INT TERM; '
mut18 m-p18-int-verzeichnis-zu-frueh "p18-sigint" $'\nabbruch_int() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 130; }' $'\nabbruch_int() { trap - EXIT INT TERM; raeume_auf; sleep 5; beende_kinder; exit 130; }'
mut18 m-p18-term-verzeichnis-zu-frueh "p18-sigterm" $'\nabbruch_term() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 143; }' $'\nabbruch_term() { trap - EXIT INT TERM; raeume_auf; sleep 5; beende_kinder; exit 143; }'

wait
fehler=0
if [ -s "$arbeit/fehler" ]; then cat "$arbeit/fehler" >&2; fehler=1; fi
if [ "$(wc -l < "$arbeit/gefahren")" -eq 0 ]; then echo "ci-gegenprobe: keine Mutation und keine Kopie gefahren (GP_MUT passt auf nichts)" >&2; fehler=1; fi
if [ "$fehler" -eq 0 ]; then echo "ci-gegenprobe: grün ($(wc -l < "$arbeit/gefahren") Läufe)"; fi
exit "$fehler"
