# Review-Report: slice-version-option — 2026-10-22

**Review-Art:** Code — Code-Review gegen Plan, Spezifikation und Hard Rules (`AGENTS.md`).

**Gegenstand:** Zweig `impl-version-option`, Diff `49da54b6993287b9d1902b30aacef66e8afc4703..1935d6b5` (Commits `61279514`, `ea6d1524`, `5804f486`, `1935d6b5`).

**Skill:** `.harness/skills/reviewer.md` @ `1935d6b5` (Fassung nach `slice-harness-regeln-rueckbau`)
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-22

> **Zitier-Form.** Kennung, nicht Adresse: Slices als `slice-<Kennung>`, Sensoren als `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code.

**Eingangs-Kontext:**

- Slice-Plan `slice-version-option` (§1 bis §8)
- Spezifikation `LH-FA-01.a` (Absätze *Hilfe vor Prüfung*, *Ende der Optionen*, *Version*), `LH-FA-14.a`, `SPEC-014`
- `LH-FA-01`, `LH-FA-14`
- `AGENTS.md` §3.9 bis §3.13 (Fassung im Worktree), §3.14; keine ADR im Diff genannt
- `v6.18.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice

Gefahren: gebautes Binary (Stufe `dist`, `VERSION=9.9.9-rv`) und Binary der Basis, Differenzlauf über 16276 Argumentfolgen (Länge 0 bis 3, 29 Symbole): genau zwei Abweichungen, `--version` und `-v` (die globale Hilfe um den neuen Block bereinigt). 14 Mutationen von `cli.go` in Kopien (`--version` ohne `-v`, ohne `--version`, weiteres Argument zugelassen, `-V`, Präfix, Block der Hilfe entfernt, Option in der Hilfe von `record`, Name geändert, Prüfung vor `len(args)==0` entfernt, Umgebung/Datei/Standarddatei gelesen): alle rot außer einer äquivalenten (Prüfung vor der Hilfe-Prüfung: `versionsAngabe` verlangt genau ein Argument, das keine Hilfe-Angabe sein kann). Die vier neuen E2E-Tests grün gegen das neue Binary, `TestE2EVersionsAngabe` und `TestE2EVersionHilfe` rot gegen das Basis-Binary. `make docs-check`, `make kopf-check`, `make abdeckung-check`, `make meldungskatalog-check`, `make a-check` je Exit 0.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | §7 *Abweichungen* nennt als Abweichungen vom Plan „E2E oder `bootstrap_test.go`“ und „drei Stellen im Handbuch statt des einen Absatzes“; §3 des Plans führt beide Dateien und alle drei Handbuchstellen bereits. Die Abweichungen beschreiben eine frühere Fassung des Plans. | `AGENTS.md` §3.9 (ohne Failure-Szenario, Skill: LOW) | `docs/plan/planning/in-progress/slice-version-option.md` · „geliefert sind beide“ | nein — Lesen von §3 gegen §7 | Plan-Abweichung beschreibt überholten Planstand |
| F-2 | INFO | Der Handbuch-Kopf nennt Version 0.11 und Stand 22.10.2026, aber Software-Version und Gültigkeitsbereich weiter `0.1.1`; das Release 0.1.1 kennt `--version` nicht. Gelebte Praxis: Das Dokument zählt je Änderung, Software-Version und Gültigkeitsbereich ziehen erst nach dem Tag nach (`1f501a6c`, `9879253a`; `0add50ae` und die Cancel-Änderungen beschrieben Verhalten vor dem Tag 0.1.1). Kein neues Muster. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „Software-Version: `0.1.1`“ | nein | Handbuch-Kopf eilt der Software-Version voraus |
| F-3 | INFO | Drei Tests decken die Zusage „liest weder Umgebungsvariable noch Konfigurationsdatei“ (`TestParseVersion`, `TestRunVersionsAngabe`, `TestE2EVersionsAngabe`). Jeder der Mutanten „liest Umgebung“, „liest Datei“, „liest Standarddatei“ färbt die ersten beiden rot; der E2E-Test deckt zusätzlich die Verdrahtung im gebauten Binary und die Zeile mit echter Version, der Bootstrap-Test die eingespeiste Version `1.2.3`. Leichte Überdeckung, kein Mangel. | Maintainability | `internal/bootstrap/bootstrap_test.go` · „TestRunVersionsAngabe“ | ja — Mutationen in Kopie | Überdeckung zweier Teststufen |
| F-4 | INFO | `TestParseKeineVersionsAngabe` führt `--version=` und `-v=` (leerer Wert) über die Spezifikation hinaus; sie fallen unter „ein Argument mit `=`“ der Spezifikation und sind belegt (Binary: `unbekanntes Kommando`, Exit 2). Kein Verhalten ohne Spezifikation. | `LH-FA-01.a` | `internal/adapters/driving/cli/cli_test.go` · „"--version=", "-v="“ | ja | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driving/cli/cli.go` (`versionsAngabe`, `parse`, `usage`) | geprüft, ohne Befund: jede Zeile des Absatzes *Version* abgedeckt (kein, ein, zwei Argumente, `--` davor, leeres Argument, Wert einer Option, Groß-/Kleinschreibung, `-vv`, `-vh`, `--help --version`); Ausgabe über dasselbe Kommando `version`, keine zweite Quelle der Zeile |
| `-v` in `internal/adapters/driving/cli`, `cmd/` | geprüft, ohne Befund: sonst nirgends als Option vergeben |
| `internal/bootstrap`, Core-Reinheit, `.a-check.yml` | geprüft, ohne Befund: keine neue Abhängigkeit, `make a-check` Exit 0 (Hinweis auf 44 ungeschichtete Dateien ist Bestand) |
| Tests `cli_test.go`, `version_e2e_test.go` | geprüft, ohne Befund: je Zusage und je Element ein Fall, keine Tautologie, Fehlerfälle prüfen Meldung und Exit 2; Mutationen rot (siehe oben) |
| `docs/user/benutzerhandbuch.md`, Hilfetext | geprüft, ohne Befund: jeder Satz gegen das gebaute Binary bestätigt, keine Chronik, kein Verweis auf Spezifikation oder Slices; Hilfetext stimmt mit Handbuch überein |
| `docs/user/abdeckung-*.md`, `spec/spezifikation.md` | geprüft, ohne Befund: Tabellen gleich den Deklarationen, Code folgt der Spezifikation |
| Plan §1, §3, §6 gegen Diff | geprüft; abweichend nur F-1 |
| Größe (`AGENTS.md` §3.14) | gemessen: 273 Zeilen in den Pfaden des Plans, 295 mit `internal/bootstrap/bootstrap_test.go`, Marke 1500 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Plan-Abweichung beschreibt überholten Planstand · Handbuch-Kopf eilt der Software-Version voraus · Überdeckung zweier Teststufen

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Nicht geprüft: die volle Integrationssuite und `make gates` (nur die vier neuen E2E-Tests mit dem gebauten Binary, ohne PostgreSQL); macOS und Windows (Grenze der Spezifikation).
