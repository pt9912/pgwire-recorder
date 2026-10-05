# Validierung: slice-extended-query-replay, Extended Query im Replay — 2026-10-05

**Rolle:** Validator (Modul 8). Hier geht es darum, ob das Richtige gebaut wurde, gemessen am Bedarf. Ob richtig gebaut wurde, hat die [Verifikation](2026-10-05-verifikation-slice-extended-query-replay.md) geprüft, zuletzt in ihrer Nachverifikation zu `5a1da90` mit allen drei DoD-Liefer-Punkten bestätigt; ihr Ergebnis ist hier Eingang, kein Beleg.

**Gegenstand:** der Replay-Modus für das Extended Query Protocol, Stand HEAD `5a1da90`, aus der Sicht eines Anwenders mit einem verbreiteten Go-Treiber: pgx v5 im Standardmodus (`QueryExecModeCacheStatement`, Prepared Statements), am Treiber nur Host und Port umgestellt. Gefragt ist, ob dieser Anwender bekommt, was [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) und [Abnahmeszenario 7](../../spec/lastenheft.md#abnahmeszenario-7--extended-query) versprechen, und ob das Verhalten bei Abweichung, Fehlerantwort und Herunterfahren für ihn brauchbar ist. Drei bekannte Grenzen aus Plan §6 sind gegen den Bedarf zu bewerten: Herunterfahren ohne Frist (F-330), Diagnose ohne SQL-Befehle wie `DEALLOCATE` (F-348), synchrones Replay ohne Full-Duplex.

**Bedarfsquellen:**

- Lastenheft §1 und §2 (Stakeholder: Entwickler und Tester, Teams mit CI-Umgebungen)
- [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) Beschreibung („genügt … die Umstellung von Host und Port“), Happy, Boundary, Negative
- [LH-FA-04](../../spec/lastenheft.md#lh-fa-04--transparenz-für-clients), [LH-QA-02](../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung)
- [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [LH-FA-11](../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers), [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [LH-FA-14](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben)
- [Abnahmeszenario 3](../../spec/lastenheft.md#abnahmeszenario-3--keine-anwendungscode-integration), [4](../../spec/lastenheft.md#abnahmeszenario-4--replay-abweichung), [6](../../spec/lastenheft.md#abnahmeszenario-6--fehlerreplay) und [7](../../spec/lastenheft.md#abnahmeszenario-7--extended-query)
- aus der Spezifikation `LH-FA-18.a` (§Replay, §Mismatch), `LH-FA-10.a`, `LH-FA-13.a`
- das [Benutzerhandbuch](../user/benutzerhandbuch.md) und die [README](../../README.md)

**Grenze dieses Belegs:** Er misst gegen den dokumentierten Bedarf und gegen das, was eine Go-Anwendung mit pgx gewöhnlich tut. Eine Aussage des Auftraggebers liegt nicht vor. Wo dieser Beleg eine Lastenheft-Änderung empfiehlt, entscheidet der Auftraggeber. Nach Modul 8 liegt ein Validierungsbeleg gewöhnlich außerhalb des Repos; dass er hier unter `docs/reviews/` steht, hat der Auftrag so festgelegt.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

---

## 1. Beobachtetes Verhalten

**Sonde.** Kopie per `git archive HEAD` im Scratchpad. Das Binary ist aus der Stufe `deps` des `Dockerfile` gebaut, unverändert und netzlos. Daneben liegt ein eigenes Sondenprogramm, nur in der Kopie, mit pgx v5.11.0, `pgxpool` und `database/sql` über `pgx/v5/stdlib`. Für `pgxpool` war ein Netzabruf nötig, um `puddle` zu laden. Gelaufen ist alles in einem eigenen Docker-Netz, mit PostgreSQL aus dem gepinnten Image von `harness/mk/integration.mk`. Aufgezeichnet wurde gegen PostgreSQL, wiedergegeben mit gestoppter PostgreSQL. Danach waren Images, Container, Netz und Kopie gelöscht. Produktcode und Arbeitsbaum sind unverändert.

Der Ablauf „App“ ist breiter als `TestE2EReplayExtendedPgx`. Er nutzt ein `pgx.Conn`, legt eine Temp-Tabelle mit `int`, `text`, `numeric`, `jsonb`, `bool`, `timestamptz`, `bytea` und `uuid` an und fügt dreimal mit demselben Prepared Statement ein. Dazu kommen eine Schlüsselverletzung, `QueryRow` mit Treffer und mit `ErrNoRows`, ein `Query` über alle Typen, eine Transaktion mit Fehler und `Rollback`, `pgx.BeginFunc` mit `Commit` und ein Batch, in dem nach `22012` die folgende Anweisung verworfen wird.

| # | Lage | Beobachtung |
|---|---|---|
| S1 | „App“ aufgezeichnet, dann zweimal wiedergegeben | Die Ausgabe der Anwendung ist in beiden Läufen **byte-gleich** mit der Ausgabe beim Aufzeichnen. Exit-Code `0`, keine Warnung. |
| S2 | wie S1, im zweiten `INSERT` ein anderer Wert (`berta` statt `bert`) | pgx erhält `XX000 Replay [PGR-E5001]: Session 1, Interaktion 4, Gruppe 1, Nachricht 1: erwartet bind, empfangen bind, abweichend in params, Anweisung erwartet "INSERT …", empfangen "INSERT …"`. Danach ist die Verbindung geschlossen (`conn closed`). Auf `stderr` stehen dieselbe Meldung und `PGR-W2001` (23 von 26 nicht verbraucht), Exit-Code `5`. Welcher der acht Parameter abweicht, sagt die Meldung nicht. |
| S3 | Fehlerantworten aus S1 | `23505` und `22012` kommen mit Text an. In der Aufzeichnung stehen alle Felder der `ErrorResponse`, darunter Detail, Constraint-, Schema- und Tabellenname, dazu der Transaktionsstatus `E` nach dem Fehler in der Transaktion. |
| S4 | `pgxpool` (`MaxConns=1`), drei Anfragen nacheinander ohne Pause, aufgezeichnet; wiedergegeben ohne Pause | grün, Exit-Code `0` |
| S5 | wie S4, wiedergegeben mit 1,5 s Pause zwischen den Anfragen | **rot.** pgxpool prüft eine Verbindung, die länger als 1 s geruht hat, mit der einfachen Anfrage `-- ping`. Das Replay meldet `erwartet Client-Nachricht bind (Anweisung "SELECT $1::int + 1"), empfangen Anfrage "-- ping"`, Exit-Code `5`. Der Pool verwirft die Verbindung und öffnet neue; jede scheitert mit `PGR-E5003`. |
| S6 | `pgxpool` mit 1,5 s Pause aufgezeichnet (zwei `-- ping` in der Datei), wiedergegeben ohne Pause | **rot**, spiegelbildlich: `erwartet Anfrage "-- ping", empfangen Client-Nachricht bind …`, Exit-Code `5`. Mit derselben Pause wiedergegeben: grün. |
| S7 | `database/sql` mit `pgx/v5/stdlib`, ohne Pause aufgezeichnet, mit 1,5 s Pause wiedergegeben | **rot** wie S5. `ResetSession` sendet `-- ping`, sobald die Verbindung länger als 1 s geruht hat. |
| S8 | `pgxpool` mit vier Goroutinen gleichzeitig | in fünf von fünf Läufen rot (`PGR-E5001` „abweichend in params“, dann `PGR-E5003`). Das ist dokumentiert („Nutzen Sie die Verbindungen nacheinander“) und in `slice-v1-abschluss-sessions` §6 geführt. |
| S9 | Client verbunden und im Leerlauf, `docker stop -t 10` | Ende nach 0,3 s, Exit-Code `0`, `PGR-W2001` |
| S10 | Client sendet `Parse`/`Bind`/`Execute`/`Flush`, liest die Antwort und pausiert vor dem `Sync`; `docker stop -t 10` | Nach 10 s kommt `SIGKILL`, Exit-Code `137`. **Auf `stderr` steht keine Zeile** nach `replay gestartet`. Der Client sieht `unexpected EOF`. |
| S11 | wie S10, der Client sendet das `Sync` 2 s nach dem Signal | Der Client erhält `ReadyForQuery`, danach ist die Verbindung zu. Ende ohne `SIGKILL`, Exit-Code `0`. |
| S12 | `record` gegen das PostgreSQL-Image mit Standard-Anmeldung (SCRAM-SHA-256) | `PGR-E6001` „Anmeldeverfahren … nicht vermittelt“. Alle Sonden liefen deshalb mit `POSTGRES_HOST_AUTH_METHOD=trust`, wie der Integrationstest. |

**Aus dem pgx-Quelltext** (v5.11.0, gelesen im Modul-Cache):

- Die Namen im Statement-Cache sind deterministisch (`stmtcache_` + SHA-256 des SQL), also über Läufe stabil.
- `pgxpool` pingt beim Ausleihen, wenn die Verbindung länger als `time.Second` geruht hat (Default von `ShouldPing`). `stdlib.ResetSession` tut dasselbe. Der Ping ist `pgConn.Exec(ctx, "-- ping")`, also eine einfache Anfrage.
- Den eigenen Statement-Cache räumt pgx über die Protokoll-Nachricht `Close` (`pgConn.Deallocate`) auf, nicht per SQL. `deallocate all` sendet pgx nur auf ausdrücklichen Aufruf von `DeallocateAll`.
- Ein großer Schreibvorgang liest bei pgx parallel mit (`flushWithPotentialWriteReadDeadlock`).

---

## 2. Urteile

### Frage 1 — Abnahmeszenario 7 mit `pgx.Conn` im Standardmodus

**Urteil: trifft den Bedarf.**

- S1 zeigt das versprochene Verhalten für einen Ablauf, wie ihn eine Anwendung tatsächlich hat, nicht nur für den Testablauf: gängige Typen, Transaktionen samt Fehler und Rollback, `ErrNoRows`, Batch mit Verwerfen bis zum `Sync`. Am Treiber geändert waren nur Host und Port (und `sslmode=disable`, das das Handbuch nennt). Die Ausgabe ist über zwei Wiedergaben gleich.
- Die Boundary aus [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), dasselbe Prepared Statement mit verschiedenen Werten, trägt in S1: dreimal derselbe `INSERT`, und jede Ausführung bekommt ihre Position.

**Empfehlung:** keine Änderung. **Zielort:** keiner.

### Frage 2 — Connection-Pool und `database/sql`: Lebendprüfung `-- ping`

**Urteil: trifft den Bedarf nicht, ist aber kein Fehler dieses Slice.** Die Verifikation ist hier grün, die Validierung rot: Der Plan ist erfüllt, aber die Zusage des Lastenhefts „genügt … die Umstellung von Host und Port“ hält für die häufigste Form, in der Go-Anwendungen pgx nutzen, nicht verlässlich.

- Eine Anwendung im Test nutzt selten ein einzelnes `pgx.Conn`; meist nutzt sie `pgxpool` oder `database/sql`. Beide senden mit Default-Konfiguration je nach Leerlaufzeit eine Lebendprüfung (S5 bis S7). Ob sie in der Aufzeichnung steht und ob sie beim Wiedergeben kommt, hängt an Pausen über 1 s, also am Zeitverhalten: an einer langsamen CI-Maschine, einem Haltepunkt im Debugger, einem `time.Sleep` im Test oder einem Test, der wartet. Weil die Wiedergabe ohne Datenbank meist schneller ist als die Aufnahme, trifft S6 den Normalfall: Pings, die beim Aufzeichnen fielen, fehlen beim Wiedergeben.
- Für den Anwender sieht das wie ein flackernder Test aus. Die Meldung nennt `-- ping` zwar deutlich, aber eine Abhilfe kennt das Handbuch nicht. Die naheliegende Abhilfe (`ShouldPing` auf `false` in `pgxpool.Config`, `stdlib.ShouldPing` bei `database/sql`) ist eine Änderung der Treiberkonfiguration über Host und Port hinaus. Genau die schließen [Abnahmeszenario 7](../../spec/lastenheft.md#abnahmeszenario-7--extended-query) und [LH-QA-02](../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung) aus.
- Die Ursache liegt nicht im Extended-Pfad: `-- ping` ist eine einfache Anfrage, und strenges Replay ([ADR-0007](../plan/adr/0007-strict-replay.md)) lehnt jede Anfrage außer der Reihe ab. Abnahmeszenario 7 ist mit `pgx.Conn` formal erfüllt. Weil aber erst dieser Slice pgx-Anwender überhaupt zum Replay bringt, trifft es hier zum ersten Mal auf den realen Bedarf. In Spezifikation, ADRs und Planung führt den Fall bisher niemand.
- Gegenargument, geprüft: Mit konstantem Zeitverhalten (S4, S6 mit gleicher Pause) ist es grün. Ein Test, der ohne Pausen durchläuft, ist also nicht betroffen. Der Bedarf entscheidet sich aber nicht am Glücksfall: Ein Werkzeug für reproduzierbare Tests, dessen Ergebnis von der Rechengeschwindigkeit abhängt, verfehlt seinen Zweck ([LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay)).

**Empfehlung** (Entscheidung beim Auftraggeber, zwei Wege):

1. **Lebendprüfungen tolerieren** (empfohlen): Eine einfache Anfrage, die nur aus Leerraum und Kommentar besteht (PostgreSQL antwortet darauf mit `EmptyQueryResponse`), ist keine fachliche Interaktion. `record` zeichnet sie nicht als Interaktion auf, oder markiert sie. `replay` beantwortet sie an jeder Stelle zwischen zwei Interaktionen mit `EmptyQueryResponse` und `ReadyForQuery`, ohne den Cursor zu bewegen. Das lockert [ADR-0007](../plan/adr/0007-strict-replay.md) eng und braucht deshalb eine eigene ADR. Davon getrennt ist zu entscheiden, ob eine Aufzeichnung mit aufgezeichneten Pings rückwärts verträglich bleibt.
2. **Zusage verengen:** [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) bzw. [Abnahmeszenario 7](../../spec/lastenheft.md#abnahmeszenario-7--extended-query) nennen dann ausdrücklich, dass Lebendprüfungen eines Pools abzuschalten sind, und das Handbuch zeigt, wie (`ShouldPing`). Das ist billiger, bricht aber die Zusage „nur Host und Port“ für den häufigsten Einsatz.

**Zielort:**

- **Lastenheft:** Entscheidung zwischen 1 und 2 an [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) bzw. [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay). Das ist eine Spec-Validierung vor der Umsetzung.
- **Umsetzung (Weg 1):** ein eigener Slice. Er passt fachlich zu `slice-v1-abschluss-sessions`, weil dort schon das Verhalten von Connection-Pools steht (§6). Ob er dort hineinpasst, entscheidet der Schnitt.
- **Dieser Slice:** nicht die Umsetzung. Vor der Closure aber erstens der Befund als Beobachtung ins Register (oder als benannte Spec-Lücke in den Lerneintrag). Zweitens ein Hinweis im Handbuch, Abschnitt „Mit einem Datenbanktreiber arbeiten“, für die Zwischenzeit: Pools und `database/sql` senden nach einer Ruhezeit über 1 s eine Lebendprüfung, die das strenge Replay als Abweichung meldet; die Abhilfe ist `ShouldPing`.

### Frage 3 — Verhalten bei Abweichung

**Urteil: trifft den Bedarf. Eine Einschränkung betrifft die Diagnose bei Parametern.**

- Abweichung erkannt, nicht beantwortet, Verbindung zu, Exit-Code `5` (S2): Das erfüllen [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage) und [Abnahmeszenario 4](../../spec/lastenheft.md#abnahmeszenario-4--replay-abweichung). Für den Anwender wichtig: Die Meldung kommt **im Test selbst** als Fehler des Treibers an, mit Stelle und SQL. Er muss nicht erst das Log des Werkzeugs suchen. Die Art-Abweichung zeigt S5 ebenfalls lesbar, mit beiden Seiten und dem SQL der Anweisung.
- Die Folgefehler in der Anwendung (`conn closed`, `failed to deallocate cached statement(s)`) und die `PGR-E5003` neuer Pool-Verbindungen sind Rauschen. Die erste Meldung bleibt aber eindeutig die Ursache. Das trägt.
- **Einschränkung:** Bei `abweichend in params` und identischem SQL weiß der Anwender nicht, welcher von acht Parametern abweicht. Werte gehören nicht in die Diagnose (`SPEC-033`), der **Index** aber verrät nichts. Laut `LH-FA-18.a` §Mismatch erscheinen die Werte erst mit Log-Level `debug`; diese Option liefert erst `slice-replay-semantik-meldungscodes`.

**Empfehlung:** Die Diagnose nennt den Index des ersten abweichenden Parameters (etwa `params[2]`), ohne Wert. Für Formate gilt dasselbe. **Zielort:** `LH-FA-18.a` §Mismatch, umgesetzt in `slice-replay-semantik-meldungscodes` zusammen mit `--log-level`. Nicht blockierend.

### Frage 4 — Fehlerantwort (Abnahmeszenario 6, Extended)

**Urteil: trifft den Bedarf.** Code, Text und alle Felder der `ErrorResponse` kommen an, auch Constraint- und Tabellenname, die Anwendungen oft auswerten (S3). Ebenso kommen der Transaktionsstatus nach dem Fehler und das Verwerfen im Batch an (S1). **Empfehlung:** keine. **Zielort:** keiner.

### Frage 5 — bekannte Grenze 1: `replay` wartet nach `SIGTERM` unbegrenzt auf einen pausierenden Client (F-330)

**Urteil: für die Zwischenzeit akzeptabel. Die Bindung an `slice-v1-abschluss-betrieb` trägt.**

- Anders als bei `record` geht hier nichts verloren: `replay` schreibt nichts, und der Test hat sein Ergebnis schon beim Client. Schaden entsteht nur durch Exit-Code `137` statt eines Codes der Spezifikation und durch ein stummes Log (S10).
- Bei pgx im Standardmodus kommt die Lage praktisch nicht vor: Jede Folge endet sofort mit `Sync`, und ein Client im Leerlauf beendet sofort (S9). Ein Client, der zwischen `Flush` und `Sync` hängt, ist ein Test, der ohnehin hängt. Wer den Ablauf abschließt, erhält die Antwort (S11). Das ist genau das Verhalten, das `LH-FA-13.a` zusagt.
- Der Folge-Slice nimmt den Fall schon richtig auf: Frist ab dem Signal, `PGR-E4006` bei unvollständiger Interaktion, Info-Zeile beim Beginn des Herunterfahrens. Diese Zeile behebt auch das stumme Log.

**Empfehlung:** Im Handbuch steht die Lage für `replay` heute nur im Treiber-Abschnitt („wartet, bis … mit ihrem `Sync` abgeschlossen ist“). Die Folge im Container (`SIGKILL` nach der Stopp-Frist, Exit-Code `137`, keine Meldung) sollte dort für die Zwischenzeit ebenfalls stehen, wie für `record`. **Zielort:** dieser Slice, bei der Closure. Das Risiko F-330 hat seinen Ausgang („eingetreten“ → `slice-v1-abschluss-betrieb`) schon.

### Frage 6 — bekannte Grenze 2: Diagnose bildet `DEALLOCATE`, `DISCARD ALL` usw. nicht nach (F-348)

**Urteil: trifft den Bedarf.** Die Erkennung der Abweichung ist nicht berührt; ungenau kann nur der SQL-Text in der Meldung sein. Den eigenen Cache räumt pgx über die Protokoll-Nachricht `Close` auf, und die bildet die Diagnose nach. SQL-seitiges `deallocate all` sendet pgx nur auf ausdrücklichen Aufruf. `DISCARD ALL` stammt typischerweise von externen Poolern wie PgBouncer, die im Testaufbau zwischen Anwendung und Recorder nicht vorkommen. Selbst dann zeigt die Meldung die Stelle richtig an, nur das SQL ist veraltet. **Empfehlung:** keine Änderung. **Zielort:** keiner; das Risiko kann bei der Closure den Ausgang „entfallen“ (vom Bedarf getragen) bekommen.

### Frage 7 — bekannte Grenze 3: `replay` bleibt synchron

**Urteil: trifft den Bedarf.** PostgreSQL bedient einen Client ebenso abwechselnd: Es liest eine Nachricht, schreibt die Antwort und blockiert, wenn der Client nicht liest. Ein Client, der gegen PostgreSQL nicht verklemmt, verklemmt deshalb auch gegen `replay` nicht. Der Record braucht Full-Duplex ([ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md)), weil er **zwischen** zwei Gegenstellen sitzt und keine von beiden aufhalten darf. Für das Replay als Server gilt das nicht. pgx liest beim Schreiben mit; den Gegendruck-Ablauf hat die Verifikation grün gesehen (`TestE2EReplayExtendedGegendruck`). Der ungeprüfte Fall aus F-329 verklemmt gegen PostgreSQL genauso. **Empfehlung:** keine. **Zielort:** keiner.

---

## 3. Nebenbefunde

- **Aufnahme gegen eine Standard-PostgreSQL (S12):** Mit dem Default-Image (SCRAM-SHA-256) ist Abnahmeszenario 7 heute nicht durchführbar. Der Anwender muss die Datenbank auf `trust` stellen. Das ist bekannt: Die README sagt, die Passwort-Anmeldung fehle noch, und `slice-v1-abschluss-anmeldung` liefert sie. Erwähnt ist es hier, weil es die erste Hürde für genau den Anwender dieses Belegs ist. **Zielort:** keiner neu. Vorschlag: Das Handbuch nennt im Abschnitt „Eine Anwendung aufzeichnen“ die Voraussetzung `trust` ausdrücklich, solange der Slice offen ist.
- **Prozess:** Die Verifikation (Nachverifikation zu `5a1da90`) ist grün, die Validierung in diesem Beleg ebenfalls, mit Ausnahme von Frage 2. Damit liegt keine Prozess-Drift vor (Verifikation rot, Validierung grün). Die Art-Abweichung, die die Verifikation als V-25 zunächst offen hatte, ist in S5 und S6 im realen Ablauf lesbar.

## 4. Übergabe an den Planner

| Nr. | Gegenstand | Urteil | Zielort |
|---|---|---|---|
| 1 | Abnahmeszenario 7 mit `pgx.Conn` | grün | keiner |
| 2 | `pgxpool`/`database/sql`: Lebendprüfung `-- ping` abhängig von der Zeit | **rot** (Plan erfüllt, Bedarf nicht) | Lastenheft [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol)/[LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), Entscheidung des Auftraggebers zwischen „tolerieren“ (ADR gegen [ADR-0007](../plan/adr/0007-strict-replay.md), eigener Slice, Nähe zu `slice-v1-abschluss-sessions`) und „Zusage verengen“; in diesem Slice Register oder Spec-Lücke im Lerneintrag sowie ein Hinweis im Handbuch |
| 3 | Abweichung: Diagnose ohne Parameter-Index | grün mit Einschränkung | `LH-FA-18.a` §Mismatch, `slice-replay-semantik-meldungscodes` |
| 4 | Fehlerantwort Extended | grün | keiner |
| 5 | F-330: Herunterfahren ohne Frist | akzeptabel bis zur Folge | `slice-v1-abschluss-betrieb` (gebunden); in diesem Slice Hinweis im Handbuch zu `replay` im Container |
| 6 | F-348: SQL-Befehle in der Diagnose | grün | keiner; Ausgang „entfallen“ |
| 7 | synchrones Replay | grün | keiner |
| — | Aufnahme nur mit `trust` | bekannt | `slice-v1-abschluss-anmeldung`; optional Hinweis im Handbuch |

Keiner dieser Punkte blockiert die Closure dieses Slice, sofern Punkt 2 vorher als Beobachtung oder benannte Spec-Lücke mit Weg zum Auftraggeber festgehalten ist und das Handbuch die Abhilfe für die Zwischenzeit nennt. Vor der nächsten Welle, die Replay-Semantik oder Sessions anfasst, sollte Punkt 2 entschieden sein.
