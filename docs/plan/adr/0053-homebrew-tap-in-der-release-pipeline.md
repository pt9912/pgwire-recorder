# ADR-0053: Die Formel des Homebrew-Taps entsteht in einem eigenen Job der Release-Pipeline nach dem GitHub-Release, mit eigenem Token und ohne Probe-Tag

**Status:** Accepted

**Datum:** 2026-10-20

**Autor:** Architect; die Annahme ist die Entscheidung des Nutzers und steht aus. Vom Nutzer ist nur der Name des Geheimnisses `HOMEBREW_TAP_GITHUB_TOKEN` gesetzt (Mitteilung vom 2026-10-20), Art und Rechte des Tokens hat er nicht genannt.

**Bezug:** [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [ADR-0051](0051-veroeffentlichung-auf-github-actions.md), [ADR-0052](0052-scan-des-images-vor-dem-push.md), [ADR-0050](0050-ci-auf-github-actions.md)

**Schärft:** [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`v0.1.0` ist veröffentlicht, ohne Formel. Der Tap `pt9912/homebrew-pgwire-recorder` ist angelegt und leer. Die Formel soll bei jedem stabilen Release entstehen (`LH-FA-19`) und das Binary des Release-Tags mit der Summe aus `SHA256SUMS` installieren. Die Pipeline ([ADR-0051](0051-veroeffentlichung-auf-github-actions.md)) hat `permissions: {}`, gepinnte Aktionen, Gegenproben ohne Netz und genau sieben Anhänge am Release. Ein Wiederanlauf des Jobs `release` baut das Image neu und kann an einem anderen Digest scheitern (ADR-0051, akzeptiertes Negativ). Das Schwesterprojekt u-boot macht dasselbe mit einem Skeleton, einem Fill-Skript, der Formel als Release-Asset und einem Job `tap` (dort ADR-0016); der Nutzer nennt es als Vorbild.

Ein Probe-Tag ist hier ausgeschlossen: Er wäre eine Vorabversion, und eine Vorabversion ändert den Tap nicht, kann also die Formel gar nicht erzeugen; er wäre außerdem ein weiterer Tag mit Veröffentlichung in beiden Registries (ADR-0051).

## Entscheidung

Wir wählen **Option A: ein zweiter Job `homebrew` in `release.yml`, der nach dem Job `release` die Formel aus der Vorlage im Repository und der Datei `SHA256SUMS` des GitHub-Release erzeugt und in den Tap pusht; Probe ist ein lokales Bare-Repository, nicht ein Tag.**

1. **Stelle.** Eigener Job mit `needs: release`, nur bei einem Push eines Tags. Ein Fehler im Tap lässt Release und Registries unberührt, und „Re-run failed jobs“ wiederholt nur diesen Job, ohne das Image neu zu bauen.
2. **Quelle.** Die Summen kommen aus `SHA256SUMS` am veröffentlichten Release (das, was der Anwender lädt); die Formel entsteht aus einer Vorlage im Repository. Sie reist nicht als achter Anhang: SPEC-067 Punkt 9 verlangt genau sieben.
3. **Umfang.** Nur stabile Tags; eine Vorabversion ändert nichts. Der Tap trägt eine Formel für die höchste stabile Version; eine ältere Version setzt ihn nicht zurück, eine gleiche mit anderem Text (korrigierte Vorlage) schreibt ihn fort.
4. **Zugang.** Das Repository-Geheimnis `HOMEBREW_TAP_GITHUB_TOKEN` steht nur im Schritt des Jobs `homebrew` und nicht im Job `release`; `git` bekommt ihn über `GIT_ASKPASS`, nie als Argument, in der URL oder im Log. **Empfehlung:** ein Fine-grained Personal Access Token nur für das Repository `pt9912/homebrew-pgwire-recorder` mit dem Recht *Contents: Read and write* und Ablaufdatum, sonst nichts; die Bestätigung, dass der gesetzte Token dazu passt, ist die Entscheidung des Nutzers.
5. **Fehlendes Geheimnis ist rot.** Anders als bei u-boot (Warnung und Überspringen) endet der Job mit Fehler: Er hängt hinter dem Release und gefährdet nichts, und ein Überspringen ließe den Tap bei einem stabilen Tag unbemerkt zurück.
6. **Nachtrag.** Derselbe Schritt lässt sich von Hand aus einem Klon des Hauptzweigs mit den zwei Token in der Umgebung fahren; so bekommt `v0.1.0`, dessen Lauf den Job nicht hatte, seine Formel, und ein Fehler in Skript oder Vorlage braucht keinen neuen Tag.

Die Regeln im Einzelnen (Text von `release.yml` und der Formel, Reihenfolge der Gründe, Verhalten des Taps, Fälle) stehen in [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md) Punkt 1c und 19.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — eigener Job, Summen aus dem Release, Vorlage im Repository** | Tap-Fehler ohne Image-Neubau wiederholbar; der Token ist vom Job mit den Registry-Rechten getrennt; sieben Anhänge bleiben | ein zweiter Job im Text der Datei (mehr Zeilen für die Gegenprobe) |
| B — Schritt im Job `release` | kein zweiter Job | ein roter Tap verlangt Wiederanlauf des Jobs mit Neubau des Images (anderer Digest möglich: höherer Tag für einen Fehler, der nichts Veröffentlichtes betrifft); der Token läge neben den Registry-Rechten |
| C — u-boot-Muster: Formel als Release-Asset, Job lädt sie | die Formel ist am Release prüfbar | ein achter Anhang ändert SPEC-067 Punkt 9 und seine Fälle; zwei Orte für dieselbe Formel |
| D — eigener Workflow auf `release: published` oder `workflow_run` | getrennte Datei | ein Release, den `GITHUB_TOKEN` anlegt, löst keinen weiteren Workflow aus; `workflow_run` ist in der CI-Härtung ausgeschlossen ([ADR-0050](0050-ci-auf-github-actions.md)) |
| E — nur von Hand | nichts zu bauen | die Formel entsteht nicht bei jedem Release (`LH-FA-19`) |
| F — Deploy-Key oder GitHub App statt Token | Deploy-Key ohne Ablauf und an ein Repository gebunden | der Nutzer hat ein Token-Geheimnis gesetzt; Deploy-Key braucht `ssh` und gepinnte Host-Schlüssel, App eine Aktion von Dritten (ADR-0051: keine Aktion außer dem Checkout) |
| G — Probe-Tag mit echtem Tap | belegt mehr | eine Vorabversion kann den Tap nicht ändern (Widerspruch); weiterer Tag mit Veröffentlichung; ein Tag braucht ein neues Ja |

## Konsequenzen

- Positiv: Die Formel entsteht bei jedem stabilen Tag; ein Fehler im Tap kostet keinen Tag; der Push-Token liegt getrennt von den Registry-Rechten.
- Negativ: Das Geheimnis fehlt oder trägt zu wenig Rechte erst im Lauf eines stabilen Tags auf; der Trockenlauf führt den Job nicht.
- Akzeptiertes Negativ: Die Gegenprobe belegt den Push nur gegen ein lokales Probe-Tap, nicht die Rechte des Tokens, nicht die Annahme der Formel durch Homebrew und nicht das Binary auf macOS. Grund: kein Netz und kein Geheimnis in der Gegenprobe; das belegt erst der Nachweis mit dem echten Tap (Abnahmeszenario 11).
- Akzeptiertes Negativ: Das Skript ist nur unter Linux zugesagt. Grund: Pipeline und Gegenprobe laufen dort.
- Folgepflicht: [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md); `docs/maintainer/releasing.md` §3, §4, §6, §7, §11; `harness/README.md` §Sensors (`release-gegenprobe`); Lastenheft unverändert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, git, Attrappen | Text von `release.yml` mit dem Job, Reihenfolge der Gründe, Formel, Verhalten des Probe-Taps, Token | `make release-gegenprobe` |
| bash, Docker, `actionlint`, d-check | Rechte, Zeitgrenzen, Syntax und Pins | `make workflows-check`, `make docs-check` |

## Re-Evaluierungs-Trigger

Wenn der Push in den echten Tap an einer Stelle scheitert, die ein Probe-Tap nicht zeigt und nicht der Token ist; wenn der Nutzer einen anderen Zugang (Deploy-Key, App) will; wenn der Tap Formeln für mehr als ein Werkzeug trägt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-20 | Proposed (Empfehlung des Architect) | — |
| 2026-10-20 | Accepted (Entscheidung des Nutzers) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
