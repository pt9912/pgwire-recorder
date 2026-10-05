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

**Der Slice ist im Kern ein Test-Nachweis.** Record und Replay führen jede Serverantwort einer Interaktion schon als geordnete Liste (unten, „Bereits geliefert“); Produktionscode wird nicht erwartet. Neu als Vertrag sind nur drei Festlegungen der Spezifikation, die bestehendes Verhalten festschreiben (§6, Zeilen 1, 9 und 10): Weitergabe an den Client (`LH-FA-02.b`), Nachrichten zwischen Interaktionen (`LH-FA-05.a`), Diagnosefelder (`LH-FA-11.a`). Jede bekommt ihren Test mit Mutation (`AGENTS.md` §3.10). Zeigt ein Test eine Abweichung von der Spezifikation, ist die Korrektur im betroffenen Adapter Teil dieses Slice und wird in §3 nachgetragen.

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
- [ ] Die drei Festlegungen aus §6 (Zeilen 1, 9, 10) mit je einem Test, der unter ihrer Mutation rot wird: Weitergabe (`LH-FA-02.b`) — `SELECT 1; COPY (SELECT 1) TO STDOUT` liefert dem Client kein Ergebnis, nur `PGR-E6001` (E2E); Nachrichten zwischen Interaktionen (`LH-FA-05.a`) — Unit-Test am Upstream-Adapter mit Fake-Server; Diagnosefelder (`LH-FA-11.a`) — Unit-Test am Upstream-Adapter mit rohen Bytes einer `ErrorResponse` mit leerem Feld und `P` = `0`.
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
| `test/integration/` (neue Datei für einfache Anfragen) | add | DoD 1 und 2 sowie die Weitergabe aus DoD 3: Ablauf über record, replay und direkt gegen PostgreSQL, Sicht des Clients als Text verglichen (Muster `extendedAblauf`), mit Abdeckungs-Deklarationen |
| `internal/adapters/driven/postgres/upstream_test.go` | update | DoD 3: Nachrichten zwischen Interaktionen und Diagnosefelder am Fake-Server |
| `docs/user/abdeckung-*.md` | update | von `make abdeckung` aus den neuen Deklarationen geschrieben |
| Produktionscode | — | keiner erwartet; zeigt ein Test eine Abweichung von der Spezifikation, wird die Korrektur hier mit Datei nachgetragen |

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
| 10 | Diagnosefeld mit leerem Wert, Zahlenfeld `0` oder ohne Zahl, Reihenfolge der Felder | gilt als fehlend, wird weder aufgezeichnet noch gesendet; Reihenfolge nicht aufgezeichnet. Die Bibliothek unterscheidet beides nicht; das ist eine bisher nicht dokumentierte Lockerung der Folge aus [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md) („nicht verlustfrei heißt nicht unterstützt“), die nicht erkennbare Fälle betrifft und mit dieser Zeile dokumentiert ist | `LH-FA-11.a` §Diagnosefelder | Unit, Upstream-Adapter (DoD 3) |
| 11 | `ErrorResponse` (etwa `FATAL`), nach der der Upstream die Verbindung ohne `ReadyForQuery` beendet | **offen — Nutzerentscheid vor dem Code.** Heute gilt `LH-FA-02.b`: unerwartetes Verbindungsende, `PGR-E4003`, Interaktion nicht übernommen. Alternative nach der Negative-Zeile von `LH-FA-11`: nicht unterstützt, `PGR-E6001`. Bis zur Entscheidung berührt der Implementer den Fall nicht | `LH-FA-02.b` (heute) | nach Entscheid in DoD 3 |

**Akzeptierte Negative** — entschieden, kein eigener Folgeauftrag:

- Eine leere Meldung (`M`), etwa aus `RAISE EXCEPTION ''`, erreicht den Client ohne das Pflichtfeld `M`, im Record- wie im Replay-Modus (Zeile 10): Die Bibliothek sendet leere Felder nicht. pgx liest das als leere Meldung; ein eigener Encoder wäre gegen [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md).
- Unbekannte Feldcodes sendet die Bibliothek in beliebiger Reihenfolge; PostgreSQL sendet nur bekannte Codes, und die Reihenfolge ist nicht Teil der Aufzeichnung (Zeile 10).

**Risiken:**

- Nachrichten, die die Bibliothek nicht verlustfrei abbildet, sind als nicht unterstützt zu klassifizieren — für Diagnosefelder in `LH-FA-11.a` entschieden (Zeile 10), für Nachrichtentypen Bestand (`PGR-E6001`) — **Ausgang:** offen bis Closure.
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
