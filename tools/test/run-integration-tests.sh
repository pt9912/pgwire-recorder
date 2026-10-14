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
# PGR_PW_SCRAM, PGR_PW_MD5 und PGR_PW_KLARTEXT an den Testlauf. Für play_reject
# lehnt pg_hba.conf jede Verbindung ab.
# Es schreibt nichts in den Arbeitsbaum;
# die Abdeckungstabellen erzeugt `make abdeckung`.
#
# Ohne die folgenden Variablen läuft die volle Suite (SPEC-053 Punkt 8). Das Ziel
# test-postgres-versionen (tools/test/run-postgres-versionen.sh) setzt sie:
#   PGR_LAUF          Kennung für die Namen von Container, Netz und Volume
#   PGR_BAUEN=0       das Testimage nicht bauen (es liegt vor)
#   PGR_ROLLEN=scram  nur den Benutzer mit scram-sha-256 anlegen
#   PGR_TESTS_1/_2    Namen der Tests der Phase 1 und 2, durch Leerzeichen getrennt;
#                     nur sie laufen, und jeder braucht eine PASS-Zeile
#   PGR_HAUPTVERSION  Hauptversion, die server_version des Servers tragen muss
#   PGR_PHASE_DATEI   Datei, in die der Runner aufbau, 1, 2 oder fertig schreibt
#
# Ausgang: der Exit-Code des Testlaufs; ungleich 0, wenn der Aufbau scheitert.
set -euo pipefail

docker="${DOCKER:-docker}"
postgres_image="${POSTGRES_IMAGE:?POSTGRES_IMAGE fehlt}"
lauf="${PGR_LAUF:-pgr-it-$$}"
netz="$lauf-net"
pg="$lauf-pg"
tests="$lauf-tests"
daten="$lauf-daten"
version="${PGR_HAUPTVERSION:-Referenz}"
werk="$(mktemp -d)"

# Passwörter der Benutzer mit festem Anmeldeverfahren; sie tragen GEHEIM, damit
# die Tests prüfen können, dass keine Ausgabe sie nennt.
pw_scram='GEHEIM scram $%41'
pw_md5='GEHEIMmd5'
pw_klartext='GEHEIMklartext'

aufraeumen() {
  "$docker" rm -f "$tests" "$pg" >/dev/null 2>&1 || true
  "$docker" network rm "$netz" >/dev/null 2>&1 || true
  "$docker" volume rm "$daten" >/dev/null 2>&1 || true
  rm -rf "$werk"
}
trap aufraeumen EXIT # S-EXIT
trap 'trap "" TERM; exit 143' TERM # S-TERM

phase() { if [ -n "${PGR_PHASE_DATEI:-}" ]; then printf '%s\n' "$1" > "$PGR_PHASE_DATEI"; fi; }

# Lange Aufrufe laufen im Hintergrund, damit ein Signal das wait beendet und die
# Falle EXIT läuft. Ein Hintergrundjob einer nicht interaktiven bash ignoriert SIGINT:
# Als Kind von run-postgres-versionen.sh endet dieser Runner auf SIGTERM, das jenes
# Skript weiterreicht.
vordergrund() { "$@" & wait $!; }

# auswahl <name>… — das Muster für -test.run; ein Name `Test/Sub` wählt Sub.
auswahl() {
  local n top=() sub=() muster
  for n; do
    top+=("${n%%/*}")
    if [ "$n" != "${n%%/*}" ]; then sub+=("${n#*/}"); fi
  done
  muster="^($(IFS='|'; echo "${top[*]}"))\$"
  if [ "${#sub[@]}" -gt 0 ]; then muster="$muster/^($(IFS='|'; echo "${sub[*]}"))\$"; fi
  echo "$muster"
}

# testlauf <ausgabe> <argument>… — docker run des Testimages; die Ausgabe geht an
# stdout und in die Datei.
testlauf() {
  local aus="$1"
  shift
  { "$docker" run --rm --name "$tests" --network "$netz" "$@" 2>&1 | tee "$aus"; } &
  wait $!
}

# erwartet <ausgabe> <phase> <name>… — jeder Name braucht eine PASS-Zeile.
erwartet() {
  local aus="$1" ph="$2" name rot=0
  shift 2
  if ! grep -Eq '^ *--- PASS: ' "$aus"; then # P3-treffer
    echo "run-integration-tests: Version $version, Phase $ph: Auswahl ohne Treffer, kein Test bestanden" >&2
    return 1
  fi
  for name; do
    if grep -Eq "^ *--- PASS: $name \(" "$aus"; then continue; fi
    if grep -Eq "^ *--- SKIP: $name \(" "$aus"; then # P3-skip
      echo "run-integration-tests: Version $version, Phase $ph: Test $name übersprungen" >&2
    else
      echo "run-integration-tests: Version $version, Phase $ph: Test $name ohne PASS-Zeile (umbenannt oder nicht gewählt)" >&2
    fi
    rot=1 # P3-fehlt
  done
  return "$rot"
}

phase aufbau
if [ "${PGR_BAUEN:-1}" != 0 ]; then
  vordergrund "$docker" build --progress=plain --target integration -t pgwire-recorder:integration .
fi
"$docker" network create --internal "$netz" >/dev/null
"$docker" volume create "$daten" >/dev/null
vordergrund "$docker" run -d --name "$pg" --network "$netz" \
  -e POSTGRES_HOST_AUTH_METHOD=trust "$postgres_image" >/dev/null

# Erst der Server auf TCP ist die Instanz; der Server der Initialisierung hört
# nur auf dem Unix-Socket.
bereit=0
for _ in $(seq 1 120); do
  if "$docker" exec "$pg" pg_isready -q -h 127.0.0.1 -U postgres; then bereit=1; break; fi
  sleep 0.5
done
if [ "$bereit" -ne 1 ]; then
  echo "run-integration-tests: Version $version: PostgreSQL nicht bereit" >&2
  exit 1 # E-bereit
fi
if [ -n "${PGR_HAUPTVERSION:-}" ]; then
  gemeldet="$("$docker" exec "$pg" psql -U postgres -At -c 'SHOW server_version' | tr -d '\r')" || gemeldet=""
  if [ "${gemeldet%%[!0-9]*}" != "$PGR_HAUPTVERSION" ]; then # P4-version
    echo "run-integration-tests: Version $PGR_HAUPTVERSION: der Server meldet server_version '$gemeldet'" >&2
    exit 1
  fi
fi
# pg_hba.conf: Die Zeilen der Benutzer stehen vor der Zeile für alle.
hba='"host all play_scram all scram-sha-256"'
sql="SET password_encryption = 'scram-sha-256';
CREATE ROLE play_scram LOGIN PASSWORD '$pw_scram';"
if [ "${PGR_ROLLEN:-alle}" != scram ]; then # P5-rollen
  hba="$hba "'"host all play_md5 all md5" "host all play_pw all password" "host all play_reject all reject"'
  sql="$sql
SET password_encryption = 'md5';
CREATE ROLE play_md5 LOGIN PASSWORD '$pw_md5';
CREATE ROLE play_pw LOGIN PASSWORD '$pw_klartext';"
fi
skript='
  set -e
  hba="$(psql -U postgres -At -c "SHOW hba_file")"
  { printf "%s\n" '"$hba"'; cat "$hba"; } > "$hba.neu"
  cat "$hba.neu" > "$hba"
  rm "$hba.neu"
  psql -U postgres -v ON_ERROR_STOP=1 -q'
"$docker" exec -i -u postgres "$pg" sh -c "$skript" >/dev/null <<< "$sql
SELECT pg_reload_conf();"

read -r -a tests_1 <<< "${PGR_TESTS_1:-}"
read -r -a tests_2 <<< "${PGR_TESTS_2:-}"
wahl_1=()
wahl_2=(-test.run '^TestE2EOhnePostgres')
if [ "${#tests_1[@]}" -gt 0 ]; then wahl_1=(-test.run "$(auswahl "${tests_1[@]}")"); fi
if [ "${#tests_2[@]}" -gt 0 ]; then wahl_2=(-test.run "$(auswahl "${tests_2[@]}")"); fi

status=0
phase 1
testlauf "$werk/phase-1" -v "$daten":/daten -e PGR_DATEN=/daten \
  -e PGR_UPSTREAM="$pg:5432" -e PGR_PW_SCRAM="$pw_scram" -e PGR_PW_MD5="$pw_md5" \
  -e PGR_PW_KLARTEXT="$pw_klartext" pgwire-recorder:integration "${wahl_1[@]}" || status=$?

if [ "$status" -ne 0 ]; then
  echo "run-integration-tests: Testlauf rot (Exit-Code $status)" >&2
  "$docker" logs "$pg" 2>&1 | tail -n 20 >&2 || true
  exit "$status"
fi
if [ "${#tests_1[@]}" -gt 0 ]; then erwartet "$werk/phase-1" 1 "${tests_1[@]}" || exit 1; fi

"$docker" stop "$pg" >/dev/null
phase 2
testlauf "$werk/phase-2" -v "$daten":/daten -e PGR_DATEN=/daten \
  -e PGR_UPSTREAM="$pg:5432" -e PGR_OHNE_POSTGRES=1 \
  pgwire-recorder:integration "${wahl_2[@]}" || status=$?
if [ "$status" -ne 0 ]; then
  echo "run-integration-tests: Phase ohne PostgreSQL rot (Exit-Code $status)" >&2
  exit "$status"
fi
if [ "${#tests_2[@]}" -gt 0 ]; then erwartet "$werk/phase-2" 2 "${tests_2[@]}" || exit 1; fi
phase fertig
echo "run-integration-tests: gruen"
