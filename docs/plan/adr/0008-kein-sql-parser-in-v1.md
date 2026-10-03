# ADR-0008: Kein SQL-Parser in v1

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay)

**Schärft:** [`LH-FA-09.a`](../../../spec/spezifikation.md#lh-fa-09a--strict-sequential-matching)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Lastenheft schließt SQL-semantische Analyse als Voraussetzung für Record und Replay aus.

## Entscheidung

Wir wählen: SQL ist fachlich Payload des Requests und wird nicht geparst.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — SQL parsen und normalisieren | tolerantes Matching möglich | hoher Aufwand; Dialektabdeckung; widerspricht der Abgrenzung von v1 |
| B — Regelbasiertes Matching (Regex) | flexibel | Regeln werden schwer wartbar; verdeckt Abweichungen |
| **C — SQL als Payload (gewählt)** | einfach, deterministisch | semantisch äquivalente Queries gelten als verschieden |

## Konsequenzen

- Positiv: kleiner, nachvollziehbarer Matcher.
- Negativ: Textänderungen an SQL erzwingen neue Recordings.
- Folgepflicht: Verhalten im Anwenderhandbuch benennen.

## Fitness Function (falls maschinell prüfbar)

Nicht maschinell prüfbar.

## Re-Evaluierungs-Trigger

Wenn SQL-Normalisierung als eigene Anforderung ins Lastenheft aufgenommen wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
