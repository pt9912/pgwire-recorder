# Die Schätzung des Diffs in Zeilen liegt unter der Messung, vor allem bei den Tests

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Architect schätzt vor dem Code den Umfang eines Slice in Zeilen mit Tests. Gemessen liegt der Diff
um die Hälfte und mehr darüber, weil die Tests des neuen Vertrags aus Tabellen mit rohem SQL bestehen und
Nacharbeiten zu Entscheidungen, Review und Verifikation hinzukommen. Die Zählung von Liefer-Punkten und
Schichten bleibt eingehalten; die Schwelle der Rückführung im Plan wird überschritten, ohne dass die
Schätzung es angekündigt hätte.

Abgrenzung zu `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`: Dort zeigt erst die
Schätzung des Architect, dass der Plan zu groß ist, obwohl die Zählung es nicht zeigte. Hier lag die
Schätzung vor dem Code und war zu niedrig.
