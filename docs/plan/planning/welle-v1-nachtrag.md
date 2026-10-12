# Welle welle-v1-nachtrag: v1-Nachtrag: Homebrew und CancelRequest ohne Schlüssel

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M5.

**Verantwortlich:** pt9912. **Datum:** 2026-10-11.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die zwei Anforderungsreste von v1, die hinter dem ersten Release liegen, sind umgesetzt: Die Homebrew-Formel liegt im Tap, und Abnahmeszenario 11 ([`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)) ist auf macOS und Linux nachgewiesen; ein `CancelRequest` ohne Abbruchschlüssel wird wie jeder `CancelRequest` behandelt ([`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)). Damit sind alle Anforderungen (MUSS und SOLL) umgesetzt. Die Slices standen in `welle-v1-abschluss` (Homebrew, `CancelRequest`) und `welle-erster-release` (Nachweis); der Nutzer hat sie am 2026-10-11 hinter den ersten Release gezogen, damit M4 ohne sie erreicht wird.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-erster-release](welle-erster-release.md) done (Entscheidung des Nutzers vom 2026-10-11).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Abnahmeszenario 11 (Homebrew) ist auf macOS und Linux mit dem echten Tap nachgewiesen — das *Mehr* gegenüber den Slice-DoDs.
- `make gates` grün.
- Closure-Notiz in `welle-v1-nachtrag-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-v1-abschluss-cancel-ohne-schluessel | CancelRequest ohne Schlüssel | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) |
| slice-v1-abschluss-homebrew | Homebrew-Bereitstellung | [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |
| slice-erster-release-homebrew-nachweis | Homebrew-Nachweis (Abnahmeszenario 11) | [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: —
- Wird blockiert von: Welle [welle-erster-release](welle-erster-release.md).
- Innerhalb der Welle: `slice-v1-abschluss-cancel-ohne-schluessel` setzt `slice-v1-abschluss-protokollrand` aus [welle-v1-abschluss](welle-v1-abschluss.md) voraus (der übrige Protokollrand); `slice-v1-abschluss-homebrew` setzt `slice-v1-abschluss-container` aus [welle-v1-abschluss](welle-v1-abschluss.md) voraus (Release-Artefakte); `slice-erster-release-homebrew-nachweis` setzt `slice-v1-abschluss-homebrew` (das Verfahren der Formel) und `slice-erster-release-freigabe` aus [welle-erster-release](welle-erster-release.md) voraus (erster stabiler Tag).

**Reihenfolge** (Entscheidung des Nutzers vom 2026-10-11; WIP-Limit 1):

1. `slice-v1-abschluss-cancel-ohne-schluessel` — `CancelRequest` ohne Schlüssel in `record` und `replay`.
2. `slice-v1-abschluss-homebrew` — Verfahren der Formel, belegt mit einem Probe-Tag und einem Probe-Tap.
3. `slice-erster-release-homebrew-nachweis` — Abnahmeszenario 11 mit dem echten Tap (setzt 2 voraus).

Vor dieser Welle liegen `slice-harness-integration-volumes`, `slice-harness-abdeckung-gate` und `slice-harness-coverage` (M3) und [welle-erster-release](welle-erster-release.md) (M4). Danach folgen die drei wellenlosen Slices der Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; Entscheidung des Nutzers vom 2026-10-12, Ziel: ein Release), dann die übrigen Harness-Slices, beginnend mit `slice-harness-commit-struktur-id`; ihre Reihenfolge steht in deren §4.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Die Aufnahme in das Standard-Repository von Homebrew — Out-of-Scope von LH-FA-19.
- Weitere Plattformen oder Registries — nicht zugesichert.
- Der erste Release und seine Pipeline — `welle-erster-release`; diese Welle setzt den veröffentlichten ersten stabilen Tag voraus und ändert ihn nicht.
- Übriger Protokollrand — `slice-v1-abschluss-protokollrand` in `welle-v1-abschluss`.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-v1-nachtrag-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
