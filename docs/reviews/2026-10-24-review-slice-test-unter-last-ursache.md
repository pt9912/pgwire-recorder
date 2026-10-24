# Review-Report: slice-test-unter-last-ursache — 2026-10-24

**Review-Art:** Code (Diff aus Messbelegen in §7 und einer Regelzeile gegen Plan und Konventionen; kein Produkt-Code im Diff)

**Gegenstand:** slice-test-unter-last-ursache, Zweig impl-last1, Commits 89a13500 (Messung, §7) und 00bf41a6 (Zwischenregel), Basis ffdacf7d

**Skill:** `.harness/skills/reviewer.md` @ Stand des Zweigs impl-last1
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24

> **Zitier-Form.** Kennung, nicht Adresse: `slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext:**

- slice-test-unter-last-ursache (Plan, §7 Belege des Implementers)
- keine aktive ADR; ADR-0055 (Zeitbudget der Gegenprobe) und ADR-0056 (Metrik) als Randbezug
- `LH-FA-20.a` (nur als Hinweis: auf main in 36b6387c geschärft, nicht im Diff)
- `AGENTS.md` §3.9, §3.11, §3.14 (Hard Rules)
- gelesen zur Plausibilität: `internal/adapters/driven/postgres/einspielen.go` (`Anfrage`, `schliesse`), `internal/bootstrap/play_test.go` (`TestRunPlayZweitesSignal`)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | §7 *Größe* nennt Summe 0 „am letzten Commit vor diesem Eintrag“; am Kopf des Zweigs liefert der Befehl aus §3.14 über die Pfade von §4 (`.claude` eingeschlossen) 1 (die Regelzeile in 00bf41a6). DoD-Punkt 4 prüft `numstat` nur über `internal test tools harness`, nicht über die Pfade von §4. | `AGENTS.md` §3.14, §3.9 | `slice-test-unter-last-ursache.md` · „lieferte am letzten Commit vor diesem Eintrag 0“ | ja — Befehl aus §3.14, Basis ffdacf7d, bis HEAD: 1 | Größen-Messzeile veraltet gegenüber Endstand |
| F-2 | LOW | §7 (d) zählt A und B mit „je 60 Läufen unter Last“; die Tabelle (b) nennt 15 + 15, und (a) sagt, die Läufe des ersten Versuchs (ein Container) „zählen nicht“. Zählt man sie nicht, sind es 30 je Test. Die Aussage „nicht reproduziert in N Läufen bei Last L“ trägt zwei verschiedene N. | `AGENTS.md` §3.11, Plan DoD 1 (d) | `slice-test-unter-last-ursache.md` · „je 60 Läufen unter Last“ | nein — Zählung der Läufe ist nicht aus dem Repo reproduzierbar (Skripte im Scratchpad) | Lauf-Zahl in Belegen inkonsistent |
| F-3 | LOW | Die Zwischenregel nennt als Anlass „(Frist von 30 s unter CPU-Last)“ für A und B; §7 (d) hält für beide fest, dass Last allein das Rot nicht erzeugt (0 von 60 bzw. 30) und die Ursache offen bleibt. Der Klammersatz sagt mehr zu, als §7 belegt. Der Wortlaut ist vom Nutzer bestätigt; der Befund betrifft die Reichweite der Aussage, nicht die Freigabe. | `AGENTS.md` §3.11 | `.claude/commands/implement-slice.md` · „(Frist von 30 s unter CPU-Last)“ | nein | Zusage im Kommentar weiter als Prüfung |
| F-4 | LOW | Plan §6 Randform 3 verlangt je Test mindestens 15 Läufe unter Last; für den Signalfall der Gegenprobe stehen 8 Läufe in §7, ohne dass der Plan die Abweichung nennt (Grund ist die Laufzeit von je 27 s). DoD 2 verlangt „Zahl der Läufe“, nicht 15; die Randform ist damit im Code-Commit still unterschritten. | `AGENTS.md` §3.9 (Plan folgt dem Lauf) | `slice-test-unter-last-ursache.md` · „Nicht reproduziert in 8 Läufen bei Last 16 bis 32“ | nein | Plan-Randform und Beleg weichen ab |
| F-5 | LOW | DoD 1 (a) verlangt den Aufruf „wörtlich und wiederholbar“. §7 enthält Platzhalter (`<serie>`, `K`, `<netz>`, `<pg>`), Kommentare in der Befehlszeile der Tests und für den Start von PostgreSQL nur „wie `make test-integration`“; die Skripte liegen im Scratchpad. Last vor und nach ist als Bereich (16 bis 32, 18 bis 48), nicht als Wert vor und nach je Serie genannt. | `AGENTS.md` §3.11; Plan DoD 1 (a) | `slice-test-unter-last-ursache.md` · „wie `make test-integration`: `postgres:17-alpine@sha256:b0f9560a…`“ | nein | Beleg nicht wörtlich wiederholbar |
| F-6 | INFO | §7 (d) zu C sagt, ob Code oder Test die Zusage bricht, „entscheide ich nicht“ und gehe an den Architect; auf main hat 36b6387c `LH-FA-20.a` geschärft (nach vollständigem Schreiben wird `Terminate` gesendet). Damit ist die Lesart entschieden, und der Satz im Plan ist beim Zusammenführen überholt. Die Zwischenregel („Wettlauf im Code, reproduzierbar“) deckt sich mit dem geschärften Text; ich habe den Code gegen ihn gelesen: `schliesse` sendet `Terminate` nur bei `TryLock`-Erfolg, `Anfrage` hält die Sperre bis nach `Flush`, dies bestätigt P1/P2. Behoben wird in einem Folge-Slice, nicht hier. | `AGENTS.md` §3.9 | `slice-test-unter-last-ursache.md` · „entscheide ich nicht: Das ist eine Randform der Spezifikation“ | nein | Plan-Satz nach Spec-Schärfung überholt |
| F-7 | INFO | Plausibel und gegen den Code gelesen: P1 (300 ms nach `Flush`: 5 von 5 rot allein) und P2 (`Lock` statt `TryLock`: 0 von 5) trennen Wettlauf von Frist belastbar; die Last erhöht nur die Fensterwahrscheinlichkeit (8 und 9 von 20). Für A und B bleibt die Ursache offen, die Risiko-Ausgänge dafür (weiter offen → `BEO-REPO/test-faellt-unter-cpu-last-aus`) tragen. Die Schleife-je-Container-Erkenntnis (cgroup-Verteilung) steht als Beleg, nicht als Zusage. | Maintainability | `slice-test-unter-last-ursache.md` · „P2: wie P1, zusätzlich `s.schreiben.TryLock()`“ | nein | — |

## Prüfung der drei Fragen

1. **Belege in §7.** Vorhanden: Last wörtlich (mit F-5), Last gemessen, je Test allein/unter Last, Proben P1 bis P4 mit dem, was sie trennen, Ursache je Test (C: Wettlauf im Code; A und B: „nicht reproduziert“ mit N und L), Signalfall der Gegenprobe „nicht reproduziert in 8 Läufen bei Last 16 bis 32“ mit Hinweis, dass ein Rot durch das Zeitbudget nicht auftrat. Die Ursachen-Aussage zu C stimmt mit `einspielen.go` und `play_test.go` überein (Rot-Bild: erste Frist hält, `typenBis` wartet auf `X`). Abweichungen: F-2, F-4, F-5.
2. **Zwischenregel gegen bestätigten Wortlaut.** Der Text in `implement-slice.md` Schritt 16 ist bis auf die Einbettung (Einrückung, „(seit slice-test-unter-last-ursache, Entscheidung des Nutzers vom 2026-10-24)“, „Die Regel gilt für den Implementer, nicht für Reviewer und Verifier“) mit dem Wortlaut in Plan §3 identisch; Wiederholung nur für die zwei Unterfälle, `TestRunPlayZweitesSignal` ausgenommen und als Befund geführt. Reichweite nach §3.11: nur F-3.
3. **Plan nachgezogen.** §1 (Ziel 3, Wortlaut-Absatz), §3 (Wortlaut, Reihenfolge), §6 (Randform 7, 8, Risiko 4), §8 (Nachzählung: 3 Punkte, 2 Schichten) und der DoD-Kopf folgen dem Zweig-Stand; Move-Commits rein, Inhalt in eigenen Commits (§3.3). Ausnahme: Größen-Messzeile (F-1) und überholter Satz (F-6). Marke 1500 mit Pfaden `.claude internal test tools harness` und den zwei Ausschlüssen: der Befehl liefert am Kopf 1 (nicht 0), weil `.claude` zählt; die Marke ist weit entfernt.
4. **Spezifikations-Hinweis** (36b6387c): kein Befund gegen diesen Slice; Code-Behebung bleibt dem Folge-Slice (Planner, Closure).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/`, `test/`, `tools/`, `harness/` (Diff) | geprüft, ohne Befund (numstat leer; kein Code, keine Frist, kein Test geändert) |
| `internal/adapters/driven/postgres/einspielen.go`, `internal/bootstrap/play_test.go` (gelesen) | geprüft, ohne Befund |
| `docs/plan/planning/` (Lifecycle, Basis-Hash) | geprüft, ohne Befund (Basis-Hash `ffdacf7d41a3…` stimmt mit dem Move-Commit überein) |
| `.claude/commands/` (eine Zeile) | geprüft, Befund nur F-3 |
| Spec-Strata, ADRs | geprüft, ohne Befund (nicht im Diff) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Größen-Messzeile veraltet gegenüber Endstand · Lauf-Zahl in Belegen inkonsistent · Zusage im Kommentar weiter als Prüfung · Plan-Randform und Beleg weichen ab · Beleg nicht wörtlich wiederholbar · Plan-Satz nach Spec-Schärfung überholt

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; die fünf LOW-Findings betreffen Beleg-Genauigkeit und Reichweite von Sätzen, keine Freigabe der Regel oder der Ursache von C.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Der Report ersetzt keine Verifikation gegen die DoD (Verifier).
