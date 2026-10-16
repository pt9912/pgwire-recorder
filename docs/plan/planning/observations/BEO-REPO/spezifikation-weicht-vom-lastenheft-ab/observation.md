# Die Spezifikation sagt zu derselben Lage etwas anderes als das Lastenheft

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Aussage der Spezifikation (Rang 2) widerspricht dem Wortlaut einer Anforderung im Lastenheft
(Rang 1): der Code folgt der Spezifikation, das Lastenheft nennt für dieselbe Lage ein anderes
Ergebnis. Ein Test prüft nur die Spezifikation; erst ein Review oder eine Verifikation, die beide
Texte nebeneinander legt, findet es. Welche Seite gilt, ist keine Entscheidung des Agenten: Das
Lastenheft rangiert höher, und eine Abweichung davon entscheidet der Nutzer.

Abgrenzung zu `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut`: Dort widersprechen sich zwei
Stellen derselben Schicht, und der Architect gleicht sie an; hier liegen die Stellen in zwei Straten,
und die Angleichung nach unten ist nicht seine Entscheidung. Abgrenzung zu
`BEO-REPO/spec-randform-erst-im-review-entschieden`: Dort ist eine Randform unentschieden; hier ist sie
entschieden, in den Straten verschieden.
