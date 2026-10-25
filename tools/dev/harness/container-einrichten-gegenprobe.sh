#!/usr/bin/env bash
# container-einrichten-gegenprobe — prüft tools/harness/container-einrichten.sh und die Fragmente
# harness/mk/container-einrichten.mk und harness/mk/container-einrichten-gegenprobe.mk gegen SPEC-069
# (harness/vertraege/SPEC-069-container-einrichten.md) OHNE CONTAINER, ohne Docker, ohne Netz und ohne echtes
# ssh (Punkt 17, ADR-0057). Fälle gibt es zu den Punkten 1 bis 17; Punkt 18 ist der Nachweis, Punkt 19 die
# Plattform und Punkt 20 sind die Grenzen, ohne Fall. Je Fall eine Welt in einem Temp-Verzeichnis (mktemp -d): ein
# Host-Repo mit einem Commit und lokaler Identität, $PGR_DEV mit u-boot.yaml und ssh/id_ed25519, das
# Verzeichnis welt/workspaces/<Projektname>/ und ein Pfad, in dem ssh, make, docker, u-boot und devcontainer
# Attrappen sind, `id` eine Attrappe ist und git und rm das echte Werkzeug mit einem Protokoll davor sind.
# Die Attrappe von ssh führt den
# übergebenen Befehl in einer lokalen Shell in der Welt aus und protokolliert Argumente, Umgebung und
# Standardeingang; so laufen git init --bare, git clone, git push (über git-receive-pack), git config und
# git switch als ECHTES git, und die Zusagen der Punkte 5 bis 11 sind an echtem git-Verhalten geprüft.
# GRENZE (Punkt 17, akzeptiertes Negativ): Dass ein echtes ssh die Argumente so annimmt, dass der Befehl im
# Container läuft, und dass make hooks-install und make traeger-fetch im Klon gelingen, belegt die Gegenprobe
# nicht; das belegt der Nachweis aus Punkt 18 (Beleg in §7 des Slice-Plans).
# GRENZE (SPEC-060 Punkt 17, akzeptiertes Negativ): Diese Gegenprobe setzt SIGINT beim Start als nicht ignoriert
# voraus. Ein Start mit ignoriertem SIGINT (Hintergrundstart mit & oder nohup, auch setsid -f allein) ist für
# ihre Fälle mit SIGINT nicht zugesagt. Aufruf im Vordergrund oder mit setsid -f env --default-signal=INT,QUIT
# <Befehl>; Rezept, Messung und Grund stehen in SPEC-060 Punkt 17.
# Der Name jedes Falls trägt die Nummer seines Punktes in SPEC-069; die Fälle des Zeitbudgets heißen
# p17-budget-* und p17-ordner-* (SPEC-069 Punkt 17 verlangt sie, ihre Form steht in SPEC-060 Punkt 16 (9)).
# Die unveränderte Kopie besteht jeden Fall, und bei jeder Mutation (genau eine Textersetzung in einer Kopie;
# trifft sie nichts oder mehr als eine Stelle, bricht die Gegenprobe ab) werden von den Fällen ihrer Auswahl
# genau die benannten rot.
# Zeitbudget (SPEC-060 Punkt 16, ADR-0059): höchstens 1200 s für den vollen Lauf; die Zeile am Ende nennt Zeit,
# Obergrenze und Last, eine Überschreitung ist rot.
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Mutationen, deren Name passt (die unveränderte Kopie nur, wenn
# das Muster `ok` trifft); GP_NUR=<Muster> wählt die Fälle (der
# unveränderten Kopie und je Mutation; die Fälle p17-* geben einer Kopie damit einen einzigen billigen Fall); GP_ENTDECKE=1 nennt je
# Mutation die roten Fälle; GP_PARALLEL=<n> (Vorgabe 4) setzt die Zahl der gleichzeitigen Mutationen;
# GP_BUDGET=<s> setzt die Obergrenze (0 schaltet die Prüfung ab).
# Braucht bash ≥ 4.4 (mit `wait -n`), git, make, sed, awk, grep, find und die GNU-Werkzeuge von coreutils
# (dazu diff mit GP_ENTDECKE=1); als root entfällt der Fall mit Dateirechten (die Gegenprobe nennt ihn am Ende).
# Ausgang: 0, wenn jeder Fall und jede Mutation wie erwartet endet und die Zeit unter der Obergrenze liegt,
# sonst 1.
set -uo pipefail
export LC_ALL=C
T0=$SECONDS
GPN=container-einrichten-gegenprobe
LOADAVG=/proc/loadavg
GP_SELF=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/$(basename "${BASH_SOURCE[0]}")
quelle=${GP_QUELLE:-$(cd "$(dirname "$GP_SELF")/../../.." && pwd -P)}

# SPEC-060 Punkt 16 (4): GP_BUDGET setzt die Obergrenze in Sekunden, 0 schaltet die Prüfung ab; ohne die
# Variable gilt die Obergrenze aus (1), und das nur bei einem Lauf ohne GP_MUT. Ein ungültiger Wert endet vor
# jedem Fall mit Exit 2 (Reihenfolge der Gründe (a)).
budget=
if [ -n "${GP_BUDGET+x}" ]; then
  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] || { echo "$GPN: GP_BUDGET: keine ganze Zahl" >&2; exit 2; }
  budget=$GP_BUDGET
elif [ -z "${GP_MUT-}" ]; then
  budget=1200
fi

arbeit=$(mktemp -d) || { echo "$GPN: ROT — mktemp -d schlug fehl" >&2; exit 1; }
# SPEC-060 Punkt 16 (5): Bei Ende, SIGINT und SIGTERM entfällt der ganze Arbeitsordner; der Ordner einer
# Mutation entfällt, sobald ihre Fälle bewertet sind, und der eines Falls, sobald der Fall endet.
raeumen() { [ ! -d "$arbeit" ] || rm -rf -- "$arbeit"; return 0; }
trap raeumen EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
: > "$arbeit/fehler"
melde() { echo "$GPN: ROT — $*" >> "$arbeit/fehler"; }

fehlt=0
for t in git make sed; do command -v "$t" > /dev/null 2>&1 || { echo "$GPN: ROT — $t fehlt im Pfad (Voraussetzung)" >&2; fehlt=1; }; done
[ "$fehlt" = 0 ] || exit 1

# Pfade der Welten: $sys trägt die echten Werkzeuge (ohne ssh und ohne make — die sind Attrappen in der Welt),
# $arbeit/ohne-git und $arbeit/ohne-make dieselbe Menge ohne ein Werkzeug (für die Fälle, in denen es fehlt).
sys=$arbeit/sys
mkdir -p "$sys" "$arbeit/ohne-git" "$arbeit/ohne-make"
for t in bash sh git git-receive-pack git-upload-pack make sed grep awk cat ls rm mkdir rmdir printf env find \
  sort cksum head tail tr wc date sleep kill timeout seq cp mv ln chmod id true od dirname basename \
  mktemp cut; do
  p=$(command -v "$t") || continue
  ln -s "$p" "$sys/$t"
  [ "$t" = git ] || ln -s "$p" "$arbeit/ohne-git/$t"
  [ "$t" = make ] || ln -s "$p" "$arbeit/ohne-make/$t"
done
cat > "$arbeit/attrappe" <<'ATTRAPPE'
#!/usr/bin/env bash
# Attrappe von ssh, make, docker, u-boot und devcontainer; git und rm sind das echte Werkzeug mit Protokoll.
n=${0##*/}
{ printf '%s' "$n"; (( $# == 0 )) || printf ' %s' "${@//$'\n'/ }"; echo; } >> "$AT_LOG"
for v in CLAUDE_CODE_OAUTH_TOKEN ANTHROPIC_API_KEY GIT_TOKEN; do [ -z "${!v+x}" ] || echo "ZUGANG $n $v" >> "$AT_LOG"; done
[ -z "${GIT_SSH_COMMAND+x}" ] || echo "UMGEBUNG $n GIT_SSH_COMMAND=$GIT_SSH_COMMAND" >> "$AT_LOG"
warten() { mkdir -p "$AT_STATE"; echo $$ > "$AT_STATE/pid"; echo $PPID > "$AT_STATE/ppid"; exec sleep 20; }
# Im Container liegen Bare-Repo und Klon unter /workspaces/…; in der Welt liegen sie unter $AT_WELT/workspaces/.
# Die Attrappe setzt deshalb den Pfad, den git in die Konfiguration des Klons schrieb, auf die Form des
# Containers zurück und für einen fetch, der seine Adresse aus der Konfiguration nimmt, vorher auf die Welt.
spiegel() {
  local f
  for f in "$AT_WELT"/workspaces/*/pgwire-recorder/.git/config; do
    [ -f "$f" ] || continue
    # Nur schreiben, wenn der Pfad wirklich die andere Form trägt: ein `sed -i` mit gleichem Inhalt
    # ändert die Änderungszeit, und die Momentaufnahme der Welt (snap) sieht sie.
    if [ "$1" = welt ]; then grep -q "= /workspaces/" "$f" && sed -i "s|= /workspaces/|= $AT_WELT/workspaces/|" "$f"
    else grep -q "= $AT_WELT/workspaces/" "$f" && sed -i "s|= $AT_WELT/workspaces/|= /workspaces/|" "$f"; fi
  done
  return 0
}
case $n in
  git | rm) exec "$AT_ECHT/$n" "$@" ;;
  make)
    echo "MAKE $PWD $*" >> "$AT_LOG"
    [[ -n ${AT_LAUT-} && $* == ${AT_LAUT-} ]] && { echo "WERKZEUG-STDOUT"; echo "WERKZEUG-STDERR" >&2; }
    [[ -n ${AT_MAKE_WARTEN-} && $* == ${AT_MAKE_WARTEN-} ]] && warten
    [[ -n ${AT_MAKE_ROT-} && $* == ${AT_MAKE_ROT-} ]] && { echo "make: $* schlug fehl" >&2; exit 2; }
    exit 0 ;;
  ssh)
    cmd=() host=
    while (( $# )); do
      case $1 in
        -o | -i | -p | -F | -l | -b | -c | -E | -J | -L | -R | -D | -W | -e | -m | -Q | -S | -w | -B | -I) shift 2 ;;
        -*) shift ;;
        *) if [ -z "$host" ]; then host=$1; else cmd+=("$1"); fi; shift ;;
      esac
    done
    [ -n "$host" ] && (( ${#cmd[@]} )) || { echo "ssh: Rechner oder Befehl fehlt im Aufruf" >&2; exit 255; }
    c=${cmd[*]}; c=${c//\/workspaces\//$AT_WELT/workspaces/}
    [[ -n ${AT_SSH_LAUT-} && $c == ${AT_SSH_LAUT-} ]] && echo "SSH-WARNUNG" >&2
    [[ -n ${AT_SSH_WARTEN-} && $c == ${AT_SSH_WARTEN-} ]] && warten
    [[ -n ${AT_SSH_ROT-} && $c == ${AT_SSH_ROT-} ]] && { echo "${AT_SSH_MELDUNG-ssh: connect to host 127.0.0.1 port 2222: Connection refused}" >&2; exit "${AT_SSH_CODE-255}"; }
    cd "$AT_WELT" || exit 1
    export PATH=${AT_FERN_PATH:-$PATH}
    [[ $c == *-pack\ * ]] && exec sh -c "$c"
    ein=
    while IFS= read -r z; do printf 'EIN %s\n' "$z" >> "$AT_LOG"; ein+=$z$'\n'; done
    [[ $c == *fetch* ]] && spiegel welt
    printf '%s' "$ein" | sh -c "$c"
    r=$?
    spiegel container
    exit $r ;;
esac
exit 0
ATTRAPPE
chmod 755 "$arbeit/attrappe"

# --- Rahmen der Fälle ---------------------------------------------------------------------------
pn=pgr-ce-dev
BP=/workspaces/$pn/austausch.git
KP=/workspaces/$pn/pgwire-recorder
IDN='Die Testerin' IDM=testerin@example.org
root=0; [ "$(id -u)" != 0 ] || root=1
SEL='.*' fn= basis= sk= baum= gpk= nf=0 GRUEN=() ROOTWEG=()
declare -A ROT=()
# GP_ENTDECKE nennt zusätzlich den Fall, in dem ein roter Name gesetzt wurde (die Namen p4-* setzt ex()
# aus jedem Fall mit Exit 2 bis 4; ohne diese Zeile ist nicht zu sehen, welcher es war).
rot() { ROT[${1:-$fn}]=1; [ -z "${GP_ENTDECKE-}" ] || echo "  rot ${1:-$fn} (im Fall $fn, Exit ${code-?})" >&2; }
FN=() FC=()
fall() { FN+=("$1"); FC+=("$2"); }
sel() { [[ $1 =~ ^($SEL)$ ]]; }

welt() {
  local t
  nf=$((nf + 1)); w=$basis/f$nf
  dev=$w/$pn weltdir=$w/welt repo=$w/repo home=$w/home
  bp=$weltdir/workspaces/$pn/austausch.git kp=$weltdir/workspaces/$pn/pgwire-recorder
  mkdir -p "$w/bin" "$home/.ssh" "$dev/ssh" "$weltdir/workspaces/$pn" "$repo"
  : > "$w/log"; echo hostschluessel > "$home/.ssh/known_hosts"
  printf 'schemaVersion: 1\nproject:\n    name: %s\n' "$pn" > "$dev/u-boot.yaml"
  echo schluessel > "$dev/ssh/id_ed25519"
  for t in ssh make git rm docker u-boot devcontainer id; do ln -s "$arbeit/attrappe" "$w/bin/$t"; done
  git init -q --initial-branch=main "$repo"
  git -C "$repo" config user.name "$IDN"
  git -C "$repo" config user.email "$IDM"
  git -C "$repo" commit -q --allow-empty -m erst
  zw=thema cwd=$repo pfad= nodev= nozw= lauf_cmd=()
}
# snap und snapw vergleichen Pfad, Rechte, Art, Link-Ziel, Größe und Änderungszeit (Nanosekunden) — alles
# aus EINEM find-Aufruf ohne weiteren Prozess je Datei. Ein `-exec cksum {} +` davor starb unter Last
# gelegentlich an einem fehlgeschlagenen fork, und die fehlende Zeile färbte p4-nichts-geaendert rot
# (lastabhängiger Fehlbefund, BEO-REPO/test-faellt-unter-cpu-last-aus). Die Änderungszeit eines
# Verzeichnisses bleibt außen vor (git schreibt in .git, und der Vergleich gilt den Dateien), ebenso
# `.git/index`: `git status --porcelain` im Lesebefehl schreibt ihn neu, wenn die Stat-Daten im selben
# Sekundenfenster liegen (racily clean) — das ist eine Eigenschaft des gerufenen `git`, kein Zustand, den
# Punkt 4 meint. Ebenso außen vor bleibt jede Datei `*.lock`: `git` startet nach einem Commit seine
# Wartung im Hintergrund (`.git/objects/maintenance.lock`) und entfernt die Sperre wieder; unter Last liegt
# sie in der ersten Momentaufnahme und in der zweiten nicht mehr (gemessen am 2026-10-21).
snap() { (cd "$w" 2> /dev/null && find . ! -path ./bin ! -path './bin/*' ! -path ./st ! -path './st/*' ! -path ./kind ! -path './kind/*' \
  ! -path '*/.git/index' ! -name '*.lock' ! -name log ! -name 'k*' ! -name out ! -name err \( -type f -printf '%p %m %y %l %s %T@\n' -o -printf '%p %m %y %l\n' \) 2> /dev/null) | sort; }
snapw() { local d=$1; shift; (cd "$d" && find . ! -path '*/.git/index' ! -name '*.lock' "$@" \( -type f -printf '%p %m %y %l %s %T@\n' -o -printf '%p %m %y %l\n' \) 2> /dev/null) | sort; }
lauf() {
  local b=(PATH="${pfad:-$w/bin:$sys}" HOME="$home" AT_LOG="$w/log" AT_STATE="$w/st" AT_WELT="$weltdir" AT_ECHT="$sys" \
    GIT_CONFIG_GLOBAL="$home/.gitconfig" GIT_CONFIG_SYSTEM="$home/gitsystem") cmd=("$BASH" "$sk")
  [ -n "$nodev" ] || b+=(PGR_DEV="$dev")
  [ -n "$nozw" ] || b+=(PGR_ZWEIG="$zw")
  (( ${#lauf_cmd[@]} == 0 )) || cmd=("${lauf_cmd[@]}")
  : > "$w/log"; s0=$(snap)
  (cd "$cwd" && timeout 120 env -i "${b[@]}" "$@" "${cmd[@]}" > "$w/out" 2> "$w/err"); code=$?
}
gruen() { lauf "$@"; [ "$code" = 0 ] || rot; grep -qxF "container-einrichten: eingerichtet $KP $zw" "$w/out" || rot; }
# ex <Exit-Code> [<Text>…]: Exit-Code und Texte auf stderr; bei Exit 2 bis 4 gelten Punkt 4 (nichts geändert,
# stdout leer, von git nur lesende Kommandos, höchstens ein ssh-Aufruf) und die Form der Meldung.
ex() {
  local x=$1 t
  shift
  [ "$code" = "$x" ] || rot
  for t; do grep -qF -- "$t" "$w/err" || rot; done
  [[ $code =~ ^[234]$ ]] || return 0
  if sel p4-nichts-geaendert; then
    local sn; sn=$(snap)
    { [ "$s0" = "$sn" ] && [ ! -s "$w/out" ] && [ "$(grep -c '^ssh ' "$w/log")" -le 1 ] \
      && ! grep -qvE '^ssh |^git (rev-parse|config --get|remote get-url|check-ref-format|-C [^ ]+ (rev-parse|remote get-url|status|for-each-ref|symbolic-ref))' "$w/log"; } \
      || { rot p4-nichts-geaendert
           [ -z "${GP_ENTDECKE-}" ] || { echo "  p4-nichts-geaendert in $fn: $(diff <(printf '%s\n' "$s0") <(printf '%s\n' "$sn") | head -4 | tr '\n' '|')" >&2; } }
  fi
  sel p4-meldung-form && { ! grep -qvE '^container-einrichten: (eingaben|werkzeuge|zustand): ' "$w/err" || rot p4-meldung-form; }
  return 0
}
zeilen() { [ "$(wc -l < "$w/err")" = "$1" ] || rot; }
ein() { local nm=$1 tx=$2; shift 2; lauf "$@"; ex 2 "eingaben: $nm" "$tx"; zeilen 1; }
log() { grep -qE -- "$1" "$w/log" || rot; }
# folge: die Folge der treibenden Aufrufe (ein ssh-Aufruf je Schritt, der Push am Host) in einer Zeile
folge() { sed -nE '/^git push /{s/.*/push/p;b}; /^ssh /{s/.*(bare=fehlt|init --bare|git clone|config --local user.name|hooks-install|switch -c|traeger-fetch).*/\1/p}' "$w/log" | tr '\n' ' '; }
nicht() { ! grep -qE -- "$1" "$w/log" || rot; }
# Welt des Containers: Bare-Repo, gefüllter main, Klon mit dem Pfad des Containers als origin.
wbare() { git init -q --bare --initial-branch="${1:-main}" "$bp"; }
wpush() { git -C "$repo" push -q "$bp" refs/heads/main:refs/heads/main; }
wklon() { wbare; wpush; git clone -q "$bp" "$kp"; git -C "$kp" remote set-url origin "$BP"; }
# sig <Signal> <Code> <Muster> <Schritt>: Signal, sobald das Kind des Schrittes wartet (Punkt 15).
sig() {
  local p
  ( for p in $(seq 300); do [ -s "$w/st/ppid" ] && break; sleep 0.05; done
    kill -"$1" "$(cat "$w/st/ppid" 2> /dev/null)" 2> /dev/null ) &
  lauf AT_SSH_WARTEN="$3"
  wait
  [ "$code" = "$2" ] || rot
  grep -qxF "container-einrichten: $4: abgebrochen" "$w/err" || rot
  ! kill -0 "$(cat "$w/st/pid" 2> /dev/null)" 2> /dev/null || rot
}
# mt: Temp-Baum mit Makefile, Fragmenten und Skript, selbst ein Host-Repo mit main (für die Aufrufarten von make).
mt() {
  mkdir -p "$w/t/harness/mk" "$w/t/tools/harness"
  cp "$baum/Makefile" "$w/t/"
  cp "$baum"/harness/mk/container-einrichten.mk "$baum"/harness/mk/container-einrichten-gegenprobe.mk \
    "$baum"/harness/mk/gates-schnell.mk "$w/t/harness/mk/"
  cp "$sk" "$w/t/tools/harness/container-einrichten.sh"
  git init -q --initial-branch=main "$w/t"
  git -C "$w/t" config user.name "$IDN"; git -C "$w/t" config user.email "$IDM"
  git -C "$w/t" commit -q --allow-empty -m t
  cwd=$w/t
}
# mk [<Flag>…] [<ZUWEISUNG>…]: das Ziel über das echte make im Temp-Baum (in der Welt ist `make` eine Attrappe)
mk() { local fl=(); while [[ ${1-} == -* ]]; do fl+=("$1"); shift; done; lauf_cmd=("$sys/make" "${fl[@]}" container-einrichten); lauf "$@"; lauf_cmd=(); }
# dbank <Variable>: der Wert einer Make-Variablen im Temp-Baum
dbank() { make -C "$w/t" -pq help 2> /dev/null | awk -v z="$1" '$1 == z && ($2 == ":=" || $2 == "=") { $1 = ""; $2 = ""; print; exit }'; }
# vorauss: die Zahl der Wörter hinter dem Namen in jeder Zeile `container-einrichten:` der Datenbank von make
# im Temp-Baum, je Zeile eine Zahl (0: keine Voraussetzung, auch keine der Form `| …`)
vorauss() { make -C "$w/t" -pq help 2> /dev/null | awk '$1 == "container-einrichten:" { print NF - 1 }' | sort -u | tr '\n' ' '; }
# gp_kind <GP_MUT> [<sed-Ausdruck>…] -- [<ZUWEISUNG>…]: eine Kopie dieser Gegenprobe in eigenem Arbeitsordner
# ($kw); ihre Quellen sind die der Mutation, ihre Fälle p17-* sind abgeschaltet (GP_KIND).
gp_kind() {
  local s=() m=$1
  shift
  while [ "${1-}" != -- ]; do s+=(-e "$1"); shift; done
  shift
  kw=$w/kind; mkdir -p "$kw"
  cp "$gpk" "$w/kgp.sh"
  (( ${#s[@]} == 0 )) || sed -i "${s[@]}" "$w/kgp.sh"
  (cd "$w" && env -i PATH="$sys" HOME="$w" TMPDIR="$kw" GP_KIND=1 GP_PARALLEL=1 GP_QUELLE="$baum" GP_MUT="$m" \
    GP_NUR=p2-vorgaben "$@" "$BASH" "$w/kgp.sh" > "$w/kout" 2> "$w/kerr"); kcode=$?
}
# gp_sig <Signal> <Code>: eine Kopie dieser Gegenprobe mit einer Warteschleife, die das Signal annimmt.
gp_sig() {
  local p i
  kw=$w/kind; mkdir -p "$kw"
  cp "$gpk" "$w/kgp.sh"
  sed -i 's|^fehlt=0$|for _i in $(seq 400); do sleep 0.05; done; fehlt=0|' "$w/kgp.sh"
  ( cd "$w" && exec env -i PATH="$sys" HOME="$w" TMPDIR="$kw" GP_KIND=1 GP_QUELLE="$baum" GP_MUT='^m-p2-vorgabe-remote$' \
    GP_NUR=p2-vorgaben "$BASH" "$w/kgp.sh" > "$w/kout" 2> "$w/kerr" ) & p=$!
  for i in $(seq 100); do [ -n "$(ls -A "$kw" 2> /dev/null)" ] && break; sleep 0.05; done
  kill -"$1" "$p" 2> /dev/null
  wait "$p"; kcode=$?
  [ "$kcode" = "$2" ] || rot
  ! grep -q ': Zeit ' "$w/kerr" || rot
  [ -z "$(ls -A "$kw" 2> /dev/null)" ] || rot
}

# --- Fälle --------------------------------------------------------------------------------------
# Punkt 1: Gegenstand und Ort — je Element der zwei Aufzählungen ein Fall.
fall p1-sshcommand 'gruen; { [ -z "$(git -C "$repo" config --local --get core.sshCommand)" ] && ! grep -q core.sshCommand "$w/log"; } || rot'
fall p1-andere-einstellung 'a=$(git -C "$repo" config --local --list | grep -v "^remote\."); gruen; [ "$a" = "$(git -C "$repo" config --local --list | grep -v "^remote\.")" ] || rot'
fall p1-ssh-verzeichnis 'a=$(snapw "$home/.ssh"); gruen; [ "$a" = "$(snapw "$home/.ssh")" ] || rot'
fall p1-known-hosts 'a=$(cat "$home/.ssh/known_hosts"); gruen; [ "$a" = "$(cat "$home/.ssh/known_hosts")" ] || rot'
fall p1-arbeitsbaum 'echo quelle > "$repo/datei.txt"; a=$(snapw "$repo" ! -path "./.git/*"); gruen; [ "$a" = "$(snapw "$repo" ! -path "./.git/*")" ] || rot'
fall p1-global-host 'a=$(snapw "$home" -maxdepth 1); gruen; [ "$a" = "$(snapw "$home" -maxdepth 1)" ] || rot'
fall p1-global-aufruf 'gruen; nicht "config .*--(global|system)"'
fall p1-pgr-dev 'a=$(snapw "$dev"); gruen; [ "$a" = "$(snapw "$dev")" ] || rot'
NIE=('docker:^docker ' 'u-boot:^u-boot ' 'devcontainer:^devcontainer ' 'push-force:^git push .*--force( |$)'
  'push-force-with-lease:--force-with-lease'
  'push-mirror:^git push .*--mirror' 'push-delete:^git push .*--delete' 'reset:^git (-C [^ ]+ )?reset'
  'checkout-f:^git (-C [^ ]+ )?checkout .*-f' 'clean:^git (-C [^ ]+ )?clean' 'merge:^git (-C [^ ]+ )?merge'
  'pull:^git (-C [^ ]+ )?pull' 'rebase:^git (-C [^ ]+ )?rebase' 'remote-remove:^git remote remove'
  'remote-set-url:^git remote set-url' 'branch-d:^git (-C [^ ]+ )?branch -D' 'rm-rf:^rm .*-rf')
for i in "${!NIE[@]}"; do fall "p1-nie-${NIE[i]%%:*}" "gruen; nicht '${NIE[i]#*:}'"; done

# Punkt 2: Eingaben — je Form ein Fall; alle Verstöße in der Reihenfolge der Tabelle.
fall p2-dev-nicht-gesetzt 'nodev=1; lauf; ex 2 "PGR_DEV ist nicht gesetzt"; zeilen 1'
fall p2-dev-leer 'ein PGR_DEV "ist leer" PGR_DEV='
fall p2-dev-relativ 'ein PGR_DEV "ist nicht absolut" PGR_DEV=rel/pgr-ce-dev'
fall p2-dev-metazeichen 'mkdir -p "$w/a;id/pgr-ce-dev"; ein PGR_DEV "Zeichen außerhalb" "PGR_DEV=$w/a;id/pgr-ce-dev"
  nicht "^id"'
fall p2-dev-leerzeichen 'mkdir -p "$w/a b/pgr-ce-dev"; ein PGR_DEV "Zeichen außerhalb" "PGR_DEV=$w/a b/pgr-ce-dev"'
fall p2-dev-fehlt 'ein PGR_DEV "nirgends fehlt" PGR_DEV=$w/nirgends'
fall p2-dev-datei 'ein PGR_DEV "ist kein Verzeichnis" PGR_DEV=$dev/u-boot.yaml'
for p in leer: abc:abc 02222:02222 plus22:+22 null:0 ueber:65536; do
  fall "p2-port-${p%%:*}" "ein PGR_SSH_PORT 'keine ganze Zahl von 1 bis 65535' PGR_SSH_PORT=${p#*:}"
done
for p in 1 65535; do fall "p2-port-grenze-$p" "gruen PGR_SSH_PORT=$p; log 'ssh .*-p $p vscode@'"; done
fall p2-remote-leer 'ein PGR_REMOTE "ist leer, beginnt nicht" PGR_REMOTE='
fall p2-remote-slash 'ein PGR_REMOTE "Zeichen außerhalb" PGR_REMOTE=a/b'
fall p2-remote-leerzeichen 'ein PGR_REMOTE "Zeichen außerhalb" "PGR_REMOTE=a b"'
fall p2-remote-minus 'ein PGR_REMOTE "beginnt nicht mit A-Za-z0-9" PGR_REMOTE=-x'
fall p2-remote-65 'ein PGR_REMOTE "länger als 64" "PGR_REMOTE=$(printf a%.0s $(seq 65))"'
fall p2-remote-64 'gruen "PGR_REMOTE=$(printf a%.0s $(seq 64))"'
fall p2-zweig-nicht-gesetzt 'nozw=1; lauf; ex 2 "PGR_ZWEIG ist nicht gesetzt"; zeilen 1'
fall p2-zweig-leer 'ein PGR_ZWEIG "ist leer" PGR_ZWEIG='
fall p2-zweig-main 'ein PGR_ZWEIG "main ist nicht zugelassen" PGR_ZWEIG=main'
fall p2-zweig-punkte 'ein PGR_ZWEIG "kein Zweigname" PGR_ZWEIG=a..b'
fall p2-zweig-minus 'ein PGR_ZWEIG "beginnt mit -" PGR_ZWEIG=-x'
fall p2-zweig-leerzeichen 'ein PGR_ZWEIG "kein Zweigname" "PGR_ZWEIG=x y"'
fall p2-zweig-refs 'ein PGR_ZWEIG "beginnt mit refs/" PGR_ZWEIG=refs/heads/x'
fall p2-zweig-a 'zw=a; gruen'
fall p2-zweig-kennung 'zw=slice-harness-devcontainer-einrichten; gruen'
fall p2-opts-nicht-gesetzt 'gruen; log "ConnectTimeout=5 -p 2222 vscode@"'
fall p2-opts-leer 'gruen PGR_SSH_OPTS=; log "ConnectTimeout=5 -p 2222 vscode@"'
fall p2-opts-eine 'gruen PGR_SSH_OPTS=-v; log "ConnectTimeout=5 -v -p 2222"'
fall p2-opts-zwei 'gruen "PGR_SSH_OPTS=-v -q"; log "ConnectTimeout=5 -v -q -p 2222"'
fall p2-opts-tabulator 'gruen "PGR_SSH_OPTS=$(printf -- "-v\t-q")"; log "ConnectTimeout=5 -v -q -p 2222"'
fall p2-opts-zeilenumbruch 'ein PGR_SSH_OPTS "Steuerzeichen" "PGR_SSH_OPTS=$(printf -- "-v\n-q")"'
# Die Zusage "nie von einer Shell ausgewertet" gilt den ssh-Aufrufen des Ziels (Punkt 2); der Lauf endet hier
# vor `einliefern`, weil git die Zeichenkette GIT_SSH_COMMAND aus Punkt 7 notwendig an eine Shell gibt.
fall p2-opts-keine-shell 'wbare master; lauf "PGR_SSH_OPTS=-oSendEnv=\$(docker)"; ex 4 "HEAD auf refs/heads/master"
  nicht "^docker"; log "ConnectTimeout=5 -oSendEnv=.[(]docker[)] -p 2222"'
# Ein Muster im Wert bleibt ein Argument: keine Pfadnamen-Expansion (Punkt 2).
# Der Lauf endet vor `einliefern`, weil git die Zeichenkette aus Punkt 7 an eine Shell gibt (Punkt 20 (7)).
fall p2-opts-muster ': > "$repo/kh-a"; : > "$repo/kh-b"; wbare master; lauf "PGR_SSH_OPTS=-F kh-*"
  ex 4 "HEAD auf refs/heads/master"; log "ConnectTimeout=5 -F kh-[*] -p 2222"'
fall p2-vorgaben 'gruen; log "ssh .*-p 2222 vscode@127.0.0.1"; [ "$(git -C "$repo" remote get-url container)" = "ssh://vscode@127.0.0.1:2222$BP" ] || rot'
fall p2-alle-gemeldet 'lauf PGR_DEV=rel PGR_SSH_PORT=0 PGR_REMOTE=-x PGR_ZWEIG=main "PGR_SSH_OPTS=$(printf -- "a\nb")"; ex 2; zeilen 5'
fall p2-reihenfolge 'lauf PGR_DEV=rel PGR_SSH_PORT=0 PGR_REMOTE=-x PGR_ZWEIG=main "PGR_SSH_OPTS=$(printf -- "a\nb")"; ex 2
  [ "$(sed -n "s/^container-einrichten: eingaben: \(PGR_[A-Z_]*\) .*/\1/p" "$w/err" | tr "\n" " ")" = "PGR_DEV PGR_SSH_PORT PGR_REMOTE PGR_ZWEIG PGR_SSH_OPTS " ] || rot'

# Punkt 3: Werkzeuge am Host und Projektname aus u-boot.yaml.
fall p3-fehlt-ssh 'rm "$w/bin/ssh"; lauf; ex 3 "werkzeuge: ssh fehlt im Pfad"; zeilen 1'
fall p3-fehlt-git 'rm "$w/bin/git"; pfad=$w/bin:$arbeit/ohne-git; lauf; ex 3 "werkzeuge: git fehlt im Pfad"; zeilen 1'
fall p3-fehlt-beide 'rm "$w/bin/ssh" "$w/bin/git"; pfad=$w/bin:$arbeit/ohne-git; lauf; ex 3
  [ "$(sed "s/^container-einrichten: werkzeuge: //;s/ fehlt im Pfad$//" "$w/err" | tr "\n" " ")" = "ssh git " ] || rot'
fall p3-fehlt-git-zweig-ungueltig 'rm "$w/bin/git"; pfad=$w/bin:$arbeit/ohne-git; lauf PGR_ZWEIG=a..b; ex 3 "git fehlt im Pfad"'
fall p3-name-datei-fehlt 'rm "$dev/u-boot.yaml"; lauf; ex 4 "u-boot.yaml fehlt" "kein Gerüst" "make container-start"'
fall p3-name-zeile-fehlt 'printf "schemaVersion: 1\n" > "$dev/u-boot.yaml"; lauf; ex 4 "trägt keine Zeile" "kein Gerüst" "make container-start"'
fall p3-name-gross 'printf "project:\n    name: Pgr-ce-dev\n" > "$dev/u-boot.yaml"; lauf; ex 4 "u-boot.yaml nennt den Namen" "Pgr-ce-dev" "außerhalb von a-z0-9-"'
fall p3-name-unterstrich 'printf "project:\n    name: pgr_ce_dev\n" > "$dev/u-boot.yaml"; lauf; ex 4 "außerhalb von a-z0-9-"'
fall p3-name-ungleich-basis 'printf "project:\n    name: andersrum\n" > "$dev/u-boot.yaml"; lauf; ex 4 "nennt den Namen andersrum, nicht den Basisnamen $pn"'
fall p3-name-zwei-zeilen 'printf "project:\n    name: %s\n    name: andersrum\n" "$pn" > "$dev/u-boot.yaml"; gruen'
fall p3-name-anderer-einzug 'printf "project:\n  name: %s\n" "$pn" > "$dev/u-boot.yaml"; lauf; ex 4 "trägt keine Zeile"'

# Punkt 4: Ablauf, Meldungen, Ströme. p4-nichts-geaendert und p4-meldung-form prüft ex() in jedem Fall mit
# Exit 2 bis 4; sie haben keinen eigenen Aufbau.
fall p4-stdout-nur-erfolg 'wklon; git -C "$kp" branch --no-track thema origin/main; gruen; [ "$(wc -l < "$w/out")" = 1 ] || rot'
fall p4-werkzeug-stdout 'gruen AT_LAUT=hooks-install; grep -q WERKZEUG-STDOUT "$w/out" || rot'
fall p4-werkzeug-stderr 'gruen AT_LAUT=hooks-install; grep -q WERKZEUG-STDERR "$w/err" || rot'
fall p4-werkzeug-stderr-ssh 'gruen "AT_SSH_LAUT=*bare=fehlt*"; grep -qx SSH-WARNUNG "$w/err" || rot'
fall p4-reihenfolge-der-schritte 'gruen
  [ "$(folge)" = "bare=fehlt init --bare push git clone config --local user.name hooks-install switch -c traeger-fetch " ] || rot'

# Punkt 5: Zustand, ein Lesezug, alles vor der ersten Änderung.
fall p5-schluessel-fehlt 'rm "$dev/ssh/id_ed25519"; lauf; ex 4 "ssh/id_ed25519 fehlt"'
fall p5-schluessel-unlesbar 'chmod 000 "$dev/ssh/id_ed25519"; lauf; ex 4 "ssh/id_ed25519 ist nicht lesbar"'
fall p5-kein-repo 'cwd=$w; lauf; ex 4 "liegt nicht in einem Git-Repo"'
fall p5-kein-main 'git -C "$repo" branch -m haupt; lauf; ex 4 "refs/heads/main fehlt im Host-Repo"'
fall p5-detached 'git -C "$repo" checkout -q --detach; gruen'
fall p5-anderer-zweig 'git -C "$repo" switch -q -c nebenbei; gruen'
fall p5-ident-name-fehlt 'git -C "$repo" config --unset user.name; lauf; ex 4 "Identität des Host-Repos fehlt: user.name"'
fall p5-ident-mail-fehlt 'git -C "$repo" config --unset user.email; lauf; ex 4 "fehlt: user.email"'
fall p5-ident-beide-fehlen 'git -C "$repo" config --unset user.name; git -C "$repo" config --unset user.email
  lauf; ex 4 "fehlt: user.name user.email"; zeilen 1'
fall p5-ident-leer 'git -C "$repo" config user.name ""; lauf; ex 4 "fehlt: user.name"'
fall p5-ident-leerraum 'printf "[user]\n\tname = \"  \"\n" >> "$repo/.git/config"; lauf; ex 4 "fehlt: user.name"'
fall p5-ident-global 'git -C "$repo" config --unset user.name; git config --file "$home/.gitconfig" user.name "$IDN"
  gruen; [ "$(git -C "$kp" config --local --get user.name)" = "$IDN" ] || rot'
for p in tabulator:'"A\\tB"' zeilenumbruch:'"A\\nB"' cr:'"A\rB"'; do
  fall "p5-ident-${p%%:*}" "printf '[user]\n\tname = ${p#*:}\n' >> \"\$repo/.git/config\"; lauf; ex 4 'user.name des Host-Repos trägt ein Steuerzeichen'
    ! grep -q 'A.B' \"\$w/err\" || rot"
done
fall p5-ident-nicht-utf8 'printf "[user]\n\tname = A\xffB\n" >> "$repo/.git/config"; lauf; ex 4 "user.name des Host-Repos trägt ein Steuerzeichen oder kein gültiges UTF-8"
  { ! grep -q "$(printf "A\xffB")" "$w/err" && ! grep -q "der Wert" "$w/out"; } || rot'
fall p5-ident-mail-nicht-utf8 'printf "[user]\n\temail = A\xffB\n" >> "$repo/.git/config"; lauf; ex 4 "user.email des Host-Repos trägt ein Steuerzeichen"
  ! grep -q "$(printf "A\xffB")" "$w/err" || rot'
fall p5-ident-utf8-gut 'git -C "$repo" config user.name "Käthe 日 😀"; gruen'
fall p5-remote-andere-url 'git -C "$repo" remote add container ssh://vscode@127.0.0.1:9/x.git
  lauf; ex 4 "zeigt auf ssh://vscode@127.0.0.1:9/x.git" "nicht auf ssh://vscode@127.0.0.1:2222$BP"'
fall p5-remote-gleich 'git -C "$repo" remote add container "ssh://vscode@127.0.0.1:2222$BP"; gruen; nicht "remote add"'
fall p5-remote-nichts-gesetzt 'git -C "$repo" remote add container ssh://x/y; a=$(git -C "$repo" remote get-url container)
  lauf; ex 4; [ "$a" = "$(git -C "$repo" remote get-url container)" ] || rot'
fall p5-remote-anderer-name 'gruen PGR_REMOTE=zwei; [ "$(git -C "$repo" remote get-url zwei)" = "ssh://vscode@127.0.0.1:2222$BP" ] || rot'
fall p5-lesezug-ein-aufruf 'gruen; [ "$(grep -c "^ssh .*bare=fehlt" "$w/log")" = 1 ] || rot'
fall p5-lesezug-aufruf 'gruen; log "^ssh -i $dev/ssh/id_ed25519 -o IdentitiesOnly=yes -o BatchMode=yes -o ConnectTimeout=5 -p 2222 vscode@127.0.0.1 "'
fall p5-lesezug-kein-hostschluessel 'gruen; nicht "StrictHostKeyChecking|UserKnownHostsFile"'
fall p5-lesezug-schreibt-nichts 'wbare master; a=$(snapw "$weltdir"); lauf; ex 4 "HEAD auf refs/heads/master"
  [ "$a" = "$(snapw "$weltdir")" ] || rot'
fall p5-nicht-erreichbar 'lauf "AT_SSH_ROT=*bare=fehlt*"; ex 4 "vscode@127.0.0.1, Port 2222" "Connection refused" "make container-start"'
fall p5-fern-git-fehlt 'lauf AT_FERN_PATH=$arbeit/ohne-git; ex 4 "git fehlt im Pfad des Containers vscode@127.0.0.1:2222"'
fall p5-fern-make-fehlt 'lauf AT_FERN_PATH=$arbeit/ohne-make; ex 4 "make fehlt im Pfad des Containers"'
fall p5-bare-master 'wbare master; lauf; ex 4 "$BP hat HEAD auf refs/heads/master"'
fall p5-bare-fremd 'mkdir -p "$bp"; echo x > "$bp/datei"; lauf; ex 4 "$BP besteht und ist kein Bare-Repo"'
fall p5-klon-schmutzig 'wklon; echo neu > "$kp/datei"; lauf; ex 4 "$KP trägt eine uncommittete Änderung" "committe oder stashe"'
fall p5-klon-unversioniert 'wklon; echo neu > "$kp/unversioniert"; lauf; ex 4 "uncommittete Änderung"'
fall p5-klon-anderer-origin 'wklon; git -C "$kp" remote set-url origin "$w/fremd.git"; lauf; ex 4 "hat origin $w/fremd.git, nicht $BP"'
fall p5-klon-fremd 'mkdir -p "$kp"; echo x > "$kp/datei"; lauf; ex 4 "$KP besteht und ist kein Git-Repo"'

# Punkt 6: Bare-Repo.
fall p6-fehlt 'gruen; log "^ssh .*git init --bare --initial-branch=main $BP$"; [ "$(git -C "$bp" symbolic-ref -q HEAD)" = refs/heads/main ] || rot'
fall p6-vorhanden 'wbare; wpush; gruen; nicht "init --bare"'
fall p6-fehler 'lauf "AT_SSH_ROT=*init --bare*"; ex 5 "bare: git init --bare $BP endete mit Exit 255"; nicht "git push"'

# Punkt 7: Remote und Einlieferung.
fall p7-remote-fehlt 'gruen; log "^git remote add container ssh://vscode@127.0.0.1:2222$BP$"'
# Der Port der URL stammt aus PGR_SSH_PORT (Randform 8: mehrere Instanzen): die URL des Remotes trägt ihn, und
# der Push ruft `ssh` mit genau diesem Port; die Attrappe protokolliert die Argumente, auch das -p, das git setzt.
fall p7-url-port 'gruen PGR_SSH_PORT=2253; [ "$(git -C "$repo" remote get-url container)" = "ssh://vscode@127.0.0.1:2253$BP" ] || rot
  log "^ssh .*-p 2253 vscode@127.0.0.1 git-receive-pack"; nicht "-p 2222"'
fall p7-push-aufruf 'gruen; log "^git push container refs/heads/main:refs/heads/main$"; nicht "git push .*--force"'
fall p7-push-angekommen 'gruen; [ "$(git -C "$bp" rev-parse main)" = "$(git -C "$repo" rev-parse main)" ] || rot'
fall p7-git-ssh-command 'gruen; log "^UMGEBUNG ssh GIT_SSH_COMMAND=ssh -i $dev/ssh/id_ed25519 -o IdentitiesOnly=yes -o BatchMode=yes$"
  [ -z "$(git -C "$repo" config --get core.sshCommand)" ] || rot'
fall p7-git-ssh-command-opts 'gruen PGR_SSH_OPTS=-v; log "^UMGEBUNG ssh GIT_SSH_COMMAND=ssh -i .* -o BatchMode=yes -v$"'
fall p7-abgewiesen 'wbare; wpush; git -C "$repo" commit -q --amend --allow-empty -m anders
  lauf; ex 5 "einliefern: der Remote container trägt einen main mit anderer Historie" "git fetch container"
  [ "$(grep -c "^git push" "$w/log")" = 1 ] || rot; nicht "git clone"'
fall p7-fehler 'lauf "AT_SSH_ROT=*receive-pack*"; ex 5
  grep -qE "^container-einrichten: einliefern: git push container endete mit Exit [0-9]+$" "$w/err" || rot
  nicht "git clone"'

# Punkt 8: Klon.
fall p8-fehlt 'gruen; log "^ssh .*git clone $BP $KP$"; [ -d "$kp/.git" ] || rot
  [ "$(git -C "$kp" remote get-url origin)" = "$BP" ] || rot'
fall p8-besteht 'wklon; gruen; nicht "git clone"; log "fetch --prune origin"'
fall p8-fehler 'wklon; lauf "AT_SSH_ROT=*fetch --prune*"; ex 5 "klon: der Klon $KP endete mit Exit 255"; nicht "config --local"'

# Punkt 9: Identität des Klons.
fall p9-lokal 'gruen; { [ "$(git -C "$kp" config --local --get user.name)" = "$IDN" ] && [ "$(git -C "$kp" config --local --get user.email)" = "$IDM" ]; } || rot
  nicht "config .*--(global|system)"; [ ! -s "$home/.gitconfig" ] || rot'
fall p9-reihenfolge 'gruen; [ "$(grep "^EIN " "$w/log" | head -2 | tr "\n" "|")" = "EIN $IDN|EIN $IDM|" ] || rot'
fall p9-kein-argument 'gruen; nicht "^ssh .*$IDN"'
fall p9-fehler 'lauf "AT_SSH_ROT=*config --local*"; ex 5 "identitaet: die Identität des Klons $KP endete mit Exit 255"; nicht "hooks-install"'
ZEICHEN=('A B' ' A B ' 'A"B' "A'B" 'A`docker gefahr`B' 'A$(docker gefahr)B' 'A\B' 'A%sB' 'A;docker gefahr;B' 'A&&docker gefahr&&B'
  'A|docker gefahr|B' 'A<B' 'A>B' 'A*B' 'A?B' 'A#B' 'A!B' 'A(B' 'A)B' 'A{B' 'A}B' 'A[B' 'A]B' 'A~B' 'AäB' 'A日B' 'A😀B')
ZNAME=(leerzeichen leerraum-am-rand anfuehrungszeichen apostroph backtick dollar backslash prozent semikolon und pipe kleiner groesser
  stern frage raute ausruf klammer-auf klammer-zu geschweift-auf geschweift-zu eckig-auf eckig-zu tilde umlaut cjk emoji)
zeichen() {
  local v=${ZEICHEN[$1]}
  git -C "$repo" config user.name "$v"
  gruen
  [ "$(git -C "$kp" config --local --get user.name)" = "$v" ] || rot
  nicht '^docker '
}
# Die Werte, die eine Shell verändern würde (Punkt 9): sie färben jede Mutation, die den Wert in eine
# Zeichenkette setzt, die eine Shell auswertet.
ZSHELL=(p9-zeichen-anfuehrungszeichen p9-zeichen-apostroph p9-zeichen-backslash p9-zeichen-backtick
  p9-zeichen-dollar p9-zeichen-groesser p9-zeichen-klammer-auf p9-zeichen-klammer-zu p9-zeichen-kleiner
  p9-zeichen-leerraum-am-rand p9-zeichen-pipe p9-zeichen-semikolon p9-zeichen-und)
ZALLE=()
for i in "${!ZEICHEN[@]}"; do fall "p9-zeichen-${ZNAME[i]}" "zeichen $i"; ZALLE+=("p9-zeichen-${ZNAME[i]}"); done

# Punkt 10: Commit-Träger.
fall p10-aufruf 'gruen; log "^MAKE $kp hooks-install$"'
fall p10-zweiter-aufruf 'gruen; gruen; log "^MAKE $kp hooks-install$"'
fall p10-fehler 'lauf AT_MAKE_ROT="hooks-install"; ex 5 "hooks: make hooks-install im Klon endete mit Exit 2"; nicht "switch"'

# Punkt 11: Zweig.
fall p11-fehlt 'gruen; log "switch -c "; [ "$(git -C "$kp" rev-parse --abbrev-ref HEAD)" = thema ] || rot'
fall p11-besteht 'wklon; git -C "$kp" branch thema origin/main; gruen; nicht "switch -c"; log "switch "
  [ "$(git -C "$kp" rev-parse --abbrev-ref HEAD)" = thema ] || rot; nicht "(reset|--force|--discard-changes)"'
fall p11-aus-origin-main 'wklon; git -C "$kp" switch -q -c main-alt; git -C "$kp" branch -f main HEAD~0
  git -C "$repo" commit -q --allow-empty -m zweit; gruen
  [ "$(git -C "$kp" rev-parse thema)" = "$(git -C "$repo" rev-parse main)" ] || rot'
fall p11-sonderzeichen 'zw="x\$(docker)y"; gruen; nicht "^docker"
  [ "$(git -C "$kp" rev-parse --abbrev-ref HEAD)" = "$zw" ] || rot'
fall p11-main-unbewegt 'wklon; a=$(git -C "$kp" rev-parse main); git -C "$repo" commit -q --allow-empty -m zweit
  gruen; [ "$(git -C "$kp" rev-parse main)" = "$a" ] || rot
  [ "$(git -C "$kp" rev-parse origin/main)" != "$a" ] || rot'
fall p11-fehler 'lauf "AT_SSH_ROT=*switch*"; ex 5 "zweig: der Zweig thema im Klon endete mit Exit 255"; nicht "traeger-fetch"'

# Punkt 12: Träger des Harness, letzter Schritt.
fall p12-aufruf 'gruen; log "^MAKE $kp traeger-fetch$"'
fall p12-stelle 'gruen; [ "$(folge | grep -oE "(switch -c|traeger-fetch)" | tr "\n" " ")" = "switch -c traeger-fetch " ] || rot'
fall p12-fehler 'lauf AT_MAKE_ROT="traeger-fetch"; ex 5 "traeger: make traeger-fetch im Klon endete mit Exit 2"
  ! grep -q eingerichtet "$w/out" || rot; [ "$(git -C "$kp" rev-parse --abbrev-ref HEAD)" = thema ] || rot'

# Punkt 13: Erfolg, Wiederholbarkeit, Zugang.
fall p13-erfolg 'gruen; [ "$(grep -c "^container-einrichten:" "$w/out")" = 1 ] || rot'
fall p13-zweiter-aufruf 'gruen; gruen; nicht "(init --bare|git clone|remote add)"; log "fetch --prune"; log "hooks-install"; log "switch "; log "traeger-fetch"'
fall p13-umgebung 'gruen CLAUDE_CODE_OAUTH_TOKEN=sentinel-xyz ANTHROPIC_API_KEY=sentinel-xyz GIT_TOKEN=sentinel-xyz
  nicht "^ZUGANG "; ! grep -rqF sentinel-xyz "$w/out" "$w/err" "$w/log" "$weltdir" "$repo" "$home" || rot'

# Punkt 14: Reihenfolge der Gründe — je Nachbarpaar ein Fall, in dem beide Gründe zutreffen.
fall p14-eingabe-vor-werkzeug 'rm "$w/bin/ssh"; lauf PGR_ZWEIG=main; ex 2 "PGR_ZWEIG main"'
fall p14-werkzeug-vor-geruest 'rm "$w/bin/ssh" "$dev/u-boot.yaml"; lauf; ex 3 "ssh fehlt im Pfad"'
fall p14-geruest-vor-repo 'rm "$dev/u-boot.yaml"; cwd=$w; lauf; ex 4 "kein Gerüst"'
fall p14-repo-vor-identitaet 'git -C "$repo" branch -m haupt; git -C "$repo" config --unset user.name; lauf; ex 4 "refs/heads/main fehlt"'
fall p14-identitaet-vor-remote 'git -C "$repo" config --unset user.name; git -C "$repo" remote add container ssh://x/y; lauf; ex 4 "fehlt: user.name"'
fall p14-remote-vor-lesezug 'git -C "$repo" remote add container ssh://x/y; lauf "AT_SSH_ROT=*bare=fehlt*"; ex 4 "zeigt auf ssh://x/y"'
fall p14-lesezug-vor-werkzeugen 'lauf "AT_SSH_ROT=*bare=fehlt*" AT_FERN_PATH=$arbeit/ohne-git; ex 4 "antwortet nicht"'
fall p14-werkzeuge-vor-bare 'wbare master; lauf AT_FERN_PATH=$arbeit/ohne-make; ex 4 "make fehlt im Pfad des Containers"'
fall p14-bare-vor-klon 'wklon; git -C "$bp" symbolic-ref HEAD refs/heads/master; echo neu > "$kp/datei"; lauf; ex 4 "HEAD auf refs/heads/master"'
fall p14-klon-vor-schritt 'wklon; echo neu > "$kp/datei"; lauf "AT_SSH_ROT=*fetch --prune*"; ex 4 "uncommittete Änderung"'

# Punkt 15: Abbruch.
fall p15-sigint 'sig INT 130 "*bare=fehlt*" zustand'
fall p15-sigterm 'sig TERM 143 "*bare=fehlt*" zustand'
fall p15-sigint-zweig 'sig INT 130 "*switch*" zweig'
fall p15-sigterm-traeger 'sig TERM 143 "*traeger-fetch*" traeger'
fall p15-kein-rueckbau 'sig INT 130 "*switch*" zweig; [ -d "$kp/.git" ] && [ -d "$bp" ] || rot; nicht "(rm |git reset|remote remove)"'
fall p15-wiederholung 'sig INT 130 "*switch*" zweig; gruen; [ "$(git -C "$kp" rev-parse --abbrev-ref HEAD)" = thema ] || rot'

# Punkt 16: Aufrufarten des Treibers.
fall p16-phony 'mt; : > "$w/t/container-einrichten"; mk; [ "$code" = 0 ] || rot; grep -q "eingerichtet" "$w/out" || rot'
# Der Wortlaut von Punkt 16: das Ziel hat keine Voraussetzung, auch keine vorhandene Datei ohne Regel.
fall p16-keine-voraussetzung 'mt; [ "$(vorauss)" = "0 " ] || rot'
for f in n q t; do
  fall "p16-$f" "mt; mk -$f; [ ! -s \"\$w/log\" ] || rot; ! grep -q eingerichtet \"\$w/out\" || rot
    [ -z \"\$(git -C \"\$w/t\" remote)\" ] || rot"
  fall "p16-makeflags-$f" "mt; mk MAKEFLAGS=$f; [ ! -s \"\$w/log\" ] || rot; ! grep -q eingerichtet \"\$w/out\" || rot"
done
fall p16-i 'mt; mk -i "AT_SSH_ROT=*bare=fehlt*"; ! grep -q eingerichtet "$w/out" || rot; grep -q "container-einrichten: zustand: " "$w/err" || rot'
fall p16-makeflags-i 'mt; mk MAKEFLAGS=i "AT_SSH_ROT=*bare=fehlt*"; ! grep -q eingerichtet "$w/out" || rot
  grep -q "container-einrichten: zustand: " "$w/err" || rot'
fall p16-makeflags-k 'mt; mk MAKEFLAGS=k; [ "$code" = 0 ] || rot; grep -q eingerichtet "$w/out" || rot'
fall p16-makeflags-j 'mt; mk MAKEFLAGS=j; [ "$code" = 0 ] || rot; grep -q eingerichtet "$w/out" || rot'
fall p16-n-und-i 'mt; mk -n -i; [ ! -s "$w/log" ] || rot; ! grep -q eingerichtet "$w/out" || rot'
fall p16-k 'mt; mk -k; [ "$code" = 0 ] || rot; grep -q eingerichtet "$w/out" || rot'
fall p16-j 'mt; mk -j4; [ "$code" = 0 ] || rot; grep -q eingerichtet "$w/out" || rot'

# Punkt 17: Listen und Zeitbudget (SPEC-060 Punkte 2 und 3 sowie Punkt 16 (9)).
fall p17-liste-container 'mt; [[ " $(dbank GEGENPROBEN_CONTAINER) " == *" container-einrichten-gegenprobe "* ]] || rot'
fall p17-nicht-host 'mt; [[ " $(dbank GEGENPROBEN_HOST) " != *" container-einrichten-gegenprobe "* ]] || rot'
fall p17-kein-gate 'mt; a=" $(dbank GATE_CHECKS) "; [[ $a != *" container-einrichten "* && $a != *" container-einrichten-gegenprobe "* ]] || rot'
fall p17-keine-ausnahme 'mt; [[ " $(dbank GATE_SCHNELL_AUSNAHME) " != *" container-einrichten-gegenprobe "* ]] || rot'
fall p17-phony 'mt; a=" $(make -C "$w/t" -pq help 2> /dev/null | sed -n "s/^\.PHONY://p" | tr "\n" " ") "
  [[ $a == *" container-einrichten "* ]] || rot
  [[ $a == *" container-einrichten-gegenprobe "* ]] || rot'
fall p17-budget-ueber 'gp_kind "^m-p2-vorgabe-remote$" "s|^fehlt=0$|sleep 2; fehlt=0|" -- GP_BUDGET=1
  { [ "$kcode" = 1 ] && grep -qE "^$GPN: ROT — Zeitbudget: [0-9]+ s über 1 s$" "$w/kerr"; } || rot'
fall p17-budget-unter 'gp_kind "^m-p2-vorgabe-remote$" -- GP_BUDGET=3600
  { [ "$kcode" = 0 ] && grep -qE "^$GPN: Zeit [0-9]+ s von höchstens 3600 s, Last [0-9]+[.][0-9]+$" "$w/kerr" \
    && ! grep -q "ROT — Zeitbudget" "$w/kerr"; } || rot'
for z in 999999 10; do
  fall "p17-budget-zulaessig-$z" "gp_kind '^m-p2-vorgabe-remote\$' -- GP_BUDGET=$z
    { [ \"\$kcode\" = 0 ] && grep -qE \"^\$GPN: Zeit [0-9]+ s von höchstens $z s, Last [0-9]+[.][0-9]+\$\" \"\$w/kerr\" \
      && ! grep -q \"ROT — Zeitbudget\" \"\$w/kerr\"; } || rot"
done
fall p17-budget-null 'gp_kind "^m-p2-vorgabe-remote$" "s|^fehlt=0$|sleep 2; fehlt=0|" -- GP_BUDGET=0
  { [ "$kcode" = 0 ] && ! grep -q "ROT — Zeitbudget" "$w/kerr" && grep -qE "^$GPN: Zeit [0-9]+ s von höchstens 0 s, Last [0-9]+[.][0-9]+$" "$w/kerr"; } || rot'
for p in leer: minus:-1 bruch:1.5 abc:abc 08:08 09:09 00:00 007:007 sieben-stellen:1000000 zwanzig-stellen:99999999999999999999 plus:+5 leerraum-vorn:~5 leerraum-hinten:5~; do
  v=${p#*:}; v=${v//"~"/ }
  fall "p17-budget-ungueltig-${p%%:*}" "gp_kind '^m-p2-vorgabe-remote\$' -- 'GP_BUDGET=$v'
    { [ \"\$kcode\" = 2 ] && grep -qx \"\$GPN: GP_BUDGET: keine ganze Zahl\" \"\$w/kerr\" && [ \$(wc -l < \"\$w/kerr\") = 1 ] \
      && [ -z \"\$(ls -A \"\$kw\")\" ]; } || rot"
done
fall p17-budget-vorgabe-nur-ohne-mut 'gp_kind "^m-p2-vorgabe-remote$" --
  { [ "$kcode" = 0 ] && grep -qE "^$GPN: Zeit [0-9]+ s, Last [0-9]+[.][0-9]+$" "$w/kerr" && ! grep -q "ROT — Zeitbudget" "$w/kerr"; } || rot'
fall p17-budget-last-fehlt 'gp_kind "^m-p2-vorgabe-remote$" "s|^LOADAVG=/proc/loadavg$|LOADAVG=/proc/nirgends-loadavg|" -- GP_BUDGET=3600
  grep -qE "^$GPN: Zeit [0-9]+ s von höchstens 3600 s, Last -$" "$w/kerr" || rot'
fall p17-budget-gruende 'gp_kind "^m-p2-vorgabe-remote$" "s|^fehlt=0$|sleep 2; fehlt=0|" "s|p2-vorgaben$|p2-gibtsnicht|" -- GP_BUDGET=1
  { [ "$kcode" = 1 ] && [ "$(grep -cE "^$GPN: (ROT|Zeit)" "$w/kerr")" = 3 ] \
    && [ "$(grep -nE "^$GPN: ROT — Mutation" "$w/kerr" | cut -d: -f1)" -lt "$(grep -nE "^$GPN: ROT — Zeitbudget" "$w/kerr" | cut -d: -f1)" ] \
    && [ "$(grep -nE "^$GPN: ROT — Zeitbudget" "$w/kerr" | cut -d: -f1)" -lt "$(grep -nE "^$GPN: Zeit " "$w/kerr" | cut -d: -f1)" ]; } || rot'
fall p17-budget-sigint 'gp_sig INT 130'
fall p17-budget-sigterm 'gp_sig TERM 143'
fall p17-ordner-mutation 'gp_kind "^m-p2-vorgabe-remote$" "s|^trap raeumen EXIT$|trap : EXIT|" --
  { [ "$kcode" = 0 ] && [ -z "$(find "$kw" -maxdepth 2 -name "m[0-9]*" -print -quit)" ]; } || rot'
fall p17-ordner-fall 'gp_kind "^m-p2-vorgabe-remote$" "s|^trap raeumen EXIT$|trap : EXIT|" "s|^    rm -rf -- \"\$basis\"$|    :|" --
  { [ "$kcode" = 0 ] && [ -n "$(find "$kw" -maxdepth 3 -name q -print -quit)" ] \
    && [ -z "$(find "$kw" -maxdepth 3 -name "f[0-9]*" -print -quit)" ]; } || rot'

faelle() {
  local i
  ROT=()
  for i in "${!FN[@]}"; do
    fn=${FN[i]}
    sel "$fn" || continue
    [[ -n ${GP_KIND-} && $fn == p17-* ]] && continue
    [[ $root = 1 && $fn = p5-schluessel-unlesbar ]] && continue
    welt
    eval "${FC[i]}" > /dev/null
    rm -rf -- "$w"
  done
  for i in "${!ROT[@]}"; do echo "$i"; done
}

# --- Mutationen ---------------------------------------------------------------------------------
maxjobs=${GP_PARALLEL:-4}
wartejob() { while [ "$(jobs -rp | wc -l)" -ge "$maxjobs" ]; do wait -n || true; done; }
MARKE='# --- Tabelle der Mut'"ationen ---"
declare -A FILE=([s]=tools/harness/container-einrichten.sh [m]=harness/mk/container-einrichten.mk \
  [g]=harness/mk/container-einrichten-gegenprobe.mk [k]=harness/mk/gates-schnell.mk [p]=gp.sh)
vorbereiten() {
  local d=$1
  mkdir -p "$d/tools/harness" "$d/harness/mk"
  cp "$quelle/Makefile" "$d/"
  cp "$quelle/tools/harness/container-einrichten.sh" "$d/tools/harness/"
  cp "$quelle"/harness/mk/container-einrichten.mk "$quelle"/harness/mk/container-einrichten-gegenprobe.mk \
    "$quelle"/harness/mk/gates-schnell.mk "$d/harness/mk/"
  cp "$GP_SELF" "$d/gp.sh"
  sk=$d/tools/harness/container-einrichten.sh baum=$d gpk=$d/gp.sh
}
# mut <Name> <Datei> <Auswahl der Fälle> <alt> <neu> [<roter Fall>…]
mut() {
  local name=$1 key=$2 s=$3 alt=$4 neu=$5 k k2=()
  shift 5
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  # Als root entfällt der Fall mit Dateirechten (faelle überspringt ihn); eine Mutation, deren rote Menge
  # dadurch leer würde, entfällt mit ihm und steht am Ende als übersprungen.
  if (( root && $# )); then
    for k; do [ "$k" = p5-schluessel-unlesbar ] || k2+=("$k"); done
    (( ${#k2[@]} )) || { ROOTWEG+=("$name"); return 0; }
    set -- "${k2[@]}"
  fi
  (( $# )) || GRUEN+=("$name")
  wartejob
  (
    local basis=$arbeit/m$BASHPID c r nach ist soll k
    mkdir -p "$basis"
    vorbereiten "$basis/q"
    f=$basis/q/${FILE[$key]}
    c=$(cat "$f"; echo x); c=${c%x}
    nach=; [ "$key" != p ] || { nach=$'\n'$MARKE${c#*$'\n'$MARKE}; c=${c%%$'\n'$MARKE*}; }
    r=${c//"$alt"/}
    if (( ${#c} - ${#r} != ${#alt} )); then
      melde "Mutation $name: die Ersetzung trifft nicht genau eine Stelle"
      rm -rf -- "$basis"; exit 0
    fi
    printf '%s%s' "${c/"$alt"/"$neu"}" "$nach" > "$f"
    SEL=${GP_NUR:-$s}
    ist=$(faelle | sort | tr '\n' ' ')
    soll=$(for k; do echo "$k"; done | sort -u | tr '\n' ' ')
    [ -z "${GP_ENTDECKE-}" ] || echo "$name => ${ist:-nichts}" >&2
    [ "$ist" = "$soll" ] || melde "Mutation $name: rot sind [${ist:-nichts}], erwartet [${soll:-nichts}]"
    rm -rf -- "$basis"
  ) &
}
if [[ ok =~ ${GP_MUT:-.} ]]; then
  wartejob
  (
    basis=$arbeit/ok
    mkdir -p "$basis"
    vorbereiten "$basis/q"
    SEL=${GP_NUR:-'.*'}
    ist=$(faelle | sort | tr '\n' ' ')
    [ -z "$ist" ] || melde "unveränderte Kopie ist nicht grün: $ist"
    rm -rf -- "$basis"
  ) &
fi

# SPEC-060 Punkt 16 (2) und (3): die Zeile mit Zeit, Obergrenze und Last; eine Überschreitung ist rot, und
# die Reihenfolge der Gründe ist (c) rote Mutation vor (d) Zeitbudget.
abschluss() {
  local T=$((SECONDS - T0)) last=-
  [ ! -r "$LOADAVG" ] || read -r last _ < "$LOADAVG"
  if [ -n "$budget" ] && [ "$budget" != 0 ] && (( T > budget )); then
    echo "$GPN: ROT — Zeitbudget: $T s über $budget s" >&2
    fehler=1
  fi
  if [ -n "$budget" ]; then
    echo "$GPN: Zeit $T s von höchstens $budget s, Last $last" >&2
  else
    echo "$GPN: Zeit $T s, Last $last" >&2
  fi
}

# --- Tabelle der Mutationen ---
A1='p1-.*|p4-.*' A2='p2-.*|p4-.*' A3='p3-.*|p4-.*|p14-.*' A5='p5-.*|p4-.*|p14-.*' A9='p9-.*|p4-.*'
A15='p15-.*' A16='p16-.*'
# Die Fälle von Punkt 17 fahren eine Kopie dieser Gegenprobe und sind teuer; die Auswahl je Mutation nennt
# deshalb nur die Familie, die sie trifft: Listen (am Fragment), Zeitbudget und Aufräumen (am Skript).
A17='p17-(liste-container|nicht-host|kein-gate|keine-ausnahme|phony)' A17B='p17-budget-.*'
A17O='p17-budget-(sigint|sigterm|ungueltig-.*)|p17-ordner-.*'
# Punkt 1: die beiden Aufzählungen
mut m-p1-sshcommand s "$A1" 's_hooks() {' 's_hooks() {
  git config --local core.sshCommand ssh' p1-sshcommand p1-andere-einstellung
mut m-p1-sshcommand-gelesen s "$A1" 's_hooks() {' 's_hooks() {
  git config --get core.sshCommand > /dev/null' p1-sshcommand
mut m-p1-andere-einstellung s "$A1" 's_hooks() {' 's_hooks() {
  git config --local gc.auto 0' p1-andere-einstellung
mut m-p1-ssh-verzeichnis s "$A1" 's_hooks() {' 's_hooks() {
  : > "$HOME/.ssh/config"' p1-ssh-verzeichnis
mut m-p1-known-hosts s "$A1" 's_hooks() {' 's_hooks() {
  echo neu >> "$HOME/.ssh/known_hosts"' p1-known-hosts p1-ssh-verzeichnis
mut m-p1-arbeitsbaum s "$A1" 's_hooks() {' 's_hooks() {
  : > datei.txt' p1-arbeitsbaum
mut m-p1-global-host s "$A1" 's_hooks() {' 's_hooks() {
  git config --global gc.auto 0' p1-global-host p1-global-aufruf
mut m-p1-global-fern s "$A1" 's_hooks() {' 's_hooks() {
  fern "git config --global gc.auto 0"' p1-global-aufruf p1-global-host
mut m-p1-pgr-dev s "$A1" 's_hooks() {' 's_hooks() {
  : > "$PGR_DEV/spur"' p1-pgr-dev
NIEMUT=('docker:docker info' 'u-boot:u-boot --version' 'devcontainer:devcontainer up'
  'push-force:git push --force "$PGR_REMOTE" refs/heads/main' 'push-mirror:git push --mirror "$PGR_REMOTE"'
  'push-delete:git push --delete "$PGR_REMOTE" alt' 'reset:git reset --hard' 'checkout-f:git checkout -f'
  'push-force-with-lease:git push --force-with-lease "$PGR_REMOTE" refs/heads/main'
  'clean:git clean -fd' 'merge:git merge --ff-only main' 'pull:git pull' 'rebase:git rebase main'
  'remote-remove:git remote remove alt' 'remote-set-url:git remote set-url "$PGR_REMOTE" "$url"'
  'branch-d:git branch -D alt' 'rm-rf:rm -rf "$PGR_DEV/nichts"')
for i in "${!NIEMUT[@]}"; do
  k=${NIEMUT[i]%%:*}
  # Ein zusätzlicher Push erscheint auch in der Folge der treibenden Aufrufe (Punkt 4).
  z=(); [[ $k != push-* ]] || z=(p4-reihenfolge-der-schritte)
  mut "m-p1-nie-$k" s "p1-nie-$k|p4-.*" 's_hooks() {' "s_hooks() {
  ${NIEMUT[i]#*:} > /dev/null 2>&1" "p1-nie-$k" ${z[@]+"${z[@]}"}
done
# Punkt 2: Eingaben
mut m-p2-dev-nicht-gesetzt s "$A2" 'if [ -z "${PGR_DEV+x}" ]; then meld "PGR_DEV ist nicht gesetzt"' 'if false; then meld "PGR_DEV ist nicht gesetzt"' p2-dev-nicht-gesetzt
mut m-p2-dev-leer s "$A2" 'elif [ -z "$PGR_DEV" ]; then' 'elif false; then' p2-dev-leer
mut m-p2-dev-relativ s "$A2" 'elif [[ $PGR_DEV != /* ]]; then' 'elif false; then' p2-dev-relativ
mut m-p2-dev-zeichen s "$A2" 'elif [[ ! $PGR_DEV =~ ^/[A-Za-z0-9._+/-]*$ ]]; then' 'elif false; then' p2-dev-metazeichen p2-dev-leerzeichen
mut m-p2-dev-fehlt s "$A2" 'elif [ ! -e "$PGR_DEV" ]; then' 'elif false; then' p2-dev-fehlt
mut m-p2-dev-verzeichnis s "$A2" 'elif [ ! -d "$PGR_DEV" ]; then' 'elif false; then' p2-dev-datei
mut m-p2-port-leer s "$A2" 'PGR_SSH_PORT=${PGR_SSH_PORT-2222}' 'PGR_SSH_PORT=${PGR_SSH_PORT:-2222}' p2-port-leer
mut m-p2-zahl-leer s "$A2" '  [[ -n $1 ]] || return 1' '  :' p2-port-leer
mut m-p2-zahl-ziffern s "$A2" '  [[ $1 != *[!0-9]* ]] || return 1' '  :' p2-port-abc p2-port-plus22
mut m-p2-zahl-plus s "$A2" '  [[ $1 != *[!0-9]* ]] || return 1' '  [[ $1 != *[!0-9+]* ]] || return 1' p2-port-plus22
mut m-p2-zahl-fuehrende-null s "$A2" '  [[ $1 != 0?* ]] || return 1' '  :' p2-port-02222
mut m-p2-zahl-null s "$A2" '  [[ $1 != 0 ]] || return 1' '  :' p2-port-null p2-alle-gemeldet p2-reihenfolge
mut m-p2-zahl-eins s "$A2" '  [[ $1 != 0 ]] || return 1' '  [[ $1 != [01] ]] || return 1' p2-port-grenze-1
mut m-p2-zahl-oben s "$A2" '(( ${#1} == ${#2} )) && [[ ! $1 > $2 ]]' 'true' p2-port-ueber
mut m-p2-zahl-grenze-oben s "$A2" '[[ ! $1 > $2 ]]' '[[ $1 < $2 ]]' p2-port-grenze-65535
mut m-p2-remote-leer s "$A2" 'PGR_REMOTE=${PGR_REMOTE-container}' 'PGR_REMOTE=${PGR_REMOTE:-container}' p2-remote-leer
mut m-p2-remote-slash s "$A2" '^[A-Za-z0-9][A-Za-z0-9._-]*$' '^[A-Za-z0-9][A-Za-z0-9._/-]*$' p2-remote-slash
mut m-p2-remote-leerzeichen s "$A2" '^[A-Za-z0-9][A-Za-z0-9._-]*$' '^[A-Za-z0-9][A-Za-z0-9._[:space:]-]*$' p2-remote-leerzeichen
mut m-p2-remote-minus s "$A2" '^[A-Za-z0-9][A-Za-z0-9._-]*$' '^[A-Za-z0-9-][A-Za-z0-9._-]*$' p2-remote-minus p2-alle-gemeldet p2-reihenfolge
mut m-p2-remote-laenge s "$A2" '(( ${#PGR_REMOTE} > 64 ))' '(( ${#PGR_REMOTE} > 65 ))' p2-remote-65
mut m-p2-remote-laenge-63 s "$A2" '(( ${#PGR_REMOTE} > 64 ))' '(( ${#PGR_REMOTE} > 63 ))' p2-remote-64
mut m-p2-zweig-nicht-gesetzt s "$A2" 'if [ -z "${PGR_ZWEIG+x}" ]; then meld "PGR_ZWEIG ist nicht gesetzt"' 'if false; then meld "PGR_ZWEIG ist nicht gesetzt"' p2-zweig-nicht-gesetzt
mut m-p2-zweig-leer s "$A2" 'elif [ -z "$PGR_ZWEIG" ]; then' 'elif false; then' p2-zweig-leer
mut m-p2-zweig-main s "$A2" 'elif [ "$PGR_ZWEIG" = main ]; then' 'elif false; then' p2-zweig-main p2-alle-gemeldet p2-reihenfolge
mut m-p2-zweig-minus s "$A2" 'elif [[ $PGR_ZWEIG == -* ]]; then' 'elif false; then' p2-zweig-minus
mut m-p2-zweig-refs s "$A2" 'elif [[ $PGR_ZWEIG == refs/* ]]; then' 'elif false; then' p2-zweig-refs
mut m-p2-zweig-check-ref-format s "$A2" 'elif command -v git > /dev/null 2>&1 && ! git check-ref-format --branch "$PGR_ZWEIG" > /dev/null 2>&1; then' 'elif false; then' p2-zweig-punkte p2-zweig-leerzeichen
# Punkt 14, Ausnahme: ohne `git` bleibt allein die Prüfung von check-ref-format aus (Entscheidung des
# Architect, 2026-10-21); ohne den Vorbehalt meldet `eingaben` auch ohne `git` den Zweignamen.
mut m-p2-zweig-ohne-git s "$A3" 'elif command -v git > /dev/null 2>&1 && ! git check-ref-format --branch "$PGR_ZWEIG" > /dev/null 2>&1; then' 'elif ! git check-ref-format --branch "$PGR_ZWEIG" > /dev/null 2>&1; then' p3-fehlt-git p3-fehlt-beide p3-fehlt-git-zweig-ungueltig
mut m-p2-zweig-kurz s "$A2" 'elif [[ $PGR_ZWEIG == -* ]]; then' 'elif [[ $PGR_ZWEIG == -* || $PGR_ZWEIG = a ]]; then' p2-zweig-a
mut m-p2-zweig-lang s "$A2" 'elif [[ $PGR_ZWEIG == refs/* ]]; then' 'elif [[ $PGR_ZWEIG == refs/* ]] || (( ${#PGR_ZWEIG} > 20 )); then' p2-zweig-kennung
mut m-p2-opts-steuerzeichen s "$A2" "[[ \$PGR_SSH_OPTS != *[\$'\\001'-\$'\\010'\$'\\012'-\$'\\037\\177']* ]]" 'true' p2-opts-zeilenumbruch p2-alle-gemeldet p2-reihenfolge
mut m-p2-opts-tabulator s "$A2" "\$'\\012'-\$'\\037\\177']* ]]" "\$'\\011'-\$'\\037\\177']* ]]" p2-opts-tabulator
mut m-p2-opts-ein-argument s "$A2" '  sshopts=($PGR_SSH_OPTS)' '  sshopts=("$PGR_SSH_OPTS")' p2-opts-tabulator p2-opts-leer p2-opts-nicht-gesetzt
mut m-p2-opts-shell s "$A2" '  sshopts=($PGR_SSH_OPTS)' '  eval "sshopts=($PGR_SSH_OPTS)"' p2-opts-keine-shell p4-nichts-geaendert
mut m-p2-opts-glob s "$A2" '  set -f
  sshopts=($PGR_SSH_OPTS)
  set +f' '  sshopts=($PGR_SSH_OPTS)' p2-opts-muster
mut m-p2-vorgabe-port s "$A2" 'PGR_SSH_PORT=${PGR_SSH_PORT-2222}' 'PGR_SSH_PORT=${PGR_SSH_PORT-2223}' p2-vorgaben p2-opts-nicht-gesetzt p2-opts-leer p2-opts-eine p2-opts-zwei p2-opts-tabulator p2-opts-keine-shell p2-opts-muster
mut m-p2-vorgabe-remote s "$A2" 'PGR_REMOTE=${PGR_REMOTE-container}' 'PGR_REMOTE=${PGR_REMOTE-behaelter}' p2-vorgaben
mut m-p2-reihenfolge s "$A2" 'for v in PGR_DEV PGR_SSH_PORT PGR_REMOTE PGR_ZWEIG PGR_SSH_OPTS' 'for v in PGR_SSH_PORT PGR_DEV PGR_REMOTE PGR_ZWEIG PGR_SSH_OPTS' p2-reihenfolge
mut m-p2-erster-verstoss s "$A2" '  (( bad == 0 )) || exit 2' '  (( bad == 0 )) || { (( bad > 1 )) && exit 2; exit 2; }' # äquivalent (§7: Sonde ist die folgende Mutation)
# Sonde zu m-p2-erster-verstoss: dieselbe Stelle, aber ein Zweig, den ein Fall erreicht, endet mit einem
# anderen Exit-Code (`meld` setzt bad nur auf 1, ein Zweig für bad > 1 wird nie erreicht); die Fälle mit
# PGR_ZWEIG=main sehen den Code.
mut m-p2-erster-verstoss-anderer-exit s "$A2" '  (( bad == 0 )) || exit 2' '  (( bad == 0 )) || { [ "${PGR_ZWEIG-}" = main ] && exit 9; exit 2; }' p2-zweig-main p2-alle-gemeldet p2-reihenfolge
mut m-p2-alle-gemeldet s "$A2" '  for v in PGR_DEV PGR_SSH_PORT PGR_REMOTE PGR_ZWEIG PGR_SSH_OPTS; do "pruefe_$v"; done' '  for v in PGR_DEV PGR_SSH_PORT PGR_REMOTE PGR_ZWEIG PGR_SSH_OPTS; do "pruefe_$v"; (( bad == 0 )) || break; done' p2-alle-gemeldet p2-reihenfolge
# Punkt 3: Werkzeuge und Projektname
mut m-p3-ssh s "$A3" '  for t in ssh git; do' '  for t in git; do' p3-fehlt-ssh p3-fehlt-beide p14-werkzeug-vor-geruest
mut m-p3-git s "$A3" '  for t in ssh git; do' '  for t in ssh; do' p3-fehlt-git p3-fehlt-beide p3-fehlt-git-zweig-ungueltig
mut m-p3-reihenfolge s "$A3" '  for t in ssh git; do' '  for t in git ssh; do' p3-fehlt-beide
mut m-p3-name-aus-basisnamen s "$A3" '  while IFS= read -r z; do [[ $z == '"'"'    name: '"'"'* ]] && { projekt=${z#'"'"'    name: '"'"'}; break; }; done < "$y"' '  projekt=$b' p3-name-zeile-fehlt p3-name-gross p3-name-unterstrich p3-name-ungleich-basis p3-name-anderer-einzug
mut m-p3-name-einzug s "$A3" "[[ \$z == '    name: '* ]]" "[[ \$z == *'name: '* ]]" p3-name-anderer-einzug
mut m-p3-name-zeichen s "$A3" '[[ $projekt =~ ^[a-z0-9][a-z0-9-]*$ ]]' '[[ $projekt =~ ^[A-Za-z0-9][A-Za-z0-9_-]*$ ]]' p3-name-gross p3-name-unterstrich
mut m-p3-name-vergleich s "$A3" '  [ "$projekt" = "$b" ] ||' '  [ -n "$projekt" ] ||' p3-name-ungleich-basis
mut m-p3-name-zweite-zeile s "$A3" '; break; }; done < "$y"' '; }; done < "$y"' p3-name-zwei-zeilen
mut m-p3-datei-fehlt s "$A3" '  [ -f "$y" ] || ende 4 "$y fehlt — kein Gerüst; make container-start legt es an"' '  :' p3-name-datei-fehlt p4-meldung-form
# Punkt 4: Ströme und Reihenfolge der Schritte
mut m-p4-stdout s "$A1" 's_hooks() {' 's_hooks() {
  echo zwischenstand' p4-stdout-nur-erfolg
mut m-p4-werkzeug-stdout s 'p4-.*' '  fern "cd $klonpfad && make hooks-install"' '  fern "cd $klonpfad && make hooks-install" > /dev/null' p4-werkzeug-stdout
mut m-p4-werkzeug-stderr s 'p4-.*' '  fern "cd $klonpfad && make hooks-install"' '  fern "cd $klonpfad && make hooks-install" 2> /dev/null' p4-werkzeug-stderr
mut m-p4-lesezug-stderr s 'p4-.*' '[ -z "$z" ] || printf' '[ -n "$z" ] || printf' p4-werkzeug-stderr-ssh
mut m-p4-meldung-ohne-schritt s "$A3" 'echo "container-einrichten: $schritt: $*" >&2' 'echo "container-einrichten: $*" >&2' p4-meldung-form p3-fehlt-ssh p3-fehlt-git p3-fehlt-beide
mut m-p4-schritte s 'p4-.*|p12-.*' 'for schritt in eingaben werkzeuge zustand bare einliefern klon identitaet hooks zweig traeger' 'for schritt in eingaben werkzeuge zustand bare einliefern klon identitaet hooks traeger zweig' p4-reihenfolge-der-schritte p12-stelle p12-fehler
mut m-p14-eingaben-nach-werkzeuge s "$A3" 'for schritt in eingaben werkzeuge zustand' 'for schritt in werkzeuge eingaben zustand' p14-eingabe-vor-werkzeug
# Punkt 5: Zustand
mut m-p5-schluessel-vorhanden s "$A5" '  [ -e "$PGR_DEV/ssh/id_ed25519" ] || ende 4 "$PGR_DEV/ssh/id_ed25519 fehlt"' '  :' p5-schluessel-fehlt
mut m-p5-schluessel-lesbar s "$A5" '  [ -r "$PGR_DEV/ssh/id_ed25519" ] ||' '  [ -e "$PGR_DEV/ssh/id_ed25519" ] ||' p5-schluessel-unlesbar
mut m-p5-repo s "$A5" '  git rev-parse --show-toplevel > /dev/null 2>&1 ||' '  true ||' p5-kein-repo
mut m-p5-main s "$A5" '  git rev-parse -q --verify refs/heads/main > /dev/null 2>&1 ||' '  true ||' p5-kein-main p14-repo-vor-identitaet
mut m-p5-main-als-head s "$A5" 'git rev-parse -q --verify refs/heads/main' 'git rev-parse -q --verify HEAD' p5-kein-main p14-repo-vor-identitaet
# Randform 12: das Ziel hängt nicht am ausgecheckten Zweig
mut m-p5-ausgecheckter-zweig s "$A5" '  git rev-parse -q --verify refs/heads/main > /dev/null 2>&1 ||' '  { [ "$(git rev-parse --abbrev-ref HEAD)" = main ] && git rev-parse -q --verify refs/heads/main > /dev/null 2>&1; } ||' p5-detached p5-anderer-zweig
mut m-p5-ident-name s "$A5" '  [ -n "${iname//[ $'"'"'\t'"'"']/}" ] || f+=(user.name)' '  :' p5-ident-name-fehlt p5-ident-beide-fehlen p5-ident-leer p5-ident-leerraum p14-identitaet-vor-remote
mut m-p5-ident-mail s "$A5" '  [ -n "${imail//[ $'"'"'\t'"'"']/}" ] || f+=(user.email)' '  :' p5-ident-mail-fehlt p5-ident-beide-fehlen
mut m-p5-ident-leerraum s "$A5" '${iname//[ $'"'"'\t'"'"']/}' '$iname' p5-ident-leerraum
mut m-p5-ident-beide s "$A5" '  (( ${#f[@]} == 0 )) || ende 4 "die Identität des Host-Repos fehlt: ${f[*]}"' '  (( ${#f[@]} == 0 )) || ende 4 "die Identität des Host-Repos fehlt: ${f[0]}"' p5-ident-beide-fehlen
mut m-p5-ident-steuerzeichen s "$A5" '[\x20-\x7e]' '[\x09-\x0d\x20-\x7e]' p5-ident-tabulator p5-ident-zeilenumbruch p5-ident-cr
mut m-p5-ident-tabulator s "$A5" '[\x20-\x7e]' '[\x09\x20-\x7e]' p5-ident-tabulator
mut m-p5-ident-utf8 s "$A5" '|\xf4[\x80-\x8f][\x80-\xbf][\x80-\xbf])*$' '|[\x80-\xff])*$' p5-ident-nicht-utf8 p5-ident-mail-nicht-utf8
mut m-p5-ident-nur-ascii s "$A5|p9-zeichen-(umlaut|cjk|emoji)" '[\xc2-\xdf][\x80-\xbf]|' '' p5-ident-utf8-gut p9-zeichen-umlaut
mut m-p5-ident-nur-lokal s "$A5" 'iname=$(git config --get user.name 2> /dev/null)' 'iname=$(git config --local --get user.name 2> /dev/null)' p5-ident-global p4-nichts-geaendert
mut m-p5-ident-mail-geprueft s "$A5" '  [[ $imail =~ $re ]] ||' '  true ||' p5-ident-mail-nicht-utf8
mut m-p5-ident-wert-in-meldung s "$A5" 'ende 4 "user.name des Host-Repos trägt ein Steuerzeichen oder kein gültiges UTF-8 (der Wert steht nicht in dieser Meldung)"' 'ende 4 "user.name $iname des Host-Repos trägt ein Steuerzeichen oder kein gültiges UTF-8"' p5-ident-tabulator p5-ident-zeilenumbruch p5-ident-cr p5-ident-nicht-utf8 p4-meldung-form
mut m-p5-ident-wert-in-meldung-mail s "$A5" 'ende 4 "user.email des Host-Repos trägt ein Steuerzeichen oder kein gültiges UTF-8 (der Wert steht nicht in dieser Meldung)"' 'ende 4 "user.email des Host-Repos trägt ein Steuerzeichen oder kein gültiges UTF-8 (der Wert steht nicht in dieser Meldung) $imail"' p5-ident-mail-nicht-utf8
mut m-p5-remote-url s "$A5" '    [ "$ist" = "$url" ] ||' '    [ -n "$ist" ] ||' p5-remote-andere-url p14-remote-vor-lesezug p5-remote-nichts-gesetzt
mut m-p5-remote-name-fest s "$A5|p7-.*" 'git remote add "$PGR_REMOTE" "$url"' 'git remote add container "$url"' p5-remote-anderer-name
mut m-p7-url-ohne-port s "$A5|p7-.*|p2-vorgaben" 'url=ssh://vscode@127.0.0.1:$PGR_SSH_PORT/workspaces/$projekt/austausch.git' 'url=ssh://vscode@127.0.0.1/workspaces/$projekt/austausch.git' p2-vorgaben p5-remote-andere-url p5-remote-anderer-name p5-remote-gleich p7-remote-fehlt p7-url-port
# Der Port der URL stammt aus der Eingabe: ein fest eingesetzter 2222 färbt allein den Fall mit einem anderen Port.
mut m-p7-url-port-fest s "$A5|p7-.*|p2-vorgaben|p2-port-.*" 'url=ssh://vscode@127.0.0.1:$PGR_SSH_PORT/workspaces/$projekt/austausch.git' 'url=ssh://vscode@127.0.0.1:2222/workspaces/$projekt/austausch.git' p7-url-port
mut m-p5-remote-gesetzt s "$A5|p1-nie-remote-set-url" '  if ist=$(git remote get-url "$PGR_REMOTE" 2> /dev/null); then' '  if ist=$(git remote set-url "$PGR_REMOTE" "$url" 2> /dev/null); then' p5-remote-andere-url p5-remote-nichts-gesetzt p5-remote-gleich p14-remote-vor-lesezug p1-nie-remote-set-url p4-nichts-geaendert
mut m-p5-lesezug-batchmode s "$A5" ' -o BatchMode=yes -o ConnectTimeout=5' ' -o ConnectTimeout=5' p5-lesezug-aufruf
mut m-p5-lesezug-timeout s "$A5" ' -o ConnectTimeout=5' '' p5-lesezug-aufruf
mut m-p5-lesezug-hostschluessel s "$A5" ' -o BatchMode=yes -o ConnectTimeout=5' ' -o BatchMode=yes -o StrictHostKeyChecking=no -o ConnectTimeout=5' p5-lesezug-kein-hostschluessel p5-lesezug-aufruf
mut m-p5-lesezug-zwei-aufrufe s "$A5" '  fern -a "$(lesebefehl)"' '  fern -a "$(lesebefehl)"; fern -a "$(lesebefehl)"' p5-lesezug-ein-aufruf p4-nichts-geaendert p4-reihenfolge-der-schritte
mut m-p5-lesezug-schreibt s "$A5" "'if [ ! -e \"\$b\" ] && [ ! -L \"\$b\" ]; then echo bare=fehlt'" "': > \"\$b.spur\"; if [ ! -e \"\$b\" ] && [ ! -L \"\$b\" ]; then echo bare=fehlt'" p5-lesezug-schreibt-nichts p4-nichts-geaendert
mut m-p5-nicht-erreichbar s "$A5" '  if [ "$kcode" != 0 ]; then' '  if false; then' p5-nicht-erreichbar p14-lesezug-vor-werkzeugen p4-meldung-form
mut m-p5-meldung-ohne-ausgabe s "$A5" '(ssh: $s) — make container-start startet ihn' '— make container-start startet ihn' p5-nicht-erreichbar
mut m-p5-fern-git s "$A5" '  for t in git make; do
    [ "${Z[wz-$t]-}" = ja ]' '  for t in make; do
    [ "${Z[wz-$t]-}" = ja ]' p5-fern-git-fehlt
mut m-p5-fern-make s "$A5" '  for t in git make; do
    [ "${Z[wz-$t]-}" = ja ]' '  for t in git; do
    [ "${Z[wz-$t]-}" = ja ]' p5-fern-make-fehlt p14-werkzeuge-vor-bare
mut m-p5-bare-head s "$A5" '    head:*) ende 4' '    head:*) : || ende 4' p5-bare-master p5-lesezug-schreibt-nichts p14-bare-vor-klon
mut m-p5-bare-fremd s "$A5" '    *) ende 4 "$barepfad besteht und ist kein Bare-Repo" ;;' '    *) ;;' p5-bare-fremd
mut m-p5-klon-schmutzig s "$A5" '      [ "${Z[status]-}" = sauber ] ||' '      true ||' p5-klon-schmutzig p5-klon-unversioniert p14-klon-vor-schritt
mut m-p5-klon-origin s "$A5" '      [ "${Z[origin]-}" = "$barepfad" ] ||' '      true ||' p5-klon-anderer-origin
mut m-p5-klon-fremd s "$A5" '    *) ende 4 "$klonpfad besteht und ist kein Git-Repo" ;;' '    *) ;;' p5-klon-fremd
mut m-p5-klon-vor-bare s "$A5" '  zust_werkzeuge_fern; zust_bare; zust_klon' '  zust_werkzeuge_fern; zust_klon; zust_bare' p14-bare-vor-klon
mut m-p5-identitaet-nach-remote s "$A5" '  zust_geruest; zust_repo; zust_identitaet; zust_remote; zust_lesezug' '  zust_geruest; zust_repo; zust_remote; zust_identitaet; zust_lesezug' p14-identitaet-vor-remote
mut m-p5-geruest-nach-repo s "$A5" '  zust_geruest; zust_repo;' '  zust_repo; zust_geruest;' p14-geruest-vor-repo
mut m-p5-repo-nach-identitaet s "$A5" '  zust_geruest; zust_repo; zust_identitaet;' '  zust_geruest; zust_identitaet; zust_repo;' p14-repo-vor-identitaet p5-kein-repo
mut m-p5-remote-nach-lesezug s "$A5" ' zust_remote; zust_lesezug' ' zust_lesezug; zust_remote' p14-remote-vor-lesezug
mut m-p5-werkzeuge-vor-lesezug s "$A5" '  zust_werkzeuge_fern; zust_bare' '  zust_bare' p5-fern-git-fehlt p5-fern-make-fehlt p14-werkzeuge-vor-bare
# Punkt 6 bis 12: die handelnden Schritte
mut m-p6-initial-branch s 'p6-.*|p8-.*' ' --initial-branch=main' '' p6-fehlt
mut m-p6-immer s 'p6-.*|p13-.*' '  [ "${Z[bare]}" = fehlt ] || return 0' '  :' p6-vorhanden p13-zweiter-aufruf
mut m-p6-fehler s 'p6-.*' '  [ "$kcode" = 0 ] || ende 5 "git init --bare $barepfad endete mit Exit $kcode"' '  :' p6-fehler
mut m-p7-remote-add s 'p7-.*|p13-.*' '  (( rneu == 0 )) || git remote add' '  (( rneu == 1 )) || git remote add' p7-remote-fehlt p7-push-aufruf p7-push-angekommen p7-git-ssh-command p7-git-ssh-command-opts p7-abgewiesen p7-url-port p13-erfolg p13-zweiter-aufruf p13-umgebung
mut m-p7-force s 'p7-.*|p1-nie-push-force' '  kind -e env "GIT_SSH_COMMAND=$gsc" git push "$PGR_REMOTE"' '  kind -e env "GIT_SSH_COMMAND=$gsc" git push --force "$PGR_REMOTE"' p7-push-aufruf p7-abgewiesen p1-nie-push-force
mut m-p7-refspec s 'p7-.*' ' refs/heads/main:refs/heads/main' ' main' p7-push-aufruf
mut m-p7-sshcommand-konfiguration s 'p7-.*|p1-.*' '  kind -e env "GIT_SSH_COMMAND=$gsc" git push' '  git config --local core.sshCommand "$gsc"; kind -e git push' p7-git-ssh-command p7-git-ssh-command-opts p1-sshcommand p1-andere-einstellung
mut m-p7-abgewiesen-nicht-erkannt s 'p7-.*' '  [[ $kausg != *'"'"'[rejected]'"'"'* ]] ||' '  true ||' p7-abgewiesen
mut m-p7-abgewiesen-immer s 'p7-.*' '  [[ $kausg != *'"'"'[rejected]'"'"'* ]] ||' '  false ||' p7-fehler
mut m-p7-zweiter-versuch s 'p7-.*' '  [ "$kcode" = 0 ] && return 0' '  [ "$kcode" = 0 ] || kind -e env "GIT_SSH_COMMAND=$gsc" git push "$PGR_REMOTE" refs/heads/main:refs/heads/main
  [ "$kcode" = 0 ] && return 0' p7-abgewiesen
mut m-p7-fehler s 'p7-.*' '  ende 5 "git push $PGR_REMOTE endete mit Exit $kcode"' '  :' p7-fehler
mut m-p8-clone-immer s 'p8-.*|p13-.*' '  if [ "${Z[klon]}" = fehlt ]; then' '  if true; then' p8-besteht p8-fehler p13-zweiter-aufruf
mut m-p8-clone-ziel s 'p8-.*|p13-.*' 'git clone $barepfad $klonpfad' 'git clone $barepfad $klonpfad-zwei' p8-fehlt p13-erfolg p13-umgebung p13-zweiter-aufruf
mut m-p8-prune s 'p8-.*' ' fetch --prune origin' ' fetch origin' p8-besteht p8-fehler
mut m-p8-pull s 'p8-.*' 'git -C $klonpfad fetch --prune origin' 'git -C $klonpfad pull' p8-besteht p8-fehler
mut m-p8-fehler s 'p8-.*' '  [ "$kcode" = 0 ] || ende 5 "der Klon $klonpfad endete mit Exit $kcode"' '  :' p8-fehler
mut m-p9-global s "$A9" 'config --local user.name "$n"' 'config --global user.name "$n"' p9-lokal p4-reihenfolge-der-schritte "${ZALLE[@]}"
mut m-p9-system s "$A9" 'config --local user.name "$n"' 'config --system user.name "$n"' p9-lokal p4-reihenfolge-der-schritte "${ZALLE[@]}"
mut m-p9-wert-aus-skript s "$A9" 'IFS= read -r n; IFS= read -r e;' 'n=Fester\ Name; IFS= read -r x; IFS= read -r e;' p9-lokal "${ZALLE[@]}"
mut m-p9-als-argument s "$A9" 'IFS= read -r n; IFS= read -r e; git -C '"'"'"$klonpfad"'"'"' config --local user.name "$n"' 'git -C '"'"'"$klonpfad"'"'"' config --local user.name '"'"'"$iname"'"'"'' p9-kein-argument p9-lokal p9-zeichen-leerzeichen "${ZSHELL[@]}"
mut m-p9-read-ohne-r s "$A9" 'IFS= read -r n; IFS= read -r e;' 'IFS= read n; IFS= read -r e;' p9-zeichen-backslash
mut m-p9-read-ohne-ifs s "$A9" 'IFS= read -r n; IFS= read -r e;' 'read -r n; IFS= read -r e;' p9-zeichen-leerraum-am-rand
mut m-p9-vertauscht s "$A9" '  kein=$iname$'"'"'\n'"'"'$imail$'"'"'\n'"'"'' '  kein=$imail$'"'"'\n'"'"'$iname$'"'"'\n'"'"'' p9-lokal p9-reihenfolge "${ZALLE[@]}"
mut m-p9-in-shell-zeichenkette s "$A9" 'config --local user.name "$n"' 'config --local user.name "$(eval echo $n)"' "${ZSHELL[@]}"
mut m-p9-fehler s "$A9" '  [ "$kcode" = 0 ] || ende 5 "die Identität des Klons $klonpfad endete mit Exit $kcode"' '  :' p9-fehler
mut m-p10-ziel s 'p10-.*|p13-.*' 'fern "cd $klonpfad && make hooks-install"' 'fern "cd $klonpfad && make hooks"' p10-aufruf p10-zweiter-aufruf p13-zweiter-aufruf p10-fehler
mut m-p10-verzeichnis s 'p10-.*' '  fern "cd $klonpfad && make hooks-install"' '  fern "cd / && make hooks-install"' p10-aufruf p10-zweiter-aufruf
mut m-p10-fehler s 'p10-.*' '  [ "$kcode" = 0 ] || ende 5 "make hooks-install im Klon endete mit Exit $kcode"' '  :' p10-fehler
mut m-p11-quelle s 'p11-.*' 'switch -c "$z" origin/main' 'switch -c "$z" main' p11-aus-origin-main
mut m-p11-immer-c s 'p11-.*' '  if (( zweigda )); then' '  if false; then' p11-besteht
mut m-p11-nie-c s 'p11-.*' '  if (( zweigda )); then' '  if true; then' p11-fehlt p11-aus-origin-main p11-main-unbewegt p11-sonderzeichen
mut m-p11-main-bewegt s 'p11-.*' '  [ "$kcode" = 0 ] || ende 5 "der Zweig $PGR_ZWEIG im Klon endete mit Exit $kcode"' '  [ "$kcode" = 0 ] || ende 5 "der Zweig $PGR_ZWEIG im Klon endete mit Exit $kcode"
  fern "git -C $klonpfad branch -f main origin/main"' p11-main-unbewegt
mut m-p11-force s 'p11-.*' 'git -C '"'"'"$klonpfad"'"'"' switch "$z"' 'git -C '"'"'"$klonpfad"'"'"' switch --force "$z"' p11-besteht
mut m-p11-kein-switch s 'p11-.*|p12-.*' '  if (( zweigda )); then
    fern '"'"'IFS= read -r z; git -C '"'"'"$klonpfad"'"'"' switch "$z"'"'"'' '  if (( zweigda )); then
    fern '"'"'IFS= read -r z; git -C '"'"'"$klonpfad"'"'"' checkout -f "$z"'"'"'' p11-besteht
mut m-p11-zweig-als-argument s 'p11-.*' 'IFS= read -r z; git -C '"'"'"$klonpfad"'"'"' switch -c "$z" origin/main' 'git -C '"'"'"$klonpfad"'"'"' switch -c '"'"'"$PGR_ZWEIG"'"'"' origin/main' p11-sonderzeichen
mut m-p11-fehler s 'p11-.*' '  [ "$kcode" = 0 ] || ende 5 "der Zweig $PGR_ZWEIG im Klon endete mit Exit $kcode"' '  :' p11-fehler
mut m-p12-ziel s 'p12-.*' 'fern "cd $klonpfad && make traeger-fetch"' 'fern "cd $klonpfad && make traeger"' p12-aufruf p12-stelle p12-fehler
mut m-p12-fehler s 'p12-.*' '  [ "$kcode" = 0 ] || ende 5 "make traeger-fetch im Klon endete mit Exit $kcode"' '  :' p12-fehler
mut m-p13-erfolgszeile s 'p13-.*' 'echo "container-einrichten: eingerichtet $klonpfad $PGR_ZWEIG"' 'echo "container-einrichten: eingerichtet"' p13-erfolg p13-zweiter-aufruf p13-umgebung
for v in CLAUDE_CODE_OAUTH_TOKEN ANTHROPIC_API_KEY GIT_TOKEN; do
  mut "m-p13-umgebung-$v" s 'p13-.*' "unset $v" ": $v" p13-umgebung
done
# Punkt 15: Abbruch
mut m-p15-int s "$A15" "trap 'abbruch 130' INT" ':' p15-sigint p15-sigint-zweig p15-kein-rueckbau p15-wiederholung
mut m-p15-term s "$A15" "trap 'abbruch 143' TERM" ':' p15-sigterm p15-sigterm-traeger
mut m-p15-exit-int s "$A15" "trap 'abbruch 130' INT" "trap 'abbruch 1' INT" p15-sigint p15-sigint-zweig p15-kein-rueckbau p15-wiederholung
mut m-p15-exit-term s "$A15" "trap 'abbruch 143' TERM" "trap 'abbruch 1' TERM" p15-sigterm p15-sigterm-traeger
mut m-p15-kind-lebt s "$A15" 'abbruch() { [ -z "$kpid" ] || kill "$kpid" 2> /dev/null;' 'abbruch() {' p15-sigint p15-sigterm p15-sigterm-traeger p15-sigint-zweig p15-kein-rueckbau p15-wiederholung
mut m-p15-ohne-meldung s "$A15" 'ende "$1" abgebrochen' 'exit "$1"' p15-sigint p15-sigterm p15-sigterm-traeger p15-sigint-zweig p15-kein-rueckbau p15-wiederholung
mut m-p15-rueckbau s "$A15" 'abbruch() { [ -z "$kpid" ] || kill "$kpid" 2> /dev/null;' 'abbruch() { [ -z "$kpid" ] || kill "$kpid" 2> /dev/null; rm -rf "$HOME/.ssh";' p15-kein-rueckbau
# Punkt 16: Aufrufarten des Treibers (Fragment)
mut m-p16-rezept-make m "$A16" '	@bash tools/harness/container-einrichten.sh' '	@$(MAKE) -f /dev/null -q; bash tools/harness/container-einrichten.sh' p16-n p16-q p16-t p16-makeflags-n p16-makeflags-q p16-makeflags-t p16-n-und-i
mut m-p16-nicht-phony m "$A16|p17-phony" '.PHONY: container-einrichten' '.PHONY: container-einrichten-nicht' p16-phony p17-phony
mut m-p16-voraussetzung m "$A16" 'container-einrichten: ##' 'container-einrichten: Makefile ##' p16-keine-voraussetzung
mut m-p16-voraussetzung-fehlt m "$A16" 'container-einrichten: ##' 'container-einrichten: gibtsnicht ##' p16-phony p16-k p16-j p16-i p16-makeflags-i p16-makeflags-j p16-makeflags-k p16-keine-voraussetzung
# Punkt 17: Listen und Zeitbudget
mut m-p17-liste-fehlt g "$A17" 'GEGENPROBEN_CONTAINER += container-einrichten-gegenprobe' '# keine Liste' p17-liste-container
mut m-p17-liste-host g "$A17" 'GEGENPROBEN_CONTAINER += container-einrichten-gegenprobe' 'GEGENPROBEN_HOST += container-einrichten-gegenprobe' p17-liste-container p17-nicht-host
mut m-p17-gate g "$A17" 'GEGENPROBEN_CONTAINER += container-einrichten-gegenprobe' 'GEGENPROBEN_CONTAINER += container-einrichten-gegenprobe
GATE_CHECKS += container-einrichten-gegenprobe' p17-kein-gate
mut m-p17-gate-ziel m "$A17" '.PHONY: container-einrichten' 'GATE_CHECKS += container-einrichten
.PHONY: container-einrichten' p17-kein-gate
mut m-p17-ausnahme k "$A17" 'GATE_SCHNELL_AUSNAHME := test-container test-integration test-postgres-versionen' 'GATE_SCHNELL_AUSNAHME := test-container test-integration test-postgres-versionen container-einrichten-gegenprobe' p17-keine-ausnahme
mut m-p17-budget-pruefung p "$A17B" '  if [ -n "$budget" ] && [ "$budget" != 0 ] && (( T > budget )); then' '  if false; then' p17-budget-ueber p17-budget-gruende
mut m-p17-budget-null p "$A17B" '[ "$budget" != 0 ] && ' '' p17-budget-null
mut m-p17-budget-zeile p "$A17B" '    echo "$GPN: Zeit $T s von höchstens $budget s, Last $last" >&2' '    :' p17-budget-unter p17-budget-null p17-budget-last-fehlt p17-budget-gruende p17-budget-zulaessig-10 p17-budget-zulaessig-999999
mut m-p17-budget-ohne-mut p "$A17B" 'elif [ -z "${GP_MUT-}" ]; then
  budget=1200' 'elif true; then
  budget=1200' p17-budget-vorgabe-nur-ohne-mut
mut m-p17-budget-last p "$A17B" 'read -r last _ < "$LOADAVG"' 'last=0' p17-budget-unter p17-budget-null p17-budget-vorgabe-nur-ohne-mut p17-budget-zulaessig-10 p17-budget-zulaessig-999999
mut m-p17-budget-last-fehlt p "$A17B" '[ ! -r "$LOADAVG" ] ||' '[ -r "$LOADAVG" ] ||' p17-budget-unter p17-budget-null p17-budget-vorgabe-nur-ohne-mut p17-budget-zulaessig-10 p17-budget-zulaessig-999999
mut m-p17-budget-ungueltig p "$A17B" '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] ||' '  [[ $GP_BUDGET =~ ^[0-9.-]*$ ]] ||' p17-budget-ungueltig-leer p17-budget-ungueltig-minus p17-budget-ungueltig-bruch p17-budget-ungueltig-08 p17-budget-ungueltig-09 p17-budget-ungueltig-00 p17-budget-ungueltig-007 p17-budget-ungueltig-sieben-stellen p17-budget-ungueltig-zwanzig-stellen
mut m-p17-budget-fuehrende-null-erlaubt p "$A17B" '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] ||' '  [[ $GP_BUDGET =~ ^[0-9]{1,6}$ ]] ||' p17-budget-ungueltig-08 p17-budget-ungueltig-09 p17-budget-ungueltig-00 p17-budget-ungueltig-007
mut m-p17-budget-stellen-unbegrenzt p "$A17B" '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] ||' '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]*)$ ]] ||' p17-budget-ungueltig-sieben-stellen p17-budget-ungueltig-zwanzig-stellen
mut m-p17-budget-null-abgelehnt p "$A17B" '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] ||' '  [[ $GP_BUDGET =~ ^[1-9][0-9]{0,5}$ ]] ||' p17-budget-null
mut m-p17-budget-sechs-stellen-abgelehnt p "$A17B" '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] ||' '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,4})$ ]] ||' p17-budget-zulaessig-999999
mut m-p17-budget-zehn-abgelehnt p "$A17B" '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] ||' '  [[ $GP_BUDGET =~ ^(0|[1-9]|[1-9][0-9]{2,5})$ ]] ||' p17-budget-zulaessig-10
mut m-p17-budget-vorzeichen-plus-erlaubt p "$A17B" '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] ||' '  [[ $GP_BUDGET =~ ^[+]?(0|[1-9][0-9]{0,5})$ ]] ||' p17-budget-ungueltig-plus
mut m-p17-budget-leerraum-vorn-erlaubt p "$A17B" '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] ||' '  [[ $GP_BUDGET =~ ^\ ?(0|[1-9][0-9]{0,5})$ ]] ||' p17-budget-ungueltig-leerraum-vorn
mut m-p17-budget-leerraum-hinten-erlaubt p "$A17B" '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] ||' '  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})\ ?$ ]] ||' p17-budget-ungueltig-leerraum-hinten
mut m-p17-budget-vor-faellen p "$A17B" 'if [ -n "${GP_BUDGET+x}" ]; then
  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] || { echo "$GPN: GP_BUDGET: keine ganze Zahl" >&2; exit 2; }' 'if [ -n "${GP_BUDGET+x}" ]; then
  [[ $GP_BUDGET =~ ^(0|[1-9][0-9]{0,5})$ ]] || { nachher=1; }' p17-budget-ungueltig-leer p17-budget-ungueltig-minus p17-budget-ungueltig-bruch p17-budget-ungueltig-abc p17-budget-ungueltig-08 p17-budget-ungueltig-09 p17-budget-ungueltig-00 p17-budget-ungueltig-007 p17-budget-ungueltig-sieben-stellen p17-budget-ungueltig-zwanzig-stellen p17-budget-ungueltig-plus p17-budget-ungueltig-leerraum-vorn p17-budget-ungueltig-leerraum-hinten
mut m-p17-signal-int p "$A17O" "trap 'exit 130' INT" ':' p17-budget-sigint
mut m-p17-signal-term p "$A17O" "trap 'exit 143' TERM" ':'
mut m-p17-ordner-ende p "$A17O" 'raeumen() { [ ! -d "$arbeit" ] || rm -rf -- "$arbeit"; return 0; }' 'raeumen() { return 0; }' p17-budget-sigint p17-budget-sigterm
mut m-p17-ordner-mutation p "$A17O" 'erwartet [${soll:-nichts}]"
    rm -rf -- "$basis"' 'erwartet [${soll:-nichts}]"' p17-ordner-mutation
mut m-p17-ordner-fall p "$A17O" '    eval "${FC[i]}" > /dev/null
    rm -rf -- "$w"' '    eval "${FC[i]}" > /dev/null' p17-ordner-fall

wait
fehler=0
if [ -s "$arbeit/fehler" ]; then cat "$arbeit/fehler" >&2; fehler=1; fi
abschluss
[ "$fehler" = 0 ] || exit 1
echo "$GPN: gruen — ${#FN[@]} Fälle zu SPEC-069 Punkt 1 bis 17, darunter die zu SPEC-060 Punkt 16 (9); je Mutation genau die benannten roten Fälle, ausgenommen die ohne roten Fall (äquivalent, §7 des Plans): ${GRUEN[*]-keine}$([ "$root" = 0 ] || echo "; als root entfällt der Fall p5-schluessel-unlesbar und mit ihm die Mutation ${ROOTWEG[*]-keine}")"
