# Ein Test schreibt an eine Verbindung, die der Server schon geschlossen hat

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test sendet eine Nachricht mit großem Rest (9992 Bytes), die der Server aus dem Kopf ablehnt und danach schließt. Schreibt der Test noch, wenn der Server schließt, endet das Schreiben mit `connection reset by peer`, und der Fall wird rot, ohne dass der Code falsch ist; das Rot ist selten (ein Lauf in dreimal 100 Läufen der Fälle). Die Ursache liegt im Test, der das Ende der Schreibseite nicht als zulässigen Ausgang kannte.

Abgrenzung zu `BEO-REPO/test-faellt-unter-cpu-last-aus`: Dort fällt ein Test mit Frist unter Last aus; hier hängt das Rot an der Reihenfolge von Schreiben und Schließen, nicht an einer gerissenen Frist.
