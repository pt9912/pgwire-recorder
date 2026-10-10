# Ein asynchron wirkender Aufruf hat im Runner keinen Wartepunkt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Runner ruft etwas auf, das asynchron wirkt (`pg_reload_conf()` für die `pg_hba.conf`), und
startet den Testlauf, ohne auf die Wirkung zu warten. Bis dahin gilt der alte Stand; ein Test,
der die neue Wirkung erwartet, schlüge in diesem Fenster fehl. Nicht aufgetreten, solange der
Lauf nach dem Aufruf einen neuen Container startet.
