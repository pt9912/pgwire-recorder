#!/usr/bin/env bash
# dateiliste — Funktionen der Dateiliste des vollen Laufs (Form: harness/vertraege/SPEC-052-gate-kette.md
# Punkt 11). Wird mit `source` geladen von working-tree-hash.sh, record-gates.sh
# und dem Stop-Hook; Schreiber, Hash und Leser teilen sich damit diese eine Logik.
# Die Funktionen laufen im Wurzelverzeichnis des Repos.
#
# Voraussetzung (nicht zugesagt, ohne Fall): bash ab Version 4 (`local -A`) und GNU-Werkzeuge
# (`readlink -z`, `sort -z`), wie das Repo GNU `make` voraussetzt. Fehlt sie, bricht der Lauf ab
# (SPEC-052 Punkt 20).
#
# GEPRÜFT DURCH tools/dev/harness/stop-hook-gegenprobe.sh.

# Pfad des Nachweises des vollen Laufs: die Dateiliste, neben gates-passed.diffsha.
DATEILISTE_STAND=".harness/state/gates-passed.dateiliste"

# dateiliste_erzeugen — schreibt die Liste des Arbeitsbaums auf stdout: je Pfad drei
# NUL-terminierte Felder Art, Wert, Pfad. Art F (Datei; Wert = SHA-256 des Inhalts),
# L (Symlink; Wert = Ziel), G (im Index, im Arbeitsbaum fehlend) oder X (anderer Typ),
# Wert bei G und X leer. Getrackte und untracked Pfade ohne ignorierte, mit `LC_ALL=C sort -zu`
# sortiert und dedupliziert (SPEC-052 Punkt 11: Form ohne eigenen Fall). Endet mit Status
# ungleich 0, wenn ein Schritt scheitert.
dateiliste_erzeugen() (
  set -euo pipefail
  git ls-files -z --cached --others --exclude-standard | LC_ALL=C sort -zu \
    | while IFS= read -r -d '' pfad; do
      if [ -L "$pfad" ]; then
        printf 'L\0'
        readlink -z -- "$pfad"
        printf '%s\0' "$pfad"
      elif [ -f "$pfad" ]; then
        hash="$(sha256sum < "$pfad")"
        printf 'F\0%s\0%s\0' "${hash%% *}" "$pfad"
      elif [ -e "$pfad" ]; then
        printf 'X\0\0%s\0' "$pfad"
      else
        printf 'G\0\0%s\0' "$pfad"
      fi
    done
)

# dateiliste_pruefen <datei> — Status 0, wenn die Datei eine unbeschädigte Liste ist.
# Beschädigt: leer oder nicht lesbar, Feldzahl nicht durch 3 teilbar, unbekannte Art,
# F mit Wert, der nicht 64 Hexzeichen ist, G oder X mit Wert, doppelter Pfad.
dateiliste_pruefen() {
  local datei="$1" art wert pfad
  local -A gesehen=()
  [ -s "$datei" ] || return 1
  while IFS= read -r -d '' art || [ -n "$art" ]; do
    IFS= read -r -d '' wert || return 1
    IFS= read -r -d '' pfad || return 1
    case "$art" in
      F) [[ "$wert" =~ ^[0-9a-fA-F]{64}$ ]] || return 1 ;;
      L) ;;
      G | X) [ -z "$wert" ] || return 1 ;;
      *) return 1 ;;
    esac
    [ -z "${gesehen["p:$pfad"]+x}" ] || return 1
    gesehen["p:$pfad"]=1
  done < "$datei"
}
