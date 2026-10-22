#!/usr/bin/env bash
# archiv-bezug-gegenprobe — prüft die Zusagen von SPEC-068 (harness/vertraege/SPEC-068-archiv-bezug.md, Punkte 1 bis 12) an tools/harness/archiv-bezug.sh und
# harness/mk/archiv-bezug.mk in Temp-Bäumen mit selbst gebauten Archiven (zip), ohne Docker. Ein Fall besteht, wenn alle seine Prüfungen (Exit-Code, Meldung, Inhalt und
# Zustand der Stubs) halten; sein Name trägt die Nummer des Punktes. Je Mutation (Textersatz in einer Kopie von Skript oder Fragment) müssen genau die benannten Fälle rot werden; die unveränderte Kopie besteht jeden Fall. Jeder Fall ist in mindestens einer Liste roter Fälle genannt.
# Zum Arbeiten: GP_MUT=<Muster> fährt nur Läufe, deren Name passt (die unveränderte Kopie nur, wenn das Muster `ok` trifft); GP_ENTDECKE=1 nennt je Mutation die roten
# Fälle; GP_PARALLEL=<n> (Vorgabe 8) setzt die Zahl der gleichzeitigen Läufe; GP_FALL=<name> gibt
# für die unveränderte Kopie das Protokoll dieses Falls aus. Als root entfallen die Fälle mit Dateirechten (Liste ENTFALLEN im Skript) und ihre Mutationen laut, ebenso ohne die Locale en_US.UTF-8 der Fall p8-sprache.
# Braucht bash ≥ 4.3 (`wait -n`), make, zip, unzip, GNU sed, grep, awk, sort, cmp, stat, find, md5sum und die Werkzeuge von coreutils. Kein Fall ändert Dateirechte außerhalb seines Baums.
# Nicht geprüft: gleichzeitige Aufrufe (Punkt 9), Abbruch durch Signale (kein Teil des Vertrags), ein Lauf gegen d-check (Punkt 11 (b)).
set -uo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="$(cd "$hier/../../.." && pwd)"
arbeit="$(mktemp -d)"
raeume_auf() { chmod -R u+rwX "${arbeit:?}" 2> /dev/null; rm -rf "${arbeit:?}"; }
trap raeume_auf EXIT
trap 'kill $(jobs -p) 2> /dev/null; exit 130' INT
trap 'kill $(jobs -p) 2> /dev/null; exit 143' TERM
: > "$arbeit/fehler"; : > "$arbeit/gefahren"; : > "$arbeit/rote"
melde() { echo "archiv-bezug-gegenprobe: ROT — $*" >> "$arbeit/fehler"; }

for t in zip unzip sed grep awk sort cmp stat find md5sum make mktemp ln cp chmod touch tr ls; do
  command -v "$t" > /dev/null 2>&1 || { echo "archiv-bezug-gegenprobe: $t fehlt" >&2; exit 2; }
done
for f in tools/harness/archiv-bezug.sh harness/mk/archiv-bezug.mk; do
  [ -f "$repo/$f" ] || { echo "archiv-bezug-gegenprobe: $f fehlt im Arbeitsbaum" >&2; exit 2; }
done
ENTFALLEN=(); ENT_MUT=()
if [ "$(id -u)" -eq 0 ]; then
  ENTFALLEN=(p3-datei-nicht-beschreibbar p3-verzeichnis-nicht-beschreibbar p3-nicht-beschreibbar-gleich p3-erster-grund p3-reihenfolge-grund-d-e); ENT_MUT+=(nicht-beschreibbar)
  echo "archiv-bezug-gegenprobe: als root entfallen ${ENTFALLEN[*]} und ihre Mutationen" >&2
fi
LOC=en_US.UTF-8
if ! locale -a 2> /dev/null | grep -qix 'en_US.utf-\?8'; then
  ENTFALLEN+=(p8-sprache); ENT_MUT+=(sprache); echo "archiv-bezug-gegenprobe: ohne die Locale en_US.UTF-8 entfallen p8-sprache und seine Mutation" >&2
fi

# --- Bausteine --------------------------------------------------------------------------------------------------------------------------------------
DONE=docs/plan/planning/done
STUB=$'# slice- — NAME: Titel\n\n> **ARCHIVIERT** — Volltext: `unzip -p x`\n\n**Welle:** [w](w.md)\n**Hervorgegangen:** [`BEO-REPO/a`](b.md) · ADR-0001'
# bin_ohne_unzip — ein PATH-Verzeichnis mit den Werkzeugen des Skripts, aber ohne unzip.
mkdir -p "$arbeit/bin"
for t in mktemp rm basename awk grep sort cut mv; do ln -s "$(command -v "$t")" "$arbeit/bin/$t"; done
baum() { w="$run/t.$fallname"; rm -rf "$w" "$w.s"; mkdir -p "$w/$DONE" "$w.s/$DONE"; vor=""; LAUFDIR=""; P="$PATH"; ME=(); LC=C; }
# sv <welle> <name> [<Zeilen des Kopfs>] [<Stub>] — Stub im Baum, Volltext (Kopf, SCHLUSS = Zeile nach dem Kopf, ---, Rumpf mit LH-RB-77) in der Bühne.
sv() {
  mkdir -p "$w/$DONE/$1" "$w.s/$DONE/$1"
  printf '%s\n' "${4-$STUB}" > "$w/$DONE/$1/$2"
  printf '# Slice %s\n\n%b\n%b---\n## 1. Ziel\nRumpf LH-RB-77\n' "$2" "${3-**Bezug:** LH-FA-01}" "${SCHLUSS-**Berührte Spec-Stellen:** SPEC-001\n}" > "$w.s/$DONE/$1/$2"
}
fertig() {
  local d n
  for d in "$w.s/$DONE"/*/; do n="$(basename "$d")"; (cd "$w.s" && zip -q -r ${ZIPOPT-} "$w/$DONE/$n/archiv.zip" "$DONE/$n"); done
  vor="$(snap)"
}
snap() { find "$w" -type f -exec stat -c '%n %i %s %.9Y %a' {} + 2> /dev/null | LC_ALL=C sort | md5sum; }
lauf() {
  ( cd "${LAUFDIR:-$w}" && env LC_ALL="${LC:-C}" PATH="$P" "$BASH" "$skript" "$w" "$@" ) > "$w.out" 2> "$w.err"
  code=$?; err="$(cat "$w.err")"
}
ok() { if ! "$@"; then bad=1; printf 'FAIL: %s\nExit %s, stderr:\n%s\n' "$*" "${code-}" "${err-}" >> "$run/log.$fallname"; fi; }
rc() { [ "$code" -eq "$1" ]; }
err_hat() { [[ $err == *"$1"* ]]; }
out_hat() { [[ "$(cat "$w.out")" == *"$1"* ]]; }
# befund_hat <welle> <stub> <Stichwort> — eine Zeile `archiv-bezug: <Pfad des Stubs>: <Stichwort>` (SPEC-068 Punkt 2, Meldungen).
befund_hat() { err_hat "archiv-bezug: $DONE/$1/$2: $3"; }
# einmal <welle> <stub> — der Stub kommt in der Fehlerausgabe genau einmal vor.
einmal() { [ "$(grep -c "$DONE/$1/$2" "$w.err")" -eq 1 ]; }
unveraendert() { [ "$(snap)" = "$vor" ]; }
inhalt() { [ "$(cat "$1"; printf x)" = "$2"$'\n'x ]; }
# bezug <welle> <name> <Text> — der Stub ist der Standard-Stub mit der Zeile **Bezug:** <Text> dahinter.
bezug() { inhalt "$w/$DONE/$1/$2" "$STUB"$'\n'"**Bezug:** $3"; }
# eins <Kopf> <Text> — eine Welle, ein Stub; Exit 0 und die Zeile mit <Text>.
eins() { baum; sv welle-a slice-a.md "$1"; fertig; lauf; ok rc 0; ok bezug welle-a slice-a.md "$2"; }
stubdatei() { printf '%b' "$3" > "$w/$DONE/$1/$2"; }
mbaum() {
  baum; mkdir -p "$w/harness/mk" "$w/tools/harness"
  cp "$fragment" "$w/harness/mk/archiv-bezug.mk"; cp "$skript" "$w/tools/harness/archiv-bezug.sh"
  printf 'include harness/mk/archiv-bezug.mk\n' > "$w/Makefile"
}
mk() { ( cd "$w" && env -u MAKEFLAGS -u MFLAGS -u MAKELEVEL ${ME[@]+"${ME[@]}"} make "$@" ) > "$w.out" 2> "$w.err"; code=$?; err="$(cat "$w.err")"; }
lauf_roh() { ( cd "$w" && env PATH="$P" "$BASH" "$skript" "$@" ) > "$w.out" 2> "$w.err"; code=$?; err="$(cat "$w.err")"; }
rest_snap() { find "$w" -type f ! -name slice-a.md -exec stat -c '%n %i %s %.9Y %a' {} + | LC_ALL=C sort | md5sum; }
# Attrappe von mv: scheitert für slice-b.md, sonst das echte mv.
mkdir -p "$arbeit/att"
printf '#!/bin/bash\ncase "$*" in *slice-b.md*) exit 1 ;; esac\nexec %s "$@"\n' "$(command -v mv)" > "$arbeit/att/mv"; chmod +x "$arbeit/att/mv"
nicht_root() { if [ "$(id -u)" -eq 0 ]; then exit 0; fi; }

# --- Punkt 1: Lauf-Menge ----------------------------------------------------------------------------------------------------------------------------
f_p1_alle() { baum; sv welle-a slice-a.md; sv welle-b slice-b.md; fertig; lauf; ok rc 0; ok eval '! out_hat "keine Stubs"'; ok bezug welle-a slice-a.md LH-FA-01; ok bezug welle-b slice-b.md LH-FA-01; }
f_p1_eine() { baum; sv welle-a slice-a.md; sv welle-b slice-b.md; fertig; lauf welle-a; ok rc 0; ok bezug welle-a slice-a.md LH-FA-01; ok inhalt "$w/$DONE/welle-b/slice-b.md" "$STUB"; }
f_p1_ohne_stub() { baum; sv welle-a slice-a.md; mkdir -p "$w/$DONE/welle-leer"; fertig; lauf; ok rc 0; ok bezug welle-a slice-a.md LH-FA-01; }
f_p1_leere_menge() { baum; mkdir -p "$w/$DONE/welle-leer"; vor="$(snap)"; lauf; ok rc 0; ok out_hat "archiv-bezug: keine Stubs"; ok unveraendert; }
f_p1_eine_ohne_stub() { baum; sv welle-leer slice-x.md; fertig; rm "$w/$DONE/welle-leer/slice-x.md"; vor="$(snap)"; lauf welle-leer; ok rc 0; ok out_hat "archiv-bezug: keine Stubs"; ok unveraendert; }
f_p1_welle_flach() {
  baum; sv welle-a slice-a.md; printf 'flach\n' > "$w/$DONE/welle-x-results.md"; printf 'flach\n' > "$w/$DONE/slice-flach.md"; fertig; lauf
  ok rc 0; ok bezug welle-a slice-a.md LH-FA-01; ok inhalt "$w/$DONE/welle-x-results.md" flach; ok inhalt "$w/$DONE/slice-flach.md" flach
}
f_p1_welle_datei() { baum; sv welle-a slice-a.md; printf 'welle\n' > "$w/$DONE/welle-a/welle-a.md"; fertig; lauf; ok rc 0; ok inhalt "$w/$DONE/welle-a/welle-a.md" welle; }
f_p1_make_welle() { mbaum; sv welle-a slice-a.md; sv welle-b slice-b.md; fertig; mk archiv-bezug WELLE=welle-a; ok rc 0; ok bezug welle-a slice-a.md LH-FA-01; ok inhalt "$w/$DONE/welle-b/slice-b.md" "$STUB"; }

# --- Punkt 2: Eingaben und Reihenfolge der Gründe ---------------------------------------------------------------------------------------------------
form() { baum; lauf "$1"; ok rc 2; ok err_hat "archiv-bezug: Form der Welle ungültig: $1"; }
f_p2_form_welle_praefix() { form foo-bar; }
f_p2_form_welle_zeichen() { form welle-Ab; }
f_p2_form_welle_slash() { form welle-a/b; }
f_p2_form_welle_punkt() { form welle-a..b; }
f_p2_form_welle_nur_praefix() { form welle-; }
f_p2_welle_leer() { baum; sv welle-a slice-a.md; fertig; lauf ""; ok rc 0; ok bezug welle-a slice-a.md LH-FA-01; }
f_p2_wurzel_fehlt() { baum; lauf_roh; ok rc 2; ok err_hat "archiv-bezug: Wurzel fehlt"; }
f_p2_wurzel_leer() { baum; lauf_roh ""; ok rc 2; ok err_hat "archiv-bezug: Wurzel fehlt"; }
f_p2_wurzel_kein_verzeichnis() { baum; lauf_roh "$w/nope"; ok rc 2; ok err_hat "archiv-bezug: Wurzel fehlt oder ist kein Verzeichnis: $w/nope"; }
f_p2_reihenfolge_0_1() { baum; lauf_roh "$w/nope" welle-Bad; ok rc 2; ok err_hat "archiv-bezug: Wurzel fehlt"; }
f_p2_welle_ohne_stub_ohne_archiv() { baum; mkdir -p "$w/$DONE/welle-leer"; vor="$(snap)"; lauf welle-leer; ok rc 0; ok out_hat "archiv-bezug: keine Stubs"; ok unveraendert; }
f_p2_welle_unbekannt() { baum; lauf welle-nope; ok rc 2; ok err_hat "archiv-bezug: Verzeichnis fehlt: $DONE/welle-nope"; }
f_p2_kein_unzip() { baum; sv welle-a slice-a.md; fertig; P="$arbeit/bin"; lauf; ok rc 2; ok err_hat "archiv-bezug: Werkzeug fehlt: unzip"; ok unveraendert; }
f_p2_archiv_fehlt() { baum; sv welle-a slice-a.md; vor="$(snap)"; lauf; ok rc 2; ok err_hat "archiv-bezug: archiv.zip fehlt: welle-a"; ok unveraendert; }
f_p2_archiv_beschaedigt() { baum; sv welle-a slice-a.md; printf 'kein zip' > "$w/$DONE/welle-a/archiv.zip"; vor="$(snap)"; lauf; ok rc 2; ok err_hat "archiv-bezug: archiv.zip beschädigt: welle-a/archiv.zip"; ok unveraendert; }
f_p2_archiv_crc() {
  baum; sv welle-a slice-a.md; ZIPOPT=-0 fertig; LC_ALL=C sed -i 's/Rumpf/Rumpg/' "$w/$DONE/welle-a/archiv.zip"; vor="$(snap)"; lauf
  ok rc 2; ok err_hat "archiv-bezug: archiv.zip beschädigt: welle-a/archiv.zip"; ok unveraendert
}
f_p2_reihenfolge_1_2() { baum; lauf welle-Bad; ok rc 2; ok err_hat "archiv-bezug: Form der Welle ungültig: welle-Bad"; }
f_p2_reihenfolge_2_3() { baum; P="$arbeit/bin"; lauf welle-nope; ok rc 2; ok err_hat "archiv-bezug: Verzeichnis fehlt: $DONE/welle-nope"; }
f_p2_reihenfolge_3_4() { baum; sv welle-a slice-a.md; P="$arbeit/bin"; lauf; ok rc 2; ok err_hat "archiv-bezug: Werkzeug fehlt: unzip"; }
f_p2_reihenfolge_4_5() {
  baum; sv welle-a slice-a.md; sv welle-b slice-b.md; fertig; printf 'kein zip' > "$w/$DONE/welle-a/archiv.zip"; rm "$w/$DONE/welle-b/archiv.zip"; lauf
  ok rc 2; ok err_hat "archiv-bezug: archiv.zip fehlt: welle-b"
}
f_p2_keine_teilschreibung() { baum; sv welle-a slice-a.md; sv welle-b slice-b.md; fertig; rm "$w/$DONE/welle-b/archiv.zip"; vor="$(snap)"; lauf; ok rc 2; ok unveraendert; }

# --- Punkt 3: Prüfung aller Stubs vor dem ersten Schreiben ------------------------------------------------------------------------------------------
f_p3_stub_ohne_eintrag() {
  baum; sv welle-a slice-a.md; sv welle-a slice-b.md; sv welle-a slice-z.md; rm "$w.s/$DONE/welle-a/slice-z.md"; fertig; lauf
  ok rc 1; ok befund_hat welle-a slice-z.md "kein Eintrag im Archiv"; ok unveraendert
}
f_p3_alle_befunde() {
  baum; sv welle-a slice-a.md; sv welle-a slice-b.md; sv welle-a slice-c.md; rm "$w.s/$DONE/welle-a/slice-a.md"; stubdatei welle-a slice-b.md "$STUB\r\n"; fertig; lauf
  ok rc 1; ok befund_hat welle-a slice-a.md "kein Eintrag im Archiv"; ok befund_hat welle-a slice-b.md "Stub enthält CR"; ok unveraendert
}
f_p3_datei_nicht_beschreibbar() { nicht_root; baum; sv welle-a slice-a.md; sv welle-a slice-b.md; fertig; chmod a-w "$w/$DONE/welle-a/slice-b.md"; vor="$(snap)"; lauf; ok rc 1; ok befund_hat welle-a slice-b.md "nicht beschreibbar"; ok unveraendert; }
f_p3_verzeichnis_nicht_beschreibbar() { nicht_root; baum; sv welle-a slice-a.md; fertig; chmod a-w "$w/$DONE/welle-a"; vor="$(snap)"; lauf; ok rc 1; ok befund_hat welle-a slice-a.md "nicht beschreibbar"; ok unveraendert; }
f_p3_nicht_beschreibbar_gleich() { nicht_root; baum; sv welle-a slice-a.md; stubdatei welle-a slice-a.md "$STUB\n**Bezug:** LH-FA-01\n"; chmod a-w "$w/$DONE/welle-a/slice-a.md"; fertig; lauf; ok rc 0; ok unveraendert; }
f_p3_erster_grund() { nicht_root; baum; sv welle-a slice-a.md; fertig; rm "$w.s/$DONE/welle-a/slice-a.md"; zip -q -d "$w/$DONE/welle-a/archiv.zip" "$DONE/welle-a/slice-a.md"; chmod a-w "$w/$DONE/welle-a/slice-a.md"; vor="$(snap)"; lauf; ok rc 1; ok befund_hat welle-a slice-a.md "kein Eintrag im Archiv"; ok einmal welle-a slice-a.md; ok unveraendert; }
f_p3_reihenfolge_grund_a_b() { baum; sv welle-a slice-a.md; rm "$w.s/$DONE/welle-a/slice-a.md"; fertig; lauf; ok rc 1; ok befund_hat welle-a slice-a.md "kein Eintrag im Archiv"; ok einmal welle-a slice-a.md; ok unveraendert; }
f_p3_reihenfolge_grund_b_c() { baum; sv welle-a slice-a.md; printf '# x\nText\n' > "$w.s/$DONE/welle-a/slice-a.md"; stubdatei welle-a slice-a.md "$STUB\r\n"; fertig; lauf; ok rc 1; ok befund_hat welle-a slice-a.md "Eintrag im Archiv ohne Kopf"; ok einmal welle-a slice-a.md; ok unveraendert; }
f_p3_reihenfolge_grund_c_c() { mit_stub '# Stub\r\nText\r\n'; ok rc 1; ok befund_hat welle-a slice-a.md "Stub enthält CR"; ok einmal welle-a slice-a.md; ok unveraendert; }
f_p3_reihenfolge_grund_c_d() { mit_stub "$STUB\r\n**Bezug:** LH-FA-01\r\n**Bezug:** LH-FA-01\r\n"; ok rc 1; ok befund_hat welle-a slice-a.md "Stub enthält CR"; ok einmal welle-a slice-a.md; ok unveraendert; }
f_p3_reihenfolge_grund_d_e() { nicht_root; baum; sv welle-a slice-a.md; stubdatei welle-a slice-a.md "$STUB\n**Bezug:** LH-FA-99\n"; fertig; chmod a-w "$w/$DONE/welle-a/slice-a.md"; vor="$(snap)"; lauf; ok rc 1; ok befund_hat welle-a slice-a.md "Zeile weicht ab"; ok einmal welle-a slice-a.md; ok unveraendert; }
f_p3_schreiben_scheitert() {
  baum; sv welle-a slice-a.md; sv welle-a slice-b.md; sv welle-a slice-c.md; fertig; P="$arbeit/att:$PATH"; lauf
  ok rc 2; ok befund_hat welle-a slice-b.md "Schreiben fehlgeschlagen"; ok bezug welle-a slice-a.md LH-FA-01; ok inhalt "$w/$DONE/welle-a/slice-c.md" "$STUB"
}
f_p3_schreiben_temp_aufgeraeumt() {
  baum; sv welle-a slice-a.md; sv welle-a slice-b.md; fertig; P="$arbeit/att:$PATH"; lauf
  ok eval "[ \"\$(LC_ALL=C ls -A '$w/$DONE/welle-a' | tr '\\n' ' ')\" = 'archiv.zip slice-a.md slice-b.md ' ]"
}
f_p3_nur_stubs() {
  local r
  baum; sv welle-a slice-a.md; sv welle-a slice-b.md; stubdatei welle-a slice-b.md "$STUB\n**Bezug:** LH-FA-01\n"; fertig
  touch -d 2020-01-01 "$w/$DONE/welle-a/archiv.zip" "$w/$DONE/welle-a/slice-b.md"; r="$(rest_snap)"; lauf
  ok rc 0; ok bezug welle-a slice-a.md LH-FA-01; ok eval '[ "$(rest_snap)" = "$r" ]'
}
f_p3_fast_gleicher_pfad() { baum; sv welle-a slice-a.md; mv "$w.s/$DONE/welle-a/slice-a.md" "$w.s/$DONE/welle-a/slice-a.md.orig"; fertig; lauf; ok rc 1; ok befund_hat welle-a slice-a.md "kein Eintrag im Archiv"; ok einmal welle-a slice-a.md; ok unveraendert; }
f_p3_kein_kopf() { baum; sv welle-a slice-a.md; printf '# x\n**Bezug:** LH-FA-01\nText\n' > "$w.s/$DONE/welle-a/slice-a.md"; fertig; lauf; ok rc 1; ok befund_hat welle-a slice-a.md "Eintrag im Archiv ohne Kopf"; ok unveraendert; }
f_p3_reihenfolge_aufruf() {
  baum; sv welle-a slice-a.md; sv welle-b slice-b.md; rm "$w.s/$DONE/welle-a/slice-a.md"; fertig; rm "$w/$DONE/welle-b/archiv.zip"; lauf
  ok rc 2; ok err_hat "archiv-bezug: archiv.zip fehlt: welle-b"; ok eval '! err_hat "kein Eintrag"'
}
f_p3_temp_und_mv() {
  local ino
  baum; sv welle-a slice-a.md; fertig; ino="$(stat -c %i "$w/$DONE/welle-a/slice-a.md")"; lauf
  ok rc 0; ok bezug welle-a slice-a.md LH-FA-01
  ok eval "[ \"\$(stat -c %i '$w/$DONE/welle-a/slice-a.md')\" != '$ino' ]"
  ok eval "[ \"\$(LC_ALL=C ls -A '$w/$DONE/welle-a' | tr '\\n' ' ')\" = 'archiv.zip slice-a.md ' ]"
}

# --- Punkt 4: Kopf und gelesene Kennungen -----------------------------------------------------------------------------------------------------------
f_p4_ein_token() { eins '**Bezug:** LH-FA-01' LH-FA-01; }
f_p4_mehrzeilig() { eins '**Bezug:** LH-FA-01\nFortsetzung LH-QA-02 und\n\nLH-RB-03\n**Berührte Spec-Stellen:** LH-FA-77' 'LH-FA-01, LH-QA-02, LH-RB-03'; }
f_p4_ende_kopf() { SCHLUSS= eins '**Bezug:** LH-FA-01\nweiter LH-QA-02' 'LH-FA-01, LH-QA-02'; }
f_p4_kopfende_genau() { SCHLUSS= eins '**Bezug:** LH-FA-01\n--- kein Ende\n-----\nLH-QA-02' 'LH-FA-01, LH-QA-02'; }
f_p4_suffix_punkt() { eins '**Bezug:** LH-FA-05.e' LH-FA-05; }
f_p4_suffix_buchstabe() { eins '**Bezug:** LH-QA-03a' LH-QA-03; }
f_p4_kurzform() { eins '**Bezug:** LH-QA-03 und -04' LH-QA-03; }
f_p4_anker_klein() { eins '**Bezug:** [Lastenheft](../spec/lastenheft.md#lh-fa-14--diagnoseausgaben)' '— keine —'; }
f_p4_vorn_buchstabe() { eins '**Bezug:** aLH-FA-05 ZLH-FA-06 (LH-QA-04)' LH-QA-04; }
f_p4_vorn_ziffer() { eins '**Bezug:** 5LH-FA-05 (LH-QA-04)' LH-QA-04; }
f_p4_vorn_strich() { eins '**Bezug:** -LH-FA-05 (LH-QA-04)' LH-QA-04; }
f_p4_kein_bezug() { eins '**Welle:** x' '— keine —'; }
f_p4_kein_token() { eins '**Bezug:** nur Text' '— keine —'; }
f_p4_leer_wert() { eins '**Welle:** x' '— keine —'; vor="$(snap)"; lauf; ok rc 0; ok unveraendert; }
f_p4_zwei_bezug() { eins '**Bezug:** LH-FA-01\n**Berührte Spec-Stellen:** SPEC-001\n**Bezug:** LH-QA-02' LH-FA-01; }

# --- Punkt 5: die Zeile -----------------------------------------------------------------------------------------------------------------------------
f_p5_ordnung() { eins '**Bezug:** LH-RB-02, LH-QA-03, LH-FA-05, LH-RB-01, LH-QA-01' 'LH-FA-05, LH-QA-01, LH-QA-03, LH-RB-01, LH-RB-02'; }
f_p5_zahl() { eins '**Bezug:** LH-FA-10 LH-FA-9 LH-FA-100' 'LH-FA-9, LH-FA-10, LH-FA-100'; }
f_p5_gleiche_zahl() { eins '**Bezug:** LH-FA-04 LH-FA-4' 'LH-FA-4, LH-FA-04'; }
f_p5_duplikat() { eins '**Bezug:** LH-FA-01, LH-FA-01 und LH-FA-01.e' LH-FA-01; }

# --- Punkt 6: vorhandene Zeile, Idempotenz ----------------------------------------------------------------------------------------------------------
mit_zeilen() { baum; sv welle-a slice-a.md; stubdatei welle-a slice-a.md "$STUB\n$1\n"; touch -d 2020-01-01 "$w/$DONE/welle-a/slice-a.md"; fertig; lauf; }
f_p6_gleich() { mit_zeilen '**Bezug:** LH-FA-01'; ok rc 0; ok unveraendert; }
f_p6_abweichend() { mit_zeilen '**Bezug:** LH-FA-99'; ok rc 1; ok befund_hat welle-a slice-a.md "Zeile weicht ab"; ok err_hat "LH-FA-99"; ok err_hat "LH-FA-01"; ok unveraendert; }
f_p6_gleich_falsche_stelle() { baum; sv welle-a slice-a.md; stubdatei welle-a slice-a.md '# Stub\n**Bezug:** LH-FA-01\n**Hervorgegangen:** x\nText\n'; fertig; lauf; ok rc 1; ok befund_hat welle-a slice-a.md "Zeile an falscher Stelle"; ok unveraendert; }
f_p6_zwei_zeilen() { mit_zeilen '**Bezug:** LH-FA-01\n**Bezug:** LH-FA-01'; ok rc 1; ok befund_hat welle-a slice-a.md "mehr als eine Zeile **Bezug:**"; ok unveraendert; }
f_p6_zweiter_lauf() { eins '**Bezug:** LH-FA-01' LH-FA-01; vor="$(snap)"; lauf; ok rc 0; ok unveraendert; }

# --- Punkt 7: Stelle der Zeile und Form des Stubs ---------------------------------------------------------------------------------------------------
mit_stub() { baum; sv welle-a slice-a.md; stubdatei welle-a slice-a.md "$1"; fertig; lauf; }
f_p7_ende() { eins '**Bezug:** LH-FA-01' LH-FA-01; }
f_p7_mitte() { mit_stub "$STUB\nWeitere Zeile\n**Letzte:** x\n"; ok rc 0; ok inhalt "$w/$DONE/welle-a/slice-a.md" "$STUB"$'\n**Bezug:** LH-FA-01\nWeitere Zeile\n**Letzte:** x'; }
f_p7_ohne_umbruch() { mit_stub "$STUB"; ok rc 0; ok bezug welle-a slice-a.md LH-FA-01; }
f_p7_ohne_umbruch_weitere() { mit_stub "$STUB\nrest"; ok rc 0; ok inhalt "$w/$DONE/welle-a/slice-a.md" "$STUB"$'\n**Bezug:** LH-FA-01\nrest'; }
f_p7_ohne_hervorgegangen() { mit_stub '# Stub\n\nText\n'; ok rc 1; ok befund_hat welle-a slice-a.md "Stub hat nicht genau eine Zeile **Hervorgegangen:**"; ok unveraendert; }
f_p7_zweimal_hervorgegangen() { mit_stub "$STUB\n**Hervorgegangen:** x\n"; ok rc 1; ok befund_hat welle-a slice-a.md "Stub hat nicht genau eine Zeile **Hervorgegangen:**"; ok unveraendert; }
f_p7_rest_unveraendert() {
  local o="$run/orig.$fallname"
  printf '%b' "$STUB \t\n\n   \nZeile mit Leerzeichen am Ende  \n" > "$o"
  mit_stub "$STUB \t\n\n   \nZeile mit Leerzeichen am Ende  \n"; ok rc 0
  ok eval "LC_ALL=C grep -av '^\\*\\*Bezug:\\*\\*' '$w/$DONE/welle-a/slice-a.md' | cmp -s - '$o'"
}

# --- Punkt 8: Zeilenende und Kodierung --------------------------------------------------------------------------------------------------------------
f_p8_stub_crlf() { mit_stub "$STUB\r\n"; ok rc 1; ok befund_hat welle-a slice-a.md "Stub enthält CR"; ok unveraendert; }
f_p8_volltext_crlf() { baum; sv welle-a slice-a.md '**Bezug:** LH-FA-01, LH-QA-02'; sed -i 's/$/\r/' "$w.s/$DONE/welle-a/slice-a.md"; fertig; lauf; ok rc 0; ok bezug welle-a slice-a.md 'LH-FA-01, LH-QA-02'; }
f_p8_umlaute() {
  local o="$run/orig.$fallname"
  printf '%b' "$STUB\nGröße äöü \xff\n" > "$o"; mit_stub "$STUB\nGröße äöü \xff\n"; ok rc 0
  ok eval "LC_ALL=C grep -av '^\\*\\*Bezug:\\*\\*' '$w/$DONE/welle-a/slice-a.md' | cmp -s - '$o'"; ok eval "LC_ALL=C grep -aq '^\\*\\*Bezug:\\*\\* LH-FA-01\$' '$w/$DONE/welle-a/slice-a.md'"
}
f_p8_sprache() { baum; LC="$LOC" lauf welle-ä; ok rc 2; ok err_hat "archiv-bezug: Form der Welle ungültig: welle-ä"; }

# --- Punkt 9 und 10 ---------------------------------------------------------------------------------------------------------------------------------
f_p9_wurzel() { baum; sv welle-a slice-a.md; fertig; LAUFDIR="$w/docs"; lauf; ok rc 0; ok bezug welle-a slice-a.md LH-FA-01; }
f_p10_n() { mbaum; sv welle-a slice-a.md; fertig; mk -n archiv-bezug; ok rc 0; ok eval '[[ $(cat "$w.out") == *archiv-bezug.sh* ]]'; ok unveraendert; }
f_p10_q() { mbaum; sv welle-a slice-a.md; fertig; mk -q archiv-bezug; ok rc 1; ok unveraendert; }
f_p10_t() { mbaum; sv welle-a slice-a.md; fertig; mk -t archiv-bezug; ok rc 0; ok unveraendert; }
f_p10_i() { mbaum; sv welle-a slice-a.md; sv welle-a slice-z.md; rm "$w.s/$DONE/welle-a/slice-z.md"; fertig; mk -i archiv-bezug; ok rc 0; ok befund_hat welle-a slice-z.md "kein Eintrag im Archiv"; ok unveraendert; }
f_p10_makeflags() { mbaum; sv welle-a slice-a.md; fertig; ME=(MAKEFLAGS=n); mk archiv-bezug; ok rc 0; ok unveraendert; }

# --- Punkt 11: kein Gate, an keiner Kette ------------------------------------------------------------------------------------------------------------
f_p11_keine_kette() {
  mbaum; printf 'GATE_CHECKS :=\nGEGENPROBEN_CONTAINER :=\nGEGENPROBEN_HOST :=\ninclude harness/mk/archiv-bezug.mk\nzeige: ; @echo "[$(GATE_CHECKS)][$(GEGENPROBEN_CONTAINER)][$(GEGENPROBEN_HOST)]"\n' > "$w/Makefile"
  mk zeige; ok rc 0; ok eval '[ "$(cat "$w.out")" = "[][][]" ]'
}

# --- Punkt 12: nur Anforderungen --------------------------------------------------------------------------------------------------------------------
f_p12_adr() { eins '**Bezug:** [ADR-0011](x) LH-FA-01' LH-FA-01; }
f_p12_spec() { eins '**Bezug:** SPEC-034 LH-FA-01' LH-FA-01; }
f_p12_arc() { eins '**Bezug:** ARC-001 LH-FA-01' LH-FA-01; }

# --- Läufe ------------------------------------------------------------------------------------------------------------------------------------------
mapfile -t FAELLE < <(declare -F | awk '$3 ~ /^f_p[0-9]/ { print substr($3, 3) }' | tr _ - | LC_ALL=C sort)
# faelle <Gruppe> — schreibt die Namen der roten Fälle der Gruppe (p<n>, oder alle), einen je Zeile.
faelle() {
  local n e
  for n in "${FAELLE[@]}"; do
    [[ $1 == alle || $n == "$1"-* ]] || continue
    for e in ${ENTFALLEN[@]+"${ENTFALLEN[@]}"}; do [ "$e" != "$n" ] || continue 2; done
    fallname="$n"; mkdir -p "$run"
    ( bad=0; "f_${n//-/_}"; exit "$bad" ) > /dev/null 2>> "$run/stderr" || echo "$n"
  done
}
maxjobs="${GP_PARALLEL:-8}"
wartejob() { while [ "$(jobs -rp | wc -l)" -ge "$maxjobs" ]; do wait -n || true; done; }
# bewerte <name> <rote Fälle> <Gruppe oder alle>
bewerte() {
  local ist erwartet
  ist="$(faelle "$3" | LC_ALL=C sort -u | tr '\n' ' ')"
  erwartet="$(printf '%s\n' $2 | grep -vxF -f <(printf '%s\n' x ${ENTFALLEN[@]+"${ENTFALLEN[@]}"}) | LC_ALL=C sort -u | tr '\n' ' ')"
  if [ -n "${GP_ENTDECKE-}" ]; then echo "$1 => ${ist:-nichts}" >&2; fi
  echo "$1" >> "$arbeit/gefahren"; printf '%s\n' $2 >> "$arbeit/rote"
  [ "$ist" = "$erwartet" ] || melde "Mutation '$1': erwartet rot: ${2:-nichts}; rot sind: ${ist:-nichts}"
}
if [[ ok =~ ${GP_MUT:-.} ]]; then
  (
    run="$arbeit/ok"; mkdir -p "$run"; skript="$repo/tools/harness/archiv-bezug.sh"; fragment="$repo/harness/mk/archiv-bezug.mk"
    ist="$(faelle alle | LC_ALL=C sort -u | tr '\n' ' ')"
    echo ok >> "$arbeit/gefahren"
    [ -z "$ist" ] || melde "unveränderte Kopie ist nicht grün: $ist"
    if [ -n "${GP_FALL-}" ] && [ -f "$run/log.$GP_FALL" ]; then cat "$run/log.$GP_FALL" >&2; fi
  ) &
fi
# mut <name> <rote Fälle> <skript|fragment> <alt> <neu> [<alt> <neu>…] — Textersatz in einer Kopie; rot sind genau die benannten Fälle (alle Fälle laufen).
mut() {
  local name="$1" rote="$2" ziel="$3" quelle inhalt neu alt rest g
  shift 3
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  for alt in ${ENT_MUT[@]+"${ENT_MUT[@]}"}; do [[ $name != *"$alt"* ]] || return 0; done
  wartejob
  (
    run="$arbeit/$name"; mkdir -p "$run"; skript="$repo/tools/harness/archiv-bezug.sh"; fragment="$repo/harness/mk/archiv-bezug.mk"
    if [ "$ziel" = skript ]; then quelle="$skript"; else quelle="$fragment"; fi
    inhalt="$(cat "$quelle"; printf x)"; inhalt="${inhalt%x}"; neu="$inhalt"
    while [ "$#" -ge 2 ]; do
      alt="$1"
      if [[ $neu != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
      rest="${neu#*"$alt"}"
      if [[ $rest == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
      neu="${neu/"$alt"/"$2"}"
      shift 2
    done
    if [[ $neu == "$inhalt" ]]; then melde "Mutation '$name': Ersatz ändert nichts"; exit 0; fi
    printf '%s' "$neu" > "$run/mutant"
    if [ "$ziel" = skript ]; then skript="$run/mutant"; else fragment="$run/mutant"; fi
    bewerte "$name" "$rote" alle
  ) &
}

LF=$'\n'
L_FORM='[[ $welle =~ ^welle-[a-z0-9-]+$ ]] || sperre "Form der Welle ungültig: $welle"'
L_DIR='[ -d "$done_dir/$welle" ] || sperre "Verzeichnis fehlt: $rel/$welle"'
L_UNZIP='command -v unzip > /dev/null 2>&1 || sperre "Werkzeug fehlt: unzip"'
L_FEHLT='for w in "${wellen[@]}"; do [ -f "$done_dir/$w/archiv.zip" ] || sperre "archiv.zip fehlt: $w"; done'
L_BESCH='for w in "${wellen[@]}"; do unzip -tq "$done_dir/$w/archiv.zip" > /dev/null 2>&1 || sperre "archiv.zip beschädigt: $w/archiv.zip"; done'
FORMEN='p2-form-welle-praefix p2-form-welle-zeichen p2-form-welle-slash p2-form-welle-punkt p2-form-welle-nur-praefix p8-sprache'
# Fälle, die den vollen Inhalt eines geschriebenen Stubs vergleichen (Stelle und Form der Zeile)
SCHREIBT='p1-alle p1-eine p1-make-welle p1-ohne-stub p1-welle-flach p12-adr p12-arc p12-spec p2-welle-leer p3-temp-und-mv p4-anker-klein p4-ein-token p4-ende-kopf p4-kein-bezug p4-kein-token p4-kopfende-genau p4-kurzform p4-leer-wert p4-mehrzeilig p4-suffix-buchstabe p4-suffix-punkt p4-vorn-buchstabe p4-vorn-strich p4-vorn-ziffer p4-zwei-bezug p5-duplikat p5-ordnung p5-gleiche-zahl p5-zahl p6-zweiter-lauf p7-ende p7-mitte p7-ohne-umbruch p7-ohne-umbruch-weitere p8-volltext-crlf p9-wurzel'

# Punkt 1
mut m-p1-welle-ignoriert "p1-eine p1-make-welle" skript $'if [ -n "$welle" ]; then\n  alle=' $'if false; then\n  alle='
mut m-p1-welle-ohne-stub-zaehlt "p1-ohne-stub p1-leere-menge p1-eine-ohne-stub p2-welle-ohne-stub-ohne-archiv" skript '  compgen -G "$done_dir/$w/slice-*.md" > /dev/null && wellen+=("$w")' '  wellen+=("$w")'
mut m-p1-welle-datei-zaehlt "p1-welle-datei" skript '"$done_dir/$w/slice-*.md"' '"$done_dir/$w/"*.md' '"$d"/slice-*.md' '"$d"/*.md'
mut m-p1-done-selbst-zaehlt "p1-welle-flach" skript 'for d in "$done_dir"/*/; do' 'for d in "$done_dir"/*/ "$done_dir/."/; do'
mut m-p1-nur-erste-welle "p1-alle p2-keine-teilschreibung p2-reihenfolge-4-5 p3-reihenfolge-aufruf" skript 'do alle+=("$(basename "$d")"); done' 'do [ "${#alle[@]}" -eq 0 ] && alle+=("$(basename "$d")"); done'
mut m-p1-make-ohne-welle "p1-make-welle" fragment '"$(WELLE)"' '""'
# Punkt 2
mut m-p2-form-gestrichen "$FORMEN p2-reihenfolge-1-2" skript "$L_FORM" ':'
mut m-p2-form-ohne-praefix "p2-form-welle-praefix p2-form-welle-nur-praefix" skript '^welle-[a-z0-9-]+$' '^[a-z0-9-]+$'
mut m-p2-form-grossbuchstaben "p2-form-welle-zeichen p2-reihenfolge-1-2" skript '^welle-[a-z0-9-]+$' '^welle-[A-Za-z0-9-]+$'
mut m-p2-form-slash "p2-form-welle-slash" skript '^welle-[a-z0-9-]+$' '^welle-[a-z0-9/-]+$'
mut m-p2-form-punkt "p2-form-welle-punkt" skript '^welle-[a-z0-9-]+$' '^welle-[a-z0-9.-]+$'
mut m-p2-form-nichtleer "p2-form-welle-nur-praefix" skript '^welle-[a-z0-9-]+$' '^welle-[a-z0-9-]*$'
mut m-p2-welle-leer-ist-form "p2-welle-leer p10-i" skript $'if [ -n "$welle" ]; then\n  [[' $'if [ "${2+x}" = x ]; then\n  [['
mut m-p2-dir-gestrichen "p2-welle-unbekannt p2-reihenfolge-2-3" skript "$L_DIR" ':'
mut m-p2-unzip-gestrichen "p2-kein-unzip p2-reihenfolge-3-4" skript "$L_UNZIP" ':'
mut m-p2-archiv-fehlt-gestrichen "p2-archiv-fehlt p2-reihenfolge-4-5 p3-reihenfolge-aufruf" skript "$L_FEHLT" ':'
mut m-p2-archiv-beschaedigt-gestrichen "p2-archiv-beschaedigt p2-archiv-crc" skript "$L_BESCH" ':'
mut m-p2-test-nur-liste "p2-archiv-crc" skript 'unzip -tq "$done_dir/$w/archiv.zip"' 'unzip -Z1 "$done_dir/$w/archiv.zip"'
mut m-p2-reihenfolge-1-2 "$FORMEN p2-reihenfolge-1-2" skript "  $L_FORM$LF  $L_DIR" "  $L_DIR$LF  $L_FORM"
mut m-p2-reihenfolge-2-3 "p2-reihenfolge-2-3" skript "  $L_DIR" "  $L_UNZIP$LF  $L_DIR"
mut m-p2-reihenfolge-3-4 "p2-reihenfolge-3-4" skript "$L_UNZIP$LF" '' "$L_FEHLT" "$L_FEHLT$LF$L_UNZIP"
mut m-p2-reihenfolge-4-5 "p2-archiv-fehlt p2-reihenfolge-4-5 p3-reihenfolge-aufruf" skript "$L_FEHLT$LF$L_BESCH" "$L_BESCH$LF$L_FEHLT"
# je Welle geprüft statt vor allem Lesen: die Prüfung 4 und 5 steht in der Schleife über die Wellen
PRO_WELLE=("$L_FEHLT$LF" '' "$L_BESCH$LF" '' $'  unzip -Z1 "$d/archiv.zip" > "$tmp/liste"' $'  [ -f "$d/archiv.zip" ] || sperre "archiv.zip fehlt: $w"\n  unzip -tq "$d/archiv.zip" > /dev/null 2>&1 || sperre "archiv.zip beschädigt: $w/archiv.zip"\n  unzip -Z1 "$d/archiv.zip" > "$tmp/liste"')
mut m-p2-geprueft-je-welle "p2-reihenfolge-4-5 p3-reihenfolge-aufruf" skript "${PRO_WELLE[@]}"
mut m-p2-schreiben-je-welle "p2-reihenfolge-4-5 p3-reihenfolge-aufruf p2-keine-teilschreibung p3-stub-ohne-eintrag p3-alle-befunde p3-datei-nicht-beschreibbar p10-i p3-schreiben-scheitert p3-schreiben-temp-aufgeraeumt" skript "${PRO_WELLE[@]}" $'  done\ndone\nif [ "$gesehen"' $'  done\n  for stub in "${todo[@]}"; do schreibe "$stub" "${neu[$stub]}"; done; todo=()\ndone\nif [ "$gesehen"'
# Punkt 3
mut m-p3-schreiben-vor-pruefung "p3-stub-ohne-eintrag p3-alle-befunde p3-datei-nicht-beschreibbar p10-i p3-schreiben-scheitert p3-schreiben-temp-aufgeraeumt" skript 'todo+=("$stub"); neu["$stub"]="$soll"' 'schreibe "$stub" "$soll"'
mut m-p3-nur-erster-befund "p3-alle-befunde" skript 'echo "archiv-bezug: $1: $2" >&2; befunde=1' 'echo "archiv-bezug: $1: $2" >&2; exit 1'
mut m-p3-ohne-eintrag-uebersprungen "p3-stub-ohne-eintrag p3-alle-befunde p3-erster-grund p10-i p3-reihenfolge-grund-a-b p3-fast-gleicher-pfad" skript '{ befund "$pfad" "kein Eintrag im Archiv"; continue; }' '{ continue; }'
mut m-p3-kein-abbruch-bei-befund "p3-stub-ohne-eintrag p3-alle-befunde p3-erster-grund p3-datei-nicht-beschreibbar p3-verzeichnis-nicht-beschreibbar p3-kein-kopf p6-abweichend p6-zwei-zeilen p7-ohne-hervorgegangen p7-zweimal-hervorgegangen p8-stub-crlf p10-i p3-reihenfolge-grund-a-b p3-reihenfolge-grund-b-c p3-reihenfolge-grund-c-c p3-reihenfolge-grund-c-d p3-reihenfolge-grund-d-e p6-gleich-falsche-stelle p3-fast-gleicher-pfad" skript '[ "$befunde" -eq 0 ] || exit 1' ':'
mut m-p3-mv-direkt "p3-temp-und-mv p3-schreiben-scheitert" skript 'mv "$t" "$stub"' 'cat "$t" > "$stub"; rm -f "$t"'
mut m-p3-temp-bleibt "p3-temp-und-mv p3-nur-stubs p3-schreiben-scheitert p3-schreiben-temp-aufgeraeumt" skript 'mv "$t" "$stub"' 'cp "$t" "$stub"'
mut m-p3-eintrag-als-teilzeichenkette "p3-fast-gleicher-pfad" skript 'grep -Fxq --' 'grep -Fq --'
mut m-p3-kein-kopf-geduldet "p3-kein-kopf p3-reihenfolge-grund-b-c" skript '|| { befund "$pfad" "Eintrag im Archiv ohne Kopf (keine Zeile ---)"; continue; }' '|| true'
mut m-p3-datei-nicht-beschreibbar-geduldet "p3-datei-nicht-beschreibbar" skript '[ ! -w "$stub" ] || ' ''
mut m-p3-verzeichnis-nicht-beschreibbar-geduldet "p3-verzeichnis-nicht-beschreibbar" skript ' [ ! -w "$d" ]' ' false'
# Punkt 4
mut m-p4-kopfende-nicht-erkannt "p4-ende-kopf p4-kopfende-genau" skript '{ gef = 1; exit }' '{ gef = 1 }'
mut m-p4-kopfende-ungenau "p4-kopfende-genau" skript '$0 == "---"' 'index($0, "---") == 1'
mut m-p4-fortsetzung-endet-an-leerzeile "p4-mehrzeilig" skript 'else if (lies == 1 && index($0, "**") == 1)' 'else if (lies == 1 && ($0 == "" || index($0, "**") == 1))'
mut m-p4-keine-fortsetzung "p4-mehrzeilig p4-ende-kopf p4-kopfende-genau" skript 'else if (lies == 1 && index($0, "**") == 1)' 'else if (lies == 1)'
mut m-p4-fortsetzung-ohne-ende-bei-sternen "p4-mehrzeilig p4-zwei-bezug" skript 'else if (lies == 1 && index($0, "**") == 1) lies = 2' ''
mut m-p4-suffix-punkt "p4-suffix-punkt p5-duplikat" skript '[0-9]+/)' '[0-9]+(\.[a-z])?/)'
mut m-p4-suffix-buchstabe "p4-suffix-buchstabe" skript '[0-9]+/)' '[0-9]+[a-z]?/)'
mut m-p4-kurzform-gelesen "p4-kurzform" skript '/LH-(FA|QA|RB)-[0-9]+/' '/(LH-(FA|QA|RB))?-[0-9]+/'
mut m-p4-anker-klein-gelesen "p4-anker-klein" skript 't = " " $0' 't = " " toupper($0)'
mut m-p4-vorn-grossbuchstabe-erlaubt "p4-vorn-buchstabe" skript '[A-Za-z0-9-]' '[a-z0-9-]'
mut m-p4-vorn-kleinbuchstabe-erlaubt "p4-vorn-buchstabe" skript '[A-Za-z0-9-]' '[A-Z0-9-]'
mut m-p4-vorn-ziffer-erlaubt "p4-vorn-ziffer" skript '[A-Za-z0-9-]' '[A-Za-z-]'
mut m-p4-vorn-strich-erlaubt "p4-vorn-strich" skript '[A-Za-z0-9-]' '[A-Za-z0-9]'
mut m-p4-kein-leerwert "p4-kein-bezug p4-kein-token p4-leer-wert p4-anker-klein" skript '${z:-— keine —}' '${z}'
mut m-p4-alle-bezug-zeilen "p4-zwei-bezug" skript 'else if (lies == 1 && index($0, "**") == 1) lies = 2' 'else if (lies == 1 && index($0, "**") == 1) lies = 0'
# Punkt 5
mut m-p5-ordnung-bytewise "p5-zahl p5-gleiche-zahl" skript '-k2,2n' '-k2,2'
mut m-p5-duplikat-bleibt "p5-duplikat" skript ' -u |' ' |'
mut m-p5-fa-nach-qa "p5-ordnung p4-mehrzeilig p4-ende-kopf p4-kopfende-genau p8-volltext-crlf" skript '== "FA" ? 1 : substr(t, RSTART + 3, 2) == "QA" ? 2 : 3' '== "FA" ? 2 : substr(t, RSTART + 3, 2) == "QA" ? 1 : 3'
mut m-p5-qa-nach-rb "p5-ordnung p4-mehrzeilig" skript '== "FA" ? 1 : substr(t, RSTART + 3, 2) == "QA" ? 2 : 3' '== "FA" ? 1 : substr(t, RSTART + 3, 2) == "QA" ? 3 : 2'
mut m-p5-gleiche-zahl-bytewise "p5-gleiche-zahl" skript ' -k4,4n' ''
mut m-p5-gleiche-zahl-verschmolzen "p5-gleiche-zahl" skript '-k2,2n -k4,4n -k3,3 -u' '-k2,2n -u'
# Punkt 6
mut m-p6-gleiche-zeile-neu "p6-gleich p6-zweiter-lauf p4-leer-wert p3-nicht-beschreibbar-gleich p3-nur-stubs" skript 'if [ "$anzahl" -eq 0 ]; then' 'if [ "$anzahl" -ge 0 ]; then'
mut m-p6-abweichend-geduldet "p6-abweichend p3-reihenfolge-grund-d-e" skript '[ "$ist" = "$soll" ] || { befund "$pfad" "Zeile weicht ab: ist \"$ist\", soll \"$soll\""; continue; }' ':'
mut m-p6-abweichend-ohne-ist "p6-abweichend" skript 'ist \"$ist\", ' ''
mut m-p6-abweichend-ohne-soll "p6-abweichend" skript ', soll \"$soll\""' '"'
mut m-p6-zwei-zeilen-geduldet "p6-zwei-zeilen" skript '[ "$anzahl" -gt 1 ]' '[ "$anzahl" -gt 2 ]'
# Punkt 7
mut m-p7-zeile-vor-hervorgegangen "$SCHREIBT p3-nur-stubs p3-schreiben-scheitert" skript '{ print } index($0, "**Hervorgegangen:**") == 1 { print ENVIRON["Z"] }' 'index($0, "**Hervorgegangen:**") == 1 { print ENVIRON["Z"] } { print }'
mut m-p7-zeile-am-ende "p7-mitte p7-ohne-umbruch-weitere" skript 'index($0, "**Hervorgegangen:**") == 1 { print ENVIRON["Z"] }' 'END { print ENVIRON["Z"] }'
mut m-p7-kein-umbruch-ergaenzt "p7-ohne-umbruch-weitere p7-mitte p7-rest-unveraendert p8-umlaute" skript "awk '{ print }" 'awk -v L="$(grep -c "" "$stub")" '"'"'{ printf "%s", $0; if (FNR < L) printf "\n" }' 'index($0, "**Hervorgegangen:**") == 1 { print ENVIRON["Z"] }' 'index($0, "**Hervorgegangen:**") == 1 { printf "\n%s\n", ENVIRON["Z"] }'
mut m-p7-ohne-hervorgegangen-geduldet "p7-ohne-hervorgegangen" skript '-eq 1 ] || { befund "$pfad" "Stub hat nicht' '-le 1 ] || { befund "$pfad" "Stub hat nicht'
mut m-p7-zweimal-hervorgegangen-geduldet "p7-zweimal-hervorgegangen" skript '-eq 1 ] || { befund "$pfad" "Stub hat nicht' '-ge 1 ] || { befund "$pfad" "Stub hat nicht'
mut m-p7-leerzeichen-entfernt "p7-rest-unveraendert" skript "awk '{ print }" "awk '{ sub(/ +\$/, \"\"); print }"
mut m-p7-leerzeilen-geleert "p7-rest-unveraendert" skript "awk '{ print }" "awk '{ sub(/^ +\$/, \"\"); print }"
# Punkt 8
mut m-p8-stub-cr-geduldet "p8-stub-crlf p3-alle-befunde p3-reihenfolge-grund-c-c p3-reihenfolge-grund-c-d" skript "if grep -q \$'\\r' \"\$stub\"; then" 'if false; then'
mut m-p8-volltext-cr-nicht-entfernt "p8-volltext-crlf" skript '{ sub(/\r$/, "") } ' ''
mut m-p8-nichtascii-entfernt "p8-umlaute" skript "' \"\$stub\" > \"\$t\" && mv" "' \"\$stub\" | tr -d '\\377' > \"\$t\" && mv"
mut m-p8-sprache-der-umgebung "p8-sprache" skript $'export LC_ALL=C\n' ''
# Punkt 9
mut m-p9-wurzel-ist-arbeitsverzeichnis "p9-wurzel p2-reihenfolge-0-1 p2-wurzel-fehlt p2-wurzel-kein-verzeichnis p2-wurzel-leer" skript 'wurzel="${1-}"' 'wurzel="$PWD"'
# Punkt 10
mut m-p10-plus-vor-rezept "p10-n p10-q p10-t p10-makeflags" fragment $'\t@bash tools' $'\t+@bash tools'
mut m-p10-make-im-rezept "p10-n p10-q p10-t p10-makeflags" fragment $'\t@bash tools' $'\t@: $(MAKE); bash tools'
mut m-p10-skript-nicht-unter-i "p10-i" fragment '@bash tools/harness/archiv-bezug.sh "$(CURDIR)" "$(WELLE)"' '@$(if $(findstring i,$(firstword $(MAKEFLAGS))),:,bash tools/harness/archiv-bezug.sh "$(CURDIR)" "$(WELLE)")'
# Punkt 11
for k in GATE_CHECKS GEGENPROBEN_CONTAINER GEGENPROBEN_HOST; do mut "m-p11-${k,,}" "p11-keine-kette" fragment ".PHONY: archiv-bezug$LF" ".PHONY: archiv-bezug$LF$k += archiv-bezug$LF"; done
# Entscheidungen des Architect zu SPEC-068 (leere Menge, Schreibbarkeit, erster Grund, Ordnung, Stichwörter)
mut m-p1-keine-stubs-still "p1-leere-menge p1-eine-ohne-stub p2-welle-ohne-stub-ohne-archiv" skript 'echo "archiv-bezug: keine Stubs"; ' ''
mut m-p1-keine-stubs-auf-stderr "p1-leere-menge p1-eine-ohne-stub p2-welle-ohne-stub-ohne-archiv" skript 'echo "archiv-bezug: keine Stubs";' 'echo "archiv-bezug: keine Stubs" >&2;'
mut m-p1-keine-stubs-immer "p1-alle" skript 'if [ "$gesehen" -eq 0 ]; then echo "archiv-bezug: keine Stubs"; exit 0; fi' 'echo "archiv-bezug: keine Stubs"; if [ "$gesehen" -eq 0 ]; then exit 0; fi'
mut m-p3-nicht-beschreibbar-fuer-alle "p3-nicht-beschreibbar-gleich" skript $'    if [ "$anzahl" -eq 0 ]; then\n      if [ ! -w "$stub" ] || [ ! -w "$d" ]; then befund "$pfad" "nicht beschreibbar"; continue; fi\n' $'    if [ ! -w "$stub" ] || [ ! -w "$d" ]; then befund "$pfad" "nicht beschreibbar"; continue; fi\n    if [ "$anzahl" -eq 0 ]; then\n'
mut m-p3-erster-grund-nicht-beschreibbar-und-weitere "p3-erster-grund p3-reihenfolge-grund-a-b p3-fast-gleicher-pfad" skript '{ befund "$pfad" "kein Eintrag im Archiv"; continue; }' '{ befund "$pfad" "kein Eintrag im Archiv"; }'
mut m-p5-laengere-zuerst "p5-gleiche-zahl" skript ' -k4,4n' ' -k4,4nr'
mut m-stichwort-praefix-eingabe "$FORMEN p2-welle-unbekannt p2-kein-unzip p2-archiv-fehlt p2-archiv-beschaedigt p2-archiv-crc p2-reihenfolge-1-2 p2-reihenfolge-2-3 p2-reihenfolge-3-4 p2-reihenfolge-4-5 p3-reihenfolge-aufruf p2-reihenfolge-0-1 p2-wurzel-fehlt p2-wurzel-kein-verzeichnis p2-wurzel-leer p2-reihenfolge-0-1 p2-wurzel-fehlt p2-wurzel-kein-verzeichnis p2-wurzel-leer" skript 'echo "archiv-bezug: $*"' 'echo "$*"'
mut m-stichwort-praefix-befund "p10-i p3-alle-befunde p3-datei-nicht-beschreibbar p3-erster-grund p3-kein-kopf p3-stub-ohne-eintrag p3-verzeichnis-nicht-beschreibbar p6-abweichend p6-zwei-zeilen p7-ohne-hervorgegangen p7-zweimal-hervorgegangen p8-stub-crlf p3-reihenfolge-grund-a-b p3-reihenfolge-grund-b-c p3-reihenfolge-grund-c-c p3-reihenfolge-grund-c-d p3-reihenfolge-grund-d-e p6-gleich-falsche-stelle p3-fast-gleicher-pfad" skript 'echo "archiv-bezug: $1: $2" >&2' 'echo "$1: $2" >&2'
mut m-stichwort-form "$FORMEN p2-reihenfolge-1-2" skript 'Form der Welle ungültig' 'Form ungültig'
mut m-stichwort-verzeichnis "p2-welle-unbekannt p2-reihenfolge-2-3" skript 'Verzeichnis fehlt' 'Ordner fehlt'
mut m-stichwort-werkzeug "p2-kein-unzip p2-reihenfolge-3-4" skript 'Werkzeug fehlt' 'Tool fehlt'
mut m-stichwort-archiv-fehlt "p2-archiv-fehlt p2-reihenfolge-4-5 p3-reihenfolge-aufruf" skript 'archiv.zip fehlt' 'Archiv fehlt'
mut m-stichwort-archiv-beschaedigt "p2-archiv-beschaedigt p2-archiv-crc" skript 'archiv.zip beschädigt' 'Archiv beschädigt'
mut m-stichwort-kein-eintrag "p3-stub-ohne-eintrag p3-alle-befunde p3-erster-grund p10-i p3-reihenfolge-grund-a-b p3-fast-gleicher-pfad" skript 'kein Eintrag im Archiv' 'kein Eintrag'
mut m-stichwort-ohne-kopf "p3-kein-kopf p3-reihenfolge-grund-b-c" skript 'Eintrag im Archiv ohne Kopf' 'Archiv-Eintrag ohne Kopf'
mut m-stichwort-cr "p3-alle-befunde p8-stub-crlf p3-reihenfolge-grund-c-c p3-reihenfolge-grund-c-d" skript 'Stub enthält CR' 'CR im Stub'
mut m-stichwort-hervorgegangen "p7-ohne-hervorgegangen p7-zweimal-hervorgegangen" skript 'Stub hat nicht genau eine Zeile **Hervorgegangen:**' 'Hervorgegangen falsch'
mut m-stichwort-mehr-als-eine-zeile "p6-zwei-zeilen" skript 'mehr als eine Zeile **Bezug:**' 'zwei Bezug-Zeilen'
mut m-stichwort-weicht-ab "p6-abweichend p3-reihenfolge-grund-d-e" skript 'Zeile weicht ab' 'Bezug anders'
mut m-stichwort-nicht-beschreibbar "p3-datei-nicht-beschreibbar p3-verzeichnis-nicht-beschreibbar" skript '"nicht beschreibbar"' '"gesperrt"'

# Review (Runde 1) und Entscheidungen arch-dt3
zeile_von() { grep -F -m1 -- "$1" "$repo/tools/harness/archiv-bezug.sh"; }
L_WURZEL="$(zeile_von 'Wurzel fehlt oder')"; L_A="$(zeile_von 'kein Eintrag im Archiv')"; L_U="$(zeile_von 'unzip -p "$d/archiv.zip"')"; L_B="$(zeile_von 'ohne Kopf')"
L_ANZ="$(zeile_von 'anzahl="$(grep -c')"; L_C1="$(zeile_von 'Stub enthält CR')"; L_C2="$(zeile_von 'nicht genau eine Zeile')"; L_D1="$(zeile_von 'mehr als eine Zeile')"; L_W="$(zeile_von 'nicht beschreibbar')"
mut m-p2-wurzel-gestrichen "p2-wurzel-fehlt p2-wurzel-leer p2-wurzel-kein-verzeichnis p2-reihenfolge-0-1" skript "$L_WURZEL$LF" ''
mut m-p2-reihenfolge-0-1 "p2-reihenfolge-0-1" skript "$L_WURZEL$LF" '' 'command -v unzip' "$L_WURZEL${LF}command -v unzip"
mut m-p3-schreiben-scheitert-exit-0 "p3-schreiben-scheitert" skript 'Schreiben fehlgeschlagen" >&2; exit 2;' 'Schreiben fehlgeschlagen" >&2; exit 0;'
mut m-p3-schreiben-scheitert-temp-bleibt "p3-schreiben-temp-aufgeraeumt" skript 'rm -f "$stub.$$"; ' ''
mut m-p3-schreiben-scheitert-weiter "p3-schreiben-scheitert" skript 'Schreiben fehlgeschlagen" >&2; exit 2;' 'Schreiben fehlgeschlagen" >&2;'
mut m-p3-grund-b-vor-a "p3-reihenfolge-grund-a-b p10-i p3-alle-befunde p3-erster-grund p3-stub-ohne-eintrag p3-fast-gleicher-pfad" skript "$L_A$LF$L_U$LF$L_B" "$L_U$LF$L_B$LF$L_A"
mut m-p3-grund-c-vor-b "p3-reihenfolge-grund-b-c" skript "$L_B$LF$L_C1" "$L_C1$LF$L_B"
mut m-p3-grund-hervorgegangen-vor-cr "p3-reihenfolge-grund-c-c" skript "$L_C1$LF$L_C2" "$L_C2$LF$L_C1"
mut m-p3-grund-d-vor-c "p3-reihenfolge-grund-c-d" skript "$L_D1$LF" '' "$L_C1" "$L_ANZ$LF$L_D1$LF$L_C1"
mut m-p3-grund-nicht-beschreibbar-vor-d "p3-reihenfolge-grund-d-e p3-nicht-beschreibbar-gleich" skript "$L_W$LF" '' "$L_D1" "$L_W$LF$L_D1"
mut m-p3-schreiben-extra-datei-done "p3-nur-stubs" skript 'for stub in "${todo[@]}"; do' 'for stub in "${todo[@]}"; do : > "$done_dir/extra.log";'
mut m-p3-schreiben-extra-datei-wurzel "p3-nur-stubs" skript 'for stub in "${todo[@]}"; do' 'for stub in "${todo[@]}"; do : > "$wurzel/extra.log";'
mut m-p3-schreiben-beruehrt-archiv "p3-nur-stubs" skript 'for stub in "${todo[@]}"; do' 'for stub in "${todo[@]}"; do touch "$done_dir"/*/archiv.zip;'
mut m-p3-schreiben-beruehrt-alle-stubs "p3-nur-stubs" skript 'for stub in "${todo[@]}"; do' 'for stub in "${todo[@]}"; do touch "$done_dir"/*/slice-*.md;'
mut m-p6-falsche-stelle-geduldet "p6-gleich-falsche-stelle" skript '|| { befund "$pfad" "Zeile an falscher Stelle"; continue; }' '|| true'
mut m-stichwort-falsche-stelle "p6-gleich-falsche-stelle" skript 'Zeile an falscher Stelle' 'falsch'
mut m-stichwort-schreiben-fehlgeschlagen "p3-schreiben-scheitert" skript 'Schreiben fehlgeschlagen' 'Schreibfehler'
mut m-stichwort-wurzel "p2-wurzel-fehlt p2-wurzel-leer p2-wurzel-kein-verzeichnis p2-reihenfolge-0-1" skript 'Wurzel fehlt oder' 'Basis fehlt oder'
mut m-name-wurzel "p2-wurzel-kein-verzeichnis" skript ': ${wurzel:-(leer)}' ''
mut m-name-form "$FORMEN p2-reihenfolge-1-2" skript 'Form der Welle ungültig: $welle' 'Form der Welle ungültig'
mut m-name-verzeichnis "p2-welle-unbekannt p2-reihenfolge-2-3" skript 'Verzeichnis fehlt: $rel/$welle' 'Verzeichnis fehlt'
mut m-name-werkzeug "p2-kein-unzip p2-reihenfolge-3-4" skript 'Werkzeug fehlt: unzip' 'Werkzeug fehlt'
mut m-name-archiv-beschaedigt "p2-archiv-beschaedigt p2-archiv-crc" skript 'archiv.zip beschädigt: $w/archiv.zip' 'archiv.zip beschädigt'
mut m-p4-art-fa-nicht-gelesen "p1-alle p1-eine p1-make-welle p1-ohne-stub p1-welle-flach p12-adr p12-arc p12-spec p2-welle-leer p3-nicht-beschreibbar-gleich p3-nur-stubs p3-schreiben-scheitert p3-temp-und-mv p4-ein-token p4-ende-kopf p4-kopfende-genau p4-mehrzeilig p4-suffix-punkt p4-zwei-bezug p5-duplikat p5-gleiche-zahl p5-ordnung p5-zahl p6-abweichend p6-gleich p6-gleich-falsche-stelle p6-zweiter-lauf p7-ende p7-mitte p7-ohne-umbruch p7-ohne-umbruch-weitere p8-umlaute p8-volltext-crlf p9-wurzel" skript '/LH-(FA|QA|RB)-[0-9]+/' '/LH-(QA|RB)-[0-9]+/'
mut m-p4-art-qa-nicht-gelesen "p4-ende-kopf p4-kopfende-genau p4-kurzform p4-mehrzeilig p4-suffix-buchstabe p4-vorn-buchstabe p4-vorn-strich p4-vorn-ziffer p5-ordnung p8-volltext-crlf" skript '/LH-(FA|QA|RB)-[0-9]+/' '/LH-(FA|RB)-[0-9]+/'
mut m-p4-art-rb-nicht-gelesen "p4-mehrzeilig p5-ordnung" skript '/LH-(FA|QA|RB)-[0-9]+/' '/LH-(FA|QA)-[0-9]+/'
# Punkt 12
mut m-p12-adr "p12-adr" skript '/LH-(FA|QA|RB)-[0-9]+/' '/(LH-(FA|QA|RB)|ADR)-[0-9]+/'
mut m-p12-spec "p12-spec" skript '/LH-(FA|QA|RB)-[0-9]+/' '/(LH-(FA|QA|RB)|SPEC)-[0-9]+/'
mut m-p12-arc "p12-arc" skript '/LH-(FA|QA|RB)-[0-9]+/' '/(LH-(FA|QA|RB)|ARC)-[0-9]+/'

wait
fehler=0
if [ -z "${GP_MUT-}" ]; then
  for n in "${FAELLE[@]}"; do
    grep -qxF -- "$n" "$arbeit/rote" || { for e in ${ENTFALLEN[@]+"${ENTFALLEN[@]}"}; do [ "$e" = "$n" ] && continue 2; done; melde "Fall '$n' steht in keiner Liste roter Fälle"; }
  done
fi
if [ -s "$arbeit/fehler" ]; then cat "$arbeit/fehler" >&2; fehler=1; fi
if [ "$(wc -l < "$arbeit/gefahren")" -eq 0 ]; then echo "archiv-bezug-gegenprobe: keine Mutation und keine Kopie gefahren (GP_MUT passt auf nichts)" >&2; fehler=1; fi
if [ "$fehler" -eq 0 ]; then echo "archiv-bezug-gegenprobe: grün ($(wc -l < "$arbeit/gefahren") Läufe)"; fi
exit "$fehler"
