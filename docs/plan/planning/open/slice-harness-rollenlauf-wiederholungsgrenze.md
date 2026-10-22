# Slice slice-harness-rollenlauf-wiederholungsgrenze: Nach dem ersten vollen Lauf einer Gegenprobe je Rollenlauf sperrt der Wächter einen weiteren ohne `GP_NUR` oder `GP_MUT`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-harness-regeln, **hinter** `slice-harness-rollenlauf-zeitgrenze-haertung` (§4 und §5 der Welle). **Zuordnung und Stelle sind Empfehlung des Planners vom 2026-10-22**, nicht Entscheidung des Nutzers (beide Vorgänger ändern Hook und Gegenprobe; ein gleichzeitiger Lauf erzeugt einen Konflikt). Entschieden hat der Nutzer die **Sache** am 2026-10-21 (die Wiederholungs-Grenze für volle Gegenprobe-Läufe, freigegeben, aber nicht im Slice des Wächters) und am 2026-10-22 die Stunde samt Wächter. Angelegt am 2026-10-22 in `open/`, **sobald der Vertrag des Wächters entschieden war**: Die Wiederholungs-Grenze erbt dessen Datenquelle, Anker und Fail-Richtung ([ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md), Folgepflicht, eigene Entscheidung); den Slice, den `slice-harness-rollenlauf-zeitgrenze` in §1 als Folge-Slice nannte, ohne ihn zu benennen.

**Bezug:** Kein Anspruch des Produkts (`LH-*`) — der Slice erweitert einen Wächter über Agentenläufe. Anlass: 46 volle Läufe einer Gegenprobe zu je etwa 175 s in einem Rollenlauf, während der Fehler in zwei Fällen saß; ein voller Lauf nach dem ersten ist selten nötig, und `GP_NUR` und `GP_MUT` fahren genau die Fälle. Aktive ADRs: [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) (Folgepflicht: erbt Datenquelle, Anker und Fail-Richtung), [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Zeitbudget je voller Lauf, 600 s), [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenprobe hängt an `GEGENPROBEN_CONTAINER`).

**Berührte Spec-Stellen:** [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) (die Wiederholungs-Grenze erbt Strom, Anker, Befehl aus `tool_input.command` und Fail-Richtung; ob der Architect sie als neue Punkte dort oder als eigene Kennung fasst, entscheidet er vor dem Code) · [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) (gelesen: Namen von `GP_NUR` und `GP_MUT`, Zeitbudget, Punkte 4 bis 12). **Kein Vertragstext steht in diesem Plan.**

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Rollenlauf fährt eine Gegenprobe voll **einmal**; ein weiterer voller Lauf derselben oder einer anderen Gegenprobe in demselben Lauf wird vom Wächter gesperrt, solange der Befehl weder `GP_NUR` noch `GP_MUT` trägt; die Meldung nennt den Weg (der enge Fall). Die Zählung folgt aus dem Strom der Erfassung des laufenden Agenten (die Dauer `duration_ms` der Aufrufe von `make`, die eine Gegenprobe voll fuhren) und nutzt den Strom, den Anker und die Fail-Richtung des Wächters von `slice-harness-rollenlauf-zeitgrenze`; der Slice erweitert den Hook, er legt keinen zweiten an.

**Sendung von `slice-harness-rollenlauf-zeitgrenze`** (dessen §1 *Ausdrücklich NICHT*, Klasse 1; eingetragen vom Planner im Commit des Schnitts, `AGENTS.md` §3.13): Die Wiederholungs-Grenze ist dieser Slice; der Wächter der Stunde sagt dazu nichts zu. `grep` nach der Kennung des Gebers findet die Sendung in §1 und DoD.

**Kein Vertragstext, keine Antworten des Planners.** Die Randformen sind **Fragen an den Architect**, zu entscheiden vor dem ersten Code-Commit (`AGENTS.md` §3.12); der Plan stellt sie in §6 und beantwortet sie nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der **Wächter der Stunde** (Hook, Leser, Eintrag, Gegenprobe-Gerüst, Nachweis im Lauf) und seine **Härtung** — Klasse 1: `slice-harness-rollenlauf-zeitgrenze` und `slice-harness-rollenlauf-zeitgrenze-haertung` (beide vor diesem Slice in `done/`); dieser Slice erweitert den Hook und ändert die Zusagen der Stunde nicht.
- Die **Regel** in `AGENTS.md` und den sechs Rollen — Klasse 1: `slice-harness-rollenlauf-zeitgrenze-regel`. Ob die Wiederholungs-Grenze einen Satz in `AGENTS.md` §6 bekommt, ist eine Entscheidung des Nutzers nach der Lieferung; eine Adresse gibt es nicht.
- Ein **Anker für die Hauptsitzung** und eine Zählung außerhalb von Rollenläufen — Klasse 3 und Grenze von `SPEC-070` Punkt 12 (1); eine Adresse gibt es nicht.
- Eine **Änderung des Zeitbudgets** von 600 s je voller Lauf ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)) und der Namen `GP_NUR` und `GP_MUT` ([`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md)) — `AGENTS.md` §3.6: Die Wiederholungs-Grenze ist eine zusätzliche Grenze, keine Senkung einer bestehenden.
- Ein **Abbruch** eines laufenden Befehls — Grenze des Wächters (kein Hook innerhalb eines Vordergrund-Befehls) und Entscheidung des Nutzers.
- Die **Erfassungsschicht**, der Stop-Hook, die Gate-Ketten, Produkt-Code und Nutzer-Doku — Schicht-Abgrenzung (Klasse 4): Der Wächter liest nur.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (die Wiederholungs-Grenze im Wächter):** Der Hook `.claude/hooks/pretooluse-zeitgrenze.sh` sperrt in einem Rollenlauf nach dem ersten vollen Lauf einer Gegenprobe jeden weiteren ohne `GP_NUR` oder `GP_MUT`, mit einer Meldung, die den engen Fall nennt (Sendung von `slice-harness-rollenlauf-zeitgrenze`, §1); er erfüllt die Zusagen, die der Architect vor dem Code festhält (Quelle der Zählung, was „voll“ heißt, woran `GP_NUR` und `GP_MUT` erkannt werden, was für eine andere Gegenprobe gilt). Jede Randform aus §6 ist entschieden, bevor die erste Zeile Code entsteht.
- [ ] **Liefer-Punkt 2 (die Gegenprobe):** `make zeitgrenze-gegenprobe` trägt die Fälle und Mutationen der Wiederholungs-Grenze, **je Zusage und je Element ihrer Aufzählungen ein Fall und eine Mutation, die nur ihn rot färbt** (`AGENTS.md` §3.10); die unveränderte Kopie ist grün, **kein Fall wartet**, der volle Lauf liegt **unter 600 s** ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md); die gemessene Zeit steht in §7), und die Fälle der Stunde bleiben unverändert grün.
- [ ] `make docs-check` und `make kopf-check` grün.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/`, `harness/` und `tools/` ändern).
- [ ] Die Gegenprobe ist am Endstand **in einem Zug** grün, mit ihrer Zeit in §7 ([`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16); während der Arbeit fährt der Implementer nur den engsten Fall mit `GP_MUT`.
- [ ] Der Wächter ist mit der neuen Zusage **einmal im Agentenlauf gesehen**, soweit sie die Eingabe oder die Form des Blocks berührt (Beleg als Datei unter `docs/reviews/`, `SPEC-070` Punkt 13); berührt sie beide nicht, steht in §7, warum der Nachweis der Stunde trägt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`): die gemessene Zeit als Beleg in `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`; ein weiterer Beleg in `BEO-REPO/rollenlauf-ohne-zeitgrenze` oder „keine Beobachtung angefallen“ in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der nächsten Welle-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-070-…md` oder eine eigene Kennung | neu oder update, **durch den Architect** vor dem Code | Zusagen der Wiederholungs-Grenze und jede Randform aus §6 entschieden (`AGENTS.md` §3.12) |
| `.claude/hooks/pretooluse-zeitgrenze.sh`, `tools/harness/hook-felder.awk` | update | Liefer-Punkt 1: die Zählung aus dem Strom und die Sperre |
| `tools/dev/harness/zeitgrenze-gegenprobe.sh` | update | Liefer-Punkt 2: Fälle und Mutationen |
| `docs/plan/planning/observations/BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/` | update bei Closure | gemessene Zeit (kein Liefer-Punkt) |

- **Reihenfolge:** (1) der Architect entscheidet die Randformen aus §6 in der Spezifikation; (2) der Hook; (3) die Gegenprobe, je Zusage eine Mutation, vom Implementer selbst rot gesehen; (4) Closure. Jeder Commit nennt die Kennung dieses Slice.
- **Was der Implementer nicht entscheidet:** Trifft er eine Randform, die der Vertrag nicht nennt, hält er an und gibt sie dem Architect zurück (`AGENTS.md` §3.12).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): WIP-Limit frei; `slice-harness-rollenlauf-zeitgrenze` und `slice-harness-rollenlauf-zeitgrenze-haertung` liegen in `done/` (der Hook, sein Strom, Anker und Fail-Richtung stehen); der Architect hat die Randformen aus §6 entschieden (`AGENTS.md` §3.12). Der Übergang `open` → `next` setzt `Verantwortlich:`.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (der Befehl steht in `AGENTS.md` §3.14 nach `slice-harness-groessen-metrik`): Summe aus Einfügungen und Löschungen über die Pfade `.claude`, `harness` und `tools`, mit den Ausschlüssen des Befehls; Basis ist der Commit, der diese Datei nach `in-progress/` bewegt. **Vorläufig Marke 700, Meldeschwelle 550**; der Architect bestätigt oder ändert beide mit dem Vertrag. Schätzung des Planners **vor dem Vertrag** (2026-10-22, in dieser Metrik, mit dem Faktor 3 bis 4,6 für Gegenprobe-Code): Hook 40 bis 80, Gegenprobe 120 bis 370, zusammen **160 bis 450, das sind 23 bis 64 % der Marke**; sie ist eine Spanne ohne Vertrag und wird mit ihm neu gerechnet.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über der Marke, oder der Vertrag verlangt einen dritten Liefer-Punkt.
- `in-progress` → `open` (blockiert): Der Strom trägt die Dauer nicht, die die Zählung braucht (die Feldliste der Erfassung nennt `duration_ms`, der Befund wäre ein Vorgang am Werkzeug), oder der Architect kann eine Randform nicht ohne eine Entscheidung des Nutzers entscheiden.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` grün; die Gegenprobe am Endstand in einem Zug grün mit ihrer Zeit unter 600 s in §7; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen — Fragen an den Architect, vor dem ersten Code-Commit zu entscheiden** (`AGENTS.md` §3.12: in der Spezifikation; der Planner stellt die Fragen und beantwortet sie nicht):

1. **Was „voll“ heißt und woran `GP_NUR` und `GP_MUT` erkannt werden:** der Befehl aus `tool_input.command` (Wortgrenze, Zeilenumbruch, `ssh`-Aufruf, Variable aus der Umgebung statt im Befehl); ein Aufruf mit beiden, mit leerem Wert, mit einem Wort in einer Commit-Message.
2. **Was gezählt wird:** je Rollenlauf alle Gegenproben zusammen oder je Gegenprobe; ob eine Gegenprobe, die rot endete oder abgebrochen wurde, als voller Lauf zählt; welcher Aufruf der Aufruf einer Gegenprobe ist (das Ziel `make <name>-gegenprobe`, `make gegenproben`).
3. **Quelle der Zählung:** `duration_ms` und der Befehl im Strom (das Feld `program` trägt nur das erste Wort); ob ein Aufruf mit unbekannter Dauer zählt; wie ein Lauf zählt, der über die Stunde hinaus lief (Anker des Wächters, Fortsetzung).
4. **Fail-Richtung und Selbstauskunft:** Der Wächter der Stunde failt offen mit Zeile auf stderr ([ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md)); die Wiederholungs-Grenze erbt das. Je Lage (kein Strom, Zählung nicht lesbar) eine Antwort.
5. **Die Meldung:** Wortlaut, der den engen Fall nennt, und was frei bleibt.
6. **Weg über die Grenze:** keiner im Hook (wie [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) Punkt 9); ob ein zweiter voller Lauf nach einer Fortsetzung des Agenten wieder frei ist.
7. **Aufrufarten, Umgebung, Plattform, Stelle im Ablauf:** wie der Wächter der Stunde (Punkt 8 und 11 dort).

**Risiken:**

- **Die Sperre nimmt einem Rollenlauf einen berechtigten zweiten vollen Lauf** (ein Lauf nach einem Fehler in einer anderen Gegenprobe) — **Ausgang:** offen bis Closure (entfallen, wenn der Vertrag je Gegenprobe zählt oder die Meldung den Weg nennt, den der Nutzer bestätigt).
- **Die Zählung irrt** (ein Aufruf mit `GP_MUT` in der Umgebung statt im Befehl gilt als voll) — **Ausgang:** weiter offen → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, falls ein Lauf trotz der Sperre 46 volle Läufe fährt.
- **Der Wächter wächst über Marke oder Meldeschwelle** (Schätzung 160 bis 450 vor dem Vertrag) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen: *Belege des Implementers* (Mutationstabelle je Zusage: Zusage · Mutation · roter Fall), die Zeit des vollen Gegenprobe-Laufs, die Größe in der Metrik von §4 mit dem Ergebnis der Schätzung gegen die Messung, und der Beleg des Registers.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (14×, `verkörpert`) — der Anlass: 46 volle Läufe in einem Rollenlauf; `BEO-REPO/rollenlauf-ohne-zeitgrenze` (1×) — die Datenquelle und der Wächter; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/spec-randform-erst-im-review-entschieden` (beide `verkörpert`) — je Zusage ein Fall und eine Mutation, Randformen vor dem Code. Eine Schwelle wird mit diesem Slice nicht neu erreicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-22, im Commit, der diesen Plan anlegt und die Sendung von `slice-harness-rollenlauf-zeitgrenze` einträgt):** Liefer-Punkte 2 (die Wiederholungs-Grenze im Wächter | die Gegenprobe dazu). Schichten 1 (Harness); Planung, Vertrag und Entscheidung zählen nicht. Schätzung 160 bis 450 in der Metrik von §4 mit dem Faktor 3 bis 4,6 für Gegenprobe-Code, **vor dem Vertrag**; Marke 700 und Meldeschwelle 550 vorläufig. Die Sendung legt keinen Liefer-Punkt und keine Schicht an.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Ortsangabe aus `slice-harness-gegenproben-umzug-dev` einträgt): Die Gegenproben liegen seit jenem Slice unter `tools/dev/<bereich>/`, die Pfade der Gegenproben in diesem Plan sind nachgezogen (`tools/dev/` ist wie `tools/harness/` Harness-Werkzeug, keine neue Schicht). Ortsangabe der Gegenprobe geändert, Liefer-Punkte (2) und Schichten (1) unverändert.
