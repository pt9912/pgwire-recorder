# Die Range-Prüfungen von d-check laufen in einem Git-Worktree nicht

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`make doc-commits` und `make doc-immutable` binden nur das Arbeitsverzeichnis in das Image von d-check ein (`-v "$(CURDIR):/repo:ro"`). In einem Worktree ist `.git` eine Datei, die auf das Verzeichnis des Hauptklons zeigt; das Image meldet `kein lesbares git-Repository unter /repo: lstat /repo/.git/config: not a directory`. Rollen, die in einem Worktree arbeiten (Review, Verifikation, Implementer im Worktree), können die Belege der Range-Prüfungen dort nicht wiederholen und weichen auf einen frischen Klon aus.

Abgrenzung zu `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`: Dort ist der Hauptbaum der Grund für den Worktree; hier trägt der Worktree ein Werkzeug nicht.
