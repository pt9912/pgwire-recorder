# Review-Report: slice-harness-blackbox-pgwire — 2026-10-07

**Review-Art:** Code-Review (Testumbau, kein Produkt-Code) gegen Plan, Entscheidungen (`SPEC-049` Punkt 5, 7 in der Fassung `f0ea7be` und 8, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5) und Hard Rules (Modul 10). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git show b37d6e7`. Die drei Testdateien unter `internal/adapters/driving/pgwire` sind jetzt im Paket `pgwire_test`. Neu ist die Brücke `internal/adapters/driving/pgwire/export_test.go`. `server.go` ist unverändert. Schwerpunkte laut Auftrag: die Brücke nach Punkt 7, der umgeschriebene `TestWeiterlesenNachWecken`, `sendeMu` → `nebenher(fe)`, die zwei grünen Mutanten aus §7, Prüfung im Umbau verloren und eigene Mutationen.

**Skill:** `.harness/skills/reviewer.md` (unverändert bis `b37d6e7`)
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-harness-blackbox-pgwire.md` (ganz, Stand `b37d6e7`), dazu `git log --follow` auf den Plan. §7 *Belege des Implementers* habe ich gelesen, den Bericht des Implementers habe ich nicht.
- `spec/spezifikation.md` `SPEC-049` Punkt 1 bis 10 und Grenze; der Abschnitt *Record* (Session-Ende, „ein blockiertes Schreiben an den Client endet damit“); [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)
- `AGENTS.md` §3.2, §3.3, §3.6, §3.7, §3.9 bis §3.12; `.claude/commands/implement-slice.md` Schritt 19
- Produkt-Code `internal/adapters/driving/pgwire/server.go` (`handle`, `recordSitzung`, `richtungen`, `wecke`, `weiterlesen`, `beende`, `clientRichtung`, `fail`, `note`); Testdateien an `b37d6e7^` und `b37d6e7`; `5d666db` (Herkunft von `sendeMu`)
- Vorherige Reports: die Reviews zu `slice-harness-blackbox-kern` und `slice-harness-blackbox-driven`, dazu für die Klasse „Plan folgt Korrektur nicht“ die Reviews zu `slice-harness-vertraege-spezifikation` (F-425, LOW) und `slice-harness-lint-werkzeug` (F-430, LOW). Die höchste vergebene Nummer vor diesem Lauf war F-447.

**Ausgeführte Läufe:**

- `make lint`: Exit 2, `54 issues:`. Unter `internal/adapters/driving/pgwire` stehen 14 Befunde, alle in `server.go`, keiner in einer `*_test.go` und keine `lint:`-Zeile.
- `make kopf-check` und `make abdeckung-check`: beide Exit 0.
- Ich habe `b37d6e7^` und `b37d6e7` per `git archive` in je eine frische Kopie unter dem Scratch-Pfad entpackt. Die Läufe dort liefen in einem eigenen Image der Stufe `deps` (eigener Tag), ohne Netz, per Bind-Mount. `go test -list .` ergibt vorher und nachher 40 Tests, der `diff` ist leer.
- In der Kopie von `b37d6e7` ist `gofmt -l` leer, `go vet` ist grün und `go test` des Pakets ist grün.
- `TestWeiterlesenNachWecken` lief zweimal mit `-count=300`: einmal mit `--cpus=0.3` und `GOMAXPROCS=1`, einmal mit `--cpus=1` und `GOMAXPROCS=4`. Beide Läufe waren grün (600 von 600).
- In den geänderten Zeilen des Diffs steht keine Kommentarzeile `Abdeckung:`.
- `go test -race` geht im gepinnten Alpine-Image nicht (`CGO_ENABLED=0`, kein gcc). Die Synchronisation habe ich deshalb lesend beurteilt (Schwerpunkt 3).
- Sieben eigene Läufe (sechs Mutationen und ein Eingriff in den Zeitablauf), je in einer frischen Kopie (`cp -r` ohne `-p`, Ankunft der Mutation per Prüfsumme kontrolliert), Tests per Bind-Mount mit `go test -count=1 -run` (Tabelle unten). Danach habe ich die Kopien gelöscht. Das Repo blieb bis auf diese Datei unverändert.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-448 | MEDIUM | §6 widerspricht sich und dem Diff. Der Punkt *Export-Test-Brücke* sagt, entschieden in `SPEC-049` Punkt 7: „eine Frist für einen Test verkürzt sie nur an einem übergebenen Wert“. Der Punkt *White-Box-Zugriffe im Bestand* nennt dieselbe Frage weiter offen: „das ist die offene Frage, ob die Brücke schreiben darf“. Der Code-Commit beantwortet sie: `MeldeFrist` ist eine Konstante und wird nur gelesen, und kein Test setzt eine Frist (§7). `b37d6e7` ändert im Plan aber nur §7. Failure-Szenario: Wer an §6 misst, also Verifier, Closure oder der nächste Umstellungs-Slice `slice-harness-blackbox-einstieg`, liest eine offene Randform. Nach `AGENTS.md` §3.12 hätte der Implementer bei einer offenen Randform anhalten müssen. Je nach Leser entsteht daraus ein falscher Befund gegen den Implementer, oder die beantwortete Frage geht erneut an den Architect. Die Klasse war schon zweimal LOW (F-425, F-430), deshalb MEDIUM. | `AGENTS.md` §3.9; `AGENTS.md` §3.12; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“; `BEO-REPO/plan-folgt-korrektur-nicht` | `docs/plan/planning/in-progress/slice-harness-blackbox-pgwire.md` · „das ist die offene Frage, ob die Brücke schreiben darf“ | nein (Lesen) | Plan folgt Korrektur nicht |
| F-449 | LOW | Eine Prüfung ist im Umbau weggefallen, und §7 nennt sie nicht. Der alte `TestWeiterlesenNachWecken` prüfte den Rückgabewert von `weiterlesen` („Wecken nicht gemeldet“, „Wecken ohne Signal gemeldet“). Der neue Test kann ihn nicht sehen, weil `clientRichtung` den Wert verwirft (`r.weiterlesen()` als Anweisung). Meine Mutation R1 (`return true` → `return false` in `weiterlesen`) ist an `b37d6e7^` rot und an `b37d6e7` in allen Tests des Pakets grün. Über die Schnittstelle ist der Mutant äquivalent, nach `SPEC-049` Punkt 7, letzter Fall, ist der Wegfall also zulässig. Plan §6 verlangt aber: „eine Mutation, die bisher ein Test auf eine unexportierte Funktion fing, muss auch nach dem Umbau fangen; sonst fehlt eine Prüfung“. Schritt 19 verlangt, einen grünen Mutanten in §7 einzuordnen. §7 tut weder das eine noch das andere. Die Mutationstabelle führt nur Mutanten, die rot werden. | Plan §6 *Mutationstests auf unexportierte Teile* und Risiko *Prüfung geht im Umbau verloren*; `.claude/commands/implement-slice.md` Schritt 19 („Ein grüner Mutant wird in §7 eingeordnet“) | `internal/adapters/driving/pgwire/server.go` · „func (r *richtungen) weiterlesen() (geweckt bool)“; `docs/plan/planning/in-progress/slice-harness-blackbox-pgwire.md` · „*Umgeschrieben:* nur `TestWeiterlesenNachWecken`“ | ja (R1) | Prüfung im Umbau weggefallen, nicht eingeordnet |
| F-450 | INFO | Ist die erwartete Folge in `TestWeiterlesenNachWecken` deterministisch? Lesend ja, soweit es die Session betrifft. `wecke` und `weiterlesen` protokollieren unter derselben Sperre `frist`. Das erste `abgelaufen` kann erst nach dem gesetzten `frist` kommen. Das Signal in `weck` ist nach dem ersten `weiterlesen` verbraucht, also folgt genau ein zweites `abgelaufen` und dann `zurück`. Diese Folge kommt auch dann zustande, wenn `cancel()` vor dem nächsten `Receive` liegt. Eine verdeckte Abhängigkeit gibt es: Der Startphasen-Wächter in `handle` muss `fertig` schon gesehen haben, bevor `ctx` endet. Zwischen beiden besteht keine happens-before-Kante. Mit 20 ms Verzögerung im Wächter (R7) wurden 6 von 20 Läufen rot, mit `"frist abgelaufen abgelaufen zurück frist abgelaufen zurück"`. Ohne Eingriff liegen zwischen `close(fertig)` und `cancel()` mehrere Pipe-Rundläufe, und 600 Läufe unter CPU-Drossel waren grün. Ich erwarte keine Flakes. Für den Verifier. | Plan §6 Risiko *Prüfung geht im Umbau verloren*; Rolle Verifier | `internal/adapters/driving/pgwire/server_extended_test.go` · „got != \"frist abgelaufen abgelaufen zurück\"“ | ja (R7; 600 Läufe) | — (Hinweis) |
| F-451 | INFO | Die Einordnung des grünen Mutanten `meldeFrist = 3 * time.Second` (R5 bestätigt grün) folgt Schritt 19: Er ändert das Verhalten. Es gibt eine Test-Idee, und die Grenze steht in §7. Einen Wert für diese Frist nennt die Spezifikation tatsächlich nicht. Der Abschnitt *Record* sagt nur „Endet eine Session, schließt der Recorder die Client-Verbindung; ein blockiertes Schreiben an den Client endet damit“. Dass die Fehlerantwort vorher höchstens eine begrenzte Zeit geschrieben wird, bevor die Session an einem nicht lesenden Client endet, entscheidet nur der Code, und zwar schon seit `slice-extended-query-record` (V-22). Die Test-Idee „Schranke als Literal“ würde einen nicht spezifizierten Wert festschreiben. Sie braucht zuerst eine Stelle in der Spezifikation. Die Adresse, die die Closure vergibt, gehört deshalb dem Architect und nicht nur einem Test-Slice. | `AGENTS.md` §3.12; Rolle Architect (Spec-Lücke) | `internal/adapters/driving/pgwire/server.go` · „const meldeFrist = time.Second“ | ja (R5) | — (Hinweis) |
| F-452 | INFO | `sendeMu` → `nebenher(fe)` ist gleichwertig. Der ursprüngliche Kommentar aus `5d666db` lautete: „sendeMu hält den Schreibpuffer des Frontends“. Gemeint war also eine Sperre je Frontend. Jeder der 21 Aufrufe von `nebenher` legt in seinem Test oder Testfall die Funktion für genau ein `fe` an, und kein Test ruft `t.Parallel`. Die globale Sperre hat Verbindungen nur zusätzlich über Testgrenzen hinweg gekoppelt: Ein hängender Flush eines Tests hätte den nächsten Test blockiert. Der grüne Mutant ohne Sperre (R6 bestätigt grün) ist richtig eingeordnet. Ohne `-race` ist er nur über einen Test des Testhelfers selbst fangbar, und den schließt §1 aus („Neue Fälle“). Vorher war die Sperre genauso ungeprüft. | Plan §6 Risiko *Befund verdeckt statt behoben*; `SPEC-049` Punkt 5 (`gochecknoglobals`) | `internal/adapters/driving/pgwire/server_extended_test.go` · „func nebenher(fe *pgproto3.Frontend) func(msgs ...pgproto3.FrontendMessage)“ | ja (R6) | — (Hinweis) |
| F-453 | INFO | Zur Frage, ob die Mutationen für §3.10 reichen. Die Tabelle in §7 deckt jede Brücken-Funktion und jede Konstante ab. Für den umgeschriebenen Test nennt sie vier Mutanten. Meine Mutationen R2 bis R4 schließen Lücken: `weiterlesen` in `clientRichtung`, das Wecken über das Ende von `ctx` in `recordSitzung` und eine Frist in `wecke`, die nicht sofort abläuft. Alle drei sind rot. Grün sind nur R1 (F-449), R5 (F-451), R6 (F-452) und der Eingriff R7 (F-450). Nach meinem Urteil reicht die Abdeckung bis auf F-449. Ob sie die DoD trägt, prüft der Verifier. | `AGENTS.md` §3.10; Rolle Verifier | `docs/plan/planning/in-progress/slice-harness-blackbox-pgwire.md` · „*Mutationen* (`AGENTS.md` §3.10, Risiko *Prüfung geht im Umbau verloren*)“ | ja (R1 bis R7) | — (Hinweis) |

**Eigene Mutationen** (frische Kopie je Mutation; ohne Mutation grün):

| ID | Mutation | Weg im Test | Ergebnis |
|---|---|---|---|
| R1 | `weiterlesen` meldet nach dem Wecken `false` (`return true` → `return false`) | `Handle` | an `b37d6e7` grün (alle Tests des Pakets); an `b37d6e7^` rot, `TestWeiterlesenNachWecken` („Wecken nicht gemeldet“), siehe F-449 |
| R2 | `clientRichtung` ruft bei abgelaufener Lesefrist `weiterlesen` nicht | `Handle`, `fristSpion` | rot, `TestWeiterlesenNachWecken` („Lesefrist nicht zurückgesetzt“) |
| R3 | der Wächter in `recordSitzung` weckt beim Ende von `ctx` nicht | `Handle`, `fristSpion` | rot, `TestWeiterlesenNachWecken` (0 Ereignisse) |
| R4 | `wecke` setzt die Lesefrist eine Stunde in die Zukunft statt sofort | `Handle`, `fristSpion` | rot, `TestWeiterlesenNachWecken` (1 Ereignis) |
| R5 | `meldeFrist = 3 * time.Second` | Konstante `MeldeFrist` | grün (alle Tests des Pakets), wie §7 angibt, siehe F-451 |
| R6 | `nebenher` ohne `Lock` und `Unlock` (Testcode) | — | grün (`-count=5`, alle Tests des Pakets), wie §7 angibt, siehe F-452 |
| R7 | Eingriff in den Zeitablauf, keine Mutation des Verhaltens: der Startphasen-Wächter in `handle` beginnt 20 ms später | `Handle`, `fristSpion` | 6 von 20 rot, siehe F-450 |

## Antwort auf die Schwerpunkte

1. **Brücke nach Punkt 7.**
   - `export_test.go` gehört zum Paket `pgwire` und ist dort die einzige Testdatei. Sie enthält fünf Konstanten und fünf Funktionen. Es gibt keine Variable, keinen Typ-Alias, keine Funktion `Test…`/`Benchmark…`/`Example…`/`Fuzz…`, keinen Zustand auf Paketebene und keinen Wert eines unexportierten Typs.
   - `Handle`, `Fail` und `Note` rufen genau eine Methode des übergebenen `*Server` auf und reichen ihre Argumente unverändert weiter. `ToClientMessage` und `ToMessage` reichen ihr Argument unverändert weiter. Signaturen und Rückgaben bestehen nur aus exportierten Typen (`net.Conn`, `*pgproto3.Backend`, `model.Meldung`, `model.ClientMessage`, `pgproto3.BackendMessage`).
   - `richtungen` entsteht nur in `recordSitzung`. Das ist Produkt-Code und läuft über `Handle` → `handle`, denselben Weg, den `Serve` nimmt. `Handle` übernimmt dabei nicht `wg.Add`/`wg.Done` von `Serve`. Für die Tests spielt das keine Rolle, weil sie `Serve` auf diesem Weg nicht nutzen.
   - `MeldeFrist` ist eine Konstante. Schreiben ließe sie sich gar nicht, ohne dass der Compiler das ablehnt. Gelesen wird sie nur in `TestFehlerantwortMitFrist`.
   - `NewRecordServer(nil, log)` ersetzt `&Server{log: …}` feldgleich: `recorder` und `replayer` sind nil, `log` ist gesetzt.
2. **`TestWeiterlesenNachWecken`.**
   - Vorher stand der Test auf einem selbst angelegten `richtungen`-Literal mit `net.Pipe`, ohne Session. Er prüfte (a) die Rückgabe `true` nach dem Wecken, (b) dass das nächste Lesen an der Frist scheitert, (c) die Rückgabe `false` ohne Signal und (d) dass das Lesen danach nicht mehr an der Frist scheitert.
   - Der neue Test prüft (b) als erstes `abgelaufen`. (d) prüft er schärfer: als `zurück` und als `c:parse` nach dem Zurücksetzen. Dazu prüft er, dass das Signal verbraucht wird (genau ein zweites `abgelaufen`). Das ist in der Session zu sehen, die das Produkt erzeugt.
   - Weggefallen sind (a) und (c), also der Rückgabewert, den das Produkt nicht liest (F-449).
   - Zur Frage, ob die Folge deterministisch ist, siehe F-450.
3. **`sendeMu` → `nebenher(fe)`.** Gleichwertig, eine geteilte Sperre zwischen Verbindungen war nicht beabsichtigt, siehe F-452.
4. **Grüne Mutanten nach Schritt 19.** Beide sind richtig eingeordnet, und beide habe ich grün nachgestellt (R5, R6). Der Wert von `meldeFrist` ist eine Lücke in der Spezifikation, die schon vor diesem Slice bestand, siehe F-451. Nicht eingeordnet ist ein dritter grüner Mutant, R1 (F-449).
5. **Prüfung im Umbau verloren.**
   - Stichproben gegen `b37d6e7^`: `TestWeiterlesenNachWecken`, `TestFehlerantwortJeKlasse`, `TestNoteGleichrangig`, `TestFailGleichrangig`, `TestErsterFehlerZaehlt`, `TestFehlerantwortMitFrist`, `TestSSLUndGSSMitN`, `TestCancelRequest`, `TestUnbekannteSonderanfrage`, `TestToClientMessage`, `TestReplayExtended`, `TestExtendedEreignisse` und `TestReplayHerunterfahrenSpaeteFrist`.
   - Bis auf F-449 sind Erwartungen, Fristen und Abläufe gleich geblieben. Geändert sind nur Qualifizierung, Brücke, Parameterreihenfolge (`ctx` zuerst) und `sende := nebenher(fe)`.
   - Testliste gleich (40 Tests).
6. **Eigene Mutationen.** Siehe F-453 und die Tabelle oben.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driving/pgwire/export_test.go` | geprüft, ohne Befund nach `SPEC-049` Punkt 7. |
| `internal/adapters/driving/pgwire/server_test.go`, `server_meldung_test.go` | geprüft, ohne Befund. Paket `pgwire_test`, nur Qualifizierung und Brücke, kein Befund von `make lint`. |
| `internal/adapters/driving/pgwire/server_extended_test.go` | geprüft. Befund F-449, Hinweise F-450 und F-452. Keine Globale, `ctx` steht als erster Parameter. Kein neues `_ =` vor einem Fehler, den `errcheck` meldet: `_ = fe.Flush()` ist aus `sende` übernommen, und im neuen Test wird der Fehler von `SetDeadline` geprüft. |
| `internal/adapters/driving/pgwire/server.go` | geprüft, ohne Befund. Unverändert in `b37d6e7`. Zur Frist siehe F-451. |
| `docs/plan/planning/in-progress/slice-harness-blackbox-pgwire.md` | geprüft. Befund F-448 (§6), dazu F-449 (§7). §1 und §3 stimmen mit dem Diff überein. |
| Hard Rule 3.2, 3.6 (Suppression, Gate-Lockerung) | geprüft, ohne Befund. Kein `//nolint`, `.golangci.yml` ist unverändert. |
| Hard Rule 3.3 | geprüft, ohne Befund. `71e2a31` und `a4cc4da` sind reine Renames (100 %). `b37d6e7` enthält keinen Move. |
| Hard Rule 3.7, 3.11 (Kommentare) | geprüft, ohne Befund. Die Kommentare an der Brücke, an `nebenher`, an `fristSpion` und an `TestWeiterlesenNachWecken` beschreiben, was da ist. Der Satz zur Sperre in `nebenher` hat keinen Test, er stand aber schon so an `sende` (F-452). |
| Hard Rule 3.12 | geprüft, ohne Befund im Code-Commit. Neu in §6 ist keine Randform. Der Diff entscheidet keine Randform, die §6 nicht nennt. Zum Widerspruch in §6 siehe F-448. |
| Commit-Message `b37d6e7` | geprüft, ohne Befund. Sie nennt `slice-harness-blackbox-pgwire`, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und `LH-QA-07`, aber keine Struktur-ID. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 4 |

**Summary:** 0 HIGH · 1 MEDIUM (F-448: §6 nennt die Frage, ob die Brücke schreiben darf, weiter offen, obwohl Punkt *Export-Test-Brücke* und Diff sie beantworten) · 1 LOW (F-449: Rückgabewert von `weiterlesen` im Umbau ungeprüft, grüner Mutant nicht in §7 eingeordnet) · 4 INFO (F-450 Folge in `TestWeiterlesenNachWecken` praktisch deterministisch, eine verdeckte Abhängigkeit vom Startphasen-Wächter; F-451 Wert von `meldeFrist` ist eine Spec-Lücke aus dem Bestand; F-452 `nebenher` gleichwertig; F-453 Mutationen reichen bis auf F-449). Wiederkehrende Klasse: „Plan folgt Korrektur nicht“ (`BEO-REPO/plan-folgt-korrektur-nicht`, drittes Auftreten in Reviews).

**Finding-Klassen dieses Laufs:** Plan folgt Korrektur nicht · Prüfung im Umbau weggefallen, nicht eingeordnet

## Verdikt

**Merge-blockierend:** nein. Die Brücke hält Punkt 7 ein, `server.go` ist unverändert, und die Tests prüfen bis auf einen Rückgabewert ohne Produkt-Leser dasselbe wie vorher. F-448 und F-449 sind Nacharbeit am Plan, nicht am Code.

**Übergabe:** F-448 und F-449 gehen an den Implementer (§6 und §7 nachziehen). F-451 geht an den Architect, weil es eine Spec-Lücke zur Frist der Fehlerantwort am Session-Ende ist. F-450 und F-453 gehen an den Verifier. F-452 erwartet keine Aktion. Die Finding-Klassen gehen in §7 des Slice und von dort in den Zähler. Dieser Report ersetzt keine Verifikation.
