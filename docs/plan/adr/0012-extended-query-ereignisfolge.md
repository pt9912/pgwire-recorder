# ADR-0012: Extended Query als ereignisbasierte Interaktion

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-04`](../../../spec/lastenheft.md#lh-fa-04--transparenz-für-clients), [ADR-0007](0007-strict-replay.md)

**Schärft:** [`LH-FA-18.a`](../../../spec/spezifikation.md#lh-fa-18a--extended-query-ablauf-aufzeichnung-matching), [`SPEC-041`](../../../spec/spezifikation.md#spec-041--recording-extended-interaktion)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Verbreitete Treiber verwenden standardmäßig das Extended Query Protocol. Ohne dessen Unterstützung genügt die Umstellung von Host und Port nicht ([`LH-FA-04`](../../../spec/lastenheft.md#lh-fa-04--transparenz-für-clients)). Anders als beim Simple Query Protocol besteht eine Interaktion aus mehreren Nachrichten (`Parse`, `Bind`, `Describe`, `Execute`, `Sync`); der Client darf sie ohne Warten auf Antworten senden, und nach einem Fehler verwirft der Server Nachrichten bis zum nächsten `Sync`.

## Entscheidung

Wir wählen: Eine Extended-Interaktion wird als geordnete Ereignisfolge aus Client- und Server-Nachrichten aufgezeichnet und strict sequential wiedergegeben. Jedes eingehende Client-Ereignis muss in allen Feldern dem erwarteten entsprechen; die folgenden Server-Ereignisse werden bis zum nächsten noch nicht empfangenen Client-Ereignis freigegeben.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Extended in Simple umschreiben (Parameter einsetzen) | nur ein Interaktionsmodell | verändert, was der Client sendet; verliert Parametertypen, Format-Codes und Portal-Semantik; Aufzeichnung gibt nicht wieder, was die Anwendung getan hat |
| B — Je Nachricht eine eigene Anfrage-Antwort-Interaktion | nah am Simple-Modell | passt nicht zu Pipelining, `Flush` und Fehlerverwerfung bis `Sync`; Antworten lassen sich keiner einzelnen Nachricht zuordnen |
| **C — Ereignisfolge je `Sync`-Interaktion (gewählt)** | bildet Pipelining und Fehlerverhalten ab; gleiches strict-sequential-Prinzip | neuer Interaktionstyp im Modell und im Recording-Format |

## Konsequenzen

- Positiv: Clients mit Standardtreibern genügt die Umstellung von Host und Port.
- Negativ: Treiber mit nichtdeterministischen Statement-Namen passen nicht zur Aufzeichnung; das ist eine bekannte Grenze.
- Folgepflicht: Domain-Modell, beide PGWire-Adapter, Matcher und Recording-Schema tragen den neuen Interaktionstyp.

## Fitness Function (falls maschinell prüfbar)

Nicht maschinell prüfbar.

## Re-Evaluierungs-Trigger

Wenn verbreitete Treiber nichtdeterministische Namen verwenden und ein toleranteres Matching als eigene Anforderung aufgenommen wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
