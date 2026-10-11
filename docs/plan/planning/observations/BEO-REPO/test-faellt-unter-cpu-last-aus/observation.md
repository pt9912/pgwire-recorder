# Ein Test fällt unter hoher CPU-Last aus

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test, der nicht auf Zeit prüft und eine Frist von 30 s trägt, fällt aus, wenn der Rechner stark
belastet ist: Die erwartete Nachricht kommt nicht binnen der Frist. Das Rot hat keine Ursache im
Diff des Slice, an dem es auffällt; es kann eine Eigenschaft des Tests (Frist, Reihenfolge) oder
eine Wettlaufbedingung im Code sein, und beides ist nicht untersucht.

Abgrenzung zu `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`: Dort fehlt die Frist; hier ist sie da
und wird unter Last gerissen. Abgrenzung zu `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients`:
Dort entscheidet die Pause des Clients über den Ablauf der Aufzeichnung; hier wartet ein Test auf
den Prozess.
