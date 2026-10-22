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
| `test/integration/ende_extended_e2e_test.go` | neu | Fälle `Sync` und `Flush`: `record` gegen PostgreSQL 17 mit Roh-Client (`rohVerbindung`, `parseSync`), Vergleich der Aufzeichnung mit der eingebetteten Fixture über die Hilfen von `ende_e2e_test.go`, `replay` der echten Aufzeichnung. Die Wiedergabe ohne PostgreSQL läuft in der Phase ohne PostgreSQL des Runners (`tools/test/run-integration-tests.sh` wählt dort `^TestE2EOhnePostgres`; die Aufzeichnung liegt unter `PGR_DATEN`, wie bei `TestE2EOhnePostgresEnde`) — `LH-FA-25`, `LH-FA-18` |
| `test/integration/testdata` (`ende-extended-sync.yaml`, `ende-extended-flush.yaml`) und `replay_ende_extended_e2e_test.go` | update, nur bei Abweichung; der Kommentar der `go:embed`-Variablen immer | Form des Servers, wo die handgeschriebene abweicht, mit ihr die Tests dort, soweit sie die Form prüfen. Der Kommentar „gehört genau einem Test“ über `endeExtendedSyncYAML` und `endeExtendedFlushYAML` stimmt nicht mehr, sobald der neue Test sie auch einbettet (`AGENTS.md` §3.7) |
| `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabelle aus den Deklarationen der Tests |

Kein Produktcode, keine Spezifikation, keine Gates.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Voraussetzungen liegen in `done/`: `slice-v1-abschluss-verbindungsende-upstream`, `-record`, `-replay-extended` und `-replay-extended-flush`, dazu [welle-erster-release](../done/welle-erster-release/welle-erster-release.md) und `slice-harness-regeln-rueckbau`. Der Nutzer hat den Start der Produkt-Slices am 2026-10-22 mit „leg los“ freigegeben (wörtlich, so vom Koordinator übermittelt). Der Architect prüft vor dem Code die Frage nach dem Weg, das Ende auszulösen (§6). WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik und Messbefehl: `AGENTS.md` §3.14; Pfade des Plans: `test docs/user`); die Basis ist der Commit des Übergangs nach `in-progress/`, sein Hash steht beim Übergang hier. Dann hält der Implementer an und meldet dem Planner; Schnitt `Sync` | `Flush` (Zuschnitt: Planner).
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

**Randformen.** Der Slice liefert keine neue Wire-, Datei- oder Formatgrenze, sondern belegt Zusagen der Spezifikation (`LH-FA-25.a` *Aufzeichnen*, *Stelle*, `LH-FA-18.a` *Interaktion*, *Abschluss*) gegen den Server; eine Randformen-Liste nach `AGENTS.md` §3.12 ist nicht verlangt. Eine Frage entscheidet der Architect vor dem ersten Code-Commit, weil der Test von ihr abhängt; der Implementer entscheidet sie nicht:

- *Weg, das Ende auszulösen.* Die Fixtures zeigen in der Gruppe mit Ende `parse` und `sync` (`flush`), als Antwort `parse_complete` und die Fehlerantwort. Eine Gruppe aus `Parse` und `Sync` führt keine Anfrage aus; `pg_terminate_backend(pg_backend_pid())` als SQL der Gruppe löst das Ende erst mit `Bind` und `Execute` aus, und die Form wäre dann länger als die Fixture (Lesart des Planners, nicht gemessen). Der andere Weg, eine zweite Verbindung, die das Backend beendet, während der Server auf `Sync` oder `Flush` wartet, hängt von der Zeit ab, wenn der Test das Warten des Backends nicht beobachtet. Entscheidung: welcher Weg, mit der Bedingung, dass der Fall nicht von der Zeit abhängt (sonst §4, Rückführung nach `open/`). Weicht die Form von der Fixture ab, wird die Fixture berichtigt (§1).

**Risiken:**

- Die Form des Servers weicht von einer Fixture ab (etwa zusätzliche Nachrichten vor der Fehlerantwort); die Releases erschienen mit den handgeschriebenen Fixtures (`BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`, 6×) — **Ausgang:** offen bis Closure (erwartet: eingetreten und mit diesem Slice berichtigt, oder entfallen mit dem Beleg der Übereinstimmung; verlangt die Abweichung Produktcode: eingetreten → Folge-Slice mit ID; das Register geht in beiden Fällen auf `gestrichen` oder `verkörpert`, DoD).
- Das Verhalten des Servers ist nur gegen die Referenzversion 17 geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, Stand `verkörpert`; akzeptiertes Negativ wie in `slice-v1-abschluss-verbindungsende-upstream` §6) — **Ausgang:** weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`.
- Ein Fall hängt von der Zeit ab (Verhältnis von Signal des Servers und `Sync` des Clients) und flackert — **Ausgang:** offen bis Closure (§4: Rückführung nach `open/`, Entscheidung des Nutzers).
- `slice-v1-abschluss-verbindungsende-doku` beschreibt das Ende einer Extended-Interaktion gegen PostgreSQL, bevor dieser Slice es belegt — **Ausgang:** entfallen, weil `-doku` in `done/` liegt und das Handbuch dazu nichts sagt (`grep` vom 2026-10-22, §1); eine neue Aussage im Handbuch bräuchte einen eigenen Slice.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (6×, Auslöser; Ausgang `geplant` mit dieser Kennung) und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (§6; Stand `verkörpert`). Keine weiteren Treffer.

**Zählung** (Planner, 2026-10-22, `AGENTS.md` §3.13; Geber `slice-harness-regeln-rueckbau`, dessen Reihenfolge den Start in §4 und die Marke 1500 setzt; Randformen-Liste, Mutationspflicht, Schätzung und Meldeschwelle entfallen): **Liefer-Punkte: 2** (Gruppe mit `Sync`, Gruppe mit `Flush`; Vergleich, Wiedergabe und die Sendungen aus `-replay-extended` und `-replay-extended-flush` gehören zu jedem Punkt). **Schichten: eine** (die Integrationstests am Binary, `test/integration`), mit der Markierung der erzeugten Tabelle `docs/user/abdeckung-e2e.md` (Nutzer- und Wartungs-Doku, `AGENTS.md` §3.13) zwei. Kein Produktcode, keine Spezifikation.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
