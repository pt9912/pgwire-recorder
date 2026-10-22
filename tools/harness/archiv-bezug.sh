#!/usr/bin/env bash
# archiv-bezug — schreibt die Anforderungs-Kennungen aus dem Kopf des archivierten Volltexts als Zeile `**Bezug:**` in den Stub
# (Vertrag: harness/vertraege/SPEC-068-archiv-bezug.md, Punkte 1 bis 12). Werkzeug, kein Gate; es committet nichts.
# Aufruf: archiv-bezug.sh <wurzel> [<welle-id>]  (ohne <welle-id> alle Verzeichnisse done/<welle-id>/ mit einer Datei slice-*.md)
# Exit 0: fertig (bei leerer Menge mit der Zeile `archiv-bezug: keine Stubs` auf stdout); 1: mindestens ein Befund (Punkt 3), nichts geschrieben; 2: Eingabe, Werkzeug oder Schreibfehler (Punkte 2 und 3).
# GEPRÜFT DURCH tools/dev/harness/archiv-bezug-gegenprobe.sh, Fälle nach Punkt (p<n>-*).
export LC_ALL=C

wurzel="${1-}"
welle="${2-}"
done_dir="$wurzel/docs/plan/planning/done"
rel="docs/plan/planning/done"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

sperre() { echo "archiv-bezug: $*" >&2; exit 2; }

# Punkt 2: Eingaben, in dieser Reihenfolge
[ -d "$wurzel" ] || sperre "Wurzel fehlt oder ist kein Verzeichnis: ${wurzel:-(leer)}"
if [ -n "$welle" ]; then
  [[ $welle =~ ^welle-[a-z0-9-]+$ ]] || sperre "Form der Welle ungültig: $welle"
  [ -d "$done_dir/$welle" ] || sperre "Verzeichnis fehlt: $rel/$welle"
fi
command -v unzip > /dev/null 2>&1 || sperre "Werkzeug fehlt: unzip"

# Punkt 1: Lauf-Menge
alle=()
if [ -n "$welle" ]; then
  alle=("$welle")
else
  for d in "$done_dir"/*/; do alle+=("$(basename "$d")"); done
fi
wellen=()
for w in "${alle[@]}"; do
  compgen -G "$done_dir/$w/slice-*.md" > /dev/null && wellen+=("$w")
done
for w in "${wellen[@]}"; do [ -f "$done_dir/$w/archiv.zip" ] || sperre "archiv.zip fehlt: $w"; done
for w in "${wellen[@]}"; do unzip -tq "$done_dir/$w/archiv.zip" > /dev/null 2>&1 || sperre "archiv.zip beschädigt: $w/archiv.zip"; done

# Punkt 4 und 5: die Kennungen aus dem Kopf, sortiert (Gruppe, Zahl, Länge der Schreibweise, Schreibweise), ohne Duplikat
kennungen() {
  awk '
    {
      if (lies == 0) { if (index($0, "**Bezug:**") != 1) next; lies = 1 }
      else if (lies == 1 && index($0, "**") == 1) lies = 2
      if (lies != 1) next
      t = " " $0
      while (match(t, /LH-(FA|QA|RB)-[0-9]+/)) {
        if (substr(t, RSTART - 1, 1) !~ /[A-Za-z0-9-]/) print (substr(t, RSTART + 3, 2) == "FA" ? 1 : substr(t, RSTART + 3, 2) == "QA" ? 2 : 3), substr(t, RSTART + 6, RLENGTH - 6), substr(t, RSTART, RLENGTH), RLENGTH
        t = substr(t, RSTART + RLENGTH - 1)
      }
    }
  ' "$1" | sort -k1,1n -k2,2n -k4,4n -k3,3 -u | cut -d" " -f3
}
zeile() {
  local k z=""
  while read -r k; do z+="${z:+, }$k"; done < <(kennungen "$1")
  echo "**Bezug:** ${z:-— keine —}"
}

schreibe() {
  local stub="$1" z="$2" t="$1.$$"
  Z="$z" awk '{ print } index($0, "**Hervorgegangen:**") == 1 { print ENVIRON["Z"] }' "$stub" > "$t" && mv "$t" "$stub"
}

# Punkt 3: alle Stubs prüfen, bevor einer geschrieben wird
befunde=0
befund() { echo "archiv-bezug: $1: $2" >&2; befunde=1; }
todo=()
gesehen=0
declare -A neu
for w in "${wellen[@]}"; do
  d="$done_dir/$w"
  unzip -Z1 "$d/archiv.zip" > "$tmp/liste"
  for stub in "$d"/slice-*.md; do
    gesehen=$((gesehen + 1))
    name="${stub##*/}"
    pfad="$rel/$w/$name"
    grep -Fxq -- "$pfad" "$tmp/liste" || { befund "$pfad" "kein Eintrag im Archiv"; continue; }
    unzip -p "$d/archiv.zip" "$pfad" > "$tmp/voll"
    awk '{ sub(/\r$/, "") } $0 == "---" { gef = 1; exit } { print } END { exit gef ? 0 : 1 }' "$tmp/voll" > "$tmp/kopf" || { befund "$pfad" "Eintrag im Archiv ohne Kopf (keine Zeile ---)"; continue; }
    if grep -q $'\r' "$stub"; then befund "$pfad" "Stub enthält CR"; continue; fi
    [ "$(grep -c '^\*\*Hervorgegangen:\*\*' "$stub")" -eq 1 ] || { befund "$pfad" "Stub hat nicht genau eine Zeile **Hervorgegangen:**"; continue; }
    soll="$(zeile "$tmp/kopf")"
    anzahl="$(grep -c '^\*\*Bezug:\*\*' "$stub")"
    if [ "$anzahl" -gt 1 ]; then befund "$pfad" "mehr als eine Zeile **Bezug:**"; continue; fi
    if [ "$anzahl" -eq 1 ]; then
      ist="$(grep '^\*\*Bezug:\*\*' "$stub")"
      [ "$ist" = "$soll" ] || { befund "$pfad" "Zeile weicht ab: ist \"$ist\", soll \"$soll\""; continue; }
      [ "$(awk 'f { print; exit } index($0, "**Hervorgegangen:**") == 1 { f = 1 }' "$stub")" = "$soll" ] || { befund "$pfad" "Zeile an falscher Stelle"; continue; }
    fi
    if [ "$anzahl" -eq 0 ]; then
      if [ ! -w "$stub" ] || [ ! -w "$d" ]; then befund "$pfad" "nicht beschreibbar"; continue; fi
      todo+=("$stub"); neu["$stub"]="$soll"
    fi
  done
done
if [ "$gesehen" -eq 0 ]; then echo "archiv-bezug: keine Stubs"; exit 0; fi
[ "$befunde" -eq 0 ] || exit 1

for stub in "${todo[@]}"; do
  schreibe "$stub" "${neu[$stub]}" || { rm -f "$stub.$$"; echo "archiv-bezug: $rel/${stub#"$done_dir"/}: Schreiben fehlgeschlagen" >&2; exit 2; }
done
exit 0
