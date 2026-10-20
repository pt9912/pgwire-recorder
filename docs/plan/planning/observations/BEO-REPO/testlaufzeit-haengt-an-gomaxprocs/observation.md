# Die Laufzeit eines Testpakets hängt an der Zahl der Kerne

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test hat rund hundert Unterfälle, die je bis zu 2 s auf „keine Antwort“ warten und mit `t.Parallel()` laufen; die Zeit des Pakets steht dann nicht fest, sondern wächst mit sinkender Kernzahl (gemessen: rund 20 s mit 20 Kernen, rund 96 s je Lauf mit `GOMAXPROCS=2`). Die Wartezeit ist der Beleg einer Abwesenheit; eine Beobachtung, die sie ersetzte, gibt es nicht. Keine Zeitgrenze von `make test` oder eines Gates ist dafür belegt; die Grenze des Jobs der CI ist 30 Minuten.

Abgrenzung zu `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`: Dort wächst die Zeit des Gates durch eine Gegenprobe; hier hängt sie an der Umgebung. Abgrenzung zu `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test`: Dort ersetzt eine feste Zeit eine mögliche Beobachtung; hier ist die Zeit der Beleg einer fehlenden Antwort.
