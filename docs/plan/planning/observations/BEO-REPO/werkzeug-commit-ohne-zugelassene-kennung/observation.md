# Ein Repo-Werkzeug committet ohne zugelassene Kennung

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Werkzeug des Repos erzeugt Commits, deren Message der Commit-Hook ablehnt,
weil die Kennungsmenge des Hooks und die Kennungen des Repos auseinanderliegen:
Der Hook kannte nur Slice-Kennungen mit Ziffern, das Repo vergibt Namen.

## Für Wellen

`make archive-welle` committet nach `harness/mk/hooks-install.mk` mit dem Welle-Namen; eine Welle-Kennung gehört nicht zur Kennungsmenge des Trägers. Bei der Closure von `welle-walking-skeleton` ist das eingetreten (Beleg unter `evidence/`).
