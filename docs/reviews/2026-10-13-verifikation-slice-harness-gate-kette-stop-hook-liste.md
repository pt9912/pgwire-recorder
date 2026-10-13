# Verifikation: slice-harness-gate-kette-stop-hook-liste — 2026-10-13

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-13-review-slice-harness-gate-kette-stop-hook-liste.md`, Commit `074285b`).

**Gegenstand:** Slice-Plan `slice-harness-gate-kette-stop-hook-liste` (Kopf, §1 bis §8, mit den Entscheidungen des Architect vom 2026-10-12/13) gegen `9662e13..23f9381`. Commits: `b259381` (Code), `e758865` (Belege), `32a5a47` (Urteil des Planners zur Marke), `c005555` (Architect: Rückgaben), `e92db76`/`78b35ab`/`9a1cc9d` (Nacharbeit Punkt 21/22), `074285b` (Review), `d1ae3b6` (Architect: F-704, F-707 b, F-710), `4aff64b`/`cc336f2`/`876fc01`/`d16d816`/`23f9381` (Nacharbeit zum Review), `b397d6f` (ADR-0040/MR-002, nicht dieser Slice).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-13

**Eingang:** Plan ganz, `SPEC-052` Teil 2 (Punkte 9 bis 22, `spec/spezifikation.md` §11), [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) (Accepted), Review-Report, Diff `9662e13..HEAD -- tools harness .claude spec`. Den Bericht des Implementers habe ich nicht gelesen; ich lese seine Belege in §7 des Plans.

## 1. Selbst gefahrene Sensoren

Ergebnis vorweg: kein HIGH, kein MEDIUM. Der Hook sperrt in jeder Probe, die ich von Hand gefahren habe, wie `SPEC-052` es sagt; die Gegenprobe wird unter neun eigenen Mutationen rot, an den Fällen, die der Plan nennt. Der Arbeitsbaum war danach unverändert (`git status --short` leer; die Proben liefen in Temp-Repos und einer Kopie im Scratchpad).

| Nr. | Sensor / Probe | Ergebnis |
|---|---|---|
| 1 | `make stop-hook-gegenprobe` | Exit 0, 31,1 s: „66 Fälle der Punkte 1, 10, 11, 14, 16 bis 22 je Mutant rot, unveränderte Kopie grün, Anschluss an GATE_CHECKS“ (§7: 31 s) |
| 2 | `make gates-schnell-gegenprobe` | Exit 0, 21,7 s |
| 3 | `make gates-schnell` | Exit 0, 8,7 s (Cache warm); Stand `.harness/state/gates-schnell-passed.diffsha` gleich der Ausgabe von `tools/harness/working-tree-hash.sh` (`99dd0650…5d34`) |
| 4 | `make docs-check` | 685 Dateien, 0 Befunde |
| 5 | `make kopf-check`, `make hook-gegenprobe` | Exit 0 (`hook-gegenprobe`: „gruen“) |
| 6 | Gegenprobe als root (`docker run -u 0`, Image `ai-harness-init:host`, bash 5.2.37, Arbeitsbaum read-only) | Exit 0; `p20-stand-unlesbar` und Mutant `m-unlesbar-als-unveraendert` laut übersprungen (F-705 behoben) |
| 7 | Hook von Hand gegen ein Temp-Repo (Kopie von Hook, `dateiliste.sh`, `record-gates.sh`; Eingabe `{}`) | siehe unten |
| 8 | `record-gates.sh` unter `make` an einer Attrappe (zwei Gates, eines mit `FAILCMD`), vorher grüner Lauf als älterer Nachweis, danach Änderung von `a.txt` | siehe unten |
| 9 | `make -n gates` | enthält `stop-hook-gegenprobe` einmal; `make help` führt das Ziel |
| 10 | Ausnahmeliste `GATE_SCHNELL_AUSNAHME` | acht Einträge: sieben des Gebers plus `stop-hook-gegenprobe`; die Gegenprobe der Kette führt dieselben acht (`LANGSAM`) und ist grün |
| 11 | Größe nach der Metrik des Planners (`git diff --numstat 9662e13..<c> -- tools harness .claude`, Einfügungen plus Löschungen) | `d16d816`: 979 + 86 = **1065** (stimmt mit §7); `HEAD`: 1020 + 87 = 1107, davon 40 + 2 aus `MR-002` und dem Index-Eintrag in `harness/conventions.md` (Commit `b397d6f`, ADR-0040, nicht dieser Slice); `HEAD` ohne die `MR-002`-Datei: 980 + 87 = 1067 |

**Zu 7, Hook von Hand** (jeweils `{"decision":…}`; Temp-Repo mit gitignoriertem `.harness/state/`):

| Fall | Antwort |
|---|---|
| kein Stand, sauberer Baum, `stop_hook_active` fehlt / `false` | `block` („file list missing“) |
| kein Stand, sauberer Baum, `stop_hook_active: true` | `approve` (Schleifenschutz, Punkt 20) |
| nach `record-gates.sh`, Baum unverändert | `approve` |
| dasselbe aus Unterverzeichnis `sub/` und unter `env -i "$(command -v bash)"` | `approve` |
| geänderte Datei (Wurzel, Unterverzeichnis, `env -i`) | `block` (alle drei) |
| untracked, gelöschte, umbenannte Datei | `block` |
| Pfade ` lead`, `trail `, Tabulator, Backslash, Zeilenumbruch, Bytes `\xff\xfe`, je neu angelegt | `block`; nach dem Entfernen wieder `approve` |
| Pfad ` lead` bereits im Stand: unverändert / Inhalt geändert / entfernt | `approve` / `block` / `block` |
| Konflikt im Index bei frisch geschriebenem Stand (identische Liste) / mit `stop_hook_active: true` | `block` („unresolved merge entries“) / `approve` |
| Stand `F\0abc\0a.txt\0` / leere Datei / Datei gelöscht | `block` (beschädigt) / `block` / `block` („missing“) |
| leere Standardeingabe (`</dev/null`) bei identischer Liste | `approve` |

**Zu 8, Nachweis unter `make`** (Hash über Liste plus Hash-Datei vorher und nachher):

| Aufruf | Exit | Nachweis |
|---|---|---|
| `make gates` grün | 0 | neu geschrieben |
| `-i` mit rotem Gate / `-i` grün / `-ki` / `--ignore-errors` / `-j4 -i` / `MAKEFLAGS=i` | 0 | beide Dateien entfernt |
| `MAKEFLAGS=ni`, `-n`, `-ni` | 0 | byte-gleich |
| `-q`, `-qi` | 1 | byte-gleich |
| `-t` | 0 | byte-gleich |
| rotes Gate ohne `-i` (auch `-k`) | 2 | älterer Nachweis bleibt (akzeptiertes Negativ, Punkt 21) |

### Mutationen (eigene Kopie im Scratchpad, `git archive HEAD`, je Mutant `bash tools/harness/stop-hook-gegenprobe.sh`)

Jeder Lauf endete mit Exit 1. Die Zeile `Kopie 'ok' ist nicht grün` nennt die roten Fälle der Kopie `ok`:

| Nr. | Zusage | Mutation | roter Fall der Kopie `ok` |
|---|---|---|---|
| V1 (M6) | Punkt 11: führendes/abschließendes Leerzeichen, Tabulator | `IFS=` im `read` von `dateiliste_erzeugen` entfernt | `p11-fuehrendes-leerzeichen`, `p11-abschliessendes-leerzeichen`, `p11-tabulator` (nur diese drei) |
| V2 (M7) | Punkt 11: Backslash | `read` ohne `-r` | `p11-backslash` (nur dieser) |
| V3 | Punkt 20: Schleifenschutz vor `git ls-files -u` | Schutz hinter die Zeile `konflikt="$(git ls-files -u)"` | `p20-schleifenschutz-kein-repo`, `-ohne-git`, `-git-u`, `-lib-fehlt` |
| V4 | Punkt 19: fehlender Stand bei sauberem Baum | bei sauberem Baum `frei` | `p19-ohne-stand-sauber`, `p20-schleifenschutz-ohne-stand` |
| V5 | Punkt 21: `-i` | Zweig `*i*)` in `record-gates.sh` unwirksam (`*nie*)`) | `p21-i-rot`, `-gruen`, `-env`, `-ki`, `-lang`, `-j` |
| V6 | Punkt 20: Abbruch antwortet `block` | `trap sicherung EXIT` entfernt | `p20-abbruch`, `p20-lib-fehlt`, `p20-kein-repo`, `p20-ohne-git`, dazu die Fälle, deren Schritt scheitert |
| V7 | Punkt 1: Liste nur bei grünem Lauf | `record-gates: $(GATE_CHECKS)` im `Makefile` ohne Voraussetzungen | `p1-liste-nicht-bei-rot` (nur dieser) |
| V8 | Punkt 20: Schleifenschutz bleibt | Muster `"stop_hook_active": true` unwirksam | `p20-schleifenschutz` und alle acht `p20-schleifenschutz-*` |
| V9 | Punkt 14: Inhalt, Link-Ziel, gelöschte Datei | Vergleich nur der Pfadfelder | `p14-inhalt`, `p14-link-ziel`, `p14-geloescht` (und `p11-*` durch die Leerfelder) |

Bei jedem Mutanten melden außerdem die eigenen Mutanten der Gegenprobe („ändert nichts“, „färbt auch … rot“), weil meine Mutation den Hook oder die Liste schon verändert hat; das ist die erwartete Folge der Mutation einer Kopie, die die Gegenprobe selbst nochmals mutiert. Maßgeblich ist die Rotfärbung der Kopie `ok` an den genannten Fällen. Kein Mutant blieb grün. V1 und V2 bestätigen, dass die in F-703 als grün gefundenen Mutanten M6 und M7 jetzt je einen eigenen roten Fall haben; V3 bestätigt dasselbe für F-704.

## 2. DoD Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Liefer-Punkt 1: Nachweis des vollen Laufs ist die Dateiliste; `record-gates.sh` schreibt nur nach grünem Lauf; `working-tree-hash.sh` liefert dieselbe Liste/den Hash; eine Logik (`dateiliste.sh`); `-i` schreibt nichts und entfernt Älteres; `-n`/`-q`/`-t` lassen stehen | erfüllt | `record-gates.sh`, `working-tree-hash.sh` und der Hook laden `tools/harness/dateiliste.sh` mit `source`; `working-tree-hash.sh` hat 12/20 Zeilen, keine eigene Logik. Probe 8 (alle sechs Zeilen), V5, V7; Sensor 3 (Hash der schnellen Kette bleibt; Stand gleich Ausgabe von `working-tree-hash.sh`). `gates-passed.diffsha` ist der Hash der Liste (`record-gates.sh`: `sha256sum` über die geschriebene Liste; `working-tree-hash.sh`: `dateiliste_erzeugen | sha256sum`, gleiche Bytes). |
| Liefer-Punkt 2: Hook gibt nur bei identischer Liste frei; fehlender/beschädigter Stand auch bei sauberem Baum, nicht lesbarer Vergleich, Konflikt, Abbruch → `block`; Schleifenschutz bleibt; Kopf beschreibt das neue Verhalten; Konflikt-Prüfung zuerst nach dem Schleifenschutz; keine Umgebungsvariable | erfüllt | Hook (`.claude/hooks/stop-require-gates.sh`): Schutz, dann `git rev-parse`, `source`, `git ls-files -u`, Stand fehlt, `dateiliste_pruefen`, `cmp`. Probe 7 in allen Zeilen; V3, V4, V6, V8. Der Kopf sagt die Sperre bei fehlendem Stand auch bei sauberem Baum zu und nennt „Restlücke“ nicht mehr. Der Hook liest keine Umgebungsvariable (`env -i`-Probe, `p22-leere-umgebung`). |
| Liefer-Punkt 3: Gegenprobe macht je Randform einen Fall, nur im inerten Fall freigebend (dazu `p21-*`, `p22-*`), hängt an `GATE_CHECKS`; Zeile in `harness/README.md` §Sensors mit Bindung [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md); Träger-Zeile in `harness/conventions.md`; Name `stop-hook-gegenprobe` | erfüllt | Probe 1, 9, 10; V1 bis V9. `harness/README.md` Zeile 79: Bindung ADR-0038 und „seit slice-harness-gate-kette-stop-hook-liste“, Punkte „1, 10, 11, 14, 16 bis 22“ stimmen mit der Meldung der Gegenprobe überein. `harness/conventions.md` nennt in der Zeile `grundlagen-durchsetzungsschicht.md` Hook, `dateiliste.sh`, `stop-hook-gegenprobe.mk` mit `.sh`. Der Name ist ein anderer als `hook-gegenprobe`. |
| `make gates` grün | erfüllt laut Aufrufer; nicht von mir gefahren | Aussage von Implementer und Stop-Hook; §7 nennt Exit 0, 5 min 48 s auf `b397d6f` mit „66 Fälle“. Dass der Stop-Hook im Arbeitsbaum jetzt `block` antwortet (Plan-Commit `23f9381` nach dem Lauf, §7 sagt es), ist erwartet und belegt, dass der Hook die Änderung seit dem Lauf sieht. |
| Review durchgeführt, Report liegt vor | erfüllt | Report vom 2026-10-13: 0 HIGH, 2 MEDIUM, 3 LOW, 5 INFO; fail-closed ausdrücklich geprüft (Negativbefunde). |
| Mutationen als Tabelle in §7, selbst gefahren (§3.10) | erfüllt | §7 trägt drei Tabellen (28 + 10 + 10 Zeilen, mehrere Mutanten je Zeile); das Skript führt 36 `sedmut`-Aufrufe und weitere Mutanten aus Schleifen. Meine neun Mutationen bestätigen die roten Fälle, die ich nachgefahren habe. Eine Zählung „70+ Mutanten“ konnte ich nicht nachvollziehen: ich komme auf rund 60 Zeilen-Mutanten in §7; nicht jede Zeile habe ich nachgefahren. |
| Closure-Notiz mit Lerneintrag | offen (Closure) | §7 trägt bis zur Closure nur *Belege des Implementers* |
| Beobachtungs-Register: Evidence zu `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (3×) | offen (Closure) | Datei `evidence/slice-harness-gate-kette-stop-hook-liste.md` noch nicht angelegt |
| Risiken aus §6 mit Ausgang | offen (Closure) | §6 trägt je Risiko einen vorgeschlagenen Ausgang; die Closure trägt sie ein, der Planner den Namen der Dateiliste in `slice-harness-gate-kette-stop-hook-freigabe` |
| Drei Paarungen | offen (Closure) | |

## 3. Entscheidungs-Konformität

| Entscheidung | Urteil | Beleg |
|---|---|---|
| [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) Punkt 3/4: Nachweis als Dateiliste, Hook mit Gegenprobe | erfüllt | Probe 7, Sensor 1; Positivliste bewusst dem Folge-Slice überlassen (Hook sperrt jeden Unterschied) |
| Nutzer 2026-10-12: fehlender Stand bei sauberem Baum sperrt (Änderung der Zusage) | erfüllt | Probe 7 (Zeile 1), V4; Hook-Kopf sagt es zu |
| Architect 2026-10-13 (`c005555`): `-i` im Nachweis; `-n`/`-q`/`-t` unberührt; Aufrufarten (Punkt 22) | erfüllt | Probe 8 vollständig, V5; Probe 7 (`env -i`, Unterverzeichnis); `bash -e` und `CLAUDE_PROJECT_DIR` als Nicht-Gegenstand gesagt und ohne Fall |
| Architect `d1ae3b6`: Schleifenschutz als akzeptierte Grenze (Punkt 20); leere Umgebung statt „liest keine Variable“; bash ab 4 als Voraussetzung (Punkt 11) | erfüllt | Hook-Kopf, Kopf der Gegenprobe und Kopf von `dateiliste.sh` sagen es; `SPEC-052` Punkte 11, 20, 22 tragen den Wortlaut; V3, V8 |
| Pfad mit `..` bleibt beim Folge-Slice (Entscheidung 3) | erfüllt | kein Fall, keine Mutation; Punkt 20 sagt, dass der Stand über die Nicht-Identität sperrt |
| Punkt 13 als Regel statt Aufzählung | erfüllt | Wortlaut im Spec; die Temp-Repos der Gegenprobe tragen nur Makefile, Fragment, Skripte und Hook (`fixture()`) |
| Keine Lockerung (§3.6) | erfüllt | Hook sperrt zusätzlich bei fehlendem Stand, Konflikt, unlesbarer oder beschädigter Liste, `git`-Fehler; gleiche Menge im Vergleich wie der alte Hash (`--cached --others --exclude-standard`, Link mit Ziel, gelöschte Datei). `GATE_CHECKS` führt ein Gate mehr, keines weniger. Siehe V-333 zur Reihenfolge des Schleifenschutzes. |
| Aufruf-Voraussetzung (nacheinander, Punkt 8) | erfüllt | Hook-Kopf, Kopf der Gegenprobe; kein Fall (so vorgesehen) |
| Ausnahmeliste: sieben + `stop-hook-gegenprobe` = acht | erfüllt | Sensor 10; `SPEC-052` Punkt 2/4/13 zählen keine Zahl und bleiben richtig |

## 4. Plan gegen Diff (`AGENTS.md` §3.9) und Kommentare (§3.11)

- **Kopf, §1, §3, §6, §8:** folgen dem Diff. §3 nennt `record-gates.sh`, `working-tree-hash.sh`, Hook, `dateiliste.sh`, Gegenprobe samt Fragment, `gates-schnell.mk`, `gates-schnell-gegenprobe.sh`, `harness/README.md`, `harness/conventions.md` und die drei Commands (Aufzählungspunkt „Nachgezogen nach dem Review“); der Diff berührt genau diese und `MR-002` (nicht dieser Slice). Kein Produkt-Code, kein Handbuch, keine Tests unter `internal/` oder `test/`. §8 zählt 3 Liefer-Punkte, 1 Schicht; der Diff stimmt. `AGENTS.md` ändert der Diff nicht, und `grep` nach `gates-passed|working-tree-hash|record-gates|Content-Hash|Stempel` über `AGENTS.md`, `harness/*.md`, `.claude`, `docs/user`, `tools` findet keine Stelle mehr, die den alten Hash als Nachweis nennt (F-706 behoben).
- **F-703 bis F-712:** F-703 behoben (Fälle `p11-*`, V1/V2 rot). F-704 behoben und vom Architect entschieden (Fälle `p20-schleifenschutz-*`, V3). F-705 behoben (root-Lauf, Probe 6). F-706 behoben (drei Commands nennen Dateiliste, fehlenden Nachweis bei sauberem Baum und `make -i gates`). F-707 (a) behoben („Signal“ steht nicht mehr im Kommentar), (b) vom Architect entschieden (Punkt 22). F-708, F-709, F-711, F-712 als INFO belassen; F-710 vom Architect als Voraussetzung gefasst und im Kopf von `dateiliste.sh` genannt.
- **Kommentare:** Hook-Kopf: jede Zusage hat einen Fall (Konflikt `p17-konflikt-identisch`, fehlender Stand `p19-*`, Eintragsarten `p14-*`/`p18-*`, Schritte `p20-*`, Schleifenschutz `p20-schleifenschutz-*`); Kopf von `record-gates.sh`: `-i` `p21-i-*`, `-n`/`-q`/`-t` `p21-n`/`-q`/`-t`/`-ni`. Kopf der Gegenprobe nennt Grenzen (Attrappen, kein Fall für Punkt 8, root, kein `bash -e`) und passt zu dem, was ich gefahren habe. Kopf von `gates-schnell.mk` (nennt die Ausnahmeliste ohne Zahl) und das Fragment `stop-hook-gegenprobe.mk` sagen nur zu, was ein Fall prüft; die Aufzählung der Werkzeuge im Fragment (`awk`, `sed`, `cmp`, `comm`, …) ist eine Voraussetzungsliste, keine Zusage. Zwei Stellen sagen etwas weiter, siehe V-332 und V-334. Kein `//nolint`.
- **Commands:** `implement-slice.md`, `plan-welle.md`, `close-welle.md` stimmen mit dem Code (Probe 7, 8): Dateiliste als Nachweis, identisch → frei, fehlender Nachweis sperrt auch bei sauberem Baum, `make -i gates` schreibt keinen und entfernt einen älteren, `git mv` ist Umbenennung und sperrt.
- **Commits:** reine Inhalts-Commits, kein `git mv`; die Messages der Slice-Commits nennen Slice-Kennung und ADR-0038. `b397d6f` nennt nur ADR-0040 und gehört nicht zum Slice (siehe V-330).

## 5. Befunde

Kein HIGH, kein MEDIUM.

| ID | Schwere | Befund | Beleg |
|---|---|---|---|
| V-330 | LOW | Die Größe nach der Metrik des Planners ist am Stand der Belege richtig (1065 auf `d16d816`, Probe 11), am `HEAD` nicht: 1107, weil `b397d6f` (ADR-0040, `MR-002`) 42 Zeilen unter `harness/` in `9662e13..HEAD` legt. Wer die Metrik am `HEAD` fährt, bekommt eine Zahl, die nicht zum Slice gehört. Der Plan sagt „Commit `d16d816`“ nicht, nur die Formel mit `HEAD`. Die Closure nennt die Zahl am Stand des letzten Slice-Commits (`d16d816`, 1065 ; netto 893) oder schließt `harness/conventions/` aus. | Probe 11; Plan §7 Zeile 555 |
| V-331 | LOW | Die Entscheidung des Nutzers vom 2026-10-13 zur Fortsetzung bei 1065 steht nicht im Plan. §4 *Urteil zur Marke* gilt dem Stand 778; §7 nennt 920 und 1065 als Zahlen, aber kein Urteil, und der Reviewer hat nur bei 920 geurteilt (F-711). Die Closure-Notiz trägt die Entscheidung nach (Quelle: Auftrag an den Verifier), sonst hängt die Fortsetzung bei 1065 an einer Aussage außerhalb des Repos. | Plan §4 Zeilen 213 bis 228, §7 Zeilen 430 bis 433, 555 |
| V-332 | INFO | Die Erfolgsmeldung der Gegenprobe („66 Fälle … je Mutant rot“) zählt auch als root, obwohl dann ein Fall und ein Mutant übersprungen sind (Probe 6: zwei Zeilen „übersprungen“, danach „gruen — 66 Fälle“). Der Kopf der Gegenprobe sagt das Überspringen zu; die Schlusszeile ist enger zu lesen als „65 geprüft, einer laut übersprungen“. Kein Befund im Gewicht, solange die Übersprung-Zeilen auf stderr stehen. | Probe 6 |
| V-333 | INFO | Der Schleifenschutz steht jetzt vor `git rev-parse` und `source`: Bei `stop_hook_active: true` antwortet der Hook auch in einem Verzeichnis ohne Repo oder ohne `git` mit `approve`. Der alte Hook brach dort vor der Eingabe mit Fehlerstatus ohne Antwort ab. Die Wirkung ist keine Lockerung, weil Claude Code einen Abbruch ohne Antwort nicht als `block` liest, und der Architect hat es als akzeptierte Grenze in `SPEC-052` Punkt 20 gefasst (Fälle `p20-schleifenschutz-kein-repo`, `-ohne-git`). Ich nenne es, weil der Plan es „Bestand des alten Hooks“ nennt; für diese Teilmenge ist es die Fassung des Architect, nicht der Bestand. | Hook; Alter Hook `git show 9662e13:.claude/hooks/stop-require-gates.sh`; Plan §6 |
| V-334 | INFO | `SPEC-052` Punkt 11 und der Kopf von `dateiliste.sh` sagen: Fehlt bash 4, „bricht der Lauf ab“ (Punkt 20: `block`). Ein Fall dafür fehlt (so entschieden, nicht zugesagt). Gelesen: `local -A` in `dateiliste_pruefen` steht in der Position links von `||` im Hook (`dateiliste_pruefen "$stand" || voll …`), dort gilt `set -e` nicht; ein Fehlschlag von `local -A` bricht die Funktion dann nicht notwendig ab. Fail-closed bleibt, weil „frei“ ein `cmp` zwischen Stand und frisch erzeugter Liste voraussetzt, und eine beschädigte Liste nie mit der erzeugten übereinstimmt. Der Satz „bricht ab“ ist damit enger zu fassen („sperrt“); mangels bash 3 nicht gefahren. | `dateiliste.sh` Zeilen 7 bis 9, Hook letzte Zeilen |

## 6. Zu den Prüfpunkten des Auftrags

- **Größe:** 1065 gegen die Marke von 700 (Metrik des Planners, am Stand `d16d816` gemessen, Probe 11; Verlauf 778, 920, 1065). 769 davon sind die tabellengetriebene Gegenprobe (72 Prozent). Mein Urteil: **prüfbar in einer Sitzung, am Rand**. Der Hook (65/45 Zeilen), `dateiliste.sh` (60), `record-gates.sh` (34/6), `working-tree-hash.sh` und die Commands sind in einem Durchgang lesbar; die Gegenprobe liest man gegen die Mutationstabelle, und genau so habe ich sie geprüft (neun eigene Mutanten, Hook von Hand). Dass die erste Runde zwei MEDIUM-Lücken in der Gegenprobe fand (F-703, F-704), zeigt, dass die Prüfung ihrer Breite an die Grenze kam; eine weitere Ausweitung in diesem Slice würde ich nicht mehr als prüfbar tragen. Die Schätzung des Architect (440 bis 620) lag um 445 daneben, im Verlauf 158, 300 und 445 darunter: derselbe Fehler wie in `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`.
- **Metrik:** stimmt am Stand der Belege (Probe 11); am `HEAD` +42 durch `MR-002` (V-330).
- **Folgt `SPEC-052` Teil 2 dem Code:** ja. Eingangssatz (Punkte, die kein Fall nennt, sind nicht zugesagt) trifft: die Gegenprobe nennt 1, 10, 11, 14, 16 bis 22; die Punkte 9, 12, 13, 15 prüft der Folge-Slice und stehen weiterhin ohne Fall. Punkt 11 (bash ab 4), 20 (Grenze des Schleifenschutzes), 21, 22 („leere Umgebung“) stimmen mit meinen Proben.
- **Keine Lockerung (§3.6):** bestätigt (Abschnitt 3, V-333).

## 7. Offen und nicht gefahren

- Offen, weil die Closure sie führt: Closure-Notiz mit Lerneintrag, Beobachtungs-Register (Evidence zu `groesse-erst-in-der-pruefung-des-architect-sichtbar`), Ausgang jedes Risikos aus §6, Name der Dateiliste im Folge-Slice mit Nachzählung, drei Paarungen, `git mv` nach `done/`.
- Nicht gefahren: `make gates` (auf Anweisung; Stand laut §7 `b397d6f`, der Hook antwortet im Arbeitsbaum nach den Plan-Commits wieder `block`, wie es sein muss); `make test-integration`; `make lint-gegenprobe` einzeln; `make a-check-negativ`; `make abdeckung-check`; `make baseline-verify`. Der Hook wurde in Temp-Repos gefahren, nicht in der Umgebung von Claude Code (das Risiko ist im Plan *weiter offen*). Ein Lauf unter bash 3 oder mit BSD-Werkzeugen fehlt (V-334). Nicht jede Zeile der Mutationstabellen in §7 ist nachgefahren (neun eigene Mutationen, plus der Lauf des Skripts selbst mit seinen 36 Mutanten). Die Fälle als root habe ich im Container gefahren (Probe 6); `p20-stand-unlesbar` als Nutzer lief im Aufruf `make stop-hook-gegenprobe` (Nutzer 1000).

## Verdikt

Die DoD-Punkte 1 bis 3, Review und Mutationen sind erfüllt, die Entscheidungs-Konformität ist gegeben, der Hook sperrt nirgends weniger als der alte, und `SPEC-052` Teil 2 folgt dem Code. `make gates` ist laut Aufrufer grün, nicht von mir gefahren. Es bleiben die Closure-Pflichten und zwei Hinweise für die Closure-Notiz (V-330 Messstand der Größe, V-331 Entscheidung des Nutzers zur Fortsetzung bei 1065). Kein Rückgabegrund an den Implementer oder den Architect.
