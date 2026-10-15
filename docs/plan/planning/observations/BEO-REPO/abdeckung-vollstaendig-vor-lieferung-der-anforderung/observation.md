# Die Abdeckungstabelle nennt eine Anforderung „vollständig“, bevor sie geliefert ist

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`make abdeckung` schreibt `vollständig` und nimmt die Anforderung in `abdeckung-vollstaendig.md`
auf, sobald Happy, Boundary und Negative gedeckt sind, gleich aus welcher Nachweisart
(`SPEC-048`). Liefert eine Anforderung mehrere Slices, steht sie nach dem ersten als
`vollständig`, obwohl Aufzeichnen, Wiedergabe, Einspielen und SQLite fehlen; `make doc-trace`
liest die Datei. Die Tabelle sagt etwas über die Pfade der Deklarationen, nicht über die Lieferung.
