# Slice slice-v1-abschluss-verbindungsende-replay-extended-flush: Verbindungsende im Replay (Extended-Interaktion, Gruppe mit Flush)

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

**Ziel:** `replay` gibt das aufgezeichnete Ende einer Extended-Interaktion wieder, deren Gruppe mit dem Ende mit `Flush` abschließt: Der Replay-Service liefert die Server-Nachrichten dieser Gruppe bis zur `error_response` des Endes bei der abschließenden `Flush`-Nachricht, danach schließt der Recorder die Verbindung zum Client, ohne `ReadyForQuery`, ohne Meldung und ohne Fehler; die Interaktion ist mit dem Senden verbraucht, die Session beendet; eine Abweichung vor dem `Flush` ist `PGR-E5001` ohne Ende (`LH-FA-25.a` *Replay*, *Form*, `LH-FA-03.a` Schritt 5, `LH-FA-03.b` *Verbraucht*, `LH-FA-18.a`). Nachweis am Binary mit handgeschriebener Aufzeichnung, ohne PostgreSQL (Abnahmeszenario 18, Replay-Teil, für die Extended-Interaktion mit `Flush`).

**Geschnitten aus `slice-v1-abschluss-verbindungsende-replay-extended`** (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap): Jener Slice (Schritt 37b) liefert den Pfad `ClientMessage` für die Gruppe mit dem Ende, die mit `Sync` abschließt, samt Gruppen davor, Verbrauch, Abweichung, Zuordnung, Adapter-Test `sync` und E2E; dieser Slice (Schritt 37c) setzt ihn voraus und fügt die Gruppe mit `Flush` hinzu.

**Übernommen aus `slice-v1-abschluss-verbindungsende-replay-extended`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-16): die Gruppe mit dem Ende, die mit `Flush` abschließt.

**Gebunden durch `slice-v1-abschluss-verbindungsende-replay-extended`** (Sendung des Planners, 2026-10-16): Die Randformen aus dessen §6 gelten hier unverändert (Stelle von `connection_closed` in der Gruppe, Client-Nachrichten der Gruppe, `Terminate` und Abbruch, Gruppen davor, Herunterfahren, Fixtures je Fall isoliert); dieser Slice entscheidet keine davon neu. Er benutzt den Pfad, den jener Slice liefert; zeigt sich, dass die `Flush`-Gruppe eine Änderung dieses Pfads, des Ports `Replayer` oder der Abbildung im Adapter (aus `slice-v1-abschluss-verbindungsende-replay`) verlangt, hält der Implementer an und gibt es dem Architect zurück. Die Zuordnung bei `first-request` und `connection` bleibt die aus `slice-v1-abschluss-sessions` (Nahtstelle 1 und 2, über jene Slices); der eine Schritt `ohneLebendpruefungen` bleibt die Entscheidung, was eine Interaktion ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Gruppe mit dem Ende, die mit `Sync` abschließt, die Gruppen davor, Verbrauch, Zuordnung, Herunterfahren und der Pfad `ClientMessage` selbst — `slice-v1-abschluss-verbindungsende-replay-extended` (Klasse 1, Voraussetzung dieses Slice): sonst läge der Plan bei etwa 87 % der Marke 700 und darüber.
- Die einfache Interaktion mit Ende, Lebendprüfung mit Ende, Meldung im Port und Abbildung im Adapter — `slice-v1-abschluss-verbindungsende-replay` (Klasse 1, in `done/`): geliefert, dieser Slice ändert sie nicht.
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

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Der Replay-Service liefert bei einer Extended-Interaktion, deren Gruppe mit dem Ende mit `Flush` abschließt, die Server-Nachrichten dieser Gruppe bis zur `error_response` bei `Flush` (ohne `ReadyForQuery`) und meldet danach, dass die Verbindung endet; die Interaktion ist mit dem Senden verbraucht, die Session beendet; eine Abweichung vor dem abschließenden `Flush` ist `PGR-E5001` ohne Ende, die Interaktion bleibt unverbraucht; `Terminate` oder Abbruch vor dieser Gruppe lässt sie unverbraucht (`PGR-W2001`, mit `--fail-on-unconsumed` `PGR-E5002`); die Randformen aus §6 von `slice-v1-abschluss-verbindungsende-replay-extended` gelten (Test je Fall, Fixture je Fall).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Am Binary: Ein Client mit Extended Query erhält für die Gruppe mit dem Ende, die mit `Flush` abschließt, die aufgezeichneten Antworten und danach das Ende der Verbindung; die Gruppen davor sind unverändert (Integrationstest mit handgeschriebener Aufzeichnung, eine Datei je Fall; Test des PGWire-Adapters mit Fake-Replayer, Lauf `flush`, zugesagt). Die Zeilen je E2E-Fall mit Testdaten stehen in §7 (Messung).
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
| `internal/hexagon/services/replay.go` | update (nur falls nötig) | Pfad `ClientMessage` für die Gruppe mit Ende, die mit `Flush` abschließt; der Pfad selbst liefert `slice-v1-abschluss-verbindungsende-replay-extended` |
| `internal/adapters/driving/pgwire/server.go` | update (nur falls nötig) | Senden der Gruppe vor dem Schließen im Pfad `replayAntwort`, Lauf `flush` |
| `internal/hexagon/services`, `internal/adapters/driving/pgwire` (Tests) | update | Randformen aus §6 für die `Flush`-Gruppe; Fake-Replayer im Adapter |
| `test/integration` (`replay_ende_extended_e2e_test.go`, `testdata`, eine Datei je Fall) | update | Fall `Flush` in der Datei, die `slice-v1-abschluss-verbindungsende-replay-extended` anlegt; Binary mit handgeschriebener Extended-Aufzeichnung, ohne PostgreSQL |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-replay-extended` liegt in `done/` (Voraussetzung: Pfad `ClientMessage` für die Gruppe mit Ende, Verbrauch, Zuordnung, Datei des E2E-Tests); dieser Slice ist Schritt 37c der Reihenfolge in §5 von welle-v1-abschluss, geschnitten am 2026-10-16. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis:** vom Planner beim Anspruch einzutragen, danach nicht mehr geändert); **über 500** geht der Slice zurück zur Zerlegung. **Meldeschwelle 400:** Erreicht die Summe 400, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik beim Schnitt vor dem Anspruch (Faktor 10, gemessen 10,2): 180 bis 335, Mitte 260 und damit 52 % der Marke (§8), das obere Ende bei 67 %; beim Anspruch ersetzt sie das gemessene Maß je E2E-Fall aus §7 von `slice-v1-abschluss-verbindungsende-replay-extended`. Zuschnitt eines weiteren Schnitts: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`).
- `in-progress` → `next`: Die `Flush`-Gruppe verlangt eine Änderung des Pfads aus `slice-v1-abschluss-verbindungsende-replay-extended`, des Ports `Replayer` oder der Abbildung im Adapter — zurück zur Zerlegung (Architect und Planner), nicht still im Code.
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

**Randformen** (`AGENTS.md` §3.12), vor dem Code in der Spezifikation entschieden oder vom Architect bei der Prüfung entschieden. Ort: `LH-FA-25.a` *Replay*, `LH-FA-03.a` Schritt 5, `LH-FA-03.b`, `LH-FA-18.a`. Die Randformen aus §6 von `slice-v1-abschluss-verbindungsende-replay-extended` gelten für beide Gruppenarten; hier stehen die, die der `Flush` ändert.

- *Stelle der Nachricht `connection_closed`, Gruppe mit dem Ende, abgeschlossen mit `Flush`* (`AGENTS.md` §3.12, Stelle einer Nachricht): (a) *Letzte Server-Nachricht der Gruppe:* die einzige zulässige Stelle (`PGR-E3003`, `Session.ValidateEnd`). (b) *Letzte Gruppe der Interaktion und der Session:* immer. (c) *`Flush`-Gruppe wie `Sync`-Gruppe:* die Server-Nachrichten vor dem Ende (etwa `ParseComplete`, `BindComplete`) werden gesendet, `ReadyForQuery` gibt es nicht, danach schließt die Verbindung. Entschieden, `LH-FA-25.a` *Replay*, *Form* und *Stelle von `connection_closed` im Replay*: „in einer `Sync`- wie in einer `Flush`-Gruppe die letzte“.
- *Client-Nachrichten der Gruppe:* geliefert wird beim abschließenden `Flush` (`LH-FA-18.a`); eine Abweichung davor ist `PGR-E5001`, es gibt kein Ende, die Interaktion bleibt unverbraucht. Entschieden, `LH-FA-25.a` *Reihenfolge der Gründe im Replay*.
- *Ein `Sync` nach dem `Flush` mit Ende:* gehört nicht zur Aufzeichnung, das Ende ist das letzte Ereignis der Session (`LH-FA-25.a` *Form*); der Adapter liest nach dem Ende nichts mehr.
- *Fixtures je Fall isoliert und Zeilen je E2E-Fall:* wie in §6 von `slice-v1-abschluss-verbindungsende-replay-extended`.
- *Ergänzungen des Architect aus `slice-v1-abschluss-verbindungsende-replay-extended` §6* (Sendung des Architect, 2026-10-16, `AGENTS.md` §3.13; Ort `LH-FA-25.a`): Sie gelten für die `Flush`-Gruppe gespiegelt, je ein Fall und eine Mutation: Abweichung bei der abschließenden Nachricht (`Sync` statt `Flush`), eine einfache Anfrage und eine Lebendprüfung mitten in der Gruppe, `Terminate` oder Abbruch vor dem `Flush`, Ende-Gruppe als einzige Gruppe, Verbrauch beider Seiten. Der Service liest die Art der abschließenden Nachricht nicht anders; der Adapter-Test fährt den Lauf `flush`. Der E2E-Fall ist ein Fall in der Datei des Gebers.
- *Wortlaut von `LH-FA-18.a` gegenüber einer Interaktion ohne `Sync`* (Sendung aus der Verifikation V-5 von `slice-v1-abschluss-verbindungsende-replay-extended`, Planner, 2026-10-16, `AGENTS.md` §3.13): `LH-FA-18.a` nennt eine `Query` mitten in einer Extended-Interaktion „vor deren `Sync`“ als Abweichung (`PGR-E5001`), und `LH-FA-25.a` *Reihenfolge der Gründe im Replay* übernimmt „vor dem `Sync` der Interaktion“. Eine Interaktion mit Ende, deren Gruppe mit `Flush` schließt, hat kein `Sync`; der Code entscheidet den Fall unabhängig davon (`mitten`: `gruppe > 0 || nachricht > 0`), der Fall „Lebendprüfung mitten in der Gruppe“ oben bleibt Fall und Mutation dieses Slice. Ob die Spezifikation den Wortlaut für diese Interaktion enger oder weiter fasst, entscheidet der Architect **vor dem ersten Code-Commit** dieses Slice (Frage des Planners, nicht seine Entscheidung); ist sie beim Anspruch offen, nennt der Implementer sie dem Architect, bevor er den Fall schreibt (`AGENTS.md` §3.12).

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`): die Rückführung in §4 benennt Marke und Meldeschwelle — **Ausgang:** offen bis Closure.
- Der Pfad aus `slice-v1-abschluss-verbindungsende-replay-extended` trägt die `Flush`-Gruppe nicht (er unterscheidet die Gruppenart im Service oder im Adapter): die zweite Rückführung in §4 — **Ausgang:** offen bis Closure.
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

Nachgezählt beim Anlegen (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Schnitt von `slice-v1-abschluss-verbindungsende-replay-extended`): Liefer-Punkte: 2 (Replay-Service für die `Flush`-Gruppe mit Randformen; Binary: Adapter-Test und E2E). Schichten: zwei nach der Auffassung des Architect (Kern mit Service, `ARC-002`; PGWire-Adapter, `ARC-006`, vor allem Tests), Abdeckungstabellen in `docs/user` keine Schicht; **konservativ gezählt mit `docs/user` als dritter Schicht: drei** — markiert, bis der Nutzer entscheidet. Schätzung in der Metrik von §4 (Faktor 10, gemessen 10,2): Produktcode 0 bis 10 (der Pfad besteht, der Architect fand ihn beim Geber durchreichend; nur die Gruppenart ist neu), Tests im Service 50 bis 110, im Adapter 40 bis 70, E2E mit Testdaten 70 bis 110 (rund 68 Zeilen je Fall in `-replay`), Tabellen 10 bis 20: **180 bis 335, Mitte 260, 52 % der Marke 500**, oberes Ende 67 %; Meldeschwelle 400. Die Summe aus beiden Teilen (Mitte 410 und 260) liegt bei 670 gegen 610 ungeteilt; der Mehraufwand ist die zweite Fixture und der zweite Adapter-Lauf.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Eintragen der Ergänzungen des Architect (2026-10-16, `AGENTS.md` §3.13, Quelle `slice-v1-abschluss-verbindungsende-replay-extended`): Liefer-Punkte weiter 2, Schichten zwei (drei konservativ mit `docs/user/`, Nutzerentscheidung steht aus). Die Ergänzungen sind Fälle derselben Tests, keine neue Datei und kein Produktcode: Schätzung 170 bis 325 (Produktcode 0 bis 10 statt 10 bis 25), Mitte 245, 49 % der Marke 500, oberes Ende 65 %; Meldeschwelle 400 bleibt.

Nachgezählt beim Eintragen der Sendung V-5 (Planner, 2026-10-16, `AGENTS.md` §3.13, Quelle `slice-v1-abschluss-verbindungsende-replay-extended`): Sie ist eine Frage an den Architect zum Wortlaut der Spezifikation und verlangt keinen Code, keinen neuen Test und keine Datei; der Fall „Lebendprüfung mitten in der Gruppe“ stand in den Ergänzungen schon. Liefer-Punkte weiter 2, Schichten zwei (drei konservativ mit `docs/user/`, Nutzerentscheidung steht aus), Schätzung unverändert.

Vermerk bei der Closure des Gebers (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Schätzung und Plan bleiben unverändert): `slice-v1-abschluss-verbindungsende-replay-extended` maß 532 in der Metrik von §4 (Schätzung des Architect 360 bis 545, Produktcode 0, Service-Tests 322 für neun Tests, Adapter-Tests 87, E2E 61 und Fixture 36, Tabellen 26; Nacharbeit nach dem Review +30). Das Maß je E2E-Fall war 27 Zeilen Fall und 36 Fixture, rund 66 mit der Tabellenzeile; die Hilfen des Gebers (33) nutzt dieser Slice wieder. Gerechnet für diesen Slice, der den Pfad nur spiegelt: Service-Tests als weitere Zeilen derselben Tabellen und Fälle 60 bis 130, Adapter-Test Lauf `flush` 30 bis 60, E2E ein Fall 63 bis 70, Tabellen 6 bis 12, Nacharbeit nach dem Review bis 10 vom Hundert: **rund 160 bis 300, Mitte 230, 46 % der Marke 500**, oberes Ende 60 %; das liegt in der Spanne 170 bis 325 oben und in deren unterer Hälfte. Marke 500 und Meldeschwelle 400 bleiben.

