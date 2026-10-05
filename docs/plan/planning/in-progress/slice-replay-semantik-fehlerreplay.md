# Slice slice-replay-semantik-fehlerreplay: Fehlerantworten, Resultsets und Transaktionen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-replay-semantik.

**Bezug:** [`LH-FA-11`](../../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers), [`LH-FA-06`](../../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [ADR-0006](../../adr/0006-kanonisches-domain-model.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md)

**Berührte Spec-Stellen:** `LH-FA-11.a` · `LH-FA-02.b` · `LH-FA-05.a` · `LH-FA-05.d` · `LH-FA-05.e` · `LH-FA-09.a` · `LH-FA-18.a` · `SPEC-002` · `SPEC-003` · `SPEC-032` · `SPEC-039`

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Abnahmeszenario 6 ist für Simple Query end-to-end nachgewiesen, und die Randformen der Fehler- und Hinweisantworten stehen in der Spezifikation: Fehlerantworten, mehrzeilige und leere Ergebnismengen, Befehle ohne Zeilen, mehrere Ergebnisse einer Anfrage, `NoticeResponse`, `ParameterStatus` nach `SET` und der Transaktionsstatus im `ReadyForQuery` werden aufgezeichnet und in der aufgezeichneten Reihenfolge reproduziert.

**Der Slice ist im Kern ein Test-Nachweis** (Zuschnitt vom Nutzer am 2026-10-05 bestätigt). Record und Replay führen jede Serverantwort einer Interaktion schon als geordnete Liste (unten, „Bereits geliefert“). Neu als Vertrag sind drei Festlegungen der Spezifikation, die bestehendes Verhalten festschreiben (§6, Zeilen 1, 9 und 10): Weitergabe an den Client (`LH-FA-02.b`), Nachrichten zwischen Interaktionen (`LH-FA-05.a`), Diagnosefelder (`LH-FA-11.a`). Dazu kommt eine Festlegung mit Produktionscode (§6, Zeile 11, Nutzerentscheid vom 2026-10-05): Eine Fehlerantwort, nach der der Upstream ohne `ReadyForQuery` abbricht, ist `PGR-E6001` statt `PGR-E4003` (`LH-FA-02.b` §Fehlerantwort vor dem Abbruch). Das ändert den Upstream-Adapter und den Record-Service. Jede Festlegung bekommt ihren Test mit Mutation (`AGENTS.md` §3.10). Zeigt ein Test eine Abweichung von der Spezifikation, ist die Korrektur im betroffenen Adapter Teil dieses Slice und wird in §3 nachgetragen.

**Bereits geliefert** — mit Test:

- Geordnete Liste der Serverantworten je einfacher Interaktion, alle Typen aus `LH-FA-05.a` im Domain Model und in beiden PGWire-Adaptern abgebildet: `TestOpenUndQuery`, `TestReplayStrictSequential`, `TestE2EReplaySelect1` (Spalten, Typ-OIDs, Zeilen, Befehlsabschluss), `TestE2ERecordMehrereInteraktionen` (Reihenfolge, NULL).
- `ErrorResponse` im Extended Query Protocol samt Verwerfen bis `Sync`: `TestReplayExtendedFehlerantwort`, `TestE2EReplayExtendedPgx` (deklariert `LH-FA-11/Happy`).
- `ErrorResponse` einer einfachen Anfrage in einer Transaktion samt `ReadyForQuery` im Status `T`, `E` und `I` (`BEGIN`, `SELECT 1/0`, `ROLLBACK`), Record und Replay gleich PostgreSQL: `TestE2EReplayLebendpruefungWiePostgres` — deklariert nur für die Lebendprüfung; dieser Slice deklariert den Fall für `LH-FA-11` in einem eigenen Ablauf.
- Binärwerte als `base64` und ihr Roundtrip (`SPEC-003`): `TestRoundtrip`; Spalten im Binärformat über pgx: `TestE2EReplayExtendedPgx`.
- `COPY` nicht unterstützt (`LH-FA-05.e`): `TestE2ERecordNichtUnterstuetzt`; nicht unterstützte Serverantwort: `TestNichtVermittelbar`, `TestRecordSessionNichtUnterstuetzt`.
- Leere einfache Anfrage (`EmptyQueryResponse`) als Lebendprüfung (`LH-FA-09.a`): `TestE2EReplayLebendpruefungWiePostgres`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Meldungscodes — `slice-replay-semantik-meldungscodes`.
- Änderungen an Domain Model, Recording-Format und Abbildung in den Adaptern — Bestand bleibt bewusst stehen: alle Antworttypen aus `LH-FA-05.a`, `Value` mit Bytes und `base64` sind geliefert (oben).
- Streaming großer Ergebnismengen einfacher Anfragen — Bestand bleibt bewusst stehen: `SPEC-032` ist ein Ziel ohne Zusicherung, [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md) lässt einfache Anfragen auf ihrem synchronen Weg; `LH-FA-02.b` hält die Grenze fest.
- `COPY`, `FunctionCall`, `NotificationResponse` als unterstützte Interaktion — Bestand bleibt bewusst stehen: nicht unterstützt nach `LH-FA-05.e` und `SPEC-039`, geprüft (oben).
- Fehlerantwort des Upstreams im Verbindungsaufbau beim Record (`BEO-REPO/record-fehlerantwort-im-aufbau-ungeregelt`) — anderer Vorgang: Der Aufbau ist keine Interaktion, und `LH-FA-11` betrifft Fehlerantworten auf Anfragen.
- Mehrere Sessions (Boundary von `LH-FA-12`) — `slice-v1-abschluss-sessions`.
- Verhalten gegen andere PostgreSQL-Versionen als die gepinnte — `slice-v1-abschluss-postgres-versionen`; die Abläufe hier vergleichen mit derselben Instanz, gegen die sie aufzeichnen.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-11`](../../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers): Abnahmeszenario 6 für Simple Query — eine einzelne Fehlerantwort, ein Fehler in der zweiten Anweisung einer Anfrage nach dem Ergebnis der ersten und ein Fehler in einer Transaktion mit `ReadyForQuery` im Status `E` und anschließendem `ROLLBACK` zeigen im Replay dieselbe Sicht des Clients (Ergebnisse, SQLSTATE, Meldung, Transaktionsstatus) wie beim Aufzeichnen und wie direkt gegen PostgreSQL; Abdeckung `LH-FA-11/Happy`, `LH-FA-11/Boundary` (E2E).
- [ ] [`LH-FA-06`](../../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen): Simple Query — mehrzeilige und leere Ergebnismenge, Befehle ohne Zeilen (`CREATE TABLE`, `INSERT`, `UPDATE`), mehrere Ergebnismengen einer Anfrage, zwei `NoticeResponse` in ihrer Reihenfolge, `ParameterStatus` nach `SET` und ein Binärwert über einen Binär-Cursor (`SPEC-003`) erscheinen im Replay in Reihenfolge und Inhalt wie beim Aufzeichnen; Abdeckung `LH-FA-06/Happy`, `LH-FA-06/Boundary`, `LH-FA-12/Happy` (E2E).
- [ ] Die drei Festlegungen aus §6 (Zeilen 1, 9, 10) mit je einem Test, der unter ihrer Mutation rot wird: Weitergabe (`LH-FA-02.b`) — `SELECT 1; COPY (SELECT 1) TO STDOUT` liefert dem Client kein Ergebnis, nur `PGR-E6001` (E2E); Nachrichten zwischen Interaktionen (`LH-FA-05.a`) — Unit-Test am Upstream-Adapter mit Fake-Server; Diagnosefelder (`LH-FA-11.a`) — Unit-Test am Upstream-Adapter mit rohen Bytes einer `ErrorResponse` mit leerem Feld und `P` = `0`; Fehlerantwort vor dem Abbruch (`LH-FA-02.b`) — `PGR-E6001` mit dem SQLSTATE des Servers statt `PGR-E4003`, Session nicht übernommen, auch im Extended Query Protocol (Unit-Tests am Upstream-Adapter mit Fake-Server und am Record-Service, Vorgabe in §6).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update (Architect, vor dem Code) | `LH-FA-02.b` Weitergabe an den Client, `LH-FA-05.a` Nachrichten zwischen Interaktionen, `LH-FA-11.a` Diagnosefelder — liegt im Commit des Architect |
| `test/integration/einfach_e2e_test.go` | add | DoD 1 und 2 sowie die Weitergabe aus DoD 3: Ablauf direkt gegen PostgreSQL, über record und über replay; die Sicht des Clients sind die Server-Nachrichten nach dem Verbindungsaufbau in Reihenfolge, mit vollständigem Inhalt, Felder von Fehler- und Hinweisantworten nach Feldcode sortiert; der Test liest PGWire selbst (`rohAblauf`), weil die PGWire-Bibliothek nach `.a-check.yml` den Adaptern gehört; mit Abdeckungs-Deklarationen |
| `internal/adapters/driven/postgres/upstream.go` | update | DoD 3: die Session merkt sich die letzte `ErrorResponse` seit dem letzten `ReadyForQuery`; endet die Verbindung davor, liefern `Query` und `Receive` `PGR-E6001` mit SQLSTATE und Meldung des Servers statt `PGR-E4003` |
| `internal/adapters/driven/postgres/upstream_test.go` | update | DoD 3: Nachrichten zwischen Interaktionen, Diagnosefelder und Fehlerantwort vor dem Abbruch (Tests a bis f) am Fake-Server; das Lesen der Session ist dort auf fünf Sekunden begrenzt |
| `internal/hexagon/services/record.go` | update | DoD 3: `AwaitServer` markiert die Session bei `PGR-E6001` aus `Receive` als nicht übernehmbar, wie `Query` es schon tut (`LH-FA-05.a`); beide über dieselbe Hilfsfunktion |
| `internal/hexagon/services/record_extended_test.go` | update | DoD 3: `Receive` mit `PGR-E6001` nach einer übernommenen Interaktion — die Session wird nicht übernommen (Test g) |
| `internal/hexagon/services/record_test.go` | update | Fake-Session: `Receive` liefert auf Wunsch einen Fehler (für Test g) |
| `docs/plan/planning/open/slice-v1-abschluss-einspielen.md`, `docs/plan/planning/open/slice-v1-abschluss-antwortvergleich.md` | update | §1 beziehungsweise §6 vermerken die seit Zeile 11e unerreichbaren Teile von `LH-FA-20.a` und `LH-FA-24.a` (Risiko 2) |
| `docs/user/abdeckung-*.md` | update | von `make abdeckung` aus den neuen Deklarationen geschrieben, darunter `LH-FA-11/Negative` |
| weiterer Produktionscode | — | keiner; kein Test zeigte eine Abweichung von der Spezifikation |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-extended-query` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: der Slice verlangt mehr als drei Liefer-Punkte — zurück zur Zerlegung.
- `in-progress` → `open`: Die PGWire-Bibliothek repräsentiert eine Nachricht nicht verlustfrei — Carveout.


## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — je Randform die Entscheidung und ihr Ort:

| # | Randform | Entscheidung | Ort | Nachweis |
|---|---|---|---|---|
| 1 | `NoticeResponse` oder `ParameterStatus` asynchron, zwischen zwei Interaktionen | gehört zur nächsten Interaktion, vor deren übrigen Antworten; nach der letzten Interaktion nicht aufgezeichnet; das Replay gibt sie mit der Interaktion wieder, bei einer Lebendprüfung nicht | `LH-FA-05.a` §Nachrichten zwischen Interaktionen | Unit, Upstream-Adapter (DoD 3) |
| 2 | `ParameterStatus` mitten in einer Session, etwa nach `SET` | Teil der Serverantworten der Interaktion, in Reihenfolge wiedergegeben | `LH-FA-05.a` §Serverantworten | E2E (DoD 2) |
| 3 | `ReadyForQuery` mit Status `E` nach einem Fehler in einer Transaktion | aufgezeichnete Antwort, keine eigene Transaktionslogik im Replay | `LH-FA-05.d`, `LH-FA-11.a` | E2E (DoD 1); Bestand in `TestE2EReplayLebendpruefungWiePostgres` |
| 4 | `COPY` | nicht unterstützt, `PGR-E6001` | `LH-FA-05.e`, `SPEC-039` | `TestE2ERecordNichtUnterstuetzt` (Bestand) |
| 5 | leere Ergebnismenge; `EmptyQueryResponse` | leere Ergebnismenge ohne Sonderform (`RowDescription`, `CommandComplete`); leere einfache Anfrage ist eine Lebendprüfung, in einer Extended-Interaktion Teil der Gruppe | `LH-FA-05.a`, `LH-FA-09.a`, `LH-FA-18.a` | E2E (DoD 2); Lebendprüfung Bestand |
| 6 | sehr große Ergebnismenge einer einfachen Anfrage | bis zum `ReadyForQuery` im Speicher, keine Zusage früherer Ankunft beim Client | `LH-FA-02.b` §Weitergabe, `SPEC-032` | kein Test — akzeptiertes Negativ: Die Session liegt bis zum Schreiben ohnehin im Speicher (`SPEC-032`), ein Puffer bis `ReadyForQuery` vergrößert den Bedarf nicht |
| 7 | Binärformat der Spalten | Bytes, im YAML als `base64`, in SQLite als Wire-Bytes | `SPEC-003`, `SPEC-043` | `TestRoundtrip` (Bestand); E2E mit Binär-Cursor (DoD 2) |
| 8 | mehrere Anweisungen in einer Anfrage, Fehler in der zweiten | eine Interaktion; Ergebnis der ersten, `ErrorResponse`, `ReadyForQuery` in Reihenfolge | `LH-FA-05.a`, `LH-FA-02.b`, `LH-FA-11.a` | E2E (DoD 1) |
| 9 | Antworten auf dem Weg zum Client beim Aufzeichnen, auch `NoticeResponse`; nicht unterstützte Antwort nach Ergebnissen | Reihenfolge des Upstreams; ist eine Antwort nicht unterstützt, erhält der Client keine der Interaktion, nur `PGR-E6001` | `LH-FA-02.b` §Weitergabe | E2E (DoD 2, DoD 3) |
| 10 | Diagnosefeld mit leerem Wert, Zahlenfeld `0` oder ohne Zahl, Reihenfolge der Felder | bei den 18 Feldcodes, die die Bibliothek benennt, gilt ein leerer Wert als fehlend, ebenso ein Zahlenfeld `0` oder ohne Zahl — weder aufgezeichnet noch gesendet; ein anderer Feldcode wird auch mit leerem Wert aufgezeichnet und gesendet (Rückgabe des Implementers, 2026-10-05); Reihenfolge nicht aufgezeichnet. Bei den benannten Codes unterscheidet die Bibliothek leer und fehlend nicht; das ist eine bisher nicht dokumentierte Lockerung der Folge aus [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md) („nicht verlustfrei heißt nicht unterstützt“), die nicht erkennbare Fälle betrifft und mit dieser Zeile dokumentiert ist | `LH-FA-11.a` §Diagnosefelder | Unit, Upstream-Adapter (DoD 3): `TestDiagnosefelder` prüft zusätzlich einen unbekannten Code mit leerem Wert, der als `""` in `Fields` steht; Mutation: leere Werte im Zweig der unbekannten Codes von `noticeFields` überspringen → rot |
| 11 | `ErrorResponse`, nach der der Upstream die Verbindung ohne `ReadyForQuery` beendet | nicht unterstützt, `PGR-E6001`, Exit-Code 6 (Nutzerentscheid 2026-10-05); ohne `ErrorResponse` davor bleibt es `PGR-E4003` | `LH-FA-02.b` §Fehlerantwort vor dem Abbruch, `LH-FA-11.a` | Unit (DoD 3, Vorgabe unten) |
| 11a | Schweregrad `FATAL` gegenüber `ERROR` (auch `PANIC`) | zählt nicht, alle gleich | `LH-FA-02.b` §Fehlerantwort vor dem Abbruch, *Schweregrad* | Test b |
| 11b | `ErrorResponse` mitten in der Antwortliste, nach Ergebnissen; als erste Antwort aus der Zeit zwischen zwei Interaktionen; `ErrorResponse` mit folgendem `ReadyForQuery` vor einer späteren Interaktion | jede Stelle zählt; zwischen Interaktionen gehört sie zur nächsten; nach einem `ReadyForQuery` zählt sie für spätere Interaktionen nicht; ohne weitere Anfrage endet die Session regulär | *Stelle* | Tests a, e, f |
| 11c | Was der Client vor dem Abbruch erhält | einfache Anfrage: keine ihrer Antworten, nur `PGR-E6001`; Extended: schon weitergegebene Server-Nachrichten bleiben beim Client, danach `PGR-E6001`, soweit die Verbindung sie annimmt; der Meldungstext nennt SQLSTATE und Meldung der letzten `ErrorResponse` | *Client*, *Diagnose* | Tests a, d |
| 11d | Was aufgezeichnet wird | nichts aus der Session, auch nicht ihre vorherigen Interaktionen (anders als bei `PGR-E4003`) | *Aufzeichnung*, `LH-FA-05.a` | Test g |
| 11f | Upstream sendet zwischen zwei Interaktionen eine `ErrorResponse` und schließt mit TCP-Reset; das Senden der nächsten Anfrage scheitert | `PGR-E4003`: maßgeblich ist eine gelesene `ErrorResponse`, nach einem Sendefehler liest der Recorder nicht weiter (Rückgabe des Implementers, 2026-10-05) | *Fehlerantwort vor dem Abbruch*, erster Absatz | keine Code-Änderung, kein Test — akzeptiertes Negativ: Nach einem Reset verwirft der Kernel ungelesene Empfangsdaten in der Regel, ein Lesen nach dem Sendefehler fände die `ErrorResponse` nicht verlässlich und machte das Ergebnis vom Zeitpunkt abhängig; PostgreSQL schließt nach einem `FATAL` geordnet, dann gilt Test f (`PGR-E6001`) |
| 11e | Replay einer solchen Interaktion | kommt nicht vor; eine Interaktion, die nicht mit `ready_for_query` endet, macht die Aufzeichnung beschädigt (`PGR-E3003`) | *Replay* | Bestand (`model.Interaction.Validate` beim Lesen) |

**Vorgabe für Zeile 11 an den Implementer** — Verhalten, Test, Mutation (`AGENTS.md` §3.10):

- *Verhalten.* Der Upstream-Adapter merkt sich je Session die letzte `ErrorResponse` seit dem letzten `ReadyForQuery` und vergisst sie mit jedem `ReadyForQuery`. Endet die Verbindung in `Query` oder `Receive` mit gemerkter `ErrorResponse`, liefert er `PGR-E6001` (Klasse nicht unterstützt) mit SQLSTATE und Meldung im Text, sonst wie bisher `PGR-E4003`. `RecordService.AwaitServer` setzt bei `PGR-E6001` die Session auf nicht übernehmbar, wie `Query` es tut.
- *Tests* (Fake-Server mit `pgproto3.Backend` in `upstream_test.go`, außer g):
  a. Einfache Anfrage: `RowDescription`, `DataRow`, `CommandComplete`, `ErrorResponse` mit `FATAL`/`57P01`, dann Schließen — `Query` liefert `PGR-E6001`, der Text enthält `57P01`.
  b. Wie a, Schweregrad `ERROR` — ebenfalls `PGR-E6001`.
  c. Schließen ohne `ErrorResponse` — `PGR-E4003`.
  d. Extended: `Receive` liefert die `ErrorResponse`; der nächste `Receive` nach dem Schließen liefert `PGR-E6001` mit `57P01`.
  e. Interaktion 1 mit `ErrorResponse` und `ReadyForQuery`, Interaktion 2 ohne Antwort geschlossen — `PGR-E4003`.
  f. `ErrorResponse` vor der Anfrage gesendet (zwischen Interaktionen), dann Schließen — `Query` liefert `PGR-E6001`.
  g. `record_extended_test.go`: Fake-`Receive` mit `PGR-E6001` nach einer übernommenen Interaktion — die Session steht nicht in der Aufzeichnung.
- *Mutationen*, je selbst rot gesehen: Einordnung entfernt (immer `PGR-E4003`) → a, b, d, f rot · Prüfung auf `FATAL` beschränkt → b rot · Merker beim `ReadyForQuery` nicht zurückgesetzt → e rot · SQLSTATE nicht im Text → a, d rot · Markierung in `AwaitServer` entfernt → g rot.

**Akzeptierte Negative** — entschieden, kein eigener Folgeauftrag:

- Eine leere Meldung (`M`), etwa aus `RAISE EXCEPTION ''`, erreicht den Client ohne das Pflichtfeld `M`, im Record- wie im Replay-Modus (Zeile 10): Die Bibliothek sendet leere Felder nicht. pgx liest das als leere Meldung; ein eigener Encoder wäre gegen [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md).
- Unbekannte Feldcodes sendet die Bibliothek in beliebiger Reihenfolge; PostgreSQL sendet nur bekannte Codes, und die Reihenfolge ist nicht Teil der Aufzeichnung (Zeile 10).

**Risiken:**

- Nachrichten, die die Bibliothek nicht verlustfrei abbildet, sind als nicht unterstützt zu klassifizieren — für Diagnosefelder in `LH-FA-11.a` entschieden (Zeile 10), für Nachrichtentypen Bestand (`PGR-E6001`) — **Ausgang:** offen bis Closure.
- `LH-FA-20.a` Schritt 3 und `LH-FA-24.a` §Unvollständige Aufzeichnung beschreiben Einspielen und Vergleich für aufgezeichnete Interaktionen ohne `ReadyForQuery`; nach Zeile 11e schreibt der Recorder keine, und der Leser lehnt sie ab. Die Zweige sind damit unerreichbar und gehören bereinigt, wenn `play` entsteht (`slice-v1-abschluss-einspielen` §1, `slice-v1-abschluss-antwortvergleich` §6 vermerken es); dieser Slice ändert sie nicht — **Ausgang:** offen bis Closure.
- Ein Ablauf vergleicht Meldungstexte des Servers; ändern sie sich zwischen Läufen (etwa mit der Sprache oder einem Zeitstempel in der Meldung), wird der Vergleich instabil — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand). Zwei offene Einträge berühren den Gegenstand, je 1×: `BEO-REPO/record-fehlerantwort-im-aufbau-ungeregelt` — der Aufbau ist hier ausgeschlossen (§1), der Zähler steigt durch diesen Slice nicht; `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` — die Abläufe vergleichen gegen die gepinnte Instanz (§1), der Zähler steigt nur, wenn ein Review hier eine Versionsabhängigkeit findet.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Repo-Default, `harness/conventions.md`).
