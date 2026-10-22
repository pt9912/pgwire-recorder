# Slice slice-harness-regeln-gate-konfiguration: Regel — die Konfiguration eines Gates ist eine Zusage und wird gegen das echte Werkzeug gemessen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Vorschlag des Planners, nicht eingetragen: welle-harness-regeln, hinter deren fünf Regel-Slices; die Welle und das Drift-Log der Roadmap ändern sich erst mit der Entscheidung des Nutzers). Angelegt am 2026-10-18 bei der Closure von `slice-harness-ci-range-konfiguration`, die das Register auf 3× hob; Auftrag der Koordination, **Vorschlag des Planners**.

**Bezug:** kein Anspruch des Produkts — der Slice ändert Regeln des Harness. Bindung an das Register: `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (3×: `slice-walking-skeleton-build-gates`, `slice-harness-lint-werkzeug`, `slice-harness-ci-range-konfiguration`); Regelstellen: `AGENTS.md` §3.10 (seit welle-extended-query) und `.claude/commands/implement-slice.md` Schritt 19. Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** Planner (Vorschlag). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Regel zu `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` steht im Briefing und im Command des Implementers: (1) `AGENTS.md` §3.10 sagt, dass die Konfiguration eines Gates (`.d-check.yml`, `.a-check.yml`, `.golangci.yml`), die als Zusage gelesen wird, mit dem echten Werkzeug geprüft wird, je Eintrag ein freier und ein roter Fall, und dass die rote Menge einer Mutation der Konfiguration vor dem Festschreiben der Erwartung in einem Entdeckungslauf gemessen wird; (2) `.claude/commands/implement-slice.md` Schritt 19 verlangt diesen Lauf vor der Mutationstabelle in §7.

**Entwurf des Wortlauts und Entscheidung des Nutzers.** Regeländerungen in `AGENTS.md` und in `.claude/commands/` sind eine Entscheidung des Nutzers. Dieser Slice ist ein **Vorschlag des Planners**; weder seine Anlage noch der Ausgang `geplant` im Register ist vom Nutzer bestätigt, und den Wortlaut hat er nicht entschieden. Der Implementer legt den Wortlaut aus §3 dem Nutzer vor und committet beide Dateien erst nach der Bestätigung (Nutzerentscheidung über den Wortlaut vor der Umsetzung); was der Nutzer nicht bestätigt, nimmt er nicht auf und nennt es in §7.

**Beleg-Stand** (Stand 2026-10-18, gemergter Stand des Registers): (a) `slice-walking-skeleton-build-gates`: Eine Regel mit Liste beider Adapter erlaubte nur einen; die Gegenprobe mit fünf Fällen zeigte es. (b) `slice-harness-lint-werkzeug`: vier Stellen von `.golangci.yml` lasen sich anders, als sie wirken (V-72, V-74, F-433). (c) `slice-harness-ci-range-konfiguration`: vier Erwartungen vor dem ersten Lauf der Mutationen von `.d-check.yml` waren falsch, weil `vcs` die Konfiguration anders auswertet (`head-allow` prüft jede unveränderliche Datei, ein `status-line`-Muster ohne Treffer färbt fast alle Fälle); erst der Entdeckungslauf `GP_ENTDECKE=1` mit dem echten Image zeigte die Mengen. Die Belege stehen im Register, nicht hier.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Sensor, der prüft, dass jede Konfiguration eines Gates eine Gegenprobe mit echtem Werkzeug hat — anderer Vorgang (Werkzeug statt Regel): ein Sensor bekäme einen Vertrag in `harness/vertraege/` und eine eigene Gegenprobe, und was „Eintrag“ einer Konfiguration heißt, ist je Werkzeug eine Auslegung. Kein Folge-Slice benannt; ein vierter Beleg nach der Regel ist der Anlass.
- Das Nachrüsten von Gegenproben für Konfigurationen im Bestand (übrige Blöcke und Klassen von `.d-check.yml`, `.a-check.yml`, `.golangci.yml`) — Bestand bleibt bewusst stehen: Keine Beobachtung hat ihn getroffen, und ein Nachrüsten ohne Befund erzeugte Fälle ohne Anlass; gezählt wird neu, wenn ein Slice eine dieser Dateien ändert (dann gilt die Regel).
- Änderung an `.harness/skills/reviewer.md`, `.claude/agents/` und den Commands des Planners — Schicht-Abgrenzung und Größe: Eine dritte Datei brächte einen dritten Liefer-Punkt und eine dritte Schicht; ob der Reviewer die Regel prüft, ist ein Folgeanlass, wenn der Review sie dreimal nicht trifft.
- Ein neuer Schalter für den Entdeckungslauf in bestehenden Gegenproben (`GP_ENTDECKE` gibt es heute nur in `range-konfiguration-gegenprobe`) — das wäre Arbeit am Werkzeug; die Regel verlangt den Lauf, nicht seinen Namen.
- Produkt-Code, Gates, Spezifikation, Benutzerhandbuch — Schicht-Abgrenzung: **kein Produktcode**, zwei Dokumente des Harness.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Briefing):** `AGENTS.md` §3.10 trägt einen Absatz zur Konfiguration eines Gates als Zusage: Prüfung mit dem echten Werkzeug, je Eintrag ein freier und ein roter Fall, die rote Menge einer Mutation der Konfiguration vor dem Festschreiben der Erwartung gemessen; mit Falsch/Richtig, Begründung und Retirement-Check `seit slice-harness-regeln-gate-konfiguration`. Der Wortlaut ist vom Nutzer bestätigt, bevor die Datei committet wird (§7 nennt die Bestätigung).
- [ ] **Liefer-Punkt 2 (Command):** `.claude/commands/implement-slice.md` Schritt 19 verlangt den Entdeckungslauf vor der Mutationstabelle in §7 für Slices, die eine Konfiguration eines Gates liefern oder ändern; der Wortlaut ist vom Nutzer bestätigt.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `AGENTS.md` und `.claude/` ändern) und `make docs-check` ohne Befund.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Ausgang von `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` von `geplant` auf `verkörpert` mit Zielort und Herkunfts-Anker.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.10 | update | Absatz „Die Konfiguration eines Gates ist eine Zusage“ (Empfehlung des Wortlauts unten), Herkunfts-Anker `seit slice-harness-regeln-gate-konfiguration` |
| `.claude/commands/implement-slice.md` Schritt 19 | update | Entdeckungslauf vor der Mutationstabelle bei einer Konfiguration eines Gates |
| `docs/plan/planning/observations/BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen/state.md` | update (Closure) | `geplant` → `verkörpert` mit Zielort und Anker |

**Empfehlung des Wortlauts (nicht entschieden):** Wer die Konfiguration eines Werkzeugs als Zusage trägt (`.d-check.yml`, `.a-check.yml`, `.golangci.yml`), prüft jeden Eintrag mit dem echten Werkzeug, je Eintrag ein freier und ein roter Fall. Die rote Menge einer Mutation der Konfiguration wird in einem Lauf gemessen, der sie zeigt, bevor sie als Erwartung steht; aus dem Lesen der Konfiguration allein wird sie nicht abgeleitet. *Falsch:* Die Erwartung „`head-allow` färbt nur den Fall mit Änderung“ steht vor dem Lauf; das Werkzeug prüft jede unveränderliche Datei im Kopf, und die Mutation färbt vier Fälle. *Richtig:* Der Entdeckungslauf zeigt die Menge, die Erwartung folgt ihr, und ein Eintrag, der alle Fälle färbt, bekommt den Grund dazu. *Begründung:* Eine Konfiguration liest sich als Zusage, das Werkzeug wertet sie anders aus; erst der Lauf mit dem echten Werkzeug zeigt es (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): **pausiert seit dem 2026-10-22** (Entscheidung des Nutzers, Drift-Log der Roadmap und §5 von welle-harness-regeln): Die Priorisierung wartet auf `slice-harness-regeln-rueckbau` in `done/`; sonst gilt, was folgt. Der Nutzer hat den Slice priorisiert und den Ausgang `geplant` des Registers bestätigt (oder `verkörpert` mit anderem Zielort entschieden, dann entfällt dieser Slice nach `done/` mit `Gegenstand: entfallen`); der Wortlaut ist vor dem ersten Commit bestätigt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: mehr als zwei Liefer-Punkte oder mehr als zwei Schichten, oder der Diff ist nicht in einer Review-Sitzung prüfbar. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- AGENTS.md .claude/commands`, liegt über 200; Basis: Elternteil des Commits, der diese Datei nach `in-progress/` bewegt. Schätzung des Planners in dieser Metrik: 40 bis 70, mit dem Doppelten für Prosa mit wörtlichen Dateiinhalten (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, Beleg `slice-harness-devcontainer-u-boot-0-10`) 80 bis 140.
- `in-progress` → `open`: Der Nutzer bestätigt den Wortlaut nicht oder entscheidet eine andere Verkörperung.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Review- und Verifikations-Report liegen vor, Wortlaut vom Nutzer bestätigt, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Nutzer entscheidet eine andere Verkörperung oder keine** (der Ausgang `geplant` ist ein Vorschlag). — **Ausgang (vorgesehen):** eingetreten → der Slice geht mit `Gegenstand: entfallen` nach `done/`, das Register trägt den Ausgang des Nutzers.
- **Die Regel liest sich als Pflicht für jede Konfiguration, auch ohne Zusage** und erzeugt Fälle ohne Anlass. — **Ausgang (vorgesehen):** entfallen, wenn der Wortlaut „als Zusage trägt“ nennt und das Review ihn daran misst.
- **Der Entdeckungslauf hängt am Werkzeug der jeweiligen Gegenprobe** (es gibt heute nur einen Schalter in `range-konfiguration-gegenprobe`). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, falls die Regel an einer Gegenprobe ohne Schalter scheitert; sonst entfallen.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln.

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (3×, Gegenstand dieses Slice); `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (`AGENTS.md` §3.10, verkörpert; der neue Absatz steht daneben und ersetzt keinen Satz); `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (7×, ohne Ausgang; Nachbar-Regel an derselben Stelle, vom Nutzer zu entscheiden, nicht Gegenstand hier). Kein Eintrag erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Der Slice trägt und nimmt keine Sendung. Liefer-Punkte 2 (Briefing | Command). Schichten 2 (Briefing, Commands; so zählen `welle-harness-regeln` und `slice-harness-regeln-nach-v1-abschluss`); Planung und Register sind keine Schicht, Handbuch und Wartungs-Doku bleiben unberührt. Damit auf der Grenze von zwei Schichten; ein dritter Liefer-Punkt oder eine dritte Datei schnitte den Slice.
