# Slice slice-v1-abschluss-zeitangaben-uhr: Uhr-Port und Systemuhr

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

**Berührte Spec-Stellen:** `LH-FA-21.a` · `ARC-004` · `ARC-015` · `ARC-001`

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

**Ziel:** Der Kern liest die Zeit und wartet auf einen Zeitpunkt nur über einen Driven Port „Uhr“ ([ADR-0015](../../adr/0015-uhr-port.md)): Der Port, eine Fake-Uhr für Tests, die Systemuhr im Composition Root und die Architekturregel, die den Kern an den Port bindet, liegen vor, bevor ein Slice sie benutzt. Beobachtbar ist der Vertrag an Tests (Warten bis zum Zeitpunkt, Abbruch durch den Kontext, Vorrücken der Fake-Uhr) und am Gate `make a-check`; ein Verhalten des Binaries ändert sich nicht.

**Übernommen aus `slice-v1-abschluss-zeitangaben`** (dort §3, Zeile *Kern, Uhr-Port* und DoD-Punkt 2, die Fake-Uhr): der Port und die Uhren. Die Folgepflicht aus [ADR-0015](../../adr/0015-uhr-port.md) (Architekturregeln halten den Kern frei von Zeitaufrufen der Standardbibliothek) ist hier der dritte Liefer-Punkt; sie war bisher kein Slice zugewiesen.

**Herkunft:** Teil des Schnitts von `slice-v1-abschluss-zeitangaben` vom 2026-10-11 (Entscheidung des Nutzers vom 2026-10-11: vor dem Anspruch zu schneiden, weil der Plan vier Schichten berührte — Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter — und damit über der Grenze von zwei lag, `AGENTS.md` §3.13; den Zuschnitt in vier Teile nach Lieferwert entschied der Planner mit seiner Empfehlung). Die vier Teile: `slice-v1-abschluss-zeitangaben-pruefen` (Lesen prüft die Werte), `slice-v1-abschluss-zeitangaben-uhr` (Uhr-Port und Systemuhr), `slice-v1-abschluss-zeitangaben-aufzeichnen` (`--record-timing`) und `slice-v1-abschluss-zeitangaben` (`--keep-timing`, schließt Abnahmeszenario 13); Handbuch und README zu allen vier liegen gemeinsam im Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` hinter dem letzten Teil (wie `slice-v1-abschluss-sqlite-format-doku` hinter dem Schreiben und dem Lesen: Zwischen den Code-Slices und dem Doku-Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen).

**Lieferwert:** Dieser Slice gibt den beiden Slices, die Zeit brauchen (`slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben`), einen gemeinsamen, getesteten Port, statt ihn in einen von beiden zu legen: Wer ihn in einen legte, brächte dort Kern und Bootstrap mit Logik zu den Schichten der Option dazu — drei statt zwei. Er ist ein Schnitt nach der gemeinsamen Voraussetzung, kein Schicht-Schnitt: Er liefert einen Vertrag mit Gegenprobe, den ein Gate hält.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Messen von `offset_ms` im Record und Warten im Einspielen — `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` (dort §1 und DoD, mit der Kennung dieses Slice): Dieser Slice liefert den Port, nicht seine Benutzer.
- Optionen und Konfiguration — dieselben beiden Slices: Der Port hat keine Option.
- Prüfung der Werte einer Aufzeichnung — `slice-v1-abschluss-zeitangaben-pruefen`: ein anderer Vorgang (Daten, nicht Zeit).
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` (dort §1 und DoD), hinter dem letzten der vier Teile in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über der Grenze der Schichten. Dieser Slice ändert kein Verhalten des Binaries; das Handbuch hat nichts zu beschreiben.
- Eine Zeit, die der Kern aus der Aufzeichnung liest oder in sie schreibt (Zeitstempel in der Aufzeichnung) — Re-Evaluierungs-Trigger von [ADR-0015](../../adr/0015-uhr-port.md), nicht dieser Slice.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Driven Port „Uhr“ im Kern mit zwei Leistungen — die monotone Zeit lesen, bis zu einem Zeitpunkt warten, abbrechbar durch den Kontext — und eine Fake-Uhr, die Zeit nur auf Anweisung vorrückt ([ADR-0015](../../adr/0015-uhr-port.md)); Tests des Vertrags: Warten kehrt beim Erreichen des Zeitpunkts zurück, sofort bei einem Zeitpunkt in der Vergangenheit, mit dem Fehler des Kontexts beim Abbruch (vor und während des Wartens), und die Fake-Uhr weckt alle Wartenden, deren Zeitpunkt erreicht ist (Happy, Boundary, Negativ).
- [ ] Systemuhr im Composition Root (`internal/bootstrap`), die den Port erfüllt: Zeit monoton, Warten abbrechbar; getestet an der Systemuhr selbst (Untergrenze der Wartezeit, Abbruch), nicht am Binary.
- [ ] Architekturregel in `.a-check.yml`, die den Kern an den Port bindet, soweit a-check sie ausdrücken kann, mit Gegenprobe in `make a-check-negativ` (Kern mit Zeitaufruf abgelehnt, Adapter und Bootstrap zugelassen); kann a-check die Regel nicht ausdrücken, entfällt der Punkt mit Begründung in §7 und ein Eintrag im Register (Randform *Regel* in §6).
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
| `internal/hexagon/ports/driven` | neu | Port „Uhr“ (Signatur nach Entscheidung des Architect, §6) |
| `internal/hexagon/...` Testhilfe | neu | Fake-Uhr, Vorrücken nur auf Anweisung; Ort nach Entscheidung des Architect |
| `internal/bootstrap` | update | Systemuhr mit Logik (Warten mit Abbruch); wird hier nicht verdrahtet, nur bereitgestellt |
| `.a-check.yml`, `harness/mk/` Gegenprobe | update | Regel „Kern ruft Zeit nur über den Port“ nach [ADR-0015](../../adr/0015-uhr-port.md), Fitness Function |
| Tests | neu | Vertrag des Ports an der Fake-Uhr und an der Systemuhr, Gegenprobe der Regel |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Voraussetzung: `slice-v1-abschluss-einspielen` ist `done` (Schritt 8 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)); WIP-Limit 1. Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar, oder der Slice braucht eine dritte Schicht (Systemuhr in einem Adapter statt im Bootstrap); die Gegenprobe der Regel ist dann ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): a-check kann die Regel nur mit einer Änderung am Werkzeug ausdrücken; dann zuerst die Entscheidung, ob ein Carveout (Modul 7, [ADR-0015](../../adr/0015-uhr-port.md) Folgepflicht offen) oder ein Slice am Werkzeug sie trägt.

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

- *Einheit und Typ im Port:* Entschieden in [ADR-0015](../../adr/0015-uhr-port.md) ist nur der Zuschnitt (zwei Leistungen); die Signatur nicht. Empfehlung des Planners: ganze Millisekunden seit einem Ursprung der Uhr (kein `time`-Typ im Kern); dann ist „der Kern importiert `time` nicht“ als Import-Regel ausdrückbar. `offset_ms` ist eine ganze Zahl in Millisekunden (`LH-FA-21.a`). Entscheidung des Architect, offen.
- *Zeitpunkt in der Vergangenheit oder jetzt:* Warten kehrt sofort zurück — offen.
- *Abbruch:* Der Kontext endet vor oder während des Wartens; der Aufruf kehrt mit dessen Fehler zurück — entschieden in [ADR-0015](../../adr/0015-uhr-port.md) (abbrechbar durch den Kontext); der Fehler ist der des Kontexts, kein neuer Code.
- *Monotone Zeit:* `LH-FA-21.a` verlangt eine monotone Uhr; eine Wanduhr, die zurückspringt, darf den Port nicht erreichen — entschieden dort.
- *Mehrere Wartende gleichzeitig* (parallele Sessions im Record, ein Abbruch weckt alle): die Fake-Uhr weckt jeden, dessen Zeitpunkt erreicht ist — offen.
- *Regel:* Ob a-check Tests der Kern-Packages ausnimmt (sie importieren `time` heute, `services/*_test.go`) und ob die Regel nach Import oder nach Aufruf greift — offen, Entscheidung des Architect mit der ersten Probe.

**Risiken:**

- a-check kann nach Import, nicht nach Aufruf unterscheiden; die Regel verbietet dann mehr, als [ADR-0015](../../adr/0015-uhr-port.md) sagt, oder weniger — **Ausgang:** offen bis Closure.
- Der Port steht vor seinem ersten Benutzer, und der Benutzer braucht eine andere Signatur (`BEO-REPO/folge-slice-adresse-nimmt-nicht-an`, verkörpert in `AGENTS.md` §3.13) — **Ausgang:** offen bis Closure; Gegenmittel: `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` nennen die Signatur in §1, nicht eine eigene.

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

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10) — der Port ist ein neuer Vertrag, jede Zusage eine Mutation.
- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (1×, offen) — betrifft Zeitverhalten im Replay, nicht den Port; keine Neuberührung.
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×, offen) — dieser Schnitt lässt keine Hälfte an der Grenze: drei kleine Liefer-Punkte, zwei Schichten; die Größe misst der Architect vor dem Code.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 3. Schichten: zwei (Kern mit dem Port, Bootstrap mit der Systemuhr — Logik: Warten mit Abbruch zählt, `AGENTS.md` §3.13); das Gate (`.a-check.yml`) ist ein Liefer-Punkt, keine Produkt-Schicht, wie bei `slice-v1-abschluss-sqlite-format`. Nachgezählt beim Schnitt (2026-10-11). Sitzt die Systemuhr in einem Adapter statt im Bootstrap, sind es zwei Schichten mit Adapter und Kern und derselbe Befund.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
