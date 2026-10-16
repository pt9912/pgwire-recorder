# Slice slice-v1-abschluss-verbindungsende-replay-extended: Verbindungsende im Replay (Extended-Interaktion, Gruppe mit Sync)

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

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `replay` gibt das aufgezeichnete Ende einer Extended-Interaktion wieder, deren Gruppe mit dem Ende mit `Sync` abschließt: Die Gruppen vor der Gruppe mit dem Ende laufen wie sonst; für die Gruppe, in der das Ende liegt, liefert der Replay-Service die Server-Nachrichten bis zur `error_response` des Endes, danach schließt der Recorder die Verbindung zum Client, ohne `ReadyForQuery`, ohne Meldung und ohne Fehler; die Interaktion ist mit dem Senden verbraucht, die Session beendet (`LH-FA-25.a` *Replay*, *Form*, `LH-FA-03.a` Schritt 5, `LH-FA-03.b` *Verbraucht*, `LH-FA-18.a`). Nachweis am Binary mit handgeschriebenen Aufzeichnungen, ohne PostgreSQL (Abnahmeszenario 18, Replay-Teil, für die Extended-Interaktion mit `Sync`).

**Geschnitten vor dem Anspruch in `-replay-extended` (dieser Slice, Schritt 37b) und `slice-v1-abschluss-verbindungsende-replay-extended-flush` (Schritt 37c)** (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap): Die Schätzung des Planners mit dem gemessenen Faktor 10,2 lag ungeteilt bei 376 bis 844, Mitte 610, 87 % der Marke 700, über der Grenze von 70 %. Schnitt nach Lieferwert: Dieser Slice liefert das Ende in der Gruppe mit `Sync` (der Pfad der Gruppen, die Verbrauchs-, Abweichungs- und Zuordnungsfälle, Adapter und E2E); die Gruppe mit dem Ende, die mit `Flush` abschließt, liefert `slice-v1-abschluss-verbindungsende-replay-extended-flush`, der diesen Pfad voraussetzt.

**Geschnitten aus `slice-v1-abschluss-verbindungsende-replay`** (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap): Jener Slice (Schritt 37a) liefert das Ende der einfachen Interaktion, die Meldung des Endes im Port `Replayer`, die Abbildung von `connection_closed` im PGWire-Adapter und die Zuordnung der Session; dieser Slice (Schritt 37b) setzt ihn voraus und wendet sie auf den Pfad `ClientMessage` an.

**Übernommen aus `slice-v1-abschluss-verbindungsende-replay`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-16): die Extended-Interaktion mit Ende; davon die gleiche Gruppe mit `Flush` weiter an `slice-v1-abschluss-verbindungsende-replay-extended-flush`.

**Gebunden durch `slice-v1-abschluss-sessions`** (Nahtstelle 1 und 2, über `slice-v1-abschluss-verbindungsende-replay` §1 und §6): `ClientMessage` ordnet mit denselben Schritten zu wie `Query` (`first-request`, `connection`, `Unassigned`); die Entscheidung, was eine Interaktion ist, bleibt der eine Schritt `ohneLebendpruefungen`; dieser Slice baut keine zweite Entscheidung daneben, und das Schließen nach dem Ende ändert die Zuordnung bei `connection` nicht.

**Gebunden durch `slice-v1-abschluss-verbindungsende-replay`** (Sendung des Planners, 2026-10-16): Dieser Slice benutzt die Meldung des Endes im Port und die Abbildung im Adapter, die jener liefert (Name und Form wählt der Architect bei jenem Slice); er ändert beide nicht. Zeigt sich, dass der Pfad der Gruppen eine andere Form des Ports verlangt, hält der Implementer an und gibt es dem Architect zurück, statt den Port zu ändern.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Gruppe mit dem Ende, die mit `Flush` abschließt (Server-Nachrichten bis zur `error_response` bei `Flush`, Abweichung vor dem `Flush`, Adapter-Fall und E2E-Fall dazu) — `slice-v1-abschluss-verbindungsende-replay-extended-flush` (Klasse 1, setzt diesen Slice voraus; dessen §1 und DoD nehmen es an): sonst läge dieser Plan bei etwa 87 % der Marke, mit dem oberen Ende darüber. Gruppen **vor** dem Ende, die mit `Flush` abschließen, laufen wie sonst und gehören zu den Gruppen davor.
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

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Der Replay-Service liefert bei einer Extended-Interaktion mit Ende die Server-Nachrichten der Gruppe mit dem Ende, die mit `Sync` abschließt (ohne `ReadyForQuery`), nach den Gruppen davor (auch solchen mit `Flush`), und meldet danach, dass die Verbindung endet; die Interaktion ist mit dem Senden verbraucht, die Session beendet; eine Abweichung vor dem abschließenden `Sync` ist `PGR-E5001` ohne Ende; `Terminate` oder Abbruch vor der Gruppe mit dem Ende lässt die Interaktion unverbraucht (`PGR-W2001`, mit `--fail-on-unconsumed` `PGR-E5002`); eine Zuordnung bei `first-request` und `connection` verhält sich wie bei der einfachen Interaktion (Bindung aus `slice-v1-abschluss-sessions`; Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Am Binary: Ein Client mit Extended Query erhält für die Gruppe mit dem Ende, die mit `Sync` abschließt, die aufgezeichneten Antworten und danach das Ende der Verbindung; die Gruppen davor sind unverändert; die nächste Verbindung erhält die nächste Session (Integrationstest mit handgeschriebener Aufzeichnung, eine Datei je Fall; Test des PGWire-Adapters mit Fake-Replayer für das Senden der Gruppe vor dem Schließen, `sync`, zugesagt). Die Zeilen je E2E-Fall mit Testdaten stehen in §7 (Messung).
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
| `internal/hexagon/services/replay.go` | update (nur falls nötig) | Pfad `ClientMessage`: Gruppe mit Ende (`Sync`) liefern, Verbrauch, Verbindung als beendet melden |
| `internal/adapters/driving/pgwire/server.go` | update (nur falls nötig) | Senden der Gruppe vor dem Schließen im Pfad `replayAntwort`; die Abbildung selbst liefert `-replay` |
| `internal/hexagon/services`, `internal/adapters/driving/pgwire` (Tests) | update | Randformen aus §6; Fake-Replayer im Adapter |
| `test/integration` (`replay_ende_extended_e2e_test.go`, `testdata`, eine Datei je Fall) | neu | Binary mit handgeschriebener Extended-Aufzeichnung, Gruppe mit `Sync`, ohne PostgreSQL |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-replay` liegt in `done/` (Voraussetzung: Meldung des Endes im Port, Abbildung im Adapter, einfache Interaktion); dieser Slice ist Schritt 37b der Reihenfolge in §5 von welle-v1-abschluss, geschnitten am 2026-10-16. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis:** `1896db3ea5103e462a7aec895ea877f43b7b6338`, Stand des Hauptzweigs beim Anspruch; die Commits danach ändern nur Pläne und die Verzeichnis-Position, keinen Pfad der Metrik); **über 700** geht der Slice zurück zur Zerlegung. **Meldeschwelle 550:** Erreicht die Summe 550, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik, gerechnet mit dem Faktor 10 (gemessen 10,2 in `-replay`) und beim Schnitt vor dem Anspruch: 260 bis 560, Mitte 410 und damit 59 % der Marke (§8), unter der Grenze von 70 %, das obere Ende bei 80 %; Zuschnitt eines weiteren Schnitts: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`).
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

- *Stelle der Nachricht `connection_closed`, Gruppe mit dem Ende* (`AGENTS.md` §3.12, Stelle einer Nachricht): (a) *Letzte Server-Nachricht der Gruppe:* die einzige zulässige Stelle; das Format lehnt jede andere beim Laden ab (`PGR-E3003`, `Session.ValidateEnd`). (b) *Gruppe mit dem Ende ist die letzte Gruppe der Interaktion:* immer; Gruppen nach dem Ende gibt es nicht. (c) *Gruppe mit dem Ende ist nicht die letzte der Session:* gibt es nicht, das Ende ist das letzte Ereignis seiner Session. (d) *In einer `Sync`-Gruppe, in einer `Flush`-Gruppe:* die `Sync`-Gruppe ist Gegenstand dieses Slice: gesendet wie aufgezeichnet, danach schließt die Verbindung ohne `ReadyForQuery`; die `Flush`-Gruppe entscheidet `LH-FA-25.a` ebenso, geliefert und geprüft wird sie in `slice-v1-abschluss-verbindungsende-replay-extended-flush` (`LH-FA-25.a` *Replay*, *Form*). (e) *Mitte der Gruppe:* Server-Nachrichten hinter dem Ende gibt es nicht; vor dem Ende liegende (etwa `ParseComplete`, `BindComplete`) werden gesendet. Entschieden, `LH-FA-25.a`: (a) bis (e) tragen, besonders (d): die Flush-Gruppe mit Ende wie die Sync-Gruppe; dieser Slice erfindet für sie nichts.
- *Client-Nachrichten der Gruppe mit dem Ende:* Die Gruppe wird erst beim abschließenden `Sync` (hier) oder `Flush` (`-flush`) geliefert (`LH-FA-18.a`); das Ende kommt mit dieser Lieferung, nicht bei einer früheren Nachricht der Gruppe. Entschieden, `LH-FA-25.a`: Eine Abweichung vor dem abschließenden `Sync` oder `Flush` ist `PGR-E5001`, es gibt kein Ende, die Interaktion bleibt unverbraucht (Reihenfolge der Gründe: Abweichung vor Ende).
- *`Terminate` oder Abbruch vor der Gruppe mit dem Ende:* die Interaktion ist nicht verbraucht, `PGR-W2001`; mit `--fail-on-unconsumed` `PGR-E5002` — `LH-FA-03.b` *Zeitpunkte* 1.
- *Gruppen davor, nach dem Ende:* Gruppen davor laufen wie ohne Ende (Verbrauch der Gruppe mit dem Senden); nach dem Ende liest der Adapter nichts mehr, eine Folge-Nachricht des Clients erreicht den Service nicht — Bindung an `slice-v1-abschluss-verbindungsende-replay`.
- *Herunterfahren während der Gruppen vor dem Ende* (`Shutdown`, `Forced`): unverändert, `LH-FA-13.a`; ein Ende ändert das Zwangsende nicht, weil die Interaktion bis zum Senden der Gruppe mit dem Ende nicht verbraucht ist. Entschieden, `LH-FA-25.a`: Herunterfahren wie `LH-FA-13.a`.
- *Fixtures je Fall isoliert:* Jeder E2E-Fall und jeder Fall des Services baut seine Aufzeichnung selbst (eine Datei je E2E-Fall, kein gemeinsamer Zustand zwischen Fällen), damit eine Mutation nur den Fall färbt, den sie bricht (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Erfahrung `-replay`, Review F-1 bis F-7); die Zeilen je E2E-Fall mit Testdaten meldet der Implementer in §7 (Messung, Maß für die Schätzung von `-flush`).
- *Session-Zuordnung, Lebendprüfung, `--record-timing`, `PANIC` und `FATAL`, Schreibfehler:* wie in `slice-v1-abschluss-verbindungsende-replay` §6 entschieden; hier wird nichts neu entschieden.

**Risiken:**

- Zwischenstand: Nach diesem Slice gibt `replay` das Ende einer Extended-Interaktion in der Gruppe mit `Sync` wieder; die Gruppe mit dem Ende, die mit `Flush` abschließt, ist nicht zugesagt und nicht getestet (der Service reicht beide Gruppen gleich durch, der Test fehlt) — **Ausgang:** offen bis Closure (erwartet: *eingetreten*, Folge-Slice `slice-v1-abschluss-verbindungsende-replay-extended-flush`, nimmt an: dessen §1 und DoD nennen die `Flush`-Gruppe mit Ende).

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×; Tests der Gruppen sind teuer): vor dem Anspruch geschnitten (`Flush`-Gruppe in `-replay-extended-flush`); die Rückführung in §4 benennt Marke und Meldeschwelle — **Ausgang:** offen bis Closure.
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

Vermerk bei der Closure des Gebers (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Schätzung und Plan bleiben unverändert): `slice-v1-abschluss-verbindungsende-replay` maß 968 in der Metrik von §4 bei einer Schätzung des Architect von 500 bis 850 (97 % der Marke 1000); Faktor Tests (Unit und E2E) zu Produktcode rund 10,2 (777 zu 76), dazu Testdaten und Tabellen 115. Die Schätzung dieses Plans (Faktor mindestens 5) gerechnet mit 10,2 auf den Produktcode von 30 bis 70: 306 bis 714 Zeilen Tests, mit Produktcode, Testdaten und Tabellen (40 bis 60) **376 bis 844, Mitte 610 und damit 87 % der Marke 700**, über der Grenze von 70 % und über der Meldeschwelle 550. Gegenrechnung: Der Architect entschied, dass der Service durchreicht (Produktcode kleiner als 30 bis 70 möglich), aber die Gruppen tragen den größeren Teil der Tests, und E2E mit Testdaten kostete im Geber rund 68 Zeilen je Fall (273 für vier Fälle); hier mindestens zwei (`Sync`- und `Flush`-Gruppe). Der Adapter-Test des Gebers fährt den Lauf `sync` mit dem Fake schon (`TestReplayEndeSendetDavorUndSchliesst`, nicht zugesagt); dieser Slice sagt ihn zu. Ob vor dem Anspruch geschnitten wird, entscheidet der Planner beim Anspruch mit der Prüfung des Architect (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`); Liefer-Punkte weiter 2, Schichten weiter zwei.

Nachgezählt beim Schnitt vor dem Anspruch (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14): Ungeteilt, mit dem Faktor 10 (gemessen 10,2) auf Produktcode 30 bis 70: 376 bis 844, Mitte 610, 87 % der Marke 700, über der Grenze von 70 %. Geschnitten nach Lieferwert (Gruppe mit `Sync` | Gruppe mit `Flush`), nicht nach Schichten. **Dieser Teil (37b):** Produktcode 20 bis 45 (der Service reicht die Gruppe mit dem Ende durch, `server.go` nur falls nötig), Tests 130 bis 270 im Service (Gruppe mit `Sync` und Ende, Gruppen davor mit `Sync` und `Flush`, Abweichung vor dem `Sync`, `Terminate` und Abbruch davor, `--fail-on-unconsumed`, `first-request` und `connection`, Herunterfahren) plus 40 bis 80 im Adapter (`sync`, Fake-Replayer), E2E mit Testdaten 70 bis 140 (ein Fall nach dem Maß von rund 68 Zeilen je Fall in `-replay`; zweiter Fall `nächste Verbindung erhält die nächste Session`), Tabellen 15 bis 25: **260 bis 560, Mitte 410, 59 % der Marke 700**, oberes Ende 80 %; die Meldeschwelle 550 bleibt. Liefer-Punkte: 2 (Replay-Service mit Randformen; Binary: Adapter-Test und E2E). Schichten: zwei nach der Auffassung des Architect (Kern mit Service, `ARC-002`; PGWire-Adapter, `ARC-006`, vor allem Tests), Abdeckungstabellen in `docs/user` keine Schicht (die Nutzerentscheidung dazu, ob `AGENTS.md` §3.13 es festhält, steht aus); **konservativ gezählt mit `docs/user` als dritter Schicht: drei** — markiert, bis der Nutzer entscheidet. Die Gruppe mit `Flush` liegt in `slice-v1-abschluss-verbindungsende-replay-extended-flush` (Schritt 37c, setzt diesen Slice voraus; dessen §1 nennt die Bindung an diese Randformen).
