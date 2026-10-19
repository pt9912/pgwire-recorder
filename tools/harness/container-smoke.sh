#!/usr/bin/env bash
# container-smoke — Smoke im Container (Vertrag: harness/vertraege/SPEC-055-container-und-binaries.md Punkt 7): version, Aufzeichnen, Wiedergeben mit einem Image.
# Aufruf: bash tools/harness/container-smoke.sh <Image-Referenz>; gerufen von tools/harness/run-container.sh (Image des Laufs) und von tools/harness/release.sh
# (Image aus einer Registry, harness/vertraege/SPEC-067-release.md Punkt 7).
# Umgebung, alle verlangt: VERSION (erwartete Ausgabe von `version`), POSTGRES_IMAGE, SMOKE_LAUF (Namensbestandteil der Container und des Netzes), SMOKE_WERK (vorhandenes
# Verzeichnis für die Dateien des Smokes), SMOKE_NAME (Präfix der Meldungen); DOCKER (Vorgabe docker). Das Skript räumt nichts auf: Container, Netz und Verzeichnis sind Sache des Aufrufers.
# Ausgang: 0 grün; ungleich 0 bei einem roten Schritt, die Meldung `<SMOKE_NAME>: ROT — …` nennt Schritt und Wert. SIGINT und SIGTERM beenden den laufenden Aufruf.
set -euo pipefail

docker="${DOCKER:-docker}"
for v in VERSION POSTGRES_IMAGE SMOKE_LAUF SMOKE_WERK SMOKE_NAME; do [ -n "${!v-}" ] || { echo "container-smoke: $v fehlt" >&2; exit 2; }; done
[ "$#" -eq 1 ] && [ -n "$1" ] || { echo "container-smoke: genau eine Image-Referenz verlangt" >&2; exit 2; }
tag="$1"
pg_image="$POSTGRES_IMAGE"
lauf="$SMOKE_LAUF"
werk="$SMOKE_WERK"
netz="$lauf-net"
pg="$lauf-pg" rec="$lauf-rec" rep="$lauf-rep"
laeuft=""

beenden() {
  trap '' INT TERM
  if [ -n "$laeuft" ]; then
    kill -TERM "$laeuft" 2> /dev/null || true
    wait "$laeuft" 2> /dev/null || true
  fi
  exit "$1"
}
trap 'beenden 130' INT
trap 'beenden 143' TERM

# l <Befehl>… — der Befehl läuft im Hintergrund, damit ein Signal das wait beendet.
l() {
  local s=0
  "$@" &
  laeuft=$!
  wait "$laeuft" || s=$?
  laeuft=""
  return "$s"
}
rot() {
  echo "$SMOKE_NAME: ROT — $*" >&2
  exit 1
}

bereit() { # <Container> — wartet auf die Startmeldung im Log
  local i
  for i in $(seq 1 120); do
    if "$docker" logs "$1" 2>&1 | grep -q 'gestartet'; then return 0; fi
    sleep 0.5
  done
  rot "$1 meldet keinen Start"
}
anhalten() { # <Container> <Name> — SIGTERM, Exit-Code 0
  local c
  l "$docker" stop "$1" > /dev/null
  c="$("$docker" inspect --format '{{.State.ExitCode}}' "$1")"
  [ "$c" = 0 ] || rot "Exit-Code $c des Containers mit $2 (erwartet 0)" # P7exit
}
client() { # <Ziel> <Datei> — eine Anfrage über psql aus dem PostgreSQL-Image
  l "$docker" run --rm --name "$lauf-psql" --network "$netz" "$pg_image" \
    psql -h "$1" -U postgres -c 'SELECT 1 AS eins' > "$2" || rot "Anfrage an $1 scheitert"
}

l "$docker" network create --internal "$netz" > /dev/null # P7intern
# a. version
l "$docker" run --rm --name "$lauf-version" "$tag" version > "$werk/version" || rot "version endet mit einem Exit-Code ungleich 0" # P7a-exit
printf 'pgwire-recorder %s\n' "$VERSION" > "$werk/version-soll"
cmp -s "$werk/version" "$werk/version-soll" || rot "version gibt '$(cat "$werk/version")' aus statt 'pgwire-recorder $VERSION' und LF" # P7a-ausgabe
# b. Aufzeichnen
l "$docker" run -d --name "$pg" --network "$netz" -e POSTGRES_HOST_AUTH_METHOD=trust "$pg_image" > /dev/null
for i in $(seq 1 120); do
  if l "$docker" exec "$pg" pg_isready -q -h 127.0.0.1 -U postgres; then break; fi
  sleep 0.5
  [ "$i" -lt 120 ] || rot "PostgreSQL nicht bereit"
done
l "$docker" run -d --name "$rec" --network "$netz" "$tag" record --listen 0.0.0.0:5432 --upstream "$pg:5432" --output /tmp/aufzeichnung.yaml > /dev/null
bereit "$rec"
client "$rec" "$werk/sicht-a"
[ -s "$werk/sicht-a" ] || rot "Sicht A ist leer" # P7a-leer
grep -q eins "$werk/sicht-a" || rot "Sicht A nennt 'eins' nicht" # P7a-eins
anhalten "$rec" record
l "$docker" cp "$rec:/tmp/aufzeichnung.yaml" "$werk/aufzeichnung.yaml" || rot "die Aufzeichnung ist dem Container nicht zu entnehmen" # P7cp
[ -s "$werk/aufzeichnung.yaml" ] || rot "die Aufzeichnung ist leer" # P7aufz
# c. Wiedergeben ohne PostgreSQL
l "$docker" rm -fv "$pg" > /dev/null # P9pg-replay
if "$docker" container inspect "$pg" > /dev/null 2>&1; then rot "der PostgreSQL-Container lebt beim Start des Wiedergebens noch"; fi # P7pg
l "$docker" create --name "$rep" --network "$netz" "$tag" replay --listen 0.0.0.0:5432 --input /tmp/aufzeichnung.yaml > /dev/null
l "$docker" cp "$werk/aufzeichnung.yaml" "$rep:/tmp/aufzeichnung.yaml"
l "$docker" start "$rep" > /dev/null
bereit "$rep"
client "$rep" "$werk/sicht-b"
cmp -s "$werk/sicht-a" "$werk/sicht-b" || rot "Sicht B ist ungleich Sicht A" # P7gleich
anhalten "$rep" replay
