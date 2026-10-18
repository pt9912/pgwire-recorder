# ADR-0047: Die Gegenprobe des Container-Starts ist ein Ziel für den Host und kein Gate

**Status:** Proposed

**Datum:** 2026-10-18

**Autor:** Architect. Die Sache hat der Nutzer am 2026-10-18 entschieden („Wir können Gegenproben aufteilen. Einen speziellen, der nur im Host (u-boot und devcontainer) läuft“); die Annahme dieser Fassung steht aus. Die Einzelheiten sind bis dahin Empfehlung.

**Bezug:** [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [ADR-0045](0045-container-start-als-werkzeug-im-repo.md), [ADR-0046](0046-gegenprobe-ueberspringt-im-dev-container.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md)

**Schärft:** [`SPEC-058`](../../../harness/vertraege/SPEC-058-container-start.md)

**Supersedes:** [ADR-0046](0046-gegenprobe-ueberspringt-im-dev-container.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0045](0045-container-start-als-werkzeug-im-repo.md) macht die Gegenprobe des Container-Starts zu einem Gate der Kette und nimmt an, dass `make gates` auf einem Host ohne `u-boot` und `devcontainer` rot ist. [ADR-0046](0046-gegenprobe-ueberspringt-im-dev-container.md) löst das für den Dev-Container durch Erkennung und Überspringen. Beides hängt an einer Erkennung, die irren kann, und lässt Fälle im Container ungeprüft. Der Nutzer hat entschieden, die Gegenproben aufzuteilen: Die Gegenproben, die ohne diese Werkzeuge laufen, bilden ein Sammelziel für den Dev-Container; die des Container-Starts, die sie braucht, ein Ziel nur für den Host. Alle Gegenproben verlassen dabei `GATE_CHECKS` (eigener Vorgang, dessen Entscheidung hier nicht liegt).

## Entscheidung

Empfehlung: `container-start-gegenprobe` ist ein Ziel für den Host. Sie hängt nicht an `GATE_CHECKS`, steht nicht in der Ausnahmeliste der schnellen Kette und gehört nicht zum Sammelziel für den Dev-Container; ihr Sammelziel für den Host heißt `gegenproben-host`, das für den Dev-Container `gegenproben` (Namen sind Vorschlag). Fehlt am Host ein Werkzeug, ist sie rot und nennt es; Erkennung und Überspringen im Dev-Container entfallen. [ADR-0046](0046-gegenprobe-ueberspringt-im-dev-container.md) wird ersetzt. Von [ADR-0045](0045-container-start-als-werkzeug-im-repo.md) bleibt die Entscheidung (Ort im Repo, Host-Voraussetzungen); überholt ist nur die Aussage, die Gegenprobe sei ein Gate der Kette. Die Einzelregeln stehen in [`SPEC-058`](../../../harness/vertraege/SPEC-058-container-start.md) Punkte 15 und 16.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — Host-Ziel, kein Gate** | keine Erkennung, die irren kann; `make gates` im Dev-Container braucht kein `u-boot`; die Fälle laufen am Host vollständig | kein Hook und keine Kette erzwingt den Lauf; er hängt am Bericht des Slice |
| B — Überspringen im Dev-Container ([ADR-0046](0046-gegenprobe-ueberspringt-im-dev-container.md)) | die Kette bleibt eine | Erkennung mit zwei Marken, Fälle im Container ungeprüft |
| C — Werkzeuge in das Gerüst des Containers bauen | überall dieselben Fälle | ein Umbau des Gerüsts außerhalb dieses Vorhabens |

## Konsequenzen

- Positiv: `make gates` ist im Dev-Container ohne `u-boot` und `devcontainer` grün; die Gegenprobe prüft ihre Voraussetzung am Host, wo sie da ist.
- Akzeptiertes Negativ: Der Lauf der Gegenprobe ist nicht erzwungen; ein Slice, der Teil I oder Teil II ändert, belegt ihn im Bericht, und das Review prüft den Beleg. Grund: Die Entscheidung des Nutzers nimmt sie aus der Kette.
- Folgepflicht: `harness/README.md` führt die Gegenprobe nicht mehr unter *Sensors*, sondern unter *Werkzeuge* (kein Gate, Host); [ADR-0045](0045-container-start-als-werkzeug-im-repo.md) bleibt Bindung. Mit der Annahme dieser ADR trägt [ADR-0046](0046-gegenprobe-ueberspringt-im-dev-container.md) im Status und in der Geschichte den Vermerk der Ersetzung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, Attrappen, Docker | Eingaben, Voraussetzungen, Zustand, Gerüst und (Teil II) Container, Bereitschaft, Abbruch werden an Fällen und Mutationen geprüft; am Host | `make container-start-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn `u-boot` und `devcontainer` in jedem Dev-Container liegen, wenn der Lauf am Host nicht mehr belegt wird oder wenn die Aufteilung der Gegenproben zurückgenommen wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-18 | Proposed (Empfehlung des Architect nach Entscheidung des Nutzers in der Sache) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
