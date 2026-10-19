# ADR-0052: Das Image wird vor dem Push mit Trivy gescannt; behebbare Befunde ab HIGH sperren, Ausnahmen gibt es nicht

**Status:** Proposed

**Datum:** 2026-10-19

**Autor:** Architect; Empfehlung, bis der Nutzer sie annimmt.

**Bezug:** [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [ADR-0051](0051-veroeffentlichung-auf-github-actions.md), [ADR-0050](0050-ci-auf-github-actions.md)

**Schärft:** [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Image `ghcr.io/pt9912/pgwire-recorder` und `docker.io/pt9912/pgwire-recorder` ist öffentlich und nicht zurückzunehmen: ein veröffentlichter Tag wird nie verändert ([ADR-0051](0051-veroeffentlichung-auf-github-actions.md)). Ein bekannter, behebbarer Fehler im Image (im Basis-Image `distroless/static` oder in einem Go-Modul des Binaries) wäre mit dem ersten Tag veröffentlicht und nur durch einen höheren Tag zu ersetzen. Die Pipeline lässt dafür einen Platz zwischen Bau und Push frei. Die Frage ist, mit welchem Werkzeug, ab welcher Schwere und mit welchen Ausnahmen der Schritt sperrt; eine Schwelle ist nach `AGENTS.md` §3.6 eine Entscheidung mit Gründen.

## Entscheidung

Wir empfehlen **Option A: Trivy als per Digest gepinntes Image scannt das OCI-Layout jeder Plattform vor dem Push; ein Befund mit Schweregrad `CRITICAL` oder `HIGH` und verfügbarer Behebung beendet den Lauf; Ausnahmen gibt es nicht.**

1. **Werkzeug.** Trivy (`aquasec/trivy`), Pin per Digest mit Tag-Kommentar, wie der Builder des Images. Es ist dasselbe Werkzeug, das d-check für seine veröffentlichten Images einsetzt (dort gemessen: ein Go-Binary auf `distroless/static` ist der Fund-Raum, ein grüner Lauf sagt „nichts Bekanntes in diesem Raum“).
2. **Schwelle.** `CRITICAL` und `HIGH`, nur mit verfügbarer Behebung. Ein Befund ohne Behebung wird gezählt und ausgegeben, sperrt aber nicht: gegen ihn gibt es mit dem gepinnten Basis-Image keinen Weg, und eine Sperre ohne Weg wäre keine Schwelle, sondern ein Stillstand.
3. **Keine Ausnahmen.** Kein Ignorierfile, keine Liste, kein Schalter. Ein Befund mit Behebung wird behoben (Basis-Image oder Modul heben); eine andere Schwelle ist eine neue Entscheidung (ADR). Das entspricht dem Suppression-Verbot für den Code (`AGENTS.md` §3.2): eine Ausnahme, die nur das Werkzeug kennt, sieht kein Leser.
4. **Fail-closed.** Ein Scan, der nicht lief (Werkzeug scheitert, leere Prüfmenge, Datenbank älter als 72 Stunden oder nicht lesbar, falsche Plattform gescannt), ist rot, nie grün und nie übersprungen.
5. **Stelle.** Zwischen `Image` und `Veröffentlichen`, auch im Trockenlauf; er liest das Layout, das der Push veröffentlicht, sodass das Gescannte das Veröffentlichte ist.

Die Regeln im Einzelnen (Aufruf, Auswertung, Exit-Codes, Fälle, Grenzen) stehen in [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md) Punkt 16.

Zur Entscheidung des Nutzers vorgelegt (Empfehlung des Architect, keine Annahme): die Annahme dieser ADR; Trivy statt eines anderen Scanners; die Schwelle `HIGH` mit Behebung; der Verzicht auf Ausnahmen.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — Trivy, `CRITICAL`/`HIGH` mit Behebung, keine Ausnahmen** | ein Werkzeug, das im Schwesterprojekt d-check gemessen ist; sperrt, was ein Mensch beheben kann; keine stille Ausblendung | ein Fehlalarm oder eine Behebung, die nicht erreichbar ist, hält den Tag auf, bis die Schwelle per ADR geändert ist |
| B — Grype | gleicher Zweck, anderer Datenbestand | nicht im Repo oder Schwesterprojekt gemessen; ein zweiter Scanner im Wartungsumfang |
| C — Schwelle nur `CRITICAL` | seltener rot | ein behebbares `HIGH` wird veröffentlicht und nur durch einen höheren Tag ersetzt |
| D — alle Befunde, auch ohne Behebung | strengste Aussage | Stillstand ohne Weg; der Tag wäre von einer fremden Datenbank abhängig |
| E — Ausnahmeliste mit Begründung und Ablauf | ein Fehlalarm blockiert nicht | eine Liste ist ein Ort für Dauer-Ausnahmen; das Suppression-Verbot gilt dem Sinn nach; mehr Vertrag, mehr Fälle |
| F — Scan nach dem Push | kein Eingriff in die Pipeline | der Fehler ist veröffentlicht, bevor der Scan ihn meldet; das Gegenteil des Zwecks |
| G — kein Scan | nichts zu bauen | ein bekannter, behebbarer Fehler im ersten Release |

## Konsequenzen

- Positiv: Was der Scan an behebbarem Bekanntem sieht, erreicht keine Registry; das Gescannte ist das Veröffentlichte.
- Negativ: Der Ausgang hängt vom Stand der Datenbank am Tag des Laufs; derselbe Quellstand kann an zwei Tagen verschieden bewertet werden. Der Lauf des ersten Tags kann am Scan scheitern, ohne dass das Repo sich geändert hat; dann hilft ein Wiederanlauf nach Heben der Ursache oder ein höherer Tag.
- Akzeptiertes Negativ: Ein grüner Scan sagt „nichts Behebbares mit `HIGH` oder `CRITICAL` bekannt“, nicht „sicher“. Grund: Das ist die Reichweite jedes Scanners dieser Art.
- Akzeptiertes Negativ: Befunde ohne Behebung sperren nicht. Grund: kein Weg, sie mit dem Basis-Image zu beheben.
- Folgepflicht: [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md) Punkt 16; `docs/maintainer/releasing.md` §5 und §10 nennen den Scan als Beleg der Checkliste; die Pflege des Pins (Version und Digest von Trivy) ist eine Handlung wie die der übrigen Pins.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, Attrappen | Aufruf, Auswertung, Fail-closed, Marke und Reihenfolge des Schritts `scan` | `make release-gegenprobe` |

Ob der Scan am Tag des Laufs grün ist, prüft kein Gate; der Beleg ist der Lauf.

## Re-Evaluierungs-Trigger

Wenn ein behebbarer Befund unter der Schwelle (`MEDIUM`) veröffentlicht wurde, der ein Problem war; wenn Fehlalarme den Tag wiederholt aufhalten (dann Option E per neuer ADR); wenn der Scan häufiger am Netz oder an der Datenbank als an Befunden scheitert (dann ein gepinnter Datenbankstand mit Alter, neue ADR); wenn d-check den Scanner wechselt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-19 | Proposed (Empfehlung des Architect vor dem Code von Slice F) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
