# CO-001: Die Pflicht zum vollen Lauf der Gegenproben am getaggten Commit ruht, bis die Gegenproben unter zehn Minuten laufen

**Status:** Aktiv seit 2026-10-21. Der Nutzer bekam den Entwurf dieser Datei mit seinem Wortlaut vorgelegt und antwortete am 2026-10-21 (Entscheidung des Nutzers, wörtlich): „ich gebe alles frei, damit wir jetzt endlich diesen brew release machen können“ und „bitte auch den tag fahren“. Die Antwort nennt `CO-001` nicht beim Namen; als Ja zu diesem Carveout gilt sie, weil sie auf den vorgelegten Entwurf folgte und ausdrücklich alles freigab, was für den Tag `v0.1.1` offen war. Eine Frage an den Nutzer bleibt in diesem Text offen (Abschnitt *Was der Nutzer noch entscheiden kann*): ob er einen dritten Ersatzbeleg verlangt.

**Datum angelegt:** 2026-10-21. **Letzte Prüfung:** 2026-10-21. **Nächste Prüfung:** vor jedem weiteren Tag, spätestens am 2027-01-19 (90 Tage, `modul-07-carveouts.md` §Carveout-Audit-Slice).

**Betroffenes Gate:** Punkte 2 und 3 der Freigabe-Checkliste (`SPEC-067` Punkt 17): `make gegenproben` (Gegenproben-Container) und `make gegenproben-host` (Host) am getaggten Commit, geregelt in `SPEC-060` Punkt 14. Die Schwelle (jede Gegenprobe grün am getaggten Commit) setzen `SPEC-060` Punkt 14 und [ADR-0048](../adr/0048-gegenproben-verlassen-die-gate-kette.md), die Obergrenze von zehn Minuten [ADR-0055](../adr/0055-gegenproben-haben-ein-zeitbudget.md); der Carveout setzt keine Schwelle, er nimmt die Punkte 2 und 3 für die genannten Tags aus. Sichtbar wird er im Eintrag der Freigabe, der ihn je Punkt nennt (Abschnitt *Ersatz*); kein Gate-Output ändert sich, `make gegenproben` und `make gegenproben-host` bleiben unverändert.

**Geltungsbereich (Frist):** Der Carveout endet mit dem Tag `v0.1.2` oder, wenn dieser nicht nötig wird, mit `v0.1.1`; er trägt keinen dritten Tag. Gegenstand: der Eintrag der Freigabe `docs/maintainer/freigaben/0.1.1.md`, Punkte 2 und 3, und der Folge-Tag `v0.1.2`, falls er nach dem Ja zu `v0.1.1` nötig wird (Entscheidung des Nutzers vom 2026-10-20). Kein anderer Tag; ein weiterer Tag braucht ein neues Ja und, solange die Auflösung unten nicht eingetreten ist, eine Verlängerung dieses Carveouts durch den Nutzer.

**Folge-Slice:** `slice-harness-gegenproben-schneller` (nimmt `release-gegenprobe` und `test-container-gegenprobe` auf) und `slice-harness-gegenproben-grenze-rest` (`slice-harness-gegenproben-grenze-rest`; angelegt vom Planner am 2026-10-21: misst die übrigen Gegenproben und nimmt sie auf). Beide tragen die Sendung in §1 und DoD mit der Kennung `CO-001` (`grep` nach `CO-001` findet sie, `AGENTS.md` §3.13).

Regeln: Baseline-Regelwerk `modul-07-carveouts.md` §Ziel-Form: Carveout.

---

## Begründung

Der Nutzer hat am 2026-10-21 entschieden (Entscheidung des Nutzers, wörtlich): „die gegenproben werden zeitlich begrenzt. höchstens 10 minuten“ ([ADR-0055](../adr/0055-gegenproben-haben-ein-zeitbudget.md), Accepted). Gemessen am Stand vor der Lieferung braucht allein `make release-gegenprobe` 25 bis 30 Minuten auf ruhigem Host und Stunden unter Last; ein voller Lauf von `make gegenproben` liegt darüber. Die Pflicht aus `SPEC-060` Punkt 14 und der Checkliste verlangt diesen vollen Lauf am getaggten Commit. Bis die Gegenproben aufgenommen sind (Tabelle in `SPEC-060` Punkt 16), widersprechen sich Pflicht und Zeitgrenze; der Vertrag senkt die Schwelle nicht still (`AGENTS.md` §3.6), die Ausnahme steht hier mit Trigger.

## Ersatz für `v0.1.1` und `v0.1.2` (Vorschlag des Architect; vom Nutzer mit der Freigabe vom 2026-10-21 nicht geändert)

Der Eintrag der Freigabe belegt statt der vollen Läufe, je Beleg mit Befehl, Hash, Datum, Dauer und Ausgang:

1. **Ablage des Gegenproben-Containers** (nichts zu warten): die Dateien `gegenproben-*-<Hash>` per `ssh`, wörtlich mit `ergebnis=`, `beginn`, `ende` oder der Marke `….laeuft`. Ein grünes Ergebnis für den Hash erfüllt Punkt 2 wie bisher; ein laufendes oder fehlendes wird nicht abgewartet.
2. **Die unveränderte Kopie der Gegenprobe des Tags:** `GP_MUT='^ok$' make release-gegenprobe` (rund zwei Minuten auf ruhigem Host) hält `release.yml` und `release.sh` des getaggten Commits gegen jeden Fall; `GP_MUT='^kein-treffer$' make test-container-gegenprobe` (rund 22 s) fährt allein die unveränderte Kopie. Beide am getaggten Hash, ruhiger Host (Last im Eintrag).
3. **Host (Punkt 3):** Der Diff der Dateien von `container-start` (`tools/harness/container-start*`, `harness/mk/container-start*`, `harness/vertraege/SPEC-058-container-start.md`) zwischen dem Hash der letzten grünen Fahrt von `make gegenproben-host` (Datum und Hash aus dem Bericht dieser Fahrt) und dem getaggten Hash ist leer. Ist er es nicht, oder gibt es keine belegte grüne Fahrt, steht das als Grenze im Eintrag und der Nutzer entscheidet vor dem Tag.

**Was der Eintrag schreibt:** Punkt 2: „nicht voll gefahren, Carveout CO-001 (Aktiv seit 2026-10-21); Ersatz: <Belege 1 und 2>; Grenze: die Mutationen der Gegenproben (Wirksamkeit) und die übrigen Gegenproben am Hash sind nicht gelaufen.“ Punkt 3: „nicht gefahren, Carveout CO-001; Ersatz: <Diff-Ausgabe>; Grenze: ein Lauf am Host am Hash steht aus.“ Nie „grün“ ohne Beleg.

## Auflösungs-Trigger

Die Tabelle in `SPEC-060` Punkt 16 führt jede Gegenprobe als `aufgenommen` (keine im Zustand `nicht gemessen`, keine gemessen über 600 s mit offenem Nehmer), und jeder Slice, den die Tabelle als Nehmer nennt, liegt in `done/`: `slice-harness-gegenproben-schneller`, `slice-harness-gegenproben-grenze-rest` und die Slices, die dieser für Gegenproben über 600 s anlegt. Dann fährt die Freigabe die Punkte 2 und 3 wieder voll, in unter zehn Minuten je Gegenprobe, und dieser Carveout wandert nach `done/` (reiner `git mv`; der Index `README.md` zieht im Commit davor nach). Prüf-Termin: bei jedem Wellen-Audit und vor jedem weiteren Tag.

## Was der Nutzer noch entscheiden kann

- Ob er zu den Ersatzbelegen oben einen dritten verlangt (etwa einen Lauf einer ausgewählten Gegenprobe, deren Werkzeug seit der letzten grünen Fahrt geändert wurde).
