#!/usr/bin/env bash
# gates-schnell — fährt die Ziele der schnellen Gate-Kette und schreibt bei grünem
# Lauf deren Stand (Vertrag: spec/spezifikation.md SPEC-052, Punkte 5 bis 7).
# Aufruf: gates-schnell.sh <make> <ziel>…
# Ein älterer Stand entfällt vor dem ersten Ziel; ein roter oder abgebrochener Lauf
# lässt deshalb keinen zurück. Der Stand ist die unveränderte Ausgabe von
# working-tree-hash.sh zum Ende des Laufs. Den Nachweis des vollen Laufs
# (.harness/state/gates-passed.diffsha) berührt dieses Skript nicht.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

make_bin="$1"
shift
stand=.harness/state/gates-schnell-passed.diffsha

mkdir -p .harness/state
rm -f "$stand"
"$make_bin" "$@"
hash="$(bash tools/harness/working-tree-hash.sh)"
printf '%s\n' "$hash" > "$stand"
