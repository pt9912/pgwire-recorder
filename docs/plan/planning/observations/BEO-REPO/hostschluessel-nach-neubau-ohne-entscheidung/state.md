**Stand:** offen

Zweites Auftreten: `slice-harness-devcontainer-u-boot-0-10` (Beleg `evidence/slice-harness-devcontainer-u-boot-0-10.md`). Das Gerüst auf `u-boot` 0.10.0 ändert den Hostschlüssel bei jedem Neubau des Images, auch bei einer Änderung von `egress.allow`; Neustart und neuer Container tun es nicht. Zähler 2, unter der Schwelle, `offen`; der Nehmer `slice-harness-devcontainer-einrichten` trägt den Befund in Randform 9.

Risiko getragen, kein neues Auftreten (Urteil des Planners, `slice-harness-devcontainer-einrichten-start`, 2026-10-18): Der Slice trug das Risiko A9 (Hostschlüssel nach Neubau) als Randform; die Probe von `container-start` setzt `StrictHostKeyChecking=no` und `UserKnownHostsFile=/dev/null` nur für sich und schreibt nie `known_hosts` (Fall `p11-known-hosts-unveraendert`). Ein Neubau mit anderem Hostschlüssel lief in diesem Slice nicht; kein Evidence, Zähler bleibt bei 2, `offen`.
