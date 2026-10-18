#!/usr/bin/env bash
# range-konfiguration-gegenprobe — prüft die Blöcke `commits` und `vcs` in .d-check.yml gegen die Punkte 11 und 12 von SPEC-064
# (harness/vertraege/SPEC-064-ci.md, Punkt 16). Je Fall ein Git-Repository in einem Temp-Verzeichnis mit Kopien von .d-check.yml, d-check.mk,
# harness/mk/doc-gate.mk und tools/harness/history-range-guard.sh (die Treiber unverändert) und einem Makefile, das doc-gate.mk einbindet;
# `make doc-commits` (Fälle p11-*) bzw. `make doc-immutable` (Fälle p12-*) läuft mit dem Image aus DCHECK_REF (übergibt `make`; der erste Lauf
# zieht es) über die Range des Falls. Ein Fall besteht, wenn Exit-Code, stdout (ohne die Zeile des Vorlaufs) und stderr vollständig dem erwarteten Text
# entsprechen; SHAs, Dateizahl und Zeilennummer in d-check.mk sind normiert. Je Mutation (Textersatz in der Kopie von .d-check.yml) müssen genau die
# benannten Fälle rot werden; die unveränderte Kopie besteht jeden Fall. Eine Mutation läuft über die Fälle ihrer Gruppe (p11 oder p12), die beiden
# Mutationen `Block fehlt` über alle Fälle. Der Name jedes Falls trägt die Nummer des Punktes. Ein Name in der Fallliste ohne Funktion `f_<name>` zählt als roter Fall.
# Abbruch (Punkt 18): Bei Ende, SIGINT und SIGTERM beendet das Skript erst seine Nachkommen und wartet auf sie, entfernt dann sein Temp-Verzeichnis
# und endet mit 130 bzw. 143; die Fälle p18-sigint und p18-sigterm starten dafür eine Kopie des Skripts mit eigenem TMPDIR (GP_INNEN=1 lässt die Kopie die p18-Fälle
# auslassen) und senden das Signal, sobald dort der erste Fall läuft; Mutationen der Kopie (Textersatz im Skript): die Falle fehlt, die Nachkommen werden nicht beendet,
# das Verzeichnis wird vor dem Ende der Nachkommen entfernt. Ein Container von d-check, der beim Signal schon läuft, endet von selbst (`--rm`); SIGKILL und SIGHUP sind nicht zugesagt.
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Mutationen, deren Name passt (die unveränderte Kopie nur, wenn das Muster `ok` trifft);
# GP_ENTDECKE=1 nennt je Mutation die roten Fälle; GP_PARALLEL=<n> (Vorgabe 4) setzt die Zahl der gleichzeitigen Läufe;
# GP_FALL=<name> gibt für die unveränderte Kopie das Protokoll dieses Falls aus.
# Braucht bash (mit `wait -n`), make, git, GNU sed, grep, pgrep, setsid, Docker und die Werkzeuge von coreutils (mktemp, sort, tr, wc, cat, cp, chmod, env); eine Version prüft das Skript nicht, zugesagt ist die Plattform von SPEC-064. `d-check.mk` ruft `docker` aus dem PATH; die Gegenprobe stellt dem Lauf eine Hülle voran, die $DOCKER mit dem PATH des Aufrufers ausführt und `-v` als `--mount` schreibt.
# Kein Fall hängt von Dateirechten ab; die Temp-Bäume sind für alle lesbar (chmod -R a+rX), weil das Image als Nicht-root liest.
# Aufrufarten von `make` (-n, -q, -t, -i, MAKEFLAGS) haben keinen Fall: kein Hook und kein Stand liest das Ergebnis (SPEC-064 Punkt 14).
set -uo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="$(cd "$hier/../.." && pwd)"
cfg0="$repo/.d-check.yml"
DOCKER="${DOCKER:-docker}"
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
melde() { echo "range-konfiguration-gegenprobe: ROT — $*" >> "$arbeit/fehler"; }

for t in make git sed grep mktemp sort tr wc cat cp chmod env; do
  command -v "$t" > /dev/null 2>&1 || { echo "range-konfiguration-gegenprobe: $t fehlt" >&2; exit 2; }
done
if [ -z "${DCHECK_REF:-}" ]; then
  echo "range-konfiguration-gegenprobe: DCHECK_REF fehlt (aufrufen über make range-konfiguration-gegenprobe)" >&2
  exit 2
fi
for f in .d-check.yml d-check.mk harness/mk/doc-gate.mk tools/harness/history-range-guard.sh; do
  [ -f "$repo/$f" ] || { echo "range-konfiguration-gegenprobe: $f fehlt im Arbeitsbaum" >&2; exit 2; }
done
# Vorlauf: das Image ist lokal oder wird hier gezogen, bevor die gleichzeitigen Läufe starten.
"$DOCKER" run --rm --network none "$DCHECK_REF" --help > /dev/null 2>&1 || {
  echo "range-konfiguration-gegenprobe: d-check ($DCHECK_REF) ist mit $DOCKER nicht lauffähig" >&2
  exit 2
}
mkdir "$arbeit/bin"
docker_echt="$(command -v "$DOCKER")" || { echo "range-konfiguration-gegenprobe: $DOCKER fehlt" >&2; exit 2; }
# Die Hülle schreibt `-v QUELLE:ZIEL:ro` als `--mount type=bind,…,readonly` (Docker legt bei `-v` einen fehlenden Quellpfad als root an, bei `--mount` meldet es einen Fehler).
cat > "$arbeit/bin/docker" << EOT
#!/bin/sh
PATH="\$GP_PFAD"
n=\$#
while [ "\$n" -gt 0 ]; do
  a="\$1"; shift; n=\$((n - 1))
  if [ "\$a" = -v ] && [ "\$n" -gt 0 ]; then
    v="\$1"; shift; n=\$((n - 1))
    q="\${v%%:*}"; z="\${v#*:}"; z="\${z%%:*}"
    set -- "\$@" --mount "type=bind,src=\$q,dst=\$z,readonly"
  else
    set -- "\$@" "\$a"
  fi
done
exec "$docker_echt" "\$@"
EOT
chmod +x "$arbeit/bin/docker"

export GP_PFAD="$PATH"
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1
export GIT_AUTHOR_NAME=gp GIT_AUTHOR_EMAIL=gp@example.invalid GIT_COMMITTER_NAME=gp GIT_COMMITTER_EMAIL=gp@example.invalid

FAELLE_P11=(
  p11-ohne-kennung p11-adr p11-lh p11-mr p11-slice p11-welle p11-adr-knapp p11-lh-knapp p11-mr-knapp p11-slice-knapp p11-welle-knapp
  p11-revert p11-merge-betreff p11-kennung-im-rumpf p11-kennung-im-kommentar p11-exempt-nicht-am-anfang p11-struktur-allein
)
FAELLE_P18=(p18-sigint p18-sigterm)
FAELLE_P12=(
  p12-kern p12-geschichte p12-ersetzt p12-ersetzt-anders p12-vorschlag p12-geloescht p12-link-ziel p12-index
)

# --- Bausteine der Fälle ---------------------------------------------------------------------------------------------------------------------------
git_() { git -C "$w" "$@" > /dev/null; }
# neu [zwei] — legt das Repository des Falls an (mit `zwei` auch die Proposed-ADR 0002) und committet die Basis; w ist seine Wurzel, base der Basis-Commit, cfg die Kopie von .d-check.yml.
neu() {
  mkdir -p "$run"
  w="$(mktemp -d "$run/t.XXXXXX")" || exit 1
  mkdir -p "$w/harness/mk" "$w/tools/harness" "$w/docs/plan/adr" "$w/docs/plan/planning/open" "$w/docs/plan/planning/in-progress"
  cp "$cfg" "$w/.d-check.yml"
  cp "$repo/d-check.mk" "$w/d-check.mk"
  cp "$repo/harness/mk/doc-gate.mk" "$w/harness/mk/doc-gate.mk"
  cp "$repo/tools/harness/history-range-guard.sh" "$w/tools/harness/history-range-guard.sh"
  printf 'GATE_CHECKS :=\ninclude harness/mk/doc-gate.mk\n' > "$w/Makefile"
  printf 'Ziel\n' > "$w/docs/plan/planning/open/s.md"
  printf '# ADR-Index\n\nzeile\n' > "$w/docs/plan/adr/README.md"
  printf '# ADR-0001: Eins\n\n**Status:** Accepted\n\nSiehe [Ziel](../planning/open/s.md).\n\n## Entscheidung\n\nText A.\n\n## Geschichte\n\n| Datum | Änderung |\n' > "$w/docs/plan/adr/0001-eins.md"
  if [ "${1-}" = zwei ]; then printf '# ADR-0002: Zwei\n\n**Status:** Proposed\n\nText A.\n' > "$w/docs/plan/adr/0002-zwei.md"; fi
  git -C "$w" init -q
  git_ add -A
  git_ commit -q -m "ADR-0001 Basis"
  base="$(git -C "$w" rev-parse HEAD)"
}
kommit() { git_ commit -q --allow-empty --cleanup=verbatim -m "$1"; }
adr1() { echo "$w/docs/plan/adr/0001-eins.md"; }
buche() { git_ add -A; git_ commit -q -m "ADR-0001 Änderung"; }

# erwarte <ziel> [<Befundzeile>…] — führt make im Baum aus und vergleicht Exit-Code, stdout und stderr vollständig mit dem erwarteten Text.
norm() { sed -E 's/[0-9a-f]{40}/<base>/g; s/\b[0-9a-f]{7}\b/<sha>/g; s/[0-9]+ Datei\(en\)/N Datei(en)/; s/d-check\.mk:[0-9]+/d-check.mk:N/'; }
erwarte() {
  local ziel="$1" n=$(($# - 1)) ecode=0 eout="" eerr b
  shift
  chmod -R a+rX "$w"
  ( cd "$w" && env -u MAKEFLAGS -u MFLAGS -u MAKELEVEL LC_ALL=C PATH="$arbeit/bin:$PATH" make -s "$ziel" RANGE="$base..HEAD" DCHECK_REF="$DCHECK_REF" ) > "$w.out" 2> "$w.err"
  code=$?
  out="$(grep -v '^history-range-guard:' "$w.out" | norm)"
  err="$(norm < "$w.err")"
  for b in "$@"; do eout+="${eout:+$'\n'}$b"; done
  eerr="d-check: N Datei(en) geprüft, $n Befund(e)"
  if [ "$n" -gt 0 ]; then ecode=2; eerr+=$'\n'"make: *** [d-check.mk:N: $ziel] Error 1"; fi
  if [ "$code" -ne "$ecode" ] || [ "$out" != "$eout" ] || [ "$err" != "$eerr" ]; then
    bad=1
    printf 'erwartet: Exit %s\nstdout:\n%s\nstderr:\n%s\ngelesen: Exit %s\nstdout:\n%s\nstderr:\n%s\n' "$ecode" "$eout" "$eerr" "$code" "$out" "$err" >> "$run/log.$fallname"
  fi
}
T=$'\t'
untraceable() { echo "<sha>:1${T}<sha>${T}commit-untraceable${T}$1"; }
KERN="docs/plan/adr/0001-eins.md:3${T}docs/plan/adr/0001-eins.md${T}core-drift-vcs${T}"

# --- Punkt 11: commits ----------------------------------------------------------------------------------------------------------------------------
f_p11_ohne_kennung() { neu; kommit foo; erwarte doc-commits "$(untraceable foo)"; }
f_p11_adr() { neu; kommit "Siehe ADR-0002 x"; erwarte doc-commits; }
f_p11_lh() { neu; kommit "Siehe LH-FA-01 x"; erwarte doc-commits; }
f_p11_mr() { neu; kommit "Siehe MR-001 x"; erwarte doc-commits; }
f_p11_slice() { neu; kommit "Siehe slice-aa-bb x"; erwarte doc-commits; }
f_p11_welle() { neu; kommit "Siehe welle-aa x"; erwarte doc-commits; }
f_p11_adr_knapp() { neu; kommit "Siehe ADR-123 x"; erwarte doc-commits "$(untraceable 'Siehe ADR-123 x')"; }
f_p11_lh_knapp() { neu; kommit "Siehe LH-FA-1 x"; erwarte doc-commits "$(untraceable 'Siehe LH-FA-1 x')"; }
f_p11_mr_knapp() { neu; kommit "Siehe MR-12 x"; erwarte doc-commits "$(untraceable 'Siehe MR-12 x')"; }
f_p11_slice_knapp() { neu; kommit "Siehe slice- x"; erwarte doc-commits "$(untraceable 'Siehe slice- x')"; }
f_p11_welle_knapp() { neu; kommit "Siehe welle- x"; erwarte doc-commits "$(untraceable 'Siehe welle- x')"; }
f_p11_revert() { neu; kommit 'Revert "x"'; erwarte doc-commits; }
f_p11_merge_betreff() { neu; kommit "Merge x"; erwarte doc-commits; }
f_p11_kennung_im_rumpf() { neu; kommit $'x\n\nSiehe MR-001'; erwarte doc-commits; }
f_p11_kennung_im_kommentar() { neu; kommit $'x\n\n# MR-001'; erwarte doc-commits "$(untraceable x)"; }
f_p11_exempt_nicht_am_anfang() { neu; kommit "x Revert y"; erwarte doc-commits "$(untraceable 'x Revert y')"; }
f_p11_struktur_allein() { neu; kommit "SPEC-001 x"; erwarte doc-commits "$(untraceable 'SPEC-001 x')"; }

# --- Punkt 12: vcs --------------------------------------------------------------------------------------------------------------------------------
f_p12_kern() { neu; sed -i 's/Text A/Text B/' "$(adr1)"; buche; erwarte doc-immutable "${KERN}Core einer immutablen Datei hat sich über die Commit-Range geändert"; }
f_p12_geschichte() { neu; printf '| 2026-10-18 | x |\n' >> "$(adr1)"; buche; erwarte doc-immutable; }
f_p12_ersetzt() { neu; sed -i 's/^\*\*Status:\*\* Accepted/**Status:** Superseded by [ADR-0002](0002-zwei.md)/' "$(adr1)"; buche; erwarte doc-immutable; }
f_p12_ersetzt_anders() {
  neu; sed -i 's/^\*\*Status:\*\* Accepted/**Status:** Rejected/' "$(adr1)"; buche
  erwarte doc-immutable "${KERN}unzulässiger Status-Übergang einer immutablen Datei"
}
f_p12_vorschlag() { neu zwei; sed -i 's/Text A/Text B/' "$w/docs/plan/adr/0002-zwei.md"; buche; erwarte doc-immutable; }
f_p12_geloescht() {
  neu; git_ rm -q docs/plan/adr/0001-eins.md; buche
  erwarte doc-immutable "docs/plan/adr/0001-eins.md:1${T}docs/plan/adr/0001-eins.md${T}core-drift-vcs${T}immutable Datei gelöscht oder umbenannt — der Pfad einer immutablen Datei ist stabil"
}
f_p12_link_ziel() {
  neu; git_ mv docs/plan/planning/open/s.md docs/plan/planning/in-progress/s.md
  sed -i 's#planning/open/s.md#planning/in-progress/s.md#' "$(adr1)"; buche; erwarte doc-immutable
}
f_p12_index() { neu; printf 'zeile2\n' >> "$w/docs/plan/adr/README.md"; buche; erwarte doc-immutable; }

# --- Punkt 18: Abbruch ------------------------------------------------------------------------------------------------------------------------------
warte_auf() { local i; for i in $(seq 1 300); do if "$@"; then return 0; fi; sleep 0.1; done; return 1; }
fall_laeuft() { [ "$(find "$t" -mindepth 3 -maxdepth 3 -name 't.*' 2> /dev/null | wc -l)" -ge 3 ]; }
reste_da() { [ -n "$(ls -A "$t" 2> /dev/null)" ] || pgrep -f -- "$t" > /dev/null 2>&1; }
# Die Kopie fährt mit einer Attrappe von docker (eine Sekunde Wartezeit, Exit 0): ihre Fälle werden rot, es zählt nur, wie sie endet.
printf '#!/bin/sh\nsleep 1\n' > "$arbeit/attrappe"
chmod +x "$arbeit/attrappe"
# abbruchfall <Signal> <Exit> — fährt $skript18 mit eigenem TMPDIR in eigener Sitzung, sendet das Signal, sobald dort der erste Fall läuft, und
# prüft danach Exit, leeren TMPDIR und Prozesse mit dem Temp-Pfad in der Befehlszeile.
abbruchfall() {
  local sig="$1" soll="$2" pid=0 rc=0 bg ok=1 reste
  mkdir -p "$run"
  t="$(mktemp -d "$run/tmp.XXXXXX")" || exit 1
  ( trap - INT TERM; exec env TMPDIR="$t" GP_INNEN=1 DOCKER="$arbeit/attrappe" DCHECK_REF="$DCHECK_REF" setsid bash -c 'echo $$ > "$1"; exec bash "$2"' _ "$t.pid" "$skript18" > "$t.out" 2> "$t.err" ) &
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
# warte_auf_reste — bis zu 8 s auf ein Verzeichnis oder einen Prozess, den ein nicht beendeter Nachkomme hinterlässt.
warte_auf_reste() { local i; for i in $(seq 1 80); do if reste_da; then return 0; fi; sleep 0.1; done; return 0; }
f_p18_sigint() { abbruchfall INT 130; }
f_p18_sigterm() { abbruchfall TERM 143; }

# --- Läufe ----------------------------------------------------------------------------------------------------------------------------------------
# faelle <gruppe> — schreibt die Namen der roten Fälle der Gruppe (p11, p12, p18, alle = p11 und p12, ok = alle und p18), einen je Zeile.
# FALL_ZUSATZ nennt einen weiteren Namen (Probe: ein Name ohne Funktion ist ein roter Fall).
faelle() {
  local n
  case "$1" in
    p11) set -- "${FAELLE_P11[@]}" ;;
    p12) set -- "${FAELLE_P12[@]}" ;;
    p18) set -- "${FAELLE_P18[@]}" ;;
    ok) set -- "${FAELLE_P11[@]}" "${FAELLE_P12[@]}"; [ -n "${GP_INNEN-}" ] || set -- "$@" "${FAELLE_P18[@]}" ;;
    *) set -- "${FAELLE_P11[@]}" "${FAELLE_P12[@]}" ;;
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

# mut <name> <rote Fälle> <Gruppe: p11|p12|alle> <alt> <neu> [<alt> <neu>…] — Textersatz in der Kopie von .d-check.yml; rot sind genau die benannten Fälle.
mut() {
  local name="$1" rote="$2" gruppe="$3" inhalt neu ist erwartet alt rest
  shift 3
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  (
    run="$arbeit/$name"; mkdir -p "$run"; cfg="$run/d-check.yml"
    inhalt="$(cat "$cfg0"; printf x)"; inhalt="${inhalt%x}"; neu="$inhalt"
    while [ "$#" -ge 2 ]; do
      alt="$1"
      if [[ $neu != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
      rest="${neu#*"$alt"}"
      if [[ $rest == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
      neu="${neu/"$alt"/"$2"}"
      shift 2
    done
    if [[ $neu == "$inhalt" ]]; then melde "Mutation '$name': Ersatz ändert nichts"; exit 0; fi
    printf '%s' "$neu" > "$cfg"
    ist="$(faelle "$gruppe" | LC_ALL=C sort -u | tr '\n' ' ')"
    erwartet="$(printf '%s\n' $rote | LC_ALL=C sort -u | tr '\n' ' ')"
    if [ -n "${GP_ENTDECKE-}" ]; then echo "$name => ${ist:-nichts}" >&2; fi
    echo "$name" >> "$arbeit/gefahren"
    [ "$ist" = "$erwartet" ] || melde "Mutation '$name': erwartet rot: ${rote:-nichts}; rot sind: ${ist:-nichts}"
  ) &
}

if [[ ok =~ ${GP_MUT:-.} ]]; then
  (
    run="$arbeit/ok"; mkdir -p "$run"; cfg="$cfg0"
    skript18="$selbst"
    ist="$(faelle ok | LC_ALL=C sort -u | tr '\n' ' ')"
    echo ok >> "$arbeit/gefahren"
    [ -z "$ist" ] || melde "unveränderte Kopie ist nicht grün: $ist"
    protokoll_zeigen
    if [ -n "${GP_FALL-}" ]; then sed 's/^/stderr: /' "$run/stderr" >&2; fi
  ) &
fi

# mut18 <name> <rote Fälle> <alt> <neu> — Textersatz im Skript (Kopie in einem Baum mit den Dateien, die es braucht); rot sind genau die benannten p18-Fälle.
mut18() {
  local name="$1" rote="$2" alt="$3" inhalt neu ist erwartet f m
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  (
    run="$arbeit/$name"; m="$run/baum"; mkdir -p "$m/tools/harness" "$m/harness/mk"
    for f in .d-check.yml d-check.mk harness/mk/doc-gate.mk tools/harness/history-range-guard.sh; do cp "$repo/$f" "$m/$f"; done
    inhalt="$(cat "$selbst"; printf x)"; inhalt="${inhalt%x}"
    if [[ $inhalt != *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $alt"; exit 0; fi
    neu="${inhalt/"$alt"/"$4"}"
    if [[ ${inhalt#*"$alt"} == *"$alt"* ]]; then melde "Mutation '$name': Textstelle nicht eindeutig: $alt"; exit 0; fi
    printf '%s' "$neu" > "$m/tools/harness/range-konfiguration-gegenprobe.sh"
    skript18="$m/tools/harness/range-konfiguration-gegenprobe.sh"; cfg="$cfg0"
    ist="$(faelle p18 | LC_ALL=C sort -u | tr '\n' ' ')"
    erwartet="$(printf '%s\n' $rote | LC_ALL=C sort -u | tr '\n' ' ')"
    if [ -n "${GP_ENTDECKE-}" ]; then echo "$name => ${ist:-nichts}" >&2; fi
    echo "$name" >> "$arbeit/gefahren"
    [ "$ist" = "$erwartet" ] || melde "Mutation '$name': erwartet rot: ${rote:-nichts}; rot sind: ${ist:-nichts}"
  ) &
}

# Probe der Fallliste: ein Name ohne Funktion ist rot.
if [[ fall-ohne-funktion =~ ${GP_MUT:-.} ]]; then
  (
    run="$arbeit/fall-ohne-funktion"; mkdir -p "$run"; cfg="$cfg0"
    ist="$(FALL_ZUSATZ=p11-tippfehler faelle p11 | LC_ALL=C sort -u | tr '\n' ' ')"
    echo fall-ohne-funktion >> "$arbeit/gefahren"
    [ "$ist" = "p11-tippfehler " ] || melde "Fall ohne Funktion: erwartet rot: p11-tippfehler; rot sind: ${ist:-nichts}"
  ) &
fi

# --- Mutationen -------------------------------------------------------------------------------------------------------------------------------------
BLOCK_COMMITS=$'commits:\n  id-patterns:\n    - \'ADR-[0-9]{4}\'\n    - \'LH-[A-Z]{2}-[0-9]{2}\'\n    - \'MR-[0-9]{3}\'\n    - \'(slice|welle)-[a-z0-9]+(-[a-z0-9]+)*\'\n  exempt-pattern: \'^(Merge |Revert )\'\n'
BLOCK_VCS=$'vcs:\n  paths: ["docs/plan/adr/[0-9]*.md"]\n  immutable-when: \'^\\*\\*Status:\\*\\* Accepted\'\n  exclude-sections: [Geschichte]\n  status-line: \'^\\*\\*Status:\\*\\*\'\n  head-allow: \'^\\*\\*Status:\\*\\* (Accepted|Superseded by )\'\n  ignore-link-targets: true\n'
P_ADR=$'    - \'ADR-[0-9]{4}\'\n'
P_LH=$'    - \'LH-[A-Z]{2}-[0-9]{2}\'\n'
P_MR=$'    - \'MR-[0-9]{3}\'\n'
P_SW=$'    - \'(slice|welle)-[a-z0-9]+(-[a-z0-9]+)*\'\n'
E_AUS=$'  exempt-pattern: \'^(Merge |Revert )\'\n'
V_ACC='immutable-when: '"'"'^\*\*Status:\*\* Accepted'"'"
V_EXCL=$'  exclude-sections: [Geschichte]\n  status-line'
V_STAT='status-line: '"'"'^\*\*Status:\*\*'"'"
V_PATHZ='  paths: ["docs/plan/adr/[0-9]*.md"]'$'\n'
V_ACCZ='  immutable-when: '"'"'^\*\*Status:\*\* Accepted'"'"$'\n'
V_STATZ='  status-line: '"'"'^\*\*Status:\*\*'"'"$'\n'
V_HEAD='head-allow: '"'"'^\*\*Status:\*\* (Accepted|Superseded by )'"'"
V_IGN=$'  ignore-link-targets: true\n'
V_PATHS='paths: ["docs/plan/adr/[0-9]*.md"]'

# Punkt 11
mut m-p11-block-fehlt "p11-ohne-kennung p11-kennung-im-kommentar p11-exempt-nicht-am-anfang p11-struktur-allein p11-adr-knapp p11-lh-knapp p11-mr-knapp p11-slice-knapp p11-welle-knapp" alle "$BLOCK_COMMITS" ''
mut m-p11-muster-adr-fehlt "p11-adr" p11 "$P_ADR" ''
mut m-p11-muster-lh-fehlt "p11-lh" p11 "$P_LH" ''
mut m-p11-muster-mr-fehlt "p11-mr p11-kennung-im-rumpf" p11 "$P_MR" ''
mut m-p11-muster-slice-welle-fehlt "p11-slice p11-welle" p11 "$P_SW" ''
mut m-p11-alternative-slice-fehlt "p11-slice" p11 '(slice|welle)-' '(welle)-'
mut m-p11-alternative-welle-fehlt "p11-welle" p11 '(slice|welle)-' '(slice)-'
mut m-p11-muster-adr-geweitet "p11-adr-knapp" p11 "'ADR-[0-9]{4}'" "'ADR-[0-9]'"
mut m-p11-muster-lh-geweitet "p11-lh-knapp" p11 "'LH-[A-Z]{2}-[0-9]{2}'" "'LH-[A-Z]{1,2}-[0-9]'"
mut m-p11-muster-mr-geweitet "p11-mr-knapp" p11 "'MR-[0-9]{3}'" "'MR-[0-9]'"
mut m-p11-muster-slice-welle-geweitet "p11-slice-knapp p11-welle-knapp" p11 "'(slice|welle)-[a-z0-9]+(-[a-z0-9]+)*'" "'(slice|welle)-'"
mut m-p11-muster-adr-verankert "p11-adr" p11 "'ADR-[0-9]{4}'" "'^ADR-[0-9]{4}'"
mut m-p11-muster-lh-verankert "p11-lh" p11 "'LH-[A-Z]{2}-[0-9]{2}'" "'^LH-[A-Z]{2}-[0-9]{2}'"
mut m-p11-muster-mr-verankert "p11-mr p11-kennung-im-rumpf" p11 "'MR-[0-9]{3}'" "'^MR-[0-9]{3}'"
mut m-p11-muster-slice-welle-verankert "p11-slice p11-welle" p11 "'(slice|welle)-[a-z0-9]+(-[a-z0-9]+)*'" "'^(slice|welle)-[a-z0-9]+(-[a-z0-9]+)*'"
mut m-p11-muster-struktur "p11-struktur-allein" p11 "$P_SW" "${P_SW}    - 'SPEC-[0-9]{3}'"$'\n'
mut m-p11-ausnahme-fehlt "p11-revert p11-merge-betreff" p11 "$E_AUS" ''
mut m-p11-ausnahme-merge-fehlt "p11-merge-betreff" p11 '^(Merge |Revert )' '^(Revert )'
mut m-p11-ausnahme-revert-fehlt "p11-revert" p11 '^(Merge |Revert )' '^(Merge )'
mut m-p11-ausnahme-ohne-anker "p11-exempt-nicht-am-anfang" p11 '^(Merge |Revert )' '(Merge |Revert )'
# Punkt 12
mut m-p12-block-fehlt "p12-kern p12-ersetzt-anders p12-geloescht" alle "$BLOCK_VCS" ''
mut m-p12-paths-leer "p12-kern p12-ersetzt-anders p12-geloescht" p12 "$V_PATHS" 'paths: []'
mut m-p12-paths-fehlt "p12-kern p12-ersetzt-anders p12-geloescht" p12 "$V_PATHZ" ''
mut m-p12-immutable-when-jeder-status "p12-vorschlag" p12 "$V_ACC" "immutable-when: '^\\*\\*Status:\\*\\*'"
mut m-p12-immutable-when-anderer-status "p12-kern p12-ersetzt-anders p12-geloescht" p12 "$V_ACC" "immutable-when: '^\\*\\*Status:\\*\\* Rejected'"
mut m-p12-immutable-when-fehlt "p12-kern p12-geschichte p12-ersetzt p12-ersetzt-anders p12-vorschlag p12-geloescht p12-link-ziel p12-index" p12 "$V_ACCZ" ''
mut m-p12-abschnitt-fehlt "p12-geschichte" p12 "$V_EXCL" '  status-line'
mut m-p12-abschnitt-falsch "p12-geschichte" p12 "$V_EXCL" $'  exclude-sections: [Historie]\n  status-line'
mut m-p12-status-zeile-fehlt "p12-kern p12-geschichte p12-ersetzt p12-ersetzt-anders p12-vorschlag p12-link-ziel p12-index" p12 "$V_STATZ" ''
mut m-p12-status-zeile-eng "p12-ersetzt p12-ersetzt-anders" p12 "$V_STAT" "status-line: '^\\*\\*Status:\\*\\* Accepted'"
mut m-p12-head-allow-fehlt "p12-ersetzt-anders" p12 "  $V_HEAD"$'\n' ''
mut m-p12-head-allow-eng "p12-ersetzt" p12 "$V_HEAD" "head-allow: '^\\*\\*Status:\\*\\* Accepted'"
mut m-p12-head-allow-ohne-accepted "p12-geschichte p12-index p12-link-ziel p12-vorschlag" p12 "$V_HEAD" "head-allow: '^\\*\\*Status:\\*\\* Superseded by '"
mut m-p12-link-ziel-fehlt "p12-link-ziel" p12 "$V_IGN" ''
mut m-p12-link-ziel-aus "p12-link-ziel" p12 "$V_IGN" $'  ignore-link-targets: false\n'
# Punkt 18
mut18 m-p18-int-falle-fehlt "p18-sigint" $'\ntrap abbruch_int INT\n' $'\n'
mut18 m-p18-term-falle-fehlt "p18-sigterm" $'\ntrap abbruch_term TERM\n' $'\n'
mut18 m-p18-int-kinder-bleiben "p18-sigint" $'\nabbruch_int() { trap - EXIT INT TERM; beende_kinder; ' $'\nabbruch_int() { trap - EXIT INT TERM; '
mut18 m-p18-term-kinder-bleiben "p18-sigterm" $'\nabbruch_term() { trap - EXIT INT TERM; beende_kinder; ' $'\nabbruch_term() { trap - EXIT INT TERM; '
mut18 m-p18-int-verzeichnis-zu-frueh "p18-sigint" $'\nabbruch_int() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 130; }' $'\nabbruch_int() { trap - EXIT INT TERM; raeume_auf; sleep 5; beende_kinder; exit 130; }'
mut18 m-p18-term-verzeichnis-zu-frueh "p18-sigterm" $'\nabbruch_term() { trap - EXIT INT TERM; beende_kinder; raeume_auf; exit 143; }' $'\nabbruch_term() { trap - EXIT INT TERM; raeume_auf; sleep 5; beende_kinder; exit 143; }'

wait
fehler=0
if [ -s "$arbeit/fehler" ]; then cat "$arbeit/fehler" >&2; fehler=1; fi
if [ "$(wc -l < "$arbeit/gefahren")" -eq 0 ]; then echo "range-konfiguration-gegenprobe: keine Mutation und keine Kopie gefahren (GP_MUT passt auf nichts)" >&2; fehler=1; fi
if [ "$fehler" -eq 0 ]; then echo "range-konfiguration-gegenprobe: grün ($(wc -l < "$arbeit/gefahren") Läufe)"; fi
exit "$fehler"
