#!/usr/bin/env bash
# range-konfiguration-gegenprobe — prüft die Blöcke `commits` und `vcs` in .d-check.yml gegen die Punkte 11 und 12 von SPEC-064
# (harness/vertraege/SPEC-064-ci.md, Punkt 16). Je Fall ein Git-Repository in einem Temp-Verzeichnis mit Kopien von .d-check.yml, d-check.mk,
# harness/mk/doc-gate.mk und tools/harness/history-range-guard.sh (die Treiber unverändert) und einem Makefile, das doc-gate.mk einbindet;
# `make doc-commits` (Fälle p11-*) bzw. `make doc-immutable` (Fälle p12-*) läuft mit dem Image aus DCHECK_REF (übergibt `make`; der erste Lauf
# zieht es) über die Range des Falls. Ein Fall besteht, wenn Exit-Code, stdout (ohne die Zeile des Vorlaufs) und stderr vollständig dem erwarteten Text
# entsprechen; SHAs, Dateizahl und Zeilennummer in d-check.mk sind normiert. Je Mutation (Textersatz in der Kopie von .d-check.yml) müssen genau die
# benannten Fälle rot werden; die unveränderte Kopie besteht jeden Fall. Eine Mutation läuft über die Fälle ihrer Gruppe (p11 oder p12), die beiden
# Mutationen `Block fehlt` über alle Fälle. Der Name jedes Falls trägt die Nummer des Punktes.
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Mutationen, deren Name passt (die unveränderte Kopie nur, wenn das Muster `ok` trifft);
# GP_ENTDECKE=1 nennt je Mutation die roten Fälle; GP_PARALLEL=<n> (Vorgabe 4) setzt die Zahl der gleichzeitigen Läufe;
# GP_FALL=<name> gibt für die unveränderte Kopie das Protokoll dieses Falls aus.
# Braucht bash ≥ 4, make, git ≥ 2.32, GNU sed, grep, Docker und die Werkzeuge von coreutils (mktemp, sort, tr, wc, cat, cp, chmod, env). `d-check.mk` ruft `docker` aus dem PATH; die Gegenprobe stellt dem Lauf eine Hülle voran, die $DOCKER mit dem PATH des Aufrufers ausführt.
# Kein Fall hängt von Dateirechten ab; die Temp-Bäume sind für alle lesbar (chmod -R a+rX), weil das Image als Nicht-root liest.
# Aufrufarten von `make` (-n, -q, -t, -i, MAKEFLAGS) haben keinen Fall: kein Hook und kein Stand liest das Ergebnis (SPEC-064 Punkt 14).
set -uo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="$(cd "$hier/../.." && pwd)"
cfg0="$repo/.d-check.yml"
DOCKER="${DOCKER:-docker}"
arbeit="$(mktemp -d)"
trap 'chmod -R u+rwX "${arbeit:?}" 2> /dev/null; rm -rf "${arbeit:?}"' EXIT
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
printf '#!/bin/sh\nPATH="$GP_PFAD"\nexec "%s" "$@"\n' "$docker_echt" > "$arbeit/bin/docker"
chmod +x "$arbeit/bin/docker"

export GP_PFAD="$PATH"
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1
export GIT_AUTHOR_NAME=gp GIT_AUTHOR_EMAIL=gp@example.invalid GIT_COMMITTER_NAME=gp GIT_COMMITTER_EMAIL=gp@example.invalid

FAELLE_P11=(
  p11-ohne-kennung p11-adr p11-lh p11-mr p11-slice p11-welle p11-revert p11-merge-betreff
  p11-kennung-im-rumpf p11-kennung-im-kommentar p11-exempt-nicht-am-anfang p11-struktur-allein
)
FAELLE_P12=(
  p12-kern p12-geschichte p12-ersetzt p12-ersetzt-anders p12-vorschlag p12-geloescht p12-link-ziel p12-index
)

# --- Bausteine der Fälle ---------------------------------------------------------------------------------------------------------------------------
git_() { git -C "$w" "$@" > /dev/null; }
# neu [zwei] — legt das Repository des Falls an (mit `zwei` auch die Proposed-ADR 0002) und committet die Basis; w ist seine Wurzel, base der Basis-Commit, cfg die Kopie von .d-check.yml.
neu() {
  w="$(mktemp -d "$run/t.XXXXXX")"
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
f_p11_adr() { neu; kommit "ADR-0002 x"; erwarte doc-commits; }
f_p11_lh() { neu; kommit "LH-FA-01 x"; erwarte doc-commits; }
f_p11_mr() { neu; kommit "MR-001 x"; erwarte doc-commits; }
f_p11_slice() { neu; kommit "slice-aa-bb x"; erwarte doc-commits; }
f_p11_welle() { neu; kommit "welle-aa x"; erwarte doc-commits; }
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

# --- Läufe ----------------------------------------------------------------------------------------------------------------------------------------
# faelle <gruppe> — schreibt die Namen der roten Fälle der Gruppe (p11, p12 oder alle), einen je Zeile.
faelle() {
  local n
  case "$1" in
    p11) set -- "${FAELLE_P11[@]}" ;;
    p12) set -- "${FAELLE_P12[@]}" ;;
    *) set -- "${FAELLE_P11[@]}" "${FAELLE_P12[@]}" ;;
  esac
  for n in "$@"; do
    fallname="$n"
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
    ist="$(faelle alle | LC_ALL=C sort -u | tr '\n' ' ')"
    echo ok >> "$arbeit/gefahren"
    [ -z "$ist" ] || melde "unveränderte Kopie ist nicht grün: $ist"
    protokoll_zeigen
    if [ -n "${GP_FALL-}" ]; then sed 's/^/stderr: /' "$run/stderr" >&2; fi
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
V_HEAD='head-allow: '"'"'^\*\*Status:\*\* (Accepted|Superseded by )'"'"
V_IGN=$'  ignore-link-targets: true\n'
V_PATHS='paths: ["docs/plan/adr/[0-9]*.md"]'

# Punkt 11
mut m-p11-block-fehlt "p11-ohne-kennung p11-kennung-im-kommentar p11-exempt-nicht-am-anfang p11-struktur-allein" alle "$BLOCK_COMMITS" ''
mut m-p11-muster-adr-fehlt "p11-adr" p11 "$P_ADR" ''
mut m-p11-muster-lh-fehlt "p11-lh" p11 "$P_LH" ''
mut m-p11-muster-mr-fehlt "p11-mr p11-kennung-im-rumpf" p11 "$P_MR" ''
mut m-p11-muster-slice-welle-fehlt "p11-slice p11-welle" p11 "$P_SW" ''
mut m-p11-alternative-slice-fehlt "p11-slice" p11 '(slice|welle)-' '(welle)-'
mut m-p11-alternative-welle-fehlt "p11-welle" p11 '(slice|welle)-' '(slice)-'
mut m-p11-muster-verankert "p11-kennung-im-rumpf" p11 \
  "'ADR-[0-9]{4}'" "'^ADR-[0-9]{4}'" "'LH-[A-Z]{2}-[0-9]{2}'" "'^LH-[A-Z]{2}-[0-9]{2}'" \
  "'MR-[0-9]{3}'" "'^MR-[0-9]{3}'" "'(slice|welle)-[a-z0-9]+(-[a-z0-9]+)*'" "'^(slice|welle)-[a-z0-9]+(-[a-z0-9]+)*'"
mut m-p11-muster-struktur "p11-struktur-allein" p11 "$P_SW" "${P_SW}    - 'SPEC-[0-9]{3}'"$'\n'
mut m-p11-ausnahme-fehlt "p11-revert p11-merge-betreff" p11 "$E_AUS" ''
mut m-p11-ausnahme-merge-fehlt "p11-merge-betreff" p11 '^(Merge |Revert )' '^(Revert )'
mut m-p11-ausnahme-revert-fehlt "p11-revert" p11 '^(Merge |Revert )' '^(Merge )'
mut m-p11-ausnahme-ohne-anker "p11-exempt-nicht-am-anfang" p11 '^(Merge |Revert )' '(Merge |Revert )'
# Punkt 12
mut m-p12-block-fehlt "p12-kern p12-ersetzt-anders p12-geloescht" alle "$BLOCK_VCS" ''
mut m-p12-paths-leer "p12-kern p12-ersetzt-anders p12-geloescht" p12 "$V_PATHS" 'paths: []'
# Der Index wird unveränderlich; head-allow prüft jede unveränderliche Datei im Kopf und lehnt eine ohne Status-Zeile ab, deshalb sind alle acht Fälle rot.
mut m-p12-paths-index "p12-kern p12-geschichte p12-ersetzt p12-ersetzt-anders p12-vorschlag p12-geloescht p12-link-ziel p12-index" p12 "$V_PATHS" 'paths: ["docs/plan/adr/*.md"]' "$V_ACC" "immutable-when: '^(\\*\\*Status:\\*\\* Accepted|# ADR-Index)'"
mut m-p12-immutable-when-jeder-status "p12-vorschlag" p12 "$V_ACC" "immutable-when: '^\\*\\*Status:\\*\\*'"
mut m-p12-immutable-when-anderer-status "p12-kern p12-ersetzt-anders p12-geloescht" p12 "$V_ACC" "immutable-when: '^\\*\\*Status:\\*\\* Rejected'"
mut m-p12-abschnitt-fehlt "p12-geschichte" p12 "$V_EXCL" '  status-line'
mut m-p12-abschnitt-falsch "p12-geschichte" p12 "$V_EXCL" $'  exclude-sections: [Historie]\n  status-line'
mut m-p12-status-zeile-eng "p12-ersetzt p12-ersetzt-anders" p12 "$V_STAT" "status-line: '^\\*\\*Status:\\*\\* Accepted'"
mut m-p12-head-allow-fehlt "p12-ersetzt-anders" p12 "  $V_HEAD"$'\n' ''
mut m-p12-head-allow-eng "p12-ersetzt" p12 "$V_HEAD" "head-allow: '^\\*\\*Status:\\*\\* Accepted'"
mut m-p12-head-allow-ohne-accepted "p12-geschichte p12-index p12-link-ziel p12-vorschlag" p12 "$V_HEAD" "head-allow: '^\\*\\*Status:\\*\\* Superseded by '"
mut m-p12-link-ziel-fehlt "p12-link-ziel" p12 "$V_IGN" ''
mut m-p12-link-ziel-aus "p12-link-ziel" p12 "$V_IGN" $'  ignore-link-targets: false\n'

wait
fehler=0
if [ -s "$arbeit/fehler" ]; then cat "$arbeit/fehler" >&2; fehler=1; fi
if [ "$(wc -l < "$arbeit/gefahren")" -eq 0 ]; then echo "range-konfiguration-gegenprobe: keine Mutation und keine Kopie gefahren (GP_MUT passt auf nichts)" >&2; fehler=1; fi
if [ "$fehler" -eq 0 ]; then echo "range-konfiguration-gegenprobe: grün ($(wc -l < "$arbeit/gefahren") Läufe)"; fi
exit "$fehler"
