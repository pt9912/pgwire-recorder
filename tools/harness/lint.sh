#!/usr/bin/env bash
# lint — die Stufe `lint` des Dockerfile: drei eigene Prüfungen, Schema-Prüfung
# des Profils und golangci-lint nach dem Profil .golangci.yml.
# Vertrag: spec/spezifikation.md SPEC-049 (Punkte 6 bis 9).
# Aufruf: in der Wurzel des Moduls, im Image von golangci-lint (bash, find, grep,
# awk, sort, mktemp); schreibt nichts in den Arbeitsbaum.
#
# Jede Prüfung läuft, auch wenn eine andere einen Befund hat. Die Stufe schreibt
# je Befund einer eigenen Prüfung eine Zeile `lint: <pfad>:<zeile>: <befund>`,
# für das Profil (fehlt, vom Schema abgelehnt, Regel ohne Befund) eine Zeile
# `lint: .golangci.yml: <befund>`; golangci-lint schreibt seine Befunde im
# Textformat. Ausgang 1 bei mindestens einer `lint:`-Zeile oder einem Ausgang
# ungleich 0 von golangci-lint, sonst 0.
#
# Grenze: Ob ein `Why:` zutrifft und ob die Brücke export_test.go nur
# weiterreicht, prüft dieses Skript nicht.
set -uo pipefail
export LC_ALL=C

profil=".golangci.yml"
befund=0

go_dateien() {
  find cmd internal test -type f -name "$1" 2>/dev/null | sort
}

# (6) Kein `//nolint`: `//` oder `/*`, danach Leerzeichen, Tabs oder `/`, dann
# das Wort `nolint` in beliebiger Schreibweise.
tab=$'\t'
nolint_muster="(//|/\\*)[ ${tab}/]*[Nn][Oo][Ll][Ii][Nn][Tt]([^A-Za-z0-9_]|\$)"
while IFS= read -r datei; do
  while IFS= read -r treffer; do
    echo "lint: ${datei}:${treffer%%:*}: Direktive nolint"
    befund=1
  done < <(grep -nE "$nolint_muster" "$datei")
done < <(go_dateien '*.go')

# (7) Export-Test-Brücke: keine Funktion Test…, Benchmark…, Example…, Fuzz… in
# einer Datei export_test.go.
while IFS= read -r datei; do
  while IFS= read -r treffer; do
    echo "lint: ${datei}:${treffer%%:*}: Testfunktion in der Brücke export_test.go"
    befund=1
  done < <(grep -nE "^func[ ${tab}]+(Test|Benchmark|Example|Fuzz)" "$datei")
done < <(go_dateien 'export_test.go')

# (8) Über jeder Regel unter linters.exclusions.rules steht unmittelbar ein
# Kommentarblock, dessen erste Zeile mit `# Why:` beginnt.
if [ -f "$profil" ]; then
  ohne_why="$(awk -v profil="$profil" '
    function einzug(s) { match(s, /^ */); return RLENGTH }
    { zeile[NR] = $0 }
    /^[ \t]*(#|$)/ { next }
    {
      e = einzug($0)
      strich_zeile = ($0 ~ /^ *- /)
      if (in_regeln && strich < 0) {
        if (strich_zeile && e >= regel_einzug) strich = e
        else in_regeln = 0
      }
      if (in_regeln && (e < strich || (e == strich && !strich_zeile))) in_regeln = 0
      if (in_regeln && e == strich) {
        i = NR - 1
        while (i > 0 && zeile[i] ~ /^[ \t]*#/) i--
        if (i == NR - 1 || zeile[i + 1] !~ /^[ \t]*# Why:/)
          print profil ":" NR ": Regel ohne Kommentarblock \"# Why:\" unmittelbar darüber"
      }
      # Schlüssel je Einzug: linters (0) > exclusions (2) > rules.
      for (k in pfad) if (k + 0 >= e) delete pfad[k]
      if (match($0, /^ *[A-Za-z0-9_-]+:/)) {
        name = substr($0, e + 1, RLENGTH - e - 1)
        pfad[e] = name
        if (name == "rules" && pfad[2] == "exclusions" && pfad[0] == "linters") {
          in_regeln = 1; strich = -1; regel_einzug = e
        }
      }
    }
  ' "$profil")"
  if [ -n "$ohne_why" ]; then
    while IFS= read -r z; do echo "lint: $z"; done <<< "$ohne_why"
    befund=1
  fi
fi

# (9) Profil vorhanden, Schema, Befunde, ungenutzte Regel. golangci-lint liest
# das Profil nur über `-c`; fehlt es, laufen Schema-Prüfung und golangci-lint nicht.
if [ ! -f "$profil" ]; then
  echo "lint: ${profil}: fehlt"
  befund=1
else
  if ! verify="$(golangci-lint config verify -c "$profil" 2>&1)"; then
    echo "lint: ${profil}: von golangci-lint abgelehnt"
    printf '%s\n' "$verify"
    befund=1
  fi

  fehler="$(mktemp)"
  golangci-lint run -c "$profil" ./... 2> >(tee "$fehler" >&2)
  status=$?
  wait
  if [ "$status" -ne 0 ]; then
    befund=1
  fi
  # Warnung von warn-unused, logfmt-gequotet: [Text: "…", Path: "…", Linters: "…"];
  # die Werte sind darin noch einmal gequotet.
  ungenutzt="$(awk -v profil="$profil" '
    function unq(t) { gsub(/\\\\/, "\001", t); gsub(/\\"/, "\"", t); gsub(/\001/, "\\", t); return t }
    /\[runner\/exclusion_rules\] Skipped 0 issues by rules: \[/ {
      r = $0
      sub(/.*Skipped 0 issues by rules: \[/, "", r)
      sub(/\]"?[ \t]*$/, "", r)
      r = unq(r)
      delete wert
      while (match(r, /[A-Za-z]+: "([^"\\]|\\.)*"/)) {
        feld = substr(r, RSTART, RLENGTH)
        r = substr(r, RSTART + RLENGTH)
        k = feld; sub(/:.*/, "", k)
        v = feld; sub(/^[A-Za-z]+: "/, "", v); sub(/"$/, "", v)
        wert[k] = unq(v)
      }
      z = ""
      if ("Linters" in wert) z = z ", Linter: " wert["Linters"]
      if ("Path" in wert) z = z ", Pfad: " wert["Path"]
      if ("PathExcept" in wert) z = z ", Pfad außer: " wert["PathExcept"]
      if ("Text" in wert) z = z ", Text: " wert["Text"]
      if ("Source" in wert) z = z ", Quelle: " wert["Source"]
      print "lint: " profil ": Regel ohne Befund: " substr(z, 3)
    }
  ' "$fehler")"
  rm -f "$fehler"
  if [ -n "$ungenutzt" ]; then
    printf '%s\n' "$ungenutzt"
    befund=1
  fi
fi

exit "$befund"
