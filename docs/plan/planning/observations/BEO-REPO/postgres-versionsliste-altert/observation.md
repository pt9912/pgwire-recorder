# Die Liste der unterstützten PostgreSQL-Versionen altert ohne Sensor

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Liste der Hauptversionen (`harness/mk/postgres-versionen.mk`, später die Spezifikation) folgt der Entscheidung „fünf jüngste erschienene Hauptversionen“, hebt sich aber nicht selbst an; mit jeder neuen Hauptversion und jedem Ende einer Version fällt sie zurück, und kein Gate merkt es, weil ein Sensor Netz zur Laufzeit bräuchte.
