# CO-001: Die Pflicht zum vollen Lauf der Gegenproben am getaggten Commit ruht, bis die Gegenproben unter zehn Minuten laufen

**Status:** Entwurf. Empfehlung des Architect; die Entscheidung des Nutzers steht aus. Mit seinem Ja (Datum und Wortlaut hier nachzutragen) wird der Carveout Aktiv und in den Index `docs/plan/carveouts/README.md` aufgenommen.

**Datum angelegt:** 2026-10-21. **Letzte Prüfung:** 2026-10-21.

**Betroffenes Gate:** Punkte 2 und 3 der Freigabe-Checkliste (`SPEC-067` Punkt 17): `make gegenproben` (Gegenproben-Container) und `make gegenproben-host` (Host) am getaggten Commit, geregelt in `SPEC-060` Punkt 14. Die Schwelle selbst (jede Gegenprobe grün am getaggten Commit) bleibt im Vertrag stehen; der Carveout ist die befristete Ausnahme für `v0.1.1`.

**Geltungsbereich:** der Eintrag der Freigabe `docs/maintainer/freigaben/0.1.1.md`, Punkte 2 und 3, und der Folge-Tag `v0.1.2`, falls er nach dem Ja zu `v0.1.1` nötig wird (Entscheidung des Nutzers vom 2026-10-20). Kein anderer Tag.

**Folge-Slice:** [`slice-harness-gegenproben-schneller`](../planning/open/slice-harness-gegenproben-schneller.md) (nimmt `release-gegenprobe` und `test-container-gegenprobe` auf) und `slice-harness-gegenproben-grenze-rest` (Vorschlag des Architect, vom Planner anzulegen: nimmt die übrigen Gegenproben auf). Beide tragen die Sendung in §1 und DoD mit der Kennung `CO-001`, bevor der Carveout Aktiv wird (`AGENTS.md` §3.13).

Regeln: Baseline-Regelwerk `modul-07-carveouts.md` §Ziel-Form: Carveout.

---

## Begründung

Der Nutzer hat am 2026-10-21 entschieden: „die gegenproben werden zeitlich begrenzt. höchstens 10 minuten“ ([ADR-0055](../adr/0055-gegenproben-haben-ein-zeitbudget.md), Accepted). Gemessen am Stand vor der Lieferung braucht allein `make release-gegenprobe` 25 bis 30 Minuten auf ruhigem Host und Stunden unter Last; ein voller Lauf von `make gegenproben` liegt darüber. Die Pflicht aus `SPEC-060` Punkt 14 und der Checkliste verlangt diesen vollen Lauf am getaggten Commit. Bis die Gegenproben aufgenommen sind (Tabelle in `SPEC-060` Punkt 16), widersprechen sich Pflicht und Zeitgrenze; der Vertrag senkt die Schwelle nicht still (`AGENTS.md` §3.6), die Ausnahme steht hier mit Trigger.

## Ersatz für `v0.1.1` (Empfehlung des Architect)

Der Eintrag der Freigabe belegt statt der vollen Läufe, je Beleg mit Befehl, Hash, Datum, Dauer und Ausgang:

1. **Ablage des Gegenproben-Containers** (nichts zu warten): die Dateien `gegenproben-*-<Hash>` per `ssh`, wörtlich mit `ergebnis=`, `beginn`, `ende` oder der Marke `….laeuft`. Ein grünes Ergebnis für den Hash erfüllt Punkt 2 wie bisher; ein laufendes oder fehlendes wird nicht abgewartet.
2. **Die unveränderte Kopie der Gegenprobe des Tags:** `GP_MUT='^ok$' make release-gegenprobe` (rund zwei Minuten auf ruhigem Host) hält `release.yml` und `release.sh` des getaggten Commits gegen jeden Fall; `GP_MUT='^kein-treffer$' make test-container-gegenprobe` (rund 22 s) fährt allein die unveränderte Kopie. Beide am getaggten Hash, ruhiger Host (Last im Eintrag).
3. **Host (Punkt 3):** Der Diff der Dateien von `container-start` (`tools/harness/container-start*`, `harness/mk/container-start*`, `harness/vertraege/SPEC-058-container-start.md`) zwischen dem Hash der letzten grünen Fahrt von `make gegenproben-host` (Datum und Hash aus dem Bericht dieser Fahrt) und dem getaggten Hash ist leer. Ist er es nicht, oder gibt es keine belegte grüne Fahrt, steht das als Grenze im Eintrag und der Nutzer entscheidet vor dem Tag.

**Was der Eintrag schreibt:** Punkt 2: „nicht voll gefahren, Carveout CO-001 (Aktiv seit <Datum der Entscheidung des Nutzers>); Ersatz: <Belege 1 und 2>; Grenze: die Mutationen der Gegenproben (Wirksamkeit) und die übrigen Gegenproben am Hash sind nicht gelaufen.“ Punkt 3: „nicht gefahren, Carveout CO-001; Ersatz: <Diff-Ausgabe>; Grenze: ein Lauf am Host am Hash steht aus.“ Nie „grün“ ohne Beleg.

## Auflösungs-Trigger

Die Tabelle in `SPEC-060` Punkt 16 führt keine Gegenprobe mehr im Zustand `nicht gemessen`, und beide Folge-Slices liegen in `done/`. Dann fährt die Freigabe die Punkte 2 und 3 wieder voll, in unter zehn Minuten je Gegenprobe, und dieser Carveout wandert nach `done/`. Prüf-Termin: bei jedem Wellen-Audit und vor jedem weiteren Tag.

## Was der Nutzer entscheiden muss

- Ob der Carveout für `v0.1.1` (und `v0.1.2`) gilt (Empfehlung: ja). Die Alternativen: der volle Lauf trotz der Stunden, oder der Tag erst nach der Lieferung der beiden Folge-Slices.
- Ob die Ersatzbelege oben genügen oder er einen dritten verlangt (etwa einen Lauf einer ausgewählten Gegenprobe, deren Werkzeug seit der letzten grünen Fahrt geändert wurde).
