#!/usr/bin/env bash
# release — Schritte der Veröffentlichung bei einem Tag (harness/vertraege/SPEC-067-release.md, Fassung A): eingaben, binaries, github-release.
# Aufruf: bash tools/harness/release.sh <Schritt> im Wurzelverzeichnis des Klons; die Pipeline ruft es nie unter make (Punkt 13).
# Umgebung: RELEASE_EVENT, RELEASE_REF, RELEASE_SHA, RELEASE_DEFAULT_BRANCH, RELEASE_REPOSITORY (gesetzt von .github/workflows/release.yml),
# RUNNER_TEMP (Ort des Zustands zwischen den Schritten: $RUNNER_TEMP/release mit den Dateien eingaben.ok, binaries.ok und dist/), GH_TOKEN (github-release).
# Ausgang: 0 grün; 1 ein Schritt scheitert (Bau, Summe, Dateiliste, Release); 2 Eingabe oder Vorbedingung (ein Schritt, dessen Vorgänger nicht grün war). Eine abgelehnte Eingabe
# steht als `release: <Schlüsselwort>: …` auf stderr, das Schlüsselwort ist der Grund (ereignis, tag-form, version-datei, tag-ungleich, notes, standard-zweig); ein Hinweis beginnt mit `release: hinweis: `.
# Das Skript gibt keinen Wert eines Geheimnisses aus (Punkt 10): kein `set -x`, kein `env`.
# Braucht bash, git, make, sha256sum, awk, cmp, grep, sort, pgrep, mktemp und gh (nur github-release); eine Version prüft das Skript nicht, zugesagt ist die Plattform von SPEC-067.
# GEPRÜFT DURCH tools/harness/release-gegenprobe.sh (Fälle p2-* bis p13-*, je Zusage eine Mutation dieser Datei).
set -uo pipefail
export LC_ALL=C

schritt="${1-}"
rc=0
melde() { echo "release: $schritt: $*" >&2; }
eingabe() { echo "release: $1: $2" >&2; exit 2; }
hinweis() { echo "release: hinweis: $*" >&2; }
rot() { melde "ROT — $*"; exit 1; }

# --- Aufräumen und Abbruch (Punkt 13) ----------------------------------------------------------------------------------------------------------------
scratch="$(mktemp -d)" || { echo "release: Temp-Verzeichnis nicht anlegbar" >&2; exit 2; }
raeume_auf() { rm -rf "${scratch:?}"; }
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
# fuehre <Befehl>… — der Befehl läuft im Hintergrund und das Skript wartet mit wait; so läuft die Falle für INT und TERM sofort. Setzt rc.
fuehre() { "$@" & wait "$!"; rc=$?; }

# --- Zustand zwischen den Schritten -----------------------------------------------------------------------------------------------------------------
zdir=""
zustand() {
  [ -n "${RUNNER_TEMP-}" ] || eingabe zustand "RUNNER_TEMP fehlt"
  zdir="$RUNNER_TEMP/release"
  mkdir -p "$zdir"
}

# --- Ereignis, Tag, Version -------------------------------------------------------------------------------------------------------------------------
trockenlauf=0
tag=""
# tag_ableiten — der Tag des Laufs aus dem Ref, der des Trockenlaufs ist v<Inhalt der ersten Zeile von VERSION>.
tag_ableiten() {
  if [ "${RELEASE_EVENT-}" = workflow_dispatch ]; then
    trockenlauf=1
    tag="v$(head -n 1 VERSION 2> /dev/null)"
  else
    tag="${RELEASE_REF-}"
    tag="${tag#refs/tags/}"
  fi
}
# aussenwirkung_erlaubt — ein Schritt mit Außenwirkung läuft nur bei einem Push von refs/tags/v… (Punkt 2).
aussenwirkung_erlaubt() {
  [ "${RELEASE_EVENT-}" = push ] || eingabe ereignis "Ereignis '${RELEASE_EVENT-}' ist kein Push eines Tags" # P2ereignis
  [[ ${RELEASE_REF-} == refs/tags/v* ]] || eingabe ereignis "Ref '${RELEASE_REF-}' ist kein Tag v…" # P2ref
}

NUM='(0|[1-9][0-9]*)'
BEZ='(0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)' # P3bez
KERN="$NUM\\.$NUM\\.$NUM" # P3kern
VOR="(-$BEZ(\\.$BEZ)*)?" # P3vor
SEMVER="$KERN$VOR"
VERSION_RE="^$SEMVER"$'\n'"\$" # P3vdatei
HANDBUCH=docs/user/benutzerhandbuch.md

# (1) Ereignis und Ref
pruefe_ereignis() {
  if [ "${RELEASE_EVENT-}" = workflow_dispatch ]; then return 0; fi
  aussenwirkung_erlaubt
}
# (2) Form des Tags
pruefe_tag_form() {
  [ "${#tag}" -le 100 ] || eingabe tag-form "der Tag ist länger als 100 Zeichen" # P3laenge
  [[ $tag =~ ^v$SEMVER$ ]] || eingabe tag-form "der Tag ${tag@Q} hat nicht die Form v<SemVer>" # P3form
}
# (3) Form der Datei VERSION
version=""
pruefe_version_datei() {
  local inhalt
  inhalt="$(cat VERSION 2> /dev/null; printf x)"
  inhalt="${inhalt%x}"
  [[ $inhalt =~ $VERSION_RE ]] || eingabe version-datei "VERSION ist nicht genau eine Zeile <SemVer> ohne v mit LF" # P3vform
  version="${inhalt%$'\n'}"
}
# (4) Tag gleich v plus Inhalt von VERSION
pruefe_tag_gleich() {
  [ "$tag" = "v$version" ] || eingabe tag-ungleich "der Tag ${tag@Q} ist nicht v$version (VERSION)" # P3gleich
}
# notiz_der_version <Version> — `Z` und die dritte Zelle der Zeile der Version in der Tabelle des Abschnitts, dessen Überschrift auf Änderungshistorie endet; `X`, wenn die Zeile nicht genau drei Zellen hat (ein `|` im Text, auch `\|`); nichts ohne Zeile.
notiz_der_version() {
  awk -F'|' -v v="$1" '
    function trim(s) { gsub(/^[ \t]+|[ \t]+$/, "", s); return s }
    /^#+[ \t].*Änderungshistorie[ \t]*$/ { ab = 1; next }
    /^#/ { ab = 0 }
    ab && /^\|/ && trim($2) == v { if (NF != 5) print "X"; else print "Z" trim($4); exit }
  ' "$HANDBUCH" 2> /dev/null
}
# (5) Zeile der Version in der Änderungshistorie mit nicht leerer dritter Zelle; im Trockenlauf ist eine fehlende Zeile ein Hinweis und eine vorhandene wird als Hinweis gezeigt
pruefe_notes() {
  local z
  z="$(notiz_der_version "$version")"
  if [ -z "$z" ]; then
    if [ "$trockenlauf" -eq 1 ]; then # P3zeiletrocken
      hinweis "die Änderungshistorie des Benutzerhandbuchs hat keine Zeile für $version"
      return 0
    fi
    eingabe notes "die Änderungshistorie des Benutzerhandbuchs hat keine Zeile für $version" # P3notes
  fi
  [[ $z != X* ]] || eingabe notes "die Zeile für $version hat nicht genau drei Zellen (ein | im Text)" # P3pipe
  [ -n "${z#Z}" ] || eingabe notes "die Zeile für $version hat eine leere dritte Zelle" # P3zelle
  if [ "$trockenlauf" -eq 1 ]; then hinweis "Änderungshistorie $version: ${z#Z}"; fi # P3zeiletext
}
# (6) Commit des Tags im Standard-Zweig; im Trockenlauf ein Hinweis
pruefe_standard_zweig() {
  local ref="origin/${RELEASE_DEFAULT_BRANCH-}" ok=1
  git merge-base --is-ancestor "${RELEASE_SHA-}" "$ref" 2> /dev/null || ok=0 # P3anc
  if [ "$trockenlauf" -eq 1 ]; then
    if [ "$ok" -eq 1 ]; then hinweis "der Commit liegt im Standard-Zweig $ref"; else hinweis "der Commit liegt nicht im Standard-Zweig $ref"; fi
  elif [ "$ok" -eq 0 ]; then
    eingabe standard-zweig "der Commit ${RELEASE_SHA-} liegt nicht im Standard-Zweig $ref"
  fi
}

schritt_eingaben() {
  zustand
  rm -f "$zdir/eingaben.ok" "$zdir/binaries.ok"
  tag_ableiten
  pruefe_ereignis
  pruefe_tag_form
  pruefe_version_datei
  pruefe_tag_gleich
  pruefe_notes
  pruefe_standard_zweig
  : > "$zdir/eingaben.ok"
  echo "release: eingaben: grün für $tag ($([ "$trockenlauf" -eq 1 ] && echo Trockenlauf || echo Lauf))"
}

# erwartete_namen <Version> — die sechs Binaries von SPEC-055 Punkt 3 und SHA256SUMS, einer je Zeile.
erwartete_namen() {
  local os a e
  for os in darwin linux windows; do
    for a in amd64 arm64; do
      e=""
      if [ "$os" = windows ]; then e=.exe; fi
      echo "pgwire-recorder_${1}_${os}_${a}$e"
    done
  done
  echo SHA256SUMS
}

schritt_binaries() {
  local dist
  zustand
  [ -f "$zdir/eingaben.ok" ] || eingabe vorbedingung "der Schritt eingaben war nicht grün" # P11eingaben
  rm -f "$zdir/binaries.ok"
  tag_ableiten
  dist="$zdir/dist"
  fuehre make binaries "VERSION=${tag#v}" "DIST=$dist" # P4aufruf
  [ "$rc" -eq 0 ] || rot "make binaries endet mit Exit $rc"
  (cd "$dist" && sha256sum -c SHA256SUMS) > "$scratch/summen" 2>&1 || { sed 's/^/release: binaries: /' "$scratch/summen" >&2; rot "die Summen von SHA256SUMS stimmen nicht"; } # P4summe
  [ "$(ls -A "$dist" | sort)" = "$(erwartete_namen "${tag#v}" | sort)" ] || rot "die Dateien in $dist sind nicht genau die sechs Binaries und SHA256SUMS" # P4liste
  : > "$zdir/binaries.ok"
  echo "release: binaries: $(erwartete_namen "${tag#v}" | wc -l) Dateien in $dist"
}

schritt_github_release() {
  local dist v n a vorhanden=0 prerelease=() anhaenge=() fehlend=()
  aussenwirkung_erlaubt
  zustand
  [ -f "$zdir/binaries.ok" ] || eingabe vorbedingung "der Schritt binaries war nicht grün" # P11binaries
  [ -n "${GH_TOKEN-}" ] || eingabe token "GH_TOKEN fehlt" # P9token
  tag_ableiten
  v="${tag#v}"
  dist="$zdir/dist"
  n="$(notiz_der_version "$v")"
  if [[ $n == Z* ]]; then n="${n#Z}"; else n=""; fi
  [ -n "$n" ] || eingabe notes "die Änderungshistorie des Benutzerhandbuchs hat keine Zeile mit Text für $v"
  printf '%s\n' "$n" > "$scratch/notes.md"
  [[ $v != *-* ]] || prerelease=(--prerelease) # P9pre
  while IFS= read -r a; do anhaenge+=("$dist/$a"); done < <(erwartete_namen "$v") # P9anhaenge
  fuehre gh release view "$tag" --repo "$RELEASE_REPOSITORY" --json assets --jq '.assets[].name' > "$scratch/assets" 2> "$scratch/gh.err"
  if [ "$rc" -eq 0 ]; then
    vorhanden=1
  elif ! grep -q 'release not found' "$scratch/gh.err"; then
    sed 's/^/release: github-release: /' "$scratch/gh.err" >&2
    rot "gh release view endet mit Exit $rc"
  fi
  if [ "$vorhanden" -eq 0 ]; then
    fuehre gh release create "$tag" --repo "$RELEASE_REPOSITORY" --verify-tag --title "$tag" --notes-file "$scratch/notes.md" ${prerelease[@]+"${prerelease[@]}"} "${anhaenge[@]}" # P9aufruf
    [ "$rc" -eq 0 ] || rot "gh release create endet mit Exit $rc"
    echo "release: github-release: $tag angelegt"
    return 0
  fi
  for a in "${anhaenge[@]}"; do grep -qxF "${a##*/}" "$scratch/assets" || fehlend+=("$a"); done
  if grep -qxF SHA256SUMS "$scratch/assets"; then
    mkdir "$scratch/remote"
    fuehre gh release download "$tag" --repo "$RELEASE_REPOSITORY" --pattern SHA256SUMS --dir "$scratch/remote"
    [ "$rc" -eq 0 ] || rot "gh release download endet mit Exit $rc"
    cmp -s "$scratch/remote/SHA256SUMS" "$dist/SHA256SUMS" || rot "SHA256SUMS am Release weicht von der lokalen Datei ab; der Weg ist ein höherer Tag" # P9vergleich
  fi
  if [ "${#fehlend[@]}" -gt 0 ]; then
    fuehre gh release upload "$tag" --repo "$RELEASE_REPOSITORY" "${fehlend[@]}" # P9upload
    [ "$rc" -eq 0 ] || rot "gh release upload endet mit Exit $rc"
  fi
  echo "release: github-release: $tag vorhanden, ${#fehlend[@]} Anhänge ergänzt"
}

case "$schritt" in
  eingaben) schritt_eingaben ;;
  binaries) schritt_binaries ;;
  github-release) schritt_github_release ;;
  *) echo "release: unbekannter Schritt '$schritt'" >&2; exit 2 ;;
esac
