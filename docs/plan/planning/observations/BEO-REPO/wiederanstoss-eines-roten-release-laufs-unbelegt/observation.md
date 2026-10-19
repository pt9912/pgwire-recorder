# Der Wiederanstoß eines roten Release-Laufs auf GitHub ist nie gelaufen

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`docs/maintainer/releasing.md` §11 beschreibt den Wiederanstoß eines roten Laufs („Re-run failed jobs“ derselben Ausführung oder ein erneuter Lauf desselben Tags) und sagt, die Schritte seien so geschrieben, dass der zweite Lauf nichts verändert, was der erste hinterließ. Die Gegenprobe zeigt das an Attrappen (`p11-*`); auf GitHub mit echten Registries ist der Weg nie gegangen worden, weil der erste Lauf eines Tags grün endete. Die Zusage trägt, solange kein Lauf rot endet; ob die Tabelle der Fehler dann stimmt, zeigt erst dieser Fall.

Abgrenzung zu `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder`: Dort altert ein Pin oder Werkzeug; hier ist ein beschriebener Weg ungegangen.
