# Ein Push auf den Hauptzweig erfolgt, bevor `make gates` für den Stand grün war

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Stop-Hook bindet das Ende einer Sitzung an einen Nachweis des vollen Laufs von `make gates` (oder die Positivliste inerter Pfade); der Push auf `main` ist nicht gebunden. Wird ein Stand gepusht, der nicht inert ist, bevor die Gates für ihn grün waren, ist das Ergebnis der Gates erst nachträglich bekannt (`ci.yml` läuft nach dem Push). Befund, keine Schuldzuweisung: Die Gates waren danach grün, aber der Push sagte vorher etwas, was noch nicht belegt war.

Abgrenzung zu `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`: Dort blockiert der Hook einen Implementer im Hauptbaum; hier fehlt die Bindung des Pushs an dasselbe Ergebnis. Abgrenzung zu `BEO-REPO/ci-ohne-branch-schutz`: Dort geht es um die Einstellung auf GitHub; hier um den Push vor dem lokalen Ergebnis.
