# ADR-0020: Antwortvergleich beim Einspielen

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0008](0008-kein-sql-parser-in-v1.md), [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md)

**Schärft:** [`LH-FA-24.a`](../../../spec/spezifikation.md#lh-fa-24a--vergleich-der-antworten), [`LH-FA-20.a`](../../../spec/spezifikation.md#lh-fa-20a--einspielen), [`SPEC-034`](../../../spec/spezifikation.md#spec-034--meldungscodes), [`ARC-002`](../../../spec/architecture.md#1-komponenten-übersicht)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md) verwirft die Serverantworten beim Einspielen und nennt als Re-Evaluierungs-Trigger einen geforderten Vergleich. Der Trigger ist eingetreten: Wer `play` gegen einen Adapter oder ein Testsystem fährt, will wissen, ob dessen Antworten der Aufzeichnung entsprechen. Ein Vergleich der Zeilenwerte wäre bei jeder Abweichung der Datenlage ein Fehler und verlangte Deutung des Inhalts ([ADR-0008](0008-kein-sql-parser-in-v1.md)). Die Vorgabe von 0017 (verwerfen) bleibt der Default; diese ADR ergänzt sie um einen Wunsch.

## Entscheidung

Wir wählen: `play --compare-responses` vergleicht nach jeder Interaktion die Struktur der Serverantwort und Fehler mit der Aufzeichnung: Nachrichtentypen in Reihenfolge, Spaltenbeschreibung (Anzahl, Name, Typ-OID), den Befehl im `CommandComplete`, den SQLSTATE einer `ErrorResponse` und den Transaktionsstatus. Zeilen (`data_row`), Zeilenzahlen, Hinweise und Parameterstatus werden nicht verglichen. Eine Abweichung ist `PGR-E5004` mit Exit-Code `5` und folgt der Regel für Fehlerantworten ([ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md)): sie beendet das Einspielen, mit `--continue-on-error` läuft es weiter. Der Vergleich ist reine Fachlogik im Play-Service.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun, Antworten verwerfen | einfach; keine Deutung des Inhalts | ein Adapter oder Testsystem lässt sich nicht gegen die Aufzeichnung prüfen |
| B — Vergleich der gesamten Antwort einschließlich Zeilenwerten | strengste Prüfung | jede Abweichung der Datenlage, jeder Zeitstempel, jede generierte Id ist ein Fehler |
| C — Vergleich nur der Fehler | einfach, robust | übersieht abweichende Struktur (Spalten, Befehl) |
| **D — Vergleich von Struktur und Fehlern, ohne Zeilenwerte** | prüft Vertragstreue der Antworten; robust gegen die Datenlage | Abweichungen in Zeilenwerten bleiben unerkannt |

## Konsequenzen

- Positiv: `play` taugt als Prüfer für Adapter und Testsysteme; Exit-Code `5` trennt Abweichungen von Serverfehlern (Exit-Code `4`).
- Negativ: Antworten werden bis zum Ende der Interaktion gehalten; die Normalisierung (was nicht verglichen wird) ist Teil der Zusage und schwer zu ändern.
- Folgepflicht: Handbuch beschreibt, was verglichen wird; Tests decken jede Art der Abweichung ab, auch für Extended-Gruppen ([ADR-0012](0012-extended-query-gruppen.md)).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | die PGWire-Bibliothek ist nur in den beiden PGWire-Adaptern erlaubt (`tech`-Regel), der Vergleich im Core importiert sie nicht | `make a-check` |

## Re-Evaluierungs-Trigger

Wenn ein Vergleich von Zeilenwerten oder Toleranzregeln für einzelne Felder gefordert wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
