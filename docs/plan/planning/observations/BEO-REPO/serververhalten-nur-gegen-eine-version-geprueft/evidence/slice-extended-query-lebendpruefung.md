**Vorgang:** slice-extended-query-lebendpruefung
**Fund:** Die Erkennung der Lebendprüfung zählte `\v` zum Leerraum, wie PostgreSQL 17 es tut; PostgreSQL 14 bis 16 beantworten eine solche Anfrage mit einem Syntaxfehler, und das Replay lieferte gegen eine Aufzeichnung von 16 eine Antwort, die nicht die des Servers war (Review F-350, HIGH, Sonde S1).
**Fund (Verifikation):** Das Gate fährt die Integrationstests nur gegen 17; eine zu niedrige Grenze im Produktcode ist dort Ende zu Ende grün und nur gegen 16 rot (V-33). Adresse: `slice-v1-abschluss-postgres-versionen`.
