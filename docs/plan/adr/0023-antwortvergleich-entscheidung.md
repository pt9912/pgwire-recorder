# ADR-0023: Antwortvergleich beim Einspielen ersetzt die Fehlerregel

**Status:** Proposed

**Datum:** 2026-10-04

**Autor:** pt9912

**Bezug:** [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0020](0020-antwortvergleich-beim-einspielen.md)

**Schärft:** [`LH-FA-24.a`](../../../spec/spezifikation.md#lh-fa-24a--vergleich-der-antworten), [`LH-FA-20.a`](../../../spec/spezifikation.md#lh-fa-20a--einspielen)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Beim Einspielen kann derselbe Serverfehler zwei Ursachen haben: die Fehlerregel ohne Vergleich ([ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md)) und den Antwortvergleich ([ADR-0020](0020-antwortvergleich-beim-einspielen.md)). Ohne Entscheidung, welche gilt, sind Meldungscode und Exit-Code nicht unterscheidbar. Die Regeln im Einzelnen sind Sache der Spezifikation und ändern sich mit ihr; diese ADR hält nur fest, was entschieden wurde und warum. Die vorige ADR zu diesem Thema (siehe Geschichte) wiederholte Einzelregeln und wurde dadurch ungenau.

## Entscheidung

Wir wählen: Mit `--compare-responses` **entscheidet der Vergleich über Fehlerantworten** und ersetzt dafür die Fehlerregel von [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md). Eine Abweichung hat einen eigenen Meldungscode und einen eigenen Exit-Code, getrennt vom Verbindungsfehler. Unvollständige Aufzeichnungen werden verglichen, soweit sie reichen. Verglichen wird die Struktur der Antworten und Fehler, nicht ihre Werte. Alle Einzelregeln (Ausnahmen, Rangfolge, Meldungsinhalt) stehen in der Spezifikation und gelten dort; bei Abweichung gilt die Spezifikation.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts entscheiden, die Umsetzung wählt | kein Aufwand | Exit-Codes hängen vom Zufall der Reihenfolge ab; nicht testbar |
| B — Fehlerregel von 0017 bleibt, Vergleich kommt zusätzlich | keine Änderung an der bestehenden Regel | derselbe Fehler kann zwei Codes tragen; die Exit-Codes sind nicht mehr unterscheidbar |
| C — Alle Regeln in der ADR festschreiben | Entscheidung und Regeln an einer Stelle | jede Korrektur an der Spezifikation lässt die ADR veralten; angenommene ADRs lassen sich nicht anpassen |
| **D — Vergleich entscheidet, Einzelregeln in der Spezifikation** | eindeutige Codes; die ADR bleibt gültig, wenn die Spezifikation Einzelheiten schärft | wer die Regeln sucht, liest die Spezifikation |

## Konsequenzen

- Positiv: Mit Vergleich bedeutet der Exit-Code der Abweichung etwas anderes als der des Verbindungsfehlers; Änderungen an Einzelregeln brauchen keine neue ADR.
- Negativ: Die ADR allein genügt nicht, um das Verhalten zu kennen.
- Folgepflicht: Die Tests des Slice für den Antwortvergleich tragen die Regeln der Spezifikation.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | nicht maschinell prüfbar | — |

## Re-Evaluierungs-Trigger

Wenn Anwender einzelne aufgezeichnete Fehler getrennt zulassen wollen oder alle Abweichungen einer Interaktion gemeldet werden sollen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed; Supersedes ADR-0022 | [ADR-0022](0022-vergleichsregeln-beim-einspielen-praezisiert.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
