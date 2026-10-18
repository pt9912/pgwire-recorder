# ADR-0046: Die Gegenprobe des Container-Starts überspringt im Dev-Container, was dort fehlt

**Status:** Accepted

**Datum:** 2026-10-18

**Autor:** Architect. Angenommen durch Entscheidung des Nutzers am 2026-10-18 (Dialogantwort „Im Container überspringen“, danach „Annehmen“ zu dieser Fassung).

**Bezug:** [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [ADR-0045](0045-container-start-als-werkzeug-im-repo.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md)

**Schärft:** [`SPEC-058`](../../../harness/vertraege/SPEC-058-container-start.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0045](0045-container-start-als-werkzeug-im-repo.md) nimmt als Preis an, dass `make gates` auf einem Host ohne `u-boot` und `devcontainer` rot ist. Der Dev-Container, den dieses Werkzeug erzeugt, führt beide nicht; ein frisch gebauter Container wäre für `make gates` rot, obwohl dort gearbeitet wird. Eine Lockerung eines Gates ist eine Entscheidung (`AGENTS.md` §3.6), und [ADR-0045](0045-container-start-als-werkzeug-im-repo.md) ist angenommen und wird nicht überschrieben.

## Entscheidung

Wir wählen: Die Gegenprobe `container-start-gegenprobe` überspringt im Dev-Container die Fälle, die `u-boot` oder `devcontainer` echt brauchen, wenn eines davon fehlt, meldet es laut und endet mit Exit 0, wenn alles Übrige grün ist. Am Host bleibt sie rot bei jedem fehlenden Werkzeug. Berechtigt zum Überspringen sind nur diese beiden Werkzeuge; jede andere Voraussetzung bleibt Pflicht. Die Erkennung des Dev-Containers, die Fälle und Mutationen stehen in [`SPEC-058`](../../../harness/vertraege/SPEC-058-container-start.md) Punkt 15. [ADR-0045](0045-container-start-als-werkzeug-im-repo.md) bleibt in Kraft; diese ADR ergänzt seine Konsequenz „`make gates` ist dort rot“ für den Dev-Container.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — im Dev-Container überspringen, am Host rot** | `make gates` läuft im Container, wo der Nutzer arbeitet; am Host bleibt der Beleg Pflicht | im Container sind die übersprungenen Fälle ungeprüft; die Erkennung kann irren |
| B — überall rot lassen ([ADR-0045](0045-container-start-als-werkzeug-im-repo.md)) | kein Gate wird gelockert | jeder frische Container braucht eine Handarbeit (`u-boot` hineinkopieren) |
| C — `u-boot` und `devcontainer` in das Image des Containers bauen | überall dieselben Fälle | der Container wäre nicht mehr das Gerüst, das `u-boot` erzeugt; ein Umbau des Gerüsts außerhalb dieses Vorhabens |

## Konsequenzen

- Positiv: `make gates` ist im Dev-Container grün, ohne dass ein Werkzeug von Hand hineinkommt.
- Akzeptiertes Negativ: Die übersprungenen Fälle (der Fall gegen das echte `u-boot`, in Teil II die Fälle gegen `devcontainer`) laufen im Container nicht; sie laufen am Host. Grund: Der Container führt die Werkzeuge nicht, und sie hineinzubringen verlangte Option C.
- Akzeptiertes Negativ: Ein Host in einem Container mit beiden Marken würde überspringen. Grund: Beide Marken legt nur ein Gerüst von `u-boot` an.
- Folgepflicht: `harness/README.md` §Sensors nennt das Überspringen bei der Gegenprobe mit Bindung an diese ADR.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, Attrappen | Erkennung an Attrappen, Überspringen nur für `u-boot` und `devcontainer`, Host rot | `make container-start-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn `u-boot` und `devcontainer` im Gerüst des Containers liegen, wenn die Erkennung am Host falsch positiv wird oder wenn ein weiteres Werkzeug im Dev-Container fehlt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-18 | Proposed (Empfehlung des Architect nach Entscheidung des Nutzers in der Sache) | — |
| 2026-10-18 | Accepted: Entscheidung des Nutzers, Fassung mit der Erkennung durch `/.dockerenv` und `u-boot-docker-init` | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
