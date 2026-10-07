# Ein roter Lauf hängt bis zum Zeitlimit des Pakets

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test wartet ohne eigene Frist: auf einen Kanal, auf einen Prozess oder auf ein
Ereignis, das unter einem Fehler nicht eintritt. Ist der Lauf rot, endet er nicht mit
einer Zusicherung, sondern erst an der Zeitgrenze von `go test` (ohne `-timeout` 10
Minuten) mit einem Goroutinen-Dump. Das Rot ist richtig, aber teuer, und ein
Mutationslauf mit einer Zeitgrenze je Mutant zählt es womöglich anders.
