#!/usr/bin/env bash
# kopf-check — prüft, dass der Kopf jedes lebenden Slice-Plans die Kennungen
# führt, die §1 oder §2 des Plans nennt (ADR-0032).
#
# ZUSAGE (je Punkt hält ihn tools/harness/kopf-check-gegenprobe.sh):
#   - Geprüft werden Dateien slice-*.md flach in docs/plan/planning/open/, next/
#     und in-progress/; done/ und andere Dateien nicht.
#   - Kennungen sind LH-XX-NN, LH-XX-NN.x, SPEC-NNN und ARC-NNN als ganzes Wort
#     (Wortzeichen: Buchstabe, Ziffer, Unterstrich); Auszeichnung ist ohne Belang,
#     eine andere Schreibweise ist keine Kennung.
#   - `SPEC-NNN bis SPEC-MMM` (ebenso ARC, Backticks erlaubt, NNN < MMM) steht für
#     jede Kennung dazwischen, im Kopf wie in §1 und §2, auch über einen
#     Zeilenumbruch im Absatz.
#   - Gleichheit ist exakt: Haupt- und Unterkennung decken einander nicht.
#   - Kopf sind die Zeilen, die mit **Bezug:** und **Berührte Spec-Stellen:**
#     beginnen, vor der ersten Zeile `## `, je mit ihren Folgezeilen bis zur nächsten
#     Leerzeile; geprüft wird gegen ihre Vereinigung. Eine Leerzeile darf
#     Leerzeichen und Tabs tragen.
#   - §1 und §2 reichen von der Zeile `## 1.` bzw. `## 2.` bis zur nächsten Zeile
#     `## `; der Absatz, der mit „Regeln dieser Sektion“ beginnt, zählt nicht.
#   - Fehlt ein Kopf-Feld oder einer der beiden Abschnitte, ist das ein Befund.
#
# GRENZE. Ob eine Kennung existiert, in welchem Kopf-Feld sie steht und was
# außerhalb von §1 und §2 steht, prüft das Skript nicht; eine Zeile `## ` in einem
# Codeblock beendet den Abschnitt. Eine Feldmarke zählt an jedem Zeilenanfang im
# Kopf, auch mitten in einem Absatz, den ADR-0032 Nr. 5 nennt; ein nicht lesbarer
# Plan endet je nach Position mit Exit 0 oder 2. Beides liegt beim Architect
# (Slice-Plan §6, Offen beim Architect).
#
# Aufruf: kopf-check.sh [<wurzel>]
# Ausgabe: je Befund eine Zeile auf stderr, `<pfad>: <abschnitt>: <befund>`,
# sortiert nach Pfad, Abschnitt und Kennung; jede Kennung je Abschnitt einmal.
# Ausgang: 0 ohne Befund, 1 bei mindestens einem Befund, 2 wenn
# docs/plan/planning/ fehlt.
set -euo pipefail
export LC_ALL=C

wurzel="${1:-.}"
ablage="docs/plan/planning"
if [ ! -d "$wurzel/$ablage" ]; then
  echo "kopf-check: $ablage fehlt unter '$wurzel'" >&2
  exit 2
fi
cd "$wurzel"

befunde="$(
  for lc in open next in-progress; do
    [ -d "$ablage/$lc" ] || continue
    find "$ablage/$lc" -mindepth 1 -maxdepth 1 -type f -name 'slice-*.md'
  done | sort | while IFS= read -r plan; do
    awk -v pfad="$plan" '
      function wortzeichen(c) { return c ~ /[A-Za-z0-9_]/ }
      # vorzeichen(s, st) — das Zeichen vor Position st in s, leer am Anfang.
      function vorzeichen(s, st) { return (st > 1) ? substr(s, st - 1, 1) : "" }
      # kennungen(s, menge) — trägt jede Kennung aus s in menge ein, Bereiche aufgelöst.
      function kennungen(s, menge,    off, st, t, k, n, teil, kl, kl2, a, b, i) {
        off = 0
        while (match(substr(s, off + 1), /`?(SPEC|ARC)-[0-9][0-9][0-9]`?[ \t]+bis[ \t]+`?(SPEC|ARC)-[0-9][0-9][0-9]`?/)) {
          st = off + RSTART
          t = substr(s, st, RLENGTH)
          k = vorzeichen(s, st)
          n = substr(s, st + RLENGTH, 1)
          gsub(/`/, "", t)
          split(t, teil, /[ \t]+bis[ \t]+/)
          kl = teil[1]; sub(/-.*/, "", kl)
          kl2 = teil[2]; sub(/-.*/, "", kl2)
          if (!wortzeichen(k) && !wortzeichen(n) && kl == kl2) {
            a = substr(teil[1], length(kl) + 2) + 0
            b = substr(teil[2], length(kl) + 2) + 0
            if (a < b) for (i = a; i <= b; i++) menge[sprintf("%s-%03d", kl, i)] = 1
          }
          off = st
        }
        off = 0
        while (match(substr(s, off + 1), /(LH-[A-Z][A-Z]-[0-9][0-9](\.[a-z])?|SPEC-[0-9][0-9][0-9]|ARC-[0-9][0-9][0-9])/)) {
          st = off + RSTART
          t = substr(s, st, RLENGTH)
          k = vorzeichen(s, st)
          n = substr(s, st + RLENGTH, 1)
          if (wortzeichen(n) && t ~ /\.[a-z]$/) {
            t = substr(t, 1, length(t) - 2)
            n = "."
          }
          if (!wortzeichen(k) && !wortzeichen(n)) menge[t] = 1
          off = st
        }
      }
      function befund(abschnitt, text) { printf "%s\t%s\t%s\n", pfad, abschnitt, text }
      BEGIN { imkopf = 1; feld = ""; abschnitt = ""; absatzanfang = 1; regel = 0 }
      {
        leer = ($0 ~ /^[ \t]*$/)
        if ($0 ~ /^## /) {
          imkopf = 0; feld = ""; regel = 0; absatzanfang = 1
          abschnitt = ""
          if ($0 ~ /^## 1\./) { abschnitt = "1"; hat1 = 1 }
          if ($0 ~ /^## 2\./) { abschnitt = "2"; hat2 = 1 }
          if (abschnitt != "") text[abschnitt] = text[abschnitt] " " $0
          next
        }
        if (imkopf) {
          if (leer) { feld = ""; next }
          if (index($0, "**Bezug:**") == 1) { feld = "b"; hatbezug = 1 }
          if (index($0, "**Berührte Spec-Stellen:**") == 1) { feld = "s"; hatstellen = 1 }
          if (feld != "") kopf = kopf " " $0
          next
        }
        if (abschnitt == "") next
        if (leer) { regel = 0; absatzanfang = 1; next }
        if (absatzanfang && index($0, "Regeln dieser Sektion") == 1) regel = 1
        absatzanfang = 0
        if (!regel) text[abschnitt] = text[abschnitt] " " $0
      }
      END {
        if (!hatbezug) befund("Kopf", "Feld Bezug fehlt")
        if (!hatstellen) befund("Kopf", "Feld Berührte Spec-Stellen fehlt")
        if (!hat1) befund("§1", "Abschnitt ## 1. fehlt")
        if (!hat2) befund("§2", "Abschnitt ## 2. fehlt")
        kennungen(kopf, imk)
        for (a = 1; a <= 2; a++) {
          split("", genannt)
          kennungen(text[a], genannt)
          for (id in genannt) if (!(id in imk)) befund("§" a, id " fehlt im Kopf")
        }
      }
    ' "$plan"
  done
)"

if [ -z "$befunde" ]; then
  exit 0
fi
printf '%s\n' "$befunde" | sort -t "$(printf '\t')" -k1,1 -k2,2 -k3,3 \
  | awk -F '\t' '{ printf "kopf-check: %s: %s: %s\n", $1, $2, $3 }' >&2
exit 1
