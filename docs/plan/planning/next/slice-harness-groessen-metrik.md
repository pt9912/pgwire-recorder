# Slice slice-harness-groessen-metrik: `AGENTS.md` §3.14 trägt den berichtigten Befehl der Größen-Metrik, und die Entscheidung dazu wird angenommen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-harness-regeln (**Empfehlung des Planners vom 2026-10-21, vom Nutzer nicht entschieden**; die Welle führt die fünf Regel-Slices des Harness, und ihr Closure-Trigger verlangt `make docs-check` und `make gates` grün auf dem Endstand mit allen Regeln eingetragen). Angelegt am 2026-10-21 nach Entscheidung des Nutzers vom selben Tag (wörtlich, vom Koordinator übermittelt: „ok deine Empfehlungen“ auf eine Vorlage mit drei Punkten, deren zweiter die Annahme von [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) und dieser Slice war). Eingeordnet in `next/`, weil er mit einer Stelle in der Reihenfolge eingeplant ist (§4 *Start*, Drift-Log der Roadmap).

**Bezug:** kein Anspruch des Produkts — der Slice ändert eine Regel des Harness und nimmt eine Entscheidung an. Entscheidung: [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (Empfehlung des Architect vom 2026-10-21, Status `Proposed`; dieser Slice setzt sie auf `Accepted`). Bindung an das Register: `BEO-REPO/groessenprobe-ohne-messbaren-stand` (5×, Stand verkörpert), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; `BEO-REPO/adr-ohne-zeile-im-index` (vom Planner am 2026-10-21 angelegt, Befund des Nutzers vom selben Tag; *benannt, nicht gezählt*, §6). Regelstelle: `AGENTS.md` §3.14 (seit slice-harness-gate-kette-stop-hook-liste).

**Berührte Spec-Stellen:** — *(keine; die Formel steht in `AGENTS.md`, nicht in einem Vertrag — [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) Punkt 6)*

**Verantwortlich:** pt9912

**Autor:** Planner. **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Zwei Leser erhalten für dieselbe Arbeit dieselbe Zahl: `AGENTS.md` §3.14 trägt den **berichtigten Befehl** der Größen-Metrik samt Basis, Messpunkt, den zwei Ausschlüssen und dem Rang-Zeiger („§4 eines Plans nennt Pfade, Marke und Meldeschwelle und zeigt auf diese Stelle; wiederholt ein Plan den Befehl, gilt der von hier“), und [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) wird mit Status und Index angenommen. Die Schwellen selbst (Marken und Meldeschwellen der Pläne) bleiben unverändert; der Slice ändert, **wie** gemessen wird, nicht **wie viel** erlaubt ist.

**Was der Befehl berichtigt**, in den Worten der Entscheidung ([ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) Punkte 1 bis 7): Der Zustandswechsel eines Slice zählte als Arbeit (die Plandatei stand am Endstand zweimal im Diff, 671 Zeilen für eine Datei), und die Formel zählte mit, was die Größenregel des Baseline-Regelwerks ausnimmt (Review- und Verifikationsbericht, Closure-Notiz, Register). Gemessen an der Closure von `slice-erster-release-homebrew-nachweis`: 1148 nach der alten, 244 nach der neuen Formel, und die neue Zahl steht über Review, Verifikation, Nacharbeit, Closure und Move still (223 bis 244).

**Bindung an `slice-harness-schaetzung-planner-faktor` (Sendung dieses Slice, einzutragen im selben Commit, `AGENTS.md` §3.13):** Jener Slice ergänzt die **Schätz-Hälfte** von §3.14 (Faktor Tests zu Produktcode, Mindestwert 5, Schnitt ab 70 % der Marke, Meldeschwelle in §4, datierte Entscheidung nach einer Meldung, Zählung nach der Obergrenze). Seine Definition of Done verlangt einen Diff **ohne Löschzeile in der Regel selbst** und nimmt diese Berichtigung darum nicht an, solange sie nicht eingetragen ist. Dieser Slice läuft deshalb **vor** ihm und trägt die Bindung in dessen §1 und dessen DoD ein, mit der Kennung dieses Slice, und zählt dessen Liefer-Punkte und Schichten im selben Commit nach (dort §8). Erst danach ist jener Slice startbar, ohne gegen seine eigene DoD zu arbeiten.

**Wortlaut und Entscheidung des Nutzers.** Der Entwurf des Wortlauts steht in §3. Eine Regeländerung in `AGENTS.md` ist eine Entscheidung des Nutzers: Der Implementer legt den Wortlaut vor und committet erst nach der Bestätigung; bis dahin ist er eine Empfehlung. Entschieden hat der Nutzer am 2026-10-21 die **Sache** (Annahme von [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) mit ihrem Befehl, Liefer-Punkte, Marke und Stelle in der Reihenfolge), nicht den Satzbau in §3.14; Datum und Fundstelle der Bestätigung stehen in §7.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die **Schätz-Hälfte** von §3.14 (Faktor, Mindestwert 5, Schnitt ab 70 %, Meldeschwelle in §4, Entscheidung nach einer Meldung, Zählung nach der Obergrenze) — `slice-harness-schaetzung-planner-faktor` (Klasse 1, `next/`; nimmt die Bindung in §1 und DoD an, siehe oben). Grund: Dort liegt der Liefergegenstand schon ausgeschrieben, und zwei Slices am selben Absatz erzeugen einen Konflikt statt einer Regel.
- Ein **Sensor**, der die Form der Marke in §4 eines Plans prüft (Pfade, Marke, Meldeschwelle, Zeiger auf §3.14) — anderer Vorgang (Klasse 3: Werkzeug statt Regel), mit Vertrag in `harness/vertraege/` und Gegenprobe. Grund: [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) *Fitness Function* nennt ihn urteilsfrei prüfbar, aber erst fällig, wenn `BEO-REPO/groessenprobe-ohne-messbaren-stand` **trotz** der Berichtigung ein weiteres Auftreten zählt; ein Sensor auf einer Formel, die gerade berichtigt wird, prüft den alten Stand. Eine Adresse gibt es nicht; der Planner legt den Slice dann an.
- Ein **Sensor für die Zeile einer ADR im Index** (`BEO-REPO/adr-ohne-zeile-im-index`, §6) — anderer Vorgang (Klasse 3), und **Entscheidung des Nutzers**, nicht dieses Slice: ob der Befund eine Prüfung bekommt und wo sie hängt, ist offen (§6, Meldung des Planners). Grund: Der Slice trägt die Zeile für [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) selbst nach, was eine Prüfung nicht ersetzt; ein Gate zu `.d-check.yml` oder ein eigenes Skript ist Arbeit am Werkzeug.
- Das **Nachziehen des Bestands** der lebenden Pläne auf den berichtigten Befehl — Bestand bleibt bewusst stehen (Klasse 2), [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) Punkt 7. Grund: Der neue Befehl liefert für jeden lebenden Plan eine Zahl, die nicht über der alten liegt; keine Marke greift dadurch früher, also ist kein Plan rückwirkend zu groß. Ein Plan, der die alte Zeile trägt, wird mit dem Befehl aus §3.14 gelesen und beim nächsten Anfassen nachgezogen.
- Eine **Änderung der Marken oder Meldeschwellen** irgendeines Plans — Klasse 3 (anderer Vorgang) und `AGENTS.md` §3.6. Grund: Die Zahlen sind in der Summe gesetzt und bleiben; eine Senkung wäre eine Lockerung mit eigener Entscheidung.
- **Produkt-Code, Gates, Gegenproben und Verträge** (`cmd/`, `internal/`, `tools/`, `harness/mk/`, `harness/vertraege/`) — Schicht-Abgrenzung (Klasse 4). Grund: Der Slice ändert `AGENTS.md`, eine ADR und den ADR-Index; kein Werkzeug und kein Gate verhält sich danach anders. Damit ist beim Review sofort prüfbar, was er nicht angefasst hat.
- `.claude/commands/plan-welle.md` und die Rollen-Typen unter `.claude/agents/` — Klasse 1, `slice-harness-schaetzung-planner-faktor` (dessen §3 führt `plan-welle.md` als Liefer-Punkt). Grund: §3.14 ist der eine Träger der Formel ([ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) Punkt 6); ein zweiter Ort für denselben Befehl driftet gegen ihn.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die Closure-Pflichten darunter zählen nicht mit.

Liefer-Punkte (zwei):

- [ ] **Liefer-Punkt 1 — `AGENTS.md` §3.14 trägt den berichtigten Befehl:** den Befehl mit `--find-renames=10%` und den zwei Ausschlüssen (`':(exclude)docs/plan/planning'`, `':(exclude)docs/reviews'`), die **Basis** (der Commit, der die Plandatei nach `in-progress/` bewegt; seinen Hash trägt §4 beim Übergang ein), den **Messpunkt** (der committete Stand des Zweigs zu jeder Zeit der Arbeit, nicht erst der letzte Commit), die Lesart der **Pfade** (die des Plans; die zwei Ausschlüsse stehen in jedem Fall im Befehl) und den **Rang-Zeiger** (§4 eines Plans nennt Pfade, Marke und Meldeschwelle und zeigt auf §3.14; wiederholt ein Plan den Befehl, gilt der von dort). Dazu ein Paar Falsch/Richtig und der Retirement-Check mit `seit slice-harness-groessen-metrik`. **Im selben Commit** (`AGENTS.md` §3.13) trägt der Slice die Bindung in §1 und die DoD von `slice-harness-schaetzung-planner-faktor` ein, mit der Kennung dieses Slice, und zählt dort in §8 dessen Liefer-Punkte und Schichten nach; `grep` nach `slice-harness-groessen-metrik` findet beide Stellen in jenem Plan. Der Wortlaut ist vom Nutzer bestätigt, **bevor** der Commit entsteht; Datum und Fundstelle stehen in §7.
- [ ] **Liefer-Punkt 2 — [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) ist angenommen:** Status `Accepted`, die Autor-Zeile nennt die Entscheidung des Nutzers vom 2026-10-21 mit ihrem Wortlaut, der Satz der Entscheidung steht nicht mehr als Empfehlung, und die **Geschichte**-Tabelle der ADR trägt die Zeile der Annahme. Der ADR-Index (`docs/plan/adr/README.md`) führt `0056` in der **Haupttabelle** mit Status `Accepted` **und** in der Tabelle der **Beziehungen**, beide in Nummernfolge (`BEO-REPO/adr-ohne-zeile-im-index`, §6).

Gates und Closure:

- [ ] `make docs-check` grün (der Slice ändert nur Dokumente; der Stop-Hook gibt den Pfad `docs/plan/planning/` frei, `AGENTS.md` und `docs/plan/adr/` nicht — der volle Lauf `make gates` ist deshalb nötig).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert. Für `BEO-REPO/adr-ohne-zeile-im-index` gilt: Trägt dieser Slice die Zeile für die ADR in **beide** Tabellen des Index ein (Liefer-Punkt 2), ist das **kein** Auftreten und erzeugt keinen Beleg; §7 vermerkt den Gegenbeleg. Fehlt eine der beiden Zeilen im gelieferten Stand, ist es ein Auftreten und bekommt `evidence/slice-harness-groessen-metrik.md`.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.14 | update | Liefer-Punkt 1: der berichtigte Befehl, Basis, Messpunkt, Pfade mit den zwei Ausschlüssen, Rang-Zeiger, Paar Falsch/Richtig, Retirement-Check |
| [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) | update | Liefer-Punkt 2: Status `Accepted`, Autor-Zeile mit der Entscheidung des Nutzers, Entscheidungssatz, Geschichte-Zeile |
| `docs/plan/adr/README.md` | update | Liefer-Punkt 2: Zeile `0056` in Haupttabelle (Status) und Beziehungen, Nummernfolge |
| `docs/plan/planning/next/slice-harness-schaetzung-planner-faktor.md` | update | Sendung dieses Slice (`AGENTS.md` §3.13): Bindung in §1 und DoD mit der Kennung dieses Slice, Nachzählung in dessen §8 — im selben Commit wie Liefer-Punkt 1 |
| `docs/plan/planning/observations/BEO-REPO/adr-ohne-zeile-im-index/` | update | bei der Closure: Gegenbeleg in `state.md`, wenn beide Index-Zeilen stehen, sonst ein Beleg in `evidence/` (kein Liefer-Punkt, DoD) |

**Entwurf des Wortlauts für `AGENTS.md` §3.14** (Empfehlung des Planners; der Nutzer bestätigt ihn, §1). Der erste Absatz der Regel wird ersetzt, die Schätz-Hälfte („Der Architect gibt Schätzungen in dieser Metrik an …“) bleibt Wort für Wort stehen:

> Jede Rückführungs-Marke in §4 eines Plans („zu groß, über N Zeilen mit Tests“) nennt ihre Metrik. Gemessen wird mit
>
> ```
> git diff --numstat --find-renames=10% <Basis>..HEAD -- <Pfade des Plans> \
>   ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \
>   | awk '{a+=$1; b+=$2} END {print a+b+0}'
> ```
>
> **Basis** ist der Commit, der die Plandatei nach `in-progress/` bewegt; seinen Hash trägt §4 beim Übergang ein. Was davor liegt — Vertrag und Entscheidung des Architect vor dem Start —, zählt nicht. **Messpunkt** ist der committete Stand des Zweigs zu jeder Zeit der Arbeit, nicht erst der letzte Commit; nicht gezählt ist, was noch nicht committet ist. **Pfade** sind die des Plans, Tests eingeschlossen; die zwei Ausschlüsse stehen in jedem Fall im Befehl, damit eine weite Pfadangabe (`docs`, `.`) den Planungsbaum und die Reports nicht wieder einfängt. Nicht gezählt sind damit Plan, Closure-Notiz, Roadmap, Welle-Plan, Beobachtungs-Register, Review- und Verifikationsbericht; gezählt bleiben Entscheidungen (`docs/plan/adr/`, `docs/plan/carveouts/`) sowie Nutzer- und Wartungs-Doku. **Der Befehl steht hier, einmal:** §4 eines Plans nennt Pfade, Marke und Meldeschwelle und zeigt auf diese Stelle; wiederholt ein Plan den Befehl, gilt der von hier.

Dazu ein Paar: *Falsch:* §4 nennt „Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit>`“ über `tools harness spec docs README.md`; am Endstand ergibt dieselbe Arbeit 1148, weil die Plandatei nach der Closure zweimal im Diff steht (`464 0` unter `done/`, `0 207` unter `next/`) und Review- und Verifikationsbericht mit 179 Zeilen mitzählen. *Richtig:* Derselbe Stand ergibt mit dem Befehl oben 244, und die Zahl steht über Review, Verifikation, Nacharbeit, Closure und Move still. Begründung: `BEO-REPO/groessenprobe-ohne-messbaren-stand`; Retirement-Check: Ist die Beobachtung seit slice-harness-groessen-metrik wieder aufgetreten?

- Der Implementer ändert **nichts** an der Schätz-Hälfte des Absatzes und an den Zahlen eines Plans; der Diff von `AGENTS.md` ersetzt den ersten Absatz der Regel und ergänzt ihn um Basis, Messpunkt, Pfade und Rang-Zeiger. Dass die Ersetzung Löschzeilen enthält, ist beabsichtigt und der Grund für die Reihenfolge gegen `slice-harness-schaetzung-planner-faktor` (§1).
- Reihenfolge der Commits: (1) Liefer-Punkt 1 mit der Sendung in den Nehmer (ein Commit, `AGENTS.md` §3.13), (2) Liefer-Punkt 2. Jeder Commit nennt eine zugelassene Kennung (die Kennung dieses Slice oder die Nummer der ADR).
- Zählung nach `AGENTS.md` §3.13 siehe §8.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): WIP-Limit frei, und der Nutzer hat den Wortlaut für §3.14 bestätigt oder beauftragt den Implementer, ihn vorzulegen (§1). Stelle in der Reihenfolge (**Empfehlung des Planners vom 2026-10-21, vom Nutzer nicht entschieden**): **vor** `slice-harness-schaetzung-planner-faktor`, weil jener Slice an derselben Regelstelle arbeitet und seine DoD einen Diff ohne Löschzeile verlangt (§1); im Übrigen innerhalb der Reihe von `welle-harness-regeln` (deren §5). Kein anderer Slice ist Vorbedingung; [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) liegt vor (`Proposed`) und wird von diesem Slice angenommen, ist also **kein** Start-Trigger.

**Größe.** Metrik in der Fassung, die dieser Slice liefert ([ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) Punkt 1; §3.14 trägt sie erst nach Liefer-Punkt 1, bis dahin gilt sie aus der ADR):

```
git diff --numstat --find-renames=10% <Basis>..HEAD -- AGENTS.md docs/plan/adr \
  ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \
  | awk '{a+=$1; b+=$2} END {print a+b+0}'
```

Basis ist der Commit, der diese Datei nach `in-progress/` bewegt (sein Hash steht beim Übergang hier). Die Pfade sind `AGENTS.md` und `docs/plan/adr`; die Sendung in `slice-harness-schaetzung-planner-faktor` und der Register-Beleg liegen unter `docs/plan/planning` und zählen nicht — das ist die Lesart, die der Slice liefert, nicht eine Ausnahme für ihn.

**Marke 200**, **Meldeschwelle 140**. Erreicht die Summe 140, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`); über 200 geht der Slice zurück zur Zerlegung (§Rückführungen). **Schätzung des Planners (2026-10-21): 60 bis 120** — §3.14 rund 30 bis 55 (ersetzter Absatz mit Befehl, Basis, Messpunkt, Pfaden, Rang-Zeiger, Paar und Retirement-Check, Löschungen gezählt), die ADR rund 15 bis 35 (Status, Autor-Zeile, Entscheidungssatz, Geschichte-Zeile), der Index 4 bis 8 (zwei Zeilen, Löschungen gezählt), Reserve 10 bis 20. Faktor Tests zu Produktcode entfällt: kein Code, kein Test, keine Gegenprobe (§1). Die Mitte der Schätzung (90) liegt bei 45 % der Marke.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik oben liegt **über 200**, oder die Bestätigung des Wortlauts durch den Nutzer verlangt einen dritten Liefer-Punkt (etwa den Sensor aus §1 oder das Nachziehen des Bestands). Dann schneidet der Planner, Vorschlag: §3.14 zuerst, die Annahme der ADR mit Index als zweiter Slice.
- `in-progress` → `open` (blockiert): Der Nutzer bestätigt den Wortlaut nicht oder verwirft den Befehl aus [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (dann ist die ADR nicht anzunehmen und der Befund geht an den Architect), oder `slice-harness-schaetzung-planner-faktor` liegt schon in `in-progress/` und beide schreiben denselben Absatz (dann entscheidet der Planner die Reihenfolge neu).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben. Zwei beobachtbare Kriterien: (1) `grep` nach `slice-harness-groessen-metrik` findet die Bindung in §1 und in der DoD von `slice-harness-schaetzung-planner-faktor`, und dessen §8 trägt eine Nachzählung mit dieser Kennung; (2) der Befehl aus §3.14, auf die Spanne der Closure von `slice-erster-release-homebrew-nachweis` angewandt (Basis `2c56dabb`, Pfade `tools harness spec docs README.md`), liefert über alle sechs Stände 223 bis 244 statt 907 bis 1148 — der Implementer rechnet es nach und nennt abweichende Ergebnisse dem Nutzer. Ein Mutationsprotokoll (`AGENTS.md` §3.10) entfällt: Der Slice liefert keinen Vertrag, kein Werkzeug und keine Zusage, die ein Test brechen könnte — die Regel ist Prosa ohne Sensor (§1, §6).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Der Slice liefert **keinen neuen Vertrag**: kein Werkzeug, kein Gate, keine Zusage mit Randformen, also kein Fall nach `AGENTS.md` §3.12 und keine Mutation nach §3.10. Die Randformen des Befehls selbst sind in [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) als *akzeptierte Negative* entschieden und stehen unten als Risiken mit Ausgang, nicht als offene Fragen.

- **Der Nutzer bestätigt den Wortlaut nicht oder ändert ihn** (§1, §3). — **Ausgang:** offen bis zur Closure; erwartet *entfallen* mit Datum und Fundstelle der Bestätigung in §7, sonst Rückführung nach `open` (§4).
- **Der Diff von `AGENTS.md` ersetzt Text in der Regel selbst**, und `slice-harness-schaetzung-planner-faktor` verlangt in seiner DoD einen Diff ohne Löschzeile in §3.14. — **Ausgang:** *entfallen*, wenn die Bindung in dessen §1 und DoD eingetragen ist (Liefer-Punkt 1) und seine DoD-Zeile die Berichtigung ausnimmt; sonst *eingetreten* → Befund an den Planner, der die Reihenfolge der beiden Slices neu entscheidet.
- **Der Nehmer `slice-harness-schaetzung-planner-faktor` liegt bei der Lieferung nicht mehr in `next/`** (zum Beispiel schon in `in-progress/` oder in `done/`): Dann ist er keine Adresse (`AGENTS.md` §3.13), und die Bindung hat keinen Ort. — **Ausgang:** *entfallen*, wenn er bei der Lieferung in `open/` oder `next/` liegt und die Bindung aufnimmt; sonst *eingetreten* → Befund an den Planner (die Berichtigung ist dann in dessen Closure-Notiz zu lesen, nicht in seinem Plan).
- **Ein Slice, dessen Liefer-Punkte im Planungsbaum liegen, misst 0** (akzeptiertes Negativ der ADR, Punkt *Konsequenzen*): Der Ausschluss ist eine Pfadangabe und nimmt auch eine ausdrücklich genannte Datei darunter weg. Dieser Slice ist selbst ein Grenzfall: Zwei seiner fünf Dateien liegen unter `docs/plan/planning` und zählen nicht. — **Ausgang:** *entfallen* mit der gemessenen Summe über `AGENTS.md docs/plan/adr`, die den Umfang trägt (Schätzung 60 bis 120); der Wortlaut in §3.14 nennt den Fall nicht als Ausnahme, weil §4 eines solchen Plans ihn selbst sagt.
- **Eine Datei, die aus einem ausgeschlossenen Pfad in einen genannten wandert, zählt als vollständige Einfügung** (akzeptiertes Negativ der ADR). Kein lebender Plan liefert das heute. — **Ausgang:** *entfallen*, wenn der Diff dieses Slice keine solche Bewegung enthält (zu belegen in §7), sonst *weiter offen* → `BEO-REPO/groessenprobe-ohne-messbaren-stand`.
- **Die Berichtigung wirkt nicht, und der Zähler steigt weiter** (`BEO-REPO/groessenprobe-ohne-messbaren-stand`, 5×, Stand verkörpert): Die Regel ist Prosa ohne Sensor; ein Plan, der den Befehl falsch abtippt oder einen Ausschluss weglässt, bleibt unbemerkt, bis das Review ihn findet. — **Ausgang:** *weiter offen* → `BEO-REPO/groessenprobe-ohne-messbaren-stand` (ein Auftreten **trotz** dieser Regel ist der Re-Evaluierungs-Trigger der ADR und macht den Sensor aus §1 fällig).
- **Die Zeile einer neuen ADR fehlt in der Haupttabelle des Index** — Befund des Nutzers vom 2026-10-21: Die Zeilen für fünf ADRs (Nummern 0051 bis 0055) fehlten in der Haupttabelle von `docs/plan/adr/README.md`, zwei Zeilen der Beziehungen waren als Tabelle zerbrochen; `make docs-check` sieht das nicht, weil es nur prüft, dass eine ADR-Kennung als Link auflöst, nicht dass jede ADR-Datei eine Zeile hat. Behoben vom Nutzer im Commit `4d9fd4ae`. Der Register-Eintrag `BEO-REPO/adr-ohne-zeile-im-index` liegt seit dem 2026-10-21 (vom Planner angelegt); der Befund steht dort unter *Benannt, nicht gezählt*, weil ihn kein abgeschlossener Vorgang trägt, und der Zähler ist 0. Dieser Slice trägt die Zeile für [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) in beide Tabellen (Liefer-Punkt 2); **ob und wo eine Prüfung entsteht, ist offen und Entscheidung des Nutzers** (§1, Meldung des Planners vom 2026-10-21). — **Ausgang:** *weiter offen* → `BEO-REPO/adr-ohne-zeile-im-index` (bei der Closure: Gegenbeleg, wenn beide Zeilen stehen, sonst ein Beleg, DoD).
- **Eine Binärdatei trägt `-` statt Zahlen und zählt 0**, und der Diff-Algorithmus ist nicht festgelegt (akzeptierte Negative der ADR; gemessen zwei Zeilen Unterschied auf 10158). — **Ausgang:** *entfallen* mit den Gründen der ADR: Zeilen sind bei einer Binärdatei keine Größe, und der Unterschied der Algorithmen liegt unter jeder Auflösung, mit der eine Marke gesetzt wird.
- **Die Messung setzt voraus, dass die Basis ein Vorfahr des gemessenen Stands ist und der Zweig keine fremde Arbeit aufgenommen hat** (akzeptiertes Negativ der ADR). — **Ausgang:** *entfallen*, wenn der Zweig dieses Slice aus dem Übergangs-Commit entsteht und keinen Merge aus `main` trägt (zu belegen in §7).
- **Die Meldeschwelle bleibt ohne Artefakt** (`BEO-REPO/meldeschwelle-ohne-artefakt`): Erreicht die Summe 140, gehört die Entscheidung des Planners als datierte Zeile in §7; diese Pflicht liefert erst `slice-harness-schaetzung-planner-faktor`. — **Ausgang:** *weiter offen* → `BEO-REPO/meldeschwelle-ohne-artefakt`, oder *entfallen*, wenn die Schwelle nicht erreicht wird.

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie ist so grob wie das Repo, trägt aber eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-21, gemergter Stand). Treffer: `BEO-REPO/groessenprobe-ohne-messbaren-stand` (5×, Stand verkörpert — die Regel, die dieser Slice berichtigt, ist ihr Zielort; mit diesem Slice erreicht sie kein neues Auftreten, sie wird geschärft), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (9×; von diesem Slice nicht berührt, er liefert keinen Test und keine Gegenprobe — der Faktor liegt bei `slice-harness-schaetzung-planner-faktor`), `BEO-REPO/meldeschwelle-ohne-artefakt` (2×, §6), `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (4×; nicht berührt, die Schätzung dieses Slice liegt bei 45 % der Marke), `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` (Regelzeile seit `slice-v1-abschluss-einspielen-tls-doku` — deshalb tragen Welle-Zuordnung und Stelle in der Reihenfolge hier ausdrücklich „Empfehlung des Planners“, Kopf und §4). Dazu `BEO-REPO/adr-ohne-zeile-im-index` (vom Planner am 2026-10-21 angelegt; Zähler 0, der Befund steht unter *Benannt, nicht gezählt*, weil kein abgeschlossener Vorgang ihn trägt — §6). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel.

alle berührten Sub-Areas GF (`REPO`: Entscheidung und Regel gehen dem Werkzeug voraus; der Slice ändert eine Regel und nimmt eine Entscheidung an, kein Code folgt).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-21, im Commit, der diesen Plan anlegt):** **Liefer-Punkte 2** — (1) der berichtigte Befehl in `AGENTS.md` §3.14 samt der Sendung in `slice-harness-schaetzung-planner-faktor`, (2) die Annahme von [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) mit Status und Index. Nicht gezählt: Gate-Läufe, Review-Report, Closure-Notiz, Register-Eintrag und -Beleg, Risiko-Ausgänge und die drei Paarungen (Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice). Die Sendung in den Nehmer ist Teil von Liefer-Punkt 1 und kein eigener: Sie trägt keinen neuen Gegenstand, sondern die Kennung dieses Slice in einen fremden Plan, und die Planung ist keine Schicht (`AGENTS.md` §3.13). **Schichten 1** — Regeln und Entscheidungen des Harness (`AGENTS.md`, `docs/plan/adr/`); die Planung (Pläne, Register) zählt nicht, Produkt-Code, Werkzeuge, Gegenproben und Verträge sind ausgeschlossen (§1). **Schätzung 60 bis 120** in der Metrik von §4, Mitte 90 bei **Marke 200** und Meldeschwelle 140 (45 % der Marke). Damit unter drei Liefer-Punkten und unter zwei Schichten.
