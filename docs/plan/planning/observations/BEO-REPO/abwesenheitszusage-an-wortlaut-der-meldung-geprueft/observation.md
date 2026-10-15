# Eine Abwesenheitszusage über eine Meldung wird gegen einen Wortlaut geprüft

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Zusage sagt, dass eine Meldung etwas nicht nennt (hier: keine Version). Der Test prüft nur, dass ein
bestimmter Wortlaut fehlt (`Version <Zahl>`); eine andere Form (`Fassung <Zahl>`, `version=<Zahl>`) lässt die
Mutanten grün. Behoben wird mit einer Prüfung der Form des Werts (die Ziffernfolge steht nirgends in der
Meldung) und einem Beleg, der mindestens zwei andere Wortlaute fährt.

Abgrenzung zu `BEO-REPO/test-prueft-teilstring-gegen-meldung-mit-zufallspfad`: Dort entscheidet ein Pfad
mit Zufallszahl über einen Fehlalarm; hier ist der Test zu eng und lässt eine Verletzung der Zusage durch.
Abgrenzung zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: Dort reicht ein Satz weiter als die
Prüfung; hier ist die Prüfung selbst auf einen Wortlaut verengt.
