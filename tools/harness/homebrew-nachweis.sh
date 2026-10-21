#!/usr/bin/env bash
# homebrew-nachweis — Nachweis der Installation aus dem Tap (Abnahmeszenario 11).
# Vertrag: harness/vertraege/SPEC-067-release.md Punkt 20; was aus einem Nachweis an Zusage folgt und was
# fuer jede Kombination ohne Nachweis Grenze bleibt, steht in spec/spezifikation.md (LH-FA-19.a).
# Aufruf: bash tools/harness/homebrew-nachweis.sh --schritte
# Unter Linux faehrt `make homebrew-nachweis` genau diesen Aufruf in einem Container des per Digest
# gepinnten Image homebrew/brew (harness/mk/homebrew-nachweis.mk). `--schritte` AM HOST ist allein die
# Handlung des Betreibers auf seinem eigenen Rechner (`brew` am Host ist ein Paketmanager,
# AGENTS.md §3.1); ein Agentenlauf faehrt sie nicht, sein Weg ist `make homebrew-nachweis` im
# Container. Das ist eine Regel an den Leser, kein Mechanismus: Der Waechter
# .claude/hooks/pretooluse-command-guard.sh blockt ein Segment, dessen Kopf ein Paketmanager ist, liest
# aber den Inhalt einer Datei nicht und verhindert diesen Aufruf darum nicht (SPEC-067 Punkt 20 (5)).
# Auf macOS laeuft bash 3.2: ohne mapfile, ohne assoziative Arrays, ohne ${v^^}.
# Die Schritte sind auf beiden Plattformen dieselben, damit die Belege vergleichbar sind: (a) brew tap,
# (b) brew install, (c) pgwire-recorder version, (d) brew test, (e) der Replay-Lauf aus Punkt 20 (3).
# `brew audit` gehoert nicht dazu (Punkt 20 (2): die Zeile `version` ist dort ein bekannter Befund und
# sperrt die Installation nicht).
# Das Skript gibt jeden Befehl vor seiner Ausfuehrung auf stdout aus; es wertet keine Ausgabe aus,
# vergleicht keinen Text und legt keine Marke an. Die eine Ausnahme ist das Warten auf den Start des
# Replay-Laufs: Es prueft die GROESSE der Fehlerausgabe (`[ -s … ]`), nie ihren Inhalt, im Abstand von
# einer Sekunde, Grenze 60 Sekunden; laeuft sie ab, ist das Warten der gescheiterte Schritt und das
# Skript endet mit Exit-Code 7, auch wenn der Lauf selbst 0 liefert (SPEC-067 Punkt 20 (2), (3)). Das
# Warten auf das Ende NACH dem SIGTERM hat keine Grenze: Dass der Lauf nach dem Signal endet, ist
# Gegenstand des Nachweises, und ein Haengen ist selbst der Befund.
# Sonst endet das Skript mit dem Exit-Code des ersten Schrittes, der scheitert, sonst mit 0. Was aus
# der Ausgabe zaehlt, liest der Mensch: der Beleg steht in §7 des Slice-Plans, der ihn fuehrt
# (Punkt 20 (6)).
# Der Kopf nennt Datum, `uname -s` und `uname -m` — die Angaben, aus denen der Beleg je Kombination
# aus Betriebssystem und Architektur besteht (Punkt 20 (6)); unter Linux nennt das Fragment dazu Tag
# und Digest des Image.
# Die letzte Zeile der Ausgabe nennt den Exit-Code des Skripts selbst:
# `homebrew-nachweis: Skript endet mit Exit-Code <n>`. Wer die Ausgabe liest, erkennt daran, wessen
# Code der Lauf meldet — fehlt die Zeile, kam das Skript nicht zum Ende, und der Code stammt vom
# Treiber davor (unter Linux von der Docker-CLI, harness/mk/homebrew-nachweis.mk).
set -u

tap='pt9912/pgwire-recorder'
formel='pt9912/pgwire-recorder/pgwire-recorder'
adresse='127.0.0.1:55432'

tmp=''
replay=''

# aufraeumen — beendet einen laufenden Replay-Lauf und entfernt das Temp-Verzeichnis; mehrfach aufrufbar.
aufraeumen() {
  trap '' INT TERM
  if [ -n "$replay" ]; then
    kill -TERM "$replay" 2> /dev/null
    wait "$replay" 2> /dev/null
    replay=''
  fi
  if [ -n "$tmp" ]; then
    rm -rf "$tmp"
    tmp=''
  fi
}
trap 'aufraeumen; exit 130' INT
trap 'aufraeumen; exit 143' TERM
# Die EXIT-Falle raeumt auf und nennt den Exit-Code, mit dem dieses Skript endet; sie aendert ihn nicht.
trap 'c=$?; aufraeumen; echo "homebrew-nachweis: Skript endet mit Exit-Code $c"' EXIT

# schritt <Befehl>… — nennt den Befehl auf stdout und fuehrt ihn aus.
schritt() {
  echo "+ $*"
  "$@"
}

if [ "$#" -ne 1 ] || [ "$1" != '--schritte' ]; then
  echo "homebrew-nachweis: genau ein Argument verlangt, --schritte" >&2
  exit 2
fi

schritt date -u '+%Y-%m-%dT%H:%M:%SZ' || exit $?
schritt uname -s || exit $?
schritt uname -m || exit $?

schritt brew tap "$tap" || exit $?
schritt brew install "$formel" || exit $?
schritt pgwire-recorder version || exit $?
schritt brew test pgwire-recorder || exit $?

tmp="$(mktemp -d)" || exit $?
aufzeichnung="$tmp/aufzeichnung.yaml"
fehlerausgabe="$tmp/replay.stderr"
cat > "$aufzeichnung" <<'YAML'
format: pgwire-recorder
version: 2
sessions:
  - id: 1
    interactions:
      - sequence: 1
        request: {type: query, sql: "SELECT 1;"}
        responses:
          - {type: command_complete, tag: SELECT 1}
          - {type: ready_for_query, tx_status: I}
YAML

echo "+ pgwire-recorder replay --listen $adresse --input $aufzeichnung (Fehlerausgabe in eine Datei)"
pgwire-recorder replay --listen "$adresse" --input "$aufzeichnung" 2> "$fehlerausgabe" &
replay=$!

# Warten auf den Start: allein die Groesse der Fehlerausgabe, nie ihr Inhalt; Abstand eine Sekunde,
# hoechstens 60 Durchlaeufe (Grenze 60 Sekunden, fest).
echo "+ warten, bis der Replay-Lauf seine erste Zeile geschrieben hat (Groesse der Fehlerausgabe, Abstand 1 s, Grenze 60 s)"
grenze_ab=0
i=0
while [ ! -s "$fehlerausgabe" ]; do
  if [ "$i" -ge 60 ]; then
    grenze_ab=1
    echo "homebrew-nachweis: die Grenze von 60 s ist abgelaufen, die Fehlerausgabe des Replay-Laufs ist leer geblieben"
    break
  fi
  sleep 1
  i=$((i + 1))
done

echo "+ kill -TERM <Replay-Lauf>"
kill -TERM "$replay" 2> /dev/null
wait "$replay"
ende=$?
replay=''
echo "--- Fehlerausgabe des Replay-Laufs ---"
cat "$fehlerausgabe"
echo "--- Exit-Code des Replay-Laufs: $ende"
if [ "$grenze_ab" = 1 ]; then
  echo "homebrew-nachweis: das Warten auf den Start ist der gescheiterte Schritt; der Nachweis gilt als nicht erbracht (Exit-Code 7)"
  exit 7
fi
exit "$ende"
