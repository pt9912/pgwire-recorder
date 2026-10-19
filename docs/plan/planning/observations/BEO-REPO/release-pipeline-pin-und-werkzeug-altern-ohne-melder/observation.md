# Pin und Werkzeug einer Release-Pipeline altern, ohne dass etwas es meldet

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Schritt der Release-Pipeline hängt an einem Pin (Image des Builders per Digest) und an Werkzeugen, die das Runner-Image mitbringt (`skopeo`, `gh`, `buildx`). Der Pin wird bei der Pflege gehoben, die Werkzeuge ändern sich mit dem Runner-Image; beides kann Verhalten oder Meldungen ändern, auf die der Schritt sich stützt (Wortlaute von „Tag fehlt“, Form des Layouts). Kein Gate und kein Nachtlauf meldet die Alterung: `slice-harness-ci-freshness-images` prüft die Pins in `Dockerfile` und `harness/mk/*.mk` und nimmt die Pins in Workflows ausdrücklich nicht an. Erst ein Lauf der Pipeline, im Zweifel der Lauf eines echten Tags, zeigt die Änderung.

Abgrenzung zu `BEO-REPO/workflow-pin-pruefung-am-gepinnten-d-check`: Dort geht es um die Prüfung der Pin-Form durch das Doku-Gate; hier um die Aktualität des Pins und die Version der Werkzeuge im Runner.
