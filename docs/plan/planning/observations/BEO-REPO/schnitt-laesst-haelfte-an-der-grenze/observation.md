# Ein Schnitt lässt eine Hälfte an der Grenze einer Review-Sitzung stehen

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Rückführung *zu groß* schneidet einen Slice in zwei. Für eine der beiden Hälften
liegt die Schätzung schon beim Schnitt an der Grenze einer Review-Sitzung; statt sie gleich
weiter zu schneiden, nennt der neue Plan einen zweiten Schnitt als vorab benannte
Rückführung. Diese tritt nach den ersten Liefer-Commits ein, und der Slice geht ein zweites
Mal zurück zur Zerlegung.

Abgrenzung zu `BEO-REPO/slice-waechst-durch-uebernahmen`: Dort wächst ein Slice nach seiner
Anlage durch Übernahmen. Hier war die Größe beim Schnitt bekannt, und der Schnitt hat sie
nicht aufgelöst.
