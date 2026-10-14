#!/usr/bin/env bash
# run-postgres-versionen — fährt die Teilmenge der Integrationstests gegen jede Version
# der Liste außer der Referenz (Vertrag: harness/vertraege/SPEC-053-postgres-versionen.md).
# Aufruf über `make test-postgres-versionen`, das die Liste und die Teilmenge aus
# harness/mk/postgres-versionen.mk als Umgebung übergibt: POSTGRES_VERSIONEN
# (`<Hauptversion>=<Image>`, durch Leerraum getrennt), POSTGRES_REFERENZ, POSTGRES_IMAGE,
# POSTGRES_TESTS_PHASE1, POSTGRES_TESTS_PHASE2 und DOCKER. Je Version läuft
# tools/test/run-integration-tests.sh mit Teilmenge, einem Benutzer mit scram-sha-256 und
# der Prüfung der Hauptversion; das Testimage wird einmal gebaut.
#
# Ausgang: ungleich 0 bei einem Fehler der Eingabe (vor dem ersten Aufruf von docker), bei einer
# roten Version und bei einem roten Aufbau des Testimages, sonst 0. SIGINT und SIGTERM gibt das
# Skript als SIGTERM an den laufenden Runner weiter und endet, wenn der Runner aufgeräumt hat.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
docker="${DOCKER:-docker}"
liste="${POSTGRES_VERSIONEN-}"
ref="${POSTGRES_REFERENZ-}"
fehler=0
eingabe() { echo "test-postgres-versionen: Eingabe: $*" >&2; fehler=1; }

# Punkt 2: neun Formen der Eingabe, je eine Zeile; die Marke am Zeilenende benennt die Form.
read -r -a eintraege <<< "$liste"
gesehen=" "
andere=0
refimage=""
[ "${#eintraege[@]}" -gt 0 ] || eingabe "die Liste der Versionen ist leer" # P2a
for e in "${eintraege[@]}"; do
  case "$e" in *=*) ;; *) eingabe "Eintrag ohne '=': $e"; continue ;; esac # P2e1
  haupt="${e%%=*}"
  image="${e#*=}"
  [ -n "$haupt" ] || eingabe "leere Hauptversion im Eintrag $e" # P2e2
  [[ -z "$haupt" || "$haupt" =~ ^[0-9]+$ ]] || eingabe "Hauptversion nicht aus Ziffern im Eintrag $e" # P2e3
  [[ "$haupt" != 0* ]] || eingabe "Hauptversion mit führender Null im Eintrag $e" # P2e4
  [[ "$gesehen" != *" $haupt "* || "$haupt" == "$ref" ]] || eingabe "Hauptversion doppelt im Eintrag $e" # P2c
  [[ "$gesehen" != *" $haupt "* || "$haupt" != "$ref" ]] || eingabe "Hauptversion der Referenz doppelt im Eintrag $e" # P2c-ref
  gesehen+="$haupt "
  name="${image%%@*}"
  name="${name##*/}"
  [[ "$name" =~ ^[^:]+:[^:]+$ ]] || eingabe "Image ohne Tag im Eintrag $e" # P2f
  [[ "$image" == *@* ]] || eingabe "Image ohne Digest im Eintrag $e" # P2g
  [[ "$image" != *@* || "${image#*@}" =~ ^sha256:[0-9a-f]{64}$ ]] || eingabe "Digest ungültig im Eintrag $e" # P2h
  if [ "$haupt" = "$ref" ]; then refimage="$image"; else andere=$((andere + 1)); fi
done
[[ "$gesehen" == *" $ref "* ]] || eingabe "die Referenz '$ref' fehlt in der Liste" # P2b
[ "$andere" -gt 0 ] || eingabe "die Liste nennt keine andere Version als die Referenz" # P2d
[[ -z "$refimage" || "$refimage" == "${POSTGRES_IMAGE-}" ]] || eingabe "Image der Referenz '$refimage' ungleich POSTGRES_IMAGE '${POSTGRES_IMAGE-}'" # P2i
if [ "$fehler" -ne 0 ]; then exit 2; fi

werk="$(mktemp -d)"
lauf="pgr-pv-$$"
laeuft=""
beenden() {
  trap '' INT TERM
  if [ -n "$laeuft" ]; then
    kill -TERM "$laeuft" 2> /dev/null || true
    wait "$laeuft" 2> /dev/null || true
  fi
  exit "$1"
}
trap 'rm -rf "$werk"' EXIT
trap 'beenden 130' INT # S-INT
trap 'beenden 143' TERM # S-TERM

"$docker" build --progress=plain --target integration -t pgwire-recorder:integration . &
laeuft=$!
if ! wait "$laeuft"; then
  echo "test-postgres-versionen: ROT — das Testimage ist nicht zu bauen" >&2
  exit 1
fi

rot=()
for e in "${eintraege[@]}"; do
  haupt="${e%%=*}"
  [ "$haupt" != "$ref" ] || continue # P6-referenz
  PGR_LAUF="$lauf-$haupt" PGR_BAUEN=0 PGR_ROLLEN=scram PGR_HAUPTVERSION="$haupt" \
    PGR_PHASE_DATEI="$werk/phase-$haupt" PGR_TESTS_1="${POSTGRES_TESTS_PHASE1-}" PGR_TESTS_2="${POSTGRES_TESTS_PHASE2-}" \
    POSTGRES_IMAGE="${e#*=}" DOCKER="$docker" bash "$hier/run-integration-tests.sh" &
  laeuft=$!
  status=0
  wait "$laeuft" || status=$?
  laeuft=""
  if [ "$status" -eq 0 ]; then
    echo "test-postgres-versionen: Version $haupt grün"
  else
    case "$(cat "$werk/phase-$haupt" 2> /dev/null || true)" in
      1) ort="Phase 1" ;;
      2) ort="Phase 2" ;;
      *) ort="Aufbau" ;;
    esac
    rot+=("$haupt ($ort)") # P6-phase
  fi
done
laeuft=""
if [ "${#rot[@]}" -gt 0 ]; then
  echo "test-postgres-versionen: ROT — ${rot[*]}" >&2
  exit 1 # P6-exit
fi
echo "test-postgres-versionen: grün"
