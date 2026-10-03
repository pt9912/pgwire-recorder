# ADR-0005: Recording Store ist Driven Adapter

**Status:** Accepted

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-07`](../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-QA-06`](../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats)

**Schärft:** [`ARC-004`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-008`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-011`](../../../spec/architecture.md#3-externe-abhängigkeiten), [`ARC-013`](../../../spec/architecture.md#3-externe-abhängigkeiten)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Recordings werden persistiert und später wiedergegeben; das Format ist versioniert, YAML ist das v1-Format.

## Entscheidung

Wir wählen: Persistenz wird über einen `RecordingRepository`-Port abstrahiert; der YAML-Adapter implementiert ihn.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — YAML-Serialisierung direkt in den Services | wenig Code | YAML-Typen im Core; Formatwechsel berührt die Fachlogik |
| B — Festes Datenbankformat (z. B. SQLite) | Abfragen möglich | nicht diff-freundlich; schwer versionierbar |
| **C — Repository-Port + YAML-Adapter (gewählt)** | Format austauschbar; Core ohne YAML | zusätzliche Abstraktion |

## Konsequenzen

- Positiv: YAML ist austauschbar und nicht Teil der Domain.
- Negativ: Persistenz-DTOs im Adapter nötig, sobald Domain- und YAML-Modell auseinanderlaufen.
- Folgepflicht: Roundtrip-Tests des Adapters.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check mit `.a-check.yml` | YAML-Bibliothek nur im Recording-Adapter (`tech`-Regel) | — |

## Re-Evaluierungs-Trigger

Wenn ein zweites Recording-Backend gebraucht wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |
| 2026-10-03 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
