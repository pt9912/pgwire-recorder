# Verifikation: slice-lint-bestand-driving — 2026-10-07

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft (F-465 bis F-470). F-470 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-lint-bestand-driving.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `b7d4555..dd5d13d`. Darin liegen `78f28a1` (Kontexte und Receiver, Befunde 9 bis 16), `e76b25f` (Komplexität, Befunde 1 bis 8), `76ce087` (Belege in §7), `934add1` ([Review](2026-10-07-review-slice-lint-bestand-driving.md)) und `dd5d13d` (Nacharbeit F-465 bis F-467). `dd5d13d` ändert im Go-Code nur die Doc-Kommentare von `replayLesefehler` und `replayZustellen`.

**Eingang:**

- Plan §1, §2 (DoD), §3, §4, §6 *Randformen* und *Risiken*, §7 *Belege des Implementers*, §8; die Commit-Messages der sechs Commits
- `spec/spezifikation.md` `SPEC-049` (`git diff b7d4555 dd5d13d -- spec/ docs/plan/adr` ist leer); [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 5; [ADR-0001](../plan/adr/0001-hexagonale-architektur.md); [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md); [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- `AGENTS.md` §3.6, §3.9 bis §3.13
- Produkt-Code an `b7d4555` und `dd5d13d`: `internal/adapters/driving/pgwire/server.go`, `internal/bootstrap/bootstrap.go`, `internal/adapters/driving/cli/cli.go`; dazu `server_test.go` (`TestReplaySent`, `replayVerbindung`, `fakeReplayer`) und `server_extended_test.go` um Zeile 846
- `docs/plan/planning/open/slice-tests-ueberlebende-mutanten.md` §1 bis §4; §1 von `slice-harness-mutation`, `slice-harness-integration-wait`, `slice-harness-coverage`, `slice-v1-abschluss-protokollrand` und `slice-v1-abschluss-betrieb`, §6 von `slice-harness-mutation`
- Der Review-Report ganz, den Bericht des Implementers nicht

Bei meinem Start war der Arbeitsbaum sauber, HEAD `dd5d13d`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

**Wie die Proben gebaut wurden.** `b7d4555` und `dd5d13d` habe ich per `git archive` in je eine Kopie in einem eigenen Unterverzeichnis des Scratchpads entpackt, ohne `cp -p`. Jede Mutation lief in einer frischen Kopie aus `git archive dd5d13d`. Das Ersetzungsskript bricht ab, wenn der Suchtext nicht genau so oft trifft wie angegeben, danach `touch` auf die Datei. Unit-Läufe liefen im Image der Stufe `deps` unter eigenem Tag, per Bind-Mount (nur lesend), mit `--network=none`, `go vet ./...` und `go test -count=1`. Lint lief je Kopie mit `docker build --progress=plain --target lint` aus dem eigenen Kontext und eigenem Tag. Der eine Integrationslauf hatte ein eigenes Image, ein eigenes internes Netz, einen eigenen PostgreSQL-Container (gepinntes Image aus `harness/mk/integration.mk`) und ein eigenes Volume, beide Phasen wie `make test-integration`. Danach habe ich Container, Netz, Volume, Images, Cache-Volume und Kopien entfernt. Hostweit habe ich nichts aufgeräumt, und kein Lauf hat in den Arbeitsbaum geschrieben. `make lint` und `make gates` liefen im Arbeitsbaum.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Komplexität unter den Schwellen, ohne Verhaltensänderung. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| die vier Funktionen unter den Schwellen aus `SPEC-049` Punkt 4 | `make lint` an `dd5d13d`: `0 issues.` In einer Kopie von `dd5d13d` mit den drei Schwellen auf 1: `replaySitzung` 17/10/10, `replayWaechter` 2/3/3, `replayLesefehler` –/2/2, `replayAntwort` 7/7/8, `replayZustellen` –/2/2, `clientRichtung` 13/10/11, `nachricht` 5/6/7, `startup` 12/9/10, `startkopf` 3/8/9, `toMessage` 3/10/11, `toExtendedMessage` –/7/8, `rowDescription` –/2/2, `dataRow` 3/3/3 (`gocognit`/`gocyclo`/`cyclop`). Jeder Wert stimmt mit der Tabelle in §7 überein. Am nächsten an einer Schwelle liegt `replaySitzung` mit `gocognit` 17 von 20. | bestätigt |
| Testliste vorher und nachher gleich | `go test -list .` je Paket an `b7d4555` und `dd5d13d`: `postgres` 14, `recording` 14, `cli` 19, `pgwire` 40, `bootstrap` 12, `model` 10, `services` 53, zusammen 162. Die sortierten Namen sind an beiden Ständen byte-gleich. `test/integration`: `-list` scheitert ohne laufendes PostgreSQL am `TestMain`; die 39 Funktionen `func Test…(t *testing.T)` sind an beiden Ständen gleich (sortiert, gleicher SHA-256). | bestätigt |
| keine Erwartung geändert | `git diff b7d4555 dd5d13d --stat -- '*_test.go' test/ docs/user .golangci.yml Dockerfile .dockerignore go.mod go.sum` ist leer. Geändert sind nur `server.go`, `bootstrap.go`, `cli.go`, der Plan und der Review-Report. `export_test.go` ist unberührt. | bestätigt |
| `make test` und `make test-integration` grün | im Gate-Lauf an `dd5d13d` (Abschnitt 5); dazu mein Unit-Lauf X0 an `dd5d13d` mit `-count=1`: alle sieben Pakete `ok` | bestätigt |

### Punkt 2: Kontexte ohne Änderung des Herunterfahrens. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| die sechs Befunde (`containedctx`, 3× `contextcheck` in `pgwire`, `contextcheck` in `bootstrap.go`, `noctx`) behoben | an `b7d4555` stehen sie in der Lint-Ausgabe, an `dd5d13d` nicht (Punkt 3) | bestätigt |
| Form nach §6 *Kontexte* | Gelesen: Feld `ctx` entfällt; `serverRichtung(context.WithoutCancel(ctx))`; `clientRichtung` bildet `uc := context.WithoutCancel(ctx)` und fragt `Err()` allein auf `ctx`; `beende`, `schreibe`, `fehler`, `wartenAufEnde` und `nachricht` bekommen den gelösten Kontext. `bootstrap.go`: `Finish(context.WithoutCancel(ctx))`, `Listen(ctx, …)` in `record` und `replay`. `Listen`: `(&net.ListenConfig{}).Listen(context.WithoutCancel(ctx), "tcp", address)`, Code `PGR-E4001` bleibt. Das bestätigt F-470. | bestätigt |
| Herunterfahren hängt weiter an `ctx.Err()` in `clientRichtung` | V1: `ctx.Err()` → `uc.Err()` an `dd5d13d` ist **rot**, `TestExtendedHerunterfahren`, `TestHerunterfahrenWeckenNichtVerloren` (wie K1 in §7) | bestätigt |
| die Mutation „`WithoutCancel` entfernt“ an 10 bis 13 ist äquivalent (akzeptiertes Negativ, §6) | V5: `service.Finish(ctx)` statt `Finish(context.WithoutCancel(ctx))`: Unit `bootstrap` grün, **und** beide Integrationsphasen grün (37 PASS in Phase 1, darunter `TestE2ERecordExtendedSigtermNachBlockade`, `TestE2ERecordExtendedSigtermBeimPipelining`, `TestE2EReplayExtendedSigtermMittenInFolge`, `TestE2EReplayNichtVerbrauchtHerunterfahren`). Mit beendetem Kontext wird die Aufzeichnung also weiter geschrieben, weil `YAML.Write` ihn nicht liest. V6: `serverRichtung(ctx)` ungelöst: Unit grün. Beide bestätigen die Äquivalenz, die §6 behauptet, statt einer Lücke. | bestätigt |
| Tests und Integrationstests zum Abbruchsignal unverändert grün | Testdateien unverändert (Punkt 1), Gate-Lauf grün | bestätigt |

### Punkt 3: `make lint` ohne Befund, modulweit. **Bestätigt — Vorbedingung für `slice-harness-lint` erfüllt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| **modulweit 0 Befunde** | `make lint` im Arbeitsbaum an `dd5d13d`: Exit 0, Ausgabe der Stufe `0 issues.`, keine `lint:`-Zeile der eigenen Prüfungen. Die Stufe lief neu (kein `CACHED`, 10,1 s). | bestätigt |
| vorher 16, alle unter den drei Pfaden | Kopie `b7d4555`: Exit ungleich 0, `16 issues:`, alle 16 unter `internal/adapters/driving/pgwire`, `internal/adapters/driving/cli` und `internal/bootstrap`. Sortiert byte-gleich mit dem Block *Vorher* in §7. | bestätigt |
| ohne `//nolint` | Das Muster aus `tools/harness/lint.sh` über `cmd`, `internal` und `test` an `dd5d13d`: 0 Treffer. | bestätigt |
| ohne neue Ausnahme in `.golangci.yml` | `git diff b7d4555 dd5d13d -- .golangci.yml` leer, ebenso `tools/harness/lint.sh` und `Dockerfile`. | bestätigt |
| Zeilen vorher und nachher im Bericht | §7 nennt die 16 Zeilen vorher, die acht nach `78f28a1` als Beschreibung und nachher keine. | bestätigt |

### Punkt 4: `make gates` grün. **Bestätigt** (Abschnitt 5).

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| Review mit Report | liegt vor (F-465 bis F-470) |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen: §7 trägt Platzhalter, alle vier Risiken in §6 tragen „— (bei Closure)“. Das ist Sache der Closure; Vorschläge in Abschnitt 6. Der Ausgang von *Verhalten von `replaySitzung` ändert sich unbemerkt* braucht eine Folge-Slice-Kennung, die es noch nicht gibt (Abschnitt 4). |

### Weitere Schwerpunkte des Auftrags

- **Abdeckungs-Deklarationen gleich.** An `dd5d13d` stehen 137 Zeilen `Abdeckung:` unter `internal`, `test` und `cmd`, so viele wie am Stand von `slice-lint-bestand-kern-driven`. Da keine Testdatei und nichts unter `docs/user/` geändert ist, sind sie gleich. `make abdeckung-check` ist im Gate-Lauf grün.

---

## 2. Ohne Verhaltensänderung: Stichproben gegen `b7d4555`

Grundlauf X0 an `dd5d13d` ohne Mutation: alle Unit-Pakete grün.

| ID | Zusage | Mutation an `dd5d13d` | Ergebnis |
|---|---|---|---|
| V1 | Herunterfahren an `ctx.Err()` in `clientRichtung` (K1) | `ctx.Err()` → `uc.Err()` | **rot**, `TestExtendedHerunterfahren`, `TestHerunterfahrenWeckenNichtVerloren` |
| V5 | `Finish` mit gelöstem Kontext ist gleich `context.Background()` | `service.Finish(ctx)` | grün in Unit und beiden Integrationsphasen: äquivalent, wie §6 sagt |
| V6 | gelöster Kontext der Server-Richtung | `serverRichtung(ctx)` | grün: äquivalent, wie §6 sagt |
| V7 | `startkopf`: unbekannte Sonderanfrage ist `PGR-E6001` | `CodeUnsupported` → `CodeProtocolVersion` | **rot**, `TestUnbekannteSonderanfrage` |
| V8 | `startkopf`: andere Protokollversion ist `PGR-E6002` | `CodeProtocolVersion` → `CodeUnsupported` | **rot**, `TestAndereProtokollversion` |
| V9 | `nachricht`: nicht unterstützte Nachricht endet mit `EndUnsupported` | → `EndFailed` | **rot**, `TestExtendedNichtUnterstuetzt`, `TestFehlerantwortMitFrist`, `TestNichtUnterstuetzteNachricht` |
| V10 | Wächter der Replay-Sitzung weckt das Lesen (W) | `SetReadDeadline(time.Now())` entfernt, `-timeout 60s` | rot nur über die Zeitüberschreitung, „running tests: `TestReplayHerunterfahrenSpaeteFrist` (58s)“ (F-468 bestätigt) |
| V11 | nach einem Lesefehler endet die Replay-Sitzung (R4) | `return` → `continue` | grün, wie §7 |
| V12 | Verbindungsende im Replay ist regulär (R2) | `verbindungsende(err) && false` | grün, wie §7 |
| V13 | `replayWaechter`: „Ein Schließen von fertig beendet ihn.“ | `case <-fertig:` entfernt | **grün** in `pgwire` und `bootstrap` (V-93) |

Gelesen gegen `b7d4555`: Die Fall-Reihenfolge in `replayAntwort` und `nachricht` ist die alte, `continue` wurde zu `return true`, jedes `return` zu `return false`. `startkopf` prüft die Länge vor dem Startcode (S1 rot laut §7, Grund wie in §1). `ReadyForQuery` wandert in `toMessage` vor den `default`, der Switch läuft über disjunkte Konstanten. `weiterlesen()` wird weiter gerufen und sein Wert verworfen. Herunterfahren: Wächter in `recordSitzung` unverändert, `replayWaechter` sieht denselben `ctx` wie vorher. F-470 ist damit bestätigt.

---

## 3. F-465 und die Kommentare

**F-465: umgesetzt.** An `dd5d13d` sagen die beiden Kommentare nur zu, was eine rote Mutation fängt:

- `replayLesefehler`: „behandelt einen Lesefehler der Replay-Sitzung.“ Der Satz nennt den Gegenstand und sagt kein Verhalten zu. Die Sätze, die R2, R3 und R4 nicht fangen, sind entfallen.
- `replayZustellen`: „schreibt Antworten an den Client und meldet sie dem Use Case danach als gesendet, nicht, wenn das Schreiben scheitert.“ Jede Hälfte selbst nachgefahren:
  - V4 (`Sent` entfernt, X8): **rot**, `TestReplaySent`.
  - V3 (`Sent` vor `send`, Z2): **rot**, `TestReplaySent` („Sent nach gescheitertem Senden gemeldet (1)“).
  - V2, eine eigene Mutation, die nur die verneinende Hälfte bricht: `Sent` zusätzlich im Fehlerzweig. **Rot**, `TestReplaySent`, `server_test.go:730`, aus dem richtigen Grund. Der zweite Teil von `TestReplaySent` erzwingt über `net.Pipe` einen gescheiterten Versand.
  - „scheitert das Schreiben, merkt es den Fehler“ ist entfallen.

**Ein Rest im selben Diff, außerhalb von F-465 (V-93).** Der neue Doc-Kommentar von `replayWaechter` (aus `e76b25f`) schließt mit „Ein Schließen von fertig beendet ihn.“ V13 entfernt `case <-fertig:`, und alle Tests bleiben grün. Die Goroutine läuft dann weiter, bis `ctx` endet, je Replay-Verbindung eine. Die übrigen neuen Kommentare habe ich stichprobenweise gefahren: `startkopf` (V7, V8, S1), `nachricht` (V9, C1, C2), `dataRow` (T2) und `toExtendedMessage` (`TestToMessageErfindetNichts`). Sie sind gedeckt. „setzt er die Lesefrist von conn auf jetzt“ ist über W nur durch die Zeitüberschreitung gedeckt, aber deterministisch (F-468).

---

## 4. Grüne Mutanten G1 bis G6, Nehmer und F-468

### Einordnung, Test-Idee und Grenze

| Mutant | Einordnung in §7 | Meine Probe | Urteil |
|---|---|---|---|
| G1 (K4, `Listen` ohne `WithoutCancel`) | verhaltensändernd, nur mit Namen als Adresse sichtbar | nicht nachgefahren. Gelesen: Mit IP-Adresse liest `net` den Kontext nicht. Die Grenze ist plausibel und in §6 vorab genannt. | stimmt |
| G2 (R2, R4) | fangbar über `Handle` mit Replay-Fake, Client schließt ohne Terminate | Temporärer Test in einer Kopie, je 5 Läufe: ohne Mutant grün. Unter R2 rot („erster Fehler "PGR-E6001", erwartet leer“). Unter R4 rot („Handle kehrt nicht zurück“ nach 2 s, `ctx` = `Background`). Die nach F-467 ergänzte Grenze (nur solange `ctx` nicht endet) stimmt mit dem Code überein. | stimmt |
| G3 (R3) | fangbar über einen Nachrichtentyp, den `pgproto3` nicht kennt | Temporärer Test: Typ `'!'`, Länge 4. Ohne Mutant ErrorResponse `0A000` und `PGR-E6001`; unter R3 rot („unexpected EOF“). | stimmt |
| G4 (T3, T5) | Durchreichen ohne Logik, fangbar über `ToMessage` | nicht nachgefahren; gelesen: `rowDescription` setzt beide Felder, kein Test vergleicht sie | stimmt |
| G5 (X4) | fangbar; Test-Idee mit Wrapper um `net.Conn`, Grenze „ein geschlossener TCP-Peer lässt den ersten Schreibvorgang oft gelingen“ | Temporärer Test, der den zweiten Teil von `TestReplaySent` (`net.Pipe`, Client schließt nach der Anfrage) um `s.FirstErrorCode() == PGR-E4003` ergänzt, 50 Durchläufe: ohne Mutant grün. Unter X4 rot im ersten Durchlauf („erster Fehler "", erwartet PGR-E4003“). | Einordnung stimmt. Test-Idee und Grenze unnötig schwer (V-94) |
| G6 (Z3) | nur über eine nicht abbildbare Antwort fangbar | gelesen: Bei `net.Pipe` scheitert nach dem Versandfehler auch das nächste Lesen, Z3 ist dort nicht sichtbar. Die Grenze stimmt. | stimmt |

Dass G2 bis G6 schon vor dem Umbau grün waren, belegt §7 an `78f28a1`. Der Code von `replaySitzung` ist dort derselbe wie an `b7d4555`, und das Review hat X4 an `b7d4555` grün gefahren. Ich habe es nicht wiederholt.

### Nehmer der grünen Mutanten — Empfehlung an den Planner

**Kein bestehender Slice nimmt an** (`AGENTS.md` §3.13, geprüft an §1 *Ausdrücklich NICHT* und DoD):

- `slice-tests-ueberlebende-mutanten`: §1 schließt „Der PGWire-Adapter (`internal/adapters/driving/pgwire`) und `test/integration`“ als Schicht-Abgrenzung aus. G1 liegt im Bootstrap, auch das ist eine dritte Schicht neben Kern und Driven. Seine *Sammelregel* sagt für diesen Fall selbst: Der Planner schneidet vor dem Start einen zweiten Slice ab. Die Zeile „seine Funde nimmt dieser Slice nach dieser Regel an“ gilt also nur über diesen Schnitt (F-469 b).
- `slice-harness-mutation`: §1 schließt „Neue Tests für Mutanten, die heute überleben“ aus.
- `slice-harness-coverage`: §1 schließt neue Tests über dem gemessenen Stand aus (F-459).
- `slice-harness-integration-wait`: nur `test/integration`, keine neuen Tests oder Fristen.
- `slice-v1-abschluss-protokollrand`: Gegenstand ist `LH-FA-05.e` (Startphase, Sonderanfragen), nicht die Replay-Sitzung.
- `slice-v1-abschluss-betrieb`: baut `replaySitzung` für die Frist um, hat aber schon mehr als drei Liefer-Punkte im Gegenstand. Er ist eine Kopplung, kein Nehmer (V-95).

**Damit ist ein zweiter Sammel-Slice zu schneiden. Das stimmt.** Vorschlag: `slice-tests-ueberlebende-mutanten-driving`, Schichten Driving Adapter und Bootstrap (zwei, wie dieser Slice), vor `slice-harness-mutation` und mit derselben Sammelregel. Nach Lieferwert gruppiert, drei Liefer-Punkte:

1. Lesefehler der Replay-Sitzung: G2 und G3 (Pakete `pgwire_test`, Fake und `net.Pipe` sind vorhanden).
2. Versand der Replay-Sitzung: G5 und G6 (G5 als Erweiterung von `TestReplaySent`, V-94).
3. Abbildung und Öffnen am Rand: G4 (`ToMessage`) und G1 (`bootstrap.Run` mit `localhost:0`).

Punkt 3 bündelt zwei Pakete mit dünnem gemeinsamen Lieferwert. Fasst der Planner G1 nicht so, bleibt nur ein vierter Punkt oder ein dritter Slice. G4 beantwortet zugleich die offene Grenze von P8/P2 im ersten Sammel-Slice (F-469 c); dort sollte der Rückverweis stehen. Nach §3.13 trägt der Planner im selben Commit die Sendung in §1 und DoD des neuen Slice mit der Kennung `slice-lint-bestand-driving` ein, und dieser Plan nennt dann die Kennung im Risiko-Ausgang.

### F-468 — Mutant W nur über die Paket-Zeitüberschreitung rot

V10 bestätigt den Befund: `TestReplayHerunterfahrenSpaeteFrist` blockiert bei `<-conn.gesetzt` (`server_extended_test.go:846`) ohne eigene Frist. Empfohlene Adresse, zweiteilig:

- **Wie ein Gate ihn zählt:** `slice-harness-mutation` §6, Randform *Laufzeit* („Zeitgrenze je Mutant (ein Mutant, der eine Schleife endlos macht)“). Sie ist offen und nennt genau diesen Fall. Die Adresse nimmt an: Weder §1 noch DoD schließen ihn aus. Einzutragen ist W dort als Beispiel mit der Kennung dieses Slice.
- **Die Frist im Test selbst** (ein `select` mit Frist statt `<-conn.gesetzt`) ist eine Teständerung im PGWire-Adapter. `slice-harness-integration-wait` verfolgt dasselbe Ziel („rot statt hängend“), schließt aber alles außer `test/integration` und neue Fristen aus. Passend ist der neue Sammel-Slice oben, als Teil von Liefer-Punkt 1 oder als Randform, wenn der Architect bei *Laufzeit* entscheidet, dass Tests eine eigene Frist tragen. Die Entscheidung liegt beim Architect, weil sie für alle Tests mit Kanal-Warten gilt, nicht nur für diesen.

---

## 5. Gate-Ergebnis

- `make gates` an `dd5d13d` im Arbeitsbaum: **grün**, Exit 0. Im Log: `run-integration-tests: gruen` (47 `--- PASS`, 2 `--- SKIP`, kein `--- FAIL`), `kopf-check-gegenprobe`, `abdeckung-gegenprobe`, `a-check-negativ` und `commit-msg-gegenprobe` grün, `d-check: 327 Datei(en) geprüft, 0 Befund(e)`. Die Stufe `test` kam aus dem Build-Cache. Docker speichert nur erfolgreiche Schritte, ein Treffer heißt also, dass derselbe Inhalt schon grün war. Mein Lauf X0 mit `-count=1` hat die Unit-Tests an `dd5d13d` zusätzlich frisch gefahren.
- `make lint` an `dd5d13d`: Exit 0, `0 issues.` (Abschnitt 1, Punkt 3).
- Commit-Messages: Alle sechs nennen `slice-lint-bestand-driving` und [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), keine eine `SPEC-` oder `ARC-`Kennung.

---

## 6. Plan gegen Code

- **Kopf:** `Bezug` ([`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md)) und `Berührte Spec-Stellen` (`SPEC-049`, unverändert) stimmen. `kopf-check` ist im Gate-Lauf grün.
- **§1:** Das Ziel ist erreicht. Die 16 Befunde sind die aus der Tabelle (byte-gleich), und modulweit bleibt keiner. Die Abgrenzung hält: `internal/hexagon/` ist unberührt, keine Testdatei, kein `//nolint`, keine Ausnahme, keine Spec- oder Harness-Änderung. Die geänderte exportierte Signatur `pgwire.Listen(ctx, address)` rufen nur `record` und `replay` in `internal/bootstrap`. Keine Testdatei ruft sie.
- **§3:** Jede Zeile hat ihre Änderung im Diff, und umgekehrt gibt es keine Änderung ohne Zeile. Die Reihenfolge der Commits (Kontexte und Receiver, dann Komplexität) entspricht dem Absatz *Größe*.
- **§4:** Keine Rückführung ausgelöst. Der Umbau von `replaySitzung` war in einer Review-Sitzung prüfbar, `contextcheck` nimmt `WithoutCancel` an, und kein Befund verlangte eine Verhaltensänderung.
- **§5:** `make gates` grün und `make lint` modulweit ohne Befund sind gemessen.
- **§6 *Randformen*:** Alle entschiedenen hält der Code ein (Kontexte, `Listen`, Receiver, `weiterlesen`). Der Diff entscheidet keine neue. Die beiden Äquivalenz-Behauptungen habe ich mit V5 (auch Integration) und V6 bestätigt.
- **§6 *Risiken*, Vorschläge für die Ausgänge** (setzen muss sie die Closure):
  - *`contextcheck` nimmt `context.WithoutCancel` nicht an:* **entfallen**. Begründung: Nach `78f28a1` meldet Lint keinen `contextcheck`, an `dd5d13d` modulweit `0 issues.`
  - *Herunterfahren ändert sich unbemerkt:* **entfallen**. Begründung: K1/V1 rot vor und nach dem Umbau, Signaltests der Integration grün, `WithoutCancel` an 10 bis 13 äquivalent (V5 einschließlich Integration, V6).
  - *Verhalten von `replaySitzung` ändert sich unbemerkt:* **eingetreten** → Folge-Slice mit Kennung, sobald der Planner ihn angelegt hat (Abschnitt 4). Vorher kann der Slice nicht nach `done/`.
  - *Aufzählung weicht vom Lauf ab:* **entfallen**. Begründung: Der Lauf an `b7d4555` ist byte-gleich mit §1 und §7.
- **§7 *Belege*:** Vollständig für das, was die DoD verlangt: Lint-Zeilen vorher und nachher, Komplexitätstabelle, Testliste, Diff ohne Testdatei, Läufe, Mutationstabelle, Kommentar-Mutationen nach F-465, grüne Mutanten. Jede Zahl, die ich nachgemessen habe, stimmt. `make gates` an `dd5d13d` steht nicht in §7 (der Kopf nennt `76ce087` als letzten Gate-Stand). Den Lauf habe ich selbst gefahren (Hinweis, kein Befund).
- **§8:** Sichtung am Stand `0d464b4`. Für die Closure: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` ist in diesem Slice mit F-465 und V-93 getroffen.

---

## 7. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-93 | LOW (Kommentar-Zusage ohne Test, gleiche Klasse wie F-465) | Der neue Doc-Kommentar von `replayWaechter` sagt „Ein Schließen von fertig beendet ihn.“ Ohne `case <-fertig:` (V13) bleiben alle Tests in `pgwire` und `bootstrap` grün. Der Wächter läuft dann je Replay-Verbindung weiter, bis `ctx` endet. Nach außen wirkt das nicht, man sieht es nur an der Zahl der Goroutinen. `AGENTS.md` §3.11 verlangt eine engere Fassung oder einen Test. Vorschlag: den Satz streichen, denn die übrigen zwei Sätze tragen (W). Adresse: Implementer, nur ein Kommentar, vor der Closure. | `internal/adapters/driving/pgwire/server.go` · „Ein Schließen von fertig beendet ihn.“ | V13 grün |
| V-94 | Hinweis (Planner, Test-Idee G5) | Die Test-Idee von G5 verlangt einen Wrapper um `net.Conn` und begründet das mit einem TCP-Peer, der den ersten Schreibvorgang gelingen lässt. Das Testgeschirr des Pakets nutzt aber `net.Pipe`, und der zweite Teil von `TestReplaySent` erzwingt den gescheiterten Versand schon deterministisch. Mit `s.FirstErrorCode() == PGR-E4003` an dieser Stelle wird X4 rot, und ohne Mutant ist der Test in 50 von 50 Durchläufen grün. Der Nehmer sollte die einfachere Test-Idee übernehmen und die Grenze auf `net.Pipe` beziehen. Den Pfad über die nicht abbildbare Antwort deckt G6. | `docs/plan/planning/in-progress/slice-lint-bestand-driving.md` · „Wrapper um `net.Conn`“ | temporärer Test in einer Kopie, mit und ohne X4 |
| V-95 | Hinweis (Planner, Kopplung) | `slice-v1-abschluss-betrieb` ändert, wie `replaySitzung` ein Verbindungsende behandelt: Schließt die Frist die Verbindung, merkt er für eine unvollständige Interaktion `PGR-E4006`. Die Test-Idee von G2 hält „`FirstErrorCode()` bleibt leer“ für ein Ende, das der Client auslöst. Die Grenze von G2 im neuen Sammel-Slice sollte nennen, dass ein Ende durch die Frist nicht gemeint ist. Sonst muss `slice-v1-abschluss-betrieb` den Test ändern, statt ihn zu ergänzen. | `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` · „dieses Ende behandelt `replaySitzung` heute als regulär“ | Lesen |

Kein Befund gegen die DoD-Liefer-Punkte.

---

## 8. Urteil

Die drei Liefer-Punkte sind erfüllt, und ich habe sie selbst gemessen:

- **Lint:** Modulweit fallen die Befunde von 16 auf **0**. `make lint` an `dd5d13d` endet mit Exit 0 und `0 issues.`, ohne `lint:`-Zeile, ohne `//nolint`, und `.golangci.yml` ist unverändert. **Die Vorbedingung für `slice-harness-lint` (Bestand ohne Befund) ist erfüllt.**
- **Komplexität:** Alle vier Funktionen und ihre neun Hilfsfunktionen liegen unter den Schwellen. Die Werte stimmen mit §7 überein.
- **Kontexte:** Die sechs Befunde sind behoben, in der Form, die §6 entscheidet. Das Herunterfahren hängt weiter an `ctx.Err()` (V1 rot). Die als äquivalent erklärten Mutationen sind äquivalent, `Finish` auch im Integrationslauf.
- **Ohne Verhaltensänderung:** Die Testliste ist gleich (162 Unit-Tests, 39 Integrationstests). Keine Testdatei und keine Abdeckungs-Deklaration ist geändert. `make gates` ist grün, und die Stichproben V7 bis V9 sind rot.
- **F-465** ist umgesetzt. Beide Kommentare sind Satz für Satz von roten Mutationen gedeckt, auch die verneinende Hälfte (V2). Ein gleichartiger Rest steht in `replayWaechter` (V-93).
- **G1 bis G6** sind richtig eingeordnet: Jeder ändert das Verhalten und ist fangbar, G2, G3 und G5 habe ich mit temporären Tests gezeigt. Kein bestehender Slice nimmt sie an. Der zweite Sammel-Slice nach der Sammelregel ist der richtige Weg (Abschnitt 4).
- **F-468** gehört für die Zählung an die Randform *Laufzeit* von `slice-harness-mutation` und für die Frist im Test an den neuen Sammel-Slice, nach Entscheidung des Architect.

Kein Befund blockiert die DoD. Vor der Closure sind offen: V-93 (Implementer, ein Satz), der neue Sammel-Slice des Planners mit Kennung für den Risiko-Ausgang, V-94 und V-95 als Eintrag in diesem Slice, und die Closure-Pflichten (Abschnitt 6).
