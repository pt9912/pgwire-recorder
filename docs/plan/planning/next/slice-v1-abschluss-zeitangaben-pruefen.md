# Slice slice-v1-abschluss-zeitangaben-pruefen: Zeitangaben prüfen: absteigende Werte beim Lesen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen)

**Berührte Spec-Stellen:** `LH-FA-21.a` · `SPEC-043` · `SPEC-034`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-11.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Beim Laden einer Aufzeichnung — `replay`, `play`, in YAML wie in SQLite — sind `offset_ms`, die innerhalb einer Session absteigen, `PGR-E3003` (Exit-Code 3): die Aufzeichnung ist beschädigt. Gleiche Werte hintereinander und ein kleinerer Wert in der nächsten Session sind gültig. Das ist beobachtbar an einer Datei von Hand, ohne dass `--record-timing` oder `--keep-timing` existieren.

**Übernommen aus `slice-v1-abschluss-zeitangaben`** (dort bis zum Schnitt DoD-Punkt 1, der letzte Halbsatz „ein ungültiger Wert ist `PGR-E3003` (Test)“): die Prüfung der Werte beim Lesen. Typ und Vorzeichen (`offset_ms` weder `NULL` noch eine ganze Zahl ≥ 0 in SQLite, negativ in YAML) liefern die Leser schon (`slice-v1-abschluss-sqlite-format-lesen`, dort §1: Die Spalte hat einen geprüften Typ; das Verhalten mit Zeitangaben — hier das Absteigen — ist der andere Liefer-Punkt, den jener Slice an `slice-v1-abschluss-zeitangaben` abgab).

**Herkunft:** Teil des Schnitts von `slice-v1-abschluss-zeitangaben` vom 2026-10-11 (Entscheidung des Nutzers vom 2026-10-11: vor dem Anspruch zu schneiden, weil der Plan vier Schichten berührte — Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter — und damit über der Grenze von zwei lag, `AGENTS.md` §3.13; den Zuschnitt in vier Teile nach Lieferwert entschied der Planner mit seiner Empfehlung). Die vier Teile: `slice-v1-abschluss-zeitangaben-pruefen` (Lesen prüft die Werte), `slice-v1-abschluss-zeitangaben-uhr` (Uhr-Port und Systemuhr), `slice-v1-abschluss-zeitangaben-aufzeichnen` (`--record-timing`) und `slice-v1-abschluss-zeitangaben` (`--keep-timing`, schließt Abnahmeszenario 13); Handbuch und README zu allen vier liegen gemeinsam im Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` hinter dem letzten Teil (wie `slice-v1-abschluss-sqlite-format-doku` hinter dem Schreiben und dem Lesen: Zwischen den Code-Slices und dem Doku-Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen).


**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `--record-timing` und das Schreiben von `offset_ms` mit Messpunkt, Uhr und Option — `slice-v1-abschluss-zeitangaben-aufzeichnen` (dort §1 und DoD, mit der Kennung dieses Slice): Dieser Slice prüft, was eine Datei trägt, nicht wie sie entsteht; das Schreiben setzt die Prüfung voraus, damit seine Roundtrip-Probe die geprüften Leser benutzt.
- `--keep-timing` und das Warten beim Einspielen — `slice-v1-abschluss-zeitangaben` (dort §1 und DoD): Das Warten setzt geprüfte, nicht absteigende Werte voraus; die Prüfung liegt hier, weil sie auch ohne Option gilt (§6, Randform *Prüfung ohne Wunsch*).
- Der Uhr-Port und die Systemuhr — `slice-v1-abschluss-zeitangaben-uhr`: Die Prüfung vergleicht Werte der Datei und braucht keine Uhr.
- Fehlt `offset_ms` einer Interaktion, ist das hier kein Fehler — Bestand und Zusage der Spezifikation: Nur zeitgetreues Einspielen verlangt die Werte (`PGR-E2003`, `slice-v1-abschluss-zeitangaben`).
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` (dort §1 und DoD), hinter dem letzten der vier Teile in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über der Grenze der Schichten.
- Produkt-Code im Kern, im PGWire-Adapter, im CLI-Adapter und im Bootstrap — Schicht-Abgrenzung: Die Prüfung liegt in den Lesern des Recording-Adapters (YAML und SQLite) und, wenn beide denselben Vergleich brauchen, in einem Helfer des Domain Model; mehr Schichten berührt der Slice nicht (Rückführung in §4).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen): Eine Aufzeichnung, in der `offset_ms` innerhalb einer Session absteigen, lädt in YAML und in SQLite nicht: `PGR-E3003`, Exit-Code 3, die Meldung nennt Session und Interaktion (Test je Format, Negativ); gleiche Werte hintereinander und ein kleinerer Wert am Anfang der nächsten Session laden (Test je Format, Boundary).
- [ ] Die Prüfung gilt beim Laden, nicht erst mit einer Option: `replay` und `play` auf einer solchen Datei enden beim Start mit `PGR-E3003` (Test am Binary, je ein Aufruf); `make abdeckung` ist gefahren, `make abdeckung-check` grün.
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
| `internal/adapters/driven/recording` (`yaml.go`, `sqlite_lesen_baum.go`) | update | Vergleich der `offset_ms` je Session nach dem Laden der Interaktionen; Meldung ohne Wert aus der Datei |
| `internal/hexagon/model` | update (nur wenn beide Leser denselben Vergleich brauchen) | ein Helfer statt zwei Kopien; zählt als zweite Schicht (§8) |
| Tests im Recording-Adapter und `test/integration` | neu/update | Negativ und Boundary je Format nach `LH-FA-21`; `replay` und `play` am Binary |
| `docs/user/abdeckung-*.md` | update | von `make abdeckung` erzeugt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-sqlite-format-lesen` liegt in `done/` (Voraussetzung: das Lesen der SQLite-Spalte; Schritt 17 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md); WIP-Limit 1). Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar, oder die Prüfung braucht eine dritte Schicht (Kern, PGWire-Adapter, CLI oder Bootstrap mit Logik); Schnitt dann: YAML und SQLite getrennt.
- `in-progress` → `open` (blockiert — Carveout?): Die Spezifikation lässt die Randform *Prüfung ohne Wunsch* (§6) offen, und der Nutzer hat sie nicht entschieden; dann zuerst die Entscheidung des Nutzers.

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

**Randformen** (`AGENTS.md` §3.12) — Liste für den Architect, geschrieben beim Schnitt vom 2026-10-11. Entschieden ist eine Randform erst, wenn die Spezifikation (Technik-Stratum) sie vor dem ersten Code-Commit trägt; die Planung hat die Spezifikation nicht geändert. Je Punkt steht, wo er entschieden ist, sonst *offen*: Der Architect entscheidet vor dem Code oder gibt die Frage an den Nutzer (Lastenheft-Bezug). Der Implementer entscheidet keine offene Randform (`.claude/commands/implement-slice.md`, Randform-Rückgabe).

- *Absteigend:* Werte, die innerhalb einer Session absteigen, machen die Aufzeichnung beschädigt (`PGR-E3003`) — entschieden in `LH-FA-21.a` *Aufzeichnen*; der Vergleich gilt je Session, nicht über die Sessions hinweg — ebenda („innerhalb einer Session“).
- *Gleicher Wert:* zwei Interaktionen mit demselben `offset_ms` — offen (der Wortlaut „absteigen“ lässt gleiche Werte zu; Empfehlung des Planners: gültig, nur ein kleinerer Wert als der vorige ist absteigend).
- *Mischform:* eine Interaktion ohne `offset_ms` zwischen zwei mit Wert — offen (Empfehlung des Planners: Lesen erlaubt, verglichen werden die vorhandenen Werte in der Reihenfolge der Interaktionen; ohne Wert wird nichts verlangt, das tut nur `--keep-timing`, `PGR-E2003`).
- *Prüfung ohne Wunsch:* „beschädigt“ (Laden) und „ohne `--keep-timing` ignoriert, im Replay immer“ (Verwenden) stehen in `LH-FA-21.a` nebeneinander; offen, ob die Prüfung auch ohne `--keep-timing` und im Replay greift (Empfehlung des Planners: ja, beim Laden — wie die Prüfung negativer Werte im Bestand, `yaml.go`; „ignoriert“ meint das Warten). Entscheidung der Spezifikation, bei Lastenheft-Bezug des Nutzers.
- *Form des Werts in YAML:* Fließkomma (`1.5`, `1.0`), Zeichenkette, Hex, Exponent, `null`, Wert über dem Bereich von `int64` — offen; ohne Entscheidung entscheidet die YAML-Bibliothek still.
- *Meldung:* nennt Session und Interaktion, nie einen Wert aus der Datei und keine Zeilennummer (Regel für `sqlite`, `SPEC-043`; für `yaml` ohne Regel, `BEO-REPO/yaml-meldung-nennt-wert-und-zeile`) — die eigene Prüfung dieses Slice hält die Regel in beiden Formaten; mehrere Fehler in einer Datei: welcher gemeldet wird — offen.
- *Eine Session mit einer oder ohne Interaktion:* nichts zu vergleichen, gültig.

**Risiken:**

- Zwei Leser, zwei Kopien des Vergleichs: eine bricht, die andere nicht (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, verkörpert in `AGENTS.md` §3.10) — **Ausgang:** offen bis Closure; Gegenmittel: Mutation je Format und je Randform, vom Implementer selbst rot gesehen.
- Die Entscheidung zur Prüfung ohne Wunsch kommt nicht vor dem Code (`BEO-REPO/spec-randform-erst-im-review-entschieden`, verkörpert in `AGENTS.md` §3.12) — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`, gelesen am Stand des Schnitts, 2026-10-11). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10) — die Prüfung ist ein neuer Vertrag, jede Zusage eine Mutation (Risiko in §6).
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12) — Randformen stehen in §6.
- `BEO-REPO/yaml-meldung-nennt-wert-und-zeile` (2×, offen) — die Meldung der Prüfung nennt keinen Wert (Randform *Meldung*); sie berührt den Eintrag nicht neu.
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×, offen) und `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) — dieser Schnitt lässt keine Hälfte an der Grenze: zwei Liefer-Punkte, eine Schicht, kleiner Diff; die Größe misst der Architect vor dem Code, und ein Befund dort wäre ein weiterer Beleg.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten: eine (Recording-Adapter); zwei, wenn ein gemeinsamer Helfer im Domain Model liegt — höchstens zwei. Nachgezählt beim Schnitt (2026-10-11, `AGENTS.md` §3.13); die Sendung aus `slice-v1-abschluss-zeitangaben` (DoD-Halbsatz) steht in DoD-Punkt 1.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
