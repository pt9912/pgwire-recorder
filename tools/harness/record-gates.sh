#!/usr/bin/env bash
# record-gates — Nachweis schreiben, dass `make gates` den aktuellen Arbeitsbaum-Zustand
# abgedeckt hat: die Dateiliste (.harness/state/gates-passed.dateiliste, Form und Erzeugung
# in tools/harness/dateiliste.sh) und deren Hash (.harness/state/gates-passed.diffsha, die
# Ausgabe von working-tree-hash.sh). Läuft als LETZTER gates-Prerequisite, also nur nach
# grünen Gates. Der Stop-Hook gibt nur frei, wenn die Dateiliste des Arbeitsbaums mit
# dieser identisch ist.
set -euo pipefail
wurzel="$(git rev-parse --show-toplevel)"
cd "$wurzel"
# shellcheck source=tools/harness/dateiliste.sh
. tools/harness/dateiliste.sh

mkdir -p .harness/state
liste="$(mktemp .harness/state/.dateiliste.XXXXXX)"
trap 'rm -f -- "$liste"' EXIT
dateiliste_erzeugen > "$liste"
hash="$(sha256sum < "$liste")"
printf '%s\n' "${hash%% *}" > .harness/state/gates-passed.diffsha
mv -f -- "$liste" "$DATEILISTE_STAND"
