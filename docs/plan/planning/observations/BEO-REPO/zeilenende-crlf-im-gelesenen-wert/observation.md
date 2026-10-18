# Ein Gate liest das Zeilenende als Teil des Werts, und der Vertrag nennt es nicht

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Skript liest eine Datei zeilenweise und vergleicht Werte; trägt die Datei `\r\n`, ist das `\r` Teil des Werts, und der Befund täuscht (`permissions-form … gelesen: {}`, der gelesene Wert sieht richtig aus). Das Ergebnis ist falsch-positiv und damit sicher, die Meldung führt aber in die Irre, und der Vertrag sagt nicht, ob das Zeilenende eine Grenze oder ein Befund ist. `AGENTS.md` §3.12 zählt das Zeilenende bei plattformabhängigen Zusagen zu den Randformen vor dem Code; hier hing die Zusage nicht an einer Plattform des Produkts, sondern an der Eingabedatei, die ein Checkout unter Windows mit `autocrlf` erzeugt.

Abgrenzung zu `BEO-REPO/verhalten-nur-unter-linux-geprueft`: Dort fehlt eine Plattform bei einer Eigenschaft des Dateisystems; hier ist die Eingabe eines Gates anders kodiert, als der Vertrag annimmt.
