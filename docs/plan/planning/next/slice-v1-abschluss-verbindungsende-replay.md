# Slice slice-v1-abschluss-verbindungsende-replay: Verbindungsende im Replay

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-03.a` · `LH-FA-03.b` · `LH-FA-09.a` · `LH-FA-12.a` · `LH-FA-18.a` · `SPEC-002` · `ARC-002` · `ARC-006`

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

**Ziel:** `replay` gibt ein aufgezeichnetes Ende des Servers wieder: Nach den aufgezeichneten Server-Nachrichten der Interaktion, deren letzte die `error_response` des Endes ist, schließt der Recorder die Verbindung zum Client, ohne `ReadyForQuery`, ohne Meldung und ohne Fehler; die Session ist damit beendet und die Interaktion verbraucht. Bei einer Extended-Interaktion geschieht das nach den Server-Nachrichten der Gruppe, in der das Ende liegt, die Gruppen davor laufen wie sonst. Eine aufgezeichnete Lebendprüfung mit Ende ist eine Interaktion wie jede andere (`LH-FA-25.a` *Replay*, `LH-FA-09.a` *Lebendprüfung mit Ende*, `LH-FA-03.a` Schritt 5, `LH-FA-03.b` *Verbraucht*). Nachweis am Binary mit handgeschriebenen Aufzeichnungen, ohne PostgreSQL (Abnahmeszenario 18, Replay-Teil).

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): die Wiedergabe einer Aufzeichnung mit Ende. Das Laden liefert jener Slice; dieser setzt ihn voraus (Fixtures: handgeschriebene Aufzeichnungen in YAML, `test/integration/testdata`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Record-Pfad des PGWire-Adapters (Antworten bis zur Fehlerantwort zustellen, beide Verbindungen schließen) und der Use Case `RecordService` — `slice-v1-abschluss-verbindungsende-record` (dort §1 und DoD): eigener Code-Pfad im Adapter (`richtungen`, nicht der Replay-Pfad) und eine andere Schicht des Kerns; mit ihm läge dieser Plan über der Grenze einer Review-Sitzung. Die Abbildung von `connection_closed` auf „keine PGWire-Nachricht“ entsteht hier und wird dort benutzt.
- Das Einlesen einer SQLite-Datei mit Ende als Fixture — `slice-v1-abschluss-verbindungsende-sqlite` (liegt davor, Schritt 36): Der E2E-Test dieses Slice benutzt YAML; die Lesbarkeit der SQLite-Datei belegt jener Slice.
- Einspielen und Vergleich — `-einspielen`, `-vergleich`: `play` benutzt den Replay-Service nicht.
- Eine Ende-Antwort ohne `error_response` im Replay — Out-of-Scope von `LH-FA-25` („Wiedergabe eines Verbindungsendes ohne Fehlerantwort“); das Format lässt sie nicht zu (`PGR-E3003`).
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Der Replay-Service liefert die aufgezeichneten Antworten einer einfachen Interaktion mit Ende und meldet danach, dass die Verbindung endet; bei einer Extended-Interaktion liefert er die Server-Nachrichten der Gruppe mit dem Ende, in einer `Sync`-Gruppe wie in einer `Flush`-Gruppe, nach den Gruppen davor; die Interaktion ist mit dem Senden verbraucht; die Session ist beendet. Eine aufgezeichnete Lebendprüfung mit Ende wird nicht übersprungen, zählt für die Zuordnung (`first-request` und `connection`) und für die nicht verbrauchten Interaktionen, und eine eingehende Lebendprüfung mit demselben SQL-Text ist diese Anfrage (Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Der PGWire-Adapter bildet `connection_closed` nicht auf eine PGWire-Nachricht ab; er sendet die Nachrichten davor und schließt die Verbindung zum Client, ohne `ReadyForQuery`, ohne Meldung, ohne Fehler und ohne den Exit-Code zu berühren; ein Schreibfehler davor bleibt wie bisher (Test mit Fake-Replayer).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus): Am Binary: Ein Client erhält für die Anfrage mit Ende die aufgezeichneten Antworten und danach das Ende der Verbindung; die nächste Verbindung erhält die nächste Session; hat der Client vor der Ende-Interaktion `Terminate` gesendet, ist sie nicht verbraucht (`PGR-W2001`, mit `--fail-on-unconsumed` `PGR-E5002`) (Integrationstest mit handgeschriebener Aufzeichnung).
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
| `internal/hexagon/services/replay.go`, `lebendpruefung.go` | update | Ende liefern und die Verbindung als beendet melden; Verbrauch; Lebendprüfung mit Ende |
| `internal/hexagon/ports/driving/replay.go` | update | Port-Vertrag: wie der Use Case das Ende meldet (Name und Form wählt der Architect vor dem Code) |
| `internal/adapters/driving/pgwire/server.go` | update | `connection_closed` ist keine PGWire-Nachricht (`toMessage`, `send`); Replay-Pfad schließt die Verbindung (`replayZustellen`) |
| `internal/hexagon/services`, `internal/adapters/driving/pgwire` (Tests) | update | Randformen aus §6; Fake-Replayer im Adapter |
| `test/integration` (`replay_ende_e2e_test.go`, `testdata`) | neu | Binary mit handgeschriebener Aufzeichnung, ohne PostgreSQL |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` liegt in `done/` (Voraussetzung; das Laden einer Aufzeichnung mit Ende); dieser Slice ist Schritt 37 der Reihenfolge in §5 von welle-v1-abschluss und steht nach `slice-v1-abschluss-verbindungsende-sqlite`, damit seine Fixtures beide Formate wählen können. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Prüfung des Architect vor dem Code schätzt den Diff mit Tests auf über 1000 Zeilen (Schätzung des Planners 600 bis 1000) — Schnitt vor dem Anspruch: Kern des Replay | Abbildung im Adapter und E2E, Zuschnitt: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`).
- `in-progress` → `next`: Die Meldung des Endes verlangt eine Änderung des Ports, die eine dritte Schicht trägt (etwa ein Eingriff in `SessionEnd` oder in den Recorder-Port) — zurück zur Zerlegung; der Port `Replayer` liegt in dieser Zählung beim Kern.
- `in-progress` → `open`: Eine Entscheidung des Nutzers ändert die Wiedergabe (etwa „Replay sendet `ReadyForQuery` statt zu schließen“) — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect; „Lebendprüfung mit Ende wird im Replay als normale Interaktion wiedergegeben“ bestätigte der Nutzer am 2026-10-12). Ort: `LH-FA-25.a` *Replay*, `LH-FA-03.a`, `LH-FA-03.b`, `LH-FA-09.a`, `LH-FA-12.a`.

- *Ende bei `PANIC` und `FATAL`:* beide werden wie aufgezeichnet gesendet, der Schweregrad ändert die Wiedergabe nicht — `LH-FA-25.a` *Replay*.
- *Ende in der letzten Gruppe einer Extended-Interaktion:* Gruppen davor laufen wie sonst; die Verbindung schließt nach den Server-Nachrichten der Gruppe mit dem Ende, in einer `Sync`-Gruppe wie in einer `Flush`-Gruppe — `LH-FA-25.a` *Replay*, *Form*, `LH-FA-03.a` Schritt 5.
- *`Terminate` des Clients vor der Ende-Interaktion:* die Interaktion ist nicht verbraucht, `PGR-W2001`; mit `--fail-on-unconsumed` `PGR-E5002` — `LH-FA-03.b` *Zeitpunkte* 1.
- *`--fail-on-unconsumed` bei verbrauchter Interaktion mit Ende:* verbraucht ist sie, sobald ihre Antworten gesendet sind; die Verbindung schließt danach ohne Meldung — `LH-FA-03.b` *Verbraucht*.
- *Lebendprüfung mit Ende:* wird nicht übersprungen, zählt für Zuordnung und nicht verbrauchte Interaktionen, Antworten werden verwendet; eine eingehende Lebendprüfung mit gleichem SQL-Text ist diese Anfrage, mit anderem Text gilt Schritt 1 — `LH-FA-09.a` *Lebendprüfung mit Ende*.
- *Session-Zuordnung:* eine Session mit Ende wird wie jede zugeordnet (`first-request`, `connection`); das Ende ist die letzte Interaktion ihrer Session — `LH-FA-12.a`.
- *Ende mit `--record-timing`:* `offset_ms` ist für das Replay keine Voraussetzung und ändert die Wiedergabe nicht — `LH-FA-21.a`, `SPEC-002`.
- *Schreibfehler an den Client beim Senden der Antworten des Endes:* bleibt wie bei jeder Interaktion (Interaktion nicht verbraucht, `PGR-W2001`); kein neuer Fall — `LH-FA-03.b` *Verbraucht*. Offen für den Architect: bestätigen, dass die Spezifikation hier nichts Eigenes verlangt.

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×): die Rückführung in §4 benennt den Schnitt — **Ausgang:** offen bis Closure.
- Der Schnitt zu `slice-v1-abschluss-verbindungsende-record` lässt die Abbildung von `connection_closed` im Adapter in diesem Slice und die Benutzung im Record-Pfad dort: ändert der Architect die Abbildung (etwa `send` ruft die Schließung selbst), trifft das beide — **Ausgang:** offen bis Closure.
- Der E2E-Test stützt sich auf handgeschriebene Aufzeichnungen, bevor `record` ein Ende erzeugen kann: Sie bilden nur ab, was die Spezifikation sagt, nicht was PostgreSQL liefert; die Entsprechung belegt der E2E-Test in `slice-v1-abschluss-verbindungsende-upstream` — **Ausgang:** offen bis Closure (erwartet: weiter offen oder entfallen mit jenem Beleg).

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×) — der Plan nennt seine Schätzung mit Tests und den Schnitt vorab (§4); tritt die Rückführung ein, ist es das dritte Auftreten von `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` und der Eintrag eine Lücke. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13): Liefer-Punkte: 3 (Replay-Service mit Randformen, Abbildung im PGWire-Adapter, E2E am Binary). Schichten: zwei (Kern mit Service und Port, `ARC-002`, `ARC-003`; PGWire-Adapter, `ARC-006`). Schätzung des Planners: 600 bis 1000 Zeilen mit Tests.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
