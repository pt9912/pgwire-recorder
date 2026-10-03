# ADR-0007: Strict Replay

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-10`](../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus)

**Schärft:** [`LH-FA-09.a`](../../../spec/spezifikation.md#lh-fa-09a--strict-sequential-matching), [`SPEC-011`](../../../spec/spezifikation.md#3-defaults-und-konstanten)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Replay muss deterministisch sein und darf bei Abweichungen nicht stillschweigend eine unpassende Antwort liefern. Annahme: Ein Recording gehört zu einem festen Testablauf.

## Entscheidung

Wir wählen: v1 verwendet exaktes sequenzielles Matching der Anfragen.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Lookup nach SQL-Text unabhängig von der Reihenfolge | robust gegen veränderte Reihenfolge | Mehrdeutigkeit bei doppelten Queries; Zustand der Session wird ignoriert |
| B — Normalisiertes oder unscharfes Matching | tolerant | nicht deterministisch nachvollziehbar; verdeckt Abweichungen |
| **C — Strict sequential (gewählt)** | deterministisch, Abweichungen sofort sichtbar | empfindlich gegen jede Textänderung |

## Konsequenzen

- Positiv: deterministische Tests, keine implizite SQL-Normalisierung.
- Negativ: Tests brechen bei Whitespace- oder Reihenfolge-Änderungen.
- Folgepflicht: Mismatch-Diagnose gemäß Spezifikation.

## Fitness Function (falls maschinell prüfbar)

Nicht maschinell prüfbar.

## Re-Evaluierungs-Trigger

Wenn Anwender regelmäßig an der Strenge scheitern oder ein toleranter Modus als eigene Anforderung aufgenommen wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
