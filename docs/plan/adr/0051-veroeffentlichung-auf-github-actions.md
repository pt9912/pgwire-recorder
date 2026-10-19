# ADR-0051: Die Veröffentlichung läuft auf GitHub Actions bei einem Tag, mit Trockenlauf statt Probe-Tag und ohne Veränderung eines Tags

**Status:** Accepted

**Datum:** 2026-10-19

**Autor:** Architect; angenommen durch Entscheidung des Nutzers am 2026-10-19 (Option A; kein Probe-Tag; Repository-Geheimnisse; Tag-Satz `<Version>` und `latest`).

**Bezug:** [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus), [ADR-0043](0043-binaries-und-container-gate.md), [ADR-0050](0050-ci-auf-github-actions.md)

**Schärft:** [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der erste Release steht an. Die Binaries entstehen reproduzierbar mit `make binaries`, das Image mit `make build`; eine Manifestliste für zwei Plattformen und eine Veröffentlichung gibt es nicht ([ADR-0043](0043-binaries-und-container-gate.md), Option D: der Standard-Treiber kann keine Manifestliste bauen). Die CI läuft auf GitHub Actions und trägt Pins, leere Rechte und ein Gate für Workflow-Dateien ([ADR-0050](0050-ci-auf-github-actions.md)). Der Nutzer hat am 2026-10-18 den ersten Tag `v0.1.0` und das Veröffentlichen erlaubt, Probe-Tags und weitere Tags nicht ohne neues Ja, Zugangsdaten nur als Geheimnisse von GitHub.

Ein Tag und ein Push in eine öffentliche Registry sind nicht zurückzunehmen: ein gesetzter Tag wird nicht verändert (`docs/maintainer/releasing.md` §11), ein veröffentlichtes Image wird über seinen Digest identifiziert. Ein Probe-Tag würde eine Vorabversion dauerhaft in `ghcr.io`, `docker.io` und am GitHub-Release hinterlassen und zählt als weiterer Tag.

## Entscheidung

Wir wählen **Option A: GitHub Actions veröffentlicht bei einem Tag `v*`; dieselbe Datei läuft als Trockenlauf auf einem Zweig und tut alles außer dem, was nach außen wirkt; kein Probe-Tag.**

1. **Träger und Auslöser.** `.github/workflows/release.yml` auf `ubuntu-24.04`, ausgelöst durch den Push eines Tags `v*` (`ci.yml` lässt Tags aus) und durch `workflow_dispatch` als Trockenlauf. Werkzeuge: `docker buildx`, `skopeo` und `gh` des Runner-Images, keine Aktion außer dem Checkout; die Logik liegt in einem Skript, das die Gegenprobe mit Attrappen ohne Netz prüft.
2. **Fail-closed vor der Außenwirkung.** Form des Tags (SemVer), Gleichheit mit der Datei `VERSION`, die Zeile der Änderungshistorie, die Geheimnisse und (beim Lauf) die Lage des Tag-Commits im Standard-Zweig werden geprüft, bevor etwas in eine Registry geschrieben wird. Die erste Außenwirkung ist der Push nach `ghcr.io`.
3. **Ein Image, zwei Registries.** Das Image wird einmal gebaut, als Manifestliste für `linux/amd64` und `linux/arm64`, und genau dieses Layout wird in beide Registries kopiert; der Index-Digest ist in beiden gleich und wird aus den rohen Bytes gelesen. Das Binary im Image muss die Summe aus `SHA256SUMS` tragen. Eine Stelle für den Scan liegt zwischen Bau und Push.
4. **Tags nie überschrieben.** Ein vorhandener Image-Tag mit gleichem Digest wird übersprungen (Wiederanlauf), mit anderem Digest ist der Lauf rot und der Weg ein höherer Tag. `latest` setzt nur ein stabiler Tag, der die höchste stabile Version ist, nach grünem Vergleich und Smoke.
5. **Kein Probe-Tag.** Der Trockenlauf belegt Eingaben, Bau, Prüfung des Images und das Vorhandensein der Geheimnisse; er belegt nicht, dass die Token schreiben dürfen und das Paket öffentlich ist. Das zeigt erst der erste Tag, und ein Fehler dort ist behebbar, ohne einen Tag zu verändern (Wiederanlauf desselben Tags oder ein höherer Tag).
6. **Zugangsdaten.** Repository-Geheimnisse `DOCKERHUB_USERNAME` und `DOCKERHUB_TOKEN`, `GITHUB_TOKEN` für `ghcr.io` und den GitHub-Release; Anmeldung nur über Standardeingabe, nie im Log.

Die Regeln im Einzelnen (Text der Datei, Eingaben, Reihenfolge, Fehlertabelle, Fälle) stehen in [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md).

Entscheidung des Nutzers vom 2026-10-19 (Empfehlung des Architect angenommen): die Annahme dieser ADR; **kein Probe-Tag** (die Pipeline wird vor `v0.1.0` durch den Trockenlauf belegt, sonst bliebe ein `v0.1.0-rc.1` dauerhaft in beiden Registries und am GitHub-Release); Repository-Geheimnisse statt eines Environments `release`; Tag-Satz nur `<Version>` und `latest`, kein `<MAJOR>.<MINOR>`.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — Tag-Lauf mit Trockenlauf, kein Probe-Tag** | der erste Tag ist der erste Lauf mit Außenwirkung; keine Vorabversion bleibt zurück; Fehler nach dem Push sind behebbar | Schreibrecht, Sichtbarkeit des Pakets und das Anlegen des Releases zeigen sich erst am echten Tag |
| B — Probe-Tag `v0.1.0-rc.1` zuerst | belegt die ganze Pipeline, auch das Schreiben und die Sichtbarkeit | ein weiterer Tag (braucht ein neues Ja); die Vorabversion bleibt dauerhaft in beiden Registries und am GitHub-Release (Tags werden nicht gelöscht) |
| C — Probe in einem eigenen Namensraum (anderes Paket, anderes Repository) | die echten Namen bleiben sauber | belegt die Rechte der echten Namen nicht, braucht Konten und Geheimnisse für den Namensraum, zwei Konfigurationen |
| D — Veröffentlichen von Hand vom Rechner des Maintainers | kein Workflow | Zugangsdaten auf dem Rechner, keine Gegenprobe, nicht wiederholbar; widerspricht [ADR-0050](0050-ci-auf-github-actions.md) |
| E — Image aus `dist/` statt Neubau | Gleichheit mit dem Anhang am Release ist trivial | das `Dockerfile` müsste das Binary aus `dist/` aufnehmen, ein Eingriff in die Gate-Stufe; der Neubau ist reproduzierbar und seine Summe wird geprüft |
| F — `ubuntu-latest`, Aktionen von Dritten für Buildx und Login | weniger Skript | Härtung nach [ADR-0050](0050-ci-auf-github-actions.md) verlangt feste Runner und gepinnte Aktionen; die Logik bliebe ohne Gegenprobe |

## Konsequenzen

- Positiv: Der erste Tag löst eine Pipeline aus, deren Logik mit Attrappen je Zusage gezeigt ist; ein Fehler hinterlässt einen benannten Zustand; kein Tag wird verändert.
- Negativ: Der erste Lauf mit Außenwirkung kann an einer Einstellung des Betreibers scheitern (Token, Paketsichtbarkeit, Repository auf Docker Hub). Der Zustand ist dann ein teilweise veröffentlichter Tag in einer Registry; der Weg ist der Wiederanlauf desselben Tags nach Behebung, bei abweichendem Digest ein höherer Tag.
- Akzeptiertes Negativ: Das Image ist nicht als reproduzierbar zugesagt. Grund: `LH-QA-01` verlangt es nur für das Binary; der Wiederanlauf nach einem halben Push scheitert, wenn ein Neubau einen anderen Digest ergäbe, und endet dann mit einem höheren Tag.
- Akzeptiertes Negativ: Ein Tag auf einem Commit mit verändertem `release.yml` führt die veränderte Datei aus. Grund: Verhalten von GitHub; wer taggen darf, bestimmt die Konfiguration des Repositorys.
- Akzeptiertes Negativ: Das Image für `linux/arm64` ist nach Architektur und Summe des Binaries belegt, nicht ausgeführt. Grund: kein Emulator vorausgesetzt (wie bei `make test-container`).
- Folgepflicht: [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md); `docs/maintainer/releasing.md` §3, §6, §8, §9 und §11 beschreiben das ausgeführte Verfahren; `harness/README.md` §Sensors führt `release-gegenprobe`; der Scan vor dem Push und der erste echte Tag kommen mit dem Folge-Slice.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, git, make, Attrappen | Text von `release.yml`, Eingaben, Reihenfolge, Digest-Vergleich, `latest`, Anhänge, Zugangsdaten, Wiederanlauf | `make release-gegenprobe` |
| bash, Docker, `actionlint`, d-check | Rechte, Zeitgrenzen, Syntax und Pins der Datei | `make workflows-check`, `make docs-check` |

Dass ein Lauf auf GitHub grün endet und dass die Token schreiben dürfen, prüft kein Gate; der Beleg ist der Lauf.

## Re-Evaluierungs-Trigger

Wenn der erste Lauf mit einem Tag an einer Stelle scheitert, die ein Probe-Tag oder ein eigener Namensraum gefangen hätte (dann Option B oder C für die nächste Änderung der Pipeline); wenn der Neubau des Images den Wiederanlauf mehrfach verhindert (dann reproduzierbarer Bau oder Image aus `dist/`); wenn die Registries Tags `<MAJOR>.<MINOR>` verlangen; wenn der Nutzer ein Environment für die Geheimnisse will.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-19 | Proposed (Empfehlung des Architect nach dem Auftrag des Nutzers, den Release voranzubringen) | — |
| 2026-10-19 | Accepted: Entscheidung des Nutzers, Option A ohne Probe-Tag, mit Repository-Geheimnissen und Tag-Satz `<Version>` und `latest` | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
