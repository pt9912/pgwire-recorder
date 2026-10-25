# Review-Report: slice-play-zweites-signal-terminate-wettlauf — 2026-10-25 (r2)

**Review-Art:** Code (Folgereview; Diff gegen Plan, Spezifikationsstelle `LH-FA-20.a` *Abbruchsignal* und die Hard Rules)

**Gegenstand:** slice-play-zweites-signal-terminate-wettlauf, Zweig impl-wettlauf1, Diff dced44df..c4f9c60c (Nacharbeit 9295e31f, 36e51415, Entscheidung des Architect fa8a8479, 9b56c33c, c4f9c60c); erster Review: `docs/reviews/2026-10-25-review-slice-play-zweites-signal-terminate-wettlauf.md` (F-1 MEDIUM, F-2 bis F-5 LOW)

**Skill:** `.harness/skills/reviewer.md` @ Stand des Zweigs impl-wettlauf1
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-25

> **Zitier-Form.** Kennung, nicht Adresse: `slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext:**

- slice-play-zweites-signal-terminate-wettlauf (§1, §3, §6 Randformen 1 bis 8, §7 Belege des Implementers; der Bericht des Implementers lag nicht vor)
- keine aktive ADR; Randformen 1 bis 5 und 8 stehen in `spec/spezifikation.md` (`LH-FA-20.a` *Abbruchsignal*)
- `LH-FA-20`, `LH-FA-20.a`
- `AGENTS.md` §3.7, §3.9, §3.10, §3.11, §3.12 (Hard Rules)
- Vorheriger Report (r1) mit F-1 bis F-8
- gelesen: Diff `internal/adapters/driven/postgres/einspielen.go`, `export_test.go`, `einspielen_terminate_test.go` ganz, Port `einspielziel.go`, `spec/spezifikation.md`, `docs/user/abdeckung-unit.md`, Plan §1 bis §7; nicht ausgeführt: `go test`, `-race`, `make gates` (Auftrag: kein Lauf)

---

## Prüfung der Fragen des Auftrags

**(a) F-1 gegen die neue Konstruktion — behoben.** `langeFrist` ist eine Sekunde (`NeueEinspielSessionMitFrist`). `TestEinspielSchliesseGesamtfrist` wartet nicht mehr auf `time.Sleep` ab dem Aufruf, sondern auf ein Ereignis (eine Gruppe wird abgelehnt, also hat `schliesse` `geschlossen` gesetzt und seine Frist schon berechnet), schläft dann 200 ms und gibt den Schreibvorgang frei; Spielraum bis zum Reißen der Frist rund 800 ms statt 20 ms. Die Tests `…SendetNachVollstaendigerAnfrage` und `…Gruppe` synchronisieren über einen Kanal, den der Fake beim ersten `Write` schließt, nicht über Schlafen; `endeSchreiben` muss binnen einer Sekunde nach `Schliesse` laufen. Die Schranke in `…Gesamtfrist` (`fristen[0]` nicht nach `ende + langeFrist`) vergleicht zwei Werte von `time.Now()` mit monotonem Anteil; sie hängt nicht an der Wanduhr und nicht an der Last: `frist` entsteht vor `ende`, also ist sie immer kleiner. Die Mutation „Schreibfrist nach dem Warten neu gesetzt“ liegt immer darüber. `TestEinspielWartePunktlich` hängt nicht an Zeit außer den 5 s des Hängenschutzes (SPEC-038: Frist für den Wert, den der Prüfling bestimmt, vorhanden); die 200 Wiederholungen gegen die Vorab-Prüfung beruhen darauf, dass ein Timer mit abgelaufener Frist unter `go 1.27` im `select` bereit ist, und §7 belegt den roten Lauf der Mutation. Die Mutanten mit 3 s Verzögerung sind die richtigen: Sie liegen 2 s über der Test-Frist und weit unter dem Hängenschutz von 5 s. Rest-Risiko: Überschlafen von `time.Sleep(200 ms)` um mehr als 800 ms (siehe F-4, INFO).

**(b) F-2 bis F-5.** F-2 behoben: Die Zeilen in `abdeckung-unit.md` sagen „sobald der Schreibvorgang endet“ und „200 ms, nachdem das Schließen begonnen hat“; beides trägt der Test (200 ms ist Untergrenze). F-3 behoben: §3 nennt `einspielen_terminate_test.go` neu und lässt die Bestandstests unverändert. F-4 behoben: Randform 8 in §6 mit Verweis auf fa8a8479, Entscheidung in `LH-FA-20.a` *Abbruchsignal* (kein Schreiben, `PGR-E4003`, ohne Meldung), Zeile im Änderungsverzeichnis der Spezifikation; Test `TestEinspielAnfrageNachSchliessen` für beide Arten. Der Spec-Satz nennt weder ADR noch Slice. F-5 behoben: Fall und Mutation für die Vorab-Prüfung stehen in §7. Neue Rest-Befunde dazu: F-1 bis F-3 unten.

**(c) Neuer Code.** Die Frist als Session-Feld mit `neueEinspielSessionMitFrist` ändert den Produktionscode nur um einen Konstruktor, der von `neueEinspielSession` mit `terminateFrist` gerufen wird; es gibt keinen weiteren Konstruktor der Session (`grep`), ein Nullwert der Frist kann nicht entstehen. Die Alternative (Paketvariable) wäre unter `-race` und parallelen Tests schlechter; der Eingriff ist angemessen. Der Export ist test-only (`export_test.go`), `WartePunktlich` reicht nur weiter. `defer s.endeSchreiben()` in `schreibeGruppe` schließt die in r1 genannte Kante (F-7 r1): ein Panic in `Send`/`Flush` lässt `Schliesse` nicht bis zur Frist warten.

**(d) Nebenläufigkeit durch Lesen (kein `-race`), ohne Befund.** `schreibEnde` wird nur unter `zustand` gelesen und geschrieben; `beginneSchreiben` läuft nur unter `schreiben`, also nie zweimal ohne dazwischenliegendes `endeSchreiben`, und jeder Kanal wird genau einmal geschlossen und danach auf `nil` gesetzt. `schliesse` liest `geschlossen` und `schreibEnde` in einem Schritt: Der Schnappschuss ist der letzte Vorgang, danach beginnt keiner mehr (`beginneSchreiben` liefert falsch), sodass `Terminate` nie neben einen Schreibvorgang schreibt; die Reihenfolge `close(schreibEnde)` → Rückkehr des Wartens ordnet das `Flush` vor das `Terminate`. In `Anfrage` läuft `endeSchreiben` (zuletzt registriertes `defer`) vor `schreiben.Unlock`, in `sendeGruppe` über `schreibeGruppe` ebenfalls; ein Deadlock entsteht nicht, weil `schliesse` `schreiben` nicht nimmt und jedes Warten an die Frist gebunden ist. Der Fehlerpfad von `sendeGruppe` (Lesen von `geschlossen` nach `endeSchreiben`) setzt `sendefehler` nur bei offener Session; unverändert gegenüber r1.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | §7 widerspricht dem Stand nach der Nacharbeit an zwei Stellen: Der Absatz „Beobachtungen für die Closure“ sagt, `TestEinspielSchliesseGesamtfrist` hänge an einem Spielraum von 80 ms und könne unter Last rot werden; die Nacharbeit (selber Abschnitt) hat ihn auf rund 800 ms gebracht. Der Absatz „Kommentare“ sagt, der Port-Kommentar „soweit die Verbindung es sofort annimmt“ bleibe; der Diff hat ihn durch „binnen einer kurzen Gesamtfrist“ ersetzt. | `AGENTS.md` §3.9 (Plan folgt jeder Korrektur), §3.7 (Zustand statt Chronik) | `docs/plan/planning/in-progress/slice-play-zweites-signal-terminate-wettlauf.md` · „hängt an einem Spielraum von 80 ms“ und „bleibt; er deckt die Gesamtfrist nicht ausdrücklich“ | nein — Urteil | Plan folgt Korrektur nicht (BEO-REPO/plan-folgt-korrektur-nicht) |
| F-2 | LOW | §3 beschreibt die Zeile `einspielen.go` noch mit dem Zustand „teilweise geschrieben“, „gesetzt für die Dauer des Schreibens einer Nachricht“; der Code und Randform 2 (entschieden: je Schreibvorgang, `Flush`) führen „schreibt“ je `Flush` von `Anfrage` oder Gruppe. Die Zeile zu `spec/spezifikation.md` sagt „unverändert“, der Diff ändert sie (Randform 8, Text erklärt es in derselben Zeile). | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-play-zweites-signal-terminate-wettlauf.md` · „gesetzt für die Dauer des Schreibens einer Nachricht“ | nein — Urteil | Plan folgt Korrektur nicht (BEO-REPO/plan-folgt-korrektur-nicht) |
| F-3 | LOW | Randform 8 entscheidet in der Spezifikation drei Zusagen: nicht geschrieben, Fehlercode `PGR-E4003`, ohne Meldung. §7 hat die Mutation nur für „nicht geschrieben“ (Prüfung von `geschlossen` entfällt); für den Code (anderer Code, etwa der Ausgang des Schreibens auf die geschlossene Verbindung) und für „ohne Meldung“ steht keine Mutation und kein Fall. Der Test prüft den Code (`code(err)`), eine Mutation dazu ist nicht belegt. Ohne Failure-Szenario: der Code verhält sich wie die Spezifikation. | `AGENTS.md` §3.10 (Aufzählung je Element) | `internal/adapters/driven/postgres/einspielen.go` · „Anfrage nach dem Schließen der Verbindung“ | ja — Mutation „anderer Code in `Anfrage` nach dem Schließen“ in einer Kopie, Lauf des Pakets | Mutation je Element der Zusage fehlt |
| F-4 | INFO | Rest-Spielraum der Gate-Tests: Überschläft `time.Sleep(200 ms)` in `…Gesamtfrist` oder die Goroutine zwischen `Write` und `endeSchreiben` die Frist von einer Sekunde, bleibt `Terminate` aus und der Test ist rot. §7 belegt 2 × 20 Läufe ohne Rot bei K=32 und K=48 und nennt selbst die Grenze („0 rot belegt keinen Spielraum“); die Konstruktion (4-facher Spielraum, ereignisgesteuert) trägt sie. | Maintainability (`BEO-REPO/test-faellt-unter-cpu-last-aus`) | `internal/adapters/driven/postgres/einspielen_terminate_test.go` · „time.Sleep(200 * time.Millisecond)“ | ja — viele Läufe unter der Last aus `slice-test-unter-last-ursache` | — |
| F-5 | INFO | Die Spezifikation wurde mit fa8a8479 nach dem letzten vollen Lauf von `make gates` (Code-Commit 9295e31f, §7) geändert; §7 nennt keinen Lauf am Endstand. Doku-Gates (Links, Kopf) laufen über die Spezifikation. Verifier und Gate-Hook des Zweigs fangen das; ob der Stop-Hook eine Spezifikationsänderung freigibt, hängt an seiner Positivliste (`spec/` steht nicht darauf). | Maintainability | `docs/plan/planning/in-progress/slice-play-zweites-signal-terminate-wettlauf.md` · „erneut nach der Nacharbeit (Code-Commit `9295e31f`)“ | ja — `make gates` am Endstand | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/postgres/einspielen.go` (Nebenläufigkeit, Deadlock, Kanäle, Zustand, `defer`) | geprüft, ohne Befund (durch Lesen, kein `-race`) |
| `internal/adapters/driven/postgres/export_test.go` (Test-only-Export) | geprüft, ohne Befund |
| `internal/adapters/driven/postgres/einspielen_terminate_test.go` (Spielraum, Wanduhr, Mutanten mit 3 s) | geprüft, ohne Befund außer F-3, F-4 |
| `internal/hexagon/ports/driven/einspielziel.go` (Kommentar `Schliesse`) | geprüft, ohne Befund (Kommentar beschreibt den Ist-Zustand) |
| `spec/spezifikation.md` (Stratum-Regeln, ADR-/Slice-Bezug, Änderungsverzeichnis) | geprüft, ohne Befund |
| `docs/user/abdeckung-unit.md` | geprüft, ohne Befund (F-2 aus r1 behoben) |
| Core-Reinheit (`internal/hexagon/**`) | geprüft, ohne Befund (nur Kommentar) |
| Suppression (`nolint`) im Diff | geprüft, ohne Befund |
| Commit-Regel §3.3 (Move und Inhalt getrennt) | geprüft, ohne Befund (kein Move in dced44df..c4f9c60c außer einer Link-Zeile in `done/slice-test-unter-last-ursache.md`, die vor dieser Range liegt) |
| Handbuch und README | geprüft, ohne Befund (unverändert, laut §1 ausgeschlossen) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Plan folgt Korrektur nicht · Mutation je Element der Zusage fehlt

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. F-1 aus r1 (MEDIUM) ist behoben, F-2 bis F-5 aus r1 sind behoben. Offen sind drei LOW (Plan-Nachzug in §3 und §7, eine fehlende Mutation für Code und „ohne Meldung“ der Randform 8) und zwei INFO; sie gehen an den Implementer und blockieren nicht. Die Klassen gehen in die Slice-Closure §7 und von dort in den Zähler; „Plan folgt Korrektur nicht“ ist die Klasse `BEO-REPO/plan-folgt-korrektur-nicht`.

**Übergabe:** Findings gehen an den Implementer; der Report ersetzt keine Verifikation (DoD-Konformität, insbesondere die Läufe unter Last und `make gates` am Endstand, prüft der Verifier).
