#!/usr/bin/env bash
# gegenproben — fährt die Einzelziele einer Liste nacheinander und fasst zusammen (Vertrag:
# harness/vertraege/SPEC-060-gegenproben.md, Punkte 4 bis 12 und 15).
# Aufruf, nur aus harness/mk/gegenproben.mk: gegenproben.sh <make> <sammelziel> <ziel>…
# Kein Gate: kein Stand, kein Nachweis, nichts im Arbeitsbaum geschrieben. Reihenfolge der Gründe
# (Punkt 12): n, q, t · i · leere Liste · die Einzelziele. Die Buchstaben im ersten Wort von MAKEFLAGS,
# wenn es nicht mit `-` beginnt, entscheiden (wie tools/harness/gates-schnell.sh); k und j ändern nichts.
# GEPRÜFT DURCH tools/harness/gegenproben-gegenprobe.sh.
make_bin="$1"
sammel="$2"
shift 2

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
[ "$trocken" -eq 0 ] || exit 0
if [ "$ueberlesen" -eq 1 ]; then
  echo "$sammel: -i ist nicht zugesagt (ein rotes Einzelziel gäbe Exit 0); kein Einzelziel gefahren" >&2
  exit 2
fi
if [ "$#" -eq 0 ]; then
  echo "$sammel: Liste leer" >&2
  exit 2
fi

# beende <pid> <signal> — das Signal an den Prozess und an alle seine Nachkommen, Nachkommen zuerst.
beende() { local k; for k in $(pgrep -P "$1" 2> /dev/null); do beende "$k" "$2"; done; kill -s "$2" "$1" 2> /dev/null || true; }
# abbruch <signal> <exit> — beendet das laufende Einzelziel, startet kein weiteres, schreibt keine Zusammenfassung.
abbruch() { trap '' INT TERM; if [ -n "$pid" ]; then beende "$pid" "$1"; wait "$pid" 2> /dev/null; fi; exit "$2"; }
pid=""
trap 'abbruch INT 130' INT
trap 'abbruch TERM 143' TERM

zeilen=()
rot=0
for z in "$@"; do
  ( exec "$make_bin" "$z" ) &
  pid=$!
  wait "$pid"; rc=$?
  pid=""
  if [ "$rc" -eq 0 ]; then zeilen+=("$sammel: $z: grün"); else zeilen+=("$sammel: $z: rot (Exit $rc)"); rot=1; fi
done
printf '%s\n' "${zeilen[@]}" >&2
exit "$rot"
