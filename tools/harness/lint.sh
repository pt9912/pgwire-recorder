#!/usr/bin/env bash
# lint — die Stufe `lint` des Dockerfile: drei eigene Prüfungen, Schema-Prüfung
# des Profils und golangci-lint nach dem Profil .golangci.yml. Als Gate läuft sie
# über `make lint` an der Gate-Kette von `make gates`.
# Vertrag: spec/spezifikation.md SPEC-049 (Punkte 6 bis 9 für dieses Skript); wie
# ein Lauf zu lesen ist: harness/sensors/lint.md.
# Aufruf: in der Wurzel des Moduls, im Image von golangci-lint (bash, find, grep,
# awk, sort, mktemp).
#
# GEPRÜFT DURCH tools/harness/lint-gegenprobe.sh (Fälle je Punkt von SPEC-049; was
# offen ist, nennt ihr Kopf):
#   (1) Gegenstand — p1-*
#   (2) Werkzeug und Umgebung — kein eigener Fall; offen im Kopf der Gegenprobe
#   (3) Linter — p3-*, p4-*-rot, p5-revive-*, p5-contextcheck-test, p7-whitebox-*
#   (4) Schwellen — p4-*
#   (5) Einstellungen — p5-*, p1-cmd
#   (6) Kein `//nolint` — p6-* (auch in einer Testdatei), p9-profil-fehlt-nolint
#   (7) Export-Test-Brücke — p7-*, p9-profil-fehlt-bruecke
#   (8) Ausnahmen — p8-*, p0-grundlauf (keine ungenutzte Regel im Bestand)
#   (9) Ausgabe und Ausgang — p9-*, p8-form-* (Ausgang 1), p0-grundlauf (Ausgang 0)
#   (10) Werkzeug, dann Gate — p10-*
#   Grenze — neg-main-test, neg-generiert, p9-laden, p8-form-eintrag-2
#
# Jede Prüfung läuft, auch wenn eine andere einen Befund hat. Die Stufe schreibt
# je Befund einer eigenen Prüfung eine Zeile `lint: <pfad>:<zeile>: <befund>`,
# für das Profil (fehlt, vom Schema abgelehnt, Regel ohne Befund) eine Zeile
# `lint: .golangci.yml: <befund>`; golangci-lint schreibt seine Befunde im
# Textformat. Ausgang 1 bei mindestens einer `lint:`-Zeile oder einem Ausgang
# ungleich 0 von golangci-lint, sonst 0.
#
# Grenze: Ob ein `Why:` zutrifft und ob die Brücke export_test.go nur
# weiterreicht, prüft dieses Skript nicht. Die Prüfung nach (8) liest das Profil
# nur in der festen Form, die SPEC-049 Punkt 8 nennt; jede andere meldet sie als
# `Form nicht erkannt`.
set -uo pipefail
export LC_ALL=C

profil=".golangci.yml"
befund=0

go_dateien() {
  find cmd internal test -type f -name "$1" 2>/dev/null | sort
}

# (6) Kein `//nolint`: `//` oder `/*`, danach Leerzeichen, Tabs oder `/`, dann
# das Wort `nolint` in beliebiger Schreibweise. Maßgeblich ist jedes `//` und
# `/*` der Zeile: `// x //nolint` ist ein Befund, `// siehe nolint` nicht.
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
# Kommentarblock, dessen erste Zeile mit `# Why:` beginnt. Die Prüfung liest das
# Profil in seiner festen Form: `linters:` ohne Einzug, `exclusions:` mit 2,
# `rules:` mit 4, jede Regel `- ` mit 6 Leerzeichen, alle drei in Blockform.
# Ein Schlüssel `exclusions` an anderer Stelle oder in Flussform, ein `rules`
# darunter mit anderem Einzug oder in Flussform und unter `rules` jede Zeile, die
# weder Kommentar, Leerzeile, Eintrag `- ` mit 6 noch dessen Fortsetzung mit
# mindestens 8 ist, ergeben `Form nicht erkannt` mit ihrer Zeile.
if [ -f "$profil" ]; then
  ohne_why="$(awk -v profil="$profil" '
    function form(n) { print profil ":" n ": Form nicht erkannt" }
    { zeile[NR] = $0 }
    /^[ \t]*(#|$)/ { next }
    {
      t = $0
      sub(/[ \t]+#.*$/, "", t)
      match(t, /^ */); e = RLENGTH
      rest = substr(t, e + 1)
      sub(/[ \t]+$/, "", rest)
      strich = (rest ~ /^-( |$)/)
      if (e == 0) {
        oben = rest; sub(/:.*/, "", oben)
        in_ex = 0; in_regeln = 0
      }
      if (in_ex && e <= ex_e && !(in_regeln && strich)) { in_ex = 0; in_regeln = 0 }
      if (in_regeln) {
        # Unter `rules`: Eintrag `- ` mit 6, Fortsetzung mit mindestens 8 unter
        # einem Eintrag; jede andere Zeile, auch ein `- ` mit anderem Einzug, ist
        # `Form nicht erkannt`. Eine Zeile ohne `- ` mit höchstens 4 beendet `rules`.
        if (!strich && e <= 4) in_regeln = 0
        else if (e >= 8 && eintrag) { }
        else if (!strich || e != 6) form(NR)
        else {
          eintrag = 1
          i = NR - 1
          while (i > 0 && zeile[i] ~ /^[ \t]*#/) i--
          if (i == NR - 1 || zeile[i + 1] !~ /^[ \t]*# Why:/)
            print profil ":" NR ": Regel ohne Kommentarblock \"# Why:\" unmittelbar darüber"
        }
      }
      if (t ~ /(^|[{,]|- ) *["\047]?exclusions["\047]? *:/) {
        ex_ok = (e == 2 && oben == "linters" && rest == "exclusions:")
        if (!ex_ok) form(NR)
        in_ex = 1; ex_e = e; in_regeln = 0
      } else if (in_ex && t ~ /(^|[{,]|- ) *["\047]?rules["\047]? *:/) {
        if (ex_ok && e == 4 && rest == "rules:") { in_regeln = 1; eintrag = 0 }
        else form(NR)
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
  # die Werte sind darin noch einmal gequotet. Die Zeile nennt die Felder in der
  # Reihenfolge Linter, Pfad, Pfad außer, Text, Quelle, nur die vorhandenen; ohne
  # eines davon lautet sie `Regel ohne Befund: Felder nicht erkannt`.
  ungenutzt="$(awk -v profil="$profil" '
    function unq(t) { gsub(/\\\\/, "\001", t); gsub(/\\"/, "\"", t); gsub(/\001/, "\\", t); return t }
    /\[runner\/exclusion_rules\] Skipped 0 issues by rules: \[/ {
      r = $0
      sub(/.*Skipped 0 issues by rules: \[/, "", r)
      sub(/\]"?[ \t]*$/, "", r)
      r = unq(r)
      delete wert
      while (match(r, /[A-Za-z][A-Za-z ]*: "([^"\\]|\\.)*"/)) {
        feld = substr(r, RSTART, RLENGTH)
        r = substr(r, RSTART + RLENGTH)
        k = feld; sub(/:.*/, "", k)
        v = feld; sub(/^[A-Za-z][A-Za-z ]*: "/, "", v); sub(/"$/, "", v)
        wert[k] = unq(v)
      }
      z = ""
      if ("Linters" in wert) z = z ", Linter: " wert["Linters"]
      if ("Path" in wert) z = z ", Pfad: " wert["Path"]
      if ("Path Except" in wert) z = z ", Pfad außer: " wert["Path Except"]
      if ("Text" in wert) z = z ", Text: " wert["Text"]
      if ("Source" in wert) z = z ", Quelle: " wert["Source"]
      if (z == "") print "lint: " profil ": Regel ohne Befund: Felder nicht erkannt"
      else print "lint: " profil ": Regel ohne Befund: " substr(z, 3)
    }
  ' "$fehler")"
  rm -f "$fehler"
  if [ -n "$ungenutzt" ]; then
    printf '%s\n' "$ungenutzt"
    befund=1
  fi
fi

exit "$befund"
