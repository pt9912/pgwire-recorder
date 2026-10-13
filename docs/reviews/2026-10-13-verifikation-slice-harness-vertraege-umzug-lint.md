# Verifikation: slice-harness-vertraege-umzug-lint

**Rolle:** Verifier. **Datum:** 2026-10-13. **Spanne:** `f4b6d6d` bis `1cdec2c` (HEAD). **Gegenstand:** [slice-harness-vertraege-umzug-lint](../plan/planning/in-progress/slice-harness-vertraege-umzug-lint.md) gegen DoD (§2), [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), [MR-002](../../harness/conventions.md#mr-002) und Spezifikation. Review: [Review-Report](2026-10-13-review-slice-harness-vertraege-umzug-lint.md).

Alle Belege unten sind in diesem Lauf selbst erzeugt, nicht aus §7 übernommen.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren) |
|---|---|---|
| Liefer-Punkt 1: Datei [`SPEC-049`](../../harness/vertraege/SPEC-049-lint-profil.md), wortgleich, README führt Vertrag und neun Historie-Zeilen | bestätigt | `git show 231089e:spec/spezifikation.md \| sed -n '3012,3163p' \| sed '1s/^### /# /' \| diff - harness/vertraege/SPEC-049-lint-profil.md`: leer, Exit 0; Datei hat 152 Zeilen; Zeile 3164 am Stand `231089e` ist leer (nur Zeilenende). Zeilen 3514 bis 3522: jede einzeln mit `grep -cxF` im `README.md` des Ordners gefunden (je 1, neun von neun). Der Text nennt keinen Pfad (Randform 3). |
| Liefer-Punkt 2: Spezifikation | bestätigt | `git diff 231089e HEAD --stat -- spec`: 347 Löschungen, keine Einfügung; Hunks: `2880,336` (Ordner-Slice und dieser Vertrag), `3514,9`, `3524,2`; `grep -c "harness/vertraege" spec/spezifikation.md` = 0, kein Link auf `harness/`; einzige verbliebene Nennung von `SPEC-049` ist die Historie-Zeile vom 2026-10-11 (mehrere Änderungen, nicht Teil der neun, Randform 2). `make docs-check`: 715 Dateien, 0 Befunde. |
| Liefer-Punkt 3: Verweise | bestätigt | `grep -rn "spezifikation.md#spec-049" --exclude-dir=.git --exclude-dir=.harness .`: Treffer nur Text (Beispiel-Anker in `SPEC-050`, wortgleich; Ordner-Slice in `done/`; dieser Plan; Review-Report), kein Link. `grep -rn "spezifikation.md SPEC-049" --exclude-dir=.git .`: nur Plan und Review-Report (Zitat der Greps). Kommentare: `.golangci.yml` Zeile 2 (eine Zeile), `lint.sh` Zeile 5, `lint-gegenprobe.sh` Zeile 4 nennen `harness/vertraege/SPEC-049-lint-profil.md`; der Diff dieser drei Dateien enthält nur diese drei Kommentarzeilen. `AGENTS.md` §3.12: `grep -n "SPEC-049" AGENTS.md` ohne Treffer; nennt `SPEC-052` mit §11. |
| `make gates` grün | bestätigt | Eigener Lauf auf HEAD: Exit 0, `lint-gegenprobe: gruen — Grundlauf ohne Befund, jeder Fall je Punkt von SPEC-049 wie erwartet`, `stop-hook-gegenprobe: gruen — 138 Fälle`. `make kopf-check`: Exit 0, keine Ausgabe. Arbeitsbaum danach sauber. |
| Review durchgeführt | bestätigt | Report liegt vor (0 HIGH, 0 MEDIUM, 2 LOW, 3 INFO); F-2 behoben, F-1 siehe unten. |
| Handbuch und README unberührt | bestätigt | `git diff f4b6d6d..HEAD --name-only`: weder `docs/user/` noch Wurzel-`README.md`, kein Code unter `internal/`, `cmd/`, `test/`. |
| Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge (drei §6-Risiken stehen auf „offen bis Closure“), drei Paarungen | offen (erwartet) | Closure-Schritte vor dem `git mv`; §7 trägt bisher nur Belege des Implementers. |

## Bewusstes Brechen (selbst nachgefahren, danach `git checkout`, Baum sauber)

| Zusage | Mutation | Ergebnis |
|---|---|---|
| wortgleich | `X` vor Zeile 3 von `SPEC-049-lint-profil.md` | Vergleich meldet `3c3` mit beiden Zeilen: rot aus dem richtigen Grund |
| Link auf den Vertrag zeigt auf den neuen Ort | `harness/sensors/lint.md:8` zurück auf `../../spec/spezifikation.md#spec-049…` | `docs-check`: 1 Befund, `anchor-missing` (deckt sich mit §7) |
| Spezifikation verweist nicht auf `harness/` | Link `../harness/vertraege/README.md` am Ende von `spec/spezifikation.md` | `docs-check`: `matrix-forbidden — spec-straten → aussen`, 1 Befund |
| Spezifikation führt `SPEC-049` nicht mehr | Zeile `### SPEC-049 — Lint-Profil (`lint`)` angehängt | `grep -c "^### SPEC-049"` meldet 1 statt 0 (grep-Beleg, kein Gate; so benannt) |
| `AGENTS.md` §3.12 nennt `SPEC-049` nicht mehr | „für die Verträge `SPEC-049` und“ in Zeile 305 eingesetzt | `grep -n "SPEC-049" AGENTS.md` meldet Zeile 305 (grep-Beleg) |
| neun Historie-Zeilen im README | letzte Zeile (2026-10-08) gelöscht | Summe der `grep -cxF`-Treffer 8 statt 9 |

## Review-Findings

- **F-2 (LOW, Zahlen 153/3164/17):** behoben. Commit `1cdec2c` zieht 152 Zeilen, 3012 bis 3163 und 18 Dateien nach; `grep` in Plan zeigt §1 und DoD mit 152 und 3163. (Kontrolle: Vertragsdatei hat 152 Zeilen.)
- **F-1 (LOW, Nehmer [slice-harness-lint-warten-ohne-frist](../plan/planning/open/slice-harness-lint-warten-ohne-frist.md) ohne Nachzählung in §8):** nicht behoben, §8 des Nehmers enthält weiterhin keine Zählung (`grep -n "Nachgez\|Schicht" ` trifft nur Fremdstellen). Bewertung:
  - **Kein Blocker für diesen Slice.** Die Sendung (Ort des Vertrags mit Kennung des Gebers in Kopf, §1, §3, §6) ist getragen. Der Ortswechsel `spec/` nach `harness/` fügt dem Nehmer keinen Liefer-Punkt und keine Schicht hinzu; die Überschreitung (nach Zählung des Planners 3 Liefer-Punkte und 3 Schichten, nach `AGENTS.md` §3.13 zählen Tests der Hexagon-Schichten mit) bestand schon vor dieser Sendung: der Bestand-Punkt (Tests in Kern und Driving-Adapter) kam mit dem Plan. Die Verifikation misst diesen Slice gegen seine DoD; seine DoD sagt nichts zur Größe des Nehmers.
  - **Closure-Ausgang, nicht stilles Vergessen.** `AGENTS.md` §3.13 sagt „läge der Nehmer darüber, geht der Punkt an den Planner“. Der Nehmer liegt in `open/` und muss vor seinem `open → next` ohnehin geschnitten werden. Der Ausgang gehört in die Closure-Notiz: ein Risiko-Ausgang *eingetreten* oder ein Hinweis an den Planner mit Kennung des Folge-Schnitts, oder (Register) *weiter offen* mit `BEO-REPO/slice-waechst-durch-uebernahmen`. Eine Kennung für den Folge-Schnitt gibt es noch nicht; sie vergibt der Planner. Bis dahin ist der Hinweis in der Closure-Notiz das Mindeste (Empfehlung: im selben Commit ein §8-Absatz im Nehmer, der die Zählung des Planners nennt, damit die Überschreitung am Plan sichtbar ist).

## Plan-vs-Code-Diff

Alle Zeilen aus §3 sind im Diff (27 Dateien: Vertrag, README, Spezifikation, `harness/README.md`, `harness/sensors/lint.md`, `AGENTS.md`, drei Kommentardateien, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) nur Link-Pfade, drei Pläne in `open/`, acht in `done/`, vier Berichte, dieser Plan, Review-Report; kein Plan in `next/`). Keine Datei außerhalb von §3 und der Pläne; `.d-check.yml`, `internal/`, `harness/mk/lint.mk`, Roadmap unberührt. Größe in der Metrik von `AGENTS.md` §3.14 (Pfade des Plans, `f4b6d6d..HEAD`): 316 + 206 = **522**, Marke 600 nicht erreicht, über der Schätzung 370 bis 440 (Review-Report, Berichtigung und §7-Belege kamen hinzu; ohne Report 438 laut Review bis `9a5c80d`).

## Befunde

- **HIGH/MEDIUM:** keine.
- **LOW (L-1):** F-1 des Reviews offen (siehe oben); Ausgang gehört in die Closure.
- **INFO (I-1):** Größe in §7 (393) steht vor Review und Berichtigung; die Closure kann 522 nennen.
- **INFO (I-2):** Die Zusagen „Spezifikation ohne `SPEC-049`“ und „§3.12 ohne `SPEC-049`“ tragen nur über `grep`, kein Gate; so deklariert.

## Fazit

Die drei Liefer-Punkte, `make gates` (Exit 0), `make docs-check` und `make kopf-check`, Review und die Unberührtheit von Handbuch und README sind bestätigt; alle sechs Mutationen sind aus dem richtigen Grund rot. F-1 blockiert nicht, braucht aber einen benannten Closure-Ausgang (Folge-Schnitt durch den Planner). Freigabe zur Closure.
