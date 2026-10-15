# Slice slice-v1-abschluss-verbindungsende-vergleich: Verbindungsende im Antwortvergleich (Abnahmeszenario 16, Satz 3)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0023](../../adr/0023-antwortvergleich-entscheidung.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-24.a` · `LH-FA-20.a` · `LH-FA-17.a` · `SPEC-034` · `SPEC-029` · `ARC-002` · `ARC-005`

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

**Ziel:** `play --compare-responses` kennt das Ende des Servers (`LH-FA-24.a` *Ende der Aufzeichnung*): Zeigt die Aufzeichnung an der Stelle ein Ende und sendet der Server eine `error_response` `FATAL` oder `PANIC` mit demselben SQLSTATE, ist das Ende erwartet — keine Abweichung, keine Meldung, die Session endet, auch ohne `--continue-on-error` läuft die nächste; jedes andere Ende ist die Abweichung `PGR-E5004` (Art Fehler beziehungsweise, bei aufgezeichnetem Ende und normaler Antwort, Nachrichtenart). `play` fügt der Folge `connection_closed` an und wartet nicht auf das Verbindungsende; `verglichen` zählt die Interaktion bis zur Fehlerantwort. Der Hilfetext von `--allow-recorded-errors` nennt das Ende. Abnahmeszenario 16, Satz 3, ist gegen PostgreSQL automatisiert.

**Übernommen aus `slice-v1-abschluss-verbindungsende-einspielen`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): das Ende mit `--compare-responses`, der Hilfetext von `--allow-recorded-errors` und `--compare-responses` im Binary und Abnahmeszenario 16, Satz 3. Jener Slice liefert die Einstufung ohne Vergleich und das Schließen ohne `Terminate` (Port `Einspielziel` und Einspiel-Adapter); dieser setzt ihn voraus und ändert beides nicht.

**Übernommen aus `slice-v1-abschluss-postgres-versionen-zusage`** (dort §6, Randform „Szenarien nur gegen 17“; Sendung des Architect, 2026-10-15, `AGENTS.md` §3.13): `SPEC-029` Punkt 2 nennt das erwartete Verbindungsende beim Einspielen (Abnahmeszenario 16, dritter Satz) heute nicht, weil kein Test es prüft. Im Commit, der den Integrationstest zu Szenario 16 Satz 3 liefert, ergänzt dieser Slice Punkt 2 um den Satz „Das erwartete Verbindungsende beim Einspielen (Abnahmeszenario 16, dritter Satz) ist nur dort geprüft“ (dort = Referenz 17); vorher steht er nicht (`AGENTS.md` §3.11).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Einstufung ohne Vergleich, das Schließen ohne `Terminate` und der Einspiel-Adapter — `slice-v1-abschluss-verbindungsende-einspielen` (Vorgänger): andere Verzweigung des Service und andere Schicht.
- Ein Verbindungsverlust ohne Fehlerantwort bei aufgezeichnetem Ende — bleibt `PGR-E4003`, keine Abweichung (`LH-FA-24.a`); ein anderer Vorgang als das Ende des Servers.
- Aufzeichnen und Wiedergeben des Endes — `-record`, `-upstream`, `-replay`.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13), dort auch die Sätze zu `FATAL`/`PANIC` und `PGR-E5004` aus `slice-v1-abschluss-antwortvergleich-doku`. Die Handbuch-Lücke bis dahin (der Satz zu `FATAL`/`PANIC` für eine handgeschriebene Aufzeichnung mit `connection_closed` stimmt nicht mehr) ist akzeptiert, Entscheidung des Nutzers vom 2026-10-12 (dort §1).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Mit Vergleich ist ein Ende erwartet, wenn die Aufzeichnung an der Position eine `error_response` `FATAL` oder `PANIC` mit `connection_closed` dahinter zeigt und der SQLSTATE (Feld `C`) beider Seiten gleich ist: keine Abweichung, keine Meldung, die Session endet, die nächste läuft ohne `--continue-on-error`. Jedes andere Ende ist `PGR-E5004`, sofort festgestellt — anderer SQLSTATE, gleicher SQLSTATE wie ein aufgezeichnetes `ERROR` ohne Ende, Aufzeichnung ohne Fehlerantwort an der Stelle (aufgezeichneter SQLSTATE leer) —; mit `--continue-on-error` endet dann die Session, die nächste läuft. Zeigt die Aufzeichnung ein Ende und antwortet der Server normal (`ERROR` mit `ReadyForQuery`), vergleicht sich die Folge nach den gewöhnlichen Regeln (Test, einfache Anfrage und Extended-Interaktion).
- [ ] [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen): `play` fügt der Folge des Servers `connection_closed` an und liest nicht über die Fehlerantwort hinaus; `verglichen` zählt die Interaktion bis zur Fehlerantwort, `abweichend` höchstens einmal je Interaktion, `eingespielt` jede gespielte; `--allow-recorded-errors` bleibt mit Vergleich ohne Wirkung, ohne dass die Kombination ein Fehler ist. Der Hilfetext von `--allow-recorded-errors` nennt das aufgezeichnete Ende ohne Vergleich und die Wirkungslosigkeit mit Vergleich, soweit ein Test es prüft (Test im CLI-Adapter).
- [ ] [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Abnahmeszenario 16, Satz 3, gegen PostgreSQL am Binary, dazu der Satz in `SPEC-029` Punkt 2, dass es nur gegen die Referenz 17 geprüft ist (Sendung aus `slice-v1-abschluss-postgres-versionen-zusage`, im Commit des Tests): Eine Aufzeichnung mit Ende wird mit Vergleich gegen den Server eingespielt, der das Ende an derselben Stelle mit demselben Fehlercode herbeiführt (`pg_terminate_backend`, `57P01`): Der Lauf geht mit der nächsten Session weiter und endet mit Exit-Code 0; bei anderem Fehlercode endet er mit `PGR-E5004` und Exit-Code 5; die Zeile zum Ende des Laufs trägt `eingespielt`, `verglichen` und `abweichend` (Integrationstest).
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
| `internal/hexagon/services/play.go`, `play_vergleich.go` | update | Ende im Vergleich: erwartet oder Abweichung, `connection_closed` anfügen, Zähler; `fatal` (heute stets Abweichung) |
| `internal/adapters/driving/cli/cli.go` | update | Hilfetext von `--allow-recorded-errors` (und `--compare-responses`, soweit ein Test es prüft) |
| `internal/hexagon/services`, `internal/adapters/driving/cli` (Tests) | update | Randformen aus §6; die Tests, die die alte Regel prüfen (`TestVergleichFatal`, `TestVergleichFatalRand`, `TestVergleichFatalFortsetzung`, `TestVergleichFatalExtendedGruppen`), werden hier nachgezogen |
| `test/integration` (`play_vergleich_e2e_test.go`, `play_ende_e2e_test.go`) | update | Szenario 16 Satz 3 gegen PostgreSQL; `TestE2EPlayVergleichFatal` prüft die alte Regel und wird nachgezogen |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-einspielen` liegt in `done/` (Voraussetzung: Schließen ohne `Terminate`, Einstufung ohne Vergleich); dieser Slice ist Schritt 39 der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Prüfung des Architect vor dem Code schätzt den Diff mit Tests auf über 1000 Zeilen (Schätzung des Planners 600 bis 1000) — Schnitt vor dem Anspruch: Service | Hilfetext und E2E, Zuschnitt: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`).
- `in-progress` → `next`: Der Einspiel-Adapter oder der Port `Einspielziel` muss sich noch einmal ändern — dritte Schicht, zurück zur Zerlegung (der Vorgänger sollte ihn abschließen).
- `in-progress` → `open`: Der Nutzer ändert die Regel „gleiche Stelle, gleicher SQLSTATE“ — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect; „gleiche Stelle = gleicher SQLSTATE“ und „`--allow-recorded-errors` ohne Vergleich erlaubt ein aufgezeichnetes Ende“ bestätigte der Nutzer am 2026-10-12). Ort: `LH-FA-24.a` *Ende der Aufzeichnung*, *Meldung*, *Randformen des Vergleichs*, `LH-FA-20.a` (Tabelle der Fehlerregeln).

- *Ende bei gleichem SQLSTATE:* erwartet — keine Abweichung, keine Meldung, die Session endet, die nächste läuft auch ohne `--continue-on-error` — `LH-FA-24.a` *Ende der Aufzeichnung*.
- *Ende bei anderem SQLSTATE:* Abweichung `PGR-E5004`, Art Fehler, sofort; ohne `--continue-on-error` Abbruch, mit ihm endet die Session und die nächste läuft — `LH-FA-24.a`, Tabelle der Fehlerregeln.
- *Aufgezeichnet `ERROR` (ohne Ende), der Server sendet `FATAL` mit gleichem SQLSTATE:* Abweichung, Art Fehler — `LH-FA-24.a` *Ende der Aufzeichnung* (zweiter Punkt).
- *Ende ohne Vorbild (die Aufzeichnung hat an der Stelle keine `ErrorResponse`) und ohne `FATAL` in der Aufzeichnung:* Abweichung, aufgezeichneter SQLSTATE leer, Art Fehler — ebenda.
- *Aufgezeichnetes Ende, der Server antwortet normal:* die Folge vergleicht sich nach den gewöhnlichen Regeln, bei `ERROR` mit gleichem SQLSTATE ist die Abweichung die Position des `connection_closed` mit der Art Nachrichtenart — ebenda (dritter Punkt).
- *`FATAL` in einer früheren Interaktion als der aufgezeichneten:* Abweichung dort, die erste Abweichung gilt — `LH-FA-24.a` *Meldung* (je Interaktion die erste).
- *`--continue-on-error`:* setzt nach einer Abweichung fort; ein erwartetes Ende bricht nicht ab — Tabelle der Fehlerregeln.
- *`--allow-recorded-errors` mit Vergleich:* ohne Wirkung, ohne dass die Kombination ein Fehler ist — `LH-FA-20.a` Optionstabelle.
- *Verbindungsverlust ohne Fehlerantwort bei aufgezeichnetem Ende:* `PGR-E4003`, keine Abweichung — `LH-FA-24.a` *Ende der Aufzeichnung*.
- *Zähler:* `verglichen` zählt bis zum `ReadyForQuery` oder bis zur Fehlerantwort `FATAL`/`PANIC`, `abweichend` höchstens einmal je Interaktion, `eingespielt` jede — `LH-FA-24.a` *Meldung*.
- *Extended:* Ende in der letzten Gruppe; der Vergleich der Gruppen davor bleibt, ab der Fehlerantwort gilt die Regel des Endes — `LH-FA-24.a`.

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×; das Vorbild `slice-v1-abschluss-antwortvergleich` maß rund 1450 Zeilen, davon etwa 1100 Tests): die Rückführung in §4 benennt den Schnitt — **Ausgang:** offen bis Closure.
- Szenario 16 Satz 3 stützt sich auf `pg_terminate_backend` und `57P01`; andere Serverversionen sind nicht geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×). **Akzeptiertes Negativ** (Entscheidung des Nutzers vom 2026-10-12): Szenario 16 Satz 3 läuft nur gegen PostgreSQL 17; die Matrix aus `slice-v1-abschluss-postgres-versionen-matrix` fährt es nicht mit (dort §1 als Abgrenzung, `AGENTS.md` §3.13) — **Ausgang:** weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×), `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×) und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×; §6) — der Plan nennt Schätzung und Schnitt vorab (§4). Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13): Liefer-Punkte: 3 (Ende im Vergleich mit Zählern, Hilfetext, E2E von Szenario 16 Satz 3). Schichten: zwei (Kern mit dem Service, `ARC-002`; CLI-Adapter, `ARC-005`, nur der Hilfetext). Schätzung des Planners: 600 bis 1000 Zeilen mit Tests.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-postgres-versionen-zusage` (Architect, 2026-10-15, `AGENTS.md` §3.13): ein Satz in `SPEC-029` Punkt 2 im Commit des Tests von Szenario 16 Satz 3; er gehört zu DoD-Punkt 3 und fügt weder Liefer-Punkt noch Schicht hinzu: Liefer-Punkte weiter 3, Schichten weiter zwei (Kern mit dem Service; CLI-Adapter); Schätzung unverändert, rund zwei Zeilen mehr.
