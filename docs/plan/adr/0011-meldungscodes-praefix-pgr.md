# ADR-0011: Meldungscodes mit Präfix PGR

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-FA-14`](../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)

**Schärft:** [`SPEC-034`](../../../spec/spezifikation.md#spec-034--meldungscodes)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Lastenheft verlangt unterscheidbare Fehlerklassen und verständliche Diagnoseausgaben. Die Spezifikation bildet Fehlerklassen auf Exit Codes ab; mehrere Ursachen teilen sich einen Exit Code, und Warnungen haben keine Kennung. Skripte und CI-Auswertungen brauchen eine stabile Kennung je Ursache, die unabhängig vom Fehlertext bleibt. Die Schwesterprojekte verwenden dafür Codes mit Schwere-Buchstabe und vierstelliger Nummer.

## Entscheidung

Wir wählen: Meldungscodes der Form `PGR-<S><NNNN>` (`S` = `E` Fehler, `W` Warnung, `I` reserviert). Bei Fehlern ist die erste Ziffer der Exit Code der Klasse, `000` der Rückfall; bei Warnungen ist sie der Bereich. Der Fehlertext beginnt mit `<klasse> [<code>]: `. Codes werden nie neu belegt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nur Exit Codes | keine neue Kennung | mehrere Ursachen je Code; Warnungen ohne Kennung |
| B — Freier Fehlertext | maximal flexibel | nicht stabil auswertbar; Textänderung bricht Auswertungen |
| **C — Meldungscodes `PGR-…` (gewählt)** | stabile Kennung je Ursache und Warnung; Klasse folgt aus dem Code; Prozessausgang unberührt | Pflege der Code-Tabelle und des Katalogs |

## Konsequenzen

- Positiv: Auswertbare Fehler und Warnungen; kein Eingriff in die Exit-Code-Semantik.
- Negativ: Code-Tabelle im Quelltext und Katalog in der Betriebsdokumentation müssen gleich bleiben.
- Folgepflicht: Ein Gate, das Code-Tabelle und Katalog abgleicht, sobald beide existieren.

## Fitness Function (falls maschinell prüfbar)

Noch kein Gate vorhanden; der Abgleich von Code-Tabelle und Katalog wird mit der Umsetzung als Sensor festgelegt.

## Re-Evaluierungs-Trigger

Wenn Codes mehrerer Klassen in einer Fehlerkette auftreten und die Vorrangfolge festgelegt werden muss, oder wenn die Zahl der Ursachen eine vierstellige Nummer je Klasse übersteigt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
