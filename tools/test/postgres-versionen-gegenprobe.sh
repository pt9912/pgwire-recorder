#!/usr/bin/env bash
# postgres-versionen-gegenprobe — prüft die Versionsmatrix (harness/mk/postgres-versionen.mk,
# tools/test/run-postgres-versionen.sh, tools/test/run-integration-tests.sh) gegen die Punkte 1
# bis 9 und 12 von SPEC-053 (harness/vertraege/SPEC-053-postgres-versionen.md) in Kopien des
# Arbeitsbaums in einem Temp-Verzeichnis; der Arbeitsbaum bleibt unberührt. docker ist eine
# Attrappe (über DOCKER): Sie protokolliert jeden Aufruf, führt Buch über Container, Netz und
# Volume und liefert Version und Testausgabe nach Vorgabe des Falls. Ein sleep in PATH
# protokolliert die Wartezeit, ohne zu warten. Der Name jedes Falls trägt die Nummer seines
# Punktes; die unveränderte Kopie besteht jeden Fall, und bei jeder Mutation (Textersatz in der
# Kopie) werden von den Fällen ihrer Gruppe (`nur`) genau die benannten rot.
#
# GRENZE der Gegenprobe: Mit der Attrappe läuft kein PostgreSQL und kein Test; ob die Teilmenge
# gegen einen echten Server besteht und ob `-test.run` die Tests trifft, zeigt nur
# `make test-postgres-versionen`. Punkt 11 (Test der Lebendprüfung) ist ein Go-Test ohne Fall
# hier, Punkt 10 und 14 sind Bindungen ohne Fall. Die Signale gehen an das Skript (`skript`)
# oder an die Prozessgruppe (`gruppe`) des Ziels, nicht an `make`; sie setzen voraus, dass diese
# Gegenprobe SIGINT nicht ignoriert (kein Start als Hintergrundjob).
#
# Ausgang: 0, wenn jeder Fall und jede Mutation wie erwartet endet, sonst 1.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
wurzel="$(cd "$hier/../.." && pwd)"
arbeit="$(mktemp -d)"
trap 'rm -rf "${arbeit:?}"' EXIT
fehler=0
melde() { echo "postgres-versionen-gegenprobe: ROT — $*" >&2; fehler=1; }

bin="$arbeit/bin" log="$arbeit/log" st="$arbeit/zustand" aus="$arbeit/aus" schlaf="$arbeit/schlaf"
mkdir -p "$bin"

# Attrappe von docker. Vorgaben des Falls (Umgebung): AT_ROT="<phase>-<version> …",
# AT_UMBENANNT und AT_SKIP=<test>, AT_KEINE=1, AT_NICHT_BEREIT, AT_PULL_ROT und AT_HBA_ROT=<version>,
# AT_VERSION_<n>=<server_version>, AT_BAUEN_ROT=1, AT_HANGT=<datei> (der Testlauf
# wartet, bis sein Container entfernt ist).
cat > "$bin/docker" <<'ATTRAPPE'
#!/usr/bin/env bash
printf '%s' "$*" | tr '\n' ' ' >> "$AT_LOG"; echo >> "$AT_LOG"
sub="$1"; shift
major() { local n="${1%-pg}"; n="${n%-tests}"; echo "${n##*-}"; }
case "$sub" in
  build) [ -z "${AT_BAUEN_ROT:-}" ] ;;
  network | volume)
    if [ "$1" = create ]; then : > "$AT_STATE/${sub:0:1}-${*: -1}"; else rm -f "${AT_STATE:?}/${sub:0:1}-${2:?}"; fi ;;
  rm) for n in "${@:2}"; do rm -f "${AT_STATE:?}/c-${n:?}"; done ;;
  stop | logs) ;;
  exec)
    if [ "$1" = -i ]; then
      cat > "$AT_LOG.sql-$4"
      [[ " ${AT_HBA_ROT:-} " != *" $(major "$4") "* ]]
    elif [ "$2" = pg_isready ]; then [[ " ${AT_NICHT_BEREIT:-} " != *" $(major "$1") "* ]]
    else m="$(major "$1")"; v="AT_VERSION_$m"; echo "${!v-$m.1}"; fi ;;
  run)
    name="$3"; m="$(major "$name")"
    if [ "$1" = -d ]; then
      [[ " ${AT_PULL_ROT:-} " != *" $m "* ]] || exit 1
      if ls "$AT_STATE"/c-*-pg > /dev/null 2>&1; then echo PARALLEL >> "$AT_LOG"; fi
      echo "$*" > "$AT_STATE/c-$name"; exit 0
    fi
    phase=1; [[ "$*" != *PGR_OHNE_POSTGRES=1* ]] || phase=2
    echo "$*" > "$AT_STATE/c-$name"
    if [ -n "${AT_HANGT:-}" ]; then
      : > "$AT_HANGT"
      for _ in $(seq 200); do [ -e "$AT_STATE/c-$name" ] || break; /bin/sleep 0.05; done; exit 1
    fi
    if [ -n "${AT_KEINE:-}" ]; then echo "testing: warning: no tests to run"; echo PASS
    else
      v="AT_NAMEN_$phase"
      for t in ${!v:-}; do
        if [ "$t" = "${AT_UMBENANNT:-}" ]; then continue; fi
        if [ "$t" = "${AT_SKIP:-}" ]; then echo "--- SKIP: $t (0.00s)"; continue; fi
        echo "--- PASS: ${t%%/*} (0.00s)"
        if [ "$t" != "${t%%/*}" ]; then echo "    --- PASS: $t (0.00s)"; fi
      done
    fi
    rm -f "${AT_STATE:?}/c-${name:?}"
    [[ " ${AT_ROT:-} " != *" $phase-$m "* ]] || { echo "--- FAIL: Test (0.00s)"; exit 1; } ;;
esac
ATTRAPPE
printf '#!/usr/bin/env bash\necho "$1" >> "%s"\n' "$schlaf" > "$bin/sleep"
chmod +x "$bin/docker" "$bin/sleep"

# make_variable <wurzel> <name> — Wörter der Make-Variable aus der Datenbank von make.
make_variable() {
  { make -s --no-print-directory -C "$1" -pq help 2> /dev/null || true; } \
    | awk -v n="$2" '$1 == n && ($2 == ":=" || $2 == "=") { for (i = 3; i <= NF; i++) print $i }'
}
namen1="$(make_variable "$wurzel" POSTGRES_TESTS_PHASE1 | tr '\n' ' ')"
namen2="$(make_variable "$wurzel" POSTGRES_TESTS_PHASE2 | tr '\n' ' ')"
refbild="$(make_variable "$wurzel" POSTGRES_IMAGE)"
D="$(printf 'a%.0s' {1..64})"
v14="14=postgres:14-alpine@sha256:$D" v15="15=postgres:15-alpine@sha256:$D"
v16="16=postgres:16-alpine@sha256:$D" v17="17=$refbild" v18="18=postgres:18-alpine@sha256:$D"
alle="$v14 $v15 $v16 $v17 $v18"
liste=""
q="'"
tests=($namen1 $namen2)
FIXTURE_ZAEHLER=0

# fixture — legt eine Kopie an und setzt w auf ihre Wurzel.
fixture() {
  FIXTURE_ZAEHLER=$((FIXTURE_ZAEHLER + 1))
  w="$arbeit/baum-$FIXTURE_ZAEHLER"
  mkdir -p "$w/harness/mk" "$w/tools/test" "$w/.harness/state"
  cp "$wurzel/Makefile" "$w/Makefile"
  cp "$wurzel/harness/mk/integration.mk" "$wurzel/harness/mk/postgres-versionen.mk" "$w/harness/mk/"
  cp "$wurzel/tools/test/run-integration-tests.sh" "$wurzel/tools/test/run-postgres-versionen.sh" "$w/tools/test/"
}

# lauf [VAR=wert…] -- <befehl>… — führt den Befehl in $w aus; setzt code, Ausgabe in $aus.
lauf() {
  local e=()
  while [ "$1" != -- ]; do e+=("$1"); shift; done
  shift
  rm -rf "${st:?}"; mkdir -p "$st"; : > "$log"; : > "$schlaf"; rm -f "$log".sql-*
  set +e
  {
    (cd "$w" && env -u MAKELEVEL -u MFLAGS MAKEFLAGS="${mf:-}" PATH="$bin:$PATH" AT_LOG="$log" AT_STATE="$st" \
      AT_NAMEN_1="$namen1" AT_NAMEN_2="$namen2" "${e[@]}" "$@") > "$aus" 2>&1
    code=$?
  } 2> /dev/null
  set -e
}
# matrix [VAR=wert…] [-- <make-flag>…] — make test-postgres-versionen mit der Liste $liste (sonst alle).
matrix() {
  local e=()
  while [ $# -gt 0 ] && [ "$1" != -- ]; do e+=("$1"); shift; done
  if [ $# -gt 0 ]; then shift; fi
  lauf "${e[@]}" -- make -s --no-print-directory "$@" test-postgres-versionen \
    POSTGRES_VERSIONEN="${liste:-$alle}" DOCKER="$bin/docker"
}
in_aus() { grep -qF -- "$1" "$aus"; }
anzahl() { grep -cE -- "$1" "$log" || true; }
zeile() { grep -nE -- "$1" "$log" | head -1 | cut -d: -f1; }
leer() { [ -z "$(ls -A "$st")" ]; }
schlafsumme() { awk '{ s += $1 } END { print s + 0 }' "$schlaf"; }
dateien() { (cd "$w" && find . -type f | sort); }
sqldatei() { echo "$log".sql-pgr-pv-*-"$1"-pg; }
execzeile() { grep -E "^exec -i .*-$1-pg " "$log"; }

# sel <fall> — läuft der Fall? ($nur: Namensanfänge, leer = alle)
sel() {
  local p
  if [ -z "${nur:-}" ]; then return 0; fi
  for p in $nur; do if [[ "$1" == "$p"* ]]; then return 0; fi; done
  return 1
}
# ist <fall> <bedingung> — schreibt den Namen, wenn der Fall läuft und die Bedingung nicht gilt.
ist() { if sel "$1"; then eval "$2" || echo "$1"; fi; return 0; }
# lauft <fall>… — läuft einer von ihnen?
lauft() { local f; for f; do if sel "$f"; then return 0; fi; done; return 1; }

eingaben=(
  "p2a-leere-liste|Liste der Versionen ist leer|"
  "p2b-referenz-fehlt|Referenz '17' fehlt|$v14 $v15"
  "p2c-doppelt|Hauptversion doppelt im Eintrag 15=|$v14 $v15 15=postgres:15-alpine@sha256:$D $v17"
  "p2c-referenz-doppelt|Hauptversion der Referenz doppelt|$v14 $v17 $v17"
  "p2d-keine-andere|keine andere Version|$v17"
  "p2e1-ohne-gleich|Eintrag ohne '=': 14postgres|14postgres $v17"
  "p2e2-leere-hauptversion|leere Hauptversion|=postgres:14-alpine@sha256:$D $v17"
  "p2e3-nicht-ziffern|nicht aus Ziffern|x14=postgres:14-alpine@sha256:$D $v17"
  "p2e4-fuehrende-null|führender Null|014=postgres:14-alpine@sha256:$D $v17"
  "p2f-ohne-tag|Image ohne Tag|14=postgres@sha256:$D $v17"
  "p2g-ohne-digest|Image ohne Digest|14=postgres:14-alpine $v17"
  "p2h-digest-kurz|Digest ungültig|14=postgres:14-alpine@sha256:${D:1} $v17"
  "p2h-digest-lang|Digest ungültig|14=postgres:14-alpine@sha256:${D}a $v17"
  "p2h-digest-gross|Digest ungültig|14=postgres:14-alpine@sha256:${D/a/A} $v17"
  "p2h-digest-kein-hex|Digest ungültig|14=postgres:14-alpine@sha256:${D/a/g} $v17"
  "p2h-digest-verfahren|Digest ungültig|14=postgres:14-alpine@sha512:$D $v17"
  "p2i-referenz-image|ungleich POSTGRES_IMAGE|$v14 17=postgres:17-alpine@sha256:$D"
)
p2namen=()
for a in "${eingaben[@]}"; do p2namen+=("${a%%|*}"); done

# muster <name>… — das Muster von -test.run, wie der Runner es bilden soll.
muster() {
  local n top="" sub=""
  for n; do
    top+="${n%%/*}|"
    if [ "$n" != "${n%%/*}" ]; then sub+="${n#*/}|"; fi
  done
  printf '^(%s)$' "${top%|}"
  if [ -n "$sub" ]; then printf '/^(%s)$' "${sub%|}"; fi
}
auswahl_ok() {
  grep -qF -- "-test.run $(muster $namen1)" "$log" && grep -qF -- "-test.run $(muster $namen2)" "$log"
}
# phasenzeile <1|2> — der Aufruf des Testlaufs der Phase.
phasenzeile() { if [ "$1" = 2 ]; then grep -E '^run --rm .*PGR_OHNE_POSTGRES=1' "$log"; else grep -E '^run --rm .*PGR_PW_SCRAM' "$log"; fi; }
# elementfall <test> — fehlt seine PASS-Zeile, ist die Version rot, und die Auswahl nennt ihn.
elementfall() {
  local t="$1" i=1
  if [[ " $namen2 " == *" $t "* ]]; then i=2; fi
  liste="$v17 $v18"
  matrix AT_UMBENANNT="$t"
  liste=""
  in_aus "Version 18, Phase $i: Test $t ohne PASS-Zeile" && in_aus "ROT — 18 (Phase $i)" \
    && phasenzeile "$i" | grep -qF -- "${t%%/*}" \
    && { [ "$t" = "${t%%/*}" ] || phasenzeile "$i" | grep -qF -- "${t#*/}"; }
}
# vor_naechster — Container, Netz und Volume einer Version sind entfernt, bevor die nächste startet.
vor_naechster() {
  local a b
  for p in "14 15" "15 16" "16 18"; do
    read -r a b <<< "$p"
    [ "$(zeile "^rm -f pgr-pv-[0-9]+-$a-tests")" -lt "$(zeile "^run -d --name pgr-pv-[0-9]+-$b-pg")" ] \
      && [ "$(zeile "^network rm pgr-pv-[0-9]+-$a-net")" -lt "$(zeile "^network create --internal pgr-pv-[0-9]+-$b-net")" ] \
      && [ "$(zeile "^volume rm pgr-pv-[0-9]+-$a-daten")" -lt "$(zeile "^volume create pgr-pv-[0-9]+-$b-daten")" ] || return 1
  done
}
# gestartet <version>… — jede Version hat genau einen Server gestartet.
gestartet() {
  local v
  for v; do [ "$(anzahl "^run -d --name [^ ]*-$v-pg")" -eq 1 ] || return 1; done
}
namen_ok() {
  local v
  for v in 14 15 16 18; do
    [ "$(anzahl "^network create --internal pgr-pv-[0-9]+-$v-net$")" -eq 1 ] \
      && [ "$(anzahl "^volume create pgr-pv-[0-9]+-$v-daten$")" -eq 1 ] \
      && [ "$(anzahl "^run -d --name pgr-pv-[0-9]+-$v-pg ")" -eq 1 ] || return 1
  done
}
# signalfall <INT|TERM> <gruppe|skript> — das Ziel hängt im Testlauf; das Signal geht an die
# Prozessgruppe seiner eigenen Sitzung oder nur an das Skript.
signalfall() {
  local marke="$arbeit/hangt" pidf="$arbeit/pid" wp
  rm -f "$marke" "$pidf"
  (
    until [ -e "$marke" ] && [ -s "$pidf" ]; do /bin/sleep 0.05; done
    if [ "$2" = gruppe ]; then kill -"$1" -- "-$(cat "$pidf")"; else kill -"$1" "$(cat "$pidf")"; fi
  ) &
  wp=$!
  lauf POSTGRES_VERSIONEN="$alle" POSTGRES_REFERENZ=17 POSTGRES_IMAGE="$refbild" POSTGRES_TESTS_PHASE1="$namen1" \
    POSTGRES_TESTS_PHASE2="$namen2" DOCKER="$bin/docker" AT_HANGT="$marke" -- \
    setsid bash -c 'echo $$ > "$1"; shift; exec "$@"' x "$pidf" bash tools/test/run-postgres-versionen.sh
  wait "$wp" || true
}
# p5_<…> <version> — Anmeldung je Version.
p5_benutzer() { grep -qF "CREATE ROLE play_scram LOGIN PASSWORD 'GEHEIM scram \$%41'" "$(sqldatei "$1")"; }
p5_methode() { grep -qF "SET password_encryption = 'scram-sha-256'" "$(sqldatei "$1")" && execzeile "$1" | grep -qF 'play_scram all scram-sha-256'; }
p5_hba_vor() { execzeile "$1" | grep -qE 'play_scram all scram-sha-256".*; cat "\$hba"'; }
p5_kein_md5() { ! execzeile "$1" | grep -qE 'play_md5|play_pw|play_reject' && ! grep -qiE 'md5|play_pw' "$(sqldatei "$1")"; }

# pruefe <baum> — schreibt die Namen der Fälle, die nicht bestehen (einen je Zeile); $nur wählt aus.
pruefe() {
  w="$1"
  local a n text lst t v
  # Punkt 1: die Liste des Fragments ist gültig und trägt vier Versionen neben der Referenz.
  ist p1-liste-gueltig 'lauf -- make -s --no-print-directory test-postgres-versionen DOCKER="$bin/docker"; [ "$code" -eq 0 ] && [ "$(anzahl "^run -d")" -eq 4 ]'

  # Punkt 2: Exit 2, kein Aufruf von docker, der Text in der Ausgabe.
  for a in "${eingaben[@]}"; do
    IFS='|' read -r n text lst <<< "$a"
    if sel "$n"; then
      liste="${lst:- }"
      matrix
      { [ "$code" -ne 0 ] && [ ! -s "$log" ] && in_aus "$text"; } || echo "$n"
      liste=""
    fi
  done

  # Punkte 3, 5, 6, 7, 9: ein grüner Lauf über 14, 15, 16 und 18.
  if lauft p3-auswahl p5-benutzer p5-methode p5-hba-vor-allgemeiner p5-kein-md5-klartext p6-gruen p6-reihenfolge \
    p6-referenz-nicht p6-testimage-einmal p6-nacheinander p7-container p7-netz p7-volume p7-vor-naechster p7-namen p9-kein-stand; then
    local vor nach
    vor="$(dateien)"
    matrix
    nach="$(dateien)"
    ist p6-gruen '[ "$code" -eq 0 ] && in_aus "grün"'
    ist p6-reihenfolge '[ "$(zeile "run -d.*postgres:14")" -lt "$(zeile "run -d.*postgres:15")" ] && [ "$(zeile "run -d.*postgres:15")" -lt "$(zeile "run -d.*postgres:16")" ] && [ "$(zeile "run -d.*postgres:16")" -lt "$(zeile "run -d.*postgres:18")" ]'
    ist p6-referenz-nicht '[ "$(anzahl "^run -d")" -eq 4 ] && ! grep -qF "postgres:17" "$log"'
    ist p6-testimage-einmal '[ "$(anzahl "^build ")" -eq 1 ]'
    ist p6-nacheinander '! grep -q PARALLEL "$log"'
    ist p7-container '! ls "$st" | grep -q "^c-"'
    ist p7-netz '! ls "$st" | grep -q "^n-"'
    ist p7-volume '! ls "$st" | grep -q "^v-"'
    ist p7-vor-naechster 'vor_naechster 2> /dev/null'
    ist p7-namen 'namen_ok'
    ist p9-kein-stand '[ "$vor" = "$nach" ]'
    ist p3-auswahl 'auswahl_ok'
    for v in 14 15 16 18; do
      ist p5-benutzer "p5_benutzer $v"
      ist p5-methode "p5_methode $v"
      ist p5-hba-vor-allgemeiner "p5_hba_vor $v"
      ist p5-kein-md5-klartext "p5_kein_md5 $v"
    done
  fi

  # Punkt 3: je Test der Teilmenge — fehlt seine PASS-Zeile, ist die Version rot, und die Auswahl nennt ihn.
  for t in "${tests[@]}"; do ist "p3-element-$t" "elementfall '$t'"; done
  liste="$v17 $v18"
  ist p3-skip 'matrix AT_SKIP=TestE2EVorbereitungOhnePostgres; in_aus "Test TestE2EVorbereitungOhnePostgres übersprungen" && in_aus "ROT — 18 (Phase 1)"'
  ist p3-ohne-treffer 'matrix AT_KEINE=1; in_aus "Phase 1: Auswahl ohne Treffer" && in_aus "ROT — 18 (Phase 1)"'
  liste="$v17 $v16"
  ist p3-fehlt-phase2 'matrix AT_UMBENANNT=TestE2EOhnePostgresReplay; in_aus "Version 16, Phase 2: Test TestE2EOhnePostgresReplay ohne PASS-Zeile" && in_aus "ROT — 16 (Phase 2)" && ! in_aus "Phase 1"'
  liste=""

  # Punkt 4: der Server.
  ist p4-nicht-bereit 'matrix AT_NICHT_BEREIT=15; in_aus "Version 15: PostgreSQL nicht bereit" && in_aus "15 (Aufbau)"'
  ist p4-60-sekunden 'matrix AT_NICHT_BEREIT=16; [ "$(schlafsumme)" = 60 ] && in_aus "16 (Aufbau)"'
  ist p4-start 'matrix AT_PULL_ROT=16; in_aus "16 (Aufbau)" && [ "$(anzahl "run -d .*postgres:18")" -eq 1 ]'
  ist p4-version-abweichend 'matrix AT_VERSION_15=14.2; in_aus "Version 15: der Server meldet server_version ${q}14.2${q}" && in_aus "15 (Aufbau)"'
  ist p4-version-unlesbar 'matrix AT_VERSION_15=; in_aus "Version 15: der Server meldet server_version ${q}${q}" && in_aus "15 (Aufbau)"'
  ist p4-version-ohne-ziffer 'matrix AT_VERSION_15=x15; in_aus "Version 15: der Server meldet server_version ${q}x15${q}" && in_aus "15 (Aufbau)"'

  # Punkt 5: scheitert die Einrichtung, ist die Version rot.
  ist p5-einrichtung-scheitert 'matrix AT_HBA_ROT=16; in_aus "ROT — 16 (Aufbau)"'

  # Punkt 6: Ablauf und Meldung.
  ist p6-build-rot 'matrix AT_BAUEN_ROT=1; [ "$code" -ne 0 ] && in_aus "das Testimage ist nicht zu bauen" && [ "$(anzahl "^build ")" -eq 1 ] && [ "$(grep -c . "$log")" -eq 1 ]'
  ist p6-meldung-phase1 'matrix AT_ROT=1-15; in_aus "ROT — 15 (Phase 1)"'
  ist p6-keine-phase2 'matrix AT_ROT=1-15; [ "$(anzahl "^run --rm --name [^ ]*-15-tests .*PGR_OHNE_POSTGRES=1")" -eq 0 ] && [ "$(anzahl "^run --rm --name [^ ]*-16-tests .*PGR_OHNE_POSTGRES=1")" -eq 1 ]'
  ist p6-alle-versionen 'matrix AT_ROT=1-14; gestartet 15 16 18'
  ist p6-meldung-phase2 'matrix AT_ROT=2-16; in_aus "ROT — 16 (Phase 2)"'
  ist p6-meldung-aufbau 'matrix AT_PULL_ROT=16; in_aus "ROT — 16 (Aufbau)"'
  ist p6-zusammenfassung 'matrix AT_ROT="1-14 2-16"; in_aus "ROT — 14 (Phase 1) 16 (Phase 2)"'
  ist p6-exit 'matrix AT_ROT=1-15; [ "$code" -ne 0 ] && in_aus "ROT"'

  # Punkt 7: das Ende jedes Falls räumt auf.
  liste="$v17 $v15"
  ist p7-ende-rot 'matrix AT_ROT=1-15; leer'
  ist p7-ende-bereit 'matrix AT_NICHT_BEREIT=15; leer'
  ist p7-ende-pull 'matrix AT_PULL_ROT=15; leer'
  ist p7-ende-gruen 'matrix; leer'
  liste=""
  for a in INT:gruppe TERM:gruppe INT:skript TERM:skript; do
    ist "p7-sig${a%%:*}-${a##*:}" 'signalfall "${a%%:*}" "${a##*:}"; [ "$code" -ne 0 ] && leer && [ "$(anzahl "^run -d")" -eq 1 ]'
  done

  # Punkt 8: der Runner ohne Auswahl wie zuvor.
  if lauft p8-phase1 p8-phase2 p8-rollen p8-ohne-pruefung; then
    lauf POSTGRES_IMAGE=x DOCKER="$bin/docker" AT_NAMEN_1= AT_NAMEN_2= -- bash tools/test/run-integration-tests.sh
    ist p8-phase1 '[ "$(anzahl "^run --rm .*PGR_PW_KLARTEXT=[^ ]* pgwire-recorder:integration$")" -eq 1 ]'
    ist p8-phase2 '[ "$(anzahl "^run --rm .*PGR_OHNE_POSTGRES=1 pgwire-recorder:integration -test.run \^TestE2EOhnePostgres$")" -eq 1 ]'
    ist p8-rollen 'grep -qF "CREATE ROLE play_md5" "$log".sql-* && grep -qF "CREATE ROLE play_pw" "$log".sql-* && grep -E "^exec -i " "$log" | grep -qE "play_md5 all md5.*play_pw all password.*play_reject all reject"'
    ist p8-ohne-pruefung '[ "$code" -eq 0 ] && [ "$(anzahl "SHOW server_version")" -eq 0 ] && [ "$(anzahl "^build ")" -eq 1 ]'
  fi

  # Punkt 9: Aufrufarten des Treibers.
  ist p9-n 'matrix -- -n; [ "$code" -eq 0 ] && [ ! -s "$log" ]'
  ist p9-q 'matrix -- -q; [ ! -s "$log" ]'
  ist p9-t 'matrix -- -t; [ "$code" -eq 0 ] && [ ! -s "$log" ]'
  ist p9-makeflags-n 'mf=n matrix; [ "$code" -eq 0 ] && [ ! -s "$log" ]'
  ist p9-n-und-i 'matrix -- -n -i; [ "$code" -eq 0 ] && [ ! -s "$log" ]'
  ist p9-i 'matrix AT_ROT=1-14 -- -i; [ "$code" -eq 0 ] && in_aus "ROT — 14 (Phase 1)" && [ "$(anzahl "^run -d")" -eq 4 ]'

  # Punkt 12: beide Ziele hängen an GATE_CHECKS.
  ist p12-gate-checks-ziel '[[ " $(make_variable "$w" GATE_CHECKS | tr "\n" " ") " == *" test-postgres-versionen "* ]]'
  ist p12-gate-checks-gegenprobe '[[ " $(make_variable "$w" GATE_CHECKS | tr "\n" " ") " == *" test-postgres-versionen-gegenprobe "* ]]'
}

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
  if [[ "$alt" == %* ]]; then set -- "$1" "$2" "$3" "" "${@:4}"; neu=""; fi
  nur="$5"
  shift 5
  fixture
  mutiere "$datei" "$alt" "$neu"
  ist="$(pruefe "$w" | sort -u | tr '\n' ' ')"
  soll="$(printf '%s\n' "$@" | sort -u | tr '\n' ' ')"
  if [ "$ist" != "$soll" ]; then melde "Mutation '$name': erwartet rot: ${soll:-nichts}; rot sind: ${ist:-nichts}"; fi
}

D_=tools/test/run-postgres-versionen.sh R_=tools/test/run-integration-tests.sh F_=harness/mk/postgres-versionen.mk

# --- Die unveränderte Kopie besteht jeden Fall; das Gate hängt im Arbeitsbaum an GATE_CHECKS ---
nur=""
fixture
ist="$(pruefe "$w" | sort -u | tr '\n' ' ')"
if [ -n "$ist" ]; then melde "die unveränderte Kopie besteht nicht: $ist"; fi
gates=" $(make_variable "$wurzel" GATE_CHECKS | tr '\n' ' ') "
for g in test-postgres-versionen test-postgres-versionen-gegenprobe; do
  if [[ "$gates" != *" $g "* ]]; then melde "$g hängt im Arbeitsbaum nicht an GATE_CHECKS"; fi
done

# --- Punkt 1 und 2: Liste und Eingabe ---
mut p1-liste "$F_" '14-alpine@sha256:4ea9' '14-alpine@sha256:4EA9' p1 p1-liste-gueltig
mut p2a "$D_" %P2a p2 p2a-leere-liste
mut p2b "$D_" %P2b p2 p2b-referenz-fehlt
mut p2c "$D_" %P2c p2 p2c-doppelt
mut p2c-referenz "$D_" %P2c-ref p2 p2c-referenz-doppelt
mut p2d "$D_" %P2d p2 p2d-keine-andere
mut p2e1 "$D_" %P2e1 p2 p2e1-ohne-gleich
mut p2e2 "$D_" %P2e2 p2 p2e2-leere-hauptversion
mut p2e3 "$D_" %P2e3 p2 p2e3-nicht-ziffern
mut p2e4 "$D_" %P2e4 p2 p2e4-fuehrende-null
mut p2f "$D_" %P2f p2 p2f-ohne-tag
mut p2g "$D_" %P2g p2 p2g-ohne-digest
mut p2h-kurz "$D_" '{64}' '{1,64}' p2 p2h-digest-kurz
mut p2h-lang "$D_" '{64}' '{64,}' p2 p2h-digest-lang
mut p2h-gross "$D_" '[0-9a-f]{64}' '[0-9a-fA-F]{64}' p2 p2h-digest-gross
mut p2h-kein-hex "$D_" '[0-9a-f]{64}' '[0-9a-z]{64}' p2 p2h-digest-kein-hex
mut p2h-verfahren "$D_" '^sha256:' '^sha[0-9]+:' p2 p2h-digest-verfahren
mut p2i "$D_" %P2i p2 p2i-referenz-image
mut p2-vor-docker "$D_" 'if [ "$fehler" -ne 0 ]; then exit 2; fi' 'if [ "$fehler" -ne 0 ]; then :; fi' p2 "${p2namen[@]}"

# --- Punkt 3: Teilmenge ---
for t in "${tests[@]}"; do
  mut "p3-element-${t//\//-}" "$F_" " $t" "" "p3-element-$t" "p3-element-$t"
done
mut p3-skip "$R_" 'if grep -Eq "^ *--- SKIP: $name \(" "$aus"; then # P3-skip' 'if false; then # P3-skip' "p3-skip p3-ohne-treffer" p3-skip
mut p3-treffer "$R_" "if ! grep -Eq '^ *--- PASS: ' \"\$aus\"; then # P3-treffer" 'if false; then # P3-treffer' "p3-skip p3-ohne-treffer" p3-ohne-treffer
mut p3-fehlt "$R_" 'rot=1 # P3-fehlt' 'rot=0 # P3-fehlt' "p3-skip p3-element p3-fehlt" p3-skip p3-fehlt-phase2 "${tests[@]/#/p3-element-}"
mut p3-auswahl-top "$R_" 'top+=("${n%%/*}")' ':' "p3-auswahl p3-element" p3-auswahl "${tests[@]/#/p3-element-}"
mut p3-phase-fest "$R_" 'local aus="$1" ph="$2" name rot=0' 'local aus="$1" ph=1 name rot=0' p3 p3-fehlt-phase2 p3-element-TestE2EOhnePostgresReplay p3-element-TestE2EOhnePostgresExtendedReplay
mut p3-phase-fest-16 "$R_" 'local aus="$1" ph="$2" name rot=0' 'local aus="$1" ph="$2" name rot=0; if [ "$version" = 16 ]; then ph=1; fi' p3 p3-fehlt-phase2
mut p3-auswahl-sub "$R_" 'sub+=("${n#*/}")' ':' "p3-auswahl p3-element" p3-auswahl p3-element-TestE2ERecordAnmeldung/scram-sha-256

# --- Punkt 4: Server ---
mut p4-version "$R_" 'if [ "${gemeldet%%[!0-9]*}" != "$PGR_HAUPTVERSION" ]; then # P4-version' 'if false; then # P4-version' p4-version p4-version-abweichend p4-version-unlesbar p4-version-ohne-ziffer
mut p4-unlesbar "$R_" '!= "$PGR_HAUPTVERSION" ]; then # P4-version' '!= "$PGR_HAUPTVERSION" ] && [ -n "$gemeldet" ]; then # P4-version' p4-version p4-version-unlesbar
mut p4-ziffer "$R_" '${gemeldet%%[!0-9]*}' '${gemeldet//[!0-9]/}' p4-version p4-version-ohne-ziffer
mut p4-60a "$R_" 'seq 1 120' 'seq 1 60' "p4-60 p4-nicht-bereit" p4-60-sekunden
mut p4-60b "$R_" 'sleep 0.5' 'sleep 1' "p4-60 p4-nicht-bereit" p4-60-sekunden
mut p4-bereit "$R_" 'exit 1 # E-bereit' ': # E-bereit' "p4-60 p4-nicht-bereit" p4-nicht-bereit p4-60-sekunden
mut p4-start "$R_" '"$postgres_image" >/dev/null' '"$postgres_image" >/dev/null || true' p4-start p4-start
# Stufe „Aufbau“: je Form eine Mutation, die die Stufe auf Phase 1 fallen lässt und nur ihren Fall rot färbt.
mut p4-aufbau-bereit "$R_" 'exit 1 # E-bereit' 'phase 1; exit 1 # E-bereit' "p4-60 p4-nicht-bereit" p4-nicht-bereit p4-60-sekunden
mut p4-aufbau-bereit-16 "$R_" 'exit 1 # E-bereit' '[ "$version" = 16 ] && phase 1; exit 1 # E-bereit' "p4-60 p4-nicht-bereit" p4-60-sekunden
mut p4-aufbau-start "$R_" '"$postgres_image" >/dev/null' '"$postgres_image" >/dev/null || { phase 1; exit 1; }' p4-start p4-start
mut p4-aufbau-abweichend "$R_" "server_version '\$gemeldet'\" >&2" "server_version '\$gemeldet'\" >&2; case \"\$gemeldet\" in [0-9]*) phase 1 ;; esac" p4-version p4-version-abweichend
mut p4-aufbau-unlesbar "$R_" "server_version '\$gemeldet'\" >&2" "server_version '\$gemeldet'\" >&2; case \"\$gemeldet\" in '') phase 1 ;; esac" p4-version p4-version-unlesbar
mut p4-aufbau-ohne-ziffer "$R_" "server_version '\$gemeldet'\" >&2" "server_version '\$gemeldet'\" >&2; case \"\$gemeldet\" in [!0-9]?*) phase 1 ;; esac" p4-version p4-version-ohne-ziffer

# --- Punkt 5: Anmeldung ---
mut p5-benutzer "$R_" 'CREATE ROLE play_scram LOGIN' 'CREATE ROLE play_scram NOLOGIN' p5 p5-benutzer
mut p5-methode "$R_" 'play_scram all scram-sha-256' 'play_scram all md5' p5 p5-methode p5-hba-vor-allgemeiner
mut p5-passwort-verschluesselung "$R_" "SET password_encryption = 'scram-sha-256';" "SET password_encryption = 'md5';" p5 p5-methode p5-kein-md5-klartext
mut p5-hba-nach-allgemeiner "$R_" '{ printf "%s\n" '"'"'"$hba"'"'"'; cat "$hba"; }' '{ cat "$hba"; printf "%s\n" '"'"'"$hba"'"'"'; }' p5 p5-hba-vor-allgemeiner
mut p5-rollen "$R_" '!= scram ]; then # P5-rollen' '!= nie ]; then # P5-rollen' p5 p5-kein-md5-klartext
mut p5-einrichtung "$R_" 'SELECT pg_reload_conf();"' 'SELECT pg_reload_conf();" || true' p5 p5-einrichtung-scheitert
mut p5-aufbau-einrichtung "$R_" 'SELECT pg_reload_conf();"' 'SELECT pg_reload_conf();" || { phase 1; exit 1; }' p5 p5-einrichtung-scheitert

# --- Punkt 6: Ablauf und Meldung ---
mut p6-referenz "$D_" %P6-referenz "p6-referenz-nicht p6-alle-versionen" p6-referenz-nicht
mut p6-reihenfolge "$D_" 'for e in "${eintraege[@]}"; do
  haupt="${e%%=*}"
  [ "$haupt" != "$ref" ]' 'for e in $(printf "%s\n" "${eintraege[@]}" | tac); do
  haupt="${e%%=*}"
  [ "$haupt" != "$ref" ]' "p6-reihenfolge p6-zusammenfassung p7-vor-naechster" p6-reihenfolge p6-zusammenfassung p7-vor-naechster
mut p6-nacheinander "$R_" '"$docker" rm -f "$tests" "$pg" >/dev/null 2>&1 || true' '"$docker" rm -f "$tests" >/dev/null 2>&1 || true' "p6-nacheinander p7-container" p6-nacheinander p7-container
mut p6-phase2 "$R_" '  exit "$status"' '  :' "p6-meldung-phase1 p6-keine-phase2 p6-exit p6-zusammenfassung" p6-meldung-phase1 p6-keine-phase2 p6-zusammenfassung
mut p6-weiter "$D_" 'rot+=("$haupt ($ort)") # P6-phase' 'rot+=("$haupt ($ort)"); break # P6-phase' "p6-alle-versionen p6-keine-phase2 p6-zusammenfassung" p6-alle-versionen p6-keine-phase2 p6-zusammenfassung
mut p6-phase-fehlt "$D_" 'rot+=("$haupt ($ort)")' 'rot+=("$haupt")' "p6-meldung p6-zusammenfassung p6-exit" p6-meldung-phase1 p6-meldung-phase2 p6-meldung-aufbau p6-zusammenfassung
mut p6-ort1 "$D_" '1) ort="Phase 1"' '1) ort="Phase 2"' "p6-meldung p6-zusammenfassung p6-exit" p6-meldung-phase1 p6-zusammenfassung
mut p6-ort2 "$D_" '2) ort="Phase 2"' '2) ort="Phase 1"' "p6-meldung p6-zusammenfassung p6-exit" p6-meldung-phase2 p6-zusammenfassung
mut p6-ort-aufbau "$D_" '*) ort="Aufbau"' '*) ort="Phase 1"' "p6-meldung p6-zusammenfassung p6-exit" p6-meldung-aufbau
mut p6-exit "$D_" 'exit 1 # P6-exit' 'exit 0 # P6-exit' "p6-meldung p6-zusammenfassung p6-exit" p6-exit
mut p6-bauen "$D_" 'PGR_BAUEN=0' 'PGR_BAUEN=1' p6-testimage-einmal p6-testimage-einmal
mut p6-build-rot "$D_" 'if ! wait "$laeuft"; then' 'if false; then' p6-build-rot p6-build-rot
mut p6-build-rot-aufruf "$D_" 'nicht zu bauen" >&2' 'nicht zu bauen" >&2; "$docker" network create --internal x' p6-build-rot p6-build-rot
mut p6-build-rot-exit "$D_" 'nicht zu bauen" >&2
  exit 1' 'nicht zu bauen" >&2
  exit 0' p6-build-rot p6-build-rot

# --- Punkt 7: Aufräumen ---
mut p7-container "$R_" '"$docker" rm -f "$tests" "$pg" >/dev/null 2>&1 || true' ':' "p7-container p7-netz p7-volume p7-ende p7-vor-naechster p7-sig" p7-container p7-ende-rot p7-ende-bereit p7-ende-gruen p7-vor-naechster p7-sigINT-gruppe p7-sigTERM-gruppe p7-sigINT-skript p7-sigTERM-skript
mut p7-netz "$R_" '"$docker" network rm "$netz" >/dev/null 2>&1 || true' ':' "p7-container p7-netz p7-volume p7-ende p7-vor-naechster p7-sig" p7-netz p7-ende-rot p7-ende-bereit p7-ende-pull p7-ende-gruen p7-vor-naechster p7-sigINT-gruppe p7-sigTERM-gruppe p7-sigINT-skript p7-sigTERM-skript
mut p7-volume "$R_" '"$docker" volume rm "$daten" >/dev/null 2>&1 || true' ':' "p7-container p7-netz p7-volume p7-ende p7-vor-naechster p7-sig" p7-volume p7-ende-rot p7-ende-bereit p7-ende-pull p7-ende-gruen p7-vor-naechster p7-sigINT-gruppe p7-sigTERM-gruppe p7-sigINT-skript p7-sigTERM-skript
mut p7-ende-rot "$R_" '  exit "$status"' '  kill -KILL $$' p7-ende p7-ende-rot
mut p7-ende-bereit "$R_" 'exit 1 # E-bereit' 'kill -KILL $$ # E-bereit' p7-ende p7-ende-bereit
mut p7-ende-pull "$R_" '"$postgres_image" >/dev/null' '"$postgres_image" >/dev/null || kill -KILL $$' p7-ende p7-ende-pull
mut p7-ende-gruen "$R_" 'echo "run-integration-tests: gruen"' 'kill -KILL $$' "p7-ende p7-container p7-netz p7-volume p7-vor-naechster" p7-ende-gruen p7-container p7-netz p7-volume p7-vor-naechster
mut p7-int-treiber "$D_" "trap 'beenden 130' INT # S-INT" "trap 'kill -KILL \$\$' INT # S-INT" p7-sig p7-sigINT-gruppe p7-sigINT-skript
mut p7-term-treiber "$D_" "trap 'beenden 143' TERM # S-TERM" "trap 'kill -KILL \$\$' TERM # S-TERM" p7-sig p7-sigTERM-gruppe p7-sigTERM-skript
mut p7-term-runner "$R_" "trap 'trap \"\" TERM; exit 143' TERM # S-TERM" "trap 'kill -KILL \$\$' TERM # S-TERM" p7-sig p7-sigINT-gruppe p7-sigINT-skript p7-sigTERM-gruppe p7-sigTERM-skript

# --- Punkt 8: der Runner ohne Auswahl ---
mut p8-phase1 "$R_" 'wahl_1=()' "wahl_1=(-test.run '^TestE2EOhnePostgres')" p8 p8-phase1
mut p8-phase2 "$R_" "wahl_2=(-test.run '^TestE2EOhnePostgres')" "wahl_2=(-test.run '^TestE2E')" p8 p8-phase2
mut p8-rollen "$R_" '!= scram ]; then # P5-rollen' '!= scram ] && false; then # P5-rollen' p8 p8-rollen
mut p8-pruefung "$R_" 'if [ "${#tests_1[@]}" -gt 0 ]; then erwartet' 'if true; then erwartet' p8 p8-ohne-pruefung p8-phase2

# --- Punkt 9: Aufrufarten des Treibers ---
mut p9-make "$F_" '@POSTGRES_VERSIONEN=' '@$(MAKE) --version >/dev/null; POSTGRES_VERSIONEN=' "p9-n p9-q p9-t p9-makeflags-n p9-i" p9-n p9-q p9-t p9-makeflags-n p9-n-und-i
mut p9-plus "$F_" $'\t@POSTGRES_VERSIONEN=' $'\t+@POSTGRES_VERSIONEN=' "p9-n p9-q p9-t p9-makeflags-n p9-i" p9-n p9-q p9-t p9-makeflags-n p9-n-und-i
mut p9-stand "$F_" 'bash tools/test/run-postgres-versionen.sh' 'bash tools/test/run-postgres-versionen.sh && touch .harness/state/stand' p9-kein-stand p9-kein-stand

# --- Punkt 12: Anschluss an GATE_CHECKS ---
mut p12-ziel "$F_" 'GATE_CHECKS += test-postgres-versionen ' 'GATE_CHECKS += ' p12 p12-gate-checks-ziel
mut p12-gegenprobe "$F_" 'GATE_CHECKS += test-postgres-versionen test-postgres-versionen-gegenprobe' 'GATE_CHECKS += test-postgres-versionen' p12 p12-gate-checks-gegenprobe

if [ "$fehler" -ne 0 ]; then
  exit 1
fi
echo "postgres-versionen-gegenprobe: gruen — Punkte 1 bis 9 und 12: unveränderte Kopie grün, je Mutation genau die benannten Fälle rot"
