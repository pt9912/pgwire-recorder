# ADR-0009: Implementierungssprache Go

**Status:** Accepted

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung)

**Schärft:** [`SPEC-036`](../../../spec/spezifikation.md#spec-036--technische-leitentscheidungen)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Werkzeug ist ein eigenständiges Kommandozeilenwerkzeug, das lokal, in CI und in Containern laufen soll, ohne Laufzeitabhängigkeit auf eine PostgreSQL-Installation.

## Entscheidung

Wir wählen: v1 wird in Go implementiert.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Skriptsprache (z. B. Python) | schnell entwickelt | Laufzeit und Abhängigkeiten müssen mitgeliefert werden; Binary-Bereitstellung aufwendiger |
| B — Rust | statisch gelinkte Binaries, hohe Kontrolle | höherer Entwicklungsaufwand; weniger passende Protokollbibliothek |
| **C — Go (gewählt)** | statisches Binary, Container-tauglich, passende PGWire-Bibliothek | Mapping zwischen Typen von Hand |

## Konsequenzen

- Positiv: ein eigenständiges Binary ohne Laufzeitabhängigkeit.
- Negativ: Sprachwahl bindet Tooling und Bibliotheken.
- Folgepflicht: `.a-check.yml` ist Go-spezifisch.

## Fitness Function (falls maschinell prüfbar)

Nicht maschinell prüfbar.

## Re-Evaluierungs-Trigger

Wenn die Protokollbibliothek wegfällt oder die Plattformanforderungen sich ändern.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |
| 2026-10-03 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
