#!/usr/bin/env bash
# stop-require-gates — wählt die Kette für den Claude-Stop: Der Stop ist frei, wenn die Dateiliste
# des Arbeitsbaums mit der Dateiliste des letzten grünen `make gates`-Laufs identisch ist, oder
# wenn jeder Unterschied zu ihr auf der Positivliste liegt (`inert`) und der Stand der schnellen
# Kette zum Arbeitsbaum passt; jeder andere Fall verlangt `make gates` (Antwort `block`).
# Liste und Vergleich: tools/harness/dateiliste.sh (dieselben Funktionen wie record-gates),
# Vertrag: spec/spezifikation.md SPEC-052.
#
# Zugesagt, soweit tools/harness/stop-hook-gegenprobe.sh es je Fall prüft:
# - Ein nicht aufgelöster Index-Eintrag (`git ls-files -u`) sperrt, als erste Prüfung nach
#   dem Schleifenschutz, auch bei sonst identischer Liste.
# - Ein fehlender, nicht lesbarer oder beschädigter Stand sperrt, auch bei sauberem Baum
#   (ein frischer Klon sperrt den ersten Stop, bis `make gates` lief).
# - Ein Eintrag, der fehlt, neu ist oder abweicht (Inhalt, Link-Ziel, gelöschte, umbenannte,
#   untracked Datei), sperrt, sofern sein Pfad nicht auf der Positivliste liegt.
# - Positivliste (`inert`): `docs/plan/planning/`, `docs/reviews/`, `docs/user/` und die
#   Wurzel-`README.md`, bytegenau und case-sensitiv; die Grenze steht mit dem Schrägstrich.
#   Ein Pfad mit absolutem Anfang oder mit Segment `.`, `..` oder leerem Segment im Unterschied
#   ist nie inert; ein geänderter Eintrag der Art X ist nie inert.
# - Ein Symlink im Unterschied ist inert, wenn Link-Pfad und das lexikalisch aufgelöste Ziel
#   auf der Positivliste liegen (ein absolutes Ziel, ein Ziel über die Wurzel hinaus und ein
#   Verzeichnis der Liste selbst sind es nicht); ein Link außerhalb der Positivliste, dessen
#   Ziel auf ihr liegt, ein Verzeichnis von ihr oder ein Vorfahr davon ist, macht jede
#   Änderung nicht inert.
# - Der Stand der schnellen Kette passt, wenn der Inhalt von
#   .harness/state/gates-schnell-passed.diffsha die Ausgabe von tools/harness/working-tree-hash.sh
#   ist; fehlend, leer, ungleich oder ein scheiterndes Skript sperrt. Fehlt der Stand, obwohl
#   jede Änderung inert ist, nennt der `reason` `make gates-schnell` (p10-fehlt); fehlt er bei
#   einer nicht inerten Änderung, nennt der `reason` nur `make gates` (p9-stand-fehlt-nicht-inert).
#   Der Wortlaut des `reason` bei einem Stand eines anderen Arbeitsbaums, einem leeren und einem
#   nicht lesbaren Stand ist nicht zugesagt.
# - Scheitert ein Schritt (git, Lesen einer Datei) oder bricht der Hook ab, antwortet er
#   `block`, nie ohne Antwort.
# - Der Schleifenschutz (`stop_hook_active: true`) antwortet `approve` vor jeder Prüfung, die zu
#   `block` führt (Konflikt, fehlender oder beschädigter Stand, scheiternder Schritt, fehlende
#   dateiliste.sh, fehlendes working-tree-hash.sh); `false` oder ein fehlendes Feld lässt
#   jede dieser Prüfungen gelten.
# Die Zusagen gelten für Läufe nacheinander: Änderung, `make gates` und Hook kommen nicht
# gleichzeitig.
set -euo pipefail

antwort=0
aktuell=""

# Die Antwort geht über Dateikennung 3 an den ursprünglichen stdout: Ein Abbruch mitten in
# einem Befehl mit Umleitung (`> "$aktuell"`) hielte sie sonst in dieser Umleitung fest.
exec 3>&1

# Jedes Ende ohne Antwort (Abbruch durch set -e, fehlende Datei) endet hier mit
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

# voll <grund> [<hinweis>] — der Hinweis nennt `make gates`, wenn kein anderer gegeben ist.
voll() {
  local hinweis='Run `make gates`.'
  [ $# -lt 2 ] || hinweis="$2"
  antwort=1
  printf '{"decision":"block","reason":"%s %s"}\n' "$1" "$hinweis" >&3
  exit 0
}

# Schleifen-Schutz: Hat dieser Hook den Stop bereits einmal blockiert (stop_hook_active),
# nicht erneut blockieren — sonst Endlosschleife bei dauerhaft rotem Gate.
input="$(cat || true)"
if [[ "$input" =~ \"stop_hook_active\"[[:space:]]*:[[:space:]]*true ]]; then
  frei
fi

export LC_ALL=C
wurzel="$(git rev-parse --show-toplevel)"
cd "$wurzel"
# shellcheck source=tools/harness/dateiliste.sh
. tools/harness/dateiliste.sh

declare -A ART=() WERT=()
SP=()
AP=()
DIFF=()
STAPEL=()

# inert <pfad> — Status 0, wenn der Pfad auf der Positivliste liegt. Die Positivliste steht nur hier.
inert() {
  case "$1" in
    docs/plan/planning/*) return 0 ;;
    docs/reviews/*) return 0 ;;
    docs/user/*) return 0 ;;
    README.md) return 0 ;;
  esac
  return 1
}

# zulaessig <pfad> — Status 1 für einen absoluten Pfad und für ein leeres Segment, `.` oder `..`.
zulaessig() {
  case "$1" in /*) return 1 ;; esac
  case "$1" in "" | *//* | */) return 1 ;; esac
  case "/$1/" in */./*) return 1 ;; esac
  case "/$1/" in */../*) return 1 ;; esac
  return 0
}

# stapeln <pfad> — wendet die Segmente des Pfades auf STAPEL an, ohne das Dateisystem zu befragen;
# Status 1, wenn `..` den Anfang überschreitet.
stapeln() {
  local rest="$1" seg fertig=0
  while [ "$fertig" -eq 0 ]; do
    case "$rest" in
      */*) seg="${rest%%/*}"; rest="${rest#*/}" ;;
      *) seg="$rest"; fertig=1 ;;
    esac
    case "$seg" in
      "" | .) ;;
      ..) [ "${#STAPEL[@]}" -gt 0 ] || return 1; unset 'STAPEL[${#STAPEL[@]}-1]' ;;
      *) STAPEL+=("$seg") ;;
    esac
  done
}

# ziel_lage <link-pfad> <ziel> — löst das Ziel eines Links lexikalisch auf. Setzt LAGE (innen oder
# aussen: außerhalb des Repos oder über die Wurzel hinaus), ABS (1 bei absolutem Ziel) und ZIEL (Pfad
# relativ zur Wurzel, leer für die Wurzel).
ziel_lage() {
  local link="$1" ziel="$2" n i s
  local -a ws
  LAGE=innen
  ABS=0
  ZIEL=""
  STAPEL=()
  case "$ziel" in
    /*)
      ABS=1
      stapeln "$wurzel"
      ws=("${STAPEL[@]}")
      STAPEL=()
      stapeln "$ziel" || { LAGE=aussen; return 0; }
      n=${#ws[@]}
      [ "${#STAPEL[@]}" -ge "$n" ] || { LAGE=aussen; return 0; }
      for ((i = 0; i < n; i++)); do
        [ "${STAPEL[i]}" = "${ws[i]}" ] || { LAGE=aussen; return 0; }
      done
      STAPEL=("${STAPEL[@]:n}")
      ;;
    *)
      if [[ "$link" == */* ]]; then stapeln "${link%/*}"; fi
      stapeln "$ziel" || { LAGE=aussen; return 0; }
      ;;
  esac
  for s in ${STAPEL[@]+"${STAPEL[@]}"}; do
    ZIEL="${ZIEL:+$ZIEL/}$s"
  done
}

# link_inert <link-pfad> <ziel> — Status 0, wenn das Ziel innerhalb des Repos, relativ und auf der
# Positivliste liegt; ein Verzeichnis der Liste selbst liegt nicht auf ihr.
link_inert() {
  ziel_lage "$1" "$2"
  [ "$LAGE" = innen ] || return 1
  [ "$ABS" -eq 0 ] || return 1
  inert "$ZIEL"
}

# spiegelt <pfad> — Status 0, wenn ein Ziel (Pfad relativ zur Wurzel) auf der Positivliste liegt,
# eines ihrer Verzeichnisse ist oder ein Vorfahr davon.
spiegelt() {
  inert "$1" && return 0
  case "$1" in
    docs/plan/planning | docs/reviews | docs/user) return 0 ;;
  esac
  case "$1" in
    docs | docs/plan) return 0 ;;
  esac
  case "$1" in
    "") return 0 ;;
  esac
  return 1
}

# einlesen <datei> <s|a> — liest die Einträge einer geprüften Liste in ART und WERT (Schlüssel
# `<s|a>:<pfad>`) und die Pfade in SP (Stand) oder AP (Arbeitsbaum).
einlesen() {
  local datei="$1" seite="$2" art wert pfad
  while IFS= read -r -d '' art; do
    IFS= read -r -d '' wert
    IFS= read -r -d '' pfad
    ART["$seite:$pfad"]="$art"
    WERT["$seite:$pfad"]="$wert"
    if [ "$seite" = s ]; then SP+=("$pfad"); else AP+=("$pfad"); fi
  done < "$datei"
}

# diff_bilden — DIFF: die Pfade, deren Eintrag im Arbeitsbaum neu ist oder abweicht, und die Pfade,
# die dort fehlen.
diff_bilden() {
  local p
  DIFF=()
  for p in ${AP[@]+"${AP[@]}"}; do
    [ "${ART["s:$p"]-}" = "${ART["a:$p"]-}" ] && [ "${WERT["s:$p"]-}" = "${WERT["a:$p"]-}" ] || DIFF+=("$p")
  done
  for p in ${SP[@]+"${SP[@]}"}; do
    [ -n "${ART["a:$p"]+x}" ] || DIFF+=("$p")
  done
}

# spiegel — ein Link außerhalb der Positivliste, dessen Ziel auf ihr liegt, macht jede Änderung nicht inert.
spiegel() {
  local p
  for p in ${AP[@]+"${AP[@]}"}; do
    [ "${ART["a:$p"]}" = L ] || continue
    ! inert "$p" || continue
    ziel_lage "$p" "${WERT["a:$p"]}"
    [ "$LAGE" = innen ] || continue
    ! spiegelt "$ZIEL" || voll "A symlink outside the inert paths points into them."
  done
  return 0
}

# pfadregel — gibt nur zurück, wenn jeder Unterschied inert ist; sonst `voll`.
pfadregel() {
  local p seite art
  diff_bilden
  for p in ${DIFF[@]+"${DIFF[@]}"}; do
    zulaessig "$p" || voll "The recorded file list is unreadable or damaged."
  done
  for p in ${DIFF[@]+"${DIFF[@]}"}; do
    inert "$p" || voll "A changed path is not on the inert list."
    for seite in s a; do
      art="${ART["$seite:$p"]-}"
      [ "$art" != X ] || voll "A changed entry of type X is never inert."
      [ "$art" != L ] || link_inert "$p" "${WERT["$seite:$p"]}" || voll "A changed symlink is not inert."
    done
  done
  spiegel
}

# schnell_gruen — gibt nur zurück, wenn der Stand der schnellen Kette der Hash des Arbeitsbaums ist.
schnell_gruen() {
  local hash gespeichert
  hash="$(bash tools/harness/working-tree-hash.sh)"
  [ -n "$hash" ] || voll "The working tree hash is empty."
  gespeichert="$(cat -- .harness/state/gates-schnell-passed.diffsha 2> /dev/null || :)"
  [ "$hash" = "$gespeichert" ] || voll "The fast gates chain has no passed run for this working tree." 'Run `make gates-schnell` (or `make gates`).'
}

# freigabe — Positivliste und Stand der schnellen Kette; endet immer in `frei` oder `voll`.
freigabe() {
  einlesen "$stand" s
  einlesen "$aktuell" a
  pfadregel
  schnell_gruen
  frei
}

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
freigabe
