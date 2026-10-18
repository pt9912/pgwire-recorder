#!/usr/bin/env bash
# container-start — legt am Host das Gerüst des Dev-Containers an und startet den Container (Schritte
# eingaben, werkzeuge, zustand, geruest, schluessel, override, start, warten); endet mit der Zeile `bereit`.
# Vertrag: harness/vertraege/SPEC-058-container-start.md (Punkte 1 bis 14). Kein Gate. Eingaben aus der
# Umgebung: PGR_DEV (Pflicht), PGR_NAME, PGR_SSH_PORT, PGR_WARTEN, PGR_CPUS, PGR_MEM.
# Exit-Codes: 0 bereit · 2 eingaben · 3 werkzeuge · 4 zustand · 5 Fehler eines Schrittes · 130 SIGINT · 143 SIGTERM.
#
# GEPRÜFT DURCH tools/harness/container-start-gegenprobe.sh (Fälle p<Punkt>-*).
set -uo pipefail
export LC_ALL=C
# Punkt 8: Kein Kindprozess (u-boot, mk-override.sh, docker, devcontainer, ssh) erbt einen Zugang.
unset CLAUDE_CODE_OAUTH_TOKEN
unset ANTHROPIC_API_KEY
unset GIT_TOKEN

schritt=eingaben
fehl() { echo "container-start: $schritt: $*" >&2; }
ende() { local c=$1; shift; fehl "$@"; exit "$c"; }
bad=0
meld() { fehl "$@"; bad=1; }

# Punkt 12: SIGINT und SIGTERM beenden das laufende Kind und den Lauf; nichts wird zurückgebaut.
kpid=
abbruch() { [ -z "$kpid" ] || kill "$kpid" 2> /dev/null; ende "$1" abgebrochen; }
trap 'abbruch 130' INT
trap 'abbruch 143' TERM
# kind [-m] <Befehl…>: Kindprozess im Hintergrund, damit ein Signal sofort wirkt. Seine Standardausgabe geht
# unverändert nach stdout und steht in $kausg (mit -m samt Standardfehler nur in $kausg); $kcode ist sein Exit-Code.
kind() {
  local m=0 z r fd; [ "${1-}" != -m ] || { m=1; shift; }
  kausg=
  if (( m )); then exec {fd}< <(exec "$@" 2>&1 < /dev/null); else exec {fd}< <(exec "$@" < /dev/null); fi
  kpid=$!
  while :; do
    IFS= read -r z <&"$fd"; r=$?
    (( r == 0 )) && z+=$'\n'
    kausg+=$z; (( m )) || printf '%s' "$z"
    (( r == 0 )) || break
  done
  exec {fd}<&-
  wait "$kpid"; kcode=$?; kpid=
}

# zahl <Wert> <Höchstwert>: Dezimalziffern, kein Vorzeichen, keine führende Null, 1 bis Höchstwert.
zahl() {
  [[ -n $1 ]] || return 1
  [[ $1 != *[!0-9+-]* ]] || return 1
  [[ $1 != +* ]] || return 1
  [[ $1 != -* ]] || return 1
  [[ $1 != ?*+* ]] || return 1
  [[ $1 != ?*-* ]] || return 1
  [[ $1 != 0?* ]] || return 1
  [[ $1 != 0 ]] || return 1
  (( ${#1} < ${#2} )) || { (( ${#1} == ${#2} )) && [[ ! $1 > $2 ]]; }
}

# Punkt 2: eingaben. Nicht gesetzt ergibt die Vorgabe, gesetzt und leer ist ungültig.
PGR_NAME=${PGR_NAME-pgr-dev}
PGR_SSH_PORT=${PGR_SSH_PORT-2222}
PGR_WARTEN=${PGR_WARTEN-60}
PGR_CPUS=${PGR_CPUS-8}
PGR_MEM=${PGR_MEM-8g}
pruefe_PGR_DEV() {
  local w=${PGR_DEV-} q zeichen='^/[A-Za-z0-9._+/-]*$'
  q=$(printf '%q' "$w")
  if [ -z "${PGR_DEV+x}" ]; then meld "PGR_DEV ist nicht gesetzt"
  elif [ -z "$w" ]; then meld "PGR_DEV ist leer"
  elif [[ $w != /* ]]; then meld "PGR_DEV $q ist nicht absolut"
  elif [[ ! $w =~ $zeichen ]]; then meld "PGR_DEV $q enthält ein Zeichen außerhalb von A-Za-z0-9._+-"
  elif [[ $w == */ ]]; then meld "PGR_DEV $q endet auf /"
  elif [[ $w == *//* ]]; then meld "PGR_DEV $q hat ein leeres Segment"
  elif [[ $w/ == */./* ]]; then meld "PGR_DEV $q hat das Segment ."
  elif [[ $w/ == */../* ]]; then meld "PGR_DEV $q hat das Segment .."
  elif [[ ! ${w##*/} =~ ^[a-z0-9][a-z0-9-]*$ ]]; then meld "PGR_DEV $q hat einen Basisnamen außerhalb von a-z0-9-"
  fi
}
pruefe_PGR_NAME() {
  [[ $PGR_NAME =~ ^[a-z0-9][a-z0-9_.-]*$ ]] || meld "PGR_NAME $(printf '%q' "$PGR_NAME") ist ungültig"
}
pruefe_PGR_SSH_PORT() {
  zahl "$PGR_SSH_PORT" 65535 || meld "PGR_SSH_PORT $(printf '%q' "$PGR_SSH_PORT") ist keine ganze Zahl von 1 bis 65535"
}
pruefe_PGR_WARTEN() {
  zahl "$PGR_WARTEN" 600 || meld "PGR_WARTEN $(printf '%q' "$PGR_WARTEN") ist keine ganze Zahl von 1 bis 600"
}
s_eingaben() {
  for v in PGR_DEV PGR_NAME PGR_SSH_PORT PGR_WARTEN; do pruefe_$v; done
  (( bad == 0 )) || exit 2
  name=${PGR_DEV##*/}
}

# Punkt 3: werkzeuge. Alle Verstöße des Schrittes werden gemeldet, in dieser Reihenfolge.
s_werkzeuge() {
  for t in docker u-boot devcontainer node ssh ssh-keygen; do
    command -v "$t" > /dev/null 2>&1 || { meld "$t fehlt im Pfad"; continue; }
    case $t in
      u-boot)
        v=$(u-boot --version 2> /dev/null) || v=
        if [[ $v =~ ([0-9]+)\.([0-9]+)\.([0-9]+) ]]; then
          x=${BASH_REMATCH[1]} y=${BASH_REMATCH[2]}
          (( 10#$x > 0 )) || (( 10#$y >= 10 )) || meld "u-boot ${BASH_REMATCH[0]} ist kleiner als 0.10.0"
        else
          meld "u-boot: Version unlesbar"
        fi ;;
      docker)
        if docker info > /dev/null 2>&1; then
          sv=$(docker version --format '{{.Server.Version}}' 2> /dev/null) || sv=
          if [[ $sv =~ ^([0-9]+) ]]; then
            (( 10#${BASH_REMATCH[1]} >= 24 )) || meld "docker: der Server $sv ist älter als 24"
          else
            meld "docker: Version des Servers unlesbar"
          fi
        else
          meld "docker: der Daemon ist nicht erreichbar"
        fi ;;
    esac
  done
  (( bad == 0 )) || exit 3
}

# Punkte 5 und 6: zustand liest nur. Reihenfolge der Gründe: Lage im Repo, Verzeichnis, Schlüsselpaar.
zust_repo() {
  wurzel=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)
  aufgeloest=$(realpath -m -- "$PGR_DEV")
  if [[ $aufgeloest == "$wurzel" || $aufgeloest == "$wurzel"/* ]]; then
    ende 4 "$PGR_DEV liegt im Repo $wurzel"
  fi
}
zust_verzeichnis() {
  eltern=$(dirname -- "$PGR_DEV")
  [ -d "$eltern" ] || ende 4 "der Elternpfad $eltern fehlt oder ist kein Verzeichnis"
  zv=neu
  if [ -e "$PGR_DEV" ] || [ -L "$PGR_DEV" ]; then
    [ -d "$PGR_DEV" ] || ende 4 "$PGR_DEV ist kein Verzeichnis"
    if [ -n "$(ls -A -- "$PGR_DEV")" ]; then
      if grep -qxF "    name: $name" "$PGR_DEV/u-boot.yaml" 2> /dev/null; then
        zv=fortsetzen
        [ -f "$PGR_DEV/.devcontainer/devcontainer.json" ] && zv=vorhanden
      else
        ende 4 "$PGR_DEV ist fremd (kein u-boot.yaml mit dem Namen $name), nichts wird überschrieben"
      fi
    fi
  fi
}
zust_schluessel() {
  sd=$PGR_DEV/ssh priv=$PGR_DEV/ssh/id_ed25519
  pub=$priv.pub
  if { [ -e "$sd" ] || [ -L "$sd" ]; } && [ ! -d "$sd" ]; then ende 4 "$sd ist kein Verzeichnis"; fi
  zs=vorhanden
  hp=0 hb=0
  { [ -e "$priv" ] || [ -L "$priv" ]; } && hp=1
  { [ -e "$pub" ] || [ -L "$pub" ]; } && hb=1
  if (( hp + hb == 0 )); then
    zs=erzeugen
  elif (( hp + hb == 1 )); then
    ende 4 "Schlüsselpaar $sd: nur eine der beiden Dateien liegt, nichts wird überschrieben"
  else
    mangel=()
    (( 8#$(stat -c %a -- "$sd") & 077 )) && mangel+=("(a) $sd ist für Gruppe oder andere zugänglich")
    (( 8#$(stat -c %a -- "$priv") & 077 )) && mangel+=("(a) $priv ist für Gruppe oder andere zugänglich")
    if pk=$(ssh-keygen -y -P '' -f "$priv" 2> /dev/null < /dev/null); then
      read -r t1 k1 _ <<< "$pk"
      read -r t2 k2 _ < "$pub"
      [ "$t1 $k1" = "$t2 $k2" ] || mangel+=("(c) die öffentliche Datei $pub gehört nicht zur privaten")
    else
      mangel+=("(b) $priv ist ohne Passphrase nicht lesbar")
    fi
    (( ${#mangel[@]} == 0 )) || ende 4 "Schlüsselpaar abgelehnt, nichts wird verändert: ${mangel[*]}"
  fi
}
# Punkt 9: Der Container des Gerüsts trägt das Label devcontainer.local_folder gleich $PGR_DEV. Reihenfolge der Ablehnungen
# (alle Exit 4): Kollision des Namens, anderer Name, Port, Zustand, Verzeichnis ohne Gerüst. zc ist keiner, gestoppt oder laeuft.
# Ein Fehler von docker ist nie "kein Container".
zust_container() {
  local lab nm n i bind re=${PGR_NAME//./\\.}
  zc=keiner
  lab=$(docker ps -a --filter "label=devcontainer.local_folder=$PGR_DEV" --format '{{.Names}}') || ende 4 "docker ps (Label) endete mit Exit ungleich 0"
  nm=$(docker ps -a --filter "name=^/$re\$" --format '{{.Names}}') || ende 4 "docker ps (Name) endete mit Exit ungleich 0"
  if [ -n "$nm" ] && [[ $'\n'$lab$'\n' != *$'\n'$PGR_NAME$'\n'* ]]; then
    ende 4 "der Container $nm trägt das Label des Gerüsts nicht (Kollision), nichts wird angefasst"
  fi
  while IFS= read -r n; do
    [ -z "$n" ] || [ "$n" = "$PGR_NAME" ] || ende 4 "der Container des Gerüsts heißt $n, nicht $PGR_NAME"
  done <<< "$lab"
  [ -n "$lab" ] || return 0
  i=$(docker inspect --format '{{.State.Status}} {{range (index .HostConfig.PortBindings "2222/tcp")}}{{.HostIp}}:{{.HostPort}} {{end}}' "$PGR_NAME") || ende 4 "docker inspect endete mit Exit ungleich 0"
  bind=${i#* }
  [ "$bind" = "127.0.0.1:$PGR_SSH_PORT " ] || ende 4 "der Container $PGR_NAME veröffentlicht Port 2222 als ${bind:-keinen}, nicht als 127.0.0.1:$PGR_SSH_PORT"
  case "${i%% *}" in
    running) zc=laeuft ;;
    exited|created) zc=gestoppt ;;
    *) ende 4 "der Container $PGR_NAME hat den Zustand ${i%% *}, nichts wird angefasst" ;;
  esac
  [ -f "$PGR_DEV/.devcontainer/devcontainer.json" ] || ende 4 "der Container $PGR_NAME besteht, aber $PGR_DEV trägt kein Gerüst (.devcontainer/devcontainer.json), nichts wird angefasst"
}
# Punkt 10: Läuft der Container nicht, darf am Port keine Verbindung zustande kommen.
zust_port() {
  [ "$zc" != laeuft ] || return 0
  ! ( exec 3<> "/dev/tcp/127.0.0.1/$PGR_SSH_PORT" ) 2> /dev/null || ende 4 "der Port $PGR_SSH_PORT ist belegt"
}
s_zustand() { zust_repo; zust_verzeichnis; zust_schluessel; zust_container; zust_port; }

# Punkt 7: geruest. Die Ausgabe von u-boot erscheint unverändert; ein Fehler beendet den Lauf mit Exit 5.
ub() { kind u-boot "$@"; [ "$kcode" = 0 ] || ende 5 "u-boot $* endete mit Exit $kcode"; }
s_geruest() {
  [ "$zv" != vorhanden ] || return 0
  [ -d "$PGR_DEV" ] || mkdir -- "$PGR_DEV" || ende 5 "$PGR_DEV ließ sich nicht anlegen"
  cd -- "$PGR_DEV" || ende 5 "$PGR_DEV ist als Arbeitsverzeichnis des Gerüsts nicht erreichbar"
  [ "$zv" = neu ] && ub init "$name" --devcontainer --sandbox --no-git
  ub config set devcontainer.features.node.enabled true
  ub config set devcontainer.features.sshd.enabled true
  ub config set devcontainer.features.claude-code.enabled true
  ub config set devcontainer.sandbox.nestedRuntime docker
  ub config set devcontainer.sandbox.egress.enabled true
  ub config set devcontainer.sandbox.egress.allow gcr.io,pkg-containers.githubusercontent.com
  ub generate devcontainer
}

# Punkt 6: schluessel legt ein fehlendes Paar ohne Passphrase an und fasst ein vorhandenes nie an.
s_schluessel() {
  [ "$zs" = erzeugen ] || return 0
  { mkdir -p -- "$sd" && chmod 700 -- "$sd"; } || ende 5 "$sd ließ sich nicht anlegen"
  ssh-keygen -q -t ed25519 -N '' -C "$name" -f "$priv" > /dev/null < /dev/null || ende 5 "ssh-keygen endete mit Exit $?"
}

# Punkt 8: override. Ein vorhandenes mk-override.sh bleibt unberührt; es läuft mit den vier Variablen.
s_override() {
  [ "$zc" != laeuft ] || return 0
  cd -- "$PGR_DEV" || ende 5 "$PGR_DEV ist als Ort von mk-override.sh nicht erreichbar"
  if [ ! -e mk-override.sh ]; then
    cat > mk-override.sh <<'MKOVERRIDE'
#!/bin/bash
# erzeugt .devcontainer/override.json aus der von u-boot erzeugten devcontainer.json
set -euo pipefail
cd "$(dirname "$0")"
node -e '
const fs = require("fs");
const j = JSON.parse(fs.readFileSync(".devcontainer/devcontainer.json", "utf8").replace(/^\s*\/\/.*$/mg, ""));
const port = process.env.PGR_SSH_PORT || "2222";
j.runArgs = [...j.runArgs, "--name=" + (process.env.PGR_NAME || "pgr-dev"),
  "--cpus=" + (process.env.PGR_CPUS || "8"), "--memory=" + (process.env.PGR_MEM || "8g")];
j.appPort = ["127.0.0.1:" + port + ":2222"];
j.mounts = [...j.mounts, "type=bind,source=${localWorkspaceFolder}/ssh/id_ed25519.pub,target=/home/vscode/.ssh/authorized_keys,readonly"];
j.remoteEnv.CLAUDE_CODE_OAUTH_TOKEN = "${localEnv:CLAUDE_CODE_OAUTH_TOKEN}";
fs.writeFileSync(".devcontainer/override.json", JSON.stringify(j, null, 2) + "\n");
'
MKOVERRIDE
    chmod 755 mk-override.sh
  fi
  export PGR_NAME PGR_SSH_PORT PGR_CPUS PGR_MEM
  bash ./mk-override.sh || ende 5 "mk-override.sh endete mit Exit $?"
}

# Punkt 9: start, im Verzeichnis, das override betreten hat. Ein laufender Container bleibt stehen; sonst legt
# devcontainer up ihn an oder startet ihn.
# Erfolg ist Exit 0 und die Ausgabe "outcome":"success"; ohne Rückbau bei einem Fehler.
s_start() {
  [ "$zc" != laeuft ] || return 0
  kind devcontainer up --workspace-folder . --override-config .devcontainer/override.json
  [ "$kcode" = 0 ] || ende 5 "devcontainer up endete mit Exit $kcode"
  [[ $kausg == *'"outcome":"success"'* ]] || ende 5 'devcontainer up meldete kein "outcome":"success"'
}

# Punkt 11: warten. ssh als vscode mit dem Schlüssel des Gerüsts, jede Sekunde, bis PGR_WARTEN Sekunden verstrichen
# sind. Die zwei Optionen für den Hostschlüssel gelten nur für diese Probe; das Ziel schreibt kein known_hosts.
s_warten() {
  local t0=$SECONDS
  while :; do
    kind -m ssh -i "$PGR_DEV/ssh/id_ed25519" -o IdentitiesOnly=yes -o BatchMode=yes -o ConnectTimeout=5 \
      -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR -p "$PGR_SSH_PORT" vscode@127.0.0.1 true
    if [ "$kcode" = 0 ]; then echo "container-start: bereit $PGR_NAME 127.0.0.1:$PGR_SSH_PORT"; return 0; fi
    (( SECONDS - t0 < PGR_WARTEN )) || ende 5 "nach $PGR_WARTEN Sekunden keine Anmeldung, letzte Ausgabe von ssh: ${kausg%$'\n'}"
    kt=$SECONDS; while (( SECONDS == kt )); do sleep 0.1; done # die nächste Probe beginnt mit der nächsten Sekunde von $SECONDS
  done
}

# Punkt 13: Die Reihenfolge der Gründe ist die Reihenfolge der Schritte; der erste Verstoß beendet den Lauf.
for schritt in eingaben werkzeuge zustand geruest schluessel override start warten; do "s_$schritt"; done
