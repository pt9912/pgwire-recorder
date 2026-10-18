# ADR-0048: Die Gegenproben verlassen die Gate-Kette, zwei Sammelziele fahren sie, eine Pflicht und ein eigener Container tragen den Beleg

**Status:** Accepted

**Datum:** 2026-10-18

**Autor:** Architect; angenommen durch Entscheidung des Nutzers am 2026-10-18 (Richtung im Dialog, danach „Annehmen“ zu dieser Fassung; `make gates-schnell` bleibt vorerst, der Wortlaut der Pflicht in `AGENTS.md` §6 ist bestätigt).

**Bezug:** [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0038](0038-gate-kette-voll-und-schnell.md), [ADR-0047](0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md), [ADR-0033](0033-gate-nachweise-in-der-abdeckung.md)

**Schärft:** [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0038](0038-gate-kette-voll-und-schnell.md) hat gemessen: `lint-gegenprobe` allein 150 s bei einer Last von im Mittel 57 auf 20 Kernen, `make gates` insgesamt etwa 4,5 Minuten, und seither sind weitere Gegenproben mit Docker dazugekommen. Die Positivliste dieser ADR hilft nur bei Änderungen an Planung, Reports und Dokumentation; jede Änderung an Code, Gate, Skript oder Fragment verlangt den vollen Lauf mit allen Gegenproben, und die Last hält die anderen Läufe auf demselben Host auf. Die Gegenproben prüfen die Gates, nicht das Produkt.

[ADR-0038](0038-gate-kette-voll-und-schnell.md) hat „Gegenproben aus `make gates` herausnehmen, nur in einem eigenen Ziel“ als Alternative D verworfen: eine Lockerung von `make gates` (`AGENTS.md` §3.6), und „ein Gate, das niemand fährt, ist eine Behauptung“ (`modul-13-quality-gates.md`, Harness-Lüge). Der Nutzer hat diese Alternative am 2026-10-18 gewählt, mit einem Gegenmittel. Diese ADR schreibt die Abwägung offen, statt sie zu übergehen.

## Entscheidung

Wir empfehlen **Option C: alle Gegenproben und `a-check-negativ` verlassen die Gate-Kette; zwei Sammelziele fahren sie; eine Pflicht zum Fahren und ein eigener Container tragen den Beleg.**

Vom Nutzer entschieden, im Dialog vom 2026-10-18 (Wortlaut vom Koordinator übermittelt): „Das dauert zu lange und ich verstehe auch nicht, warum das nötig ist, wenn sich nichts geändert hat“, „Der schnelle fix besteht darin es aus make gates herauszunehmen“, „Alle Gegenproben raus“, „Man kann es explizit in einem eigenen Dev-Container fahren, dann werden die anderen nicht von der eigentlichen Arbeit abgehalten“ und „Wir können Gegenproben aufteilen. Einen speziellen, der nur im Host (u-boot und devcontainer) läuft (zwei make targets)“. Daraus folgt:

1. Jedes Gate, dessen Name auf `-gegenprobe` endet, und `a-check-negativ` verlassen `GATE_CHECKS`. Jedes steht in genau einer von zwei Listen: der für den Dev-Container oder der für den Host (`u-boot`, `devcontainer`). Je ein Sammelziel fährt eine Liste nacheinander und sammelt rote Einzelziele.
2. Die Pflicht zum Fahren: bei einer Änderung am geprüften Werkzeug oder Gate, an der Gegenprobe oder ihrem Vertrag die betroffene Gegenprobe; am Ende jeder Welle und vor dem Release beide Sammelziele (Entscheidung des Nutzers vom 2026-10-18, so vom Koordinator übermittelt; der Wortlaut liegt dem Architect nicht vor).
3. Der automatische Lauf liegt nur in einem eigenen Gegenproben-Container, der auf dem eingelieferten Stand eines Zweigs fährt, während am Host und im Arbeitscontainer weitergearbeitet wird.

**Entscheidung des Nutzers vom 2026-10-18 (Empfehlung des Architect angenommen):**

- Die Namen `gegenproben` (Dev-Container) und `gegenproben-host`, die Variablen der Listen und die Namensregel stehen in [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md).
- Kein Stand der Gegenproben und kein Hook, der ihn liest. Der Stop-Hook glaubt weiter allein den Nachweis von `make gates`, und der schließt die Gegenproben nicht ein.
- `make gates-schnell` bleibt neben `make gates`. Beide Ketten fahren nach dem Umbau keine Gegenprobe; sie unterscheiden sich durch `test-container`, `test-integration` und `test-postgres-versionen`. Der Hook wählt für inerte Änderungen die schnelle Kette und spart damit diese drei; ob das die zweite Kette trägt, zeigt die gemessene Dauer.
- Die Wahl der Liste trifft der Schreiber einer Gegenprobe mit einer Zeile im Fragment; die Gegenprobe der Kette prüft die Eindeutigkeit, nicht die Wahl.

**Von [ADR-0038](0038-gate-kette-voll-und-schnell.md) ersetzt** sind: Punkt 2 der Entscheidung (die schnelle Kette fährt „nur die Gates ohne Gegenproben und ohne Integrationstests“: ab jetzt fahren beide Ketten keine Gegenprobe, die schnelle zusätzlich keines der drei Docker-Gates), die Aussage in Punkt 1, `make gates` sei der einzige Nachweis für alles, was nicht freigegeben ist (der Nachweis schließt die Gegenproben nicht mehr ein), und die Verwerfung der Alternative D. Unverändert bleiben die Positivliste, die Dateiliste des Nachweises, die Wahl am geänderten Pfad und der Rückfall auf die volle Kette. [ADR-0033](0033-gate-nachweise-in-der-abdeckung.md) liest sich ab jetzt so: „hängt an der Gate-Kette“ heißt „steht in einer der beiden Listen“, die Nachweisart Gate bleibt. Die ADRs [0025](0025-benannte-slice-kennungen-im-commit-hook.md), [0028](0028-abdeckung-je-anforderung-und-pfad.md), [0029](0029-benannte-welle-kennungen-im-commit-hook.md), [0032](0032-kopf-sensor-fuer-slice-plaene.md), [0034](0034-lint-gate-mit-solid-nahem-profil.md), [0042](0042-versionsmatrix-der-integrationstests.md), [0043](0043-binaries-und-container-gate.md), [0044](0044-schwelle-der-anweisungsabdeckung.md) und [0045](0045-container-start-als-werkzeug-im-repo.md) nennen eine Gegenprobe als Teil von `make gates`; sie bleiben unverändert und gelten, soweit diese ADR nichts anderes sagt. Die Einzelregeln stehen in [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) und [`SPEC-052`](../../../harness/vertraege/SPEC-052-gate-kette.md).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Gegenproben nur bei geänderten Eingaben (Alternative A der ADR 0038, erster Zuschnitt) | `make gates` bleibt ein Aufruf mit allen Gegenproben | Eine Eingabeliste je Gegenprobe ist unvollständig, sobald ein Quelltext ein Befund wird, und ein Fehler ist ein stilles Grün ([ADR-0038](0038-gate-kette-voll-und-schnell.md), Befund 1 und 2); vom Nutzer verworfen („Alle Gegenproben raus“) |
| B — unverändert, alle Gegenproben in der Kette | Nichts bleibt ungeprüft; der Nachweis sagt alles | Laufzeit und Last halten jeden Lauf auf demselben Host auf; der Nutzer sieht den Zweck nicht, wenn sich nichts geändert hat |
| **C — heraus aus der Kette, Pflicht, eigener Container** | `make gates` ohne die schwersten Läufe; die Gegenproben laufen im eigenen Container, ohne die Arbeit zu bremsen; der Beleg liegt vor, wo er gebraucht wird (Welle, Release) | Lockerung der Kette (`AGENTS.md` §3.6); nichts erzwingt den Lauf auf dem Hauptbaum; ein Bruch bleibt bis zur nächsten Fahrt unbemerkt |
| D — heraus, ohne Pflicht und ohne Container | wie C, ohne Aufwand | Genau die Harness-Lüge: Ein Gate, das niemand fährt, belegt nichts, und der Nachweis grün sagt trotzdem „alles geprüft“ |
| E — nur die langsamen Gegenproben heraus | kleinere Lockerung | Die Grenze ist willkürlich (die Last ist die Summe), und der Nutzer hat „alle“ gesagt |
| F — wie C, dazu ein Stand der Gegenproben und ein Hook, der ihn für Änderungen an Gates und Werkzeugen verlangt | Der Lauf wäre erzwungen | Bringt die Laufzeit zurück, die die Entscheidung streicht; erst zu erwägen, wenn die Pflicht nicht trägt |

## Konsequenzen

- Positiv: `make gates` fährt keine Gegenprobe; die Dauer wird im Slice gemessen und genannt. Die Gegenproben laufen im eigenen Container, ohne dass Implementer, Planner und Hook warten.
- Negativ, offen geschrieben: `make gates` grün sagt nichts mehr über die Gegenproben, und der Stop-Hook gibt auf diesen Nachweis frei. Ein Bruch eines Gates oder seiner Gegenprobe fällt erst auf, wenn jemand die Gegenprobe fährt. Das ist die Lockerung nach `AGENTS.md` §3.6, die diese ADR dokumentiert.
- Gegenmittel und hinterlegter Beleg: die Pflicht zum Fahren (Änderung, Ende der Welle, Release), das Ergebnis der Fahrt im Gegenproben-Container nach der Einlieferung und die Prüfung dieses Ergebnisses in der Closure von Welle und Release. Die Pflicht trägt nur, wenn jemand sie liest; ein Sensor für „geprüftes Werkzeug geändert“ fehlt (`BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`).
- Akzeptiertes Negativ: Nichts erzwingt den Lauf auf dem Hauptbaum. Grund: Ein Erzwingen wäre Option F, ein zweiter Nachweis, der die Laufzeit zurückbrächte, und der Nutzer will die Kette lockern, nicht ersetzen.
- Akzeptiertes Negativ: Die Abdeckung führt Gegenproben weiter als Nachweisart Gate ([ADR-0033](0033-gate-nachweise-in-der-abdeckung.md)); belegt ist sie durch die Fahrt bei Pflicht, nicht durch jeden Lauf von `make gates`. Grund: `make abdeckung-check` liest Deklarationen und prüft schon heute nicht, ob eine Gegenprobe an der Kette hängt.
- Akzeptiertes Negativ: Eine Container-Gegenprobe in der Liste für den Host läuft am Host unbemerkt weiter. Grund: Der Fehler in der anderen Richtung (Host-Gegenprobe in der Container-Liste) zeigt sich bei der ersten Fahrt im Container.
- Folgepflicht: [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) und die Nachzüge in [`SPEC-052`](../../../harness/vertraege/SPEC-052-gate-kette.md) bis `SPEC-058`; `harness/README.md` §Sensors führt die Listen und Sammelziele und die Bindung jeder Gegenprobe; `AGENTS.md` §6 Schritt 6 trägt die Pflicht (Wortlaut: Empfehlung des Planners, der Nutzer sagt beim Start des Vorhabens ausdrücklich Ja); die Rollen-Anweisungen und der Gegenproben-Container folgen mit dem Vorhaben danach.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, make, Attrappen | Namensregel und Eindeutigkeit der Listen; Ausnahmeliste der schnellen Kette auf drei Gates | `make gates-schnell-gegenprobe` |
| bash, make, Attrappen | Sammelziele: Reihenfolge, rotes Einzelziel, Abbruch, leere Liste, Aufrufarten von `make`, Reihenfolge der Gründe | `make gegenproben-gegenprobe` |

Die Pflicht selbst prüft kein Gate; sie ist Regel an den Aufrufer ([`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 14).

## Re-Evaluierungs-Trigger

Wenn eine Regression eines geprüften Werkzeugs unbemerkt blieb, die eine nicht gefahrene Gegenprobe gefangen hätte (dann Option F oder B); wenn die Closure einer Welle oder eines Releases kein Ergebnis des Gegenproben-Containers findet; wenn `make gates` ohne Gegenproben und `make gates-schnell` sich in der Dauer nicht mehr unterscheiden (dann entfällt die zweite Kette); wenn der Nutzer die Option A wieder will.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-18 | Proposed (Empfehlung des Architect nach der Richtung des Nutzers im Dialog) | — |
| 2026-10-18 | Accepted: Entscheidung des Nutzers, Option C mit Pflicht und eigenem Container; `make gates-schnell` bleibt vorerst | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
