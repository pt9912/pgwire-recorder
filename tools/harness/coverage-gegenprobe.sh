#!/usr/bin/env bash
# coverage-gegenprobe — prüft die Schwelle der Anweisungsabdeckung (Prüfschritt
# tools/harness/coverage.sh und Stufe `test` des Dockerfile) gegen den Vertrag
# harness/vertraege/SPEC-056-testabdeckung-schwelle.md (Entscheidung ADR-0044).
# Aufruf: über `make coverage-gegenprobe` (harness/mk/coverage.mk), das DOCKER übergibt;
# auf dem Host mit bash, make, docker, awk, sed, grep und coreutils.
#
# Abdeckung: LH-QA-07/Messung-2-von-3 — Teil 2, Anweisungsabdeckung der Unit-Tests: ein Stand unter der Schwelle wird abgelehnt und einer genau an der Schwelle angenommen (p5-*), Gewicht nach Anweisungen (p3-*), Pakete ohne Test (p4-*), Eingabeformen des Profils und der Schwelle (p6-*), Ausgabe (p7-*), die Stufe test in Kopien des Arbeitsbaums (p1-*, p2-*, p8-*) und die Aufrufarten von make (p9-*)
#
# Der Name jedes Falls trägt die Nummer des Punktes von SPEC-056, den er prüft.
# Fälle ohne Docker: der Prüfschritt an Profilen, die dieses Skript schreibt (p3 bis p7),
# und `make test` mit einer Attrappe als `docker` (p1-wert-*, p1-zeile-*, p9-*) in
# Miniatur-Kopien von Makefile und harness/. Fälle mit Docker (p1-schwelle-kommt-an,
# p1-zeile-fehlt-stufe, p2-paket-ohne-test, p8-*): je eine Kopie des Arbeitsbaums ohne
# .git, `docs/user/` und README.md unter `mktemp -d`, kopiert mit `cp -R` ohne `-p`; jede
# geschriebene Datei bekommt ein `touch`, damit der Build-Kontext sie neu überträgt. Der
# Lauf der Stufe ist `make test` in der Kopie; ein Wrapper um DOCKER ersetzt das Image-Tag
# `pgwire-recorder:test` durch `pgr-cgp-<PID>-<fall>:test`, das am Ende per Namen entfernt
# wird (kein prune, kein Eingriff in ein Objekt mit anderem Namen).
# Die Gegenprobe liest `docs/user/` und README.md des Arbeitsbaums nicht (SPEC-056 Punkt 11);
# dass die Stufe sie in der Kopie nicht braucht, zeigt p8-grundlauf.
#
# GEPRÜFT WIRD — Fälle je Punkt von SPEC-056:
#   (1) Schwelle — p1-wert-aus-datei (die Zeile in der Datei der Kopie, nicht eine feste
#       Zahl), p1-wert-94, p1-zeile-doppelt (`94 94`), p1-zeile-fehlt (leerer Wert),
#       p1-zeile-mitten-im-text (eine Zeile, die nicht mit `**Schwelle:**` beginnt, zählt nicht);
#       in der Stufe: p1-schwelle-kommt-an (100 in der Datei der Kopie wird rot),
#       p1-zeile-fehlt-stufe
#   (2) Messung — p2-paket-ohne-test (ein Paket ohne Test, nur von einem Test eines
#       anderen Pakets ausgeführt, zählt mit 0 %)
#   (3) Gesamtwert — p3-gewicht-gruen, p3-gewicht-rot (Anweisungen, nicht Mittel der
#       Pakete und nicht Blöcke), p3-zaehler-ist-nicht-anzahl
#   (4) Pakete — (a) p4-main-ohne-test, (b) p4-hilfspaket-ohne-test, (c) p4-anzahl-null
#   (5) Vergleich — p5-genau (grün), p5-eins-darunter (rot, abgerundet angezeigt),
#       p5-abgerundet
#   (6) Eingabe — (a) p6-a-fehlt, p6-a-ohne-argumente, p6-a-verzeichnis, p6-a-nicht-lesbar;
#       (b) p6-b-leer; (c) p6-c-ohne-mode, p6-c-count, p6-c-atomic, p6-c-zusatz,
#       p6-c-leere-erste-zeile; (d) p6-d-<form> je Form einer Zeile an drei Stellen
#       (p6-d-pos-*); (e) p6-e-*; (f) p6-f-*; (g) p6-g-*; (h) p6-h-*, dazu die Grenzen
#       1 und 100 als grüne Fälle
#   (7) Ausgabe — p7-gruen-ausgabe, p7-rot-ausgabe (Liste, Reihenfolge, Gesamtzeile
#       zuletzt, Differenz nur bei rot); in der Stufe (awk des Images):
#       p7-gesamtzeile-zuletzt-in-der-stufe
#   (8) Stufe — p8-grundlauf (die unveränderte Kopie ist grün), p8-format, p8-vet-integration,
#       p8-test-rot (jeder bleibt bei erfüllter Schwelle rot)
#   (9) Aufrufarten — p9-ohne-flag, p9-flag-n, p9-flag-q, p9-flag-t, p9-makeflags-n,
#       p9-makeflags-q, p9-makeflags-t (kein Aufruf von docker)
#   (11) Ausnahmeliste — p11-gate-checks, p11-ausnahmeliste
#
# OFFEN — Zusagen ohne Fall:
#   - Punkt 4 (d), (Paket ohne Anweisung steht nicht in der Ausgabe) und der Hinweis auf
#     generierten Code: der Prüfschritt liest nur das Profil, ein Paket ohne Block kommt
#     in ihm nicht vor.
#   - Punkt 7, Ausgabe nur bei laufender Schicht: BuildKit entscheidet über den Cache.
#   - Punkt 8, netzlos und ohne Schreiben in den Arbeitsbaum: die Stufe läuft mit
#     `--network=none`; kein Fall vergleicht den Arbeitsbaum vor und nach einem Lauf.
#   - Aufräumen: dass die Image-Tags `pgr-cgp-<PID>-*` am Ende entfernt werden und die Kopien
#     verschwinden, prüft kein Fall.
#   - Punkt 9, `-i`: wie bei jedem Gate (SPEC-052 Punkt 21), kein eigener Fall.
#   - Punkt 12: Grenzen.
#   - Die Reihenfolge, in der mehrere Gründe einer Eingabe gemeldet werden, ist nicht
#     zugesagt; jeder Fall hat genau einen Grund.
#
# Ausgang: 0, wenn jeder Fall wie erwartet endet, sonst 1.
set -euo pipefail
export LC_ALL=C

wurzel="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
docker_echt="${DOCKER:-docker}"
praefix="pgr-cgp-$$"
parallel=2
arbeit="$(mktemp -d)"
fehler=0
faelle=0
pruef="$wurzel/tools/harness/coverage.sh"
tags="$arbeit/tags"
: > "$tags"

aufraeumen() {
  local status=$? j t
  trap - EXIT
  j="$(jobs -p)"
  if [ -n "$j" ]; then kill $j 2>/dev/null || true; wait 2>/dev/null || true; fi
  while read -r t; do
    [ -z "$t" ] || "$docker_echt" image rm "$t" > /dev/null 2>&1 || true
  done < "$tags"
  rm -rf "${arbeit:?}"
  if [ "$status" -ne 0 ] && [ "$fehler" -eq 0 ]; then
    echo "coverage-gegenprobe: ROT — Abbruch mit Ausgang $status" >&2
  fi
  exit "$status"
}
trap aufraeumen EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

melde() { echo "coverage-gegenprobe: ROT — $*" >&2; fehler=1; }
fall() { faelle=$((faelle + 1)); }

# --- Fälle ohne Docker: der Prüfschritt an Profilen -------------------------
mkdir -p "$arbeit/p"
st=0 aus="" err=""

# pf <name> — schreibt stdin als Profil <arbeit>/p/<name>; setzt $pfd.
pf() { pfd="$arbeit/p/$1"; cat > "$pfd"; }

# lauf <argument>… — fährt den Prüfschritt; setzt st, aus (stdout) und err (stderr).
lauf() {
  st=0
  sh "$pruef" "$@" > "$arbeit/aus" 2> "$arbeit/err" || st=$?
  aus="$(cat "$arbeit/aus")"
  err="$(cat "$arbeit/err")"
}

# gruen <fall> — erwartet Ausgang 0.
gruen() {
  fall
  [ "$st" -eq 0 ] || melde "Fall '$1': erwartet Ausgang 0, Ausgang $st ($err)"
}

# rot <fall> <text> — erwartet Ausgang 1, den Text auf stderr und keine Zeile auf stdout.
rot() {
  fall
  [ "$st" -eq 1 ] || { melde "Fall '$1': erwartet Ausgang 1, Ausgang $st"; return; }
  grep -qF -- "$2" <<< "$err" || melde "Fall '$1': erwartet auf stderr '$2', gefunden '$err'"
  [ -z "$aus" ] || melde "Fall '$1': Ablehnung der Eingabe mit Ausgabe auf stdout: '$aus'"
}

# aus_hat <fall> <text> / aus_hat_nicht <fall> <text> — Text in stdout (fest, ganze Zeichenfolge).
aus_hat() { fall; grep -qF -- "$2" <<< "$aus" || melde "Fall '$1': erwartet in der Ausgabe '$2', gefunden '$aus'"; }
aus_hat_nicht() { fall; if grep -qF -- "$2" <<< "$aus"; then melde "Fall '$1': unerwartet in der Ausgabe '$2'"; fi; }

# aus_ist <fall> <erwartet> — die Ausgabe auf stdout ist genau der Text.
aus_ist() {
  fall
  [ "$aus" = "$2" ] || melde "Fall '$1': erwartete Ausgabe '$2', gefunden '$aus'"
}

B='m/a/x.go:1.1,2.2'   # Block-Anfang eines gültigen Profils

# (3) Gesamtwert nach Anweisungen
pf p3-gewicht-gruen <<EOF
mode: set
m/a/x.go:1.1,2.2 1 0
m/b/x.go:1.1,2.2 99 1
EOF
lauf "$pfd" 90; gruen p3-gewicht-gruen; aus_hat p3-gewicht-gruen 'gesamt 99.00 %'
pf p3-gewicht-rot <<EOF
mode: set
m/a/x.go:1.1,2.2 1 1
m/b/x.go:1.1,2.2 99 0
EOF
lauf "$pfd" 50; st_alt=$st; fall
[ "$st_alt" -eq 1 ] || melde "Fall 'p3-gewicht-rot': erwartet Ausgang 1, Ausgang $st_alt"
aus_hat p3-gewicht-rot 'gesamt 1.00 %'
pf p3-zaehler-ist-nicht-anzahl <<EOF
mode: set
m/a/x.go:1.1,2.2 10 5
m/a/x.go:3.1,4.2 10 0
EOF
lauf "$pfd" 50; gruen p3-zaehler-ist-nicht-anzahl; aus_hat p3-zaehler-ist-nicht-anzahl 'gesamt 50.00 %'

# (4) Pakete
pf p4-main-ohne-test <<EOF
mode: set
m/lib/a.go:1.1,2.2 100 1
m/cmd/app/main.go:1.1,2.2 100 0
EOF
lauf "$pfd" 60; fall
[ "$st" -eq 1 ] || melde "Fall 'p4-main-ohne-test': erwartet Ausgang 1, Ausgang $st"
aus_hat p4-main-ohne-test 'coverage: m/cmd/app 0.00 % (100 Anweisungen)'
aus_hat p4-main-ohne-test 'gesamt 50.00 %'
pf p4-hilfspaket-ohne-test <<EOF
mode: set
m/lib/a.go:1.1,2.2 100 1
m/testpki/p.go:1.1,2.2 100 0
EOF
lauf "$pfd" 60; fall
[ "$st" -eq 1 ] || melde "Fall 'p4-hilfspaket-ohne-test': erwartet Ausgang 1, Ausgang $st"
aus_hat p4-hilfspaket-ohne-test 'coverage: m/testpki 0.00 % (100 Anweisungen)'
pf p4-anzahl-null <<EOF
mode: set
m/lib/a.go:1.1,2.2 100 1
m/leer/doc.go:1.1,2.2 0 0
m/lib/a.go:3.1,4.2 0 0
EOF
lauf "$pfd" 100; gruen p4-anzahl-null
aus_hat_nicht p4-anzahl-null 'm/leer'
aus_hat p4-anzahl-null 'coverage: m/lib 100.00 % (100 Anweisungen)'

# (5) Vergleich und Rundung
pf p5-genau <<EOF
mode: set
$B 9400 1
m/a/x.go:3.1,4.2 600 0
EOF
lauf "$pfd" 94; gruen p5-genau; aus_hat p5-genau 'gesamt 94.00 %'
pf p5-eins-darunter <<EOF
mode: set
$B 9399 1
m/a/x.go:3.1,4.2 601 0
EOF
lauf "$pfd" 94; fall
[ "$st" -eq 1 ] || melde "Fall 'p5-eins-darunter': erwartet Ausgang 1, Ausgang $st"
aus_hat p5-eins-darunter 'gesamt 93.99 %'
aus_hat p5-eins-darunter 'Differenz 0.01 Prozentpunkte'
pf p5-abgerundet <<EOF
mode: set
$B 2 1
m/a/x.go:3.1,4.2 1 0
EOF
lauf "$pfd" 60; gruen p5-abgerundet
aus_hat p5-abgerundet 'coverage: m/a 66.66 % (3 Anweisungen)'
aus_hat p5-abgerundet 'gesamt 66.66 %'
aus_hat_nicht p5-abgerundet '66.67'

# (6) Eingabe des Prüfschritts
G="$arbeit/p/gueltig"
printf 'mode: set\n%s 5 1\n' "$B" > "$G"
lauf "$arbeit/p/fehlt" 94; rot p6-a-fehlt 'Profil fehlt oder ist nicht lesbar'
lauf; rot p6-a-ohne-argumente 'Profil fehlt oder ist nicht lesbar'
mkdir "$arbeit/p/verzeichnis"; lauf "$arbeit/p/verzeichnis" 94; rot p6-a-verzeichnis 'Profil fehlt oder ist nicht lesbar'
cp "$G" "$arbeit/p/gesperrt"; chmod 000 "$arbeit/p/gesperrt"
if [ -r "$arbeit/p/gesperrt" ]; then
  echo "coverage-gegenprobe: p6-a-nicht-lesbar übersprungen (der Aufrufer liest auch gesperrte Dateien)" >&2
else
  lauf "$arbeit/p/gesperrt" 94; rot p6-a-nicht-lesbar 'Profil fehlt oder ist nicht lesbar'
fi
: > "$arbeit/p/leer"; lauf "$arbeit/p/leer" 94; rot p6-b-leer 'Profil ist leer'
pf p6-c-ohne-mode <<EOF
$B 5 1
EOF
lauf "$pfd" 94; rot p6-c-ohne-mode "'mode: set'"
pf p6-c-count <<EOF
mode: count
$B 5 1
EOF
lauf "$pfd" 94; rot p6-c-count "'mode: set'"
pf p6-c-atomic <<EOF
mode: atomic
$B 5 1
EOF
lauf "$pfd" 94; rot p6-c-atomic "'mode: set'"
pf p6-c-zusatz <<EOF
mode: set x
$B 5 1
EOF
lauf "$pfd" 94; rot p6-c-zusatz "'mode: set'"
printf '\nmode: set\n%s 5 1\n' "$B" > "$arbeit/p/leere-erste"
lauf "$arbeit/p/leere-erste" 94; rot p6-c-leere-erste-zeile "'mode: set'"

# (d): je Form eine Zeile, die nicht <Datei>:<Bereich> <Anzahl> <Zähler> hat; in der Mitte
# eines gültigen Profils (p6-d-<form>) und je an erster und letzter Stelle (p6-d-pos-*).
d_fall() { # name zeile
  { printf 'mode: set\n%s 5 1\n' "$B"; printf '%s\n' "$2"; printf 'm/b/x.go:1.1,2.2 5 1\n'; } > "$arbeit/p/$1"
  lauf "$arbeit/p/$1" 94; rot "$1" 'nicht die Form'
}
d_fall p6-d-ohne-doppelpunkt 'm/a/x.go.1.1,2.2 5 1'
d_fall p6-d-datei-leer ':1.1,2.2 5 1'
d_fall p6-d-bereich-leer 'm/a/x.go: 5 1'
d_fall p6-d-ein-feld 'm/a/x.go:1.1,2.2'
d_fall p6-d-zwei-felder 'm/a/x.go:1.1,2.2 5'
d_fall p6-d-vier-felder 'm/a/x.go:1.1,2.2 5 1 1'
d_fall p6-d-anzahl-text 'm/a/x.go:1.1,2.2 x 1'
d_fall p6-d-anzahl-dezimal 'm/a/x.go:1.1,2.2 1.5 1'
d_fall p6-d-anzahl-negativ 'm/a/x.go:1.1,2.2 -1 1'
d_fall p6-d-zaehler-text 'm/a/x.go:1.1,2.2 5 x'
d_fall p6-d-zaehler-dezimal 'm/a/x.go:1.1,2.2 5 1.5'
d_fall p6-d-zaehler-negativ 'm/a/x.go:1.1,2.2 5 -1'
d_fall p6-d-doppeltes-leerzeichen 'm/a/x.go:1.1,2.2  5 1'
d_fall p6-d-leerzeichen-am-ende 'm/a/y.go:1.1,2.2 5 1 '
d_fall p6-d-leerzeichen-am-anfang ' m/a/y.go:1.1,2.2 5 1'
d_fall p6-d-leerzeile ''
d_fall p6-d-tab "$(printf 'm/a/y.go:1.1,2.2\t5 1')"
printf 'mode: set\nm/a/x.go:1.1,2.2 5\n%s 5 1\n' "$B" > "$arbeit/p/d-pos-erste"
lauf "$arbeit/p/d-pos-erste" 94; rot p6-d-pos-erste 'nicht die Form'
printf 'mode: set\n%s 5 1\nm/a/x.go:3.1,4.2 5\n' "$B" > "$arbeit/p/d-pos-letzte"
lauf "$arbeit/p/d-pos-letzte" 94; rot p6-d-pos-letzte 'nicht die Form'

# (e): ein Block (Datei und Bereich) zweimal; verschiedene Bereiche oder Dateien sind keiner.
e_fall() { # name zeile-zwei
  printf 'mode: set\n%s 5 1\n%s\n' "$B" "$2" > "$arbeit/p/$1"
  lauf "$arbeit/p/$1" 94
}
e_fall p6-e-doppelt "$B 5 1";                 rot p6-e-doppelt 'zweimal'
e_fall p6-e-anderer-zaehler "$B 5 0";         rot p6-e-anderer-zaehler 'zweimal'
e_fall p6-e-andere-anzahl "$B 7 1";           rot p6-e-andere-anzahl 'zweimal'
printf 'mode: set\n%s 5 1\nm/b/x.go:1.1,2.2 5 1\n%s 5 1\n' "$B" "$B" > "$arbeit/p/e-getrennt"
lauf "$arbeit/p/e-getrennt" 94; rot p6-e-getrennt 'zweimal'
e_fall p6-e-anderer-bereich "m/a/x.go:3.1,4.2 5 1"; gruen p6-e-anderer-bereich
e_fall p6-e-andere-datei "m/b/x.go:1.1,2.2 5 1";   gruen p6-e-andere-datei

# (f): Summe der Anweisungen 0
printf 'mode: set\n' > "$arbeit/p/f-nur-kopf"; lauf "$arbeit/p/f-nur-kopf" 94; rot p6-f-nur-kopf 'Summe der Anweisungen'
printf 'mode: set\n%s 0 0\nm/b/x.go:1.1,2.2 0 1\n' "$B" > "$arbeit/p/f-null"
lauf "$arbeit/p/f-null" 94; rot p6-f-nur-null-bloecke 'Summe der Anweisungen'

# (g), (h): Schwelle
lauf "$G"; rot p6-g-fehlt 'Schwelle fehlt'
lauf "$G" ''; rot p6-g-leer 'Schwelle fehlt'
for w in '94.5' '094' '0' '101' '94 94' '-1' '+94' 'abc' ' 94' '94 ' '1e2' '00' '1000'; do
  lauf "$G" "$w"; rot "p6-h-[$w]" 'keine ganze Zahl'
done
lauf "$G" "$(printf '94\n94')"; rot p6-h-zeilenumbruch 'keine ganze Zahl'
printf 'mode: set\n%s 5 1\n' "$B" > "$arbeit/p/voll"
lauf "$arbeit/p/voll" 1; gruen p6-h-grenze-1
lauf "$arbeit/p/voll" 9; gruen p6-h-grenze-9
lauf "$arbeit/p/voll" 10; gruen p6-h-grenze-10
lauf "$arbeit/p/voll" 100; gruen p6-h-grenze-100
printf 'mode: set\n%s 5 0\n' "$B" > "$arbeit/p/leer-ausgefuehrt"
lauf "$arbeit/p/leer-ausgefuehrt" 1; fall
[ "$st" -eq 1 ] || melde "Fall 'p6-h-grenze-1-rot': erwartet Ausgang 1 bei 0 % und Schwelle 1, Ausgang $st"

# (7) Ausgabe: Liste, Reihenfolge, Gesamtzeile zuletzt, Differenz nur bei rot
pf p7 <<EOF
mode: set
m/z/z.go:1.1,2.2 4 1
m/a/b/b.go:1.1,2.2 3 0
m/a/a.go:1.1,2.2 2 1
m/a/a.go:3.1,4.2 2 0
m/a-b/c.go:1.1,2.2 3 1
m/a-b/c.go:3.1,4.2 0 0
m/B/x.go:1.1,2.2 5 1
EOF
paketzeilen='coverage: m/B 100.00 % (5 Anweisungen)
coverage: m/a 50.00 % (4 Anweisungen)
coverage: m/a-b 100.00 % (3 Anweisungen)
coverage: m/a/b 0.00 % (3 Anweisungen)
coverage: m/z 100.00 % (4 Anweisungen)'
lauf "$pfd" 70; gruen p7-gruen-ausgabe
aus_ist p7-gruen-ausgabe "$paketzeilen
coverage: gesamt 73.68 % (14 von 19 Anweisungen), Schwelle 70 %"
lauf "$pfd" 80; fall
[ "$st" -eq 1 ] || melde "Fall 'p7-rot-ausgabe': erwartet Ausgang 1, Ausgang $st"
aus_ist p7-rot-ausgabe "$paketzeilen
coverage: gesamt 73.68 % (14 von 19 Anweisungen), Schwelle 80 %, Differenz 6.32 Prozentpunkte, unter der Schwelle"

# --- Fälle mit make und einer Attrappe als docker (ohne Docker) ---------------
attrappe="$arbeit/attrappe"
cat > "$attrappe" <<'ATTRAPPE'
#!/bin/sh
printf '%s\n' "$@" >> "${CGP_LOG:?}"
ATTRAPPE
chmod +x "$attrappe"
SPEC=harness/vertraege/SPEC-056-testabdeckung-schwelle.md

# miniatur <name> — Makefile und harness/ des Arbeitsbaums; setzt m (Wurzel der Kopie).
miniatur() {
  m="$arbeit/mini-$1"
  mkdir -p "$m"
  cp -R "$wurzel/Makefile" "$wurzel"/*.mk "$wurzel/harness" "$m/"
  [ -f "$m/$SPEC" ] || melde "Fall 'p1': $SPEC fehlt im Arbeitsbaum"
}

# make_attrappe <kopie> <variablen…> -- <make-argumente…> — fährt make in der Kopie mit der
# Attrappe; die Umgebung von make des Aufrufers (MAKEFLAGS, MAKELEVEL) gilt nicht.
make_attrappe() {
  local k="$1"; shift
  local umgebung=()
  while [ "$1" != -- ]; do umgebung+=("$1"); shift; done
  shift
  rm -f "$arbeit/cgp.log"
  env -u MAKEFLAGS -u MFLAGS -u MAKELEVEL "${umgebung[@]}" CGP_LOG="$arbeit/cgp.log" \
    make -C "$k" --no-print-directory "$@" DOCKER="$attrappe" > "$arbeit/make.log" 2>&1 || true
}

# log_hat <fall> <zeile> — die Attrappe hat <zeile> als ganzes Argument erhalten.
log_hat() {
  fall
  [ -f "$arbeit/cgp.log" ] && grep -qxF -- "$2" "$arbeit/cgp.log" \
    || melde "Fall '$1': erwartet das Argument '$2', Aufrufe: $(cat "$arbeit/cgp.log" 2>/dev/null | tr '\n' ' ')"
}
# kein_aufruf <fall> — die Attrappe wurde nicht gerufen.
kein_aufruf() {
  fall
  if [ -f "$arbeit/cgp.log" ]; then melde "Fall '$1': docker wurde gerufen: $(tr '\n' ' ' < "$arbeit/cgp.log")"; fi
}

miniatur grund
make_attrappe "$m" -- test
log_hat p1-wert-94 'COVERAGE_SCHWELLE=94'
log_hat p9-ohne-flag 'test'
log_hat p9-ohne-flag '--build-arg'
log_hat p9-ohne-flag 'COVERAGE_SCHWELLE=94'
make_attrappe "$m" -- -n test; kein_aufruf p9-flag-n
make_attrappe "$m" -- -q test; kein_aufruf p9-flag-q
make_attrappe "$m" -- -t test; kein_aufruf p9-flag-t
make_attrappe "$m" MAKEFLAGS=n -- test; kein_aufruf p9-makeflags-n
make_attrappe "$m" MAKEFLAGS=q -- test; kein_aufruf p9-makeflags-q
make_attrappe "$m" MAKEFLAGS=t -- test; kein_aufruf p9-makeflags-t

miniatur wert
sed -i 's/^\( *\)\*\*Schwelle:\*\* `94`/\1**Schwelle:** `77`/' "$m/$SPEC"; touch "$m/$SPEC"
make_attrappe "$m" -- test; log_hat p1-wert-aus-datei 'COVERAGE_SCHWELLE=77'
miniatur doppelt
awk '{ print } /^ *\*\*Schwelle:\*\* `/ { print }' "$m/$SPEC" > "$m/spec.neu" && mv "$m/spec.neu" "$m/$SPEC"
make_attrappe "$m" -- test; log_hat p1-zeile-doppelt 'COVERAGE_SCHWELLE=94 94'
miniatur mitte
printf 'Im Text steht **Schwelle:** `77` mitten in der Zeile.\n- **Schwelle:** `66` als Listenpunkt.\n' >> "$m/$SPEC"
make_attrappe "$m" -- test; log_hat p1-zeile-mitten-im-text 'COVERAGE_SCHWELLE=94'
miniatur fehlt
sed -i '/^ *\*\*Schwelle:\*\* `/d' "$m/$SPEC"
make_attrappe "$m" -- test; log_hat p1-zeile-fehlt 'COVERAGE_SCHWELLE='

# (11) Ausnahmeliste und Anschluss an GATE_CHECKS: Wörter der Make-Variablen aus der Datenbank von make.
miniatur liste
variable() {
  env -u MAKEFLAGS -u MFLAGS -u MAKELEVEL make -s --no-print-directory -C "$m" -pq help 2>/dev/null \
    | awk -v n="$1" '$1 == n && ($2 == ":=" || $2 == "=") { for (i = 3; i <= NF; i++) print $i }' || true
}
fall; variable GATE_CHECKS | grep -x coverage-gegenprobe > /dev/null || melde "Fall 'p11-gate-checks': coverage-gegenprobe steht nicht in GATE_CHECKS"
fall; variable GATE_SCHNELL_AUSNAHME | grep -x coverage-gegenprobe > /dev/null || melde "Fall 'p11-ausnahmeliste': coverage-gegenprobe steht nicht in GATE_SCHNELL_AUSNAHME"

# --- Fälle mit Docker: die Stufe test in Kopien -------------------------------
dockerwrap="$arbeit/docker"
cat > "$dockerwrap" <<'WRAP'
#!/usr/bin/env bash
a=()
for x in "$@"; do
  if [ "$x" = pgwire-recorder:test ]; then x="$CGP_TAG"; fi
  a+=("$x")
done
exec "$CGP_DOCKER" "${a[@]}"
WRAP
chmod +x "$dockerwrap"

# kopie <name> — Kopie des Arbeitsbaums ohne .git, docs/user und README.md; k ist ihre Wurzel.
kopie() {
  k="$arbeit/$1"
  mkdir -p "$k"
  local e
  for e in "$wurzel"/* "$wurzel"/.[!.]*; do
    [ -e "$e" ] || continue
    case "${e##*/}" in .git | README.md) continue ;; esac
    cp -R "$e" "$k/"
  done
  rm -rf "$k/docs/user"
}

# datei <pfad> — schreibt stdin nach <pfad> in der Kopie k.
datei() {
  mkdir -p "$(dirname "$k/$1")"
  cat > "$k/$1"
  touch "$k/$1"
}

# starte <name> — fährt `make test` in der Kopie <name> im Hintergrund; Ausgabe nach
# <name>.log, Ausgang nach <name>.exit.
starte() {
  local name="$1"
  echo "$praefix-$name:test" >> "$tags"
  while [ "$(jobs -rp | wc -l)" -ge "$parallel" ]; do wait -n || true; done
  (
    cd "$arbeit/$name"
    set +e
    CGP_TAG="$praefix-$name:test" CGP_DOCKER="$docker_echt" \
      env -u MAKEFLAGS -u MFLAGS -u MAKELEVEL make --no-print-directory test DOCKER="$dockerwrap" \
      > "$arbeit/$name.log" 2>&1 &
    kind=$!
    trap 'kill "$kind" 2>/dev/null' TERM
    wait "$kind"
    { echo $? > "$arbeit/$name.exit"; } 2>/dev/null
  ) &
}

# docker_rot <fall> <text>… — erwartet Ausgang ungleich 0 und jeden Text im Log.
docker_rot() {
  local f="$1" t; shift
  fall
  [ "$(cat "$arbeit/$f.exit")" -ne 0 ] || melde "Fall '$f': erwartet Ausgang ungleich 0 der Stufe, Lauf endete mit 0"
  for t in "$@"; do
    grep -qF -- "$t" "$arbeit/$f.log" || melde "Fall '$f': erwartet im Lauf '$t'"
  done
}
# docker_nicht <fall> <text> — der Text steht nicht im Log.
docker_nicht() {
  fall
  if grep -qF -- "$2" "$arbeit/$1.log"; then melde "Fall '$1': unerwartet im Lauf '$2'"; fi
}

kopie p8-grundlauf; starte p8-grundlauf

kopie p1-schwelle-kommt-an
sed -i 's/^\( *\)\*\*Schwelle:\*\* `94`/\1**Schwelle:** `100`/' "$k/$SPEC"
starte p1-schwelle-kommt-an

kopie p1-zeile-fehlt-stufe
sed -i '/^ *\*\*Schwelle:\*\* `/d' "$k/$SPEC"
starte p1-zeile-fehlt-stufe

# p2: 301 Anweisungen in einem Paket ohne Test; ein Test eines anderen Pakets führt sie aus.
kopie p2-paket-ohne-test
{
  printf 'package cgpviel\n\nfunc F(a int) int {\n'
  for i in $(seq 1 300); do printf '\ta += %d\n' "$i"; done
  printf '\treturn a\n}\n'
} | datei internal/cgpviel/viel.go
datei internal/testuhr/cgp_test.go <<'GO'
package testuhr_test

import (
	"testing"

	"github.com/pt9912/pgwire-recorder/internal/cgpviel"
)

func TestCgpViel(t *testing.T) {
	if cgpviel.F(0) == 0 {
		t.Fatal("F(0) ist 0")
	}
}
GO
starte p2-paket-ohne-test

kopie p8-format
datei internal/cgpformat/f.go <<'GO'
package cgpformat

func F( ) {
}
GO
starte p8-format

kopie p8-vet-integration
datei internal/cgpvet/v.go <<'GO'
//go:build integration

package cgpvet

import "fmt"

func F() string {
	return fmt.Sprintf("%d", "x")
}
GO
starte p8-vet-integration

kopie p8-test-rot
datei internal/cgprot/doc.go <<'GO'
package cgprot
GO
datei internal/cgprot/rot_test.go <<'GO'
package cgprot

import "testing"

func TestRot(t *testing.T) {
	t.Fatal("rot")
}
GO
starte p8-test-rot

wait || true

fall
[ "$(cat "$arbeit/p8-grundlauf.exit")" -eq 0 ] \
  || melde "Fall 'p8-grundlauf': erwartet Ausgang 0 der Stufe, Lauf endete mit $(cat "$arbeit/p8-grundlauf.exit")"

docker_rot p1-schwelle-kommt-an 'Schwelle 100 %' 'unter der Schwelle'
fall
grep -E 'coverage: (github\.com|gesamt)' "$arbeit/p1-schwelle-kommt-an.log" | tail -n 1 | grep -F 'coverage: gesamt' > /dev/null \
  || melde "Fall 'p7-gesamtzeile-zuletzt-in-der-stufe': die letzte Zeile des Prüfschritts im Lauf ist nicht die Gesamtzeile"
docker_rot p1-zeile-fehlt-stufe 'coverage: Schwelle fehlt'
docker_rot p2-paket-ohne-test 'coverage: github.com/pt9912/pgwire-recorder/internal/cgpviel 0.00 % (301 Anweisungen)' 'unter der Schwelle'
docker_rot p8-format 'nicht gofmt-formatiert' 'internal/cgpformat/f.go'
docker_nicht p8-format 'coverage: gesamt'
docker_rot p8-vet-integration 'format %d has arg'
docker_nicht p8-vet-integration 'coverage: gesamt'
docker_rot p8-test-rot '--- FAIL: TestRot'
docker_nicht p8-test-rot 'coverage: gesamt'

if [ "$fehler" -ne 0 ]; then
  echo "coverage-gegenprobe: ROT ($faelle Prüfungen)" >&2
  exit 1
fi
echo "coverage-gegenprobe: $faelle Prüfungen wie erwartet"
