# Mutationen: slice-extended-query-record — 2026-10-05

**Zweck:** Beleg der Implementer-Rolle für Reviewer und Verifier: welche Änderung am geprüften Code jede Zusage des Slice rot macht, und dass sie einmal rot gesehen wurde.

**Vorgehen:** Jede Mutation einzeln im Arbeitsbaum auf genau eine Textstelle angewandt (E3 und E4: zwei Stellen gemeinsam), danach die Datei aus der Sicherung zurückgeschrieben. M-Zeilen: `go test -count=1` des betroffenen Pakets im gepinnten Go-Image des `Dockerfile` (dasselbe Image wie `make test`). E-Zeilen: `make test-integration`. Ein Kontrolllauf ohne Mutation (`make test`, `make test-integration`) war grün. „Rot“ heißt: der genannte Test schlägt fehl, kein Übersetzungsfehler.

**Stand:** Erste Runde (M, E1 bis E4) gegen `5d666db`. Zweite Runde (Wiederholung, N, E5) gegen den Umbau nach [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) (Review F-301 bis F-307), den Commit, der sie einträgt. Mit dem Umbau entfallen `ServerMessage` und die Ereignisschleife `sitzung`; die Zeilen M5, M6, M8b bis M12b treffen Code, den es nicht mehr gibt, und ihre Zusagen prüfen jetzt N-Zeilen (Zuordnung N11/N12, Fehler des Aufrufers `TestRecordExtendedAufruferfehler`, Ende N7/N8/N13 bis N17, Herunterfahren N9/N10/N14).

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

## Zweite Runde (Umbau nach ADR-0030)

Wiederholt gegen den neuen Code, jeweils rot: M1 (zusätzlich `TestRecordGegendruck`, `TestRecordCloseBeendetWartende`), M2, M2c, M3, M4 (Zeitüberschreitung: die Query wartet auf die offene Interaktion, statt abgelehnt zu werden), M7 (in `AwaitServer`; `TestRecordExtendedSyncGruppe`, `TestRecordExtendedPipelining`), M13 bis M23 wie in der ersten Runde.

| # | Zusage | Mutation | Ort | Rot in |
|---|---|---|---|---|
| N1 | Empfang beginnt, sobald eine Gruppe an `Send` übergeben ist, nicht erst nach dessen Rückkehr (F-301) | `gesendet` erst nach `Send` setzen | `services/record.go` (`ClientMessage`) | `TestRecordGegendruck` |
| N1b | `ClientMessage` hält keine Sperre über `Send` | Sperre über `Send` gehalten | `services/record.go` (`ClientMessage`) | `TestRecordGegendruck`, `TestRecordCloseBeendetWartende` |
| N2 | `CloseSession` weckt jeden Wartenden | `Broadcast` beim Beenden entfernt | `services/record.go` (`CloseSession`) | Zeitüberschreitung in `TestRecordQueryWartetAufExtended` |
| N3 | `CloseSession` beendet ein am Upstream blockiertes `Send` | Upstream nicht geschlossen | `services/record.go` (`CloseSession`) | `TestRecordCloseBeendetWartende` |
| N4 | eine einfache Anfrage nach einem `Sync` wartet auf dessen Interaktion | Warten auf laufende Interaktionen entfernt | `services/record.go` (`Query`) | `TestRecordQueryWartetAufExtended`, `TestRecordCloseBeendetWartende` |
| N5 | … und auf die Zustellung ihres `ReadyForQuery` | Warten auf `empfaengt` entfernt | `services/record.go` (`Query`) | `TestRecordQueryWartetAufExtended` |
| N6 | eine laufende einfache Anfrage hält das Herunterfahren auf | `einfach` aus `laeuft` entfernt | `services/record.go` (`laeuft`) | `TestRecordHerunterfahren` |
| N7 | Verbindungsende oder Terminate während einer laufenden Interaktion ist PGR-E4003 (F-304: Entscheidung im Service) | Einstufung abgeschaltet | `services/record.go` (`CloseSession`) | `TestRecordExtendedEnde` |
| N8 | bei einem Schreibfehler entfällt die nicht zugestellte Interaktion | nichts verworfen | `services/record.go` (`CloseSession`) | `TestRecordExtendedEnde`, `TestRecordSessionVerbindungsende` |
| N9 | Herunterfahren wartet auf eine laufende Interaktion | `laeuft` nicht geprüft | `services/record.go` (`Shutdown`) | `TestRecordHerunterfahren` |
| N9b | … und auf die Zustellung des letzten `ReadyForQuery` | `unzugestellt` nicht geprüft | `services/record.go` (`Shutdown`) | `TestRecordHerunterfahren` |
| N10 | `Delivered` gibt das Ende nur beim Herunterfahren frei | `herunterfahren` nicht geprüft | `services/record.go` (`Delivered`) | `TestRecordHerunterfahren` |
| N11 | nach einem `Sync` beginnt die nächste Client-Nachricht eine neue Interaktion (Pipelining) | an die synchronisierte Interaktion angehängt | `services/record.go` (`ClientMessage`) | `TestRecordExtendedPipelining` |
| N12 | Server-Nachrichten gehören der ältesten laufenden Interaktion | der jüngsten zugeordnet | `services/record.go` (`AwaitServer`) | `TestRecordExtendedPipelining` |
| N13 | der Adapter meldet jede geschriebene Antwort als zugestellt | `Delivered` nicht gerufen | `pgwire/server.go` (`schreibe`) | `TestExtendedSyncGruppe`, `TestExtendedZweiRichtungen`, `TestExtendedHerunterfahren`, `TestQueryUndTerminate` |
| N14 | eine Lesefrist beim Herunterfahren beendet nur, wenn der Use Case es freigibt | sofort `EndShutdown` | `pgwire/server.go` (`clientRichtung`) | `TestExtendedHerunterfahren` (laufende Interaktion) |
| N15 | beide Richtungen laufen gleichzeitig (F-301) | Server-Richtung erst nach der Client-Richtung gestartet | `pgwire/server.go` (`recordSitzung`) | `TestExtendedZweiRichtungen`, `TestExtendedSyncGruppe`, `TestExtendedFlush`, `TestExtendedEreignisse`, `TestExtendedHerunterfahren` |
| N16 | endet die Session aus der Server-Richtung, endet auch das Lesen vom Client | Lesefrist beim Beenden nicht gesetzt | `pgwire/server.go` (`beende`) | `TestExtendedEreignisse` (Fehler aus AwaitServer) |
| N17 | ein Schreibfehler wird als `EndWriteFailed` gemeldet | als `EndClosed` gemeldet | `pgwire/server.go` (`endBeiSchreibfehler`) | `TestExtendedEreignisse`, `TestAntwortNichtZugestellt` |
| N18 | `Close` blockiert nicht hinter einem wartenden `Send` | `Lock` statt `TryLock` | `postgres/upstream.go` (`Close`) | `TestCloseBeendetWartende` |
| N19 | `Send` und `Receive` laufen gleichzeitig | `Receive` nimmt die Schreibsperre | `postgres/upstream.go` (`Receive`) | `TestSendUndReceiveGleichzeitig` |
| N20 | `Close` beendet wartendes `Send` und `Receive` | Verbindung nicht geschlossen | `postgres/upstream.go` (`Close`) | `TestCloseBeendetWartende` |

**Grün gebliebene Mutationen, Code bereinigt:** N5 und N6 in ihrer ersten Fassung (`unzugestellt > 0` im Warten von `Query`; `einfach` im Warten von `AwaitServer`) und `!empfaengt` in `Shutdown` blieben grün, weil andere Bedingungen sie abdecken: `empfaengt` hält `Query` bis zur Zustellung, und `AwaitServer` braucht eine laufende Extended-Interaktion, die nur dieselbe Richtung beginnt, die gerade in `Query` steht. Die drei Bedingungen sind entfernt; der Kommentar an `Query` nennt den Grund. N9b blieb zuerst grün; `TestRecordHerunterfahren` prüft seither `Shutdown` vor der Zustellung.

| # | Zusage | Mutation | Rot in |
|---|---|---|---|
| E5 | ein Batch mit großer Ausgabe und großem Parameter läuft durch; nach dem Schließen eines Clients mit blockierter Gruppe endet der Prozess auf SIGTERM (F-301) | wie N1 | `TestE2ERecordExtendedGegendruck` (Zeitüberschreitung), `TestE2ERecordExtendedSigtermNachBlockade` (endet nicht nach SIGTERM) |

**Gegenprobe gegen `5d666db`:** In einer Kopie von `5d666db` mit den neuen E2E-Tests sind `TestE2ERecordExtendedGegendruck` (Zeitüberschreitung nach 60 s) und `TestE2ERecordExtendedSigtermNachBlockade` („Recorder endet nicht nach SIGTERM“) rot; mit dem Umbau grün.

**Race-Detector** (`go test -race -count=30 ./internal/...`, kein Gate) nach dem Umbau: grün.
