#!/usr/bin/env bash
# commit-msg-gegenprobe — prueft den Traeger .githooks/commit-msg (ADR-0025,
# ADR-0029) mit Message-Dateien in einem Temp-Verzeichnis; der Arbeitsbaum bleibt
# unberuehrt.
#
#   Message                                                     erwartet
#   ohne jede Kennung                                           abgelehnt
#   mit erfundener Slice-Kennung                                abgelehnt
#   mit dem Namen des Lifecycle-Werkzeugs allein                abgelehnt
#   Kennung eines Slice in open/ im Betreff                     angenommen
#   Kennung eines Slice in done/ im Betreff                     angenommen
#   Kennung eines vorhandenen Slice nur im Rumpf                angenommen
#   Kennung mit angehaengtem Suffix                             abgelehnt
#   Kennung mit vorangestelltem Wortteil                        abgelehnt
#   Kennung nur in einer Kommentarzeile                         abgelehnt
#   Kennung nur unter der Scissors-Zeile von `git commit -v`    abgelehnt
#   mit einer Lastenheft-Kennung                                angenommen
#   Kennung einer offenen Welle (flach in planning/)            angenommen
#   Kennung einer geschlossenen Welle (done/ oder Archiv)       angenommen
#   mit erfundener Welle-Kennung                                abgelehnt
#   Name einer Closure-Notiz einer Welle (-results)             abgelehnt
#
# Slices und Wellen der Positivfaelle legt die Gegenprobe selbst in einem Temp-Repo an
# (feste Namen unter docs/plan/planning/, im Index); sie liest den Planungsbaum des
# Arbeitsbaums nicht. Der Traeger und die Pruefung der Traceability werden aus dem
# Arbeitsbaum dorthin kopiert. Ausgang: 0, wenn jeder Fall wie erwartet endet.
set -euo pipefail

arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT

repo="$arbeit/repo"
mkdir -p "$repo/.githooks" "$repo/tools/harness"
cp .githooks/commit-msg "$repo/.githooks/commit-msg"
cp tools/harness/commit-msg-traceability.sh "$repo/tools/harness/commit-msg-traceability.sh"
planung="$repo/docs/plan/planning"
mkdir -p "$planung/open" "$planung/done"
offen="slice-gegenprobe-offen"
fertig="slice-gegenprobe-fertig"
welle_offen="welle-gegenprobe-offen"
welle_fertig="welle-gegenprobe-fertig"
notiz="$welle_fertig-results"
: > "$planung/open/$offen.md"
: > "$planung/done/$fertig.md"
: > "$planung/$welle_offen.md"
: > "$planung/done/$welle_fertig.md"
: > "$planung/done/$notiz.md"
git -C "$repo" init -q
git -C "$repo" add -A
traeger="$repo/.githooks/commit-msg"

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

schere='# ------------------------ >8 ------------------------'

fall ab "Arbeit ohne Kennung"
fall ab "slice-erfunden-gegenprobe: Arbeit"
fall ab "slice-mv: Arbeit"
fall an "$offen: Arbeit"
fall an "$fertig: Arbeit"
fall an "$(printf 'Arbeit ohne Kennung im Betreff\n\nBetrifft %s.' "$offen")"
fall ab "$offen-zusatz: Arbeit"
fall ab "x$offen: Arbeit"
fall ab "$(printf 'Arbeit ohne Kennung\n# %s' "$offen")"
fall ab "$(printf 'Arbeit ohne Kennung\n%s\ndiff --git a/x b/x\n+%s' "$schere" "$offen")"
fall an "Arbeit (LH-QA-04)"
fall an "$welle_offen: Arbeit"
fall an "$welle_fertig: Arbeit"
fall ab "welle-erfunden-gegenprobe: Arbeit"
fall ab "$notiz: Arbeit"

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "commit-msg-gegenprobe: gruen — vorhandene Slice- und Welle-Kennungen angenommen, erfundene, verlaengerte, kommentierte, Closure-Notizen und fehlende abgelehnt"
