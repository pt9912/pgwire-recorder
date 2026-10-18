#!/usr/bin/env bash
# container-start-gegenprobe — prüft tools/harness/container-start.sh und harness/mk/container-start.mk
# gegen Teil I von SPEC-058 (harness/vertraege/SPEC-058-container-start.md: Punkte 1 bis 8, die Teile der
# Punkte 13 und 14 ohne Container, Punkte 15 und 16). Ziel für den Host, kein Gate (ADR-0047). Gruppe Attrappen: Temp-Verzeichnisse mit Attrappen
# von docker, u-boot, devcontainer und ssh in einem vorangestellten Pfad, echtes node und ssh-keygen,
# ohne Docker und ohne Netz; dazu der Fall p7-echt-u-boot gegen das echte u-boot. Jede Attrappe
# protokolliert Aufruf und Argumente. Der Name jedes Falls trägt die Nummer seines Punktes; die
# unveränderte Kopie besteht jeden Fall, und bei jeder Mutation (genau eine Textersetzung in der Kopie
# einer Datei; trifft sie nichts oder mehr als eine Stelle, bricht die Gegenprobe ab) werden von den
# Fällen ihrer Gruppe (SEL) genau die benannten rot.
#
# Voraussetzungen (SPEC-058 Punkt 15): make, GNU-Werkzeuge, node, ssh-keygen und u-boot im Pfad; die Gegenprobe prüft sie
# zuerst und endet bei einem fehlenden mit einer Meldung (Fälle p15-voraussetzung-*).
# Verzeichnisse tragen das Präfix pgr-cs- (mktemp -d); entfernt wird nur dieses Verzeichnis, bei Ende, SIGINT und
# SIGTERM (Fälle p15-*). GRENZE: Die Mutanten p15-ende und p15-exit-trap hinterlassen leere Verzeichnisse der Sonde;
# die Gegenprobe entfernt am Ende nur diese, nur wenn leer, nur per rmdir, und prüft, dass keines bleibt.
# GRENZE: p15-aufraeumen-int und -term setzen voraus, dass die Gegenprobe SIGINT nicht ignoriert (kein Start
# als Hintergrundjob einer Shell ohne Jobsteuerung). Es gibt keinen Container und keinen Fall zu Punkt 9
# bis 12 (Teil II).
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
PSDS=()
sonden_raeumen() { local d; for d in "${PSDS[@]}"; do [[ ${d##*/} == pgr-cs-* ]] && [ -d "$d" ] && rmdir -- "$d" 2> /dev/null; done; return 0; }
aufraeumen() { [ -z "$bg" ] || kill "$bg" 2> /dev/null; sonden_raeumen; [[ ${arbeit##*/} == pgr-cs-* ]] && [ -d "$arbeit" ] && rm -rf -- "$arbeit"; }
trap aufraeumen EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
if [ -n "${GP_PROBE-}" ]; then # Sondenmodus: macht das Aufräumen von außen beobachtbar
  printf '%s\n%s\n' "$arbeit" "$$" > "$GP_PROBE_DATEI"
  [ "$GP_PROBE" = ende ] && exit 0
  sleep 60 & bg=$!
  wait "$bg"
  exit 0
fi

fehler=0
melde() { echo "container-start-gegenprobe: ROT — $*" >&2; fehler=1; }
declare -A ROT FILE=([s]=tools/harness/container-start.sh [m]=harness/mk/container-start.mk [g]=harness/mk/container-start-gegenprobe.mk [k]=harness/mk/gates-schnell.mk [p]=gp.sh)
rot() { ROT[${1:-$fn}]=1; [ -z "${GP_DEBUG-}" ] || echo "ROT ${1:-$fn} $w" >&2; }

fehlt=0
for t in make node ssh-keygen u-boot; do command -v "$t" > /dev/null 2>&1 || { melde "$t fehlt im Pfad (Voraussetzung der Gegenprobe)"; fehlt=1; }; done
[ "$fehlt" = 0 ] || exit 1
node=$(command -v node) sshkeygen=$(command -v ssh-keygen) realub=$(command -v u-boot)
sys=$arbeit/sys att=$arbeit/att
mkdir -p "$sys" "$att"
for t in bash dirname realpath ls grep stat mkdir cat chmod make env; do ln -s "$(command -v "$t")" "$sys/$t"; done
cat > "$att/attrappe" <<'ATTRAPPE'
#!/usr/bin/env bash
n=${0##*/}
echo "$n $*" >> "$AT_LOG"
for v in CLAUDE_CODE_OAUTH_TOKEN ANTHROPIC_API_KEY GIT_TOKEN; do [ -z "${!v+x}" ] || echo "ENV! $n $v" >> "$AT_LOG"; done
case $n in
  docker) case $1 in info) exit "${AT_DK_INFO-0}" ;; version) echo "${AT_DK_VER-29.9.1}"; exit "${AT_DK_VEXIT-0}" ;; esac ;;
  u-boot)
    [ "$1" = --version ] && { echo "${AT_UB_AUS-u-boot version ${AT_UB_VER-0.10.0}}"; exit "${AT_UB_EXIT-0}"; }
    mkdir -p "$AT_STATE"
    k=$(( $(cat "$AT_STATE/k" 2> /dev/null || echo 0) + 1 )); echo $k > "$AT_STATE/k"
    echo "UB-STDOUT-$k"; echo "UB-STDERR-$k" >&2
    [ "${AT_UB_WARN-}" = 1 ] && echo "Warning: Lockerung" >&2
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
P4S='p2-dev-leer|p2-port-abc|p3-fehlt-ssh|p3-u-boot-0.9.0|p5-fremd-datei|p5-elternpfad-fehlt|p6-nur-privat|p6-passphrase|p4-nichts-geaendert|p4-meldung-form'
ZW='container-start: override: Teil II fehlt, Gerüst liegt, kein Container gestartet'
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
gruen() { lauf "$@"; ex 6 "$ZW"; }
# ex <Exit-Code> [<Text>…]: Exit-Code und Texte auf stderr; bei Exit 2 bis 4 gelten Punkt 4 (nichts verändert,
# nur lesende Aufrufe, stdout leer) und die Form der Meldung.
ex() {
  local x=$1 t; shift
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
sonde() { # <Modus> [<Signal>]: startet die Gegenprobe im Sondenmodus; setzt pcode und psd
  local d=$w/sonde k
  if [ -n "${2-}" ]; then
    ( for k in $(seq 100); do [ -s "$d" ] && [ "$(wc -l < "$d")" -ge 2 ] && break; sleep 0.1; done; kill -"$2" "$(sed -n 2p "$d")" ) &
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

cases() {
  local t k v a p
  # Punkt 1: Gegenstand und Ort
  c p1-keine-zerstoerung && { gruen; ! grep -qE '^docker (stop|rm|rmi|kill|volume|network|image|container|system)|prune' "$w/log" || rot; }
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
    for a in --name=pgr-dev --cpus=8 --memory=12g 127.0.0.1:2222:2222; do grep -qF -- "$a" "$v" || rot; done; }
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
  for v in PGR_NAME=pgr-dev PGR_SSH_PORT=2222 PGR_CPUS=8 PGR_MEM=12g; do
    c "p8-skript-env-${v%%=*}" && { gruen; printf '#!/bin/bash\nenv > "$(dirname "$0")/envdump"\n' > "$dev/mk-override.sh"
      gruen; grep -qx "$v" "$dev/envdump" 2> /dev/null || rot; }
  done
  for v in CLAUDE_CODE_OAUTH_TOKEN ANTHROPIC_API_KEY GIT_TOKEN; do
    c "p8-umgebung-$v" && { gruen; printf '#!/bin/bash\nenv > "$(dirname "$0")/envdump"\n' > "$dev/mk-override.sh"
      gruen CLAUDE_CODE_OAUTH_TOKEN=SENTINEL-1 ANTHROPIC_API_KEY=SENTINEL-2 GIT_TOKEN=SENTINEL-3
      { ! grep -q "^ENV! .* $v\$" "$w/log" && [ -s "$dev/envdump" ] && ! grep -q "^$v=" "$dev/envdump"; } || rot; }
  done
  c p8-keine-werte && { gruen CLAUDE_CODE_OAUTH_TOKEN=SENTINEL-1 ANTHROPIC_API_KEY=SENTINEL-2 GIT_TOKEN=SENTINEL-3; ! grep -rq SENTINEL "$w" || rot; }

  # Punkt 13: Reihenfolge der Gründe, Zwischenstand
  c p13-zwischenstand && { gruen; ! grep -q bereit "$w/out" "$w/err" && ! grep -qv '^UB-STDOUT-' "$w/out" || rot; }
  c p13-eingabe-vor-werkzeug && { rm "$w/bin/u-boot"; lauf PGR_NAME=X; ex 2 'eingaben: PGR_NAME '; ! grep -q werkzeuge "$w/err" || rot; }
  c p13-werkzeug-vor-repo && { rm "$w/bin/docker"; lauf PGR_DEV="$root/pgr-cs-dev"; ex 3 'docker fehlt'; }
  c p13-repo-vor-fremd && { mkdir "$root/pgr-cs-xdir"; echo x > "$root/pgr-cs-xdir/datei"; lauf PGR_DEV="$root/pgr-cs-xdir"; ex 4 'liegt im Repo'; ! grep -q 'ist fremd' "$w/err" || rot; }
  c p13-fremd-vor-schluessel && { mkdir -p "$w/f/ssh"; echo x > "$w/f/datei"; "$sshkeygen" -q -t ed25519 -N '' -f "$w/f/ssh/id_ed25519" < /dev/null; rm "$w/f/ssh/id_ed25519.pub"
    lauf PGR_DEV="$w/f"; ex 4 'ist fremd'; ! grep -q 'nur eine' "$w/err" || rot; }

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
  c p14-i && { mt; mk -i PGR_NAME=X; [ "$code" = 0 ] && grep -qF 'container-start: eingaben: PGR_NAME ' "$w/err" && ! grep -q bereit "$w/out" "$w/err" || rot; }
  c p14-makeflags-i && { mt; mk MAKEFLAGS=-i PGR_NAME=X; [ "$code" = 0 ] && grep -qF 'container-start: eingaben: PGR_NAME ' "$w/err" && ! grep -q bereit "$w/out" "$w/err" || rot; }
  for p in k:-k j:-j4 makeflags-k:MAKEFLAGS=-k makeflags-j:MAKEFLAGS=-j4; do
    c "p14-${p%%:*}" && { mt; mk "${p#*:}"; [ "$code" != 0 ] && grep -qF "$ZW" "$w/err" && ! grep -q bereit "$w/out" "$w/err" || rot; }
  done

  # Punkt 15: Rahmen — Aufräumen und Präfix
  c p15-praefix && { sonde ende; [[ ${psd##*/} == pgr-cs-* ]] || rot; }
  c p15-aufraeumen-ende && { sonde ende; { [ "$pcode" = 0 ] && [ ! -e "$psd" ]; } || rot; }
  c p15-aufraeumen-int && { sonde warten INT; { [ "$pcode" = 130 ] && [ ! -e "$psd" ]; } || rot; }
  c p15-aufraeumen-term && { sonde warten TERM; { [ "$pcode" = 143 ] && [ ! -e "$psd" ]; } || rot; }

  for t in u-boot node ssh-keygen; do
    c "p15-voraussetzung-$t" && { kind "$t"; [ "$code" != 0 ] && grep -qF "$t fehlt im Pfad" "$w/err" && ! grep -q übersprungen "$w/err" || rot; }
  done
  c p15-rekursion && { (GP_KIND=1; eval "$(sed -n '/^kind() {/,/^}/p' "$gpk")"; kind nichts) > /dev/null 2> "$w/err"; [ $? = 1 ] && grep -q Rekursion "$w/err" || rot; }
  c p15-sonden-raeumen && { kind nichts '^p15-ende$'; [ "$code" = 0 ] || rot; }

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
mut p2-vorgabe-speicher s 'PGR_MEM=${PGR_MEM-12g}' 'PGR_MEM=${PGR_MEM-11g}' 'p2-vorgaben' p2-vorgaben
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
mut p4-meldung-form s 'echo "container-start: $schritt: $*" >&2' 'echo "$schritt: $*" >&2' "$P4S|p13-zwischenstand|p3-.*" p13-zwischenstand p3-alle-fehlen p3-docker-23 p3-docker-24 p3-docker-daemon p3-docker-exit1 p3-docker-unlesbar p3-fehlt-devcontainer p3-fehlt-docker p3-fehlt-node p3-fehlt-ssh p3-fehlt-ssh-keygen p3-fehlt-u-boot p3-u-boot-0.10.0 p3-u-boot-0.100.0 p3-u-boot-0.9.0 p3-u-boot-1.0.0 p3-u-boot-exit1 p3-u-boot-unlesbar p4-meldung-form p6-nur-privat p6-passphrase
mut p4-meldung-stdout s 'echo "container-start: $schritt: $*" >&2' 'echo "container-start: $schritt: $*"' "$P4S|p13-zwischenstand" p13-zwischenstand p2-dev-leer p2-port-abc p3-fehlt-ssh p3-u-boot-0.9.0 p4-nichts-geaendert p5-elternpfad-fehlt p5-fremd-datei p6-nur-privat p6-passphrase
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
mut p5-fortsetzen-mit-init s '[ "$zv" = neu ] && ub init' 'ub init' 'p5-.*|p7-fortsetzen-ohne-init' p7-fortsetzen-ohne-init
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
mut p7-fehler-uebersehen s '|| { c=$?; ende 5 "u-boot $* endete mit Exit $c"; }; }' '|| true; }' 'p7-.*' p7-fehler-1 p7-fehler-2 p7-fehler-3 p7-fehler-4 p7-fehler-5 p7-fehler-6 p7-fehler-7 p7-fehler-8
mut p7-fehler-exit s 'ende 5 "u-boot $* endete mit Exit $c"' 'ende 1 "u-boot $* endete mit Exit $c"' 'p7-.*' p7-fehler-1 p7-fehler-2 p7-fehler-3 p7-fehler-4 p7-fehler-5 p7-fehler-6 p7-fehler-7 p7-fehler-8
mut p7-warnung-als-fehler s 'ende 5 "u-boot $* endete mit Exit $c"; }; }' 'ende 5 "u-boot $* endete mit Exit $c"; }; [ -z "${AT_UB_WARN-}" ] || ende 5 "Warnung"; }' 'p7-.*' p7-warnung
mut p7-stdout-verschluckt s 'ub() { u-boot "$@" ||' 'ub() { u-boot "$@" > /dev/null ||' 'p7-.*' p7-ausgabe-stdout
mut p7-stderr-verschluckt s 'ub() { u-boot "$@" ||' 'ub() { u-boot "$@" 2> /dev/null ||' 'p7-.*' p7-ausgabe-stderr p7-fehler-1 p7-fehler-2 p7-fehler-3 p7-fehler-4 p7-fehler-5 p7-fehler-6 p7-fehler-7 p7-fehler-8 p7-warnung
mut p7-kein-cd s '  cd -- "$PGR_DEV" || ende 5 "$PGR_DEV ist als Arbeitsverzeichnis des Gerüsts nicht erreichbar"' '  :' 'p7-im-verzeichnis' p7-im-verzeichnis
# Punkt 8
mut p8-name s '"--name=" + (process.env.PGR_NAME || "pgr-dev"),' '' 'p8-.*' p8-name
mut p8-name-wert s 'process.env.PGR_NAME || "pgr-dev"' '"pgr-dev"' 'p8-.*' p8-name
mut p8-cpus s '"--cpus=" + (process.env.PGR_CPUS || "8"), ' '' 'p8-.*' p8-cpus
mut p8-speicher s ', "--memory=" + (process.env.PGR_MEM || "12g")' '' 'p8-.*' p8-speicher
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
mut p8-behalte-oauth s 'unset CLAUDE_CODE_OAUTH_TOKEN' ':' 'p8-.*' p8-umgebung-CLAUDE_CODE_OAUTH_TOKEN
mut p8-behalte-apikey s 'unset ANTHROPIC_API_KEY' ':' 'p8-.*' p8-umgebung-ANTHROPIC_API_KEY
mut p8-behalte-gittoken s 'unset GIT_TOKEN' ':' 'p8-.*' p8-umgebung-GIT_TOKEN
mut p8-ausgabe-apikey s 'unset ANTHROPIC_API_KEY' 'echo "k=${ANTHROPIC_API_KEY-}"' 'p8-.*' p8-keine-werte p8-umgebung-ANTHROPIC_API_KEY
mut p8-ausgabe-oauth s 'unset CLAUDE_CODE_OAUTH_TOKEN' 'echo "t=${CLAUDE_CODE_OAUTH_TOKEN-}"' 'p8-.*' p8-keine-werte p8-umgebung-CLAUDE_CODE_OAUTH_TOKEN
# Punkt 13
mut p13-eingabe-nach-werkzeug s 'for schritt in eingaben werkzeuge zustand geruest schluessel override; do' 'for schritt in werkzeuge eingaben zustand geruest schluessel override; do' 'p13-.*' p13-eingabe-vor-werkzeug
mut p13-werkzeug-nach-zustand s 'for schritt in eingaben werkzeuge zustand geruest schluessel override; do' 'for schritt in eingaben zustand werkzeuge geruest schluessel override; do' 'p13-.*' p13-werkzeug-vor-repo
mut p13-verzeichnis-vor-repo s 's_zustand() { zust_repo; zust_verzeichnis; zust_schluessel; }' 's_zustand() { zust_verzeichnis; zust_repo; zust_schluessel; }' 'p13-.*' p13-repo-vor-fremd
mut p13-schluessel-vor-verzeichnis s 's_zustand() { zust_repo; zust_verzeichnis; zust_schluessel; }' 's_zustand() { zust_repo; zust_schluessel; zust_verzeichnis; }' 'p13-.*' p13-fremd-vor-schluessel
mut p13-exit-null s 'ende 6 "Teil II fehlt, Gerüst liegt, kein Container gestartet"' 'ende 0 "Teil II fehlt, Gerüst liegt, kein Container gestartet"' 'p13-zwischenstand|p14-(k|j|makeflags-k|makeflags-j)' p13-zwischenstand p14-j p14-k p14-makeflags-j p14-makeflags-k
mut p13-bereit s 'ende 6 "Teil II fehlt, Gerüst liegt, kein Container gestartet"' 'ende 6 "bereit"' 'p13-zwischenstand|p14-(k|j|makeflags-k|makeflags-j)' p13-zwischenstand p14-j p14-k p14-makeflags-j p14-makeflags-k
# Punkt 14
mut p14-plus m $'\t@bash tools/harness/container-start.sh' $'\t+@bash tools/harness/container-start.sh' 'p14-.*' p14-makeflags-n p14-makeflags-n-i p14-makeflags-q p14-makeflags-t p14-n p14-n-i p14-q p14-rezept p14-t
mut p14-make m $'\t@bash tools/harness/container-start.sh' $'\t@bash tools/harness/container-start.sh # $(MAKE)' 'p14-.*' p14-makeflags-n p14-makeflags-n-i p14-makeflags-q p14-makeflags-t p14-n p14-n-i p14-q p14-rezept p14-t
mut p14-phony m '.PHONY: container-start' '.PHONY: nichts' 'p14-.*' p14-makeflags-t p14-phony p14-t
# Punkt 15
mut p15-ende p $'[ "$GP_PROBE" = ende ] && exit 0' $'[ "$GP_PROBE" = ende ] && kill -KILL $$' 'p15-(praefix|aufraeumen-.*)' p15-aufraeumen-ende
mut p15-exit-trap p 'trap aufraeumen EXIT' ':' 'p15-(praefix|aufraeumen-.*)' p15-aufraeumen-ende p15-aufraeumen-int p15-aufraeumen-term
mut p15-int p $'trap \'exit 130\' INT' ':' 'p15-(praefix|aufraeumen-.*)'
mut p15-term p $'trap \'exit 143\' TERM' ':' 'p15-(praefix|aufraeumen-.*)'
mut p15-praefix p 'mktemp -d -t pgr-cs-XXXXXX' 'mktemp -d "${GP_PROBE_DATEI%/*}/xyz-XXXXXX"' 'p15-(praefix|aufraeumen-.*)' p15-aufraeumen-ende p15-aufraeumen-int p15-aufraeumen-term p15-praefix
mut p15-schutz p '[[ ${arbeit##*/} == pgr-cs-* ]] && ' '' 'p15-(praefix|aufraeumen-.*)'
mut p15-rekursion-schutz p '  [ -z "${GP_KIND-}" ] || { melde "Rekursion: eine Kind-Gegenprobe ruft kind erneut auf"; exit 1; }' '' 'p15-rekursion' p15-rekursion
mut p15-sonden-bleiben p 'rmdir -- "$d" 2> /dev/null' ':' 'p15-sonden-raeumen' p15-sonden-raeumen
mut p15-voraussetzung-ohne-node p 'for t in make node ssh-keygen u-boot; do' 'for t in make ssh-keygen u-boot; do' 'p15-voraussetzung-.*' p15-voraussetzung-node
mut p15-voraussetzung-ohne-ssh-keygen p 'for t in make node ssh-keygen u-boot; do' 'for t in make node u-boot; do' 'p15-voraussetzung-.*' p15-voraussetzung-ssh-keygen
mut p15-voraussetzung-ohne-u-boot p 'for t in make node ssh-keygen u-boot; do' 'for t in make node ssh-keygen; do' 'p15-voraussetzung-.*' p15-voraussetzung-u-boot
# Punkt 16
mut p16-gate-checks g '.PHONY: container-start-gegenprobe' $'GATE_CHECKS += container-start-gegenprobe\n.PHONY: container-start-gegenprobe' 'p16-.*' p16-GATE_CHECKS
mut p16-scharf g $'\t@bash tools/harness/container-start-gegenprobe.sh' $'\t-@bash tools/harness/container-start-gegenprobe.sh' 'p16-.*' p16-scharf
mut p16-ausnahme k 'GATE_SCHNELL_AUSNAHME := test-container ' 'GATE_SCHNELL_AUSNAHME := container-start-gegenprobe test-container ' 'p16-.*' p16-GATE_SCHNELL_AUSNAHME
mut p16-ziel-gate m '.PHONY: container-start' $'GATE_CHECKS += container-start\n.PHONY: container-start' 'p16-.*' p16-ziel-kein-gate
baue "$arbeit/orig"
sk=$arbeit/orig/tools/harness/container-start.sh root=$arbeit/orig gpk=$GP_SELF SEL=${GP_NUR:-'.*'}
cp -a "$root" "$root.rein"; ROT=(); cases
for k in "${!ROT[@]}"; do melde "Fall $k ist an der unveränderten Kopie rot"; done
sonden_raeumen
for d in "${PSDS[@]}"; do [ ! -e "$d" ] || melde "Verzeichnis der Sonde bleibt: $d"; done
[ "$fehler" = 0 ] || exit 1
echo "container-start-gegenprobe: gruen — Teil I von SPEC-058: Fälle der Punkte 1 bis 8, 13 bis 16 und je Zusage eine Mutation"
