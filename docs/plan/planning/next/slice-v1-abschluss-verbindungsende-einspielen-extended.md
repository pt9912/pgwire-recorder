# Slice slice-v1-abschluss-verbindungsende-einspielen-extended: Verbindungsende beim Einspielen ohne Vergleich (Extended-Interaktion und die Grenzen am Binary)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-20.a` · `LH-FA-18.a` · `SPEC-034` · `ARC-002` · `ARC-004` · `ARC-007`

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

**Ziel:** `play` ohne `--compare-responses` behandelt eine Fehlerantwort `FATAL` oder `PANIC` des Servers auch in einer **Extended-Interaktion** (`LH-FA-20.a` *Fehlerregeln*, *Gruppen*, *Ende einer Session*; `LH-FA-25.a` *Form*): Trägt die aufgezeichnete Interaktion ein Ende (in der letzten Gruppe, die mit `Sync` oder `Flush` abschließt) und ist `--allow-recorded-errors` gesetzt, ist das Ende erwartet — keine Meldung, `play` sendet kein `Terminate`, sondern schließt die Verbindung, die Session endet, die nächste läuft, auch ohne `--continue-on-error`; sonst bleibt es `PGR-E4003`. Nach einer Gruppe mit `Flush` endet das Warten auf die Antworten mit der Fehlerantwort (`erwarteteAntworten`). Nachweis gegen PostgreSQL am Binary mit handgeschriebenen Aufzeichnungen (Gegenstand: Einspielen einer Aufzeichnung mit Ende gegen einen Server; das Ende bewirkt der Server, `play` löst es nicht aus), dazu die Grenzen der Einstufung am Binary: ohne `--allow-recorded-errors` und mit einer Aufzeichnung ohne Ende an der Stelle ist es `PGR-E4003` mit Exit-Code 4.

**Geschnitten aus `slice-v1-abschluss-verbindungsende-einspielen`** (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap): Jener Slice (Schritt 38a) liefert die Einstufung für die einfache Anfrage, das Schließen ohne `Terminate` im Port `Einspielziel` und im Einspiel-Adapter und einen Fall am Binary (erwartet, nächste Session); dieser Slice (Schritt 38b) setzt ihn voraus und wendet die Einstufung auf die Extended-Interaktion an. Die Methode `fehlerantwort` ist derselbe Pfad für beide Formen; neu sind die Erkennung des aufgezeichneten Endes in der letzten Gruppe, die Tests der Gruppen und die Fälle am Binary.

**Übernommen aus `slice-v1-abschluss-verbindungsende-einspielen`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-16): die Extended-Interaktion mit Ende und die Fälle am Binary, in denen das Ende nicht erwartet ist (ohne Option, aufgezeichnet ohne Ende an der Stelle).

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): das Einspielen einer Aufzeichnung mit Ende, soweit es die Extended-Interaktion betrifft.

**Gebunden durch `slice-v1-abschluss-verbindungsende-einspielen`** (§1, §6): Port, Adapter und Schließen ohne `Terminate` liefert jener Slice; dieser Slice ändert sie nicht (Rückführung in §4). Die Randformen zur Stelle des Endes, zu den Gründen der Ablehnung und zum Sendefehler nach der gelesenen Fehlerantwort gelten hier unverändert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Einstufung der einfachen Anfrage, das Schließen ohne `Terminate`, Port und Einspiel-Adapter — `slice-v1-abschluss-verbindungsende-einspielen` (Vorgänger, Klasse 1): sonst läge dieser Plan bei etwa 119 % der Marke.
- Das Ende mit `--compare-responses` (Extended eingeschlossen), der Hilfetext von `--allow-recorded-errors` und `--compare-responses` — `slice-v1-abschluss-verbindungsende-vergleich` (dort §1 und DoD; Klasse 1): andere Verzweigung des Service (`play_vergleich.go`) und der CLI-Adapter als zweite Schicht.
- Das Aufzeichnen und das Wiedergeben des Endes — `-record`, `-upstream`, `-upstream-extended`, `-replay`, `-replay-extended`, `-replay-extended-flush`: `play` löst das Ende nicht aus, es sendet die aufgezeichneten Anfragen, das Ende bewirkt der Server (`LH-FA-25`, Out-of-Scope). Der Beleg, dass der Upstream-Adapter das Ende einer Extended-Interaktion so erkennt, wie die handgeschriebene Aufzeichnung es zeigt, liegt in `slice-v1-abschluss-verbindungsende-upstream-extended`; die Fixtures dieses Slice bilden nur ab, was die Spezifikation sagt und was PostgreSQL beim Einspielen liefert.
- Ein Verbindungsende ohne Fehlerantwort (Netzwerkabbruch) als erwartetes Ende — Out-of-Scope von `LH-FA-25`; es bleibt `PGR-E4003`.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13). Die Handbuch-Lücke bis dahin (der Satz zu `FATAL`/`PANIC` für eine handgeschriebene Aufzeichnung mit `connection_closed` stimmt nicht mehr) ist akzeptiert, Entscheidung des Nutzers vom 2026-10-12 (`slice-v1-abschluss-verbindungsende-doku` §1).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): `play` ohne Vergleich stuft eine gelesene Fehlerantwort `FATAL` oder `PANIC` in einer Extended-Interaktion ein wie bei der einfachen Anfrage: mit `--allow-recorded-errors` und einem aufgezeichneten Ende in der letzten Gruppe ist sie erwartet (keine Meldung; die Session endet, die nächste läuft, auch ohne `--continue-on-error`; `eingespielt` zählt die Interaktion), sonst `PGR-E4003`, sofort und auch mit `--continue-on-error`; das gilt für das Ende in einer Gruppe mit `Sync` wie in einer Gruppe mit `Flush` (dort endet das Warten auf die Antworten mit der Fehlerantwort), für Gruppen davor, für eine Fehlerantwort in einer früheren Gruppe derselben Interaktion und für eine normale Antwort auf eine Interaktion mit aufgezeichnetem Ende; nach dem Ende sendet `play` weder `Terminate` noch weitere Gruppen (Test im Service und im Einspiel-Adapter).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Gegen PostgreSQL am Binary: Eine handgeschriebene Extended-Aufzeichnung mit Ende, deren Ausführung den Backend beendet (`FATAL`, SQLSTATE `57P01`), endet mit `--allow-recorded-errors` erwartet, mit Exit-Code 0 und mit der nächsten Session, in einer Gruppe mit `Sync` wie in einer Gruppe mit `Flush`; ohne die Option und mit einer Aufzeichnung ohne Ende an der Stelle (aufgezeichnet `ERROR`, der Server sendet `FATAL`) ist es `PGR-E4003` mit Exit-Code 4 (Integrationstest, eine Aufzeichnung als eigene Datei je Fall).
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
| `internal/hexagon/services/play.go` | update (Umfang entscheidet der Architect) | Erkennung des aufgezeichneten Endes in der letzten Gruppe (Extended-Seite); `extended`, `fehlerantwort` und `session` führen den Pfad der einfachen Anfrage schon durch (`errSessionEnde` läuft aus `extended` zurück); Produktcode 0 bis 25 |
| `internal/adapters/driven/postgres/einspielen.go` | keine Änderung erwartet | Port und Schließen liefert `slice-v1-abschluss-verbindungsende-einspielen`; ein Befund der Tests an dieser Stelle ist die zweite Rückführung in §4 |
| `internal/hexagon/services`, `internal/adapters/driven/postgres` (Tests) | update, neu | Randformen aus §6; Fälle der Gruppen im Service (Zeilen derselben Tabellen und eigene Tests), Adapter-Test für das Senden der Gruppen vor dem Schließen; die Hilfen des Gebers werden benutzt |
| `test/integration` (`play_ende_extended_e2e_test.go`, `testdata`, eine Datei je Fall) | neu | Binary gegen PostgreSQL mit handgeschriebenen Aufzeichnungen: Extended mit `Sync`, Extended mit `Flush`, einfache Anfrage ohne Option, aufgezeichnet ohne Ende an der Stelle; Name der Datei für die Grenzen (`play_ende_grenzen_e2e_test.go`) wählt der Architect |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-einspielen` liegt in `done/` (Voraussetzung: Einstufung, Port und Adapter); dieser Slice ist Schritt 38b der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis:** vom Planner beim Anspruch einzutragen, danach nicht mehr geändert); **über 1000** geht der Slice zurück zur Zerlegung. **Meldeschwelle 800:** Erreicht die Summe 800, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik beim Schnitt vor dem Anspruch (Faktor 10, gemessen 10,2 in `-replay`; Maße: Service-Test rund 36 Zeilen, Adapter-Test rund 44, E2E-Fall mit Testdaten rund 66): **430 bis 820, Mitte 625 und damit 63 % der Marke** (§8), unter der Grenze von 70 %, das obere Ende bei 82 %; beim Anspruch ersetzt durch das gemessene Maß des Gebers (§7 dort); Zuschnitt eines weiteren Schnitts: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`).
- `in-progress` → `next`: Die Extended-Interaktion verlangt eine Änderung des Ports `Einspielziel` oder des Einspiel-Adapters, die `slice-v1-abschluss-verbindungsende-einspielen` geliefert hat — zurück zur Zerlegung (Architect und Planner), nicht still im Code.
- `in-progress` → `next`: Die Extended-Interaktion verlangt eine Änderung am Bootstrap, die Logik trägt (`AGENTS.md` §3.13), oder am CLI-Adapter — dritte Schicht, zurück zur Zerlegung.
- `in-progress` → `open`: Der Nutzer ändert, was `--allow-recorded-errors` ohne Vergleich erlaubt — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code. Ort: `LH-FA-20.a` (Tabelle der Fehlerregeln, *Interaktion*, *Gruppen*, *Ende einer Session*), `LH-FA-25.a` *Form*, `LH-FA-18.a`. Die Randformen der einfachen Anfrage (Ende bei `FATAL` und `PANIC`, ohne Vergleich mit und ohne Option, ohne Ende an der Stelle, normale Antwort, `--continue-on-error`, Sendefehler nach gelesener `FATAL`, kein `Terminate`, Zähler) stehen in `slice-v1-abschluss-verbindungsende-einspielen` §6 und gelten hier unverändert; hier steht, was die Extended-Interaktion hinzufügt. Der Architect prüft vor dem Code, ob jede Zeile entschieden ist; als Setzung des Planners zur Prüfung vorgelegt sind (d) und (e) unten.

- *Gegenstand und Fixtures:* wie im Geber §6: Das Ende bewirkt der Server, `play` löst es nicht aus. Jeder Fall des Services und jeder E2E-Fall baut seine Aufzeichnung selbst (eine Datei je E2E-Fall, kein gemeinsamer Zustand), damit eine Mutation nur den Fall färbt, den sie bricht; die Hilfen des Gebers werden benutzt; die Zeilen je E2E-Fall mit Testdaten meldet der Implementer in §7 (Messung).
- *Stelle der Nachricht (`AGENTS.md` §3.12), Extended-Interaktion:* (a) *Aufgezeichnetes `connection_closed`:* letzte Server-Nachricht der letzten Gruppe der Interaktion; das Format lehnt jede andere Stelle beim Laden ab (`PGR-E3003`, `Session.ValidateEnd`). (b) *Gruppe mit dem Ende ist nicht die letzte der Interaktion, Interaktion mit Ende nicht die letzte der Session:* gibt es nicht. (c) *Gelesene Fehlerantwort `FATAL` oder `PANIC` in der Gruppe mit dem Ende:* in der `Sync`-Gruppe nach den Antworten auf Nachrichten davor (etwa `ParseComplete`, `BindComplete`), in der `Flush`-Gruppe ebenso; beide erwartet, `play` liest und sendet nichts mehr. (d) *In einer früheren Gruppe als der mit dem Ende (Gruppe mit `Flush`, deren Antworten `play` vor der nächsten Gruppe liest, oder Gruppe mit `Sync`):* die Einstufung gilt für die Interaktion, nicht für die Gruppe — erwartet, wenn die Aufzeichnung der Interaktion ein Ende trägt, gleich an welcher Stelle (`LH-FA-20.a` *Interaktion*, Setzung des Planners zur Prüfung des Architect, weil die Gruppen nach der Fehlerantwort nicht gesendet werden). (e) *Nach einer Fehlerantwort `ERROR`, die nicht abbricht, folgt `FATAL` in einer späteren Gruppe:* die Gruppen werden nach dem `ERROR` ohne Warten dazwischen gesendet (`fehlerGesehen`); die `FATAL` wird bis zum `ReadyForQuery` gelesen und wie in (c) eingestuft. (f) *Reihenfolge der Gründe der Ablehnung:* ein Grund (das Ende ist nicht erwartet), keine Reihenfolge. Entschieden, `LH-FA-20.a`, `LH-FA-25.a`: (a) bis (c) und (f); (d) und (e) legt der Architect vor dem Code fest und trägt sie in die Spezifikation.
- *Antworten nach einer `Flush`-Gruppe:* `erwarteteAntworten` zählt die erwarteten Antworten; die Fehlerantwort des Endes beendet das Zählen (`n = 0`), und der Service kehrt mit dem Ende der Session zurück, ohne weitere Gruppe zu senden — `LH-FA-20.a` *Gruppen*, `LH-FA-25.a` *Form*.
- *Sendefehler zwischen Gruppen und nach der gelesenen `FATAL`:* ändert die Einstufung nicht; maßgeblich ist der Lesestand (Geber §6). Der Sender des Adapters sendet Gruppen nacheinander unabhängig vom Lesen; das Schließen ohne `Terminate` bricht ihn ab, ein Fehler danach bleibt ohne Folge (Port `Einspielziel`, `Gruppe`, *Fehler nach Schliesse*).
- *Zähler und Nachricht:* `eingespielt` zählt die Interaktion mit Ende; `verglichen` und `abweichend` bleiben 0 — `LH-FA-24.a` *Meldung*.
- *Fälle am Binary:* Ein Server, der den Backend während der Ausführung beendet, sendet `FATAL` `57P01`; die Aufzeichnung der Fälle ohne Ende an der Stelle trägt ein `ERROR` an der Stelle und kein Ende. Die Version des Servers ist die Referenz 17 (Risiko unten).

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×; Tests der Gruppen sind teuer): vor dem Anspruch geschnitten; die Rückführung in §4 benennt Marke und Meldeschwelle — **Ausgang:** offen bis Closure.
- Port oder Adapter aus `slice-v1-abschluss-verbindungsende-einspielen` tragen den Pfad der Gruppen nicht (Schließen während der Sender eine Gruppe sendet): die zweite Rückführung in §4 — **Ausgang:** offen bis Closure.
- Randform (d) und (e) sind Setzungen des Planners, die der Architect vor dem Code prüft; entscheidet er anders, ändert sich die Einstufung einer Fehlerantwort in einer früheren Gruppe — **Ausgang:** offen bis Closure.
- Handgeschriebene Extended-Aufzeichnungen bilden nur ab, was die Spezifikation sagt; die Fälle am Binary belegen das Einspielen gegen PostgreSQL, nicht das Aufzeichnen (`BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`, 3×; Aufzeichnen: `slice-v1-abschluss-verbindungsende-upstream-extended`) — **Ausgang:** offen bis Closure (erwartet: weiter offen, keine neue Adresse; der Beleg liegt hier für das Einspielen).
- Das Verhalten des Servers (`57P01` bei `pg_terminate_backend` in einer Extended-Ausführung) ist nur gegen die Referenzversion 17 geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×) — **Ausgang:** offen bis Closure (erwartet: weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`).

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×), `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×), `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (3×) und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×; §6) — der Plan nennt Schätzung und Schnitt vorab (§4); tritt die Rückführung ein, ist es das dritte Auftreten von `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`. Keine weiteren Treffer. Der Zähler-Stand ist beim Anspruch aus dem Register zu lesen.

Nachgezählt beim Anlegen (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Schnitt von `slice-v1-abschluss-verbindungsende-einspielen`): Liefer-Punkte: 2 (Einstufung der Extended-Interaktion im Service mit Randformen, Test im Service und im Adapter; Binary gegen PostgreSQL). Schichten: zwei (Kern mit Service, `ARC-002`; Einspiel-Adapter, `ARC-007`, vor allem Tests); die Abdeckungstabellen in `docs/user` zählen nicht als Schicht (Begründung in `slice-v1-abschluss-verbindungsende-replay` §8); **konservativ gezählt mit `docs/user` als dritter Schicht: drei** — markiert, bis der Nutzer entscheidet. Schätzung in der Metrik von §4 (Faktor 10): Produktcode 0 bis 25 (Erkennung des Endes in der letzten Gruppe; der Pfad der einfachen Anfrage besteht), Tests im Service 190 bis 320 (acht bis zehn Fälle: `Sync` und `Flush` mit Ende, ohne Option, ohne Ende an der Stelle, Gruppen davor, Fehlerantwort in einer früheren Gruppe, normale Antwort, kein weiteres Senden nach dem Ende, Sendefehler zwischen Gruppen; als Zeilen derselben Tabellen, wo der Geber Tabellen führt, sonst eigene Tests zu je rund 30 bis 36 Zeilen), im Adapter 40 bis 90 (ein bis zwei Tests), E2E mit Testdaten vier Fälle zu je 63 bis 70, also 250 bis 280, Tabellen 20 bis 35: **430 bis 820, Mitte 625, 63 % der Marke 1000**, oberes Ende 82 %; Meldeschwelle 800. Die Summe aus beiden Teilen (Mitte 690 und 625) liegt bei 1315 gegen 1185 ungeteilt; der Mehraufwand sind die zweite Gruppe von Hilfen und der zweite Adapter-Lauf.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
