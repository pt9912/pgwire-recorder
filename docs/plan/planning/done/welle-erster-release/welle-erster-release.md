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

Der erste Release ist veröffentlicht: die Binaries für die Zielplattformen und das Docker/OCI-Image in `ghcr.io` und `docker.io`. Vor der Veröffentlichung läuft die Integration des Repos in GitHub Actions (CI auf jeden Pull-Request und Push; Auftrag des Nutzers vom 2026-10-18), und die Workflow-Dateien, darunter die der Veröffentlichung, bestehen ein Gate. Die Nachtläufe (Gegenproben, upstream) gehören nicht dazu: Sie folgen nach dem Release in [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) (Entscheidung des Nutzers vom 2026-10-18, „Erst das Release voranbringen“). Die Homebrew-Formel im Tap und Abnahmeszenario 11 ([`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)) gehören seit der Entscheidung des Nutzers vom 2026-10-11 zu [welle-v1-nachtrag](../welle-v1-nachtrag.md), die auf diese Welle folgt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-v1-abschluss](welle-v1-abschluss/welle-v1-abschluss.md) done, und `slice-harness-coverage` liegt in `done/` (Abnahmeszenario 17, M3 vor M4; Entscheidung des Nutzers vom 2026-10-08), dazu `slice-harness-devcontainer-u-boot-0-10` in `done/` (wellenlos; Entscheidung des Nutzers vom 2026-10-17: lang laufende Slices laufen im Dev-Container, damit `main` nicht durch Implementer-Läufe gesperrt ist; `slice-harness-devcontainer-einrichten` ist nicht Bedingung, Entscheidung des Nutzers vom 2026-10-17).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Der erste stabile Tag ist nach der Freigabe-Checkliste veröffentlicht (Binaries mit `SHA256SUMS` am Release, Image in beiden Registries mit gleichem Index-Digest, Smoke aus Abnahmeszenario 9 bestanden) — das *Mehr* gegenüber den Slice-DoDs, deren Pipeline nur der Trockenlauf belegt (ein Probe-Tag ist nicht vorgesehen, Entscheidung des Nutzers vom 2026-10-19). Die Pipeline-Läufe sind nachgewiesen: ein grüner Lauf von `ci.yml` auf dem Stand des Tags und der Lauf der Veröffentlichung; ein Nachtlauf ist keine Bedingung (Nachtläufe: [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md)). Abnahmeszenario 11 ist keine Bedingung dieser Welle; [welle-v1-nachtrag](../welle-v1-nachtrag.md) weist es nach.
- `make gates` grün.
- Closure-Notiz in `welle-erster-release-results.md`; sie trägt den ausgefüllten Eintrag der Freigabe-Checkliste des ersten Releases (Entscheidung des Nutzers vom 2026-10-09, geliefert von `slice-erster-release-freigabe`).

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-harness-ci-workflows | Gate der Workflows: Pins, Rechte und Zeitgrenzen der Dateien unter `.github/workflows/` | kein Anspruch des Produkts |
| slice-harness-ci-range-konfiguration | Die Range-Prüfungen wirken: `.d-check.yml` führt `commits` und `vcs`, eine Gegenprobe zeigt es | kein Anspruch des Produkts |
| slice-harness-ci | Die CI auf GitHub: `ci.yml` fährt `make gates` und prüft die Commit-Range | kein Anspruch des Produkts |
| slice-erster-release-veroeffentlichung | Veröffentlichung von Tag und Binaries | [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |
| slice-erster-release-veroeffentlichung-image | Veröffentlichung des Images | [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität) |
| slice-erster-release-freigabe | Freigabe und erster echter Tag | [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-v1-nachtrag](../welle-v1-nachtrag.md) und Welle [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) (deren Start-Trigger ist diese Welle in `done/`).
- Wird blockiert von: Welle [welle-v1-abschluss](welle-v1-abschluss/welle-v1-abschluss.md) und den wellenlosen Slices `slice-harness-abdeckung-gate`, `slice-harness-coverage` und `slice-harness-devcontainer-u-boot-0-10` (nicht an `slice-harness-devcontainer-einrichten`). Für die ersten zwei Slices der CI zusätzlich `slice-harness-gegenproben-aus-gates` (Start von `slice-harness-ci-workflows` und `slice-harness-ci`); er gehört zu [welle-gegenproben-und-mutanten](../welle-gegenproben-und-mutanten.md) und steht dort vorn (Drift-Log der Roadmap, 2026-10-18); er ist Bedingung der genannten Slices, nicht des Starts dieser Welle. `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung` ist keine Bedingung (Entscheidung des Nutzers vom 2026-10-18).
- Innerhalb der Welle: `slice-harness-ci-range-konfiguration` setzt `slice-harness-ci-workflows` voraus (Reihenfolge, WIP-Limit 1), `slice-harness-ci` setzt `slice-harness-ci-workflows` (das Gate prüft `ci.yml`) und `slice-harness-ci-range-konfiguration` (die Range-Prüfungen müssen etwas prüfen) voraus, `slice-erster-release-veroeffentlichung` setzt `slice-harness-ci-workflows` und `slice-harness-ci` voraus (Pins und Härtung von `release.yml`), `slice-erster-release-veroeffentlichung-image` setzt `slice-erster-release-veroeffentlichung` voraus (die Endfassung von `release.yml` ergänzt dessen Fassung um Zeilen), `slice-erster-release-freigabe` setzt `slice-erster-release-veroeffentlichung` und `slice-erster-release-veroeffentlichung-image` voraus (der Scan läuft in deren Pipeline zwischen Bau und Push des Images).

**Reihenfolge** (Entscheidung des Nutzers vom 2026-10-18, „Erst das Release voranbringen“ mit der Antwort „Minimal: ci-workflows, ci, veroeffentlichung, freigabe“: der kritische Pfad zum ersten Release besteht aus diesen vier Slices, seit dem Schnitt von `slice-harness-ci` am 2026-10-18 (Empfehlung des Planners, vom Nutzer nicht bestätigt) mit `slice-harness-ci-range-konfiguration` aus fünf; Auftrag des Nutzers vom 2026-10-18: die CI vor die Veröffentlichung; Schnitt von `slice-harness-ci-workflows` und `slice-harness-ci`: Empfehlung des Planners, nicht entschieden; Entscheidungen des Nutzers vom 2026-10-08: Wellen vor Harness, M3 vor M4; Schnitt von `slice-erster-release-veroeffentlichung` vom 2026-10-09 und in Tag und Binaries gegen Image vom 2026-10-19 (Empfehlung des Planners, vom Nutzer nicht bestätigt; Drift-Log der Roadmap); Umplanung von `slice-erster-release-homebrew-nachweis` nach `welle-v1-nachtrag` vom 2026-10-11; WIP-Limit 1):

1. `slice-harness-ci-workflows` — Gate der Workflows (Start: `slice-harness-gegenproben-aus-gates` in `done/`).
2. `slice-harness-ci-range-konfiguration` — `.d-check.yml` führt `commits` und `vcs`, mit Gegenprobe (setzt 1 voraus; Schnitt aus `slice-harness-ci`, Drift-Log der Roadmap, 2026-10-18).
3. `slice-harness-ci` — `ci.yml` mit `make gates` und Commit-Range (setzt 1 und 2 voraus).
4. `slice-erster-release-veroeffentlichung` — Pipeline für Tag und Binaries (Fassung A des Vertrags `SPEC-067`), belegt mit dem Trockenlauf (setzt 1 bis 3 voraus).
5. `slice-erster-release-veroeffentlichung-image` — Pipeline für das Image (Fassung B), belegt mit dem Trockenlauf (setzt 4 voraus; Schnitt vom 2026-10-19, Drift-Log der Roadmap).
6. `slice-erster-release-freigabe` — Scan, Freigabe-Checkliste und der erste echte Tag (das Freshness-Audit ist optional und kein Punkt der Checkliste, Entscheidung des Nutzers vom 2026-10-19) (setzt 4 und 5 voraus).

Die Nachtläufe (`slice-harness-ci-gegenproben`, `slice-harness-ci-freshness`, `slice-harness-ci-freshness-images`) liegen seit dem 2026-10-18 in [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) hinter dieser Welle.

Vor dieser Welle liegen die Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; wellenlos, Entscheidung des Nutzers vom 2026-10-13), `slice-harness-integration-volumes`, `slice-harness-abdeckung-gate` und `slice-harness-coverage`: Mit Coverage wird Abnahmeszenario 17 nachweisbar, M3 ist damit vor M4 erreicht. Hinter Coverage stehen wellenlos `slice-harness-devcontainer-u-boot-0-10` (vor der Closure von welle-v1-abschluss und vor den Slices dieser Welle) und danach `slice-harness-devcontainer-einrichten` (hinter jenem, ohne Bindung an Closure und Release; Entscheidung des Nutzers vom 2026-10-17). Danach folgt [welle-v1-nachtrag](../welle-v1-nachtrag.md), dann die übrigen Harness-Slices, beginnend mit `slice-harness-commit-struktur-id`; ihre Reihenfolge steht in deren §4.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Die Homebrew-Formel im Tap und ihr Nachweis (Abnahmeszenario 11) — `slice-erster-release-homebrew-nachweis` in [welle-v1-nachtrag](../welle-v1-nachtrag.md) (Entscheidung des Nutzers vom 2026-10-11); der Release ist ohne sie vollständig, weil kein Release-Artefakt von der Formel abhängt.
- Weitere Plattformen oder Registries — nicht zugesichert.
- Branch-Schutz („required status check“), Geheimnisse und Sichtbarkeit des Repos auf GitHub — Einstellungen des Betreibers, aus dem Klon nicht auditierbar; die Welle nennt sie als Handlung, ohne sie zu behaupten.
- Aktualität und Gültigkeit der Action-SHA und der Werkzeug-Versionen in den Workflows — kein Slice der Welle nimmt sie an (`slice-harness-ci-workflows` §6 Frage 3).
- Die Nachtläufe auf GitHub (Gegenproben, Audit gegen upstream, Image-Pins) — [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) (Entscheidung des Nutzers vom 2026-10-18); das Freshness-Audit ist optional, von Hand und kein Punkt der Freigabe-Checkliste (Entscheidung des Nutzers vom 2026-10-19), ein Nachtlauf ist weder Beleg noch Bedingung.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-erster-release-results.md` als Geschwister im Ruheort `done/`.
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
