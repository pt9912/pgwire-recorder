#!/usr/bin/env bash
# gates-schnell-gegenprobe — prüft die schnelle Gate-Kette (harness/mk/gates-schnell.mk,
# tools/harness/gates-schnell.sh) gegen die Punkte 1 bis 7 von SPEC-052
# (harness/vertraege/SPEC-052-gate-kette.md) an Kopien in einem Temp-Verzeichnis; der Arbeitsbaum
# bleibt unberührt. Eine Kopie trägt Makefile, Fragment der Kette, Skripte der Kette
# und ein Fragment mit Gates, die ihren Namen in eine Logdatei schreiben; die
# unveränderte Kopie muss jede Prüfung bestehen, und je Mutant (sed auf der Kopie)
# muss die benannte Prüfung rot werden. Die Deckung der Ausnahmeliste (Punkt 4) läuft
# zusätzlich gegen den Arbeitsbaum selbst: ein neues Gate dieser Namen macht die
# Gegenprobe rot, bis jemand es in die Liste aufnimmt oder bewusst schnell lässt.
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
LANGSAM="lint-gegenprobe coverage-gegenprobe container-start-gegenprobe kopf-check-gegenprobe abdeckung-gegenprobe hook-gegenprobe integration-volumes-gegenprobe meldungskatalog-check-gegenprobe stop-hook-gegenprobe a-check-negativ test-container test-container-gegenprobe test-integration test-postgres-versionen test-postgres-versionen-gegenprobe"

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

# deckung <wurzel> — Punkt 4: schreibt die Namen der verletzten Prüfungen.
deckung() {
  local ausnahme gates g
  mapfile -t ausnahme < <(gate_variable "$1" GATE_SCHNELL_AUSNAHME)
  mapfile -t gates < <(gate_variable "$1" GATE_CHECKS)
  for g in "${ausnahme[@]}"; do
    in_liste "$g" "${gates[@]}" || echo "p4-eintrag"
  done
  for g in "${gates[@]}"; do
    case "$g" in
      *-gegenprobe | a-check-negativ | test-container | test-integration | test-postgres-versionen)
        in_liste "$g" "${ausnahme[@]}" || echo "p4-name" ;;
    esac
  done
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
    printf 'FAKE := $(FAKE_SCHNELL) $(FAKE_LANGSAM)\n.PHONY: $(FAKE)\n$(FAKE):\n'
    printf '\t@echo $@ >> "$(LOG)"; if [ "$(FAIL)" = "$@" ]; then exit 1; fi; if [ "$(ABBRUCH)" = "$@" ]; then kill -KILL 0; fi; if [ "$(AENDERE)" = "$@" ]; then echo x >> aenderung; fi\n'
    printf 'GATE_CHECKS += $(FAKE)\n'
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
    case "$g" in
      gates-schnell-gegenprobe) grep -q 'tools/harness/gates-schnell-gegenprobe.sh' "$arbeit/out" || echo "p1-voll" ;;
      *) grep -qF "echo $g >> " "$arbeit/out" || echo "p1-voll" ;;
    esac
  done
  grep -q 'tools/harness/gates-schnell.sh' "$arbeit/out" && echo "p1-voll"
  grep -q 'tools/harness/record-gates.sh' "$arbeit/out" || echo "p1-voll"

  # Anschluss: die Gegenprobe hängt an GATE_CHECKS, beide Ziele enden bei rotem Skript ≠ 0.
  in_liste gates-schnell-gegenprobe "${gates[@]}" || echo "anschluss-gegenprobe"
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

# --- Unveränderte Kopie, und Deckung im Arbeitsbaum --------------------------------
fixture ok; gruen ok
ist="$(deckung "$wurzel" | sort -u | tr '\n' ' ')"
[ -z "$ist" ] || melde "Arbeitsbaum: Ausnahmeliste und GATE_CHECKS decken sich nicht: $ist"

# --- Punkt 2: kein langsames Gate ---------------------------------------------------
fixture m-lint-gegenprobe
sed -i "${AUSNAHME_ZEILE}s/ lint-gegenprobe//" "$w/harness/mk/gates-schnell.mk"; rot m-lint-gegenprobe p2-langsam
fixture m-test-integration
sed -i "${AUSNAHME_ZEILE}s/ test-integration//" "$w/harness/mk/gates-schnell.mk"; rot m-test-integration p2-langsam
fixture m-test-container
sed -i "${AUSNAHME_ZEILE}s/ test-container / /" "$w/harness/mk/gates-schnell.mk"; rot m-test-container p2-langsam; rot m-test-container p4-name
fixture m-test-container-gegenprobe
sed -i "${AUSNAHME_ZEILE}s/ test-container-gegenprobe//" "$w/harness/mk/gates-schnell.mk"; rot m-test-container-gegenprobe p2-langsam
fixture m-test-postgres-versionen
sed -i "${AUSNAHME_ZEILE}s/ test-postgres-versionen / /" "$w/harness/mk/gates-schnell.mk"; rot m-test-postgres-versionen p2-langsam; rot m-test-postgres-versionen p4-name
fixture m-test-postgres-versionen-gegenprobe
sed -i "${AUSNAHME_ZEILE}s/ test-postgres-versionen-gegenprobe//" "$w/harness/mk/gates-schnell.mk"; rot m-test-postgres-versionen-gegenprobe p2-langsam

# --- Punkt 2: jedes andere Gate -----------------------------------------------------
fixture m-schnelles-gate-in-liste
sed -i "${AUSNAHME_ZEILE}s/\$/ kopf-check/" "$w/harness/mk/gates-schnell.mk"; rot m-schnelles-gate-in-liste p2-schnell

# --- Punkt 3: neues Gate läuft in beiden Ketten -------------------------------------
fixture m-positivliste
sed -i "s|\$(filter-out \$(GATE_SCHNELL_AUSNAHME),\$(GATE_CHECKS))|$SCHNELL|" "$w/harness/mk/gates-schnell.mk"
rot m-positivliste p3-neu

# --- Punkt 4: Deckung der Ausnahmeliste ---------------------------------------------
fixture m-eintrag-ohne-gate
sed -i "${AUSNAHME_ZEILE}s/\$/ phantom/" "$w/harness/mk/gates-schnell.mk"; rot m-eintrag-ohne-gate p4-eintrag exakt
fixture m-gate-ohne-eintrag
sed -i 's/^FAKE_LANGSAM := /FAKE_LANGSAM := x-gegenprobe /' "$w/harness/mk/fake.mk"; rot m-gate-ohne-eintrag p4-name

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

# --- Anschluss: GATE_CHECKS und Exit ≠ 0 bei rotem Skript ---------------------------
fixture m-ohne-anschluss
sed -i '/^GATE_CHECKS += gates-schnell-gegenprobe$/d' "$w/harness/mk/gates-schnell.mk"; rot m-ohne-anschluss anschluss-gegenprobe
fixture m-schnell-exit-ignoriert
sed -i 's|^\t@bash tools/harness/gates-schnell.sh|\t-@bash tools/harness/gates-schnell.sh|' "$w/harness/mk/gates-schnell.mk"
rot m-schnell-exit-ignoriert anschluss-exit-schnell
fixture m-gegenprobe-exit-ignoriert
sed -i 's|^\t@bash tools/harness/gates-schnell-gegenprobe.sh|\t-@bash tools/harness/gates-schnell-gegenprobe.sh|' "$w/harness/mk/gates-schnell.mk"
rot m-gegenprobe-exit-ignoriert anschluss-exit-gegenprobe

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "gates-schnell-gegenprobe: gruen — Punkte 1 bis 7 (Punkt 5 mit Ende des Laufs, Punkt 6 mit Make-Flags) und Anschluss je Mutant rot, unveränderte Kopie grün, Ausnahmeliste deckt GATE_CHECKS"
