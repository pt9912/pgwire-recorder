#!/usr/bin/env bash
# gates-schnell-gegenprobe — prüft die schnelle Gate-Kette (harness/mk/gates-schnell.mk,
# tools/harness/gates-schnell.sh) gegen die Punkte 1 bis 8 von SPEC-052
# (spec/spezifikation.md) an Kopien in einem Temp-Verzeichnis; der Arbeitsbaum
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
# Aufruf-Voraussetzung (Punkt 8) ist eine Bindung der Zusage und hat keinen Fall.
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
LANGSAM="lint-gegenprobe kopf-check-gegenprobe abdeckung-gegenprobe hook-gegenprobe a-check-negativ test-integration"

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
      *-gegenprobe | a-check-negativ | test-integration)
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
    "$wurzel/tools/harness/record-gates.sh" "$w/tools/harness/"
  printf '.harness/state/\n' > "$w/.gitignore"
  {
    printf 'FAKE_SCHNELL := %s neu-x\n' "$SCHNELL"
    printf 'FAKE_LANGSAM := %s\n' "$LANGSAM"
    printf 'FAKE := $(FAKE_SCHNELL) $(FAKE_LANGSAM)\n.PHONY: $(FAKE)\n$(FAKE):\n'
    printf '\t@echo $@ >> "$(LOG)"; if [ "$(FAIL)" = "$@" ]; then exit 1; fi; if [ "$(ABBRUCH)" = "$@" ]; then kill -KILL 0; fi\n'
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
    env -u MAKEFLAGS -u MAKELEVEL -u MFLAGS LOG="$arbeit/log" FAIL="${fail:-}" ABBRUCH="${abbruch:-}" \
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
  rm -f "$stand"; echo VOLL > "$voll"; fail=""; abbruch=""
  lauf "$w" gates-schnell
  [ "$code" -eq 0 ] || echo "lauf-gruen"
  for g in $LANGSAM; do grep -qxF "$g" "$arbeit/log" && echo "p2-langsam"; done
  for g in $SCHNELL; do grep -qxF "$g" "$arbeit/log" || echo "p2-schnell"; done
  grep -qxF neu-x "$arbeit/log" || echo "p3-neu"
  if [ ! -f "$stand" ] || [ "$(cat "$stand")" != "$(cd "$w" && bash tools/harness/working-tree-hash.sh)" ]; then
    echo "p5-stand"
  fi
  [ "$(cat "$voll" 2> /dev/null)" = VOLL ] || echo "p7-nachweis"

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
sed -i '/^rm -f "\$stand"$/d' "$w/tools/harness/gates-schnell.sh"; rot m-stand-bleibt-nach-rot p6-rot
fixture m-entfernen-erst-bei-rot
sed -i "s|^rm -f \"\\\$stand\"\$|trap 'rm -f \"\$stand\"' ERR|" "$w/tools/harness/gates-schnell.sh"
rot m-entfernen-erst-bei-rot p6-abbruch exakt

# --- Punkt 7: der volle Nachweis bleibt getrennt ------------------------------------
fixture m-stand-am-pfad-des-vollen-nachweises
sed -i 's/gates-schnell-passed/gates-passed/' "$w/tools/harness/gates-schnell.sh"
rot m-stand-am-pfad-des-vollen-nachweises p7-nachweis-neu

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
echo "gates-schnell-gegenprobe: gruen — Punkte 1 bis 7 und Anschluss je Mutant rot, unveränderte Kopie grün, Ausnahmeliste deckt GATE_CHECKS"
