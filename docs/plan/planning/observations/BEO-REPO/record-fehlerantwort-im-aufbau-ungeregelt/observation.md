# Fehlerantwort des Servers im Verbindungsaufbau ist für record ungeregelt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Beendet der Upstream den Verbindungsaufbau mit einer Fehlerantwort (etwa eine
fehlende Datenbank), gibt `record` sie an den Client weiter, zählt sie aber
nicht als Verbindungsfehler nach LH-FA-13.b; die Spezifikation regelt den Fall
für `record` nicht.
