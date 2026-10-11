# Slice slice-v1-abschluss-zeitangaben: Zeitgetreues Einspielen (`--keep-timing`)

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

**Berührte Spec-Stellen:** `LH-FA-21.a` · `LH-FA-17.a` · `LH-FA-20` · `SPEC-034` · `ARC-002` · `ARC-005`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `play --keep-timing` stellt den zeitlichen Abstand der aufgezeichneten Interaktionen her, ohne ihn zu verkürzen: `--timing-mode relative` (Default) oder `absolute`, bei `absolute` mit `--timing-reference connect` (Default) oder `first-request`; dazu Umgebungsvariable und Schlüssel im Abschnitt `play:`. Ungültige Verwendung ist `PGR-E2001`, eine Aufzeichnung, der bei `--keep-timing` irgendeiner Interaktion `offset_ms` fehlt, `PGR-E2003` (Exit-Code 2, Startfehler). Der Slice schließt Abnahmeszenario 13: aufgezeichnet mit `--record-timing`, eingespielt mit `--keep-timing`.

**Übernommen aus `slice-v1-abschluss-einspielen`** (dort §1, Abgrenzung): Zeitangaben und zeitgetreues Einspielen; hier die Hälfte beim Einspielen (das Ziel oben), die Hälfte beim Aufzeichnen liegt in `slice-v1-abschluss-zeitangaben-aufzeichnen`. Die Optionen `--keep-timing`, `--timing-mode` und `--timing-reference` sind bis zu diesem Slice unbekannt (`PGR-E2001`), ihre Schlüssel im Abschnitt `play:` unbekannt (`PGR-E2004`) (dort §6, *Optionen der Folge-Slices*; `slice-v1-abschluss-upstream-verbinden` ließ die Schlüssel `keep_timing` und `timing_mode` bei diesem Slice).

**Übernommen aus `slice-v1-abschluss-sqlite-format-lesen`** (dort §1, Abgrenzung): die Zeitangaben in der SQLite-Datei beim Einspielen. Das Lesen der Spalte `interaction.offset_ms` liefert jener Slice (nur Spalte mit geprüftem Typ); die Prüfung absteigender Werte liegt in `slice-v1-abschluss-zeitangaben-pruefen`, das Schreiben mit `--record-timing` in `slice-v1-abschluss-zeitangaben-aufzeichnen`. Dieser Slice spielt Aufzeichnungen mit Zeitangaben aus beiden Formaten ein.

**Abgegeben beim Schnitt vom 2026-10-11** (je mit der Kennung dieses Slice im Nehmer, `AGENTS.md` §3.13): das Aufzeichnen mit `--record-timing` an `slice-v1-abschluss-zeitangaben-aufzeichnen`, die Prüfung absteigender Werte beim Lesen (der Halbsatz „ein ungültiger Wert ist `PGR-E3003`“ aus DoD-Punkt 1) an `slice-v1-abschluss-zeitangaben-pruefen`, der Uhr-Port mit Fake-Uhr und Systemuhr an `slice-v1-abschluss-zeitangaben-uhr`. Dieser Slice benutzt den Port, wie `-uhr` ihn liefert, und die geprüften Leser.

**Herkunft:** Teil des Schnitts von `slice-v1-abschluss-zeitangaben` vom 2026-10-11 (Entscheidung des Nutzers vom 2026-10-11: vor dem Anspruch zu schneiden, weil der Plan vier Schichten berührte — Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter — und damit über der Grenze von zwei lag, `AGENTS.md` §3.13; den Zuschnitt in vier Teile nach Lieferwert entschied der Planner mit seiner Empfehlung). Die vier Teile: `slice-v1-abschluss-zeitangaben-pruefen` (Lesen prüft die Werte), `slice-v1-abschluss-zeitangaben-uhr` (Uhr-Port und Systemuhr), `slice-v1-abschluss-zeitangaben-aufzeichnen` (`--record-timing`) und `slice-v1-abschluss-zeitangaben` (`--keep-timing`, schließt Abnahmeszenario 13); Handbuch und README zu allen vier liegen gemeinsam im Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` hinter dem letzten Teil (wie `slice-v1-abschluss-sqlite-format-doku` hinter dem Schreiben und dem Lesen: Zwischen den Code-Slices und dem Doku-Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen).


**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `--record-timing` und das Schreiben von `offset_ms` — `slice-v1-abschluss-zeitangaben-aufzeichnen` (dort §1 und DoD): Ein anderer Weg, der Schreiben, Messpunkt und eine Option des Befehls `record` trägt; die Hälfte beim Einspielen braucht ihn nicht, weil Tests Aufzeichnungen mit Zeitangaben von Hand bauen.
- Die Prüfung absteigender Werte beim Lesen — `slice-v1-abschluss-zeitangaben-pruefen` (dort §1 und DoD): Sie gilt auch ohne diesen Slice (`replay`, `play` ohne Option) und liegt in einer anderen Schicht, dem Recording-Adapter.
- Der Uhr-Port, die Fake-Uhr und die Systemuhr — `slice-v1-abschluss-zeitangaben-uhr` (dort §1 und DoD): Dieser Slice ändert die Signatur nicht; verlangt das Warten eine andere, ist das ein Plan-Wechsel jenes Slice.
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` (dort §1 und DoD), hinter dem letzten der vier Teile in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über der Grenze der Schichten.
- Zeitgetreue Antworten im Replay-Modus — Out-of-Scope von LH-FA-21.
- Paralleles Einspielen — Out-of-Scope von LH-FA-20 und LH-FA-21.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen): `--keep-timing` hält die Pausen in beiden Modi (`relative`, `absolute`) und beiden Bezugspunkten (`connect`, `first-request`) ein, ohne sie zu verkürzen (Abnahmeszenario 13: Test mit Fake-Uhr im Service, dazu ein E2E am Binary, aufgezeichnet mit `--record-timing`, eingespielt mit `--keep-timing`, geprüft auf die Untergrenze der Pausen); ohne `--keep-timing` werden vorhandene Zeitangaben ignoriert, und `replay` ignoriert sie immer (Test, Boundary).
- [ ] Optionen `--keep-timing`, `--timing-mode` und `--timing-reference` für `play` mit Umgebungsvariablen und Schlüsseln im Abschnitt `play:` nach der Optionstabelle (`LH-FA-17.a`; Default `false`, `relative`, `connect`; Priorität Option vor Variable vor Datei); bei `record` und `replay` unbekannt; ungültige Verwendung — `--timing-mode` oder `--timing-reference` ohne `--keep-timing`, `--timing-reference` zusammen mit `relative`, ein Wert außerhalb der Aufzählungen — ist `PGR-E2001`, Exit-Code 2, Startfehler; `config show` und `--help` nennen sie (Test am Binary, Negativ je Fall).
- [ ] Eine Aufzeichnung, der bei `--keep-timing` irgendeiner Interaktion `offset_ms` fehlt, wird als ungültige Verwendung abgelehnt: `PGR-E2003`, Exit-Code 2, Startfehler, vor der ersten Anfrage (Test, Negativ).
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
| `internal/hexagon/services` (Play-Service) | update | Warten vor jeder Interaktion über die Uhr des Ports: `relative` und `absolute`, Bezugspunkt `connect` (Beginn des Verbindungsaufbaus zum Server, vor TLS und Anmeldung) oder `first-request`; Prüfung auf `offset_ms` aller Interaktionen vor dem Start (`PGR-E2003`); Abbruch einer Wartezeit durch das Abbruchsignal |
| `internal/hexagon/ports/driving` (Player) | update (nur wenn die Optionen den Port ändern) | Entscheidung des Architect |
| `internal/adapters/driving/cli` | update | Optionen `--keep-timing`, `--timing-mode`, `--timing-reference`, Umgebungsvariablen, Schlüssel im Abschnitt `play:`, Verwendungsfehler `PGR-E2001`, `config show`, Hilfe |
| `internal/bootstrap` | update | Uhr an den Play-Service reichen (reine Verdrahtung, keine Schicht, `AGENTS.md` §3.13) |
| Tests, `test/integration` | neu/update | Service mit Fake-Uhr nach `LH-FA-21` (Happy je Modus und Bezugspunkt, Boundary, Negativ), Optionen am Binary, E2E aufzeichnen → einspielen; `docs/user/abdeckung-*.md` von `make abdeckung` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-uhr` und `slice-v1-abschluss-zeitangaben-aufzeichnen` liegen in `done/`, dazu `slice-v1-abschluss-einspielen`, `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` (Schritt 22 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter `-aufzeichnen`; WIP-Limit 1). Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar, oder die Wartelogik braucht eine dritte Schicht (PGWire-Adapter, Recording-Adapter, Bootstrap mit Logik); Schnitt dann: `relative` mit `connect` hier, `absolute` mit seinen Bezugspunkten in einem eigenen Slice.
- `in-progress` → `open` (blockiert — Carveout?): Die Fake-Uhr bildet die monotone Uhr nicht ab — Carveout.

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

- *Entschieden in `LH-FA-21.a` *Einspielen*:* `relative` wartet vor jeder Interaktion außer der ersten einer Session, bis seit dem Beginn der vorigen mindestens der Unterschied ihrer beiden `offset_ms` vergangen ist; eine Verspätung wird nicht aufgeholt. `absolute` wartet bis zum Bezugspunkt plus aufgezeichnetem Abstand; eine Verspätung wird aufgeholt, der Abstand zweier Interaktionen kann kürzer sein. `connect` ist der Beginn des Verbindungsaufbaus zum Server, vor TLS und Anmeldung (aufgezeichnet gilt `offset_ms` = 0, Aufbau und Anmeldung zählen als Verspätung); `first-request` ist der Beginn der ersten Interaktion der Session (aufgezeichnet gilt deren `offset_ms`). Jede Session beginnt mit eigener Uhr, die Sessions laufen nacheinander.
- *Entschieden ebenda, ungültige Verwendung:* `--timing-mode` oder `--timing-reference` ohne `--keep-timing`, `--timing-reference` zusammen mit `relative` und ein Wert außerhalb der Aufzählungen sind `PGR-E2001`; fehlt irgendeiner Interaktion der Aufzeichnung `offset_ms`, ist die Verwendung `PGR-E2003` (Exit-Code 2, Startfehler). Ohne `--keep-timing` werden Zeitangaben ignoriert, der Replay-Modus ignoriert sie immer.
- *Abbruch einer Wartezeit:* Die Spezifikation bricht eine Wartezeit mit dem Abbruchsignal ab (`LH-FA-20.a`); was `--finish-session-on-interrupt` mit einer Wartezeit tut — offen (Empfehlung des Planners: Die Wartezeit endet, die laufende Interaktion gibt es nicht; die Session endet wie ohne Wartezeit).
- *Wartezeit nach einem Fehler mit `--continue-on-error`:* Der Bezug „Beginn der vorigen Interaktion“ bleibt der der vorigen aufgezeichneten Interaktion — offen.
- *Beginn einer Interaktion:* bei einer Extended-Interaktion der Zeitpunkt, zu dem `play` ihre erste Nachricht sendet — offen (Wortlaut: „seit dem Beginn der vorigen Interaktion“).
- *Sehr großer Wert:* Ein `offset_ms`, der eine Wartezeit von Stunden ergibt, blockiert den Lauf bis zum Signal; eine Obergrenze oder eine Warnung — offen, Entscheidung der Spezifikation.
- *Sessions ohne Interaktion:* keine Prüfung auf `offset_ms`, keine Wartezeit — folgt aus „irgendeiner Interaktion“.
- *Verbindungsaufbau scheitert:* kein Warten, der Aufbau-Fehler gilt wie ohne Option (`LH-FA-20.a`) — offen, ob `connect` bis dahin gemessen wird.
- *Option, Variable, Schlüssel:* Namen, Defaults und Wertemenge stehen in der Optionstabelle und in `LH-FA-17.a` — entschieden; Groß- und Kleinschreibung der Werte wie bei den anderen Aufzählungen (nur Kleinbuchstaben) — entschieden ebenda.

**Risiken:**

- Zeitabstände auf einem ausgelasteten Rechner schwanken; der Test prüft nur die Untergrenze (Pausen nie kürzer), das E2E nur diese — **Ausgang:** offen bis Closure.
- Das Warten hängt am Zeitverhalten des Clients oder Servers im Replay (`BEO-REPO/replay-haengt-am-zeitverhalten-des-clients`, 1×, offen): Der E2E-Test wartet auf Zeit, die nur die Fake-Uhr des Service-Tests deterministisch macht — **Ausgang:** offen bis Closure.

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

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12) — neuer Vertrag, Randformen in §6.
- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (1×, offen) — berührt, weil das Einspielen auf Zeit wartet (Risiko in §6).
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×, offen) — Dieser Plan ist der Rest des Schnitts vom 2026-10-11 mit drei Liefer-Punkten und zwei Schichten; die Größe misst der Architect vor dem Code, und die Rückführung in §4 nennt die Bedingung vorab.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 3. Schichten: zwei (Kern mit dem Play-Service, CLI-Adapter); der Bootstrap reicht nur die Uhr weiter und zählt als reine Verdrahtung nicht (`AGENTS.md` §3.13). Nachgezählt beim Schnitt (2026-10-11): Der Plan stand vorher bei drei Liefer-Punkten und vier Schichten (Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter; die Dokumentation ausgenommen); nach dem Schnitt hat jeder der vier Teile höchstens drei Liefer-Punkte und höchstens zwei Schichten (`-pruefen` 2 und 1, mit Helfer im Kern 2; `-uhr` 3 und 2; `-aufzeichnen` 2 und 2; dieser Slice 3 und 2). Die Nachzählungen früherer Sendungen — Schichtzählung aus `slice-doku-ist-stand` (2026-10-10) und Sendung aus `slice-v1-abschluss-sqlite-format-lesen` (2026-10-11) — galten dem ungeschnittenen Plan; die Sendung aus `slice-v1-abschluss-sqlite-format-lesen` steht jetzt in diesem Plan (§1, Einspielen aus beiden Formaten) und in `-aufzeichnen` (Schreiben, Roundtrip) und `-pruefen`.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
