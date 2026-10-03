# ADR-0003: PGWire Server ist Driving Adapter

**Status:** Accepted

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-04`](../../../spec/lastenheft.md#lh-fa-04--transparenz-für-clients), [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-03`](../../../spec/lastenheft.md#lh-fa-03--replay-modus)

**Schärft:** [`ARC-006`](../../../spec/architecture.md#1-komponenten-übersicht)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der PostgreSQL-Client stößt im Record- und im Replay-Modus die Verarbeitung an; der Recorder ist aus Sicht des Clients ein PostgreSQL-Endpunkt.

## Entscheidung

Wir wählen: Der PostgreSQL-Client ist ein Driving Actor. Der Serverteil des Recorders ist deshalb ein Driving Adapter.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — PGWire-Server als Teil des Cores | weniger Schichten | Netzwerkcode im Core; Replay nur mit Sockets testbar |
| B — PGWire-Server als Driven Adapter, vom Core gesteuert | Core bestimmt den Ablauf | Gegenläufig zum Kontrollfluss: der Client stößt an, nicht der Core |
| **C — Driving Adapter (gewählt)** | entspricht dem tatsächlichen Kontrollfluss; Core ohne Netzwerk | Übersetzung in Domain-Typen nötig |

## Konsequenzen

- Positiv: Netzwerk-/PGWire-Code ruft Inbound Ports auf und enthält keine Replay-Fachlogik.
- Negativ: Startup und `SSLRequest` liegen technisch im Adapter und müssen dort getestet werden.
- Folgepflicht: Adapter-Contract-Tests für Startup und `SSLRequest`.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check mit `.a-check.yml` | `adapters-driving` darf nur `ports-driving` und `model` importieren | — |

## Re-Evaluierungs-Trigger

Wenn ein weiterer Eingang (z. B. eine API) dieselben Use Cases anspricht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |
| 2026-10-03 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
