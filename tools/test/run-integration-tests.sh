#!/usr/bin/env bash
# run-integration-tests — startet eine PostgreSQL-Instanz und das Testimage
# (Dockerfile, Stufe integration) in einem eigenen, nach außen abgeschlossenen
# Docker-Netz, führt die Integrationstests aus und räumt Container und Netz
# danach auf, auch nach einem Abbruch. Beide Phasen teilen ein Docker-Volume
# (PGR_DATEN): Die erste schreibt dort Aufzeichnungen über record, danach
# stoppt der Runner PostgreSQL, und die zweite führt TestE2EOhnePostgres* gegen
# diese Aufzeichnungen aus (PGR_OHNE_POSTGRES=1).
# Auf der Instanz legt der Runner drei Benutzer an, deren Anmeldung pg_hba.conf
# auf scram-sha-256, md5 und password festlegt; ihre Passwörter reichen
# PGR_PW_SCRAM, PGR_PW_MD5 und PGR_PW_KLARTEXT an den Testlauf.
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
daten="$lauf-daten"

# Passwörter der Benutzer mit festem Anmeldeverfahren; sie tragen GEHEIM, damit
# die Tests prüfen können, dass keine Ausgabe sie nennt.
pw_scram='GEHEIM scram $%41'
pw_md5='GEHEIMmd5'
pw_klartext='GEHEIMklartext'

aufraeumen() {
  "$docker" rm -f "$tests" "$pg" >/dev/null 2>&1 || true
  "$docker" network rm "$netz" >/dev/null 2>&1 || true
  "$docker" volume rm "$daten" >/dev/null 2>&1 || true
}
trap aufraeumen EXIT

"$docker" build --progress=plain --target integration -t pgwire-recorder:integration .
"$docker" network create --internal "$netz" >/dev/null
"$docker" volume create "$daten" >/dev/null
"$docker" run -d --name "$pg" --network "$netz" \
  -e POSTGRES_HOST_AUTH_METHOD=trust "$postgres_image" >/dev/null

# Erst der Server auf TCP ist die Instanz; der Server der Initialisierung hört
# nur auf dem Unix-Socket.
bereit=0
for _ in $(seq 1 120); do
  if "$docker" exec "$pg" pg_isready -q -h 127.0.0.1 -U postgres; then bereit=1; break; fi
  sleep 0.5
done
if [ "$bereit" -ne 1 ]; then
  echo "run-integration-tests: PostgreSQL nicht bereit" >&2
  exit 1
fi
# pg_hba.conf: Die Zeilen der drei Benutzer stehen vor der Zeile für alle.
"$docker" exec -i -u postgres "$pg" sh -c '
  set -e
  hba="$(psql -U postgres -At -c "SHOW hba_file")"
  { printf "%s\n" "host all play_scram all scram-sha-256" "host all play_md5 all md5" "host all play_pw all password"; cat "$hba"; } > "$hba.neu"
  cat "$hba.neu" > "$hba"
  rm "$hba.neu"
  psql -U postgres -v ON_ERROR_STOP=1 -q' >/dev/null <<SQL
SET password_encryption = 'scram-sha-256';
CREATE ROLE play_scram LOGIN PASSWORD '$pw_scram';
SET password_encryption = 'md5';
CREATE ROLE play_md5 LOGIN PASSWORD '$pw_md5';
CREATE ROLE play_pw LOGIN PASSWORD '$pw_klartext';
SELECT pg_reload_conf();
SQL

status=0
"$docker" run --rm --name "$tests" --network "$netz" -v "$daten":/daten -e PGR_DATEN=/daten \
  -e PGR_UPSTREAM="$pg:5432" -e PGR_PW_SCRAM="$pw_scram" -e PGR_PW_MD5="$pw_md5" \
  -e PGR_PW_KLARTEXT="$pw_klartext" pgwire-recorder:integration || status=$?

if [ "$status" -ne 0 ]; then
  echo "run-integration-tests: Testlauf rot (Exit-Code $status)" >&2
  "$docker" logs "$pg" 2>&1 | tail -n 20 >&2 || true
  exit "$status"
fi

"$docker" stop "$pg" >/dev/null
"$docker" run --rm --name "$tests" --network "$netz" -v "$daten":/daten -e PGR_DATEN=/daten \
  -e PGR_UPSTREAM="$pg:5432" -e PGR_OHNE_POSTGRES=1 \
  pgwire-recorder:integration -test.run '^TestE2EOhnePostgres' || status=$?
if [ "$status" -ne 0 ]; then
  echo "run-integration-tests: Phase ohne PostgreSQL rot (Exit-Code $status)" >&2
  exit "$status"
fi
echo "run-integration-tests: gruen"
