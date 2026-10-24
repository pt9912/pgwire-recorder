# CO-003: Die Pflicht zum vollen Lauf der Gegenproben am Wellenende ruht, bis die Gegenproben unter zehn Minuten laufen

**Status:** Aktiv seit 2026-10-24. **Entscheidung des Nutzers vom 2026-10-24**, auf die Frage nach der Abweichung am Ende von `welle-v1-nachtrag` wörtlich: „CO-003 anlegen (Recommended)“. Die Wiedergabe des Wortlauts stammt von der Koordination und ist nicht nachgefahren.

**Datum angelegt:** 2026-10-24. **Letzte Prüfung:** 2026-10-24. **Nächste Prüfung:** bei jeder Welle-Closure vor dem Trigger, spätestens am 2027-01-22 (90 Tage, `modul-07-carveouts.md` §Carveout-Audit-Slice).

**Betroffenes Gate:** `make gegenproben` (Gegenproben-Container) und `make gegenproben-host` (Host) am Ende einer Welle, geregelt in `SPEC-060` Punkt 14 und `AGENTS.md` §6 Punkt 6. Die Schwelle (jede Gegenprobe grün am Wellenende) setzen `SPEC-060` Punkt 14 und [ADR-0048](../adr/0048-gegenproben-verlassen-die-gate-kette.md), die Obergrenze von zehn Minuten [ADR-0055](../adr/0055-gegenproben-haben-ein-zeitbudget.md); der Carveout setzt keine Schwelle, er nimmt die Pflicht am Wellenende für den genannten Bereich aus. Kein Gate-Output ändert sich, `make gegenproben` und `make gegenproben-host` bleiben unverändert.

**Geltungsbereich (Frist):** das Wellenende von `welle-v1-nachtrag` und jedes weitere Wellenende, das vor dem Auflösungs-Trigger unten schließt. **Nicht gedeckt** sind die Tags: Die Pflicht am getaggten Commit (Punkte 2 und 3 der Freigabe-Checkliste, `SPEC-067` Punkt 17) deckt allein ein Carveout je Tag, wie `CO-001` und `CO-002` es tun; dieser Carveout verlängert sie nicht und ersetzt sie nicht.

**Folge-Slice:** `slice-harness-gegenproben-schneller` (nimmt `release-gegenprobe` und `test-container-gegenprobe` auf) und `slice-harness-gegenproben-grenze-rest` (misst die übrigen und nimmt sie auf). Beide tragen die Sendung in §1 und DoD mit der Kennung `CO-003` (`grep` nach `CO-003` findet sie, `AGENTS.md` §3.13); ihre Ausschlüsse treffen sie nicht (gelesen am 2026-10-24: kein Ausschluss nennt das Wellenende, beide schließen das Auflösen von `CO-001` vor dem Trigger und das Setzen eines Tags aus, nicht diesen Carveout). Neue Slices legt dieser Carveout nicht an; die Slices, die `slice-harness-gegenproben-grenze-rest` für Gegenproben über 600 s anlegt, tragen `CO-001` und lösen damit auch diesen Carveout mit auf.

Regeln: Baseline-Regelwerk `modul-07-carveouts.md` §Ziel-Form: Carveout.

---

## Begründung

Die Pflicht aus `SPEC-060` Punkt 14 verlangt am Ende jeder Welle den vollen Lauf beider Sammelziele. Die Zeitgrenze aus [ADR-0055](../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Entscheidung des Nutzers vom 2026-10-21, höchstens zehn Minuten je Gegenprobe) und `AGENTS.md` §6 Punkt 6 (Entscheidung des Nutzers vom 2026-10-22: nur der engste Fall mit `GP_MUT`, auch beim Verifier) machen ihn unmöglich, solange `release-gegenprobe` und `test-container-gegenprobe` nicht aufgenommen sind und die übrigen Gegenproben nicht gemessen sind (Tabelle in `SPEC-060` Punkt 16): Gemessen brauchen allein `make release-gegenprobe` 25 bis 30 Minuten auf ruhigem Host und Stunden unter Last und `make container-start-gegenprobe` 19 min 38 s. Der Vertrag senkt die Schwelle nicht still (`AGENTS.md` §3.6); die Ausnahme steht hier mit Trigger. Die Closure von `welle-v1-nachtrag` hat die Läufe nicht gefahren und nannte das eine Abweichung ohne Carveout; dieser Carveout benennt sie.

## Ersatz am Wellenende

Die Closure-Notiz der Welle belegt statt der vollen Läufe, je Beleg mit Befehl, Hash, Datum, Dauer und Ausgang:

1. **`make gates` auf `main`** am Stand nach dem Merge der Closure (Exit-Code, Dauer). Die Gegenproben gehören nicht zu dieser Kette (`SPEC-052`, [ADR-0048](../adr/0048-gegenproben-verlassen-die-gate-kette.md)); der Beleg trägt die Gates, nicht die Gegenproben.
2. **Die Gegenproben, die ein Slice der Welle änderte:** die Notiz nennt je Slice, welche Gegenprobe er änderte und wie er sie fuhr (§7 des Slice: voll oder mit `GP_MUT`). Eine Gegenprobe, die keiner änderte und die nicht gefahren ist, nennt die Notiz als nicht gefahren.
3. **Die Zeile der Notiz** („Gegenproben am Wellenende“) nennt den Carveout beim Namen und führt die Grenze: Die Mutationen der Gegenproben (Wirksamkeit) und die nicht geänderten Gegenproben am Hash sind nicht gelaufen. Nie „grün“ ohne Beleg.

## Auflösungs-Trigger

Derselbe wie bei [`CO-001`](CO-001-gegenproben-am-getaggten-commit.md): Die Tabelle in `SPEC-060` Punkt 16 führt jede Gegenprobe als `aufgenommen` (keine im Zustand `nicht gemessen`, keine gemessen über 600 s mit offenem Nehmer), und jeder Slice, den die Tabelle als Nehmer nennt, liegt in `done/`: `slice-harness-gegenproben-schneller`, `slice-harness-gegenproben-grenze-rest` und die Slices, die dieser für Gegenproben über 600 s anlegt. Dann fährt die Welle am Ende wieder voll, in unter zehn Minuten je Gegenprobe, und dieser Carveout wandert nach `done/` (reiner `git mv`; der Index `README.md` zieht im Commit davor nach), zusammen mit `CO-001`. Prüf-Termin: bei jeder Welle-Closure (Carveout-Audit) und vor jedem Tag.

## Geltungs-Konfiguration

Keine: Der Carveout ruht auf einem Verfahren (Closure-Text), nicht auf einer Gate-Konfiguration.
