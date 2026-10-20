#!/usr/bin/env bash
# release — Schritte der Veröffentlichung bei einem Tag (harness/vertraege/SPEC-067-release.md): eingaben, binaries, image, scan, veroeffentlichen, vergleichen, smoke, latest, github-release, homebrew.
# Aufruf: bash tools/harness/release.sh <Schritt> im Wurzelverzeichnis des Klons; die Pipeline ruft es nie unter make (Punkt 13).
# Umgebung: RELEASE_EVENT, RELEASE_REF, RELEASE_SHA, RELEASE_DEFAULT_BRANCH, RELEASE_REPOSITORY (gesetzt von .github/workflows/release.yml),
# RUNNER_TEMP (Ort des Zustands zwischen den Schritten: $RUNNER_TEMP/release mit den Dateien eingaben.ok, binaries.ok, scan.ok, vergleichen.ok, smoke.ok und den Verzeichnissen dist/ und image/), GH_TOKEN (github-release),
# RELEASE_DOCKERHUB_USER und RELEASE_DOCKERHUB_TOKEN (eingaben, veroeffentlichen, latest), RELEASE_GHCR_USER und RELEASE_GHCR_TOKEN (veroeffentlichen, latest), RELEASE_TAP_TOKEN und GH_TOKEN (homebrew), RELEASE_TAP_URL (homebrew, Vorgabe der Tap des Projekts; nur für Gegenprobe und Nachtrag).
# Ausgang: 0 grün; 1 ein Schritt scheitert (Bau, Summe, Dateiliste, Release, ein Befund des Scans); 2 Eingabe oder Vorbedingung (ein Schritt, dessen Vorgänger nicht grün war) oder ein Scan, der nicht lief. Eine abgelehnte Eingabe
# steht als `release: <Schlüsselwort>: …` auf stderr, das Schlüsselwort ist der Grund (ereignis, tag-form, version-datei, tag-ungleich, notes, standard-zweig, geheimnis); der Scan meldet `release: scan: <Stichwort>: …`
# (schwachstelle, scan-gescheitert, pruefmenge, datenbank, plattform); der Schritt homebrew `release: homebrew: <Stichwort>: …` bei Exit 1 (summen-laden, summen, vorlage, tap-klon, tap-zweig, tap-pfad, tap-formel, tap-push); ein Hinweis beginnt mit `release: hinweis: `.
# Das Skript gibt keinen Wert eines Geheimnisses aus (Punkt 10): kein `set -x`, kein `env`.
# Braucht bash, git, make, sha256sum, awk, cmp, grep, sed, tar, sort (GNU, mit -V), pgrep, mktemp, docker mit buildx (image) und Zugriff auf das Image des Scanners (scan), jq und date (GNU, scan), skopeo (veroeffentlichen, vergleichen, latest), gh (github-release, homebrew) und GNU sort mit -V (latest, homebrew); eine Version prüft das Skript nicht, zugesagt ist die Plattform von SPEC-067.
# GEPRÜFT DURCH tools/harness/release-gegenprobe.sh (Fälle p2-* bis p19-*, je Zusage eine Mutation dieser Datei).
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

# (7) Geheimnisse der Registries gesetzt und nicht leer; die Meldung nennt die Namen, nie einen Wert
pruefe_geheimnisse() {
  local fehlt=""
  [ -n "${RELEASE_DOCKERHUB_USER-}" ] || fehlt+=" DOCKERHUB_USERNAME" # P3secret-user
  [ -n "${RELEASE_DOCKERHUB_TOKEN-}" ] || fehlt+=" DOCKERHUB_TOKEN" # P3secret-token
  [ -z "$fehlt" ] || eingabe geheimnis "gesetzt und nicht leer verlangt:$fehlt"
}

schritt_eingaben() {
  zustand
  rm -f "$zdir/eingaben.ok" "$zdir/binaries.ok"
  rm -rf "$zdir/vergleichen.ok" "$zdir/image" # P11loescht-b
  rm -f "$zdir/smoke.ok" # P11loescht-smoke
  rm -f "$zdir/scan.ok" # P16loescht-eingaben
  tag_ableiten
  pruefe_ereignis
  pruefe_tag_form
  pruefe_version_datei
  pruefe_tag_gleich
  pruefe_notes
  pruefe_standard_zweig
  pruefe_geheimnisse
  : > "$zdir/eingaben.ok"
  echo "release: eingaben: grün für $tag ($([ "$trockenlauf" -eq 1 ] && echo Trockenlauf || echo Lauf))"
}

# erwartete_namen — die sechs Binaries von SPEC-055 Punkt 3 und SHA256SUMS, einer je Zeile.
erwartete_namen() {
  local os a e
  for os in darwin linux windows; do
    for a in amd64 arm64; do
      e=""
      if [ "$os" = windows ]; then e=.exe; fi
      echo "pgwire-recorder-${os}-${a}$e"
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
  [ "$(ls -A "$dist" | sort)" = "$(erwartete_namen | sort)" ] || rot "die Dateien in $dist sind nicht genau die sechs Binaries und SHA256SUMS" # P4liste
  : > "$zdir/binaries.ok"
  echo "release: binaries: $(erwartete_namen | wc -l) Dateien in $dist"
}

schritt_github_release() {
  local dist v n a vorhanden=0 prerelease=() anhaenge=() fehlend=()
  aussenwirkung_erlaubt
  zustand
  [ -f "$zdir/binaries.ok" ] || eingabe vorbedingung "der Schritt binaries war nicht grün" # P11binaries
  [ -f "$zdir/smoke.ok" ] || eingabe vorbedingung "der Schritt smoke war nicht grün" # P9registries
  [ -n "${GH_TOKEN-}" ] || eingabe token "GH_TOKEN fehlt" # P9token
  tag_ableiten
  v="${tag#v}"
  dist="$zdir/dist"
  n="$(notiz_der_version "$v")"
  if [[ $n == Z* ]]; then n="${n#Z}"; else n=""; fi
  [ -n "$n" ] || eingabe notes "die Änderungshistorie des Benutzerhandbuchs hat keine Zeile mit Text für $v"
  printf '%s\n' "$n" > "$scratch/notes.md"
  [[ $v != *-* ]] || prerelease=(--prerelease) # P9pre
  while IFS= read -r a; do anhaenge+=("$dist/$a"); done < <(erwartete_namen) # P9anhaenge
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

# --- Image, Registries (Punkte 5 bis 8) ---------------------------------------------------------------------------------------------------------------
IMAGE_NAME=pt9912/pgwire-recorder
BUILDKIT_IMAGE='moby/buildkit@sha256:cec9f139f45e93c5c69c60f8b07cfad9f43f4ef6b6a6cd917527fea5ff2e3dea' # v0.33.1
FEHLT_RE='manifest unknown|name unknown|not found' # Meldungen von skopeo inspect, die „der Tag fehlt“ heißen (ohne Beachtung der Schreibweise)
REGISTRIES=(ghcr.io docker.io)
# verlange <Datei unter dem Zustand> <Schritt> — Exit 2, wenn der Vorgänger nichts hinterlassen hat.
verlange() { [ -e "$zdir/$1" ] || eingabe vorbedingung "der Schritt $2 war nicht grün"; }
# manifeste <Index> — Zahl der Manifeste; plattformen <Index> — ihre Plattformen `os/arch`, sortiert, je mit Leerzeichen dahinter.
manifeste() { tr -d ' \t\n' < "$1" | grep -o '"digest":"sha256:[0-9a-f]*"' | wc -l; }
plattformen() {
  local p o a
  tr -d ' \t\n' < "$1" | grep -o '"platform":{[^}]*}' | while IFS= read -r p; do
    o="${p#*\"os\":\"}"; a="${p#*\"architecture\":\"}"
    echo "${o%%\"*}/${a%%\"*}"
  done | sort | tr '\n' ' '
}
# index_pruefen <Index> <Ort> — genau zwei Manifeste, für linux/amd64 und linux/arm64.
index_pruefen() {
  [ "$(manifeste "$1")" -eq 2 ] || rot "$2 hat nicht genau zwei Manifeste" # P5zahl
  [ "$(plattformen "$1")" = "linux/amd64 linux/arm64 " ] || rot "$2 nennt die Plattformen '$(plattformen "$1")' statt linux/amd64 linux/arm64" # P5plattformen
}
# anmelden <Registry> — Anmeldung per --password-stdin, der Token steht nie in einem Argument.
anmelden() {
  local u t
  if [ "$1" = ghcr.io ]; then u="${RELEASE_GHCR_USER-}"; t="${RELEASE_GHCR_TOKEN-}"; else u="${RELEASE_DOCKERHUB_USER-}"; t="${RELEASE_DOCKERHUB_TOKEN-}"; fi
  printf '%s' "$t" | skopeo login --authfile "$scratch/auth.json" --username "$u" --password-stdin "$1" > /dev/null 2>&1 || rot "Anmeldung bei $1 scheitert" # P10stdin
}
# zustand_tag <Tag> — je Registry der Digest der rohen Bytes des Tags oder -.
zustand_tag() {
  local reg out=""
  for reg in "${REGISTRIES[@]}"; do
    if skopeo inspect --raw --authfile "$scratch/auth.json" "docker://$reg/$IMAGE_NAME:$1" > "$scratch/zroh" 2> /dev/null; then
      out+=" $reg:$1=sha256:$(sha256sum "$scratch/zroh" | cut -d' ' -f1)"
    else
      out+=" $reg:$1=-"
    fi
  done
  echo "Zustand:$out"
}

schritt_image() {
  local dist v zeit arch soll ist idx hex
  local -a bargs
  zustand
  verlange binaries.ok binaries # P13image
  tag_ableiten
  v="${tag#v}"
  dist="$zdir/dist"
  rm -rf "$zdir/image"; mkdir -p "$zdir/image"
  rm -f "$zdir/scan.ok" # P16loescht-image
  zeit="$(git log -1 --format=%ct "${RELEASE_SHA-}")" || rot "die Commit-Zeit von ${RELEASE_SHA-} ist nicht zu lesen"
  fuehre docker buildx create --name pgr-release --driver docker-container --driver-opt "image=$BUILDKIT_IMAGE" --bootstrap
  [ "$rc" -eq 0 ] || rot "docker buildx create endet mit Exit $rc"
  bargs=(buildx build --builder pgr-release --target runtime --build-arg "VERSION=$v" --build-arg "SOURCE_DATE_EPOCH=$zeit" --provenance=false --sbom=false)
  fuehre docker "${bargs[@]}" --platform linux/amd64,linux/arm64 --output "type=oci,dest=$zdir/image/layout.tar,rewrite-timestamp=true" . # P5aufruf
  [ "$rc" -eq 0 ] || rot "docker buildx build endet mit Exit $rc"
  tar -xOf "$zdir/image/layout.tar" index.json > "$scratch/index.json" 2> /dev/null || rot "das Layout enthält keine index.json"
  [ "$(manifeste "$scratch/index.json")" -eq 1 ] || rot "das Layout enthält nicht genau einen Index" # P5ein-index
  grep -q '"manifests":\[{"mediaType":"application/vnd.oci.image.index.v1+json"' <(tr -d ' \t\n' < "$scratch/index.json") || rot "der Eintrag des Layouts ist kein Index" # P5index-typ
  hex="$(grep -o '"digest": *"sha256:[0-9a-f]*"' "$scratch/index.json" | grep -o '[0-9a-f]\{64\}')"
  tar -xOf "$zdir/image/layout.tar" "blobs/sha256/$hex" > "$scratch/index" 2> /dev/null || rot "das Layout enthält den Index sha256:$hex nicht"
  index_pruefen "$scratch/index" "der Index des Layouts" # P5manifeste
  for arch in amd64 arm64; do
    fuehre docker "${bargs[@]}" --platform "linux/$arch" --output "type=local,dest=$scratch/fs-$arch" .
    [ "$rc" -eq 0 ] || rot "docker buildx build für linux/$arch endet mit Exit $rc"
    soll="$(awk -v n="pgwire-recorder-linux-$arch" '$2 == n { print $1 }' "$dist/SHA256SUMS")"
    ist="$(sha256sum "$scratch/fs-$arch/pgwire-recorder" 2> /dev/null | cut -d' ' -f1)"
    [ -n "$soll" ] && [ "$ist" = "$soll" ] || rot "linux/$arch: das Binary im Image hat die Summe '$ist' statt '$soll'" # P5summe
  done
  idx="sha256:$(sha256sum "$scratch/index" | cut -d' ' -f1)" # P5digest
  printf '%s\n' "$idx" > "$zdir/image/digest"
  echo "release: image: Index-Digest $idx"
}

# --- Scan (Punkt 16) -----------------------------------------------------------------------------------------------------------------------------------
TRIVY_IMAGE='aquasec/trivy@sha256:af6acf9a6b85dfe389a1941505c0ce9efef52a4719635e1a962f022a3d855daa' # 0.75.0
DB_GRENZE=$((72 * 3600)) # P16db-grenze
SCAN_PLATTFORMEN=(linux/amd64 linux/arm64)
# Ein Eintrag über der Schwelle hat den Schweregrad CRITICAL oder HIGH; er ist behebbar, wenn FixedVersion nicht leer ist (fehlt, null und "" sind ohne Behebung).
SCAN_EINTRAG='.Results[] | (.Vulnerabilities // [])[] | select(.Severity == "CRITICAL" or .Severity == "HIGH")' # P16schwelle
# Der Eintrag des Index, den Trivy für eine Plattform liest: das Manifest mit os linux und der Architektur ($a).
SCAN_WAHL='[.manifests[] | select(.platform.os == "linux" and .platform.architecture == $a)]' # P16wahl
SCAN_BEHOBEN='(.FixedVersion // "") != ""' # P16fix
# scan_fehler <Stichwort> <Text> — der Scan ist nicht gelaufen oder nicht auswertbar: Exit 2, kein Überspringen.
scan_fehler() { echo "release: scan: $1: $2" >&2; exit 2; }
# db_alter <Cache-Verzeichnis> — setzt alter auf das Alter der Schwachstellen-Datenbank in Sekunden (Feld UpdatedAt von db/metadata.json); Exit 2 `datenbank`, wenn es nicht lesbar ist oder mehr als 5 Minuten in der Zukunft liegt.
alter=0
db_alter() {
  local t ts
  t="$(jq -r '.UpdatedAt | strings' "$1/db/metadata.json" 2> /dev/null)"
  [[ $t =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?(Z|[+-][0-9]{2}:[0-9]{2})$ ]] || scan_fehler datenbank "das Alter der Datenbank ist nicht lesbar (db/metadata.json, UpdatedAt ${t@Q})" # P16db-zeitpunkt
  ts="$(date -u -d "$t" +%s 2> /dev/null)" || scan_fehler datenbank "UpdatedAt ${t@Q} ist kein Zeitpunkt" # P16db-datum
  alter=$(( $(date +%s) - ts ))
  [ "$alter" -ge -300 ] || scan_fehler datenbank "UpdatedAt ${t@Q} liegt mehr als 5 Minuten in der Zukunft" # P16db-zukunft
}

# waehle_manifest <Plattform> — legt $scratch/layout-<Architektur> an (oci-layout, blobs/ als Hardlinks, index.json mit genau dem Deskriptor des Manifests der Plattform) und setzt pdir; Exit 2 `plattform`, wenn der Index das Manifest nicht oder mehr als einmal nennt. Trivy liest ein Layout nur als Verzeichnis und bei zwei Manifesten das erste.
pdir=""; lay=""; indexblob=""
waehle_manifest() {
  local a="${1#linux/}" n
  n="$(jq --arg a "$a" "$SCAN_WAHL | length" "$indexblob" 2> /dev/null)"
  [ "$n" = 1 ] || scan_fehler plattform "der Index des Layouts nennt das Manifest für $1 ${n:-0}-mal statt einmal" # P16genau-eines
  pdir="$scratch/layout-$a"
  mkdir -p "$pdir"
  cp "$lay/oci-layout" "$pdir/oci-layout" # P16oci-layout
  cp -al "$lay/blobs" "$pdir/blobs" # P16blobs
  jq -c --arg a "$a" "{schemaVersion: 2, manifests: $SCAN_WAHL}" "$indexblob" > "$pdir/index.json"
}

schritt_scan() {
  local plat arch cache json erg b bef="" n m=0 z dg
  zustand
  verlange image/digest image # P16image
  [ -f "$zdir/image/layout.tar" ] || eingabe vorbedingung "das Layout des Schrittes image fehlt" # P16layout
  rm -f "$zdir/scan.ok" # P16marke-loescht
  cache="$scratch/trivy"; mkdir -p "$cache"
  lay="$scratch/layout"; mkdir -p "$lay"
  tar -xf "$zdir/image/layout.tar" -C "$lay" 2> /dev/null || scan_fehler scan-gescheitert "image/layout.tar ist nicht entpackbar" # P16entpacken
  dg="$(jq -r '.manifests[0].digest // empty' "$lay/index.json" 2> /dev/null)"
  indexblob="$lay/blobs/${dg/://}"
  { [ -n "$dg" ] && jq -e . "$indexblob" > /dev/null 2>&1; } || scan_fehler scan-gescheitert "index.json oder der Index-Blob des Layouts ist nicht lesbar" # P16index
  for plat in "${SCAN_PLATTFORMEN[@]}"; do
    json="$scratch/scan-${plat#linux/}.json"
    waehle_manifest "$plat" # P16auswahl
    fuehre docker run --rm --user "$(id -u):$(id -g)" -v "$pdir:/layout:ro" -v "$cache:/cache" "$TRIVY_IMAGE" \
      image --input /layout --format json --exit-code 0 --severity HIGH,CRITICAL --timeout 10m --cache-dir /cache > "$json" 2> "$scratch/trivy.err" # P16aufruf
    if [ "$rc" -ne 0 ]; then
      tail -n 20 "$scratch/trivy.err" | sed 's/^/release: trivy: /' >&2
      scan_fehler scan-gescheitert "Trivy für $plat endet mit Exit $rc" # P16exit
    fi
    jq -e '(.Results // []) | (type == "array") and (length > 0)' "$json" > /dev/null 2>&1 || scan_fehler pruefmenge "die Ausgabe von Trivy für $plat ist kein JSON oder trägt kein Ergebnis" # P16pruefmenge
    db_alter "$cache"
    [ "$alter" -le "$DB_GRENZE" ] || scan_fehler datenbank "die Datenbank ist $((alter / 3600)) Stunden alt, die Grenze sind 72" # P16db-alt
    arch="$(jq -r '.Metadata.ImageConfig.architecture // ""' "$json" 2> /dev/null)"
    [ "$arch" = "${plat#linux/}" ] || scan_fehler plattform "gescannt wurde die Architektur '$arch' statt ${plat#linux/}" # P16plattform
    erg="$(jq -r --arg p "$plat" "$SCAN_EINTRAG | if $SCAN_BEHOBEN"' then "F \($p) \(.VulnerabilityID) \(.PkgName) \(.InstalledVersion) -> \(.FixedVersion) \(.Severity)" else "O" end' "$json" 2> /dev/null)" || scan_fehler pruefmenge "die Ausgabe von Trivy für $plat ist nicht auswertbar" # P16auswertung
    b="$(grep '^F ' <<< "$erg" | cut -c3-)"
    m=$((m + $(grep -c '^O$' <<< "$erg" || true))) # P16zaehlung
    [ -z "$b" ] || bef+="$b"$'\n'
  done
  if [ -n "$bef" ]; then
    n="$(grep -c . <<< "$bef")"
    while IFS= read -r z; do [ -z "$z" ] || echo "release: scan: schwachstelle: $z" >&2; done <<< "$bef" # P16alle
    echo "release: scan: $n Befund(e) über der Schwelle, $m ohne Behebung" >&2 # P16summe
    exit 1
  fi
  : > "$zdir/scan.ok" # P16marke
  echo "release: scan: 0 Befund(e) über der Schwelle, $m ohne Behebung" # P16gruen-zeile
}

schritt_veroeffentlichen() {
  local reg v d have
  aussenwirkung_erlaubt # P2veroeffentlichen
  zustand
  verlange image/digest image # P13veroeffentlichen
  verlange scan.ok scan # P16veroeffentlichen
  tag_ableiten
  v="${tag#v}"
  d="$(cat "$zdir/image/digest")"
  for reg in "${REGISTRIES[@]}"; do # P6reg
    anmelden "$reg"
    fuehre skopeo inspect --raw --authfile "$scratch/auth.json" "docker://$reg/$IMAGE_NAME:$v" > "$scratch/roh" 2> "$scratch/skopeo.err"
    if [ "$rc" -ne 0 ]; then sed 's/^/release: skopeo: /' "$scratch/skopeo.err" >&2; fi # P6meldung
    if [ "$rc" -eq 0 ]; then
      have="sha256:$(sha256sum "$scratch/roh" | cut -d' ' -f1)"
      [ "$have" = "$d" ] || rot "$reg trägt $v mit dem Digest $have statt $d; der Weg ist ein höherer Tag" # P6anders
      echo "release: veroeffentlichen: $reg trägt $v schon mit $d"
      continue # P6gleich
    fi
    grep -qiE "$FEHLT_RE" "$scratch/skopeo.err" || rot "skopeo inspect für $reg endet mit Exit $rc" # P6fehler
    fuehre skopeo copy --all --preserve-digests --authfile "$scratch/auth.json" "oci-archive:$zdir/image/layout.tar" "docker://$reg/$IMAGE_NAME:$v" # P6copy
    [ "$rc" -eq 0 ] || rot "skopeo copy nach $reg endet mit Exit $rc; $(zustand_tag "$v")"
    echo "release: veroeffentlichen: $v nach $reg kopiert"
  done
}

schritt_vergleichen() {
  local reg v d dg dd
  aussenwirkung_erlaubt # P2vergleichen
  zustand
  verlange image/digest image # P13vergleichen
  rm -f "$zdir/vergleichen.ok" # P11vergleichen-marke
  tag_ableiten
  v="${tag#v}"
  d="$(cat "$zdir/image/digest")"
  for reg in "${REGISTRIES[@]}"; do
    fuehre skopeo inspect --raw --no-creds "docker://$reg/$IMAGE_NAME:$v" > "$scratch/roh-$reg" 2> "$scratch/skopeo.err" # P7anonym
    if [ "$rc" -ne 0 ]; then
      sed 's/^/release: skopeo: /' "$scratch/skopeo.err" >&2 # P7meldung
      if grep -qiE "$FEHLT_RE" "$scratch/skopeo.err"; then rot "$reg trägt den Tag $v nicht"; fi # P7tag
      rot "$reg ist ohne Anmeldung nicht lesbar (Exit $rc); das Paket muss öffentlich sein (Handlung des Betreibers)"
    fi
    index_pruefen "$scratch/roh-$reg" "der Index aus $reg" # P7plattformen
  done
  dg="sha256:$(sha256sum "$scratch/roh-ghcr.io" | cut -d' ' -f1)"
  dd="sha256:$(sha256sum "$scratch/roh-docker.io" | cut -d' ' -f1)"
  [ "$dg" = "$d" ] || rot "ghcr.io weicht ab: $dg statt $d" # P7ghcr
  [ "$dd" = "$d" ] || rot "docker.io weicht ab: $dd statt $d" # P7docker
  : > "$zdir/vergleichen.ok"
  echo "release: vergleichen: $d in ghcr.io und docker.io"
}

schritt_smoke() {
  local reg v d pg
  aussenwirkung_erlaubt # P2smoke
  zustand
  verlange vergleichen.ok vergleichen # P13smoke
  rm -f "$zdir/smoke.ok" # P11smoke-marke
  tag_ableiten
  v="${tag#v}"
  d="$(cat "$zdir/image/digest")"
  pg="$(awk '/^POSTGRES_IMAGE[ \t]*\?=/ { sub(/^[^=]*=[ \t]*/, ""); print; exit }' harness/mk/integration.mk 2> /dev/null)"
  [ -n "$pg" ] || eingabe vorbedingung "POSTGRES_IMAGE steht nicht in harness/mk/integration.mk"
  for reg in "${REGISTRIES[@]}"; do
    mkdir "$scratch/smoke-$reg"
    fuehre env VERSION="$v" POSTGRES_IMAGE="$pg" SMOKE_LAUF="pgr-rel-$$-${reg%%.*}" SMOKE_WERK="$scratch/smoke-$reg" SMOKE_NAME="release: smoke" \
      bash tools/harness/container-smoke.sh "$reg/$IMAGE_NAME@$d" # P7smoke
    [ "$rc" -eq 0 ] || rot "der Smoke aus $reg endet mit Exit $rc"
    echo "release: smoke: grün aus $reg"
  done
  : > "$zdir/smoke.ok" # P7smoke-ok
}

schritt_latest() {
  local reg v d hoechste
  aussenwirkung_erlaubt # P2latest
  zustand
  verlange smoke.ok smoke # P8vergleich
  tag_ableiten
  v="${tag#v}"
  d="$(cat "$zdir/image/digest")"
  if [[ $v == *-* ]]; then echo "release: latest: $tag ist eine Vorabversion, latest bleibt unverändert"; return 0; fi # P8stabil
  hoechste="$({ git tag -l 'v*' | grep -E "^v$KERN\$" | sed 's/^v//'; echo "$v"; } | sort -V | tail -n 1)" # P8hoechste
  if [ "$hoechste" != "$v" ]; then echo "release: latest: $hoechste ist höher als $v, latest bleibt unverändert"; return 0; fi
  for reg in "${REGISTRIES[@]}"; do # P8reg
    anmelden "$reg"
    fuehre skopeo copy --all --preserve-digests --authfile "$scratch/auth.json" "docker://$reg/$IMAGE_NAME@$d" "docker://$reg/$IMAGE_NAME:latest" # P8copy
    [ "$rc" -eq 0 ] || rot "latest in $reg nicht gesetzt (Exit $rc); $(zustand_tag latest)"
    echo "release: latest: $reg zeigt auf $d"
  done
}

# --- Homebrew (Punkt 19) ---------------------------------------------------------------------------------------------------------------------------------
TAP_URL_VORGABE='https://github.com/pt9912/homebrew-pgwire-recorder.git'
TAP_NAME='github-actions[bot]'
TAP_MAIL='41898282+github-actions[bot]@users.noreply.github.com'
FORMEL=Formula/pgwire-recorder.rb
VORLAGE=tools/harness/homebrew-formula.rb.tmpl
PLATTFORMEN_TAP=(darwin-amd64 darwin-arm64 linux-amd64 linux-arm64)
# homebrew_fehler <Stichwort> <Text> — Exit 1 mit dem Stichwort (Punkt 19).
homebrew_fehler() { echo "release: homebrew: $1: $2" >&2; exit 1; }
# tapgit <Argument>… — git gegen den Tap: der Token nur über GIT_ASKPASS, keine Abfrage am Terminal, kein Credential-Helper.
tapgit() { GIT_ASKPASS="$scratch/askpass.sh" GIT_TERMINAL_PROMPT=0 git -c credential.helper= "$@"; } # P19kein-prompt
# git_meldung — die Fehlerausgabe von git aus $scratch/git.err mit dem Präfix `release: git: `.
git_meldung() { sed 's/^/release: git: /' "$scratch/git.err" >&2; }

schritt_homebrew() {
  local v sums f k kk s c zweig tv hoch
  local -A summe
  aussenwirkung_erlaubt # P19ereignis
  tag_ableiten
  pruefe_tag_form # P19tag-form
  v="${tag#v}"
  if [[ $tag == *-* ]]; then echo "release: homebrew: $tag ist eine Vorabversion, der Tap bleibt unverändert"; return 0; fi # P19vorabversion
  [ -n "${GH_TOKEN-}" ] || eingabe token "GH_TOKEN fehlt" # P19gh-token
  [ -n "${RELEASE_TAP_TOKEN-}" ] || eingabe geheimnis "HOMEBREW_TAP_GITHUB_TOKEN ist nicht gesetzt oder leer" # P19tap-token
  # (3) Summen aus dem Release
  mkdir "$scratch/sums"
  fuehre gh release download "$tag" --repo "$RELEASE_REPOSITORY" --pattern SHA256SUMS --dir "$scratch/sums"
  [ "$rc" -eq 0 ] || homebrew_fehler summen-laden "gh release download endet mit Exit $rc" # P19summen-laden
  sums="$scratch/sums/SHA256SUMS"
  for k in "${PLATTFORMEN_TAP[@]}"; do
    c="$(awk -v n="pgwire-recorder-$k" '$2 == n { c++ } END { print c + 0 }' "$sums" 2> /dev/null)"
    [ "$c" -eq 1 ] || homebrew_fehler summen "SHA256SUMS nennt pgwire-recorder-$k ${c:-0}-mal statt einmal"
    s="$(awk -v n="pgwire-recorder-$k" 'NF == 2 && $2 == n { print $1; exit }' "$sums")"
    [[ $s =~ ^[0-9a-f]{64}$ ]] || homebrew_fehler summen "die Summe für pgwire-recorder-$k ist keine SHA-256-Summe" # P19summen-form
    summe[$k]="$s"
  done
  # (4) Formel aus der Vorlage
  f="$(cat "$VORLAGE" 2> /dev/null; printf x)"
  f="${f%x}"
  [ -n "$f" ] || homebrew_fehler vorlage "die Vorlage $VORLAGE fehlt oder ist leer"
  f="${f//@@VERSION@@/$v}"
  for k in "${PLATTFORMEN_TAP[@]}"; do
    kk="${k^^}"
    f="${f//@@SHA256_${kk//-/_}@@/${summe[$k]-}}" # P19summe
  done
  printf '%s' "$f" > "$scratch/formel.rb"
  # (5) Tap
  printf '%s\n' '#!/bin/sh' 'case "$1" in Username*) printf "%s\n" x-access-token ;; *) printf "%s\n" "$RELEASE_TAP_TOKEN" ;; esac' > "$scratch/askpass.sh"
  chmod +x "$scratch/askpass.sh"
  fuehre tapgit clone --quiet "${RELEASE_TAP_URL:-$TAP_URL_VORGABE}" "$scratch/tap" 2> "$scratch/git.err" # P19klon
  [ "$rc" -eq 0 ] || { git_meldung; homebrew_fehler tap-klon "der Tap ist nicht klonbar (Exit $rc)"; } # P19klon-pruefung
  cd "$scratch/tap" || exit 1
  if tapgit rev-parse -q --verify HEAD > /dev/null 2>&1; then
    zweig="$(tapgit symbolic-ref -q --short HEAD 2> /dev/null)"
    [ "$zweig" = main ] || homebrew_fehler tap-zweig "der Tap ist auf '${zweig:-HEAD ohne Zweig}' ausgecheckt statt auf main" # P19zweig
  else
    tapgit symbolic-ref HEAD refs/heads/main # P19leer
  fi
  if [ -e Formula ] || [ -L Formula ]; then
    { [ -d Formula ] && [ ! -L Formula ]; } || homebrew_fehler tap-pfad "Formula im Tap ist kein Verzeichnis" # P19pfad-verzeichnis
  fi
  if [ -e "$FORMEL" ] || [ -L "$FORMEL" ]; then
    { [ -f "$FORMEL" ] && [ ! -L "$FORMEL" ]; } || homebrew_fehler tap-pfad "$FORMEL im Tap ist keine reguläre Datei" # P19pfad-datei
    mapfile -t vz < <(grep -a -E '^  version "[^"]*"$' "$FORMEL")
    [ "${#vz[@]}" -eq 1 ] || homebrew_fehler tap-formel "die Formel im Tap hat ${#vz[@]} Zeilen version statt einer" # P19formel-zeilen
    tv="${vz[0]#  version \"}"
    tv="${tv%\"}"
    [[ $tv =~ ^$KERN$ ]] || homebrew_fehler tap-formel "die Version ${tv@Q} der Formel im Tap ist keine stabile Version" # P19formel-stabil
    if [ "$tv" != "$v" ]; then
      hoch="$(printf '%s\n%s\n' "$tv" "$v" | sort -V | tail -n 1)" # P19sort
      if [ "$hoch" = "$tv" ]; then echo "release: homebrew: der Tap trägt $tv, höher als $v, und bleibt unverändert"; return 0; fi # P19hoeher
    elif cmp -s "$scratch/formel.rb" "$FORMEL"; then
      echo "release: homebrew: der Tap trägt $v schon (schon aktuell)"
      return 0
    fi
  fi
  mkdir -p Formula
  cp "$scratch/formel.rb" "$FORMEL"
  tapgit add -- "$FORMEL" || rot "git add im Tap scheitert"
  GIT_AUTHOR_NAME="$TAP_NAME" GIT_AUTHOR_EMAIL="$TAP_MAIL" GIT_COMMITTER_NAME="$TAP_NAME" GIT_COMMITTER_EMAIL="$TAP_MAIL" tapgit commit --quiet -m "pgwire-recorder $v" || rot "git commit im Tap scheitert" # P19commit
  fuehre tapgit push origin main 2> "$scratch/git.err" # P19push
  [ "$rc" -eq 0 ] || { git_meldung; homebrew_fehler tap-push "der Tap weist den Push ab (Exit $rc)"; } # P19push-meldung
  echo "release: homebrew: Formel $v im Tap abgelegt"
}

case "$schritt" in
  eingaben) schritt_eingaben ;;
  binaries) schritt_binaries ;;
  image) schritt_image ;;
  scan) schritt_scan ;;
  veroeffentlichen) schritt_veroeffentlichen ;;
  vergleichen) schritt_vergleichen ;;
  smoke) schritt_smoke ;;
  latest) schritt_latest ;;
  github-release) schritt_github_release ;;
  homebrew) schritt_homebrew ;;
  *) echo "release: unbekannter Schritt '$schritt'" >&2; exit 2 ;;
esac
