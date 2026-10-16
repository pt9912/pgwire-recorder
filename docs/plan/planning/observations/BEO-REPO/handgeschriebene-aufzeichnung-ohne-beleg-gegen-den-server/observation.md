# Handgeschriebene Aufzeichnungen bilden ab, was die Spezifikation sagt, nicht was der Server liefert

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Tests der Slices zum Verbindungsende stützen sich auf handgeschriebene Aufzeichnungen und Fakes,
bevor `record` ein Ende erzeugen kann. Sie belegen, dass der Recorder die Spezifikation einhält; ob
PostgreSQL dieselbe Nachrichtenfolge liefert, belegt erst der E2E-Test des Abnahmeszenarios am Ende der
Reihe.

Abgrenzung zu `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`: Dort ist das Verhalten gegen
eine Version geprüft, hier gegen keinen Server.
