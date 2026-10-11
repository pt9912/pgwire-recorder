# `play` mit einer Aufzeichnung ohne Sitzung endet ohne Meldung, `replay` meldet `PGR-E3004`

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Aufzeichnung ohne jede Session (`record` ohne Verbindung) ist gültig. `replay` lehnt sie beim Start
ab (`PGR-E3004`, Exit 3). `play` endet mit Exit 0 und den Zeilen „play gestartet“ und „play beendet“,
ohne Warnung. Die Spezifikation sagt nur, `play` führe Sessions ohne Interaktion nicht aus; `PGR-E3004` ist
dort für `replay` definiert (`LH-FA-03.a`), für `play` steht weder Fehler noch Warnung, noch ausdrücklich
„kein Fehler“. Das ist eine benannte Spezifikationslücke, kein Fehler des Binaries: Das Handbuch sagt über
`play` bei einer Aufzeichnung ohne Sitzung nichts (Review F-634). Wer die Lücke schließt, entscheidet, ob
`play` wie `replay` ablehnt, warnt oder bewusst schweigt, und trägt es in der Spezifikation ein, bevor ein
Slice, der `play` berührt, es liefert.
