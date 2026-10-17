#!/usr/bin/env bash
# abdeckung — bildet die Abdeckungstabellen aus den Abdeckungs-Deklarationen der Tests
# und der Gegenproben.
# Vertrag: harness/vertraege/SPEC-048-abdeckung.md; wie ein Lauf zu lesen ist:
# harness/sensors/abdeckung-check.md.
# Aufruf: abdeckung.sh [--check] [<wurzel>] (ohne --check schreibt es die Tabellen,
# mit --check prüft es nur; ohne Wurzel das aktuelle Verzeichnis)
#
# Grenze: Das Skript liest Deklarationen, es führt weder Tests noch Gegenproben aus. Es
# prüft nicht, ob eine Gegenprobe an GATE_CHECKS hängt (das bleibt Review), nicht, ob die
# Gesamtzahl n einer geteilten Messung stimmt (nur, dass sie je Anforderung einheitlich
# ist), und nicht, ob ein Teil zur richtigen Messmethode gehört.
#
# GEPRÜFT DURCH tools/test/abdeckung-gegenprobe.sh (Fallnamen dort):
#   (1) Deklaration — nicht-gelesen, verschachtelt, leerzeile, hilfsfunktion, doppelt,
#       dateiende, eingerueckt, ohne-trenner, gueltig-maskierung
#   (2) Anforderung und Pfad — pfad, qa-happy, fa-messung, unbekannt,
#       ueberschrift-ebene, inhalt-gesamt, teil-zahlform-*, teil-fa
#   (3) Nachweisart nach Ort — inhalt-e2e, inhalt-unit, ohne-deklaration,
#       ohne-deklaration-unit, nicht-gelesen, e2e-tiefer-*
#   (4) Tabellen — inhalt-e2e, inhalt-unit, inhalt-gate-leer, inhalt-gesamt,
#       inhalt-vollstaendig, gueltig-vollstaendig, gueltig-maskierung, gate-gutform*,
#       teile*, teil-fehlt-*, teil-neun-voll
#   (5) Schreiben und Prüfen — schreiben-fehlend, schreiben-nur-geaendert,
#       schreiben-rechte, schreiben-unberuehrt, gueltig-rechte, check-aktuell,
#       check-veraltet, check-zwei-veraltet, check-schreibt-nicht, check-gate-*
#   (6) Ausgang — gueltig, jeder Fall von abgelehnt (Exit 1, Datei und Zeile auf
#       stderr), ohne-deklaration-zeile, check-veraltet, ohne-lastenheft,
#       lastenheft-ohne-anforderung, ohne-docs-user, ausgang2-*, grund-*, widerspruch-*
#   (7) Gegenprobe — gate-*, gegenprobe-ohne-deklaration
#   (8) Geteilte Messung — teil-*, widerspruch-*, messung-mehrfach
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

# Ausgang 2: nichts geprüft, nichts angelegt (Punkt 6, in dieser Reihenfolge).
if [ ! -f "$lastenheft" ]; then
  echo "abdeckung: das Lastenheft fehlt: $lastenheft" >&2
  exit 2
fi
if ! grep -q '^### LH-' "$lastenheft"; then
  echo "abdeckung: $lastenheft führt keine Überschrift ### LH-…" >&2
  exit 2
fi
if [ ! -d "$ziel" ]; then
  echo "abdeckung: das Verzeichnis $ziel fehlt" >&2
  exit 2
fi

arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT

{ grep -oE '^### LH-(FA|QA|RB)-[0-9]+ ' "$lastenheft" || true; } | sed 's/^### //; s/ $//' | sort -u > "$arbeit/bekannt.txt"

# 1. Deklarationen sammeln. Die Awk-Programme geben je Element eine Zeile
#    „ELEM\tart\tanforderung\tpfad\ttext\ttest\tdatei\tzeile\tfolge“ und je Fehlform eine Zeile
#    „FEHLER\t<datei>:<zeile>: <Meldung>“ aus. Die Gründe eines Elements stehen in der Folge
#    des Vertrags (Punkt 6): Schreibweise einschließlich Paarung, Wertebereich, Anforderung
#    unbekannt; der Widerspruch zwischen Deklarationen (5) folgt in einem zweiten Durchgang.
gemeinsam='
BEGIN {
  while ((getline id < bekannt) > 0) ok[id] = 1
  pform = "^LH-(QA|RB)-[0-9]+/(Messung|Messung-[1-9]-von-[1-9])$"
}
function fehler(m, z) { printf "FEHLER\t%s:%d: %s\n", datei, z, m }
function elemente(ids, text, name, z,    n, teile, i, t, ap, nn, a, b) {
  n = split(ids, teile, /, */)
  for (i = 1; i <= n; i++) {
    t = teile[i]
    if (art == "Gate") {
      if (t ~ /^LH-FA-/) { fehler("„" t "“ im Kopf der Gegenprobe: Ein Gate-Nachweis belegt nur LH-QA und LH-RB", z); continue }
      if (t !~ pform) {
        fehler("„" t "“ im Kopf der Gegenprobe ist nicht LH-QA/RB-<NN>/Messung oder LH-QA/RB-<NN>/Messung-<i>-von-<n>", z); continue
      }
    } else if (t !~ /^LH-FA-[0-9]+\/(Happy|Boundary|Negative)$/ && t !~ pform) {
      fehler("„" t "“ über " name " ist nicht LH-FA-<NN>/<Happy|Boundary|Negative> oder LH-QA/RB-<NN>/Messung[-<i>-von-<n>]", z); continue
    }
    split(t, ap, "/")
    if (ap[2] ~ /^Messung-/) {
      split(substr(ap[2], 9), nn, "-von-")
      a = nn[1] + 0; b = nn[2] + 0
      if (a > b) { fehler("„" t "“: Teil " a " liegt über der Gesamtzahl " b, z); continue }
      if (b < 2) { fehler("„" t "“: eine geteilte Messung hat mindestens zwei Teile", z); continue }
    }
    if (!(ap[1] in ok)) { fehler(ap[1] " steht nicht im Lastenheft", z); continue }
    printf "ELEM\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\n", art, ap[1], ap[2], text, name, datei, z, i
  }
}
function teile_text(decl, name, z,    trenner, ids, text) {
  trenner = index(decl, " — ")
  if (trenner == 0) { fehler("Deklaration ohne „ — “" (name == "-" ? "" : " über " name), z); return }
  ids = substr(decl, 1, trenner - 1)
  text = substr(decl, trenner + length(" — "))
  gsub(/\|/, "\\|", text)
  elemente(ids, text, name, z)
}
'

test_awk='
/^[ \t]+\/\/ Abdeckung:/ { fehler("eingerückte Abdeckungs-Deklaration", NR); next }
/^\/\/ Abdeckung: / { if (offen) fehler("Deklaration ohne Test", NR); decl = substr($0, 15); offen = 1; next }
offen && /^\/\/ ?/ { zeile = $0; sub(/^\/\/ ?/, "", zeile); decl = decl " " zeile; next }
/^func TestE2E/ && !offen { name = $2; sub(/\(.*/, "", name); fehler(name " ohne Abdeckungs-Deklaration", NR); next }
offen && /^func Test/ {
  name = $2; sub(/\(.*/, "", name)
  teile_text(decl, name, NR)
  offen = 0; next
}
offen { fehler("Deklaration ohne direkt folgenden Test", NR); offen = 0 }
END { if (offen) fehler("Deklaration ohne Test am Dateiende", NR) }
'

# Kopf einer Gegenprobe: die zusammenhängenden #-Zeilen ab Zeile 1.
gate_awk='
function abschluss() { if (offen) { teile_text(decl, "-", dz); offen = 0 } }
BEGIN { kopf = 1 }
{
  if (kopf && $0 !~ /^#/) { abschluss(); kopf = 0 }
  if ($0 ~ /^[ \t]+# Abdeckung:/) fehler("eingerückte Abdeckungs-Deklaration", NR)
  if (kopf) {
    if ($0 ~ /^# Abdeckung: /) { abschluss(); decl = substr($0, 14); dz = NR; offen = 1 }
    else if (offen && $0 ~ /^# / && $0 !~ /^# *$/) { zeile = $0; sub(/^# ?/, "", zeile); decl = decl " " zeile }
    else abschluss()
  } else if ($0 ~ /^# Abdeckung: /) fehler("Abdeckungs-Deklaration nach dem Kopf der Gegenprobe", NR)
}
END { abschluss() }
'

{
  find . -path ./.git -prune -o -path ./.harness -prune -o -name '*_test.go' -print \
    | sed 's|^\./||' | sort | while read -r datei; do
    art=Unit
    case "$datei" in test/integration/*) art=E2E ;; esac
    awk -v datei="$datei" -v art="$art" -v bekannt="$arbeit/bekannt.txt" "$gemeinsam$test_awk" "$datei"
  done
  if [ -d tools ]; then
    find tools -type f -name '*-gegenprobe.sh' | sort | while read -r datei; do
      awk -v datei="$datei" -v art=Gate -v bekannt="$arbeit/bekannt.txt" "$gemeinsam$gate_awk" "$datei"
    done
  fi
} > "$arbeit/p1.tsv"

# Zweiter Durchgang: Widerspruch zwischen Deklarationen (Grund 5), nur über Elemente ohne
# Grund aus (1) bis (4); ein Fehler steht an der späteren Fundstelle (Datei bytegenau, Zeile)
# und nennt die frühere. decl.tsv: art \t anforderung \t pfad \t text \t test \t datei.
{ grep '^FEHLER' "$arbeit/p1.tsv" || true; } > "$arbeit/decl.tsv"
{ grep '^ELEM' "$arbeit/p1.tsv" || true; } \
  | sort -t$'\t' -k7,7 -k8,8n -k9,9n \
  | awk -F'\t' '
    {
      id = $3; pf = $4; ort = $7 ":" $8
      geteilt = (pf ~ /^Messung-/)
      n = substr(pf, length(pf))
      if (geteilt && (id in fn) && fn[id] != n) {
        printf "FEHLER\t%s: %s/%s widerspricht %s/%s (früher: %s)\n", ort, id, pf, id, fpf[id], fort[id]; next
      }
      if (geteilt && (id in fp)) {
        printf "FEHLER\t%s: %s/%s steht neben %s/Messung (früher: %s)\n", ort, id, pf, id, fp[id]; next
      }
      if (pf == "Messung" && (id in fn)) {
        printf "FEHLER\t%s: %s/Messung steht neben %s/%s (früher: %s)\n", ort, id, id, fpf[id], fort[id]; next
      }
      if (geteilt && !(id in fn)) { fn[id] = n; fpf[id] = pf; fort[id] = ort }
      if (pf == "Messung" && !(id in fp)) fp[id] = ort
      printf "%s\t%s\t%s\t%s\t%s\t%s\n", $2, $3, $4, $5, $6, $7
    }' >> "$arbeit/decl.tsv"

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

{
  cat <<'EOF'
# Gate-Abdeckung je Anforderung und Pfad

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`) aus den
Abdeckungs-Deklarationen im Kopf der Gegenproben unter `tools/`. Die Datei ist eine
Abdeckungs-Deklaration, kein Lauf-Beleg; dass die Gegenproben laufen, sichert
`make gates`. Die Gesamtsicht steht in `abdeckung-gesamt.md`.

| Anforderung | Pfad | Kurzbeschreibung | Nachweis |
| --- | --- | --- | --- |
EOF
  awk -F'\t' '$1 == "Gate"' "$arbeit/decl.tsv" \
    | sort -t$'\t' -k2,2V -k3,3 -k6,6 \
    | while IFS=$'\t' read -r _ id pfad text _ datei; do
        printf '| %s | %s | %s | `%s` |\n' "$(link "$id")" "$pfad" "$text" "$datei"
      done
} > "$arbeit/abdeckung-gate.md"

# 2. Gesamtsicht und vollständige Anforderungen. Vollständig ist eine LH-FA mit
#    Happy, Boundary und Negative, eine LH-QA oder LH-RB mit Messung oder mit allen Teilen
#    1 bis n einer geteilten Messung. Kopplung: abdeckung-vollstaendig.md liest
#    `make doc-trace` (trace.coverage in .d-check.yml).
awk -F'\t' '{ print $2 "\t" $3 "\t" $1 }' "$arbeit/decl.tsv" | sort -u > "$arbeit/pfade.tsv"
cut -f1 "$arbeit/pfade.tsv" | sort -uV > "$arbeit/ids.txt"

arten() { awk -F'\t' -v id="$1" -v pfad="$2" '$1 == id && $2 == pfad { print $3 }' "$arbeit/pfade.tsv" | sort -u | paste -sd, - | sed 's/,/, /g'; }

# gesamtzahl <id> — das n der geteilten Messung der Anforderung, leer ohne eine.
gesamtzahl() { awk -F'\t' -v id="$1" '$1 == id && $2 ~ /^Messung-/ { print substr($2, length($2)); exit }' "$arbeit/pfade.tsv"; }

# messung <id> — die Spalte Messung: die Nachweisarten des Pfads `Messung` oder alle Teile 1 bis n.
messung() {
  local id="$1" n i a s=""
  n="$(gesamtzahl "$id")"
  if [ -z "$n" ]; then
    arten "$id" Messung
    return
  fi
  for ((i = 1; i <= n; i++)); do
    a="$(arten "$id" "Messung-$i-von-$n")"
    s="${s:+$s; }$i/$n ${a:-—}"
  done
  printf '%s' "$s"
}

vollstaendig() {
  local id="$1" n i
  case "$id" in
    LH-FA-*) [ -n "$(arten "$id" Happy)" ] && [ -n "$(arten "$id" Boundary)" ] && [ -n "$(arten "$id" Negative)" ] ;;
    *)
      n="$(gesamtzahl "$id")"
      if [ -z "$n" ]; then
        [ -n "$(arten "$id" Messung)" ]
        return
      fi
      for ((i = 1; i <= n; i++)); do
        [ -n "$(arten "$id" "Messung-$i-von-$n")" ] || return 1
      done
      ;;
  esac
}

{
  cat <<'EOF'
# Abdeckung je Anforderung

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`) aus den
Abdeckungs-Deklarationen aller Tests und der Gegenproben. Je Anforderung zeigt die
Tabelle, welche Nachweisart welchen Pfad belegt: bei funktionalen Anforderungen die
drei Akzeptanzkriterien des Lastenhefts, bei Qualitätsanforderungen und
Randbedingungen die Messung, bei einer geteilten Messung jeden ihrer Teile.
Vollständig belegte Anforderungen stehen zusätzlich in
`abdeckung-vollstaendig.md`, die `make doc-trace` liest.

| Anforderung | Happy | Boundary | Negative | Messung | Stand |
| --- | --- | --- | --- | --- | --- |
EOF
  while read -r id; do
    h="$(arten "$id" Happy)"; b="$(arten "$id" Boundary)"; n="$(arten "$id" Negative)"; m="$(messung "$id")"
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
Anforderung erst, wenn Tests und Gegenproben alle ihre Pfade belegen: bei
funktionalen Anforderungen Happy, Boundary und Negative, bei Qualitätsanforderungen
und Randbedingungen die Messung oder alle Teile einer geteilten Messung. Die Spalte
nennt die beteiligten Nachweisarten. Diese Datei liest
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
for f in abdeckung-e2e.md abdeckung-unit.md abdeckung-gate.md abdeckung-gesamt.md abdeckung-vollstaendig.md; do
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
