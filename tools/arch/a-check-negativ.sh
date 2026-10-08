#!/usr/bin/env bash
# a-check-negativ — Gegenprobe des Architektur-Gates in Kopien des Arbeitsbaums.
# Der Arbeitsbaum selbst bleibt unberuehrt. Zwoelf Faelle:
#
#   1. internal/hexagon/model importiert pgproto3      → a-check meldet die Datei
#   2. internal/adapters/driven/recording importiert
#      pgproto3                                         → a-check meldet tech-leak
#   3. internal/adapters/driven/recording importiert
#      crypto/tls                                       → a-check meldet tech-leak
#   4. internal/adapters/driven/postgres importiert
#      pgproto3 und crypto/tls                          → a-check meldet nichts
#   5. internal/adapters/driving/pgwire importiert
#      pgproto3 und crypto/tls                          → a-check meldet nichts
#   6. internal/hexagon/model importiert go.yaml.in/yaml/v3 → a-check meldet die Datei
#   7. internal/adapters/driven/postgres importiert
#      gopkg.in/yaml.v3                                 → a-check meldet tech-leak
#   8. internal/adapters/driving/cli importiert
#      go.yaml.in/yaml/v3 und gopkg.in/yaml.v3          → a-check meldet nichts
#   9. internal/adapters/driving/pgwire importiert
#      go.yaml.in/yaml/v3                               → a-check meldet tech-leak
#  10. internal/adapters/driving/pgwire importiert
#      gopkg.in/yaml.v3                                 → a-check meldet tech-leak
#  11. internal/hexagon/model importiert gopkg.in/yaml.v3 → a-check meldet die Datei
#  12. internal/adapters/driven/postgres importiert
#      go.yaml.in/yaml/v3                               → a-check meldet tech-leak
#
# Die Faelle 4 und 5 halten fest, dass beide PGWire-Adapter beide Bibliotheken
# nutzen duerfen (ADR-0010); die Faelle 1 bis 3, dass ein anderer Ort abgelehnt
# wird. Die Faelle 6 bis 12 gelten den beiden YAML-Modulpfaden (ADR-0027,
# ADR-0036): Fall 8 laesst sie im CLI-Adapter zu, die Faelle 6, 7 und 9 bis 12
# lehnen jeden Modulpfad im Domain Model, im Postgres-Adapter und im
# PGWire-Adapter ab. Andere Regeln von .a-check.yml prueft die Gegenprobe nicht.
#
# Ausgang: 0, wenn alle Faelle das erwartete Ergebnis liefern, sonst 1.
# Aufruf ueber `make a-check-negativ` (harness/mk/arch-negativ.mk), das Image und
# Runtime aus a-check.mk uebernimmt.
set -euo pipefail

image="${A_CHECK_IMAGE:?A_CHECK_IMAGE fehlt}"
docker="${DOCKER:-docker}"

arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT

fehler=0

# fall <name> <datei> <package> <imports> <erwartung> <muster>
#   erwartung: "rot" (a-check muss scheitern und <muster> ausgeben) oder "gruen"
fall() {
  local name="$1" datei="$2" paket="$3" imports="$4" erwartung="$5" muster="$6"
  local kopie="$arbeit/$name" log="$arbeit/$name.log"
  mkdir -p "$kopie"
  tar --exclude=./.git -cf - . | tar -xf - -C "$kopie"
  {
    printf 'package %s\n\nimport (\n' "$paket"
    for i in $imports; do printf '\t_ "%s"\n' "$i"; done
    printf ')\n'
  } > "$kopie/$datei"

  local ergebnis=gruen
  "$docker" run --rm --network none -v "$kopie":/src:ro "$image" /src >"$log" 2>&1 || ergebnis=rot

  if [ "$ergebnis" != "$erwartung" ]; then
    echo "a-check-negativ: ROT in Fall '$name' — erwartet $erwartung, a-check war $ergebnis" >&2
    cat "$log" >&2
    fehler=1
  elif [ "$erwartung" = rot ] && ! grep -q -- "$muster" "$log"; then
    echo "a-check-negativ: ROT in Fall '$name' — a-check scheitert, nennt aber '$muster' nicht" >&2
    cat "$log" >&2
    fehler=1
  fi
}

fall model internal/hexagon/model/negativprobe.go model \
  "github.com/jackc/pgx/v5/pgproto3" rot "internal/hexagon/model/negativprobe.go"
fall recording internal/adapters/driven/recording/negativprobe.go recording \
  "github.com/jackc/pgx/v5/pgproto3" rot "tech-leak"
fall recording-tls internal/adapters/driven/recording/negativprobe.go recording \
  "crypto/tls" rot "tech-leak"
fall postgres internal/adapters/driven/postgres/negativprobe.go postgres \
  "github.com/jackc/pgx/v5/pgproto3 crypto/tls" gruen ""
fall pgwire internal/adapters/driving/pgwire/negativprobe.go pgwire \
  "github.com/jackc/pgx/v5/pgproto3 crypto/tls" gruen ""
fall model-yaml internal/hexagon/model/negativprobe.go model \
  "go.yaml.in/yaml/v3" rot "internal/hexagon/model/negativprobe.go"
fall postgres-yaml internal/adapters/driven/postgres/negativprobe.go postgres \
  "gopkg.in/yaml.v3" rot "tech-leak"
fall cli-yaml internal/adapters/driving/cli/negativprobe.go cli \
  "go.yaml.in/yaml/v3 gopkg.in/yaml.v3" gruen ""
fall pgwire-yaml internal/adapters/driving/pgwire/negativprobe.go pgwire \
  "go.yaml.in/yaml/v3" rot "tech-leak"
fall pgwire-gopkg-yaml internal/adapters/driving/pgwire/negativprobe.go pgwire \
  "gopkg.in/yaml.v3" rot "tech-leak"
fall model-gopkg-yaml internal/hexagon/model/negativprobe.go model \
  "gopkg.in/yaml.v3" rot "internal/hexagon/model/negativprobe.go"
fall postgres-yamlin internal/adapters/driven/postgres/negativprobe.go postgres \
  "go.yaml.in/yaml/v3" rot "tech-leak"

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "a-check-negativ: gruen — PGWire- und TLS-Bibliothek nur in den PGWire-Adaptern, YAML nur im Recording- und CLI-Adapter"
