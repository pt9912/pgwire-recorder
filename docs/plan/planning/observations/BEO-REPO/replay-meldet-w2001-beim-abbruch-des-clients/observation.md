# `replay` meldet beim Beenden `PGR-W2001`, wenn der Client allein abbricht

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Bricht ein Client im `replay` selbst ab (libpq `require_auth` mit einem Verfahren, das `replay` nicht anbietet) und
verbraucht die einzige Sitzung der Aufzeichnung nicht, endet `replay` mit Exit-Code 0 und meldet beim Beenden
`PGR-W2001` (Sitzung nie zugeordnet). Die Sitzungszahl entscheidet nicht, sondern ob die Sitzung verbraucht wurde
(Probe P15 (a), `done/slice-v1-abschluss-protokollrand-doku.md` §7). Das Handbuch nennt den Fall nicht ausdrücklich;
die allgemeinen Sätze zu nie zugeordneten Sitzungen und die Tabellenzeile `PGR-W2001` decken ihn, kein Satz ist falsch.

Abgrenzung zu `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`: Dort fehlt ein Teil, den ein Slice schuldet. Hier schuldet
ihn kein Slice; es ist eine benannte Lücke des Handbuchs ohne Adresse.
