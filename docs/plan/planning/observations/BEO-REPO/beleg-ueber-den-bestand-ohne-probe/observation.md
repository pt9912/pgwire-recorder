# Ein Beleg über den Bestand nennt keine Probe, die ihn zu einem Stand wieder prüft

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Entscheidung stützt sich auf eine Aussage über den Bestand („kein langsames Gate liest diese
Pfade“) und führt sie als Beleg in der Spezifikation. Der Beleg ist zu einem Stand wahr gewesen oder nie
gefahren worden; eine Probe, die ihn hält, steht nicht dabei. Das Review sucht den Leser und findet ihn
(`SPEC-052` Punkt 13: `make hook-gegenprobe` las den Planungsbaum aus dem Index; F-713, MEDIUM, gefahren
durch Verschieben der Slices in einer Kopie).

Abgrenzung zu `BEO-REPO/beleg-traegt-zusage-aus-anderem-grund` (ein Test oder Mutant trägt die Zusage aus
einem anderen Grund) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (eine Randform des Formats ist
offen): Hier ist die Aussage über den Bestand der Beleg selbst, und sie war falsch oder lückenhaft.
