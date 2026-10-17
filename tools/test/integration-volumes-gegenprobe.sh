#!/usr/bin/env bash
# integration-volumes-gegenprobe — prüft das Aufräumen der Integrationstests
# (tools/test/integration-aufraeumen.sh und tools/test/run-integration-tests.sh) gegen die Punkte
# 1 bis 8 von SPEC-057 (harness/vertraege/SPEC-057-integration-aufraeumen.md) mit realem Docker
# (über DOCKER). Container und Volumes tragen das Präfix pgr-ivg-<PID> und werden per Namen
# entfernt; kein prune, kein Eingriff in ein Objekt mit anderem Namen. Die unveränderte Kopie
# besteht jeden Fall; bei jeder Mutation (Textersatz in einer Kopie, genau eine Stelle) werden von
# den Fällen ihrer Gruppe (`nur`) genau die benannten rot. Der Name jedes Falls trägt die Nummer
# seines Punktes.
#
# Ein Fall prüft nach Namen: Er liest die Volumes seiner Container mit einem eigenen
# `docker container inspect` und fragt nach dem Aufräumen je Name `docker volume inspect`;
# die Zahl der Volumes des Hosts zählt er nie. Ein "leckendes" docker (Skript `leck-docker`)
# reicht jeden Aufruf an DOCKER weiter und ersetzt nur `rm -fv` durch `rm -f`; damit entsteht ein
# Rest, ohne dass der Fall den Code des Runners verändert.
#
# GRENZE: Das Testimage wird nicht gebraucht (PGR_BAUEN=0); ob die Tests laufen, zeigt
# `make test-integration`. Die Fälle p5-* senden das Signal an den Runner, solange er auf PostgreSQL
# wartet, und setzen voraus, dass diese Gegenprobe SIGINT nicht ignoriert (kein Start als
# Hintergrundjob). Ein roter Lauf und ein Fehler des Aufbaus gehen durch dieselbe Falle wie die
# Signale und haben keinen eigenen Fall (SPEC-057 Punkt 5).
#
# Zum Arbeiten: GP_MUT=<Muster> fährt nur die Mutationen, deren Name passt; GP_ENTDECKE=1 nennt je
# Mutation die roten Fälle; GP_FALL=<Muster> fährt nur die Fälle der unveränderten Kopie, deren Name passt.
#
# Ausgang: 0, wenn jeder Fall und jede Mutation wie erwartet endet, sonst 1.
set -euo pipefail

hier="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
docker_echt="${DOCKER:-docker}"
image="${POSTGRES_IMAGE:?POSTGRES_IMAGE fehlt}"
arbeit="$(mktemp -d)"
praefix="pgr-ivg-$$"
fehler=0
melde() { echo "integration-volumes-gegenprobe: ROT — $*" >&2; fehler=1; }

# Was die Gegenprobe angelegt hat, steht in diesen Dateien und wird am Ende per Namen entfernt.
c_liste="$arbeit/container" v_liste="$arbeit/volumes" n_liste="$arbeit/netze"
: > "$c_liste" > "$v_liste" > "$n_liste"
aus="$arbeit/aus" err="$arbeit/err"
export IVG_DOCKER="$docker_echt" IVG_SICHER="$v_liste"

endreinigung() {
  local n
  while read -r n; do [ -z "$n" ] || "$docker_echt" rm -fv "$n" > /dev/null 2>&1 || true; done < "$c_liste"
  while read -r n; do [ -z "$n" ] || "$docker_echt" network rm "$n" > /dev/null 2>&1 || true; done < "$n_liste"
  while read -r n; do [ -z "$n" ] || "$docker_echt" volume rm "$n" > /dev/null 2>&1 || true; done < "$v_liste"
  rm -rf "${arbeit:?}"
}
trap endreinigung EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

# leck-docker: wie DOCKER, aber `rm -fv` wird zu `rm -f` (der Rest, den der Runner finden muss).
# Mit IVG_NACH_RUN=<Sekunden> wartet `run` nach dem Start des Containers so lange: Der Runner steht
# dann in seinem `wait`, wenn das Signal kommt (die Falle für TERM entscheidet dort den Exit-Code).
leck="$arbeit/leck-docker"
cat > "$leck" <<'LECK'
#!/usr/bin/env bash
a=()
for x in "$@"; do
  if [ "$x" = -fv ]; then x=-f; fi
  a+=("$x")
done
if [ "$1" = run ] && [ -n "${IVG_NACH_RUN-}" ]; then
  "$IVG_DOCKER" "${a[@]}"
  r=$?
  sleep "$IVG_NACH_RUN"
  exit "$r"
fi
exec "$IVG_DOCKER" "${a[@]}"
LECK
chmod +x "$leck"

ZAEHLER_DATEI="$arbeit/zaehler"
echo 0 > "$ZAEHLER_DATEI"
v_existiert() { "$docker_echt" volume inspect "$1" > /dev/null 2>&1; }
# paar — frische Namen für einen Fall; setzt lauf, pg, tests, daten.
paar() {
  ZAEHLER=$(($(cat "$ZAEHLER_DATEI") + 1))
  echo "$ZAEHLER" > "$ZAEHLER_DATEI"
  lauf="$praefix-$ZAEHLER"
  pg="$lauf-pg" tests="$lauf-tests" daten="$lauf-daten"
  printf '%s\n%s\n' "$pg" "$tests" >> "$c_liste"
  echo "$lauf-net" >> "$n_liste"
  echo "$daten" >> "$v_liste"
}
# anlegen — legt $pg und $tests aus dem Image an (Zustand created, das Image deklariert VOLUME;
# $tests bindet außerdem das benannte Volume $daten ein, wie im Runner); setzt probe auf die Namen
# der anonymen Volumes, gelesen mit einem eigenen docker inspect.
anlegen() {
  local c n
  probe=()
  "$docker_echt" volume create "$daten" > /dev/null
  for c in "$pg" "$tests"; do
    if [ "$c" = "$tests" ]; then
      "$docker_echt" create --name "$c" -v "$daten":/daten "$image" > /dev/null
    else
      "$docker_echt" create --name "$c" "$image" > /dev/null
    fi
    for n in $("$docker_echt" container inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}' "$c"); do
      if [ "$n" = "$daten" ]; then continue; fi
      probe+=("$n")
      echo "$n" >> "$v_liste"
    done
  done
  [ "${#probe[@]}" -ge 2 ] || { melde "das Image legt kein Volume an: ${probe[*]-}"; return 1; }
}
# fremde — legt ein benanntes und ein anonymes Volume an, die kein Container des Laufs nutzt.
fremde() {
  fremd_benannt="$lauf-fremd"
  "$docker_echt" volume create "$fremd_benannt" > /dev/null
  fremd_anonym="$("$docker_echt" volume create)"
  echo "$fremd_benannt" >> "$v_liste"
  echo "$fremd_anonym" >> "$v_liste"
}
# lauf_in <Verzeichnis> <docker> <Befehl>… — der Befehl läuft in einem Unterprozess, der die
# Funktionsdatei des Verzeichnisses lädt; setzt code, Ausgabe in $aus und $err.
lauf_in() {
  local d="$1" dk="$2"
  shift 2
  set +e
  ( docker="$dk"; . "$d/integration-aufraeumen.sh"; "$@" ) > "$aus" 2> "$err"
  code=$?
  set -e
}
alle_weg() { local n; for n in "$@"; do if v_existiert "$n"; then return 1; fi; done; }
bleiben() { local n; for n in "$@"; do if ! v_existiert "$n"; then return 1; fi; done; }
in_err() { grep -qF -- "$1" "$err"; }

sel() {
  local p
  if [ -z "${nur:-}" ]; then return 0; fi
  for p in $nur; do if [[ "$1" == "$p"* ]]; then return 0; fi; done
  return 1
}
# fall <Name> <Funktion> <Verzeichnis> — schreibt den Namen, wenn die Funktion nicht 0 liefert.
fall() {
  if sel "$1"; then "$2" "$3" || echo "$1"; fi
  return 0
}

f_p2_grundfall() {
  paar; anlegen; fremde
  lauf_in "$1" "$docker_echt" integration_ende 0
  [ "$code" -eq 0 ] && [ ! -s "$err" ] && alle_weg "${probe[@]}" && bleiben "$fremd_benannt" "$fremd_anonym" "$daten"
}
f_p2_ohne_v() {
  paar; anlegen
  lauf_in "$1" "$docker_echt" integration_entfernen
  alle_weg "${probe[@]}"
}
f_p3_angelegt() {
  paar; anlegen
  lauf_in "$1" "$leck" integration_ende 0
  [ "$code" -eq 1 ] && in_err "${probe[0]}" && in_err "${probe[1]}" && alle_weg "${probe[@]}"
}
f_p3_kein_container() {
  paar
  lauf_in "$1" "$docker_echt" integration_ende 0
  [ "$code" -eq 0 ] && [ ! -s "$err" ] && [ ! -s "$aus" ]
}
f_p4_fremd() {
  paar; anlegen; fremde
  lauf_in "$1" "$leck" integration_ende 0
  bleiben "$fremd_benannt" "$fremd_anonym"
}
f_p4_rest_rot() {
  paar; anlegen
  lauf_in "$1" "$leck" integration_ende 0
  [ "$code" -eq 1 ]
}
f_p4_vorrang() {
  paar; anlegen
  lauf_in "$1" "$leck" integration_ende 143
  [ "$code" -eq 143 ]
}
# p5 <INT|TERM> <Verzeichnis> <erwarteter Code> — der Runner (Kopie) startet bis zum Container $pg;
# das Signal geht an ihn (nicht an seine Prozessgruppe), sobald $pg existiert. Mit dem leckenden docker bleibt das Volume von $pg
# ohne die Prüfung zurück.
f_p5() {
  local sg="$1" d="$2" soll="$3" marke="$arbeit/p5-gesendet" pidf="$arbeit/p5-pid" namen="$arbeit/p5-namen" wp rc n
  paar
  rm -f "$marke" "$pidf" "$namen"
  (
    limit=$((SECONDS + 60))
    until "$docker_echt" container inspect "$pg" > /dev/null 2>&1 && [ -s "$pidf" ]; do
      if [ "$SECONDS" -ge "$limit" ]; then exit 0; fi
      /bin/sleep 0.05
    done
    "$docker_echt" container inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}' "$pg" > "$namen"
    kill -"$sg" "$(cat "$pidf")"
    : > "$marke"
  ) &
  wp=$!
  set +e
  {
    env DOCKER="$leck" IVG_NACH_RUN=3 POSTGRES_IMAGE="$image" PGR_BAUEN=0 PGR_LAUF="$lauf" \
      setsid bash -c 'echo $$ > "$1"; shift; exec "$@"' x "$pidf" bash "$d/run-integration-tests.sh" > "$aus" 2> "$err"
    rc=$?
  } 2> /dev/null
  set -e
  wait "$wp" || true
  for n in $(cat "$namen" 2> /dev/null); do echo "$n" >> "$v_liste"; done
  if [ -n "${GP_ENTDECKE-}" ]; then echo "p5-$sg: rc=$rc soll=$soll marke=$([ -e "$marke" ] && echo ja) namen=$(cat "$namen" 2>/dev/null) err: $(cat "$err")" >&2; fi
  [ -e "$marke" ] && [ -s "$namen" ] && [ "$rc" -eq "$soll" ] \
    && ! "$docker_echt" container inspect "$pg" > /dev/null 2>&1 \
    && alle_weg $(cat "$namen") "$daten" \
    && in_err "$(cut -d' ' -f1 "$namen")"
}
f_p5_sigterm() { f_p5 TERM "$1" 143; }
f_p5_sigint() { f_p5 INT "$1" 130; }

# pruefe <Verzeichnis> — schreibt die Namen der Fälle, die nicht bestehen; $nur wählt aus.
pruefe() {
  fall p2-grundfall f_p2_grundfall "$1"
  fall p2-ohne-v f_p2_ohne_v "$1"
  fall p3-angelegt f_p3_angelegt "$1"
  fall p3-kein-container f_p3_kein_container "$1"
  fall p4-fremd f_p4_fremd "$1"
  fall p4-rest-rot f_p4_rest_rot "$1"
  fall p4-vorrang f_p4_vorrang "$1"
  fall p5-sigterm f_p5_sigterm "$1"
  fall p5-sigint f_p5_sigint "$1"
}

BAUM=0
# fixture — legt eine Kopie von Funktionsdatei und Runner an; setzt w auf ihr Verzeichnis.
fixture() {
  BAUM=$((BAUM + 1))
  w="$arbeit/baum-$BAUM"
  mkdir -p "$w"
  cp "$hier/integration-aufraeumen.sh" "$hier/run-integration-tests.sh" "$w/"
}
# mutiere <Datei> <alt> <neu> — ersetzt genau eine Stelle (wörtlich) in der Kopie; trifft die
# Ersetzung nichts oder mehr als eine Stelle, bricht die Gegenprobe ab.
mutiere() {
  local f="$w/$1" inhalt rest n
  inhalt="$(cat "$f"; printf x)"; inhalt="${inhalt%x}"
  rest="${inhalt//"$2"/}"
  n=$(((${#inhalt} - ${#rest}) / ${#2}))
  if [ "$n" -ne 1 ]; then echo "integration-volumes-gegenprobe: Mutation trifft $n Stellen in $1: $2" >&2; exit 1; fi
  printf '%s' "${inhalt/"$2"/"$3"}" > "$f"
}
# mut <name> <datei> <alt> <neu> <nur> <erwartet rot>… — bei der Mutation sind genau diese Fälle rot.
mut() {
  local name="$1" datei="$2" alt="$3" neu="$4" soll ist
  if [ -n "${GP_MUT-}" ] && ! [[ "$1" =~ $GP_MUT ]]; then return 0; fi
  nur="$5"
  shift 5
  fixture
  mutiere "$datei" "$alt" "$neu"
  ist="$(pruefe "$w" | sort -u | tr '\n' ' ')"
  soll=""
  if [ $# -gt 0 ]; then soll="$(printf '%s\n' "$@" | sort -u | tr '\n' ' ')"; fi
  if [ -n "${GP_ENTDECKE-}" ]; then echo "$name => ${ist:-nichts}" >&2; fi
  if [ "$ist" != "$soll" ]; then melde "Mutation $name: rot sind [${ist}], erwartet [${soll}]"; fi
  nur=""
}

F_=integration-aufraeumen.sh R_=run-integration-tests.sh

# Die unveränderte Kopie besteht jeden Fall.
nur="${GP_FALL-}"
if [ -z "${GP_MUT-}" ] || [ -n "${GP_BASIS-}" ]; then
  fixture
  ist="$(pruefe "$w" | sort -u | tr '\n' ' ')"
  if [ -n "$ist" ]; then melde "die unveränderte Kopie besteht nicht: ${ist}"; fi
fi

# Punkt 2: ohne -v bleibt ein Rest.
mut p2-ohne-v "$F_" '"$docker" rm -fv "$tests" "$pg"' '"$docker" rm -f "$tests" "$pg"' "p2" p2-grundfall p2-ohne-v
mut p2-ohne-tests "$F_" '"$docker" rm -fv "$tests" "$pg"' '"$docker" rm -fv "$pg"' "p2" p2-grundfall p2-ohne-v
mut p2-ohne-pg "$F_" '"$docker" rm -fv "$tests" "$pg"' '"$docker" rm -fv "$tests"' "p2" p2-grundfall p2-ohne-v
# Punkt 1: Gegenstand der Namen.
mut p1-nur-pg "$F_" 'for c in "$tests" "$pg"; do' 'for c in "$pg"; do' "p3-angelegt" p3-angelegt
mut p1-nur-tests "$F_" 'for c in "$tests" "$pg"; do' 'for c in "$tests"; do' "p3-angelegt" p3-angelegt
mut p1-mit-daten "$F_" 'if [ "$name" != "$daten" ]; then' 'if true; then' "p2-grundfall p3-angelegt" p2-grundfall
# Die Art des Mounts ändert die Namen nicht: Ein Mount ohne Volume trägt keinen Namen (Mutant äquivalent).
mut p1-ohne-typ "$F_" '{{if eq .Type "volume"}}{{.Name}} {{end}}' '{{.Name}} ' "p2-grundfall p3-angelegt"
# Punkt 3: die Namen werden vor dem Entfernen gelesen, in jedem Zustand.
mut p3-nach-entfernen "$F_" $'  namen="$(integration_volumen_lesen)" || return 1\n  integration_entfernen\n' $'  integration_entfernen\n  namen="$(integration_volumen_lesen)" || return 1\n' "p3 p4" p3-angelegt p4-rest-rot
mut p3-leer-fehler "$F_" '2> /dev/null)" || continue' '2> /dev/null)" || return 1' "p2-grundfall p3 p5" p3-kein-container p5-sigint p5-sigterm
# Punkt 4: Rest, Exit-Code, Vorrang, nur eigene Namen.
mut p4-alles-entfernen "$F_" 'for n in $rest; do "$docker" volume rm "$n" > /dev/null 2>&1 || true; done' '"$docker" volume ls -q -f dangling=true | grep -Fx -f "$IVG_SICHER" | xargs -r "$docker" volume rm > /dev/null 2>&1 || true' "p4-fremd" p4-fremd
mut p4-nicht-entfernen "$F_" 'for n in $rest; do "$docker" volume rm "$n" > /dev/null 2>&1 || true; done' ':' "p3-angelegt p4" p3-angelegt
mut p4-ohne-namen "$F_" 'noch vorhanden:$rest" >&2' 'noch vorhanden" >&2' "p3-angelegt" p3-angelegt
mut p4-pruefen-0 "$F_" '  return 1
}

# integration_aufraeumen' '  return 0
}

# integration_aufraeumen' "p3-angelegt p4" p3-angelegt p4-rest-rot
mut p4-code-verworfen "$F_" 'integration_aufraeumen || ende_rest=1' 'integration_aufraeumen || ende_rest=0' "p3-angelegt p4" p3-angelegt p4-rest-rot
mut p4-immer-1 "$F_" 'if [ "$ende_code" -ne 0 ]; then return "$ende_code"; fi' 'if [ "$ende_code" -ne 0 ]; then return 1; fi' "p4" p4-vorrang
# Punkt 5: Ende durch SIGTERM und SIGINT, aus der Falle EXIT, mit Prüfung und Code.
mut p5-exit "$R_" 'trap aufraeumen EXIT # S-EXIT' ':' "p5" p5-sigint p5-sigterm
mut p5-ohne-pruefung "$R_" 'integration_ende "$code" || ende=$?' 'integration_entfernen' "p5" p5-sigint p5-sigterm
mut p5-ausgang-0 "$R_" 'exit "$ende"' 'exit 0' "p5" p5-sigint p5-sigterm
mut p5-term "$R_" $'trap \'trap "" TERM; exit 143\' TERM # S-TERM' $'trap \'kill -KILL $$\' TERM # S-TERM' "p5" p5-sigterm
# Ohne die Falle für TERM läuft die Falle EXIT trotzdem, und bash endet mit 143: Der Mutant ist
# äquivalent, kein Fall ist rot (gemessen mit bash 5.2.21; SPEC-057 Punkt 7 nennt ihn als roten Mutanten).
mut p5-term-entfaellt "$R_" $'trap \'trap "" TERM; exit 143\' TERM # S-TERM' ':' "p5"
mut p5-int "$R_" $'trap \'trap "" INT; exit 130\' INT # S-INT' $'trap \'kill -KILL $$\' INT' "p5" p5-sigint

if [ "$fehler" -ne 0 ]; then exit 1; fi
echo "integration-volumes-gegenprobe: gruen — Punkte 1 bis 8 von SPEC-057: unveränderte Kopie grün, je Mutation genau die benannten Fälle rot"
