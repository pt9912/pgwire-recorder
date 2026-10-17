#!/usr/bin/env bash
# binaries — baut das Binary jeder Zielplattform und schreibt es mit SHA256SUMS nach DIST
# (Vertrag: harness/vertraege/SPEC-055-container-und-binaries.md Punkt 3). Aufruf über
# `make binaries` und aus tools/harness/run-container.sh, mit der Umgebung VERSION, DIST,
# BINARY_PLATTFORMEN (durch Leerraum getrennt) und DOCKER; OHNE_CACHE=1 umgeht den Cache
# der Stufen `build` und `dist`.
#
# Ausgang: ungleich 0 bei einem Fehler der Eingabe (vor dem ersten Aufruf von docker, auch bei einem
# DIST, das nicht leer ist) und bei einem gescheiterten Build (DIST bleibt leer), sonst 0.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
docker="${DOCKER:-docker}"
name=binaries
fehler=0
# shellcheck source=tools/harness/container-eingabe.sh
. "$hier/container-eingabe.sh"

read -r -a plattformen <<< "${BINARY_PLATTFORMEN-}"
liste_pruefen "der Binaries" "$ERF_BINARIES" "${plattformen[@]}"
version_pruefen
dist="${DIST-}"
[ -n "$dist" ] || eingabe "DIST ist leer" # P3distleer
[ -z "$(ls -A "$dist" 2> /dev/null)" ] || eingabe "DIST '$dist' ist nicht leer" # P3dist
if [ "$fehler" -ne 0 ]; then exit 2; fi
mkdir -p "$dist"

stage="$(mktemp -d)"
laeuft=""
beenden() {
  trap '' INT TERM
  if [ -n "$laeuft" ]; then
    kill -TERM "$laeuft" 2> /dev/null || true
    wait "$laeuft" 2> /dev/null || true
  fi
  exit "$1"
}
trap 'rm -rf "$stage"' EXIT
trap 'beenden 130' INT
trap 'beenden 143' TERM

cache=()
if [ "${OHNE_CACHE-}" = 1 ]; then cache=(--no-cache-filter build --no-cache-filter dist); fi # P4cache

namen=()
for p in "${plattformen[@]}"; do
  os="${p%/*}"
  arch="${p#*/}"
  datei="pgwire-recorder_${VERSION}_${os}_${arch}" # P3name
  if [ "$os" = windows ]; then datei+=".exe"; fi # P3exe
  "$docker" build --progress=plain --target dist --platform "$p" --build-arg "VERSION=$VERSION" \
    ${cache[@]+"${cache[@]}"} --output "type=local,dest=$stage/$os-$arch" . > "$stage/log-$os-$arch" 2>&1 &
  laeuft=$!
  status=0
  wait "$laeuft" || status=$?
  laeuft=""
  if [ "$status" -ne 0 ]; then
    tail -n 20 "$stage/log-$os-$arch" >&2 || true
    echo "binaries: ROT — Build für $p gescheitert" >&2
    exit 1 # P3rot
  fi
  mv "$stage/$os-$arch/pgwire-recorder" "$stage/$datei"
  namen+=("$datei")
done

mapfile -t namen < <(printf '%s\n' "${namen[@]}" | LC_ALL=C sort) # P3sort
mv "${namen[@]/#/$stage/}" "$dist/"
(cd "$dist" && sha256sum "${namen[@]}" > SHA256SUMS)
echo "binaries: ${#namen[@]} Binaries und SHA256SUMS in $dist"
