# Review-Report: slice-erster-release-veroeffentlichung — 2026-10-19

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (Fassung A: Tag und Binaries).

**Gegenstand:** `slice-erster-release-veroeffentlichung`, Diff `17a8c618..fc63818c` auf Zweig `rev-rel-a` (Implementer-Commits 44835211 bis fc63818c, dazu Architect-Merge 5f9fddb5 mit SPEC-067 Punkten 3, 9, 13).

**Skill:** `.harness/skills/reviewer.md` @ fc63818c
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-19

> **Zitier-Form.** Dieser Report friert ein; was er zitiert, bewegt sich weiter. Deshalb: Kennung, nicht Adresse — `slice-<Kennung>` statt seines Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill).

**Eingangs-Kontext:**

- Slice-Plan `slice-erster-release-veroeffentlichung` (§1, §3, §4, §6, §7 Belege des Implementers)
- `SPEC-067` Fassung 1a (= A), [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md)
- `LH-FA-19.a`, `LH-QA-01`, `LH-QA-03` (über den Plan-Kopf)
- `AGENTS.md` (Hard Rules, §3.9 bis §3.14)

**Eigene Läufe (engster Sensor, keine Wiederholung der vollen Läufe des Implementers):** `bash -n` für `release.sh` und `release-gegenprobe.sh` grün; `make workflows-check` (2 Dateien, kein Befund) und `make docs-check` (0 Befunde) grün; unveränderte Kopie der Gegenprobe (`GP_MUT=^ok$`) grün; fünf eigene Mutanten von `release.sh` in einer Kopie unter dem Scratch-Verzeichnis gegen die unveränderte Gegenprobe (Ergebnisse in F-1, F-4); Sonde der Äquivalenz 1 von Hand nachgefahren. Kein Push, kein Netz, kein Tag.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Zusage „nur die Meldung `release not found` heißt ‚nicht vorhanden‘“ hat nur die Mutation „jeder Fehler gilt als nicht gefunden“. Ein Mutant, der das Muster auf `not found` (ohne Schreibweise) weitet, färbt keinen der Fälle: `GP_MUT=^ok$` bleibt grün (selbst gemessen). Ein Fehler von `gh release view` wie `HTTP 404: Not Found` würde dann als „Release fehlt“ gelesen und führte zu `release create`. Der Kopf von `release.sh` sagt „je Zusage eine Mutation dieser Datei“. | `AGENTS.md` §3.10 (Zusage ohne Mutation, Aufzählung/Eingrenzung je Element), §3.11 | `tools/harness/release.sh` · `grep -q 'release not found' "$scratch/gh.err"`; `tools/harness/release-gegenprobe.sh` · `f_p9_view_anderer_fehler` | ja — `make release-gegenprobe` mit dem Mutanten bliebe grün | Zusage mit Eingrenzung ohne Grenzmutation |
| F-2 | MEDIUM | Die Zusage „ein Tag wird nie durch einen späteren abgebrochen“ (`concurrency: cancel-in-progress: false`) trägt nach der dokumentierten Semantik von GitHub nur für den laufenden Lauf: Ein wartender Lauf derselben Gruppe wird durch einen neueren wartenden ersetzt. Der Trockenlauf teilt die Gruppe `release`; ein von Hand gestarteter Trockenlauf, während `v0.1.0` läuft und ein zweiter Tag wartet, kann den wartenden Tag-Lauf verdrängen. Kein Fall und keine Mutation kann das fassen (`p1-*` hält nur den Text). | `SPEC-067` Punkt 13 (Gleichzeitigkeit); `AGENTS.md` §3.11 | `.github/workflows/release.yml` · `cancel-in-progress: false`; `SPEC-067-release.md` · „ein Tag wird nie durch einen späteren abgebrochen“; Plan §1 · „ein Tag wird nie durch einen späteren abgebrochen“ | nein — nur ein Lauf auf GitHub zeigt es; Wortlaut des Vertrags ist Sache des Architect | Zusage über das Verhalten einer Plattform ohne Prüfung |
| F-3 | LOW | `docs/maintainer/releasing.md` §3 definiert den Probe-Tag weiter als Tag, „mit dem das Release-Verfahren geprüft wird“, während der neue §6 sagt „Ein Probe-Tag gibt es nicht“ (Entscheidung des Nutzers vom 2026-10-19). Der Leser findet beide Aussagen im selben Dokument, ohne dass §3 sagt, dass der Begriff nur noch für das Homebrew-Verfahren gilt. | Maintainability; [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md) (kein Probe-Tag) | `docs/maintainer/releasing.md` · „Ein **Probe-Tag** ist ein Vorabversions-Tag“ | nein | Dokument widerspricht sich nach Teilanpassung |
| F-4 | LOW | Die Aufzählung der Anhänge (sechs Binaries und `SHA256SUMS`, Punkt 9) hat in der Mutationstabelle nur „ein achter Anhang“; es gibt keine Mutation je Element, die nur ihn bricht. Der Mutant „`windows`/`arm64` fehlt in `erwartete_namen`“ färbt gemessen 20 Fälle statt eines Falls der Aufzählung. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `tools/harness/release.sh` · `erwartete_namen() {`; Plan §7 · „9: Anhänge genau die sieben Dateien“ | ja — Mutant je Element, Lauf der Gegenprobe | Aufzählung ohne Mutation je Element |
| F-5 | LOW | Die Zelle der Änderungshistorie wird mit `awk -F'|'` gelesen; enthält der Text der dritten Zelle ein `|`, wird er am ersten `|` abgeschnitten, und der Release trägt einen verkürzten Text ohne Meldung. §6 des Plans und `SPEC-067` Punkt 9 nennen die Randform (Zelle mit `|`) nicht. | `AGENTS.md` §3.12 (Randform nicht genannt, im Code entschieden) | `tools/harness/release.sh` · `awk -F'\|' -v v="$1"` in `notiz_der_version` | ja — Fall mit `|` in der Zelle | Randform im Code still entschieden |
| F-6 | INFO | Größe: Metrik nach `AGENTS.md` §3.14 laut Plan §7 1393 (Marke 1300, Meldeschwelle 1000, Schätzung 820 bis 1070); die datierte Zeile „weiter“ oder „schneiden“ in §7 steht aus. Die Entscheidung liegt nicht beim Review. Die Zahl nicht nachgemessen; Dateigrößen im Diff (Gegenprobe 1016, `release.sh` 230) stimmen mit dem Beleg überein. | Plan §4 (Rückführung) | Plan · „Entscheidung „weiter“ oder „schneiden“ steht aus“ | ja — `git diff --numstat` über die Pfade des Plans | — |
| F-7 | INFO | DoD-Belege, die erst der Verifier bestätigt: Trockenlauf auf GitHub nicht gefahren (Lauf-Kennung und Dauer fehlen in §7); der Wortlaut `release not found` von `gh` ist im Dev-Container nicht gemessen (akzeptierte Grenze in `SPEC-067` Punkt 9). Rolle: Verifier. | `SPEC-067` Punkt 9, 12 | Plan §7 · „Trockenlauf auf GitHub:** nicht gefahren“ | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.github/workflows/release.yml` (Text gegen Punkt 1a, Pin, `permissions: {}` oben, Job nur `contents: write`, Ereignisse nur `push` auf `v*` und `workflow_dispatch`, kein `pull_request*`, kein `environment:`, Geheimnis nur `GITHUB_TOKEN` im letzten Schritt, Ref und SHA über `env:` statt in `run:`, `persist-credentials: false`, Außenwirkung nur bei `push`) | geprüft, ohne Befund (außer F-2 zur Gleichzeitigkeit) |
| `tools/harness/release.sh` (Reihenfolge der Gründe (1) bis (6), Regeln der SemVer-Form, `VERSION`-Form, Marken zwischen den Schritten, Aufräumen bei Ende/`SIGINT`/`SIGTERM`, keine Ausgabe von Geheimnissen, kein `set -x`/`env`, `--repo` an jedem `gh`-Aufruf, `GH_TOKEN` vor erstem `gh`, Trockenlauf ohne `gh`) | geprüft, ohne weiteren Befund |
| Äquivalente Mutanten in Plan §7 (`m-p3-form-v-optional`, `m-p3-lauf-fehlende-zeile-ohne-eigene-pruefung`) | geprüft: Invariante, Ort und Sonde tragen; die Sonde 1 selbst nachgefahren (Ref `refs/tags/1.2.3` endet mit `release: ereignis:`, Exit 2), `-01` abgelehnt, `-0a` zugelassen; kein Befund |
| `tools/harness/release-gegenprobe.sh` (Stichproben: Fälle `p4-*`, `p9-*`, Mutanten SHA256SUMS-Bedingung und `--is-ancestor` gegen `HEAD` jeweils rot) | geprüft, Lücken siehe F-1, F-4 |
| `harness/mk/release.mk`, `harness/README.md` (Zeile von `make release-gegenprobe`), `VERSION` (`0.1.0-dev`, LF) | geprüft, ohne Befund |
| `docs/maintainer/releasing.md` (Fehlertabelle gegen `SPEC-067` Punkt 11, Zeilen A) | geprüft; Probe-Tag-Widerspruch siehe F-3 |
| `docs/user/benutzerhandbuch.md`, `README.md` (unverändert; Ist-Zustand, §3.11) | geprüft, ohne Befund |
| Plan folgt dem Diff (§3.9: §1, §3, §6, Kopf), Commit-Regeln (§3.3: kein Move im Diff; Kennung in den Commit-Messages), Adressen (§3.13: keine neue Sendung im Diff) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage mit Eingrenzung ohne Grenzmutation · Zusage über das Verhalten einer Plattform ohne Prüfung · Dokument widerspricht sich nach Teilanpassung · Aufzählung ohne Mutation je Element · Randform im Code still entschieden

## Verdikt

**Merge-blockierend:** ja — zwei MEDIUM. F-1 liegt beim Implementer. F-2 betrifft den Vertragstext (Architect): entweder die Zusage enger fassen oder eine Prüfung benennen; der Implementer entscheidet das nicht. Ein Finding wird nicht herabgestuft, weil es widersprochen wird; bei Widerspruch zu F-2 läuft der Konflikt-Pfad über den Architect.

**Übergabe:** F-1, F-3, F-4, F-5 an den Implementer; F-2 an den Architect; F-6 an die Koordination (Größenentscheidung); F-7 an den Verifier. Die Finding-Klassen gehen in die Slice-Closure §7. Der Report ersetzt keine Verifikation gegen die DoD.
