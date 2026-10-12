# Slice slice-v1-abschluss-antwortvergleich-doku: Handbuch und README zu: Antwortvergleich beim Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-antwortvergleich`)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-10.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-antwortvergleich` liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-antwortvergleich`, abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-antwortvergleich` über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt jenen Slice voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt dahinter. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit beiden.

**Übernommen aus `slice-v1-abschluss-umgebung-anderes-kommando`** (dort §7, offene Frage zu V-3; Planner, 2026-10-12, `AGENTS.md` §3.13): Das Handbuch nennt `--compare-responses` erst mit `slice-v1-abschluss-antwortvergleich` und nennt dann zu jeder Option ihre Umgebungsvariable mit den Kommandos, bei denen sie gilt; die Variable beim anderen Kommando ist `PGR-E2001` (Satz zur Meldung `PGR-E2001` im Handbuch, §5 *Einstellungen*). Das Handbuch beschreibt die Variable, nachdem das gebaute Binary sie liest.

**Gebunden durch `slice-v1-abschluss-antwortvergleich`** (dort §6 *Randformen* und §8; Sendung des Architect, 2026-10-12, `AGENTS.md` §3.13, §3.11): Das Handbuch trägt erst mit diesem Slice, was das Binary nach jenem Slice zeigt, und nur das, was ein Test jenes Slice prüft: `--compare-responses` bei `play` (Variable `PGWIRE_RECORDER_COMPARE_RESPONSES`, Schlüssel `compare_responses` unter `play:`); verglichen werden Art und Reihenfolge der Antworten, Spaltenbeschreibung (Anzahl, Name, Typ), Beschreibung der Parameter, der Befehl ohne abschließende Zahlen, der SQLSTATE eines Fehlers und der Transaktionsstatus, nicht Zeilenwerte, Zeilenzahlen, Hinweise und Einstellungen des Servers; eine Abweichung ist `PGR-E5004` mit Exit-Code 5 und nennt Session, Sequenznummer, Position und Art (bei einem Fehler die SQLSTATE beider Seiten, nie den Text), je Interaktion die erste; ein Serverfehler ohne Vorbild und eine Fehlerantwort `FATAL` oder `PANIC` sind `PGR-E5004`, ein Netzwerkabbruch ohne Fehlerantwort ist `PGR-E4003`; mit Vergleich ersetzt dieser die Fehlerregel (`PGR-E4004` entfällt), `--allow-recorded-errors` bleibt ohne Wirkung, `--continue-on-error` setzt nach einer Abweichung fort; die Zeile zum Ende des Laufs trägt `eingespielt`, `verglichen` und `abweichend` auf der Stufe `info`; ohne Option findet kein Vergleich statt. Welche dieser Sätze ins Handbuch kommen, bestimmt die Probe am gebauten Binary; was sie nicht zeigt, bleibt draußen. Zum Handbuch-Teil in DoD-Punkt 1; kein neuer Liefer-Punkt, Schichten weiter eine.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-antwortvergleich` liefert.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.
- Das Verhalten bei einem aufgezeichneten Ende des Servers (`FATAL`/`PANIC` an derselben Stelle mit demselben SQLSTATE ist erwartet, nicht `PGR-E5004`) — `slice-v1-abschluss-verbindungsende-doku` (dort §1, *Übernommen aus `slice-v1-abschluss-antwortvergleich-doku`*; Entscheidung des Nutzers vom 2026-10-12, Option B, [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md); Sendung des Planners, `AGENTS.md` §3.13): Der Satz in „Gebunden durch“, `FATAL` und `PANIC` seien `PGR-E5004`, beschreibt den Ist-Zustand, den dieser Slice nach `slice-v1-abschluss-antwortvergleich` vorfindet; `slice-v1-abschluss-verbindungsende-vergleich` ändert das Binary dahinter, und der Doku-Slice dazu ersetzt den Satz. Dieser Slice schreibt ihn so, wie das Binary heute zeigt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [ ] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-antwortvergleich` liefert, im Ist-Zustand des gebauten Binaries (zu `--compare-responses` mit Umgebungsvariable und dem Verhalten beim anderen Kommando, aus `slice-v1-abschluss-umgebung-anderes-kommando`; Probe am gebauten Binary): Die Abschnitte, die Kommandos, Optionen, Meldungscodes, Beispiele und Grenzen dieses Slice betreffen, stehen im Handbuch, und die Grenz-Sätze, die er überholt, sind ersetzt. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [ ] `README.md` nennt, was `slice-v1-abschluss-antwortvergleich` liefert, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-antwortvergleich` liefert: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt |
| `README.md` | update | Sätze, die `slice-v1-abschluss-antwortvergleich` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5; die Variablen-Zeile lautet `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`, und mit ihr ist die Ausgabe leer (Beleg und Mutanten: `slice-v1-abschluss-einspielen-anmeldung-doku` §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-antwortvergleich` liegt in `done/` (Voraussetzung; Schritt 24 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar; Schnitt dann: Handbuch und README getrennt.
- `in-progress` → `open` (blockiert — Carveout?): Eine Aussage lässt sich nur mit einer Änderung am Binary wahr machen, und der Nutzer will sie nicht streichen; dann zuerst die Entscheidung des Nutzers, welcher Slice das Verhalten liefert.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, die Proben und `grep`-Läufe in §7, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect prüft vor dem ersten Commit, ob `slice-v1-abschluss-antwortvergleich` eine weitere Randform liefert, die dort nicht steht.

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** offen bis Closure.
- Der Slice liefert seinen Teil von Handbuch und README nicht vollständig, weil nur die DoD ihn trägt und das Review ihn übersieht (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** offen bis Closure.

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

**Stand.** Geändert sind `docs/user/benutzerhandbuch.md` (Version 0.5: neuer Abschnitt *Antworten beim Einspielen vergleichen*; Einleitung, Vorspann der *Hinweise* beim Einspielen, zwei Sätze zum Exit-Code nach einem Signal, Satz und Tabellenzeile in §5, Exit-Code 5, Zeilen `PGR-E2004` und `PGR-E4004`, neue Zeile `PGR-E5004`, Glossar), `README.md` und dieser Plan (§7, §8). Produkt-Code, Tests, Gates, Hilfetexte und `docs/user/abdeckung-*.md` sind unverändert; `make abdeckung-check` endet mit Exit 0. Commits `ef0c823` (Handbuch, README) und `6f10a41` (Nachzug des Satzes zum Signal im Abschnitt zu den Zeitangaben).

**Probe-Aufbau.** `make build` am Stand `6c4d53c` (Image `pgwire-recorder:dev`, `sha256:e6c5ef3178…`). Proben in einem Scratch-Verzeichnis außerhalb des Repos; Docker-Netz `impl-av-net` mit dem gepinnten PostgreSQL-Image aus `harness/mk/integration.mk` (`postgres:17-alpine@sha256:b0f9560a…`, Anmeldung `trust`, Netzalias `postgres`, damit die Beispiele des Handbuchs wörtlich laufen), je Probe eine Datenbank: `probe` (`t(n int, s text)`), `probe2` (andere Werte, drei Zeilen), `d_name` (Spalte `label`), `d_type` (`n bigint`), `d_cnt` (Spalte `x`), `dup` (Zeile mit dem Schlüssel, den die Aufzeichnung einfügt). Aufzeichnungen entstanden mit `record` und `psql` (einfach, `\bind` für das erweiterte Protokoll). Aufzeichnungen mit einer Abweichung, die kein Server auslöst, sind aus ihnen von Hand geändert (`sed`/`awk`: Tag, SQLSTATE, `tx_status`, Reihenfolge, SQL, entfernte `notice_response`/`parameter_status`/Server-Parameter) und in der Tabelle als *editiert* gekennzeichnet; `FATAL` entstand durch `pg_terminate_backend` aus einer zweiten Sitzung während `pg_sleep(8)`, der Verbindungsabbruch durch `docker kill` des Servers in `pg_sleep(8)`, das Signal durch `docker kill -s TERM`. Container, Netz und Scratch-Dateien sind danach entfernt.

**Probe-Tabelle** (Handbuch = `docs/user/benutzerhandbuch.md`; `Z` = die Zeile `play beendet`; jeder Aufruf mit `play --upstream … --input …`, wo nichts anderes steht):

| Nr. | Aussage (Handbuch/README) | Probe | Ergebnis |
|---|---|---|---|
| 1 | Beispiel *Vorgehen* (Option), wörtlich | `--upstream postgres:5432 --input ./recordings/users.yaml --compare-responses` | Exit 0, `Z`: `eingespielt=3 verglichen=3 abweichend=0` |
| 2 | Variable `PGWIRE_RECORDER_COMPARE_RESPONSES=true` genügt | dasselbe ohne Option, Variable gesetzt | Exit 0, `verglichen=3` |
| 3 | Schlüssel `compare_responses: true` im Abschnitt `play:` genügt, Beispieldatei | `--config ./compare.yaml` (Datei wie im Handbuch) | Exit 0, `verglichen=3`; `config show` nennt den Schlüssel |
| 4 | Abweichung: `PGR-E5004`, Exit 5, die beiden Zeilen des Beispiels | `p2` gegen `d_cnt` | Exit 5, `level=ERROR … code=PGR-E5004 … Position 1, Art Spaltenbeschreibung`, `Z`: `eingespielt=1 verglichen=1 abweichend=1` |
| 5 | Spalten: Anzahl, Name, Typ verglichen | `SELECT * FROM t` gegen `d_cnt` (Anzahl), `d_name` (Name), `d_type` (Typ) | je Exit 5, Art `Spaltenbeschreibung`; gegen `probe` Exit 0 |
| 6 | Parameter verglichen (erweitertes Protokoll) | editiert: Describe auf `statement`, `parameter_description [23]`; gegen `probe` und `d_type` (`bigint`) | `probe` Exit 0; `d_type` Exit 5, Position 3, Art `Spaltenbeschreibung` |
| 7 | Befehl ohne abschließende Zahlen: `UPDATE 1` gleich `UPDATE 7`, `UPDATE 1` nicht gleich `DELETE 1` | editiert: Tags `UPDATE 7`, `SELECT 99`; Tag `DELETE 1` | Exit 0 (`verglichen=3`); Exit 5, `Interaktion 2`, Position 1, Art `Befehl` |
| 8 | Fehler: SQLSTATE verglichen, Text nicht | `SELECT 1/0` gegen `probe`, Text der Aufzeichnung editiert; SQL zu `SELECT * FROM nichtda` editiert | Exit 0; Exit 5, `Position 1, Art Fehler (SQLSTATE aufgezeichnet „22012“, Server „42P01“)`, kein Text |
| 9 | Transaktionsstatus verglichen | editiert: `tx_status` von `BEGIN` auf `I` | Exit 5, Position 2, Art `Transaktionsstatus` |
| 10 | Art und Reihenfolge der Antworten (Art `Nachrichtenart`) | editiert: `UPDATE …` durch `SELECT 1` ersetzt; erweitert: `parse_complete`/`bind_complete` vertauscht | je Exit 5, Position 1, Art `Nachrichtenart` |
| 11 | Nicht verglichen: Werte und Zahl der Zeilen | `p1` gegen `probe2` (drei Zeilen, andere Werte), einfach und erweitert (`\bind`) | je Exit 0, `verglichen` gezählt |
| 12 | Nicht verglichen: Hinweise (`NOTICE`) | `SELECT f()` (Funktion mit `RAISE NOTICE`) gegen `probe2` (ohne), und umgekehrt, einfach; erweitert gegen `probe2` | Exit 0 |
| 13 | Nicht verglichen: Meldungen über Einstellungen | `SET application_name`, `parameter_status` aus der Aufzeichnung entfernt, einfach und erweitert; Server-Parameter der Datei editiert | Exit 0 |
| 14 | Position zählt ab 1 ohne Zeilen, Hinweise, Einstellungen | editiert: Tag der Interaktion mit zwei Zeilen; mit `NOTICE`; mit `parameter_status` und geändertem `tx_status` | Position 2 in allen drei Fällen (Art `Befehl`, `Befehl`, `Transaktionsstatus`) |
| 15 | Meldung nennt Sitzung und Interaktion; erweitert: Nummer der Aufzeichnung | Probe 7 (`Interaktion 2`), erweitert `Interaktion 1` | wie in den Zeilen |
| 16 | Je Interaktion die erste Abweichung | editiert: Position 1 (Spalten) und Tag zugleich verschieden, zwei solche Interaktionen mit `--continue-on-error` | zwei Zeilen `PGR-E5004`, je `Position 1`; `abweichend=2` |
| 17 | Ohne `--continue-on-error` endet `play` nach der Interaktion mit der Abweichung | Probe 4 (zwei Interaktionen) | `eingespielt=1` |
| 18 | Mit `--continue-on-error` weiter, jede abweichende gemeldet, Exit 5 | zwei abweichende Interaktionen, auch mit Variablen statt Optionen | zwei Zeilen, Exit 5, `eingespielt=2 abweichend=2` |
| 19 | Fehler ohne Vorbild: `PGR-E5004` (Art `Nachrichtenart`), nicht `PGR-E4004`; ohne Vergleich `PGR-E4004` | `INSERT` auf vorhandenen Schlüssel, einfach und erweitert, mit und ohne Option, mit `--continue-on-error` | mit Option Exit 5, Position 1 (einfach) und 4 (erweitert); ohne Option Exit 4, `PGR-E4004`; mit beiden Exit 5 |
| 20 | Fehler mit demselben SQLSTATE: keine Abweichung, nächste Interaktion läuft, ohne `--continue-on-error` und ohne `--allow-recorded-errors` | `SELECT 1/0; SELECT 2;` | Exit 0, `eingespielt=2 verglichen=2 abweichend=0`; erweitert Exit 0; ohne Option derselbe Lauf `PGR-E4004`, Exit 4 |
| 21 | Anderer SQLSTATE: Art `Fehler` mit beiden SQLSTATE | Probe 8, erweitert (editiert) | Exit 5, `aufgezeichnet „42P01“, Server „22012“` |
| 22 | Mit Vergleich ist `--allow-recorded-errors` ohne Wirkung | Fehler ohne Vorbild, anderer SQLSTATE, gleicher SQLSTATE, je mit der Option; ohne Vergleich mit der Option und anderem SQLSTATE | Exit 5, Exit 5, Exit 0 (wie ohne Option); ohne Vergleich Exit 0 |
| 23 | `FATAL`: mit Vergleich `PGR-E5004`, auch bei gleichem SQLSTATE; Art `Fehler` an der Position der Fehlerantwort; eine frühere Abweichung geht vor; ohne Vergleich `PGR-E4003` | `pg_sleep(8)`, Server beendet die Sitzung (`57P01`): Aufzeichnung ohne Fehler (`row_description` dann `FATAL`), mit aufgezeichnetem `ERROR 57P01` an Position 2 (editiert), mit aufgezeichnetem Fehler an Position 1 (frühere Abweichung); beides mit `--allow-recorded-errors` und `--continue-on-error`; ohne Vergleich | Exit 5, Position 2, Art `Fehler` (aufgezeichnet leer bzw. `57P01`); mit Fehler an Position 1 `Position 1, Art Nachrichtenart`; ohne Vergleich Exit 4, `PGR-E4003` |
| 24 | `PANIC` wie `FATAL` | kein Server löst `PANIC` aus: `TestVergleichFatal` (Fall `V PANIC`, `S PANIC ohne V`, einfach und erweitert) in `internal/hexagon/services/play_vergleich_test.go` | Test grün im Lauf von `make gates` |
| 25 | Verbindungsverlust ohne Fehlerantwort: `PGR-E4003`, Exit 4, auch nach einer Abweichung, beide Meldungen im Log | `docker kill` des Servers in `pg_sleep(8)`, ohne und nach einer Abweichung (`--continue-on-error`) | Exit 4; eine Zeile `PGR-E4003`; nach Abweichung zusätzlich `PGR-E5004`; `Z`: `verglichen=0` bzw. `eingespielt=2 verglichen=1 abweichend=1` |
| 26 | Signal: laufende Interaktion mit Vergleich beendet, Abweichung darin `PGR-E5004`, Exit 5 auch bei früherer Abweichung | `SIGTERM` in `pg_sleep(8)`: ohne Abweichung; Abweichung in der laufenden Interaktion; frühere Abweichung (`--continue-on-error`); in der Wartezeit von `--keep-timing` nach Abweichung | Exit 0 mit `verglichen=1`; Exit 5, Art `Befehl`; Exit 5; Exit 5 mit `eingespielt=1` |
| 27 | Zähler: `eingespielt`, `verglichen`, `abweichend`; unterbrochene Interaktion zählt nicht als verglichen | Proben 1, 4, 25 | wie dort |
| 28 | Ohne `--compare-responses`: kein Vergleich, `verglichen`/`abweichend` 0 | `p1` und `p2` gegen `probe2`/`d_cnt` ohne Option; Option `=false` gegen Schlüssel `true`; Variable `false` gegen Schlüssel `true` | je Exit 0, `verglichen=0` |
| 29 | `Z` hat die Stufe `info`, mit `--log-level warn` fehlt sie | `--log-level warn`, mit und ohne Abweichung | mit Abweichung nur die Zeile `PGR-E5004`, ohne Abweichung keine Ausgabe |
| 30 | `record`, `replay`: Option `PGR-E2001`, Variable `PGR-E2001`, Schlüssel `PGR-E2004` | je Kommando | `flag provided but not defined`, Exit 2; `gehört zu einer Option, die record/replay nicht kennt`, Exit 2; `record.compare_responses: unbekannter Schlüssel`, `replay.…`, Exit 2 |
| 31 | Anderer Wert als `true`/`false`: Variable `PGR-E2001`, Schlüssel `PGR-E2004` | `…=vielleicht`; `compare_responses: yes` | Exit 2, `erlaubt sind true und false`; Exit 2, `play.compare_responses: erlaubt sind true und false` |
| 32 | Vorspann der *Hinweise* beim Einspielen: die drei Hinweise gelten ohne Vergleich | Proben 19, 20, 22, 23 gegen die Läufe ohne Option (`PGR-E4004`, Exit 4; `--allow-recorded-errors` entschuldigt jeden SQLSTATE, Exit 0; `PGR-E4003`) | wie dort |
| 33 | README: Sätze zu `play` | Proben 4, 5, 7, 8, 9, 10, 11, 19 | wie dort |
| 34 | Tabelle §5: Betriebsart `play`, Standard `false` | `--help` der vier Kommandos | `play --help` nennt `--compare-responses`; `grep -c compare` bei `record`, `replay`, `config show`: 0, 0, 0 |
| 35 | Exit-Code 5, Zeilen `PGR-E2004`, `PGR-E4004`, `PGR-E5004`, Glossar | Proben 4, 31, 19 | wie dort |

**Stand der Gates.** `make gates` am Stand `485aa94` (Arbeitsbaum sauber): Exit 0 (Laufzeit 5 min 9 s). Der Commit danach ändert nur diese Zeile; `make gates` auf dem Stand nach dem letzten Commit ist die Messung, die der Bericht nennt.

**Neue und geänderte Sätze neben Bestand.** Der Vorspann der *Hinweise* beim Einspielen grenzt die drei Sätze zu Fehlerantworten (`PGR-E4004`, `--continue-on-error`, `--allow-recorded-errors`, `FATAL`) gegen den Vergleich ab (Proben 19, 20, 22, 23, 32). Die beiden Sätze „Der Exit-Code ist 0, nach einem früheren Fehler 4“ (Einspielen, Zeitangaben) nennen jetzt die Abweichung (5; Probe 26). Die Zeile des Exit-Codes 5 und die Beispiele in §5 nennen den Vergleich (Proben 30, 4). Unverändert bleiben die Sätze zu TLS, Anmeldung, COPY und `--finish-session-on-interrupt`: Der Vergleich ändert nichts an ihnen, und keine Probe dieses Slice liest sie neu.

**Ungenannt geblieben** (keine Probe am gebauten Binary oder nur Spezifikation): Verhalten je Serverversion (kein Ausgleich); Tag mit leerem Wort (`SELECT ` gegen `SELECT`); Fehlerantwort ohne Feld `C`; Gruppengrenzen beim erweiterten Protokoll; die folgende Sitzung nach `FATAL` mit `--continue-on-error`; `--keep-timing` neben `--compare-responses` (nur der Exit-Code nach Signal, Probe 26); Wert der Option außer `true`/`false`; die Rangfolge Option, Variable, Schlüssel für diese Option (Probe 28 zeigt sie, der Satz steht allgemein in §5 und bleibt dort); Lauf ohne eine Interaktion. Das Handbuch sagt nichts zu einem aufgezeichneten Verbindungsende (§1; Entscheidung des Nutzers vom 2026-10-12): Das Binary zeigt heute `FATAL`/`PANIC` mit Vergleich immer als `PGR-E5004`, und so steht es im Handbuch.

**Befund an den Planner.** (a) Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen, Codes) ist am Stand `6c4d53c` ohne meine Änderung rot: Die Variablen-Zeile mit der benannten Ausnahme `PGWIRE_RECORDER_PASSWORD` nennt `PGWIRE_RECORDER_BOGUS`, eine Variable in §5 des Handbuchs (Satz zu unbeachteten Variablen, eingeführt mit `slice-v1-abschluss-umgebung-anderes-kommando`, Commit `b91131c`), die zu keiner Option gehört. §3 dieses Plans verlangt „mit ihr ist die Ausgabe leer“; das stimmt nicht. Ich habe den Satz nicht geändert (er ist Bestand einer früheren Entscheidung und steht nicht in meinem Diff); die Ausnahme in `implement-slice.md` Schritt 17 und in §3 bräuchte `PGWIRE_RECORDER_BOGUS` als zweite benannte Ausnahme, oder der Satz im Handbuch nennt die Variable nicht. (b) Die Roadmap führt in Zeile 23 „In Arbeit: nichts“; ich habe sie nicht angefasst (Planung).

**Befehlsfolge aus `slice-doku-ist-stand` §7** (Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD`), am Stand `6f10a41`:

| Prüfung | Ausgabe |
|---|---|
| Optionen des Handbuchs, die `--help` nicht kennt | `--h`, `--help`, `--version` (Nicht-Optionen im Handbuch), `--name`, `--rm` (`docker run`) |
| Variablen ohne Option, ohne `PGWIRE_RECORDER_PASSWORD` | `PGWIRE_RECORDER_BOGUS` (Befund (a)) |
| Codes außerhalb des Katalogs | leer |

`--compare-responses` und `PGWIRE_RECORDER_COMPARE_RESPONSES` stehen nicht in der Ausgabe; `PGR-E5004` ist im Katalog. Die Befehlsfolge ist nicht geändert, kein Wächter eingeführt; ihre Mutanten stehen in `done/slice-v1-abschluss-einspielen-anmeldung-doku` §7 und werden nicht wiederholt. `grep -n -i -E "spec/|lastenheft|slice|welle|review|LH-|SPEC-|ADR|noch nicht|kommt|geplant|künftig"` über das Handbuch findet nur Bestand ohne Bezug (`noch nicht` in §3 „Zieldatei existiert noch nicht“, `kommt` in der Ruhefrist von `record` und in §4 zu TLS); `make docs-check`: 0 Befunde.

**Mutation (§3.10)** entfällt: Der Slice liefert keinen neuen Vertrag, nur Beschreibung; die Probe-Tabelle ersetzt sie (rot wäre ein Satz ohne bestandene Probe).

**Kommentare (Schritt 20)** entfallen: kein Code, kein Skript und keine Konfiguration geändert.

**Rückgaben an den Architect:** keine; keine Randform entschieden. **Nachzählen:** Liefer-Punkte 2, Schichten 1 (Nutzerdokumentation), wie in §8.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) — genau der Gegenstand dieses Slice: Das Handbuch sagt nur zu, was das Binary zeigt.
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-antwortvergleich`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code. Nachgezählt beim Eintragen der Übernahme aus `slice-v1-abschluss-umgebung-anderes-kommando` (Planner, 2026-10-12): Sie gehört zu DoD-Punkt 1 (Handbuch), kein neuer Liefer-Punkt; Liefer-Punkte weiter 2, Schichten weiter eine.

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-antwortvergleich` (Architect, 2026-10-12): Sie gehört zu DoD-Punkt 1 (Handbuch); Liefer-Punkte weiter 2, Schichten weiter eine.

Nachgezählt beim Eintragen der Abgrenzung zum aufgezeichneten Ende des Servers (Planner, 2026-10-12, `AGENTS.md` §3.13): Sie nimmt dem Slice nichts und fügt nichts hinzu; Liefer-Punkte weiter 2, Schichten weiter eine.

Nachgezählt beim Umsetzen (Implementer, 2026-10-12): Der Diff ändert `docs/user/benutzerhandbuch.md`, `README.md` und diesen Plan, keinen Code; Liefer-Punkte 2, Schichten 1 (Nutzerdokumentation; Planung zählt nicht).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
