# Eine eigene harte Regel ändert einen Satz einer angenommenen ADR, ohne den Nachzug zu nennen

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Änderung von `AGENTS.md` (Rang 8) ändert, was eine angenommene ADR (Rang 4) in einem Punkt festlegt, und weder die Regel noch der Plan nennt, dass der Punkt sich ändert und was mit der ADR geschieht (kein Folge-ADR, ein Satz in der Regel, Folge-ADR mit `Supersedes`). Der Leser der ADR liest weiter den alten Punkt; nach `AGENTS.md` §3.5 darf die ADR nicht überschrieben werden, und die Rangfolge der Quellen verlangt, dass die höher rangierende gewinnt, der Widerspruch aber benannt ist.

Abgrenzung zu `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`: Dort ersetzt eine Regel eine andere. Hier bleibt die Regel stehen, und ihr Inhalt widerspricht dem Punkt einer ADR.
