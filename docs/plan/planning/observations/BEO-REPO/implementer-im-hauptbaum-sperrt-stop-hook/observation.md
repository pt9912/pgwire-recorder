# Ein Implementer-Lauf im Hauptbaum sperrt den Stop-Hook über Stunden

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Stop-Hook am Host (`.claude/hooks/stop-require-gates.sh`) verlangt für Änderungen außerhalb der Positivliste einen vollen `make gates` auf dem Arbeitsbaum. Arbeitet ein Implementer-Lauf im Hauptbaum an einem lang laufenden Slice, ist der Baum über Stunden uneinheitlich, und jede andere Sitzung am Host (Planner, Architect, Nutzer) endet mit gesperrtem Stop oder muss auf den Gate-Lauf des Implementers warten.

## Benannt, nicht gezählt

- 2026-10-17: Der Implementer-Lauf an `slice-harness-coverage` (`Dockerfile`, `harness/`, `tools/` im Hauptbaum) hält den Hook für die Planungsarbeit am selben Tag; die Sperre dauerte über zwei Stunden (Beleg `evidence/slice-harness-coverage.md`). Anlass der Entscheidung des Nutzers vom 2026-10-17, lang laufende Slices im Dev-Container zu fahren (`slice-harness-devcontainer-u-boot-0-10`, `slice-harness-devcontainer-einrichten`).
