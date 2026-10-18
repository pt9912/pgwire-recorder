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

- [welle-erster-release](welle-erster-release.md) und [welle-v1-nachtrag](welle-v1-nachtrag.md) sind `done` (Wellen vor Harness, Entscheidung des Nutzers vom 2026-10-08 und vom 2026-10-16: die Harness-Reihe steht hinter beiden).

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
| slice-harness-schaetzung-planner-faktor | Der Planner schätzt in der Metrik von §3.14, mit benanntem Faktor, und nennt die Meldeschwelle | kein Anspruch des Produkts; Register-Einträge in deren Bezug |
| slice-harness-gegenprobe-matrix-p7-flanke | Eine Mutation der Gruppe `p7-sig` färbt nur Fälle ihrer Gruppe, oder die Voraussetzung dafür steht im Vertrag | kein Anspruch des Produkts |
| slice-harness-commit-struktur-id | Commit-Träger lehnt Struktur-Kennungen ab | kein Anspruch des Produkts |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md) (deren Start-Trigger verlangt diese Welle in `done/`). Die Bindung auf Slice-Ebene `slice-tests-ueberlebende-mutanten` setzt `slice-harness-commit-struktur-id` voraus ergibt sich damit aus der Reihenfolge der Wellen; sie bleibt im Plan des Slice stehen und ändert sich nicht.
- Wird blockiert von: Welle [welle-v1-nachtrag](welle-v1-nachtrag.md) (Start-Trigger).
- Innerhalb der Welle: keine Abhängigkeit zwischen den Slices außer der Reihenfolge unten.

**Reihenfolge** (nach den Start-Triggern der Slices; WIP-Limit 1; die Welle ordnet nicht neu; Reihenfolge der drei Wellen strikt: diese, [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md), [welle-lh-links](welle-lh-links.md), Entscheidung des Nutzers vom 2026-10-18):

1. Die Reihe der fünf wellenlosen Slices hinter dem Release: `slice-harness-slice-mv-exakter-treffer`, `slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-schaetzung-planner-faktor`, `slice-harness-gegenprobe-matrix-p7-flanke`.
2. `slice-harness-commit-struktur-id` — erster der übrigen Harness-Slices.

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
