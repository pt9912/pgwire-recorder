# Verifikation: slice-harness-integration-wait — 2026-10-08

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft ([Review](2026-10-08-review-slice-harness-integration-wait.md), F-480 bis F-484). F-483 und F-484 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-integration-wait.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `86db551..9fec9e6`. Architect: `2ebe058` (Absatz *Warten in Tests* in `SPEC-038`, §6 *Randformen*), `2a6cd91` (F-480 bis F-482). Implementer: `6deea56` (Umbau), `9fec9e6` (Nacharbeit). Review: `d3e456b`.

**Eingang:**

- Plan ganz, besonders §1, §2 (DoD), §3, §6 (*Randformen*, *Mutationen* mit der Vorgabe nach F-481, *Risiken*) und §7 (*Belege des Implementers* mit *Nacharbeit nach dem Review* und *Grüne Mutanten, eingeordnet*)
- `spec/spezifikation.md` `SPEC-038`, Absatz *Warten in Tests*; [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md)
- `AGENTS.md` §3.9 bis §3.12
- Am Stand `9fec9e6`: `test/integration/record_e2e_test.go` (Typ `recorder`, `startProzess`, `stop`, `signal`, `warteEnde`, `nachKill`, `pruefeExit`), die beiden SIGTERM-Tests, `harness/mk/integration.mk`, `tools/test/run-integration-tests.sh`, die Stufe `integration` im `Dockerfile`, der Block `herunterfahren` in `ClientMessage` (`internal/hexagon/services/record.go`)
- `docs/plan/planning/next/slice-v1-abschluss-herunterfahren.md` §1, §2, §3, §6 (Erwartungen an die Helfer)
- Der Review-Report ganz, den Bericht des Implementers nicht

Bei meinem Start war der Arbeitsbaum sauber, HEAD `9fec9e6`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-08

**Wie die Proben gebaut wurden.** Je Probe eine frische Kopie per `git archive <stand> | tar -x -m` in einem eigenen Unterverzeichnis des Scratchpads, nie im Repo. Die Mutation per Skript, das abbricht, wenn der Suchtext nicht trifft (I2: genau einmal), danach `touch` auf die geänderte Datei. Je Probe ein eigenes Image-Tag der Stufe `integration`; je Lauf ein eigenes `--internal`-Netz, ein eigener PostgreSQL-Container (Image und Digest aus `harness/mk/integration.mk`) und ein eigenes Volume, nach jedem Lauf entfernt, nach der Probe Image und Kopie entfernt. Danach standen keine Container, Netze, Volumes oder Images mit dem Lauf-Präfix mehr. Gestartet gezielt mit `-test.run '^<Test>$' -test.timeout 60s`. `make gates`, `make lint` und `make docs-check` liefen im Arbeitsbaum; `git status` war danach sauber.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Ein `Wait` je Prozess. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| vorher vier Stellen | `git grep -n "\.Wait()" 2ebe058 -- test/integration`: `extended_e2e_test.go:410`, `extended_replay_e2e_test.go:297`, `record_e2e_test.go:373`, `record_e2e_test.go:396` — wie §1 *Bestand* und §7 | bestätigt |
| nachher genau eine | an `9fec9e6`: nur `record_e2e_test.go:381` (`r.waitErr = cmd.Wait()` in der Goroutine aus `startProzess`) — wie §7 | bestätigt |
| `stop`, beide Tests und Cleanup lesen das Ergebnis | gelesen: `stop` → `signal`, `warteEnde`, `pruefeExit`; beide SIGTERM-Tests → `signal`, `warteEnde`, `pruefeExit`; der Cleanup fragt `beendet` ab, ruft sonst `Kill` und wartet höchstens 5 s auf `beendet` | bestätigt |
| Leser erst nach dem Schließen (§6 *Ergebnis nach dem Ende*, *Mehrere Leser*) | `waitErr` und `ProcessState` schreibt die Goroutine vor `close(r.beendet)`. Gelesen werden sie nur in `signal` im Zweig `<-r.beendet`, in `pruefeExit` (an allen drei Aufrufstellen nach `warteEnde`). Kein `sync.Once`, kein zweiter Wert. Wie F-483. | bestätigt |

### Punkt 2: Rot statt hängend, Mutationen aus §6. **Bestätigt** (mit V-104).

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| vorher mindestens ein Hänger | §7: `2ebe058`, 1 von 10 (Lauf 8). Eigene Messung `2ebe058` mit I2: **2 von 10 hängend** (Läufe 1 und 9, `panic: test timed out after 1m0s`, 60,6 s), 8 rot nach 6,11–6,12 s Testdauer (Lauf gesamt 6,6–7,4 s). Der Stack des Hängers steht in `os/exec.(*Cmd).awaitGoroutines` unter `record_e2e_test.go:373` (Cleanup), vorher die Meldung des Tests (V-88). Mit Review (2 von 20) sind es 5 von 40. | bestätigt |
| nachher in jedem von mindestens zehn Läufen rot mit der Meldung, weit vor dem Zeitlimit | §7: zehn Läufe am Arbeitsbaum des Umbaus, nach der Nacharbeit ein Lauf (V-104). Eigene Messung **am Endstand `9fec9e6`: 10 von 10 rot, 0 hängend**, Testdauer 6,11–6,12 s, Lauf gesamt 6,6–9,2 s, Zeitlimit 60 s. Meldung `extended_e2e_test.go:407: Recorder endet nicht binnen 5 s nach SIGTERM, obwohl der Client weiter pipelinet`, darunter `stderr` des Recorders. Mit Review (0 von 30 an `6deea56`) sind es 0 von 50 nach dem Umbau. | bestätigt |
| Mutation *Cleanup* (Vorgabe nach F-481) | §7 *Nacharbeit*: Lauf 1 FAIL 5,11 s mit der Meldung des Cleanups, Lauf 2 (`t.Errorf` → `t.Logf`) PASS, Lauf 3 (ohne Frist) Abbruch an der Zeitgrenze. Status, Dauer und Meldung je Lauf vorhanden. Den Gegenlauf zu Lauf 1 mit derselben Provokation und **mit** `Kill` führt §7 nur in der älteren Form mit `t.Fatal`; eigener Lauf: `if true { return }` nach `startRecorder`, Cleanup unverändert → **PASS nach 0,10 s**, keine Meldung des Cleanups. Damit unterscheidet Lauf 1 auch den Mutanten *`Kill` entfernt*. | bestätigt |
| Mutation *Signal an einen beendeten Prozess* (F-480) | §7: FAIL 1,11 s mit „Prozess endete vor dem Signal …“ und `stderr`; Mutant FAIL nur mit `os: process already finished`. **Nachgefahren**, beides gleich: Provokation in `TestE2ERecordSelect1` → FAIL 1,11 s, `record_e2e_test.go:75: Prozess endete vor dem Signal terminated (signal: killed, Wait: signal: killed)`, darunter `stderr` („record gestartet“). Mutant (Warten auf `beendet` in `signal` entfernt) → FAIL 1,11 s, `record_e2e_test.go:75: Signal terminated: os: process already finished`, ohne Zustand und `stderr`. | bestätigt |
| Mutation *Lausch-Frist* (F-482) | §7: FAIL 10,05 s mit „Recorder lauscht nicht …“ und `stderr`, Prozess beendet. Nicht nachgefahren; der Code führt den Zweig über `nachKill`, der auch unter I2 läuft (gesehen). Dass `stderr` nicht **vor** dem Ende gelesen wird, fängt ohne `-race` kein Lauf — gelesen: `nachKill` liest `stderr` nur im Zweig `<-r.beendet`. | bestätigt |
| Beleg-Form (Quellstand, Läufe, rot, hängend, Dauer, Zeitlimit) | Tabelle in §7 mit allen Spalten; „Dauer je Lauf“ als Spanne, Hänger mit Laufnummer | bestätigt |

### Punkt 3: Testliste und Deklarationen unverändert. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| dieselben Tests | `git grep -h "^func Test"` unter `test/integration` an `2ebe058` und `9fec9e6`: je 40, sortiert `diff` leer | bestätigt |
| keine Zeile `Abdeckung:` geändert | `git diff 2ebe058 9fec9e6 -U0 -- test/integration \| grep -c Abdeckung` → 0 | bestätigt |
| `make abdeckung-check` grün ohne neu geschriebene Tabellen | Teil von `make gates` (Exit 0), `git status` danach sauber | bestätigt |

### Punkt 4: `make gates` grün. **Bestätigt** (Abschnitt 4).

### Übrige DoD-Punkte (konstant je Slice)

- Review: liegt vor (`d3e456b`), F-480 bis F-482 entschieden (`2a6cd91`) und umgesetzt (`9fec9e6`, Abschnitt 2).
- Closure-Notiz, Register, Risiko-Ausgänge, Paarungen: in §7 noch Platzhalter, §6 *Risiken* ohne Ausgang. Das ist Closure-Arbeit nach dieser Verifikation, kein Befund. Vorschläge für die Ausgänge in Abschnitt 5.

## 2. Review-Befunde F-480 bis F-484: nachgefahren

| Finding | Umsetzung | Urteil |
|---|---|---|
| F-480 — Signal an einen schon beendeten Prozess | §6 *Prozess ist beim Signal schon beendet* entschieden vor dem Code (`2a6cd91`). `Process.Signal` steht unter `test/integration` nur in `signal` (`record_e2e_test.go:424`), `stop` und die beiden Tests rufen `signal`. Liefert `Signal` einen Fehler, höchstens 5 s auf `beendet`, dann „Prozess endete vor dem Signal“ mit `ProcessState.String()`, Fehler von `Wait` und `stderr`, sonst der Fehler von `Signal`. Keine Vorab-Abfrage. Provokation und Mutant nachgefahren (Abschnitt 1). Zum Ablauf-Zweig siehe V-103. | umgesetzt |
| F-481 — Einordnung der Mutanten | drei Cleanup-Läufe gefahren, *Cleanup wartet ohne Frist* nicht mehr unter *äquivalent*, `t.Errorf` → `t.Logf` aufgenommen. Die Einordnung ist uneinheitlich (V-102). | umgesetzt, V-102 |
| F-482 — `stderr` vor dem Ende gelesen | Lausch-Frist über `t.Fatal(r.nachKill(…))`; `nachKill` ruft `Kill`, wartet höchstens 5 s auf `beendet` und liest `stderr` nur danach; sonst „endet auch nach Kill nicht binnen 5 s“ ohne `stderr`. `warteEnde` nutzt dieselbe Funktion. Kein Leser im Geschirr liest `stderr` vor dem Schließen; die Leser in den Tests (etwa `record_e2e_test.go:308`) lesen nach `stop` oder bleiben als akzeptiertes Negativ nach §6. | umgesetzt |
| F-483 — happens-before, ein `Wait`, gleiche Prüfungen | am Endstand wiederholt (Abschnitt 1, Punkt 1; Abschnitt 3) | bestätigt |
| F-484 — Messung trägt schwach | Die eigenen 20 Läufe ergänzen: vorher zusammen 5 von 40 hängend, nachher 0 von 50. Der strukturelle Beleg (ein `Wait`) trägt wie im Review. Die Rückführung aus §4 ist nicht eingetreten. | bestätigt |

## 3. `SPEC-038` *Warten in Tests* im neuen Code

| Regel | Code an `9fec9e6` | Urteil |
|---|---|---|
| eigene Frist, als Literal im Test | `stop` `15*time.Second`, Pipelining-Test `5*time.Second`, Replay-Test `10*time.Second` je am Aufruf von `warteEnde`; `signal`, `nachKill` und Cleanup je `5 * time.Second`; Lausch-Frist `10 * time.Second`. `warteEnde` nimmt die Frist als Parameter, der Literal steht am Aufruf. | eingehalten |
| höchstens 60 s | längste Kette `stop`: 15 s + Nachfrist 5 s, danach der Cleanup höchstens 5 s — 25 s | eingehalten |
| rot statt Hängen, die Meldung nennt das ausgebliebene Ereignis | `warteEnde` mit dem Text des Aufrufers („endet nicht …“), `nachKill` ergänzt „endet auch nach Kill nicht binnen 5 s“, Cleanup „endet im Cleanup auch nach Kill nicht binnen 5 s“. **Ausnahme:** der Ablauf-Zweig von `signal` meldet nur „Signal %v: %v“ (V-103). | eingehalten bis auf V-103 |
| ein `Wait` je Prozess, weitere Leser lesen das Ergebnis | Abschnitt 1, Punkt 1 | eingehalten |

## 4. Gate-Ergebnis

- `make gates` im Arbeitsbaum an `9fec9e6`: **Exit 0**. Darin u. a. `run-integration-tests: gruen` (beide Phasen), `d-check: … 0 Befund(e)`, `a-check-negativ: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `abdeckung-gegenprobe: gruen`, `abdeckung-check` ohne Ausgabe und ohne geänderte Tabellen.
- `make lint`: **Exit 0**. Die Stufe kam aus dem Cache (`#14 CACHED`), also aus einem früheren Lauf derselben Quellen, der grün war; ein roter `RUN` wird nicht gecacht. Eine Zeile `0 issues.` steht darum in meinem Log nicht.
- `make docs-check` vor dem Commit dieses Berichts: grün.

## 5. Plan gegen Code

| Plan-Stelle | Code | Urteil |
|---|---|---|
| §1 Ziel, Abgrenzung „kein Produkt-Code“ | `git diff --name-only 86db551 9fec9e6`: nur der Plan, `slice-tests-ueberlebende-mutanten-driving`, der Review-Report, `spec/spezifikation.md` und die drei Testdateien unter `test/integration`. I2 nur in Kopien. | konform |
| §1 Fristen der drei Tests bleiben (15 s, 5 s, 10 s), Meldungen gleich | wörtlich gleich (Diff gelesen) | konform |
| §1 `defer pc.Close` bleibt | unverändert; hängt in keinem meiner 20 I2-Läufe | konform |
| §3 Tabelle | `recorder` mit `beendet` und `waitErr`, `startProzess` mit einer Goroutine, `warteEnde`, `pruefeExit`, `signal`, Lausch-Frist über `nachKill`. `nachKill` nennt §3 nicht beim Namen, §6 *`stderr` erst nach dem Ende lesen* erlaubt die gemeinsame Funktion ausdrücklich. | konform |
| §6 *Frist und Ablauf* | Kill, Nachfrist 5 s, `t.Fatal` mit Text und `stderr` bzw. Hinweis | konform |
| §6 *Test scheitert vorher* | Cleanup nach `Start` und Goroutine registriert, bei geschlossenem Kanal nichts, sonst `Kill` und 5 s, `t.Errorf` | konform |
| §6 *Prozess ist beim Signal schon beendet* | wie entschieden; der Ablauf-Zweig siehe V-103 | konform |
| §6 *Risiken* — Vorschläge für die Ausgänge bei Closure | *Messung unterscheidet nicht*: entfallen — vorher 5 von 40 hängend (§7, Review, Verifikation), nachher 0 von 50. *Mutant kommt im Build-Kontext nicht an*: entfallen für diesen Slice — frische Kopie mit `tar -x -m` und `touch` je Probe, jede I2-Probe zeigte den Mutanten (rot). *Erste Phase verdeckt die zweite*: entfallen — gezielt mit `-test.run`. *Cleanup nach fehlgeschlagenem Start oder nach `Fatalf`*: entfallen — Cleanup-Läufe aus F-481 und mein Gegenlauf. *Neuer Fund*: entfallen — `defer pc.Close` hing in keinem I2-Lauf; die Hänger vorher stehen im Cleanup, nicht in `Close` (Stack gelesen). Die Wahl trifft der Planner. | — |

Keine Abweichung von §6 *Randformen*; kein Diff entscheidet eine Randform, die §6 nicht nennt.

## 6. Vorbereitung für `slice-v1-abschluss-herunterfahren`

Dessen §1 und DoD erwarten „Tests mit Signal und Frist je Modus über die Helfer“, ein zweites Signal, `--shutdown-timeout 0` und Exit-Code 4. Was der Helfer trägt:

- **Signal:** `signal(t, sig)` nimmt jedes `os.Signal` (`SIGINT` als `os.Interrupt`, `SIGTERM`); erstes und zweites Signal laufen über ihn. Ein zweites Signal an einen schon beendeten Prozess meldet den Befund (F-480); ein Test, der ein Ende vor dem zweiten Signal **erwartet**, liest `beendet` (§6).
- **Frist:** `warteEnde(t, frist, meldung)` mit Frist als Literal am Aufruf, bis 60 s, und eigener Meldung.
- **Exit-Code:** `pruefeExit(t, 4)` mit `ProcessState.String()` und `stderr` bei Abweichung.
- **Prozess endet noch nicht** (`--shutdown-timeout 0`): `select` auf `rec.beendet` mit eigener Frist; der Cleanup beendet den Prozess danach mit Frist.
- **Nicht getragen:** ein wettlauffreies Lesen von `stderr`, **während** der Prozess läuft. Wer vor dem zweiten Signal auf die Info-Zeile mit der Zahl der Sessions oder auf die Log-Zeile einer abgebrochenen Session wartet, müsste `stderr` pollen — genau die Leser, die §6 hier nur als Bestand akzeptiert (V-105).

Der Helfer trägt also, was der Folge-Slice braucht, bis auf das Lesen von Log-Zeilen zur Laufzeit.

## 7. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-102 | LOW (DoD Punkt 2; `AGENTS.md` §3.10; Rest von F-481) | Die Einordnung unter *Grüne Mutanten, eingeordnet* misst mit zwei Maßen. *Cleanup wartet ohne Frist* heißt „grün im Gate, gefangen nur zusammen mit `Kill` entfernt“ (Lauf 1 gegen Lauf 3). *`t.Errorf` → `t.Logf`* heißt „gefangen, kein grüner Mutant mehr“ (Lauf 1 gegen Lauf 2) — dabei haben Lauf 1 und Lauf 2 **beide** `Kill` entfernt; der Mutant ist also genauso nur unter Provokation und zusammen mit `Kill` entfernt gefangen. Ebenso *kein Warten auf `beendet` in `signal`*: gefangen nur unter der Provokation (Kill vor `stop`) und nur über den Wortlaut der Meldung, der Status ist mit und ohne Mutant FAIL. Im Gate ohne Provokation sind alle Mutanten des Geschirrs grün; §1 schließt Mutationen als Gate aus, das ist in Ordnung. Nur die Wörter müssen gleich messen: je Mutant „grün im Gate; unterschieden unter Provokation <P> (mit <zweiter Mutation>) über <Status \| Meldung>“. Adresse Implementer, nur §7. | `docs/plan/planning/in-progress/slice-harness-integration-wait.md` · „gefangen (F-481 Lauf 1 FAIL gegen Lauf 2 PASS), kein grüner Mutant mehr“; · „grün im Gate, gefangen nur zusammen mit `Kill` entfernt“ | Lesen; §7 Zeilen F-481 Lauf 1 bis 3 |
| V-103 | LOW (`SPEC-038` *Warten in Tests*; `AGENTS.md` §3.11) | Der Ablauf-Zweig von `signal` (Signal liefert einen Fehler, `beendet` schließt nicht binnen 5 s) meldet nur `Signal %v: %v`. `SPEC-038` verlangt, dass die Meldung bei Ablauf einer Frist das ausgebliebene Ereignis nennt — hier das Ende des Prozesses. Der Code folgt §6, der genau diesen Text entschieden hat; die Abweichung liegt also zwischen Plan und Spezifikation. Dazu sagt der Kommentar an `signal` diesen Zweig zu, ohne dass ein Lauf ihn zeigt (im Geschirr liefert `Signal` an einen lebenden Prozess keinen Fehler). Praktisch kaum erreichbar, darum LOW. Adresse Architect: Meldung um „Prozess endet nicht binnen 5 s“ ergänzen oder in §6 begründen, warum diese Wartezeit nicht unter `SPEC-038` fällt; den nicht herstellbaren Zweig in §7 als akzeptiertes Negativ nennen. | `test/integration/record_e2e_test.go` · `t.Fatalf("Signal %v: %v", sig, err)`; Plan §6 · „sonst mit `t.Fatalf` und dem Fehler von `Signal`“ | Lesen |
| V-104 | Hinweis (DoD Punkt 2, Beleg) | Die zehn I2-Läufe in §7 liefen am Arbeitsbaum des Umbaus; die Nacharbeit hat den gemessenen Pfad geändert (`nachKill` aus `warteEnde` gezogen, `signal` neu), und am Endstand belegt §7 nur einen Lauf. Die DoD nennt den Stand nicht, die Lücke ist also keine Verletzung. Geschlossen durch meine Messung am Endstand `9fec9e6`: 10 von 10 rot, 0 hängend (Abschnitt 1). Für die Closure-Notiz: Stand je Messung nennen. | `docs/plan/planning/in-progress/slice-harness-integration-wait.md` · „rot statt hängend, Gegenlauf nach der Nacharbeit“ | eigene 10 Läufe |
| V-105 | Hinweis (Architect von `slice-v1-abschluss-herunterfahren`, `AGENTS.md` §3.12) | Der Helfer liest `stderr` erst nach dem Ende wettlauffrei. Tests des Folge-Slice, die vor dem zweiten Signal auf eine Log-Zeile warten (Info-Zeile mit der Zahl der Sessions, Zeile je abgebrochener Session), bräuchten ein Lesen zur Laufzeit; §6 dieses Slice akzeptiert solche Leser nur als Bestand. Eine Randform für §6 des Folge-Slice vor dem Code: Log-Zeilen erst nach dem Ende prüfen, oder ein synchronisierter Puffer im Geschirr. | `docs/plan/planning/next/slice-v1-abschluss-herunterfahren.md` · „ein zweites Signal lässt die Frist sofort ablaufen; beim Beginn steht eine Info-Zeile mit der Zahl der Sessions“ | Lesen |

Kein Befund gegen DoD-Punkt 1, 3 und 4. Kein Produkt-Code im Slice.

---

## 8. Urteil

**Fertig für die Closure; kein Befund blockiert.**

- **Ein `Wait` je Prozess (Punkt 1): bestätigt.** Vier Stellen vorher, eine nachher (`record_e2e_test.go:381`); alle anderen lesen den geschlossenen Kanal, Ergebnis und Zustand sind vor dem Schließen geschrieben.
- **Rot statt hängend (Punkt 2): bestätigt und selbst gemessen.** Vorher 2 von 10 hängend an `awaitGoroutines` im Cleanup, am Endstand 10 von 10 rot nach 6,1 s, keiner hängend. F-480-Provokation und -Mutant nachgefahren, Ergebnis wie in §7; Gegenlauf zur F-481-Provokation ergänzt (PASS, 0,10 s). Die Einordnung der grünen Mutanten misst mit zwei Maßen (V-102).
- **Testliste und Deklarationen (Punkt 3): bestätigt.** 40 Tests, `diff` leer, keine Zeile `Abdeckung:` geändert.
- **`make gates` (Punkt 4): grün**, `make lint` Exit 0.
- **F-480 bis F-482 umgesetzt**: nur `signal` ruft `Process.Signal`; `nachKill` liest `stderr` erst nach dem Ende.
- **`SPEC-038`** im neuen Code eingehalten bis auf die Meldung im kaum erreichbaren Ablauf-Zweig von `signal` (V-103).
- **Folge-Slice:** Der Helfer trägt Signal, zweites Signal, Frist, Exit-Code und „endet noch nicht“; das Lesen von Log-Zeilen zur Laufzeit trägt er nicht (V-105).

**Vor der Closure offen:** V-102 (Wortlaut in §7, Implementer), V-103 (Meldung oder Begründung, Architect). V-104 für die Closure-Notiz, V-105 an den Architect des Folge-Slice. Dazu die Closure-Pflichten mit den Vorschlägen aus Abschnitt 5.
