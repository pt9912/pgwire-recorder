# Ein Gate übergeht einen Eintrag der Ablage still

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Gate, das Dateien einer Ablage liest, sieht einen Eintrag nicht, den es lesen
müsste — ein nicht lesbares Verzeichnis, eine nicht lesbare Datei, einen Symlink —
und endet grün oder mit einem Ausgang, der etwas anderes meint, statt den Eintrag als
Befund zu melden.
