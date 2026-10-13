#!/usr/bin/env bash
# stop-require-gates — gibt den Claude-Stop nur frei, wenn die Dateiliste des Arbeitsbaums
# mit der Dateiliste des letzten grünen `make gates`-Laufs identisch ist; jeder andere
# Fall verlangt `make gates` (Antwort `block`). Liste und Vergleich: tools/harness/dateiliste.sh
# (dieselben Funktionen wie record-gates), Vertrag: spec/spezifikation.md SPEC-052.
#
# Zugesagt, soweit tools/harness/stop-hook-gegenprobe.sh es je Fall prüft:
# - Ein nicht aufgelöster Index-Eintrag (`git ls-files -u`) sperrt, als erste Prüfung nach
#   dem Schleifenschutz, auch bei sonst identischer Liste.
# - Ein fehlender, nicht lesbarer oder beschädigter Stand sperrt, auch bei sauberem Baum
#   (ein frischer Klon sperrt den ersten Stop, bis `make gates` lief).
# - Ein Eintrag, der fehlt, neu ist oder abweicht (Inhalt, Link-Ziel, gelöschte, umbenannte,
#   untracked Datei), sperrt.
# - Scheitert ein Schritt (git, Lesen einer Datei) oder bricht der Hook ab, antwortet er
#   `block`, nie ohne Antwort.
# - Der Schleifenschutz (`stop_hook_active`) antwortet `approve`, damit ein dauerhaft rotes
#   Gate keine Endlosschleife erzeugt.
# Die Zusagen gelten für Läufe nacheinander: Änderung, `make gates` und Hook kommen nicht
# gleichzeitig.
set -euo pipefail

antwort=0
aktuell=""

# Die Antwort geht über Dateikennung 3 an den ursprünglichen stdout: Ein Abbruch mitten in
# einem Befehl mit Umleitung (`> "$aktuell"`) hielte sie sonst in dieser Umleitung fest.
exec 3>&1

# Jedes Ende ohne Antwort (Abbruch durch set -e, fehlende Datei, Signal) endet hier mit
# Antwort `block` und Status 0.
sicherung() {
  if [ "$antwort" -eq 0 ]; then
    printf '%s\n' '{"decision":"block","reason":"The stop hook aborted. Run `make gates`."}' >&3
  fi
  if [ -n "$aktuell" ]; then
    rm -f -- "$aktuell" || :
  fi
  exit 0
}
trap sicherung EXIT

frei() {
  antwort=1
  printf '%s\n' '{"decision":"approve"}' >&3
  exit 0
}

voll() {
  antwort=1
  printf '{"decision":"block","reason":"%s Run `make gates`."}\n' "$1" >&3
  exit 0
}

# Schleifen-Schutz: Hat dieser Hook den Stop bereits einmal blockiert (stop_hook_active),
# nicht erneut blockieren — sonst Endlosschleife bei dauerhaft rotem Gate.
input="$(cat || true)"
if [[ "$input" =~ \"stop_hook_active\"[[:space:]]*:[[:space:]]*true ]]; then
  frei
fi

wurzel="$(git rev-parse --show-toplevel)"
cd "$wurzel"
# shellcheck source=tools/harness/dateiliste.sh
. tools/harness/dateiliste.sh

konflikt="$(git ls-files -u)"
[ -z "$konflikt" ] || voll "There are unresolved merge entries in the index."

stand="$DATEILISTE_STAND"
[ -f "$stand" ] || voll "No recorded successful make gates run (file list missing)."
dateiliste_pruefen "$stand" || voll "The recorded file list is unreadable or damaged."

aktuell="$(mktemp)"
dateiliste_erzeugen > "$aktuell"
if cmp -s "$stand" "$aktuell"; then
  frei
fi
voll "The working tree differs from the last recorded gates run (commits without a gates run count too)."
