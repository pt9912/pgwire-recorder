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

---

# Runde 2 — Nachprüfung der Behebung (slice-harness-ci)

**Gegenstand:** `impl-ci` Kopf `fea80886` (vorher `5f88de70`), Commits `7d80d9f2` und `fea80886`. Geprüft nur die Behebung von F-1, F-2, F-3, F-5; F-4 und F-6 gehören dem Architect. Reviewer-Skill `.harness/skills/reviewer.md`, Datum 2026-10-19. Engster Sensor: keine volle `ci-gegenprobe`, kein `make gates`.

## Findings (neu)

- `kategorie`: LOW
- `quelle`: Maintainability (`AGENTS.md` §3.9, Plan folgt dem Diff)
- `pfad`: `docs/plan/planning/in-progress/slice-harness-ci.md` · „mit `GP_ENTDECKE=1` gesehen, dann der volle Lauf grün (144 Läufe: 141 Mutationen“
- `befund`: Die Überschrift davor nennt 144 Mutationen, der Satz selbst noch den alten Stand (144 Läufe, 141 Mutationen); nach §7 „Läufe nach dem Review“ sind es 147 Läufe und 144 Mutationen. Die Zählung im Skript stimmt mit 144 (12 `mut1` + 38 `mutr` + 7 `mut18` + 87 aus der Schleife) und der Tabelle überein; nur dieser Satz ist nicht nachgezogen.
- `verifizierbar`: nein (kein Gate liest den Satz)
- `klasse`: Zählung im Plan nach Korrektur nicht vollständig nachgezogen

## Behebung je Finding

| Finding | behoben | Messung |
|---|---|---|
| F-1 `\|\| fehler` an `git rev-list --count` | ja | `GP_MUT` auf `ok` und `m-p8-rev-list-fehler-ignoriert` u. a.: `ci-gegenprobe: grün (11 Läufe)`, `m-p8-rev-list-fehler-ignoriert => p8-rev-list-fehler` (genau dieser Fall). Eigener Mutant (`sed` auf `\|\| true` in Kopie von `ci-range.sh`, Repo mit gelöschtem Objekt des mittleren Commits, selbst gebaut): Original Exit 2 mit Präfix `ci-range:`, Mutant Exit 0 und gibt die Range aus (dazu `[: : integer expression expected`). Der Fall ist also trennscharf. `m-p7-merge-kopf-abgelehnt` färbt jetzt zusätzlich `p8-rev-list-fehler` (§7 begründet es, gemessen: `p7-pr-merge-kopf p8-rev-list-fehler`). Die Invariante zu `m-p8-git-fehler-ignoriert` ist auf den Aufruf `--is-shallow-repository` eingegrenzt und verweist auf den neuen Fall; Invariante, Ort und Sonde stehen weiter dort. |
| F-2 Aufräumen bei normalem Ende | ja | `m-p18-ende-falle-fehlt => p18-ende` (genau dieser Fall); `ok` grün, also besteht `p18-ende` im Original (Exit 0, `TMPDIR` leer). Der Fall prüft Exit-Code und leeres Verzeichnis; der Mutant ohne `trap raeume_auf EXIT` lässt `tmp.*` liegen. |
| F-3 `CI_DEFAULT_BRANCH` bei unbekanntem Vorgänger | ja | `m-p6-default-vorgaenger-ungeprueft => p6-default-vorgaenger-fehlt-fehlt p6-default-vorgaenger-fehlt-ungueltig p7-push-vorgaenger-fehlt p8-stdout-eine-zeile`. Die zwei neuen Fälle sind in `ok` grün. Eigene Sonde: Push mit unbekanntem Vorgänger und ohne `CI_DEFAULT_BRANCH` ergibt Exit 2, Hinweis und `ci-range: CI_DEFAULT_BRANCH fehlt oder ist ungültig`. Die neuen Fälle färben erwartungsgemäß die Mutanten `m-p6-default-*`, `m-p7-vorgaenger-fehlt-exit-2`, `m-p7-push-nimmt-pr-basis`, `m-p8-hinweis-auf-stdout`; die Mengen sind im Lauf gemessen und stimmen mit der Tabelle in §7 überein. |
| F-5 Kommentare enger | ja | Kopf von `ci-gegenprobe.sh` „Kein Fall ändert Dateirechte“ beschreibt den Code (`grep chmod` findet nur die Werkzeugliste, die Vorprüfung und `raeume_auf` auf dem eigenen Arbeitsverzeichnis, keinen Fall; INFO: der Satz sagt „Kein Fall“, das Aufräumen ändert Rechte, trägt aber keinen Fall); `ci.mk` und `harness/README.md` „ruft weder Docker noch `make` auf“: im Skript stehen `make`/`docker` nur in Daten (Text von `ci.yml`) und Kommentaren, kein Aufruf. §7 Schritt 20 benennt beide Sätze mit ihrem Beleg. Die Zeile „ohne Docker und ohne make“ in der Werkzeugliste des Skriptkopfs ist dieselbe Aussage und trägt. |
| F-4, F-6 | nicht bewertet | `git diff 5f88de70..fea80886 -- harness/vertraege` ist leer; der Implementer hat den Vertragstext nicht angefasst. F-4 steht in §7 als Rückgabe an den Architect. |

## Zählungen und Prüfläufe

- Fälle 148 (99 + 26 + 10 + 4 + 6 + 3) und Mutationen 144 (99 + 38 + 7) in §3 und §7: gegen das Skript nachgezählt (`grep` je Hilfsfunktion: `mut1` 12, `mutr` 38, `mut18` 7; `FAELLE_*` im Skript), stimmig. Tabelle „`ci-range.sh` (38)“ und „Abbruch (7)“ stimmen.
- Metrik: `git diff --numstat 39f1ce8a..fea80886 -- .github tools harness .d-check.yml Makefile` ergibt 6+1, 39, 1, 11, 14, 624, 75 = **771**; gleich der Angabe des Implementers (Marke 900, Meldeschwelle 700 überschritten, keine Rückführung). §7 zählt bis `7d80d9f2`; `fea80886` ändert nur den Plan, die Zahl bleibt.
- `bash -n tools/harness/ci-gegenprobe.sh` grün; `make docs-check` 0 Befunde; `make abdeckung-check kopf-check meldungskatalog-check` Exit 0.
- Neue Randformen/Zusagen ohne Fall (§3.10/§3.11): keine. Die neuen Sätze sind Beschreibungen des Codes, die neuen Zusagen (rev-list-Fehler, Ende, Vorgänger mit Standard-Zweig) haben je Fall und Mutation. Der Test `p18-ende` nutzt `GP_MUT='^fall-ohne-funktion$'` als schnellen Lauf; ein Mutant, der nur bei anderen Läufen Reste ließe, wäre nicht gedeckt, das ist keine Zusage des Vertrags.
- §3.9: §1, §3 (Zeile Gegenprobe), §7 (Fälle, Mutationen, Läufe, Metrik, Kommentare, Nacharbeit) folgen dem Diff, bis auf das LOW oben.
- Bewertung, ob ein erneutes `make gates` nötig ist: nein, für den Inhalt. Seit dem grünen Lauf (`a31f4a7f`) änderten sich die Gegenprobe (in keiner Kette), Kommentarzeilen in `ci.mk`, eine Zeile in `harness/README.md` und der Plan; die Gates, die diese Dateien lesen (`docs-check`, `kopf-check`, `abdeckung-check`, `meldungskatalog-check`), habe ich gefahren, alle grün. Dass der Stop-Hook den Nachweis dennoch verlangen kann (`harness/README.md` steht nicht auf der Positivliste), ist Sache des Hooks und des Orchestrators, kein inhaltliches Argument.

## Negativbefunde (Runde 2)

| Bereich | Ergebnis |
|---|---|
| `tools/harness/ci-range.sh` (unverändert gegenüber Runde 1) | geprüft, ohne Befund |
| `harness/vertraege/` (Diff leer) | geprüft, ohne Befund |
| `harness/mk/ci.mk`, `harness/README.md` Zeile `make ci-gegenprobe` | geprüft, ohne Befund |
| `tools/harness/ci-gegenprobe.sh` neue Fälle, neue Mutationen, Kopf | geprüft, ohne Befund |

## Summary (Runde 2)

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen:** Zählung im Plan nach Korrektur nicht vollständig nachgezogen.

## Verdikt (Runde 2)

**Merge-blockierend:** nein. F-1, F-2, F-3, F-5 sind behoben und gemessen; das neue LOW (ein veralteter Zählsatz in §7) blockiert nicht. Offen bleiben F-4 und F-6 beim Architect (Vertragstext `SPEC-064`), nicht Gegenstand dieser Runde.
