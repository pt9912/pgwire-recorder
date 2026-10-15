# Ein Test prüft eine Meldung gegen Werte, die ein Zufallspfad enthalten kann

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test lehnt eine Meldung ab, die einen Wert als Teilstring enthält (hier die Ziffernfolgen
`4709`, `4710`, `4300`, `4900`, `4800`), obwohl die Meldung den Dateipfad nennt und
`t.TempDir()` eine Zufallszahl im Pfad trägt. Die Wirkung ist ein seltener Fehlalarm im Gate
ohne Ursache im Diff; er sah wie ein Lastproblem aus. Abgrenzung zu
`BEO-REPO/test-faellt-unter-cpu-last-aus`: Dort reißt eine Frist unter Last; hier entscheidet der
Pfad, nicht die Zeit.
