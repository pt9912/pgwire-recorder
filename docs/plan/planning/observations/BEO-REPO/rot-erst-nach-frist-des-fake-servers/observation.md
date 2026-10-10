# Ein Mutant wird erst nach der Frist des Fake-Servers rot

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test mit einem Fake-Server wartet auf ein Ereignis, das der Prüfling unter einer Mutation
nicht auslöst, mit einer eigenen Frist (Literal nach `SPEC-038`). Der Test ist aus dem
richtigen Grund rot, aber erst nach der Frist, bei mehreren Fällen und Mutanten nach
Minuten. Ein Fake-Server, der das Ende der Verbindung als Ereignis meldet, macht dieselbe
Mutation sofort rot.

Abgrenzung zu `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`: Dort fehlt die Frist, der Lauf
endet am Zeitlimit des Pakets; hier ist die Frist da und zu lang für die Rückmeldung.
