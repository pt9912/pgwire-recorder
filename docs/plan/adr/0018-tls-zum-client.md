# ADR-0018: TLS-Terminierung zum Client

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-23`](../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client), [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [ADR-0003](0003-pgwire-server-ist-driving-adapter.md), [ADR-0010](0010-verwendung-von-pgproto3.md)

**Schärft:** [`LH-FA-23.a`](../../../spec/spezifikation.md#lh-fa-23a--tls-zum-client), [`LH-FA-05.c`](../../../spec/spezifikation.md#lh-fa-05c--tls), [`ARC-006`](../../../spec/architecture.md#1-komponenten-übersicht)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Clients, die TLS verlangen, können sich bisher nicht mit dem Recorder verbinden: Er lehnt `SSLRequest` ab. Testumgebungen, in denen die Anwendung oder ein vorgeschalteter Dienst nur verschlüsselt spricht, brauchen TLS zum Client. Ein TLS-Passthrough ohne Protokolleinsicht genügt nicht, weil der Recorder die Nachrichten lesen muss. Zum Upstream im Record-Modus bleibt die Verbindung unverschlüsselt.

## Entscheidung

Wir wählen: Der PGWire-Server-Adapter terminiert TLS **auf Wunsch**, mit Zertifikat und Schlüssel des Anwenders (`--tls-cert`, `--tls-key`). Ohne Konfiguration gilt das bisherige Verhalten (`N`). Mit Konfiguration werden unverschlüsselte Clients abgewiesen (`PGR-E6003`), sofern der Anwender sie nicht mit `--allow-plaintext` zulässt. Der Core und die Aufzeichnung kennen TLS nicht; dieselbe Aufzeichnung ist über beide Verbindungsarten wiedergebbar. Prüfung von Client-Zertifikaten und TLS zum Upstream sind nicht Teil der Entscheidung.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun, `N` beibehalten | keine Zertifikatsverwaltung | Clients mit Pflicht-TLS sind ausgeschlossen |
| B — TLS-Passthrough | keine Schlüssel im Recorder | der Recorder sieht die Nachrichten nicht; kein Recording, kein Replay |
| C — TLS durch einen vorgeschalteten Dienst (Sidecar) | keine Änderung am Recorder | zusätzlicher Dienst in jeder Testumgebung; kein Beitrag zu „eine Binärdatei" |
| **D — Terminierung im PGWire-Adapter, auf Wunsch** | ein Werkzeug genügt; Core bleibt TLS-frei; Default unverändert | Zertifikat und Schlüssel müssen bereitgestellt werden; zusätzliche Fehlerfälle |

## Konsequenzen

- Positiv: Verschlüsselnde Clients lassen sich aufzeichnen und bedienen; die Aufzeichnung bleibt unabhängig von der Verbindungsart.
- Negativ: Der Schlüssel liegt im Zugriff des Recorders (Dateirechte beim Anwender); der Recorder bietet keine Zertifikatsverwaltung.
- Folgepflicht: Handbuch beschreibt die Optionen und den Betrieb mit selbst erzeugten Zertifikaten; Tests decken Aushandlung, Ablehnung unverschlüsselter Clients und Handshake-Fehler.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | `crypto/tls` nur in den beiden PGWire-Adaptern (Server und Upstream), nie im Core | `make a-check` |

## Re-Evaluierungs-Trigger

Wenn Client-Zertifikate geprüft werden müssen oder TLS zum Upstream im Record-Modus gefordert wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
