# ADR-0043: Binaries je Zielplattform und Container-Image als ein Gate mit Smoke

**Status:** Proposed

**Datum:** 2026-10-17

**Autor:** Architect

**Bezug:** [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md)

**Schärft:** [`SPEC-055`](../../../harness/vertraege/SPEC-055-container-und-binaries.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0026](0026-build-und-test-im-multistage-dockerfile.md) lässt dieselbe Datei das Produkt-Image bauen. Das Lastenheft verlangt mehr: das Binary für sechs Plattformen, reproduzierbar (zwei Builds desselben Quellstands, dieselbe Prüfsumme), und ein Image für `linux/amd64` und `linux/arm64`, das Aufzeichnen und Wiedergeben besteht. Heute baut `make build` das Image für die Plattform des Hosts; weder die Summen noch die anderen Plattformen noch ein Lauf des Images stehen unter einem Gate. Das Release-Verfahren soll dieselben Binaries mit `SHA256SUMS` an den Release hängen und das Image veröffentlichen; es braucht ein Ziel, das beides liefert, statt es ein zweites Mal zu bauen.

## Entscheidung

Wir wählen **Option C**: Ein Werkzeug `make binaries` baut die sechs Binaries mit `SHA256SUMS` aus einer eigenen Stufe des `Dockerfile`. Ein Gate `make test-container` baut sie zweimal ohne Cache der Bau-Stufen und vergleicht die Summen, prüft den Kopf jedes Binaries, baut das Image je Linux-Plattform mit demselben Binary wie `binaries` und fährt auf der Plattform des Docker-Servers den Smoke (`version`, Aufzeichnen und Wiedergeben ohne PostgreSQL in einem internen Netz). Die Einzelregeln stehen in [`SPEC-055`](../../../harness/vertraege/SPEC-055-container-und-binaries.md). Reproduzierbar ist das Binary, nicht das Image. Die Manifestliste und die Registries gehören zur Veröffentlichung, nicht zu diesem Gate.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Summen nur im Release-Lauf prüfen | kein Zusatz im Gate-Lauf | die Zusage von `LH-QA-01` hätte bis zum ersten Tag kein Gate (`AGENTS.md` §3.11); ein Fehler zeigte sich beim Release |
| B — Image reproduzierbar machen (feste Zeitstempel, gleicher Digest) | strengere Aussage | das Lastenheft verlangt es nur für das Binary; Aufwand und Pflege ohne Anforderung |
| **C — ein Gate: Summen über zwei Builds, Kopf, Image je Plattform, Smoke; Binaries als Werkzeug** | jede Zusage hat ein Gate; die Veröffentlichung baut nicht ein zweites Mal; das Image trägt dasselbe Binary | zwölf Builds und ein Smoke im Gate-Lauf; die Manifestliste ist erst bei der Veröffentlichung belegt |
| D — Manifestliste lokal bauen (buildx mit Container-Treiber) | belegt die Liste vor dem Release | der Standard-Treiber kann es nicht; ein zusätzlicher Builder im Gate-Lauf |

## Konsequenzen

- Positiv: `LH-QA-01` (Bereitstellung) und Abnahmeszenario 9 haben ein Gate; das Release-Verfahren ruft `make binaries`.
- Negativ: Die Gate-Laufzeit wächst um einige Minuten (zwölf Builds der Bau-Stufe, zwei Image-Builds, ein Smoke).
- Akzeptiertes Negativ: Das Gate vergleicht zwei Builds auf demselben Docker-Server und prüft die Pins der Eingangs-Images nicht; der Grund: Anhebungen sind bewusste Commits ([ADR-0026](0026-build-und-test-im-multistage-dockerfile.md)), und ein Lauf auf zweitem Host braucht eine zweite Umgebung ohne Anforderung. Die Manifestliste belegt erst der Lauf der Veröffentlichung.
- Folgepflicht: `harness/README.md` §Sensors führt die Ziele mit Bindung an diese ADR; die Ausnahmeliste der schnellen Kette ([ADR-0038](0038-gate-kette-voll-und-schnell.md)) nennt `test-container` und seine Gegenprobe.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Docker | zwei Builds ohne Cache der Bau-Stufen liefern dieselben Summen; Image und Smoke bestehen | `make test-container` |
| Gegenprobe | Eingabefehler, abweichende Summe, Cache-Umgehung, Kopf, Image, Smoke, Aufräumen und Aufrufarten des Treibers werden rot | `make test-container-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn die Laufzeit des Ziels die Kette dominiert, wenn das Lastenheft das Image reproduzierbar verlangt, wenn eine Plattform hinzukommt oder wegfällt oder wenn ein zweiter Bau-Host zur Abnahme gehört.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-17 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
