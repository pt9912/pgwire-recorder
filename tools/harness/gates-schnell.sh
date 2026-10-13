#!/usr/bin/env bash
# gates-schnell — fährt die Ziele der schnellen Gate-Kette und schreibt bei grünem
# Lauf deren Stand (Vertrag: spec/spezifikation.md SPEC-052, Punkte 5 bis 7).
# Aufruf: gates-schnell.sh <make> <ziel>…
# Ein älterer Stand entfällt vor dem ersten Ziel; ein roter oder abgebrochener Lauf
# lässt deshalb keinen zurück. Der Stand ist die unveränderte Ausgabe von
# working-tree-hash.sh zum Ende des Laufs. Den Nachweis des vollen Laufs
# (.harness/state/gates-passed.diffsha) berührt dieses Skript nicht.
# Make-Flags (Punkt 6): Die Buchstaben im ersten Wort von MAKEFLAGS, wenn es nicht mit
# `-` beginnt, entscheiden vor dem Entfernen: bei n, q oder t läuft kein Gate wirklich,
# und der Stand bleibt unberührt (nichts entfernt, nichts geschrieben); bei i überliest
# make rote Gates, der Stand wird entfernt und nicht geschrieben. Andere Flags ändern nichts.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

make_bin="$1"
shift
stand=.harness/state/gates-schnell-passed.diffsha

flags="${MAKEFLAGS:-}"
erstes="${flags%% *}"
case "$erstes" in
  -*) buchstaben="" ;;
  *) buchstaben="$erstes" ;;
esac
trocken=0
ueberlesen=0
case "$buchstaben" in *n*) trocken=1 ;; esac
case "$buchstaben" in *q*) trocken=1 ;; esac
case "$buchstaben" in *t*) trocken=1 ;; esac
case "$buchstaben" in *i*) ueberlesen=1 ;; esac

mkdir -p .harness/state
if [ "$trocken" -eq 0 ]; then
  rm -f "$stand"
fi
"$make_bin" "$@"
if [ "$trocken" -eq 1 ] || [ "$ueberlesen" -eq 1 ]; then
  exit 0
fi
hash="$(bash tools/harness/working-tree-hash.sh)"
printf '%s\n' "$hash" > "$stand"
