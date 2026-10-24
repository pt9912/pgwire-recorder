# Review-Report: slice-tests-ueberlebende-mutanten-driving (Delta) — 2026-10-24

**Review-Art:** Code (Delta-Review gegen Plan, `SPEC-038` und Hard Rules)

**Gegenstand:** `slice-tests-ueberlebende-mutanten-driving`, Commits `0950edcc` (Architect, `SPEC-038`), `9d4ddda9` (Frist an `<-danach`), `a3c120fc` und `1b522283` (Plan-Nachzug); Kopf `1b522283`. Der erste Review und die Verifikation (bestanden, V-1/V-2) liegen vor und sind nicht Gegenstand.

**Skill:** `.harness/skills/reviewer.md` @ `1b522283`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24

**Eingangs-Kontext:**

- `slice-tests-ueberlebende-mutanten-driving` (§1, §2, §3, §6, §7, §8, Kopf)
- `SPEC-038` *Warten in Tests* (`spec/spezifikation.md`, Stand `0950edcc`)
- `LH-FA-13`, `LH-FA-20` (berührt über `play_anmeldung_test.go`)
- `AGENTS.md` (Hard Rules §3.9, §3.10, §3.11)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | §3 (Plan) führt keine Zeile für `play_anmeldung_test.go`, obwohl der Commit sie ändert; die Bootstrap-Zeile nennt nur `bootstrap_test.go:353`, `tls_test.go:184` und `:100`. §1, §2, §6, §7, §8 sind nachgezogen, §3 nicht. | `AGENTS.md` §3.9 | `slice-tests-ueberlebende-mutanten-driving.md` · „`internal/bootstrap/bootstrap_test.go` | update" | nein (Urteil; `make kopf-check` prüft nur den Kopf) | Plan folgt Korrektur nicht |
| F-2 | LOW | §7 *Vollständige Durchsicht* zählt `play_anmeldung_test.go` (F-3) weiter unter die Empfänge ohne Frist, die „nach `SPEC-038` nicht unter die Regel fallen“. Das widerspricht der Entscheidung von `0950edcc` und §6 (Frist Pflicht). | `AGENTS.md` §3.9, §3.11 | `slice-tests-ueberlebende-mutanten-driving.md` · „dazu `play_anmeldung_test.go` (F-3). Alle anderen" | nein | Plan folgt Korrektur nicht |
| F-3 | INFO | Zeilenangaben sind nach dem Commit veraltet: `<-danach` steht jetzt bei 172 und 249, die Stellen mit `time.After` bei 197 und 285. Der Satz in der §6-Tabelle „Zeile 192 und 275 tragen `time.After`“ stimmt so nicht mehr. Die Zeile ist Lesehilfe; Anker ist der Name der Stelle. | Maintainability | `slice-tests-ueberlebende-mutanten-driving.md` · „Zeile 192 und 275 tragen `time.After`“ | nein | Zeilenangabe veraltet |
| F-4 | INFO | Der Fake-Server setzt `conn.SetDeadline(30 s)` ab Accept. Hält `play` die Verbindung offen, läuft die Deadline des Doubles zuerst ab (vor der Frist des Tests), der Double füllt `danach` mit den bisher gelesenen Typen, und die neue Frist greift nicht. Sie greift nur, wenn der Double nie füllt. Die Mutation in §7 trifft genau diesen Fall; ein „Prüfling beendet die Verbindung nie“ färbt der Wert- und Exit-Code-Vergleich, nicht diese Frist. Keine Aktion nötig. | Maintainability | `internal/bootstrap/play_anmeldung_test.go` · `_ = conn.SetDeadline(time.Now().Add(30 * time.Second))` | nein | — |

## Prüfpunkte des Auftrags

- **(a) Literal und Meldung:** beide Stellen (`TestRunPlayPasswort`, `TestRunPlayErstesSignalInAnmeldung`) tragen `case <-time.After(30 * time.Second)` als Literal, unter 60 s. Die Meldung „der Fake-Server sieht binnen 30 s kein Ende der Verbindung“ nennt das ausgebliebene Ereignis. Die Zeitgrenze von `go test` wird nicht als Frist benutzt. Ohne Befund.
- **(b) Mutation „Double füllt `danach` nie“:** Die Mutation (beide `typen <- string(gelesen)` der Lesestufe entfernt) bricht den Double, nicht den Prüfling. Im Code ist `typen` der Kanal, den `danach` liefert; ohne Füllen bliebe der Empfang ohne Frist bis zum Paket-Timeout (250 s) stehen. Der Beleg nennt Rot nach 30,03 s und 30,54 s, vor dem Paket-Timeout, mit der Frist-Meldung. Das ist der richtige Grund. Gegen den Test-Code gelesen stimmt das: der Empfang ist die einzige Stelle, die auf den Kanal wartet. Ohne Befund. Zur Grenze des Belegs siehe F-4.
- **(c) Plan folgt Korrektur:** §1 (Zeile 96), §2 (DoD *Frist beim Warten*), §6 (Tabelle, Satz zum Nicht-Anfassen), §7 (Nacharbeit, F-2/F-3 des ersten Reviews) und §8 (*Nachzählung*) nennen die vierte Stelle mit der Kennung des Architect-Commits. Der Kopf führt `SPEC-038` schon. Die Zählung 3 Liefer-Punkte und 2 Schichten stimmt: zwei Zeilen in einer Testdatei der Schicht Bootstrap, kein neuer Liefer-Punkt. Die Verneinung der Sendung an `slice-harness-lint-warten-ohne-frist` ist begründet (§1 dort schließt Stellen im Bootstrap aus). Lücken: F-1 (§3) und F-2 (§7).
- **(d) `SPEC-038` gegen alle Wartestellen:** `grep` nach `<-` in allen `*_test.go` von `internal/adapters/driving/pgwire` und `internal/bootstrap`. Jeder Empfang auf Testseite steht in einem `select` mit `time.After` oder in einer Hilfsfunktion mit Frist (`zurueckBinnen`, `anmeldePasswort`, `warte`). Fehlend wäre eine Stelle ohne Frist ohne Ausnahme: keine gefunden. Ausnahmen, die der Text von `SPEC-038` trägt: die Test-Doubles `server_test.go` (121, 129, 217, 233, 730) und `server_extended_test.go:790` (Freigabe durch den Test) sowie vier `select … default:` (`play_diff_all_test.go:54`, `play_test.go:266` und `:306`, `play_zeit_test.go:53`), die nicht warten. `.Wait()` hat keinen Treffer; `range` in `leser_test.go:30` läuft über eine Map. Der neue Absatz fügt sich in den Text von `SPEC-038` ein, ohne den Absatz davor (Kanal nur des Testgeschirrs, Frist erlaubt) zu widersprechen. Wo die Grenze liegt („mit dem, was der Prüfling gesendet hat“), nennt der Satz Beispiel und Kriterium.
- **(e) §3.11/§3.10:** Die Frist ist eine Zusage der Spezifikation; Zusage, Mutation und roter Test stehen in §7 (V-2). Kein Kommentar im Diff sagt über die Prüfung hinaus etwas zu. Der Satz im Änderungsprotokoll der Spezifikation ist Chronik der Tabelle, nicht eines Zustandsfeldes. Ohne Befund.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/bootstrap/play_anmeldung_test.go` (Diff `9d4ddda9`) | geprüft, ohne Befund |
| `spec/spezifikation.md` `SPEC-038` (Diff `0950edcc`) | geprüft, ohne Befund |
| `internal/adapters/driving/pgwire/*_test.go` (Wartestellen) | geprüft, ohne Befund |
| `internal/bootstrap/*_test.go` (Wartestellen) | geprüft, ohne Befund |
| Kopf und §1, §2, §6, §7, §8 des Slice-Plans | geprüft, Befunde F-1 bis F-3 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Plan folgt Korrektur nicht · Zeilenangabe veraltet

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Die zwei LOW betreffen nur den Plan (§3, §7).

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation.
