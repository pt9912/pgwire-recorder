# Folgeabweichungen nach einer verschobenen Folge machen die Meldung von `--diff-all` lang und irreführend

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`play --compare-responses --diff-all` ordnet die Positionen der aufgezeichneten und der
beobachteten Folge Position für Position zu, ohne die Folgen abzugleichen. Fehlt eine
Nachricht, weichen alle folgenden Positionen ab; die Meldung nennt bis zu 20 Positionen und
„und N weitere“. Ob die Praxis einen Abgleich der Folgen verlangt, ist offen; wäre es nötig,
ist es ein Folge-Slice, kein Teil der Option.
