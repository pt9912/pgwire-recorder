#!/usr/bin/env bash
# working-tree-hash — Hash über die Dateiliste des Arbeitsbaums (getrackte und untracked
# Dateien ohne ignorierte, Inhalt, Link-Ziel und Pfad je Eintrag; Form und Erzeugung in
# tools/harness/dateiliste.sh). Die Ausgabe wechselt genau dann, wenn sich die Liste
# ändert; ein Commit ohne Inhaltsänderung lässt sie gleich. Der Stand der schnellen Kette
# (gates-schnell.sh) und der Hash neben dem Nachweis des vollen Laufs (record-gates.sh)
# nutzen diese Ausgabe; die Logik der Liste steht nur in dateiliste.sh.
set -euo pipefail
wurzel="$(git rev-parse --show-toplevel)"
cd "$wurzel"
# shellcheck source=tools/harness/dateiliste.sh
. tools/harness/dateiliste.sh

hash="$(dateiliste_erzeugen | sha256sum)"
printf '%s\n' "${hash%% *}"
