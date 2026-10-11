# Die Spezifikation nennt `version: !!int 1` nicht, das Binary lehnt es mit `PGR-E3003` ab

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`SPEC-001` *Rangfolge der Prüfung* nennt für `version` die Schreibweisen, die keine kanonische Dezimalzahl
sind (etwa `"1"`, `null`, `01`, `1.0`), nicht aber einen ausdrücklichen Typ am Wert (`!!int 1`). Das Binary lehnt
ihn ab (`PGR-E3003`, der Stil des Skalars ist nicht der einfache); bei `offset_ms` nennt die Spezifikation
`!!int 15` ausdrücklich. Die Randform ist weder zugesagt noch geprüft; ein Handbuch oder eine Probe darf sie
nicht als geprüft führen (Verifikation V-185, INFO). Wer die Lücke schließt, ergänzt die Form in `SPEC-001`
und gibt ihr einen Fall in `TestZeitangabenVersionVorForm`.
