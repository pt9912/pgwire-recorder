# Review-Report: slice-harness-ci — 2026-10-19

**Review-Art:** Code — Code-Review gegen Plan (`slice-harness-ci`), `SPEC-063`, `SPEC-064`, [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md) und `AGENTS.md` §3.9 bis §3.14 (Maintainability; keine DoD-Verifikation).

**Gegenstand:** `slice-harness-ci`, Diff `39f1ce8a..5f88de70` (Zweig `impl-ci`)

**Skill:** `.harness/skills/reviewer.md` @ 5f88de70 (Stand des Worktrees)
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-19

> **Zitier-Form.** Kennung, nicht Adresse: `slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-harness-ci`, §7 (Belege des Implementers, Mutationstabellen)
- `SPEC-064` (Teile 1 und 2, Punkte 14 und 18), `SPEC-063`
- [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md), gelesen: [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md), [ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md)
- keine `LH-*`-Kennung (der Slice ändert die Prüfung der Integration, nicht das Binary)
- `AGENTS.md` (Hard Rules, §3.9 bis §3.14)

Gelaufen: `bash -n` auf beide Skripte (grün); `make docs-check` (1196 Dateien, 0 Befunde); Teilläufe von `tools/harness/ci-gegenprobe.sh` mit `GP_MUT` (16 Läufe grün: `m-p9-*`, `m-p7-leer*`, `m-p1-pin*`, `m-p1-zeile-25*`, `rekursion-schutz`, `fall-ohne-funktion`); eigene Mutanten in Kopien (Sonden zu F-1 bis F-3). Metrik `git diff --numstat 39f1ce8a..5f88de70 -- .github tools harness .d-check.yml Makefile`: Summe 747 (wie gemeldet; unter der Marke 900, über der Meldeschwelle 700). Nicht wiederholt: volles `make gates`, volle `make ci-gegenprobe`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Zusage „Exit 2 auch bei einem `git`-Fehler“ (Punkt 8) hat nur den Fall `p8-git-fehler` (kein Repository, trifft die Prüfung auf den flachen Klon). Das `\|\| fehler` am `git rev-list --count` hat weder Fall noch Mutation: Der Mutant `\|\| true` dort bleibt grün; bei auflösbarer Basis und auflösbarem Kopf, aber fehlendem Zwischen-Commit (Objekt gelöscht), gibt er `[: : integer expression expected` und dennoch Exit 0 mit der Range aus, das Original Exit 2 (gemessen). §7 nennt „Kein Mutant blieb grün“ und für `m-p8-git-fehler-ignoriert` die Invariante „jeder Fehler endet mit Exit 2 und Präfix“, die für diesen Aufruf nicht gilt. | `AGENTS.md` §3.10 (Zusage je Mutation) · `SPEC-064` Punkt 8 | `tools/harness/ci-range.sh` · `anzahl="$(git rev-list --count "$basis..$kopf"` | ja — Mutant `\|\| true` an dieser Stelle, `make ci-gegenprobe` bleibt grün | Zusage mit Ortsangabe ohne Mutation (git-Fehler nur an einer Stelle gefangen) |
| F-2 | MEDIUM | Punkt 18 sagt „Bei Ende, bei `SIGINT` und bei `SIGTERM`“ zu; Fälle und Mutationen gibt es nur für die beiden Signale. Ohne `trap raeume_auf EXIT` bleibt die Gegenprobe grün (gemessen, `GP_MUT='^ok$'`), und das Temp-Verzeichnis `tmp.*` liegt nach dem regulären Ende im `TMPDIR`. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) · `SPEC-064` Punkt 18 | `tools/harness/ci-gegenprobe.sh` · `trap raeume_auf EXIT` | ja — Mutant ohne die Zeile, danach `ls $TMPDIR` | Element einer Aufzählung ohne Fall und Mutation (Abbruch: „Ende“) |
| F-3 | MEDIUM | Punkt 6 liest `CI_DEFAULT_BRANCH` „bei `push` mit neuem Zweig oder unbekanntem Vorgänger“. Die Fälle `p6-default-*` laufen nur mit Nullen; der Zweig „vorgänger nicht im Klon“ hat keinen Fall mit ungültigem Standard-Zweig. Der Mutant ohne `pruefe_standard` im `elif`-Zweig (gemessen: Exit 1 durch `set -u` statt Exit 2) lässt jeden vorhandenen Fall grün (`p7-push-vorgaenger-fehlt` setzt einen gültigen Wert). | `AGENTS.md` §3.10 (Ortsangabe/Aufzählung je Element) · `SPEC-064` Punkt 6 | `tools/harness/ci-range.sh` · `pruefe_standard` im Zweig `elif ! basis="$(aufloesen "$basis_roh")"` | ja — Mutant ohne den Aufruf; Probe: Push, Vorgänger nicht im Klon, `CI_DEFAULT_BRANCH=a..b` | Element einer Aufzählung ohne Fall (zweiter Anlass für `CI_DEFAULT_BRANCH`) |
| F-4 | LOW | `SPEC-064` Punkt 18 und Randform 18 im Plan sagen, bei `ci-gegenprobe` entfalle die Mutation „Unterprozesse nicht beendet“ (keine Unterprozesse). Das Skript führt sie (`m-p18-int-kinder-bleiben`, `m-p18-term-kinder-bleiben`; 6 Mutationen, §7 zählt 6); Skript, Vertrag und Plan sagen Verschiedenes. | `AGENTS.md` §3.9 (Plan folgt dem Diff) | `harness/vertraege/SPEC-064-ci.md` · `bei `ci-gegenprobe` ohne Unterprozesse entfällt diese Mutation` | nein (Abgleich von Hand) | Plan/Vertrag widerspricht dem Skript |
| F-5 | LOW | Zwei Sätze sagen weiter zu, als ein Test prüft: „Kein Fall hängt von Dateirechten ab“ (Kopf von `ci-gegenprobe.sh`; §7 belegt ihn mit einer einmaligen Sonde als root, nicht mit einem Test) und „weder Docker noch make zum Pruefen“ (`ci.mk`; kein Fall läuft ohne beide). | `AGENTS.md` §3.11 | `tools/harness/ci-gegenprobe.sh` · `Kein Fall hängt von Dateirechten ab`; `harness/mk/ci.mk` · `weder Docker noch make zum Pruefen` | nein | Kommentar-Zusage weiter als die Prüfung |
| F-6 | INFO | `SPEC-064` Punkte 6, 8 und 9 stehen in Spannung: Punkt 6 sagt `ci-range: <Name> fehlt oder ist ungültig` wörtlich zu, Punkt 8 nennt die Meldung zur Variable „Lesart, kein Fall“, die Fälle `p6-*` prüfen den Namen aber als Teiltext; Punkt 9 nimmt „mehrere Gründe derselben Stufe“ aus, Stufe 2 ist aber geordnet (`p9-head-vor-prbase`). Für den Architect, kein Befund gegen den Diff. | Maintainability | `harness/vertraege/SPEC-064-ci.md` · `welchen zuerst, ist nicht zugesagt` | nein | Vertrag intern uneinheitlich |
| F-7 | INFO | Ausnahme „kein erneutes `make gates`“: `git diff a31f4a7f..5f88de70 --stat` ändert drei Dateien (Plan, `SPEC-064`, `ci-gegenprobe.sh`); die Gegenprobe steht in keiner Kette ([ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md)), die Kette liest von Vertrag und Plan nur `docs-check` (hier grün) und `kopf-check` (nicht wiederholt, vom Implementer grün gemeldet). Die Ausnahme trägt. | Maintainability | `docs/plan/planning/in-progress/slice-harness-ci.md` · `make gates` lief hier nicht erneut | ja — `make kopf-check` | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.github/workflows/ci.yml` gegen `SPEC-064` Punkt 1 (Zeile für Zeile: SHA-Pin mit Tag-Kommentar `# v7.0.1`, `permissions: {}`, `contents: read` im Job, `timeout-minutes: 30`, `ubuntu-24.04`, `fetch-depth: 0`, `persist-credentials: false`, `!cancelled()`, kein `concurrency`/`workflow_dispatch`/`secrets.`) | geprüft, ohne Befund |
| `tools/harness/ci-range.sh`: fünf Variablen, flacher Klon, Nullen-Basis, Vorgänger, leere Range, Reihenfolge der Gründe (1) bis (5); Meldungen nur nach Präfix, Exit und Ort verglichen (Entscheidung des Architect) | geprüft, Befund nur F-1 und F-3 |
| `harness/mk/ci.mk`, `harness/README.md` (Zeile `make ci-gegenprobe`, Name `ci-gates`, Branch-Schutz als Handlung), `.d-check.yml` (`workflows` in `modules:`, `workflows.dir`) | geprüft, ohne Befund außer F-5 |
| Zählung 144 Läufe / 141 Mutationen: `p1` 99 (12 + 2 × 39 + 9), `p6` 24, `p7` 10, `p8` 3, `p9` 6, `p18` 2; Mutationen 99 + 36 (`mutr`) + 6 (`mut18`); Tabellen in §7 stimmen in den Stichproben mit Namen und roten Mengen im Skript überein | geprüft, ohne Befund |
| Äquivalenter Mutant `m-p8-git-fehler-ignoriert` (nur der Aufruf `is-shallow-repository`): Invariante, Ort (`case "$flach"`, Zweig `*)`) und Sonde in §7 genannt; sondiert (leere Ausgabe führt in `*)`: Exit 2 mit Präfix), für diesen Aufruf äquivalent; die Formulierung der Invariante reicht weiter, siehe F-1 | geprüft, ohne eigenen Befund |
| Abbruch (Punkt 18) für `SIGINT`/`SIGTERM`, Rekursions-Schutz (`env -u GP_MUT -u GP_ENTDECKE -u GP_FALL`, `GP_INNEN`, Lauf `rekursion-schutz`), root/umask (Sonden in §7, nicht wiederholt) | geprüft, Befund nur F-2 |
| §3.9/§3.13 (Plan folgt dem Diff, Zählung in §8); bekannte Zuweisungen an den Planner (`ubuntu-latest` in §1, `concurrency` in der DoD, „offen“ in der Einleitung von §6) | geprüft, bis auf F-4 ohne Befund |
| §3.14 Metrik: 747 gegen Marke 900, Metrik in §4 genannt | geprüft, ohne Befund |
| `make docs-check`, `bash -n` | grün |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage mit Ortsangabe/Aufzählung ohne Fall und Mutation (git-Fehler, „Ende“, unbekannter Vorgänger) · Plan/Vertrag widerspricht dem Skript · Kommentar-Zusage weiter als die Prüfung

## Verdikt

**Merge-blockierend:** ja — drei MEDIUM (F-1 bis F-3): je eine Zusage oder ein Element ihrer Aufzählung ohne Fall und Mutation, mit grünem Mutanten (gemessen).

**Übergabe:** Findings gehen an den Implementer; F-6 an den Architect, F-4 zusätzlich an Planner/Architect (Vertragstext). Die Finding-Klassen gehen in die Slice-Closure §7.
