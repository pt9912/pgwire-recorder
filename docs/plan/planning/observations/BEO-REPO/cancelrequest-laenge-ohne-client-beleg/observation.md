# Die Behandlung einer CancelRequest-Länge, die kein bekannter Client sendet, hat keinen Beleg eines Clients

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Spezifikation führt die Längen 8 bis 11, 13 bis 15 und 269 bis zur Höchstlänge als `PGR-E6001` (Exit-Code 6) mit der Begründung, kein Client sende sie; die Fälle des Unit-Tests pinnen die Behandlung je Länge. Belegt ist die Annahme über Clients nicht: Der Integrationstest nutzt `pgx` mit der Länge 12, die Längen 16 bis 268 liest die Bibliothek, und für die übrigen gibt es weder eine Probe mit einem Client noch einen Fall aus einer Aufzeichnung. Sendet ein Client eine dieser Längen, endet ein Lauf mit Exit-Code 6, wo er den `CancelRequest` wie jeden anderen behandeln könnte.

Abgrenzung zu `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: Die Fälle sind da; was fehlt, ist der Beleg der Annahme über das Verhalten der Clients.
