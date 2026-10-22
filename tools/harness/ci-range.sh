#!/usr/bin/env bash
# ci-range — berechnet die Commit-Range eines Ereignisses der CI (harness/vertraege/SPEC-064-ci.md, Teil 2, Punkte 6 bis 9).
# Eingaben (Umgebung): CI_EVENT (pull_request oder push), CI_HEAD, CI_PR_BASE (nur pull_request), CI_PUSH_BEFORE (nur push),
# CI_DEFAULT_BRANCH (nur push mit neuem Zweig oder unbekanntem Vorgänger). Ausgabe: auf stdout genau eine Zeile `<Basis>..<Kopf>` mit
# aufgelösten SHAs; Hinweise und Fehler stehen mit dem Präfix `ci-range:` auf stderr. Exit 0: nicht leere Range; 1: leere Range;
# 2: Eingabe, Klon oder Auflösung nicht möglich. Bei Exit ungleich 0 bleibt stdout leer. Läuft im Klon des Aufrufers (Arbeitsverzeichnis).
# Braucht bash und git; eine Version prüft das Skript nicht, zugesagt ist die Plattform von SPEC-064.
# GEPRÜFT DURCH tools/dev/harness/ci-gegenprobe.sh (Fälle p6-* bis p9-*, je Zusage eine Mutation dieser Datei).
set -uo pipefail
export LC_ALL=C

fehler() { echo "ci-range: $*" >&2; exit 2; }
ungueltig() { fehler "$1 fehlt oder ist ungültig"; }
aufloesen() { git rev-parse --verify --quiet "$1^{commit}" 2> /dev/null; }
HEX='^[0-9a-f]{40}$'
NULLEN='^0{40}$'
# pruefe_standard — prüft CI_DEFAULT_BRANCH und setzt standard.
pruefe_standard() {
  standard="${CI_DEFAULT_BRANCH-}"
  [ -n "$standard" ] && [ "$(git check-ref-format --branch "$standard" 2> /dev/null)" = "$standard" ] || ungueltig CI_DEFAULT_BRANCH
}
# vom_standard — setzt basis auf die Merge-Base von CI_HEAD und dem Standard-Zweig (nicht auf dessen Spitze); Exit 2, wenn der Zweig fehlt oder kein gemeinsamer Vorfahr besteht.
vom_standard() {
  basis="$(git merge-base "$kopf_roh" "refs/remotes/origin/$standard" 2> /dev/null)" || fehler "keine Merge-Base von $kopf_roh und refs/remotes/origin/$standard (Zweig fehlt im Klon oder kein gemeinsamer Vorfahr)"
}

# (1) Ereignis
event="${CI_EVENT-}"
case "$event" in
  pull_request | push) ;;
  *) ungueltig CI_EVENT ;;
esac

# (2) Variablen, die das Ereignis liest
kopf_roh="${CI_HEAD-}"
[[ $kopf_roh =~ $HEX ]] || ungueltig CI_HEAD
neu=0
if [ "$event" = pull_request ]; then
  basis_roh="${CI_PR_BASE-}"
  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PR_BASE
else
  basis_roh="${CI_PUSH_BEFORE-}"
  [[ $basis_roh =~ $HEX ]] || ungueltig CI_PUSH_BEFORE
  if [[ $basis_roh =~ $NULLEN ]]; then
    neu=1
    pruefe_standard
  fi
fi

# (3) flacher Klon
flach="$(git rev-parse --is-shallow-repository 2> /dev/null)" || fehler "git rev-parse --is-shallow-repository schlug fehl"
case "$flach" in
  true) fehler "flacher Klon, der Checkout braucht fetch-depth: 0" ;;
  false) ;;
  *) fehler "git rev-parse --is-shallow-repository ergab: $flach" ;;
esac

# (4) Auflösung von Basis und Kopf
if [ "$neu" -eq 1 ]; then
  vom_standard
elif ! basis="$(aufloesen "$basis_roh")"; then
  [ "$event" = push ] || fehler "Basis $basis_roh nicht im Klon"
  echo "ci-range: Vorgänger $basis_roh nicht im Klon, Basis ist der Standard-Zweig" >&2
  pruefe_standard
  vom_standard
fi
kopf="$(aufloesen "$kopf_roh")" || fehler "Kopf $kopf_roh nicht im Klon"

# (5) leere Range
anzahl="$(git rev-list --count "$basis..$kopf" 2> /dev/null)" || fehler "git rev-list --count $basis..$kopf schlug fehl"
if [ "$anzahl" -eq 0 ]; then
  echo "ci-range: Range $basis..$kopf ist leer" >&2
  exit 1
fi
echo "$basis..$kopf"
