#!/usr/bin/env bash
# abdeckung — bildet die Abdeckungstabellen aus den Abdeckungs-Deklarationen der
# Tests. Eine Deklaration steht direkt über `func Test…` und beginnt mit
#
#   // Abdeckung: <Anforderung>/<Pfad>, … — <Kurzbeschreibung>
#
# Folgezeilen mit `//` setzen sie fort. Pfad ist Happy, Boundary oder Negative
# (die drei Akzeptanzkriterien des Lastenhefts). Die Nachweisart folgt aus dem
# Ort: test/integration/ ist E2E, alles andere Unit.
#
# Geschrieben werden unter docs/user/:
#   e2e-abdeckung.md          je E2E-Deklaration und Anforderung eine Zeile
#   unit-abdeckung.md         dasselbe für Unit-Tests
#   abdeckung-gesamt.md       je Anforderung die drei Pfade mit ihren Nachweisarten
#   abdeckung-vollstaendig.md nur Anforderungen mit allen drei Pfaden belegt;
#                             diese Datei liest `make doc-trace` (trace.coverage)
#
# Aufruf: abdeckung.sh           schreibt die Tabellen
#         abdeckung.sh --check   prüft, dass sie aktuell sind, und schreibt nichts
# Ausgang: 0 bei Erfolg; 1 bei einer fehlerhaften Deklaration oder, mit --check,
# bei einer veralteten Tabelle.
set -euo pipefail

modus="${1:-}"
ziel="docs/user"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT

# 1. Deklarationen sammeln: art \t anforderung \t pfad \t text \t test \t datei
find . -path ./.git -prune -o -path ./.harness -prune -o -name '*_test.go' -print \
  | sed 's|^\./||' | sort | while read -r datei; do
  art=Unit
  case "$datei" in test/integration/*) art=E2E ;; esac
  awk -v datei="$datei" -v art="$art" '
    /^\/\/ Abdeckung: / { decl = substr($0, 15); offen = 1; next }
    offen && /^\/\/ ?/ { zeile = $0; sub(/^\/\/ ?/, "", zeile); decl = decl " " zeile; next }
    offen && /^func Test/ {
      name = $2; sub(/\(.*/, "", name)
      trenner = index(decl, " — ")
      if (trenner == 0) { printf "FEHLER\t%s: Deklaration ohne „ — “ über %s\n", datei, name; offen = 0; next }
      ids = substr(decl, 1, trenner - 1)
      text = substr(decl, trenner + length(" — "))
      n = split(ids, teile, /, */)
      for (i = 1; i <= n; i++) {
        if (teile[i] !~ /^LH-(FA|QA|RB)-[0-9]+\/(Happy|Boundary|Negative)$/) {
          printf "FEHLER\t%s: „%s“ über %s ist nicht <Anforderung>/<Happy|Boundary|Negative>\n", datei, teile[i], name
          continue
        }
        split(teile[i], ap, "/")
        printf "%s\t%s\t%s\t%s\t%s\t%s\n", art, ap[1], ap[2], text, name, datei
      }
      offen = 0; next
    }
    { offen = 0 }
  ' "$datei"
done > "$arbeit/decl.tsv"

if grep -q '^FEHLER' "$arbeit/decl.tsv"; then
  grep '^FEHLER' "$arbeit/decl.tsv" | cut -f2 >&2
  exit 1
fi

link() { printf '[`%s`](../../spec/lastenheft.md)' "$1"; }

kopf_art() {
  cat <<EOF
# $1-Abdeckung je Anforderung und Pfad

Erzeugt von \`make abdeckung\` (\`tools/test/abdeckung.sh\`) aus den
Abdeckungs-Deklarationen über den Tests ($2). Die Datei ist eine
Abdeckungs-Deklaration, kein Lauf-Beleg; dass die Tests grün laufen, sichert
\`make gates\`. Die Gesamtsicht steht in \`abdeckung-gesamt.md\`.

| Anforderung | Pfad | Kurzbeschreibung | Nachweis | Ort |
| --- | --- | --- | --- | --- |
EOF
}

tabelle_art() {
  local art="$1" titel="$2" quelle="$3"
  {
    kopf_art "$titel" "$quelle"
    awk -F'\t' -v art="$art" '$1 == art' "$arbeit/decl.tsv" \
      | sort -t$'\t' -k2,2V -k3,3 -k5,5 \
      | while IFS=$'\t' read -r _ id pfad text name datei; do
          printf '| %s | %s | %s | `%s` | `%s` |\n' "$(link "$id")" "$pfad" "$text" "$name" "$datei"
        done
  } > "$arbeit/$4"
}

tabelle_art E2E E2E "\`test/integration/\`, gegen eine reale PostgreSQL-Instanz" e2e-abdeckung.md
tabelle_art Unit Unit "\`internal/\`" unit-abdeckung.md

# 2. Gesamtsicht und vollständige Anforderungen.
awk -F'\t' '{ print $2 "\t" $3 "\t" $1 }' "$arbeit/decl.tsv" | sort -u > "$arbeit/pfade.tsv"
cut -f1 "$arbeit/pfade.tsv" | sort -uV > "$arbeit/ids.txt"

arten() { awk -F'\t' -v id="$1" -v pfad="$2" '$1 == id && $2 == pfad { print $3 }' "$arbeit/pfade.tsv" | sort -u | paste -sd, - | sed 's/,/, /g'; }

{
  cat <<'EOF'
# Abdeckung je Anforderung

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`) aus den
Abdeckungs-Deklarationen aller Tests. Je Anforderung zeigt die Tabelle, welche
Nachweisart welchen der drei Pfade des Lastenhefts belegt. Vollständig ist eine
Anforderung, wenn alle drei Pfade belegt sind; nur diese stehen in
`abdeckung-vollstaendig.md`, die `make doc-trace` liest.

| Anforderung | Happy | Boundary | Negative | Stand |
| --- | --- | --- | --- | --- |
EOF
  while read -r id; do
    h="$(arten "$id" Happy)"; b="$(arten "$id" Boundary)"; n="$(arten "$id" Negative)"
    stand=teilweise
    [ -n "$h" ] && [ -n "$b" ] && [ -n "$n" ] && stand=vollständig
    printf '| %s | %s | %s | %s | %s |\n' "$(link "$id")" "${h:-—}" "${b:-—}" "${n:-—}" "$stand"
  done < "$arbeit/ids.txt"
} > "$arbeit/abdeckung-gesamt.md"

{
  cat <<'EOF'
# Vollständig abgedeckte Anforderungen

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`). Hier steht eine
Anforderung erst, wenn Tests alle drei Pfade des Lastenhefts (Happy, Boundary,
Negative) belegen; die Spalte nennt die beteiligten Nachweisarten. Diese Datei
liest `make doc-trace` (`trace.coverage`); Teilabdeckung steht in
`abdeckung-gesamt.md`.

| Anforderung | Nachweisarten |
| --- | --- |
EOF
  while read -r id; do
    h="$(arten "$id" Happy)"; b="$(arten "$id" Boundary)"; n="$(arten "$id" Negative)"
    if [ -n "$h" ] && [ -n "$b" ] && [ -n "$n" ]; then
      alle="$(awk -F'\t' -v id="$id" '$1 == id { print $3 }' "$arbeit/pfade.tsv" | sort -u | paste -sd, - | sed 's/,/, /g')"
      printf '| %s | %s |\n' "$(link "$id")" "$alle"
    fi
  done < "$arbeit/ids.txt"
} > "$arbeit/abdeckung-vollstaendig.md"

# 3. Schreiben oder prüfen.
veraltet=0
for f in e2e-abdeckung.md unit-abdeckung.md abdeckung-gesamt.md abdeckung-vollstaendig.md; do
  if ! cmp -s "$arbeit/$f" "$ziel/$f"; then
    if [ "$modus" = "--check" ]; then
      echo "abdeckung: $ziel/$f ist nicht aktuell — make abdeckung schreibt sie neu" >&2
      veraltet=1
    else
      cp "$arbeit/$f" "$ziel/$f"
      chmod 0644 "$ziel/$f"
      echo "abdeckung: $ziel/$f geschrieben"
    fi
  fi
done
exit "$veraltet"
