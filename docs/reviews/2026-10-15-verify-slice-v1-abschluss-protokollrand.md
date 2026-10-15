# Verifikation slice-v1-abschluss-protokollrand

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-15. **Spanne:** `5c4b3ae` (Basis, Übergang nach `in-progress/`) bis `eccdbcb` (HEAD). **Gegenstand:** Slice `slice-v1-abschluss-protokollrand` (Plan unter `docs/plan/planning/in-progress/`) gegen DoD (§2), Plan §1 bis §8 (§6: Entscheidungen des Architect vom 2026-10-15), `LH-FA-05.b` *Auswahl im Replay*, `LH-FA-05.e` *Randformen des Protokollrands*, `LH-FA-12.a` *Zählung der Verbindungen*, `LH-FA-18.a`, `LH-FA-23.a` *Unverschlüsselte Clients* und `AGENTS.md` §3.9 bis §3.14. Review: `docs/reviews/2026-10-15-review-slice-v1-abschluss-protokollrand.md` (F-1 bis F-7; Nacharbeit `0347c8f`, `eccdbcb`); gegen ihn ist nicht verifiziert, die Nacharbeit ist unabhängig nachgeprüft.

Eine Verifier-Vorlage führt die vendored Baseline nicht; Aufbau wie der Report von `slice-v1-abschluss-postgres-versionen-zusage`. Belege sind in diesem Lauf selbst erzeugt. Mutationen liefen in Kopien von HEAD (`git archive`, Scratchpad, danach gelöscht); der Arbeitsbaum blieb unberührt (`git status` leer). Läufe nacheinander, kein `pkill`; ein hängender Testcontainer (Mutant im Upstream-Adapter, Test blieb stehen) wurde per Container-ID beendet. Danach 0 Container und 0 Netze.

**Verdikt: PASS mit Befunden.** Kein Liefer-Punkt der DoD scheitert an einem Lauf; `make gates` ist grün (zweiter Lauf, siehe unten). 0 HIGH, 1 MEDIUM, 3 LOW, 3 INFO. Der MEDIUM-Befund (V-1) verlangt eine Rückgabe an den Architect vor der Closure: Er betrifft die Reihenfolge `PGR-E6001` (`_pq_.`) vor `PGR-E6003`, die der Code entscheidet und die Spezifikation im Wortlaut nicht trägt.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren) |
|---|---|---|
| 1: Jede Zeile der Tabelle in `LH-FA-05.e` durch einen Test belegt (Happy/Boundary/Negative), Deckung, `replay` mit einem Fall oder Grenze | bestätigt | Zeilen gegen Tests gelesen: Version, Sonderanfrage, `CancelRequest` Länge 8, Folgen von `SSLRequest`/`GSSENCRequest`, Rumpf, `_pq_.`, Client-Nachrichten, 0 bis 7 Bytes, gleichzeitig (`server_protokollrand_test.go`, beide Modi, mit und ohne TLS); Serverantworten (`upstream_protokollrand_test.go`, `record` nur, Grenze steht in der Spezifikation *Modus*); `LISTEN`/`NOTIFY` und TCP-Probe und `_pq_.` ohne Upstream (E2E). `make test-integration`: Exit 0, 2 min 13 s, 251 `--- PASS`, 0 FAIL, 2 SKIP (Phasenwechsel), die fünf neuen E2E-Tests PASS. `make abdeckung-check` und `make kopf-check` ohne Befund. |
| 2: Replay-Handshake emuliert feste Auswahl, Test belegt sie samt Randformen | bestätigt | `TestE2EReplayAuswahl` und Unit-Zeilen (Benutzer/Datenbank/`options` ungeprüft, `p` → `PGR-E6001`, Ende nach `AuthenticationOk`, ohne und mit TLS im Unit-Test); in den Mutanten unten rot (`ende-io-eof` an `p` ohne Nullbyte). Grenze `require_auth` nur mit `pgconn`, nicht `psql` (in §6 benannt). |
| 3: `make gates` grün | bestätigt (nach Wiederholung) | Erster Lauf: Exit 2, `postgres-versionen-gegenprobe` rot, Mutation `p7-int-treiber` (erwartet rot: `p7-sigINT-gruppe p7-sigINT-skript`, rot waren zusätzlich `p7-sigTERM-skript`). Der Diff berührt `tools/`, `harness/`, `.claude/` und `Makefile` nicht (`git diff --stat 5c4b3ae..HEAD -- tools harness Makefile .claude` leer); `make test-postgres-versionen-gegenprobe` allein: Exit 0 (2 min 26 s). Zweiter voller Lauf `make gates`: Exit 0, 10 min 26 s, Versionsmatrix 14 bis 18 grün, Stop-Hook-Gegenprobe 139 Fälle, `docs-check` 779 Dateien 0 Befunde. Ein Signal-Zeitfenster der Gegenprobe ist also nicht stabil (Hinweis V-8). |
| 4: Review durchgeführt, Report liegt vor | bestätigt | Report liegt vor; Ausgang je Finding unten. |
| 5 bis 8: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | §7 trägt „Wird bei Closure gefüllt"; kein Häkchen in §2; §6-Risiken tragen „offen bis Closure". Handlung des Planners. |

## Ausgang jedes Review-Findings (Nacharbeit unabhängig nachgeprüft)

| Finding | Ausgang laut Plan | Ergebnis | Urteil |
|---|---|---|---|
| F-1 MEDIUM (Prüfreihenfolge Länge vor Code ohne Mutation) | `TestProtokollrandPruefreihenfolge` | Mutant `(laenge < 8 \|\| laenge > max) && code != codeProtocol30` → rot (3.0, Länge 4/7/10001, je 4 Läufe); dasselbe für `codeCancelRequest`; `laenge > max` → `>=` rot (`TestTLSSonderanfrageMitFalscherLaenge`); `laenge < 8` → `<= 8` rot (Startnachricht); `laenge > max` → `laenge < 8`: 7 Codes × Länge 10001 rot. Der Rest wird nicht gelesen: `code == 4<<16` und `65535<<16\|1` an pgproto3 weitergereicht → `Länge 40 ohne Rest` rot | trägt |
| F-2 MEDIUM (Zählregel `LH-FA-12.a` im Wortlaut zweimal verschieden) | Spezifikation angeglichen (*Angenommene Verbindung*, *Zählung der Verbindungen*, *Zuordnung bei `connection`*), `TestProtokollrandAbgelehnteVerbindungZaehltNicht` | Wortlaut „gelesen und vom Recorder angenommen" an allen drei Stellen und in der Zeile *Modus*; Lastenheft unberührt. Mutanten: `annehmen` setzt nach `PGR-E6003` fort (`return true`) → `TestProtokollrandAbgelehnteVerbindungZaehltNicht` und `TestTLSUnverschluesselt` rot; nach `_pq_.` fortsetzen → `…ZaehltNicht`, `Startnachricht`, `ErweiterungVorKlartext` rot; im Replay `OpenConnection` vor der Fehlerantwort der Version (E2E) → `BeideModi/replay/Version_3.2` rot | trägt |
| F-3 LOW (Rückverweis auf nicht vorhandene Überschrift) | auf *Startnachricht 3.0* | Verweis in *Auswahl im Replay* zeigt auf *Startnachricht 3.0*, die Überschrift existiert | trägt |
| F-4 LOW (Kommentar von `startup`) | neu gesetzt | Kommentar liest sich als ein Absatz | trägt |
| F-5 LOW (Spalte „nur dieser" der Mutationstabelle) | Spalte gemessen, Kopf auf „genannte Tests" | Kopf ist berichtigt; **Zeile `rest-gelesen-<code>` behält „nur diese"** und stimmt nicht (V-2) | trägt nur teilweise |
| F-6 INFO | keine Aktion | Grenzen in §6 gelesen; Ergänzung durch V-4 | — |
| F-7 INFO (TCP-Probe nur Stichprobe 0, 1, 7) | E2E läuft 0 bis 7, geschlossen und Reset | Code gelesen (`for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 7}`); Mutant `bytes-5-e2e` (bei 5 gepufferten Bytes `PGR-E6001` statt stilles Ende) → `TestE2EProtokollrandTCPProbe` rot, Exit 2, Meldung im Lauf. Die Zeile der Mutationstabelle nennt noch „0, 1 und 7" (V-6) | trägt |

## Mutanten dieses Laufs

Klassen laut Auftrag; alle Unit-Läufe `go test ./internal/adapters/driving/pgwire/` (und je genanntem Paket) in Kopie von HEAD, E2E über `make test-integration` mit `PGR_TESTS_1` = die drei E2E-Tests des Slice. Basis jeweils grün.

| Klasse | Mutant | Ergebnis |
|---|---|---|
| Rumpf-Ablehnung | `ReceiveStartupMessage`-Fehler endet still (ohne `fail`) | rot: `TestProtokollrandStartnachricht` (Länge 8, Länge 9, ohne Nullbyte, Kopf ohne Rest, `CancelRequest` Länge 8; alle vier Läufe je Modus/TLS) |
| Versions-Ablehnung | `code == codeProtocol30` → `code>>16 == 3` | rot: Zeilen 3.1 und 3.2 ×4, `TestAndereProtokollversion` |
| Versions-Ablehnung E2E | derselbe Mutant | rot: `BeideModi/record/Version_3.2` und `…/replay/Version_3.2` |
| Zählung abgelehnter Verbindung | Replay: `OpenConnection` vor der Versionsablehnung (E2E) | rot: `BeideModi/replay/Version_3.2` (reguläre Verbindung erhält die zweite Session) |
| Zählung abgelehnter Verbindung | `annehmen` setzt nach `PGR-E6003` und nach `_pq_.` fort | rot (Unit), siehe F-2 |
| `_pq_.`-Prüfung | Prüfung abgeschaltet (`if false && …`) | rot: `Startnachricht`, `…ZaehltNicht`, `ErweiterungVorKlartext`; E2E: `BeideModi/replay/Parameter__pq_.` und `TestE2ERecordErweiterungOhneUpstream`. `BeideModi/record/Parameter__pq_.` bleibt grün (V-7) |
| Längenprüfung vor Code | siehe F-1 (5 Mutanten) | rot |
| Rest nicht gelesen | `rest-gelesen-4.0`, `rest-gelesen-65535.1` | rot, aber nicht „nur diese" (V-2) |
| Leser-Wrapper EOF vs. Rumpf | `verbindungsende` zählt zusätzlich `io.EOF` | rot: `TestAnmeldungAndereNachricht/Rumpf_ohne_Nullbyte`, `TestProtokollrandClientNachrichten` (`p` und `Parse` ohne Nullbyte, 4 Läufe) |
| Leser-Wrapper | `io.EOF` aus `endeDesStroms` entfernt | rot: Anmeldung, Extended, Startnachricht, Client-Nachrichten, Folgen |
| Leser-Wrapper | `syscall.ECONNRESET` entfernt | rot: `TestAnmeldungEndeDesClients` |
| Leser-Wrapper | `io.ErrUnexpectedEOF`, `net.ErrClosed`, `os.ErrDeadlineExceeded` je einzeln entfernt | **grün** (V-4) |
| TCP-Probe | siehe F-7 | rot |
| Serverantworten | `toResponse` bildet `CopyOutResponse` als `NoData` ab | rot: `TestServerNachrichtenNichtUnterstuetzt/CopyOutResponse/Anfrage` |
| Zusage „schließt die Verbindung zum Server" | `CloseSession` ruft `upstream.Close()` nicht | rot im Bestand (`TestRecordZwangsende`, `TestRecordCloseBeendetWartende`); der Slice fügt keinen Test dafür hinzu |

## Findings

### V-1 MEDIUM

- `quelle`: `AGENTS.md` §3.12 (Randform ist vor dem Code in der Spezifikation entschieden), `LH-FA-23.a` *Unverschlüsselte Clients*, Plan §7 *Randformen und Lesarten*
- `pfad`: `spec/spezifikation.md` · `LH-FA-23.a` *Unverschlüsselte Clients*: „Das gilt nur für eine `StartupMessage` der Version 3.0; eine andere Version ist `PGR-E6002`, eine unbekannte Sonderanfrage `PGR-E6001` und eine erste Nachricht, die keine Startnachricht ist, `PGR-W3003`"; `LH-FA-05.e` *Startnachricht 3.0*: „Ein Parameter mit dem Präfix `_pq_.` … ist `PGR-E6001` … in beiden Modi"; `internal/adapters/driving/pgwire/server.go` · `annehmen` (erst `erweiterung`, dann `s.tls != nil && !s.tls.klartext`)
- `befund`: Eine unverschlüsselte Startnachricht 3.0 mit `_pq_.` an einen Recorder mit TLS-Pflicht ist im Code `PGR-E6001`, nicht `PGR-E6003`. Die Spezifikation sagt es an keiner Stelle: Der Absatz in `LH-FA-23.a` nennt `PGR-E6001` nur für die unbekannte Sonderanfrage und liest sich für jede Startnachricht 3.0 als `PGR-E6003`; *Startnachricht 3.0* nennt die Reihenfolge zu `PGR-E6003` nicht; die Zeile *Modus* sagt „mit und ohne TLS, soweit `LH-FA-23.a` nichts anderes regelt", was auf `PGR-E6003` zeigt. Plan §7 (*Randformen und Lesarten*) belegt „Keine Randform wurde im Code entschieden" mit „steht in `LH-FA-23.a` … (E6003 hinter E6002, E6001 und W3003) und im Drift-Log"; das Drift-Log (Zeile 2026-10-13) benennt `_pq_.` nicht, die Zeile 2026-10-15 nennt `PGR-E6003` nicht. Die Reihenfolge steht nur in der Deklaration von `TestProtokollrandErweiterungVorKlartext`, in Plan §1 und im Kommentar von `annehmen`. Der Test belegt den Code (Mutant `pq-nach-klartext`), nicht die Spezifikation.
- `verifizierbar`: ja (`grep -n "_pq_" spec/spezifikation.md` zeigt keine Stelle zur Reihenfolge; `sed -n 2333,2342p`)
- `klasse`: Randform im Code entschieden, Spezifikation sagt es nicht (`BEO-REPO/spec-randform-erst-im-review-entschieden`, `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`)
- Rückgabe an den Architect: in *Startnachricht 3.0* den Satz „vor `PGR-E6003`" ergänzen (mit Verweis auf `LH-FA-23.a`) und in `LH-FA-23.a` *Unverschlüsselte Clients* `_pq_.` neben der Sonderanfrage nennen; Drift-Log-Zeile ergänzen. Plan §7 berichtigt dann den Beleg.

### V-2 LOW

- `quelle`: `AGENTS.md` §3.11 (Quantor „nur"), Plan §7 Nacharbeit zu F-5
- `pfad`: `docs/plan/planning/in-progress/slice-v1-abschluss-protokollrand.md` · Mutationszeile „der Rest einer Startnachricht mit abgelehntem Code wird nicht gelesen … je die Zeile `<Code>, Länge 40 ohne Rest`, nur diese"
- `befund`: Selbst gefahren: `rest-gelesen-4.0` färbt zusätzlich `TestProtokollrandStartnachricht/Version_4.0` (×4), `TestProtokollrandAbgelehnteVerbindungZaehltNicht/Version_4.0` (×2), `TestAndereProtokollversion`; `rest-gelesen-65535.1` zusätzlich `TestProtokollrandStartnachricht/Version_65535.1` (×4). „Nur diese" ist falsch. Die Zusage (die Zeile färbt) trägt; die Nacharbeit zu F-5 hatte den Wortlaut für die übrigen Zeilen zurückgenommen, diese Zeile blieb.
- `verifizierbar`: ja
- `klasse`: Zusage in Belegen weiter als die Messung (Wiederkehr von F-5)

### V-3 LOW

- `quelle`: `AGENTS.md` §3.11 (Satz in der Spezifikation sagt nur zu, was ein Test oder der Code zeigt)
- `pfad`: `spec/spezifikation.md` · `LH-FA-05.e` *Sonderanfrage*: „Beim `CancelRequest` ist die Länge 16 zugesagt (Prozess-ID und Schlüssel, 4 Bytes); der Recorder liest beides nicht."
- `befund`: `startkopf` gibt für den Code `CancelRequest` `true` zurück, danach ruft `startup` `ReceiveStartupMessage` auf, die Bibliothek liest und decodiert Prozess-ID und Schlüssel (deshalb lehnt sie Länge 12 ab). „Liest beides nicht" ist im Wortlaut falsch; gemeint ist „wertet nicht aus". Kein Test hängt daran.
- `verifizierbar`: ja (Code-Pfad `server.go` `startup`, `case *pgproto3.CancelRequest`)
- `klasse`: Aussage der Spezifikation weiter als der Code
- Vorschlag an den Architect: „leitet beides nicht weiter und wertet es nicht aus".

### V-4 LOW

- `quelle`: `AGENTS.md` §3.10 (Aufzählung zählt je Element), §3.11
- `pfad`: `internal/adapters/driving/pgwire/leser.go` · `streamFehler` „(Schließen, Abbruch, Frist)", `endeDesStroms`
- `befund`: Die Liste der fünf Fehler, die das Ende des Datenstroms kennzeichnen, ist mit dem Leser neu im Slice. Selbst gefahren: Je Element entfernt (`io.ErrUnexpectedEOF`, `net.ErrClosed`, `os.ErrDeadlineExceeded`) bleibt `go test` des Pakets grün; `io.EOF` und `syscall.ECONNRESET` färben. `io.ErrUnexpectedEOF` ist äquivalent (der Wrapper kennzeichnet schon das `io.EOF` darunter); für `net.ErrClosed` und die Frist ist offen, ob sie äquivalent sind oder eine Randform ohne Fall (das Zwangsende setzt Lesefristen, `server.go` Zeile 706 prüft die Frist vor `verbindungsende`). Plan §6 *Grenzen der Tests* nennt diese Mutanten nicht.
- `verifizierbar`: ja (`errors.Is(err, net.ErrClosed)` → `(… && false)` in `endeDesStroms`)
- `klasse`: Aufzählung ohne brechende Mutation je Element (verschobener Bestand)

### V-5 LOW

- `quelle`: `AGENTS.md` §3.13 (Sendung nimmt an, bevor sie genannt wird), §3.11
- `pfad`: `docs/plan/planning/next/slice-v1-abschluss-protokollrand-doku.md` · §1 *Ziel*, §2 erste DoD-Zeile
- `befund`: Der Doku-Folge-Slice nennt als Liefer-Inhalt nur „was `slice-v1-abschluss-protokollrand` liefert" und lässt den Architect vor dem ersten Commit prüfen, ob der Code-Slice eine weitere Randform liefert. Er nennt keinen der neuen Inhalte: `_pq_.` in beiden Modi `PGR-E6001`, nicht lesbarer Rumpf `PGR-E6001` statt stilles Ende (Handbuch-Grenze zu `PGR-E6001`), Zählung abgelehnter Verbindungen, `NotificationResponse` und `LISTEN`/`NOTIFY` in einer Session → `PGR-E6001`, feste Auswahl im Replay. `grep -n "_pq_\|NotificationResponse\|Zählung\|abgelehnte"` im Plan: keine Treffer. Der Code-Slice führt unter „Gebundene Sendungen" nur eingehende. Der Handbuch- und README-Stand nennt diese Verhalten heute nicht (`grep` leer); das Handbuch hinkt wie vorgesehen hinter dem Binary her. Die allgemeine Zeile trägt die Sendung; die Aufzählung fehlt, wie sie `slice-doku-ist-stand` §6 für Randformen verlangt, und ob der Architect sie rechtzeitig liefert, hängt an einem Satz in §4.
- `verifizierbar`: ja
- `klasse`: Sendung ohne benannten Inhalt beim Nehmer
- Vorschlag: der Planner trägt die Liste in §1 und die erste DoD-Zeile des Nehmers ein, mit der Kennung des Gebers, und zählt den Nehmer nach (aktuell eine Schicht, zwei Liefer-Punkte).

### V-6 INFO

- `quelle`: `AGENTS.md` §3.9
- `befund`: Plan §7, Zeile „0 bis 7 Bytes enden still": „E2E `TCPProbe` für 0, 1 und 7" ist nach der Nacharbeit überholt (der E2E läuft 0 bis 7); die Zeile für `bytes-3-e2e` weiter unten nennt die frühere Stichprobe richtig als Vorzustand. Sonst folgt der Plan dem Diff: Größe 1172 (selbst gemessen: `git diff --numstat 5c4b3ae..HEAD -- internal test spec docs/user` = 1137 + 35; Spezifikation 123, Produkt-Code 117, Tests 903, Tabellen 29), unter Meldeschwelle 1250 und Marke 1400; §8 zählt 2 Liefer-Punkte und 2 Schichten (PGWire-Adapter mit Code und Tests, Upstream-Adapter nur Tests), der Diff berührt nur diese; Kern, Ports, CLI, Bootstrap, Handbuch und README unberührt, kein `//nolint`, kein `pgproto3` in `test/integration`.
- `klasse`: Plan folgt dem Diff nicht in einer Zeile

### V-7 INFO

- `befund`: Unter dem Mutanten „`_pq_.`-Prüfung aus" bleibt `TestE2EProtokollrandBeideModi/record/Parameter__pq_.` grün: Der echte Server antwortet mit `NegotiateProtocolVersion`, die Aufzeichnung meldet das als `PGR-E6001`. Die Zusage „vor dem Aufbau zum Upstream" im Record trägt allein `TestE2ERecordErweiterungOhneUpstream` (rot, Exit 2), und dort ist es der richtige Grund. Plan §6 *Grenzen der Tests* nennt das nicht; keine Aktion nötig, nur Hinweis für die Closure.

### V-8 INFO

- `befund`: (a) Die Skripte und Rohausgaben der Mutanten liegen laut Plan §7 im Scratchpad der Sitzung und sind kein Beleg im Repo; die Zeilen der Tabelle sind aber ohne sie nachfahrbar (Mutation ist im Text genau beschrieben), und alle hier gefahrenen stimmen mit der Spalte „roter Test" überein, bis auf V-2. (b) Der erste `make gates` des Verifiers war an `postgres-versionen-gegenprobe` (Mutation `p7-int-treiber`, zusätzlich rot `p7-sigTERM-skript`) rot, beim Einzellauf und im zweiten Gate-Lauf grün; der Slice berührt die Gate-Kette nicht. Das ist eine Beobachtung zu einer zeitabhängigen Gegenprobe, kein Befund gegen diesen Slice.

## Entscheidungen des Architect (§1 und §6) gegen den Diff

| Entscheidung | Urteil |
|---|---|
| `_pq_.` beide Modi `PGR-E6001`, vor dem Aufbau zum Upstream | eingehalten (`annehmen`, `handle` ruft `oeffne` danach; E2E ohne Upstream); Reihenfolge zu `PGR-E6003` siehe V-1 |
| Ungültiger Rumpf einer vollständig gelesenen Client-Nachricht `PGR-E6001` statt stilles Ende, auch in der Anmeldung | eingehalten (Leser, `verbindungsende`; Mutanten `ende-io-eof` rot) |
| `LISTEN`/`NOTIFY` berichtigt, beides in einer Session `PGR-E6001` | eingehalten; E2E gegen PostgreSQL 17 PASS (`TestE2ERecordListenNotify`); andere Versionen als Grenze in der Spezifikation benannt |
| Zählung der Verbindungen (F-2): erst gelesen und angenommen zählt | eingehalten, Lastenheft und Code unverändert (`git diff 5c4b3ae..HEAD -- spec/lastenheft.md` leer) |
| Auswahl im Replay fest | eingehalten (`AuthenticationOk` Code 0 zuerst, kein Verfahren, nichts geprüft) |
| Je Zeile ein Fall und eine Mutation, Aufzählungen je Element (`AGENTS.md` §3.10) | eingehalten bis auf V-4 (Leser) und die in §6 benannten Grenzen |

## Quantoren (`AGENTS.md` §3.11) und Sendungen (§3.13)

- Geprüft: „jede andere Version … schon bei Länge 8", „jede Folge aus zwei Anfragen", „jede andere [Nebennummer] vor dem Lesen des Rests, mit und ohne TLS", „eine Länge unter 8 oder über 10000 … gleich welcher Code", „keine Meldung/ohne Session" der Randverbindungen: je Element ein Fall und ein roter Mutant. Offen: V-3 („liest beides nicht"), V-4 (Leser-Liste).
- Sendungen: `CancelRequest` Länge 12 ist bei `slice-v1-abschluss-cancel-ohne-schluessel` in §1 und DoD mit der Kennung dieses Gebers angenommen (Unit-Test mit 12 Byte); die Spezifikation sagt „andere Längen nicht zugesagt". Handbuch und README bei `slice-v1-abschluss-protokollrand-doku`: angenommen, aber ohne benannten Inhalt (V-5). `slice-v1-abschluss-verbindungsende-format` ist Reihenfolge, kein Gegenstand. Keine neue offene Adresse im Diff.
- Die zwei offenen Nutzerfragen (weitere Anmeldeverfahren im Replay, Aufzeichnung von `NotificationResponse`): im Code nicht umgesetzt (`grep` auf `AuthenticationCleartext`, `AuthenticationMD5`, `SASL`, `noch nicht`, `später`, `geplant` im Diff von `internal` und `test` ohne Treffer; die Spezifikation führt beide als „nicht zugesagt"/„nicht aufgezeichnet"), in keinem Kommentar als Zusage versteckt.

## Ausgänge der Risiken aus §6 (Kandidaten, Urteil beim Planner)

| Risiko | Kandidat | Beleg |
|---|---|---|
| Inventur sprengt Größe oder Review-Sitzung | *entfallen*: 1172, Meldeschwelle 1250 und Marke 1400 nicht erreicht, zwei Schichten | Messung oben |
| `LISTEN`/`NOTIFY` in einer Session nicht aufgezeichnet, Zeile zu berichtigen | *entfallen*: Zeile berichtigt, Lastenheft unberührt, E2E gegen 17 grün; Verhalten der übrigen Versionen als Grenze benannt (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` bleibt offen) | Spezifikation, E2E |
| Auswahl im Replay genügt einem Treiber mit Anmeldung nicht | *weiter offen* oder *entfallen* mit der Nutzerfrage (1): Grenze in der Spezifikation, Fall im E2E (`require_auth`); die Entscheidung des Nutzers steht aus | Spezifikation, `TestE2EReplayAuswahl` |
| Reihenfolge `_pq_.` zu `PGR-E6003` (neu aus V-1) | *eingetreten* → Rückgabe an den Architect | V-1 |

## Gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make test-integration` (HEAD) | Exit 0, 2 min 13 s, 251 PASS, 0 FAIL, 2 SKIP |
| `go test ./internal/adapters/driving/pgwire/` (Basis HEAD) und Unit-Mutanten | Basis `ok`; Mutanten wie oben (grün nur `io.ErrUnexpectedEOF`, `net.ErrClosed`, `os.ErrDeadlineExceeded`) |
| `go test ./internal/adapters/driven/postgres/` mit Mutant | rot (`CopyOutResponse/Anfrage`) |
| `go test ./internal/hexagon/services/` mit Mutant | rot (Bestand) |
| E2E-Mutanten (`make test-integration` mit `PGR_TESTS_1`) | `bytes-5-e2e`, `v32-verbraucht`, `pq-weg`, `v32-durch` alle rot an der genannten Stelle |
| `make abdeckung-check`, `make kopf-check`, `make docs-check` | ohne Befund (779 Dateien) |
| `make test-postgres-versionen-gegenprobe` | Exit 0 (Einzellauf nach rotem Erstlauf) |
| `make gates` | Lauf 1 Exit 2 (`postgres-versionen-gegenprobe`, V-8), Lauf 2 Exit 0 (10 min 26 s) |
| Aufräumen | 0 Container, 0 Netze, `git status` leer |
