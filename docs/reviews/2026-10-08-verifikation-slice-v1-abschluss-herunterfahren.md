# Verifikation: slice-v1-abschluss-herunterfahren — 2026-10-08

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-08-review-slice-v1-abschluss-herunterfahren.md`, F-485 bis F-495); F-489 und F-493 bis F-495 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-herunterfahren.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `d1bb8e4..1b6e77c`. Architect: `7437aec`, `d1bb8e4`, `de79b71`, `9a5a87e`. Implementer: `17e4925`, `64f48b3`, `61ad648`, `5e20a87`, `c644fc3`. Review: `5843201`. Nacharbeit: `77e203c`, `38c4a9c`, `1b6e77c`.

**Eingang:**

- Plan ganz, besonders §1, §2 (DoD), §3, §6 (*Randformen*, *Risiken*), §7 (*Belege des Implementers*) und den Block *Nacharbeit zum Review* (steht in §8, siehe V-109)
- `spec/lastenheft.md` `LH-FA-13`; `spec/spezifikation.md` `LH-FA-13.a`, `LH-FA-13.b`, `LH-FA-14.a`, `LH-FA-17.a` *Dauer*, `LH-FA-03.b` *Zeitpunkte* und *Fehlerebene*, `SPEC-038` *Warten in Tests*, `SPEC-046`, `SPEC-051`
- `AGENTS.md` §3.9 bis §3.12
- Am Stand `1b6e77c` der Produkt-Diff ganz (`cmd/pgwire-recorder/main.go`, `internal/bootstrap/bootstrap.go`, `internal/adapters/driving/cli/cli.go`, `internal/adapters/driving/pgwire/server.go`, `internal/hexagon/{model,ports/driving,services}`), die neuen Tests (`server_zwangsende_test.go`, `record_zwangsende_test.go`, `replay_zwangsende_test.go`, `internal/bootstrap/frist_test.go`, `exitcode_test.go`, `internal/adapters/driving/cli/frist_test.go`, `test/integration/herunterfahren_e2e_test.go`) und die Helfer in `test/integration/record_e2e_test.go`
- `docs/user/benutzerhandbuch.md` (Diff), Register `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`
- Der Review-Report ganz, den Bericht des Implementers nicht

Bei meinem Start war der Arbeitsbaum sauber, HEAD `1b6e77c`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-08

**Wie die Proben gebaut wurden.** Je Probe eine frische Kopie per `cp -r` ohne `-p` (neue mtime, neuer Pfad je Mutant) in einem eigenen Unterverzeichnis des Scratchpads, `.git` darin entfernt, nie im Repo. Die Mutation per Skript, das abbricht, wenn der Suchtext nicht genau einmal trifft. Unit-Mutanten: Stufe `source` der Kopie mit eigenem Tag, `go test -count=1 -timeout 120s` in einem Container mit `--network none`. Integrations-Mutanten und der Betriebslauf: Stufe `integration` bzw. `runtime` der Kopie mit eigenem Tag, eigenes `--internal`-Netz, eigener PostgreSQL-Container (Image und Digest aus `harness/mk/integration.mk`), eigenes Volume. Danach Container, Netz, Volume, Image und Kopie entfernt; danach stand nichts mit meinen Lauf-Namen mehr. `make gates` und `make lint` liefen im Arbeitsbaum; `git status` war danach sauber.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `record` mit Frist. **Bestätigt bis auf V-106 und V-108.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| wartet höchstens `--shutdown-timeout` | Betriebslauf (Abschnitt 3): Frist `3s` → Ende nach 3,37 s, Standard unter `docker stop` → Ende nach 5,32 s, Exit-Code 4, nicht 137. Im Test nur „endet überhaupt“: R18/R37 (Zeitgeber nie gesetzt bzw. Frist 0) rot; **V2 (zehnfache Frist) und V4 (ein Zehntel) grün** in allen Unit- und den fünf neuen Integrationstests | Verhalten bestätigt, Zusage ungeprüft (V-106) |
| ab dem ersten Signal | Code: Zeitgeber entsteht in `betreiben` nach `<-ctx.Done()`. **V1 (Zeitgeber vor dem Warten auf das Signal) grün** in allen Unit- und den fünf Integrationstests; §7 führt keine Zeile zu dieser Zusage | Verhalten gelesen, Zusage ungeprüft (V-106) |
| `0` ohne Frist | R19 (Tabelle §7); `TestRunRecordZweitesSignalOhneFrist` hält „läuft nach 1 s noch“; Betriebslauf mit `0`: nach 2 s läuft der Prozess noch | bestätigt |
| Zwangsende wie Abbruch nach `LH-FA-02.b`, Aufzeichnung geschrieben | Betriebslauf: Aufzeichnung trägt `SELECT 1` als Sequenz 1, `pg_sleep` fehlt; E2E `nurErsteInteraktion`; R39 | bestätigt |
| Meldung je Verbindung (Session und verworfene Interaktion) | Betriebslauf: `level=ERROR … code=PGR-E4006 … Interaktion 2 nicht abgeschlossen und verworfen; die Session wird als Session 1 geschrieben`, psql erhält `FATAL: Netzwerk [PGR-E4006]: …`. R02, R03, R05, M1; eigene Probe V3 (Zwangsende trifft nur das erste Ziel): rot (`TestZwangsendeNebeneinander`) | bestätigt |
| `PGR-E4006` nur bei begonnener Interaktion | R01; Betriebslauf ohne Verbindung ohne `ERROR`; `TestRunRecordFristLaeuftAb` (Aufbau, Exit-Code 0, keine `ERROR`-Zeile) | bestätigt |
| Exit-Code 4 | E2E `TestE2ERecordFristLaeuftAb`; Betriebslauf | bestätigt |
| zweites Signal lässt die Frist sofort ablaufen | R20, R35, R36; Betriebslauf (`0`, zweites `SIGTERM` nach 2 s, Ende 0,33 s danach, Exit-Code 4) | bestätigt |
| weitere Signale ohne Wirkung (Test mit drei Signalen) | **R34 nachgefahren: 3 von 3 Läufen rot**, je `Exit-Code -1, erwartet 4 (signal: terminated)`, 2,4 s; ohne Mutant im Gate grün (3,24 s) | bestätigt |
| Info-Zeile `sessions` | R22 bis R24; Betriebslauf `sessions=1` bzw. `sessions=0` | bestätigt |
| ohne offene Verbindung sofort, `record` schreibt | R21; Betriebslauf mit Frist `60s`: Ende nach 0,22 s, Exit-Code 0, Datei mit `sessions: []` | bestätigt |
| Randfälle `LH-FA-13.a` (Startphase, Aufbau, Was die Frist begrenzt) | R09 bis R15, R25, R42; G01 als Grenze in `LH-FA-13.a` *Startphase* (Architect `9a5a87e`) | bestätigt |
| Form der Dauer (`LH-FA-17.a`) | R27 bis R33; **M11 nachgefahren: rot** (`05s`, `00s` abgelehnt statt 5 s bzw. 0); `00` → `PGR-E2001` in `TestParseShutdownTimeoutWerte` | bestätigt |
| `meldeFrist` = `SPEC-051`, Schranke als Literal | `TestFehlerantwortMitFrist` mit `2 * time.Second` als Literal, Konstante `meldeFrist = time.Second`; R17 | bestätigt |
| Unit-Test im Kern für die Einstufung | `TestRecordZwangsende` u. a. in `internal/hexagon/services` | bestätigt |
| Nach `SIGINT` oder `SIGTERM` | Kein Test sendet `SIGINT` (`grep` über `test/` und `internal/`: nur `cli.go` Hilfetext und `main.go`) | Zusage ungeprüft (V-108) |

### Punkt 2: `replay` mit Frist. **Bestätigt bis auf V-106.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Frist begrenzt das Warten | P18 (Frist 0); `replay` läuft durch dasselbe `betreiben`, V1, V2 und V4 gelten gleich (`TestE2EReplayFristLaeuftAb` 10,25 s unter V2, grün) | wie V-106 |
| schließt die Client-Verbindung, unvollständige Interaktion ist `PGR-E4006`, Exit-Code 4 | E2E `TestE2EReplayFristLaeuftAb`; P01 bis P15 | bestätigt |
| `PGR-E4006` vor `PGR-E5002`, Exit-Code 4 | **P19 nachgefahren (Unit): rot** — `TestReplayZwangsendeVorNichtVerbraucht` (`erster Fehler "PGR-E5002"`), dazu vier weitere Tests über die Reihenfolge `forced close`; E2E `TestE2EReplayFristVorNichtVerbraucht` im Gate grün | bestätigt |
| Test über `pgwire.Handle` mit Replay-Fake (V-95) | `TestReplayZwangsende` u. a. in `server_zwangsende_test.go` (Paket `pgwire_test`) | bestätigt |
| Unit-Test im Kern für die Einstufung | `replay_zwangsende_test.go` | bestätigt |

### Punkt 3: Exit-Code und Handbuch. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Zuordnung Code → Exit-Code je Klasse aus `SPEC-013` bis `SPEC-019` | `TestExitCodeJeKlasse`: `""`→0, `PGR-E1000`→1, `PGR-E2001`→2, `PGR-E3001`→3, `PGR-E4006`/`PGR-E4003`→4, `PGR-E5001`→5, `PGR-E6001`→6 | bestätigt |
| 3 mit Vorrang, sonst gemerkte Klasse | X04; **X05 nachgefahren: rot** (`TestRunRecordSchreibfehlerVorGemerkterKlasse`, Exit-Code 0 statt 4) | bestätigt |
| je Exit-Code 0, 3, 4, 5, 6 ein Integrationstest nach Signal | gelesen: alle sieben in §7 genannten Tests enden über `stop` (SIGTERM, `warteEnde` 15 s, `pruefeExit`) mit 0, 3, 4, 4, 5, 5/0, 6; neu mit `PGR-E4006` die beiden E2E | bestätigt |
| Handbuch: Herunterfahren mit Frist in `record` und `replay`, `replay` im Container | Abschnitt *Herunterfahren mit Frist*, Hinweise unter *Eine Anwendung aufzeichnen* und *Wiedergeben*, Optionstabelle, Log-Stufen, Exit-Codes, Meldungstabelle; Container im Replay als Hinweis (F-488). Der Betriebslauf verhält sich wie dort beschrieben (Abschnitt 3) | bestätigt, V-111 |
| Beleg in §7 je Zusage: Zusage · Mutation · roter Test | Tabellen R01–R42, P01–P19, X01–X05 (ohne X03, akzeptiertes Negativ in §6), H1–H4, M1–M3, M11 | vorhanden bis auf V-106, V-108 |

### Punkt 4: `make gates` grün. **Bestätigt** (Abschnitt 4).

### Übrige DoD-Punkte (konstant je Slice)

- Review: liegt vor (`5843201`); F-485 bis F-492 haben einen Ausgang (Abschnitt 2).
- Closure-Notiz, Register, Risiko-Ausgänge, Paarungen: in §7 noch Platzhalter, §6 *Risiken* mit „offen bis Closure“. Das ist Closure-Arbeit nach dieser Verifikation, kein Befund. Vorschläge in Abschnitt 6; V-107 betrifft das Register.

---

## 2. Review-Befunde F-485 bis F-495: nachgefahren

| Befund | Ausgang im Plan | Nachgefahren | Urteil |
|---|---|---|---|
| F-485 Tests ohne eigene Frist | H1 bis H4 | **H1 rot nach 2,00 s** (`OpenSession nach der Startnachricht bleibt binnen 2s aus`, beide Tests), kein Hänger; alle fünf Stellen aus F-485 warten jetzt in `select` mit Literal-Frist (gelesen) | umgesetzt |
| F-486 Zweig „anderen Fehler nur merken“ | M2, M3, `TestRecordZwangsendeAndererFehler` | **M2 rot** (`erster Fehler ""`) | umgesetzt |
| F-487 dritte Meldungsform | `LH-FA-13.a` *Meldung* (Architect `9a5a87e`), `nichtGeschrieben`, M1 | gelesen: zwei Formen im Code wie in der Spezifikation | umgesetzt |
| F-488 Handbuch weiter als Prüfung | `38c4a9c` | gelesen: Bedingung für `PGR-E4006` und Exit-Code genannt, Container im Replay als Hinweis | umgesetzt, Rest V-111 |
| F-489 Startphasen-Fenster | Grenze in `LH-FA-13.a` *Startphase*, §6, G01 in §7 als „Fenster offen, als Grenze entschieden“ | gelesen | umgesetzt |
| F-490 führende Nullen | `LH-FA-17.a` *Dauer* (`de79b71`), M11 | **M11 rot** | umgesetzt |
| F-491 Upstream im Aufbau | Grenze in `LH-FA-13.a` *Zwangsende*, §6 | gelesen | umgesetzt |
| F-492 X03 äquivalent | §6, akzeptiertes Negativ | gelesen | umgesetzt |
| F-493 Hexagon | — | Abschnitt 5 | bestätigt |
| F-494 Nebenläufigkeit | — | R34 3 von 3 rot; Betriebslauf mit zweitem Signal | bestätigt |
| F-495 Größe | — | drei Liefer-Punkte, zwei Schichten (Driving mit Composition Root, Kern), kein Driven-Adapter im Diff; Rückführung aus §4 nicht eingetreten | bestätigt |

---

## 3. Betriebslauf gegen PostgreSQL (`LH-FA-13`, `LH-FA-02`)

Produkt-Image (Stufe `runtime`) einer frischen Kopie, `record --listen 0.0.0.0:5432 --upstream <pg>:5432 --output /rec/<name>.yaml`, Client `psql -c "SELECT 1" -c "SELECT pg_sleep(30)"` aus dem PostgreSQL-Image (Standard `sslmode=prefer`), Signal 2 s nach dem Start des Clients. Gemessen vom Signal bis zum Ende des Containers.

| Szenario | Ende | Exit-Code | `stderr` | Client | Aufzeichnung |
|---|---|---|---|---|---|
| `--shutdown-timeout 3s`, ein `SIGTERM` | 3,37 s | 4 | `sessions=1`; nach 3,01 s `ERROR code=PGR-E4006 … Interaktion 2 … als Session 1 geschrieben`; `record beendet` | `FATAL: Netzwerk [PGR-E4006]: …` | Session 1, Sequenz 1 `SELECT 1` |
| Standard, `docker stop` (10 s bis `SIGKILL`) | 5,32 s | 4 (nicht 137) | wie oben, nach 5,02 s | wie oben | wie oben |
| `--shutdown-timeout 60s`, keine Verbindung | 0,22 s | 0 | `sessions=0`, `record beendet` | — | `sessions: []` |
| `--shutdown-timeout 0`, zwei `SIGTERM` im Abstand von 2 s | 2,33 s (nach 2 s läuft er noch) | 4 | wie oben | wie oben | wie oben |

Das Binary verhält sich so, wie das Handbuch es einem Betreiber beschreibt: Ende spätestens nach der Frist, die abgeschlossene Anfrage bleibt, die laufende fehlt, der Client erhält `PGR-E4006`, ohne Verbindung endet es sofort und schreibt, das zweite Signal beendet das Warten ohne Frist. Der Standardwert hält unter der Stopp-Frist von Docker: Exit-Code 4 statt 137.

---

## 4. Gate-Ergebnis

Am Stand `1b6e77c` im Arbeitsbaum: `make gates` Exit 0 (2 min 40 s), darin `make test-integration` grün (52 Zeilen `--- PASS`, keine `--- FAIL`, die fünf neuen E2E mit 1,23–3,24 s), `make a-check` 0 Befunde, `make a-check-negativ` grün, `make docs-check` 0 Befunde, `make abdeckung-check`, `make kopf-check` und die Gegenproben grün. `make lint` Exit 0.

---

## 5. Plan gegen Code

- **§1 Schicht-Abgrenzung:** Der Diff berührt `cmd/pgwire-recorder`, `internal/bootstrap`, `internal/adapters/driving/{cli,pgwire}` und `internal/hexagon/{model,ports/driving,services}`; keine Zeile in `internal/adapters/driven/`. Wie §1.
- **§3:** jede Zeile hat ihre Änderung im Diff; nichts Produktives außerhalb von §3. Die Spezifikation hat nur der Architect geändert.
- **Hexagon (F-493):** `PGR-E4006` entsteht im Use Case — `RecordService.CloseSession` aus `EndForced` und `laeuft()`, `ReplayService.Forced` aus `ungesendet` und `mitten()`. Der Adapter meldet das Ereignis (`EndForced`, `Forced`) und liest den Code nur, um zwischen Zustellen und Merken zu wählen. Keine neuen Imports im Kern; `make a-check` 0 Befunde, `make a-check-negativ` grün.
- **§6 Randformen:** jede Randform im Code so wie entschieden (gelesen an `betreiben`, `Zwangsende`, `oeffne`, `richtungen.zwangsende`, `replaySitzung`, `dauer.Set`). Den Beginn der Frist hält der Code; ein Test fehlt (V-106).
- **§7 Grüne Mutanten:** G01, die zwei äquivalenten, P10 und X03 wie eingeordnet; G01 jetzt als Grenze, nicht als unerreichbar (F-489).
- **`SPEC-038` *Warten in Tests* in den neuen Tests:** Jede Stelle, die auf ein Ereignis des Prüflings wartet, wartet in `select` mit Literal-Frist (`ereignisBinnen`, `zurueckBinnen`, `endetBinnen`, `laeuftNoch`, `warteEnde`, `warteAbgelehnt`, Lesefristen 5 s), oder über einen `context.WithTimeout` (pgconn in den E2E). Kanäle, auf die ein Fake schreibt, sind gepuffert (`empfang` 16, `server` 1). Fünf Prüfungen „keine weitere Nachricht“ lesen mit `fe.Receive()` gegen die Pipe-Frist von 10 s und nehmen deren Ablauf als Erfolg; jede steht neben einem `zurueckBinnen` mit 2 s, das das Ausbleiben des Endes rot macht. Kein Befund. Keiner der neuen Tests liest `stderr`, solange der Prozess läuft (§6).
- **Die Hilfsfunktion `session` in `internal/hexagon/services/record_test.go`:** siehe V-110.

---

## 6. Vorschläge für die Closure (Planner)

- **Risiko *Herunterfahren ohne Obergrenze* (Record)** und **(Replay):** *entfallen*, Grund: Frist verdrahtet, E2E je Modus und der Betriebslauf hier — erst nachdem V-106 einen Test hat, der die Obergrenze hält.
- **Risiko *Test mit Signal hängt bis zum Zeitlimit*:** *eingetreten* — F-485 (Mutant M-hang hing bis `panic: test timed out after 20s`), dazu der erste Anlauf zu H4 in `session`. Ausgang über das Register, siehe V-107.
- **Register:** V-107; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-486, F-488) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (F-487, F-489 bis F-491) aus der Summary-Zeile des Reviews.

---

## 7. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-106 | MEDIUM (DoD Punkt 1 und 2; `LH-FA-13` *Boundary* „spätestens nach einer einstellbaren Frist“; `LH-FA-13.a` „zählt ab dem ersten Signal“; `AGENTS.md` §3.10) | Beginn und Länge der Frist hält kein Test. Die Tests prüfen nur, dass der Lauf überhaupt endet: `TestRunRecordFristLaeuftAb` gibt 200 ms Frist 3 s Zeit, die E2E geben 1 s Frist 15 s. Drei eigene Mutanten bleiben in allen Unit-Paketen und in allen fünf neuen Integrationstests grün: **V1** Zeitgeber vor `<-ctx.Done()` (Frist ab Start von `betreiben` statt ab dem ersten Signal), **V2** `time.NewTimer(10 * frist)` (`TestE2ERecordFristLaeuftAb` 10,24 s, `TestE2EReplayFristLaeuftAb` 10,25 s — mit dem Standardwert hieße das 50 s, also `SIGKILL` nach der Stopp-Frist von Docker; gerechnet, nicht gefahren), **V4** `time.NewTimer(frist / 10)`. §7 führt zu „ab dem ersten Signal“ keine Zeile; R18 und R37 fangen nur „Zeitgeber nie gesetzt“ bzw. „Frist 0“. Das Verhalten selbst ist richtig (Abschnitt 3: 3,37 s bei `3s`, 5,32 s beim Standard). Vorschlag: ein Test in `internal/bootstrap` mit einer Frist von einigen 100 ms, der vor dem Signal länger wartet als die Frist, danach „läuft noch“ für einen Teil der Frist und „endet binnen“ Frist plus Rand prüft; rot unter V1, V2 und V4. Adresse Implementer, vor der Closure. | `internal/bootstrap/frist_test.go` · `r.endetBinnen(t, 3*time.Second, "Frist 200ms")`; `test/integration/herunterfahren_e2e_test.go` · `rec.warteEnde(t, 15*time.Second, "Recorder endet nicht nach Ablauf der Frist von 1s")`; `internal/bootstrap/bootstrap.go` · `t := time.NewTimer(frist)` | V1, V2, V4 je Unit und Integration, grün |
| V-107 | MEDIUM (Closure, Register; Plan §8) | `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` steht bei 2× (`evidence/` mit zwei Dateien). §8 sagt: „Hinge ein Test dieses Slice bis zum Zeitlimit, erreichte der Eintrag 3×.“ Das ist eingetreten: Review-Mutant M-hang ließ `TestRecordZwangsendeImAufbau` bis `panic: test timed out after 20s` hängen (F-485), und der erste Anlauf zu H4 hing in `session` (§8 *Nacharbeit*). Mit der Evidenz dieses Slice erreicht der Eintrag die Schwelle: keine Notiz mehr, sondern eine Lücke, die einen eigenen Ausgang braucht. Die Regel `SPEC-038` *Warten in Tests* gibt es seit `slice-harness-integration-wait`; F-485 ist ihr erster Bruch danach, gefunden vom Review, nicht von einem Sensor. Adresse Planner bei der Closure. | `docs/plan/planning/observations/BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit/evidence/`; Plan §8 · „erreichte der Eintrag 3×“ | Lesen; Review F-485 |
| V-108 | LOW (DoD Punkt 1 „Nach `SIGINT` oder `SIGTERM`“; `AGENTS.md` §3.10) | Kein Test sendet `SIGINT`. Der Slice hat die Signalbehandlung in `main` neu geschrieben (`signal.Notify(signale, os.Interrupt, syscall.SIGTERM)`); ein Mutant, der `os.Interrupt` streicht, bleibt grün, denn `grep` findet `SIGINT`/`os.Interrupt` in `test/` und `internal/` nur im Hilfetext. Der Bestand vor dem Slice prüfte es ebenso wenig. Vorschlag: in einem der E2E das erste Signal als `SIGINT`, oder in §7 als akzeptiertes Negativ mit Grund. Adresse Implementer. | `cmd/pgwire-recorder/main.go` · `signal.Notify(signale, os.Interrupt, syscall.SIGTERM)` | `grep -rn "SIGINT\|os.Interrupt" test/ internal/ cmd/` |
| V-109 | LOW (Plan-Form, Closure-Notiz) | Der Block *Nacharbeit zum Review* mit H1 bis M11 und den Läufen am Stand `38c4a9c` steht in §8, hinter „**Modus-Begründungsblock:** alle berührten Sub-Areas GF.“, nicht in §7. `1b6e77c` wollte ihn „in der Closure-Notiz“. Die Belege einer Rolle gehören dorthin, wo die DoD sie sucht. Adresse Implementer: nach §7 hinter *Belege des Implementers* verschieben, kein Inhalt ändert sich. | `docs/plan/planning/in-progress/slice-v1-abschluss-herunterfahren.md` · „**Nacharbeit zum Review**“ | Lesen |
| V-110 | Hinweis (`SPEC-038` *Warten in Tests*; Bestand) | Einordnung der Hilfsfunktion `session` in `internal/hexagon/services/record_test.go`: Bestand seit `slice-walking-skeleton-record` (`bef7015`), in diesem Slice unverändert und von §1 nicht erfasst, darum kein Befund gegen ihn. Sie ruft den Prüfling **synchron** ohne Frist (`OpenSession`, `Query`, `Delivered`); hält ein Mutant `Query` an, hängt jeder Test, der sie nutzt, bis zur Zeitgrenze (erster Anlauf zu H4). `SPEC-038` nennt als Warten auf den Prüfling einen Wert auf einem Kanal, dessen Schließen und das Ende eines Prozesses; ob ein blockierender synchroner Aufruf darunter fällt, sagt es nicht — eine offene Randform. Vorgeschlagene Adresse: zuerst der Architect in `SPEC-038` (synchroner Aufruf: Frist ja oder nein), danach der Code in `slice-tests-ueberlebende-mutanten` (in `open/`, Kern-Tests in `internal/hexagon/services`, eine der zwei Schichten dort). Weil sie die Laufzeit eines Mutationslaufs trifft, ist die Randform *Laufzeit* von `slice-harness-mutation` die Rückfall-Adresse. Die Entscheidung trifft der Planner, zusammen mit V-107. | `internal/hexagon/services/record_test.go` · `func session(t *testing.T, s *services.RecordService, queries ...string) model.SessionID {` | Lesen; §8 *Nacharbeit* zu H4 |
| V-111 | Hinweis (`AGENTS.md` §3.11; Rest von F-488) | *Herunterfahren mit Frist* sagt „…erhält die Fehlermeldung `PGR-E4006`, wenn sie sie binnen einer Sekunde annimmt, und der Exit-Code ist 4“. Nach `LH-FA-13.b` ist er die Klasse des ersten gemerkten Fehlers, 3 bei einem Schreibfehler am Ende; der Abschnitt *Exit-Codes* sagt das richtig. Vorschlag: „…ist der Exit-Code 4, wenn vorher kein anderer Fehler auftrat“. Adresse Implementer, optional. | `docs/user/benutzerhandbuch.md` · „sie binnen einer Sekunde annimmt, und der Exit-Code ist 4“ | Lesen |
| V-112 | Hinweis (Beleg R34) | `TestE2ERecordDrittesSignal` stellt den Stau zum Client über `time.Sleep(2 * time.Second)` her und sendet das dritte Signal 300 ms nach dem zweiten. Staut sich die Ausgabe nicht, endet der Prozess vielleicht vor dem dritten Signal; dann wird der Test über `signal` rot („Prozess endete vor dem Signal“), er hängt nicht. Am Stand `1b6e77c` im Gate grün, R34 3 von 3 rot. Für die Closure-Notiz: Grenze des Tests neben der in §6 genannten. | `test/integration/herunterfahren_e2e_test.go` · `time.Sleep(300 * time.Millisecond)` | Lesen; Gate und R34 |

Kein Befund gegen DoD-Punkt 3 und 4, gegen §1, §3 und gegen das Hexagon.

---

## 8. Urteil

**Fertig für die Closure, sobald V-106 einen Test hat; V-107 braucht bei der Closure einen Ausgang.**

- **`record` (Punkt 1):** bestätigt und im Betriebslauf gegen PostgreSQL gesehen — Ende nach der Frist mit Exit-Code 4, `PGR-E4006` an Client und Log, abgeschlossene Anfrage aufgezeichnet, laufende verworfen, ohne Verbindung sofortiges Ende mit geschriebener Datei, zweites Signal bei Frist `0`. R34 (drittes Signal) 3 von 3 rot, M11 rot, H1 rot nach 2 s statt Hänger, M2 rot. Offen: Beginn und Länge der Frist ohne Test (V-106), `SIGINT` ohne Test (V-108).
- **`replay` (Punkt 2):** bestätigt; P19 rot (`PGR-E4006` vor `PGR-E5002`). V-106 gilt für `replay` gleich.
- **Exit-Code und Handbuch (Punkt 3):** bestätigt; X05 rot. Das Handbuch beschreibt, was der Betriebslauf zeigt; eine Formulierung ist enger zu fassen (V-111).
- **`make gates`:** grün am Stand `1b6e77c`; `make lint` Exit 0.
- **Review:** F-485 bis F-492 umgesetzt bzw. entschieden, F-493 bis F-495 bestätigt.
- **Hexagon:** Einstufung im Use Case, Adapter meldet Ereignisse, `make a-check` 0 Befunde.
- **`SPEC-038`:** in den neuen Tests eingehalten; die Bestands-Hilfsfunktion `session` braucht eine Adresse (V-110).

**Vor der Closure offen:** V-106 (Test für Beginn und Länge der Frist, Implementer), V-108 und V-109 (Implementer), V-107 und V-110 (Planner, Register und Adresse), V-111 und V-112 optional.
