# Ein DoD-Beleg auf einer fremden Plattform nennt ein Belegverfahren, das ihn nicht erzeugen kann

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die DoD eines Slice verlangt einen Beleg von einer Plattform, die der Implementer nicht erreicht (Lauf auf GitHub), und nennt im selben Punkt das Verfahren, mit dem er entstehen soll. Das Verfahren erzeugt den Beleg nicht: Ein Trockenlauf ist grün gebaut und kann den Wiederanstoß eines roten Laufs nicht zeigen. Das merkt erst die Verifikation, und der Slice ist bis dahin nicht schließbar, ohne dass der Beleg entfallen oder weitergegeben ist.

Abgrenzung zu `BEO-REPO/beleg-traegt-zusage-aus-anderem-grund`: Dort trägt ein vorhandener Beleg die Zusage nicht; hier gibt es das Verfahren, aber nicht den Beleg, und die DoD nennt ihn trotzdem.
