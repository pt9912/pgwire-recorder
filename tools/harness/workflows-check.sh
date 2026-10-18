#!/usr/bin/env bash
# workflows-check — prüft jede Workflow-Datei unter .github/workflows nach
# harness/vertraege/SPEC-063-workflows-check.md: Syntax und Ausdrücke durch actionlint im Image
# aus WORKFLOWS_ACTIONLINT, Rechte und Zeitgrenzen durch die Regeln dieses Skripts, die Bindung der
# Pin-Prüfung in .d-check.yml. Was die Regeln lesen, steht in Punkt 5 des Vertrags.
# Aufruf: workflows-check.sh [<wurzel>] (ohne Wurzel das Verzeichnis zwei Ebenen über dem Skript)
# Umgebung: DOCKER (Vorgabe docker), WORKFLOWS_ACTIONLINT (Image mit Digest, harness/mk/workflows.mk)
# Ausgang: 0 ohne Befund und ohne Prüfbereich, 1 mit Befund, 2 bei Umgebungsfehler.
#
# GEPRÜFT DURCH tools/harness/workflows-check-gegenprobe.sh, Fälle nach Punkt (p<n>-*).
set -o pipefail
export LC_ALL=C

WF=".github/workflows"
DOCKER="${DOCKER:-docker}"

sperre() { echo "workflows-check: $*" >&2; exit 2; }

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
wurzel="${1:-$(cd "$hier/../.." && pwd)}"
wurzel="$(cd -- "$wurzel" 2> /dev/null && pwd -P)" || sperre "Wurzel nicht lesbar"
dir="$wurzel/$WF"

keine() { echo "workflows-check: keine Workflow-Dateien unter .github/workflows"; exit 0; }

[ -d "$dir" ] || keine
{ [ -r "$dir" ] && [ -x "$dir" ]; } || sperre "$WF nicht lesbar"

# Punkt 1: Prüfbereich.
kand=()
geordnet=()
hinweise=()
alle=0
while IFS= read -r -d '' p; do
  n="${p##*/}"
  alle=$((alle + 1))
  case "$n" in *.yml | *.yaml) endung=1 ;; *) endung=0 ;; esac
  if [ -L "$p" ]; then
    if [ "$endung" -eq 1 ]; then geordnet+=("s:$n"); else hinweise+=("hinweis: $WF/$n wird nicht geprüft (Endung)"); fi
  elif [ -d "$p" ]; then
    while IFS= read -r -d '' q; do
      hinweise+=("hinweis: ${q#"$wurzel/"} wird nicht geprüft (Unterverzeichnis)")
    done < <(find "$p" ! -type d -print0 | sort -z)
  elif [ -f "$p" ] && [ "$endung" -eq 1 ]; then
    kand+=("$n")
    geordnet+=("k:$n")
  else
    hinweise+=("hinweis: $WF/$n wird nicht geprüft (Endung)")
  fi
done < <(find "$dir" -mindepth 1 -maxdepth 1 -print0 | sort -z)
for h in "${hinweise[@]}"; do echo "$h"; done

# Punkt 2: kein Prüfbereich.
if [ "$alle" -eq 0 ]; then keine; fi
if [ "${#geordnet[@]}" -eq 0 ]; then keine; fi

for n in "${kand[@]}"; do
  [ -r "$dir/$n" ] || sperre "$WF/$n nicht lesbar"
done

# Punkt 5 bis 7: Regeln je Datei. Ein Satz je Befund: Zeile, Tabulator, Regel, Tabulator, Text.
read -r -d '' REGELN << 'EOT'
function wert(s,   v) {
  v = s
  sub(/^[ ]+/, "", v)
  if (v ~ /^#/) v = ""
  sub(/[ \t]+#.*$/, "", v)
  sub(/[ \t]+$/, "", v)
  gsub(/\t/, " ", v)
  return v
}
function form(l, t) { nf++; fl[nf] = l; ft[nf] = t }
function regel(l, r, t) { nr++; rl[nr] = l; rr[nr] = r; rt[nr] = t }
function jobende(   z, bad) {
  if (injob && hasruns) {
    if (!tseen) {
      regel(jline, "timeout-fehlt", "Job " jname ": timeout-minutes fehlt")
    } else {
      if (tval ~ /^(0|[1-9][0-9]*)$/) { z = tval + 0; bad = (length(tval) > 3 || z < 1 || z > 60) }
      else bad = 1
      if (bad) regel(tline, "timeout-wert", "Job " jname ": timeout-minutes " (tval == "" ? "ohne Wert" : tval) " ist keine Ganzzahl von 1 bis 60")
    }
  }
  injob = 0; inperm = 0
}
BEGIN { sect = ""; firstseen = 0; sawperm = 0 }
{
  if ($0 ~ /^[ \t]*$/) next
  if ($0 ~ /^ *#/) next
  if ($0 ~ /^ *\t/) { form(NR, "Tabulator im Einzug"); next }
  match($0, /^ */); ind = RLENGTH
  rest = substr($0, ind + 1)
  haskey = (rest !~ /^-/ && rest ~ /^[^ ][^:]*:([ ]|$)/)
  if (ind == 0) {
    if (!haskey) next
    jobende()
    key = substr(rest, 1, index(rest, ":") - 1)
    val = wert(substr(rest, index(rest, ":") + 1))
    sect = key
    if (key == "permissions") {
      sawperm = 1
      if (val != "{}") regel(NR, "permissions-form", "oben ist nur {} zulässig, gelesen: " (val == "" ? "ein Block" : val))
    } else if (key == "jobs") {
      firstseen = 0
      if (val != "") form(NR, "jobs: mit Wert in der Zeile: " val)
    }
    next
  }
  if (sect != "jobs") next
  if (!firstseen) {
    firstseen = 1
    if (ind != 2) form(NR, "Job-Schlüssel mit Einzug " ind " statt 2")
  }
  if (ind == 2 && haskey) {
    jobende()
    injob = 1; hasruns = 0; tseen = 0; jline = NR; jname = substr(rest, 1, index(rest, ":") - 1)
    next
  }
  if ((ind == 1 || ind == 3) && haskey) { form(NR, "Job-Schlüssel mit Einzug " ind " statt 2"); next }
  if (ind <= 4) inperm = 0
  if (!injob) next
  if (ind == 4 && haskey) {
    key = substr(rest, 1, index(rest, ":") - 1)
    val = wert(substr(rest, index(rest, ":") + 1))
    if (key == "runs-on") hasruns = 1
    else if (key == "timeout-minutes") { if (!tseen) { tseen = 1; tval = val; tline = NR } }
    else if (key == "permissions") {
      if (val == "{}") { }
      else if (val == "") inperm = 1
      else if (val ~ /^\{/) regel(NR, "job-permissions-form", "Job " jname ": permissions als Flow-Mapping mit Einträgen: " val)
      else regel(NR, "job-permissions-form", "Job " jname ": permissions als Skalar: " val)
    }
    next
  }
  if (inperm && ind > 4 && haskey) {
    key = substr(rest, 1, index(rest, ":") - 1)
    val = wert(substr(rest, index(rest, ":") + 1))
    if (val !~ /^(read|write|none)$/) regel(NR, "job-permissions-wert", "Job " jname ": " key ": " (val == "" ? "ohne Wert" : val) " ist nicht read, write oder none")
  }
}
END {
  jobende()
  if (nf > 0) { for (i = 1; i <= nf; i++) print fl[i] "\tform\t" ft[i]; exit }
  if (!sawperm) regel(1, "permissions-fehlt", "oben fehlt permissions: {}")
  for (i = 1; i <= nr; i++) print rl[i] "\t" rr[i] "\t" rt[i]
}
EOT

# Punkt 8: Bindung der Pin-Prüfung. Ausgabe: drei Ziffern für modules, workflows: und dir:.
read -r -d '' BINDUNG << 'EOT'
BEGIN { a = 0; b = 0; c = 0; inb = 0 }
/^modules:[ \t]*\[/ {
  s = $0; sub(/^modules:[ \t]*\[/, "", s); sub(/\].*$/, "", s)
  n = split(s, e, ",")
  for (i = 1; i <= n; i++) { x = e[i]; gsub(/^[ \t]+|[ \t]+$/, "", x); if (x == "workflows") a = 1 }
  inb = 0; next
}
/^workflows:[ \t]*(#.*)?$/ { b = 1; inb = 1; next }
/^[^ \t#]/ { inb = 0 }
inb && /^  dir:[ \t]*["']?\.github\/workflows["']?[ \t]*(#.*)?$/ { c = 1 }
END { print a b c }
EOT

bef_datei=()
for g in "${geordnet[@]}"; do
  n="${g#?:}"
  rel="$WF/$n"
  if [ "${g%%:*}" = "s" ]; then
    bef_datei+=("$rel:1: symlink: ist ein Symlink (Ziel $(readlink -- "$dir/$n")) und wird nicht gefolgt")
    continue
  fi
  while IFS=$'\t' read -r zeile regelname text; do
    [ -n "$regelname" ] || continue
    bef_datei+=("$rel:$zeile: $regelname: $text")
  done < <(awk "$REGELN" "$dir/$n" | sort -t "$(printf '\t')" -k1,1n -s)
done

bef_global=()
for k in actionlint.yaml actionlint.yml; do
  if [ -e "$wurzel/.github/$k" ] || [ -L "$wurzel/.github/$k" ]; then
    bef_global+=(".github/$k:1: konfiguration: actionlint läse diese Datei und könnte Befunde ausblenden")
  fi
done
if [ "${#kand[@]}" -gt 0 ]; then
  if [ -f "$wurzel/.d-check.yml" ]; then
    bindung="$(awk "$BINDUNG" "$wurzel/.d-check.yml")"
    fehlt=""
    if [ "${bindung:0:1}" != "1" ]; then fehlt="modules: führt workflows nicht"; fi
    if [ "${bindung:1:1}" != "1" ]; then
      fehlt="${fehlt:+$fehlt; }Zeile workflows: fehlt"
    elif [ "${bindung:2:1}" != "1" ]; then
      fehlt="${fehlt:+$fehlt; }dir: .github/workflows unter workflows: fehlt"
    fi
  else
    fehlt=".d-check.yml fehlt"
  fi
  if [ -n "$fehlt" ]; then
    bef_global+=("$WF/${kand[0]}:1: pins-ungeprueft: die Pins prüft make docs-check, wenn .d-check.yml es bindet; $fehlt")
  fi
fi

# Punkt 3: actionlint, genau einmal, wenn es Dateien gibt.
bef_al=()
if [ "${#kand[@]}" -gt 0 ]; then
  [ -n "${WORKFLOWS_ACTIONLINT:-}" ] || sperre "WORKFLOWS_ACTIONLINT ist leer"
  dateien=()
  for n in "${kand[@]}"; do dateien+=("$WF/$n"); done
  errf="$(mktemp)" || sperre "keine temporäre Datei"
  trap 'rm -f "$errf"' EXIT
  al="$("$DOCKER" run --rm --network none -v "$wurzel:/repo:ro" -w /repo "$WORKFLOWS_ACTIONLINT" -oneline "${dateien[@]}" 2> "$errf")"
  rc=$?
  case "$rc" in
    0) ;;
    1)
      while IFS= read -r z; do
        [ -n "$z" ] && bef_al+=("$z")
      done <<< "$al"
      if [ "${#bef_al[@]}" -eq 0 ]; then bef_al+=("workflows-check: actionlint endete mit Exit 1 ohne Zeile"); fi
      ;;
    *)
      echo "workflows-check: actionlint nicht lauffähig (Exit $rc)" >&2
      cat "$errf" >&2
      if [ "$rc" -eq 3 ] && grep -q 'permission denied' "$errf"; then
        echo "workflows-check: Hinweis: der Container liest als Nicht-root; Wurzel und .github/workflows brauchen Leserechte für alle" >&2
      fi
      exit 2
      ;;
  esac
fi

# Punkt 9: Ausgabe.
for b in "${bef_al[@]}" "${bef_datei[@]}" "${bef_global[@]}"; do
  echo "$b" >&2
done
nbef=$((${#bef_al[@]} + ${#bef_datei[@]} + ${#bef_global[@]}))
if [ "$nbef" -gt 0 ]; then
  echo "workflows-check: $nbef Befund(e)" >&2
  exit 1
fi
echo "workflows-check: ${#kand[@]} Datei(en) geprüft, kein Befund"
