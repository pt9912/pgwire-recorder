# Verifikation: slice-v1-abschluss-einspielen (Kern von `play`) — 2026-10-09

**Rolle:** Verifier (Modul 11). Geprüft wird, ob der Kern Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-09-review-slice-v1-abschluss-einspielen.md`, F-547 bis F-554). F-553 und F-554 hat er an mich übergeben, F-550 prüfe ich auf Auftrag mit.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen.md`, Kopf und §1 bis §8, gegen den Gesamt-Diff `e4963ee..cfcdcb5`. Architect vor dem Code: `f7c9c79`, `36b9d56`, `19f5512`, `e4963ee`. Implementer: `8025485`, Nacharbeit `cfcdcb5`. Review: `a980d4a`.

**Abgrenzung:** Anmeldung mit Passwort, TLS, die Optionen der Laufsteuerung und das Einspielen von Extended-Interaktionen sind abgegeben (§1, *Abgegeben*). Dass sie fehlen, ist kein Befund; geprüft ist der *Zwischenstand* aus §6.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-09

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- der Plan ganz, besonders §1 (*Abgegeben*, *Ausdrücklich NICHT*), §2 (DoD), §3, §4, §6 (*Randformen* mit *Zwischenstand*, *Risiken*) und §7 (*Belege des Implementers*, *Belege der Nacharbeit zum Review*, *Beobachtungen für Review und Closure*)
- `spec/spezifikation.md` `LH-FA-20.a` ganz (mit *Randformen je Schritt*), `LH-FA-14.a`, `LH-FA-17.a` *Wirkung einer URL*, `LH-FA-24.a` (Abgrenzung der Log-Zeile mit Zählern), `SPEC-038` *Warten in Tests*; [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-FA-14`](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler)
- [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md), [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md), Abschnitt *Entscheidung*
- am Stand `cfcdcb5`: `services/play.go`, `ports/driven/einspielziel.go`, `ports/driving/play.go`, `postgres/einspielen.go` ganz; Diff von `bootstrap.go`, `cli.go`, `upstream.go`, `verbindung.go`, `fehler.go`, Handbuch; `test/integration/play_e2e_test.go` ganz; die Hilfsfunktionen des Fake-Servers in `postgres/einspielen_test.go`
- die Nehmer `slice-v1-abschluss-einspielen-laufsteuerung` (§1, §2, §3, §4, §6) und `slice-v1-abschluss-antwortvergleich` (§1, DoD)
- `AGENTS.md` §3.9 bis §3.13
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen

Beim Start war der Arbeitsbaum sauber, HEAD `cfcdcb5`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.** Alles lief unter dem Scratch-Verzeichnis mit Präfix `ver-einspielen-`, nie im Repo.

- **Binary:** frische Kopie per `git archive cfcdcb5 | tar -x`, daraus Stufe `runtime` mit eigenem Tag. Ein internes Docker-Netz (`--internal`), ein Volume für Aufzeichnungen und Konfigurationsdateien, zwei PostgreSQL-Container aus dem gepinnten Image von `make test-integration` mit `log_connections`, `log_disconnections` und `log_statement=all` (Quelle für `record`, frische Zielinstanz für `play`), ein dritter mit Passwort (SCRAM). Für Fälle, die ein echter Server nicht liefert, ein Fake-Server aus der Standardbibliothek, gebaut im gepinnten Go-Image, der jedes Byte des Clients nach dem Startup protokolliert. Jeder `play`-Lauf mit `timeout 60` bzw. `docker wait` unter `timeout 40`. Client war `psql` aus dem PostgreSQL-Image.
- **Diskriminator für `Terminate`:** PostgreSQL protokolliert „unexpected EOF on client connection with an open transaction“, wenn ein Client in offener Transaktion ohne `Terminate` geht. Kalibriert mit `psql`, der in offener Transaktion per `SIGKILL` endet: Die Zeile erscheint. Die Probe-Aufzeichnungen öffnen deshalb vor dem geprüften Ende ein `BEGIN`.
- **Unit-Mutanten:** Image der Stufe `deps` aus der frischen Kopie; Grundlauf `gofmt -l` leer, `go vet`, `go test -count=1 ./internal/...` ohne Netz grün. Je Mutant eine neue Kopie von `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/` per `cp -r` ohne `-p`; ein Skript ersetzt genau eine Stelle und bricht ab, wenn der Suchtext nicht genau einmal trifft; `gofmt -l` auf der Datei leer; `go test -count=1 -timeout 120s -run <Test> <Paket>` im Image `deps` per Bind-Mount, `--network none`; danach nur das Verzeichnis des Mutanten gelöscht.
- **Danach:** Container, Netz, Volume, die beiden eigenen Images und alle Kopien entfernt; `docker ps -a`, `docker network ls`, `docker volume ls` und `docker images` führen nichts mit dem Präfix. `make gates` lief im Arbeitsbaum.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Einspielen einfacher Anfragen, Optionen, Vorrang, U8, Zwischenstand, Handbuch §5 ([`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung)). **Bestätigt; Beleg-Form siehe V-133.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Abnahmeszenario 12: DDL und DML als einfache Anfragen gegen eine leere Instanz ohne Passwort und TLS, Wirkung in der Datenbank | Binary: `record` mit `psql` über den Recorder, zwei Sessions (`CREATE TABLE kunden`, `INSERT` zwei Zeilen, `UPDATE`; `CREATE TABLE bestellungen` mit Fremdschlüssel, `INSERT`, `DELETE`). `play` gegen die frische Zielinstanz: Exit 0, Zeilen `info` Start und Ende, danach `kunden` = `1|anna`, `bestellungen` = `10|1`. Der Server protokolliert die sechs Anweisungen in Reihenfolge. `TestE2EPlayDDLDML` im eigenen `make gates` grün | bestätigt |
| mehrere Sessions über eigene Verbindungen nacheinander | Binary: zwei `connection received` mit verschiedenen Client-Ports, die zweite erst nach `disconnection` der ersten; `application_name: psql` aus der Aufzeichnung gesendet | bestätigt |
| Optionen aus Kommandozeile, Umgebung, Abschnitt `play:` | Binary: Datei mit `log_level: warn` und `play: {upstream: mituser, input: …, database: optdb}`, `play --config …` ohne Option: Exit 0, keine Zeile `info`, Server sieht `user=urluser database=optdb`. Mit `PGWIRE_RECORDER_DATABASE=aufgezdb` und `PGWIRE_RECORDER_LOG_LEVEL=info` daneben: `database=aufgezdb`, Zeilen `info` erscheinen | bestätigt |
| `--user`, `--database` vor URL vor Aufzeichnung | Binary, Server-Log je Lauf (`connection authorized`): URL `urluser@…/urldb` ohne Option → `urluser/urldb`; mit `--user optuser` → `optuser/urldb`; mit `--database optdb` → `urluser/optdb`; `host:port` → `aufgez/aufgezdb` aus der Aufzeichnung; `PGWIRE_RECORDER_USER`/`_DATABASE` bei `host:port` → `optuser/optdb`; URL ohne Benutzer → `aufgez/urldb`; `PGWIRE_RECORDER_USER` vor URL → `optuser/urldb`. Startzeile nennt nur `upstream=host:port` und `input`. Mutanten M01, M06, M12 rot (Abschnitt 5) | bestätigt |
| U8: nicht gesetzte Variable in Benutzer, Passwort, Datenbank `PGR-E2005` | Binary: `${VER_PW}` im Passwort nicht gesetzt → `PGR-E2005 … connections.mitpw: Umgebungsvariable VER_PW …`, Exit 2. `${VER_U}` im Benutzer nicht gesetzt, obwohl `--user` ihn überschreibt → `PGR-E2005` mit `VER_U` (erste in der Reihenfolge der URL). Gesetztes Passwort gegen einen Server ohne Passwort: Exit 0, Wert in keiner Ausgabe | bestätigt |
| gewöhnliches Argument `PGR-E2001` | Binary: `play … extra` und `play … -- extra` je `PGR-E2001` „unerwartetes Argument“, Exit 2 | bestätigt |
| `--fail-on-unconsumed` unbekannt, Variable unbeachtet auch ungültig | Binary: `--fail-on-unconsumed` → `PGR-E2001`; `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=vielleicht` (zusammen mit ungültigen Werten in `…_UPSTREAM_TLS` und `…_CONTINUE_ON_ERROR`) → Exit 0, eingespielt | bestätigt |
| `--log-level` mit Wertemenge, Strenge, Schwelle | Binary: `--log-level INFO` → `PGR-E2001`; `PGWIRE_RECORDER_LOG_LEVEL=warning` neben gültigem `--log-level error` → `PGR-E2001` an der Umgebungsvariable; `--log-level error` → keine Zeile `info`; Standard `info` zeigt Start und Ende | bestätigt |
| Zwischenstand: `--upstream-tls`, `--upstream-ca`, Optionen der Laufsteuerung unbekannt; Schlüssel in `play:` unbekannt | Binary: `--upstream-tls`, `--upstream-ca`, `--continue-on-error` je `PGR-E2001`; `continue_on_error: true` in `play:` → `PGR-E2004 … play.continue_on_error: unbekannter Schlüssel` | bestätigt |
| Zwischenstand: `sslmode=require` `PGR-E2004` | Binary: „`connections.tls: sslmode=require ist bei play ungültig, play verbindet ohne TLS zum Upstream`“, Exit 2 | bestätigt |
| Zwischenstand: Extended-Interaktion Startfehler `PGR-E6001` ohne Verbindung | Binary: Aufzeichnung über den Recorder mit `psql` (`\bind 4711`), Session 1 einfach (`CREATE TABLE ext_vorher`), Session 2 Extended. `play`: Exit 6, genau die Zeile beim Prozessende „`PGR-E6001: Session 2, Interaktion 1: eine Interaktion der Art extended spielt play nicht ein`“, keine Log-Zeile, kein `connection received` an der Zielinstanz, `ext_vorher` fehlt dort; `4711` in keiner Ausgabe. Mutant M16 rot | bestätigt |
| Zwischenstand: Passwort-Anforderung `PGR-E4005` ohne gesendete Bytes | Fake-Server, Code 3, 5, 10 und 7, jeweils mit Passwort aus dem Platzhalter und `PGWIRE_RECORDER_PASSWORD` gesetzt: `PGR-E4005 … (Code n), das play nicht unterstützt`, Exit 4; der Fake zählt nach dem Startup **0 Bytes** bis EOF. Gegen echtes PostgreSQL mit SCRAM: `PGR-E4005 (Code 10)`, der Server protokolliert nur `connection received`. Kein Passwort in `stdout` oder `stderr` | bestätigt |
| Aufzeichnung ohne Session mit Interaktion | Binary: `sessions: []` → Exit 0, kein `connection received`. (Eine Session mit `interactions: []` lehnt der Leser als beschädigt ab, `PGR-E3003`, Bestand) | bestätigt |
| Handbuch §5 *Konfigurationsdatei*: `play:` in Beispiel und Liste, Beispiel als Datei startet `config show` ohne Meldung (V-125) | Beispiel wörtlich als `.pgwire-recorder.yaml`, `config show` im Produkt-Image ohne Netz: Exit 0, `stderr` 0 Bytes, `stdout` = Pfad und Datei. `play` mit der Datei: `PGR-E3001` (Aufzeichnung fehlt), also nach dem Laden. Der neue Absatz zur Wirkung bei `play` deckt sich mit den Vorrang-Läufen oben | bestätigt |

### Punkt 2: Fehlerantwort, Verbindungsfehler, Aufbau-Fehler, Nachrichten und Anmelde-Codes im Aufbau (Test). **Bestätigt; V-130.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Fehlerantwort: `PGR-E4004`, Exit 4, keine weitere Nachricht, `Terminate` | Binary: `BEGIN`, `INSERT`, `SELECT * FROM gibtsnicht`, `CREATE TABLE nach_fehler`, Session 2 `CREATE TABLE zweite_session`: Exit 4, eine Zeile `error` `code=PGR-E4004 … Session 1, Interaktion 3: Fehlerantwort des Servers 42P01 „relation "gibtsnicht" does not exist“`. Server: nur die drei Anweisungen, keine weitere Verbindung, `disconnection` **ohne** „unexpected EOF … open transaction“ → `Terminate` gesendet; die `INSERT`-Zeile ist zurückgesetzt, beide Tabellen fehlen. Mutanten M17, M21 rot | bestätigt |
| `FATAL` sofort `PGR-E4003` | Binary: `SELECT pg_terminate_backend(pg_backend_pid())`, danach `CREATE TABLE nach_fatal`: Exit 4, `PGR-E4003 … 57P01 „terminating connection due to administrator command“, die Verbindung endet`; zweite Anweisung nie beim Server. Mutanten M04, M05, M20 rot | bestätigt |
| Verbindungsverlust ohne Fehlerantwort `PGR-E4003` | Binary: Server per `SIGKILL` während `pg_sleep(20)`: Exit 4, `PGR-E4003 … Interaktion 2: Verbindung zum Server beendet: unexpected EOF` | bestätigt |
| `Copy…Response` `PGR-E6001` | Binary: `COPY kunden TO STDOUT`, danach `CREATE TABLE nach_copy`: Exit 6, `PGR-E6001`, zweite Anweisung nicht ausgeführt | bestätigt |
| nicht erreichbar `PGR-E4002` | Binary: geschlossener Port und unbekannter Name je `PGR-E4002 … nicht erreichbar`, Exit 4 | bestätigt |
| Fehlerantwort im Aufbau: außer Klasse 28 `PGR-E4002`, Klasse 28 `PGR-E4005`, mit SQLSTATE und `M` | Binary: `--database gibtsnichtdb` → `PGR-E4002 … Fehlerantwort im Aufbau 3D000 „database "gibtsnichtdb" does not exist“`; `--user niemand` → `PGR-E4005 … 28000 „role "niemand" does not exist“`; je Exit 4. Mutanten M02, M09 rot | bestätigt |
| Nachrichten und Anmelde-Codes im Aufbau nach §6 | Unit, gegengeprüft mit Mutanten M03 (`A` nicht verworfen), M18 (`ReadyForQuery` vor `AuthenticationOk`), M19 (Code 12 keine Fortsetzung), X03 (Länge unter 4): alle rot am genannten Teilfall | bestätigt |
| Abbruch im Aufbau ohne `Terminate`, Verbindung geschlossen | Binary/Fake: nach Passwort-Anforderung 0 Bytes bis EOF. Das **Schließen** selbst hält keine Zusicherung eines Tests: V-130 | bestätigt, V-130 |
| Beleg in §7 je Zusage Zusage · Mutation · roter Test | beide Tabellen gelesen; 21 Stichproben nachgefahren, alle rot aus dem in §7 genannten Grund (Abschnitt 5) | bestätigt |

### Punkt 3: Abbruchsignal, zweites Signal, Exit-Code 0, Verbindungsende nach dem ersten `ReadyForQuery` (Test). **Bestätigt; V-131.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| erstes Signal: Ende nach der laufenden Interaktion, keine weitere Session | Binary: `BEGIN`, `SELECT pg_sleep(3)`, `CREATE TABLE sig_nach`, Session 2; `SIGINT` bei 1,5 s: `pg_sleep` läuft zu Ende, Exit 0 nach 3,5 s, Zeilen `info` „Abbruchsignal …“ und „play beendet“, weder `sig_nach` noch Session 2; `disconnection` ohne „unexpected EOF“ → `Terminate`. Mutanten M14, M21, X01 rot | bestätigt |
| erstes Signal im Aufbau: Aufbau läuft zu Ende, keine Interaktion | Fake verzögert `AuthenticationOk` 3 s, `SIGINT` bei 1 s: Exit 0 nach 3,5 s, der Fake erhält nach dem Aufbau genau 5 Bytes Typ `X` (`Terminate`), keine `Query`. Kontrolle: Fake ohne Antwort, nur ein `SIGINT` → `play` wartet bis zum zweiten Signal | bestätigt |
| zweites Signal: sofortiges Ende, Exit 0, unterbrochene Interaktion kein Fehler, weitere Signale ohne Wirkung | Binary: `pg_sleep(20)` in offener Transaktion, `SIGINT` 1,5 s, `SIGTERM` 2,5 s, `SIGINT` 3,0 s: Ende unmittelbar nach dem zweiten Signal, Exit 0, keine Zeile `error`. Der Server beendet die Session nach `pg_sleep` **ohne** „unexpected EOF“ → `Terminate` angenommen. Mutante X06 rot | bestätigt |
| zweites Signal im Aufbau: kein `Terminate`, Verbindung geschlossen | Fake ohne Antwort, Signale bei 1 s und 2 s: Exit 0 nach 2,6 s, Fake zählt 0 Bytes bis EOF | bestätigt |
| zweites Signal im Verbindungsversuch | Binary gegen eine unbelegte Adresse im internen Netz: ohne Signal `PGR-E4002 no route to host` nach 3 s; mit Signalen bei 0,8 s und 1,3 s Exit 0 nach 1,9 s, keine Zeile `error`. Mutante M08 rot | bestätigt |
| `Terminate` beim blockierten Senden, Schreibfrist | Mutanten M07, M10 rot. Die Sperre in `Anfrage` hält kein Test: V-131 | bestätigt, V-131 |
| jedes Verbindungsende nach dem ersten `ReadyForQuery` ohne Vergleich `PGR-E4003`; Rangfolge mit Vergleich bei `slice-v1-abschluss-antwortvergleich` | `FATAL` und Verbindungsverlust oben. Der Nehmer führt die Rangfolge in seiner DoD (Punkt 3: „ein Abbruchsignal endet nach einer Abweichung mit 5“); die Adresse nimmt an | bestätigt |

### Punkt 4: `make gates` grün. **Bestätigt durch eigenen Lauf; Beleg in §7 fehlt (V-132).**

Eigener Lauf im Arbeitsbaum am Stand `cfcdcb5`: Exit 0. Darin `baseline-verify: v6.16.0 OK`, `d-check: 425 Datei(en) geprüft, 0 Befund(e)`, a-check `gesamt: 0 Befund(e)`, `a-check-negativ: gruen`, `run-integration-tests: gruen` mit `TestE2EPlayDDLDML`, `TestE2EPlayFehlerantwort`, `TestE2EPlayAufbau`, `TestE2EPlayKonfiguration` je `PASS`, `abdeckung-gegenprobe`, `kopf-check-gegenprobe`, `lint-gegenprobe`, `commit-msg-gegenprobe` je `gruen`. Die Stufe `test` kam dort aus dem Build-Cache; unabhängig davon lief mein Grundlauf `gofmt -l`, `go vet`, `go test -count=1 ./internal/...` ohne Cache und ohne Netz grün.

### Übrige DoD-Punkte (konstant je Slice)

- **Review:** liegt vor (`a980d4a`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen:** offen, wie es die Reihenfolge will. Für die Closure aus meiner Sicht: Das zweite Risiko (`FATAL` erst im Test belegbar) ist jetzt auch am Binary belegt (57P01 → `PGR-E4003`, oben). Die Rückführung aus §4 trat nicht ein (F-554). Die Klassen aus der Summary des Reviews (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/folge-slice-adresse-nimmt-nicht-an`) gehören in die Register-Fortschreibung; V-129 bis V-131 sind weitere Ausprägungen derselben drei.

---

## 2. Übergaben des Reviews

| Befund | Prüfung | Urteil |
|---|---|---|
| F-547 zweites Signal beim blockierten Senden | `TestEinspielSchliesseBeimSenden` und `TestEinspielSchliesseOhneAnnahme` neu; M07 (`Lock` statt `TryLock`) und M10 (ohne Schreibfrist) rot mit den Meldungen aus §7. Rest: V-131 | erledigt; V-131 |
| F-548 Abbruch des Verbindungsversuchs | `TestEinspielVersuchAbgebrochen`; M08 rot; am Binary bestätigt (Punkt 3) | erledigt |
| F-549 SQLSTATE und `M` im Aufbau | `TestEinspielAufbauFehlerMeldung` und `TestE2EPlayAufbau` prüfen beides; M09 rot; am Binary bestätigt | erledigt |
| F-550 Form der Übergabe an den Play-Service | Der Bootstrap konvertiert `cli.Einspielvorgaben` in `services.PlayOptions`, ohne ein Feld zu nennen; M06 rot. Der Nehmer `slice-v1-abschluss-einspielen-laufsteuerung` §3 nennt die Felder in beiden Typen. Für die **Optionen** hält sein Ausschluss „kein Code im Bootstrap“ jetzt. Der vom Implementer gemeldete offene Punkt zum **Exit-Code** steht dagegen nur in §7 des Kerns: V-129 | Optionen erledigt; Exit-Code V-129 |
| F-551 Handbuch, erster Teil von F-533 | §5 *Einstellungen* und *Konfigurationsdatei* sagen jetzt Verbindung vor Aufzeichnung; sieben Vorrang-Läufe am Binary decken sich mit dem Text. Plan §1 führt den ersten Teil von F-533 als geliefert | erledigt |
| F-552 Schweregrad in der Meldung | Text ohne Schweregrad; M20 (Schweregrad zurück) rot; am Binary: `FATAL`-Meldung nennt SQLSTATE und `M` | erledigt |
| F-553 Handbuch §4 im Zielstand; Zeile `PGR-E4004` | Die Zeile `PGR-E4004` in §7 nennt jetzt Sitzung, Nummer, SQLSTATE und Meldung, nicht die Anfrage; das deckt sich mit `LH-FA-20.a` *Meldungen* und dem Binary. §4 (*Eine Aufzeichnung in eine Datenbank einspielen*) und die Optionstabelle in §5 beschreiben weiter den Zielstand mit Optionen, die heute `PGR-E2001` sind (am Binary gesehen). Die DoD verlangt nur §5 *Konfigurationsdatei*; jede dieser Optionen hat einen Nehmer in `welle-v1-abschluss`. Kein Befund gegen diesen Slice; das Urteil über einen Hinweis bis dahin bleibt beim Planner, wie §7 sagt | eingeordnet |
| F-554 Größe | §7 *Größe* misst jetzt gegen §6 und nennt keine Zahl ohne Artefakt; +2806 −57 im Gesamt-Diff, davon rund 630 Produktzeilen; die Rückführung aus §4 trat nicht ein | erledigt |

---

## 3. Hexagon und Architektur-Gate

- `make a-check` im eigenen `make gates`: `gesamt: 0 Befund(e)`, `a-check-negativ: gruen`.
- `services/play.go` importiert `context`, `errors`, `fmt`, `model`, `ports/driven`; die Ports nur `context` und `model`. `pgproto3` liegt nur in `driven/postgres`. `record` nutzt unverändert `postgres.Upstream`; `upstream.go` des Postgres-Adapters ist nicht im Diff.
- [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md): Sessions nacheinander in der Reihenfolge der Datei, erste Fehlerantwort bricht ab — am Binary gesehen. [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md): im Zwischenstand jedes Verfahren `PGR-E4005` ohne Senden, übriger Aufbau `PGR-E4002` — am Binary und am Fake gesehen.

---

## 4. Plan gegen Code

- **§1:** Geliefert ist der Kern: einfache Anfragen, Abbruch bei jedem Fehler, Signal ohne Option, Optionen aus drei Quellen, Wirkung der Verbindung bei `play` (Host, Port, Benutzer, Datenbank), U8, Handbuch §5. Abgegebenes fehlt, wie §1 sagt; der Zwischenstand hält am Binary.
- **§3:** Jede geänderte Datei steht dort, dazu `postgres/export_test.go` als Test-Hilfe. Die Zeile Bootstrap („ohne Feld zu nennen“, Zeile `info` ohne Zeitpunkt) folgt dem Code der Nacharbeit.
- **§4:** keine Rückführung eingetreten.
- **§6:** Jede Randform, die ich am Binary oder Fake gefahren habe, folgt ihrem Ort in `LH-FA-20.a` bzw. `LH-FA-17.a`: Optionen, Variablen aller Teile, Start (Ladefehler, ohne Session, Extended vor jeder Verbindung, Startfehler als Zeile beim Prozessende), Startup-Daten, Aufbau (nicht erreichbar, Klasse 28 und andere, Verfahren ohne Senden), Interaktion (Copy, `FATAL`, Verbindungsende, kein Weiterlesen nach `PGR-E4004`), Ende einer Session mit `Terminate`, alle Signal-Fälle. Keine Randform außerhalb von §6 ist im Code entschieden.
- **§7:** Beide Tabellen decken die Zusagen in der Form *Zusage · Mutation · roter Test*; jede nachgefahrene Zeile deckt sich. Es fehlt der Beleg für `make gates` (V-132). Der Hinweis an die Laufsteuerung (dritter Punkt der *Beobachtungen*) steht nur hier (V-129).

---

## 5. Gefahrene Mutanten

Alle am Stand `cfcdcb5`, je in einer frischen Kopie, `go test -count=1 -timeout 120s -run <Test> <Paket>`. M01 bis M21 sind Stichproben gegen Zeilen der Tabellen in §7, X01 bis X09 eigene Stellen.

| Mutant | Änderung | Ergebnis |
|---|---|---|
| basis | keine | grün (`./internal/...`) |
| M01 | `services`: `--user` ersetzt `user` nicht | rot: `TestPlayStartup` („erwartet … user:neu“) |
| M02 | `postgres`: Klasse 28 als Präfix `2` | rot: `TestEinspielAufbauFehler/Klasse_2_ähnlich` |
| M03 | `postgres`: `A` nicht in der Liste des Aufbaus | rot: `TestEinspielAufbauVerworfen` („Nachricht 'A' ist im Aufbau nicht vorgesehen“) |
| M04 | `services`: nur `FATAL`, nicht `PANIC` | rot: `TestPlayFehlerantwort/V_PANIC`, `/S_PANIC_ohne_V` |
| M05 | `services`: Schweregrad aus `S` statt `V` | rot: `TestPlayFehlerantwort/V_ERROR_vor_S_FATAL` u. a. |
| M06 | `bootstrap`: `services.PlayOptions{}` | rot: `TestRunPlayOptionen` („Startup map[user:u]“) |
| M07 | `postgres`: `Lock` statt `TryLock` in `Schliesse` | rot: `TestEinspielSchliesseBeimSenden` („Schliesse endet binnen 5 s nicht“) |
| M08 | `postgres`: Wählen mit `context.Background()` | rot: `TestEinspielVersuchAbgebrochen` |
| M09 | `postgres`: SQLSTATE aus der Meldung im Aufbau | rot: `TestEinspielAufbauFehlerMeldung` (28P01, 28000, 3D000) |
| M10 | `postgres`: ohne Schreibfrist in `Schliesse` | rot: `TestEinspielSchliesseOhneAnnahme` |
| M11 | `cli`: `play` nicht in `leserKommandos` | rot: `TestLeserAlleOptionen`, `TestDateiUngueltig` („play: unbekannter Schlüssel“) |
| M12 | `cli`: Benutzer der URL überschreibt `--user` | rot: `TestPlayVerbindung` |
| M13 | `cli`: `sslmode=require` bei `play` nicht geprüft | rot: `TestPlayVariablen` |
| M14 | `services`: `abbruch` aus `ctx` statt `WithoutCancel(ctx)` | rot: `TestPlayErstesSignal`, `TestPlayErstesSignalImAufbau` („ctx des Aufbaus beendet: true“) |
| M15 | `bootstrap`: Zeile `info` zum ersten Signal entfernt | rot: `TestRunPlaySignalImStart` |
| M16 | `services`: Prüfung auf Extended wirkungslos | rot: `TestPlayStart` („Extended: <nil>“) |
| M17 | `services`: Fehlerantwort ohne Schweregrad liefert keinen Fehler | rot: `TestPlayFehlerantwort/ERROR`, `/ohne_Schweregrad` |
| M18 | `postgres`: `ReadyForQuery` vor `AuthenticationOk` angenommen | rot: `TestEinspielAufbauFehler/ReadyForQuery_vor_AuthenticationOk` |
| M19 | `postgres`: Code 12 keine Fortsetzung | rot: `TestEinspielAufbauFehler/SASLFinal` |
| M20 | `services`: Schweregrad zurück in die Meldung (F-552) | rot: `TestPlayFehlerantwort/V_FATAL` u. a. |
| M21 | `services`: Prüfung des Signals je Interaktion entfernt | rot: `TestPlayErstesSignal` („anfrage B“) |
| X01 | `services`: Prüfung des Signals je Session entfernt | rot: `TestPlayErstesSignal`, `TestPlayErstesSignalImAufbau` (zweites `verbinde`) |
| X02 | `postgres`: im Fehlerzweig von `Verbinde` kein `Close` | **rot nur über die Gesamtfrist** von `go test` (2 min); jeder Teilfall einzeln **grün nach 20 s** — V-130 |
| X03 | `postgres`: Länge unter 4 angenommen | rot: `TestEinspielAufbauFehler/Länge_unter_4` |
| X04 | `postgres`: `Anfrage` ohne die Sperre `schreiben` | **grün** — V-131 |
| X05 | `bootstrap`: Zeile zum Signal auch ohne Signal | rot: `TestRunPlay`, `TestRunPlayFehlerantwort` |
| X06 | `services`: unterbrochene Session beim zweiten Signal als Fehler | rot: `TestPlayZweitesSignal` („PGR-E4003 … geschlossen“) |
| X07 | `postgres`: Abbruch des Aufbaus durch `ctx` nicht gesondert gemeldet | grün; äquivalent: Der Aufbau scheitert dann am geschlossenen Socket mit `PGR-E4002`, und der Play-Service wertet jeden Fehler nach dem zweiten Signal nicht. Im Rennfall (Aufbau fertig, als `ctx` endet) schließt die `AfterFunc` der Session sofort, und `ctx` des ersten Signals ist dann ebenfalls beendet; beobachtbar ist kein Unterschied |
| X08 | `bootstrap`: Startzeile ohne `input` | rot: `TestRunPlay` |
| X09 | `services`: Fehler beim Senden übergangen | rot: `TestPlayFehler` („Senden: <nil>“) |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-129 | MEDIUM | Der offene Punkt des Implementers für `slice-v1-abschluss-einspielen-laufsteuerung` steht nur in §7 des Kerns, nicht im Nehmer. Der Bootstrap bildet den Exit-Code aus der **ersten** Meldung und schreibt die Zeilen `error` erst nach dem Einspielen; `LH-FA-20.a` *Exit-Code* verlangt nach einem früheren `PGR-E4004` den Code des **abbrechenden** Fehlers (etwa 6). Ohne Code im Bootstrap geht das nur, wenn der Play-Service die Meldungen außer der zeitlichen Reihenfolge liefert; Reihenfolge und Zeitpunkt der Zeilen `error` legt weder `LH-FA-20.a` *Meldungen* noch `LH-FA-14.a` fest. Der Nehmer prüft in §4 nur „die Form, in der der Kern Optionen und Signale an den Play-Service reicht“, sein Risiko in §6 nennt ebenso nur Optionen und Signale, und sein §1 setzt voraus, dass der Bootstrap den Exit-Code aus dem Ergebnis des Play-Service bildet. *Failure-Szenario:* Der Nehmer startet ohne die Frage, liefert `--continue-on-error` mit Exit-Code 4 statt 6 nach `PGR-E4004` und `PGR-E6001`, oder er trifft beim Code seinen eigenen Ausschluss und die Rückführung in §4. | `AGENTS.md` §3.13, §3.12; `LH-FA-20.a` *Exit-Code*, *Meldungen*; Plan §7 *Beobachtungen für Review und Closure* | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen.md` §7 · „Hinweis zur Prüfung des Architect vor dessen Code, dort §4“; `internal/bootstrap/bootstrap.go` · `if code == "" {` | nein (Lesen des Nehmers §1, §4, §6 gegen §7 des Kerns) | Sendung an einen Nehmer nur beim Geber notiert |
| V-130 | LOW | Dass ein Fehler im Aufbau die Verbindung schließt, hält keine Zusicherung: Mutant X02 (kein `Close` im Fehlerzweig von `Verbinde`) lässt jeden Teilfall von `TestEinspielAufbauFehler` nach 20 s **grün**, weil der Fake-Server nach seiner Frist von 20 s `io.ReadAll` beendet und den Lauf ohne Bytes als Erfolg meldet. Rot wird der Paketlauf erst an der Zeitgrenze von `go test`. Im Produkt endet der Prozess kurz danach und schließt den Socket; die Zusage im Port und in der Abdeckung des Tests („schließt die Verbindung“) ist dennoch ungeprüft. | `LH-FA-20.a` *Abbruch im Aufbau*; `SPEC-038` *Warten in Tests* („die Zeitgrenze von `go test` ist keine Frist“); `AGENTS.md` §3.10, §3.11 | `internal/adapters/driven/postgres/einspielen_test.go` · `_ = conn.SetDeadline(time.Now().Add(20 * time.Second))` und `rest, _ := io.ReadAll(conn)`; `internal/hexagon/ports/driven/einspielziel.go` · „die Verbindung ist dann geschlossen“ | ja (Mutant X02, `go test -run 'TestEinspielAufbauFehler$/Klasse_28_FATAL$'` grün nach 20 s) | Frist im Test-Double macht ein ausbleibendes Ereignis grün |
| V-131 | LOW | Dass `Schliesse` und `Anfrage` nicht gleichzeitig schreiben, hält kein Test: Mutant X04 (`Anfrage` ohne die Sperre `schreiben`) bleibt in allen Tests von `postgres` grün. Ohne die Sperre schreibt `Schliesse` beim zweiten Signal in denselben Puffer des `pgproto3.Frontend`, den `Anfrage` gerade sendet; das ist ein Datenrennen, und die Tests laufen ohne `-race`. M07 hält nur die Hälfte in `Schliesse` (`TryLock`). | `LH-FA-20.a` *Abbruchsignal*; `AGENTS.md` §3.10, §3.11 | `internal/adapters/driven/postgres/einspielen.go` · „schreiben hält das Senden von Anfrage und das Terminate von Schliesse auseinander“ | ja (Mutant X04, `go test ./internal/adapters/driven/postgres/` grün) | Zusage nur für einen Teil ihrer Fälle von einer Mutation gehalten |
| V-132 | LOW | §7 belegt `make gates` nicht: beide Abschnitte *Läufe* sagen „`make gates` vor der Übergabe (Ergebnis im Bericht)“. Der Bericht des Implementers ist kein Artefakt; der Beleg, den die DoD verlangt, fehlt damit in §7. Mein eigener Lauf ist grün (Punkt 4). | Auftrag des Verifiers (Belege in §7, seit slice-harness-blackbox-kern); Plan §2 DoD „`make gates` grün“ | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen.md` §7 · „`make gates` vor der Übergabe (Ergebnis im Bericht)“ | ja (§7 lesen) | DoD-Beleg nur im Bericht, nicht in §7 |
| V-133 | LOW | DoD-Punkt 1 schließt mit „(Integrationstest)“. Die Integrationstests decken Abnahmeszenario 12, Fehlerantwort, Aufbau und Konfiguration mit `play:` und `--fail-on-unconsumed`. Vorrang von `--user` und `--database` vor der URL, U8, `--log-level`, die Umgebungsvariable von `--fail-on-unconsumed` und der ganze Zwischenstand (Extended `PGR-E6001` ohne Verbindung, `--upstream-tls`, `sslmode=require`) sind nur in Unit-Tests von CLI und Play-Service geprüft. Am Binary halten alle (Punkt 1); die Beleg-Form weicht vom Wortlaut der DoD ab. | Plan §2 DoD-Punkt 1 | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen.md` §2 · „bis zu den Folge-Slices gilt der Zwischenstand aus §6 (Integrationstest)“ | ja (`grep` in `test/integration/play_e2e_test.go`) | Beleg-Form schwächer als in der DoD zugesagt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Passwörter und Parameterwerte in Ausgaben | geprüft, ohne Befund: In keinem Lauf erscheinen die gesetzten Passwörter (`VER_PW`, `PGWIRE_RECORDER_PASSWORD`) oder der Parameterwert `4711` der Extended-Aufzeichnung in `stdout` oder `stderr`; Startzeile ohne Benutzer und Datenbank |
| Exit-Codes je Klasse | geprüft, ohne Befund: 0, 2 (`PGR-E2001`, `PGR-E2004`, `PGR-E2005`), 3 (`PGR-E3001`, `PGR-E3003`), 4 (`PGR-E4002`, `PGR-E4003`, `PGR-E4004`, `PGR-E4005`), 6 (`PGR-E6001` beim Start und in der Interaktion) am Binary |
| Log-Zeilen nach `LH-FA-14.a` | geprüft, ohne Befund: `info` Start, erstes Signal, Ende; `error` je Fehler mit `code` und `error`; Startfehler nur als Zeile beim Prozessende |
| `record` unverändert | geprüft, ohne Befund: Aufzeichnungen dieses Laufs mit `record` aus derselben Binary-Kopie; `postgres/upstream.go` nicht im Diff |
| Mutant X07 | äquivalent, siehe Abschnitt 5 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 0 |

- V-129: Der offene Punkt zum Exit-Code nach einem früheren `PGR-E4004` steht nur beim Geber; der Nehmer prüft ihn nicht.
- V-130: Das Schließen der Verbindung nach einem Fehler im Aufbau hält keine Zusicherung; der Fake-Server macht es nach 20 s grün.
- V-131: Die Sperre in `Anfrage` hält kein Test (X04 grün).
- V-132: §7 belegt `make gates` nicht.
- V-133: DoD-Punkt 1 sagt „Integrationstest“, mehrere Teile sind nur in Unit-Tests geprüft.

**Finding-Klassen dieses Laufs:** Sendung an einen Nehmer nur beim Geber notiert · Frist im Test-Double macht ein ausbleibendes Ereignis grün · Zusage nur für einen Teil ihrer Fälle von einer Mutation gehalten · DoD-Beleg nur im Bericht, nicht in §7 · Beleg-Form schwächer als in der DoD zugesagt

Zuordnung für die Closure: V-129 ist eine Ausprägung von `BEO-REPO/folge-slice-adresse-nimmt-nicht-an`, V-130 und V-131 von `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.

## Verdikt

**Urteil je DoD-Liefer-Punkt:** Punkt 1 bestätigt (Beleg-Form V-133), Punkt 2 bestätigt (V-130), Punkt 3 bestätigt (V-131), Punkt 4 durch eigenen Lauf bestätigt (Beleg in §7 fehlt, V-132). Das gelieferte Verhalten hält am echten Binary gegen PostgreSQL und den Fake-Server in jedem gefahrenen Fall; 21 Stichproben gegen §7 und 6 von 9 eigenen Mutanten sind am genannten Test rot aus dem richtigen Grund; X02 wird nur über die Gesamtfrist rot (V-130), X04 bleibt grün (V-131), X07 ist äquivalent.

**Closure-blockierend:** V-129 ja, bis die Frage im Nehmer steht; V-130 bis V-133 nein, sie gehen vor der Closure in §7 oder werden dort mit Ausgang eingetragen.

**Übergabe:**

- V-129 an den Planner: die Frage in `slice-v1-abschluss-einspielen-laufsteuerung` §4 (Prüfauftrag des Architect) und §6 (Risiko) mit der Kennung `slice-v1-abschluss-einspielen` eintragen (`AGENTS.md` §3.13); danach an den Architect: Reihenfolge und Zeitpunkt der Zeilen `error` bei `play` als Randform vor dem Code des Nehmers entscheiden (§3.12).
- V-130 und V-131 an den Implementer.
- V-132 an den Implementer: den Lauf mit Stand und Exit-Code in §7 eintragen.
- V-133 an den Planner: Wortlaut von DoD-Punkt 1 auf die Teile engen, die der Integrationstest trägt, oder an den Implementer: fehlende Teile als Integrationstest.
