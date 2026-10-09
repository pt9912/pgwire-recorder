# ADR-0037: YAML-Bibliothek im CLI-Adapter auch für die Anzeige der Konfigurationsdatei

**Status:** Proposed

**Datum:** 2026-10-09

**Autor:** pt9912

**Bezug:** [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0027](0027-yaml-bibliothek.md), [ADR-0036](0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md)

**Schärft:** [`ARC-005`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-013`](../../../spec/architecture.md#3-externe-abhängigkeiten), [§2 *Schichten und Constraints*](../../../spec/architecture.md#2-schichten-und-constraints)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0036](0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) lässt die YAML-Bibliothek im CLI-Adapter zu, „nur dort, zum Lesen der Konfigurationsdatei“, und die Sicht sagt dasselbe (`ARC-013`, §2 *Zusätzliche Einschränkungen*). `config show` gibt nach `LH-FA-17.a` (*Anzeige*) den Inhalt derselben Datei als YAML aus, mit zwei Leerzeichen Einzug, in der Reihenfolge der Datei und ohne Kommentare. Diese Ausgabe ist Schreiben von YAML, kein Lesen; sie liegt außerhalb des Wortlauts von [ADR-0036](0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), aber nicht außerhalb seines Zwecks: Der Re-Evaluierungs-Trigger dort nennt „etwas anderes als die Konfigurationsdatei“.

## Entscheidung

Die YAML-Bibliothek ist im CLI-Adapter auch zulässig, um die gelesene Konfigurationsdatei für `config show` wieder als YAML auszugeben; Eingabe des Kodierers ist nur der Baum, den das Lesen derselben Datei geliefert hat. Für jeden anderen Zweck bleibt sie dort unzulässig. Sonst gilt [ADR-0036](0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) unverändert, auch dass kein Typ der Bibliothek den CLI-Adapter verlässt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — Kodierer der Bibliothek für die Anzeige** | eine Stelle für YAML im Adapter; Quoting, Einzug und Reihenfolge kommen aus derselben Bibliothek, die die Datei gelesen hat | der zugelassene Zweck wird weiter |
| B — eigener Ausgeber für den Baum, ohne die Bibliothek | Wortlaut von [ADR-0036](0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) bleibt | eigene Regeln für Quoting und Sonderzeichen, eine zweite YAML-Implementierung im Repo mit eigenem Fehlerbild |
| C — Text der Datei ausgeben, Kommentare entfernen | keine Serialisierung | Kommentare lassen sich im Text nicht ohne YAML-Lexer erkennen (`#` in Anführungszeichen); verletzt die Form aus `LH-FA-17.a` (Einzug) |

## Konsequenzen

- Positiv: Die Anzeige folgt der Form aus `LH-FA-17.a` ohne zweite YAML-Implementierung.
- Negativ: Wie bei [ADR-0036](0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) prüft kein Gate den Zweck eines Imports im CLI-Adapter; das bleibt Review.
- Folgepflicht: Mit `Accepted` nennt die Sicht in `ARC-013` und in §2 *Zusätzliche Einschränkungen* für `ARC-005` Lesen und Anzeigen der Konfigurationsdatei, im Slice, der `config show` liefert, vor seiner Closure.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | unverändert: `go.yaml.in/yaml` und `gopkg.in/yaml` nur im Recording-Adapter und im CLI-Adapter (`tech`-Regel) | `make a-check` |

## Re-Evaluierungs-Trigger

Wenn der CLI-Adapter YAML für etwas anderes als Lesen und Anzeigen der Konfigurationsdatei braucht, oder wenn die Anzeige einen Baum ausgeben soll, der nicht aus dem Lesen der Datei stammt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-09 | Proposed | slice-v1-abschluss-konfigurationsdatei |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
