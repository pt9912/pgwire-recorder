#!/usr/bin/env bash
# gates-schnell-gegenprobe — prüft die schnelle Gate-Kette (harness/mk/gates-schnell.mk,
# tools/harness/gates-schnell.sh) gegen die Punkte 1 bis 7 von SPEC-052
# (harness/vertraege/SPEC-052-gate-kette.md) und die Namens- und Zuordnungsregel der Gegenproben
# (Punkte 1 bis 3 von SPEC-060, tools/harness/gegenproben-regel.sh) an Kopien in einem Temp-Verzeichnis;
# der Arbeitsbaum bleibt unberührt. Eine Kopie trägt Makefile, Fragment der Kette, Skripte der Kette
# und ein Fragment mit Gates, die ihren Namen in eine Logdatei schreiben; die
# unveränderte Kopie muss jede Prüfung bestehen, und je Mutant (sed auf der Kopie)
# muss die benannte Prüfung rot werden. Die Deckung der Ausnahmeliste (Punkt 4) läuft
# zusätzlich gegen den Arbeitsbaum selbst: ein neues Gate dieser Namen macht die
# Gegenprobe rot, bis jemand es in die Liste aufnimmt oder bewusst schnell lässt. Die Namensregel
# läuft ebenfalls gegen den Arbeitsbaum: ein Ziel auf -gegenprobe oder a-check-negativ in
# GATE_CHECKS, in beiden oder in keiner Liste macht sie rot (Fall p2-arbeitsbaum: eine Kopie der
# Gegenprobe als eigener Arbeitsbaum mit einer Verletzung wird rot, und ohne die Zeile der Prüfung bliebe
# sie grün). Ihre dreizehn Fälle der Regel (p2-*) legt je eine Kopie an, dazu p1-liste-aus-variable
# und p2-locale; je Fall mutiert eine Kopie der Regel die Prüfung, die ihn fängt (rot sind dann genau
# die benannten Fälle). GP_NUR_ARBEITSBAUM=1 fährt nur die Prüfung des Arbeitsbaums.
# Das Abbrechen eines Gates (Punkt 6) fährt ein Gate, das seine Prozessgruppe
# beendet; der Lauf der Kopie steht dafür in einer eigenen Sitzung (setsid).
#
# GRENZE der Gegenprobe: Die Gates der Kopie sind Attrappen; ob die echten Gates in
# der echten Reihenfolge laufen, prüft `make gates-schnell` selbst, nicht sie. Die
# Aufruf-Voraussetzung (Punkt 8) ist eine Bindung der Zusage und hat keinen Fall. Sie
# fährt `make gates` nie, nur `make -n gates`; was der volle Lauf mit dem Stand der
# schnellen Kette tut, ist nicht zugesagt und hat keinen Fall.
#
# Ausgang: 0, wenn jede Prüfung wie erwartet endet, sonst 1.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
wurzel="$(cd "$hier/../.." && pwd)"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT
fehler=0

melde() { echo "gates-schnell-gegenprobe: ROT — $*" >&2; fehler=1; }

SCHNELL="docs-check kopf-check baseline-verify abdeckung-check lint test build a-check"
LANGSAM="test-container test-integration test-postgres-versionen"
# Gegenproben der Kopie: Liste für den Container und für den Host (SPEC-060).
GP_CONTAINER="lint-gegenprobe coverage-gegenprobe kopf-check-gegenprobe abdeckung-gegenprobe hook-gegenprobe integration-volumes-gegenprobe meldungskatalog-check-gegenprobe stop-hook-gegenprobe a-check-negativ test-container-gegenprobe test-postgres-versionen-gegenprobe"
GP_HOST="container-start-gegenprobe"

# shellcheck source=tools/harness/gegenproben-regel.sh
. "$hier/gegenproben-regel.sh"

# gate_variable <wurzel> <name> — Wörter der Make-Variable aus der Datenbank von make.
gate_variable() {
  { make -s --no-print-directory -C "$1" -pq help 2>/dev/null || true; } \
    | awk -v n="$2" '$1 == n && ($2 == ":=" || $2 == "=") { for (i = 3; i <= NF; i++) print $i }'
}

# in_liste <wort> <wort>… — Exit 0, wenn das erste Wort unter den übrigen steht.
in_liste() {
  local g="$1" x
  shift
  for x in "$@"; do [ "$x" = "$g" ] && return 0; done
  return 1
}

# deckung <wurzel> — Punkt 4 von SPEC-052 und die Namensregel (SPEC-060 Punkt 2): schreibt die
# Namen der verletzten Prüfungen.
deckung() {
  local ausnahme gates g
  mapfile -t ausnahme < <(gate_variable "$1" GATE_SCHNELL_AUSNAHME)
  mapfile -t gates < <(gate_variable "$1" GATE_CHECKS)
  for g in "${ausnahme[@]}"; do
    in_liste "$g" "${gates[@]}" || echo "p4-eintrag"
  done
  for g in test-container test-integration test-postgres-versionen; do
    in_liste "$g" "${ausnahme[@]}" || echo "p4-name"
  done
  namensregel "$1"
}

# fixture <name> — legt eine Kopie an und setzt w auf ihre Wurzel.
fixture() {
  w="$arbeit/$1"
  mkdir -p "$w/harness/mk" "$w/tools/harness" "$w/.harness/state"
  cp "$wurzel/Makefile" "$w/Makefile"
  cp "$wurzel/harness/mk/gates-schnell.mk" "$wurzel/harness/mk/enforce.mk" "$w/harness/mk/"
  cp "$wurzel/tools/harness/gates-schnell.sh" "$wurzel/tools/harness/working-tree-hash.sh" \
    "$wurzel/tools/harness/dateiliste.sh" "$wurzel/tools/harness/record-gates.sh" "$w/tools/harness/"
  printf '.harness/state/\n' > "$w/.gitignore"
  {
    printf 'FAKE_SCHNELL := %s neu-x\n' "$SCHNELL"
    printf 'FAKE_LANGSAM := %s\n' "$LANGSAM"
    printf 'FAKE := $(sort $(FAKE_SCHNELL) $(FAKE_LANGSAM) %s %s)\n.PHONY: $(FAKE)\n$(FAKE):\n' "$GP_CONTAINER" "$GP_HOST"
    printf '\t@echo $@ >> "$(LOG)"; if [ "$(FAIL)" = "$@" ]; then exit 1; fi; if [ "$(ABBRUCH)" = "$@" ]; then kill -KILL 0; fi; if [ "$(AENDERE)" = "$@" ]; then echo x >> aenderung; fi\n'
    printf 'GATE_CHECKS += $(FAKE_SCHNELL) $(FAKE_LANGSAM)\n'
    printf 'GEGENPROBEN_CONTAINER += %s\n' "$GP_CONTAINER"
    printf 'GEGENPROBEN_HOST += %s\n' "$GP_HOST"
  } > "$w/harness/mk/fake.mk"
  git -C "$w" init -q
}

# lauf <wurzel> <make-argument>… — setzt code; die Logdatei der Gates ist $arbeit/log.
lauf() {
  local w="$1"
  shift
  : > "$arbeit/log"
  set +e
  {
    env -u MAKELEVEL -u MFLAGS MAKEFLAGS="${mf:-}" LOG="$arbeit/log" FAIL="${fail:-}" ABBRUCH="${abbruch:-}" \
      AENDERE="${aendere:-}" \
      setsid -w make -s --no-print-directory -C "$w" "$@" > "$arbeit/out" 2>&1
    code=$?
  } 2> /dev/null
  set -e
}

# pruefe <wurzel> — schreibt die Namen der verletzten Prüfungen, einen je Zeile.
pruefe() {
  local w="$1" stand="$1/.harness/state/gates-schnell-passed.diffsha" voll="$1/.harness/state/gates-passed.diffsha"
  local g gates
  deckung "$w"

  # Punkte 2, 3, 5, 7: grüner Lauf neben einem vorhandenen vollen Nachweis.
  rm -f "$stand"; echo VOLL > "$voll"; fail=""; abbruch=""; mf=""; aendere=""
  lauf "$w" gates-schnell
  [ "$code" -eq 0 ] || echo "lauf-gruen"
  for g in $LANGSAM; do grep -qxF "$g" "$arbeit/log" && echo "p2-langsam"; done
  for g in $SCHNELL; do grep -qxF "$g" "$arbeit/log" || echo "p2-schnell"; done
  grep -qxF neu-x "$arbeit/log" || echo "p3-neu"
  for g in $GP_CONTAINER $GP_HOST; do grep -qxF "$g" "$arbeit/log" && echo "p2-gegenprobe-in-schnell"; done
  if [ ! -f "$stand" ] || [ "$(cat "$stand")" != "$(cd "$w" && bash tools/harness/working-tree-hash.sh)" ]; then
    echo "p5-stand"
  fi
  [ "$(cat "$voll" 2> /dev/null)" = VOLL ] || echo "p7-nachweis"

  # Punkt 5: ändert ein Gate des Laufs den Arbeitsbaum, gilt der Hash nach dem Lauf.
  rm -f "$stand"; aendere=docs-check
  lauf "$w" gates-schnell
  aendere=""
  if [ ! -f "$stand" ] || [ "$(cat "$stand")" != "$(cd "$w" && bash tools/harness/working-tree-hash.sh)" ]; then
    echo "p5-ende"
  fi

  # Punkt 7: ohne vollen Nachweis entsteht keiner.
  rm -f "$voll" "$stand"
  lauf "$w" gates-schnell
  [ ! -e "$voll" ] || echo "p7-nachweis-neu"

  # Punkt 6: rotes Gate bei vorhandenem Stand.
  echo ALT > "$stand"; fail=kopf-check
  lauf "$w" gates-schnell
  { [ "$code" -ne 0 ] && [ ! -e "$stand" ]; } || echo "p6-rot"
  [ ! -e "$voll" ] || echo "p7-nachweis-neu"

  # Punkt 6: Lauf mitten im Gate abgebrochen, bei vorhandenem Stand.
  echo ALT > "$stand"; fail=""; abbruch=lint
  lauf "$w" gates-schnell
  { [ "$code" -ne 0 ] && [ ! -e "$stand" ]; } || echo "p6-abbruch"
  abbruch=""

  # Punkt 6: Make-Flags, mit und ohne vorhandenen Stand.
  flagfall "$w" n "" -n gates-schnell
  flagfall "$w" nenv n gates-schnell
  flagfall "$w" q "" -q gates-schnell
  flagfall "$w" t "" -t gates-schnell
  # n zusammen mit i: n gilt, der Stand bleibt unberührt.
  flagfall "$w" ni ni gates-schnell
  # -i: Fehler werden überlesen; der Lauf endet mit 0, ein älterer Stand entfällt, keiner entsteht.
  rm -f "$stand"; mf=""; fail=""
  lauf "$w" -i gates-schnell
  [ ! -e "$stand" ] || echo "p6-i-schreibt"
  echo ALT > "$stand"; fail=kopf-check
  lauf "$w" -i gates-schnell
  fail=""
  [ ! -e "$stand" ] || echo "p6-i-entfernt"
  echo ALT > "$stand"
  lauf "$w" -i gates-schnell
  [ ! -e "$stand" ] || echo "p6-i-entfernt"
  # andere Flags ändern nichts.
  for f in -k -j2; do
    rm -f "$stand"
    lauf "$w" $f gates-schnell
    [ -f "$stand" ] || echo "p6-andere"
  done

  # Punkt 1: make gates fährt jedes Gate von GATE_CHECKS, nichts daneben.
  mapfile -t gates < <(gate_variable "$w" GATE_CHECKS)
  lauf "$w" -n gates
  for g in "${gates[@]}"; do
    grep -qF "echo $g >> " "$arbeit/out" || echo "p1-voll"
  done
  for g in $GP_CONTAINER $GP_HOST; do
    grep -qF "echo $g >> " "$arbeit/out" && echo "p1-voll"
  done
  grep -q 'tools/harness/gates-schnell.sh' "$arbeit/out" && echo "p1-voll"
  grep -q 'tools/harness/record-gates.sh' "$arbeit/out" || echo "p1-voll"

  # Anschluss: beide Ziele enden bei rotem Skript ≠ 0; dass die Gegenprobe in der Liste für den Container
  # steht, hält die Namensregel (deckung).
  rm -rf "$w-rot"
  cp -r "$w" "$w-rot"
  printf 'exit 1\n' > "$w-rot/tools/harness/gates-schnell.sh"
  printf 'exit 1\n' > "$w-rot/tools/harness/gates-schnell-gegenprobe.sh"
  lauf "$w-rot" gates-schnell
  [ "$code" -ne 0 ] || echo "anschluss-exit-schnell"
  lauf "$w-rot" gates-schnell-gegenprobe
  [ "$code" -ne 0 ] || echo "anschluss-exit-gegenprobe"
}

# flagfall <wurzel> <name> <MAKEFLAGS der Umgebung> <make-argument>… — ein Aufruf, der kein Gate
# wirklich fährt (n, q, t), schreibt keinen Stand und entfernt keinen.
flagfall() {
  local w="$1" name="$2" stand="$1/.harness/state/gates-schnell-passed.diffsha"
  mf="$3"; fail=""; abbruch=""
  shift 3
  rm -f "$stand"
  lauf "$w" "$@"
  [ ! -e "$stand" ] || echo "p6-$name-schreibt"
  echo ALT > "$stand"
  lauf "$w" "$@"
  if [ ! -e "$stand" ]; then
    echo "p6-$name-entfernt"
  elif [ "$(cat "$stand")" != ALT ]; then
    echo "p6-$name-schreibt"
  fi
  mf=""
}

# gruen <name> — die unveränderte Kopie besteht jede Prüfung.
gruen() {
  local ist
  ist="$(pruefe "$w" | sort -u | tr '\n' ' ')"
  [ -z "$ist" ] || melde "Kopie '$1' ist nicht grün: $ist"
}

# rot <name> <prüfung> [exakt] — der Mutant färbt die Prüfung rot, mit „exakt“ nur sie.
rot() {
  local ist
  ist="$(pruefe "$w" | sort -u | tr '\n' ' ')"
  case " $ist" in
    *" $2 "*) ;;
    *) melde "Mutant '$1' färbt '$2' nicht rot, rot sind: ${ist:-nichts}" ;;
  esac
  if [ "${3:-}" = exakt ] && [ "$ist" != "$2 " ]; then
    melde "Mutant '$1' färbt nur '$2' rot, rot sind: $ist"
  fi
}

AUSNAHME_ZEILE='/^GATE_SCHNELL_AUSNAHME/'

# --- Deckung und Namensregel im Arbeitsbaum, unveränderte Kopie ----------------------
ist="$(deckung "$wurzel" | sort -u | tr '\n' ' ')"
[ -z "$ist" ] || melde "Arbeitsbaum: Ausnahmeliste, GATE_CHECKS und Listen der Gegenproben stimmen nicht: $ist"
if [ -n "${GP_NUR_ARBEITSBAUM-}" ]; then exit "$fehler"; fi
fixture ok; gruen ok

# Fall p2-arbeitsbaum: Eine Kopie dieser Gegenprobe ist ein eigener Arbeitsbaum, in dem lint-gegenprobe in
# GATE_CHECKS steht; die Prüfung des Arbeitsbaums färbt sie rot. Mutation: ohne die Zeile der Prüfung bleibt sie grün.
selbst="$arbeit/selbst"
mkdir -p "$selbst/tools/harness"
cp "$wurzel/Makefile" "$selbst/"
cp -r "$wurzel/harness" "$selbst/harness"
cp "$hier/gates-schnell-gegenprobe.sh" "$hier/gegenproben-regel.sh" "$selbst/tools/harness/"
sed -i 's/^GATE_CHECKS += lint$/GATE_CHECKS += lint lint-gegenprobe/' "$selbst/harness/mk/lint.mk"
if (GP_NUR_ARBEITSBAUM=1 bash "$selbst/tools/harness/gates-schnell-gegenprobe.sh" > /dev/null 2>&1); then
  melde "Fall 'p2-arbeitsbaum': der Arbeitsbaum mit lint-gegenprobe in GATE_CHECKS bleibt grün"
fi
cp "$selbst/tools/harness/gates-schnell-gegenprobe.sh" "$arbeit/selbst-alt.sh"
sed -i 's/^ist="\$(deckung "\$wurzel" .*$/ist=""/' "$selbst/tools/harness/gates-schnell-gegenprobe.sh"
if cmp -s "$arbeit/selbst-alt.sh" "$selbst/tools/harness/gates-schnell-gegenprobe.sh"; then
  melde "Mutation 'p2-arbeitsbaum-gestrichen' ändert die Gegenprobe nicht"
elif ! (GP_NUR_ARBEITSBAUM=1 bash "$selbst/tools/harness/gates-schnell-gegenprobe.sh" > /dev/null 2>&1); then
  melde "Mutation 'p2-arbeitsbaum-gestrichen': die Gegenprobe ohne Prüfung des Arbeitsbaums ist rot"
fi

# --- Punkt 2: kein langsames Gate ---------------------------------------------------
fixture m-test-integration
sed -i "${AUSNAHME_ZEILE}s/ test-integration//" "$w/harness/mk/gates-schnell.mk"; rot m-test-integration p2-langsam
fixture m-test-container
sed -i "${AUSNAHME_ZEILE}s/ test-container / /" "$w/harness/mk/gates-schnell.mk"; rot m-test-container p2-langsam; rot m-test-container p4-name
fixture m-test-postgres-versionen
sed -i "${AUSNAHME_ZEILE}s/ test-postgres-versionen\$//" "$w/harness/mk/gates-schnell.mk"; rot m-test-postgres-versionen p2-langsam; rot m-test-postgres-versionen p4-name

# --- Punkt 2: jedes andere Gate -----------------------------------------------------
fixture m-schnelles-gate-in-liste
sed -i "${AUSNAHME_ZEILE}s/\$/ kopf-check/" "$w/harness/mk/gates-schnell.mk"; rot m-schnelles-gate-in-liste p2-schnell

# --- Punkt 2: keine Gegenprobe in der schnellen Kette ---------------------------------
fixture m-liste-in-schnell
sed -i 's/\$(filter-out \$(GATE_SCHNELL_AUSNAHME),\$(GATE_CHECKS))/& $(GEGENPROBEN_CONTAINER)/' "$w/harness/mk/gates-schnell.mk"
rot m-liste-in-schnell p2-gegenprobe-in-schnell

# --- Punkt 3: neues Gate läuft in beiden Ketten -------------------------------------
fixture m-positivliste
sed -i "s|\$(filter-out \$(GATE_SCHNELL_AUSNAHME),\$(GATE_CHECKS))|$SCHNELL|" "$w/harness/mk/gates-schnell.mk"
rot m-positivliste p3-neu

# --- Punkt 4: Deckung der Ausnahmeliste ---------------------------------------------
fixture m-eintrag-ohne-gate
sed -i "${AUSNAHME_ZEILE}s/\$/ phantom/" "$w/harness/mk/gates-schnell.mk"; rot m-eintrag-ohne-gate p4-eintrag exakt
fixture m-gegenprobe-in-ausnahme
sed -i "${AUSNAHME_ZEILE}s/\$/ lint-gegenprobe/" "$w/harness/mk/gates-schnell.mk"; rot m-gegenprobe-in-ausnahme p4-eintrag

# --- Punkt 5: ein grüner Lauf schreibt den Stand ------------------------------------
fixture m-kein-stand
sed -i '/^printf .*> "\$stand"$/d' "$w/tools/harness/gates-schnell.sh"; rot m-kein-stand p5-stand

# --- Punkt 6: roter oder abgebrochener Lauf lässt keinen Stand ----------------------
fixture m-stand-bleibt-nach-rot
sed -i 's|^  rm -f "\$stand"$|  :|' "$w/tools/harness/gates-schnell.sh"; rot m-stand-bleibt-nach-rot p6-rot
fixture m-entfernen-erst-bei-rot
sed -i "s|^  rm -f \"\\\$stand\"\$|  trap 'rm -f \"\$stand\"' ERR|" "$w/tools/harness/gates-schnell.sh"
rot m-entfernen-erst-bei-rot "p6-abbruch p6-i-entfernt" exakt

# --- Punkt 5: der Hash gilt zum Ende des Laufs --------------------------------------
fixture m-hash-vor-dem-lauf
sed -i -e '/^hash="\$(bash tools\/harness\/working-tree-hash.sh)"$/d' \
  -e 's|^"\$make_bin" "\$@"$|hash="$(bash tools/harness/working-tree-hash.sh)"\n&|' "$w/tools/harness/gates-schnell.sh"
rot m-hash-vor-dem-lauf p5-ende exakt

# --- Punkt 6: Make-Flags ------------------------------------------------------------
fixture m-flags-gestrichen
sed -i '/^case "\$buchstaben" in/d' "$w/tools/harness/gates-schnell.sh"
rot m-flags-gestrichen p6-n-schreibt
fixture m-n-nicht-erkannt
sed -i '/^case "\$buchstaben" in \*n\*)/d' "$w/tools/harness/gates-schnell.sh"
rot m-n-nicht-erkannt p6-n-schreibt
fixture m-q-nicht-erkannt
sed -i '/^case "\$buchstaben" in \*q\*)/d' "$w/tools/harness/gates-schnell.sh"
rot m-q-nicht-erkannt p6-q-entfernt
fixture m-t-nicht-erkannt
sed -i '/^case "\$buchstaben" in \*t\*)/d' "$w/tools/harness/gates-schnell.sh"
rot m-t-nicht-erkannt p6-t-schreibt
fixture m-trocken-entfernt-doch
sed -i 's|^if \[ "\$trocken" -eq 0 \]; then$|if true; then|' "$w/tools/harness/gates-schnell.sh"
rot m-trocken-entfernt-doch p6-n-entfernt
fixture m-trocken-schreibt-doch
sed -i 's|^if \[ "\$trocken" -eq 1 \] \|\| |if |' "$w/tools/harness/gates-schnell.sh"
rot m-trocken-schreibt-doch p6-n-schreibt
fixture m-i-nicht-erkannt
sed -i '/^case "\$buchstaben" in \*i\*)/d' "$w/tools/harness/gates-schnell.sh"
rot m-i-nicht-erkannt p6-i-schreibt
fixture m-i-wie-n
sed -i 's|\*i\*) ueberlesen=1|*i*) trocken=1|' "$w/tools/harness/gates-schnell.sh"
rot m-i-wie-n p6-i-entfernt
fixture m-i-vor-n
sed -i 's|\*i\*) ueberlesen=1 ;;|*i*) ueberlesen=1; trocken=0 ;;|' "$w/tools/harness/gates-schnell.sh"
rot m-i-vor-n p6-ni-entfernt exakt
fixture m-k-wie-n
sed -i 's|^case "\$buchstaben" in \*n\*)|case "$buchstaben" in *n* \| *k*)|' "$w/tools/harness/gates-schnell.sh"
rot m-k-wie-n p6-andere
fixture m-ganzes-makeflags
sed -i 's|^  \*) buchstaben="\$erstes" ;;|  *) buchstaben="$flags" ;;|' "$w/tools/harness/gates-schnell.sh"
rot m-ganzes-makeflags p5-stand

# --- Punkt 7: der volle Nachweis bleibt getrennt ------------------------------------
fixture m-stand-am-pfad-des-vollen-nachweises
sed -i 's/gates-schnell-passed/gates-passed/' "$w/tools/harness/gates-schnell.sh"
rot m-stand-am-pfad-des-vollen-nachweises p7-nachweis-neu
fixture m-loescht-den-vollen-nachweis
sed -i 's|^  rm -f "\$stand"$|  rm -f "$stand" .harness/state/gates-passed.diffsha|' "$w/tools/harness/gates-schnell.sh"
rot m-loescht-den-vollen-nachweis p7-nachweis

# --- Punkt 1: make gates enthält jedes Gate von GATE_CHECKS, nichts daneben ---------
fixture m-gate-aus-gates
sed -i 's/^record-gates: \$(GATE_CHECKS)$/record-gates: $(filter-out kopf-check,$(GATE_CHECKS))/' "$w/Makefile"
rot m-gate-aus-gates p1-voll
fixture m-kette-in-gates
sed -i 's/^record-gates: \$(GATE_CHECKS)$/record-gates: $(GATE_CHECKS) gates-schnell/' "$w/Makefile"
rot m-kette-in-gates p1-voll

fixture m-gegenprobe-in-gates
sed -i 's/^record-gates: \$(GATE_CHECKS)$/record-gates: $(GATE_CHECKS) lint-gegenprobe/' "$w/Makefile"
rot m-gegenprobe-in-gates p1-voll

# --- Anschluss: Liste und Exit ≠ 0 bei rotem Skript ---------------------------------
fixture m-ohne-anschluss
sed -i '/^GEGENPROBEN_CONTAINER += gates-schnell-gegenprobe$/d' "$w/harness/mk/gates-schnell.mk"; rot m-ohne-anschluss p2-v3-n1
fixture m-schnell-exit-ignoriert
sed -i 's|^\t@bash tools/harness/gates-schnell.sh|\t-@bash tools/harness/gates-schnell.sh|' "$w/harness/mk/gates-schnell.mk"
rot m-schnell-exit-ignoriert anschluss-exit-schnell
fixture m-gegenprobe-exit-ignoriert
sed -i 's|^\t@bash tools/harness/gates-schnell-gegenprobe.sh|\t-@bash tools/harness/gates-schnell-gegenprobe.sh|' "$w/harness/mk/gates-schnell.mk"
rot m-gegenprobe-exit-ignoriert anschluss-exit-gegenprobe

# --- SPEC-060 Punkte 1 und 2: die Namensregel, ein Fall je Verletzung und je Fall eine Mutation der Regel ---
# Jede Kopie hat beide Listen mit Einträgen (wie die Fragmente mit +=); die Regel fragt make in der Kopie.
declare -A RERW=() RDIR=() RENV=()
RNAMEN=()
# regelname <name> <erwarteter Befund> — ein Fall auf der Kopie in $w.
regelname() {
  RNAMEN+=("$1")
  RDIR[$1]="$w"
  RERW[$1]="${2:+$2 }"
}
# regelfall <name> <erwarteter Befund> <sed-Argument>… — eine Kopie, geändert durch sed auf fake.mk
# (oder auf der Datei in $rdatei); trifft sed keine Zeile, ist der Fall wirkungslos und die Gegenprobe rot.
regelfall() {
  local datei="$arbeit/r-$1"
  fixture "r-$1"
  cp "$w/${rdatei:-harness/mk/fake.mk}" "$datei.vor"
  sed -i "${@:3}" "$w/${rdatei:-harness/mk/fake.mk}"
  if cmp -s "$datei.vor" "$w/${rdatei:-harness/mk/fake.mk}"; then melde "Fall '$1' ändert die Kopie nicht"; fi
  regelname "$1" "$2"
}
# regelrot <regeldatei> — die Fälle, deren Befund vom erwarteten abweicht, sortiert in einer Zeile.
regelrot() {
  local n
  for n in "${RNAMEN[@]}"; do
    # shellcheck disable=SC2086
    [ "$( ( export ${RENV[$n]:-LC_ALL=C}; . "$1"; namensregel "${RDIR[$n]}" ) | sort -u | tr '\n' ' ')" = "${RERW[$n]}" ] || echo "$n"
  done | sort | tr '\n' ' '
}
# regelmut <name> <alt> <neu> <rote Fälle> — eine Kopie der Regel mit der Änderung; rot sind genau die benannten Fälle.
regelmut() {
  local datei="$arbeit/regel-$1.sh" inhalt neu ist
  inhalt="$(cat "$hier/gegenproben-regel.sh"; printf x)"; inhalt="${inhalt%x}"
  neu="${inhalt/"$2"/"$3"}"
  if [ "$neu" = "$inhalt" ]; then melde "Mutation '$1' ändert die Regel nicht: $2"; return; fi
  printf '%s' "$neu" > "$datei"
  ist="$(regelrot "$datei")"
  if [ "$ist" != "${4:+$4 }" ]; then melde "Mutation '$1': erwartet rot: ${4:-nichts}; rot sind: ${ist:-nichts}"; fi
}

fixture r-ok; regelname ok ""
regelfall p2-v1-n1 p2-v1-n1 -e 's/^GATE_CHECKS += .*/& lint-gegenprobe/'
regelfall p2-v1-n2 p2-v1-n2 -e 's/^GATE_CHECKS += .*/& a-check-negativ/'
regelfall p2-v2-n1 p2-v2-n1 -e 's/^GEGENPROBEN_HOST += .*/& lint-gegenprobe/'
regelfall p2-v2-n2 p2-v2-n2 -e 's/^GEGENPROBEN_HOST += .*/& a-check-negativ/'
regelfall p2-v3-n1 p2-v3-n1 -e 's/^\(GEGENPROBEN_CONTAINER += .*\) coverage-gegenprobe/\1/'
regelfall p2-v3-n2 p2-v3-n2 -e 's/^\(GEGENPROBEN_CONTAINER += .*\) a-check-negativ/\1/'
regelfall p2-v4-n1 p2-v4-n1 -e 's/^GEGENPROBEN_CONTAINER += .*/& lint-gegenprobe/'
regelfall p2-v4-n2 p2-v4-n2 -e 's/^GEGENPROBEN_CONTAINER += .*/& a-check-negativ/'
regelfall p2-v4-host p2-v4-n1 -e 's/^GEGENPROBEN_HOST += .*/& container-start-gegenprobe/'
regelfall p2-v6-container p2-v6-container -e 's/^GEGENPROBEN_CONTAINER += .*/& phantom/'
regelfall p2-v6-host p2-v6-host -e 's/^GEGENPROBEN_HOST += .*/& phantom/'
regelfall p2-andere-namen "" -e 's/^GATE_CHECKS += .*/& gegenproben gegenproben-host x-negativ x-gegenprobe-alt gegenprobe/' \
  -e '$a\gegenproben gegenproben-host x-negativ x-gegenprobe-alt gegenprobe: ; @:'
# Ein Ziel mit Namen auf -gegenprobe, dessen Rezept nicht aus einem Fragment unter harness/mk/ stammt, ist keine Gegenprobe im Sinn der Regel.
rdatei=Makefile regelfall p2-rezept-ausserhalb-fragment "" -e '$a\x-ausserhalb-gegenprobe: ; @:'
# Punkt 1: Ein Eintrag, den ein Fragment aus einer Variable bildet, steht in der Liste (ausgewertet).
regelfall p1-liste-aus-variable "" -e 's/^\(GEGENPROBEN_CONTAINER += .*\) coverage-gegenprobe/\1/' \
  -e '$a\GP_VAR = coverage-gegenprobe' -e '$a\GEGENPROBEN_CONTAINER += $(GP_VAR)'
# Die Rezept-Zeile von make -p steht je Locale in anderer Sprache; der Fall erzwingt die deutsche Ausgabe
# (LANGUAGE=de gilt, sobald die Locale nicht C ist) und verlangt den Befund der Kopie ohne Locale.
LOC="$(locale -a 2> /dev/null | grep -iE '^(de_DE|en_US)\.utf-?8$' | head -n 1 || true)"
fixture r-locale; regelname p2-locale ""
RENV[p2-locale]="LC_ALL=${LOC:-C} LANGUAGE=de"
if [ -z "$LOC" ] || ! (export LC_ALL="$LOC" LANGUAGE=de; { make -s -C "$w" -pq help 2> /dev/null || true; } | grep -q 'Rezept (aus'); then
  melde "Fall 'p2-locale': die Locale ${LOC:-(keine de_DE oder en_US installiert)} erzwingt nicht die deutsche Ausgabe von make -p"
fi

ist="$(regelrot "$hier/gegenproben-regel.sh")"
[ -z "$ist" ] || melde "Regel: Fälle rot, die grün sein müssen: $ist"

for v in 1 2 3 4; do
  # Der Fall p2-v4-host verletzt ebenfalls die Klasse n1 (sein Ziel endet auf -gegenprobe).
  rote="p2-v$v-n1"; if [ "$v" = 4 ]; then rote="p2-v4-host p2-v4-n1"; fi
  regelmut "p2-v$v-n1" "echo \"p2-v$v-\$k\"" "{ [ \"\$k\" = n2 ] && echo \"p2-v$v-\$k\"; true; }" "$rote"
  regelmut "p2-v$v-n2" "echo \"p2-v$v-\$k\"" "{ [ \"\$k\" = n1 ] && echo \"p2-v$v-\$k\"; true; }" "p2-v$v-n2"
done
regelmut p2-v4-nur-host 'if [ "$(gp_anzahl "$g" "${lc[@]}")" -gt 1 ] || ' 'if false || ' "p2-v4-n1 p2-v4-n2"
regelmut p2-v4-nur-container ' || [ "$(gp_anzahl "$g" "${lh[@]}")" -gt 1 ]; then' '; then' "p2-v4-host"
regelmut p2-v6-container 'echo "p2-v6-container"' ':' "p2-v6-container"
regelmut p2-v6-host 'echo "p2-v6-host"' ':' "p2-v6-host"
# p2-andere-namen: der Teilstring statt des Endes, dann je Name eine Mutation, die nur ihn als Gegenprobe liest.
regelmut p2-andere-teilstring '*-gegenprobe) echo n1' '*-gegenprobe*) echo n1' "p2-andere-namen"
for name in gegenproben gegenproben-host x-gegenprobe-alt gegenprobe; do
  regelmut "p2-andere-$name" '*-gegenprobe) echo n1' "*-gegenprobe | $name) echo n1" "p2-andere-namen"
done
regelmut p2-andere-x-negativ 'a-check-negativ) echo n2' 'a-check-negativ | x-negativ) echo n2' "p2-andere-namen"
regelmut p2-rezept-ohne-herkunft '.harness\/mk\/' '' "p2-rezept-ausserhalb-fragment"
# Punkt 1: die Regel liest den unausgewerteten Text von make -p.
ALT_LISTE='gp_make "$1" --eval "gp-lesen: ; @printf '"'"'%s\\n'"'"' \$($2)" gp-lesen | grep -v '"'"'^$'"'"' || true'
NEU_LISTE='gp_make "$1" -pq help | awk -v n="$2" '"'"'$1 == n && ($2 == ":=" || $2 == "=") { for (i = 3; i <= NF; i++) print $i }'"'"''
regelmut p1-liste-unausgewertet "$ALT_LISTE" "$NEU_LISTE" "p1-liste-aus-variable"
# Die Sprache von make ist nicht auf C gesetzt.
regelmut p2-locale-nicht-c 'LC_ALL=C env -u MAKEFLAGS' 'env -u MAKEFLAGS' "p2-locale"

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "gates-schnell-gegenprobe: gruen — Punkte 1 bis 7 von SPEC-052 (Punkt 5 mit Ende des Laufs, Punkt 6 mit Make-Flags) und Anschluss je Mutant rot, unveränderte Kopie grün, Ausnahmeliste deckt GATE_CHECKS; Namensregel von SPEC-060 Punkt 2 mit dreizehn Fällen der Regel, dem Fall p1-liste-aus-variable, dem Fall p2-locale und je Fall einer Mutation der Regel"
