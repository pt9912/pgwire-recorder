# Ein Vertrag nennt den Exit-Code des Skripts, das Ziel läuft über `make` und gibt einen anderen

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Vertrag legt für ein Make-Ziel Exit-Codes fest (1 bei rotem Einzelziel, 2 bei `-i` und bei leerer Liste), die Schnittstelle ist das Ziel, und der Code liegt im Skript, das das Rezept ruft. `make` gibt jedem scheiternden Rezept Exit 2 und überliest ein rotes Rezept unter `-i` mit Exit 0. Der Vertrag stimmte dann für das Skript und nicht für das, was ein Aufrufer über `make` sieht; Implementer und Review fanden es erst am Code (`SPEC-060` Punkte 6, 9 und 10 vor dem Absatz zu den Exit-Codes).

Abgrenzung zu `BEO-REPO/spec-randform-erst-im-review-entschieden`: Dort fehlt eine Randform im Vertrag; hier steht sie, an einer Schnittstelle, die der Aufrufer nicht sieht.
