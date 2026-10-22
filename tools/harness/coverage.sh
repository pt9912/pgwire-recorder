#!/bin/sh
# coverage.sh — Prüfschritt der Stufe `test` des Dockerfile: vergleicht die Anweisungsabdeckung
# eines Go-Profils (`-covermode=set`) mit der Schwelle.
# Aufruf: sh tools/harness/coverage.sh <Profil> <Schwelle>
# Vertrag: harness/vertraege/SPEC-056-testabdeckung-schwelle.md (Punkt 3 bis 7);
# Entscheidung ADR-0044. Läuft mit der POSIX-Shell und dem awk des Images der Stufe.
# GEPRÜFT DURCH tools/dev/harness/coverage-gegenprobe.sh.
#
# Ausgang 0: Gesamtwert nicht unter der Schwelle. Ausgang 1: darunter, oder eine Eingabe
# hat nicht die Form des Vertrags (Punkt 6); die Meldung auf stderr nennt die Form.
# Stdout trägt je Paket eine Zeile und zuletzt die Gesamtzeile (Punkt 7).
set -eu
export LC_ALL=C

fehler() {
  echo "coverage: $*" >&2
  exit 1
}

profil="${1-}"
schwelle="${2-}"

[ -f "$profil" ] && [ -r "$profil" ] || fehler "Profil fehlt oder ist nicht lesbar: '$profil'"
[ -s "$profil" ] || fehler "Profil ist leer: '$profil'"
[ "$(head -n 1 "$profil")" = "mode: set" ] || fehler "erste Zeile des Profils ist nicht 'mode: set'"
[ -n "$schwelle" ] || fehler "Schwelle fehlt"
case "$schwelle" in
  [1-9] | [1-9][0-9] | 100) ;;
  *) fehler "Schwelle ist keine ganze Zahl von 1 bis 100 ohne führende Null: '$schwelle'" ;;
esac

awk -v schwelle="$schwelle" '
function fehl(text) {
  print "coverage: " text | "cat 1>&2"
  close("cat 1>&2")
  bad = 1
  exit 1
}
# anzeige — Anteil a von n in Prozent, auf zwei Stellen abgerundet, in Hundertsteln.
function hundertstel(a, n) { return (a * 10000 - (a * 10000) % n) / n }
function text(q) { return sprintf("%d.%02d", int(q / 100), q % 100) }
BEGIN { FS = "[ ]" }
NR == 1 { next }
{
  if ($0 !~ /^[^ ]+:[^ ]+ [0-9]+ [0-9]+$/) {
    fehl("Zeile " NR " hat nicht die Form <Datei>:<Bereich> <Anzahl> <Zähler>: \"" $0 "\"")
  }
  if ($1 in gesehen) { fehl("Block kommt zweimal vor (Zeile " NR "): " $1) }
  gesehen[$1] = 1
  datei = $1
  sub(/:[^:]*$/, "", datei)
  paket = datei
  if (sub(/\/[^\/]*$/, "", paket) == 0) { paket = "." }
  if ($2 > 0) {
    anzahl[paket] += $2
    gesamt += $2
    if ($3 > 0) { ausgefuehrt[paket] += $2; ausgefuehrt_gesamt += $2 }
  }
}
END {
  if (bad) { exit 1 }
  if (gesamt == 0) { fehl("Summe der Anweisungen im Profil ist 0") }
  sortiert = "sort | cut -f2-"
  for (p in anzahl) {
    printf "%s\tcoverage: %s %s %% (%d Anweisungen)\n", p, p, text(hundertstel(ausgefuehrt[p], anzahl[p])), anzahl[p] | sortiert
  }
  close(sortiert)
  q = hundertstel(ausgefuehrt_gesamt, gesamt)
  gruen = (ausgefuehrt_gesamt * 100 >= schwelle * gesamt)
  zeile = sprintf("coverage: gesamt %s %% (%d von %d Anweisungen), Schwelle %d %%", text(q), ausgefuehrt_gesamt, gesamt, schwelle)
  if (!gruen) { zeile = zeile sprintf(", Differenz %s Prozentpunkte, unter der Schwelle", text(schwelle * 100 - q)) }
  print zeile
  exit (gruen ? 0 : 1)
}
' "$profil"
