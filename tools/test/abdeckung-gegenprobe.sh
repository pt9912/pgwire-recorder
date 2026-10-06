#!/usr/bin/env bash
# abdeckung-gegenprobe — prüft tools/test/abdeckung.sh gegen seinen Vertrag
# (spec/spezifikation.md SPEC-048) an kleinen Bäumen in einem Temp-Verzeichnis; der
# Arbeitsbaum bleibt unberührt. Abgelehnt werden müssen, je mit Exit 1 und einer
# Zeile auf stderr, die die Testdatei mit Zeile nennt: eine Leerzeile zwischen
# Deklaration und Test, eine Deklaration über einer Hilfsfunktion, zwei
# Deklarationen über einem Test, eine Deklaration am Dateiende, eine eingerückte
# Deklaration, ein unbekannter Pfad, ein Pfad der falschen Anforderungsart (LH-QA mit
# Happy, LH-FA mit Messung), eine Anforderung, die nicht im Lastenheft steht, ein
# TestE2E ohne Deklaration unter test/integration/ (mit der Zeile des Tests), anderswo
# und unter einem .harness/ tiefer im Baum, und eine Deklaration ohne „ — “; abgelehnt
# wird auch eine Anforderung, die das Lastenheft nur als Überschrift zweiter Ebene
# führt. Ohne Lastenheft Exit 2 mit Meldung, mit einem Lastenheft ohne Überschrift
# ### LH-… Exit 1. Angenommen werden muss ein Baum, dessen einzige Fehlform unter
# .git/ und .harness/ an der Wurzel liegt, und ein gültiger Baum; dessen
# vollständige Anforderungen, das maskierte `|`, die Fortsetzungszeile, die
# Dateirechte 0644 und das Verhalten von --check (prüft, schreibt nicht, nennt jede
# veraltete Tabelle, Exit 1) werden geprüft. Ein Baum prüft das Schreiben: nur
# abweichende oder fehlende Tabellen, je eine Zeile auf stdout, Rechte 0644 nur für
# geschriebene. Ein Baum mit Unit- und E2E-Tests prüft die Tabellenzeilen aller vier
# Tabellen Zeile für Zeile, mit je einer LH-FA, der Happy, Boundary bzw. Negative
# fehlt.
#
# Ausgang: 0, wenn jeder Fall wie erwartet endet, sonst 1.
set -euo pipefail

skript="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/abdeckung.sh"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT
fehler=0

melde() { echo "abdeckung-gegenprobe: ROT — $*" >&2; fehler=1; }

# baum <name> <datei> <inhalt> [<datei> <inhalt>]… — legt einen Baum mit Lastenheft
# und den Testdateien an.
baum() {
  local wurzel="$arbeit/$1"
  shift
  mkdir -p "$wurzel/spec" "$wurzel/docs/user"
  printf '### LH-FA-01 — Eins\n\n### LH-FA-02 — Zwei\n\n### LH-FA-03 — Drei\n\n### LH-FA-04 — Vier\n\n### LH-FA-05 — Fünf\n\n### LH-QA-01 — Qualität\n\n### LH-RB-01 — Rand\n' > "$wurzel/spec/lastenheft.md"
  while [ "$#" -ge 2 ]; do
    mkdir -p "$wurzel/$(dirname "$1")"
    printf '%s\n' "$2" > "$wurzel/$1"
    shift 2
  done
  echo "$wurzel"
}

# abgelehnt <name> <datei> <inhalt> — erwartet Exit 1 und eine Zeile auf stderr,
# die <datei> mit Zeilennummer nennt.
abgelehnt() {
  local name="$1" datei="$2" inhalt="$3" wurzel code
  wurzel="$(baum "$name" "$datei" "$inhalt")"
  set +e
  bash "$skript" "$wurzel" >/dev/null 2>"$arbeit/err"
  code=$?
  set -e
  if [ "$code" -ne 1 ]; then
    melde "Fall '$name' erwartet Exit 1, bekam $code"
  elif ! grep -qE "^abdeckung: $datei:[0-9]+: " "$arbeit/err"; then
    melde "Fall '$name' nennt '$datei:<zeile>' nicht auf stderr: $(cat "$arbeit/err")"
  fi
}

abgelehnt leerzeile internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy — x\n\nfunc TestA(t *testing.T) {}')"
abgelehnt hilfsfunktion internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy — x\nfunc hilfe() {}')"
abgelehnt doppelt internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy — x\n// Abdeckung: LH-FA-01/Boundary — y\nfunc TestA(t *testing.T) {}')"
abgelehnt dateiende internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy — x')"
abgelehnt eingerueckt internal/a_test.go "$(printf 'package a\n\n\t// Abdeckung: LH-FA-01/Happy — x\nfunc TestA(t *testing.T) {}')"
abgelehnt pfad internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Irgendwie — x\nfunc TestA(t *testing.T) {}')"
abgelehnt qa-happy internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-QA-01/Happy — x\nfunc TestA(t *testing.T) {}')"
abgelehnt fa-messung internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Messung — x\nfunc TestA(t *testing.T) {}')"
abgelehnt unbekannt internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-99/Happy — x\nfunc TestA(t *testing.T) {}')"
abgelehnt ohne-deklaration test/integration/a_test.go "$(printf 'package a\n\nfunc TestE2EA(t *testing.T) {}')"
# Bei einem TestE2E ohne Deklaration nennt die Zeile die Zeile des Tests.
grep -qF 'abdeckung: test/integration/a_test.go:3: ' "$arbeit/err" \
  || melde "Fall 'ohne-deklaration-zeile': stderr nennt nicht Zeile 3 des Tests: $(cat "$arbeit/err")"
abgelehnt ohne-deklaration-unit internal/a_test.go "$(printf 'package a\n\nfunc TestE2EA(t *testing.T) {}')"
abgelehnt ohne-trenner internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy x\nfunc TestA(t *testing.T) {}')"
# Nur .git/ und .harness/ an der Wurzel sind ausgenommen; tiefer im Baum wird gelesen.
abgelehnt verschachtelt sub/.harness/b_test.go "$(printf 'package a\n\nfunc TestE2EA(t *testing.T) {}')"

# Eine Anforderung zählt nur als Überschrift dritter Ebene des Lastenhefts.
wurzel="$(baum ueberschrift-ebene internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-09/Happy — x\nfunc TestA(t *testing.T) {}')")"
printf '\n## LH-FA-09 — Neun\n' >> "$wurzel/spec/lastenheft.md"
if bash "$skript" "$wurzel" >/dev/null 2>&1; then
  melde "Fall 'ueberschrift-ebene': Anforderung aus einer Überschrift zweiter Ebene angenommen"
fi

# Testdateien unter .git/ und .harness/ werden nicht gelesen: ein TestE2E ohne
# Deklaration dort bleibt grün.
wurzel="$(baum nicht-gelesen \
  .git/a_test.go "$(printf 'package a\n\nfunc TestE2EA(t *testing.T) {}')" \
  .harness/x/a_test.go "$(printf 'package a\n\nfunc TestE2EA(t *testing.T) {}')")"
bash "$skript" "$wurzel" >/dev/null 2>"$arbeit/err" \
  || melde "Fall 'nicht-gelesen': Testdatei unter .git/ oder .harness/ gelesen: $(cat "$arbeit/err")"

# Fehlt das Lastenheft: Exit 2 mit einer Meldung auf stderr.
wurzel="$(baum ohne-lastenheft internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy — x\nfunc TestA(t *testing.T) {}')")"
rm "$wurzel/spec/lastenheft.md"
set +e
bash "$skript" "$wurzel" >/dev/null 2>"$arbeit/err"
code=$?
set -e
{ [ "$code" -eq 2 ] && [ -s "$arbeit/err" ]; } \
  || melde "Fall 'ohne-lastenheft' erwartet Exit 2 mit Meldung, bekam $code: $(cat "$arbeit/err")"

# Grenze: Führt das Lastenheft keine Überschrift ### LH-…, endet die Prüfung mit 1.
wurzel="$(baum lastenheft-ohne-anforderung internal/a_test.go 'package a')"
printf '# Lastenheft\n' > "$wurzel/spec/lastenheft.md"
set +e
bash "$skript" "$wurzel" >/dev/null 2>/dev/null
code=$?
set -e
[ "$code" -eq 1 ] \
  || melde "Fall 'lastenheft-ohne-anforderung' erwartet Exit 1, bekam $code"

# --- Gültiger Baum: Fortsetzung, Maskierung, Rechte, --check ---------------
gueltig="$(printf '%s\n' \
  'package a' '' \
  '// Abdeckung: LH-FA-01/Happy, LH-FA-01/Boundary — a | b' \
  '// Fortsetzung' \
  'func TestA(t *testing.T) {}' '' \
  '// Abdeckung: LH-FA-01/Negative, LH-QA-01/Messung — c' \
  'func TestB(t *testing.T) {}')"
wurzel="$(baum gueltig internal/a_test.go "$gueltig")"
if ! bash "$skript" "$wurzel" >/dev/null 2>&1; then
  melde "Fall 'gueltig': gültiger Baum abgelehnt"
else
  voll="$wurzel/docs/user/abdeckung-vollstaendig.md"
  grep -q 'LH-FA-01' "$voll" && grep -q 'LH-QA-01' "$voll" \
    || melde "Fall 'gueltig-vollstaendig': vollständige Anforderungen fehlen"
  grep -qF 'a \| b Fortsetzung' "$wurzel/docs/user/abdeckung-unit.md" \
    || melde "Fall 'gueltig-maskierung': | nicht maskiert oder Fortsetzung fehlt"
  [ "$(stat -c %a "$voll")" = 644 ] \
    || melde "Fall 'gueltig-rechte': Dateirechte nicht 0644"
  bash "$skript" --check "$wurzel" >/dev/null 2>&1 \
    || melde "Fall 'check-aktuell': --check auf aktuellem Baum rot"
  echo veraendert > "$voll"
  set +e
  bash "$skript" --check "$wurzel" >/dev/null 2>"$arbeit/err"
  code=$?
  set -e
  [ "$code" -eq 1 ] \
    || melde "Fall 'check-veraltet': --check auf veraltetem Baum erwartet Exit 1, bekam $code"
  grep -qF 'docs/user/abdeckung-vollstaendig.md' "$arbeit/err" \
    || melde "Fall 'check-veraltet': stderr nennt die veraltete Tabelle nicht: $(cat "$arbeit/err")"
  [ "$(cat "$voll")" = veraendert ] \
    || melde "Fall 'check-schreibt-nicht': --check hat geschrieben"
  # Je veralteter Tabelle eine Zeile: zwei veraltete Tabellen, zwei Zeilen.
  echo veraendert > "$wurzel/docs/user/abdeckung-unit.md"
  set +e
  bash "$skript" --check "$wurzel" >/dev/null 2>"$arbeit/err"
  code=$?
  set -e
  { [ "$code" -eq 1 ] && grep -qF 'docs/user/abdeckung-unit.md' "$arbeit/err" \
      && grep -qF 'docs/user/abdeckung-vollstaendig.md' "$arbeit/err"; } \
    || melde "Fall 'check-zwei-veraltet': erwartet Exit 1 und beide Tabellen auf stderr, bekam $code: $(cat "$arbeit/err")"
fi

# --- Schreiben: nur abweichende oder fehlende Tabellen ----------------------
# Der erste Lauf schreibt alle vier und meldet je eine Zeile auf stdout. Danach
# tragen alle vier Rechte 0600 und eine ist verändert: Der zweite Lauf schreibt nur
# sie, setzt ihre Rechte auf 0644 und lässt die übrigen samt Rechten stehen.
wurzel="$(baum schreiben internal/a_test.go "$gueltig")"
U="$wurzel/docs/user"
alle="$(printf 'abdeckung: docs/user/%s geschrieben\n' abdeckung-e2e.md abdeckung-unit.md abdeckung-gesamt.md abdeckung-vollstaendig.md)"
if ! aus="$(bash "$skript" "$wurzel" 2>/dev/null)"; then
  melde "Fall 'schreiben': gültiger Baum abgelehnt"
elif [ "$aus" != "$alle" ]; then
  melde "Fall 'schreiben-fehlend' erwartet auf stdout:"$'\n'"$alle"$'\n'"bekam:"$'\n'"$aus"
else
  cp "$U/abdeckung-gesamt.md" "$arbeit/gesamt.soll"
  chmod 0600 "$U"/abdeckung-*.md
  echo veraendert > "$U/abdeckung-gesamt.md"
  aus="$(bash "$skript" "$wurzel" 2>/dev/null)" || true
  [ "$aus" = "abdeckung: docs/user/abdeckung-gesamt.md geschrieben" ] \
    || melde "Fall 'schreiben-nur-geaendert' erwartet genau die Zeile für abdeckung-gesamt.md, bekam: $aus"
  cmp -s "$U/abdeckung-gesamt.md" "$arbeit/gesamt.soll" \
    || melde "Fall 'schreiben-nur-geaendert': abdeckung-gesamt.md nicht neu geschrieben"
  [ "$(stat -c %a "$U/abdeckung-gesamt.md")" = 644 ] \
    || melde "Fall 'schreiben-rechte': geschriebene Tabelle trägt nicht 0644"
  for t in abdeckung-e2e.md abdeckung-unit.md abdeckung-vollstaendig.md; do
    [ "$(stat -c %a "$U/$t")" = 600 ] \
      || melde "Fall 'schreiben-unberuehrt': aktuelle Tabelle $t in ihren Rechten geändert"
  done
fi

# --- Inhalt der vier Tabellen ----------------------------------------------
# Unit-Tests unter internal/, E2E-Tests unter test/integration/. LH-FA-01 ist
# teilweise belegt (Boundary fehlt), LH-FA-02 vollständig aus zwei Nachweisarten,
# LH-QA-01 und LH-RB-01 je mit Messung; LH-FA-03 steht im Lastenheft ohne
# Deklaration. Verglichen werden die Tabellenzeilen (Zeilen mit `|` am Anfang).
unit="$(printf '%s\n' \
  'package a' '' \
  '// Abdeckung: LH-FA-01/Happy — u1' \
  'func TestU(t *testing.T) {}' '' \
  '// Abdeckung: LH-FA-02/Boundary, LH-FA-02/Negative, LH-RB-01/Messung — u2' \
  'func TestV(t *testing.T) {}' '' \
  '// Abdeckung: LH-FA-04/Boundary, LH-FA-04/Negative, LH-FA-05/Happy, LH-FA-05/Boundary — u3' \
  'func TestW(t *testing.T) {}')"
e2e="$(printf '%s\n' \
  'package e' '' \
  '// Abdeckung: LH-FA-01/Happy, LH-QA-01/Messung — e1' \
  'func TestE2EB(t *testing.T) {}' '' \
  '// Abdeckung: LH-FA-01/Negative, LH-FA-01/Happy, LH-FA-02/Happy — e2' \
  'func TestE2EA(t *testing.T) {}')"
wurzel="$(baum inhalt internal/a_test.go "$unit" test/integration/e_test.go "$e2e")"
L='](../../spec/lastenheft.md)'
# zeilen <fall> <tabelle> <zeile>… — vergleicht die Tabellenzeilen der Datei mit den Zeilen.
zeilen() {
  local fall="$1" tabelle="$2" ist soll
  shift 2
  ist="$(grep '^|' "$wurzel/docs/user/$tabelle" || true)"
  soll="$(printf '%s\n' "$@")"
  [ "$ist" = "$soll" ] \
    || melde "Fall '$fall' in $tabelle erwartet:"$'\n'"$soll"$'\n'"bekam:"$'\n'"$ist"
}
if ! bash "$skript" "$wurzel" >/dev/null 2>&1; then
  melde "Fall 'inhalt': gültiger Baum abgelehnt"
else
  zeilen inhalt-e2e abdeckung-e2e.md \
    '| Anforderung | Pfad | Kurzbeschreibung | Nachweis | Ort |' \
    '| --- | --- | --- | --- | --- |' \
    "| [\`LH-FA-01\`$L | Happy | e2 | \`TestE2EA\` | \`test/integration/e_test.go\` |" \
    "| [\`LH-FA-01\`$L | Happy | e1 | \`TestE2EB\` | \`test/integration/e_test.go\` |" \
    "| [\`LH-FA-01\`$L | Negative | e2 | \`TestE2EA\` | \`test/integration/e_test.go\` |" \
    "| [\`LH-FA-02\`$L | Happy | e2 | \`TestE2EA\` | \`test/integration/e_test.go\` |" \
    "| [\`LH-QA-01\`$L | Messung | e1 | \`TestE2EB\` | \`test/integration/e_test.go\` |"
  zeilen inhalt-unit abdeckung-unit.md \
    '| Anforderung | Pfad | Kurzbeschreibung | Nachweis | Ort |' \
    '| --- | --- | --- | --- | --- |' \
    "| [\`LH-FA-01\`$L | Happy | u1 | \`TestU\` | \`internal/a_test.go\` |" \
    "| [\`LH-FA-02\`$L | Boundary | u2 | \`TestV\` | \`internal/a_test.go\` |" \
    "| [\`LH-FA-02\`$L | Negative | u2 | \`TestV\` | \`internal/a_test.go\` |" \
    "| [\`LH-FA-04\`$L | Boundary | u3 | \`TestW\` | \`internal/a_test.go\` |" \
    "| [\`LH-FA-04\`$L | Negative | u3 | \`TestW\` | \`internal/a_test.go\` |" \
    "| [\`LH-FA-05\`$L | Boundary | u3 | \`TestW\` | \`internal/a_test.go\` |" \
    "| [\`LH-FA-05\`$L | Happy | u3 | \`TestW\` | \`internal/a_test.go\` |" \
    "| [\`LH-RB-01\`$L | Messung | u2 | \`TestV\` | \`internal/a_test.go\` |"
  zeilen inhalt-gesamt abdeckung-gesamt.md \
    '| Anforderung | Happy | Boundary | Negative | Messung | Stand |' \
    '| --- | --- | --- | --- | --- | --- |' \
    "| [\`LH-FA-01\`$L | E2E, Unit | — | E2E | n/a | teilweise |" \
    "| [\`LH-FA-02\`$L | E2E | Unit | Unit | n/a | vollständig |" \
    "| [\`LH-FA-04\`$L | — | Unit | Unit | n/a | teilweise |" \
    "| [\`LH-FA-05\`$L | Unit | Unit | — | n/a | teilweise |" \
    "| [\`LH-QA-01\`$L | n/a | n/a | n/a | E2E | vollständig |" \
    "| [\`LH-RB-01\`$L | n/a | n/a | n/a | Unit | vollständig |"
  zeilen inhalt-vollstaendig abdeckung-vollstaendig.md \
    '| Anforderung | Nachweisarten |' \
    '| --- | --- |' \
    "| [\`LH-FA-02\`$L | E2E, Unit |" \
    "| [\`LH-QA-01\`$L | E2E |" \
    "| [\`LH-RB-01\`$L | Unit |"
fi

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "abdeckung-gegenprobe: gruen — Fehlformen abgelehnt, gültige Bäume mit dem erwarteten Inhalt der vier Tabellen, --check schreibt nicht"
