#!/usr/bin/env bash
# abdeckung-gegenprobe — prüft tools/test/abdeckung.sh an kleinen Bäumen in einem
# Temp-Verzeichnis; der Arbeitsbaum bleibt unberührt. Abgelehnt werden müssen:
# eine Leerzeile zwischen Deklaration und Test, eine Deklaration über einer
# Hilfsfunktion, eine eingerückte Deklaration, ein unbekannter Pfad, eine
# Anforderung, die nicht im Lastenheft steht, ein TestE2E ohne Deklaration und
# eine Deklaration ohne „ — “. Angenommen werden muss ein gültiger Baum; dessen
# vollständige Anforderungen, das maskierte `|`, die Dateirechte 0644 und das
# Verhalten von --check (prüft, schreibt nicht) werden geprüft.
#
# Ausgang: 0, wenn jeder Fall wie erwartet endet, sonst 1.
set -euo pipefail

skript="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/abdeckung.sh"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT
fehler=0

# baum <name> <datei> <inhalt> — legt einen Baum mit Lastenheft und einer Testdatei an.
baum() {
  local wurzel="$arbeit/$1"
  mkdir -p "$wurzel/spec" "$wurzel/docs/user" "$wurzel/$(dirname "$2")"
  printf '### LH-FA-01 — Eins\n\n### LH-QA-01 — Qualität\n' > "$wurzel/spec/lastenheft.md"
  printf '%s\n' "$3" > "$wurzel/$2"
  echo "$wurzel"
}

abgelehnt() {
  local name="$1" datei="$2" inhalt="$3" wurzel
  wurzel="$(baum "$name" "$datei" "$inhalt")"
  if bash "$skript" "$wurzel" >/dev/null 2>&1; then
    echo "abdeckung-gegenprobe: ROT — Fall '$name' angenommen" >&2
    fehler=1
  fi
}

abgelehnt leerzeile internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy — x\n\nfunc TestA(t *testing.T) {}')"
abgelehnt hilfsfunktion internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy — x\nfunc hilfe() {}')"
abgelehnt eingerueckt internal/a_test.go "$(printf 'package a\n\n\t// Abdeckung: LH-FA-01/Happy — x\nfunc TestA(t *testing.T) {}')"
abgelehnt pfad internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Irgendwie — x\nfunc TestA(t *testing.T) {}')"
abgelehnt unbekannt internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-99/Happy — x\nfunc TestA(t *testing.T) {}')"
abgelehnt ohne-deklaration test/integration/a_test.go "$(printf 'package a\n\nfunc TestE2EA(t *testing.T) {}')"
abgelehnt ohne-trenner internal/a_test.go "$(printf 'package a\n\n// Abdeckung: LH-FA-01/Happy x\nfunc TestA(t *testing.T) {}')"

gueltig="$(printf '%s\n' \
  'package a' '' \
  '// Abdeckung: LH-FA-01/Happy, LH-FA-01/Boundary — a | b' \
  '// Fortsetzung' \
  'func TestA(t *testing.T) {}' '' \
  '// Abdeckung: LH-FA-01/Negative, LH-QA-01/Messung — c' \
  'func TestB(t *testing.T) {}')"
wurzel="$(baum gueltig internal/a_test.go "$gueltig")"
if ! bash "$skript" "$wurzel" >/dev/null 2>&1; then
  echo "abdeckung-gegenprobe: ROT — gültiger Baum abgelehnt" >&2
  fehler=1
else
  voll="$wurzel/docs/user/abdeckung-vollstaendig.md"
  grep -q 'LH-FA-01' "$voll" && grep -q 'LH-QA-01' "$voll" \
    || { echo "abdeckung-gegenprobe: ROT — vollständige Anforderungen fehlen" >&2; fehler=1; }
  grep -qF 'a \| b Fortsetzung' "$wurzel/docs/user/abdeckung-unit.md" \
    || { echo "abdeckung-gegenprobe: ROT — | nicht maskiert oder Fortsetzung fehlt" >&2; fehler=1; }
  [ "$(stat -c %a "$voll")" = 644 ] \
    || { echo "abdeckung-gegenprobe: ROT — Dateirechte nicht 0644" >&2; fehler=1; }
  bash "$skript" --check "$wurzel" >/dev/null 2>&1 \
    || { echo "abdeckung-gegenprobe: ROT — --check auf aktuellem Baum rot" >&2; fehler=1; }
  echo veraendert > "$voll"
  if bash "$skript" --check "$wurzel" >/dev/null 2>&1; then
    echo "abdeckung-gegenprobe: ROT — --check auf veraltetem Baum grün" >&2
    fehler=1
  fi
  [ "$(cat "$voll")" = veraendert ] \
    || { echo "abdeckung-gegenprobe: ROT — --check hat geschrieben" >&2; fehler=1; }
fi

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "abdeckung-gegenprobe: gruen — Fehlformen abgelehnt, gültiger Baum vollständig, --check schreibt nicht"
