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

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-02.b` · `LH-FA-18.a` · `LH-FA-05.a` · `SPEC-034` · `SPEC-029` · `ARC-004` · `ARC-007`

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

**Übernommen aus `slice-v1-abschluss-verbindungsende-format-schreiber`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-15, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-format`): der YAML-Schreiber mit Version nach dem Inhalt. Abnahmeszenario 18 zeichnet in YAML auf und lädt die Datei wieder (DoD-Punkt 3); den Schreiber liefert jener Slice, dieser setzt ihn voraus.

**Übernommen aus `slice-v1-abschluss-verbindungsende-record`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): der Upstream-Adapter und Abnahmeszenario 18. Der Use Case und die Zustellung liegen in jenem Slice; dieser setzt ihn voraus und liefert, was dessen Fakes bisher nur abbilden.

**Übernommen aus `slice-v1-abschluss-verbindungsende-sqlite-schreiber`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12 aus `slice-v1-abschluss-verbindungsende-sqlite`, mit dem Schnitt vom 2026-10-15 an den Schreiber gegangen): der E2E-Beleg, dass `record --format sqlite` ein echtes Ende aufzeichnet und die Datei `meta.version` 2 trägt, und aus `slice-v1-abschluss-verbindungsende-replay` (dort §3 und §6, Risiko): dass die handgeschriebenen Aufzeichnungen dessen Fixtures dem entsprechen, was PostgreSQL liefert.

**Gebunden durch `slice-v1-abschluss-sessions`** (dort §1, *Nahtstelle zu den Verbindungsende-Slices*, Punkt 4; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): *Eine Verbindung wird erst mit eröffneter Session eine Session.* Ein Ende des Servers vor dem ersten `ReadyForQuery` bleibt ein Aufbaufehler (`PGR-E4002`/`PGR-E4005`) und erzeugt keine Session, auch mit `--record-empty-sessions`; `Open` liefert dann den Fehler und der Dienst die Kennung 0. Dieser Slice ändert `Open` nicht und erkennt ein Ende nur in der laufenden Interaktion nach dem ersten `ReadyForQuery`.

**Übernommen aus `slice-v1-abschluss-postgres-versionen-zusage`** (dort §6, Randform „Szenarien nur gegen 17“; Sendung des Architect, 2026-10-15, `AGENTS.md` §3.13): `SPEC-029` Punkt 2 nennt Abnahmeszenario 18 heute nicht, weil kein Test es prüft. Im Commit, der den Integrationstest zu Szenario 18 liefert, ergänzt dieser Slice Punkt 2 um den Satz „Das Abnahmeszenario 18 ist nur dort geprüft“ (dort = Referenz 17); vorher steht er nicht (`AGENTS.md` §3.11), und die Handbuch-Sätze bleiben beim Doku-Slice. Dazu berichtigt der erste der beiden Nehmer, der Punkt 2 berührt, den Numerus des vorhandenen Satzes „Die Anmeldung mit Klartext-Passwort und mit MD5 sind nur dort geprüft“ (Verifikation `slice-v1-abschluss-postgres-versionen-zusage`, I-1, INFO): „Die Anmeldung mit Klartext-Passwort und die mit MD5 sind nur dort geprüft“; der zweite prüft, dass es geschehen ist. Ein Wort im selben Satz, kein Liefer-Punkt. Geschehen: Die Berichtigung steht im Commit `179ab424` (Architect, bei der Prüfung von `slice-v1-abschluss-verbindungsende-hilfetext-binary`, Schritt 39b, vor diesem Slice); dieser Slice prüft nur, dass sie steht, und berichtigt nicht noch einmal.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Use Case, `CloseSession`, die Übernahme der Session und die Zustellung im PGWire-Adapter — `slice-v1-abschluss-verbindungsende-record` (Vorgänger): andere Schichten.
- Fehlerantwort `FATAL` im Verbindungsaufbau (etwa eine abgelehnte Anmeldung) — `Open` bleibt unverändert; das Ende vor dem ersten `ReadyForQuery` ist kein Ende nach `LH-FA-25` (Out-of-Scope; `BEO-REPO/record-fehlerantwort-im-aufbau-ungeregelt` bleibt offen).
- Ein Netzwerkabbruch ohne Fehlerantwort als Ende — Out-of-Scope von `LH-FA-25`; er bleibt `PGR-E4003` (Teil von Abnahmeszenario 18, hier belegt).
- Die Matrix der PostgreSQL-Versionen — `slice-v1-abschluss-postgres-versionen-matrix` fährt die Szenarien 1, 2 und 7; Szenario 18 läuft hier gegen die Referenzversion 17. Ob die Matrix es mitfährt, ist eine Frage an den Nutzer (§6).
- Der Beleg gegen PostgreSQL für die Extended-Fixtures `ende-extended-sync.yaml` und `ende-extended-flush.yaml` (Extended-Client, Gruppe mit `Sync` und mit `Flush`) — `slice-v1-abschluss-verbindungsende-upstream-extended` (Folge-Slice, wellenlos hinter dem Release und vor welle-v1-nachtrag, setzt diesen Slice voraus; Entscheidung des Nutzers vom 2026-10-16 zum Slice und zur Stelle hinter dem Release, Register `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`): DoD-Punkt 3 fährt eine einfache Anfrage, ein Extended-Fall ließe den Punkt wachsen (`BEO-REPO/slice-waechst-durch-uebernahmen`). Der Adapter-Test mit Fake-Server für `Send` und `Receive` einer Extended-Interaktion (DoD-Punkt 1) bleibt hier.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus): Der Upstream-Adapter liefert bei `Query` die Antworten bis einschließlich einer gelesenen Fehlerantwort `FATAL`/`PANIC` und danach `connection_closed`, ohne weiter zu lesen, und liefert dieselbe Einstufung bei `Receive` und `Send` einer Extended-Interaktion (`Send` nach einer gelesenen `FATAL`/`PANIC` liefert keinen Fehler, auch wenn das Schreiben scheitert; der Schweregrad ist `V`, ein leeres `V` gilt als fehlend, dann `S`, genau `FATAL` und `PANIC`); ein Ende ohne gelesene `FATAL`/`PANIC` ist `PGR-E4003`, nach einer anderen gelesenen `ErrorResponse` (`ERROR` und anschließender Abbruch) `PGR-E6001`; der Adapter wartet nicht auf das Verbindungsende (ein Server, der nach der Fehlerantwort offen bleibt, hält den Aufruf nicht auf) (Test mit Fake-Server).
- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Ende im Leerlauf: Sendet der Client nach dem Ende des Servers noch eine Anfrage, gehört die Fehlerantwort zu ihr (nach den Nachrichten zwischen den Interaktionen); scheitert das Senden dieser Anfrage vor einer gelesenen Fehlerantwort (Reset des Upstreams), ist es `PGR-E4003` und die Fehlerantwort im Puffer zählt nicht (`LH-FA-25.a` *Leerlauf*, Berichtigung des Architect vom 2026-10-16); sendet der Client keine mehr, endet die Session regulär ohne Aufzeichnung der Fehlerantwort (Test mit Fake-Server und Integrationstest mit `idle_session_timeout`).
- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Abnahmeszenario 18 am Binary gegen PostgreSQL 17, dazu der Formvergleich (die aufgezeichnete Interaktion mit Ende trägt als letzte zwei Antworten eine `error_response` mit denselben Feldern `V`, `S`, `C`, `M` und `connection_closed` wie die handgeschriebenen Fixtures `ende-fatal.yaml`/`ende-terminate.yaml` für `57P01` und `ende-lebend.yaml` für `57P05`; `F`, `L`, `R` und die Antworten davor zählen nicht), dazu der Satz in `SPEC-029` Punkt 2, dass es nur dort geprüft ist (Sendung aus `slice-v1-abschluss-postgres-versionen-zusage`, im Commit des Tests) (`pg_terminate_backend(pg_backend_pid())`, `57P01`): Die Anwendung erhält im Record-Modus die Antworten und danach das Ende ihrer Verbindung; die Aufzeichnung enthält das Ende und trägt `version: 2` (YAML) beziehungsweise `meta.version` 2 (SQLite); im Replay ohne PostgreSQL erhält die Anwendung dieselben Antworten und danach das Ende; ein Netzwerkabbruch ohne Fehlerantwort wird nicht aufgezeichnet, sondern als `PGR-E4003` gemeldet; der Exit-Code des Ende-Laufs bleibt 0 (Integrationstest, Abdeckungstabelle `abdeckung-e2e.md`).
- [x] Je Zusage des neuen Vertrags eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/postgres/upstream.go` | update | `Query`, `Send`, `Receive`, `lies` und neu `beendetVerbindung`, `endeGelesen`: Ende erkennen, `connection_closed` liefern, nicht weiterlesen; die Einstufung `PGR-E4003`/`PGR-E6001` in `verbindungsende` bleibt |
| `internal/hexagon/ports/driven/upstream.go` | update (nur Kommentar, keine Logik, keine Schicht nach `AGENTS.md` §3.13) | Dokumentation des Port-Vertrags: wann `connection_closed` kommt und dass danach nichts mehr gelesen wird (Name und Form hat der Vorgänger festgelegt) |
| `internal/adapters/driven/postgres` (Tests: neu `upstream_ende_test.go`, nachgezogen `upstream_test.go`) | update | Randformen aus §6; die Tests, die die alte Regel prüfen (`TestFehlerantwortVorDemAbbruch` für `FATAL` und `PANIC`, `TestSendNachFehlerantwort`), werden hier nachgezogen: `FATAL` und `PANIC` sind ein Ende (`TestEnde*`, `TestSendNachDemEnde`), die bestehenden Tests prüfen `ERROR` |
| `test/integration` (`ende_e2e_test.go`; `TestE2ERecordExtendedAbbruch` bleibt für `ERROR`) | neu/update | Abnahmeszenario 18 gegen PostgreSQL 17 in YAML und SQLite (`TestE2ERecordEndeTerminate`, Replay ohne PostgreSQL in `TestE2EOhnePostgresEnde`, Formvergleich), Leerlauf (`TestE2ERecordEndeLeerlauf`), Netzwerkabbruch (`TestE2ERecordNetzwerkabbruch`) |
| `spec/spezifikation.md` (`SPEC-029` Punkt 2) | update | der Satz „Das Abnahmeszenario 18 ist nur dort geprüft.“ im Commit des Tests (Sendung aus `slice-v1-abschluss-postgres-versionen-zusage`, §1) |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-record`, `slice-v1-abschluss-verbindungsende-sqlite`, `slice-v1-abschluss-verbindungsende-sqlite-schreiber` und `slice-v1-abschluss-verbindungsende-format-schreiber` liegen in `done/` (Voraussetzung; der Schreiber für die Aufzeichnung in YAML mit `version: 2`); `slice-v1-abschluss-anmeldung` liegt davor, weil es die Aufbau-Schleife im Upstream-Adapter umbaut. Dieser Slice ist Schritt 41 der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal/adapters/driven/postgres internal/hexagon/ports/driven test/integration docs/user` liegt über 1200 (`AGENTS.md` §3.14; **Basis:** `c5089444`, der Commit vor dem Anspruch; die Plan- und der Move-Commit des Anspruchs berühren keinen der Pfade, die Basis ändert sich danach nicht). Die Marke liegt über der Obergrenze der Schätzung (1050) mit Spielraum, weil die Maße der Geber (873, 968, 732) in je einer Review-Sitzung geprüft wurden (Entscheidung des Orchestrators beim Anspruch 2026-10-16, Vorgabe des Nutzers: Release zuerst, weniger Slices). Schätzung des Architect vom 2026-10-16 in dieser Metrik: Spanne 780 bis 1050, Mitte 910, Obergrenze unter der Marke 1200; Schnitt vor dem Anspruch: nein (Entscheidung nach der Vorgabe des Nutzers: Release zuerst, Mitte unter der Marke). **Meldeschwellen:** (1) Gesamt: über 900 meldet der Implementer dem Architect. (2) Adapter-Teil: Liegt die Summe nach dem Adapter-Teil (DoD-Punkt 1 und 2 mit Fake-Server, Mutationen gefahren, vor dem ersten Commit unter `test/integration`) über 650, meldet der Implementer das dem Architect, bevor er den E2E-Teil beginnt; die Linie des Schnitts ist dann `-upstream` (Adapter, DoD 1 und 2 mit Fake) | `-upstream-e2e` (DoD 3 und der Integrationstest zum Leerlauf), Zuschnitt: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`).
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
- *Ende vor dem ersten `ReadyForQuery`:* bleibt `PGR-E4002`/`PGR-E4005` (`Open` unverändert) und erzeugt keine Session, auch mit `--record-empty-sessions` (Bindung aus `slice-v1-abschluss-sessions`, Nahtstelle 4); Prüfung: `TestRecordAufbauFehlerKeineSession` und am Binary `TestE2ERecordKeineVerbindung` bleiben grün — `LH-FA-25.a` *Kein Ende*, `LH-FA-12.a` *Angenommene Verbindung*.
- *Ende mit `--record-timing` und `--record-empty-sessions`:* nicht Sache des Adapters — `slice-v1-abschluss-verbindungsende-record`.
- *Nach dem Ende gelesene Bytes:* werden nicht gelesen — `LH-FA-25.a` *Aufzeichnen*.
- *Schweregrad:* `V`, leer gilt als fehlend, dann `S`; genau `FATAL`/`PANIC` (`fatal` ist kein Ende) — `LH-FA-25.a`, `SPEC-002` *Randformen des Endes in der Datei*.
- *Sendefehler im Leerlauf:* `PGR-E4003`, die Fehlerantwort im Puffer zählt nicht; aufzeichenbar ist das Ende im Leerlauf, wenn das erste Senden den Upstream erreicht (PostgreSQL 17 schließt ohne Reset: gemessen am 2026-10-16, die Fehlerantwort `57P05` ist beim Senden schon gesendet und bleibt lesbar) — `LH-FA-25.a` *Leerlauf* (Berichtigung; die frühere Zeile „kein Fehler“ widersprach `LH-FA-02.b` *Gelesen*).
- *`Send` nach gelesener `FATAL`/`PANIC`:* kein Fehler, kein Lesen; vorher `PGR-E4003` — `LH-FA-25.a` *Senden nach dem Ende*.
- *Aufruf von `Query`, `Send` oder `Receive` nach `connection_closed`:* Der Port-Kommentar (seit `cf06cfdd`) sagt: `Send` nach gelesener `FATAL`/`PANIC` liefert keinen Fehler, auch wenn das Schreiben scheitert, und liest nichts (`TestSendNachDemEnde`, `TestSendLiestNichtNachDemEnde`; `ClientMessage` ruft `Send` neben dem lesenden `Receive` auf, darum ist es zugesagt). `Query` und `Receive` nach `connection_closed` sind nicht zugesagt — akzeptiertes Negativ: kein Test, der Use Case ruft sie danach nicht — `LH-FA-25.a` *Senden nach dem Ende*.
- *Ende ohne Verbindungsende des Servers (Fake bleibt offen) und Nachrichten hinter dem Ende:* der Adapter wartet nicht und liest nicht — `LH-FA-25.a` *Aufzeichnen*.
- *`PANIC` gegen PostgreSQL:* nicht herstellbar (kein SQL löst ein `PANIC` aus, ohne den Server zu beschädigen); `PANIC` ist nur mit Fake-Server und in der handgeschriebenen `ende-panic.yaml` geprüft — **Akzeptiertes Negativ** (Architect, 2026-10-16: der Schweregrad wird in einem Ausdruck mit `FATAL` gelesen, ein Fall je Wert im Fake-Test genügt).
- *Netzwerkabbruch im E2E:* Mittel ist ein TCP-Proxy im Testprozess (`net.Listen` auf `127.0.0.1:0`, Ziel `PGR_UPSTREAM`), der nach der weitergeleiteten Anfrage beide Verbindungen schließt, ohne eine Antwort zu senden; Zusage `PGR-E4003`, Exit-Code 4, keine Interaktion mit Ende in der Aufzeichnung, `version` 1 — entschieden hier statt im Code (kein Docker-Kommando, kein Container anhalten).

**Risiken:**

- Das Verhalten des Servers (`57P01`, `idle_session_timeout` ab PostgreSQL 14) ist nur gegen die Referenzversion 17 geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×); **Akzeptiertes Negativ** (Entscheidung des Nutzers vom 2026-10-12): Abnahmeszenario 18 läuft nur gegen PostgreSQL 17; die Matrix aus `slice-v1-abschluss-postgres-versionen-matrix` fährt es nicht mit (dort §1 als Abgrenzung, `AGENTS.md` §3.13) — **Ausgang:** weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`.
- Ein Netzwerkabbruch ohne Fehlerantwort braucht im E2E einen Weg, die Verbindung zu kappen: Mittel entschieden (Proxy im Testprozess, siehe Randformen), die Zusage ist `PGR-E4003` statt Ende — **Ausgang:** entfallen — Mittel geliefert (Proxy im Testprozess), `TestE2ERecordNetzwerkabbruch` mit EM6 rot, von der Verifikation erneut gesehen.
- Der Diff liegt an der Grenze einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×): die Rückführung in §4 benennt den Schnitt — **Ausgang:** entfallen — gemessen 941 unter der Marke 1200, eine Review-Sitzung und ein Verifikationslauf trugen den Diff.

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

Die Belege des Implementers stehen zuerst, die Closure des Planners danach (bei der Closure gefüllt, vor dem `git mv` nach `done/`).

### Belege des Implementers

**Stand.** Adapter, Tests und Port-Kommentar an `37e78f76`; die E2E-Tests, der Satz in `SPEC-029` Punkt 2 und die Tabellen an `51193b48`; der Kommentar des Proxys an `d0b77bf4`; danach nur dieser Abschnitt. Geändert sind `internal/adapters/driven/postgres/upstream.go` (`Query`, `Send`, `Receive`, `lies`; neu `beendetVerbindung`, `endeGelesen`), `internal/hexagon/ports/driven/upstream.go` (nur Kommentare), `upstream_test.go` (nachgezogen), neu `upstream_ende_test.go` und `test/integration/ende_e2e_test.go`, `spec/spezifikation.md` (ein Satz in `SPEC-029` Punkt 2), `docs/user/abdeckung-e2e.md` und `abdeckung-unit.md` (`make abdeckung`), dieser Plan (§3). Handbuch und README: nicht berührt (§1, `AGENTS.md` §3.11); die Durchsicht nach Quantoren im Handbuch entfällt. `AGENTS.md` und `.claude` unberührt. Die Berichtigung des Numerus in `SPEC-029` Punkt 2 steht, wie §1 sagt, seit `179ab424` („die Anmeldung mit Klartext-Passwort und die mit MD5“); dieser Slice hat sie nicht noch einmal angefasst, nur den Zeilenumbruch hinter dem neuen Satz gesetzt.

**Umsetzung der Entscheidungen des Architect, soweit sie eine Form im Code haben.** (a) *Schweregrad:* `beendetVerbindung` liest `V`, ohne `V` `S`; `noticeFields` lässt leere Felder weg, darum gilt ein leeres `V` wie ein fehlendes, ohne eigenen Zweig; genau `FATAL` und `PANIC` (`TestEndeSchweregrad`, zehn Zeilen). (b) *Letzte entscheidet:* `lies` setzt den Merker `fehler` bei jeder `ErrorResponse` neu und stuft jede einzeln ein, `ReadyForQuery` leert ihn; `endeGelesen` leitet das Ende aus dem Merker ab, es gibt kein zweites Flag. (c) *`Send`:* Scheitert das Schreiben, liefert `Send` `nil`, wenn `endeGelesen` gilt (der Schreibversuch findet statt), sonst läuft der bisherige Weg (`PGR-E4003` ohne gelesene Fehlerantwort, `PGR-E6001` nach einer anderen). (d) *Nicht auf das Verbindungsende warten, nichts dahinter lesen:* `Query` und `Receive` kehren mit der Fehlerantwort und `connection_closed` zurück, bevor sie weiter lesen; bei `Receive` geht der Rückkehr der Prüfung `ReadBufferLen` voran. (e) *Bestandstests:* Die Zeilen `FATAL` und `PANIC` aus `TestFehlerantwortVorDemAbbruch` sind jetzt ein Ende (`TestEndeSchweregrad`); der Rest der Tests prüft `ERROR`: `TestFehlerantwortVorDemAbbruch` (alle Zeilen, ein neuer Hilfsaufruf `fehlerAntwort`), `TestSendNachFehlerantwort`, und die Zeile „unbekannter Typ nach ErrorResponse“ in `TestNichtLesbareServerantwort` (hinter einer `FATAL` liest der Adapter nicht weiter, die Zeile prüfte sonst nichts). Die Zeile „ReadyForQuery nach der Fehlerantwort“ trägt `ERROR` statt `FATAL`: hinter einer `FATAL` liest der Adapter das `ReadyForQuery` nie (`LH-FA-25.a` *Aufzeichnen*: „was der Upstream danach sendet, liest er nicht“). (f) *Reset mit Sendefehler im Leerlauf:* Der Test schließt die Schreibhälfte der Verbindung (`CloseWrite`); gegen PostgreSQL 17 ist der Fall nicht herstellbar, der Server schließt ohne Reset (§6, gemessen vom Architect), darum kein E2E-Fall dazu.

**Läufe.** Unit-Tests von `postgres`, `model` und `services` im Image `deps` ohne Netz: grün; `-count=10` des Pakets `postgres`: grün (21 s), keine Flocke. `make lint` grün (nach einer Nacharbeit: `gochecknoglobals` traf drei Tabellen der Tests, sie sind jetzt Funktionen, keine Ausnahme im Profil, `AGENTS.md` §3.2). `make docs-check` und `make kopf-check` grün (0 Befunde), `make abdeckung` geschrieben, `make abdeckung-check` grün. E2E-Teilmenge (`tools/test/run-integration-tests.sh` mit `PGR_TESTS_1` und `PGR_TESTS_2`, PostgreSQL 17): `TestE2ERecordEndeTerminate` (`rec.yaml`, `rec.sqlite`), `TestE2ERecordEndeLeerlauf`, `TestE2ERecordNetzwerkabbruch` und in Phase 2 `TestE2EOhnePostgresEnde`: grün. **Voller `make gates`, ein Lauf an `d0b77bf4` (über die Hintergrund-Ausführung des Werkzeugs, nicht aus einer Subshell mit `&`): Exit 0, 13 min 7 s** (22:48:12 bis 23:01:19); `test-integration` mit den vier Tests der Teilmenge `PASS` (Phase 2: `TestE2EOhnePostgresEnde` `PASS`, in Phase 1 `SKIP` wie vorgesehen), `test-postgres-versionen` Version 14 bis 18 grün, `stop-hook-gegenprobe` grün. Danach nur dieser Abschnitt geändert (Planung, Positivliste des Stop-Hooks).

**Wege der Mutanten.** *Unit:* je Mutant eine Kopie von `go.mod`, `go.sum`, `cmd/` und `internal/` im Scratchpad; ein Textersatz auf `upstream.go` (`assert`: genau eine Fundstelle), `go test -count=1 ./internal/adapters/driven/postgres/` im Image `deps` mit `--network=none` und Bind-Mount; nach jedem Mutanten wird die Datei aus dem Original zurückgeschrieben. Stand `37e78f76` (`internal/` ist bis `d0b77bf4` unverändert). Die Ausgabe nennt den roten Test samt Zeile. *E2E:* je Mutant ein frisches Verzeichnis aus `git archive HEAD` (kein mtime-Problem; die Datei bekommt `os.utime`), ein Textersatz, dann `tools/test/run-integration-tests.sh` mit den vier Tests; Stand `51193b48` für EM1 bis EM12, `d0b77bf4` für RM1 und RM2. Mutanten liegen nicht im Repo. Kurznamen: Zeile der Tabelle `Schweregrad` = `TestEndeSchweregrad/<Zeile>/<Weg>`, Wege `Query`, `Receive nach Flush`, `Receive nach Sync` (die letzten zwei verhalten sich gleich: der Adapter unterscheidet `Flush` und `Sync` nicht).

*Unit (Adapter):*

| Zusage | Mutation | roter Test |
|---|---|---|
| `FATAL` ist ein Ende | `FATAL` aus `beendetVerbindung` gestrichen (S1) | `TestEndeSchweregrad` (Zeilen `FATAL`, `V leer, S FATAL`, `V FATAL, S lokalisiert`, je drei Wege), `TestEndeNachFehlerAntwort`, `TestEndeLiestNichtWeiter`, `TestEndeImLeerlauf` (je drei Wege), `TestSendNachDemEnde` |
| `PANIC` ist ein Ende | `PANIC` gestrichen (S2) | `TestEndeSchweregrad` (`PANIC`, `V leer, S PANIC`, je drei Wege), `TestSendNachDemEnde` (beide, `Schreiben scheitert`) |
| ohne `V` zählt `S`, `FATAL` (leeres `V` wie fehlendes) | `S` nie gelesen (S3a), nur `FATAL` aus `S` ausgenommen (S3b) | `TestEndeSchweregrad` (`V leer, S FATAL`; bei S3a auch `V leer, S PANIC`), `TestSendNachDemEnde` |
| ohne `V` zählt `S`, `PANIC` | nur `PANIC` aus `S` ausgenommen (S3c) | `TestEndeSchweregrad` (nur `V leer, S PANIC`), `TestSendNachDemEnde` (nur diese Zeile) |
| `V` geht vor `S` | nur `S` gelesen (S4); `V` oder `S` (S6) | S4: `TestEndeSchweregrad` (`V FATAL, S lokalisiert`, `V ERROR, S FATAL`, `V klein geschrieben`); S6: `V ERROR, S FATAL`, `V klein geschrieben` |
| genau `FATAL` und `PANIC` (`fatal` ist kein Ende) | Groß-/Kleinschreibung ignoriert (S5); leerer Schweregrad zählt als Ende (S7) | S5: `V klein geschrieben`, `V leer, S klein geschrieben`; S7: `V leer, S leer` (je drei Wege und `TestSendNachDemEnde`) |
| die letzte Fehlerantwort entscheidet (`ERROR`, dann `FATAL`; Lesen) | Ende nur, wenn es die erste Fehlerantwort ist (E1) | `TestEndeNachFehlerAntwort` (drei Wege) |
| … (Senden) | der Merker behält die erste Fehlerantwort (E2) | `TestSendNachDemEnde/ERROR dann FATAL/Schreiben scheitert`, `TestFehlerantwortVorDemAbbruch/zwei Fehlerantworten, die letzte zählt` |
| `ERROR` und Abbruch ist `PGR-E6001` mit SQLSTATE und Meldung | der Merker bleibt unbeachtet, immer `PGR-E4003` (V1) | `TestFehlerantwortVorDemAbbruch` (alle `ERROR`-Zeilen), `TestEndeSchweregrad` (die Zeilen ohne Ende), `TestSendNachFehlerantwort/ErrorResponse gelesen` |
| ohne Fehlerantwort `PGR-E4003` | immer `PGR-E6001` (V2) | `TestFehlerantwortVorDemAbbruch/ohne Fehlerantwort`, `/ReadyForQuery nach der Fehlerantwort`, `TestReceiveFehler/Verbindungsende`, `TestSendNachFehlerantwort/ohne ErrorResponse`, `TestEndeSendefehlerImLeerlauf` (Receive) |
| der Merker wird bei `ReadyForQuery` zurückgesetzt | Rücksetzen entfernt (M1) | `TestFehlerantwortVorDemAbbruch/ReadyForQuery nach der Fehlerantwort`, `TestSendNachFehlerantwort/ErrorResponse und ReadyForQuery gelesen` |
| `Query` wartet nicht auf das Verbindungsende (der Fake bleibt offen) | `Query` liest nach dem Ende bis `EOF` (Q1); liest bis zu einem `ReadyForQuery` weiter (Q2) | `TestEndeSchweregrad` (die fünf Ende-Zeilen)/`Query`, `TestEndeNachFehlerAntwort/Query`, `TestEndeLiestNichtWeiter/Query`, `TestEndeImLeerlauf/Query` (Frist zwei Sekunden) |
| … `Receive` | dasselbe in `Receive` (R1, R2) | dieselben Zeilen mit den Wegen `Receive nach Flush`, `Receive nach Sync` |
| `Query` liest nichts hinter der Fehlerantwort (ein `ReadyForQuery` und eine nicht lesbare Nachricht im selben Schreibvorgang) | `Query` liest den Lesepuffer leer (Q5) | `TestEndeLiestNichtWeiter/Query` (nur diese Zeile; die Nachricht hinter dem `ReadyForQuery` ist nicht lesbar und lässt das Lesen scheitern) |
| … `Receive` (Lesepuffer) | `Receive` liest den Lesepuffer leer (R5) | `TestEndeLiestNichtWeiter/Receive nach Flush`, `/Receive nach Sync` (nur diese) |
| die Antworten enden mit `connection_closed` (`Query`, `Receive`) | `Query` liefert die Fehlerantwort ohne `connection_closed` (Q4); `Receive` (R4) | dieselben Zeilen wie Q1/R1, je im eigenen Weg (`TestEndeSchweregrad`, `TestEndeNachFehlerAntwort`, `TestEndeLiestNichtWeiter`, `TestEndeImLeerlauf`) |
| `Query`, `Receive`, `Send` je eigener Fall | Ende-Zweig nur in `Query` entfernt (Q3); nur in `Receive` (R3); `Send` ohne `endeGelesen` (D1) | Q3: nur die Zeilen `/Query`; R3: nur die Zeilen `/Receive …`; D1: nur `TestSendNachDemEnde` |
| Leerlauf: `NoticeResponse` und `ParameterStatus` vor der Fehlerantwort bleiben (`Query`) | `Query` verwirft sie beim Ende (N1) | `TestEndeImLeerlauf/Query` (nur diese Zeile) |
| … (`Receive`) | `Receive` verwirft sie (N2) | `TestEndeImLeerlauf/Receive nach Flush`, `/Receive nach Sync` |
| `Send` liest nach dem Ende nichts (eine Nachricht des Servers wartet, das Schreiben scheitert) | `Send` liest nach der gelesenen `FATAL` mit 30 ms Frist eine Nachricht und verwirft sie (X1); dasselbe ohne Frist (X2) | `TestSendLiestNichtNachDemEnde` (nur dieser Test; `Read`-Zähler, Frist zwei Sekunden) |
| `Send` nach gelesener `FATAL`/`PANIC` liefert keinen Fehler, auch wenn das Schreiben scheitert | die Prüfung entfällt (D1); `Send` liefert dann `PGR-E4003` (D6) | `TestSendNachDemEnde/<Zeile mit Ende>/Schreiben scheitert` (`FATAL`, `PANIC`, `V leer, S FATAL`, `V leer, S PANIC`, `V FATAL, S lokalisiert`, `ERROR dann FATAL`) |
| nach einer anderen gelesenen Fehlerantwort bleibt der Sendefehler `PGR-E6001` | jede gelesene Fehlerantwort macht den Sendefehler zum Nicht-Fehler (D3) | `TestSendNachDemEnde` (Zeilen ohne Ende, `Schreiben scheitert`), `TestSendNachFehlerantwort/ErrorResponse gelesen` |
| ohne gelesene Fehlerantwort bleibt der Sendefehler `PGR-E4003` (`Send`, im Leerlauf mit gepufferter `FATAL`) | `Send` liefert bei jedem Sendefehler `nil` (D2); ohne gelesene Fehlerantwort `nil` (D7); `Send` liest die gepufferte Fehlerantwort nach dem Sendefehler (D4) | D2/D7: `TestEndeSendefehlerImLeerlauf/Receive nach Flush`, `/Receive nach Sync`, `TestSendNachFehlerantwort/ohne ErrorResponse`, `/ErrorResponse und ReadyForQuery gelesen`; D4: `TestEndeSendefehlerImLeerlauf/Receive nach Flush`, `/Receive nach Sync` |
| … (`Query`): die gepufferte Fehlerantwort zählt nicht | `Query` liest nach dem Sendefehler und liefert das Ende (D5) | `TestEndeSendefehlerImLeerlauf/Query` (nur diese Zeile) |

*E2E (Binary gegen PostgreSQL 17; RM nur in Phase 2 rot):*

| Zusage | Mutation | roter Test |
|---|---|---|
| die Aufzeichnung trägt `meta.version` 2 (SQLite) | Version wird nicht gesetzt (EM1) | `TestE2ERecordEndeTerminate/rec.sqlite` (nur diese Zeile) |
| … `version: 2` (YAML) | der Schreiber wählt Version 1 (EM2) | `TestE2ERecordEndeTerminate/rec.yaml`, `TestE2ERecordEndeLeerlauf/mit_Anfrage_danach` |
| der Client erhält die Antworten bis einschließlich der Fehlerantwort | der Adapter liefert nur die Fehlerantwort (EM9) | `TestE2ERecordEndeTerminate` (beide Formate) |
| die Aufzeichnung trägt `connection_closed` hinter der Fehlerantwort | der Adapter liefert kein `connection_closed` (EM5) | `TestE2ERecordEndeTerminate` (beide), `TestE2ERecordEndeLeerlauf/mit_Anfrage_danach` |
| der Exit-Code bleibt 0, ohne Meldung | der Adapter liefert zum Ende einen Fehler (EM10) | `TestE2ERecordEndeTerminate` (beide), `TestE2ERecordEndeLeerlauf/mit_Anfrage_danach` |
| Formvergleich: Feld `V`, `S`, `C`, `M` der Fehlerantwort der Aufzeichnung wie in den Fixtures | der YAML-Schreiber lässt das Feld weg (EM4v, EM4s, EM4c, EM4m; die Leitung bleibt unberührt) | `TestE2ERecordEndeTerminate/rec.yaml` (`gleicheForm`), `TestE2ERecordEndeLeerlauf/mit_Anfrage_danach`; EM4c auch `…/rec.sqlite` (derselbe Abbilder) |
| … dasselbe auf der Leitung | der Adapter lässt das Feld weg (EM3a für `V`, EM3b für `S`, EM3c für `C`: beide Formate durch `pruefeBeendeAntworten`; EM3d für `M`: nur `rec.yaml`, die Leitung prüft `M` nicht) | `TestE2ERecordEndeTerminate`, `TestE2ERecordEndeLeerlauf/mit_Anfrage_danach` |
| Leerlauf: die Fehlerantwort gehört zur nachfolgenden Anfrage | `Query` verwirft vorab gesendete Nachrichten (EM7) | `TestE2ERecordEndeLeerlauf/mit_Anfrage_danach` (nur diese Zeile) |
| ein Netzwerkabbruch wird nicht als Ende aufgezeichnet | der Adapter wertet das Verbindungsende ohne Fehlerantwort als Ende (EM6) | `TestE2ERecordNetzwerkabbruch` (nur dieser Test) |
| Replay ohne PostgreSQL: der Client erhält das Ende seiner Verbindung | Replay schließt die Verbindung nach dem Ende nicht (RM1) | `TestE2EOhnePostgresEnde/rec.yaml`, `/rec.sqlite` (Frist des Lesens 10 s) |
| … dieselben Antworten | Replay lässt die Datenzeile aus (RM2) | `TestE2EOhnePostgresEnde/rec.yaml`, `/rec.sqlite` |

**Grüne Mutanten, Äquivalente, Nicht belegt (Grenze).** (1) *EM12* (`gekappt` im Proxy des Tests entfernt): grün; die Verbindungen werden unmittelbar nach dem Schreiben geschlossen, das Flag schließt nur ein Zeitfenster. Der Kommentar nennt es als Grenze und sagt nicht zu, dass eine Antwort nicht zurückgeht (`AGENTS.md` §3.11). (2) *`Send` liest nichts nach dem Ende:* kein Äquivalent (Review F-1, die frühere Aussage „äquivalent“ ist gestrichen); `TestSendLiestNichtNachDemEnde` zählt die `Read`-Aufrufe der Verbindung, lässt das Schreiben scheitern, lässt eine weitere Nachricht des Servers warten und gibt `Send` zwei Sekunden (`fertigBinnen`); X1 und X2 rot (Tabelle). (3) *Leerlauf ohne nachfolgende Anfrage* (`TestE2ERecordEndeLeerlauf/ohne_Anfrage_danach`): Der Test beobachtet das Nichtgeschehen; ein Mutant im Adapter ist nicht herstellbar, ohne dass der Adapter außerhalb von `Query` und `Receive` liest, und das Ergebnis ginge nur durch einen Use Case in die Aufzeichnung (Ort: `services/record.go`, Geber `-record`). Der Test hält das heutige Verhalten fest; er ist keine Mutation wert. (4) *`PANIC` gegen PostgreSQL:* nicht herstellbar, akzeptiertes Negativ (§6); `PANIC` ist in `TestEndeSchweregrad`, `TestSendNachDemEnde` und `ende-panic.yaml` (Replay) geprüft. (5) *Reset des Upstreams im Leerlauf:* nur im Adapter-Test (Schreibhälfte geschlossen), nicht gegen PostgreSQL (§6, Punkt Sendefehler im Leerlauf). (6) *Extended gegen PostgreSQL:* gehört `slice-v1-abschluss-verbindungsende-upstream-extended` (§1); der Adapter ist für `Send` und `Receive` mit Fake-Server geprüft. (7) *Die Probe der Version in SQLite* liest Bytes: `version2` darf stehen und `version1` nicht; sie hängt daran, dass SQLite das Feld beim `UPDATE` an Ort und Stelle überschreibt (am Stand gemessen, EM1 lässt `version1` stehen und ist rot). (8) *Replay der Datei:* verglichen werden Typ und Rumpf von `RowDescription` und `DataRow` und `V`, `S`, `C`, `M` der Fehlerantwort; Reihenfolge und Zusatzfelder (`F`, `L`, `R`) zählen nicht (Vorgabe des Architect).

**Messung** (Metrik aus §4: `git diff --numstat c5089444..12d7be91 -- internal/adapters/driven/postgres internal/hexagon/ports/driven test/integration docs/user`, Einfügungen plus Löschungen). Nach der Nacharbeit zum Review (F-1, F-2): **941** (884 Einfügungen, 57 Löschungen), **über der Meldeschwelle 900** (dem Architect gemeldet), unter der Marke 1200. Vor der Nacharbeit 875 an `d0b77bf4`. Meldeschwelle des Adapter-Teils: nach dem Adapter-Teil (Commit `37e78f76`) **483** (428 Einfügungen, 55 Löschungen; im ersten Stand stand 475, Review F-3), unter 650, nicht gemeldet. Aufschlüsselung jetzt: Produktcode (`upstream.go`) 65, Port-Kommentar 14, Tests des Adapters 469 (`upstream_ende_test.go` 391, `upstream_test.go` 78), E2E 370, Tabellen 23. Schätzung des Architect: Mitte 910, Spanne 780 bis 1050.

**Kommentare** (Schritt 20, `AGENTS.md` §3.7, §3.11). Jeder neue oder geänderte Kommentar beschreibt den Ist-Zustand im Indikativ; zu jeder Zusage darin steht ein Test: Port-Kommentar (Ende des Servers, `Query`, `Receive`, `Send`: `TestEndeSchweregrad`, `TestEndeLiestNichtWeiter`, `TestSendNachDemEnde`, `TestEndeSendefehlerImLeerlauf`; der Satz „nur noch `Close`“ ist gestrichen (Review F-2: `ClientMessage` ruft `Send` neben dem lesenden `Receive`); der Kommentar sagt zu, was getestet ist, und von `Query` und `Receive` nach `connection_closed` nichts. Hinweis an den Architect: Derselbe Satz steht in `LH-FA-25.a` *Senden nach dem Ende*; ich habe die Spezifikation nicht geändert), `Query`, `Send`, `Receive`, `lies`, `beendetVerbindung` (`TestEndeSchweregrad`), `endeGelesen` (`TestSendNachDemEnde`), Kommentare der Tests (je ein Test) und der Proxy (Grenze statt Zusage, siehe EM12). Altes Wort: „bis einschließlich ReadyForQuery“ steht in `ports/driving/record.go` Zeile 39 für den Use Case, nicht für den Upstream-Port, und bleibt; sonst kein Treffer in `internal/` für den geänderten Satz des Port-Kommentars.

**Beobachtungen.** (1) Sieben Bestandstests (`TestFehlerantwortVorDemAbbruch` mit drei Zeilen, „zwischen zwei Interaktionen“, „Extended“, „ReadyForQuery nach der Fehlerantwort“, `TestSendNachFehlerantwort`, `TestNichtLesbareServerantwort`) nutzten `FATAL` und `PANIC` als Beispiel für einen *gewöhnlichen* Fehler; mit dem Vertrag des Endes prüften sie etwas anderes oder nichts. Das Gate zeigte sie rot (außer der Zeile in `TestNichtLesbareServerantwort`, die grün geblieben wäre und nichts mehr geprüft hätte: erst das Lesen der Zeile fiel auf). Klasse für den Planner: ein Beispiel im Bestand benutzt eine Ausprägung, die ein neuer Vertrag umdeutet; Zuordnung zu einer Klasse des Registers beim Planner (etwa `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, Gegenseite). (2) Die Mutanten EM3a bis EM3c (Feld auf der Leitung weggelassen) färbten beide Formate rot, nicht nur den Formvergleich: Die Leitungsprüfung des Tests (`pruefeBeendeAntworten`) deckt `V`, `S`, `C` ab; erst die Mutanten im Schreiber (EM4) belegen den Formvergleich allein. (3) `gochecknoglobals` traf drei Tabellen der Tests; sie sind Funktionen geworden (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` nicht berührt, kein neuer Eintrag). (4) Drei Wege (`Query`, `Receive nach Flush`, `Receive nach Sync`) in einer Tabelle halfen gegen `dupl`; `Flush` und `Sync` verhalten sich im Adapter gleich, die Zeilen belegen es.

**Lauf nach dem Review.** Voller `make gates`, ein Lauf an `bf6487cb` (Hintergrund-Ausführung des Werkzeugs): Exit 0, 13 min 8 s (23:19:39 bis 23:32:47). Ein erster Lauf an `5430db56` endete nach einer Sekunde rot an `abdeckung-check`, weil ich den Kommentar des Tests nach `make abdeckung` geändert hatte; behoben mit `bf6487cb` (Tabelle neu).

**Ausgang der Befunde des Reviews** (`docs/reviews/2026-10-16-review-slice-v1-abschluss-verbindungsende-upstream.md`): F-1 behoben (`cf06cfdd`, `12d7be91`, X1/X2 rot). F-2 behoben im Port-Kommentar, Spezifikation unberührt (Hinweis an den Architect). F-3 berichtigt (483). F-4 akzeptiert: Der Plan folgte der Zeile von `upstream.go` einen Commit später; keine Aktion. F-5 Hinweis: Die Schweregrad-Regel steht an drei Stellen (Modell, `play.go`, Adapter); der Adapter darf das Modell nicht importieren, die Zeilen von `TestEndeSchweregrad` pinnen seine Kopie. F-6 Hinweis: Die Probe der SQLite-Version liest Bytes (siehe Grenze 7). F-7 Hinweis: Der Formvergleich gegen die Fixtures läuft nur für YAML; die Fehlerantwort der SQLite-Aufzeichnung prüft der Replay in Phase 2 über `endeSicht` (V, S, C, M).

### Closure des Planners (2026-10-16)

**Gate-Nachweis.** Voller `make gates`, Exit 0, an `bf6487cb` (Implementer, Belege oben). `git diff --name-only bf6487cb..8dd4d982` nennt vor dieser Closure nur den Plan und den Verifikationsbericht (Positivliste); die Verifikation fuhr `make gates-schnell` (Exit 0) und `make abdeckung-check`. `make kopf-check` und `make docs-check` laufen nach der Closure nacheinander.

**Gegenstand:** geliefert (kein Fall von „übernommen“ oder „entfallen“): der Upstream-Adapter stuft `FATAL`/`PANIC` als Ende ein (`Query`, `Send`, `Receive`), Ende im Leerlauf, Abnahmeszenario 18 am Binary gegen PostgreSQL 17 (YAML, SQLite, Replay ohne PostgreSQL, Netzwerkabbruch), der Satz in `SPEC-029` Punkt 2.

**Review und Verifikation.** Review (`docs/reviews/2026-10-16-review-slice-v1-abschluss-verbindungsende-upstream.md`): 0 HIGH, 1 MEDIUM, 3 LOW, 3 INFO. Verifikation (`docs/reviews/2026-10-16-verify-slice-v1-abschluss-verbindungsende-upstream.md`): bestanden, V-1 LOW, V-2 bis V-4 INFO; sieben Mutanten selbst gefahren, alle rot wie in §7. Ausgang je Befund (Nutzerentscheid vom 2026-10-16: LOW und INFO ohne Zusage ohne Test sind „akzeptiert, mit Ausgang in §7“):

| Befund | Ausgang |
|---|---|
| F-1 (MEDIUM) `Send` liest nichts: Äquivalenz falsch | behoben `cf06cfdd`, `12d7be91` (`TestSendLiestNichtNachDemEnde`, X1 und X2 rot); Register `gruener-mutant-faelschlich-aequivalent` |
| F-2 (LOW) Port-Kommentar „nur noch `Close`“ | behoben `cf06cfdd`: der Kommentar sagt zu, was getestet ist |
| V-2 (INFO) derselbe Satz in der Spezifikation | **Hinweis an den Architect, keine Änderung durch den Planner, kein Folge-Slice:** `spec/spezifikation.md`, `LH-FA-25.a` *Senden nach dem Ende*: „Nach dem Ende ruft der Use Case den Upstream nur noch zum Schließen auf“ hat keinen Test und widerspricht F-2 (`ClientMessage` ruft `Send` neben `Receive` auf). Der Architect fasst den Satz enger oder liefert den Test (`AGENTS.md` §3.11); bis dahin sagt das Handbuch (Schritt 42) nichts dazu |
| F-3 (LOW) Zahl 475 statt 483 | berichtigt in §7 |
| F-4 (LOW) §3 folgte dem Code einen Commit später | akzeptiert: behoben im folgenden Commit, keine Aktion; Register `plan-folgt-korrektur-nicht` zählt diesen Slice einmal (mit V-1) |
| V-1 (LOW) Plan §6 „nur `Close`“ | behoben mit dieser Closure: Zeile in §6 auf den Stand des Port-Kommentars gezogen (`Send` zugesagt, `Query`/`Receive` nach `connection_closed` nicht zugesagt) |
| F-5 (INFO), V-4 Schweregrad-Regel an drei Stellen | akzeptiert: der Adapter darf das Modell nicht importieren (`make a-check`), `TestEndeSchweregrad` pinnt seine Kopie; Hinweis an den Architect, kein Folge-Slice |
| F-6 (INFO) Versionsprobe über Rohbytes | akzeptiert: Grenze 7 in §7, EM1 rot gesehen |
| F-7 (INFO), V-3 Formvergleich nur für YAML | akzeptiert: die SQLite-Fehlerantwort prüft der Replay der Phase 2 (`endeSicht`: `V`, `S`, `C`, `M`); DoD-Punkt 3 gedeckt (Verifikation) |

**Entscheidung nach der Meldung (Meldeschwelle 900).** 2026-10-16: Der Implementer meldete 941 an den Architect; der Orchestrator entschied Fortsetzung ohne Schnitt (Vorgabe des Nutzers: Release zuerst, weniger Slices). Die Marke 1200 wurde nicht erreicht, die Rückführung trat nicht ein. Gemessen 941 (Produktcode 65, Port-Kommentar 14, Tests des Adapters 469, E2E 370, Tabellen 23).

**Risiken aus §6 — Ausgang je Risiko** (die Zeilen stehen in §6 am Risiko): (1) Verhalten nur gegen PostgreSQL 17 *weiter offen* → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (Akzeptiertes Negativ, Entscheidung des Nutzers vom 2026-10-12; Grenze in `SPEC-029` Punkt 2); (2) Weg für den Netzwerkabbruch *entfallen*: Proxy im Testprozess geliefert, `TestE2ERecordNetzwerkabbruch` mit EM6 rot (Implementer und Verifikation); (3) Diff an der Grenze *entfallen*: 941 < 1200, eine Review-Sitzung trug ihn. Ob die Ausgänge tragen: (1) an der Zeile im Register und an der Grenze in der Spezifikation, (2) am Test, (3) an der Messung.

**Register** (Zähler = Dateien unter `evidence/`): `gruener-mutant-faelschlich-aequivalent` 4 → 5 (F-1; Stand `offen`, über der Schwelle ohne Verkörperung und ohne Zuweisung; **die Welle-Closure von welle-v1-abschluss muss ihm einen Ausgang geben**, Nutzerentscheidung: der Planner verkörpert und plant nichts); `plan-folgt-korrektur-nicht` 55 → 56 (V-1; `verkörpert`, bleibt); `zusage-im-kommentar-weiter-als-pruefung` 67 → 68 (F-2, V-2; `verkörpert`, bleibt); `serververhalten-nur-gegen-eine-version-geprueft` 5 → 6 (Risiko 1, *weiter offen*; `verkörpert`, bleibt). Vermerke ohne Zuwachs: `handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (Szenario 18, einfache Anfrage, jetzt gegen PostgreSQL 17 belegt; Extended gegen PostgreSQL bleibt `slice-v1-abschluss-verbindungsende-upstream-extended`, `geplant`), `schaetzung-in-zeilen-unterschaetzt-tests` (941 in der Spanne 780 bis 1050: kein Auftreten), `meldeschwelle-ohne-artefakt` (Schwelle stand in §4, die Entscheidung steht als Zeile mit Datum oben: kein Auftreten). Kein Eintrag erreicht neu 3× ohne Verkörperung und ohne Zuweisung.

**Lerneintrag — keine geänderte Regel, zwei Beobachtungen.** (1) Eine Äquivalenz-Aussage, deren genannte Invariante die Zusage des Vertrags selbst ist („`Send` enthält keinen Lesezugriff“ gegen „es wird nichts mehr gelesen“), ist keine Äquivalenz, sondern eine ungetestete Zusage; fünfter Fall der Klasse `gruener-mutant-faelschlich-aequivalent`, erstmals bei einer Zusage der Spezifikation statt einer Eigenschaft des Modells. Die Nutzerfrage dazu (Invariante und Ort nennen, sonst Fall) steht im Eintrag und gehört zur Welle-Closure. (2) Ein Satz über den Aufrufer („ruft nur noch `Close`“) in Port-Kommentar und Spezifikation ist eine Zusage ohne Test, wenn kein Gate den Aufrufer prüft; der Kommentar wurde auf das Getestete enger gefasst, die Spezifikation bleibt beim Architect (`AGENTS.md` §3.11 trägt). Neuer Sensor: keiner. Spec-Lücke benannt: V-2.

**Paarungen.** (a) Anker: Der Lerneintrag trägt kein Feld `liegt in`; nichts zu prüfen. (b) Folge-Slice: Kein Risiko hat den Ausgang *eingetreten*; die Abgrenzung der Klasse 1 nennt `slice-v1-abschluss-verbindungsende-upstream-extended` (`next/`, nennt die Sendung und diese Kennung in §1 und DoD-Punkt 1 und 2) und `slice-v1-abschluss-verbindungsende-doku` (`next/`, Voraussetzung Schritt 41). (c) Register: Alle genannten `BEO-REPO/…` bestehen als Verzeichnis mit nicht leerem `evidence/`.

**Schätzung des nächsten Schritts, `slice-v1-abschluss-verbindungsende-doku` (Schritt 42):** zwei Liefer-Punkte (Handbuch, README), eine Schicht (Nutzer- und Wartungs-Doku); Schätzung des Planners 150 bis 300 Zeilen Text, mit den Beispiel-Dateien der Probe und der Nacharbeit rund 200 bis 450, Mitte 300. Eine Marke steht in §4 noch nicht; der Planner trägt beim Anspruch Metrik (`AGENTS.md` §3.14, Pfade `docs/user README.md` und die Beispiel-Dateien) und Marke 600 mit Meldeschwelle 480 ein; die Obergrenze liegt darunter, die Nutzerpraxis schneidet nur, wenn die Schätzung die Marke überschreitet: kein Schnitt. Ein Vorlauf des Architect ist nicht nötig: kein neuer Vertrag, §6 nennt keine Randform (`AGENTS.md` §3.12), die Probe am Binary ist Sache des Implementers. Zu beachten: Handbuch und README sagen zu `Send` nach dem Ende und zu „nur noch Schließen“ nichts (V-2), und zum Aufzeichnen einer Extended-Interaktion mit Ende gegen PostgreSQL nichts (Bindung aus `-upstream-extended`).

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

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-sessions` (Nahtstelle 4, Planner, 2026-10-12, `AGENTS.md` §3.13): Sie verlangt keinen neuen Code (`Open` bleibt unverändert) und keinen neuen Test, nur dass vorhandene grün bleiben; sie gehört zu DoD-Punkt 1 und fügt weder Liefer-Punkt noch Schicht hinzu: Liefer-Punkte weiter 3, Schichten weiter zwei (Upstream-Adapter; Kern nur als Vertrag des Ports).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-postgres-versionen-zusage` (Architect, 2026-10-15, `AGENTS.md` §3.13): ein Satz in `SPEC-029` Punkt 2 im Commit des Tests von Szenario 18; er gehört zu DoD-Punkt 3 und fügt weder Liefer-Punkt noch Schicht hinzu (die Spezifikation zählt nicht als Schicht, Planung ebenso nicht): Liefer-Punkte weiter 3, Schichten weiter zwei (Upstream-Adapter; Kern nur als Vertrag des Ports); Schätzung unverändert, rund zwei Zeilen mehr.

Nachgezählt beim Eintragen der Berichtigung I-1 aus `slice-v1-abschluss-postgres-versionen-zusage` (Planner, 2026-10-15, `AGENTS.md` §3.13): ein Wort im Satz von `SPEC-029` Punkt 2, den dieser Slice ohnehin im Commit seines Tests schreibt; sie gehört zu DoD-Punkt 3 und fügt weder Liefer-Punkt noch Schicht hinzu (die Spezifikation zählt nicht als Schicht, Planung ebenso nicht). Liefer-Punkte weiter 3, Schichten weiter zwei (Upstream-Adapter; Kern nur als Vertrag des Ports); Schätzung unverändert.

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-verbindungsende-format-schreiber` (Planner, 2026-10-15, `AGENTS.md` §3.13): Der YAML-Schreiber ist Voraussetzung und keine Lieferung dieses Slice; er gehört zu DoD-Punkt 3 (Abnahmeszenario 18) und fügt weder Liefer-Punkt noch Schicht hinzu. Liefer-Punkte weiter 3, Schichten weiter zwei (Upstream-Adapter; Kern nur als Vertrag des Ports); Schätzung unverändert.

Nachgezählt beim Abgeben an `slice-v1-abschluss-verbindungsende-upstream-extended` (Planner, 2026-10-16, `AGENTS.md` §3.13): Der Extended-Beleg gegen PostgreSQL liegt nicht in diesem Slice (§1); die Sendung aus `-replay` (Fixtures) bleibt für die einfache Interaktion in DoD-Punkt 3. Es entsteht kein Liefer-Punkt und keine Schicht: Liefer-Punkte weiter 3, Schichten weiter zwei; die Schätzung 700 bis 1200 bleibt unverändert (sie enthielt den Extended-Fall nicht).

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Planner, 2026-10-15, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-sqlite`): Die Sendung (E2E-Beleg der SQLite-Aufzeichnung, `meta.version` 2) bleibt in DoD-Punkt 3, nur der Geber heißt jetzt anders; die Voraussetzung des Starts nennt beide Hälften. Kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 3, Schichten weiter zwei.

Nachgezählt beim Nachziehen der Stelle von `slice-v1-abschluss-verbindungsende-upstream-extended` (§1; Planner, 2026-10-16, `AGENTS.md` §3.13, Verschiebung hinter das Release, Entscheidung des Nutzers): Der Folge-Slice nimmt die Sendung weiter an (nicht in `done/`, DoD-Punkt 1 und 2 und §1 nennen sie); dieser Slice trägt nichts Neues ein: Liefer-Punkte weiter 3, Schichten weiter zwei, Schätzung 700 bis 1200 unverändert.

Nachgezählt bei der Closure von `slice-v1-abschluss-verbindungsende-record` (Planner, 2026-10-16, `AGENTS.md` §3.13): Das Risiko „Fakes statt echtem Upstream“ geht mit dem Ausgang *eingetreten* an diesen Slice; die Sendung steht schon in §1 (*Übernommen aus `slice-v1-abschluss-verbindungsende-record`*: liefert, was dessen Fakes bisher nur abbilden) und in DoD-Punkt 1 und 3. Die Test-Idee „Schreibfehler in `serverRichtung`“ und `EndTerminate` nach dem Ende gehen **nicht** hierher (dort gestrichen, kein Kommentar sagt zu): kein Zuwachs. Liefer-Punkte weiter 3, Schichten weiter zwei (Upstream-Adapter; Kern nur als Vertrag des Ports); Schätzung 700 bis 1200 unverändert. Maße des Gebers für den Anspruch: 873, Tests zu Produktcode 5,2; der Planner trägt beim Anspruch die Metrik aus `AGENTS.md` §3.14 und die Meldeschwelle in §4 ein.

Nachgezählt bei der Prüfung des Architect vor dem Anspruch (2026-10-16, `AGENTS.md` §3.13, §3.14): Die Sendungen aus 37 bis 40 stehen in §1 und DoD-Punkt 3 (handgeschriebene Aufzeichnung ohne Beleg: der Formvergleich in DoD-Punkt 3; der Satz in `SPEC-029` Punkt 2 zu Szenario 18 im Commit des Tests, die Berichtigung des Numerus steht seit `179ab424`). Liefer-Punkte weiter 3 (Ende erkennen; Leerlauf; Abnahmeszenario 18). Schichten zwei: Upstream-Adapter und `docs/user` (die Abdeckungstabellen, konservativ gezählt); der Kommentar im Port ist keine Logik und keine Schicht (wie reine Verdrahtung im Bootstrap, `AGENTS.md` §3.13), `test/integration` zählt als Tests des Adapters. Schätzung in der Metrik aus `AGENTS.md` §3.14 (Pfade siehe §4): Produktcode und Port-Kommentar rund 85, Tests des Adapters rund 500 (rund 350 neu, rund 120 für die umgeschriebenen `TestFehlerantwortVorDemAbbruch` und `TestSendNachFehlerantwort`, Hilfen des Fake-Servers eingeschlossen), E2E rund 315 (Szenario 18 in beiden Formaten 75, Formvergleich 60, Leerlauf 60, Netzwerkabbruch mit Proxy 90, Hilfen 30), Tabellen rund 12: Mitte rund 910, Spanne 780 bis 1050. Schnitt: nein (§4).
