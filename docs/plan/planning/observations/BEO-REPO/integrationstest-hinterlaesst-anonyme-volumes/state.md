**Stand:** verkörpert

Dritter Beleg mit `slice-harness-integration-volumes` (Beleg `evidence/slice-harness-integration-volumes.md`).

Zielort: `tools/test/integration-aufraeumen.sh` (Lesen der Namen, Prüfung nach Namen, Entfernen des Rests), `tools/test/run-integration-tests.sh` und `tools/harness/run-container.sh` (`docker rm -fv`), `harness/vertraege/SPEC-057-integration-aufraeumen.md` und `harness/vertraege/SPEC-055-container-und-binaries.md` Punkt 9; Sensor `make integration-volumes-gegenprobe` und `make test-container-gegenprobe` (Fälle `p9-*-volume`) — Herkunfts-Anker `seit slice-harness-integration-volumes`.

Grenze, weiter offen: ein neues `docker run` außerhalb der beiden Läufer ohne `--rm` und ohne `-v` sieht kein Gate (`SPEC-057` Punkt 11); geführt unter `BEO-REPO/docker-run-ohne-rm-ausserhalb-der-laeufer`.

Kein Auftreten, Vermerk (`slice-harness-devcontainer-einrichten`, Urteil des Planners, 2026-10-21): Der Nachweis nach `SPEC-069` Punkt 18 ließ nach `docker rm -fv` zwei **benannte** Volumes und ein Image-Paar des Gerüsts stehen; das sind keine anonymen Volumes, und kein Gate fährt den Nachweis. Der Vertrag nennt jetzt das Entfernen per Name (Punkt 18, Schritt 5), und der Beleg zeigt beide Zustände. Kein Evidence, der Zähler bleibt bei 3.
