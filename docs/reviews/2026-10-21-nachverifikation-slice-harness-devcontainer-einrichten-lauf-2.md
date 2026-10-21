# Nachverifikation, zweiter Lauf: slice-harness-devcontainer-einrichten — 2026-10-21

**Rolle:** Verifier (Modul 11), frischer Kontext, enger Auftrag (Zeitbudget 20 Minuten). **Gegenstand:** Arbeitsbaum `/Development/pgwire-recorder-int`, Zweig `int-kette`, Stand **`2c73122a`** (Merge von `e772fc39`), sauber. **Vorgänger:** `2026-10-21-nachverifikation-slice-harness-devcontainer-einrichten` (Stand `a8fa0f43`, HIGH N-1, LOW N-2 und N-3, INFO N-4 und N-5; Lauf-Beleg, unverändert).

**Maßstab:** DoD-Punkt 4 des Plans `docs/plan/planning/in-progress/slice-harness-devcontainer-einrichten.md` („`make container-einrichten-gegenprobe` am Endstand in einem Zug grün, mit ihrer Zeit im Bericht"), [`SPEC-060`](../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 3, 16 und 17, [`SPEC-069`](../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 17, [ADR-0055](../plan/adr/0055-gegenproben-haben-ein-zeitbudget.md), [ADR-0056](../plan/adr/0056-groessen-metrik-misst-den-gelieferten-diff.md).

## Ergebnis in einem Satz

**Der volle Lauf der Gegenprobe ist am Stand `2c73122a` in einem Zug grün: Exit 0, 227 Fälle, Zeit 183 s von höchstens 600 s, Last 9.95 am Ende (3.72 zum Start); N-1 ist behoben, N-2 bis N-4 stehen im Plan und im Skript.** Der DoD-Punkt „am Endstand in einem Zug grün" ist erfüllt. Kein HIGH und kein MEDIUM offen.

## Selbst gefahren

Lauf in `pgr-gp` (Benutzer `vscode`, Port 2225, bash 5.2.37, 20 Kerne) im Klon `/workspaces/pgr-gp/nv-pruefung` aus einem `git bundle` von `2c73122a` (ohne Push, ohne Austausch-Repo; `git rev-parse HEAD` im Klon gleich dem Stand; `.harness/state` fehlte: `ls -a .harness` nannte `.gitattributes`, `.gitignore`, `baseline`, `skills`). Zeiten UTC, Last = erster Wert von `/proc/loadavg` (im Container gleich dem des Hosts). Start im Vordergrund der `ssh`-Sitzung mit `env --default-signal=INT,QUIT` davor ([`SPEC-060`](../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 17); ein Lauf, kein Startfehler, keine Wiederholung.

| # | Befehl | Stand | Start (UTC), Last zum Start | Ergebnis |
|---|---|---|---|---|
| 1 | `make container-einrichten-gegenprobe` (voller Lauf in einem Zug, frischer Klon ohne `.harness/state`) | `2c73122a` | 18:30:44, **3.72** (zuvor gelesen 3.87; ein Warten war nicht nötig) | **Exit 0, grün.** Ende 18:33:47, Wanduhr **183 s**. Zeile `container-einrichten-gegenprobe: Zeit 183 s von höchstens 600 s, Last 9.95` (die Last am Ende trägt der Lauf selbst, `GP_PARALLEL` 4). Schlusszeile: `gruen — 227 Fälle zu SPEC-069 Punkt 1 bis 17, darunter die zu SPEC-060 Punkt 16 (9); je Mutation genau die benannten roten Fälle, ausgenommen die ohne roten Fall (äquivalent, §7 des Plans): m-p2-erster-verstoss m-p17-signal-term`. Keine Zeile `ROT`, keine Abweichung (am Stand `a8fa0f43` stand hier `ROT — Mutation m-p7-remote-add …`). Das übrige auf stderr (`error: No such remote`, `fatal: cannot change to`, `Terminated`) ist das Rauschen erwartet roter Fälle. |

Gelesen, nicht gefahren: Protokoll des Koordinators `gates-int-endstand2.log`: `2c73122a` im Kopf, **`exit=0`, `real 4m19,926s`**, `a-check` 0 Befunde, `baseline-verify` OK (54 Dateien), `test-container` grün, `d-check: 1344 Datei(en) geprüft, 0 Befund(e)`, `workflows-check` kein Befund; `.harness/state/gates-passed.diffsha` (`c142e407…`) ist gleich `tools/harness/working-tree-hash.sh` des Baums; `git status --short` leer.

## N-1 bis N-4 im Skript und im Plan (gelesen, keine eigenen Läufe)

| Befund | Stand am Endstand | Beleg |
|---|---|---|
| **N-1 (HIGH)** `m-p7-remote-add` färbt `p7-url-port` | **behoben** | `tools/harness/container-einrichten-gegenprobe.sh` Zeile 831: die rote Menge der Mutation nennt jetzt `p7-url-port` (hinter `p7-abgewiesen`); `git diff --numstat a8fa0f43 HEAD -- tools`: **3 Einfügungen, 3 Löschungen**, das sind die Zeile 831 und zwei Zeilen des Kopfkommentars (Voraussetzungsliste); mein Lauf 1 hat keine Abweichung. Der Plan benennt in *Nacharbeit nach der Nachverifikation* die Ursache (Gruppe des neuen Falls, 55 Mutationen mechanisch bestimmt und eng gefahren, Exit 0, 49 s) und eine Beobachtung für die Closure. |
| **N-2 (LOW)** Plan folgt dem Endstand nicht | **behoben** | Kopf: „`SPEC-060` Punkte 2, 3, 8, 16 und 17 (…Punkt 17 … und sein Verweis in Punkt 8 — beide vom Architect am 2026-10-21 nach F-V10 eingetragen)"; §3 Zeile `SPEC-060` nennt Punkt 17, die Zeile Z6 und den Verweis in Punkt 8 (der Satz „kein anderer Punkt ändert sich" steht jetzt hinter dieser Aufzählung und trifft); §3 Zeile `harness/README.md` nennt den Satz im Vorspann von §Sensors; der Satz „§1, §3 und §6 bleiben" am Kopf von *Nacharbeit nach der Verifikation* ist berichtigt („§1 bleibt; … §3 … und §6 … sind … nachgezogen"). Offen bleibt als INFO die fehlende Zeile für `harness/vertraege/README.md` in §3 (der Index des Architect, eine Zeile). |
| **N-3 (LOW)** §7 trägt `make gates` und die Größe nur für `7f2dd9ed` | **behoben für `a8fa0f43`, für `2c73122a` noch nicht (siehe offene Punkte)** | §7 nennt den roten Lauf (Exit 2, 200 s, Last 3.04) und `make gates` auf `a8fa0f43` (Exit 0, 7 min 55 s); die Größe ist mit **1407** neu genannt, und ich messe an `2c73122a` dieselben **1407** (Einfügungen 1387, Löschungen 20). |
| **N-4 (INFO)** `awk` fehlt in der Voraussetzungsliste | **behoben** | Kopf des Skripts („bash ≥ 4.4 …, git, make, sed, awk, grep, find und die GNU-Werkzeuge von coreutils (dazu diff mit GP_ENTDECKE=1)"), Kommentar des Fragments `harness/mk/container-einrichten-gegenprobe.mk` („bash, git, make, sed, awk, grep, find und die GNU-Werkzeuge") und die Zeile der Gegenprobe in `harness/README.md` („bash, git, make, sed, awk, grep, find …") nennen es; `diff` steht allein im Kopf des Skripts (es wird nur mit `GP_ENTDECKE=1` gerufen, Zeile 208). |
| **N-5 (INFO)** Flake von `m-p17-signal-int` | **keine neue Messung** (Auftrag) | Stand aus dem ersten Nachverifikationsbericht: 0 von 15 bei Last bis 6.3. Der volle Lauf 1 dieses Berichts ist unter Last 3.72 bis 9.95 grün, einschließlich `p17-budget-sigint` (kein Fall von `m-p17-*` abweichend). |

## DoD-Stand am Endstand `2c73122a`

| Punkt | Verdikt |
|---|---|
| 1: Ziel `container-einrichten`, Nachweis | **erfüllt, unverändert.** Skript `tools/harness/container-einrichten.sh` (Blob `be225508…`) und Fragment `harness/mk/container-einrichten.mk` (Blob `319dafa5…`) sind an `HEAD` byteweise die des Nachweises am Stand `afe27802` (`git rev-parse` gelesen). |
| 2: Gegenprobe | **erfüllt:** 227 Fälle, jede Mutation mit genau der benannten roten Menge (Lauf 1), Zeit 183 s von höchstens 600 s. |
| 3: `make gates` grün | **erfüllt am Endstand** laut Protokoll (`exit=0`, 4 min 20 s, Hash gleich dem Stand); nicht von mir gefahren. |
| **4: `make container-einrichten-gegenprobe` am Endstand in einem Zug grün, Zeit im Bericht, frischer Container ohne `make gates`** | **ERFÜLLT** (Lauf 1: Exit 0, 227 Fälle, **183 s von höchstens 600 s, Last 9.95**; Klon ohne `.harness/state`). |
| 5: Review durchgeführt | erfüllt, unverändert (Review, Verifikation und zwei Nachverifikationen liegen vor). |
| 6 bis 9: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, drei Paarungen | **offen, gehören zur Closure** (Planner). |

## Offene Punkte (keine Befunde gegen den Code; nichts davon blockiert die Closure, jeder gehört in die Closure-Notiz oder §7)

1. **§7 trägt die beiden Belege am Endstand `2c73122a` noch nicht:** den vollen Lauf der Gegenprobe (183 s, Last 9.95, 227 Fälle, dieser Bericht) und `make gates` (Exit 0, 4 min 20 s, Protokoll des Koordinators). Der Plan sagt selbst, dass beides der Nachweis für den neuen Stand ist und nicht von der Nacharbeit gefahren wird; sie gehören bei der Closure in §7 (`slice-harness-blackbox-kern`: der Beleg steht in §7).
2. **Größe (F-V11, an den Planner):** **1407** nach [ADR-0056](../plan/adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) gegen die Marke 900 (Meldeschwelle 700), unverändert gegen den Stand `a8fa0f43`; Ausgang des Risikos *Größe über Schätzung*: eingetreten, Befund an den Planner (Entscheidung des Nutzers vom 2026-10-21: kein weiterer Schnitt).
3. **INFO:** `harness/vertraege/README.md` (eine Zeile des Architect) fehlt in der Tabelle von §3.
4. **Closure-Pflichten:** Closure-Notiz mit Steering-Loop-Lerneintrag; Register (Finding-Klassen für den Zähler: *die Erwartung einer anderen Mutation* — der Plan nennt die Form und `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`; Plan folgt der Korrektur nicht, `BEO-REPO/plan-folgt-korrektur-nicht`; Last-Flake, `BEO-REPO/test-faellt-unter-cpu-last-aus`); Risiko-Ausgänge aus §6; die drei Paarungen; der Übergang nach `done/` als reiner `git mv`.

## Verdikt

**DoD-konform am Endstand `2c73122a`.** Der DoD-Punkt „am Endstand in einem Zug grün" ist erfüllt: Exit 0, 227 Fälle, 183 s von höchstens 600 s, Last 9.95 am Ende, im frischen Klon ohne `.harness/state`, ohne Push. N-1 bis N-4 stehen im Skript (Zeile 831, Kopf) und im Plan (Kopf, §3, §7). Was fehlt, ist Closure-Arbeit (offene Punkte 1 bis 4), nichts am Code.

**Nicht geprüft:** der Nachweis nach Punkt 18 (am Skript unverändert, Blobs gleich), die Fahrt unter fremdem Benutzer und Last, die Sonde zu `m-p17-signal-term` (Text unverändert), eine erneute Flake-Messung (nicht beauftragt), `make docs-check` über diese Datei (fährt der Koordinator).
