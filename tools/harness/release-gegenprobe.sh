#!/usr/bin/env bash
# release-gegenprobe — prüft die Zusagen von SPEC-067 Fassung A (harness/vertraege/SPEC-067-release.md, Punkt 14): den Text von .github/workflows/release.yml
# (Fälle p1-*, geprüft von der Funktion pruefe_release, die den Text hält) und tools/harness/release.sh (Fälle p2-* bis p13-*): je Fall ein Temp-Baum mit einem
# Git-Repository (Commits, Zweig, refs/remotes/origin/main und origin/trunk), dem Auszug des Benutzerhandbuchs, einem PATH, in dem `make`, `gh`, `docker` und
# `skopeo` Attrappen sind, die ihre Aufrufe aufzeichnen; `git`, `sha256sum`, `awk` und `cmp` sind die echten Werkzeuge. Ohne Docker, ohne Netz, ohne make. Ein Fall besteht,
# wenn Exit-Code, Kennwort der Meldung und die Aufrufe der Attrappen dem erwarteten Text entsprechen; jeder Fall hat seinen Temp-Baum. Je Mutation (Textersatz in
# pruefe_release bzw. in einer Kopie von release.sh) müssen genau die benannten Fälle rot werden; die unveränderte Kopie besteht jeden Fall. Der Name jedes Falls trägt
# die Nummer des Punktes; ein Name in der Fallliste ohne Funktion f_<name> zählt als roter Fall.
# Abbruch (Punkt 13): Bei Ende, SIGINT und SIGTERM beendet das Skript erst seine Nachkommen und wartet auf sie, entfernt dann sein Temp-Verzeichnis und endet mit 130 bzw.
# 143; die Fälle p13-gegenprobe-sigint und -sigterm starten dafür eine Kopie des Skripts mit eigenem TMPDIR (die Kopie startet ohne GP_MUT, GP_ENTDECKE und GP_FALL und mit
# GP_PARALLEL=2, und GP_INNEN=1 schaltet in jedem ihrer Läufe die Fälle p13-gegenprobe-* und die Mutationen m-p13-gegenprobe-* ab; ein Lauf mit GP_INNEN=1 und
# GP_MUT=m-p13-gegenprobe fährt nichts und endet mit Exit 1, Lauf `rekursion-schutz`) und senden das Signal, sobald dort ein Fall läuft; `env --default-signal` stellt SIGINT und SIGTERM
# für diese Prozesse auf die Voreinstellung, auch wenn die Gegenprobe selbst in einer Umgebung läuft, die sie ignoriert (etwa als Hintergrundauftrag einer nicht interaktiven Shell). Dieselben Fälle p13-sigint und
# p13-sigterm prüfen release.sh: ein `make` als Attrappe, das Nachkommen startet, bekommt das Signal mit dem Skript. SIGKILL und SIGHUP sind nicht zugesagt.
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Läufe, deren Name passt (die unveränderte Kopie nur, wenn das Muster `ok` trifft); GP_ENTDECKE=1 nennt je Mutation die
# roten Fälle; GP_PARALLEL=<n> (Vorgabe 4) setzt die Zahl der gleichzeitigen Läufe; GP_FALL=<name> gibt für die unveränderte Kopie das Protokoll dieses Falls aus.
# Braucht bash (mit `wait -n`), git, GNU sed, grep, awk, cmp, pgrep, setsid, timeout, sha256sum und die Werkzeuge von coreutils; ohne Docker, ohne make und ohne Netz.
# Aufrufarten von `make` (-n, -q, -t, -i, MAKEFLAGS) haben keinen Fall: kein Hook und kein Stand liest das Ergebnis, und release.sh läuft nie unter `make` (SPEC-067 Punkt 13).
set -uo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="$(cd "$hier/../.." && pwd)"
arbeit="$(mktemp -d)"
selbst="$hier/$(basename "${BASH_SOURCE[0]}")"
raeume_auf() { chmod -R u+rwX "${arbeit:?}" 2> /dev/null; rm -rf "${arbeit:?}"; }
# nachkommen <pid> — alle Nachkommen, Eltern vor Kindern.
nachkommen() { local k; for k in $(pgrep -P "$1" 2> /dev/null); do echo "$k"; nachkommen "$k"; done; }
# beende_kinder — hält die Nachkommen an (bis die Menge feststeht), beendet sie und wartet, bis keiner mehr läuft.
beende_kinder() {
  local i p="" n k
  for i in 1 2 3 4 5; do
    n="$(nachkommen $$ | sort -u | tr '\n' ' ')"
    [ "$n" != "$p" ] || break
    p="$n"
    # shellcheck disable=SC2086
    kill -s STOP $p 2> /dev/null
  done
  # shellcheck disable=SC2086
  [ -z "$p" ] || { kill -s TERM $p 2> /dev/null; kill -s CONT $p 2> /dev/null; }
  for k in $p; do for i in $(seq 1 100); do kill -0 "$k" 2> /dev/null || break; sleep 0.1; done; done
  wait 2> /dev/null
}
abbruch_int() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 130; }
abbruch_term() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 143; }
trap raeume_auf EXIT
trap abbruch_int INT
trap abbruch_term TERM
: > "$arbeit/fehler"
: > "$arbeit/gefahren"
melde() { echo "release-gegenprobe: ROT — $*" >> "$arbeit/fehler"; }

for t in git sed grep awk cmp mktemp sort tr wc cat cp chmod rm env find ls seq sleep mkdir dirname basename head pgrep setsid timeout sha256sum; do
  command -v "$t" > /dev/null 2>&1 || { echo "release-gegenprobe: $t fehlt" >&2; exit 2; }
done
for f in .github/workflows/release.yml tools/harness/release.sh; do
  [ -f "$repo/$f" ] || { echo "release-gegenprobe: $f fehlt im Arbeitsbaum" >&2; exit 2; }
done
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1
export GIT_AUTHOR_NAME=gp GIT_AUTHOR_EMAIL=gp@example.invalid GIT_COMMITTER_NAME=gp GIT_COMMITTER_EMAIL=gp@example.invalid

# --- Punkt 1: der Text -------------------------------------------------------------------------------------------------------------------------------
# SOLL_ROH ist der Text von Punkt 1a, Zeile für Zeile, mit der Herkunft (A oder B) vor dem ersten `|`; HERKUNFT nennt die Herkunft, die diese Fassung führt. Die Zeile
# mit <SHA> und <Tag> prüft PIN_RE. Die Gegenprobe trägt den Text wörtlich als eigene Kopie.
HERKUNFT=A
mapfile -t SOLL_ROH << 'EOT'
A|name: release
A|on:
A|  push:
A|    tags:
A|      - 'v*'
A|  workflow_dispatch:
A|permissions: {}
A|concurrency:
A|  group: release
A|  cancel-in-progress: false
A|jobs:
A|  release:
A|    name: release
A|    runs-on: ubuntu-24.04
A|    timeout-minutes: 60
A|    permissions:
A|      contents: write
A|    env:
A|      RELEASE_EVENT: ${{ github.event_name }}
A|      RELEASE_REF: ${{ github.ref }}
A|      RELEASE_SHA: ${{ github.sha }}
A|      RELEASE_DEFAULT_BRANCH: ${{ github.event.repository.default_branch }}
A|      RELEASE_REPOSITORY: ${{ github.repository }}
A|    steps:
A|      - name: Checkout
A|        uses: actions/checkout@<SHA> # <Tag>
A|        with:
A|          fetch-depth: 0
A|          persist-credentials: false
A|      - name: Eingaben
A|        run: bash tools/harness/release.sh eingaben
A|      - name: Binaries
A|        run: bash tools/harness/release.sh binaries
A|      - name: GitHub-Release
A|        if: ${{ github.event_name == 'push' }}
A|        env:
A|          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
A|        run: bash tools/harness/release.sh github-release
EOT
SOLL=()
for z in "${SOLL_ROH[@]}"; do if [[ $HERKUNFT == *"${z%%|*}"* ]]; then SOLL+=("${z#*|}"); fi; done
N=${#SOLL[@]}
PIN=-1
BEDINGUNG=-1
for ((i = 0; i < N; i++)); do
  if [[ ${SOLL[i]} == *"actions/checkout@<SHA>"* ]]; then PIN=$i; fi
  if [[ ${SOLL[i]} == *"if: \${{ github.event_name == 'push' }}"* ]]; then BEDINGUNG=$i; fi
done
PIN_RE="^${SOLL[PIN]%%<SHA>*}[0-9a-f]{40} # [^ ].*\$"
SHA_PROBE=0123456789abcdef0123456789abcdef01234567
# Anfang der drei Schritte nach dem Checkout (Zeilen in SOLL)
for ((i = 0; i < N; i++)); do
  case "${SOLL[i]}" in
    *"name: Eingaben") S_E=$i ;;
    *"name: Binaries") S_B=$i ;;
    *"name: GitHub-Release") S_G=$i ;;
  esac
done

# pruefe_release <datei> — Exit 0, wenn die Datei außer Leerzeilen und Kommentarzeilen genau aus den Zeilen von SOLL besteht (Pin der Zeile PIN nach PIN_RE).
pruefe_release() {
  local datei="$1" z i j=0
  local -a ist=()
  while IFS= read -r z || [ -n "$z" ]; do
    [[ $z =~ ^\ *(#|$) ]] && continue
    ist+=("$z")
  done < "$datei"
  for ((i = 0; i < ${#SOLL[@]}; i++)); do
    if [ "$i" -eq "$PIN" ]; then
      [[ ${ist[j]-} =~ $PIN_RE ]] || return 1
    else
      [ "${ist[j]-}" = "${SOLL[i]}" ] || return 1
    fi
    j=$((j + 1))
  done
  [ "${#ist[@]}" -eq "$j" ]
}

# --- Bausteine der Fälle ---------------------------------------------------------------------------------------------------------------------------
# text — der Text von Punkt 1 mit einem gültigen Pin, eine Zeile je Zeile.
text() {
  local i
  for i in "${!SOLL[@]}"; do
    if [ "$i" -eq "$PIN" ]; then echo "${SOLL[i]/<SHA> # <Tag>/$SHA_PROBE # v7.0.1}"; else echo "${SOLL[i]}"; fi
  done
}
# urteil <ok|rot> — schreibt stdin nach $w/release.yml und vergleicht das Urteil von pruefe_release.
urteil() {
  local soll="$1" r=0
  mkdir -p "$w"
  cat > "$w/release.yml"
  pruefe_release "$w/release.yml" || r=$?
  if { [ "$soll" = ok ] && [ "$r" -ne 0 ]; } || { [ "$soll" = rot ] && [ "$r" -eq 0 ]; }; then
    bad=1
    printf 'erwartet: %s\ngelesen: Exit %s\n' "$soll" "$r" >> "$run/log.$fallname"
  fi
}
neu() { w="$run/t.$fallname"; mapfile -t t < <(text); }
zeilen() { printf '%s\n' "$@"; }
f_p1_gleich() { neu; mkdir -p "$w"; pruefe_release "$repo/.github/workflows/release.yml" || { bad=1; echo "release.yml des Arbeitsbaums weicht von Punkt 1 ab" >> "$run/log.$fallname"; }; }
f_p1_text_gleich() { neu; urteil ok < <(zeilen "${t[@]}"); }
p1_fehlt() { neu; unset "t[$(($1 - 1))]"; urteil rot < <(zeilen "${t[@]}"); }
p1_geaendert() {
  local l lead
  neu
  l="${t[$(($1 - 1))]}"; lead="${l%%[! ]*}"
  t[$(($1 - 1))]="${lead}Z${l:${#lead}+1}"
  urteil rot < <(zeilen "${t[@]}")
}
for ((n = 1; n <= N; n++)); do eval "f_p1_fehlt_$n() { p1_fehlt $n; }; f_p1_geaendert_$n() { p1_geaendert $n; }"; done
# Zusätze am Ende der Datei: <Wort, das die Mutation durchlässt>|<Zeile>
ZUSATZ=(
  "pull_request_target:|pull_request_target:"
  "secrets.ACTIONS_RUNNER_DEBUG|      X: \${{ secrets.ACTIONS_RUNNER_DEBUG }}"
  "actions/cache@|      - uses: actions/cache@$SHA_PROBE # v4"
  "environment:|    environment: release"
  "continue-on-error: true|        continue-on-error: true"
  "always()|        if: \${{ always() }}"
  "packages: write|      packages: write"
)
ZUSATZ_NAMEN=(pull-request-target secrets-debug schritt-cache environment continue-on-error always packages-write)
zusatz() { neu; urteil rot < <(zeilen "${t[@]}" "$1"); }
for k in "${!ZUSATZ[@]}"; do eval "f_p1_zusatz_${ZUSATZ_NAMEN[k]//-/_}() { zusatz \"\${ZUSATZ[$k]#*|}\"; }"; done
pin() { neu; t[PIN]="        $1"; urteil rot < <(zeilen "${t[@]}"); }
f_p1_pin_kurz() { pin "uses: actions/checkout@${SHA_PROBE:1} # v7.0.1"; }
f_p1_pin_gross() { pin "uses: actions/checkout@${SHA_PROBE^^} # v7.0.1"; }
f_p1_pin_tag_leer() { pin "uses: actions/checkout@$SHA_PROBE # "; }
f_p1_pin_ohne_tag() { pin "uses: actions/checkout@$SHA_PROBE"; }
f_p1_pin_andere_action() { pin "uses: actions/cache@$SHA_PROBE # v4"; }
f_p1_bedingung_github_release() { neu; unset "t[$BEDINGUNG]"; urteil rot < <(zeilen "${t[@]}"); }
# umgestellt <E|B|G>… — der Text mit den Schritten in der genannten Reihenfolge (E: Eingaben, B: Binaries, G: GitHub-Release; der Checkout bleibt vorn).
umgestellt() {
  local b
  zeilen "${t[@]:0:S_E}"
  for b in "$@"; do
    case "$b" in
      E) zeilen "${t[@]:S_E:S_B-S_E}" ;;
      B) zeilen "${t[@]:S_B:S_G-S_B}" ;;
      G) zeilen "${t[@]:S_G}" ;;
    esac
  done
}
f_p1_reihenfolge_eingaben_binaries() { neu; urteil rot < <(umgestellt B E G); }
f_p1_reihenfolge_binaries_release() { neu; urteil rot < <(umgestellt E G B); }
f_p1_leerzeile() { neu; urteil ok < <(zeilen "${t[@]:0:3}" "" "${t[@]:3}"); }
f_p1_kommentarzeile() { neu; urteil ok < <(zeilen "${t[@]:0:3}" "# Kommentar" "${t[@]:3}"); }
f_p1_kommentarzeile_eingerueckt() { neu; urteil ok < <(zeilen "${t[@]:0:3}" "    # Kommentar" "${t[@]:3}"); }

# --- Welt, Attrappen, Hilfen für release.sh -------------------------------------------------------------------------------------------------------------
TOKEN=ghs_GeheimerTokenWert0123456789
# attrappen <dir> — make, gh, docker und skopeo als Skripte, die ihre Aufrufe nach $GP_LOG/<name> schreiben; das Verhalten steuern GP_MAKE und GP_GH*.
attrappen() {
  local b="$1" n
  mkdir -p "$b"
  cat > "$b/make" << 'EOT'
#!/usr/bin/env bash
echo "make $*" >> "$GP_LOG/make"
[ "${1-}" = binaries ] || exit 0
dist=""; v=""
for a in "$@"; do case "$a" in DIST=*) dist="${a#DIST=}" ;; VERSION=*) v="${a#VERSION=}" ;; esac; done
if [ -e "$dist" ]; then echo "dist-vorher: ja" >> "$GP_LOG/make"; else echo "dist-vorher: nein" >> "$GP_LOG/make"; fi
case "${GP_MAKE-ok}" in
  haengt)
    d="$(ls -d "$TMPDIR"/tmp.* | head -n 1)"
    bash -c 'while :; do mkdir -p "$1"; : > "$1/spur"; sleep 0.1; done' _ "$d/spur" &
    : > "$GP_LOG/bereit"
    sleep 60 &
    wait
    exit 0 ;;
esac
mkdir -p "$dist"
for os in darwin linux windows; do for a in amd64 arm64; do
  e=""; [ "$os" = windows ] && e=.exe
  printf 'binary %s %s %s\n' "$v" "$os" "$a" > "$dist/pgwire-recorder_${v}_${os}_${a}$e"
done; done
[ "${GP_MAKE-ok}" != datei-fehlt ] || rm "$dist/pgwire-recorder_${v}_linux_arm64"
[ "${GP_MAKE-ok}" != datei-zu-viel ] || echo extra > "$dist/extra"
(cd "$dist" && sha256sum -- * > SHA256SUMS)
[ "${GP_MAKE-ok}" != summe-falsch ] || echo manipuliert >> "$dist/pgwire-recorder_${v}_linux_amd64"
[ "${GP_MAKE-ok}" != scheitert ] || exit 1
exit 0
EOT
  cat > "$b/gh" << 'EOT'
#!/usr/bin/env bash
echo "gh $*" >> "$GP_LOG/gh"
case "$1 $2" in
  "release view")
    case "${GP_GH-fehlt}" in
      fehlt) echo "release not found" >&2; exit 1 ;;
      fehler) echo "HTTP 502: Bad Gateway" >&2; exit 1 ;;
      da) for n in ${GP_GH_ASSETS-}; do echo "$n"; done; exit 0 ;;
    esac ;;
  "release create")
    while [ $# -gt 0 ]; do if [ "$1" = --notes-file ]; then cp "$2" "$GP_LOG/notes"; fi; shift; done
    exit 0 ;;
  "release download")
    dir=""; while [ $# -gt 0 ]; do if [ "$1" = --dir ]; then dir="$2"; fi; shift; done
    case "${GP_GH_SUMS-gleich}" in
      gleich) cp "$RUNNER_TEMP/release/dist/SHA256SUMS" "$dir/SHA256SUMS" ;;
      abweichend) echo "0000000000000000000000000000000000000000000000000000000000000000  pgwire-recorder_x" > "$dir/SHA256SUMS" ;;
    esac
    exit 0 ;;
esac
exit 0
EOT
  for n in docker skopeo; do
    printf '#!/usr/bin/env bash\necho "%s $*" >> "$GP_LOG/%s"\nexit 0\n' "$n" "$n" > "$b/$n"
  done
  chmod +x "$b"/*
}
# welt <dir> — legt quelle/ an (main: c1 c2; Zweig side mit s1 aus c1; refs/remotes/origin/main auf c1, origin/trunk auf s1) und bin/ mit den Attrappen.
welt() {
  local q="$1/quelle"
  mkdir -p "$q"
  git -C "$q" init -q -b main
  git -C "$q" commit -q --allow-empty -m c1
  c1="$(git -C "$q" rev-parse HEAD)"
  git -C "$q" checkout -q -b side
  git -C "$q" commit -q --allow-empty -m s1
  s1="$(git -C "$q" rev-parse HEAD)"
  git -C "$q" checkout -q main
  git -C "$q" commit -q --allow-empty -m c2
  c2="$(git -C "$q" rev-parse HEAD)"
  git -C "$q" update-ref refs/remotes/origin/main "$c1"
  git -C "$q" update-ref refs/remotes/origin/trunk "$s1"
  attrappen "$1/bin"
}
# handbuch <Version> <Modus> — der Auszug des Benutzerhandbuchs nach $w/repo. Modi: normal (eine Zeile der Version) · zeile (Zeilen mit Vorabversion davor und Nachfolger danach) ·
# abschnitt (dieselbe Version in einem anderen Abschnitt davor) · nur-anderer (die Version nur in einem anderen Abschnitt davor) · ueberschrift (die Version nur unter einer
# Überschrift, die nicht auf Änderungshistorie endet) · spaeter (die Version nur in einem Abschnitt danach) · leer (Zeile der Version mit leerer dritter Zelle) · ohne (keine Zeile).
handbuch() {
  local v="$1" m="${2:-normal}" f="$w/repo/docs/user/benutzerhandbuch.md"
  mkdir -p "$w/repo/docs/user"
  {
    printf '# Benutzerhandbuch\n\n'
    case "$m" in abschnitt | nur-anderer) printf '## 10. Support\n\n| %s | 2020-01-01 | Falscher-Abschnitt |\n\n' "$v" ;; esac
    if [ "$m" = ueberschrift ]; then printf '## 11. Änderungshistorie (Entwurf)\n\n'; else printf '## 11. Änderungshistorie\n\n'; fi
    printf '| Version | Datum | Änderung |\n|---|---|---|\n'
    case "$m" in
      normal | abschnitt | ueberschrift) printf '| %s | Datum-%s | Text-%s |\n' "$v" "$v" "$v" ;;
      zeile) printf '| %s-rc.1 | 2026-10-01 | Text-Vorabversion |\n| %s | Datum-%s | Text-%s |\n| %s0 | 2026-10-03 | Text-Nachfolger |\n' "$v" "$v" "$v" "$v" "$v" ;;
      leer) printf '| %s | Datum-%s | |\n' "$v" "$v" ;;
      nur-anderer | spaeter | ohne) printf '| 9.9.9 | 2026-10-02 | Anderes |\n' ;;
    esac
    if [ "$m" = spaeter ]; then printf '\n## 12. Anhang\n\n| Version | Datum | Änderung |\n|---|---|---|\n| %s | 2020-01-01 | Falscher-Abschnitt |\n' "$v"; fi
  } > "$f"
}
# kopie — der Baum des Falls: repo/ (Kopie der Welt), log/, runner/, tmp/.
kopie() { w="$run/t.$fallname"; rm -rf "$w"; mkdir -p "$w/log" "$w/runner" "$w/tmp"; cp -a "$run/welt/quelle" "$w/repo"; }
# fix [Version [Handbuch-Modus]] — setzt V, legt den Baum an, schreibt VERSION und das Handbuch.
fix() { V="${1:-0.1.0}"; kopie; printf '%s\n' "$V" > "$w/repo/VERSION"; handbuch "$V" "${2:-normal}"; }
# umg <lauf|trocken> [NAME=Wert | -NAME]… — setzt E: die Umgebung von release.sh, dann die Abweichungen.
umg() {
  local art="$1" a
  local -A v=([RELEASE_EVENT]=push [RELEASE_REF]="refs/tags/v$V" [RELEASE_SHA]="$c1" [RELEASE_DEFAULT_BRANCH]=main [RELEASE_REPOSITORY]=pt9912/pgwire-recorder
    [RUNNER_TEMP]="$w/runner" [TMPDIR]="$w/tmp" [GP_LOG]="$w/log" [GH_TOKEN]="$TOKEN")
  shift
  if [ "$art" = trocken ]; then v[RELEASE_EVENT]=workflow_dispatch; v[RELEASE_REF]=refs/heads/impl; fi
  for a in "$@"; do if [[ $a == -* ]]; then unset "v[${a#-}]"; else v[${a%%=*}]="${a#*=}"; fi; done
  E=()
  for a in "${!v[@]}"; do E+=("$a=${v[$a]}"); done
}
umg_befehl() { ENVV=(env --default-signal=INT,TERM -i "PATH=$run/welt/bin:/usr/local/bin:/usr/bin:/bin" GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 "${E[@]}"); }
# lauf <Schritt> [NAME=Wert]… — führt release.sh mit E (und den Zusätzen) im Klon des Falls aus und setzt code, out, err (stdout und stderr unverändert).
lauf() {
  local s="$1"
  shift
  umg_befehl
  ( cd "$w/repo" && "${ENVV[@]}" "$@" bash "$releasesh" "$s" ) > "$w.out" 2> "$w.err"
  code=$?
  out="$(cat "$w.out"; printf x)"; out="${out%x}"
  err="$(cat "$w.err"; printf x)"; err="${err%x}"
}
# sei <Beschreibung> <Befehl>… — der Befehl muss Exit 0 haben, sonst ist der Fall rot und das Protokoll nennt die Beschreibung mit Exit, stdout und stderr.
sei() {
  local b="$1"
  shift
  if ! "$@"; then
    bad=1
    printf 'verletzt: %s\nExit %s\nstdout:\n%s\nstderr:\n%s\n' "$b" "${code-?}" "${out-}" "${err-}" >> "$run/log.$fallname"
  fi
}
enthaelt() { [[ $1 == *"$2"* ]]; }
ohne() { [[ $1 != *"$2"* ]]; }
# grund <Exit> <Schritt> <Kennwort> — Exit-Code und eine Meldung `release: <Schritt>: <Kennwort>:` auf stderr.
grund() { sei "Exit $1" test "$code" -eq "$1"; sei "Meldung release: $3:" enthaelt "$err" "release: $3:"; }
# beginnt <Schlüsselwort> — die erste Zeile von stderr beginnt mit `release: <Schlüsselwort>: `.
beginnt() { [[ ${err%%$'\n'*} == "release: $1: "* ]]; }
# hat_hinweis <Text> — eine Zeile von stderr beginnt mit `release: hinweis: ` und nennt den Text.
hat_hinweis() { local z; while IFS= read -r z; do if [[ $z == "release: hinweis: "*"$1"* ]]; then return 0; fi; done <<< "$err"; return 1; }
gruen() { sei "Exit 0" test "$code" -eq 0; }
# bis <Schritt>… — fährt die Schritte nacheinander mit E; jeder endet mit Exit 0.
bis() { local s; for s in "$@"; do lauf "$s"; sei "Schritt $s grün" test "$code" -eq 0; done; }
aufrufe() { cat "$w/log/$1" 2> /dev/null || true; }
nie() { [ ! -s "$w/log/$1" ]; }
DIST() { echo "$w/runner/release/dist"; }
NAMEN() { local o a e; for o in darwin linux windows; do for a in amd64 arm64; do e=""; [ "$o" = windows ] && e=.exe; echo "pgwire-recorder_${1}_${o}_${a}$e"; done; done; echo SHA256SUMS; }
# anhaenge_der_zeile <Zeile des Aufrufs von gh> — die Wörter, die unter dem Verzeichnis des Baus liegen, sortiert, einer je Zeile.
anhaenge() { tr ' ' '\n' <<< "$1" | grep "^$(DIST)/" | LC_ALL=C sort; }

# Punkt 2
AUSSEN=(github-release)
f_p2_trockenlauf_ohne_aussenwirkung() {
  fix 0.3.0-dev; umg trocken
  bis eingaben binaries
  sei "gh nie gerufen" nie gh; sei "skopeo nie gerufen" nie skopeo; sei "docker nie gerufen" nie docker
  sei "make binaries gerufen" enthaelt "$(aufrufe make)" "make binaries VERSION=0.3.0-dev"
}
f_p2_aussenwirkung_ohne_push() {
  local s
  for s in "${AUSSEN[@]}"; do
    fix 0.1.0; umg lauf; bis eingaben binaries
    umg lauf RELEASE_EVENT=workflow_dispatch; lauf "$s"
    grund 2 "$s" ereignis; sei "gh nie gerufen ($s)" nie gh
  done
}
f_p2_ref_kein_tag() {
  fix 0.1.0; umg lauf RELEASE_REF=refs/heads/main; lauf eingaben; grund 2 eingaben ereignis
  fix 0.1.0; umg lauf; bis eingaben binaries
  umg lauf RELEASE_REF=refs/heads/main; lauf github-release; grund 2 github-release ereignis; sei "gh nie gerufen" nie gh
}
# Punkt 3
gueltig() { # <Tag>
  fix "${1#v}"; umg lauf; lauf eingaben; gruen
}
f_p3_tag_gueltig() { local t; for t in v1.2.3 v0.1.0-rc.1 v1.0.0-x.7.z.92 v1.2.3-0a v0.0.0 v10.20.30 v1.2.3-rc-1 v1.2.3-RC.1; do gueltig "$t"; done; }
f_p3_laenge_100() { gueltig "v1.2.3-$(printf 'a%.0s' $(seq 1 93))"; }
# abgelehnt <Kennwort> <Tag> <Inhalt von VERSION (printf %b)> — Lauf: eingaben endet mit Exit 2 und dem Kennwort; das Handbuch trägt die Zeile der Version des Tags.
abgelehnt() {
  fix "${2#v}"; printf '%b' "$3" > "$w/repo/VERSION"
  umg lauf "RELEASE_REF=refs/tags/$2"; lauf eingaben; grund 2 eingaben "$1"
}
f_p3_ohne_v() { abgelehnt ereignis 1.2.3 '1.2.3\n'; }
f_p3_zwei_teile() { abgelehnt tag-form v1.2 '1.2\n'; }
f_p3_vier_teile() { abgelehnt tag-form v1.2.3.4 '1.2.3.4\n'; }
f_p3_fuehrende_null() { abgelehnt tag-form v01.2.3 '01.2.3\n'; }
f_p3_fuehrende_null_minor() { abgelehnt tag-form v1.02.3 '1.02.3\n'; }
f_p3_fuehrende_null_patch() { abgelehnt tag-form v1.2.03 '1.2.03\n'; }
f_p3_vorabversion_fuehrende_null() { abgelehnt tag-form v1.2.3-01 '1.2.3-01\n'; gueltig v1.2.3-0a; }
f_p3_leere_vorabversion() { abgelehnt tag-form v1.2.3- '1.2.3-\n'; }
f_p3_leerer_bezeichner() { abgelehnt tag-form v1.2.3-a..b '1.2.3-a..b\n'; }
f_p3_ungueltiges_zeichen() { abgelehnt tag-form v1.2.3-a_b '1.2.3-a_b\n'; }
f_p3_build_metadaten() { abgelehnt tag-form v1.2.3+x '1.2.3+x\n'; }
f_p3_zu_lang() { local t="v1.2.3-$(printf 'a%.0s' $(seq 1 94))"; abgelehnt tag-form "$t" "${t#v}\\n"; }
f_p3_ereignis_unbekannt() { fix 0.1.0; umg lauf RELEASE_EVENT=pull_request; lauf eingaben; grund 2 eingaben ereignis; }
f_p3_ereignis_fehlt() { fix 0.1.0; umg lauf -RELEASE_EVENT; lauf eingaben; grund 2 eingaben ereignis; }
vdatei() { abgelehnt version-datei v1.2.3 "$1"; }
f_p3_version_datei_form_fehlt() { fix 1.2.3; rm "$w/repo/VERSION"; umg lauf; lauf eingaben; grund 2 eingaben version-datei; }
f_p3_version_datei_form_zwei_zeilen() { vdatei '1.2.3\n1.2.3\n'; }
f_p3_version_datei_form_v() { vdatei 'v1.2.3\n'; }
f_p3_version_datei_form_leerzeichen_vorn() { vdatei ' 1.2.3\n'; }
f_p3_version_datei_form_leerzeichen_hinten() { vdatei '1.2.3 \n'; }
f_p3_version_datei_form_crlf() { vdatei '1.2.3\r\n'; }
f_p3_version_datei_form_kein_lf() { vdatei '1.2.3'; }
f_p3_tag_ungleich_datei() { abgelehnt tag-ungleich v1.2.3 '1.2.4\n'; }
f_p3_tag_ungleich_gross() { abgelehnt tag-ungleich v1.2.3-RC.1 '1.2.3-rc.1\n'; }
f_p3_tag_ungleich_tag_laenger() { abgelehnt tag-ungleich v1.2.3-rc.1 '1.2.3\n'; }
f_p3_tag_ungleich_tag_kuerzer() { abgelehnt tag-ungleich v1.2.3 '1.2.3-rc.1\n'; }
notes_fehlen() { fix 0.1.0 "$1"; umg lauf; lauf eingaben; grund 2 eingaben notes; }
f_p3_notes_fehlen() { notes_fehlen ohne; }
f_p3_notes_abschnitt() { notes_fehlen nur-anderer; }
f_p3_notes_ueberschrift_endet() { notes_fehlen ueberschrift; }
f_p3_notes_abschnitt_ende() { notes_fehlen spaeter; }
f_p3_nicht_im_standard_zweig() { fix 0.1.0; umg lauf "RELEASE_SHA=$s1"; lauf eingaben; grund 2 eingaben standard-zweig; }
f_p3_nicht_im_standard_zweig_vorwaerts() { fix 0.1.0; umg lauf "RELEASE_SHA=$c2"; lauf eingaben; grund 2 eingaben standard-zweig; }
f_p3_standard_zweig_aus_umgebung() { fix 0.1.0; umg lauf "RELEASE_SHA=$s1" RELEASE_DEFAULT_BRANCH=trunk; lauf eingaben; gruen; }
f_p3_trockenlauf_gueltig() {
  fix 0.1.0-dev; umg trocken; lauf eingaben; gruen
  sei "Meldung nennt den Tag aus VERSION" enthaelt "$out" "v0.1.0-dev"
}
f_p3_trockenlauf_im_standard_zweig() {
  fix 0.1.0; umg trocken "RELEASE_SHA=$c1"; lauf eingaben; gruen
  sei "Hinweis: liegt im Standard-Zweig" hat_hinweis "der Commit liegt im Standard-Zweig origin/main"
  sei "Hinweis ohne „nicht“" ohne "$err" "liegt nicht im"
}
f_p3_trockenlauf_nicht_im_standard_zweig() {
  fix 0.1.0; umg trocken "RELEASE_SHA=$s1"; lauf eingaben; gruen
  sei "Hinweis: liegt nicht im Standard-Zweig" hat_hinweis "der Commit liegt nicht im Standard-Zweig origin/main"
}
f_p3_notes_zelle_leer() { notes_fehlen leer; }
f_p3_trockenlauf_notes_fehlen() {
  fix 0.1.0 ohne; umg trocken; lauf eingaben; gruen
  sei "Hinweis: keine Zeile" hat_hinweis "keine Zeile"
}
f_p3_trockenlauf_notes_zelle_leer() { fix 0.1.0 leer; umg trocken; lauf eingaben; grund 2 eingaben notes; }
f_p3_trockenlauf_notes_vorhanden() {
  fix 0.1.0; umg trocken; lauf eingaben; gruen
  sei "Hinweis mit dem Text der Zeile" hat_hinweis "Text-0.1.0"
  sei "kein Hinweis „keine Zeile“" ohne "$err" "keine Zeile"
}
# Schlüsselwörter: Jede Meldung eines Eingabefehlers beginnt mit `release: <Schlüsselwort>: `.
f_p3_schluesselwort_ereignis() { fix 0.1.0; umg lauf RELEASE_EVENT=pull_request; lauf eingaben; sei "Schlüsselwort ereignis" beginnt ereignis; }
f_p3_schluesselwort_tag_form() { abgelehnt tag-form v1.2 '1.2\n'; sei "Schlüsselwort tag-form" beginnt tag-form; }
f_p3_schluesselwort_version_datei() { vdatei '1.2.3'; sei "Schlüsselwort version-datei" beginnt version-datei; }
f_p3_schluesselwort_tag_ungleich() { abgelehnt tag-ungleich v1.2.3 '1.2.4\n'; sei "Schlüsselwort tag-ungleich" beginnt tag-ungleich; }
f_p3_schluesselwort_notes() { notes_fehlen ohne; sei "Schlüsselwort notes" beginnt notes; }
f_p3_schluesselwort_standard_zweig() { fix 0.1.0; umg lauf "RELEASE_SHA=$s1"; lauf eingaben; sei "Schlüsselwort standard-zweig" beginnt standard-zweig; }
f_p3_trockenlauf_version_datei_form() { fix 0.1.0; printf '0.1.0' > "$w/repo/VERSION"; umg trocken; lauf eingaben; grund 2 eingaben version-datei; }
# Reihenfolge: je Paar benachbarter Gründe ein Fall mit beiden Verstößen; die Meldung nennt den früheren.
f_p3_reihenfolge_1_2() { fix 0.1.0; umg lauf RELEASE_REF=refs/tags/x; lauf eingaben; grund 2 eingaben ereignis; }
f_p3_reihenfolge_2_3() { fix 0.1.0; printf 'v1\n' > "$w/repo/VERSION"; umg lauf RELEASE_REF=refs/tags/v01.2.3; lauf eingaben; grund 2 eingaben tag-form; }
f_p3_reihenfolge_3_4() { fix 1.2.3; printf '1.2.4\r\n' > "$w/repo/VERSION"; umg lauf; lauf eingaben; grund 2 eingaben version-datei; }
f_p3_reihenfolge_4_5() { fix 1.2.3 ohne; printf '1.2.4\n' > "$w/repo/VERSION"; umg lauf; lauf eingaben; grund 2 eingaben tag-ungleich; }
f_p3_reihenfolge_5_6() { fix 0.1.0 ohne; umg lauf "RELEASE_SHA=$s1"; lauf eingaben; grund 2 eingaben notes; }
# Punkt 4
f_p4_aufruf() {
  local m d
  fix 0.2.0; umg lauf; bis eingaben binaries
  m="$(aufrufe make)"
  sei "Aufruf: make binaries" enthaelt "$m" "make binaries "
  sei "Aufruf: VERSION=0.2.0" enthaelt "$m" " VERSION=0.2.0 "
  d="$(grep -o ' DIST=[^ ]*' <<< "$m" | head -n 1)"; d="${d# DIST=}"
  sei "Aufruf: DIST außerhalb des Arbeitsbaums" test -n "$d" -a "${d#"$w/repo/"}" = "$d"
  sei "Aufruf: DIST ist neu" enthaelt "$m" "dist-vorher: nein"
  fix 0.3.0-dev; umg trocken; bis eingaben binaries
  sei "Trockenlauf: VERSION=0.3.0-dev" enthaelt "$(aufrufe make)" " VERSION=0.3.0-dev "
}
f_p4_summe_falsch() { fix 0.1.0; umg lauf; bis eingaben; lauf binaries GP_MAKE=summe-falsch; sei "Exit 1" test "$code" -eq 1; sei "ROT" enthaelt "$err" "release: binaries: ROT"; sei "keine Marke" test ! -e "$w/runner/release/binaries.ok"; }
f_p4_datei_fehlt() { fix 0.1.0; umg lauf; bis eingaben; lauf binaries GP_MAKE=datei-fehlt; sei "Exit 1" test "$code" -eq 1; sei "ROT" enthaelt "$err" "release: binaries: ROT"; }
f_p4_datei_zu_viel() { fix 0.1.0; umg lauf; bis eingaben; lauf binaries GP_MAKE=datei-zu-viel; sei "Exit 1" test "$code" -eq 1; sei "ROT" enthaelt "$err" "release: binaries: ROT"; }
f_p4_make_scheitert() { fix 0.1.0; umg lauf; bis eingaben; lauf binaries GP_MAKE=scheitert; sei "Exit 1" test "$code" -eq 1; sei "ROT" enthaelt "$err" "release: binaries: ROT"; }
# Punkt 9
ALLE7="pgwire-recorder_0.1.0_darwin_amd64 pgwire-recorder_0.1.0_darwin_arm64 pgwire-recorder_0.1.0_linux_amd64 pgwire-recorder_0.1.0_linux_arm64 pgwire-recorder_0.1.0_windows_amd64.exe pgwire-recorder_0.1.0_windows_arm64.exe SHA256SUMS"
# release_neu [Version [Modus]] — Lauf bis github-release ohne vorhandenen Release; setzt c (die Zeile des Aufrufs `release create`).
release_neu() { fix "${1:-0.1.0}" "${2:-normal}"; umg lauf; bis eingaben binaries; lauf github-release GP_GH=fehlt; gruen; c="$(grep '^gh release create' "$w/log/gh" 2> /dev/null || true)"; }
f_p9_anhaenge() {
  local erwartet
  release_neu
  erwartet="$(for n in $ALLE7; do echo "$(DIST)/$n"; done | LC_ALL=C sort)"
  sei "genau die sieben Anhänge" test "$(anhaenge "$c")" = "$erwartet"
  sei "sonst nichts: Wörter ohne Flags" test "$(tr ' ' '\n' <<< "$c" | grep -v -e "^$(DIST)/" -e '^--' -e "^pt9912/" -e "^$w/tmp/" -e '^v0.1.0$' | tr '\n' ' ')" = "gh release create "
}
f_p9_aufruf() {
  release_neu
  sei "release create" test -n "$c"
  sei "--verify-tag" enthaelt "$c" " --verify-tag "
  sei "--title v0.1.0" enthaelt "$c" " --title v0.1.0 "
  sei "--notes-file" enthaelt "$c" " --notes-file "
  sei "kein Entwurf" ohne "$c" "--draft"
}
f_p9_stabil_ohne_prerelease() { release_neu 0.1.0; sei "kein --prerelease" ohne "$c" "--prerelease"; }
f_p9_vorabversion_prerelease() { release_neu 0.1.0-rc.1; sei "--prerelease" enthaelt "$c" " --prerelease "; sei "genau einmal" test "$(grep -o -e '--prerelease' <<< "$c" | wc -l)" -eq 1; }
f_p9_notes_aus_zelle() { release_neu 0.1.0 normal; sei "Notes aus der dritten Zelle" test "$(cat "$w/log/notes")" = "Text-0.1.0"; }
f_p9_notes_zeile_der_version() { release_neu 0.1.0 zeile; sei "Notes der Zeile der Version" test "$(cat "$w/log/notes")" = "Text-0.1.0"; }
f_p9_notes_aus_abschnitt() { release_neu 0.1.0 abschnitt; sei "Notes aus der Änderungshistorie" test "$(cat "$w/log/notes")" = "Text-0.1.0"; }
vorhanden() { # <Assets> <Summe> — Release ist da; setzt code/out/err des Schritts github-release
  fix 0.1.0; umg lauf; bis eingaben binaries
  lauf github-release GP_GH=da "GP_GH_ASSETS=$1" "GP_GH_SUMS=$2"
}
f_p9_release_existiert_gleich() {
  local u
  vorhanden "pgwire-recorder_0.1.0_linux_amd64 SHA256SUMS" gleich; gruen
  sei "kein create" ohne "$(aufrufe gh)" "release create"
  u="$(grep '^gh release upload' "$w/log/gh" 2> /dev/null || true)"
  sei "upload der fehlenden fünf" test "$(anhaenge "$u")" = "$(for n in $ALLE7; do [ "$n" = pgwire-recorder_0.1.0_linux_amd64 ] || [ "$n" = SHA256SUMS ] || echo "$(DIST)/$n"; done | LC_ALL=C sort)"
  sei "kein --clobber" ohne "$(aufrufe gh)" "--clobber"
}
f_p9_release_existiert_abweichend() {
  vorhanden "pgwire-recorder_0.1.0_linux_amd64 SHA256SUMS" abweichend
  sei "Exit 1" test "$code" -eq 1; sei "ROT" enthaelt "$err" "release: github-release: ROT"
  sei "kein upload" ohne "$(aufrufe gh)" "release upload"; sei "kein create" ohne "$(aufrufe gh)" "release create"; sei "kein --clobber" ohne "$(aufrufe gh)" "--clobber"
}
f_p9_view_anderer_fehler() {
  fix 0.1.0; umg lauf; bis eingaben binaries
  lauf github-release GP_GH=fehler
  sei "Exit 1" test "$code" -eq 1; sei "kein create" ohne "$(aufrufe gh)" "release create"
}
# jeder Aufruf von gh trägt --repo
repo_in_jedem_aufruf() {
  local z n=0 alle=1
  while IFS= read -r z; do n=$((n + 1)); if [[ $z != *" --repo pt9912/pgwire-recorder"* ]]; then alle=0; fi; done < "$w/log/gh"
  sei "gh gerufen" test "$n" -ge 2; sei "jeder Aufruf von gh mit --repo" test "$alle" -eq 1
}
f_p9_repo_argument() {
  release_neu; repo_in_jedem_aufruf
  vorhanden "pgwire-recorder_0.1.0_linux_amd64 SHA256SUMS" gleich; repo_in_jedem_aufruf
  sei "view, download, upload gerufen" enthaelt "$(aufrufe gh)" "gh release upload"
}
f_p9_gh_token_fehlt() {
  fix 0.1.0; umg lauf; bis eingaben binaries
  umg lauf -GH_TOKEN; lauf github-release; grund 2 github-release token; sei "gh nie gerufen" nie gh
}
# Punkt 10
kein_token() { sei "kein Token in stdout und stderr" ohne "$out$err" "$TOKEN"; }
f_p10_nicht_im_log() {
  fix 0.1.0; umg lauf
  lauf eingaben; kein_token; lauf binaries; kein_token; lauf github-release GP_GH=fehlt; kein_token
  fix 0.1.0; umg lauf; bis eingaben binaries
  lauf github-release GP_GH=da "GP_GH_ASSETS=SHA256SUMS" GP_GH_SUMS=abweichend; sei "rot" test "$code" -eq 1; kein_token
  lauf github-release GP_GH=fehler; sei "rot" test "$code" -eq 1; kein_token
  fix 0.1.0; umg lauf; lauf eingaben; lauf binaries GP_MAKE=summe-falsch; sei "rot" test "$code" -eq 1; kein_token
  fix 0.1.0; umg lauf RELEASE_REF=refs/tags/v01.2.3; lauf eingaben; sei "rot" test "$code" -eq 2; kein_token
}
# Punkt 13: Marken zwischen den Schritten
f_p13_binaries_ohne_eingaben() { fix 0.1.0; umg lauf; lauf binaries; sei "Exit 2" test "$code" -eq 2; sei "make nie gerufen" nie make; }
f_p13_release_ohne_binaries() { fix 0.1.0; umg lauf; bis eingaben; lauf github-release; sei "Exit 2" test "$code" -eq 2; sei "gh nie gerufen" nie gh; }
# Punkt 11
f_p11_wiederanlauf_release_teilweise() {
  local u
  vorhanden "pgwire-recorder_0.1.0_linux_amd64 pgwire-recorder_0.1.0_darwin_arm64" gleich; gruen
  sei "kein create" ohne "$(aufrufe gh)" "release create"
  u="$(grep '^gh release upload' "$w/log/gh" 2> /dev/null || true)"
  sei "upload der fehlenden fünf, SHA256SUMS dabei" test "$(anhaenge "$u")" = "$(for n in $ALLE7; do [ "$n" = pgwire-recorder_0.1.0_linux_amd64 ] || [ "$n" = pgwire-recorder_0.1.0_darwin_arm64 ] || echo "$(DIST)/$n"; done | LC_ALL=C sort)"
  sei "kein download ohne SHA256SUMS am Release" ohne "$(aufrufe gh)" "release download"
  sei "kein --clobber" ohne "$(aufrufe gh)" "--clobber"
}
f_p11_wiederanlauf_release_vollstaendig() {
  vorhanden "$ALLE7" gleich; gruen
  sei "kein create" ohne "$(aufrufe gh)" "release create"; sei "kein upload" ohne "$(aufrufe gh)" "release upload"
}
f_p11_reihenfolge_gruende() {
  fix 0.1.0; umg lauf RELEASE_REF=refs/tags/v01.2.3
  lauf eingaben; sei "eingaben rot" test "$code" -eq 2
  lauf binaries; sei "binaries: Exit 2" test "$code" -eq 2; sei "make nie gerufen" nie make
  lauf github-release; sei "github-release: Exit 2" test "$code" -eq 2; sei "gh nie gerufen" nie gh
  fix 0.1.0; umg lauf; bis eingaben; lauf binaries GP_MAKE=summe-falsch; sei "binaries rot" test "$code" -eq 1
  lauf github-release; sei "github-release: Exit 2" test "$code" -eq 2; sei "gh nie gerufen" nie gh
}
f_p11_alte_marke() {
  fix 0.1.0; umg lauf; bis eingaben binaries
  umg lauf RELEASE_REF=refs/tags/v01.2.3; lauf eingaben; sei "eingaben rot" test "$code" -eq 2
  umg lauf; lauf binaries; sei "binaries: Exit 2 trotz alter Marke" test "$code" -eq 2
  fix 0.1.0; umg lauf; bis eingaben binaries
  lauf binaries GP_MAKE=summe-falsch; sei "binaries rot" test "$code" -eq 1
  lauf github-release; sei "github-release: Exit 2 trotz alter Marke" test "$code" -eq 2; sei "gh nie gerufen" nie gh
}

# --- Punkt 13: Abbruch ------------------------------------------------------------------------------------------------------------------------------
warte_auf() { local i; for i in $(seq 1 300); do if "$@"; then return 0; fi; sleep 0.1; done; return 1; }
fall_laeuft() { [ "$(find "$t" -mindepth 3 -maxdepth 3 -name 't.*' 2> /dev/null | wc -l)" -ge 1 ]; }
reste_da() { [ -n "$(ls -A "$t" 2> /dev/null)" ] || pgrep -f -- "$t" > /dev/null 2>&1; }
# warte_auf_reste — bis zu 8 s auf ein Verzeichnis oder einen Prozess, den ein nicht beendeter Nachkomme hinterlässt.
warte_auf_reste() { local i; for i in $(seq 1 80); do if reste_da; then return 0; fi; sleep 0.1; done; return 0; }
# abbruchfall <Signal> <Exit> — fährt $skript13 mit eigenem TMPDIR in eigener Sitzung, sendet das Signal, sobald dort ein Fall läuft, und prüft danach Exit, leeren TMPDIR
# und Prozesse mit dem Temp-Pfad in der Befehlszeile.
abbruchfall() {
  local sig="$1" soll="$2" pid=0 rc=0 bg ok=1 reste
  mkdir -p "$run"
  t="$(mktemp -d "$run/tmp.XXXXXX")" || exit 1
  ( trap - INT TERM; exec env --default-signal=INT,TERM -u GP_MUT -u GP_ENTDECKE -u GP_FALL TMPDIR="$t" GP_PARALLEL=2 GP_INNEN=1 setsid bash -c 'echo $$ > "$1"; exec bash "$2"' _ "$t.pid" "$skript13" > "$t.out" 2> "$t.err" ) &
  bg=$!
  warte_auf fall_laeuft || ok=0
  pid="$(cat "$t.pid" 2> /dev/null || echo 0)"
  if [ "$pid" -gt 1 ]; then kill -s "$sig" "$pid" 2> /dev/null; else ok=0; fi
  wait "$bg" || rc=$?
  warte_auf_reste
  [ "$rc" -eq "$soll" ] || ok=0
  reste="$(ls -A "$t" 2> /dev/null)"
  ! reste_da || ok=0
  if [ "$ok" -eq 0 ]; then
    printf 'erwartet: Exit %s, TMPDIR leer, kein Prozess mit dem Pfad\ngelesen: Exit %s, TMPDIR: %s\nstderr:\n%s\n' "$soll" "$rc" "$reste" "$(cat "$t.err" 2> /dev/null)" >> "$run/log.$fallname"
    bad=1
  fi
  if [ "$pid" -gt 1 ]; then kill -s KILL -- "-$pid" 2> /dev/null || true; fi
}
# Normales Ende der Gegenprobe: eine Kopie mit einem einzigen schnellen Lauf endet mit Exit 0 und lässt nichts in ihrem TMPDIR.
f_p13_gegenprobe_ende() {
  local rc=0 reste
  mkdir -p "$run"
  t="$(mktemp -d "$run/tmp.XXXXXX")" || exit 1
  env -u GP_FALL -u GP_ENTDECKE TMPDIR="$t" GP_INNEN=1 GP_MUT='^fall-ohne-funktion$' bash "$skript13" > "$t.out" 2> "$t.err" || rc=$?
  reste="$(ls -A "$t")"
  if [ "$rc" -ne 0 ] || [ -n "$reste" ]; then
    bad=1; printf 'erwartet: Exit 0, TMPDIR leer\ngelesen: Exit %s, TMPDIR: %s\n' "$rc" "$reste" >> "$run/log.$fallname"
  fi
}
f_p13_gegenprobe_sigint() { abbruchfall INT 130; }
f_p13_gegenprobe_sigterm() { abbruchfall TERM 143; }
# release.sh: Ende, SIGINT und SIGTERM in `binaries` mit einem make als Attrappe, das Nachkommen startet und in das Temp-Verzeichnis des Skripts schreibt.
f_p13_ende() {
  fix 0.1.0; umg lauf; bis eingaben
  sei "TMPDIR leer nach dem Ende" test -z "$(ls -A "$w/tmp")"
}
skript_abbruch() { # <Signal> <Exit>
  local sig="$1" soll="$2" pid=0 rc=0 bg ok=1 reste
  fix 0.1.0; umg lauf; bis eingaben
  rm -rf "${w:?}/tmp" && mkdir "$w/tmp"
  umg_befehl
  ( trap - INT TERM; cd "$w/repo" && exec "${ENVV[@]}" GP_MAKE=haengt setsid bash -c 'echo $$ > "$1"; shift; exec bash "$@"' _ "$w/pid" "$releasesh" binaries > "$w/o" 2> "$w/e" ) &
  bg=$!
  warte_auf test -e "$w/log/bereit" || ok=0
  pid="$(cat "$w/pid" 2> /dev/null || echo 0)"
  if [ "$pid" -gt 1 ]; then kill -s "$sig" "$pid" 2> /dev/null; else ok=0; fi
  wait "$bg" || rc=$?
  t="$w/tmp"
  warte_auf_reste
  [ "$rc" -eq "$soll" ] || ok=0
  reste="$(ls -A "$t" 2> /dev/null)"
  ! reste_da || ok=0
  if [ "$ok" -eq 0 ]; then
    printf 'erwartet: Exit %s, TMPDIR leer, kein Prozess mit dem Pfad\ngelesen: Exit %s, TMPDIR: %s\nstderr:\n%s\n' "$soll" "$rc" "$reste" "$(cat "$w/e" 2> /dev/null)" >> "$run/log.$fallname"
    bad=1
  fi
  if [ "$pid" -gt 1 ]; then kill -s KILL -- "-$pid" 2> /dev/null || true; fi
}
f_p13_sigint() { skript_abbruch INT 130; }
f_p13_sigterm() { skript_abbruch TERM 143; }

# --- Läufe ------------------------------------------------------------------------------------------------------------------------------------------
FAELLE_P1=(p1-gleich p1-text-gleich p1-reihenfolge-eingaben-binaries p1-reihenfolge-binaries-release p1-leerzeile p1-kommentarzeile p1-kommentarzeile-eingerueckt
  p1-pin-kurz p1-pin-gross p1-pin-ohne-tag p1-pin-tag-leer p1-pin-andere-action p1-bedingung-github-release)
for ((n = 1; n <= N; n++)); do FAELLE_P1+=("p1-fehlt-$n" "p1-geaendert-$n"); done
for k in "${ZUSATZ_NAMEN[@]}"; do FAELLE_P1+=("p1-zusatz-$k"); done
FAELLE_SH=(p2-trockenlauf-ohne-aussenwirkung p2-aussenwirkung-ohne-push p2-ref-kein-tag
  p3-tag-gueltig p3-laenge-100 p3-ohne-v p3-zwei-teile p3-vier-teile p3-fuehrende-null p3-fuehrende-null-minor p3-fuehrende-null-patch p3-vorabversion-fuehrende-null
  p3-leere-vorabversion p3-leerer-bezeichner p3-ungueltiges-zeichen p3-build-metadaten p3-zu-lang p3-ereignis-unbekannt p3-ereignis-fehlt
  p3-version-datei-form-fehlt p3-version-datei-form-zwei-zeilen p3-version-datei-form-v p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-leerzeichen-hinten
  p3-version-datei-form-crlf p3-version-datei-form-kein-lf p3-tag-ungleich-datei p3-tag-ungleich-gross p3-tag-ungleich-tag-laenger p3-tag-ungleich-tag-kuerzer
  p3-notes-fehlen p3-notes-zelle-leer p3-notes-abschnitt p3-notes-ueberschrift-endet p3-notes-abschnitt-ende p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts
  p3-standard-zweig-aus-umgebung p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-zelle-leer p3-trockenlauf-notes-vorhanden
  p3-trockenlauf-version-datei-form p3-schluesselwort-ereignis p3-schluesselwort-tag-form p3-schluesselwort-version-datei p3-schluesselwort-tag-ungleich p3-schluesselwort-notes p3-schluesselwort-standard-zweig p3-reihenfolge-1-2 p3-reihenfolge-2-3 p3-reihenfolge-3-4 p3-reihenfolge-4-5 p3-reihenfolge-5-6
  p4-aufruf p4-summe-falsch p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert
  p9-anhaenge p9-aufruf p9-stabil-ohne-prerelease p9-vorabversion-prerelease p9-notes-aus-zelle p9-notes-zeile-der-version p9-notes-aus-abschnitt
  p9-release-existiert-gleich p9-release-existiert-abweichend p9-view-anderer-fehler p9-repo-argument p9-gh-token-fehlt p13-binaries-ohne-eingaben p13-release-ohne-binaries p10-nicht-im-log
  p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p11-reihenfolge-gruende p11-alte-marke)
FAELLE_P13=(p13-ende p13-sigint p13-sigterm)
FAELLE_GP13=(p13-gegenprobe-sigint p13-gegenprobe-sigterm p13-gegenprobe-ende)

# faelle <gruppe> — schreibt die Namen der roten Fälle der Gruppe (p1, sh, p13, gp13, ok = alle), einen je Zeile. FALL_ZUSATZ nennt einen weiteren Namen (Probe: ein Name ohne
# Funktion ist ein roter Fall).
faelle() {
  local n
  case "$1" in
    p1) set -- "${FAELLE_P1[@]}" ;;
    sh) set -- "${FAELLE_SH[@]}" ;;
    p13) set -- "${FAELLE_P13[@]}" ;;
    gp13) if [ -n "${GP_INNEN-}" ]; then set --; else set -- "${FAELLE_GP13[@]}"; fi ;;
    ok) set -- "${FAELLE_P1[@]}" "${FAELLE_SH[@]}" "${FAELLE_P13[@]}"; [ -n "${GP_INNEN-}" ] || set -- "$@" "${FAELLE_GP13[@]}" ;;
  esac
  for n in "$@" ${FALL_ZUSATZ-}; do
    fallname="$n"
    mkdir -p "$run"
    declare -F "f_${n//-/_}" > /dev/null || { echo "$n"; continue; }
    ( bad=0; "f_${n//-/_}"; exit "$bad" ) > /dev/null 2>> "$run/stderr" || echo "$n"
  done
}
protokoll_zeigen() { [ -z "${GP_FALL-}" ] || [ ! -f "$run/log.$GP_FALL" ] || cat "$run/log.$GP_FALL" >&2; }

maxjobs="${GP_PARALLEL:-4}"
wartejob() { while [ "$(jobs -rp | wc -l)" -ge "$maxjobs" ]; do wait -n || true; done; }
# lauf_welt — legt die Welt des Laufs an (für die Gruppen sh, p13 und ok).
lauf_welt() { mkdir -p "$run/welt"; welt "$run/welt"; }
# bewerte <name> <rote Fälle> <Gruppe> — fährt die Gruppe und vergleicht die roten Fälle mit den erwarteten.
bewerte() {
  local ist erwartet
  ist="$(faelle "$3" | LC_ALL=C sort -u | tr '\n' ' ')"
  erwartet="$([ -z "$2" ] || printf '%s\n' $2 | LC_ALL=C sort -u | tr '\n' ' ')"
  if [ -n "${GP_ENTDECKE-}" ]; then echo "$1 => ${ist:-nichts}" >&2; fi
  echo "$1" >> "$arbeit/gefahren"
  [ "$ist" = "$erwartet" ] || melde "Mutation '$1': erwartet rot: ${2:-nichts}; rot sind: ${ist:-nichts}"
}
# fn_ersetze <alt> <neu> — Textersatz in pruefe_release (Quelltext aus declare -f), wirkt im aktuellen Lauf.
fn_ersetze() {
  local src alt="$1" rest
  src="$(declare -f pruefe_release)"
  if [[ $src != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
  rest="${src#*"$alt"}"
  if [[ $rest == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
  eval "${src/"$alt"/"$2"}"
}
# mut1 <name> <rote Fälle> <Code> — Mutation von pruefe_release: der Code läuft im Lauf vor den Fällen (fn_ersetze, Zuweisung an PIN_RE, SOLL).
mut1() {
  name="$1"
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  ( run="$arbeit/$name"; mkdir -p "$run"; eval "$3"; bewerte "$name" "$2" p1 ) &
}
# mutr <name> <rote Fälle> <Gruppe> <alt> <neu> [<alt> <neu>…] — Mutation einer Kopie von release.sh; rot sind genau die benannten Fälle der Gruppe. Beginnt <alt> mit `Z:`, ist es
# eine ganze Zeile (mit Zeilenende), die entfällt; beginnt es mit `L:`, entfällt die Zeile, die den folgenden Text (eine Marke) enthält.
mutr() {
  local name="$1" rote="$2" gruppe="$3" inhalt neu alt rest
  shift 3
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  (
    run="$arbeit/$name"; mkdir -p "$run"; releasesh="$run/release.sh"
    inhalt="$(cat "$repo/tools/harness/release.sh"; printf x)"; inhalt="${inhalt%x}"; neu="$inhalt"
    while [ "$#" -ge 2 ]; do
      alt="$1"; ersatz="$2"
      if [[ $alt == Z:* ]]; then alt="${alt#Z:}"$'\n'; ersatz=""; fi
      if [[ $alt == L:* ]]; then
        alt="${alt#L:}"
        if [ "$(grep -cF -- "$alt" <<< "$neu")" -ne 1 ]; then melde "Mutation '$name': Marke nicht genau einmal gefunden: $alt"; exit 0; fi
        if [ -z "$ersatz" ]; then
          neu="$(grep -vF -- "$alt" <<< "$neu"; printf x)"; neu="${neu%x}"
        else
          neu="$(while IFS= read -r z; do if [[ $z == *"$alt"* ]]; then printf '%s\n' "$ersatz"; else printf '%s\n' "$z"; fi; done <<< "$neu"; printf x)"; neu="${neu%x}"
        fi
        shift 2
        continue
      fi
      if [[ $neu != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
      rest="${neu#*"$alt"}"
      if [[ $rest == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
      neu="${neu/"$alt"/"$ersatz"}"
      shift 2
    done
    if [[ $neu == "$inhalt" ]]; then melde "Mutation '$name': Ersatz ändert nichts"; exit 0; fi
    printf '%s' "$neu" > "$releasesh"
    lauf_welt
    bewerte "$name" "$rote" "$gruppe"
  ) &
}

if [[ ok =~ ${GP_MUT:-.} ]]; then
  (
    run="$arbeit/ok"; mkdir -p "$run"; releasesh="$repo/tools/harness/release.sh"; skript13="$selbst"
    lauf_welt
    ist="$(faelle ok | LC_ALL=C sort -u | tr '\n' ' ')"
    echo ok >> "$arbeit/gefahren"
    [ -z "$ist" ] || melde "unveränderte Kopie ist nicht grün: $ist"
    protokoll_zeigen
    if [ -n "${GP_FALL-}" ]; then sed 's/^/stderr: /' "$run/stderr" >&2; fi
  ) &
fi

# mut13g <name> <rote Fälle> <alt> <neu> — Textersatz im Skript (Kopie in einem Baum mit den Dateien, die es braucht); rot sind genau die benannten Fälle p13-gegenprobe-*.
mut13g() {
  local name="$1" rote="$2" alt="$3" inhalt neu m f
  [ -z "${GP_INNEN-}" ] || return 0
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  (
    run="$arbeit/$name"; m="$run/baum"; mkdir -p "$m/tools/harness" "$m/.github/workflows"
    for f in .github/workflows/release.yml tools/harness/release.sh; do cp "$repo/$f" "$m/$f"; done
    inhalt="$(cat "$selbst"; printf x)"; inhalt="${inhalt%x}"
    if [[ $inhalt != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
    if [[ ${inhalt#*"$alt"} == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
    neu="${inhalt/"$alt"/"$4"}"
    printf '%s' "$neu" > "$m/tools/harness/release-gegenprobe.sh"
    skript13="$m/tools/harness/release-gegenprobe.sh"
    bewerte "$name" "$rote" gp13
  ) &
}

# Schutz gegen Rekursion: Eine Kopie (GP_INNEN=1), deren GP_MUT auf m-p13-gegenprobe passt, fährt nichts und endet mit Exit 1.
if [[ rekursion-schutz =~ ${GP_MUT:-.} ]] && [ -z "${GP_INNEN-}" ]; then
  (
    run="$arbeit/rekursion-schutz"; mkdir -p "$run"
    rc=0; ausgabe="$(env -u GP_FALL GP_INNEN=1 GP_MUT=m-p13-gegenprobe TMPDIR="$run" timeout 60 bash "$selbst" 2>&1)" || rc=$?
    echo rekursion-schutz >> "$arbeit/gefahren"
    [[ $ausgabe == *"keine Mutation und keine Kopie gefahren"* ]] && [ "$rc" -eq 1 ] || melde "Rekursionsschutz: erwartet Exit 1 ohne Lauf; Exit $rc: $ausgabe"
  ) &
fi

# Probe der Fallliste: ein Name ohne Funktion ist rot.
if [[ fall-ohne-funktion =~ ${GP_MUT:-.} ]]; then
  (
    run="$arbeit/fall-ohne-funktion"; mkdir -p "$run"
    ist="$(FALL_ZUSATZ=p1-tippfehler faelle p1 | LC_ALL=C sort -u | tr '\n' ' ')"
    echo fall-ohne-funktion >> "$arbeit/gefahren"
    [ "$ist" = "p1-tippfehler " ] || melde "Fall ohne Funktion: erwartet rot: p1-tippfehler; rot sind: ${ist:-nichts}"
  ) &
fi

# --- Mutationen von pruefe_release (Punkt 1) --------------------------------------------------------------------------------------------------------
# Je Zeile n zwei Mutationen: die Zeile ist optional (fehlt sie, rückt die Prüfung nicht weiter) und ihr Inhalt wird nicht verglichen.
for ((n = 1; n <= N; n++)); do
  i=$((n - 1))
  if [ "$i" -eq "$PIN" ]; then igel='[[ ${ist[j]-} =~ $PIN_RE ]]'; else igel='[ "${ist[j]-}" = "${SOLL[i]}" ]'; fi
  rote="p1-fehlt-$n"
  if [ "$i" -eq "$BEDINGUNG" ]; then rote+=" p1-bedingung-github-release"; fi
  mut1 "m-p1-zeile-$n-optional" "$rote" "fn_ersetze 'if [ \"\$i\" -eq \"\$PIN\" ]; then' 'if [ \"\$i\" -eq $i ] && ! $igel; then continue; fi; if [ \"\$i\" -eq \"\$PIN\" ]; then'"
  rote="p1-geaendert-$n"
  if [ "$i" -eq "$PIN" ]; then rote+=" p1-pin-kurz p1-pin-gross p1-pin-ohne-tag p1-pin-tag-leer p1-pin-andere-action"; fi
  mut1 "m-p1-zeile-$n-nicht-verglichen" "$rote" "fn_ersetze 'if [ \"\$i\" -eq \"\$PIN\" ]; then' 'if [ \"\$i\" -eq $i ]; then j=\$((j + 1)); continue; fi; if [ \"\$i\" -eq \"\$PIN\" ]; then'"
done
for k in "${!ZUSATZ[@]}"; do
  wort="${ZUSATZ[k]%%|*}"
  mut1 "m-p1-zusatz-${ZUSATZ_NAMEN[k]}-erlaubt" "p1-zusatz-${ZUSATZ_NAMEN[k]}" "fn_ersetze 'ist+=(\"\$z\")' '[[ \$z == *\"$wort\"* ]] || ist+=(\"\$z\")'"
done
mut1 m-p1-zusatzzeilen-erlaubt "$(printf 'p1-zusatz-%s ' "${ZUSATZ_NAMEN[@]}")" "fn_ersetze '[ \"\${#ist[@]}\" -eq \"\$j\" ]' 'true'"
mut1 m-p1-pin-grossbuchstaben "p1-pin-gross" 'PIN_RE="${PIN_RE/0-9a-f/0-9a-fA-F}"'
mut1 m-p1-pin-39-ziffern "p1-pin-kurz" 'PIN_RE="${PIN_RE/"{40}"/"{39,40}"}"'
mut1 m-p1-pin-ohne-kommentar "p1-pin-ohne-tag" 'PIN_RE="${PIN_RE/ # \[^ \].*\$/( # [^ ].*)?\$}"'
mut1 m-p1-pin-kommentar-ohne-text "p1-pin-tag-leer" 'PIN_RE="${PIN_RE/"[^ ]"/}"'
mut1 m-p1-pin-andere-action "p1-pin-andere-action" 'PIN_RE="${PIN_RE/actions?checkout@/actions/[a-z-]+@}"'
mut1 m-p1-kommentarzeilen-nicht-ignoriert "p1-kommentarzeile p1-kommentarzeile-eingerueckt" "fn_ersetze '^\\ *(#|\$)' '^\\ *\$'"
mut1 m-p1-leerzeilen-nicht-ignoriert "p1-leerzeile" "fn_ersetze '^\\ *(#|\$)' '^\\ *#'"
mut1 m-p1-kommentar-irgendwo-in-der-zeile "p1-gleich p1-text-gleich p1-leerzeile p1-kommentarzeile p1-kommentarzeile-eingerueckt" "fn_ersetze '^\\ *(#|\$)' '^.*#|^\\ *\$'"
mut1 m-p1-kommentar-nur-am-zeilenanfang "p1-kommentarzeile-eingerueckt" "fn_ersetze '^\\ *(#|\$)' '^(#|\$)'"
# Reihenfolge: die Prüfung nimmt zusätzlich die Folge mit vertauschten Schritten an.
PERM_BE=("${SOLL[@]:0:S_E}" "${SOLL[@]:S_B:S_G-S_B}" "${SOLL[@]:S_E:S_B-S_E}" "${SOLL[@]:S_G}")
PERM_GB=("${SOLL[@]:0:S_B}" "${SOLL[@]:S_G}" "${SOLL[@]:S_B:S_G-S_B}")
tauschbar='src="$(declare -f pruefe_release)"; eval "${src/pruefe_release/pruefe_release_orig}"; pruefe_release() { pruefe_release_orig "$1" && return 0; local -a s=("${SOLL[@]}"); local r; SOLL=("${%s[@]}"); pruefe_release_orig "$1"; r=$?; SOLL=("${s[@]}"); return $r; }'
# shellcheck disable=SC2059
mut1 m-p1-reihenfolge-eingaben-binaries-tauschbar "p1-reihenfolge-eingaben-binaries" "$(printf "$tauschbar" PERM_BE)"
# shellcheck disable=SC2059
mut1 m-p1-reihenfolge-binaries-release-tauschbar "p1-reihenfolge-binaries-release" "$(printf "$tauschbar" PERM_GB)"

# --- Mutationen von release.sh (Punkte 2 bis 11) ----------------------------------------------------------------------------------------------------
# Punkt 2
mutr m-p2-trockenlauf-ruft-skopeo "p2-trockenlauf-ohne-aussenwirkung" sh '  : > "$zdir/binaries.ok"' $'  [ "$trockenlauf" -eq 0 ] || skopeo inspect docker://x\n  : > "$zdir/binaries.ok"'
mutr m-p2-github-release-ohne-ereignispruefung "p2-aussenwirkung-ohne-push p2-ref-kein-tag" sh $'  aussenwirkung_erlaubt\n  zustand\n' $'  zustand\n'
mutr m-p2-ereignis-nicht-geprueft "p2-aussenwirkung-ohne-push p3-ereignis-fehlt p3-ereignis-unbekannt p3-schluesselwort-ereignis" sh \
  'Z:  [ "${RELEASE_EVENT-}" = push ] || eingabe ereignis "Ereignis '"'"'${RELEASE_EVENT-}'"'"' ist kein Push eines Tags" # P2ereignis' ""
mutr m-p2-ref-nicht-geprueft "p2-ref-kein-tag p3-ohne-v p3-reihenfolge-1-2" sh \
  'Z:  [[ ${RELEASE_REF-} == refs/tags/v* ]] || eingabe ereignis "Ref '"'"'${RELEASE_REF-}'"'"' ist kein Tag v…" # P2ref' ""
mutr m-p2-ref-ohne-v-praefix "p3-ohne-v p3-reihenfolge-1-2" sh '[[ ${RELEASE_REF-} == refs/tags/v* ]]' '[[ ${RELEASE_REF-} == refs/tags/* ]]'
mutr m-p2-trockenlauf-wie-lauf "p2-trockenlauf-ohne-aussenwirkung p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-trockenlauf-version-datei-form p4-aufruf" sh \
  'Z:  if [ "${RELEASE_EVENT-}" = workflow_dispatch ]; then return 0; fi' ""
# Punkt 3
TAGFORM="p3-zwei-teile p3-vier-teile p3-fuehrende-null p3-fuehrende-null-minor p3-fuehrende-null-patch p3-vorabversion-fuehrende-null p3-leere-vorabversion p3-leerer-bezeichner p3-ungueltiges-zeichen p3-build-metadaten p3-zu-lang p3-reihenfolge-2-3 p3-schluesselwort-tag-form"
VDATEI="p3-version-datei-form-fehlt p3-version-datei-form-zwei-zeilen p3-version-datei-form-v p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-leerzeichen-hinten p3-version-datei-form-crlf p3-version-datei-form-kein-lf p3-schluesselwort-version-datei"
UNGLEICH="p3-tag-ungleich-datei p3-tag-ungleich-gross p3-tag-ungleich-tag-laenger p3-tag-ungleich-tag-kuerzer p3-schluesselwort-tag-ungleich"
NOTES="p3-notes-fehlen p3-notes-zelle-leer p3-notes-abschnitt p3-notes-ueberschrift-endet p3-notes-abschnitt-ende p3-schluesselwort-notes"
# GRUEN: die Fälle, die eine grüne Prüfung der Eingaben brauchen oder deren Prüfung hinter (4) liegt; sie werden rot, wenn (4) nichts Gleiches mehr vergleicht.
# BIN: die Fälle, die `binaries` grün fahren und an der Dateiliste des Baus hängen; sie werden rot, wenn der Aufruf von make nicht die Version oder ein Verzeichnis außerhalb des Baums nennt.
BIN="p9-gh-token-fehlt p9-repo-argument p10-nicht-im-log p11-alte-marke p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p4-aufruf p9-anhaenge p9-aufruf p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease"
GRUEN="p13-release-ohne-binaries p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-schluesselwort-standard-zweig p3-vorabversion-fuehrende-null p10-nicht-im-log p11-alte-marke p11-reihenfolge-gruende p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-reihenfolge-5-6 p3-standard-zweig-aus-umgebung p3-tag-gueltig p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p4-aufruf p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert p4-summe-falsch p9-anhaenge p9-aufruf p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease"
# Je Prüfung (1) bis (6) eine Mutation, die den Aufruf entfernt.
mutr m-p3-1-ereignis-entfernt "p2-ref-kein-tag p3-ereignis-fehlt p3-ereignis-unbekannt p3-ohne-v p3-reihenfolge-1-2 p3-schluesselwort-ereignis" sh $'  pruefe_ereignis\n' ''
mutr m-p3-2-tag-form-entfernt "$TAGFORM" sh $'  pruefe_tag_form\n' ''
mutr m-p3-3-version-datei-entfernt "p10-nicht-im-log p11-alte-marke p11-reihenfolge-gruende p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p13-release-ohne-binaries p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-3-4 p3-reihenfolge-5-6 p3-schluesselwort-notes p3-schluesselwort-standard-zweig p3-schluesselwort-version-datei p3-standard-zweig-aus-umgebung p3-tag-gueltig p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-trockenlauf-version-datei-form p3-version-datei-form-crlf p3-version-datei-form-fehlt p3-version-datei-form-kein-lf p3-version-datei-form-leerzeichen-hinten p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-v p3-version-datei-form-zwei-zeilen p3-vorabversion-fuehrende-null p4-aufruf p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert p4-summe-falsch p9-anhaenge p9-aufruf p9-gh-token-fehlt p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-repo-argument p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease" sh $'  pruefe_version_datei\n' ''
mutr m-p3-4-tag-gleich-entfernt "$UNGLEICH p3-reihenfolge-4-5" sh $'  pruefe_tag_gleich\n' ''
mutr m-p3-5-notes-entfernt "p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-5-6 p3-schluesselwort-notes p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer" sh $'  pruefe_notes\n' ''
mutr m-p3-6-standard-zweig-entfernt "p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-schluesselwort-standard-zweig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig" sh $'  pruefe_standard_zweig\n' ''
# Reihenfolge: je Paar benachbarter Prüfungen vertauscht.
mutr m-p3-reihenfolge-1-2-vertauscht "p3-reihenfolge-1-2 p3-ohne-v p2-ref-kein-tag" sh $'  pruefe_ereignis\n  pruefe_tag_form\n' $'  pruefe_tag_form\n  pruefe_ereignis\n'
mutr m-p3-reihenfolge-2-3-vertauscht "${TAGFORM/ p3-zu-lang/}" sh $'  pruefe_tag_form\n  pruefe_version_datei\n' $'  pruefe_version_datei\n  pruefe_tag_form\n'
mutr m-p3-reihenfolge-3-4-vertauscht "p10-nicht-im-log p11-alte-marke p11-reihenfolge-gruende p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p13-release-ohne-binaries p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-3-4 p3-reihenfolge-5-6 p3-schluesselwort-notes p3-schluesselwort-standard-zweig p3-schluesselwort-version-datei p3-standard-zweig-aus-umgebung p3-tag-gueltig p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-trockenlauf-version-datei-form p3-version-datei-form-crlf p3-version-datei-form-fehlt p3-version-datei-form-kein-lf p3-version-datei-form-leerzeichen-hinten p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-v p3-version-datei-form-zwei-zeilen p3-vorabversion-fuehrende-null p4-aufruf p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert p4-summe-falsch p9-anhaenge p9-aufruf p9-gh-token-fehlt p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-repo-argument p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease" sh $'  pruefe_version_datei\n  pruefe_tag_gleich\n' $'  pruefe_tag_gleich\n  pruefe_version_datei\n'
mutr m-p3-reihenfolge-4-5-vertauscht "$UNGLEICH p3-reihenfolge-4-5" sh $'  pruefe_tag_gleich\n  pruefe_notes\n' $'  pruefe_notes\n  pruefe_tag_gleich\n'
mutr m-p3-reihenfolge-5-6-vertauscht "p3-reihenfolge-5-6" sh $'  pruefe_notes\n  pruefe_standard_zweig\n' $'  pruefe_standard_zweig\n  pruefe_notes\n'
# Form des Tags, je Element.
mutr m-p3-form-entfernt "${TAGFORM/ p3-zu-lang/}" sh 'L:# P3form' ''
mutr m-p3-laenge-entfernt "p3-zu-lang" sh 'L:# P3laenge' ''
mutr m-p3-laenge-101 "p3-zu-lang" sh '"${#tag}" -le 100' '"${#tag}" -le 101'
mutr m-p3-laenge-99 "p3-laenge-100" sh '"${#tag}" -le 100' '"${#tag}" -le 99'
mutr m-p3-null-major "p3-fuehrende-null p3-reihenfolge-2-3" sh 'KERN="$NUM\\.$NUM\\.$NUM" # P3kern' 'KERN="[0-9]+\\.$NUM\\.$NUM" # P3kern'
mutr m-p3-null-minor "p3-fuehrende-null-minor" sh 'KERN="$NUM\\.$NUM\\.$NUM" # P3kern' 'KERN="$NUM\\.[0-9]+\\.$NUM" # P3kern'
mutr m-p3-null-patch "p3-fuehrende-null-patch" sh 'KERN="$NUM\\.$NUM\\.$NUM" # P3kern' 'KERN="$NUM\\.$NUM\\.[0-9]+" # P3kern'
mutr m-p3-kern-zwei-teile "p3-schluesselwort-tag-form p3-zwei-teile" sh 'KERN="$NUM\\.$NUM\\.$NUM" # P3kern' 'KERN="$NUM\\.$NUM(\\.$NUM)?" # P3kern'
mutr m-p3-kern-vier-teile "p3-vier-teile" sh 'KERN="$NUM\\.$NUM\\.$NUM" # P3kern' 'KERN="$NUM\\.$NUM\\.$NUM(\\.$NUM)?" # P3kern'
mutr m-p3-vor-leer-erlaubt "p3-leere-vorabversion" sh 'VOR="(-$BEZ(\\.$BEZ)*)?" # P3vor' 'VOR="(-($BEZ(\\.$BEZ)*)?)?" # P3vor'
mutr m-p3-vor-leerer-bezeichner "p3-leerer-bezeichner" sh 'VOR="(-$BEZ(\\.$BEZ)*)?" # P3vor' 'VOR="(-$BEZ(\\.$BEZ?)*)?" # P3vor'
mutr m-p3-vor-build-metadaten "p3-build-metadaten" sh 'VOR="(-$BEZ(\\.$BEZ)*)?" # P3vor' 'VOR="(-$BEZ(\\.$BEZ)*)?(\\+[0-9A-Za-z.-]+)?" # P3vor'
mutr m-p3-vor-ein-bezeichner "p3-tag-gueltig p3-tag-ungleich-gross p3-tag-ungleich-tag-laenger p3-tag-ungleich-tag-kuerzer p9-vorabversion-prerelease" sh 'VOR="(-$BEZ(\\.$BEZ)*)?" # P3vor' 'VOR="(-$BEZ)?" # P3vor'
mutr m-p3-bez-fuehrende-null "p3-vorabversion-fuehrende-null" sh "BEZ='(0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)' # P3bez" "BEZ='[0-9A-Za-z-]+' # P3bez"
mutr m-p3-bez-zeichen "p3-ungueltiges-zeichen" sh "BEZ='(0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)' # P3bez" "BEZ='(0|[1-9][0-9]*|[0-9A-Za-z_-]*[A-Za-z_-][0-9A-Za-z_-]*)' # P3bez"
mutr m-p3-bez-nur-zahlen "p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-tag-gueltig p3-tag-ungleich-gross p3-tag-ungleich-tag-kuerzer p3-tag-ungleich-tag-laenger p3-trockenlauf-gueltig p3-vorabversion-fuehrende-null p4-aufruf p9-vorabversion-prerelease" sh "BEZ='(0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)' # P3bez" "BEZ='(0|[1-9][0-9]*)' # P3bez"
mutr m-p3-bez-ohne-ziffern "p3-tag-gueltig p3-vorabversion-fuehrende-null" sh "BEZ='(0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)' # P3bez" "BEZ='(0|[1-9][0-9]*|[A-Za-z-]*[A-Za-z-][A-Za-z-]*)' # P3bez"
mutr m-p3-bez-ohne-grossbuchstaben "p3-tag-gueltig p3-tag-ungleich-gross" sh "BEZ='(0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)' # P3bez" "BEZ='(0|[1-9][0-9]*|[0-9a-z-]*[a-z-][0-9a-z-]*)' # P3bez"
mutr m-p3-bez-ohne-kleinbuchstaben "p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-tag-gueltig p3-tag-ungleich-gross p3-tag-ungleich-tag-kuerzer p3-tag-ungleich-tag-laenger p3-trockenlauf-gueltig p3-vorabversion-fuehrende-null p4-aufruf p9-vorabversion-prerelease" sh "BEZ='(0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)' # P3bez" "BEZ='(0|[1-9][0-9]*|[0-9A-Z-]*[A-Z-][0-9A-Z-]*)' # P3bez"
mutr m-p3-bez-ohne-bindestrich "p3-tag-gueltig" sh "BEZ='(0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)' # P3bez" "BEZ='(0|[1-9][0-9]*|[0-9A-Za-z]*[A-Za-z][0-9A-Za-z]*)' # P3bez"
# Der Präfix v der Form ist durch (1) schon gesichert (Ref refs/tags/v…; im Trockenlauf trägt der Tag das v der Konstruktion): kein Fall wird rot (äquivalent, §7).
mutr m-p3-form-v-optional "" sh '[[ $tag =~ ^v$SEMVER$ ]]' '[[ $tag =~ ^v?$SEMVER$ ]]'
# Form der Datei VERSION, je Element.
mutr m-p3-vform-entfernt "$VDATEI p3-reihenfolge-3-4 p3-trockenlauf-version-datei-form" sh 'L:# P3vform' ''
mutr m-p3-vdatei-kein-lf "p3-schluesselwort-version-datei p3-trockenlauf-version-datei-form p3-version-datei-form-kein-lf" sh $'VERSION_RE="^$SEMVER"$\'\\n\'"\\$" # P3vdatei' $'VERSION_RE="^$SEMVER"$\'\\n\'"?\\$" # P3vdatei'
mutr m-p3-vdatei-crlf "p3-version-datei-form-crlf p3-reihenfolge-3-4" sh $'VERSION_RE="^$SEMVER"$\'\\n\'"\\$" # P3vdatei' $'VERSION_RE="^$SEMVER"$\'\\r\'?$\'\\n\'"\\$" # P3vdatei'
mutr m-p3-vdatei-leerzeichen "p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-leerzeichen-hinten" sh $'VERSION_RE="^$SEMVER"$\'\\n\'"\\$" # P3vdatei' $'VERSION_RE="^ *$SEMVER *"$\'\\n\'"\\$" # P3vdatei'
mutr m-p3-vdatei-v "p3-version-datei-form-v" sh $'VERSION_RE="^$SEMVER"$\'\\n\'"\\$" # P3vdatei' $'VERSION_RE="^v?$SEMVER"$\'\\n\'"\\$" # P3vdatei'
mutr m-p3-vdatei-mehrere-zeilen "p3-version-datei-form-zwei-zeilen" sh $'VERSION_RE="^$SEMVER"$\'\\n\'"\\$" # P3vdatei' $'VERSION_RE="^$SEMVER("$\'\\n\'"$SEMVER)*"$\'\\n\'"\\$" # P3vdatei'
# Tag gleich v plus VERSION.
mutr m-p3-gleich-ohne-v "p10-nicht-im-log p11-alte-marke p11-reihenfolge-gruende p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p13-release-ohne-binaries p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-5-6 p3-schluesselwort-notes p3-schluesselwort-standard-zweig p3-standard-zweig-aus-umgebung p3-tag-gueltig p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-vorabversion-fuehrende-null p4-aufruf p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert p4-summe-falsch p9-anhaenge p9-aufruf p9-gh-token-fehlt p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-repo-argument p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease" sh '"$tag" = "v$version"' '"$tag" = "$version"'
mutr m-p3-gleich-gross "p3-tag-ungleich-gross" sh '"$tag" = "v$version"' '"${tag,,}" = "v${version,,}"'
mutr m-p3-gleich-praefix-tag-laenger "p3-tag-ungleich-tag-laenger" sh '[ "$tag" = "v$version" ]' '[[ $tag == "v$version"* ]]'
mutr m-p3-gleich-praefix-tag-kuerzer "p3-tag-ungleich-tag-kuerzer" sh '[ "$tag" = "v$version" ]' '[[ "v$version" == "$tag"* ]]'
# Zeile der Version in der Änderungshistorie.
mutr m-p3-notes-abschnitt-ignoriert "p3-notes-abschnitt p3-notes-ueberschrift-endet p3-notes-abschnitt-ende p9-notes-aus-abschnitt" sh 'ab && /^\|/ &&' '/^\|/ &&'
mutr m-p3-notes-ueberschrift-ende-ignoriert "p3-notes-ueberschrift-endet" sh 'Änderungshistorie[ \t]*$/ { ab = 1; next }' 'Änderungshistorie/ { ab = 1; next }'
mutr m-p3-notes-abschnitt-ende-ignoriert "p3-notes-abschnitt-ende" sh $'    /^#/ { ab = 0 }\n' ''
mutr m-p3-notes-praefix "p9-notes-zeile-der-version" sh 'trim($2) == v' 'index(trim($2), v) == 1'
mutr m-p3-notes-zelle-datum "p3-notes-zelle-leer p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version" sh 'print "Z" trim($4)' 'print "Z" trim($3)'
# Commit im Standard-Zweig.
mutr m-p3-anc-ohne-is-ancestor "p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-schluesselwort-standard-zweig p3-trockenlauf-nicht-im-standard-zweig" sh 'git merge-base --is-ancestor "${RELEASE_SHA-}" "$ref"' 'git merge-base "${RELEASE_SHA-}" "$ref"'
mutr m-p3-anc-umgekehrt "p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-schluesselwort-standard-zweig p3-trockenlauf-nicht-im-standard-zweig" sh 'git merge-base --is-ancestor "${RELEASE_SHA-}" "$ref"' 'git merge-base --is-ancestor "$ref" "${RELEASE_SHA-}"'
mutr m-p3-standard-zweig-main "p3-standard-zweig-aus-umgebung" sh 'local ref="origin/${RELEASE_DEFAULT_BRANCH-}"' 'local ref="origin/main"'
mutr m-p3-lauf-nur-hinweis "p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-schluesselwort-standard-zweig" sh '  elif [ "$ok" -eq 0 ]; then' '  elif false; then'
mutr m-p3-trockenlauf-fehler "p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig" sh $'  if [ "$trockenlauf" -eq 1 ]; then\n    if [ "$ok" -eq 1 ]' $'  if false; then\n    if [ "$ok" -eq 1 ]'
mutr m-p3-hinweis-immer-ja "p3-trockenlauf-nicht-im-standard-zweig" sh 'else hinweis "der Commit liegt nicht im Standard-Zweig $ref"; fi' 'else hinweis "der Commit liegt im Standard-Zweig $ref"; fi'
mutr m-p3-hinweis-fehlt "p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig" sh $'    if [ "$ok" -eq 1 ]; then hinweis "der Commit liegt im Standard-Zweig $ref"; else hinweis "der Commit liegt nicht im Standard-Zweig $ref"; fi' '    :'
mutr m-p3-trockenlauf-tag-aus-ref "p2-trockenlauf-ohne-aussenwirkung p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-trockenlauf-version-datei-form p4-aufruf" sh \
  $'tag="v$(head -n 1 VERSION 2> /dev/null)"' $'tag="${RELEASE_REF-}"; tag="${tag#refs/tags/}"'
# Punkt 4
mutr m-p4-version-fehlt "$BIN" sh ' "VERSION=${tag#v}" "DIST=$dist"' ' "DIST=$dist"'
mutr m-p4-dist-im-baum "$BIN" sh '"DIST=$dist"' '"DIST=dist"'
mutr m-p4-dist-vorher-angelegt "p4-aufruf" sh $'  dist="$zdir/dist"\n  fuehre make' $'  dist="$zdir/dist"\n  mkdir -p "$dist"\n  fuehre make'
mutr m-p4-summe-entfernt "p10-nicht-im-log p11-alte-marke p11-reihenfolge-gruende p4-summe-falsch" sh '(cd "$dist" && sha256sum -c SHA256SUMS)' '(cd "$dist" && true)'
mutr m-p4-liste-entfernt "p4-datei-fehlt p4-datei-zu-viel" sh 'L:# P4liste' ''
mutr m-p4-make-fehler-ignoriert "p4-make-scheitert" sh $'  [ "$rc" -eq 0 ] || rot "make binaries endet mit Exit $rc"\n' ''
# Punkt 9
mutr m-p9-achter-anhang "p9-anhaenge p9-release-existiert-gleich p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig" sh \
  'while IFS= read -r a; do anhaenge+=("$dist/$a"); done < <(erwartete_namen "$v") # P9anhaenge' $'while IFS= read -r a; do anhaenge+=("$dist/$a"); done < <(erwartete_namen "$v") # P9anhaenge\n  anhaenge+=("$dist/SHA256SUMS.sig")'
mutr m-p9-prerelease-immer "p9-stabil-ohne-prerelease" sh 'L:# P9pre' $'  prerelease=(--prerelease)'
mutr m-p9-prerelease-nie "p9-vorabversion-prerelease" sh 'L:# P9pre' ''
mutr m-p9-ohne-verify-tag "p9-aufruf" sh ' --verify-tag' ''
mutr m-p9-ohne-title "p9-aufruf" sh ' --title "$tag"' ''
mutr m-p9-ohne-notes-file "p9-aufruf p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version" sh ' --notes-file "$scratch/notes.md"' ''
mutr m-p9-entwurf "p9-aufruf" sh ' --verify-tag' ' --verify-tag --draft'
mutr m-p9-vergleich-entfernt "p10-nicht-im-log p9-release-existiert-abweichend" sh 'L:# P9vergleich' ''
mutr m-p9-ueberschreiben-statt-vergleichen "p10-nicht-im-log p9-release-existiert-abweichend" sh \
  '|| rot "SHA256SUMS am Release weicht von der lokalen Datei ab; der Weg ist ein höherer Tag"' $'|| fuehre gh release upload "$tag" --repo "$RELEASE_REPOSITORY" --clobber "$dist/SHA256SUMS"'
mutr m-p9-upload-clobber "p9-release-existiert-gleich p11-wiederanlauf-release-teilweise" sh \
  'fuehre gh release upload "$tag" --repo "$RELEASE_REPOSITORY" "${fehlend[@]}"' 'fuehre gh release upload "$tag" --repo "$RELEASE_REPOSITORY" --clobber "${fehlend[@]}"'
mutr m-p9-create-trotz-vorhandenem-release "p10-nicht-im-log p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p9-release-existiert-abweichend p9-release-existiert-gleich p9-repo-argument" sh \
  'if [ "$vorhanden" -eq 0 ]; then' 'if true; then'
mutr m-p9-upload-alle "p9-release-existiert-gleich p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig" sh \
  $'for a in "${anhaenge[@]}"; do grep -qxF "${a##*/}" "$scratch/assets" || fehlend+=("$a"); done' $'fehlend=("${anhaenge[@]}")'
mutr m-p9-view-anderer-fehler-wie-fehlend "p10-nicht-im-log p9-view-anderer-fehler" sh '  elif ! grep -q '"'"'release not found'"'"' "$scratch/gh.err"; then' '  elif false; then'
# Zeile der Änderungshistorie im Trockenlauf und leere dritte Zelle
mutr m-p3-zelle-leer-angenommen-lauf "p3-notes-zelle-leer" sh '[ -n "${z#Z}" ] || eingabe notes' '[ -n "${z#Z}" ] || [ "$trockenlauf" -eq 0 ] || eingabe notes'
mutr m-p3-zelle-leer-angenommen-trockenlauf "p3-trockenlauf-notes-zelle-leer" sh '[ -n "${z#Z}" ] || eingabe notes' '[ -n "${z#Z}" ] || [ "$trockenlauf" -eq 1 ] || eingabe notes'
mutr m-p3-zelle-leer-angenommen "p3-notes-zelle-leer p3-trockenlauf-notes-zelle-leer" sh 'L:# P3zelle' ''
mutr m-p3-trockenlauf-zeile-nicht-ausgewertet "p3-trockenlauf-notes-vorhanden" sh 'L:# P3zeiletext' ''
mutr m-p3-trockenlauf-fehlende-zeile-fehler "p3-trockenlauf-notes-fehlen" sh 'if [ "$trockenlauf" -eq 1 ]; then # P3zeiletrocken' 'if false; then # P3zeiletrocken'
# Im Lauf ist eine fehlende Zeile wie eine leere Zelle ein Fehler `notes`: äquivalent (§7)
mutr m-p3-lauf-fehlende-zeile-ohne-eigene-pruefung "" sh 'Z:    eingabe notes "die Änderungshistorie des Benutzerhandbuchs hat keine Zeile für $version" # P3notes' ""
# Schlüsselwörter
declare -A SW_ROT=([ereignis]="p2-aussenwirkung-ohne-push p2-ref-kein-tag p3-ereignis-fehlt p3-ereignis-unbekannt p3-ohne-v p3-reihenfolge-1-2 p3-schluesselwort-ereignis" [tag-form]="p3-build-metadaten p3-fuehrende-null p3-fuehrende-null-minor p3-fuehrende-null-patch p3-leere-vorabversion p3-leerer-bezeichner p3-reihenfolge-2-3 p3-schluesselwort-tag-form p3-ungueltiges-zeichen p3-vier-teile p3-vorabversion-fuehrende-null p3-zu-lang p3-zwei-teile" [version-datei]="p3-reihenfolge-3-4 p3-schluesselwort-version-datei p3-trockenlauf-version-datei-form p3-version-datei-form-crlf p3-version-datei-form-fehlt p3-version-datei-form-kein-lf p3-version-datei-form-leerzeichen-hinten p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-v p3-version-datei-form-zwei-zeilen" [tag-ungleich]="p3-reihenfolge-4-5 p3-schluesselwort-tag-ungleich p3-tag-ungleich-datei p3-tag-ungleich-gross p3-tag-ungleich-tag-kuerzer p3-tag-ungleich-tag-laenger" [notes]="p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-5-6 p3-schluesselwort-notes p3-trockenlauf-notes-zelle-leer" [standard-zweig]="p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-schluesselwort-standard-zweig")
for kw in ereignis tag-form version-datei tag-ungleich notes standard-zweig; do
  mutr "m-p3-schluesselwort-$kw-weicht-ab" "${SW_ROT[$kw]}" sh 'eingabe() { echo "release: $1: $2" >&2; exit 2; }' "eingabe() { local k=\"\$1\"; [ \"\$k\" != $kw ] || k=\"\${k}x\"; echo \"release: \$k: \$2\" >&2; exit 2; }"
done
# Punkt 9
mutr m-p9-repo-fehlt-view "p9-repo-argument" sh 'release view "$tag" --repo "$RELEASE_REPOSITORY" --json' 'release view "$tag" --json'
mutr m-p9-repo-fehlt-create "p9-repo-argument" sh 'release create "$tag" --repo "$RELEASE_REPOSITORY" --verify-tag' 'release create "$tag" --verify-tag'
mutr m-p9-repo-fehlt-download "p9-repo-argument" sh 'release download "$tag" --repo "$RELEASE_REPOSITORY" --pattern' 'release download "$tag" --pattern'
mutr m-p9-repo-fehlt-upload "p9-repo-argument" sh 'release upload "$tag" --repo "$RELEASE_REPOSITORY" "${fehlend[@]}"' 'release upload "$tag" "${fehlend[@]}"'
mutr m-p9-gh-ohne-token-pruefung "p9-gh-token-fehlt" sh 'L:# P9token' ''
# Punkt 10
mutr m-p10-token-bei-fehler "p10-nicht-im-log" sh '|| rot "SHA256SUMS am Release' '|| { echo "$GH_TOKEN" >&2; false; } || rot "SHA256SUMS am Release'
mutr m-p10-token-bei-fehler-binaries "p10-nicht-im-log" sh '|| { sed '"'"'s/^/release: binaries: /'"'"' "$scratch/summen" >&2;' '|| { echo "$GH_TOKEN" >&2; sed '"'"'s/^/release: binaries: /'"'"' "$scratch/summen" >&2;'
mutr m-p10-set-x "p10-nicht-im-log p3-schluesselwort-ereignis p3-schluesselwort-notes p3-schluesselwort-standard-zweig p3-schluesselwort-tag-form p3-schluesselwort-tag-ungleich p3-schluesselwort-version-datei" sh $'set -uo pipefail\nexport LC_ALL=C\n' $'set -uo pipefail\nset -x\nexport LC_ALL=C\n'
mutr m-p10-env-bei-fehler "p10-nicht-im-log" sh '|| rot "SHA256SUMS am Release' '|| { env >&2; false; } || rot "SHA256SUMS am Release'
# Punkt 11
mutr m-p11-binaries-ohne-eingaben-marke "p11-alte-marke p11-reihenfolge-gruende p13-binaries-ohne-eingaben" sh 'L:# P11eingaben' ''
mutr m-p11-release-ohne-binaries-marke "p11-alte-marke p11-reihenfolge-gruende p13-release-ohne-binaries" sh 'L:# P11binaries' ''
mutr m-p11-eingaben-loescht-marken-nicht "p11-alte-marke" sh $'  rm -f "$zdir/eingaben.ok" "$zdir/binaries.ok"\n' ''
mutr m-p11-binaries-loescht-marke-nicht "p11-alte-marke" sh $'  rm -f "$zdir/binaries.ok"\n  tag_ableiten\n  dist=' $'  tag_ableiten\n  dist='
# Punkt 13: Abbruch von release.sh
mutr m-p13-ende-falle-fehlt "p13-ende" p13 $'\ntrap raeume_auf EXIT\n' $'\n'
mutr m-p13-int-falle-fehlt "p13-sigint" p13 $'\ntrap abbruch_int INT\n' $'\n'
mutr m-p13-term-falle-fehlt "p13-sigterm" p13 $'\ntrap abbruch_term TERM\n' $'\n'
mutr m-p13-int-kinder-bleiben "p13-sigint" p13 $'\nabbruch_int() { trap - EXIT INT TERM; beende_kinder; ' $'\nabbruch_int() { trap - EXIT INT TERM; '
mutr m-p13-term-kinder-bleiben "p13-sigterm" p13 $'\nabbruch_term() { trap - EXIT INT TERM; beende_kinder; ' $'\nabbruch_term() { trap - EXIT INT TERM; '
mutr m-p13-int-verzeichnis-zu-frueh "p13-sigint" p13 $'\nabbruch_int() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 130; }' $'\nabbruch_int() { trap - EXIT INT TERM; raeume_auf; sleep 5; beende_kinder; exit 130; }'
mutr m-p13-term-verzeichnis-zu-frueh "p13-sigterm" p13 $'\nabbruch_term() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 143; }' $'\nabbruch_term() { trap - EXIT INT TERM; raeume_auf; sleep 5; beende_kinder; exit 143; }'
# Punkt 13: Abbruch der Gegenprobe selbst
mut13g m-p13-gegenprobe-ende-falle-fehlt "p13-gegenprobe-ende" $'\ntrap raeume_auf EXIT\n' $'\n'
mut13g m-p13-gegenprobe-int-falle-fehlt "p13-gegenprobe-sigint" $'\ntrap abbruch_int INT\n' $'\n'
mut13g m-p13-gegenprobe-term-falle-fehlt "p13-gegenprobe-sigterm" $'\ntrap abbruch_term TERM\n' $'\n'
mut13g m-p13-gegenprobe-int-kinder-bleiben "p13-gegenprobe-sigint" $'\nabbruch_int() { trap - EXIT INT TERM; beende_kinder; ' $'\nabbruch_int() { trap - EXIT INT TERM; '
mut13g m-p13-gegenprobe-term-kinder-bleiben "p13-gegenprobe-sigterm" $'\nabbruch_term() { trap - EXIT INT TERM; beende_kinder; ' $'\nabbruch_term() { trap - EXIT INT TERM; '
mut13g m-p13-gegenprobe-int-verzeichnis-zu-frueh "p13-gegenprobe-sigint" $'\nabbruch_int() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 130; }' $'\nabbruch_int() { trap - EXIT INT TERM; raeume_auf; sleep 5; beende_kinder; exit 130; }'
mut13g m-p13-gegenprobe-term-verzeichnis-zu-frueh "p13-gegenprobe-sigterm" $'\nabbruch_term() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 143; }' $'\nabbruch_term() { trap - EXIT INT TERM; raeume_auf; sleep 5; beende_kinder; exit 143; }'
wait
fehler=0
if [ -s "$arbeit/fehler" ]; then cat "$arbeit/fehler" >&2; fehler=1; fi
if [ "$(wc -l < "$arbeit/gefahren")" -eq 0 ]; then echo "release-gegenprobe: keine Mutation und keine Kopie gefahren (GP_MUT passt auf nichts)" >&2; fehler=1; fi
if [ "$fehler" -eq 0 ]; then echo "release-gegenprobe: grün ($(wc -l < "$arbeit/gefahren") Läufe)"; fi
exit "$fehler"
