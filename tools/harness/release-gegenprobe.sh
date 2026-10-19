#!/usr/bin/env bash
# release-gegenprobe — prüft die Zusagen von SPEC-067 (harness/vertraege/SPEC-067-release.md, Punkt 14): den Text von .github/workflows/release.yml
# (Fälle p1-*, geprüft von der Funktion pruefe_release, die den Text hält) und tools/harness/release.sh (Fälle p2-* bis p16-*): je Fall ein Temp-Baum mit einem
# Git-Repository (Commits, Zweig, refs/remotes/origin/main und origin/trunk), dem Auszug des Benutzerhandbuchs, einem PATH, in dem `make`, `gh`, `docker`, `skopeo` und
# `gp-layout` (schreibt ein OCI-Layout) Attrappen sind, die ihre Aufrufe aufzeichnen; `skopeo` führt Buch über zwei Registries in einem Verzeichnis, `tools/harness/container-smoke.sh`
# des Baums ist eine Attrappe; `git`, `sha256sum`, `awk`, `tar` und `cmp` sind die echten Werkzeuge. Ohne Docker, ohne Netz, ohne make. Ein Fall besteht,
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
# Braucht bash (mit `wait -n`), git, GNU sed, grep, awk, cmp, tar, pgrep, setsid, timeout, sha256sum, jq, date (GNU) und die Werkzeuge von coreutils; ohne Docker, ohne make und ohne Netz.
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

for t in git sed grep awk cmp tar mktemp sort tr wc cat cp chmod rm env find ls seq sleep mkdir dirname basename head pgrep setsid timeout sha256sum; do
  command -v "$t" > /dev/null 2>&1 || { echo "release-gegenprobe: $t fehlt" >&2; exit 2; }
done
for f in .github/workflows/release.yml tools/harness/release.sh; do
  [ -f "$repo/$f" ] || { echo "release-gegenprobe: $f fehlt im Arbeitsbaum" >&2; exit 2; }
done
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1
export GIT_AUTHOR_NAME=gp GIT_AUTHOR_EMAIL=gp@example.invalid GIT_COMMITTER_NAME=gp GIT_COMMITTER_EMAIL=gp@example.invalid

# --- Punkt 1: der Text -------------------------------------------------------------------------------------------------------------------------------
# SOLL_ROH ist der Text von Punkt 1b, Zeile für Zeile, mit der Herkunft (A oder B) vor dem ersten `|`; HERKUNFT nennt die Herkunft, die diese Fassung führt. Die Zeile
# mit <SHA> und <Tag> prüft PIN_RE. Die Gegenprobe trägt den Text wörtlich als eigene Kopie.
HERKUNFT=ABF
mapfile -t SOLL_ROH << 'EOT'
A|name: release
A|on:
A|  push:
A|    tags:
A|      - 'v*'
A|  workflow_dispatch:
A|permissions: {}
A|concurrency:
A|  group: release-${{ github.event_name }}
A|  cancel-in-progress: false
A|jobs:
A|  release:
A|    name: release
A|    runs-on: ubuntu-24.04
A|    timeout-minutes: 60
A|    permissions:
A|      contents: write
B|      packages: write
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
B|        env:
B|          RELEASE_DOCKERHUB_USER: ${{ secrets.DOCKERHUB_USERNAME }}
B|          RELEASE_DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}
A|        run: bash tools/harness/release.sh eingaben
A|      - name: Binaries
A|        run: bash tools/harness/release.sh binaries
B|      - name: Image
B|        run: bash tools/harness/release.sh image
F|      - name: Scan
F|        run: bash tools/harness/release.sh scan
B|      - name: Veröffentlichen
B|        if: ${{ github.event_name == 'push' }}
B|        env:
B|          RELEASE_GHCR_USER: ${{ github.actor }}
B|          RELEASE_GHCR_TOKEN: ${{ secrets.GITHUB_TOKEN }}
B|          RELEASE_DOCKERHUB_USER: ${{ secrets.DOCKERHUB_USERNAME }}
B|          RELEASE_DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}
B|        run: bash tools/harness/release.sh veroeffentlichen
B|      - name: Vergleichen
B|        if: ${{ github.event_name == 'push' }}
B|        run: bash tools/harness/release.sh vergleichen
B|      - name: Smoke
B|        if: ${{ github.event_name == 'push' }}
B|        run: bash tools/harness/release.sh smoke
B|      - name: Latest
B|        if: ${{ github.event_name == 'push' }}
B|        env:
B|          RELEASE_GHCR_USER: ${{ github.actor }}
B|          RELEASE_GHCR_TOKEN: ${{ secrets.GITHUB_TOKEN }}
B|          RELEASE_DOCKERHUB_USER: ${{ secrets.DOCKERHUB_USERNAME }}
B|          RELEASE_DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}
B|        run: bash tools/harness/release.sh latest
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
# Die Schritte nach dem Checkout: STEP_AB[k] ist die Zeile `- name:` des Schritts k in SOLL, STEP_ID[k] sein Kürzel (Name des Falls); BEDING[j] ist die Zeile `if:` des Schritts BEDING_ID[j].
declare -A SCHRITT_ID=([Eingaben]=eingaben [Binaries]=binaries [Image]=image [Scan]=scan [Veröffentlichen]=veroeffentlichen [Vergleichen]=vergleichen [Smoke]=smoke [Latest]=latest [GitHub-Release]=github-release)
STEP_AB=() STEP_ID=() BEDING=() BEDING_ID=()
schritt_re='^      - name: (.*)$'
for ((i = 0; i < N; i++)); do
  if [[ ${SOLL[i]} =~ $schritt_re && ${BASH_REMATCH[1]} != Checkout ]]; then STEP_AB+=("$i"); STEP_ID+=("${SCHRITT_ID[${BASH_REMATCH[1]}]}"); fi
  if [[ ${SOLL[i]} == *"actions/checkout@<SHA>"* ]]; then PIN=$i; fi
  if [[ ${SOLL[i]} == *"if: \${{ github.event_name == 'push' }}"* ]]; then BEDING+=("$i"); BEDING_ID+=("${STEP_ID[-1]}"); fi
done
K=${#STEP_AB[@]}
STEP_AB+=("$N")
PIN_RE="^${SOLL[PIN]%%<SHA>*}[0-9a-f]{40} # [^ ].*\$"
SHA_PROBE=0123456789abcdef0123456789abcdef01234567

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
)
ZUSATZ_NAMEN=(pull-request-target secrets-debug schritt-cache environment continue-on-error always)
zusatz() { neu; urteil rot < <(zeilen "${t[@]}" "$1"); }
for k in "${!ZUSATZ[@]}"; do eval "f_p1_zusatz_${ZUSATZ_NAMEN[k]//-/_}() { zusatz \"\${ZUSATZ[$k]#*|}\"; }"; done
pin() { neu; t[PIN]="        $1"; urteil rot < <(zeilen "${t[@]}"); }
f_p1_pin_kurz() { pin "uses: actions/checkout@${SHA_PROBE:1} # v7.0.1"; }
f_p1_pin_gross() { pin "uses: actions/checkout@${SHA_PROBE^^} # v7.0.1"; }
f_p1_pin_tag_leer() { pin "uses: actions/checkout@$SHA_PROBE # "; }
f_p1_pin_ohne_tag() { pin "uses: actions/checkout@$SHA_PROBE"; }
f_p1_pin_andere_action() { pin "uses: actions/cache@$SHA_PROBE # v4"; }
# Je Schritt mit Bedingung auf das Ereignis `push` ein Fall, in dem die Zeile `if:` fehlt.
for j in "${!BEDING[@]}"; do eval "f_p1_bedingung_${BEDING_ID[j]//-/_}() { neu; unset \"t[${BEDING[j]}]\"; urteil rot < <(zeilen \"\${t[@]}\"); }"; done
# umgestellt <k>… — der Text mit den Schritten in der genannten Reihenfolge (Zahl: Stelle des Schritts im Text; der Checkout bleibt vorn).
umgestellt() {
  local k
  zeilen "${t[@]:0:STEP_AB[0]}"
  for k in "$@"; do zeilen "${t[@]:STEP_AB[k]:STEP_AB[k+1]-STEP_AB[k]}"; done
}
# tausch <k> <l> — die Reihenfolge der Schritte mit den Stellen k und l vertauscht.
tausch() { local o=() j; for ((j = 0; j < K; j++)); do o+=("$j"); done; o[$1]=$2; o[$2]=$1; echo "${o[*]}"; }
# Paare vertauschter Schritte: <Name> <Stelle> <Stelle>
TAUSCHE=("eingaben-binaries 0 1" "binaries-release 1 $((K - 1))" "binaries-image 1 2" "image-scan 2 3" "scan-veroeffentlichen 3 4" "image-veroeffentlichen 2 4" "veroeffentlichen-vergleichen 4 5" "vergleichen-smoke 5 6" "smoke-latest 6 7" "latest-release 7 8")
for tp in "${TAUSCHE[@]}"; do read -r tn tk tl <<< "$tp"; eval "f_p1_reihenfolge_${tn//-/_}() { neu; urteil rot < <(umgestellt \$(tausch $tk $tl)); }"; done
f_p1_leerzeile() { neu; urteil ok < <(zeilen "${t[@]:0:3}" "" "${t[@]:3}"); }
f_p1_kommentarzeile() { neu; urteil ok < <(zeilen "${t[@]:0:3}" "# Kommentar" "${t[@]:3}"); }
f_p1_kommentarzeile_eingerueckt() { neu; urteil ok < <(zeilen "${t[@]:0:3}" "    # Kommentar" "${t[@]:3}"); }

# --- Welt, Attrappen, Hilfen für release.sh -------------------------------------------------------------------------------------------------------------
TOKEN=ghs_GeheimerTokenWert0123456789
TOKEN_D=dckr_pat_GeheimerDockerWert0123
TOKEN_G=ghs_GeheimerRegistryWert0123
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
      http404) echo "HTTP 404: Not Found" >&2; exit 1 ;;
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
  cat > "$b/gp-layout" << 'EOT'
#!/usr/bin/env bash
# gp-layout <Variante> <Tar> — schreibt ein OCI-Layout als Tar und die Bytes des Index nach <Tar>.index
v="$1"; out="$2"; d="$(mktemp -d)"; trap 'rm -rf "$d"' EXIT
h() { printf '%s' "$1" | sha256sum | cut -d' ' -f1; }
m() { printf '{"mediaType":"application/vnd.oci.image.manifest.v1+json","digest":"sha256:%s","size":1,"platform":{"architecture":"%s","os":"linux"}}' "$(h "m-$1")" "$1"; }
case "$v" in
  drei) l="$(m amd64),$(m arm64),{\"mediaType\":\"application/vnd.oci.image.manifest.v1+json\",\"digest\":\"sha256:$(h m-x)\",\"size\":1}" ;;
  ohne-amd64) l="$(m arm64),$(m s390x)" ;;
  ohne-arm64) l="$(m amd64),$(m s390x)" ;;
  eine) l="$(m amd64)" ;;
  ohne-eintrag) l="" ;;
  *) l="$(m amd64),$(m arm64)" ;;
esac
idx="{\"schemaVersion\":2,\"mediaType\":\"application/vnd.oci.image.index.v1+json\",\"manifests\":[$l]}"
x="$(h "$idx")"
mkdir -p "$d/blobs/sha256"; printf '%s' "$idx" > "$d/blobs/sha256/$x"; printf '%s' "$idx" > "$out.index"
e="{\"mediaType\":\"application/vnd.oci.image.index.v1+json\",\"digest\":\"sha256:$x\",\"size\":${#idx}}"
[ "$v" != zwei-eintraege ] || e="$e,$e"
[ "$v" != manifest-statt-index ] || e="${e/image.index.v1/image.manifest.v1}"
[ "$v" != ohne-eintrag ] || e=""
printf '{"schemaVersion":2,"manifests":[%s]}' "$e" > "$d/index.json"
printf '{"imageLayoutVersion":"1.0.0"}' > "$d/oci-layout"
tar -cf "$out" -C "$d" oci-layout index.json blobs
EOT
  cat > "$b/docker" << 'EOT'
#!/usr/bin/env bash
echo "docker $*" >> "$GP_LOG/docker"
if [ "${1-}" = run ]; then
  shift; plat=""; cache=""
  while [ $# -gt 0 ]; do
    case "$1" in --platform) plat="$2" ;; -v) case "$2" in *:/cache) cache="${2%:/cache}" ;; esac ;; esac
    shift
  done
  a="${plat#linux/}"
  modus() { local x; x="$(tr ' ' '\n' <<< "${!1-}" | sed -n "s/^$a=//p" | head -n 1)"; echo "${x:-$2}"; }
  u() { date -u -d "@$(( $(date +%s) - $1 * 3600 ))" +%Y-%m-%dT%H:%M:%S.123456789Z; }
  mkdir -p "$cache/db"
  case "$(modus GP_DB frisch)" in
    frisch) printf '{"Version":2,"UpdatedAt":"%s","DownloadedAt":"%s"}' "$(u 1)" "$(u 0)" > "$cache/db/metadata.json" ;;
    h71) printf '{"Version":2,"UpdatedAt":"%s","DownloadedAt":"%s"}' "$(u 71)" "$(u 0)" > "$cache/db/metadata.json" ;;
    h73) printf '{"Version":2,"UpdatedAt":"%s","DownloadedAt":"%s"}' "$(u 73)" "$(u 0)" > "$cache/db/metadata.json" ;;
    fehlt) rm -f "$cache/db/metadata.json" ;;
    kein-json) printf 'kein json' > "$cache/db/metadata.json" ;;
    ohne-feld) printf '{"Version":2}' > "$cache/db/metadata.json" ;;
    leer) printf '{"UpdatedAt":""}' > "$cache/db/metadata.json" ;;
    zahl) printf '{"UpdatedAt":5}' > "$cache/db/metadata.json" ;;
    jetzt) printf '{"UpdatedAt":"now"}' > "$cache/db/metadata.json" ;;
    kein-datum) printf '{"UpdatedAt":"2026-13-45T00:00:00Z"}' > "$cache/db/metadata.json" ;;
  esac
  V() { local f=""; case "$4" in -) ;; leer) f=',"FixedVersion":""' ;; null) f=',"FixedVersion":null' ;; *) f=",\"FixedVersion\":\"$4\"" ;; esac
    printf '{"VulnerabilityID":"%s","PkgName":"%s","InstalledVersion":"1.0","Severity":"%s"%s}' "$1" "$3" "$2" "$f"; }
  R() { printf '{"Target":"t","Type":"debian","Vulnerabilities":[%s]}' "$1"; }
  ar="$(modus GP_ARCH "$a")"; ic="\"architecture\":\"$ar\","; [ "$ar" != - ] || ic=""
  J() { printf '{"SchemaVersion":2,"Metadata":{"ImageConfig":{%s"os":"linux"}},"Results":%s}\n' "$ic" "$1"; }
  sauber='{"Target":"a"},{"Target":"b","Vulnerabilities":null}'
  case "$(modus GP_TRIVY sauber)" in
    sauber) J "[$sauber]" ;;
    exit1) J "[$sauber]"; echo "FATAL scan" >&2; exit 1 ;;
    exit125) exit 125 ;;
    kein-json) echo "kein json" ;;
    leer) : ;;
    ergebnis-leer) J "[]" ;;
    ergebnis-null) J null ;;
    ergebnis-fehlt) printf '{"SchemaVersion":2,"Metadata":{"ImageConfig":{%s"os":"linux"}}}\n' "$ic" ;;
    critical-fix) J "[$(R "$(V CVE-1 CRITICAL pkg-c 1.2.3)")]" ;;
    high-fix) J "[$(R "$(V CVE-1 HIGH pkg-h 1.2.3)")]" ;;
    medium-fix) J "[$(R "$(V CVE-1 MEDIUM pkg-m 1.2.3)")]" ;;
    low-fix) J "[$(R "$(V CVE-1 LOW pkg-l 1.2.3)")]" ;;
    unknown-fix) J "[$(R "$(V CVE-1 UNKNOWN pkg-u 1.2.3)")]" ;;
    high-ohne-fix) J "[$(R "$(V CVE-1 HIGH pkg-h -)")]" ;;
    high-ohne-fix-leer) J "[$(R "$(V CVE-1 HIGH pkg-h leer)")]" ;;
    high-ohne-fix-null) J "[$(R "$(V CVE-1 HIGH pkg-h null)")]" ;;
    critical-ohne-fix) J "[$(R "$(V CVE-1 CRITICAL pkg-c -)")]" ;;
    zweites-ergebnis) J "[{\"Target\":\"a\"},$(R "$(V CVE-Z HIGH pkg-z 7.0)")]" ;;
    befunde-a) J "[$(R "$(V CVE-A1 CRITICAL pkg-a1 2.0)"),$(R "$(V CVE-A2 HIGH pkg-a2 3.0)")]" ;;
    befunde-b) J "[$(R "$(V CVE-B1 HIGH pkg-b1 4.0),$(V CVE-B2 MEDIUM pkg-b2 5.0),$(V CVE-B3 HIGH pkg-b3 -)")]" ;;
    summe) J "[$(R "$(V S1 CRITICAL p 1),$(V S2 HIGH p 2),$(V S3 HIGH p -),$(V S4 CRITICAL p leer),$(V S5 MEDIUM p 9),$(V S6 LOW p -)")]" ;;
  esac
  exit 0
fi
[ "${1-}" = buildx ] || exit 0
case "${2-}" in
  create) [ "${GP_BUILDX-}" != create-rot ] || exit 1 ;;
  build)
    shift 2; plat=""; out=""
    while [ $# -gt 0 ]; do
      case "$1" in --platform) plat="$2" ;; --output) out="$2" ;; esac
      shift
    done
    [ "${GP_BUILDX-}" != build-rot ] || exit 1
    dest="$(sed -E 's/.*dest=([^,]*).*/\1/' <<< "$out")"
    if [[ $dest == *.tar ]]; then
      gp-layout "${GP_IMAGE-ok}" "$dest"
    else
      mkdir -p "$dest"; a="${plat#linux/}"
      printf 'binary %s linux %s\n' "$(head -n 1 VERSION)" "$a" > "$dest/pgwire-recorder"
      [ "${GP_OHNE_BINARY-}" != "$a" ] || rm "$dest/pgwire-recorder"
      [ "${GP_SUMME-}" != "$a" ] || echo manipuliert >> "$dest/pgwire-recorder"
    fi ;;
esac
exit 0
EOT
  cat > "$b/skopeo" << 'EOT'
#!/usr/bin/env bash
echo "skopeo $*" >> "$GP_LOG/skopeo"
cmd="$1"; shift
raw=0 anon=0 pws=0 user="" args=()
while [ $# -gt 0 ]; do
  case "$1" in
    --raw) raw=1 ;;
    --no-creds) anon=1 ;;
    --authfile | --password) shift ;;
    --username) user="$2"; shift ;;
    --password-stdin) pws=1 ;;
    --all | --preserve-digests) ;;
    *) args+=("$1") ;;
  esac
  shift
done
# teile <Referenz> — setzt reg und ref (Tag oder sha256:…) einer Referenz docker://<Registry>/pt9912/pgwire-recorder[:@]<ref>
teile() { local r="${1#docker://}"; reg="${r%%/*}"; r="${r#*/pt9912/pgwire-recorder}"; ref="${r:1}"; }
case "$cmd" in
  login)
    reg="${args[-1]}"
    if [ "$pws" = 1 ]; then cat > "$GP_LOG/stdin.$reg"; else : > "$GP_LOG/stdin.$reg"; fi
    echo "$user" > "$GP_LOG/user.$reg"
    [ "${GP_LOGIN_ROT-}" != "$reg" ] || exit 1 ;;
  inspect)
    teile "${args[0]}"
    if [ "$anon" = 1 ] && [[ " ${GP_ANON-} " == *" $reg "* ]]; then echo "unauthorized: authentication required" >&2; exit 1; fi
    case "${GP_INSPECT-}" in
      "$reg") echo "received unexpected HTTP status: 502 Bad Gateway" >&2; exit 1 ;;
      "denied:$reg") echo "requested access to the resource is denied" >&2; exit 1 ;;
    esac
    f="$GP_REG/$reg/$ref"
    if [ ! -f "$f" ]; then
      case "${GP_FEHLT-manifest}" in
        manifest) echo "manifest unknown" >&2 ;;
        name) echo "name unknown: repository name not known to registry" >&2 ;;
        notfound) echo "Error: 404 Not Found" >&2 ;;
      esac
      exit 1
    fi
    if [ "$raw" = 1 ]; then cat "$f"; else printf '{"Name":"x","Digest":"sha256:%s"}\n' "$(tr -d ' \n' < "$f" | sed 's/x/y/' | sha256sum | cut -d' ' -f1)"; fi ;;
  copy)
    src="${args[0]}"; teile "${args[1]}"; dreg="$reg"; dref="$ref"
    [[ " ${GP_COPY_ROT-} " != *" $dreg "* ]] || { echo "copy nach $dreg verweigert" >&2; exit 1; }
    mkdir -p "$GP_REG/$dreg"
    case "$src" in
      oci-archive:*)
        t="${src#oci-archive:}"
        x="$(tar -xOf "$t" index.json | grep -o '[0-9a-f]\{64\}' | head -n 1)"
        tar -xOf "$t" "blobs/sha256/$x" > "$GP_REG/$dreg/$dref" ;;
      docker://*)
        teile "$src"
        for g in "$GP_REG/$reg"/*; do
          if [ "sha256:$(sha256sum "$g" | cut -d' ' -f1)" = "$ref" ]; then cp "$g" "$GP_REG/$dreg/$dref"; break; fi
        done
        [ -f "$GP_REG/$dreg/$dref" ] || exit 1 ;;
    esac ;;
esac
exit 0
EOT
  chmod +x "$b"/*
}
# welt <dir> — legt quelle/ an (main: c1 c2; Zweig side mit s1 aus c1; refs/remotes/origin/main auf c1, origin/trunk auf s1) und bin/ mit den Attrappen.
welt() {
  local q="$1/quelle"
  mkdir -p "$q"
  git -C "$q" init -q -b main
  GIT_AUTHOR_DATE="1700000000 +0000" GIT_COMMITTER_DATE="1700000000 +0000" git -C "$q" commit -q --allow-empty -m c1
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
# Überschrift, die nicht auf Änderungshistorie endet) · pipe, pipe-maskiert (ein weiteres `|` im Text, roh und als `\|`) · spaeter (die Version nur in einem Abschnitt danach) · leer (Zeile der Version mit leerer dritter Zelle) · ohne (keine Zeile).
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
      pipe) printf '| %s | Datum-%s | Text | mehr |\n' "$v" "$v" ;;
      pipe-maskiert) printf '| %s | Datum-%s | Text \\| mehr |\n' "$v" "$v" ;;
      nur-anderer | spaeter | ohne) printf '| 9.9.9 | 2026-10-02 | Anderes |\n' ;;
    esac
    if [ "$m" = spaeter ]; then printf '\n## 12. Anhang\n\n| Version | Datum | Änderung |\n|---|---|---|\n| %s | 2020-01-01 | Falscher-Abschnitt |\n' "$v"; fi
  } > "$f"
}
# kopie — der Baum des Falls: repo/ (Kopie der Welt), log/, runner/, tmp/.
kopie() { w="$run/t.$fallname"; rm -rf "$w"; mkdir -p "$w/log" "$w/runner" "$w/tmp" "$w/registry"; cp -a "$run/welt/quelle" "$w/repo"; }
# fix [Version [Handbuch-Modus]] — setzt V, legt den Baum an, schreibt VERSION und das Handbuch.
PGIMG='postgres:test@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef'
fix() {
  V="${1:-0.1.0}"; kopie; printf '%s\n' "$V" > "$w/repo/VERSION"; handbuch "$V" "${2:-normal}"
  mkdir -p "$w/repo/harness/mk" "$w/repo/tools/harness"
  printf 'POSTGRES_IMAGE ?= %s\n' "$PGIMG" > "$w/repo/harness/mk/integration.mk"
  cat > "$w/repo/tools/harness/container-smoke.sh" << 'EOT'
#!/usr/bin/env bash
echo "ref=$1 VERSION=$VERSION POSTGRES_IMAGE=$POSTGRES_IMAGE SMOKE_NAME=$SMOKE_NAME LAUF=$SMOKE_LAUF WERK=$([ -d "$SMOKE_WERK" ] && echo da || echo fehlt)" >> "$GP_LOG/smoke"
[[ " ${GP_SMOKE_ROT-} " != *" ${1%%/*} "* ]] || { echo "smoke rot" >&2; exit 1; }
EOT
}
# umg <lauf|trocken> [NAME=Wert | -NAME]… — setzt E: die Umgebung von release.sh, dann die Abweichungen.
umg() {
  local art="$1" a
  local -A v=([RELEASE_EVENT]=push [RELEASE_REF]="refs/tags/v$V" [RELEASE_SHA]="$c1" [RELEASE_DEFAULT_BRANCH]=main [RELEASE_REPOSITORY]=pt9912/pgwire-recorder
    [RUNNER_TEMP]="$w/runner" [TMPDIR]="$w/tmp" [GP_LOG]="$w/log" [GH_TOKEN]="$TOKEN" [GP_REG]="$w/registry"
    [RELEASE_DOCKERHUB_USER]=dockeruser [RELEASE_DOCKERHUB_TOKEN]="$TOKEN_D" [RELEASE_GHCR_USER]=ghuser [RELEASE_GHCR_TOKEN]="$TOKEN_G")
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
  if [ "$s" = github-release ] && [ -z "${ohne_vergl-}" ]; then mkdir -p "$w/runner/release"; : > "$w/runner/release/smoke.ok"; fi
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
f_p3_notes_pipe() {
  local m
  for m in pipe pipe-maskiert; do
    fix 0.1.0 "$m"; umg lauf; lauf eingaben; grund 2 eingaben notes
    fix 0.1.0 "$m"; umg trocken; lauf eingaben; grund 2 eingaben notes
  done
}
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
  fix 0.1.0; umg lauf; bis eingaben binaries
  lauf github-release GP_GH=http404
  sei "HTTP 404: Exit 1" test "$code" -eq 1; sei "HTTP 404: kein create" ohne "$(aufrufe gh)" "release create"
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

# --- Fassung B: Image, Registries, latest (Punkte 2, 3 (7), 5 bis 11) ----------------------------------------------------------------------------------
IMG() { echo "$w/runner/release/image"; }
REG() { echo "$w/registry/$1"; }
# rotgrund <Schritt> [Text] — Exit 1 und eine Meldung `release: <Schritt>: ROT`, die den Text nennt.
rotgrund() { sei "Exit 1" test "$code" -eq 1; sei "ROT in $1" enthaelt "$err" "release: $1: ROT"; if [ -n "${2-}" ]; then sei "Meldung nennt: $2" enthaelt "$err" "$2"; fi; }
# ibild [Version] — Baum und Lauf bis image, setzt D; bild — dazu der Scan.
ibild() { fix "${1:-0.1.0}"; umg lauf; bis eingaben binaries image; D="$(cat "$(IMG)/digest")"; }
bild() { ibild "$@"; bis scan; }
bis_vergl() { bild "$@"; bis veroeffentlichen vergleichen; }
# tamper <Registry> — ändert die Bytes des Tags, ohne den Index zu ändern (Leerzeichen am Ende).
tamper() { printf ' ' >> "$(REG "$1")/0.1.0"; }
# verben — die Anmeldungen und Kopien in der Reihenfolge des Aufrufs, `login <Registry>` und `copy <Registry>`.
verben() { awk '$2 == "login" { print "login", $NF } $2 == "copy" { d = $NF; sub("docker://", "", d); sub("/.*", "", d); print "copy", d }' "$w/log/skopeo" | tr '\n' ' '; }
kopien() { grep -c "^skopeo copy .* docker://$1/" "$w/log/skopeo" || true; }
# Punkt 2
f_p2_trockenlauf_image() {
  fix 0.3.0-dev; umg trocken; bis eingaben binaries image
  sei "skopeo nie gerufen" nie skopeo; sei "gh nie gerufen" nie gh; sei "kein login" ohne "$(aufrufe docker)" "login"
  sei "buildx gerufen" enthaelt "$(aufrufe docker)" "buildx build"; sei "Index-Digest geschrieben" test -s "$(IMG)/digest"
}
# aussen_b <Umgebung…> — jeder Schritt mit Außenwirkung von B lehnt ab, bevor skopeo, gh oder der Smoke gerufen werden.
aussen_b() {
  local s
  for s in veroeffentlichen vergleichen smoke latest; do
    bild; umg lauf "$@"; lauf "$s"; grund 2 "$s" ereignis
    sei "skopeo, gh und der Smoke nie gerufen ($s)" test ! -s "$w/log/skopeo" -a ! -s "$w/log/gh" -a ! -s "$w/log/smoke"
  done
}
f_p2_aussenwirkung_ohne_push_b() { aussen_b RELEASE_EVENT=workflow_dispatch; }
f_p2_ref_kein_tag_b() { aussen_b RELEASE_REF=refs/heads/main; }
# Punkt 3 (7): Lauf und Trockenlauf; die Meldung nennt den Namen, nie einen Wert
secret() { # <Direktive für umg> <Name>
  local art
  for art in lauf trocken; do
    fix 0.1.0; umg "$art" "$1"; lauf eingaben; grund 2 eingaben geheimnis
    sei "nennt $2" enthaelt "$err" "$2"; sei "nennt keinen Wert" ohne "$out$err" "$TOKEN_D"
  done
}
f_p3_secret_fehlt_user() { secret -RELEASE_DOCKERHUB_USER DOCKERHUB_USERNAME; }
f_p3_secret_fehlt_token() { secret -RELEASE_DOCKERHUB_TOKEN DOCKERHUB_TOKEN; }
f_p3_secret_leer_user() { secret RELEASE_DOCKERHUB_USER= DOCKERHUB_USERNAME; }
f_p3_secret_leer_token() { secret RELEASE_DOCKERHUB_TOKEN= DOCKERHUB_TOKEN; }
f_p3_secret_fehlen_beide() {
  fix 0.1.0; umg lauf -RELEASE_DOCKERHUB_USER -RELEASE_DOCKERHUB_TOKEN; lauf eingaben; grund 2 eingaben geheimnis
  sei "nennt beide Namen" enthaelt "$err" "DOCKERHUB_USERNAME"; sei "nennt DOCKERHUB_TOKEN" enthaelt "$err" "DOCKERHUB_TOKEN"; sei "eine Meldung" test "$(grep -c . <<< "$err")" -eq 1
}
f_p3_schluesselwort_geheimnis() { fix 0.1.0; umg lauf -RELEASE_DOCKERHUB_TOKEN; lauf eingaben; sei "Schlüsselwort geheimnis" beginnt geheimnis; }
f_p3_reihenfolge_6_7() { fix 0.1.0; umg lauf "RELEASE_SHA=$s1" -RELEASE_DOCKERHUB_TOKEN; lauf eingaben; grund 2 eingaben standard-zweig; }
# Punkt 5
bild_mit() { fix 0.1.0; umg lauf; bis eingaben binaries; lauf image "$@"; }
f_p5_aufruf() {
  local b c x
  bild
  b="$(grep '^docker buildx build ' "$w/log/docker" | grep -e 'type=oci' || true)"
  for x in "--builder pgr-release" "--platform linux/amd64,linux/arm64" "--target runtime" "--build-arg VERSION=0.1.0" "--build-arg SOURCE_DATE_EPOCH=1700000000" "--provenance=false" "--sbom=false" \
    "--output type=oci,dest=$(IMG)/layout.tar,rewrite-timestamp=true"; do
    sei "Aufruf: $x" enthaelt "$b" " $x"
  done
  c="$(grep '^docker buildx create ' "$w/log/docker" || true)"
  sei "Treiber docker-container" enthaelt "$c" " --driver docker-container "
  sei "Image des Builders: der Pin aus release.sh" enthaelt "$c" " --driver-opt image=$(pin_zeile | sed "s/^BUILDKIT_IMAGE='\([^']*\)'.*/\1/") "
}
f_p5_zwei_manifeste() { bild_mit GP_IMAGE=drei; rotgrund image "nicht genau zwei Manifeste"; sei "kein Digest" test ! -e "$(IMG)/digest"; }
f_p5_plattform_fehlt_amd64() { bild_mit GP_IMAGE=ohne-amd64; rotgrund image "statt linux/amd64 linux/arm64"; sei "kein Digest" test ! -e "$(IMG)/digest"; }
f_p5_plattform_fehlt_arm64() { bild_mit GP_IMAGE=ohne-arm64; rotgrund image "statt linux/amd64 linux/arm64"; sei "kein Digest" test ! -e "$(IMG)/digest"; }
f_p5_layout_form() {
  local v
  for v in zwei-eintraege:"nicht genau einen Index" ohne-eintrag:"nicht genau einen Index" manifest-statt-index:"kein Index"; do
    bild_mit "GP_IMAGE=${v%%:*}"; rotgrund image "${v#*:}"; sei "kein Digest ($v)" test ! -e "$(IMG)/digest"
  done
}
f_p5_ohne_binaries() { fix 0.1.0; umg lauf; bis eingaben; lauf image; grund 2 image vorbedingung; sei "docker nie gerufen" nie docker; }
# der Pin des Builders steht in release.sh (Kopie des Falls)
pin_zeile() { grep '^BUILDKIT_IMAGE=' "$releasesh"; }
f_p5_builder_pin_form() { sei "Form: moby/buildkit@sha256:<64 Hexziffern> und Tag-Kommentar" test -n "$(pin_zeile | grep -E "^BUILDKIT_IMAGE='moby/buildkit@sha256:[0-9a-f]{64}' # v[0-9]+\.[0-9]+\.[0-9]+\$" || true)"; }
f_p5_builder_platzhalter() { sei "der Digest ist kein Platzhalter aus Nullen" test -z "$(pin_zeile | grep -E "@sha256:0{64}'" || true)"; }
f_p5_epoch_committer() {
  local c3
  fix 0.1.0
  GIT_AUTHOR_DATE="1600000000 +0000" GIT_COMMITTER_DATE="1700000500 +0000" git -C "$w/repo" commit -q --allow-empty -m c3
  c3="$(git -C "$w/repo" rev-parse HEAD)"; git -C "$w/repo" update-ref refs/remotes/origin/main "$c3"
  umg lauf "RELEASE_SHA=$c3"; bis eingaben binaries image
  sei "SOURCE_DATE_EPOCH ist die Committer-Zeit" enthaelt "$(grep '^docker buildx build ' "$w/log/docker" | grep -e 'type=oci')" " --build-arg SOURCE_DATE_EPOCH=1700000500 "
}
f_p5_summe_amd64() { bild_mit GP_SUMME=amd64; rotgrund image "linux/amd64"; sei "kein Digest" test ! -e "$(IMG)/digest"; }
f_p5_summe_arm64() { bild_mit GP_SUMME=arm64; rotgrund image "linux/arm64"; sei "nicht amd64" ohne "$err" "linux/amd64: das Binary"; sei "kein Digest" test ! -e "$(IMG)/digest"; }
f_p5_summe_eintrag_fehlt() {
  fix 0.1.0; umg lauf; bis eingaben binaries
  grep -v '_linux_arm64$' "$(DIST)/SHA256SUMS" > "$w/s"; mv "$w/s" "$(DIST)/SHA256SUMS"
  lauf image GP_OHNE_BINARY=arm64; rotgrund image "linux/arm64"; sei "kein Digest" test ! -e "$(IMG)/digest"
}
f_p5_digest() {
  bild
  sei "D ist der Hash des Index-Blobs" test "$D" = "sha256:$(sha256sum "$(IMG)/layout.tar.index" | cut -d' ' -f1)"
  sei "D ist nicht der Hash des Archivs" test "$D" != "sha256:$(sha256sum "$(IMG)/layout.tar" | cut -d' ' -f1)"
  bis veroeffentlichen
  sei "die Registry trägt die Bytes mit D" test "sha256:$(sha256sum "$(REG ghcr.io)/0.1.0" | cut -d' ' -f1)" = "$D"
}
f_p5_bau_scheitert() {
  bild_mit GP_BUILDX=create-rot; rotgrund image "buildx create"
  bild_mit GP_BUILDX=build-rot; rotgrund image "buildx build"
}
# Punkt 6
f_p6_reihenfolge() { bild; bis veroeffentlichen; sei "login ghcr, copy ghcr, login docker, copy docker" test "$(verben)" = "login ghcr.io copy ghcr.io login docker.io copy docker.io "; }
f_p6_all_und_digests() {
  local z r
  bild; bis veroeffentlichen
  z="$(grep '^skopeo copy ' "$w/log/skopeo")"
  for r in ghcr.io docker.io; do
    sei "$r: --all" enthaelt "$z" " --all "; sei "$r: --preserve-digests" enthaelt "$z" " --preserve-digests "
    sei "$r: Quelle ist das Layout" enthaelt "$z" " oci-archive:$(IMG)/layout.tar docker://$r/pt9912/pgwire-recorder:0.1.0"
  done
}
f_p6_tag_vorhanden_gleich() {
  bild; mkdir -p "$(REG ghcr.io)"; cp "$(IMG)/layout.tar.index" "$(REG ghcr.io)/0.1.0"
  lauf veroeffentlichen; gruen
  sei "kein Copy nach ghcr.io" test "$(kopien ghcr.io)" -eq 0; sei "Copy nach docker.io" test "$(kopien docker.io)" -eq 1
}
f_p6_tag_vorhanden_anders() {
  bild; mkdir -p "$(REG ghcr.io)"; printf '{"anders":true}' > "$(REG ghcr.io)/0.1.0"; cp "$(REG ghcr.io)/0.1.0" "$w/vorher"
  lauf veroeffentlichen; rotgrund veroeffentlichen "ghcr.io"
  sei "Meldung nennt beide Digests" enthaelt "$err" "sha256:$(sha256sum "$w/vorher" | cut -d' ' -f1) statt $D"
  sei "kein Copy" test "$(kopien ghcr.io)" -eq 0; sei "der Tag bleibt, wie er war" cmp -s "$w/vorher" "$(REG ghcr.io)/0.1.0"
}
tag_fehlt() { bild; lauf veroeffentlichen "GP_FEHLT=$1"; gruen; sei "kopiert nach ghcr.io" test "$(kopien ghcr.io)" -eq 1; sei "kopiert nach docker.io" test "$(kopien docker.io)" -eq 1; }
f_p6_tag_fehlt_manifest_unknown() { tag_fehlt manifest; }
f_p6_tag_fehlt_name_unknown() { tag_fehlt name; }
f_p6_tag_fehlt_not_found() { tag_fehlt notfound; }
f_p6_inspect_anderer_fehler() {
  local f
  for f in ghcr.io denied:ghcr.io; do
    bild; lauf veroeffentlichen "GP_INSPECT=$f"; rotgrund veroeffentlichen; sei "$f: kein Copy" test "$(kopien ghcr.io)" -eq 0
  done
}
f_p6_inspect_meldung_im_log() {
  bild; lauf veroeffentlichen GP_INSPECT=ghcr.io; sei "Fehler im Log (veroeffentlichen)" enthaelt "$err" "release: skopeo: received unexpected HTTP status: 502 Bad Gateway"
  bild; lauf veroeffentlichen; gruen; sei "Tag fehlt im Log (veroeffentlichen)" enthaelt "$err" "release: skopeo: manifest unknown"
  bild; bis veroeffentlichen; lauf vergleichen GP_ANON=ghcr.io; sei "Fehler im Log (vergleichen)" enthaelt "$err" "release: skopeo: unauthorized: authentication required"
  bild; bis veroeffentlichen; rm "$(REG docker.io)/0.1.0"; lauf vergleichen; sei "Tag fehlt im Log (vergleichen)" enthaelt "$err" "release: skopeo: manifest unknown"
}
f_p6_zweite_registry_faellt() {
  bild; lauf veroeffentlichen GP_COPY_ROT=docker.io; rotgrund veroeffentlichen "ghcr.io:0.1.0=$D docker.io:0.1.0=-"
  sei "ghcr.io bleibt beschrieben" test -f "$(REG ghcr.io)/0.1.0"
}
f_p6_tag_satz() {
  bild; bis veroeffentlichen
  sei "ghcr.io trägt nur 0.1.0" test "$(ls "$(REG ghcr.io)")" = 0.1.0; sei "docker.io trägt nur 0.1.0" test "$(ls "$(REG docker.io)")" = 0.1.0
  sei "jedes Ziel hat den Tag 0.1.0" test "$(grep '^skopeo copy ' "$w/log/skopeo" | grep -vc ':0.1.0$' || true)" -eq 0
}
f_p6_anmeldung_scheitert() { bild; lauf veroeffentlichen GP_LOGIN_ROT=ghcr.io; rotgrund veroeffentlichen "Anmeldung bei ghcr.io"; sei "kein Copy" test "$(kopien ghcr.io)" -eq 0; }
# Punkt 7
f_p7_gleich() {
  bis_vergl
  sei "Marke vergleichen.ok" test -e "$w/runner/release/vergleichen.ok"
  sei "ghcr.io ohne Anmeldung, rohe Bytes" test "$(grep -c '^skopeo inspect --raw --no-creds docker://ghcr.io/pt9912/pgwire-recorder:0.1.0$' "$w/log/skopeo" || true)" -eq 1
  sei "docker.io ohne Anmeldung, rohe Bytes" test "$(grep -c '^skopeo inspect --raw --no-creds docker://docker.io/pt9912/pgwire-recorder:0.1.0$' "$w/log/skopeo" || true)" -eq 1
}
f_p7_rohbytes() { bis_vergl; sei "gleich, weil die rohen Bytes gleich sind" test -e "$w/runner/release/vergleichen.ok"; }
digest_abweichend() { # <Registry>
  bild; bis veroeffentlichen; tamper "$1"; lauf vergleichen; rotgrund vergleichen "$1 weicht ab"; sei "keine Marke" test ! -e "$w/runner/release/vergleichen.ok"
}
f_p7_digest_ghcr() { digest_abweichend ghcr.io; sei "docker.io nicht genannt" ohne "$err" "docker.io weicht ab"; }
f_p7_digest_docker() { digest_abweichend docker.io; sei "ghcr.io nicht genannt" ohne "$err" "ghcr.io weicht ab"; }
f_p7_digest_lokal() {
  bild; bis veroeffentlichen; printf 'sha256:%064d\n' 0 > "$(IMG)/digest"; lauf vergleichen; rotgrund vergleichen "weicht ab"; sei "keine Marke" test ! -e "$w/runner/release/vergleichen.ok"
}
anonym() { bild; bis veroeffentlichen; lauf vergleichen "GP_ANON=$1"; rotgrund vergleichen "$1 ist ohne Anmeldung nicht lesbar"; sei "nennt die Handlung" enthaelt "$err" "öffentlich"; sei "keine Marke" test ! -e "$w/runner/release/vergleichen.ok"; }
f_p7_anonym_ghcr() { anonym ghcr.io; }
f_p7_anonym_docker() { anonym docker.io; }
f_p7_eine_plattform() {
  local r
  for r in ghcr.io docker.io; do
    bild; bis veroeffentlichen; "$run/welt/bin/gp-layout" eine "$w/eine.tar"; cp "$w/eine.tar.index" "$(REG "$r")/0.1.0"
    lauf vergleichen; rotgrund vergleichen "der Index aus $r hat nicht genau zwei Manifeste"; sei "keine Marke" test ! -e "$w/runner/release/vergleichen.ok"
  done
}
f_p7_tag_fehlt() {
  local r
  for r in ghcr.io docker.io; do
    bild; bis veroeffentlichen; rm "$(REG "$r")/0.1.0"; lauf vergleichen; rotgrund vergleichen "$r trägt den Tag 0.1.0 nicht"
  done
}
f_p7_smoke_aufruf() {
  local z
  bis_vergl; lauf smoke; gruen
  z="$(aufrufe smoke)"
  sei "zwei Aufrufe" test "$(wc -l <<< "$z")" -eq 2
  sei "erst ghcr.io, mit @D" test "$(head -n 1 <<< "$z" | cut -d' ' -f1)" = "ref=ghcr.io/pt9912/pgwire-recorder@$D"
  sei "dann docker.io, mit @D" test "$(tail -n 1 <<< "$z" | cut -d' ' -f1)" = "ref=docker.io/pt9912/pgwire-recorder@$D"
  sei "VERSION=0.1.0 in beiden" test "$(grep -c ' VERSION=0.1.0 ' <<< "$z")" -eq 2
  sei "POSTGRES_IMAGE aus integration.mk in beiden" test "$(grep -c " POSTGRES_IMAGE=$PGIMG " <<< "$z")" -eq 2
  sei "Meldungspräfix und Verzeichnis in beiden" test "$(grep -c 'SMOKE_NAME=release: smoke LAUF=.* WERK=da$' <<< "$z")" -eq 2
  sei "Namen der Läufe verschieden" test "$(grep -o ' LAUF=[^ ]*' <<< "$z" | sort -u | wc -l)" -eq 2
}
f_p7_postgres_fehlt() { bis_vergl; rm "$w/repo/harness/mk/integration.mk"; lauf smoke; grund 2 smoke vorbedingung; sei "Smoke nie gerufen" nie smoke; }
f_p7_smoke_ohne_vergleichen() { bild; bis veroeffentlichen; lauf smoke; grund 2 smoke vorbedingung; sei "Smoke nie gerufen" nie smoke; }
f_p7_smoke_ok_marke() {
  bis_vergl; lauf smoke; gruen; sei "Marke smoke.ok nach grünem Smoke" test -e "$w/runner/release/smoke.ok"
  lauf smoke GP_SMOKE_ROT=ghcr.io; sei "Smoke rot" test "$code" -eq 1; sei "keine Marke nach rotem Smoke" test ! -e "$w/runner/release/smoke.ok"
}
smoke_rot() { bis_vergl; lauf smoke "GP_SMOKE_ROT=$1"; rotgrund smoke "der Smoke aus $1"; sei "kein latest" ohne "$(aufrufe skopeo)" ":latest"; }
f_p7_smoke_ghcr() { smoke_rot ghcr.io; }
f_p7_smoke_docker() { smoke_rot docker.io; }
# Punkt 8
latest_von() { # <Version> [Tag…] — die Tags liegen im Klon; Lauf bis smoke, setzt D
  local t
  fix "$1"; shift
  for t in "$@"; do git -C "$w/repo" tag "$t" "$c1"; done
  umg lauf; bis eingaben binaries image scan veroeffentlichen vergleichen smoke
  D="$(cat "$(IMG)/digest")"
}
gesetzt() { local r; for r in ghcr.io docker.io; do sei "$r: latest zeigt auf D" test "sha256:$(sha256sum "$(REG $r)/latest" 2> /dev/null | cut -d' ' -f1)" = "$D"; done; }
ungesetzt() { sei "latest in keiner Registry" test ! -e "$(REG ghcr.io)/latest" -a ! -e "$(REG docker.io)/latest"; sei "kein Copy nach latest" ohne "$(aufrufe skopeo)" ":latest"; }
f_p8_stabil_hoechste() {
  local z
  latest_von 1.1.0 v1.0.0 v1.1.0-rc.1; lauf latest; gruen; gesetzt
  z="$(grep ':latest$' "$w/log/skopeo")"
  sei "ghcr.io zuerst" test "$(awk '{ d = $NF; sub("docker://", "", d); sub("/.*", "", d); printf "%s ", d }' <<< "$z")" = "ghcr.io docker.io "
  sei "--all und --preserve-digests" test "$(grep -c '^skopeo copy --all --preserve-digests ' <<< "$z")" -eq 2
  sei "Quelle ist @D in derselben Registry" enthaelt "$z" " docker://docker.io/pt9912/pgwire-recorder@$D docker://docker.io/pt9912/pgwire-recorder:latest"
}
f_p8_vorabversion() {
  local v
  for v in 1.2.0-rc.1 1.2.0-beta.1 1.2.0-alpha.1; do latest_von "$v" v1.1.0; lauf latest; gruen; ungesetzt; sei "Meldung ($v)" enthaelt "$out" "Vorabversion"; done
}
f_p8_tag_satz() {
  local z r
  latest_von 1.1.0 v1.0.0; lauf latest; gruen
  z="$(grep '^skopeo copy ' "$w/log/skopeo")"
  sei "jedes Ziel trägt 1.1.0 oder latest" test "$(grep -vcE ' docker://[^ ]*:(1\.1\.0|latest)$' <<< "$z" || true)" -eq 0
  for r in ghcr.io docker.io; do sei "$r trägt nur 1.1.0 und latest" test "$(ls "$(REG $r)" | tr '\n' ' ')" = "1.1.0 latest "; done
}
f_p8_aelterer_zweig() { latest_von 1.0.1 v1.1.0; lauf latest; gruen; ungesetzt; }
f_p8_erster_stabiler() { latest_von 0.1.0; lauf latest; gruen; gesetzt; }
f_p8_vorabversion_hoeher_zaehlt_nicht() { latest_von 1.1.1 v1.1.0 v1.2.0-rc.1; lauf latest; gruen; gesetzt; }
f_p8_semver_numerisch() {
  latest_von 1.10.0 v1.9.0; lauf latest; gruen; gesetzt
  latest_von 1.9.0 v1.10.0; lauf latest; gruen; ungesetzt
}
f_p8_ohne_smoke() { bis_vergl; lauf latest; grund 2 latest vorbedingung; ungesetzt; }
f_p8_fremder_tag() { latest_von 1.0.0 vfoo v1.2; lauf latest; gruen; gesetzt; }
f_p8_nach_rotem_vergleich() { bild; bis veroeffentlichen; tamper ghcr.io; lauf vergleichen; lauf latest; grund 2 latest vorbedingung; ungesetzt; }
f_p8_zweite_registry_faellt() {
  latest_von 1.1.0 v1.0.0; lauf latest GP_COPY_ROT=docker.io; rotgrund latest "ghcr.io:latest=$D docker.io:latest=-"
  sei "ghcr.io trägt latest" test -f "$(REG ghcr.io)/latest"
}
# Punkt 9
f_p9_nicht_vor_registries() { bis_vergl; ohne_vergl=1; lauf github-release; grund 2 github-release vorbedingung; sei "gh nie gerufen" nie gh; }
# Punkt 10
kein_geheimnis() { sei "kein Token in stdout und stderr" ohne "$out$err" "$TOKEN"; sei "kein Token der Registry ghcr.io" ohne "$out$err" "$TOKEN_G"; sei "kein Token von docker.io" ohne "$out$err" "$TOKEN_D"; }
f_p10_nicht_im_log_image() {
  local s
  fix 0.1.0; umg lauf
  for s in eingaben binaries image veroeffentlichen vergleichen smoke latest; do lauf "$s"; kein_geheimnis; done
  bild; mkdir -p "$(REG ghcr.io)"; printf '{"anders":true}' > "$(REG ghcr.io)/0.1.0"; lauf veroeffentlichen; sei "rot" test "$code" -eq 1; kein_geheimnis
  bild; bis veroeffentlichen; tamper ghcr.io; lauf vergleichen; sei "rot" test "$code" -eq 1; kein_geheimnis
  bild; lauf veroeffentlichen GP_COPY_ROT=docker.io; sei "rot" test "$code" -eq 1; kein_geheimnis
  bild; lauf veroeffentlichen GP_LOGIN_ROT=docker.io; sei "rot" test "$code" -eq 1; kein_geheimnis
  fix 0.1.0; umg lauf -RELEASE_DOCKERHUB_USER; lauf eingaben; sei "rot" test "$code" -eq 2; kein_geheimnis
}
f_p10_stdin() {
  local r l
  latest_von 1.1.0; lauf latest; gruen
  for r in ghcr.io docker.io; do
    sei "$r: der Token kommt über stdin" test "$(cat "$w/log/stdin.$r")" = "$([ "$r" = ghcr.io ] && echo "$TOKEN_G" || echo "$TOKEN_D")"
    sei "$r: Benutzer" test "$(cat "$w/log/user.$r")" = "$([ "$r" = ghcr.io ] && echo ghuser || echo dockeruser)"
  done
  l="$(grep '^skopeo login ' "$w/log/skopeo")"
  sei "vier Anmeldungen (zwei Schritte, zwei Registries)" test "$(wc -l <<< "$l")" -eq 4
  sei "jede mit --password-stdin" test "$(grep -c ' --password-stdin ' <<< "$l")" -eq 4
  sei "kein Token in einem Argument" ohne "$l" "$TOKEN_G"; sei "kein Token von docker.io in einem Argument" ohne "$l" "$TOKEN_D"
}
# Punkt 11
f_p11_wiederanlauf_ghcr_vorhanden() {
  bild; lauf veroeffentlichen GP_COPY_ROT=docker.io; sei "erster Lauf rot" test "$code" -eq 1
  lauf veroeffentlichen; gruen
  sei "ghcr.io nur einmal kopiert" test "$(kopien ghcr.io)" -eq 1; sei "docker.io kopiert" test "$(kopien docker.io)" -eq 2
  sei "beide tragen D" test "sha256:$(sha256sum "$(REG docker.io)/0.1.0" | cut -d' ' -f1)" = "$D"
}
f_p11_wiederanlauf_latest() { latest_von 1.1.0 v1.0.0; lauf latest GP_COPY_ROT=docker.io; sei "erster Lauf rot" test "$code" -eq 1; lauf latest; gruen; gesetzt; }
f_p11_reihenfolge_gruende_b() {
  local s
  fix 0.1.0; umg lauf; bis eingaben
  for s in image veroeffentlichen vergleichen smoke latest; do lauf "$s"; grund 2 "$s" vorbedingung; done
  sei "skopeo, docker und der Smoke nie gerufen" test ! -s "$w/log/skopeo" -a ! -s "$w/log/docker" -a ! -s "$w/log/smoke"
}
f_p11_alte_marke_b() {
  bis_vergl; lauf smoke; umg lauf RELEASE_REF=refs/tags/v01.2.3; lauf eingaben; sei "eingaben rot" test "$code" -eq 2
  : > "$w/runner/release/binaries.ok"; umg lauf; ohne_vergl=1; lauf github-release; grund 2 github-release vorbedingung; sei "gh nie gerufen" nie gh
  lauf smoke; grund 2 smoke vorbedingung
  lauf veroeffentlichen; grund 2 veroeffentlichen vorbedingung
  bis_vergl; tamper ghcr.io; lauf vergleichen; sei "vergleichen rot" test "$code" -eq 1; sei "keine Marke eines früheren Laufs" test ! -e "$w/runner/release/vergleichen.ok"
  bild; lauf image GP_SUMME=amd64; sei "image rot" test "$code" -eq 1; lauf veroeffentlichen; grund 2 veroeffentlichen vorbedingung
}

# --- Fassung F: Scan des Images vor dem Push (Punkt 16) ---------------------------------------------------------------------------------------------------
runs() { grep '^docker run ' "$w/log/docker" 2> /dev/null || true; }
ZD() { echo "$w/runner/release"; }
# sscan [NAME=Wert]… — Baum und Lauf bis image, dann scan mit den Zusätzen (GP_TRIVY, GP_DB, GP_ARCH: Liste `<arch>=<Modus>`).
sscan() { ibild; lauf scan "$@"; }
in_beiden() { local z n=0 alle=1; while IFS= read -r z; do n=$((n + 1)); [[ $z == *"$2"* ]] || alle=0; done < <(runs); sei "Trivy gerufen" test "$n" -ge 1; sei "$1" test "$alle" -eq 1; }
in_keinem() { sei "Trivy gerufen" test "$(runs | wc -l)" -ge 1; sei "ohne $1" ohne "$(runs)" "$1"; }
f_p16_trockenlauf() {
  fix 0.3.0-dev; umg trocken; bis eingaben binaries image scan
  sei "skopeo nie gerufen" nie skopeo; sei "gh nie gerufen" nie gh; sei "kein login" ohne "$(aufrufe docker)" "login"
  sei "Trivy zweimal gerufen" test "$(runs | wc -l)" -eq 2; sei "Marke scan.ok" test -e "$(ZD)/scan.ok"
}
f_p16_gruen() {
  sscan; gruen
  sei "Marke scan.ok" test -e "$(ZD)/scan.ok"; sei "Summenzeile" enthaelt "$out" "release: scan: 0 Befund(e) über der Schwelle, 0 ohne Behebung"
  sei "skopeo nie gerufen" nie skopeo; sei "gh nie gerufen" nie gh
}
# Aufruf, je Element
f_p16_e1() { sscan; gruen; in_beiden "docker run --rm" "docker run --rm "; }
f_p16_e2() { sscan; gruen; in_beiden "Mount :ro" ":/layout.tar:ro "; }
f_p16_e2_quelle() {
  local z
  sscan; gruen; in_beiden "Quelle des Mounts ist das Layout des Schrittes image" " -v $(IMG)/layout.tar:"
  lauf veroeffentlichen; gruen; z="$(grep '^skopeo copy ' "$w/log/skopeo" | head -n 1)"
  sei "der Push liest dieselbe Datei" enthaelt "$z" " oci-archive:$(IMG)/layout.tar "
}
f_p16_e3() { sscan; gruen; in_beiden "das Image des Pins aus release.sh" " $(grep '^TRIVY_IMAGE=' "$releasesh" | sed "s/^TRIVY_IMAGE='\([^']*\)'.*/\1/") "; }
f_p16_e4() { sscan; gruen; in_beiden "Unterbefehl image" " image --input "; }
f_p16_e4_input() { sscan; gruen; in_beiden "--input /layout.tar" " --input /layout.tar "; }
f_p16_e5() { sscan; gruen; in_beiden "--format json" " --format json "; }
f_p16_e6() { sscan; gruen; in_beiden "--exit-code 0" " --exit-code 0 "; }
f_p16_e7() { sscan; gruen; in_beiden "--severity HIGH,CRITICAL" " --severity HIGH,CRITICAL "; }
f_p16_e8() { sscan; gruen; in_beiden "--timeout 10m" " --timeout 10m "; }
f_p16_e9_amd64() { sscan; gruen; sei "genau ein Lauf für linux/amd64" test "$(runs | grep -cE ' --platform linux/amd64( |$)')" -eq 1; }
f_p16_e9_arm64() { sscan; gruen; sei "genau ein Lauf für linux/arm64" test "$(runs | grep -cE ' --platform linux/arm64( |$)')" -eq 1; }
f_p16_e9_genau_zwei() { sscan; gruen; sei "genau zwei Läufe" test "$(runs | wc -l)" -eq 2; }
for f in skip-db-update skip-files skip-dirs ignorefile ignore-policy vex; do eval "f_p16_e10_${f//-/_}() { sscan; gruen; in_keinem \"--$f\"; }"; done
f_p16_e10_mount() {
  local z
  sscan; gruen
  while IFS= read -r z; do sei "genau zwei Mounts je Lauf (Layout und Cache)" test "$(grep -o ' -v ' <<< "$z" | wc -l)" -eq 2; done < <(runs)
}
f_p16_user() { sscan; gruen; in_beiden "--user <uid>:<gid> des Aufrufers" " --user $(id -u):$(id -g) "; }
f_p16_cache_dir() { sscan; gruen; in_beiden "--cache-dir /cache" " --cache-dir /cache"; }
f_p16_cache_kein_zwischenspeicher() { sscan; gruen; in_beiden "der Cache liegt im Temp-Verzeichnis des Laufs" " -v $w/tmp/"; }
f_p16_pin_form() { sei "Form: aquasec/trivy@sha256:<64 Hexziffern> und Tag-Kommentar" test -n "$(grep -E "^TRIVY_IMAGE='aquasec/trivy@sha256:[0-9a-f]{64}' # [0-9]+\.[0-9]+\.[0-9]+\$" "$releasesh" || true)"; }
f_p16_pin_platzhalter() { sei "der Digest ist kein Platzhalter aus Nullen" test -z "$(grep -E "^TRIVY_IMAGE=.*@sha256:0{64}'" "$releasesh" || true)"; }
# Auswertung
# befund <Modi> <Meldung…> — scan endet mit Exit 1 und nennt jede Meldung (Text nach `release: scan: `); keine Marke.
befund() {
  local modi="$1" z
  shift; ibild; lauf scan "GP_TRIVY=$modi"
  sei "Exit 1" test "$code" -eq 1; sei "keine Marke" test ! -e "$(ZD)/scan.ok"
  for z in "$@"; do sei "Meldung: $z" enthaelt "$err" "release: scan: $z"; done
}
f_p16_critical_mit_fix() { befund amd64=critical-fix "schwachstelle: linux/amd64 CVE-1 pkg-c 1.0 -> 1.2.3 CRITICAL" "1 Befund(e) über der Schwelle, 0 ohne Behebung"; }
f_p16_high_mit_fix() { befund amd64=high-fix "schwachstelle: linux/amd64 CVE-1 pkg-h 1.0 -> 1.2.3 HIGH" "1 Befund(e) über der Schwelle, 0 ohne Behebung"; }
# darunter <Modus> <ohne Behebung> — grün trotz eines Eintrags; die Summe zählt ihn, die Meldung nennt keinen Befund.
darunter() { ibild; lauf scan "GP_TRIVY=amd64=$1"; gruen; sei "Marke" test -e "$(ZD)/scan.ok"; sei "Summe" enthaelt "$out" "0 Befund(e) über der Schwelle, $2 ohne Behebung"; sei "kein Befund genannt" ohne "$err" "schwachstelle"; }
f_p16_medium_mit_fix() { darunter medium-fix 0; }
f_p16_low_mit_fix() { darunter low-fix 0; }
f_p16_unknown_mit_fix() { darunter unknown-fix 0; }
f_p16_high_ohne_fix() { darunter high-ohne-fix 1; }
f_p16_high_ohne_fix_leer() { darunter high-ohne-fix-leer 1; }
f_p16_high_ohne_fix_null() { darunter high-ohne-fix-null 1; }
f_p16_critical_ohne_fix() { darunter critical-ohne-fix 1; }
f_p16_alle_befunde() {
  befund "amd64=befunde-a arm64=befunde-b" "schwachstelle: linux/amd64 CVE-A1 pkg-a1 1.0 -> 2.0 CRITICAL" "schwachstelle: linux/amd64 CVE-A2 pkg-a2 1.0 -> 3.0 HIGH" \
    "schwachstelle: linux/arm64 CVE-B1 pkg-b1 1.0 -> 4.0 HIGH" "3 Befund(e) über der Schwelle, 1 ohne Behebung"
  sei "MEDIUM nicht genannt" ohne "$err" CVE-B2
}
f_p16_plattform_genannt() { befund "amd64=befunde-a arm64=befunde-b" "schwachstelle: linux/amd64 CVE-A1" "schwachstelle: linux/arm64 CVE-B1"; }
f_p16_befund_nur_amd64() { befund amd64=high-fix "schwachstelle: linux/amd64 CVE-1"; sei "arm64 nicht genannt" ohne "$err" "linux/arm64"; }
f_p16_befund_nur_arm64() { befund arm64=high-fix "schwachstelle: linux/arm64 CVE-1"; sei "amd64 nicht genannt" ohne "$err" "linux/amd64"; }
f_p16_befund_im_zweiten_ergebnis() { befund amd64=zweites-ergebnis "schwachstelle: linux/amd64 CVE-Z pkg-z 1.0 -> 7.0 HIGH"; }
f_p16_summe() { befund amd64=summe "2 Befund(e) über der Schwelle, 2 ohne Behebung"; }
# Gescheitert
scanfehler() { sei "Exit 2" test "$code" -eq 2; sei "Meldung release: scan: $1:" enthaelt "$err" "release: scan: $1:"; sei "keine Marke" test ! -e "$(ZD)/scan.ok"; }
f_p16_trivy_exit() { sscan GP_TRIVY=amd64=exit1; scanfehler scan-gescheitert; }
f_p16_trivy_exit_arm64() { sscan GP_TRIVY=arm64=exit1; scanfehler scan-gescheitert; }
f_p16_trivy_exit_125() { sscan GP_TRIVY=amd64=exit125; scanfehler scan-gescheitert; }
f_p16_kein_json() { sscan GP_TRIVY=amd64=kein-json; scanfehler pruefmenge; }
f_p16_kein_json_arm64() { sscan GP_TRIVY=arm64=kein-json; scanfehler pruefmenge; }
f_p16_kein_json_leer() { sscan GP_TRIVY=amd64=leer; scanfehler pruefmenge; }
f_p16_leeres_ergebnis() { sscan GP_TRIVY=amd64=ergebnis-leer; scanfehler pruefmenge; }
f_p16_leeres_ergebnis_arm64() { sscan GP_TRIVY=arm64=ergebnis-leer; scanfehler pruefmenge; }
f_p16_leeres_ergebnis_fehlt() { sscan GP_TRIVY=amd64=ergebnis-fehlt; scanfehler pruefmenge; }
f_p16_leeres_ergebnis_null() { sscan GP_TRIVY=amd64=ergebnis-null; scanfehler pruefmenge; }
f_p16_db_alt() { sscan GP_DB=amd64=h73; scanfehler datenbank; sei "das Alter steht in der Meldung" enthaelt "$err" "73 Stunden"; }
f_p16_db_alt_arm64() { sscan GP_DB=arm64=h73; scanfehler datenbank; sei "das Alter steht in der Meldung" enthaelt "$err" "73 Stunden"; }
f_p16_db_frisch() { sscan "GP_DB=amd64=h71 arm64=h71"; gruen; sei "Marke" test -e "$(ZD)/scan.ok"; }
db_unlesbar() { sscan "GP_DB=$1"; scanfehler datenbank; sei "kein Alter behauptet" ohne "$err" "Stunden alt"; }
f_p16_db_unlesbar() { db_unlesbar amd64=fehlt; }
f_p16_db_unlesbar_arm64() { db_unlesbar arm64=fehlt; }
f_p16_db_unlesbar_kein_json() { db_unlesbar amd64=kein-json; }
f_p16_db_unlesbar_ohne_feld() { db_unlesbar amd64=ohne-feld; }
f_p16_db_unlesbar_zahl() { db_unlesbar amd64=zahl; }
f_p16_db_unlesbar_leer() { db_unlesbar amd64=leer; }
f_p16_db_unlesbar_jetzt() { db_unlesbar amd64=jetzt; }
f_p16_db_unlesbar_kein_datum() { db_unlesbar amd64=kein-datum; }
f_p16_plattform_abweichend() { sscan GP_ARCH=amd64=arm64; scanfehler plattform; }
f_p16_plattform_abweichend_arm64() { sscan GP_ARCH=arm64=amd64; scanfehler plattform; }
f_p16_plattform_fehlt() { sscan GP_ARCH=amd64=-; scanfehler plattform; }
f_p16_ohne_image() { fix 0.1.0; umg lauf; bis eingaben binaries; lauf scan; grund 2 scan vorbedingung; sei "Trivy nie gerufen" test -z "$(runs)"; }
f_p16_ohne_image_mit_layout() { bild_mit GP_SUMME=amd64; lauf scan; grund 2 scan vorbedingung; sei "das Layout liegt da" test -f "$(IMG)/layout.tar"; sei "Trivy nie gerufen" test -z "$(runs)"; }
f_p16_ohne_layout() { ibild; rm "$(IMG)/layout.tar"; lauf scan; grund 2 scan vorbedingung; sei "Trivy nie gerufen" test -z "$(runs)"; }
# Zustand
f_p16_scan_ok_marke() {
  sscan; gruen; sei "Marke nach grünem Scan" test -e "$(ZD)/scan.ok"
  sscan GP_TRIVY=amd64=critical-fix; sei "keine Marke nach Befund" test ! -e "$(ZD)/scan.ok"
  sscan GP_TRIVY=amd64=exit125; sei "keine Marke nach gescheitertem Scan" test ! -e "$(ZD)/scan.ok"
}
f_p16_alte_marke() { sscan; gruen; lauf scan GP_TRIVY=amd64=critical-fix; sei "Exit 1" test "$code" -eq 1; sei "die Marke des früheren Laufs ist weg" test ! -e "$(ZD)/scan.ok"; }
f_p16_eingaben_loescht_marke() { sscan; gruen; lauf eingaben; gruen; sei "die Marke ist weg" test ! -e "$(ZD)/scan.ok"; }
f_p16_image_loescht_marke() { sscan; gruen; lauf image; gruen; sei "die Marke ist weg" test ! -e "$(ZD)/scan.ok"; }
f_p16_veroeffentlichen_ohne_scan() { ibild; lauf veroeffentlichen; grund 2 veroeffentlichen vorbedingung; sei "scan genannt" enthaelt "$err" "scan"; sei "skopeo nie gerufen" nie skopeo; }
f_p16_ueber_schwelle_keine_registry() {
  sscan GP_TRIVY=amd64=critical-fix; sei "Scan rot" test "$code" -eq 1
  lauf veroeffentlichen; grund 2 veroeffentlichen vorbedingung; sei "skopeo nie gerufen" nie skopeo; sei "keine Registry beschrieben" test -z "$(ls -A "$w/registry")"
}
f_p16_d_unveraendert() {
  local a
  ibild; a="$(sha256sum "$(IMG)/layout.tar" "$(IMG)/digest")"
  lauf scan; gruen; sei "Layout und Digest unverändert" test "$(sha256sum "$(IMG)/layout.tar" "$(IMG)/digest")" = "$a"
}
f_p13_veroeffentlichen_ohne_image_bei_scan_ok() { bild; rm -rf "$(IMG)"; lauf veroeffentlichen; grund 2 veroeffentlichen vorbedingung; sei "skopeo nie gerufen" nie skopeo; }
f_p11_image_rot_loescht_digest_b() {
  bild; lauf image GP_SUMME=amd64; sei "image rot" test "$code" -eq 1
  : > "$(ZD)/scan.ok"; lauf veroeffentlichen; grund 2 veroeffentlichen vorbedingung; sei "skopeo nie gerufen" nie skopeo
}
f_p11_alte_marke_image_b() {
  bild; umg lauf RELEASE_REF=refs/tags/v01.2.3; lauf eingaben; sei "eingaben rot" test "$code" -eq 2
  : > "$(ZD)/scan.ok"; umg lauf; lauf veroeffentlichen; grund 2 veroeffentlichen vorbedingung; sei "skopeo nie gerufen" nie skopeo
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
FAELLE_P1=(p1-gleich p1-text-gleich p1-leerzeile p1-kommentarzeile p1-kommentarzeile-eingerueckt
  p1-pin-kurz p1-pin-gross p1-pin-ohne-tag p1-pin-tag-leer p1-pin-andere-action)
for ((n = 1; n <= N; n++)); do FAELLE_P1+=("p1-fehlt-$n" "p1-geaendert-$n"); done
for k in "${ZUSATZ_NAMEN[@]}"; do FAELLE_P1+=("p1-zusatz-$k"); done
for j in "${BEDING_ID[@]}"; do FAELLE_P1+=("p1-bedingung-$j"); done
for tp in "${TAUSCHE[@]}"; do FAELLE_P1+=("p1-reihenfolge-${tp%% *}"); done
FAELLE_SH=(p2-trockenlauf-ohne-aussenwirkung p2-aussenwirkung-ohne-push p2-ref-kein-tag
  p3-tag-gueltig p3-laenge-100 p3-ohne-v p3-zwei-teile p3-vier-teile p3-fuehrende-null p3-fuehrende-null-minor p3-fuehrende-null-patch p3-vorabversion-fuehrende-null
  p3-leere-vorabversion p3-leerer-bezeichner p3-ungueltiges-zeichen p3-build-metadaten p3-zu-lang p3-ereignis-unbekannt p3-ereignis-fehlt
  p3-version-datei-form-fehlt p3-version-datei-form-zwei-zeilen p3-version-datei-form-v p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-leerzeichen-hinten
  p3-version-datei-form-crlf p3-version-datei-form-kein-lf p3-tag-ungleich-datei p3-tag-ungleich-gross p3-tag-ungleich-tag-laenger p3-tag-ungleich-tag-kuerzer
  p3-notes-fehlen p3-notes-zelle-leer p3-notes-pipe p3-notes-abschnitt p3-notes-ueberschrift-endet p3-notes-abschnitt-ende p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts
  p3-standard-zweig-aus-umgebung p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-zelle-leer p3-trockenlauf-notes-vorhanden
  p3-trockenlauf-version-datei-form p3-schluesselwort-ereignis p3-schluesselwort-tag-form p3-schluesselwort-version-datei p3-schluesselwort-tag-ungleich p3-schluesselwort-notes p3-schluesselwort-standard-zweig p3-reihenfolge-1-2 p3-reihenfolge-2-3 p3-reihenfolge-3-4 p3-reihenfolge-4-5 p3-reihenfolge-5-6
  p4-aufruf p4-summe-falsch p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert
  p9-anhaenge p9-aufruf p9-stabil-ohne-prerelease p9-vorabversion-prerelease p9-notes-aus-zelle p9-notes-zeile-der-version p9-notes-aus-abschnitt
  p9-release-existiert-gleich p9-release-existiert-abweichend p9-view-anderer-fehler p9-repo-argument p9-gh-token-fehlt p13-binaries-ohne-eingaben p13-release-ohne-binaries p10-nicht-im-log
  p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p11-reihenfolge-gruende p11-alte-marke)
FAELLE_B=(p2-trockenlauf-image p2-aussenwirkung-ohne-push-b p2-ref-kein-tag-b p3-secret-fehlt-user p3-secret-fehlt-token p3-secret-leer-user p3-secret-leer-token p3-secret-fehlen-beide p3-schluesselwort-geheimnis p3-reihenfolge-6-7
  p5-aufruf p5-builder-pin-form p5-builder-platzhalter p5-epoch-committer p5-ohne-binaries p5-zwei-manifeste p5-plattform-fehlt-amd64 p5-plattform-fehlt-arm64 p5-layout-form p5-summe-amd64 p5-summe-arm64 p5-summe-eintrag-fehlt p5-digest p5-bau-scheitert
  p6-reihenfolge p6-all-und-digests p6-tag-vorhanden-gleich p6-tag-vorhanden-anders p6-tag-fehlt-manifest-unknown p6-tag-fehlt-name-unknown p6-tag-fehlt-not-found p6-inspect-anderer-fehler p6-inspect-meldung-im-log p6-zweite-registry-faellt p6-tag-satz p6-anmeldung-scheitert
  p7-gleich p7-rohbytes p7-digest-ghcr p7-digest-docker p7-digest-lokal p7-anonym-ghcr p7-anonym-docker p7-eine-plattform p7-tag-fehlt p7-smoke-aufruf p7-smoke-ghcr p7-smoke-docker p7-postgres-fehlt p7-smoke-ohne-vergleichen p7-smoke-ok-marke
  p8-stabil-hoechste p8-vorabversion p8-aelterer-zweig p8-erster-stabiler p8-vorabversion-hoeher-zaehlt-nicht p8-semver-numerisch p8-nach-rotem-vergleich p8-ohne-smoke p8-fremder-tag p8-tag-satz p8-zweite-registry-faellt
  p9-nicht-vor-registries p10-nicht-im-log-image p10-stdin p11-wiederanlauf-ghcr-vorhanden p11-wiederanlauf-latest p11-reihenfolge-gruende-b p11-alte-marke-b
  p13-veroeffentlichen-ohne-image-bei-scan-ok p11-alte-marke-image-b p11-image-rot-loescht-digest-b)
FAELLE_F=(p16-trockenlauf p16-gruen p16-e1 p16-e2 p16-e2-quelle p16-e3 p16-e4 p16-e4-input p16-e5 p16-e6 p16-e7 p16-e8 p16-e9-amd64 p16-e9-arm64 p16-e9-genau-zwei
  p16-e10-skip-db-update p16-e10-skip-files p16-e10-skip-dirs p16-e10-ignorefile p16-e10-ignore-policy p16-e10-vex p16-e10-mount p16-user p16-cache-dir p16-cache-kein-zwischenspeicher p16-pin-form p16-pin-platzhalter
  p16-critical-mit-fix p16-high-mit-fix p16-medium-mit-fix p16-low-mit-fix p16-unknown-mit-fix p16-high-ohne-fix p16-high-ohne-fix-leer p16-high-ohne-fix-null p16-critical-ohne-fix
  p16-alle-befunde p16-plattform-genannt p16-befund-nur-amd64 p16-befund-nur-arm64 p16-befund-im-zweiten-ergebnis p16-summe
  p16-trivy-exit p16-trivy-exit-arm64 p16-trivy-exit-125 p16-kein-json p16-kein-json-arm64 p16-kein-json-leer p16-leeres-ergebnis p16-leeres-ergebnis-arm64 p16-leeres-ergebnis-fehlt p16-leeres-ergebnis-null
  p16-db-alt p16-db-alt-arm64 p16-db-frisch p16-db-unlesbar p16-db-unlesbar-arm64 p16-db-unlesbar-kein-json p16-db-unlesbar-ohne-feld p16-db-unlesbar-zahl p16-db-unlesbar-leer p16-db-unlesbar-jetzt p16-db-unlesbar-kein-datum
  p16-plattform-abweichend p16-plattform-abweichend-arm64 p16-plattform-fehlt p16-ohne-image p16-ohne-image-mit-layout p16-ohne-layout
  p16-scan-ok-marke p16-alte-marke p16-eingaben-loescht-marke p16-image-loescht-marke p16-veroeffentlichen-ohne-scan p16-ueber-schwelle-keine-registry p16-d-unveraendert)
FAELLE_P13=(p13-ende p13-sigint p13-sigterm)
FAELLE_GP13=(p13-gegenprobe-sigint p13-gegenprobe-sigterm p13-gegenprobe-ende)

# faelle <gruppe> — schreibt die Namen der roten Fälle der Gruppe (p1, sh, b, f, p13, gp13, ok = alle), einen je Zeile. FALL_ZUSATZ nennt einen weiteren Namen (Probe: ein Name ohne
# Funktion ist ein roter Fall).
faelle() {
  local n
  case "$1" in
    p1) set -- "${FAELLE_P1[@]}" ;;
    sh) set -- "${FAELLE_SH[@]}" ;;
    b) set -- "${FAELLE_B[@]}" ;;
    f) set -- "${FAELLE_F[@]}" ;;
    p13) set -- "${FAELLE_P13[@]}" ;;
    gp13) if [ -n "${GP_INNEN-}" ]; then set --; else set -- "${FAELLE_GP13[@]}"; fi ;;
    ok) set -- "${FAELLE_P1[@]}" "${FAELLE_SH[@]}" "${FAELLE_B[@]}" "${FAELLE_F[@]}" "${FAELLE_P13[@]}"; [ -n "${GP_INNEN-}" ] || set -- "$@" "${FAELLE_GP13[@]}" ;;
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
  for j in "${!BEDING[@]}"; do if [ "$i" -eq "${BEDING[j]}" ]; then rote+=" p1-bedingung-${BEDING_ID[j]}"; fi; done
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
soll_umgestellt() { local k; printf '%s\n' "${SOLL[@]:0:STEP_AB[0]}"; for k in "$@"; do printf '%s\n' "${SOLL[@]:STEP_AB[k]:STEP_AB[k+1]-STEP_AB[k]}"; done; }
tauschbar='src="$(declare -f pruefe_release)"; eval "${src/pruefe_release/pruefe_release_orig}"; pruefe_release() { pruefe_release_orig "$1" && return 0; local -a s=("${SOLL[@]}"); local r; SOLL=("${%s[@]}"); pruefe_release_orig "$1"; r=$?; SOLL=("${s[@]}"); return $r; }'
for tp in "${TAUSCHE[@]}"; do
  read -r tn tk tl <<< "$tp"
  mapfile -t "PERM_${tn//-/_}" < <(soll_umgestellt $(tausch "$tk" "$tl"))
  # shellcheck disable=SC2059
  mut1 "m-p1-reihenfolge-$tn-tauschbar" "p1-reihenfolge-$tn" "$(printf "$tauschbar" "PERM_${tn//-/_}")"
done

# --- Mutationen von release.sh (Punkte 2 bis 11) ----------------------------------------------------------------------------------------------------
# Punkt 2
mutr m-p2-trockenlauf-ruft-skopeo "p2-trockenlauf-ohne-aussenwirkung" sh '  : > "$zdir/binaries.ok"' $'  [ "$trockenlauf" -eq 0 ] || skopeo inspect docker://x\n  : > "$zdir/binaries.ok"'
mutr m-p2-github-release-ohne-ereignispruefung "p2-aussenwirkung-ohne-push p2-ref-kein-tag" sh $'  aussenwirkung_erlaubt\n  zustand\n' $'  zustand\n'
mutr m-p2-ereignis-nicht-geprueft "p2-aussenwirkung-ohne-push p3-ereignis-fehlt p3-ereignis-unbekannt p3-schluesselwort-ereignis" sh \
  'Z:  [ "${RELEASE_EVENT-}" = push ] || eingabe ereignis "Ereignis '"'"'${RELEASE_EVENT-}'"'"' ist kein Push eines Tags" # P2ereignis' ""
mutr m-p2-ref-nicht-geprueft "p2-ref-kein-tag p3-ohne-v p3-reihenfolge-1-2" sh \
  'Z:  [[ ${RELEASE_REF-} == refs/tags/v* ]] || eingabe ereignis "Ref '"'"'${RELEASE_REF-}'"'"' ist kein Tag v…" # P2ref' ""
mutr m-p2-ref-ohne-v-praefix "p3-ohne-v p3-reihenfolge-1-2" sh '[[ ${RELEASE_REF-} == refs/tags/v* ]]' '[[ ${RELEASE_REF-} == refs/tags/* ]]'
mutr m-p2-trockenlauf-wie-lauf "p2-trockenlauf-ohne-aussenwirkung p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-trockenlauf-version-datei-form p4-aufruf p3-notes-pipe" sh \
  'Z:  if [ "${RELEASE_EVENT-}" = workflow_dispatch ]; then return 0; fi' ""
# Punkt 3
TAGFORM="p3-zwei-teile p3-vier-teile p3-fuehrende-null p3-fuehrende-null-minor p3-fuehrende-null-patch p3-vorabversion-fuehrende-null p3-leere-vorabversion p3-leerer-bezeichner p3-ungueltiges-zeichen p3-build-metadaten p3-zu-lang p3-reihenfolge-2-3 p3-schluesselwort-tag-form"
VDATEI="p3-version-datei-form-fehlt p3-version-datei-form-zwei-zeilen p3-version-datei-form-v p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-leerzeichen-hinten p3-version-datei-form-crlf p3-version-datei-form-kein-lf p3-schluesselwort-version-datei"
UNGLEICH="p3-tag-ungleich-datei p3-tag-ungleich-gross p3-tag-ungleich-tag-laenger p3-tag-ungleich-tag-kuerzer p3-schluesselwort-tag-ungleich"
NOTES="p3-notes-fehlen p3-notes-zelle-leer p3-notes-pipe p3-notes-abschnitt p3-notes-ueberschrift-endet p3-notes-abschnitt-ende p3-schluesselwort-notes"
# GRUEN: die Fälle, die eine grüne Prüfung der Eingaben brauchen oder deren Prüfung hinter (4) liegt; sie werden rot, wenn (4) nichts Gleiches mehr vergleicht.
# BIN: die Fälle, die `binaries` grün fahren und an der Dateiliste des Baus hängen; sie werden rot, wenn der Aufruf von make nicht die Version oder ein Verzeichnis außerhalb des Baums nennt.
BIN="p9-gh-token-fehlt p9-repo-argument p10-nicht-im-log p11-alte-marke p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p4-aufruf p9-anhaenge p9-aufruf p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease"
GRUEN="p3-notes-pipe p13-release-ohne-binaries p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-schluesselwort-standard-zweig p3-vorabversion-fuehrende-null p10-nicht-im-log p11-alte-marke p11-reihenfolge-gruende p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-reihenfolge-5-6 p3-standard-zweig-aus-umgebung p3-tag-gueltig p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p4-aufruf p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert p4-summe-falsch p9-anhaenge p9-aufruf p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease"
# Je Prüfung (1) bis (6) eine Mutation, die den Aufruf entfernt.
mutr m-p3-1-ereignis-entfernt "p2-ref-kein-tag p3-ereignis-fehlt p3-ereignis-unbekannt p3-ohne-v p3-reihenfolge-1-2 p3-schluesselwort-ereignis" sh $'  pruefe_ereignis\n' ''
mutr m-p3-2-tag-form-entfernt "$TAGFORM" sh $'  pruefe_tag_form\n' ''
mutr m-p3-3-version-datei-entfernt "p10-nicht-im-log p11-alte-marke p11-reihenfolge-gruende p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p13-release-ohne-binaries p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-3-4 p3-reihenfolge-5-6 p3-schluesselwort-notes p3-schluesselwort-standard-zweig p3-schluesselwort-version-datei p3-standard-zweig-aus-umgebung p3-tag-gueltig p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-trockenlauf-version-datei-form p3-version-datei-form-crlf p3-version-datei-form-fehlt p3-version-datei-form-kein-lf p3-version-datei-form-leerzeichen-hinten p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-v p3-version-datei-form-zwei-zeilen p3-vorabversion-fuehrende-null p4-aufruf p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert p4-summe-falsch p9-anhaenge p9-aufruf p9-gh-token-fehlt p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-repo-argument p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease p3-notes-pipe" sh $'  pruefe_version_datei\n' ''
mutr m-p3-4-tag-gleich-entfernt "$UNGLEICH p3-reihenfolge-4-5" sh $'  pruefe_tag_gleich\n' ''
mutr m-p3-5-notes-entfernt "p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-5-6 p3-schluesselwort-notes p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-notes-pipe" sh $'  pruefe_notes\n' ''
mutr m-p3-6-standard-zweig-entfernt "p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-schluesselwort-standard-zweig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig" sh $'  pruefe_standard_zweig\n' ''
# Reihenfolge: je Paar benachbarter Prüfungen vertauscht.
mutr m-p3-reihenfolge-1-2-vertauscht "p3-reihenfolge-1-2 p3-ohne-v p2-ref-kein-tag" sh $'  pruefe_ereignis\n  pruefe_tag_form\n' $'  pruefe_tag_form\n  pruefe_ereignis\n'
mutr m-p3-reihenfolge-2-3-vertauscht "${TAGFORM/ p3-zu-lang/}" sh $'  pruefe_tag_form\n  pruefe_version_datei\n' $'  pruefe_version_datei\n  pruefe_tag_form\n'
mutr m-p3-reihenfolge-3-4-vertauscht "p10-nicht-im-log p11-alte-marke p11-reihenfolge-gruende p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p13-release-ohne-binaries p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-3-4 p3-reihenfolge-5-6 p3-schluesselwort-notes p3-schluesselwort-standard-zweig p3-schluesselwort-version-datei p3-standard-zweig-aus-umgebung p3-tag-gueltig p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-trockenlauf-version-datei-form p3-version-datei-form-crlf p3-version-datei-form-fehlt p3-version-datei-form-kein-lf p3-version-datei-form-leerzeichen-hinten p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-v p3-version-datei-form-zwei-zeilen p3-vorabversion-fuehrende-null p4-aufruf p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert p4-summe-falsch p9-anhaenge p9-aufruf p9-gh-token-fehlt p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-repo-argument p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease p3-notes-pipe" sh $'  pruefe_version_datei\n  pruefe_tag_gleich\n' $'  pruefe_tag_gleich\n  pruefe_version_datei\n'
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
mutr m-p3-gleich-ohne-v "p10-nicht-im-log p11-alte-marke p11-reihenfolge-gruende p11-wiederanlauf-release-teilweise p11-wiederanlauf-release-vollstaendig p13-release-ohne-binaries p2-aussenwirkung-ohne-push p2-ref-kein-tag p2-trockenlauf-ohne-aussenwirkung p3-laenge-100 p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-5-6 p3-schluesselwort-notes p3-schluesselwort-standard-zweig p3-standard-zweig-aus-umgebung p3-tag-gueltig p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-vorabversion-fuehrende-null p4-aufruf p4-datei-fehlt p4-datei-zu-viel p4-make-scheitert p4-summe-falsch p9-anhaenge p9-aufruf p9-gh-token-fehlt p9-notes-aus-abschnitt p9-notes-aus-zelle p9-notes-zeile-der-version p9-release-existiert-abweichend p9-release-existiert-gleich p9-repo-argument p9-stabil-ohne-prerelease p9-view-anderer-fehler p9-vorabversion-prerelease p3-notes-pipe" sh '"$tag" = "v$version"' '"$tag" = "$version"'
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
mutr m-p3-trockenlauf-tag-aus-ref "p2-trockenlauf-ohne-aussenwirkung p3-trockenlauf-gueltig p3-trockenlauf-im-standard-zweig p3-trockenlauf-nicht-im-standard-zweig p3-trockenlauf-notes-fehlen p3-trockenlauf-notes-vorhanden p3-trockenlauf-notes-zelle-leer p3-trockenlauf-version-datei-form p4-aufruf p3-notes-pipe" sh \
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
declare -A SW_ROT=([ereignis]="p2-aussenwirkung-ohne-push p2-ref-kein-tag p3-ereignis-fehlt p3-ereignis-unbekannt p3-ohne-v p3-reihenfolge-1-2 p3-schluesselwort-ereignis" [tag-form]="p3-build-metadaten p3-fuehrende-null p3-fuehrende-null-minor p3-fuehrende-null-patch p3-leere-vorabversion p3-leerer-bezeichner p3-reihenfolge-2-3 p3-schluesselwort-tag-form p3-ungueltiges-zeichen p3-vier-teile p3-vorabversion-fuehrende-null p3-zu-lang p3-zwei-teile" [version-datei]="p3-reihenfolge-3-4 p3-schluesselwort-version-datei p3-trockenlauf-version-datei-form p3-version-datei-form-crlf p3-version-datei-form-fehlt p3-version-datei-form-kein-lf p3-version-datei-form-leerzeichen-hinten p3-version-datei-form-leerzeichen-vorn p3-version-datei-form-v p3-version-datei-form-zwei-zeilen" [tag-ungleich]="p3-reihenfolge-4-5 p3-schluesselwort-tag-ungleich p3-tag-ungleich-datei p3-tag-ungleich-gross p3-tag-ungleich-tag-kuerzer p3-tag-ungleich-tag-laenger" [notes]="p3-notes-pipe p3-notes-abschnitt p3-notes-abschnitt-ende p3-notes-fehlen p3-notes-ueberschrift-endet p3-notes-zelle-leer p3-reihenfolge-5-6 p3-schluesselwort-notes p3-trockenlauf-notes-zelle-leer" [standard-zweig]="p3-nicht-im-standard-zweig p3-nicht-im-standard-zweig-vorwaerts p3-schluesselwort-standard-zweig")
for kw in ereignis tag-form version-datei tag-ungleich notes standard-zweig; do
  mutr "m-p3-schluesselwort-$kw-weicht-ab" "${SW_ROT[$kw]}" sh 'eingabe() { echo "release: $1: $2" >&2; exit 2; }' "eingabe() { local k=\"\$1\"; [ \"\$k\" != $kw ] || k=\"\${k}x\"; echo \"release: \$k: \$2\" >&2; exit 2; }"
done
# Punkt 9
mutr m-p9-repo-fehlt-view "p9-repo-argument" sh 'release view "$tag" --repo "$RELEASE_REPOSITORY" --json' 'release view "$tag" --json'
mutr m-p9-repo-fehlt-create "p9-repo-argument" sh 'release create "$tag" --repo "$RELEASE_REPOSITORY" --verify-tag' 'release create "$tag" --verify-tag'
mutr m-p9-repo-fehlt-download "p9-repo-argument" sh 'release download "$tag" --repo "$RELEASE_REPOSITORY" --pattern' 'release download "$tag" --pattern'
mutr m-p9-repo-fehlt-upload "p9-repo-argument" sh 'release upload "$tag" --repo "$RELEASE_REPOSITORY" "${fehlend[@]}"' 'release upload "$tag" "${fehlend[@]}"'
mutr m-p9-gh-ohne-token-pruefung "p9-gh-token-fehlt" sh 'L:# P9token' ''
mutr m-p9-view-muster-geweitet "p9-view-anderer-fehler" sh "grep -q 'release not found'" "grep -qi 'not found'"
# Anhänge: je Element von erwartete_namen eine Mutation
for el in darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 windows_amd64 windows_arm64; do
  r="$BIN"; if [ "$el" = linux_arm64 ]; then r+=" p4-datei-fehlt"; fi
  mutr "m-p9-name-fehlt-$el" "$r" sh '      echo "pgwire-recorder_${1}_${os}_${a}$e"' "      [ \"\${os}_\${a}\" = $el ] || echo \"pgwire-recorder_\${1}_\${os}_\${a}\$e\""
done
mutr m-p9-name-fehlt-sha256sums "$BIN" sh 'L:  echo SHA256SUMS' ''
mutr m-p3-pipe-abgeschnitten-angenommen "p3-notes-pipe" sh 'L:# P3pipe' ''
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
# --- Mutationen von release.sh: Fassung B (Gruppe b) -------------------------------------------------------------------------------------------------
# Punkt 2: je Schritt mit Außenwirkung prüft das Ereignis; der Trockenlauf ruft weder skopeo noch eine Anmeldung
for st in veroeffentlichen vergleichen smoke latest; do
  mutr "m-p2-$st-ohne-ereignispruefung" "p2-aussenwirkung-ohne-push-b p2-ref-kein-tag-b" b "L:# P2$st" ''
done
mutr m-p2-trockenlauf-image-ruft-skopeo "p2-trockenlauf-image" b '  echo "release: image: Index-Digest $idx"' $'  [ "$trockenlauf" -eq 0 ] || skopeo inspect docker://x\n  echo "release: image: Index-Digest $idx"'
mutr m-p2-trockenlauf-image-meldet-an "p2-trockenlauf-image" b '  echo "release: image: Index-Digest $idx"' $'  [ "$trockenlauf" -eq 0 ] || docker login x\n  echo "release: image: Index-Digest $idx"'
# Punkt 3 (7)
mutr m-p3-secret-user-entfernt "p10-nicht-im-log-image p3-secret-fehlen-beide p3-secret-fehlt-user p3-secret-leer-user" b 'L:# P3secret-user' ''
mutr m-p3-secret-token-entfernt "p3-schluesselwort-geheimnis p3-secret-fehlen-beide p3-secret-fehlt-token p3-secret-leer-token" b 'L:# P3secret-token' ''
mutr m-p3-secret-user-leer-gilt-als-gesetzt "p3-secret-leer-user" b '[ -n "${RELEASE_DOCKERHUB_USER-}" ] ||' '[ "${RELEASE_DOCKERHUB_USER+x}" ] ||'
mutr m-p3-secret-token-leer-gilt-als-gesetzt "p3-secret-leer-token" b '[ -n "${RELEASE_DOCKERHUB_TOKEN-}" ] ||' '[ "${RELEASE_DOCKERHUB_TOKEN+x}" ] ||'
mutr m-p3-7-geheimnisse-entfernt "p10-nicht-im-log-image p3-schluesselwort-geheimnis p3-secret-fehlen-beide p3-secret-fehlt-token p3-secret-fehlt-user p3-secret-leer-token p3-secret-leer-user" b $'  pruefe_geheimnisse\n' ''
mutr m-p3-reihenfolge-6-7-vertauscht "p3-reihenfolge-6-7" b $'  pruefe_standard_zweig\n  pruefe_geheimnisse\n' $'  pruefe_geheimnisse\n  pruefe_standard_zweig\n'
mutr m-p3-schluesselwort-geheimnis-weicht-ab "p3-schluesselwort-geheimnis p3-secret-fehlen-beide p3-secret-fehlt-token p3-secret-fehlt-user p3-secret-leer-token p3-secret-leer-user" b 'eingabe() { echo "release: $1: $2" >&2; exit 2; }' 'eingabe() { local k="$1"; [ "$k" != geheimnis ] || k="${k}x"; echo "release: $k: $2" >&2; exit 2; }'
mutr m-p3-secret-wert-in-meldung "p10-nicht-im-log-image p3-secret-fehlt-user p3-secret-leer-user" b 'eingabe geheimnis "gesetzt und nicht leer verlangt:$fehlt"' 'eingabe geheimnis "gesetzt und nicht leer verlangt:$fehlt (${RELEASE_DOCKERHUB_TOKEN-})"'
# Punkt 5: je Argument des Baus
mutr m-p5-plattformen-entfallen "p5-aufruf" b ' --platform linux/amd64,linux/arm64 --output "type=oci' ' --output "type=oci'
mutr m-p5-target-entfaellt "p5-aufruf" b ' --target runtime' ''
mutr m-p5-version-entfaellt "p5-aufruf" b ' --build-arg "VERSION=$v"' ''
mutr m-p5-source-date-epoch-entfaellt "p5-aufruf p5-epoch-committer" b ' --build-arg "SOURCE_DATE_EPOCH=$zeit"' ''
mutr m-p5-provenance-entfaellt "p5-aufruf" b ' --provenance=false' ''
mutr m-p5-sbom-entfaellt "p5-aufruf" b ' --sbom=false' ''
mutr m-p5-rewrite-timestamp-entfaellt "p5-aufruf" b ',rewrite-timestamp=true' ''
mutr m-p5-output-nicht-oci "p5-aufruf p5-epoch-committer" b 'type=oci,dest=' 'type=docker,dest='
mutr m-p5-treiber-nicht-docker-container "p5-aufruf" b '--driver docker-container' '--driver docker'
mutr m-p5-manifeste-entfernt "p5-zwei-manifeste p5-plattform-fehlt-amd64 p5-plattform-fehlt-arm64" b 'L:# P5manifeste' ''
mutr m-p5-zahl-entfernt "p5-zwei-manifeste p7-eine-plattform" b 'L:# P5zahl' ''
mutr m-p5-plattform-amd64-nicht-verlangt "p5-plattform-fehlt-amd64" b '[ "$(plattformen "$1")" = "linux/amd64 linux/arm64 " ]' '[[ "$(plattformen "$1")" == *linux/arm64* ]]'
mutr m-p5-plattform-arm64-nicht-verlangt "p5-plattform-fehlt-arm64" b '[ "$(plattformen "$1")" = "linux/amd64 linux/arm64 " ]' '[[ "$(plattformen "$1")" == *linux/amd64* ]]'
mutr m-p5-summe-amd64-nicht-geprueft "p11-alte-marke-b p5-summe-amd64" b '[ -n "$soll" ] && [ "$ist" = "$soll" ] ||' '[ -n "$soll" ] && [ "$ist" = "$soll" ] || [ "$arch" = amd64 ] ||'
mutr m-p5-summe-arm64-nicht-geprueft "p5-summe-arm64 p5-summe-eintrag-fehlt" b '[ -n "$soll" ] && [ "$ist" = "$soll" ] ||' '[ -n "$soll" ] && [ "$ist" = "$soll" ] || [ "$arch" = arm64 ] ||'
mutr m-p5-digest-aus-archiv "p10-stdin p11-alte-marke-b p11-wiederanlauf-ghcr-vorhanden p11-wiederanlauf-latest p5-digest p6-tag-vorhanden-gleich p6-zweite-registry-faellt p7-digest-docker p7-gleich p7-postgres-fehlt p7-rohbytes p7-smoke-aufruf p7-smoke-docker p7-smoke-ghcr p7-smoke-ok-marke p8-aelterer-zweig p8-erster-stabiler p8-fremder-tag p8-ohne-smoke p8-semver-numerisch p8-stabil-hoechste p8-tag-satz p8-vorabversion p8-vorabversion-hoeher-zaehlt-nicht p8-zweite-registry-faellt p9-nicht-vor-registries" b 'idx="sha256:$(sha256sum "$scratch/index" | cut -d'"'"' '"'"' -f1)" # P5digest' 'idx="sha256:$(sha256sum "$zdir/image/layout.tar" | cut -d'"'"' '"'"' -f1)" # P5digest'
mutr m-p5-builder-fehler-ignoriert "p5-bau-scheitert" b $'  [ "$rc" -eq 0 ] || rot "docker buildx create endet mit Exit $rc"\n' ''
mutr m-p5-bau-fehler-ignoriert "p5-bau-scheitert" b $'  [ "$rc" -eq 0 ] || rot "docker buildx build endet mit Exit $rc"\n' ''
# Punkt 6
mutr m-p6-reihenfolge-vertauscht "p6-reihenfolge p6-zweite-registry-faellt" b 'for reg in "${REGISTRIES[@]}"; do # P6reg' 'for reg in docker.io ghcr.io; do # P6reg'
mutr m-p6-all-entfernt "p6-all-und-digests" b 'fuehre skopeo copy --all --preserve-digests --authfile "$scratch/auth.json" "oci-archive:' 'fuehre skopeo copy --preserve-digests --authfile "$scratch/auth.json" "oci-archive:'
mutr m-p6-preserve-digests-entfernt "p6-all-und-digests" b 'fuehre skopeo copy --all --preserve-digests --authfile "$scratch/auth.json" "oci-archive:' 'fuehre skopeo copy --all --authfile "$scratch/auth.json" "oci-archive:'
mutr m-p6-anderer-digest-angenommen "p10-nicht-im-log-image p6-tag-vorhanden-anders" b 'L:# P6anders' ''
mutr m-p6-vorhandener-tag-ueberschrieben "p10-nicht-im-log-image p6-tag-vorhanden-anders" b \
  '[ "$have" = "$d" ] || rot "$reg trägt $v mit dem Digest $have statt $d; der Weg ist ein höherer Tag" # P6anders' '[ "$have" = "$d" ] || { rc=1; echo "manifest unknown" > "$scratch/skopeo.err"; } # P6anders' \
  '      continue # P6gleich' '      [ "$rc" -ne 0 ] || continue # P6gleich'
mutr m-p6-gleicher-tag-rot "p6-tag-vorhanden-gleich p11-wiederanlauf-ghcr-vorhanden" b '      continue # P6gleich' '      rot "gleich" # P6gleich'
mutr m-p6-gleicher-tag-kopiert "p6-tag-vorhanden-gleich p11-wiederanlauf-ghcr-vorhanden" b '      continue # P6gleich' '      : # P6gleich'
mutr m-p6-zusaetzlicher-tag "p11-wiederanlauf-ghcr-vorhanden p6-reihenfolge p6-tag-fehlt-manifest-unknown p6-tag-fehlt-name-unknown p6-tag-fehlt-not-found p6-tag-satz p6-tag-vorhanden-gleich p8-tag-satz" b '    echo "release: veroeffentlichen: $v nach $reg kopiert"' $'    fuehre skopeo copy --all --preserve-digests --authfile "$scratch/auth.json" "docker://$reg/$IMAGE_NAME@$d" "docker://$reg/$IMAGE_NAME:${v%.*}"\n    echo "release: veroeffentlichen: $v nach $reg kopiert"'
mutr m-p6-zustand-fehlt "p6-zweite-registry-faellt" b '; $(zustand_tag "$v")"' '"'
mutr m-p6-zweite-registry-ignoriert "p6-zweite-registry-faellt p10-nicht-im-log-image p11-wiederanlauf-ghcr-vorhanden" b '[ "$rc" -eq 0 ] || rot "skopeo copy nach $reg endet mit Exit $rc; $(zustand_tag "$v")"' ':'
mutr m-p6-anmeldung-fehler-ignoriert "p6-anmeldung-scheitert p10-nicht-im-log-image" b '|| rot "Anmeldung bei $1 scheitert" # P10stdin' '|| true # P10stdin'
# Punkt 7
mutr m-p7-ghcr-vergleich-entfernt "p10-nicht-im-log-image p11-alte-marke-b p7-digest-ghcr" b 'L:# P7ghcr' ''
mutr m-p7-docker-vergleich-entfernt "p7-digest-docker" b 'L:# P7docker' ''
mutr m-p7-vergleich-gegen-d-entfernt "p7-digest-ghcr p7-digest-lokal" b '  [ "$dg" = "$d" ] || rot "ghcr.io weicht ab: $dg statt $d" # P7ghcr' $'  d="$dg"\n  [ "$dg" = "$d" ] || rot "ghcr.io weicht ab: $dg statt $d" # P7ghcr'
mutr m-p7-anonym-mit-anmeldung "p6-inspect-meldung-im-log p7-anonym-docker p7-anonym-ghcr p7-gleich" b ' --no-creds "docker://$reg/$IMAGE_NAME:$v"' ' "docker://$reg/$IMAGE_NAME:$v"'
mutr m-p7-plattformen-nicht-geprueft "p7-eine-plattform" b 'L:# P7plattformen' ''
mutr m-p7-digest-aus-anzeige "p10-stdin p11-alte-marke-b p11-wiederanlauf-latest p7-digest-docker p7-gleich p7-postgres-fehlt p7-rohbytes p7-smoke-aufruf p7-smoke-docker p7-smoke-ghcr p7-smoke-ok-marke p8-aelterer-zweig p8-erster-stabiler p8-fremder-tag p8-ohne-smoke p8-semver-numerisch p8-stabil-hoechste p8-tag-satz p8-vorabversion p8-vorabversion-hoeher-zaehlt-nicht p8-zweite-registry-faellt p9-nicht-vor-registries" b '  dg="sha256:$(sha256sum "$scratch/roh-ghcr.io" | cut -d'"'"' '"'"' -f1)"' '  dg="$(skopeo inspect --no-creds "docker://ghcr.io/$IMAGE_NAME:$v" | sed -E '"'"'s/.*"Digest":"([^"]*)".*/\1/'"'"')"' \
  '  dd="sha256:$(sha256sum "$scratch/roh-docker.io" | cut -d'"'"' '"'"' -f1)"' '  dd="$(skopeo inspect --no-creds "docker://docker.io/$IMAGE_NAME:$v" | sed -E '"'"'s/.*"Digest":"([^"]*)".*/\1/'"'"')"'
mutr m-p7-marke-zu-frueh "p11-alte-marke-b p7-anonym-docker p7-anonym-ghcr p7-digest-docker p7-digest-ghcr p7-digest-lokal p7-eine-plattform" b '  rm -f "$zdir/vergleichen.ok" # P11vergleichen-marke' '  : > "$zdir/vergleichen.ok" # P11vergleichen-marke'
mutr m-p7-smoke-nur-eine-registry "p7-smoke-aufruf p7-smoke-docker" b $'  for reg in "${REGISTRIES[@]}"; do\n    mkdir "$scratch/smoke-$reg"' $'  for reg in ghcr.io; do\n    mkdir "$scratch/smoke-$reg"'
mutr m-p7-smoke-fehler-ignoriert "p7-smoke-docker p7-smoke-ghcr p7-smoke-ok-marke" b $'    [ "$rc" -eq 0 ] || rot "der Smoke aus $reg endet mit Exit $rc"\n' ''
mutr m-p7-smoke-referenz-ohne-digest "p7-smoke-aufruf" b 'container-smoke.sh "$reg/$IMAGE_NAME@$d" # P7smoke' 'container-smoke.sh "$reg/$IMAGE_NAME:$v" # P7smoke'
mutr m-p7-smoke-version-fehlt "p7-smoke-aufruf" b 'fuehre env VERSION="$v" POSTGRES_IMAGE=' 'fuehre env POSTGRES_IMAGE='
mutr m-p7-smoke-postgres-image-fehlt "p7-smoke-aufruf" b ' POSTGRES_IMAGE="$pg" SMOKE_LAUF=' ' SMOKE_LAUF='
mutr m-p7-smoke-lauf-gleich "p7-smoke-aufruf" b 'SMOKE_LAUF="pgr-rel-$$-${reg%%.*}"' 'SMOKE_LAUF="pgr-rel-$$"'
mutr m-p7-smoke-werk-fehlt "p7-smoke-aufruf" b 'SMOKE_WERK="$scratch/smoke-$reg"' 'SMOKE_WERK="$scratch/fehlt-$reg"'
mutr m-p7-smoke-name-weicht-ab "p7-smoke-aufruf" b 'SMOKE_NAME="release: smoke"' 'SMOKE_NAME="smoke"'
mutr m-p11-vergleichen-loescht-marke-nicht "p11-alte-marke-b" b 'L:# P11vergleichen-marke' ''
# Punkt 8
mutr m-p8-vorabversion-setzt-latest "p8-vorabversion" b 'L:# P8stabil' ''
mutr m-p8-vorabversionen-zaehlen "p8-fremder-tag p8-stabil-hoechste p8-vorabversion-hoeher-zaehlt-nicht" b 'grep -E "^v$KERN\$" | sed' 'grep -E "^v" | sed'
mutr m-p8-hoechste-entfaellt "p8-aelterer-zweig p8-semver-numerisch" b 'L:# P8hoechste' '  hoechste="$v"'
mutr m-p8-lexikographisch "p8-semver-numerisch" b ' | sort -V | tail -n 1)" # P8hoechste' ' | sort | tail -n 1)" # P8hoechste'
mutr m-p8-reihenfolge-vertauscht "p8-stabil-hoechste p8-zweite-registry-faellt" b 'for reg in "${REGISTRIES[@]}"; do # P8reg' 'for reg in docker.io ghcr.io; do # P8reg'
mutr m-p8-all-entfernt "p8-stabil-hoechste" b 'copy --all --preserve-digests --authfile "$scratch/auth.json" "docker://' 'copy --preserve-digests --authfile "$scratch/auth.json" "docker://'
mutr m-p8-preserve-digests-entfernt "p8-stabil-hoechste" b 'copy --all --preserve-digests --authfile "$scratch/auth.json" "docker://' 'copy --all --authfile "$scratch/auth.json" "docker://'
mutr m-p8-quelle-tag "p10-stdin p11-wiederanlauf-latest p8-erster-stabiler p8-fremder-tag p8-semver-numerisch p8-stabil-hoechste p8-tag-satz p8-vorabversion-hoeher-zaehlt-nicht p8-zweite-registry-faellt" b '"docker://$reg/$IMAGE_NAME@$d" "docker://$reg/$IMAGE_NAME:latest"' '"docker://$reg/$IMAGE_NAME:$v" "docker://$reg/$IMAGE_NAME:latest"'
mutr m-p8-zweite-registry-ignoriert "p8-zweite-registry-faellt p11-wiederanlauf-latest" b '[ "$rc" -eq 0 ] || rot "latest in $reg nicht gesetzt (Exit $rc); $(zustand_tag latest)"' ':'
mutr m-p8-zustand-fehlt "p8-zweite-registry-faellt" b '; $(zustand_tag latest)"' '"'
# Punkt 9 (B), 13
mutr m-p5-builder-arg-entfaellt "p5-aufruf" b 'bargs=(buildx build --builder pgr-release' 'bargs=(buildx build'
mutr m-p5-summe-waechter-entfernt "p5-summe-eintrag-fehlt" b '[ -n "$soll" ] && [ "$ist" = "$soll" ] ||' '[ "$ist" = "$soll" ] ||'
mutr m-p6-meldung-verworfen "p6-inspect-meldung-im-log" b 'L:# P6meldung' ''
mutr m-p7-meldung-verworfen "p6-inspect-meldung-im-log" b 'L:# P7meldung' ''
mutr m-p7-plattformen-docker-nicht-geprueft "p7-eine-plattform" b 'index_pruefen "$scratch/roh-$reg" "der Index aus $reg" # P7plattformen' '[ "$reg" = docker.io ] || index_pruefen "$scratch/roh-$reg" "der Index aus $reg" # P7plattformen'
mutr m-p7-plattformen-ghcr-nicht-geprueft "p7-eine-plattform" b 'index_pruefen "$scratch/roh-$reg" "der Index aus $reg" # P7plattformen' '[ "$reg" = ghcr.io ] || index_pruefen "$scratch/roh-$reg" "der Index aus $reg" # P7plattformen'
mutr m-p8-vorabversion-muster-auf-rc "p8-vorabversion" b 'if [[ $v == *-* ]]; then echo "release: latest:' 'if [[ $v == *-rc* ]]; then echo "release: latest:'
mutr m-p8-zusaetzlicher-tag "p8-tag-satz" b '    echo "release: latest: $reg zeigt auf $d"' $'    fuehre skopeo copy --all --preserve-digests --authfile "$scratch/auth.json" "docker://$reg/$IMAGE_NAME@$d" "docker://$reg/$IMAGE_NAME:${v%.*}"\n    echo "release: latest: $reg zeigt auf $d"'
mutr m-p5-builder-ohne-digest "p5-builder-pin-form" b "@sha256:cec9f139f45e93c5c69c60f8b07cfad9f43f4ef6b6a6cd917527fea5ff2e3dea' # v0.33.1" ":latest' # v0.33.1"
mutr m-p5-builder-platzhalter "p5-builder-platzhalter" b "@sha256:cec9f139f45e93c5c69c60f8b07cfad9f43f4ef6b6a6cd917527fea5ff2e3dea'" "@sha256:0000000000000000000000000000000000000000000000000000000000000000'"
mutr m-p5-builder-ohne-tag-kommentar "p5-builder-pin-form" b "dea' # v0.33.1" "dea'"
mutr m-p5-builder-pin-nicht-uebergeben "p5-aufruf" b ' --driver-opt "image=$BUILDKIT_IMAGE"' ''
mutr m-p5-zeit-aus-head "p5-aufruf" b 'git log -1 --format=%ct "${RELEASE_SHA-}"' 'git log -1 --format=%ct HEAD'
mutr m-p5-epoch-autor-zeit "p5-epoch-committer" b 'git log -1 --format=%ct' 'git log -1 --format=%at'
mutr m-p5-ein-index-entfernt "p5-layout-form" b 'L:# P5ein-index' ''
mutr m-p5-index-typ-entfernt "p5-layout-form" b 'L:# P5index-typ' ''
mutr m-p3-secret-nur-der-erste "p3-secret-fehlen-beide" b '|| fehlt+=" DOCKERHUB_TOKEN"' '|| [ -n "$fehlt" ] || fehlt+=" DOCKERHUB_TOKEN"'
mutr m-p6-fehler-wie-fehlend "p6-inspect-anderer-fehler p7-anonym-docker p7-anonym-ghcr" b "FEHLT_RE='manifest unknown|name unknown|not found'" "FEHLT_RE='.'"
mutr m-p6-muster-geweitet "p6-inspect-anderer-fehler" b "FEHLT_RE='manifest unknown|name unknown|not found'" "FEHLT_RE='manifest unknown|name unknown|not found|denied'"
mutr m-p6-wortlaut-manifest-unknown-nicht-erkannt "p10-nicht-im-log-image p10-stdin p11-alte-marke-b p11-wiederanlauf-ghcr-vorhanden p11-wiederanlauf-latest p5-digest p6-all-und-digests p6-inspect-meldung-im-log p6-reihenfolge p6-tag-fehlt-manifest-unknown p6-tag-satz p6-tag-vorhanden-gleich p6-zweite-registry-faellt p7-anonym-docker p7-anonym-ghcr p7-digest-docker p7-digest-ghcr p7-digest-lokal p7-eine-plattform p7-gleich p7-postgres-fehlt p7-rohbytes p7-smoke-aufruf p7-smoke-docker p7-smoke-ghcr p7-smoke-ohne-vergleichen p7-smoke-ok-marke p7-tag-fehlt p8-aelterer-zweig p8-erster-stabiler p8-fremder-tag p8-nach-rotem-vergleich p8-ohne-smoke p8-semver-numerisch p8-stabil-hoechste p8-tag-satz p8-vorabversion p8-vorabversion-hoeher-zaehlt-nicht p8-zweite-registry-faellt p9-nicht-vor-registries" b "FEHLT_RE='manifest unknown|name unknown|not found'" "FEHLT_RE='name unknown|not found'"
mutr m-p6-wortlaut-name-unknown-nicht-erkannt "p6-tag-fehlt-name-unknown" b "FEHLT_RE='manifest unknown|name unknown|not found'" "FEHLT_RE='manifest unknown|not found'"
mutr m-p6-wortlaut-not-found-nicht-erkannt "p6-tag-fehlt-not-found" b "FEHLT_RE='manifest unknown|name unknown|not found'" "FEHLT_RE='manifest unknown|name unknown'"
mutr m-p6-schreibweise-beachtet "p6-tag-fehlt-not-found" b 'grep -qiE "$FEHLT_RE" "$scratch/skopeo.err" || rot "skopeo inspect' 'grep -qE "$FEHLT_RE" "$scratch/skopeo.err" || rot "skopeo inspect'
mutr m-p7-postgres-quelle-fehlt-angenommen "p7-postgres-fehlt" b '  [ -n "$pg" ] || eingabe vorbedingung "POSTGRES_IMAGE steht nicht in harness/mk/integration.mk"' ':'
mutr m-p7-smoke-marke-bei-rot "p7-smoke-ok-marke" b '  rm -f "$zdir/smoke.ok" # P11smoke-marke' '  : > "$zdir/smoke.ok" # P11smoke-marke'
mutr m-p7-smoke-marke-entfaellt "p10-stdin p11-wiederanlauf-latest p7-smoke-ok-marke p8-aelterer-zweig p8-erster-stabiler p8-fremder-tag p8-semver-numerisch p8-stabil-hoechste p8-tag-satz p8-vorabversion p8-vorabversion-hoeher-zaehlt-nicht p8-zweite-registry-faellt" b 'L:# P7smoke-ok' ''
mutr m-p11-smoke-loescht-marke-nicht "p7-smoke-ok-marke" b 'L:# P11smoke-marke' ''
mutr m-p11-eingaben-loescht-smoke-nicht "p11-alte-marke-b" b 'L:# P11loescht-smoke' ''
mutr m-p8-smoke-marke-nicht-verlangt "p11-reihenfolge-gruende-b p8-nach-rotem-vergleich p8-ohne-smoke" b 'L:# P8vergleich' ''
mutr m-p9-smoke-marke-nicht-verlangt "p11-alte-marke-b p9-nicht-vor-registries" b 'L:# P9registries' ''
mutr m-p8-fremder-tag-zaehlt "p8-fremder-tag p8-stabil-hoechste p8-vorabversion-hoeher-zaehlt-nicht" b 'grep -E "^v$KERN\$" | sed' 'grep -E "^v[0-9]" | sed'
mutr m-p13-image-ohne-binaries "p11-reihenfolge-gruende-b p5-ohne-binaries" b 'L:# P13image' ''
mutr m-p13-veroeffentlichen-ohne-image "p11-alte-marke-image-b p11-image-rot-loescht-digest-b p13-veroeffentlichen-ohne-image-bei-scan-ok" b 'L:# P13veroeffentlichen' ''
mutr m-p13-vergleichen-ohne-image "p11-reihenfolge-gruende-b" b 'L:# P13vergleichen' ''
mutr m-p13-smoke-ohne-vergleichen "p11-alte-marke-b p11-reihenfolge-gruende-b p7-smoke-ohne-vergleichen" b 'L:# P13smoke' ''
# Punkt 10 (B)
mutr m-p10-token-bei-fehler-registry "p10-nicht-im-log-image" b '[ "$have" = "$d" ] || rot "$reg trägt' '[ "$have" = "$d" ] || { echo "$RELEASE_DOCKERHUB_TOKEN" >&2; false; } || rot "$reg trägt'
mutr m-p10-env-bei-fehler-registry "p10-nicht-im-log-image" b '[ "$have" = "$d" ] || rot "$reg trägt' '[ "$have" = "$d" ] || { env >&2; false; } || rot "$reg trägt'
mutr m-p10-set-x-image "p10-nicht-im-log-image p3-schluesselwort-geheimnis p3-secret-fehlen-beide p3-secret-fehlt-user p3-secret-leer-user" b $'set -uo pipefail\nexport LC_ALL=C\n' $'set -uo pipefail\nset -x\nexport LC_ALL=C\n'
mutr m-p10-token-als-argument "p10-stdin" b 'printf '"'"'%s'"'"' "$t" | skopeo login --authfile "$scratch/auth.json" --username "$u" --password-stdin "$1"' 'skopeo login --authfile "$scratch/auth.json" --username "$u" --password "$t" "$1"'
# Punkt 11 (B): kein Schritt läuft mit dem Zustand eines früheren Laufs
mutr m-p11-eingaben-loescht-vergleichen-nicht "p11-alte-marke-b" b 'rm -rf "$zdir/vergleichen.ok" "$zdir/image" # P11loescht-b' 'rm -rf "$zdir/image" # P11loescht-b'
mutr m-p11-eingaben-loescht-image-nicht "p11-alte-marke-image-b" b 'rm -rf "$zdir/vergleichen.ok" "$zdir/image" # P11loescht-b' 'rm -rf "$zdir/vergleichen.ok" # P11loescht-b'
mutr m-p11-image-loescht-nicht "p11-image-rot-loescht-digest-b" b 'rm -rf "$zdir/image"; mkdir -p "$zdir/image"' 'mkdir -p "$zdir/image"'
# Punkt 16: Scan (Fälle der Gruppe f)
TD=af6acf9a6b85dfe389a1941505c0ce9efef52a4719635e1a962f022a3d855daa
mf() { mutr "$1" "$2" f "${@:3}"; }
# ausser <Fall>… — alle Fälle der Gruppe f außer den genannten (für Mutationen, die fast jeden Fall färben, weil der Scan dann nicht mehr grün endet).
ausser() { printf '%s\n' "${FAELLE_F[@]}" | grep -vxF -f <(printf '%s\n' "$@") | tr '\n' ' '; }
mf m-p16-pin-ohne-digest "p16-pin-form" "@sha256:$TD' # 0.75.0" ":0.75.0' # 0.75.0"
mf m-p16-pin-platzhalter "p16-pin-platzhalter" "@sha256:$TD'" "@sha256:$(printf '0%.0s' {1..64})'"
mf m-p16-pin-ohne-tag-kommentar "p16-pin-form" "$TD' # 0.75.0" "$TD'"
mf m-p16-e3-anderes-image "p16-e3" '-v "$cache:/cache" "$TRIVY_IMAGE" \' '-v "$cache:/cache" aquasec/trivy:latest \'
mf m-p16-e1-rm-entfaellt "p16-e1" 'docker run --rm --user' 'docker run --user'
mf m-p16-e2-ro-entfaellt "p16-e2" ':/layout.tar:ro"' ':/layout.tar"'
mf m-p16-e2-quelle-kopie "p16-e2-quelle" '-v "$zdir/image/layout.tar:/layout.tar:ro"' '-v "$scratch/layout.tar:/layout.tar:ro"'
mf m-p16-e4-unterbefehl-entfaellt "p16-e4" $'      image --input /layout.tar' $'      --input /layout.tar'
mf m-p16-e4-input-anderer-pfad "p16-e4-input" '--input /layout.tar --format' '--input /layout --format'
mf m-p16-e5-format-entfaellt "p16-e5" ' --format json --exit-code' ' --exit-code'
mf m-p16-e5-format-table "p16-e5" '--format json' '--format table'
mf m-p16-e6-exit-code-entfaellt "p16-e6" ' --exit-code 0 --severity' ' --severity'
mf m-p16-e6-exit-code-1 "p16-e6" '--exit-code 0' '--exit-code 1'
mf m-p16-e7-severity-entfaellt "p16-e7" ' --severity HIGH,CRITICAL --timeout' ' --timeout'
mf m-p16-e7-nur-critical "p16-e7" '--severity HIGH,CRITICAL' '--severity CRITICAL'
mf m-p16-e7-mit-medium "p16-e7" '--severity HIGH,CRITICAL' '--severity HIGH,CRITICAL,MEDIUM'
mf m-p16-e8-timeout-entfaellt "p16-e8" ' --timeout 10m --platform' ' --platform'
mf m-p16-e8-timeout-60m "p16-e8" '--timeout 10m' '--timeout 60m'
mf m-p16-e9-platform-entfaellt "$(ausser p16-pin-form p16-pin-platzhalter p16-plattform-abweichend p16-plattform-abweichend-arm64 p16-plattform-fehlt p16-ohne-image p16-ohne-image-mit-layout p16-ohne-layout p16-veroeffentlichen-ohne-scan)" ' --platform "$plat" --cache-dir' ' --cache-dir'
mf m-p16-e9-nur-amd64 "p16-alle-befunde p16-befund-nur-arm64 p16-db-alt-arm64 p16-db-unlesbar-arm64 p16-e9-arm64 p16-e9-genau-zwei p16-kein-json-arm64 p16-leeres-ergebnis-arm64 p16-plattform-abweichend-arm64 p16-plattform-genannt p16-trivy-exit-arm64 p16-trockenlauf" 'SCAN_PLATTFORMEN=(linux/amd64 linux/arm64)' 'SCAN_PLATTFORMEN=(linux/amd64)'
mf m-p16-e9-nur-arm64 "$(ausser p16-gruen p16-e1 p16-e2 p16-e2-quelle p16-e3 p16-e4 p16-e4-input p16-e5 p16-e6 p16-e7 p16-e8 p16-e9-arm64 p16-e10-skip-db-update p16-e10-skip-files p16-e10-skip-dirs p16-e10-ignorefile p16-e10-ignore-policy p16-e10-vex p16-e10-mount p16-user p16-cache-dir p16-cache-kein-zwischenspeicher p16-pin-form p16-pin-platzhalter p16-medium-mit-fix p16-low-mit-fix p16-unknown-mit-fix p16-befund-nur-arm64 p16-trivy-exit-arm64 p16-kein-json-arm64 p16-leeres-ergebnis-arm64 p16-db-alt-arm64 p16-db-frisch p16-db-unlesbar-arm64 p16-plattform-abweichend-arm64 p16-ohne-image p16-ohne-image-mit-layout p16-ohne-layout p16-eingaben-loescht-marke p16-image-loescht-marke p16-veroeffentlichen-ohne-scan p16-d-unveraendert)" 'SCAN_PLATTFORMEN=(linux/amd64 linux/arm64)' 'SCAN_PLATTFORMEN=(linux/arm64)'
mf m-p16-e9-dritte-plattform "p16-e9-genau-zwei p16-trockenlauf" 'SCAN_PLATTFORMEN=(linux/amd64 linux/arm64)' 'SCAN_PLATTFORMEN=(linux/amd64 linux/arm64 linux/s390x)'
for fl in skip-db-update skip-files skip-dirs ignorefile ignore-policy vex; do mf "m-p16-e10-$fl" "p16-e10-$fl" ' --cache-dir /cache > "$json"' " --cache-dir /cache --$fl > \"\$json\""; done
mf m-p16-e10-dritter-mount "p16-e10-mount" '-v "$cache:/cache" ' '-v "$cache:/cache" -v /ignore:/ignore '
mf m-p16-user-entfaellt "p16-user" ' --user "$(id -u):$(id -g)"' ''
mf m-p16-cache-dir-entfaellt "p16-cache-dir" ' --cache-dir /cache >' ' >'
mf m-p16-cache-ueber-laeufe "p16-cache-kein-zwischenspeicher" 'cache="$scratch/trivy"' 'cache="$zdir/trivy"'
mf m-p16-schwelle-medium "p16-alle-befunde p16-medium-mit-fix p16-summe" 'select(.Severity == "CRITICAL" or .Severity == "HIGH")' 'select(.Severity == "CRITICAL" or .Severity == "HIGH" or .Severity == "MEDIUM")'
mf m-p16-schwelle-nur-critical "p16-alle-befunde p16-befund-im-zweiten-ergebnis p16-befund-nur-amd64 p16-befund-nur-arm64 p16-high-mit-fix p16-high-ohne-fix p16-high-ohne-fix-leer p16-high-ohne-fix-null p16-plattform-genannt p16-summe" 'select(.Severity == "CRITICAL" or .Severity == "HIGH")' 'select(.Severity == "CRITICAL")'
mf m-p16-schwelle-nur-high "p16-alle-befunde p16-alte-marke p16-critical-mit-fix p16-critical-ohne-fix p16-plattform-genannt p16-scan-ok-marke p16-summe p16-ueber-schwelle-keine-registry" 'select(.Severity == "CRITICAL" or .Severity == "HIGH")' 'select(.Severity == "HIGH")'
mf m-p16-schwelle-jede "p16-alle-befunde p16-low-mit-fix p16-medium-mit-fix p16-summe p16-unknown-mit-fix" 'select(.Severity == "CRITICAL" or .Severity == "HIGH")' 'select(.Severity != "")'
mf m-p16-ohne-fix-sperrt "p16-alle-befunde p16-critical-ohne-fix p16-high-ohne-fix p16-high-ohne-fix-leer p16-high-ohne-fix-null p16-summe" "SCAN_MIT_FIX='select((.FixedVersion // \"\") != \"\")'" "SCAN_MIT_FIX='select(true)'"
mf m-p16-fix-null-gilt-als-fix "p16-alle-befunde p16-critical-ohne-fix p16-high-ohne-fix p16-high-ohne-fix-null p16-summe" "SCAN_MIT_FIX='select((.FixedVersion // \"\") != \"\")'" "SCAN_MIT_FIX='select(.FixedVersion != \"\")'"
mf m-p16-fix-leer-gilt-als-fix "p16-high-ohne-fix-leer p16-summe" "SCAN_MIT_FIX='select((.FixedVersion // \"\") != \"\")'" "SCAN_MIT_FIX='select(.FixedVersion != null)'"
mf m-p16-ohne-fix-zaehlt-alle "p16-alle-befunde p16-critical-mit-fix p16-high-mit-fix p16-summe" "SCAN_OHNE_FIX='select((.FixedVersion // \"\") == \"\")'" "SCAN_OHNE_FIX='select(true)'"
mf m-p16-nur-erster-befund "p16-alle-befunde p16-summe" $'bef+="$b"$\'\\n\'' $'bef+="${b%%$\'\\n\'*}"$\'\\n\''
mf m-p16-nur-erstes-ergebnis "p16-alle-befunde p16-befund-im-zweiten-ergebnis" "SCAN_EINTRAG='.Results[] |" "SCAN_EINTRAG='.Results[0] |"
mf m-p16-arm64-nicht-ausgewertet "p16-alle-befunde p16-befund-nur-arm64 p16-plattform-genannt" '[ -z "$b" ] || bef+=' '[ -z "$b" ] || [ "$plat" = linux/arm64 ] || bef+='
mf m-p16-amd64-nicht-ausgewertet "p16-alle-befunde p16-alte-marke p16-befund-im-zweiten-ergebnis p16-befund-nur-amd64 p16-critical-mit-fix p16-high-mit-fix p16-plattform-genannt p16-scan-ok-marke p16-summe p16-ueber-schwelle-keine-registry" '[ -z "$b" ] || bef+=' '[ -z "$b" ] || [ "$plat" = linux/amd64 ] || bef+='
mf m-p16-plattform-nicht-genannt "p16-alle-befunde p16-befund-im-zweiten-ergebnis p16-befund-nur-amd64 p16-befund-nur-arm64 p16-critical-mit-fix p16-high-mit-fix p16-plattform-genannt" '\($p) \(.VulnerabilityID)' '\(.VulnerabilityID)'
mf m-p16-summe-entfaellt "p16-alle-befunde p16-critical-mit-fix p16-high-mit-fix p16-summe" 'L:# P16summe' ''
mf m-p16-summe-zaehlt-unbehebbare-mit "p16-alle-befunde p16-summe" 'n="$(grep -c . <<< "$bef")"' 'n="$(( $(grep -c . <<< "$bef") + m ))"'
mf m-p16-befunde-nicht-genannt "p16-alle-befunde p16-befund-im-zweiten-ergebnis p16-befund-nur-amd64 p16-befund-nur-arm64 p16-critical-mit-fix p16-high-mit-fix p16-plattform-genannt" 'L:# P16alle' ''
mf m-p16-auswertung-aus-anzeige "p16-alle-befunde p16-befund-im-zweiten-ergebnis p16-befund-nur-amd64 p16-befund-nur-arm64 p16-critical-mit-fix p16-critical-ohne-fix p16-high-mit-fix p16-high-ohne-fix p16-high-ohne-fix-leer p16-high-ohne-fix-null p16-plattform-genannt p16-summe" 'L:# P16auswertung' $'    b="$(grep -E \'CRITICAL|HIGH\' "$json" || true)" # P16auswertung'
mf m-p16-exit-gilt-als-gruen "p16-trivy-exit p16-trivy-exit-125 p16-trivy-exit-arm64" $'if [ "$rc" -ne 0 ]; then\n      tail' $'if false; then\n      tail'
mf m-p16-exit-nur-1 "p16-trivy-exit-125" $'if [ "$rc" -ne 0 ]; then\n      tail' $'if [ "$rc" -eq 1 ]; then\n      tail'
mf m-p16-exit-arm64-ignoriert "p16-trivy-exit-arm64" $'if [ "$rc" -ne 0 ]; then\n      tail' $'if [ "$rc" -ne 0 ] && [ "$plat" = linux/amd64 ]; then\n      tail'
mf m-p16-pruefmenge-leer-gruen "p16-leeres-ergebnis p16-leeres-ergebnis-arm64" '(length > 0)' '(length >= 0)'
mf m-p16-pruefmenge-entfernt "p16-kein-json p16-kein-json-arm64 p16-kein-json-leer p16-leeres-ergebnis p16-leeres-ergebnis-arm64" 'L:# P16pruefmenge' ''
mf m-p16-pruefmenge-arm64-ignoriert "p16-kein-json-arm64 p16-leeres-ergebnis-arm64" "jq -e '(.Results // [])" "[ \"\$plat\" = linux/arm64 ] || jq -e '(.Results // [])"
mf m-p16-db-nicht-geprueft "p16-db-alt p16-db-alt-arm64" 'L:# P16db-alt' ''
mf m-p16-db-7-tage "p16-db-alt p16-db-alt-arm64" 'DB_GRENZE=$((72 * 3600))' 'DB_GRENZE=$((168 * 3600))'
mf m-p16-db-70-stunden "p16-db-frisch" 'DB_GRENZE=$((72 * 3600))' 'DB_GRENZE=$((70 * 3600))'
mf m-p16-db-74-stunden "p16-db-alt p16-db-alt-arm64" 'DB_GRENZE=$((72 * 3600))' 'DB_GRENZE=$((74 * 3600))'
mf m-p16-db-downloadedat "p16-db-alt p16-db-alt-arm64" "'.UpdatedAt | strings'" "'.DownloadedAt | strings'"
mf m-p16-db-arm64-nicht-geprueft "p16-db-alt-arm64 p16-db-unlesbar-arm64" '    db_alter "$cache"' '    [ "$plat" = linux/arm64 ] || db_alter "$cache"'
mf m-p16-db-zeitpunkt-ungeprueft "p16-db-unlesbar p16-db-unlesbar-arm64 p16-db-unlesbar-jetzt p16-db-unlesbar-kein-json p16-db-unlesbar-leer p16-db-unlesbar-ohne-feld p16-db-unlesbar-zahl" 'L:# P16db-zeitpunkt' ''
mf m-p16-db-datum-ungeprueft "p16-db-unlesbar-kein-datum" ' || scan_fehler datenbank "UpdatedAt ${t@Q} ist kein Zeitpunkt" # P16db-datum' ' # P16db-datum'
mf m-p16-plattform-nicht-verglichen "p16-plattform-abweichend p16-plattform-abweichend-arm64 p16-plattform-fehlt" 'L:# P16plattform' ''
mf m-p16-plattform-arm64-nicht-verglichen "p16-plattform-abweichend-arm64" '[ "$arch" = "${plat#linux/}" ] ||' '[ "$arch" = "${plat#linux/}" ] || [ "$plat" = linux/arm64 ] ||'
mf m-p16-plattform-nur-nicht-leer "p16-plattform-abweichend p16-plattform-abweichend-arm64" '[ "$arch" = "${plat#linux/}" ]' '[ -n "$arch" ]'
mf m-p16-image-marke-nicht-verlangt "p16-ohne-image-mit-layout" 'L:# P16image' ''
mf m-p16-layout-nicht-verlangt "p16-ohne-layout" 'L:# P16layout' ''
mf m-p16-vorbedingung-beide-entfernt "p16-ohne-image p16-ohne-image-mit-layout p16-ohne-layout" 'L:# P16image' '' 'L:# P16layout' ''
mf m-p16-marke-bei-rot "$(ausser p16-trockenlauf p16-gruen p16-e1 p16-e2 p16-e2-quelle p16-e3 p16-e4 p16-e4-input p16-e5 p16-e6 p16-e7 p16-e8 p16-e9-amd64 p16-e9-arm64 p16-e9-genau-zwei p16-e10-skip-db-update p16-e10-skip-files p16-e10-skip-dirs p16-e10-ignorefile p16-e10-ignore-policy p16-e10-vex p16-e10-mount p16-user p16-cache-dir p16-cache-kein-zwischenspeicher p16-pin-form p16-pin-platzhalter p16-medium-mit-fix p16-low-mit-fix p16-unknown-mit-fix p16-high-ohne-fix p16-high-ohne-fix-leer p16-high-ohne-fix-null p16-critical-ohne-fix p16-db-frisch p16-ohne-image p16-ohne-image-mit-layout p16-ohne-layout p16-eingaben-loescht-marke p16-image-loescht-marke p16-veroeffentlichen-ohne-scan p16-d-unveraendert)" 'rm -f "$zdir/scan.ok" # P16marke-loescht' ': > "$zdir/scan.ok"'
mf m-p16-marke-entfaellt "p16-critical-ohne-fix p16-db-frisch p16-e2-quelle p16-gruen p16-high-ohne-fix p16-high-ohne-fix-leer p16-high-ohne-fix-null p16-low-mit-fix p16-medium-mit-fix p16-scan-ok-marke p16-trockenlauf p16-unknown-mit-fix" 'L:: > "$zdir/scan.ok" # P16marke' ''
mf m-p16-marke-loescht-entfaellt "p16-alte-marke" 'L:# P16marke-loescht' ''
mf m-p16-eingaben-loescht-scan-nicht "p16-eingaben-loescht-marke" 'L:# P16loescht-eingaben' ''
mf m-p16-image-loescht-scan-nicht "p16-image-loescht-marke" 'L:# P16loescht-image' ''
mf m-p16-veroeffentlichen-ohne-scan-marke "p16-ueber-schwelle-keine-registry p16-veroeffentlichen-ohne-scan" 'L:# P16veroeffentlichen' ''
mf m-p16-scan-ruft-skopeo "p16-gruen p16-trockenlauf" ': > "$zdir/scan.ok" # P16marke' $'skopeo inspect docker://x\n  : > "$zdir/scan.ok" # P16marke'
mf m-p16-trockenlauf-ruft-skopeo "p16-trockenlauf" ': > "$zdir/scan.ok" # P16marke' $'[ "${RELEASE_EVENT-}" = push ] || skopeo inspect docker://x\n  : > "$zdir/scan.ok" # P16marke'
mf m-p16-scan-verlangt-push "p16-trockenlauf" $'  local plat arch cache json b bef="" n m=0 k z\n' $'  local plat arch cache json b bef="" n m=0 k z\n  aussenwirkung_erlaubt\n'
mf m-p16-digest-veraendert "p16-d-unveraendert" ': > "$zdir/scan.ok" # P16marke' $'echo x >> "$zdir/image/layout.tar"\n  : > "$zdir/scan.ok" # P16marke'
mf m-p16-stichwort-schwachstelle "p16-alle-befunde p16-befund-im-zweiten-ergebnis p16-befund-nur-amd64 p16-befund-nur-arm64 p16-critical-mit-fix p16-high-mit-fix p16-plattform-genannt" 'release: scan: schwachstelle: $z' 'release: scan: befund: $z'
mf m-p16-stichwort-scan-gescheitert "p16-trivy-exit p16-trivy-exit-125 p16-trivy-exit-arm64" 'scan_fehler scan-gescheitert' 'scan_fehler scan-fehler'
mf m-p16-stichwort-pruefmenge "p16-kein-json p16-kein-json-arm64 p16-kein-json-leer p16-leeres-ergebnis p16-leeres-ergebnis-arm64 p16-leeres-ergebnis-fehlt p16-leeres-ergebnis-null" 'scan_fehler pruefmenge "die Ausgabe von Trivy für $plat ist kein JSON' 'scan_fehler leer "die Ausgabe von Trivy für $plat ist kein JSON'
mf m-p16-stichwort-datenbank-alt "p16-db-alt p16-db-alt-arm64" 'scan_fehler datenbank "die Datenbank ist' 'scan_fehler db "die Datenbank ist'
mf m-p16-stichwort-datenbank-lesbar "p16-db-unlesbar p16-db-unlesbar-arm64 p16-db-unlesbar-jetzt p16-db-unlesbar-kein-json p16-db-unlesbar-leer p16-db-unlesbar-ohne-feld p16-db-unlesbar-zahl" 'scan_fehler datenbank "das Alter' 'scan_fehler db "das Alter'
mf m-p16-stichwort-plattform "p16-plattform-abweichend p16-plattform-abweichend-arm64 p16-plattform-fehlt" 'scan_fehler plattform' 'scan_fehler arch'
mf m-p16-fehler-exit-1 "$(ausser p16-trockenlauf p16-gruen p16-e1 p16-e2 p16-e2-quelle p16-e3 p16-e4 p16-e4-input p16-e5 p16-e6 p16-e7 p16-e8 p16-e9-amd64 p16-e9-arm64 p16-e9-genau-zwei p16-e10-skip-db-update p16-e10-skip-files p16-e10-skip-dirs p16-e10-ignorefile p16-e10-ignore-policy p16-e10-vex p16-e10-mount p16-user p16-cache-dir p16-cache-kein-zwischenspeicher p16-pin-form p16-pin-platzhalter p16-critical-mit-fix p16-high-mit-fix p16-medium-mit-fix p16-low-mit-fix p16-unknown-mit-fix p16-high-ohne-fix p16-high-ohne-fix-leer p16-high-ohne-fix-null p16-critical-ohne-fix p16-alle-befunde p16-plattform-genannt p16-befund-nur-amd64 p16-befund-nur-arm64 p16-befund-im-zweiten-ergebnis p16-summe p16-db-frisch p16-ohne-image p16-ohne-image-mit-layout p16-ohne-layout p16-scan-ok-marke p16-alte-marke p16-eingaben-loescht-marke p16-image-loescht-marke p16-veroeffentlichen-ohne-scan p16-ueber-schwelle-keine-registry p16-d-unveraendert)" 'echo "release: scan: $1: $2" >&2; exit 2; }' 'echo "release: scan: $1: $2" >&2; exit 1; }'
mf m-p16-befund-exit-2 "p16-alle-befunde p16-alte-marke p16-befund-im-zweiten-ergebnis p16-befund-nur-amd64 p16-befund-nur-arm64 p16-critical-mit-fix p16-high-mit-fix p16-plattform-genannt p16-summe p16-ueber-schwelle-keine-registry" $'    exit 1\n  fi\n  : > "$zdir/scan.ok" # P16marke' $'    exit 2\n  fi\n  : > "$zdir/scan.ok" # P16marke'
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
