# Slice slice-v1-abschluss-verbindungsende-upstream-extended: Das Ende einer Extended-Interaktion gegen PostgreSQL belegen (Fixtures der Gruppen mit Sync und mit Flush)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag (Entscheidung des Nutzers vom 2026-10-18; Schritt 0 der Reihenfolge in [ihrem Plan](../welle-v1-nachtrag.md)). Der Gegenstand liefert Tests und Testdaten, keinen Produktcode und keine Anforderung.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-18.a` · `LH-FA-03.a`

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-16; nach `slice-harness-regeln-rueckbau` nachgezogen am 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Fixtures `ende-extended-sync.yaml` und `ende-extended-flush.yaml` in `test/integration/testdata` sind von Hand geschrieben und bilden ab, was die Spezifikation sagt (`LH-FA-25.a` *Aufzeichnen*, *Stelle*, `LH-FA-18.a`). Dieser Slice belegt gegen einen echten PostgreSQL (Referenzversion 17), dass der Server die Nachrichtenfolge liefert, die sie zeigen: Ein Roh-Client mit Extended Query, dessen Gruppe mit `Sync` und dessen Gruppe mit `Flush` vom Server mit `FATAL` (`57P01`) beendet wird, wird mit `record` aufgezeichnet; die Aufzeichnung (YAML, Version 2) hat in der Gruppe mit Ende dieselbe Form wie die Fixture (Server-Nachrichten bis einschließlich der Fehlerantwort, danach `connection_closed`, kein `ready_for_query`), und `replay` gibt die echte Aufzeichnung ohne PostgreSQL wieder. Weicht die Form des Servers von einer Fixture ab, wird die Fixture berichtigt und mit ihr die Tests von `slice-v1-abschluss-verbindungsende-replay-extended` und `-replay-extended-flush`, soweit sie die Form prüfen; ändert sich dazu Produktcode, ist das ein Befund an den Planner (§4). Der Slice liefert Tests und Testdaten, keinen Produktcode.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-16 („ja“ zu den Vorschlägen des Planners) im Register `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`: Die Extended-Hälfte der Beobachtung bekommt einen eigenen Slice statt eines vierten Extended-Falls in `slice-v1-abschluss-verbindungsende-upstream` (`BEO-REPO/slice-waechst-durch-uebernahmen`). Der Ausgang des Registers ist `geplant` mit dieser Kennung; die Adresse nimmt an, solange der Slice nicht in `done/` liegt.

**Übernommen aus `slice-v1-abschluss-verbindungsende-replay-extended`** (dort §6, Risiko, und Register, Zähler 2; Sendung des Planners, 2026-10-16, `AGENTS.md` §3.13): dass die handgeschriebene Extended-Aufzeichnung mit einer Gruppe mit `Sync` (`ende-extended-sync.yaml`) dem entspricht, was PostgreSQL liefert. **Übernommen aus `slice-v1-abschluss-verbindungsende-replay-extended-flush`** (dort §6, Risiko, und Register, Zähler 3; Sendung des Planners, 2026-10-16): dasselbe für die Gruppe mit `Flush` (`ende-extended-flush.yaml`). Beide Geber liegen in `done/`; die Sendung steht in ihren Risiko-Ausgängen und im Register, in diesem Slice als Adresse.

**Voraussetzung, keine Lieferung — `slice-v1-abschluss-verbindungsende-upstream`** (in `done/`): Der Upstream-Adapter erkennt das Ende in `Send` und `Receive` einer Extended-Interaktion (dort DoD-Punkt 1), und das E2E gegen PostgreSQL für eine einfache Anfrage (Abnahmeszenario 18, `TestE2ERecordEndeTerminate` in `test/integration/ende_e2e_test.go`) liegt dort. Dieser Slice nutzt deren Hilfen (`gleicheForm`, `endeFelder`, `felderDerFehlerantwort` in `ende_e2e_test.go`; `rohNachricht`, `parseSync` in `replay_ende_extended_e2e_test.go`) und stellt keine Anforderung an sie.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Upstream-Adapter, der Record-Service, der Replay-Service und der PGWire-Adapter (Produktcode) — `slice-v1-abschluss-verbindungsende-upstream`, `-record`, `-replay-extended` und `-replay-extended-flush` (alle in `done/`): Schicht-Abgrenzung; der Slice belegt, er baut nicht. Zeigt der Server ein Verhalten, das Produktcode ändern müsste, geht der Befund an den Planner (§4).
- Die einfache Anfrage und die Fixtures `ende-fatal.yaml`, `ende-lebend.yaml`, `ende-panic.yaml`; Ende im Leerlauf, Netzwerkabbruch ohne Fehlerantwort, `ERROR` mit anschließendem Abbruch — `slice-v1-abschluss-verbindungsende-upstream` (DoD-Punkte 2 und 3, Sendung aus `-replay`): Der Beleg dort führt die Entsprechung für die einfache Interaktion; diese Hälfte der Beobachtung hat dort ihre Adresse.
- Die Aufzeichnung in SQLite — der Beleg für `meta.version` 2 gegen echtes Ende liegt in `slice-v1-abschluss-verbindungsende-upstream` (DoD-Punkt 3); das Format ist von der Art der Interaktion unabhängig, und beide Schreiber sind in `-format-schreiber` und `-sqlite-schreiber` geprüft. Bestand bleibt bewusst stehen.
- `play-ende-extended-sync.yaml` und `play-ende-extended-flush.yaml` (Eingaben von `play`, aus `slice-v1-abschluss-verbindungsende-einspielen-extended`, Zähler 4 im Register) — Bestand bleibt bewusst stehen (Klasse 2). Grund: Beide sind handgeschrieben und weichen vom Server ab (für `execute` mit Ende keine `data_row`, ein anderer Gruppenschnitt als die echte Aufzeichnung); `play` liest ohne `--compare-responses` keine aufgezeichneten Antworten, die Abweichung ist harmlos, solange kein Test sie vergleichend einspielt. Ein Beleg gegen den Server wäre ein anderer Vorgang (Einspielen, nicht Aufzeichnen); §6 führt es als Risiko mit Ausgang.
- Andere PostgreSQL-Versionen — Szenario 18 läuft nur gegen die Referenz 17 (akzeptiertes Negativ aus `slice-v1-abschluss-verbindungsende-upstream` §6, Entscheidung des Nutzers vom 2026-10-12); die Matrix aus `slice-v1-abschluss-postgres-versionen-matrix` fährt Szenarien 1, 2 und 7. Der Beleg hier gilt für 17 und sagt nichts über andere.
- Eine Erweiterung von Abnahmeszenario 18 im Lastenheft auf Extended Query — anderer Vorgang: Eine Änderung des Lastenhefts ist ein Change Request (`harness/conventions.md` §Versionierung); die Tests tragen `LH-FA-25` und `LH-FA-18`, nicht das Szenario.
- Benutzerhandbuch und `README.md` — kein neues Verhalten (`AGENTS.md` §3.11): Das Handbuch sagt zum Ende einer Extended-Interaktion beim Aufzeichnen nichts (geprüft 2026-10-22 mit `grep`; `slice-v1-abschluss-verbindungsende-doku` liegt in `done/`), und dieser Slice liefert keinen Satz dazu. Ob nach seinem `done/` einer entsteht, ist eine offene Frage an den Nutzer (Bericht des Planners).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1** — [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Gruppe mit `Sync`. Ein Roh-Client mit Extended Query, dessen Gruppe mit Ende mit `Sync` schließt (Gruppen davor mit `Sync` und mit `Flush`, wie `ende-extended-sync.yaml`), wird am Binary gegen PostgreSQL 17 aufgezeichnet; die Aufzeichnung enthält in der Gruppe mit Ende die Server-Nachrichten bis einschließlich der Fehlerantwort `FATAL` `57P01` und danach `connection_closed`, kein `ready_for_query`, und trägt `version: 2`; die Typenfolge je Gruppe und Schweregrad und SQLSTATE der Fehlerantwort stimmen mit der Fixture überein (Vergleich im Test gegen die eingebettete Datei, nicht in Prosa; der Vergleich setzt die Folge nicht vorab, die Abweichung ist der Befund); `replay` gibt die echte Aufzeichnung ohne PostgreSQL wieder (dieselben Antworten, danach Ende; Exit-Code 0 mit `--fail-on-unconsumed`). Weicht der Server ab, ist die Fixture berichtigt und ihre Tests nachgezogen (Integrationstest, Abdeckungstabelle `abdeckung-e2e.md`).
- [ ] **Liefer-Punkt 2** — [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Gruppe mit `Flush`. Dasselbe für die Gruppe mit Ende, die mit `Flush` schließt (`ende-extended-flush.yaml`), mit der Besonderheit, dass auch die Gruppe davor mit `Flush` schließt: Der Server sendet keine `ReadyForQuery`; die echte Aufzeichnung entspricht der Fixture oder die Fixture ist berichtigt, und `replay` gibt sie wieder (Integrationstest, Abdeckungstabelle).
- [ ] Der Vergleich ist an beiden Fixtures belegt: Der Test benennt die Fixture, die er prüft; eine Fixture ohne Vergleich gibt es nicht (`AGENTS.md` §3.11: der Satz „entspricht dem Server“ gilt nur, wo ein Test ihn prüft). Ein Test je Zusage der Punkte 1 und 2, ohne Mutation (Tests und Testdaten, `AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` von `geplant` auf `gestrichen` (Begründung: der Beleg gegen den Server liegt vor, Verweis auf die Tests) oder, wo die Form abwich, `verkörpert`/neue Beobachtung mit Beleg.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/ende_extended_e2e_test.go` | neu | Fälle `Sync` und `Flush`: `record` gegen PostgreSQL 17 mit Roh-Client (`rohVerbindung`, `parseSync`), Vergleich der Aufzeichnung mit der eingebetteten Fixture (`gleicheTypenfolge`, dazu `gleicheForm` aus `ende_e2e_test.go`), `replay` der echten Aufzeichnung. Der Ablauf des Clients steht je Fall einmal (`endeSyncAblauf`, `endeFlushAblauf`) und wird beim Aufzeichnen und beim Wiedergeben gesendet. Die Wiedergabe ohne PostgreSQL läuft in der Phase ohne PostgreSQL des Runners (`tools/test/run-integration-tests.sh` wählt dort `^TestE2EOhnePostgres`; die Aufzeichnung liegt unter `PGR_DATEN`, wie bei `TestE2EOhnePostgresEnde`) — `LH-FA-25`, `LH-FA-18` |
| `test/integration/testdata/ende-extended-sync.yaml` und `replay_ende_extended_e2e_test.go` | update | Die Gruppe mit Ende der Fixture mit `Sync` ist die gemessene Form (§6); `TestE2EReplayEndeExtendedSync` sendet sie über `endeSyncAblauf` (neue Hilfen `sendeParse`, `parseBindExecuteSync`) und erwartet `12DE`. `ende-extended-flush.yaml` bleibt. Der Kommentar „gehört genau einem Test“ über `endeExtendedSyncYAML` und `endeExtendedFlushYAML` nennt beide Tests je Fixture beim Namen (`AGENTS.md` §3.7) |
| `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabelle aus den Deklarationen der Tests |

Kein Produktcode, keine Spezifikation, keine Gates.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Voraussetzungen liegen in `done/`: `slice-v1-abschluss-verbindungsende-upstream`, `-record`, `-replay-extended` und `-replay-extended-flush`, dazu [welle-erster-release](../done/welle-erster-release/welle-erster-release.md) und `slice-harness-regeln-rueckbau`. Der Nutzer hat den Start der Produkt-Slices am 2026-10-22 mit „leg los“ freigegeben (wörtlich, so vom Koordinator übermittelt). Der Architect prüft vor dem Code die Frage nach dem Weg, das Ende auszulösen (§6). WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik und Messbefehl: `AGENTS.md` §3.14; Pfade des Plans: `test docs/user`); die Basis ist der Commit des Übergangs nach `in-progress/`, sein Hash: `642f5f51ceed85ff06588e45e07065add06eecd6`, Zweig `impl-upstream-ext`. Dann hält der Implementer an und meldet dem Planner; Schnitt `Sync` | `Flush` (Zuschnitt: Planner).
- `in-progress` → `next`: Die echte Form des Servers verlangt eine Änderung an Produktcode (Upstream-Adapter, Record-, Replay-Service, PGWire-Adapter) — Befund an den Planner, zurück zur Zerlegung; ein Folge-Slice trägt die Änderung (§1: der Slice liefert keinen Produktcode).
- `in-progress` → `open`: Das Ende einer Extended-Gruppe ist gegen PostgreSQL nicht deterministisch herstellbar (die Antwort des Servers hängt von der Zeit ab, ein Fall ist flackernd) — der Nutzer entscheidet, ob die Fixture als handgeschrieben akzeptiert bleibt (Carveout, Register bleibt offen).

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

**Randformen.** Der Slice liefert keine neue Wire-, Datei- oder Formatgrenze, sondern belegt Zusagen der Spezifikation (`LH-FA-25.a` *Aufzeichnen*, *Stelle*, `LH-FA-18.a` *Interaktion*, *Abschluss*) gegen den Server; eine Randformen-Liste nach `AGENTS.md` §3.12 ist nicht verlangt. Der Architect hat den Weg, das Ende auszulösen, am 2026-10-22 vor dem ersten Code-Commit entschieden, nach einer Messung am Binary gegen PostgreSQL 17.11 (`Flush` 20 von 20 Läufen, `Sync` 20 von 20 und 15 weitere mit einzeln gesendeten Nachrichten; `replay` je einer echten Aufzeichnung mit `--fail-on-unconsumed`: Exit-Code 0). Der Ort ist dieser Absatz (Mechanik und Vergleich der Tests, keine Zusage des Produkts); die Spezifikation bleibt, die Messung widerlegt keine ihrer Aussagen.

- *Gruppe mit `Sync`: der Backend beendet sich selbst, keine zweite Verbindung.* Die Gruppe mit Ende ist `Parse` (`beendeSQL`), `Bind`, `Execute`, `Sync`; der Server antwortet `parse_complete`, `bind_complete`, `data_row`, die Fehlerantwort (`FATAL`, `57P01`), dann das Ende, kein `ready_for_query`. `Parse` und `Sync` allein sind kein Weg: Der Recorder sendet eine Gruppe erst mit ihrem `Sync` an den Server (`RecordService.ClientMessage`), der Server antwortet `ready_for_query`; die Form der Fixture (`parse_complete`, Fehlerantwort) entstünde nur, wenn das Signal zwischen den beiden Nachrichten einträfe, das ist nicht herstellbar. Deshalb wird `ende-extended-sync.yaml` berichtigt: Gruppe 2 der Interaktion 2 trägt die vier Client-Nachrichten und `parse_complete`, `bind_complete`, `data_row`, Fehlerantwort, `connection_closed` (`bind` und `execute` wie in `play-ende-extended-sync.yaml`, `data_row` mit `values: [{text: t}]`); die Gruppe davor (`Flush`) und die Interaktion 1 (`Sync`) bleiben. `TestE2EReplayEndeExtendedSync` sendet die berichtigte Gruppe und erwartet vor dem Ende drei Nachrichten und die Fehlerantwort.
- *Gruppe mit `Flush`: zweite Verbindung, ohne Wartezeit.* Der Client liest `parse_complete` der Gruppe mit Ende, erst dann beendet eine zweite Verbindung (`verbinde` auf `PGR_UPSTREAM`) den Backend, den sie in `pg_stat_activity` am Text der Anfrage dieser Gruppe findet (`count(pg_terminate_backend(pid))` muss 1 sein, sonst endet der Test; kein Schlafen, keine Wiederholung). Die Fehlerantwort ist dann die nächste Nachricht; die Zuordnung hängt nur vom Zeitverhalten ab, wenn der Client nicht wartet (`LH-FA-18.a`), und er wartet. `ende-extended-flush.yaml` stimmt mit dem Server überein und bleibt.
- *Vergleich.* Typenfolge ist die Folge der `type:`-Werte der Session 1 in Dateireihenfolge (Block- wie Flussschreibweise); die Fixture wird vor `- id: 2` abgeschnitten, der Lauf legt keine zweite Session an. Dazu `V`, `S`, `C`, `M` der Fehlerantwort über `gleicheForm`; `F`, `L`, `R`, Werte, SQL-Text, `startup` und `server_parameters` zählen nicht. Die Session 2 der Fixture mit `Sync` (einfache Anfrage ohne `row_description` und `data_row`) dient dem Replay-Test für die Zuordnung der nächsten Session, sagt nichts über das Ende und bleibt unverglichen (die einfache Anfrage ist nicht Gegenstand, §1).
- *Wiedergabe und Einbettung.* Die Wiedergabe der echten Aufzeichnungen läuft in Phase 2 des Runners: ein Test, dessen Name mit `TestE2EOhnePostgres` beginnt, liest sie aus `PGR_DATEN` wie `TestE2EOhnePostgresEnde`, fährt `replay --fail-on-unconsumed` (Exit-Code 0, keine Meldung) und vergleicht die Sicht des Clients mit der beim Aufzeichnen (`endeSicht`). Die neuen Tests nutzen `endeExtendedSyncYAML` und `endeExtendedFlushYAML` ohne zweite Einbettung; der Kommentar darüber nennt statt „gehört genau einem Test“ beide Tests beim Namen (Eingabe des Replay-Tests, Vorbild des Vergleichs im Record-Test).
- *Mutation, Plattform.* Keine Mutation (`AGENTS.md` §3.10): Der Slice liefert keine Zusage neu und ändert keine; der Vergleich bricht an der Fixture, das Brechen für die DoD zeigt der Verifier mit einer geänderten Zeile. Keine Zusage hängt vom Betriebssystem ab (Form der Nachrichten; die Tests laufen im Linux-Container): keine Grenze nach §3.12.
- *Akzeptiertes Negativ.* `play-ende-extended-sync.yaml` und `-flush.yaml` zeigen für `execute` mit Ende keine `data_row`, der Server sendet eine (gemessen). Beide sind Eingabe von `play` ohne `--compare-responses`, das die aufgezeichneten Antworten nicht liest (Zeile `3/0/0`); harmlos, solange kein Test sie mit `--compare-responses` einspielt. Nicht Gegenstand dieses Slice (§1).

**Risiken:**

- Die Form des Servers weicht von einer Fixture ab (etwa zusätzliche Nachrichten vor der Fehlerantwort); die Releases erschienen mit den handgeschriebenen Fixtures (`BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`, 6×) — **Ausgang:** offen bis Closure (erwartet: eingetreten und mit diesem Slice berichtigt, oder entfallen mit dem Beleg der Übereinstimmung; verlangt die Abweichung Produktcode: eingetreten → Folge-Slice mit ID; das Register geht in beiden Fällen auf `gestrichen` oder `verkörpert`, DoD).
- Das Verhalten des Servers ist nur gegen die Referenzversion 17 geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, Stand `verkörpert`; akzeptiertes Negativ wie in `slice-v1-abschluss-verbindungsende-upstream` §6) — **Ausgang:** weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`.
- Ein Fall hängt von der Zeit ab (Verhältnis von Signal des Servers und `Sync` des Clients) und flackert — **Ausgang:** offen bis Closure (§4: Rückführung nach `open/`, Entscheidung des Nutzers).
- `slice-v1-abschluss-verbindungsende-doku` beschreibt das Ende einer Extended-Interaktion gegen PostgreSQL, bevor dieser Slice es belegt — **Ausgang:** entfallen, weil `-doku` in `done/` liegt und das Handbuch dazu nichts sagt (`grep` vom 2026-10-22, §1); eine neue Aussage im Handbuch bräuchte einen eigenen Slice.
- `play-ende-extended-sync.yaml` und `-flush.yaml` bleiben handgeschriebene Eingaben, die vom Server abweichen (akzeptiertes Negativ oben, §1); ein Test, der sie mit `--compare-responses` einspielte, würde rot (Review F-1, Verifikation V-1) — **Ausgang:** weiter offen → `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (Zähler 4).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
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

### Belege des Implementers

Beginn: `date` am 2026-10-22 18:22:17 CEST (Zweig `impl-upstream-ext`, Basis `642f5f51`); Ende der Läufe 18:39 CEST. Server: PostgreSQL 17.11 (`POSTGRES_IMAGE` aus `harness/mk/integration.mk`).

**Beleg je DoD-Punkt**

| DoD | Befehl | Ausgabe |
|---|---|---|
| Liefer-Punkt 1 (Sync) und 2 (Flush), Vergleich an beiden Fixtures | `make test-integration` (voll) | `--- PASS:` für `TestE2ERecordEndeExtendedSync`, `TestE2ERecordEndeExtendedFlush`, `TestE2EReplayEndeExtendedSync`, `TestE2EReplayEndeExtendedFlush` (Phase 1) und `TestE2EOhnePostgresEndeExtendedSync`, `TestE2EOhnePostgresEndeExtendedFlush` (Phase 2); `run-integration-tests: gruen`, `EXIT=0 DAUER=193s`, 281 PASS, kein FAIL |
| dasselbe, engster Lauf | `PGR_TESTS_1="TestE2ERecordEndeExtendedSync TestE2ERecordEndeExtendedFlush TestE2EReplayEndeExtendedSync TestE2EReplayEndeExtendedFlush" PGR_TESTS_2="TestE2EOhnePostgresEndeExtendedSync TestE2EOhnePostgresEndeExtendedFlush" bash tools/test/run-integration-tests.sh` | die sechs `--- PASS:`, `run-integration-tests: gruen` |
| kein Flackern (§4, `in-progress` → `open`) | die zwei Aufzeichnungs-Tests mit `-test.count=30` gegen eine Instanz 17.11 | je 30 `--- PASS:`, kein FAIL |
| Test wird rot, wenn die Fixture falsch ist | je Runde die Fixtures geändert (`touch`), engster Lauf, danach zurückgesetzt | siehe unten |
| `make docs-check` | `make docs-check` | `d-check: 1404 Datei(en) geprüft, 0 Befund(e)` |
| `make kopf-check`, `make abdeckung-check` | beide | Exit 0 |
| Teile von `make gates` | `make lint` (46 s), `make test` (gofmt, `go vet -tags integration`, Unit-Tests, Abdeckung) | beide Exit 0; `make gates` fährt der Koordinator am Endstand |

**Fixtures falsch** (Lauf je Runde nach Änderung von `ende-extended-sync.yaml` und `ende-extended-flush.yaml`):

- Runde a: Sync-Fixture in der Fassung vor diesem Slice (Gruppe 2: `parse sync`, `parse_complete`, Fehlerantwort), Flush-Fixture mit `C: 57P02`. `TestE2ERecordEndeExtendedSync` rot: `Typenfolge der Session 1:` aufgezeichnet `… parse bind execute sync parse_complete bind_complete data_row error_response connection_closed`, Fixture `… parse sync parse_complete error_response connection_closed`. `TestE2ERecordEndeExtendedFlush` rot: `Feld C der Aufzeichnung "57P01", der Fixture "57P02"`. Dazu die zwei Replay-Tests rot: `TestE2EReplayEndeExtendedFlush` an `C57P02`; `TestE2EReplayEndeExtendedSync` mit `Verbindung endet nicht sauber: … connection reset by peer`, nicht mit `parse weicht in sql ab` (berichtigt nach Verifikation V-3: die Datei passt nicht zum Ablauf, die Stelle nennt der Aufzeichnungs-Test mit der Typenfolge).
- Runde b: Sync-Fixture ohne `data_row`, Flush-Fixture mit einem zusätzlichen `bind_complete` vor der Fehlerantwort. Beide Aufzeichnungs-Tests rot an der Typenfolge (`… bind_complete error_response …` fehlt `data_row` bzw. trägt ein `bind_complete` mehr), beide Replay-Tests rot (`12E` statt `12DE` bzw. `12E` statt `1E`).
- Zurückgesetzt: `diff` gegen die Sicherung ohne Unterschied; danach der Lauf grün (`EXIT=0`).

**Größe** (Befehl aus `AGENTS.md` §3.14, Pfade `test docs/user`, Basis `642f5f51`, Stand `e21c89a4`): 280 Zeilen (Marke 1500); davon `ende_extended_e2e_test.go` 202, `replay_ende_extended_e2e_test.go` 54, `abdeckung-e2e.md` 18, Fixture 6.

**Abweichungen vom Plan**

- Der Text der Parse-Nachricht der Gruppe mit Ende im Flush-Fall ist `SELECT 3 /* ende-extended-flush */` (`endeFlushSQL`), nicht `SELECT 3` der Fixture: die zweite Verbindung findet den Backend am Text (`query =` und `pid <> pg_backend_pid()`); der Text zählt im Vergleich nicht (§6). `TestE2EReplayEndeExtendedFlush` sendet weiter `SELECT 3`.
- In der Sync-Fixture trägt die Gruppe mit Ende `beendeSQL` statt `SELECT 3` (zur `data_row` mit `t` passend). Das ist keine Abweichung: §6 nennt `beendeSQL` als SQL der Gruppe (berichtigt nach Review F-2 und Verifikation V-2).
- Der Ablauf des Clients steht je Fall einmal (`endeSyncAblauf`, `endeFlushAblauf`); `TestE2EReplayEndeExtendedSync` nutzt ihn, damit Aufzeichnen und Wiedergeben dieselben Bytes senden (§3 nachgezogen).

**Grenzen des Vergleichs**

- Die Typenfolge enthält keine Gruppengrenzen (`client:`, `server:`), nur die `type:`-Werte (so §6): Eine Client-Nachricht, die über eine Grenze ohne Server-Nachricht dazwischen wandert, ändert die Folge nicht. In beiden Fixtures trägt jede Gruppe Server-Nachrichten; ein Test dafür fehlt. Die Interaktionsgrenzen enthält die Folge (`extended` steht in ihr; Verifikation V-4, Lauf i).
- Im Flush-Fall deckt die Sicht des Clients (`endeSicht`) nur die Fehlerantwort; die Bestätigungen vor ihr liest `bisTyp` (Typ `1`), den Rumpf (leer) vergleicht kein Test.
- Alle Läufe gegen 17.11; andere Versionen sind nicht geprüft (§1, §6 Risiko 2).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (6×, Auslöser; Ausgang `geplant` mit dieser Kennung) und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (§6; Stand `verkörpert`). Keine weiteren Treffer.

**Zählung** (Planner, 2026-10-22, `AGENTS.md` §3.13; Geber `slice-harness-regeln-rueckbau`, dessen Reihenfolge den Start in §4 und die Marke 1500 setzt; Randformen-Liste, Mutationspflicht, Schätzung und Meldeschwelle entfallen): **Liefer-Punkte: 2** (Gruppe mit `Sync`, Gruppe mit `Flush`; Vergleich, Wiedergabe und die Sendungen aus `-replay-extended` und `-replay-extended-flush` gehören zu jedem Punkt). **Schichten: eine** (die Integrationstests am Binary, `test/integration`), mit der Markierung der erzeugten Tabelle `docs/user/abdeckung-e2e.md` (Nutzer- und Wartungs-Doku, `AGENTS.md` §3.13) zwei. Kein Produktcode, keine Spezifikation.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
