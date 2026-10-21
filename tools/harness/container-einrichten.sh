#!/usr/bin/env bash
# container-einrichten — richtet am Host für einen laufenden Dev-Container in einem Aufruf ein: Bare-Repo,
# eingelieferter main, Arbeitsklon, Identität des Klons, Commit-Träger und Zweig des Slice. Schritte:
# eingaben, werkzeuge, zustand, bare, einliefern, klon, identitaet, hooks, zweig, traeger; Erfolg ist die
# Zeile `container-einrichten: eingerichtet <Klon> <Zweig>` auf stdout.
# Vertrag: harness/vertraege/SPEC-069-container-einrichten.md (Punkte 1 bis 20). KEIN GATE: das Ziel steht im
# Abschnitt Werkzeuge von harness/README.md, und kein Hook und kein Stand glaubt sein Ergebnis.
# Eingaben aus der Umgebung: PGR_DEV und PGR_ZWEIG (Pflicht), PGR_SSH_PORT, PGR_REMOTE, PGR_SSH_OPTS.
# Voraussetzung am Host: bash ab 4.4 (`wait` auf eine Prozess-Substitution über `$!`, Punkt 19 und Kopf der
# Kennung), `ssh` und `git` im Pfad (Punkt 3).
# Exit-Codes: 0 eingerichtet · 2 eingaben · 3 werkzeuge am Host · 4 zustand · 5 Fehler eines Schrittes ·
# 130 SIGINT · 143 SIGTERM.
# GRENZEN (Punkt 20): keine Versionsprüfung; keine Zeitgrenze außer ConnectTimeout=5 je ssh-Aufruf; keine
# Einstellung für den Hostschlüssel; gleichzeitige Aufrufe mit demselben PGR_DEV sind nicht zugesagt; eine
# Option in PGR_SSH_OPTS, die Leerraum oder ein Shell-Metazeichen enthält, ist im Push nicht zugesagt
# (Punkt 20 (7)); der Exit-Code von `git push` bei gescheiterter Verbindung ist nicht gebunden
# (Punkt 20 (11)). Der Zeichenvorrat von PGR_DEV ist dagegen geprüft (Punkt 2), weil der Pfad in Punkt 7 in
# GIT_SSH_COMMAND reist.
#
# GEPRÜFT DURCH tools/harness/container-einrichten-gegenprobe.sh (Fälle p<Punkt>-*), ohne Container.
set -uo pipefail
export LC_ALL=C
# Punkt 13: Kein Kindprozess (ssh, git, make im Container) erbt einen Zugang von Claude Code.
unset CLAUDE_CODE_OAUTH_TOKEN
unset ANTHROPIC_API_KEY
unset GIT_TOKEN

schritt=eingaben
fehl() { echo "container-einrichten: $schritt: $*" >&2; }
ende() { local c=$1; shift; fehl "$@"; exit "$c"; }
bad=0
meld() { fehl "$@"; bad=1; }
q() { printf '%q' "$1"; }

# Punkt 15: SIGINT und SIGTERM beenden das laufende Kind und den Lauf; nichts wird zurückgebaut.
kpid=
abbruch() { [ -z "$kpid" ] || kill "$kpid" 2> /dev/null; ende "$1" abgebrochen; }
trap 'abbruch 130' INT
trap 'abbruch 143' TERM
exec 4>&1
kein=
# kind [-a|-e] <Befehl…>: Kindprozess im Hintergrund, damit ein Signal sofort wirkt. Ohne Schalter erscheinen
# beide Ströme unverändert dort, wo das Werkzeug sie schrieb (Punkt 4); -a hält beide in $kausg (der Lesezug
# ist Eingabe des Ziels, Punkt 5 — was darin nicht zu den Zeilen des Lesebefehls gehört, schreibt
# zust_lesezug auf stderr), -e hält stderr in $kausg und schreibt es danach unverändert auf stderr.
# $kein ist der Standardeingang (leer: nichts), $kcode der Exit-Code.
kind() {
  local m=${1-} fd z
  case $m in -a | -e) shift ;; *) m= ;; esac
  kausg=
  case $m in
    -a) exec {fd}< <(exec "$@" 2>&1 < <(printf '%s' "$kein")) ;;
    -e) exec {fd}< <(exec "$@" 2>&1 1>&4 < <(printf '%s' "$kein")) ;;
    *) "$@" < <(printf '%s' "$kein") & ;;
  esac
  kpid=$!
  [ -z "$m" ] || { while IFS= read -r z; do kausg+=$z$'\n'; done <&"$fd"; exec {fd}<&-; }
  wait "$kpid"; kcode=$?; kpid=
  [ "$m" != -e ] || printf '%s' "$kausg" >&2
}
# fern [-a] <Befehl im Container>: der ssh-Aufruf aus Punkt 5 (5) — BatchMode=yes in jedem Aufruf, keine
# Einstellung für den Hostschlüssel, PGR_SSH_OPTS als einzelne Argumente.
fern() {
  local m=${1-}; case $m in -a) shift ;; *) m= ;; esac
  kind ${m:+$m} ssh -i "$PGR_DEV/ssh/id_ed25519" -o IdentitiesOnly=yes -o BatchMode=yes -o ConnectTimeout=5 \
    ${sshopts[@]+"${sshopts[@]}"} -p "$PGR_SSH_PORT" vscode@127.0.0.1 "$1"
}

# zahl <Wert> <Höchstwert>: Dezimalziffern, kein Vorzeichen, keine führende Null, 1 bis Höchstwert.
zahl() {
  [[ -n $1 ]] || return 1
  [[ $1 != *[!0-9]* ]] || return 1
  [[ $1 != 0?* ]] || return 1
  [[ $1 != 0 ]] || return 1
  (( ${#1} < ${#2} )) || { (( ${#1} == ${#2} )) && [[ ! $1 > $2 ]]; }
}

# Punkt 2: eingaben. Nicht gesetzt ergibt die Vorgabe, gesetzt und leer ist ungültig; PGR_SSH_OPTS ist die
# Ausnahme der Tabelle (gesetzt und leer gilt wie nicht gesetzt).
PGR_SSH_PORT=${PGR_SSH_PORT-2222}
PGR_REMOTE=${PGR_REMOTE-container}
PGR_SSH_OPTS=${PGR_SSH_OPTS-}
pruefe_PGR_DEV() {
  if [ -z "${PGR_DEV+x}" ]; then meld "PGR_DEV ist nicht gesetzt"
  elif [ -z "$PGR_DEV" ]; then meld "PGR_DEV ist leer"
  elif [[ $PGR_DEV != /* ]]; then meld "PGR_DEV $(q "$PGR_DEV") ist nicht absolut"
  elif [[ ! $PGR_DEV =~ ^/[A-Za-z0-9._+/-]*$ ]]; then meld "PGR_DEV $(q "$PGR_DEV") trägt ein Zeichen außerhalb von A-Za-z0-9._+- und /"
  elif [ ! -e "$PGR_DEV" ]; then meld "PGR_DEV $(q "$PGR_DEV") fehlt"
  elif [ ! -d "$PGR_DEV" ]; then meld "PGR_DEV $(q "$PGR_DEV") ist kein Verzeichnis"
  fi
}
pruefe_PGR_SSH_PORT() {
  zahl "$PGR_SSH_PORT" 65535 || meld "PGR_SSH_PORT $(q "$PGR_SSH_PORT") ist keine ganze Zahl von 1 bis 65535"
}
pruefe_PGR_REMOTE() {
  if [[ ! $PGR_REMOTE =~ ^[A-Za-z0-9][A-Za-z0-9._-]*$ ]]; then meld "PGR_REMOTE $(q "$PGR_REMOTE") ist leer, beginnt nicht mit A-Za-z0-9 oder trägt ein Zeichen außerhalb von A-Za-z0-9._-"
  elif (( ${#PGR_REMOTE} > 64 )); then meld "PGR_REMOTE $(q "$PGR_REMOTE") ist länger als 64 Zeichen"
  fi
}
# Fehlt git, bleibt die Prüfung von check-ref-format aus; den Mangel meldet der Schritt werkzeuge (Punkt 3).
pruefe_PGR_ZWEIG() {
  if [ -z "${PGR_ZWEIG+x}" ]; then meld "PGR_ZWEIG ist nicht gesetzt"
  elif [ -z "$PGR_ZWEIG" ]; then meld "PGR_ZWEIG ist leer"
  elif [ "$PGR_ZWEIG" = main ]; then meld "PGR_ZWEIG main ist nicht zugelassen: der Klon arbeitet nie auf main"
  elif [[ $PGR_ZWEIG == -* ]]; then meld "PGR_ZWEIG $(q "$PGR_ZWEIG") beginnt mit -"
  elif [[ $PGR_ZWEIG == refs/* ]]; then meld "PGR_ZWEIG $(q "$PGR_ZWEIG") beginnt mit refs/ und benennt im Klon nicht den gemeinten Zweig"
  elif command -v git > /dev/null 2>&1 && ! git check-ref-format --branch "$PGR_ZWEIG" > /dev/null 2>&1; then
    meld "PGR_ZWEIG $(q "$PGR_ZWEIG") ist kein Zweigname (git check-ref-format --branch)"
  fi
}
# Tabulator ist Trenner und zulässig; Zeilenumbruch und jedes andere Steuerzeichen sind ungültig.
pruefe_PGR_SSH_OPTS() {
  [[ $PGR_SSH_OPTS != *[$'\001'-$'\010'$'\012'-$'\037\177']* ]] || meld "PGR_SSH_OPTS trägt ein Steuerzeichen (Tabulator ausgenommen)"
}
s_eingaben() {
  local v
  for v in PGR_DEV PGR_SSH_PORT PGR_REMOTE PGR_ZWEIG PGR_SSH_OPTS; do "pruefe_$v"; done
  (( bad == 0 )) || exit 2
  # Punkt 2: an Leerraum geteilt und als einzelne Argumente übergeben; `set -f` hält die
  # Pfadnamen-Expansion der Shell aus dem Wert heraus.
  set -f
  sshopts=($PGR_SSH_OPTS)
  set +f
}

# Punkt 3: werkzeuge am Host. Alle Verstöße des Schrittes werden gemeldet, in dieser Reihenfolge; Versionen
# prüft das Ziel nicht (Punkt 20 (1)).
s_werkzeuge() {
  local t
  for t in ssh git; do command -v "$t" > /dev/null 2>&1 || meld "$t fehlt im Pfad"; done
  (( bad == 0 )) || exit 3
}

# Punkt 5: zustand liest, bevor etwas geändert wird — am Host und mit genau EINEM ssh-Aufruf im Container.
# Reihenfolge der Gründe (Punkt 14): Gerüst, Host-Repo, Identität, Remote, Lesezug, Werkzeuge im Container,
# Bare-Repo, Klon. Jeder Grund endet den Lauf mit Exit 4.
# Punkt 3: Der Projektname kommt aus der ersten Zeile `    name: <Name>` von $PGR_DEV/u-boot.yaml, nicht aus
# dem Wert der Variablen; er bestimmt den Pfad im Container.
zust_geruest() {
  local z y=$PGR_DEV/u-boot.yaml b=${PGR_DEV##*/}
  projekt=
  [ -f "$y" ] || ende 4 "$y fehlt — kein Gerüst; make container-start legt es an"
  while IFS= read -r z; do [[ $z == '    name: '* ]] && { projekt=${z#'    name: '}; break; }; done < "$y"
  [ -n "$projekt" ] || ende 4 "$y trägt keine Zeile '    name: <Name>' — kein Gerüst; make container-start legt es an"
  [[ $projekt =~ ^[a-z0-9][a-z0-9-]*$ ]] || ende 4 "$y nennt den Namen $(q "$projekt") außerhalb von a-z0-9-"
  [ "$projekt" = "$b" ] || ende 4 "$y nennt den Namen $projekt, nicht den Basisnamen $b von PGR_DEV"
  [ -e "$PGR_DEV/ssh/id_ed25519" ] || ende 4 "$PGR_DEV/ssh/id_ed25519 fehlt"
  [ -r "$PGR_DEV/ssh/id_ed25519" ] || ende 4 "$PGR_DEV/ssh/id_ed25519 ist nicht lesbar"
  barepfad=/workspaces/$projekt/austausch.git
  klonpfad=/workspaces/$projekt/pgwire-recorder
  url=ssh://vscode@127.0.0.1:$PGR_SSH_PORT/workspaces/$projekt/austausch.git
}
# Ein detached HEAD und ein anderer ausgecheckter Zweig sind kein Grund: eingeliefert wird immer
# refs/heads/main (Randform 12).
zust_repo() {
  git rev-parse --show-toplevel > /dev/null 2>&1 || ende 4 "$PWD liegt nicht in einem Git-Repo"
  git rev-parse -q --verify refs/heads/main > /dev/null 2>&1 || ende 4 "refs/heads/main fehlt im Host-Repo"
}
# Die Quelle darf global sein (Randform 7); das Ziel setzt die Werte später nur --local (Punkt 9). Die Meldung
# nennt den Namen der Einstellung, nie den Wert — er kann der Name einer Person sein und geht in Logs.
zust_identitaet() {
  local f=() re=$'^([\x20-\x7e]|[\xc2-\xdf][\x80-\xbf]|\xe0[\xa0-\xbf][\x80-\xbf]|[\xe1-\xec\xee\xef][\x80-\xbf][\x80-\xbf]|\xed[\x80-\x9f][\x80-\xbf]|\xf0[\x90-\xbf][\x80-\xbf][\x80-\xbf]|[\xf1-\xf3][\x80-\xbf][\x80-\xbf][\x80-\xbf]|\xf4[\x80-\x8f][\x80-\xbf][\x80-\xbf])*$'
  iname=$(git config --get user.name 2> /dev/null) || iname=
  imail=$(git config --get user.email 2> /dev/null) || imail=
  [ -n "${iname//[ $'\t']/}" ] || f+=(user.name)
  [ -n "${imail//[ $'\t']/}" ] || f+=(user.email)
  (( ${#f[@]} == 0 )) || ende 4 "die Identität des Host-Repos fehlt: ${f[*]}"
  [[ $iname =~ $re ]] || ende 4 "user.name des Host-Repos trägt ein Steuerzeichen oder kein gültiges UTF-8 (der Wert steht nicht in dieser Meldung)"
  [[ $imail =~ $re ]] || ende 4 "user.email des Host-Repos trägt ein Steuerzeichen oder kein gültiges UTF-8 (der Wert steht nicht in dieser Meldung)"
}
# Nichts wird gesetzt oder entfernt (Randform 3); fehlt das Remote, legt Punkt 7 es an.
zust_remote() {
  local ist
  rneu=1
  if ist=$(git remote get-url "$PGR_REMOTE" 2> /dev/null); then
    [ "$ist" = "$url" ] || ende 4 "der Remote $PGR_REMOTE zeigt auf $ist, nicht auf $url"
    rneu=0
  fi
}
# Der Lesebefehl schreibt nichts und gibt Zeilen <schlüssel>=<wert> aus. Der Zweig kommt als Liste der
# Zweige zurück und wird am Host verglichen, damit kein Wert von außen in eine Shell des Containers gerät.
lesebefehl() {
  printf '%s\n' "b=$barepfad" "k=$klonpfad" \
    'for t in git make; do if command -v $t > /dev/null 2>&1; then echo "wz-$t=ja"; else echo "wz-$t=nein"; fi; done' \
    'if [ ! -e "$b" ] && [ ! -L "$b" ]; then echo bare=fehlt' \
    'elif [ "$(git -C "$b" rev-parse --is-bare-repository 2> /dev/null)" = true ]; then echo "bare=head:$(git -C "$b" symbolic-ref -q HEAD)"' \
    'else echo bare=fremd; fi' \
    'if [ ! -e "$k" ] && [ ! -L "$k" ]; then echo klon=fehlt' \
    'elif [ "$(git -C "$k" rev-parse --is-inside-work-tree 2> /dev/null)" != true ]; then echo klon=fremd' \
    'else echo klon=repo; echo "origin=$(git -C "$k" remote get-url origin 2> /dev/null)"' \
    'if [ -n "$(git -C "$k" status --porcelain)" ]; then echo status=schmutzig; else echo status=sauber; fi' \
    'git -C "$k" for-each-ref --format="zweig=%(refname:short)" refs/heads/; fi'
}
zust_lesezug() {
  local z s
  fern -a "$(lesebefehl)"
  if [ "$kcode" != 0 ]; then
    s=${kausg%$'\n'}; s=${s##*$'\n'}
    ende 4 "der Container antwortet nicht auf vscode@127.0.0.1, Port $PGR_SSH_PORT (ssh: $s) — make container-start startet ihn"
  fi
  zweige=$'\n'
  Z=()
  while IFS= read -r z; do
    case $z in
      zweig=*) zweige+=${z#zweig=}$'\n' ;;
      wz-git=* | wz-make=* | bare=* | klon=* | origin=* | status=*) Z[${z%%=*}]=${z#*=} ;;
      # Punkt 4: Was das gerufene Werkzeug neben den Zeilen des Lesebefehls schrieb, erscheint auf stderr.
      *) [ -z "$z" ] || printf '%s\n' "$z" >&2 ;;
    esac
  done <<< "$kausg"
}
# Beide liegen im Image, und ein Eingriff in das Image ist nicht Gegenstand dieses Ziels.
zust_werkzeuge_fern() {
  local t
  for t in git make; do
    [ "${Z[wz-$t]-}" = ja ] || ende 4 "$t fehlt im Pfad des Containers vscode@127.0.0.1:$PGR_SSH_PORT"
  done
}
# Ein Bare-Repo mit HEAD auf master ließe git clone ohne Arbeitskopie enden (Randform 1).
zust_bare() {
  case ${Z[bare]-} in
    fehlt | head:refs/heads/main) ;;
    head:*) ende 4 "$barepfad hat HEAD auf ${Z[bare]#head:}, nicht auf refs/heads/main" ;;
    *) ende 4 "$barepfad besteht und ist kein Bare-Repo" ;;
  esac
}
# Das Ziel überschreibt keine Arbeit (Randform 2).
zust_klon() {
  case ${Z[klon]-} in
    fehlt) ;;
    repo)
      [ "${Z[status]-}" = sauber ] || ende 4 "$klonpfad trägt eine uncommittete Änderung — committe oder stashe sie im Klon"
      [ "${Z[origin]-}" = "$barepfad" ] || ende 4 "$klonpfad hat origin ${Z[origin]:-keinen}, nicht $barepfad"
      ;;
    *) ende 4 "$klonpfad besteht und ist kein Git-Repo" ;;
  esac
  zweigda=0
  [[ $zweige != *$'\n'$PGR_ZWEIG$'\n'* ]] || zweigda=1
}
s_zustand() {
  declare -gA Z=()
  zust_geruest; zust_repo; zust_identitaet; zust_remote; zust_lesezug
  zust_werkzeuge_fern; zust_bare; zust_klon
}

# Punkt 6: bare. Ist es vorhanden und in Ordnung, folgt kein Aufruf (kein zweites git init, Randform 1).
s_bare() {
  [ "${Z[bare]}" = fehlt ] || return 0
  fern "git init --bare --initial-branch=main $barepfad"
  [ "$kcode" = 0 ] || ende 5 "git init --bare $barepfad endete mit Exit $kcode"
}

# Punkt 7: einliefern, am Host. core.sshCommand liest und schreibt das Ziel nie — die Einstellung gilt dem
# ganzen Repo, GIT_SSH_COMMAND genau diesem Aufruf. Kein --force, kein anderer Refspec.
s_einliefern() {
  local gsc="ssh -i $PGR_DEV/ssh/id_ed25519 -o IdentitiesOnly=yes -o BatchMode=yes${PGR_SSH_OPTS:+ $PGR_SSH_OPTS}"
  (( rneu == 0 )) || git remote add "$PGR_REMOTE" "$url" || ende 5 "git remote add $PGR_REMOTE $url endete mit Exit $?"
  kind -e env "GIT_SSH_COMMAND=$gsc" git push "$PGR_REMOTE" refs/heads/main:refs/heads/main
  [ "$kcode" = 0 ] && return 0
  [[ $kausg != *'[rejected]'* ]] || ende 5 "der Remote $PGR_REMOTE trägt einen main mit anderer Historie; hole ihn mit git fetch $PGR_REMOTE und führe zusammen — nichts wird erzwungen"
  ende 5 "git push $PGR_REMOTE endete mit Exit $kcode"
}

# Punkt 8: klon. Ein vorhandener Klon bekommt allein fetch --prune: er kann Arbeit tragen, und der Stand, aus
# dem Punkt 11 den Zweig bildet, ist origin/main.
s_klon() {
  if [ "${Z[klon]}" = fehlt ]; then
    fern "git clone $barepfad $klonpfad"
  else
    fern "git -C $klonpfad fetch --prune origin"
  fi
  [ "$kcode" = 0 ] || ende 5 "der Klon $klonpfad endete mit Exit $kcode"
}

# Punkt 9: identitaet, nur --local im Klon, nie --global und nie --system, nie ein Name aus dem Skript. Die
# beiden Werte reisen über den Standardeingang des ssh-Aufrufs, je einer in einer Zeile; im Container liest
# sie `IFS= read -r` und gibt sie als Argument an git config. Kein Wert steht in einer Kommandozeile oder in
# einer Zeichenkette, die eine Shell auswertet.
s_identitaet() {
  kein=$iname$'\n'$imail$'\n'
  fern 'IFS= read -r n; IFS= read -r e; git -C '"$klonpfad"' config --local user.name "$n" && git -C '"$klonpfad"' config --local user.email "$e"'
  kein=
  [ "$kcode" = 0 ] || ende 5 "die Identität des Klons $klonpfad endete mit Exit $kcode"
}

# Punkt 10: hooks. Läuft bei jedem Aufruf, auch beim zweiten (core.hooksPath gilt je Klon).
s_hooks() {
  fern "cd $klonpfad && make hooks-install"
  [ "$kcode" = 0 ] || ende 5 "make hooks-install im Klon endete mit Exit $kcode"
}

# Punkt 11: zweig. Quelle ist immer origin/main nach dem fetch aus Punkt 8, nie der lokale main des Klons;
# ein vorhandener Zweig wird gewechselt, nicht erzwungen (Randform 16). Der Name reist über den
# Standardeingang, damit kein Wert von außen in eine Shell des Containers gerät.
s_zweig() {
  kein=$PGR_ZWEIG$'\n'
  if (( zweigda )); then
    fern 'IFS= read -r z; git -C '"$klonpfad"' switch "$z"'
  else
    fern 'IFS= read -r z; git -C '"$klonpfad"' switch -c "$z" origin/main'
  fi
  kein=
  [ "$kcode" = 0 ] || ende 5 "der Zweig $PGR_ZWEIG im Klon endete mit Exit $kcode"
}

# Punkt 12: traeger, als LETZTER Schritt. Der Schritt braucht Netz im Container, und das Egress des Gerüsts
# ist eingeschränkt; scheitert er, stehen Bare-Repo, Klon, Identität, Träger-Pfad und Zweig schon, und die
# Wiederholung ist ein zweiter Aufruf des Ziels (Randform 10).
s_traeger() {
  fern "cd $klonpfad && make traeger-fetch"
  [ "$kcode" = 0 ] || ende 5 "make traeger-fetch im Klon endete mit Exit $kcode"
}

# Punkt 4 und 14: Die Reihenfolge der Gründe ist die Reihenfolge der Schritte; der erste Verstoß beendet den
# Lauf. Punkt 13: Erfolg ist eine Zeile auf stdout und Exit 0, beim ersten wie beim zweiten Aufruf.
# Die Stelle für die Erweiterung um einen Schritt `zugang` ist vor `eingaben` (Punkt 13, Randform 25).
for schritt in eingaben werkzeuge zustand bare einliefern klon identitaet hooks zweig traeger; do "s_$schritt"; done
echo "container-einrichten: eingerichtet $klonpfad $PGR_ZWEIG"
