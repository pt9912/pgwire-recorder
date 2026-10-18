#!/usr/bin/env bash
# gegenproben-nach-einlieferung — Inhalt des post-receive-Hooks im Bare-Repo des Gegenproben-Containers (Vertrag:
# harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md, Punkte 1 bis 14 und 16). Kein Gate.
# Aufruf: als Kopie (`install -m 755`) unter <Bare-Repo>/hooks/post-receive; git liefert die Zeilen `<alt> <neu> <ref>` auf stdin.
# Je eingeliefertem refs/heads/*-Zweig fährt ein Hintergrundlauf (dieselbe Datei mit --lauf, eigene Sitzung) `make gegenproben` auf
# genau dem Hash und legt bei beendetem `make` die Ergebnisdatei ab. Der Hook endet in jedem Fall mit Exit 0 und braucht keine
# andere Datei des Repos. GEPRÜFT DURCH tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh.
sig="gegenproben-nach-einlieferung"
melde() { echo "$sig: $1: $2" >&2; }

# kodiere <Zweig> — jedes Byte außer A-Za-z0-9._- als % und zwei Großbuchstaben-Hexziffern (Punkt 8).
kodiere() {
  local s="$1" i c out=""
  local LC_ALL=C
  for ((i = 0; i < ${#s}; i++)); do
    c="${s:i:1}"
    case "$c" in
      [A-Za-z0-9._-]) out+="$c" ;;
      *) out+="$(printf '%%%02X' "'$c")" ;;
    esac
  done
  printf '%s' "$out"
}

haupt() {
  local self bare marke zeilen ablage t probe zeile neu ref rest zweig name n v args=()
  self="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/$(basename "${BASH_SOURCE[0]}")"
  bare="$(cd "${GIT_DIR:-.}" 2> /dev/null && pwd -P)"
  for v in $(compgen -v GIT_); do unset "$v"; done
  unset MAKEFLAGS MFLAGS GNUMAKEFLAGS MAKELEVEL MAKEOVERRIDES
  # (1) Marke des Containers (Punkt 2)
  marke="$bare/gegenproben-ablage"
  mapfile -t zeilen 2> /dev/null < "$marke"
  [ "${#zeilen[@]}" -eq 1 ] || { melde marke "keine einzelne Zeile in $marke"; return 0; }
  ablage="${zeilen[0]}"
  case "$ablage" in /*) ;; *) melde marke "kein absoluter Pfad: $ablage"; return 0 ;; esac
  [ -d "$ablage" ] || { melde marke "kein Verzeichnis: $ablage"; return 0; }
  # (2) Werkzeuge, (3) Ablage (Punkt 13)
  for t in git flock setsid make pgrep; do command -v "$t" > /dev/null 2>&1 || { melde werkzeug "$t fehlt"; return 0; }; done
  probe="$(mktemp "$ablage/.probe.XXXXXX" 2> /dev/null)" || { melde ablage "nicht schreibbar: $ablage"; return 0; }; rm -f "${probe:?}"
  # (4) je Zeile (Punkte 3, 4, 8)
  while IFS= read -r zeile || [ -n "$zeile" ]; do
    [[ -n "${zeile//[[:space:]]/}" ]] || continue
    read -r _ neu ref rest <<< "$zeile"
    [ -n "$ref" ] || { melde zeile "Feld fehlt: $zeile"; continue; }
    [ -z "$rest" ] || { melde zeile "mehr als drei Felder: $zeile"; continue; }
    case "$ref" in refs/heads/*) ;; *) continue ;; esac
    [[ $neu =~ ^(0{40}|0{64})$ ]] && continue
    [[ $neu =~ ^([0-9a-f]{40}|[0-9a-f]{64})$ ]] || { melde zeile "Form des Hashes: $zeile"; continue; }
    zweig="${ref#refs/heads/}"
    name="$(kodiere "$zweig")"
    [ "${#name}" -le 120 ] || { melde name "länger als 120 Zeichen: $name"; continue; }
    n="gegenproben-$name-$neu"
    rm -f "${ablage:?}/${n:?}"
    : > "$ablage/$n.laeuft"
    melde gestartet "$zweig $neu"
    args+=("$n" "$zweig" "$neu")
  done
  if [ "${#args[@]}" -gt 0 ]; then
    # `setsid -f` statt `cmd &`: ein einfacher Hintergrundbefehl einer Shell ohne Jobsteuerung startet mit ignoriertem SIGINT und SIGQUIT, und seine Kinder erben es.
    # `env --default-signal=INT,QUIT` setzt beide für den Lauf zurück, auch wenn der Aufrufer des Hooks sie ignoriert (coreutils ab 8.31).
    # Das Kind schließt vor dem Lauf jeden Deskriptor über 2, den git an den Hook reicht (sonst hält der Lauf die Leitung des Pushs offen).
    setsid -f env --default-signal=INT,QUIT "$BASH" -c 'for f in /proc/self/fd/*; do n="${f##*/}"; [ "$n" -gt 2 ] 2> /dev/null && eval "exec $n>&-"; done; exec "$BASH" "$@"' \
      _ "$self" --lauf "$bare" "$ablage" "${args[@]}" < /dev/null > /dev/null 2>&1
  fi
  return 0
}

# --- Hintergrundlauf (Punkte 6, 9, 10, 12) -------------------------------------------------------------
pid=""
wt=""
# beende <pid> <signal> — das Signal an den Prozess und alle Nachkommen, Nachkommen zuerst.
beende() { local k; for k in $(pgrep -P "$1" 2> /dev/null); do beende "$k" "$2"; done; kill -s "$2" "$1" 2> /dev/null || true; }
# raeume — entfernt den Arbeitsbaum des laufenden Laufs.
raeume() {
  [ -z "$wt" ] || git --git-dir="$bare" worktree remove --force "$wt" > /dev/null 2>&1 || rm -rf "${wt:?}"
  git --git-dir="$bare" worktree prune > /dev/null 2>&1 || true
  wt=""
}
abbruch() { trap '' TERM; if [ -n "$pid" ]; then beende "$pid" TERM; wait "$pid" 2> /dev/null; fi; raeume; exit 143; }

lauf() {
  local n="$1" zweig="$2" hash="$3" log="$ablage/$1.log" beginn ende rc erg rot z tmp aus
  # (5) Sperre
  exec 9> "$ablage/gegenproben.lock" || { echo "$sig: sperre: Sperre nicht angelegt" > "$log"; return 0; }
  flock 9 & pid=$!
  wait "$pid"; rc=$?; pid=""
  [ "$rc" -eq 0 ] || { echo "$sig: sperre: flock endete mit $rc" > "$log"; exec 9>&-; return 0; }
  # (6) Arbeitsbaum
  wt="$(mktemp -d "$ablage/arbeitsbaum.XXXXXX" 2> /dev/null)" || { echo "$sig: arbeitsbaum: Verzeichnis nicht angelegt" > "$log"; wt=""; exec 9>&-; return 0; }
  if ! aus="$(git --git-dir="$bare" worktree add --detach "$wt" "$hash" 2>&1)"; then
    { echo "$aus"; echo "$sig: arbeitsbaum: Arbeitsbaum nicht angelegt"; } > "$log"
    rm -rf "${wt:?}"; wt=""; exec 9>&-; return 0
  fi
  # (7) der Lauf
  beginn="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  (cd "$wt" && exec make gegenproben) > "$log" 2>&1 9>&- &
  pid=$!
  wait "$pid"; rc=$?; pid=""
  ende="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  erg=rot; [ "$rc" -eq 0 ] && erg=gruen
  rot=""
  while IFS= read -r z || [ -n "$z" ]; do
    [[ $z =~ ^gegenproben:\ ([^:\ ]+):\ rot\ \(Exit\ [0-9]+\)$ ]] && rot+="${rot:+ }${BASH_REMATCH[1]}"
  done < "$log"
  tmp="$(mktemp "$ablage/.$n.XXXXXX")"
  printf 'zweig=%s\nhash=%s\nexit=%s\nergebnis=%s\nbeginn=%s\nende=%s\nrot=%s\n' "$zweig" "$hash" "$rc" "$erg" "$beginn" "$ende" "$rot" > "$tmp" \
    && mv -f "$tmp" "$ablage/$n"
  rm -f "${ablage:?}/${n:?}.laeuft"
  raeume
  exec 9>&-
}

lauf_haupt() {
  bare="$2"; ablage="$3"; shift 3
  trap abbruch TERM
  while [ "$#" -ge 3 ]; do lauf "$1" "$2" "$3"; shift 3; done
}

if [ "${1:-}" = "--lauf" ]; then lauf_haupt "$@"; exit 0; fi
haupt
exit 0
