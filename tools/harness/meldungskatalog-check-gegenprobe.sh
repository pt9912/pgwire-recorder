#!/usr/bin/env bash
# meldungskatalog-check-gegenprobe — prüft tools/harness/meldungskatalog-check.sh gegen seinen Vertrag
# (harness/vertraege/SPEC-054-meldungskatalog.md) an kleinen Bäumen in einem Temp-Verzeichnis. Je Fall
# ein Baum; erwartet ist Exit 0 ohne Ausgabe (grün), Exit 1 mit genau den Befund-Zeilen auf stderr
# (rot) oder Exit 2 mit der Zeile der Sperre. Die Fälle sind nach den Punkten des Vertrags gruppiert
# (Präfix nr<n>-). Regel an den Schreiber dieser Gegenprobe: Sie baut jeden Baum selbst und liest
# aus dem Arbeitsbaum nur das Skript, das Makefile und die Fragmente (Punkt 11); Handbuch,
# Spezifikation und Planungsbaum des Arbeitsbaums liest sie nicht.
#
# GRENZE der Gegenprobe: Als root bleibt eine Datei mit Rechten 000 lesbar; dann sind die Fälle
# nr8-unlesbar-* nicht herstellbar und werden laut übersprungen. Der Fall nr9-bytefolge braucht die
# Locale de_DE.UTF-8 und wird ohne sie laut übersprungen. Eine Sperre `internal/ fehlt` ist nicht
# erreichbar, weil die Code-Tabelle unter internal/ liegt und ihr Fehlen zuerst gilt.
#
# Ausgang: 0, wenn jeder Fall wie erwartet endet, sonst 1.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
skript="$hier/meldungskatalog-check.sh"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT
fehler=0
nr=0
F="internal/hexagon/model/fehler.go"
S="spec/spezifikation.md"
B="docs/user/benutzerhandbuch.md"
BC="PGR-E1000 PGR-E2001 PGR-W2001"
BK="PGR-E1000 PGR-E2000 PGR-E2001 PGR-E3000 PGR-E4000 PGR-E5000 PGR-E6000 PGR-W2001"
BF="PGR-E1000 PGR-E2001"
BW="PGR-W2001"

melde() { echo "meldungskatalog-check-gegenprobe: ROT — $*" >&2; fehler=1; }

tabelle() { local c; for c; do printf '| `%s` | Klasse | Text |\n' "$c"; done; }

# bau [<c> <k> <hf> <hw>] — legt einen Baum $w an; die Listen sind die Codes je Quelle.
bau() {
  local i=0 c
  nr=$((nr + 1)); w="$arbeit/$nr"
  mkdir -p "$w/cmd/x" "$w/internal/hexagon/model" "$w/spec" "$w/docs/user"
  { printf 'package model\n\n// Meldungscodes\nconst (\n'
    for c in ${1-$BC}; do i=$((i + 1)); printf '\tCode%d = "%s"\n' "$i" "$c"; done
    printf ')\n'; } > "$w/$F"
  { printf '# Spezifikation\n\n### SPEC-034 — Meldungscodes\n\n| Code | Klasse | Bedeutung |\n|---|---|---|\n'
    tabelle ${2-$BK}; } > "$w/$S"
  { printf '# Handbuch\n\n### Fehlercodes\n\n| Code | Bedeutung | Lösung |\n|---|---|---|\n'
    tabelle ${3-$BF}
    printf '\n### Warnungen\n\n| Code | Bedeutung | Hinweis |\n|---|---|---|\n'
    tabelle ${4-$BW}; } > "$w/$B"
  printf 'package main\n' > "$w/cmd/x/main.go"
}
czeile() { sed -i '$d' "$w/$F"; printf '%s\n)\n' "$1" >> "$w/$F"; }   # Zeile ans Ende der Konstanten
sk() { printf '%s\n' "$1" >> "$w/$S"; }                              # Zeile ans Ende von SPEC-034
z() { grep -nF -- "$2" "$w/$1" | sed -n "${3:-1}p" | cut -d: -f1; }   # Zeile <n> mit dem Muster

lauf() {
  set +e
  bash "$skript" "$1" >"$arbeit/out" 2>"$arbeit/err"
  code=$?
  set -e
  aus="$(cat "$arbeit/out")"; err="$(cat "$arbeit/err")"
}
gruen() {
  lauf "$2"
  if [ "$code" -ne 0 ] || [ -n "$aus$err" ]; then melde "Fall '$1' erwartet grün, Exit $code: $err"; fi
}
# rot <name> <wurzel> <zeile>… — Exit 1, leeres stdout, genau die Zeilen ohne Präfix in dieser Reihenfolge.
rot() {
  local name="$1" wurzel="$2" soll
  shift 2
  lauf "$wurzel"
  soll="$(printf 'meldungskatalog-check: %s\n' "$@")"
  if [ "$code" -ne 1 ]; then melde "Fall '$name' erwartet Exit 1, bekam $code: $err"
  elif [ -n "$aus" ]; then melde "Fall '$name' schreibt auf stdout: $aus"
  elif [ "$err" != "$soll" ]; then melde "Fall '$name' erwartet:"$'\n'"$soll"$'\n'"bekam:"$'\n'"$err"; fi
}
# sp <name> <wurzel> <grund> — Exit 2, leeres stdout, genau die Zeile der Sperre.
sp() {
  lauf "$2"
  if [ "$code" -ne 2 ]; then melde "Fall '$1' erwartet Exit 2, bekam $code: $err"
  elif [ -n "$aus" ]; then melde "Fall '$1' schreibt auf stdout: $aus"
  elif [ "$err" != "meldungskatalog-check: $3" ]; then melde "Fall '$1' erwartet '$3', bekam: $err"; fi
}
# Hilfen für die häufigen Formen: eine Konstante, die gelesen wird (Befunde nicht im Handbuch und
# nicht im Katalog) und eine, die nicht gelesen wird (grün).
c_gelesen() {
  bau; czeile "$2"
  rot "$1" "$w" "$F:$(z $F PGR-E9001): PGR-E9001: nicht im Handbuch" "$F:$(z $F PGR-E9001): PGR-E9001: nicht im Katalog"
}
c_nicht() { bau; czeile "$2"; gruen "$1" "$w"; }

# --- Basis -------------------------------------------------------------------
bau; gruen basis "$w"

# --- Nr. 1: Quellen ------------------------------------------------------------
c_gelesen nr1-c-einzug-leerzeichen '    CodeS = "PGR-E9001"'
c_gelesen nr1-c-ohne-einzug 'CodeS = "PGR-E9001"'
c_gelesen nr1-c-leerraum-tab $'\tCodeS\t=\t"PGR-E9001"'
c_gelesen nr1-c-nur-code $'\tCode = "PGR-E9001"'
c_gelesen nr1-c-wortzeichen $'\tCodeA_1b = "PGR-E9001"'
c_nicht nr1-c-kommentar $'\t// CodeK = "PGR-E9001"'
c_nicht nr1-c-anderer-name $'\tKlasse1 = "PGR-E9001"\n\tNotCode = "PGR-E9001"'
c_nicht nr1-c-typisiert $'\tCodeT Typ = "PGR-E9001"'
c_nicht nr1-c-backtick $'\tCodeB = `PGR-E9001`'
bau; sk '| PGR-E9001 | ohne Backticks |'; sk '| Text | `PGR-E9001` |'; gruen nr1-k-keine-zeile "$w"
bau; sed -i '1i | `PGR-E9001` | davor |' "$w/$S"; gruen nr1-k-davor "$w"
bau; sk $'### SPEC-035 — Folge\n\n| `PGR-E9001` | danach |'; gruen nr1-k-danach "$w"
bau; sed -i '1i | `PGR-E9001` | davor |' "$w/$B"; gruen nr1-h-davor "$w"
bau; sed -i '/^### Warnungen$/i ### Anderes\n\n| `PGR-E9001` | x |\n' "$w/$B"; gruen nr1-h-zwischen "$w"
bau; printf '\n### Anderes\n\n| `PGR-E9001` | x |\n' >> "$w/$B"; gruen nr1-h-danach "$w"

# --- Nr. 2: Form ---------------------------------------------------------------
bau "$BC PGR-E123"; rot nr2-form-c "$w" "$F:$(z $F PGR-E123): PGR-E123: Form"
bau "$BC" "$BK PGR-X0001"; rot nr2-form-k "$w" "$S:$(z $S PGR-X0001): PGR-X0001: Form"
bau "$BC" "$BK" "$BF PGR-E<k>000"; rot nr2-form-h-platzhalter "$w" "$B:$(z $B 'PGR-E<k>000'): PGR-E<k>000: Form"
bau; sk '| `PGR-E2002`, `PGR-E2003` | zwei Codes |'
rot nr2-form-zwei-codes "$w" "$S:$(z $S PGR-E2002): PGR-E2002\`, \`PGR-E2003: Form"
bau "$BC" "$BK PGR-E20011"; rot nr2-form-zu-lang "$w" "$S:$(z $S PGR-E20011): PGR-E20011: Form"
bau "$BC" "$BK xPGR-E2001"; rot nr2-form-vorsatz "$w" "$S:$(z $S xPGR-E2001): xPGR-E2001: Form"
bau "$BC" "$BK PGR-I0001"; rot nr2-schwere-i "$w" "$S:$(z $S PGR-I0001): PGR-I0001: nicht in der Code-Tabelle"

# --- Nr. 3: doppelt ------------------------------------------------------------
bau "$BC PGR-E2001"; rot nr3-doppelt-c "$w" "$F:$(z $F PGR-E2001 2): PGR-E2001: doppelt"
bau "$BC" "$BK PGR-E2001"; rot nr3-doppelt-k "$w" "$S:$(z $S PGR-E2001 2): PGR-E2001: doppelt"
bau "$BC" "$BK" "$BF PGR-E2001"; rot nr3-doppelt-h "$w" "$B:$(z $B PGR-E2001 2): PGR-E2001: doppelt"
bau "$BC" "$BK" "$BF" "$BW PGR-E2001"; rot nr3-doppelt-h-abschnitte "$w" "$B:$(z $B PGR-E2001 2): PGR-E2001: doppelt"
bau "$BC PGR-E2001 PGR-E2001"; rot nr3-doppelt-dreifach "$w" "$F:$(z $F PGR-E2001 2): PGR-E2001: doppelt"

# --- Nr. 4: Mengen -------------------------------------------------------------
bau "$BC PGR-E2003" "$BK" "$BF PGR-E2003"; rot nr4a-nicht-im-katalog "$w" "$F:$(z $F PGR-E2003): PGR-E2003: nicht im Katalog"
bau "$BC" "$BK PGR-E2002"; rot nr4b-nicht-in-c "$w" "$S:$(z $S PGR-E2002): PGR-E2002: nicht in der Code-Tabelle"
bau "$BC" "$BK PGR-W3000"; gruen nr4b-rueckfall-w "$w"
bau "$BC" "$BK PGR-E2100"; rot nr4b-endung-00 "$w" "$S:$(z $S PGR-E2100): PGR-E2100: nicht in der Code-Tabelle"
bau "$BC" "$BK" "PGR-E1000"; rot nr4c-fehler "$w" "$F:$(z $F PGR-E2001): PGR-E2001: nicht im Handbuch"
bau "$BC" "$BK" "$BF" ""; rot nr4c-warnung "$w" "$F:$(z $F PGR-W2001): PGR-W2001: nicht im Handbuch"
bau "$BC" "$BK" "$BF PGR-E2000"; rot nr4d-in-k "$w" "$B:$(z $B PGR-E2000): PGR-E2000: nicht in der Code-Tabelle"
bau "$BC" "$BK" "$BF PGR-E9999"; rot nr4d-nicht-in-k "$w" "$B:$(z $B PGR-E9999): PGR-E9999: nicht in der Code-Tabelle"
bau "$BC" "$BK" "$BF" "$BW PGR-W9999"; rot nr4d-warnung "$w" "$B:$(z $B PGR-W9999): PGR-W9999: nicht in der Code-Tabelle"

# --- Nr. 5: Rückfall je Klasse ---------------------------------------------------
bau "$BC" "${BK/PGR-E1000/}"
rot nr5-klasse-1 "$w" "$F:$(z $F PGR-E1000): PGR-E1000: nicht im Katalog" "$S: PGR-E1000: Rückfall fehlt"
for k in 2 3 4 5 6; do
  bau "$BC" "${BK/PGR-E${k}000/}"; rot "nr5-klasse-$k" "$w" "$S: PGR-E${k}000: Rückfall fehlt"
done

# --- Nr. 6: Code außerhalb der Code-Tabelle ------------------------------------
lit() { bau; mkdir -p "$(dirname "$w/$2")"; printf 'package p\n\nvar a = %s\n' "$3" > "$w/$2"; }
lit x cmd/x/lit.go '"PGR-E2001: Text"'; rot nr6-cmd "$w" "cmd/x/lit.go:3: PGR-E2001: Code außerhalb der Code-Tabelle"
lit x internal/adapter/a.go 'fmt.Sprint("PGR-W9999")'; rot nr6-internal "$w" "internal/adapter/a.go:3: PGR-W9999: Code außerhalb der Code-Tabelle"
lit x internal/adapter/fehler.go '"PGR-E2001"'; rot nr6-anderes-fehler-go "$w" "internal/adapter/fehler.go:3: PGR-E2001: Code außerhalb der Code-Tabelle"
lit x cmd/x/lit.go '1 // "PGR-E2001"'; rot nr6-kommentar-mit-anfuehrung "$w" "cmd/x/lit.go:3: PGR-E2001: Code außerhalb der Code-Tabelle"
lit x cmd/x/lit_test.go '"PGR-E2001"'; mkdir -p "$w/internal/adapter"; printf 'package p\nvar a = "PGR-E2001"\n' > "$w/internal/adapter/a_test.go"; gruen nr6-test "$w"
lit x cmd/x/lit.go '`PGR-E2001`'; gruen nr6-backtick "$w"
lit x cmd/x/lit.go '1 // siehe PGR-E2001'; gruen nr6-kommentar-ohne-anfuehrung "$w"
lit x test/lit.go '"PGR-E2001"'; gruen nr6-test-verzeichnis "$w"
lit x internal/adapter/hinweis.txt '"PGR-E2001"'; gruen nr6-nicht-go "$w"
c_nicht nr6-fehler-go-anderer-name $'\tAndere = "PGR-E9001"'

# --- Nr. 8: Sperren und ihre Reihenfolge ---------------------------------------
bau; rm "$w/$F"; sp nr8-c-fehlt "$w" "$F fehlt"
bau; rm "$w/$S"; sp nr8-k-fehlt "$w" "$S fehlt"
bau; rm "$w/$B"; sp nr8-h-fehlt "$w" "$B fehlt"
bau; rm -r "$w/cmd"; sp nr8-cmd-fehlt "$w" "cmd/ fehlt"
for q in F S B; do
  bau; chmod 000 "$w/${!q}"
  if [ -r "$w/${!q}" ]; then echo "meldungskatalog-check-gegenprobe: Fall 'nr8-unlesbar-$q' übersprungen — Datei mit Rechten 000 ist lesbar (root?)" >&2
  else sp "nr8-unlesbar-$q" "$w" "${!q} nicht lesbar"; chmod 644 "$w/${!q}"; fi
done
bau; sed -i 's/^### SPEC-034 /### SPEC-0341 /' "$w/$S"; sp nr8-abschnitt-k "$w" "Abschnitt ### SPEC-034 fehlt"
bau; sed -i 's/^### Fehlercodes$/### Fehlercodes und mehr/' "$w/$B"; sp nr8-abschnitt-fehlercodes "$w" "Abschnitt ### Fehlercodes fehlt"
bau; sed -i 's/^### Warnungen$/#### Warnungen/' "$w/$B"; sp nr8-abschnitt-warnungen "$w" "Abschnitt ### Warnungen fehlt"
bau ""; sp nr8-c-ohne-code "$w" "$F ohne Code"
bau "$BC" ""; sp nr8-k-ohne-codes "$w" "$S ohne Codes"
bau "$BC" "$BK" "" ""; sp nr8-h-ohne-codes "$w" "$B ohne Codes"
bau "PGR-E1000 PGR-E2001" "PGR-E1000 PGR-E2000 PGR-E2001 PGR-E3000 PGR-E4000 PGR-E5000 PGR-E6000" "$BF" ""; gruen nr8-h-nur-fehlercodes "$w"
bau "PGR-W2001" "PGR-E1000 PGR-E2000 PGR-E3000 PGR-E4000 PGR-E5000 PGR-E6000 PGR-W2001" "" "$BW"; gruen nr8-h-nur-warnungen "$w"
bau "$BC PGR-E123"; sed -i 's/^### SPEC-034 /### SPEC-0341 /' "$w/$S"; sp nr8-sperre-vor-befund "$w" "Abschnitt ### SPEC-034 fehlt"
bau; rm "$w/$F" "$w/$S"; sp nr8-reihenfolge-c-vor-k "$w" "$F fehlt"
bau; rm "$w/$S" "$w/$B"; sp nr8-reihenfolge-k-vor-h "$w" "$S fehlt"
bau; rm "$w/$B"; rm -r "$w/cmd"; sp nr8-reihenfolge-datei-vor-verzeichnis "$w" "$B fehlt"
bau; rm -r "$w/cmd"; sed -i 's/^### SPEC-034 /### SPEC-0341 /' "$w/$S"; sp nr8-reihenfolge-verzeichnis-vor-abschnitt "$w" "cmd/ fehlt"
bau ""; sed -i 's/^### SPEC-034 /### SPEC-0341 /' "$w/$S"; sp nr8-reihenfolge-abschnitt-vor-quelle "$w" "Abschnitt ### SPEC-034 fehlt"
bau; sed -i 's/^### SPEC-034 /### SPEC-0341 /' "$w/$S"; sed -i 's/^### Fehlercodes$/### Fehler/' "$w/$B"; sp nr8-reihenfolge-abschnitt-k-vor-f "$w" "Abschnitt ### SPEC-034 fehlt"
bau; sed -i 's/^### Fehlercodes$/### Fehler/; s/^### Warnungen$/### Warn/' "$w/$B"; sp nr8-reihenfolge-abschnitt-f-vor-w "$w" "Abschnitt ### Fehlercodes fehlt"
bau "" ""; sp nr8-reihenfolge-quelle-c-vor-k "$w" "$F ohne Code"
bau "$BC" "" "" ""; sp nr8-reihenfolge-quelle-k-vor-h "$w" "$S ohne Codes"

# --- Nr. 9: Ausgabe und Ausgang ------------------------------------------------
bau "$BC PGR-E2003" "$BK PGR-E2002" "$BF PGR-E2004"; mkdir -p "$w/cmd/y"; printf 'package y\nvar a = "PGR-E2005"\n' > "$w/cmd/y/lit.go"
rot nr9-sortiert "$w" "cmd/y/lit.go:2: PGR-E2005: Code außerhalb der Code-Tabelle" \
  "$B:$(z $B PGR-E2004): PGR-E2004: nicht in der Code-Tabelle" \
  "$F:$(z $F PGR-E2003): PGR-E2003: nicht im Handbuch" "$F:$(z $F PGR-E2003): PGR-E2003: nicht im Katalog" \
  "$S:$(z $S PGR-E2002): PGR-E2002: nicht in der Code-Tabelle"
bau "PGR-W2001 PGR-E2001 PGR-E1000" "PGR-E1000 PGR-E2000 PGR-E3000 PGR-E4000 PGR-E5000 PGR-E6000"
rot nr9-code-vor-zeile "$w" "$F:$(z $F PGR-E2001): PGR-E2001: nicht im Katalog" "$F:$(z $F PGR-W2001): PGR-W2001: nicht im Katalog"
bau "PGR-E2001 PGR-E1000 PGR-W2001 PGR-E2001" "$BK" "PGR-E1000"
rot nr9-befund-vor-zeile "$w" "$F:$(z $F PGR-E2001 2): PGR-E2001: doppelt" "$F:$(z $F PGR-E2001): PGR-E2001: nicht im Handbuch"
bau "PGR-E1000 PGR-E2001 PGR-W2001 PGR-W3001 PGR-E123 PGR-E123" "$BK PGR-W3001" "$BF" "$BW PGR-W3001"
rot nr9-zeile-numerisch "$w" "$F:9: PGR-E123: Form" "$F:10: PGR-E123: Form"
if locale -a | grep -qix 'de_DE.utf-\?8'; then
  bau "$BC" "$BK PGR-E2002 abc"; LANG=de_DE.UTF-8 rot nr9-bytefolge "$w" "$S:$(z $S PGR-E2002): PGR-E2002: nicht in der Code-Tabelle" "$S:$(z $S abc): abc: Form"
else
  echo "meldungskatalog-check-gegenprobe: Fall 'nr9-bytefolge' übersprungen — Locale de_DE.UTF-8 fehlt" >&2
fi

# --- Nr. 10: Aufrufarten -------------------------------------------------------
bau "$BC PGR-E2003"
(cd "$w" && bash "$skript") >"$arbeit/out" 2>"$arbeit/err" && code=0 || code=$?
[ "$code" -eq 1 ] && grep -qxF "meldungskatalog-check: $F:$(z $F PGR-E2003): PGR-E2003: nicht im Katalog" "$arbeit/err" \
  || melde "Fall 'nr10-wurzel-aktuell' erwartet Exit 1 mit Befund, bekam $code"
bestand() { (cd "$1" && find . | sort && find . -type f -exec cksum {} + | sort); }
for form in gruen rot; do
  if [ "$form" = gruen ]; then bau; else bau "$BC PGR-E2003"; fi
  vorher="$(bestand "$w")"; lauf "$w"
  [ "$vorher" = "$(bestand "$w")" ] || melde "Fall 'nr10-schreibt-nichts-$form': der Lauf hat den Baum verändert"
done
for flag in n q t i; do
  bau "$BC PGR-E2003"; vorher="$(MAKEFLAGS="-$flag" bash "$skript" "$w" 2>&1 || true)"
  zeile="$F:$(z $F PGR-E2003): PGR-E2003"
  [ "$vorher" = "$(printf 'meldungskatalog-check: %s\n' "$zeile: nicht im Handbuch" "$zeile: nicht im Katalog")" ] \
    || melde "Fall 'nr10-makeflags-$flag': MAKEFLAGS=-$flag ändert die Aussage: $vorher"
done

# --- Nr. 11: Anschluss ohne Stufung --------------------------------------------
# Makefile und Fragmente des Arbeitsbaums in einem Temp-Baum: Beide Ziele stehen in GATE_CHECKS, die
# Prüfung in der schnellen Kette, die Gegenprobe in der Ausnahmeliste; beide enden über make bei
# rotem Skript mit Exit ungleich 0, und make -n, -q, -t, -i hinterlassen nichts.
anschluss() {
  bau "$BC$1"
  cp "$hier/../../Makefile" "$w/Makefile"
  mkdir -p "$w/harness/mk" "$w/tools/harness"
  cp "$hier/../../harness/mk/meldungskatalog-check.mk" "$hier/../../harness/mk/gates-schnell.mk" "$w/harness/mk/"
  cp "$skript" "$w/tools/harness/meldungskatalog-check.sh"
  printf 'exit 1\n' > "$w/tools/harness/meldungskatalog-check-gegenprobe.sh"
}
anschluss ""; gruen_baum="$w"
anschluss " PGR-E2003"
datenbank="$(make -C "$w" -pq help 2>/dev/null || true)"
for ziel in meldungskatalog-check meldungskatalog-check-gegenprobe; do
  printf '%s\n' "$datenbank" | awk -v z="$ziel" '/^GATE_CHECKS :?= / { for (i = 3; i <= NF; i++) if ($i == z) t = 1 } END { exit !t }' \
    || melde "Fall 'nr11-gate-checks': $ziel hängt nicht an GATE_CHECKS"
  if make -s -C "$w" "$ziel" >/dev/null 2>&1; then melde "Fall 'nr11-scharf': make $ziel endet bei rotem Skript mit Exit 0"; fi
done
printf '%s\n' "$datenbank" | awk '/^GATE_SCHNELL_AUSNAHME :?= / { for (i = 3; i <= NF; i++) if ($i == "meldungskatalog-check-gegenprobe") t = 1 } END { exit !t }' \
  || melde "Fall 'nr11-ausnahme-gegenprobe': meldungskatalog-check-gegenprobe fehlt in GATE_SCHNELL_AUSNAHME"
! printf '%s\n' "$datenbank" | awk '/^GATE_SCHNELL_AUSNAHME :?= / { for (i = 3; i <= NF; i++) if ($i == "meldungskatalog-check") t = 1 } END { exit !t }' \
  || melde "Fall 'nr11-schnelle-kette': meldungskatalog-check steht in GATE_SCHNELL_AUSNAHME"
make -s -C "$gruen_baum" meldungskatalog-check >/dev/null 2>&1 || melde "Fall 'nr11-ziel-gruen': make meldungskatalog-check endet im grünen Baum ungleich 0"
vorher="$(bestand "$w")"
for flag in -n -q -t -i; do
  make -s -C "$w" "$flag" meldungskatalog-check >/dev/null 2>&1 || true
  [ "$vorher" = "$(bestand "$w")" ] || melde "Fall 'nr11-make$flag': make $flag hat den Baum verändert"
done
MAKEFLAGS=-n make -s -C "$w" meldungskatalog-check >/dev/null 2>&1 || true
[ "$vorher" = "$(bestand "$w")" ] || melde "Fall 'nr11-makeflags-umgebung': MAKEFLAGS=-n hat den Baum verändert"

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "meldungskatalog-check-gegenprobe: gruen — Fehlformen je Punkt von SPEC-054 abgelehnt, Gutfälle angenommen"
