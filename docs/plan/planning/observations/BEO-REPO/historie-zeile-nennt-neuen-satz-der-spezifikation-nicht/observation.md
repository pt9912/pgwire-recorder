# Die Änderungshistorie der Spezifikation nennt einen neuen Satz nicht

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Commit fügt der Spezifikation einen Satz hinzu (eine Aussage, die vorher nicht stand), und die Historie-Zeile desselben Tages zur selben Stelle sagt „Berichtigung, keine Aussage geändert“. Die Historie ist dann für die Stelle falsch, und kein Sensor und keine Regel verlangt, dass sie mit jedem Satz in `spec/` wächst; das Review findet es per Lesen.

Abgrenzung zu `BEO-REPO/plan-folgt-korrektur-nicht`: Dort folgt der Slice-Plan dem Code nicht; hier folgt die Historie der Spezifikation dem Satz nicht.
