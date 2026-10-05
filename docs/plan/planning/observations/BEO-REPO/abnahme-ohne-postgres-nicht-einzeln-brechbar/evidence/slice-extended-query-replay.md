**Vorgang:** slice-extended-query-replay
**Fund (Verifikation, V-30):** Abnahmeszenario 7 (`TestE2EOhnePostgresExtendedReplay`) ließ sich nur über getrennte Phasen von Hand bewusst brechen; der Gate blieb richtig rot, die Zuordnung „welcher Test fängt was“ fehlte. Vorschlag der Verifikation: eine Option des Runners, nur die zweite Phase gegen ein anderes Image zu fahren.
