# Ein Pfad mit Prosa steht nicht auf der Positivliste des Stop-Hooks

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Positivliste des Stop-Hooks (`inert` in `.claude/hooks/stop-require-gates.sh`, Vertrag `SPEC-052` Punkt 12) nennt `docs/plan/planning/`, `docs/reviews/`, `docs/user/`, `docs/developer/` und die Wurzel-`README.md`. Ein Pfad, der Prosa trägt und den kein Gate der Ausnahmeliste liest, aber nicht auf der Liste steht, verlangt für jede Änderung einen vollen `make gates` vor dem Push. Beobachtet an `docs/developer/` (vier volle Läufe für ein Dokument), jetzt aufgenommen; andere Pfade (`docs/maintainer/`, `harness/README.md`) stehen nicht darauf, ihre Aufnahme verlangt je einen Beleg nach `SPEC-052` Punkt 13 und eine Entscheidung des Nutzers.
