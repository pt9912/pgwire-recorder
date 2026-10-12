# Slice slice-v1-abschluss-verbindungsende-upstream: Verbindungsende aus dem Upstream erkennen (Abnahmeszenario 18)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-02.b` · `LH-FA-18.a` · `LH-FA-05.a` · `SPEC-034` · `ARC-004` · `ARC-007`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Upstream-Adapter stuft das Ende des Servers ein (`LH-FA-25.a` *Aufzeichnen*, *Stelle*, *Leerlauf*, *Kein Ende*): Liest er in der laufenden Interaktion eine `ErrorResponse` mit dem Schweregrad `FATAL` oder `PANIC`, ist die Interaktion abgeschlossen — er liefert die Antworten bis einschließlich dieser Fehlerantwort und danach `connection_closed`, ohne auf das Verbindungsende zu warten und ohne zu lesen, was danach kommt; bei einer einfachen Anfrage wie in einer Extended-Interaktion. Ein Ende im Leerlauf ist nur mit nachfolgender Anfrage aufzeichenbar; ein Verbindungsende ohne gelesene `FATAL`/`PANIC` bleibt `PGR-E4003`, nach einer anderen gelesenen `ErrorResponse` `PGR-E6001`, vor dem ersten `ReadyForQuery` `PGR-E4002` beziehungsweise `PGR-E4005`. Mit diesem Slice zeichnet `record` ein Ende gegen PostgreSQL auf und `replay` gibt es wieder: Abnahmeszenario 18 ist automatisiert.

**Übernommen aus `slice-v1-abschluss-verbindungsende-record`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): der Upstream-Adapter und Abnahmeszenario 18. Der Use Case und die Zustellung liegen in jenem Slice; dieser setzt ihn voraus und liefert, was dessen Fakes bisher nur abbilden.

**Übernommen aus `slice-v1-abschluss-verbindungsende-sqlite`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12): der E2E-Beleg, dass `record --format sqlite` ein echtes Ende aufzeichnet und die Datei `meta.version` 2 trägt, und aus `slice-v1-abschluss-verbindungsende-replay` (dort §3 und §6, Risiko): dass die handgeschriebenen Aufzeichnungen dessen Fixtures dem entsprechen, was PostgreSQL liefert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Use Case, `CloseSession`, die Übernahme der Session und die Zustellung im PGWire-Adapter — `slice-v1-abschluss-verbindungsende-record` (Vorgänger): andere Schichten.
- Fehlerantwort `FATAL` im Verbindungsaufbau (etwa eine abgelehnte Anmeldung) — `Open` bleibt unverändert; das Ende vor dem ersten `ReadyForQuery` ist kein Ende nach `LH-FA-25` (Out-of-Scope; `BEO-REPO/record-fehlerantwort-im-aufbau-ungeregelt` bleibt offen).
- Ein Netzwerkabbruch ohne Fehlerantwort als Ende — Out-of-Scope von `LH-FA-25`; er bleibt `PGR-E4003` (Teil von Abnahmeszenario 18, hier belegt).
- Die Matrix der PostgreSQL-Versionen — `slice-v1-abschluss-postgres-versionen` fährt die Szenarien 1, 2 und 7; Szenario 18 läuft hier gegen die Referenzversion 17. Ob die Matrix es mitfährt, ist eine Frage an den Nutzer (§6).
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus): Der Upstream-Adapter liefert bei `Query` die Antworten bis einschließlich einer gelesenen Fehlerantwort `FATAL`/`PANIC` und danach `connection_closed`, ohne weiter zu lesen, und liefert dieselbe Einstufung bei `Receive` und `Send` einer Extended-Interaktion; ein Ende ohne gelesene `FATAL`/`PANIC` ist `PGR-E4003`, nach einer anderen gelesenen `ErrorResponse` (`ERROR` und anschließender Abbruch) `PGR-E6001`; der Adapter wartet nicht auf das Verbindungsende (ein Server, der nach der Fehlerantwort offen bleibt, hält den Aufruf nicht auf) (Test mit Fake-Server).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Ende im Leerlauf: Sendet der Client nach dem Ende des Servers noch eine Anfrage, gehört die Fehlerantwort zu ihr (nach den Nachrichten zwischen den Interaktionen), und ein Sendefehler an den geschlossenen Upstream ist kein Fehler und ändert die Einstufung nicht; ein Sendefehler vor einer gelesenen Fehlerantwort bleibt `PGR-E4003`; sendet der Client keine mehr, endet die Session regulär ohne Aufzeichnung der Fehlerantwort (Test mit Fake-Server und Integrationstest mit `idle_session_timeout`).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Abnahmeszenario 18 am Binary gegen PostgreSQL 17 (`pg_terminate_backend(pg_backend_pid())`, `57P01`): Die Anwendung erhält im Record-Modus die Antworten und danach das Ende ihrer Verbindung; die Aufzeichnung enthält das Ende und trägt `version: 2` (YAML) beziehungsweise `meta.version` 2 (SQLite); im Replay ohne PostgreSQL erhält die Anwendung dieselben Antworten und danach das Ende; ein Netzwerkabbruch ohne Fehlerantwort wird nicht aufgezeichnet, sondern als `PGR-E4003` gemeldet; der Exit-Code des Ende-Laufs bleibt 0 (Integrationstest, Abdeckungstabelle `abdeckung-e2e.md`).
- [ ] Je Zusage des neuen Vertrags eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/postgres/upstream.go` | update | `Query`, `Send`, `Receive`, `lies`, `verbindungsende`: Ende erkennen, `connection_closed` liefern, nicht weiterlesen; Einstufung `PGR-E4003`/`PGR-E6001` |
| `internal/hexagon/ports/driven/upstream.go` | update (Vertrag) | Dokumentation des Port-Vertrags: wann `connection_closed` kommt und dass danach nichts mehr gelesen wird (Name und Form hat der Vorgänger festgelegt) |
| `internal/adapters/driven/postgres` (Tests) | update | Randformen aus §6; die Tests, die die alte Regel prüfen (`TestFehlerantwortVorDemAbbruch` für `FATAL` und `PANIC`, `TestSendNachFehlerantwort`), werden hier nachgezogen |
| `test/integration` (`ende_e2e_test.go`; `TestE2ERecordExtendedAbbruch` bleibt für `ERROR`) | neu/update | Abnahmeszenario 18 gegen PostgreSQL 17 in YAML und SQLite, Leerlauf, Netzwerkabbruch |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-record` und `slice-v1-abschluss-verbindungsende-sqlite` liegen in `done/` (Voraussetzung); `slice-v1-abschluss-anmeldung` liegt davor, weil es die Aufbau-Schleife im Upstream-Adapter umbaut. Dieser Slice ist Schritt 41 der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Prüfung des Architect vor dem Code schätzt den Diff mit Tests und E2E auf über 1000 Zeilen (Schätzung des Planners 700 bis 1200) — Schnitt vor dem Anspruch: Adapter | E2E von Szenario 18 mit Leerlauf, Zuschnitt: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`).
- `in-progress` → `next`: Der Adapter erkennt das Ende nur mit einer Änderung des Use Case oder des Ports — dritte Schicht (der Kern mit dem Port zählt als eine), zurück zur Zerlegung.
- `in-progress` → `open`: Der Leerlauf-Fall ist gegen PostgreSQL nicht deterministisch herstellbar (`idle_session_timeout` und die nachfolgende Anfrage) — der Nutzer entscheidet, ob der Fall mit Fake-Server ausreicht (Carveout).

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect; „Leerlauf nur mit nachfolgender Anfrage“ und „aufgezeichnet wird genau: Server beendet in laufender Interaktion mit `FATAL`/`PANIC`“ bestätigte der Nutzer am 2026-10-12). Ort: `LH-FA-25.a` *Aufzeichnen*, *Stelle*, *Leerlauf*, *Kein Ende*, `LH-FA-02.b` *Gelesen*, `LH-FA-05.a`.

- *Ende bei `PANIC` und `FATAL`:* beide — `LH-FA-25.a`.
- *`ERROR`, dann `FATAL`:* die letzte gelesene Fehlerantwort entscheidet über das Ende; die frühere ist Teil der Antworten — `LH-FA-25.a`, `LH-FA-02.b`.
- *`ERROR` und dann Abbruch:* bleibt `PGR-E6001` mit SQLSTATE und Meldung der letzten Fehlerantwort — `LH-FA-25.a` *Kein Ende*, `LH-FA-02.b`.
- *Fehlerantwort gelesen oder nicht (Lesestand) bei Sendefehler und Verbindungsende:* eine beim Bemerken des Endes nicht gelesene Fehlerantwort zählt nicht (`PGR-E4003`); eine gelesene `FATAL`/`PANIC` macht den Sendefehler danach zum Teil des Endes — `LH-FA-25.a` *Aufzeichnen*, `LH-FA-02.b` *Gelesen*.
- *Leerlauf mit und ohne nachfolgende Anfrage; Sendefehler an den geschlossenen Upstream:* siehe DoD-Punkt 2 — `LH-FA-25.a` *Leerlauf*, `LH-FA-05.a`.
- *Ende in einer `Sync`-Gruppe, in einer `Flush`-Gruppe, in einer unvollständigen Gruppe:* der Adapter meldet das Ende, wo es liegt; die Zuordnung zur Gruppe und `PGR-E6001` für die unvollständige trifft der Use Case (`slice-v1-abschluss-verbindungsende-record`) — `LH-FA-25.a` *Stelle*, `LH-FA-18.a`; bei einer `Flush`-Gruppe gilt die zeitabhängige Grenze von `LH-FA-18.a`.
- *Ende vor dem ersten `ReadyForQuery`:* bleibt `PGR-E4002`/`PGR-E4005` (`Open` unverändert) — `LH-FA-25.a` *Kein Ende*.
- *Ende mit `--record-timing` und `--record-empty-sessions`:* nicht Sache des Adapters — `slice-v1-abschluss-verbindungsende-record`.
- *Nach dem Ende gelesene Bytes:* werden nicht gelesen — `LH-FA-25.a` *Aufzeichnen*.

**Risiken:**

- Das Verhalten des Servers (`57P01`, `idle_session_timeout` ab PostgreSQL 14) ist nur gegen die Referenzversion 17 geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×); ob die Matrix aus `slice-v1-abschluss-postgres-versionen` Abnahmeszenario 18 mitfährt, ist eine Frage an den Nutzer, nicht in diesem Plan entschieden — **Ausgang:** offen bis Closure (erwartet: weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`).
- Ein Netzwerkabbruch ohne Fehlerantwort braucht im E2E einen Weg, die Verbindung zu kappen (Proxy, Container anhalten); das Mittel wählt der Implementer, die Zusage ist `PGR-E4003` statt Ende — **Ausgang:** offen bis Closure.
- Der Diff liegt an der Grenze einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×): die Rückführung in §4 benennt den Schnitt — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×), `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×), `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×; §6) und `BEO-REPO/record-fehlerantwort-im-aufbau-ungeregelt` (1×; Fehlerantwort im Aufbau bleibt ausgenommen, §1) — der Plan nennt Schätzung und Schnitt vorab (§4). Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13), einschließlich der Sendungen aus `slice-v1-abschluss-verbindungsende-record`, `-sqlite` und `-replay` (sie gehören zu DoD-Punkt 3 und fügen keinen Liefer-Punkt hinzu): Liefer-Punkte: 3 (Ende erkennen bei Query, Send und Receive; Leerlauf; Abnahmeszenario 18). Schichten: zwei (Upstream-Adapter, `ARC-007`; Kern nur als Vertrag des Ports, `ARC-004`). Schätzung des Planners: 700 bis 1200 Zeilen mit Tests.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
