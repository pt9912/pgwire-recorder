**Vorgang:** slice-extended-query-replay
**Fund:** Test-Kommentar und Abdeckungszeile sagten „die Verbindung endet“ zu, geprüft war es nicht (Review F-322); die Kommentare an `replaySitzung` beschrieben nur den einfachen Fall (F-324).
**Fund (Folge-Reviews):** Die Kopplung zur Frist nannte den Anschluss, nicht dessen Einstufung als `PGR-E4006` (F-332); der Port-Kommentar sagte an `Shutdown` einen Zustandswechsel zu, den die Implementierung nicht hat (F-333); Test-Kommentare und Abdeckungszeilen versprachen Werte und Fälle, die der Test nicht vergleicht (F-335, F-340, F-342).
