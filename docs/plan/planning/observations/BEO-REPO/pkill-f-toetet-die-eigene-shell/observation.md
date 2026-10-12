# Ein `pkill -f` im selben Shell-Befehl tötet die eigene Shell

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Steht das Muster von `pkill -f` in der Befehlszeile desselben Shell-Aufrufs, passt es auf diese Shell selbst; sie endet mit Exit 144, und der Rest des Befehls läuft nicht. Betrifft Aufräumschritte nach einer Messung mit Hintergrundprozess (Sampler).
