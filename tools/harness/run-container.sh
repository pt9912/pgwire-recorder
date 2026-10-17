#!/usr/bin/env bash
# run-container — Gate test-container (Vertrag: harness/vertraege/SPEC-055-container-und-binaries.md).
# Aufruf über `make test-container`, das VERSION, BINARY_PLATTFORMEN, IMAGE_PLATTFORMEN,
# POSTGRES_IMAGE und DOCKER aus harness/mk/container.mk als Umgebung übergibt. Die Punkte laufen
# in der Ordnung 1, 2 (Eingabe), 4 mit 3 (zwei Builds, gleiche Summen), 5 (Kopf), 6 (Image je
# Plattform), 7 (Smoke); der erste rote Punkt beendet den Lauf. Die Marken `# P…` am Zeilenende
# benennen die Form für die Gegenprobe.
#
# Ausgang: 2 bei einem Fehler der Eingabe (vor dem ersten Aufruf von docker), 1 bei einem roten
# Punkt, sonst 0. SIGINT und SIGTERM beenden den laufenden Aufruf und räumen auf.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
docker="${DOCKER:-docker}"
name=test-container
fehler=0
# shellcheck source=tools/harness/container-eingabe.sh
. "$hier/container-eingabe.sh"

read -r -a bin_pl <<< "${BINARY_PLATTFORMEN-}"
read -r -a img_pl <<< "${IMAGE_PLATTFORMEN-}"
liste_pruefen "der Binaries" "$ERF_BINARIES" "${bin_pl[@]}"
liste_pruefen "der Images" "$ERF_IMAGES" "${img_pl[@]}"
images_in_binaries "${BINARY_PLATTFORMEN-}" "${img_pl[@]}"
version_pruefen
if [ "$fehler" -ne 0 ]; then exit 2; fi # P8vor-docker

pg_image="${POSTGRES_IMAGE-}"
werk="$(mktemp -d -t pgr-container.XXXXXX)"
lauf="pgr-ct-$$"
netz="$lauf-net"
pg="$lauf-pg" rec="$lauf-rec" rep="$lauf-rep" img="$lauf-img"
laeuft=""

aufraeumen() {
  "$docker" rm -f "$pg" "$rec" "$rep" "$img" "$lauf-version" "$lauf-psql" > /dev/null 2>&1 || true # P9container
  "$docker" network rm "$netz" > /dev/null 2>&1 || true # P9netz
  for p in "${img_pl[@]}"; do "$docker" rmi "$lauf:${p%/*}-${p#*/}" > /dev/null 2>&1 || true; done # P9tags
  rm -rf "$werk" # P9werk
}
beenden() {
  trap '' INT TERM
  if [ -n "$laeuft" ]; then
    kill -TERM "$laeuft" 2> /dev/null || true
    wait "$laeuft" 2> /dev/null || true
  fi
  exit "$1"
}
trap aufraeumen EXIT # P9exit
trap 'beenden 130' INT # P9int
trap 'beenden 143' TERM # P9term

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
  echo "test-container: ROT — $*" >&2
  exit 1
}
hex() { od -An -tx1 -v -j "$2" -N "$3" "$1" | tr -d ' \n'; }

binaer() { # <os> <arch> — Name des Binaries
  local n="pgwire-recorder_${VERSION}_$1_$2"
  if [ "$1" = windows ]; then n+=".exe"; fi
  echo "$n"
}

# Punkt 4 (mit 3): zwei Builds in zwei Verzeichnissen, der zweite ohne Cache der Stufen build und dist.
punkt4() {
  local n oc abw
  for n in a b; do
    oc=""
    if [ "$n" = b ]; then oc=1; fi # P4zweiter
    l env DIST="$werk/$n" OHNE_CACHE="$oc" bash "$hier/binaries.sh" > "$werk/bin-$n.log" 2>&1 \
      || { cat "$werk/bin-$n.log" >&2; rot "Bau $n der Binaries gescheitert"; }
  done
  if ! cmp -s "$werk/a/SHA256SUMS" "$werk/b/SHA256SUMS"; then # P4cmp
    abw="$(diff <(cat "$werk/a/SHA256SUMS") <(cat "$werk/b/SHA256SUMS") | sed -n 's/^< [0-9a-f]*  //p' | tr '\n' ' ')" || true # P4namen
    rot "die Summen zweier Builds weichen ab: ${abw% }"
  fi
  echo "test-container: Punkt 4 grün — zwei Builds, dieselben Summen"
}

# Punkt 5: Kopf jedes Binaries.
kopf() { # <os> <arch> <datei> — Exit 0, sonst der gelesene Wert auf stdout
  local os="$1" arch="$2" f="$3" erw g o
  case "$os" in
    linux)
      g="$(hex "$f" 0 4)"
      [ "$g" = 7f454c46 ] || { echo "ELF-Magic $g"; return 1; } # P5elf-magic
      case "$arch" in amd64) erw=3e00 ;; arm64) erw=b700 ;; esac # P5elf-arch
      g="$(hex "$f" 18 2)"
      [ "$g" = "$erw" ] || { echo "ELF-Maschine $g"; return 1; } # P5elf-maschine
      ;;
    darwin)
      g="$(hex "$f" 0 4)"
      [ "$g" = cffaedfe ] || { echo "Mach-O-Magic $g"; return 1; } # P5macho-magic
      case "$arch" in amd64) erw=07000001 ;; arm64) erw=0c000001 ;; esac # P5macho-arch
      g="$(hex "$f" 4 4)"
      [ "$g" = "$erw" ] || { echo "Mach-O-CPU-Typ $g"; return 1; } # P5macho-cpu
      ;;
    windows)
      g="$(hex "$f" 0 2)"
      [ "$g" = 4d5a ] || { echo "PE-Magic $g"; return 1; } # P5pe-magic
      o="$(hex "$f" 60 4)"
      o=$((16#${o:6:2}${o:4:2}${o:2:2}${o:0:2}))
      case "$arch" in amd64) erw=6486 ;; arm64) erw=64aa ;; esac # P5pe-arch
      g="$(hex "$f" $((o + 4)) 2)"
      [ "$g" = "$erw" ] || { echo "PE-Maschine $g"; return 1; } # P5pe-maschine
      ;;
  esac
}
punkt5() {
  local p g
  for p in "${bin_pl[@]}"; do
    g="$(kopf "${p%/*}" "${p#*/}" "$werk/a/$(binaer "${p%/*}" "${p#*/}")")" \
      || rot "Kopf des Binaries für $p: $g"
  done
  echo "test-container: Punkt 5 grün — Kopf von ${#bin_pl[@]} Binaries"
}

# Punkt 6: Image je Plattform, mit dem Binary aus Punkt 4.
punkt6() {
  local p os arch tag g soll
  for p in "${img_pl[@]}"; do
    os="${p%/*}"
    arch="${p#*/}"
    tag="$lauf:$os-$arch"
    l "$docker" build --progress=plain --target runtime --platform "$p" --build-arg "VERSION=$VERSION" -t "$tag" . \
      > "$werk/image-$os-$arch.log" 2>&1 \
      || { tail -n 20 "$werk/image-$os-$arch.log" >&2; rot "Image für $p nicht zu bauen"; } # P6bau
    g="$("$docker" image inspect --format '{{.Os}}' "$tag")"
    [ "$g" = "$os" ] || rot "Image $p nennt das Betriebssystem '$g'" # P6os
    g="$("$docker" image inspect --format '{{.Architecture}}' "$tag")"
    [ "$g" = "$arch" ] || rot "Image $p nennt die Architektur '$g'" # P6arch
    l "$docker" create --name "$img" --platform "$p" "$tag" > /dev/null
    l "$docker" cp "$img:/pgwire-recorder" "$werk/image-$os-$arch"
    l "$docker" rm -f "$img" > /dev/null
    soll="$(awk -v n="$(binaer "$os" "$arch")" '$2 == n { print $1 }' "$werk/a/SHA256SUMS")"
    g="$(sha256sum "$werk/image-$os-$arch" | cut -d' ' -f1)"
    [ "$g" = "$soll" ] || rot "Image $p trägt ein Binary mit der Summe $g statt $soll" # P6summe
  done
  echo "test-container: Punkt 6 grün — Image je Plattform mit dem Binary aus Punkt 4"
}

# Punkt 7: Smoke mit dem Image der Plattform des Docker-Servers.
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
punkt7() {
  local arch tag i
  arch="$("$docker" version --format '{{.Server.Arch}}')"
  case "$arch" in amd64 | arm64) ;; *) rot "Plattform $arch des Docker-Servers ist weder amd64 noch arm64" ;; esac # P7server
  tag="$lauf:linux-$arch"
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
  l "$docker" rm -f "$pg" > /dev/null
  if "$docker" container inspect "$pg" > /dev/null 2>&1; then rot "der PostgreSQL-Container lebt beim Start des Wiedergebens noch"; fi # P7pg
  l "$docker" create --name "$rep" --network "$netz" "$tag" replay --listen 0.0.0.0:5432 --input /tmp/aufzeichnung.yaml > /dev/null
  l "$docker" cp "$werk/aufzeichnung.yaml" "$rep:/tmp/aufzeichnung.yaml"
  l "$docker" start "$rep" > /dev/null
  bereit "$rep"
  client "$rep" "$werk/sicht-b"
  cmp -s "$werk/sicht-a" "$werk/sicht-b" || rot "Sicht B ist ungleich Sicht A" # P7gleich
  anhalten "$rep" replay
  echo "test-container: Punkt 7 grün — Smoke auf linux/$arch"
}

punkt4
punkt5
punkt6
punkt7
echo "test-container: grün"
