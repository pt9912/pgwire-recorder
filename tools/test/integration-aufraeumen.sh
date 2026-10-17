#!/usr/bin/env bash
# integration-aufraeumen — Funktionen des Aufräumens von tools/test/run-integration-tests.sh
# (Vertrag: harness/vertraege/SPEC-057-integration-aufraeumen.md). Der Runner und
# tools/test/integration-volumes-gegenprobe.sh laden diese Datei; sie setzen vorher
# docker, pg, tests und daten (Namen der Container und des benannten Volumes des Laufs).
# Die Funktionen schreiben nichts außer Meldungen auf stderr und entfernen nur Container und Volumes,
# deren Namen sie so gelesen haben (Punkt 4).

# integration_volumen_lesen — Namen der Volumes, die `docker inspect` an $tests und $pg nennt
# (Mounts vom Typ volume, außer $daten), einer je Zeile. Ein Container, der nicht existiert,
# nennt keine Namen (Punkt 3).
integration_volumen_lesen() {
  local c n name
  for c in "$tests" "$pg"; do
    n="$("$docker" container inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}' "$c" 2> /dev/null)" || continue
    for name in $n; do
      if [ "$name" != "$daten" ]; then printf '%s\n' "$name"; fi
    done
  done
}

# integration_entfernen — entfernt $tests und $pg samt ihrer anonymen Volumes (Punkt 2).
integration_entfernen() {
  "$docker" rm -fv "$tests" "$pg" > /dev/null 2>&1 || true
}

# integration_pruefen <Name>… — nennt ein Volume, das nach dem Entfernen noch existiert, entfernt es
# per Namen und endet mit 1 (Punkt 4). Ohne Namen oder ohne Rest endet sie mit 0 und schweigt.
integration_pruefen() {
  local n rest=""
  for n; do
    if "$docker" volume inspect "$n" > /dev/null 2>&1; then rest="$rest $n"; fi
  done
  if [ -z "$rest" ]; then return 0; fi
  echo "run-integration-tests: Volumes nach dem Entfernen noch vorhanden:$rest" >&2
  for n in $rest; do "$docker" volume rm "$n" > /dev/null 2>&1 || true; done
  return 1
}

# integration_aufraeumen — Namen lesen, entfernen, prüfen; 1 bei einem Rest (Punkte 2 bis 4).
integration_aufraeumen() {
  local namen
  namen="$(integration_volumen_lesen)" || return 1
  integration_entfernen
  # shellcheck disable=SC2086
  integration_pruefen $namen
}

# integration_ende <Code des Laufs> — räumt auf und endet mit dem Code des Laufs, wenn er ungleich 0
# ist, sonst mit 1 bei einem Rest und sonst mit 0 (Punkt 4, Vorrang des Codes).
integration_ende() {
  local ende_code="$1" ende_rest=0
  integration_aufraeumen || ende_rest=1
  if [ "$ende_code" -ne 0 ]; then return "$ende_code"; fi
  return "$ende_rest"
}
