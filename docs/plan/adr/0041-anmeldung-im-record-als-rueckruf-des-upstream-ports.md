# ADR-0041: Anmeldung im Record-Modus als Rückruf des Upstream-Ports

**Status:** Accepted

**Datum:** 2026-10-14

**Autor:** pt9912

**Bezug:** [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus), [ADR-0003](0003-pgwire-server-ist-driving-adapter.md), [ADR-0004](0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](0010-verwendung-von-pgproto3.md), [ADR-0030](0030-full-duplex-im-record-pfad.md)

**Schärft:** [`LH-FA-05.b`](../../../spec/spezifikation.md#lh-fa-05b--authentifizierung)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Im Record-Modus verlangt der Server beim Verbindungsaufbau gegebenenfalls ein Passwort, das nur der Client kennt ([`LH-FA-05.b`](../../../spec/spezifikation.md#lh-fa-05b--authentifizierung)). Der Aufbau zum Server läuft im Upstream-Adapter, die Verbindung zum Client im PGWire-Adapter, und der Core kennt nach [ADR-0003](0003-pgwire-server-ist-driving-adapter.md) und [ADR-0004](0004-postgresql-upstream-ist-driven-adapter.md) keine PGWire-Typen. Der Austausch ist ein Hin und Her mit unbestimmter Länge (SCRAM braucht drei Schritte), das mitten im Aufbau liegt.

## Entscheidung

Wir wählen: **Der Upstream-Port nimmt beim Öffnen einen Rückruf entgegen, über den der Upstream-Adapter jede Anmelde-Anforderung des Servers an den Client weitergibt und dessen Antwort zurückerhält.** Der Record-Service stellt den Rückruf bereit und reicht die Nachrichten als fachliche, für ihn undurchsichtige Werte an den PGWire-Adapter weiter, der sie mit dem Client spricht. Der Core liest den Inhalt nicht. Die Regeln im Einzelnen (Verfahren, Reihenfolge, Fehler, Abbruch) stehen in der Spezifikation.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Der Core kennt die Anmelde-Nachrichten als PGWire-Typen | einfacher Weg im Code | bricht die Richtung der Abhängigkeiten ([ADR-0001](0001-hexagonale-architektur.md), [ADR-0004](0004-postgresql-upstream-ist-driven-adapter.md)) |
| B — Der Upstream-Adapter spricht selbst mit der Verbindung des Clients | kein Rückruf | der Driven Adapter hängt am Driving Adapter; Architektur-Gate rot |
| C — `Open` liefert Zwischenzustände und wird mehrfach aufgerufen | keine Rückrufe | Zustand des Aufbaus wandert in den Service und den Port; jeder Verfahrensschritt wird ein Aufruf, die Gleichzeitigkeit je Session ([ADR-0030](0030-full-duplex-im-record-pfad.md)) wird schwerer zu sichern |
| **D — Rückruf in `Open` mit undurchsichtigen Nachrichten (gewählt)** | der Core bleibt frei von PGWire; der Ablauf des Servers bestimmt die Länge des Austauschs; ein Verfahren mehr braucht keine Änderung am Core | der Port bekommt einen Vertrag für Aufrufe nacheinander und für den Abbruch beim Schließen der Client-Verbindung |

## Konsequenzen

- Positiv: Klartext, MD5 und SCRAM-SHA-256 laufen ohne Kenntnis der Verfahren im Core; ein neues Verfahren des Servers wird nur im Adapter benannt.
- Negativ: Der Austausch läuft im Aufbau-Aufruf des Ports; das erste Signal darf ihn nicht abbrechen, das Zwangsende muss ihn über die Client-Verbindung beenden.
- Folgepflicht: Spezifikation und Port-Kommentar nennen den Vertrag der Aufrufe (nacheinander, nie gleichzeitig); Tests mit einem Fake-Rückruf je Zusage.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check mit `.a-check.yml` | der Core importiert keine PGWire-Bibliothek, der Driven Adapter keinen Driving Adapter | `make a-check` |
| Integrationstest gegen PostgreSQL | Anmeldung mit Klartext, MD5 und SCRAM-SHA-256 über den Recorder | `make test-integration` |

## Re-Evaluierungs-Trigger

Wenn ein Verfahren Zustand über den Aufbau hinaus braucht (etwa eine erneute Anmeldung während der Session) oder der Replay-Modus denselben Austausch führen soll.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-14 | Proposed | — |
| 2026-10-14 | Accepted (Entscheidung des Nutzers vom 2026-10-14) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
