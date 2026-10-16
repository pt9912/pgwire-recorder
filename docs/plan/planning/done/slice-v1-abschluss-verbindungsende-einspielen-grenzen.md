# Slice slice-v1-abschluss-verbindungsende-einspielen-grenzen: Verbindungsende beim Einspielen ohne Vergleich (die Grenzen der Einstufung am Binary)

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

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-20.a`

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

**Ziel:** Die Grenzen der Einstufung eines Verbindungsendes beim Einspielen ohne `--compare-responses` sind am Binary gegen PostgreSQL belegt (`LH-FA-20.a` *Interaktion*, `LH-FA-25.a` *Form*): Ohne `--allow-recorded-errors` ist die Fehlerantwort `FATAL` des Servers auf eine einfache Anfrage `PGR-E4003` mit Exit-Code 4, auch wenn die Aufzeichnung ein Ende trägt; mit der Option, aber mit einer Aufzeichnung ohne Ende an der Stelle (aufgezeichnet `ERROR`, der Server sendet `FATAL`), ist es ebenfalls `PGR-E4003` mit Exit-Code 4. Der Slice liefert zwei Fälle am Binary mit je einer handgeschriebenen Aufzeichnung (Gegenstand: Einspielen einer Aufzeichnung mit Ende gegen einen Server; das Ende bewirkt der Server, `play` löst es nicht aus) und keinen Produktcode: Das Verhalten liefert `slice-v1-abschluss-verbindungsende-einspielen`; die Fälle belegen es am Binary, und ihre Mutationen färben sie rot.

**Geschnitten aus `slice-v1-abschluss-verbindungsende-einspielen-extended`** (Planner, 2026-10-16, beim Anspruch, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap): Die Schätzung des Planners beim Anspruch lag mit den Maßen des Gebers bei 490 bis 815, Mitte 650 (65 %), am Geber kalibriert 732 (73 %), das obere Ende über der Meldeschwelle 800; die beiden Fälle prüfen Verhalten der einfachen Anfrage, das das Binary seit Schritt 38a erfüllt, und hängen an keinem Code der Extended-Interaktion. Einzeln lieferbar: Dieser Slice setzt nur `slice-v1-abschluss-verbindungsende-einspielen` (`done/`) voraus, `-einspielen-extended` wartet nicht auf ihn und er nicht auf `-einspielen-extended`.

**Übernommen aus `slice-v1-abschluss-verbindungsende-einspielen-extended`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-16): die Fälle am Binary, in denen das Ende nicht erwartet ist (ohne Option, aufgezeichnet ohne Ende an der Stelle). Ursprung: `slice-v1-abschluss-verbindungsende-einspielen` (dort §1) und `slice-v1-abschluss-verbindungsende-format` (dort §1).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Einstufung selbst, das Schließen ohne `Terminate`, Port und Einspiel-Adapter — `slice-v1-abschluss-verbindungsende-einspielen` (Geber, `done/`; Klasse 4): geliefert; ein Befund der Fälle an dieser Stelle ist die Rückführung in §4, nicht eine Änderung hier (Schicht-Abgrenzung: dieser Slice ändert keinen Produktcode).
- Die Extended-Interaktion mit Ende (Gruppe mit `Sync` oder `Flush`), Gruppen, Antworten nach einer `Flush`-Gruppe, Sendefehler zwischen Gruppen, und die Fälle am Binary dazu — `slice-v1-abschluss-verbindungsende-einspielen-extended` (Klasse 1): andere Verzweigung des Service und Code.
- Das Ende mit `--compare-responses`, der Hilfetext von `--allow-recorded-errors` und `--compare-responses` — `slice-v1-abschluss-verbindungsende-vergleich` (Klasse 1): andere Verzweigung des Service (`play_vergleich.go`) und der CLI-Adapter als zweite Schicht.
- Die Kombination mit `--continue-on-error` am Binary — bleibt bei den Fällen des Service im Geber (`TestPlayEndeNichtErwartet`): sie ist dort je Zusage mit Mutation gefangen; ein dritter Fall am Binary belegte keine weitere Zusage des Binaries (Bestand bleibt bewusst stehen, Klasse 2).
- Das Aufzeichnen und Wiedergeben des Endes, ein Verbindungsende ohne Fehlerantwort — wie in `-einspielen-extended` §1 (`LH-FA-25`, Out-of-Scope; Klasse 3).
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13; Klasse 1).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): Gegen PostgreSQL am Binary: Eine handgeschriebene Aufzeichnung einer einfachen Anfrage mit Ende, deren Ausführung den Backend beendet (`FATAL`, SQLSTATE `57P01`), endet **ohne** `--allow-recorded-errors` mit `PGR-E4003` und Exit-Code 4; eine Aufzeichnung ohne Ende an der Stelle (aufgezeichnet `ERROR`, der Server sendet `FATAL`) endet **mit** der Option mit `PGR-E4003` und Exit-Code 4 (Integrationstest, eine Aufzeichnung als eigene Datei je Fall).
- [x] Je Zusage des Vertrags eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
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
| `test/integration` (`play_ende_grenzen_e2e_test.go`, `testdata`, eine Datei je Fall) | neu | Binary gegen PostgreSQL mit handgeschriebenen Aufzeichnungen: einfache Anfrage ohne Option, aufgezeichnet ohne Ende an der Stelle; Aufzeichnungen unter `testdata/` je Fall eine: `play-ende-grenze-ohne-option.yaml`, `play-ende-grenze-ohne-ende.yaml` (Namen des Architect, 2026-10-16; die Fixtures des Gebers heißen `play-ende-erwartet.yaml` und werden nicht geteilt); die Hilfen des Gebers in `play_ende_e2e_test.go` werden benutzt |
| `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabelle aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-einspielen` liegt in `done/` (Voraussetzung: Einstufung, Port und Adapter); dieser Slice ist Schritt 38c der Reihenfolge in §5 von welle-v1-abschluss und braucht `-einspielen-extended` nicht. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis:** `fa149c68`, Stand des Hauptzweigs beim Anspruch, danach nicht mehr geändert; die Commits danach tragen `Verantwortlich:`, den Move und die Roadmap ein und ändern keinen Pfad der Metrik); **über 400** geht der Slice zurück zur Zerlegung. **Meldeschwelle 300:** Erreicht die Summe 300, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik (gemessen im Geber: E2E-Fall mit Testdaten 65): **140 bis 190, Mitte 165 und damit 41 % der Marke** (zwei Fälle zu 65 bis 80, Tabelle 10 bis 30).
- `in-progress` → `next`: Ein Fall am Binary verlangt eine Änderung an Produktcode (Service, Port, Adapter, Bootstrap) — das Verhalten gehört `slice-v1-abschluss-verbindungsende-einspielen`; zurück zur Zerlegung (Architect und Planner), nicht still im Code.
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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code. Ort: `LH-FA-20.a` (Tabelle der Fehlerregeln, *Interaktion*), `LH-FA-25.a` *Form*. Die Randformen der einfachen Anfrage stehen in `slice-v1-abschluss-verbindungsende-einspielen` §6 und gelten hier unverändert; der Slice fügt keinen Vertrag hinzu, er belegt zwei Zusagen des Gebers am Binary. Der Architect prüft vor dem Code, ob jede Zeile entschieden ist.

- *Gegenstand und Fixtures:* Das Ende bewirkt der Server, `play` löst es nicht aus. Jeder E2E-Fall baut seine Aufzeichnung selbst (eine Datei je Fall, kein gemeinsamer Zustand), damit eine Mutation nur den Fall färbt, den sie bricht; Fixtures `testdata/play-ende-grenze-ohne-option.yaml` (trägt das Ende an der Stelle; der Aufruf ohne Option) und `testdata/play-ende-grenze-ohne-ende.yaml` (trägt `ERROR` an der Stelle; der Aufruf mit Option), Tests in `play_ende_grenzen_e2e_test.go`; jeder Fall prüft Exit-Code 4 und `code=PGR-E4003` in der Ausgabe, mit der Stelle der Interaktion (Session und Interaktion) wie die Meldung sie nennt. Die Zeilen je Fall mit Testdaten meldet der Implementer in §7 (Messung).
- *Stelle der Nachricht (`AGENTS.md` §3.12), einfache Anfrage:* (a) *Aufgezeichnetes `connection_closed`:* letzter Eintrag der Antworten der Interaktion; das Format lehnt jede andere Stelle beim Laden ab (`PGR-E3003`). (b) *Fall ohne Option:* die Aufzeichnung trägt das Ende an der Stelle; die Option fehlt; `FATAL` des Servers ist nicht erwartet, `PGR-E4003`, Exit 4, auch wenn das Ende aufgezeichnet ist — `LH-FA-20.a` *Interaktion*. (c) *Fall ohne Ende an der Stelle:* aufgezeichnet `ERROR`, der Server sendet `FATAL`, die Option ist gesetzt; `PGR-E4003`, Exit 4 — `LH-FA-20.a` *Interaktion*. (d) *Reihenfolge der Gründe:* ein Grund (das Ende ist nicht erwartet), keine Reihenfolge.
- *Zähler und Nachricht:* keine Zusage hier; `PGR-E4003` und Exit 4 sind die Zusage, die Zähler der Meldung sind nicht Gegenstand (`LH-FA-24.a` *Meldung* gilt für den Vergleich).
- *Fälle am Binary:* Ein Server, der den Backend während der Ausführung beendet, sendet `FATAL` `57P01`; die Aufzeichnung ohne Ende an der Stelle trägt ein `ERROR` an der Stelle und kein Ende. Die Version des Servers ist die Referenz 17 (Risiko unten).
- *Mutationen:* Zusage 1 (ohne Option nicht erwartet): Mutant `play.go` behandelt das Ende als erwartet, auch ohne `--allow-recorded-errors` → roter Fall `ohne-option`. Zusage 2 (ohne Ende an der Stelle nicht erwartet): Mutant erwartet die `FATAL`, sobald die Option gesetzt ist, ohne das aufgezeichnete Ende zu prüfen → roter Fall `ohne-ende`. Jeder Mutant färbt nur den Fall seiner Zusage (der erste ändert nichts, wo die Option gesetzt und das Ende nicht aufgezeichnet ist, der zweite nichts, wo die Option fehlt); der Implementer prüft das. Die E2E-Mutanten liegen im Arbeitsbaum, werden vor dem Fahren von `make test-integration` committed (das Image baut aus dem Arbeitsbaum, der Stand der Gates hängt am Hash) und danach per `git checkout` zurückgenommen. Die Mutanten sind Wegwerf-Änderungen, nicht Teil der Lieferung.

**Risiken:**

- Die Fälle sind bei Beginn grün (das Binary erfüllt sie), ein grüner Test belegt ohne Mutation nichts — **Ausgang:** *entfallen*: Mutanten 1 bis 6 (§7 und Verifikation) rot gesehen, die beiden Zusagen der Fälle färben je nur ihren Fall (Isolation von Review und Verifikation nachgefahren).
- Handgeschriebene Aufzeichnungen bilden nur ab, was die Spezifikation sagt (`BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`, 3×); die Fälle belegen das Einspielen gegen PostgreSQL, nicht das Aufzeichnen (`slice-v1-abschluss-verbindungsende-upstream-extended`) — **Ausgang:** *weiter offen* → `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (Zähler 5, Stand `geplant` bleibt; keine neue Adresse).
- Das Verhalten des Servers (`57P01`) ist nur gegen die Referenzversion 17 geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 3× bei Beginn) — **Ausgang:** *weiter offen* → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (Zähler 4).
- Ein Befund verlangt Produktcode (zweite Rückführung in §4) — **Ausgang:** *entfallen*: `git diff fa149c68..HEAD -- internal cmd` ist leer, Review und Verifikation bestätigen es, keine Rückführung.

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

### Belege des Implementers

**Stand:** Basis `fa149c68`, Lieferung `670ec850`; Server `postgres:17-alpine` (Referenz). Kein Produktcode geändert (`git diff fa149c68..HEAD -- internal` leer); kein Befund an `play.go`, keine Rückführung, keine Randform außerhalb von §6.

**Läufe:** `make test-integration` mit der Lieferung grün (beide Fälle `PASS`); `make abdeckung` schrieb 4 Zeilen in `docs/user/abdeckung-e2e.md`; `make gates` ein Lauf, Exit 0 (darin `abdeckung-check`, `a-check`, `lint`, `test-integration`, `test-postgres-versionen`, Gegenproben), Nachweis geschrieben.

**Mutationen** (Wegwerf im Arbeitsbaum, nach jedem Lauf `git checkout`, `touch` vor dem Lauf; Stellen in `internal/hexagon/services/play.go`; je ein voller `make test-integration`, Exit 2):

| Zusage | Mutation | roter Test |
|---|---|---|
| Ohne `--allow-recorded-errors` ist das Ende nicht erwartet, auch wenn es aufgezeichnet ist (`PGR-E4003`, Exit 4) | Zeile 261 `l.ende = traegtEnde(in)` (Option nicht geprüft) | `TestE2EPlayEndeGrenzeOhneOption`; `ohne-ende` bleibt grün (Isolation) |
| Mit Option, aber ohne aufgezeichnetes Ende an der Stelle ist das Ende nicht erwartet | Zeile 261 `l.ende = s.optionen.AllowRecordedErrors` (Aufzeichnung nicht geprüft) | `TestE2EPlayEndeGrenzeOhneEnde` (Ausgabe: Exit 0, `play beendet eingespielt=2`); `ohne-option` bleibt grün (Isolation) |
| Die Meldung nennt die Stelle (Session und Interaktion) | Zeile 203 ohne `Interaktion %d` | beide Grenzfälle rot (dazu die Bestandstests der Fehlerstelle, Folge der Mutation an einer gemeinsamen Meldung) |
| Die Meldung trägt `PGR-E4003` | Zeile 417 `CodeServerError` statt `CodeConnectionLost` | beide Grenzfälle rot (dazu `TestE2EPlayFortsetzungAbbruch`) |
| Exit-Code 4 | `internal/hexagon/model/fehler.go`, `exitCode`: `PGR-E4003` liefert 5 (Wurzel der Zuordnung, `Error.ExitCode`) | beide Grenzfälle rot (Ausgabe: „Exit-Code 5, erwartet 4“); dazu `TestE2ERecordExtendedAbbruch`, `TestE2ERecordExtendedSigtermNachBlockade`, `TestE2EPlayFortsetzungAbbruch` |

Alle Zusagen tragen eine gesehene Mutation (Exit 4 nachgeholt auf Anweisung des Koordinators, Mutant im Arbeitsbaum, `touch`, danach `git checkout`; `git diff fa149c68..HEAD -- internal` ist leer, daher kein neuer voller `make gates`).

**Messung** (Metrik §4, `git diff --numstat fa149c68..670ec850 -- internal test docs/user`): Summe **92** (Marke 400, Meldeschwelle 300 nicht erreicht; Schätzung 140 bis 190, darunter): `play_ende_grenzen_e2e_test.go` 52, zwei Fixtures je 18, Abdeckungstabelle 4. Je Fall mit Testdaten 46 (92 durch zwei Fälle; die Hilfe `spieleEndeNichtErwartet` und der Kopf der Testdatei, rund 30 Zeilen, sind geteilt). Gegenüber 65 bis 70 der Geber niedriger, weil die Fälle eine gemeinsame Hilfe `spieleEndeNichtErwartet` teilen.

**Beobachtungen:** Der Plan schätzte 140 bis 190; gemessen 92, die Hilfe über zwei Fälle senkt den Satz je Fall (Beleg für `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` in der Gegenrichtung, Planner entscheidet über Aufnahme). Die E2E-Mutanten kosten einen vollen `make test-integration` (rund 2,5 Minuten) je Mutant.

### Closure des Planners (2026-10-16)

**Gate-Nachweis.** Voller `make gates` ein Lauf grün, Exit 0 (Stand der Lieferung `670ec850`, Belege oben). Die Commits danach ändern nur Pfade der Positivliste (`docs/plan/planning/`, `docs/reviews/`, `docs/user/`); `git diff fa149c68..HEAD -- internal cmd` ist leer; die Verifikation fuhr `make gates-schnell` grün. `make kopf-check` und `make docs-check` laufen nach der Closure nacheinander.

**Gegenstand:** geliefert (kein Fall von „übernommen“ oder „entfallen“): zwei Fälle am Binary gegen PostgreSQL, kein Produktcode.

**Review und Verifikation.** Review (`docs/reviews/2026-10-16-review-slice-v1-abschluss-verbindungsende-einspielen-grenzen.md`): 0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO (F-1 bis F-4); fünf Mutanten der Tabelle selbst nachgefahren. Verifikation (`docs/reviews/2026-10-16-verify-slice-v1-abschluss-verbindungsende-einspielen-grenzen.md`): bestanden, V-1 bis V-4 INFO; sieben Mutanten gefahren (Nr. 1 bis 6 rot am erwarteten Grund und isoliert, Nr. 7 grün). Ausgang je Befund: F-1/V-2 akzeptiertes Negativ (SQLSTATE und Meldung der Fehlerantwort sind keine Zusage dieses Slice, §6 *Zähler und Nachricht*; den Wortlaut trägt der Test des Service, von der Verifikation gelesen, nicht nachgefahren), Register `BEO-REPO/e2e-meldung-bindung-weicht-von-zusage-ab` (neu, 1×); F-2 akzeptiert, derselbe Eintrag (Gegenrichtung); F-3/V-3 akzeptiert, die Zeile sagt zu, was der Test prüft (kein Zuwachs bei `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`); F-4/V-1 kein Befund am Diff, Frage an den Nutzer unten, Register `BEO-REPO/spezifikation-weicht-vom-lastenheft-ab` (neu, 1×); V-4 Datenpunkt unten.

**Größe.** 92 in der Metrik von §4 (Messung oben), 23 % der Marke 400, Meldeschwelle 300 nicht erreicht; Schätzung des Planners 140 bis 190, Mitte 165: **unter der Untergrenze** (66 % von 140, 56 % der Mitte), die Obergrenze ist nicht überschritten. Nach der Zählregel vom 2026-10-16 (Obergrenze überschritten = Auftreten) kein Auftreten bei `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Vermerk als Gegenrichtung: Die zwei Fälle teilen die Hilfe `spieleEndeNichtErwartet` und den Kopf (rund 30 Zeilen), der E2E-Fall mit Testdaten maß 46 statt 65 bis 70 (Geber) und 70,5 (`-einspielen-extended`). Faktor Tests zu Produktcode entfällt (kein Produktcode). Kein Auftreten bei `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`: Die Rückführung trat nicht ein.

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge; die Zeilen stehen in §6 am Risiko selbst):

| Risiko | Ausgang | Folge-Slice bzw. Begründung |
|---|---|---|
| Grüner Test ohne Mutation | entfallen | Mutanten 1 bis 6 rot gesehen (Implementer, Review, Verifikation); die zwei Zusagen der Fälle färben je nur ihren Fall; Mutant 7 (SQLSTATE und Meldung) ist keine Zusage des Slice |
| Handgeschriebene Aufzeichnungen ohne Beleg gegen den Server | weiter offen | Register `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (Zähler 5, Stand `geplant`: `slice-v1-abschluss-verbindungsende-upstream` für die einfache Hälfte, `-upstream-extended` für die Extended-Hälfte) |
| `57P01` nur gegen die Referenz 17 | weiter offen | Register `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (Zähler 4) |
| Ein Befund verlangt Produktcode | entfallen | `git diff fa149c68..HEAD -- internal cmd` leer; Rückführung nicht eingetreten |

Ob die Ausgänge tragen: *Entfallen* (Mutation) trägt an der Verifikation, die die Mutanten selbst nachbaute, *entfallen* (Produktcode) am Diff, nicht am Bericht. *Weiter offen* zu den handgeschriebenen Aufzeichnungen trägt, weil die Adresse besteht: Die einfache Hälfte liegt bei `slice-v1-abschluss-verbindungsende-upstream` (§1 *Übernommen aus* `-replay`, Fixtures `ende-fatal.yaml`, `ende-lebend.yaml`); ob deren Form die der beiden Fixtures dieses Slice trägt (`FATAL` mit `connection_closed`; `ERROR` mit `ready_for_query`), urteilt die Verifikation jenes Slice, nicht dieser Eintrag. Der Eintrag steht bei 5× mit Zuweisung, die Schwelle 3× ist überschritten und durch die Zuweisung getragen.

**Sendungen und Bindungen.** Der Slice weist keinem Nehmer etwas zu. Die Übernahme aus `slice-v1-abschluss-verbindungsende-einspielen-extended` (§1) ist mit der Lieferung getragen. Keine Nachzählung bei einem anderen Slice.

**Register** (Zähler = Dateien unter `evidence/`; Vermerke zählen nicht):

| Eintrag | Zähler | Fund | Stand |
|---|---|---|---|
| `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` | 5 (vorher 4) | Risiko aus §6, *weiter offen* | geplant (`slice-v1-abschluss-verbindungsende-upstream-extended`) |
| `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` | 4 (vorher 3) | Risiko aus §6, *weiter offen* | offen, ohne Verkörperung (Frage unten) |
| `BEO-REPO/e2e-meldung-bindung-weicht-von-zusage-ab` | 1 (neu) | F-1, F-2, V-2: Meldung als Teilstrings des stderr gebunden, SQLSTATE und Meldung nicht, Wortlaut des Codes gebunden | offen |
| `BEO-REPO/spezifikation-weicht-vom-lastenheft-ab` | 1 (neu) | F-4, V-1: `LH-FA-25.a` *Negative* gegen `LH-FA-20.a` *Interaktion* | offen, Frage unten |
| `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` | 3 | Vermerk: V-1 zählt hier nicht mit (zwei Straten, nicht zwei Stellen); Frage an den Nutzer | offen |
| `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` | 9 | kein Auftreten: 92 unter der Untergrenze 140; Gegenrichtung als Vermerk | geplant, Vermerk |
| `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` | 2 | kein Auftreten: Rückführung nicht eingetreten, 23 % der Marke | offen, Vermerk |
| `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | 67 | kein Auftreten: F-3 ist Lesbarkeit der Abdeckungs-Zeile | verkörpert, Vermerk |
| `BEO-REPO/test-prueft-teilstring-gegen-meldung-mit-zufallspfad` | 1 | kein Auftreten: kein Pfad in den Teilstrings | offen, Vermerk |

*Einordnung der Befunde (Urteil des Planners):* F-1 und V-2 sind eine eigene Klasse (Bindung der Meldung am Binary), nicht `test-prueft-teilstring-gegen-meldung-mit-zufallspfad` (kein Pfad) und nicht `abwesenheitszusage-an-wortlaut-der-meldung-geprueft` (keine Abwesenheitszusage). Kein Zuwachs bei `BEO-REPO/plan-folgt-korrektur-nicht`: §1, §3, §4, §6 und §7 folgen dem Diff (Review: Negativbefund zu §3.9). Kein Zuwachs bei `BEO-REPO/liste-gruener-mutanten-unvollstaendig`: Mutant 7 ist keine Mutation einer Zusage des Slice, sondern die Sonde der Verifikation. Kein Eintrag für `BEO-REPO/meldeschwelle-ohne-artefakt`: Die Schwelle wurde nicht erreicht.

**Lerneintrag — benannte Spec-Lücke und eine Beobachtung, keine geänderte Regel.** Geändert wird keine Regel, kein `liegt in` (`AGENTS.md` und `.claude/` ändert der Planner nicht). (1) *Benannte Spec-Lücke: Lastenheft und Spezifikation sagen für ein Ende mit anderem SQLSTATE bei `--allow-recorded-errors` Verschiedenes (F-4, V-1).* `LH-FA-25.a` *Negative* nennt für ein Ende, das die Aufzeichnung „mit einem anderen Fehlercode“ zeigt, einen Verbindungsfehler; `LH-FA-20.a` *Interaktion* lässt das Ende mit der Option erwartet sein, „gleich welcher SQLSTATE“. Das Binary folgt der Spezifikation; beide Fixtures tragen `57P01` und lassen die Frage offen. Das Lastenheft rangiert höher; die Abweichung entscheidet der Nutzer (Frage unten). Gefunden hat es das Review beim Vergleich der beiden Texte, nicht ein Test: Ein Slice, der „die Grenzen“ am Binary belegt, belegt nur die Grenzen, die die Spezifikation zieht. (2) *Was funktionierte:* Die Mutationstabelle nannte je Mutant, welchen der beiden Fälle er färbt und welchen nicht; Review und Verifikation prüften die Isolation, nicht nur das Rotwerden, und fanden sie bestätigt. Die Verifikation baute zusätzlich Mutant 3 (Fehlerantwort erwartet, sobald die Aufzeichnung eine trägt), der am selben Fall fällt, und Mutant 7, der grün bleibt und damit die Grenze des Slice zeigt. (3) *Beobachtung zur Schätzung:* Der Satz je E2E-Fall hängt an den geteilten Hilfen; siehe den Vermerk bei `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`. (4) *Neuer Sensor: keiner.*

**Fragen an den Nutzer.** (a) *Lastenheft-Abweichung (V-1/F-4).* Das Lastenheft ist `Draft`; vor `Accepted` ist es ohne Change Request änderbar, der Planner ändert es nicht. Drei Wege: (i) Lastenheft angleichen: In `LH-FA-25.a` *Negative* nennt der Satz zum Ende „mit einem anderen Fehlercode“ den Fall ohne Vergleich nicht als Verbindungsfehler, wenn der Anwender `--allow-recorded-errors` gesetzt hat (Empfehlung des Planners: Die Option sagt ausdrücklich, dass aufgezeichnete Fehler als erwartet gelten; der Wortlaut der Spezifikation „gleich welcher SQLSTATE“ ist dazu die Lesart, die das Binary seit Schritt 38a zeigt); (ii) Spezifikation und Binary angleichen: mit der Option ist ein Ende mit anderem SQLSTATE `PGR-E4003`; das ist eine Verhaltensänderung des Service, ein eigener Slice mit Mutation und Fall am Binary, und die Fixtures dieses Slice wären um einen dritten Fall zu ergänzen; (iii) beides stehen lassen und die Lesart in `LH-FA-25.a` klarstellen. Eine Entscheidung steht aus; bis dahin bleibt der Eintrag `BEO-REPO/spezifikation-weicht-vom-lastenheft-ab` offen. (b) Zählt V-1 bei `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` mit (dann 4×, kein Strata-Unterschied)? Der Planner zählt es nicht mit. (c) Die Zählung von `docs/user/` (Abdeckungstabelle) als Schicht steht weiter aus: Dieser Slice zählt eine Schicht, konservativ zwei (§8).

**Schätzung des nächsten Schritts, `slice-v1-abschluss-verbindungsende-vergleich` (Schritt 39):** Der Plan nennt als Rückführung die Marke 1000 in der Metrik von `AGENTS.md` §3.14 (Summe aus Einfügungen und Löschungen über `internal test docs/user`), Schätzung des Planners im Plan 600 bis 1000; eine Meldeschwelle und eine Basis stehen dort nicht (die Basis trägt der Planner beim Anspruch ein). Gerechnet mit den Maßen der Geber: Produktcode (Ende im Vergleich in `play.go` und `play_vergleich.go`, Hilfetext, drei Kommentare) 50 bis 90; Tests im Service (acht Randformen aus §6, zwei Fälle der letzten `Flush`-Gruppe, die vier nachzuziehenden Bestandstests) 300 bis 560; CLI-Test 20 bis 40; E2E (zwei Fälle zu 46 bis 70, dazu `TestE2EPlayVergleichFatal` nachziehen) 110 bis 180; Tabellen (`abdeckung-unit.md`, `abdeckung-e2e.md`) 20 bis 40; Zeilen für die Parametrisierung der Bestandstests unter dem Lint-Profil `dupl` (Geber: 100 in `-einspielen`, Nacharbeit +89 in `-einspielen-extended`) 0 bis 100. Summe **570 bis 1010, Mitte rund 790 (79 % der Marke)**; das obere Ende erreicht die Marke. Die 15 bis 25 Zeilen des Architect für den Einleitungssatz von `LH-FA-24.a` stehen in `spec/` und zählen in der Metrik nicht. Drei Liefer-Punkte, zwei Schichten (Kern, CLI-Adapter; `docs/user` konservativ als dritte markiert). **Ob zu schneiden ist:** Nach Zahl der Liefer-Punkte und Schichten ist der Plan nicht zu groß; nach der Spanne liegt der Mittelwert unter der Marke, das obere Ende bei ihr. Empfehlung des Planners: Nicht jetzt schneiden, die Prüfung des Architect vor dem Code entscheidet beim Anspruch (wie bei 38a, 38b); schätzt er über 900, ist die im Plan genannte Linie *Service | Hilfetext und E2E* zu ziehen (der erste Teil liefert das Verhalten gegen den Fake, der zweite belegt es am Binary mit Szenario 16 Satz 3, wie 38c es für die Grenzen tat), und der Planner trägt dann Meldeschwelle 800 und Basis ein. Die Entscheidung dazu fällt beim Anspruch.

**Paarungen.** (a) Anker: Der Lerneintrag trägt kein Feld `liegt in`, keine Verkörperung; nichts zu prüfen. (b) Folge-Slice: Kein Risiko trägt den Ausgang *eingetreten*; genannt sind `slice-v1-abschluss-verbindungsende-upstream` und `slice-v1-abschluss-verbindungsende-upstream-extended` als bestehende Adressen des Registers (`geplant`), beide in `next/`, nicht geschlossen. (c) Register: Alle genannten Kennungen `BEO-REPO/…` bestehen als Verzeichnis, jedes mit nicht leerem `evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (3×), `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×; §6) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (Anlass dieses Schnitts: geschnitten beim Anspruch, nicht in `in-progress/` zurückgeführt). Keine weiteren Treffer. Der Zähler-Stand ist beim Anspruch aus dem Register zu lesen.

Nachgezählt beim Anlegen (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Schnitt von `slice-v1-abschluss-verbindungsende-einspielen-extended`): Liefer-Punkte: 1 (zwei Fälle am Binary mit Aufzeichnungen). Schichten: eine (Integrationstests, `test/integration`); **konservativ gezählt mit `docs/user` (Abdeckungstabelle) als zweiter Schicht: zwei** — markiert, bis der Nutzer entscheidet; kein Produktcode, Bootstrap und CLI-Adapter unberührt. Die Sendung (zwei Fälle) liegt unter 3 Liefer-Punkten und 2 Schichten. Schätzung in der Metrik von §4: E2E zwei Fälle zu 65 bis 80 (130 bis 160), Tabelle 10 bis 30: **140 bis 190, Mitte 165, 41 % der Marke 400**; Meldeschwelle 300.

Geprüft beim Anspruch (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14): Maße des Gebers `-einspielen` (E2E-Fall 65, Test 37 und Aufzeichnung 28) und des Nehmers `-einspielen-extended` (70,5 je Fall, Aufzeichnung 43, geteilter Kopf 36) gegen den Bestand `play_ende_e2e_test.go` gelesen: Die Hilfen (`endeDatei`, `leereDatenbank`, `starteBisFrist`, `zaehler`, `wert`) liegen im Paket und werden nicht neu geschrieben, der Kopf mit zwei Einbettungen kostet rund 15 Zeilen, zwei Aufzeichnungen mit einer Session und ein bis zwei Interaktionen 20 bis 30 je Stück. Schätzung **bestätigt: 140 bis 190, Mitte 165 (41 % der Marke 400)**, das untere Ende rechnerisch bei rund 110, das obere unter der Meldeschwelle 300. Zählung bestätigt: Liefer-Punkte 1, Schichten 1, konservativ 2 mit `docs/user`; die Nutzerentscheidung dazu steht aus (Regel-Slice `slice-harness-regeln-nach-v1-abschluss`). Mutierbarkeit gegen `play.go` gelesen (`l.ende = AllowRecordedErrors && traegtEnde(in)`): Der Mutant „Ende immer erwartet“ (`traegtEnde(in)` allein) färbt nur `ohne-option`, „Option allein genügt“ (`AllowRecordedErrors` allein) nur `ohne-ende`. Die Nachzählung im Nehmer `-vergleich` (§4: Metrik der Marke) berührt dessen Größe nicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
