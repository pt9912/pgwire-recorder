# Review-Report: slice-harness-ci-workflows — 2026-10-18

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; das ist die Frage des Verifiers).

**Gegenstand:** `slice-harness-ci-workflows`, Zweig `rev-ci-workflows`, Kopf 486b45d0 (Zweig `impl-ci-workflows`; der Worktree stand auf 88066ffe, das Delta `88066ffe..486b45d0` wurde mit `git diff` gelesen), Diff `79b82a0e..486b45d0` (Code-Commits 7a66cc39, a207e17c, 0c30b767, Symlink-Fall c92ed918; Plan-Belege; SPEC-063-Nachtrag 2890e31f des Architect als Basis des Urteils; Symlink-Entscheidung Punkt 2 F4 und Punkt 9).

**Skill:** `.harness/skills/reviewer.md` @ 88066ffe ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

> **Zitier-Form.** Dieser Report friert ein; Slices werden mit Kennung zitiert, Gates als `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code.

**Eingangs-Kontext:**

- `slice-harness-ci-workflows` (Plan in `in-progress/`, §7 *Belege des Implementers*: Mutationstabelle, Läufe, Metrik; der Bericht des Implementers lag nicht vor)
- [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md) (Accepted 2026-10-18), [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md), [ADR-0047](../plan/adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md)
- SPEC-063 (alle Punkte), SPEC-060 Teil 1 (Listen, Namensregel), SPEC-052 Punkte 1 bis 4 (Aufrufarten, Ausnahmeliste)
- Nehmer `slice-harness-ci` (§1, DoD, §8: Sendung `.d-check.yml` aus Punkt 8, gelesen)
- `AGENTS.md` (Hard Rules §3.2, §3.7 bis §3.12, §3.14)

**Gefahrene Sensoren (nur der engste):** `bash -n` über beide Skripte (ok), `make docs-check` (0 Befunde), `bash tools/harness/workflows-check.sh` am Arbeitsbaum (Exit 0, Meldung „keine Workflow-Dateien“), `git diff --numstat a98d0b9e..88066ffe -- .github tools harness Makefile` (Summe am Kopf 486b45d0: 849, stimmt mit §7; am Kopf 88066ffe 830). Dazu 10 zusätzliche Mutationen in einer Kopie unter dem Scratchpad (nicht im Repo), gefahren mit `GP_MUT` und `GP_ENTDECKE=1` ohne die echten Fälle (`bash tools/harness/workflows-check-gegenprobe.sh` mit Platzhaltern für die Images). Kein `make gates`, keine volle Gegenprobe, kein Lauf der echten Fälle.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | SPEC-063 Punkt 8 (c) bindet `dir:` „darunter, vor dem nächsten Schlüssel in Spalte 1“, und (a) verlangt das Element `workflows` in `modules:`. Kein Fall hat ein `dir: .github/workflows` unter einem anderen Block (der Fall `p8-dir` nutzt `dir: docs`), und kein Fall hat ein Element, das nur `workflows` enthält. Die Mutationen „Rücksetzen bei Schlüssel in Spalte 1 gestrichen“ und „Element als Muster `~ /workflows/`“ blieben in der Kopie grün (kein Fall rot). Folge: das Gate nähme `modules: [workflows-x]` oder ein `dir:` eines fremden Blocks als Bindung an, obwohl `make docs-check` die Pins dann nicht prüft. | `AGENTS.md` §3.10 (Ortsangabe und Aufzählung zählen je Element) · SPEC-063 Punkt 8 | `tools/harness/workflows-check.sh` · `/^[^ \t#]/ { inb = 0 }` und `if (x == "workflows") a = 1` | ja — die beiden Mutationen in einer Kopie der Gegenprobe | Ortsangabe und Element einer Zusage ohne trennenden Fall |
| F-2 | MEDIUM | SPEC-063 Punkt 9 nennt „ein nicht lesbares Verzeichnis“ als Exit 2 neben der nicht lesbaren Datei. Das Skript hat die Prüfung `[ -r "$dir" ] && [ -x "$dir" ]`, die Gegenprobe hat dazu weder Fall noch Mutation (`p9-nicht-lesbar` prüft nur die Datei). Die Mutation „Prüfung gestrichen“ blieb grün; ohne die Prüfung liest `find` ein Verzeichnis ohne Rechte als leer und das Gate endet mit Exit 0 und „keine Workflow-Dateien“, also grün über einem ungeprüften Bereich. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) · SPEC-063 Punkt 9 | `tools/harness/workflows-check.sh` · `{ [ -r "$dir" ] && [ -x "$dir" ]; } \|\| sperre "$WF nicht lesbar"` | ja — Mutation „Zeile ersetzt durch `:`“ in einer Kopie, kein Fall rot | Element einer Aufzählung ohne Fall und Mutation |
| F-3 | MEDIUM | Die Gegenprobe überspringt `p9-nicht-lesbar` als root laut (`return 0`), die Mutation `m-p9-nicht-lesbar` erwartet aber den roten Fall `p9-nicht-lesbar`. Als root bleibt der Fall grün, die Mutation bleibt grün, und die Gegenprobe meldet „erwartet rot“; SPEC-063 Punkt 9 sagt „als root nicht herstellbar, dann laut übersprungen“. Läuft die Gegenprobe als root (etwa im CI-Container der folgenden Slices), ist sie rot, ohne dass das Skript falsch ist. | SPEC-063 Punkt 9 · `AGENTS.md` §3.11 (der Kommentar sagt „laut übersprungen“, die Gegenprobe tut es nur zur Hälfte) | `tools/harness/workflows-check-gegenprobe.sh` · `mut m-p9-nicht-lesbar "p9-nicht-lesbar"` | ja — `bash tools/harness/workflows-check-gegenprobe.sh` mit `GP_MUT=m-p9-nicht-lesbar` als root (am Host nicht gefahren, es läuft als uid 1000; aus dem Code gelesen) | Überspringregel gilt für den Fall, nicht für seine Mutation |
| F-4 | LOW | SPEC-063 Punkt 5 sagt, bei einem Befund `form` laufen die Regeln der Punkte 6 und 7 für die Datei nicht. Die Fälle `p5-tab`, `p5-einzug` und `p5-jobs-flow` sind so gebaut, dass die Regeln ohnehin nichts fänden (kein Job mit erkennbarem `runs-on`); die Mutation „`exit` im `END` nach den Form-Befunden entfällt“ blieb grün. Die Zusage „laufen nicht“ ist damit nicht von „laufen und finden nichts“ zu unterscheiden. | `AGENTS.md` §3.10 · SPEC-063 Punkt 5 | `tools/harness/workflows-check.sh` · `if (nf > 0) { for (i = 1; i <= nf; i++) print fl[i] "\tform\t" ft[i]; exit }` | ja — Mutation „`; exit }` durch `}`“ in einer Kopie, kein Fall rot | Zusage „Regel läuft nicht“ ohne trennenden Fall |
| F-5 | LOW | Mehrere Zusagen haben einen Fall, aber keine eigene Mutation in `mut` und keine Zeile in der Tabelle von §7: der Hinweis `hinweis: … wird nicht geprüft` samt Grund (Endung, Unterverzeichnis), die Zeilen von `actionlint` unverändert in der Ausgabe (Exit 1), der Fehlerausgabe-Anhang bei Exit außer 0 und 1. Die vier Mutationen, die ich dazu fuhr, wurden von den vorhandenen Fällen rot gefärbt; es fehlt die dokumentierte, vom Implementer gefahrene Mutation je Zusage. | `AGENTS.md` §3.10 (Bericht nennt je Zusage: Zusage · Mutation · roter Test) | `docs/plan/planning/done/slice-harness-ci-workflows.md` · `Zusage · Mutation · roter Test.` | nein — Lücke im Beleg, kein Gate-Lauf | Zusage mit Fall, aber ohne gefahrene Mutation im Beleg |
| F-6 | LOW | Die Zeile von `make workflows-check-gegenprobe` in `harness/README.md` sagt „Mutationen …, die nur die Fälle ihrer Gruppe rot färben“. Drei Mutationen (`m-p1-nur-yaml`, `m-p6-job-pflicht`, `m-p9-summe-nur-al`) sind „breit“: sie verlangen nur den Kern rot und einzelne Fälle grün und setzen die übrigen roten Fälle als erwartet (`rote="$ist"`). Die Aussage geht weiter als die Prüfung; §7 nennt die Breite, die Zeile nicht. | `AGENTS.md` §3.11 | `harness/README.md` · `Mutationen von Skript und Fragment, die nur die Fälle ihrer Gruppe rot färben` | nein | Prosa über ein Gate weiter als die Gegenprobe |
| F-7 | LOW | Der Plan folgt dem Diff nicht in zwei Stellen: §6 leitet die Randformen mit „alle offen für den Architect, zu entscheiden vor dem ersten Code-Commit“ ein, obwohl jede Randform darunter als „Entschieden“ steht; DoD Liefer-Punkt 1 sagt „nennt je Befund Datei, Job und Regel“, die Ausgabe nennt Pfad, Zeile und Regel (der Job steht nur im Text der Regeln für Job-Eigenschaften). | `AGENTS.md` §3.9 | `docs/plan/planning/done/slice-harness-ci-workflows.md` · `**alle offen für den Architect**` und `nennt je Befund Datei, Job und Regel` | nein | Plan folgt der Korrektur nicht |
| F-8 | INFO | Größenmetrik: Summe aus Einfügungen und Löschungen `a98d0b9e..486b45d0` über `.github tools harness Makefile` ist 849 (830 bei 88066ffe), wie §7 sagt (Meldeschwelle 700 überschritten, Marke 900 nicht). Die Gegenprobe trägt 542, das Skript 240, das Fragment 24, die Gate-Zeilen in `harness/README.md` 2, die Vertragsdatei 22; nichts im Diff liegt außerhalb des Slice. Die Schätzung des Architect (545 bis 695) lag unter der Messung. | `AGENTS.md` §3.14 | `docs/plan/planning/done/slice-harness-ci-workflows.md` · `Metrik gegen die Marke von §4` | ja — `git diff --numstat a98d0b9e..486b45d0 -- .github tools harness Makefile` | Schätzung unterhalb der Messung |
| F-9 | INFO | Das Fragment führt `export DOCKER` global; damit steht `DOCKER` in der Umgebung aller Rezepte der Sitzung, nicht nur dieses Gates. Ich fand kein Skript, das daran Schaden nimmt (`harness/mk/build.mk` setzt `DOCKER ?= docker`), nenne es als Kopplung für den, der ein Skript ändert, das `DOCKER` aus der Umgebung liest. | Maintainability | `harness/mk/workflows.mk` · `export DOCKER` | nein | globaler Export aus Fragment |
| F-10 | INFO | Punkt 7 liest „dezimale Ganzzahl von 1 bis 60“; das Skript nimmt `010` als 10 an und lehnt `0060` ab (Länge über 3). Führende Nullen sind in §6 keine genannte Randform. Der Fall „Symlink mit anderer Endung“ (Hinweis `(Endung)`) hat ebenfalls keinen Fall; die Mutation „Hinweis entfällt“ blieb grün. Beides für den Architect, kein Befund gegen den Implementer. | `AGENTS.md` §3.12 | `tools/harness/workflows-check.sh` · `bad = (length(tval) > 3 \|\| z < 1 \|\| z > 60)` | ja — Probelauf mit `timeout-minutes: 010` | Randform nicht genannt |

Der Fall „Prüfbereich nur ein Symlink“ ist inzwischen entschieden (SPEC-063 Punkt 2 F4 und Punkt 9: Befund `symlink`, Exit 1, kein Aufruf von `docker`, Punkte 4 und 8 nicht geprüft). Das unveränderte Skript setzt das um; Fall `p2-nur-symlink` und die Mutationen `m-p2-f4-wie-f3` und `m-p2-f4-docker` trennen die Zweige (Mutationen gelesen, nicht gefahren). Kein Befund.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/workflows-check.sh` · Reihenfolge der Gründe (Punkt 9), Exit-Codes 0/1/2, Summe aus beiden Quellen, Umgebungsfehler vor Gruppen (b) und (c) | geprüft, ohne Befund (gegen Vertrag gelesen; Fälle und Mutationen vorhanden) |
| `tools/harness/workflows-check.sh` · Aufruf von `actionlint` (Punkt 3, E1 bis E8), kein `-ignore`, `-shellcheck`, `-pyflakes`, `--network none`, `:ro`, kein `--name` | geprüft, ohne Befund |
| `tools/harness/workflows-check.sh` · Suppression, Netz- und Schreibzugriffe | geprüft, ohne Befund (kein `-ignore`, keine Konfiguration von `actionlint` zugelassen; schreibt nur in `mktemp` außerhalb des Baums, `trap` räumt auf) |
| `harness/mk/workflows.mk` · Rezept ohne `$(MAKE)` und ohne `+`, `GATE_CHECKS`, `GEGENPROBEN_CONTAINER`, Digest-Form (Index-Digest, 64 Hexziffern, Tag als Kommentar) | geprüft, ohne Befund |
| `harness/mk/gates-schnell.mk` · `workflows-check` nicht in der Ausnahmeliste der schnellen Kette | geprüft, ohne Befund (Datei im Diff unverändert, Eintrag fehlt dort) |
| Aufrufarten `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` (Punkt 10) | geprüft, ohne Befund (Fall je Flag, Mutation `+` und Zweitaufruf vorhanden) |
| `tools/harness/workflows-check-gegenprobe.sh` · Punkte 1, 2, 3, 4, 6, 7, 10 und 12 | geprüft, ohne Befund außer den Lücken in F-1, F-2, F-4, F-5 (Zahlen der Mutationen: 80, davon 10 eigene Probemutationen in der Kopie) |
| `tools/harness/workflows-check-gegenprobe.sh` · Kommentare (§3.7) | geprüft, ohne Befund (Zusage, Kopplung, Grenze; keine verworfene Alternative, kein Verweis auf abwesenden Text) |
| `harness/README.md` · Gate-Zeile `make workflows-check` | geprüft, ohne Befund (Regel an den Schreiber, Verweis auf `SPEC-063` und `.d-check.yml`; die Aussage über die Pins nennt `make docs-check` als allein zuständig, die Fälle `p8-pins-nicht-hier` und `real-pins` stützen sie) |
| `harness/vertraege/SPEC-063-workflows-check.md` · Änderungen im Diff | geprüft, ohne Befund (Rückgabe Image-Herkunft und Rechtehinweis vom Architect entschieden, vor dem Code-Commit 7a66cc39 geschrieben) |
| Delta `88066ffe..486b45d0` (SPEC-063 Punkt 2 und 9, Fall `p2-nur-symlink`, zwei Mutationen, §6 Randform 1, §7) | geprüft, ohne Befund (Vertrag, Fall, Mutationen und Tabelle in §7 stimmen überein; Skript unverändert) |
| `docs/plan/planning/done/slice-harness-ci-workflows.md` · Kopf, §3, §4, §8, Sendung an `slice-harness-ci` | geprüft, ohne Befund außer F-7 (Nehmer führt die Sendung in §1, DoD und §8 mit der Kennung des Gebers und zählt nach) |
| `make docs-check` am Endstand | geprüft, ohne Befund (0 Befunde) |
| Core-Reinheit, `internal/**`, `cmd/**`, ADRs, Lastenheft, Spezifikation, Architektur | nicht berührt, geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Ortsangabe und Element einer Zusage ohne trennenden Fall · Element einer Aufzählung ohne Fall und Mutation · Überspringregel gilt für den Fall, nicht für seine Mutation · Zusage „Regel läuft nicht“ ohne trennenden Fall · Zusage mit Fall, aber ohne gefahrene Mutation im Beleg · Prosa über ein Gate weiter als die Gegenprobe · Plan folgt der Korrektur nicht · Schätzung unterhalb der Messung

## Verdikt

**Merge-blockierend:** ja — F-1 bis F-3 (MEDIUM) blockieren: F-1 und F-2 sind Zusagen von SPEC-063 ohne Fall und Mutation (grün über einer falschen Bindung bzw. einem nicht lesbaren Bereich), F-3 färbt die Gegenprobe als root rot. F-4 bis F-7 sind LOW und blockieren nicht; die Reihenfolge der Gründe, die Exit-Codes, die Aufrufarten und die Kette (`GATE_CHECKS`, `GEGENPROBEN_CONTAINER`, nicht in der Ausnahmeliste) sind wie im Vertrag umgesetzt.

**Übergabe:** Findings gehen an den Implementer; F-10 und der offene Fall „nur ein Symlink“ gehen an den Architect. Die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; die Konformität gegen die DoD prüft der Verifier.
