# Slice slice-v1-abschluss-verbindungsende-replay-extended: Verbindungsende im Replay (Extended-Interaktion)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-03.a` · `LH-FA-03.b` · `LH-FA-18.a` · `ARC-002` · `ARC-006`

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `replay` gibt das aufgezeichnete Ende einer Extended-Interaktion wieder: Die Gruppen vor der Gruppe mit dem Ende laufen wie sonst; für die Gruppe, in der das Ende liegt, liefert der Replay-Service die Server-Nachrichten bis zur `error_response` des Endes, in einer `Sync`-Gruppe wie in einer `Flush`-Gruppe, danach schließt der Recorder die Verbindung zum Client, ohne `ReadyForQuery`, ohne Meldung und ohne Fehler; die Interaktion ist mit dem Senden verbraucht, die Session beendet (`LH-FA-25.a` *Replay*, *Form*, `LH-FA-03.a` Schritt 5, `LH-FA-03.b` *Verbraucht*, `LH-FA-18.a`). Nachweis am Binary mit handgeschriebenen Aufzeichnungen, ohne PostgreSQL (Abnahmeszenario 18, Replay-Teil, für die Extended-Interaktion).

**Geschnitten aus `slice-v1-abschluss-verbindungsende-replay`** (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap): Jener Slice (Schritt 37a) liefert das Ende der einfachen Interaktion, die Meldung des Endes im Port `Replayer`, die Abbildung von `connection_closed` im PGWire-Adapter und die Zuordnung der Session; dieser Slice (Schritt 37b) setzt ihn voraus und wendet sie auf den Pfad `ClientMessage` an.

**Übernommen aus `slice-v1-abschluss-verbindungsende-replay`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-16): die Extended-Interaktion mit Ende.

**Gebunden durch `slice-v1-abschluss-sessions`** (Nahtstelle 1 und 2, über `slice-v1-abschluss-verbindungsende-replay` §1 und §6): `ClientMessage` ordnet mit denselben Schritten zu wie `Query` (`first-request`, `connection`, `Unassigned`); die Entscheidung, was eine Interaktion ist, bleibt der eine Schritt `ohneLebendpruefungen`; dieser Slice baut keine zweite Entscheidung daneben, und das Schließen nach dem Ende ändert die Zuordnung bei `connection` nicht.

**Gebunden durch `slice-v1-abschluss-verbindungsende-replay`** (Sendung des Planners, 2026-10-16): Dieser Slice benutzt die Meldung des Endes im Port und die Abbildung im Adapter, die jener liefert (Name und Form wählt der Architect bei jenem Slice); er ändert beide nicht. Zeigt sich, dass der Pfad der Gruppen eine andere Form des Ports verlangt, hält der Implementer an und gibt es dem Architect zurück, statt den Port zu ändern.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die einfache Interaktion mit Ende, Lebendprüfung mit Ende, Meldung im Port und Abbildung im Adapter — `slice-v1-abschluss-verbindungsende-replay` (Klasse 1, Voraussetzung dieses Slice): sonst läge dieser Plan bei etwa 100 % der Marke.
- Der Record-Pfad (Zustellen der Gruppe an den Client, Schließen beider Verbindungen) und der Use Case `RecordService` — `slice-v1-abschluss-verbindungsende-record`: anderer Code-Pfad im Adapter (`richtungen`) und andere Schicht des Kerns.
- Das Einspielen einer Extended-Aufzeichnung mit Ende und der Vergleich — `slice-v1-abschluss-verbindungsende-einspielen`, `-vergleich`: `play` benutzt den Replay-Service nicht.
- Eine Ende-Antwort ohne `error_response` im Replay — Out-of-Scope von `LH-FA-25`; das Format lässt sie nicht zu (`PGR-E3003`).
- Weitere Anmeldeverfahren und `NotificationResponse` im Replay — Entscheidung: nein, wie in `slice-v1-abschluss-verbindungsende-replay` §1 begründet; die Spezifikation trägt sie.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Der Replay-Service liefert bei einer Extended-Interaktion mit Ende die Server-Nachrichten der Gruppe mit dem Ende, in einer `Sync`-Gruppe wie in einer `Flush`-Gruppe, nach den Gruppen davor, und meldet danach, dass die Verbindung endet; die Interaktion ist mit dem Senden verbraucht, die Session beendet; `Terminate` oder Abbruch vor der Gruppe mit dem Ende lässt die Interaktion unverbraucht (`PGR-W2001`, mit `--fail-on-unconsumed` `PGR-E5002`); eine Zuordnung bei `first-request` und `connection` verhält sich wie bei der einfachen Interaktion (Bindung aus `slice-v1-abschluss-sessions`; Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Am Binary: Ein Client mit Extended Query erhält für die Gruppe mit dem Ende die aufgezeichneten Antworten und danach das Ende der Verbindung, in einer `Sync`- wie in einer `Flush`-Gruppe; die Gruppen davor sind unverändert; die nächste Verbindung erhält die nächste Session (Integrationstest mit handgeschriebener Aufzeichnung; Test des PGWire-Adapters mit Fake-Replayer für das Senden der Gruppe vor dem Schließen).
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
| `internal/hexagon/services/replay.go` | update | Pfad `ClientMessage`: Gruppe mit Ende liefern, Verbrauch, Verbindung als beendet melden |
| `internal/adapters/driving/pgwire/server.go` | update (nur falls nötig) | Senden der Gruppe vor dem Schließen im Pfad `replayAntwort`; die Abbildung selbst liefert `-replay` |
| `internal/hexagon/services`, `internal/adapters/driving/pgwire` (Tests) | update | Randformen aus §6; Fake-Replayer im Adapter |
| `test/integration` (`replay_ende_extended_e2e_test.go`, `testdata`) | neu | Binary mit handgeschriebener Extended-Aufzeichnung, ohne PostgreSQL |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-replay` liegt in `done/` (Voraussetzung: Meldung des Endes im Port, Abbildung im Adapter, einfache Interaktion); dieser Slice ist Schritt 37b der Reihenfolge in §5 von welle-v1-abschluss, geschnitten am 2026-10-16. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis:** vom Planner beim Anspruch einzutragen); **über 700** geht der Slice zurück zur Zerlegung. **Meldeschwelle 550:** Erreicht die Summe 550, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik: 300 bis 600, Mitte 450 und damit 64 % der Marke (§8), unter der Grenze von 70 %; Zuschnitt eines weiteren Schnitts: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`).
- `in-progress` → `next`: Der Pfad der Gruppen verlangt eine Änderung des Ports `Replayer` oder der Abbildung im Adapter, die `slice-v1-abschluss-verbindungsende-replay` geliefert hat — zurück zur Zerlegung (Architect und Planner), nicht still im Code.
- `in-progress` → `open`: Eine Entscheidung des Nutzers ändert die Wiedergabe — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), vor dem Code in der Spezifikation entschieden oder vom Architect bei der Prüfung entschieden (die Spalte sagt, was schon steht). Ort: `LH-FA-25.a` *Replay*, `LH-FA-03.a` Schritt 5, `LH-FA-03.b`, `LH-FA-18.a`.

- *Stelle der Nachricht `connection_closed`, Gruppe mit dem Ende* (`AGENTS.md` §3.12, Stelle einer Nachricht): (a) *Letzte Server-Nachricht der Gruppe:* die einzige zulässige Stelle; das Format lehnt jede andere beim Laden ab (`PGR-E3003`, `Session.ValidateEnd`). (b) *Gruppe mit dem Ende ist die letzte Gruppe der Interaktion:* immer; Gruppen nach dem Ende gibt es nicht. (c) *Gruppe mit dem Ende ist nicht die letzte der Session:* gibt es nicht, das Ende ist das letzte Ereignis seiner Session. (d) *In einer `Sync`-Gruppe, in einer `Flush`-Gruppe:* beide wie aufgezeichnet gesendet, danach schließt die Verbindung ohne `ReadyForQuery` (`LH-FA-25.a` *Replay*, *Form*). (e) *Mitte der Gruppe:* Server-Nachrichten hinter dem Ende gibt es nicht; vor dem Ende liegende (etwa `ParseComplete`, `BindComplete`) werden gesendet. Entschieden, `LH-FA-25.a`: (a) bis (e) tragen, besonders (d): die Flush-Gruppe mit Ende wie die Sync-Gruppe.
- *Client-Nachrichten der Gruppe mit dem Ende:* Die Gruppe wird erst beim abschließenden `Flush` oder `Sync` geliefert (`LH-FA-18.a`); das Ende kommt mit dieser Lieferung, nicht bei einer früheren Nachricht der Gruppe. Entschieden, `LH-FA-25.a`: Eine Abweichung vor dem abschließenden `Sync` oder `Flush` ist `PGR-E5001`, es gibt kein Ende, die Interaktion bleibt unverbraucht (Reihenfolge der Gründe: Abweichung vor Ende).
- *`Terminate` oder Abbruch vor der Gruppe mit dem Ende:* die Interaktion ist nicht verbraucht, `PGR-W2001`; mit `--fail-on-unconsumed` `PGR-E5002` — `LH-FA-03.b` *Zeitpunkte* 1.
- *Gruppen davor, nach dem Ende:* Gruppen davor laufen wie ohne Ende (Verbrauch der Gruppe mit dem Senden); nach dem Ende liest der Adapter nichts mehr, eine Folge-Nachricht des Clients erreicht den Service nicht — Bindung an `slice-v1-abschluss-verbindungsende-replay`.
- *Herunterfahren während der Gruppen vor dem Ende* (`Shutdown`, `Forced`): unverändert, `LH-FA-13.a`; ein Ende ändert das Zwangsende nicht, weil die Interaktion bis zum Senden der Gruppe mit dem Ende nicht verbraucht ist. Entschieden, `LH-FA-25.a`: Herunterfahren wie `LH-FA-13.a`.
- *Session-Zuordnung, Lebendprüfung, `--record-timing`, `PANIC` und `FATAL`, Schreibfehler:* wie in `slice-v1-abschluss-verbindungsende-replay` §6 entschieden; hier wird nichts neu entschieden.

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×; Tests der Gruppen sind teuer): die Rückführung in §4 benennt Marke und Meldeschwelle — **Ausgang:** offen bis Closure.
- Der Port aus `slice-v1-abschluss-verbindungsende-replay` trägt den Pfad der Gruppen nicht (Meldung des Endes pro Antwortliste statt pro Interaktion): die zweite Rückführung in §4 — **Ausgang:** offen bis Closure.
- Handgeschriebene Extended-Aufzeichnungen bilden nur ab, was die Spezifikation sagt, nicht was PostgreSQL liefert; die Entsprechung belegt der E2E-Test in `slice-v1-abschluss-verbindungsende-upstream` — **Ausgang:** offen bis Closure (erwartet: weiter offen oder entfallen mit jenem Beleg).

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` — der Plan nennt seine Schätzung mit Tests und den Schnitt vorab (§4); tritt die Rückführung ein, ist es ein weiteres Auftreten von `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`. Keine weiteren Treffer. Der Zähler-Stand ist beim Anspruch aus dem Register zu lesen.

Nachgezählt beim Anlegen (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14): Liefer-Punkte: 2 (Replay-Service für die Extended-Interaktion mit Randformen, E2E und Adapter-Test). Schichten: zwei (Kern mit Service, `ARC-002`, `ARC-003`; PGWire-Adapter, `ARC-006`, vor allem Tests); die Abdeckungstabellen in `docs/user` zählen nicht als Schicht (Begründung in `slice-v1-abschluss-verbindungsende-replay` §8). Schätzung des Planners in der Metrik von §4: Produktcode 30 bis 70 (Pfad `ClientMessage`, Gruppe mit Ende), Tests mit dem Faktor von mindestens 5 (gemessene Datenpunkte 8 bis 11) 230 bis 480 (Gruppen mit `Sync` und `Flush`, Gruppen davor, `Terminate`, `--fail-on-unconsumed`, Zuordnung; Adapter-Fall; E2E), Testdaten und Tabellen 40 bis 60: **300 bis 600, Mitte 450, 64 % der Marke 700**, oberes Ende 86 %; die Meldeschwelle 550 trägt, bei Erreichen hält der Implementer an.

Nachgezählt beim Eintragen der Entscheidungen des Architect (Planner, 2026-10-16, `AGENTS.md` §3.13, Quelle `slice-v1-abschluss-verbindungsende-replay`): Die drei Sätze „Offen für den Architect“ sind in `LH-FA-25.a` entschieden; es kommt kein Code und kein Test hinzu. Liefer-Punkte weiter 2, Schichten weiter zwei (Kern, PGWire-Adapter), Schätzung 300 bis 600, Mitte 450, unverändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
