# Der Dateiname der Review-Reports weicht von der Vorlagen-Form ab

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Vorlage des Review-Reports nennt `<YYYY-MM-DD>-<slice-Kennung>.md`; dieses Repo legt
Reports als `<datum>-<rolle>-<slice-Kennung>.md` ab (`-review-`, `-verifikation-`,
`-mutationen-`, `-validierung-`). Der Wortlaut des Regelwerks („die volle Slice-Kennung im
Dateinamen“) verbietet den Rollen-Infix nicht, und kein Sensor liest den Namen; die Abweichung
ist nur im Plan eines Slice begründet. Geprüft würde sie erst von einem künftigen Slice, der
das d-check-Modul `reviews` aktiviert (`match: name` gegen die Namen); dieser Slice ist nicht
angelegt, der Eintrag ist deshalb keine Zuweisung.
