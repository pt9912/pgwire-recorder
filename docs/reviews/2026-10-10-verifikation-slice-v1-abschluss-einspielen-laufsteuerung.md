# Verifikation: slice-v1-abschluss-einspielen-laufsteuerung — 2026-10-10

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-10-review-slice-v1-abschluss-einspielen-laufsteuerung.md`, F-564 bis F-569). F-566 hat er an mich übergeben; F-564 bis F-567 prüfe ich nach, F-568 und F-569 ordne ich nur ein.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-laufsteuerung.md`, Kopf und §1 bis §8, gegen den Gesamt-Diff `85caec4..785fca9`. Architect: `85caec4` (Randformen in `LH-FA-20.a`), `ee23fad` (Rückgabe des Implementers). Implementer: `1b63ef8`, `de00222`, `eda0bc9`, `785fca9`. Review: `5d22326`.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingang:**

- der Plan ganz am Stand `785fca9`, besonders §1 (*Übernimmt*, *Ausdrücklich NICHT*), §2 (DoD), §3, §4, §6 (*Randformen*, *Akzeptierte Negative* (a) bis (e), *Risiken*) und §7 (*Belege des Implementers*, *Nacharbeit zum Review*)
- `spec/spezifikation.md` `LH-FA-20.a` ganz (Optionstabelle, Schritte 6 und 7, Tabelle *Fehlerregeln*, *Randformen je Schritt* mit *Interaktion*, *Meldungen*, *Abbruchsignal*, *Exit-Code*); `LH-FA-17.a` Absatz zu Umgebungsvariablen und booleschen Werten und Abschnitt *Fehler*
- [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md), Abschnitt *Entscheidung*; [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler)
- am Stand `785fca9`: der Diff von `internal/adapters/driving/cli/cli.go`, `internal/hexagon/ports/driving/play.go` und `internal/hexagon/services/play.go`, dazu `Play`, `session` und `interaktion` ganz; unverändert, aber vorausgesetzt: `play` in `internal/bootstrap/bootstrap.go`
- `docs/user/benutzerhandbuch.md` §4 *Eine Aufzeichnung in eine Datenbank einspielen* ganz, §5 Optionstabelle und *Exit-Codes*, §7 Zeile `PGR-E4004`; `README.md` *Was kann ich heute tun?*
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen
- der Nehmer der Abgrenzung zum Vergleich, `slice-v1-abschluss-antwortvergleich`, §1 und DoD; dazu `AGENTS.md` §3.9 bis §3.13

Beim Start war der Arbeitsbaum sauber, HEAD `785fca9`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.** Alles lief unter dem Präfix `ver-lauf-` im Scratchpad, nie im Repo.

- **Binary:** eine frische Kopie per `git archive 785fca9 | tar -x`, ohne `.git` und ohne `cp -p`. Aus ihr die Stufe `runtime` des `Dockerfile` mit eigenem Tag. Ein internes Docker-Netz und ein PostgreSQL-Container mit dem gepinnten Image aus `harness/mk/integration.mk` (`trust`, ohne Passwort und ohne TLS). `play` lief je Fall als eigener Container im Netz, die Aufzeichnung schreibgeschützt eingebunden; die Wirkung las ich mit `psql` im PostgreSQL-Container aus einer frisch angelegten Datenbank je Fall. Signale sendete `docker kill --signal`, erst nachdem `pg_stat_activity` die laufende `pg_sleep`-Anfrage zeigte; das zweite Signal erst nach der Log-Zeile zum Abbruchsignal. Für das Signal im Aufbau hielt `docker pause` den Server an, sodass der Aufbau hing, bis `docker unpause` ihn freigab. Jeder `docker`-Aufruf hatte eine Frist (60 s, `play` 120 s, `docker wait` 60 s).
- **Unit-Mutanten:** je Mutant eine neue Kopie per `cp -r` ohne `-p`. Die Ersetzung machte ein Skript, das abbricht, wenn der Suchtext nicht genau einmal trifft. Danach im Image der Stufe `deps` (eigener Tag, ohne Netz) `gofmt -l ./internal` leer, `go vet ./internal/...` und `go test -count=1` für `internal/adapters/driving/cli` und `internal/hexagon/services`; danach nur das Verzeichnis des Mutanten gelöscht. Der Grundlauf ohne Mutation war grün.
- **Integrations-Mutanten:** eine Kopie wie oben, daraus die Stufe `integration` mit eigenem Tag (das Image `pgwire-recorder:integration` des Arbeitsbaums blieb unberührt), dann das Testbinary mit `-test.run '^TestE2EPlay(Fortsetzung|FortsetzungAbbruch|ErwarteterFehler|FinishSession|ZweitesSignal)$'` in meinem Netz gegen mein PostgreSQL, Frist 900 s. Der Grundlauf ohne Mutation war mit allen 17 Unterfällen grün.
- **Danach:** Container, Netz, Volumes, Kopien und eigene Images (`runtime`, `deps`, `integration` je Mutant) entfernt; `docker ps -a`, `docker network ls`, `docker volume ls` und `docker images` mit `ver-lauf` und `pgr-it` danach leer. `make gates` lief im Arbeitsbaum; `git status` war danach sauber.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Die drei Optionen aus Kommandozeile, Umgebung und `play:`; Fortsetzung mit Exit-Code 4, je `PGR-E4004` eine Zeile `error`; Verbindungsfehler bricht ab, nach früherem `PGR-E4004` mit dem Exit-Code des abbrechenden (Integrationstest). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `--continue-on-error` wirkt aus allen drei Quellen | Binary: Aufzeichnung mit Session 1 `CREATE TABLE t`, `SELECT * FROM fehlt`, `INSERT 1`, `SELECT 1/0` und Session 2 `INSERT 2`. Ohne Option: eine Zeile `PGR-E4004` (Interaktion 2), Exit 4, `t` leer. Mit `--continue-on-error`, mit `PGWIRE_RECORDER_CONTINUE_ON_ERROR=true` und mit `play: continue_on_error: true` je gleich: zwei Zeilen `PGR-E4004` (Interaktion 2 `42P01`, dann Interaktion 4 `22012`), Exit 4, `t` = `1,2` (die nächste Interaktion und die nächste Session liefen) | bestätigt |
| `--allow-recorded-errors` und `--finish-session-on-interrupt` wirken aus allen drei Quellen | Binary: siehe Punkt 2 und 3, je Kommandozeile, Umgebung und Datei | bestätigt |
| Vorrang Kommandozeile vor Umgebung vor Datei | Binary: `--continue-on-error=false` neben Umgebung `true`: Abbruch nach der ersten Zeile; Umgebung `false` neben Datei `true`: Abbruch; `--continue-on-error` ohne Wert neben Datei `false`: Fortsetzung; leere Umgebungsvariable neben Datei `true`: Fortsetzung (leer gilt als nicht gesetzt) | bestätigt |
| Boolesche Werte nach `LH-FA-17.a` | Binary: `--continue-on-error=1` und `=TRUE` je `PGR-E2001`, Exit 2. `PGWIRE_RECORDER_CONTINUE_ON_ERROR=1` auch neben gesetzter Option `PGR-E2001` mit dem Namen der Variable; ebenso `…ALLOW_RECORDED_ERRORS=yes` und `…FINISH_SESSION_ON_INTERRUPT=GEHEIMWERT`, ohne den Wert. In der Datei `continue_on_error: 1`, `allow_recorded_errors: ja`, `finish_session_on_interrupt: GEHEIMWERT` je `PGR-E2004` an `play.<schlüssel>`, ohne den Wert; `continue_on_error` auf der obersten Ebene `PGR-E2004` „unbekannter Schlüssel“. Die Meldung zu `--continue-on-error=1` nennt den Wert; das ist Bestand des allgemeinen Lesers (gleich bei `replay --fail-on-unconsumed=…`) und nach `LH-FA-17.a` *Fehler* zulässig, weil der Aufruf ihn selbst enthält | bestätigt |
| Je `PGR-E4004` eine Zeile `error`, nach dem Ende, in der Reihenfolge des Auftretens | Binary: Zeilenfolge `INFO play gestartet` · `ERROR PGR-E4004` (S1/I2) · `ERROR PGR-E4004` (S1/I4) · `INFO play beendet`; jede Zeile mit `code` und `error`, darin Session, Interaktion, SQLSTATE und `M`, keine weiteren Felder | bestätigt |
| Abbrechender `PGR-E6001` nach früheren `PGR-E4004`: Exit 6, seine Zeile zuerst | Binary: Session 2 mit `SELECT 1/0` und `COPY (SELECT 1) TO STDOUT` nach `fehlt` in Session 1, `--continue-on-error`: Zeilen `PGR-E6001` (S2/I2), `PGR-E4004` (S1/I2), `PGR-E4004` (S2/I1), Exit 6; `t` = `1`, `INSERT 2` und Session 3 liefen nicht. Ohne Option bricht schon der erste `PGR-E4004` ab, Exit 4 | bestätigt |
| Abbrechender `PGR-E4003` nach früherem `PGR-E4004`: Exit 4, seine Zeile zuerst | Binary: `SELECT pg_terminate_backend(pg_backend_pid())` in Session 2 (`FATAL` `57P01`): Zeilen `PGR-E4003`, dann `PGR-E4004`, Exit 4, Session 3 lief nicht | bestätigt |
| Verbindungsfehler bricht auch mit der Option ab | Binary: Session 2 mit `database: gibtsnicht`, `--continue-on-error`: Zeilen `PGR-E4002` (S2, `3D000`), dann `PGR-E4004` (S1/I2), Exit 4, Session 3 lief nicht | bestätigt |
| Integrationstest | `TestE2EPlayFortsetzung` (drei Quellen) und `TestE2EPlayFortsetzungAbbruch` (`PGR-E6001`, `PGR-E4003`) grün im Grundlauf und im eigenen `make gates`. Mutant **I1** (`errors.Join(append(frueher, err)...)`): **rot**, `TestE2EPlayFortsetzungAbbruch/PGR-E6001` „Exit-Code 4, Zeilen … PGR-E4004 \| ERROR PGR-E6001 …, erwartet 6“ und `/PGR-E4003` (Reihenfolge) | bestätigt |

### Punkt 2: Mit `--allow-recorded-errors` ist eine Fehlerantwort erwartet, wenn die aufgezeichnete Interaktion eine `error_response` trägt: keine Meldung, kein Abbruch, kein Beitrag zum Exit-Code; sonst `PGR-E4004` wie ohne Option (Test). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Erwartet: keine Meldung, kein Abbruch, Exit 0 | Binary: `SELECT * FROM fehlt` mit aufgezeichneter `error_response` (`XX000`), Option über Kommandozeile, `PGWIRE_RECORDER_ALLOW_RECORDED_ERRORS=true` und `play: allow_recorded_errors: true`: je keine Zeile `error`, Exit 0, `t` = `1,2` | bestätigt |
| Nicht SQLSTATE: aufgezeichnet `22012`, Server `42P01` | Binary: Exit 0, ohne Zeile | bestätigt |
| Aufgezeichnet, aber ohne Option: `PGR-E4004` | Binary: Exit 4, eine Zeile | bestätigt |
| Nicht aufgezeichnet, mit Option: `PGR-E4004` wie ohne Option | Binary: Exit 4, Abbruch; mit `--continue-on-error` dazu: Exit 4, Fortsetzung, `t` = `1,2` | bestätigt |
| Kombination: erwarteter und nicht erwarteter Fehler mit `--continue-on-error` | Binary: der erwartete ohne Zeile, der andere (`22012`) eine Zeile, Exit 4 | bestätigt |
| `FATAL` bleibt `PGR-E4003`, auch erwartet | Binary: `pg_terminate_backend` mit aufgezeichneter `error_response`, mit und ohne `--continue-on-error`: Zeile `PGR-E4003`, Exit 4, Abbruch | bestätigt |
| Test | `TestPlayErwarteterFehler`, `TestE2EPlayErwarteterFehler` grün. Mutanten **V7** (erwartet vor `FATAL`), **V8** (erwartet ohne Option), **V9** (nur die erste aufgezeichnete Antwort), **V16** (erwarteter Fehler beendet das Lesen vor `ReadyForQuery`): alle **rot**, `TestPlayErwarteterFehler` mit dem Unterfall des Mutanten (V7 „Meldungen [], erwartet [… PGR-E4003 …]“; V16 Ablauf ohne das zweite `naechste`). Integrations-Mutant **I3** (`erwartet := AllowRecordedErrors`): **rot**, `TestE2EPlayErwarteterFehler` „ohne aufgezeichnete error_response: Exit-Code 0“ | bestätigt |

### Punkt 3: `SIGINT` und `SIGTERM` mit `--finish-session-on-interrupt` nach der laufenden Session, ein erstes Signal im Aufbau nach der ganzen Session; Exit-Code 4 nach früherem Fehler, sonst 0 (Test); Beleg in §7 je Zusage. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Ohne Option: nach der laufenden Interaktion | Binary: Session 1 `CREATE`, `INSERT 1`, `pg_sleep(4)`, `INSERT 2`, Session 2 `INSERT 3`; `SIGINT` während `pg_sleep`: Zeilen `play gestartet` · `Abbruchsignal, play endet vorzeitig` · `play beendet`, Exit 0, `t` = `1` | bestätigt |
| Mit Option: nach der laufenden Session, keine neue Session | Binary: `SIGINT` mit `--finish-session-on-interrupt`, `SIGTERM` mit `PGWIRE_RECORDER_FINISH_SESSION_ON_INTERRUPT=true`, `SIGTERM` mit `play: finish_session_on_interrupt: true`: je Exit 0, `t` = `1,2`, `INSERT 3` lief nie | bestätigt |
| Erstes Signal im Aufbau: mit Option die ganze Session, ohne keine Interaktion | Binary: Server per `docker pause` angehalten, `play` gestartet, `SIGINT`, Zeile zum Abbruchsignal, dann `docker unpause`. Mit Option: Exit 0, `t` = `1,2` (Session 1 ganz, also lief das Signal im Aufbau, nicht im Start), Session 2 nicht; ohne Option: Exit 0, Tabelle `t` nicht angelegt | bestätigt |
| Exit 4 nach früherem Fehler | Binary: `fehlt` vor `pg_sleep`, `--continue-on-error`, mit und ohne Option: Zeile zum Abbruchsignal, danach `PGR-E4004`, Exit 4 (`t` = `1,2` mit, `1` ohne Option) | bestätigt |
| Fehler in der zu Ende laufenden Session zählen wie ohne Signal | Binary: mit Option nach dem Signal `COPY … TO STDOUT`: `PGR-E6001`, Exit 6 (nicht die Zeile *Abbruchsignal*), `INSERT 2` lief nicht; `fehlt` nach dem Signal: ohne `--continue-on-error` Exit 4 und Abbruch, mit ihm Exit 4 und `t` = `1,2` | bestätigt |
| Zweites Signal, ohne und mit Option: sofort, kein Fehler | Binary: `pg_sleep(30)`, zwei `SIGINT` (bzw. zwei `SIGTERM` mit der Option aus der Umgebung): Ende 0,2 bis 0,3 s nach dem zweiten Signal, ohne Zeile `error`, Exit 0, `t` = `1`. Mit Option und früherem `PGR-E4004`: genau diese eine Zeile, Exit 4; ohne Option ebenso. Zweites Signal im angehaltenen Aufbau, mit und ohne Option: Ende 0,2 s danach, Exit 0, keine Zeile | bestätigt |
| Test | `TestPlayFinishSession`, `TestPlayFinishSessionFehler`, `TestPlayErstesSignal`, `TestPlaySignalNachFehlerantwort`, `TestPlayZweitesSignalNachFehlerantwort`, `TestPlayZweitesSignalFinishSession`, `TestE2EPlayFinishSession`, `TestE2EPlayZweitesSignal` grün. Mutanten **V10** (Option in `session` nicht beachtet), **V11** (Option auch in der Schleife der Sessions), **V12** (frühere Fehler beim Signal zwischen Sessions verworfen), **V14** (nach dem zweiten Signal alle verworfen): alle **rot** mit dem genannten Test (V11: Ablauf mit `verbinde`, `anfrage D` der zweiten Session). Integrations-Mutant **I2** (= V10): **rot**, `TestE2EPlayFinishSession/Kommandozeile` und `/Umgebung_nach_früherem_Fehler` „Wirkung "1", erwartet 1,2“ | bestätigt |
| Beleg in §7: je Zusage Zusage · Mutation · roter Test | Beide Tabellen in §7 decken die Zusagen aus Punkt 1 bis 3 und nennen je Zeile Mutation und roten Test. Selbst gefahren habe ich 14 Zeilen als Unit-Mutanten (V1 bis V14), drei als Integrations-Mutanten (I1 bis I3) und beide Zeilen der Nacharbeit (M1, M19) in Unit und E2E; alle rot mit dem genannten Test und aus dem genannten Grund, keiner über Build, Vet, gofmt oder eine Gesamtfrist. Außerhalb der Tabelle dazu V15 (Aufbaufehler ohne Code), V16 und V17 (weitergereichter Fehler ohne Ort), ebenfalls rot | bestätigt |

### Punkt 4: `make gates` grün. **Bestätigt durch eigenen Lauf.**

Eigener Lauf im Arbeitsbaum am Stand `785fca9`, sauber vor und nach dem Lauf: Exit 0. Darin `baseline-verify: v6.16.0 OK`, `d-check: 445 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)` (a-check), `a-check-negativ: gruen`, `abdeckung-gegenprobe: gruen`, `commit-msg-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, die Stufen `test` und `lint` gebaut, `run-integration-tests: gruen` mit `--- PASS` für `TestE2EPlayFortsetzung`, `TestE2EPlayFortsetzungAbbruch`, `TestE2EPlayErwarteterFehler`, `TestE2EPlayFinishSession` und `TestE2EPlayZweitesSignal` mit allen 17 Unterfällen; `make kopf-check` und `make abdeckung-check` liefen ohne Ausgabe durch. Der Beleg in §7 (*Nacharbeit*, `make gates` am Stand `eda0bc9`, d-check 445 Dateien) deckt sich damit; `785fca9` ändert gegenüber `eda0bc9` nur §7.

### Übrige DoD-Punkte (konstant je Slice)

- **Review:** Der Report liegt vor (`5d22326`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen:** noch offen; §6 trägt bei allen drei Risiken „offen bis Closure“. Das ist die Reihenfolge, kein Befund. Für die Closure aus meiner Sicht: Risiko 1 (Optionen oder Signale nur mit Bootstrap) und Risiko 2 (Reihenfolge der Zeilen verlangt Bootstrap) traten nicht ein: `internal/bootstrap` ist nicht im Diff, und die Zeilenfolge und Exit-Codes am Binary entsprechen *Meldungen* und *Exit-Code* (Punkt 1). Risiko 3 (Signal zwischen Fehlerantwort und Fortsetzung) ist mit `TestPlaySignalNachFehlerantwort` belegt; V4, V6 und V12 machen ihn rot. In die Register-Fortschreibung gehören die Klassen aus der Summary-Zeile des Reviews (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` für F-564, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` für F-564 und F-565, `BEO-REPO/plan-folgt-korrektur-nicht` für F-567), dazu der Hinweis des Reviews zu §Pflege des Reviewer-Skills.

---

## 2. Übergaben des Reviews

| Befund | Prüfung | Urteil |
|---|---|---|
| F-564 zweites Signal mit Option | Neu sind `TestPlayZweitesSignalFinishSession` (drei Unterfälle) und `TestE2EPlayZweitesSignal` (drei Unterfälle). Review-Mutant **M1** (`if abbruch.Err() != nil && !s.optionen.FinishSessionOnInterrupt {`): **rot** in Unit („Meldungen [… PGR-E4003 …: geschlossen], Exit-Code 4, erwartet [] mit 0“) und E2E (`/mit_Option` Exit 4 mit Zeile `PGR-E4003` statt 0, `/mit_Option_nach_früherem_Fehler` `PGR-E4003` vor `PGR-E4004`). **M19** (Schließen bei `abbruch` nur ohne Option): **rot** in Unit („Play endet binnen 30 s nach dem zweiten Signal nicht“) und E2E („play endet binnen 10 s nach dem zweiten Signal nicht“, beide Unterfälle mit Option). Am Binary hält das Verhalten (Punkt 3) | behoben |
| F-565 Hilfe | `play --help` am Binary: `--continue-on-error` sagt nur noch die Fortsetzung mit der nächsten Anfrage zu, keinen Exit-Code; `--allow-recorded-errors` nur eine Fehlerantwort mit dem Schweregrad `ERROR`, „wenn die aufgezeichnete Anfrage ebenfalls eine trägt“; `--finish-session-on-interrupt` das Ende nach der laufenden Session beim ersten `SIGINT` oder `SIGTERM`; dazu die drei Umgebungsvariablen. Jede Aussage hält am Binary (Punkt 1 bis 3), auch bei anderem SQLSTATE und bei `FATAL` | behoben |
| F-566 Handbuch §4 | Die Hinweise zu Fortsetzung (Exit 4, nach abbrechendem `PGR-E6001` bzw. `PGR-E4003` 6 bzw. 4), erwartetem Fehler (`ERROR`, `FATAL` bricht mit `PGR-E4003` ab) und Signalen (zweites Signal sofort, auch mit der Option, kein Fehler, Exit 0 bzw. 4) halten jede am Binary. §5 (Optionstabelle mit Umgebungsvariablen und Standard `false`, *Exit-Codes* Zeile 0) hält ebenso. Nicht gehalten ist der Absatz zu `--compare-responses` in derselben Liste, der zwei der drei Optionen nennt: V-147 | Hinweise behoben; V-147 |
| F-567 Kopf und §1 | Am Stand `785fca9` führt der Kopf `ARC-003`, §1 nennt den Kommentar des Ports und warum er keine dritte Schicht ist; `make kopf-check` im eigenen `make gates` grün. Der Code-Commit der Nacharbeit `eda0bc9` zieht §1 und §3 mit (12 Zeilen im Plan) | erledigt |
| F-568 Bindung von Anmeldung und TLS in (e) | eingeordnet: Der Hinweis betrifft das akzeptierte Negativ (e) in §6 und die Nehmer `slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls`, keinen DoD-Punkt dieses Slice. Beide nennen die Bindung weiter nicht. Am Binary trug jeder Fehler des Aufbaus einen Code (`PGR-E4002` bei fehlender Datenbank, auch im angehaltenen Aufbau kein Fehler ohne Code). Bleibt beim Architect und Planner | eingeordnet, offen beim Adressaten |
| F-569 Kommentar im Bootstrap | eingeordnet: `internal/bootstrap` liegt außerhalb des Slice (§1); der Kommentar ist unverändert. Das Verhalten, das er beschreiben soll (Exit-Code der ersten Meldung), hält am Binary | eingeordnet, offen beim Architect |

---

## 3. Hexagon und Architektur-Gate

- `make a-check` im eigenen `make gates`: siehe Punkt 4.
- `git diff --stat 85caec4..785fca9 -- internal/bootstrap internal/adapters/driven internal/adapters/driving/pgwire internal/hexagon/model cmd spec go.mod go.sum docs/plan/adr` ist leer. Außerhalb von `docs/` und `README.md` ändert der Diff nur `internal/adapters/driving/cli/`, `internal/hexagon/services/`, den Kommentar in `internal/hexagon/ports/driving/play.go` und `test/integration/play_laufsteuerung_e2e_test.go`.
- `services/play.go` importiert weiter nur `context`, `errors`, `fmt`, `model` und `ports/driven`. Signatur von `Player.Play` unverändert.
- [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md): Sessions nacheinander, Antworten außer `ErrorResponse` verworfen, Verbindungsfehler brechen immer ab, `--allow-recorded-errors` nur bei aufgezeichnetem Fehler — am Binary gesehen. Dass ein abbrechender `PGR-E6001` mit `--continue-on-error` Exit 6 statt „am Ende 4“ ergibt, schärft `LH-FA-20.a` *Exit-Code* auf der Spezifikationsebene (`AGENTS.md` §3.8); kein Widerspruch.

---

## 4. Plan gegen Code

- **§1:** Geliefert sind die drei Optionen am allgemeinen Leser mit Umgebung und `play:`, Fortsetzung, erwarteter Fehler, Session-Ende nach dem ersten Signal und die Exit-Codes nach `LH-FA-20.a`. Die Abgrenzungen halten: Extended-Interaktionen spielt `play` weiter nicht ein, `--compare-responses`, `--upstream-tls` und `--keep-timing` sind am Binary `PGR-E2001`, gegen den Server ohne Passwort und ohne TLS geprüft. Die Zusage „Die Aussagen zu den drei Optionen … im Benutzerhandbuch §4 bringt der Slice auf den gelieferten Stand“ hält für die Hinweise, nicht für den Vergleichs-Absatz (V-147).
- **§3:** Jede geänderte Datei steht dort. `internal/hexagon/services/play_test.go` (+3 −3) erweitert nur den Fake um `vorher` für die neuen Tests; das deckt die Zeile *Play-Service (Tests)*. `datei_test.go` entfernt den Fall `continue_on_error` ungültig, wie §3 sagt.
- **§4:** Keine Rückführung trat ein; kein Code im Bootstrap, der Diff war in einer Review-Sitzung prüfbar (Review *Größe*).
- **§6:** Jede Randform, die ich am Binary gefahren habe, folgt ihrem Ort in `LH-FA-20.a` bzw. `LH-FA-17.a`: *Optionen der Laufsteuerung*, *Erwarteter Fehler*, *Fortsetzung*, *Je `PGR-E4004` eine Log-Zeile*, *Erstes Signal im Aufbau*, *Exit-Code beim Abbruchsignal nach einem früheren Fehler*, *Exit-Code beim Abbruch nach einem früheren `PGR-E4004`*, *Reihenfolge der Zeilen `error`*, *Abbrechender Fehler nach dem ersten Signal*, *Fehlerregeln in der zu Ende laufenden Session; Signal zwischen Sessions*, *Mehrere Fehlerantworten* (zwei in einer Session), *Kombinationen der Optionen*, *Quellen und ungültige Werte*. *Fehlerantwort vor dem zweiten Signal* in derselben Interaktion habe ich am Binary nicht gesetzt; sie hält `TestPlayZweitesSignalNachFehlerantwort` (V14 rot). Keine Randform außerhalb von §6 ist im Code entschieden; §6 änderte nur der Architect (`85caec4`, `ee23fad`).
- **§7:** Die Belege nennen Größe, Schichten, Weg der Mutanten, beide Tabellen und die Läufe mit Stand, auch `make gates` am Stand der Nacharbeit; was die DoD verlangt, ist belegt. Der Absatz *Handbuch und Hilfe* des ersten Blocks galt am Stand `1b63ef8` nicht; die Nacharbeit sagt das selbst und nennt den Stand, ab dem er gilt.

---

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-147 | LOW | Der Hinweis zu `--compare-responses` im Handbuch §4 sagt zwei der drei Optionen zu: „mit `--continue-on-error` läuft es weiter und endet am Ende mit Exit-Code 5“ und „`--allow-recorded-errors` hat mit Vergleich keine Wirkung; die Kombination ist kein Fehler“. Am Binary ist `play --compare-responses --allow-recorded-errors` `PGR-E2001` („flag provided but not defined: -compare-responses“), Exit 2. Plan §1 sagt, der Slice bringe die Aussagen zu den drei Optionen im Handbuch §4 auf den gelieferten Stand; §7 *Nacharbeit* (F-566) nennt die geprüften Hinweise und schweigt zu diesem. Die Wirkung mit Vergleich gibt §1 an `slice-v1-abschluss-antwortvergleich` (dort DoD-Punkt 3); dessen §1, DoD und §3 nennen das Handbuch nicht, die Adresse nimmt die Handbuch-Stelle also nicht an (`AGENTS.md` §3.13). Wer den Hinweis liest, hält die Kombination für lieferbar. Dieselbe Lage haben im selben Abschnitt die Hinweise zu `--keep-timing`, `--upstream-tls`, `--upstream-ca` und das Passwort in *Vorgehen*; sie nennen keine der drei Optionen und liegen außerhalb dieses Slice. | `docs/user/benutzerhandbuch.md` §4 *Eine Aufzeichnung in eine Datenbank einspielen*, *Hinweise* · „mit `--continue-on-error` läuft es weiter und endet am Ende mit Exit-Code 5“ | Planner: die Sendung „Handbuch §4, Hinweis zu `--compare-responses` mit den beiden Optionen“ mit der Kennung dieses Slice in §1 oder §3 von `slice-v1-abschluss-antwortvergleich` eintragen und §1 dieses Slice bei der Closure auf „ohne Vergleich“ fassen, oder Implementer: den Hinweis als noch nicht geliefert kennzeichnen |

Kein HIGH, kein MEDIUM. DoD-Liefer-Punkte 1 bis 4 sind bestätigt; V-147 betrifft eine Zusage aus §1, keinen DoD-Punkt, und blockiert die Closure nicht.

**Übergabe an den Planner:** dieser Bericht mit V-147; dazu für die Closure die Ausgänge der drei Risiken aus Abschnitt 1 (*Übrige DoD-Punkte*) und die Klassen des Reviews.
