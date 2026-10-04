#!/usr/bin/env bash
# a-check-negativ — Gegenprobe des Architektur-Gates: In einer Kopie des
# Arbeitsbaums importiert internal/hexagon/model die PGWire-Bibliothek (pgproto3);
# a-check muss diese Kopie ablehnen. Der Arbeitsbaum selbst bleibt unberuehrt.
#
# Ausgang: 0, wenn a-check die Verletzung meldet; 1, wenn es sie durchlaesst.
# Aufruf ueber `make a-check-negativ` (harness/mk/arch-negativ.mk), das Image und
# Runtime aus a-check.mk uebernimmt.
set -euo pipefail

image="${A_CHECK_IMAGE:?A_CHECK_IMAGE fehlt}"
docker="${DOCKER:-docker}"

kopie="$(mktemp -d)"
trap 'rm -rf "$kopie"' EXIT

tar --exclude=./.git -cf - . | tar -xf - -C "$kopie"

cat > "$kopie/internal/hexagon/model/verletzung_negativprobe.go" <<'EOF'
package model

import _ "github.com/jackc/pgx/v5/pgproto3"
EOF

if "$docker" run --rm --network none -v "$kopie":/src:ro "$image" /src >"$kopie.log" 2>&1; then
  echo "a-check-negativ: ROT — a-check laesst den pgproto3-Import im Domain Model durch" >&2
  cat "$kopie.log" >&2
  rm -f "$kopie.log"
  exit 1
fi

if ! grep -q "pgproto3" "$kopie.log"; then
  echo "a-check-negativ: ROT — a-check scheitert, nennt aber den pgproto3-Import nicht" >&2
  cat "$kopie.log" >&2
  rm -f "$kopie.log"
  exit 1
fi

rm -f "$kopie.log"
echo "a-check-negativ: gruen — a-check lehnt den pgproto3-Import im Domain Model ab"
