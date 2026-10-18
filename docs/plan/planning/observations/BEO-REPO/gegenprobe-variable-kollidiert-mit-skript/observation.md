# Eine neue Variable im Kopf einer Gegenprobe trägt den Namen einer Variable des Skripts

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Gegenproben sind lange bash-Skripte mit vielen Variablen im Hauptteil. Eine neue Variable im Kopf, deren Name das Skript schon führt, überschreibt dessen Wert: In `postgres-versionen-gegenprobe.sh` und `container-gegenprobe.sh` stand `liste` im Kopf, das Skript benutzt `liste` selbst, und 15 Mutationen mit unbeteiligten Fällen wurden rot. Gefunden hat es der Implementer an dem unerwarteten Rot; behoben durch die Umbenennung in `gp_liste`. Es gibt keine Prüfung auf doppelte Namen (kein Sensor, `bash -n` und `shellcheck` melden sie nicht).
