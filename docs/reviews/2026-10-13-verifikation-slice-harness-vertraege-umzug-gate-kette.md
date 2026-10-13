# Verifikation: slice-harness-vertraege-umzug-gate-kette

**Rolle:** Verifier. **Datum:** 2026-10-13. **Spanne:** `8b1ede6` bis `491ef31` (HEAD). **Gegenstand:** [slice-harness-vertraege-umzug-gate-kette](../plan/planning/in-progress/slice-harness-vertraege-umzug-gate-kette.md) gegen DoD (§2), [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md), [MR-002](../../harness/conventions.md#mr-002) und Spezifikation. Review: [Review-Report](2026-10-13-review-slice-harness-vertraege-umzug-gate-kette.md).

Alle Belege unten sind in diesem Lauf selbst erzeugt, nicht aus §7 übernommen.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren) |
|---|---|---|
| Liefer-Punkt 1: Datei [`SPEC-052`](../../harness/vertraege/SPEC-052-gate-kette.md) wortgleich bis auf Punkt 12 und Punkt 22 (b); README führt Vertrag und fünf Historie-Zeilen | bestätigt | `git show 231089e:spec/spezifikation.md \| sed -n '3216,3477p' \| sed '1s/^### /# /;$d' \| diff - harness/vertraege/SPEC-052-gate-kette.md`: genau zwei Hunks, `120,122c120,122` (Punkt 12, beide Stellen im Wortlaut der Randform 2) und `247,249c247,250` (Punkt 22 (b), neuer Wortlaut wie Randform 1, der Satz „beide Aufrufarten gelten auch …“ entfällt). Datei 262 Zeilen; Zeile 3477 am Stand `231089e` ist die Leerzeile vor `## 12. Historie`, `$d` ist belegt. `#` statt `###` nur in der Überschrift. Historie: `diff <(tail -5 harness/vertraege/README.md) <(… sed -n '3576,3580p')` leer, Exit 0; zwei der fünf Zeilen tragen 2026-10-14. Zeile von `SPEC-052` in der Tabelle der Verträge (README Zeile 31). |
| Liefer-Punkt 2: Spezifikation | bestätigt | `git diff 8b1ede6 HEAD --stat -- spec`: 1 Einfügung (`## 11. Historie`), 277 Löschungen. `grep -c "SPEC-052\|^## 11\. Harness\|^## 12\. Historie" spec/spezifikation.md` = 0 (am Stand `8b1ede6`: 8). `## 11. Historie` ist die letzte Überschrift; `grep -n "§ \?12\|Harness-Werkzeug\|SPEC-052\|harness/" spec/*.md` ohne Treffer. `make docs-check`: 719 Dateien, 0 Befunde. |
| Liefer-Punkt 3: Verweise | bestätigt | DoD-`grep` (ohne `.git`, `.harness`): Treffer nur in Zeitdokumenten (Reviews, Pläne in `done/`, Register-Evidenz), in [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), im Plan selbst (Zitate), in `MR-002` (Zustandsabsatz und Geltungsbereich, Closure-Arbeit nach §5) und in `open/slice-harness-lint-warten-ohne-frist.md:26,168` (dort ist `spezifikation.md` §11 die umnummerierte *Historie*, richtig). `grep -rn "§ \?12" docs/plan/planning/open docs/plan/planning/next`: ohne Treffer. `grep -n "§11\|§ 11" AGENTS.md`: keine Zeile. `AGENTS.md` §3.12 (Zeilen 300 bis 312): nennt `harness/vertraege/` ohne Ausnahme und ohne `SPEC-052`-Sonderfall. `.claude`: nur der Kommentar des Hooks mit dem neuen Pfad. Vier Pläne (`lh-links-pflicht`, `mutation`, `coverage`, `meldungskatalog-gate`): Ort `harness/vertraege/`, Nummer 053. Acht Kommentardateien und `harness/README.md` im Diff. |
| `make gates` grün | bestätigt | Eigener Lauf auf HEAD (`491ef31`, sauberer Baum): Exit 0, 7 min 29 s; `lint-gegenprobe: gruen …`, `stop-hook-gegenprobe: gruen — 139 Fälle der Punkte 1, 9 bis 22 je Mutant rot, unveränderte Kopie grün, Anschluss an GATE_CHECKS`. `make docs-check` 0 Befunde, `make kopf-check` Exit 0. |
| Review durchgeführt | bestätigt | Report liegt vor (0 HIGH, 0 MEDIUM, 3 LOW, 3 INFO). |
| Handbuch und README unberührt | bestätigt | `git diff 8b1ede6..HEAD --name-only`: weder `docs/user/` noch Wurzel-`README.md`, nichts unter `internal/`, `cmd/`, `test/`. |
| Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge (vier §6-Risiken stehen auf „offen bis Closure“), drei Paarungen | offen (erwartet) | Closure-Schritte vor dem `git mv`; §7 trägt bisher nur Belege des Implementers. |

## Bewusstes Brechen (selbst nachgefahren, danach `git checkout`, Baum sauber)

| Zusage | Mutation | Ergebnis |
|---|---|---|
| Punkt 12: Änderung an der Datei des Vertrags verlangt den vollen Lauf (Fall `p12-selbst-vertrag`) | `harness/vertraege/*) return 0 ;;` in die Positivliste von `.claude/hooks/stop-require-gates.sh` (Zeile 101), `touch`, volle `make stop-hook-gegenprobe` (2 min 15 s) | `ROT — Kopie 'ok' ist nicht grün: p12-selbst-vertrag`, nur dieser Fall: aus dem richtigen Grund rot |
| wortgleich | `X` vor Zeile 3 von `SPEC-052-gate-kette.md` | Vergleich meldet 3 statt 2 Hunks |
| Punkt 22 (b) sagt nicht das Gegenteil | alter Satz „beide Aufrufarten gelten auch für die Freigabe nach Punkt 9“ eingefügt | `grep -c` meldet 1 statt 0 (grep-Beleg, kein Gate; so benannt) |
| Link auf den Vertrag zeigt auf den neuen Ort | `harness/README.md:78` zurück auf `../spec/spezifikation.md#spec-052…` | `docs-check`: 1 Befund, `anchor-missing` |
| Spezifikation verweist nicht auf `harness/` | Link `../harness/vertraege/README.md` am Ende von `spec/spezifikation.md` | `docs-check`: `matrix-forbidden — spec-straten → aussen`, 1 Befund |
| `AGENTS.md` §3.12 ohne Ausnahme | `spec/spezifikation.md §11` an Zeile 305 angehängt | `grep -n "§11" AGENTS.md` meldet Zeile 305 (grep-Beleg) |
| Historie der Spezifikation heißt `## 11. Historie` | `## 12. Historie` | `grep -c "^## 12. Historie"` meldet 1 statt 0 (grep-Beleg) |
| fünf Historie-Zeilen im README | letzte Zeile gelöscht | `git diff --numstat`: 0/1, die `diff`-Prüfung gegen 3576 bis 3580 wäre rot |

Der Fall `p12-selbst-vertrag` ist damit nicht nur im grünen Lauf, sondern nachweislich durch genau die Mutation rot, die seine Zusage bricht. Die Prosa-Zusagen (Kommentar-Pfade, Spezifikation ohne `SPEC-052`, §3.12) tragen nur über `grep`; das deklariert §7 selbst.

## Review-Findings

- **F-1 (LOW, Kommentare der Gegenprobe):** behoben in `c85257c`; Kopfkommentar nennt „Hook, Spezifikation und Vertrag nicht inert“, der Kommentar über `f_p12_selbst_hook` nennt den Pfad des Vertrags. Nur Kommentare im Diff (3 Einfügungen, 2 Löschungen), grüner Lauf bestätigt.
- **F-2 (LOW, Zahlen 261/3476):** behoben in `491ef31`. Der Vergleichsbefehl in §7 lautet jetzt `sed -n '3216,3477p' | sed '1s/^### /# /;$d'` und liefert am Stand `231089e` genau die zwei Hunks; Zeilenzahl 262 überall. Ich habe den Befehl wortgleich ausgeführt.
- **F-3 (LOW, `slice-harness-lh-links-pflicht` Zeile 343 „Spezifikation“, uneinheitliche Zählung):** offen. Bewertung:
  - **Kein Blocker.** Die Sendung dieses Slice (Ort neuer Kennungen mit Kennung des Gebers in Kopf, §1, §3, §6) ist getragen. Die Zählung des Nehmers (zwei Liefer-Punkte, zwei Schichten) bleibt in jeder Lesart unter der Grenze: Zeile 343 ist ein datierter Fortschreibungs-Eintrag („Nachgezählt nach der Entscheidung des Nutzers vom 2026-10-10“), und „Spezifikation“ war zu dem Zeitpunkt die richtige Schicht; Zeile 329 trägt jetzt „Vertrag“. Eine Änderung hier verfälschte die Chronik der Zählung, eine Nichtänderung lässt nur das Etikett uneinheitlich. Die DoD dieses Slice sagt nichts zur Zählweise des Nehmers.
  - **Closure-Ausgang des Planners.** Der Nehmer liegt in `open/` und wird vor `open → next` ohnehin neu gezählt. Die Closure-Notiz nennt den Hinweis (Implementer-Hinweis 2 in §7) mit Ausgang *weiter offen* (Register `BEO-REPO/slice-waechst-durch-uebernahmen` bzw. `BEO-REPO/plan-folgt-korrektur-nicht`, Evidence-Datei) oder *entfallen* mit der Begründung, dass die Stelle nach `harness/vertraege/` in der Schicht Harness liegt und der Nehmer sie schon zählt. Eine neue Kennung ist nicht nötig; der Planner entscheidet beim `open → next`.

## Plan-vs-Code-Diff

Alle Zeilen aus §3 sind im Diff (20 Dateien, `git diff 8b1ede6 HEAD --name-status`): Vertrag (neu), README des Ordners, Spezifikation, `harness/README.md`, `AGENTS.md`, `.claude/hooks/stop-require-gates.sh`, vier Skripte unter `tools/harness/`, zwei Fragmente `harness/mk/`, vier Pläne in `open/` und `next/`, `slice-harness-lint-warten-ohne-frist` (Randform 3), Plan selbst, Review-Report. Nicht im Diff, wie vorgesehen: `roadmap.md` (Stand erst bei Closure), `MR-002` (Closure, §5), `.claude/agents/architect.md` (Randform 7), `docs/reviews` und `done/` (kein Link brach, `docs-check` 0 Befunde), ADRs (kein Diff, §3.5), Handbuch, README, Produkt-Code. Keine Renames oder Löschungen. Code-Änderung außer Kommentaren: nur der Fall `p12-selbst-vertrag` samt `aufnahme`-Mutation (Randform 2, Entscheidung des Nutzers). Größe in der Metrik von `AGENTS.md` §3.14 (Pfade des Plans, `8b1ede6..491ef31`): 442 + 324 = **766**; Marke 900 und Meldeschwelle 800 nicht erreicht; Schätzung des Architect 680 bis 780 mit Fall: am oberen Rand (zwei Review-Runden und Report). Dieser Report kommt dazu und bleibt voraussichtlich unter 800.

## Befunde

- **HIGH/MEDIUM:** keine.
- **LOW (L-1):** F-3 des Reviews offen (siehe oben); Ausgang gehört in die Closure-Notiz.
- **INFO (I-1):** §7 nennt die Größe 636 (bis `ebccbeb`); die Closure kann 766 nennen und die Review-Runde (rund 20 Prozent) dazu.
- **INFO (I-2):** Offen für die Closure, wie §5 vorsieht: `MR-002` (Zustandsabsatz und Geltungsbereich), Lerneintrag, Register, Ausgänge der vier §6-Risiken (Punkt 22 (b): Zusage unverändert, Gegenprobe-Fälle `p22-*` grün, *entfallen*; fünfter Plan mit `§11`: `grep` fand nur die bekannten Treffer, *entfallen* mit Begründung; Punkt 12 und Fall: gelöst durch `p12-selbst-vertrag`; Größe 766 unter Marke).

## Fazit

Die drei Liefer-Punkte, `make gates` (Exit 0, volle Kette auf HEAD), `make docs-check`, `make kopf-check`, Review und die Unberührtheit von Handbuch und README sind bestätigt; die Mutation zu `p12-selbst-vertrag` und die Stichproben sind aus dem richtigen Grund rot. F-1 und F-2 sind behoben, F-3 blockiert nicht und ist ein Closure-Ausgang des Planners. Freigabe zur Closure.
