# Ein Rollenlauf läuft Stunden, ohne dass eine Marke ihn anhält

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Agentenlauf einer Rolle (hier der Implementer) hat keine Zeitgrenze: Kein Artefakt des Repos (Plan, `AGENTS.md`, `.claude/commands/`, Hook) nennt eine Dauer, bei der er anhält und berichtet, und kein Wächter prüft sie. Der Lauf endet, wenn die Arbeit endet, nach Stunden, und Review und Verifikation sehen erst den fertigen Stand. Gemessen ist die Dauer am Span-Strom des Laufs (`.harness/state/spans/`).

Abgrenzung zu `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`: Dort geht es um die Laufzeit von Gates und Gegenproben (Obergrenze 600 s je Gegenprobe), hier um die Dauer eines Agentenlaufs. Abgrenzung zu `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`: Dort fehlt einem Test die Frist; hier fehlt dem Lauf einer Rolle die Grenze. Abgrenzung zu `BEO-REPO/meldeschwelle-ohne-artefakt`: Dort geht es um die Größe eines Diffs, hier um Zeit.
