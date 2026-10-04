# Review-Report: slice-walking-skeleton-build-gates, Folge-Review — 2026-10-04

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Das ist ein Folge-Lauf zu [`2026-10-04-review-slice-walking-skeleton-build-gates.md`](2026-10-04-review-slice-walking-skeleton-build-gates.md) (F-204 bis F-213). Die DoD ist nicht Gegenstand dieses Reviews. Sie prüft der Verifier, und sein Beleg liegt vor: [`2026-10-04-verifikation-slice-walking-skeleton-build-gates.md`](2026-10-04-verifikation-slice-walking-skeleton-build-gates.md).

**Gegenstand:** `git diff 8143a47 a834db1`. Darin geändert: `.a-check.yml` (`tech`), `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk`, `cmd/pgwire-recorder/main.go`, die neun `doc.go` unter `internal/`, `harness/README.md` §Sensors, `README.md` und der Slice-Plan (§1 Abgrenzung, §3, §6). Der Report des Vorlaufs liegt ebenfalls im Diff, wird hier aber nicht bewertet.

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:**

- Slice-Plan `slice-walking-skeleton-build-gates` (§1 bis §6)
- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0009](../plan/adr/0009-implementierungssprache-go.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) (Entscheidung und Fitness Function: `pgproto3` nur in den beiden PGWire-Adaptern)
- `LH-QA-03`, `LH-QA-04`
- `AGENTS.md` §3 (Hard Rules), `harness/README.md`, `harness/conventions.md`
- der Vorlauf-Report (F-204 bis F-213) und der Verifikationsbeleg (nur als Abgrenzung)

**Ausgeführte Läufe:** Im Repo liefen `make a-check` (0 Befunde), `make a-check-negativ` (grün), `make build` (grün) und `make test` (grün, keine Testdateien). Die Mutationen liefen nur an Kopien der versionierten Dateien im Scratchpad, mit dem Skript der Gegenprobe bzw. a-check im gepinnten Image:

| # | Mutation an `.a-check.yml` bzw. Umgebung | Gegenprobe | a-check direkt |
|---|---|---|---|
| M1 | `tech` in zwei Einträge je Bibliothek zerlegt, Server-Adapter zuerst (Stand vor `a834db1`) | rot (Fall `postgres`) | — |
| M2 | zwei Einträge je Bibliothek, Upstream-Adapter zuerst | **grün** | `pgproto3`/`crypto/tls` im Server-Adapter: `tech-leak … außerhalb internal/adapters/driven/postgres`, rot |
| M3 | Liste nur mit `internal/adapters/driven/postgres` | **grün** | `pgproto3` im Server-Adapter: tech-leak, rot |
| M4 | `pgproto3`-Zeile entfernt | rot (Fälle `model`, `recording`) | — |
| M5 | `crypto/tls` nur im Server-Adapter | rot (Fall `postgres`) | — |
| M6 | `DOCKER=false` (Runtime scheitert) | rot (alle drei Fälle, mit Meldung) | — |
| M7 | `pgproto3`-Liste um den Recording-Adapter erweitert | rot (Fall `recording`) | — |
| — | aktueller Stand, ohne Mutation | — | `pgproto3`+`crypto/tls` im Server-Adapter grün; `pgproto3` im Model: core-impurity, in services: app-impurity, im CLI-Adapter: tech-leak; `crypto/tls` im Recording-Adapter: tech-leak |

Bei Lauf-Ende zeigte `git status` keine Änderung durch diesen Review. Zwischendurch stand eine fremde Änderung an `docs/plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md` im Arbeitsbaum, inzwischen als `3810a34`/`1b4fe66` committet. Sie gehört nicht zu diesem Review. `make gates` lief nicht.

---

## Stand der Findings aus dem Vorlauf

| ID | Vorlauf | Stand | Beleg |
|---|---|---|---|
| F-204 | HIGH | **behoben** | `tech` führt je Bibliothek einen Eintrag mit beiden PGWire-Adaptern als Liste. a-check wertet die Liste aus (Meldung `außerhalb internal/adapters/driving/pgwire\|internal/adapters/driven/postgres`). `pgproto3` und `crypto/tls` sind in beiden Adaptern zugelassen und anderswo abgelehnt (Tabelle oben). Die Plan-Zeile `.a-check.yml` steht in §3. Wie weit die Gegenprobe diesen Fix festhält, steht in F-214. |
| F-205 | MEDIUM | **bekannt, außerhalb des Slice** | Wird über eine eigene ADR entschieden ([ADR-0025](../plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md), Proposed). In diesem Review nicht bewertet. |
| F-206 | LOW | **behoben** | Die Sensors-Zeile von `make a-check-negativ` bindet [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) und [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md), und die Message von `a834db1` nennt [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md). |
| F-207 | LOW | **behoben** | Fall `model` erwartet den Pfad der Probedatei, Fall `recording` die Regelklasse `tech-leak`. Scheitert die Runtime oder die Konfiguration, wird die Probe rot (M6). |
| F-208 | LOW | **behoben** | Kopien und Logs liegen unter einem gemeinsamen `mktemp -d`, das `trap … EXIT` vollständig räumt. Explizite `rm` gibt es nicht mehr. |
| F-209 | LOW | **behoben** | Die Zustandssätze in den neun `doc.go` und in `main.go` sind entfernt. Was bleibt, sind Zweck-Zusagen je Package. |
| F-210 | INFO | **behoben für `arch-negativ.mk`**, Muster wiederholt sich an anderer Stelle | Der Kopfkommentar nennt genau die geprüften Fälle und grenzt ab („andere Regeln … prüft sie nicht“). Skript-Kopf, Hilfetext und `README.md` beschreiben dagegen mehr, als geprüft wird (F-214). |
| F-211 | INFO | **notiert** | Steht als Risiko in Plan §6. `slice-walking-skeleton-record` nennt die netzlose Quelle der Module weiterhin nicht. Die Übergabe an den Planner hat damit ein Artefakt nur in diesem Slice. Ob bei Closure eine Paarung mit dem Folge-Slice nötig ist, entscheidet der Verifier. |
| F-212 | INFO | **notiert** | Steht als Risiko in Plan §6. Eine Entscheidung des Architects gibt es nicht. Das Gate deckt `os` im Model weiterhin nicht ab, und der Reviewer bleibt dort der Wächter. |
| F-213 | INFO | **behoben** | Plan §1 grenzt `test/integration/` mit Begründung ab. `slice-walking-skeleton-record` führt `test/integration` als „neu“. |

## Neue Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-214 | LOW | Die Gegenprobe prüft die Zulassung nur im Upstream-Adapter. Laut Skript-Kopf hält Fall 3 fest, „dass beide PGWire-Adapter die Bibliothek nutzen dürfen“, und laut Hilfetext bzw. `README.md` prüft sie „pgproto3 nur in den PGWire-Adaptern“. Zerfällt die Regel wieder in zwei Einträge mit dem Upstream-Adapter zuerst (M2) oder fehlt der Server-Adapter in der Liste (M3), bleibt die Gegenprobe grün, obwohl a-check `pgproto3` im Server-Adapter als tech-leak ablehnt. Rot wird sie beim Zerfall nur in der Reihenfolge des Stands vor `a834db1` (M1). Die Sensors-Zeile („im Upstream-Adapter zugelassen“) und der Kopf von `arch-negativ.mk` beschreiben die Reichweite genau. Das ist das zweite Auftreten der Klasse aus F-210. | [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md); `AGENTS.md` §3.7; Maintainability | `tools/arch/a-check-negativ.sh` · Kopfkommentar („Fall 3 haelt fest …“) und Fall-Aufrufe; `harness/mk/arch-negativ.mk` · Hilfetext des Ziels; `README.md` · Zeile „Gates“ | ja: Gegenprobe an einer Kopie mit Mutation M2 oder M3 | Reichweite einer Gegenprobe allgemeiner beschrieben als geprüft |
| F-215 | LOW | Plan §3 beschreibt die Gegenprobe weiterhin als „rot, wenn a-check den `pgproto3`-Import im Domain Model durchlässt“. Der Diff führt drei Fälle: Model und Recording-Adapter abgelehnt, Upstream-Adapter zugelassen. Die neue Zeile `.a-check.yml` daneben gibt ihren Diff dagegen wieder. | Maintainability (Plan gibt Diff wieder) | Slice-Plan §3 · Zeile `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk` | ja (Lesen) | Plan-Zeile nach Korrektur nicht nachgezogen |
| F-216 | INFO | `crypto/tls` kommt in der Gegenprobe nur im Zulassungsfall vor (Fall 3). Dass es außerhalb der PGWire-Adapter abgelehnt wird, belegt nur ein direkter a-check-Lauf in diesem Review (Recording-Adapter: tech-leak), kein Gate. Weder die Sensors-Zeile noch der `mk`-Kopf behaupten das. Die DoD verlangt es nicht; die Bewertung liegt beim Verifier. | [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) Fitness Function; Verweis Verifier | `tools/arch/a-check-negativ.sh` · Fall-Aufrufe | ja: Mutation, die `crypto/tls` für den Recording-Adapter freigibt | Teil einer Tech-Regel ohne Gegenprobe |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.a-check.yml` `tech` gegen [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) | geprüft, ohne Befund. Beide Bibliotheken sind auf genau die beiden PGWire-Adapter beschränkt. `yaml`/`sqlite` sind unverändert. Die Kommentarzeilen der Datei sind unberührt. |
| Hard Rule 3.6: Gates nicht lockern | geprüft, ohne Befund. Die Liste gibt `internal/adapters/driven/postgres` frei, das vorher faktisch abgelehnt wurde. Das setzt die Entscheidung von [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) um und senkt keine Schwelle. Kein anderer Ort ist hinzugekommen (M7 und CLI-Probe zeigen weiter rot). |
| `tools/arch/a-check-negativ.sh`: `set -euo pipefail` mit `\|\| ergebnis=rot` | geprüft, ohne Befund. Der Docker-Lauf steht links von `\|\|` und bricht das Skript nicht ab. `tar`- und `mkdir`-Fehler brechen ab, und `trap` räumt auf. `grep` steht in einer `elif`-Bedingung und löst kein `-e` aus. `fehler` sammelt über alle drei Fälle, sodass jeder rote Fall berichtet wird. |
| `tools/arch/a-check-negativ.sh`: Quoting | geprüft, ohne Befund. Pfade und Image sind gequotet. `for i in $imports` trennt absichtlich an Leerzeichen, die Werte sind Literale ohne Glob-Zeichen. Das Muster in `grep -q --` ist ein regulärer Ausdruck, und der Punkt im Dateinamen ist für das Ergebnis unerheblich. |
| `tools/arch/a-check-negativ.sh`: Aufräumen, Seiteneffekte | geprüft, ohne Befund. Alle Kopien und Logs liegen unter einem Temp-Verzeichnis. Der Mount ist `:ro` mit `--network none`, und der Arbeitsbaum bleibt nach `make a-check-negativ` unverändert. Das Skript kopiert dreimal den ganzen Arbeitsbaum (inklusive unversionierter Dateien, ohne `.git`). Das kostet Zeit, ist aber kein Defekt. |
| Gegenprobe bei Zerfall der Regel | geprüft. In der Reihenfolge vor `a834db1` wird sie rot (M1), in umgekehrter Reihenfolge bleibt sie grün (M2), siehe F-214. Entfällt die Regel (M4), werden zwei Fälle rot. Bei teilweisem Entfall für `crypto/tls` (M5) und bei Erweiterung auf den Recording-Adapter (M7) wird die Probe rot. Fällt die Runtime aus (M6), wird sie rot. |
| `harness/mk/arch-negativ.mk` | geprüft. Der Kopf ist genau und hält Hard Rule 3.7 ein (Zusage plus Abgrenzung), das Ziel hängt an `GATE_CHECKS`. Zum Hilfetext siehe F-214. |
| Hard Rule 3.7: Kommentare in `doc.go`, `main.go`, Skript, `mk` | geprüft. Sie stehen im Indikativ und tragen Zusage, Kopplung oder Abgrenzung. Es gibt keine verworfene Alternative, keinen abwesenden Text und keinen abgebrochenen Satz. Nach dem Entfernen der Zustandssätze bleibt kein leerer `//`-Rest. Zur Zusage im Skript-Kopf siehe F-214. |
| Hard Rule 3.1: Docker-only | geprüft, ohne Befund. Es gibt keinen neuen Host-Bedarf, und a-check, Build und Test laufen im gepinnten Image. |
| Hard Rule 3.3: Move und Inhalt getrennt | geprüft, ohne Befund. `a834db1` enthält keinen Move. |
| Hard Rule 3.5: ADRs | geprüft, ohne Befund. Im Diff steht keine ADR. |
| `harness/README.md` §Sensors gegen `make help` und Skript | geprüft, ohne Befund. Die Zeile `make a-check-negativ` beschreibt genau die drei Fälle, die Bindung siehe F-206. |
| Slice-Plan §1, §6 | geprüft, ohne Befund. Die Abgrenzung `test/integration/` ist begründet. Die zwei neuen Risiken tragen einen Ausgang in der Form der übrigen. Plan §3 siehe F-215. |
| Commit-Message `a834db1` | geprüft, ohne Befund. Sie nennt `LH-QA-04`, [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) und [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md). Struktur-IDs nennt sie keine. |
| `spec/` | geprüft, ohne Befund. Nicht im Diff. |

## Summary

| Kategorie | Anzahl (neu) |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

Aus dem Vorlauf sind F-204 (HIGH), F-206 bis F-209 und F-213 behoben. F-210 ist im `mk`-Kopf behoben, sein Muster kehrt aber als F-214 wieder. F-211 und F-212 stehen als Risiken im Plan. F-205 ist bekannt und liegt außerhalb des Slice ([ADR-0025](../plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md)).

**Finding-Klassen dieses Laufs:** Reichweite einer Gegenprobe allgemeiner beschrieben als geprüft (2. Auftreten nach F-210; beim dritten greift §Pflege) · Plan-Zeile nach Korrektur nicht nachgezogen · Teil einer Tech-Regel ohne Gegenprobe

## Verdikt

**Merge-blockierend:** nein. Der HIGH-Befund F-204 ist behoben, und die Behebung ist durch eigene Mutationen bestätigt. Keine neue Stufe liegt über LOW. F-214 und F-215 gehen an den Implementer, der sie annimmt oder begründet. Eine Konflikt-Sequenz ist dafür nicht nötig.

**Übergabe:** F-214 und F-215 gehen an den Implementer. F-216 und die Paarungsfrage zu F-211 gehen an den Verifier. F-212 bleibt beim Architect offen. Die Finding-Klassen gehen in die Slice-Closure §7, und die Wiederholung von F-210/F-214 wird dort als Steering-Loop-Zähler fortgeschrieben.
