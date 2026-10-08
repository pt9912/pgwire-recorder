# ADR-0036: YAML-Bibliothek auch für die Konfigurationsdatei im CLI-Adapter

**Status:** Proposed

**Datum:** 2026-10-08

**Autor:** pt9912

**Bezug:** [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0001](0001-hexagonale-architektur.md), [ADR-0014](0014-konfigurationsdatei.md), [ADR-0027](0027-yaml-bibliothek.md)

**Schärft:** [`ARC-005`](../../../spec/architecture.md#1-komponenten-übersicht), [`ARC-013`](../../../spec/architecture.md#3-externe-abhängigkeiten), [§2 *Schichten und Constraints*](../../../spec/architecture.md#2-schichten-und-constraints)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Konfigurationsdatei ist YAML ([ADR-0014](0014-konfigurationsdatei.md)), und die Sicht gibt ihr Lesen dem CLI-Adapter: `ARC-005` führt Argumente, Umgebung und Konfigurationsdatei zu einer Konfiguration zusammen, einschließlich Platzhalter und Passwort. Dieselbe Sicht beschränkt die YAML-Bibliothek auf den Recording-Adapter (`ARC-013`, §2 *Zusätzliche Einschränkungen*), und [ADR-0027](0027-yaml-bibliothek.md) bindet die `tech`-Regel des Architektur-Gates daran. Beides zusammen lässt die Konfigurationsdatei ohne Leser: Der CLI-Adapter darf die Bibliothek nicht importieren, der Recording-Adapter darf vom CLI-Adapter nicht gerufen werden.

## Entscheidung

Die YAML-Bibliothek aus [ADR-0027](0027-yaml-bibliothek.md) ist zusätzlich im CLI-Adapter (`internal/adapters/driving/cli`) zulässig, und nur dort, zum Lesen der Konfigurationsdatei. Kein Typ der Bibliothek verlässt den CLI-Adapter. Die `tech`-Regel in `.a-check.yml` führt für beide Modulpfade beide Adapter; die Sicht nennt `ARC-005` neben `ARC-008`. Die Wahl der Bibliothek bleibt die von [ADR-0027](0027-yaml-bibliothek.md).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Datei im Recording-Adapter lesen, über einen Port | Bibliothek bleibt an einer Stelle | ein Driven Port für eine Eingabe des Driving Adapters; Port und Service ändern sich für eine Startkonfiguration (dritte Schicht) |
| B — eigener Driven Adapter für die Konfiguration | klare Zuständigkeit | ebenfalls ein Port im Kern; widerspricht `ARC-005` |
| C — Datei in der Composition Root lesen, Baum an den CLI-Adapter geben | Composition Root darf alles importieren | Prüfung und Fehlerstellen der Datei liegen dann in zwei Paketen; `ARC-005` bleibt unerfüllt |
| **D — Bibliothek auch im CLI-Adapter** | entspricht `ARC-005`; eine Stelle für Zusammenführung und Fehler | die Bibliothek steht in zwei Adaptern |

## Konsequenzen

- Positiv: Die Konfigurationsdatei wird dort gelesen, wo die Sicht die Zusammenführung verortet; Kern und Ports bleiben unberührt.
- Negativ: Die `tech`-Regel ist weiter als bisher; ein Import der Bibliothek im CLI-Adapter zu einem anderen Zweck fällt dem Gate nicht auf.
- Folgepflicht: `.a-check.yml` und die Sicht (`ARC-013`, §2 *Zusätzliche Einschränkungen*, §6) werden vor dem ersten Code-Commit, der die Bibliothek im CLI-Adapter importiert, nachgezogen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | `go.yaml.in/yaml` und `gopkg.in/yaml` nur im Recording-Adapter und im CLI-Adapter (`tech`-Regel) | `make a-check` |

## Re-Evaluierungs-Trigger

Wenn ein dritter Adapter YAML lesen muss oder der CLI-Adapter die Bibliothek für etwas anderes als die Konfigurationsdatei braucht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-08 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
