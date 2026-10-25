# Slice slice-harness-rollenlauf-zeitgrenze-haertung: Der Wächter der Stunde nennt seine Lage und ist gehärtet: strenge Zeitform des Ankers, Reihenfolge der Gründe, Gründe der Freigabe mit der Zeile auf stderr und die Fälle der Umgebung — Stufe 2

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-harness-regeln, **direkt hinter** `slice-harness-rollenlauf-zeitgrenze` (§4 und §5 der Welle). **Die Stelle in der Welle ist Empfehlung des Planners vom 2026-10-22**, nicht Entscheidung des Nutzers; die Zuordnung der Punkte zu den zwei Stufen hat der Architect im Vertrag festgehalten. Entschieden hat der Nutzer die **Sache**: die Grenze von einer Stunde je Rollenlauf mit einem Wächter im Werkzeug-Aufruf (2026-10-21, wörtlich: „Wir haben rollen (siehe .claude/agents) und commands -- max. 1 Stunde“, dazu sein „ok“ auf den Entwurf des Mechanismus). **Geschnitten am 2026-10-22** aus `slice-harness-rollenlauf-zeitgrenze`: Dessen Schätzung (Hook und Gegenprobe zusammen 703 bis 1212 in der Metrik von [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md)) lag über der Marke; der Architect hat den Schnitt nach **Wirkung** und **Härtung und Selbstauskunft** bestätigt und die Punkte und Fälle den zwei Stufen zugeordnet ([`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md), Abschnitt *Stufen*, 2026-10-22); dieser Slice ist die zweite Stufe. Angelegt in `open/`, weil er nicht priorisiert ist und hinter der ersten Stufe steht (§4).

**Bezug:** Kein Anspruch des Produkts (`LH-*`) — der Slice härtet einen Wächter über Agentenläufe. Bindung an das Register: `BEO-REPO/rollenlauf-ohne-zeitgrenze` (1×), `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test` (kein Fall wartet). Aktive ADRs: [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) (Accepted; Entscheidung und Gründe des Wächters), [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Zeitbudget der Gegenprobe, 1200 s nach [ADR-0059](../../adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md), gilt auch für die erweiterte Gegenprobe), [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenprobe hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS`).

**Berührte Spec-Stellen:** [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) (die Punkte und Fälle der zweiten Stufe nach dem Abschnitt *Stufen*) · [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) (gelesen: Zeitbudget, Aufrufarten, Signale der Gegenprobe). Der Slice ändert keinen Vertrag.

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

**Ziel:** Der Wächter, den `slice-harness-rollenlauf-zeitgrenze` liefert (Stufe 1: die Stunde wirkt, die Freigabe nur in den Lagen, die ein Fall zeigt, über stderr nichts zugesagt), nennt seine Lage und hält seine Zusagen an den Rändern — die Punkte und Fälle der zweiten Stufe nach dem Abschnitt *Stufen* von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) (rund 29 Fälle): (1) der **Bereich des Ankers** (Punkt 3 (e): `p3-form-monat-13`, `-tag-32`, `-februar-30`, `-stunde-24`, `-minute-60`, `-sekunde-60`), (2) die **Reihenfolge der Gründe** (Punkt 4, acht Fälle), (3) die **Zeile auf stderr** (Punkt 6: `p6-exit-0-nicht-geprueft`, `p6-stderr-eine-zeile`; je Lage von Punkt 3, 5 und 7 genau eine Zeile `zeitgrenze: nicht geprüft (<Grund>)`, auch `befehl-unlesbar` bei `p5-bash-unlesbar`) und die **übrigen Fälle der Freigabe** (Punkt 7: `p7-eingabe-abgeschnitten`, `-eingabe-doppelt`, `-eingabe-unicode`, `-kennung-session-leer`, `-kennung-session-fehlt`, `-kennung-lang`, `-kennung-pfad`, `-anker-unlesbar`, `-kein-date`, `-befehl-unlesbar`), (4) die **Fälle der Umgebung** (Punkt 8 (2), (3), (5): `p8-leere-umgebung`, `p8-gross`, `p8-parallel`). Dazu die Erweiterung der Gegenprobe im selben Ziel `make zeitgrenze-gegenprobe`, unter 1200 s (Punkt 14). Nach dieser Stufe gelten die Sätze „Jede Unsicherheit des Hooks ist eine Freigabe“ und „Exit 0 in jedem Zweig“ uneingeschränkt (Vertrag, *Stufen*).

**Warum eine zweite Stufe und kein Anhang der ersten:** Die erste Stufe wirkt und gibt in den Lagen frei, die ein Fall zeigt, ohne eine Zeile auf stderr; die Selbstauskunft ist die Gegenmaßnahme dafür, dass der Hook dort schweigt, wo er gebraucht wird ([ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md), Punkt 4). Hook und Gegenprobe mit allen Zusagen lagen bei 703 bis 1212, das sind 106 % der Marke 900; der Architect hat die Zuordnung der Zusagen zu zwei Lieferungen im Vertrag festgehalten (Abschnitt *Stufen*, auf den Vorschlag des Planners, mit den Änderungen, die der Vertrag dort nennt; die Freigabe in den sieben Lagen liegt in Stufe 1, der Bereich des Ankers in Stufe 2).

**Sendung von `slice-harness-rollenlauf-zeitgrenze`** (eingetragen vom Planner im Commit des Schnitts, `AGENTS.md` §3.13; nachgezogen am 2026-10-22 nach der Zuordnung des Architect): Die Punkte und Fälle der zweiten Stufe aus dem Abschnitt *Stufen* des Vertrags sind dieser Slice: der Bereich von Punkt 3 (e), Punkt 4, die Zeile auf stderr aus Punkt 6, von Punkt 7 **nur noch die übrigen Fälle** und die Zeile auf stderr in jeder Lage, Punkt 8 (2), (3) und (5). **Die Freigabe-Fälle der sieben Lagen** (`p7-kein-verzeichnis`, `-kein-strom`, `-strom-leer`, `-kein-start`, `-eingabe-leer`, `-eingabe-kein-json`, `-kein-awk`) **gehören nach Stufe 1**, nicht hierher. `grep` nach der Kennung des Gebers findet die Sendung in §1 und DoD.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der **Wächter der ersten Stufe** (Hook, Leser, Eintrag, Gegenprobe-Gerüst, die sieben Fälle der Freigabe, Nachweis im Agentenlauf) — Klasse 1: `slice-harness-rollenlauf-zeitgrenze` (Geber, vor diesem Slice in `done/`); dieser Slice erweitert ihn und ändert die Zusagen der ersten Stufe nicht.
- Die **Regel** in `AGENTS.md` und den sechs Rollen — Klasse 1: `slice-harness-rollenlauf-zeitgrenze-regel`. Grund: anderer Vorgang und andere Dateien.
- Die **Wiederholungs-Grenze** für volle Gegenprobe-Läufe — Klasse 1: `slice-harness-rollenlauf-wiederholungsgrenze`. Grund: ein zweiter Mechanismus mit eigenen Zusagen, der die Datenquelle des Wächters erbt.
- Ein **Anker für die Hauptsitzung**, eine **Zeit innerhalb eines Vordergrund-Befehls**, ein **Abbruch** des Laufs — Grenzen des Vertrags (Punkt 12) und Entscheidung des Nutzers; eine Adresse gibt es nicht.
- Die **Zeile der Gegenprobe in der Tabelle von [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16** — Vertragstext; sie steht seit der ersten Stufe dort (Architect, vor deren Closure), und dieser Slice ändert sie nicht.
- Der **Nachweis im Agentenlauf** (Punkt 13) — Klasse 1: Er gehört der ersten Stufe und gilt je Lieferung; diese Stufe wiederholt ihn nur, wenn sie den Eintrag, den Pfad von `agent_id` oder die Ausgabe der Sperre ändert (Vertrag, Punkt 13).
- Die **Erfassungsschicht**, der **Stop-Hook** und die Gate-Ketten, Produkt-Code und Nutzer-Doku — Schicht-Abgrenzung (Klasse 4): Der Wächter liest nur und urteilt über Zeit, nicht über Gates; kein Verhalten des gebauten Binary ändert sich.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (der Wächter nennt seine Lage):** `.claude/hooks/pretooluse-zeitgrenze.sh` und `tools/harness/hook-felder.awk` erfüllen die Punkte der zweiten Stufe nach dem Abschnitt *Stufen* von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md): den Bereich des Ankers (Punkt 3 (e)), die Reihenfolge der Gründe (Punkt 4), die Zeile auf stderr in jeder Lage (Punkte 6 und 7, mit „Exit 0 in jedem Zweig“ uneingeschränkt), die übrigen Fälle der Freigabe (Punkt 7), Umgebung und Aufrufe (Punkt 8 (2), (3), (5)) (Sendung von `slice-harness-rollenlauf-zeitgrenze`, §1). **Stufe 2 lockert keinen Fall der Stufe 1.** **Jeder Satz über den Hook** — die Zeile des Wächters in `harness/README.md`, der Kommentar im Kopf des Hooks und des Fragments — darf nach dieser Stufe die Selbstauskunft auf stderr und die allgemeine Freigabe nennen; der Implementer zieht sie nach, und er nennt nicht mehr (Vertrag, *Zwischen den Stufen* (d); `AGENTS.md` §3.11). Jede Randform dieser Punkte hat der Vertrag entschieden.
- [ ] **Liefer-Punkt 2 (die Gegenprobe der zweiten Stufe):** `make zeitgrenze-gegenprobe` trägt die Fälle und Mutationen der zweiten Stufe (rund 29 nach dem Vertrag), **je Zusage und je Element ihrer Aufzählungen ein Fall und eine Mutation, die nur ihn rot färbt** (`AGENTS.md` §3.10; die Namen der Fälle tragen die Nummer des Punktes). **Die Erwartung ist getrennt je Fall:** Stufe 2 ergänzt Fälle der Stufe 1 um die Erwartung an stderr (die Zeile) — in jeder Lage von Punkt 3, 5 und 7, darunter `p5-bash-unlesbar`, die sieben `p3-form-*` der Gestalt und die sieben `p7`-Fälle der Freigabe —, lockert keinen und **ändert weder Namen noch erwartete Mengen roter Fälle der Stufe 1** (Vertrag, *Zwischen den Stufen* (b)); der Beleg in §7 vergleicht die Mengen vor und nach. Die unveränderte Kopie ist grün, **kein Fall wartet**, der volle Lauf liegt **unter 1200 s**, auch nach dieser Stufe ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md), Zahl nach [ADR-0059](../../adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md); die gemessene Zeit steht in §7), und die Fälle der ersten Stufe bleiben grün.
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
| `harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md` | liegt vor; die Zuordnung der Punkte und Fälle zu den zwei Stufen steht im Abschnitt *Stufen* | Zusagen und Randformen der zweiten Stufe entschieden (`AGENTS.md` §3.12) |
| `.claude/hooks/pretooluse-zeitgrenze.sh`, `tools/harness/hook-felder.awk` | update | Liefer-Punkt 1: Bereich des Ankers, Reihenfolge der Gründe, Zeile auf stderr, übrige Fälle der Freigabe, Umgebung |
| `tools/dev/harness/zeitgrenze-gegenprobe.sh` | update | Liefer-Punkt 2: Fälle und Mutationen der zweiten Stufe |
| `harness/mk/zeitgrenze.mk` (Kommentar im Kopf), `harness/README.md` (Zeile des Wächters) | update, höchstens wenige Zeilen | das Ziel und seine Zeile bestehen seit der ersten Stufe; Sätze über den Hook nennen nach dieser Stufe die Selbstauskunft und die allgemeine Freigabe (DoD) |
| `docs/plan/planning/observations/BEO-REPO/rollenlauf-ohne-zeitgrenze/` und `…/gate-laufzeit-waechst-mit-gegenprobe/` | update bei Closure | Beleg, gemessene Zeit (kein Liefer-Punkt) |

- **Reihenfolge:** (1) die Zuordnung steht im Vertrag (erfüllt); (2) der Wächter nennt seine Lage; (3) die Gegenprobe, je Zusage eine Mutation, vom Implementer selbst rot gesehen; (4) Closure mit Register. Jeder Commit nennt die Kennung dieses Slice.
- **Was der Implementer nicht entscheidet:** Trifft er eine Randform, die `SPEC-070` nicht nennt, hält er an und gibt sie dem Architect zurück (`AGENTS.md` §3.12). Den Vertragstext ändert er nicht.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): WIP-Limit frei; `slice-harness-rollenlauf-zeitgrenze` liegt in `done/` (Hook, Leser, Gegenprobe-Gerüst und Aufnahme stehen); die Zuordnung der Punkte von `SPEC-070` zu den zwei Stufen steht im Vertrag und jede Randform der zweiten Stufe ist entschieden (`AGENTS.md` §3.12); [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) ist vom Nutzer angenommen (2026-10-22; `Accepted`). Reihenfolge: direkt hinter `slice-harness-rollenlauf-zeitgrenze` (**Empfehlung des Planners**; beide ändern Hook und Gegenprobe). Der Übergang `open` → `next` setzt `Verantwortlich:`.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (der Befehl steht in `AGENTS.md` §3.14 nach `slice-harness-groessen-metrik`): Summe aus Einfügungen und Löschungen über die Pfade `AGENTS.md`, `.claude`, `harness` und `tools`, mit den Ausschlüssen des Befehls; Basis ist der Commit, der diese Datei nach `in-progress/` bewegt. **Marke 500, Meldeschwelle 400:** Erreicht die Summe 400, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet. Schätzung des Planners (2026-10-22, nachgerechnet nach der Zuordnung der Stufen, in dieser Metrik, mit dem Faktor 3 bis 4,6 für Gegenprobe-Code gegenüber Produktcode): Wächter 50 bis 100 (Bereich, Reihenfolge der Gründe, die Zeile je Lage, Umgebung), `harness/README.md` und Kommentar 0 bis 5, Gegenprobe 150 bis 300 (29 Fälle zu 4,4 bis 7,6 Zeilen aus der Schätzung des Architect, dazu die Erwartung an stderr in den Fällen der Stufe 1), zusammen **200 bis 405, das sind 40 bis 81 % der Marke 500** (Mitte 302, 60 %); die Obergrenze liegt an der Meldeschwelle. **Die Marke ist von 700 auf 500 gesenkt**, weil die Zuordnung des Architect die Stufe kleiner machte (vorher 53 Fälle, jetzt 29); eine Marke über der Schätzung wäre keine.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über 500, oder ein Fall der Stufe 2 verlangt, einen Fall der Stufe 1 zu lockern oder umzubenennen (dann gehen beide Stufen an den Architect zurück, Vertrag, *Zwischen den Stufen* (b) und (f)).
- `in-progress` → `open` (blockiert): Die erste Stufe liegt nicht in `done/`; oder der Nutzer nimmt die Entscheidung vom 2026-10-21 zurück.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` grün; die Gegenprobe am Endstand in einem Zug grün mit ihrer Zeit unter 1200 s in §7; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen:** Der Slice liefert die zweite Stufe eines Vertrags, den der Architect entschieden hat ([`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md)); die Randformen der Punkte 3 (e), 4, 6, 7 und 8 stehen dort, jede mit ihrer Antwort, und der Plan wiederholt sie nicht. Die Zuordnung der Punkte und Fälle zu den zwei Stufen steht im Vertrag (Abschnitt *Stufen*). **Getrennte Erwartung je Fall:** Ein Fall mit Grund und Zeile prüft in Stufe 1 die Freigabe, in Stufe 2 dazu die Zeile; Stufe 2 lockert nie, ändert keinen Namen und keine erwartete Menge roter Fälle der Stufe 1 ((b) dort). **Randformen, die ein Fall der Stufe 2 aufwirft und die Stufe 1 berühren,** gibt der Implementer an den Architect zurück; beide Stufen gehen dann zurück ((f) dort). Der Implementer entscheidet keine Randform, die der Vertrag nicht nennt (`AGENTS.md` §3.12).

**Risiken:**

- **Der Wächter der ersten Stufe schweigt, wo der Strom fehlt, bis diese Stufe liefert** (über stderr ist nach Stufe 1 nichts zugesagt) — **Ausgang:** weiter offen → `BEO-REPO/rollenlauf-ohne-zeitgrenze`; Gegenmittel: Dieser Slice steht direkt hinter der ersten Stufe, in derselben Welle.
- **Stufe 2 lockert oder benennt einen Fall der Stufe 1 um** (etwa weil eine Zeile auf stderr einen erwarteten roten Fall verschiebt): Der Hook der ersten Stufe wäre nach der zweiten schwächer geprüft als vorher. — **Ausgang:** offen bis Closure (entfallen, wenn der Vergleich der Mengen roter Fälle der Stufe 1 vor und nach der Stufe in §7 gleich ist; sonst *eingetreten* → Rückgabe an den Architect, Vertrag (b) und (f)).
- **Ein Satz über den Hook bleibt nach Stufe 2 auf dem Stand der ersten Stufe oder sagt mehr zu, als geliefert ist** (`harness/README.md`, Kommentar im Kopf). — **Ausgang:** entfallen, wenn das Review die Sätze gegen den Vertrag liest (*Nach Stufe 1 allein ist zugesagt*, und die Sätze „jede Unsicherheit“ und „Exit 0 in jedem Zweig“ uneingeschränkt erst jetzt).
- **Die Gegenprobe wächst mit den Fällen über das Zeitbudget** (600 s, [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)) — **Ausgang:** offen bis Closure (die gemessene Zeit steht in §7; entfallen, wenn sie unter 600 s liegt).
- **Die Fälle der zweiten Stufe verlangen Wartezeiten** (Zeitgrenze des Hooks, gleichzeitige Aufrufe) — **Ausgang:** entfallen, wenn kein Fall wartet (die Zeit kommt aus Anker und Attrappe, nicht aus `sleep`; `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test`).
- **Größe** (Schätzung 200 bis 405 gegen Marke 500) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen: *Belege des Implementers* (Mutationstabelle je Zusage der zweiten Stufe: Zusage · Mutation · roter Fall), der Vergleich der erwarteten Mengen roter Fälle der Stufe 1 vor und nach dieser Stufe, die Zeit des vollen Gegenprobe-Laufs, die Größe in der Metrik von §4 mit dem Ergebnis der Schätzung gegen die Messung, und der Beleg des Registers.

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

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-22, im Commit, der diesen Plan anlegt und die Sendung von `slice-harness-rollenlauf-zeitgrenze` einträgt; nachgezählt am 2026-10-22 nach der Zuordnung der Stufen im Vertrag):** Liefer-Punkte 2 (der Wächter nennt seine Lage | die Gegenprobe der zweiten Stufe). Schichten 1 (Harness: Hook, Leser, Gegenprobe, eine Zeile `harness/README.md`; die Planung zählt keine Schicht, Vertrag und Entscheidung liegen vor der Basis). Schätzung **200 bis 405** in der Metrik von §4 mit dem Faktor 3 bis 4,6 für Gegenprobe-Code, rund 29 Fälle nach dem Vertrag (Marke 500, Meldeschwelle 400; vorher 273 bis 493 gegen Marke 700 bei 53 Fällen). Die Sendung (nur noch die übrigen Fälle von Punkt 7 und die Zeile auf stderr) legt keinen Liefer-Punkt und keine Schicht an.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Ortsangabe aus `slice-harness-gegenproben-umzug-dev` einträgt): Die Gegenproben liegen seit jenem Slice unter `tools/dev/<bereich>/`, die Pfade der Gegenproben in diesem Plan sind nachgezogen (`tools/dev/` ist wie `tools/harness/` Harness-Werkzeug, keine neue Schicht). Ortsangabe der Gegenprobe geändert, Liefer-Punkte (2) und Schichten (1) unverändert.
