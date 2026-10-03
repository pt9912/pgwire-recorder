# ADR-0017: Einspielen: sequenziell, Antworten verworfen, Fehlersemantik

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0007](0007-strict-replay.md), [ADR-0012](0012-extended-query-gruppen.md)

**Schärft:** [`LH-FA-20.a`](../../../spec/spezifikation.md#lh-fa-20a--einspielen), [`SPEC-034`](../../../spec/spezifikation.md#spec-034--meldungscodes), [`ARC-004`](../../../spec/architecture.md#1-komponenten-übersicht)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`play` spielt die Client-Anfragen einer Aufzeichnung gegen einen echten PostgreSQL-Server ein, zum Beispiel um Last oder Datenänderungen für einen Test zu erzeugen. Der Server antwortet auf eigener Datenlage; ein Abgleich mit der Aufzeichnung wäre bei jeder Abweichung ein Fehler und hieße, den Inhalt der Antworten zu deuten ([ADR-0008](0008-kein-sql-parser-in-v1.md)). Mehrere Sessions parallel würden die Reihenfolge der Wirkung auf der Datenbank unbestimmt machen.

## Entscheidung

Wir wählen: `play` spielt **sequenziell** ein, Sessions in der Reihenfolge ihrer `id`, Anfragen einer Session nacheinander, und wartet nach jeder Interaktion auf `ReadyForQuery`. Die Serverantworten werden gelesen und verworfen, mit Ausnahme einer `ErrorResponse`; ein Vergleich findet nicht statt. Eine `ErrorResponse` beendet das Einspielen mit `PGR-E4004` (Exit-Code `4`), außer `--continue-on-error` (weiter, am Ende Exit-Code `4`) oder `--allow-recorded-errors` (nur, wo die Aufzeichnung ebenfalls einen Fehler enthält). Verbindungsfehler beenden das Einspielen immer (`PGR-E4002`, `PGR-E4003`). Der Prozess folgt nicht den Fehlerebenen des Record-Modus.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts verwerfen, Antworten mit der Aufzeichnung vergleichen | prüft das Verhalten der Datenbank | jede Abweichung der Datenlage ist ein Fehler; verlangt Deutung des Inhalts |
| B — Sessions parallel einspielen | nähert sich der aufgezeichneten Last an | Reihenfolge der Wirkung unbestimmt; schlecht reproduzierbar |
| C — Fehler stets ignorieren | nie ein Abbruch | verdeckt, dass die Datenlage nicht zur Aufzeichnung passt |
| **D — Sequenziell, Antworten verwerfen, Fehler als Abbruch mit zwei Lockerungen** | bestimmte Reihenfolge; einfache Semantik; Fehler sichtbar | keine Prüfung der Antwortinhalte; keine Lastsimulation durch Parallelität |

## Konsequenzen

- Positiv: Der Lauf ist bei gleicher Datenlage reproduzierbar; der Exit-Code sagt, ob ein Fehler auftrat.
- Negativ: Wer Antworten prüfen oder Last parallel erzeugen will, braucht ein anderes Werkzeug oder eine spätere Erweiterung.
- Folgepflicht: Tests decken je Fehlerklasse und je Lockerung ab, auch für Extended-Gruppen ([ADR-0012](0012-extended-query-gruppen.md)).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Tests | Fehlerklassen und Exit-Codes des Einspielens entsprechen der Spezifikation | `make gates` |

## Re-Evaluierungs-Trigger

Wenn Parallelität oder ein Vergleich der Antworten gefordert wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
