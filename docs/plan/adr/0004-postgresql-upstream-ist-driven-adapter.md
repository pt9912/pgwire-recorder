# ADR-0004: PostgreSQL Upstream ist Driven Adapter

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus)

**Schärft:** [`ARC-004`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-007`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-010`](../../../spec/architecture.md#3-externe-abhängigkeiten)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Record-Modus braucht eine Verbindung zu einem realen PostgreSQL pro Client-Session; der Core darf davon nichts über TCP oder Bibliothekstypen wissen.

## Entscheidung

Wir wählen: Die reale PostgreSQL-Datenbank wird durch den Core über einen Outbound Port angesprochen; die Implementierung ist ein Driven Adapter.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Core nutzt die Verbindung direkt | kürzester Weg | Core kennt TCP und Bibliothekstypen; keine Fakes möglich |
| B — Bytes unverändert durchreichen (Byte-Proxy) | sehr einfacher Record-Pfad | Core sieht keine fachlichen Interaktionen; Aufzeichnung und Mapping entfallen nicht, sondern wandern in den Proxy |
| **C — Outbound Port + Driven Adapter (gewählt)** | Record-Service ist mit einem Fake testbar | Mapping-Aufwand |

## Konsequenzen

- Positiv: Der Record-Service kennt keine Bibliotheks- oder TCP-Typen.
- Negativ: Streaming und Backpressure müssen im Port-Design mitgedacht werden.
- Folgepflicht: Adapter-Contract-Test gegen eine reale PostgreSQL-Testinstanz.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check mit `.a-check.yml` | `adapters-driven` darf nur `ports-driven` und `model` importieren | — |

## Re-Evaluierungs-Trigger

Wenn ein weiterer Upstream-Typ (anderes Protokoll) unterstützt werden soll.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
