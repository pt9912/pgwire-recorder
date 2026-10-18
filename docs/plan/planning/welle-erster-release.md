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

Der erste Release ist veröffentlicht: die Binaries für die Zielplattformen und das Docker/OCI-Image in `ghcr.io` und `docker.io`. Vor der Veröffentlichung läuft die Integration des Repos in GitHub Actions (CI auf jeden Pull-Request und Push, nächtliche Läufe der Gegenproben und gegen upstream; Auftrag des Nutzers vom 2026-10-18), und die Workflow-Dateien, darunter die der Veröffentlichung, bestehen ein Gate. Die Homebrew-Formel im Tap und Abnahmeszenario 11 ([`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)) gehören seit der Entscheidung des Nutzers vom 2026-10-11 zu [welle-v1-nachtrag](welle-v1-nachtrag.md), die auf diese Welle folgt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-v1-abschluss](done/welle-v1-abschluss/welle-v1-abschluss.md) done, und `slice-harness-coverage` liegt in `done/` (Abnahmeszenario 17, M3 vor M4; Entscheidung des Nutzers vom 2026-10-08), dazu `slice-harness-devcontainer-u-boot-0-10` in `done/` (wellenlos; Entscheidung des Nutzers vom 2026-10-17: lang laufende Slices laufen im Dev-Container, damit `main` nicht durch Implementer-Läufe gesperrt ist; `slice-harness-devcontainer-einrichten` ist nicht Bedingung, Entscheidung des Nutzers vom 2026-10-17).

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
| slice-harness-ci-workflows | Gate der Workflows: Pins, Rechte und Zeitgrenzen der Dateien unter `.github/workflows/` | kein Anspruch des Produkts |
| slice-harness-ci | Die CI auf GitHub: `ci.yml` fährt `make gates` und prüft die Commit-Range | kein Anspruch des Produkts |
| slice-harness-ci-gegenproben | Der nächtliche Lauf der Gegenproben auf GitHub (`gegenproben.yml`) | kein Anspruch des Produkts |
| slice-harness-ci-freshness | Der Nachtlauf gegen upstream: `upstream-drift.yml`, Audit der Baseline und von d-check | kein Anspruch des Produkts |
| slice-harness-ci-freshness-images | Aktualität der Image-Pins im Nachtlauf (Tag-Achse und Digest-Achse) | kein Anspruch des Produkts |
| slice-erster-release-veroeffentlichung | Veröffentlichung von Binaries und Images | [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität) |
| slice-erster-release-freigabe | Freigabe und erster echter Tag | [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-v1-nachtrag](welle-v1-nachtrag.md).
- Wird blockiert von: Welle [welle-v1-abschluss](done/welle-v1-abschluss/welle-v1-abschluss.md) und den wellenlosen Slices `slice-harness-abdeckung-gate`, `slice-harness-coverage` und `slice-harness-devcontainer-u-boot-0-10` (nicht an `slice-harness-devcontainer-einrichten`). Für die ersten drei Slices der CI zusätzlich `slice-harness-gegenproben-aus-gates` (Start von `slice-harness-ci-workflows` und `slice-harness-ci`) und `slice-harness-gegenproben-aus-gates-sammelziele` (Start von `slice-harness-ci-gegenproben`); beide gehören zu [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md) und stehen dort vorn (Drift-Log der Roadmap, 2026-10-18); sie sind Bedingung der genannten Slices, nicht des Starts dieser Welle.
- Innerhalb der Welle: `slice-erster-release-freigabe` setzt `slice-erster-release-veroeffentlichung` voraus (der Scan läuft in dessen Pipeline vor dem Push). `slice-harness-ci` setzt `slice-harness-ci-workflows` voraus (das Gate prüft `ci.yml`), `slice-harness-ci-gegenproben` und `slice-harness-ci-freshness` setzen `slice-harness-ci` voraus, `slice-harness-ci-freshness-images` setzt `slice-harness-ci-freshness` voraus (Skript und Workflow), `slice-erster-release-veroeffentlichung` setzt `slice-harness-ci-workflows` und `slice-harness-ci` voraus (Pins und Härtung von `release.yml`), `slice-erster-release-freigabe` setzt zusätzlich `slice-harness-ci-freshness` und `slice-harness-ci-freshness-images` voraus (der Nachtlauf ist der Beleg des Freshness-Audits).

**Reihenfolge** (Auftrag des Nutzers vom 2026-10-18: die CI vor die Veröffentlichung; der Schnitt in fünf Slices und die Reihenfolge darin sind Empfehlung des Planners, nicht entschieden; Entscheidungen des Nutzers vom 2026-10-08: Wellen vor Harness, M3 vor M4; Schnitt von `slice-erster-release-veroeffentlichung` vom 2026-10-09; Umplanung von `slice-erster-release-homebrew-nachweis` nach `welle-v1-nachtrag` vom 2026-10-11; WIP-Limit 1):

1. `slice-harness-ci-workflows` — Gate der Workflows (Start: `slice-harness-gegenproben-aus-gates` in `done/`).
2. `slice-harness-ci` — `ci.yml` mit `make gates` und Commit-Range (setzt 1 voraus).
3. `slice-harness-ci-gegenproben` — `gegenproben.yml` (setzt 2 und `slice-harness-gegenproben-aus-gates-sammelziele` voraus).
4. `slice-harness-ci-freshness` — `upstream-drift.yml`, Audit der Baseline und von d-check (setzt 2 voraus).
5. `slice-harness-ci-freshness-images` — Image-Pins im Nachtlauf (setzt 4 voraus).
6. `slice-erster-release-veroeffentlichung` — Pipeline für Binaries und Image, belegt mit einem Probe-Tag (setzt 1 und 2 voraus).
7. `slice-erster-release-freigabe` — Scan, Freshness-Audit, Freigabe-Checkliste und der erste echte Tag (setzt 6, 4 und 5 voraus).

Vor dieser Welle liegen die Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; wellenlos, Entscheidung des Nutzers vom 2026-10-13), `slice-harness-integration-volumes`, `slice-harness-abdeckung-gate` und `slice-harness-coverage`: Mit Coverage wird Abnahmeszenario 17 nachweisbar, M3 ist damit vor M4 erreicht. Hinter Coverage stehen wellenlos `slice-harness-devcontainer-u-boot-0-10` (vor der Closure von welle-v1-abschluss und vor den Slices dieser Welle) und danach `slice-harness-devcontainer-einrichten` (hinter jenem, ohne Bindung an Closure und Release; Entscheidung des Nutzers vom 2026-10-17). Danach folgt [welle-v1-nachtrag](welle-v1-nachtrag.md), dann die übrigen Harness-Slices, beginnend mit `slice-harness-commit-struktur-id`; ihre Reihenfolge steht in deren §4.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Die Homebrew-Formel im Tap und ihr Nachweis (Abnahmeszenario 11) — `slice-erster-release-homebrew-nachweis` in [welle-v1-nachtrag](welle-v1-nachtrag.md) (Entscheidung des Nutzers vom 2026-10-11); der Release ist ohne sie vollständig, weil kein Release-Artefakt von der Formel abhängt.
- Weitere Plattformen oder Registries — nicht zugesichert.
- Branch-Schutz („required status check“), Geheimnisse und Sichtbarkeit des Repos auf GitHub — Einstellungen des Betreibers, aus dem Klon nicht auditierbar; die Welle nennt sie als Handlung, ohne sie zu behaupten.
- Aktualität und Gültigkeit der Action-SHA und der Werkzeug-Versionen in den Workflows — kein Slice der Welle nimmt sie an (`slice-harness-ci-workflows` §6 Frage 3).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-erster-release-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
