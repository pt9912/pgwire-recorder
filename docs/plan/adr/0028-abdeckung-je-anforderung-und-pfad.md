# ADR-0028: Abdeckung je Anforderung und Pfad

**Status:** Proposed

**Datum:** 2026-10-04

**Autor:** pt9912

**Bezug:** [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md)

**Schärft:** —

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Requirements Traceability Matrix (`make doc-trace`) zählt eine Anforderung als belegt, sobald ein Nachweis sie nennt. Ein Test belegt aber meist nur einen Teil: einen der drei Pfade einer funktionalen Anforderung (Happy, Boundary, Negative) oder nur einen Fall davon. Nachweise kommen aus mehreren Arten (E2E, Unit, später Bench oder CI-Matrix), und nicht jeder Pfad ist in jeder Art belegbar.

## Entscheidung

Wir wählen: Jeder Nachweis deklariert **Anforderung und Pfad** am Test; die Deklarationen ergeben je Nachweisart eine Tabelle und eine Gesamtsicht. Die RTM liest nur die Anforderungen, deren Pfade **alle** belegt sind, gleich aus welcher Art; Teilabdeckung bleibt sichtbar, zählt aber nicht. Die Tabellen entstehen aus den Deklarationen und werden von einem Gate auf Aktualität geprüft, nicht im Gate-Lauf geschrieben.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts ändern: eine Nennung genügt | einfach | die RTM meldet Anforderungen als belegt, die nur teilweise geprüft sind |
| B — Sub-Kennungen im Lastenheft je Kriterium | feinste Bindung | dutzende neue Kennungen; ID-Schema, Hook und Gate-Konfiguration ändern sich; dieselbe Teilung wie die Pfade |
| C — Abdeckung nur über geschlossene Slices | keine Deklarationen | sagt nichts darüber, was geprüft ist |
| **D — Deklaration je Pfad, RTM nur bei vollständiger Belegung** | nutzt die Kriterien des Lastenhefts; mehrere Nachweisarten; ehrliche RTM | ob ein Test einen Pfad ausreichend belegt, bleibt Urteil dessen, der die Deklaration schreibt |

## Konsequenzen

- Positiv: Die RTM zeigt eine Anforderung erst als belegt, wenn alle Pfade belegt sind; die Gesamtsicht zeigt, was fehlt.
- Negativ: Deklarationen können ein Kriterium verfehlen; das fängt nur das Review.
- Folgepflicht: Eine Gegenprobe hält das Skript, das die Tabellen erzeugt; `harness/README.md` nennt Werkzeug und Gates.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Skript | Tabellen entsprechen den Deklarationen; jeder E2E-Test trägt eine | `make abdeckung-check` |
| Gegenprobe | Fehlformen der Deklaration werden abgelehnt | `make abdeckung-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn eine weitere Nachweisart (Bench, CI-Matrix) hinzukommt, die keine Tests trägt, oder wenn Deklarationen regelmäßig Kriterien verfehlen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
