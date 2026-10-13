#!/usr/bin/env bash
# stop-hook-gegenprobe — prüft den Stop-Hook .claude/hooks/stop-require-gates.sh und die
# Dateiliste (tools/harness/dateiliste.sh, working-tree-hash.sh, record-gates.sh) gegen die
# Punkte von SPEC-052 (spec/spezifikation.md), die sie in diesem Stand zusagt. Der Name jedes
# Falls trägt die Nummer des Punktes (pNN-…). Der Arbeitsbaum bleibt unberührt: jeder Fall
# läuft in einer Kopie eines Temp-Repos, das Makefile, Fragment, Skripte und Hook aus dem
# Arbeitsbaum trägt. Die unveränderte Kopie muss jeden Fall bestehen; je Mutant (Änderung
# an der Kopie von Hook oder Skript) müssen die benannten Fälle rot werden, die als grün
# genannten Fälle grün bleiben.
#
# Fälle: Punkt 1 (Liste nur nach grünem `make gates`, Hash daneben), 10 (Hash wechselt genau mit dem Arbeitsbaum), 11 (Form, Pfade mit Zeilenumbruch,
# Leerzeichen innen, am Anfang und am Ende, Tabulator, Backslash und Nicht-UTF-8), 14 (Inhalt, Link-Ziel, Umbenennung, gelöschte und neue Datei), 16 (Identität gibt
# frei), 17 (Merge-Konflikt vor der Identität), 18 (untracked Datei), 19 (kein Stand, auch bei
# sauberem Baum), 20 (nicht lesbarer Vergleich, Abbruch, beschädigte Liste, Schleifenschutz vor jeder
# Prüfung, die zu `block` führt),
# 21 (Make-Flags beim Schreiben des Nachweises: -i entfernt ihn, -n, -q, -t lassen ihn stehen),
# 22 (Aufrufarten: Arbeitsverzeichnis im Unterverzeichnis, leere Umgebung).
#
# GRENZE der Gegenprobe:
# - Sie fährt den Hook mit seiner Eingabe auf stdin und liest seine Ausgabe; ob die Umgebung
#   von Claude Code die Antwort so liest, prüft sie nicht. Von den Aufrufarten prüft sie das
#   Arbeitsverzeichnis im Unterverzeichnis und die leere Umgebung (Punkt 22); `bash -e` und
#   eine fehlende Variable CLAUDE_PROJECT_DIR sind nicht Gegenstand. „kein Repo“ und „ohne
#   git“ sind Fälle des Scheiterns von git (Punkt 20).
# - Die Make-Flags -n, -q und -t prüft sie am Ergebnis (Nachweis unverändert), nicht am
#   Mechanismus: make führt das Rezept des Nachweises dort nicht aus (Punkt 21).
# - Die Gates der Kopie sind Attrappen; die echte Gate-Kette fährt sie nicht.
# - Die Fehler von git und sha256sum kommen von Wrappern vor dem echten Programm im PATH.
# - Als root bleibt eine Datei mit Rechten 000 lesbar; dann sind p20-stand-unlesbar und sein
#   Mutant m-unlesbar-als-unveraendert nicht herstellbar und werden beide mit einer Meldung
#   auf stderr übersprungen, ohne dass die Gegenprobe deshalb rot endet.
# - Sortierung und Einmaligkeit der Liste haben keinen eigenen Fall (SPEC-052 Punkt 11,
#   akzeptiertes Negativ); gleichzeitige Läufe sind nicht zugesagt (Punkt 8, ohne Fall).
#
# Ausgang: 0, wenn jeder Fall wie erwartet endet, sonst 1.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
wurzel="$(cd "$hier/../.." && pwd)"
arbeit="$(mktemp -d)"
trap 'rm -rf "${arbeit:?}"' EXIT
fehler=0

melde() { echo "stop-hook-gegenprobe: ROT — $*" >&2; fehler=1; }

export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null GIT_CEILING_DIRECTORIES="$arbeit"
export GIT_AUTHOR_NAME=gegenprobe GIT_AUTHOR_EMAIL=gegenprobe@example.invalid
export GIT_COMMITTER_NAME=gegenprobe GIT_COMMITTER_EMAIL=gegenprobe@example.invalid
REAL_GIT="$(command -v git)"
REAL_SHA="$(command -v sha256sum)"
export REAL_GIT REAL_SHA

# Wrapper vor dem echten Programm: git beendet `ls-files -u` (SHIM=u) oder die Liste (SHIM=liste),
# sha256sum (SHIM=sha) nach der echten Ausgabe mit Status 1.
mkdir -p "$arbeit/shim" "$arbeit/ohne-git" "$arbeit/leer"
cat > "$arbeit/shim/git" <<'EOT'
#!/usr/bin/env bash
case "${SHIM:-}|$*" in
  "u|ls-files -u" | "liste|ls-files -z "*) "$REAL_GIT" "$@"; exit 1 ;;
esac
exec "$REAL_GIT" "$@"
EOT
cat > "$arbeit/shim/sha256sum" <<'EOT'
#!/usr/bin/env bash
if [ "${SHIM:-}" = sha ]; then "$REAL_SHA" "$@"; exit 1; fi
exec "$REAL_SHA" "$@"
EOT
chmod +x "$arbeit/shim/git" "$arbeit/shim/sha256sum"
# PATH ohne git: die übrigen Programme, die Hook und Liste rufen.
for t in cat mktemp cmp rm sort readlink sha256sum mv; do
  ln -s "$(command -v "$t")" "$arbeit/ohne-git/$t"
done

HOOK=.claude/hooks/stop-require-gates.sh
LIB=tools/harness/dateiliste.sh
HASH=tools/harness/working-tree-hash.sh
REC=tools/harness/record-gates.sh
H64=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa

# --- Mutationswerkzeug --------------------------------------------------------------
# ersetze <datei> <zeile> <neu> — ersetzt die Zeile, die genau <zeile> ist, durch <neu>.
ersetze() {
  local d="$w/$1" vorher
  vorher="$(cat "$d")"
  M="$2" E="$3" awk '$0 == ENVIRON["M"] { print ENVIRON["E"]; next } { print }' "$d" > "$d.neu"
  mv "$d.neu" "$d"
  [ "$(cat "$d")" != "$vorher" ] || melde "Mutant '$mname' ändert nichts ($1)"
}

# sedmut <datei> <sed-skript> — sed auf der Kopie.
sedmut() {
  local d="$w/$1" vorher
  vorher="$(cat "$d")"
  sed -i "$2" "$d"
  [ "$(cat "$d")" != "$vorher" ] || melde "Mutant '$mname' ändert nichts ($1)"
}

# Zeilen des Hooks, die Mutanten ersetzen.
L_CMP='if cmp -s "$stand" "$aktuell"; then'
L_KONFLIKT='konflikt="$(git ls-files -u)"'
L_KONFLIKT2='[ -z "$konflikt" ] || voll "There are unresolved merge entries in the index."'
L_STANDFEHLT='[ -f "$stand" ] || voll "No recorded successful make gates run (file list missing)."'
L_ERZEUGEN='dateiliste_erzeugen > "$aktuell"'
L_TRAP='trap sicherung EXIT'
L_LOOP='if [[ "$input" =~ \"stop_hook_active\"[[:space:]]*:[[:space:]]*true ]]; then'
L_SCHLUSS='voll "The working tree differs from the last recorded gates run (commits without a gates run count too)."'

# Vergleichszeilen: die Felder beider Listen je Zeile (tr), zu Einträgen oder Teilen davon (awk).
IFS= read -r -d '' TPL <<'EOT' || true
if [ "$(tr '\0' '\n' < "$stand" | awk '@P@' @S@)" = "$(tr '\0' '\n' < "$aktuell" | awk '@P@' @S@)" ]; then
EOT
TPL="${TPL%$'\n'}"
vergleich() { local t="${TPL//@P@/$1}"; printf '%s' "${t//@S@/${2:-}}"; }
EINTRAG='NR%3==1{a=$0;next} NR%3==2{v=$0;next} {print a "|" v "|" $0}'
IFS= read -r -d '' TPL_MENGE <<'EOT' || true
if [ -z "$(comm @OPT@ <(tr '\0' '\n' < "$stand" | awk '@E@' | sort) <(tr '\0' '\n' < "$aktuell" | awk '@E@' @F@ | sort))" ]; then
EOT
TPL_MENGE="${TPL_MENGE%$'\n'}"
menge() { local t="${TPL_MENGE//@OPT@/$1}"; t="${t//@E@/$EINTRAG}"; printf '%s' "${t//@F@/${2:-}}"; }

# --- Temp-Repo und Läufer ------------------------------------------------------------
# fixture <name> — legt das Temp-Repo an (Makefile, Fragment, Skripte, Hook) und setzt w.
fixture() {
  mname="$1"
  w="$arbeit/$1"
  mkdir -p "$w/.claude/hooks" "$w/tools/harness" "$w/harness/mk" "$w/.harness" "$w/dir"
  cp "$wurzel/$HOOK" "$w/$HOOK"
  cp "$wurzel/$LIB" "$wurzel/$HASH" "$wurzel/$REC" "$w/tools/harness/"
  cp "$wurzel/Makefile" "$w/Makefile"
  cp "$wurzel/harness/mk/enforce.mk" "$w/harness/mk/"
  printf 'state/\n' > "$w/.harness/.gitignore"
  {
    printf 'FAKE := fake-a fake-b\n.PHONY: $(FAKE)\n$(FAKE):\n'
    printf '\t@if [ "$(FAIL)" = "$@" ]; then exit 1; fi\n'
    printf 'GATE_CHECKS += $(FAKE)\n'
  } > "$w/harness/mk/fake.mk"
  printf 'eins\n' > "$w/a.txt"
  printf 'zwei\n' > "$w/dir/b.txt"
  ln -s a.txt "$w/link"
  git -C "$w" init -q -b main
  git -C "$w" add -A
  git -C "$w" commit -q -m basis
}

# kopie <fall> — kopiert das Temp-Repo samt Mutation, bucht diese ein und setzt c.
kopie() {
  c="$w.$1"
  rm -rf "${c:?}"
  cp -a "$w" "$c"
  git -C "$c" add -A
  git -C "$c" commit -q --allow-empty -m kopie
  cwd="$c"; pfad="$PATH"; shim=""; leer=0
}

# stand — schreibt den Nachweis des vollen Laufs in c, wie `make gates` es am Ende tut.
stand() { (cd "$c" && bash "$REC") 2> /dev/null; }

# hook [<json>] — führt den Hook aus; setzt ent (decision oder ein Fehlerwort) und hcode.
hook() {
  local json='{}' zeilen
  [ $# -eq 0 ] || json="$1"
  set +e
  if [ "$leer" -eq 1 ]; then
    hout="$(cd "$cwd" && env -i "$BASH" "$c/$HOOK" <<< "$json" 2> /dev/null)"
  else
    hout="$(cd "$cwd" && env PATH="$pfad" SHIM="$shim" "$BASH" "$c/$HOOK" <<< "$json" 2> /dev/null)"
  fi
  hcode=$?
  set -e
  zeilen="$(printf '%s\n' "$hout" | wc -l)"
  ent="$(printf '%s' "$hout" | sed -n 's/^{"decision":"\([a-z]*\)".*/\1/p')"
  if [ "$hcode" -ne 0 ]; then
    ent="exit-$hcode"
  elif [ -z "$ent" ] || [ "$zeilen" -ne 1 ]; then
    ent="keine-antwort"
  fi
}

# ist <fall> <erwartet> — meldet den Fall, wenn die Antwort nicht die erwartete ist.
ist() { [ "$ent" = "$2" ] || echo "$1"; }

# stub <inhalt> — ersetzt in c den Erzeuger der Liste durch einen, der die Bytes liefert, und
# legt dieselben Bytes als Nachweis ab: zwei gleiche Listen, die beschädigt sind.
stub() {
  printf "$1" > "$c/.stub-liste"
  printf 'dateiliste_erzeugen() { cat -- .stub-liste; }\n' >> "$c/$LIB"
  mkdir -p "$c/.harness/state"
  cp "$c/.stub-liste" "$c/.harness/state/gates-passed.dateiliste"
}

# --- Fälle: jeder meldet seinen Namen, wenn er nicht wie erwartet endet ---------------
FAELLE=(
  p16-identisch
  p14-inhalt p14-link-ziel p14-umbenennung p14-geloescht p14-geloescht-index p14-neu-getrackt
  p18-untracked p11-zeilenumbruch p11-leerzeichen p11-form
  p11-fuehrendes-leerzeichen p11-abschliessendes-leerzeichen p11-tabulator p11-backslash
  p11-nicht-utf8 p11-utf8
  p17-konflikt-identisch
  p19-ohne-stand-sauber p19-ohne-stand-geaendert
  p20-stand-unlesbar
  p20-leer p20-feldzahl-4 p20-feldzahl-5 p20-feldzahl-unterminiert p20-art p20-hash-kurz
  p20-hash-nichthex p20-wert-g p20-wert-x p20-doppelt
  p20-abgeschnitten-eintragsende
  p20-git-u-scheitert p20-git-liste-scheitert p20-datei-nicht-lesbar
  p20-abbruch p20-lib-fehlt p20-kein-repo p20-ohne-git
  p20-schleifenschutz p20-schleifenschutz-false
  p20-schleifenschutz-konflikt p20-schleifenschutz-ohne-stand p20-schleifenschutz-beschaedigt
  p20-schleifenschutz-kein-repo p20-schleifenschutz-ohne-git p20-schleifenschutz-git-u
  p20-schleifenschutz-git-liste p20-schleifenschutz-lib-fehlt
  p10-hash-stabil p10-hash-wechselt
  p1-liste-bei-gruen p1-liste-nicht-bei-rot p1-hash-neben-liste
  p21-i-rot p21-i-gruen p21-i-env p21-i-ki p21-i-lang p21-i-j p21-n p21-q p21-t p21-ni
  p22-unterverzeichnis p22-leere-umgebung
)

f_p16_identisch() { kopie p16-identisch; stand; hook; ist p16-identisch approve; }

f_p14_inhalt() { kopie p14-inhalt; stand; echo x >> "$c/a.txt"; hook; ist p14-inhalt block; }
f_p14_link_ziel() { kopie p14-link-ziel; stand; ln -sfn dir/b.txt "$c/link"; hook; ist p14-link-ziel block; }
f_p14_umbenennung() { kopie p14-umbenennung; stand; git -C "$c" mv a.txt c.txt; hook; ist p14-umbenennung block; }
f_p14_geloescht() { kopie p14-geloescht; stand; rm -f "${c:?}/a.txt"; hook; ist p14-geloescht block; }
f_p14_geloescht_index() { kopie p14-geloescht-index; stand; git -C "$c" rm -q a.txt; hook; ist p14-geloescht-index block; }
f_p14_neu_getrackt() {
  kopie p14-neu-getrackt; stand; echo n > "$c/neu.txt"; git -C "$c" add neu.txt; hook; ist p14-neu-getrackt block
}
f_p18_untracked() { kopie p18-untracked; stand; echo n > "$c/neu.txt"; hook; ist p18-untracked block; }
f_p11_zeilenumbruch() {
  kopie p11-zeilenumbruch
  printf 'a\n' > "$c/x"$'\n'"y.txt"
  stand; echo b >> "$c/x"$'\n'"y.txt"; hook; ist p11-zeilenumbruch block
}
f_p11_leerzeichen() {
  kopie p11-leerzeichen
  printf 'a\n' > "$c/sp ace.txt"
  stand; echo b >> "$c/sp ace.txt"; hook; ist p11-leerzeichen block
}

# pfadform <fall> <dateiname> — ein Pfad dieser Form: die unveränderte Liste gibt frei, eine
# Änderung des Inhalts der Datei sperrt (der Pfad muss unverändert in die Liste gelangen).
pfadform() {
  local r=""
  kopie "$1"; printf 'a\n' > "$c/$2"; stand
  hook; [ "$ent" = approve ] || r=x
  echo b >> "$c/$2"; hook; [ "$ent" = block ] || r=x
  [ -z "$r" ] || echo "$1"
}
f_p11_fuehrendes_leerzeichen() { pfadform p11-fuehrendes-leerzeichen ' lead.txt'; }
f_p11_abschliessendes_leerzeichen() { pfadform p11-abschliessendes-leerzeichen 'trail.txt '; }
f_p11_tabulator() { pfadform p11-tabulator $'\ttab\tmitte.txt\t'; }
f_p11_backslash() { pfadform p11-backslash 'bs\x.txt'; }
f_p11_nicht_utf8() { pfadform p11-nicht-utf8 $'inv\xffalid.txt'; }
f_p11_utf8() { pfadform p11-utf8 'ä.txt'; }

# Form der Liste: alle vier Arten, ein Pfad mit Zeilenumbruch und einer mit Leerzeichen, in
# einem Mini-Repo; die erwarteten Bytes stehen hier, unabhängig vom Erzeuger.
f_p11_form() {
  local m erwartet ist_datei
  kopie p11-form
  m="$c/mini"
  mkdir -p "$m/dir"
  git -C "$m" init -q -b main
  printf 'eins\n' > "$m/a.txt"; printf 'zwei\n' > "$m/dir/b.txt"; printf 'g\n' > "$m/gone.txt"
  printf 'x\n' > "$m/x.txt"
  ln -s a.txt "$m/link"
  git -C "$m" add -A
  rm -f "${m:?}/gone.txt" "${m:?}/x.txt"; mkdir "$m/x.txt"; printf 'i\n' > "$m/x.txt/inner"
  printf 'n\n' > "$m/n"$'\n'"l.txt"; printf 's\n' > "$m/sp ace.txt"
  erwartet="$arbeit/erwartet"
  {
    printf 'F\0%s\0a.txt\0' "$(hash_von eins)"
    printf 'F\0%s\0dir/b.txt\0' "$(hash_von zwei)"
    printf 'G\0\0gone.txt\0'
    printf 'L\0a.txt\0link\0'
    printf 'F\0%s\0n\nl.txt\0' "$(hash_von n)"
    printf 'F\0%s\0sp ace.txt\0' "$(hash_von s)"
    printf 'X\0\0x.txt\0'
    printf 'F\0%s\0x.txt/inner\0' "$(hash_von i)"
  } > "$erwartet"
  ist_datei="$arbeit/ist"
  (cd "$m" && . "$c/$LIB" && dateiliste_erzeugen) > "$ist_datei" 2> /dev/null || true
  cmp -s "$erwartet" "$ist_datei" || echo p11-form
}
# hash_von <text> — SHA-256 des Textes mit Zeilenvorschub, 64 Hexzeichen.
hash_von() { printf '%s\n' "$1" | "$REAL_SHA" | cut -c1-64; }

# merge_konflikt — stellt in c einen nicht aufgelösten Index-Eintrag her; Status 1, wenn keiner entsteht.
merge_konflikt() {
  git -C "$c" checkout -q -b zweig
  echo zweig > "$c/a.txt"; git -C "$c" commit -qam zweig
  git -C "$c" checkout -q main
  echo haupt > "$c/a.txt"; git -C "$c" commit -qam haupt
  git -C "$c" merge -q zweig > /dev/null 2>&1 || true
  [ -n "$(git -C "$c" ls-files -u)" ]
}
f_p17_konflikt_identisch() {
  kopie p17-konflikt-identisch
  merge_konflikt || { echo p17-konflikt-identisch; return; }
  stand; hook; ist p17-konflikt-identisch block
}

f_p19_ohne_stand_sauber() {
  kopie p19-ohne-stand-sauber
  [ -z "$(git -C "$c" status --porcelain=v1)" ] || { echo p19-ohne-stand-sauber; return; }
  hook; ist p19-ohne-stand-sauber block
}
f_p19_ohne_stand_geaendert() { kopie p19-ohne-stand-geaendert; echo x >> "$c/a.txt"; hook; ist p19-ohne-stand-geaendert block; }

# unlesbar_geht — Status 0, wenn eine Datei mit Rechten 000 hier unlesbar ist (nicht als root).
unlesbar_geht() {
  local d="$arbeit/unlesbar-probe" r=0
  : > "$d"; chmod 000 "$d"
  [ ! -r "$d" ] || r=1
  rm -f -- "$d"
  return "$r"
}
f_p20_stand_unlesbar() {
  if ! unlesbar_geht; then
    echo "stop-hook-gegenprobe: p20-stand-unlesbar übersprungen (Dateien mit Rechten 000 sind lesbar)" >&2
    return
  fi
  kopie p20-stand-unlesbar; stand
  chmod 000 "$c/.harness/state/gates-passed.dateiliste"
  hook; ist p20-stand-unlesbar block
}

# Beschädigte Liste, die der Erzeuger des Arbeitsbaums gleich liefert (Stub): nur eine
# Prüfung des Stands sperrt, der Vergleich allein gäbe frei.
beschaedigt() { kopie "$1"; stub "$2"; hook; ist "$1" block; }
f_p20_leer() { beschaedigt p20-leer ''; }
f_p20_feldzahl_4() { beschaedigt p20-feldzahl-4 "F\\0$H64\\0a\\0L\\0"; }
f_p20_feldzahl_5() { beschaedigt p20-feldzahl-5 "F\\0$H64\\0a\\0L\\0x\\0"; }
f_p20_feldzahl_unterminiert() { beschaedigt p20-feldzahl-unterminiert "F\\0$H64\\0a\\0F"; }
f_p20_art() { beschaedigt p20-art 'Q\0\0a\0'; }
f_p20_hash_kurz() { beschaedigt p20-hash-kurz 'F\0abc\0a\0'; }
f_p20_hash_nichthex() { beschaedigt p20-hash-nichthex "F\\0${H64//a/z}\\0a\\0"; }
f_p20_wert_g() { beschaedigt p20-wert-g 'G\0x\0a\0'; }
f_p20_wert_x() { beschaedigt p20-wert-x 'X\0x\0a\0'; }
f_p20_doppelt() { beschaedigt p20-doppelt "F\\0$H64\\0a\\0F\\0$H64\\0a\\0"; }

# Liste des Standes nach <n> Feldern abgeschnitten (6: am Ende eines Eintrags).
abschneiden() {
  local felder
  kopie "$1"; stand
  mapfile -d '' -t felder < "$c/.harness/state/gates-passed.dateiliste"
  printf '%s\0' "${felder[@]:0:$2}" > "$c/.harness/state/gates-passed.dateiliste"
  hook; ist "$1" block
}
f_p20_abgeschnitten_eintragsende() { abschneiden p20-abgeschnitten-eintragsende 6; }

f_p20_git_u_scheitert() { kopie p20-git-u-scheitert; stand; pfad="$arbeit/shim:$PATH"; shim=u; hook; ist p20-git-u-scheitert block; }
f_p20_git_liste_scheitert() { kopie p20-git-liste-scheitert; stand; pfad="$arbeit/shim:$PATH"; shim=liste; hook; ist p20-git-liste-scheitert block; }
f_p20_datei_nicht_lesbar() { kopie p20-datei-nicht-lesbar; stand; pfad="$arbeit/shim:$PATH"; shim=sha; hook; ist p20-datei-nicht-lesbar block; }
f_p20_abbruch() { kopie p20-abbruch; stand; printf 'exit 3\n' >> "$c/$LIB"; hook; ist p20-abbruch block; }
f_p20_lib_fehlt() { kopie p20-lib-fehlt; stand; rm -f "${c:?}/${LIB:?}"; hook; ist p20-lib-fehlt block; }
f_p20_kein_repo() { kopie p20-kein-repo; stand; cwd="$arbeit/leer"; hook; ist p20-kein-repo block; }
f_p20_ohne_git() { kopie p20-ohne-git; stand; pfad="$arbeit/ohne-git"; hook; ist p20-ohne-git block; }

f_p20_schleifenschutz() {
  kopie p20-schleifenschutz; stand; echo x >> "$c/a.txt"
  hook '{"stop_hook_active": true}'; ist p20-schleifenschutz approve
}
f_p20_schleifenschutz_false() {
  kopie p20-schleifenschutz-false; stand; echo x >> "$c/a.txt"
  hook '{"stop_hook_active": false}'; ist p20-schleifenschutz-false block
}

# schutz <fall> — mit dem Zustand in c: `stop_hook_active: true` gibt frei, `false` sperrt.
schutz() {
  local r=""
  hook '{"stop_hook_active": true}'; [ "$ent" = approve ] || r=x
  hook '{"stop_hook_active": false}'; [ "$ent" = block ] || r=x
  [ -z "$r" ] || echo "$1"
}
f_p20_schleifenschutz_konflikt() {
  kopie p20-schleifenschutz-konflikt
  merge_konflikt || { echo p20-schleifenschutz-konflikt; return; }
  stand; schutz p20-schleifenschutz-konflikt
}
f_p20_schleifenschutz_ohne_stand() { kopie p20-schleifenschutz-ohne-stand; schutz p20-schleifenschutz-ohne-stand; }
f_p20_schleifenschutz_beschaedigt() {
  kopie p20-schleifenschutz-beschaedigt; stand
  printf 'Q' > "$c/.harness/state/gates-passed.dateiliste"
  schutz p20-schleifenschutz-beschaedigt
}
f_p20_schleifenschutz_kein_repo() { kopie p20-schleifenschutz-kein-repo; stand; cwd="$arbeit/leer"; schutz p20-schleifenschutz-kein-repo; }
f_p20_schleifenschutz_ohne_git() { kopie p20-schleifenschutz-ohne-git; stand; pfad="$arbeit/ohne-git"; schutz p20-schleifenschutz-ohne-git; }
f_p20_schleifenschutz_git_u() {
  kopie p20-schleifenschutz-git-u; stand; pfad="$arbeit/shim:$PATH"; shim=u; schutz p20-schleifenschutz-git-u
}
f_p20_schleifenschutz_git_liste() {
  kopie p20-schleifenschutz-git-liste; stand; pfad="$arbeit/shim:$PATH"; shim=liste; schutz p20-schleifenschutz-git-liste
}
f_p20_schleifenschutz_lib_fehlt() {
  kopie p20-schleifenschutz-lib-fehlt; stand; rm -f "${c:?}/${LIB:?}"; schutz p20-schleifenschutz-lib-fehlt
}

f_p10_hash_stabil() {
  local h1 h2
  kopie p10-hash-stabil
  h1="$(cd "$c" && bash "$HASH")"; h2="$(cd "$c" && bash "$HASH")"
  [ -n "$h1" ] && [ "$h1" = "$h2" ] || { echo p10-hash-stabil; return; }
  git -C "$c" add -A; git -C "$c" commit -q --allow-empty -m nochmal
  [ "$h1" = "$(cd "$c" && bash "$HASH")" ] || echo p10-hash-stabil
}
f_p10_hash_wechselt() {
  local h1
  kopie p10-hash-wechselt
  h1="$(cd "$c" && bash "$HASH")"
  echo x >> "$c/a.txt"
  [ "$h1" != "$(cd "$c" && bash "$HASH")" ] || echo p10-hash-wechselt
}

# make gates mit Attrappen: grün schreibt Liste und Hash, rot lässt keine Liste zurück.
make_gates() { (cd "$c" && env -u MAKELEVEL -u MFLAGS -u MAKEFLAGS FAIL="${1:-}" make -s --no-print-directory gates > /dev/null 2>&1); }
# dateiliste_gleich <repo> — der Nachweis ist die Liste des Arbeitsbaums.
dateiliste_gleich() {
  (cd "$1" && . "$1/$LIB" && dateiliste_erzeugen) | cmp -s - "$1/.harness/state/gates-passed.dateiliste"
}
f_p1_liste_bei_gruen() {
  kopie p1-liste-bei-gruen
  make_gates || { echo p1-liste-bei-gruen; return; }
  [ -s "$c/.harness/state/gates-passed.dateiliste" ] || { echo p1-liste-bei-gruen; return; }
  dateiliste_gleich "$c" || echo p1-liste-bei-gruen
}
f_p1_liste_nicht_bei_rot() {
  kopie p1-liste-nicht-bei-rot
  if make_gates fake-b; then echo p1-liste-nicht-bei-rot; return; fi
  [ ! -e "$c/.harness/state/gates-passed.dateiliste" ] || echo p1-liste-nicht-bei-rot
}
f_p1_hash_neben_liste() {
  kopie p1-hash-neben-liste
  make_gates || { echo p1-hash-neben-liste; return; }
  [ "$(cat "$c/.harness/state/gates-passed.diffsha" 2> /dev/null)" = "$(cd "$c" && bash "$HASH")" ] || echo p1-hash-neben-liste
}

# Punkt 21: ein grüner Lauf als älterer Nachweis, dann der Baum geändert (ein Schreiben
# würde den Nachweis ändern), dann `make gates` mit Flags.
# lauf21 <fail> <MAKEFLAGS aus der Umgebung oder leer> <make-Argumente…> — Ausgang unerheblich.
lauf21() {
  local fail="$1" umgebung="$2" mf=()
  shift 2
  [ -z "$umgebung" ] || mf=("MAKEFLAGS=$umgebung")
  (cd "$c" && env -u MAKELEVEL -u MFLAGS -u MAKEFLAGS "${mf[@]}" FAIL="$fail" \
    make --no-print-directory "$@" gates > /dev/null 2>&1) || :
}
# vorlauf21 — grüner Lauf, Kopie des Nachweises, dann der Baum geändert; Status 1, wenn der Lauf scheitert.
vorlauf21() {
  make_gates || return 1
  [ -s "$c/.harness/state/gates-passed.dateiliste" ] && [ -s "$c/.harness/state/gates-passed.diffsha" ] || return 1
  cp "$c/.harness/state/gates-passed.dateiliste" "$arbeit/vor.dateiliste"
  cp "$c/.harness/state/gates-passed.diffsha" "$arbeit/vor.diffsha"
  echo x >> "$c/a.txt"
}
nachweis_weg() {
  [ ! -e "$c/.harness/state/gates-passed.dateiliste" ] && [ ! -e "$c/.harness/state/gates-passed.diffsha" ]
}
nachweis_gleich() {
  cmp -s "$arbeit/vor.dateiliste" "$c/.harness/state/gates-passed.dateiliste" \
    && cmp -s "$arbeit/vor.diffsha" "$c/.harness/state/gates-passed.diffsha"
}
# i21 <fall> <fail> <umgebung> <make-Argumente…> — -i in irgendeiner Schreibweise: kein Nachweis.
i21() {
  local fall="$1" fail="$2" umgebung="$3"
  shift 3
  kopie "$fall"
  vorlauf21 || { echo "$fall"; return; }
  lauf21 "$fail" "$umgebung" "$@"
  nachweis_weg || echo "$fall"
}
# n21 <fall> <MAKEFLAGS aus der Umgebung oder leer> <make-Argumente…> — Nachweis byte-gleich.
n21() {
  local fall="$1" umgebung="$2"
  shift 2
  kopie "$fall"
  vorlauf21 || { echo "$fall"; return; }
  lauf21 "" "$umgebung" "$@"
  nachweis_gleich || echo "$fall"
}
f_p21_i_rot() { i21 p21-i-rot fake-b "" -i; }
f_p21_i_gruen() { i21 p21-i-gruen "" "" -i; }
f_p21_i_env() { i21 p21-i-env fake-b i; }
f_p21_i_ki() { i21 p21-i-ki fake-b "" -ki; }
f_p21_i_lang() { i21 p21-i-lang fake-b "" --ignore-errors; }
f_p21_i_j() { i21 p21-i-j fake-b "" -j4 -i; }
f_p21_n() { n21 p21-n "" -n; }
f_p21_q() { n21 p21-q "" -q; }
f_p21_t() { n21 p21-t "" -t; }
f_p21_ni() {
  local r
  r="$(n21 p21-ni "" -ni)$(n21 p21-ni ni)"
  [ -z "$r" ] || echo p21-ni
}

# Punkt 22: Aufrufarten. Identische Liste gibt frei, jeder Unterschied sperrt.
f_p22_unterverzeichnis() {
  local r=""
  kopie p22-unterverzeichnis; stand; cwd="$c/dir"
  hook
  [ "$ent" = approve ] || r=x
  echo x >> "$c/a.txt"; hook
  [ "$ent" = block ] || r=x
  [ -z "$r" ] || echo p22-unterverzeichnis
}
f_p22_leere_umgebung() {
  local r=""
  kopie p22-leere-umgebung; stand; leer=1
  hook
  [ "$ent" = approve ] || r=x
  echo x >> "$c/a.txt"; hook
  [ "$ent" = block ] || r=x
  [ -z "$r" ] || echo p22-leere-umgebung
}

# fall <name> — führt den Fall aus; Ausgabe: sein Name, wenn er rot ist.
fall() { "f_${1//-/_}"; }

# gruen — die unveränderte Kopie besteht jeden Fall.
gruen() {
  local f ist_rot=""
  for f in "${FAELLE[@]}"; do ist_rot+="$(fall "$f") "; done
  [ -z "${ist_rot// /}" ] || melde "Kopie '$mname' ist nicht grün: $ist_rot"
}

# rot <rote Fälle> [<grüne Fälle>] — der Mutant färbt jeden roten Fall rot und lässt die
# grünen Fälle grün.
rot() {
  local f ist
  for f in $1; do
    ist="$(fall "$f")"
    [ -n "$ist" ] || melde "Mutant '$mname' färbt '$f' nicht rot"
  done
  for f in ${2:-}; do
    ist="$(fall "$f")"
    [ -z "$ist" ] || melde "Mutant '$mname' färbt auch '$f' rot, nicht nur: $1"
  done
}

# --- Unveränderte Kopie --------------------------------------------------------------
fixture ok; gruen

# --- Punkt 16: eine identische Liste gibt frei ---------------------------------------
fixture m-immer-block
ersetze $HOOK "$L_CMP" 'if false; then'; rot p16-identisch

# --- Punkt 14: Inhalt, Link-Ziel, Umbenennung, gelöschte und neue Datei --------------
fixture m-nur-pfade
ersetze $HOOK "$L_CMP" "$(vergleich 'NR%3==0')"; rot "p14-inhalt p14-link-ziel"
fixture m-inhalt-ignoriert
ersetze $HOOK "$L_CMP" "$(vergleich 'NR%3==1{a=$0;next} NR%3==2{v=(a=="F")?"-":$0;next} {print a "|" v "|" $0}')"
rot p14-inhalt "p14-link-ziel p14-umbenennung"
fixture m-link-ziel-ignoriert
ersetze $HOOK "$L_CMP" "$(vergleich 'NR%3==1{a=$0;next} NR%3==2{v=(a=="L")?"-":$0;next} {print a "|" v "|" $0}')"
rot p14-link-ziel "p14-inhalt p14-umbenennung"
fixture m-nur-werte
ersetze $HOOK "$L_CMP" "$(vergleich 'NR%3==2' '| sort')"; rot p14-umbenennung "p14-inhalt p14-link-ziel"
fixture m-geloeschte-ignoriert
ersetze $HOOK "$L_CMP" "$(menge -13 "| grep -v '^G|'")"
rot "p14-geloescht p14-geloescht-index" "p14-inhalt p14-umbenennung"
fixture m-neue-ignoriert
ersetze $HOOK "$L_CMP" "$(menge -23)"
rot "p14-neu-getrackt p18-untracked p20-abgeschnitten-eintragsende" "p14-inhalt p14-geloescht"

# --- Punkt 18: untracked Datei -------------------------------------------------------
fixture m-ohne-others
sedmut $LIB 's/ --others//'; rot p18-untracked p14-neu-getrackt

# --- Punkt 11: Pfade und Form --------------------------------------------------------
fixture m-zeilenweise
sedmut $LIB "s/ls-files -z/ls-files/; s/sort -zu/sort -u/; s/while IFS= read -r -d '' pfad/while IFS= read -r pfad/"
rot "p11-zeilenumbruch p11-form" p11-leerzeichen
fixture m-ungequotet
sedmut $LIB '/^dateiliste_erzeugen/,/^)$/ s/"\$pfad"/$pfad/g'
rot "p11-leerzeichen p11-zeilenumbruch p11-form"
fixture m-felder-vertauscht
sedmut $LIB 's/"\${hash%% \*}" "\$pfad"/"$pfad" "${hash%% *}"/'; rot p11-form
fixture m-art-x-als-g
sedmut $LIB "s/'X\\\\0\\\\0/'G\\\\0\\\\0/"; rot p11-form

fixture m-ohne-ifs
sedmut $LIB "s/while IFS= read -r -d '' pfad/while read -r -d '' pfad/"
rot "p11-fuehrendes-leerzeichen p11-abschliessendes-leerzeichen p11-tabulator" "p11-backslash p11-nicht-utf8 p11-utf8 p11-leerzeichen p11-zeilenumbruch p11-form"
fixture m-read-ohne-r
sedmut $LIB "s/while IFS= read -r -d '' pfad/while IFS= read -d '' pfad/"
rot p11-backslash "p11-fuehrendes-leerzeichen p11-abschliessendes-leerzeichen p11-tabulator p11-nicht-utf8 p11-utf8 p11-leerzeichen p11-zeilenumbruch p11-form"
fixture m-nicht-utf8-normalisiert
sedmut $LIB 's/| LC_ALL=C sort -zu/| iconv -c -f UTF-8 -t UTF-8 | LC_ALL=C sort -zu/'
rot p11-nicht-utf8 "p11-utf8 p11-fuehrendes-leerzeichen p11-abschliessendes-leerzeichen p11-tabulator p11-backslash p11-leerzeichen p11-zeilenumbruch p11-form"

# --- Punkt 17: Merge-Konflikt vor der Identität ---------------------------------------
fixture m-konflikt-nicht-gelesen
ersetze $HOOK "$L_KONFLIKT" 'konflikt=""'; rot p17-konflikt-identisch
fixture m-identitaet-zuerst
ersetze $HOOK "$L_KONFLIKT2" ''
ersetze $HOOK "$L_SCHLUSS" "$L_KONFLIKT2"$'\n'"$L_SCHLUSS"; rot p17-konflikt-identisch

# --- Punkt 19: fehlender Stand, auch bei sauberem Baum ----------------------------------
fixture m-stand-fehlt-als-unveraendert
ersetze $HOOK "$L_STANDFEHLT" '[ -f "$stand" ] || frei'; rot "p19-ohne-stand-sauber p19-ohne-stand-geaendert"
fixture m-sauberer-baum-frei
ersetze $HOOK "$L_STANDFEHLT" '[ -f "$stand" ] || { [ -z "$(git status --porcelain=v1)" ] && frei; voll "x"; }'
rot p19-ohne-stand-sauber p19-ohne-stand-geaendert
fixture m-unlesbar-als-unveraendert
ersetze $HOOK "$L_STANDFEHLT" "$L_STANDFEHLT"$'\n''[ -r "$stand" ] || frei'
if unlesbar_geht; then
  rot p20-stand-unlesbar
else
  echo "stop-hook-gegenprobe: Mutant 'm-unlesbar-als-unveraendert' übersprungen (Dateien mit Rechten 000 sind lesbar)" >&2
fi

# --- Punkt 20: beschädigte Liste, je Fehlform ------------------------------------------
fixture m-leer-gilt
sedmut $LIB 's/^  \[ -s "\$datei" \] || return 1$/  :/'; rot p20-leer
fixture m-feldzahl-4-gilt
sedmut $LIB "s/^    IFS= read -r -d '' \\(wert\\|pfad\\) || return 1$/    IFS= read -r -d '' \\1 || :/"
rot "p20-feldzahl-4 p20-feldzahl-5"
fixture m-feldzahl-5-gilt
sedmut $LIB "s/^    IFS= read -r -d '' pfad || return 1$/    IFS= read -r -d '' pfad || :/"; rot p20-feldzahl-5 p20-feldzahl-4
fixture m-unterminiert-gilt
sedmut $LIB "s/ || \[ -n \"\\\$art\" \]; do/; do/"; rot p20-feldzahl-unterminiert
fixture m-art-gilt
sedmut $LIB 's/^      \*) return 1 ;;$/      *) ;;/'; rot p20-art
fixture m-hash-kurz-gilt
sedmut $LIB 's/{64}\$/+$/'; rot p20-hash-kurz p20-hash-nichthex
fixture m-hash-nichthex-gilt
sedmut $LIB 's/\^\[0-9a-fA-F\]/^./'; rot p20-hash-nichthex p20-hash-kurz
fixture m-wert-g-gilt
sedmut $LIB 's/^      G | X) \[ -z "\$wert" \] || return 1 ;;$/      G) ;; X) [ -z "$wert" ] || return 1 ;;/'
rot p20-wert-g p20-wert-x
fixture m-wert-x-gilt
sedmut $LIB 's/^      G | X) \[ -z "\$wert" \] || return 1 ;;$/      G) [ -z "$wert" ] || return 1 ;; X) ;;/'
rot p20-wert-x p20-wert-g
fixture m-doppelt-gilt
sedmut $LIB 's/^    \[ -z "\${gesehen\["p:\$pfad"\]+x}" \] || return 1$/    :/'; rot p20-doppelt

# --- Punkt 20: nicht lesbarer Vergleich, Abbruch ---------------------------------------
fixture m-git-u-ignoriert
ersetze $HOOK "$L_KONFLIKT" 'konflikt="$(git ls-files -u || true)"'; rot p20-git-u-scheitert
fixture m-git-liste-ignoriert
ersetze $HOOK "$L_ERZEUGEN" 'dateiliste_erzeugen > "$aktuell" || :'; rot "p20-git-liste-scheitert p20-datei-nicht-lesbar"
fixture m-pipefail-aus
sedmut $LIB 's/set -euo pipefail/set -eu; set +o pipefail/'; rot p20-git-liste-scheitert p20-datei-nicht-lesbar
fixture m-sha-ignoriert
sedmut $LIB 's/hash="\$(sha256sum < "\$pfad")"/hash="$(sha256sum < "$pfad" || true)"/'; rot p20-datei-nicht-lesbar p20-git-liste-scheitert
fixture m-ohne-trap
ersetze $HOOK "$L_TRAP" ''; rot "p20-abbruch p20-lib-fehlt p20-kein-repo p20-ohne-git"

# --- Punkt 20: Schleifenschutz -------------------------------------------------------
SCHUTZ_ALLE="p20-schleifenschutz-kein-repo p20-schleifenschutz-ohne-git p20-schleifenschutz-lib-fehlt p20-schleifenschutz-git-u
  p20-schleifenschutz-konflikt p20-schleifenschutz-ohne-stand p20-schleifenschutz-beschaedigt p20-schleifenschutz-git-liste"
fixture m-schleifenschutz-aus
ersetze $HOOK "$L_LOOP" 'if false; then'; rot "p20-schleifenschutz $SCHUTZ_ALLE" p20-schleifenschutz-false
fixture m-schleifenschutz-wert-ignoriert
ersetze $HOOK "$L_LOOP" 'if [[ "$input" =~ \"stop_hook_active\" ]]; then'
rot "p20-schleifenschutz-false $SCHUTZ_ALLE" p20-schleifenschutz

# schutz_nach <Zeile des Hooks> — der Schleifenschutz steht hinter dieser Zeile statt vor allen Prüfungen.
schutz_nach() {
  ersetze $HOOK "$L_LOOP" 'if false; then'
  ersetze $HOOK "$1" "$1"$'\n'"$L_LOOP"$'\n''  frei'$'\n''fi'
}
fixture m-schutz-nach-source
schutz_nach '. tools/harness/dateiliste.sh'
rot "p20-schleifenschutz-kein-repo p20-schleifenschutz-ohne-git p20-schleifenschutz-lib-fehlt" "p20-schleifenschutz p20-schleifenschutz-git-u p20-schleifenschutz-konflikt p20-schleifenschutz-ohne-stand p20-schleifenschutz-beschaedigt p20-schleifenschutz-git-liste"
fixture m-schutz-nach-git-u
schutz_nach "$L_KONFLIKT"
rot "p20-schleifenschutz-git-u" "p20-schleifenschutz p20-schleifenschutz-konflikt p20-schleifenschutz-ohne-stand p20-schleifenschutz-beschaedigt p20-schleifenschutz-git-liste"
fixture m-schutz-nach-konflikt
schutz_nach "$L_KONFLIKT2"
rot "p20-schleifenschutz-konflikt" "p20-schleifenschutz p20-schleifenschutz-ohne-stand p20-schleifenschutz-beschaedigt p20-schleifenschutz-git-liste"
fixture m-schutz-nach-stand-fehlt
schutz_nach "$L_STANDFEHLT"
rot "p20-schleifenschutz-ohne-stand" "p20-schleifenschutz p20-schleifenschutz-beschaedigt p20-schleifenschutz-git-liste"
fixture m-schutz-nach-pruefung
schutz_nach 'dateiliste_pruefen "$stand" || voll "The recorded file list is unreadable or damaged."'
rot "p20-schleifenschutz-beschaedigt" "p20-schleifenschutz p20-schleifenschutz-git-liste"
fixture m-schutz-nach-erzeugen
schutz_nach "$L_ERZEUGEN"
rot "p20-schleifenschutz-git-liste" "p20-schleifenschutz"

# --- Punkt 10: Hash -----------------------------------------------------------------
fixture m-hash-konstant
sedmut $HASH 's/^printf .*$/printf "%s\\n" konstant/'; rot p10-hash-wechselt p10-hash-stabil
fixture m-hash-zeit
sedmut $HASH 's/^printf .*$/printf "%s\\n" "${hash%% *}$(date +%s%N)"/'; rot p10-hash-stabil p10-hash-wechselt

# --- Punkt 1: Liste nur nach grünem Lauf, Hash daneben ------------------------------
fixture m-liste-vor-gates
sedmut Makefile 's/^record-gates: \$(GATE_CHECKS)$/record-gates:/; s/^gates: record-gates /gates: record-gates $(GATE_CHECKS) /'
rot p1-liste-nicht-bei-rot p1-liste-bei-gruen
fixture m-liste-nicht-geschrieben
sedmut $REC '/^mv -f/d'; rot p1-liste-bei-gruen p1-liste-nicht-bei-rot
fixture m-hash-fehlt
sedmut $REC '/gates-passed.diffsha$/d'; rot p1-hash-neben-liste p1-liste-bei-gruen

# --- Punkt 21: Make-Flags beim Schreiben des Nachweises -----------------------------------
FALLE_I="p21-i-rot p21-i-gruen p21-i-env p21-i-ki p21-i-lang p21-i-j"
FALLE_NQT="p21-n p21-q p21-t p21-ni"
fixture m-21-ohne-i
sedmut $REC 's/^  \*i\*)$/  \*nie\*)/'; rot "$FALLE_I" "$FALLE_NQT"
fixture m-21-literal-i
sedmut $REC '/^esac$/,$ { s/^case "\$buchstaben" in$/case "$MAKEFLAGS" in/; s/^  \*i\*)$/  *-i*)/; }'; rot "$FALLE_I" "$FALLE_NQT"
fixture m-21-wort-i
sedmut $REC '/^esac$/,$ { s/^case "\$buchstaben" in$/case "$erstes" in/; s/^  \*i\*)$/  i)/; }'
rot p21-i-ki "p21-i-rot p21-i-gruen p21-i-env p21-i-lang p21-i-j $FALLE_NQT"
fixture m-21-nur-liste-bleibt
sedmut $REC 's|^    rm -f -- .harness/state/gates-passed.diffsha "\$DATEILISTE_STAND"$|    rm -f -- .harness/state/gates-passed.diffsha|'
rot "$FALLE_I" "$FALLE_NQT"
fixture m-21-nur-hash-bleibt
sedmut $REC 's|^    rm -f -- .harness/state/gates-passed.diffsha "\$DATEILISTE_STAND"$|    rm -f -- "$DATEILISTE_STAND"|'
rot "$FALLE_I" "$FALLE_NQT"
fixture m-21-nichts-entfernt
sedmut $REC '/^    rm -f -- .harness\/state\/gates-passed.diffsha/d'
rot "$FALLE_I" "$FALLE_NQT"
fixture m-21-rezept-plus
sedmut harness/mk/enforce.mk 's/^\t@bash tools\/harness\/record-gates.sh$/\t+@bash tools\/harness\/record-gates.sh/'
rot "p21-n p21-t p21-ni" "p21-q $FALLE_I"
fixture m-21-rezept-plus-ohne-voraussetzung
sedmut harness/mk/enforce.mk 's/^\t@bash tools\/harness\/record-gates.sh$/\t+@bash tools\/harness\/record-gates.sh/'
sedmut Makefile 's/^record-gates: \$(GATE_CHECKS)$/record-gates:/; s/^gates: record-gates /gates: record-gates $(GATE_CHECKS) /'
rot "$FALLE_NQT" "$FALLE_I"

# --- Punkt 22: Aufrufarten ------------------------------------------------------------
fixture m-22-pwd
ersetze $HOOK 'wurzel="$(git rev-parse --show-toplevel)"' 'wurzel="$(pwd)"'; rot p22-unterverzeichnis p22-leere-umgebung
fixture m-22-home
ersetze $HOOK 'input="$(cat || true)"' ': "${HOME:?}"'$'\n''input="$(cat || true)"'; rot p22-leere-umgebung p22-unterverzeichnis

# --- Anschluss: das Gate hängt an GATE_CHECKS, sein Exit gilt -------------------------
gate_variable() {
  { make -s --no-print-directory -C "$1" -pq help 2> /dev/null || true; } \
    | awk -v n="$2" '$1 == n && ($2 == ":=" || $2 == "=") { for (i = 3; i <= NF; i++) print $i }'
}
anschluss() {
  local g ok=0
  for g in $(gate_variable "$1" GATE_CHECKS); do
    if [ "$g" = stop-hook-gegenprobe ]; then ok=1; fi
  done
  [ "$ok" -eq 1 ] || echo anschluss-gate-checks
}
anschluss_exit() {
  local code=0
  printf 'exit 1\n' > "$1/tools/harness/stop-hook-gegenprobe.sh"
  env -u MAKELEVEL -u MFLAGS -u MAKEFLAGS make -s --no-print-directory -C "$1" stop-hook-gegenprobe > /dev/null 2>&1 || code=$?
  [ "$code" -ne 0 ] || echo anschluss-exit
}
fixture_anschluss() {
  mname="$1"
  w="$arbeit/$1"
  mkdir -p "$w/harness/mk" "$w/tools/harness"
  cp "$wurzel/Makefile" "$w/Makefile"
  cp "$wurzel/harness/mk/enforce.mk" "$wurzel/harness/mk/stop-hook-gegenprobe.mk" "$w/harness/mk/"
}
[ -z "$(anschluss "$wurzel")" ] || melde "Arbeitsbaum: stop-hook-gegenprobe steht nicht in GATE_CHECKS"
fixture_anschluss anschluss-ok
[ -z "$(anschluss "$w")$(anschluss_exit "$w")" ] || melde "Kopie 'anschluss-ok' ist nicht grün"
fixture_anschluss m-ohne-anschluss
sedmut harness/mk/stop-hook-gegenprobe.mk '/^GATE_CHECKS += stop-hook-gegenprobe$/d'
[ -n "$(anschluss "$w")" ] || melde "Mutant 'm-ohne-anschluss' färbt 'anschluss-gate-checks' nicht rot"
fixture_anschluss m-exit-ignoriert
sedmut harness/mk/stop-hook-gegenprobe.mk 's|^\t@bash tools/harness/stop-hook-gegenprobe.sh|\t-@bash tools/harness/stop-hook-gegenprobe.sh|'
[ -n "$(anschluss_exit "$w")" ] || melde "Mutant 'm-exit-ignoriert' färbt 'anschluss-exit' nicht rot"

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "stop-hook-gegenprobe: gruen — ${#FAELLE[@]} Fälle der Punkte 1, 10, 11, 14, 16 bis 22 je Mutant rot, unveränderte Kopie grün, Anschluss an GATE_CHECKS"
