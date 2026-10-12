# Verifikation: slice-v1-abschluss-zeitangaben-doku — 2026-10-12

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-12-review-slice-v1-abschluss-zeitangaben-doku.md`).

**Gegenstand:** Slice-Plan `slice-v1-abschluss-zeitangaben-doku` (Kopf, §1 bis §8) gegen `56d7797~1..236a39e`. Commits: `56d7797` (Move), `15583ea` (Roadmap-Zeile), `c2c7272` (Handbuch 0.4, README, Belege), `445f074` (Review-Report), `8209d87` (Nacharbeit F-660 bis F-663), `aa6eed4` (Mutanten-Zeilen entfernt), `236a39e` (`.claude/commands/implement-slice.md`).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

**Eingang:** Plan ganz (§7 *Belege des Implementers*), Handbuch- und README-Diff, Review-Report. Den Bericht des Implementers habe ich nicht gelesen.

## 1. Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make gates` am Stand `236a39e`, Arbeitsbaum sauber | Exit 0 (letzte Zeile `lint-gegenprobe: gruen`) |
| `make build`, dann Befehlsfolge aus `done/slice-doku-ist-stand` §7 mit `grep -vx PGWIRE_RECORDER_PASSWORD` | Optionen: nur `--h`, `--help`, `--name`, `--rm`, `--version`; Variablen: leer; Codes: leer. Entspricht der Tabelle in §7. |
| `grep` im Handbuch auf `ADR-`, `SPEC-`, `LH-`, `[Ss]lice`, `Welle`, `Spezifikation`, `Lastenheft`, `[Rr]eview` | kein Treffer |
| `grep` auf `noch nicht`, `kommt`, `geplant`, `bisher`, `später`, `compare-responses` in Handbuch und README | Treffer nur im Bestand (Handbuch Zeilen 14, 508, 794, 905: Wörter in anderem Sinn, nicht zu den Zeitangaben); README Zeile 69 ist Bestand |
| `grep` auf `vor jeder Anfrage`, `Beginn der Verbindung`, `seit dem Beginn` | kein Treffer mehr (F-660, F-662 umgesetzt) |
| Proben am Binary ohne Datenbank (Startfehler, vor jeder Verbindung): Aufzeichnung mit fehlender und absteigender Anfrage; fehlend vor Anfrage 1; `--timing-mode` allein; `--keep-timing --timing-reference connect`; `relative` mit `--timing-reference`; `1.5` als `offset_ms`; `replay --keep-timing`; `--keep-timing=1`; `--timing-mode Absolute`; Umgebung `KEEP_TIMING=Geheim123`; Umgebung `TIMING_MODE=absolute` ohne `KEEP_TIMING`; Datei `keep_timing: yes` und `record_timing: yes` | Alle wie im Handbuch: `PGR-E3003` (Exit 3) geht `PGR-E2003` (Exit 2) vor; Kombinationen `PGR-E2001`; Datei `PGR-E2004` mit Schlüssel; die Meldung zur Umgebungsvariable nennt den Wert nicht; `replay` kennt die Option nicht (`PGR-E2001`) |

Nicht selbst gefahren: die Proben, die einen PostgreSQL brauchen (Wartezeiten, `offset_ms`-Werte, Abbruchsignal). Dafür gilt der Beleg in §7 (Probe-Matrix, Rohwerte vom Server-Log, Stand `15583ea` und Nacharbeit); er ist in sich schlüssig und deckt sich mit den Stichproben des Reviewers.

## 2. DoD Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1 Handbuch im Ist-Zustand: Abschnitte, Optionen gegen `--help`, Codes im Katalog und ausgelöst, Beispiele als Datei, keine Verweise, Lesarten aus §1 (Ordnung der Startfehler, Kombinationen, Wertemenge) | erfüllt | Neuer Abschnitt *Zeitangaben aufzeichnen und einspielen* mit *Ungültige Kombinationen*, *Fehlende Zeitangaben und Abbruch*, *Prüfung beim Laden*; Tabelle in §5 (4 Zeilen, Standards stimmen mit `--help`); Fehlercodes `PGR-E2001`/`E2003`/`E2004`/`E3003`; SQLite-Tabelle, FAQ, Glossar, Einleitung. Die Befehlsfolge ist grün (Abschnitt 1), die vier Beispielblöcke stehen als Läufe in §7. Die drei Lesarten stehen im Handbuch und sind in §7 geprobt (eigene Proben bestätigen Ordnung, Kombinationen und Wertemenge). Messpunkt (nach Startup) steht in Einleitung, Abschnittsanfang, Hinweis, Glossar und README. |
| 2 README | erfüllt | Ein Absatz (Zeilen 37 bis 43), Ist-Zustand, kein Zielstand; Verweise auf `spec/` und `docs/plan/` unverändert. Kein überholter Satz zu Zeitangaben im README gefunden. |
| 3 `make gates` grün | erfüllt, Beleg fehlt in §7 | Mein Lauf: Exit 0. §7 trägt jedoch keinen Gate-Lauf mit Stand und Ergebnis (V-203). |
| 4 Review, Report unter `docs/reviews/` | erfüllt | Report liegt vor; 0 HIGH, 0 MEDIUM, 3 LOW, 2 INFO. Das Verdikt sagt „Gate-Ergebnis siehe Bericht des Reviewers“, der Report nennt es nicht (V-205). |
| 5 Closure-Notiz mit Lerneintrag | offen (Closure führt der Implementer nach diesem Bericht) | §7 bis zur Closure unvollständig, kein Befund |
| 6 Beobachtungs-Register | offen (Closure) | Hinweis: Die Finding-Klassen des Reviews (Kettenaussage „jede“ ohne eigene Ausnahme; Zustandswort im Satz weggelassen; Bezugspunkt nur an einer Stelle genau) gehören in den Lerneintrag; `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` kann hier „nicht eingetreten“ belegt werden |
| 7 Risiken mit Ausgang | offen (Closure) | Zwei Risiken in §6, beide „offen bis Closure“ |
| 8 Paarungen | offen (Closure) | |

## 3. Probe-Tabelle in §7 gegen die Verhaltenssätze

Ich bin jeden Satz des neuen Abschnitts und der geänderten Zeilen gegen die Tabelle gegangen (Regel aus `.claude/commands/implement-slice.md`: jeder Satz mit Verhaltensaussage hat eine eigene Zeile). Jeder Satz hat eine Zeile, einschließlich der Nacharbeit: erste Anfrage ohne Wartezeit bei `relative` (F-660, Zeile mit `r2`), Exit-Code 4 nach früherem Fehler (F-661, `err20`), Bezugspunkt Startup (F-662, Roh-Client), Log-Zeilen `play gestartet`/`play beendet` (F-663, `e1`). Kein Satz ohne Zeile gefunden. Zeilen mit Gegenprobe in beiden Richtungen: Startup-Zeit zählt nicht / Zeit nach Anmeldung zählt; `relative` / `absolute` bei Verspätung; `connect` / `first-request`.

Unschärfe, kein Befund: Der Satz *„Jeder dieser Aufrufe endet mit `PGR-E2001`“* über der Liste *Ungültige Kombinationen* umfasst auch den Punkt „Die Reihenfolge der Optionen ändert nichts“, der kein Aufruf ist (V-206).

## 4. AGENTS.md §3.11

| Forderung | Urteil |
|---|---|
| Handbuch beschreibt nur das Binary, keine Chronik, kein Zielstand | erfüllt (`grep`, eigene Proben) |
| keine Verweise auf Spezifikation, Lastenheft, ADRs, Slices, Wellen, Reviews | erfüllt (`grep`: 0 Treffer) |
| README ohne Chronik, Verweise auf `spec/` ohne Aussage über einen Stand | erfüllt |
| Beispiele als Datei gegen das Binary | laut §7 gelaufen (vier Blöcke wörtlich); nicht selbst wiederholt, weil Datenbank nötig |
| Sätze nur so weit, wie eine Probe reicht | erfüllt; die nicht übernommenen Aussagen stehen in §7 *Nicht in das Handbuch übernommen* |

## 5. Review-Findings F-660 bis F-664

| Finding | Behoben? | Beleg |
|---|---|---|
| F-660 „vor jeder Anfrage“ ohne Ausnahme | ja | Einleitung des Abschnitts, *Ergebnis* und README nennen die Ausnahme der ersten Anfrage bei `relative`; `grep` ohne Treffer; Probe `r2` in §7 |
| F-661 Exit-Code nach früherem Fehler | ja | Satz lautet „Der Exit-Code ist 0, nach einem früheren Fehler 4.“; Probe `err20` in §7 |
| F-662 Bezugspunkt | ja | Einleitung, README, Abschnittsanfang, Hinweis, Glossar nennen „Startup-Nachricht der Verbindung“; Probe in §7 |
| F-663 Gliederung | ja | *Ungültige Kombinationen* (5 Punkte) und *Fehlende Zeitangaben und Abbruch* als eigene Abschnitte; *Ergebnis* mit Log-Zeilen und Probe `e1`. Rest: Die Überschrift `#### Ungültige Kombinationen` folgt ohne Leerzeile auf die Liste (V-206) |
| F-664 INFO, zwei Hinweise an den Planner | offen, an den Planner | in §7 *Hinweise an Architect und Planner*; kein Befund gegen den Slice |

## 6. Plan gegen Diff (AGENTS.md §3.9)

- **Kopf:** `Bezug` LH-FA-21 stimmt; `Berührte Spec-Stellen: —` stimmt (kein Eintrag in `spec/` im Diff); `Verantwortlich` gesetzt.
- **§1 gegen Diff:** Abgrenzung verletzt, siehe V-204: Der Diff ändert `.claude/commands/implement-slice.md`, §1 sagt „Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan“ (die Roadmap-Zeile steht in §3). Produkt-Code, Tests, Gates, `docs/user/abdeckung-*.md` unverändert: stimmt.
- **§3:** Handbuch, README und Roadmap sind genannt; `implement-slice.md` fehlt (V-204).
- **§6:** Randformen und Risiken folgen dem Diff; kein Code, keine neue Randform durch den Diff entschieden.
- **§7 *Stand*:** zählt die geänderten Dateien auf und nennt `implement-slice.md` nicht (V-204). Die Aussage zu den Mutanten (`aa6eed4`: Zeilen entfernt, „kein neuer Wächter“) stimmt: Der Slice führt keinen Wächter ein. Die Befehlsfolge ist grün; sie wurde im Plan am Stand `15583ea` gemeldet, nach der Nacharbeit `8209d87` änderte sich das Handbuch; mein Lauf am Stand `236a39e` bestätigt das Ergebnis.
- **§8:** Zählung „2 Liefer-Punkte, 1 Schicht“ stimmt für Handbuch und README. Die Änderung am Command ist eine Harness-Änderung ohne Zeile in §8.
- **Roadmap:** Zeile *In Arbeit* nennt den Slice, Schritt 23, stimmt mit Schritt 23 in `welle-v1-abschluss` überein.
- **Commits:** reiner Move `56d7797` getrennt von Inhaltsänderungen (§3.3); Messages nennen Slice-Kennung und `LH-FA-21`, keine `SPEC-`/`ARC-`-Kennung.

## 7. Befunde

| ID | Schwere | Befund | Beleg |
|---|---|---|---|
| V-203 | LOW | DoD-Punkt 3 (`make gates` grün) hat in §7 keinen Beleg des Implementers (kein Lauf mit Stand und Ergebnis). Die Regel der Implementer-Rolle verlangt die Belege in §7. Mein Lauf am Stand `236a39e` ist grün, das Urteil bleibt „erfüllt“, der Beleg fehlt. | `grep -n 'make gates'` im Plan: nur Zeilen 52 (DoD) und 97 (§5) |
| V-204 | MEDIUM | Der Plan folgt dem Diff nicht (§3.9): `236a39e` ändert `.claude/commands/implement-slice.md` (Schritt 17: Satz zu den Mutanten der Befehlsfolge). Weder §1 (Abgrenzung nennt nur Handbuch, README, Plan), §3 (Tabelle), §7 *Stand* noch §8 (Nachzählung) kennen die Datei; der Satz trägt außerdem keinen Anker `(seit slice-…)`, wie die Nachbarsätze desselben Schritts. Es ist kein Produkt-Code, aber ein Harness-Artefakt, das jeder künftige Lauf liest; er ist der Lerneintrag, der vor der Closure schon verkörpert ist. Nicht blockierend für die Dokumentationsarbeit, aber vor der Closure im Plan nachzuziehen (Zeile in §3, Satz in §1 oder §7, Anker, Zählung in §8) oder in einen eigenen Vorgang zu geben. | `git show 236a39e`; `grep -n 'implement-slice' <Plan>`: kein Treffer |
| V-205 | INFO | Das Verdikt des Review-Reports verweist für das Gate-Ergebnis auf einen „Bericht des Reviewers“; der Report nennt es nicht. | Report, Abschnitt *Verdikt* und Zeile *Gates* |
| V-206 | INFO | Zwei Kleinigkeiten im Handbuch: (a) `#### Ungültige Kombinationen` folgt ohne Leerzeile auf die letzte Listenzeile von *Hinweise zum Einspielen* (Zeile 668; das Gate ist grün, die Darstellung ist in CommonMark in Ordnung); (b) die Einleitung der Liste „Jeder dieser Aufrufe …“ umfasst „Die Reihenfolge der Optionen ändert nichts“, das kein Aufruf ist. | `docs/user/benutzerhandbuch.md` Zeilen 666 bis 681 |

## Verdikt

**Der Slice erfüllt die Liefer-Punkte der DoD (1, 2) und, nach meinem Lauf, `make gates` (3).** Handbuch und README beschreiben das Binary im Ist-Zustand; keine Verweise auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews; Befehlsfolge grün; jeder Verhaltenssatz hat eine Zeile in der Probe-Tabelle; F-660 bis F-663 sind behoben. Keine HIGH-Befunde.

**Vor der Closure nachzuziehen:** V-204 (Plan folgt dem Diff nicht: `implement-slice.md` in §1, §3, §7, §8 und der Anker), V-203 (Gate-Beleg in §7). V-205 und V-206 sind Hinweise.

**Offen für die Closure (Implementer):** Notiz mit Lerneintrag (Finding-Klassen aus dem Review), Register, Ausgang der beiden Risiken aus §6, Paarungen, dann `git mv` nach `done/`. Die Closure habe ich nicht ausgeführt.
