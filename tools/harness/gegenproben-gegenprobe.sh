#!/usr/bin/env bash
# gegenproben-gegenprobe — prüft die Sammelziele (harness/mk/gegenproben.mk, tools/harness/gegenproben.sh)
# gegen Teil 2 von SPEC-060 (harness/vertraege/SPEC-060-gegenproben.md, Punkte 4 bis 12) an Kopien in einem
# Temp-Verzeichnis; der Arbeitsbaum bleibt unberührt, Docker wird nicht gerufen. Eine Kopie trägt Makefile, Fragment,
# Skript und ein Fragment mit Attrappen als Einzelziele (za zb zc in der Liste für den Container, ha hb in der für den
# Host, g1 als Gate); die unveränderte Kopie muss jeden Fall bestehen, und je Mutation (Textersatz in der Kopie)
# müssen genau die benannten Fälle rot werden. Der Name jedes Falls trägt die Nummer des Punktes.
# Exit-Codes (0, 1, 2, 130, 143) sind die des Skripts: make gibt jedem Fehler eines Rezepts Exit 2 und jedem -i Exit 0;
# darum fahren die Fälle das Skript dort unmittelbar (mit der Attrappe als Aufruf von make und MAKEFLAGS aus dem Fall) und
# prüfen über make, was make zeigt (Ziel, Reihenfolge, Ausgabe, kein Lauf, Exit ungleich 0).
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Mutationen, deren Name passt (die unveränderte Kopie nur, wenn das Muster
# `ok` trifft); GP_ENTDECKE=1 nennt je Mutation die roten Fälle; GP_PARALLEL=<n> (Vorgabe 4) setzt die Zahl der
# gleichzeitigen Mutationen. Braucht bash ≥ 5 (das Skript selbst läuft ab bash 4), make, setsid, pgrep, timeout.
# GRENZE: Die Einzelziele sind Attrappen; ob die echten Gegenproben grün sind, sagt `make gegenproben`. Ein Signal nur an
# make (Punkt 15) hat keinen Fall.
set -uo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
wurzel="$(cd "$hier/../.." && pwd)"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT
: > "$arbeit/fehler"
melde() { echo "gegenproben-gegenprobe: ROT — $*" >> "$arbeit/fehler"; }

# fixture <name> — legt eine Kopie an und setzt w (Wurzel) und st (Ablage der Läufe).
fixture() {
  w="$arbeit/$1"
  st="$arbeit/st-$1"
  mkdir -p "$st" "$w/harness/mk" "$w/tools/harness" "$w/.harness/state"
  cp "$wurzel/Makefile" "$w/Makefile"
  cp "$wurzel/harness/mk/gegenproben.mk" "$wurzel/harness/mk/enforce.mk" "$w/harness/mk/"
  cp "$hier/gegenproben.sh" "$hier/working-tree-hash.sh" "$hier/dateiliste.sh" "$hier/record-gates.sh" "$w/tools/harness/"
  cat > "$w/tools/harness/attrappe.sh" <<'EOF'
#!/usr/bin/env bash
z="$1"
echo "beginn $z" >> "$LOG"
echo "aus $z"
echo "fehl $z" >&2
[ -z "${LIES:-}" ] || echo "eingabe $z $(head -c 1 | wc -c)" >> "$LOG"
if [ "${BLOCK:-}" = "$z" ]; then
  echo $$ > "$ARBEIT/block.pid"
  if [ -n "${KIND:-}" ]; then (trap - INT TERM; exec sleep 31) & echo $! > "$ARBEIT/kind.pid"; fi
  exec sleep 20
fi
[ -z "${WARTE:-}" ] || sleep "$WARTE"
echo "ende $z" >> "$LOG"
for p in ${EXITE:-}; do [ "${p%%=*}" = "$z" ] && exit "${p#*=}"; done
exit 0
EOF
  chmod +x "$w/tools/harness/attrappe.sh"
  cat > "$w/harness/mk/fake.mk" <<'EOF'
FAKE := za zb zc ha hb g1
.PHONY: $(FAKE)
$(FAKE):
	@bash tools/harness/attrappe.sh $@
ifneq ($(LEER),container)
GEGENPROBEN_CONTAINER += za zb zc
endif
ifneq ($(LEER),host)
GEGENPROBEN_HOST += ha hb
endif
GATE_CHECKS += g1
GATE_SCHNELL_AUSNAHME := g1
EOF
  mf=""; warte=""; leer=""; exite=""; lies=""
}

# lauf <make-Argument>… — make in der Kopie; setzt mcode; Ausgabe in $st/out und $st/err (ohne Zeilen von make), Protokoll $st/log.
lauf() {
  : > "$st/log"
  rm -f "$st/block.pid"
  mcode=0
  (cd "$w" && env -u MAKELEVEL -u MFLAGS MAKEFLAGS="$mf" LOG="$st/log" ARBEIT="$st" WARTE="$warte" EXITE="$exite" \
    LEER="$leer" LIES="$lies" BLOCK="" timeout 60 make -s "$@" > "$st/out" 2> "$st/err.raw") || mcode=$?
  grep -v '^make' "$st/err.raw" > "$st/err" || true
}

# direkt <sammelziel> <MAKEFLAGS> — das Skript unmittelbar, die Attrappe als Aufruf von make; die Liste liest make aus der Variable.
direkt() {
  local var=GEGENPROBEN_CONTAINER liste
  if [ "$1" = gegenproben-host ]; then var=GEGENPROBEN_HOST; fi
  liste="$(cd "$w" && env -u MAKELEVEL -u MFLAGS MAKEFLAGS= LEER="$leer" make -s --eval "gp-lesen: ; @echo \$($var)" gp-lesen)"
  : > "$st/log"
  scode=0
  # shellcheck disable=SC2086
  (cd "$w" && env -u MAKELEVEL -u MFLAGS MAKEFLAGS="$2" LOG="$st/log" ARBEIT="$st" WARTE="$warte" EXITE="$exite" \
    LEER="$leer" LIES="$lies" BLOCK="" timeout 60 bash tools/harness/gegenproben.sh "$w/tools/harness/attrappe.sh" "$1" $liste \
    > "$st/out" 2> "$st/err") || scode=$?
}

# summe <sammelziel> — die Zeilen der Zusammenfassung in einer Zeile; erw <Zeile>… — die erwartete.
summe() { grep "^$1: " "$st/err" | tr '\n' '|' || true; }
erw() { local p s=""; for p in "$@"; do s+="gegenproben: $p|"; done; echo "$s"; }
# reihenfolge — die Zeilen des Protokolls der Einzelziele za bis zc in einer Zeile.
reihenfolge() { grep -E ' z[abc]$' "$st/log" | tr '\n' ' ' || true; }

warte_auf() { local i; for i in $(seq 1 100); do if "$@"; then return 0; fi; sleep 0.1; done; return 1; }

# abbruchfall <name> <Signal> <gruppe|skript> <Exit> — Punkt 8: das mittlere von drei Zielen wartet.
abbruchfall() {
  local name="$1" sig="$2" wer="$3" soll="$4" bg pid rc bp kp ok=1
  : > "$st/log"
  rm -f "$st/block.pid" "$st/pid" "$st/kind.pid"
  (trap - INT TERM; cd "$w" && exec env -u MAKELEVEL -u MFLAGS MAKEFLAGS= LOG="$st/log" ARBEIT="$st" BLOCK=zb KIND=1 WARTE= EXITE= LEER= \
    setsid bash -c 'echo $$ > "$ARBEIT/pid"; exec bash tools/harness/gegenproben.sh "$1" gegenproben za zb zc' _ \
    "$w/tools/harness/attrappe.sh" > "$st/out" 2> "$st/err") &
  bg=$!
  warte_auf test -s "$st/block.pid" || ok=0
  pid="$(cat "$st/pid" 2> /dev/null || echo 0)"
  if [ "$pid" -gt 1 ]; then
    if [ "$wer" = gruppe ]; then kill -s "$sig" -- "-$pid" 2> /dev/null; else kill -s "$sig" "$pid" 2> /dev/null; fi
    warte_auf bash -c "! kill -0 $pid 2> /dev/null" || { ok=0; kill -KILL -- "-$pid" 2> /dev/null; }
  else
    ok=0
  fi
  rc=0; wait "$bg" || rc=$?
  [ "$rc" -eq "$soll" ] || ok=0
  bp="$(cat "$st/block.pid" 2> /dev/null || echo 0)"
  if [ "$bp" -gt 1 ] && kill -0 "$bp" 2> /dev/null; then ok=0; kill -KILL "$bp" 2> /dev/null; fi
  kp="$(cat "$st/kind.pid" 2> /dev/null || echo 0)"
  if [ "$kp" -gt 1 ] && kill -0 "$kp" 2> /dev/null; then ok=0; kill -KILL "$kp" 2> /dev/null; fi
  if grep -qx 'beginn zc' "$st/log"; then ok=0; fi
  if grep -q '^gegenproben: ' "$st/err"; then ok=0; fi
  if [ "$pid" -gt 1 ]; then kill -KILL -- "-$pid" 2> /dev/null || true; fi
  [ "$ok" -eq 1 ] || echo "$name"
}

# sechs <name> <EXITE> <Exit> <Zeile>… — Punkt 6 am Skript: Exit, Zusammenfassung (genau, zuletzt), alle Ziele liefen.
sechs() {
  local name="$1" soll="$3"
  exite="$2"; shift 3
  direkt gegenproben ""
  { [ "$scode" -eq "$soll" ] && [ "$(summe gegenproben)" = "$(erw "$@")" ] && [ "$(tail -n 3 "$st/err" | tr '\n' '|')" = "$(erw "$@")" ] \
    && grep -qx 'ende zc' "$st/log"; } || echo "$name"
  exite=""
}

# flagfall <name> <MAKEFLAGS der Umgebung> <MAKEFLAGS am Skript> <trocken|i> <make-Argument>… — Punkt 10.
flagfall() {
  local name="$1" menv="$2" dmf="$3" art="$4" ok=1
  shift 4
  mf="$menv"; lauf "$@" gegenproben; mf=""
  if grep -q ' z' "$st/log"; then ok=0; fi
  if [ "$art" = trocken ]; then
    if grep -q '^gegenproben' "$st/err"; then ok=0; fi
  else
    grep -q '^gegenproben: -i' "$st/err" || ok=0
    [ "$mcode" -eq 0 ] || ok=0
  fi
  direkt gegenproben "$dmf"
  if grep -q ' z' "$st/log"; then ok=0; fi
  if [ "$art" = trocken ]; then
    { [ "$scode" -eq 0 ] && [ ! -s "$st/err" ]; } || ok=0
  else
    { [ "$scode" -eq 2 ] && grep -q '^gegenproben: -i' "$st/err"; } || ok=0
  fi
  [ "$ok" -eq 1 ] || echo "$name"
}

# stand — Pfade und Inhalte des Arbeitsbaums der Kopie.
stand() { (cd "$w" && find . -type f -print0 | sort -z | xargs -0 sha256sum); }

# standfall <name> <gruen|rot|abbruch|n|q|t> — Punkt 11: der Lauf legt im Arbeitsbaum der Kopie nichts an und ändert nichts.
standfall() {
  local vor
  rm -rf "$w/.harness"; mkdir -p "$w/.harness/state"
  vor="$(stand)"
  case "$2" in
    gruen) direkt gegenproben "" ;;
    rot) exite="zb=1"; direkt gegenproben ""; exite="" ;;
    abbruch) abbruchfall "$1" INT skript 130 > /dev/null ;;
    *) direkt gegenproben "$2" ;;
  esac
  [ "$(stand)" = "$vor" ] || echo "$1"
}

# pruefe — schreibt die Namen der roten Fälle, einen je Zeile.
pruefe() {
  local gw x vor ok

  # Punkt 4
  lauf gegenproben
  { grep -qx 'beginn za' "$st/log" && grep -qx 'ende zc' "$st/log" && ! grep -q ' h' "$st/log"; } || echo p4-container
  lauf gegenproben-host
  { grep -qx 'ende hb' "$st/log" && ! grep -q ' z' "$st/log"; } || echo p4-host
  lauf -n gates
  gw="$(cd "$w" && env -u MAKELEVEL -u MFLAGS MAKEFLAGS= make -s --eval 'gp-lesen: ; @echo $(GATE_CHECKS)' gp-lesen)"
  ok=1
  if grep -q 'gegenproben' "$st/out"; then ok=0; fi
  for x in $gw; do case "$x" in gegenproben | gegenproben-host) ok=0 ;; esac; done
  case " $gw " in *" g1 "*) ;; *) ok=0 ;; esac
  [ "$ok" -eq 1 ] || echo p4-kein-gate
  gw="$(cd "$w" && env -u MAKELEVEL -u MFLAGS MAKEFLAGS= make -s --eval 'gp-lesen: ; @echo $(GATE_SCHNELL_AUSNAHME)' gp-lesen)"
  ok=1
  for x in $gw; do case "$x" in gegenproben | gegenproben-host) ok=0 ;; esac; done
  case " $gw " in *" g1 "*) ;; *) ok=0 ;; esac
  [ "$ok" -eq 1 ] || echo p4-kein-ausnahme

  # Punkt 5
  warte=0.15
  lauf gegenproben
  [ "$(reihenfolge)" = "beginn za ende za beginn zb ende zb beginn zc ende zc " ] || echo p5-reihenfolge
  lauf -j4 gegenproben
  [ "$(reihenfolge)" = "beginn za ende za beginn zb ende zb beginn zc ende zc " ] || echo p5-j
  warte=""

  # Punkt 6
  sechs p6-alle-gruen "" 0 "za: grün" "zb: grün" "zc: grün"
  lauf gegenproben
  [ "$mcode" -eq 0 ] || echo p6-alle-gruen
  sechs p6-erstes-rot "za=1" 1 "za: rot (Exit 1)" "zb: grün" "zc: grün"
  exite="za=1"; lauf gegenproben; exite=""
  [ "$mcode" -ne 0 ] || echo p6-erstes-rot
  sechs p6-mittleres-rot "zb=1" 1 "za: grün" "zb: rot (Exit 1)" "zc: grün"
  sechs p6-letztes-rot "zc=1" 1 "za: grün" "zb: grün" "zc: rot (Exit 1)"
  sechs p6-zwei-rot "za=3 zc=5" 1 "za: rot (Exit 3)" "zb: grün" "zc: rot (Exit 5)"
  sechs p6-127 "zb=127" 1 "za: grün" "zb: rot (Exit 127)" "zc: grün"

  # Punkt 7
  direkt gegenproben ""
  { grep -qx 'aus za' "$st/out" && ! grep -q 'fehl' "$st/out" && grep -qx 'fehl za' "$st/err" && ! grep -q '^aus' "$st/err"; } || echo p7-ausgabe

  lies=1; direkt gegenproben "" <<< x; lies=""
  { grep -qx 'eingabe za 0' "$st/log" && grep -qx 'eingabe zb 0' "$st/log" && grep -qx 'eingabe zc 0' "$st/log"; } || echo p7-eingabe

  # Punkt 8
  abbruchfall p8-sigint-gruppe INT gruppe 130
  abbruchfall p8-sigterm-gruppe TERM gruppe 143
  abbruchfall p8-sigint-skript INT skript 130
  abbruchfall p8-sigterm-skript TERM skript 143

  # Punkt 9
  leer=container; lauf gegenproben
  { [ "$mcode" -ne 0 ] && grep -qx 'gegenproben: Liste leer' "$st/err" && ! grep -q ' z' "$st/log"; } || echo p9-container
  direkt gegenproben ""
  { [ "$scode" -eq 2 ] && grep -qx 'gegenproben: Liste leer' "$st/err"; } || echo p9-container
  leer=host; lauf gegenproben-host
  { [ "$mcode" -ne 0 ] && grep -qx 'gegenproben-host: Liste leer' "$st/err" && ! grep -q ' h' "$st/log"; } || echo p9-host
  direkt gegenproben-host ""
  { [ "$scode" -eq 2 ] && grep -qx 'gegenproben-host: Liste leer' "$st/err"; } || echo p9-host
  leer=""

  # Punkt 10
  flagfall p10-n "" n trocken -n
  flagfall p10-q "" q trocken -q
  flagfall p10-t "" t trocken -t
  flagfall p10-n-env n n trocken
  flagfall p10-q-env q q trocken
  flagfall p10-t-env t t trocken
  flagfall p10-i "" i i -i
  flagfall p10-i-env i i i
  flagfall p10-i-lang "" i i --ignore-errors
  flagfall p10-i-k "" ik i -i -k
  flagfall p10-i-j "" "i -j2" i -i -j2
  flagfall p10-in "" in trocken -i -n
  flagfall p10-iq "" iq trocken -i -q
  flagfall p10-it "" it trocken -i -t
  exite="zb=1"; lauf -k gegenproben; exite=""
  { [ "$mcode" -ne 0 ] && grep -qx 'ende zc' "$st/log" && [ "$(summe gegenproben)" = "$(erw 'za: grün' 'zb: rot (Exit 2)' 'zc: grün')" ]; } || echo p10-k
  exite="zb=1"; warte=0.1; lauf -j4 gegenproben; exite=""; warte=""
  { [ "$mcode" -ne 0 ] && [ "$(reihenfolge)" = "beginn za ende za beginn zb ende zb beginn zc ende zc " ] \
    && [ "$(summe gegenproben)" = "$(erw 'za: grün' 'zb: rot (Exit 2)' 'zc: grün')" ]; } || echo p10-j

  # Punkt 11
  for x in gruen rot abbruch n q t; do standfall "p11-$x" "$x"; done

  # Punkt 12
  leer=container
  direkt gegenproben n
  { [ "$scode" -eq 0 ] && [ ! -s "$st/err" ]; } || echo p12-n-leer
  direkt gegenproben i
  { [ "$scode" -eq 2 ] && grep -q '^gegenproben: -i' "$st/err" && ! grep -q 'Liste leer' "$st/err"; } || echo p12-i-leer
  leer=""
}

# ergebnis — die roten Fälle sortiert in einer Zeile.
ergebnis() { pruefe | LC_ALL=C sort -u | tr '\n' ' '; }

maxjobs="${GP_PARALLEL:-4}"
wartejob() { while [ "$(jobs -rp | wc -l)" -ge "$maxjobs" ]; do wait -n || true; done; }

# mut <name> <datei> <rote Fälle> <alt> <neu> [<alt> <neu>…] — Textersatz in der Kopie; rot sind genau die benannten Fälle.
mut() {
  local name="$1" datei="$2" rote="$3" inhalt neu ist erwartet
  shift 3
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  (
    fixture "$name"
    inhalt="$(cat "$w/$datei"; printf x)"; inhalt="${inhalt%x}"
    neu="$inhalt"
    while [ "$#" -ge 2 ]; do
      if [[ $neu != *"$1"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $1"; exit 0; fi
      neu="${neu/"$1"/"$2"}"
      shift 2
    done
    printf '%s' "$neu" > "$w/$datei"
    ist="$(ergebnis)"
    erwartet="$(printf '%s\n' $rote | LC_ALL=C sort -u | tr '\n' ' ')"
    if [ -n "${GP_ENTDECKE-}" ]; then echo "$name => ${ist:-nichts}" >&2; fi
    [ "$ist" = "$erwartet" ] || melde "Mutation '$name': erwartet rot: ${rote:-nichts}; rot sind: ${ist:-nichts}"
  ) &
}

if [[ ok =~ ${GP_MUT:-.} ]]; then
  (
    fixture ok
    ist="$(ergebnis)"
    [ -z "$ist" ] || melde "unveränderte Kopie ist nicht grün: $ist"
  ) &
fi

S=tools/harness/gegenproben.sh
M=harness/mk/gegenproben.mk
NL=$'\n'
# --- Punkt 4: Ziele, Listen, kein Gate ---------------------------------------------------------------
mut m-p4-container-mit-host $M "p4-container p9-container p10-j p10-k" 'gegenproben $(GEGENPROBEN_CONTAINER)' 'gegenproben $(GEGENPROBEN_CONTAINER) $(GEGENPROBEN_HOST)'
mut m-p4-host-mit-container $M "p4-host p9-host" 'gegenproben-host $(GEGENPROBEN_HOST)' 'gegenproben-host $(GEGENPROBEN_HOST) $(GEGENPROBEN_CONTAINER)'
mut m-p4-container-haengt-an-host $M "p4-container" 'gegenproben: ## Fährt die Gegenproben der Liste für den Dev' 'gegenproben: gegenproben-host ## Fährt die Gegenproben der Liste für den Dev'
mut m-p4-host-haengt-an-container $M "p4-host" 'gegenproben-host: ## ' 'gegenproben-host: gegenproben ## '
mut m-p4-in-gate-checks $M "p4-kein-gate" '.PHONY: gegenproben gegenproben-host' ".PHONY: gegenproben gegenproben-host${NL}GATE_CHECKS += gegenproben"
mut m-p4-host-in-gate-checks $M "p4-kein-gate" '.PHONY: gegenproben gegenproben-host' ".PHONY: gegenproben gegenproben-host${NL}GATE_CHECKS += gegenproben-host"
mut m-p4-ausnahme-container $M "p4-kein-ausnahme" '.PHONY: gegenproben gegenproben-host' ".PHONY: gegenproben gegenproben-host${NL}GATE_SCHNELL_AUSNAHME += gegenproben"
mut m-p4-ausnahme-host $M "p4-kein-ausnahme" '.PHONY: gegenproben gegenproben-host' ".PHONY: gegenproben gegenproben-host${NL}GATE_SCHNELL_AUSNAHME += gegenproben-host"
# --- Punkt 5: nacheinander, in der Reihenfolge der Liste ---------------------------------------------
mut m-p5-hintergrund $S "p5-reihenfolge p5-j p6-erstes-rot p6-mittleres-rot p6-letztes-rot p6-zwei-rot p6-127 p8-sigint-gruppe p8-sigterm-gruppe p8-sigint-skript p8-sigterm-skript p10-j p10-k" '  wait "$pid"; rc=$?' '  sleep 0.01; rc=0'
mut m-p5-rueckwaerts $S "p5-reihenfolge p5-j p6-alle-gruen p6-erstes-rot p6-mittleres-rot p6-letztes-rot p6-zwei-rot p6-127 p8-sigint-gruppe p8-sigterm-gruppe p8-sigint-skript p8-sigterm-skript p10-j p10-k" 'for z in "$@"; do' 'for z in $(printf "%s\n" "$@" | tac); do'
mut m-p5-prerequisites $M "p5-j p9-container p10-i p10-i-env p10-i-lang p10-i-k p10-i-j p10-j p10-k" 'gegenproben: ## Fährt' 'gegenproben: $(GEGENPROBEN_CONTAINER) ## Fährt' "	@bash tools/harness/gegenproben.sh '\$(MAKE)' gegenproben \$(GEGENPROBEN_CONTAINER)" '	@:'
# --- Punkt 6: rotes Einzelziel, Ausgang, Zusammenfassung ---------------------------------------------
mut m-p6-abbruch-beim-ersten $S "p6-erstes-rot p6-mittleres-rot p6-zwei-rot p6-127 p10-j p10-k" 'rot=1; fi' 'rot=1; break; fi'
mut m-p6-exit0 $S "p6-erstes-rot p6-mittleres-rot p6-letztes-rot p6-zwei-rot p6-127 p10-j p10-k" 'exit "$rot"' 'exit 0'
mut m-p6-exit-des-letzten $S "p6-erstes-rot p6-mittleres-rot p6-127 p10-j p10-k" 'grün"); else' 'grün"); rot=0; else'
mut m-p6-nur-rote $S "p6-alle-gruen p6-erstes-rot p6-mittleres-rot p6-letztes-rot p6-zwei-rot p6-127 p10-j p10-k" 'zeilen+=("$sammel: $z: grün"); ' ':; '
mut m-p6-ohne-exit $S "p6-erstes-rot p6-mittleres-rot p6-letztes-rot p6-zwei-rot p6-127 p10-j p10-k" ' (Exit $rc)' ''
mut m-p6-zusammenfassung-vor-ausgabe $S "p6-alle-gruen p6-erstes-rot p6-mittleres-rot p6-letztes-rot p6-zwei-rot p6-127 p10-j p10-k" '  ( exec' "  if [ \"\$z\" = \"\${!#}\" ]; then printf '%s\\n' \"\${zeilen[@]}\" >&2; fi${NL}  ( exec"
# --- Punkt 7: Ausgabe unverändert, auf demselben Strom ------------------------------------------------
mut m-p7-stderr-nach-stdout $S "p7-ausgabe" 'exec "$make_bin" "$z" )' 'exec "$make_bin" "$z" 2>&1 )'
mut m-p7-stdout-verworfen $S "p7-ausgabe" 'exec "$make_bin" "$z" )' 'exec "$make_bin" "$z" > /dev/null )'
# --- Punkt 8: Abbruch ----------------------------------------------------------------------------------
mut m-p8-naechstes-startet $S "p8-sigint-gruppe p8-sigterm-gruppe p8-sigint-skript p8-sigterm-skript" 'exit "$2"; }' 'return 0; }'
mut m-p8-vertauscht $S "p8-sigint-gruppe p8-sigterm-gruppe p8-sigint-skript p8-sigterm-skript" "abbruch INT 130" "abbruch INT 143" "abbruch TERM 143" "abbruch TERM 130"
mut m-p8-ziel-bleibt $S "p8-sigint-skript p8-sigterm-skript" 'if [ -n "$pid" ]; then beende "$pid" "$1"; wait "$pid" 2> /dev/null; fi; ' ''
mut m-p8-zusammenfassung $S "p8-sigint-gruppe p8-sigterm-gruppe p8-sigint-skript p8-sigterm-skript" "trap '' INT TERM; if" "trap '' INT TERM; printf '%s\\n' \"\${zeilen[@]}\" >&2; if"
mut m-p7-eingabe-zero $S "p7-eingabe" 'exec "$make_bin" "$z" )' 'exec "$make_bin" "$z" < /dev/zero )'
mut m-p8-int-nach-143 $S "p8-sigint-gruppe p8-sigint-skript" "abbruch INT 130" "abbruch INT 143"
mut m-p8-term-nach-130 $S "p8-sigterm-gruppe p8-sigterm-skript" "abbruch TERM 143" "abbruch TERM 130"
mut m-p8-trap-int-gestrichen $S "p8-sigint-skript" "trap 'abbruch INT 130' INT${NL}" ""
mut m-p8-trap-term-gestrichen $S "p8-sigterm-skript" "trap 'abbruch TERM 143' TERM${NL}" ""
mut m-p8-nachkommen-bleiben $S "p8-sigint-skript p8-sigterm-skript" 'for k in $(pgrep -P "$1" 2> /dev/null); do beende "$k" "$2"; done; ' ''
# --- Punkt 9: leere Liste ------------------------------------------------------------------------------
mut m-p9-exit0 $S "p9-container p9-host" "Liste leer\" >&2${NL}  exit 2" "Liste leer\" >&2${NL}  exit 0"
mut m-p9-liste-des-anderen $M "p4-container p4-host p5-reihenfolge p5-j p6-erstes-rot p9-container p9-host p10-j p10-k" 'gegenproben $(GEGENPROBEN_CONTAINER)' 'gegenproben $(GEGENPROBEN_HOST)' 'gegenproben-host $(GEGENPROBEN_HOST)' 'gegenproben-host $(GEGENPROBEN_CONTAINER)'
# --- Punkt 10: Aufrufarten von make ---------------------------------------------------------------------
mut m-p10-n-uebersehen $S "p10-n p10-n-env p10-in p12-n-leer" 'case "$buchstaben" in *n*) trocken=1 ;; esac' ''
mut m-p10-q-uebersehen $S "p10-q p10-q-env p10-iq" 'case "$buchstaben" in *q*) trocken=1 ;; esac' ''
mut m-p10-t-uebersehen $S "p10-t p10-t-env p10-it" 'case "$buchstaben" in *t*) trocken=1 ;; esac' ''
mut m-p10-i-faellt-aus $S "p10-i p10-i-env p10-i-lang p10-i-k p10-i-j p12-i-leer" 'case "$buchstaben" in *i*) ueberlesen=1 ;; esac' ''
mut m-p10-i-wie-n $S "p10-i p10-i-env p10-i-lang p10-i-k p10-i-j p12-i-leer" '*i*) ueberlesen=1' '*i*) trocken=1'
mut m-p10-i-vor-n $S "p10-in p10-iq p10-it" '[ "$trocken" -eq 0 ] || exit 0' '[ "$trocken" -eq 0 ] || [ "$ueberlesen" -eq 1 ] || exit 0'
mut m-p10-ganzes-makeflags $S "p5-j p10-i-j p10-j" 'erstes="${flags%% *}"' 'erstes="$flags"'
mut m-p10-ohne-k $S "p6-erstes-rot p6-mittleres-rot p6-zwei-rot p6-127 p10-j" 'rot=1; fi' 'rot=1; case "$buchstaben" in *k*) ;; *) break ;; esac; fi'
mut m-p10-j-gleichzeitig $S "p5-j p10-j" '  wait "$pid"; rc=$?' '  case "$flags" in *-j*) rc=0 ;; *) wait "$pid"; rc=$? ;; esac'
# --- Punkt 11: kein Stand ---------------------------------------------------------------------------------
SCHREIBE='mkdir -p .harness/state; echo x > .harness/state/gegenproben-passed'
mut m-p11-gruen-schreibt $S "p11-gruen" 'exit "$rot"' "if [ \"\$rot\" -eq 0 ]; then ${SCHREIBE}; fi${NL}exit \"\$rot\""
mut m-p11-rot-schreibt $S "p11-rot" 'exit "$rot"' "if [ \"\$rot\" -eq 1 ]; then ${SCHREIBE}; fi${NL}exit \"\$rot\""
mut m-p11-abbruch-schreibt $S "p11-abbruch" "abbruch() { trap '' INT TERM; " "abbruch() { trap '' INT TERM; ${SCHREIBE}; "
mut m-p11-n-schreibt $S "p11-n" '*n*) trocken=1 ;;' "*n*) trocken=1; ${SCHREIBE} ;;"
mut m-p11-q-schreibt $S "p11-q" '*q*) trocken=1 ;;' "*q*) trocken=1; ${SCHREIBE} ;;"
mut m-p11-t-schreibt $S "p11-t" '*t*) trocken=1 ;;' "*t*) trocken=1; ${SCHREIBE} ;;"
# --- Punkt 12: Reihenfolge der Gründe ---------------------------------------------------------------------
LEERPRUEFUNG="if [ \"\$#\" -eq 0 ]; then echo \"\$sammel: Liste leer\" >&2; exit 2; fi"
mut m-p12-liste-vor-n $S "p12-n-leer p12-i-leer" "shift 2${NL}" "shift 2${NL}${LEERPRUEFUNG}${NL}"
mut m-p12-liste-vor-i $S "p12-i-leer" '[ "$trocken" -eq 0 ] || exit 0' "[ \"\$trocken\" -eq 0 ] || exit 0${NL}${LEERPRUEFUNG}"

wait
if [ -s "$arbeit/fehler" ]; then
  cat "$arbeit/fehler" >&2
  exit 1
fi
echo "gegenproben-gegenprobe: gruen — Punkte 4 bis 12 von SPEC-060 mit 44 Fällen, je Mutation genau die benannten Fälle rot, unveränderte Kopie grün"
