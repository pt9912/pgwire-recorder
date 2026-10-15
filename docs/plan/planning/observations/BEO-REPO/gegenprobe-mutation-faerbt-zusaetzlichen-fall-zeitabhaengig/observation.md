# Eine Mutation der Gegenprobe färbt einmal einen Fall außerhalb ihrer Gruppe rot

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Im ersten vollen `make gates` am Stand der Zusage war `postgres-versionen-gegenprobe` rot: Die Mutation
`p7-int-treiber` färbte zusätzlich den Fall `p7-sigTERM-gruppe` rot, den sie nicht färben soll
(`SPEC-053`: Mutationen färben nur die Fälle ihrer Gruppe). `make test-postgres-versionen-gegenprobe`
allein danach und der zweite volle `make gates` waren grün; der Slice änderte weder Gate noch
Gegenprobe. Die Flanke ist zeitabhängig und nicht untersucht.

Abgrenzung zu `BEO-REPO/test-faellt-unter-cpu-last-aus`: Dort fällt ein Test des Produkts unter Last aus;
hier ist es eine Gegenprobe des Gates, und der Fehler ist eine zusätzlich gefärbte Mutation. Abgrenzung zu
`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`: Dort geht es um die Dauer, hier um den Ausgang.
