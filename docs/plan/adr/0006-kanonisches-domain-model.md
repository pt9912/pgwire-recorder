# ADR-0006: Kanonisches Domain Model

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-06`](../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus)

**Schärft:** [`ARC-001`](../../../spec/architecture.md#1-komponenten-übersicht)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

PGWire-Nachrichten kommen als Bibliothekstypen; der Core soll davon unabhängig und über Versionen der Bibliothek hinweg stabil bleiben.

## Entscheidung

Wir wählen: PGWire-Bibliothekstypen werden an den Adaptergrenzen in eigene Domain-Typen übersetzt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Bibliothekstypen im Core verwenden | kein Mapping | Core an die Bibliothek gekoppelt; Versionswechsel schlägt durch |
| B — Rohe Bytes im Core | kein Mapping von Nachrichtentypen | Core muss das Protokoll selbst dekodieren |
| **C — Eigene Domain-Typen (gewählt)** | stabile innere Architektur | mehr Mapping-Code |

## Konsequenzen

- Positiv: stabiles Modell, Core ohne Drittbibliothek.
- Negativ: Mapping-Code in beiden PGWire-Adaptern.
- Folgepflicht: Tests, dass Bytes nach Load/Save identisch bleiben.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check mit `.a-check.yml` | `model` importiert keine Adapter und keine Infrastrukturbibliothek | — |

## Re-Evaluierungs-Trigger

Wenn das Mapping zum Engpass wird oder die Bibliothek gewechselt wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
