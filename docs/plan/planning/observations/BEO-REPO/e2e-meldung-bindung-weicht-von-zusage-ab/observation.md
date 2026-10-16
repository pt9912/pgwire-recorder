# Der E2E-Test bindet die Meldung nicht an das, was die Zusage nennt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test am Binary prüft die Meldung eines Fehlers als einzelne Teilstrings des ganzen Ausgabestroms.
Die Bindung weicht auf zwei Seiten von der Zusage ab: zu lose (Teilstrings nicht in derselben Zeile,
SQLSTATE und Text der Fehlerantwort nicht gebunden, ein Mutant, der sie streicht, bleibt grün) oder zu
eng (ein Wortlaut des Codes ist gebunden, den die Spezifikation nicht fixiert, eine Umformulierung
färbt den Test rot, ohne dass eine Zusage bricht).

Abgrenzung zu `BEO-REPO/test-prueft-teilstring-gegen-meldung-mit-zufallspfad`: Dort entscheidet ein
Zufallspfad über einen Fehlalarm; hier ist der Pfad nicht im Spiel. Abgrenzung zu
`BEO-REPO/abwesenheitszusage-an-wortlaut-der-meldung-geprueft`: Dort ist eine Zusage über das
Fehlen eines Werts an einen Wortlaut gebunden; hier sagt die Zusage etwas über das Vorhandene.
