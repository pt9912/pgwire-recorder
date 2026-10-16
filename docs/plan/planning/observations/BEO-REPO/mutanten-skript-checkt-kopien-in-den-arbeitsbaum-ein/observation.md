# Ein Mutanten-Skript legt Kopien im Arbeitsbaum an, und `git add -A` checkt sie ein

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Lauf fährt Mutanten mit einem selbst geschriebenen Skript, das eine Kopie des Quellbaums anlegt. Wird das Skript mit einem Namen aus mehreren Wörtern als ein Argument aufgerufen, zerfällt der Name, und die Kopien entstehen im Arbeitsbaum; ein folgendes `git add -A` nimmt sie in den Commit auf, und die History trägt Hunderte Dateien als Zugang, bis ein späterer Commit sie entfernt. Zweite Ursache derselben Art nach Angabe des Implementers im Bericht an den Aufrufer: Ein Scratchpad, das mehrere Agenten teilen, ließ ein Skript (`mut.sh`) vom anderen Lauf überschreiben. Beide Ursachen betreffen Hilfsskripte der Agenten, nicht den Produktcode.

Abgrenzung zu `BEO-REPO/pkill-f-toetet-die-eigene-shell`: Dort scheitert ein Aufräumschritt; hier bleibt Fremdes im Commit. Abgrenzung zu `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`: Dort erreicht der Mutant den Test nicht; hier erreicht eine Kopie den Commit.
