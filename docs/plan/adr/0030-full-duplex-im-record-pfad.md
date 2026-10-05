# ADR-0030: Full-Duplex im Record-Pfad mit Zustand im Record-Service

**Status:** Accepted

**Datum:** 2026-10-05

**Autor:** pt9912

**Bezug:** [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-04`](../../../spec/lastenheft.md#lh-fa-04--transparenz-für-clients), [ADR-0001](0001-hexagonale-architektur.md), [ADR-0003](0003-pgwire-server-ist-driving-adapter.md), [ADR-0004](0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0012](0012-extended-query-gruppen.md)

**Schärft:** [`ARC-002`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-003`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-004`](../../../spec/architecture.md#1-komponenten-übersicht), [Ports](../../../spec/architecture.md#23-ports), [Fehlermodelle und Resilienz](../../../spec/architecture.md#5-fehlermodelle-und-resilienz), [`LH-FA-18.a`](../../../spec/spezifikation.md#lh-fa-18a--extended-query-ablauf-aufzeichnung-matching)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Im Extended Query Protocol darf ein Client Nachrichten senden, ohne auf Antworten zu warten ([`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol)); verbreitete Treiber schreiben einen Batch und lesen zugleich die Antworten. Direkt gegen PostgreSQL trägt das, weil der Server schreibt, während der Client liest. Ein Recorder dazwischen trägt es nur, wenn er in beiden Richtungen unabhängig vermittelt: Steht das Senden an den Server, weil der Server nicht liest, und liest der Recorder währenddessen nicht vom Server, blockieren sich Server und Recorder gegenseitig, und die Session endet weder beim Schließen des Clients noch beim Herunterfahren ([`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)). [ADR-0004](0004-postgresql-upstream-ist-driven-adapter.md) verlangt, dass das Port-Design Gegendruck mitdenkt, legt aber keinen Zuschnitt fest.

Zugleich hat der Record-Pfad für Extended Query einen Interaktionszustand (Gruppe im Aufbau, gesendet, `Sync` gesendet), von dem abhängt, welcher Gruppe eine Server-Nachricht gehört, ob ein Verbindungsende ein Abbruch ist und ob das Herunterfahren wartet. Nach [ADR-0003](0003-pgwire-server-ist-driving-adapter.md) und der Architektur-Sicht trägt der PGWire-Adapter keine Record-Fachlogik; nach der Sicht ist Nebenläufigkeit je Verbindung Sache des Adapters.

Annahmen: Die PGWire-Bibliothek erlaubt gleichzeitiges Lesen und Schreiben auf einer Verbindung. Die Gruppen-Semantik aus [ADR-0012](0012-extended-query-gruppen.md) bleibt unverändert.

## Entscheidung

Wir wählen: **Der Record-Pfad vermittelt je Session in zwei unabhängigen Richtungen, und der Record-Service führt den Interaktionszustand allein.**

- Der PGWire-Adapter betreibt je Session eine Richtung Client → Upstream und eine Richtung Upstream → Client als reine Transportmechanik. Jede Richtung ruft den Record-Use-Case und blockiert nur an ihrem eigenen Ziel; Gegendruck wirkt damit über die Verbindung des jeweiligen Gegenübers, ohne dass eine Richtung auf die andere wartet und ohne Puffer über die laufende Gruppe hinaus.
- Der Record-Service hält den Zustand jeder Interaktion und trifft jede Record-Entscheidung: Zuordnung der Server-Nachrichten, Einstufung eines Verbindungsendes, Warten beim Herunterfahren, Grund des Session-Endes. Der Adapter meldet Ereignisse (Client-Nachricht, Verbindungsende, `Terminate`, Schreibfehler zum Client, Herunterfahren) und führt aus, was der Service zurückgibt; er führt keinen eigenen Interaktionszustand.
- Beide Ports tragen einen Gleichzeitigkeits- und Abbruchvertrag: Senden und Empfangen derselben Session laufen gleichzeitig, und das Schließen der Session beendet jeden wartenden Aufruf mit einem Fehler. Der Vertrag steht am Port und wird mit Fakes geprüft.
- Einfache Anfragen behalten ihren bisherigen synchronen Weg; der Service lässt eine einfache Anfrage und den Empfang der Gegenrichtung nicht gleichzeitig laufen.

Die Operationen im Einzelnen, ihre Zustände und Meldungscodes stehen in der Architektur-Sicht und in der Spezifikation; bei Abweichung gelten diese.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts entscheiden; den Empfang vor dem Senden starten | kleinste Änderung, Ports bleiben | der Ablauf steht weiter im Senden; nach einem Stapel liest wieder niemand, die Blockade bleibt; der Interaktionszustand bleibt doppelt im Adapter und im Service |
| B — Senden entkoppeln durch unbegrenzten Puffer im Upstream-Adapter | Ports und Ablauf bleiben; behebt die Blockade | kein Gegendruck: der Speicher wächst mit der Eingabe des Clients, gegen [ADR-0004](0004-postgresql-upstream-ist-driven-adapter.md); der Zustand bleibt im Adapter |
| C — Der Service betreibt die Session über einen Client-Strom-Port (eigene Lese- und Schreibschleifen im Core) | gesamte Logik im Core; Adapter nur Lesen und Schreiben; mit zwei Fakes testbar | Nebenläufigkeitsmechanik wandert in den Core, gegen die Sicht; der Core ruft den Driving Adapter zurück, gegen den Kontrollfluss aus [ADR-0003](0003-pgwire-server-ist-driving-adapter.md); Replay und einfache Anfragen müssten mit umgebaut werden |
| D — Bytes durchreichen und aus einer Kopie aufzeichnen | Full-Duplex und Gegendruck ohne eigenes Zutun | von [ADR-0004](0004-postgresql-upstream-ist-driven-adapter.md) als Byte-Proxy verworfen; Gruppierung und Formprüfung wandern in den Mitschnitt |
| **E — Zwei Richtungen im Adapter, Zustand und Entscheidungen im Service, Abbruchvertrag an beiden Ports (gewählt)** | Gegendruck je Richtung über die Verbindung, keine zusätzliche Kopie; Nebenläufigkeit bleibt Infrastruktur; Adapter ohne Fachzustand; Fakes mit begrenztem Puffer bilden die Blockade deterministisch ab; einfache Anfragen unverändert | der Service muss gleichzeitige Aufrufe je Session ordnen, ohne über einem blockierenden Aufruf zu sperren; zwei nebenläufige Abläufe je Session; neue Port-Verträge brauchen Negativtests |

## Konsequenzen

- Positiv: Pipelining mit großer Eingabe und großer Ausgabe läuft über den Recorder wie direkt gegen PostgreSQL; eine Session endet beim Schließen des Clients und beim Herunterfahren, auch während Daten fließen. Die Zustands-Entscheidungen sind ohne Netz am Service testbar.
- Negativ: Die Reihenfolge, in der der Service Ereignisse beider Richtungen sieht, bestimmt die Zuordnung später Antworten einer `Flush`-Gruppe; das ist die bekannte Grenze aus [ADR-0012](0012-extended-query-gruppen.md), nicht neu. Der Replay-Pfad übernimmt diesen Zuschnitt nicht automatisch.
- Folgepflicht: Architektur-Sicht und Spezifikation nennen die Verträge der Ports; Unit-Tests mit einem Fake-Upstream, der bei vollem Puffer nicht mehr liest, und ein Integrationstest mit großer Gruppe und großer Ausgabe über den Recorder.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Unit-Tests mit Fake-Upstream (begrenzter Puffer je Richtung) | Senden und Empfangen derselben Session laufen gleichzeitig; Schließen beendet jeden wartenden Aufruf | `make test` |
| Integrationstest gegen PostgreSQL | ein Batch mit großem Parameter und großer Ausgabe läuft über den Recorder durch; danach endet der Prozess auf Signal | `make test-integration` |
| a-check mit `.a-check.yml` | Schichtregeln unverändert: der Driving Adapter importiert keine Driven Ports | `make a-check` |

## Re-Evaluierungs-Trigger

Wenn der Replay-Pfad oder einfache Anfragen denselben Zwei-Richtungs-Ablauf brauchen und eine gemeinsame Session-Abstraktion für beide Modi weniger Code ergibt (dann Option C neu prüfen), oder wenn die PGWire-Bibliothek gleichzeitiges Lesen und Schreiben auf einer Verbindung nicht mehr zusagt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-05 | Proposed | [Review-Report](../../reviews/2026-10-05-review-slice-extended-query-record.md) |
| 2026-10-05 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
