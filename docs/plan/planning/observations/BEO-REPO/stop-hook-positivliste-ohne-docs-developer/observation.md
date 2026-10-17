# Der Stop-Hook kennt `docs/developer/` nicht auf seiner Positivliste

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`docs/developer/` steht nicht auf der Positivliste des Stop-Hooks (`docs/plan/planning/`, `docs/reviews/`, `docs/user/`,
Wurzel-`README.md`). Jede Änderung unter `docs/developer/` verlangt deshalb einen vollen `make gates` vor dem Push, auch
wenn sie nur Prosa ist. Die Positivliste zu ändern ist eine Regeländerung (Entscheidung des Nutzers).
