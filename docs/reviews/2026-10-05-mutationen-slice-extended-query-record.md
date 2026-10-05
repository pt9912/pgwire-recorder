# Mutationen: slice-extended-query-record — 2026-10-05

**Zweck:** Beleg der Implementer-Rolle für Reviewer und Verifier: welche Änderung am geprüften Code jede Zusage des Slice rot macht, und dass sie einmal rot gesehen wurde.

**Vorgehen:** Jede Mutation einzeln im Arbeitsbaum auf genau eine Textstelle angewandt (E3 und E4: zwei Stellen gemeinsam), danach die Datei aus der Sicherung zurückgeschrieben. M-Zeilen: `go test -count=1` des betroffenen Pakets im gepinnten Go-Image des `Dockerfile` (dasselbe Image wie `make test`). E-Zeilen: `make test-integration`. Ein Kontrolllauf ohne Mutation (`make test`, `make test-integration`) war grün. „Rot“ heißt: der genannte Test schlägt fehl, kein Übersetzungsfehler.

**Stand:** gegen den Code des Commits, der diese Tabelle einführt.

## Unit (Record-Service, PGWire-Adapter, Upstream-Adapter)

| # | Zusage | Mutation | Ort | Rot in |
|---|---|---|---|---|
| M1 | Client-Nachrichten einer Gruppe gehen erst mit Flush oder Sync an den Upstream | jede Nachricht sofort einzeln senden | `services/record.go` (`ClientMessage`) | `TestRecordExtendedSyncGruppe`, `TestRecordExtendedFlushGruppen` |
| M2 | eine neue Gruppe nimmt ab ihrer ersten Client-Nachricht die Server-Nachrichten auf | `ziel` nicht weiterschalten | `services/record.go` (`ClientMessage`) | `TestRecordExtendedFlushGruppen` |
| M2c | eine Server-Nachricht nach der nächsten Client-Nachricht gehört zur folgenden Gruppe | `ziel` erst beim Senden der Gruppe weiterschalten | `services/record.go` (`ClientMessage`) | `TestRecordExtendedFlushGruppen` |
| M3 | jede Interaktion wird erst nach `Validate` übernommen | `Validate`-Ergebnis ignoriert | `services/record.go` (`uebernehmen`) | `TestRecordExtendedNichtDarstellbar` (ReadyForQuery ohne Sync, Extended-Antwort in einfacher Anfrage) |
| M4 | `Query` während einer laufenden Extended-Interaktion ist PGR-E6001 | Prüfung abgeschaltet | `services/record.go` (`Query`) | `TestRecordExtendedNichtDarstellbar` (Query in laufender Interaktion) |
| M5 | Client-Nachricht nach dem Sync vor dessen ReadyForQuery ist PGR-E1000 | Prüfung abgeschaltet | `services/record.go` (`ClientMessage`) | `TestRecordExtendedAufruferfehler` |
| M6 | Server-Nachricht ohne gesendete Gruppe ist PGR-E1000 | nur „keine Interaktion“ geprüft | `services/record.go` (`ServerMessage`) | `TestRecordExtendedAufruferfehler` |
| M7 | Extended-Interaktionen zählen in der Sequenz der Session mit | `Sequence` ohne `+ 1` | `services/record.go` (`ServerMessage`) | `TestRecordExtendedSyncGruppe`, `TestRecordExtendedMehrereInteraktionen` |
| M8b | nach einem Sync liest `sitzung` erst nach dessen ReadyForQuery weiter | Bedingung `!sync` beim Anfordern entfernt | `pgwire/server.go` (`sitzung`) | `TestExtendedSyncGruppe` |
| M9 | Verbindungsende in laufender Interaktion ist PGR-E4003 | Meldung abgeschaltet | `pgwire/server.go` (`sitzung`) | `TestExtendedAbbruch` (Verbindungsende nach Parse) |
| M10 | Terminate in laufender Interaktion ist PGR-E4003 | Meldung abgeschaltet | `pgwire/server.go` (`sitzung`) | `TestExtendedAbbruch` (Terminate nach Flush) |
| M11 | nicht zugestellte Server-Nachricht vor dem ReadyForQuery endet mit EndNormal, nicht EndLost | Herabstufung abgeschaltet | `pgwire/server.go` (`sitzung`) | `TestExtendedAbbruch` (frühere Server-Nachricht nicht zugestellt) |
| M12 | beim Herunterfahren endet eine Session erst nach dem ReadyForQuery der laufenden Interaktion | Prüfung am Schleifenkopf ohne `!offen` | `pgwire/server.go` (`sitzung`) | `TestExtendedHerunterfahren` (laufende Interaktion) |
| M12b | dasselbe | `ctx.Done()` auch bei laufender Interaktion im `select` | `pgwire/server.go` (`sitzung`) | `TestExtendedHerunterfahren` (laufende Interaktion) |
| M13 | Zielart außer `S`/`P` ist PGR-E6001 | unbekannte Zielart als Statement gelesen | `pgwire/server.go` (`zielart`) | `TestExtendedNichtUnterstuetzt` (Zielart), `TestToClientMessage` |
| M14b | eine leere Liste wird nil, wie der Leser sie liefert | leere Liste unverändert zurück | `pgwire/server.go` (`kopieOderNil`) | `TestToClientMessage` |
| M15 | Parameterwerte werden aus dem Puffer der Bibliothek kopiert | Wert ohne Kopie übernommen | `pgwire/server.go` (`toClientMessage`) | `TestToClientMessage` |
| M16 | ein leerer Parameter ist nicht NULL | leerer Wert als NULL gelesen | `pgwire/server.go` (`toClientMessage`) | `TestExtendedSyncGruppe` |
| M17 | je Client-Nachricht nur die Felder ihres Typs | `execute` belegt zusätzlich `Name` | `pgwire/server.go` (`toClientMessage`) | `TestExtendedSyncGruppe`, `TestToClientMessage` |
| M18 | `parameter_description` geht mit ihren Typ-OIDs an den Client | OIDs weggelassen | `pgwire/server.go` (`toMessage`) | `TestExtendedFlush`, `TestToMessageExtended` |
| M19 | ein leerer Parameter geht nicht als NULL an den Upstream | leerer Wert als NULL gesendet | `postgres/upstream.go` (`toFrontendMessage`) | `TestSendUndReceive` |
| M20 | `parameter_description` trägt ihre Typ-OIDs | OIDs weggelassen | `postgres/upstream.go` (`toResponse`) | `TestSendUndReceive` |
| M20b | `param_types` steht nur an `parameter_description` | `no_data` mit `ParamTypes` | `postgres/upstream.go` (`toResponse`) | `TestSendUndReceive` |
| M21 | `Receive` liest nicht über ReadyForQuery hinaus | Abbruch bei ReadyForQuery entfernt | `postgres/upstream.go` (`Receive`) | `TestReceiveEndetMitReadyForQuery` |
| M22 | Verbindungsende vor ReadyForQuery ist PGR-E4003 | Code PGR-E1000 | `postgres/upstream.go` (`Receive`) | `TestReceiveFehler` (Verbindungsende) |
| M23 | eine Client-Nachricht ohne gültige Zielart wird nicht gesendet | leere Zielart durchgelassen | `postgres/upstream.go` (`toFrontendMessage`) | `TestReceiveFehler` (Zielart) |

**Äquivalente Mutationen, verworfen:** M2b (`ziel` per `defer` am Ende von `ClientMessage` gesetzt) setzt es am selben Zeitpunkt wie der Code und blieb grün; ersetzt durch M2c. M14 (`kopieOderNil` prüft `s == nil` statt `len(s) == 0`) blieb grün, weil `append` auf nil mit null Elementen nil liefert; ersetzt durch M14b. M8 (Variable `sync` ungenutzt) war ein Übersetzungsfehler; ersetzt durch M8b.

## E2E (`make test-integration`)

| # | Zusage | Mutation | Rot in |
|---|---|---|---|
| E1 | Pipelining über mehrere Sync-Gruppen läuft gegen eine reale Instanz | `sitzung` liest nach einem Sync weiter (wie M8b) | `TestE2ERecordExtendedPipeline` (PGR-E1000 aus `ClientMessage`) |
| E2 | Flush erreicht den Server; der Client erhält die Antwort ohne Sync | `Send` lässt `Flush` weg | `TestE2ERecordExtendedPipeline`, `TestE2ERecordExtendedAbbruch` (Zeitüberschreitung beim Prepare) |
| E3 | die Gruppen stehen mit ihren Server-Nachrichten in Reihenfolge in der Aufzeichnung | `ziel` nicht weiterschalten und `Validate` abschalten | `TestE2ERecordExtendedPipeline` (Reihenfolgeprüfung) |
| E4 | die Aufzeichnung lädt mit dem Leser und besteht `Validate` (`laedt`) | ReadyForQuery doppelt in die Gruppe und `Validate` im Service abschalten | `TestE2ERecordExtendedPgx`, `TestE2ERecordExtendedPipeline` (replay startet nicht, Zeilen 92 und 184) |

## Ohne Mutation

- Race-Detector (`go test -race -count=30 ./internal/...` im gepinnten Go-Image mit `gcc`), kein Gate: grün. Er fand zwei Datenrennen in Test-Fakes (`fakeUpstream.letzte`, Hilfsfunktion `sende`), beide behoben; im Produktcode keines.
