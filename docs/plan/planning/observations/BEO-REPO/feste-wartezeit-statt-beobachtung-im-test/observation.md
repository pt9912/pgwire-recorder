# Ein Integrationstest wartet eine feste Zeit, wo er eine Beobachtung abwarten könnte

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test am Binary schläft eine feste Zeit, bevor er den nächsten Schritt setzt (etwa ein Signal), und
setzt voraus, dass der Prozess bis dahin einen Zustand erreicht hat. Unter Last kann der Zustand
später eintreten; dann prüft der Test etwas anderes als vorgesehen, und das Rot hat keine Ursache im Code.
Die Wartezeit ersetzt eine Beobachtung, die der Test haben könnte (eine Abfrage am Upstream, eine Ausgabe).

Abgrenzung zu `BEO-REPO/test-faellt-unter-cpu-last-aus`: Dort fällt ein Test mit Frist unter Last aus;
hier ist der Test noch grün, und das Review findet die feste Zeit als Ursache eines künftigen Rots.
Abgrenzung zu `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`: Dort fehlt eine Frist; hier ist es eine zu kurze feste Zeit.
