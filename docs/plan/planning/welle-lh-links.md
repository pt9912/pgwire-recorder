# Welle welle-lh-links: Lastenheft- und MR-Kennungen als Links, mit Pflicht im Doku-Gate

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

Lastenheft- und MR-Kennungen stehen in den Dokumenten des Repos als Links, und das Doku-Gate verlangt sie: Der Bestand in ADRs und Wartungs-Doku ist verlinkt (`slice-harness-lh-links-bestand`), das Muster in `.d-check.yml` ist aktiv, und `make docs-check` meldet eine Kennung ohne Link (`slice-harness-lh-links-pflicht`). Der Nutzer hat am 2026-10-18 entschieden, die Wellen 2 und 3 entlang der Harness-Reihe (Entscheidung des Nutzers vom 2026-10-08) neu zu schneiden, damit sie strikt nacheinander laufen; diese Welle ist die dritte, hinter [welle-harness-regeln](welle-harness-regeln.md) und [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md). Die beiden Slices standen bis dahin in welle-harness-regeln. Es gibt kein Abnahmeszenario und keinen Meilenstein, der an ihr hängt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md) ist `done` (Neuschnitt der Wellen 2 und 3, Entscheidung des Nutzers vom 2026-10-18; die Wellen laufen strikt nacheinander und `slice-harness-lh-links-pflicht` steht in der Harness-Reihe hinter `slice-harness-mutation`, Drift-Log der Roadmap vom 2026-10-09). Damit sind auch [welle-harness-regeln](welle-harness-regeln.md), [welle-erster-release](done/welle-erster-release/welle-erster-release.md) und [welle-v1-nachtrag](welle-v1-nachtrag.md) `done`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Das *Mehr* gegenüber den Slice-DoDs (Entscheidung des Nutzers vom 2026-10-18): `make docs-check` verlangt für jede Lastenheft-Kennung einen Link und ist auf dem Endstand der Welle grün. Die DoD von `slice-harness-lh-links-pflicht` belegt die Pflicht an dem Bestand, den sie kennt; der Endstand schließt ein, was zwischen den beiden Slices neu entstand.
- `make gates` grün auf dem Endstand.
- Closure-Notiz in `welle-lh-links-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-harness-lh-links-bestand | Lastenheft- und MR-Kennungen in ADRs und Wartungs-Doku als Links | kein Anspruch des Produkts ([MR-000](../../../harness/conventions.md#mr-000--baseline-aussage)) |
| slice-harness-lh-links-pflicht | Linkpflicht für Lastenheft- und MR-Kennungen im Doku-Gate | kein Anspruch des Produkts ([MR-000](../../../harness/conventions.md#mr-000--baseline-aussage)) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: —
- Wird blockiert von: Welle [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md) (Start-Trigger; die Wellen laufen strikt nacheinander). Die Bindung auf Slice-Ebene (`slice-harness-lh-links-pflicht` steht hinter `slice-harness-mutation`) ergibt sich daraus und bleibt im Plan des Slice stehen.
- Innerhalb der Welle: `slice-harness-lh-links-pflicht` setzt `slice-harness-lh-links-bestand` voraus (unmittelbar davor).

**Reihenfolge** (nach den Start-Triggern der Slices; WIP-Limit 1; die Welle ordnet nicht neu):

1. `slice-harness-lh-links-bestand`
2. `slice-harness-lh-links-pflicht`

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Das Produkt — die Welle ändert Dokumente, das Muster des Doku-Gates und Commands des Harness, keinen Produktcode.
- Die Regeln, Werkzeuge und Gegenproben des Harness — [welle-harness-regeln](welle-harness-regeln.md) und [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md).
- Der Gate-Index um den Teil des Werkzeugs (`slice-harness-gate-index-werkzeug-teil`) — wellenlos, sein Auslöser ist ein Lauf von ai-harness-init.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-lh-links-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
