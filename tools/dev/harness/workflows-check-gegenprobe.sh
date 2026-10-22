#!/usr/bin/env bash
# workflows-check-gegenprobe — prüft tools/harness/workflows-check.sh und harness/mk/workflows.mk gegen die Punkte 1 bis 12 von SPEC-063
# (harness/vertraege/SPEC-063-workflows-check.md). Je Fall ein Baum in einem Temp-Verzeichnis, das Skript mit einer Attrappe von `docker`
# (über DOCKER), die ihre Argumente aufzeichnet und Ausgabe, Fehlerausgabe und Exit-Code vorgibt. Die unveränderte Kopie muss jeden Fall
# bestehen; je Mutation (Textersatz in einer Kopie von Skript oder Fragment, bei zwei Mutationen der echten Fälle in der Umgebung) müssen genau
# die benannten Fälle rot werden. Der Name jedes Falls trägt die Nummer des Punktes. Drei echte Fälle (real-*) laufen gegen das Image aus
# WORKFLOWS_ACTIONLINT und das Image aus DCHECK_REF (beide übergibt `make`; der erste Lauf zieht sie); sie laufen in der unveränderten Kopie und bei
# den Mutationen, deren rote Fälle einen real-* nennen.
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Mutationen, deren Name passt (die unveränderte Kopie nur, wenn das Muster `ok` trifft);
# GP_ENTDECKE=1 nennt je Mutation die roten Fälle; GP_PARALLEL=<n> (Vorgabe 4) setzt die Zahl der gleichzeitigen Läufe;
# GP_FALL=<name> gibt für die unveränderte Kopie das Protokoll dieses Falls aus.
# Braucht bash ≥ 4, make, mktemp, awk, find, sort, cksum und Docker für die echten Fälle.
# GRENZE: Als root bleibt jede Datei und jedes Verzeichnis lesbar; dann sind p9-nicht-lesbar, p9-verzeichnis-ohne-r und p9-verzeichnis-ohne-x nicht herstellbar,
# ihre Mutationen entfallen, und das Ende der Gegenprobe nennt sie auf stderr.
# Abdeckung der Punkte 11 (Grenzen) und der Lauffähigkeit auf GitHub hat die Gegenprobe nicht: es sind akzeptierte Negative des Vertrags.
set -uo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="$(cd "$hier/../../.." && pwd)"
skript0="$repo/tools/harness/workflows-check.sh"
fragment0="$repo/harness/mk/workflows.mk"
DOCKER="${DOCKER:-docker}"
arbeit="$(mktemp -d)"
trap 'chmod -R u+rwX "${arbeit:?}" 2> /dev/null; rm -rf "${arbeit:?}"' EXIT
: > "$arbeit/fehler"
: > "$arbeit/gefahren"
melde() { echo "workflows-check-gegenprobe: ROT — $*" >> "$arbeit/fehler"; }

for t in make awk find sort cksum mktemp; do
  command -v "$t" > /dev/null 2>&1 || { echo "workflows-check-gegenprobe: $t fehlt" >&2; exit 2; }
done
if [ -z "${WORKFLOWS_ACTIONLINT:-}" ] || [ -z "${DCHECK_REF:-}" ]; then
  echo "workflows-check-gegenprobe: WORKFLOWS_ACTIONLINT und DCHECK_REF fehlen (aufrufen über make workflows-check-gegenprobe)" >&2
  exit 2
fi

IMG='img/test@sha256:0000000000000000000000000000000000000000000000000000000000000000'

# Die Attrappe von docker: zeichnet jeden Aufruf in $AT_LOG auf (Zeile CALL, dann ein Argument je Zeile) und gibt AT_OUT, AT_ERR und AT_EXIT vor.
cat > "$arbeit/attrappe" << 'EOT'
#!/usr/bin/env bash
if [ -n "${AT_LOG-}" ]; then { echo CALL; printf '%s\n' "$@"; } >> "$AT_LOG"; fi
if [ -n "${AT_OUT-}" ]; then printf '%s\n' "$AT_OUT"; fi
if [ -n "${AT_ERR-}" ]; then printf '%s\n' "$AT_ERR" >&2; fi
exit "${AT_EXIT-0}"
EOT
chmod +x "$arbeit/attrappe"
attrappe="$arbeit/attrappe"

FAELLE=(
  p1-yml p1-yaml p1-andere-endung p1-unterverzeichnis p1-symlink p1-symlink-andere-endung
  p2-kein-verzeichnis p2-leer p2-nur-andere p2-ohne-4 p2-ohne-8 p2-nur-symlink
  p3-e1 p3-e2 p3-e3 p3-e4 p3-e5 p3-e6 p3-e7 p3-e8 p3-einmal
  p3-exit0 p3-exit1 p3-exit1-leer p3-exit2 p3-exit3 p3-exit3-rechte p3-exit125 p3-kein-docker p3-digest-form p3-image-leer
  p4-yaml p4-yml
  p5-tab p5-einzug p5-einzug-spaeter p5-jobs-flow p5-schritt-ist-kein-job p5-form-ohne-regeln
  p6-fehlt p6-leer-ok p6-leer-kommentar-ok p6-read-all p6-write-all p6-block p6-flow-eintrag
  p6-job-block-ok p6-job-leer-ok p6-job-ohne-ok p6-job-write-all p6-job-read-all p6-job-flow-eintrag p6-job-wert
  p7-fehlt p7-1-ok p7-60-ok p7-61 p7-0 p7-ausdruck p7-nicht-ganzzahl p7-wort p7-fuehrende-null p7-vorzeichen p7-dezimal p7-uses-job-ok
  p8-ok p8-ok-doppelt p8-ok-einfach p8-modules p8-block p8-dir p8-dir-schraeg p8-dir-davor p8-dir-danach p8-modules-praefix p8-modules-suffix p8-ohne-datei p8-pins-nicht-hier
  p9-reihenfolge p9-alle-befunde p9-umgebung-vor-befund p9-gruen p9-nicht-lesbar p9-verzeichnis-ohne-r p9-verzeichnis-ohne-x
  p10-n p10-q p10-t p10-i p10-makeflags p10-rezept p10-unterverzeichnis p10-schreibt-nichts p10-kette-gate p10-kette-gegenprobe
  real-gruen real-ausdruck real-pins
)

# --- Bausteine der Fälle ---------------------------------------------------------------------------------------------------------------------------
# neu — legt einen Baum mit gebundener Pin-Prüfung an; w ist seine Wurzel.
neu() {
  w="$(mktemp -d "$run/t.XXXXXX")"
  chmod 755 "$w"
  mkdir -p "$w/.github/workflows"
  gebunden
}
gebunden() { printf 'modules: [links, workflows]\nworkflows:\n  dir: .github/workflows\n' > "$w/.d-check.yml"; }
schreibe() { printf '%s' "$2" > "$w/.github/workflows/$1"; }
# bau <oben> <nach runs-on> — eine Workflow-Datei; oben und nach runs-on sind Zeilen mit Zeilenende oder leer.
bau() {
  printf 'name: x\non: push\n%sjobs:\n  a:\n    runs-on: ubuntu-latest\n%s    steps:\n      - run: echo hi\n' "$1" "$2"
}
T0=$'permissions: {}\n'
J0=$'    timeout-minutes: 10\n'
gut() { bau "$T0" "$J0"; }

# lauf [argument…] — führt das Skript aus (Voreinstellung: Wurzel w); setzt code, out, err.
lauf() {
  local args=("$@") img="$IMG"
  [ "$#" -gt 0 ] || args=("$w")
  if [ -n "${IMG_W+x}" ]; then img="$IMG_W"; fi
  rm -f "$w.aufrufe"
  ( cd "${CWD:-$w}" && env DOCKER="${DOCKER_W:-$attrappe}" WORKFLOWS_ACTIONLINT="$img" AT_LOG="$w.aufrufe" bash "${SKRIPT_W:-$skript}" "${args[@]}" ) > "$w.out" 2> "$w.err"
  code=$?
  out="$(cat "$w.out")"
  err="$(cat "$w.err")"
}

chk() {
  if ! "$@"; then
    bad=1
    { printf 'nicht erfüllt: %s\n' "$*"; printf 'Exit %s\nstdout:\n%s\nstderr:\n%s\n' "${code-?}" "${out-}" "${err-}"; } >> "$run/log.$fallname"
  fi
}
code_ist() { [ "$code" -eq "$1" ]; }
err_hat() { grep -qF -- "$1" <<< "$err"; }
err_hat_nicht() { ! grep -qF -- "$1" <<< "$err"; }
out_hat() { grep -qF -- "$1" <<< "$out"; }
err_leer() { [ -z "$err" ]; }
aufrufe_sind() { local n; n="$(grep -c '^CALL$' "$w.aufrufe" 2> /dev/null || true)"; [ "${n:-0}" -eq "$1" ]; }
# tags — je Zeile der Fehlerausgabe (ohne Zusammenfassung) die Regel; actionlint-Zeilen heißen al.
tags() {
  local z
  while IFS= read -r z; do
    if [[ $z =~ ^\.github/[^:]+:[0-9]+:\ ([a-z-]+):\  ]]; then echo "${BASH_REMATCH[1]}"
    elif [[ $z =~ ^workflows-check:\ [0-9]+\ Befund ]]; then :
    elif [[ $z =~ :[0-9]+:[0-9]+:\  ]]; then echo al
    elif [ -n "$z" ]; then echo "?"; fi
  done <<< "$err"
}
tags_sind() { [ "$(tags | tr '\n' ' ')" = "$* " ] || { [ "$#" -eq 0 ] && [ -z "$(tags)" ]; }; }
summe_stimmt() {
  local letzte n
  letzte="$(grep -v '^$' <<< "$err" | tail -n 1)"
  n="$(tags | grep -c .)"
  [ "$letzte" = "workflows-check: $n Befund(e)" ]
}
rot() { chk code_ist 1; chk tags_sind "$@"; chk summe_stimmt; }
gruen() { chk code_ist 0; chk err_leer; chk out_hat "kein Befund"; }
umgebung() { chk code_ist 2; }
# argumente — die Argumente des (ersten) Aufrufs der Attrappe, eines je Zeile, ohne die Zeile CALL.
argumente() { sed -n '2,$p' "$w.aufrufe" 2> /dev/null; }
hat_arg() { argumente | grep -qxF -- "$1"; }
hat_paar() { argumente | awk -v a="$1" -v b="$2" 'prev == a && $0 == b { f = 1 } { prev = $0 } END { exit !f }'; }

# --- Punkt 1: Prüfbereich ---------------------------------------------------------------------------------------------------------------------------
f_p1_yml() { neu; schreibe a.yml "$(bau "" "$J0")"; lauf; rot permissions-fehlt; chk err_hat ".github/workflows/a.yml:1: permissions-fehlt"; }
f_p1_yaml() { neu; schreibe a.yaml "$(bau "" "$J0")"; lauf; rot permissions-fehlt; chk err_hat ".github/workflows/a.yaml:1: permissions-fehlt"; }
f_p1_andere_endung() {
  neu; rm "$w/.d-check.yml"
  local b; b="$(bau "" "")"
  schreibe x.YML "$b"; schreibe x.yml.bak "$b"; schreibe x.txt "$b"
  lauf
  chk code_ist 0; chk err_leer
  chk out_hat "hinweis: .github/workflows/x.YML wird nicht geprüft (Endung)"
  chk out_hat "hinweis: .github/workflows/x.yml.bak wird nicht geprüft (Endung)"
  chk out_hat "hinweis: .github/workflows/x.txt wird nicht geprüft (Endung)"
  chk aufrufe_sind 0
}
f_p1_unterverzeichnis() {
  neu; schreibe a.yml "$(gut)"; mkdir "$w/.github/workflows/sub"; schreibe sub/b.yml "$(bau "" "")"
  lauf
  gruen; chk out_hat "hinweis: .github/workflows/sub/b.yml wird nicht geprüft (Unterverzeichnis)"
  chk out_hat "1 Datei(en) geprüft"
}
f_p1_symlink_andere_endung() { neu; schreibe a.yml "$(gut)"; ln -s a.yml "$w/.github/workflows/x.txt"; lauf; gruen; chk out_hat "hinweis: .github/workflows/x.txt wird nicht geprüft (Endung)"; }
f_p1_symlink() {
  neu; schreibe a.yml "$(gut)"; ln -s a.yml "$w/.github/workflows/l.yml"
  lauf
  rot symlink; chk err_hat ".github/workflows/l.yml:1: symlink:"
  chk argumente_ohne_l;
}
argumente_ohne_l() { ! argumente | grep -qF 'l.yml'; }

# --- Punkt 2: kein Prüfbereich ----------------------------------------------------------------------------------------------------------------------
KEINE="workflows-check: keine Workflow-Dateien unter .github/workflows"
f_p2_kein_verzeichnis() { neu; rm -rf "$w/.github" "$w/.d-check.yml"; lauf; chk code_ist 0; chk err_leer; chk out_hat "$KEINE"; chk aufrufe_sind 0; }
f_p2_leer() { neu; rm "$w/.d-check.yml"; lauf; chk code_ist 0; chk err_leer; chk out_hat "$KEINE"; chk aufrufe_sind 0; }
f_p2_nur_andere() {
  neu; rm "$w/.d-check.yml"; schreibe README.md "x"; schreibe x.txt "y"
  lauf; chk code_ist 0; chk err_leer; chk out_hat "$KEINE"; chk aufrufe_sind 0
}
f_p2_ohne_4() { neu; rm "$w/.d-check.yml"; printf 'x\n' > "$w/.github/actionlint.yaml"; lauf; chk code_ist 0; chk err_leer; chk out_hat "$KEINE"; }
f_p2_ohne_8() { neu; rm "$w/.d-check.yml"; schreibe y.txt "y"; lauf; chk code_ist 0; chk err_leer; chk out_hat "$KEINE"; }

f_p2_nur_symlink() {
  neu; rm "$w/.d-check.yml"; ln -s ziel.txt "$w/.github/workflows/l.yml"; printf 'x\n' > "$w/.github/workflows/ziel.txt"
  lauf
  chk code_ist 1; chk tags_sind symlink; chk summe_stimmt; chk err_hat "workflows-check: 1 Befund(e)"
  chk aufrufe_sind 0; chk test -z "$(out_hat "$KEINE" && echo ja)"
}

# --- Punkt 3: Aufruf von actionlint -----------------------------------------------------------------------------------------------------------------
zwei() { neu; schreibe b.yaml "$(gut)"; schreibe a.yml "$(gut)"; lauf; chk code_ist 0; }
f_p3_e1() { zwei; chk hat_arg --rm; }
f_p3_e2() { zwei; chk hat_paar --network none; }
f_p3_e3() { zwei; chk hat_paar -v "$(cd "$w" && pwd -P):/repo:ro"; }
f_p3_e4() { zwei; chk hat_arg "$IMG"; }
f_p3_e5() { zwei; chk hat_arg -oneline; }
f_p3_e6() { zwei; chk test "$(argumente | tail -n 2 | tr '\n' ' ')" = ".github/workflows/a.yml .github/workflows/b.yaml "; }
f_p3_e7() { zwei; chk test -z "$(argumente | grep -- '--name' || true)"; }
f_p3_e8() { zwei; chk test -z "$(argumente | grep -E -- '^-(ignore|shellcheck|pyflakes)' || true)"; }
f_p3_einmal() { zwei; chk aufrufe_sind 1; chk test "$(argumente | head -n 1)" = run; }
einfach() { neu; schreibe a.yml "$(gut)"; }
f_p3_exit0() { einfach; AT_EXIT=0 lauf; gruen; chk out_hat "1 Datei(en) geprüft, kein Befund"; }
f_p3_exit1() {
  einfach; AT_EXIT=1 AT_OUT=$'.github/workflows/a.yml:3:1: boom [x]\n.github/workflows/a.yml:4:1: bam [y]' lauf
  rot al al
  chk test "$err" = $'.github/workflows/a.yml:3:1: boom [x]\n.github/workflows/a.yml:4:1: bam [y]\nworkflows-check: 2 Befund(e)'
}
f_p3_exit1_leer() { einfach; AT_EXIT=1 lauf; chk code_ist 1; chk err_hat "ohne Zeile"; }
f_p3_exit2() {
  einfach; AT_EXIT=2 AT_ERR="boom permission denied" lauf
  umgebung; chk err_hat "workflows-check: actionlint nicht lauffähig (Exit 2)"; chk err_hat "boom permission denied"
  chk err_hat_nicht "Hinweis: der Container liest als Nicht-root"
}
f_p3_exit3() {
  einfach; AT_EXIT=3 AT_ERR='could not read "x.yml": open x.yml: no such file or directory' lauf
  umgebung; chk err_hat "actionlint nicht lauffähig (Exit 3)"; chk err_hat "no such file or directory"
  chk err_hat_nicht "Hinweis: der Container liest als Nicht-root"
}
f_p3_exit3_rechte() {
  einfach; AT_EXIT=3 AT_ERR='could not read ".github/workflows/a.yml": open .github/workflows/a.yml: permission denied' lauf
  umgebung; chk err_hat "actionlint nicht lauffähig (Exit 3)"
  chk err_hat "workflows-check: Hinweis: der Container liest als Nicht-root; Wurzel und .github/workflows brauchen Leserechte für alle"
}
f_p3_exit125() { einfach; AT_EXIT=125 AT_ERR="pull access denied" lauf; umgebung; chk err_hat "actionlint nicht lauffähig (Exit 125)"; chk err_hat "pull access denied"; }
f_p3_kein_docker() {
  einfach; DOCKER_W="$arbeit/gibt-es-nicht" lauf
  umgebung; chk err_hat "actionlint nicht lauffähig (Exit 127)"; chk out_hat_nicht_gruen
}
out_hat_nicht_gruen() { ! out_hat "kein Befund"; }
f_p3_digest_form() { chk grep -Eq '^WORKFLOWS_ACTIONLINT := rhysd/actionlint@sha256:[0-9a-f]{64}(#.*)?$' "$fragment"; }
f_p3_image_leer() { einfach; IMG_W="" lauf; umgebung; chk err_hat "WORKFLOWS_ACTIONLINT ist leer"; chk aufrufe_sind 0; }

# --- Punkt 4: Einstufung ----------------------------------------------------------------------------------------------------------------------------
f_p4_yaml() { einfach; printf 'x\n' > "$w/.github/actionlint.yaml"; lauf; rot konfiguration; chk err_hat ".github/actionlint.yaml:1: konfiguration:"; }
f_p4_yml() { einfach; printf 'x\n' > "$w/.github/actionlint.yml"; lauf; rot konfiguration; chk err_hat ".github/actionlint.yml:1: konfiguration:"; }

# --- Punkt 5: Lesegrammatik -------------------------------------------------------------------------------------------------------------------------
f_p5_tab() { neu; schreibe a.yml $'name: x\non: push\npermissions: {}\njobs:\n  a:\n\truns-on: ubuntu-latest\n'; lauf; rot form; chk err_hat ".github/workflows/a.yml:6: form:"; }
f_p5_einzug() { neu; schreibe a.yml $'name: x\non: push\npermissions: {}\njobs:\n    a:\n        runs-on: ubuntu-latest\n'; lauf; rot form; chk err_hat ".github/workflows/a.yml:5: form:"; }
f_p5_einzug_spaeter() {
  neu; schreibe a.yml "$(gut)"$'\n   b:\n     runs-on: ubuntu-latest\n'; lauf; rot form; chk err_hat ".github/workflows/a.yml:10: form:"
}
f_p5_jobs_flow() { neu; schreibe a.yml $'name: x\non: push\npermissions: {}\njobs: {a: {runs-on: ubuntu-latest}}\n'; lauf; rot form; chk err_hat ".github/workflows/a.yml:4: form:"; }
f_p5_schritt_ist_kein_job() {
  neu; schreibe a.yml $'name: x\non: push\npermissions: {}\njobs:\n  a:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo hi\n        timeout-minutes: 10\n'
  lauf; rot timeout-fehlt
}

f_p5_form_ohne_regeln() { neu; schreibe a.yml $'name: x\non: push\njobs:\n  a:\n    runs-on: ubuntu-latest\n\tx: 1\n'; lauf; rot form; }

# --- Punkt 6: Rechte --------------------------------------------------------------------------------------------------------------------------------
f_p6_fehlt() { neu; schreibe a.yml "$(bau "" "$J0")"; lauf; rot permissions-fehlt; }
f_p6_leer_ok() { neu; schreibe a.yml "$(gut)"; lauf; gruen; }
f_p6_leer_kommentar_ok() { neu; schreibe a.yml "$(bau $'# keine Rechte\npermissions: {}  # nichts\n' "$J0")"; lauf; gruen; }
f_p6_read_all() { neu; schreibe a.yml "$(bau $'permissions: read-all\n' "$J0")"; lauf; rot permissions-form; chk err_hat "gelesen: read-all"; }
f_p6_write_all() { neu; schreibe a.yml "$(bau $'permissions: write-all\n' "$J0")"; lauf; rot permissions-form; chk err_hat "gelesen: write-all"; }
f_p6_block() { neu; schreibe a.yml "$(bau $'permissions:\n  contents: read\n' "$J0")"; lauf; rot permissions-form; chk err_hat "gelesen: ein Block"; }
f_p6_flow_eintrag() { neu; schreibe a.yml "$(bau $'permissions: {contents: read}\n' "$J0")"; lauf; rot permissions-form; chk err_hat "gelesen: {contents: read}"; }
jobperm() { neu; schreibe a.yml "$(bau "$T0" "$J0$1")"; lauf; }
f_p6_job_block_ok() { jobperm $'    permissions:\n      contents: read\n      issues: write\n      id-token: none\n'; gruen; }
f_p6_job_leer_ok() { jobperm $'    permissions: {}\n'; gruen; }
f_p6_job_ohne_ok() { jobperm ''; gruen; }
f_p6_job_write_all() { jobperm $'    permissions: write-all\n'; rot job-permissions-form; chk err_hat "Job a: permissions als Skalar: write-all"; }
f_p6_job_read_all() { jobperm $'    permissions: read-all\n'; rot job-permissions-form; chk err_hat "als Skalar: read-all"; }
f_p6_job_flow_eintrag() { jobperm $'    permissions: {contents: read}\n'; rot job-permissions-form; chk err_hat "als Flow-Mapping mit Einträgen: {contents: read}"; }
f_p6_job_wert() { jobperm $'    permissions:\n      contents: admin\n'; rot job-permissions-wert; chk err_hat "contents: admin"; }

# --- Punkt 7: Zeitgrenzen ---------------------------------------------------------------------------------------------------------------------------
zeit() { neu; schreibe a.yml "$(bau "$T0" "$1")"; lauf; }
f_p7_fehlt() { zeit ''; rot timeout-fehlt; chk err_hat "Job a: timeout-minutes fehlt"; }
f_p7_1_ok() { zeit $'    timeout-minutes: 1\n'; gruen; }
f_p7_60_ok() { zeit $'    timeout-minutes: 60  # Obergrenze\n'; gruen; }
f_p7_61() { zeit $'    timeout-minutes: 61\n'; rot timeout-wert; chk err_hat "timeout-minutes 61"; }
f_p7_0() { zeit $'    timeout-minutes: 0\n'; rot timeout-wert; chk err_hat "timeout-minutes 0"; }
f_p7_ausdruck() { zeit $'    timeout-minutes: ${{ 10 }}\n'; rot timeout-wert; chk err_hat 'timeout-minutes ${{ 10 }}'; }
f_p7_nicht_ganzzahl() { zeit $'    timeout-minutes: 30.5\n'; rot timeout-wert; chk err_hat "timeout-minutes 30.5"; }
f_p7_wort() { zeit $'    timeout-minutes: bald\n'; rot timeout-wert; chk err_hat "timeout-minutes bald"; }
f_p7_fuehrende_null() { local v; for v in 010 0060 00; do zeit "    timeout-minutes: $v"$'\n'; rot timeout-wert; chk err_hat "timeout-minutes $v"; done; }
f_p7_vorzeichen() { zeit $'    timeout-minutes: +5\n'; rot timeout-wert; chk err_hat "timeout-minutes +5"; }
f_p7_dezimal() { zeit $'    timeout-minutes: 5.0\n'; rot timeout-wert; chk err_hat "timeout-minutes 5.0"; }
f_p7_uses_job_ok() {
  neu; schreibe a.yml "$(gut)"$'\n  c:\n    uses: ./.github/workflows/w.yml\n'; lauf; gruen
}

# --- Punkt 8: Bindung der Pin-Prüfung ---------------------------------------------------------------------------------------------------------------
bindung() { neu; printf '%s' "$1" > "$w/.d-check.yml"; schreibe a.yml "${2:-$(gut)}"; lauf; }
f_p8_ok() { bindung $'modules: [links, workflows]\nworkflows:\n  dir: .github/workflows\n'; gruen; }
f_p8_ok_doppelt() { bindung $'modules: [workflows]\nworkflows:\n  dir: ".github/workflows"\n'; gruen; }
f_p8_ok_einfach() { bindung $'modules: [workflows]\nworkflows:\n  dir: \'.github/workflows\'\n'; gruen; }
f_p8_modules() { bindung $'modules: [links]\nworkflows:\n  dir: .github/workflows\n'; rot pins-ungeprueft; chk err_hat "modules: führt workflows nicht"; chk err_hat ".github/workflows/a.yml:1: pins-ungeprueft:"; }
f_p8_block() { bindung $'modules: [workflows]\n'; rot pins-ungeprueft; chk err_hat "Zeile workflows: fehlt"; }
f_p8_dir() { bindung $'modules: [workflows]\nworkflows:\n  dir: docs\n'; rot pins-ungeprueft; chk err_hat "dir: .github/workflows unter workflows: fehlt"; }
f_p8_dir_schraeg() { bindung $'modules: [workflows]\nworkflows:\n  dir: .github/workflows/\n'; rot pins-ungeprueft; chk err_hat "dir: .github/workflows unter workflows: fehlt"; }
dbl() { bindung "$1"; rot pins-ungeprueft; chk err_hat "$2"; }
f_p8_dir_davor() { dbl $'modules: [workflows]\nandere:\n  dir: .github/workflows\nworkflows:\n  x: 1\n' "dir: .github/workflows unter workflows: fehlt"; }
f_p8_dir_danach() { dbl $'modules: [workflows]\nworkflows:\n  x: 1\nandere:\n  dir: .github/workflows\n' "dir: .github/workflows unter workflows: fehlt"; }
f_p8_modules_praefix() { dbl $'modules: [workflows-x]\nworkflows:\n  dir: .github/workflows\n' "modules: führt workflows nicht"; }
f_p8_modules_suffix() { dbl $'modules: [xworkflows]\nworkflows:\n  dir: .github/workflows\n' "modules: führt workflows nicht"; }
f_p8_ohne_datei() { neu; rm "$w/.d-check.yml"; schreibe a.yml "$(gut)"; lauf; rot pins-ungeprueft; chk err_hat ".d-check.yml fehlt"; }
f_p8_pins_nicht_hier() {
  neu; schreibe a.yml "$(bau "$T0" "$J0" | sed 's|echo hi|echo hi\n      - uses: actions/checkout@v4|')"
  lauf; gruen
}

# --- Punkt 9: Ausgabe, Exit-Codes, Reihenfolge der Gründe -------------------------------------------------------------------------------------------
f_p9_reihenfolge() {
  neu; rm "$w/.d-check.yml"
  schreibe a.yml "$(bau "$T0" "")"
  schreibe b.yml "$(bau "" "$J0")"
  AT_EXIT=1 AT_OUT='.github/workflows/a.yml:2:1: boom [x]' lauf
  rot al timeout-fehlt permissions-fehlt pins-ungeprueft
  chk test "$(grep -n 'timeout-fehlt' <<< "$err" | grep -c 'a.yml')" -eq 1
  chk test "$(grep -n 'permissions-fehlt' <<< "$err" | grep -c 'b.yml')" -eq 1
}
f_p9_alle_befunde() {
  neu; schreibe a.yml "$(bau "" $'    permissions: read-all\n')"; lauf
  rot permissions-fehlt timeout-fehlt job-permissions-form
  chk err_hat "workflows-check: 3 Befund(e)"
}
f_p9_umgebung_vor_befund() {
  neu; schreibe a.yml "$(bau "" "$J0")"; DOCKER_W="$arbeit/gibt-es-nicht" lauf
  umgebung; chk err_hat_nicht "permissions-fehlt"; chk err_hat_nicht "Befund(e)"
}
f_p9_gruen() { neu; schreibe a.yml "$(gut)"; schreibe b.yaml "$(gut)"; lauf; gruen; chk test "$out" = "workflows-check: 2 Datei(en) geprüft, kein Befund"; }
nur_nichtroot() { [ "$(id -u)" -ne 0 ] || { echo "$fallname" >> "$arbeit/uebersprungen"; return 1; }; }
verz() { nur_nichtroot || return 0; neu; schreibe a.yml "$(gut)"; chmod "$1" "$w/.github/workflows"; lauf; chmod 755 "$w/.github/workflows"; umgebung; chk err_hat ".github/workflows nicht lesbar"; }
f_p9_verzeichnis_ohne_r() { verz 311; }
f_p9_verzeichnis_ohne_x() { verz 444; }
f_p9_nicht_lesbar() {
  nur_nichtroot || return 0
  neu; schreibe a.yml "$(gut)"; chmod 000 "$w/.github/workflows/a.yml"; lauf
  chmod 644 "$w/.github/workflows/a.yml"
  umgebung; chk err_hat ".github/workflows/a.yml nicht lesbar"
}

# --- Punkt 10: Aufrufarten, Unterverzeichnis, Schreiben, Parallelität -------------------------------------------------------------------------------
# kopie <workflow-text> — Baum mit Makefile, Fragment, Skript und der Attrappe als DOCKER.
kopie() {
  neu
  mkdir -p "$w/harness/mk" "$w/tools/harness"
  cp "$fragment" "$w/harness/mk/workflows.mk"; cp "$skript" "$w/tools/harness/workflows-check.sh"
  printf 'DOCKER := %s\nDCHECK_REF := x\nGATE_CHECKS :=\nGEGENPROBEN_CONTAINER :=\ninclude harness/mk/workflows.mk\n' "$attrappe" > "$w/Makefile"
  schreibe a.yml "$1"
}
# mk <make-argument>… — make im Baum; setzt code, out, err.
mk() {
  rm -f "$w.aufrufe"
  ( cd "$w" && env -u MAKEFLAGS -u MFLAGS -u MAKELEVEL ${MKENV-} AT_LOG="$w.aufrufe" make "$@" ) > "$w.out" 2> "$w.err"
  code=$?; out="$(cat "$w.out")"; err="$(cat "$w.err")"
}
f_p10_n() { kopie "$(bau "" "$J0")"; mk -n workflows-check; chk code_ist 0; chk out_hat "bash tools/harness/workflows-check.sh"; chk aufrufe_sind 0; }
f_p10_q() { kopie "$(bau "" "$J0")"; mk -q workflows-check; chk code_ist 1; chk aufrufe_sind 0; }
f_p10_t() { kopie "$(bau "" "$J0")"; mk -t workflows-check; chk code_ist 0; chk aufrufe_sind 0; }
f_p10_i() { kopie "$(gut)"; AT_EXIT=1 AT_OUT='.github/workflows/a.yml:2:1: boom [x]' mk -i workflows-check; chk code_ist 0; chk aufrufe_sind 1; chk err_hat "boom [x]"; }
f_p10_makeflags() { kopie "$(bau "" "$J0")"; MKENV="MAKEFLAGS=n" mk workflows-check; chk code_ist 0; chk aufrufe_sind 0; }
f_p10_rezept() { kopie "$(gut)"; mk workflows-check; chk code_ist 0; chk aufrufe_sind 1; chk out_hat "kein Befund"; }
f_p10_unterverzeichnis() {
  kopie "$(gut)"; mkdir "$w/sub"
  AT_EXIT=1 AT_OUT='.github/workflows/a.yml:2:1: boom [x]' SKRIPT_W="$w/tools/harness/workflows-check.sh" lauf "$w"; local c1="$code" e1="$err"
  AT_EXIT=1 AT_OUT='.github/workflows/a.yml:2:1: boom [x]' CWD="$w/sub" SKRIPT_W="$w/tools/harness/workflows-check.sh" lauf ""; local c2="$code" e2="$err"
  code=$c2
  chk test "$c1" -eq 1; chk test "$c2" -eq 1
  chk test "$e1" = "$e2"
  chk err_hat "boom [x]"
}
baumsumme() { ( cd "$w" && { find . | sort; find . -type f -exec cksum {} + | sort; } | cksum ); }
f_p10_schreibt_nichts() {
  kopie "$(gut)"; local v1 n1 v2 n2
  v1="$(baumsumme)"; AT_EXIT=1 AT_OUT='.github/workflows/a.yml:2:1: boom [x]' SKRIPT_W="$w/tools/harness/workflows-check.sh" lauf "$w"; n1="$(baumsumme)"
  v2="$(baumsumme)"; SKRIPT_W="$w/tools/harness/workflows-check.sh" lauf "$w"; n2="$(baumsumme)"
  chk test "$v1" = "$n1"; chk test "$v2" = "$n2"
}
liste() { ( cd "$w" && make -s -f Makefile --eval "gp-lesen: ; @printf '%s\\n' \$($1)" gp-lesen 2> /dev/null ); }
f_p10_kette_gate() { kopie "$(gut)"; chk grep -qxF workflows-check <<< "$(liste GATE_CHECKS)"; }
f_p10_kette_gegenprobe() { kopie "$(gut)"; chk grep -qxF workflows-check-gegenprobe <<< "$(liste GEGENPROBEN_CONTAINER)"; }

# --- Punkt 12: echte Fälle --------------------------------------------------------------------------------------------------------------------------
echt() { DOCKER_W="${REAL_DOCKER:-$DOCKER}" IMG_W="$WORKFLOWS_ACTIONLINT" lauf; }
f_real_gruen() { neu; schreibe a.yml "$(gut)"; chmod -R a+rX "$w"; echt; gruen; }
f_real_ausdruck() {
  neu; schreibe a.yml "$(bau "$T0" "$J0" | sed 's|echo hi|echo ${{ github.evnt }}|')"; chmod -R a+rX "$w"; echt
  chk code_ist 1; chk err_hat ".github/workflows/a.yml:"; chk err_hat "evnt"
}
f_real_pins() {
  local v mod="${REAL_MODULES:-[workflows]}" sha=11bd71901bbe5b1630ceea73d27597364c9af683
  neu
  printf 'modules: %s\nworkflows:\n  dir: .github/workflows\n' "$mod" > "$w/.d-check.yml"
  v="$(bau "$T0" "$J0" | sed "s|echo hi|echo hi\n      - uses: actions/checkout@v4|")"; schreibe a.yml "$v"; chmod -R a+rX "$w"
  out="$("$DOCKER" run --rm --network none -v "$w:/repo:ro" "$DCHECK_REF" 2> /dev/null)"; code=$?; err=""
  chk code_ist 1; chk out_hat "uses-pin-missing"
  v="$(bau "$T0" "$J0" | sed "s|echo hi|echo hi\n      - uses: actions/checkout@$sha|")"; schreibe a.yml "$v"; chmod -R a+rX "$w"
  out="$("$DOCKER" run --rm --network none -v "$w:/repo:ro" "$DCHECK_REF" 2> /dev/null)"; code=$?
  chk code_ist 1; chk out_hat "uses-pin-untagged"
  v="$(bau "$T0" "$J0" | sed "s|echo hi|echo hi\n      - uses: actions/checkout@$sha # v4.2.2|")"; schreibe a.yml "$v"; chmod -R a+rX "$w"
  out="$("$DOCKER" run --rm --network none -v "$w:/repo:ro" "$DCHECK_REF" 2> /dev/null)"; code=$?
  chk code_ist 0
}

# --- Läufe ------------------------------------------------------------------------------------------------------------------------------------------
# faelle <mit echten Fällen 0|1> — schreibt die Namen der roten Fälle, einen je Zeile.
faelle() {
  local n
  for n in "${FAELLE[@]}"; do
    case "$n" in real-*) [ "$1" = 1 ] || continue ;; esac
    fallname="$n"
    ( bad=0; "f_${n//-/_}"; exit "$bad" ) > /dev/null 2>> "$run/stderr" || echo "$n"
  done
}
# lauf_mit <name> <mit echten> <skript-text-datei> — wählt Skript und Fragment der Kopie und fährt die Fälle.
protokoll_zeigen() { [ -z "${GP_FALL-}" ] || [ ! -f "$run/log.$GP_FALL" ] || cat "$run/log.$GP_FALL" >&2; }

maxjobs="${GP_PARALLEL:-4}"
wartejob() { while [ "$(jobs -rp | wc -l)" -ge "$maxjobs" ]; do wait -n || true; done; }

# mut <name> <rote Fälle> <ziel: skript|fragment|env> <alt> <neu> [<alt> <neu>…] — Textersatz in der Kopie (bei env: Variable und Wert);
# rot sind genau die benannten Fälle.
mut() {
  local name="$1" rote="$2" ziel="$3" inhalt neu ist erwartet mit=0 alt rest muss nicht okb x
  shift 3
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  [[ " $rote " == *" real-"* ]] && mit=1
  if [ "$(id -u)" -eq 0 ]; then
    for x in p9-nicht-lesbar p9-verzeichnis-ohne-r p9-verzeichnis-ohne-x; do rote="${rote//$x/}"; done
    [ -n "${rote// /}" ] || { echo "$name" >> "$arbeit/uebersprungen"; return 0; }
  fi
  wartejob
  (
    run="$arbeit/$name"; mkdir -p "$run/tools/harness" "$run/harness/mk"
    cp "$skript0" "$run/tools/harness/workflows-check.sh"; cp "$fragment0" "$run/harness/mk/workflows.mk"
    skript="$run/tools/harness/workflows-check.sh"; fragment="$run/harness/mk/workflows.mk"
    case "$ziel" in
      skript) datei="$skript" ;;
      fragment) datei="$fragment" ;;
      env) datei="" ;;
    esac
    if [ "$ziel" = env ]; then
      while [ "$#" -ge 2 ]; do export "$1=$2"; shift 2; done
    else
      inhalt="$(cat "$datei"; printf x)"; inhalt="${inhalt%x}"; neu="$inhalt"
      while [ "$#" -ge 2 ]; do
        alt="$1"
        if [[ $neu != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
        rest="${neu#*"$alt"}"
        if [[ $rest == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
        neu="${neu/"$alt"/"$2"}"
        shift 2
      done
      if [[ $neu == "$inhalt" ]]; then melde "Mutation '$name': Ersatz ändert nichts"; exit 0; fi
      printf '%s' "$neu" > "$datei"
    fi
    ist="$(faelle "$mit" | LC_ALL=C sort -u | tr '\n' ' ')"
    if [[ $rote == MUSS:* ]]; then
      # breite Mutation: die genannten Fälle müssen rot sein, die hinter NICHT: dürfen es nicht sein
      muss="${rote#MUSS:}"; muss="${muss%% NICHT:*}"; nicht=""
      [[ $rote == *" NICHT:"* ]] && nicht="${rote#* NICHT:}"
      okb=1
      for x in $muss; do [[ " $ist " == *" $x "* ]] || okb=0; done
      for x in $nicht; do [[ " $ist " != *" $x "* ]] || okb=0; done
      if [ "$okb" -eq 1 ]; then rote="$ist"; fi
    fi
    erwartet="$(printf '%s\n' $rote | LC_ALL=C sort -u | tr '\n' ' ')"
    if [ -n "${GP_ENTDECKE-}" ]; then echo "$name => ${ist:-nichts}" >&2; fi
    echo "$name" >> "$arbeit/gefahren"
    [ "$ist" = "$erwartet" ] || melde "Mutation '$name': erwartet rot: ${rote:-nichts}; rot sind: ${ist:-nichts}"
  ) &
}

if [[ ok =~ ${GP_MUT:-.} ]]; then
  (
    run="$arbeit/ok"; mkdir -p "$run"; skript="$skript0"; fragment="$fragment0"
    ist="$(faelle 1 | LC_ALL=C sort -u | tr '\n' ' ')"
    echo ok >> "$arbeit/gefahren"
    [ -z "$ist" ] || melde "unveränderte Kopie ist nicht grün: $ist"
    protokoll_zeigen
    if [ -n "${GP_FALL-}" ]; then sed 's/^/stderr: /' "$run/stderr" >&2; fi
  ) &
fi

# --- Mutationen -------------------------------------------------------------------------------------------------------------------------------------
# Punkt 1
mut m-p1-nur-yml "p1-yaml p3-e6 p9-gruen" skript '*.yml | *.yaml) endung=1' '*.yml) endung=1'
mut m-p1-nur-yaml "MUSS:p1-yml NICHT:p1-yaml p2-leer p2-nur-andere p3-digest-form" skript '*.yml | *.yaml) endung=1' '*.yaml) endung=1'
mut m-p1-teilstring "p1-andere-endung" skript '*.yml | *.yaml) endung=1' '*.yml* | *.yaml*) endung=1'
mut m-p1-rekursiv "p1-unterverzeichnis" skript 'find "$dir" -mindepth 1 -maxdepth 1 -print0' 'find "$dir" -mindepth 1 -print0'
mut m-p1-symlink-gefolgt "p1-symlink p2-nur-symlink" skript '  if [ -L "$p" ]; then' '  if false; then'
mut m-p1-reihenfolge "p3-e6 p9-reihenfolge" skript 'done < <(find "$dir" -mindepth 1 -maxdepth 1 -print0 | sort -z)' 'done < <(find "$dir" -mindepth 1 -maxdepth 1 -print0 | sort -rz)'
# Punkt 2
mut m-p2-f1-exit2 "p2-kein-verzeichnis" skript '[ -d "$dir" ] || keine' '[ -d "$dir" ] || exit 2'
mut m-p2-f2-still "p2-leer p2-ohne-4" skript 'if [ "$alle" -eq 0 ]; then keine; fi' 'if [ "$alle" -eq 0 ]; then exit 0; fi'
mut m-p2-f3-docker "p1-andere-endung p2-nur-andere p2-nur-symlink p2-ohne-8" skript 'if [ "${#geordnet[@]}" -eq 0 ]; then keine; fi' ':' \
  $'if [ "${#kand[@]}" -gt 0 ]; then\n  [ -n "${WORKFLOWS_ACTIONLINT:-}" ]' $'if true; then\n  [ -n "${WORKFLOWS_ACTIONLINT:-}" ]'
mut m-p2-f4-wie-f3 "p2-nur-symlink" skript 'if [ "${#geordnet[@]}" -eq 0 ]; then keine; fi' 'if [ "${#kand[@]}" -eq 0 ]; then keine; fi'
mut m-p2-f4-docker "p2-nur-symlink" skript $'if [ "${#kand[@]}" -gt 0 ]; then\n  [ -n "${WORKFLOWS_ACTIONLINT:-}" ]' $'if true; then\n  [ -n "${WORKFLOWS_ACTIONLINT:-}" ]'
mut m-p2-punkt4 "p2-ohne-4" skript 'keine() { echo' 'keine() { if [ -e "$wurzel/.github/actionlint.yaml" ]; then echo "workflows-check: 1 Befund(e)" >&2; exit 1; fi; echo'
mut m-p2-punkt8 "p1-andere-endung p2-kein-verzeichnis p2-leer p2-nur-andere p2-ohne-4 p2-ohne-8" skript 'keine() { echo' 'keine() { if [ ! -f "$wurzel/.d-check.yml" ]; then echo "workflows-check: 1 Befund(e)" >&2; exit 1; fi; echo'
# Punkt 3
mut m-p3-e1 "p3-e1" skript 'run --rm --network' 'run --network'
mut m-p3-e2 "p3-e2" skript '--network none -v' '-v'
mut m-p3-e3 "p3-e3" skript ':/repo:ro"' ':/repo"'
mut m-p3-e4 "p3-e4" skript '-w /repo "$WORKFLOWS_ACTIONLINT" -oneline' '-w /repo fest/image -oneline'
mut m-p3-e5 "p3-e5" skript '-oneline "${dateien[@]}"' '"${dateien[@]}"'
mut m-p3-e6 "p3-e6" skript '"${dateien[@]}" 2> "$errf"' '"${dateien[0]}" 2> "$errf"'
mut m-p3-e7 "p3-e7" skript 'run --rm --network' 'run --rm --name wfc --network'
mut m-p3-e8-ignore "p3-e8" skript '-oneline "${dateien' '-oneline -ignore . "${dateien'
mut m-p3-e8-shellcheck "p3-e8" skript '-oneline "${dateien' '-oneline -shellcheck= "${dateien'
mut m-p3-e8-pyflakes "p3-e8" skript '-oneline "${dateien' '-oneline -pyflakes= "${dateien'
mut m-p3-zweimal "p3-einmal p10-i p10-rezept" skript '  al="$("$DOCKER" run' $'  "$DOCKER" run --rm "$WORKFLOWS_ACTIONLINT" > /dev/null 2>&1\n  al="$("$DOCKER" run'
mut m-p3-exit1-als-0 "p3-exit1 p3-exit1-leer p9-reihenfolge p10-i p10-unterverzeichnis real-ausdruck" skript '    0) ;;' '    0 | 1) ;;'
mut m-p3-jeder-exit-als-1 "p3-exit2 p3-exit3 p3-exit3-rechte p3-exit125 p3-kein-docker p9-umgebung-vor-befund" skript '    1)' '    1 | 2 | 3 | 125 | 127)'
mut m-p3-docker-fehlt-als-0 "p3-kein-docker p9-umgebung-vor-befund" skript $'\n    *)\n' $'\n    127) ;;\n    *)\n'
mut m-p3-hinweis-fehlt "p3-exit3-rechte" skript 'if [ "$rc" -eq 3 ] && grep -q' 'if false && grep -q'
mut m-p3-hinweis-bei-jedem-exit3 "p3-exit3" skript 'if [ "$rc" -eq 3 ] && grep -q '"'permission denied'"' "$errf"; then' 'if [ "$rc" -eq 3 ]; then'
mut m-p3-hinweis-bei-jedem-exit "p3-exit2" skript 'if [ "$rc" -eq 3 ] && grep -q' 'if true && grep -q'
mut m-p3-digest-tag "p3-digest-form" fragment '@sha256:' ':'
mut m-p3-digest-kurz "p3-digest-form" fragment 'e9# 1.7.7' 'e# 1.7.7'
mut m-p3-digest-repository "p3-digest-form" fragment ':= rhysd/actionlint@' ':= ghcr.io/rhysd/actionlint@'
mut m-p3-image-leer "p3-image-leer" skript '  [ -n "${WORKFLOWS_ACTIONLINT:-}" ] || sperre "WORKFLOWS_ACTIONLINT ist leer"' '  :'
# Punkt 4
mut m-p1-symlink-andere-befund "p1-symlink-andere-endung" skript 'if [ "$endung" -eq 1 ]; then geordnet+=("s:$n"); else' 'if true; then geordnet+=("s:$n"); else'
mut m-p1-symlink-andere-still "p1-symlink-andere-endung" skript 'else hinweise+=("hinweis: $WF/$n wird nicht geprüft (Endung)"); fi' 'else :; fi'
mut m-p1-hinweis-endung-fehlt "p1-andere-endung" skript $'    hinweise+=("hinweis: $WF/$n wird nicht geprüft (Endung)")\n  fi' $'    :\n  fi'
mut m-p1-hinweis-grund-endung "p1-andere-endung" skript $'    hinweise+=("hinweis: $WF/$n wird nicht geprüft (Endung)")\n  fi' $'    hinweise+=("hinweis: $WF/$n wird nicht geprüft (Name)")\n  fi'
mut m-p1-hinweis-unterverz-fehlt "p1-unterverzeichnis" skript 'hinweise+=("hinweis: ${q#"$wurzel/"} wird nicht geprüft (Unterverzeichnis)")' ':'
mut m-p1-hinweis-grund-unterverz "p1-unterverzeichnis" skript 'wird nicht geprüft (Unterverzeichnis)")' 'wird nicht geprüft (Ordner)")'
mut m-p3-al-zeilen-veraendert "p3-exit1" skript 'bef_al+=("$z")' 'bef_al+=("al: $z")'
mut m-p3-stderr-anhang-fehlt "p3-exit2 p3-exit3 p3-exit125" skript '      cat "$errf" >&2' '      :'
mut m-p4-yaml-nicht "p4-yaml" skript 'for k in actionlint.yaml actionlint.yml; do' 'for k in actionlint.yml; do'
mut m-p4-yml-nicht "p4-yml" skript 'for k in actionlint.yaml actionlint.yml; do' 'for k in actionlint.yaml; do'
# Punkt 5
mut m-p5-tab "p5-tab p5-form-ohne-regeln" skript '{ form(NR, "Tabulator im Einzug"); next }' '{ next }'
mut m-p5-einzug "p5-einzug" skript 'if (ind != 2) form(NR, "Job-Schlüssel mit Einzug " ind " statt 2")' ''
mut m-p5-einzug-spaeter "p5-einzug-spaeter" skript '{ form(NR, "Job-Schlüssel mit Einzug " ind " statt 2"); next }' '{ next }'
mut m-p5-jobs-flow "p5-jobs-flow" skript 'if (val != "") form(NR, "jobs: mit Wert in der Zeile: " val)' ''
mut m-p5-regeln-laufen "p5-form-ohne-regeln" skript 'print fl[i] "\tform\t" ft[i]; exit }' 'print fl[i] "\tform\t" ft[i] }'
mut m-p5-jeder-einzug "p5-schritt-ist-kein-job p6-job-wert" skript '  if (ind == 4 && haskey) {' '  if (ind >= 4 && haskey) {'
# Punkt 6
mut m-p6-fehlt-angenommen "p6-fehlt p1-yml p1-yaml p9-alle-befunde p9-reihenfolge" skript '  if (!sawperm) regel(' '  if (0) regel('
mut m-p6-jeder-wert-oben "p6-read-all p6-write-all p6-block p6-flow-eintrag" skript 'if (val != "{}") regel(NR, "permissions-form"' 'if (0) regel(NR, "permissions-form"'
mut m-p6-skalar-am-job "p6-job-write-all p6-job-read-all p9-alle-befunde" skript 'else regel(NR, "job-permissions-form", "Job " jname ": permissions als Skalar: " val)' 'else { }'
mut m-p6-flow-am-job "p6-job-flow-eintrag" skript 'regel(NR, "job-permissions-form", "Job " jname ": permissions als Flow-Mapping mit Einträgen: " val)' '{ }'
mut m-p6-unbekannter-wert "p6-job-wert" skript 'if (val !~ /^(read|write|none)$/) regel(NR' 'if (0) regel(NR'
mut m-p6-job-pflicht "MUSS:p6-job-ohne-ok NICHT:p6-job-block-ok p6-job-leer-ok p6-job-write-all" skript 'function jobende(   z, bad) {' $'function jobende(   z, bad) {\n  if (injob && !pseen) regel(jline, "job-permissions-fehlt", "x"); pseen = 0' \
  'else if (key == "permissions") {' 'else if (key == "permissions") { pseen = 1'
# Punkt 7
mut m-p7-pflicht-entfaellt "p7-fehlt p5-schritt-ist-kein-job p9-alle-befunde p9-reihenfolge" skript '      regel(jline, "timeout-fehlt", "Job " jname ": timeout-minutes fehlt")' '      { }'
mut m-p7-obergrenze-61 "p7-61" skript 'z > 60' 'z > 61'
mut m-p7-obergrenze-360 "p7-61" skript 'length(tval) > 3 || z < 1 || z > 60' 'z < 1 || z > 360'
mut m-p7-untergrenze "p7-1-ok" skript 'z < 1' 'z < 2'
mut m-p7-null-angenommen "p7-0" skript 'z < 1' 'z < 0'
mut m-p7-ausdruck-angenommen "p7-ausdruck" skript '      else bad = 1' '      else bad = (tval !~ /^\$\{\{.*\}\}$/)'
mut m-p7-nicht-ganzzahl-angenommen "p7-nicht-ganzzahl p7-dezimal" skript '      else bad = 1' '      else bad = (tval !~ /^[0-9]+\.[0-9]+$/)'
mut m-p7-fuehrende-null-angenommen "p7-fuehrende-null" skript 'tval ~ /^(0|[1-9][0-9]*)$/' 'tval ~ /^[0-9]+$/'
mut m-p7-vorzeichen-angenommen "p7-vorzeichen" skript '      else bad = 1' '      else bad = (tval !~ /^\+[0-9]+$/)'
mut m-p7-dezimal-angenommen "p7-dezimal" skript '      else bad = 1' '      else bad = (tval !~ /^[0-9]+\.0$/)'
mut m-p7-wort-angenommen "p7-wort" skript '      else bad = 1' '      else bad = (tval !~ /^[A-Za-z]+$/)'
mut m-p7-job-ohne-runs-on "p7-uses-job-ok" skript '  if (injob && hasruns) {' '  if (injob) {'
# Punkt 8
mut m-p8-a "p8-modules p8-modules-praefix p8-modules-suffix" skript 'if [ "${bindung:0:1}" != "1" ]; then fehlt=' 'if false; then fehlt='
mut m-p8-dir-ohne-block "p8-dir-davor p8-dir-danach" skript 'inb && /^  dir:' '/^  dir:'
mut m-p8-ohne-ruecksetzung "p8-dir-danach" skript '/^[^ \t#]/ { inb = 0 }' ''
mut m-p8-modules-muster "p8-modules-praefix p8-modules-suffix" skript 'if (x == "workflows") a = 1' 'if (x ~ /workflows/) a = 1'
mut m-p8-b "p8-block" skript 'if [ "${bindung:1:1}" != "1" ]; then' 'if false; then'
mut m-p8-c "p8-dir p8-dir-schraeg p8-dir-davor p8-dir-danach" skript 'elif [ "${bindung:2:1}" != "1" ]; then' 'elif false; then'
mut m-p8-datei-fehlt-angenommen "p8-ohne-datei p9-reihenfolge" skript 'fehlt=".d-check.yml fehlt"' 'fehlt=""'
mut m-p8-anfuehrung-doppelt "p8-ok-doppelt" skript '["'"'"']?\.github\/workflows["'"'"']?[ \t]*(#.*)?$/ { c = 1 }' "'?\\.github\\/workflows'?[ \\t]*(#.*)?\$/ { c = 1 }"
mut m-p8-anfuehrung-einfach "p8-ok-einfach" skript '["'"'"']?\.github\/workflows["'"'"']?[ \t]*(#.*)?$/ { c = 1 }' '"?\.github\/workflows"?[ \t]*(#.*)?$/ { c = 1 }'
mut m-p8-schraegstrich "p8-dir-schraeg" skript '["'"'"']?\.github\/workflows["'"'"']?[ \t]*(#.*)?$/ { c = 1 }' '["'"'"']?\.github\/workflows\/?["'"'"']?[ \t]*(#.*)?$/ { c = 1 }'
mut m-p8-pins-selbst "p8-pins-nicht-hier" skript '  if ($0 ~ /^ *#/) next' $'  if ($0 ~ /^ *#/) next\n  if ($0 ~ /uses: [^.].*@v[0-9]/) regel(NR, "uses-pin-missing", "x")'
# Punkt 9
mut m-p9-datei-vor-al "p9-reihenfolge" skript 'for b in "${bef_al[@]}" "${bef_datei[@]}" "${bef_global[@]}"; do' 'for b in "${bef_datei[@]}" "${bef_al[@]}" "${bef_global[@]}"; do'
mut m-p9-global-vor-datei "p9-reihenfolge" skript 'for b in "${bef_al[@]}" "${bef_datei[@]}" "${bef_global[@]}"; do' 'for b in "${bef_al[@]}" "${bef_global[@]}" "${bef_datei[@]}"; do'
mut m-p9-ohne-sortierung "p9-alle-befunde" skript ' | sort -t "$(printf '"'\\t'"')" -k1,1n -s)' ')'
mut m-p9-erster-befund "p9-reihenfolge p9-alle-befunde p3-exit1" skript '  echo "$b" >&2' '  echo "$b" >&2; break'
mut m-p9-summe-ohne-al "p3-exit1 p3-exit1-leer p9-reihenfolge p10-unterverzeichnis" skript 'nbef=$((${#bef_al[@]} + ${#bef_datei[@]} + ${#bef_global[@]}))' 'nbef=$((${#bef_datei[@]} + ${#bef_global[@]}))'
mut m-p9-summe-nur-al "MUSS:p9-alle-befunde p9-reihenfolge NICHT:p3-exit1 p3-exit0" skript 'echo "workflows-check: $nbef Befund(e)" >&2' 'echo "workflows-check: ${#bef_al[@]} Befund(e)" >&2'
mut m-p3-exit1-leer-angenommen "p3-exit1-leer" skript 'if [ "${#bef_al[@]}" -eq 0 ]; then bef_al+=("workflows-check: actionlint endete mit Exit 1 ohne Zeile"); fi' ':'
mut m-p9-umgebung-exit1 "p3-exit2 p3-exit3 p3-exit3-rechte p3-exit125 p3-kein-docker p9-umgebung-vor-befund" skript $'      exit 2\n      ;;' $'      exit 1\n      ;;'
mut m-p9-eigene-vor-umgebung "p9-umgebung-vor-befund" skript $'      echo "workflows-check: actionlint nicht lauffähig (Exit $rc)" >&2' $'      printf \'%s\\n\' "${bef_datei[@]}" "${bef_global[@]}" >&2\n      echo "workflows-check: actionlint nicht lauffähig (Exit $rc)" >&2'
mut m-p9-verzeichnis-r-fehlt "p9-verzeichnis-ohne-r" skript '{ [ -r "$dir" ] && [ -x "$dir" ]; } ||' '{ [ -x "$dir" ]; } ||'
mut m-p9-verzeichnis-x-fehlt "p9-verzeichnis-ohne-x" skript '{ [ -r "$dir" ] && [ -x "$dir" ]; } ||' '{ [ -r "$dir" ]; } ||'
mut m-p9-nicht-lesbar "p9-nicht-lesbar" skript '[ -r "$dir/$n" ] || sperre "$WF/$n nicht lesbar"' '[ -r "$dir/$n" ] || :'
# Punkt 10
mut m-p10-plus "p10-n p10-q p10-t p10-makeflags" fragment $'\t@bash tools/harness/workflows-check.sh' $'\t+@bash tools/harness/workflows-check.sh'
mut m-p10-zweimal "p10-rezept p10-i" fragment $'\t@bash tools/harness/workflows-check.sh' $'\t@bash tools/harness/workflows-check.sh\n\t@bash tools/harness/workflows-check.sh'
mut m-p10-wurzel-ist-arbeitsverzeichnis "p10-unterverzeichnis" skript 'wurzel="${1:-$(cd "$hier/../.." && pwd)}"' 'wurzel="${1:-$(pwd)}"'
mut m-p10-schreibt "p10-schreibt-nichts" skript '  errf="$(mktemp)" || sperre "keine temporäre Datei"' $'  errf="$(mktemp)" || sperre "keine temporäre Datei"\n  : > "$wurzel/.spur"'
mut m-p10-kette-gate "p10-kette-gate" fragment 'GATE_CHECKS += workflows-check' ':'
mut m-p10-kette-gegenprobe "p10-kette-gegenprobe" fragment 'GEGENPROBEN_CONTAINER += workflows-check-gegenprobe' ':'
# Punkt 12
mut m-real-attrappe-image "real-ausdruck" env REAL_DOCKER "$attrappe"
mut m-real-modules "real-pins" env REAL_MODULES '[links]'

wait
fehler=0
if [ -s "$arbeit/uebersprungen" ]; then echo "workflows-check-gegenprobe: übersprungen (root liest jede Datei und jedes Verzeichnis): $(sort -u "$arbeit/uebersprungen" | tr '\n' ' ')" >&2; fi
if [ -s "$arbeit/fehler" ]; then cat "$arbeit/fehler" >&2; fehler=1; fi
if [ "$(wc -l < "$arbeit/gefahren")" -eq 0 ]; then echo "workflows-check-gegenprobe: keine Mutation und keine Kopie gefahren (GP_MUT passt auf nichts)" >&2; fehler=1; fi
if [ "$fehler" -eq 0 ]; then echo "workflows-check-gegenprobe: grün ($(wc -l < "$arbeit/gefahren") Läufe)"; fi
exit "$fehler"
