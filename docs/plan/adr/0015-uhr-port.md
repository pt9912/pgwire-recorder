# ADR-0015: Zeit über einen Uhr-Port

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen), [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus), [ADR-0001](0001-hexagonale-architektur.md)

**Schärft:** [`LH-FA-21.a`](../../../spec/spezifikation.md#lh-fa-21a--zeitangaben), [`ARC-004`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-015`](../../../spec/architecture.md#3-externe-abhängigkeiten)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`--record-timing` misst Abstände zwischen Client-Nachrichten, `--keep-timing` wartet beim Einspielen bis zu einem Zeitpunkt. Beides hängt von der Zeit ab. Die Hexagon-Regel ([ADR-0001](0001-hexagonale-architektur.md)) hält den Core frei von Infrastruktur; Tests des Zeitverhaltens (relativ, absolut, Verspätung) dürfen nicht real warten und nicht flackern.

## Entscheidung

Wir wählen einen **Driven Port „Uhr"** mit zwei Leistungen: die monotone Zeit lesen und bis zu einem Zeitpunkt warten (abbrechbar durch den Kontext). Der Composition Root stellt die Systemuhr bereit; Tests verwenden eine Fake-Uhr, die Zeit nur auf Anweisung vorrückt. Der Core ruft `time.Now` und `time.Sleep` nicht direkt auf.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Direkter Zugriff auf die Systemzeit im Core | kein zusätzlicher Port | Tests warten real oder flackern; Core hängt an Infrastruktur |
| B — Zeitfunktionen als Parameter der Services | kein eigener Port-Typ | Signaturen wachsen; die Zusage steht nicht an einer Stelle |
| **C — Uhr als Driven Port** | Fake-Uhr macht Zeittests deterministisch; passt zur Port-Struktur | ein weiterer Port; Abbruch beim Warten muss im Port zugesagt sein |

## Konsequenzen

- Positiv: Zeitverhalten ist ohne Wartezeit testbar; der Abbruch eines Wartens durch ein Signal ist am Port prüfbar.
- Negativ: Jede Stelle mit Zeitbezug geht durch den Port.
- Folgepflicht: Architekturregeln halten den Core frei von Zeitaufrufen der Standardbibliothek.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | der Core ruft Zeitfunktionen nur über den Uhr-Port auf | `make a-check` |

## Re-Evaluierungs-Trigger

Wenn weitere Zeitquellen nötig werden (zum Beispiel Zeitstempel in der Aufzeichnung).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
