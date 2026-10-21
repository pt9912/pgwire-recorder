# Welle welle-harness-regeln: Regeln und Werkzeuge des Harness nach v1

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug und kein Abnahmeszenario; die Welle ist Harness-Arbeit.

**Verantwortlich:** pt9912. **Datum:** 2026-10-18.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die Regeln und Werkzeuge des Harness, die aus dem Bau von v1 offen sind, sind geschärft: Die Regel-Slices tragen ihre Register-Einträge in `AGENTS.md`, `.claude/agents/` und `.claude/commands/` (`slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-schaetzung-planner-faktor`), `slice-mv` findet die Quelle am exakten Namen (`slice-harness-slice-mv-exakter-treffer`), die Gegenprobe der Versionsmatrix färbt mit einer Mutation der Gruppe `p7-sig` nur Fälle ihrer Gruppe (`slice-harness-gegenprobe-matrix-p7-flanke`), und der Commit-Träger lehnt Struktur-Kennungen ab (`slice-harness-commit-struktur-id`). Der Nutzer hat am 2026-10-18 entschieden, diese bis dahin wellenlosen Slices zu einer Welle zu gliedern, und am selben Tag, die Wellen 2 und 3 entlang der Harness-Reihe (Entscheidung des Nutzers vom 2026-10-08) neu zu schneiden, damit sie strikt nacheinander laufen: Diese Welle ist die erste der drei (danach [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md), danach [welle-lh-links](welle-lh-links.md)). Es gibt kein Abnahmeszenario und keinen Meilenstein, der an ihr hängt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [welle-erster-release](done/welle-erster-release/welle-erster-release.md) und [welle-v1-nachtrag](welle-v1-nachtrag.md) sind `done` (Wellen vor Harness, Entscheidung des Nutzers vom 2026-10-08 und vom 2026-10-16: die Harness-Reihe steht hinter beiden).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Das *Mehr* gegenüber den Slice-DoDs (Entscheidung des Nutzers vom 2026-10-18): `make docs-check` und `make gates` sind auf dem Endstand der Welle grün, mit allen Regeln eingetragen. Die DoDs der Slices belegen ihren Lauf je für sich auf dem Stand, den sie kennen, nicht den Stand, den die übrigen fünf danach schreiben.
- Für die Regel-Slices, `slice-harness-slice-mv-exakter-treffer`, `slice-harness-gegenprobe-matrix-p7-flanke` und `slice-harness-commit-struktur-id` gibt es **kein weiteres gemeinsames *Mehr***: Ihre Bündelung ist Thema (Regeln und Werkzeuge des Harness), keine Closure-Bedingung.
- Closure-Notiz in `welle-harness-regeln-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-harness-slice-mv-exakter-treffer | `slice-mv` findet die Quelle am exakten Namen, und eine Gegenprobe hält es | kein Anspruch des Produkts |
| slice-harness-regeln-nach-v1-abschluss | Regeln aus dem Verbindungsende-Block (Abdeckungstabellen zählen nicht als Schicht, Wortlaut-Suche und Arbeitsbaum vor E2E-Mutanten, offene Fragen zu §3.10) | kein Anspruch des Produkts; Register-Einträge in deren Bezug |
| slice-harness-regeln-spezifikation-satzabgleich | Der Architect gleicht vor dem Code die Sätze einer Spezifikationsstelle ab | kein Anspruch des Produkts; `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` |
| slice-harness-groessen-metrik | `AGENTS.md` §3.14 trägt den berichtigten Befehl der Größen-Metrik, und [ADR-0056](../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) wird angenommen | kein Anspruch des Produkts; `BEO-REPO/groessenprobe-ohne-messbaren-stand` |
| slice-harness-schaetzung-planner-faktor | Der Planner schätzt in der Metrik von §3.14, mit benanntem Faktor, und nennt die Meldeschwelle | kein Anspruch des Produkts; Register-Einträge in deren Bezug |
| slice-harness-gegenprobe-matrix-p7-flanke | Eine Mutation der Gruppe `p7-sig` färbt nur Fälle ihrer Gruppe, oder die Voraussetzung dafür steht im Vertrag | kein Anspruch des Produkts |
| slice-harness-commit-struktur-id | Commit-Träger lehnt Struktur-Kennungen ab | kein Anspruch des Produkts |
| slice-harness-rollenlauf-zeitgrenze | Ein Rollenlauf hält bei einer Stunde an und berichtet, erzwungen von einem deterministischen Wächter im Werkzeug-Aufruf | kein Anspruch des Produkts; neuer Register-Eintrag `BEO-REPO/rollenlauf-ohne-zeitgrenze` in seinem Bezug |

`slice-harness-rollenlauf-zeitgrenze` ist am 2026-10-21 in `open/` angelegt (Anlass: der Implementer-Lauf von `slice-harness-devcontainer-einrichten` lief 4,34 h, ohne dass eine Marke feuerte). **Entscheidung des Nutzers vom 2026-10-21** ist die **Sache** — die Grenze von einer Stunde je Rollenlauf („Wir haben rollen (siehe .claude/agents) und commands -- max. 1 Stunde"), Anhalten mit Bericht statt Abbruch und der Mechanismus (Wächter im Werkzeug-Aufruf, reines bash+awk, kein Hintergrundprozess, Span-Strom als Quelle). **Empfehlung des Planners**, nicht entschieden, sind Welle und Stelle in der Reihenfolge (§5), der Schnitt (die Wiederholungs-Grenze für volle Gegenproben bleibt draußen, dessen §1) und die Fail-Richtung (dessen §6). Zwei Punkte liegen beim Nutzer, bevor der Slice startbar ist: ob `AGENTS.md` §3.9 die Rollen-Typen unter `.claude/agents/` in seine Liste der Harness-Artefakte aufnimmt, und ob die Stunde als **zweite Rückführungs-Marke** neben der Zeilenzahl in diesen Slice oder in `slice-harness-groessen-metrik` gehört (§6 des Plans, Punkte (a) und (b)).

`slice-harness-groessen-metrik` trägt `Welle: welle-harness-regeln` nach **Entscheidung des Nutzers vom 2026-10-21** (wörtlich, vom Koordinator übermittelt: „ok deine empfehlugnen“ auf eine Vorlage mit sechs Punkten, deren erster Welle und Stelle dieses Slice war); seine Stelle in der Reihenfolge — vor `slice-harness-schaetzung-planner-faktor` — ist mit derselben Entscheidung gesetzt und steht in §5. Bis dahin waren Welle und Stelle Empfehlung des Planners (`BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet`).

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md) (deren Start-Trigger verlangt diese Welle in `done/`). Die Bindung auf Slice-Ebene `slice-tests-ueberlebende-mutanten` setzt `slice-harness-commit-struktur-id` voraus ergibt sich damit aus der Reihenfolge der Wellen; sie bleibt im Plan des Slice stehen und ändert sich nicht.
- Wird blockiert von: Welle [welle-v1-nachtrag](welle-v1-nachtrag.md) (Start-Trigger).
- Innerhalb der Welle: keine Abhängigkeit zwischen den Slices außer der Reihenfolge unten.

**Reihenfolge** (nach den Start-Triggern der Slices; WIP-Limit 1; die Welle ordnet nicht neu; Reihenfolge der drei Wellen strikt: diese, [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md), [welle-lh-links](welle-lh-links.md), Entscheidung des Nutzers vom 2026-10-18):

1. Die Reihe der fünf wellenlosen Slices hinter dem Release: `slice-harness-slice-mv-exakter-treffer`, `slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-schaetzung-planner-faktor`, `slice-harness-gegenprobe-matrix-p7-flanke`.
2. `slice-harness-groessen-metrik` steht **vor** `slice-harness-schaetzung-planner-faktor` (Entscheidung des Nutzers vom 2026-10-21, §4): Beide schreiben denselben Absatz `AGENTS.md` §3.14, und die DoD von `slice-harness-schaetzung-planner-faktor` verlangt einen Diff ohne Löschzeile in der Regel selbst, nimmt die Berichtigung also nicht an, solange sie nicht eingetragen ist (`AGENTS.md` §3.13, [ADR-0056](../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) *Folgepflicht*). Die Reihe aus Punkt 1 bleibt im Übrigen unverändert.
3. `slice-harness-commit-struktur-id` — erster der übrigen Harness-Slices.
4. `slice-harness-rollenlauf-zeitgrenze` steht **hinter** `slice-harness-groessen-metrik` (**Empfehlung des Planners vom 2026-10-21**, nicht entschieden): Der Metrik-Slice schreibt `AGENTS.md` §3.14 neu, und die neue harte Regel entsteht unmittelbar dahinter — ein gleichzeitiger Lauf erzeugt einen Konflikt statt einer Regel. Seine Stelle gegenüber den übrigen Slices der Reihe ist nicht festgelegt; sein Start-Trigger verlangt zusätzlich, dass `slice-harness-devcontainer-einrichten` abgewickelt ist (dessen Lauf hält `tools/` und `harness/`) und dass Vertrag und ADR des Architect vorliegen (dessen §4).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Das Produkt — die Welle ändert Regeln, Werkzeuge und Dokumente des Harness, keinen Produktcode.
- Tests und Gegenproben des Bestands (außer der Gegenprobe der Versionsmatrix) sowie das Mutations-Gate — [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md).
- Lastenheft- und MR-Kennungen als Links und deren Pflicht im Doku-Gate — [welle-lh-links](welle-lh-links.md).
- Der Gate-Index um den Teil des Werkzeugs (`slice-harness-gate-index-werkzeug-teil`) — der Slice ist wellenlos: Sein Auslöser ist ein Lauf von ai-harness-init, und eine Welle, die auf Fremdes wartet, schließt nie (Entscheidung des Nutzers vom 2026-10-18).
- Der Wortlaut der Regeln steht beim Nutzer: Die Regel-Slices legen ihn dem Nutzer vor und committen `AGENTS.md` und die Commands erst nach der Bestätigung (Plan der Slices); die Welle ändert das nicht.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-harness-regeln-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
