# `play` mit einer Aufzeichnung ohne Sitzung endet ohne Meldung, `replay` meldet `PGR-E3004`

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Aufzeichnung ohne jede Session (`record` ohne Verbindung) ist gültig. `replay` lehnt sie beim Start
ab (`PGR-E3004`, Exit 3). `play` endet mit Exit 0 und den Zeilen „play gestartet“ und „play beendet“,
ohne Warnung. `LH-FA-20.a` *Start* entscheidet das: Enthält die Aufzeichnung keine Session mit Interaktion,
verbindet sich `play` nicht und endet mit Exit-Code `0`; der Architect empfahl, bewusst zu schweigen.
Das Handbuch sagt über `play` bei einer Aufzeichnung ohne Sitzung nichts (Review F-634), und das ist kein Fehler.
