#!/usr/bin/env bash
# stop-hook-gegenprobe — prüft den Stop-Hook .claude/hooks/stop-require-gates.sh und die
# Dateiliste (tools/harness/dateiliste.sh, working-tree-hash.sh, record-gates.sh) gegen die
# Punkte von SPEC-052 (harness/vertraege/SPEC-052-gate-kette.md), die sie in diesem Stand zusagt. Der Name jedes
# Falls trägt die Nummer des Punktes (pNN-…). Der Arbeitsbaum bleibt unberührt: jeder Fall
# läuft in einer Kopie eines Temp-Repos, das Makefile, Fragment, Skripte und Hook aus dem
# Arbeitsbaum trägt. Die unveränderte Kopie muss jeden Fall bestehen; je Mutant (Änderung
# an der Kopie von Hook oder Skript) müssen die benannten Fälle rot werden, die als grün
# genannten Fälle grün bleiben.
#
# Fälle: Punkt 1 (Liste nur nach grünem `make gates`, Hash daneben), 9 (Freigabe bei inerten Pfaden und grünem Stand der
# schnellen Kette; Reihenfolge der Prüfungen; Pfade außerhalb der Positivliste; der `reason` bei fehlendem Stand der
# schnellen Kette nennt `make gates-schnell` bei inerter Änderung und nur `make gates` bei nicht inerter;
# der Wortlaut bei fremdem, leerem und nicht lesbarem Stand ist nicht zugesagt und hat keinen Fall), 10 (Hash wechselt genau mit dem Arbeitsbaum;
# Stand der schnellen Kette fehlt, ist leer, nicht lesbar oder fremd; das Hash-Skript fehlt, scheitert oder liefert nichts),
# 11 (Form, Pfade mit Zeilenumbruch,
# Leerzeichen innen, am Anfang und am Ende, Tabulator, Backslash und Nicht-UTF-8), 12 und 13 (Positivliste je Eintrag in jeder
# Tiefe, Grenze je Eintrag, Groß- und Kleinschreibung, unzulässige Pfade, Art X, Hook, Spezifikation und Vertrag nicht inert),
# 14 (Inhalt, Link-Ziel, Umbenennung, gelöschte und neue Datei), 15 (Symlink: Ziel, Art, Spiegelregel je Zieltyp),
# 16 (Identität gibt frei, auch ohne passenden Stand der schnellen Kette), 17 (Merge-Konflikt vor der Identität),
# 18 (untracked Datei), 19 (kein Stand, auch bei sauberem Baum), 20 (nicht lesbarer Vergleich, Abbruch, beschädigte Liste,
# Schleifenschutz vor jeder Prüfung, die zu `block` führt),
# 21 (Make-Flags beim Schreiben des Nachweises: -i entfernt ihn, -n, -q, -t lassen ihn stehen),
# 22 (Aufrufarten: Arbeitsverzeichnis im Unterverzeichnis, leere Umgebung, auch bei der Freigabe).
#
# GRENZE der Gegenprobe:
# - Der Schleifenschutz (`stop_hook_active: true`) antwortet `approve` auch bei Konflikt, fehlendem
#   oder beschädigtem Stand, scheiterndem git und fehlender dateiliste.sh: gewollte Grenze des
#   Hooks (SPEC-052 Punkt 20). Je Zustand ein Fall p20-schleifenschutz-*; `false` sperrt dort.
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
# - Drei Mutanten bleiben grün und werden als grün geprüft: `[ -f "$stand" ] || freigabe` (ohne Stand bricht die
#   Freigabe beim Lesen ab, `block`), die Streichung der Prüfung auf absolute Pfade (die Positivliste lehnt sie am
#   Anfang des Pfades ab) und `|| true` am Hash-Skript bei fehlendem Skript (die leere Ausgabe wird abgelehnt).
# - Die Seite des Arbeitsbaums hat für unzulässige Pfade keinen herstellbaren Fall (`git ls-files` liefert sie nicht).
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
L_SCHLUSS='freigabe'
L_HASH='  hash="$(bash tools/harness/working-tree-hash.sh)"'
L_GESPEICHERT='  gespeichert="$(cat -- .harness/state/gates-schnell-passed.diffsha 2> /dev/null || :)"'
L_SCHNELL_CMP='  [ "$hash" = "$gespeichert" ] || voll "The fast gates chain has no passed run for this working tree." '"'"'Run `make gates-schnell` (or `make gates`).'"'"
L_HASH_LEER='  [ -n "$hash" ] || voll "The working tree hash is empty."'
L_LINK='      [ "$art" != L ] || link_inert "$p" "${WERT["$seite:$p"]}" || voll "A changed symlink is not inert."'
L_ARTX='      [ "$art" != X ] || voll "A changed entry of type X is never inert."'
L_SEITEN='    for seite in s a; do'
L_SPIEGEL='  spiegel'
L_README='    README.md) return 0 ;;'

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
# Dateien des Temp-Repos auf der Positivliste und daneben (Pfade der Fälle der Punkte 9 bis 15).
DATEIEN_FIXTURE="docs/plan/planning/p.md docs/reviews/r.md docs/user/u.md docs/user/abdeckung-e2e.md
  docs/user/benutzerhandbuch.md README.md harness/README.md internal/x.go .golangci.yml tools/x.sh
  Dockerfile go.mod .dockerignore .d-check.yml spec/spezifikation.md
  harness/vertraege/SPEC-052-gate-kette.md"
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
  for d in $DATEIEN_FIXTURE; do
    mkdir -p "$(dirname "$w/$d")"
    printf 'inhalt\n' > "$w/$d"
  done
  ln -s u.md "$w/docs/user/ln"
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
# grund <fall> <teilstring> / kein_grund <fall> <teilstring> — meldet den Fall, wenn der `reason` der Antwort
# den Teilstring nicht (grund) oder doch (kein_grund) enthält; nur die Fälle des `reason` in Punkt 9 rufen sie.
grund() { case "$hout" in *'"reason":"'*"$2"*) ;; *) echo "$1" ;; esac; }
kein_grund() { case "$hout" in *"$2"*) echo "$1" ;; esac; }

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
  p14-geloescht-index-nicht-inert
  p9-freigabe p9-konflikt p9-stand-fehlt p9-stand-fehlt-nicht-inert p9-stand-beschaedigt p9-nicht-inert-internal
  p9-nicht-inert-makefile p9-nicht-inert-harness-mk p9-nicht-inert-golangci p9-nicht-inert-tools
  p9-nicht-inert-dockerfile p9-nicht-inert-go-mod p9-nicht-inert-dockerignore
  p9-nicht-inert-d-check p10-fehlt p10-fremd p10-leer p10-unlesbar p10-skript-fehlt
  p10-skript-scheitert p10-skript-leer p10-skript-fehlt-schleifenschutz p16-schnell-leer
  p16-schnell-fremd p12-planning-tief p12-reviews-tief p12-user-tief p12-readme p13-docs-user
  p12-grenze-planning-x p12-grenze-reviews-alt p12-grenze-users p12-grenze-datei-reviews
  p12-grenze-datei-user p12-grenze-readme-bak p12-case-planning p12-case-readme p12-harness-readme
  p12-selbst-hook p12-selbst-dateiliste p12-selbst-hash p12-selbst-spec p12-selbst-vertrag p12-unzulaessig-punktpunkt
  p12-unzulaessig-punkt p12-unzulaessig-leer p12-unzulaessig-leer-ende p12-unzulaessig-absolut p14-geloescht-inert
  p14-geloescht-nicht-inert p14-umbenennung-inert p14-umbenennung-nach-inert
  p14-umbenennung-von-inert p14-neu-nicht-inert p18-untracked-nicht-inert p12-art-x p15-link-inert
  p15-ziel-aussen p15-ziel-ueber-wurzel p15-ziel-absolut p15-ziel-absolut-aussen
  p15-ziel-ausserhalb-liste p15-ziel-verzeichnis p15-art-datei-zu-link p15-art-link-zu-datei
  p15-spiegel-liste p15-spiegel-verzeichnis p15-spiegel-vorfahr p15-spiegel-wurzel
  p15-spiegel-absolut p22-freigabe-unterverzeichnis p22-freigabe-leere-umgebung
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

# --- Positivliste und Stand der schnellen Kette (Punkte 9, 10, 12 bis 15, 18, 22) -----------
# aendere <datei…> — hängt der Datei in c eine Kommentarzeile an (legt sie bei Bedarf an).
aendere() { local d; for d; do mkdir -p "$(dirname "$c/$d")"; echo '# x' >> "$c/$d"; done; }
# schnell — schreibt den Stand der schnellen Kette in c: den Hash des Arbeitsbaums, wie die schnelle Kette ihn schreibt.
schnell() { mkdir -p "$c/.harness/state"; (cd "$c" && bash "$wurzel/$HASH") > "$c/.harness/state/gates-schnell-passed.diffsha" || :; }
# pfad_fall <fall> <erwartet> <datei…> — Stand des vollen Laufs, dann diese Dateien geändert, Stand der schnellen Kette grün.
pfad_fall() { local f="$1" e="$2"; shift 2; kopie "$f"; stand; aendere "$@"; schnell; hook; ist "$f" "$e"; }
# nach_stand <fall> <erwartet> <befehl…> — Stand des vollen Laufs, dann der Befehl (er liest c selbst), Stand der schnellen Kette grün.
nach_stand() { local f="$1" e="$2"; kopie "$f"; stand; shift 2; "$@"; schnell; hook; ist "$f" "$e"; }

f_p9_freigabe() { pfad_fall p9-freigabe approve docs/user/u.md; }
f_p9_konflikt() {
  kopie p9-konflikt
  merge_konflikt || { echo p9-konflikt; return; }
  stand; aendere docs/user/u.md; schnell; hook; ist p9-konflikt block
}
f_p9_stand_fehlt() { kopie p9-stand-fehlt; aendere docs/user/u.md; schnell; hook; ist p9-stand-fehlt block; }
# Der Stand der schnellen Kette fehlt, die Änderung ist nicht inert: der `reason` nennt die schnelle Kette nicht.
f_p9_stand_fehlt_nicht_inert() {
  kopie p9-stand-fehlt-nicht-inert; stand; aendere internal/x.go; hook
  ist p9-stand-fehlt-nicht-inert block; kein_grund p9-stand-fehlt-nicht-inert 'make gates-schnell'
}
f_p9_stand_beschaedigt() {
  local felder i
  kopie p9-stand-beschaedigt; stand
  mapfile -d '' -t felder < "$c/.harness/state/gates-passed.dateiliste"
  for ((i = 2; i < ${#felder[@]}; i += 3)); do
    [ "${felder[i]}" != docs/user/u.md ] || felder[i - 1]=abc
  done
  printf '%s\0' "${felder[@]}" > "$c/.harness/state/gates-passed.dateiliste"
  schnell; hook; ist p9-stand-beschaedigt block
}
f_p9_nicht_inert_internal() { pfad_fall p9-nicht-inert-internal block internal/x.go; }
f_p9_nicht_inert_makefile() { pfad_fall p9-nicht-inert-makefile block Makefile; }
f_p9_nicht_inert_harness_mk() { pfad_fall p9-nicht-inert-harness-mk block harness/mk/fake.mk; }
f_p9_nicht_inert_golangci() { pfad_fall p9-nicht-inert-golangci block .golangci.yml; }
f_p9_nicht_inert_tools() { pfad_fall p9-nicht-inert-tools block tools/x.sh; }
f_p9_nicht_inert_dockerfile() { pfad_fall p9-nicht-inert-dockerfile block Dockerfile; }
f_p9_nicht_inert_go_mod() { pfad_fall p9-nicht-inert-go-mod block go.mod; }
f_p9_nicht_inert_dockerignore() { pfad_fall p9-nicht-inert-dockerignore block .dockerignore; }
f_p9_nicht_inert_d_check() { pfad_fall p9-nicht-inert-d-check block .d-check.yml; }

f_p10_fehlt() {
  kopie p10-fehlt; stand; aendere docs/user/u.md; hook
  ist p10-fehlt block; grund p10-fehlt 'make gates-schnell'
}
f_p10_fremd() { kopie p10-fremd; stand; schnell; aendere docs/user/u.md; hook; ist p10-fremd block; }
f_p10_leer() {
  kopie p10-leer; stand; aendere docs/user/u.md
  mkdir -p "$c/.harness/state"; : > "$c/.harness/state/gates-schnell-passed.diffsha"
  hook; ist p10-leer block
}
f_p10_unlesbar() {
  if ! unlesbar_geht; then
    echo "stop-hook-gegenprobe: p10-unlesbar übersprungen (Dateien mit Rechten 000 sind lesbar)" >&2
    return
  fi
  kopie p10-unlesbar; stand; aendere docs/user/u.md; schnell
  chmod 000 "$c/.harness/state/gates-schnell-passed.diffsha"
  hook; ist p10-unlesbar block
}
# Das Hash-Skript der Kopie ist vor dem Stand des vollen Laufs geändert (kein Unterschied zu ihm); den
# Stand der schnellen Kette schreibt das unveränderte Skript des Arbeitsbaums.
f_p10_skript_fehlt() {
  kopie p10-skript-fehlt; rm -f "${c:?}/${HASH:?}"; stand; aendere docs/user/u.md; schnell
  hook; ist p10-skript-fehlt block
}
f_p10_skript_scheitert() {
  kopie p10-skript-scheitert; printf 'exit 1\n' >> "$c/$HASH"; stand; aendere docs/user/u.md; schnell
  hook; ist p10-skript-scheitert block
}
f_p10_skript_leer() {
  kopie p10-skript-leer; printf 'exit 0\n' > "$c/$HASH"; stand; aendere docs/user/u.md
  hook; ist p10-skript-leer block
}
f_p10_skript_fehlt_schleifenschutz() {
  kopie p10-skript-fehlt-schleifenschutz; rm -f "${c:?}/${HASH:?}"; stand; aendere docs/user/u.md; schnell
  hook '{"stop_hook_active": true}'; ist p10-skript-fehlt-schleifenschutz approve
}

# Identität (Punkt 16): ein leerer Diff gibt auch frei, wenn der Stand der schnellen Kette leer oder fremd ist.
f_p16_schnell_leer() {
  kopie p16-schnell-leer; stand; mkdir -p "$c/.harness/state"; : > "$c/.harness/state/gates-schnell-passed.diffsha"
  hook; ist p16-schnell-leer approve
}
f_p16_schnell_fremd() {
  kopie p16-schnell-fremd; stand; mkdir -p "$c/.harness/state"
  echo fremd > "$c/.harness/state/gates-schnell-passed.diffsha"
  hook; ist p16-schnell-fremd approve
}

# Positivliste (Punkt 12): jeder Eintrag, in jeder Tiefe.
f_p12_planning_tief() { pfad_fall p12-planning-tief approve docs/plan/planning/a/b/c.md; }
f_p12_reviews_tief() { pfad_fall p12-reviews-tief approve docs/reviews/a/b/c.md; }
f_p12_user_tief() { pfad_fall p12-user-tief approve docs/user/a/b/c.md; }
f_p12_readme() { pfad_fall p12-readme approve README.md; }
f_p13_docs_user() { pfad_fall p13-docs-user approve docs/user/benutzerhandbuch.md docs/user/abdeckung-e2e.md; }
# Grenze der Einträge (Punkt 12): der Pfad daneben gibt nicht frei.
f_p12_grenze_planning_x() { pfad_fall p12-grenze-planning-x block docs/plan/planning-x/n.md; }
f_p12_grenze_reviews_alt() { pfad_fall p12-grenze-reviews-alt block docs/reviews-alt/n.md; }
f_p12_grenze_users() { pfad_fall p12-grenze-users block docs/users/n.md; }
f_p12_grenze_datei_reviews() { nach_stand p12-grenze-datei-reviews block ersetze_durch_datei docs/reviews; }
f_p12_grenze_datei_user() { nach_stand p12-grenze-datei-user block ersetze_durch_datei docs/user; }
# ersetze_durch_datei <pfad> — das Verzeichnis des Pfades in c wird eine Datei.
ersetze_durch_datei() { rm -rf "${c:?}/${1:?}"; echo x > "$c/$1"; }
f_p12_grenze_readme_bak() { pfad_fall p12-grenze-readme-bak block README.md.bak; }
f_p12_case_planning() { pfad_fall p12-case-planning block docs/Plan/planning/n.md; }
f_p12_case_readme() { pfad_fall p12-case-readme block readme.md; }
f_p12_harness_readme() { pfad_fall p12-harness-readme block harness/README.md; }
# Selbständerung (Punkt 12): Hook, Skripte der Liste, Spezifikation und der Vertrag
# harness/vertraege/SPEC-052-gate-kette.md sind nicht inert.
f_p12_selbst_hook() { pfad_fall p12-selbst-hook block "$HOOK"; }
f_p12_selbst_dateiliste() { pfad_fall p12-selbst-dateiliste block "$LIB"; }
f_p12_selbst_hash() { pfad_fall p12-selbst-hash block "$HASH"; }
f_p12_selbst_spec() { pfad_fall p12-selbst-spec block spec/spezifikation.md; }
f_p12_selbst_vertrag() { pfad_fall p12-selbst-vertrag block harness/vertraege/SPEC-052-gate-kette.md; }

# Pfad der Differenz unzulässig (Punkt 12): ein Eintrag im Stand, den der Arbeitsbaum nicht hat.
unzulaessig() {
  kopie "$1"; stand
  printf 'F\0%s\0%s\0' "$H64" "$2" >> "$c/.harness/state/gates-passed.dateiliste"
  schnell; hook; ist "$1" block
}
f_p12_unzulaessig_punktpunkt() { unzulaessig p12-unzulaessig-punktpunkt 'docs/plan/planning/../../internal/x.go'; }
f_p12_unzulaessig_punkt() { unzulaessig p12-unzulaessig-punkt 'docs/plan/planning/./x.md'; }
f_p12_unzulaessig_leer() { unzulaessig p12-unzulaessig-leer 'docs/plan/planning//x.md'; }
f_p12_unzulaessig_leer_ende() { unzulaessig p12-unzulaessig-leer-ende 'docs/user/x/'; }
f_p12_unzulaessig_absolut() { unzulaessig p12-unzulaessig-absolut '/docs/user/x.md'; }

# Gelöscht, umbenannt, neu, untracked (Punkte 14, 18) und Art X (Punkt 12).
loesche() { rm -f "${c:?}/$1"; }
f_p14_geloescht_inert() { nach_stand p14-geloescht-inert approve loesche docs/user/u.md; }
f_p14_geloescht_nicht_inert() { nach_stand p14-geloescht-nicht-inert block loesche internal/x.go; }
git_rm() { git -C "$c" rm -q "$1"; }
f_p14_geloescht_index_nicht_inert() { nach_stand p14-geloescht-index-nicht-inert block git_rm internal/x.go; }
git_mv() { git -C "$c" mv "$1" "$2"; }
f_p14_umbenennung_inert() { nach_stand p14-umbenennung-inert approve git_mv docs/user/u.md docs/user/v.md; }
f_p14_umbenennung_nach_inert() { nach_stand p14-umbenennung-nach-inert block git_mv internal/x.go docs/user/x.go; }
f_p14_umbenennung_von_inert() { nach_stand p14-umbenennung-von-inert block git_mv docs/user/u.md internal/u.go; }
f_p14_neu_nicht_inert() { nach_stand p14-neu-nicht-inert block neu_getrackt internal/n.go; }
neu_getrackt() { echo n > "$c/$1"; git -C "$c" add "$1"; }
f_p18_untracked_nicht_inert() { nach_stand p18-untracked-nicht-inert block neu_untracked internal/n.go; }
neu_untracked() { echo n > "$c/$1"; }
f_p12_art_x() {
  kopie p12-art-x
  mkdir -p "$c/docs/user/sub"
  git -C "$c" update-index --add --cacheinfo "160000,$(git -C "$c" rev-parse HEAD),docs/user/sub"
  stand; rmdir "$c/docs/user/sub"; schnell; hook; ist p12-art-x block
}

# Symlink (Punkt 15): der Link docs/user/ln zeigt anders; Link-Pfad und Ziel müssen inert sein.
ziel_neu() { ln -sfn "$1" "$c/docs/user/ln"; }
f_p15_link_inert() { nach_stand p15-link-inert approve ziel_neu abdeckung-e2e.md; }
f_p15_ziel_aussen() { nach_stand p15-ziel-aussen block ziel_neu ../../../docs/user/u.md; }
f_p15_ziel_ueber_wurzel() { nach_stand p15-ziel-ueber-wurzel block ziel_neu ../../../../../docs/user/u.md; }
f_p15_ziel_absolut() { kopie p15-ziel-absolut; stand; ziel_neu "$c/docs/user/u.md"; schnell; hook; ist p15-ziel-absolut block; }
f_p15_ziel_absolut_aussen() { nach_stand p15-ziel-absolut-aussen block ziel_neu /docs/user/u.md; }
f_p15_ziel_ausserhalb_liste() { nach_stand p15-ziel-ausserhalb-liste block ziel_neu ../../internal/x.go; }
f_p15_ziel_verzeichnis() { nach_stand p15-ziel-verzeichnis block ziel_neu .; }
datei_zu_link() { rm -f "${c:?}/docs/user/u.md"; ln -s ../../internal/x.go "$c/docs/user/u.md"; }
f_p15_art_datei_zu_link() { nach_stand p15-art-datei-zu-link block datei_zu_link; }
f_p15_art_link_zu_datei() {
  kopie p15-art-link-zu-datei; ziel_neu ../../internal/x.go; stand
  rm -f "${c:?}/docs/user/ln"; echo x > "$c/docs/user/ln"; schnell; hook; ist p15-art-link-zu-datei block
}
# Spiegelregel: ein Link außerhalb der Positivliste mit Ziel darin macht die Änderung einer inerten Datei nicht inert.
spiegel_fall() {
  local ziel="$2"
  kopie "$1"; [ "$ziel" != ABS ] || ziel="$c/docs/user"
  ln -s "$ziel" "$c/$3"; stand
  aendere docs/plan/planning/p.md; schnell; hook; ist "$1" block
}
f_p15_spiegel_liste() { spiegel_fall p15-spiegel-liste docs/user/u.md sp1; }
f_p15_spiegel_verzeichnis() { spiegel_fall p15-spiegel-verzeichnis docs/user sp2; }
f_p15_spiegel_vorfahr() { spiegel_fall p15-spiegel-vorfahr docs sp3; }
f_p15_spiegel_wurzel() { spiegel_fall p15-spiegel-wurzel .. dir/up; }
f_p15_spiegel_absolut() { spiegel_fall p15-spiegel-absolut ABS sp5; }

# Aufrufarten bei der Freigabe (Punkt 22).
f_p22_freigabe_unterverzeichnis() {
  kopie p22-freigabe-unterverzeichnis; stand; aendere docs/user/u.md; schnell; cwd="$c/dir"
  hook; ist p22-freigabe-unterverzeichnis approve
}
f_p22_freigabe_leere_umgebung() {
  kopie p22-freigabe-leere-umgebung; stand; aendere docs/user/u.md; schnell; leer=1
  hook; ist p22-freigabe-leere-umgebung approve
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
ersetze $HOOK "$L_STANDFEHLT" '[ -f "$stand" ] || frei'; rot "p19-ohne-stand-sauber p19-ohne-stand-geaendert p9-stand-fehlt"
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

# --- Punkte 9, 10, 12 bis 15, 18, 22: Positivliste und Stand der schnellen Kette -----------------
UNL=""
if unlesbar_geht; then UNL=p10-unlesbar; fi

# Reihenfolge (Punkt 9): die Freigabe läuft vor der Prüfung; der Fall des Schritts allein wird rot.
# vor <zeile> — die Freigabe (Positivliste, Stand der schnellen Kette) steht vor dieser Zeile des Hooks.
vor() { ersetze $HOOK "$1" 'stand="$DATEILISTE_STAND"; aktuell="$(mktemp)"; dateiliste_erzeugen > "$aktuell"; freigabe'$'\n'"$1"; }
fixture m-pfad-vor-konflikt
vor "$L_KONFLIKT"; rot p9-konflikt "p9-stand-fehlt"
fixture m-pfad-vor-stand-fehlt
ersetze $HOOK "$L_STANDFEHLT" '[ -f "$stand" ] || { aktuell="$(mktemp)"; dateiliste_erzeugen > "$aktuell"; freigabe; }'
rot "" "p9-konflikt p9-stand-fehlt p9-stand-beschaedigt"
fixture m-pfad-vor-pruefung
vor 'dateiliste_pruefen "$stand" || voll "The recorded file list is unreadable or damaged."'
rot p9-stand-beschaedigt "p9-konflikt p9-stand-fehlt"

# `reason` bei fehlendem Stand der schnellen Kette (Punkt 9): nur dort nennt er die schnelle Kette.
fixture m-hinweis-immer
ersetze $HOOK "  local hinweis='Run \`make gates\`.'" "  local hinweis='Run \`make gates-schnell\` (or \`make gates\`).'"; rot p9-stand-fehlt-nicht-inert p10-fehlt
fixture m-hinweis-nie
ersetze $HOOK "$L_SCHNELL_CMP" '  [ "$hash" = "$gespeichert" ] || voll "The fast gates chain has no passed run for this working tree."'; rot p10-fehlt p9-stand-fehlt-nicht-inert
fixture m-stand-vor-pfadregel
ersetze $HOOK '  pfadregel' '  schnell_gruen'; sedmut $HOOK '0,/^  schnell_gruen$/! s/^  schnell_gruen$/  pfadregel/'
rot p9-stand-fehlt-nicht-inert p10-fehlt

# Stand der schnellen Kette (Punkt 10).
fixture m-schnell-nicht-geprueft
ersetze $HOOK "$L_SCHNELL_CMP" ''; rot "p10-fehlt p10-fremd p10-leer $UNL"
fixture m-schnell-nur-vorhanden
ersetze $HOOK "$L_SCHNELL_CMP" '  [ -n "$gespeichert" ] || voll "x" '"'"'Run `make gates-schnell` (or `make gates`).'"'"; rot p10-fremd "p10-fehlt p10-leer"
fixture m-schnell-leer-gilt
ersetze $HOOK "$L_GESPEICHERT" "$L_GESPEICHERT"$'\n''  [ ! -e .harness/state/gates-schnell-passed.diffsha ] || [ -s .harness/state/gates-schnell-passed.diffsha ] || gespeichert="$hash"'
rot p10-leer "p10-fehlt p10-fremd"
fixture m-schnell-unlesbar-gilt
ersetze $HOOK "$L_GESPEICHERT" '  gespeichert="$(cat -- .harness/state/gates-schnell-passed.diffsha 2> /dev/null || printf "%s" "$hash")"'
rot "p10-fehlt $UNL" "p10-leer p10-fremd"
fixture m-skript-ueberlesen
ersetze $HOOK "$L_HASH" "$L_HASH || true"; rot p10-skript-scheitert "p10-skript-fehlt p10-skript-leer"
fixture m-skript-fehlt-ersatz
ersetze $HOOK "$L_HASH" '  hash="$(bash tools/harness/working-tree-hash.sh 2> /dev/null || cat .harness/state/gates-schnell-passed.diffsha)"'
rot p10-skript-fehlt "p10-skript-leer p10-skript-scheitert"
fixture m-hash-leer-gilt
ersetze $HOOK "$L_HASH_LEER" ''; rot p10-skript-leer "p10-skript-fehlt p10-skript-scheitert"
fixture m-schutz-nach-hash
schutz_nach "$L_HASH"; rot p10-skript-fehlt-schleifenschutz
fixture m-schnell-bei-identitaet
ersetze $HOOK "$L_CMP" 'if cmp -s "$stand" "$aktuell"; then schnell_gruen'; rot "p16-identisch p16-schnell-leer p16-schnell-fremd"

# Aufrufarten bei der Freigabe (Punkt 22).
fixture m-22-hash-relativ
ersetze $HOOK "$L_HASH" '  hash="$(bash "$OLDPWD/tools/harness/working-tree-hash.sh")"'
rot p22-freigabe-unterverzeichnis "p22-freigabe-leere-umgebung p9-freigabe"
fixture m-22-hash-home
ersetze $HOOK "$L_HASH" '  hash="$(: "${HOME:?}"; bash tools/harness/working-tree-hash.sh)"'
rot p22-freigabe-leere-umgebung "p22-freigabe-unterverzeichnis p9-freigabe"

# Positivliste (Punkt 12): ein Pfad mehr auf der Liste wird rot; die Liste steht an einer Stelle.
NA=0
# aufnahme <muster> <rote Fälle> — das Muster steht zusätzlich in der Positivliste.
aufnahme() {
  NA=$((NA + 1))
  fixture "m-aufnahme-$NA"
  ersetze $HOOK "$L_README" "$L_README"$'\n'"    $1) return 0 ;;"
  rot "$2"
}
aufnahme 'spec/*' p12-selbst-spec
aufnahme 'harness/vertraege/*' p12-selbst-vertrag
aufnahme '.claude/hooks/stop-require-gates.sh' p12-selbst-hook
aufnahme 'tools/harness/dateiliste.sh' p12-selbst-dateiliste
aufnahme 'tools/harness/working-tree-hash.sh' p12-selbst-hash
aufnahme 'tools/*' "p9-nicht-inert-tools p12-selbst-dateiliste p12-selbst-hash"
aufnahme 'internal/*' p9-nicht-inert-internal
aufnahme 'Makefile' p9-nicht-inert-makefile
aufnahme 'harness/mk/*' p9-nicht-inert-harness-mk
aufnahme '.golangci.yml' p9-nicht-inert-golangci
aufnahme 'Dockerfile' p9-nicht-inert-dockerfile
aufnahme 'go.mod' p9-nicht-inert-go-mod
aufnahme '.dockerignore' p9-nicht-inert-dockerignore
aufnahme '.d-check.yml' p9-nicht-inert-d-check
# Grenze je Eintrag (Punkt 12): der Pfad daneben gibt nicht frei.
aufnahme 'docs/plan/planning-*' p12-grenze-planning-x
aufnahme 'docs/reviews-*' p12-grenze-reviews-alt
aufnahme 'docs/users/*' p12-grenze-users
aufnahme 'docs/reviews' p12-grenze-datei-reviews
aufnahme 'docs/user' p12-grenze-datei-user
aufnahme 'README.md.bak' p12-grenze-readme-bak
aufnahme 'harness/README.md' p12-harness-readme

# Je Eintrag der Positivliste (Punkt 12): ohne ihn gibt sein Pfad nicht frei; nur die erste Ebene gibt nicht in jeder Tiefe frei.
L_PLANNING='    docs/plan/planning/*) return 0 ;;'
L_REVIEWS='    docs/reviews/*) return 0 ;;'
L_USER='    docs/user/*) return 0 ;;'
fixture m-ohne-planning
ersetze $HOOK "$L_PLANNING" ''; rot p12-planning-tief
fixture m-ohne-reviews
ersetze $HOOK "$L_REVIEWS" ''; rot p12-reviews-tief
fixture m-ohne-user
ersetze $HOOK "$L_USER" ''; rot "p12-user-tief p13-docs-user p9-freigabe"
fixture m-ohne-readme
ersetze $HOOK "$L_README" ''; rot p12-readme
fixture m-planning-erste-ebene
ersetze $HOOK "$L_PLANNING" '    docs/plan/planning/*/*) ;;'$'\n'"$L_PLANNING"; rot p12-planning-tief p9-freigabe
fixture m-reviews-erste-ebene
ersetze $HOOK "$L_REVIEWS" '    docs/reviews/*/*) ;;'$'\n'"$L_REVIEWS"; rot p12-reviews-tief
fixture m-user-erste-ebene
ersetze $HOOK "$L_USER" '    docs/user/*/*) ;;'$'\n'"$L_USER"; rot p12-user-tief p9-freigabe
fixture m-case-egal
ersetze $HOOK 'export LC_ALL=C' 'export LC_ALL=C; shopt -s nocasematch'; rot "p12-case-planning p12-case-readme" p9-freigabe
fixture m-readme-jede-tiefe
ersetze $HOOK "$L_README" '    README.md | */README.md) return 0 ;;'; rot p12-harness-readme p12-readme

# Unzulässiger Pfad in der Differenz (Punkt 12): je Form ohne ihre Prüfung.
fixture m-ohne-punktpunkt
ersetze $HOOK '  case "/$1/" in */../*) return 1 ;; esac' ''; rot p12-unzulaessig-punktpunkt "p12-unzulaessig-punkt p12-unzulaessig-leer p12-unzulaessig-leer-ende p12-unzulaessig-absolut"
fixture m-ohne-punkt
ersetze $HOOK '  case "/$1/" in */./*) return 1 ;; esac' ''; rot p12-unzulaessig-punkt "p12-unzulaessig-punktpunkt p12-unzulaessig-leer p12-unzulaessig-leer-ende p12-unzulaessig-absolut"
fixture m-ohne-leeres-segment
ersetze $HOOK '  case "$1" in "" | *//* | */) return 1 ;; esac' ''; rot "p12-unzulaessig-leer p12-unzulaessig-leer-ende" "p12-unzulaessig-punktpunkt p12-unzulaessig-punkt p12-unzulaessig-absolut"
fixture m-ohne-leeres-segment-mitte
ersetze $HOOK '  case "$1" in "" | *//* | */) return 1 ;; esac' '  case "$1" in "" | */) return 1 ;; esac'; rot p12-unzulaessig-leer "p12-unzulaessig-leer-ende p12-unzulaessig-punkt p12-unzulaessig-punktpunkt p12-unzulaessig-absolut"
fixture m-ohne-leeres-segment-ende
ersetze $HOOK '  case "$1" in "" | *//* | */) return 1 ;; esac' '  case "$1" in "" | *//*) return 1 ;; esac'; rot p12-unzulaessig-leer-ende "p12-unzulaessig-leer p12-unzulaessig-punkt p12-unzulaessig-punktpunkt p12-unzulaessig-absolut"
fixture m-ohne-absolut
ersetze $HOOK '  case "$1" in /*) return 1 ;; esac' ''; rot "" "p12-unzulaessig-absolut p12-unzulaessig-punktpunkt"
fixture m-absolut-relativ
ersetze $HOOK '  case "$1" in /*) return 1 ;; esac' ''
sedmut $HOOK '/^inert() {/,/^}/ s/case "\$1" in/case "${1#\/}" in/'; rot p12-unzulaessig-absolut "p12-unzulaessig-punktpunkt p12-unzulaessig-punkt"

# Gelöscht, umbenannt, neu, untracked, Art X (Punkte 12, 14, 18).
fixture m-untracked-inert
sedmut $HOOK '/^inert() {/a\  [ -n "$(git ls-files --error-unmatch -- "$1" 2> /dev/null)" ] || return 0'
rot p18-untracked-nicht-inert p14-neu-nicht-inert
fixture m-nur-arbeitsbaum
sedmut $HOOK '/^diff_bilden() {/,/^}/ { /^  for p in \${SP/,/^  done/ d }'
rot "p14-geloescht-index-nicht-inert p14-umbenennung-nach-inert" "p14-geloescht-nicht-inert p14-umbenennung-von-inert p14-neu-nicht-inert"
fixture m-nur-stand
sedmut $HOOK '/^diff_bilden() {/,/^}/ s/for p in \${AP\[@\]+"\${AP\[@\]}"}; do/for p in ${SP[@]+"${SP[@]}"}; do/'
rot "p14-umbenennung-von-inert p14-neu-nicht-inert p18-untracked-nicht-inert" "p14-geloescht-nicht-inert p14-geloescht-index-nicht-inert p14-umbenennung-nach-inert"
fixture m-art-x-wie-datei
ersetze $HOOK "$L_ARTX" ''; rot p12-art-x

# Symlink (Punkt 15): je Element nur dieses Element gelockert.
fixture m-link-ohne-ziel
ersetze $HOOK "$L_LINK" ''
rot "p15-ziel-aussen p15-ziel-ueber-wurzel p15-ziel-absolut p15-ziel-absolut-aussen p15-ziel-ausserhalb-liste p15-ziel-verzeichnis p15-art-datei-zu-link p15-art-link-zu-datei" p15-link-inert
fixture m-verzeichnis-inert
ersetze $HOOK '  inert "$ZIEL"' '  inert "$ZIEL" || inert "$ZIEL/"'; rot p15-ziel-verzeichnis "p15-ziel-ausserhalb-liste p15-link-inert"
fixture m-ueber-wurzel-bleibt-innen
ersetze $HOOK "      ..) [ \"\${#STAPEL[@]}\" -gt 0 ] || return 1; unset 'STAPEL[\${#STAPEL[@]}-1]' ;;" "      ..) if [ \"\${#STAPEL[@]}\" -gt 0 ]; then unset 'STAPEL[\${#STAPEL[@]}-1]'; fi ;;"
rot "p15-ziel-aussen p15-ziel-ueber-wurzel" "p15-ziel-ausserhalb-liste p15-link-inert"
fixture m-absolut-wie-relativ
ersetze $HOOK '      ABS=1' ''; rot p15-ziel-absolut "p15-ziel-absolut-aussen p15-ziel-ausserhalb-liste"
fixture m-absolut-ohne-zweig
ersetze $HOOK '    /*)' '    /nichtvorhanden*)'; rot "p15-ziel-absolut p15-ziel-absolut-aussen" "p15-ziel-ausserhalb-liste p15-link-inert"
fixture m-art-nur-stand
ersetze $HOOK "$L_SEITEN" '    for seite in s; do'; rot p15-art-datei-zu-link p15-art-link-zu-datei
fixture m-art-nur-arbeitsbaum
ersetze $HOOK "$L_SEITEN" '    for seite in a; do'; rot p15-art-link-zu-datei p15-art-datei-zu-link

# Spiegelregel (Punkt 15): ganz entfernt, und je Zieltyp.
SPIEGEL_ALLE="p15-spiegel-liste p15-spiegel-verzeichnis p15-spiegel-vorfahr p15-spiegel-wurzel p15-spiegel-absolut"
fixture m-ohne-spiegel
ersetze $HOOK "$L_SPIEGEL" ''; rot "$SPIEGEL_ALLE"
fixture m-spiegel-ohne-liste
ersetze $HOOK '  inert "$1" && return 0' ''; rot p15-spiegel-liste "p15-spiegel-verzeichnis p15-spiegel-vorfahr p15-spiegel-wurzel p15-spiegel-absolut"
fixture m-spiegel-ohne-verzeichnis
ersetze $HOOK '    docs/plan/planning | docs/reviews | docs/user) return 0 ;;' ''
rot "p15-spiegel-verzeichnis p15-spiegel-absolut" "p15-spiegel-liste p15-spiegel-vorfahr p15-spiegel-wurzel"
fixture m-spiegel-ohne-vorfahr
ersetze $HOOK '    docs | docs/plan) return 0 ;;' ''; rot p15-spiegel-vorfahr "p15-spiegel-liste p15-spiegel-verzeichnis p15-spiegel-wurzel"
fixture m-spiegel-ohne-wurzel
ersetze $HOOK '    "") return 0 ;;' ''; rot p15-spiegel-wurzel "p15-spiegel-liste p15-spiegel-verzeichnis p15-spiegel-vorfahr"
fixture m-spiegel-ohne-absolut
sedmut $HOOK 's/^    \[ "\$LAGE" = innen \] || continue$/    [ "$LAGE" = innen ] \&\& [ "$ABS" -eq 0 ] || continue/'
rot p15-spiegel-absolut "p15-spiegel-liste p15-spiegel-verzeichnis p15-spiegel-vorfahr p15-spiegel-wurzel"

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
echo "stop-hook-gegenprobe: gruen — ${#FAELLE[@]} Fälle der Punkte 1, 9 bis 22 je Mutant rot, unveränderte Kopie grün, Anschluss an GATE_CHECKS"
