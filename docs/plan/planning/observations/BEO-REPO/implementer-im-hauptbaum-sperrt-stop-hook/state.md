**Stand:** offen

Erstes Auftreten: `slice-harness-coverage` (Beleg `evidence/slice-harness-coverage.md`); unter der Schwelle, `offen`. Die Gegenmaßnahme ist als Slice geplant (`slice-harness-devcontainer-u-boot-0-10`, `slice-harness-devcontainer-einrichten`); sie ist keine Zuweisung der Beobachtung, die Ausgänge folgen ab 3× oder `gestrichen` mit Begründung.

Gegenmaßnahme geliefert, Standardweg offen (Urteil des Planners, `slice-harness-devcontainer-u-boot-0-10`, 2026-10-17): Das Gerüst auf `u-boot` 0.10.0 trägt den vollen `make gates` (zwei grüne Läufe, 1701 s, höchstens 5,42 GiB). Dass lang laufende Slices im Container laufen und der Hauptbaum frei bleibt, belegt erst `slice-harness-devcontainer-einrichten`; kein Auftreten, kein Evidence.

Vermerk, kein Auftreten (Urteil des Planners, `slice-harness-devcontainer-einrichten-geruest`, 2026-10-18): Die Arbeit eines Implementers im Dev-Container `pgr-dev2` (Hostsitzung per `ssh` im Klon) trug einen Slice einschließlich Gegenprobe und Gates; der Hauptbaum am Host blieb für den Stop-Hook frei. Der Standardweg bleibt offen, bis `slice-harness-devcontainer-einrichten` ihn automatisiert.
