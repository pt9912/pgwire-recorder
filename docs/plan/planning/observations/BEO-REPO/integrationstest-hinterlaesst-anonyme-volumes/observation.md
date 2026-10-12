# Ein Integrationstest-Lauf hinterlässt anonyme Docker-Volumes

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Gate startet einen Container aus einem Image, das ein `VOLUME` deklariert (hier das
gepinnte PostgreSQL-Image), und räumt ihn ohne die Option auf, die dessen anonymes Volume
mit entfernt. Das Gate bleibt grün, der Container ist weg, das Volume liegt weiter auf dem
Host. Jeder Lauf kostet einige Dutzend MB; erst eine volle Platte macht es sichtbar, und
dann scheitert ein anderes Gate mit einer Meldung, die auf die Ursache nicht zeigt
(„PostgreSQL nicht bereit“).

Abgrenzung zu `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`: Dort fehlt einem Test die
Frist; hier räumt ein Lauf nicht auf, was er angelegt hat.
