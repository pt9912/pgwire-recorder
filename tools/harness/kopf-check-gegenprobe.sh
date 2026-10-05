#!/usr/bin/env bash
# kopf-check-gegenprobe — prüft tools/harness/kopf-check.sh (ADR-0032) an kleinen
# Bäumen in einem Temp-Verzeichnis; der Arbeitsbaum bleibt unberührt. Je Fall ein
# Baum mit einem oder mehreren Plänen; erwartet ist entweder Exit 0 ohne Ausgabe
# (grün) oder Exit 1 mit einer bestimmten Befund-Zeile auf stderr (rot). Die Fälle
# sind nach den Nummern der ADR gruppiert; die akzeptierten Negative der ADR
# (Kennung außerhalb von §1/§2, Feldwahl, Existenz der Kennung) stehen als
# grüne Fälle darunter.
#
# Ausgang: 0, wenn jeder Fall wie erwartet endet, sonst 1.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
skript="$hier/kopf-check.sh"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT
fehler=0
nr=0
P="docs/plan/planning"

melde() { echo "kopf-check-gegenprobe: ROT — $*" >&2; fehler=1; }

# plan <datei> <bezug> <stellen> <§1> <§2> [<§3>] — schreibt einen Slice-Plan.
plan() {
  mkdir -p "$(dirname "$1")"
  {
    printf '# Slice: Beispiel\n\n**Welle:** ohne Welle.\n\n'
    printf '**Bezug:** %s\n\n**Berührte Spec-Stellen:** %s\n\n' "$2" "$3"
    printf '**Verantwortlich:** —\n\n---\n\n## 1. Ziel und Abgrenzung\n\n'
    printf 'Regeln dieser Sektion: Baseline-Regelwerk LH-QA-99 und SPEC-999.\n\n'
    printf '%s\n\n## 2. Definition of Done\n\n' "$4"
    printf 'Regeln dieser Sektion: ARC-999.\n\n%s\n\n## 3. Plan\n\n%s\n' "$5" "${6:-—}"
  } > "$1"
}

# neu — legt einen leeren Baum an und setzt w auf seine Wurzel.
neu() {
  nr=$((nr + 1))
  w="$arbeit/$nr"
  mkdir -p "$w/$P/open" "$w/$P/next" "$w/$P/in-progress" "$w/$P/done"
}

# lauf <wurzel> — führt das Skript aus; setzt code, aus (stdout), err (stderr).
lauf() {
  set +e
  bash "$skript" "$1" >"$arbeit/out" 2>"$arbeit/err"
  code=$?
  set -e
  aus="$(cat "$arbeit/out")"
  err="$(cat "$arbeit/err")"
}

# gruen <name> <wurzel> — erwartet Exit 0 und keine Ausgabe.
gruen() {
  lauf "$2"
  if [ "$code" -ne 0 ] || [ -n "$aus$err" ]; then
    melde "Fall '$1' erwartet grün, Exit $code: $err"
  fi
}

# rot <name> <wurzel> <zeile>… — erwartet Exit 1, leeres stdout und genau die
# genannten Befund-Zeilen (ohne Präfix „kopf-check: “) in dieser Reihenfolge.
rot() {
  local name="$1" wurzel="$2" soll
  shift 2
  lauf "$wurzel"
  soll="$(printf 'kopf-check: %s\n' "$@")"
  if [ "$code" -ne 1 ]; then
    melde "Fall '$name' erwartet Exit 1, bekam $code: $err"
  elif [ -n "$aus" ]; then
    melde "Fall '$name' schreibt auf stdout: $aus"
  elif [ "$err" != "$soll" ]; then
    melde "Fall '$name' erwartet:"$'\n'"$soll"$'\n'"bekam:"$'\n'"$err"
  fi
}

# einzeln_gruen <name> <bezug> <stellen> <§1> <§2> [<§3>]
# einzeln_rot <name> <befunde> <bezug> <stellen> <§1> <§2> [<§3>]
# — je ein Baum mit einem Plan in open/; <befunde> sind die erwarteten Zeilen
# ohne das Pfad-Präfix, eine je Zeile.
P1="$P/open/slice-a.md"
einzeln_gruen() {
  local name="$1" w
  shift
  neu; plan "$w/$P1" "$@"; gruen "$name" "$w"
}
einzeln_rot() {
  local name="$1" befund="$2" w
  shift 2
  neu; plan "$w/$P1" "$@"
  local zeilen=()
  while IFS= read -r z; do zeilen+=("$P1: $z"); done <<<"$befund"
  rot "$name" "$w" "${zeilen[@]}"
}

# --- Vollständiger Kopf ------------------------------------------------------
einzeln_gruen vollstaendig \
  "[\`LH-FA-01\`](../../spec/lastenheft.md#lh-fa-01--x)" "\`LH-FA-02.a\` · \`SPEC-001\` · \`ARC-001\`" \
  "Liefert LH-FA-01 und LH-FA-02.a." "- [ ] SPEC-001 und ARC-001 belegt."
einzeln_gruen leerer-kopf-ohne-nennung "—" "—" "Keine Kennung." "- [ ] nichts."

# --- Nr. 2: Kennungs-Klassen, je in §1 und §2 --------------------------------
for k in LH-FA-01 LH-QA-05.b SPEC-001 ARC-002; do
  einzeln_rot "nr2-$k-§1" "§1: $k fehlt im Kopf" "—" "—" "Nennt $k." "- [ ] nichts."
  einzeln_rot "nr2-$k-§2" "§2: $k fehlt im Kopf" "—" "—" "Nichts." "- [ ] Nennt $k."
done
# Auszeichnung ohne Belang.
einzeln_rot nr2-code-span "§1: SPEC-001 fehlt im Kopf" "—" "—" "Nennt \`SPEC-001\`." "-"
einzeln_rot nr2-linktext "§1: LH-FA-01 fehlt im Kopf" "—" "—" "[\`LH-FA-01\`](x.md#lh-fa-01--y)" "-"
einzeln_rot nr2-hervorhebung "§1: ARC-001 fehlt im Kopf" "—" "—" "**ARC-001** und *ARC-001*" "-"
einzeln_rot nr2-codeblock "§2: SPEC-007 fehlt im Kopf" "—" "—" "-" "$(printf '```\nSPEC-007\n```')"
einzeln_rot nr2-satzende "§1: LH-FA-03 fehlt im Kopf" "—" "—" "Siehe LH-FA-03." "-"
einzeln_rot nr2-zeilenanfang "§1: SPEC-010 fehlt im Kopf" "—" "—" "SPEC-010 steht vorn" "-"
# Ganzes Wort mit Wortzeichen Buchstabe, Ziffer, Unterstrich: ein Punkt endet das Wort.
einzeln_rot nr2-punkt-endet-wort "§1: LH-FA-03 fehlt im Kopf" "—" "—" "LH-FA-03.ab" "-"
# Link-Anker sind keine Nennung.
einzeln_gruen nr2-anker "—" "—" "[Abschnitt](../x.md#lh-fa-01--y) und [z](a.md#spec-001)" "-"
# Andere Schreibweise ist keine Kennung — in §1/§2 …
einzeln_gruen nr2-schreibweise-abschnitt "—" "—" \
  "lh-fa-01, LH-FA-1, LH-FA-01a, LH-FA01, SPEC-01, SPEC-0001, spec-001, ARC-1, XLH-FA-01, LH-FA-01_x, SPEC-001x" "-"
# … und im Kopf.
einzeln_rot nr2-schreibweise-kopf "§1: LH-FA-01 fehlt im Kopf" "\`lh-fa-01\`" "\`LH-FA-1\` · LH-FA-01a" "Nennt LH-FA-01." "-"

# --- Nr. 3: Bereich ----------------------------------------------------------
einzeln_gruen nr3-bereich-kopf "—" "\`SPEC-013\` bis \`SPEC-019\` · ARC-001 bis ARC-003" "SPEC-016 und ARC-002" "-"
einzeln_gruen nr3-bereich-abschnitt "—" "SPEC-013 · SPEC-014 · SPEC-015" "\`SPEC-013\` bis \`SPEC-015\`" "-"
einzeln_gruen nr3-bereich-zeilenumbruch "—" "$(printf 'SPEC-013 bis\nSPEC-019')" "SPEC-015" "-"
einzeln_rot nr3-bereich-abschnitt-luecke "$(printf '§1: SPEC-014 fehlt im Kopf\n§1: SPEC-015 fehlt im Kopf')" \
  "—" "SPEC-013 · SPEC-016" "SPEC-013 bis SPEC-016" "-"
einzeln_rot nr3-absteigend "§1: SPEC-015 fehlt im Kopf" "—" "SPEC-019 bis SPEC-013" "SPEC-015" "-"
einzeln_rot nr3-gemischt "§1: SPEC-015 fehlt im Kopf" "—" "SPEC-013 bis ARC-019" "SPEC-015" "-"
einzeln_rot nr3-andere-form "§1: SPEC-015 fehlt im Kopf" "—" "SPEC-013–SPEC-019 · SPEC-013 - SPEC-019" "SPEC-015" "-"
einzeln_rot nr3-kein-ganzes-wort "§1: SPEC-015 fehlt im Kopf" "—" "XSPEC-013 bis SPEC-019 · SPEC-013 bis SPEC-0190" "SPEC-015" "-"

# --- Nr. 4: Gleichheit ist exakt ---------------------------------------------
einzeln_rot nr4-haupt-deckt-unter "§1: LH-FA-10.a fehlt im Kopf" "[\`LH-FA-10\`](x.md)" "—" "LH-FA-10.a" "-"
einzeln_rot nr4-unter-deckt-haupt "§1: LH-FA-10 fehlt im Kopf" "—" "\`LH-FA-10.a\`" "LH-FA-10" "-"

# --- Nr. 5: Kopf -------------------------------------------------------------
einzeln_gruen nr5-nur-bezug "LH-FA-01 · SPEC-001" "—" "LH-FA-01, SPEC-001" "-"
einzeln_gruen nr5-nur-stellen "—" "LH-FA-01 · SPEC-001" "LH-FA-01, SPEC-001" "-"
einzeln_gruen nr5-mehrzeilig "$(printf 'erste Zeile\nzweite Zeile SPEC-004')" "—" "SPEC-004" "-"
einzeln_gruen nr5-mehrfach "SPEC-004 · SPEC-004" "SPEC-004" "SPEC-004 und SPEC-004" "-"
einzeln_gruen nr5-kopf-ohne-nennung "SPEC-004 · LH-FA-09" "ARC-005" "Nichts." "-"
einzeln_rot nr5-nach-leerzeile "§1: SPEC-004 fehlt im Kopf" "$(printf 'x\n\nSPEC-004')" "—" "SPEC-004" "-"
einzeln_rot nr5-leerzeile-mit-leerzeichen "§1: SPEC-004 fehlt im Kopf" "$(printf 'x\n  \t\nSPEC-004')" "—" "SPEC-004" "-"
# Ein anderes Kopf-Feld zählt nicht.
neu; plan "$w/$P1" "—" "—" "SPEC-004" "-"
sed -i 's/^\*\*Welle:\*\* ohne Welle\.$/**Welle:** SPEC-004/' "$w/$P1"
rot nr5-anderes-feld "$w" "$P1: §1: SPEC-004 fehlt im Kopf"
# Ein Feld nach der ersten Zeile „## “ ist kein Kopf.
neu; plan "$w/$P1" "—" "—" "SPEC-004" "-" "$(printf '**Berührte Spec-Stellen:** SPEC-004')"
rot nr5-feld-nach-ueberschrift "$w" "$P1: §1: SPEC-004 fehlt im Kopf"

# --- Nr. 6: §1 und §2 --------------------------------------------------------
einzeln_gruen nr6-ausserhalb "—" "—" "Nichts." "- [ ] nichts." "SPEC-004 in §3"
einzeln_gruen nr6-regel-absatz "—" "—" "Regeln dieser Sektion: SPEC-004
und LH-FA-05 in der Folgezeile." "-"
einzeln_rot nr6-nach-regel-absatz "§1: SPEC-004 fehlt im Kopf" "—" "—" "$(printf 'Regeln dieser Sektion: x.\n\nSPEC-004')" "-"
einzeln_rot nr6-regel-mitten-im-absatz "§1: SPEC-004 fehlt im Kopf" "—" "—" "$(printf 'Vorher.\nRegeln dieser Sektion: SPEC-004')" "-"
einzeln_rot nr6-unterueberschrift "§1: SPEC-004 fehlt im Kopf" "—" "—" "$(printf '### Teil\n\nSPEC-004')" "-"
einzeln_rot nr6-abgrenzung "§1: LH-FA-20 fehlt im Kopf" "—" "—" "- Paralleles Einspielen — Out-of-Scope von LH-FA-20." "-"
# Erkannt an der Nummer, nicht am Titel.
neu; plan "$w/$P1" "—" "—" "SPEC-004" "ARC-004"
sed -i 's/^## 1\. Ziel und Abgrenzung$/## 1. Anderer Titel/; s/^## 2\. Definition of Done$/## 2./' "$w/$P1"
rot nr6-titel "$w" "$P1: §1: SPEC-004 fehlt im Kopf" "$P1: §2: ARC-004 fehlt im Kopf"

# --- Nr. 7: Formfehler -------------------------------------------------------
neu; plan "$w/$P1" "—" "—" "x" "-"; sed -i '/^\*\*Bezug:\*\*/d' "$w/$P1"
rot nr7-ohne-bezug "$w" "$P1: Kopf: Feld Bezug fehlt"
neu; plan "$w/$P1" "—" "—" "x" "-"; sed -i '/^\*\*Berührte Spec-Stellen:\*\*/d' "$w/$P1"
rot nr7-ohne-stellen "$w" "$P1: Kopf: Feld Berührte Spec-Stellen fehlt"
neu; plan "$w/$P1" "—" "—" "x" "-"; sed -i 's/^## 1\. .*/## Ziel/' "$w/$P1"
rot nr7-ohne-1 "$w" "$P1: §1: Abschnitt ## 1. fehlt"
neu; plan "$w/$P1" "—" "—" "x" "-"; sed -i 's/^## 2\. .*/## 20. DoD/' "$w/$P1"
rot nr7-ohne-2 "$w" "$P1: §2: Abschnitt ## 2. fehlt"
neu; : > "$w/$P1"
rot nr7-leer "$w" "$P1: Kopf: Feld Berührte Spec-Stellen fehlt" "$P1: Kopf: Feld Bezug fehlt" \
  "$P1: §1: Abschnitt ## 1. fehlt" "$P1: §2: Abschnitt ## 2. fehlt"

# --- Nr. 1: Gegenstand -------------------------------------------------------
for lc in open next in-progress; do
  neu; plan "$w/$P/$lc/slice-b.md" "—" "—" "SPEC-004" "-"
  rot "nr1-$lc" "$w" "$P/$lc/slice-b.md: §1: SPEC-004 fehlt im Kopf"
done
neu
plan "$w/$P/done/slice-b.md" "—" "—" "SPEC-004" "-"
plan "$w/$P/done/welle-x/slice-c.md" "—" "—" "SPEC-004" "-"
plan "$w/$P/open/welle-x.md" "—" "—" "SPEC-004" "-"
plan "$w/$P/in-progress/roadmap.md" "—" "—" "SPEC-004" "-"
plan "$w/$P/open/slicex-b.md" "—" "—" "SPEC-004" "-"
plan "$w/$P/open/slice-b.txt" "—" "—" "SPEC-004" "-"
plan "$w/$P/slice-flach.md" "—" "—" "SPEC-004" "-"
plan "$w/$P/open/unter/slice-d.md" "—" "—" "SPEC-004" "-"
gruen nr1-nicht-geprueft "$w"

# --- Nr. 8: Ausgabe ----------------------------------------------------------
# Alle Pläne gelesen, sortiert nach Pfad, Abschnitt, Kennung; je Abschnitt einmal.
neu
plan "$w/$P/open/slice-z.md" "—" "—" "SPEC-009 SPEC-002 SPEC-009" "SPEC-009"
plan "$w/$P/in-progress/slice-y.md" "—" "—" "ARC-001" "-"
plan "$w/$P/open/slice-m.md" "—" "—" "LH-FA-01" "-"
sed -i '/^\*\*Bezug:\*\*/d' "$w/$P/open/slice-m.md"
rot nr8-sortiert "$w" \
  "$P/in-progress/slice-y.md: §1: ARC-001 fehlt im Kopf" \
  "$P/open/slice-m.md: Kopf: Feld Bezug fehlt" \
  "$P/open/slice-m.md: §1: LH-FA-01 fehlt im Kopf" \
  "$P/open/slice-z.md: §1: SPEC-002 fehlt im Kopf" \
  "$P/open/slice-z.md: §1: SPEC-009 fehlt im Kopf" \
  "$P/open/slice-z.md: §2: SPEC-009 fehlt im Kopf"
# Ohne Ablage: Exit 2.
mkdir -p "$arbeit/ohne-ablage/docs/plan"
lauf "$arbeit/ohne-ablage"
[ "$code" -eq 2 ] || melde "Fall 'nr8-ohne-ablage' erwartet Exit 2, bekam $code"
# Ablage ohne Lifecycle-Verzeichnisse: nichts zu prüfen.
mkdir -p "$arbeit/leer/$P"
gruen nr8-leere-ablage "$arbeit/leer"
# Aufruf aus einem anderen Verzeichnis ohne Wurzel: Wurzel ist das aktuelle.
neu; plan "$w/$P1" "—" "—" "SPEC-004" "-"
set +e; (cd "$w" && bash "$skript") 2>"$arbeit/err"; code=$?; set -e
[ "$code" -eq 1 ] && grep -qxF "kopf-check: $P1: §1: SPEC-004 fehlt im Kopf" "$arbeit/err" \
  || melde "Fall 'nr8-wurzel-aktuell' erwartet Exit 1 mit Befund, bekam $code"

# --- Nr. 9: Start ohne Stufung -----------------------------------------------
# Das Gate-Ziel hängt an GATE_CHECKS; ein Befund ist Exit 1 (Fälle oben).
# make -q endet ungleich 0, wenn das Ziel nicht aktuell ist; gelesen wird nur -p.
datenbank="$(make -C "$hier/../.." -pq help 2>/dev/null || true)"
if ! printf '%s\n' "$datenbank" | awk '/^GATE_CHECKS :?= / { for (i = 3; i <= NF; i++) if ($i == "kopf-check") t = 1 } END { exit !t }'; then
  melde "Fall 'nr9-gate-checks': kopf-check hängt nicht an GATE_CHECKS"
fi

# --- Akzeptierte Negative (ADR-0032 §Konsequenzen) ---------------------------
einzeln_gruen neg-existenz "LH-ZZ-99" "SPEC-999" "LH-ZZ-99 und SPEC-999" "-"
einzeln_gruen neg-feldwahl "SPEC-001 · ARC-001" "LH-FA-01" "LH-FA-01, SPEC-001, ARC-001" "-"
neu; plan "$w/$P1" "—" "—" "x" "-" "SPEC-004 in §3"
printf '\n## 7. Closure-Notiz\n\nLH-FA-07\n' >> "$w/$P1"
sed -i 's/^\*\*Verantwortlich:\*\* —$/**Verantwortlich:** — ARC-004/' "$w/$P1"
gruen neg-ausserhalb "$w"
# Eine Zeile „## “ im Codeblock beendet den Abschnitt.
einzeln_gruen neg-codeblock-ueberschrift "—" "—" "$(printf '```\n## kein Abschnitt\nSPEC-004\n```')" "-"

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "kopf-check-gegenprobe: gruen — Fehlformen je Nummer von ADR-0032 abgelehnt, vollständiger Kopf und akzeptierte Negative angenommen"
