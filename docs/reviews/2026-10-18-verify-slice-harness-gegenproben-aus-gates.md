# Verifikations-Report: slice-harness-gegenproben-aus-gates (Kettenumbau) — 2026-10-18

**Rolle:** Verifier. **Stand:** Zweig `rev-kette`, Commit `2f9717fc` (Basis `b4ce87d3`), Worktree `/Development/pgwire-recorder-revkette`, Host (Locale `de_DE.UTF-8`, GNU make 4.3, u-boot, node, ssh-keygen, Docker vorhanden).
**Plan:** [slice-harness-gegenproben-aus-gates](../plan/planning/done/slice-harness-gegenproben-aus-gates.md) (§2 DoD, §7 Belege). **Verträge:** [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) Teil 1, [SPEC-052](../../harness/vertraege/SPEC-052-gate-kette.md), [ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md) (Accepted).
**Regel des engsten Sensors:** `make gates` wurde nicht gefahren; die Belege des Implementers (261 s grün, vorher 1630 s, zwölf Gegenproben einzeln grün) sind nicht nachgemessen und hier als Behauptung des Implementers geführt.

## Verdikt

**Nicht bestätigt: eine Bestätigung der DoD scheitert an einem Befund am Host.**
Die Kette ist wie zugesagt umgebaut, und alle Mutanten, die ich selbst gesetzt habe, färben genau das Erwartete rot. `make gates-schnell-gegenprobe` ist am Host mit der Locale des Hosts aber **rot** (V-1), und die Prüfung der Namensregel am Arbeitsbaum hat keinen roten Mutanten (V-2). Beides hat das Review als F-1 und F-2 schon gefunden ([Review-Report](2026-10-18-review-slice-harness-gegenproben-aus-gates.md), untracked im Worktree); ich bestätige beide mit eigenem Lauf.

## Befunde

| Nr | Schwere | Befund | Beleg (eigener Lauf) |
|---|---|---|---|
| V-1 | MEDIUM, DoD-blockierend | `make gates-schnell-gegenprobe` am Host mit `LANG=de_DE.UTF-8` rot. `tools/harness/gegenproben-regel.sh` (`gp_rezeptziele`) erkennt Ziele am englischen Text `#  recipe to execute (from 'harness/mk/…`. GNU make 4.3 schreibt unter deutscher Locale `#  Auszuführendes Rezept (aus „harness/mk/lint.mk“, Zeile 18):`; die Regel findet dann kein Ziel mit Rezept, und jeder Eintrag der Listen gilt als Eintrag ohne Ziel (V6). | Ohne Änderung der Umgebung: Exit 2, 13 Fälle der Regel rot (`p2-v6-container`, `p2-v6-host`), Mutanten färben falsche Mengen; `make -s -pq help \| grep -c 'recipe to execute (from'` gibt 0. Mit `LC_ALL=C make gates-schnell-gegenprobe`: grün (41 s). Der Implementer hat nur im Dev-Container (englische Locale) gemessen, deshalb ohne Befund. Behebung gehört dem Implementer (z. B. `LC_ALL=C` im Aufruf von `make -pq` in `datenbank()` und im Skript der Regel, plus Fall dazu); SPEC-060 Plattform sagt Linux-Host mit GNU-Werkzeugen und nennt keine Locale. |
| V-2 | MEDIUM | Grüner Mutant: Die Zusage „die Regel gilt im Arbeitsbaum selbst“ (SPEC-060 Punkt 2; §7 Zeile `deckung "$wurzel"`: „n/a … von Hand gesehen“) hat in `gates-schnell-gegenprobe.sh` keinen roten Fall. Ersetze ich die Zeile `ist="$(deckung "$wurzel" …` durch `ist=""`, bleibt die Gegenprobe grün (Mutant A8, Exit 0, keine Zeile ROT). Die Mutanten, die den Arbeitsbaum selbst verletzen (A1 bis A6 unten), fängt nur diese Zeile; sie selbst fängt keiner. | A8: `exit=0`, `gruen`. Gleiches gilt für den Arbeitsbaum-Teil von SPEC-052 Punkt 4 (Bestand, nicht neu). `AGENTS.md` §3.10: je Zusage eine Mutation. |
| V-3 | INFO | Fall der Gegenprobe des Container-Starts: erster Lauf am Host rot (110 Fälle „an der unveränderten Kopie rot“), weil das Temp-Verzeichnis `/tmp/pgr-cs-H2D7L6` während des Laufs verschwand (`Cannot find module '/tmp/pgr-cs-H2D7L6/pruef.js'`, `u-boot.yaml: No such file`). Zweiter Lauf desselben Stands grün. Ursache nicht geklärt: Ein anderer Lauf am Host (Implementer des Slice `-start`, Docker-Daemon dort gestartet 11:01) oder ein Aufräumen von `/tmp` ist möglich, nichts im Diff dieses Slice (`container-start-gegenprobe.sh`: nur der Anker von `p16-ausnahme`). Kein Befund gegen den Slice, aber die Gegenprobe ist gegen Eingriff in `/tmp` nicht fest. | `cs.log` (rot, 3 m 47 s), `cs2.log` (grün, 4 m 20 s). |
| V-4 | LOW | Der Vertrag nennt elf Fälle (Teil 1 Punkt 2), die Gegenprobe führt dreizehn (`p2-v4-host`, `p2-rezept-ausserhalb-fragment`), Kopf von `gates-schnell-gegenprobe.sh` sagt noch „elf“; Rückgabe 3 des Implementers ist offen, `SPEC-060` unverändert. Das Review führt es als F-5 und F-7. | Gegenprobe-Ausgabe: „dreizehn Fälle“; `SPEC-060` Punkt 2 „Fälle, je einer für N1 und N2 …“. |
| V-5 | LOW | `Makefile` Z. 2 sagt weiter „jedes [Fragment] haengt seine Checks an GATE_CHECKS“ (Rückgabe 4, Datei ist generiert, nicht im Plan). Aussage stimmt für Gegenproben nicht mehr. | `grep -n GATE_CHECKS Makefile`. |

## DoD-Punkte gegen Belege und eigene Läufe

| DoD-Punkt | Stand | Beleg |
|---|---|---|
| Liefer-Punkt 1: kein `-gegenprobe`/`a-check-negativ` an `GATE_CHECKS`, jedes an genau einer Liste, `GATE_SCHNELL_AUSNAHME` = drei Gates | **bestätigt** für die Kette | `make -pn gates`: `GATE_CHECKS := abdeckung-check a-check baseline-verify build test test-container docs-check test-integration kopf-check lint meldungskatalog-check test-postgres-versionen`; `GEGENPROBEN_CONTAINER` zwölf Einträge (inkl. `a-check-negativ`, `gates-schnell-gegenprobe`), `GEGENPROBEN_HOST = container-start-gegenprobe`; `GATE_SCHNELL_AUSNAHME := test-container test-integration test-postgres-versionen`; `record-gates` hängt an `GATE_CHECKS` ohne Gegenprobe. `grep GATE_CHECKS harness/mk/*.mk` zeigt für Gegenproben nur Kommentare. Dreizehn Ziele mit Namen auf `-gegenprobe`/`negativ` in der Datenbank von `make`, alle in genau einer Liste. |
| Liefer-Punkt 1: Fälle und Mutationen in `gates-schnell-gegenprobe` | **bestätigt mit V-1, V-2, V-4** | Mit `LC_ALL=C` grün, 13 Fälle; Mutanten siehe unten. Am Host ohne `LC_ALL=C` rot (V-1). |
| Liefer-Punkt 2: Anschluss der acht Gegenproben prüft die Liste | **durch Implementer belegt, nicht nachgefahren** | Belege in §7 (Tabelle, Läufe mit Dauer); ich habe nur die engsten Sensoren gefahren (`gates-schnell-gegenprobe`, `container-start-gegenprobe`). Stichprobe: `container-start-gegenprobe` Punkt 16 (`GP_NUR='^p16' GP_MUT='^p16' GP_ENTDECKE=1`) färbt je Mutation ihren Fall: `p16-gate-checks => p16-GATE_CHECKS`, `p16-scharf => p16-scharf`, `p16-ausnahme => p16-GATE_SCHNELL_AUSNAHME`, `p16-ziel-gate => p16-ziel-kein-gate`; grün. |
| Liefer-Punkt 3: Nachzug Aussagen, Gate-Index, `AGENTS.md` §6 | **bestätigt, Rest V-5** | Diff von `harness/README.md` und `AGENTS.md` gelesen: Zeilen aller Gegenproben nennen ihre Liste, `make gates`/`gates-schnell` neu gefasst, `AGENTS.md` §6 Schritt 6 trägt die erste Hälfte des Satzes (Wortlaut der Entscheidung des Nutzers, bis „Zuordnung in `harness/README.md` §Sensors“). `harness/sensors/*.md`, `harness/conventions.md` nachgezogen bis auf Z. 134 („jedes hängt seine Prüfungen an `GATE_CHECKS`“, in `conventions.md` unverändert zu lesen) und V-5. |
| `make gates` grün, Zeit | **nicht nachgefahren (Anweisung)** | Zahlen des Implementers, §7: 261 s grün im Dev-Container, vorher 1630 s (laut Auftrag des Koordinators), Log mit `EXIT=0`; „Zeit nicht zugesagt“ (SPEC-060 Punkt 15). |
| `make gates-schnell-gegenprobe` grün | **am Host nur mit `LC_ALL=C`** | V-1. |
| Beobachtungs-Register, Review, Closure-Notiz, Risiko-Ausgänge, Paarungen | **offen, zeitlich vor der Closure** | §7 Lerneintrag, Risiken und Paarungen stehen auf „—“; DoD-Häkchen nicht gesetzt (richtig vor Closure). Review-Report liegt vor, Verdikt dort: merge-blockierend (F-1, F-2). |

## Eigene Läufe (Host)

| Lauf | Ergebnis | Dauer |
|---|---|---|
| `make docs-check` | grün, 0 Befunde (1082 Dateien) | 1 s |
| `make kopf-check` | grün | <1 s |
| `make gates-schnell` | grün | 114 s |
| `make gates-schnell-gegenprobe` (Locale des Hosts) | **ROT**, 13 Fälle plus Mutanten | 30 s |
| `LC_ALL=C make gates-schnell-gegenprobe` | grün | 41 s |
| `make container-start-gegenprobe` ganz, erster Lauf | rot (V-3: Temp-Verzeichnis verschwunden) | 227 s |
| `make container-start-gegenprobe` ganz, zweiter Lauf | **grün** (Rückgabe 1 des Implementers damit am Host erledigt) | 260 s |
| `make -pn gates` / Namensregel am Worktree | siehe oben; `namensregel` auf der Datenbank des Hosts ohne `LC_ALL=C`: 13 Verstöße (V-1) | — |

## Mutanten der Kette (eigene, in Kopien des Stands unter dem Scratchpad, `LC_ALL=C`, je ein Lauf der Gegenprobe; erwartet rot)

Alle rot (Exit 2). Spalte „rot, weil“ nennt die erste Meldung.

| Nr | Klasse | Mutation | rot, weil |
|---|---|---|---|
| A1 | V1 N1 | `GATE_CHECKS += lint lint-gegenprobe` (`lint.mk`) | Arbeitsbaum: `p2-v1-n1` |
| A6 | V1 N2 | `GATE_CHECKS += a-check-negativ` (`arch-negativ.mk`) | Arbeitsbaum: `p2-v1-n2` |
| A2 | V2 | `kopf-check-gegenprobe` zusätzlich in `GEGENPROBEN_HOST` | Arbeitsbaum: `p2-v2-n1` |
| A3 | V3 | Zeile von `hook-gegenprobe` in der Liste gestrichen | Arbeitsbaum: `p2-v3-n1` |
| A4 | V4 | Zeile von `a-check-negativ` verdoppelt | Arbeitsbaum: `p2-v4-n2` |
| A5 | V6 | `GEGENPROBEN_CONTAINER += phantom` | Arbeitsbaum: `p2-v6-container` |
| A7 | SPEC-052 Punkt 4 | `lint-gegenprobe` in `GATE_SCHNELL_AUSNAHME` | `p4-eintrag` (Kopie `ok`, Arbeitsbaum, vier Mutanten mit unerwarteter Menge) |
| B1 | Regel, V4 | `-gt 1` zu `-gt 2` | `p2-v4-host p2-v4-n1 p2-v4-n2` grün erwartet, rot |
| B2 | Regel, Klasse | `a-check-negativ`-Zweig in `gp_klasse` gestrichen | `p2-v1-n2 p2-v2-n2 p2-v3-n2 p2-v4-n2` |
| B3 | Regel, `p2-rezept-ausserhalb-fragment` | Beschränkung auf `harness/mk/` entfernt | `p2-rezept-ausserhalb-fragment` |
| B4 | Regel, V2 | `&&` zu `\|\|` | `p2-v2-n1 p2-v2-n2` |
| B5 | Regel, V3 | `&&` zu `\|\|` | `p2-v3-n1 p2-v3-n2` |
| B6 | Regel, N1 | `*-gegenprobe` zu `*gegenprobe*` | `p2-andere-namen` |
| B7 | Regel, V1 | `-gt 0` zu `-gt 1` | `p2-v1-n1 p2-v1-n2` |
| **A8** | **Arbeitsbaum-Prüfung** | Zeile `ist="$(deckung "$wurzel" …` durch `ist=""` ersetzt | **grün (Exit 0): V-2** |

Die Mutanten A1 bis A6 färben die Gegenprobe über den Arbeitsbaum, nicht über einen Fall der Kopie. Dass dieser Zweig selbst nichts fängt (A8), ist V-2. Bei B1 bis B3 meldet die Gegenprobe zusätzlich „Mutation … ändert die Regel nicht“, weil meine Änderung die Zeile traf, an der ihre eigene Mutation ansetzt; das ist ein Artefakt meiner Mutanten, kein Befund.

## Plan-vs-Code-Diff

- Liefer-Punkte 3, Schicht 1 (Harness): wie §8. Dateien nach `git diff --stat b4ce87d3..HEAD`: 48, davon die Kette, zwölf Fragmente, `gates-schnell-gegenprobe.sh`, neue `tools/harness/gegenproben-regel.sh`, sieben Gegenproben, `harness/README.md`, `AGENTS.md`, SPEC-052 bis -058 und neu SPEC-060 (Architect, vor dem Start), Pläne (Nachzug des Planners).
- Nicht im Plan §3, im Diff: `tools/harness/gegenproben-regel.sh` (neue Datei, nur in §7 „Neu“), `harness/mk/container-start-gegenprobe.mk`, `container-start.mk` (§3 nennt sie unter „Nachzug“). Plan folgt dem Diff in §3 nicht (`AGENTS.md` §3.9; Review F-7).
- Dreizehn statt elf Fälle in Teil 1; der Vertrag nennt nur elf (V-4).
- Größe: Metrik `git diff --numstat b4ce87d3..HEAD -- tools harness AGENTS.md` laut Implementer 707, ohne den Vertrag des Architect 438; ich habe sie nicht neu gerechnet (Review F-11 hat sie nachgerechnet und bestätigt). Marke 900 nicht erreicht.
- Zeitvergleich: nur die Zahlen des Implementers (261 s gegen 1630 s). Meine Messung der schnellen Kette: `make gates-schnell` 114 s.

## Offene Punkte an Implementer, Architect, Planner

1. Implementer: V-1 beheben (`LC_ALL=C` für den Aufruf von `make -p…` in Regel und Gegenprobe) und einen Fall für die Locale, der ohne die Behebung rot wird; V-2 einen Fall, der die Prüfung im Arbeitsbaum färbt.
2. Architect: Fälle über den Vertrag hinaus (V-4), Initialisierung der Listen in `SPEC-060` Punkt 1 (Rückgabe 3), Locale als Voraussetzung nennen.
3. Planner: `Makefile` Z. 2 (V-5, Rückgabe 4) und Plan §3 (F-7).
4. V-3 beobachten: tritt die Gegenprobe des Container-Starts erneut rot ohne Code-Änderung auf, ist das eine Beobachtung für das Register.

Aufgeräumt: Alle Kopien liegen unter dem Scratchpad; kein Container und kein Volume von mir (die Gegenprobe des Container-Starts räumt `/tmp/pgr-cs-*` selbst, `ls /tmp/pgr-cs*` leer). Nichts committet.

## Nachverifikation (Stand a8bb5dd7)

Engste Sensoren; die Läufe, die der Implementer belegt (`make gates` 260 s, `coverage`, `lint`), habe ich nicht wiederholt.

| Befund | Ergebnis | Beleg |
|---|---|---|
| V-1 (Locale) | **behoben** | `LANG=de_DE.UTF-8 make gates-schnell-gegenprobe` am Host: Exit 0, 59 s, Meldung nennt den Fall `p2-locale`. Mutation von Hand in Kopie (`gegenproben-regel.sh`/Gegenprobe: `LC_ALL=C` entfällt) ist die Mutation `p2-locale-nicht-c` der Gegenprobe selbst; sie lief in diesem grünen Lauf (rot gefärbt wird genau `p2-locale`). Ob der Fall ohne die Behebung rot wäre, zeigt diese Mutation. |
| V-2 / A8 (`ist=""`) | **behoben** | Kopie mit `ist=""` statt `deckung`, ganze Gegenprobe: Exit 1, `ROT — Fall 'p2-arbeitsbaum': der Arbeitsbaum mit lint-gegenprobe in GATE_CHECKS bleibt grün`. (Die zweite Zeile `Mutation p2-arbeitsbaum-gestrichen ändert die Gegenprobe nicht` ist Artefakt meines Mutanten: die Zeile war schon mutiert.) |
| V-3 (Temp-Verzeichnis) | **nicht wieder aufgetreten, Gegenprobe grün** | `make container-start-gegenprobe` am Host einmal ganz: Exit 0, 256 s, `gruen — Teil I von SPEC-058`. Verzeichnis `/tmp/pgr-cs-RWWURK` blieb bis zum Ende. Einmal eine Zeile `grep: …/u-boot.yaml: No such file` im Log, kein Fall rot (im Lauf erwartet/folgenlos; nicht weiter verfolgt). Fremdlauf im Protokoll: ab 11:14 `bash tools/harness/container-start-gegenprobe.sh` (PID 496680, `/Development/pgwire-recorder-impl3`, Verzeichnis `/tmp/pgr-cs-496680`, später `pgr-cs-48RUAl`), lief gleichzeitig mit meinem Lauf (11:36 bis 11:40) mit demselben Präfix `pgr-cs-`; dazu zehn leere `pgr-cs-*`-Verzeichnisse von 11:19 (Rest des Fremdlaufs). Kein Eingriff in mein Verzeichnis beobachtet. Der Nachweis, dass ein gleichzeitiger Lauf unschädlich ist, ist damit ein Fall von einem Lauf, kein Beleg. |
| V-4 (13 statt 11 Fälle, Vertrag) | **offen** | `SPEC-060` Punkt 2 nennt weiter nur „je einer für N1 und N2“ ohne `p2-v4-host`, `p2-rezept-ausserhalb-fragment`, `p2-locale`; keine Locale-Voraussetzung im Vertrag (Rückgabe an den Architect, vom Implementer in §7 benannt, nicht umgesetzt). Keine Änderung am Code nötig. |
| V-5 (`Makefile` Z. 2) | **offen** | Z. 2 sagt weiter „jedes haengt seine Checks an GATE_CHECKS“ (generierte Datei, Rückgabe an den Planner). |

Weitere Prüfungen:

- **F-8-Mutationen:** `make kopf-check-gegenprobe` und `make meldungskatalog-check-gegenprobe` grün (je 2 s). Ohne die Zeile `GEGENPROBEN_CONTAINER += …` im Fragment (Kopie) werden `nr9-liste` und `nr11-liste` rot (`steht nicht in GEGENPROBEN_CONTAINER`). Für `lint` (`p10-gegenprobe-liste`) und `coverage` (`p11-gegenprobe-liste`) nur die Mengenprüfung: dieselbe `make -pq`/awk-Prüfung an einer Kopie, Zeile gestrichen, Ziel nicht mehr in der Liste. Die Skripte selbst habe ich nicht gefahren (Docker, 100 s bis 280 s; Belege des Implementers).
- **F-3 `m-liste-in-schnell`:** Oracle-Zeile (`p2-gegenprobe-in-schnell`) in Kopie gestrichen: `ROT — Mutant 'm-liste-in-schnell' färbt 'p2-gegenprobe-in-schnell' nicht rot`. Fall hat Zähne.
- **F-4 `p2-andere-namen`:** Echte Mutationen der Regel in Kopie, `*-gegenprobe | gegenprobe` und `a-check-negativ | x-negativ`: beide färben `p2-andere-namen` rot (`Regel: Fälle rot, die grün sein müssen: p2-andere-namen`). Die übrigen drei Namen (`gegenproben`, `gegenproben-host`, `x-gegenprobe-alt`) lief nur die eingebaute `regelmut`-Schleife im grünen Lauf.
- **Hinweis:** In meinen Kopien unter dem Scratchpad (ohne `.git`, ohne die übrigen Verzeichnisse) war die Prüfung des Arbeitsbaums schon an der unveränderten Kopie rot (`p2-v3-n1 p2-v3-n2`). Im echten Worktree ist sie grün; Artefakt der Teilkopie, kein Befund.
- **Sensoren (5):** `make docs-check` 0 Befunde (1084 Dateien), `make kopf-check` Exit 0, `make gates-schnell` Exit 0 (10 s, Cache).
- **Metrik (6):** `git diff --numstat b4ce87d3..HEAD -- tools harness AGENTS.md`: 635 Einfügungen + 183 Löschungen = **818**; ohne `harness/vertraege` **540**. Beide Angaben des Implementers stimmen. Marke 900 nicht erreicht.

Aufräumen: Kopien nur unter dem Scratchpad; kein Container und kein Volume von mir. Nichts committet.
