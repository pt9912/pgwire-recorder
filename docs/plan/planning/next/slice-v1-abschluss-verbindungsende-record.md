# Slice slice-v1-abschluss-verbindungsende-record: Verbindungsende im Record: Use Case und Zustellung an den Client

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-02.b` · `LH-FA-12.a` · `LH-FA-18.a` · `LH-FA-21.a` · `SPEC-002` · `ARC-002` · `ARC-006`

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

**Ziel:** Der Record-Service übernimmt eine Interaktion, deren Antworten ein Ende des Servers zeigen (`connection_closed` als letzte Antwort, aus dem Upstream-Port), als letzte Interaktion der Session — bei einer einfachen Anfrage wie bei einer Extended-Interaktion — und der PGWire-Adapter stellt dem Client die Antworten bis einschließlich der Fehlerantwort zu und schließt dann beide Verbindungen, ohne `ReadyForQuery`, ohne Meldung und mit unverändertem Exit-Code (`LH-FA-25.a` *Aufzeichnen*, *Form*, *Stelle*). Ein Ende in einer Gruppe, die der Client nicht mit `Sync` oder `Flush` abgeschlossen hat, ist `PGR-E6001` (die Session wird nicht übernommen). Nachweis mit Fake-Upstream und Fake-Recorder; der echte Upstream und Abnahmeszenario 18 folgen in `slice-v1-abschluss-verbindungsende-upstream`.

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): das Aufzeichnen des Endes im Use Case und im PGWire-Adapter. Jener Slice liefert `connection_closed` und seine Formregeln im Modell (`Validate` nimmt die Interaktion mit Ende an); dieser benutzt sie.

**Übernommen aus `slice-v1-abschluss-verbindungsende-replay`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12): der Record-Pfad des PGWire-Adapters. Die Abbildung von `connection_closed` auf „keine PGWire-Nachricht“ (`toMessage`, `send`) liefert jener Slice; dieser ruft sie auf und schließt in `richtungen` beide Verbindungen.

**Übernommen aus `slice-v1-abschluss-sessions`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12): Das Ende des Servers als Ereignis einer Session. `slice-v1-abschluss-sessions` regelt Verbindungsende und Verbindungsfehler ohne das Ende des Servers; die Reaktion des Use Case auf ein Ende nach `LH-FA-25.a` liegt hier (die Einstufung von `PGR-E4003` und `PGR-E6001` im Adapter in `slice-v1-abschluss-verbindungsende-upstream`), und `CloseSession`, `--record-empty-sessions` und die Kennung der Session (Nummer in der Reihenfolge der Sessionenden) sind nach jenem Slice (Schritt 26) unverändert zu benutzen.

**Gebunden durch `slice-v1-abschluss-sessions`** (dort §1, *Nahtstelle zu den Verbindungsende-Slices*, Punkt 3 und 4; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): Dieser Slice hält ein, was `slice-v1-abschluss-sessions` im Record-Service liefert:

3. *`CloseSession` entscheidet allein, ob und unter welcher Kennung eine Session geschrieben wird* (Kennung = Zahl der geschriebenen Sessions + 1, unter dem Sperrzustand des Services). Dieser Slice fügt die Art des Endes (Ende des Servers) hinzu und ändert nicht, dass eine Session mit nicht unterstützter Interaktion und das Zwangsende im Aufbau (`EndForced` ohne laufende und ohne abgeschlossene Interaktion, `wirdVerworfen` in `internal/hexagon/services/record.go`) nie geschrieben werden und jede andere Session mit `--record-empty-sessions` geschrieben wird. Eine Session mit Ende trägt mindestens die Interaktion mit dem Ende und ist deshalb nie eine leere Session.
4. *Eine Verbindung wird erst mit eröffneter Session eine Session:* Ein Ende des Servers vor dem ersten `ReadyForQuery` (`LH-FA-25.a` *Kein Ende*) bleibt ein Aufbaufehler und erzeugt keine Session, auch mit `--record-empty-sessions`; `OpenSession` liefert dann die Kennung 0.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Upstream-Adapter (Erkennen der Fehlerantwort `FATAL`/`PANIC`, nicht auf das Verbindungsende warten, Ende im Leerlauf mit nachfolgender Anfrage, Lesestand bei Sendefehler, Einstufung `PGR-E4003`/`PGR-E6001` im Adapter) und Abnahmeszenario 18 — `slice-v1-abschluss-verbindungsende-upstream` (dort §1 und DoD): andere Schicht (Upstream-Adapter); mit ihm läge dieser Plan über zwei Schichten. Dieser Slice behandelt, was der Port liefert.
- SQLite-Aufzeichnung eines echten Endes — `slice-v1-abschluss-verbindungsende-upstream` (E2E); der Schreiber liegt in `slice-v1-abschluss-verbindungsende-sqlite`.
- Die Wiedergabe im Replay — `slice-v1-abschluss-verbindungsende-replay`.
- Ein Netzwerkabbruch ohne Fehlerantwort, ein Ende im Verbindungsaufbau und `FATAL` im Aufbau (etwa eine abgelehnte Anmeldung) als Ende — Out-of-Scope von `LH-FA-25`; sie bleiben wie bisher (`PGR-E4003`, `PGR-E4002`, `PGR-E4005`; `slice-v1-abschluss-anmeldung`, Schritt 30, bleibt unberührt).
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus): Der Record-Service übernimmt eine einfache Interaktion, deren Antworten mit `error_response` `FATAL`/`PANIC` und `connection_closed` enden, als letzte Interaktion der Session und gibt dem Aufrufer die Antworten bis einschließlich der Fehlerantwort zurück; die Interaktion ist abgeschlossen, `CloseSession` meldet für das Ende der Verbindung weder `PGR-E4003` noch `PGR-E6001`; eine Nachricht nach dem Ende gehört nicht zur Aufzeichnung; mit `--record-timing` trägt die Interaktion `offset_ms`, mit `--record-empty-sessions` ändert das Ende nichts; `wirdVerworfen` und die Entscheidung in `CloseSession` bleiben die eine Stelle, die über das Schreiben entscheidet (Bindung aus `slice-v1-abschluss-sessions`; Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Bei einer Extended-Interaktion steht das Ende in der letzten Gruppe, die mit `Sync` oder mit `Flush` abgeschlossen ist; Interaktionen, die der Client nach dieser Gruppe in die Pipeline gestellt hat, gehören nicht zur Aufzeichnung; ein Ende in einer nicht abgeschlossenen Gruppe ist `PGR-E6001`, die Session wird nicht übernommen, auch nicht ihre vorherigen Interaktionen; `ERROR` mit anschließendem Abbruch (kein Ende) bleibt `PGR-E6001` (Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Der PGWire-Adapter schreibt die Antworten bis einschließlich der Fehlerantwort an den Client und schließt dann Client- und Upstream-Verbindung (Session-Ende über den Use Case), ohne `ReadyForQuery`, ohne Meldung, ohne den Exit-Code zu berühren; ein Schreibfehler davor bleibt `EndWriteFailed` (Test mit Fake-Recorder).
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
| `internal/hexagon/services/record.go` | update | Übernahme einer Interaktion mit Ende (`Query`, `AwaitServer`, `uebernehmen`), Session als letzte, Stelle des Endes in der Gruppe, `CloseSession` |
| `internal/hexagon/ports/driving/record.go`, `ports/driven/upstream.go` | update | Port-Verträge: wie `connection_closed` aus dem Upstream kommt und dem Adapter gemeldet wird (Name und Form wählt der Architect vor dem Code) |
| `internal/adapters/driving/pgwire/server.go` | update | Record-Pfad: `serverRichtung`, `richtungen.schreibe`, Schließen beider Verbindungen |
| `internal/hexagon/services`, `internal/adapters/driving/pgwire` (Tests) | update | Randformen aus §6; die Tests, die die alte Regel prüfen (`TestRecordExtendedFehlerantwortVorDemAbbruch`, `TestRecordExtendedSendefehlerNachFehlerantwort`, soweit sie `FATAL`/`PANIC` meinen), werden hier nachgezogen |
| `docs/user/abdeckung-unit.md` | update (`make abdeckung`) | Tabelle aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` und `slice-v1-abschluss-verbindungsende-replay` liegen in `done/` (Voraussetzung: Formregeln im Modell, Abbildung im Adapter); `slice-v1-abschluss-sessions` und `slice-v1-abschluss-anmeldung` liegen davor in der Reihenfolge, weil sie `record.go` und den Aufbau im Upstream-Adapter berühren. Dieser Slice ist Schritt 40 der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Prüfung des Architect vor dem Code schätzt den Diff mit Tests auf über 1000 Zeilen (Schätzung des Planners 700 bis 1300, an der Grenze) — Schnitt vor dem Anspruch: einfache Interaktion mit PGWire-Zustellung | Extended-Interaktion (DoD-Punkt 2) im Kern, Zuschnitt: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`).
- `in-progress` → `next`: Der Use Case braucht das Ende vom Adapter statt aus den Antworten des Upstream-Ports, oder der Upstream-Adapter muss sich ändern — dritte Schicht, zurück zur Zerlegung.
- `in-progress` → `open`: `slice-v1-abschluss-sessions` oder `slice-v1-abschluss-anmeldung` ändert die Einstufung des Verbindungsendes in `record.go` so, dass dieser Plan nicht mehr trägt — Planner prüft neu.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect; die Regel „aufgezeichnet wird genau: der Server beendet in laufender Interaktion mit `FATAL`/`PANIC`“ bestätigte der Nutzer am 2026-10-12). Ort: `LH-FA-25.a` *Aufzeichnen*, *Form*, *Stelle*, *Kein Ende*, `LH-FA-02.b`, `LH-FA-12.a`, `LH-FA-18.a`.

- *Ende bei `PANIC` und `FATAL`:* beide wie `FATAL` — `LH-FA-25.a`.
- *`ERROR`, dann `FATAL` in derselben Interaktion:* die Interaktion trägt beide Fehlerantworten in der Reihenfolge des Upstreams und das Ende; die letzte Fehlerantwort ist die des Endes — `LH-FA-25.a` *Aufzeichnen* (Antworten in Reihenfolge des Upstreams bis einschließlich der Fehlerantwort).
- *`ERROR` und dann Abbruch bleibt `PGR-E6001`:* Das Ende verlangt `FATAL`/`PANIC` — `LH-FA-25.a` *Kein Ende*, `LH-FA-02.b`.
- *Fehlerantwort gelesen oder nicht (Lesestand):* maßgeblich ist, was der Recorder gelesen hat; eine nicht gelesene zählt nicht — `LH-FA-25.a` *Aufzeichnen*, `LH-FA-02.b` *Gelesen*. Die Einstufung im Adapter liegt in `slice-v1-abschluss-verbindungsende-upstream`; hier die Reaktion des Use Case auf das, was der Port liefert.
- *Ende in einer `Sync`-Gruppe, in einer `Flush`-Gruppe, in einer unvollständigen Gruppe:* `Sync` und `Flush` schließen die Gruppe ab, die unvollständige ist `PGR-E6001` — `LH-FA-25.a` *Form*, *Stelle*.
- *Nachrichten des Clients nach der Gruppe mit dem Ende, schon in der Pipeline stehende Interaktionen:* gehören nicht zur Aufzeichnung — `LH-FA-25.a` *Form*. Offen für den Architect: bestätigen, dass dies auch für Interaktionen gilt, die schon vollständig an den Upstream gegangen sind.
- *Ende mit `--record-timing` und `--record-empty-sessions`:* die Interaktion trägt `offset_ms` wie jede; das Ende ändert die Regeln der leeren Session nicht (die Session hat die Interaktion mit Ende) — `LH-FA-21.a`, `LH-FA-12.a`.
- *Schreibfehler an den Client beim Zustellen der Antworten des Endes:* `EndWriteFailed` wie bei jeder Interaktion; die Interaktion zählt als nicht zugestellt — `LH-FA-02.b`, `LH-FA-13.a`. Offen für den Architect: bestätigen, ob die Session mit dem Ende dann übernommen wird.
- *Schreiben der Session, Zwangsende im Aufbau und Aufbaufehler* (Bindung aus `slice-v1-abschluss-sessions`, Nahtstelle 3 und 4): `CloseSession` entscheidet allein über das Schreiben und die Kennung; das Zwangsende im Aufbau und eine nicht unterstützte Interaktion verwirft die Session weiter, auch mit `--record-empty-sessions`; eine Session mit Ende ist nie leer. Ein Ende des Servers vor dem ersten `ReadyForQuery` erzeugt keine Session (Kennung 0, `CloseSession` mit 0 ist wirkungslos). Prüfung: `TestRecordZwangsendeLeereSession`, `TestRecordLeereSessionNichtUnterstuetztVerworfen`, `TestRecordAufbauFehlerKeineSession` und am Binary `TestE2ERecordKeineVerbindung` bleiben grün, dazu ein Fall `EndForced` nach einer Interaktion mit Ende (geschrieben, nicht „im Aufbau“) — `LH-FA-12.a` *Angenommene Verbindung*, `LH-FA-25.a` *Kein Ende*.
- *Aufrufe der Client-Richtung einer Session nacheinander:* Bindung aus `slice-v1-abschluss-zeitangaben-aufzeichnen` (`LH-FA-21.a` *Messpunkt*), unverändert — hier genannt.
- *Exit-Code:* ein Ende ist weder Meldung noch Fehler und ändert den Exit-Code nicht — `LH-FA-25.a` *Aufzeichnen*.

**Risiken:**

- Der Diff liegt an der Grenze einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×; `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`, 2×): die Rückführung in §4 benennt den Schnitt — **Ausgang:** offen bis Closure.
- `slice-v1-abschluss-sessions` (`CloseSession`, Nebenläufigkeit) und `slice-v1-abschluss-anmeldung` (`OpenSession`, Port) ändern `record.go` und die Ports vor diesem Slice; der Architect liest ihren Stand vor dem Code — **Ausgang:** offen bis Closure.
- Bis `slice-v1-abschluss-verbindungsende-upstream` liefert kein echter Upstream ein Ende; die Tests dieses Slice stützen sich auf Fakes, die nur abbilden, was die Spezifikation sagt — **Ausgang:** offen bis Closure (erwartet: eingetreten, `slice-v1-abschluss-verbindungsende-upstream`, dessen E2E belegt die Entsprechung).

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×), `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×), `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×) und `BEO-REPO/record-fehlerantwort-im-aufbau-ungeregelt` (1×; die Fehlerantwort im Aufbau bleibt ausgenommen, §1) — der Plan nennt Schätzung und Schnitt vorab (§4). Tritt die Rückführung ein oder misst die Prüfung des Architect über der Grenze, wird mindestens einer der drei Einträge zum dritten Mal belegt und ist eine Lücke. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13), einschließlich der Sendungen aus `slice-v1-abschluss-verbindungsende-format`, `-replay` und `slice-v1-abschluss-sessions` (sie fügen keinen Liefer-Punkt hinzu: sie gehören zu DoD-Punkt 1 und 3): Liefer-Punkte: 3 (einfache Interaktion im Use Case, Extended-Interaktion im Use Case, Zustellung im PGWire-Adapter). Schichten: zwei (Kern mit Service und Ports, `ARC-002`, `ARC-003`, `ARC-004`; PGWire-Adapter, `ARC-006`). Schätzung des Planners: 700 bis 1300 Zeilen mit Tests.

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-sessions` (Nahtstelle 3 und 4, Planner, 2026-10-12, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 1 (Record-Service) und fügt keinen Liefer-Punkt und keine Schicht hinzu: eine Bindung an vorhandenen Code (`CloseSession`, `wirdVerworfen`, `OpenSession`) und ein Testfall (Zwangsende nach einer Interaktion mit Ende). Liefer-Punkte weiter 3, Schichten weiter zwei (Kern, PGWire-Adapter); die Schätzung bleibt bei 700 bis 1300 Zeilen, der Fall kommt mit einigen Dutzend Zeilen hinzu; der Slice liegt weiter an der Grenze aus §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
