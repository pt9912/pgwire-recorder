# Ein neues `docker run` ohne `--rm` legt das Leck der anonymen Volumes wieder an

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Skript oder Make-Fragment startet einen Container aus einem Image mit `VOLUME` ohne
`--rm` und ohne `-v` beim Entfernen. Die Gegenprobe des Aufräumens prüft nur die beiden
Läufer (`tools/test/run-integration-tests.sh`, `tools/harness/run-container.sh`); kein Gate
sieht das neue `docker run`, und das Leck zeigt sich erst an einer vollen Platte.

Abgrenzung zu `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`: Dort ist die Ursache
in den zwei Läufern behoben; hier geht es um künftige Stellen, die der Slice nicht kennt.
