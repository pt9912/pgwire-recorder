#!/usr/bin/env bash
# container-gegenprobe — prüft `binaries` und `test-container` (harness/mk/container.mk,
# tools/harness/binaries.sh, run-container.sh, container-eingabe.sh) gegen die Punkte 1 bis 9,
# 11 und 12 von SPEC-055 (harness/vertraege/SPEC-055-container-und-binaries.md) in Kopien des
# Arbeitsbaums in einem Temp-Verzeichnis; der Arbeitsbaum bleibt unberührt. docker ist eine
# Attrappe (über DOCKER): Sie protokolliert jeden Aufruf, führt Buch über Container, Netz und
# Tags und liefert Binaries mit dem Kopf ihres Formats, Image-Metadaten, Anfragen und Exit-Codes
# nach Vorgabe des Falls. Ein sleep in PATH wartet nicht. Der Name jedes Falls trägt die Nummer
# seines Punktes; die unveränderte Kopie besteht jeden Fall, und bei jeder Mutation (Textersatz
# in der Kopie) werden von den Fällen ihrer Gruppe (`nur`) genau die benannten rot.
#
# GRENZE der Gegenprobe: Mit der Attrappe baut kein Dockerfile und läuft kein Binary, kein
# PostgreSQL und kein Container. Ob zwei echte Builds dieselbe Summe liefern, ob die Stufen
# `dist` und `runtime` bauen, ob das Image Aufzeichnen und Wiedergeben besteht und ob
# `--no-cache-filter` den Cache umgeht, zeigt nur `make test-container` (SPEC-055 Punkt 12).
# Ein Zeitstempel im Binary (Punkt 4) ist kein Fall: Er färbt Punkt 4 rot, und der Fall
# der abweichenden Summe belegt das. Die Signale gehen an das Skript, nicht an `make`; sie
# setzen voraus, dass diese Gegenprobe SIGINT nicht ignoriert (kein Start als Hintergrundjob).
# Punkt 10 prüft der Fall des unveränderten Arbeitsbaums; Punkt 13 ist eine Grenze ohne Fall.
#
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Mutationen, deren Name passt; GP_ENTDECKE=1 nennt je
# Mutation die roten Fälle.
#
# Ausgang: 0, wenn jeder Fall und jede Mutation wie erwartet endet, sonst 1.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
wurzel="$(cd "$hier/../.." && pwd)"
arbeit="$(mktemp -d)"
trap 'rm -rf "${arbeit:?}"' EXIT
fehler=0
melde() { echo "container-gegenprobe: ROT — $*" >&2; fehler=1; }

bin="$arbeit/bin" log="$arbeit/log" st="$arbeit/zustand" aus="$arbeit/aus" tmp="$arbeit/tmp" gesendet="$arbeit/gesendet"
mkdir -p "$bin"

# Attrappe von docker. Vorgaben des Falls (Umgebung, Plattformen als `os/arch`):
# AT_BAU_ROT=<Plattformen> (der Build der Stufe dist scheitert), AT_ABWEICH=<Plattformen> und
# AT_NICHT_REPRO=1 (ein Build mit beiden --no-cache-filter liefert andere Bytes; ohne beide
# die zwischengespeicherte Fassung), AT_KOPF=<Plattform>=arch|elf|macho|pe (anderer Kopf),
# AT_IMG_BAU_ROT, AT_IMG_OS, AT_IMG_ARCH, AT_IMG_SUMME=<Plattformen>, AT_SERVER_ARCH,
# AT_VERSION_AUS=<Text>, AT_VERSION_OHNE_LF, AT_VERSION_ROT, AT_NICHT_BEREIT, AT_KEIN_START,
# AT_SICHT_A_LEER, AT_SICHT_LEER (beide Sichten), AT_SICHT_OHNE_EINS, AT_SICHT_B_ABW, AT_PSQL_ROT, AT_EXIT_REC, AT_EXIT_REP,
# AT_CP_ROT, AT_AUFZ_LEER, AT_PG_BLEIBT, AT_HANG_ORT=build|psql mit AT_HANGT=<Marke> (die Attrappe wartet höchstens 5 s und vermerkt SIGTERM und SIGINT in der Datei `term`).
cat > "$bin/docker" <<'ATTRAPPE'
#!/usr/bin/env bash
printf '%s' "$*" | tr '\n' ' ' >> "$AT_LOG"; echo >> "$AT_LOG"
S="$AT_STATE"
sub="$1"; shift
hat() { [[ " $1 " == *" $2 "* ]]; }
hang() {
  if [ "${AT_HANG_ORT-}" = "$1" ]; then
    trap 'echo "$sub" >> "$S/term"; exit 143' TERM INT
    : > "$AT_HANGT"
    for _ in $(seq 100); do /bin/sleep 0.05; done
    exit 1
  fi
}
hexbytes() { printf "$(printf '\\x%s' "$@")"; }
nul() { head -c "$1" /dev/zero; }
kopf() { # <Format> <Architektur>
  local a=$2
  case $1 in
    elf) hexbytes 7f 45 4c 46; nul 14; if [ "$a" = amd64 ]; then hexbytes 3e 00; else hexbytes b7 00; fi; nul 8 ;;
    macho) hexbytes cf fa ed fe; if [ "$a" = amd64 ]; then hexbytes 07 00 00 01; else hexbytes 0c 00 00 01; fi; nul 8 ;;
    pe) printf MZ; nul 58; hexbytes 98 00 00 00; nul 88; printf 'PE\0\0'; if [ "$a" = amd64 ]; then hexbytes 64 86; else hexbytes 64 aa; fi; nul 8 ;;
  esac
}
binaer() { # <Plattform> <Version> <Variante>
  local os=${1%/*} a=${1#*/} f k
  case $os in linux) f=elf ;; darwin) f=macho ;; windows) f=pe ;; esac
  for k in ${AT_KOPF-}; do
    if [ "${k%%=*}" = "$1" ]; then
      case ${k#*=} in
        arch) if [ "$a" = amd64 ]; then a=arm64; else a=amd64; fi ;;
        *) f=${k#*=} ;;
      esac
    fi
  done
  kopf "$f" "$a"; echo "version=$2 plattform=$1 variante=$3"
}
bild() { for a; do [ -d "$S/i-$a" ] && { echo "$a"; return; }; done; }
case "$sub" in
  build)
    target= plat= ver=dev tag= dest= ncf=
    while [ $# -gt 0 ]; do
      case $1 in
        --target) target=$2; shift ;; --platform) plat=$2; shift ;; --build-arg) ver=${2#VERSION=}; shift ;;
        -t) tag=$2; shift ;; --output) dest=${2#*dest=}; shift ;; --no-cache-filter) ncf+=" $2"; shift ;;
      esac
      shift
    done
    if [ "$target" = dist ]; then
      hat "${AT_BAU_ROT-}" "$plat" && { echo "build gescheitert" >&2; exit 1; }
      hang build
      var=v0
      if hat "$ncf" build && hat "$ncf" dist && { hat "${AT_ABWEICH-}" "$plat" || [ -n "${AT_NICHT_REPRO-}" ]; }; then var=v1; fi
      mkdir -p "$dest"; binaer "$plat" "$ver" "$var" > "$dest/pgwire-recorder"
    else
      hat "${AT_IMG_BAU_ROT-}" "$plat" && exit 1
      mkdir -p "$S/i-$tag"; os=${plat%/*}; a=${plat#*/}
      echo "$ver" > "$S/i-$tag/ver"
      hat "${AT_IMG_OS-}" "$plat" && os=windows
      hat "${AT_IMG_ARCH-}" "$plat" && a=riscv64
      echo "$os" > "$S/i-$tag/os"; echo "$a" > "$S/i-$tag/arch"
      binaer "$plat" "$ver" v0 > "$S/i-$tag/bin"
      hat "${AT_IMG_SUMME-}" "$plat" && echo x >> "$S/i-$tag/bin"
      :
    fi ;;
  image)
    [ -d "$S/i-${*: -1}" ] || exit 1
    case "$*" in *.Os*) cat "$S/i-${*: -1}/os" ;; *.Architecture*) cat "$S/i-${*: -1}/arch" ;; esac ;;
  rmi) for t; do rm -rf "${S:?}/i-$t"; done ;;
  network)
    if [ "$1" = create ]; then : > "$S/n-${*: -1}"; else rm -f "${S:?}/n-${2:?}"; fi ;;
  create | run)
    name= img= cmd= alle="$*"
    while [ $# -gt 0 ]; do
      case $1 in --name) name=$2; shift ;; esac
      shift
    done
    set -- $alle
    img="$(bild "$@")"
    if [ "$sub" = create ]; then : > "$S/c-$name"; echo "$img" > "$S/c-$name.bild"; exit 0; fi
    case " $alle " in
      *" version "*)
        [ -n "$img" ] || exit 1
        [ -z "${AT_VERSION_ROT-}" ] || exit 1
        if [ -n "${AT_VERSION_AUS-}" ]; then printf '%s\n' "$AT_VERSION_AUS"
        elif [ -n "${AT_VERSION_OHNE_LF-}" ]; then printf 'pgwire-recorder %s' "$(cat "$S/i-$img/ver")"
        else echo "pgwire-recorder $(cat "$S/i-$img/ver")"; fi ;;
      *" psql "*)
        hang psql
        [ -z "${AT_PSQL_ROT-}" ] || exit 2
        z=${alle#* -h }; z=${z%% *}
        case $z in
          *-rec) if [ -n "${AT_SICHT_A_LEER-}" ] || [ -n "${AT_SICHT_LEER-}" ]; then :; elif [ -n "${AT_SICHT_OHNE_EINS-}" ]; then echo " x"; else printf ' eins \n------\n    1\n(1 row)\n'; fi ;;
          *-rep)
            [ -e "$S/c-$z.aufz" ] || exit 2
            [ -z "${AT_SICHT_LEER-}" ] || exit 0
            if [ -n "${AT_SICHT_B_ABW-}" ]; then printf ' eins \n------\n    2\n(1 row)\n'; else printf ' eins \n------\n    1\n(1 row)\n'; fi ;;
        esac ;;
      *" -d "*) : > "$S/c-$name" ;;
    esac ;;
  cp)
    case "$1" in
      *:/pgwire-recorder) n=${1%%:*}; cp "$S/i-$(cat "$S/c-$n.bild")/bin" "$2" ;;
      *:/tmp/aufzeichnung.yaml)
        [ -z "${AT_CP_ROT-}" ] || exit 1
        if [ -n "${AT_AUFZ_LEER-}" ]; then : > "$2"; else echo aufzeichnung > "$2"; fi ;;
      *) n=${2%%:*}; cp "$1" "$S/c-$n.aufz" ;;
    esac ;;
  rm)
    for n in "${@:2}"; do
      if [[ "$n" == *-pg && -n "${AT_PG_BLEIBT-}" ]]; then continue; fi
      rm -f "${S:?}/c-${n:?}" "$S/c-$n.bild" "$S/c-$n.aufz"
    done ;;
  container) [ -e "$S/c-$2" ] ;;
  inspect)
    n=${*: -1}; [ -e "$S/c-$n" ] || exit 1
    case $n in *-rec) echo "${AT_EXIT_REC:-0}" ;; *-rep) echo "${AT_EXIT_REP:-0}" ;; esac ;;
  stop) echo "$1" ;;
  start) : ;;
  logs) [ -n "${AT_KEIN_START-}" ] || echo 'level=INFO msg="gestartet"' ;;
  exec) [ -z "${AT_NICHT_BEREIT-}" ] ;;
  version) echo "${AT_SERVER_ARCH:-amd64}" ;;
esac
ATTRAPPE
printf '#!/usr/bin/env bash\n:\n' > "$bin/sleep"
chmod +x "$bin/docker" "$bin/sleep"

make_variable() {
  { make -s --no-print-directory -C "$1" -pq help 2> /dev/null || true; } \
    | awk -v n="$2" '$1 == n && ($2 == ":=" || $2 == "=") { for (i = 3; i <= NF; i++) print $i }'
}
ERF="$(make_variable "$wurzel" BINARY_PLATTFORMEN | tr '\n' ' ')"; ERF="${ERF% }"
IMG="$(make_variable "$wurzel" IMAGE_PLATTFORMEN | tr '\n' ' ')"; IMG="${IMG% }"
refbild="$(make_variable "$wurzel" POSTGRES_IMAGE)"
FIXTURE_ZAEHLER=0

# fixture — legt eine Kopie an und setzt w auf ihre Wurzel.
fixture() {
  FIXTURE_ZAEHLER=$((FIXTURE_ZAEHLER + 1))
  w="$arbeit/baum-$FIXTURE_ZAEHLER"
  mkdir -p "$w/harness/mk" "$w/tools/harness" "$w/.harness/state"
  cp "$wurzel/Makefile" "$wurzel/Dockerfile" "$w/"
  cp "$wurzel/harness/mk/integration.mk" "$wurzel/harness/mk/container.mk" "$w/harness/mk/"
  cp "$wurzel/tools/harness/binaries.sh" "$wurzel/tools/harness/run-container.sh" "$wurzel/tools/harness/container-eingabe.sh" "$w/tools/harness/"
  printf '#!/usr/bin/env bash\n: > marke-gegenprobe\n' > "$w/tools/harness/container-gegenprobe.sh"
}

# lauf [VAR=wert…] -- <befehl>… — führt den Befehl in $w aus; setzt code, Ausgabe in $aus.
lauf() {
  local e=()
  while [ "$1" != -- ]; do e+=("$1"); shift; done
  shift
  rm -rf "${st:?}" "${tmp:?}"; mkdir -p "$st" "$tmp"; : > "$log"
  set +e
  {
    (cd "$w" && env -u MAKELEVEL -u MFLAGS MAKEFLAGS="${mf:-}" TMPDIR="$tmp" PATH="$bin:$PATH" AT_LOG="$log" AT_STATE="$st" \
      ${e[@]+"${e[@]}"} "$@") > "$aus" 2>&1
    code=$?
  } 2> /dev/null
  set -e
}
ma=()
# gate [VAR=wert…] [-- <make-flag>…] — make test-container mit den Argumenten in $ma.
gate() {
  local e=()
  while [ $# -gt 0 ] && [ "$1" != -- ]; do e+=("$1"); shift; done
  if [ $# -gt 0 ]; then shift; fi
  lauf ${e[@]+"${e[@]}"} -- make -s --no-print-directory "$@" test-container DOCKER="$bin/docker" ${ma[@]+"${ma[@]}"}
}
# werkz [VAR=wert…] [-- <make-flag>…] — make binaries.
werkz() {
  local e=()
  while [ $# -gt 0 ] && [ "$1" != -- ]; do e+=("$1"); shift; done
  if [ $# -gt 0 ]; then shift; fi
  lauf ${e[@]+"${e[@]}"} -- make -s --no-print-directory "$@" binaries DOCKER="$bin/docker" ${ma[@]+"${ma[@]}"}
}
in_aus() { grep -qF -- "$1" "$aus"; }
nicht_aus() { ! grep -qF -- "$1" "$aus"; }
anzahl() { grep -cE -- "$1" "$log" || true; }
zeile() { grep -nE -- "$1" "$log" | head -1 | cut -d: -f1; }
dateien() { (cd "$w" && find . -type f | sort); }
ohne() { local x r=""; for x in $1; do if [ "$x" != "$2" ]; then r+="$x "; fi; done; echo "${r% }"; }
dt() { echo "${1//\//-}"; }
bname() { local n="pgwire-recorder_${3:-dev}_$1_$2"; if [ "$1" = windows ]; then n+=".exe"; fi; echo "$n"; }

sel() {
  local p
  if [ -z "${nur:-}" ]; then return 0; fi
  for p in $nur; do if [[ "$1" == "$p"* ]]; then return 0; fi; done
  return 1
}
ist() { if sel "$1"; then eval "$2" 2> /dev/null || echo "$1"; fi; return 0; }
# lauft <Fall oder Gruppe>… — läuft einer von ihnen (oder einer ihrer Fälle)?
lauft() {
  local f p
  for f; do
    if sel "$f"; then return 0; fi
    for p in ${nur:-}; do if [[ "$p" == "$f"* ]]; then return 0; fi; done
  done
  return 1
}

# efall <Name> <Text> — Fehler der Eingabe: rot, kein Aufruf von docker, der Text in der Ausgabe.
efall() {
  if sel "$1"; then
    gate
    { [ "$code" -ne 0 ] && [ ! -s "$log" ] && in_aus "$2"; } || echo "$1"
  fi
}
# gruen — ein grüner Lauf mit Meldung.
gruen() { [ "$code" -eq 0 ] && in_aus "test-container: grün"; }
# rot <Text> — rot mit dem Text.
rot() { [ "$code" -ne 0 ] && in_aus "ROT — $1"; }
# runtime_gebaut — hat das Gate ein Image gebaut?
runtime_gebaut() { [ "$(anzahl "^build .*--target runtime")" -gt 0 ]; }
# aufgeraeumt — Container, Netz, Tags und Temp-Verzeichnisse des Laufs sind entfernt.
res_container() { ! ls "$st" | grep -q '^c-'; }
res_netz() { ! ls "$st" | grep -q '^n-'; }
res_tags() { ! ls "$st" | grep -q '^i-'; }
res_werk() { ! ls -A "$tmp" | grep -q '^pgr-container\.'; }
res_stage() { [ -z "$(ls -A "$tmp" | grep -v '^pgr-container\.')" ]; }
ende() { # <Name des Endes> — vier Fälle
  ist "p9-$1-container" res_container
  ist "p9-$1-netz" res_netz
  ist "p9-$1-tags" res_tags
  ist "p9-$1-werk" res_werk
  ist "p9-$1-stage" res_stage
}
# signalfall <INT|TERM> <build|psql> [gruppe|skript] [werkzeug] — das Skript (mit `werkzeug` binaries.sh) hängt in der
# Attrappe; das Signal geht an das Skript, mit `gruppe` an seine Prozessgruppe (die Sitzung von setsid).
signalfall() {
  local marke="$arbeit/hangt" pidf="$arbeit/pid" wp
  rm -f "$marke" "$pidf" "$gesendet"
  (
    limit=$((SECONDS + 20))
    until [ -e "$marke" ] && [ -s "$pidf" ]; do
      if [ "$SECONDS" -ge "$limit" ]; then exit 0; fi
      if [ -s "$pidf" ] && ! kill -0 "$(cat "$pidf")" 2> /dev/null; then exit 0; fi
      /bin/sleep 0.05
    done
    if [ "${3-}" = gruppe ]; then kill -"$1" -- "-$(cat "$pidf")"; else kill -"$1" "$(cat "$pidf")"; fi
    : > "$gesendet"
  ) &
  wp=$!
  lauf VERSION=dev BINARY_PLATTFORMEN="$ERF" IMAGE_PLATTFORMEN="$IMG" POSTGRES_IMAGE="$refbild" DOCKER="$bin/docker" DIST=dist \
    AT_HANGT="$marke" AT_HANG_ORT="$2" -- \
    setsid bash -c 'echo $$ > "$1"; shift; exec "$@"' x "$pidf" bash "tools/harness/$([ "${4-}" = werkzeug ] && echo binaries.sh || echo run-container.sh)"
  wait "$wp" || true
}
sums_zeilen() { (cd "$w/${1:-dist}" && cat SHA256SUMS); }
namen_soll() { # <Version> — die sechs Namen, nach Name sortiert
  local p
  for p in $ERF; do bname "${p%/*}" "${p#*/}" "$1"; done | LC_ALL=C sort
}

# pruefe <baum> — schreibt die Namen der Fälle, die nicht bestehen (einen je Zeile); $nur wählt aus.
pruefe() {
  w="$1"
  local p a v n text sg ort modus ziel ende_name k1 k2
  ma=()

  # Punkt 1: Plattformen, vor dem ersten Aufruf von docker.
  ma=(BINARY_PLATTFORMEN=); efall p1a-bin-leer "die Liste der Binaries ist leer"
  ma=(IMAGE_PLATTFORMEN=); efall p1a-img-leer "die Liste der Images ist leer"
  ma=("BINARY_PLATTFORMEN=$ERF linux"); efall p1b1-ohne-slash "Eintrag ohne '/' in der Liste der Binaries: linux"
  ma=("BINARY_PLATTFORMEN=$ERF /amd64"); efall p1b2-os-leer "leeres Betriebssystem im Eintrag /amd64"
  ma=("BINARY_PLATTFORMEN=$ERF linux/"); efall p1b3-arch-leer "leere Architektur im Eintrag linux/"
  ma=("BINARY_PLATTFORMEN=$ERF freebsd/amd64"); efall p1c1-os-unbekannt "Betriebssystem 'freebsd' unbekannt im Eintrag freebsd/amd64"
  ma=("BINARY_PLATTFORMEN=$ERF linux/riscv64"); efall p1c2-arch-unbekannt "Architektur 'riscv64' unbekannt im Eintrag linux/riscv64"
  ma=("BINARY_PLATTFORMEN=$ERF linux/amd64"); efall p1d-doppelt "Eintrag doppelt in der Liste der Binaries: linux/amd64"
  for p in $ERF; do
    ma=("BINARY_PLATTFORMEN=$(ohne "$ERF" "$p")"); efall "p1e-bin-$(dt "$p")" "$p fehlt in der Liste der Binaries"
  done
  for p in $IMG; do
    ma=("IMAGE_PLATTFORMEN=$(ohne "$IMG" "$p")"); efall "p1e-img-$(dt "$p")" "$p fehlt in der Liste der Images"
  done
  ma=("BINARY_PLATTFORMEN=$(ohne "$ERF" linux/arm64)"); efall p1f-image-nicht-in-binaries "Image-Plattform linux/arm64 steht nicht in der Liste der Binaries"
  # dieselbe Prüfung vor dem ersten Aufruf von docker bei `make binaries`
  if sel p1e-binaries; then
    ma=("BINARY_PLATTFORMEN=$(ohne "$ERF" darwin/arm64)"); werkz
    { [ "$code" -ne 0 ] && [ ! -s "$log" ] && in_aus "darwin/arm64 fehlt in der Liste der Binaries"; } || echo p1e-binaries
  fi

  # Punkt 2: VERSION, vor dem ersten Aufruf von docker.
  ma=(VERSION=); efall p2-leer "VERSION '' ist ungültig"
  ma=("VERSION=a b"); efall p2-leerzeichen "VERSION 'a b' ist ungültig"
  ma=('VERSION=a$$b'); efall p2-dollar 'VERSION '"'"'a$b'"'"' ist ungültig'
  ma=("VERSION=a'b"); efall p2-einfaches-anfuehrungszeichen "VERSION 'a'b' ist ungültig"
  ma=('VERSION=a"b'); efall p2-doppeltes-anfuehrungszeichen "VERSION 'a\"b' ist ungültig"
  ma=(VERSION=-x); efall p2-fuehrendes-minus "VERSION '-x' ist ungültig"
  if sel p2-binaries; then
    ma=(VERSION=); werkz
    { [ "$code" -ne 0 ] && [ ! -s "$log" ] && in_aus "VERSION '' ist ungültig"; } || echo p2-binaries
  fi
  # Eine gültige Version geht in den Build und in den Dateinamen ein.
  if lauft p2-gueltig; then
    ma=("VERSION=1.0.0-rc.1+b_2"); gate
    ist p2-gueltig 'gruen && [ "$(anzahl "^build .*--target dist .*--build-arg VERSION=1.0.0-rc.1\+b_2 ")" -eq 12 ]'
  fi

  # Punkte 3 bis 7, 9, 10: ein grüner Lauf des Gates.
  if lauft p4-gruen p5-gruen p6-gruen p6-tag p7-gruen p7-intern p7-kein p7-pg-vor-replay p7-cp p10 p9-gruen p9-kennung; then
    local vor nach
    ma=()
    vor="$(dateien)"
    gate
    nach="$(dateien)"
    ist p4-gruen 'gruen && [ "$(anzahl "^build .*--target dist")" -eq 12 ] && [ "$(anzahl "^build .*--target dist .*--no-cache-filter build --no-cache-filter dist")" -eq 6 ]'
    ist p5-gruen 'in_aus "Punkt 5 grün — Kopf von 6 Binaries"'
    ist p6-gruen '[ "$(anzahl "^build .*--target runtime --platform linux/(amd64|arm64) ")" -eq 2 ] && in_aus "Punkt 6 grün"'
    ist p6-tag '[ "$(anzahl "^build .*--target runtime .* -t pgr-ct-[0-9]+:linux-(amd64|arm64) ")" -eq 2 ]'
    ist p7-gruen 'in_aus "Punkt 7 grün — Smoke auf linux/amd64" && [ "$(anzahl "^run --rm --name [^ ]*-version ")" -eq 1 ]'
    ist p7-intern '[ "$(anzahl "^network create --internal pgr-ct-[0-9]+-net$")" -eq 1 ]'
    ist p7-kein-terminal '! grep -E "^(run|create) " "$log" | grep -qE " -[a-zA-Z]*t[a-zA-Z]*( |$)| --tty"'
    ist p7-kein-user '! grep -E "^(run|create) " "$log" | grep -qE " --user| -u "'
    ist p7-kein-volume '! grep -E "^(run|create) " "$log" | grep -qE " -v | --volume"'
    ist p7-kein-mount '! grep -E "^(run|create) " "$log" | grep -qE " --mount"'
    ist p7-pg-vor-replay '[ "$(zeile "^rm -f pgr-ct-[0-9]+-pg$")" -lt "$(zeile "^create --name pgr-ct-[0-9]+-rep ")" ] && [ "$(zeile "^rm -f pgr-ct-[0-9]+-pg$")" -lt "$(zeile "^start pgr-ct-[0-9]+-rep$")" ]'
    ist p7-cp-entnehmen '[ "$(anzahl "^cp pgr-ct-[0-9]+-rec:/tmp/aufzeichnung.yaml ")" -eq 1 ]'
    ist p7-cp-hinein '[ "$(anzahl "^cp [^ ]*/aufzeichnung.yaml pgr-ct-[0-9]+-rep:/tmp/aufzeichnung.yaml$")" -eq 1 ] && [ "$(zeile "^cp [^ ]*/aufzeichnung.yaml pgr-ct-[0-9]+-rep:")" -lt "$(zeile "^start pgr-ct-[0-9]+-rep$")" ]'
    ist p10-arbeitsbaum '[ "$vor" = "$nach" ]'
    ist p9-kennung-ziffern '! grep -qE "pgr-ct-[^0-9]" "$log"'
    ist p9-gruen-container res_container
    ist p9-gruen-netz res_netz
    ist p9-gruen-tags res_tags
    ist p9-gruen-werk res_werk
    ist p9-gruen-stage res_stage
    k1="$(grep '^network create' "$log" || true)"
    gate
    k2="$(grep '^network create' "$log" || true)"
    ist p9-kennung-eindeutig '[ -n "$k1" ] && [ "$k1" != "$k2" ]'
  fi

  # Punkt 3: das Werkzeug `make binaries`.
  if lauft p3; then
    ma=(VERSION=7x7 DIST=out); werkz
    ist p3-gruen '[ "$code" -eq 0 ]'
    ist p3-namen '[ "$(cd "$w/out" && ls -A | LC_ALL=C sort | tr "\n" " ")" = "SHA256SUMS $(namen_soll 7x7 | tr "\n" " ")" ]'
    ist p3-exe '[ -f "$w/out/pgwire-recorder_7x7_windows_amd64.exe" ] && [ -f "$w/out/pgwire-recorder_7x7_windows_arm64.exe" ] && [ ! -e "$w/out/pgwire-recorder_7x7_linux_amd64.exe" ]'
    ist p3-version-im-build '[ "$(anzahl "^build .*--target dist .*--build-arg VERSION=7x7 ")" -eq 6 ]'
    ist p3-nur-binaries '[ "$(anzahl "^build ")" -eq 6 ] && [ "$(anzahl "no-cache-filter")" -eq 0 ] && [ "$(anzahl "^(run|create|network) ")" -eq 0 ]'
    ist p3-sums-format '[ "$(sums_zeilen out | grep -cE "^[0-9a-f]{64}  pgwire-recorder_7x7_[a-z]+_[a-z0-9]+(\.exe)?$")" -eq 6 ]'
    ist p3-sums-je-binary '[ "$(sums_zeilen out | wc -l)" -eq 6 ] && [ "$(sums_zeilen out | awk "{print \$2}" | LC_ALL=C sort | tr "\n" " ")" = "$(namen_soll 7x7 | tr "\n" " ")" ]'
    ist p3-sums-sortiert 'sums_zeilen out | awk "{print \$2}" | LC_ALL=C sort -c'
    ist p3-sums-lf '! grep -q $'"'"'\r'"'"' "$w/out/SHA256SUMS" && [ "$(tail -c1 "$w/out/SHA256SUMS" | od -An -c | tr -d " ")" = "\\n" ]'
    ist p3-sums-pruefbar '(cd "$w/out" && sha256sum -c SHA256SUMS > /dev/null 2>&1)'
    ma=(DIST=neu/dist); werkz
    ist p3-dist-fehlt '[ "$code" -eq 0 ] && [ -f "$w/neu/dist/SHA256SUMS" ]'
    mkdir -p "$w/leer"; ma=(DIST=leer); werkz
    ist p3-dist-leer '[ "$code" -eq 0 ] && [ -f "$w/leer/SHA256SUMS" ]'
    mkdir -p "$w/voll"; echo x > "$w/voll/frueher"; ma=(DIST=voll); werkz
    ist p3-dist-nicht-leer '[ "$code" -ne 0 ] && [ ! -s "$log" ] && in_aus "DIST '"'"'voll'"'"' ist nicht leer" && [ "$(ls "$w/voll")" = frueher ]'
    ma=(DIST=); werkz
    ist p3-dist-leerer-wert '[ "$code" -ne 0 ] && [ ! -s "$log" ] && in_aus "DIST ist leer"'
    echo x > "$w/datei"; ma=(DIST=datei); werkz
    ist p3-dist-datei '[ "$code" -ne 0 ] && [ ! -s "$log" ] && in_aus "DIST '"'"'datei'"'"' ist nicht leer"'
    ma=(DIST=zweimal); werkz; werkz
    ist p3-zweiter-lauf '[ "$code" -ne 0 ] && [ ! -s "$log" ] && in_aus "ist nicht leer"'
    for p in linux/amd64 windows/arm64; do
      ma=(DIST=rot); rm -rf "$w/rot"; werkz AT_BAU_ROT="$p"
      ist "p3-rot-nennt-$(dt "$p")" "[ \"\$code\" -ne 0 ] && in_aus 'binaries: ROT — Build für $p gescheitert'"
    done
    ma=(DIST=rot); rm -rf "$w/rot"; werkz AT_BAU_ROT=windows/amd64
    ist p3-rot-ohne-sums '[ ! -e "$w/rot/SHA256SUMS" ]'
    ist p3-rot-keine-teilmenge '[ -z "$(ls -A "$w/rot" | grep -v "^SHA256SUMS$")" ]'
  fi

  # Punkt 4: zwei Builds, dieselbe Summe.
  for p in $ERF; do
    if sel "p4-abweichung-$(dt "$p")"; then
      ma=(); gate AT_ABWEICH="$p"
      n="$(bname "${p%/*}" "${p#*/}")"
      { [ "$code" -ne 0 ] && grep -qxF "test-container: ROT — die Summen zweier Builds weichen ab: $n" "$aus" && ! runtime_gebaut; } || echo "p4-abweichung-$(dt "$p")"
    fi
  done
  if lauft p4-namen; then
    ma=(); gate AT_ABWEICH="linux/amd64 darwin/arm64"
    ist p4-namen '[ "$code" -ne 0 ] && grep -qxF "test-container: ROT — die Summen zweier Builds weichen ab: pgwire-recorder_dev_darwin_arm64 pgwire-recorder_dev_linux_amd64" "$aus"'
  fi
  if lauft p4-cache; then
    ma=(); gate AT_NICHT_REPRO=1
    ist p4-cache-umgehen 'rot "die Summen zweier Builds weichen ab:"'
  fi

  # Punkt 5: Kopf der Binaries, je Plattform mit der Architektur der anderen und mit dem Kopf eines anderen Formats.
  for p in $ERF; do
    local os="${p%/*}" ar="${p#*/}" wert fremd werta
    case "$os/$ar" in
      linux/amd64) werta="ELF-Maschine b700" ;; linux/arm64) werta="ELF-Maschine 3e00" ;;
      darwin/amd64) werta="Mach-O-CPU-Typ 0c000001" ;; darwin/arm64) werta="Mach-O-CPU-Typ 07000001" ;;
      windows/amd64) werta="PE-Maschine 64aa" ;; windows/arm64) werta="PE-Maschine 6486" ;;
    esac
    case "$os" in linux) fremd=macho; wert="ELF-Magic cffaedfe" ;; darwin) fremd=pe; wert="Mach-O-Magic 4d5a0000" ;; windows) fremd=elf; wert="PE-Magic 7f45" ;; esac
    if sel "p5-arch-$(dt "$p")"; then
      ma=(); gate AT_KOPF="$p=arch"
      rot "Kopf des Binaries für $p: $werta" || echo "p5-arch-$(dt "$p")"
    fi
    if sel "p5-format-$(dt "$p")"; then
      ma=(); gate AT_KOPF="$p=$fremd"
      rot "Kopf des Binaries für $p: $wert" || echo "p5-format-$(dt "$p")"
    fi
  done

  # Punkt 6: Image je Plattform.
  for p in $IMG; do
    for a in os:Os:Betriebssystem arch:ARCH:Architektur summe:SUMME:Summe; do
      v="${a%%:*}"
      if sel "p6-$v-$(dt "$p")"; then
        ma=()
        case "$v" in
          os) gate AT_IMG_OS="$p"; rot "Image $p nennt das Betriebssystem 'windows'" || echo "p6-os-$(dt "$p")" ;;
          arch) gate AT_IMG_ARCH="$p"; rot "Image $p nennt die Architektur 'riscv64'" || echo "p6-arch-$(dt "$p")" ;;
          summe) gate AT_IMG_SUMME="$p"; rot "Image $p trägt ein Binary mit der Summe" || echo "p6-summe-$(dt "$p")" ;;
        esac
      fi
    done
  done
  if lauft p6-bau; then
    ma=(); gate AT_IMG_BAU_ROT=linux/arm64
    ist p6-bau-rot 'rot "Image für linux/arm64 nicht zu bauen"'
  fi
  ist p6-runtime-ohne-run '! awk "/^FROM .* AS runtime/ {r=1; next} /^FROM / {r=0} r && /^RUN / {f=1} END {exit !f}" "$w/Dockerfile"'

  # Punkt 7: Smoke. gfall <Fall> '<AT_…=…>…' '<Bedingung>' — ein Lauf nur, wenn der Fall gewählt ist.
  gfall() {
    if sel "$1"; then
      ma=()
      gate "$2"
      eval "$3" 2> /dev/null || echo "$1"
    fi
  }
  gfall p7-server-arm64 AT_SERVER_ARCH=arm64 'in_aus "Smoke auf linux/arm64" && [ "$(anzahl "^run --rm --name [^ ]*-version pgr-ct-[0-9]+:linux-arm64 version$")" -eq 1 ]'
  gfall p7-server-fremd AT_SERVER_ARCH=riscv64 'rot "Plattform riscv64 des Docker-Servers ist weder amd64 noch arm64" && ! grep -q "^run " "$log"'
  gfall p7a-ausgabe-abweichend "AT_VERSION_AUS=pgwire-recorder_x" 'rot "version gibt ${q}pgwire-recorder_x${q} aus statt ${q}pgwire-recorder dev${q} und LF"'
  gfall p7a-zusatzzeile "AT_VERSION_AUS=pgwire-recorder dev"$'\n'"zweite Zeile" '[ "$code" -ne 0 ] && in_aus "version gibt ${q}pgwire-recorder dev"'
  gfall p7a-ohne-lf AT_VERSION_OHNE_LF=1 'rot "version gibt ${q}pgwire-recorder dev${q} aus statt"'
  gfall p7a-exit AT_VERSION_ROT=1 'rot "version endet mit einem Exit-Code ungleich 0"'
  gfall p7b-pg-nicht-bereit AT_NICHT_BEREIT=1 'rot "PostgreSQL nicht bereit"'
  gfall p7b-kein-start AT_KEIN_START=1 '[ "$code" -ne 0 ] && in_aus "-rec meldet keinen Start"'
  gfall p7b-sicht-a-leer AT_SICHT_A_LEER=1 'rot "Sicht A ist leer"'
  gfall p7b-beide-leer AT_SICHT_LEER=1 'rot "Sicht A ist leer"'
  gfall p7b-ohne-eins AT_SICHT_OHNE_EINS=1 'rot "Sicht A nennt ${q}eins${q} nicht"'
  gfall p7b-anfrage-rot AT_PSQL_ROT=1 'rot "Anfrage an"'
  gfall p7b-exit-rec AT_EXIT_REC=1 'rot "Exit-Code 1 des Containers mit record (erwartet 0)"'
  gfall p7b-cp-rot AT_CP_ROT=1 'rot "die Aufzeichnung ist dem Container nicht zu entnehmen"'
  gfall p7b-aufzeichnung-leer AT_AUFZ_LEER=1 'rot "die Aufzeichnung ist leer"'
  gfall p7c-pg-lebt AT_PG_BLEIBT=1 'rot "der PostgreSQL-Container lebt beim Start des Wiedergebens noch" && ! grep -q "^create --name [^ ]*-rep" "$log"'
  gfall p7c-sicht-ungleich AT_SICHT_B_ABW=1 'rot "Sicht B ist ungleich Sicht A"'
  gfall p7c-exit-rep AT_EXIT_REP=1 'rot "Exit-Code 1 des Containers mit replay (erwartet 0)"'

  # Punkt 8: Reihenfolge der Gründe.
  if lauft p8; then
    ma=("BINARY_PLATTFORMEN=$(ohne "$ERF" linux/arm64)"); gate AT_KOPF="linux/amd64=arch"
    ist p8-1-vor-5 'in_aus "Eingabe:" && nicht_aus "Kopf des Binaries" && [ ! -s "$log" ]'
    ma=(); gate AT_ABWEICH=linux/amd64 AT_IMG_SUMME=linux/amd64
    ist p8-4-vor-6 'in_aus "ROT — die Summen zweier Builds weichen ab" && nicht_aus "Image linux/amd64" && ! runtime_gebaut'
  fi

  # Punkt 9: Ende jedes Falls räumt auf (grün oben).
  if lauft p9-rot4 p9-rot5 p9-rot6 p9-rot7; then
    ma=(); gate AT_ABWEICH=linux/amd64; ende rot4
    gate AT_KOPF="linux/amd64=arch"; ende rot5
    gate AT_IMG_SUMME=linux/amd64; ende rot6
    gate AT_SICHT_B_ABW=1; ende rot7
  fi
  for a in INT:build:skript:gate:sigINT-build INT:psql:skript:gate:sigINT-psql TERM:build:skript:gate:sigTERM-build TERM:psql:skript:gate:sigTERM-psql INT:build:gruppe:gate:sigINTg-build INT:build:skript:werkzeug:binINT TERM:build:skript:werkzeug:binTERM; do
    IFS=: read -r sg ort modus ziel ende_name <<< "$a"
    if lauft "p9-$ende_name"; then
      signalfall "$sg" "$ort" "$modus" "$ziel"
      ist "p9-$ende_name-beendet" '[ -e "$gesendet" ] && [ "$code" -ne 0 ] && nicht_aus "test-container: grün" && nicht_aus "binaries: 6"'
      ist "p9-$ende_name-signal" '[ -s "$st/term" ]'
      ende "$ende_name"
    fi
  done

  # Punkt 11: Aufrufarten des Treibers.
  local ziel flag f k
  for ziel in gate werkz; do
    for f in -n -q -t -i "-n -i" "-q -i" "-t -i"; do
      k="${f//-/}"; k="${k// /}"
      if lauft "p11-$ziel-$k"; then
        ma=(); rm -rf "$w/dist"
        if [ "$f" = -i ]; then "$ziel" AT_BAU_ROT=linux/amd64 AT_ABWEICH=linux/amd64 -- $f; else "$ziel" -- $f; fi
        if [ "$f" = -i ]; then
          ist "p11-$ziel-i" '[ "$code" -eq 0 ] && [ "$(anzahl "^build ")" -gt 0 ] && in_aus "ROT"'
        else
          ist "p11-$ziel-$k" '[ ! -s "$log" ] && [ ! -e "$w/dist" ]'
        fi
      fi
    done
    for f in n q t i; do
      if lauft "p11-$ziel-makeflags-$f"; then
        ma=(); rm -rf "$w/dist"
        if [ "$f" = i ]; then mf=i "$ziel" AT_BAU_ROT=linux/amd64 AT_ABWEICH=linux/amd64; else mf=$f "$ziel"; fi
        if [ "$f" = i ]; then
          ist "p11-$ziel-makeflags-i" '[ "$code" -eq 0 ] && [ "$(anzahl "^build ")" -gt 0 ] && in_aus "ROT"'
        else
          ist "p11-$ziel-makeflags-$f" '[ ! -s "$log" ] && [ ! -e "$w/dist" ]'
        fi
        mf=""
      fi
    done
  done
  if lauft p11-gegenprobe; then
    ma=(); lauf -- make -n -s --no-print-directory test-container-gegenprobe
    ist p11-gegenprobe-rezept '[ "$(cat "$aus")" = "bash tools/harness/container-gegenprobe.sh" ] && [ ! -e "$w/marke-gegenprobe" ]'
  fi

  # Punkt 12: beide Ziele hängen an GATE_CHECKS.
  ist p12-gate-checks-ziel '[[ " $(make_variable "$w" GATE_CHECKS | tr "\n" " ") " == *" test-container "* ]]'
  ist p12-gate-checks-gegenprobe '[[ " $(make_variable "$w" GATE_CHECKS | tr "\n" " ") " == *" test-container-gegenprobe "* ]]'
}
q="'"

# mutiere <datei> <alt> <neu> — ersetzt das erste Vorkommen (wörtlich) in der Kopie; mit `%<marke>`
# als alt wird die Zeile ersetzt, die auf `# <marke>` endet.
mutiere() {
  local f="$w/$1" inhalt neu
  if [[ "$2" == %* ]]; then
    awk -v m="# ${2#%}" '{ l = length($0); k = length(m); if (l >= k && substr($0, l - k + 1) == m) print ":"; else print }' "$f" > "$f.neu"
    mv "$f.neu" "$f"
    return
  fi
  inhalt="$(cat "$f"; printf x)"; inhalt="${inhalt%x}"
  neu="${inhalt/"$2"/"$3"}"
  if [ "$neu" = "$inhalt" ]; then melde "Mutation wirkungslos in $1: $2"; fi
  printf '%s' "$neu" > "$f"
}
# mut <name> <datei> <alt> <neu> <nur> <erwartet rot>… — bei der Mutation sind genau diese Fälle rot.
mut() {
  local name="$1" datei="$2" alt="$3" neu="$4" soll ist
  if [ -n "${GP_MUT-}" ] && ! [[ "$1" =~ $GP_MUT ]]; then return 0; fi
  if [[ "$alt" == %* ]]; then set -- "$1" "$2" "$3" "" "${@:4}"; neu=""; fi
  nur="$5"
  shift 5
  fixture
  mutiere "$datei" "$alt" "$neu"
  ist="$(pruefe "$w" | sort -u | tr '\n' ' ')"
  soll="$(printf '%s\n' "$@" | sort -u | tr '\n' ' ')"
  if [ -n "${GP_ENTDECKE-}" ]; then echo "$name => ${ist:-nichts}" >&2; fi
  if [ "$ist" != "$soll" ]; then melde "Mutation '$name': erwartet rot: ${soll:-nichts}; rot sind: ${ist:-nichts}"; fi
}

L_=tools/harness/container-eingabe.sh B_=tools/harness/binaries.sh G_=tools/harness/run-container.sh F_=harness/mk/container.mk

# --- Die unveränderte Kopie besteht jeden Fall; beide Ziele hängen im Arbeitsbaum an GATE_CHECKS ---
nur=""
fixture
ist="$(pruefe "$w" | sort -u | tr '\n' ' ')"
if [ -n "$ist" ]; then melde "die unveränderte Kopie besteht nicht: $ist"; fi
gates=" $(make_variable "$wurzel" GATE_CHECKS | tr '\n' ' ') "
for g in test-container test-container-gegenprobe; do
  if [[ "$gates" != *" $g "* ]]; then melde "$g hängt im Arbeitsbaum nicht an GATE_CHECKS"; fi
done

#MUTATIONEN
# --- Punkt 1: Plattformen ---
P1N="p1a p1b p1c p1d p1e p1f p1-"
mut p1a "$L_" %P1a p1a p1a-bin-leer p1a-img-leer
mut p1b1 "$L_" %P1b1 p1b1 p1b1-ohne-slash
mut p1b2 "$L_" %P1b2 p1b2 p1b2-os-leer
mut p1b3 "$L_" %P1b3 p1b3 p1b3-arch-leer
mut p1c1 "$L_" %P1c1 p1c1 p1c1-os-unbekannt
mut p1c2 "$L_" %P1c2 p1c2 p1c2-arch-unbekannt
mut p1d "$L_" %P1d p1d p1d-doppelt
mut p1f "$L_" %P1f p1f p1f-image-nicht-in-binaries
for p in $ERF; do
  if [ "$p" = darwin/arm64 ]; then x=p1e-binaries; else x=; fi
  mut "p1e-bin-$p" "$L_" "ERF_BINARIES=\"$ERF\"" "ERF_BINARIES=\"$(ohne "$ERF" "$p")\"" "p1e-bin-$(dt "$p") $x" "p1e-bin-$(dt "$p")" $x
done
for p in $IMG; do
  mut "p1e-img-$p" "$L_" "ERF_IMAGES=\"$IMG\"" "ERF_IMAGES=\"$(ohne "$IMG" "$p")\"" "p1e-img-$(dt "$p")" "p1e-img-$(dt "$p")"
done
efaelle="p1a-bin-leer p1a-img-leer p1b1-ohne-slash p1b2-os-leer p1b3-arch-leer p1c1-os-unbekannt p1c2-arch-unbekannt p1d-doppelt p1f-image-nicht-in-binaries p2-leer p2-leerzeichen p2-dollar p2-einfaches-anfuehrungszeichen p2-doppeltes-anfuehrungszeichen p2-fuehrendes-minus"
for p in $ERF; do efaelle+=" p1e-bin-$(dt "$p")"; done
for p in $IMG; do efaelle+=" p1e-img-$(dt "$p")"; done
mut p8-vor-docker "$G_" %P8vor-docker "p1a p1b p1c p1d p1e p1f p2 p8" $efaelle p8-1-vor-5
mut p3-vor-docker "$B_" 'if [ "$fehler" -ne 0 ]; then exit 2; fi' 'if [ "$fehler" -ne 0 ]; then :; fi' "p1e-binaries p2-binaries p3-dist p3-zweiter" p1e-binaries p2-binaries p3-dist-nicht-leer p3-zweiter-lauf

# --- Punkt 2: VERSION ---
mut p2-leer "$L_" '[[ "${VERSION-}" =~' '[[ -z "${VERSION-}" || "${VERSION-}" =~' "p2-leer p2-binaries" p2-leer p2-binaries
mut p2-leerzeichen "$L_" 'rest="0-9A-Za-z.+_-"' 'rest="0-9A-Za-z.+_ -"' p2-leerzeichen p2-leerzeichen
mut p2-dollar "$L_" 'rest="0-9A-Za-z.+_-"' 'rest="0-9A-Za-z.+_\$-"' p2-dollar p2-dollar
mut p2-einfach "$L_" 'rest="0-9A-Za-z.+_-"' "rest=\"0-9A-Za-z.+_'-\"" p2-einfach p2-einfaches-anfuehrungszeichen
mut p2-doppelt "$L_" 'rest="0-9A-Za-z.+_-"' 'rest="0-9A-Za-z.+_\"-"' p2-doppelt p2-doppeltes-anfuehrungszeichen
mut p2-minus "$L_" 'local erst="0-9A-Za-z"' 'local erst="-0-9A-Za-z"' p2-fuehrendes p2-fuehrendes-minus
mut p2-gueltig-punkt "$L_" 'rest="0-9A-Za-z.+_-"' 'rest="0-9A-Za-z+_-"' p2-gueltig p2-gueltig
mut p2-gueltig-plus "$L_" 'rest="0-9A-Za-z.+_-"' 'rest="0-9A-Za-z._-"' p2-gueltig p2-gueltig
mut p2-gueltig-unterstrich "$L_" 'rest="0-9A-Za-z.+_-"' 'rest="0-9A-Za-z.+-"' p2-gueltig p2-gueltig
mut p2-gueltig-minus "$L_" 'rest="0-9A-Za-z.+_-"' 'rest="0-9A-Za-z.+_"' p2-gueltig p2-gueltig
mut p2-dev "$G_" "'pgwire-recorder %s\\n' \"\$VERSION\"" "'pgwire-recorder %s\\n' dev" "p2-gueltig" p2-gueltig




# --- Punkt 3: das Werkzeug ---
mut p3-dist-leer "$B_" %P3distleer "p3-dist" p3-dist-leerer-wert
mut p3-dist "$B_" %P3dist "p3-dist p3-zweiter" p3-dist-nicht-leer p3-dist-datei p3-zweiter-lauf
mut p3-exe "$B_" %P3exe p3 p3-namen p3-exe p3-sums-je-binary
mut p3-name "$B_" 'datei="pgwire-recorder_${VERSION}_${os}_${arch}"' 'datei="pgwire-recorder_${os}_${arch}"' p3 p3-exe p3-namen p3-sums-format p3-sums-je-binary
mut p3-build-arg "$B_" ' --build-arg "VERSION=$VERSION"' '' p3 p3-version-im-build
mut p3-sort "$B_" %P3sort p3 p3-sums-sortiert
mut p3-sums-binaer "$B_" 'sha256sum "${namen[@]}"' 'sha256sum --binary "${namen[@]}"' p3 p3-sums-format p3-sums-je-binary
mut p3-sums-summe "$B_" 'sha256sum "${namen[@]}"' 'sha1sum "${namen[@]}"' p3 p3-sums-format p3-sums-pruefbar
mut p3-sums-name "$B_" 'sha256sum "${namen[@]}"' 'sha256sum "${namen[@]/#/.\/}"' p3 p3-sums-format p3-sums-je-binary
mut p3-sums-lf "$B_" '> SHA256SUMS)' "| sed 's/\$/\\r/' > SHA256SUMS)" p3 p3-sums-lf p3-sums-format p3-sums-je-binary
for p in $ERF; do
  mut "p3-sums-zeile-$p" "$B_" 'sha256sum "${namen[@]}"' "sha256sum \"\${namen[@]}\" | grep -v \"_${p%/*}_${p#*/}\"" p3 p3-sums-format p3-sums-je-binary
done
mut p3-dist-fehlt "$B_" 'mkdir -p "$dist"' ':' p3 p3-gruen p3-namen p3-exe p3-sums-format p3-sums-je-binary p3-sums-sortiert p3-sums-lf p3-sums-pruefbar p3-dist-fehlt p3-zweiter-lauf
mut p3-rot "$B_" 'exit 1 # P3rot' 'break # P3rot' p3-rot p3-rot-nennt-windows-arm64 p3-rot-ohne-sums p3-rot-keine-teilmenge
mut p3-teilmenge "$B_" 'mv "$stage/$os-$arch/pgwire-recorder" "$stage/$datei"' 'mv "$stage/$os-$arch/pgwire-recorder" "$dist/$datei"; cp "$dist/$datei" "$stage/$datei"' p3-rot p3-rot-keine-teilmenge
mut p3-rot-sums "$B_" 'exit 1 # P3rot' '(cd "$stage" && sha256sum * > "$OLDPWD/$dist/SHA256SUMS" 2> /dev/null); exit 1 # P3rot' p3-rot p3-rot-ohne-sums
mut p3-nur-binaries "$B_" 'if [ "${OHNE_CACHE-}" = 1 ]; then cache=(--no-cache-filter build --no-cache-filter dist); fi # P4cache' 'cache=(--no-cache-filter build --no-cache-filter dist)' p3 p3-nur-binaries

# --- Punkt 4: zwei Builds, dieselbe Summe ---
for p in $ERF; do
  x=; if [ "$p" = linux/amd64 ]; then x=p8-4-vor-6; fi
  mut "p4-abweichung-$p" "$G_" 'if ! cmp -s "$werk/a/SHA256SUMS" "$werk/b/SHA256SUMS"; then # P4cmp' "if ! cmp -s <(grep -v \"_${p%/*}_${p#*/}\" \"\$werk/a/SHA256SUMS\") <(grep -v \"_${p%/*}_${p#*/}\" \"\$werk/b/SHA256SUMS\"); then # P4cmp" "p4-abweichung-$(dt "$p") p8-4" "p4-abweichung-$(dt "$p")" $x
done
mut p4-namen "$G_" "diff <(cat \"\$werk/a/SHA256SUMS\") <(cat \"\$werk/b/SHA256SUMS\") | sed -n 's/^< [0-9a-f]*  //p'" "cat \"\$werk/a/SHA256SUMS\" | sed -n 's/^[0-9a-f]*  //p'" p4 p4-namen p4-abweichung-darwin-amd64 p4-abweichung-darwin-arm64 p4-abweichung-linux-amd64 p4-abweichung-linux-arm64 p4-abweichung-windows-amd64 p4-abweichung-windows-arm64
mut p4-zweiter "$G_" %P4zweiter "p4-cache p4-gruen" p4-cache-umgehen p4-gruen
mut p4-filter-build "$B_" 'cache=(--no-cache-filter build --no-cache-filter dist)' 'cache=(--no-cache-filter dist)' "p4-cache p4-gruen" p4-cache-umgehen p4-gruen
mut p4-filter-dist "$B_" 'cache=(--no-cache-filter build --no-cache-filter dist)' 'cache=(--no-cache-filter build)' "p4-cache p4-gruen" p4-cache-umgehen p4-gruen
mut p4-ausserhalb "$G_" 'DIST="$werk/$n"' 'DIST="dist-$n"' p10 p10-arbeitsbaum
mut p10-schreibt "$G_" 'punkt4() {' 'punkt4() { : > geschrieben' p10 p10-arbeitsbaum


# --- Punkt 5: Kopf der Binaries ---
# Der Lauf endet beim ersten roten Kopf: Eine Mutation, die eine Plattform rot färbt, verdeckt die späteren; der
# Umfang `nur` nennt die Fälle, die die Mutation allein bricht.
mut p5-elf-magic "$G_" %P5elf-magic p5-format-linux p5-format-linux-amd64 p5-format-linux-arm64
mut p5-elf-maschine "$G_" %P5elf-maschine p5-arch-linux p5-arch-linux-amd64 p5-arch-linux-arm64
mut p5-macho-magic "$G_" %P5macho-magic p5-format-darwin p5-format-darwin-amd64 p5-format-darwin-arm64
mut p5-macho-cpu "$G_" %P5macho-cpu p5-arch-darwin p5-arch-darwin-amd64 p5-arch-darwin-arm64
mut p5-pe-magic "$G_" %P5pe-magic p5-format-windows p5-format-windows-amd64 p5-format-windows-arm64
mut p5-pe-maschine "$G_" %P5pe-maschine p5-arch-windows p5-arch-windows-amd64 p5-arch-windows-arm64
mut p5-linux-amd64 "$G_" 'amd64) erw=3e00' 'amd64) erw=b700' "p5-gruen p5-arch-linux-amd64" p5-gruen p5-arch-linux-amd64
mut p5-linux-arm64 "$G_" 'arm64) erw=b700' 'arm64) erw=3e00' "p5-gruen p5-arch-linux-arm64" p5-gruen p5-arch-linux-arm64
mut p5-darwin-amd64 "$G_" 'amd64) erw=07000001' 'amd64) erw=0c000001' "p5-gruen p5-arch-darwin-amd64" p5-gruen p5-arch-darwin-amd64
mut p5-darwin-arm64 "$G_" 'arm64) erw=0c000001' 'arm64) erw=07000001' "p5-gruen p5-arch-darwin-arm64" p5-gruen p5-arch-darwin-arm64
mut p5-windows-amd64 "$G_" 'amd64) erw=6486' 'amd64) erw=64aa' "p5-gruen p5-arch-windows-amd64" p5-gruen p5-arch-windows-amd64
mut p5-windows-arm64 "$G_" 'arm64) erw=64aa' 'arm64) erw=6486' "p5-gruen p5-arch-windows-arm64" p5-gruen p5-arch-windows-arm64
mut p5-pe-offset-fest "$G_" 'g="$(hex "$f" $((o + 4)) 2)"' 'g="$(hex "$f" 132 2)"' "p5-gruen p5-arch-windows" p5-gruen p5-arch-windows-amd64 p5-arch-windows-arm64
mut p5-pe-offset-3c "$G_" 'o="$(hex "$f" 60 4)"' 'o="$(hex "$f" 56 4)"' "p5-gruen p5-arch-windows" p5-gruen p5-arch-windows-amd64 p5-arch-windows-arm64

# --- Punkt 6: Image je Plattform ---
mut p6-os "$G_" %P6os p6-os p6-os-linux-amd64 p6-os-linux-arm64
mut p6-arch "$G_" %P6arch p6-arch p6-arch-linux-amd64 p6-arch-linux-arm64
mut p6-summe "$G_" %P6summe p6-summe p6-summe-linux-amd64 p6-summe-linux-arm64
for p in $IMG; do
  mut "p6-os-$p" "$G_" '[ "$g" = "$os" ] ||' "[ \"\$g\" = \"\$os\" ] || [ \"\$p\" = $p ] ||" "p6-os-$(dt "$p")" "p6-os-$(dt "$p")"
  mut "p6-arch-$p" "$G_" '[ "$g" = "$arch" ] ||' "[ \"\$g\" = \"\$arch\" ] || [ \"\$p\" = $p ] ||" "p6-arch-$(dt "$p")" "p6-arch-$(dt "$p")"
  mut "p6-summe-$p" "$G_" '[ "$g" = "$soll" ] ||' "[ \"\$g\" = \"\$soll\" ] || [ \"\$p\" = $p ] ||" "p6-summe-$(dt "$p")" "p6-summe-$(dt "$p")"
done
mut p6-bau "$G_" 'rot "Image für $p nicht zu bauen"; } # P6bau' 'true; } # P6bau' p6-bau p6-bau-rot
mut p6-runtime-run Dockerfile 'USER nonroot' $'RUN true\nUSER nonroot' p6-runtime p6-runtime-ohne-run

# --- Punkt 7: Smoke ---
mut p7-server "$G_" %P7server p7-server-fremd p7-server-fremd
mut p7-server-amd64 "$G_" 'amd64 | arm64) ;;' 'amd64) ;;' p7-server-arm64 p7-server-arm64
mut p7-server-arm64 "$G_" 'amd64 | arm64) ;;' 'arm64) ;;' "p7-gruen p7-server-fremd" p7-gruen
mut p7-tag "$G_" 'tag="$lauf:linux-$arch"' 'tag="$lauf:linux-amd64"' p7-server-arm64 p7-server-arm64
mut p7-intern "$G_" 'network create --internal' 'network create' p7-intern p7-intern
mut p7a-exit "$G_" '|| rot "version endet mit einem Exit-Code ungleich 0" # P7a-exit' '|| true # P7a-exit' p7a-exit p7a-exit
mut p7a-ausgabe "$G_" %P7a-ausgabe p7a p7a-ausgabe-abweichend p7a-ohne-lf p7a-zusatzzeile
mut p7a-lf "$G_" 'cmp -s "$werk/version" "$werk/version-soll"' '[ "$(cat "$werk/version")" = "$(cat "$werk/version-soll")" ]' p7a p7a-ohne-lf
mut p7a-zeile "$G_" 'cmp -s "$werk/version" "$werk/version-soll"' 'cmp -s <(head -n1 "$werk/version") "$werk/version-soll"' p7a p7a-zusatzzeile
mut p7b-leer "$G_" %P7a-leer "p7b-sicht-a-leer p7b-beide-leer" p7b-sicht-a-leer p7b-beide-leer
mut p7b-eins "$G_" %P7a-eins p7b-ohne-eins p7b-ohne-eins
mut p7b-pg-bereit "$G_" '|| rot "PostgreSQL nicht bereit"' '|| true' p7b-pg-nicht-bereit p7b-pg-nicht-bereit
mut p7b-cp "$G_" '|| rot "die Aufzeichnung ist dem Container nicht zu entnehmen" # P7cp' '|| true # P7cp' p7b-cp-rot p7b-cp-rot
mut p7b-aufz "$G_" %P7aufz p7b-aufzeichnung-leer p7b-aufzeichnung-leer
mut p7-exit-rec "$G_" '[ "$c" = 0 ] ||' '[ "$c" = 0 ] || [ "$2" = record ] ||' "p7b-exit-rec p7c-exit-rep" p7b-exit-rec
mut p7-exit-rep "$G_" '[ "$c" = 0 ] ||' '[ "$c" = 0 ] || [ "$2" = replay ] ||' "p7b-exit-rec p7c-exit-rep" p7c-exit-rep
mut p7c-pg "$G_" %P7pg p7c-pg-lebt p7c-pg-lebt
mut p7c-gleich "$G_" %P7gleich p7c-sicht-ungleich p7c-sicht-ungleich
mut p7-terminal "$G_" 'run -d --name "$rec"' 'run -d -t --name "$rec"' p7-kein-terminal p7-kein-terminal
mut p7-user "$G_" 'run -d --name "$rec"' 'run -d --user 0 --name "$rec"' p7-kein-user p7-kein-user
mut p7-volume "$G_" 'run -d --name "$rec"' 'run -d -v /tmp:/tmp --name "$rec"' p7-kein-volume p7-kein-volume
mut p7-mount "$G_" 'run -d --name "$rec"' 'run -d --mount type=tmpfs,dst=/x --name "$rec"' p7-kein-mount p7-kein-mount
mut p7-cp-hinein "$G_" 'l "$docker" cp "$werk/aufzeichnung.yaml" "$rep:/tmp/aufzeichnung.yaml"' ':' p7-cp-hinein p7-cp-hinein
mut p7-cp-entnehmen "$G_" 'l "$docker" cp "$rec:/tmp/aufzeichnung.yaml"' 'l "$docker" cp "$rec:/tmp/anders.yaml"' p7-cp-entnehmen p7-cp-entnehmen
mut p7-rm-pg "$G_" 'l "$docker" rm -f "$pg" > /dev/null' ':' p7-pg-vor-replay p7-pg-vor-replay


# --- Punkt 8: Reihenfolge der Gründe ---
mut p8-4-weiter "$G_" 'rot "die Summen zweier Builds weichen ab: ${abw% }"' 'echo "test-container: ROT — die Summen zweier Builds weichen ab: ${abw% }" >&2' "p4-abweichung p4-namen p4-cache p8-4" p4-abweichung-darwin-amd64 p4-abweichung-darwin-arm64 p4-abweichung-linux-amd64 p4-abweichung-linux-arm64 p4-abweichung-windows-amd64 p4-abweichung-windows-arm64 p4-namen p4-cache-umgehen p8-4-vor-6

# --- Punkt 9: Aufräumen ---
mut p9-container "$G_" %P9container p9 p9-gruen-container p9-rot7-container p9-sigINT-psql-container p9-sigTERM-psql-container
mut p9-netz "$G_" %P9netz p9 p9-gruen-netz p9-rot7-netz p9-sigINT-psql-netz p9-sigTERM-psql-netz
mut p9-tags "$G_" %P9tags p9 p9-gruen-tags p9-rot6-tags p9-rot7-tags p9-sigINT-psql-tags p9-sigTERM-psql-tags
mut p9-werk "$G_" %P9werk p9 p9-gruen-werk p9-rot4-werk p9-rot5-werk p9-rot6-werk p9-rot7-werk p9-sigINT-build-werk p9-sigINTg-build-werk p9-sigINT-psql-werk p9-sigTERM-build-werk p9-sigTERM-psql-werk
mut p9-exit "$G_" %P9exit p9 p9-gruen-container p9-gruen-netz p9-gruen-tags p9-gruen-werk p9-rot4-werk p9-rot5-werk p9-rot6-tags p9-rot6-werk p9-rot7-container p9-rot7-netz p9-rot7-tags p9-rot7-werk p9-sigINT-build-werk p9-sigINTg-build-werk p9-sigINT-psql-container p9-sigINT-psql-netz p9-sigINT-psql-tags p9-sigINT-psql-werk p9-sigTERM-build-werk p9-sigTERM-psql-container p9-sigTERM-psql-netz p9-sigTERM-psql-tags p9-sigTERM-psql-werk
mut p9-int "$G_" "trap 'beenden 130' INT # P9int" "trap 'kill -KILL \$\$' INT # P9int" p9-sigINT p9-sigINT-build-signal p9-sigINT-build-stage p9-sigINT-build-werk p9-sigINTg-build-signal p9-sigINTg-build-stage p9-sigINTg-build-werk p9-sigINT-psql-container p9-sigINT-psql-netz p9-sigINT-psql-signal p9-sigINT-psql-tags p9-sigINT-psql-werk
mut p9-term "$G_" "trap 'beenden 143' TERM # P9term" "trap 'kill -KILL \$\$' TERM # P9term" p9-sigTERM p9-sigTERM-build-signal p9-sigTERM-build-stage p9-sigTERM-build-werk p9-sigTERM-psql-container p9-sigTERM-psql-netz p9-sigTERM-psql-signal p9-sigTERM-psql-tags p9-sigTERM-psql-werk
mut p9-weiter "$G_" 'kill -TERM "$laeuft" 2> /dev/null || true' ':' p9 p9-sigINT-build-signal p9-sigINTg-build-signal p9-sigINT-psql-signal p9-sigTERM-build-signal p9-sigTERM-psql-signal
mut p9-stage-exit "$B_" "trap 'rm -rf \"\$stage\"' EXIT" ':' p9 p9-binINT-stage p9-binTERM-stage p9-gruen-stage p9-rot4-stage p9-rot5-stage p9-rot6-stage p9-rot7-stage p9-sigINT-build-stage p9-sigINTg-build-stage p9-sigINT-psql-stage p9-sigTERM-build-stage p9-sigTERM-psql-stage
mut p9-stage-int "$B_" "trap 'beenden 130' INT" "trap 'kill -KILL \$\$' INT" p9 p9-binINT-signal p9-binINT-stage
mut p9-stage-term "$B_" "trap 'beenden 143' TERM" "trap 'kill -KILL \$\$' TERM" p9 p9-binTERM-signal p9-binTERM-stage p9-sigINT-build-signal p9-sigINT-build-stage p9-sigINTg-build-signal p9-sigINTg-build-stage p9-sigTERM-build-signal p9-sigTERM-build-stage
mut p9-stage-weiter "$B_" 'kill -TERM "$laeuft" 2> /dev/null || true' ':' p9 p9-binINT-signal p9-binTERM-signal p9-sigINT-build-signal p9-sigINTg-build-signal p9-sigTERM-build-signal
mut p9-kennung-ziffern "$G_" 'lauf="pgr-ct-$$"' 'lauf="pgr-ct-x"' p9-kennung p9-kennung-ziffern p9-kennung-eindeutig
mut p9-kennung-eindeutig "$G_" 'lauf="pgr-ct-$$"' 'lauf="pgr-ct-1"' p9-kennung p9-kennung-eindeutig

# --- Punkt 11: Aufrufarten des Treibers ---
mut p11-make-binaries "$F_" $'\t@DOCKER=\'$(DOCKER)\' bash tools/harness/binaries.sh' $'\t@$(MAKE) --version >/dev/null; DOCKER=\'$(DOCKER)\' bash tools/harness/binaries.sh' p11-werkz p11-werkz-makeflags-n p11-werkz-makeflags-q p11-werkz-makeflags-t p11-werkz-n p11-werkz-ni p11-werkz-q p11-werkz-qi p11-werkz-t p11-werkz-ti
mut p11-plus-binaries "$F_" $'\t@DOCKER=\'$(DOCKER)\' bash tools/harness/binaries.sh' $'\t+@DOCKER=\'$(DOCKER)\' bash tools/harness/binaries.sh' p11-werkz p11-werkz-makeflags-n p11-werkz-makeflags-q p11-werkz-makeflags-t p11-werkz-n p11-werkz-ni p11-werkz-q p11-werkz-qi p11-werkz-t p11-werkz-ti
mut p11-make-gate "$F_" $'\t@POSTGRES_IMAGE=' $'\t@$(MAKE) --version >/dev/null; POSTGRES_IMAGE=' p11-gate p11-gate-makeflags-n p11-gate-makeflags-q p11-gate-makeflags-t p11-gate-n p11-gate-ni p11-gate-q p11-gate-qi p11-gate-t p11-gate-ti
mut p11-plus-gate "$F_" $'\t@POSTGRES_IMAGE=' $'\t+@POSTGRES_IMAGE=' p11-gate p11-gate-makeflags-n p11-gate-makeflags-q p11-gate-makeflags-t p11-gate-n p11-gate-ni p11-gate-q p11-gate-qi p11-gate-t p11-gate-ti
mut p11-i-binaries "$F_" $'\t@DOCKER=\'$(DOCKER)\' bash tools/harness/binaries.sh' $'\t@case "$(MAKEFLAGS)" in *i*) ;; *) DOCKER=\'$(DOCKER)\' bash tools/harness/binaries.sh ;; esac' "p11-werkz-i p11-werkz-makeflags-i" p11-werkz-i p11-werkz-makeflags-i
mut p11-i-gate "$F_" $'\t@POSTGRES_IMAGE=\'$(POSTGRES_IMAGE)\' DOCKER=\'$(DOCKER)\' bash tools/harness/run-container.sh' $'\t@case "$(MAKEFLAGS)" in *i*) ;; *) POSTGRES_IMAGE=\'$(POSTGRES_IMAGE)\' DOCKER=\'$(DOCKER)\' bash tools/harness/run-container.sh ;; esac' "p11-gate-i p11-gate-makeflags-i" p11-gate-i p11-gate-makeflags-i
mut p11-make-gegenprobe "$F_" $'\t@bash tools/harness/container-gegenprobe.sh' $'\t@$(MAKE) --version >/dev/null; bash tools/harness/container-gegenprobe.sh' p11-gegenprobe p11-gegenprobe-rezept
mut p11-plus-gegenprobe "$F_" $'\t@bash tools/harness/container-gegenprobe.sh' $'\t+@bash tools/harness/container-gegenprobe.sh' p11-gegenprobe p11-gegenprobe-rezept

# --- Punkt 12: Anschluss an GATE_CHECKS ---
mut p12-ziel "$F_" 'GATE_CHECKS += test-container ' 'GATE_CHECKS += ' p12 p12-gate-checks-ziel
mut p12-gegenprobe "$F_" 'GATE_CHECKS += test-container test-container-gegenprobe' 'GATE_CHECKS += test-container' p12 p12-gate-checks-gegenprobe

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "container-gegenprobe: gruen — Punkte 1 bis 9, 11 und 12 (Punkt 10: kein Schreiben in den Arbeitsbaum): unveränderte Kopie grün, je Mutation genau die benannten Fälle rot"
