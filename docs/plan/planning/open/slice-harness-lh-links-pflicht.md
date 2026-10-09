# Slice slice-harness-lh-links-pflicht: Linkpflicht für Lastenheft- und MR-Kennungen im Doku-Gate

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er
von der nächsten Welle-Closure. Angelegt nach Entscheidung des Nutzers vom 2026-10-09 zu
`slice-harness-d-check-v0-85`; eingeplant nach welle-v1-abschluss, nicht vorgezogen (§4
*Start*).

**Bezug:** [MR-000](../../../../harness/conventions.md#mr-000--baseline-aussage) (Baseline-Aussage mit dem ID-Schema, dessen Lastenheft- und Adaptions-Kennungen das Doku-Gate als Link verlangen soll). Keine Anforderung des Lastenhefts im Scope: Der Slice schärft das Gate der Dokumentation, nicht das Produkt.

**Berührte Spec-Stellen:** `spezifikation.md §11` (Harness-Werkzeuge: Dort entscheidet der Architect vor dem Code die Randformen der Linkpflicht aus §6, nach `AGENTS.md` §3.12 auch für Harness-Werkzeuge; die Kennung der neuen Stelle vergibt er beim Schreiben und trägt sie hier nach)

**Verantwortlich:** —

**Autor:** pt9912 (Planner). **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make docs-check` meldet eine Lastenheft-Kennung (Präfix `LH-` mit Klasse und
Nummer) und eine Kennung des Adaptions-Blocks (Präfix `MR-` mit drei Ziffern) ohne Link als
Befund (`id-unlinked`), der lebende Bestand trägt beide als Links (ADRs und Nutzer- und
Wartungs-Doku vorab durch `slice-harness-lh-links-bestand`), und der Block *Strenges
Doc-Gate* der Workflow-Commands sagt genau dieses Verhalten zu.

**Sendung aus `slice-harness-d-check-v0-85`** (`AGENTS.md` §3.13; Entscheidung des Nutzers
vom 2026-10-09 *DoD berichtigen + eigener Slice*):

- **Teil-Zusage „bares `LH-`-Token ist rot“** aus DoD-Punkt 1 des Gebers. Dort war sie vor dem
  Code nicht erfüllbar: §1 des Gebers lässt `.d-check.yml` unverändert, und das `ids`-Muster
  für Lastenheft-Kennungen steht dort nur als Kommentar; aktiv ist allein das Muster für
  `ADR-` (Gegenprobe in §7 des Gebers, unter d-check `v0.82.0` und `v0.85.0` gleich; ebenso
  `done/slice-lastenheft-pruefbarkeit.md` §1). Dieser Slice aktiviert das Muster als
  Entscheidung über das Gate — eine Verschärfung, keine Lockerung (`AGENTS.md` §3.6 greift
  nicht).
- **Prüfauftrag `MR-`:** Heute fängt das Gate eine blanke MR-Kennung nur in den drei
  Spec-Dateien (Klasse `adaptionsblock` der Matrix). Der Kommentar in `.d-check.yml` hält ein
  `ids`-Muster dafür für schädlich, weil `harness/conventions.md` ihre eigenen Kennungen blank
  nennt; die Messung unten widerlegt das für `v0.85.0` (die Datei der Definition meldet
  nichts). Deshalb gehört `MR-` mit in diesen Slice.
- **Nebenbefund des Implementers des Gebers:** Der Block *Strenges Doc-Gate* in
  `.claude/commands/implement-slice.md`, `.claude/commands/plan-welle.md` und
  `.claude/commands/close-welle.md` sagte Link-Pflicht für `LH-`, `ADR-` und `MR-` zu.
  Der Geber fasst ihn auf das heute aktive Verhalten (sein DoD-Punkt 1, `AGENTS.md` §3.11);
  dieser Slice zieht ihn auf das Verhalten nach, das er liefert.

**Geschnitten am 2026-10-09** (Planner, Review F-556 zu `slice-harness-d-check-v0-85`): Mit
der Schichtteilung in §8 berührte der Bestand in ADRs und in Nutzer- und Wartungs-Doku eine
dritte und vierte Schicht neben Spezifikation und Harness. Ihn verlinkt vorab
`slice-harness-lh-links-bestand` (Sendung dort in §1 und DoD mit dieser Kennung); dieser
Slice startet erst, wenn jener in `done/` liegt (§4).

**Messung des Planners** (2026-10-09, Stand `825b71e`, Kopie ohne `.git` im Scratchpad,
d-check `v0.85.0` mit dem Digest aus `d-check.mk`; ein zweites `ids`-Muster vor dem für `ADR-`,
Ziel `spec/lastenheft.md` bzw. `harness/conventions.md`):

| Muster | `link-policy` | ausgenommen | Befunde | Verteilung |
|---|---|---|---|---|
| `LH-` | `always` | — | 2442 in 156 Dateien | `docs/reviews/` 1399, `done/` 330, Spezifikation 245, `next/` 186, `docs/user/` 173, übrige 109 |
| `LH-` | `prose` | — | 1441 | `docs/reviews/` 1071, `docs/user/` 173, Spezifikation 135, `next/` 34, übrige 28 |
| `LH-` | `prose` | `docs/reviews/**`, `done/**`, `observations/**` | 357 | `docs/user/` 173 (davon `abdeckung-unit.md` 156, `abdeckung-e2e.md` 17), Spezifikation 135, `next/` 34, `open/` 9, `in-progress/` 5, Rest 1 |
| `MR-` | `prose` | wie oben | 0 | — |
| `MR-` | `always` | — | 50 | `done/` 22, `docs/reviews/` 21, `in-progress/` 3, `observations/` 3, `open/` 1 |

Die Zahlen in `next/` sinken bis zum Start: Die meisten Treffer liegen in Slices von
welle-v1-abschluss, die dann in `done/` liegen. Gemessen wird beim Start neu (§4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Abdeckungstabellen `docs/user/abdeckung-*.md` mit Links schreiben — ein anderer
  Vorgang: Ihre Treffer stehen in den Kurzbeschreibungen, die `make abdeckung`
  (`tools/test/abdeckung.sh`) aus den Abdeckungs-Deklarationen der Tests schreibt; Links dort
  sind Arbeit am Generator samt seiner Gegenprobe und den Deklarationen unter `internal/`,
  eine dritte Schicht. Geplant ist, die Tabellen auszunehmen; ob das trägt, ist Randform (§6).
  Entscheidet der Architect dagegen, gilt die Rückführung in §4.
- Zeitdokumente umschreiben (`docs/reviews/**`, `docs/plan/planning/done/**`,
  `docs/plan/planning/observations/**`) — Bestand bleibt bewusst stehen: Sie frieren mit
  ihrem Vorgang ein, `docs/reviews/**` ist aus der Matrix schon ausgenommen, und ein
  Nachziehen von rund 1700 Stellen änderte keine Aussage. Geplant ist, sie im Muster
  auszunehmen (Randform, §6).
- Muster für weitere Kennungen (`SPEC-`, `ARC-`, `BEO-`, `CO-`, `slice-`, `welle-`) — ein
  anderer Vorgang: `SPEC-`/`ARC-` adressieren innerhalb der Spec (`AGENTS.md` §5),
  `slice-`/`welle-` fängt die Matrix; eine Linkpflicht für sie ist je eine eigene
  Entscheidung.
- Opt-in-Module von d-check (`modules:`), Pin und `d-check.mk` — ein anderer Vorgang: Den
  Pin lieferte `slice-harness-d-check-v0-85`; die Zeile `modules:` bleibt
  `[links, anchors, ids, matrix, spans]`.
- Ein eigenes Gegenproben-Ziel für das Doku-Gate in `harness/mk/` — ein anderer Vorgang
  (neues Gate-Werkzeug); die Gegenprobe dieses Slice ist ein Beleg in §7 wie beim Geber.
- Lastenheft- und MR-Kennungen in den ADRs und in der Nutzer- und Wartungs-Doku
  (`docs/user/` außer den Abdeckungstabellen, `docs/maintainer/`) als Links schreiben —
  Schicht-Abgrenzung: Entscheidungen und Nutzer- und Wartungs-Doku sind eigene Schichten
  (§8); sie verlinkt vorab `slice-harness-lh-links-bestand`.
- Produkt-Code und Tests unter `internal/` und `test/` — Schicht-Abgrenzung: Der Slice ändert
  Spezifikation und Harness (§8); Pläne und Roadmap trägt er als Planung mit.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Muster aktiv** (Teil-Zusage aus `slice-harness-d-check-v0-85`): `.d-check.yml` führt
      unter `ids.patterns` je ein Muster für Lastenheft-Kennungen (Ziel `spec/lastenheft.md`)
      und für Kennungen des Adaptions-Blocks (Ziel `harness/conventions.md`), mit
      `link-policy` und Ausnahmen, wie die Spezifikation sie nach §6 entscheidet; der
      Kommentar zum `MR-`-Muster in `.d-check.yml` sagt, was gemessen ist. Gegenprobe in
      Kopien des Arbeitsbaums (Beleg in §7: Zusage · Mutation · Befundzeile, `AGENTS.md`
      §3.10): eine blanke Lastenheft-Kennung und eine blanke MR-Kennung in einer gescannten,
      nicht ausgenommenen Datei sind je rot (`id-unlinked`), dieselben als Link und in einer
      ausgenommenen Datei grün, und je entschiedene Randform aus §6 ein Fall.
- [ ] **Bestand verlinkt:** Jede blanke Lastenheft- und MR-Kennung in den gescannten, nicht
      ausgenommenen Dateien unter `spec/`, `harness/`, in lebenden Plänen und in der Roadmap
      ist ein Link auf ihre Definition; ADRs und Nutzer- und Wartungs-Doku verlinkt
      `slice-harness-lh-links-bestand` vorab. `make docs-check` meldet 0 Befunde mit den
      Mustern aus Punkt 1.
- [ ] **Aussage über das Gate** (Nebenbefund aus `slice-harness-d-check-v0-85`): Der Block
      *Strenges Doc-Gate* in `.claude/commands/implement-slice.md`,
      `.claude/commands/plan-welle.md` und `.claude/commands/close-welle.md` sagt die
      Linkpflicht für Lastenheft- und MR-Kennungen so zu, wie Punkt 1 sie liefert,
      Ausnahmen eingeschlossen, und nicht weiter (`AGENTS.md` §3.11); jede Zusage dort hat
      eine Zeile der Gegenprobe in §7.
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
| `spec/spezifikation.md` §11 | update | neue Stelle: Linkpflicht der Kennungen im Doku-Gate mit den Randformen aus §6 (Architect, vor dem Code) |
| `.d-check.yml` | update | `ids`-Muster für Lastenheft- und MR-Kennungen nach der Spezifikation; Kommentare über `ids` und die Klasse `adaptionsblock` auf den gemessenen Stand |
| gescannte, nicht ausgenommene `.md` unter `spec/`, `harness/`, lebende Pläne, Roadmap | update | blanke Kennungen als Links; mechanisch, keine Aussage ändert sich. ADRs und Nutzer- und Wartungs-Doku: `slice-harness-lh-links-bestand` |
| `.claude/commands/implement-slice.md`, `.claude/commands/plan-welle.md`, `.claude/commands/close-welle.md` | update | Block *Strenges Doc-Gate* auf das gelieferte Verhalten |
| dieser Plan, §7 | update | Messung beim Start, Gegenprobe, Läufe |

- Eine angenommene ADR wird nicht inhaltlich überschrieben (`AGENTS.md` §3.5). Ob ein Link
  um eine schon genannte Kennung dort zulässig ist, entscheidet der Nutzer vor
  `slice-harness-lh-links-bestand`; lautet die Antwort nein, ist das Ausnehmen der ADRs
  Randform dieses Slice (§6).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): welle-v1-abschluss liegt in `done/`,
`slice-harness-d-check-v0-85` liegt in `done/` (er liefert den Pin und die enge Fassung des
Blocks *Strenges Doc-Gate*, die dieser Slice erweitert), und `slice-harness-lh-links-bestand`
liegt in `done/` (ADRs und Nutzer- und Wartungs-Doku verlinkt); WIP-Limit 1. Platz in der
Harness-Reihe nach den Wellen (Drift-Log der Roadmap vom 2026-10-08): an ihrem Ende, nach
`slice-harness-mutation`; ein Vorziehen entscheidet der Nutzer beim Übergang `open` → `next`.
Beim Start misst der Implementer die Tabelle aus §1 neu und trägt sie in §7 ein.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung beim Start ergibt mit
  den entschiedenen Ausnahmen mehr Befunde, als ein Diff in einer Review-Sitzung trägt
  (Richtwert 400), oder der Architect entscheidet, dass die Abdeckungstabellen Links tragen
  müssen (dritte Schicht: Generator und Test-Deklarationen). Schnitt dann: Muster und
  Gegenprobe mit den Ausnahmen bleiben hier, der Bestand außerhalb der Spezifikation oder
  der Generator wird ein eigener Slice. Ebenso, wenn die Messung beim Start Treffer in ADRs
  oder in Nutzer- und Wartungs-Doku findet, die seit `slice-harness-lh-links-bestand` neu
  entstanden sind: Sie gehen an den Planner, nicht in diesen Slice (dritte Schicht).
- `in-progress` → `open` (blockiert — Carveout?): d-check kann eine entschiedene Randform
  nicht ausdrücken (etwa Ausnahmen je Muster oder die Unterpunkt-Form der
  Lastenheft-Kennung); dann zuerst die Klärung mit dem Nutzer — Change Request an d-check
  oder andere Entscheidung —, das Muster bleibt bis dahin aus.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Gegenprobe in §7 je Zusage rot gesehen, `make docs-check` 0 Befunde mit den
neuen Mustern), Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12): Der Slice liefert einen neuen Vertrag (eine Zusage des
Doku-Gates). Entschieden wird jede vor dem ersten Commit mit Gate-Konfiguration, in der
Spezifikation §11, vom Architect; offen sind heute alle:

- **`link-policy`:** `always` (auch Inline-Code, wie beim Muster für `ADR-`) oder `prose`
  (Messung in §1: 2442 gegen 1441 Befunde ohne Ausnahmen).
- **Ausnahmen:** welche Pfade (`exempt-paths` je Muster, nur bei `always` dokumentiert —
  die Messung in §1 nahm sie auch bei `prose` an; zu belegen): Zeitdokumente
  (`docs/reviews/**`, `done/**`, `observations/**`), die erzeugten Abdeckungstabellen
  `docs/user/abdeckung-*.md`, angenommene ADRs (nur, wenn der Nutzer vor
  `slice-harness-lh-links-bestand` Links dort ablehnt).
- **Unterpunkt-Form:** eine Lastenheft-Kennung mit Unterpunkt (Form `.a`) — fängt das Muster
  nur den Stamm, und worauf zeigt der Link (Datei oder Anker der Anforderung)?
- **Datei der Definition:** `spec/lastenheft.md` und `harness/conventions.md` nennen ihre
  Kennungen blank (Überschriften, Tabellen); die Messung zeigt dort keine Befunde — ob das
  eine Zusage von d-check ist, auf die sich das Gate stützen darf.
- **Kennung in Code-Blöcken und Commit-Beispielen** (etwa `AGENTS.md` §5, Hilfetexte der
  Commands): Link-Pflicht oder Ausnahme.
- **Muster-Grenzen:** Klassen `FA`, `QA`, `RB` und zweistellige Nummer — eine Kennung
  außerhalb dieser Form (Platzhalter in Vorlagen, Präfix ohne Nummer) bleibt unberührt.

**Risiken:**

- Die Ausnahmen werden so weit, dass das Gate dort grün bleibt, wo neue Kennungen entstehen
  (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`) — **Ausgang:** offen bis Closure;
  die Gegenprobe in DoD-Punkt 1 hat je Ausnahme einen Fall, der zeigt, dass ein lebender
  Pfad daneben rot wird.
- Der Bestand wächst bis zum Start (jeder Slice bis dahin schreibt Kennungen, die heute kein
  Gate prüft) — **Ausgang:** offen bis Closure; Messung beim Start (§4), Rückführung bei mehr
  als dem Richtwert.
- Ein Massen-Ersetzen trifft eine Kennung in einer Aussage, die sich dabei ändert (etwa in
  einer Tabelle mit Kennungen als Daten) — **Ausgang:** offen bis Closure; der Diff ist nach
  Datei gruppiert und je Gruppe im Review prüfbar.

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

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel
`REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `825b71e` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — Herkunft dieses
  Slice: DoD-Punkt 1 von `slice-harness-d-check-v0-85` las `.d-check.yml` als Linkpflicht
  für Lastenheft-Kennungen. Ob das der dritte Beleg ist, entscheidet die Closure des Gebers;
  hier daraus das erste Risiko in §6 und die Gegenprobe je Ausnahme.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (27×, verkörpert, `AGENTS.md` §3.11) —
  DoD-Punkt 3: jede Zusage im Block *Strenges Doc-Gate* mit einer Zeile der Gegenprobe.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (18×, verkörpert, §3.12) und
  `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos
  und liefert einen neuen Vertrag; die Randformen stehen deshalb schon in §6, und der Weg zum
  Architect vor dem Code ist der Übergang `open` → `next` (Nutzer).
- `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum` (1×, offen) — die Festlegung der
  Linkpflicht gehört in die Spezifikation §11, nicht in `.d-check.yml`-Kommentare oder die
  Commands; daher die Spec-Zeile in §3.
- `BEO-REPO/slice-waechst-durch-uebernahmen` (3×, verkörpert) und
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (5×, verkörpert) — dieser Slice ist Nehmer
  einer Sendung (§1). **Nachgezählt mit der Sendung** (`AGENTS.md` §3.13), nach der
  Schichtteilung unten neu am 2026-10-09 (Review F-556): Die Zählung beim Anlegen
  (`dfb7e88`) fasste Spezifikation, Pläne und Commands zu „Dokumentation“ zusammen und kam
  auf zwei Schichten; nach der Teilung, die auch `slice-harness-d-check-v0-85` zählt, waren
  es vier (Spezifikation; Harness: `.d-check.yml` und Commands; Entscheidungen; Nutzer- und
  Wartungs-Doku). Deshalb geschnitten (§1): Nach dem Schnitt drei Liefer-Punkte (Muster mit
  Gegenprobe, Bestand, Aussage über das Gate) und zwei Schichten (Spezifikation: neue Stelle
  in §11 und Links; Harness: `.d-check.yml`, Commands, Links in `harness/`); Pläne und
  Roadmap sind Planung. Die Abdeckungstabellen samt Generator bleiben draußen (§1), sonst
  kämen die Abdeckungs-Deklarationen der Tests unter `internal/` als weitere Schicht hinzu. Der abgeschnittene Teil,
  `slice-harness-lh-links-bestand`, ist dort nachgezählt (§8).

Keiner der Einträge erreicht mit diesem Plan neu die Schwelle 3×.

**Schichtteilung** (gleich in `slice-harness-d-check-v0-85`, `slice-harness-lh-links-bestand`
und diesem Slice; festgelegt am 2026-10-09 nach Review F-556): Eine Schicht im Sinn von
Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice ist ein Bereich mit
eigenem Maßstab im Review — die Spezifikation (`spec/`) · die Entscheidungen
(`docs/plan/adr/`) · der Harness (Gate-Konfiguration und Gate-Werkzeuge, `Dockerfile`,
Agenten-Anweisungen: `AGENTS.md`, Skills, Commands, `harness/`) · die Nutzer- und
Wartungs-Doku (`docs/user/`, `docs/maintainer/`) · im Produkt-Code je Schicht des Hexagons,
ebenso deren Tests. So zählten `slice-lastenheft-pruefbarkeit` (Harness und Lastenheft mit
Spezifikation als verschiedene Schichten) und `slice-harness-lint-werkzeug` (Harness mit
`Dockerfile` als eine Schicht, Tests je Schicht des Hexagons). Keine Schicht ist die
Planung (Pläne, Roadmap, Register): In ihr schreibt jeder Slice (eigener Plan, Sendungen,
Drift-Log); was sie an Umfang trägt, misst das dritte Kriterium, eine Review-Sitzung.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
