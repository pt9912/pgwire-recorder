# Carveouts — pgwire-recorder

Aktive Carveouts mit Auflösungs-Trigger. Aufgelöste Carveouts wandern nach
`done/` (reiner `git mv`).

## Aktive Carveouts

| ID | Titel | Gate | Trigger | Folge-Slice |
|---|---|---|---|---|
| [CO-001](CO-001-gegenproben-am-getaggten-commit.md) | Die Pflicht zum vollen Lauf der Gegenproben am getaggten Commit ruht für `v0.1.1` und `v0.1.2` | Punkte 2 und 3 der Freigabe-Checkliste (`make gegenproben`, `make gegenproben-host`; Schwelle: [ADR-0048](../adr/0048-gegenproben-verlassen-die-gate-kette.md), Obergrenze: [ADR-0055](../adr/0055-gegenproben-haben-ein-zeitbudget.md)) | Die Tabelle in `SPEC-060` Punkt 16 führt jede Gegenprobe als `aufgenommen`, und jeder Nehmer der Tabelle liegt in `done/` | `slice-harness-gegenproben-schneller`, `slice-harness-gegenproben-grenze-rest` |

## Aufgelöste Carveouts

(noch keine)

## Konventionen

- Jeder aktive Carveout braucht: Trigger, Folge-Slice, letzten Prüf-Termin.
- Bei Welle-Closure: Carveout-Audit zwingend — welche gültig, welche aufgelöst?
- Siehe Baseline-Regelwerk `modul-07-carveouts.md`.
