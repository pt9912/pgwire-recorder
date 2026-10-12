# Verifikation: slice-harness-lint-gegenprobe-last — 2026-10-12

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-12-review-slice-harness-lint-gegenprobe-last.md`).

**Gegenstand:** Slice-Plan `slice-harness-lint-gegenprobe-last` (Kopf, §1 bis §8) gegen `2d1f0ca..9c4b0a6`. Commits: `10dc4c8` (Skript), `b4ee633` (Belege in §7), `9c4b0a6` (Review-Report). Entscheidung: [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

**Eingang:** Plan ganz (§7 *Belege des Implementers*), Diff von `tools/harness/lint-gegenprobe.sh`, Review-Report. Den Bericht des Implementers habe ich nicht gelesen.

## 1. Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `grep -n 'parallel=' tools/harness/lint-gegenprobe.sh` | Zeile 125 `parallel=3`; Zweig für gesetzte Werte (`^[1-9][0-9]*$`) unverändert |
| `grep` nach `6`/`sechs` im Skript | keine Fundstelle mehr als Standard (Treffer nur Punkt-(6)-Beschriftungen und Parameter `${6:-.}`) |
| `LINT_GEGENPROBE_PARALLEL=<wert> DOCKER_BUILD=false bash tools/harness/lint-gegenprobe.sh` für `0`, leer, `-1`, `abc`, ` 3` | je Exit 2 und die feste Zeile `lint-gegenprobe: LINT_GEGENPROBE_PARALLEL ist keine positive ganze Zahl: '<wert>'`; endet vor der ersten Kopie ohne Docker |
| `grep -rn` nach `LINT_GEGENPROBE`/`PARALLEL` in `tools/`, `harness/`, `spec/`, `docs/user`, `README.md`, `AGENTS.md`, `.claude/`, `Makefile` | nur das Skript; `Makefile` Zeile 4 nennt `make -j` (anderer Gegenstand) |
| `make docs-check` | `620 Datei(en) geprüft, 0 Befund(e)` |

**Nicht selbst gefahren:** `make gates`, `make lint-gegenprobe` (150 bis 170 s, Last bis 40) und die Messungen in §7; ebenso der Fall mit nicht gesetzter Einstellung und echtem Docker-Lauf. Für das Gate gilt der Auftrag (`make gates` an `b4ee633` Exit 0), für die Messung §7. Der Standard 3 selbst ist durch keinen Lauf geprüft (akzeptiertes Negativ, siehe unten).

## 2. DoD Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1 Standard 3; Kopf nennt 3, keine Prüfzusage durch `make gates`, keine Aussage zu Last, Speicher, OOM; Satz „den Default fährt jeder Lauf von `make gates`“ gestrichen; 6 im Kopf ersetzt; 3 als Höchstzahl je Lauf des Skripts, nicht Obergrenze der Last; kein Fall für den Standard; gültiger gesetzter Wert unverändert | erfüllt | Zeile 125 `parallel=3`. Kopf Zeilen 13 bis 18: „Höchstzahl gleichzeitiger Läufe der Stufe je Lauf dieses Skripts, keine Obergrenze der Last von `make gates`. Nicht gesetzt: 3.“; `GEPRÜFT DURCH einstellung-<wert> (0, leer, -1, abc, ` 3`).` ohne den Satz zu `make gates`. Kein Wort zu OOM, Last, Speicher, Zeit. Diff berührt nur diese Zeilen; die Zweige für gesetzte und ungültige Werte sind unverändert (Probe oben). Kein neuer Fall. |
| 2 `make gates` grün | erfüllt laut Auftrag; Beleg in §7 fehlt | siehe V-1 |
| 3 Review, Report unter `docs/reviews/` | erfüllt | Report liegt vor: HIGH 0 · MEDIUM 0 · LOW 0 · INFO 3 |
| 4 Closure-Notiz mit Lerneintrag | offen (Closure) | §7 trägt bis jetzt nur *Belege des Implementers* |
| 5 Beobachtungs-Register: `evidence/slice-harness-lint-gegenprobe-last.md` unter `gate-laufzeit-waechst-mit-gegenprobe` | offen (Closure) | Ordner enthält bisher nur `slice-harness-lint.md`; die Messung steht in §7 und wartet auf den Eintrag |
| 6 Risiken mit Ausgang | offen (Closure) | §6 trägt Ausgänge vorab (zwei *weiter offen*, einer *entfallen*); der Neueintrag `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert` ist noch nicht angelegt |
| 7 Paarungen | offen (Closure) | |

## 3. Plan gegen Code

| Plan | Code | Urteil |
|---|---|---|
| §3: `parallel=6` wird `parallel=3`; Kopf nennt 3 und streicht den Satz | wie beschrieben (Diff `10dc4c8`, 9 Zeilen) | stimmt |
| §3: keine weitere Fundstelle, kein Makefile-Fragment, kein Test | Suche (oben) bestätigt; Diff berührt nur das Skript und den Plan | stimmt |
| §6 Randform *ungültiger Wert*: unverändert | Probe: Exit 2, feste Zeile | stimmt |
| §6 *gültiger gesetzter Wert* überstimmt den Standard | Zweig unverändert; `lint.mk` nicht berührt | stimmt (nicht mit Docker gefahren) |
| §6 *Lesart des Standards 3* | Kopf liest 3 als Höchstzahl je Lauf des Skripts | stimmt |
| §6 Mutation `parallel=3` zu `parallel=6` bleibt grün | akzeptiertes Negativ (Entscheidung des Nutzers, 2026-10-12); kein Finding | wie entschieden; Kopf sagt keine Prüfung zu |
| Kopf des Plans: Spec-Stelle `—`, kein Lastenheft-Bezug | Diff berührt weder `spec/` noch Doku | stimmt |

Zählung nach `AGENTS.md` §3.13 (ein Liefer-Punkt, eine Schicht): stimmt mit dem Diff. Kein Move in den Code-Commits (§3.3).

## 4. AGENTS.md §3.11

Der Kopf sagt nicht zu, dass ein Lauf den Standard prüft, und nichts über den OOM; der Plan (§1, §6, §7 *Nicht belegt*) sagt ebenfalls nicht zu, dass die Senkung den OOM behebt. Erfüllt.

## 5. Befunde

**V-1 (LOW) — §7 trägt keinen Gate-Lauf als Beleg.** DoD-Punkt 2 und §5 verlangen `make gates` grün; §7 *Belege des Implementers* nennt Messung, Mutation und Suche, aber keinen Lauf von `make gates` mit Stand und Ergebnis. Das Ergebnis steht nur im Auftrag und im Review-Report („laut Auftrag, nicht von diesem Lauf gefahren“). Nachzutragen vor der Closure: Stand (Commit) und Exit-Code. Nicht eigenständig verifiziert (nicht gefahren).

**V-2 (INFO) — Kopf Zeilen 14 und 15 sind umgebrochen wie ein Fließtext, Zeile 15 hat 136 Zeichen.** Kein Konventions-Anker, kein Befund gegen eine Regel; Lesbarkeit. Deckt sich mit INFO-3 des Reviews.

**V-3 (INFO) — Das Review-INFO-1 („etwas höher“ bei Load-Maximum 38,0 gegen 31,9) und INFO-2 (Suchumfang `docs/user` statt `docs/`) stehen unbehoben im Plan.** Meine eigene Suche über `docs/user` stimmt mit §7; die Review-Suche über das ganze Repo fand keine 6 als Standard. Für die Lerneintrag-Zeile oder die Evidence-Datei bei der Closure: die Zahl nennen, nicht „etwas“.

Kein HIGH, kein MEDIUM.

## 6. Offen für die Closure

Closure-Notiz mit Lerneintrag; Evidence-Datei `slice-harness-lint-gegenprobe-last.md` im Register (Zähler dann 2×); Neueintrag `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert`; Ausgänge der drei Risiken bestätigen; Paarungen; Gate-Beleg in §7 (V-1).

## 7. Verdikt

DoD-Punkte 1 und 3 erfüllt, Punkt 2 erfüllt laut Auftrag mit fehlendem Beleg (V-1), Punkte 4 bis 7 sind Closure. Plan und Code stimmen überein; das akzeptierte Negativ ist entschieden und kein Finding.
