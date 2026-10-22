# Slice slice-harness-rollenlauf-zeitgrenze-haertung: Der Wächter der Stunde nennt seine Lage und ist gehärtet: strenge Zeitform des Ankers, Reihenfolge der Gründe, Gründe der Freigabe mit der Zeile auf stderr und die Fälle der Umgebung — Stufe 2

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-harness-regeln, **direkt hinter** `slice-harness-rollenlauf-zeitgrenze` (§4 und §5 der Welle). **Zuordnung und Stelle sind Empfehlung des Planners vom 2026-10-22**, nicht Entscheidung des Nutzers. Entschieden hat der Nutzer die **Sache**: die Grenze von einer Stunde je Rollenlauf mit einem Wächter im Werkzeug-Aufruf (2026-10-21, wörtlich: „Wir haben rollen (siehe .claude/agents) und commands -- max. 1 Stunde“, dazu sein „ok“ auf den Entwurf des Mechanismus). **Geschnitten am 2026-10-22** aus `slice-harness-rollenlauf-zeitgrenze`: Dessen Schätzung (Hook und Gegenprobe zusammen 703 bis 1212 in der Metrik von [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md)) lag über der Marke, und die Zusagen von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) lassen sich nach **Wirkung** und **Härtung und Selbstauskunft** ordnen; dieser Slice ist die zweite Stufe. Angelegt in `open/`, weil er nicht priorisiert ist und hinter der ersten Stufe und der Stufung des Vertrags steht (§4).

**Bezug:** Kein Anspruch des Produkts (`LH-*`) — der Slice härtet einen Wächter über Agentenläufe. Bindung an das Register: `BEO-REPO/rollenlauf-ohne-zeitgrenze` (1×), `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test` (kein Fall wartet). Aktive ADRs: [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) (Proposed; Entscheidung und Gründe des Wächters), [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Zeitbudget der Gegenprobe, 600 s, gilt auch für die erweiterte Gegenprobe), [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenprobe hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS`).

**Berührte Spec-Stellen:** [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) (die Punkte der zweiten Stufe, nach der Zuordnung des Architect: Randform 7 in §6 von `slice-harness-rollenlauf-zeitgrenze`) · [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) (gelesen: Zeitbudget, Aufrufarten, Signale der Gegenprobe). Der Slice ändert keinen Vertrag.

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

**Ziel:** Der Wächter, den `slice-harness-rollenlauf-zeitgrenze` liefert (Stufe 1: die Stunde wirkt), nennt seine Lage und hält seine Zusagen an den Rändern: (1) die **Zeitform des Ankers** ist streng (`YYYY-MM-DDTHH:MM:SSZ` mit Wertebereich, Schaltjahr, Monatsende; die Form nach Punkt 3 (e) von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) mit ihren Bereichsformen), (2) die **Reihenfolge der Gründe** (Punkt 4), (3) die **Gründe der Freigabe und die Zeile auf stderr** (Punkte 6 und 7: je Grund genau eine Zeile `zeitgrenze: nicht geprüft (<Grund>)`, die Aufzählung der Gründe), (4) die **Fälle der Umgebung** (Punkt 8: leere Umgebung, Eingabe bis 1 MB unter der Zeitgrenze des Hooks, gleichzeitige Aufrufe). Dazu die Fälle und Mutationen der Gegenprobe, je Zusage und je Element ihrer Aufzählungen, im selben Ziel `make zeitgrenze-gegenprobe`.

**Warum eine zweite Stufe und kein Anhang der ersten:** Die erste Stufe wirkt, wo der Strom liegt und lesbar ist, und gibt sonst still frei; die Selbstauskunft ist die Gegenmaßnahme dafür, dass der Hook dort schweigt, wo er gebraucht wird ([ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md), Punkt 4). Hook und Gegenprobe der ersten Stufe lagen mit allen Zusagen bei 703 bis 1212; das sind 106 % der Marke 900. Die Zuordnung der Zusagen zu zwei Lieferungen, die jede einzeln lieferbar halten (jede Zusage mit ihrem Fall und ihrer Mutation, keine Zusage ohne Prüfung zwischen den Stufen), ist **Vorschlag des Planners** an den Architect; der Vertrag trägt sie, sobald der Architect sie festhält (§4, Start).

**Sendung von `slice-harness-rollenlauf-zeitgrenze`** (eingetragen vom Planner im Commit des Schnitts, `AGENTS.md` §3.13): Die Punkte der zweiten Stufe aus der Zuordnung dort (§6, Randform 7) sind dieser Slice: Punkt 3 (e) in der Fassung der Bereichsformen, Punkt 4, Punkt 6 (die Zeile auf stderr), Punkt 7 (Gründe und Selbstauskunft) und Punkt 8 (2), (3) und (5). `grep` nach der Kennung des Gebers findet die Sendung in §1 und DoD.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der **Wächter der ersten Stufe** (Hook, Leser, Eintrag, Gegenprobe-Gerüst, Nachweis im Agentenlauf) — Klasse 1: `slice-harness-rollenlauf-zeitgrenze` (Geber, vor diesem Slice in `done/`); dieser Slice erweitert ihn und ändert die Zusagen der ersten Stufe nicht.
- Die **Regel** in `AGENTS.md` und den sechs Rollen — Klasse 1: `slice-harness-rollenlauf-zeitgrenze-regel`. Grund: anderer Vorgang und andere Dateien.
- Die **Wiederholungs-Grenze** für volle Gegenprobe-Läufe — Klasse 1: `slice-harness-rollenlauf-wiederholungsgrenze`. Grund: ein zweiter Mechanismus mit eigenen Zusagen, der die Datenquelle des Wächters erbt.
- Ein **Anker für die Hauptsitzung**, eine **Zeit innerhalb eines Vordergrund-Befehls**, ein **Abbruch** des Laufs — Grenzen des Vertrags (Punkt 12) und Entscheidung des Nutzers; eine Adresse gibt es nicht.
- Die **Zeile der Gegenprobe in der Tabelle von [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16** — Vertragstext; sie steht seit der ersten Stufe dort (Architect, vor deren Closure), und dieser Slice ändert sie nicht.
- Der **Nachweis im Agentenlauf** (Punkt 13) — Klasse 1: Er gehört der ersten Stufe und gilt je Lieferung; ändert diese Stufe die Eingabe oder die Form des Blocks nicht, ist er nicht zu wiederholen.
- Die **Erfassungsschicht**, der **Stop-Hook** und die Gate-Ketten, Produkt-Code und Nutzer-Doku — Schicht-Abgrenzung (Klasse 4): Der Wächter liest nur und urteilt über Zeit, nicht über Gates; kein Verhalten des gebauten Binary ändert sich.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (der Wächter nennt seine Lage):** `.claude/hooks/pretooluse-zeitgrenze.sh` und `tools/harness/hook-felder.awk` erfüllen die Punkte der zweiten Stufe von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) (Zuordnung des Architect): strenge Zeitform des Ankers, Reihenfolge der Gründe, Gründe der Freigabe mit der Zeile auf stderr, Umgebung und Aufrufe (Sendung von `slice-harness-rollenlauf-zeitgrenze`, §1). Jede Randform dieser Punkte ist in `SPEC-070` entschieden, bevor die erste Zeile Code entsteht.
- [ ] **Liefer-Punkt 2 (die Gegenprobe der zweiten Stufe):** `make zeitgrenze-gegenprobe` trägt die Fälle und Mutationen der zweiten Stufe, **je Zusage und je Element ihrer Aufzählungen ein Fall und eine Mutation, die nur ihn rot färbt** (`AGENTS.md` §3.10; die Namen der Fälle tragen die Nummer des Punktes). Die unveränderte Kopie ist grün, **kein Fall wartet**, der volle Lauf liegt **unter 600 s** ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md); die gemessene Zeit steht in §7), und die Fälle der ersten Stufe bleiben unverändert grün.
- [ ] `make docs-check` und `make kopf-check` grün.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/`, `harness/` und `tools/` ändern).
- [ ] Die Gegenprobe ist am Endstand **in einem Zug** grün, mit ihrer Zeit in §7 ([`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16); während der Arbeit fährt der Implementer nur den engsten Fall mit `GP_MUT`.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`): ein weiterer Beleg in `BEO-REPO/rollenlauf-ohne-zeitgrenze` oder „keine Beobachtung angefallen“ in §7; die gemessene Zeit der erweiterten Gegenprobe als Beleg in `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der nächsten Welle-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md` | liegt vor; Zuordnung der Punkte zu den zwei Stufen **durch den Architect** vor dem Code | Zusagen und Randformen der zweiten Stufe entschieden (`AGENTS.md` §3.12) |
| `.claude/hooks/pretooluse-zeitgrenze.sh`, `tools/harness/hook-felder.awk` | update | Liefer-Punkt 1: Zeitform, Reihenfolge der Gründe, Gründe mit der Zeile auf stderr |
| `tools/harness/zeitgrenze-gegenprobe.sh` | update | Liefer-Punkt 2: Fälle und Mutationen der zweiten Stufe |
| `harness/mk/zeitgrenze.mk`, `harness/README.md` | nur gelesen, höchstens eine Zeile | das Ziel und seine Zeile bestehen seit der ersten Stufe |
| `docs/plan/planning/observations/BEO-REPO/rollenlauf-ohne-zeitgrenze/` und `…/gate-laufzeit-waechst-mit-gegenprobe/` | update bei Closure | Beleg, gemessene Zeit (kein Liefer-Punkt) |

- **Reihenfolge:** (1) der Architect hält die Zuordnung in `SPEC-070` fest; (2) der Wächter nennt seine Lage; (3) die Gegenprobe, je Zusage eine Mutation, vom Implementer selbst rot gesehen; (4) Closure mit Register. Jeder Commit nennt die Kennung dieses Slice.
- **Was der Implementer nicht entscheidet:** Trifft er eine Randform, die `SPEC-070` nicht nennt, hält er an und gibt sie dem Architect zurück (`AGENTS.md` §3.12). Den Vertragstext ändert er nicht.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): WIP-Limit frei; `slice-harness-rollenlauf-zeitgrenze` liegt in `done/` (Hook, Leser, Gegenprobe-Gerüst und Aufnahme stehen); der Architect hat die Zuordnung der Punkte von `SPEC-070` zu den zwei Stufen festgehalten und jede Randform der zweiten Stufe entschieden (`AGENTS.md` §3.12); [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) ist vom Nutzer angenommen (2026-10-22; `Accepted` in der Statuszeile mit dem Merge). Reihenfolge: direkt hinter `slice-harness-rollenlauf-zeitgrenze` (**Empfehlung des Planners**; beide ändern Hook und Gegenprobe). Der Übergang `open` → `next` setzt `Verantwortlich:`.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (der Befehl steht in `AGENTS.md` §3.14 nach `slice-harness-groessen-metrik`): Summe aus Einfügungen und Löschungen über die Pfade `AGENTS.md`, `.claude`, `harness` und `tools`, mit den Ausschlüssen des Befehls; Basis ist der Commit, der diese Datei nach `in-progress/` bewegt. **Marke 700, Meldeschwelle 550:** Erreicht die Summe 550, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet. Schätzung des Planners (2026-10-22, in dieser Metrik, mit dem Faktor 3 bis 4,6 für Gegenprobe-Code gegenüber Produktcode): Wächter 40 bis 90, Gegenprobe 233 bis 403 (53 der 125 Fälle der Schätzung des Architect von 550 bis 950), zusammen **273 bis 493, das sind 39 bis 70 % der Marke** (Mitte 383, 55 %).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über 700, oder die Zuordnung des Architect verlangt einen dritten Liefer-Punkt.
- `in-progress` → `open` (blockiert): Die Zuordnung der Punkte zu den Stufen liegt nicht im Vertrag oder die erste Stufe liegt nicht in `done/`; oder der Nutzer nimmt die Entscheidung vom 2026-10-21 zurück.

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

**Randformen:** Der Slice liefert die zweite Stufe eines Vertrags, den der Architect entschieden hat ([`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md)); die Randformen der Punkte 3 (e), 4, 6, 7 und 8 stehen dort, jede mit ihrer Antwort, und der Plan wiederholt sie nicht. **Offen für den Architect**, vor dem Code: die Zuordnung der Punkte zu den zwei Stufen (der Vorschlag steht in §6 von `slice-harness-rollenlauf-zeitgrenze`, Randform 7). Der Implementer entscheidet keine Randform, die der Vertrag nicht nennt (`AGENTS.md` §3.12).

**Risiken:**

- **Der Wächter der ersten Stufe schweigt, wo der Strom fehlt, bis diese Stufe liefert** (kein Grund auf stderr) — **Ausgang:** weiter offen → `BEO-REPO/rollenlauf-ohne-zeitgrenze`; Gegenmittel: Dieser Slice steht direkt hinter der ersten Stufe, in derselben Welle.
- **Die Zuordnung der Architect hält sich nicht an die Naht** (ein Punkt der zweiten Stufe liegt in der ersten oder umgekehrt) und verschiebt Fälle zwischen den Slices — **Ausgang:** offen bis Closure (entfallen, wenn Plan und Zuordnung dieselben Punkte nennen; sonst zieht der Planner beide Pläne nach, `AGENTS.md` §3.9).
- **Die Gegenprobe wächst mit den Fällen über das Zeitbudget** (600 s, [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)) — **Ausgang:** offen bis Closure (die gemessene Zeit steht in §7; entfallen, wenn sie unter 600 s liegt).
- **Die Fälle der zweiten Stufe verlangen Wartezeiten** (Zeitgrenze des Hooks, gleichzeitige Aufrufe) — **Ausgang:** entfallen, wenn kein Fall wartet (die Zeit kommt aus Anker und Attrappe, nicht aus `sleep`; `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test`).
- **Größe** (Schätzung 273 bis 493 gegen Marke 700) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen: *Belege des Implementers* (Mutationstabelle je Zusage der zweiten Stufe: Zusage · Mutation · roter Fall), die Zeit des vollen Gegenprobe-Laufs, die Größe in der Metrik von §4 mit dem Ergebnis der Schätzung gegen die Messung, und der Beleg des Registers.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: `BEO-REPO/rollenlauf-ohne-zeitgrenze` (1×) — der Auslöser; `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (14×, `verkörpert`) — die erweiterte Gegenprobe hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS`, die Kette wächst nicht; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/spec-randform-erst-im-review-entschieden` (beide `verkörpert`) — je Zusage ein Fall und eine Mutation, Randformen vor dem Code; `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (10×, `geplant`) — die Fälle der Umgebung (Punkt 8). Eine Schwelle wird mit diesem Slice nicht neu erreicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-22, im Commit, der diesen Plan anlegt und die Sendung von `slice-harness-rollenlauf-zeitgrenze` einträgt):** Liefer-Punkte 2 (der Wächter nennt seine Lage | die Gegenprobe der zweiten Stufe). Schichten 1 (Harness: Hook, Leser, Gegenprobe; die Planung zählt keine Schicht, Vertrag und Entscheidung liegen vor der Basis). Schätzung 273 bis 493 in der Metrik von §4 mit dem Faktor 3 bis 4,6 für Gegenprobe-Code (Marke 700, Meldeschwelle 550). Die Sendung legt keinen Liefer-Punkt und keine Schicht an.
