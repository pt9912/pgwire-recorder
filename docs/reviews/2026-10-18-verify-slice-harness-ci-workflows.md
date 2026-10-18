# Verifikation slice-harness-ci-workflows

Stand: Zweig `rev-ci-workflows`, Kopf `e9219f59`, Basis des Plans (Metrik) `a98d0b9e`, Basis der Änderung `79b82a0e`.
Maßstab: [Plan](../plan/planning/done/slice-harness-ci-workflows.md) §1, §2, §3, §6, §7,
[SPEC-063](../../harness/vertraege/SPEC-063-workflows-check.md), [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md), `AGENTS.md`.
Eingang: DoD-Bestätigung und Belege des Implementers (§7, *Belege des Implementers* und *Findings des Reviews*), Review-Report
[2026-10-18-review-slice-harness-ci-workflows.md](2026-10-18-review-slice-harness-ci-workflows.md) (0 HIGH, 3 MEDIUM, 4 LOW, 3 INFO).
Host-Läufe als Benutzer `db` (uid 1000, kein root), Docker am Host. Nicht gefahren, wie beauftragt: `make gates`, die volle Gegenprobe.

## Ergebnis je Punkt

| # | Punkt | Ergebnis |
|---|---|---|
| 1 | `make docs-check` am Kopf | 0 Befunde (1160 Dateien); der Link-Befund des Reports (Zeile 15) ist mit `e9219f59` behoben |
| 2 | `make workflows-check` am Arbeitsbaum | Exit 0, Meldung „keine Workflow-Dateien unter .github/workflows“ (Randform 1) |
| 3 | `make kopf-check abdeckung-check meldungskatalog-check` | grün |
| 4 | `make gates-schnell-gegenprobe` | grün, 1 min 05 s |
| 5 | Kette: `make -n gates`, `make -n gegenproben`, `make -n gates-schnell` | `workflows-check` hängt an `GATE_CHECKS` und läuft in der schnellen Kette; `workflows-check-gegenprobe` steht in `GEGENPROBEN_CONTAINER` |
| 6 | 25 Mutationen einzeln per `GP_MUT` und `GP_ENTDECKE=1`, in drei Läufen (dazu der Lauf unter root, Punkt 7) | alle mit genau der erwarteten roten Menge, jeder Lauf grün gemeldet (siehe unten) |
| 7 | Fälle unter Root (`unshare -r`, `id -u` = 0) | grün, 9 Läufe; die sechs Namen der entfallenen Fälle und Mutationen stehen auf stderr (F-3 bestätigt) |
| 8 | Unabhängige Probe des Skripts vor und nach den Fixes | siehe „Zu 6 und 8“ |
| 9 | Metrik `a98d0b9e..e9219f59 -- .github tools harness Makefile` | 899 (wie §7); von `79b82a0e` aus 891 |
| 10 | Anzahl und Namen der Mutationen gegen §7 | 99 `mut`-Zeilen, 99 Namen in der Tabelle, `diff` leer; 98 Fall-Funktionen (wie §7) |

### Zu 6 und 8 (Mutationen, Sonden)

Gefahren per `GP_MUT` (der Mechanismus der Gegenprobe verlangt je Mutation die benannte Menge, ein Unterschied färbt sie rot):

- zu F-1: `m-p8-dir-ohne-block` (`p8-dir-danach p8-dir-davor`), `m-p8-ohne-ruecksetzung` (`p8-dir-danach`), `m-p8-modules-muster` (`p8-modules-praefix p8-modules-suffix`).
- zu F-2: `m-p9-verzeichnis-r-fehlt` (`p9-verzeichnis-ohne-r`), `m-p9-verzeichnis-x-fehlt` (`p9-verzeichnis-ohne-x`).
- zu F-3: `m-p9-nicht-lesbar` bleibt als Nicht-root rot (`p9-nicht-lesbar`); als (Namespace-)root entfallen die drei Fälle und die Mutationen m-p9-nicht-lesbar, m-p9-verzeichnis-r-fehlt, m-p9-verzeichnis-x-fehlt, die Gegenprobe endet grün und nennt sie.
- zu F-4: `m-p5-regeln-laufen` (`p5-form-ohne-regeln`).
- zu F-5: `m-p3-al-zeilen-veraendert` (`p3-exit1`).
- zu F-10: `m-p7-fuehrende-null-angenommen` (`p7-fuehrende-null`), `m-p7-vorzeichen-angenommen` (`p7-vorzeichen`), `m-p7-dezimal-angenommen` (`p7-dezimal`), `m-p1-symlink-andere-befund` und `m-p1-symlink-andere-still` (je `p1-symlink-andere-endung`).
- Aufrufart: `m-p10-plus` (`p10-makeflags p10-n p10-q p10-t`).
- Punkt 9 und 3: `m-p9-datei-vor-al`, `m-p9-global-vor-datei` (je `p9-reihenfolge`), `m-p9-ohne-sortierung` (`p9-alle-befunde`), `m-p9-erster-befund`, `m-p9-umgebung-exit1`, `m-p9-eigene-vor-umgebung`, `m-p9-summe-ohne-al`, `m-p3-digest-kurz` (`p3-digest-form`).
- Echte Images: `m-real-attrappe-image` (`real-ausdruck`), `m-real-modules` (`real-pins`); die unveränderte Kopie (`ok`, mit den drei echten Fällen) grün.
- Breite Mutationen: `m-p6-job-pflicht`, `m-p9-summe-nur-al` färben den benannten Kern und lassen die genannten Fälle grün.

**Grund der Röte, unabhängig vom Mechanismus** (das Skript vor dem Fix, `1bdaaefe^`, gegen das Skript am Kopf, an Bäumen im Scratchpad mit Attrappe von `docker`):
- F-10: `timeout-minutes: 010` — vorher „1 Datei(en) geprüft, kein Befund“ (Exit 0), jetzt `timeout-wert … 010 …` (Exit 1). Der Fix ändert das Verhalten; die Fälle `p7-fuehrende-null`, `p7-vorzeichen`, `p7-dezimal` hängen daran.
- F-1 und F-2: dieselbe Probe (`.d-check.yml` mit `dir:` unter einem fremden Block; Verzeichnis mit Modus 311) liefert vor und nach dem Fix dasselbe richtige Ergebnis (Befund `pins-ungeprueft`, Exit 2 „nicht lesbar“). Die Lücke lag in der Gegenprobe, nicht im Skript; die Mutationen oben färben genau die neuen Fälle, also schließt der Fix die Lücke aus dem richtigen Grund.

## DoD-Konformität

- **Liefer-Punkt 1 (Gate):** bestätigt. Fragment `harness/mk/workflows.mk`, Skript `tools/harness/workflows-check.sh`; Rezept ohne `$(MAKE)` und `+`; Image `rhysd/actionlint@sha256:…` mit Tag als Kommentar, `--network none`, `:ro`; Ausgabe `Pfad:Zeile: Regel: Text`; Exit 0/1/2; Eintrag in `GATE_CHECKS`; `.d-check.yml` unverändert (der Eintrag gehört `slice-harness-ci`, dort in §1, DoD, §3 und §8 mit der Kennung des Gebers eingetragen und nachgezählt). Randformen in SPEC-063 vor dem ersten Code-Commit (`7a66cc39`; Nachträge `2890e31f`, `1b7d214a`, `77de7875` sind Entscheidungen des Architect zu Rückgaben und Review).
- **Liefer-Punkt 2 (Gegenprobe):** bestätigt. 99 Mutationen, 98 Fälle; je Aufrufart `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` ein Fall; Stichproben (25 Mutationen) rot mit der benannten Menge; unveränderte Kopie grün, einschließlich der drei echten Fälle.
- **Liefer-Punkt 3 (Gate-Index):** bestätigt. Zwei Zeilen in `harness/README.md` §Sensors; die Prosa nennt die Regel an den Schreiber, zeigt auf SPEC-063 und `.d-check.yml`; die Aussagen darin (Exit 0 ohne Datei, `:ro`, `--network none`, drei breite Mutationen, root) deckt je ein Fall oder ein Lauf oben.
- **`make gates` grün:** **am Endstand nicht belegt** (V-1).
- Review durchgeführt: Report liegt vor; alle zehn Findings haben in §7 einen Ausgang, F-1 bis F-8 und F-10 nachgeprüft (F-8: Zahl stimmt, siehe Metrik; F-9 bleibt bewusst).
- Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen: noch offen (Closure-Schritte, nicht Teil dieser Verifikation).

## Plan gegen Code (`AGENTS.md` §3.9)

- §1 (Ziel, Abgrenzung): Das Gate tut, was §1 sagt; keiner der Ausschlüsse ist verletzt (kein Eintrag in `.d-check.yml`, keine Workflow-Datei, keine Pin-Gültigkeit, keine Actions, kein Netz im Lauf).
- §3 (Plan): vier Zeilen, vier Artefakte, kein weiteres im Diff (`git diff --stat 79b82a0e..HEAD`: Fragment, zwei Skripte, `harness/README.md`, SPEC-063, Plan, Review-Report).
- §6: Einleitung („vom Architect entschieden“) und DoD-Wortlaut („Pfad, Zeile und Regel“) sind nachgezogen (F-7); die zwölf Randformen stehen mit Verweis auf die Punkte von SPEC-063.
- Kopf: `Berührte Spec-Stellen` führt SPEC-063, SPEC-060 Teil 1, SPEC-052; `Verantwortlich` trägt den Zweig.
- Sendung an `slice-harness-ci`: getragen (§1 Abgrenzung, DoD Liefer-Punkt 1, §3 Tabellenzeile, §8 Zählung mit der Kennung des Gebers).
- **Nicht nachgezogen:** §7 *Läufe am Endstand* und *Offene Rückgabe* (V-2).

## Zusagen in README, Kommentaren, Hilfetexten (`AGENTS.md` §3.11)

Gelesen und gegen Fall oder Lauf gehalten: Kopf von `workflows-check.sh` (Aufruf, Umgebung, Ausgang 0/1/2), Kommentar und Hilfetexte in `workflows.mk` („prüft … faehrt drei Faelle gegen die echten Images“; gedeckt durch `real-*`), Kopf der Gegenprobe, Zeilen in `harness/README.md`. Keine Zusage ohne Prüfung gefunden, mit einer Randbemerkung: die Werkzeugliste des Fragments („bash, awk, find, sort und Docker“) lässt `readlink`, `mktemp` und `sed` des Skripts weg; die Gegenprobe prüft die vorhandenen Werkzeuge (`for t in make awk find sort cksum mktemp`), das Gate nicht (V-4, gering).

## Risiken aus §6 und ihr Ausgang für die Closure

| Risiko | Ausgang heute im Plan | Was bei der Closure fehlt |
|---|---|---|
| 1 Gate prüft Klasse, nicht Lauf | „vorgesehen: weiter offen → neues `BEO-REPO/workflow-gate-gruen-lauf-rot`, wenn … sonst entfallen“ | eine der drei Formen wählen; das Register hat den Eintrag nicht (`observations/BEO-REPO/` führt keinen `workflow-…`) |
| 2 Pin-Check hängt am gepinnten d-check | „vorgesehen: weiter offen → an `slice-harness-ci` …; sonst entfallen“ | wählen; die Sendung in `slice-harness-ci` steht, die Wahl *eingetreten/weiter offen* trägt nur die Closure |
| 3 `shellcheck` auf `run:`-Blöcke | „vorgesehen: entfallen, wenn SPEC-063 die Einstufung festlegt und die Gegenprobe sie zeigt“ | Bedingung ist erfüllt (Punkt 4, Fall `p3-e8`, drei Mutationen `m-p3-e8-*`); als *entfallen* mit dieser Begründung eintragen |
| 4 Gegenprobe braucht das Image | „vorgesehen: entfallen, weil in `GEGENPROBEN_CONTAINER`“ | als *entfallen* eintragen; Anmerkung: `make workflows-check` selbst zieht das Image in `make gates` (Pull-Limit des Docker Hub, vom Vertrag und von [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md) benannt) |

Alle vier tragen nur einen *vorgesehenen* Ausgang; §7 sagt „Wird bei Closure gefüllt“. Der Slice darf erst nach `done/`, wenn jedes genau einen Ausgang trägt. Beobachtungen/Lerneintrag: die Finding-Klassen des Reviews gehen als Zähler in das Register (siehe Übergabe des Reports).

## Größenmetrik (`AGENTS.md` §3.14)

`git diff --numstat a98d0b9e..e9219f59 -- .github tools harness Makefile`: Einfügungen 877, Löschungen 22, Summe **899**.

| Datei | Einfügungen/Löschungen |
|---|---|
| `tools/harness/workflows-check-gegenprobe.sh` | 589 / 0 |
| `tools/harness/workflows-check.sh` | 240 / 0 |
| `harness/mk/workflows.mk` | 24 / 0 |
| `harness/README.md` | 2 / 0 |
| `harness/vertraege/SPEC-063-workflows-check.md` | 22 / 17 |
| `harness/vertraege/SPEC-064-ci.md` | 1 / 4 |

Marke 900 um eine Zeile unterschritten, Meldeschwelle 700 überschritten (wie §7 sagt). Nicht zum Slice gehören 8 Zeilen (SPEC-064 mit 5 und 3 Zeilen von SPEC-063): der Commit `28775212` des Architect (Stand der Entscheidung nach Annahme der [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md)) liegt vor dem Merge `79b82a0e`, aber hinter der Basis `a98d0b9e` des Plans. Von `79b82a0e` aus sind es 891. Beide Zahlen unter der Marke; keine Rückführung. Der Review-Report und §7 liegen außerhalb der gemessenen Pfade.

## Befunde

- **V-1 (mittel, vor der Closure): kein Beleg für `make gates` grün am Endstand.** Die DoD verlangt den vollen Lauf (Stop-Hook, weil `harness/` und `tools/` ändern). §7 belegt nur `make -k gates` mit rotem `docs-check`; der Nachweis des vollen Laufs wurde damit nicht geschrieben, und seitdem sind `harness/README.md` (`6f4058fe`) und der Report (`e9219f59`) hinzugekommen. Ich habe `make docs-check`, `kopf-check`, `abdeckung-check`, `meldungskatalog-check`, `workflows-check` und `gates-schnell-gegenprobe` grün gesehen, den vollen Lauf nicht gefahren (Auftrag). Der Implementer fährt `make gates` am Endstand einmal und trägt das Ergebnis ein.
- **V-2 (gering): §7 folgt dem Stand nicht (`AGENTS.md` §3.9).** Die Zeilen *Läufe am Endstand* (`make gates` endet an `docs-check`, drei Befunde `id-unlinked`, `docs-check`: 3 Befunde) und der Absatz *Offene Rückgabe* beschreiben einen Zustand, den `e9219f59` behoben hat; `docs-check` ist am Kopf 0 Befunde. Die Zeile „Stand … letzter Code-Commit `1bdaaefe`“ nennt zudem nicht den Kopf. Vor der Closure nachziehen.
- **V-3 (gering): Zeilenende CRLF.** Eine Workflow-Datei mit `\r\n` ergibt `permissions-form … gelesen: {}` (der gelesene Wert sieht richtig aus, der Befund täuscht) (das Skript liest `\r` als Teil des Werts; `actionlint` nimmt die Datei an). Falsch-positiv und damit sicher, aber die Meldung führt in die Irre. SPEC-063 Punkt 5 nennt das Zeilenende nicht (`AGENTS.md` §3.12, Randform Zeilenende bei Plattformabhängigkeit); Sache des Architect: Grenze nennen oder `\r` als Befund `form` einstufen. Kein Defekt gegen den Vertrag, wie er steht.
- **V-4 (gering): Werkzeugliste im Kopf von `workflows.mk`** nennt `readlink`, `mktemp` und `sed` (GNU) nicht; Hinweis, kein Blocker.
- **V-5 (INFO): zwei Basen der Metrik.** Der Plan misst von `a98d0b9e` (899, mit 8 fremden Zeilen), die Beauftragung von `79b82a0e` (891). Die Marke trägt in beiden Lesarten; für künftige Pläne die Basis hinter den Commit des Architect legen, der den Stand der Entscheidung setzt.

Negativbefund (geprüft, ohne Befund): Reihenfolge der Gründe, Exit-Codes, Aufruf von `actionlint` (E1 bis E8), Symlink-Fälle, Fälle `p2-*`, Normalform der Zeitgrenze, Bindung (a) bis (c), root-Verhalten, Aufrufarten (`-n -q -t -i MAKEFLAGS`), Schreibverhalten (nichts im Baum), `GATE_CHECKS`, `GEGENPROBEN_CONTAINER`, Ausnahmeliste der schnellen Kette (`workflows-check` nicht darin), Sendung an `slice-harness-ci`, Tabelle der Mutationen in §7 gegen das Skript (99 gegen 99, Namen gleich), Zahlen 100 Läufe/98 Fälle (Anzahl der Funktionen bestätigt; die Läufe selbst nicht neu gefahren).

## Urteil

**Bestanden**, mit Auflagen vor der Closure: V-1 (mittel) muss belegt sein, V-2 (gering) nachgezogen, die vier Risiken brauchen je genau einen Ausgang (Tabelle oben). V-3 bis V-5 blockieren nicht; V-3 geht als Randform-Frage an den Architect. Kein HIGH, kein Befund gegen Liefer-Punkt 1 bis 3.
