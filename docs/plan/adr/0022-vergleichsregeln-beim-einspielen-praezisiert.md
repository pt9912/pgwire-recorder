# ADR-0022: Vergleichsregeln beim Einspielen, präzisiert

**Status:** Accepted

**Datum:** 2026-10-04

**Autor:** pt9912

**Bezug:** [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0020](0020-antwortvergleich-beim-einspielen.md)

**Schärft:** [`LH-FA-24.a`](../../../spec/spezifikation.md#lh-fa-24a--vergleich-der-antworten), [`LH-FA-20.a`](../../../spec/spezifikation.md#lh-fa-20a--einspielen), [`SPEC-017`](../../../spec/spezifikation.md#4-fehler-codes-und-logging-felder), [`SPEC-018`](../../../spec/spezifikation.md#4-fehler-codes-und-logging-felder), [`SPEC-022`](../../../spec/spezifikation.md#4-fehler-codes-und-logging-felder), [`SPEC-027`](../../../spec/spezifikation.md#4-fehler-codes-und-logging-felder), [`SPEC-034`](../../../spec/spezifikation.md#spec-034--meldungscodes)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die zuvor angenommene ADR zu den Vergleichsregeln beim Einspielen (siehe Geschichte) legt sie fest, trägt aber Aussagen, die über die Spezifikation hinausgehen: Mit Vergleich gelte „jeder aufgezeichnete Fehler" als erwartet (die Spezifikation verlangt gleichen SQLSTATE an der aufgezeichneten Stelle), ein Lauf ende mit „0, 4 oder 5" (auch 2, 3 und 6 bleiben möglich), und seine Folgepflichten und die Fitness-Function-Zeile nennen Tests, die der Slice nicht vollständig trug. Die Entscheidung selbst bleibt; diese ADR ersetzt den Text.

## Entscheidung

Wir wählen: Mit `--compare-responses` **ersetzt der Vergleich die Fehlerregel** von [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md) für Fehlerantworten. Eine `ErrorResponse` mit gleichem SQLSTATE an der aufgezeichneten Stelle gilt als erwartet, jede andere ist eine Abweichung (`PGR-E5004`, Exit-Code `5`); `PGR-E4004` entfällt in diesem Modus, und `--allow-recorded-errors` hat dann keine Wirkung. Verbindungsfehler beenden immer sofort mit Exit-Code `4`. Bei Extended-Interaktionen wird die Folge aller Server-Nachrichten ohne Gruppengrenzen verglichen. Eine unvollständige Aufzeichnung wird nach der Normalisierung verglichen, soweit sie reicht; was der Server danach sendet, ist keine Abweichung. Der Befehl im `CommandComplete` zählt ohne abschließende Zahlen. Je Interaktion wird die erste Abweichung gemeldet; am Ende des Laufs steht eine Zusammenfassung mit der Zahl der verglichenen und der abweichenden Interaktionen.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Die vorige ADR unverändert lassen | kein Aufwand | die ADR sagt mehr und anderes als die Spezifikation |
| B — Fehlerregel von 0017 bleibt, Vergleich kommt zusätzlich | keine Änderung an 0017 | derselbe Fehler kann `PGR-E4004` und `PGR-E5004` heißen; die Exit-Codes sind nicht mehr unterscheidbar |
| C — Alle Abweichungen einer Interaktion einzeln melden | vollständigere Meldung | längere Ausgabe; ein Fehler zieht mehrere Folgemeldungen nach sich |
| **D — Vergleich ersetzt die Fehlerregel, erste Abweichung je Interaktion, Zusammenfassung am Ende** | eindeutige Exit-Codes; überschaubare Ausgabe; Umfang der Abweichungen im Ergebnis sichtbar | weitere Abweichungen einer Interaktion bleiben ungemeldet |

## Konsequenzen

- Positiv: Mit Vergleich bedeutet Exit-Code `5` eine Abweichung und `4` einen Verbindungsfehler; die Zusammenfassung zeigt den Umfang.
- Negativ: Mit Vergleich gilt ein Fehler nur dann als erwartet, wenn die Aufzeichnung ihn an der Stelle mit gleichem SQLSTATE enthält; wer ihn einzeln zulassen will, hat dafür keine Option.
- Folgepflicht: Das Handbuch beschreibt Rangfolge und Zusammenfassung. Der Slice für den Antwortvergleich trägt die Tests je Regel (gleicher und anderer SQLSTATE, `--allow-recorded-errors`, Signal, Verbindungsfehler nach Abweichung, `Flush`-Gruppen, unvollständige Aufzeichnung).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | nicht maschinell prüfbar; die Tests des Slice für den Antwortvergleich tragen die Regeln | — |

## Re-Evaluierungs-Trigger

Wenn Gruppen in der Meldung bestimmbar sein müssen, Anwender einzelne aufgezeichnete Fehler getrennt zulassen wollen oder alle Abweichungen einer Interaktion gemeldet werden sollen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed; Supersedes ADR-0021 | [ADR-0021](0021-vergleichsregeln-beim-einspielen.md) |
| 2026-10-04 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
