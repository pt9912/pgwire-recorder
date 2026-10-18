# Eine Anforderung ohne Zeile in der Abdeckungstabelle bleibt für den Abgleich unsichtbar

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`make abdeckung-check` hält die Tabellen `docs/user/abdeckung-*.md` gegen die Deklarationen der Tests und Gegenproben, nicht gegen die Liste der Anforderungen des Lastenhefts. Eine Anforderung, die keine Deklaration trägt (zum Beispiel `LH-FA-15`, `LH-QA-04`, `LH-FA-16`, `LH-QA-03`), hat keine Zeile, und der Abgleich meldet keine Lücke, obwohl ein Abnahmeszenario sie nennt.
