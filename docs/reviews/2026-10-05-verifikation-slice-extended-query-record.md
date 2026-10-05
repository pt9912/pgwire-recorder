# Verifikation: slice-extended-query-record — 2026-10-05

**Rolle:** Verifier (Modul 11). Geprüft wird, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-extended-query-record.md` (Kopf, §1 bis §4, §6) gegen die Commits `5d666db` (Umsetzung), `0dbcf9d` (Umbau nach [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md), F-301 bis F-310) und `09cc0e9` (Folge-Review F-311 bis F-319) bei HEAD `09cc0e9`. `09cc0e9` hat noch kein Review gesehen; er ist hier mit eigenen Mutationen und Sonden geprüft (Abschnitt 3).

**Eingang:**

- die DoD-Liefer-Punkte 1 bis 3 und die Commit-Messages
- der [Review-Report](2026-10-05-review-slice-extended-query-record.md) (F-301 bis F-310, Stand `5d666db`)
- das [Folge-Review](2026-10-05-folge-review-slice-extended-query-record.md) (F-311 bis F-320, Stand `0dbcf9d`)
- die [Mutationstabelle](2026-10-05-mutationen-slice-extended-query-record.md) des Implementers (drei Runden)

Der Arbeitsbaum war bei meinem Start sauber.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

Die Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg. Mutationen und Sonden liefen nur in Kopien von `git archive HEAD` im Scratchpad:

- Unit-Mutationen: je Mutation eine eigene Kopie, `go test -count=1` des betroffenen Pakets, netzlos in einem Container des Images `pgr-verif:source` (Stufe `source` des `Dockerfile`)
- E2E-Sonden: zusätzliche Testdatei in `test/integration`, Image `pgr-verif:<tag>` (Stufe `integration`), eigenes internes Docker-Netz `pgr-verif-*-net` und eigener PostgreSQL-Container aus dem gepinnten Image von `harness/mk/integration.mk`
- Race-Detector: Container aus `pgr-verif:source` mit `build-base` (Netz nur für `apk`), kein Gate

Images, Container, Netze und Kopien habe ich danach gelöscht. Die Images `pgwire-recorder:*` erzeugen die Make-Ziele des Repos; ihre IDs sind unverändert. `git status --short` war vor dem Anlegen dieser Datei leer.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Go-Client mit Prepared Statements über den Recorder, Aufzeichnung mit der Nachrichtenfolge (Integrationstest). **Bestätigt.**

- `TestE2ERecordExtendedPgx` grün im eigenen `make gates`. Der Test nutzt pgx im Standardmodus gegen `record` und das gepinnte PostgreSQL 17. Er prüft die Folge mit `inReihe` und startet `replay` mit der Aufzeichnung (`laedt`: Leser und `Validate` bestanden).
- Eigene Sonde **V-S1**, nicht aus der Mutationstabelle: pgx-Batch mit 400 Anfragen mit je 64 KiB Binärparameter (`bytea`) und einer Anfrage mit 24 MB Ausgabe, alles in einer Sync-Gruppe.
  - Ergebnis: grün in 0,17 s, alle Werte korrekt, Exit-Code 0.
  - Die Aufzeichnung trägt 401 `bind` und zwei `sync`: die Prepare-Runde von pgx und den Batch. Sie lädt (`laedt`).
- Großer Batch ohne Selbstblockade: `TestE2ERecordExtendedGegendruck` (32 MB Ausgabe, 16 MB Parameter, Batch und Flush-Pipeline) ist grün.

### Punkt 2 — Gruppierung nach `LH-FA-18.a` (Flush, Sync, Pipelining, Abbruch), jede aufgezeichnete Interaktion besteht `Interaction.Validate` (Unit-Tests). **Bestätigt.**

Eigene Läufe der Mutationen in Kopien von HEAD, jede rot:

| # | Mutation | Rot in |
|---|---|---|
| M1 | jede Client-Nachricht sofort senden | `TestRecordExtendedSyncGruppe`, `…FlushGruppen`, `…Pipelining`, `TestRecordHerunterfahren`, `TestRecordGegendruck` |
| M3 | Ergebnis von `Validate` ignoriert | `TestRecordExtendedNichtDarstellbar` |
| N12 | Server-Nachrichten der jüngsten statt der ältesten Interaktion zugeordnet | `TestRecordExtendedPipelining` |

Den Abbruch ohne `Sync` (PGR-E4003, unvollständige Interaktion nicht übernommen) trägt `TestE2ERecordExtendedAbbruch`, grün.

### Punkt 3 — `make gates` grün. **Bestätigt.**

- Eigener Lauf von `make gates` an HEAD, Exit 0:
  - `abdeckung-check`/`-gegenprobe` grün
  - `a-check` 0 Befunde, `a-check-negativ` grün
  - `baseline-verify` v6.13.0 OK
  - `docs-check` 153 Dateien, 0 Befunde
  - `commit-msg-gegenprobe` grün
  - `test-integration` in beiden Phasen grün, alle `TestE2ERecord*`
- Die Stufe `test` kam im Gate-Lauf aus dem Cache. In der Kopie habe ich sie darum ohne Cache gefahren (`--no-cache-filter test`): `gofmt -l` leer, `go vet` und alle Unit-Pakete `ok`.
- Race-Detector `go test -race -count=20 ./internal/adapters/... ./internal/hexagon/...`: grün.

### Prozess-Punkte der DoD (nur vermerkt)

- **Review:** Die Reports zu `5d666db` und `0dbcf9d` liegen vor. `09cc0e9` bringt neuen Nebenläufigkeits-Code:
  - `wecke`/`weiterlesen` unter der Sperre `frist`
  - `wartenAufEnde`
  - `beende` mit Fehlerantwort und dem Schließen der Verbindung

  Ein Review dazu gibt es nicht. Dieser Bericht prüft den Commit funktional. Ob ein zweites Folge-Review nötig ist, entscheidet der Planner.
- **Offen, wie bei `in-progress` zu erwarten:** Closure-Notiz, Register, Ausgänge der Risiken in §6 und die Paarungen.

---

## 2. Stand der Review-Findings (F-301 bis F-320)

| ID | Stand | Beleg (eigene Läufe an HEAD) |
|---|---|---|
| F-301 | **behoben** | `TestE2ERecordExtendedGegendruck` und `…SigtermNachBlockade` grün. Dazu V-S1 (Batch anderer Form) und V-S5 (Client liest nicht, Session endet trotzdem). |
| F-302 | **behoben** | Zuschnitt entschieden in [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) (`Accepted`). `ARC-004` nennt Senden und Empfangen getrennt. |
| F-303 | **behoben** | Am Stand `0dbcf9d` vom Folge-Review bestätigt. Die Tests stehen an HEAD unverändert, die Unit-Pakete sind grün. |
| F-304 | **behoben** | `richtungen` führt nur Transportzustand (`schreiben`, `beendet`, `einmal`, `ende`, `frist`/`weck`). Ob die Session wartet und ob eine Nachricht eine Interaktion beginnt, entscheidet der Service (`ErrShutdown`, `Shutdown`, `Delivered`). |
| F-305 | **behoben** | Entschieden in `LH-FA-13.a`. `TestExtendedHerunterfahren`, Fall „gelesene Anfrage“, prüft jetzt die Reihenfolge im Protokoll („q:SELECT 1 zugestellt shutdown“). Die Mutation H4 der Tabelle deckt das ab. |
| F-306, F-307 | **behoben** | am Stand `0dbcf9d` bestätigt, an HEAD unverändert |
| F-308 | **korrekt übergeben** | `slice-extended-query-replay` §6 trägt das Risiko mit Herkunft und Ausgang |
| F-309 | **korrekt übergeben** | Plan §6 als Frage an den Validator, Ausgang „weiter offen“ |
| F-310 | **korrekt übergeben** | `slice-v1-abschluss-cancel-ohne-schluessel` in `open/`, zwei Liefer-Punkte, Eintrag und Abhängigkeit in `welle-v1-abschluss.md` |
| F-311 | **behoben** | Mutationen X5 (eigene) und H1 rot. `TestE2ERecordExtendedSigtermBeimPipelining` grün. Sonden V-S4 und V-S6 (Abschnitt 3). |
| F-312 | **behoben** | Mutationen W1 und X6 (eigene) rot. `TestHerunterfahrenWeckenNichtVerloren` erzwingt die Verschränkung der Folge-Review-Sonde S2 (Wecken während des zweiten `Shutdown`). |
| F-313 | **behoben** | Mutationen C (`TestCloseMitSchreibfrist`), D und G (`TestRecordEndeNachErfolgreichemUpstream`) rot |
| F-314 | **behoben** | `gofmt -l` steht in der Stufe `test`. Eigene Mutation an einer Datei mit Build-Tag (`test/integration/extended_e2e_test.go`): die Stufe bricht ab und nennt die Datei. `harness/README.md` §Sensors ist nachgezogen. |
| F-315 | **behoben** | Plan §4 sagt, dass die Bedingung eintrat, die Rückführung nicht vollzogen wurde, und warum |
| F-316 | **behoben** | Der Ersatzsatz in §2.3 ist auf den Record-Modus beschränkt |
| F-317 | **korrekt übergeben** | Plan §6, zusammen mit der Validator-Frage aus F-309 |
| F-318 | **behoben** | `beende` schließt die Verbindung (V1 rot in `TestEndeSchliesstVerbindung`), die Spezifikation sagt es in `LH-FA-18.a`. Die Frist der Fehlerantwort ist ungeprüft: V-22. |
| F-319 | **behoben, Rest INFO** | §5 sagt jetzt, dass das Schreiben beim Session-Ende sitzungsübergreifend serialisiert ist. Dass `laufende` dieselbe Sperre in jedem Aufruf nimmt, also während des Schreibens jeder Aufruf jeder Session wartet, sagt der Satz nicht. Er ist wahr, aber knapp. |
| F-320 | **übergeben** | Pflege-Schritt des Registers liegt bei der Closure. `09cc0e9` hat den Fund F-313 als Evidenz nachgetragen. |

---

## 3. `09cc0e9` im Einzelnen

### Code-Lesung

**Herunterfahren.**

1. Der Wächter ruft nur noch `wecke`. Danach ruft die Client-Richtung `Shutdown` am Schleifenkopf vor jedem `Receive`.
2. Eine Nachricht, die `Receive` schon geliefert hat, wird darum verarbeitet.
3. Nach dem ersten `Shutdown` liefern `Query` und eine Nachricht, die eine neue Interaktion beginnt, `ErrShutdown`. Nachrichten der laufenden Extended-Interaktion vor ihrem `Sync` nimmt `ClientMessage` weiter an.
4. Auf `ErrShutdown` liest der Adapter nicht weiter (`wartenAufEnde`) und fragt nach jedem Wecken erneut.

Das entspricht `LH-FA-13.a` und dem Plan §1. Die Prüfreihenfolge in `Query` ist erst PGR-E6001 (Anfrage vor dem `Sync`), dann `ErrShutdown`. Sie hält die ältere Regel aus `LH-FA-18.a`.

**Sperre für Wecken und Zurücksetzen.**

- `wecke` setzt Signal und Frist unter `frist`.
- `weiterlesen` setzt die Frist unter `frist` nur zurück, wenn kein Signal liegt.

Alle Verschränkungen mit `Delivered` → `wecke` enden in einem weiteren Frist-Fehler und einem erneuten `Shutdown`. Ein verlorenes Aufwecken habe ich beim Lesen nicht gefunden. Ein veraltetes Signal kostet höchstens eine zusätzliche Runde.

**`beende`.**

1. Mit einer Meldung setzt es zuerst eine Schreibfrist von 1 s und wartet dann auf `schreiben`.
2. Steht die Server-Richtung mit der Sperre an einem Client, der nicht liest, bricht ihr Schreiben an derselben Frist ab und gibt die Sperre frei. Danach wartet ihr `beende` am `sync.Once`, ohne Verklemmung.
3. Danach folgen `closeRecord`, `conn.Close` und `close(ende)`.

**gofmt.** `gofmt -l .` läuft vor `go vet` über alle Go-Dateien im Kontext, auch über die mit Build-Tag `integration`.

### Eigene Mutationen (nicht aus der Tabelle)

| # | Mutation | Ergebnis |
|---|---|---|
| X1 | `beende`: keine Schreibfrist vor der Fehlerantwort | **grün** in allen Unit-Tests und in allen `TestE2E*` → V-22 |
| X2 | Schleifenkopf ruft `Shutdown` nicht | rot: `TestExtendedHerunterfahren`, `TestHerunterfahrenWeckenNichtVerloren` |
| X3 | `wartenAufEnde` wartet nur auf `ende`, nicht auf `weck` | rot: `TestExtendedHerunterfahren` |
| X5 | `ErrShutdown` in `ClientMessage` nur ohne laufende Interaktion (nicht nach einem `Sync`) | rot: `TestRecordHerunterfahren` |
| X6 | `schreibe` weckt nach `Delivered` == wahr nicht | rot: `TestExtendedHerunterfahren`, `TestHerunterfahrenWeckenNichtVerloren` |
| X7 | auf `ErrShutdown` sofort `EndShutdown`, ohne auf die laufende Interaktion zu warten | rot: `TestExtendedHerunterfahren` |

Aus der Tabelle (dritte Runde) nachgefahren, jede rot im genannten Test: C, D, G, H1 (Zeitüberschreitung), W1, V1, F (eigene Variante, siehe F-314).

### Eigene E2E-Sonden (gegen PostgreSQL 17 über `record`)

| # | Ablauf | Ergebnis |
|---|---|---|
| V-S1 | siehe Punkt 1 | grün, 0,17 s, Aufzeichnung lädt |
| V-S2 | Flush-Gruppe (Prepare) beantwortet, dann `SIGTERM`, 0,7 s später erst `Bind`/`Execute`/`Sync` | Prozess lebt bis zum `Sync`; Antwort korrekt; Ende 14 ms nach dem `ReadyForQuery`, Exit-Code 0; Aufzeichnung trägt Flush- und Sync-Gruppe in Reihenfolge und lädt |
| V-S3 | ruhende pgx-Verbindung nach einer Extended- und einer einfachen Interaktion, `SIGTERM` | Ende nach 25 ms, Exit-Code 0, beide Interaktionen aufgezeichnet, lädt |
| V-S4 | roher Client: `SELECT pg_sleep(1)` mit `Sync`; `SIGTERM` nach 0,2 s; nach 0,4 s eine einfache `Query` | Der Client erhält genau `ParseComplete BindComplete DataRow CommandComplete ReadyForQuery`, dann EOF. Die `Query` ist nicht weitergeleitet und nicht aufgezeichnet. Exit-Code 0, Aufzeichnung lädt. |
| V-S5 | roher Client: 64 MB Ausgabe mit `Sync`, liest nicht; nach 2 s ein `FunctionCall` | ohne Lesen ist PGR-E6001 nach **1,01 s** im Log, die Verbindung ist geschlossen, Exit-Code 6. Mit Mutation X1: **die Session endet nicht** (4 s), solange der Client nicht liest. |
| V-S6 | 20 Läufe: pgx abwechselnd Extended (`pg_sleep(0.01)`) und einfach, `SIGTERM` zu zufälligem Zeitpunkt (100–500 ms, Saat 42) | jeder Lauf endet binnen 3 s mit Exit-Code 0, jede Aufzeichnung lädt |

SIGTERM bei Blockade (`…SigtermNachBlockade`) und beim Pipelining (`…SigtermBeimPipelining`) sind grün. Bei einem Client, der lebt und nicht liest, wartet das Herunterfahren auf das Zustellen. Das ist die offene Validator-Frage F-309/F-317, kein Befund hier.

---

## 4. Plan gegen Code

### §1 Ziel und Abgrenzung

| Plan | Code | Urteil |
|---|---|---|
| Upstream-Port `Query`, `Send`, `Receive`, `Close` mit Gleichzeitigkeits- und Abbruchvertrag | `ports/driven/upstream.go`, `postgres/upstream.go` (Schreibsperre, `TryLock`, `terminateFrist`) | konform |
| Record-Use-Case `Query`, `ClientMessage`, `AwaitServer`, `Delivered`, `Shutdown`, `CloseSession`; Pipelining, `Validate`, Einstufung des Endes | `ports/driving/record.go`, `services/record.go` | konform |
| Herunterfahren: Meldung vor dem nächsten Lesen, danach keine neue Interaktion, Wecken und Zurücksetzen unter einer Sperre | `clientRichtung`, `wecke`/`weiterlesen`, `ErrShutdown` | konform |
| Adapter schließt beim Ende die Client-Verbindung, Fehlerantwort höchstens 1 s | `beende` | konform im Code (V-S5); ohne fangenden Test (V-22) |
| Gate `gofmt -l` in der Stufe `test` | `Dockerfile` | konform |
| Replay lehnt Extended mit PGR-E6001 ab | `replaySitzung`; Test „Replay“ in `server_extended_test.go` | konform |

### §3 Dateien — **Abweichung (V-23)**

Jede Zeile in §3 hat Änderungen im Diff. Umgekehrt fehlen Zeilen für Dateien aus `git diff --stat 5d666db^ 09cc0e9`:

- `docs/plan/planning/open/slice-extended-query-replay.md`, der Folge-Slice, geändert in `5d666db` und `0dbcf9d`
- die drei Evidenz-Dateien unter `docs/plan/planning/observations/BEO-REPO/*/evidence/`
- `docs/plan/adr/0030-full-duplex-im-record-pfad.md` und `docs/plan/adr/README.md` (in `0dbcf9d`)

Die Zeile `docs/reviews/` nennt nur die Mutationstabelle, mitcommittet sind auch beide Review-Reports.

### §4 Trigger

Konform; siehe F-315.

### §6 Risiken

Jede Behebung aus `0dbcf9d`/`09cc0e9` hat ein Risiko mit Beleg. Die genannten Tests existieren und sind grün. Eine Ausnahme: Das Risiko F-318 nennt `TestEndeSchliesstVerbindung` als Beleg auch für „die Fehlerantwort dazu hat eine Frist von einer Sekunde“. Der Test endet über `Terminate` und schreibt keine Fehlerantwort (V-22).

### Kommentare und Abdeckungs-Deklarationen

Die neuen Deklarationen sagen zu, was ihre Tests prüfen:

- `TestRecordHerunterfahren`
- `TestEndeSchliesstVerbindung`
- `TestE2ERecordExtendedSigtermBeimPipelining`

Die Deklaration von `TestExtendedHerunterfahren` nennt den neuen Fall „neue Interaktion nach dem Beginn“ nicht. Sie sagt also weniger zu, als der Test prüft, und ist kein Befund.

Kommentare, die mehr zusagen, als geprüft ist: der Kommentar an `beende` („höchstens meldeFrist“) und an `meldeFrist`. Beide gehören zu V-22.

---

## 5. ADR-Konformität

| ADR | Zusage | Urteil |
|---|---|---|
| [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) | zwei unabhängige Richtungen; Zustand und Entscheidungen im Service; Abbruchvertrag an beiden Ports, mit Fakes geprüft; einfache Anfragen synchron | **konform.** `wartenAufEnde` entscheidet nichts. Es führt `ErrShutdown` aus und fragt den Service. Fitness Functions: Fake-Tests mit begrenztem Puffer grün, Integrationstest mit großer Gruppe und Signal grün, `a-check` 0 Befunde. Konsequenz „endet beim Herunterfahren, auch während Daten fließen“: für pipelinende Clients jetzt gegeben (E6, V-S6); bei einem lebenden, nicht lesenden Client nicht (F-309/F-317, übergeben). |
| [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) | Gruppen je `Flush`/`Sync`, zuerst Client-, dann Server-Nachrichten | **konform** (M1, M3, N12 rot; `inReihe` in den E2E-Tests). Späte Flush-Antworten in der Folgegruppe: F-308, übergeben. |
| [ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md) | PGWire-Server ist Driving Adapter, ruft Inbound Ports | **konform.** Der Core ruft den Adapter nicht zurück. Die Ereignisse gehen über `driving.Recorder`. |
| [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md) | Upstream als Driven Adapter, Gegendruck im Port-Design | **konform.** Kein Puffer über die Gruppe hinaus. Gegendruck über die Verbindungen ist belegt (Gegendruck-E2E, V-S1, V-S5). |

---

## 6. Spec-Konformität

| Stelle | Urteil |
|---|---|
| [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) | **konform** für die Record-Hälfte (Punkt 1). Die Replay-Hälfte ist Gegenstand von `slice-extended-query-replay`. |
| [LH-FA-06](../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten) | **konform.** Anfragen und Antworten stehen in Reihenfolge in der Aufzeichnung (V-S1, V-S2, `inReihe`). |
| [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) | **konform.** Exit-Codes 0 (V-S2 bis V-S4, V-S6), 4 (`…Abbruch`, `…SigtermNachBlockade`), 6 (V-S5). |
| [LH-FA-02](../../spec/lastenheft.md#lh-fa-02--record-modus) | **konform.** Eine nicht abgeschlossene Interaktion wird nicht übernommen, die vorherigen bleiben (`…Abbruch`, `TestRecordExtendedEnde`). |
| `LH-FA-18.a` | **konform** (Gruppen, Pipelining, Abbruch, Herunterfahren bis `ReadyForQuery`, Ende schließt die Verbindung). |
| `LH-FA-13.a` | **konform** für „danach keine neue Interaktion“ (V-S4, E6) und „laufende nimmt bis `Sync` an“ (V-S2). Randform: siehe V-24. |
| `SPEC-041` | **konform.** Jede Aufzeichnung der Sonden lädt mit dem Leser und besteht `Validate`. |
| `spec/architecture.md` §2.3, §4, §4.5, §5 | **konform.** Zu §5 der Rest aus F-319 (Abschnitt 2). Hard Rule 3.4 hält: Die Änderungen in `09cc0e9` nennen keine ADR, keinen Slice und keinen Code. |

**Rückwärts (einfache Anfragen):**

- Alle `TestE2ERecord*` für einfache Anfragen sind grün. V-S3 und V-S6 mischen beide Arten.
- Geändert hat sich nur, was die Spezifikation jetzt sagt: Eine `Query` nach dem Beginn des Herunterfahrens wird nicht weitergeleitet (`LH-FA-13.a`), und das Ende einer Session schließt die Verbindung sofort (`LH-FA-18.a`, für einfache Anfragen ohne sichtbare Folge, weil `handle` die Verbindung danach ohnehin schloss).
- Das Handbuch nennt beides.

---

## 7. Befunde des Verifiers

| ID | Kategorie | Befund | Beleg |
|---|---|---|---|
| V-22 | MEDIUM | Die Zusage „eine Fehlerantwort beim Ende schreibt der Adapter höchstens eine Sekunde lang“ hat keinen fangenden Test. Sie steht in Plan §1, in §6 (Risiko F-318) und in den Kommentaren an `beende` und `meldeFrist`. Ohne die Schreibfrist (Mutation X1) bleiben alle Unit-Tests und alle `TestE2E*` grün. Die Frist trägt aber Last: Mit X1 wartet `beende` auf `schreiben`, das die Server-Richtung an einem Client hält, der nicht liest. Die Session endet dann nicht, und das Herunterfahren hinge mit (V-S5: 1,01 s an HEAD, kein Ende mit X1). Plan §6 nennt `TestEndeSchliesstVerbindung` als Beleg. Dieser Test endet über `Terminate` ohne Fehlerantwort. Das ist der dritte Fund von `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` im selben Vorgang (nach F-303 und F-313), diesmal im Korrektur-Commit. Er gehört auch zur Klasse `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. Vorschlag: ein Unit-Test wie `TestEndeSchliesstVerbindung` mit einer nicht unterstützten Nachricht (etwa `FunctionCall`) statt `Terminate`, Sitzung kehrt binnen 2 s zurück; danach §6 berichtigen. | Mutation X1; Sonde V-S5; `internal/adapters/driving/pgwire/server.go` (`beende`, `meldeFrist`); Plan §1, §6 |
| V-23 | LOW | Plan §3 führt nicht jede geänderte Datei. Es fehlen der Folge-Slice `docs/plan/planning/open/slice-extended-query-replay.md`, die drei Register-Evidenzen, die Datei von [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) und der ADR-Index. Die Zeile `docs/reviews/` nennt nur die Mutationstabelle. `AGENTS.md` §3.9 verlangt, dass der Plan dem Diff folgt. Im Vorgänger-Slice war genau das ein Finding (F-291) und wurde dort mit Zeilen für Folge-Slice und Register-Belege behoben. | `git diff --stat 5d666db^ 09cc0e9`; Plan §3 |
| V-24 | INFO | Den „Beginn des Herunterfahrens“ aus `LH-FA-13.a` setzt der Code beim ersten `Shutdown` der Client-Richtung, nicht beim Signal. Zwei Randformen folgen daraus. (a) Eine Nachricht, die `Receive` zwischen Signal und Wecken liefert, wird noch verarbeitet, auch eine neue `Query`. Das Fenster ist kurz. (b) Eine Nachricht, die schon vollständig im Lesepuffer liegt, aber noch nicht geliefert ist, wird verworfen, wenn keine Interaktion läuft. Beides liegt im Wortlaut „schon vollständig gelesen“, den F-317 als an innere Zeitpunkte gebunden meldet. Kein eigener Handlungsbedarf, nur Material für die Validator-Frage F-309/F-317. | Code-Lesung `clientRichtung`, `recordSitzung` (Wächter) |

---

## 8. Gegenprobe dieser Datei

`make docs-check` nach dem Anlegen dieser Datei: `d-check: 154 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

## 9. Gesamturteil

**Alle drei DoD-Liefer-Punkte sind bestätigt.** Grundlage sind eigene Läufe:

- `make gates` mit Exit 0
- die Stufe `test` ohne Cache
- der Race-Detector 20-fach
- 18 Mutationen: 6 eigene, 12 aus der Tabelle nachgefahren. 17 sind rot im erwarteten Test, eine (X1) bleibt grün → V-22.
- 6 eigene E2E-Sonden gegen PostgreSQL 17, alle grün an HEAD

Die Findings F-301 bis F-319 sind behoben oder mit Artefakt übergeben (F-308, F-309/F-317, F-310). F-320 liegt bei der Closure. Der Code hält [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md), [ADR-0012](../plan/adr/0012-extended-query-gruppen.md), [ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md) und [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md) ein. Er entspricht `LH-FA-18.a`, `LH-FA-13.a`, `SPEC-041` und der Architektur-Sicht. Einfache Anfragen verhalten sich unverändert, außer wo die Spezifikation es jetzt anders sagt.

**Kein Befund ist HIGH.** Vor der Closure empfohlen:

1. V-22: Test für die Frist der Fehlerantwort bei nicht lesendem Client, Plan §6 berichtigen — Implementer.
2. V-23: Plan §3 um die fehlenden Dateien ergänzen — Implementer/Planner.
3. V-24 zur Validator-Frage F-309/F-317 legen — Planner.
4. Über ein Review für `09cc0e9` entscheiden (neue Nebenläufigkeit beim Herunterfahren und beim Ende) — Planner.
