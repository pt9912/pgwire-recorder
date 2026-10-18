# Eine Gegenprobe zählt Aufrufe in einem Zeitfenster und wird unter Last grün oder rot

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Fall der Gegenprobe zählt, wie oft eine Attrappe in einem Zeitfenster gerufen wird (die Probe von `container-start` ruft `ssh` jede Sekunde bis `PGR_WARTEN`; die Zusage lautet 60 oder 61 Aufrufe bei 60 Sekunden). Unter Last verschiebt sich das Fenster: Die Vorgabe 61 lieferte 61 statt 62 Aufrufe, und zwei Mutationen (`p11-vorgabe-warten`, `p11-zu-kurz`) blieben grün, obwohl der Fall die Zusage prüfen sollte. Behoben, indem die Probe im Takt von `SECONDS` läuft (61 Aufrufe bei 60, unabhängig von der Last).

Abgrenzung zu `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`: Dort färbt eine Mutation einen Fall außerhalb ihrer Gruppe; hier bleibt eine Mutation grün. Abgrenzung zu `BEO-REPO/test-faellt-unter-cpu-last-aus`: Dort fällt ein Test des Produkts aus; hier zählt ein Fall der Gegenprobe falsch.
