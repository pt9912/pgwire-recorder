#!/usr/bin/env bash
# container-start-gegenprobe — prüft tools/harness/container-start.sh und harness/mk/container-start.mk
# gegen SPEC-058 (harness/vertraege/SPEC-058-container-start.md: Punkte 1 bis 16). Ziel für den Host, kein
# Gate (ADR-0047). Gruppe Attrappen: Temp-Verzeichnisse mit Attrappen
# von docker, u-boot, devcontainer und ssh in einem vorangestellten Pfad, echtes node und ssh-keygen,
# ohne Docker und ohne Netz; dazu der Fall p7-echt-u-boot gegen das echte u-boot. Jede Attrappe
# protokolliert Aufruf und Argumente. Gruppe Docker (Fälle p9-docker-*, p10-docker-*, p11-docker-*,
# p15-docker-*): echtes Docker, echtes devcontainer, ein Container des Präfixes pgr-cs- zugleich.
# Der Name jedes Falls trägt die Nummer seines Punktes; die
# unveränderte Kopie besteht jeden Fall, und bei jeder Mutation (genau eine Textersetzung in der Kopie
# einer Datei; trifft sie nichts oder mehr als eine Stelle, bricht die Gegenprobe ab) werden von den
# Fällen ihrer Gruppe (SEL) genau die benannten rot.
#
# Voraussetzungen (SPEC-058 Punkt 15): make, GNU-Werkzeuge, node, ssh-keygen, u-boot, docker, devcontainer und ssh im Pfad;
# die Gegenprobe prüft sie zuerst und endet bei einem fehlenden mit einer Meldung (Fälle p15-voraussetzung-*).
# Namen tragen das Präfix pgr-cs-: Verzeichnisse (mktemp -d), Container, Volumes, Images (vsc-pgr-cs-). Entfernt
# wird nur per Name, bei Ende, SIGINT und SIGTERM (Fälle p15-*); ein Name aus GP_NAME ohne das Präfix beendet die
# Gegenprobe vor dem ersten Aufruf von docker. GRENZE: Die Mutanten p15-ende und p15-exit-trap hinterlassen leere Verzeichnisse der Sonde;
# die Gegenprobe entfernt am Ende nur diese, nur wenn leer, nur per rmdir, und prüft, dass keines bleibt.
# GRENZE: p15-aufraeumen-int und -term setzen voraus, dass die Gegenprobe SIGINT nicht ignoriert (kein Start
# als Hintergrundjob einer Shell ohne Jobsteuerung). Die Fälle der Gruppe Docker laufen nicht neben dem Lauf
# eines anderen Slice auf demselben Daemon und nie gegen den Container, in dem gearbeitet wird (Regel an den Aufrufer).
# Zum Arbeiten: GP_MUT=<Muster> fährt nur Mutationen, deren Name passt; GP_ENTDECKE=1 nennt je Mutation die
# roten Fälle.
#
# Ausgang: 0, wenn jeder Fall und jede Mutation wie erwartet endet, sonst 1.
set -uo pipefail
export LC_ALL=C
shopt -u patsub_replacement 2> /dev/null || true

GP_SELF=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/$(basename "${BASH_SOURCE[0]}")
quelle=${GP_QUELLE:-$(cd "$(dirname "$GP_SELF")/../.." && pwd -P)}
arbeit=$(mktemp -d -t pgr-cs-XXXXXX)
bg=
PSDS=() DKC=() DKV=()
rest=0
nm=${GP_NAME:-pgr-cs-$$}
node=$(command -v node) sshkeygen=$(command -v ssh-keygen) realub=$(command -v u-boot)
sonden_raeumen() { local d; for d in "${PSDS[@]}"; do [[ ${d##*/} == pgr-cs-* ]] && [ -d "$d" ] && rmdir -- "$d" 2> /dev/null; done; return 0; }
# Punkt 15 (d): Container und Images nur per Name; die benannten Volumes überleben docker rm -fv und gehen per Name, ein
# anonymes Volume, das danach bleibt, wird gemeldet und entfernt.
dk_raeumen() {
  local n v i k
  [ -n "${dkb-}" ] || return 0
  for n in "${DKC[@]}"; do
    [[ $n == pgr-cs-* ]] || continue
    DKV+=($(docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}' "$n" 2> /dev/null))
    docker rm -fv "$n" > /dev/null 2>&1
  done
  for v in "${DKV[@]}"; do
    if [[ $v == pgr-cs-* ]]; then docker volume rm "$v" > /dev/null 2>&1
    elif docker volume inspect "$v" > /dev/null 2>&1; then echo "container-start-gegenprobe: ROT — Volume $v bleibt nach docker rm -fv" >&2; rest=1; docker volume rm "$v" > /dev/null 2>&1; fi
  done
  [ "${1-}" = bilder ] || return 0
  for k in 1 2; do for i in $(docker images --format '{{.Repository}}:{{.Tag}}' 2> /dev/null | grep "^vsc-$nm-"); do docker rmi "$i" > /dev/null 2>&1; done; done
}
aufraeumen() { [ -z "$bg" ] || kill "$bg" 2> /dev/null; [ -z "${lk-}" ] || kill "$lk" 2> /dev/null; dk_raeumen bilder; sonden_raeumen; [[ ${arbeit##*/} == pgr-cs-* ]] && [ -d "$arbeit" ] && rm -rf -- "$arbeit"; (( rest == 0 )) || exit 1; }
trap aufraeumen EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
[[ $nm == pgr-cs-* ]] || { echo "container-start-gegenprobe: ROT — der Name $nm trägt das Präfix pgr-cs- nicht, nichts wird angefasst" >&2; exit 1; }
# Punkt 15 (b): ein Port, der jetzt frei ist, nie 2222 und nie ein Port eines Containers ohne Präfix.
freier_port() {
  local p k
  for k in 1 2 3 4 5; do
    p=$("$node" -e 'const s=require("net").createServer().listen(+process.argv[1],"127.0.0.1",()=>{console.log(s.address().port);s.close()})' "${GP_PORT-0}")
    [ "$p" != 2222 ] && ! docker ps -a --format '{{.Names}} {{.Ports}}' | grep -v '^pgr-cs-' | grep -qF ":$p->" && { echo "$p"; return 0; }
  done
  return 1
}
dk_sonde() { # legt Image, Container, ein benanntes und ein anonymes Volume des Präfixes an, startet nichts und nennt sie
  local im=vsc-$nm-sonde c=$nm-sonde
  tar -cT /dev/null | docker import - "$im" > /dev/null || return 1
  DKC+=("$c") dkb=1
  docker create --name "$c" -v /anon -v "$c-vol:/benannt" "$im" /x > /dev/null || return 1
  printf '%s\n%s\n%s\n' "$c" "$c-vol" "$im"
  docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}}{{"\n"}}{{end}}{{end}}' "$c" | grep .
}
if [ -n "${GP_PROBE-}" ]; then # Sondenmodus: macht das Aufräumen von außen beobachtbar
  printf '%s\n%s\n' "$arbeit" "$$" > "$GP_PROBE_DATEI"
  [ "$GP_PROBE" = ende ] && exit 0
  [ "$GP_PROBE" != port ] || { freier_port >> "$GP_PROBE_DATEI"; exit 0; }
  [[ $GP_PROBE != docker-* ]] || dk_sonde >> "$GP_PROBE_DATEI"
  [ "$GP_PROBE" != docker-ende ] || exit 0
  sleep 60 & bg=$!
  wait "$bg"
  exit 0
fi

fehler=0
melde() { echo "container-start-gegenprobe: ROT — $*" >&2; fehler=1; }
declare -A ROT FILE=([s]=tools/harness/container-start.sh [m]=harness/mk/container-start.mk [g]=harness/mk/container-start-gegenprobe.mk [k]=harness/mk/gates-schnell.mk [p]=gp.sh)
rot() { ROT[${1:-$fn}]=1; [ -z "${GP_DEBUG-}" ] || echo "ROT ${1:-$fn} $w" >&2; }

fehlt=0
for t in make node ssh-keygen u-boot docker devcontainer ssh; do command -v "$t" > /dev/null 2>&1 || { melde "$t fehlt im Pfad (Voraussetzung der Gegenprobe)"; fehlt=1; }; done
[ "$fehlt" = 0 ] || exit 1
sys=$arbeit/sys att=$arbeit/att
mkdir -p "$sys" "$att"
for t in bash dirname realpath ls grep stat mkdir cat chmod make env sleep; do ln -s "$(command -v "$t")" "$sys/$t"; done
cat > "$att/attrappe" <<'ATTRAPPE'
#!/usr/bin/env bash
n=${0##*/}
echo "$n $*" >> "$AT_LOG"
[ "$n" != devcontainer ] || echo "CWD $PWD" >> "$AT_LOG"
for v in CLAUDE_CODE_OAUTH_TOKEN ANTHROPIC_API_KEY GIT_TOKEN; do [ -z "${!v+x}" ] || echo "ENV! $n $v" >> "$AT_LOG"; done
mkdir -p "$AT_STATE"
warten() { echo $$ > "$AT_STATE/pid"; echo $PPID > "$AT_STATE/ppid"; exec sleep 5; } # wartet auf ein Signal des Aufrufers
case $n in
  docker) case $1 in
      info) exit "${AT_DK_INFO-0}" ;;
      version) echo "${AT_DK_VER-29.9.1}"; exit "${AT_DK_VEXIT-0}" ;;
      ps) case "$*" in *label=*) x=${AT_DK_LABEL-} ;; *) x=${AT_DK_NAME-} ;; esac; [ -z "$x" ] || echo "$x"; exit 0 ;;
      inspect) echo "${AT_DK_INSPECT-}"; exit 0 ;;
    esac ;;
  devcontainer)
    [ -z "${AT_DC_WARTEN-}" ] || warten
    echo "DC-STDOUT"; echo "DC-STDERR" >&2; echo "${AT_DC_OUT-{\"outcome\":\"success\"\}}"; exit "${AT_DC_EXIT-0}" ;;
  ssh)
    k=$(( $(cat "$AT_STATE/ssh" 2> /dev/null || echo 0) + 1 )); echo $k > "$AT_STATE/ssh"
    case "$*" in *UserKnownHostsFile=/dev/null*) ;; *) echo "[127.0.0.1] key" >> "$HOME/.ssh/known_hosts" ;; esac
    [ -z "${AT_SSH_WARTEN-}" ] || warten
    [ "${AT_SSH_AB-1}" != 0 ] && [ "$k" -ge "${AT_SSH_AB-1}" ] && exit 0
    echo "SSH-ERR-$k" >&2; exit 255 ;;
  u-boot)
    [ "$1" = --version ] && { echo "${AT_UB_AUS-u-boot version ${AT_UB_VER-0.10.0}}"; exit "${AT_UB_EXIT-0}"; }
    k=$(( $(cat "$AT_STATE/k" 2> /dev/null || echo 0) + 1 )); echo $k > "$AT_STATE/k"
    echo "UB-STDOUT-$k"; echo "UB-STDERR-$k" >&2
    [ "${AT_UB_WARN-}" = 1 ] && echo "Warning: Lockerung" >&2
    [ "${AT_UB_WARTEN-}" != "$k" ] || warten
    [ "${AT_UB_ROT-}" = "$k" ] && exit 7
    case $1 in
      init) printf 'schemaVersion: 1\nproject:\n    name: %s\n' "$2" > u-boot.yaml ;;
      generate) mkdir -p .devcontainer
        cat > .devcontainer/devcontainer.json <<'J'
// Kommentar
{"name":"x","image":"i","runArgs":["--privileged"],"mounts":["m1"],"remoteEnv":{"A":"${localEnv:A}"}}
J
        ;;
    esac ;;
esac
exit 0
ATTRAPPE
chmod 755 "$att/attrappe"
cat > "$arbeit/pruef.js" <<'PRUEF'
const [,, a, d, nm, pt, cp, me] = process.argv;
const fs = require("fs");
const rd = (f) => JSON.parse(fs.readFileSync(d + "/.devcontainer/" + f, "utf8").replace(/^\s*\/\/.*$/mg, ""));
const j = rd("devcontainer.json"), o = rd("override.json"), n = j.runArgs.length;
const eq = (x, y) => JSON.stringify(x) === JSON.stringify(y);
const einmal = (p, v) => o.runArgs.filter((x) => x.startsWith(p)).length === 1 && o.runArgs.includes(p + v);
const t = {
  name: () => einmal("--name=", nm),
  cpus: () => einmal("--cpus=", cp),
  speicher: () => einmal("--memory=", me),
  port: () => eq(o.appPort, ["127.0.0.1:" + pt + ":2222"]),
  mount: () => eq(o.mounts, [...j.mounts, "type=bind,source=${localWorkspaceFolder}/ssh/id_ed25519.pub,target=/home/vscode/.ssh/authorized_keys,readonly"]),
  env: () => eq(o.remoteEnv, { ...j.remoteEnv, CLAUDE_CODE_OAUTH_TOKEN: "${localEnv:CLAUDE_CODE_OAUTH_TOKEN}" }),
  rest: () => { const { runArgs, mounts, remoteEnv, ...r1 } = j, { runArgs: a1, mounts: m1, remoteEnv: e1, appPort, ...r2 } = o; return eq(r1, r2) && eq(o.runArgs.slice(0, n), j.runArgs); },
};
process.exit(t[a]() ? 0 : 1);
PRUEF

# --- Rahmen der Fälle -----------------------------------------------------------------------
SEL='.*' sk= root= gpk=$GP_SELF n=0 nodev= fn= w= dev= code= s0=
lauf_cmd=()
P4S='p2-dev-leer|p2-port-abc|p3-fehlt-ssh|p3-u-boot-0.9.0|p5-fremd-datei|p5-elternpfad-fehlt|p6-nur-privat|p6-passphrase|p9-kollision|p4-nichts-geaendert|p4-meldung-form'
neu() {
  local t
  n=$((n + 1)); w=$arbeit/c$n; dev=$w/pgr-cs-dev; nodev=
  mkdir -p "$w/bin" "$w/home/.ssh"; echo k > "$w/home/.ssh/known_hosts"; : > "$w/log"
  for t in docker u-boot devcontainer ssh; do ln -s "$att/attrappe" "$w/bin/$t"; done
  ln -s "$node" "$w/bin/node"; ln -s "$sshkeygen" "$w/bin/ssh-keygen"
}
c() {
  fn=$1; [[ $1 =~ ^($SEL)$ ]] || return 1
  [[ -n ${GP_KIND-} && $1 =~ ^p15-(voraussetzung|sonden-raeumen|rekursion) ]] && return 1
  [ -z "${GP_LISTE-}" ] || echo "$1" >> "$GP_LISTE"
  neu
}
snap() { (cd "$w" && find . ! -path ./bin ! -path './bin/*' ! -path ./st ! -path './st/*' ! -name log ! -name out ! -name err -printf '%p %m %y %l\n' -type f -exec cksum {} +) | sort; }
snapw() { (cd "$1" && find . -printf '%p %m\n' -type f -exec cksum {} + | sort); }
lauf() {
  local b=(PATH="$w/bin:$sys" HOME="$w/home" AT_LOG="$w/log" AT_STATE="$w/st") cmd=("$BASH" "$sk")
  [ -n "$nodev" ] || b+=(PGR_DEV="$dev")
  (( ${#lauf_cmd[@]} )) && cmd=("${lauf_cmd[@]}")
  : > "$w/log"; s0=$(snap)
  (cd "$w" && env -i "${b[@]}" "$@" "${cmd[@]}" > "$w/out" 2> "$w/err"); code=$?
}
gruen() { lauf "$@"; ex 0; grep -q '^container-start: bereit ' "$w/out" || rot; }
# ex <Exit-Code> [<Text>…]: Exit-Code und Texte auf stderr; bei Exit 2 bis 4 gelten Punkt 4 (nichts verändert,
# nur lesende Aufrufe, stdout leer) und die Form der Meldung.
ex() {
  local x=$1 t; shift
  [[ p9-kein-update =~ ^($SEL)$ ]] && { ! grep -q '^docker update' "$w/log" || rot p9-kein-update; }
  [ "$code" = "$x" ] || rot
  for t; do grep -qF -- "$t" "$w/err" || rot; done
  [[ $code =~ ^[234]$ ]] || return 0
  [[ p4-nichts-geaendert =~ ^($SEL)$ ]] && { { [ "$s0" = "$(snap)" ] && [ ! -s "$w/out" ] && ! grep -vqE '^(u-boot --version|docker (info|version.*|ps.*|inspect.*|port.*))$' "$w/log" \
    && { [ "$code" != 2 ] || [ ! -s "$w/log" ]; }; } || rot p4-nichts-geaendert; }
  [[ p4-meldung-form =~ ^($SEL)$ ]] && { ! grep -vqE '^container-start: (eingaben|werkzeuge|zustand): ' "$w/err" || rot p4-meldung-form; }
  return 0
}
zeilen() { [ "$(wc -l < "$w/err")" = "$1" ] || rot; }
aufrufe() { grep -E '^u-boot (init|config|generate|up|doctor)' "$w/log"; }
schl() { rm -f "$dev"/ssh/id_ed25519*; "$sshkeygen" -q -t ed25519 -N "$1" -C fix -f "$dev/ssh/id_ed25519" < /dev/null; }
mt() { mkdir -p "$w/t/harness/mk" "$w/t/tools/harness"; cp "$root/Makefile" "$w/t/"; cp "$root/harness/mk/container-start.mk" "$w/t/harness/mk/"; cp "$sk" "$w/t/tools/harness/"; }
mk() { local fl=(); while [[ ${1-} == -* ]]; do fl+=("$1"); shift; done; lauf_cmd=(make -C "$w/t" "${fl[@]}" container-start); lauf "$@"; lauf_cmd=(); }
nichts() { [ ! -s "$w/log" ] && [ ! -e "$dev" ] && [ "$s0" = "$(snap)" ] || rot; }
ein() { # <Fall> <Variable> <Text> <Zuweisung>…: ein Verstoß, genau eine Zeile
  c "$1" || return 0
  local nm=$2 tx=$3 a; shift 3; a=("${@//@W@/$w}")
  lauf "${a[@]}"; ex 2 "eingaben: $nm " "$tx"; zeilen 1
}
hoere() { # Listener auf einem freien Port: lp (Port), lk (Prozess)
  coproc LST { exec "$node" -e 'const s=require("net").createServer(()=>{}).listen(0,"127.0.0.1",()=>console.log(s.address().port))' < /dev/null; }
  read -r lp <&"${LST[0]}"; lk=$LST_PID
}
up_ok() { grep -qx 'devcontainer up --workspace-folder . --override-config .devcontainer/override.json' "$w/log"; }
zerstoert() { grep -qE '^docker (stop|rm|rmi|kill|volume|network|image|container|system)|prune' "$w/log"; }
sig_lauf() { # <Schritt> <Signal>: sendet das Signal, sobald das Kind des Schrittes wartet; v ist der Stand des Verzeichnisses vor start
  local e=() p; v=
  case $1 in
    geruest) e=(AT_UB_WARTEN=3) ;;
    start) gruen; v=$(snapw "$dev"); e=(AT_DC_WARTEN=1 AT_DK_LABEL=pgr-dev AT_DK_NAME=pgr-dev 'AT_DK_INSPECT=exited 127.0.0.1:2222 ') ;;
    warten) e=(AT_SSH_WARTEN=1) ;;
  esac
  ( for p in $(seq 100); do [ -s "$w/st/ppid" ] && break; sleep 0.1; done; kill -"$2" "$(cat "$w/st/ppid")" 2> /dev/null ) &
  lauf "${e[@]}"; wait
}
dk_weg() { # entfernt, was eine rote Sonde hinterließ, per Name
  local v
  [[ ${pl[2]} == pgr-cs-* ]] && docker rm -fv "${pl[2]}" > /dev/null 2>&1
  for v in "${pl[3]}" "${pl[@]:5}"; do [[ $v == pgr-cs-* || $v =~ ^[0-9a-f]{64}$ ]] && docker volume rm "$v" > /dev/null 2>&1; done
  [[ ${pl[4]} == vsc-pgr-cs-* ]] && docker rmi "${pl[4]}" > /dev/null 2>&1
}
sonde() { # <Modus> [<Signal> [<Zeilen der Sondendatei>]]: startet die Gegenprobe im Sondenmodus; setzt pcode und psd
  local d=$w/sonde k
  if [ -n "${2-}" ]; then
    ( for k in $(seq 100); do [ -s "$d" ] && [ "$(wc -l < "$d")" -ge "${3:-2}" ] && break; sleep 0.1; done; kill -"$2" "$(sed -n 2p "$d")" ) &
  fi
  { GP_PROBE=$1 GP_PROBE_DATEI=$d bash "$gpk"; } 2> /dev/null; pcode=$?
  wait
  psd=$(sed -n 1p "$d"); [[ ${psd##*/} == pgr-cs-* ]] && PSDS+=("$psd")
}
# kind <Werkzeug> [<Muster>]: die Gegenprobe (gpk) in einem Pfad ohne <Werkzeug>
kind() {
  local f=$arbeit/ohne-$1 p b
  [ -z "${GP_KIND-}" ] || { melde "Rekursion: eine Kind-Gegenprobe ruft kind erneut auf"; exit 1; }
  if [ ! -d "$f" ]; then
    mkdir "$f"; declare -A alle=()
    while IFS= read -r p; do b=${p##*/}; [ "$b" = "$1" ] || [ -n "${alle[$b]-}" ] || alle[$b]=$p; done < <(for p in ${PATH//:/ }; do find "$p/" -maxdepth 1 \( -type f -o -type l \) -perm -u+x 2> /dev/null; done)
    ln -s -t "$f" "${alle[@]}"
  fi
  (cd "$w" && timeout 300 env -i PATH="$f" HOME="$w" TMPDIR="$w" GP_KIND=1 GP_QUELLE="$quelle" GP_MUT="${2:-^\$}" GP_LISTE="$w/liste" "$BASH" "$gpk" > "$w/out" 2> "$w/err"); code=$?
}

# Gruppe Docker (Punkte 9 bis 11 und 15): echtes Docker und devcontainer, höchstens ein Container des Gerüsts zugleich.
dk() { (cd "$arbeit" && env PGR_DEV="$dd" PGR_NAME="$nm" PGR_SSH_PORT="$dp" PGR_CPUS=2 PGR_MEM=2g PGR_WARTEN=20 "$@" "$BASH" "$sk" > "$arbeit/dk.out" 2> "$arbeit/dk.err"); code=$?; }
ins() { docker inspect --format "$1" "$nm"; }
sel() { [[ $1 =~ ^($SEL)$ ]]; }
dk_gruppe() {
  local id1= t1 v1 d2
  for t1 in p9-docker-erststart p9-docker-zweiter-aufruf p9-docker-stop-start p9-docker-rm-neustart p10-docker-port-belegt p11-docker-bereit \
    p11-docker-anderer-schluessel p15-docker-praefix p15-docker-grenzen; do sel $t1 && id1=x; done
  [ -n "$id1" ] || return 0
  [ -z "${GP_KIND-}" ] || return 0
  dd=$arbeit/$nm dp=$(freier_port) || { melde "kein freier Port"; return 0; }
  DKC=("$nm") DKV=() dkb=1
  dk
  sel p11-docker-bereit && { { [ "$code" = 0 ] && grep -qx "container-start: bereit $nm 127.0.0.1:$dp" "$arbeit/dk.out"; } || rot p11-docker-bereit; }
  sel p9-docker-erststart && { { [ "$(ins '{{.Name}}')" = "/$nm" ] && [ "$(ins '{{index .Config.Labels "devcontainer.local_folder"}}')" = "$dd" ] && [ "$(ins '{{.State.Running}}')" = true ]; } || rot p9-docker-erststart; }
  sel p15-docker-praefix && { { [[ $nm == pgr-cs-* && ${dd##*/} == pgr-cs-* && $(ins '{{.Config.Image}}') == vsc-pgr-cs-* ]]; } || rot p15-docker-praefix; }
  sel p15-docker-grenzen && { [ "$(ins '{{.HostConfig.NanoCpus}} {{.HostConfig.Memory}}')" = '2000000000 2147483648' ] || rot p15-docker-grenzen; }
  id1=$(ins '{{.Id}}'); t1=$(ins '{{.State.StartedAt}}'); v1=$(ins '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}')
  sel p9-docker-zweiter-aufruf && { dk; { [ "$code" = 0 ] && [ "$(ins '{{.State.StartedAt}}')" = "$t1" ] && [ "$(ins '{{.Id}}')" = "$id1" ]; } || rot p9-docker-zweiter-aufruf; }
  sel p10-docker-port-belegt && { d2=$dd-zwei; DKC+=("$nm-zwei"); dd=$d2 dk PGR_NAME="$nm-zwei"
    { [ "$code" = 4 ] && grep -qF "der Port $dp ist belegt" "$arbeit/dk.err" && [ -z "$(docker ps -a --filter "name=^/$nm-zwei\$" --format '{{.Names}}')" ] && [ ! -e "$d2" ]; } || rot p10-docker-port-belegt; }
  sel p11-docker-anderer-schluessel && { cp -a "$dd/ssh" "$arbeit/ssh.sav"; rm -f "$dd"/ssh/id_ed25519*; "$sshkeygen" -q -t ed25519 -N '' -f "$dd/ssh/id_ed25519" < /dev/null
    dk PGR_WARTEN=2; { [ "$code" = 5 ] && grep -qF 'container-start: warten: ' "$arbeit/dk.err"; } || rot p11-docker-anderer-schluessel
    rm -rf "$dd/ssh"; mv "$arbeit/ssh.sav" "$dd/ssh"; }
  sel p9-docker-stop-start && { docker stop "$nm" > /dev/null; dk; { [ "$code" = 0 ] && [ "$(ins '{{.Id}}')" = "$id1" ] && [ "$(ins '{{.State.Running}}')" = true ]; } || rot p9-docker-stop-start; }
  sel p9-docker-rm-neustart && { docker rm -fv "$nm" > /dev/null; dk
    { [ "$code" = 0 ] && [ "$(ins '{{.Id}}')" != "$id1" ] && [ "$(ins '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}')" = "$v1" ]; } || rot p9-docker-rm-neustart; }
  DKV+=($v1); dk_raeumen; DKC=() DKV=(); rm -rf -- "$dd" "$dd-zwei"
}

cases() {
  local t k v a p
  # Punkt 1: Gegenstand und Ort
  c p1-keine-zerstoerung && { gruen; ! grep -qE '^docker (stop|rm|rmi|kill|volume|network|image|container|system)|prune' "$w/log" || rot; }
  c p1-keine-zerstoerung-kollision && { lauf AT_DK_NAME=pgr-dev; ex 4; ! zerstoert || rot; }
  c p1-repo-unveraendert && { cp -a "$root.rein" "$w/r"; a=$sk; sk=$w/r/tools/harness/container-start.sh; v=$(snapw "$w/r"); gruen; sk=$a; [ "$v" = "$(snapw "$w/r")" ] || rot; }
  c p1-home-unveraendert && { v=$(snapw "$w/home"); gruen; [ "$v" = "$(snapw "$w/home")" ] || rot; }
  c p1-nur-unter-pgr-dev && { gruen; [ "$(ls -A "$w" | tr '\n' ' ')" = 'bin err home log out pgr-cs-dev st ' ] || rot; }

  # Punkt 2: Eingaben, ein Fall je Form
  c p2-dev-nicht-gesetzt && { nodev=1; lauf; ex 2 'PGR_DEV ist nicht gesetzt'; zeilen 1; }
  ein p2-dev-leer PGR_DEV 'ist leer' PGR_DEV=
  ein p2-dev-relativ PGR_DEV 'ist nicht absolut' PGR_DEV=rel/pgr-cs-dev
  ein p2-dev-leerzeichen PGR_DEV 'enthält ein Zeichen' 'PGR_DEV=@W@/a b'
  ein p2-dev-tabulator PGR_DEV 'enthält ein Zeichen' $'PGR_DEV=@W@/a\tb'
  ein p2-dev-zeilenumbruch PGR_DEV 'enthält ein Zeichen' $'PGR_DEV=@W@/a\nb'
  ein p2-dev-nicht-ascii PGR_DEV 'enthält ein Zeichen' 'PGR_DEV=@W@/ä'
  ein p2-dev-segment-punktpunkt PGR_DEV 'hat das Segment ..' 'PGR_DEV=@W@/x/../pgr-cs-dev'
  ein p2-dev-segment-punkt PGR_DEV 'hat das Segment .' 'PGR_DEV=@W@/./pgr-cs-dev'
  ein p2-dev-leeres-segment PGR_DEV 'hat ein leeres Segment' 'PGR_DEV=@W@//pgr-cs-dev'
  ein p2-dev-slash-am-ende PGR_DEV 'endet auf /' 'PGR_DEV=@W@/pgr-cs-dev/'
  ein p2-dev-basis-gross PGR_DEV 'Basisnamen' 'PGR_DEV=@W@/Dev'
  ein p2-dev-basis-unterstrich PGR_DEV 'Basisnamen' 'PGR_DEV=@W@/pgr_dev'
  ein p2-name-leer PGR_NAME 'ungültig' PGR_NAME=
  ein p2-name-gross PGR_NAME 'ungültig' PGR_NAME=Pgr
  for p in leer: abc:abc 02222:02222 plus22:+22 minus22:-22 2plus2:2+2 22minus:22- 0:0 65536:65536; do ein "p2-port-${p%%:*}" PGR_SSH_PORT 'keine ganze Zahl' "PGR_SSH_PORT=${p#*:}"; done
  for p in leer: abc:abc ab:ab 0:0 601:601; do ein "p2-warten-${p%%:*}" PGR_WARTEN 'keine ganze Zahl' "PGR_WARTEN=${p#*:}"; done
  for p in PGR_SSH_PORT=1 PGR_SSH_PORT=65535 PGR_WARTEN=1 PGR_WARTEN=600; do c "p2-grenze-${p//=/-}" && gruen "$p"; done
  c p2-vorgaben && { gruen; v=$dev/.devcontainer/override.json
    for a in --name=pgr-dev --cpus=8 --memory=8g 127.0.0.1:2222:2222; do grep -qF -- "$a" "$v" || rot; done; }
  c p2-alle-gemeldet && { lauf PGR_DEV=rel PGR_NAME=X PGR_SSH_PORT=0 PGR_WARTEN=0; ex 2; zeilen 4; }
  c p2-reihenfolge && { lauf PGR_DEV=rel PGR_NAME=X PGR_SSH_PORT=0 PGR_WARTEN=0; ex 2
    sed -n 's/^container-start: eingaben: \(PGR_[A-Z_]*\) .*/\1/p' "$w/err" \
      | awk 'BEGIN{r["PGR_DEV"]=1;r["PGR_NAME"]=2;r["PGR_SSH_PORT"]=3;r["PGR_WARTEN"]=4} {if (r[$1]<p) b=1; p=r[$1]} END{exit b}' || rot; }

  # Punkt 3: Voraussetzungen am Host
  for t in docker u-boot devcontainer node ssh ssh-keygen; do
    c "p3-fehlt-$t" && { rm "$w/bin/$t"; lauf; ex 3 "container-start: werkzeuge: $t fehlt im Pfad"; zeilen 1; }
  done
  c p3-alle-fehlen && { rm "$w"/bin/*; lauf; ex 3
    [ "$(sed 's/^container-start: werkzeuge: //;s/ fehlt im Pfad$//' "$w/err" | tr '\n' ' ')" = 'docker u-boot devcontainer node ssh ssh-keygen ' ] || rot; }
  c p3-u-boot-0.9.0 && { lauf AT_UB_VER=0.9.0; ex 3 'container-start: werkzeuge: u-boot 0.9.0 ist kleiner als 0.10.0'; zeilen 1; }
  for v in 0.10.0 0.100.0 1.0.0; do c "p3-u-boot-$v" && gruen AT_UB_VER=$v; done
  c p3-u-boot-unlesbar && { lauf 'AT_UB_AUS=u-boot version dev'; ex 3 'container-start: werkzeuge: u-boot: Version unlesbar'; zeilen 1; }
  c p3-u-boot-exit1 && { lauf AT_UB_EXIT=1; ex 3 'container-start: werkzeuge: u-boot: Version unlesbar'; zeilen 1; }
  c p3-docker-unlesbar && { lauf AT_DK_VER=dev; ex 3 'container-start: werkzeuge: docker: Version des Servers unlesbar'; zeilen 1; }
  c p3-docker-exit1 && { lauf AT_DK_VEXIT=1; ex 3 'container-start: werkzeuge: docker: Version des Servers unlesbar'; zeilen 1; }
  c p3-docker-daemon && { lauf AT_DK_INFO=1; ex 3 'container-start: werkzeuge: docker: der Daemon ist nicht erreichbar'; zeilen 1; }
  c p3-docker-23 && { lauf AT_DK_VER=23.0.6; ex 3 'container-start: werkzeuge: docker: der Server 23.0.6 ist älter als 24'; zeilen 1; }
  c p3-docker-24 && gruen AT_DK_VER=24.0.0

  # Punkt 5: das Verzeichnis
  c p5-im-repo && { lauf PGR_DEV="$root/pgr-cs-dev"; ex 4 'liegt im Repo'; }
  c p5-wurzel-selbst && { lauf PGR_DEV="$root"; ex 4 'liegt im Repo'; }
  c p5-im-repo-symlink && { ln -s "$root" "$w/lnk"; lauf PGR_DEV="$w/lnk/pgr-cs-dev"; ex 4 'liegt im Repo'; }
  c p5-praefix-ohne-slash && gruen PGR_DEV="$root-zwei"
  c p5-elternpfad-fehlt && { lauf PGR_DEV="$w/nirgends/pgr-cs-dev"; ex 4 'Elternpfad'; [ ! -e "$w/nirgends" ] || rot; }
  c p5-elternpfad-ist-datei && { : > "$w/datei"; lauf PGR_DEV="$w/datei/pgr-cs-dev"; ex 4 'Elternpfad'; }
  c p5-kein-verzeichnis && { : > "$dev"; lauf; ex 4 'ist kein Verzeichnis'; }
  c p5-fehlt && { gruen; [ -d "$dev" ] || rot; }
  c p5-leer && { mkdir "$dev"; gruen; aufrufe | grep -q '^u-boot init ' || rot; }
  c p5-vorhanden-ohne-aufruf && { gruen; gruen; ! aufrufe | grep -q . || rot; }
  c p5-vorhanden-unveraendert && { gruen; v=$(snap | grep -v override.json); gruen PGR_CPUS=3; [ "$v" = "$(snap | grep -v override.json)" ] || rot; }
  c p5-fortsetzen && { mkdir "$dev"; printf 'schemaVersion: 1\nproject:\n    name: pgr-cs-dev\n' > "$dev/u-boot.yaml"; gruen
    [ "$(aufrufe | grep -c '^u-boot config set ')" = 6 ] && aufrufe | grep -q '^u-boot generate devcontainer$' || rot; }
  c p5-fremd-datei && { mkdir "$dev"; echo x > "$dev/datei"; lauf; ex 4 "$dev ist fremd"; }
  c p5-fremd-anderer-name && { mkdir "$dev"; printf 'project:\n    name: pgr-cs-dev-other\n' > "$dev/u-boot.yaml"; lauf; ex 4 "$dev ist fremd"; }

  # Punkt 6: Schlüsselpaar
  c p6-modus-ordner && { gruen; [ "$(stat -c %a "$dev/ssh")" = 700 ] || rot; }
  c p6-modus-privat && { gruen; [ "$(stat -c %a "$dev/ssh/id_ed25519")" = 600 ] || rot; }
  c p6-kommentar && { gruen; [ "$(cut -d' ' -f3 "$dev/ssh/id_ed25519.pub")" = pgr-cs-dev ] || rot; }
  c p6-typ && { gruen; [ "$(cut -d' ' -f1 "$dev/ssh/id_ed25519.pub")" = ssh-ed25519 ] || rot; }
  c p6-ohne-passphrase && { gruen; "$sshkeygen" -y -P '' -f "$dev/ssh/id_ed25519" > /dev/null 2>&1 < /dev/null || rot; }
  c p6-gueltig-unveraendert && { gruen; v=$(cksum "$dev"/ssh/*); gruen; [ "$v" = "$(cksum "$dev"/ssh/*)" ] || rot; }
  c p6-kommentar-ignoriert && { gruen; sed -i 's/ pgr-cs-dev$/ anderer/' "$dev/ssh/id_ed25519.pub"; gruen; }
  c p6-modus-ordner-vorhanden && { gruen; rm -rf "$dev/ssh"; mkdir -m 755 "$dev/ssh"; echo x > "$dev/ssh/andere"; gruen; [ "$(stat -c %a "$dev/ssh")" = 700 ] && [ "$(cat "$dev/ssh/andere")" = x ] || rot; }
  c p6-ssh-ist-datei && { gruen; rm -rf "$dev/ssh"; : > "$dev/ssh"; lauf; ex 4 "$dev/ssh ist kein Verzeichnis"; }
  c p6-ssh-ist-symlink && { gruen; rm -rf "$dev/ssh"; : > "$w/datei"; ln -s "$w/datei" "$dev/ssh"; lauf; ex 4 "$dev/ssh ist kein Verzeichnis"; }
  c p6-nur-privat && { gruen; rm "$dev/ssh/id_ed25519.pub"; lauf; ex 4 'nur eine der beiden Dateien'; }
  c p6-nur-oeffentlich && { gruen; rm "$dev/ssh/id_ed25519"; lauf; ex 4 'nur eine der beiden Dateien'; }
  c p6-ordner-750 && { gruen; chmod 750 "$dev/ssh"; lauf; ex 4 '(a)'; }
  c p6-privat-640 && { gruen; chmod 640 "$dev/ssh/id_ed25519"; lauf; ex 4 '(a)'; }
  c p6-passphrase && { gruen; schl x; lauf; ex 4 '(b)'; }
  c p6-leere-datei && { gruen; : > "$dev/ssh/id_ed25519"; lauf; ex 4 '(b)'; }
  c p6-pub-anderer && { gruen; "$sshkeygen" -q -t ed25519 -N '' -f "$w/anderer" < /dev/null; cp "$w/anderer.pub" "$dev/ssh/id_ed25519.pub"; lauf; ex 4 '(c)'; }

  # Punkt 7: u-boot
  A=('u-boot init pgr-cs-dev --devcontainer --sandbox --no-git' 'u-boot config set devcontainer.features.node.enabled true'
    'u-boot config set devcontainer.features.sshd.enabled true' 'u-boot config set devcontainer.features.claude-code.enabled true'
    'u-boot config set devcontainer.sandbox.nestedRuntime docker' 'u-boot config set devcontainer.sandbox.egress.enabled true'
    'u-boot config set devcontainer.sandbox.egress.allow gcr.io,pkg-containers.githubusercontent.com' 'u-boot generate devcontainer')
  for k in 1 2 3 4 5 6 7 8; do
    c "p7-aufruf-$k" && { gruen; [ "$(aufrufe | sed -n "${k}p")" = "${A[k-1]}" ] || rot; }
    c "p7-fehler-$k" && { lauf AT_UB_ROT=$k; ex 5 'container-start: geruest: ' "UB-STDERR-$k"; [ "$(aufrufe | wc -l)" = "$k" ] || rot; }
  done
  c p7-reihenfolge && { gruen; [ "$(aufrufe | awk '{print $2}' | tr '\n' ' ')" = 'init config config config config config config generate ' ] || rot; }
  c p7-warnung && { gruen AT_UB_WARN=1; grep -qF 'Warning: Lockerung' "$w/err" || rot; }
  c p7-ausgabe-stdout && { gruen; for k in 1 2 3 4 5 6 7 8; do grep -qx "UB-STDOUT-$k" "$w/out" || rot; done; }
  c p7-ausgabe-stderr && { gruen; for k in 1 2 3 4 5 6 7 8; do grep -qx "UB-STDERR-$k" "$w/err" || rot; done; }
  c p7-kein-up-doctor && { gruen; ! grep -qE '^u-boot (up|doctor)' "$w/log" || rot; }
  c p7-fortsetzen-ohne-init && { mkdir "$dev"; printf 'project:\n    name: pgr-cs-dev\n' > "$dev/u-boot.yaml"; gruen; ! aufrufe | grep -q '^u-boot init' || rot; }
  c p7-im-verzeichnis && { gruen; [ -f "$dev/u-boot.yaml" ] && [ -f "$dev/.devcontainer/devcontainer.json" ] || rot; }
  c p7-echt-u-boot && { ln -sf "$realub" "$w/bin/u-boot"; gruen; v=$dev/u-boot.yaml
    grep -qx '    name: pgr-cs-dev' "$v" && grep -qx '        nestedRuntime: docker' "$v" && grep -qx '                - gcr.io' "$v" \
      && grep -qx '                - pkg-containers.githubusercontent.com' "$v" && [ -f "$dev/.devcontainer/devcontainer.json" ] \
      && awk '/^        [a-zA-Z-]+:$/{f=$1} /^            enabled: true$/{c[f]=1} END{exit !(c["claude-code:"] && c["node:"] && c["sshd:"])}' "$v" || rot; }

  # Punkt 8: mk-override.sh und override.json
  for a in name cpus speicher port mount env rest; do
    c "p8-$a" && { gruen PGR_NAME=pgr-cs-x PGR_SSH_PORT=2999 PGR_CPUS=2 PGR_MEM=2g; "$node" "$arbeit/pruef.js" "$a" "$dev" pgr-cs-x 2999 2 2g || rot; }
  done
  c p8-skript-modus && { gruen; [ "$(stat -c %a "$dev/mk-override.sh")" = 755 ] || rot; }
  c p8-vorhanden-bleibt && { gruen; printf '#!/bin/bash\necho x\n' > "$dev/mk-override.sh"; v=$(cksum < "$dev/mk-override.sh"); gruen; [ "$v" = "$(cksum < "$dev/mk-override.sh")" ] || rot; }
  c p8-vorhanden-laeuft && { gruen; printf '#!/bin/bash\n: > "$(dirname "$0")/lief"\n' > "$dev/mk-override.sh"; gruen; [ -e "$dev/lief" ] || rot; }
  c p8-skript-fehler && { gruen; printf '#!/bin/bash\nexit 9\n' > "$dev/mk-override.sh"; lauf; ex 5 'container-start: override: mk-override.sh endete mit Exit 9'; }
  for v in PGR_NAME=pgr-dev PGR_SSH_PORT=2222 PGR_CPUS=8 PGR_MEM=8g; do
    c "p8-skript-env-${v%%=*}" && { gruen; printf '#!/bin/bash\nenv > "$(dirname "$0")/envdump"\n' > "$dev/mk-override.sh"
      gruen; grep -qx "$v" "$dev/envdump" 2> /dev/null || rot; }
  done
  for v in CLAUDE_CODE_OAUTH_TOKEN ANTHROPIC_API_KEY GIT_TOKEN; do
    c "p8-umgebung-$v" && { gruen; printf '#!/bin/bash\nenv > "$(dirname "$0")/envdump"\n' > "$dev/mk-override.sh"
      gruen CLAUDE_CODE_OAUTH_TOKEN=SENTINEL-1 ANTHROPIC_API_KEY=SENTINEL-2 GIT_TOKEN=SENTINEL-3
      { ! grep -q "^ENV! .* $v\$" "$w/log" && [ -s "$dev/envdump" ] && ! grep -q "^$v=" "$dev/envdump"; } || rot; }
  done
  c p8-keine-werte && { gruen CLAUDE_CODE_OAUTH_TOKEN=SENTINEL-1 ANTHROPIC_API_KEY=SENTINEL-2 GIT_TOKEN=SENTINEL-3; ! grep -rq SENTINEL "$w" || rot; }

  # Punkt 9: Container und zweiter Aufruf
  DK=(AT_DK_LABEL=pgr-dev AT_DK_NAME=pgr-dev)
  GS=("${DK[@]}" 'AT_DK_INSPECT=exited 127.0.0.1:2222 ')
  c p9-kollision && { lauf AT_DK_NAME=pgr-dev; ex 4 'der Container pgr-dev trägt das Label des Gerüsts nicht (Kollision)'; }
  c p9-label-anderer-name && { lauf AT_DK_LABEL=pgr-anders; ex 4 'heißt pgr-anders, nicht pgr-dev'; }
  c p9-adresse-weicht-ab && { lauf "${DK[@]}" 'AT_DK_INSPECT=running 0.0.0.0:2222 '; ex 4 'als 0.0.0.0:2222'; }
  c p9-port-weicht-ab && { lauf "${DK[@]}" 'AT_DK_INSPECT=running 127.0.0.1:2999 '; ex 4 'als 127.0.0.1:2999'; }
  c p9-laeuft && { hoere; lauf PGR_SSH_PORT=$lp "${DK[@]}" "AT_DK_INSPECT=running 127.0.0.1:$lp "; kill "$lk"; ex 0
    { ! grep -q '^devcontainer' "$w/log" && [ ! -e "$dev/.devcontainer/override.json" ] && grep -q '^container-start: bereit ' "$w/out"; } || rot; }
  for t in exited created; do
    c "p9-$t" && { lauf "${DK[@]}" "AT_DK_INSPECT=$t 127.0.0.1:2222 "; ex 0; { grep -q '^devcontainer up ' "$w/log" && [ -f "$dev/.devcontainer/override.json" ] && ! zerstoert; } || rot; }
  done
  c p9-kein-update && { lauf PGR_MEM=3g "${DK[@]}" 'AT_DK_INSPECT=running 127.0.0.1:2222 '; ex 0; }
  c p9-keiner && { gruen; grep -q '^devcontainer up ' "$w/log" || rot; }
  c p9-up-aufruf && { gruen "${GS[@]}"; up_ok || rot; }
  c p9-up-verzeichnis && { gruen "${GS[@]}"; grep -qx "CWD $dev" "$w/log" || rot; }
  c p9-up-ohne-outcome && { lauf "${GS[@]}" AT_DC_OUT=nichts; ex 5 'container-start: start: ' 'outcome'; }
  c p9-up-exit1 && { lauf "${GS[@]}" AT_DC_EXIT=1; ex 5 'container-start: start: devcontainer up endete mit Exit 1'; }
  c p9-ausgabe-stdout && { gruen "${GS[@]}"; grep -qx DC-STDOUT "$w/out" || rot; }
  c p9-ausgabe-stderr && { gruen "${GS[@]}"; grep -qx DC-STDERR "$w/err" || rot; }
  c p9-umgebung && { gruen "${GS[@]}" CLAUDE_CODE_OAUTH_TOKEN=SENTINEL-1 ANTHROPIC_API_KEY=SENTINEL-2 GIT_TOKEN=SENTINEL-3; ! grep -q '^ENV! ' "$w/log" || rot; }

  # Punkt 10: Port belegt
  c p10-belegt && { hoere; lauf PGR_SSH_PORT=$lp; kill "$lk"; ex 4 "der Port $lp ist belegt"; }
  c p10-belegt-gestoppt && { hoere; lauf PGR_SSH_PORT=$lp "${DK[@]}" "AT_DK_INSPECT=exited 127.0.0.1:$lp "; kill "$lk"; ex 4 "der Port $lp ist belegt"; }

  # Punkt 11: Bereitschaft
  c p4-stdout-nur-bereit && { gruen; ! grep -qvE '^(UB-STDOUT-[0-9]+|DC-STDOUT|\{"outcome":"success"\}|container-start: bereit .*)$' "$w/out" || rot; }
  c p11-bereit-meldung && { gruen PGR_NAME=pgr-cs-x PGR_SSH_PORT=2999; [ "$(tail -n 1 "$w/out")" = 'container-start: bereit pgr-cs-x 127.0.0.1:2999' ] || rot; }
  c p11-erster-versuch && { gruen; [ "$(grep -c '^ssh ' "$w/log")" = 1 ] || rot; }
  c p11-nach-drei && { gruen AT_SSH_AB=3; [ "$(grep -c '^ssh ' "$w/log")" = 3 ] || rot; }
  c p11-nie && { v=$SECONDS; lauf PGR_WARTEN=2 AT_SSH_AB=0; v=$((SECONDS - v)); a=$(grep -c '^ssh ' "$w/log"); ex 5 'container-start: warten: ' "SSH-ERR-$a"
    { [ "$v" -ge 2 ] && [ "$a" -ge 2 ] && [ "$a" -le 3 ] && ! zerstoert && ! grep -q 'SSH-ERR-1' "$w/err"; } || rot; }
  c p11-vorgabe-warten && { lauf AT_SSH_AB=0; a=$(grep -c '^ssh ' "$w/log"); { [ "$a" -ge 60 ] && [ "$a" -le 61 ]; } || rot; }
  for p in schluessel identitiesonly batchmode connecttimeout stricthostkeychecking userknownhostsfile loglevel port benutzer-adresse befehl; do
    c "p11-opt-$p" && { gruen PGR_SSH_PORT=2999
      case $p in schluessel) t=" -i $dev/ssh/id_ed25519 " ;; identitiesonly) t=' -o IdentitiesOnly=yes ' ;; batchmode) t=' -o BatchMode=yes ' ;;
        connecttimeout) t=' -o ConnectTimeout=5 ' ;; stricthostkeychecking) t=' -o StrictHostKeyChecking=no ' ;; userknownhostsfile) t=' -o UserKnownHostsFile=/dev/null ' ;;
        loglevel) t=' -o LogLevel=ERROR ' ;; port) t=' -p 2999 ' ;; benutzer-adresse) t=' vscode@127.0.0.1 ' ;; befehl) t=' true ' ;; esac
      [[ " $(grep -m1 '^ssh ' "$w/log") " == *"$t"* ]] || rot; }
  done
  c p11-known-hosts-unveraendert && { v=$(snapw "$w/home"); gruen; [ "$v" = "$(snapw "$w/home")" ] || rot; }

  # Punkt 12: Abbruch
  for g in INT:130 TERM:143; do for t in geruest start warten; do
    c "p12-sig${g%%:*}-$t" && { sig_lauf "$t" "${g%%:*}"; ex "${g#*:}" "container-start: $t: abgebrochen"
      { ! kill -0 "$(cat "$w/st/pid")" 2> /dev/null && ! zerstoert && { [ -z "$v" ] || [ "$v" = "$(snapw "$dev")" ]; }; } || { rot; kill -KILL "$(cat "$w/st/pid")" 2> /dev/null; }; }
  done; done
  c p12-wiederholung && { sig_lauf geruest INT; gruen; ! aufrufe | grep -q '^u-boot init' && [ -f "$dev/.devcontainer/devcontainer.json" ] || rot; kill -KILL "$(cat "$w/st/pid")" 2> /dev/null; }

  # Punkt 13: Reihenfolge der Gründe
  c p13-eingabe-vor-werkzeug && { rm "$w/bin/u-boot"; lauf PGR_NAME=X; ex 2 'eingaben: PGR_NAME '; ! grep -q werkzeuge "$w/err" || rot; }
  c p13-werkzeug-vor-repo && { rm "$w/bin/docker"; lauf PGR_DEV="$root/pgr-cs-dev"; ex 3 'docker fehlt'; }
  c p13-repo-vor-fremd && { mkdir "$root/pgr-cs-xdir"; echo x > "$root/pgr-cs-xdir/datei"; lauf PGR_DEV="$root/pgr-cs-xdir"; ex 4 'liegt im Repo'; ! grep -q 'ist fremd' "$w/err" || rot; }
  c p13-fremd-vor-schluessel && { mkdir -p "$w/f/ssh"; echo x > "$w/f/datei"; "$sshkeygen" -q -t ed25519 -N '' -f "$w/f/ssh/id_ed25519" < /dev/null; rm "$w/f/ssh/id_ed25519.pub"
    lauf PGR_DEV="$w/f"; ex 4 'ist fremd'; ! grep -q 'nur eine' "$w/err" || rot; }
  c p13-schluessel-vor-kollision && { gruen; rm "$dev/ssh/id_ed25519.pub"; lauf AT_DK_NAME=pgr-dev; ex 4 'nur eine der beiden Dateien'; ! grep -q Kollision "$w/err" || rot; }
  c p13-kollision-vor-port && { hoere; lauf PGR_SSH_PORT=$lp AT_DK_NAME=pgr-dev; kill "$lk"; ex 4 '(Kollision)'; ! grep -q belegt "$w/err" || rot; }

  # Punkt 14: Aufrufarten von make
  c p14-rezept && { v=$(awk '/^container-start:/{f=1;next} f&&/^\t/{print;next} f{exit}' "$root/harness/mk/container-start.mk")
    [ "$(printf '%s\n' "$v" | wc -l)" = 1 ] && [[ $v != *'$(MAKE)'* && ! $v =~ ^[[:space:]]*[@-]*\+ ]] || rot; }
  c p14-phony && { grep -qE '^\.PHONY:.* container-start( |$)' "$root/harness/mk/container-start.mk" && grep -qE '^container-start:[[:space:]]*##' "$root/harness/mk/container-start.mk" || rot; }
  for k in n q t; do
    c "p14-$k" && { mt; mk "-$k"; nichts; [ "$k" != n ] || grep -qF container-start.sh "$w/out" || rot; [ ! -e "$w/t/container-start" ] || rot; }
    c "p14-makeflags-$k" && { mt; mk "MAKEFLAGS=-$k"; nichts; [ ! -e "$w/t/container-start" ] || rot; }
  done
  c p14-n-i && { mt; mk -n -i; nichts; }
  c p14-makeflags-n-i && { mt; mk MAKEFLAGS=-in; nichts; }
  for p in i:-i makeflags-i:MAKEFLAGS=-i; do
    c "p14-${p%%:*}" && { mt; mk "${p#*:}" PGR_NAME=X; [ "$code" = 0 ] && grep -qF 'container-start: eingaben: PGR_NAME ' "$w/err" && ! grep -q bereit "$w/out" "$w/err" || rot; }
    c "p14-${p%%:*}-start" && { mt; mk "${p#*:}" AT_DC_EXIT=1; [ "$code" = 0 ] && grep -qF 'container-start: start: ' "$w/err" && ! grep -q bereit "$w/out" "$w/err" || rot; }
    c "p14-${p%%:*}-warten" && { mt; mk "${p#*:}" PGR_WARTEN=1 AT_SSH_AB=0; [ "$code" = 0 ] && grep -qF 'container-start: warten: ' "$w/err" && ! grep -q bereit "$w/out" "$w/err" || rot; }
  done
  for p in k:-k j:-j4 makeflags-k:MAKEFLAGS=-k makeflags-j:MAKEFLAGS=-j4; do
    c "p14-${p%%:*}" && { mt; mk "${p#*:}"; [ "$code" = 0 ] && grep -q '^container-start: bereit ' "$w/out" || rot; }
  done

  # Punkt 15: Rahmen — Aufräumen und Präfix
  c p15-praefix && { sonde ende; [[ ${psd##*/} == pgr-cs-* ]] || rot; }
  c p15-aufraeumen-ende && { sonde ende; { [ "$pcode" = 0 ] && [ ! -e "$psd" ]; } || rot; }
  c p15-aufraeumen-int && { sonde warten INT; { [ "$pcode" = 130 ] && [ ! -e "$psd" ]; } || rot; }
  c p15-aufraeumen-term && { sonde warten TERM; { [ "$pcode" = 143 ] && [ ! -e "$psd" ]; } || rot; }

  for t in u-boot node ssh-keygen docker devcontainer ssh; do
    c "p15-voraussetzung-$t" && { kind "$t"; [ "$code" != 0 ] && grep -qF "$t fehlt im Pfad" "$w/err" && ! grep -q übersprungen "$w/err" || rot; }
  done
  c p15-rekursion && { (GP_KIND=1; eval "$(sed -n '/^kind() {/,/^}/p' "$gpk")"; kind nichts) > /dev/null 2> "$w/err"; [ $? = 1 ] && grep -q Rekursion "$w/err" || rot; }
  c p15-fremdes-verzeichnis && { mkdir "$w/tmp"; a=$(TMPDIR=$w/tmp mktemp -d -t pgr-cs-XXXXXX); TMPDIR=$w/tmp sonde ende; { [ -d "$a" ] && [ "$pcode" = 0 ]; } || rot; rmdir -- "$a" 2> /dev/null; }
  c p15-sonden-raeumen && { kind nichts '^p15-ende$'; [ "$code" = 0 ] || rot; }
  c p15-name-ohne-praefix && { (cd "$w" && env -i PATH="$PATH" HOME="$w" TMPDIR="$w" GP_KIND=1 GP_NAME=foo GP_NUR=p15-aufraeumen-ende GP_MUT='^$' GP_QUELLE="$quelle" "$BASH" "$gpk" > "$w/out" 2> "$w/err"); code=$?
    { [ "$code" != 0 ] && grep -qF 'trägt das Präfix pgr-cs- nicht' "$w/err"; } || rot; }
  c p15-port && { sonde port; a=$(sed -n 3p "$w/sonde"); { [[ $a =~ ^[0-9]+$ ]] && [ "$a" != 2222 ] && ! ( exec 3<> "/dev/tcp/127.0.0.1/$a" ) 2> /dev/null; } || rot; }
  c p15-port-2222 && { GP_PORT=2222 sonde port; [ -z "$(sed -n 3p "$w/sonde")" ] || rot; }
  c p15-port-fremd && { hoere; a=$lp; kill "$lk"; PATH="$w/bin:$PATH" AT_LOG=$w/log AT_STATE=$w/st GP_PORT=$a AT_DK_NAME="pgr-dev3 127.0.0.1:$a->2222/tcp" sonde port; [ -z "$(sed -n 3p "$w/sonde")" ] || rot; }
  for p in ende:0 INT:130 TERM:143; do
    c "p15-docker-${p%%:*}" && { case ${p%%:*} in ende) sonde docker-ende ;; *) sonde docker-warten "${p%%:*}" 7 ;; esac
      mapfile -t pl < "$w/sonde"; x=0; ! docker inspect "${pl[2]}" > /dev/null 2>&1 || x=1
      for v in "${pl[3]}" "${pl[@]:5}"; do ! docker volume inspect "$v" > /dev/null 2>&1 || x=1; done
      ! docker image inspect "${pl[4]}" > /dev/null 2>&1 || x=1
      { [ "$x" = 0 ] && [ "$pcode" = "${p#*:}" ]; } || { rot; dk_weg; }; }
  done

  # Punkt 16: Host-Ziel, kein Gate (Temp-Baum mit Makefile und Fragmenten)
  c p16-keine-kette && { mkdir -p "$w/t/harness/mk" "$w/t/tools/harness"; cp "$root/Makefile" "$w/t/"
    cp "$root"/harness/mk/{container-start,container-start-gegenprobe,gates-schnell}.mk "$w/t/harness/mk/"
    printf 'exit 1\n' > "$w/t/tools/harness/container-start-gegenprobe.sh"
    v=$(make -C "$w/t" -pq help 2> /dev/null || true)
    for t in GATE_CHECKS GATE_SCHNELL_AUSNAHME; do
      ! printf '%s\n' "$v" | awk -v z="$t" '$1 == z && ($2 == ":=" || $2 == "=") { for (i = 3; i <= NF; i++) if ($i == "container-start-gegenprobe") f = 1 } END { exit !f }' || rot "p16-$t"
    done
    ! printf '%s\n' "$v" | awk '$1 == "GATE_CHECKS" && ($2 == ":=" || $2 == "=") { for (i = 3; i <= NF; i++) if ($i == "container-start") f = 1 } END { exit !f }' || rot p16-ziel-kein-gate
    ! make -C "$w/t" container-start-gegenprobe > /dev/null 2>&1 || rot p16-scharf; }

  dk_gruppe
}

# --- Mutationen -----------------------------------------------------------------------------
MARKE='# --- Tabelle der Mut'"ationen ---"
baue() {
  mkdir -p "$1/tools/harness" "$1/harness/mk"; cp "$quelle/Makefile" "$1/"; cp "$quelle/tools/harness/container-start.sh" "$1/tools/harness/"
  cp "$quelle"/harness/mk/{container-start,container-start-gegenprobe,gates-schnell}.mk "$1/harness/mk/"; cp "$GP_SELF" "$1/gp.sh"
}
# mut <Name> <Datei> <alt> <neu> <SEL> [<roter Fall>…]
mut() {
  local name=$1 key=$2 alt=$3 repl=$4 d f c r got exp
  SEL=$5; shift 5
  [[ $name =~ ${GP_MUT:-.} ]] || return 0
  d=$arbeit/m-$name; baue "$d"
  f=$d/${FILE[$key]}; c=$(cat "$f"; echo x); c=${c%x}
  nach=; [ "$key" != p ] || { nach=$'\n'$MARKE${c#*$'\n'$MARKE}; c=${c%%$'\n'$MARKE*}; } # die Tabelle der Mutationen nennt jede Zeichenfolge selbst
  r=${c//"$alt"/}
  (( ${#c} - ${#r} == ${#alt} )) || { melde "Mutation $name: die Ersetzung trifft nicht genau eine Stelle"; exit 1; }
  printf '%s%s' "${c/"$alt"/"$repl"}" "$nach" > "$f"
  sk=$d/tools/harness/container-start.sh root=$d gpk=$d/gp.sh
  cp -a "$d" "$d.rein"; ROT=(); cases
  got=$(for k in "${!ROT[@]}"; do echo "$k"; done | sort | tr '\n' ' '); exp=$(for k; do echo "$k"; done | sort | tr '\n' ' ')
  [ -z "${GP_ENTDECKE-}" ] || echo "$name => $got" >&2
  [ "$got" = "$exp" ] || melde "Mutation $name: rot sind [$got], erwartet [$exp]"
}

# --- Tabelle der Mutationen ---
# Punkt 1
mut p1-zerstoerung s '  (( bad == 0 )) || exit 3' '  docker rm -f x; (( bad == 0 )) || exit 3' "p1-keine-zerstoerung|$P4S" p1-keine-zerstoerung p4-nichts-geaendert
mut p1-repo s '  aufgeloest=$(realpath -m -- "$PGR_DEV")' '  : > "$wurzel/x"; aufgeloest=$(realpath -m -- "$PGR_DEV")' 'p1-.*' p1-repo-unveraendert
mut p1-home s 'unset GIT_TOKEN' 'unset GIT_TOKEN; echo x >> "$HOME/.ssh/known_hosts"' "p1-home-unveraendert|$P4S" p1-home-unveraendert p4-nichts-geaendert
mut p1-ausserhalb s '  eltern=$(dirname -- "$PGR_DEV")' '  : > "$(dirname -- "$PGR_DEV")/x"; eltern=$(dirname -- "$PGR_DEV")' 'p1-nur-unter-pgr-dev' p1-nur-unter-pgr-dev
# Punkt 2
mut p2-dev-nicht-gesetzt s 'if [ -z "${PGR_DEV+x}" ]; then meld "PGR_DEV ist nicht gesetzt"' 'if false; then meld "PGR_DEV ist nicht gesetzt"' 'p2-dev-.*' p2-dev-nicht-gesetzt
mut p2-dev-leer s 'elif [ -z "$w" ]; then meld "PGR_DEV ist leer"' 'elif false; then meld "PGR_DEV ist leer"' 'p2-dev-.*' p2-dev-leer
mut p2-dev-relativ s 'elif [[ $w != /* ]]; then' 'elif false; then' 'p2-dev-.*' p2-dev-relativ
mut p2-dev-leerzeichen s '^/[A-Za-z0-9._+/-]*$' '^/[A-Za-z0-9._+/ -]*$' 'p2-dev-.*' p2-dev-leerzeichen
mut p2-dev-tabulator s '^/[A-Za-z0-9._+/-]*$' $'^/[A-Za-z0-9._+/\t-]*$' 'p2-dev-.*' p2-dev-tabulator
mut p2-dev-zeilenumbruch s '^/[A-Za-z0-9._+/-]*$' $'^/[A-Za-z0-9._+/\n-]*$' 'p2-dev-.*' p2-dev-zeilenumbruch
mut p2-dev-nicht-ascii s '^/[A-Za-z0-9._+/-]*$' '^/[A-Za-z0-9._+/ä-]*$' 'p2-dev-.*' p2-dev-nicht-ascii
mut p2-dev-slash-am-ende s 'elif [[ $w == */ ]]; then' 'elif false; then' 'p2-dev-.*' p2-dev-slash-am-ende
mut p2-dev-leeres-segment s 'elif [[ $w == *//* ]]; then' 'elif false; then' 'p2-dev-.*' p2-dev-leeres-segment
mut p2-dev-segment-punkt s 'elif [[ $w/ == */./* ]]; then' 'elif false; then' 'p2-dev-.*' p2-dev-segment-punkt
mut p2-dev-segment-punktpunkt s 'elif [[ $w/ == */../* ]]; then' 'elif false; then' 'p2-dev-.*' p2-dev-segment-punktpunkt
mut p2-dev-basis-gross s '^[a-z0-9][a-z0-9-]*$' '^[a-zA-Z0-9][a-zA-Z0-9-]*$' 'p2-dev-.*' p2-dev-basis-gross
mut p2-dev-basis-unterstrich s '^[a-z0-9][a-z0-9-]*$' '^[a-z0-9][a-z0-9_-]*$' 'p2-dev-.*' p2-dev-basis-unterstrich
mut p2-name-leer s 'PGR_NAME=${PGR_NAME-pgr-dev}' 'PGR_NAME=${PGR_NAME:-pgr-dev}' 'p2-name-.*' p2-name-leer
mut p2-name-gross s '^[a-z0-9][a-z0-9_.-]*$' '^[a-zA-Z0-9][a-zA-Z0-9_.-]*$' 'p2-name-.*' p2-name-gross
mut p2-port-leer s 'PGR_SSH_PORT=${PGR_SSH_PORT-2222}' 'PGR_SSH_PORT=${PGR_SSH_PORT:-2222}' 'p2-port-.*' p2-port-leer
mut p2-warten-leer s 'PGR_WARTEN=${PGR_WARTEN-60}' 'PGR_WARTEN=${PGR_WARTEN:-60}' 'p2-warten-.*' p2-warten-leer
mut p2-zahl-leer s '  [[ -n $1 ]] || return 1' '  :' 'p2-(port|warten)-.*' p2-port-leer p2-warten-leer
mut p2-zahl-ziffern s '  [[ $1 != *[!0-9+-]* ]] || return 1' '  :' 'p2-(port|warten)-.*' p2-port-abc p2-warten-ab
mut p2-zahl-plus s '  [[ $1 != +* ]] || return 1' '  :' 'p2-(port|warten)-.*' p2-port-plus22
mut p2-zahl-minus s '  [[ $1 != -* ]] || return 1' '  :' 'p2-(port|warten)-.*' p2-port-minus22
mut p2-zahl-plus-innen s '  [[ $1 != ?*+* ]] || return 1' '  :' 'p2-(port|warten)-.*' p2-port-2plus2
mut p2-zahl-minus-innen s '  [[ $1 != ?*-* ]] || return 1' '  :' 'p2-(port|warten)-.*' p2-port-22minus
mut p2-zahl-fuehrende-null s '  [[ $1 != 0?* ]] || return 1' '  :' 'p2-(port|warten)-.*' p2-port-02222
mut p2-zahl-null s '  [[ $1 != 0 ]] || return 1' '  :' 'p2-(port|warten)-.*' p2-port-0 p2-warten-0
mut p2-zahl-eins s '  [[ $1 != 0 ]] || return 1' '  [[ $1 != [01] ]] || return 1' 'p2-(port|warten|grenze)-.*' p2-grenze-PGR_SSH_PORT-1 p2-grenze-PGR_WARTEN-1
mut p2-zahl-obergrenze-streng s '[[ ! $1 > $2 ]]' '[[ $1 < $2 ]]' 'p2-(port|warten|grenze)-.*' p2-grenze-PGR_SSH_PORT-65535 p2-grenze-PGR_WARTEN-600
mut p2-zahl-oben s '(( ${#1} == ${#2} )) && [[ ! $1 > $2 ]]' 'true' 'p2-(port|warten|grenze)-.*' p2-port-65536 p2-warten-601
mut p2-vorgabe-name s 'PGR_NAME=${PGR_NAME-pgr-dev}' 'PGR_NAME=${PGR_NAME-pgr-x}' 'p2-vorgaben' p2-vorgaben
mut p2-vorgabe-port s 'PGR_SSH_PORT=${PGR_SSH_PORT-2222}' 'PGR_SSH_PORT=${PGR_SSH_PORT-2223}' 'p2-vorgaben' p2-vorgaben
mut p2-vorgabe-cpus s 'PGR_CPUS=${PGR_CPUS-8}' 'PGR_CPUS=${PGR_CPUS-9}' 'p2-vorgaben' p2-vorgaben
mut p2-vorgabe-speicher s 'PGR_MEM=${PGR_MEM-8g}' 'PGR_MEM=${PGR_MEM-12g}' 'p2-vorgaben' p2-vorgaben
mut p2-vorgabe-warten s 'PGR_WARTEN=${PGR_WARTEN-60}' 'PGR_WARTEN=${PGR_WARTEN-61}' 'p2-vorgaben'
mut p2-reihenfolge s 'for v in PGR_DEV PGR_NAME PGR_SSH_PORT PGR_WARTEN; do' 'for v in PGR_NAME PGR_DEV PGR_SSH_PORT PGR_WARTEN; do' 'p2-(alle-gemeldet|reihenfolge)' p2-reihenfolge
mut p2-nur-der-erste s 'meld() { fehl "$@"; bad=1; }' 'meld() { fehl "$@"; exit 2; }' 'p2-(alle-gemeldet|reihenfolge)' p2-alle-gemeldet
mut p2-exit s '(( bad == 0 )) || exit 2' '(( bad == 0 )) || exit 1' 'p2-dev-leer' p2-dev-leer
# Punkt 3
mut p3-ohne-docker s 'for t in docker u-boot devcontainer node ssh ssh-keygen; do' 'for t in u-boot devcontainer node ssh ssh-keygen; do' 'p3-(fehlt|alle).*' p3-alle-fehlen p3-fehlt-docker
mut p3-ohne-u-boot s 'for t in docker u-boot devcontainer node ssh ssh-keygen; do' 'for t in docker devcontainer node ssh ssh-keygen; do' 'p3-(fehlt|alle).*' p3-alle-fehlen p3-fehlt-u-boot
mut p3-ohne-devcontainer s 'for t in docker u-boot devcontainer node ssh ssh-keygen; do' 'for t in docker u-boot node ssh ssh-keygen; do' 'p3-(fehlt|alle).*' p3-alle-fehlen p3-fehlt-devcontainer
mut p3-ohne-node s 'for t in docker u-boot devcontainer node ssh ssh-keygen; do' 'for t in docker u-boot devcontainer ssh ssh-keygen; do' 'p3-(fehlt|alle).*' p3-alle-fehlen p3-fehlt-node
mut p3-ohne-ssh s 'for t in docker u-boot devcontainer node ssh ssh-keygen; do' 'for t in docker u-boot devcontainer node ssh-keygen; do' 'p3-(fehlt|alle).*' p3-alle-fehlen p3-fehlt-ssh
mut p3-ohne-ssh-keygen s 'for t in docker u-boot devcontainer node ssh ssh-keygen; do' 'for t in docker u-boot devcontainer node ssh; do' 'p3-(fehlt|alle).*' p3-alle-fehlen p3-fehlt-ssh-keygen
mut p3-reihenfolge s 'for t in docker u-boot devcontainer node ssh ssh-keygen; do' 'for t in docker u-boot devcontainer node ssh-keygen ssh; do' 'p3-(fehlt|alle).*' p3-alle-fehlen
mut p3-nur-das-erste s '{ meld "$t fehlt im Pfad"; continue; }' '{ meld "$t fehlt im Pfad"; break; }' 'p3-(fehlt|alle).*' p3-alle-fehlen
mut p3-version-trotz-fehlen s '{ meld "$t fehlt im Pfad"; continue; }' '{ meld "$t fehlt im Pfad"; }' 'p3-(fehlt|alle).*' p3-alle-fehlen p3-fehlt-docker p3-fehlt-u-boot
mut p3-u-boot-lexikalisch s '(( 10#$x > 0 )) || (( 10#$y >= 10 )) ||' '[[ ! $x.$y < 0.10 ]] ||' 'p3-u-boot-.*' p3-u-boot-0.9.0
mut p3-u-boot-groesser s '(( 10#$y >= 10 ))' '(( 10#$y > 10 ))' 'p3-u-boot-.*' p3-u-boot-0.10.0
mut p3-u-boot-hauptversion s '(( 10#$x > 0 )) || ' '' 'p3-u-boot-.*' p3-u-boot-1.0.0
mut p3-u-boot-nicht-geprueft s '(( 10#$x > 0 )) || (( 10#$y >= 10 )) || meld "u-boot ${BASH_REMATCH[0]} ist kleiner als 0.10.0"' ':' 'p3-u-boot-.*' p3-u-boot-0.9.0
mut p3-u-boot-unlesbar s '          meld "u-boot: Version unlesbar"' '          :' 'p3-u-boot-.*' p3-u-boot-exit1 p3-u-boot-unlesbar
mut p3-u-boot-exit s '|| v=' '|| :' 'p3-u-boot-.*' p3-u-boot-exit1
mut p3-docker-exit s '|| sv=' '|| :' 'p3-docker-.*' p3-docker-exit1
mut p3-docker-unlesbar s '            meld "docker: Version des Servers unlesbar"' '            :' 'p3-docker-.*' p3-docker-exit1 p3-docker-unlesbar
mut p3-docker-info s 'if docker info > /dev/null 2>&1; then' 'if true; then' 'p3-docker-.*' p3-docker-daemon
mut p3-docker-version s '(( 10#${BASH_REMATCH[1]} >= 24 ))' 'true' 'p3-docker-.*' p3-docker-23
mut p3-docker-groesser s '(( 10#${BASH_REMATCH[1]} >= 24 ))' '(( 10#${BASH_REMATCH[1]} > 24 ))' 'p3-docker-.*' p3-docker-24
mut p3-exit s '(( bad == 0 )) || exit 3' '(( bad == 0 )) || exit 1' 'p3-fehlt-ssh' p3-fehlt-ssh
# Punkt 4
mut p4-meldung-form s 'echo "container-start: $schritt: $*" >&2' 'echo "$schritt: $*" >&2' "$P4S|p3-.*" p3-alle-fehlen p3-docker-23 p3-docker-daemon p3-docker-exit1 p3-docker-unlesbar p3-fehlt-devcontainer p3-fehlt-docker p3-fehlt-node p3-fehlt-ssh p3-fehlt-ssh-keygen p3-fehlt-u-boot p3-u-boot-0.9.0 p3-u-boot-exit1 p3-u-boot-unlesbar p4-meldung-form
mut p4-meldung-stdout s 'echo "container-start: $schritt: $*" >&2' 'echo "container-start: $schritt: $*"' "$P4S" p2-dev-leer p2-port-abc p3-fehlt-ssh p3-u-boot-0.9.0 p4-nichts-geaendert p5-elternpfad-fehlt p5-fremd-datei p6-nur-privat p9-kollision p6-passphrase
mut p4-schreibt-vor-pruefung s '  eltern=$(dirname -- "$PGR_DEV")' '  eltern=$(dirname -- "$PGR_DEV"); echo x > "$eltern/spur"' "$P4S" p4-meldung-form p4-nichts-geaendert
mut p4-ruft-up s '  (( bad == 0 )) || exit 3' '  docker up; (( bad == 0 )) || exit 3' "$P4S" p4-nichts-geaendert
# Punkt 5
mut p5-nicht-aufgeloest s 'aufgeloest=$(realpath -m -- "$PGR_DEV")' 'aufgeloest=$PGR_DEV' 'p5-.*' p5-im-repo-symlink
mut p5-praefix s '|| $aufgeloest == "$wurzel"/* ]]' '|| $aufgeloest == "$wurzel"* ]]' 'p5-.*' p5-praefix-ohne-slash
mut p5-wurzel s '$aufgeloest == "$wurzel" || ' '' 'p5-.*' p5-wurzel-selbst
mut p5-arbeitsverzeichnis s 'wurzel=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)' 'wurzel=$(pwd -P)' 'p5-.*' p5-elternpfad-fehlt p5-elternpfad-ist-datei p5-fehlt p5-fortsetzen p5-fremd-anderer-name p5-fremd-datei p5-im-repo p5-im-repo-symlink p5-kein-verzeichnis p5-leer p5-vorhanden-ohne-aufruf p5-vorhanden-unveraendert p5-wurzel-selbst
mut p5-mkdir-p s '  [ -d "$eltern" ] || ende 4 "der Elternpfad $eltern fehlt oder ist kein Verzeichnis"' '  mkdir -p -- "$eltern"' 'p5-.*' p5-elternpfad-fehlt p5-elternpfad-ist-datei
mut p5-eltern-datei s '[ -d "$eltern" ]' '[ -e "$eltern" ]' 'p5-.*' p5-elternpfad-ist-datei
mut p5-kein-verzeichnis s '    [ -d "$PGR_DEV" ] || ende 4 "$PGR_DEV ist kein Verzeichnis"' '    :' 'p5-.*' p5-kein-verzeichnis
mut p5-leer s 'if [ -n "$(ls -A -- "$PGR_DEV")" ]; then' 'if true; then' 'p5-.*' p5-leer
mut p5-vorhanden s '[ "$zv" != vorhanden ] || return 0' ':' 'p5-.*' p5-vorhanden-ohne-aufruf
mut p5-devcontainer-json s '[ -f "$PGR_DEV/.devcontainer/devcontainer.json" ] && zv=vorhanden' ':' 'p5-.*' p5-vorhanden-ohne-aufruf
mut p5-fortsetzen-mit-init s '[ "$zv" = neu ] && ub init' 'ub init' 'p5-.*|p7-fortsetzen-ohne-init|p12-wiederholung' p12-wiederholung p7-fortsetzen-ohne-init
mut p5-fremd s 'if grep -qxF "    name: $name" "$PGR_DEV/u-boot.yaml" 2> /dev/null; then' 'if true; then' 'p5-.*|p13-.*' p13-fremd-vor-schluessel p5-fremd-anderer-name p5-fremd-datei
mut p5-fremd-vergleich s 'grep -qxF' 'grep -qF' 'p5-.*' p5-fremd-anderer-name
mut p5-nicht-anlegen s '[ -d "$PGR_DEV" ] || mkdir -- "$PGR_DEV" || ende 5 "$PGR_DEV ließ sich nicht anlegen"' ':' 'p5-fehlt' p5-fehlt
mut p5-exit s 'ende 4 "$PGR_DEV ist fremd (kein u-boot.yaml mit dem Namen $name), nichts wird überschrieben"' 'ende 1 "$PGR_DEV ist fremd (kein u-boot.yaml mit dem Namen $name), nichts wird überschrieben"' 'p5-fremd-datei' p5-fremd-datei
# Punkt 6
mut p6-ueberschreibt s '  [ "$zs" = erzeugen ] || return 0' ':' 'p6-.*' p6-gueltig-unveraendert p6-kommentar-ignoriert
mut p6-modus-ordner s '&& chmod 700 -- "$sd"' '' 'p6-.*' p6-gueltig-unveraendert p6-kommentar-ignoriert p6-modus-ordner p6-modus-ordner-vorhanden
mut p6-ssh-datei s 'if { [ -e "$sd" ] || [ -L "$sd" ]; } && [ ! -d "$sd" ]; then' 'if false; then' 'p6-.*' p6-ssh-ist-datei p6-ssh-ist-symlink
mut p6-modus-privat s '|| ende 5 "ssh-keygen endete mit Exit $?"' '|| ende 5 "ssh-keygen endete mit Exit $?"; chmod 644 "$priv"' 'p6-.*' p6-gueltig-unveraendert p6-kommentar-ignoriert p6-modus-privat p6-ohne-passphrase p6-pub-anderer
mut p6-kommentar s '-C "$name"' '-C x' 'p6-.*' p6-kommentar
mut p6-passphrase-anlegen s $'-N \'\'' '-N x' 'p6-.*' p6-gueltig-unveraendert p6-kommentar-ignoriert p6-ohne-passphrase p6-pub-anderer
mut p6-typ s '-t ed25519' '-t rsa' 'p6-.*' p6-typ
mut p6-a-ordner s '(( 8#$(stat -c %a -- "$sd") & 077 )) &&' 'false &&' 'p6-.*' p6-ordner-750
mut p6-a-privat s '(( 8#$(stat -c %a -- "$priv") & 077 )) &&' 'false &&' 'p6-.*' p6-privat-640
mut p6-b s $'if pk=$(ssh-keygen -y -P \'\' -f "$priv" 2> /dev/null < /dev/null); then' $'pk=$(ssh-keygen -y -P \'\' -f "$priv" 2> /dev/null < /dev/null); if true; then' 'p6-.*' p6-leere-datei p6-passphrase
mut p6-c s '[ "$t1 $k1" = "$t2 $k2" ] ||' 'true ||' 'p6-.*' p6-pub-anderer
mut p6-c-ganze-zeile s '[ "$t1 $k1" = "$t2 $k2" ]' '[ "$pk" = "$(<"$pub")" ]' 'p6-.*' p6-kommentar-ignoriert
mut p6-eine-datei s 'elif (( hp + hb == 1 )); then' 'elif false; then' 'p6-.*' p6-nur-oeffentlich p6-nur-privat
mut p6-exit s 'ende 4 "Schlüsselpaar abgelehnt, nichts wird verändert: ${mangel[*]}"' 'ende 1 "Schlüsselpaar abgelehnt, nichts wird verändert: ${mangel[*]}"' 'p6-passphrase' p6-passphrase
# Punkt 7
mut p7-init-devcontainer s 'ub init "$name" --devcontainer --sandbox --no-git' 'ub init "$name" --sandbox --no-git' 'p7-.*' p7-aufruf-1 p7-echt-u-boot
mut p7-init-sandbox s 'ub init "$name" --devcontainer --sandbox --no-git' 'ub init "$name" --devcontainer --no-git' 'p7-.*' p7-aufruf-1 p7-echt-u-boot
mut p7-init-no-git s 'ub init "$name" --devcontainer --sandbox --no-git' 'ub init "$name" --devcontainer --sandbox' 'p7-.*' p7-aufruf-1 p7-echt-u-boot
mut p7-init-name s 'ub init "$name" --devcontainer --sandbox --no-git' 'ub init --devcontainer --sandbox --no-git' 'p7-.*' p7-aufruf-1
mut p7-node s 'ub config set devcontainer.features.node.enabled true' 'ub config set devcontainer.features.node.enabled false' 'p7-.*' p7-aufruf-2 p7-echt-u-boot
mut p7-sshd s 'ub config set devcontainer.features.sshd.enabled true' 'ub config set devcontainer.features.sshd.enabled false' 'p7-.*' p7-aufruf-3 p7-echt-u-boot
mut p7-claude s 'ub config set devcontainer.features.claude-code.enabled true' 'ub config set devcontainer.features.claude-code.enabled false' 'p7-.*' p7-aufruf-4 p7-echt-u-boot
mut p7-runtime s 'ub config set devcontainer.sandbox.nestedRuntime docker' 'ub config set devcontainer.sandbox.nestedRuntime podman' 'p7-.*' p7-aufruf-5 p7-echt-u-boot
mut p7-egress s 'ub config set devcontainer.sandbox.egress.enabled true' 'ub config set devcontainer.sandbox.egress.enabled false' 'p7-.*' p7-aufruf-6
mut p7-allow-gcr s 'gcr.io,pkg-containers.githubusercontent.com' 'pkg-containers.githubusercontent.com' 'p7-.*' p7-aufruf-7 p7-echt-u-boot
mut p7-allow-pkg s 'gcr.io,pkg-containers.githubusercontent.com' 'gcr.io' 'p7-.*' p7-aufruf-7 p7-echt-u-boot
mut p7-generate s 'ub generate devcontainer' 'ub generate' 'p7-.*' p7-aufruf-8 p7-echt-u-boot
mut p7-up s 'ub generate devcontainer' 'ub generate devcontainer; ub up' 'p7-(kein-up-doctor|reihenfolge)' p7-kein-up-doctor p7-reihenfolge
mut p7-doctor s 'ub generate devcontainer' 'ub generate devcontainer; ub doctor' 'p7-(kein-up-doctor|reihenfolge)' p7-kein-up-doctor p7-reihenfolge
mut p7-reihenfolge s '  ub config set devcontainer.features.node.enabled true' '  ub generate devcontainer' 'p7-.*' p7-aufruf-2 p7-echt-u-boot p7-reihenfolge
mut p7-fehler-uebersehen s '[ "$kcode" = 0 ] || ende 5 "u-boot $* endete mit Exit $kcode"' ':' 'p7-.*' p7-fehler-1 p7-fehler-2 p7-fehler-3 p7-fehler-4 p7-fehler-5 p7-fehler-6 p7-fehler-7 p7-fehler-8
mut p7-fehler-exit s 'ende 5 "u-boot $* endete mit Exit $kcode"' 'ende 1 "u-boot $* endete mit Exit $kcode"' 'p7-.*' p7-fehler-1 p7-fehler-2 p7-fehler-3 p7-fehler-4 p7-fehler-5 p7-fehler-6 p7-fehler-7 p7-fehler-8
mut p7-warnung-als-fehler s 'endete mit Exit $kcode"; }' 'endete mit Exit $kcode"; [ -z "${AT_UB_WARN-}" ] || ende 5 "Warnung"; }' 'p7-.*' p7-warnung
mut p7-stdout-verschluckt s 'ub() { kind u-boot "$@";' 'ub() { kind u-boot "$@" > /dev/null;' 'p7-.*' p7-ausgabe-stdout
mut p7-stderr-verschluckt s 'ub() { kind u-boot "$@";' 'ub() { kind u-boot "$@" 2> /dev/null;' 'p7-.*' p7-ausgabe-stderr p7-fehler-1 p7-fehler-2 p7-fehler-3 p7-fehler-4 p7-fehler-5 p7-fehler-6 p7-fehler-7 p7-fehler-8 p7-warnung
mut p7-kein-cd s '  cd -- "$PGR_DEV" || ende 5 "$PGR_DEV ist als Arbeitsverzeichnis des Gerüsts nicht erreichbar"' '  :' 'p7-im-verzeichnis' p7-im-verzeichnis
# Punkt 8
mut p8-name s '"--name=" + (process.env.PGR_NAME || "pgr-dev"),' '' 'p8-.*' p8-name
mut p8-name-wert s 'process.env.PGR_NAME || "pgr-dev"' '"pgr-dev"' 'p8-.*' p8-name
mut p8-cpus s '"--cpus=" + (process.env.PGR_CPUS || "8"), ' '' 'p8-.*' p8-cpus
mut p8-speicher s ', "--memory=" + (process.env.PGR_MEM || "8g")' '' 'p8-.*' p8-speicher
mut p8-adresse s '"127.0.0.1:" + port' '"0.0.0.0:" + port' 'p8-.*' p8-port
mut p8-port s '"127.0.0.1:" + port + ":2222"' '"127.0.0.1:" + "2222" + ":2222"' 'p8-.*' p8-port
mut p8-readonly s ',readonly"' '"' 'p8-.*' p8-mount
mut p8-mount-quelle s 'source=${localWorkspaceFolder}/ssh/id_ed25519.pub' 'source=${localWorkspaceFolder}/ssh/id_ed25519' 'p8-.*' p8-mount
mut p8-mount-ziel s 'target=/home/vscode/.ssh/authorized_keys' 'target=/home/vscode/.ssh/x' 'p8-.*' p8-mount
mut p8-mount-bind s 'type=bind,' 'type=volume,' 'p8-.*' p8-mount
mut p8-verweis-wert s '"${localEnv:CLAUDE_CODE_OAUTH_TOKEN}"' 'process.env.CLAUDE_CODE_OAUTH_TOKEN' 'p8-.*' p8-env
mut p8-verweis-fehlt s 'j.remoteEnv.CLAUDE_CODE_OAUTH_TOKEN = "${localEnv:CLAUDE_CODE_OAUTH_TOKEN}";' '' 'p8-.*' p8-env
mut p8-rest s 'const port = process.env.PGR_SSH_PORT || "2222";' 'const port = process.env.PGR_SSH_PORT || "2222"; j.image = "x";' 'p8-.*' p8-rest
mut p8-ueberschreiben s 'if [ ! -e mk-override.sh ]; then' 'if true; then' 'p8-.*' p8-skript-env-PGR_CPUS p8-skript-env-PGR_MEM p8-skript-env-PGR_NAME p8-skript-env-PGR_SSH_PORT p8-skript-fehler p8-umgebung-ANTHROPIC_API_KEY p8-umgebung-CLAUDE_CODE_OAUTH_TOKEN p8-umgebung-GIT_TOKEN p8-vorhanden-bleibt p8-vorhanden-laeuft
mut p8-modus s '    chmod 755 mk-override.sh' '    :' 'p8-.*' p8-skript-modus
mut p8-export-name s 'export PGR_NAME PGR_SSH_PORT PGR_CPUS PGR_MEM' 'export PGR_SSH_PORT PGR_CPUS PGR_MEM' 'p8-.*' p8-skript-env-PGR_NAME
mut p8-export-port s 'export PGR_NAME PGR_SSH_PORT PGR_CPUS PGR_MEM' 'export PGR_NAME PGR_CPUS PGR_MEM' 'p8-.*' p8-skript-env-PGR_SSH_PORT
mut p8-export-cpus s 'export PGR_NAME PGR_SSH_PORT PGR_CPUS PGR_MEM' 'export PGR_NAME PGR_SSH_PORT PGR_MEM' 'p8-.*' p8-skript-env-PGR_CPUS
mut p8-export-speicher s 'export PGR_NAME PGR_SSH_PORT PGR_CPUS PGR_MEM' 'export PGR_NAME PGR_SSH_PORT PGR_CPUS' 'p8-.*' p8-skript-env-PGR_MEM
mut p8-skript-laeuft-nicht s '  bash ./mk-override.sh || ende 5 "mk-override.sh endete mit Exit $?"' '  :' 'p8-vorhanden-laeuft' p8-vorhanden-laeuft
mut p8-skript-fehler s '|| ende 5 "mk-override.sh endete mit Exit $?"' '|| true' 'p8-.*' p8-skript-fehler
mut p8-behalte-oauth s 'unset CLAUDE_CODE_OAUTH_TOKEN' ':' 'p8-.*|p9-umgebung' p8-umgebung-CLAUDE_CODE_OAUTH_TOKEN p9-umgebung
mut p8-behalte-apikey s 'unset ANTHROPIC_API_KEY' ':' 'p8-.*|p9-umgebung' p8-umgebung-ANTHROPIC_API_KEY p9-umgebung
mut p8-behalte-gittoken s 'unset GIT_TOKEN' ':' 'p8-.*|p9-umgebung' p8-umgebung-GIT_TOKEN p9-umgebung
mut p8-ausgabe-apikey s 'unset ANTHROPIC_API_KEY' 'echo "k=${ANTHROPIC_API_KEY-}"' 'p8-.*' p8-keine-werte p8-umgebung-ANTHROPIC_API_KEY
mut p8-ausgabe-oauth s 'unset CLAUDE_CODE_OAUTH_TOKEN' 'echo "t=${CLAUDE_CODE_OAUTH_TOKEN-}"' 'p8-.*' p8-keine-werte p8-umgebung-CLAUDE_CODE_OAUTH_TOKEN
# Punkt 13
mut p13-eingabe-nach-werkzeug s 'for schritt in eingaben werkzeuge zustand geruest schluessel override start warten; do' 'for schritt in werkzeuge eingaben zustand geruest schluessel override start warten; do' 'p13-.*' p13-eingabe-vor-werkzeug
mut p13-werkzeug-nach-zustand s 'for schritt in eingaben werkzeuge zustand geruest schluessel override start warten; do' 'for schritt in eingaben zustand werkzeuge geruest schluessel override start warten; do' 'p13-.*' p13-werkzeug-vor-repo
mut p13-verzeichnis-vor-repo s 's_zustand() { zust_repo; zust_verzeichnis; zust_schluessel; zust_container; zust_port; }' 's_zustand() { zust_verzeichnis; zust_repo; zust_schluessel; zust_container; zust_port; }' 'p13-.*' p13-repo-vor-fremd
mut p13-schluessel-vor-verzeichnis s 's_zustand() { zust_repo; zust_verzeichnis; zust_schluessel; zust_container; zust_port; }' 's_zustand() { zust_repo; zust_schluessel; zust_verzeichnis; zust_container; zust_port; }' 'p13-.*' p13-fremd-vor-schluessel
# Punkt 14
mut p14-plus m $'\t@bash tools/harness/container-start.sh' $'\t+@bash tools/harness/container-start.sh' 'p14-.*' p14-makeflags-n p14-makeflags-n-i p14-makeflags-q p14-makeflags-t p14-n p14-n-i p14-q p14-rezept p14-t
mut p14-make m $'\t@bash tools/harness/container-start.sh' $'\t@bash tools/harness/container-start.sh # $(MAKE)' 'p14-.*' p14-makeflags-n p14-makeflags-n-i p14-makeflags-q p14-makeflags-t p14-n p14-n-i p14-q p14-rezept p14-t
mut p14-phony m '.PHONY: container-start' '.PHONY: nichts' 'p14-.*' p14-makeflags-t p14-phony p14-t
# Punkt 15
mut p15-ende p $'[ "$GP_PROBE" = ende ] && exit 0' $'[ "$GP_PROBE" = ende ] && kill -KILL $$' 'p15-(praefix|aufraeumen-.*)' p15-aufraeumen-ende
mut p15-exit-trap p 'trap aufraeumen EXIT' ':' 'p15-(praefix|aufraeumen-.*|docker-.*)' p15-aufraeumen-ende p15-aufraeumen-int p15-aufraeumen-term p15-docker-INT p15-docker-TERM p15-docker-ende
mut p15-int p $'trap \'exit 130\' INT' ':' 'p15-(praefix|aufraeumen-.*|docker-.*)'
mut p15-term p $'trap \'exit 143\' TERM' ':' 'p15-(praefix|aufraeumen-.*|docker-.*)'
mut p15-praefix p 'mktemp -d -t pgr-cs-XXXXXX' 'mktemp -d "${GP_PROBE_DATEI%/*}/xyz-XXXXXX"' 'p15-(praefix|aufraeumen-.*)' p15-aufraeumen-ende p15-aufraeumen-int p15-aufraeumen-term p15-praefix
mut p15-schutz p '[[ ${arbeit##*/} == pgr-cs-* ]] && ' '' 'p15-(praefix|aufraeumen-.*)'
mut p15-rekursion-schutz p '  [ -z "${GP_KIND-}" ] || { melde "Rekursion: eine Kind-Gegenprobe ruft kind erneut auf"; exit 1; }' '' 'p15-rekursion' p15-rekursion
mut p15-sonden-per-praefix p 'local d; for d in "${PSDS[@]}"; do' 'local d; for d in "${TMPDIR:-/tmp}"/pgr-cs-*; do' 'p15-(fremdes-verzeichnis|aufraeumen-.*)' p15-fremdes-verzeichnis
mut p15-sonden-bleiben p 'rmdir -- "$d" 2> /dev/null' ':' 'p15-sonden-raeumen' p15-sonden-raeumen
mut p15-voraussetzung-ohne-node p 'for t in make node ssh-keygen u-boot docker devcontainer ssh; do' 'for t in make ssh-keygen u-boot docker devcontainer ssh; do' 'p15-voraussetzung-.*' p15-voraussetzung-node
mut p15-voraussetzung-ohne-ssh-keygen p 'for t in make node ssh-keygen u-boot docker devcontainer ssh; do' 'for t in make node u-boot docker devcontainer ssh; do' 'p15-voraussetzung-.*' p15-voraussetzung-ssh-keygen
mut p15-voraussetzung-ohne-u-boot p 'for t in make node ssh-keygen u-boot docker devcontainer ssh; do' 'for t in make node ssh-keygen docker devcontainer ssh; do' 'p15-voraussetzung-.*' p15-voraussetzung-u-boot
mut p15-voraussetzung-ohne-docker p 'for t in make node ssh-keygen u-boot docker devcontainer ssh; do' 'for t in make node ssh-keygen u-boot devcontainer ssh; do' 'p15-voraussetzung-.*' p15-voraussetzung-docker
mut p15-voraussetzung-ohne-devcontainer p 'for t in make node ssh-keygen u-boot docker devcontainer ssh; do' 'for t in make node ssh-keygen u-boot docker ssh; do' 'p15-voraussetzung-.*' p15-voraussetzung-devcontainer
mut p15-voraussetzung-ohne-ssh p 'for t in make node ssh-keygen u-boot docker devcontainer ssh; do' 'for t in make node ssh-keygen u-boot docker devcontainer; do' 'p15-voraussetzung-.*' p15-voraussetzung-ssh
# Punkt 16
mut p16-gate-checks g '.PHONY: container-start-gegenprobe' $'GATE_CHECKS += container-start-gegenprobe\n.PHONY: container-start-gegenprobe' 'p16-.*' p16-GATE_CHECKS
mut p16-scharf g $'\t@bash tools/harness/container-start-gegenprobe.sh' $'\t-@bash tools/harness/container-start-gegenprobe.sh' 'p16-.*' p16-scharf
mut p16-ausnahme k 'GATE_SCHNELL_AUSNAHME := test-container ' 'GATE_SCHNELL_AUSNAHME := container-start-gegenprobe test-container ' 'p16-.*' p16-GATE_SCHNELL_AUSNAHME
mut p16-ziel-gate m '.PHONY: container-start' $'GATE_CHECKS += container-start\n.PHONY: container-start' 'p16-.*' p16-ziel-kein-gate
# Punkt 9 bis 12 (Gruppe Attrappen)
SC='p9-(kollision|label-anderer-name|adresse-weicht-ab|port-weicht-ab|laeuft|exited|created|kein-update)|p10-belegt.*|p13-(schluessel-vor-kollision|kollision-vor-port)|p1-keine-zerstoerung-kollision'
SU='p9-(up-.*|keiner|exited|created|ausgabe-.*|umgebung)'
SW='p11-(bereit-meldung|erster-versuch|nach-drei|nie|opt-.*|known-hosts-unveraendert)'
mut p9-kollision s 'if [ -n "$nm" ] && [[ $'"'"'\n'"'"'$lab$'"'"'\n'"'"' != *$'"'"'\n'"'"'$PGR_NAME$'"'"'\n'"'"'* ]]; then' 'if false; then' "$SC" p1-keine-zerstoerung-kollision p13-kollision-vor-port p9-kollision
mut p9-kollision-zerstoert s 'ende 4 "der Container $nm trägt' 'docker rm -f "$nm"; ende 4 "der Container $nm trägt' "p1-keine-zerstoerung-kollision|$P4S" p1-keine-zerstoerung-kollision p4-nichts-geaendert
mut p9-anderer-name s 'ende 4 "der Container des Gerüsts heißt $n, nicht $PGR_NAME"' ':' "$SC" p9-label-anderer-name
mut p9-port-nicht-geprueft s '[ "$bind" = "127.0.0.1:$PGR_SSH_PORT " ] || ende 4' ':; true || ende 4' "$SC" p9-adresse-weicht-ab p9-port-weicht-ab
mut p9-port-ohne-adresse s '[ "$bind" = "127.0.0.1:$PGR_SSH_PORT " ]' '[[ $bind == *":$PGR_SSH_PORT " ]]' "$SC" p9-adresse-weicht-ab
mut p9-adresse-ohne-port s '[ "$bind" = "127.0.0.1:$PGR_SSH_PORT " ]' '[[ $bind == "127.0.0.1:"* ]]' "$SC" p9-port-weicht-ab
mut p9-laeuft-immer s '[ "${i%% *}" != running ] || zc=laeuft' 'zc=laeuft' "$SC" p10-belegt-gestoppt p9-created p9-exited
mut p9-created-laeuft s '[ "${i%% *}" != running ] || zc=laeuft' '[ "${i%% *}" = exited ] || zc=laeuft' "$SC" p9-created
mut p9-exited-laeuft s '[ "${i%% *}" != running ] || zc=laeuft' '[ "${i%% *}" = created ] || zc=laeuft' "$SC" p10-belegt-gestoppt p9-exited
mut p9-up-bei-laeuft s $'  [ "$zc" != laeuft ] || return 0\n  kind devcontainer up' '  kind devcontainer up' "$SC" p9-laeuft
mut p9-override-bei-laeuft s $'s_override() {\n  [ "$zc" != laeuft ] || return 0\n' $'s_override() {\n' "$SC" p9-laeuft
mut p9-ohne-override s $'s_override() {\n  [ "$zc" != laeuft ] || return 0\n' $'s_override() {\n  return 0\n' "$SC" p9-created p9-exited
mut p9-kein-up-ohne-container s $'  [ "$zc" != laeuft ] || return 0\n  kind devcontainer up' $'  [ "$zc" = gestoppt ] || return 0\n  kind devcontainer up' "$SU" p9-keiner
mut p9-update s '[ "${i%% *}" != running ] || zc=laeuft' '[ "${i%% *}" != running ] || { zc=laeuft; docker update --memory "$PGR_MEM" --memory-swap "$PGR_MEM" "$PGR_NAME"; }' "$SC" p9-kein-update
mut p9-outcome s $'  [[ $kausg == *\'"outcome":"success"\'* ]] || ende 5 \'devcontainer up meldete kein "outcome":"success"\'' ':' "$SU" p9-up-ohne-outcome
mut p9-up-exit s '  [ "$kcode" = 0 ] || ende 5 "devcontainer up endete mit Exit $kcode"' ':' "$SU" p9-up-exit1
mut p9-stdout-verschluckt s 'kind devcontainer up --workspace-folder . --override-config .devcontainer/override.json' 'kind devcontainer up --workspace-folder . --override-config .devcontainer/override.json > /dev/null' "$SU" p9-ausgabe-stdout
mut p9-stderr-verschluckt s 'kind devcontainer up --workspace-folder . --override-config .devcontainer/override.json' 'kind devcontainer up --workspace-folder . --override-config .devcontainer/override.json 2> /dev/null' "$SU" p9-ausgabe-stderr
mut p9-up-unterbefehl s 'kind devcontainer up ' 'kind devcontainer start ' "$SU" p9-created p9-exited p9-keiner p9-up-aufruf
mut p9-up-workspace s '--workspace-folder .' '--workspace-folder ..' "$SU" p9-up-aufruf
mut p9-up-override s '--override-config .devcontainer/override.json' '--config .devcontainer/override.json' "$SU" p9-up-aufruf
mut p9-up-override-datei s '--override-config .devcontainer/override.json' '--override-config .devcontainer/devcontainer.json' "$SU" p9-up-aufruf
mut p9-up-verzeichnis s '  kind devcontainer up ' '  cd /; kind devcontainer up ' "$SU" p9-up-verzeichnis
mut p4-stdout s '  kind devcontainer up ' '  echo start; kind devcontainer up ' 'p4-stdout-nur-bereit' p4-stdout-nur-bereit
# Punkt 10
mut p10-ohne-pruefung s '! ( exec 3<> "/dev/tcp/127.0.0.1/$PGR_SSH_PORT" ) 2> /dev/null || ende 4 "der Port $PGR_SSH_PORT ist belegt"' ':' "$SC" p10-belegt p10-belegt-gestoppt
mut p10-bei-laeuft s $'  [ "$zc" != laeuft ] || return 0\n  ! ( exec' $'  ! ( exec' "$SC" p9-laeuft
mut p10-nur-ohne-container s $'  [ "$zc" != laeuft ] || return 0\n  ! ( exec' $'  [ "$zc" = keiner ] || return 0\n  ! ( exec' "$SC" p10-belegt-gestoppt
# Punkt 11
mut p11-schluessel s 'kind -m ssh -i "$PGR_DEV/ssh/id_ed25519"' 'kind -m ssh -i "$PGR_DEV/ssh/x"' "$SW" p11-opt-schluessel
mut p11-identitiesonly s ' -o IdentitiesOnly=yes' '' "$SW" p11-opt-identitiesonly
mut p11-batchmode s ' -o BatchMode=yes' '' "$SW" p11-opt-batchmode
mut p11-connecttimeout s ' -o ConnectTimeout=5' ' -o ConnectTimeout=50' "$SW" p11-opt-connecttimeout
mut p11-stricthostkeychecking s ' -o StrictHostKeyChecking=no' '' "$SW" p11-opt-stricthostkeychecking
mut p11-userknownhostsfile s ' -o UserKnownHostsFile=/dev/null' '' "$SW" p11-opt-userknownhostsfile p11-known-hosts-unveraendert
mut p11-known-hosts-geschrieben s '  local t0=$SECONDS' '  local t0=$SECONDS; echo x >> "$HOME/.ssh/known_hosts"' "$SW" p11-known-hosts-unveraendert
mut p11-loglevel s ' -o LogLevel=ERROR' '' "$SW" p11-opt-loglevel
mut p11-port s '-p "$PGR_SSH_PORT"' '-p 2222' "$SW" p11-opt-port
mut p11-benutzer s 'vscode@' 'root@' "$SW" p11-opt-benutzer-adresse
mut p11-adresse s '@127.0.0.1 true' '@localhost true' "$SW" p11-opt-benutzer-adresse
mut p11-befehl s '@127.0.0.1 true' '@127.0.0.1 false' "$SW" p11-opt-befehl
mut p11-kein-sleep s '    sleep 1' '    :' "$SW" p11-nie
mut p11-zu-kurz s 'SECONDS - t0 < PGR_WARTEN ))' 'SECONDS - t0 < PGR_WARTEN - 1 ))' "$SW" p11-nie
mut p11-ausgabe-summiert s $'  kausg=\n  if (( m ))' '  if (( m ))' "$SW" p11-nie
mut p11-erfolg-exit s 'if [ "$kcode" = 0 ]; then echo' 'if true; then echo' "$SW" p11-nach-drei p11-nie
mut p11-meldung-name s 'bereit $PGR_NAME 127.0.0.1:' 'bereit pgr-dev 127.0.0.1:' "$SW" p11-bereit-meldung
mut p11-meldung-port s '127.0.0.1:$PGR_SSH_PORT"; return 0' '127.0.0.1:2222"; return 0' "$SW" p11-bereit-meldung
mut p11-meldung-stderr s 'echo "container-start: bereit $PGR_NAME 127.0.0.1:$PGR_SSH_PORT";' 'echo "container-start: bereit $PGR_NAME 127.0.0.1:$PGR_SSH_PORT" >&2;' 'p11-bereit-meldung' p11-bereit-meldung
mut p11-vorgabe-warten s 'PGR_WARTEN=${PGR_WARTEN-60}' 'PGR_WARTEN=${PGR_WARTEN-61}' 'p11-vorgabe-warten' p11-vorgabe-warten
# Punkt 12
SS='p12-.*'
mut p12-kind-geruest s '[ -z "$kpid" ] || kill "$kpid"' '[ -z "$kpid" ] || [ "$schritt" = geruest ] || kill "$kpid"' "$SS" p12-sigINT-geruest p12-sigTERM-geruest
mut p12-kind-start s '[ -z "$kpid" ] || kill "$kpid"' '[ -z "$kpid" ] || [ "$schritt" = start ] || kill "$kpid"' "$SS" p12-sigINT-start p12-sigTERM-start
mut p12-kind-warten s '[ -z "$kpid" ] || kill "$kpid"' '[ -z "$kpid" ] || [ "$schritt" = warten ] || kill "$kpid"' "$SS" p12-sigINT-warten p12-sigTERM-warten
mut p12-int-entfaellt s $'trap \'abbruch 130\' INT' ':' "$SS" p12-sigINT-geruest p12-sigINT-start p12-sigINT-warten
mut p12-term-entfaellt s $'trap \'abbruch 143\' TERM' ':' "$SS" p12-sigTERM-geruest p12-sigTERM-start p12-sigTERM-warten
mut p12-exit-int s $'trap \'abbruch 130\' INT' $'trap \'abbruch 1\' INT' "$SS" p12-sigINT-geruest p12-sigINT-start p12-sigINT-warten
mut p12-exit-term s $'trap \'abbruch 143\' TERM' $'trap \'abbruch 1\' TERM' "$SS" p12-sigTERM-geruest p12-sigTERM-start p12-sigTERM-warten
mut p12-ohne-meldung s '; ende "$1" abgebrochen' '; exit "$1"' "$SS" p12-sigINT-geruest p12-sigINT-start p12-sigINT-warten p12-sigTERM-geruest p12-sigTERM-start p12-sigTERM-warten
mut p12-rueckbau-docker s '; ende "$1" abgebrochen' '; docker rm -f "$PGR_NAME" > /dev/null 2>&1; ende "$1" abgebrochen' "$SS" p12-sigINT-geruest p12-sigINT-start p12-sigINT-warten p12-sigTERM-geruest p12-sigTERM-start p12-sigTERM-warten
mut p12-verzeichnis-veraendert s '; ende "$1" abgebrochen' '; : > "$PGR_DEV/spur"; ende "$1" abgebrochen' "$SS" p12-sigINT-start p12-sigTERM-start
# Punkt 13 (mit Container)
mut p13-schluessel-nach-container s 'zust_schluessel; zust_container; zust_port; }' 'zust_container; zust_schluessel; zust_port; }' 'p13-.*' p13-schluessel-vor-kollision
mut p13-container-nach-port s 'zust_container; zust_port; }' 'zc=keiner; zust_port; zust_container; }' 'p13-.*' p13-kollision-vor-port
# Punkt 15: Name, Port, Aufräumen der Container, Volumes und Images
mut p15-name p '[[ $nm == pgr-cs-* ]] || {' 'true || {' 'p15-name-ohne-praefix' p15-name-ohne-praefix
mut p15-raeumen-name p '[[ $n == pgr-cs-* ]] || continue' ':' 'p15-name-ohne-praefix'
mut p15-port-2222 p '[ "$p" != 2222 ] && ! docker' '! docker' 'p15-port.*' p15-port-2222
mut p15-port-fremd p 'grep -qF ":$p->"' 'grep -qF "nie"' 'p15-port.*' p15-port-fremd
mut p15-rm-ohne-v p 'docker rm -fv "$n"' 'docker rm -f "$n"' 'p15-docker-.*' p15-docker-INT p15-docker-TERM p15-docker-ende
mut p15-container-bleibt p 'docker rm -fv "$n" > /dev/null 2>&1' ':' 'p15-docker-.*' p15-docker-INT p15-docker-TERM p15-docker-ende
mut p15-volume-bleibt p 'if [[ $v == pgr-cs-* ]]; then docker volume rm "$v" > /dev/null 2>&1' 'if [[ $v == pgr-cs-* ]]; then :' 'p15-docker-.*' p15-docker-INT p15-docker-TERM p15-docker-ende
mut p15-volume-rest-entfernen p 'rest=1; docker volume rm "$v" > /dev/null 2>&1' 'rest=1' 'p15-docker-.*'
mut p15-image-bleibt p 'docker rmi "$i" > /dev/null 2>&1' ':' 'p15-docker-.*' p15-docker-INT p15-docker-TERM p15-docker-ende
# Gruppe Docker (echtes Docker)
mut p9-docker-neue-volumes s 'j.mounts = [...j.mounts, "type=bind' 'j.mounts = j.mounts.map((m) => m.replace("-workspace-", "-ws" + Date.now() + "-")); j.mounts = [...j.mounts, "type=bind' 'p9-docker-.*' p9-docker-rm-neustart
mut p9-docker-neustart-ersetzt s '  kind devcontainer up ' '  [ "$zc" != gestoppt ] || docker rm -f "$PGR_NAME"; kind devcontainer up ' "p9-(exited|created)|p9-docker-stop-start" p9-created p9-docker-stop-start p9-exited
mut p10-docker-ohne-pruefung s '! ( exec 3<> "/dev/tcp/127.0.0.1/$PGR_SSH_PORT" ) 2> /dev/null || ende 4 "der Port $PGR_SSH_PORT ist belegt"' ':' 'p10-docker-port-belegt' p10-docker-port-belegt
mut p10-docker-bei-laeuft s $'  [ "$zc" != laeuft ] || return 0\n  ! ( exec' $'  ! ( exec' 'p9-laeuft|p9-docker-zweiter-aufruf' p9-docker-zweiter-aufruf p9-laeuft
mut p11-docker-schluessel s 'kind -m ssh -i "$PGR_DEV/ssh/id_ed25519"' 'kind -m ssh -i "$PGR_DEV/ssh/x"' 'p11-docker-bereit' p11-docker-bereit
mut p11-docker-erfolg s 'if [ "$kcode" = 0 ]; then echo' 'if true; then echo' 'p11-docker-anderer-schluessel' p11-docker-anderer-schluessel
mut p15-docker-cpus s '(process.env.PGR_CPUS || "8")' '"8"' 'p8-cpus|p15-docker-grenzen' p8-cpus p15-docker-grenzen
mut p15-docker-speicher s '(process.env.PGR_MEM || "8g")' '"8g"' 'p8-speicher|p15-docker-grenzen' p8-speicher p15-docker-grenzen
baue "$arbeit/orig"
sk=$arbeit/orig/tools/harness/container-start.sh root=$arbeit/orig gpk=$GP_SELF SEL=${GP_NUR:-'.*'}
cp -a "$root" "$root.rein"; ROT=(); cases
for k in "${!ROT[@]}"; do melde "Fall $k ist an der unveränderten Kopie rot"; done
sonden_raeumen
for d in "${PSDS[@]}"; do [ ! -e "$d" ] || melde "Verzeichnis der Sonde bleibt: $d"; done
[ "$fehler" = 0 ] || exit 1
echo "container-start-gegenprobe: gruen — SPEC-058: Fälle der Punkte 1 bis 16 und je Zusage eine Mutation"
