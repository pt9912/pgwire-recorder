#!/usr/bin/env bash
# gegenproben-regel — Namens- und Zuordnungsregel der Gegenproben (harness/vertraege/SPEC-060-gegenproben.md,
# Punkt 2). Kein Make-Ziel und kein Gate: tools/dev/harness/gates-schnell-gegenprobe.sh lädt die Datei und ruft
# `namensregel` für den Arbeitsbaum und für jede Kopie.
# namensregel <wurzel>: Die Funktion fragt `make` im Verzeichnis <wurzel> und schreibt je Verletzung den Namen
# des Falls, einen je Zeile, und nichts, wenn die Regel gilt. Die Listen liest sie ausgewertet (ein Eintrag darf
# aus einer Variable gebildet sein, SPEC-060 Punkt 1), die Ziele mit Rezept aus `make -pq`. Jeder Aufruf von
# `make` läuft mit LC_ALL=C: `make -p` schreibt die Rezept-Zeile je Locale in anderer Sprache. Als Gegenprobe gilt jedes Ziel mit Rezept in einem
# Fragment unter harness/mk/, dessen Name auf -gegenprobe endet (Klasse n1), und a-check-negativ (Klasse n2):
#   p2-v1-<k>  das Ziel steht in GATE_CHECKS
#   p2-v2-<k>  es steht in beiden Listen
#   p2-v3-<k>  es steht in keiner Liste
#   p2-v4-<k>  es steht zweimal in derselben Liste
#   p2-v6-container, p2-v6-host  ein Eintrag der Liste hat kein Ziel mit Rezept
# GEPRÜFT DURCH tools/dev/harness/gates-schnell-gegenprobe.sh: die Fälle p2-* und je Fall eine Mutation dieser Datei.

# gp_make <wurzel> <make-argument>… — make in <wurzel>, ohne Umgebung von make, in der Sprache C.
gp_make() {
  LC_ALL=C env -u MAKEFLAGS -u MFLAGS -u MAKELEVEL make -s --no-print-directory -C "$1" "${@:2}" 2> /dev/null || true
}

# gp_variable <wurzel> <name> — Wörter der Make-Variablen, ausgewertet.
gp_variable() {
  gp_make "$1" --eval "gp-lesen: ; @printf '%s\\n' \$($2)" gp-lesen | grep -v '^$' || true
}

# gp_rezeptziele <wurzel> — Namen der Ziele, deren Rezept aus einem Fragment unter harness/mk/ stammt.
gp_rezeptziele() {
  gp_make "$1" -pq help | awk '/^[A-Za-z0-9_.%-]+:( |$)/ { z = $1; sub(/:$/, "", z); next }
       /^#  recipe to execute \(from .harness\/mk\// { if (z != "") print z }'
}

# gp_klasse <ziel> — n1 für einen Namen auf -gegenprobe, n2 für a-check-negativ, sonst nichts.
gp_klasse() {
  case "$1" in
    *-gegenprobe) echo n1 ;;
    a-check-negativ) echo n2 ;;
  esac
}

# gp_anzahl <wort> <wort>… — wie oft das erste Wort unter den übrigen steht.
gp_anzahl() {
  local g="$1" x n=0
  shift
  for x in "$@"; do
    if [ "$x" = "$g" ]; then n=$((n + 1)); fi
  done
  echo "$n"
}

namensregel() {
  local w="$1" g k e
  local -a gates lc lh ziele
  mapfile -t gates < <(gp_variable "$w" GATE_CHECKS)
  mapfile -t lc < <(gp_variable "$w" GEGENPROBEN_CONTAINER)
  mapfile -t lh < <(gp_variable "$w" GEGENPROBEN_HOST)
  mapfile -t ziele < <(gp_rezeptziele "$w" | sort -u)
  for g in "${ziele[@]}"; do
    k="$(gp_klasse "$g")"
    if [ -z "$k" ]; then continue; fi
    if [ "$(gp_anzahl "$g" "${gates[@]}")" -gt 0 ]; then echo "p2-v1-$k"; fi
    if [ "$(gp_anzahl "$g" "${lc[@]}")" -gt 0 ] && [ "$(gp_anzahl "$g" "${lh[@]}")" -gt 0 ]; then echo "p2-v2-$k"; fi
    if [ "$(gp_anzahl "$g" "${lc[@]}")" -eq 0 ] && [ "$(gp_anzahl "$g" "${lh[@]}")" -eq 0 ]; then echo "p2-v3-$k"; fi
    if [ "$(gp_anzahl "$g" "${lc[@]}")" -gt 1 ] || [ "$(gp_anzahl "$g" "${lh[@]}")" -gt 1 ]; then echo "p2-v4-$k"; fi
  done
  for e in "${lc[@]}"; do
    if [ "$(gp_anzahl "$e" "${ziele[@]}")" -eq 0 ]; then echo "p2-v6-container"; fi
  done
  for e in "${lh[@]}"; do
    if [ "$(gp_anzahl "$e" "${ziele[@]}")" -eq 0 ]; then echo "p2-v6-host"; fi
  done
  return 0
}
