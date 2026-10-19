# Slice slice-v1-abschluss-verbindungsende-upstream-extended: Das Ende einer Extended-Interaktion gegen PostgreSQL belegen (Fixtures der Gruppen mit Sync und mit Flush)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag (Entscheidung des Nutzers vom 2026-10-18; bis dahin ohne Welle). Bis zum 2026-10-16 Schritt 41b von welle-v1-abschluss; hinter das Release verschoben nach Entscheidung des Nutzers vom selben Tag (§4 *Start*, Drift-Log der Roadmap): Der Gegenstand liefert Tests und Testdaten, keinen Produktcode und keine Anforderung; er passte nicht in das Ziel von welle-v1-nachtrag (die zwei Anforderungsreste Homebrew und `CancelRequest` ohne Schlüssel); seit dem 2026-10-18 nennt [welle-v1-nachtrag](../welle-v1-nachtrag.md) ihn in §1, §4 und §5 als dritten Gegenstand (Nachtrag an Tests und Testdaten), auf Entscheidung des Nutzers.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-18.a` · `LH-FA-03.a`

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Fixtures `ende-extended-sync.yaml` und `ende-extended-flush.yaml` in `test/integration/testdata` sind von Hand geschrieben und bilden ab, was die Spezifikation sagt (`LH-FA-25.a` *Aufzeichnen*, *Stelle*, `LH-FA-18.a`). Dieser Slice belegt gegen einen echten PostgreSQL (Referenzversion 17), dass der Server die Nachrichtenfolge liefert, die sie zeigen: Ein Roh-Client mit Extended Query, dessen Gruppe mit `Sync` und dessen Gruppe mit `Flush` vom Server mit `FATAL` (`57P01`) beendet wird, wird mit `record` aufgezeichnet; die Aufzeichnung (YAML, Version 2) hat in der Gruppe mit Ende dieselbe Form wie die Fixture (Server-Nachrichten bis einschließlich der Fehlerantwort, danach `connection_closed`, kein `ready_for_query`), und `replay` gibt die echte Aufzeichnung ohne PostgreSQL wieder. Weicht die Form des Servers von einer Fixture ab, wird die Fixture berichtigt und mit ihr die Tests von `slice-v1-abschluss-verbindungsende-replay-extended` und `-replay-extended-flush`, soweit sie die Form prüfen; ändert sich dazu Produktcode, ist das ein Befund an den Planner (§4). Der Slice liefert Tests und Testdaten, keinen Produktcode.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-16 („ja“ zu den Vorschlägen des Planners) im Register `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`: Die Extended-Hälfte der Beobachtung (Zähler 2 und 3) bekommt einen eigenen Slice statt eines vierten Extended-Falls in `slice-v1-abschluss-verbindungsende-upstream`, dessen DoD-Punkt 3 eine einfache Anfrage fährt und dessen Größe (Schätzung 700 bis 1200) an der Grenze liegt (`BEO-REPO/slice-waechst-durch-uebernahmen`). Der Ausgang des Registers ist `geplant` mit dieser Kennung und bleibt es hinter dem Release (Entscheidung des Nutzers vom 2026-10-16: Der Slice liegt hinter [welle-erster-release](../done/welle-erster-release/welle-erster-release.md) und vor [welle-v1-nachtrag](../welle-v1-nachtrag.md); die Adresse nimmt an, solange der Slice nicht in `done/` liegt).

**Übernommen aus `slice-v1-abschluss-verbindungsende-replay-extended`** (dort §6, Risiko, und Register, Zähler 2; Sendung des Planners, 2026-10-16, `AGENTS.md` §3.13): dass die handgeschriebene Extended-Aufzeichnung mit einer Gruppe mit `Sync` (`ende-extended-sync.yaml`) dem entspricht, was PostgreSQL liefert. **Übernommen aus `slice-v1-abschluss-verbindungsende-replay-extended-flush`** (dort §6, Risiko, und Register, Zähler 3; Sendung des Planners, 2026-10-16): dasselbe für die Gruppe mit `Flush` (`ende-extended-flush.yaml`). Beide Geber liegen in `done/`; die Sendung steht in ihren Risiko-Ausgängen und im Register, in diesem Slice als Adresse.

**Voraussetzung, keine Lieferung — `slice-v1-abschluss-verbindungsende-upstream`** (Schritt 41): Der Upstream-Adapter erkennt das Ende in `Send` und `Receive` einer Extended-Interaktion (dort DoD-Punkt 1, Test mit Fake-Server), und das E2E gegen PostgreSQL für eine einfache Anfrage (Abnahmeszenario 18, Fixtures `ende-fatal.yaml` und `ende-lebend.yaml`) liegt dort. Dieser Slice setzt voraus, dass es in `done/` liegt, und nutzt die Hilfen seiner E2E-Datei, soweit vorhanden; er stellt keine Anforderung an sie.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Upstream-Adapter, der Record-Service, der Replay-Service und der PGWire-Adapter (Produktcode) — `slice-v1-abschluss-verbindungsende-upstream`, `-record`, `-replay-extended` und `-replay-extended-flush` (alle Voraussetzung): Schicht-Abgrenzung; der Slice belegt, er baut nicht. Zeigt der Server ein Verhalten, das Produktcode ändern müsste, geht der Befund an den Planner (§4).
- Die einfache Anfrage und die Fixtures `ende-fatal.yaml`, `ende-lebend.yaml`, `ende-panic.yaml` — `slice-v1-abschluss-verbindungsende-upstream` (DoD-Punkt 3, Sendung aus `-replay`): Der Beleg dort führt die Entsprechung für die einfache Interaktion; diese Hälfte der Beobachtung hat dort ihre Adresse.
- Ende im Leerlauf, Netzwerkabbruch ohne Fehlerantwort, `ERROR` mit anschließendem Abbruch — `slice-v1-abschluss-verbindungsende-upstream` (DoD-Punkt 2 und 3): andere Fälle des Szenarios, hier nicht Gegenstand der Fixtures.
- Die Aufzeichnung in SQLite — der Beleg für `meta.version` 2 gegen echtes Ende liegt in `slice-v1-abschluss-verbindungsende-upstream` (DoD-Punkt 3); das Format ist von der Art der Interaktion unabhängig, und beide Schreiber sind in `-format-schreiber` und `-sqlite-schreiber` geprüft. Bestand bleibt bewusst stehen.
- Andere PostgreSQL-Versionen — Szenario 18 läuft nur gegen die Referenz 17 (akzeptiertes Negativ aus `slice-v1-abschluss-verbindungsende-upstream` §6, Entscheidung des Nutzers vom 2026-10-12); die Matrix aus `slice-v1-abschluss-postgres-versionen-matrix` fährt Szenarien 1, 2 und 7. Der Beleg hier gilt für 17 und sagt nichts über andere (§6).
- Eine Erweiterung von Abnahmeszenario 18 im Lastenheft auf Extended Query — anderer Vorgang: Eine Änderung des Lastenhefts ist ein Change Request (`harness/conventions.md` §Versionierung); die Tests tragen `LH-FA-25` und `LH-FA-18`, nicht das Szenario.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (Schritt 42, `AGENTS.md` §3.11, §3.13): Der Slice liegt seit dem 2026-10-16 hinter dem Release, also nach Schritt 42; das Handbuch sagt zum Ende einer Extended-Interaktion beim Aufzeichnen gegen PostgreSQL nichts und nennt es erst, wenn dieser Slice in `done/` liegt (Bindung, dort §1 und §4). Dieser Slice liefert keinen Handbuch-Satz (kein neues Verhalten, `AGENTS.md` §3.11); ob nach seinem `done/` einer entsteht, ist eine offene Frage an den Nutzer (Bericht des Planners, 2026-10-16).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1** — [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Gruppe mit `Sync`. Ein Roh-Client mit Extended Query, dessen Gruppe mit Ende mit `Sync` schließt (Gruppen davor mit `Sync` und mit `Flush`, wie `ende-extended-sync.yaml`), wird am Binary gegen PostgreSQL 17 aufgezeichnet; die Aufzeichnung enthält in der Gruppe mit Ende die Server-Nachrichten bis einschließlich der Fehlerantwort `FATAL` `57P01` und danach `connection_closed`, kein `ready_for_query`, und trägt `version: 2`; die Typenfolge je Gruppe und Schweregrad und SQLSTATE der Fehlerantwort stimmen mit der Fixture überein (Vergleich im Test gegen die eingebettete Datei, nicht in Prosa); `replay` gibt die echte Aufzeichnung ohne PostgreSQL wieder (dieselben Antworten, danach Ende; Exit-Code 0 mit `--fail-on-unconsumed`). Weicht der Server ab, ist die Fixture berichtigt und ihre Tests nachgezogen (Integrationstest, Abdeckungstabelle `abdeckung-e2e.md`).
- [ ] **Liefer-Punkt 2** — [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Gruppe mit `Flush`. Dasselbe für die Gruppe mit Ende, die mit `Flush` schließt (`ende-extended-flush.yaml`), mit der Besonderheit, dass auch die Gruppe davor mit `Flush` schließt: Der Server sendet keine `ReadyForQuery`; die echte Aufzeichnung entspricht der Fixture oder die Fixture ist berichtigt, und `replay` gibt sie wieder (Integrationstest, Abdeckungstabelle).
- [ ] Je Zusage der beiden Punkte (Typenfolge je Gruppe; Schweregrad und SQLSTATE; kein `ready_for_query` nach dem Ende; `connection_closed` als letzte Nachricht; `version: 2`; Wiedergabe; je Element „Gruppe davor mit `Sync`“ und „mit `Flush`“ ein Fall) eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10), die Zählung der Tabelle gegen §6 vor der Übergabe.
- [ ] Der Vergleich ist an beiden Fixtures belegt: Der Test benennt die Fixture, die er prüft; eine Fixture ohne Vergleich gibt es nicht (`AGENTS.md` §3.11: der Satz „entspricht dem Server“ gilt nur, wo ein Test ihn prüft).
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
| `test/integration` (`ende_extended_e2e_test.go`, eine Datei je Fall) | neu | Fälle `Sync` und `Flush`: `record` gegen PostgreSQL 17 mit Roh-Client (Hilfen aus `replay_ende_extended_e2e_test.go` und, soweit vorhanden, der E2E-Datei von `-upstream`), Vergleich der Aufzeichnung mit der Fixture, `replay` der echten Aufzeichnung — `LH-FA-25`, `LH-FA-18` |
| `test/integration/testdata` (`ende-extended-sync.yaml`, `ende-extended-flush.yaml`) | update, nur bei Abweichung | Form des Servers, wo die handgeschriebene abweicht; mit ihr die Tests in `replay_ende_extended_e2e_test.go`, soweit sie die Form prüfen |
| `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabelle aus den Deklarationen der Tests |

Kein Produktcode, keine Spezifikation, keine Gates.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-upstream` (Schritt 41), `-record` (Schritt 40), `-replay-extended` (37b) und `-replay-extended-flush` (37c) liegen in `done/`. Der Slice steht seit dem 2026-10-16 (Entscheidung des Nutzers, Drift-Log der Roadmap) hinter [welle-erster-release](../done/welle-erster-release/welle-erster-release.md) und seit dem 2026-10-18 als erster Slice in [welle-v1-nachtrag](../welle-v1-nachtrag.md) (Entscheidung des Nutzers vom 2026-10-18); er war Schritt 41b von welle-v1-abschluss und ist keine Voraussetzung eines Schritts mehr (Schritt 42 setzt 35a bis 41 voraus). Seine Voraussetzungen liegen vor dem Release in `done/`. WIP-Limit 1.

**Größe, vor dem Anspruch geschätzt (Planner, 2026-10-16, `AGENTS.md` §3.14):** Schätzung **300 bis 500, Mitte 400**, 57 % der Marke 700 und damit unter der Marke (geschnitten wird nur, wenn die Schätzung die Marke überschreitet, Praxis des Planners seit 2026-10-16; die Meldeschwelle bleibt); Faktor Tests zu Produktcode: bei Produktcode 0 nicht bestimmbar (Beleg `slice-v1-abschluss-verbindungsende-replay-extended` und `-replay-extended-flush`, beide Produktcode 0), deshalb gerechnet in der Metrik der E2E-Fälle: rund 66 Zeilen je Fall mit Testdaten (`-replay` 68, `-replay-extended` 66, `-replay-extended-flush` 50), hier mit Vergleich und Wiedergabe der echten Aufzeichnung je Fall rund 100 (zwei Fälle, 200), dazu Hilfen für den Vergleich der Aufzeichnung mit der Fixture und für das Beenden des Backends rund 120, Fixtures und Tabellenzeilen rund 60, Nacharbeit zum Review rund 20 %. Der Mindestwert 5 aus dem Faktor-Slice (`slice-harness-schaetzung-planner-faktor`, Wortlaut vom Nutzer noch nicht bestätigt) ist ohne Produktcode nicht anwendbar.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 700** geht der Slice zurück zur Zerlegung, Schnitt `Sync` | `Flush` (Zuschnitt: Planner); schätzt die Prüfung des Architect vor dem Code die Mitte über der Marke 700, wird vor dem Anspruch geschnitten, sonst nicht. **Meldeschwelle 550:** Erreicht die Summe 550, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`); die Entscheidung des Planners steht als datierte Zeile in §7 (`weiter` oder `schneiden`).
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

**Randformen** (`AGENTS.md` §3.12). Der Slice liefert keinen neuen Vertrag, sondern einen Beleg gegen einen Vertrag, der in der Spezifikation steht (`LH-FA-25.a` *Aufzeichnen*, *Stelle*, `LH-FA-18.a` *Interaktion*, *Abschluss*). Die Randformen des Belegs entscheidet der Architect vor dem ersten Code-Commit; bis dahin sind sie offen und der Implementer entscheidet sie nicht (Randform-Rückgabe):

- *Weg, das Ende auszulösen:* die Anfrage der Gruppe selbst (`pg_terminate_backend(pg_backend_pid())` als SQL der Gruppe) oder eine zweite Verbindung, während der Server auf `Sync`/`Flush` wartet. Beides ergibt unterschiedliche Formen (Fehlerantwort nach `parse_complete` ohne weitere Nachrichten, oder nach `bind_complete`/Ergebnis der Anfrage); die Fixtures zeigen nur `parse_complete` und die Fehlerantwort. Entscheidungsort: der Architect, in der Prüfung vor dem Code, mit der Bedingung, dass der Fall nicht von der Zeit abhängt (sonst §4, Rückführung nach `open/`).
- *Stelle der Fehlerantwort in der Folge der Gruppe:* nach welchen Server-Nachrichten der Gruppe sie steht, und ob `ready_for_query` oder weitere Nachrichten folgen — `LH-FA-25.a` *Aufzeichnen*, *Stelle*; der Vergleich prüft die Typenfolge der Gruppe mit Ende gegen die Fixture, ohne sie vorab zu setzen; die Abweichung ist der Befund (eingetreten).
- *Gruppen davor:* mit `Sync` (bis `ready_for_query`) und mit `Flush` (bis `parse_complete`), je Element ein Fall — `LH-FA-18.a`; die Fixture `ende-extended-flush.yaml` hat nur Gruppen mit `Flush`, die Fixture `ende-extended-sync.yaml` beide.
- *Lesestand:* eine vom Server gesendete, vom Recorder beim Bemerken des Endes nicht gelesene Fehlerantwort zählt nicht — `LH-FA-25.a` *Aufzeichnen*, `LH-FA-02.b` *Gelesen*; der Fall „Fehlerantwort liegt vor dem Schließen im Socket“ ist der des Belegs, entschieden im Upstream-Slice (Voraussetzung).
- *Verbindungsende ohne gelesene Fehlerantwort, Leerlauf, Netzwerkabbruch, `ERROR`:* nicht Gegenstand — `slice-v1-abschluss-verbindungsende-upstream` (§1).
- *Aufruf-Voraussetzung, Aufrufarten des Treibers:* nicht berührt; der Slice liefert weder Werte über mehrere Aufrufe noch ein Ergebnis, das ein Hook oder Gate glaubt.
- *Version des Servers:* gilt nur für die Referenz 17; eine Zeile in `SPEC-029` Punkt 2 entsteht nicht (sie benennt Abnahmeszenarien), das Handbuch nennt keine Version darüber hinaus (Bindung an `slice-v1-abschluss-verbindungsende-doku`).

**Risiken:**

- Die Form des Servers weicht von einer Fixture ab (etwa zusätzliche Nachrichten vor der Fehlerantwort) und die Tests von `-replay-extended` und `-replay-extended-flush` müssen nachgezogen werden — **Ausgang:** offen bis Closure (erwartet: eingetreten und mit diesem Slice berichtigt, oder entfallen mit dem Beleg der Übereinstimmung; verlangt die Abweichung Produktcode: eingetreten → Folge-Slice mit ID).
- Das Verhalten des Servers ist nur gegen die Referenzversion 17 geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×; Akzeptiertes Negativ wie in `slice-v1-abschluss-verbindungsende-upstream` §6) — **Ausgang:** weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`.
- Ein Fall hängt von der Zeit ab (Verhältnis von Signal des Servers und `Sync` des Clients) und flackert — **Ausgang:** offen bis Closure (§4: Rückführung nach `open/`, Entscheidung des Nutzers).
- Der Diff liegt über der Schätzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 9×, Zählregel vom 2026-10-16: die Überschreitung der Obergrenze 500 zählt) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.
- `slice-v1-abschluss-verbindungsende-doku` beschreibt das Ende einer Extended-Interaktion gegen PostgreSQL, bevor dieser Slice es belegt (Schritt 42 liegt seit dem 2026-10-16 vor diesem Slice) — **Ausgang:** offen bis Closure (Gegenmittel: Bindung in `-doku` §1 und §4: das Handbuch sagt dazu nichts, bis dieser Slice in `done/` liegt; Prüfung beim Review von `-doku` durch Suche nach dem Wortlaut).
- Das Release erscheint mit handgeschriebenen Fixtures der Extended-Interaktion mit Ende, die kein Test gegen PostgreSQL belegt (Entscheidung des Nutzers vom 2026-10-16; `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`, 5×) — **Ausgang:** weiter offen → `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (Ausgang `geplant` mit dieser Kennung, Vermerk „hinter dem Release“ im Register).

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

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-16). Treffer in der Sub-Area `REPO`: `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (3×, Auslöser; Ausgang `geplant` mit dieser Kennung), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (9×, `geplant` in `slice-harness-schaetzung-planner-faktor`; der Plan nennt Schätzung, Marke und Meldeschwelle vorab, §4), `BEO-REPO/meldeschwelle-ohne-artefakt` (2×; §4) und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×; §6). Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-16, `AGENTS.md` §3.13): Liefer-Punkte: 2 (Gruppe mit `Sync`, Gruppe mit `Flush`; Vergleich und Wiedergabe gehören zu jedem Punkt). Schichten: **eine** ohne und **zwei** mit der konservativen Zählung: die Integrationstests am Binary (`test/integration`) sowie, als Markierung, die erzeugte Tabelle `docs/user/abdeckung-e2e.md` als Nutzer- und Wartungs-Doku (`AGENTS.md` §3.13 *Schichten zählen alle Pläne gleich*; ob erzeugte Abdeckungstabellen zählen, ist nicht entschieden); auch mit ihr zwei. Kein Produktcode, keine Spezifikation. Schätzung in der Metrik von §4: 300 bis 500 (Mitte 400, 57 % der Marke 700).

Nachgezählt beim Eintragen der Sendungen aus `slice-v1-abschluss-verbindungsende-replay-extended` und `-replay-extended-flush` (Planner, 2026-10-16, `AGENTS.md` §3.13): Beide gehören zu Liefer-Punkt 1 und 2 und fügen weder Liefer-Punkt noch Schicht hinzu: Liefer-Punkte weiter 2, Schichten weiter eine, mit Markierung zwei.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Verschieben hinter das Release (Planner, 2026-10-16, `AGENTS.md` §3.13): Welle-Zugehörigkeit, Start, Herkunft und zwei Risiken geändert, Gegenstand und DoD unverändert: Liefer-Punkte weiter 2, Schichten weiter eine (zwei mit der Markierung der erzeugten Tabelle); Schätzung weiter 300 bis 500 (Mitte 400, 57 % der Marke 700). Die Verschiebung nimmt der Slice keinem Nehmer ab und trägt keine Sendung ein.

Nachgezählt beim Anpassen von §4 (Planner, 2026-10-16, `AGENTS.md` §3.13): Die Prüfung des Architect vor dem Code führt nicht mehr bei einer Mitte über 490 (70 % der Marke) zur Meldung, sondern der Schnitt vor dem Anspruch folgt nur einer Schätzung über der Marke 700; Marke 700 und Meldeschwelle 550 bleiben. Kein Liefer-Punkt, keine Schicht.

Nachgezählt beim Eintragen in welle-v1-nachtrag (Planner, 2026-10-18, `AGENTS.md` §3.13): Nur Welle-Zugehörigkeit, Kopf und Start geändert, Gegenstand und DoD unverändert: Liefer-Punkte weiter 2, Schichten weiter eine (zwei mit der Markierung der erzeugten Tabelle); Schätzung weiter 300 bis 500. Die Eintragung trägt keine Sendung ein.
