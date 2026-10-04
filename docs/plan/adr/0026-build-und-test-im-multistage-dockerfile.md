# ADR-0026: Build, Test und Integration über ein Multistage-Dockerfile

**Status:** Accepted

**Datum:** 2026-10-04

**Autor:** pt9912

**Bezug:** [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [ADR-0009](0009-implementierungssprache-go.md)

**Schärft:** [`LH-FA-16.a`](../../../spec/spezifikation.md#lh-fa-16a--containerbetrieb), [`SPEC-031`](../../../spec/spezifikation.md#6-externe-verträge)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Build und Test laufen nur in Docker (`AGENTS.md` §3.1). Mit den ersten externen Go-Abhängigkeiten braucht der Build eine Quelle der Module, die den Gate-Lauf nicht vom Netz abhängig macht und Änderungen an Abhängigkeiten nachvollziehbar hält. Das Produkt-Image (`LH-FA-16`) und das Binary sollen aus demselben Build entstehen, das Binary reproduzierbar (`LH-QA-01`). Die Record-Seite lässt sich nur gegen eine echte PostgreSQL-Instanz belegen.

## Entscheidung

Wir wählen **ein Multistage-Dockerfile** im Wurzelverzeichnis:

- Eine Stufe lädt die Module aus `go.mod` und `go.sum` (mit Netz, geprüft gegen `go.sum`); Docker cacht sie, bis sich die beiden Dateien ändern.
- Build, Test und das Bauen der Integrationstests laufen in eigenen Stufen ohne Netz.
- Dieselbe Datei baut später das Produkt-Image als letzte Stufe.
- Integrationstests laufen über ein eigenes Make-Ziel gegen ein per Digest gepinntes PostgreSQL-Image in einem eigenen Docker-Netz und hängen an `make gates`.

Die Einzelheiten (Stufen, Images, Ziele) stehen in `harness/mk/` und im Dockerfile, nicht in dieser ADR.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Abhängigkeiten als `vendor/` im Repository | Build ohne jede Netzquelle; Änderungen im Diff | großes Verzeichnis im Repository; jede Aktualisierung bläht den Diff |
| B — Modul-Cache in einem Docker-Volume | kein Verzeichnis im Repository | Cache außerhalb der Versionskontrolle; zusätzlicher Download-Schritt; kein Bezug zum Produkt-Image |
| C — Build und Test mit Netz über den Go-Proxy | einfachste Form | jeder Gate-Lauf braucht Netz; Lauf hängt am Proxy |
| **D — Multistage-Dockerfile mit Download-Stufe und netzlosen Stufen** | Netz nur bei leerem Cache oder geänderten Abhängigkeiten; Prüfsummen über `go.sum`; dieselbe Datei baut das Produkt-Image | Gate-Lauf auf frischem Rechner braucht einmal Netz (wie für die Gate-Images heute schon); BuildKit nötig |

## Konsequenzen

- Positiv: Ein Dockerfile trägt Build, Test, Integration und Produkt-Image; der Gate-Lauf ist bei warmem Cache netzlos; Änderungen an Abhängigkeiten stehen in `go.mod` und `go.sum`.
- Negativ: Ein Gate-Lauf dauert länger, weil die Integrationstests eine PostgreSQL-Instanz starten; der Docker-Build-Cache überspringt einen unveränderten Testlauf.
- Folgepflicht: `harness/README.md` nennt die Ziele und ihre Bindung; die Pins von Go- und PostgreSQL-Image werden bewusst angehoben.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Docker BuildKit | Build- und Teststufen laufen mit `RUN --network=none` | `make build`, `make test` |
| Integration | Record gegen ein gepinntes PostgreSQL-Image | `make test-integration` |

## Re-Evaluierungs-Trigger

Wenn der Gate-Lauf durch die Integrationstests zu lang wird, oder wenn ein Build ganz ohne Netzquelle gefordert wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed | — |
| 2026-10-04 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
