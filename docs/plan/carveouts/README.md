# Carveouts — pgwire-recorder

Aktive Carveouts mit Auflösungs-Trigger. Aufgelöste Carveouts wandern nach
`done/` (reiner `git mv`).

## Aktive Carveouts

| ID | Titel | Gate | Trigger | Folge-Slice |
|---|---|---|---|---|
| [CO-001](CO-001-gegenproben-am-getaggten-commit.md) | Die Pflicht zum vollen Lauf der Gegenproben am getaggten Commit ruht für `v0.1.1` und `v0.1.2` | Punkte 2 und 3 der Freigabe-Checkliste (`make gegenproben`, `make gegenproben-host`; Schwelle: [ADR-0048](../adr/0048-gegenproben-verlassen-die-gate-kette.md), Obergrenze: [ADR-0055](../adr/0055-gegenproben-haben-ein-zeitbudget.md)) | Die Tabelle in `SPEC-060` Punkt 16 führt jede Gegenprobe als `aufgenommen`, und jeder Nehmer der Tabelle liegt in `done/` | `slice-harness-gegenproben-schneller`, `slice-harness-gegenproben-grenze-rest` |

## Aufgelöste Carveouts

| ID | Titel | Gate | Trigger | Folge-Slice |
|---|---|---|---|---|
| [CO-002](done/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md) | Die Pflicht zum vollen Lauf der Gegenproben am getaggten Commit ruhte für `v0.2.0` | Punkte 2 und 3 der Freigabe-Checkliste (wie `CO-001`) | Closure von `slice-release-v0-2-0` (Tag gesetzt): aufgelöst am 2026-10-24; Beleg `docs/maintainer/freigaben/0.2.0.md` Punkte 2 und 3 und Verifikationsbericht `docs/reviews/2026-10-24-verifikation-slice-release-v0-2-0.md`; die Ursache führt `CO-001` | `slice-release-v0-2-0` (verbraucht) |

## Konventionen

- Jeder aktive Carveout braucht: Trigger, Folge-Slice, letzten Prüf-Termin.
- Bei Welle-Closure: Carveout-Audit zwingend — welche gültig, welche aufgelöst?
- Siehe Baseline-Regelwerk `modul-07-carveouts.md`.
