# Slice slice-v1-abschluss-verbindungsende-einspielen: Verbindungsende beim Einspielen ohne Vergleich (`--allow-recorded-errors`, einfache Anfrage)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-20.a` · `SPEC-034` · `ARC-002` · `ARC-004` · `ARC-007`

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

**Ziel:** `play` ohne `--compare-responses` behandelt eine Fehlerantwort `FATAL` oder `PANIC` des Servers auf eine **einfache Anfrage** (`LH-FA-20.a` *Fehlerregeln*, *Ende einer Session*; `LH-FA-25.a`): Trägt die aufgezeichnete Interaktion ein Ende und ist `--allow-recorded-errors` gesetzt, ist das Ende erwartet — keine Meldung, `play` wartet nicht auf das Verbindungsende und sendet kein `Terminate`, sondern schließt die Verbindung, die Session endet, die nächste läuft, auch ohne `--continue-on-error`; sonst bleibt es `PGR-E4003`, sofort, auch mit `--continue-on-error`. Ein aufgezeichnetes Ende, auf das der Server normal antwortet, ist kein Fehler. Nachweis gegen PostgreSQL mit einer handgeschriebenen Aufzeichnung, deren Anfrage `pg_terminate_backend(pg_backend_pid())` ist. Gegenstand der Einspiel-Tests ist das Einspielen einer Aufzeichnung mit Ende gegen einen Server; das Ende des Servers ist nicht einzuspielen, `play` löst es nicht aus und sendet es nicht (`connection_closed` ist ein Eintrag der aufgezeichneten Antworten, keine Client-Nachricht).

**Geschnitten vor dem Anspruch in `-einspielen` (dieser Slice, Schritt 38a) und `slice-v1-abschluss-verbindungsende-einspielen-extended` (Schritt 38b)** (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap): Die Schätzung des Planners mit dem gemessenen Faktor 10 lag ungeteilt bei 830 bis 1540, Mitte 1185 und damit 119 % der Marke 1000, über der Grenze von 70 % (die frühere Schätzung 700 bis 1100 rechnete mit dem Faktor 5). Schnitt nach Lieferwert: Dieser Slice liefert die Einstufung, das Schließen ohne `Terminate` (Port und Adapter) und den Beleg am Binary für die **einfache Anfrage**; die **Extended-Interaktion** mit Ende (Gruppe mit `Sync` oder `Flush`) und die Grenzen am Binary (ohne Option, ohne Ende an der Stelle) liefert `slice-v1-abschluss-verbindungsende-einspielen-extended`, der diesen Slice voraussetzt.

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): das Einspielen einer Aufzeichnung mit Ende. Das Laden liefert jener Slice; dieser setzt ihn voraus. Davon geht die Extended-Interaktion und der Beleg der Grenzen am Binary an `slice-v1-abschluss-verbindungsende-einspielen-extended`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Extended-Interaktion mit Ende (Gruppe mit `Sync` oder `Flush`; Ende in der letzten Gruppe, Gruppen davor, Antworten nach einer `Flush`-Gruppe, Sendefehler zwischen Gruppen) und die Fälle am Binary, in denen das Ende nicht erwartet ist (ohne Option, aufgezeichnet ohne Ende an der Stelle) — `slice-v1-abschluss-verbindungsende-einspielen-extended` (Klasse 1, setzt diesen Slice voraus; dessen §1 und DoD nehmen es an): mit ihnen läge dieser Plan bei etwa 119 % der Marke (§8). Bis jener Slice schließt, bleibt eine Fehlerantwort `FATAL` oder `PANIC` in einer Extended-Interaktion `PGR-E4003`, auch mit `--allow-recorded-errors` und aufgezeichnetem Ende; die Extended-Seite der Erkennung des aufgezeichneten Endes hat hier keinen Code und keinen Test (§6, Zwischenstand). Die Fälle des Services, in denen das Ende der einfachen Anfrage nicht erwartet ist, gehören zur Einstufung und bleiben hier, weil jede Bedingung der Einstufung ihre Mutation braucht (`AGENTS.md` §3.10).
- Das Ende mit `--compare-responses` (erwartet bei gleichem SQLSTATE, Abweichung `PGR-E5004`, Zähler `verglichen`, Szenario 16 Satz 3), der Hilfetext von `--allow-recorded-errors` und `--compare-responses` im Binary — `slice-v1-abschluss-verbindungsende-vergleich` (dort §1 und DoD; Sendung des Planners, 2026-10-12): andere Verzweigung des Service (`play_vergleich.go`) und der CLI-Adapter als zweite Schicht; mit ihm läge dieser Plan über der Grenze (drei Schichten, Diff über einer Review-Sitzung). Bis dahin nennt der Hilfetext der Option das Ende nicht; er sagt zu, was ein Test prüft (`AGENTS.md` §3.11).
- Das Aufzeichnen und das Wiedergeben des Endes — `-record`, `-upstream`, `-replay`: `play` löst das Ende nicht aus, es sendet die aufgezeichneten Anfragen, das Ende bewirkt der Server (`LH-FA-25`, Out-of-Scope).
- Ein Verbindungsende ohne Fehlerantwort (Netzwerkabbruch) als erwartetes Ende — Out-of-Scope von `LH-FA-25`; es bleibt `PGR-E4003`.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13). Die Handbuch-Lücke bis dahin (der Satz zu `FATAL`/`PANIC` für eine handgeschriebene Aufzeichnung mit `connection_closed` stimmt nicht mehr) ist akzeptiert, Entscheidung des Nutzers vom 2026-10-12 (dort §1).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): `play` ohne Vergleich stuft eine gelesene Fehlerantwort `FATAL` oder `PANIC` auf eine einfache Anfrage ein: mit `--allow-recorded-errors` und einem aufgezeichneten Ende in der Interaktion ist sie erwartet (keine Meldung; die Session endet, die nächste läuft, auch ohne `--continue-on-error`; `eingespielt` zählt die Interaktion, `verglichen` und `abweichend` bleiben 0), sonst `PGR-E4003`, sofort und auch mit `--continue-on-error` (ohne die Option, aufgezeichnet ohne Ende an der Stelle, `FATAL` in einer früheren Interaktion; ein aufgezeichnetes Ende, auf das der Server normal antwortet, ist kein Fehler); eine Extended-Interaktion bleibt `PGR-E4003` (Test im Service).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md): Nach einem Ende des Servers, erwartet oder nicht, wartet `play` nicht auf das Verbindungsende, sendet kein `Terminate` und schließt die Verbindung mit `SchliesseOhneTerminate` (Test im Service: der Aufruf; Test im Einspiel-Adapter: kein `Terminate` auf der Leitung, ein zweites Schließen ohne Wirkung); ein Fehler beim Schließen ist keine Meldung und ändert den Exit-Code nicht (getragen von der Form des Ports ohne Fehlerrückgabe, kein eigener Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): Gegen PostgreSQL am Binary: Eine Aufzeichnung mit Ende, deren Anfrage den Backend beendet (`FATAL`, SQLSTATE `57P01`), endet mit `--allow-recorded-errors` erwartet, mit Exit-Code 0 und mit der nächsten Session (Integrationstest, ein Fall, Aufzeichnung als eigene Datei).
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
| `internal/hexagon/services/play.go` | update | Einstufung des Endes ohne Vergleich bei einer einfachen Anfrage (`fehlerantwort`, `interaktion`, `session`), Ende der Session |
| `internal/hexagon/ports/driven/einspielziel.go` | update | Port-Vertrag: neue Methode `SchliesseOhneTerminate()` an `EinspielSession` (ohne Rückgabe, jederzeit aufrufbar wie `Schliesse`; ein Schließen danach ist ohne Wirkung); Entscheidung des Architect, 2026-10-16 |
| `internal/adapters/driven/postgres/einspielen.go` | update | `SchliesseOhneTerminate`; `Schliesse` und die neue Methode teilen einen Pfad, der nur beim ersten Schließen wirkt und nur dort `Terminate` sendet |
| `internal/hexagon/services`, `internal/adapters/driven/postgres` (Tests) | update | Randformen aus §6; die Tests, die die alte Regel prüfen, werden hier nachgezogen: `TestPlayFehlerantwort` und `TestPlayExtendedFatal` (`play_extended_test.go`; erwarten nach `FATAL` und `PANIC` das Schließen ohne `Terminate`, die Extended-Interaktion bleibt `PGR-E4003`, dort mit einem Fall mit aufgezeichnetem Ende als Zwischenstand bis `-einspielen-extended`), dazu der Fake `fakeEinspiel` (hält das Schließen ohne `Terminate` in einem Feld fest, das Ereignis bleibt `schliesse`) und eine Prüfung des Aufrufs in `TestVergleichFatal` (`play_vergleich_test.go`); `TestRunPlayFehlerantwort` im Bootstrap bleibt unverändert (er prüft `ERROR`, nicht `FATAL`) |
| `test/integration` (`play_ende_e2e_test.go`, `testdata/play-ende-erwartet.yaml`) | neu | Binary gegen PostgreSQL mit handgeschriebener Aufzeichnung (Version 2, zwei Sessions), ein Fall, Aufzeichnung als eigene Datei |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` liegt in `done/` (Voraussetzung); dieser Slice ist Schritt 38a der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis:** `33e5eab`, Stand des Hauptzweigs beim Anspruch, danach nicht mehr geändert; die Commits danach ändern nur Pläne und die Verzeichnis-Position, keinen Pfad der Metrik); **über 1000** geht der Slice zurück zur Zerlegung. **Meldeschwelle 800:** Erreicht die Summe 800, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik beim Schnitt vor dem Anspruch (Faktor 10, gemessen 10,2 in `-replay`; Maße: Service-Test rund 36 Zeilen, Adapter-Test rund 44, E2E-Fall mit Testdaten rund 66): **500 bis 880, Mitte 690 und damit 69 % der Marke** (§8), unter der Grenze von 70 %, das obere Ende bei 88 %; Zuschnitt eines weiteren Schnitts: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`).
- `in-progress` → `next`: Das Schließen ohne `Terminate` verlangt eine Änderung am Bootstrap, die Logik trägt (`AGENTS.md` §3.13), oder am CLI-Adapter — dritte Schicht, zurück zur Zerlegung.
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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code. Ort: `LH-FA-20.a` (Tabelle der Fehlerregeln, *Interaktion*, *Ende einer Session*), `LH-FA-25.a`. „`--allow-recorded-errors` ohne `--compare-responses` erlaubt ein aufgezeichnetes Ende“ und „gleiche Stelle, gleicher SQLSTATE“ bestätigte der Nutzer am 2026-10-12; die übrigen sind Setzungen des Architect.

- *Gegenstand und Fixtures:* Eingespielt wird eine Aufzeichnung mit Ende gegen einen Server; das Ende bewirkt der Server (`pg_terminate_backend(pg_backend_pid())`), `play` löst es nicht aus und sendet es nicht. Jeder Fall des Services und jeder E2E-Fall baut seine Aufzeichnung selbst (eine Datei je E2E-Fall, kein gemeinsamer Zustand zwischen Fällen), damit eine Mutation nur den Fall färbt, den sie bricht (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Erfahrung `-replay`); die Zeilen je Fall meldet der Implementer in §7 (Messung).
- *Stelle der Nachricht (`AGENTS.md` §3.12), einfache Anfrage:* (a) *Aufgezeichnetes `connection_closed`:* letzter Eintrag der Antworten der Interaktion; das Format lehnt jede andere Stelle beim Laden ab (`PGR-E3003`, `Session.ValidateEnd`). (b) *Interaktion mit Ende ist nicht die letzte der Session:* gibt es nicht, das Ende ist das letzte Ereignis seiner Session. (c) *Gelesene Fehlerantwort `FATAL` oder `PANIC`:* als erste Antwort der Interaktion, nach `RowDescription` und `DataRow`: gleich behandelt, erwartet, wenn die Aufzeichnung der Interaktion ein Ende trägt, gleich an welcher Stelle (`LH-FA-20.a` *Interaktion*). (d) *Nach der Fehlerantwort:* `play` liest nichts mehr und sendet nichts mehr außer dem Schließen. (e) *Gründe der Ablehnung:* es gibt einen (das Ende ist nicht erwartet: ohne Option, ohne Ende an der Stelle, Ende des Servers in einer früheren Interaktion); eine Reihenfolge von Gründen gibt es nicht. Entschieden, `LH-FA-20.a` *Interaktion*, `LH-FA-25.a` (Architect, 2026-10-16, Commit vor dem Code): (a) bis (e) tragen. Die Setzung (d) des Planners (Stelle ist die Interaktion, nicht Gruppe oder Position) steht dort für beide Formen; in diesem Slice zählt sie für die Position in den Antworten der einfachen Anfrage, die Gruppen gehören `-einspielen-extended`.
- *Erkennung des aufgezeichneten Endes (einfache Anfrage):* die letzte Antwort der Interaktion ist `connection_closed` (`endetEinfach` in `replay.go` liefert das; keine neue Hilfe im Model). Eine Extended-Interaktion trägt keine Antworten außerhalb ihrer Gruppen und ist damit hier ohne Ende: `PGR-E4003` bleibt.
- *Verbindungsverlust ohne Fehlerantwort trotz aufgezeichnetem Ende, Server schließt ohne Antwort:* `PGR-E4003` (Tabelle der Fehlerregeln), kein erwartetes Ende; das Schließen sendet `Terminate` wie sonst, weil kein Ende des Servers gelesen ist. Ein Fall in der Tabelle des Service.
- *Fehlerantwort `ERROR` in einer Interaktion mit aufgezeichnetem Ende:* gewöhnlicher Fehler, mit Option erwartet (`LH-FA-20.a` *Interaktion*), die Interaktion läuft bis `ReadyForQuery`, die Session endet normal mit `Terminate`; ein Fall in der Tabelle.
- *Erstes und zweites Signal während einer Interaktion mit Ende:* keine neue Entscheidung, es gilt der Satz „Fehlerregeln wie ohne Signal“ (`LH-FA-20.a` *Abbruchsignal*); kein neuer Code, kein eigener Fall. Das zweite Signal sendet `Terminate`, soweit die Verbindung es sofort annimmt; ein Rennen mit dem Schließen nach der Fehlerantwort ist folgenlos (zweites Schließen ohne Wirkung).
- *Timeout:* `play` wartet ohne eigene Frist (`LH-FA-20.a` *Interaktion*); unverändert.
- *Mit `--compare-responses`:* die Einstufung (erwartet, `PGR-E5004`, Zähler) ist nicht Gegenstand (`-vergleich`). Das Schließen ohne `Terminate` gilt aber für jede gelesene Fehlerantwort `FATAL` oder `PANIC` an einem Ort, vor der Verzweigung in `fehlerantwort`, und damit auch mit Vergleich (`LH-FA-20.a` *Ende einer Session* nennt „erwartet oder nicht“): Sonst bliebe die Zusage an der Hälfte der Pfade falsch und `-vergleich` bekäme eine Sendung, die sein §1 ausschließt. Ein Fall mit Vergleich gehört in die Tests dieses Slice (ein Fall in `TestVergleichFatal` prüft den Aufruf); die übrigen Tests mit Vergleich bleiben unverändert, wenn der Fake das Schließen ohne `Terminate` zusätzlich zum Ereignis `schliesse` in einem Feld festhält.
- *Ende bei `FATAL` und `PANIC`:* gleich behandelt — `LH-FA-25.a`, `LH-FA-20.a` *Interaktion*.
- *Ohne Vergleich, mit `--allow-recorded-errors` und aufgezeichnetem Ende:* erwartet, gleich welcher SQLSTATE und an welcher Stelle der Interaktion — `LH-FA-20.a` *Interaktion* (Fehlerregeln).
- *Ohne Vergleich, ohne `--allow-recorded-errors`:* `PGR-E4003`, sofort — Tabelle der Fehlerregeln.
- *Aufzeichnung ohne Ende an der Stelle (ohne Vorbild, ohne `FATAL`; aufgezeichnet `ERROR`, der Server sendet `FATAL`):* `PGR-E4003`, auch mit `--allow-recorded-errors` — `LH-FA-20.a` *Interaktion*.
- *Aufgezeichnetes Ende, der Server antwortet normal:* kein Fehler, die Folge wird ohne Vergleich verworfen; die Session läuft mit der nächsten Interaktion weiter (die Aufzeichnung hat nach dem Ende keine) — `LH-FA-20.a` *Interaktion*, `LH-FA-25.a`.
- *`FATAL` in einer früheren Interaktion als der aufgezeichneten:* `PGR-E4003`, die Aufzeichnung zeigt dort kein Ende — `LH-FA-20.a`.
- *`--continue-on-error`:* ändert `PGR-E4003` nicht („immer sofort“); ein erwartetes Ende bricht nicht ab — Tabelle der Fehlerregeln.
- *Sendefehler und gelesene `FATAL`:* maßgeblich ist der Lesestand (`LH-FA-20.a` *Interaktion*, `LH-FA-02.b` *Gelesen*). Bei der einfachen Anfrage gibt es nach der gelesenen Fehlerantwort kein Senden außer dem Schließen (ohne Rückgabe); scheitert das Senden der Anfrage, ist das `PGR-E4003` und es gibt keine Fehlerantwort. Der Fall „Sendefehler nach gelesener `FATAL`“ gehört damit `-einspielen-extended`.
- *Kein `Terminate` nach einem Ende (erwartet oder nicht):* Der Service ruft `SchliesseOhneTerminate` sofort nach der gelesenen Fehlerantwort `FATAL` oder `PANIC` (einfache und Extended-Interaktion, mit und ohne Vergleich, ein Ort in `fehlerantwort`); das spätere `Schliesse` der Session ist ohne Wirkung. Ein Fehler beim Schließen ist ohne Meldung und ohne Folge für den Exit-Code (Port ohne Rückgabe) — `LH-FA-20.a` *Ende einer Session*. Das Ende der Session ohne `FATAL` (regulär, `ERROR`, Verbindungsverlust) schließt mit `Schliesse` und `Terminate`.
- *Extended:* Zwischenstand, siehe Risiken; die Entscheidung (Ende in der letzten Gruppe mit `Sync` oder `Flush`, Antworten nach einer `Flush`-Gruppe, letzte Gruppe mit `Flush` ohne Warten auf `ReadyForQuery`) steht in `LH-FA-20.a` *Interaktion*, *Gruppen* und `LH-FA-25.a` *Form* und wird in `slice-v1-abschluss-verbindungsende-einspielen-extended` geliefert; dieser Slice erfindet dafür nichts und prüft nur, dass die Extended-Interaktion `PGR-E4003` bleibt (ein Fall in `TestPlayExtendedFatal`, den jener Slice umkehrt).
- *Zähler:* `eingespielt` zählt die Interaktion mit Ende; `verglichen` und `abweichend` bleiben ohne Vergleich 0 — `LH-FA-24.a` *Meldung*.

**Risiken:**

- Zwischenstand: Nach diesem Slice ist die Extended-Interaktion mit aufgezeichnetem Ende beim Einspielen nicht zugesagt und nicht getestet; eine Fehlerantwort `FATAL` oder `PANIC` bleibt dort `PGR-E4003`, auch mit `--allow-recorded-errors` (das Schließen ohne `Terminate` gilt auch dort, der Port ist derselbe) — **Ausgang:** offen bis Closure (erwartet: *eingetreten*, Folge-Slice `slice-v1-abschluss-verbindungsende-einspielen-extended`, nimmt an: dessen §1 und DoD nennen die Extended-Interaktion und die Grenzen am Binary).
- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×): vor dem Anspruch geschnitten (Extended in `-einspielen-extended`); die Rückführung in §4 benennt Marke und Meldeschwelle — **Ausgang:** offen bis Closure.
- Das Schließen ohne `Terminate` braucht eine Änderung des Ports `Einspielziel` und des Adapters; die Zahl der Schichten (Kern mit Port, Einspiel-Adapter) ist zwei, und der Architect hat vor dem Code geprüft (2026-10-16), dass Bootstrap und CLI-Adapter unberührt bleiben: Der Bootstrap reicht den Adapter als `driven.Einspielziel` durch, die neue Methode steht an `EinspielSession`, kein Bootstrap-Code und kein Fake im Bootstrap implementiert sie — **Ausgang:** offen bis Closure.
- Das Verhalten des Servers (`57P01` bei `pg_terminate_backend`) ist nur gegen die Referenzversion 17 geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×); der Vertrag stellt es dar (Architect, 2026-10-16): `LH-FA-20.a` *Interaktion* nennt es als Grenze, `SPEC-029` Punkt 2 (volle Suite nur gegen die Referenz) und `SPEC-053` Punkt 14 (5) (Abnahmeszenario 18 und Szenario 16 Satz 3 nur gegen die Referenz) tragen es; die Matrix fährt den Fall nicht mit, ohne Änderung von [ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md) und `SPEC-053`. Empfehlung des Architect: nicht aufnehmen (die Teilmenge ist die Wahl des Nutzers vom 2026-10-14, der Fall hängt an einem SQLSTATE, den die Hauptversionen 14 bis 18 gleich liefern, ein Gewinn ist nicht belegt); will der Nutzer es, ist das eine Änderung von [ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md) und `SPEC-053` und eine Entscheidung des Nutzers — **Ausgang:** offen bis Closure (erwartet: weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`).

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×), `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×) und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×; §6) — der Plan nennt Schätzung und Schnitt vorab (§4); bei Eintritt der Rückführung ist `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` zum dritten Mal belegt. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13): Liefer-Punkte: 3 (Einstufung im Service, Schließen ohne `Terminate` im Port und Adapter, E2E gegen PostgreSQL). Schichten: zwei (Kern mit Service und Port, `ARC-002`, `ARC-004`; Einspiel-Adapter, `ARC-007`). Der CLI-Adapter (Hilfetext) liegt in `slice-v1-abschluss-verbindungsende-vergleich`, der Bootstrap bleibt unberührt. Schätzung des Planners: 700 bis 1100 Zeilen mit Tests.

Nachgezählt beim Schnitt vor dem Anspruch (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14): Ungeteilt, mit dem Faktor 10 (gemessen 10,2 in `-replay`, `-replay-extended` 532 mit Produktcode 0) auf Produktcode 70 bis 130 (Service mit Erkennung des Endes beider Formen, Port, Adapter) plus Testdaten und Tabellen 60 bis 110: **830 bis 1540, Mitte 1185, 119 % der Marke 1000**, über der Grenze von 70 %; Gegenrechnung von unten (Service 14 bis 17 Tests und Tabellenzeilen zu je 30 bis 36 Zeilen, Adapter 4 Tests, E2E 5 Fälle zu je 66, Nachziehen der alten Tests 50): 970 bis 1365, Mitte rund 1170. Geschnitten nach Lieferwert (einfache Anfrage | Extended-Interaktion mit den Grenzen am Binary), nicht nach Schichten. **Dieser Teil (38a):** Produktcode 60 bis 110 (Einstufung 25 bis 50, Port 10 bis 15, Adapter 20 bis 40), Tests im Service 250 bis 360 (Tabelle mit acht Zeilen: erwartet bei `FATAL`, bei `PANIC`, ohne Option, ohne Ende an der Stelle, `FATAL` früher, normale Antwort, mit und ohne `--continue-on-error`; dazu Zähler, nächste Session, Schließen ohne `Terminate`; Nachziehen der alten Tests 40 bis 60), im Adapter 100 bis 180 (Schließen ohne `Terminate`, Fehler beim Schließen, Sendefehler nach der gelesenen Fehlerantwort), E2E mit Testdaten ein Fall 63 bis 70, Hilfen und Fakes 15 bis 30, Tabellen 10 bis 20: **500 bis 880, Mitte 690, 69 % der Marke 1000**, oberes Ende 88 %; Meldeschwelle 800. Liefer-Punkte: 3 (unverändert). Schichten: zwei (Kern mit Service und Port; Einspiel-Adapter, `ARC-007`), Abdeckungstabellen in `docs/user` keine Schicht (die Nutzerentscheidung dazu, ob `AGENTS.md` §3.13 es festhält, steht aus); **konservativ gezählt mit `docs/user` als dritter Schicht: drei** — markiert, bis der Nutzer entscheidet. **Geprüft vom Architect vor dem Code (2026-10-16, Metrik `AGENTS.md` §3.14):** Zählung bestätigt (Liefer-Punkte 3; Schichten 2 nach dem Vorbild der geschlossenen Vorgänger, Kern mit Service und Port als eine Schicht, Adapter `ARC-007`; Bootstrap, CLI und Model unberührt; `docs/user` wie dort mit Vorbehalt). Produktcode nach dem Entwurf des Architect 45 bis 70 (Service 12 bis 20, Port 8 bis 12, Adapter 25 bis 40), mit dem Faktor 10 rund 450 bis 700; die Schätzung des Planners **500 bis 880, Mitte 690** bleibt als obere Rechnung stehen, keine Korrektur; Marke 1000 und Meldeschwelle 800 unverändert. Der zweite Teil liegt in `slice-v1-abschluss-verbindungsende-einspielen-extended` (Schritt 38b, setzt diesen Slice voraus; dessen §1 nennt die Bindung an die Randformen oben).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
