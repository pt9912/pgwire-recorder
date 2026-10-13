#!/usr/bin/env bash
# record-gates — Nachweis schreiben, dass `make gates` den aktuellen Arbeitsbaum-Zustand
# abgedeckt hat: die Dateiliste (.harness/state/gates-passed.dateiliste, Form und Erzeugung
# in tools/harness/dateiliste.sh) und deren Hash (.harness/state/gates-passed.diffsha, die
# Ausgabe von working-tree-hash.sh). Läuft als LETZTER gates-Prerequisite, also nur nach
# grünen Gates. Der Stop-Hook gibt frei, wenn die Dateiliste des Arbeitsbaums mit dieser
# identisch ist oder wenn jeder Unterschied auf der Positivliste des Hooks liegt und der Stand
# der schnellen Kette zum Arbeitsbaum passt (harness/vertraege/SPEC-052-gate-kette.md Punkt 9; stop-hook-gegenprobe: p16-identisch,
# p9-freigabe).
# Make-Flags (SPEC-052 Punkt 21): Die Buchstaben im ersten Wort von MAKEFLAGS, wenn es nicht
# mit `-` beginnt, entscheiden wie in gates-schnell.sh: bei i überliest make rote Gates, der
# Nachweis wird nicht geschrieben und ein älterer entfernt. -n, -q und -t prüft dieses Skript
# nicht: das Rezept in harness/mk/enforce.mk trägt weder $(MAKE) noch `+`, make führt es dort
# nicht aus (stop-hook-gegenprobe: p21-n, p21-q, p21-t, p21-ni).
set -euo pipefail
wurzel="$(git rev-parse --show-toplevel)"
cd "$wurzel"
# shellcheck source=tools/harness/dateiliste.sh
. tools/harness/dateiliste.sh

flags="${MAKEFLAGS:-}"
erstes="${flags%% *}"
case "$erstes" in
  -*) buchstaben="" ;;
  *) buchstaben="$erstes" ;;
esac
case "$buchstaben" in
  *i*)
    rm -f -- .harness/state/gates-passed.diffsha "$DATEILISTE_STAND"
    exit 0
    ;;
esac

mkdir -p .harness/state
liste="$(mktemp .harness/state/.dateiliste.XXXXXX)"
trap 'rm -f -- "$liste"' EXIT
dateiliste_erzeugen > "$liste"
hash="$(sha256sum < "$liste")"
printf '%s\n' "${hash%% *}" > .harness/state/gates-passed.diffsha
mv -f -- "$liste" "$DATEILISTE_STAND"
