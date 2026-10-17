**Stand:** verkörpert

Drittes Auftreten mit `slice-harness-gate-kette-schnell` (Beleg `evidence/slice-harness-gate-kette-schnell.md`, F-702): Die Kette, die den Lauf kürzt, bringt ihre Gegenprobe selbst in `make gates` mit. Der Hook liest den Stand der schnellen Kette seit `slice-harness-gate-kette-stop-hook-freigabe`.

Zielort: `.claude/hooks/stop-require-gates.sh` (Freigabe über die Positivliste und den Stand der schnellen Kette), die schnelle Kette `make gates-schnell` und ihr Vertrag `SPEC-052` — Herkunfts-Anker `seit slice-harness-gate-kette-stop-hook-freigabe` (Beleg `evidence/slice-harness-gate-kette-stop-hook-freigabe.md`). Zuweisung des Planners bei der Closure: Die Messung zeigt die Senkung des Weges (ein Stop nach einem Doku-Commit endet nach 8 s und einem Hook-Aufruf von 4,5 s statt nach 7 min 32 s), nicht der Laufzeit von `make gates`, die mit `stop-hook-gegenprobe` (2 min 9 s) weiter stieg. Der Retirement-Check lautet: Verlangt der Hook für einen inerten Stop wieder den vollen Lauf, oder verlängert ein neues Gate die schnelle Kette über das Maß, in dem sie einem Doku-Commit zumutbar ist?

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-v1-abschluss-postgres-versionen-matrix` — rund 195 s Zuwachs an `make gates` durch das Ziel und zwei Gegenproben; die schnelle Kette bleibt davon ausgenommen (Beleg `evidence/slice-v1-abschluss-postgres-versionen-matrix.md`).

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-v1-abschluss-container` — `make gates` 927 s am Endstand, davon Gegenprobe 3 m 40 s; sie bleibt aus der schnellen Kette (Beleg `evidence/slice-v1-abschluss-container.md`).

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-integration-volumes` — `make gates` 1038 s am Endstand, die Gegenprobe des Aufräumens 42 s; sie bleibt aus der schnellen Kette (Beleg `evidence/slice-harness-integration-volumes.md`).

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-coverage` — die Gegenprobe 89 s; sie bleibt aus der schnellen Kette (Beleg `evidence/slice-harness-coverage.md`).
