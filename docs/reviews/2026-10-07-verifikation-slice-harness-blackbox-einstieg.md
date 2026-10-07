# Verifikation: slice-harness-blackbox-einstieg — 2026-10-07

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft. F-456 bis F-458 hat er an mich übergeben, F-454 an den Implementer und F-455 an die Closure.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-blackbox-einstieg.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `f28a2df..0c1e9f7`. Darin liegen `1e4381b` (Umbau, Belege in §7) und `0c1e9f7` ([Review](2026-10-07-review-slice-harness-blackbox-einstieg.md), F-454 bis F-458). Der Go-Code an `0c1e9f7` ist gleich dem an `1e4381b`.

**Eingang:**

- DoD-Liefer-Punkte, §1, §3, §6 *Randformen* und *Risiken*, §7 *Belege des Implementers*, die Commit-Message von `1e4381b`
- `spec/spezifikation.md` `SPEC-049` Punkt 7 und 8 in der Fassung `f0ea7be` (`git diff f0ea7be 0c1e9f7 -- spec/` ist leer); [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5; [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md); [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- `.claude/commands/implement-slice.md` Schritt 19 (Einordnung grüner Mutanten)
- Produkt-Code `internal/adapters/driving/cli/cli.go` (Konstanten, Lesen der Umgebung), `internal/bootstrap/bootstrap.go` (`Run`, `logger`, `fail`), `internal/adapters/driven/recording/yaml.go` (`Unmarshal`, `fromDTO`), `internal/hexagon/services/lebendpruefung.go`, `internal/hexagon/services/record.go`
- der Review-Report und §1 von `slice-harness-mutation` (für die Adresse zu V-88)

Bei meinem Start war der Arbeitsbaum sauber, HEAD `0c1e9f7`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

**Wie die Proben gebaut wurden.** `f28a2df` (vor dem Umbau, `1e4381b^`) und `1e4381b` habe ich per `git archive` in je eine Kopie im Scratchpad entpackt. Testlisten, Unit-Tests und Unit-Mutationen liefen im Image der Stufe `deps` (eigener Tag, aus dem Arbeitsbaum gebaut), per Bind-Mount und mit `--network=none`. Die Lint-Läufe auf den Kopien liefen mit `docker build --progress=plain --target lint`. Jede Mutation lief in einer frischen Kopie (`cp -r` ohne `-p`). Das Ersetzungsskript bricht ab, wenn der Suchtext nicht genau einmal trifft, und prüft per Prüfsumme, dass sich die Datei geändert hat. Integrationsläufe: Testimage der Stufe `integration` aus der jeweiligen Kopie unter eigenem Tag, das gepinnte PostgreSQL-17-Image aus `harness/mk/integration.mk`, je Lauf ein eigenes internes Docker-Netz und ein eigenes Volume. Danach habe ich Container, Netze, Volumes und meine Images entfernt; hostweit habe ich nichts aufgeräumt. Das Repo selbst blieb unverändert. `make lint`, `make gates`, `make abdeckung-check` und `make kopf-check` liefen im Arbeitsbaum.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `_test`-Pakete, Testliste, Abdeckung, Brücke. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Testdateien gehören zu `_test`-Paketen | `cli_test.go` ist `package cli_test`, `bootstrap_test.go` ist `package bootstrap_test`, alle sieben Dateien unter `test/integration` sind `package integration_test`. Im Paket des Codes liegt je nur `export_test.go`. | bestätigt |
| Testliste vorher und nachher gleich | `go test -list .` an `f28a2df` und `1e4381b`, für `./test/integration` mit `-tags integration` und `PGR_OHNE_POSTGRES=1`: `cli` 19 und 19, `bootstrap` 12 und 12, `integration` 39 und 39, `diff` je Paket leer. `go test -count=1 -v` von `cli` und `bootstrap` liefert an beiden Ständen dieselben 31 Zeilen `--- PASS` (sortiert, ohne Zeitangabe), kein `FAIL`. | bestätigt |
| keine Prüfung entfallen | In `cli_test.go` und `bootstrap_test.go` habe ich Paketname, Präfix `cli.`/`bootstrap.` und die vier Brückennamen auf die alten Namen zurückgeführt: Der `diff` der beiden Stände zeigt dann nur die neue Import-Zeile des Pakets. In `test/integration` habe ich alle Hunks einzeln gelesen. Geändert sind nur der Paketname, die neun `Close`-Stellen, `lebendTexte`/`vtTexte`, `exec_` → `fuehreAus` und die SA4000-Bedingung. Erwartungen, Fristen und Abläufe sind gleich. | bestätigt |
| Abdeckungs-Deklarationen unverändert | An beiden Ständen 58 Deklarationen `// Abdeckung:` in den Testdateien der drei Pfade. Die vollständigen Kommentarblöcke samt Zuordnung Deklaration → nächstes `func Test…` (301 Zeilen, sortiert) sind gleich. | bestätigt |
| `make abdeckung-check` grün | Exit 0, einzeln und im Gate-Lauf (Abschnitt 5) | bestätigt |
| Unexportiertes nur über die Brücke | Der Compiler erzwingt es für die `_test`-Pakete. Die Brücken selbst: siehe Abschnitt 2. | bestätigt |

### Punkt 2: `make lint` ohne Befund in den Testdateien. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| 22 Befunde vorher in den Testdateien | Kopie `f28a2df`: `54 issues:`, 22 davon in `*_test.go`, alle unter den Pfaden dieses Slice. Je Paket `cli` 2 (1 in Testdateien), `bootstrap` 2 (1), `test/integration` 20 (20). Die 22 Zeilen sind sortiert byte-gleich mit dem Block *Vorher* in §7. | bestätigt |
| Aufteilung 9 + 13 (§1, §3) | `testpackage` 9 (1 + 1 + 7). Dazu `errcheck` 9, `gochecknoglobals` 2, `revive` 1 und `staticcheck` 1, zusammen 13. | bestätigt |
| 0 nachher | Kopie `1e4381b` und `make lint` im Arbeitsbaum (Exit 2): `32 issues:`, **keiner in einer `*_test.go` des ganzen Moduls**, keine `lint:`-Zeile der eigenen Prüfungen, auch nicht an den beiden `export_test.go`. | bestätigt |
| modulweit 54 → 32, kein neuer Befund | Sortierte Befundzeilen mit Zeile und Spalte: `comm -13` (nur nachher) ist leer, `comm -23` (nur vorher) sind genau die 22 Zeilen der Testdateien. | bestätigt |
| Zeilen vorher und nachher im Bericht | §7 nennt vorher die 22 Zeilen der Testdateien und die Zahl der Produkt-Befunde je Paket, nachher die 2 Produkt-Zeilen (`bootstrap.go:71:26` `contextcheck`, `cli.go:243:7` `revive` `unused-receiver`) mit dem Satz, es seien dieselben wie vorher. Die Gleichheit habe ich mit Zeile und Spalte gemessen. | bestätigt |
| ohne `//nolint` | `grep -rn nolint` über `internal/`, `cmd/` und `test/` an `1e4381b` ist leer. | bestätigt |
| ohne Änderung an `.golangci.yml` | `git diff f28a2df 0c1e9f7 -- .golangci.yml tools/harness/lint.sh Dockerfile .dockerignore` ist leer. | bestätigt |
| kein `_ =` vor einem gemeldeten Fehler | Alle Blank-Zuweisungen in den Testdateien der drei Pfade an beiden Ständen (Paketpräfix entfernt, sortiert): je 30, der `diff` ist leer. Die neun `Close`-Fehler gehen jetzt an `t.Error`. | bestätigt |

**Vorbedingung für `slice-lint-bestand-*` und `slice-harness-lint`, ausdrücklich bestätigt:** `make lint` meldet an `0c1e9f7` in den Testdateien des ganzen Moduls **0 Befunde**. Alle 32 Befunde liegen in Produkt-Dateien. Von den Testdateien, die der Slice umgeschrieben hat, und den beiden Brücken meldet auch `tools/harness/lint.sh` nichts.

### Punkt 3: `make gates` grün. **Bestätigt.**

`make gates` im Arbeitsbaum an `0c1e9f7`, Exit 0 (Abschnitt 5). Für diesen Punkt trägt §7 nur die Zeile „`make gates` grün am Stand dieses Commits“, ohne Ausgabe. Den Lauf habe ich deshalb selbst gefahren.

### Kein Produkt-Code geändert. **Bestätigt.**

`git diff --name-only f28a2df 1e4381b`: der Plan, die neun Testdateien und die zwei neuen `export_test.go`. `cli.go`, `bootstrap.go`, `doc.go` und alle übrigen Produkt-Dateien sind unverändert, ebenso `Dockerfile`. Die Stufe `integration` baut `package integration_test` (Gate-Lauf und eigene Images).

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| Review mit Report | liegt vor (F-454 bis F-458) |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen: §7 trägt Platzhalter, §6 *Risiken* hat drei Ausgänge „— (bei Closure)“. Das ist Sache der Closure; Vorschläge stehen in Abschnitt 6. |

---

## 2. Brücken nach `SPEC-049` Punkt 7 (Fassung `f0ea7be`)

| Bedingung aus Punkt 7 | `cli/export_test.go` | `bootstrap/export_test.go` |
|---|---|---|
| gehört zum Paket des Codes, einzige Testdatei dort | `package cli`; daneben `cli.go`, `doc.go` und `cli_test.go` (`cli_test`) | `package bootstrap`; daneben `bootstrap.go`, `doc.go` und `bootstrap_test.go` (`bootstrap_test`) |
| nur Typ-Aliase, Konstanten, weiterreichende Funktionen | zwei Konstanten, `EnvFailOnUnconsumed = envFailOnUnconsumed`, `EnvLogLevel = envLogLevel` | zwei Funktionen, `Logger` und `Fail`, je ein Aufruf von `logger` bzw. `fail` mit unveränderten Argumenten |
| kein Zustand auf Paketebene, keine Variable | keine | keine |
| keine Funktion `Test…`/`Benchmark…`/`Example…`/`Fuzz…` | keine | keine |
| kein Wert eines unexportierten Typs, weder in der Brücke noch im Test | Konstanten vom Typ `string` | Signaturen nur mit `io.Writer`, `string`, `error`, `*slog.Logger` und `int` |
| eine erzeugende Funktion ruft auch das Produkt | — | `logger` ruft `Run` zweimal (`bootstrap.go` Zeile 41 und 43), `fail` ruft es siebenmal |

Hält. Die Zuordnung in §7 (vier Zugriffe, alle über die Brücke, keiner gegen die Schnittstelle umgeschrieben) stimmt mit dem normierten `diff` aus Abschnitt 1 überein: Außer der Qualifizierung hat sich in den beiden Unit-Testdateien nichts geändert.

---

## 3. Mutationen nach `AGENTS.md` §3.10, Stichprobe gegen §7

Der Grundlauf an `1e4381b` ist grün: Unit-Tests von `cli` und `bootstrap` mit `ok` (U0), Integration ohne Mutation 37 bestanden und 2 übersprungen, Exit 0 (E0).

| ID | Zeile in §7 | Mutation (wie ausgeführt) | Ergebnis |
|---|---|---|---|
| V1 | eine Stufe zeigt ihre und die strengeren Zeilen | `Level: stufen[stufe]` → `Level: slog.LevelInfo` | rot, `TestLoggerSchwelle`, dazu `TestRunLogLevel`, `TestRunRecordLogLevel`, `TestRunRecordSchreibfehlerJeStufe` |
| V2 | je Meldung eine Zeile beim Prozessende | `fail` schreibt nur `ms[:1]` | rot, `TestFailJeKlasse` |
| V3 | Exit-Code der ersten Meldung | `fail` liefert `ms[len(ms)-1].ExitCode()` | rot, `TestFailJeKlasse` |
| V4 | Name `PGWIRE_RECORDER_LOG_LEVEL` | `envLogLevel = "PGWIRE_RECORDER_LOGLEVEL"` | rot, `TestParseLogLevelHilfe` (`cli`), `TestRunLogLevel`, `TestRunStartfehlerJeStufe` (`bootstrap`) |
| V5 | Name `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` (*Ohne roten Test*) | `envFailOnUnconsumed = "PGWIRE_RECORDER_FAILONUNCONSUMED"` | grün in allen Tests von `cli` und `bootstrap`, wie §7 angibt. Rot im Integrationstest hat das Review gezeigt (R1i); nicht wiederholt. |
| V6 | `PGWIRE_RECORDER_LOG_LEVEL` wird gelesen | `logLevelOption` liest die Variable nicht | rot, `TestParseLogLevelUmgebung`, dazu `TestParseLogLevelUmgebungUngueltig`, `TestRunLogLevel`, `TestRunStartfehlerJeStufe` |
| V7 | `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` wird gelesen | `Parse` liest die Variable nicht | rot, `TestParseFailOnUnconsumedUmgebung`, `TestParseFailOnUnconsumedUmgebungNebenOption` |
| P1 | beschädigte Aufzeichnung ist PGR-E3003 (über `fuehreAus`) | `fromDTO`: „Liste der Sessions fehlt“ mit `CodeRecordingIO` (PGR-E3001) | rot, `TestE2EReplayBeschaedigt` („Exit-Code 3“, Ausgabe ohne PGR-E3003), dazu im Paket `recording` `TestUnmarshalFehler` |
| P2 | *Ohne roten Test*, PGR-E3003 an `formVorpruefung` | Rückgabe „Aufzeichnung beschädigt“ nach `formVorpruefung` mit `CodeRecordingIO` | grün, `TestE2EReplayBeschaedigt`; rot im Paket `recording`, `TestUnmarshalExtendedFehler` und `TestVorpruefungNenntOrt`, wie §7 angibt |
| I1 | verschachtelte Blockkommentare (`lebendTexte()`) | `/*` öffnet nur auf Tiefe 0 eine Ebene | rot, `TestE2EReplayLebendpruefungWiePostgres` (PGR-E5001 an `"/* a /* b */ c */"`); `…Pgxpool` und `…DatabaseSQL` grün, wie §7 angibt |
| I2 | nach SIGTERM keine neue Interaktion beim Pipelinen | `ClientMessage` ohne die Sperre `herunterfahren` | rot, `TestE2ERecordExtendedSigtermBeimPipelining` („Recorder endet nicht binnen 5 s nach SIGTERM“); danach hängt der Test in zwei von fünf Läufen bis zum Zeitlimit, siehe V-88 |

Neun Zeilen der Tabelle in §7 habe ich nachgefahren (V1 bis V4, V6, V7, P1, I1, I2), jede ist aus dem genannten Grund rot. Die zwei Einträge unter *Ohne roten Test*, die das Produkt betreffen, habe ich nachgestellt (V5, P2); beide sind richtig eingeordnet. Zusammen mit R1 bis R3 und R1i des Reviews sehe ich keine Prüfung, die im Umbau verloren ging.

---

## 4. F-456 bis F-458

**F-456: bestätigt.** V5 ist in `cli` und `bootstrap` grün. Der Bestand verhält sich gleich: Vorher nannten die Tests in `cli` dieselbe Konstante direkt, also fing auch der alte Test den Mutanten nicht. Ein Verlust durch den Umbau ist das nicht. Die Asymmetrie zu `--log-level` (V4 rot über `TestParseLogLevelHilfe`) stammt aus dem Bestand. Ein neuer Fall wäre nach §1 ein anderer Vorgang. Einen weiteren Beleg für einen Register-Eintrag schlage ich nicht vor: Der Integrationstest hält den Namen als Literal fest.

**F-457: bestätigt, die Formulierung ist ungenau.** Gemessen (P1, P2): Die Eingabe von `TestE2EReplayBeschaedigt` (`format`, `version`, keine `sessions`) **besteht** `formVorpruefung` ohne Fehler und scheitert erst in `fromDTO` an „Liste der Sessions fehlt“. P1 an dieser Stelle ist rot, P2 an der Rückgabe nach `formVorpruefung` ist grün. Die Fakten in §7 stimmen, der Satz trägt sie an drei Stellen ungenau:

1. „derselben Mutation“: P2 ist ein anderer Mutant als die Tabellenzeile. Er sitzt an einer anderen Stelle mit anderer Meldung („Aufzeichnung beschädigt“), gleich ist nur die Ersetzung PGR-E3003 → PGR-E3001.
2. „Der erste Platz“: Diese Reihenfolge nennt §7 nirgends. In `Unmarshal` ist die Rückgabe nach `formVorpruefung` der fünfte Platz mit `CodeRecordingBroken`, der erste ist der Fehler von `yaml.Unmarshal`.
3. „erreicht die Stelle nicht“: Die Eingabe erreicht `formVorpruefung`. Sie nimmt nur den Fehlerzweig dort nicht.

Eine genauere Fassung für die Closure: „Dieselbe Ersetzung PGR-E3003 → PGR-E3001 an der Rückgabe „Aufzeichnung beschädigt“ nach `formVorpruefung` blieb in `TestE2EReplayBeschaedigt` grün. Die Eingabe besteht `formVorpruefung` und scheitert erst in `fromDTO` (Zeile oben, rot). Die Rückgabe selbst fangen `TestUnmarshalExtendedFehler` und `TestVorpruefungNenntOrt` im Paket `recording`.“ Das ist ein Hinweis, kein DoD-Mangel.

**F-458: bestätigt.** An allen neun Stellen steht der `defer` mit `Close` nach dem `defer cancel()`. Nach LIFO schließt die Verbindung also wie vorher vor dem `cancel`. Im Grundlauf E0 und im Gate-Lauf kam kein `Close`-Fehler. Ob die Meldung bei einem Fehler wirkt, prüft kein Test. §7 nennt das als Grenze, und es ist eine Zusage am Testgeschirr, nicht am Produkt. Ein weiterer Befund aus dem Bestand an demselben Test steht in V-88.

---

## 5. Gate-Ergebnis

- `make gates` an `0c1e9f7` im Arbeitsbaum: **grün** (Exit 0). Im Log grün unter anderem `abdeckung-gegenprobe`, `a-check-negativ`, `commit-msg-gegenprobe`, `run-integration-tests` (39 `--- PASS`: 37 in der ersten Phase, 2 in der zweiten; 2 `--- SKIP`) und `kopf-check-gegenprobe`, kein `--- FAIL`. Das deckt sich mit den Zahlen in §7.
- `make abdeckung-check` und `make kopf-check` einzeln: Exit 0.
- `make lint` im Arbeitsbaum: Exit 2, `32 issues:`, keiner in einer Testdatei.
- `gofmt -l` der drei Pfade leer, `go vet -tags integration` der drei Pakete grün (Kopien `f28a2df` und `1e4381b`).

---

## 6. Plan gegen Code

- **Kopf:** `Bezug` ([`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md)) und `Berührte Spec-Stellen` (`SPEC-049` Punkt 7 und 8, unverändert) stimmen. `make kopf-check` ist grün.
- **§1:** Das Ziel ist erreicht (Abschnitt 1). Die Zerlegung 9 + 13 = 22 stimmt. Die Abgrenzung hält: kein Produkt-Code, kein neuer Export im Produkt, keine neue Testfunktion, `.golangci.yml` und die Lint-Gegenprobe unverändert. Die zwei Produkt-Befunde bleiben für `slice-lint-bestand-driving`.
- **§3:** Jede Zeile hat ihre Datei im Diff und umgekehrt. `export_test.go` gibt es zweimal, in `cli` und `bootstrap`, nicht in `test/integration`. `Dockerfile` ist unverändert, wie die Zeile sagt.
- **§4:** Keine Rückführung ausgelöst. Kein Test ist umgeschrieben, alle 70 sind umgestellt. Die Brücke reichte, ohne dass der Architect entscheiden musste.
- **§5:** `make gates` grün und `make lint` ohne Befund in den Testdateien sind gemessen.
- **§6 *Randformen*:** `1e4381b` entscheidet keine Randform, die §6 nicht nennt. Alle vier White-Box-Zugriffe fallen in die Klassen des Punkts *Export-Test-Brücke*. Dass die Liste erst im Code-Commit per AST gemessen wurde, hat das Review als F-455 an die Closure gegeben. Dem folge ich.
- **§6 *Risiken*, Vorschläge für die Ausgänge** (setzen muss sie die Closure):
  - *Prüfung geht im Umbau verloren:* **entfallen**. Begründung: Die Unit-Testdateien unterscheiden sich nach Normierung nur in der Import-Zeile, die Testliste ist gleich, und alle nachgefahrenen Mutationen sind rot. Die grünen Mutanten V5 und P2 waren auch vorher grün.
  - *Abdeckung sinkt:* **entfallen** für diese Pfade. Begründung: Die Testkörper sind unverändert und rufen über die Brücke dieselben Funktionen, also kann der erreichte Code nicht sinken. Die Messung bleibt Gegenstand von `slice-harness-coverage`, und die gilt für alle vier Umstellungs-Slices.
  - *Befund verdeckt statt behoben:* **entfallen**. Begründung: Je 30 Blank-Zuweisungen an beiden Ständen, gleich. Kein `//nolint`. Die `Close`-Fehler gehen an `t.Error`. `lebendTexte()` und `vtTexte()` liefern je Aufruf ein neues Slice-Literal.
- **§7 *Belege*:** Vollständig für das, was die DoD verlangt: White-Box-Zugriffe, Testliste, Deklarationen, `make abdeckung-check`, Lint-Zeilen vorher und nachher, Mutationen. `make gates` steht nur als Zeile mit Stand, den Lauf habe ich selbst gefahren. Die Zahlen stimmen. Die Formulierung zu PGR-E3003 ist ungenau (F-457, Abschnitt 4).
- **§8:** Die Sichtung nennt den Stand `175fd2d`. Die Treffer sind plausibel, keine Schwelle ist erreicht.

---

## 7. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-88 | Hinweis (Bestand, Testgeschirr) | Wird `TestE2ERecordExtendedSigtermBeimPipelining` rot, weil der Recorder nach SIGTERM nicht endet, hängt der Test danach in einem Teil der Läufe bis zum Zeitlimit von `go test`, statt nach etwa 6 s zu scheitern. Ursache: Die Goroutine `beendet <- rec.cmd.Wait()` und der `t.Cleanup` aus `startProzess` (`Kill`, dann `cmd.Wait()`) rufen `Wait` auf demselben `exec.Cmd` nebeneinander. Den Kanal der Kopier-Goroutinen liest nur einer der beiden, der andere wartet in `awaitGoroutines` für immer (Stack im Lauf). Unter I2, je 40 s Zeitlimit: an `1e4381b` 2 von 5 Läufen hängend, an `f28a2df` 1 von 5. Der Fehler kommt also **aus dem Bestand**, der Slice hat den `Cleanup` nicht geändert. Im Grünfall wirkt er nicht, er verteuert nur einen roten Lauf (in `make test-integration` gilt das Standardlimit von 10 min). Dasselbe Muster, eine Goroutine mit `cmd.Wait()` neben dem `Cleanup`, steht in `record_e2e_test.go` (`done <- r.cmd.Wait()`) und `extended_replay_e2e_test.go` (`done <- rep.cmd.Wait()`). Kein DoD-Mangel. Adresse für den Planner: **nicht** `slice-harness-mutation`, sein §1 schließt Mutationen in den Integrationstests aus. Vorschlag: eine Zeile in §7 dieses Slice unter *Folge-Slices*, „offen, ohne Adresse“, oder ein Eintrag im Beobachtungs-Register. Ob ein eigener kleiner Slice am Testgeschirr lohnt (ein `Wait` je Prozess, das Ergebnis über einen Kanal geteilt), entscheidet der Planner. | `test/integration/record_e2e_test.go` · „_ = cmd.Wait()“; `test/integration/extended_e2e_test.go` · „go func() { beendet <- rec.cmd.Wait() }()“ | I2 an beiden Ständen, je 5 Läufe |

Zu F-457 siehe Abschnitt 4 (bestätigt, keine eigene V-Nummer).

---

## 8. Urteil

Die drei Liefer-Punkte sind erfüllt und selbst belegt:

- Die Tests laufen als `cli_test`, `bootstrap_test` und `integration_test`. Die Testliste ist gleich (`cli` 19, `bootstrap` 12, `integration` 39), dazu dieselben 31 PASS-Zeilen der Unit-Tests an beiden Ständen.
- Die 58 Abdeckungs-Deklarationen sind gleich, Wort für Wort und in der Zuordnung. `make abdeckung-check` ist grün.
- `make lint`: in den Testdateien 22 → 0, modulweit 54 → 32, kein neuer Befund, kein `//nolint`, kein neues `_ =`, `.golangci.yml` unverändert. **In keiner Testdatei des Moduls steht noch ein Befund.** Die Vorbedingung für `slice-lint-bestand-*` und `slice-harness-lint` ist erfüllt.
- Kein Produkt-Code geändert. Die Brücken halten `SPEC-049` Punkt 7 in der Fassung `f0ea7be`.
- Neun nachgefahrene Mutationen aus §7 sind aus dem genannten Grund rot. Die grünen Mutanten sind nach Schritt 19 tragfähig eingeordnet.
- `make gates` ist grün.
- F-456 und F-458 bestätigt. F-457 bestätigt: Die Fakten stimmen, die Formulierung ist ungenau.

Kein Befund blockiert die Closure. V-88 ist ein Hinweis aus dem Bestand für Folge-Slices und Register. Abschnitt 6 enthält Vorschläge für die Risiko-Ausgänge.
