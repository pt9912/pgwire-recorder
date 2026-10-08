# Ein Slice wächst durch Übernahmen, ohne dass seine Liefer-Punkte mehr werden

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Slice, der für andere Slices die Adresse ihrer Abgrenzungen ist, nimmt Übernahme um
Übernahme in seine vorhandenen Liefer-Punkte auf. Die Zahl der Punkte bleibt bei drei,
ihr Inhalt wächst, bis ein Punkt nicht mehr in eine Review-Sitzung passt oder der Plan
mehr als zwei Schichten berührt; die Größenregel zählt Punkte und meldet das nicht.
