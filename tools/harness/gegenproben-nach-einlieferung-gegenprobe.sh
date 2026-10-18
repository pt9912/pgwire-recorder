#!/usr/bin/env bash
# gegenproben-nach-einlieferung-gegenprobe — prüft tools/harness/gegenproben-nach-einlieferung.sh gegen die Punkte 1 bis 13 von SPEC-061
# (harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md). Je Gruppe von Fällen entsteht ein Aufbau in einem Temp-Verzeichnis: ein echtes Bare-Repo mit
# Marke und der Kopie des Skripts als hooks/post-receive, ein Klon mit drei Commits und eine Attrappe als `make` (im PATH vorangestellt; sie protokolliert
# Aufruf, Umgebung und Arbeitsbaum). Die unveränderte Kopie muss jeden Fall bestehen; je Mutation (Textersatz in der Kopie) müssen genau die benannten Fälle rot
# werden. Der Name jedes Falls trägt die Nummer des Punktes. Docker wird nicht gerufen; der Arbeitsbaum bleibt unberührt.
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Mutationen, deren Name passt (die unveränderte Kopie nur, wenn das Muster `ok` trifft); GP_ENTDECKE=1 nennt je
# Mutation die roten Fälle; GP_PARALLEL=<n> (Vorgabe 4) setzt die Zahl der gleichzeitigen Mutationen. Braucht bash ≥ 5, git, flock, setsid, pgrep, timeout.
# Der Fall p13-ablage-weg verlangt einen Benutzer ohne root (chmod sperrt root nicht aus); als root endet die Gegenprobe mit Exit 2.
# GRENZE: Die Attrappe ersetzt `make gegenproben`; ob die echten Gegenproben grün sind, sagt die erste Fahrt im Container. SIGKILL, ein Neustart und die
# Gleichzeitigkeit zweier Aufrufe mit gleichem Namen (Punkte 14, 16) haben keinen Fall.
set -uo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
skript="$hier/gegenproben-nach-einlieferung.sh"
if [ "$(id -u)" -eq 0 ]; then echo "gegenproben-nach-einlieferung-gegenprobe: nicht als root (p13-ablage-weg)" >&2; exit 2; fi
arbeit="$(mktemp -d)"
trap 'pkill -KILL -f -- "$arbeit" 2> /dev/null; chmod -R u+rwX "${arbeit:?}"; rm -rf "${arbeit:?}"' EXIT
: > "$arbeit/fehler"
melde() { echo "gegenproben-nach-einlieferung-gegenprobe: ROT — $*" >> "$arbeit/fehler"; }
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null
PFAD="$PATH"
Z0=0000000000000000000000000000000000000000
Z64="$(printf '0%.0s' {1..64})"

# Die Attrappe von make: legt je Aufruf ein Verzeichnis aufruf.<n> in steuer/ an (argv, env, cwd, Arbeitsbaum) und liest ihre Steuerung aus steuer/.
cat > "$arbeit/attrappe" <<'EOT'
#!/usr/bin/env bash
st="$(cd "$(dirname "$0")/../steuer" && pwd)"
n=$(($(cat "$st/zaehler" 2> /dev/null || echo 0) + 1)); echo "$n" > "$st/zaehler"
d="$st/aufruf.$n"; mkdir -p "$d"
printf '%s\n' "$@" > "$d/argv"
env | LC_ALL=C sort > "$d/env"
pwd -P > "$d/cwd"; git rev-parse HEAD > "$d/head" 2>&1; git status --porcelain > "$d/status" 2>&1
git symbolic-ref -q HEAD > "$d/symref" 2>&1; echo $? > "$d/symrc"
echo "beginn $n $(date +%s.%N)" >> "$st/protokoll"
echo aus; echo fehl >&2
if [ -e "$st/block" ]; then (trap - TERM; exec sleep 31) & echo $! > "$st/kind.pid"; echo $$ > "$st/block.pid"; exec sleep 30; fi
[ ! -s "$st/warte" ] || sleep "$(cat "$st/warte")"
[ ! -e "$st/rot.$n" ] || cat "$st/rot.$n" >&2
echo "ende $n $(date +%s.%N)" >> "$st/protokoll"
exit "$(sed -n "${n}p" "$st/exits" 2> /dev/null | grep . || echo 0)"
EOT
chmod +x "$arbeit/attrappe"

warte_auf() { local i; for i in $(seq 1 "${WA:-100}"); do if "$@"; then return 0; fi; sleep 0.1; done; return 1; }
fertig() { ! pgrep -f -- "post-receive --lauf $g/bare" > /dev/null; }
bis_fertig() { sleep 0.3; warte_auf fertig; }
zl() { printf '%s %s refs/heads/%s\n' "$Z0" "$2" "$1"; }
zs() { while [ "$#" -ge 2 ]; do zl "$1" "$2"; shift 2; done; }
gest() { echo "gegenproben-nach-einlieferung: gestartet: $1 $2"; }
snap() { (cd "$1" && { find . | LC_ALL=C sort; find . -type f -exec sha256sum {} + | LC_ALL=C sort; }); }
feld() { sed -n "s/^$2=//p" "$1"; }
res() { echo "$g/ablage/gegenproben-$1-$2"; }
zurueck() { bis_fertig; leere; rm -rf "${g:?}"/steuer/aufruf.* "${g:?}/steuer/zaehler"; }
leere() { rm -f "${g:?}"/ablage/gegenproben* "${g:?}"/ablage/.gegenproben-* "${g:?}"/steuer/protokoll; }

# aufbau <gruppe> — setzt g (Wurzel der Gruppe) und h1 bis h3 (Hashes der drei Commits; main des Bare-Repos steht auf h3).
aufbau() {
  local i
  g="$fx/$1"; mkdir -p "$g/ablage" "$g/bin" "$g/steuer"
  git init -q --bare "$g/bare"; git init -q -b main "$g/arb"
  for i in 1 2 3; do
    git -C "$g/arb" -c user.name=t -c user.email=t@t commit -q --allow-empty -m "c$i"
    eval "h$i=\$(git -C \"\$g/arb\" rev-parse HEAD)"
  done
  git -C "$g/arb" push -q "$g/bare" main
  cp "$arbeit/attrappe" "$g/bin/make"
  echo "$g/ablage" > "$g/bare/gegenproben-ablage"
  install -m 755 "$hook" "$g/bare/hooks/post-receive"
  pfad=""; hskript=""
}

# haken <Eingabe> [VAR=Wert…] — ruft den Hook wie git (Arbeitsverzeichnis Bare-Repo, GIT_DIR=., leere Umgebung); setzt hrc; Ausgabe in $g/h.out und $g/h.err.
haken() {
  local z="$1"; shift
  hrc=0
  (cd "$g/bare" && printf '%s' "$z" | env -i PATH="${pfad:-$g/bin:$PFAD}" HOME="$g" GIT_DIR=. "$@" "$(command -v timeout)" 30 "$BASH" "${hskript:-hooks/post-receive}" \
    > "$g/h.out" 2> "$g/h.err") || hrc=$?
}
nichts() { [ "$hrc" -eq 0 ] && [ ! -s "$g/h.out" ] && [ -z "$(ls -A "$g/ablage")" ] && [ ! -e "$g/steuer/zaehler" ]; }
err_ist() { [ "$(cat "$g/h.err")" = "$1" ]; }

# --- Punkt 1 --------------------------------------------------------------------------------------------------------------
g1() {
  local prc
  aufbau p1
  mkdir "$g/kopie"; install -m 755 "$hook" "$g/kopie/post-receive"; hskript="$g/kopie/post-receive"
  haken "$(zl main "$h1")"; hskript=""
  { err_ist "$(gest main "$h1")" && [ "$hrc" -eq 0 ]; } || echo p1-kopie
  bis_fertig; [ -f "$(res main "$h1")" ] || echo p1-kopie
  echo 2 > "$g/steuer/warte"
  PATH="$g/bin:$PFAD" git -C "$g/arb" push -q "$g/bare" "$h2:refs/heads/neu" > /dev/null 2>&1; prc=$?
  { [ "$prc" -eq 0 ] && [ -e "$(res neu "$h2").laeuft" ] && [ ! -e "$(res neu "$h2")" ]; } || echo p1-echter-push
  bis_fertig; [ "$(feld "$(res neu "$h2")" hash)" = "$h2" ] || echo p1-echter-push
}

# --- Punkt 2 --------------------------------------------------------------------------------------------------------------
g2() {
  aufbau p2
  marke_fall() {
    local vor; vor="$(snap "$g/bare")"
    haken "$(zl main "$h1")"; sleep 0.3
    { nichts && grep -q '^gegenproben-nach-einlieferung: marke: ' "$g/h.err" && [ "$(snap "$g/bare")" = "$vor" ]; } || echo "$1"
    zurueck
  }
  rm -f "${g:?}/bare/gegenproben-ablage"; marke_fall p2-fehlt
  : > "$g/bare/gegenproben-ablage"; marke_fall p2-leer
  mkdir "$g/bare/rel"; echo rel > "$g/bare/gegenproben-ablage"; marke_fall p2-relativ
  printf '%s\nx\n' "$g/ablage" > "$g/bare/gegenproben-ablage"; marke_fall p2-zwei-zeilen
  echo "$g/nichtda" > "$g/bare/gegenproben-ablage"; marke_fall p2-kein-verzeichnis
}

# --- Punkt 3 --------------------------------------------------------------------------------------------------------------
g3() {
  local r
  aufbau p3
  for r in tags:p3-tag notes:p3-notes headsx:p3-headsx; do
    haken "$(printf '%s %s refs/%s/x\n' "$Z0" "$h1" "${r%%:*}")"; sleep 0.3
    { nichts && [ ! -s "$g/h.err" ]; } || echo "${r##*:}"
    zurueck
  done
  echo 0.3 > "$g/steuer/warte"
  haken "$(zs a "$h1" b "$h2" c "$h3")"; bis_fertig
  { [ -f "$(res a "$h1")" ] && [ -f "$(res b "$h2")" ] && [ -f "$(res c "$h3")" ] \
    && [ "$(awk '{print $1}' "$g/steuer/protokoll" | tr '\n' ' ')" = "beginn ende beginn ende beginn ende " ] \
    && [ "$(cat "$g/steuer/aufruf.1/head" "$g/steuer/aufruf.2/head" "$g/steuer/aufruf.3/head" | tr '\n' ' ')" = "$h1 $h2 $h3 " ]; } || echo p3-mehrere
}

# --- Punkt 4 --------------------------------------------------------------------------------------------------------------
g4() {
  local gross h64 v
  aufbau p4
  haken "$h1 $Z0 refs/heads/x"; sleep 0.3; { nichts && [ ! -s "$g/h.err" ]; } || echo p4-null40; zurueck
  haken "$h1 $Z64 refs/heads/x"; sleep 0.3; { nichts && [ ! -s "$g/h.err" ]; } || echo p4-null64; zurueck
  echo alt > "$(res b "$h1")"; echo log > "$(res b "$h1").log"; : > "$(res b "$h1").laeuft"; v="$(snap "$g/ablage")"
  haken "$h1 $Z0 refs/heads/b"; sleep 0.3
  { [ "$hrc" -eq 0 ] && [ ! -s "$g/h.err" ] && [ "$(snap "$g/ablage")" = "$v" ] && [ "$(cat "$(res b "$h1")")" = alt ]; } || echo p4-ergebnis-bleibt
  leere
  gross="$(printf 'A%.0s' {1..40})"; h64="$(printf 'ab%.0s' {1..32})"
  haken "$Z0 ${h1:0:39} refs/heads/x"; sleep 0.3; { nichts && grep -q 'zeile: ' "$g/h.err"; } || echo p4-hash-kurz; zurueck
  haken "$Z0 $gross refs/heads/x"; sleep 0.3; { nichts && grep -q 'zeile: ' "$g/h.err"; } || echo p4-hash-gross; zurueck
  haken "$Z0 ${h1}a refs/heads/x"; sleep 0.3; { nichts && grep -q 'zeile: ' "$g/h.err"; } || echo p4-hash-41; zurueck
  haken "$Z0 $h64 refs/heads/x"; bis_fertig
  { err_ist "$(gest x "$h64")" && [ -e "$(res x "$h64").laeuft" ]; } || echo p4-hash-64
  leere
  haken "$(printf 'x\n%s %s\n' "$Z0" "$h1")"; sleep 0.3
  { nichts && [ "$(grep -c 'zeile: ' "$g/h.err")" -eq 2 ]; } || echo p4-feld-fehlt
  haken ""; sleep 0.3; { nichts && [ ! -s "$g/h.err" ]; } || echo p4-stdin-leer; zurueck
  haken "$(printf '%s %s refs/heads/y\n' "$Z0" "${h1:0:39}"; zl x "$h1")"; bis_fertig
  { grep -q '^gegenproben-nach-einlieferung: zeile: ' "$g/h.err" && grep -q "gestartet: x $h1\$" "$g/h.err" && [ -f "$(res x "$h1")" ]; } || echo p4-zeile-kaputt-dann-gut
}

# --- Punkt 5 --------------------------------------------------------------------------------------------------------------
g5() {
  local pgid bgpid l
  aufbau p5
  echo 2 > "$g/steuer/warte"
  haken "$(zl main "$h1")"
  { [ ! -e "$(res main "$h1")" ] && ! grep -qs ende "$g/steuer/protokoll"; } || echo p5-vor-ende
  { [ "$hrc" -eq 0 ] && err_ist "$(gest main "$h1")"; } || echo p5-gestartet-zeile
  [ ! -s "$g/h.out" ] || echo p5-stdout-leer
  warte_auf grep -qs beginn "$g/steuer/protokoll"
  bgpid="$(pgrep -of "post-receive --lauf $g/bare" || echo 0)"
  for l in 0 1 2; do [ "$(readlink "/proc/$bgpid/fd/$l" 2> /dev/null)" = /dev/null ] || echo "p5-fd-$l"; done
  bis_fertig; leere
  (cd "$g/bare" && printf '%s' "$(zl main "$h1")" | env -i PATH="$g/bin:$PFAD" GIT_DIR=. "$BASH" hooks/post-receive 2> "$g/e2" | cat > /dev/null)
  ! grep -qs ende "$g/steuer/protokoll" || echo p5-leitung-stdout
  bis_fertig; leere
  (cd "$g/bare" && printf '%s' "$(zl main "$h1")" | env -i PATH="$g/bin:$PFAD" GIT_DIR=. "$BASH" hooks/post-receive 2>&1 > /dev/null | cat > /dev/null)
  ! grep -qs ende "$g/steuer/protokoll" || echo p5-leitung-stderr
  bis_fertig; leere
  (cd "$g/bare" && printf '%s' "$(zl main "$h1")" | env -i PATH="$g/bin:$PFAD" GIT_DIR=. setsid "$BASH" -c 'echo $$ > "$1"; "$2" hooks/post-receive; sleep 15' _ "$g/pgid" "$BASH" > /dev/null 2>&1) &
  warte_auf grep -qs beginn "$g/steuer/protokoll"
  pgid="$(cat "$g/pgid" 2> /dev/null || echo 0)"
  if [ "$pgid" -gt 1 ]; then kill -TERM -- "-$pgid" 2> /dev/null; fi
  bis_fertig
  { [ "$(feld "$(res main "$h1")" exit)" = 0 ] && [ ! -e "$(res main "$h1").laeuft" ]; } || echo p5-gruppe
}

# --- Punkt 6 und p9-beginn-nach-sperre ------------------------------------------------------------------------------------
g6() {
  local a b
  aufbau p6
  echo 2 > "$g/steuer/warte"
  haken "$(zl a "$h1")"; sleep 0.5; haken "$(zl b "$h2")"; bis_fertig
  [ "$(awk '{print $1 $2}' "$g/steuer/protokoll" | tr '\n' ' ')" = "beginn1 ende1 beginn2 ende2 " ] || echo p6-nacheinander
  { [ "$(feld "$(res a "$h1")" exit)" = 0 ] && [ "$(feld "$(res b "$h2")" exit)" = 0 ]; } || echo p6-kein-abbruch
  a="$(feld "$(res a "$h1")" ende)"; b="$(feld "$(res b "$h2")" beginn)"
  { [ -n "$b" ] && [[ ! $b < $a ]]; } || echo p9-beginn-nach-sperre
}

# --- Punkt 7 --------------------------------------------------------------------------------------------------------------
g7() {
  aufbau p7
  printf '1\n0\n' > "$g/steuer/exits"
  haken "$(zl main "$h1")"; bis_fertig
  { [ "$(feld "$(res main "$h1")" exit)" = 1 ] && [ "$(feld "$(res main "$h1")" ergebnis)" = rot ]; } || echo p7-neu
  haken "$(zl main "$h1")"; bis_fertig
  { [ "$(feld "$(res main "$h1")" exit)" = 0 ] && [ "$(feld "$(res main "$h1")" ergebnis)" = gruen ]; } || echo p7-neu
  echo 1.5 > "$g/steuer/warte"
  rm -f "${g:?}/steuer/protokoll"
  haken "$(zl main "$h1")"
  [ ! -e "$(res main "$h1")" ] || echo p7-altes-weg
  warte_auf grep -qs beginn "$g/steuer/protokoll"
  [ ! -e "$(res main "$h1")" ] || echo p7-altes-weg
  bis_fertig; leere; echo 0 > "$g/steuer/warte"
  haken "$(zs a "$h1" b "$h1")"; bis_fertig
  { [ "$(feld "$(res a "$h1")" zweig)" = a ] && [ "$(feld "$(res b "$h1")" zweig)" = b ]; } || echo p7-anderer-zweig
}

# --- Punkt 8 --------------------------------------------------------------------------------------------------------------
g8() {
  local lang120 lang121 e
  aufbau p8
  lang120="$(printf 'y%.0s' {1..120})"; lang121="$(printf 'x%.0s' {1..121})"
  haken "$(zs thema/teil "$h1" 'a%b' "$h1" 'ü' "$h1" 'a+b' "$h1" a.b_c-d "$h1" .x "$h1" a/b "$h1" 'a%2Fb' "$h1" "$lang121" "$h1" "$lang120" "$h1")"
  bis_fertig
  e="$g/ablage/gegenproben"
  { [ -f "$e-thema%2Fteil-$h1" ] && [ ! -d "$e-thema" ]; } || echo p8-schraeg
  [ -f "$e-a%25b-$h1" ] || echo p8-prozent
  [ -f "$e-%C3%BC-$h1" ] || echo p8-nichtascii
  [ -f "$e-a%2Bb-$h1" ] || echo p8-sonderzeichen
  { [ -f "$e-a.b_c-d-$h1" ] && [ -f "$e-.x-$h1" ]; } || echo p8-unveraendert
  { [ "$(feld "$e-a%2Fb-$h1" zweig)" = a/b ] && [ "$(feld "$e-a%252Fb-$h1" zweig)" = 'a%2Fb' ]; } || echo p8-eindeutig
  { grep -q '^gegenproben-nach-einlieferung: name: ' "$g/h.err" && [ ! -e "$e-$lang121-$h1" ] && [ ! -e "$e-$lang121-$h1.laeuft" ]; } || echo p8-lang
  [ -f "$e-$lang120-$h1" ] || echo p8-grenze
}

# --- Punkt 9 --------------------------------------------------------------------------------------------------------------
g9() {
  local r k jetzt zeit='^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$'
  aufbau p9
  printf '0\n2\n' > "$g/steuer/exits"
  printf 'gegenproben: za: rot (Exit 1)\ngegenproben: zc: grün\ngegenproben: zb: rot (Exit 3)\n' > "$g/steuer/rot.2"
  haken "$(zs x "$h1" y "$h2")" TZ=UTC-5; bis_fertig
  r="$(res x "$h1")"
  k="$(cut -d= -f1 "$r" | tr '\n' ' ')"
  jetzt="$(date -u +%s)"
  { [ "$k" = "zweig hash exit ergebnis beginn ende rot " ] && [ "$(feld "$r" exit)" = 0 ] && [ "$(feld "$r" ergebnis)" = gruen ] && [ -z "$(feld "$r" rot)" ] \
    && [[ "$(feld "$r" beginn)" =~ $zeit ]] && [[ "$(feld "$r" ende)" =~ $zeit ]] && [[ ! "$(feld "$r" beginn)" > "$(feld "$r" ende)" ]]; } || echo p9-gruen
  for k in beginn ende; do
    k="$(date -u -d "$(feld "$r" "$k")" +%s 2> /dev/null || echo 0)"
    { [ "$((jetzt - k))" -ge -5 ] && [ "$((jetzt - k))" -le 120 ]; } || echo p9-gruen
  done
  r="$(res y "$h2")"
  { [ "$(feld "$r" exit)" = 2 ] && [ "$(feld "$r" ergebnis)" = rot ] && [ "$(feld "$r" rot)" = "za zb" ]; } || echo p9-rot
  leere; rm -rf "${g:?}"/steuer/aufruf.* "${g:?}/steuer/zaehler" "${g:?}/steuer/rot.2"; : > "$g/steuer/exits"
  echo 2 > "$g/steuer/warte"
  haken "$(zl m "$h3")"
  [ -e "$(res m "$h3").laeuft" ] || echo p9-marke
  warte_auf grep -qs beginn "$g/steuer/protokoll"
  { ! ls -A "$g/ablage" | grep -Eq '^gegenproben-.*-[0-9a-f]{40}$'; } || echo p9-nicht-vorher
  bis_fertig
  [ ! -e "$(res m "$h3").laeuft" ] || echo p9-marke
  { [ "$(grep -c '^aus$' "$(res m "$h3").log")" -eq 1 ] && grep -qx fehl "$(res m "$h3").log"; } || echo p9-log
  echo 0 > "$g/steuer/warte"
  haken "$(zl m "$h3")"; bis_fertig
  [ "$(grep -c '^aus$' "$(res m "$h3").log")" -eq 1 ] || echo p9-log
}

# lauf_direkt <Name> <Zweig> <Hash> — der Hintergrundlauf unmittelbar (--lauf) in eigener Sitzung; sein Exit steht in $g/rc.
lauf_direkt() {
  : > "$g/ablage/$1.laeuft"
  (env -i PATH="$g/bin:$PFAD" HOME="$g" setsid -w "$BASH" "$g/bare/hooks/post-receive" --lauf "$g/bare" "$g/ablage" "$1" "$2" "$3" < /dev/null > /dev/null 2>&1; echo $? > "$g/rc") &
}
tot() { ! kill -0 "$(cat "$1" 2> /dev/null || echo 0)" 2> /dev/null; }

g9s() {
  local n bgpid
  aufbau p9s
  touch "$g/steuer/block"; n="gegenproben-z-$h1"
  lauf_direkt "$n" z "$h1"
  warte_auf test -s "$g/steuer/block.pid"
  bgpid="$(pgrep -of "post-receive --lauf $g/bare" || echo 0)"
  if [ "$bgpid" -gt 1 ]; then kill -TERM "$bgpid"; fi
  WA=40 warte_auf test -s "$g/rc"
  { [ "$(cat "$g/rc" 2> /dev/null)" = 143 ] && [ ! -e "$g/ablage/$n" ] && [ -e "$g/ablage/$n.laeuft" ] && tot "$g/steuer/block.pid" && tot "$g/steuer/kind.pid" \
    && ! ls -A "$g/ablage" | grep -q '^arbeitsbaum' && fertig; } || echo p9-sigterm-lauf
  pkill -KILL -f -- "$g/" 2> /dev/null
  rm -f "${g:?}/rc" "${g:?}/steuer/block" "${g:?}/steuer/zaehler"; rm -rf "${g:?}"/steuer/aufruf.*
  flock "$g/ablage/gegenproben.lock" sleep 30 &
  sleep 0.5
  lauf_direkt "$n" z "$h1"
  sleep 0.5
  bgpid="$(pgrep -of "post-receive --lauf $g/bare" || echo 0)"
  if [ "$bgpid" -gt 1 ]; then kill -TERM "$bgpid"; fi
  WA=40 warte_auf test -s "$g/rc"
  { [ "$(cat "$g/rc" 2> /dev/null)" = 143 ] && [ ! -e "$g/ablage/$n" ] && [ -e "$g/ablage/$n.laeuft" ] && [ ! -e "$g/steuer/zaehler" ] && fertig; } || echo p9-sigterm-wartet
  pkill -KILL -f -- "$g/" 2> /dev/null
}

# --- Punkt 10 -------------------------------------------------------------------------------------------------------------
g10() {
  aufbau p10
  haken "$(zl main "$h1")"; bis_fertig
  [ "$(tr '\n' ' ' < "$g/steuer/aufruf.1/argv" 2> /dev/null)" = "gegenproben " ] || echo p10-argv
}

# --- Punkt 11 -------------------------------------------------------------------------------------------------------------
g11() {
  local e v
  aufbau p11; mkdir "$g/quar"
  haken "$(zl main "$h1")" MAKEFLAGS=-i MFLAGS=-i MAKELEVEL=3 MAKEOVERRIDES=x=1; bis_fertig
  e="$g/steuer/aufruf.1/env"
  if [ -s "$e" ]; then
    for v in MAKEFLAGS MFLAGS MAKELEVEL MAKEOVERRIDES; do ! grep -q "^$v=" "$e" || echo "p11-$(echo "$v" | tr 'A-Z' 'a-z')"; done
    ! grep -q '^GIT_DIR=' "$e" || echo p11-git-dir
  else
    echo p11-makeflags p11-mflags p11-makelevel p11-makeoverrides p11-git-dir | tr ' ' '\n'
  fi
  zurueck
  haken "$(zl main "$h1")" GIT_QUARANTINE_PATH="$g/quar"; bis_fertig
  e="$g/steuer/aufruf.1/env"
  { [ -s "$e" ] && ! grep -q '^GIT_QUARANTINE_PATH=' "$e"; } || echo p11-git-quarantaene
}

# --- Punkt 12 -------------------------------------------------------------------------------------------------------------
g12() {
  local a vor
  aufbau p12
  vor="$(git -C "$g/bare" for-each-ref; cat "$g/bare/config" "$g/bare/HEAD")"
  haken "$(zs main "$h1" main "$h2")"; bis_fertig
  a="$g/steuer/aufruf.1"
  { [[ "$(cat "$a/cwd")" == "$(cd "$g/ablage" && pwd -P)"/arbeitsbaum.* ]] && [ "$(cat "$a/head")" = "$h1" ] && [ ! -s "$a/status" ] && [ "$(cat "$a/symrc")" = 1 ] \
    && [ "$(cat "$g/steuer/aufruf.2/head")" = "$h2" ] && [ "$(cat "$a/cwd")" != "$(cat "$g/steuer/aufruf.2/cwd")" ]; } || echo p12-stand
  ! ls -A "$g/ablage" | grep -q '^arbeitsbaum' || echo p12-entfernt-gruen
  rm -rf "${g:?}"/steuer/aufruf.* "${g:?}/steuer/zaehler"; printf '0\n3\n' > "$g/steuer/exits"
  haken "$(zl main "$h3")"; bis_fertig
  haken "$(zl rot "$h1")"; bis_fertig
  { [ "$(feld "$(res rot "$h1")" exit)" = 3 ] && ! ls -A "$g/ablage" | grep -q '^arbeitsbaum'; } || echo p12-entfernt-rot
  [ "$(git -C "$g/bare" for-each-ref; cat "$g/bare/config" "$g/bare/HEAD")" = "$vor" ] || echo p12-bare-unveraendert
  mkdir "$g/ablage/arbeitsbaum.alt"; echo x > "$g/ablage/arbeitsbaum.alt/datei"
  haken "$(zl main "$h2")"; bis_fertig
  { [ -f "$g/ablage/arbeitsbaum.alt/datei" ] && [[ "$(cat "$g/steuer/aufruf.3/cwd")" != */arbeitsbaum.alt ]]; } || echo p12-rest
}

# --- Punkt 13 -------------------------------------------------------------------------------------------------------------
g13() {
  local t w h40
  aufbau p13
  h40="$(printf 'c%.0s' {1..40})"
  mkdir "$g/pfad"
  for t in git flock setsid mktemp rm pgrep date; do ln -s "$(command -v "$t")" "$g/pfad/$t"; done
  ln -s "$g/bin/make" "$g/pfad/make"
  for w in git flock setsid make; do
    pfad="$g/pfad-$w"; cp -a "$g/pfad" "$pfad"; rm -f "${pfad:?}/${w:?}"
    haken "$(zl main "$h1")"; sleep 0.3
    { nichts && grep -q "^gegenproben-nach-einlieferung: werkzeug: $w " "$g/h.err"; } || echo "p13-werkzeug-$w"
    zurueck
  done
  pfad=""
  chmod 555 "$g/ablage"; haken "$(zl main "$h1")"; sleep 0.3; chmod 755 "$g/ablage"
  { [ "$hrc" -eq 0 ] && grep -q '^gegenproben-nach-einlieferung: ablage: ' "$g/h.err" && [ -z "$(ls -A "$g/ablage")" ] && [ ! -e "$g/steuer/zaehler" ]; } || echo p13-ablage-weg; zurueck
  mkdir "$g/ablage/gegenproben.lock"
  haken "$(zl main "$h1")"; bis_fertig
  { [ "$hrc" -eq 0 ] && grep -q '^gegenproben-nach-einlieferung: sperre: ' "$(res main "$h1").log" && [ ! -e "$(res main "$h1")" ] && [ -e "$(res main "$h1").laeuft" ] \
    && [ ! -e "$g/steuer/zaehler" ]; } || echo p13-sperre
  rmdir "$g/ablage/gegenproben.lock"; leere
  haken "$(zl main "$h40")"; bis_fertig
  { [ "$hrc" -eq 0 ] && grep -q '^gegenproben-nach-einlieferung: arbeitsbaum: ' "$(res main "$h40").log" && [ ! -e "$(res main "$h40")" ] && [ -e "$(res main "$h40").laeuft" ] \
    && [ ! -e "$g/steuer/zaehler" ]; } || echo p13-hash-fehlt
  leere; echo 127 > "$g/steuer/exits"
  haken "$(zl main "$h1")"; bis_fertig
  { [ "$(feld "$(res main "$h1")" exit)" = 127 ] && [ "$(feld "$(res main "$h1")" ergebnis)" = rot ]; } || echo p13-make-127
  leere; rm -f "${g:?}/bare/gegenproben-ablage"; pfad="$g/pfad-git"
  haken "$(zl main "$h1")"; sleep 0.3
  { grep -q 'marke: ' "$g/h.err" && ! grep -q 'werkzeug: ' "$g/h.err"; } || echo p13-marke-vor-werkzeug; zurueck
  echo "$g/ablage" > "$g/bare/gegenproben-ablage"
  chmod 555 "$g/ablage"; haken "$(zl main "$h1")"; sleep 0.3; chmod 755 "$g/ablage"
  { grep -q 'werkzeug: ' "$g/h.err" && ! grep -q 'ablage: ' "$g/h.err"; } || echo p13-werkzeug-vor-ablage; zurueck
  pfad=""
  chmod 555 "$g/ablage"; haken "kaputt"; sleep 0.3; chmod 755 "$g/ablage"
  { grep -q 'ablage: ' "$g/h.err" && ! grep -q 'zeile: ' "$g/h.err"; } || echo p13-ablage-vor-zeile
}

# pruefe — schreibt die Namen der roten Fälle, einen je Zeile; die Gruppen laufen gleichzeitig.
pruefe() {
  local f
  for f in g1 g2 g3 g4 g5 g6 g7 g8 g9 g9s g10 g11 g12 g13; do ("$f" > "$fx/rot.$f" 2> "$fx/err.$f") & done
  wait
  cat "$fx"/rot.g*
  pkill -KILL -f -- "$fx/" 2> /dev/null
}
ergebnis() { pruefe | LC_ALL=C sort -u | tr '\n' ' '; }

maxjobs="${GP_PARALLEL:-4}"
wartejob() { while [ "$(jobs -rp | wc -l)" -ge "$maxjobs" ]; do wait -n || true; done; }

# mut <name> <rote Fälle|?> <alt> <neu> [<alt> <neu>…] — Textersatz im Skript der Kopie; rot sind genau die benannten Fälle.
mut() {
  local name="$1" rote="$2" inhalt neu ist erwartet
  shift 2
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  wartejob
  (
    fx="$arbeit/$name"; hook="$fx/hook.sh"; mkdir -p "$fx"
    inhalt="$(cat "$skript"; printf x)"; inhalt="${inhalt%x}"
    neu="$inhalt"
    while [ "$#" -ge 2 ]; do
      if [[ $neu != *"$1"* ]]; then melde "Mutation '$name': Textstelle nicht gefunden: $1"; exit 0; fi
      neu="${neu/"$1"/"$2"}"
      shift 2
    done
    printf '%s' "$neu" > "$hook"; chmod 755 "$hook"
    ist="$(ergebnis)"
    erwartet="$(printf '%s\n' $rote | LC_ALL=C sort -u | tr '\n' ' ')"
    if [ -n "${GP_ENTDECKE-}" ]; then echo "$name => ${ist:-nichts}" >&2; fi
    [ "$rote" = "?" ] || [ "$ist" = "$erwartet" ] || melde "Mutation '$name': erwartet rot: ${rote:-nichts}; rot sind: ${ist:-nichts}"
  ) &
}

if [[ ok =~ ${GP_MUT:-.} ]]; then
  (
    fx="$arbeit/ok"; hook="$fx/hook.sh"; mkdir -p "$fx"; cp "$skript" "$hook"; chmod 755 "$hook"
    ist="$(ergebnis)"
    [ -z "$ist" ] || melde "unveränderte Kopie ist nicht grün: $ist"
  ) &
fi

NL=$'\n'
SIGL='  [ "${#zeilen[@]}" -eq 1 ] || { melde marke "keine einzelne Zeile in $marke"; return 0; }'
ABSL='  case "$ablage" in /*) ;; *) melde marke "kein absoluter Pfad: $ablage"; return 0 ;; esac'
DIRL='  [ -d "$ablage" ] || { melde marke "kein Verzeichnis: $ablage"; return 0; }'
WERK='  for t in git flock setsid make; do command -v "$t" > /dev/null 2>&1 || { melde werkzeug "$t fehlt"; return 0; }; done'
PROBE='  probe="$(mktemp "$ablage/.probe.XXXXXX" 2> /dev/null)" || { melde ablage "nicht schreibbar: $ablage"; return 0; }; rm -f "${probe:?}"'
HASHL='    [[ $neu =~ ^([0-9a-f]{40}|[0-9a-f]{64})$ ]] || { melde zeile "Form des Hashes: $zeile"; continue; }'
NULLL='    [[ $neu =~ ^(0{40}|0{64})$ ]] && continue'
REFL='    case "$ref" in refs/heads/*) ;; *) continue ;; esac'
ALTL='    rm -f "${ablage:?}/${n:?}"'
MARKE_WEG='  rm -f "${ablage:?}/${n:?}.laeuft"'
FLOCK='  flock 9 & pid=$!'
START='< /dev/null > /dev/null 2>&1 &'
KODE=$'      *) out+="$(printf \'%%%02X\' "\'$c")" ;;'
WTADD='git --git-dir="$bare" worktree add --detach "$wt" "$hash"'
MAKEL='(cd "$wt" && exec make gegenproben)'
ABBR='abbruch() { trap '"''"' TERM; '
WTFEHL='    rm -rf "${wt:?}"; wt=""; exec 9>&-; return 0'
SPERRE='"$sig: sperre: Sperre nicht angelegt" > "$log";'
# --- Punkt 1
mut m-p1-laedt-nach "p1-kopie p3-headsx p3-notes p3-tag p4-ergebnis-bleibt p4-hash-64 p4-null40 p4-null64 p4-stdin-leer p5-gestartet-zeile" 'sig="gegenproben-nach-einlieferung"' $'sig="gegenproben-nach-einlieferung"\n. "$(dirname "${BASH_SOURCE[0]}")/bibliothek.sh"'
mut m-p1-liest-stdin-nicht "p1-echter-push p1-kopie p10-argv p11-git-dir p11-git-quarantaene p11-makeflags p11-makelevel p11-makeoverrides p11-mflags p12-entfernt-rot p12-stand p13-hash-fehlt p13-make-127 p13-sperre p3-mehrere p4-feld-fehlt p4-hash-41 p4-hash-64 p4-hash-gross p4-hash-kurz p4-zeile-kaputt-dann-gut p5-fd-0 p5-fd-1 p5-fd-2 p5-gestartet-zeile p5-gruppe p6-kein-abbruch p6-nacheinander p7-anderer-zweig p7-neu p8-eindeutig p8-grenze p8-lang p8-nichtascii p8-prozent p8-schraeg p8-sonderzeichen p8-unveraendert p9-beginn-nach-sperre p9-gruen p9-log p9-marke p9-rot" 'while IFS= read -r zeile || [ -n "$zeile" ]; do' 'while IFS= read -r zeile < /dev/null || [ -n "$zeile" ]; do'
# --- Punkt 2
mut m-p2-zeilenzahl "p2-zwei-zeilen" "$SIGL" '  :'
mut m-p2-absolut "p2-relativ" "$ABSL" '  :'
mut m-p2-verzeichnis "p2-kein-verzeichnis" "$DIRL" '  :'
mut m-p2-nur-existenz "p2-kein-verzeichnis p2-leer p2-relativ p2-zwei-zeilen" "$SIGL" '  [ -e "$marke" ] || { melde marke "fehlt"; return 0; }' "$ABSL" '  :' "$DIRL" '  :'
mut m-p2-keine-marke "p13-marke-vor-werkzeug p2-fehlt p2-kein-verzeichnis p2-leer p2-relativ p2-zwei-zeilen" "$SIGL" '  :' "$ABSL" '  :' "$DIRL" '  :'
mut m-p2-legt-an "p2-fehlt p2-kein-verzeichnis p2-leer p2-relativ p2-zwei-zeilen" "$SIGL" $'  echo x >> "$bare/gegenproben-spur"\n'"$SIGL"
# --- Punkt 3
mut m-p3-refs-stern "p3-headsx p3-notes p3-tag" "$REFL" '    case "$ref" in refs/*) ;; *) continue ;; esac'
mut m-p3-heads-ohne-schraeg "p3-headsx" "$REFL" '    case "$ref" in refs/heads*) ;; *) continue ;; esac'
mut m-p3-gleichzeitig "p3-mehrere p6-nacheinander p9-beginn-nach-sperre p9-sigterm-lauf p9-sigterm-wartet" 'do lauf "$1" "$2" "$3"; shift 3; done' 'do lauf "$1" "$2" "$3" & sleep 0.1; shift 3; done; wait' "$FLOCK" '  : & pid=$!'
mut m-p3-rueckwaerts "p12-stand p3-mehrere p9-gruen p9-rot" '    args+=("$n" "$zweig" "$neu")' '    args=("$n" "$zweig" "$neu" "${args[@]}")'
mut m-p3-nur-erste "p12-stand p3-mehrere p7-anderer-zweig p8-eindeutig p8-grenze p8-lang p8-nichtascii p8-prozent p8-sonderzeichen p8-unveraendert p9-rot" '    args+=("$n" "$zweig" "$neu")' '    args+=("$n" "$zweig" "$neu"); break'
mut m-p3-meldet-fremde-refs "p3-headsx p3-notes p3-tag" "$REFL" '    case "$ref" in refs/heads/*) ;; *) melde zeile "kein Zweig"; continue ;; esac'
# --- Punkt 4
mut m-p4-null-laeuft "p4-ergebnis-bleibt p4-null40 p4-null64" "$NULLL" '    :'
mut m-p4-null-nur-40 "p4-null64" "$NULLL" '    [[ $neu =~ ^0{40}$ ]] && continue'
mut m-p4-keine-form "p4-hash-41 p4-hash-gross p4-hash-kurz p4-zeile-kaputt-dann-gut" "$HASHL" '    :'
mut m-p4-nur-40 "p4-hash-64" '^([0-9a-f]{40}|[0-9a-f]{64})$' '^[0-9a-f]{40}$'
mut m-p4-grossbuchstaben "p4-hash-gross" '^([0-9a-f]{40}|[0-9a-f]{64})$' '^([0-9a-fA-F]{40}|[0-9a-f]{64})$'
mut m-p4-41-stellen "p4-hash-41" '^([0-9a-f]{40}|[0-9a-f]{64})$' '^([0-9a-f]{40,41}|[0-9a-f]{64})$'
mut m-p4-kaputte-zeile-beendet "p4-zeile-kaputt-dann-gut" '"Form des Hashes: $zeile"; continue; }' '"Form des Hashes: $zeile"; return 0; }'
mut m-p4-feld-ohne-pruefung "p4-feld-fehlt" '    [ -n "$ref" ] || { melde zeile "Feld fehlt: $zeile"; continue; }' '    :'
mut m-p4-feld-beendet "p4-feld-fehlt" '"Feld fehlt: $zeile"; continue; }' '"Feld fehlt: $zeile"; return 0; }'
mut m-p4-loeschung-entfernt-ergebnis "p4-ergebnis-bleibt" "$NULLL" '    [[ $neu =~ ^(0{40}|0{64})$ ]] && { rm -f "${ablage:?}"/gegenproben-b-*; continue; }'
mut m-p4-loeschung-meldet "p4-ergebnis-bleibt p4-null40 p4-null64" "$NULLL" '    [[ $neu =~ ^(0{40}|0{64})$ ]] && { melde zeile "geloescht"; continue; }'
mut m-p4-leere-eingabe-meldet "p3-headsx p3-notes p3-tag p4-ergebnis-bleibt p4-feld-fehlt p4-null40 p4-null64 p4-stdin-leer" '  if [ "${#args[@]}" -gt 0 ]; then' $'  [ "${#args[@]}" -gt 0 ] || melde zeile "keine Eingabe"\n  if [ "${#args[@]}" -gt 0 ]; then'
# --- Punkt 5
mut m-p5-wartet "p1-echter-push p5-fd-0 p5-fd-1 p5-fd-2 p5-leitung-stderr p5-leitung-stdout p5-vor-ende p7-altes-weg p9-marke p9-nicht-vorher" "$START" '< /dev/null > /dev/null 2>&1'
mut m-p5-erbt-stdout "p1-echter-push p5-fd-1 p5-leitung-stdout" "$START" '< /dev/null >&1 2> /dev/null &'
mut m-p5-erbt-stderr "p1-echter-push p5-fd-2 p5-leitung-stderr" "$START" '< /dev/null > /dev/null 2>&2 &'
mut m-p5-erbt-stdin "p5-fd-0" "$START" '<&0 > /dev/null 2>&1 &'
mut m-p5-gleiche-sitzung "p5-gruppe" 'setsid "$BASH" "$self"' '"$BASH" "$self"'
mut m-p5-stdout-zeile "p5-stdout-leer" '    melde gestartet "$zweig $neu"' $'    melde gestartet "$zweig $neu"\n    echo hallo'
mut m-p5-gestartet-fehlt "p1-kopie p4-hash-64 p4-zeile-kaputt-dann-gut p5-gestartet-zeile" '    melde gestartet "$zweig $neu"' '    :'
# --- Punkt 6
mut m-p6-ohne-sperre "p6-nacheinander p9-beginn-nach-sperre p9-sigterm-wartet" "$FLOCK" '  : & pid=$!'
mut m-p6-sperre-nicht-erwartet "p6-kein-abbruch p6-nacheinander p9-beginn-nach-sperre p9-sigterm-wartet" "$FLOCK" '  flock -n 9 & pid=$!'
mut m-p6-zweiter-beendet-ersten "p6-kein-abbruch p6-nacheinander" "$FLOCK" $'  p="$(pgrep -of "post-receive --lauf $bare")"; [ "$p" = "$$" ] || kill -TERM "$p"\n'"$FLOCK"
# --- Punkt 7
mut m-p7-lauf-entfaellt "p7-altes-weg p7-neu" "$ALTL" '    [ -e "$ablage/$n" ] && continue'
mut m-p7-altes-bleibt "p7-altes-weg" "$ALTL" '    :'
mut m-p7-name-ohne-zweig "p1-echter-push p1-kopie p12-entfernt-rot p13-hash-fehlt p13-make-127 p13-sperre p3-mehrere p4-hash-64 p4-zeile-kaputt-dann-gut p5-gruppe p6-kein-abbruch p7-anderer-zweig p7-neu p8-eindeutig p8-grenze p8-nichtascii p8-prozent p8-schraeg p8-sonderzeichen p8-unveraendert p9-beginn-nach-sperre p9-gruen p9-log p9-marke p9-rot" '    n="gegenproben-$name-$neu"' '    n="gegenproben-$neu"'
# --- Punkt 8
mut m-p8-keine-kodierung "p8-eindeutig p8-nichtascii p8-prozent p8-schraeg p8-sonderzeichen" "$KODE" '      *) out+="$c" ;;'
mut m-p8-prozent-bleibt "p8-eindeutig p8-prozent" "$KODE" $'      %) out+="$c" ;;\n'"$KODE"
mut m-p8-nur-schraeg "p8-eindeutig p8-nichtascii p8-prozent p8-sonderzeichen" $'      [A-Za-z0-9._-]) out+="$c" ;;\n'"$KODE" $'      /) out+="%2F" ;;\n      *) out+="$c" ;;'
mut m-p8-hex-klein "p8-eindeutig p8-nichtascii p8-schraeg p8-sonderzeichen" "'%%%02X'" "'%%%02x'"
mut m-p8-ohne-punkt "p8-unveraendert" '[A-Za-z0-9._-])' '[A-Za-z0-9_-])'
mut m-p8-ohne-unterstrich "p8-unveraendert" '[A-Za-z0-9._-])' '[A-Za-z0-9.-])'
mut m-p8-ohne-bindestrich "p8-unveraendert" '[A-Za-z0-9._-])' '[A-Za-z0-9._])'
mut m-p8-ohne-ziffern "p8-eindeutig" '[A-Za-z0-9._-])' '[A-Za-z._-])'
mut m-p8-ohne-grossbuchstaben "p8-eindeutig" '[A-Za-z0-9._-])' '[a-z0-9._-])'
mut m-p8-ohne-kleinbuchstaben "p1-echter-push p1-kopie p12-entfernt-rot p13-hash-fehlt p13-make-127 p13-sperre p3-mehrere p4-hash-64 p4-zeile-kaputt-dann-gut p5-gruppe p6-kein-abbruch p7-anderer-zweig p7-neu p8-eindeutig p8-grenze p8-prozent p8-schraeg p8-sonderzeichen p8-unveraendert p9-beginn-nach-sperre p9-gruen p9-log p9-marke p9-rot" '[A-Za-z0-9._-])' '[A-Z0-9._-])'
mut m-p8-grenze-121 "p8-lang" '"${#name}" -le 120' '"${#name}" -le 121'
mut m-p8-grenze-119 "p8-grenze" '"${#name}" -le 120' '"${#name}" -le 119'
mut m-p8-lang-ohne-meldung "p8-lang" 'melde name "länger als 120 Zeichen: $name"; continue; }' ': ; continue; }'
# --- Punkt 9
mut m-p9-gruen-immer "p13-make-127 p7-neu p9-rot" '  erg=rot; [ "$rc" -eq 0 ] && erg=gruen' '  erg=gruen'
mut m-p9-exit1-gruen "p7-neu" '[ "$rc" -eq 0 ] && erg=gruen' '[ "$rc" -le 1 ] && erg=gruen'
mut m-p9-rot-nur-erstes "p9-rot" 'rot+="${rot:+ }${BASH_REMATCH[1]}"' 'rot="${rot:-${BASH_REMATCH[1]}}"'
mut m-p9-rot-rueckwaerts "p9-rot" 'rot+="${rot:+ }${BASH_REMATCH[1]}"' 'rot="${BASH_REMATCH[1]}${rot:+ }$rot"'
mut m-p9-zweig-hash-vertauscht "p1-echter-push p7-anderer-zweig p8-eindeutig p9-gruen" "'zweig=%s\\nhash=%s\\n" "'hash=%s\\nzweig=%s\\n"
mut m-p9-exit-ergebnis-vertauscht "p12-entfernt-rot p13-make-127 p5-gruppe p6-kein-abbruch p7-neu p9-gruen p9-rot" '\nexit=%s\nergebnis=%s\n' '\nergebnis=%s\nexit=%s\n'
mut m-p9-beginn-ende-vertauscht "p9-gruen" '\nbeginn=%s\nende=%s\n' '\nende=%s\nbeginn=%s\n'
mut m-p9-ohne-rot-zeile "p9-gruen p9-rot" '\nrot=%s\n' '\n' '"$ende" "$rot" >' '"$ende" >'
mut m-p9-beginn-lokal "p9-gruen" 'beginn="$(date -u ' 'beginn="$(date '
mut m-p9-ende-lokal "p9-gruen" 'ende="$(date -u ' 'ende="$(date '
mut m-p9-marke-bleibt "p5-gruppe p9-marke" "$MARKE_WEG" '  :'
mut m-p9-marke-fehlt "p1-echter-push p13-hash-fehlt p13-sperre p4-hash-64 p9-marke" '    : > "$ablage/$n.laeuft"' '    :'
mut m-p9-ergebnis-vorher "p7-altes-weg p9-nicht-vorher p9-sigterm-lauf p9-sigterm-wartet" '  # (7) der Lauf' '  : > "$ablage/$n"'
mut m-p9-ergebnis-bei-sigterm "p9-sigterm-lauf p9-sigterm-wartet" "$ABBR" "$ABBR"$'printf \'ergebnis=gruen\\n\' > "$ablage/$n"; '
mut m-p9-marke-weg-bei-sigterm "p9-sigterm-lauf p9-sigterm-wartet" "$ABBR" "$ABBR"$'rm -f "${ablage:?}"/*.laeuft; '
mut m-p9-nachkommen-bleiben "p9-sigterm-lauf" 'for k in $(pgrep -P "$1" 2> /dev/null); do beende "$k" "$2"; done; ' ''
mut m-p9-sigterm-wartet-unbehandelt "p9-sigterm-wartet" $'  flock 9 & pid=$!\n  wait "$pid"; rc=$?; pid=""' '  flock 9; rc=$?'
mut m-p9-exit0 "p9-sigterm-lauf p9-sigterm-wartet" 'raeume; exit 143; }' 'raeume; exit 0; }'
mut m-p9-sigterm-ohne-trap "p9-sigterm-lauf" '  trap abbruch TERM' '  :'
mut m-p9-arbeitsbaum-bleibt-bei-sigterm "p9-sigterm-lauf" 'wait "$pid" 2> /dev/null; fi; raeume; exit 143; }' 'wait "$pid" 2> /dev/null; fi; exit 143; }'
mut m-p9-log-ohne-stderr "p9-log p9-rot" '> "$log" 2>&1 9>&- &' '> "$log" 2> /dev/null 9>&- &'
mut m-p9-log-ohne-stdout "p9-log" '> "$log" 2>&1 9>&- &' '2> "$log" 9>&- &'
mut m-p9-log-anhaengen "p9-log" '> "$log" 2>&1 9>&- &' '>> "$log" 2>&1 9>&- &'
mut m-p9-beginn-vor-sperre "p9-beginn-nach-sperre" $'  beginn="$(date -u +%Y-%m-%dT%H:%M:%SZ)"\n  (cd' $'  (cd' '  # (5) Sperre' $'  beginn="$(date -u +%Y-%m-%dT%H:%M:%SZ)"'
# --- Punkt 10
mut m-p10-host "p10-argv" "$MAKEL" '(cd "$wt" && exec make gegenproben-host)'
mut m-p10-k "p10-argv" "$MAKEL" '(cd "$wt" && exec make -k gegenproben)'
mut m-p10-j "p10-argv" "$MAKEL" '(cd "$wt" && exec make -j4 gegenproben)'
mut m-p10-i "p10-argv" "$MAKEL" '(cd "$wt" && exec make -i gegenproben)'
mut m-p10-n "p10-argv" "$MAKEL" '(cd "$wt" && exec make -n gegenproben)'
mut m-p10-variable "p10-argv" "$MAKEL" '(cd "$wt" && exec make gegenproben X=1)'
mut m-p10-einzelziel "p10-argv" "$MAKEL" '(cd "$wt" && exec make a-check-negativ)'
# --- Punkt 11
mut m-p11-makeflags-bleibt "p11-makeflags" 'unset MAKEFLAGS MFLAGS' 'unset MFLAGS'
mut m-p11-mflags-bleibt "p11-mflags" 'unset MAKEFLAGS MFLAGS MAKELEVEL' 'unset MAKEFLAGS MAKELEVEL'
mut m-p11-makelevel-bleibt "p11-makelevel" 'MFLAGS MAKELEVEL MAKEOVERRIDES' 'MFLAGS MAKEOVERRIDES'
mut m-p11-makeoverrides-bleibt "p11-makeoverrides" 'MAKELEVEL MAKEOVERRIDES' 'MAKELEVEL'
mut m-p11-git-dir-bleibt "p11-git-dir p12-stand p3-mehrere" 'do unset "$v"; done' 'do [ "$v" = GIT_DIR ] || unset "$v"; done'
mut m-p11-quarantaene-bleibt "p11-git-quarantaene" 'do unset "$v"; done' 'do [ "$v" = GIT_QUARANTINE_PATH ] || unset "$v"; done'
mut m-p11-nur-git-dir "p11-git-quarantaene" 'for v in $(compgen -v GIT_); do unset "$v"; done' 'unset GIT_DIR'
# --- Punkt 12
mut m-p12-im-bare "p12-stand p3-mehrere" "$MAKEL" '(cd "$bare" && exec make gegenproben)'
mut m-p12-spitze "p12-entfernt-rot p12-stand p13-hash-fehlt p13-make-127 p3-mehrere p4-zeile-kaputt-dann-gut p6-kein-abbruch p6-nacheinander p7-anderer-zweig p8-eindeutig p8-grenze p8-nichtascii p8-prozent p8-schraeg p8-sonderzeichen p8-unveraendert p9-beginn-nach-sperre p9-gruen p9-log p9-marke p9-rot p9-sigterm-lauf" "$WTADD" 'git --git-dir="$bare" worktree add --detach "$wt" "refs/heads/$zweig"'
mut m-p12-nicht-losgeloest "p11-git-quarantaene p12-bare-unveraendert p12-stand p5-gruppe p7-neu p9-log" "$WTADD" 'git --git-dir="$bare" worktree add -b "gp-$n" "$wt" "$hash"'
mut m-p12-bleibt-nach-gruen "p12-entfernt-gruen p12-entfernt-rot" "$MARKE_WEG"$'\n  raeume' "$MARKE_WEG"$'\n  [ "$rc" -eq 0 ] && wt=""\n  raeume'
mut m-p12-bleibt-nach-rot "p12-entfernt-rot" "$MARKE_WEG"$'\n  raeume' "$MARKE_WEG"$'\n  [ "$rc" -ne 0 ] && wt=""\n  raeume'
mut m-p12-fester-name "p12-stand" 'mktemp -d "$ablage/arbeitsbaum.XXXXXX"' 'echo "$ablage/arbeitsbaum.fest"; mkdir -p "$ablage/arbeitsbaum.fest"'
mut m-p12-loescht-reste "p12-rest" '  wt="$(mktemp -d' $'  rm -rf "${ablage:?}"/arbeitsbaum.*\n  wt="$(mktemp -d'
mut m-p12-aendert-ref "p12-bare-unveraendert" "$MARKE_WEG"$'\n  raeume' "$MARKE_WEG"$'\n  git --git-dir="$bare" update-ref refs/heads/gp-neu "$hash"\n  raeume'
mut m-p12-aendert-config "p12-bare-unveraendert" "$MARKE_WEG"$'\n  raeume' "$MARKE_WEG"$'\n  git --git-dir="$bare" config gp.x y\n  raeume'
# --- Punkt 13
mut m-p13-ohne-git "p13-werkzeug-git p13-werkzeug-vor-ablage" "$WERK" '  for t in flock setsid make; do command -v "$t" > /dev/null 2>&1 || { melde werkzeug "$t fehlt"; return 0; }; done'
mut m-p13-ohne-flock "p13-werkzeug-flock" "$WERK" '  for t in git setsid make; do command -v "$t" > /dev/null 2>&1 || { melde werkzeug "$t fehlt"; return 0; }; done'
mut m-p13-ohne-setsid "p13-werkzeug-setsid" "$WERK" '  for t in git flock make; do command -v "$t" > /dev/null 2>&1 || { melde werkzeug "$t fehlt"; return 0; }; done'
mut m-p13-ohne-make "p13-werkzeug-make" "$WERK" '  for t in git flock setsid; do command -v "$t" > /dev/null 2>&1 || { melde werkzeug "$t fehlt"; return 0; }; done'
mut m-p13-ohne-ablage-probe "p13-ablage-vor-zeile p13-ablage-weg" "$PROBE" '  :'
mut m-p13-werkzeug-vor-marke "p13-marke-vor-werkzeug" '  # (1) Marke des Containers (Punkt 2)' "$WERK"
mut m-p13-ablage-vor-werkzeug "p13-werkzeug-vor-ablage" '  # (2) Werkzeuge, (3) Ablage (Punkt 13)' "$PROBE"
mut m-p13-zeile-vor-ablage "p13-ablage-vor-zeile" "$PROBE" '  :' "$HASHL" "$HASHL"$'\n'"$PROBE"
mut m-p13-exit1 "p1-kopie p13-ablage-weg p13-hash-fehlt p13-sperre p13-werkzeug-flock p13-werkzeug-git p13-werkzeug-make p13-werkzeug-setsid p2-fehlt p2-kein-verzeichnis p2-leer p2-relativ p2-zwei-zeilen p3-headsx p3-notes p3-tag p4-ergebnis-bleibt p4-feld-fehlt p4-hash-41 p4-hash-gross p4-hash-kurz p4-null40 p4-null64 p4-stdin-leer p5-gestartet-zeile" $'haupt\nexit 0' $'haupt\nexit 1'
mut m-p13-gruen-bei-sperre "p13-sperre" "$SPERRE" "$SPERRE"$' printf \'ergebnis=gruen\\n\' > "$ablage/$n";'
mut m-p13-gruen-bei-arbeitsbaum "p13-hash-fehlt" "$WTFEHL" $'    printf \'ergebnis=gruen\\n\' > "$ablage/$n"; rm -rf "${wt:?}"; wt=""; exec 9>&-; return 0'
mut m-p13-marke-weg-bei-sperre "p13-sperre" "$SPERRE" "$SPERRE"$' rm -f "${ablage:?}/${n:?}.laeuft";'
mut m-p13-marke-weg-bei-arbeitsbaum "p13-hash-fehlt p4-hash-64" "$WTFEHL" $'    rm -f "${ablage:?}/${n:?}.laeuft"; rm -rf "${wt:?}"; wt=""; exec 9>&-; return 0'
mut m-p13-127-ohne-ergebnis "p13-make-127" '  ende="$(date -u' $'  [ "$rc" -eq 127 ] && { exec 9>&-; return 0; }\n  ende="$(date -u'

wait
if [ -s "$arbeit/fehler" ]; then
  cat "$arbeit/fehler" >&2
  exit 1
fi
echo "gegenproben-nach-einlieferung-gegenprobe: gruen — Punkte 1 bis 13 von SPEC-061 mit 74 Fällen und 113 Mutationen, je Mutation genau die benannten Fälle rot, unveränderte Kopie grün"
