# Verifikation: slice-replay-semantik-fehlerreplay — 2026-10-06

**Rolle:** Verifier (Modul 11). Ich prüfe, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-replay-semantik-fehlerreplay.md` (Kopf, §1 bis §6, §8) gegen den Gesamt-Diff `13bb028..ce50a10` ohne die in `5ee8116` und `20f9044` angelegten fremden Slice-Pläne. Spezifikation: `a8c490e`, `186f476`, `8b399ba`, `9e93b2f`. Code: `212a3d4`, `af07614`, `ce50a10`. Review-Report: `016b993`.

**Eingang:**

- die DoD-Liefer-Punkte und die Commit-Messages
- der [Review-Report](2026-10-05-review-slice-replay-semantik-fehlerreplay.md) (F-390 bis F-398, Stand `af07614`). Die Nacharbeit `ce50a10` hat kein eigenes Review gesehen. Ich habe sie mit eigenen Mutationen geprüft (Abschnitt 2).
- [`LH-FA-02`](../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-05`](../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-06`](../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [`LH-FA-10`](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-11`](../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers), [`LH-FA-12`](../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen) mit den Abnahmeszenarien 4 und 6. In der Spezifikation: [`LH-FA-02.b`](../../spec/spezifikation.md#lh-fa-02b--abschluss-einer-interaktion), [`LH-FA-05.a`](../../spec/spezifikation.md#lh-fa-05a--unterstützter-umfang-startup-simple-query-und-extended-query), [`LH-FA-11.a`](../../spec/spezifikation.md#lh-fa-11a--aufzeichnung-von-fehlerantworten)
- der Welle-Plan `docs/plan/planning/welle-replay-semantik.md` §3 und die Folge-Slices `docs/plan/planning/open/slice-v1-abschluss-einspielen.md` und `docs/plan/planning/open/slice-v1-abschluss-antwortvergleich.md`

Bei meinem Start war der Arbeitsbaum sauber. Während meines Laufs kamen fremde Änderungen an acht Planungsdateien hinzu (Roadmap, sieben Slice-Pläne in `open/`), die nicht von mir stammen. Ich habe sie nicht angefasst, und dieser Commit enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

Die Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg.

**Wie die Mutanten gebaut wurden.** Alles lief unter `verify-fr/` im Scratchpad. Für jeden Mutanten habe ich eine frische Kopie aus `git archive ce50a10` angelegt. Die Ersetzung erfolgte über ein Skript, das abbricht, wenn die Stelle nicht genau einmal vorkommt. Den Docker-Build-Kontext habe ich umgangen (implement-slice Schritt 19, Bind-Mount). Ein Image der Stufe `deps` des `Dockerfile` (`verify-fr-deps`) trägt nur die Module und keinen Quellstand.

- **Unit-Tests:** Die Kopie ist schreibgeschützt eingebunden. Darin läuft netzlos `go test -count=1`.
- **E2E:** Binary und Integrationstests entstehen im Container aus der eingebundenen Kopie. Sie laufen gegen `postgres:17-alpine` mit dem Digest aus `harness/mk/integration.mk`, in einem eigenen internen Docker-Netz, mit `-v`.

Ohne Mutation ist die Kopie auf beiden Wegen grün. Die drei E2E-Tests dieses Slice liefen dreimal hintereinander grün (`-count=3`).

Eigene Container, Netze, das Image `verify-fr-deps` und die Kopien habe ich danach gelöscht. Die Images, die `make gates` unter den Namen des Repos baut (`pgwire-recorder:dev`, `:test`, `:integration`), sind nicht meine und bleiben stehen.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — [LH-FA-11](../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers): Abnahmeszenario 6 für Simple Query. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| eine einzelne Fehlerantwort (`SELECT 1/0`) | `TestE2EFehlerreplayEinfach`, `enthaeltInFolge` verlangt `C=22012`, `M=division by zero`, `ReadyForQuery I` | bestätigt |
| Fehler in der zweiten Anweisung nach dem Ergebnis der ersten | `SELECT 2 AS b; SELECT 1/0`: RowDescription, DataRow `"2"`, CommandComplete, ErrorResponse in Reihenfolge; VF15 rot | bestätigt |
| Fehler in einer Transaktion mit `ReadyForQuery` im Status `E`, danach `ROLLBACK` | `BEGIN` (T), `SELECT 1/0` (E), `SELECT 3` (`25P02`, E), `ROLLBACK` (I); VF12 und VF17 rot | bestätigt |
| im Replay dieselbe Sicht wie beim Aufzeichnen und wie direkt gegen PostgreSQL | `dreiSichten`: record gleich direkt und replay gleich direkt, also auch replay gleich record. Eine Mutation nur im Replay wird bei „Sicht im Replay“ rot (VF12, VF14, VF15). Eine Mutation nur in der Weitergabe beim Aufzeichnen wird bei „Sicht beim Aufzeichnen“ rot (VF17). | bestätigt |
| Sicht umfasst Ergebnisse, SQLSTATE, Meldung, Transaktionsstatus | Vergleich der vollständigen Server-Nachrichten nach dem Aufbau, Felder nach Code sortiert; VF14 (SQLSTATE), VF15 (Ergebnisse), VF12 (Status) rot | bestätigt |
| [LH-FA-11](../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers) Boundary: Position in einer Folge mehrerer Interaktionen | Der Ablauf beginnt mit `SELECT 1 AS a` und endet mit `SELECT 4 AS c`, beide fehlerfrei. Die Deklaration sagt das jetzt so (F-394 behoben). | bestätigt |
| Abdeckung `LH-FA-11/Happy`, `LH-FA-11/Boundary` (E2E) | `docs/user/abdeckung-e2e.md` Zeilen 33 und 34; `make abdeckung-check` grün | bestätigt |

### Punkt 2 — [LH-FA-06](../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [LH-FA-12](../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen): Ergebnisarten der Simple Query. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| mehrzeilige Ergebnismenge | `generate_series(1, 3)`: drei DataRow, `SELECT 3` | bestätigt |
| leere Ergebnismenge | `WHERE false`: RowDescription, `SELECT 0` | bestätigt |
| Befehle ohne Zeilen (`CREATE TABLE`, `INSERT`, `UPDATE`) | drei CommandComplete in Reihenfolge | bestätigt |
| mehrere Ergebnismengen einer Anfrage | zwei RowDescription/DataRow/CommandComplete-Folgen, mit NULL und `numeric` | bestätigt |
| zwei `NoticeResponse` in ihrer Reihenfolge | `M=eins` vor `M=zwei`; VF13 (Hinweis weggelassen) und VF19 (beide vertauscht) rot | bestätigt |
| `ParameterStatus` nach `SET` | `ParameterStatus application_name=ergebnisarten` nach `CommandComplete SET`; VF16 rot | bestätigt |
| Binärwert über einen Binär-Cursor ([`SPEC-003`](../../spec/spezifikation.md)) | `format=1`, DataRow `"\x00\x00\x01\x02" "\x00\xff"`. `\x00\xff` ist kein gültiges UTF-8 und geht deshalb über `base64` in die YAML-Aufzeichnung. VF20 (base64-Wert beim Lesen um ein Byte gekürzt) rot | bestätigt |
| im Replay in Reihenfolge und Inhalt wie beim Aufzeichnen | `dreiSichten`, wie Punkt 1; VF13, VF19, VF20 bei „Sicht im Replay“ rot, VF16 bei „Sicht beim Aufzeichnen“ | bestätigt |
| Abdeckung `LH-FA-06/Happy`, `LH-FA-06/Boundary`, `LH-FA-12/Happy` (E2E) | `docs/user/abdeckung-e2e.md` Zeilen 21, 22 und 36 | bestätigt |

### Punkt 3 — Die vier Festlegungen aus §6 (Zeilen 1, 9, 10 und 11 mit 11a bis 11h). **Bestätigt.**

| Festlegung und Teilbehauptung | Beleg | Urteil |
|---|---|---|
| **Weitergabe** ([`LH-FA-02.b`](../../spec/spezifikation.md#lh-fa-02b--abschluss-einer-interaktion), Zeile 9): `SELECT 1; COPY (SELECT 1) TO STDOUT` liefert dem Client kein Ergebnis, nur `PGR-E6001` (E2E) | `TestE2ERecordWeitergabeNichtUnterstuetzt`; VF18 rot („erhalten 1 Ergebnisse“). Session verworfen (`sessions: []`), Exit-Code 6. | bestätigt |
| **Nachrichten zwischen Interaktionen** ([`LH-FA-05.a`](../../spec/spezifikation.md#lh-fa-05a--unterstützter-umfang-startup-simple-query-und-extended-query), Zeile 1): Zuordnung zur nächsten Interaktion | `TestNachrichtenZwischenInteraktionen` (Review M6 rot, Test seitdem unverändert) | bestätigt |
| … nach der letzten Interaktion nicht aufgezeichnet | `TestRecordNachDerLetztenInteraktion`; VF10 rot (`Server-Nachricht "notice_response" ohne laufende Interaktion`) | bestätigt (F-392 behoben) |
| … Hinweise einer Lebendprüfung gibt das Replay nicht wieder | `TestReplayLebendpruefungOhneHinweise`; VF11 rot (Hinweise der Lebendprüfung erscheinen vor `SELECT 2`) | bestätigt (F-392 behoben) |
| **Diagnosefelder** ([`LH-FA-11.a`](../../spec/spezifikation.md#lh-fa-11a--aufzeichnung-von-fehlerantworten), Zeile 10): rohe Bytes mit leerem Feld, `P` = `0`, `p` ohne Zahl, leerem `L` | `TestDiagnosefelder` (Review M5 rot) | bestätigt |
| … unbekannter Feldcode mit leerem Wert steht als `""` in den Feldern | `rohFelder` trägt `Y\x00`, `want` verlangt `"Y": ""`; VF09 rot | bestätigt (Test aus F-397 nachgereicht) |
| **Fehlerantwort vor dem Abbruch** ([`LH-FA-02.b`](../../spec/spezifikation.md#lh-fa-02b--abschluss-einer-interaktion), Zeile 11): `PGR-E6001` mit SQLSTATE statt `PGR-E4003` | `TestFehlerantwortVorDemAbbruch` a, b, d, f; VF01, VF07 rot | bestätigt |
| … Schweregrad zählt nicht (11a) | VF06 (nur `FATAL` gemerkt) rot bei `ERROR` und `PANIC` | bestätigt |
| … Merker mit `ReadyForQuery` zurückgesetzt (11b) | VF02 rot bei Test e und h3 | bestätigt |
| … SQLSTATE und Meldung der **letzten** `ErrorResponse` (11c) | Fall „zwei Fehlerantworten, die letzte zählt“; VF03 rot | bestätigt (F-393 behoben) |
| … Session nicht übernommen (11d) | Test g; VF08 und VF22 rot | bestätigt |
| … auch beim Senden (11g): `Send` stuft nach dem Merker ein | `TestSendNachFehlerantwort` h1 bis h3 mit `CloseWrite`; VF04 rot bei h1, VF02 rot bei h3 | bestätigt (F-390 behoben) |
| … unabhängig von der Reihenfolge der Richtungen (11g) | `TestRecordExtendedSendefehlerNachFehlerantwort` i1 und i2. VF05 (Markierung in `ClientMessage` entfernt) rot bei i1: „Session übernommen“. VF08 rot bei i2. | bestätigt (F-390 behoben), Abweichung der Reihenfolge in i1 siehe V-53 |
| … eine nicht lesbare Serverantwort ist stets `PGR-E6001` (11h) | `TestNichtLesbareServerantwort` j1 bis j3. VF21 rot bei j1 bis j3, VF23 rot bei j2. Die Sonde zeigt den Text „Serverantwort des Upstreams nicht lesbar: unknown message type: Y“ bzw. „… ReadyForQuery body must have length of 1, but it is 0“, ohne `57P01`. | bestätigt (F-391 behoben) |
| Unit-Tests am Upstream-Adapter mit Fake-Server und am Record-Service | wie oben | bestätigt |

### Punkt 4 — `make gates` grün. **Bestätigt.**

`make gates` habe ich selbst gefahren, in einem frischen Klon des Repos, ausgecheckt auf `ce50a10`. **Exit-Code 0.**

- `a-check`: 0 Befunde
- `a-check-negativ`, `abdeckung-gegenprobe`, `commit-msg-gegenprobe`, `kopf-check-gegenprobe`: grün
- `baseline-verify`: v6.13.0 OK, 54 Dateien
- d-check: 224 Dateien, 0 Befunde
- Unit-Tests in der Stufe `test`: alle Pakete `ok`
- `run-integration-tests`: grün, beide Phasen
- `make kopf-check` und `make abdeckung-check` einzeln: Exit-Code 0, auf `ce50a10` und auf HEAD `20f9044`

### Prozess-Punkte der DoD (nur vermerkt, kein Häkchen gesetzt)

Der Review-Report liegt vor. Closure-Notiz, Register, Risiko-Ausgänge und Paarungen stehen aus. §7 ist leer, und vor der Closure ist das richtig. Die drei Risiken in §6 tragen „offen bis Closure“. Der Ausgang wird bei der Closure entschieden.

---

## 2. Bewusstes Brechen

Je Mutant eine frische Kopie. Die Rot-Spalte gibt die Meldung des Tests an, damit der Grund prüfbar ist. 22 von 23 Mutationen sind rot. Grün bleibt VF24, und das betrifft nicht diesen Slice, sondern den Closure-Trigger der Welle (V-52).

| # | Mutation | Rot in | Grund (Testmeldung) |
|---|---|---|---|
| VF01 | **E6001 gegen E4003:** `verbindungsende` liefert immer `PGR-E4003` | `TestFehlerantwortVorDemAbbruch`, `TestSendNachFehlerantwort` | Code `PGR-E4003`, erwartet `PGR-E6001` |
| VF02 | **Merker-Rücksetzen:** `merke(nil)` beim `ReadyForQuery` entfernt | `…VorDemAbbruch/ReadyForQuery_nach_der_Fehlerantwort` (e), `TestSendNachFehlerantwort/ErrorResponse_und_ReadyForQuery_gelesen` (h3) | `PGR-E6001`, erwartet `PGR-E4003` |
| VF03 | **erste statt letzte** `ErrorResponse` gemerkt | `…VorDemAbbruch/zwei_Fehlerantworten,_die_letzte_zählt` | „Text nennt die erste Fehlerantwort“ |
| VF04 | `Send` liest den Merker nicht (immer `PGR-E4003`) | `TestSendNachFehlerantwort/ErrorResponse_gelesen` (h1) | `PGR-E4003`, erwartet `PGR-E6001` |
| VF05 | **F-390:** Markierung in `ClientMessage` entfernt | `TestRecordExtendedSendefehlerNachFehlerantwort/Client-Richtung_zuerst` (i1) | „Session übernommen“, mit der vorherigen Interaktion `SELECT 1` |
| VF06 | nur `FATAL` gemerkt | `…VorDemAbbruch/Schweregrad_ERROR`, `…/Schweregrad_PANIC` | `PGR-E4003`, erwartet `PGR-E6001` |
| VF07 | SQLSTATE nicht im Meldungstext | sechs Fälle von `…VorDemAbbruch`, `TestSendNachFehlerantwort/ErrorResponse_gelesen` | Text ohne `57P01` |
| VF08 | **Session verwerfen:** Markierung in `AwaitServer` entfernt | `TestRecordExtendedFehlerantwortVorDemAbbruch` (g), `…SendefehlerNachFehlerantwort/Server-Richtung_zuerst` (i2) | Session in der Aufzeichnung |
| VF09 | unbekannter Feldcode mit leerem Wert übersprungen | `TestDiagnosefelder` | Felder ohne `"Y": ""` |
| VF10 | `AwaitServer` liest ohne laufende Interaktion (Wartebedingung `len(offen) > 0 && …`) | `TestRecordNachDerLetztenInteraktion` | „AwaitServer ohne laufende Interaktion zurückgekehrt: … [PGR-E1000]“ |
| VF11 | Hinweise einer aufgezeichneten Lebendprüfung beim Laden der nächsten Interaktion vorangestellt | `TestReplayLebendpruefungOhneHinweise` | „Interaktion nach der Lebendprüfung“: `notice_response`, `parameter_status` vor `command_complete B` |
| VF12 | **E2E, Replay:** `ReadyForQuery E` als `I` geliefert | `TestE2EFehlerreplayEinfach` | „Sicht im Replay“ ungleich direkt; `TestE2EErgebnisartenEinfach` grün, erwartet (kein `E`) |
| VF13 | **E2E, Replay:** `NoticeResponse` weggelassen | `TestE2EErgebnisartenEinfach` | „Sicht im Replay“; `TestE2EFehlerreplayEinfach` grün, erwartet |
| VF14 | **E2E, Replay:** SQLSTATE jeder `ErrorResponse` als `XX000` | `TestE2EFehlerreplayEinfach` | „Sicht im Replay“ |
| VF15 | **E2E, Replay:** Antworten vor einer `ErrorResponse` weggelassen | `TestE2EFehlerreplayEinfach` | „Sicht im Replay“ |
| VF16 | **E2E, Record:** an den Client ohne `ParameterStatus`, die Aufzeichnung vollständig | `TestE2EErgebnisartenEinfach` | „Sicht beim Aufzeichnen“ ungleich direkt |
| VF17 | **E2E, Record:** an den Client `ReadyForQuery E` als `I` | `TestE2EFehlerreplayEinfach` | „Sicht beim Aufzeichnen“ |
| VF18 | **Weitergabe:** Driving-Adapter schreibt die Teilantworten vor `PGR-E6001` | `TestE2ERecordWeitergabeNichtUnterstuetzt` | „erwartet nur PGR-E6001, erhalten 1 Ergebnisse“ |
| VF19 | **E2E, Replay:** zwei aufeinanderfolgende `NoticeResponse` vertauscht | `TestE2EErgebnisartenEinfach` | „Sicht im Replay“ |
| VF20 | `base64`-Wert beim Lesen der YAML-Aufzeichnung um ein Byte gekürzt | `TestE2EErgebnisartenEinfach` | „Sicht im Replay“ |
| VF21 | **F-391:** jeder Lesefehler gilt als Verbindungsende | `TestNichtLesbareServerantwort` j1, j2, j3 | `PGR-E4003` bzw. Text mit `57P01` |
| VF22 | `CloseSession` liest `unsupported` nicht | `TestRecordExtendedNichtDarstellbar`, `…FehlerantwortVorDemAbbruch`, beide i, `TestRecordSessionNichtUnterstuetzt` | Session in der Aufzeichnung |
| VF23 | **F-391:** Lesefehler ohne Verbindungsende nach dem Merker beschriftet | `TestNichtLesbareServerantwort/unbekannter_Typ_nach_ErrorResponse` (j2) | Text nennt die Fehlerantwort statt „nicht lesbar“ |
| VF24 | **Szenario 4:** Bei abweichender Anfrage liefert das Replay die Antworten der aufgezeichneten Anfrage ohne `ReadyForQuery`, der Driving-Adapter sendet sie vor `PGR-E5001` | **grün**: `TestE2EReplayAbweichung`, alle Unit-Tests in `services` und `driving/pgwire` | Die Sonde im Mutanten zeigt: Der Client erhält auf `SELECT 2;` ein Ergebnis mit der Zeile `1`, danach `PGR-E5001`. Siehe V-52. |

Das Sendefehler-Fenster aus 11f und 11g habe ich nicht mutiert. Es ist dort ein akzeptiertes Negativ ohne Test: Eine `ErrorResponse`, die beim Sendefehler noch ungelesen ist, zählt nicht.

---

## 3. Closure-Trigger der Welle: Abnahmeszenario 4 und 6 für Simple Query

- **Szenario 6** steckt in `TestE2EFehlerreplayEinfach` und trägt (Punkt 1). Extended deckt `TestE2EReplayExtendedPgx` aus `welle-extended-query` ab.
- **Szenario 4** steckt in einem Test: `TestE2EReplayAbweichung` in `test/integration/replay_e2e_test.go` (LH-FA-10/Happy). Er nutzt `pgconn.Exec`, also das Simple Query Protocol, und prüft `PGR-E5001` und Exit-Code 5. Seine Deklaration sagt aber auch „und keine Antwort“. Das prüft er nicht, die Ergebnisse von `ReadAll` verwirft er (`_`). VF24 bleibt grün. Damit hält der Test die Hälfte von Szenario 4 nicht, die das Lastenheft ausdrücklich nennt: „statt eine unpassende aufgezeichnete Antwort zu verwenden“ (LH-FA-10 Negative). Siehe V-52.

**Urteil:** Szenario 6 trägt den Trigger. Szenario 4 trägt ihn nur zur Hälfte. Der Mangel liegt außerhalb dieses Slice (Bestand aus dem Walking Skeleton). Für die Welle-Closure ist er aber offen.

---

## 4. Plan gegen Code (`AGENTS.md` §3.9)

| Stelle | Ergebnis |
|---|---|
| Kopf `Bezug` | nennt jetzt LH-FA-02 und LH-FA-05 (F-396 behoben). Diese Ergänzung hält kein Sensor, siehe V-55. |
| Kopf `Berührte Spec-Stellen` | deckt die Stellen, die der Diff in der Spezifikation ändert (`LH-FA-02.b`, `LH-FA-05.a`, `LH-FA-11.a`). `make kopf-check` ist grün. |
| §1 | Ziel, Zuschnitt (drei Festlegungen ohne Code, eine mit Code in Upstream-Adapter und Record-Service) und Abgrenzung stimmen mit dem Diff überein. Weiterer Produktionscode liegt nicht im Diff. |
| §2 | DoD 3 zählt vier Festlegungen (F-396 behoben) und nennt 11a bis 11h. |
| §3 | Jede Zeile hat ihren Gegenstand im Diff: Merker unter eigenem Mutex, `Send` über `verbindungsende`, `istVerbindungsende`, Markierung in `ClientMessage` und `AwaitServer` über `merkeNichtUnterstuetzt`, Tests a bis f, h1 bis h3, i1/i2, j1 bis j3, die Fünf-Sekunden-Frist des Fake-Servers, das neue Feld `empfangsFehler` der Fake-Session, die Folge-Slices und die Abdeckung. |
| §6 Zeilen 10, 11f, 11g, 11h | Der Code folgt jeder Entscheidung. Zeile 10: VF09. 11f: `Query` liefert bei gescheitertem `Flush` unverändert `PGR-E4003`. 11g: VF04, VF05, VF08. 11h: VF21, VF23, die Sonde zum Text. Die Vorgabe für i1 beschreibt eine andere Reihenfolge als der Test (V-53). |
| Spezifikation | `LH-FA-02.b` *Gelesen* und *Verbindungsende* geben 11g und 11h wieder, die Historie trägt die Zeilen. |
| Folge-Slice `slice-v1-abschluss-einspielen` | §1 vermerkt die Bereinigung von `LH-FA-20.a` Schritt 3. Die DoD sagt „eine aufgezeichnete Interaktion ohne `ReadyForQuery` wird wie jede andere eingespielt“ nicht mehr zu. Der Kopf nennt `LH-FA-02.b`. Stimmig (F-395 behoben). |
| Folge-Slice `slice-v1-abschluss-antwortvergleich` | §6 trägt das Risiko zu `LH-FA-24.a`. Aus der DoD sind unvollständige Aufzeichnung, aufgezeichnetes Verbindungsende und „Aufzeichnung ohne `ReadyForQuery`“ gestrichen. Der Kopf nennt `LH-FA-02.b`. Stimmig (F-395 behoben). |
| Welle-Plan §4 | Die Zeile des Slice nennt weiter nur LH-FA-11, LH-FA-06 und LH-FA-12 (V-54). |

---

## 5. Status der Findings aus dem Review

| Finding | Status | Beleg |
|---|---|---|
| F-390 | behoben | Code (`Send` über `verbindungsende`, Markierung in `ClientMessage`), Spezifikation *Gelesen*, Tests h1 bis h3, i1, i2; VF04, VF05, VF08 rot |
| F-391 | behoben | `istVerbindungsende`, Spezifikation *Verbindungsende*, Tests j1 bis j3; VF21, VF23 rot |
| F-392 | behoben | `TestRecordNachDerLetztenInteraktion`, `TestReplayLebendpruefungOhneHinweise`; VF10, VF11 rot. §6 Zeile 1 nennt beide Tests. |
| F-393 | behoben | Fall „zwei Fehlerantworten“; VF03 rot |
| F-394 | behoben | Die Deklaration sagt jetzt: „in einer Folge von Interaktionen, die mit fehlerfreien beginnt und endet“. Die Abdeckung ist neu geschrieben, `make abdeckung-check` grün. |
| F-395 | behoben | DoD und Kopf beider Folge-Slices nachgezogen (Abschnitt 4) |
| F-396 | behoben | `Bezug` um LH-FA-02 und LH-FA-05 ergänzt, DoD 3 zählt vier Festlegungen |
| F-397 | erledigt (INFO) | Test zum unbekannten Feldcode nachgereicht; VF09 rot |
| F-398 | erledigt (INFO) | `lesbar` fängt ein Panic beim Zerlegen ab und beendet den Test mit Diagnose, `leseNachricht` lehnt Längen unter 4 ab |

---

## 6. Befunde

| ID | Befund | Ort | Empfehlung an den Planner |
|---|---|---|---|
| V-52 | **Szenario 4 hält „keine Antwort“ nicht.** `TestE2EReplayAbweichung` deklariert „erhält einen eindeutigen Fehler mit PGR-E5001 und keine Antwort“, verwirft aber die Ergebnisse von `ReadAll`. Kein Test hält für eine einfache Anfrage die Zusage aus LH-FA-10 Negative, dass keine aufgezeichnete Antwort einer anderen Anfrage geliefert wird. `TestReplayMismatch` deklariert LH-FA-10/Negative, prüft aber nur Code und Cursor. VF24 bleibt in E2E und Unit grün. Das Verhalten im heutigen Code ist richtig: `Query` liefert bei Abweichung `nil`, und der Adapter sendet bei Fehler nichts. Nur der Nachweis fehlt. Klasse `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. Dieser Slice ist nicht der Ort, der Closure-Trigger der Welle aber hängt daran. | `test/integration/replay_e2e_test.go:102`–`:118`; `internal/hexagon/services/replay_test.go` (`TestReplayMismatch`); `docs/user/abdeckung-e2e.md:31` | Vor der Welle-Closure einen Test nachziehen, der unter VF24 rot wird (E2E prüft `len(results) == 0`, Unit prüft `out == nil` bei Abweichung). Der Ort ist ein Folge-Slice der Welle, etwa `slice-replay-semantik-meldungscodes`, oder ein eigener kleiner Slice. |
| V-53 | **Vorgabe i1 und Test in anderer Reihenfolge.** §6 *Vorgabe für die Zeilen 11g und 11h*, i1: „Erst danach bekommt `Receive` `PGR-E6001`, und `CloseSession(EndFailed)` folgt.“ Der Test ruft `CloseSession` vor dem Fehler aus `Receive`, und `AwaitServer` liefert dann `ErrSessionEnded`. Die Reihenfolge des Tests ist die strengere, denn sie trifft genau das Rennen aus F-390 (`CloseSession` entscheidet, bevor die Server-Richtung markiert), und VF05 färbt sie rot. Der Plan beschreibt aber nicht, was geprüft ist. | `docs/plan/planning/in-progress/slice-replay-semantik-fehlerreplay.md` §6, i1; `internal/hexagon/services/record_extended_test.go:329`–`:349` | i1 in §6 an den Test angleichen (§3.9). Am Code ist nichts zu ändern. |
| V-54 | **Welle-Plan nennt die neuen Bezüge nicht.** Die Zeile des Slice in §4 von `welle-replay-semantik` nennt LH-FA-11, LH-FA-06 und LH-FA-12. Der Kopf des Slice nennt seit F-396 auch LH-FA-02 und LH-FA-05. | `docs/plan/planning/welle-replay-semantik.md` §4 | Bei der Closure die Zeile nachziehen oder bewusst bei den drei Kern-Anforderungen lassen. |
| V-55 | **Die Ergänzung aus F-396 hält kein Sensor.** Als Sonde habe ich LH-FA-02 in einer Kopie aus `Bezug` entfernt, und `kopf-check` bleibt grün. §1 und §2 nennen nur `LH-FA-02.b`, und das steht in `Berührte Spec-Stellen`. Das liegt innerhalb der dokumentierten Grenze des Sensors (Haupt- und Unterkennung decken einander nicht, geprüft wird gegen die Vereinigung beider Felder). Ein Mangel des Sensors ist es nicht. | `tools/harness/kopf-check.sh` (Kopfkommentar, ZUSAGE und GRENZE) | Nur vermerkt. Einen eigenen Folgeauftrag braucht es nicht. |

**Grenzen dieser Verifikation:**

- Den Race-Detector habe ich nicht gefahren. Die Stufen bauen mit `CGO_ENABLED=0`, und das Image enthält keinen C-Compiler. Dass `fehler` nur unter `merker` gelesen und geschrieben wird, habe ich am Code gelesen (`merke`, `verbindungsende`), nicht gemessen.
- Exit-Code 6 für die Fehlerantwort vor dem Abbruch ist nur über die Klasse belegt (`note` und Code `PGR-E6001`), nicht in einem E2E-Lauf. Die DoD sagt das auch nicht zu.

---

## 7. Gesamturteil

Alle drei DoD-Liefer-Punkte und `make gates` sind **bestätigt**. Jede Teilbehauptung hat einen Test, der unter seiner Mutation aus dem richtigen Grund rot wird. F-390 bis F-396 sind behoben. Die Vorgaben aus §6 in den Zeilen 10 und 11f bis 11h hält der Code. Der Plan folgt dem Code bis auf die Reihenfolge in i1 (V-53). Für die Closure der Welle ist V-52 offen: Szenario 4 ist für Simple Query nur zur Hälfte nachgewiesen.

**Summary:** 3/3 Liefer-Punkte bestätigt · `make gates` grün auf `ce50a10` · 22 von 23 Mutationen rot, VF24 grün (V-52, außerhalb des Slice) · Befunde V-52 (Closure-Trigger der Welle), V-53 (Plan §6 i1), V-54, V-55 (vermerkt).
