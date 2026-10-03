# ADR-0013: SQLite als zweites Aufzeichnungsformat

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), [ADR-0005](0005-recording-store-ist-driven-adapter.md)

**Schärft:** [`LH-FA-22.a`](../../../spec/spezifikation.md#lh-fa-22a--aufzeichnungsformat), [`SPEC-043`](../../../spec/spezifikation.md#spec-043--recording-sqlite-format)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Standardformat der Aufzeichnung ist YAML: lesbar, diff-freundlich, versionierbar. Das Schreiben ersetzt die Zieldatei atomar als Ganzes, damit sie nie unvollständig ist. Bei großen Aufzeichnungen, wie sie beim Einspielen einer Last anfallen, kostet das bei jeder beendeten Session das Neuschreiben der ganzen Datei, und der gesamte Bestand steht im Speicher. Der Port `RecordingRepository` ([ADR-0005](0005-recording-store-ist-driven-adapter.md)) erlaubt ein zweites Backend.

## Entscheidung

Wir wählen: Der Recorder speichert Aufzeichnungen wahlweise (`--format sqlite`) als SQLite-Datei; YAML bleibt der Standard. Beide Formate tragen dasselbe logische Modell, und Replay und Einspielen erkennen das Format der Datei selbst. Jede beendete Session wird in einer Transaktion ergänzt. Die Tabellenform steht im neutralen Schema-Format von d-migrate (`tools/schema/schema.yaml`); das SQL für SQLite wird daraus erzeugt und nicht von Hand geschrieben, wie bei den Schwester-Projekten.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nur YAML | ein Format, diff-freundlich | Neuschreiben der ganzen Datei je Session; ganzer Bestand im Speicher |
| D — SQLite mit handgeschriebenem DDL | keine zusätzliche Schema-Datei | zwei Quellen für dieselbe Schema-Form (DDL und Spezifikation), Abweichung unbemerkt |
| B — Zeilenorientierter Textstrom (ein Dokument je Session oder Interaktion, nur anhängen) | streambar, diff-freundlich, ohne Zusatzabhängigkeit | ein abgebrochener Schreibvorgang hinterlässt ein unvollständiges letztes Stück; ohne Transaktionsschutz und ohne wahlfreien Zugriff |
| **C — SQLite als zweites Format (gewählt)** | inkrementell und transaktional ohne Neuschreiben; wahlfreier Zugriff; ein Format, das sich unabhängig prüfen lässt | binär und nicht diff-freundlich; zusätzliche Bibliothek im Recording-Adapter; Pflege zweier Formate |

## Konsequenzen

- Positiv: Große Aufzeichnungen lassen sich ohne Neuschreiben und zu jedem Zeitpunkt konsistent speichern; YAML bleibt für alle, die ihre Aufzeichnungen versionieren wollen.
- Negativ: Zwei Formate müssen dasselbe Modell tragen; eine SQLite-Bibliothek kommt in den Recording-Adapter, bei plattformübergreifenden Binaries (Linux, macOS, Windows) bevorzugt ohne native Abhängigkeit.
- Folgepflicht: Ein Test, dass beide Formate im Replay und beim Einspielen dasselbe Verhalten liefern (Roundtrip-Gleichheit); die Importregel für die SQLite-Bibliothek im Architektur-Gate.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check mit `.a-check.yml` | die SQLite-Bibliothek nur im Recording-Adapter (`tech`-Regel, sobald die Bibliothek feststeht) | — |
| d-migrate | das neutrale Schema `tools/schema/schema.yaml` ist gültig | `make schema-validate` |

## Re-Evaluierungs-Trigger

Wenn Aufzeichnungen so groß werden, dass auch SQLite nicht genügt, oder wenn ein weiteres Format gefordert wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
