# Eine Abnahme ohne PostgreSQL lässt sich mit dem Gate nicht bewusst brechen

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Integrations-Runner fährt erst die Phase gegen PostgreSQL, dann die ohne. Jede
Mutation, die das Replay bricht, macht schon die erste Phase rot, und der Runner
endet dort; welcher Test der zweiten Phase eine DoD-Behauptung fängt, zeigt
`make test-integration` allein nicht.
