# Review-Report: slice-extended-query-record — 2026-10-05

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** Commit `5d666db` (`git show 5d666db`). Die mitcommittete Mutationstabelle [`2026-10-05-mutationen-slice-extended-query-record.md`](2026-10-05-mutationen-slice-extended-query-record.md) ist nur als Eingang gelesen. Inhaltliche Schwerpunkte: Port-Zuschnitt, Nebenläufigkeit in `sitzung`, Protokolltreue der Gruppierung, Tests je Zusage, Rückwärtsverhalten einfacher Anfragen.

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-record.md` (Kopf, §1 bis §4, §6), `docs/plan/planning/open/slice-extended-query-replay.md` §1
- `spec/lastenheft.md` [LH-FA-02](../../spec/lastenheft.md), [LH-FA-04](../../spec/lastenheft.md), [LH-FA-05](../../spec/lastenheft.md), [LH-FA-06](../../spec/lastenheft.md), [LH-FA-13](../../spec/lastenheft.md), [LH-FA-15](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md)
- `spec/spezifikation.md`: LH-FA-02.b, LH-FA-05.a, LH-FA-05.e (Tabelle Protokollrand), LH-FA-13.a, LH-FA-13.b, LH-FA-18.a samt den Sätzen und der Historienzeile vom 2026-10-05, `SPEC-002`, `SPEC-041`
- `spec/architecture.md` `ARC-002` bis `ARC-004` mit §2.3 (Ports)
- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0012](../plan/adr/0012-extended-query-gruppen.md)
- `AGENTS.md` (Hard Rules 3.3 bis 3.9); Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/spec-randform-erst-im-review-entschieden`
- Vorherige Reports am Modul: [`2026-10-04-folge-review-slice-extended-query-modell.md`](2026-10-04-folge-review-slice-extended-query-modell.md) (bis F-300)
- Quelltext von `pgproto3` v5.11.0 (`frontend.go`, `backend.go`, `cancel_request.go`), gelesen im Modul-Cache der Stufe `deps`

**Ausgeführte Läufe im Repo:** `make docs-check` nach dem Anlegen dieser Datei (siehe unten). `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert (`git status --short` leer). Den Race-Detector-Lauf aus der Mutationstabelle habe ich nicht wiederholt, weil er eine Host-fremde Toolchain (`gcc`) im Image braucht.

**Mutationen und Sonden:** Nur in Kopien (`git archive 5d666db`) im Scratchpad. Images trugen die Tags `pgr-review*`, das Netz und der PostgreSQL-Container (gepinntes Image aus `harness/mk/integration.mk`) eigene Namen. Images, Container, Netze und Kopien sind gelöscht.

*Sonden* (zusätzliche Testdatei in `test/integration`, Stufe `integration`, gegen PostgreSQL 17):

| # | Ablauf | direkt gegen PostgreSQL | über `record` |
|---|---|---|---|
| P1 | pgx im Standardmodus, `SendBatch` mit zwei Anfragen in **einer** Sync-Gruppe: `SELECT repeat('x', 1000000) FROM generate_series(1, 50)` (≈ 50 MB Ausgabe) und `SELECT length($1::text)` mit 20 MB Parameter | grün nach 0,13 s | **hängt**: Zeitüberschreitung nach 25 s bei der ersten Anfrage |
| P2 | `pgconn`-Pipeline: Flush-Gruppe mit verzögerter großer Ausgabe (`pg_sleep(2)`, 50 MB), danach ohne Warten eine Sync-Gruppe mit 20 MB Parameter | grün nach 2,1 s | **hängt**: Zeitüberschreitung nach 40 s |
| P2a | wie P2, erste Gruppe klein (`SELECT 1`), zweite 20 MB | — | grün |
| P2b | wie P2, erste Gruppe 50 MB, zweite klein | — | grün |
| P3 | wie P1; nach der Zeitüberschreitung schließt der Client die Verbindung, dann `SIGTERM` an den Recorder | — | **der Prozess endet nicht** innerhalb von 15 s |

Im Log von P3 steht zudem `PGR-E6001 … Startnachricht nicht lesbar: cancel request too short` (siehe F-310).

*Mutationen* (einzeln):

| # | Mutation | Ergebnis |
|---|---|---|
| R1 | `sitzung` setzt nach dem `ReadyForQuery` `gesendet` nicht zurück; nach jeder Interaktion läuft ein weiteres `AwaitServer` | Unit grün; E2E rot (`TestE2ERecordExtendedPgx`, Zeitüberschreitung bei der folgenden einfachen Anfrage) |
| R2 | `session.Close` schließt die Verbindung nicht, sendet nur `Terminate` | Unit grün (E2E nicht gelaufen; gegen einen realen Server schließt der Server nach `Terminate` selbst) |
| R7 | `Receive` liest immer bis `ReadyForQuery` und ignoriert den Lesepuffer | Unit rot (Paket `postgres`, Zeitüberschreitung) |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-301 | HIGH | Der Record-Pfad blockiert sich selbst, wenn eine Gruppe groß ist und der Server zugleich viel ausgibt. `ClientMessage` sendet die Gruppe synchron (`Send`, blockierendes Schreiben). Erst danach startet `sitzung` das `AwaitServer`, und jedes `AwaitServer` liest nur einen Stapel. Solange `Send` blockiert, liest niemand vom Upstream. Der Server blockiert beim Schreiben, liest nicht weiter, und `Send` kommt nie zurück. Betroffen ist pgx im Standardmodus mit einem Batch, also eine einzige Sync-Gruppe (P1), ebenso eine Flush-Gruppe mit nachgeschickter Gruppe (P2); direkt gegen PostgreSQL laufen beide in Sekundenbruchteilen. Weil `sitzung` in `ClientMessage` steht, endet die Session auch nicht, wenn der Client schließt, und der Prozess endet nicht nach `SIGTERM` (P3). LH-FA-18.a sagt zu, dass der Recorder alle Nachrichten weiterleitet, „auch wenn der Client mehrere Nachrichten sendet, ohne auf Antworten zu warten“. [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md) nennt als Konsequenz: „Streaming und Backpressure müssen im Port-Design mitgedacht werden.“ Kein Test und kein Risiko in §6 nennt den Fall. Risiko 5 nennt nur die Gleichzeitigkeit von `Send` und `Receive`, nicht dass sie fehlt, solange `Send` blockiert. | [LH-FA-18](../../spec/lastenheft.md) Happy, [LH-FA-04](../../spec/lastenheft.md), [LH-FA-13](../../spec/lastenheft.md) Boundary; LH-FA-18.a „Record“; [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md) Konsequenzen; Reviewer-Skill HIGH „Korrektheitsfehler im kritischen Pfad“ | `internal/hexagon/services/record.go:152`; `internal/adapters/driving/pgwire/server.go:204`, `:288`; `internal/adapters/driven/postgres/upstream.go:102`, `:119` | ja — Sonden P1, P2, P3 | Gegendruck fehlt: Senden blockiert, solange niemand liest |
| F-302 | MEDIUM | Plan §4 nennt die Rückführung „Der Upstream-Adapter braucht eine neue Port-Operation — Entscheidung klären“. Ausgelöst wurde sie nicht; Begründung in §1 und §3: `ARC-004` nenne die Operation schon. Für `Send` trägt diese Lesart: `ARC-004` sagt „je Gruppe … die Client-Nachrichten senden“. Für `Receive` trägt sie nicht. `ARC-004` sagt „je Gruppe … die Server-Nachrichten liefern“, also eine an die Gruppe gebundene Antwort. `Receive` ist von der Gruppe gelöst, liefert „die nächsten“ Nachrichten und hat einen eigenen Vertrag zur Gleichzeitigkeit („darf gleichzeitig mit Send und Close laufen, nicht mit Query“). Dazu kommt auf der Driving-Seite `AwaitServer`, das gleichzeitig mit `ClientMessage` laufen darf. Der Zuschnitt entscheidet, wie Gegendruck zwischen Senden und Empfangen entsteht; das ist genau die offene Frage aus [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md), und F-301 ist ihre Folge. Die Entscheidung liegt beim Architect, nicht bei der Implementation. | Plan §4 (Rückführung); `ARC-004` (`spec/architecture.md` §2.3); [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md) | `docs/plan/planning/in-progress/slice-extended-query-record.md:88`, `:34`, `:70`; `internal/hexagon/ports/driven/upstream.go:22`–`:30`; `internal/hexagon/ports/driving/record.go:28`–`:32`; `spec/architecture.md:199` | ja (Lesen) | Rückführungs-Trigger durch eigene Lesart entschärft |
| F-303 | MEDIUM | Port-Zusagen ohne fangenden Test. (a) `UpstreamSession.Receive`: „Close beendet es mit einem Fehler“ und „darf gleichzeitig mit Send … laufen“. Kein Test ruft `Close` oder `Send`, während `Receive` wartet. R2 (Close schließt die Verbindung nicht) bleibt in allen Unit-Tests grün. (b) Der Adapter verlässt sich darauf, dass `CloseSession` ein wartendes `AwaitServer` beendet; sonst bleibt dessen Goroutine stehen (`server.go:205`). Der Driving-Port sagt das nicht zu: `AwaitServer` „darf gleichzeitig mit ClientMessage laufen“, von `CloseSession` steht dort nichts. `fakeRecorder.AwaitServer` bildet es nicht ab. In `TestExtendedAbbruch` („Terminate nach Flush“) und `TestExtendedHerunterfahren` bleibt das gestartete `AwaitServer` deshalb nach Testende im Kanal hängen, ohne dass ein Test das bemerkt. (c) Ein zusätzliches `AwaitServer` nach dem `ReadyForQuery` (R1) fängt kein Unit-Test, nur die E2E-Folge mit einer einfachen Anfrage. | `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`; Reviewer-Skill MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag“; `ARC-003`, `ARC-004` | `internal/hexagon/ports/driven/upstream.go:25`–`:29`; `internal/hexagon/ports/driving/record.go:28`–`:31`, `:39`–`:41`; `internal/adapters/driven/postgres/upstream_extended_test.go`; `internal/adapters/driving/pgwire/server_test.go:45` | ja — R1, R2 | Vertragszusage ohne fangenden Test |
| F-304 | MEDIUM | Der Driving Adapter führt eine eigene Abbildung des Interaktionszustands (`offen`, `gesendet`, `sync`), parallel zu `laufend.offen` im Record-Service. Auf dieser Abbildung trifft er Record-Entscheidungen: (1) ob ein Verbindungsende oder `Terminate` `PGR-E4003` ist (LH-FA-18.a „Abbruch“); (2) ob das Herunterfahren wartet (LH-FA-13.a); (3) ob `EndLost` zu `EndNormal` herabgestuft wird. Die Begründung dafür steht nur im Kommentar des Adapters („Die laufende Interaktion verwirft CloseSession ohnehin“). `ARC-003` sagt: „Der PGWire-Adapter übernimmt keine Record- oder Replay-Fachlogik.“ Die beiden Zustände stimmen heute nur deshalb überein, weil beide denselben Aufrufen folgen. Der Replay-Slice übernimmt `sitzung` laut seinem §1 und damit dieselbe Doppelung. | `ARC-003` (`spec/architecture.md` §2.3); [ADR-0001](../plan/adr/0001-hexagonale-architektur.md); Maintainability | `internal/adapters/driving/pgwire/server.go:189`–`:193`, `:196`, `:211`–`:214`, `:240`–`:245`, `:263`–`:272`; `internal/hexagon/services/record.go:33`–`:35` | ja (Lesen) | Fachzustand doppelt geführt, Regel im Driving Adapter |
| F-305 | LOW | Beim Herunterfahren verwirft `sitzung` eine Client-Nachricht, die der Leser schon gelesen hat. Das geschieht, wenn sie beim Prüfen am Schleifenkopf in `clientCh` liegt oder wenn `select` `stopp` wählt. Das gilt auch für eine einfache `Query`: Der Client erhält keine Antwort und keine Fehlerantwort, und es wird kein Fehler gemerkt. Vor dem Commit lief eine vor der Lesefrist vollständig gelesene `Query` zu Ende. Ob eine empfangene, noch nicht übergebene Nachricht „laufende Interaktion“ im Sinne von LH-FA-13.a ist, regelt die Spezifikation nicht; der Code entscheidet es still. | LH-FA-13.a; `BEO-REPO/spec-randform-erst-im-review-entschieden`; Reviewer-Skill MEDIUM/LOW-Grenze (zeitabhängig, nur beim Herunterfahren) | `internal/adapters/driving/pgwire/server.go:196`–`:199`, `:211`–`:219` | nein (zeitabhängig; Code-Lesung gegen `git show 5d666db^:internal/adapters/driving/pgwire/server.go`) | Randform still im Code entschieden |
| F-306 | LOW | Die Abdeckungs-Deklaration von `TestRecordExtendedSyncGruppe` trägt nach einer Teilersetzung ein doppeltes Prädikat („mit dem ReadyForQuery steht die Interaktion steht sie als type extended …“). Über `make abdeckung` steht der Satz wörtlich in der Nutzerdoku. Er bleibt lesbar und trägt die Klasse Zusage; der HIGH-Fall des Skills (abgebrochener Satz ohne Klasse) liegt nicht vor. | `AGENTS.md` §3.7; Maintainability | `internal/hexagon/services/record_extended_test.go:74`; `docs/user/abdeckung-unit.md:40` | ja (Lesen) | Teilersetzung hinterlässt Satzrest |
| F-307 | LOW | Plan §1 begründet das Verhalten im Replay mit abwesendem Code: „wie zuvor die Schleife des Adapters“. Die Schleife gibt es nach diesem Commit nicht mehr. Der Satz beschreibt Herkunft statt Zustand; Herkunft hält `git`. | `AGENTS.md` §3.7 (sinngemäß für Planungstext), §3.9; Maintainability | `docs/plan/planning/in-progress/slice-extended-query-record.md:38` | ja (Lesen) | Planungstext beschreibt abwesenden Code |
| F-308 | INFO | Die neue Regel „eine danach eintreffende [Server-Nachricht] gehört zur folgenden Gruppe“ ist in sich stimmig und widerspricht keinem Satz des Lastenhefts. Sie hat aber eine Folge für das Replay, die die Spezifikation nicht nennt. Antworten auf Flush-Gruppe N, die im Record nach der ersten Nachricht von N+1 eintrafen, gibt das Replay erst nach dem `Flush`/`Sync` von N+1 frei. Ein Client, der zwischen Beginn und Ende von N+1 auf die Antwort auf N wartet, hängt im Replay mit genau der Aufzeichnung, die er selbst erzeugt hat. Das liegt innerhalb der bekannten Grenze „Flush ohne wartenden Client“. Der Satz „die Wiedergabe einer vorhandenen Aufzeichnung bleibt deterministisch“ ist wahr, sagt aber nicht, dass sie den Client anhalten kann. Zudem sagt [ADR-0012](../plan/adr/0012-extended-query-gruppen.md), eine Gruppe enthalte die Server-Nachrichten, „die auf sie antworten“; späte Antworten auf N stehen nun in N+1. Hinweis an Planner und `slice-extended-query-replay`. | LH-FA-18.a „Gruppen“; [ADR-0012](../plan/adr/0012-extended-query-gruppen.md); [LH-FA-18](../../spec/lastenheft.md) Happy | `spec/spezifikation.md:703`–`:708` | ja (Lesen) | Spec-Regel mit unbenannter Folge im anderen Modus |
| F-309 | INFO | Beim Herunterfahren wartet eine Session ohne Frist auf das `ReadyForQuery` einer laufenden Extended-Interaktion (LH-FA-18.a, neu). Das ist mit LH-FA-13.a vereinbar und als Risiko 4 in §6 benannt. Ein Client, der nach `Flush` ruht, etwa eine Verbindung in einem Pool, hält den Prozess also beliebig lange. [LH-FA-13](../../spec/lastenheft.md) Boundary setzt voraus, dass das kontrollierte Beenden einen Prozessstatus liefert. Ob ein unbegrenztes Warten den Bedarf in CI trifft ([LH-FA-15](../../spec/lastenheft.md)), ist eine Frage an den Validator, nicht an dieses Review. | LH-FA-13.a; LH-FA-18.a „Abbruch“; [LH-FA-13](../../spec/lastenheft.md), [LH-FA-15](../../spec/lastenheft.md) | `spec/spezifikation.md:671`–`:673`; Plan §6 Risiko 4 | nein | Wartezeit ohne Obergrenze beim Herunterfahren |
| F-310 | INFO | Außerhalb des Diffs, in Sonde P3 beobachtet. Der Recorder gibt `BackendKeyData` nicht an den Client weiter. pgx sendet bei Kontext-Abbruch deshalb einen `CancelRequest` ohne Schlüssel (12 Byte). `pgproto3` v5.11.0 lehnt ihn mit „cancel request too short“ ab, und der Recorder meldet `PGR-E6001` mit FATAL-Antwort und merkt sich die Klasse 6. LH-FA-05.e verlangt für `CancelRequest` die Warnung `PGR-W3001` und keinen Verbindungsfehler. `TestCancelRequest` sendet 16 Byte (mit Schlüssel) und deckt den Fall nicht ab. Betroffen ist jeder pgx-Client mit Zeitüberschreitung, einfach wie Extended. Übergabe an Implementer und Planner als eigener Gegenstand. | LH-FA-05.e (Tabelle Protokollrand, `PGR-W3001`); [LH-FA-13](../../spec/lastenheft.md) Negative | `internal/adapters/driving/pgwire/server.go:425`; `internal/adapters/driven/postgres/upstream.go:55`; `internal/adapters/driving/pgwire/server_test.go:278` | ja — Sonde P3 (Log) | Randform des Protokollrands nicht am realen Client geprüft |

## Antwort auf die Schwerpunkte

1. **Architektur (Rückführung).** Die Lesart trägt für `Send`, für `Receive` und `AwaitServer` nicht (F-302). Die Core-Reinheit ist gewahrt: Ports und Service importieren nur Domain und Ports, `pgproto3` liegt nur in den beiden PGWire-Adaptern. Die Verantwortungsgrenze aus `ARC-003` verletzt der Adapter mit seinem gespiegelten Interaktionszustand (F-304).
2. **Nebenläufigkeit.** Datenrennen fand ich beim Lesen keines. Die Felder von `laufend` berührt nur die Goroutine von `sitzung`. `AwaitServer` liest nur `upstream`, über `laufende` unter `mu`. `Frontend.Send`/`Flush` (`wbuf`, `encodeError`) und `Receive` (`cr`, `bodyLen`, `msgType`, `partialMsg`, Flyweights) haben in v5.11.0 getrennte Felder; Gleiches gilt für `Backend` mit Leser-Goroutine und `send`/`fail`. Goroutinen: Der Client-Leser endet über `conn.Close` in `handle`, `AwaitServer` über `upstream.Close` in `CloseSession` (siehe F-303 b). `anfrage` ist ungepuffert, und `sitzung` fordert erst an, wenn das vorige Ergebnis verbraucht ist; das ist verklemmungsfrei. Es gibt aber einen Deadlock zwischen Senden und Empfangen am Upstream (F-301). Kontext: `AwaitServer` läuft mit `WithoutCancel`, `Send` und `Receive` ignorieren `ctx`. Damit gelingt das Senden des `Sync` auch nach dem Abbruch, und die Herunterfahren-Regel trägt (`TestExtendedHerunterfahren`). Zur Randform beim Herunterfahren siehe F-305.
3. **Protokolltreue.** Pipelining über mehrere Syncs: `sitzung` liest nach einem `Sync` erst nach dessen `ReadyForQuery` weiter. Die Zuordnung folgt dem Draht und hängt nicht vom Zeitverhalten ab, solange kein Deadlock eintritt (F-301). Flush-Gruppen: Server-Nachrichten gehen erst nach `ServerMessage` an den Client. Ein Client, der wartet, kann darum keine Nachricht vor der Antwort senden, und die Aufzeichnung ist für ihn zeitunabhängig, wie LH-FA-18.a zusagt. Fehler mitten in einer Gruppe: Die verworfenen Client-Nachrichten stehen in der Aufzeichnung, Flush-Gruppen ohne Antwort als `server: []` (`TestRecordExtendedFlushGruppen`, `TestE2ERecordExtendedPipeline`). `Query` während einer Extended-Interaktion ist `PGR-E6001`, ohne Weiterleitung. Eine Zielart außer `S`/`P` ist `PGR-E6001`, bevor der Use Case etwas sieht. Die neuen Sätze der Spezifikation sind in sich stimmig und widersprechen keinem Satz des Lastenhefts. Zur unbenannten Folge für das Replay siehe F-308, zur Wartezeit beim Herunterfahren F-309.
4. **Tests je Zusage.** Die 26 Zeilen der Mutationstabelle und die Kommentare der Tests passen zueinander, soweit ich sie gelesen habe. Ohne fangenden Test bleiben die Zusagen aus F-303. Zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` keine neue Belegstelle: Die Abdeckungs-Deklarationen sagen zu, was ihre Tests prüfen, F-306 betrifft nur den Wortlaut. Randformen, die weiter offen sind: F-305 und F-310.
5. **Rückwärts.** Einfache Anfragen gehen jetzt auch durch `Validate`. Abgelehnt werden nur Antworttypen, die der Upstream-Adapter vorher schon als `PGR-E6001` abwies (`simpleResponses`). Meldung an den Client und Klasse bleiben gleich, der Text der Meldung ändert sich. Ausgabe der Aufzeichnung und Exit-Codes einfacher Anfragen sind unverändert (E2E `TestE2ERecord*` laut Implementer grün; R1 zeigt, dass `TestE2ERecordExtendedPgx` die Folge Extended → einfach prüft). Abweichung nur beim Herunterfahren (F-305).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/ports/driven`, `internal/hexagon/ports/driving` | geprüft. Fachliche Typen, keine Bibliothekstypen ([ADR-0001](../plan/adr/0001-hexagonale-architektur.md), `ARC-004`). Befund F-302, F-303. |
| `internal/hexagon/services` | geprüft. Gruppierung nach LH-FA-18.a; `Validate` vor jeder Übernahme; `Query` in laufender Interaktion `PGR-E6001`; Sequenz zählt fortlaufend über beide Arten; `CloseSession` übernimmt keine offene Interaktion. Unit-Tests mit Mutationen M1 bis M7 der Tabelle. Befund F-301 (synchrones `Send`), F-306. |
| `internal/adapters/driving/pgwire` | geprüft. Abbildung `toClientMessage` und `toMessage` belegt je Typ genau die Felder aus `SPEC-041`, kopiert Werte; leerer Parameter bleibt von NULL getrennt. Replay lehnt Extended weiter mit `PGR-E6001` ab. Befund F-301, F-303, F-304, F-305. |
| `internal/adapters/driven/postgres` | geprüft. `toResponse` und `toFrontendMessage` vollständig für die Nachrichten aus LH-FA-18.a; `Receive` endet bei `ReadyForQuery` (R7 rot). Befund F-301, F-303; außerhalb des Diffs F-310. |
| `test/integration` | geprüft. Drei E2E-Tests mit Abdeckungs-Deklaration; die Reihenfolgeprüfung `inReihe` und `laedt` stützen die Zusagen. Kein Test mit großer Gruppe oder großer Ausgabe (F-301). |
| `spec/` — Strata, Hard Rule 3.4 | geprüft, ohne Befund zu Strata. Nur `spezifikation.md` geändert; die neuen Sätze und die Historienzeile nennen Spec-Kennungen und Lastenheft-Kennungen, keine ADR, keinen Slice, keinen Commit. `architecture.md` und `lastenheft.md` unverändert. Inhaltlich F-308, F-309. |
| `docs/user/` | geprüft. Abdeckungstabellen entsprechen den Deklarationen; Handbuch nennt `PGR-E6001` für einfache Anfrage in laufender Folge und das Warten beim Beenden. Befund F-306. |
| `docs/plan/planning/` — Hard Rule 3.9 | geprüft. Kopf, §1, §3, §6 und der Folge-Slice folgen dem Code. Befund F-302 (§4), F-307. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. `5d666db` enthält keinen Move; die Moves liegen in eigenen Commits davor. |
| Hard Rule 3.5, 3.8 — ADRs | geprüft, ohne Befund. Keine ADR im Diff. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund. Keine Gate-Konfiguration, kein `.a-check.yml`, kein `harness/mk/*` im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft. Neue Kommentare in Code und Tests sind Zusagen oder Kopplungen im Indikativ; keine verworfene Alternative, kein abgebrochener Satz. Befund F-306, F-307. |
| Commit-Message `5d666db` | geprüft, ohne Befund. Nennt `slice-extended-query-record`, [LH-FA-18](../../spec/lastenheft.md), [LH-FA-06](../../spec/lastenheft.md), [LH-FA-13](../../spec/lastenheft.md), [ADR-0012](../plan/adr/0012-extended-query-gruppen.md), [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md); keine `SPEC-*`- oder `ARC-*`-Kennung. |
| Register `BEO-REPO/*` | geprüft. F-303 ist ein weiterer Fall von `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, F-305 von `BEO-REPO/spec-randform-erst-im-review-entschieden` (hier für einen Protokollrand, nicht ein Format). |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:**

- Gegendruck fehlt: Senden blockiert, solange niemand liest (F-301)
- Rückführungs-Trigger durch eigene Lesart entschärft (F-302)
- Vertragszusage ohne fangenden Test (F-303; Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`)
- Fachzustand doppelt geführt, Regel im Driving Adapter (F-304)
- Randform still im Code entschieden (F-305; Register `BEO-REPO/spec-randform-erst-im-review-entschieden`)
- Teilersetzung hinterlässt Satzrest (F-306)
- Planungstext beschreibt abwesenden Code (F-307)
- Spec-Regel mit unbenannter Folge im anderen Modus (F-308)
- Wartezeit ohne Obergrenze beim Herunterfahren (F-309)
- Randform des Protokollrands nicht am realen Client geprüft (F-310)

## Verdikt

**Merge-blockierend:** ja, wegen F-301. pgx im Standardmodus mit einem Batch aus großer Ausgabe und großem Parameter hängt über `record` dauerhaft, direkt gegen PostgreSQL nicht. Der Prozess endet danach auch nicht mehr nach `SIGTERM`. Das trifft die Zusage „Pipelining“ aus LH-FA-18.a und das kontrollierte Beenden nach [LH-FA-13](../../spec/lastenheft.md).

Abgesehen davon ist der Diff in sich geschlossen. Gruppierung, Zuordnung und Abbruchregeln folgen LH-FA-18.a, die neuen Spec-Sätze sind stimmig, und fast jede neue Zusage hat einen Test, den eine Mutation rot färbt.

**Übergabe:**

- F-301 und F-302 hängen zusammen. Die Behebung von F-301 berührt den Zuschnitt von `Receive`/`AwaitServer`, also genau die Frage, die Plan §4 an eine Entscheidung bindet. Übergabe an den Architect mit diesem Report als Artefakt. Widerspricht der Implementer der HIGH-Einstufung von F-301, läuft der Konflikt-Pfad als Sequenz über den Architect (Modul 8): Implementer-Stellungnahme als Datei unter `docs/reviews/`, Entscheidung des Architects als ADR oder Plan-Änderung. Ein Herabstufen allein auf Widerspruch hin findet nicht statt.
- F-303, F-304, F-305, F-306, F-307 an den Implementer; Plan §6 nimmt F-301 und F-305 auf (Hard Rule 3.9).
- F-308 an den Planner und an `slice-extended-query-replay`; F-309 an den Validator; F-310 an Implementer und Planner als eigener Gegenstand außerhalb dieses Slice.
- Die Finding-Klassen gehen in die Closure-Notiz §7. Die DoD-Konformität prüft der Verifier separat.
