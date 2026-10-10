# Schichten werden je Plan verschieden geteilt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Größenregel zählt höchstens zwei Schichten, die Baseline nennt Schichten aber nur am
Beispiel (Adapter, Service, UI, DB-Schema). Geber und Nehmer einer Sendung zählten im selben
Commit mit verschiedenen Teilungen: Der Geber führte die Spezifikation als eigene Schicht, der
Nehmer fasste Spezifikation, Pläne und Commands zu „Dokumentation“ zusammen und kam so auf
zwei statt vier Schichten. Eine gemeinsame Teilung steht seither nur in §8 der beteiligten
Pläne, nicht an einem Ort, den jeder Plan liest.

Abgrenzung zu `BEO-REPO/slice-waechst-durch-uebernahmen`: Dort wurde nicht nachgezählt oder
die Größe wuchs nach der Anlage. Hier wurde gezählt, aber mit einem Maß, das ein anderer Plan
nicht teilt.
