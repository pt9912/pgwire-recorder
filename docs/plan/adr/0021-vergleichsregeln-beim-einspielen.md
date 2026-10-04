# ADR-0021: Vergleichsregeln beim Einspielen

**Status:** Superseded by [ADR-0023](0023-antwortvergleich-entscheidung.md)

**Datum:** 2026-10-04

**Autor:** pt9912

**Bezug:** [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0020](0020-antwortvergleich-beim-einspielen.md)

**Schärft:** [`LH-FA-24.a`](../../../spec/spezifikation.md#lh-fa-24a--vergleich-der-antworten), [`LH-FA-20.a`](../../../spec/spezifikation.md#lh-fa-20a--einspielen), [`SPEC-018`](../../../spec/spezifikation.md#4-fehler-codes-und-logging-felder), [`SPEC-027`](../../../spec/spezifikation.md#4-fehler-codes-und-logging-felder), [`SPEC-034`](../../../spec/spezifikation.md#spec-034--meldungscodes)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0020](0020-antwortvergleich-beim-einspielen.md) legt fest, *dass* und *was* beim Einspielen verglichen wird. Die Wechselwirkung mit der Fehlersemantik von [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md) und mit unvollständigen Aufzeichnungen bleibt dort offen: Ein Serverfehler kann zugleich `PGR-E4004` (Abbruch, Exit-Code 4) und eine Abweichung (`PGR-E5004`, Exit-Code 5) sein; eine Aufzeichnung kann mitten in der Antwort enden; bei Extended-Interaktionen ohne `Sync` ist die Zuordnung der Antworten zu einer Gruppe zeitabhängig.

## Entscheidung

Wir wählen: Mit `--compare-responses` **ersetzt der Vergleich die Fehlerregel** von [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md) für Fehlerantworten: Ein Fehler mit gleichem SQLSTATE wie aufgezeichnet gilt als erwartet, jeder andere ist eine Abweichung (`PGR-E5004`); `PGR-E4004` entfällt in diesem Modus. Der Exit-Code am Ende ist `5`; Verbindungsfehler beenden immer mit `4`. Bei Extended-Interaktionen wird die Folge aller Server-Nachrichten ohne Gruppengrenzen verglichen. Eine unvollständige Aufzeichnung wird verglichen, soweit sie reicht. Der Befehl im `CommandComplete` zählt ohne abschließende Zahlen.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts festlegen, Umsetzung entscheidet | kein Aufwand | Exit-Codes und Meldungen hängen vom Zufall der Reihenfolge ab; nicht testbar |
| B — Fehlerregel von 0017 bleibt, Vergleich kommt zusätzlich | keine Änderung an 0017 | derselbe Fehler kann `PGR-E4004` und `PGR-E5004` heißen; „zuerst aufgetretene Ursache" hebt die Unterscheidbarkeit der Exit-Codes auf |
| C — Unvollständige Aufzeichnungen nicht vergleichen | einfach | der Lauf kann mit Exit-Code 0 enden, obwohl nicht geprüft wurde |
| **D — Vergleich ersetzt die Fehlerregel, Vergleich bis zum Ende der Aufzeichnung** | eindeutige Exit-Codes (5 für Abweichung, 4 für Verbindung); nichts bleibt ungeprüft | `--allow-recorded-errors` ist mit Vergleich überflüssig; bei Extended ist die Gruppe in der Meldung nicht bestimmbar |

## Konsequenzen

- Positiv: Ein Lauf mit Vergleich endet mit 0, 4 oder 5, und jede Zahl hat genau eine Bedeutung.
- Negativ: Mit Vergleich gilt jeder aufgezeichnete Fehler als erwartet, ohne dass der Anwender ihn einzeln zulässt.
- Folgepflicht: Tests decken Fehler mit gleichem und anderem SQLSTATE, `Flush`-Gruppen und abgebrochene Aufzeichnungen ab; das Handbuch beschreibt die Rangfolge.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | nicht maschinell prüfbar; die Tests des Slice für den Antwortvergleich tragen die Regeln | — |

## Re-Evaluierungs-Trigger

Wenn Gruppen in der Meldung bestimmbar sein müssen oder Anwender einzelne aufgezeichnete Fehler weiter getrennt zulassen wollen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed | — |
| 2026-10-04 | Accepted | — |
| 2026-10-04 | Superseded by ADR-0022 | [ADR-0022](0022-vergleichsregeln-beim-einspielen-praezisiert.md) |
| 2026-10-04 | Superseded by ADR-0023 (über ADR-0022) | [ADR-0023](0023-antwortvergleich-entscheidung.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
