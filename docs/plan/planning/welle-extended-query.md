# Welle welle-extended-query: Extended Query Protocol

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M2.

**Verantwortlich:** pt9912. **Datum:** 2026-10-03.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Das Extended Query Protocol ist im Record- und Replay-Modus umgesetzt; ein Client mit Prepared Statements wird aufgezeichnet und ohne PostgreSQL wiedergegeben ([`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), Abnahmeszenario 7).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-walking-skeleton](done/welle-walking-skeleton/welle-walking-skeleton.md) done.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Abnahmeszenario 7 und der Fehlerfall von Abnahmeszenario 6 sind end-to-end mit einem verbreiteten Go-Client im Standardmodus automatisiert nachgewiesen, und der Walking-Skeleton-Smoke läuft weiter — das *Mehr* gegenüber den Slice-DoDs.
- `make gates` grün.
- Closure-Notiz in `welle-extended-query-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-extended-query-modell | Spezifikation und Domain-Modell | [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) |
| slice-extended-query-record | Extended Query im Record | [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-06`](../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten) |
| slice-extended-query-replay | Extended Query im Replay | [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-10`](../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage) |
| slice-extended-query-lebendpruefung | Lebendprüfungen im Replay | [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-QA-02`](../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-replay-semantik](welle-replay-semantik.md).
- Wird blockiert von: Welle [welle-walking-skeleton](done/welle-walking-skeleton/welle-walking-skeleton.md).
- Innerhalb der Welle: `slice-extended-query-record` und `slice-extended-query-replay` setzen `slice-extended-query-modell` voraus. `slice-extended-query-lebendpruefung` setzt `slice-extended-query-replay` voraus.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- `COPY`, Replikationsprotokoll, TLS — nicht Teil dieser Welle.
- SQL-Normalisierung oder Parametermatching über den exakten Vergleich hinaus — Folge-Anforderung, falls gewünscht.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-extended-query-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
