#!/usr/bin/env bash
# run-integration-tests — startet eine PostgreSQL-Instanz und das Testimage
# (Dockerfile, Stufe integration) in einem eigenen, nach außen abgeschlossenen
# Docker-Netz, führt die Integrationstests aus und räumt Container und Netz
# danach auf, auch nach einem Abbruch. In einer zweiten Phase stoppt es
# PostgreSQL und führt die Tests TestE2EOhnePostgres* aus (PGR_OHNE_POSTGRES=1).
# Es schreibt nichts in den Arbeitsbaum;
# die Abdeckungstabellen erzeugt `make abdeckung`.
#
# Ausgang: der Exit-Code des Testlaufs; ungleich 0, wenn der Aufbau scheitert.
set -euo pipefail

docker="${DOCKER:-docker}"
postgres_image="${POSTGRES_IMAGE:?POSTGRES_IMAGE fehlt}"
lauf="pgr-it-$$"
netz="$lauf-net"
pg="$lauf-pg"
tests="$lauf-tests"

aufraeumen() {
  "$docker" rm -f "$tests" "$pg" >/dev/null 2>&1 || true
  "$docker" network rm "$netz" >/dev/null 2>&1 || true
}
trap aufraeumen EXIT

"$docker" build --progress=plain --target integration -t pgwire-recorder:integration .
"$docker" network create --internal "$netz" >/dev/null
"$docker" run -d --name "$pg" --network "$netz" \
  -e POSTGRES_HOST_AUTH_METHOD=trust "$postgres_image" >/dev/null

status=0
"$docker" run --rm --name "$tests" --network "$netz" -e PGR_UPSTREAM="$pg:5432" \
  pgwire-recorder:integration || status=$?

if [ "$status" -ne 0 ]; then
  echo "run-integration-tests: Testlauf rot (Exit-Code $status)" >&2
  "$docker" logs "$pg" 2>&1 | tail -n 20 >&2 || true
  exit "$status"
fi

"$docker" stop "$pg" >/dev/null
"$docker" run --rm --name "$tests" --network "$netz" -e PGR_UPSTREAM="$pg:5432" -e PGR_OHNE_POSTGRES=1 \
  pgwire-recorder:integration -test.run '^TestE2EOhnePostgres' || status=$?
if [ "$status" -ne 0 ]; then
  echo "run-integration-tests: Phase ohne PostgreSQL rot (Exit-Code $status)" >&2
  exit "$status"
fi
echo "run-integration-tests: gruen"
