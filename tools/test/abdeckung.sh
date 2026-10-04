#!/usr/bin/env bash
# abdeckung — bildet die Abdeckungstabellen aus den Abdeckungs-Deklarationen der
# Tests. Eine Deklaration steht direkt über `func Test…` (ohne Leerzeile) und
# beginnt mit
#
#   // Abdeckung: <Anforderung>/<Pfad>, … — <Kurzbeschreibung>
#
# Folgezeilen mit `//` setzen sie fort. Pfad ist bei funktionalen Anforderungen
# (LH-FA) eines der drei Akzeptanzkriterien Happy, Boundary oder Negative, bei
# Qualitätsanforderungen (LH-QA) und Randbedingungen (LH-RB) die Messung. Jede
# Anforderung muss im Lastenheft als Überschrift stehen. Die Nachweisart folgt
# aus dem Ort: test/integration/ ist E2E, alles andere Unit. Jeder `TestE2E…`
# trägt eine Deklaration.
#
# Geschrieben werden unter docs/user/:
#   abdeckung-e2e.md          je E2E-Deklaration und Anforderung eine Zeile
#   abdeckung-unit.md         dasselbe für Unit-Tests
#   abdeckung-gesamt.md       je Anforderung die Pfade mit ihren Nachweisarten
#   abdeckung-vollstaendig.md nur Anforderungen mit allen Pfaden belegt; diese
#                             Datei liest `make doc-trace` (trace.coverage)
#
# Aufruf: abdeckung.sh [--check] [<wurzel>]
#   ohne --check schreibt es die Tabellen, mit --check prüft es nur
# Ausgang: 0 bei Erfolg; 1 bei einer fehlerhaften Deklaration oder, mit --check,
# bei einer veralteten Tabelle.
set -euo pipefail
export LC_ALL=C

modus=""
if [ "${1:-}" = "--check" ]; then
  modus="--check"
  shift
fi
wurzel="${1:-.}"
cd "$wurzel"

ziel="docs/user"
lastenheft="spec/lastenheft.md"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT

grep -oE '^### LH-(FA|QA|RB)-[0-9]+ ' "$lastenheft" | sed 's/^### //; s/ $//' | sort -u > "$arbeit/bekannt.txt"

# 1. Deklarationen sammeln: art \t anforderung \t pfad \t text \t test \t datei;
#    Fehler als Zeile „FEHLER\t<Meldung>“.
find . -path ./.git -prune -o -path ./.harness -prune -o -name '*_test.go' -print \
  | sed 's|^\./||' | sort | while read -r datei; do
  art=Unit
  case "$datei" in test/integration/*) art=E2E ;; esac
  awk -v datei="$datei" -v art="$art" -v bekannt="$arbeit/bekannt.txt" '
    BEGIN { while ((getline id < bekannt) > 0) ok[id] = 1 }
    function fehler(m) { printf "FEHLER\t%s:%d: %s\n", datei, NR, m }
    /^[ \t]+\/\/ Abdeckung:/ { fehler("eingerückte Abdeckungs-Deklaration"); next }
    /^\/\/ Abdeckung: / { if (offen) fehler("Deklaration ohne Test"); decl = substr($0, 15); offen = 1; next }
    offen && /^\/\/ ?/ { zeile = $0; sub(/^\/\/ ?/, "", zeile); decl = decl " " zeile; next }
    /^func TestE2E/ && !offen { name = $2; sub(/\(.*/, "", name); fehler(name " ohne Abdeckungs-Deklaration"); next }
    offen && /^func Test/ {
      name = $2; sub(/\(.*/, "", name)
      trenner = index(decl, " — ")
      if (trenner == 0) { fehler("Deklaration ohne „ — “ über " name); offen = 0; next }
      ids = substr(decl, 1, trenner - 1)
      text = substr(decl, trenner + length(" — "))
      gsub(/\|/, "\\|", text)
      n = split(ids, teile, /, */)
      for (i = 1; i <= n; i++) {
        if (teile[i] !~ /^LH-FA-[0-9]+\/(Happy|Boundary|Negative)$/ && teile[i] !~ /^LH-(QA|RB)-[0-9]+\/Messung$/) {
          fehler("„" teile[i] "“ über " name " ist nicht LH-FA-<NN>/<Happy|Boundary|Negative> oder LH-QA/RB-<NN>/Messung")
          continue
        }
        split(teile[i], ap, "/")
        if (!(ap[1] in ok)) { fehler(ap[1] " über " name " steht nicht im Lastenheft"); continue }
        printf "%s\t%s\t%s\t%s\t%s\t%s\n", art, ap[1], ap[2], text, name, datei
      }
      offen = 0; next
    }
    offen { fehler("Deklaration ohne direkt folgenden Test"); offen = 0 }
    END { if (offen) fehler("Deklaration ohne Test am Dateiende") }
  ' "$datei"
done > "$arbeit/decl.tsv"

if grep -q '^FEHLER' "$arbeit/decl.tsv"; then
  grep '^FEHLER' "$arbeit/decl.tsv" | cut -f2 | sed 's/^/abdeckung: /' >&2
  exit 1
fi

link() { printf '[`%s`](../../spec/lastenheft.md)' "$1"; }

tabelle_art() {
  local art="$1" quelle="$2" datei="$3"
  {
    cat <<EOF
# $art-Abdeckung je Anforderung und Pfad

Erzeugt von \`make abdeckung\` (\`tools/test/abdeckung.sh\`) aus den
Abdeckungs-Deklarationen über den Tests ($quelle). Die Datei ist eine
Abdeckungs-Deklaration, kein Lauf-Beleg; dass die Tests grün laufen, sichert
\`make gates\`. Die Gesamtsicht steht in \`abdeckung-gesamt.md\`.

| Anforderung | Pfad | Kurzbeschreibung | Nachweis | Ort |
| --- | --- | --- | --- | --- |
EOF
    awk -F'\t' -v art="$art" '$1 == art' "$arbeit/decl.tsv" \
      | sort -t$'\t' -k2,2V -k3,3 -k5,5 \
      | while IFS=$'\t' read -r _ id pfad text name datei; do
          printf '| %s | %s | %s | `%s` | `%s` |\n' "$(link "$id")" "$pfad" "$text" "$name" "$datei"
        done
  } > "$arbeit/$datei"
}

tabelle_art E2E "\`test/integration/\`, gegen eine reale PostgreSQL-Instanz" abdeckung-e2e.md
tabelle_art Unit "\`internal/\`" abdeckung-unit.md

# 2. Gesamtsicht und vollständige Anforderungen. Vollständig ist eine LH-FA mit
#    Happy, Boundary und Negative, eine LH-QA oder LH-RB mit Messung.
awk -F'\t' '{ print $2 "\t" $3 "\t" $1 }' "$arbeit/decl.tsv" | sort -u > "$arbeit/pfade.tsv"
cut -f1 "$arbeit/pfade.tsv" | sort -uV > "$arbeit/ids.txt"

arten() { awk -F'\t' -v id="$1" -v pfad="$2" '$1 == id && $2 == pfad { print $3 }' "$arbeit/pfade.tsv" | sort -u | paste -sd, - | sed 's/,/, /g'; }

vollstaendig() {
  case "$1" in
    LH-FA-*) [ -n "$(arten "$1" Happy)" ] && [ -n "$(arten "$1" Boundary)" ] && [ -n "$(arten "$1" Negative)" ] ;;
    *) [ -n "$(arten "$1" Messung)" ] ;;
  esac
}

{
  cat <<'EOF'
# Abdeckung je Anforderung

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`) aus den
Abdeckungs-Deklarationen aller Tests. Je Anforderung zeigt die Tabelle, welche
Nachweisart welchen Pfad belegt: bei funktionalen Anforderungen die drei
Akzeptanzkriterien des Lastenhefts, bei Qualitätsanforderungen die Messung.
Vollständig belegte Anforderungen stehen zusätzlich in
`abdeckung-vollstaendig.md`, die `make doc-trace` liest.

| Anforderung | Happy | Boundary | Negative | Messung | Stand |
| --- | --- | --- | --- | --- | --- |
EOF
  while read -r id; do
    h="$(arten "$id" Happy)"; b="$(arten "$id" Boundary)"; n="$(arten "$id" Negative)"; m="$(arten "$id" Messung)"
    stand=teilweise
    vollstaendig "$id" && stand=vollständig
    case "$id" in
      LH-FA-*) m="n/a" ;;
      *) h="n/a"; b="n/a"; n="n/a" ;;
    esac
    printf '| %s | %s | %s | %s | %s | %s |\n' "$(link "$id")" "${h:-—}" "${b:-—}" "${n:-—}" "${m:-—}" "$stand"
  done < "$arbeit/ids.txt"
} > "$arbeit/abdeckung-gesamt.md"

{
  cat <<'EOF'
# Vollständig abgedeckte Anforderungen

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`). Hier steht eine
Anforderung erst, wenn Tests alle ihre Pfade belegen: bei funktionalen
Anforderungen Happy, Boundary und Negative, bei Qualitätsanforderungen die
Messung. Die Spalte nennt die beteiligten Nachweisarten. Diese Datei liest
`make doc-trace` (`trace.coverage`); Teilabdeckung steht in
`abdeckung-gesamt.md`.

| Anforderung | Nachweisarten |
| --- | --- |
EOF
  while read -r id; do
    if vollstaendig "$id"; then
      alle="$(awk -F'\t' -v id="$id" '$1 == id { print $3 }' "$arbeit/pfade.tsv" | sort -u | paste -sd, - | sed 's/,/, /g')"
      printf '| %s | %s |\n' "$(link "$id")" "$alle"
    fi
  done < "$arbeit/ids.txt"
} > "$arbeit/abdeckung-vollstaendig.md"

# 3. Schreiben oder prüfen.
veraltet=0
for f in abdeckung-e2e.md abdeckung-unit.md abdeckung-gesamt.md abdeckung-vollstaendig.md; do
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
