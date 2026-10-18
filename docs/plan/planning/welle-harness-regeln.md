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

Die Regeln und Werkzeuge des Harness, die aus dem Bau von v1 offen sind, sind geschärft: Die Regel-Slices tragen ihre Register-Einträge in `AGENTS.md`, `.claude/agents/` und `.claude/commands/` (`slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-schaetzung-planner-faktor`), der Commit-Träger lehnt Struktur-Kennungen ab (`slice-harness-commit-struktur-id`), `slice-mv` findet die Quelle am exakten Namen (`slice-harness-slice-mv-exakter-treffer`), Lastenheft- und MR-Kennungen sind Links, und das Doku-Gate verlangt sie (`slice-harness-lh-links-bestand`, `slice-harness-lh-links-pflicht`), und der Gate-Index trägt den Teil des Werkzeugs, sobald ai-harness-init ihn schreibt (`slice-harness-gate-index-werkzeug-teil`). Der Nutzer hat am 2026-10-18 entschieden, diese bis dahin wellenlosen Slices zu einer Welle zu gliedern. Es gibt kein Abnahmeszenario und keinen Meilenstein, der an ihr hängt.

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

- Alle Slices der Welle liegen in `done/`. Für `slice-harness-gate-index-werkzeug-teil` setzt das voraus, dass sein Auslöser eingetreten ist (ein Lauf von ai-harness-init schreibt einen Teil des Gate-Index); tritt er bis zur Closure nicht ein, entscheidet der Nutzer, ob der Slice §4 verlässt (Drift-Log der Roadmap) oder die Welle wartet (**offene Frage**, keine Entscheidung des Planners).
- Das *Mehr* gegenüber den Slice-DoDs (**Empfehlung des Planners, vom Nutzer nicht entschieden**): `make docs-check` ist auf dem Endstand der Welle grün, mit der Linkpflicht aus `slice-harness-lh-links-pflicht`, die für alle Dokumente gilt, die die übrigen Slices der Welle geändert haben (Commands, Pläne, `AGENTS.md`). Die DoD des Slice belegt die Pflicht am Bestand, den er kennt, nicht an dem, was die anderen sieben danach geschrieben haben.
- Für die drei Regel-Slices und die Werkzeug-Slices (Commit-Träger, `slice-mv`, Gate-Index) gibt es **kein gemeinsames *Mehr***: Ihre Bündelung ist Thema (Regeln und Werkzeuge des Harness), keine Closure-Bedingung. Verwirft der Nutzer das Kriterium oben, gilt das für die ganze Welle.
- `make gates` grün auf dem Endstand.
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
| slice-harness-commit-struktur-id | Commit-Träger lehnt Struktur-Kennungen ab | kein Anspruch des Produkts |
| slice-harness-lh-links-bestand | Lastenheft- und MR-Kennungen in ADRs und Wartungs-Doku als Links | kein Anspruch des Produkts ([MR-000](../../../harness/conventions.md#mr-000--baseline-aussage)) |
| slice-harness-lh-links-pflicht | Linkpflicht für Lastenheft- und MR-Kennungen im Doku-Gate | kein Anspruch des Produkts ([MR-000](../../../harness/conventions.md#mr-000--baseline-aussage)) |
| slice-harness-gate-index-werkzeug-teil | Gate-Index um den Teil des Werkzeugs teilen | kein Anspruch des Produkts ([MR-000](../../../harness/conventions.md#mr-000--baseline-aussage)) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: —
- Wird blockiert von: Welle [welle-v1-nachtrag](welle-v1-nachtrag.md) (Start-Trigger). Auf Slice-Ebene außerdem von `slice-harness-mutation` aus [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md): `slice-harness-lh-links-pflicht` steht in der Harness-Reihe hinter ihm (Drift-Log der Roadmap, 2026-10-09). Die Welle verschiebt diese Bindung nicht; sie liegt im Plan des Slice.
- Innerhalb der Welle: `slice-harness-lh-links-pflicht` setzt `slice-harness-lh-links-bestand` voraus (unmittelbar davor); `slice-harness-gate-index-werkzeug-teil` hängt an keinem Slice, nur an seinem Auslöser.

**Reihenfolge** (nach den Start-Triggern der Slices; WIP-Limit 1; die Welle ordnet nicht neu):

1. Die Reihe der fünf wellenlosen Slices hinter dem Release, soweit in dieser Welle: `slice-harness-slice-mv-exakter-treffer`, `slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-schaetzung-planner-faktor` (der fünfte, `slice-harness-gegenprobe-matrix-p7-flanke`, steht in [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md)).
2. `slice-harness-commit-struktur-id` — erster der übrigen Harness-Slices.
3. `slice-harness-lh-links-bestand`, danach `slice-harness-lh-links-pflicht` — am Ende der Harness-Reihe, hinter `slice-harness-mutation`.
4. `slice-harness-gate-index-werkzeug-teil` — außerhalb der Reihe, startet mit seinem Auslöser, frühestens nach den Wellen.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Das Produkt — die Welle ändert Regeln, Werkzeuge und Dokumente des Harness, keinen Produktcode.
- Tests und Gegenproben der Gates sowie das Mutations-Gate — [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md).
- Der Wortlaut der Regeln steht beim Nutzer: Die Regel-Slices legen ihn dem Nutzer vor und committen `AGENTS.md` und die Commands erst nach der Bestätigung (Plan der Slices); die Welle ändert das nicht.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-harness-regeln-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
