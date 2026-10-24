# Verifikationsbericht: slice-test-unter-last-ursache — 2026-10-24

**Rolle:** Verifier (Frage: Bauen wir es richtig? — gegen Plan, DoD, bestätigten Wortlaut und `LH-FA-20.a`)
**Gegenstand:** Zweig impl-last1, Stand 136499c4, Basis ffdacf7d (§4)
**Modell:** claude-sonnet-5-5
**Eingang:** Plan §7 *Belege des Implementers*, Review-Report `2026-10-24-review-slice-test-unter-last-ursache.md`
**Randbezug:** [ADR-0055](../plan/adr/0055-gegenproben-haben-ein-zeitbudget.md), [ADR-0056](../plan/adr/0056-groessen-metrik-misst-den-gelieferten-diff.md)

## 1. DoD-Konformität

| DoD-Punkt | Stand im Plan | Befund des Verifiers |
|---|---|---|
| 1 Ursache der Go-Tests (a) Last wörtlich (b) Läufe (c) Proben (d) Ursache | [x] | erfüllt mit Einschränkung, siehe unten |
| 2 Signalfall `p17-budget-sigint` | [x] | erfüllt: „nicht reproduziert in 8 Läufen bei Last 16 bis 32 auf 16 Kernen“; Abweichung von Randform 3 (8 statt 15) in §7 und Plan benannt |
| 3 Zwischenregel in Schritt 16 | [ ] (Kästchen offen) | Inhalt geliefert (Commit 00bf41a6); Kästchen sollte der Implementer/Planner bei Closure setzen |
| 4 kein Code/Frist/Test geändert | [x] | bestätigt: `git diff --numstat ffdacf7d..HEAD -- internal test tools harness` leer |
| `make gates` grün | [ ] | **kein Beleg in §7**; nicht gefahren (Auftrag: nicht nötig). Bleibt offen für die Closure, der Stop-Hook verlangt ihn wegen `.claude/` |
| Review | [ ] | Report liegt vor (0 HIGH/MEDIUM, 5 LOW, 2 INFO); F-1 bis F-6 nachgezogen, siehe 4 |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | [ ] | Closure-Pflichten, noch nicht Gegenstand dieser Verifikation |

**DoD 1 im Einzelnen.**
- (a) Der Aufruf der Last steht in §7 mit Erzeuger (eine Rechenschleife je Container), Anzahl (K = 32 und 48), Dauer (bis Serienende), Messung (`/proc/loadavg`). Er ist **nicht wörtlich wiederholbar**: Platzhalter (`<serie>`, `K`, `<netz>`, `<pg>`), Bereiche statt Werte vor und nach je Serie, PostgreSQL-Start nur „wie `make test-integration`“; Messhost abgeschaltet, Skripte im Scratchpad. §7 sagt das selbst („Grenze der Wiederholbarkeit“, Review F-5). Die DoD verlangt „wörtlich und wiederholbar“: nur als Beschreibung erfüllt. Ehrlich ausgewiesen, daher kein Verstoß gegen §3.11, aber eine **Teilerfüllung**, die der Planner bei Closure sehen muss.
- (b) Tabelle je Test mit allein/K=32/K=48, Läufen und Rotzahl vorhanden; Zählung in (d) (je 30 Läufe unter Last für A und B, die 30 des ersten Versuchs zählen nicht) ist mit (a) und (b) konsistent (Review F-2 behoben).
- (c) P1 bis P4 mit „trennt“-Spalte und Ergebnis vorhanden.
- (d) C: Wettlauf im Code; A und B: „nicht reproduziert in je 30 Läufen unter Last 18 bis 48“, Ursache offen. Diese Form ist in DoD 1 (d) und Randform 3 ausdrücklich zugelassen.

## 2. Plausibilität der Ursache (selbst nachgefahren)

Code gelesen: `einspielen.go` `Anfrage` hält `schreiben` bis nach `Flush` (`defer Unlock`); `schliesse` sendet `Terminate` nur bei `s.schreiben.TryLock()`. `play_test.go` `TestRunPlayZweitesSignal`: wartet auf `Q`, löst beide Signale, wartet auf `X` (Frist 30 s). Der Mechanismus ist aus dem Code ableitbar.

Nachgefahren in einer Kopie außerhalb des Repos (Scratchpad `…/scratchpad/p`, Arbeitsbaum des Repos unberührt, Image `pgwire-recorder:test`, `--network none`, kein anderer Docker-Lauf, ohne künstliche Last):
- **P1** (300 ms `time.Sleep` nach `Flush` in `Anfrage`): 5 Läufe, **5 rot**, Meldung `play_test.go:341: binnen 30 s keine Client-Nachricht 'X', bisher ""` — wie in §7.
- **P2** (zusätzlich `TryLock` durch `Lock` ersetzt; meine Variante ruft `Lock()` bedingungslos bei `terminate`): 5 Läufe, **0 rot** — wie in §7.

Damit trennt die Probe Wettlauf von Frist belastbar; das Rot ist ohne Last herstellbar. Nicht nachgefahren: die Lastserien (K = 32/48) und A/B (Zeit- und Hostbudget; §7 ist dort Behauptung mit Zahlen, nicht von mir bestätigt). Meine Läufe sagen nichts über A und B.

Lesart der Zusage: `LH-FA-20.a` auf main (36b6387c) schärft „zweites Signal sendet Terminate nach vollständig geschriebener Nachricht“; der Code bricht das, der Befund „Wettlauf im Code“ deckt sich. Behebung ist ausdrücklich nicht Gegenstand (§1) und braucht einen Folge-Slice, den der Planner anlegt (Risiko-Ausgang *eingetreten*, noch ohne Kennung — Closure-Aufgabe).

## 3. Zwischenregel gegen bestätigten Wortlaut

`implement-slice.md` Schritt 16 gegen den Wortlaut in Plan §3 (Nutzer 2026-10-24, „Vorschlag 2 (Recommended)“): Wortlaut von „Endet `make gates` rot …“ bis „… bis ein Folge-Slice es behebt.“ **zeichengleich**. Zusätze: Einleitung „Zwischenregel (seit slice-test-unter-last-ursache, Entscheidung des Nutzers vom 2026-10-24)“ und der Schlusssatz „Die Regel gilt für den Implementer, nicht für Reviewer und Verifier.“ Beide sind durch Plan §1 und §6 Randform 7 gedeckt (Geltung nur für den Implementer), stehen aber nicht in dem Text, den der Nutzer bestätigt hat; Hinweis, keine Abweichung der Regel. Wiederholung nur für die zwei Unterfälle A/B; `TestRunPlayZweitesSignal` ausgenommen und als Befund geführt, wie bestätigt. Der Klammersatz „(Frist von 30 s unter CPU-Last)“ ist im Wortlaut bestätigt; §7 hält ausdrücklich fest, dass Last allein das Rot für A/B nicht erzeugte (Review F-3, Reichweite benannt).

## 4. Plan-vs-Diff (AGENTS.md §3.9)

Diff Basis..HEAD: 3 Dateien — `implement-slice.md` (+1), Plan, Review-Report. Kein Code, Test, Gate. §1 (Ziel 3 und Wortlaut-Absatz), §3 (Wortlaut), §6 (Randformen 3, 7, 8), §7 und §8 folgen dem Stand. Review-Findings: F-1 (Größe 1 genannt), F-2 (Laufzahlen), F-3 (Reichweite), F-4 (8 Läufe benannt), F-5 (Grenze der Wiederholbarkeit), F-6 (Satz „entscheide ich nicht“ entfernt, Lesart nach 36b6387c) — alle in §7/Plan behoben, bis auf F-5, das **dokumentiert statt behoben** ist (nicht behebbar: Host abgeschaltet). F-7 INFO ohne Handlung. Kopf: `Berührte Spec-Stellen: —` passt zu §1 (`LH-FA-20.a` wird nur als Hinweis genannt, nicht berührt); `make kopf-check` grün.

## 5. Größe (Marke 1500, AGENTS.md §3.14)

`git diff --numstat --find-renames=10% ffdacf7d..HEAD -- .claude internal test tools harness ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' | awk …` → **1**. Mit §7 übereinstimmend; Marke nicht erreicht.

## 6. Gates

- `make docs-check`: 1475 Dateien, 0 Befunde (am Stand 136499c4 vor Anlage dieses Berichts).
- `make kopf-check`: Exit 0.
- `make gates`, Gegenproben: nicht gefahren (Auftrag). Der DoD-Punkt „`make gates` grün“ bleibt unbelegt.

## 7. Offene Punkte für Planner und Closure

1. DoD „`make gates` grün“ ohne Beleg; vor `done/` fahren.
2. DoD 1 (a) nur als Beschreibung erfüllt (nicht wörtlich wiederholbar) — Planner entscheidet, ob das als erfüllt gilt oder als benannte Lücke in die Closure geht.
3. Folge-Slice für den Wettlauf in `schliesse`/`Anfrage` ohne Kennung (Risiko-Ausgang *eingetreten*), vom Planner anzulegen; der Test C bleibt bis dahin rot-reproduzierbar.
4. A/B: Ursache offen (Risiko *weiter offen* → `BEO-REPO/test-faellt-unter-cpu-last-aus`); Register steht noch auf `geplant`, Umstellung bei Closure.
5. Schlusssatz der Regelzeile (Geltung nur Implementer) ist gegenüber dem bestätigten Text ein Zusatz; Nutzer kann ihn bestätigen oder streichen lassen.

## Verdikt

**Bestanden mit Vorbehalten.** Die Liefer-Punkte 1 bis 3 sind belegt, die Ursache von C ist von mir unabhängig bestätigt (P1 5/5 rot, P2 0/5 rot), die Zwischenregel entspricht dem bestätigten Wortlaut (plus Geltungssatz), Marke nicht berührt, Review-Findings abgearbeitet. Vorbehalte: Punkt 1 und 2 der Liste in Abschnitt 7 (kein Gate-Beleg, Last nicht wörtlich wiederholbar). Nicht nachgefahren: A/B und die Lastserien.
