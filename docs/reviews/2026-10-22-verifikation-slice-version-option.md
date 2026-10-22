# Verifikationsbericht: slice-version-option — 2026-10-22

**Rolle:** Verifier (Modul 11), frischer Kontext. **Frage:** Bauen wir es richtig, gegen Plan und DoD? Beginn 20:25:58, Ende 20:32 (`date`).

**Gegenstand:** Worktree `/Development/pgwire-recorder-vo`, Zweig `impl-version-option`, Stand **`17047e7c`** (Basis der Metrik `49da54b6993287b9d1902b30aacef66e8afc4703`, reiner Move `R100`). Lieferung: `ea6d1524` (Architect, `LH-FA-01.a` *Version*, §6), `5804f486` (Code, Tests, Handbuch, Tabellen, §3), `1935d6b5` (Belege in §7), Review `17047e7c`.

**Maßstab:** Plan (Kopf, §1 bis §8), `spec/spezifikation.md` `LH-FA-01.a` *Version*, `LH-FA-14.a`, `AGENTS.md` §3.3, §3.9 bis §3.14, Baseline `v6.18.0` `modul-11-verification.md`. Eingang: [Review-Report](2026-10-22-review-slice-version-option.md) (0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO).

**Grenze dieses Laufs:** nur diese Datei im Worktree; kein `make gates`, keine volle Suite (fährt die Koordination im Hauptbaum nach dem Merge); kein Push; Kopien unter `scratchpad/vo-vf/copies/` (`git archive HEAD`).

## Ergebnis in einem Satz

Beide Liefer-Punkte sind am Endstand **erfüllt**: Unit-, Bootstrap- und vier E2E-Tests grün, fünf Mutationen (a bis e) bis auf e aus dem richtigen Grund rot, Größe 273 gegen Marke 1500, `make docs-check`, `make kopf-check`, `make abdeckung-check` Exit 0; **offen** sind `make gates` (Hauptbaum) und die Closure-Pflichten (Planner); **ein LOW** (V-1, bestätigt F-1), **zwei INFO**, kein HIGH, kein MEDIUM.

## Selbst nachgefahren

| # | Lauf | Ergebnis |
|---|---|---|
| 1 | Kopie `base`, Image der Stufe `deps` mit Bind-Mount: `go test ./internal/adapters/driving/cli/ ./internal/bootstrap/ -run 'Version\|Usage\|Parse' -v` | Exit 0; `TestParseVersion`, `TestParseKeineVersionsAngabe`, `TestParseVersionMitWeiteremArgument`, `TestParseVersionHilfe`, `TestUsageNenntVersionsAngabe`, `TestRunVersion`, `TestRunVersionLogLevel`, `TestRunVersionsAngabe` je `PASS` |
| 2 | Runner `tools/test/run-integration-tests.sh` (Kopie mit eigenem Image-Tag und Präfix `pgr-vf-`), `PGR_TESTS_1` mit den vier Tests | Phase 1: `TestE2EVersionsAngabe`, `TestE2EKeineVersionsAngabe`, `TestE2EVersionNachKommando`, `TestE2EVersionHilfe` je `PASS`, `PASS` des Pakets. Der Runner endete mit Exit 1, weil ich `PGR_TESTS_2` leer ließ und er dann alle `TestE2EOhnePostgres*` ohne Aufzeichnungen der Phase 1 fährt; die vier Tests haben keine Phase 2, der Exit ist ein Artefakt meiner Auswahl, kein Befund |
| 3 | gebautes Binary (Stufe `dist` der Kopie, `VERSION=dev`) | `version`, `--version`, `-v`: Exit 0, `pgwire-recorder dev`, stderr leer; `-V`, `--version=x`, `--version record`: Exit 2, `PGR-E2001` (`unbekanntes Kommando`); `-- --version`: `unbekanntes Kommando "--"`; `record -v`: `PGR-E2001` (unbekannte Option); `--version --help`: Exit 0, globale Hilfe mit Block `Programmversion, nur als einziges Argument:`; `record --help` nennt `version` nicht |
| 4 | `git diff --numstat --find-renames=10% 49da54b6..HEAD -- internal/adapters/driving/cli test docs/user ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \| awk …` | **273** (`cli.go` 17/1, `cli_test.go` 87/7, `version_e2e_test.go` 115/0, Handbuch 23/8, Tabellen 13/2); gleich der Meldung; Marke 1500 nicht erreicht |
| 5 | `git diff --name-status -M 49da54b6..HEAD`; `--stat` über `README.md cmd internal/hexagon harness tools Makefile .claude .harness AGENTS.md spec/lastenheft.md internal/bootstrap/bootstrap.go` | elf Dateien: Plan, Review-Report, drei Tabellen, Handbuch, `cli.go`, `cli_test.go`, `bootstrap_test.go`, `spec/spezifikation.md`, `version_e2e_test.go`; die zweite Liste ist leer |
| 6 | `make docs-check`, `make kopf-check`, `make abdeckung-check` im Worktree | je Exit 0 (`d-check: 1407 Datei(en) geprüft, 0 Befund(e)`) |
| 7 | `grep -rn -- '--version' spec docs/user README.md internal test`, gefiltert auf Verneinung | ein Treffer, `abdeckung-unit.md:23` („keine Versions-Angabe“, Wortlaut der Spezifikation, keine Verneinung der Option); `README.md` 0 Treffer |
| 8 | Aufräumen | kein Container, Netz, Volume oder Image mit Präfix `pgr-vf-`; `git status` des Worktrees leer vor diesem Bericht |

## DoD-Stand (§2)

| Punkt | Verdikt |
|---|---|
| **Liefer-Punkt 1 — die Option** (ein Test je Zusage und je Element, ohne Mutation; `TestParseVersion` und Deklaration neu; `make abdeckung-check` grün) | **erfüllt.** Lauf 1 bis 3 und 6; `TestParseVersion` prüft `--version` und `-v` als `Command{Name: "version"}`, die Deklaration nennt das neue Verhalten, die alte Verneinung steht nirgends mehr (Lauf 7). Bewusstes Brechen: a, b, c, d |
| **Liefer-Punkt 2 — Doku im Ist-Zustand** (Hilfetext, Handbuch, jeder Satz mit Lauf gegen das Binary, keine Chronik, kein Verweis nach oben) | **erfüllt.** Lauf 2 und 3: jeder Satz des Handbuchs (drei Formen, eine Zeile, stderr leer, Exit 0, Schreibweisen, sechs Beispiele, `record`/`replay`/`play`, Hilfe, nur globale Hilfe) deckt sich mit dem Binary und ist in §7 je Zeile einem E2E-Test zugeordnet; der Handbuch-Diff enthält keinen Verweis auf Spezifikation, ADR, Slice oder Welle und keine Chronik. **Grenze:** Die Zuordnung ist Lesart, kein Gate (V-2) |
| `make gates` grün | **offen** (Hauptbaum, Koordination). Gefahren und grün: Läufe 1 bis 3, 6 |
| Review-Report liegt vor | **erfüllt.** Anderer Kontext, 0 HIGH, 0 MEDIUM; seit `1935d6b5` ändert nur der Report den Zweig |
| Closure-Notiz mit Lerneintrag · Register · Risiken-Ausgänge · drei Paarungen | **offen (Planner).** Beleg für das Risiko 2 (§6, `grep`): Lauf 7, keine Verneinung am Endstand; Risiko 1 (`-v`): `record -v` bleibt `PGR-E2001` (Lauf 3) |

## Bewusstes Brechen

Jede Mutation in eigener Kopie `scratchpad/vo-vf/copies/mut-<x>/`; Unit mit Bind-Mount, E2E über den Runner.

| Lauf | Mutation | Ergebnis |
|---|---|---|
| a | `-v` aus `versionsAngabe` entfernt | **rot:** `TestParseVersion` (`["-v"]: erwartet das Kommando version`), `TestRunVersionsAngabe` (Exit 2, `unbekanntes Kommando "-v"`), E2E `TestE2EVersionsAngabe` |
| b | `len(args) >= 1` statt `== 1` (`--version record` akzeptiert) | **rot:** `TestParseVersionMitWeiteremArgument` (sechs Beispiele, `["--version" "record"]: erwartet PGR-E2001`), E2E `TestE2EKeineVersionsAngabe` (`Exit-Code 0`); `TestE2EVersionNachKommando` bleibt grün (`record --version` trifft die Mutation nicht, erwartet) |
| c | Ausgabe `pgwire-recorder:` statt `pgwire-recorder` | **rot:** `TestRunVersion`, `TestRunVersionLogLevel`, `TestRunVersionsAngabe` (`stdout "pgwire-recorder: 1.2.3\n"`), E2E `TestE2EVersionsAngabe` |
| d | Exit-Code 1 für `version` | **rot:** dieselben drei Bootstrap-Tests (`Exit-Code 1`), E2E `TestE2EVersionsAngabe` |
| e | Handbuch: „Exit-Code 0“ → „Exit-Code 1“ und `-V` aus der Liste der Schreibweisen gestrichen | **grün in allen Gates:** `make docs-check` Exit 0 (1407 Dateien, 0 Befunde); `make abdeckung-check` Exit 0. Kein Sensor bindet einen Handbuchsatz an seinen Lauf (V-2) |

## Entscheidungs-Konformität

| Prüfpunkt | Befund |
|---|---|
| Diff nur CLI-Adapter, Tests, Handbuch, Tabellen, Plan, Spezifikation (Architect) | **konform** (Lauf 5); `internal/bootstrap/bootstrap.go` und `internal/hexagon` unberührt, `bootstrap_test.go` ist nur Test |
| Spezifikation vor dem ersten Code-Commit | **konform:** `ea6d1524` (Architect) liegt vor `5804f486` (Code); Lastenheft, Verträge, ADRs, Gates unberührt |
| `README.md` unberührt (§1) | **konform** (Lauf 5 und 7) |
| §1 *Ausdrücklich NICHT* (`version` bleibt, kein Kurzoptions-Schema, `record -v` bleibt `PGR-E2001`) | **konform** (Lauf 3; `TestE2EVersionNachKommando`) |
| Entscheidungen des Nutzers: `-v` belegt, nur als einziges Argument | **konform** (Lauf 3, Mutationen a und b) |
| `AGENTS.md` §3.3 (Move rein) | **konform:** `49da54b6` ist `R100`, Inhalt in eigenen Commits |
| `AGENTS.md` §3.10 (keine Mutationspflicht, §6 *Randformen*) | **nicht ausgelöst**; die Mutationen oben sind das Brechen für die DoD |
| `AGENTS.md` §3.11 (Handbuch Ist-Zustand, Deklarationen) | **konform** (Lauf 2, 3, 6); Hilfetext und Handbuch stimmen mit dem Binary überein |
| `AGENTS.md` §3.14 | **konform:** 273 gegen 1500, Befehl wörtlich (Lauf 4) |
| §4 Rückführungs-Bedingungen | **nicht eingetreten:** Marke nicht erreicht, die Spezifikation kam vor dem Start, kein Kurzoptions-Schema |

## Plan gegen Code (§3.9)

**V-1 (LOW, bestätigt F-1 des Reviews).** §3 des Plans führt seit `5804f486` den Stand der Lieferung (`bootstrap_test.go`, `version_e2e_test.go` mit vier Tests, drei Handbuchstellen: `git diff ea6d1524 5804f486 -- docs/plan` ersetzt vier Zeilen der Tabelle). §7 *Abweichungen* (`1935d6b5`) nennt als Abweichungen vom Plan „E2E-Fall oder `bootstrap_test.go`“ und „drei Stellen statt des einen Absatzes“; beides steht in §3 schon, die Abweichung beschreibt den Planstand vor `5804f486`. Der Code folgt §3 und §6. Der Planner trägt es bei der Closure nach (§7 *Abweichungen* streichen oder auf das beschränken, was §3 nicht führt: die zwei Zusatzfälle `--version=` und `-v=`, die unter „Argument mit `=`“ der Spezifikation fallen).

**V-2 (INFO).** Handbuchsätze zur Version sind nur über die Tabelle in §7 an ihre E2E-Probe gebunden; Lauf e zeigt, dass ein falscher Satz kein Gate rot färbt. Das ist die gelebte Praxis des Repos für Handbuchsätze (`AGENTS.md` §3.11 verlangt den Lauf gegen das Binary, kein Gate dafür); kein Befund gegen den Slice.

**V-3 (INFO).** Der Handbuch-Kopf nennt Version 0.11, Software-Version und Gültigkeitsbereich weiter `0.1.1`, obwohl das Release 0.1.1 die Option nicht kennt (F-2 des Reviews, gleiche Lesart: Kopf zieht beim nächsten Release nach).

## Nicht geprüft

`make gates` und die volle Integrationssuite; macOS und Windows (Grenze der Spezifikation, `\n` nur unter Linux); Phase 2 des Runners (die vier Tests haben keine); Gegenproben (kein Gate und kein Vertrag berührt).

## Offen für den Planner

Closure-Notiz und Lerneintrag; Register (V-1 als Klasse „Plan-Abweichung beschreibt überholten Planstand“); Ausgänge der zwei Risiken aus §6; Paarungen; `make gates` am Endstand im Hauptbaum.
