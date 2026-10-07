# Dass die Session aus `Open` den Lesepuffer des Aufbaus trägt, prüft kein Test

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`Open` in `internal/adapters/driven/postgres/upstream.go` liest den Verbindungsaufbau über
das Frontend einer Session und gibt an ReadyForQuery diese Session zurück. Gibt `Open`
stattdessen eine neue Session mit neuem Frontend zurück, gehen Bytes verloren, die der
Server mit dem ReadyForQuery geschickt hat und die schon im Lesepuffer des Aufbaus liegen.
Diese Mutation ändert das Verhalten und bleibt in allen Tests von `postgres` grün.
