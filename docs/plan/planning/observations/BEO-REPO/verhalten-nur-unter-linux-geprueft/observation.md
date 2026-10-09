# Verhalten des Dateisystems nur unter Linux geprüft

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Das Produkt sagt eine Eigenschaft des Dateisystems zu, die von der Plattform abhängt (etwa
das atomare Ersetzen einer Datei); geliefert wird es auch für macOS und Windows, Gate und
Tests laufen aber nur unter Linux in einem Dateisystem, und die Spezifikation sagt nicht,
was auf einer Plattform gilt, die die Eigenschaft nicht bietet.
