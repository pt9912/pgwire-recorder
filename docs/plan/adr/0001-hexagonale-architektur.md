# ADR-0001: Hexagonale Architektur

**Status:** Accepted

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität)

**Schärft:** [`ARC-001`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-002`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-003`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-004`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-005`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-006`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-007`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-008`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-009`](../../../spec/architecture.md#1-komponenten-übersicht), [Schichten und Constraints](../../../spec/architecture.md#2-schichten-und-constraints)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Fachlogik (Record/Replay, Matching) und Infrastruktur (PGWire, PostgreSQL, Dateisystem) müssen getrennt testbar sein; Replay muss ohne Netzwerk und Dateisystem geprüft werden können. Annahme: Das Produkt bleibt ein Testwerkzeug, kein Produktions-Proxy.

## Entscheidung

Wir wählen: Ports & Adapters werden als verbindlicher Architekturstil verwendet.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Klassische Schichtenarchitektur | einfach, bekannt | Infrastruktur-Abhängigkeiten wandern leicht in die Fachlogik; Core nur mit Infrastruktur testbar |
| B — Ein Package ohne Strukturvorgabe | geringster Anfangsaufwand | Replay-Logik und PGWire-Code verschmelzen; kaum Tests ohne Netzwerk |
| **C — Hexagonal (gewählt)** | Core ohne Infrastruktur, Fakes für Ports, austauschbare Adapter | mehr Mapping-Code und mehr Packages |

## Konsequenzen

- Positiv: Infrastrukturabhängigkeiten liegen außerhalb des Application Core.
- Negativ: zusätzlicher Mapping-Aufwand an den Adaptergrenzen.
- Folgepflicht: Abhängigkeitsregeln maschinell prüfen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check mit `.a-check.yml` | Layer-Rollen und erlaubte Kanten gemäß Schichtentabelle | — |

## Re-Evaluierungs-Trigger

Wenn die Adaptergrenzen mehr Aufwand erzeugen, als sie an Testbarkeit einbringen, oder wenn ein zweiter Core nötig wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |
| 2026-10-03 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
