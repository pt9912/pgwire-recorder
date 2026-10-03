# ADR-0010: Verwendung von pgproto3

**Status:** Accepted

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol)

**Schärft:** [`SPEC-030`](../../../spec/spezifikation.md#6-externe-verträge), [`ARC-012`](../../../spec/architecture.md#3-externe-abhängigkeiten)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Verarbeitung von PGWire-Nachrichten soll nicht selbst geschrieben werden, soweit eine gepflegte Bibliothek sie abdeckt.

## Entscheidung

Wir wählen: Für PGWire-Nachrichten wird `github.com/jackc/pgx/v5/pgproto3` verwendet, ausschließlich als Infrastrukturabhängigkeit in den beiden PGWire-Adaptern.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Eigener PGWire-Codec | volle Kontrolle | hoher Aufwand und Fehlerrisiko |
| B — Umfassendes PostgreSQL-Server-Framework | viel Fertiges | bringt Verhalten mit, das der Recorder nicht will; Kontrolle über Nachrichtenfluss geringer |
| **C — pgproto3 (gewählt)** | nachrichtengenau, schlank | deckt nicht jede Nachricht ab |

## Konsequenzen

- Positiv: bewährte Kodierung der Nachrichten.
- Negativ: Nachrichten, die die Bibliothek nicht verlustfrei repräsentiert, sind nicht unterstützt.
- Folgepflicht: Import-Regel im Architektur-Gate.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check mit `.a-check.yml` | `pgproto3` nur in `adapters-driving/pgwire` und `adapters-driven/postgres` (`tech`-Regel) | — |

## Re-Evaluierungs-Trigger

Wenn die Bibliothek benötigte Nachrichtentypen nicht abdeckt oder nicht mehr gepflegt wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |
| 2026-10-03 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
