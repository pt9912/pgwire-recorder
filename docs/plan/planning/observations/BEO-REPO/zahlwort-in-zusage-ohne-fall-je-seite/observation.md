# Ein Zahlwort in einer Zusage hat nicht auf jeder Seite einen Fall

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Zusage mit Zahlwort („genau eine“, „jede andere Zahl“, „alle Zeilen“, „nur“) nennt
mehrere Seiten: weniger, genau so viele, mehr. Geprüft wird eine Seite, die andere hat
keinen Fall, und der Mutant, der die Bedingung auf diese Seite verengt (`n != 1` zu
`n == 0`), bleibt grün. Das Zahlwort ist keine Aufzählung, die `AGENTS.md` §3.10 je Element
zählt; es zählt keine Regel je Seite.

Abgrenzung zu `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: Dort fehlt ein Element der
Aufzählung oder ein Schritt der Ortsangabe. Hier steht ein Zahlwort, dessen Seiten kein
Element der Zusage einzeln nennt. Abgrenzung zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`:
Dort ist der Satz weiter als die Prüfung im Allgemeinen; dieser Eintrag hält die Form fest,
in der der Satz ein Zahlwort trägt.
