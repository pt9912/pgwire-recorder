#!/usr/bin/env bash
# abdeckung-gegenprobe — prüft tools/test/abdeckung.sh gegen seinen Vertrag
# (harness/vertraege/SPEC-048-abdeckung.md) an kleinen Bäumen in einem Temp-Verzeichnis; der
# Arbeitsbaum bleibt unberührt. Abgelehnt werden müssen, je mit Exit 1 und einer
# Zeile auf stderr, die die Datei mit Zeile nennt: eine Leerzeile zwischen
# Deklaration und Test, eine Deklaration über einer Hilfsfunktion, zwei
# Deklarationen über einem Test, eine Deklaration am Dateiende, eine eingerückte
# Deklaration, ein unbekannter Pfad, ein Pfad der falschen Anforderungsart (LH-QA mit
# Happy, LH-FA mit Messung), eine Anforderung, die nicht im Lastenheft steht, ein
# TestE2E ohne Deklaration unter test/integration/ (mit der Zeile des Tests), anderswo
# und unter einem .harness/ tiefer im Baum, und eine Deklaration ohne „ — “; abgelehnt
# wird auch eine Anforderung, die das Lastenheft nur als Überschrift zweiter Ebene
# führt. Ohne Lastenheft, mit einem Lastenheft ohne Überschrift ### LH-… und ohne
# docs/user/ Exit 2 mit einer Meldung, die nennt, was fehlt; in dieser Reihenfolge, ohne
# etwas anzulegen und auch bei einer Fehlform im Baum. Angenommen werden muss ein Baum,
# dessen einzige Fehlform unter .git/ und .harness/ an der Wurzel liegt, und ein
# gültiger Baum; dessen vollständige Anforderungen, das maskierte `|`, die
# Fortsetzungszeile, die Dateirechte 0644 und das Verhalten von --check (prüft, schreibt
# nicht, nennt jede veraltete Tabelle, Exit 1) werden geprüft. Ein Baum prüft das
# Schreiben: nur abweichende oder fehlende Tabellen, je eine Zeile auf stdout, Rechte
# 0644 nur für geschriebene. Ein Baum mit Unit- und E2E-Tests prüft die Tabellenzeilen
# aller Tabellen Zeile für Zeile, mit je einer LH-FA, der Happy, Boundary bzw. Negative
# fehlt.
#
# Gegenproben (Punkt 7) — Fälle gate-*: gelesen wird der Kopf der Dateien auf
# `-gegenprobe.sh` unter tools/ an der Wurzel in jeder Tiefe, nicht tools/x.sh, nicht
# `gegenprobe.sh`, nicht `sub/tools/` und nicht außerhalb von tools/; die Form des Kopfes
# je Fehlform und je Gutform ein Fall (eingerückt, ohne Trenner, nach dem Kopf, nach einer
# Leerzeile hinter der Shebang-Zeile, Fortsetzung, Ende an `#` allein, an `#` ohne
# Leerzeichen, an der nächsten Deklaration, am Kopfende, mehrere Deklarationen), LH-FA
# an einer Gate-Deklaration. Eine Gegenprobe ohne Deklaration erzeugt weder Fehler noch Zeile.
#
# Geteilte Messung (Punkt 8) — Fälle teil-*, widerspruch-*, grund-*: je Zahlform an i und
# n ein Fall, i > n, n = 1, verschiedene n, `Messung` neben einem Teil in beiden
# Reihenfolgen, an Tests und Gegenproben; die Reihenfolge der Gründe, soweit die Klasse
# der Meldung sie unterscheidet, der Widerspruch an der späteren Fundstelle nach Dateipfad (bytegenau)
# und Zeile; die Darstellung in der Gesamtsicht und die Vollständigkeit.
#
# Gelesen werden nur reguläre Dateien (ein Symlink auf eine Gegenprobe nicht). Ein Widerspruch
# nennt die erste Deklaration, die das n festlegt, bzw. die erste des anderen Pfads; ein
# widersprechendes Element legt nichts fest (widerspruch-erste-*, widerspruch-dritter). Eine
# eingerückte Deklaration ist an jeder Stelle der Datei eine Fehlform, mit der Zeile selbst; ein
# Lastenheft nur mit `### LH-XX-01` besteht den Ausgang 2 und scheitert an (4) mit Exit 1.
#
# Nicht gezeigt, weil über die Schnittstelle nicht unterscheidbar: die Reihenfolge von i > n und
# n = 1 untereinander (der Vertrag lässt offen, welcher Grund genannt wird); das Ende der
# Fortsetzung am Kopfende (die erste Zeile ohne `#` beendet sie selbst); der Zweig für `LH-FA` in
# der Gate-Deklaration allein (das Muster der Pfade lehnt sie mit einer anderen Meldung ab).
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

# rot <name> <wurzel> — erwartet Exit 1; stderr liegt danach in $arbeit/err.
rot() {
  local code
  set +e
  bash "$skript" "$2" >/dev/null 2>"$arbeit/err"
  code=$?
  set -e
  [ "$code" -eq 1 ] || melde "Fall '$1' erwartet Exit 1, bekam $code: $(cat "$arbeit/err")"
}

# angenommen <name> <wurzel> — erwartet Exit 0 und nichts auf stderr.
angenommen() {
  local code
  set +e
  bash "$skript" "$2" >/dev/null 2>"$arbeit/err"
  code=$?
  set -e
  { [ "$code" -eq 0 ] && [ ! -s "$arbeit/err" ]; } \
    || melde "Fall '$1' erwartet Exit 0 ohne Zeile auf stderr, bekam $code: $(cat "$arbeit/err")"
}

# stderr_klasse <name> <anzahl> <muster> — stderr hat genau <anzahl> Zeilen, jede passt auf
# das Muster; die Klasse des Grundes erkennt der Fall an der Meldung.
stderr_klasse() {
  local name="$1" anzahl="$2" muster="$3" ist
  ist="$(wc -l < "$arbeit/err")"
  [ "$ist" -eq "$anzahl" ] || melde "Fall '$name' erwartet $anzahl Zeile(n) auf stderr, bekam $ist: $(cat "$arbeit/err")"
  if [ -s "$arbeit/err" ] && grep -qvE -- "$muster" "$arbeit/err"; then
    melde "Fall '$name' erwartet auf jeder Zeile von stderr das Muster '$muster': $(cat "$arbeit/err")"
  fi
}

# stderr_ist <name> <zeile>… — stderr besteht genau aus diesen Zeilen.
stderr_ist() {
  local name="$1" soll
  shift
  soll="$(printf '%s\n' "$@")"
  [ "$(cat "$arbeit/err")" = "$soll" ] \
    || melde "Fall '$name' erwartet auf stderr:"$'\n'"$soll"$'\n'"bekam:"$'\n'"$(cat "$arbeit/err")"
}

# tabelle <name> <wurzel> <tabelle> <zeile>… — vergleicht die Tabellenzeilen der Datei.
tabelle() {
  local fall="$1" w="$2" t="$3" ist soll
  shift 3
  ist="$(grep '^|' "$w/docs/user/$t" || true)"
  soll="$(printf '%s\n' "$@")"
  [ "$ist" = "$soll" ] \
    || melde "Fall '$fall' in $t erwartet:"$'\n'"$soll"$'\n'"bekam:"$'\n'"$ist"
}

# t <deklaration> — Inhalt einer Testdatei: Deklaration in Zeile 3, Test in Zeile 4.
t() { printf 'package a\n\n// Abdeckung: %s — x\nfunc TestA(t *testing.T) {}' "$1"; }

# gate <zeile>… — Inhalt einer Gegenprobe: Shebang-Zeile, die Zeilen, danach `set -e`.
gate() { printf '%s\n' '#!/usr/bin/env bash' "$@" 'set -e'; }
D='# Abdeckung: '

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
{ [ "$code" -eq 2 ] && grep -qF 'spec/lastenheft.md' "$arbeit/err"; } \
  || melde "Fall 'ohne-lastenheft' erwartet Exit 2 mit Meldung zum Lastenheft, bekam $code: $(cat "$arbeit/err")"

# Führt das Lastenheft keine Überschrift ### LH-…, endet die Prüfung mit 2 und einer Meldung,
# die das Lastenheft nennt; das ist der zweite Grund des Ausgangs 2.
wurzel="$(baum lastenheft-ohne-anforderung internal/a_test.go 'package a')"
printf '# Lastenheft\n' > "$wurzel/spec/lastenheft.md"
set +e
bash "$skript" "$wurzel" >/dev/null 2>"$arbeit/err"
code=$?
set -e
{ [ "$code" -eq 2 ] && grep -qF 'spec/lastenheft.md' "$arbeit/err"; } \
  || melde "Fall 'lastenheft-ohne-anforderung' erwartet Exit 2 mit Meldung zum Lastenheft, bekam $code: $(cat "$arbeit/err")"

# Fehlt docs/user/: Exit 2 mit Meldung, die das Verzeichnis nennt; nichts geprüft (auch eine
# Fehlform im Baum bleibt ohne Exit 1) und nichts angelegt, mit und ohne --check.
wurzel="$(baum ohne-docs-user internal/a_test.go "$(t 'LH-FA-99/Happy')")"
rm -r "$wurzel/docs"
for m in "" --check; do
  set +e
  bash "$skript" ${m:+"$m"} "$wurzel" >/dev/null 2>"$arbeit/err"
  code=$?
  set -e
  { [ "$code" -eq 2 ] && grep -qF 'docs/user' "$arbeit/err" && [ "$(wc -l < "$arbeit/err")" -eq 1 ]; } \
    || melde "Fall 'ohne-docs-user${m:+ $m}' erwartet Exit 2 mit einer Meldung zu docs/user, bekam $code: $(cat "$arbeit/err")"
  [ ! -e "$wurzel/docs" ] || melde "Fall 'ohne-docs-user${m:+ $m}': docs/ angelegt"
done

# Reihenfolge der Gründe des Ausgangs 2: Lastenheft vor Überschrift vor docs/user/.
wurzel="$(baum ausgang2-lastenheft-zuerst internal/a_test.go 'package a')"
rm -r "$wurzel/spec" "$wurzel/docs"
set +e
bash "$skript" "$wurzel" >/dev/null 2>"$arbeit/err"
code=$?
set -e
{ [ "$code" -eq 2 ] && grep -qF 'fehlt: spec/lastenheft.md' "$arbeit/err"; } \
  || melde "Fall 'ausgang2-lastenheft-zuerst' erwartet Exit 2 mit der Meldung zum fehlenden Lastenheft, bekam $code: $(cat "$arbeit/err")"
wurzel="$(baum ausgang2-ueberschrift-vor-docs internal/a_test.go 'package a')"
printf '# Lastenheft\n' > "$wurzel/spec/lastenheft.md"
rm -r "$wurzel/docs"
set +e
bash "$skript" "$wurzel" >/dev/null 2>"$arbeit/err"
code=$?
set -e
{ [ "$code" -eq 2 ] && grep -qF 'Überschrift' "$arbeit/err" && ! grep -qF 'docs/user' "$arbeit/err"; } \
  || melde "Fall 'ausgang2-ueberschrift-vor-docs' erwartet Exit 2 mit der Meldung zur Überschrift, bekam $code: $(cat "$arbeit/err")"

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
# Der erste Lauf schreibt alle fünf und meldet je eine Zeile auf stdout. Danach
# tragen alle fünf Rechte 0600 und eine ist verändert: Der zweite Lauf schreibt nur
# sie, setzt ihre Rechte auf 0644 und lässt die übrigen samt Rechten stehen.
wurzel="$(baum schreiben internal/a_test.go "$gueltig")"
U="$wurzel/docs/user"
alle="$(printf 'abdeckung: docs/user/%s geschrieben\n' abdeckung-e2e.md abdeckung-unit.md abdeckung-gate.md abdeckung-gesamt.md abdeckung-vollstaendig.md | sort)"
if ! aus="$(bash "$skript" "$wurzel" 2>/dev/null | sort)"; then
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
  for t in abdeckung-e2e.md abdeckung-unit.md abdeckung-gate.md abdeckung-vollstaendig.md; do
    [ "$(stat -c %a "$U/$t")" = 600 ] \
      || melde "Fall 'schreiben-unberuehrt': aktuelle Tabelle $t in ihren Rechten geändert"
  done
fi

# --- Inhalt der Tabellen ----------------------------------------------
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
  zeilen inhalt-gate-leer abdeckung-gate.md \
    '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' \
    '| --- | --- | --- | --- |'
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

# --- Nachweisart nach Ort: ein E2E-Verzeichnis tiefer im Baum ergibt Unit ----
wurzel="$(baum e2e-tiefer sub/test/integration/x_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy — x\nfunc TestE2EX(t *testing.T) {}')")"
if bash "$skript" "$wurzel" >/dev/null 2>&1; then
  tabelle e2e-tiefer-e2e "$wurzel" abdeckung-e2e.md \
    '| Anforderung | Pfad | Kurzbeschreibung | Nachweis | Ort |' '| --- | --- | --- | --- | --- |'
  tabelle e2e-tiefer-unit "$wurzel" abdeckung-unit.md \
    '| Anforderung | Pfad | Kurzbeschreibung | Nachweis | Ort |' '| --- | --- | --- | --- | --- |' \
    "| [\`LH-FA-01\`$L | Happy | x | \`TestE2EX\` | \`sub/test/integration/x_test.go\` |"
else
  melde "Fall 'e2e-tiefer': gültiger Baum abgelehnt"
fi

# --- Gegenproben: gelesen werden Dateien auf -gegenprobe.sh unter tools/ --------
# Gutform: Zeile 2 des Kopfs, ein Eintrag in der Tabelle der Nachweisart Gate, keiner in den
# Tabellen der Tests, ein `|` im Text maskiert.
wurzel="$(baum gate-gutform tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung — a | b")")"
angenommen gate-gutform "$wurzel"
tabelle gate-gutform "$wurzel" abdeckung-gate.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' '| --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | Messung | a \\| b | \`tools/x-gegenprobe.sh\` |"
tabelle gate-gutform-unit "$wurzel" abdeckung-unit.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis | Ort |' '| --- | --- | --- | --- | --- |'
tabelle gate-gutform-gesamt "$wurzel" abdeckung-gesamt.md \
  '| Anforderung | Happy | Boundary | Negative | Messung | Stand |' '| --- | --- | --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | n/a | n/a | n/a | Gate | vollständig |"
tabelle gate-gutform-vollstaendig "$wurzel" abdeckung-vollstaendig.md \
  '| Anforderung | Nachweisarten |' '| --- | --- |' "| [\`LH-QA-01\`$L | Gate |"

# Randbedingung ebenso; tiefer unter tools/ und ohne Shebang-Zeile gelesen.
wurzel="$(baum gate-tiefer tools/a/b/c-gegenprobe.sh "$(printf '%s\n' "${D}LH-RB-01/Messung — r" 'set -e')")"
angenommen gate-tiefer "$wurzel"
tabelle gate-tiefer "$wurzel" abdeckung-gate.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' '| --- | --- | --- | --- |' \
  "| [\`LH-RB-01\`$L | Messung | r | \`tools/a/b/c-gegenprobe.sh\` |"

# Nicht gelesen, auch mit einer Fehlform in der Deklaration: tools/x.sh, `gegenprobe.sh` ohne
# Bindestrich, ein anderer Endung, eine Gegenprobe außerhalb von tools/ und unter einem
# tools/ tiefer im Baum.
for f in tools/x.sh tools/gegenprobe.sh tools/x-gegenprobe.sh.bak harness/x-gegenprobe.sh sub/tools/x-gegenprobe.sh; do
  wurzel="$(baum "gate-nicht-gelesen-$(basename "$(dirname "$f")")-$(basename "$f")" "$f" "$(gate "${D}LH-FA-99/Happy ohne Trenner")")"
  angenommen "gate-nicht-gelesen $f" "$wurzel"
  tabelle "gate-nicht-gelesen $f" "$wurzel" abdeckung-gate.md \
    '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' '| --- | --- | --- | --- |'
done

# Ein Symlink auf eine Datei mit gültiger Deklaration wird nicht gelesen (nur reguläre Dateien).
wurzel="$(baum gate-nicht-gelesen-symlink tools/ziel.sh "$(gate "${D}LH-QA-01/Messung — x")")"
ln -s ziel.sh "$wurzel/tools/x-gegenprobe.sh"
angenommen gate-nicht-gelesen-symlink "$wurzel"
tabelle gate-nicht-gelesen-symlink "$wurzel" abdeckung-gate.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' '| --- | --- | --- | --- |'

# Eine Gegenprobe ohne Deklaration im Muster erzeugt weder Fehler noch Zeile (Bindung aus
# slice-harness-integration-volumes): Kopf mit Fließtext, ohne Zeile `Abdeckung:`.
wurzel="$(baum gegenprobe-ohne-deklaration tools/test/x-gegenprobe.sh "$(gate '# x-gegenprobe — prüft y.' '#' '# GRENZE: Abdeckung wird hier nicht deklariert.' '#' '# Ausgang: 0, wenn jeder Fall wie erwartet endet, sonst 1.')" internal/a_test.go "$(t LH-FA-01/Happy)")"
angenommen gegenprobe-ohne-deklaration "$wurzel"
tabelle gegenprobe-ohne-deklaration "$wurzel" abdeckung-gate.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' '| --- | --- | --- | --- |'
tabelle gegenprobe-ohne-deklaration-gesamt "$wurzel" abdeckung-gesamt.md \
  '| Anforderung | Happy | Boundary | Negative | Messung | Stand |' '| --- | --- | --- | --- | --- | --- |' \
  "| [\`LH-FA-01\`$L | Unit | — | — | n/a | teilweise |"

# Form des Kopfes: mehrere Deklarationen; Fortsetzung mit `# ` und Text; Ende an `#` allein,
# an `#` ohne Leerzeichen, an der nächsten Deklaration und am Kopfende.
wurzel="$(baum gate-kopf tools/x-gegenprobe.sh "$(gate \
  "${D}LH-QA-01/Messung-1-von-3 — eins" '# zwei' '# ' '# frei' \
  "${D}LH-QA-01/Messung-2-von-3 — a" '# b' '#c' \
  "${D}LH-QA-01/Messung-3-von-3 — d" '# e' "${D}LH-RB-01/Messung — f" '# g' \
  '#' '# frei' )")"
printf '%s\n' '# nach dem Kopf' >> "$wurzel/tools/x-gegenprobe.sh"
angenommen gate-kopf "$wurzel"
tabelle gate-kopf "$wurzel" abdeckung-gate.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' '| --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | Messung-1-von-3 | eins zwei | \`tools/x-gegenprobe.sh\` |" \
  "| [\`LH-QA-01\`$L | Messung-2-von-3 | a b | \`tools/x-gegenprobe.sh\` |" \
  "| [\`LH-QA-01\`$L | Messung-3-von-3 | d e | \`tools/x-gegenprobe.sh\` |" \
  "| [\`LH-RB-01\`$L | Messung | f g | \`tools/x-gegenprobe.sh\` |"
# Am Kopfende endet die Fortsetzung: Die Zeile nach dem Kopf wird nicht angefügt.
wurzel="$(baum gate-kopfende tools/x-gegenprobe.sh "$(printf '%s\n' '#!/usr/bin/env bash' "${D}LH-QA-01/Messung — eins" '# zwei' 'set -e' '# drei')")"
angenommen gate-kopfende "$wurzel"
tabelle gate-kopfende "$wurzel" abdeckung-gate.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' '| --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | Messung | eins zwei | \`tools/x-gegenprobe.sh\` |"

# Fehlformen des Kopfs, je ein Fall; die Meldung nennt Datei und Zeile.
abgelehnt gate-eingerueckt tools/x-gegenprobe.sh "$(gate " ${D}LH-QA-01/Messung — x")"
abgelehnt gate-eingerueckt-tab tools/x-gegenprobe.sh "$(gate $'\t'"${D}LH-QA-01/Messung — x")"
abgelehnt gate-ohne-trenner tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung x")"
abgelehnt gate-nach-dem-kopf tools/x-gegenprobe.sh "$(printf '%s\n' '#!/usr/bin/env bash' 'set -e' "${D}LH-QA-01/Messung — x")"
abgelehnt gate-nach-leerzeile tools/x-gegenprobe.sh "$(printf '%s\n' '#!/usr/bin/env bash' '' "${D}LH-QA-01/Messung — x")"
abgelehnt gate-fa tools/x-gegenprobe.sh "$(gate "${D}LH-FA-01/Happy — x")"
abgelehnt gate-fa-messung tools/x-gegenprobe.sh "$(gate "${D}LH-FA-01/Messung — x")"
abgelehnt gate-fa-teil tools/x-gegenprobe.sh "$(gate "${D}LH-FA-01/Messung-1-von-2 — x")"
abgelehnt gate-pfad-happy tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Happy — x")"
abgelehnt gate-unbekannt tools/x-gegenprobe.sh "$(gate "${D}LH-QA-99/Messung — x")"
abgelehnt gate-teil-schreibweise tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung-01-von-2 — x")"
abgelehnt gate-teil-bereich tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung-3-von-2 — x")"
# Die Zeile der Meldung ist die der Deklaration (Zeile 3 im zweiten Eintrag).
abgelehnt gate-zeile tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung — x" "${D}LH-QA-01/Happy — y")"
grep -qF 'abdeckung: tools/x-gegenprobe.sh:3: ' "$arbeit/err" \
  || melde "Fall 'gate-zeile': stderr nennt nicht Zeile 3 der Deklaration: $(cat "$arbeit/err")"
# Eingerückt hinter dem Kopf (Leerzeichen, Tab) ist ebenfalls eine Fehlform; die Meldung nennt die Zeile selbst.
abgelehnt gate-eingerueckt-nach-kopf tools/x-gegenprobe.sh "$(gate 'echo x' " ${D}LH-QA-01/Messung — x")"
grep -qF 'abdeckung: tools/x-gegenprobe.sh:3: ' "$arbeit/err" \
  || melde "Fall 'gate-eingerueckt-nach-kopf': stderr nennt nicht Zeile 3: $(cat "$arbeit/err")"
abgelehnt gate-eingerueckt-nach-kopf-tab tools/x-gegenprobe.sh "$(gate 'echo x' $'\t'"${D}LH-QA-01/Messung — x")"
grep -qF 'abdeckung: tools/x-gegenprobe.sh:3: ' "$arbeit/err" \
  || melde "Fall 'gate-eingerueckt-nach-kopf-tab': stderr nennt nicht Zeile 3: $(cat "$arbeit/err")"
# Die eingerückte Deklaration im Kopf nennt ihre Zeile (2), die nach dem Kopf ebenso (3).
abgelehnt gate-eingerueckt-zeile tools/x-gegenprobe.sh "$(gate " ${D}LH-QA-01/Messung — x")"
grep -qF 'abdeckung: tools/x-gegenprobe.sh:2: ' "$arbeit/err" \
  || melde "Fall 'gate-eingerueckt-zeile': stderr nennt nicht Zeile 2: $(cat "$arbeit/err")"
abgelehnt gate-nach-kopf-zeile tools/x-gegenprobe.sh "$(printf '%s\n' '#!/usr/bin/env bash' 'set -e' "${D}LH-QA-01/Messung — x")"
grep -qF 'abdeckung: tools/x-gegenprobe.sh:3: ' "$arbeit/err" \
  || melde "Fall 'gate-nach-kopf-zeile': stderr nennt nicht Zeile 3: $(cat "$arbeit/err")"
# Eine mehrzeilige Deklaration nennt ihren Anfang (Zeile 2), nicht die Fortsetzung.
abgelehnt gate-zeile-mehrzeilig tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Happy — x" '# Fortsetzung' '# noch eine')"
grep -qF 'abdeckung: tools/x-gegenprobe.sh:2: ' "$arbeit/err" \
  || melde "Fall 'gate-zeile-mehrzeilig': stderr nennt nicht Zeile 2: $(cat "$arbeit/err")"
# Eine Zeile `Abdeckung:` am Zeilenanfang in einer Gegenprobe außerhalb des Kopfs ist eine
# Fehlform, auch mitten in einem Fixture (dort entstehen sie mit printf in einer Zeile).
abgelehnt gate-nach-dem-kopf-code tools/x-gegenprobe.sh "$(printf '%s\n' '#!/usr/bin/env bash' 'cat <<EOF' "${D}LH-QA-01/Messung — x" 'EOF')"

# Gate-Fehlform in einer Datei und gültige Deklaration in einer anderen: beide Dateien gelesen.
wurzel="$(baum gate-zwei-dateien tools/a-gegenprobe.sh "$(gate "${D}LH-FA-01/Happy — x")" tools/b-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung — y")")"
rot gate-zwei-dateien "$wurzel"
stderr_klasse gate-zwei-dateien 1 '^abdeckung: tools/a-gegenprobe.sh:2: '

# --- Geteilte Messung: Schreibweise und Zahlform (Punkt 2, 8) -----------------
# Je Zahlform ein Fall, an i und an n; die Meldung nennt die Schreibweise.
for zf in 'Messung-10-von-12' 'Messung-10-von-2' 'Messung-1-von-10' 'Messung-01-von-2' 'Messung-1-von-02' \
          'Messung-+1-von-2' 'Messung-1-von-+2' 'Messung- 1-von-2' 'Messung-1-von- 2' \
          'Messung-0-von-2' 'Messung-1-von-0' 'Messung-1von-2' 'Messung-1-von' 'Messung-' 'Messung-a-von-2'; do
  fall="teil-zahlform-$zf"
  abgelehnt "$fall" internal/a_test.go "$(t "LH-QA-01/$zf")"
  stderr_klasse "$fall" 1 ' ist nicht '
  # Dieselbe Schreibweise im Kopf einer Gegenprobe.
  abgelehnt "gate-$fall" tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/$zf — x")"
  stderr_klasse "gate-$fall" 1 ' ist nicht '
done
# Dieselbe Schreibweise als Test an einer LH-FA und als Gate-Deklaration.
abgelehnt teil-fa internal/a_test.go "$(t 'LH-FA-01/Messung-1-von-2')"
stderr_klasse teil-fa 1 ' ist nicht '

# --- Wertebereich (Punkt 8): i > n, n = 1; an Tests und Gegenproben ---------------
abgelehnt teil-i-ueber-n internal/a_test.go "$(t 'LH-QA-01/Messung-3-von-2')"
stderr_klasse teil-i-ueber-n 1 'Gesamtzahl'
abgelehnt teil-n-eins internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-1')"
stderr_klasse teil-n-eins 1 'mindestens zwei'
abgelehnt teil-i-ueber-n-rb internal/a_test.go "$(t 'LH-RB-01/Messung-9-von-8')"
stderr_klasse teil-i-ueber-n-rb 1 'Gesamtzahl'
abgelehnt gate-teil-i-ueber-n tools/x-gegenprobe.sh "$(gate "${D}LH-RB-01/Messung-9-von-8 — x")"
stderr_klasse gate-teil-i-ueber-n 1 'Gesamtzahl'
abgelehnt gate-teil-n-eins tools/x-gegenprobe.sh "$(gate "${D}LH-RB-01/Messung-1-von-1 — x")"
stderr_klasse gate-teil-n-eins 1 'mindestens zwei'

# --- Reihenfolge der Gründe je Paar (Punkt 6) ----------------------------------
# (1) Form vor allem: Eine Deklaration ohne „ — “ nennt nur die Form, nicht die Elemente.
abgelehnt grund-1-vor-2-bis-4 internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-QA-99/Messung-01-von-2 x\nfunc TestA(t *testing.T) {}')"
stderr_klasse grund-1-vor-2-bis-4 1 'ohne „ — “'
# (2) Schreibweise vor Wertebereich, (2) vor (4), (3) vor (4).
abgelehnt grund-2-vor-3 internal/a_test.go "$(t 'LH-QA-01/Messung-0-von-1')"
stderr_klasse grund-2-vor-3 1 ' ist nicht '
abgelehnt grund-2-vor-4 internal/a_test.go "$(t 'LH-QA-99/Messung-01-von-2')"
stderr_klasse grund-2-vor-4 1 ' ist nicht '
abgelehnt grund-3-vor-4 internal/a_test.go "$(t 'LH-QA-99/Messung-3-von-2')"
stderr_klasse grund-3-vor-4 1 'Gesamtzahl'
abgelehnt grund-4 internal/a_test.go "$(t 'LH-QA-99/Messung-1-von-2')"
stderr_klasse grund-4 1 'steht nicht im Lastenheft'
# Ein Element nennt einen Grund: zwei Elemente mit je einem Grund, zwei Zeilen.
abgelehnt grund-je-element internal/a_test.go "$(t 'LH-QA-99/Messung-01-von-2, LH-QA-01/Messung-3-von-2')"
stderr_klasse grund-je-element 2 '^abdeckung: internal/a_test.go:4: '

# (5) betrachtet nur Elemente ohne einen Grund aus (1) bis (4): Jeder dieser Elemente steht
# neben einem gültigen mit anderem n und erzeugt keinen Widerspruch.
wurzel="$(baum grund-5-ohne-4 internal/a_test.go "$(t 'LH-QA-99/Messung-1-von-2')" internal/b_test.go "$(t 'LH-QA-99/Messung-1-von-3')")"
rot grund-5-ohne-4 "$wurzel"
stderr_klasse grund-5-ohne-4 2 'steht nicht im Lastenheft'
wurzel="$(baum grund-5-ohne-3 internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-2')" internal/b_test.go "$(t 'LH-QA-01/Messung-1-von-1')")"
rot grund-5-ohne-3 "$wurzel"
stderr_klasse grund-5-ohne-3 1 'mindestens zwei'
wurzel="$(baum grund-5-ohne-2 internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-2')" internal/b_test.go "$(t 'LH-QA-01/Messung-01-von-3')")"
rot grund-5-ohne-2 "$wurzel"
stderr_klasse grund-5-ohne-2 1 ' ist nicht '
wurzel="$(baum grund-5-ohne-1 internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-2')" internal/b_test.go "$(printf 'package a\n\n// Abdeckung: LH-QA-01/Messung-1-von-3 x\nfunc TestA(t *testing.T) {}')")"
rot grund-5-ohne-1 "$wurzel"
stderr_klasse grund-5-ohne-1 1 'ohne „ — “'

# --- Widerspruch zwischen Deklarationen (Punkt 8) --------------------------------
# Verschiedene n: Der Fehler steht an der späteren Fundstelle und nennt die frühere.
wurzel="$(baum widerspruch-n internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-2')" internal/b_test.go "$(t 'LH-QA-01/Messung-2-von-3')")"
rot widerspruch-n "$wurzel"
stderr_ist widerspruch-n 'abdeckung: internal/b_test.go:4: LH-QA-01/Messung-2-von-3 widerspricht LH-QA-01/Messung-1-von-2 (früher: internal/a_test.go:4)'
# Die Reihenfolge ist die der Dateipfade, bytegenau (B vor a), nicht die der Funde.
wurzel="$(baum widerspruch-bytegenau internal/B_test.go "$(t 'LH-QA-01/Messung-1-von-2')" internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-3')")"
rot widerspruch-bytegenau "$wurzel"
stderr_ist widerspruch-bytegenau 'abdeckung: internal/a_test.go:4: LH-QA-01/Messung-1-von-3 widerspricht LH-QA-01/Messung-1-von-2 (früher: internal/B_test.go:4)'
# Dateipfad vor Zeile, die Zeile numerisch (4 vor 12).
spaet="$(printf 'package a\n\n// Abdeckung: LH-QA-01/Messung-1-von-2 — x\nfunc TestA(t *testing.T) {}\n\n\n\n\n\n\n// Abdeckung: LH-QA-01/Messung-1-von-3 — y\nfunc TestB(t *testing.T) {}')"
wurzel="$(baum widerspruch-zeile internal/a_test.go "$spaet")"
rot widerspruch-zeile "$wurzel"
stderr_ist widerspruch-zeile 'abdeckung: internal/a_test.go:12: LH-QA-01/Messung-1-von-3 widerspricht LH-QA-01/Messung-1-von-2 (früher: internal/a_test.go:4)'
# Die Elemente einer Deklaration in ihrer Folge.
wurzel="$(baum widerspruch-element internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-2, LH-QA-01/Messung-1-von-3')")"
rot widerspruch-element "$wurzel"
stderr_ist widerspruch-element 'abdeckung: internal/a_test.go:4: LH-QA-01/Messung-1-von-3 widerspricht LH-QA-01/Messung-1-von-2 (früher: internal/a_test.go:4)'
wurzel="$(baum widerspruch-element-rueckwaerts internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-3, LH-QA-01/Messung-1-von-2')")"
rot widerspruch-element-rueckwaerts "$wurzel"
stderr_ist widerspruch-element-rueckwaerts 'abdeckung: internal/a_test.go:4: LH-QA-01/Messung-1-von-2 widerspricht LH-QA-01/Messung-1-von-3 (früher: internal/a_test.go:4)'
# Ein widersprechender Eintrag legt das n nicht fest: Ein dritter mit dem n des zweiten widerspricht
# dem ersten und steht ebenfalls im Fehler.
wurzel="$(baum widerspruch-dritter internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-2')" internal/b_test.go "$(t 'LH-QA-01/Messung-1-von-3')" internal/c_test.go "$(t 'LH-QA-01/Messung-2-von-3')")"
rot widerspruch-dritter "$wurzel"
stderr_klasse widerspruch-dritter 2 'widerspricht LH-QA-01/Messung-1-von-2 \(früher: internal/a_test.go:4\)$'
# Die frühere Deklaration ist die erste, die das n festlegt: zwei frühere mit n = 3, danach n = 2.
wurzel="$(baum widerspruch-erste-n internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-3')" internal/b_test.go "$(t 'LH-QA-01/Messung-2-von-3')" internal/c_test.go "$(t 'LH-QA-01/Messung-1-von-2')")"
rot widerspruch-erste-n "$wurzel"
stderr_ist widerspruch-erste-n 'abdeckung: internal/c_test.go:4: LH-QA-01/Messung-1-von-2 widerspricht LH-QA-01/Messung-1-von-3 (früher: internal/a_test.go:4)'
# Bei `Messung` neben einem Teil die erste Deklaration des anderen Pfads, in beiden Richtungen.
wurzel="$(baum widerspruch-erste-messung internal/a_test.go "$(t 'LH-QA-01/Messung')" internal/b_test.go "$(t 'LH-QA-01/Messung')" internal/c_test.go "$(t 'LH-QA-01/Messung-1-von-2')")"
rot widerspruch-erste-messung "$wurzel"
stderr_ist widerspruch-erste-messung 'abdeckung: internal/c_test.go:4: LH-QA-01/Messung-1-von-2 steht neben LH-QA-01/Messung (früher: internal/a_test.go:4)'
wurzel="$(baum widerspruch-erster-teil internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-2')" internal/b_test.go "$(t 'LH-QA-01/Messung-2-von-2')" internal/c_test.go "$(t 'LH-QA-01/Messung')")"
rot widerspruch-erster-teil "$wurzel"
stderr_ist widerspruch-erster-teil 'abdeckung: internal/c_test.go:4: LH-QA-01/Messung steht neben LH-QA-01/Messung-1-von-2 (früher: internal/a_test.go:4)'

# Ein Lastenheft nur mit einer fremden Kennung besteht den Ausgang 2; die Deklaration scheitert
# an (4) mit Exit 1.
wurzel="$(baum lastenheft-fremde-kennung internal/a_test.go "$(t 'LH-FA-01/Happy')")"
printf '### LH-XX-01 — Fremd\n' > "$wurzel/spec/lastenheft.md"
rot lastenheft-fremde-kennung "$wurzel"
stderr_klasse lastenheft-fremde-kennung 1 'steht nicht im Lastenheft'

# Test und Gegenprobe, gleich in welcher Datei: tools/ liegt vor zz/ und hinter internal/.
wurzel="$(baum widerspruch-gate-frueher tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung-1-von-2 — g")" zz/z_test.go "$(t 'LH-QA-01/Messung-1-von-3')")"
rot widerspruch-gate-frueher "$wurzel"
stderr_ist widerspruch-gate-frueher 'abdeckung: zz/z_test.go:4: LH-QA-01/Messung-1-von-3 widerspricht LH-QA-01/Messung-1-von-2 (früher: tools/x-gegenprobe.sh:2)'
wurzel="$(baum widerspruch-gate-spaeter tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung-1-von-2 — g")" internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-3')")"
rot widerspruch-gate-spaeter "$wurzel"
stderr_ist widerspruch-gate-spaeter 'abdeckung: tools/x-gegenprobe.sh:2: LH-QA-01/Messung-1-von-2 widerspricht LH-QA-01/Messung-1-von-3 (früher: internal/a_test.go:4)'
# Zwei Gegenproben.
wurzel="$(baum widerspruch-gate-gate tools/a-gegenprobe.sh "$(gate "${D}LH-RB-01/Messung-1-von-2 — g")" tools/b-gegenprobe.sh "$(gate "${D}LH-RB-01/Messung-2-von-3 — h")")"
rot widerspruch-gate-gate "$wurzel"
stderr_ist widerspruch-gate-gate 'abdeckung: tools/b-gegenprobe.sh:2: LH-RB-01/Messung-2-von-3 widerspricht LH-RB-01/Messung-1-von-2 (früher: tools/a-gegenprobe.sh:2)'

# `Messung` neben einem Teil derselben Anforderung, in beiden Reihenfolgen und über Dateien.
wurzel="$(baum widerspruch-messung-dann-teil internal/a_test.go "$(t 'LH-QA-01/Messung')" internal/b_test.go "$(t 'LH-QA-01/Messung-1-von-2')")"
rot widerspruch-messung-dann-teil "$wurzel"
stderr_ist widerspruch-messung-dann-teil 'abdeckung: internal/b_test.go:4: LH-QA-01/Messung-1-von-2 steht neben LH-QA-01/Messung (früher: internal/a_test.go:4)'
wurzel="$(baum widerspruch-teil-dann-messung internal/a_test.go "$(t 'LH-QA-01/Messung-1-von-2')" internal/b_test.go "$(t 'LH-QA-01/Messung')")"
rot widerspruch-teil-dann-messung "$wurzel"
stderr_ist widerspruch-teil-dann-messung 'abdeckung: internal/b_test.go:4: LH-QA-01/Messung steht neben LH-QA-01/Messung-1-von-2 (früher: internal/a_test.go:4)'
wurzel="$(baum widerspruch-messung-gate tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung — g")" zz/z_test.go "$(t 'LH-QA-01/Messung-2-von-2')")"
rot widerspruch-messung-gate "$wurzel"
stderr_ist widerspruch-messung-gate 'abdeckung: zz/z_test.go:4: LH-QA-01/Messung-2-von-2 steht neben LH-QA-01/Messung (früher: tools/x-gegenprobe.sh:2)'

# Gutformen: verschiedene Anforderungen mit verschiedenem n, `Messung` mehrfach, derselbe
# Teil zweimal (zwei Zeilen), i = n, n = 9, die Pfade an Tests im Unit- und im E2E-Verzeichnis.
wurzel="$(baum teil-gut internal/a_test.go "$(t 'LH-QA-01/Messung-2-von-2, LH-RB-01/Messung-9-von-9')" internal/b_test.go "$(t 'LH-QA-01/Messung-2-von-2')" internal/c_test.go "$(t 'LH-FA-01/Happy')")"
angenommen teil-gut "$wurzel"
wurzel="$(baum messung-mehrfach internal/a_test.go "$(t 'LH-QA-01/Messung')" tools/x-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung — g")")"
angenommen messung-mehrfach "$wurzel"

# --- Tabellen mit einer geteilten Messung (Punkt 4) ----------------------------------
# QA-01 vollständig aus drei Nachweisarten, RB-01 mit Teil 1 aus allen drei und ohne Teil 2.
wurzel="$(baum teile \
  tools/b-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung-1-von-3 — q1" "${D}LH-RB-01/Messung-1-von-2 — aa")" \
  tools/a-gegenprobe.sh "$(gate "${D}LH-RB-01/Messung-1-von-2 — zz" "${D}LH-QA-01/Messung-2-von-3 — q2")" \
  internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-QA-01/Messung-3-von-3, LH-RB-01/Messung-1-von-2 — u\nfunc TestU(t *testing.T) {}')" \
  test/integration/e_test.go "$(printf 'package e\n\n// Abdeckung: LH-QA-01/Messung-2-von-3, LH-RB-01/Messung-1-von-2 — e\nfunc TestE2EE(t *testing.T) {}')")"
angenommen teile "$wurzel"
tabelle teile-gate "$wurzel" abdeckung-gate.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' '| --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | Messung-1-von-3 | q1 | \`tools/b-gegenprobe.sh\` |" \
  "| [\`LH-QA-01\`$L | Messung-2-von-3 | q2 | \`tools/a-gegenprobe.sh\` |" \
  "| [\`LH-RB-01\`$L | Messung-1-von-2 | zz | \`tools/a-gegenprobe.sh\` |" \
  "| [\`LH-RB-01\`$L | Messung-1-von-2 | aa | \`tools/b-gegenprobe.sh\` |"
tabelle teile-unit "$wurzel" abdeckung-unit.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis | Ort |' '| --- | --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | Messung-3-von-3 | u | \`TestU\` | \`internal/a_test.go\` |" \
  "| [\`LH-RB-01\`$L | Messung-1-von-2 | u | \`TestU\` | \`internal/a_test.go\` |"
tabelle teile-e2e "$wurzel" abdeckung-e2e.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis | Ort |' '| --- | --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | Messung-2-von-3 | e | \`TestE2EE\` | \`test/integration/e_test.go\` |" \
  "| [\`LH-RB-01\`$L | Messung-1-von-2 | e | \`TestE2EE\` | \`test/integration/e_test.go\` |"
tabelle teile-gesamt "$wurzel" abdeckung-gesamt.md \
  '| Anforderung | Happy | Boundary | Negative | Messung | Stand |' '| --- | --- | --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | n/a | n/a | n/a | 1/3 Gate; 2/3 E2E, Gate; 3/3 Unit | vollständig |" \
  "| [\`LH-RB-01\`$L | n/a | n/a | n/a | 1/2 E2E, Gate, Unit; 2/2 — | teilweise |"
tabelle teile-vollstaendig "$wurzel" abdeckung-vollstaendig.md \
  '| Anforderung | Nachweisarten |' '| --- | --- |' "| [\`LH-QA-01\`$L | E2E, Gate, Unit |"

# Derselbe Teil in zwei Nachweisarten zählt einmal; mit Teil 2 ist RB-01 vollständig.
wurzel="$(baum teile-voll \
  tools/a-gegenprobe.sh "$(gate "${D}LH-RB-01/Messung-1-von-2 — g")" \
  internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-RB-01/Messung-1-von-2, LH-RB-01/Messung-2-von-2 — u\nfunc TestU(t *testing.T) {}')")"
angenommen teile-voll "$wurzel"
tabelle teile-voll-gesamt "$wurzel" abdeckung-gesamt.md \
  '| Anforderung | Happy | Boundary | Negative | Messung | Stand |' '| --- | --- | --- | --- | --- | --- |' \
  "| [\`LH-RB-01\`$L | n/a | n/a | n/a | 1/2 Gate, Unit; 2/2 Unit | vollständig |"
tabelle teile-voll-vollstaendig "$wurzel" abdeckung-vollstaendig.md \
  '| Anforderung | Nachweisarten |' '| --- | --- |' "| [\`LH-RB-01\`$L | Gate, Unit |"

# Die Teile in der Reihenfolge ihrer Nummer, auch wenn sie anders deklariert sind.
wurzel="$(baum teile-reihenfolge tools/a-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung-3-von-3, LH-QA-01/Messung-1-von-3, LH-QA-01/Messung-2-von-3 — t")")"
angenommen teile-reihenfolge "$wurzel"
tabelle teile-reihenfolge "$wurzel" abdeckung-gate.md \
  '| Anforderung | Pfad | Kurzbeschreibung | Nachweis |' '| --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | Messung-1-von-3 | t | \`tools/a-gegenprobe.sh\` |" \
  "| [\`LH-QA-01\`$L | Messung-2-von-3 | t | \`tools/a-gegenprobe.sh\` |" \
  "| [\`LH-QA-01\`$L | Messung-3-von-3 | t | \`tools/a-gegenprobe.sh\` |"
tabelle teile-reihenfolge-gesamt "$wurzel" abdeckung-gesamt.md \
  '| Anforderung | Happy | Boundary | Negative | Messung | Stand |' '| --- | --- | --- | --- | --- | --- |' \
  "| [\`LH-QA-01\`$L | n/a | n/a | n/a | 1/3 Gate; 2/3 Gate; 3/3 Gate | vollständig |"

# Jeder Teil zählt: fehlt der erste, ein mittlerer, der letzte oder der neunte von neun, ist die
# Anforderung teilweise und steht nicht in der vollständigen Tabelle.
teilweise() {
  local fall="$1" n="$2" fehlt="$3" ids="" i darst="" w
  for ((i = 1; i <= n; i++)); do
    [ "$i" -eq "$fehlt" ] && { darst="${darst:+$darst; }$i/$n —"; continue; }
    ids="${ids:+$ids, }LH-QA-01/Messung-$i-von-$n"
    darst="${darst:+$darst; }$i/$n Gate"
  done
  w="$(baum "$fall" tools/a-gegenprobe.sh "$(gate "${D}$ids — t")")"
  angenommen "$fall" "$w"
  tabelle "$fall-gesamt" "$w" abdeckung-gesamt.md \
    '| Anforderung | Happy | Boundary | Negative | Messung | Stand |' '| --- | --- | --- | --- | --- | --- |' \
    "| [\`LH-QA-01\`$L | n/a | n/a | n/a | $darst | teilweise |"
  tabelle "$fall-vollstaendig" "$w" abdeckung-vollstaendig.md '| Anforderung | Nachweisarten |' '| --- | --- |'
}
teilweise teil-fehlt-erster 3 1
teilweise teil-fehlt-mittlerer 3 2
teilweise teil-fehlt-letzter 3 3
teilweise teil-fehlt-neunter 9 9
teilweise teil-fehlt-erster-von-neun 9 1

# n = 9 vollständig.
ids=""
for i in 1 2 3 4 5 6 7 8 9; do ids="${ids:+$ids, }LH-RB-01/Messung-$i-von-9"; done
wurzel="$(baum teil-neun-voll tools/a-gegenprobe.sh "$(gate "${D}$ids — t")")"
angenommen teil-neun-voll "$wurzel"
tabelle teil-neun-voll "$wurzel" abdeckung-vollstaendig.md '| Anforderung | Nachweisarten |' '| --- | --- |' "| [\`LH-RB-01\`$L | Gate |"

# --check kennt die Tabelle der Nachweisart Gate: veraltet oder fehlend, Exit 1 und ihr Pfad.
wurzel="$(baum check-gate tools/a-gegenprobe.sh "$(gate "${D}LH-QA-01/Messung — t")")"
bash "$skript" "$wurzel" >/dev/null 2>&1 || melde "Fall 'check-gate': gültiger Baum abgelehnt"
for ver in veraendert fehlt; do
  if [ "$ver" = veraendert ]; then echo veraendert > "$wurzel/docs/user/abdeckung-gate.md"; else rm "$wurzel/docs/user/abdeckung-gate.md"; fi
  set +e
  bash "$skript" --check "$wurzel" >/dev/null 2>"$arbeit/err"
  code=$?
  set -e
  { [ "$code" -eq 1 ] && grep -qF 'docs/user/abdeckung-gate.md' "$arbeit/err"; } \
    || melde "Fall 'check-gate-$ver' erwartet Exit 1 mit dem Pfad der Tabelle, bekam $code: $(cat "$arbeit/err")"
  bash "$skript" "$wurzel" >/dev/null 2>&1
done

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "abdeckung-gegenprobe: gruen — Fehlformen abgelehnt, gültige Bäume mit dem erwarteten Inhalt der fünf Tabellen, --check schreibt nicht"
