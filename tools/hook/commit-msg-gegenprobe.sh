#!/usr/bin/env bash
# commit-msg-gegenprobe — prueft den Traeger .githooks/commit-msg (ADR-0025) mit
# Message-Dateien in einem Temp-Verzeichnis; der Arbeitsbaum bleibt unberuehrt.
#
#   Message                                         erwartet
#   ohne jede Kennung                               abgelehnt
#   mit erfundener Slice-Kennung                    abgelehnt
#   mit dem Namen des Lifecycle-Werkzeugs allein    abgelehnt
#   mit der Kennung eines vorhandenen Slice         angenommen
#   mit einer Lastenheft-Kennung                    angenommen
#   Kennung nur in einer Kommentarzeile             abgelehnt
#
# Der vorhandene Slice ist der erste unter docs/plan/planning/*/; ohne Slice
# bricht die Gegenprobe ab. Ausgang: 0, wenn jeder Fall wie erwartet endet.
set -euo pipefail

traeger=".githooks/commit-msg"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT

vorhanden="$(find docs/plan/planning -mindepth 2 -maxdepth 2 -name 'slice-*.md' | sort | head -n 1)"
if [ -z "$vorhanden" ]; then
  echo "commit-msg-gegenprobe: kein Slice unter docs/plan/planning/*/ gefunden" >&2
  exit 1
fi
vorhanden="$(basename "$vorhanden" .md)"

fehler=0
nr=0

# fall <erwartung: an|ab> <message>
fall() {
  local erwartung="$1" message="$2" ergebnis=an
  nr=$((nr + 1))
  printf '%s\n' "$message" > "$arbeit/msg$nr"
  bash "$traeger" "$arbeit/msg$nr" >/dev/null 2>&1 || ergebnis=ab
  if [ "$ergebnis" != "$erwartung" ]; then
    echo "commit-msg-gegenprobe: ROT — erwartet $erwartung, Traeger war $ergebnis: $message" >&2
    fehler=1
  fi
}

fall ab "Arbeit ohne Kennung"
fall ab "slice-erfunden-gegenprobe: Arbeit"
fall ab "slice-mv: Arbeit"
fall an "$vorhanden: Arbeit"
fall an "Arbeit (LH-QA-04)"
fall ab "$(printf 'Arbeit ohne Kennung\n# %s' "$vorhanden")"

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "commit-msg-gegenprobe: gruen — vorhandene Slice-Kennung angenommen, erfundene und fehlende abgelehnt"
