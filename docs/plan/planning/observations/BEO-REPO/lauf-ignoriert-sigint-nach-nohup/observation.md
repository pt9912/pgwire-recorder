# Ein mit `nohup` gestarteter `make gates` ist rot, weil `SIGINT` ignoriert wird

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Wird `make gates` mit `nohup` im Hintergrund der Shell gestartet, erbt der Lauf die Disposition
„`SIGINT` ignoriert“. Gegenproben, die `SIGINT` an einen Kindprozess schicken und seine Reaktion
erwarten (Fälle `p9-*INT` in `container-gegenprobe`), werden rot, ohne dass der Diff etwas
damit zu tun hat. Derselbe Stand war grün, als der Lauf über das Werkzeug mit
`run_in_background` gestartet wurde.
