#!/usr/bin/env bash
# run-integration-tests — startet eine PostgreSQL-Instanz und das Testimage in
# einem eigenen Docker-Netz, fuehrt die Integrationstests aus und raeumt danach
# auf. Danach schreibt es docs/user/e2e-abdeckung.md aus den Zeilen
# `// Abdeckung: <Kennungen> — <Text>` ueber jedem `func TestE2E*` unter
# test/integration/, aber nur, wenn sich der Inhalt aendert.
#
# Ausgang: der Exit-Code des Testlaufs; 1, wenn Aufbau oder Abdeckung scheitern.
set -euo pipefail

docker="${DOCKER:-docker}"
postgres_image="${POSTGRES_IMAGE:?POSTGRES_IMAGE fehlt}"
lauf="pgr-it-$$"
netz="$lauf-net"
pg="$lauf-pg"

aufraeumen() {
  "$docker" rm -f "$pg" >/dev/null 2>&1 || true
  "$docker" network rm "$netz" >/dev/null 2>&1 || true
}
trap aufraeumen EXIT

"$docker" build --progress=plain --target integration -t pgwire-recorder:integration .
"$docker" network create --internal "$netz" >/dev/null
"$docker" run -d --name "$pg" --network "$netz" \
  -e POSTGRES_HOST_AUTH_METHOD=trust "$postgres_image" >/dev/null

status=0
"$docker" run --rm --network "$netz" -e PGR_UPSTREAM="$pg:5432" \
  pgwire-recorder:integration || status=$?

if [ "$status" -ne 0 ]; then
  echo "run-integration-tests: Testlauf rot (Exit-Code $status)" >&2
  "$docker" logs "$pg" 2>&1 | tail -n 20 >&2 || true
  exit "$status"
fi

# Abdeckungstabelle aus den Deklarationen der Tests.
ziel="docs/user/e2e-abdeckung.md"
neu="$(mktemp)"
trap 'rm -f "$neu"; aufraeumen' EXIT
{
  cat <<'KOPF'
# E2E-Abdeckung je Anforderung

Erzeugt von `make test-integration` über `tools/test/run-integration-tests.sh`
aus den Zeilen `// Abdeckung: …` direkt über jedem `func TestE2E*` unter
`test/integration/`. Die Datei ändert sich nur, wenn sich eine Deklaration oder
ihr Ort ändert. Sie ist eine Abdeckungs-Deklaration, kein Lauf-Beleg.

| Anforderung | Kurzbeschreibung | Nachweis | Ort |
| --- | --- | --- | --- |
KOPF
  for datei in $(find test/integration -name '*_test.go' | sort); do
    awk -v datei="$datei" '
      /^\/\/ Abdeckung: / { deklaration = substr($0, 15); zeile = 0; next }
      /^func TestE2E/ {
        if (deklaration == "") { printf "FEHLT %s:%d\n", datei, NR; next }
        name = $2; sub(/\(.*/, "", name)
        trenner = index(deklaration, " — ")
        kennungen = substr(deklaration, 1, trenner - 1)
        text = substr(deklaration, trenner + length(" — "))
        n = split(kennungen, ids, /, */)
        links = ""
        for (i = 1; i <= n; i++) {
          links = links (i > 1 ? ", " : "") "[`" ids[i] "`](../../spec/lastenheft.md)"
        }
        printf "| %s | %s | `%s` | `%s:%d` |\n", links, text, name, datei, NR
        deklaration = ""
        next
      }
      { if ($0 !~ /^\/\//) deklaration = "" }
    ' "$datei"
  done
} > "$neu"

if grep -q '^FEHLT ' "$neu"; then
  echo "run-integration-tests: TestE2E ohne Abdeckungs-Zeile:" >&2
  grep '^FEHLT ' "$neu" >&2
  exit 1
fi

if ! cmp -s "$neu" "$ziel"; then
  cp "$neu" "$ziel"
  chmod 0644 "$ziel"
  echo "run-integration-tests: $ziel aktualisiert"
fi
echo "run-integration-tests: gruen"
