# ADR-0044: Schwelle der Anweisungsabdeckung der Unit-Tests in der Stufe test

**Status:** Proposed

**Datum:** 2026-10-17

**Autor:** Architect; Empfehlung, bis der Nutzer entschieden hat

**Bezug:** [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md), [ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0033](0033-gate-nachweise-in-der-abdeckung.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md)

**Schärft:** [`SPEC-056`](../../../harness/vertraege/SPEC-056-testabdeckung-schwelle.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`LH-QA-07` verlangt in Messmethode 2, dass die Anweisungsabdeckung der Unit-Tests nicht unter einer festgelegten Schwelle liegt und dass die Prüfung nachweislich rot wird. Heute misst kein Lauf sie. Die Stufe `test` ([ADR-0026](0026-build-und-test-im-multistage-dockerfile.md)) fährt `gofmt`, `go vet` und die Unit-Tests netzlos und schreibt nichts in den Arbeitsbaum. Die Abdeckung je Anforderung und Pfad ([ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md)) ist ein anderer Gegenstand.

**Messung.** Quellstand `69aaf3f7`, Stufe `source` des `Dockerfile`, `go test -covermode=set -coverprofile` über alle Pakete, netzlos: 4201 von 4448 Anweisungen, 94,45 % (`go tool cover -func`: 94,4 %). Je Paket: `adapters/driven/postgres` 97,7 %, `adapters/driven/recording` 94,2 %, `adapters/driving/cli` 98,0 %, `adapters/driving/pgwire` 96,4 %, `bootstrap` 96,8 %, `bootstrap/tlsproxy` 49,1 % (55 Anweisungen), `hexagon/model` 98,3 %, `hexagon/services` 98,7 %, `testuhr` und `testzeit` 100 %; ohne Test und mit 0 %: `cmd/pgwire-recorder` (12), `bootstrap/offsetlese` (13) und `testpki` (64). Mit `-coverpkg=./...` stiege die Zahl auf 95,8 %, weil Tests anderer Pakete mitzählten.

## Entscheidung

Wir wählen **Option A**: Die Stufe `test` misst die Unit-Tests mit Bordmitteln von Go (Profil, `-covermode=set`, ohne `-coverpkg`) und ein kleiner Prüfschritt vergleicht die Summe aller ausgeführten mit der Summe aller Anweisungen, über alle Pakete, mit einer **Schwelle für den Gesamtwert**. Die Schwelle ist der gemessene Wert, **abgerundet auf die ganze Zahl**: **94 %**. Die Zahl steht in der Spezifikation des Gates ([`SPEC-056`](../../../harness/vertraege/SPEC-056-testabdeckung-schwelle.md)), dort wird sie angehoben; jede Anhebung folgt derselben Regel (gemessener Wert zum Quellstand des Commits, abgerundet), und eine Senkung ist eine neue ADR (`AGENTS.md` §3.6). Alle Pakete mit Anweisungen zählen, auch `cmd/` und Pakete ohne Test; die Ausgabe nennt jedes Paket, ohne dass eines eine eigene Schwelle trägt. Die Einzelregeln stehen in der Spezifikation.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — Gesamtwert nach Anweisungen, Bordmittel, nur Unit-Tests, abgerundete ganze Zahl** | trägt auf dem heutigen Stand ohne neue Tests; kein neues Modul, kein Image; die Zahl ist ein Anteil, den jeder nachrechnet | ein schwaches Paket hinter starken ist nur in der Ausgabe sichtbar |
| B — Schwelle je Paket | strenger, kein Paket verdeckt | `bootstrap/tlsproxy` (49,1 %), `cmd/` und zwei weitere Pakete hielten keine gemeinsame Zahl; verlangte neue Tests vor dem Gate |
| C — Integrationstests über `GOCOVERDIR` mitzählen | misst auch den Prozess des Binaries | braucht ein `-cover`-Binary und einen Weg aus dem Container der Integrationstests, den [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md) nicht hat; Messung über zwei Läufe hinweg; `LH-QA-07` verlangt nur die Unit-Tests |
| D — `-coverpkg=./...` | höhere Zahl (95,8 %) | Tests anderer Pakete decken ein Paket, das keinen eigenen Test hat; die Zahl sagt weniger |
| E — Fremdwerkzeug (Dienst oder Modul für Berichte und Differenzen) | Berichte, Verlauf | Netz oder Pin in `deps`, ein weiterer Anbieter; kein Bedarf, den `go tool cover` nicht deckt |
| F — Spielraum größer als die Rundung (etwa 90 %) | weniger Anlass für Anhebungen | zehn Prozentpunkte (rund 440 Anweisungen) ungeprüft; das Gate hätte keinen Biss |

## Konsequenzen

- Positiv: `LH-QA-07` Messmethode 2 hat ein Gate in `make test`, ohne neue Tests, ohne Netz und ohne Zusatz im Arbeitsbaum.
- Negativ: Der Spielraum ist klein (19 Anweisungen); ein Umbau, der Code ohne Test einführt, braucht Tests im selben Commit.
- Akzeptiertes Negativ: `bootstrap/tlsproxy`, `cmd/` und Hilfspakete zählen mit ihrem 0 % in den Gesamtwert und heben die Zahl nicht; der Grund: eine Ausnahmeliste wäre ein zweiter Ort für Entscheidungen, und die Schwelle folgt der Messung. Die Zahl sagt, dass Anweisungen laufen, nicht, dass etwas geprüft wird; das bleibt `AGENTS.md` §3.10.
- Folgepflicht: `harness/README.md` §Sensors führt die Schwelle beim Vertrag von `make test` und die Gegenprobe mit Bindung an diese ADR; die Ausnahmeliste der schnellen Kette ([ADR-0038](0038-gate-kette-voll-und-schnell.md)) nennt `coverage-gegenprobe`; die Abdeckungstabellen führen `LH-QA-07` Teil 2.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Profil in der Stufe `test` | Gesamtwert nicht unter der Schwelle aus `SPEC-056` | `make test` |
| Gegenprobe | Stand unter der Schwelle, Randformen des Profils und der Schwelle, Paket ohne Test, Stufe bleibt bei Format- und Testfehler rot | `make coverage-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn ein Paket mit eigener Logik dauerhaft hinter dem Gesamtwert liegt und das Review es nicht fängt (dann Option B in eigenem Slice), wenn das Lastenheft die Integrationstests in die Messung einbezieht, wenn `-race` in die Stufe `test` kommt (dann `atomic` und eine neue Spezifikationsstelle) oder wenn Code generiert wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-17 | Proposed: Empfehlung des Architect | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
