# ADR-0002: Driving/Driven für Adapter, Inbound/Outbound für Ports

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** —

**Schärft:** [`ARC-003`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-004`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-005`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-006`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-007`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-008`](../../../spec/architecture.md#1-komponenten-übersicht), [Schichten und Constraints](../../../spec/architecture.md#2-schichten-und-constraints)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die hexagonale Architektur (ADR-0001) braucht eindeutige Begriffe, damit Package-Namen, a-check-Rollen und Dokumentation dieselbe Sprache sprechen.

## Entscheidung

Wir wählen: Adapter heißen **Driving** (rufen den Core) oder **Driven** (vom Core aufgerufen); Ports heißen **Inbound** (Use Cases des Cores) oder **Outbound** (Anforderungen des Cores). Die Package-Namen `ports/driving` und `ports/driven` entsprechen Inbound und Outbound.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Einheitlich Inbound/Outbound für Ports und Adapter | ein Begriffspaar | Adapter- und Port-Rolle verschwimmen |
| B — Primary/Secondary | verbreitet | Rollen nicht aus dem Namen ableitbar |
| **C — Driving/Driven + Inbound/Outbound (gewählt)** | Rolle jeder Komponente aus dem Namen ablesbar | zwei Begriffspaare nötig |

## Konsequenzen

- Positiv: Package-Namen spiegeln die Rolle eindeutig wider.
- Negativ: Neue Mitwirkende müssen zwei Begriffspaare lernen.
- Folgepflicht: a-check-Rollen und -Richtungen entsprechend benennen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check mit `.a-check.yml` | `role` und `direction` je Layer | — |

## Re-Evaluierungs-Trigger

Permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
