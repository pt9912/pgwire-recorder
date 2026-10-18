# Die Prüfung der Pins hängt am gepinnten d-check und liegt außerhalb des Gates der Workflows

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Das Gate der Workflows (`make workflows-check`) prüft die Pins der `uses:`-Referenzen nicht selbst; das Modul `workflows` von d-check prüft sie in `make docs-check`, und das Gate hält nur, dass `.d-check.yml` das Modul führt. Was das Modul tut, hängt an der Version des gepinnten d-check: Es meldet ohne Datei und ohne `uses:` einen Befund (leere Prüfmenge), der Eintrag konnte deshalb nicht vor der ersten Workflow-Datei stehen, wie der erste Plan annahm. Ob das Modul `.yaml`-Dateien und die Blockform von `uses:` so liest, wie `SPEC-063` Punkt 8 es voraussetzt, ist nur an der echten Datei gezeigt (ein Fall der Gegenprobe gegen das echte Image deckt `.yml` mit Tag-Kommentar, fehlender Version und fehlendem Kommentar).

Abgrenzung zu `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`: Dort wertet ein Werkzeug eine Regel seiner Konfiguration anders aus; hier steht die Prüfung in einem anderen Gate als dem, das sie verspricht, und ihre Wirkung hängt an der Version des Moduls.
