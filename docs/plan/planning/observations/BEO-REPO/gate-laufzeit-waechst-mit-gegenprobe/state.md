**Stand:** verkörpert

Drittes Auftreten mit `slice-harness-gate-kette-schnell` (Beleg `evidence/slice-harness-gate-kette-schnell.md`, F-702): Die Kette, die den Lauf kürzt, bringt ihre Gegenprobe selbst in `make gates` mit. Der Hook liest den Stand der schnellen Kette seit `slice-harness-gate-kette-stop-hook-freigabe`.

Zielort: `.claude/hooks/stop-require-gates.sh` (Freigabe über die Positivliste und den Stand der schnellen Kette), die schnelle Kette `make gates-schnell` und ihr Vertrag `SPEC-052` — Herkunfts-Anker `seit slice-harness-gate-kette-stop-hook-freigabe` (Beleg `evidence/slice-harness-gate-kette-stop-hook-freigabe.md`). Zuweisung des Planners bei der Closure: Die Messung zeigt die Senkung des Weges (ein Stop nach einem Doku-Commit endet nach 8 s und einem Hook-Aufruf von 4,5 s statt nach 7 min 32 s), nicht der Laufzeit von `make gates`, die mit `stop-hook-gegenprobe` (2 min 9 s) weiter stieg. Der Retirement-Check lautet: Verlangt der Hook für einen inerten Stop wieder den vollen Lauf, oder verlängert ein neues Gate die schnelle Kette über das Maß, in dem sie einem Doku-Commit zumutbar ist?

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-v1-abschluss-postgres-versionen-matrix` — rund 195 s Zuwachs an `make gates` durch das Ziel und zwei Gegenproben; die schnelle Kette bleibt davon ausgenommen (Beleg `evidence/slice-v1-abschluss-postgres-versionen-matrix.md`).

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-v1-abschluss-container` — `make gates` 927 s am Endstand, davon Gegenprobe 3 m 40 s; sie bleibt aus der schnellen Kette (Beleg `evidence/slice-v1-abschluss-container.md`).

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-integration-volumes` — `make gates` 1038 s am Endstand, die Gegenprobe des Aufräumens 42 s; sie bleibt aus der schnellen Kette (Beleg `evidence/slice-harness-integration-volumes.md`).

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-coverage` — die Gegenprobe 89 s; sie bleibt aus der schnellen Kette (Beleg `evidence/slice-harness-coverage.md`).

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-devcontainer-einrichten-geruest` — die Gegenprobe des Container-Starts 4 min 18 s; sie hängt nicht an `make gates` (Host-Ziel, [ADR-0047](../../../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md)), die Laufzeit der Kette wächst nicht (Beleg `evidence/slice-harness-devcontainer-einrichten-geruest.md`).

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-gegenproben-aus-gates` (Plan, Beleg `evidence/slice-harness-gegenproben-aus-gates.md`): Der Nutzer nimmt die Gegenproben und `a-check-negativ` aus `make gates` und aus der Kette ganz, statt sie zu kürzen; die Gegenprobe der schnellen Kette und die Ausnahmeliste schrumpfen mit. Der Zuschnitt „nur bei geänderten Eingaben“ ist als möglicher späterer Slice vermerkt, nicht angelegt.

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-gegenproben-aus-gates` (Beleg `evidence/slice-harness-gegenproben-aus-gates.md`, Messung bei der Closure): Die Gegenproben und `a-check-negativ` verlassen `make gates`, 1630 s auf 260 s im Dev-Container; die Laufzeit der Kette wächst mit einer neuen Gegenprobe nicht mehr, die Fahrt der Gegenproben hängt an der Pflicht (`BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`). Retirement-Check unverändert; zusätzlich: wächst `make gates` erneut über das Maß, in dem der Lauf vor dem Handoff nicht übersprungen wird?

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-gegenproben-aus-gates-sammelziele` (Beleg `evidence/slice-harness-gegenproben-aus-gates-sammelziele.md`, Messung der Fahrten): `make gates` 267 s, `make gegenproben` im Container 1095 s, `make gegenproben-host` 255 s, `make gegenproben-gegenprobe` 42 s. Retirement-Check unverändert; die Gegenprobe der Sammelziele hängt nicht an `GATE_CHECKS`.
