# Review-Report: slice-play-zweites-signal-terminate-wettlauf — 2026-10-25

**Review-Art:** Code (Diff gegen Plan, Spezifikationsstelle `LH-FA-20.a` *Abbruchsignal* mit den Randformen aus f17d9027 und die Hard Rules)

**Gegenstand:** slice-play-zweites-signal-terminate-wettlauf, Zweig impl-wettlauf1, Diff dced44df..impl-wettlauf1 (Commits b0a7e10a Code und Tests, b8b39d09 §7 und DoD)

**Skill:** `.harness/skills/reviewer.md` @ Stand des Zweigs impl-wettlauf1
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-25

> **Zitier-Form.** Kennung, nicht Adresse: `slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext:**

- slice-play-zweites-signal-terminate-wettlauf (Plan mit §6 Randformen 1 bis 7 und §7 Belege des Implementers; der Bericht des Implementers lag nicht vor)
- keine aktive ADR; Randformen 2 bis 5 stehen in `spec/spezifikation.md` (`LH-FA-20.a` *Abbruchsignal*, f17d9027)
- `LH-FA-20`, `LH-FA-20.a`
- `AGENTS.md` §3.7, §3.9, §3.10, §3.11, §3.12 (Hard Rules)
- gelesen: `internal/adapters/driven/postgres/einspielen.go` ganz, `einspielen_terminate_test.go` ganz, Port `internal/hexagon/ports/driven/einspielziel.go` (Kommentar zu `Schliesse`); nicht ausgeführt: `go test`, `-race`, `make gates` (Auftrag: kein Lauf)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | `TestEinspielSchliesseGesamtfrist` hat effektiv rund 20 ms Spielraum, nicht 80 ms: `Schliesse` setzt seine Frist von 100 ms am Aufruf, der Test schläft 80 ms und schließt dann `frei`; überschläft `time.Sleep` unter CPU-Last um mehr als etwa 20 ms, reißt die Frist vor dem Ende des Vorgangs, es gibt keine Schreibfrist und der Test ist rot (`len(conn.fristen) != 1`). Zusätzlich verlangt `TestEinspielSchliesseSendetNachVollstaendigerAnfrage` (und der Gruppen-Fall), dass `endeSchreiben` binnen 100 ms nach dem Write läuft; eine Verzögerung der Goroutine darüber lässt `Terminate` aus. Die Last-Läufe in §7 (20 je Last) fahren nur `TestRunPlayZweitesSignal`, nicht diese Tests. | Maintainability (Muster `BEO-REPO/test-faellt-unter-cpu-last-aus`, im Register 7×; der Slice soll dieses Muster beheben) | `internal/adapters/driven/postgres/einspielen_terminate_test.go` · „time.Sleep(80 * time.Millisecond)“ | ja — der Test unter der Last aus slice-test-unter-last-ursache §7 (K=32, K=48), viele Läufe | Test mit Zeitgrenze ohne Spielraum unter Last |
| F-2 | LOW | Die Abdeckungs-Zeile von `TestEinspielSchliesseSendetNachVollstaendigerAnfrage` sagt zu, der Schreibvorgang gelte nicht als laufend („ohne dass der Schreibvorgang als laufend gilt“); im Code gilt er bis `endeSchreiben` nach dem Rückkehr von `Flush` als laufend, der Test besteht, weil `Schliesse` darauf wartet. Ebenso nennt die Zeile von `TestEinspielSchliesseGesamtfrist` „80 ms nach dem Aufruf“, gemessen wird `ende` nach dem Sleep. Ein Test prüft den Satz nicht. | `AGENTS.md` §3.11 | `docs/user/abdeckung-unit.md` · „ohne dass der Schreibvorgang als laufend gilt“ (Quelle: Kommentar in `einspielen_terminate_test.go`) | nein — Lesart, kein Gate | Zusage im Kommentar weiter als Prüfung |
| F-3 | LOW | Plan §3 nennt als geänderte Test-Dateien `einspielen_ende_test.go` und `einspielen_test.go` („update“); der Diff legt `einspielen_terminate_test.go` neu an und ändert die beiden genannten nicht. §3 ist im selben Commit nicht nachgezogen. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-play-zweites-signal-terminate-wettlauf.md` · „einspielen_ende_test.go`, `einspielen_test.go`“ | nein — Urteil | Plan folgt Korrektur nicht (BEO-REPO/plan-folgt-korrektur-nicht) |
| F-4 | LOW | `Anfrage` nach dem Schließen liefert `PGR-E4003` ohne zu schreiben (neue Verhaltensänderung, mit `TestEinspielAnfrageNachSchliessen`); die Stelle „Anfrage nach dem Schließen“ steht weder in §6 noch in der Spezifikation. §7 hält fest „Keine Randform außerhalb von §6 berührt“. Fehlercode und Ergebnis auf einer echten Verbindung sind wie zuvor (der Write auf die geschlossene Verbindung scheiterte mit demselben Code), nur die Meldung ändert sich; daher ohne Failure-Szenario. | `AGENTS.md` §3.12 | `internal/adapters/driven/postgres/einspielen.go` · „Anfrage nach dem Schließen der Verbindung“ | nein — Urteil | Randform im Code-Commit entschieden |
| F-5 | LOW | In `wartePunktlich` entscheidet die Vorab-Prüfung `select { case <-ende: return true; default: }`, ob bei bereits abgelaufener Frist und geschlossenem `ende` deterministisch gesendet wird; ohne sie wählt das folgende `select` bei beiden bereiten Kanälen zufällig. §7 führt dafür keine Mutation und keinen Fall; sie gilt dort auch nicht als äquivalent. Die Zusage „Warten und Schreiben teilen die Gesamtfrist“ nennt die Ortsangabe „Ende vor der Frist“ als Element. | `AGENTS.md` §3.10 (Aufzählung und Ortsangabe je Element) | `internal/adapters/driven/postgres/einspielen.go` · „case <-ende:\n\t\treturn true\n\tdefault:“ (Vorab-Prüfung in `wartePunktlich`) | ja — Mutation „Vorab-Prüfung entfernen“ in einer Kopie, Lauf des Pakets | Mutation je Element der Ortsangabe fehlt |
| F-6 | INFO | Zwei Mutationen in §7 (Warten entfällt, `Terminate` trotz Frist) färben ihren Test über einen Panic des Fakes (`close of closed channel` in `einspielen_gruppe_test.go`), der das Testbinary beendet und die übrigen Ergebnisse des Laufs verdeckt; der Implementer nennt es offen in §7. Rot ist es dennoch aus dem Grund der Zusage (zweites gleichzeitiges `Write`). | Maintainability | `docs/plan/planning/in-progress/slice-play-zweites-signal-terminate-wettlauf.md` · „Zwei der Fälle (Warten entfällt, Terminate trotz Frist)“ | ja | — |
| F-7 | INFO | Die Nebenläufigkeit habe ich durch Lesen geprüft, nicht mit `-race` (kein Gate der Kette fährt `-race`; Docker-Lauf war ausgeschlossen). Ergebnis ohne Befund: `schreibEnde` wird nur unter `zustand` gesetzt und gelesen; `beginneSchreiben` läuft nur unter der Sperre `schreiben`, also nie zweimal ohne `endeSchreiben` dazwischen, jeder Kanal wird genau einmal geschlossen und danach auf `nil` gesetzt; `schliesse` hält `zustand` nicht während des Wartens; `geschlossen` und `schreibEnde` werden in einem Schritt gelesen, sodass nach dem Schließen kein Vorgang beginnt und `Terminate` nie in einen Vorgang schreibt (Happens-before über `close(schreibEnde)` bzw. über `zustand`); die Frist beendet jedes Warten, ein Deadlock mit `schreiben` entsteht nicht, weil `schliesse` die Sperre nicht nimmt. Einzige Kante: `sendeGruppe` ruft `endeSchreiben` nicht per `defer`; ein Panic in `Send` oder `Flush` ließe `Schliesse` bis zur Frist warten (kein Deadlock). | Maintainability | `internal/adapters/driven/postgres/einspielen.go` · „func (s *einspielSession) schliesse(terminate bool)“ | nein — `-race` fährt kein Gate | — |
| F-8 | INFO | Der Rest-Wettlauf bleibt dort, wo die Spezifikation ihn setzt: Wird die Goroutine von `Anfrage` zwischen der Rückkehr von `Flush` und `endeSchreiben` länger als die Gesamtfrist aufgehalten, bleibt `Terminate` aus. Der Port-Kommentar „soweit die Verbindung es sofort annimmt“ nennt das Warten bis 100 ms nicht; der Implementer führt das in §7 selbst, die Spezifikation trägt es. Ein Beleg unter Last für die Gate-Tests fehlt (siehe F-1). | `LH-FA-20.a` | `internal/hexagon/ports/driven/einspielziel.go` · „Schliesse sendet Terminate, soweit die Verbindung es sofort annimmt“ | nein | — |

## Prüfung der Fragen des Auftrags

- **Zustand statt Sperre, `TryLock`-Wegfall:** `schliesse` fragt `schreibEnde` (Zustand), nicht die Sperre; Granularität je `Flush` (`Anfrage`, `sendeGruppe`) wie in der Spezifikation; Warten und Schreiben teilen `frist`, am Aufruf gesetzt, das Schreiben benutzt `SetWriteDeadline(frist)` dieselbe Frist.
- **Reihenfolge der Gründe:** bereits geschlossen (`erstes`) → `terminate`-Parameter → Frist (`wartePunktlich`); „teilweise geschrieben“ ist kein eigener Grund, wie die Spezifikation es entscheidet. Je Grund ein Fall in §7; die Reihenfolge der letzten zwei ist laut Spezifikation nicht beobachtbar (Architect-Entscheidung, kein äquivalenter Mutant des Implementers).
- **Äquivalente Mutanten:** §7 führt keinen als äquivalent außer der von der Spezifikation entschiedenen Reihenfolge Warten/Schreiben; sondiert: beide Wege enden ohne `Terminate` und mit `Close`, die Sonde trennt sie nicht. Offen ist nur F-5 (nicht als äquivalent geführt).
- **Plattform:** Grenze für macOS und Windows steht an der Zusage in der Spezifikation; kein Handbuchsatz.
- **Kommentare (§3.7):** `einspielSession`, `beginneSchreiben`, `endeSchreiben`, `Schliesse`, `schliesse`, `wartePunktlich` beschreiben den Ist-Zustand; Befund nur F-2 (Abdeckungs-Zeile).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/postgres/einspielen.go` (Nebenläufigkeit, Deadlock, doppeltes Schließen von Kanälen, Zustand) | geprüft, ohne Befund (F-7 nur Lesung, ohne `-race`) |
| Core-Reinheit (`internal/hexagon/**`) | geprüft, ohne Befund (Diff berührt sie nicht) |
| Suppression (`nolint`) im Diff | geprüft, ohne Befund |
| `spec/` (Stratum-Regeln, ADR-Immutabilität) | geprüft, ohne Befund (Diff ändert die Spezifikation nicht) |
| Commit-Regel §3.3 (Move und Inhalt getrennt) | geprüft, ohne Befund (beide Commits ohne Move) |
| `docs/user/abdeckung-unit.md` (Tabelle) | geprüft, ohne Befund außer F-2 |
| Handbuch und README | geprüft, ohne Befund (unverändert, laut §1 ausgeschlossen) |
| Größe §3.14 | geprüft, ohne Befund (252 laut §7, Marke 1500, nicht nachgemessen) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Test mit Zeitgrenze ohne Spielraum unter Last · Zusage im Kommentar weiter als Prüfung · Plan folgt Korrektur nicht · Randform im Code-Commit entschieden · Mutation je Element der Ortsangabe fehlt

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM): ein Slice, der den Last-Flake von `TestRunPlayZweitesSignal` behebt, bringt zwei neue Tests mit Zeitgrenzen ein, deren Spielraum unter derselben Last unterschritten wird. F-2 bis F-5 sind LOW und blockieren nicht. Der Code selbst (Nebenläufigkeit) hat keinen Befund.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Der Report ersetzt keine Verifikation (DoD-Konformität, insbesondere die 20 Läufe je Last, prüft der Verifier). Dateiname nach Skill-Schema (`<Datum>-review-<Kennung>.md`), nicht nach dem im Auftrag genannten.
