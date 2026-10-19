# Welle welle-ci-nachtlaeufe: CI-Nachtläufe auf GitHub

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-10-18.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Zwei Nachtläufe laufen auf GitHub außerhalb von `make gates`: die Gegenproben (`gegenproben.yml`, über die Sammelziele `make gegenproben` und `make gegenproben-host`) und der Lauf gegen upstream (`upstream-drift.yml`: Audit der vendored Baseline, des gepinnten d-check und der Image-Pins nach Tag und Digest). Die Welle entstand am 2026-10-18 aus welle-erster-release (Entscheidung des Nutzers vom 2026-10-18, „Erst das Release voranbringen“): Der kritische Pfad zum ersten Release enthält die Nachtläufe nicht; das Freshness-Audit der Freigabe-Checkliste läuft für den ersten Tag von Hand (`slice-erster-release-freigabe`).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-erster-release](welle-erster-release.md) done (Entscheidung des Nutzers vom 2026-10-18).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Je Workflow, `gegenproben.yml` und `upstream-drift.yml`, liegt ein grüner Nachtlauf (Auslöser `schedule`) auf `main` vor, mit Lauf-Kennung in der Closure-Notiz — das *Mehr* gegenüber den Slice-DoDs, deren Beleg ein Lauf per `workflow_dispatch` sein darf.
- `make gates` grün.
- Closure-Notiz in `welle-ci-nachtlaeufe-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-harness-ci-gegenproben | Der nächtliche Lauf der Gegenproben auf GitHub (`gegenproben.yml`) | kein Anspruch des Produkts |
| slice-harness-ci-freshness | Der Nachtlauf gegen upstream: `upstream-drift.yml`, Audit der Baseline und von d-check | kein Anspruch des Produkts |
| slice-harness-ci-freshness-images | Aktualität der Image-Pins im Nachtlauf (Tag-Achse und Digest-Achse) | kein Anspruch des Produkts |

Die drei Slices tragen im Kopf die Welle-Zeile `welle-ci-nachtlaeufe (Entscheidung des Nutzers vom 2026-10-18)`; sie verließen damit §4 von welle-erster-release.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle.
- Wird blockiert von: Welle [welle-erster-release](welle-erster-release.md) (Start) und den Slices `slice-harness-ci` (liefert `ci.yml` und die Konventionen der Workflows) und `slice-harness-gegenproben-aus-gates-sammelziele` (liefert `make gegenproben`), beide schon Voraussetzung von welle-erster-release bzw. `done/`.
- Innerhalb der Welle: `slice-harness-ci-freshness` setzt `slice-harness-ci` voraus, `slice-harness-ci-freshness-images` setzt `slice-harness-ci-freshness` voraus (Skript, Workflow, Vertrag), `slice-harness-ci-gegenproben` setzt `slice-harness-ci` und `slice-harness-gegenproben-aus-gates-sammelziele` voraus.

**Reihenfolge** (wie in welle-erster-release bis zur Entscheidung des Nutzers vom 2026-10-18; Schnitt und Reihenfolge der drei Slices: Empfehlung des Planners, `slice-harness-ci-workflows` §6 Frage 1; WIP-Limit 1):

1. `slice-harness-ci-gegenproben` — `gegenproben.yml` (setzt `slice-harness-ci` und die Sammelziele voraus).
2. `slice-harness-ci-freshness` — `upstream-drift.yml`, Audit der Baseline und von d-check (setzt `slice-harness-ci` voraus).
3. `slice-harness-ci-freshness-images` — Image-Pins im Nachtlauf (setzt 2 voraus).

Die Stellung der Welle gegenüber [welle-v1-nachtrag](welle-v1-nachtrag.md) und den Harness-Wellen ist nicht entschieden: Beide folgen auf welle-erster-release; keine Kante zwischen ihnen und dieser Welle (Frage an den Nutzer).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- `ci.yml`, das Gate der Workflows und `release.yml` — `slice-harness-ci-workflows`, `slice-harness-ci`, `slice-erster-release-veroeffentlichung` und `slice-erster-release-veroeffentlichung-image` in [welle-erster-release](welle-erster-release.md); diese Welle ändert sie nicht.
- Das Freshness-Audit vor dem ersten Tag — `slice-erster-release-freigabe` führt es von Hand; ein Nachtlauf dieser Welle ist dafür weder Beleg noch Bedingung (Entscheidung des Nutzers vom 2026-10-18).
- Das Fahren der Gegenproben im Gegenproben-Container vor einem Release — Pflicht nach `AGENTS.md` §6, bleibt; ein Nachtlauf beweist den Stand des Standard-Zweigs, nicht den eines Releases.
- Branch-Schutz („required status check“), Geheimnisse und Sichtbarkeit des Repos auf GitHub — Einstellungen des Betreibers, aus dem Klon nicht auditierbar; die Welle nennt sie als Handlung, ohne sie zu behaupten.
- Aktualität und Gültigkeit der Action-SHA und der Werkzeug-Versionen (`u-boot`, `devcontainer`) in den Workflows — kein Slice der Welle nimmt sie an (`slice-harness-ci-workflows` §6 Frage 3, `slice-harness-ci-gegenproben` §6 Frage 3).
- Benachrichtigung bei rotem Nachtlauf über die Anzeige von GitHub hinaus — nicht zugesichert.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-ci-nachtlaeufe-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
