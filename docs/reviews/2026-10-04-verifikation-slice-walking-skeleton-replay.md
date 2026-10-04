# Verifikation: slice-walking-skeleton-replay — 2026-10-04

**Rolle:** Verifier (Modul 11). Geprüft wird, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-walking-skeleton-replay.md` <!-- d-check:ignore (Pfad beim geprüften Stand) --> (§1, §3, §5, §6) gegen `git diff c30b3f1 HEAD` bei HEAD `6636818`. Das sind 5 Commits, 27 Dateien, +1061/−61.

**Eingang:** DoD-Liefer-Punkte 1 bis 3 mit dem Hinweis des Implementers zur zweiten Runner-Phase (`tools/test/run-integration-tests.sh`, `TestE2EOhnePostgresReplay`, `test/integration/testdata/select1.yaml`). Dazu der Review-Report `2026-10-04-review-slice-walking-skeleton-replay.md` (F-273 bis F-282, Stand `c66cd55`) und der Korrektur-Commit `6636818`.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

Alle Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg. Mutationen und Sonden liefen nur in einer Kopie von `git archive HEAD` im Scratchpad, nie im Repo. In der Kopie hießen Runner, Image-Tags und Netze anders (`pgr-verif:{integration,test}`, `pgr-vf-*`). Die beiden Images habe ich danach gelöscht, Container und Netze ebenso. Die Images `pgwire-recorder:*` erzeugen die Make-Ziele des Repos; sie lagen schon vorher vor, ihre IDs sind unverändert. `git status --short` war vor dem Anlegen dieser Datei leer.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — LH-FA-03, LH-FA-09: Mit gestoppter PostgreSQL liefert Replay dem Client für `SELECT 1;` dasselbe Ergebnis wie im Record (End-to-End-Test). **Bestätigt**, mit einer Grenze des Test-Orakels (V-9).

**Lauf.** Eigener Lauf von `make gates`, Teil `test-integration`, gegen das gepinnte `postgres:17-alpine@sha256:b0f9560a…` in einem internen Docker-Netz:

| Phase | Test | Ergebnis |
|---|---|---|
| 1 (PostgreSQL läuft) | `TestE2EReplaySelect1`, `TestE2EReplayAbweichung`, `TestE2EReplayBeschaedigt` und zehn `TestE2ERecord*` | PASS |
| 1 | `TestE2EOhnePostgresReplay` | SKIP (gewollt, nur Phase 2) |
| 2 (`docker stop` auf den PostgreSQL-Container) | `TestE2EOhnePostgresReplay` | PASS |
| Runner | | `run-integration-tests: gruen` |

**Tragen die beiden Tests den Punkt gemeinsam?** Ja, aber nur zusammen mit der Fixture-Gleichheit, die ich unten eigens belegt habe.

- `TestE2EReplaySelect1` belegt „wie im Record“. Die Aufzeichnung entsteht im selben Lauf über `record` gegen die reale Instanz und wird mit einem frischen Prozess abgespielt (das deckt auch den Rest aus V-3 des Record-Verifikationsbelegs). Der Upstream läuft dabei noch, `replay` hat aber keine Upstream-Option.
- `TestE2EOhnePostgresReplay` belegt „mit gestoppter PostgreSQL“. Der Test prüft die Vorbedingung selbst: Ist `PGR_UPSTREAM` erreichbar, ist er rot. Abgespielt wird eine Fixture, nicht die Aufzeichnung aus Phase 1.
- Die Brücke ist, dass die Fixture dem entspricht, was `record` schreibt. Das sichert kein Test, also habe ich es selbst geprüft.

**Eigene Sonde zur Fixture** (Probe-Test in der Kopie, Phase 1). Sie zeichnet `SELECT 1;` über `record` auf und hält drei Client-Sichten nebeneinander: das Ergebnis im Record, im Replay der eben geschriebenen Aufzeichnung und im Replay von `select1.yaml`. Gemessen werden Spalte mit OID, Zeilen und Befehlsabschluss:

```
RECORD / REPLAY-REC / REPLAY-FIX (alle drei gleich):
spalte ?column? 23
zeile ["1"]
befehl SELECT 1
```

Die Interaktion von Session 1 in der Fixture ist inhaltlich gleich mit der, die `record` schreibt. Der einzige Unterschied ist die Quotierung von `'?column?'`. **Byte-gleich mit einer Aufzeichnung von `record` ist die Fixture nicht** (V-11):

- `server_parameters` führt 6 der 14 Parameter, die `record` gegen dieses Image schreibt. Es fehlen `IntervalStyle`, `TimeZone`, `application_name`, `default_transaction_read_only`, `in_hot_standby`, `is_superuser`, `scram_iterations` und `session_authorization`.
- Session 2 (`SELECT 2;` → nur `command_complete` mit Tag `SELECT 1`, ohne `row_description`/`data_row`) ist von Hand gebaut. `record` schriebe so etwas nie. Sie ist nur dazu da, die Warnung `PGR-W2001` am Laufende auszulösen.

**Eigener Welle-Smoke** (der Closure-Trigger aus §5, Ablauf wörtlich: Record, PostgreSQL stoppen, Replay) mit dem Binary aus der Stufe `integration` und `psql` 17 aus dem gepinnten Image:

1. `record` gegen PostgreSQL, `psql -c 'SELECT 1;'`, SIGTERM → Exit 0
2. PostgreSQL-Container gestoppt und entfernt
3. zweimal `replay --input` mit **derselben** Datei, je `psql -c 'SELECT 1;'`, SIGTERM → Exit 0, 0

Ergebnis: Die Ausgabe von `psql` (`?column? / 1 / (1 row)`) ist in Record, Replay 1 und Replay 2 byte-gleich (`cmp`).

**Mutationsprobe** (in der Kopie, Runner mit Phase 1 auf `TestE2EReplaySelect1` eingeschränkt, Phase 2 unverändert):

| # | Mutation | Ergebnis |
|---|---|---|
| MA | Runner ohne `docker stop` vor Phase 2 | **rot:** `PostgreSQL unter …-pg:5432 ist noch erreichbar` |
| MC | Replay liefert jede Spalte als `X` mit Typ-OID 25 statt `?column?`/23 (Service, durchgängig gleich) | **grün** in Phase 1, Phase 2 und `make test`-Stufe |

MA zeigt, dass die Vorbedingung „gestoppt“ tatsächlich hergestellt und geprüft wird. MC zeigt die Grenze des Orakels: Beide E2E-Tests prüfen nur `zeile ["1"]` und `befehl SELECT 1`, nicht Spaltenname und Typ, und vergleichen nicht mit dem Ergebnis, das der Client im Record sah. Der Punkt ist durch Sonde und Smoke bestätigt. Der Test selbst trägt „dasselbe Ergebnis“ nur für Wert und Befehlsabschluss (V-9).

### Punkt 2 — LH-QA-01: Zwei aufeinanderfolgende Replay-Läufe mit demselben Recording zeigen identisches Verhalten. **Bestätigt.**

`TestE2EReplaySelect1` startet zehn Replay-Prozesse nacheinander mit derselben Aufzeichnung und denselben zwei Anfragen (`SELECT 1 AS eins;`, `SELECT 'a' AS text, NULL AS leer;`). Verglichen wird die beobachtete Sicht jedes Laufs (Spalten mit OID, Zeilen inklusive NULL, Befehlsabschluss) mit der des ersten Laufs. Zehn Läufe sind mehr als die zwei der DoD und erfüllen die Messmethode von LH-QA-01 („mindestens zehn“). Im eigenen Lauf von `make gates`: PASS.

| # | Mutation | Ergebnis |
|---|---|---|
| MB | Replay benennt die Spalten nur in jedem zweiten Prozess um (abhängig von der Startzeit) | **rot:** `Lauf 2 weicht ab` |

Der Welle-Smoke oben hat außerdem zwei Replay-Läufe mit `psql` byte-gleich gezeigt. Grenze: Der Vergleich deckt keine `ParameterStatus` des Handshakes ab und keine Fehlerpfade. Das genügt für die DoD-Zeile.

### Punkt 3 — `make gates` grün. **Bestätigt.**

Eigener Lauf von `make gates` am Stand `6636818`, Exit 0. Im Log:

- `build` (Stufe `runtime`) und `test` (Stufe `test`: `go vet -tags integration ./... && go test ./...`)
- `abdeckung-check` (gesamt: 0 Befunde), `abdeckung-gegenprobe` grün
- `a-check` 0 Befunde, `a-check-negativ` grün
- `docs-check` mit 127 Dateien und 0 Befunden, `commit-msg-gegenprobe` grün
- `baseline-verify` v6.13.0 OK (54 Dateien)
- `test-integration` grün in beiden Phasen
- zuletzt `record-gates`

### Prozess-Punkte der DoD (nicht Teil des Auftrags, nur vermerkt)

- Review: Ein Report liegt vor, er erfasst aber den Stand `c66cd55`. Den Korrektur-Commit `6636818` hat kein Review gesehen (V-14).
- Closure-Notiz, Beobachtungs-Register, Ausgänge der Risiken aus §6 und die drei Paarungen: noch offen, wie bei `in-progress` zu erwarten.

---

## 2. Stand der Review-Findings (F-273 bis F-282)

Geprüft am Code von `6636818`. Die Mutationen des Reviews habe ich in der Kopie wiederholt, soweit sie im Review grün geblieben waren.

| ID | Stand | Beleg |
|---|---|---|
| F-273 | **getragen** | (1) Handshake der nächsten freien Session: `TestReplayHandshake`. Meine Wiederholung von R3 (immer die letzte Session) ist **rot**. (2) Diagnose nach LH-FA-10.a: `TestReplayMismatch` prüft Session, Interaktion, erwartete und empfangene Anfrage; R12 **rot**. (3) `PGR-W2001` für nie zugeordnete Sessions: `TestE2EOhnePostgresReplay` über Session 2 der Fixture; R11 (Bootstrap loggt `Unassigned` nicht) in Phase 2 **rot**. Sortierung: R15 mit fünf Parametern **rot** (`TestReplayHandshake`, `TestReplayStrictSequential`); nicht mehr nur zufällig. |
| F-274 | **getragen**, mit Rest | §1 grenzt jetzt Extended-Mismatch, `--session-assignment connection` und `--fail-on-unconsumed` aus und nennt das Gelieferte. `Bezug` nennt LH-FA-10, -12, -13, `Berührte Spec-Stellen` nennt `SPEC-018`, `SPEC-027`, `SPEC-034` und LH-FA-03.a/.b, -09.a, -10.a, -12.a. `welle-walking-skeleton.md` §Abgrenzung sowie die beiden Folge-Slices sind nachgezogen. Rest: §3 (V-12) und die DoD des Folge-Slice (V-13). |
| F-275 | **getragen** | `TestFelderDerVersion1`: R14a (negatives `offset_ms` angenommen), R14c (`OffsetMS` nicht ins Modell), R14d (`EmptySessions` nicht ins Modell) sind **rot**. R14b (Kennzeichnung ignoriert) fängt der Fall „leere Session mit Kennzeichnung“ im selben Test; ich habe ihn gelesen, nicht mutiert. |
| F-276 | **getragen** | Felder von `Server` sind nicht exportiert; `NewRecordServer` und `NewReplayServer` setzen genau einen Use Case. Ein Nullwert entsteht nur noch im Paket selbst (`TestErsterFehlerZaehlt`). |
| F-277 | **getragen** | Eigener Typ `model.Warning` ohne `ExitCode`/`Error`; der Port gibt bei `CloseConnection` `*model.Warning` zurück, `Unassigned` ebenso. Entspricht `SPEC-034` („trägt keine Klasse“). |
| F-278 | **getragen über den Kommentar** | Der Typ `SessionID` bleibt für beide Räume; der Kommentar nennt jetzt beide Bedeutungen. Für ein LOW vertretbar. |
| F-279 | **offen** (Validator) | Plan §6, Risiko 1 steht weiter auf „offen bis Closure“. Mein Smoke fügt `psql` 17 ohne `sslmode`-Angabe hinzu, keinen weiteren Treiber. |
| F-280 | **getragen** | Zweite Runner-Phase, Vorbedingung im Test geprüft (MA **rot**). `LH-FA-03/Happy` ist von `TestE2EReplaySelect1` auf `TestE2EOhnePostgresReplay` gewandert. Siehe Punkt 1 und V-9, V-11. |
| F-281 | **getragen** | `Query` prüft `erwartet.Request.Type != model.RequestQuery`. Der Ausgang von Risiko 2 aus §6 ist bei der Closure zu setzen. |
| F-282 | **teilweise** | `slice-replay-semantik-mismatch` §1 nennt jetzt `--fail-on-unconsumed` (`PGR-E5002`), DoD und §3 dort aber nicht (V-13). |

---

## 3. Plan gegen Code-Diff

### §1 Ziel und Abgrenzung

| Plan | Code | Urteil |
|---|---|---|
| `replay` beantwortet `SELECT 1;` aus dem Recording, ohne PostgreSQL, strict sequential | Port `Replayer`, `ReplayService`, Replay-Modus des PGWire-Adapters, Kommando `replay` | konform |
| Extended-Mismatch nicht hier; Mismatch einfacher Anfragen mit `PGR-E5001` und Exit 5 schon hier | `Query` → `PGR-E5001`; `FirstErrorCode` → 5 (`TestE2EReplayAbweichung`) | konform |
| Fehlerreplay nicht hier | kein Pfad für `error_response` im Replay über das hinaus, was gespeicherte Antworten ohnehin liefern | konform |
| `--session-assignment connection`, `--fail-on-unconsumed` nicht hier; `first-request`, `PGR-W2001` als Warnung | `parseReplay` kennt nur `--listen`, `--input`; `CloseConnection`/`Unassigned` → `model.Warning` | konform |

### §3 Dateien

| Plan-Zeile | Im Diff | Urteil |
|---|---|---|
| `ports/driving` (`Replayer`), `services` (Replay-Service) | `replay.go` in beiden | konform |
| `adapters/driving/pgwire` | `server.go` (Replay-Modus, Konstruktoren) | konform |
| `adapters/driving/cli` | `cli.go`, `cli_test.go` | konform |
| `adapters/driven/recording` (`offset_ms`, `empty_sessions`) | `yaml.go`, `yaml_test.go` | konform |
| `bootstrap` | `bootstrap.go` | konform |
| Unit- und E2E-Tests | `replay_test.go`, `server_test.go`, `replay_e2e_test.go` | konform |
| — | `internal/hexagon/model/fehler.go` (`Warning`, Codes), `model/recording.go` (`OffsetMS`, `EmptySessions`, `SessionID`-Kommentar) | **nicht im Plan** (V-12) |
| — | `tools/test/run-integration-tests.sh` (Phase 2), `Dockerfile` (`PGR_FIXTURES`), `test/integration/testdata/select1.yaml`, `TestMain` in `record_e2e_test.go` | **nicht im Plan** (V-12) |
| — | `docs/user/abdeckung-*.md`, `README.md` | nicht im Plan; abgeleitete bzw. Überblicks-Dokumente, unkritisch |

### §5 Closure-Trigger

„Der Welle-Smoke (Record, PostgreSQL stoppen, Replay) ist durchlaufen“: Kein automatischer Lauf verkettet Record, Stopp und Replay **derselben** Datei, Phase 2 spielt die Fixture ab. Der Smoke in §1 dieses Belegs ist durchlaufen und kann bei der Closure zitiert werden (V-15).

### §6 Risiken

| Risiko | Stand im Code | Ausgang bei Closure |
|---|---|---|
| 1 — Authentifizierung typischer Treiber | Handshake `AuthenticationOk` + `ParameterStatus` + `ReadyForQuery`, ohne `BackendKeyData`; belegt mit `pgconn` und `psql` | offen; Urteil beim Validator (F-279) |
| 2 — `offset_ms`, `empty_sessions`, `type: extended` | die ersten beiden gelesen und getestet; `extended` gezielt als `PGR-E3003` abgelehnt, Matcher prüft den Typ zusätzlich | trägt den Ausgang „entfallen“ bzw. „eingetreten und gelöst“; vom Implementer bei der Closure zu setzen |

---

## 4. ADR-Konformität

| ADR | Zusage | Urteil |
|---|---|---|
| [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) | Schichten des Hexagons, Import-Richtungen | **konform.** `make a-check` 0 Befunde, `make a-check-negativ` grün. `model` importiert nur `fmt`; der Service nur `context`, `sort`, `sync`, `model`, `ports/driven`. |
| [ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md) | PGWire-Server ist Driving Adapter ohne Fachlogik | **konform.** Cursor, Zuordnung und Vergleich liegen im Service; der Adapter übersetzt, schreibt `ErrorResponse` und loggt die Warnung. |
| [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md) | Kanonisches Domain-Modell an den Ports | **konform.** `Replayer` nutzt `model.SessionID`, `model.Response`, `model.Warning`; `OffsetMS`/`EmptySessions` sind Modellfelder ohne Bibliothekstyp. |
| [ADR-0007](../plan/adr/0007-strict-replay.md) | Strict Replay | **konform.** Vergleich byte-genau (`!=`), Cursor nur bei Gleichheit, keine Suche, nach dem Ende `PGR-E5001`. |
| [ADR-0008](../plan/adr/0008-kein-sql-parser-in-v1.md) | Kein SQL-Parser | **konform.** Kein Parser, keine Normalisierung. |
| [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) | Build und Test im Multistage-Dockerfile; Integration gegen gepinntes PostgreSQL im eigenen Netz, an `make gates` | **konform.** Phase 2 läuft im selben Runner, Netz `--internal`, Image per Digest; `PGR_FIXTURES` im Dockerfile, keine Netzstufe hinzugekommen. |
| [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) | Deklaration je Anforderung und Pfad; RTM nur bei vollständiger Belegung | **mechanisch konform** (`abdeckung-check` grün, Tabellen aktuell). Inhaltlich siehe §6: Grenze bei LH-QA-01 (V-10), Orakel bei LH-FA-03/Happy (V-9). |

---

## 5. Spec-Konformität

| Stelle | Zusage | Urteil |
|---|---|---|
| LH-FA-03.a | Pflicht `--listen`, `--input`; kein Upstream; Laden vor dem Start, Fehler → 3; ohne zuordenbare Session `PGR-E3004`; Cursor je Verbindung; Anfrage ohne freie Session `PGR-E5003` | **konform.** `parseReplay`; `NewReplayService` vor `Listen`; `TestReplayStartfehler`, `TestE2EReplayBeschaedigt`; Cursor in `verbindung`. Die optionalen Optionen fehlen, das grenzt §1 aus. |
| LH-FA-03.b | Rest einer Session oder nie zugeordnete Session → mindestens `PGR-W2001` | **konform.** `CloseConnection` und `Unassigned`; der Teil „nie zugeordnet“ ist in Phase 2 belegt (R11 rot). |
| LH-FA-09.a | Exakter Stringvergleich, Position zählt, jede Interaktion höchstens einmal | **konform** (`TestReplayStrictSequential`, `TestReplayMismatch`). |
| LH-FA-10.a | Diagnose mit Kontext, Nummer, erwarteter und empfangener Query; kein Sprung; Anfrage nach dem Ende ist Mismatch; Exit 5 | **konform.** Text `Session %d, Interaktion %d: erwartet %q, empfangen %q`, getestet (R12 rot); nach dem Ende eigener Text ohne „erwartete Query“, denn es gibt keine. |
| LH-FA-12.a (Replay) | `first-request`: n-te Verbindung mit Anfrage → n-te Session mit Interaktion; leere übersprungen; Handshake aus der nächsten freien, sonst der letzten Session | **konform.** `frei` nur mit Interaktionen; Zuordnung bei der ersten Anfrage; `TestReplayHandshake`. Die Spec sagt „Startup-Daten“, der Code nimmt die aufgezeichneten `server_parameters`. Das ist der Teil, den der Server im Handshake sendet. |
| LH-FA-13.b | Startfehler vor dem ersten Accept; Verbindungsfehler mit `ErrorResponse` und Code im Text, erster gemerkt; Prozessende: nie zugeordnete Sessions als Warnung, Exit der gemerkten Klasse, sonst 0 | **konform.** Exit 0 trotz `PGR-W2001` in Phase 2; Exit 5 in `TestE2EReplayAbweichung`. |
| `SPEC-011` | strict sequential einziges Matching | **konform.** |
| `SPEC-018` | Exit 5 für `PGR-E5001`, `PGR-E5003` | **konform.** `ExitCode` aus der ersten Ziffer; E2E für `PGR-E5001`. |
| `SPEC-027` | Replay-Mismatch → 5 mit `PGR-E5001`/`PGR-E5003` | **konform.** |
| `SPEC-034` | Codes nach Tabelle; Fehlertext `<klasse> [<code>]: …`; Warnung als Attribut `code`, Text ohne Code, Exit unverändert | **konform.** Codes in `fehler.go` stimmen mit der Tabelle überein; `model.Warning` ohne Klasse; Log `log.Warn(w.Msg, "code", w.Code)`. |

---

## 6. Ehrlichkeit der Abdeckung ([ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md))

Gelesen gegen die Akzeptanzkriterien im Lastenheft. `make doc-trace` meldet: LH-FA-09, LH-FA-10 und LH-QA-01 `Tests | ok`, LH-FA-03 `WAISE`, LH-FA-12 ohne Tests.

| Anforderung / Pfad | Deklariert über | Urteil |
|---|---|---|
| LH-FA-03/Happy | `TestE2EOhnePostgresReplay` | **trifft die Vorbedingung** (kein erreichbarer Server, MA rot). Das Orakel „aufgezeichnetes, aus Sicht des Clients relevantes Verhalten“ prüft nur Wert und Befehlsabschluss (V-9). |
| LH-FA-03/Boundary | — | nicht deklariert; ehrlich, denn kein Test beobachtet, dass beim Start kein Verbindungsversuch erfolgt. LH-FA-03 steht als „teilweise“. |
| LH-FA-03/Negative | `TestE2EReplayBeschaedigt`, `TestReplayStartfehler` | trifft. |
| LH-FA-07/Happy | `TestE2EReplaySelect1` | trifft (Record-Prozess beendet, neuer Prozess lädt und spielt ab). |
| LH-FA-09 Happy/Boundary/Negative | `TestE2EReplaySelect1` / `TestReplayStrictSequential` / `TestReplayMismatch` | trifft; „vollständig“ ist ehrlich (MB rot). |
| LH-FA-10 Happy/Boundary/Negative | `TestE2EReplayAbweichung` / `TestReplayMismatch` (Whitespace) / `TestReplayMismatch` (Cursor bleibt, keine Antwort) | trifft; „vollständig“ ist ehrlich. |
| LH-FA-13/Negative | `TestE2EReplayAbweichung` (Klasse 5) | trifft. |
| LH-FA-12 | — | nicht deklariert, obwohl `TestReplaySessionZuordnung` und `TestReplayHandshake` Teile des Replay-Pfads von Boundary belegen; zurückhaltend, nicht unehrlich. |
| LH-QA-01/Messung | `TestE2EReplaySelect1` | trifft die **Messmethode** (zehn Läufe). LH-QA-01 trägt aber eine zweite Zusage („Bereitstellung“: zwei Builds desselben Quellstands mit derselben Prüfsumme), die nichts misst. Trotzdem steht LH-QA-01 als „vollständig“ und in der RTM als `ok` (V-10). |

---

## 7. Befunde des Verifiers

| ID | Kategorie | Befund | Beleg |
|---|---|---|---|
| V-9 | MEDIUM | Das Orakel der E2E-Tests ist schwächer als die Zusage „dasselbe Ergebnis wie im Record“ (DoD 1) und „aus Sicht des Clients relevantes Verhalten“ (LH-FA-03/Happy). `TestE2EReplaySelect1` und `TestE2EOhnePostgresReplay` prüfen nur `zeile ["1"]` und `befehl SELECT 1`. Keiner vergleicht die Sicht des Clients im Replay mit der im Record, obwohl `beobachten` Spalten und OIDs schon erfasst. Ein Replay, das Spaltennamen und Typ verfälscht, bleibt grün. Vorschlag: In `TestE2EReplaySelect1` die Sicht während `aufnehmen` mit `beobachten` festhalten und Lauf 1 dagegen vergleichen; in Phase 2 die volle erwartete Sicht (`spalte ?column? 23 …`) prüfen. | Mutation MC: Phase 1, Phase 2 und Unit-Stufe grün |
| V-10 | LOW | LH-QA-01 steht in `abdeckung-gesamt.md`/`abdeckung-vollstaendig.md` als vollständig und in der RTM als `Tests \| ok`. Belegt ist nur der Replay-Determinismus (Messmethode). Die Bereitstellungs-Zusage (reproduzierbares Binary, gleiche Prüfsumme) misst kein Sensor. Das Pfadmodell von [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) kennt für QA nur „Messung“ und kann den zweiten Teil nicht ausdrücken. Frage an den Planner: Deklaration zurücknehmen, bis ein Prüfsummen-Vergleich existiert, oder die Grenze in der Abdeckungs-Sicht nennen. | `make doc-trace`; `grep` nach Prüfsummen-Vergleich in `tools/`, `harness/mk/`: kein Treffer |
| V-11 | LOW | Die Fixture `select1.yaml` ist eine von Hand reduzierte und erweiterte Aufzeichnung, keine Ausgabe von `record`: 6 von 14 Serverparametern, und Session 2 ist so gebaut, wie `record` nie schreiben würde. Die Interaktion `SELECT 1;` ist heute gleich mit der echten (eigene Sonde). Nichts hält diese Gleichheit fest, wenn sich `record` oder das PostgreSQL-Image ändert. Phase 2 belegt „wie im Record“ also nur über diese ungeprüfte Brücke. Vorschlag: Phase 1 schreibt die Aufzeichnung in ein Volume, Phase 2 spielt genau sie ab (dann ist der Welle-Smoke automatisiert). Session 2 bekäme dann eine eigene, ausgewiesene Fixture, oder die Fixture trägt einen Kopfkommentar, der sagt, dass sie von Hand gebaut ist. | Probe `TestVerifProbe` (Abschnitt 1) |
| V-12 | LOW | Plan §3 nennt die Dateien des Korrektur-Commits nicht. Es fehlen `internal/hexagon/model/fehler.go` und `recording.go`, `tools/test/run-integration-tests.sh` (Phase 2), `Dockerfile` (`PGR_FIXTURES`), `test/integration/testdata/select1.yaml` und `TestMain`. Die E2E-Zeile von §3 erwähnt die Phase ohne PostgreSQL nicht. Gleiche Klasse wie F-274 („Plan-Zeile nach Korrektur nicht nachgezogen“), jetzt an §3. | `git diff c30b3f1 HEAD --stat` gegen §3 |
| V-13 | INFO | `slice-replay-semantik-mismatch`: Die DoD-Punkte zu LH-FA-10 (Mismatch, Query nach dem Ende, keine fremde Antwort) und LH-FA-13 (`ErrorResponse`, Klasse 5) hat dieser Slice bereits geliefert und belegt. `--fail-on-unconsumed` (`PGR-E5002`) steht dort nur in §1, nicht in DoD und §3. F-282 ist damit halb getragen. Hinweis an den Planner. | `docs/plan/planning/open/slice-replay-semantik-mismatch.md` §1–§3 |
| V-14 | INFO | Den Commit `6636818` hat kein Review gesehen; der Report erfasst `c66cd55`. Er ändert die Port-Signatur (`CloseConnection` → `*model.Warning`), die Konstruktion des Servers und den Runner. Funktional habe ich ihn hier verifiziert (Abschnitt 2). Ob ein Folge-Review nötig ist, entscheidet der Planner vor der Closure. | `git log c66cd55..HEAD` |
| V-15 | INFO | Den Closure-Trigger „Welle-Smoke (Record, PostgreSQL stoppen, Replay)“ erfüllt kein automatischer Lauf wörtlich. Der Smoke dieses Belegs (Abschnitt 1, `psql`, derselbe Recording-Pfad, zwei Replays, `cmp` gleich) kann dafür zitiert werden. Mit dem Vorschlag aus V-11 würde der Runner ihn tragen. | Abschnitt 1, Welle-Smoke |

---

## 8. Gegenprobe dieser Datei

`make docs-check` nach dem Anlegen dieser Datei: `d-check: 128 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

## 9. Gesamturteil

**Die DoD-Liefer-Punkte 1 bis 3 sind bestätigt.** Grundlage sind eigene Läufe:

- `make gates` mit Exit 0, beide Runner-Phasen grün
- `make doc-trace`
- eine Sonde, die Record, Replay der echten Aufzeichnung und Replay der Fixture vergleicht
- ein Welle-Smoke mit `psql`: Record, PostgreSQL entfernt, zwei Replays byte-gleich
- 10 Mutanten: MA, MB, R3, R11, R12, R14a, R14c, R14d und R15 sind rot. MC bleibt grün, das ist V-9.

**Zu Punkt 1 im Einzelnen:** `TestE2EReplaySelect1` (Record → Replay, PostgreSQL läuft noch) und `TestE2EOhnePostgresReplay` (PostgreSQL gestoppt und geprüft, Fixture) tragen den Punkt nur zusammen mit der Gleichheit der Fixture. Diese hält kein Test fest. Ich habe sie für `SELECT 1;` selbst belegt: Die Client-Sicht ist gleich. Die Datei selbst ist nicht das, was `record` schreibt (V-11), und das Orakel beider Tests ist schmaler als die Zusage (V-9).

Die Review-Findings F-273, F-275 bis F-278, F-280 und F-281 trägt `6636818`. F-274 und F-282 sind nur zum Teil getragen (V-12, V-13), F-279 bleibt beim Validator.

Der Code hält [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md), [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md), [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0008](../plan/adr/0008-kein-sql-parser-in-v1.md) und [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) ein. Er entspricht LH-FA-03.a/.b, LH-FA-09.a, LH-FA-10.a, LH-FA-12.a (Replay-Teil), LH-FA-13.b und `SPEC-011`, `SPEC-018`, `SPEC-027`, `SPEC-034`. [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) ist mechanisch gehalten. Inhaltlich ehrlich sind LH-FA-09 und LH-FA-10. LH-QA-01 steht wegen der ungemessenen Bereitstellungs-Zusage als vollständig, obwohl es das nur zur Hälfte ist (V-10).

**Kein Befund blockiert die DoD.** Vor der Closure empfohlen:

1. V-9: das Orakel schärfen (Record-Sicht gegen Replay-Sicht, volle Sicht in Phase 2).
2. V-11 und V-15: Phase 2 die Aufzeichnung aus Phase 1 abspielen lassen, oder die Fixture ausdrücklich als von Hand gebaut ausweisen.
3. V-10: LH-QA-01 zurücknehmen oder seine Grenze nennen (Planner).
4. V-12 und V-13: Plan §3 und den Folge-Slice nachziehen (Planner).
5. V-14: über ein Folge-Review für `6636818` entscheiden (Planner).
