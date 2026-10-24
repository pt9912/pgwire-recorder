# Der Tag darf auf einem Nachfahren des Release-Commits liegen, und `releasing.md` sagt es nicht

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`docs/maintainer/releasing.md` §5 Schritt 3 und §10 nennen als getaggten Commit den, der `VERSION` und die Zeile der Änderungshistorie trägt. Bleibt der Trockenlauf auf diesem Commit rot (der Scan fand einen Befund mit Behebung), kommt der Pin der Behebung als weiterer Commit dazu, und der Tag liegt auf dem Nachfahren, der beides trägt. Die Checkliste, der Eintrag und die Ersatzbelege eines Carveouts nennen dann einen anderen Hash als den „Release-Commit“; erst der Plan des Slice und der Carveout haben es als Korrektur nachgezogen. Die Regel für diesen Fall (der Tag liegt auf dem Commit, auf dem der grüne Trockenlauf stand und der `VERSION` und die Zeile trägt; Belege 1 bis 6 gelten für genau diesen Hash) steht nirgends.
