**Vorgang:** slice-extended-query-lebendpruefung
**Fund:** Der Kommentar an `leerraum` sagte „die Zeichen, die der Scanner von PostgreSQL als Leerraum liest“ ohne Serverversion zu; geprüft war nur gegen PostgreSQL 17, und `\v` ist darunter kein Leerraum (Review F-350).
**Fund (Folge-Review):** Abdeckungs-Deklaration („das Replay folgt der Version in der Aufzeichnung“) und Folge-Slice `slice-v1-abschluss-postgres-versionen` sagten für die Versionsgrenze am Replay mehr zu, als der Integrationstest prüfte; er prüfte sie nur von oben (F-360, Mutation IM1 grün).
