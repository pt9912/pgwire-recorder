# Ein grüner Mutant wird als äquivalent eingestuft, obwohl eine Eingabe ihn unterscheidet

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Implementer fährt je Zusage eine Mutation, findet einen grünen Mutanten und stuft ihn
in §7 als äquivalent ein, oft mit der Folge, dass die Prüfung als überflüssig entfernt
wird. Die Begründung stützt sich auf die Eingaben, an die er gedacht hat; eine andere
Eingabe unterscheidet Original und Mutant, und erst Review oder Verifikation finden sie
mit einer Sonde.

Abgrenzung zu `BEO-REPO/liste-gruener-mutanten-unvollstaendig`: Dort sind die grünen
Mutanten eines Umbaus im Bestand richtig eingeordnet, aber nicht alle gefunden. Hier ist
ein gefundener Mutant falsch eingeordnet, im neuen Code.
