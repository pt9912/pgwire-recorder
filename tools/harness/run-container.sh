#!/usr/bin/env bash
# run-container — Gate test-container (Vertrag: harness/vertraege/SPEC-055-container-und-binaries.md).
# Aufruf über `make test-container`, das VERSION, BINARY_PLATTFORMEN, IMAGE_PLATTFORMEN,
# POSTGRES_IMAGE und DOCKER aus harness/mk/container.mk als Umgebung übergibt. Die Punkte laufen
# in der Ordnung 1, 2 (Eingabe), 4 mit 3 (zwei Builds, gleiche Summen), 5 (Kopf), 6 (Image je
# Plattform), 7 (Smoke, tools/harness/container-smoke.sh); der erste rote Punkt beendet den Lauf. Die Marken `# P…` am Zeilenende
# benennen die Form für die Gegenprobe.
#
# Ausgang: ungleich 0 bei einem Fehler der Eingabe (vor dem ersten Aufruf von docker) und bei einem
# roten Punkt, sonst 0. SIGINT und SIGTERM beenden den laufenden Aufruf und räumen auf.
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
[ "${VERSION-}" != dev ] || eingabe "VERSION 'dev' ist am Gate unzulässig" # P2dev
if [ "$fehler" -ne 0 ]; then exit 2; fi # P8vor-docker

pg_image="${POSTGRES_IMAGE-}"
werk="$(mktemp -d -t pgr-container.XXXXXX)"
lauf="pgr-ct-$$"
netz="$lauf-net"
pg="$lauf-pg" rec="$lauf-rec" rep="$lauf-rep" img="$lauf-img"
laeuft=""

aufraeumen() {
  "$docker" rm -fv "$pg" "$rec" "$rep" "$img" "$lauf-version" "$lauf-psql" > /dev/null 2>&1 || true # P9container
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
  local n="pgwire-recorder-$1-$2"
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
    l "$docker" rm -fv "$img" > /dev/null
    soll="$(awk -v n="$(binaer "$os" "$arch")" '$2 == n { print $1 }' "$werk/a/SHA256SUMS")"
    g="$(sha256sum "$werk/image-$os-$arch" | cut -d' ' -f1)"
    [ "$g" = "$soll" ] || rot "Image $p trägt ein Binary mit der Summe $g statt $soll" # P6summe
  done
  echo "test-container: Punkt 6 grün — Image je Plattform mit dem Binary aus Punkt 4"
}

# Punkt 7: Smoke mit dem Image der Plattform des Docker-Servers; die Schritte a bis c stehen in container-smoke.sh.
punkt7() {
  local arch
  arch="$("$docker" version --format '{{.Server.Arch}}')"
  case "$arch" in amd64 | arm64) ;; *) rot "Plattform $arch des Docker-Servers ist weder amd64 noch arm64" ;; esac # P7server
  l env VERSION="$VERSION" POSTGRES_IMAGE="$pg_image" DOCKER="$docker" SMOKE_LAUF="$lauf" SMOKE_WERK="$werk" SMOKE_NAME="$name" \
    bash "$hier/container-smoke.sh" "$lauf:linux-$arch"
  echo "test-container: Punkt 7 grün — Smoke auf linux/$arch"
}

punkt4
punkt5
punkt6
punkt7
echo "test-container: grün"
