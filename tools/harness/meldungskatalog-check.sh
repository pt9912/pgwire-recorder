#!/usr/bin/env bash
# meldungskatalog-check — gleicht die Meldungscodes PGR-… der Code-Tabelle im Quelltext (C), des
# Katalogs in SPEC-034 (K) und der Tabellen des Benutzerhandbuchs (H) ab.
# Vertrag: harness/vertraege/SPEC-054-meldungskatalog.md (Punkte 1 bis 11 und Grenze); wie ein Lauf
# zu lesen ist: harness/sensors/meldungskatalog-check.md.
# Aufruf: meldungskatalog-check.sh [<wurzel>] (ohne Wurzel das aktuelle Verzeichnis)
#
# GEPRÜFT DURCH tools/harness/meldungskatalog-check-gegenprobe.sh, Fälle nach Punkt (nr<n>-*):
#   (1) Quellen, (2) Form, (3) doppelt, (4) Mengen (nr4a bis nr4d), (5) Rückfall je Klasse,
#   (6) Code außerhalb der Code-Tabelle, (8) Sperren und ihre Reihenfolge, (9) Ausgabe und Ausgang,
#   (10) Aufrufarten, (11) Anschluss ohne Stufung
set -euo pipefail
export LC_ALL=C

C="internal/hexagon/model/fehler.go"
K="spec/spezifikation.md"
H="docs/user/benutzerhandbuch.md"

sperre() { echo "meldungskatalog-check: $*" >&2; exit 2; }

wurzel="${1:-.}"
cd -- "$wurzel" 2>/dev/null || sperre "$C fehlt"
for datei in "$C" "$K" "$H"; do
  [ -f "$datei" ] || sperre "$datei fehlt"
  [ -r "$datei" ] || sperre "$datei nicht lesbar"
done
for verzeichnis in cmd internal; do
  [ -d "$verzeichnis" ] || sperre "$verzeichnis/ fehlt"
done

# Der Baum für Punkt 6: Go-Dateien unter cmd/ und internal/ ohne Tests und ohne die Code-Tabelle.
# Ein nicht lesbares Verzeichnis oder eine nicht lesbare Datei des Baums ist eine Sperre; gemeldet
# wird der erste Pfad in Bytefolge.
baum=()
unlesbar=()
while IFS= read -r -d '' pfad; do
  if [ -d "$pfad" ]; then
    { [ -r "$pfad" ] && [ -x "$pfad" ]; } || unlesbar+=("$pfad")
  else
    [ -r "$pfad" ] && baum+=("$pfad") || unlesbar+=("$pfad")
  fi
done < <(find cmd internal \( -type d -o \( -type f -name '*.go' ! -name '*_test.go' ! -path "$C" \) \) -print0 2>/dev/null)
if [ "${#unlesbar[@]}" -gt 0 ]; then
  sperre "$(printf '%s\n' "${unlesbar[@]}" | sort | head -n 1) nicht lesbar"
fi

# Je Befund eine Zeile pfad<TAB>code<TAB>befund<TAB>zeile; awk endet mit 2 bei einer Sperre.
if ! befunde="$(awk -v c="$C" -v k="$K" -v h="$H" '
  function sperre(grund) { print "meldungskatalog-check: " grund > "/dev/stderr"; exit 2 }
  function ausgabe(pfad, code, befund, zeile) { printf "%s\t%s\t%s\t%s\n", pfad, code, befund, zeile }
  # tabellenzeile(s) — setzt zelle auf den Text zwischen den Backticks der ersten Zelle.
  function tabellenzeile(s,    i) {
    if (substr(s, 1, 1) != "|") return 0
    s = substr(s, 2)
    i = index(s, "|")
    if (i > 0) s = substr(s, 1, i - 1)
    sub(/^[ \t]+/, "", s)
    sub(/[ \t]+$/, "", s)
    if (length(s) < 2 || substr(s, 1, 1) != "`" || substr(s, length(s), 1) != "`") return 0
    zelle = substr(s, 2, length(s) - 2)
    return 1
  }
  function eintrag(q, pfad, zeile, code) {
    anzahl[q]++
    if (code !~ /^PGR-[EWI][0-9][0-9][0-9][0-9]$/) { ausgabe(pfad, code, "Form", zeile); return }
    if ((q, code) in erste) {
      if (!((q, code) in gemeldet)) { gemeldet[q, code] = 1; ausgabe(pfad, code, "doppelt", zeile) }
      return
    }
    erste[q, code] = zeile
    liste[q, ++n[q]] = code
  }
  FNR == 1 {
    art = (FILENAME == c) ? "C" : (FILENAME == k) ? "K" : (FILENAME == h) ? "H" : "T"
    imabschnitt = 0
  }
  art == "C" {
    if ($0 ~ /^[ \t]*Code[A-Za-z0-9_]*[ \t]+([A-Za-z0-9_.]+[ \t]+)?=[ \t]+"[^"]*"/) {
      match($0, /"[^"]*"/)
      eintrag("C", c, FNR, substr($0, RSTART + 1, RLENGTH - 2))
    }
    next
  }
  art == "K" {
    if (index($0, "### SPEC-034 ") == 1) { imabschnitt = 1; hatk = 1 }
    else if (substr($0, 1, 1) == "#") imabschnitt = 0
    else if (imabschnitt && tabellenzeile($0)) eintrag("K", k, FNR, zelle)
    next
  }
  art == "H" {
    if ($0 == "### Fehlercodes") { imabschnitt = 1; hatf = 1 }
    else if ($0 == "### Warnungen") { imabschnitt = 1; hatw = 1 }
    else if (substr($0, 1, 1) == "#") imabschnitt = 0
    else if (imabschnitt && tabellenzeile($0)) eintrag("H", h, FNR, zelle)
    next
  }
  {
    p = index($0, "\"PGR-")
    if (p > 0) {
      rest = substr($0, p + 5)
      match(rest, /^[A-Za-z0-9_]*/)
      ausgabe(FILENAME, "PGR-" substr(rest, 1, RLENGTH), "Code außerhalb der Code-Tabelle", FNR)
    }
  }
  END {
    if (!hatk) sperre("Abschnitt ### SPEC-034 fehlt")
    if (!hatf) sperre("Abschnitt ### Fehlercodes fehlt")
    if (!hatw) sperre("Abschnitt ### Warnungen fehlt")
    if (!anzahl["C"]) sperre(c " ohne Code")
    if (!anzahl["K"]) sperre(k " ohne Codes")
    if (!anzahl["H"]) sperre(h " ohne Codes")
    for (i = 1; i <= n["C"]; i++) {
      code = liste["C", i]
      if (!(("K", code) in erste)) ausgabe(c, code, "nicht im Katalog", erste["C", code])
      if (!(("H", code) in erste)) ausgabe(c, code, "nicht im Handbuch", erste["C", code])
    }
    for (i = 1; i <= n["K"]; i++) {
      code = liste["K", i]
      if (!(("C", code) in erste) && code !~ /000$/) ausgabe(k, code, "nicht in der Code-Tabelle", erste["K", code])
    }
    for (i = 1; i <= n["H"]; i++) {
      code = liste["H", i]
      if (!(("C", code) in erste)) ausgabe(h, code, "nicht in der Code-Tabelle", erste["H", code])
    }
    for (i = 1; i <= 6; i++) {
      code = "PGR-E" i "000"
      if (!(("K", code) in erste)) ausgabe(k, code, "Rückfall fehlt", "")
    }
  }
' "$C" "$K" "$H" ${baum[@]+"${baum[@]}"})"; then
  exit 2
fi

if [ -z "$befunde" ]; then
  exit 0
fi
printf '%s\n' "$befunde" | sort -t "$(printf '\t')" -k1,1 -k2,2 -k3,3 -k4,4n \
  | awk -F '\t' '{
      if ($4 == "") printf "meldungskatalog-check: %s: %s: %s\n", $1, $2, $3
      else printf "meldungskatalog-check: %s:%s: %s: %s\n", $1, $4, $2, $3
    }' >&2
exit 1
