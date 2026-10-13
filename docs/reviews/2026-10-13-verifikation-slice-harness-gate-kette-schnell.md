# Verifikation: slice-harness-gate-kette-schnell — 2026-10-13

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-13-review-slice-harness-gate-kette-schnell.md`).

**Gegenstand:** Slice-Plan `slice-harness-gate-kette-schnell` (Kopf, §1 bis §8, mit den Entscheidungen des Architect vom 2026-10-12/13) gegen `d075d3f..8708f21`. Commits: `cb37a34` (`SPEC-052`), `f633694` (Code), `a203246` (Belege), `e09b4f0`/`bee6785` (Review-Report), `b5cb7b3` (Entscheidungen des Architect), `0e8fa93`/`14c72f9` (Nacharbeit), `6e22c4b` (`-n`/`-q`/`-t` mit `-i`), `68bfa16`/`da865e9`/`8708f21` (Fall `p6-ni`, Belege, Stand der Gates).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-13

**Eingang:** Plan ganz (DoD, §1, §3, §4, §6, §7 *Belege des Implementers*, §8), `SPEC-052` (Teil 1 und 2), [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) (Accepted), Diff, Review-Report (Findings-Tabelle). Den Bericht des Implementers habe ich nicht gelesen.

## 1. Selbst gefahrene Sensoren

Ergebnis vorweg: kein HIGH, kein MEDIUM; die Zusagen der Kette halten gegen meine Proben und gegen elf eigene Mutationen. Der Arbeitsbaum war danach unverändert (`git status --short` leer).

| Nr. | Sensor / Probe | Ergebnis |
|---|---|---|
| 1 | `make gates-schnell-gegenprobe` | grün, 18,9 s (§7: rund 17 s); Meldung „Punkte 1 bis 7 (Punkt 5 mit Ende des Laufs, Punkt 6 mit Make-Flags) …“ |
| 2 | `make gates-schnell` (Cache warm) | Exit 0, 7,2 s (§7: 7,2 bis 7,8 s); Stand `.harness/state/gates-schnell-passed.diffsha` gleich der Ausgabe von `tools/harness/working-tree-hash.sh` (`d5856278…408f`) |
| 3 | `make -n gates-schnell`, `MAKEFLAGS=n make gates-schnell`, je ohne vorhandenen Stand | Exit 0, kein Stand entstanden (F-694 behoben) |
| 4 | `make -n gates-schnell` und `MAKEFLAGS=ni make gates-schnell` bei vorhandenem Stand | Stand byteweise unberührt (`cmp`); gilt auch für `make -q` (Exit 1, Stand da) und `make -t` (Exit 0, Stand da, keine Datei berührt) |
| 5 | Kopie im Scratchpad (Makefile, Fragmente, Skripte, ein Fragment mit rotem und grünem Attrappen-Gate): `make gates-schnell` mit rotem Gate bei vorhandenem Stand; `-i`; `--ignore-errors`; `-ki`; `-j4 -i`; `--keep-going`; `--dry-run` | rot: Exit 2, Stand weg; `-i`, `--ignore-errors`, `-ki`, `-j4 -i`: Exit 0, Stand weg; `--keep-going`: Exit 2, Stand weg; `--dry-run`: Exit 0, Stand bleibt. Grüner Lauf in der Kopie: Stand gleich Hash nach dem Lauf. |
| 6 | Roter Lauf am Arbeitsbaum (F-701): `// probe` ohne Leerzeile an `cmd/pgwire-recorder/main.go` angehängt, `make gates-schnell` | Exit 2 nach 3,8 s, Stufe `test` meldet „nicht gofmt-formatiert:“ (Log Zeile 102), kein Stand zurückgeblieben; Datei danach `git checkout` |
| 7 | Grüner Lauf mit geänderter Go-Datei: `\n// probe\n` angehängt, `make gates-schnell` | Exit 0 nach 97,8 s (§7: 102 s); Stand gleich Hash; Datei zurückgesetzt, Kette danach erneut grün |
| 8 | `make help` | führt `gates-schnell` und `gates-schnell-gegenprobe`; kein `stop-hook-gegenprobe` (Teil 2 von `SPEC-052` damit nicht zugesagt) |
| 9 | `make -n gates` | enthält `gates-schnell-gegenprobe`, nicht `gates-schnell` (ein Treffer: `bash tools/harness/gates-schnell-gegenprobe.sh`); `Makefile`, `harness/mk/enforce.mk`, `record-gates.sh` unverändert |
| 10 | Ausnahmeliste gegen `GATE_CHECKS` | sieben Einträge (`abdeckung-gegenprobe a-check-negativ gates-schnell-gegenprobe hook-gegenprobe kopf-check-gegenprobe lint-gegenprobe test-integration`); decken sich mit den Namen auf `-gegenprobe` plus `a-check-negativ` und `test-integration` in `GATE_CHECKS`; Punkt 13, `harness/README.md` und Gegenprobe zählen gleich |
| 11 | `make docs-check` | 674 Dateien, 0 Befunde |
| 12 | `git diff d075d3f..HEAD` für `.claude/`, `tools/harness/record-gates.sh`, `tools/harness/working-tree-hash.sh`, `Makefile`, `harness/mk/enforce.mk`, `.golangci.yml`, `.a-check.yml`, `.d-check.yml` | leer; Hook und Nachweis-Skripte unverändert, keine Gate-Lockerung (`AGENTS.md` §3.6): `make gates` fährt ein Gate mehr (die Gegenprobe), keines weniger |
| 13 | Größe `git diff --shortstat d075d3f HEAD -- . ':!docs'` | 563 Zeilen eingefügt, 1 entfernt (`SPEC-052` 157, Gegenprobe 338, Fragment 23, Skript 42, Doku 3); wie in §7 berichtet |

### Mutationen (eigene Kopie im Scratchpad, `git archive HEAD`, je Mutant `bash tools/harness/gates-schnell-gegenprobe.sh`, danach `git checkout`)

Die Gegenprobe endete bei jeder Mutation mit Exit 1. Die Kopie `ok` nennt die roten Fälle:

| Nr. | Zusage (`SPEC-052`) | Mutation | roter Fall der Kopie `ok` |
|---|---|---|---|
| M1 | Punkt 2: kein langsames Gate | `lint-gegenprobe` aus `GATE_SCHNELL_AUSNAHME` genommen | `p2-langsam`, `p4-name`; Arbeitsbaum-Prüfung rot |
| M2 | Punkt 4: Gate `x-gegenprobe` ohne Eintrag (Fragment `zz.mk`) | | Arbeitsbaum: `p4-name` |
| M3 | Punkt 4: Eintrag ohne Gate | `phantom` in die Liste | `p4-eintrag` |
| M4 | Punkt 2: jedes andere Gate | `kopf-check` in die Liste | `p2-schnell`, `p6-rot` |
| M5 | Punkt 6: roter Lauf lässt keinen Stand | `rm -f "$stand"` gestrichen | `lauf-gruen p2-schnell p3-neu p5-ende p5-stand p6-abbruch p6-andere p6-i-entfernt p6-rot` |
| M6 | Punkt 5: Hash zum Ende des Laufs | Hash vor `"$make_bin" "$@"` gebildet | `p5-ende` (nur dieser) |
| M7 | Punkt 6: `-n`/`-q`/`-t` vor `-i` | `*i*) ueberlesen=1; trocken=0` | `p6-ni-entfernt` (nur dieser) |
| M8 | Punkt 1/4: Gegenprobe hängt an `GATE_CHECKS` | `GATE_CHECKS += gates-schnell-gegenprobe` gestrichen | `anschluss-gegenprobe`, `p4-eintrag` |
| M9 | Punkt 7: nur der schnelle Stand | Stand an `gates-passed.diffsha` geschrieben | `p7-nachweis`, `p7-nachweis-neu` (und `p5-stand`, `p6-rot` u. a.) |
| M10 | Punkt 6: `-q` | Zeile für `q` gestrichen | `p6-q-entfernt` (nur dieser) |
| M11 | Punkt 1: `make gates` fährt die Kette nicht | `record-gates: $(GATE_CHECKS) gates-schnell` im `Makefile` | `p1-voll` (nur dieser) |

Mehrteilige Bedingung (Punkt 4): M2 und M3, je Hälfte eine. Verneinend: M5 (kein Stand bei Rot), M7/M10 (kein Stand und kein Entfernen bei `n`/`q`/`t`), M11 (`make gates` fährt die Kette nicht). Kein Mutant blieb grün. (Bei M2 meldet mein Hilfsaufruf „wirkungslos“, weil die neue Datei untracked ist und `git diff` sie nicht zeigt; die Gegenprobe wurde rot.)

**Nicht selbst gefahren:** `make gates` (volle Kette, auf Anweisung); `make lint-gegenprobe` und `make test-integration` einzeln; `make hook-gegenprobe`, `make a-check-negativ`. Zum Stand des vollen Laufs: `.harness/state/gates-passed.diffsha` (`d5856278…408f`) ist gleich dem Hash des Arbeitsbaums am sauberen HEAD `8708f21`; damit liegt ein voller Lauf zu diesem Baum vor (Hash-Gleichheit, nicht mein Lauf). Die in §7 genannten Hash-Anfänge `23d0f621…` (Baum des Commits davor) kann ich nicht nachprüfen. Nicht gefahren sind außerdem die Randformen von Teil 2 (Hook), weil Teil 2 nicht zugesagt ist.

## 2. DoD Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Liefer-Punkt 1: `SPEC-052` in §11, erster Commit vor jedem Code, Übertragung aus §6, Punkte nummeriert, Teil 2 gekennzeichnet, §12 Historie, Kennung löst auf | erfüllt | `cb37a34` ist der erste Commit nach `d075d3f`, Code folgt erst in `f633694`. `SPEC-052` Punkte 1 bis 20; Historie Zeile vom 2026-10-13 in `spezifikation.md` §12. Eingangssatz von Teil 2 wörtlich gleich der Zeile *Zwischenstand* in §6 („Die Zusagen dieses Teils gelten, soweit das Gate `stop-hook-gegenprobe` sie prüft, … ist der Hook-Teil insgesamt nicht zugesagt.“). Kein ADR-, Slice- oder Wellen-Verweis im neuen Text. `make docs-check`: 0 Befunde. Übertragung: die Zeilen aus §6 stehen als Punkte 1 bis 20; Ergänzungen sind Rückbezüge, der Absatz „Der Name jedes Falls …“ (aus *Zwischenstand*) und der Satz „Die Punkte 1 bis 7 prüft …“; Punkt 6 trägt die Entscheidung zu `-n`/`-q`/`-t` mit `-i` wörtlich. Kein Widerspruch zu [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md): Punkt 1 bis 4 entsprechen *Entscheidung* 1 bis 3, die Positivliste mit `docs/user/` und `README.md` ist die Setzung des Nutzers vom 2026-10-13 (die ADR lässt die Aufnahme mit Beleg als Änderung der Spezifikation zu); Punkt 19 folgt der ADR wörtlich (fehlender Stand verlangt `make gates`). |
| Liefer-Punkt 2: `gates-schnell` fährt `GATE_CHECKS` ohne Ausnahmeliste; Stand nur bei grünem Lauf; älterer entfällt; `-n`/`-q`/`-t`/`-i`; schreibt den vollen Nachweis nie; `make gates` unverändert | erfüllt | Proben 2 bis 5, 9, 10, 12; Mutationen M1 bis M7, M9 bis M11. Eingehalten ist auch die Vorrangregel (Probe 4 `MAKEFLAGS=ni`, M7). |
| Liefer-Punkt 3: Gegenprobe macht je Zusage einen roten Fall; hängt an `GATE_CHECKS`, `gates-schnell` nicht; Gate-Zeilen in `harness/README.md` mit Bindung [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md); Träger-Zeile in `harness/conventions.md` | erfüllt | Probe 1, M1 bis M11; Probe 9. `harness/README.md` führt zwei Zeilen (Kette, Gegenprobe) mit Bindung [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) und „seit slice-harness-gate-kette-schnell“; ihr Inhalt stimmt mit dem Code (Ausnahmeliste, Stand, `-i`, `-n`/`-q`/`-t` unberührt, Gegenprobe in der Ausnahmeliste, ohne Docker, `make gates` fährt die Kette nicht). `harness/conventions.md` nennt `harness/mk/gates-schnell.mk` mit `tools/harness/gates-schnell.sh` in der Zeile zu `grundlagen-durchsetzungsschicht.md`; die Gegenprobe-Datei hängt am Fragment. Die Auslegung 1 (nur die Gegenprobe hängt an `GATE_CHECKS`) ist in der DoD nachgezogen und vom Architect bestätigt. Eine Zeile zu `stop-hook-gegenprobe` fehlt, wie vorgesehen (Hook-Slice). |
| `make gates` grün | erfüllt laut Aufrufer; nicht von mir gefahren | Aussage von Implementer und Stop-Hook; Hash-Gleichheit des Nachweises mit dem sauberen HEAD (siehe oben). §7 nennt den Stand der Gates am Baum des Commits davor. |
| Review, Report unter `docs/reviews/` | erfüllt | Report vom 2026-10-13: 0 HIGH, 3 MEDIUM, 2 LOW, 4 INFO. |
| Mutationen als Tabelle in §7 (Zusage · Mutation · roter Fall), selbst gefahren | erfüllt | §7 trägt die Tabelle (15 Zeilen), die Nacharbeit-Tabelle (13 Zeilen) und die Zeile `m-i-vor-n`. Meine elf Mutationen bestätigen die roten Fälle, die ich nachgefahren habe; ich habe nicht jede Zeile nachgefahren. |
| Closure-Notiz mit Lerneintrag | offen (Closure) | §7 trägt bis zur Closure *Belege des Implementers*, keine Notiz |
| Beobachtungs-Register | offen (Closure) | keine Zeile gesetzt; die Klassen des Reviews (Stand ohne Lauf, Zusage ohne fangende Mutation, Zusage im Satz weiter als die Prüfung) gehören in den Lerneintrag |
| Risiken aus §6 mit Ausgang | offen (Closure) | vier Risiken mit vorgeschlagenem Ausgang in §6 (zweimal *entfallen*, einmal *eingetreten* mit den Hook-Slices, einmal *weiter offen* mit `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`); die Closure trägt sie |
| Drei Paarungen | offen (Closure) | |

## 3. Entscheidungs-Konformität

| Entscheidung | Urteil | Beleg |
|---|---|---|
| Namen: `gates-schnell`, `gates-schnell-gegenprobe`, `.harness/state/gates-schnell-passed.diffsha` (Architect, 2026-10-12) | erfüllt | `make help`, Fragment, Skript |
| Zugehörigkeit als Abzug, Ausnahmeliste im Fragment, Wortlaut `-gegenprobe` (F-700) | erfüllt | M1 bis M4; Muster `*-gegenprobe` in der Gegenprobe |
| Auslegung 1 und 2, sieben Einträge, Beleg gilt auch für `gates-schnell-gegenprobe` | erfüllt | Probe 9 und 10; Punkt 13 nennt die Kopien der Gegenprobe |
| Stand = Hash zum Ende des Laufs, unveränderte Ausgabe von `working-tree-hash.sh` | erfüllt | Probe 2, 5, 7; M6 |
| Roter oder abgebrochener Lauf: kein Stand, älterer entfällt; Ergebnis, nicht Mechanismus | erfüllt | Probe 5, 6; M5 |
| Make-Flags (F-694): `-n`/`-q`/`-t` unberührt, `-i` kein Stand und älterer entfällt, Vorrang von `-n`/`-q`/`-t` (Entscheidung vom 2026-10-13) | erfüllt | Probe 3 bis 5, M7, M10; die Langformen `--ignore-errors`, `--dry-run` verhalten sich gleich |
| Akzeptierte Negative (F-695): nicht in `SPEC-052` zugesagt (nur die schnelle Kette schreibt; der volle Lauf löscht den Stand nicht) | erfüllt | Punkt 5 bis 7 sagen es nicht zu; der Kopf der Gegenprobe nennt, dass sie `make gates` nie fährt |
| Teil 2 nicht zugesagt, bis `stop-hook-gegenprobe` existiert | erfüllt | `make help` führt kein solches Ziel; Eingangssatz wörtlich |
| Hook, `record-gates.sh`, `working-tree-hash.sh` unverändert | erfüllt | Probe 12 |
| Keine Lockerung eines Gates (§3.6) | erfüllt | Probe 12 und 9 |
| Aufruf-Voraussetzung (nacheinander) in Punkt 8, ohne Fall | erfüllt | Punkt 8; Kopf der Gegenprobe |

## 4. Plan gegen Diff (`AGENTS.md` §3.9) und Kommentare (§3.11)

- **Kopf, §1, §3, §6, §8:** folgen dem Diff. §3 nennt Spezifikation, Fragment, Skript, Gegenprobe, `harness/README.md`, `harness/conventions.md`; der Diff berührt genau diese und den Plan. Kein Produkt-Code, kein Handbuch, kein Test unter `internal/` oder `test/`. §8 zählt 3 Liefer-Punkte und 2 Schichten; der Diff stimmt.
- **F-694 bis F-702:** F-694 behoben (`gates-schnell.sh` prüft `MAKEFLAGS`, Probe 3 bis 5); F-695 behoben (Fälle `p5-ende`, `p7-nachweis`) oder durch den Architect aus `SPEC-052` genommen (Punkt 5 „nur die schnelle Kette“, Punkt 7 „volle Lauf berührt sie nicht“); F-696 behoben (`SPEC-052`: „Die Punkte 1 bis 7 prüft …“, Kopf der Gegenprobe, Gate-Zeile und Schlusszeile sagen 1 bis 7); F-697 behoben (Punkt 13 verweist auf Punkt 2, sieben Einträge); F-698 behoben (Kopf von `gates-schnell.mk` nennt Docker für die Gates der Kette, die Gegenprobe ruft es nicht); F-699 und F-700 entschieden und in Punkt 3, 4 nachgezogen (Überschuss wird gefangen, Probe M4); F-701 mit Ursache nachgefahren (Probe 6, gofmt in Stufe `test`); F-702 INFO, keine Aktion.
- **Kommentare:** Kopf von `gates-schnell.mk`, `gates-schnell.sh` und `gates-schnell-gegenprobe.sh`: jede Zusage hat einen Fall (Ausnahmeliste M1 bis M4, Stand M5 bis M9, Flags M7/M10, Anschluss M8/M11; Grenzen der Gegenprobe benannt: Attrappen, Punkt 8 ohne Fall, `make gates` nie gefahren). Die Behauptung „Die Gegenprobe ruft Docker nicht“ im Kopf des Fragments hat keinen eigenen Fall (Review F-698 meinte denselben Punkt); ich habe sie im Skript gegengelesen (`docker` kommt nicht vor), eine Mutation dazu nicht gefahren. Kein `//nolint`.
- **Commits:** reine Inhalts-Commits, kein `git mv`, jede Message nennt Slice-Kennung und [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md). (`6e22c4b` nennt zusätzlich `SPEC-052`; `AGENTS.md` §5 Regel 1 zählt das als Hinweis, kein Sensor prüft es, kein Befund im Gewicht.)

## 5. Befunde

Kein HIGH, kein MEDIUM.

| ID | Schwere | Befund | Beleg |
|---|---|---|---|
| V-320 | INFO | Die Größe liegt mit 563 Zeilen 48 über der Schätzung des Architect von 475 bis 515, unter der Rückführungsgrenze von 650 (87 darunter). §7 berichtet es offen (129 mehr als der erste Stand, Mehr in der Gegenprobe +99). Eine Review-Sitzung hat getragen. | Probe 13; §7 *Größe* |
| V-321 | INFO | In §7 steht die Überschrift „Auslegung, die ich dem Architect zur Bestätigung vorlege“, obwohl der Absatz darunter die Bestätigung nennt. Der Text ist in sich stimmig, die Überschrift liest sich als offen. Die Closure kann sie beim Aufräumen von §7 berichtigen. | Plan §7, Zeilen 371 bis 374 |
| V-322 | INFO | Die Zusage „die Gegenprobe ruft Docker nicht“ (Kopf von `gates-schnell.mk`) ist durch Lesen des Skripts belegt, nicht durch einen Fall (§3.11). Sie steht im Fragment als Beschreibung der Grenze, nicht als Zusage der Spezifikation. Kein Handlungsbedarf, solange der Satz nur beschreibt. | Probe 1; Skript |

## 6. Zu den Prüfpunkten des Auftrags

- **Laufzeit der Kette:** 7,2 s warm (§7: 7,2 bis 7,8 s), 97,8 s mit geänderter Go-Datei (§7: 102 s). Die Zeit ist keine Zusage der Spezifikation.
- **Roter Lauf mit geänderter Go-Datei (F-701):** belegt, Ursache gofmt in Stufe `test`, kein Stand danach.
- **Teil 2 nicht zugesagt:** ja, `make help` führt kein `stop-hook-gegenprobe`.
- **`GATE_CHECKS`:** nur die Gegenprobe der Kette hängt an; die Kette selbst nicht (`make -n gates`).

## 7. Offen und nicht gefahren

- Offen, weil die Closure sie führt: Closure-Notiz mit Lerneintrag, Beobachtungs-Register, Ausgang jedes Risikos aus §6, drei Paarungen, `git mv` nach `done/`.
- Nicht gefahren: `make gates`, `make lint-gegenprobe`, `make test-integration`, `make hook-gegenprobe`, `make a-check-negativ`; nicht jede Zeile der Mutationstabellen in §7 (elf eigene Mutationen, mehrere Zeilen damit nachgefahren). Der Hook (`.claude/hooks/stop-require-gates.sh`) wurde nicht gefahren, er ist unverändert.

## Verdikt

Die DoD-Punkte 1 bis 3 und die Entscheidungs-Konformität sind erfüllt; `make gates` ist laut Aufrufer grün und mit dem Hash des Nachweises zum sauberen HEAD vereinbar. Es bleiben nur die Closure-Pflichten. Kein Rückgabegrund an den Implementer oder den Architect.
