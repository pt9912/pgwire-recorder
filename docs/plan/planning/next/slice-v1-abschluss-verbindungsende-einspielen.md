# Slice slice-v1-abschluss-verbindungsende-einspielen: Verbindungsende beim Einspielen ohne Vergleich (`--allow-recorded-errors`)

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

**Ziel:** `play` ohne `--compare-responses` behandelt eine Fehlerantwort `FATAL` oder `PANIC` des Servers nach der Aufzeichnung (`LH-FA-20.a` *Fehlerregeln*, *Ende einer Session*; `LH-FA-25.a`): Trägt die aufgezeichnete Interaktion ein Ende und ist `--allow-recorded-errors` gesetzt, ist das Ende erwartet — keine Meldung, `play` wartet nicht auf das Verbindungsende und sendet kein `Terminate`, sondern schließt die Verbindung, die Session endet, die nächste läuft, auch ohne `--continue-on-error`; sonst bleibt es `PGR-E4003`, sofort, auch mit `--continue-on-error`. Ein aufgezeichnetes Ende, auf das der Server normal antwortet, ist kein Fehler. Nachweis gegen PostgreSQL mit einer handgeschriebenen Aufzeichnung, deren Anfrage `pg_terminate_backend(pg_backend_pid())` ist.

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): das Einspielen einer Aufzeichnung mit Ende. Das Laden liefert jener Slice; dieser setzt ihn voraus.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Ende mit `--compare-responses` (erwartet bei gleichem SQLSTATE, Abweichung `PGR-E5004`, Zähler `verglichen`, Szenario 16 Satz 3), der Hilfetext von `--allow-recorded-errors` und `--compare-responses` im Binary — `slice-v1-abschluss-verbindungsende-vergleich` (dort §1 und DoD; Sendung des Planners, 2026-10-12): andere Verzweigung des Service (`play_vergleich.go`) und der CLI-Adapter als zweite Schicht; mit ihm läge dieser Plan über der Grenze (drei Schichten, Diff über einer Review-Sitzung). Bis dahin nennt der Hilfetext der Option das Ende nicht; er sagt zu, was ein Test prüft (`AGENTS.md` §3.11).
- Das Aufzeichnen und das Wiedergeben des Endes — `-record`, `-upstream`, `-replay`: `play` löst das Ende nicht aus, es sendet die aufgezeichneten Anfragen, das Ende bewirkt der Server (`LH-FA-25`, Out-of-Scope).
- Ein Verbindungsende ohne Fehlerantwort (Netzwerkabbruch) als erwartetes Ende — Out-of-Scope von `LH-FA-25`; es bleibt `PGR-E4003`.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): `play` ohne Vergleich stuft eine gelesene Fehlerantwort `FATAL` oder `PANIC` ein: mit `--allow-recorded-errors` und einem aufgezeichneten Ende in der Interaktion ist sie erwartet (keine Meldung; die Session endet, die nächste läuft, auch ohne `--continue-on-error`; `eingespielt` zählt die Interaktion), sonst `PGR-E4003`, sofort und auch mit `--continue-on-error`. Das gilt für eine einfache Anfrage und für eine Extended-Interaktion (Ende in der letzten Gruppe, `Sync` oder `Flush`) (Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md): Nach einem Ende des Servers, erwartet oder nicht, wartet `play` nicht auf das Verbindungsende, sendet kein `Terminate` und schließt die Verbindung; ein Sendefehler nach der gelesenen Fehlerantwort ändert die Einstufung nicht; ein Fehler beim Schließen ist keine Meldung und ändert den Exit-Code nicht (Test im Service und im Einspiel-Adapter).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): Gegen PostgreSQL am Binary: Eine Aufzeichnung mit Ende, deren Anfrage den Backend beendet (`FATAL`, SQLSTATE `57P01`), endet mit `--allow-recorded-errors` erwartet und mit der nächsten Session; ohne die Option, mit einer Aufzeichnung ohne Ende an der Stelle und mit aufgezeichnetem `ERROR` ist es `PGR-E4003` mit Exit-Code 4 (Integrationstest).
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
| `internal/hexagon/services/play.go` | update | Einstufung des Endes ohne Vergleich (`fehlerantwort`, `interaktion`, `session`), Ende der Session |
| `internal/hexagon/ports/driven/einspielziel.go` | update | Port-Vertrag: Schließen nach einem Ende ohne `Terminate` (Name und Form wählt der Architect vor dem Code) |
| `internal/adapters/driven/postgres/einspielen.go` | update | `Schliesse` ohne `Terminate` nach einem Ende; Verhalten bei Sendefehler nach der gelesenen Fehlerantwort |
| `internal/hexagon/services`, `internal/adapters/driven/postgres` (Tests) | update | Randformen aus §6; die Tests, die die alte Regel prüfen (`TestPlayFehlerantwort`, `TestPlayExtendedFatal`, `TestRunPlayFehlerantwort` im Bootstrap), werden hier nachgezogen |
| `test/integration` (`play_ende_e2e_test.go`) | neu | Binary gegen PostgreSQL mit handgeschriebener Aufzeichnung |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` liegt in `done/` (Voraussetzung); dieser Slice ist Schritt 38 der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Prüfung des Architect vor dem Code schätzt den Diff mit Tests auf über 1000 Zeilen (Schätzung des Planners 700 bis 1100; die Zahl der Randformen in §6 treibt die Tests) — Schnitt vor dem Anspruch: Einstufung im Service | Einspiel-Adapter und E2E, Zuschnitt: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`).
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

- *Ende bei `FATAL` und `PANIC`:* gleich behandelt — `LH-FA-25.a`, `LH-FA-20.a` *Interaktion*.
- *Ohne Vergleich, mit `--allow-recorded-errors` und aufgezeichnetem Ende:* erwartet, gleich welcher SQLSTATE und an welcher Stelle der Interaktion — `LH-FA-20.a` *Interaktion* (Fehlerregeln).
- *Ohne Vergleich, ohne `--allow-recorded-errors`:* `PGR-E4003`, sofort — Tabelle der Fehlerregeln.
- *Aufzeichnung ohne Ende an der Stelle (ohne Vorbild, ohne `FATAL`; aufgezeichnet `ERROR`, der Server sendet `FATAL`):* `PGR-E4003`, auch mit `--allow-recorded-errors` — `LH-FA-20.a` *Interaktion*.
- *Aufgezeichnetes Ende, der Server antwortet normal:* kein Fehler, die Folge wird ohne Vergleich verworfen; die Session läuft mit der nächsten Interaktion weiter (die Aufzeichnung hat nach dem Ende keine) — `LH-FA-20.a` *Interaktion*, `LH-FA-25.a`.
- *`FATAL` in einer früheren Interaktion als der aufgezeichneten:* `PGR-E4003`, die Aufzeichnung zeigt dort kein Ende — `LH-FA-20.a`.
- *`--continue-on-error`:* ändert `PGR-E4003` nicht („immer sofort“); ein erwartetes Ende bricht nicht ab — Tabelle der Fehlerregeln.
- *Sendefehler nach gelesener `FATAL`:* ändert die Einstufung nicht; maßgeblich ist der Lesestand — `LH-FA-25.a` *Aufzeichnen* (Lesestand), `LH-FA-02.b` *Gelesen*, Setzung des Architect.
- *Kein `Terminate` nach einem Ende (erwartet oder nicht):* Schließen ohne `Terminate`, ein Fehler dabei ohne Meldung und ohne Folge für den Exit-Code — `LH-FA-20.a` *Ende einer Session*.
- *Extended:* Ende in der letzten Gruppe mit `Sync` oder `Flush`; die Erwartung der Antworten nach einer `Flush`-Gruppe (`erwarteteAntworten`) endet mit der Fehlerantwort — `LH-FA-20.a` *Gruppen*, `LH-FA-25.a` *Form*.
- *Zähler:* `eingespielt` zählt die Interaktion mit Ende; `verglichen` und `abweichend` bleiben ohne Vergleich 0 — `LH-FA-24.a` *Meldung*.

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×): die Rückführung in §4 benennt den Schnitt — **Ausgang:** offen bis Closure.
- Das Schließen ohne `Terminate` braucht eine Änderung des Ports `Einspielziel` und des Adapters; die Zahl der Schichten (Kern mit Port, Einspiel-Adapter) ist zwei, und der Architect prüft vor dem Code, dass das Bootstrap und der CLI-Adapter unberührt bleiben — **Ausgang:** offen bis Closure.
- Das Verhalten des Servers (`57P01` bei `pg_terminate_backend`) ist nur gegen die Referenzversion 17 geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×); ob die Matrix aus `slice-v1-abschluss-postgres-versionen` Abnahmeszenario 18 mitfährt, ist eine Frage an den Nutzer, nicht in diesem Plan entschieden — **Ausgang:** offen bis Closure (erwartet: weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`).

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

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
