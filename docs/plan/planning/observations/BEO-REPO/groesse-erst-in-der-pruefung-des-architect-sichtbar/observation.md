# Die Größe eines Slice zeigt sich erst in der Schätzung des Architect, nicht in der Zählung der Liefer-Punkte

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Plan zählt höchstens drei Liefer-Punkte und höchstens zwei Schichten und hält die
Größenregel damit ein. Erst die Prüfung des Architect vor dem Code schätzt den Diff auf
2200 bis 3000 Zeilen mit Tests, über einer Review-Sitzung, und der Slice geht zurück zur
Zerlegung. Die Zählung misst Punkte, nicht Umfang; ein Punkt wie *Roundtrip* oder *Kern*
trägt beliebig viel.

Abgrenzung zu `BEO-REPO/slice-waechst-durch-uebernahmen`: Dort wächst ein Slice nach seiner
Anlage durch Übernahmen. Hier war der Plan beim Anlegen so groß, wie er blieb; die Schätzung
in Zeilen fehlte.
