# Ein Handbuchsatz mit „soweit“ trägt mehrere Zweige, und ein Fall deckt nur einen

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Satz des Handbuchs nennt eine Angabe „soweit sie abweicht“ (oder „wenn verschieden“),
und der Test fährt Fälle, in denen mehrere Zweige zusammen die Angabe tragen; ein Mutant, der
einen Zweig allein bricht (hier „Name nur bei gleicher Typ-OID“), bleibt grün. Der Satz ist
eng genug, um nicht falsch zu sein, aber weiter als jeder einzelne Fall. `AGENTS.md` §3.10
verlangt je Element einer Aufzählung in der Spezifikation einen Fall; für einen Satz des
Handbuchs verlangt §3.11 einen Test je Satz, nicht je Zweig.

Abgrenzung zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: Dort ist der Wortlaut
weiter als der Test insgesamt; hier belegt der Test den Satz im Gutfall, und nur ein Zweig
des „soweit“ bleibt ohne Fall.
