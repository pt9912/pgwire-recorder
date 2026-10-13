# Welle welle-erster-release: Erster Release

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M4.

**Verantwortlich:** pt9912. **Datum:** 2026-10-03.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Der erste Release ist veröffentlicht: die Binaries für die Zielplattformen und das Docker/OCI-Image in `ghcr.io` und `docker.io`. Die Homebrew-Formel im Tap und Abnahmeszenario 11 ([`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)) gehören seit der Entscheidung des Nutzers vom 2026-10-11 zu [welle-v1-nachtrag](welle-v1-nachtrag.md), die auf diese Welle folgt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-v1-abschluss](welle-v1-abschluss.md) done, und `slice-harness-coverage` liegt in `done/` (Abnahmeszenario 17, M3 vor M4; Entscheidung des Nutzers vom 2026-10-08).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Der erste stabile Tag ist nach der Freigabe-Checkliste veröffentlicht (Binaries mit `SHA256SUMS` am Release, Image in beiden Registries mit gleichem Index-Digest, Smoke aus Abnahmeszenario 9 bestanden) — das *Mehr* gegenüber den Slice-DoDs, deren Pipeline nur ein Probe-Tag belegt. Abnahmeszenario 11 ist keine Bedingung dieser Welle; [welle-v1-nachtrag](welle-v1-nachtrag.md) weist es nach.
- `make gates` grün.
- Closure-Notiz in `welle-erster-release-results.md`; sie trägt den ausgefüllten Eintrag der Freigabe-Checkliste des ersten Releases (Entscheidung des Nutzers vom 2026-10-09, geliefert von `slice-erster-release-freigabe`).

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-erster-release-veroeffentlichung | Veröffentlichung von Binaries und Images | [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität) |
| slice-erster-release-freigabe | Freigabe und erster echter Tag | [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-v1-nachtrag](welle-v1-nachtrag.md).
- Wird blockiert von: Welle [welle-v1-abschluss](welle-v1-abschluss.md) und den wellenlosen Slices `slice-harness-abdeckung-gate` und `slice-harness-coverage`.
- Innerhalb der Welle: `slice-erster-release-freigabe` setzt `slice-erster-release-veroeffentlichung` voraus (der Scan läuft in dessen Pipeline vor dem Push).

**Reihenfolge** (Entscheidungen des Nutzers vom 2026-10-08: Wellen vor Harness, M3 vor M4; Schnitt von `slice-erster-release-veroeffentlichung` vom 2026-10-09; Umplanung von `slice-erster-release-homebrew-nachweis` nach `welle-v1-nachtrag` vom 2026-10-11; WIP-Limit 1):

1. `slice-erster-release-veroeffentlichung` — Pipeline für Binaries und Image, belegt mit einem Probe-Tag.
2. `slice-erster-release-freigabe` — Scan, Freshness-Audit, Freigabe-Checkliste und der erste echte Tag (setzt 1 voraus).

Vor dieser Welle liegen die Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; wellenlos, Entscheidung des Nutzers vom 2026-10-13), `slice-harness-integration-volumes`, `slice-harness-abdeckung-gate` und `slice-harness-coverage`: Mit Coverage wird Abnahmeszenario 17 nachweisbar, M3 ist damit vor M4 erreicht. Danach folgt [welle-v1-nachtrag](welle-v1-nachtrag.md), dann die übrigen Harness-Slices, beginnend mit `slice-harness-commit-struktur-id`; ihre Reihenfolge steht in deren §4.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Die Homebrew-Formel im Tap und ihr Nachweis (Abnahmeszenario 11) — `slice-erster-release-homebrew-nachweis` in [welle-v1-nachtrag](welle-v1-nachtrag.md) (Entscheidung des Nutzers vom 2026-10-11); der Release ist ohne sie vollständig, weil kein Release-Artefakt von der Formel abhängt.
- Weitere Plattformen oder Registries — nicht zugesichert.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-erster-release-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
