# E2E-Abdeckung je Anforderung

Erzeugt von `make test-integration` über `tools/test/run-integration-tests.sh`
aus den Zeilen `// Abdeckung: …` direkt über jedem `func TestE2E*` unter
`test/integration/`. Die Datei ändert sich nur, wenn sich eine Deklaration oder
ihr Ort ändert. Sie ist eine Abdeckungs-Deklaration, kein Lauf-Beleg.

| Anforderung | Kurzbeschreibung | Nachweis | Ort |
| --- | --- | --- | --- |
| [`LH-FA-02`](../../spec/lastenheft.md), [`LH-FA-06`](../../spec/lastenheft.md), [`LH-FA-07`](../../spec/lastenheft.md) | ein Client führt `SELECT 1;` über `record` gegen eine reale PostgreSQL-Instanz aus und erhält deren Ergebnis; nach dem Beenden des Laufs steht die Interaktion geordnet in einer Aufzeichnung mit Formatkennung und Version. | `TestE2ERecordSelect1` | `test/integration/record_e2e_test.go:44` |
| [`LH-FA-02`](../../spec/lastenheft.md), [`LH-FA-13`](../../spec/lastenheft.md) | ist der Upstream nicht erreichbar, erhält der Client eine Fehlerantwort mit PGR-E4002; der Lauf geht weiter und endet beim Beenden mit Exit-Code 0 und einer gültigen Aufzeichnung ohne Session. | `TestE2ERecordUpstreamNichtErreichbar` | `test/integration/record_e2e_test.go:101` |
