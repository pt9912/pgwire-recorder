# Die White-Box-Liste eines Umstellungs-Slice steht vor dem Code nur als Namenssuche

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Slice stellt Tests auf Black-Box-Pakete um. §6 nennt die Zugriffe der Tests auf
Unexportiertes vor dem Code nur aus einer Suche nach Namen, ungemessen und mit
Fehltreffern. Gemessen (per AST) wird erst im Code-Commit. Was die Messung neu zeigt,
prüfen Architect und Plan vor dem Code nicht. Liegt ein neuer Zugriff außerhalb der vorab
entschiedenen Klassen, wird er im Code entschieden.
