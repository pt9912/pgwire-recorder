# Ein Fall einer Gegenprobe startet die Gegenprobe selbst und läuft ohne Ende

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Gegenprobe, die ihr Werkzeug in Kopien aufruft, kann einen Fall enthalten, der ein Kind-Werkzeug oder eine Kopie der Gegenprobe ausführt. Führt die Kopie dieselben Fälle wieder aus, startet sie sich selbst, und der Lauf verzweigt, bis der Rechner voll ist. In `container-start-gegenprobe` rief die erste Fassung der Nacharbeit sich über den Fall `p15-sonden-raeumen` selbst auf (die Kind-Gegenprobe führte den Fall wieder aus). Behoben durch die Ausnahme der Fälle, die eine Kind-Gegenprobe starten, und einen Schutz im Kind (Fall `p15-rekursion`, Mutation `p15-rekursion-schutz`); alle Läufe seitdem mit `timeout` und mit Beobachtung der Prozesszahl.

Abgrenzung zu `BEO-REPO/mutanten-skript-checkt-kopien-in-den-arbeitsbaum-ein`: Dort schreibt das Skript Kopien in den Arbeitsbaum; hier startet es sich selbst. Abgrenzung zu `BEO-REPO/pkill-f-toetet-die-eigene-shell`: Dort trifft ein Befehl die eigene Shell; hier vermehrt sich der Lauf.
