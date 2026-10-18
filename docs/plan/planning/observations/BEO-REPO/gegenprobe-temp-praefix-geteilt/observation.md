# Zwei Läufe derselben Gegenprobe teilen das Präfix ihres Temp-Verzeichnisses

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`container-start-gegenprobe` legt ihre Verzeichnisse unter `/tmp` mit dem Präfix `pgr-cs-` an und räumt Verzeichnisse dieses Präfixes am Ende auf. Läuft eine zweite Fahrt am selben Host gleichzeitig (ein anderer Worktree, ein anderer Agent), kann die eine die Verzeichnisse der anderen entfernen oder die Aufräumprüfung der anderen stören: Im ersten Lauf der Verifikation von `slice-harness-gegenproben-aus-gates` verschwand das Temp-Verzeichnis während des Laufs (`Cannot find module '/tmp/pgr-cs-…/pruef.js'`), und 110 Fälle wurden „an der unveränderten Kopie rot“. Der zweite Lauf desselben Stands war grün; im Prozessprotokoll lief zur selben Zeit ein Lauf des Skripts aus einem anderen Worktree mit demselben Präfix. Der Vertrag der Gegenprobe sagt nichts zu gleichzeitigen Läufen am selben Host.
